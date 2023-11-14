// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"net/http"
	"strconv"

	"github.com/ddatdt12/kapo-play-ws-server/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var messagesStore = make([]models.Message, 0)

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
	// Rooms and clients
	Rooms map[*RoomSocket]map[*Client]bool

	// Rooms and RoomSockets
	RoomSocket map[int]*RoomSocket

	// Inbound Messages from the clients.
	Broadcast chan []byte

	// Inbound Messages from the clients.
	Messages chan models.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[*RoomSocket]map[*Client]bool),
		RoomSocket: make(map[int]*RoomSocket),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Messages:   make(chan models.Message),
		// clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			roomSocket, ok := h.RoomSocket[client.Room.Id]
			if !ok {
				roomSocket = &RoomSocket{Info: client.Room, Send: make(chan models.Message, 256)}
				h.RoomSocket[client.Room.Id] = roomSocket
				h.Rooms[roomSocket] = make(map[*Client]bool)
			}
			log.Info().Interface("h.RoomSocket", h.RoomSocket).Msg("RoomSocket")

			h.Rooms[roomSocket][client] = true
		case client := <-h.Unregister:
			if RoomSocket, ok := h.RoomSocket[client.Room.Id]; ok {
				if _, ok := h.Rooms[RoomSocket]; ok {
					close(client.Send)
					delete(h.Rooms[RoomSocket], client)
				}
			}
		case message := <-h.Messages:
			RoomSocket, ok := h.RoomSocket[message.RoomID]
			log.Info().Interface("broadcast to members in Room", message).Msg("broadcast")
			if !ok {
				continue
			}
			members, ok := h.Rooms[RoomSocket]
			if !ok {
				continue
			}

			messagesStore = append(messagesStore, message)

			log.Info().Interface("messagesStore", messagesStore).Msg("messagesStore")

			for member := range members {
				select {
				case member.Send <- message:
				default:
					close(member.Send)
					delete(h.Rooms[RoomSocket], member)
				}
			}
		}

	}
}

// serveWs handles websocket requests from the peer.
func (hub *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return
	}

	log.Info().Msg("new client connected")
	log.Debug().Interface("URL Query", r.URL.Query()).Msg("")
	// get room id from query string
	roomId, error := strconv.Atoi(r.URL.Query().Get("room_id"))
	if error != nil {
		log.Error().Err(error).Msg("error")
		conn.Close()
		return
	}
	username := r.URL.Query().Get("username")
	if error != nil {
		log.Error().Err(error).Msg("error")
		conn.Close()
		return
	}
	// Get room from database
	room := models.GetRoom(roomId)
	user := models.CreateUser(username)

	if room == nil {
		log.Error().Err(error).Msg("error")
		conn.Close()
		return
	}

	client := &Client{Hub: hub, Conn: conn, Room: room, User: user, Send: make(chan models.Message, 256)}
	log.Info().Interface("client", map[string]any{
		"Room": client.Room,
		"User": client.User,
	}).Msg("client connected")
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
