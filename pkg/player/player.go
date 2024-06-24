package player

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const (
	logsDir = "Legion_TD_TeamOZE"

	playerDirPattern = `(\w+#\d+)`
)

type Player struct {
	BattleTag   string `json:"battleTag"`
	LogsDirPath string `json:"logsDirPath"`
}

func NewService(customMapDataPath string) *Service {
	return &Service{
		customMapDataPath: customMapDataPath,
	}
}

type Service struct {
	customMapDataPath string
}

func (s *Service) GetAll() ([]*Player, error) {
	logsDirPath := filepath.Join(s.customMapDataPath, logsDir)
	items, err := os.ReadDir(logsDirPath)
	if err != nil {
		return nil, fmt.Errorf("could not read logs directory: %w", err)
	}

	players := make([]*Player, 0, len(items)-1) // We know that there will be at least one garbage file (UserSettings.pld)
	for _, item := range items {
		if item.IsDir() && isBattleTag(item.Name()) {
			players = append(players, &Player{
				BattleTag:   item.Name(),
				LogsDirPath: filepath.Join(logsDirPath, item.Name()),
			})
		}
	}

	return players, nil
}

func isBattleTag(path string) bool {
	return regexp.MustCompile(playerDirPattern).MatchString(path)
}
