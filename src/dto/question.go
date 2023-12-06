package dto

import (
	"encoding/json"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/src/models"
)

type QuestionRes struct {
	ID        uint                `json:"id"`
	Content   string              `json:"content"`
	Type      models.QuestionType `json:"type"`
	LimitTime uint                `json:"limitTime"`
	Points    uint                `json:"points"`
	Choices   []QuestionChoiceRes `json:"choices"`
	StartAt   time.Time           `json:"startAt"`
	EndAt     time.Time           `json:"endAt"`
}

type QuestionResult struct {
	Question          *QuestionRes              `json:"question"`
	Answer            *models.Answer            `json:"answer"`
	QuestionStatistic *models.QuestionStatistic `json:"questionStatistic"`
}

type QuestionChoiceRes struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"isCorrect,omitempty"`
}

func NewQuestionRes(question *models.Question) *QuestionRes {
	choices := make([]QuestionChoiceRes, 0)
	for _, choice := range question.Choices {
		choices = append(choices, QuestionChoiceRes{
			ID:      choice.ID,
			Content: choice.Content,
		})
	}

	return &QuestionRes{
		ID:        question.ID,
		Content:   question.Content,
		Type:      question.Type,
		LimitTime: question.LimitTime,
		Points:    question.Points,
		Choices:   choices,
		StartAt:   question.StartAt.Time,
		EndAt:     question.StartAt.Time.Add(5 * time.Second),
	}
}

func NewQuestionResult(question *models.Question, anwser *models.Answer, statistic *models.QuestionStatistic) *QuestionResult {
	choices := make([]QuestionChoiceRes, 0)
	for _, choice := range question.Choices {
		choices = append(choices, QuestionChoiceRes{
			ID:        choice.ID,
			Content:   choice.Content,
			IsCorrect: choice.IsCorrect,
		})
	}

	return &QuestionResult{
		Question: &QuestionRes{
			ID:        question.ID,
			Content:   question.Content,
			Type:      question.Type,
			LimitTime: question.LimitTime,
			Points:    question.Points,
			Choices:   choices,
			StartAt:   question.StartAt.Time,
			EndAt:     question.StartAt.Time.Add(time.Duration(question.LimitTime) * time.Second),
		},
		Answer:            anwser,
		QuestionStatistic: statistic,
	}

}

type AnswerQuestionReq struct {
	QuestionOffset uint  `json:"question_offset"`
	Answers        []any `json:"answers"`
}

func (m *QuestionRes) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *QuestionChoiceRes) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
