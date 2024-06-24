package history

import "github.com/ktsivkov/ltd-he/pkg/game_stats"

type GameHistory struct {
	Outcome game_stats.Outcome `json:"outcome"`
	EloDiff int                `json:"eloDiff"`
	Stats   *game_stats.Stats  `json:"stats"`
	Date    string             `json:"date"`
}
