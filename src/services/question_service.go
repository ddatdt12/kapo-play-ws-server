package services

import (
	"context"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/dto"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/repositories"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	ErrQuestionNotFound = errors.New("question not found")
)

var UserAnwsersMap = make(map[uint]*[]models.Answer, 0)
var QuestionAnswersMap = make(map[uint]*[]models.Answer, 0)

type IQuestionService interface {
	GetQuestion(ctx context.Context, gameID uint, questionOffset int64) (*models.Question, error)
	AnwserQuestion(ctx context.Context, gameID uint, user *models.User, answerDto dto.AnswerQuestionReq) (*models.Answer, error)
	StatisticAnswersOfQuestion(ctx context.Context, gameID uint, questionID uint) (*models.QuestionStatistic, error)
}

type QuestionService struct {
	userService  IUserService
	questionRepo repositories.IQuestionRepository
	leaderboard  ILeaderboardService
}

func NewQuestionService(
	questionRepo repositories.IQuestionRepository,
	userService IUserService,
	leaderboard ILeaderboardService,
) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
		userService:  userService,
		leaderboard:  leaderboard,
	}
}

func (s *QuestionService) GetQuestion(ctx context.Context, gameID uint, questionOffset int64) (*models.Question, error) {
	question, err := s.questionRepo.GetByOffset(ctx, gameID, questionOffset)

	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrQuestionNotFound
		}
		return nil, errors.Wrap(err, "GetByOffset")
	}

	return question, nil
}

func (s *QuestionService) AnwserQuestion(
	ctx context.Context,
	gameID uint,
	user *models.User,
	answerDto dto.AnswerQuestionReq) (*models.Answer, error) {
	question, err := s.questionRepo.GetByOffset(ctx, gameID, int64(answerDto.QuestionOffset))
	anwserAt := time.Now()
	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}

	isCorrect := question.VerifyAnswers(answerDto.Answers)
	var answerTime time.Duration = anwserAt.Sub(question.StartAt.Time)
	answer := models.Answer{
		Choices:    answerDto.Answers,
		QuestionID: question.ID,
		AnswerTime: answerTime.Milliseconds(),
		IsCorrect:  isCorrect,
		Point:      calculateScore(answerTime.Milliseconds(), question),
	}
	storeAnswer(gameID, user, answer)
	if isCorrect {
		err = s.leaderboard.IncrementPointsLeaderboard(ctx, gameID, user.Username, 100)
		if err != nil {
			return nil, errors.Wrap(err, "IncrementPointsLeaderboard")
		}
	}

	return &answer, nil
}

func (s *QuestionService) StatisticAnswersOfQuestion(
	ctx context.Context,
	gameID uint,
	questionID uint,
) (*models.QuestionStatistic, error) {
	question, err := s.questionRepo.GetByID(ctx, gameID, questionID)

	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}
	// mock
	questionStatistic := &models.QuestionStatistic{
		QuestionID:         question.ID,
		AnswerCountMap:     make(map[uint]int),
		TotalAnswer:        10,
		TotalCorrectAnswer: 5,
		TotalWrongAnswer:   5,
	}

	questionStatistic.AnswerCountMap[1] = 5
	questionStatistic.AnswerCountMap[2] = 5
	questionStatistic.AnswerCountMap[3] = 2
	questionStatistic.AnswerCountMap[4] = 3

	// answers := QuestionAnswersMap[questionID]
	// log.Info().Interface("answers", answers).Msg("StatisticAnswersOfQuestion")
	// for _, answer := range *answers {
	// 	if answer.IsCorrect {
	// 		questionStatistic.TotalCorrectAnswer++
	// 	} else {
	// 		questionStatistic.TotalWrongAnswer++
	// 	}
	// 	for _, choice := range answer.Choices {
	// 		questionStatistic.AnswerCountMap[choice]++
	// 	}
	// }

	return questionStatistic, nil
}

func calculateScore(answerTime int64, question *models.Question) int64 {
	// var limitTime int64 = int64(question.LimitTime * 1000)
	// if answerTime > limitTime {
	// 	return 0
	// }
	// percent := float64(answerTime) / float64(limitTime)
	// return int64(float64(question.Points) * percent)

	return 100
}

func storeAnswer(gameID uint, user *models.User, answer models.Answer) {
	if UserAnwsersMap[gameID] == nil {
		UserAnwsersMap[gameID] = &[]models.Answer{}
	}
	*UserAnwsersMap[gameID] = append(*UserAnwsersMap[gameID], answer)
	if QuestionAnswersMap[answer.QuestionID] == nil {
		QuestionAnswersMap[answer.QuestionID] = &[]models.Answer{}
	}
	*QuestionAnswersMap[answer.QuestionID] = append(*QuestionAnswersMap[answer.QuestionID], answer)

	log.Info().Interface("UserAnwsersMap", UserAnwsersMap).Msg("storeAnswer")
	log.Info().Interface("QuestionAnswersMap", QuestionAnswersMap).Msg("storeAnswer")
}
