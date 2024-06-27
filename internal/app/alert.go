package app

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AlertType string

const (
	alertEvent = "alert"

	AlertSuccess AlertType = "success"
	AlertError   AlertType = "error"
	AlertWarning AlertType = "warning"
	AlertInfo    AlertType = "info"
)

type alert struct {
	Type AlertType `json:"type"`
	Msg  string    `json:"message"`
}

func EmitAlert(ctx context.Context, alertType AlertType, msg string) {
	runtime.EventsEmit(ctx, alertEvent, &alert{
		Type: alertType,
		Msg:  msg,
	})
}
