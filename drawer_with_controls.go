package fromage

import (
	"time"

	"gioui.org/unit"
)

// DrawerWithControls combines a drawer with radio button controls to change its position
type DrawerWithControls struct {
	*Window
	drawer           *Drawer
	radioGroup       *RadioButtonGroup
	currentPos       DrawerPosition
	onPositionChange func(DrawerPosition)
}

// NewDrawerWithControls creates a new drawer with position controls
func (w *Window) NewDrawerWithControls() *DrawerWithControls {
	dwc := &DrawerWithControls{
		Window:     w,
		drawer:     w.NewDrawer(),
		currentPos: DrawerLeft,
	}

	// Create radio button group for position selection
	dwc.radioGroup = w.NewRadioButtonGroup().
		SetLayout(LayoutVertical).
		AddButton("Left", true).
		AddButton("Right", false).
		AddButton("Top", false).
		AddButton("Bottom", false)

	// Set up the drawer's initial position
	dwc.drawer.Position(dwc.currentPos)

	// Set up radio group change handler
	dwc.radioGroup.SetOnChange(func(index int, label string) {
		var newPos DrawerPosition
		switch label {
		case "Left":
			newPos = DrawerLeft
		case "Right":
			newPos = DrawerRight
		case "Top":
			newPos = DrawerTop
		case "Bottom":
			newPos = DrawerBottom
		}

		if newPos != dwc.currentPos {
			// Hide current drawer with animation
			dwc.drawer.Hide()

			// Update position after a short delay to allow hide animation
			go func() {
				// Wait for hide animation to complete
				time.Sleep(350 * time.Millisecond)

				// Update position and show with animation
				dwc.currentPos = newPos
				dwc.drawer.Position(newPos)
				dwc.drawer.Show()

				// Trigger callback if set
				if dwc.onPositionChange != nil {
					dwc.onPositionChange(newPos)
				}
			}()
		}
	})

	return dwc
}

// Content sets the content widget for the drawer
func (dwc *DrawerWithControls) Content(content W) *DrawerWithControls {
	dwc.drawer.Content(content)
	return dwc
}

// Width sets the width of the drawer (for left/right positions)
func (dwc *DrawerWithControls) Width(width unit.Dp) *DrawerWithControls {
	dwc.drawer.Width(width)
	return dwc
}

// Height sets the height of the drawer (for top/bottom positions)
func (dwc *DrawerWithControls) Height(height unit.Dp) *DrawerWithControls {
	dwc.drawer.Height(height)
	return dwc
}

// OnClose sets the callback function when the drawer is closed
func (dwc *DrawerWithControls) OnClose(fn func()) *DrawerWithControls {
	dwc.drawer.OnClose(fn)
	return dwc
}

// OnPositionChange sets the callback function when the drawer position changes
func (dwc *DrawerWithControls) OnPositionChange(fn func(DrawerPosition)) *DrawerWithControls {
	dwc.onPositionChange = fn
	return dwc
}

// Blocking sets whether this drawer blocks interaction with content behind it
func (dwc *DrawerWithControls) Blocking(blocking bool) *DrawerWithControls {
	dwc.drawer.Blocking(blocking)
	return dwc
}

// Show makes the drawer visible with animation
func (dwc *DrawerWithControls) Show() {
	dwc.drawer.Show()
}

// Hide makes the drawer invisible with animation
func (dwc *DrawerWithControls) Hide() {
	dwc.drawer.Hide()
}

// Toggle toggles the drawer visibility
func (dwc *DrawerWithControls) Toggle() {
	dwc.drawer.Toggle()
}

// IsVisible returns whether the drawer is currently visible
func (dwc *DrawerWithControls) IsVisible() bool {
	return dwc.drawer.IsVisible()
}

// GetCurrentPosition returns the current drawer position
func (dwc *DrawerWithControls) GetCurrentPosition() DrawerPosition {
	return dwc.currentPos
}

// SetPosition programmatically sets the drawer position
func (dwc *DrawerWithControls) SetPosition(position DrawerPosition) *DrawerWithControls {
	if position != dwc.currentPos {
		// Update radio button selection
		var index int
		switch position {
		case DrawerLeft:
			index = 0
		case DrawerRight:
			index = 1
		case DrawerTop:
			index = 2
		case DrawerBottom:
			index = 3
		}
		dwc.radioGroup.SetSelected(index)

		// Update position immediately for synchronous calls
		dwc.currentPos = position
		dwc.drawer.Position(position)

		// Trigger callback if set
		if dwc.onPositionChange != nil {
			dwc.onPositionChange(position)
		}
	}
	return dwc
}

// Layout renders the drawer with controls
func (dwc *DrawerWithControls) Layout(gtx C) D {
	// Layout the drawer (it handles its own visibility)
	return dwc.drawer.Layout(gtx)
}

// LayoutControls renders just the radio button controls
func (dwc *DrawerWithControls) LayoutControls(gtx C) D {
	return dwc.Theme.VFlex().
		SpaceEvenly().
		Rigid(func(gtx C) D {
			return dwc.Theme.H6("Drawer Position").
				Color(dwc.Theme.Colors.OnSurface()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return dwc.radioGroup.Layout(gtx)
		}).
		Layout(gtx)
}

// LayoutWithControls renders both the drawer and the controls in a layout
func (dwc *DrawerWithControls) LayoutWithControls(gtx C) D {
	return dwc.Theme.VFlex().
		SpaceEvenly().
		Rigid(func(gtx C) D {
			// Controls section
			return dwc.Theme.NewCard(func(gtx C) D {
				return dwc.Window.Inset(16, func(gtx C) D {
					return dwc.LayoutControls(gtx)
				}).Fn(gtx)
			}).Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			// Toggle button
			return dwc.Theme.TextButton("Toggle Drawer").
				OnClick(func() {
					dwc.Toggle()
				}).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			// Drawer (overlays everything when visible)
			return dwc.Layout(gtx)
		}).
		Layout(gtx)
}
