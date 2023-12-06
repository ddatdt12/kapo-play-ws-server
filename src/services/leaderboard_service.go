package services

import (
	"context"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrLeaderboardNotFound = errors.New("game not found")
)

type ILeaderboardService interface {
	GetLeaderboard(ctx context.Context, gameID uint) (*models.Leaderboard, error)
	GetUserRank(ctx context.Context, gameID uint, username string) (*models.LeaderboardUser, error)
	IncrementPointsLeaderboard(ctx context.Context, gameID uint, username string, points uint) error
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

func (s *LeaderboardService) GetLeaderboard(ctx context.Context, gameID uint) (*models.Leaderboard, error) {
	// zResult, err := s.redis.DB().ZRangeWithScores(ctx, buildLeaderboardKey(gameID), 0, -1).Result()

	// if err != nil {
	// 	return nil, errors.Wrap(err, "get game")
	// }

	// leaderboardItems := buildLeaderboardUsers(zResult)
	// leaderboard := &models.Leaderboard{
	// 	ID:     gameID,
	// 	GameID: gameID,
	// 	Items:  leaderboardItems,
	// }

	leaderboard := &models.Leaderboard{
		GameID: gameID,
		Items: []*models.LeaderboardUser{
			{
				User: models.User{
					Username: "dat",
				},
				Username: "dat",
				Points:   100,
				Rank:     1,
			},
			{
				User: models.User{
					Username: "dat2",
				},
				Username: "dat2",
				Points:   60,
				Rank:     2,
			},
			{
				User: models.User{
					Username: "dat23",
				},
				Username: "dat23",
				Points:   40,
				Rank:     3,
			},
		},
	}

	return leaderboard, nil
}

func (s *LeaderboardService) GetUserRank(ctx context.Context, gameID uint, username string) (*models.LeaderboardUser, error) {
	userRank := models.LeaderboardUser{
		User: models.User{
			Username: "dat23",
		},
		Username: "dat23",
		Points:   40,
		Rank:     3,
	}

	return &userRank, nil
}

func (s *LeaderboardService) IncrementPointsLeaderboard(ctx context.Context, gameID uint, username string, points uint) error {
	err := s.redis.DB().ZIncrBy(ctx, buildLeaderboardKey(gameID), float64(points), username).Err()

	if err != nil {
		return errors.Wrap(err, "increment points leaderboard")
	}

	return nil
}

func buildLeaderboardKey(gameID uint) string {
	return "leaderboard:" + string(gameID)
}

func buildLeaderboardUsers(zResult []redis.Z) []*models.LeaderboardUser {
	items := make([]*models.LeaderboardUser, len(zResult))
	for _, z := range zResult {
		items = append(items, &models.LeaderboardUser{
			Username: z.Member.(string),
			Points:   uint(z.Score),
		})
	}
	return items
}
