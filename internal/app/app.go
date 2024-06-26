package app

import (
	"context"
	"github.com/ktsivkov/ltd-he/pkg/backup"
	"github.com/ktsivkov/ltd-he/pkg/history"
	"github.com/ktsivkov/ltd-he/pkg/player"
	"log/slog"
)

type App struct {
	ctx            context.Context
	logger         *slog.Logger
	playerService  *player.Service
	historyService *history.Service
	backupService  *backup.Service
}

func New(logger *slog.Logger, playerService *player.Service, historyService *history.Service, backupService *backup.Service) *App {
	return &App{
		logger:         logger,
		playerService:  playerService,
		historyService: historyService,
		backupService:  backupService,
	}
}

func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListPlayers() ([]*player.Player, error) {
	players, err := a.playerService.LoadAll(a.ctx)
	if err != nil {
		a.logger.Error("Could not get all players!", "error", err)
		return nil, err
	}

	return players, nil
}

func (a *App) LoadHistory(p *player.Player) (history.History, error) {
	playerHistory, err := a.historyService.Load(a.ctx, p)
	if err != nil {
		a.logger.Error("Could not get all players!", "error", err)
		return nil, err
	}

	return playerHistory, nil
}

func (a *App) Rollback(game *history.GameHistory) error {
	b, err := a.backupService.Backup(a.ctx, game.Account)
	if err != nil {
		a.logger.Error("Could not create a backup!", "error", err)
		return err
	}
	a.logger.Info("Backup successfully created!", "backup_file", b.File)

	return nil
}
