package game_stats

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ktsivkov/ltd-he/pkg/player"
)

const (
	dataPldFile = "Data.pld"
	dataTxtFile = "Data.txt"
)

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) Load(_ context.Context, p *player.Player, gameId int) (*Stats, error) {
	filePath := filepath.Join(p.LogsPathAbsolute, getStatsFileName(gameId))
	payload, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", DataLoadingError, err)
	}

	stats := &Stats{
		File:    filePath,
		Payload: payload,
	}

	if err := stats.hydrate(); err != nil {
		return nil, stats.descriptiveError(fmt.Errorf("%w: %w", DataParsingError, err))
	}

	stats.GameId, err = stats.gameId()
	if err != nil {
		return nil, stats.descriptiveError(fmt.Errorf("%w: %w", DataParsingError, err))
	}

	return stats, nil
}

func (s *Service) Delete(_ context.Context, p *player.Player, gameId int) error {
	filePath := filepath.Join(p.LogsPathAbsolute, getStatsFileName(gameId))
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete file %s: %w", filePath, err)
	}

	return nil
}

func (s *Service) Rollback(_ context.Context, p *player.Player, g *Stats) error {
	if err := s.rollbackFile(p.LogsPathAbsolute, dataPldFile, g.Payload); err != nil {
		return err
	}

	if err := s.rollbackFile(p.LogsPathAbsolute, dataTxtFile, g.Payload); err != nil {
		return err
	}

	return nil
}

func (s *Service) rollbackFile(path string, file string, data []byte) error {
	if err := os.WriteFile(filepath.Join(path, file), data, os.ModePerm); err != nil {
		return fmt.Errorf("could not rollback %s file: %w", dataPldFile, err)
	}
	return nil
}

func getStatsFileName(lastGameId int) string {
	return fmt.Sprintf("%s%d%s", statsFilePrefix, lastGameId, statsFileSuffix)
}
