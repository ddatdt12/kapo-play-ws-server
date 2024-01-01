package models

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type GameState struct {
	//Current Question Offset
	CurrentQuestionOffset int `json:"current_question_offset"`
	// GameStage when game is playing
	GameStage GameStage `json:"game_stage"`
	// Status
	Status GameStatus `json:"status"`
	// Answer
	Answer *Answer `json:"answer"`
	//Question
	Question *Question `json:"question"`

	// Game State Changed event function
	OnGameStateChanged func(gameState *GameState) `json:"-"`
	// GameStatusChanged event function
	OnGameStatusChanged func(newValue GameStatus, oldValue GameStatus, gameState *GameState) `json:"-"`
}

func NewGameState() *GameState {
	return &GameState{
		CurrentQuestionOffset: -1,
		GameStage:             GameStageNil,
		Status:                GameStatusWaiting,
		Answer:                nil,
		Question:              nil,
	}
}

func (gameState *GameState) Start() {
	gameState.Status = GameStatusPlaying
	if gameState.OnGameStateChanged != nil {
		gameState.OnGameStateChanged(gameState)
	}
}
func (gameState *GameState) SetGameStage(gameStage GameStage) {
	gameState.GameStage = gameStage
	if gameState.OnGameStateChanged != nil {
		gameState.OnGameStateChanged(gameState)
	}
}

func (gameState *GameState) SetStatus(status GameStatus) {
	if gameState.Status == status {
		return
	}

	if gameState.OnGameStatusChanged != nil {
		log.Info().Interface("OnGameStatusChanged", gameState).Msg("gameState")
		gameState.OnGameStatusChanged(status, gameState.Status, gameState)
	}
	gameState.Status = status
	if gameState.OnGameStateChanged != nil {
		gameState.OnGameStateChanged(gameState)
	}
}

func (gameState *GameState) SetQuestion(question *Question) {
	gameState.Question = question
	if gameState.OnGameStateChanged != nil {
		gameState.OnGameStateChanged(gameState)
	}
}

func (gameState *GameState) SetAnswer(answer *Answer) {
	gameState.Answer = answer
	if gameState.OnGameStateChanged != nil {
		gameState.OnGameStateChanged(gameState)
	}
}

func (gameState *GameState) SetCurrentQuestionOffset(currentQuestionOffset int) {
	gameState.CurrentQuestionOffset = currentQuestionOffset
	if gameState.OnGameStateChanged != nil {
		gameState.OnGameStateChanged(gameState)
	}
}
func (gameState *GameState) NextQuestion() {
	gameState.CurrentQuestionOffset++
	if gameState.OnGameStateChanged != nil {
		gameState.OnGameStateChanged(gameState)
	}
}

func (gameState *GameState) Reset() {
	gameState.CurrentQuestionOffset = 0
	gameState.GameStage = GameStageNil
	gameState.Status = GameStatusWaiting
	gameState.Answer = nil
	gameState.Question = nil
	if gameState.OnGameStateChanged != nil {
		gameState.OnGameStateChanged(gameState)
	}
}

func (g GameState) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func (g *GameState) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, g)
}
