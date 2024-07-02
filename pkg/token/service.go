package token

const (
	eloMaxLvlRound = 9
	payloadSize    = 28
	yd             = 1
	uppercaseToken = false
	encodePlayer   = true
)

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) Token(player string, totalGames int, wins int, elo int, gamesLeftEarly int, winsStreak int, highestWinStreak int, mvp int, wasWin bool) string {
	token := &Token{
		player:           player,
		totalGames:       totalGames,
		wins:             wins,
		elo:              elo,
		gamesLeftEarly:   gamesLeftEarly,
		winsStreak:       winsStreak,
		highestWinStreak: highestWinStreak,
		mvp:              mvp,
		// Defaults
		eloMaxLvlRound: eloMaxLvlRound,
		payloadSize:    payloadSize,
		yd:             yd,
		wasWin:         wasWin,
		uppercaseToken: uppercaseToken,
		encodePlayer:   encodePlayer,
	}

	return token.Token()
}
