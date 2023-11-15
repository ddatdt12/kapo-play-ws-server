package models

import (
	"encoding/json"
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/internal/utils/types"
)

type GameType string
type GameStatus string

const (
	GameTypeQuizGame    GameType = "quiz_game"
	GameTypeInteractive GameType = "interactive_game"
)

const (
	GameStatusWaiting GameStatus = "waiting"
	GameStatusPlaying GameStatus = "playing"
	GameStatusEnded   GameStatus = "ended"
)

type Game struct {
	ID         uint
	Code       string
	Name       string
	Status     GameStatus
	StartTime  types.NullTime
	EndTime    types.NullTime
	TemplateID uint
	Template   Template
	Type       GameType
	HostID     uint
	Settings   GameSettings
	CreatedAt  time.Time
}

type GameSettings struct {
	RandomizeQuestions bool `json:"randomize_questions"`
	RandomizeAnswers   bool `json:"randomize_answers"`
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
