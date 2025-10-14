# Drawer Component

The Drawer component provides a sliding panel that can appear from any of the four sides of the screen with smooth animations.

## Features

- **Four Positions**: Slide from left, right, top, or bottom
- **Smooth Animations**: Eased slide-in/slide-out transitions
- **Radio Button Controls**: Built-in controls to change drawer position
- **Scrim Overlay**: Semi-transparent background when drawer is open
- **Click Outside to Close**: Click on the scrim to close the drawer
- **Customizable Content**: Add any widgets as drawer content
- **Responsive Sizing**: Configurable width/height for different positions

## Basic Usage

### Simple Drawer

```go
// Create a basic drawer
drawer := win.NewDrawer().
    Position(fromage.DrawerLeft).
    Width(unit.Dp(300)).
    Content(func(gtx layout.Context) layout.Dimensions {
        return th.Body1("Drawer Content").Layout(gtx)
    })

// Show/hide the drawer
drawer.Show()
drawer.Hide()
drawer.Toggle()
```

### Drawer with Position Controls

```go
// Create a drawer with radio button controls
drawerWithControls := win.NewDrawerWithControls().
    Width(unit.Dp(300)).
    Height(unit.Dp(200)).
    Content(func(gtx layout.Context) layout.Dimensions {
        return th.VFlex().
            SpaceEvenly().
            Rigid(func(gtx layout.Context) layout.Dimensions {
                return th.H6("Drawer Content").Layout(gtx)
            }).
            Rigid(func(gtx layout.Context) layout.Dimensions {
                return th.TextButton("Close").
                    OnClick(func() {
                        drawerWithControls.Hide()
                    }).
                    Layout(gtx)
            }).
            Layout(gtx)
    }).
    OnPositionChange(func(pos fromage.DrawerPosition) {
        log.Printf("Drawer moved to position: %v", pos)
    })

// Layout the controls and drawer
drawerWithControls.LayoutControls(gtx) // Just the radio buttons
drawerWithControls.Layout(gtx)         // The drawer itself
```

## Drawer Positions

- `fromage.DrawerLeft` - Slides from the left edge
- `fromage.DrawerRight` - Slides from the right edge  
- `fromage.DrawerTop` - Slides from the top edge
- `fromage.DrawerBottom` - Slides from the bottom edge

## Convenience Methods

```go
// Create drawers for specific positions
leftDrawer := win.LeftDrawer()
rightDrawer := win.RightDrawer()
topDrawer := win.TopDrawer()
bottomDrawer := win.BottomDrawer()
```

## API Reference

### Drawer Methods

- `Position(pos DrawerPosition) *Drawer` - Set slide position
- `Width(width unit.Dp) *Drawer` - Set width (for left/right drawers)
- `Height(height unit.Dp) *Drawer` - Set height (for top/bottom drawers)
- `Content(content W) *Drawer` - Set drawer content
- `OnClose(fn func()) *Drawer` - Set close callback
- `Blocking(blocking bool) *Drawer` - Set blocking behavior
- `Show()` - Show drawer with animation
- `Hide()` - Hide drawer with animation
- `Toggle()` - Toggle drawer visibility
- `IsVisible() bool` - Check if drawer is visible

### DrawerWithControls Methods

- `Content(content W) *DrawerWithControls` - Set drawer content
- `Width(width unit.Dp) *DrawerWithControls` - Set drawer width
- `Height(height unit.Dp) *DrawerWithControls` - Set drawer height
- `OnClose(fn func()) *DrawerWithControls` - Set close callback
- `OnPositionChange(fn func(DrawerPosition)) *DrawerWithControls` - Set position change callback
- `SetPosition(pos DrawerPosition) *DrawerWithControls` - Programmatically set position
- `GetCurrentPosition() DrawerPosition` - Get current position
- `LayoutControls(gtx C) D` - Layout just the radio button controls
- `Layout(gtx C) D` - Layout the drawer
- `LayoutWithControls(gtx C) D` - Layout both controls and drawer

## Animation Details

- **Duration**: 300ms for slide animations
- **Easing**: Cubic ease-out curve for smooth motion
- **Scrim**: 50% opacity black overlay
- **Transform**: Hardware-accelerated slide transforms

## Example

See `cmd/drawer/main.go` for a complete working example that demonstrates:

- Creating a drawer with position controls
- Changing drawer position with radio buttons
- Smooth animations between positions
- Custom drawer content
- Show/hide/toggle functionality

## Running the Example

```bash
go run cmd/drawer/main.go
```

This will open a window with:
- A single "Open Drawer" button to show the drawer
- Radio buttons inside the drawer to select position (Left, Right, Top, Bottom)
- Smooth animations when changing positions (current state slides out, new state slides in)
- Click outside the drawer (on the scrim) or use the close button to hide it

Click the "Open Drawer" button to show the drawer, then use the radio buttons inside to change its position and watch the smooth slide animations.
