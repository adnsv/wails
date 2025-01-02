//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit

#include <Cocoa/Cocoa.h>

int getMonitorCount() {
    return [[NSScreen screens] count];
}

void getScreenInfo(int index, int* x, int* y, int* width, int* height,
                  int* workX, int* workY, int* workWidth, int* workHeight) {

	NSArray<NSScreen *> *screens = [NSScreen screens];
	NSScreen* screen = [screens objectAtIndex:index];

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

	count := C.getMonitorCount()

	// Get main screen height for Y-coordinate conversion
	var mainX, mainY, mainWidth, mainHeight C.int
	var mainWorkX, mainWorkY, mainWorkWidth, mainWorkHeight C.int
	C.getScreenInfo(0, &mainX, &mainY, &mainWidth, &mainHeight,
		&mainWorkX, &mainWorkY, &mainWorkWidth, &mainWorkHeight)

	// Iterate through screens
	for index := C.int(0); index < count; index++ {
		var x, y, width, height C.int
		var workX, workY, workWidth, workHeight C.int

		// Get screen information
		C.getScreenInfo(index,
			&x, &y, &width, &height,
			&workX, &workY, &workWidth, &workHeight)

		// Convert Y coordinates by subtracting from main screen height
		convertedY := int(mainHeight) - (int(y) + int(height))
		convertedWorkY := int(mainHeight) - (int(workY) + int(workHeight))

		// Convert to our format
		monitorInfo := frontend.MonitorInfo{
			Bounds: frontend.ScreenRect{
				X:      int(x),
				Y:      convertedY,
				Width:  int(width),
				Height: int(height),
			},
			WorkArea: frontend.ScreenRect{
				X:      int(workX),
				Y:      convertedWorkY,
				Width:  int(workWidth),
				Height: int(workHeight),
			},
			Scale: 1.0, // Screen units match logical units on macOS
		}

		monitors = append(monitors, monitorInfo)
	}

	return monitors, nil
}
