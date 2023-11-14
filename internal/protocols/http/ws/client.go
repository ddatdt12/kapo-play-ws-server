// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"encoding/json"
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

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub  *Hub
	Room *models.Room
	User *models.User

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan models.Message
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, rawMessage, err := c.Conn.ReadMessage()
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

			// NOTE:save message to database

			c.Hub.Messages <- *message
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Error().Msg("The hub closed the channel")
				return
			}

			err := c.Conn.WriteJSON(dto.NewMessageTransfer(dto.NewMessage, message, nil))
			if err != nil {
				log.Error().Stack().Err(err).Msg("error")
				return
			}

			// Add queued chat messages to the current websocket message.
			for i := 0; i < len(c.Send); i++ {
				queueMessage := <-c.Send
				c.Conn.WriteJSON(dto.NewMessageTransfer(dto.NewMessage, queueMessage, nil))
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
