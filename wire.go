//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/internal/protocols/http"
	"github.com/ddatdt12/kapo-play-ws-server/internal/protocols/http/ws"
	"github.com/ddatdt12/kapo-play-ws-server/src/repositories"
	"github.com/ddatdt12/kapo-play-ws-server/src/services"
	"github.com/google/wire"
)

var gameRepo = wire.NewSet(
	repositories.NewGameRepository,
	wire.Bind(
		new(repositories.IGameRepository),
		new(*repositories.GameRepository),
	),
)

var gameSvc = wire.NewSet(
	services.NewGameService,
	wire.Bind(
		new(services.IGameService),
		new(*services.GameService),
	),
)

var userRepo = wire.NewSet(
	repositories.NewUserRepository,
	wire.Bind(
		new(repositories.IUserRepository),
		new(*repositories.UserRepository),
	),
)

var userSvc = wire.NewSet(
	services.NewUserService,
	wire.Bind(
		new(services.IUserService),
		new(*services.UserService),
	),
)

func InitHttpProtocol() *http.HttpImpl {
	wire.Build(
		db.NewRedisClient,
		http.NewHttpProtocol,
		ws.NewHub,
		gameSvc,
		userSvc,
		gameRepo,
		userRepo,
	)

	return &http.HttpImpl{}
}
