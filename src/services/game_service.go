package services

import (
	"context"
	"time"

	gameclient "github.com/ddatdt12/kapo-play-ws-server/clients/game"
	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/types"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/repositories"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	ErrGameNotFound = errors.New("game not found")
)

type IGameService interface {
	GetGame(ctx context.Context, code string) (*models.Game, error)
	GetGameState(ctx context.Context, code string) (*models.GameState, error)
	UpdateGameState(ctx context.Context, code string, gameState *models.GameState) error
	StartGame(ctx context.Context, code string) error
	EndGame(ctx context.Context, code string) error
	PlayAgain(ctx context.Context, code string) error
	Update(ctx context.Context, game *models.Game) (*models.Game, error)
}

type GameService struct {
	gameRepo   repositories.IGameRepository
	redis      *db.RedisImpl
	gameClient gameclient.IGameClient
}

func NewGameService(
	gameRepo repositories.IGameRepository,
	redis *db.RedisImpl,
	gameClient gameclient.IGameClient,
) *GameService {
	return &GameService{
		gameRepo:   gameRepo,
		redis:      redis,
		gameClient: gameClient,
	}
}

func (s *GameService) GetGame(ctx context.Context, code string) (*models.Game, error) {
	game, err := s.gameRepo.GetByCode(ctx, code)

	if err != nil {
		return nil, errors.Wrap(err, "get game")
	}

	return game, nil
}

func (s *GameService) StartGame(ctx context.Context, code string) error {
	game, err := s.GetGame(ctx, code)

	if err != nil {
		return errors.Wrap(err, "get game")
	}

	game.Status = models.GameStatusPlaying
	game.StartTime = types.NewNullableTime(time.Now())

	go func() {
		err := s.gameClient.Start(game.ID)
		log.Info().Msgf("StartGame: %v", err)
	}()

	return nil
}

func (s *GameService) EndGame(ctx context.Context, code string) error {
	game, err := s.GetGame(ctx, code)

	if err != nil {
		return errors.Wrap(err, "get game")
	}

	game.Status = models.GameStatusEnded
	game.EndTime = types.NewNullableTime(time.Now())

	go func() {
		err := s.gameClient.End(game.ID)
		log.Info().Msgf("EndGame: %v", err)
	}()
	return nil
}

func (s *GameService) GetGameState(ctx context.Context, code string) (*models.GameState, error) {
	gameState, err := s.gameRepo.GetGameState(ctx, code)

	if err != nil {
		return nil, errors.Wrap(err, "GetGameState")
	}

	return gameState, nil
}

func (s *GameService) UpdateGameState(ctx context.Context, code string, gameState *models.GameState) error {
	err := s.gameRepo.UpdateGameState(ctx, code, gameState)

	if err != nil {
		return errors.Wrap(err, "UpdateGameState")
	}

	return nil
}

func (s *GameService) PlayAgain(ctx context.Context, code string) error {
	game, err := s.GetGame(ctx, code)
	if err != nil {
		return errors.Wrap(err, "get game")
	}

	if err = s.gameRepo.Reset(ctx, code); err != nil {
		return errors.Wrap(err, "update game")
	}

	if err = s.gameClient.PlayAgain(game.ID); err != nil {
		return errors.Wrap(err, "play again")
	}

	return nil
}

func (s *GameService) Update(ctx context.Context, game *models.Game) (*models.Game, error) {
	game, err := s.gameRepo.Update(ctx, game)

	// if game.Status == models.GameStatusFinished {
	// 	err = s.redis.DB().Del(ctx, game.Code).Err()
	// }

	if err != nil {
		return nil, errors.Wrap(err, "update game")
	}

	return game, nil
}
