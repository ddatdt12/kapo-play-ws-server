package repositories

import (
	"context"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
)

type IGameRepository interface {
	GetByCode(ctx context.Context, code string) (*models.Game, error)
	Update(ctx context.Context, game *models.Game) (*models.Game, error)
}

type GameRepository struct {
	redis *db.RedisImpl
}

func NewGameRepository(redis *db.RedisImpl) *GameRepository {
	return &GameRepository{
		redis: redis,
	}
}

func (r *GameRepository) GetByCode(ctx context.Context, code string) (*models.Game, error) {
	var game models.Game
	// err := r.redis.DB().Get(ctx, buildGameKey(code)).Scan(&game)
	// game.TotalQuestions = r.redis.DB().LLen(ctx, buildQuestionKey(code)).Val()

	// if err != nil {
	// 	if err == redis.Nil {
	// 		return nil, db.ErrNotFound
	// 	}
	// 	return nil, err
	// }

	game = models.Game{
		Code:           code,
		Name:           "Game 1",
		Status:         models.GameStatusWaiting,
		TotalQuestions: 10,
		Host: models.User{
			ID:       1,
			Username: "Dat",
		},
		Settings: models.GameSettings{
			RandomizeQuestions: true,
			RandomizeAnswers:   true,
		},
	}

	return &game, nil
}

func (r *GameRepository) Update(ctx context.Context, game *models.Game) (*models.Game, error) {
	// err := r.redis.DB().Set(ctx, buildGameKey(game.Code), game, 0).Err()

	// if err != nil {
	// 	if err == redis.Nil {
	// 		return nil, db.ErrNotFound
	// 	}
	// 	return nil, err
	// }

	return game, nil
}
