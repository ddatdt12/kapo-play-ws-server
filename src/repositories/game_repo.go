package repositories

import (
	"context"
	"log"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/redis/go-redis/v9"
)

type IGameRepository interface {
	GetByCode(ctx context.Context, code string) (*models.Game, error)
}

type GameRepository struct {
	redis *db.RedisImpl
}

func NewGameRepository(
	db *db.RedisImpl,
) *GameRepository {
	return &GameRepository{
		redis: db,
	}
}

func (r *GameRepository) GetByCode(ctx context.Context, code string) (*models.Game, error) {
	var game models.Game
	log.Println("GetByCode", code)
	err := r.redis.DB().Get(ctx, buildGameKey(code)).Scan(&game)

	if err != nil {
		if err == redis.Nil {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	return &game, nil
}
