package game_stats

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ktsivkov/ltd-he/pkg/player"
)

const (
	dataPldFile     = "Data.pld"
	dataTxtFile     = "Data.txt"
	statsFileFormat = "DataBU%d.pld"
)

var GameFileNotFoundErr = errors.New("game file not found")

func NewService(storageDriver StorageDriver) *Service {
	return &Service{
		mu:            &sync.Mutex{},
		storageDriver: storageDriver,
	}
}

type Service struct {
	mu            *sync.Mutex
	storageDriver StorageDriver
}

func (s *Service) NewStats(player string, totalGames int, wins int, elo int, gamesLeftEarly int, winsStreak int, highestWinStreak int, mvp int, token string, timestamp time.Time, gameVersion string) *Stats {
	return &Stats{
		TotalGames:       totalGames,
		Wins:             wins,
		Elo:              elo,
		TotalLosses:      totalGames - wins - gamesLeftEarly,
		GamesLeftEarly:   gamesLeftEarly,
		WinsStreak:       winsStreak,
		HighestWinStreak: highestWinStreak,
		Mvp:              mvp,
		Token:            token,
		Player:           player,
		GameVersion:      gameVersion,
		Timestamp:        timestamp,
	}
}

func (s *Service) Load(_ context.Context, p *player.Player, gameId int) (*Stats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.loadFile(p, getStatsFileName(gameId))
}

func (s *Service) Delete(_ context.Context, p *player.Player, gameId int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.storageDriver.DeletePath(p.LogsPathAbsolute, getStatsFileName(gameId)); err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}

	return nil
}

func (s *Service) Rollback(_ context.Context, p *player.Player, g *Stats) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	payload := g.GenerateFileContents()
	if err := s.storeFile(p, dataPldFile, payload); err != nil {
		return err
	}

	if err := s.storeFile(p, dataTxtFile, payload); err != nil {
		return err
	}

	return nil
}

func (s *Service) Insert(_ context.Context, p *player.Player, stats *Stats) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	payload := stats.GenerateFileContents()

	if err := s.storeFile(p, getStatsFileName(stats.TotalGames), payload); err != nil {
		return err
	}

	if err := s.storeFile(p, dataPldFile, payload); err != nil {
		return err
	}

	if err := s.storeFile(p, dataTxtFile, payload); err != nil {
		return err
	}

	return nil
}

func (s *Service) ClearStats(_ context.Context, p *player.Player) error {
	if err := s.storageDriver.DeletePath(p.LogsPathAbsolute); err != nil {
		return fmt.Errorf("could not delete logs path: %w", err)
	}

	if err := s.storageDriver.CreateDir(p.LogsPathAbsolute); err != nil {
		return fmt.Errorf("could not create logs path: %w", err)
	}

	return nil
}

func (s *Service) loadFile(p *player.Player, file string) (*Stats, error) {
	payload, err := s.storageDriver.ReadFile(p.LogsPathAbsolute, file)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", file, err)
	}

	stats := &Stats{}

	if err := stats.ParseFileContents(payload); err != nil {
		return nil, stats.descriptiveError(fmt.Errorf("could not parse game file: %w", err))
	}

	return stats, nil
}

func (s *Service) storeFile(p *player.Player, file string, data []byte) error {
	if err := s.storageDriver.WriteFile(data, p.LogsPathAbsolute, file); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func getStatsFileName(gameId int) string {
	return fmt.Sprintf(statsFileFormat, gameId)
}
