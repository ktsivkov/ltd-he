package history

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ktsivkov/ltd-he/pkg/game_stats"
	"github.com/ktsivkov/ltd-he/pkg/player"
	"github.com/ktsivkov/ltd-he/pkg/report"
	"github.com/ktsivkov/ltd-he/pkg/token"
)

const dateFormat string = "02-01-2006 15:04:05"

func NewService(reportService *report.Service, gameStatsService *game_stats.Service, tokenService *token.Service) *Service {
	return &Service{
		reportService:    reportService,
		gameStatsService: gameStatsService,
		tokenService:     tokenService,
	}
}

type Service struct {
	reportService    *report.Service
	gameStatsService *game_stats.Service
	tokenService     *token.Service
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
		return fmt.Errorf("could not rollback: target_game_id %d >= last_game_id %d", game.GameId, r.LastGameId)
	}

	if err := s.reportService.Update(ctx, game.Account, r, game.GameId, game.Token); err != nil {
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

func (s *Service) Insert(ctx context.Context, p *player.Player, req *InsertRequest) error {
	r, err := s.reportService.Load(ctx, p)
	if err != nil {
		return fmt.Errorf("could not load report: %w", err)
	}

	lastGame, err := s.gameStatsService.Load(ctx, p, r.LastGameId)
	if err != nil {
		return fmt.Errorf("could not load last game: %w", err)
	}

	if lastGame.Elo == req.Elo {
		return fmt.Errorf("cannot create game with the same elo as the last one")
	}

	if req.Elo > 3000 {
		return errors.New("elo > 3000 is not supported")
	}

	if req.Elo < 1000 {
		return errors.New("elo < 1000 is not supported")
	}

	gameId := r.LastGameId + 1
	totalGames := lastGame.TotalGames + 1
	wins := lastGame.Wins
	winsStreak := 0
	highestWinStreak := lastGame.HighestWinStreak
	totalLosses := lastGame.TotalLosses + 1
	gamesLeftEarly := lastGame.GamesLeftEarly
	wasWin := false
	if req.Elo > lastGame.Elo {
		wasWin = true
		totalLosses = lastGame.TotalLosses
		wins++
		winsStreak = lastGame.WinsStreak + 1
		if winsStreak > highestWinStreak {
			highestWinStreak = winsStreak
		}
	}
	mvp := lastGame.Mvp
	if req.Mvp {
		mvp++
	}

	now := time.Now()
	t, err := s.tokenService.Token(p.BattleTag, lastGame.TotalGames+1, wins, req.Elo, lastGame.GamesLeftEarly, winsStreak, highestWinStreak, mvp, now, wasWin)
	if err != nil {
		return fmt.Errorf("could not generate token: %w", err)
	}

	ok, err := s.tokenService.ValidateToken(p.BattleTag, t)
	if err != nil {
		return fmt.Errorf("could not validate generated token: %w", err)
	}

	if !ok {
		return fmt.Errorf("could generated a valid token: %w", err)
	}

	if err := s.gameStatsService.Insert(ctx, p, gameId, totalGames, wins, req.Elo, totalLosses, gamesLeftEarly, winsStreak, highestWinStreak, mvp, t, now); err != nil {
		return fmt.Errorf("could not insert game stats: %w", err)
	}

	if err := s.reportService.Update(ctx, p, r, gameId, t); err != nil {
		return fmt.Errorf("could not update report: %w", err)
	}

	return nil
}
