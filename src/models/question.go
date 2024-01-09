package models

import (
	"encoding/json"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/constants"
)

type QuestionType string

type QuestionTypeGroup string
type QuestionStatus string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeTrueFalse      QuestionType = "true_false"
	QuestionTypeTypeAnswer     QuestionType = "type_answer"
	QuestionOpenEnded          QuestionType = "open_ended"
	QuestionTypePoll           QuestionType = "poll"
	QuestionWordCloud          QuestionType = "word_cloud"
)

const (
	QuestionStatusWaiting QuestionStatus = "waiting"
	QuestionStatusPlaying QuestionStatus = "playing"
	QuestionStatusEnded   QuestionStatus = "ended"
)

const (
	AcceptableDelayTime = 1
)

type Question struct {
	ID         uint              `json:"id"`
	Offset     uint              `json:"offset"`
	Content    string            `json:"content"`
	Type       QuestionType      `json:"type"`
	TemplateID uint              `json:"templateId"`
	Game       *Game             `json:"game"`
	LimitTime  uint              `json:"limitTime"`
	Points     uint              `json:"points"`
	Choices    []*QuestionChoice `json:"choices"`
	Status     QuestionStatus    `json:"status"`
	StartedAt  *time.Time        `json:"startedAt"`
	EndedAt    *time.Time        `json:"endedAt"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type QuestionChoice struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	IsCorrect bool   `json:"isCorrect"`
}

func (s Question) GetEndedTime() time.Time {
	return s.StartedAt.Add(time.Duration(s.LimitTime) * time.Second)
}
func (q *Question) Start() {
	startTime := time.Now().Add(constants.WaitingTimeBeforeStart * time.Second)
	q.StartedAt = &startTime
	q.Status = QuestionStatusPlaying
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
	if q.Type == QuestionTypeMultipleChoice || q.Type == QuestionTypeTrueFalse {
		if len(answers) == 0 {
			return false
		}
		for _, choice := range q.Choices {
			if choice.IsCorrect {
				for _, answer := range answers {
					answerId, ok := answer.(float64)

					if ok && uint(answerId) == choice.ID {
						return true
					}
				}
			}
		}
	} else if q.Type == QuestionTypeTypeAnswer {
		if len(answers) == 0 {
			return false
		}
		answer, ok := answers[0].(string)
		if !ok {
			return false
		}

		for _, choice := range q.Choices {
			if choice.IsCorrect && choice.Content == answer {
				return true
			}
		}
	} else if q.Type == QuestionTypePoll || q.Type == QuestionWordCloud || q.Type == QuestionOpenEnded {
		return true
	}

	return false
}

func (q Question) CalculatePoints(time float64, correctRatio float64) uint {
	limitTime := float64(q.LimitTime) + AcceptableDelayTime
	if time < 0 || time > limitTime || correctRatio <= 0 || correctRatio > 1 {
		return 0
	}

	if time > float64(q.LimitTime) {
		time = float64(q.LimitTime)
	}

	points := uint(float64(q.Points) * (2 - time/float64(q.LimitTime)) * correctRatio)
	return points
}

func (ques *Question) UnmarshalBinary(data []byte) error {
	ques.Status = QuestionStatusWaiting
	if ques.StartedAt != nil {
		ques.Status = QuestionStatusPlaying
	}
	if ques.EndedAt != nil {
		ques.Status = QuestionStatusEnded
	}

	return json.Unmarshal(data, ques)
}

func (ques *Question) MarshalBinary() (data []byte, err error) {
	ques.Status = QuestionStatusWaiting
	if ques.StartedAt != nil {
		ques.Status = QuestionStatusPlaying
	}
	if ques.EndedAt != nil {
		ques.Status = QuestionStatusEnded
	}

	return json.Marshal(ques)
}

func (m *QuestionChoice) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *QuestionChoice) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func IsQuestionTypeGroupMultipleChoice(questionType QuestionType) bool {
	return questionType == QuestionTypeMultipleChoice || questionType == QuestionTypeTrueFalse
}

func IsQuestionTypeGroupTypeAnswer(questionType QuestionType) bool {
	return questionType == QuestionTypeTypeAnswer || questionType == QuestionOpenEnded || questionType == QuestionTypePoll || questionType == QuestionWordCloud
}
