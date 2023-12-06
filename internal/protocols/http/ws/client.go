// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/src/dto"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub        *Hub
	GameSocket *GameSocket
	Game       *models.Game
	User       *models.User
	IsHost     bool

	// The websocket connection.
	Conn *websocket.Conn

	// Context
	ConnectionCtx *ConnectionContext

	// Buffered channel of outbound messages.
	Send chan dto.MessageTransfer

	//Information related to game
	QuestionAnswersMap map[uint]*models.Answer
}

func NewClient(hub *Hub, conn *websocket.Conn, gameSocket *GameSocket, game *models.Game, user *models.User) *Client {
	return &Client{
		Hub:                hub,
		Conn:               conn,
		ConnectionCtx:      NewConnectionContext(),
		GameSocket:         gameSocket,
		Game:               game,
		User:               user,
		Send:               make(chan dto.MessageTransfer, 256),
		QuestionAnswersMap: map[uint]*models.Answer{},
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, rawMessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Stack().Err(err).Msg("websocket close error")
			}
			break
		}
		rawMessage = bytes.TrimSpace(bytes.Replace(rawMessage, newline, space, -1))

		// parse message byte to json
		var messageObj dto.MessageTransfer
		err = json.Unmarshal(rawMessage, &messageObj)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error")
			continue
		}

		log.Info().Interface("Comming message", messageObj).Msg("messageObj")
		router(c, &messageObj)
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
				log.Error().Msgf("The hub closed the channel of client: %s", c.User.Username)
				return
			}
			log.Info().Interface("WRITE message: ", message).Msg("message")
			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Error().Stack().Err(err).Msg("error")
				return
			}

			// Add queued chat messages to the current websocket message.
			for i := 0; i < len(c.Send); i++ {
				queueMessage := <-c.Send
				c.Conn.WriteJSON(queueMessage)
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) CleanUp() {
	log.Info().Msgf("CleanUp client: %s", c.User.Username)
	close(c.Send)
	c.Conn.Close()
	c.ConnectionCtx.Cancel()
}

func (c *Client) Notify(message dto.MessageTransfer) {
	log.Info().Interface("message zxczx", message).Msg("message")
	c.Send <- message
}
