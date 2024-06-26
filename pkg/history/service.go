package history

import (
	"context"
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

func (s *Service) Load(ctx context.Context, p *player.Player) (History, error) {
	r, err := s.reportService.Load(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("could not load history: %v", err)
	}

	var lastGameStats *game_stats.Stats

	history := make(History, r.LastGameId)
	lastElemId := r.LastGameId - 1
	for i := range r.LastGameId {
		stats, err := s.gameStatsService.Load(ctx, p, i+1)
		if err != nil {
			return nil, fmt.Errorf("could not load game game stats: %v", err)
		}

		history[lastElemId-i] = &GameHistory{
			GameId:  i + 1,
			IsLast:  i == lastElemId,
			Outcome: stats.Outcome(lastGameStats),
			EloDiff: stats.EloDiff(lastGameStats),
			Date:    stats.Timestamp.Format(dateFormat),
			Account: p,
			Stats:   stats,
		}

		lastGameStats = stats
	}

	return history, err
}
