package stats

import (
	"fmt"
	"io"
	"ltd-he/pkg/report"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type Outcome string

const OutcomeLeave Outcome = "LEAVE"
const OutcomeWin Outcome = "WIN"
const OutcomeLoss Outcome = "LOSS"
const OutcomeDraw Outcome = "DRAW"

const dateFormat = "02-01-2006 15:04:05"

type History struct {
	Games []*Stats `json:"games"`
}

func LoadHistory(path string, player string) (*History, error) {
	files, err := os.ReadDir(filepath.Join(path, logsDir, player))
	if err != nil {
		return nil, fmt.Errorf("could not read log directory: %w", err)
	}

	statsFiles := make([]string, 0, len(files)-2) // In the directory we know for sure that we will have at least two non-relatable files Data.pld and Data.txt
	for _, file := range files {
		if isStatsFilename(file.Name()) {
			if file.IsDir() {
				return nil, fmt.Errorf("directory conflicting stats file found: %s", file.Name())
			}
			statsFiles = append(statsFiles, file.Name())
		}
	}

	games := make([]*Stats, len(statsFiles))
	for i, file := range statsFiles {
		games[i], err = Load(path, player, file)
		if err != nil {
			return nil, fmt.Errorf("could not load history: %w", err)
		}
	}

	slices.SortFunc(games, func(a, b *Stats) int {
		return a.TotalGames - b.TotalGames
	})

	return &History{
		Games: games,
	}, nil
}

func (h *History) Validate(player string, rp *report.Report) []error {
	errors := make([]error, 0, len(h.Games))

	latestGameId := 0
	for _, game := range h.Games {
		if err := game.Validate(player); err != nil {
			errors = append(errors, err)
		}

		if game.TotalGames > latestGameId {
			latestGameId = game.TotalGames
		}
	}

	if rp.LastGameId != latestGameId {
		errors = append(errors, fmt.Errorf("%w: %w", DataValidationError, fmt.Errorf("last found game mismatch report: found=%d != report=%d", latestGameId, rp.LastGameId)))
	}

	if len(h.Games) != rp.LastGameId {
		errors = append(errors, fmt.Errorf("%w: %w", DataValidationError, fmt.Errorf("total found games mismatch report: total=%d != report=%d", latestGameId, rp.LastGameId)))
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (h *History) Presentation() HistoryPresentation {
	lastGame := GetDefaultGameStats()

	res := make(HistoryPresentation, len(h.Games))
	for i, game := range h.Games {
		res[i] = &HistoryGamePresentation{
			Game:     game,
			LastGame: lastGame,
		}
		lastGame = game
	}
	return res
}

type HistoryPresentation []*HistoryGamePresentation

func (h HistoryPresentation) CliPresent(w io.Writer) error {
	const (
		winEloSign  = "+"
		lossEloSign = "-"
		drawEloSign = "~"
	)

	maxGameNumberWidth := len(strconv.Itoa(h[len(h)-1].Game.GameId))

	for _, p := range h {
		outcome := p.Outcome()

		var eloSign string
		switch outcome {
		case OutcomeWin:
			eloSign = winEloSign
		case OutcomeLoss, OutcomeLeave:
			eloSign = lossEloSign
		default:
			eloSign = drawEloSign
		}

		if _, err := fmt.Fprintln(w, fmt.Sprintf("[%s] Game: %*d Elo: %d (%s%02d) %s", p.Game.Timestamp.Format(dateFormat), maxGameNumberWidth, p.Game.GameId, p.Game.Elo, eloSign, p.EloDiff(), outcome)); err != nil {
			return fmt.Errorf("could not write cli output: %w", err)
		}
	}

	return nil
}

type HistoryGamePresentation struct {
	Game     *Stats
	LastGame *Stats
}

func (p *HistoryGamePresentation) Outcome() Outcome {
	if p.Game.GamesLeftEarly > p.LastGame.GamesLeftEarly {
		return OutcomeLeave
	}
	if p.Game.Elo > p.LastGame.Elo {
		return OutcomeWin
	}
	if p.Game.Elo < p.LastGame.Elo {
		return OutcomeLoss
	}

	return OutcomeDraw
}

func (p *HistoryGamePresentation) EloDiff() int {
	return int(math.Abs(float64(p.Game.Elo - p.LastGame.Elo)))
}
