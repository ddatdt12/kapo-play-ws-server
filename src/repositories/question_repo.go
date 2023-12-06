package repositories

import (
	"context"
	"errors"

	"github.com/ddatdt12/kapo-play-ws-server/src/models"
)

var (
	ErrQuestionNotFound = errors.New("question not found")
)

type IQuestionRepository interface {
	GetByOffset(ctx context.Context, gameID uint, offset int64) (*models.Question, error)
	GetByID(ctx context.Context, gameID uint, questionID uint) (*models.Question, error)
}

type QuestionRepository struct {
}

func NewQuestionRepository() *QuestionRepository {
	return &QuestionRepository{}
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

func (r *QuestionRepository) GetByOffset(ctx context.Context, gameID uint, offset int64) (*models.Question, error) {
	// var question models.Question
	// log.Info().Msgf("GetByOffset: %v %v", buildQuestionKey(gameCode), offset)
	// err := r.redis.DB().LIndex(ctx, buildQuestionKey(gameCode), offset).Scan(&question)

	// if err != nil {
	// 	if err == redis.Nil {
	// 		return nil, db.ErrNotFound
	// 	}
	// 	return nil, err
	// }

	if offset >= int64(len(questions)) {
		return nil, nil
	}
	question := &questions[offset]
	return question, nil
}

func (r *QuestionRepository) GetByID(ctx context.Context, gameID uint, questionID uint) (*models.Question, error) {
	for _, question := range questions {
		if question.ID == questionID {
			return &question, nil
		}
	}

	return nil, ErrQuestionNotFound
}

func buildQuestionKey(code string) string {
	return "game:" + code + ":questions"
}
