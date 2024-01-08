// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"fmt"

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
	// Host
	Host *Client

	GameState *models.GameState
}

func NewGameSocket(gameInfo *models.Game, gameState *models.GameState) *GameSocket {
	if gameState == nil {
		gameState = models.NewGameState()
	}

	return &GameSocket{
		Info:      gameInfo,
		Send:      make(chan dto.MessageTransfer, 256),
		ClientSet: make(map[*Client]bool),
		GameState: gameState,
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

func (game *GameSocket) Cancel() {
	for client := range game.ClientSet {
		close(client.Send)
		delete(game.ClientSet, client)
	}
	close(game.Send)
}

func (game *GameSocket) AddClient(client *Client) {
	game.ClientSet[client] = true
}

func (game *GameSocket) RemoveClient(client *Client) {
	if _, ok := game.ClientSet[client]; ok {
		client.CleanUp()
		delete(game.ClientSet, client)
	}
}

func (game *GameSocket) GetHost() *Client {
	return game.Host
}

func (game *GameSocket) SetHost(client *Client) {
	game.Host = client
}

func (game *GameSocket) GetUsers() []*models.User {
	users := make([]*models.User, 0)
	for client := range game.ClientSet {
		users = append(users, client.User)
	}
	return users
}

func (game *GameSocket) NotifyMembers(message dto.MessageTransfer) {
	game.Send <- message
}

func (game *GameSocket) NotifyTo(client *Client, message dto.MessageTransfer) {
	client.Send <- message
}

func (game *GameSocket) NotifyHost(message dto.MessageTransfer) {
	game.Host.Send <- message
}

func (game *GameSocket) NotifyAll(message dto.MessageTransfer) {
	game.Send <- message
	if game.Host != nil {
		game.Host.Send <- message
	}
}

func (game *GameSocket) NotifyUpdatedListPlayers() {
	message := dto.MessageTransfer{
		Type: dto.MessagePlayerList,
		Data: game.GetUsers(),
	}
	game.NotifyAll(message)
}

func (c *Client) String() string {
	return fmt.Sprintf("Client %v", *c.User)
}
