package models

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

type Answer struct {
	Values     []any       `json:"values"`
	QuestionID uint        `json:"questionId"`
	IsCorrect  bool        `json:"isCorrect"`
	AnswerTime float64     `json:"answerTime"`
	Points     int64       `json:"points"`
	User       *User       `json:"user"`
	GameID     uint        `json:"gameId"`
	Username   string      `json:"username"`
	AnsweredAt time.Time   `json:"answeredAt"`
	Report     *UserReport `json:"report"`
}

type UserReport struct {
	Points int64 `json:"points"`
	Rank   int   `json:"rank"`
}

func (a *Answer) ToJSON() string {
	answerJson, err := json.Marshal(a)
	if err != nil {
		log.Error().Err(err).Msg("Answer.ToJSON")
		return ""
	}
	return string(answerJson)
}

func (m *Answer) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Answer) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
