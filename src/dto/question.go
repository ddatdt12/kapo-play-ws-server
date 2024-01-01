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
	StartedAt time.Time           `json:"startedAt"`
	EndedAt   time.Time           `json:"endedAt"`
}

type QuestionResult struct {
	Question          *QuestionRes              `json:"question"`
	Answer            *models.Answer            `json:"answer"`            // Answer of user
	QuestionStatistic *models.QuestionStatistic `json:"questionStatistic"` // Statistic of question for host
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
		StartedAt: *question.StartedAt,
		EndedAt:   question.StartedAt.Add(time.Duration(question.LimitTime) * time.Second),
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
			StartedAt: *question.StartedAt,
			EndedAt:   question.StartedAt.Add(time.Duration(question.LimitTime) * time.Second),
		},
		Answer:            anwser,
		QuestionStatistic: statistic,
	}

}

type AnswerQuestionReq struct {
	QuestionOffset uint      `json:"questionOffset"`
	Answers        []any     `json:"answers"`
	AnsweredAt     time.Time `json:"answeredAt"`
}

func (m *QuestionRes) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *QuestionChoiceRes) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
