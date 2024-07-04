package history

import (
	"github.com/ktsivkov/ltd-he/pkg/game_stats"
	"github.com/ktsivkov/ltd-he/pkg/player"
)

type History []*GameHistory

type GameHistory struct {
	Outcome game_stats.Outcome `json:"outcome"`
	EloDiff int                `json:"eloDiff"`
	Date    string             `json:"date"`
	IsLast  bool               `json:"isLast"`
	Account *player.Player     `json:"account"`
	*game_stats.Stats
}
