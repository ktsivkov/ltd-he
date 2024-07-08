package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ktsivkov/ltd-he/pkg/backup"
	"github.com/ktsivkov/ltd-he/pkg/history"
	"github.com/ktsivkov/ltd-he/pkg/player"
)

type App struct {
	ctx            context.Context
	logger         *slog.Logger
	playerService  *player.Service
	historyService *history.Service
	backupService  *backup.Service
	mu             *sync.Mutex
}

func New(logger *slog.Logger, playerService *player.Service, historyService *history.Service, backupService *backup.Service) *App {
	return &App{
		logger:         logger,
		playerService:  playerService,
		historyService: historyService,
		backupService:  backupService,
		mu:             &sync.Mutex{},
	}
}

func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListPlayers() ([]*player.Player, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	players, err := a.playerService.LoadAll(a.ctx)
	if err != nil {
		a.logger.Error("Could not list all players!", "error", err)
		EmitAlert(a.ctx, AlertError, fmt.Sprintf("Could not list all players!\nError: %s", err))
		return nil, err
	}

	return players, nil
}

func (a *App) LoadHistory(p *player.Player) (history.History, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	playerHistory, err := a.historyService.Load(a.ctx, p)
	if err != nil {
		a.logger.Error("Could not load player history!", "error", err)
		EmitAlert(a.ctx, AlertError, fmt.Sprintf("Could not load player history!\nError: %s", err))
		return nil, err
	}

	return playerHistory, nil
}

func (a *App) Rollback(game *history.GameHistory) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	b, err := a.backupService.Backup(a.ctx, game.Account)
	if err != nil {
		a.logger.Error("Could not create a backup!", "error", err)
		EmitAlert(a.ctx, AlertError, fmt.Sprintf("Could not create a backup!\nError: %s", err))
		return err
	}
	a.logger.Info("Backup successfully created!", "backup_file", b.File)
	EmitAlert(a.ctx, AlertInfo, fmt.Sprintf("Successful backup!\nBackup location: %s", b.File))

	if err := a.historyService.Rollback(a.ctx, game); err != nil {
		a.logger.Error("Could not rollback history!", "error", err)
		EmitAlert(a.ctx, AlertError, fmt.Sprintf("Could not rollback history!\nError: %s", err))
		return err
	}
	a.logger.Info("History rollback was successful.", "rollback_target", game.TotalGames)
	EmitAlert(a.ctx, AlertSuccess, fmt.Sprintf("Successful rollback!"))

	return nil
}

func (a *App) Append(p *player.Player, req *history.AppendRequest) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	b, err := a.backupService.Backup(a.ctx, p)
	if err != nil {
		a.logger.Error("Could not create a backup!", "error", err)
		EmitAlert(a.ctx, AlertError, fmt.Sprintf("Could not create a backup!\nError: %s", err))
		return err
	}
	a.logger.Info("Backup successfully created!", "backup_file", b.File)
	EmitAlert(a.ctx, AlertInfo, fmt.Sprintf("Successful backup!\nBackup location: %s", b.File))

	if err := a.historyService.Append(a.ctx, p, req); err != nil {
		a.logger.Error("Could not append game!", "error", err)
		EmitAlert(a.ctx, AlertError, fmt.Sprintf("Could not append game!\nError: %s", err))
		return err
	}
	a.logger.Info("Game was appended successfully.", "insert_request", req, "target_player", p)
	EmitAlert(a.ctx, AlertSuccess, fmt.Sprintf("Successful game insertion!"))

	return nil
}

func (a *App) Insert(p *player.Player, req *history.InsertRequest) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	b, err := a.backupService.Backup(a.ctx, p)
	if err != nil {
		a.logger.Error("Could not create a backup!", "error", err)
		EmitAlert(a.ctx, AlertError, fmt.Sprintf("Could not create a backup!\nError: %s", err))
		return err
	}
	a.logger.Info("Backup successfully created!", "backup_file", b.File)
	EmitAlert(a.ctx, AlertInfo, fmt.Sprintf("Successful backup!\nBackup location: %s", b.File))

	if err := a.historyService.Insert(a.ctx, p, req); err != nil {
		a.logger.Error("Could not insert game!", "error", err)
		EmitAlert(a.ctx, AlertError, fmt.Sprintf("Could not insert game!\nError: %s", err))
		return err
	}
	a.logger.Info("Game was inserted successfully.", "insert_request", req, "target_player", p)
	EmitAlert(a.ctx, AlertSuccess, fmt.Sprintf("Successful game insertion!"))

	return nil
}

func (a *App) BackupFolder(p *player.Player) string {
	return a.backupService.BackupFolder(p)
}
