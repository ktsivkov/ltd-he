package player

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

const (
	errorLogsDir     = "War3ErrorLogs"
	reportLogsDir    = "ReportLogs"
	logsDir          = "Legion_TD_TeamOZE"
	customMapDataDir = "CustomMapData"

	playerDirPattern = `(\w+#\d+)`
)

func NewService(wc3Path string) *Service {
	return &Service{
		wc3Path: wc3Path,
		mu:      &sync.Mutex{},
	}
}

type Service struct {
	wc3Path string
	mu      *sync.Mutex
}

func (s *Service) LoadAll(_ context.Context) ([]*Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	logsPathRelative := filepath.Join(customMapDataDir, logsDir)
	logsPathAbsolute := filepath.Join(s.wc3Path, logsPathRelative)
	items, err := os.ReadDir(logsPathAbsolute)
	if err != nil {
		return nil, fmt.Errorf("could not read logs directory: %w", err)
	}

	players := make([]*Player, 0, len(items)-1) // We know that there will be at least one garbage file (UserSettings.pld)
	for _, item := range items {
		if item.IsDir() && isBattleTag(item.Name()) {
			reportPathRelative := filepath.Join(customMapDataDir, errorLogsDir, reportLogsDir, item.Name())
			players = append(players, &Player{
				BattleTag:              item.Name(),
				LogsPathAbsolute:       filepath.Join(logsPathAbsolute, item.Name()),
				LogsPathRelative:       filepath.Join(logsPathRelative, item.Name()),
				ReportFilePathAbsolute: filepath.Join(s.wc3Path, reportPathRelative),
				ReportFilePathRelative: reportPathRelative,
			})
		}
	}

	return players, nil
}

func isBattleTag(path string) bool {
	return regexp.MustCompile(playerDirPattern).MatchString(path)
}
