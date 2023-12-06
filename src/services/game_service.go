package services

import (
	"context"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/types"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/repositories"
	"github.com/pkg/errors"
)

var (
	ErrGameNotFound = errors.New("game not found")
)

type IGameService interface {
	GetGame(ctx context.Context, code string) (*models.Game, error)
	StartGame(ctx context.Context, code string) error
}

type GameService struct {
	gameRepo repositories.IGameRepository
	redis    *db.RedisImpl
}

func NewGameService(
	gameRepo repositories.IGameRepository,
	redis *db.RedisImpl,
) *GameService {
	return &GameService{
		gameRepo: gameRepo,
		redis:    redis,
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

	_, err = s.gameRepo.Update(ctx, game)

	if err != nil {
		return errors.Wrap(err, "update game")
	}

	return nil
}
