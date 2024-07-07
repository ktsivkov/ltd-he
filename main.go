package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ktsivkov/ltd-he/pkg/backup"
	"github.com/ktsivkov/ltd-he/pkg/game_stats"
	"github.com/ktsivkov/ltd-he/pkg/history"
	"github.com/ktsivkov/ltd-he/pkg/player"
	"github.com/ktsivkov/ltd-he/pkg/report"
	"github.com/ktsivkov/ltd-he/pkg/storage"
	"github.com/ktsivkov/ltd-he/pkg/token"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/ktsivkov/ltd-he/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

const DocumentsDir = "Documents"
const Wc3Dir = "Warcraft III"
const AppDir = "LegionTD History Editor"
const LogFile = "runtime.log"

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not determine user home directory!", "error", err)
	}

	documentsPath := filepath.Join(homeDir, DocumentsDir)
	wc3Path := filepath.Join(documentsPath, Wc3Dir)
	appPath := filepath.Join(documentsPath, AppDir)

	if _, err := os.Stat(appPath); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal("Could not determine user home directory!", "error", err)
		}
		if err := os.MkdirAll(appPath, os.ModePerm); err != nil {
			log.Fatal("Could not determine user home directory!", "error", err)
		}
	}

	logFile, err := os.OpenFile(filepath.Join(appPath, LogFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer logFile.Close()
	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	storageDriver := storage.New()
	playerService := player.NewService(wc3Path, storageDriver)
	reportService := report.NewService(storageDriver)
	gameStatsService := game_stats.NewService(storageDriver)
	tokenService := token.NewService()
	historyService := history.NewService(reportService, gameStatsService, tokenService, storageDriver)
	backupService := backup.NewService(appPath)

	appInstance := app.New(logger, playerService, historyService, backupService)
	if err := wails.Run(&options.App{
		Title:  "LegionTD History Editor - by ktsivkov",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		DisableResize:    true,
		OnStartup:        appInstance.OnStartup,
		Bind: []interface{}{
			appInstance,
		},
	}); err != nil {
		logger.Error("Wails application failed with an error!", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}
