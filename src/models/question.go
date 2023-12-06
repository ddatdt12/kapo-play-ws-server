package models

import (
	"encoding/json"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/types"
	"github.com/rs/zerolog/log"
)

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeTrueFalse      QuestionType = "true_false"
	QuestionTypeTypeAnswer     QuestionType = "type_answer"
	QuestionOpenEnded          QuestionType = "open_ended"
	QuestionTypePoll           QuestionType = "poll"
	QuestionWordCloud          QuestionType = "word_cloud"
)

type Question struct {
	ID         uint
	Content    string
	Type       QuestionType
	TemplateID uint
	LimitTime  uint
	Points     uint
	Choices    []*QuestionChoice
	StartAt    types.NullableTime
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type QuestionChoice struct {
	ID        uint
	Content   string
	IsCorrect bool
}

func (s QuestionChoice) ValidTypes() []QuestionType {
	return []QuestionType{QuestionTypeMultipleChoice, QuestionTypeTypeAnswer, QuestionTypeTrueFalse, QuestionTypePoll, QuestionWordCloud}
}

func (s QuestionType) IsValid() bool {
	switch s {
	case QuestionTypeMultipleChoice, QuestionTypeTrueFalse, QuestionTypeTypeAnswer, QuestionOpenEnded, QuestionTypePoll, QuestionWordCloud:
		return true
	}

	return false
}

func (q Question) VerifyAnswers(answers []any) bool {
	log.Info().Msgf("VerifyAnswers: %v", answers)
	log.Info().Msgf("VerifyAnswers - Choices: %v", q.Choices)
	if len(answers) == 0 {
		return false
	}

	if q.Type == QuestionTypeMultipleChoice || q.Type == QuestionTypeTrueFalse {
		for _, choice := range q.Choices {
			if choice.IsCorrect {
				for _, answer := range answers {
					answerID, ok := answer.(uint)
					if !ok {
						return false
					}

					if choice.ID == answerID {
						return true
					}
				}
			}
		}
	} else if q.Type == QuestionTypeTypeAnswer {
		actualUserAnswer := answers[0].(string)
		for _, choice := range q.Choices {
			if choice.IsCorrect && choice.Content == actualUserAnswer {
				return true
			}
		}
	} else if q.Type == QuestionTypePoll || q.Type == QuestionWordCloud {
		return true
	}

	return false
}

func (m *Question) UnmarshalBinary(data []byte) error {
	log.Info().Msgf("UnmarshalBinary Question: %v", string(data))
	return json.Unmarshal(data, m)
}

func (m *Question) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *QuestionChoice) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *QuestionChoice) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
