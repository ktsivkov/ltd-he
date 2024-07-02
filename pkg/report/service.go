package report

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/ktsivkov/ltd-he/pkg/player"
	"github.com/ktsivkov/ltd-he/pkg/utils"
)

const (
	lastGameIdTokenContainerPattern = `BlzSetAbilityTooltip\('A017', "([^"]+)", 0\)`
	lastGameIdPattern               = `^(\d+)-`
	tokenPattern                    = `-([^"]+)`
)

const reportFileName = "BnetLogs.pld"

func NewService() *Service {
	return &Service{
		mu: &sync.Mutex{},
	}
}

type Service struct {
	mu *sync.Mutex
}

func (s *Service) Load(_ context.Context, p *player.Player) (*Report, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	reportFilePath := reportFile(p)
	payload, err := os.ReadFile(reportFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading report file: %w", err)
	}

	content := string(payload)
	report := &Report{
		File:    reportFilePath,
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

func (s *Service) Update(_ context.Context, p *player.Player, report *Report, gameId int, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.WriteFile(reportFile(p), bytes.ReplaceAll(
		report.Payload,
		[]byte(getLastGameIdToken(report.LastGameId, report.Token)),
		[]byte(getLastGameIdToken(gameId, token)),
	), os.ModePerm); err != nil {
		return fmt.Errorf("could not edit report file: %w", err)
	}

	return nil
}

func reportFile(p *player.Player) string {
	return filepath.Join(p.ReportFilePathAbsolute, reportFileName)
}

func getLastGameIdToken(gameId int, token string) string {
	return fmt.Sprintf("%d-%s", gameId, token)
}
