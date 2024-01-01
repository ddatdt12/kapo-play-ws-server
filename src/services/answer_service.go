package services

import (
	"context"

	answerclient "github.com/ddatdt12/kapo-play-ws-server/clients/answer"
	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/utils"
	"github.com/rs/zerolog/log"
)

type IAnswerService interface {
	SaveToUser(ctx context.Context, gameCode string, username string, answer *models.Answer) error
	SaveToQuestion(ctx context.Context, gameCode string, answer *models.Answer) error
	GetAnswersOfQuestion(ctx context.Context, gameCode string, questionID uint) ([]*models.Answer, error)
	GetAnswersOfUser(ctx context.Context, gameCode string, username string) ([]*models.Answer, error)
}

type AnswerService struct {
	redis        *db.RedisImpl
	leaderboard  ILeaderboardService
	answerClient answerclient.IAnswerClient
}

func NewAnswerService(
	redis *db.RedisImpl,
	leaderboard ILeaderboardService,
	answerClient answerclient.IAnswerClient,
) *AnswerService {
	return &AnswerService{
		redis:        redis,
		leaderboard:  leaderboard,
		answerClient: answerClient,
	}
}

func (s *AnswerService) GetAnswersOfQuestion(ctx context.Context, gameCode string, questionID uint) ([]*models.Answer, error) {
	var answers []*models.Answer
	result := s.redis.DB().LRange(ctx, utils.BuildGameQuestionAnswers(gameCode, questionID), 0, -1)
	if result.Err() != nil {
		return nil, result.Err()
	}
	log.Info().Msgf("GetAnswersOfQuestion: %v", result.Val())
	err := result.ScanSlice(&answers)
	return answers, err
}

func (s *AnswerService) GetAnswersOfUser(ctx context.Context, gameCode string, username string) ([]*models.Answer, error) {
	var answers []*models.Answer
	result := s.redis.DB().LRange(ctx, utils.BuildUserAnswersKey(gameCode, username), 0, -1)
	if result.Err() != nil {
		return nil, result.Err()
	}

	err := result.ScanSlice(&answers)
	return answers, err
}

func (s *AnswerService) SaveToUser(ctx context.Context, gameCode string, username string, answer *models.Answer) error {
	err := s.redis.DB().LPush(ctx, utils.BuildUserAnswersKey(gameCode, username), answer.ToJSON()).Err()
	if err != nil {
		return err
	}

	userRank, err := s.leaderboard.GetUserRank(ctx, gameCode, username)
	if err != nil {
		return err
	}

	answer.Report = &models.UserReport{
		Points: int64(userRank.Points),
		Rank:   int(userRank.Rank),
	}
	log.Info().Interface("answer out side", answer).Msg("SaveToUser out side")
	// Send answer to answer service
	go func() {
		answer := answerclient.CreateAnswer{
			QuestionID: answer.QuestionID,
			Values:     answer.Values,
			GameID:     answer.GameID,
			IsCorrect:  answer.IsCorrect,
			Points:     answer.Points,
			AnswerTime: answer.AnswerTime,
			AnsweredAt: answer.AnsweredAt,
			Username:   answer.Username,
			Report: answerclient.UserReport{
				Points: int64(userRank.Points),
				Rank:   int(userRank.Rank),
			},
		}

		log.Info().Interface("answer", answer).Msg("SaveToUser")

		err = s.answerClient.CreateAnswer(&answer)
		if err != nil {
			log.Error().Msgf("Error sending answer to answer service: %v", err)
		}
		log.Info().Msgf("CreateAnswer: %v", err)
	}()

	return nil
}
func (s *AnswerService) SaveToQuestion(ctx context.Context, gameCode string, answer *models.Answer) error {
	err := s.redis.DB().LPush(ctx, utils.BuildGameQuestionAnswers(gameCode, answer.QuestionID), answer.ToJSON()).Err()
	if err != nil {
		return err
	}

	return nil
}
