package repositories

import (
	"context"
	"errors"

	"github.com/ddatdt12/kapo-play-ws-server/infras/db"
	"github.com/ddatdt12/kapo-play-ws-server/src/models"
	"github.com/ddatdt12/kapo-play-ws-server/src/utils"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	ErrQuestionNotFound = errors.New("question not found")
)

type IQuestionRepository interface {
	GetByOffset(ctx context.Context, gameCode string, offset uint) (*models.Question, error)
	GetByID(ctx context.Context, gameCode string, questionID uint) (*models.Question, error)
	Update(ctx context.Context, game *models.Question) error
}

type QuestionRepository struct {
	redis *db.RedisImpl
}

func NewQuestionRepository(redis *db.RedisImpl) *QuestionRepository {
	return &QuestionRepository{
		redis: redis,
	}
}

var questions = []models.Question{
	{
		ID:         1,
		Content:    "What is your name?",
		Type:       models.QuestionTypeMultipleChoice,
		TemplateID: 1,
		LimitTime:  10,
		Points:     10,
		Choices: []*models.QuestionChoice{
			{
				ID:        1,
				Content:   "A",
				IsCorrect: true,
			},
			{
				ID:      2,
				Content: "B",
			},
			{
				ID:      3,
				Content: "C",
			},
			{
				ID:      4,
				Content: "D",
			},
		},
	},
	{
		ID:      2,
		Content: "How old I am?",
		Type:    models.QuestionTypeTypeAnswer,
		Choices: []*models.QuestionChoice{
			{
				ID:        1,
				Content:   "20",
				IsCorrect: true,
			},
			{
				ID:        2,
				Content:   "22",
				IsCorrect: true,
			},
		},
		TemplateID: 1,
		LimitTime:  15,
		Points:     50,
	},
}

func (r *QuestionRepository) GetByOffset(ctx context.Context, gameCode string, offset uint) (*models.Question, error) {
	var question models.Question
	length, err := r.redis.DB().LLen(ctx, utils.BuildQuestionKey(gameCode)).Result()
	if err != nil {
		return nil, err
	}

	if offset >= uint(length) {
		return nil, nil
	}

	err = r.redis.DB().LIndex(ctx, utils.BuildQuestionKey(gameCode), int64(offset)).Scan(&question)
	log.Info().Msgf("GetByOffset: %v - %v", offset, question)

	if err != nil {
		if err == redis.Nil {
			return nil, db.ErrNotFound
		}
		return nil, err
	}
	question.Offset = offset

	return &question, nil
}

func (r *QuestionRepository) GetByID(ctx context.Context, gameCode string, questionID uint) (*models.Question, error) {
	for _, question := range questions {
		if question.ID == questionID {
			return &question, nil
		}
	}

	return nil, ErrQuestionNotFound
}

func (r *QuestionRepository) Update(ctx context.Context, question *models.Question) error {
	log.Info().Interface("Repo update question", question).Msg("Update question")
	err := r.redis.DB().LSet(ctx, utils.BuildQuestionKey(question.Game.Code), int64(question.Offset), question).Err()

	if err != nil {
		if err == redis.Nil {
			return db.ErrNotFound
		}
		return err
	}

	return nil
}
