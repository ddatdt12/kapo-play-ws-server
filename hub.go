// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strconv"

	"github.com/ddatdt12/kapo-play-ws-server/models"
	"github.com/rs/zerolog/log"
)

var messagesStore = make([]models.Message, 0)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Rooms and clients
	rooms map[*RoomSocket]map[*Client]bool

	// Rooms and RoomSockets
	roomSocket map[int]*RoomSocket

	// Inbound messages from the clients.
	broadcast chan []byte

	// Inbound messages from the clients.
	messages chan models.Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		rooms:      make(map[*RoomSocket]map[*Client]bool),
		roomSocket: make(map[int]*RoomSocket),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		messages:   make(chan models.Message),
		// clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			roomSocket, ok := h.roomSocket[client.room.Id]
			if !ok {
				roomSocket = &RoomSocket{hub: h, info: client.room, send: make(chan models.Message, 256)}
				h.roomSocket[client.room.Id] = roomSocket
				h.rooms[roomSocket] = make(map[*Client]bool)
			}
			log.Info().Interface("roomSocket", h.roomSocket).Msg("roomSocket")
			log.Info().Interface("h.roomSocket", h.roomSocket).Msg("roomSocket")

			h.rooms[roomSocket][client] = true
		case client := <-h.unregister:
			if roomSocket, ok := h.roomSocket[client.room.Id]; ok {
				if _, ok := h.rooms[roomSocket]; ok {
					close(client.send)
					delete(h.rooms[roomSocket], client)
				}
			}
		case message := <-h.messages:
			roomSocket, ok := h.roomSocket[message.RoomID]
			log.Info().Interface("broadcast to members in room", message).Msg("broadcast")
			if !ok {
				continue
			}
			members, ok := h.rooms[roomSocket]
			if !ok {
				continue
			}

			messagesStore = append(messagesStore, message)

			log.Info().Interface("messagesStore", messagesStore).Msg("messagesStore")

			for member := range members {
				select {
				case member.send <- message:
				default:
					close(member.send)
					delete(h.rooms[roomSocket], member)
				}
			}
		}

	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
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

	client := &Client{hub: hub, conn: conn, room: room, user: user, send: make(chan models.Message, 256)}
	log.Info().Interface("client", *client).Msg("client")
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
