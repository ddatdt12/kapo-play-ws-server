// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/ddatdt12/kapo-play-ws-server/configs"
	"github.com/ddatdt12/kapo-play-ws-server/internal/logger"
)

func main() {
	flag.Parse()
	logger.InitLogger()
	configs.InitEnvConfigs()

	httpProtocol := InitHttpProtocol()

	httpProtocol.Listen()
}
