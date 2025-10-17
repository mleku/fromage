package fromage

import (
	"image"
	"time"

	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/unit"
)

// Scrollbar is a scrollbar widget for indicating scroll position and allowing scrolling
type Scrollbar struct {
	// Viewport represents the visible portion of the content (0-1)
	viewport float32
	// Position represents the scroll position within the content (0-1)
	position float32
	// Orientation determines if the scrollbar is horizontal or vertical
	orientation Orientation
	// Width is the width of the scrollbar (defaults to 1 text height)
	width unit.Dp

	eventHandler *EventHandler
	changed      bool
	changeHook   func(float32)
	dragging     bool

	// Animation for smooth track clicks
	animating     bool
	animStartTime time.Time
	animStartPos  float32
	animTargetPos float32

	// Long press tracking for track clicks
	trackPressed    bool
	trackPressStart time.Time
	trackPressSide  int // -1 for left/up, 1 for right/down

	// Track dimensions for event handling
	trackLength float32
	thumbSize   int
	// Store gtx context for event handling
	gtx layout.Context
}

// Orientation represents the scrollbar orientation
type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

// NewScrollbar creates a new scrollbar
func (t *Theme) NewScrollbar(orientation Orientation) *Scrollbar {
	sb := &Scrollbar{
		changeHook:  func(float32) {},
		orientation: orientation,
		width:       t.TextSize, // Default to 1 text height wide
		viewport:    0.5,        // Default to showing 50% of content
		position:    0.0,        // Default to start position
		dragging:    false,
	}

	// Create event handler with callbacks
	sb.eventHandler = NewEventHandler(func(event string) {
		// Log scrollbar events if needed
	}).SetOnPress(func(e pointer.Event) {
		sb.handlePress(e)
	}).SetOnDrag(func(e pointer.Event) {
		sb.handleDrag(e)
	}).SetOnRelease(func(e pointer.Event) {
		sb.handleRelease(e)
	}).SetOnScroll(func(distance float32) {
		sb.handleScroll(distance)
	})

	return sb
}

// SetViewport sets the viewport size (0-1, where 1 means all content is visible)
func (s *Scrollbar) SetViewport(viewport float32) *Scrollbar {
	if viewport < 0 {
		viewport = 0
	} else if viewport > 1 {
		viewport = 1
	}
	s.viewport = viewport
	return s
}

// SetPosition sets the scroll position (0-1, where 0 is start, 1 is end)
func (s *Scrollbar) SetPosition(position float32) *Scrollbar {
	if position < 0 {
		position = 0
	} else if position > 1 {
		position = 1
	}
	s.position = position
	return s
}

// SetWidth sets the width of the scrollbar
func (s *Scrollbar) SetWidth(width unit.Dp) *Scrollbar {
	s.width = width
	return s
}

// SetHook sets the change callback
func (s *Scrollbar) SetHook(fn func(float32)) *Scrollbar {
	s.changeHook = fn
	return s
}

// Viewport returns the current viewport size
func (s *Scrollbar) Viewport() float32 {
	return s.viewport
}

// Position returns the current scroll position
func (s *Scrollbar) Position() float32 {
	return s.position
}

// handlePress handles press events
func (s *Scrollbar) handlePress(e pointer.Event) {
	// Check if click is on thumb or track
	var clickPos float32
	if s.orientation == Horizontal {
		clickPos = e.Position.X / s.trackLength
	} else {
		clickPos = e.Position.Y / s.trackLength
	}

	// Calculate thumb position and size in normalized coordinates
	thumbRatio := float32(s.thumbSize) / s.trackLength
	thumbStart := s.position * (1 - thumbRatio)
	thumbEnd := thumbStart + thumbRatio

	if clickPos >= thumbStart && clickPos <= thumbEnd {
		// Click is on thumb - start dragging
		s.dragging = true
	} else {
		// Click is on track - start tracking for long press
		s.trackPressed = true
		s.trackPressStart = time.Now()
		if clickPos < thumbStart {
			s.trackPressSide = -1 // Left/up side
		} else {
			s.trackPressSide = 1 // Right/down side
		}
		// Also do immediate scroll
		s.handleTrackClick(clickPos, thumbStart, thumbEnd, int(s.trackLength), s.gtx)
	}
}

// handleDrag handles drag events
func (s *Scrollbar) handleDrag(e pointer.Event) {
	if s.dragging {
		// Update position based on drag position
		var newPos float32
		if s.orientation == Horizontal {
			newPos = e.Position.X / s.trackLength
		} else {
			newPos = e.Position.Y / s.trackLength
		}

		// Adjust for thumb size
		thumbRatio := float32(s.thumbSize) / s.trackLength
		if newPos < thumbRatio/2 {
			newPos = 0
		} else if newPos > 1-thumbRatio/2 {
			newPos = 1
		} else {
			newPos = (newPos - thumbRatio/2) / (1 - thumbRatio)
		}

		if newPos < 0 {
			newPos = 0
		} else if newPos > 1 {
			newPos = 1
		}
		s.position = newPos
		s.changed = true
		s.changeHook(s.position)
	}
}

