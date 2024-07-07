package history

import "time"

type AppendRequest struct {
	Elo int  `json:"elo"`
	Mvp bool `json:"mvp"`
}

type InsertRequest struct {
	TotalGames       int       `json:"totalGames"`
	Wins             int       `json:"wins"`
	Elo              int       `json:"elo"`
	GamesLeftEarly   int       `json:"gamesLeftEarly"`
	WinsStreak       int       `json:"winsStreak"`
	HighestWinStreak int       `json:"highestWinStreak"`
	Mvp              int       `json:"mvp"`
	Timestamp        time.Time `json:"timestamp"`
}
