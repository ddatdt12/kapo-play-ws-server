package services

import (
	"context"
	"errors"

	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/repositories"
)

var (
	ErrGameNotFound = errors.New("game not found")
)

type IGameService interface {
	GetGame(ctx context.Context, code string) (*models.Game, error)
}

type GameService struct {
	gameRepo *repositories.GameRepository
}

func NewGameService(
	gameRepo *repositories.GameRepository,
) *GameService {
	return &GameService{
		gameRepo: gameRepo,
	}
}

func (s *GameService) GetGame(ctx context.Context, code string) (*models.Game, error) {
	game, err := s.gameRepo.GetByCode(ctx, code)

	if err != nil {
		return nil, err
	}

	return game, nil
}
