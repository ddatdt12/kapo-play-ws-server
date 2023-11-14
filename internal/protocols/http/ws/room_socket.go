// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"github.com/ddatdt12/kapo-play-ws-server/models"
)

// Client is a middleman between the websocket connection and the hub.
type RoomSocket struct {
	Info *models.Room
	// Buffered channel of outbound messages.
	Send chan models.Message
	// Registered clients.
	Clients map[*Client]bool
}