// handleRelease handles release events
func (s *Scrollbar) handleRelease(e pointer.Event) {
	s.dragging = false
	s.trackPressed = false
}

// handleScroll handles scroll events
func (s *Scrollbar) handleScroll(distance float32) {
	// Adjust position based on scroll distance
	scrollAmount := distance / 100.0 // Scale down the scroll amount
	newPos := s.position + scrollAmount

	if newPos < 0 {
		newPos = 0
	} else if newPos > 1 {
		newPos = 1
	}

	s.position = newPos
	s.changed = true
	s.changeHook(s.position)
}

// Layout renders the scrollbar
func (s *Scrollbar) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	// Use full available space for track
	var trackSpace int
	if s.orientation == Horizontal {
		trackSpace = gtx.Constraints.Min.X
		if trackSpace < gtx.Dp(unit.Dp(50)) {
			trackSpace = gtx.Dp(unit.Dp(50)) // Minimum track space
		}
	} else {
		trackSpace = gtx.Constraints.Min.Y
		if trackSpace < gtx.Dp(unit.Dp(50)) {
			trackSpace = gtx.Dp(unit.Dp(50)) // Minimum track space
		}
	}

	// Handle animation
	s.updateAnimation(gtx)

	// Calculate thumb size and position
	var thumbSize, thumbPos int
	var trackLength float32 = float32(trackSpace)

	if s.orientation == Horizontal {
		thumbSize = int(float32(trackSpace) * s.viewport)
		if thumbSize < gtx.Dp(unit.Dp(20)) {
			thumbSize = gtx.Dp(unit.Dp(20)) // Minimum thumb size
		}
		thumbPos = int(float32(trackSpace-thumbSize) * s.position)
	} else {
		thumbSize = int(float32(trackSpace) * s.viewport)
		if thumbSize < gtx.Dp(unit.Dp(20)) {
			thumbSize = gtx.Dp(unit.Dp(20)) // Minimum thumb size
		}
		thumbPos = int(float32(trackSpace-thumbSize) * s.position)
	}

	// Ensure minimum thumb size for visibility
	if thumbSize < gtx.Dp(unit.Dp(20)) {
		thumbSize = gtx.Dp(unit.Dp(20))
	}

	// Store track dimensions and context for event handlers
	s.trackLength = trackLength
	s.thumbSize = thumbSize
	s.gtx = gtx

	// Register event handler for this scrollbar area
	s.eventHandler.AddToOps(gtx.Ops)
	s.eventHandler.ProcessEvents(gtx)

	// Check for long press on track (1 second)
	if s.trackPressed {
		elapsed := gtx.Now.Sub(s.trackPressStart)
		if elapsed >= 1*time.Second {
			// Long press detected - scroll to end
			var targetPos float32
			if s.trackPressSide == -1 {
				targetPos = 0 // Scroll to start
			} else {
				targetPos = 1 // Scroll to end
			}
			s.startAnimation(targetPos, gtx)
			s.trackPressed = false // Stop tracking to prevent repeated triggers
		} else {
			// Request next frame to continue checking
			gtx.Execute(op.InvalidateCmd{})
		}
	}

	// Layout the scrollbar track only
	if s.orientation == Horizontal {
		// Horizontal layout: just the track
		gtx.Constraints.Min = image.Pt(trackSpace, gtx.Dp(s.width))
		return s.layoutTrack(gtx, th, trackSpace, thumbSize, thumbPos)
	} else {
		// Vertical layout: just the track
		gtx.Constraints.Min = image.Pt(gtx.Dp(s.width), trackSpace)
		return s.layoutTrack(gtx, th, trackSpace, thumbSize, thumbPos)
	}
}

