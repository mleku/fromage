package fromage

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

// DrawerPosition specifies which side the drawer slides from
type DrawerPosition int

const (
	DrawerLeft DrawerPosition = iota
	DrawerRight
	DrawerTop
	DrawerBottom
)

// Drawer represents a sliding drawer that can appear from any side
type Drawer struct {
	*Window
	content           W
	position          DrawerPosition
	width             unit.Dp
	height            unit.Dp
	scrimClickable    *widget.Clickable
	onClose           func()
	blocking          bool
	animationProgress float32   // 0.0 = invisible, 1.0 = fully visible
	animationStart    time.Time // When the animation started
	isAnimating       bool      // Whether animation is in progress
	animationStarted  bool      // Whether animation has ever been started
	isFadingOut       bool      // Whether we're fading out (true) or fading in (false)
	isVisible         bool      // Whether the drawer should be visible
}

// NewDrawer creates a new drawer
func (w *Window) NewDrawer() *Drawer {
	return &Drawer{
		Window:            w,
		position:          DrawerLeft,
		width:             unit.Dp(280), // Default width for left/right drawers
		height:            unit.Dp(200), // Default height for top/bottom drawers
		scrimClickable:    w.Theme.Pool.GetClickable(),
		blocking:          true,
		animationProgress: 0.0,
		animationStart:    time.Time{},
		isAnimating:       false,
		animationStarted:  false,
		isFadingOut:       false,
		isVisible:         false,
	}
}

// Content sets the content widget for the drawer
func (d *Drawer) Content(content W) *Drawer {
	d.content = content
	return d
}

// Position sets which side the drawer slides from
func (d *Drawer) Position(position DrawerPosition) *Drawer {
	d.position = position
	return d
}

// Width sets the width of the drawer (for left/right positions)
func (d *Drawer) Width(width unit.Dp) *Drawer {
	d.width = width
	return d
}

// Height sets the height of the drawer (for top/bottom positions)
func (d *Drawer) Height(height unit.Dp) *Drawer {
	d.height = height
	return d
}

// OnClose sets the callback function when the drawer is closed
func (d *Drawer) OnClose(fn func()) *Drawer {
	d.onClose = fn
	return d
}

// Blocking sets whether this drawer blocks interaction with content behind it
func (d *Drawer) Blocking(blocking bool) *Drawer {
	d.blocking = blocking
	return d
}

// Show makes the drawer visible with animation
func (d *Drawer) Show() {
	if d.isVisible {
		return
	}
	d.isVisible = true
	d.startAnimation(time.Now())
}

// Hide makes the drawer invisible with animation
func (d *Drawer) Hide() {
	if !d.isVisible {
		return
	}
	d.isVisible = false
	d.startFadeOut()
}

// Toggle toggles the drawer visibility
func (d *Drawer) Toggle() {
	if d.isVisible {
		d.Hide()
	} else {
		d.Show()
	}
}

// IsVisible returns whether the drawer is currently visible
func (d *Drawer) IsVisible() bool {
	return d.isVisible
}

