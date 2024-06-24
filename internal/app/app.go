package app

import (
	"context"
	"github.com/ktsivkov/ltd-he/pkg/history"
	"github.com/ktsivkov/ltd-he/pkg/player"
	"github.com/labstack/gommon/log"
	"log/slog"
)

type App struct {
	ctx            context.Context
	logger         *slog.Logger
	playerService  *player.Service
	historyService *history.Service
}

func New(logger *slog.Logger, playerService *player.Service, historyService *history.Service) *App {
	return &App{
		logger:         logger,
		playerService:  playerService,
		historyService: historyService,
	}
}

func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListPlayers() ([]*player.Player, error) {
	players, err := a.playerService.GetAll()
	if err != nil {
		log.Error("Could not get all players!", "error", err)
		return nil, err
	}

	return players, nil
}

func (a *App) LoadHistory(p *player.Player) ([]*history.GameHistory, error) {
	playerHistory, err := a.historyService.LoadHistory(p)
	if err != nil {
		log.Error("Could not get all players!", "error", err)
		return nil, err
	}

	return playerHistory, nil
}
