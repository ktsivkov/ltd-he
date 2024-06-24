package history

import (
	"fmt"
	"github.com/ktsivkov/ltd-he/pkg/game_stats"
	"github.com/ktsivkov/ltd-he/pkg/player"
	"github.com/ktsivkov/ltd-he/pkg/report"
)

const dateFormat string = "02-01-2006 15:04:05"

func NewService(reportService *report.Service, gameStatsService *game_stats.Service) *Service {
	return &Service{
		reportService:    reportService,
		gameStatsService: gameStatsService,
	}
}

type Service struct {
	reportService    *report.Service
	gameStatsService *game_stats.Service
}

func (s *Service) LoadHistory(p *player.Player) ([]*GameHistory, error) {
	r, err := s.reportService.Load(p)
	if err != nil {
		return nil, fmt.Errorf("could not load history: %v", err)
	}

	var lastGameStats *game_stats.Stats

	history := make([]*GameHistory, r.LastGameId)
	lastElemId := r.LastGameId - 1
	for i := range r.LastGameId {
		stats, err := s.gameStatsService.LoadGameStats(p, i+1)
		if err != nil {
			return nil, fmt.Errorf("could not load game game_stats: %v", err)
		}

		history[lastElemId-i] = &GameHistory{
			GameId:  i + 1,
			Outcome: stats.Outcome(lastGameStats),
			EloDiff: stats.EloDiff(lastGameStats),
			Stats:   stats,
			Date:    stats.Timestamp.Format(dateFormat),
		}

		lastGameStats = stats
	}

	return history, err
}
