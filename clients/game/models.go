package game

import (
	"time"

	"github.com/ddatdt12/kapo-play-ws-server/src/models"
)

type GameRequest struct {
	ID        uint              `json:"id"`
	Code      string            `json:"code"`
	Status    models.GameStatus `json:"status"`
	StartTime *time.Time        `json:"startTime"`
	EndTime   *time.Time        `json:"endTime"`
}

type GameStateRequest struct {
	GameRequest
	Metadata map[string]interface{} `json:"metadata"`
}
