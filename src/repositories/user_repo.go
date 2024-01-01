package repositories

import (
	"context"
	"encoding/json"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/utils"
)

type IUserRepository interface {
	UsernameExist(ctx context.Context, gameCode string, username string) (bool, error)
	Add(ctx context.Context, gameCode string, user *models.User) error
	Remove(ctx context.Context, gameCode string, username string) error
}

type UserRepository struct {
	redis *db.RedisImpl
}

func NewUserRepository(
	redis *db.RedisImpl,
) *UserRepository {
	return &UserRepository{
		redis: redis,
	}
}

func (r *UserRepository) UsernameExist(ctx context.Context, gameCode string, username string) (bool, error) {
	return r.redis.DB().HExists(ctx, utils.BuildGameKeyMembers(gameCode), username).Result()
}

func (r *UserRepository) Add(ctx context.Context, gameCode string, user *models.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = r.redis.DB().HSet(ctx, utils.BuildGameKeyMembers(gameCode), user.Username, userJson).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Remove(ctx context.Context, gameCode string, username string) error {
	err := r.redis.DB().HDel(ctx, utils.BuildGameKeyMembers(gameCode), username).Err()
	if err != nil {
		return err
	}

	return nil
}
