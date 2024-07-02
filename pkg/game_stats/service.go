package game_stats

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ktsivkov/ltd-he/pkg/player"
)

const (
	dataPldFile = "Data.pld"
	dataTxtFile = "Data.txt"
)

var GameFileNotFoundErr = errors.New("game file not found")

func NewService() *Service {
	return &Service{
		mu: &sync.Mutex{},
	}
}

type Service struct {
	mu *sync.Mutex
}

func (s *Service) Load(_ context.Context, p *player.Player, gameId int) (*Stats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.loadFile(p, getStatsFileName(gameId))
}

func (s *Service) Delete(_ context.Context, p *player.Player, gameId int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(p.LogsPathAbsolute, getStatsFileName(gameId))
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not delete file %s: %w", filePath, err)
	}

	return nil
}

func (s *Service) Rollback(_ context.Context, p *player.Player, g *Stats) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.storeFile(p.LogsPathAbsolute, dataPldFile, g.Payload); err != nil {
		return err
	}

	if err := s.storeFile(p.LogsPathAbsolute, dataTxtFile, g.Payload); err != nil {
		return err
	}

	return nil
}

func (s *Service) Insert(_ context.Context, p *player.Player, gameId int, totalGames int, wins int, elo int, totalLosses int, gamesLeftEarly int, winsStreak int, highestWinStreak int, mvp int, token string, timestamp time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats, err := s.loadFile(p, dataPldFile)
	if err != nil {
		return err
	}

	stats.GameId = gameId
	stats.TotalGames = totalGames
	stats.Wins = wins
	stats.Elo = elo
	stats.TotalLosses = totalLosses
	stats.GamesLeftEarly = gamesLeftEarly
	stats.WinsStreak = winsStreak
	stats.HighestWinStreak = highestWinStreak
	stats.Mvp = mvp
	stats.Token = token
	stats.Timestamp = timestamp
	stats.payloadUpdate()

	if err := s.storeFile(p.LogsPathAbsolute, getStatsFileName(stats.GameId), stats.Payload); err != nil {
		return err
	}

	if err := s.storeFile(p.LogsPathAbsolute, dataPldFile, stats.Payload); err != nil {
		return err
	}

	if err := s.storeFile(p.LogsPathAbsolute, dataTxtFile, stats.Payload); err != nil {
		return err
	}

	return nil
}

func (s *Service) loadFile(p *player.Player, file string) (*Stats, error) {
	fp := filepath.Join(p.LogsPathAbsolute, file)
	payload, err := os.ReadFile(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, GameFileNotFoundErr
		}

		return nil, fmt.Errorf("could not read game file: %w", err)
	}

	stats := &Stats{
		File:    fp,
		Payload: payload,
	}

	if err := stats.hydrate(); err != nil {
		return nil, stats.descriptiveError(fmt.Errorf("could not parse game file: %w", err))
	}

	if file == dataPldFile || file == dataTxtFile {
		stats.GameId = stats.TotalGames
		return stats, nil
	}

	stats.GameId, err = stats.gameId()
	if err != nil {
		return nil, stats.descriptiveError(fmt.Errorf("could not parse game file: %w", err))
	}

	return stats, nil
}

func (s *Service) storeFile(path string, file string, data []byte) error {
	if err := os.WriteFile(filepath.Join(path, file), data, os.ModePerm); err != nil {
		return fmt.Errorf("could not store %s file: %w", dataPldFile, err)
	}
	return nil
}

func getStatsFileName(gameId int) string {
	return fmt.Sprintf("%s%d%s", statsFilePrefix, gameId, statsFileSuffix)
}
