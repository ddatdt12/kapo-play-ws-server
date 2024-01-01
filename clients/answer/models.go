package answer

import "time"

type CreateAnswer struct {
	Values     []any      `json:"values"`
	QuestionID uint       `json:"questionId"`
	GameID     uint       `json:"gameId"`
	IsCorrect  bool       `json:"isCorrect"`
	AnswerTime float64    `json:"answerTime"`
	Points     int64      `json:"points"`
	Username   string     `json:"username"`
	AnsweredAt time.Time  `json:"answeredAt"`
	Report     UserReport `json:"report"`
}

type UserReport struct {
	Points      int64 `json:"points"`
	Rank        int   `json:"rank"`
	StreakCount int   `json:"streakCount"`
}
