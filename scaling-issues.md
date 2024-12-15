# Scaling Issues

## Windows

The Windows implementation currently has inconsistent handling of window coordinates and DPI scaling. Here are the identified issues:

### Coordinate System Inconsistencies

1. **Window Size Operations**
   - Basic size methods use logical coordinates without DPI scaling:
     ```go
     // Current implementation ignores DPI scaling
     func (w *Window) SetMinSize(minWidth int, minHeight int) {
         w.minWidth = minWidth   // Stored as logical coordinates
         w.minHeight = minHeight // Should be scaled by DPI
         w.Form.SetMinSize(minWidth, minHeight)
     }
     ```
   - Window creation uses raw values without explicit scaling:
     ```go
     // Raw values from options used directly
     err := wails.Run(&options.App{
         Width:     width,      // Should consider DPI
         Height:    height,     // Should consider DPI
         MinWidth:  minWidth,   // Should consider DPI
         MinHeight: minHeight,  // Should consider DPI
     })
     ```
   - DPI changed event handling is incomplete:
     ```go
     case w32.WM_DPICHANGED:
         // Only handles window resizing, missing constraint updates
         newWindowSize := (*w32.RECT)(unsafe.Pointer(lparam))
         w32.SetWindowPos(w.Handle(),
             0,
             int(newWindowSize.Left),
             int(newWindowSize.Top),
             int(newWindowSize.Right-newWindowSize.Left),
             int(newWindowSize.Bottom-newWindowSize.Top),
             w32.SWP_NOZORDER|w32.SWP_NOACTIVATE)
     ```
   - Window constraints (min/max sizes) not consistently scaled across all window states
   - DPI scaling only applied to constraints in maximized state

2. **Window Position Operations**
   - Position changes don't handle DPI scaling:
     ```go
     // Current implementation uses raw coordinates
     func (w *Window) SetPosition(x int, y int) {
         w32.SetWindowPos(w.Handle(),
             0, x, y, 0, 0,
             w32.SWP_NOSIZE|w32.SWP_NOZORDER)
     }
     ```
   - Monitor transitions not handled:
     ```go
     case w32.WM_MOVE, w32.WM_MOVING:
         // Only notifies webview, no DPI handling
         w.chromium.NotifyParentWindowPositionChanged()
     ```
   - Position notifications to Chromium likely expect pixel coordinates
   - No tracking of monitor changes during window movement
   - Missing DPI adjustment when moving between monitors
   - Inconsistent work area handling between GetPosition and SetPosition:
     ```go
     // Center method correctly uses work area
     info := getMonitorInfo(fm.hwnd)
     workRect := info.RcWork  // Uses work area for calculations
     
     // But SetPosition doesn't consider work area
     func (w *Window) SetPosition(x int, y int) {
         w32.SetWindowPos(w.Handle(),
             0, x, y, 0, 0,
             w32.SWP_NOSIZE|w32.SWP_NOZORDER)  // Raw coordinates used
     }
     ```

3. **Monitor/Display Handling**
   - Maximized window handling shows correct approach:
     ```go
     // Good: Uses monitor DPI for constraints
     monitor := w32.MonitorFromRect(rgrc, w32.MONITOR_DEFAULTTONULL)
     var dpiX, dpiY uint
     w32.GetDPIForMonitor(monitor, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)
     maxWidth := int32(winc.ScaleWithDPI(maxWidth, dpiX))
     ```
   - But missing in regular window operations:
     ```go
     // Missing: Regular window operations don't use monitor DPI
     func (w *Window) Center() {
         // Should consider monitor DPI
         w32.CenterWindow(w.Handle())
     }
     ```
   - Other window positioning code lacks consistent DPI awareness
   - Monitor work area only handled in maximized state
   - No proper handling of window movement between monitors with different DPIs

4. **DPI Change Event Handling**
   - Current implementation misses critical updates:
     ```go
     case w32.WM_DPICHANGED:
         // Missing:
         // - Update internal DPI state
         // - Scale window constraints
         // - Update webview DPI
         // - Handle content scaling
         newWindowSize := (*w32.RECT)(unsafe.Pointer(lparam))
         w32.SetWindowPos(w.Handle(), ...)
     ```
   - `WM_DPICHANGED` only handles window resizing
   - Internal min/max size constraints not updated on DPI changes
   - No notification to webview about DPI changes
   - Incomplete per-monitor DPI awareness implementation
   - Missing DPI scale factor updates for webview content

### Impact
- Inconsistent window sizes on high-DPI displays
- Potential positioning issues in multi-monitor setups with different DPI settings
- Unexpected behavior when moving windows between monitors with different DPI settings
- Window constraints may not scale properly with DPI changes
- Webview content may not scale correctly with DPI changes

### Windows OS Implementation Details

1. **DPI Awareness**
   - Application correctly declares DPI awareness in manifest:
     ```xml
     <asmv3:application>
         <asmv3:windowsSettings>
             <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true/pm</dpiAware>
             <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">permonitorv2,permonitor</dpiAwareness>
         </asmv3:windowsSettings>
     </asmv3:application>
     ```

