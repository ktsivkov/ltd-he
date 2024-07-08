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

func NewService(reportService *report.Service, gameStatsService *game_stats.Service, tokenService *token.Service, storageDriver StorageDriver) *Service {
	return &Service{
		reportService:    reportService,
		gameStatsService: gameStatsService,
		tokenService:     tokenService,
		storageDriver:    storageDriver,
	}
}

type Service struct {
	reportService    *report.Service
	gameStatsService *game_stats.Service
	tokenService     *token.Service
	storageDriver    StorageDriver
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

	if game.TotalGames >= r.LastGameId {
		return fmt.Errorf("could not rollback: target_game_id %d >= last_game_id %d", game.TotalGames, r.LastGameId)
	}

	if err := s.reportService.Update(ctx, game.Account, game.TotalGames, game.Token); err != nil {
		return fmt.Errorf("could not rollback report: %w", err)
	}

	if err := s.gameStatsService.Rollback(ctx, game.Account, game.Stats); err != nil {
		return fmt.Errorf("could not rollback game stats: %w", err)
	}

	for i := game.TotalGames + 1; i <= r.LastGameId; i++ {
		if err := s.gameStatsService.Delete(ctx, game.Account, i); err != nil {
			return fmt.Errorf("could not delete game stats: %w", err)
		}
	}

	return nil
}

func (s *Service) Insert(ctx context.Context, p *player.Player, req *InsertRequest) error {

	if err := req.Validate(); err != nil {
		return fmt.Errorf("could not validate request: %w", err)
	}

	t, err := s.tokenService.Token(p.BattleTag, req.TotalGames, req.Wins, req.Elo, req.GamesLeftEarly, req.WinsStreak, req.HighestWinStreak, req.Mvp, req.Timestamp, req.WinsStreak <= 1)
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

	stats := s.gameStatsService.NewStats(p.BattleTag, req.TotalGames, req.Wins, req.Elo, req.GamesLeftEarly, req.WinsStreak, req.HighestWinStreak, req.Mvp, t, req.Timestamp, game_stats.DefaultGameVersion)

	if err := s.gameStatsService.ClearStats(ctx, p); err != nil {
		return fmt.Errorf("could not insert game stats: %w", err)
	}

	if err := s.gameStatsService.Insert(ctx, p, stats); err != nil {
		return fmt.Errorf("could not insert game stats: %w", err)
	}

	if err := s.reportService.Update(ctx, p, stats.TotalGames, t); err != nil {
		return fmt.Errorf("could not update report: %w", err)
	}

	return nil
}

func (s *Service) Append(ctx context.Context, p *player.Player, req *AppendRequest) error {
	r, err := s.reportService.Load(ctx, p)
	if err != nil {
		return fmt.Errorf("could not load report: %w", err)
	}

	lastGame, err := s.gameStatsService.Load(ctx, p, r.LastGameId)
	if err != nil {
		return fmt.Errorf("could not load last game: %w", err)
	}

	if err := req.Validate(lastGame.Elo); err != nil {
		return fmt.Errorf("could not validate request: %w", err)
	}

	gameId := r.LastGameId + 1
	totalGames := lastGame.TotalGames + 1
	wins := lastGame.Wins
	winsStreak := 0
	highestWinStreak := lastGame.HighestWinStreak
	gamesLeftEarly := lastGame.GamesLeftEarly
	gameVersion := lastGame.GameVersion
	wasWin := false
	if req.Elo > lastGame.Elo {
		wasWin = true
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

	stats := s.gameStatsService.NewStats(p.BattleTag, totalGames, wins, req.Elo, gamesLeftEarly, winsStreak, highestWinStreak, mvp, t, now, gameVersion)

	if err := s.gameStatsService.Insert(ctx, p, stats); err != nil {
		return fmt.Errorf("could not insert game stats: %w", err)
	}

	if err := s.reportService.Update(ctx, p, gameId, t); err != nil {
		return fmt.Errorf("could not update report: %w", err)
	}

	return nil
}
