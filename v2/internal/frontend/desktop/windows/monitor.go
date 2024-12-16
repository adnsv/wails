//go:build windows
// +build windows

package windows

import (
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

type monitorEnumData struct {
	monitors []frontend.MonitorInfo
}

func monitorEnumProc(handle w32.HMONITOR, _, _ *w32.RECT, enumData *monitorEnumData) uintptr {
	monInfo, err := GetMonitorInfo(handle)
	if err != nil {
		return w32.TRUE // continue enumeration
	}

	// Get DPI for monitor
	var dpiX, dpiY w32.UINT
	w32.GetDPIForMonitor(handle, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)

	// Convert monitor info to our format
	monitor := frontend.MonitorInfo{
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
		Scale: float64(dpiX) / 96.0, // Convert DPI to scale factor (96 DPI = 1.0)
	}

	enumData.monitors = append(enumData.monitors, monitor)
	return w32.TRUE
}

func (f *Frontend) MonitorGetAll() ([]frontend.MonitorInfo, error) {
	var enumData monitorEnumData

	// Enumerate monitors
	dc := w32.GetDC(0)
	defer w32.ReleaseDC(0, dc)

	succeeded := w32.EnumDisplayMonitors(dc, nil, syscall.NewCallback(monitorEnumProc), unsafe.Pointer(&enumData))
	if !succeeded {
		return nil, nil
	}

	return enumData.monitors, nil
}
