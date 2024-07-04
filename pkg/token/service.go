package token

import (
	_ "embed"
	"errors"
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"
)

//go:embed generator.lua
var generatorScript string

//go:embed validator.lua
var validatorScript string

const (
	eloMaxLvlRound = 31
	payloadSize    = 28
	yd             = 1
)

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) Token(player string, totalGames int, wins int, elo int, gamesLeftEarly int, winsStreak int, highestWinStreak int, mvp int, timestamp time.Time, wasWin bool) (string, error) {
	token := &tableContents{
		totalGames:       totalGames,
		wins:             wins,
		elo:              elo,
		gamesLeftEarly:   gamesLeftEarly,
		winsStreak:       winsStreak,
		highestWinStreak: highestWinStreak,
		mvp:              mvp,
		timestamp:        timestamp,
		// Defaults
		wasWin:         wasWin,
		eloMaxLvlRound: eloMaxLvlRound,
		payloadSize:    payloadSize,
		yd:             yd,
	}

	luaVm := lua.NewState()
	defer luaVm.Close()

	luaTableContents := luaVm.NewTable()
	for _, v := range token.slice() {
		luaTableContents.Append(lua.LNumber(v))
	}

	luaVm.SetGlobal("TableContents", luaTableContents)
	luaVm.SetGlobal("Player", lua.LString(player))

	if err := luaVm.DoString(generatorScript); err != nil {
		return "", fmt.Errorf("token generator script failed: %w", err)
	}

	luaGlobalString := luaVm.GetGlobal("result")
	tokenStr, ok := luaGlobalString.(lua.LString)
	if !ok {
		return "", errors.New("could not extract token from the lua generator")
	}

	return string(tokenStr), nil
}

func (s *Service) ValidateToken(player string, token string) (bool, error) {
	luaVm := lua.NewState()
	defer luaVm.Close()

	luaVm.SetGlobal("Player", lua.LString(player))
	luaVm.SetGlobal("Token", lua.LString(token))

	if err := luaVm.DoString(validatorScript); err != nil {
		return false, fmt.Errorf("token validator script failed: %w", err)
	}

	luaGlobalString := luaVm.GetGlobal("result")
	validBool, ok := luaGlobalString.(lua.LBool)
	if !ok {
		return false, errors.New("could not extract result of token validator")
	}

	return bool(validBool), nil
}
