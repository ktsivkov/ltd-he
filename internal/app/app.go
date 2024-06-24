package app

import (
	"context"
	"github.com/labstack/gommon/log"
	"log/slog"
	"ltd-he/pkg/player"
)

type App struct {
	ctx           context.Context
	logger        *slog.Logger
	playerService *player.Service
}

func New(logger *slog.Logger, playerService *player.Service) *App {
	return &App{
		logger:        logger,
		playerService: playerService,
	}
}

func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListPlayers() []*player.Player {
	players, err := a.playerService.GetAll()
	if err != nil {
		log.Error("Could not get all players!", "error", err)
	}

	return players
}
