package models

import (
	"encoding/json"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/types"
)

type GameType string
type GameStatus string
type GameStage string

const (
	GameTypeQuizGame    GameType = "quiz_game"
	GameTypeInteractive GameType = "interactive_game"
)

const (
	GameStatusWaiting GameStatus = "waiting"
	GameStatusPlaying GameStatus = "playing"
	GameStatusEnded   GameStatus = "ended"
)

var (
	GameStageShowQuestion GameStage = "stage/show_question"
	GameStageShowAnswer   GameStage = "stage/show_answer"
)

type Game struct {
	ID             uint               `json:"id"`
	Code           string             `json:"code"`
	Name           string             `json:"name"`
	Status         GameStatus         `json:"status"`
	StartTime      types.NullableTime `json:"startTime"`
	EndTime        types.NullableTime `json:"endTime"`
	TemplateID     uint               `json:"templateId"`
	Template       Template           `json:"template"`
	Type           GameType           `json:"type"`
	HostID         uint               `json:"hostId"`
	Host           User               `json:"host"`
	Settings       GameSettings       `json:"settings"`
	TotalQuestions int64              `json:"totalQuestions"`
	CreatedAt      time.Time          `json:"createdAt"`
}

type GameSettings struct {
	RandomizeQuestions bool `json:"randomizeQuestions"`
	RandomizeAnswers   bool `json:"randomizeAnswers"`
}

func (s GameType) IsValid() bool {
	switch s {
	case GameTypeQuizGame, GameTypeInteractive:
		return true
	}

	return false
}

func (m *Game) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Game) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
