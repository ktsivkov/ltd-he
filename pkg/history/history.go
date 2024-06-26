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
	GameId  int                `json:"gameId"`
	IsLast  bool               `json:"isLast"`
	Account *player.Player     `json:"account"`
	*game_stats.Stats
}

type Update struct {
	History History
	Err     error
}