2. **Window Position Handling**
   - Current implementation uses raw coordinates without DPI awareness
   - Work area handling is inconsistent between different window operations

## Linux

The Linux implementation using GTK3 and WebKit2GTK has its own set of scaling-related challenges:

### Coordinate System Inconsistencies

1. **Monitor Scale Factor Handling**
   - Scale factor queried but not consistently applied:
     ```c
     // Current implementation only gets scale factor but doesn't use it
     static int getCurrentMonitorScaleFactor(GtkWindow *window) {
         GdkMonitor *monitor = getCurrentMonitor(window);
         return gdk_monitor_get_scale_factor(monitor);
     }
     ```
   - Window geometry operations don't account for scale factor
   - No dynamic updates when scale factor changes

2. **Window Size Operations**
   - Min/max size constraints don't account for monitor scale factor:
     ```c
     // Current implementation ignores DPI scaling
     void SetMinMaxSize(GtkWindow *window, int min_width, int min_height, int max_width, int max_height) {
         GdkGeometry size;
         size.min_width = min_width;   // Should be scaled by monitor DPI
         size.min_height = min_height; // Should be scaled by monitor DPI
         gtk_window_set_geometry_hints(window, NULL, &size, GDK_HINT_MIN_SIZE);
     }
     ```
   - Window creation uses raw values without scaling
   - Geometry hints not adjusted for DPI

3. **Window Position Operations**
   - Position setting uses monitor-relative coordinates without scaling:
     ```c
     // Current implementation misses DPI scaling
     void SetPosition(void *window, int x, int y) {
         GdkRectangle monitorDimensions = getCurrentMonitorGeometry(window);
         args->x = monitorDimensions.x + x;  // Should consider monitor scale
         args->y = monitorDimensions.y + y;  // Should consider monitor scale
         gtk_window_move(window, args->x, args->y);
     }
     ```
   - Missing DPI adjustment when moving between monitors

4. **Monitor/Display Handling**
   - Basic monitor detection implemented but lacks scaling awareness:
     ```c
     // Current implementation doesn't handle DPI changes
     static GdkMonitor *getCurrentMonitor(GtkWindow *window) {
         GdkDisplay *display = gtk_widget_get_display(GTK_WIDGET(window));
         GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(window));
         return gdk_display_get_monitor_at_window(display, gdk_window);
     }
     ```
   - No handling of runtime monitor configuration changes

### Impact
- Window sizes may not match specified values on scaled displays
- Window position may be incorrect when moving between monitors
- Inconsistent appearance across different scale factors
- Potential issues with window decorations and constraints

## macOS

The macOS implementation using Cocoa and WebKit has several scaling-related considerations and issues:

### Coordinate System Inconsistencies

1. **Screen Resolution Handling**
   - Native resolution detection uses complex fallback chain:
     ```objc
     // Current implementation with multiple fallbacks
     CGDirectDisplayID sid = [screen.deviceDescription[@"NSScreenNumber"] unsignedIntegerValue];
     CFArrayRef modes = CGDisplayCopyAllDisplayModes(sid, NULL);
     // Fallback if native mode not found
     NSRect pSize = [screen convertRectToBacking:screen.frame];
     ```
   - No dynamic handling of resolution changes
   - Missing notification system for resolution updates

2. **Window Size Operations**
   - Window size methods use points without scale factor consideration:
     ```objc
     // Current implementation ignores backing scale factor
     void SetSize(void* inctx, int width, int height) {
         WailsContext *ctx = (__bridge WailsContext*) inctx;
         [ctx SetSize:width :height];  // Raw values used directly
     }
     ```
   - Min/max constraints stored without scale factor:
     ```objc
     // Constraints not scaled for display
     void SetMinSize(void* inctx, int width, int height) {
         WailsContext *ctx = (__bridge WailsContext*) inctx;
         [ctx SetMinSize:width :height];  // Should consider backingScaleFactor
     }
     ```

3. **Window Position Operations**
   - Position calculations require manual coordinate system conversion:
     ```objc
     // Current implementation with Y-coordinate inversion
     const char* GetPosition(void *inctx) {
         NSRect windowFrame = [ctx.mainWindow frame];
         NSRect screenFrame = [screen visibleFrame];
         int x = windowFrame.origin.x - screenFrame.origin.x;
         int y = screenFrame.origin.y - screenFrame.origin.y;
         y = screenFrame.size.height - y - windowFrame.size.height;  // Manual flip
         return [NSString stringWithFormat:@"%d,%d", x, y].UTF8String;
     }
     ```
   - Missing scale factor adjustments for multi-screen setups

4. **Display/Screen Management**
   - Screen enumeration lacks proper scale factor handling:
     ```objc
     // Missing scale factor in screen info
     Screen GetNthScreen(int nth, void *inctx) {
         Screen returnScreen;
         returnScreen.width = (int) nthScreen.frame.size.width;   // Raw size
         returnScreen.height = (int) nthScreen.frame.size.height; // Raw size
         // No scale factor information stored
         return returnScreen;
     }
     ```

### Impact
- Potential scaling issues on retina displays
- Window size inconsistencies across different scale factors
- Position calculation issues in multi-monitor setups
- Missing dynamic display configuration updates