// layoutTrack renders the track and thumb with perfect centering
func (s *Scrollbar) layoutTrack(gtx layout.Context, th *Theme, trackSpace, thumbSize, thumbPos int) layout.Dimensions {
	// Calculate perfect center points
	centerX := gtx.Constraints.Min.X / 2
	centerY := gtx.Constraints.Min.Y / 2

	// Define the full background area (same size as buttons)
	backgroundRect := image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(gtx.Constraints.Min.X, gtx.Constraints.Min.Y),
	}

	// Event handling is now done by the EventHandler in the main Layout method

	// Draw track FIRST to ensure it's visible
	var trackRect image.Rectangle
	if s.orientation == Horizontal {
		// Track line is thin in height, centered within the expanded area, widened 1px downwards
		trackHeight := gtx.Dp(unit.Dp(4))
		trackRect = image.Rectangle{
			Min: image.Pt(0, centerY-trackHeight/2),
			Max: image.Pt(gtx.Constraints.Min.X, centerY+trackHeight/2+1),
		}
	} else {
		// Track line is thin in width, centered within the expanded area, widened 1px rightwards
		trackWidth := gtx.Dp(unit.Dp(4))
		trackRect = image.Rectangle{
			Min: image.Pt(centerX-trackWidth/2, 0),
			Max: image.Pt(centerX+trackWidth/2+1, gtx.Constraints.Min.Y),
		}
	}

	defer clip.Rect(trackRect).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, th.Colors.OutlineVariant()) // Same color as float slider track

	// Draw thumb SECOND - rectangle matching track cross-axis width
	var thumbRect image.Rectangle
	if s.orientation == Horizontal {
		// Horizontal thumb: length along axis = thumbSize, height = track height (thin), widened 1px downwards
		trackHeight := gtx.Dp(unit.Dp(4))
		thumbRect = image.Rectangle{
			Min: image.Pt(thumbPos, centerY-trackHeight/2),
			Max: image.Pt(thumbPos+thumbSize, centerY+trackHeight/2+1),
		}
	} else {
		// Vertical thumb: length along axis = thumbSize, width = track width (thin), widened 1px rightwards
		trackWidth := gtx.Dp(unit.Dp(4))
		thumbRect = image.Rectangle{
			Min: image.Pt(centerX-trackWidth/2, thumbPos),
			Max: image.Pt(centerX+trackWidth/2+1, thumbPos+thumbSize),
		}
	}

	defer clip.Rect(thumbRect).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, th.Colors.Primary()) // Use theme primary color

	// Draw transparent background (this creates the full-width area)
	defer clip.RRect{Rect: backgroundRect, NW: 0, NE: 0, SW: 0, SE: 0}.Push(gtx.Ops).Pop()
	// No paint.Fill here - this is transparent background area

	return layout.Dimensions{Size: gtx.Constraints.Min}
}

// handleTrackClick handles clicks on the track (non-thumb area) with smooth animation
func (s *Scrollbar) handleTrackClick(clickPos, thumbStart, thumbEnd float32, trackSpace int, gtx layout.Context) {
	// Calculate thumb size
	thumbSize := int(float32(trackSpace) * s.viewport)
	if thumbSize < 20 { // Minimum thumb size
		thumbSize = 20
	}

	// Calculate thumb width as a fraction of the track
	thumbWidth := float32(thumbSize) / float32(trackSpace)

	var targetPos float32
	if clickPos < thumbStart {
		// Click is before thumb - scroll backward by thumb width or to start
		targetPos = s.position - thumbWidth
		if targetPos < 0 {
			targetPos = 0
		}
	} else if clickPos > thumbEnd {
		// Click is after thumb - scroll forward by thumb width or to end
		targetPos = s.position + thumbWidth
		if targetPos > 1 {
			targetPos = 1
		}
	} else {
		// Click is on thumb - no animation needed
		return
	}

	// Start smooth animation
	s.startAnimation(targetPos, gtx)
}

// startAnimation starts a smooth animation to the target position
func (s *Scrollbar) startAnimation(targetPos float32, gtx layout.Context) {
	s.animating = true
	s.animStartTime = time.Now()
	s.animStartPos = s.position
	s.animTargetPos = targetPos
	// Request immediate frame update for animation
	gtx.Execute(op.InvalidateCmd{})
}

// updateAnimation updates the animation progress
func (s *Scrollbar) updateAnimation(gtx layout.Context) {
	if !s.animating {
		return
	}

	now := gtx.Now
	elapsed := float32(now.Sub(s.animStartTime).Seconds())

	// Determine animation duration based on target position
	var duration float32
	if s.animTargetPos == 0 || s.animTargetPos == 1 {
		// Fast animation for end positions (100ms)
		duration = float32(0.1)
	} else {
		// Normal animation for track clicks (250ms)
		duration = float32(0.25)
	}

	if elapsed >= duration {
		// Animation complete
		s.position = s.animTargetPos
		s.animating = false
		s.changed = true
		s.changeHook(s.position)
	} else {
		// Calculate progress with easing (ease-out cubic)
		progress := elapsed / duration
		easedProgress := 1 - (1-progress)*(1-progress)*(1-progress)

		// Interpolate position
		s.position = s.animStartPos + (s.animTargetPos-s.animStartPos)*easedProgress
		s.changed = true
		s.changeHook(s.position)

		// Request next frame
		gtx.Execute(op.InvalidateCmd{})
	}
}

// Changed returns true if the position has changed since last call
func (s *Scrollbar) Changed() bool {
	changed := s.changed
	s.changed = false
	return changed
}
