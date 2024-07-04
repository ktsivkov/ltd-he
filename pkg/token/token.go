package token

import (
	"time"
)

type tableContents struct {
	totalGames       int
	wins             int
	elo              int
	gamesLeftEarly   int
	winsStreak       int
	highestWinStreak int
	mvp              int
	eloMaxLvlRound   int
	payloadSize      int
	yd               int
	wasWin           bool
	timestamp        time.Time
}

func (t *tableContents) slice() []int {
	tc := make([]int, t.payloadSize)
	// Salts
	tc[1] = getRandomInt(131070, 4194500)
	//tc[1] = 131071
	tc[4] = getRandomInt(65000, 2097150)
	//tc[4] = 65001
	tc[7] = getRandomInt(131070, 2097150)
	//tc[7] = 131071
	tc[9] = getRandomInt(65000, 2097150)
	//tc[9] = 65001
	tc[16] = getRandomInt(65000, 1048150)
	//tc[16] = 65001
	tc[24] = getRandomInt(65000, 1048570)
	//tc[24] = 65001
	tc[27] = getRandomInt(65000, 1548570)
	//tc[24] = 65001
	// Values
	tc[0] = t.totalGames
	tc[2] = t.wins
	tc[3] = t.gamesLeftEarly
	tc[5] = t.elo
	tc[6] = t.winsStreak
	tc[8] = t.highestWinStreak
	tc[10] = t.mvp
	tc[11] = t.yd
	tc[12] = t.eloMaxLvlRound
	if t.wasWin { // should be a bool, my best guess so far is that it says if the game was win or lose
		tc[13] = 1
	} else {
		tc[13] = 2
	}
	tc[18] = tc[13] // Changed -- might be better to leave it as 0 always
	tc[19], tc[20], tc[21] = int(t.timestamp.Month()), t.timestamp.Day(), t.timestamp.Year()
	tc[22], tc[23] = t.timestamp.Hour(), t.timestamp.Minute()
	// Not actually needed since in Go primitives (such as int) cannot be nil
	//tc[14] = 0
	//tc[15] = 0
	//tc[17] = 0
	//tc[25] = 0
	//tc[26] = 0
	return tc
}
