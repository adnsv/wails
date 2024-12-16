package runtime

import (
	"context"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

type Screen = frontend.Screen

// ScreenGetAll returns all screens
func ScreenGetAll(ctx context.Context) ([]Screen, error) {
	appFrontend := getFrontend(ctx)
	return appFrontend.ScreenGetAll()
}

// ScreenRect represents a rectangle in screen units
type ScreenRect = frontend.ScreenRect

// MonitorInfo provides information about a monitor and its coordinate spaces
type MonitorInfo = frontend.MonitorInfo

// MonitorGetAll returns information about all connected monitors
func MonitorGetAll(ctx context.Context) []MonitorInfo {
	if f := getFrontend(ctx); f != nil {
		monitors, _ := f.MonitorGetAll()
		return monitors
	}
	return nil
}
