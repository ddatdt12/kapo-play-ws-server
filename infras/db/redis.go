package db

import (
	"context"
	"fmt"

	"github.com/ddatdt12/kapo-play-ws-server/configs"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisImpl struct {
	db *redis.Client
}

func NewRedisClient() *RedisImpl {
	log.Info().Msg("Connecting to redis")
	host := fmt.Sprintf("%v:%v", configs.EnvConfigs.REDIS_HOST, configs.EnvConfigs.REDIS_PORT)

	options := redis.Options{
		Addr:       host,
		Password:   configs.EnvConfigs.REDIS_PASSWORD,
		DB:         configs.EnvConfigs.REDIS_DB,
		MaxRetries: 3,
	}
	log.Info().Msgf("Redis options: %v", options)

	rdb := redis.NewClient(&options)
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connect to redis")
	}

	log.Info().Msg("Connected to redis")
	return &RedisImpl{
		db: rdb,
	}
}

func (c RedisImpl) DB() *redis.Client {
	return c.db
}

func (c RedisImpl) Close() error {
	return c.db.Close()
}
