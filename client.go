// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/dto"
	"github.com/ddatdt12/kapo-play-ws-server/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub  *Hub
	room *models.Room
	user *models.User

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan models.Message
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Stack().Err(err).Msg("error")
			}
			break
		}
		rawMessage = bytes.TrimSpace(bytes.Replace(rawMessage, newline, space, -1))
		log.Info().Msg("new rawMessage: " + string(rawMessage))

		// parse message byte to json
		var messageObj dto.MessageTransfer
		err = json.Unmarshal(rawMessage, &messageObj)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error")
			continue
		}

		if !dto.VerifyMessageType(messageObj.Type) {
			log.Error().Any("messageObj.Type", messageObj.Type).Msg("Error message type")
			continue
		}

		log.Info().Interface("messageObj", messageObj).Msg("messageObj")
		if messageObj.Type == dto.SendMessage {
			log.Info().Interface("messageObj.Type", messageObj.Type).Msg("messageObj.Type")
			message := models.NewMessage(messageObj.Data)

			// save message to database

			c.hub.messages <- *message
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Error().Msg("The hub closed the channel")
				return
			}

			err := c.conn.WriteJSON(dto.NewMessageTransfer(dto.NewMessage, message, nil))
			if err != nil {
				log.Error().Stack().Err(err).Msg("error")
				return
			}

			// Add queued chat messages to the current websocket message.
			for i := 0; i < len(c.send); i++ {
				queueMessage := <-c.send
				c.conn.WriteJSON(dto.NewMessageTransfer(dto.NewMessage, queueMessage, nil))
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
