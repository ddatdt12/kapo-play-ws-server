package repositories

import (
	"context"
	"fmt"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/utils"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

type IGameRepository interface {
	GetByCode(ctx context.Context, code string) (*models.Game, error)
	Update(ctx context.Context, game *models.Game) (*models.Game, error)
	GetGameState(ctx context.Context, code string) (*models.GameState, error)
	UpdateGameState(ctx context.Context, code string, gameState *models.GameState) error
	Reset(ctx context.Context, code string) error
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
	err := r.redis.DB().Get(ctx, utils.BuildGameKey(code)).Scan(&game)

	if err != nil {
		if err == redis.Nil {
			return nil, db.ErrNotFound
		}
		return nil, err
	}
	game.TotalQuestions = r.redis.DB().LLen(ctx, utils.BuildQuestionKey(code)).Val()

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

func (r *GameRepository) GetGameState(ctx context.Context, code string) (*models.GameState, error) {
	var gameState models.GameState
	err := r.redis.DB().Get(ctx, utils.BuildGameStateKey(code)).Scan(&gameState)

	if err != nil {
		if err == redis.Nil {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	return &gameState, nil
}

func (r *GameRepository) UpdateGameState(ctx context.Context, code string, gameState *models.GameState) error {
	err := r.redis.DB().Set(ctx, utils.BuildGameStateKey(code), gameState, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

func getIgnoreKeys(code string) []string {
	return []string{
		utils.BuildGameKey(code),
		utils.BuildQuestionKey(code),
	}
}

func (r *GameRepository) Reset(ctx context.Context, code string) error {
	keys := r.redis.DB().Keys(ctx, fmt.Sprintf("%s:*", utils.BuildGameKey(code))).Val()

	if len(keys) == 0 {
		return nil
	}

	leftKeys, _ := lo.Difference(keys, getIgnoreKeys(code))

	err := r.redis.DB().Del(ctx, leftKeys...).Err()

	if err != nil {
		return err
	}

	return nil
}
