package report

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/ktsivkov/ltd-he/pkg/player"
	"github.com/ktsivkov/ltd-he/pkg/utils"
)

const (
	reportFileName                  = "BnetLogs.pld"
	lastGameIdTokenContainerPattern = `BlzSetAbilityTooltip\('A017', "([^"]+)", 0\)`
	lastGameIdPattern               = `^(\d+)-`
	tokenPattern                    = `-([^"]+)`
	gameIdHolder                    = "__GAME_ID__"
	tokenHolder                     = "__TOKEN__"
)

//go:embed "BnetLogs.template"
var bnetLogsTemplate string

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

func (s *Service) Load(_ context.Context, p *player.Player) (*Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	payload, err := s.storageDriver.ReadFile(p.ReportFilePathAbsolute, reportFileName)
	if err != nil {
		return nil, fmt.Errorf("error reading report file: %w", err)
	}

	content := string(payload)
	report := &Report{
		Payload: payload,
	}
	var container string
	container, err = utils.RegexpMatchString(regexp.MustCompile(lastGameIdTokenContainerPattern), content)
	if err != nil {
		return nil, fmt.Errorf("could not match total games: %w", err)
	}

	report.LastGameId, err = utils.RegexpMatchInt(regexp.MustCompile(lastGameIdPattern), container)
	if err != nil {
		return nil, fmt.Errorf("could not match total games: %w", err)
	}

	report.Token, err = utils.RegexpMatchString(regexp.MustCompile(tokenPattern), container)
	if err != nil {
		return nil, fmt.Errorf("could not match total games: %w", err)
	}

	return report, nil
}

func (s *Service) Update(_ context.Context, p *player.Player, gameId int, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.storageDriver.WriteFile([]byte(strings.ReplaceAll(
		strings.ReplaceAll(
			bnetLogsTemplate,
			gameIdHolder,
			fmt.Sprintf("%d", gameId),
		),
		tokenHolder,
		token,
	)), p.ReportFilePathAbsolute, reportFileName); err != nil {
		return fmt.Errorf("could not edit report file: %w", err)
	}

	return nil
}

func reportFile(p *player.Player) string {
	return filepath.Join()
}
