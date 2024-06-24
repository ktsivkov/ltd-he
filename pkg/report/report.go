package report

import (
	"fmt"
	"github.com/ktsivkov/ltd-he/pkg/player"
	"github.com/ktsivkov/ltd-he/pkg/utils"
	"os"
	"path/filepath"
	"regexp"
)

const (
	lastGameIdTokenContainerPattern = `BlzSetAbilityTooltip\('A017', "([^"]+)", 0\)`
	lastGameIdPattern               = `^(\d+)-`
	tokenPattern                    = `-([^"]+)`
)

const reportFile = "BnetLogs.pld"

func NewService() *Service {
	return &Service{}
}

type Report struct {
	File       string `json:"file"`
	LastGameId int    `json:"lastGameId"`
	Token      string `json:"token"`
	Payload    []byte `json:"payload"`
}

type Service struct{}

func (s *Service) Load(p *player.Player) (*Report, error) {
	reportFilePath := filepath.Join(p.ReportFilePath, reportFile)
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
