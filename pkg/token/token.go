package token

import (
	"strconv"
	"strings"
	"time"
)

const encryptionKey = "lhECop@tFxnbvkdr,;)+V3L^56aR1_gW7u&qJzTASB*N}iU|fDm>c=Q[2y~ZKH!OMwP0(9s`.<%8G:${?Xe]#4j'YI"

type Token struct {
	player           string
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
	uppercaseToken   bool
	encodePlayer     bool
}

func (s *Token) Token() string {
	return formatToken(s.generateToken())
}

func (s *Token) tableContents() []int {
	tc := make([]int, s.payloadSize)
	tc[0] = s.totalGames
	tc[1] = getRandomInt(131070, 4194500)
	tc[2] = s.wins
	tc[3] = s.gamesLeftEarly
	tc[4] = getRandomInt(65000, 2097150)
	tc[5] = s.elo
	tc[6] = s.winsStreak
	tc[7] = getRandomInt(131070, 2097150)
	tc[8] = s.highestWinStreak
	tc[9] = getRandomInt(65000, 2097150)
	tc[10] = s.mvp
	tc[11] = s.yd
	tc[12] = s.eloMaxLvlRound
	if s.wasWin { // TODO should be bool
		tc[13] = 2
	} else {
		tc[13] = 1
	}
	tc[14] = 0
	tc[15] = 0
	tc[16] = getRandomInt(65000, 1048150)
	tc[17] = 0

	tc[18] = tc[13] // Changed
	now := time.Now()
	tc[19], tc[20], tc[21] = int(now.Month()), now.Day(), now.Year()
	tc[22], tc[23] = now.Hour(), now.Minute()
	tc[24] = getRandomInt(65000, 1048570)
	tc[25] = 0
	tc[26] = 0
	tc[27] = getRandomInt(65000, 1548570)
	return tc
}

func (s *Token) encryptionKey() string {
	if s.uppercaseToken {
		return strings.ToUpper(encryptionKey)
	}
	return encryptionKey
}

func (s *Token) calculateStringValue(inputString string) int {
	value := 0
	index := 0

	if s.encodePlayer {
		for index < len(s.player) {
			value += encodeCharacter(s.player[index : index+1])
			index++
		}
	}

	index = 0
	for index < len(inputString) {
		value += encodeCharacter(inputString[index : index+1])
		index++
	}

	return value
}

func (s *Token) generateToken() string {
	tc := s.tableContents()

	var _a, _b, _c, _d, _e, _f = 0, 0, 0, 0, 0, len(s.encryptionKey())
	var _g = make([]int, 100)
	var _h, _i = "", ""
	var _j, _k = 0, 1000000
	var _l = "0123456789"

	_a = 0
	for {
		_a++
		if _a > s.payloadSize {
			break
		}
		_h += strconv.Itoa(tc[_a-1]) + "-"
	}
	_h += strconv.Itoa(s.calculateStringValue(_h))

	if tc[0] == 0 {
		_h = "-" + _h
	}
	_a = 0
	for {
		_g[_a] = 0
		_a++
		if _a >= 100 {
			break
		}
	}
	_e, _a = 0, 0
	for {
		_b = 0
		for {
			_g[_b] *= 11
			_b++
			if _b > _e {
				break
			}
		}
		_d = 0
		_i = _h[_a : _a+1]
		for {
			if _l[_d:_d+1] == _i {
				break
			}
			_d++
			if _d > 9 {
				break
			}
		}
		_g[0] += _d
		_b = 0
		for {
			_c = _g[_b] / _k
			_g[_b] = _g[_b] - _c*_k
			_g[_b+1] += _c
			_b++
			if _b > _e {
				break
			}
		}
		if _c > 0 {
			_e++
		}
		_a++
		if _a >= len(_h) {
			break
		}
	}
	_h = ""
	for {
		if _e < 0 {
			break
		}
		_b = _e
		for {
			if _b <= 0 {
				break
			}
			_c = _g[_b] / _f
			_g[_b-1] += (_g[_b] - _c*_f) * _k
			_g[_b] = _c
			_b--
		}
		_c = _g[_b] / _f
		_a = _g[_b] - _c*_f
		_h += s.encryptionKey()[_a : _a+1]
		_g[_b] = _c
		if _g[_e] == 0 {
			_e--
		}
	}
	_a = len(_h)
	_j, _i = 0, ""
	for {
		_a--
		_i += _h[_a : _a+1]
		_j++
		if _j == 80 && _a > 0 {
			_i += "-"
			_j = 0
		}
		if _a <= 0 {
			break
		}
	}
	return _i
}