// Layout renders the drawer
func (d *Drawer) Layout(gtx C) D {
	if !d.isVisible && !d.isAnimating {
		return D{}
	}

	// Start animation if not already started
	if !d.animationStarted && d.isVisible {
		d.startAnimation(gtx.Now)
	}

	// Update animation progress
	d.updateAnimation(gtx)

	// Don't render if completely hidden
	if d.animationProgress <= 0.0 && !d.isVisible {
		return D{}
	}

	// Create scrim color with animation progress
	scrimAlpha := uint8(255 * 0.5 * d.animationProgress) // 50% opacity when fully visible
	scrimColor := color.NRGBA{
		R: 0,
		G: 0,
		B: 0,
		A: scrimAlpha,
	}

	// Calculate drawer dimensions and position
	drawerSize, drawerOffset := d.calculateDrawerLayout(gtx)

	// Use a stack layout to properly position the scrim and drawer
	return layout.Stack{}.Layout(gtx,
		// First layer: Fill entire screen with scrim and handle clicks
		layout.Expanded(func(gtx C) D {
			// Fill with scrim color
			paint.Fill(gtx.Ops, scrimColor)

			// Layout clickable area over the entire scrim
			d.scrimClickable.Layout(gtx, func(gtx C) D {
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})

			// Handle scrim clicks (click outside to close)
			if d.scrimClickable.Clicked(gtx) {
				if d.onClose != nil {
					d.onClose()
				} else {
					d.Hide()
				}
			}

			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		// Second layer: Layout the drawer content
		layout.Stacked(func(gtx C) D {
			// Apply animation transform to the drawer
			defer d.applyDrawerTransform(gtx, drawerOffset, drawerSize).Pop()

			// Create a clip area for the drawer to prevent overflow
			defer clip.Rect(image.Rectangle{Max: drawerSize}).Push(gtx.Ops).Pop()

			// Constrain the drawer to its calculated size
			gtx.Constraints = layout.Exact(drawerSize)

			// Create a clickable area for the drawer content to prevent clicks from reaching the scrim
			contentClickable := &widget.Clickable{}
			return contentClickable.Layout(gtx, func(gtx C) D {
				// Fill drawer background
				paint.Fill(gtx.Ops, d.Theme.Colors.Surface())

				if d.content != nil {
					return d.content(gtx)
				}
				// Default content if none provided
				return d.Theme.Body1("Drawer Content").
					Color(d.Theme.Colors.OnSurface()).
					Layout(gtx)
			})
		}),
	)
}

// calculateDrawerLayout calculates the size and offset for the drawer based on its position
func (d *Drawer) calculateDrawerLayout(gtx C) (image.Point, image.Point) {
	screenSize := gtx.Constraints.Max
	drawerSize := image.Point{}
	drawerOffset := image.Point{}

	switch d.position {
	case DrawerLeft:
		drawerSize = image.Pt(gtx.Dp(d.width), screenSize.Y)
		drawerOffset = image.Pt(-drawerSize.X, 0)
	case DrawerRight:
		drawerSize = image.Pt(gtx.Dp(d.width), screenSize.Y)
		drawerOffset = image.Pt(screenSize.X, 0)
	case DrawerTop:
		drawerSize = image.Pt(screenSize.X, gtx.Dp(d.height))
		drawerOffset = image.Pt(0, -drawerSize.Y)
	case DrawerBottom:
		drawerSize = image.Pt(screenSize.X, gtx.Dp(d.height))
		drawerOffset = image.Pt(0, screenSize.Y)
	}

	return drawerSize, drawerOffset
}

// applyDrawerTransform applies the slide animation transform to the drawer
func (d *Drawer) applyDrawerTransform(gtx C, offset image.Point, size image.Point) op.TransformStack {
	// Calculate the current position based on animation progress
	// When progress is 0, drawer is at offset position (off-screen)
	// When progress is 1, drawer is at final position (on-screen)

	var finalOffset image.Point
	switch d.position {
	case DrawerLeft:
		finalOffset = image.Pt(0, 0)
	case DrawerRight:
		finalOffset = image.Pt(gtx.Constraints.Max.X-size.X, 0)
	case DrawerTop:
		finalOffset = image.Pt(0, 0)
	case DrawerBottom:
		finalOffset = image.Pt(0, gtx.Constraints.Max.Y-size.Y)
	}

	// Interpolate between offset and finalOffset based on animation progress
	currentOffset := image.Pt(
		offset.X+int(float32(finalOffset.X-offset.X)*d.animationProgress),
		offset.Y+int(float32(finalOffset.Y-offset.Y)*d.animationProgress),
	)

	// Apply the transform
	transform := op.Offset(currentOffset).Push(gtx.Ops)
	return transform
}

// startAnimation begins the slide-in animation
func (d *Drawer) startAnimation(now time.Time) {
	d.animationStart = now
	d.isAnimating = true
	d.animationStarted = true
	d.isFadingOut = false
}

// startFadeOut begins the slide-out animation
func (d *Drawer) startFadeOut() {
	d.animationStart = time.Now()
	d.isAnimating = true
	d.isFadingOut = true
}

// updateAnimation updates the animation progress
func (d *Drawer) updateAnimation(g C) {
	if !d.isAnimating {
		return
	}

	const animationDuration = 300 * time.Millisecond
	elapsed := g.Now.Sub(d.animationStart)

	if elapsed >= animationDuration {
		if d.isFadingOut {
			d.animationProgress = 0.0
			d.isAnimating = false
		} else {
			d.animationProgress = 1.0
			d.isAnimating = false
		}
		return
	}

	progress := float32(elapsed) / float32(animationDuration)
	if progress > 1.0 {
		progress = 1.0
	}

	// Use ease-out curve: 1 - (1-t)^3 for smoother animation
	easedProgress := 1.0 - (1.0-progress)*(1.0-progress)*(1.0-progress)

	if d.isFadingOut {
		d.animationProgress = 1.0 - easedProgress
	} else {
		d.animationProgress = easedProgress
	}

	if d.animationProgress > 1.0 {
		d.animationProgress = 1.0
	}
	if d.animationProgress < 0.0 {
		d.animationProgress = 0.0
	}

	g.Execute(op.InvalidateCmd{})
}

// Convenience methods for common drawer configurations

// LeftDrawer creates a drawer that slides from the left
func (w *Window) LeftDrawer() *Drawer {
	return w.NewDrawer().Position(DrawerLeft)
}

// RightDrawer creates a drawer that slides from the right
func (w *Window) RightDrawer() *Drawer {
	return w.NewDrawer().Position(DrawerRight)
}

// TopDrawer creates a drawer that slides from the top
func (w *Window) TopDrawer() *Drawer {
	return w.NewDrawer().Position(DrawerTop)
}

// BottomDrawer creates a drawer that slides from the bottom
func (w *Window) BottomDrawer() *Drawer {
	return w.NewDrawer().Position(DrawerBottom)
}
