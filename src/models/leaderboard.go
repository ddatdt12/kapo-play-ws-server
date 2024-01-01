package models

type Leaderboard struct {
	ID       uint               `json:"id"`
	GameID   uint               `json:"gameId"`
	GameCode string             `json:"gameCode"`
	Items    []*LeaderboardUser `json:"items"`
}

type LeaderboardUser struct {
	LeaderboardID uint   `json:"leaderboardId"`
	Username      string `json:"username"`
	User          User   `json:"user"`
	Points        uint   `json:"points"`
	Rank          uint   `json:"rank"`
}
