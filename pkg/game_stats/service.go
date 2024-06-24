package game_stats

import (
	"fmt"
	"github.com/ktsivkov/ltd-he/pkg/player"
	"os"
	"path/filepath"
)

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) LoadGameStats(p *player.Player, gameId int) (*Stats, error) {
	filePath := filepath.Join(p.LogsDirPath, getStatsFileName(gameId))
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

func getStatsFileName(lastGameId int) string {
	return fmt.Sprintf("%s%d%s", statsFilePrefix, lastGameId, statsFileSuffix)
}
