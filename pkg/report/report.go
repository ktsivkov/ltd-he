package report

import (
	"fmt"
	"ltd-he/pkg/utils"
	"os"
	"path/filepath"
	"regexp"
)

const (
	lastGameIdTokenContainerPattern = `BlzSetAbilityTooltip\('A017', "([^"]+)", 0\)`
	lastGameIdPattern               = `^(\d+)-`
	tokenPattern                    = `-([^"]+)`
)

const errorLogsDir = "War3ErrorLogs"
const reportLogsDir = "ReportLogs"
const reportFile = "BnetLogs.pld"

type Report struct {
	LastGameId int    `json:"lastGameId"`
	Token      string `json:"token"`
	Payload    []byte `json:"payload"`
}

func Load(path string, player string) (*Report, error) {
	reportFilePath := filepath.Join(path, errorLogsDir, reportLogsDir, player, reportFile)
	payload, err := os.ReadFile(reportFilePath)
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
