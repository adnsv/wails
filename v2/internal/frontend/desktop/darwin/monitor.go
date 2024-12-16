//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit

#include <Cocoa/Cocoa.h>

void getScreenInfo(NSScreen* screen, int* x, int* y, int* width, int* height,
                  int* workX, int* workY, int* workWidth, int* workHeight) {
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
*/
import "C"
import (
	"github.com/wailsapp/wails/v2/internal/frontend"
)

func (f *Frontend) MonitorGetAll() ([]frontend.MonitorInfo, error) {
	var monitors []frontend.MonitorInfo

	// Get NSScreen array
	screens := C.NSScreen.screens()
	count := C.NSArray_count(screens)

	// Iterate through screens
	for i := C.NSUInteger(0); i < count; i++ {
		screen := C.NSArray_objectAtIndex(screens, i)
		if screen == nil {
			continue
		}

		var x, y, width, height C.int
		var workX, workY, workWidth, workHeight C.int

		// Get screen information
		C.getScreenInfo(screen,
			&x, &y, &width, &height,
			&workX, &workY, &workWidth, &workHeight)

		// Convert to our format
		monitorInfo := frontend.MonitorInfo{
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
			Scale: 1.0, // Screen units match logical units on macOS
		}

		monitors = append(monitors, monitorInfo)
	}

	return monitors, nil
}
