package history

import (
	"fmt"
	"time"

	"github.com/ktsivkov/ltd-he/pkg/game_stats"
)

type AppendRequest struct {
	Elo int  `json:"elo"`
	Mvp bool `json:"mvp"`
}

func (r *AppendRequest) Validate(currentElo int) error {
	if r.Elo == currentElo {
		return fmt.Errorf("elo cannot be the same as the current one %d", currentElo)
	}

	if r.Elo < game_stats.MinElo {
		return fmt.Errorf("elo cannot be less than %d", game_stats.MinElo)
	}

	if r.Elo > game_stats.MaxElo {
		return fmt.Errorf("elo cannot be less than %d", game_stats.MaxElo)
	}

	return nil
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

func (r *InsertRequest) Validate() error {
	if r.TotalGames < 1 {
		return fmt.Errorf("total games cannot be more than 1")
	}

	if r.Wins < 0 {
		return fmt.Errorf("wins cannot be less than 0")
	}

	if r.Wins > r.TotalGames {
		return fmt.Errorf("wins cannot be more than total games")
	}

	if r.GamesLeftEarly < 0 {
		return fmt.Errorf("games left early cannot be less than 0")
	}

	if r.GamesLeftEarly > r.TotalGames-r.Wins {
		return fmt.Errorf("games left early cannot be more than total games - wins")
	}

	if r.HighestWinStreak < 0 {
		return fmt.Errorf("highest win streak cannot be less than 0")
	}

	if r.HighestWinStreak > r.Wins {
		return fmt.Errorf("highest wins streak cannot be more than wins")
	}

	if r.WinsStreak < 0 {
		return fmt.Errorf("wins streak cannot be less than 0")
	}

	if r.WinsStreak > r.HighestWinStreak {
		return fmt.Errorf("wins streak cannot be more than highest wins streak")
	}

	if r.Mvp < 0 {
		return fmt.Errorf("mvp cannot be less than 0")
	}

	if r.Mvp > r.TotalGames-r.GamesLeftEarly {
		return fmt.Errorf("mvp cannot be more than total games - games left early")
	}

	if r.Elo < game_stats.MinElo {
		return fmt.Errorf("elo cannot be less than %d", game_stats.MinElo)
	}

	if r.Elo > game_stats.MaxElo {
		return fmt.Errorf("elo cannot be less than %d", game_stats.MaxElo)
	}

	if r.Timestamp.After(time.Now()) {
		return fmt.Errorf("timestamp cannot be in the future")
	}

	return nil
}
