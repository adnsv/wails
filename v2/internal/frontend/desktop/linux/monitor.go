//go:build linux
// +build linux

package linux

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// #include <gdk/gdk.h>
// #include "window.h"
import "C"

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
)

func (f *Frontend) MonitorGetAll() ([]frontend.MonitorInfo, error) {
	var monitors []frontend.MonitorInfo

	// Get GDK display
	display := C.gdk_display_get_default()
	if display == nil {
		return nil, nil
	}

	// Get GDK screen
	screen := C.gdk_display_get_default_screen(display)
	if screen == nil {
		return nil, nil
	}

	// Get number of monitors
	numMonitors := int(C.gdk_display_get_n_monitors(display))

	// Iterate through monitors
	for i := 0; i < numMonitors; i++ {
		monitor := C.gdk_display_get_monitor(display, C.int(i))
		if monitor == nil {
			continue
		}

		// Get monitor geometry
		var geometry C.GdkRectangle
		C.gdk_monitor_get_geometry(monitor, &geometry)

		// Get work area
		var workArea C.GdkRectangle
		C.gdk_monitor_get_workarea(monitor, &workArea)

		// Get scale factor
		scale := float64(C.gdk_monitor_get_scale_factor(monitor))

		// Convert to our format
		monitorInfo := frontend.MonitorInfo{
			Bounds: frontend.ScreenRect{
				X:      int(geometry.x),
				Y:      int(geometry.y),
				Width:  int(geometry.width),
				Height: int(geometry.height),
			},
			WorkArea: frontend.ScreenRect{
				X:      int(workArea.x),
				Y:      int(workArea.y),
				Width:  int(workArea.width),
				Height: int(workArea.height),
			},
			Scale: scale,
		}

		monitors = append(monitors, monitorInfo)
	}

	return monitors, nil
}
