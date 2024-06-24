package stats

import (
	"fmt"
	"ltd-he/pkg/report"
	"ltd-he/pkg/utils"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	DataFileTxt = "Data.txt"
	DataFilePld = "Data.pld"

	logsDir                 = "Legion_TD_TeamOZE"
	statsFilePrefix         = "DataBU"
	statsFileSuffix         = ".pld"
	statsFilePattern        = `DataBU(\d+).pld`
	totalGamesPattern       = `Total Games: (\d+)`
	winsPattern             = `Wins: (\d+)`
	eloPattern              = `ELO: (\d+)`
	totalLosesPattern       = `Total Losses: (\d+)`
	gamesLeftEarlyPattern   = `Games Left early: (\d+)`
	winsStreakPattern       = `Wins Streak: (\d+)`
	highestWinStreakPattern = `Highest Win Streak: (\d+)`
	mvpPattern              = `MVP: (\d+)`
	playerPattern           = `Player: (\w+#\d+)`
	tokenPattern            = `BlzSetAbilityTooltip\('A017', "([^"]+)", 0\)`
	timestampPattern        = `Time Stamp: (\d{1,2})\/(\d{1,2})\/(\d{4}) - (\d{1,2}):(\d{1,2}):(\d{1,2})`
	timestampFormat         = "1/2/2006 - 15:04:05"

	defaultElo = 1500
)

type Stats struct {
	GameId           int
	File             string    `json:"file"`
	TotalGames       int       `json:"totalGames"`
	Wins             int       `json:"wins"`
	Elo              int       `json:"elo"`
	TotalLosses      int       `json:"totalLosses"`
	GamesLeftEarly   int       `json:"gamesLeftEarly"`
	WinsStreak       int       `json:"winsStreak"`
	HighestWinStreak int       `json:"highestWinStreak"`
	Mvp              int       `json:"mvp"`
	Token            string    `json:"token"`
	Player           string    `json:"player"`
	Timestamp        time.Time `json:"timestamp"`
	Payload          []byte    `json:"payload"`
}

func LoadFromReport(path string, player string, rp *report.Report) (*Stats, error) {
	return Load(path, player, getStatsFileName(rp.LastGameId))
}

func Load(path string, player string, file string) (*Stats, error) {
	filePath := filepath.Join(path, logsDir, player, file)
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

func (s *Stats) Validate(player string) error {
	if s.Player != player {
		return s.descriptiveError(fmt.Errorf("player name does not match with the file: expected %s, got %s", player, s.Player))
	}

	gameId, err := s.gameId()
	if err != nil {
		return s.descriptiveError(err)
	}

	if gameId != s.TotalGames {
		return s.descriptiveError(fmt.Errorf("total number of games does not match the gameId of the file: game=%d, total=%d", gameId, s.TotalGames))
	}

	return nil
}
func (s *Stats) hydrate() error {
	var (
		err     error
		content = string(s.Payload)
	)

	s.TotalGames, err = utils.RegexpMatchInt(regexp.MustCompile(totalGamesPattern), content)
	if err != nil {
		return fmt.Errorf("could not match total games: %w", err)
	}

	s.Wins, err = utils.RegexpMatchInt(regexp.MustCompile(winsPattern), content)
	if err != nil {
		return fmt.Errorf("could not match wins: %w", err)
	}

	s.Elo, err = utils.RegexpMatchInt(regexp.MustCompile(eloPattern), content)
	if err != nil {
		return fmt.Errorf("could not match elo: %w", err)
	}

	s.TotalLosses, err = utils.RegexpMatchInt(regexp.MustCompile(totalLosesPattern), content)
	if err != nil {
		return fmt.Errorf("could not match total losses: %w", err)
	}

	s.GamesLeftEarly, err = utils.RegexpMatchInt(regexp.MustCompile(gamesLeftEarlyPattern), content)
	if err != nil {
		return fmt.Errorf("could not match games left early: %w", err)
	}

	s.WinsStreak, err = utils.RegexpMatchInt(regexp.MustCompile(winsStreakPattern), content)
	if err != nil {
		return fmt.Errorf("could not match wins streak: %w", err)
	}

	s.HighestWinStreak, err = utils.RegexpMatchInt(regexp.MustCompile(highestWinStreakPattern), content)
	if err != nil {
		return fmt.Errorf("could not match highest win streak: %w", err)
	}

	s.Mvp, err = utils.RegexpMatchInt(regexp.MustCompile(mvpPattern), content)
	if err != nil {
		return fmt.Errorf("could not match mvp: %w", err)
	}

	s.Player, err = utils.RegexpMatchString(regexp.MustCompile(playerPattern), content)
	if err != nil {
		return fmt.Errorf("could not match player: %w", err)
	}

	s.Token, err = utils.RegexpMatchString(regexp.MustCompile(tokenPattern), content)
	if err != nil {
		return fmt.Errorf("could not match token: %w", err)
	}

	var timestampUnits []int
	timestampUnits, err = utils.RegexpMatchAllInt(regexp.MustCompile(timestampPattern), content)
	if err != nil {
		return fmt.Errorf("could not match timestamp: %w", err)
	}

	s.Timestamp, err = parseTimestamp(timestampUnits)
	if err != nil {
		return fmt.Errorf("could not parse timestamp: %w", err)
	}

	return nil
}

func (s *Stats) descriptiveError(err error) error {
	return fmt.Errorf("game=%d, file=%s: %w", s.TotalGames, s.File, err)
}

func (s *Stats) gameId() (int, error) {
	gameId, err := utils.RegexpMatchInt(regexp.MustCompile(statsFilePattern), s.File)
	if err != nil {
		return 0, fmt.Errorf("could not match game id: %w", err)
	}

	return gameId, nil
}

func getStatsFileName(lastGameId int) string {
	return fmt.Sprintf("%s%d%s", statsFilePrefix, lastGameId, statsFileSuffix)
}

func isStatsFilename(filename string) bool {
	return regexp.MustCompile(statsFilePattern).MatchString(filename)
}

func parseTimestamp(timestampUnits []int) (time.Time, error) {
	if len(timestampUnits) != 6 {
		return time.Time{}, fmt.Errorf("invalid number of timestamp units: %d", len(timestampUnits))
	}

	return time.Date(timestampUnits[2], time.Month(timestampUnits[0]), timestampUnits[1], timestampUnits[3], timestampUnits[4], timestampUnits[5], 0, time.Local), nil
}

func GetDefaultGameStats() *Stats {
	return &Stats{
		File:             "",
		TotalGames:       0,
		Wins:             0,
		Elo:              defaultElo,
		TotalLosses:      0,
		GamesLeftEarly:   0,
		WinsStreak:       0,
		HighestWinStreak: 0,
		Mvp:              0,
		Token:            "",
		Player:           "",
		Timestamp:        time.Unix(0, 0).UTC(),
		Payload:          nil,
	}
}
