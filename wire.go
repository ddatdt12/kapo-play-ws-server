//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ddatdt12/kapo-play-ws-server/internal/protocols/http"

	"github.com/google/wire"
)

func InitHttpProtocol() *http.HttpImpl {
	wire.Build(
		http.NewHttpProtocol,
		// cache.NewRedisClient,
	)

	return &http.HttpImpl{}
}
