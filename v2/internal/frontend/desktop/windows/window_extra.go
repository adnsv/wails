package windows

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

func (f *Frontend) WindowGetPlacement() (bounds frontend.ScreenRect, monitor frontend.MonitorInfo) {
	// Get window rect - actual position in pixels, which matches screen units on Windows systems
	rect := w32.GetWindowRect(f.mainWindow.Handle())

	// Convert to our types
	bounds = frontend.ScreenRect{
		X:      int(rect.Left),
		Y:      int(rect.Top),
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	}

	// Get monitor
	handle := w32.MonitorFromWindow(f.mainWindow.Handle(), w32.MONITOR_DEFAULTTONEAREST)
	monInfo, err := GetMonitorInfo(handle)
	if err != nil {
		return
	}

	// Get DPI for monitor
	var dpiX, dpiY w32.UINT
	w32.GetDPIForMonitor(handle, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)

	monitor = frontend.MonitorInfo{
		Bounds: frontend.ScreenRect{
			X:      int(monInfo.RcMonitor.Left),
			Y:      int(monInfo.RcMonitor.Top),
			Width:  int(monInfo.RcMonitor.Right - monInfo.RcMonitor.Left),
			Height: int(monInfo.RcMonitor.Bottom - monInfo.RcMonitor.Top),
		},
		WorkArea: frontend.ScreenRect{
			X:      int(monInfo.RcWork.Left),
			Y:      int(monInfo.RcWork.Top),
			Width:  int(monInfo.RcWork.Right - monInfo.RcWork.Left),
			Height: int(monInfo.RcWork.Bottom - monInfo.RcWork.Top),
		},
		Scale: float64(dpiX) / 96.0,
	}

	return
}

func (f *Frontend) WindowSetBounds(bounds frontend.ScreenRect) {
	f.mainWindow.Invoke(func() {
		w32.SetWindowPos(f.mainWindow.Handle(), 0,
			bounds.X, bounds.Y, bounds.Width, bounds.Height,
			w32.SWP_NOZORDER|w32.SWP_NOACTIVATE)
	})
}
