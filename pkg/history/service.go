package history

import (
	"context"
	"errors"
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
		return nil, fmt.Errorf("could not load report: %w", err)
	}

	var lastGameStats *game_stats.Stats

	history := make(History, r.LastGameId)
	lastElemId := r.LastGameId - 1
	for i := range r.LastGameId {
		stats, err := s.gameStatsService.Load(ctx, p, i+1)
		if err != nil {
			if errors.Is(err, game_stats.GameFileNotFoundErr) {
				continue
			}
			return nil, fmt.Errorf("could not load game game stats: %w", err)
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

func (s *Service) Rollback(ctx context.Context, game *GameHistory) error {
	r, err := s.reportService.Load(ctx, game.Account)
	if err != nil {
		return fmt.Errorf("could not load report: %w", err)
	}

	if game.GameId >= r.LastGameId {
		return fmt.Errorf("could rollback: target_game_id %d >= last_game_id %d", game.GameId, r.LastGameId)
	}

	if err := s.reportService.Rollback(ctx, game.Account, r, game.GameId, game.Token); err != nil {
		return fmt.Errorf("could not rollback report: %w", err)
	}

	if err := s.gameStatsService.Rollback(ctx, game.Account, game.Stats); err != nil {
		return fmt.Errorf("could not rollback game stats: %w", err)
	}

	for i := game.GameId + 1; i <= r.LastGameId; i++ {
		if err := s.gameStatsService.Delete(ctx, game.Account, i); err != nil {
			return fmt.Errorf("could not delete game stats: %w", err)
		}
	}

	return nil
}
