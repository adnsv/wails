//go:build linux
// +build linux

package linux

// #cgo pkg-config: gtk+-3.0 webkit2gtk-4.0
// #include <gtk/gtk.h>
// #include <gdk/gdk.h>
// #include "window.h"
import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

func (f *Frontend) WindowGetPlacement() (bounds frontend.ScreenRect, monitor frontend.MonitorInfo) {
	// Get window geometry
	var x, y, width, height C.int
	C.gtk_window_get_position(f.gtkWindow, &x, &y)
	C.gtk_window_get_size(f.gtkWindow, &width, &height)

	bounds = frontend.ScreenRect{
		X:      int(x),
		Y:      int(y),
		Width:  int(width),
		Height: int(height),
	}

	// Get monitor for window
	gdkWindow := C.gtk_widget_get_window(C.GtkWidget(unsafe.Pointer(f.gtkWindow)))
	if gdkWindow == nil {
		return
	}

	display := C.gdk_display_get_default()
	if display == nil {
		return
	}

	gmonitor := C.gdk_display_get_monitor_at_window(display, gdkWindow)
	if gmonitor == nil {
		return
	}

	// Get monitor geometry
	var geometry, workArea C.GdkRectangle
	C.gdk_monitor_get_geometry(gmonitor, &geometry)
	C.gdk_monitor_get_workarea(gmonitor, &workArea)

	// Get scale factor
	scale := float64(C.gdk_monitor_get_scale_factor(gmonitor))

	monitor = frontend.MonitorInfo{
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

	return
}

func (f *Frontend) WindowSetBounds(bounds frontend.ScreenRect) {
	C.gtk_window_move(f.gtkWindow, C.int(bounds.X), C.int(bounds.Y))
	C.gtk_window_resize(f.gtkWindow, C.int(bounds.Width), C.int(bounds.Height))
}
