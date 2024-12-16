//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit

#include <Cocoa/Cocoa.h>

void getWindowInfo(void* nsWindow, int* x, int* y, int* width, int* height) {
    NSWindow* window = (NSWindow*)nsWindow;
    NSRect frame = [window frame];

    *x = frame.origin.x;
    *y = frame.origin.y;
    *width = frame.size.width;
    *height = frame.size.height;
}

void getScreenInfoForWindow(void* nsWindow, int* x, int* y, int* width, int* height,
                          int* workX, int* workY, int* workWidth, int* workHeight) {
    NSWindow* window = (NSWindow*)nsWindow;
    NSScreen* screen = [window screen];
    if (screen == nil) {
        screen = [NSScreen mainScreen];
    }

    NSRect frame = [screen frame];
    NSRect visibleFrame = [screen visibleFrame];

    *x = frame.origin.x;
    *y = frame.origin.y;
    *width = frame.size.width;
    *height = frame.size.height;

    *workX = visibleFrame.origin.x;
    *workY = visibleFrame.origin.y;
    *workWidth = visibleFrame.size.width;
    *workHeight = visibleFrame.size.height;
}

void setWindowBounds(void* nsWindow, int x, int y, int width, int height) {
    NSWindow* window = (NSWindow*)nsWindow;
    NSRect frame = NSMakeRect(x, y, width, height);
    [window setFrame:frame display:YES animate:NO];
}
*/
import "C"
import (
	"github.com/wailsapp/wails/v2/internal/frontend"
)

func (f *Frontend) WindowGetPlacement() (bounds frontend.ScreenRect, monitor frontend.MonitorInfo) {
	var x, y, width, height C.int
	C.getWindowInfo(f.nsWindow, &x, &y, &width, &height)

	bounds = frontend.ScreenRect{
		X:      int(x),
		Y:      int(y),
		Width:  int(width),
		Height: int(height),
	}

	var workX, workY, workWidth, workHeight C.int

	C.getScreenInfoForWindow(f.nsWindow,
		&x, &y, &width, &height,
		&workX, &workY, &workWidth, &workHeight)

	monitor = frontend.MonitorInfo{
		Bounds: frontend.ScreenRect{
			X:      int(x),
			Y:      int(y),
			Width:  int(width),
			Height: int(height),
		},
		WorkArea: frontend.ScreenRect{
			X:      int(workX),
			Y:      int(workY),
			Width:  int(workWidth),
			Height: int(workHeight),
		},
		Scale: 1.0,
	}

	return
}

func (f *Frontend) WindowSetBounds(bounds frontend.ScreenRect) {
	C.setWindowBounds(f.nsWindow,
		C.int(bounds.X), C.int(bounds.Y),
		C.int(bounds.Width), C.int(bounds.Height))
}
