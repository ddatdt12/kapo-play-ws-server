package services

import (
	"context"
	"fmt"
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

type IQuestionService interface {
	GetQuestion(ctx context.Context, gameCode string, questionOffset uint) (*models.Question, error)
	AnwserQuestion(ctx context.Context, gameCode string, user *models.User, answerDto dto.AnswerQuestionReq) (*models.Answer, error)
	StatisticAnswersOfQuestion(ctx context.Context, gameCode string, questionOffset uint) (*models.QuestionStatistic, error)
	Update(ctx context.Context, gameCode string, question *models.Question) error
}

type QuestionService struct {
	userService   IUserService
	questionRepo  repositories.IQuestionRepository
	leaderboard   ILeaderboardService
	answerService IAnswerService
	gameService   IGameService
}

func NewQuestionService(
	questionRepo repositories.IQuestionRepository,
	gameService IGameService,
	userService IUserService,
	leaderboard ILeaderboardService,
	answerService IAnswerService,
) *QuestionService {
	return &QuestionService{
		gameService:   gameService,
		questionRepo:  questionRepo,
		userService:   userService,
		leaderboard:   leaderboard,
		answerService: answerService,
	}
}

func (s *QuestionService) GetQuestion(ctx context.Context, gameCode string, questionOffset uint) (*models.Question, error) {
	question, err := s.questionRepo.GetByOffset(ctx, gameCode, questionOffset)

	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrQuestionNotFound
		}
		return nil, errors.Wrap(err, "GetByOffset")
	}

	return question, nil
}

func (s *QuestionService) Update(ctx context.Context, gameCode string, question *models.Question) error {
	question.Game = &models.Game{
		Code: gameCode,
	}
	err := s.questionRepo.Update(ctx, question)
	if err != nil {
		return errors.Wrap(err, "Update")
	}

	return nil
}

func (s *QuestionService) AnwserQuestion(
	ctx context.Context,
	gameCode string,
	user *models.User,
	answerDto dto.AnswerQuestionReq) (*models.Answer, error) {
	game, err := s.gameService.GetGame(ctx, gameCode)
	if err != nil {
		return nil, err
	}

	question, err := s.questionRepo.GetByOffset(ctx, gameCode, answerDto.QuestionOffset)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}
	answeredAt := time.Now()
	log.Info().Interface("answeredAt", answeredAt).Msg("AnwserQuestion")
	log.Info().Interface("question.startedAt", question.StartedAt).Msg("AnwserQuestion")
	// isTimeOut := answerDto.AnsweredAt.After(question.GetEndedTime())

	// log.Info().Interface("isTimeOut", isTimeOut).Msg("isTimeOut")

	if question.StartedAt == nil {
		startedAt := answeredAt.Add(-time.Duration(question.LimitTime) * time.Second)
		question.StartedAt = &startedAt
	}

	isCorrect := question.VerifyAnswers(answerDto.Answers)
	var answerTime time.Duration = answeredAt.Sub(*question.StartedAt)
	answer := models.Answer{
		Values:     answerDto.Answers,
		QuestionID: question.ID,
		AnswerTime: answerTime.Seconds(),
		IsCorrect:  isCorrect,
		Points:     calculateScore(answerTime.Milliseconds(), question),
		User:       user,
		GameID:     game.ID,
		Username:   user.Username,
		AnsweredAt: answeredAt,
	}

	if isCorrect {
		err = s.leaderboard.IncrementPointsLeaderboard(ctx, gameCode, user.Username, uint(answer.Points))
		if err != nil {
			return nil, errors.Wrap(err, "IncrementPointsLeaderboard")
		}
	}

	err = s.answerService.SaveToUser(ctx, gameCode, user.Username, &answer)
	if err != nil {
		log.Error().Err(err).Msg("SaveToUser")
	}

	err = s.answerService.SaveToQuestion(ctx, gameCode, &answer)
	if err != nil {
		log.Error().Err(err).Msg("SaveToQuestion")
	}

	return &answer, nil
}

func (s *QuestionService) StatisticAnswersOfQuestion(
	ctx context.Context,
	gameCode string,
	questionOffset uint,
) (*models.QuestionStatistic, error) {
	question, err := s.questionRepo.GetByOffset(ctx, gameCode, questionOffset)

	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}

	if question == nil {
		return nil, ErrQuestionNotFound
	}

	answers, err := s.answerService.GetAnswersOfQuestion(ctx, gameCode, question.ID)

	if err != nil {
		log.Error().Err(err).Msg("GetAnswersOfQuestion")
	}

	log.Info().Interface("answers of question", answers).Msg("StatisticAnswersOfQuestion")

	totalAnswered := 0
	totalCorrect := 0

	for _, answer := range answers {
		if answer.IsCorrect {
			totalCorrect++
		}
		totalAnswered += len(answer.Values)
	}

	questionStatistic := &models.QuestionStatistic{
		QuestionID:         question.ID,
		ChoiceStatistics:   buildChoiceStatistic(question, answers),
		TotalAnswer:        totalAnswered,
		TotalCorrectAnswer: totalCorrect,
	}

	return questionStatistic, nil
}

func calculateScore(answerTime int64, question *models.Question) int64 {
	return 10
	var limitTime int64 = int64(question.LimitTime * 1000)
	if answerTime > limitTime {
		return 0
	}
	percent := float64(answerTime) / float64(limitTime)
	return int64(float64(question.Points) * percent)
}

func buildAnswerCountMap(answers []*models.Answer) map[string]int {
	answerCountMap := make(map[string]int)
	for _, answer := range answers {
		for _, value := range answer.Values {
			key := fmt.Sprint(value)
			answerCountMap[key]++
		}
	}
	return answerCountMap
}

func buildChoiceStatistic(question *models.Question, answers []*models.Answer) map[string]int {
	answerCountMap := buildAnswerCountMap(answers)
	choicesStatistic := make(map[string]int)
	for _, choice := range question.Choices {
		key := choice.Content
		if models.IsQuestionTypeGroupMultipleChoice(question.Type) {
			key = fmt.Sprint(choice.ID)
		}
		choicesStatistic[key] = 0
	}

	for key, answerCount := range answerCountMap {
		choicesStatistic[key] = answerCount
	}
	return choicesStatistic
}
