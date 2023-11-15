// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"github.com/ddatdt12/kapo-play-ws-server/src/dto"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
)

// Client is a middleman between the websocket connection and the hub.
type GameSocket struct {
	Info *models.Game
	// Buffered channel of outbound messages.
	Send chan dto.MessageTransfer
	// Registered clients.
	ClientSet map[*Client]bool
}

func NewGameSocket(gameInfo *models.Game) *GameSocket {
	return &GameSocket{
		Info:      gameInfo,
		Send:      make(chan dto.MessageTransfer, 256),
		ClientSet: make(map[*Client]bool),
	}
}

func (game *GameSocket) Run() {
	for message := range game.Send {
		for client := range game.ClientSet {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(game.ClientSet, client)
			}
		}
	}
}
