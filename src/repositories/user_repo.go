package repositories

import (
	"context"
	"encoding/json"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
)

type IUserRepository interface {
	UsernameExist(ctx context.Context, gameCode string, username string) (bool, error)
	AddUser(ctx context.Context, gameCode string, user models.User) (*models.User, error)
}

type UserRepository struct {
	redis *db.RedisImpl
}

func NewUserRepository(
	db *db.RedisImpl,
) *UserRepository {
	return &UserRepository{
		redis: db,
	}
}

func (r *UserRepository) UsernameExist(ctx context.Context, gameCode string, username string) (bool, error) {
	return r.redis.DB().HExists(ctx, buildGameKeyMembers(gameCode), username).Result()
}

func (r *UserRepository) AddUser(ctx context.Context, gameCode string, user models.User) (*models.User, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = r.redis.DB().HSet(ctx, buildGameKeyMembers(gameCode), userJson).Err()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func buildGameKey(code string) string {
	return "game:" + code
}

func buildGameKeyMembers(code string) string {
	return "game:" + code + ":members"
}
