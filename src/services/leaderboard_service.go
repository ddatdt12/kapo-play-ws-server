package services

import (
	"context"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	ErrLeaderboardNotFound = errors.New("leaderboard not found")
)

type ILeaderboardService interface {
	GetLeaderboard(ctx context.Context, gameCode string) (*models.Leaderboard, error)
	GetUserRank(ctx context.Context, gameCode string, username string) (*models.LeaderboardUser, error)
	IncrementPointsLeaderboard(ctx context.Context, gameCode string, username string, points uint) error
}

type LeaderboardService struct {
	redis *db.RedisImpl
}

func NewLeaderboardService(
	redis *db.RedisImpl,
) *LeaderboardService {
	return &LeaderboardService{
		redis: redis,
	}
}

func (s *LeaderboardService) GetLeaderboard(ctx context.Context, gameCode string) (*models.Leaderboard, error) {
	zResult, err := s.redis.DB().ZRevRangeWithScores(ctx, buildLeaderboardKey(gameCode), 0, -1).Result()

	if err != nil {
		return nil, errors.Wrap(err, "get game")
	}

	leaderboardItems := buildLeaderboardUsers(zResult)
	leaderboard := &models.Leaderboard{
		GameCode: gameCode,
		Items:    leaderboardItems,
	}

	return leaderboard, nil
}

func (s *LeaderboardService) GetUserRank(ctx context.Context, gameCode string, username string) (*models.LeaderboardUser, error) {
	zResult, err := s.redis.DB().ZRevRankWithScore(ctx, buildLeaderboardKey(gameCode), username).Result()
	userRank := models.LeaderboardUser{
		User: models.User{
			Username: username,
		},
		Username: username,
		Points:   0,
		Rank:     0,
	}
	if err != nil {
		if err != redis.Nil {
			log.Error().Err(err).Msg("get user rank")
		}
		return &userRank, nil
	}

	userRank.Points = uint(zResult.Score)
	userRank.Rank = uint(zResult.Rank + 1)

	return &userRank, nil
}

func (s *LeaderboardService) IncrementPointsLeaderboard(ctx context.Context, gameCode string, username string, points uint) error {
	err := s.redis.DB().ZIncrBy(ctx, buildLeaderboardKey(gameCode), float64(points), username).Err()

	if err != nil {
		return errors.Wrap(err, "increment points leaderboard")
	}

	return nil
}

func buildLeaderboardKey(gameCode string) string {
	return "game:" + gameCode + ":leaderboard"
}

func buildLeaderboardUsers(zResult []redis.Z) []*models.LeaderboardUser {
	log.Info().Interface("zResult", zResult).Msg("buildLeaderboardUsers")
	items := make([]*models.LeaderboardUser, len(zResult))
	for i, z := range zResult {
		if z.Member == nil {
			continue
		}
		items[i] = &models.LeaderboardUser{
			Username: z.Member.(string),
			Points:   uint(z.Score),
			Rank:     uint(i + 1),
		}
	}
	return items
}
