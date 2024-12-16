# Scaling Fixes

This document outlines the fixes for scaling issues identified in
`scaling-issues.md`.

## Introduction

Window positioning and scaling present unique challenges across different
platforms due to their distinct coordinate systems and DPI handling approaches:

- **Windows**: Uses physical pixels for most Win32 APIs. Window positions are in
  screen coordinates, and DPI scaling must be manually handled for proper
  positioning. The OS provides per-monitor DPI awareness and scaling
  information.

- **Linux**: Window management is handled through the window manager and GTK.

  Note: Fractional scaling support on Wayland is not a goal for this
  implementation as it requires GTK4, which is not yet supported. Only integer
  scaling will be properly supported.

- **macOS**: Uses points (logical pixels) consistently throughout its APIs. The
  system handles most scaling automatically, providing
  the conversion to physical pixels when needed.

## Terminology

Because all the platforms have their own understanding of DPI, pixels, etc., We
will try to come up with a terminology that is consistent across all platforms.

- **Physical Pixels**: The smallest addressable unit on the display.

- **Screen Units**: Units that are used by the display manager to position
  windows on the screen.

- **Logical Units**: Units that provide resolution-independent way of describing
  window size and position.

On Windows and Linux, screen units and physical pixels are the same. Window
positioning, monitor boundaries and mouse cursor movements are done in physical
pixel coordinates. Monitors, however, may have display scale factor. That factor
provides a mapping between (screen=physical) units and logical units. When
windows are dragged between monitors with different scale factors, the window
manager tries to resize and rescale the window when it is dragged between
monitors featuring different scale factors.

MacOS is different. In MacOS terminology, our screen units correspont to screen
points. On regular resolution displays, a screen unit is the same as physical
pixel. On Retina displays, a screen unit is 2x2 physical pixels. There is
scaling factors from screen units to logical units.

Inderstanding "Effective" DPI terminology:

- The effective DPI set by display managers is a logical construct to ensure
  consistent UI scaling and does not directly correspond to the monitor's
  physical DPI
- On Windows and linux:
  - Monitors with 100% scaling factor have 96 effective DPI resolution.
  - Monitors with 200% scaling factor have 192 effective DPI resolution.
- On macOS:
  - It is assumed that non-retina displays have 72 effective DPI resolution.
  - Retina displays have 144 effective DPI resolution.

## Implementation Strategy

### Phase 1: Minimal Invasion

1. **New Coordinate Management Layer**

   - Add new functions for handling window positioning without modifying
     existing code
   - Implement platform-specific coordinate conversion utilities
   - Maintain backward compatibility with current behavior

2. **Position Management**

   - Windows: Add work area and DPI-aware positioning
   - Linux: Handle window manager positioning through GTK
   - macOS: Ensure consistent point-based positioning

3. **Monitor Handling**
   - Track active monitor for each window
   - Handle DPI changes from:
     - System changing current monitor's DPI
     - Window moving to monitor with different DPI
     - Display reconfiguration (monitor added/removed)
   - Adjust window size and position on DPI changes

### Phase 2: Constraint Management

1. **Size Constraints**

   - Store constraints in logical pixels
   - Add proper scaling for physical constraints
   - Implement per-monitor constraint adjustment

2. **Platform-Specific Improvements**

   - Windows: Proper DPI scaling for min/max sizes
   - Linux: Proper GTK window manager integration
   - macOS: Backing scale factor consideration

3. **API Consistency**
   - Unified coordinate system across platforms
   - Clear documentation of coordinate spaces
   - Helper functions for coordinate conversion

## Implementation Steps

For window positioning, we will use screen units by default. This approach
aligns with display manager expectations and ensures consistent window placement
by avoiding rounding errors.

### API for querying monitor information

```go
// ScreenRect represents a rectangle in screen units
type ScreenRect struct {
    X int
    Y int
    Width int
    Height int
}

// MonitorInfo provides information about a monitor and its coordinate spaces
type MonitorInfo struct {
    // Monitor bounds in screen units
    Bounds ScreenRect

    // Work area bounds in screen units
    WorkArea ScreenRect

    // Scale factor for mapping from screen units to logical units
    // - Windows: DPI scale (96 DPI = 1.0, 192 DPI = 2.0)
    // - Linux: Scale factor (1, 2, etc.)
    // - macOS: Always 1.0
    //
    // To convert:
    // screen = logical * Scale
    // logical = screen / Scale
    Scale float64
}

// MonitorGetAll returns all monitors and their information
func MonitorGetAll() ([]MonitorInfo, error)
```

### API for querying window placement

```go
// WindowGetPlacement returns the current placement of the window
func WindowGetPlacement() (bounds ScreenRect, monitor MonitorInfo)
```

### API for setting placement for a window

```go
// WindowSetBounds sets the bounds of the window, in screen units
func WindowSetBounds(bounds ScreenRect)
```

### Implementation Phase 1

- Add new types to frontend.go:

  - `ScreenRect`
  - `MonitorInfo`

- Add new function to the FrontEnd interface in frontend.go:

```go
func MonitorGetAll() ([]MonitorInfo, error)
func WindowGetPlacement() (bounds ScreenRect, monitor MonitorInfo)
func WindowSetBounds(bounds ScreenRect)
```

- Implement support for monitor enumeration (new methods):

  - windows: v2/internal/frontend/desktop/windows/monitor.go
  - linux: v2/internal/frontend/desktop/linux/monitor.go
  - darwin: v2/internal/frontend/desktop/darwin/monitor.go

- Implement support for window placement:
  - windows: v2/internal/frontend/desktop/windows/window_extra.go
  - linux: v2/internal/frontend/desktop/linux/window_extra.go
  - darwin: v2/internal/frontend/desktop/darwin/window_extra.go
