package game_stats

import (
	_ "embed"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/ktsivkov/ltd-he/pkg/utils"
)

//go:embed Data.template
var payloadTemplate string

type Outcome string

const (
	OutcomeLeave Outcome = "LEAVE"
	OutcomeWin   Outcome = "WIN"
	OutcomeLoss  Outcome = "LOSS"
	OutcomeDraw  Outcome = "DRAW"

	// payload Patterns
	totalGamesPattern       = `Total Games: (\d+)`
	winsPattern             = `Wins: (\d+)`
	eloPattern              = `ELO: (\d+)`
	totalLossesPattern      = `Total Losses: (\d+)`
	gamesLeftEarlyPattern   = `Games Left early: (\d+)`
	winsStreakPattern       = `Wins Streak: (\d+)`
	highestWinStreakPattern = `Highest Win Streak: (\d+)`
	mvpPattern              = `MVP: (\d+)`
	playerPattern           = `Player: (\w+#\d+)`
	tokenPattern            = `BlzSetAbilityTooltip\('A017', "([^"]+)", 0\)`
	timestampPattern        = `Time Stamp: (\d{1,2})\/(\d{1,2})\/(\d{4}) - (\d{1,2}):(\d{1,2}):(\d{1,2})`
	gameVersionPattern      = `LTD TeamOZE Game Version: (\w+?\.\w+)`
	// template formats
	totalGamesFormat       = "%d"
	winsFormat             = "%d"
	eloFormat              = "%d"
	totalLossesFormat      = "%d"
	gamesLeftEarlyFormat   = "%d"
	winsStreakFormat       = "%d"
	highestWinStreakFormat = "%d"
	mvpFormat              = "%d"
	playerFormat           = "%s"
	tokenFormat            = `%s`
	timestampFormat        = "%d/%d/%d - %d:%d:%d"
	gameVersionFormat      = "%s"
	// template holders
	totalGamesHolder       = "__TOTAL_GAMES__"
	timestampHolder        = "__TIMESTAMP__"
	winsHolder             = "__WINS__"
	eloHolder              = "__ELO__"
	totalLossesHolder      = "__TOTAL_LOSSES__"
	gamesLeftEarlyHolder   = "__GAMES_LEFT_EARLY__"
	winsStreakHolder       = "__WINS_STREAK__"
	highestWinStreakHolder = "__HIGHEST_WIN_STREAK__"
	mvpHolder              = "__MVP__"
	playerHolder           = "__PLAYER__"
	gameVersionHolder      = "__GAME_VERSION__"
	tokenHolder            = "__TOKEN__"

	defaultElo         = 1500
	DefaultGameVersion = "11.0g"
)

type Stats struct {
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
	GameVersion      string    `json:"gameVersion"`
	Timestamp        time.Time `json:"timestamp"`
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

func (s *Stats) GenerateFileContents() []byte {
	payload := strings.ReplaceAll(payloadTemplate, totalGamesHolder, fmt.Sprintf(totalGamesFormat, s.TotalGames))
	payload = strings.ReplaceAll(payload, totalGamesHolder, fmt.Sprintf(totalGamesFormat, s.TotalGames))
	payload = strings.ReplaceAll(payload, winsHolder, fmt.Sprintf(winsFormat, s.Wins))
	payload = strings.ReplaceAll(payload, eloHolder, fmt.Sprintf(eloFormat, s.Elo))
	payload = strings.ReplaceAll(payload, totalLossesHolder, fmt.Sprintf(totalLossesFormat, s.TotalLosses))
	payload = strings.ReplaceAll(payload, gamesLeftEarlyHolder, fmt.Sprintf(gamesLeftEarlyFormat, s.GamesLeftEarly))
	payload = strings.ReplaceAll(payload, winsStreakHolder, fmt.Sprintf(winsStreakFormat, s.WinsStreak))
	payload = strings.ReplaceAll(payload, highestWinStreakHolder, fmt.Sprintf(highestWinStreakFormat, s.HighestWinStreak))
	payload = strings.ReplaceAll(payload, mvpHolder, fmt.Sprintf(mvpFormat, s.Mvp))
	payload = strings.ReplaceAll(payload, playerHolder, fmt.Sprintf(playerFormat, s.Player))
	payload = strings.ReplaceAll(payload, gameVersionHolder, fmt.Sprintf(gameVersionFormat, s.GameVersion))
	payload = strings.ReplaceAll(payload, tokenHolder, fmt.Sprintf(tokenFormat, s.Token))
	payload = strings.ReplaceAll(payload, timestampHolder, fmt.Sprintf(timestampFormat, s.Timestamp.Month(), s.Timestamp.Day(), s.Timestamp.Year(), s.Timestamp.Hour(), s.Timestamp.Minute(), s.Timestamp.Second()))
	return []byte(payload)
}

func (s *Stats) ParseFileContents(payload []byte) error {
	var (
		err     error
		content = string(payload)
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

	s.TotalLosses, err = utils.RegexpMatchInt(regexp.MustCompile(totalLossesPattern), content)
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

	s.GameVersion, err = utils.RegexpMatchString(regexp.MustCompile(gameVersionPattern), content)
	if err != nil {
		return fmt.Errorf("could not match game version: %w", err)
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
	return fmt.Errorf("game=%d: %w", s.TotalGames, err)
}

func getDefaultGameStats() *Stats {
	return &Stats{
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
	}
}
