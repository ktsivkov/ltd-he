package main

import (
	"embed"
	"github.com/ktsivkov/ltd-he/pkg/game_stats"
	"github.com/ktsivkov/ltd-he/pkg/history"
	"github.com/ktsivkov/ltd-he/pkg/player"
	"github.com/ktsivkov/ltd-he/pkg/report"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/ktsivkov/ltd-he/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

const DocumentsDir = "Documents"
const Wc3Dir = "Warcraft III"
const CustomMapDataDir = "CustomMapData"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Could not determine user home directory!", "error", err)
		os.Exit(1)
	}

	customMapDataPath := filepath.Join(homeDir, DocumentsDir, Wc3Dir, CustomMapDataDir)

	playerService := player.NewService(customMapDataPath)
	reportService := report.NewService()
	gameStatsService := game_stats.NewService()
	historyService := history.NewService(reportService, gameStatsService)

	appInstance := app.New(logger, playerService, historyService)
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
