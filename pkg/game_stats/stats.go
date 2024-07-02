package game_stats

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/ktsivkov/ltd-he/pkg/utils"
)

type Outcome string

const (
	OutcomeLeave Outcome = "LEAVE"
	OutcomeWin   Outcome = "WIN"
	OutcomeLoss  Outcome = "LOSS"
	OutcomeDraw  Outcome = "DRAW"

	statsFilePrefix  = "DataBU"
	statsFileSuffix  = ".pld"
	statsFilePattern = `DataBU(\d+).pld`

	// Payload Patterns
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
	// Payload Formats
	totalGamesFormat       = "Total Games: %d"
	winsFormat             = "Wins: %d"
	eloFormat              = "ELO: %d"
	totalLosesFormat       = "Total Losses: %d"
	gamesLeftEarlyFormat   = "Games Left early: %d"
	winsStreakFormat       = "Wins Streak: %d"
	highestWinStreakFormat = "Highest Win Streak: %d"
	mvpFormat              = "MVP: %d"
	playerFormat           = "Player: %s"
	tokenFormat            = `BlzSetAbilityTooltip('A017', "%s", 0)`
	timestampFormat        = "Time Stamp: %d/%d/%d - %d:%d:%d"

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

func (s *Stats) payloadUpdate() {
	s.Payload = regexp.MustCompile(totalGamesPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(totalGamesFormat, s.TotalGames)))
	s.Payload = regexp.MustCompile(winsPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(winsFormat, s.Wins)))
	s.Payload = regexp.MustCompile(eloPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(eloFormat, s.Elo)))
	s.Payload = regexp.MustCompile(totalLosesPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(totalLosesFormat, s.TotalLosses)))
	s.Payload = regexp.MustCompile(gamesLeftEarlyPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(gamesLeftEarlyFormat, s.GamesLeftEarly)))
	s.Payload = regexp.MustCompile(winsStreakPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(winsStreakFormat, s.WinsStreak)))
	s.Payload = regexp.MustCompile(highestWinStreakPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(highestWinStreakFormat, s.HighestWinStreak)))
	s.Payload = regexp.MustCompile(mvpPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(mvpFormat, s.Mvp)))
	s.Payload = regexp.MustCompile(playerPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(playerFormat, s.Player)))
	s.Payload = regexp.MustCompile(tokenPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(tokenFormat, s.Token)))
	s.Payload = regexp.MustCompile(timestampPattern).ReplaceAll(s.Payload, []byte(fmt.Sprintf(timestampFormat, s.Timestamp.Month(), s.Timestamp.Day(), s.Timestamp.Year(), s.Timestamp.Hour(), s.Timestamp.Minute(), s.Timestamp.Second())))
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

func (s *Stats) Outcome(lastGame *Stats) Outcome {
	if lastGame == nil {
		lastGame = getDefaultGameStats()
	}

	if s.GamesLeftEarly > lastGame.GamesLeftEarly {
		return OutcomeLeave
	}
	if s.Elo > lastGame.Elo {
		return OutcomeWin
	}
	if s.Elo < lastGame.Elo {
		return OutcomeLoss
	}

	return OutcomeDraw
}

func (s *Stats) EloDiff(lastGame *Stats) int {
	if lastGame == nil {
		lastGame = getDefaultGameStats()
	}

	return int(math.Abs(float64(s.Elo - lastGame.Elo)))
}

func getDefaultGameStats() *Stats {
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
