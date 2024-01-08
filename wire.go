//go:build wireinject
// +build wireinject

package main

import (
	answerclient "github.com/ddatdt12/kapo-play-ws-server/clients/answer"
	gameclient "github.com/ddatdt12/kapo-play-ws-server/clients/game"
	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/internal/protocols/http"
	"github.com/ddatdt12/kapo-play-ws-server/internal/protocols/http/ws"
	"github.com/ddatdt12/kapo-play-ws-server/src/handlers"
	"github.com/ddatdt12/kapo-play-ws-server/src/repositories"
	"github.com/ddatdt12/kapo-play-ws-server/src/services"

	"github.com/google/wire"
)

var questionRepo = wire.NewSet(
	repositories.NewQuestionRepository,
	wire.Bind(
		new(repositories.IQuestionRepository),
		new(*repositories.QuestionRepository),
	),
)

var questionSvc = wire.NewSet(
	services.NewQuestionService,
	wire.Bind(
		new(services.IQuestionService),
		new(*services.QuestionService),
	),
)

var leaderboardSvc = wire.NewSet(
	services.NewLeaderboardService,
	wire.Bind(
		new(services.ILeaderboardService),
		new(*services.LeaderboardService),
	),
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

var answerSvc = wire.NewSet(
	services.NewAnswerService,
	wire.Bind(
		new(services.IAnswerService),
		new(*services.AnswerService),
	),
)

var answerClient = wire.NewSet(
	answerclient.NewAnswerClient,
	wire.Bind(
		new(answerclient.IAnswerClient),
		new(*answerclient.AnswerClient),
	),
)

var gameClient = wire.NewSet(
	gameclient.NewGameClient,
	wire.Bind(
		new(gameclient.IGameClient),
		new(*gameclient.GameClient),
	),
)

func InitHttpProtocol() *http.HttpImpl {
	wire.Build(
		http.NewHttpProtocol,
		handlers.NewHttpHandler,
		db.NewRedisClient,
		ws.NewHub,
		leaderboardSvc,
		gameSvc,
		userSvc,
		answerSvc,
		questionSvc,
		gameRepo,
		userRepo,
		questionRepo,
		answerClient,
		gameClient,
	)

	return &http.HttpImpl{}
}
