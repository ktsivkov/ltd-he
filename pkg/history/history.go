package history

import "github.com/ktsivkov/ltd-he/pkg/game_stats"

type GameHistory struct {
	Outcome game_stats.Outcome `json:"outcome"`
	EloDiff int                `json:"eloDiff"`
	Date    string             `json:"date"`
	GameId  int                `json:"gameId"`
	*game_stats.Stats
}
