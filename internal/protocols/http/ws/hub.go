// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/json"
	"net/http"

	"github.com/ddatdt12/kapo-play-ws-server/src/dto"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/services"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var messagesStore = make([]dto.MessageTransfer, 0)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Hub maintains the set of active clients and broadcasts Messages to the
// clients.
type Hub struct {
	// Games and GameSockets
	GameSocketMap map[string]*GameSocket

	// Inbound Messages from the clients.
	Broadcast chan []byte

	// Inbound Messages from the clients.
	Messages chan dto.MessageTransfer

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// game service
	GameService services.IGameService

	// user service
	UserService services.IUserService
}

func NewHub(gameService services.IGameService, userService services.IUserService) *Hub {
	return &Hub{
		GameSocketMap: make(map[string]*GameSocket),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Messages:      make(chan dto.MessageTransfer),
		GameService:   gameService,
		UserService:   userService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			gameSocket, ok := h.GameSocketMap[client.Game.Code]
			if !ok {
				gameSocket = NewGameSocket(client.Game)
				h.GameSocketMap[client.Game.Code] = gameSocket
				go gameSocket.Run()
			}
			client.GameSocket = gameSocket
			gameSocket.ClientSet[client] = true
			log.Info().Interface("GameSocket.Info", gameSocket.Info).Msg("GameSocket.Info")
			// h.UserService.JoinGame(client.Ctx, client.Game.Code, client.User)
		case client := <-h.Unregister:
			// h.UserService.QuitGame(client.Ctx, client.Game.Code, client.User)
			if GameSocket, ok := h.GameSocketMap[client.Game.Code]; ok {
				if _, ok := GameSocket.ClientSet[client]; ok {
					client.CleanUp()
					delete(GameSocket.ClientSet, client)
				}
			}
		case message := <-h.Messages:
			messagesStore = append(messagesStore, message)
			log.Info().Interface("messagesStore", messagesStore).Msg("messagesStore")

			log.Info().Interface("broadcast to all games", message).Msg("broadcast")
			for _, gameSocket := range h.GameSocketMap {
				select {
				case gameSocket.Send <- message:
				default:
					close(gameSocket.Send)
					delete(h.GameSocketMap, gameSocket.Info.Code)
				}
			}
		}

	}
}

// serveWs handles websocket requests from the peer.
func (hub *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	log.Debug().Interface("URL Query", r.URL.Query()).Msg("URL Query")
	gameCode := r.URL.Query().Get("game_code")
	username := r.URL.Query().Get("username")

	// get game from database
	game, err := hub.GameService.GetGame(r.Context(), gameCode)

	if err != nil {
		log.Error().Err(err).Msg("Error when joining game")
		renderJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Error when joining game",
			"details": err,
		})
		return
	}

	// check if username exist
	exist, err := hub.UserService.UsernameExist(r.Context(), gameCode, username)

	if err != nil {
		log.Error().Err(err).Msg("Error when joining game")
		renderJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Error when joining game",
			"details": err,
		})
		return
	} else if exist {
		renderJSON(w, http.StatusBadRequest, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Error when joining game",
			"details": err,
		})
		return
	}

	log.Info().Msg("new client connected")

	user := models.User{
		Username: username,
	}
	client := NewClient(hub, conn, r.Context(), nil, game, &user)
	log.Info().Interface("client", map[string]any{
		"Game": client.Game,
		"User": client.User,
	}).Msg("client connected")
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
