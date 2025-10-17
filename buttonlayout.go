package fromage

import (
	"image"
	"image/color"
	"math"

	"lol.mleku.dev/log"

	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// ButtonLayout provides a button with background, rounded corners, and embedded widget
type ButtonLayout struct {
	// Theme reference
	theme *Theme
	// Background color
	background color.NRGBA
	// Corner radius
	cornerRadius unit.Dp
	// Event handler for handling interactions
	eventHandler *EventHandler
	// Embedded widget
	widget W
	// Corner flags
	corners int
	// Disabled state
	disabled bool
	// Click handler
	onClick func()
	// Disable inking effect
	disableInking bool
	// Internal state for button behavior
	pressed bool
	hovered bool
	clicked bool
}

// NewButtonLayout creates a new button layout
func (t *Theme) NewButtonLayout() *ButtonLayout {
	bl := &ButtonLayout{
		theme:        t,
		background:   t.Colors.Primary(),
		cornerRadius: unit.Dp(float32(t.TextSize) * 0.25), // Scale corner radius based on text size
		corners:      CornerAll,
		disabled:     false,
		pressed:      false,
		hovered:      false,
		clicked:      false,
	}

	// Create event handler with callbacks
	bl.eventHandler = NewEventHandler(func(event string) {
		log.I.F("[ButtonLayout] %s", event)
	}).SetOnClick(func(e pointer.Event) {
		bl.clicked = true
		if bl.onClick != nil {
			bl.onClick()
		}
	}).SetOnPress(func(e pointer.Event) {
		bl.pressed = true
	}).SetOnRelease(func(e pointer.Event) {
		bl.pressed = false
	}).SetOnHover(func(hovered bool) {
		bl.hovered = hovered
	})

	return bl
}

// Background sets the background color
func (b *ButtonLayout) Background(color color.NRGBA) *ButtonLayout {
	b.background = color
	return b
}

// CornerRadius sets the corner radius as a scale factor of the theme text size
func (b *ButtonLayout) CornerRadius(scale float32) *ButtonLayout {
	b.cornerRadius = unit.Dp(float32(b.theme.TextSize) * scale)
	return b
}

// PillRadius sets the corner radius to create a pill shape (half the button height)
func (b *ButtonLayout) PillRadius() *ButtonLayout {
	// For pill shape, we'll use a very large radius that will be clamped to half the height
	// This will be handled in the drawBackground method
	b.cornerRadius = unit.Dp(1000) // Large value to indicate pill shape
	return b
}

// Corners sets which corners should be rounded
func (b *ButtonLayout) Corners(corners int) *ButtonLayout {
	b.corners = corners
	return b
}

// Widget sets the embedded widget
func (b *ButtonLayout) Widget(widget W) *ButtonLayout {
	b.widget = widget
	return b
}

// Disabled sets the disabled state
func (b *ButtonLayout) Disabled(disabled bool) *ButtonLayout {
	b.disabled = disabled
	return b
}

// OnClick sets the click handler
func (b *ButtonLayout) OnClick(handler func()) *ButtonLayout {
	b.onClick = handler
	return b
}

// DisableInking disables the ink animation effect
func (b *ButtonLayout) DisableInking(disable bool) *ButtonLayout {
	b.disableInking = disable
	return b
}

// Clicked returns true if the button was clicked
func (b *ButtonLayout) Clicked(g C) bool {
	clicked := b.clicked
	b.clicked = false // Reset after checking
	if clicked {
		log.I.F("[TRACE] Button clicked at position: %v", g.Now)
	}
	return clicked
}

// Hovered returns true if the button is being hovered
func (b *ButtonLayout) Hovered() bool {
	return b.hovered
}

// Pressed returns true if the button is being pressed
func (b *ButtonLayout) Pressed() bool {
	return b.pressed
}

// Layout renders the button layout
func (b *ButtonLayout) Layout(g C) D {
	// Handle disabled state
	if b.disabled {
		g = g.Disabled()
	}

	// Get the dimensions of the embedded widget
	var widgetDims D
	if b.widget != nil {
		widgetDims = b.widget(g)
	} else {
		// If no widget, use minimum button size
		widgetDims = D{Size: image.Pt(100, 40)}
	}

	// Respect the constraints passed to the button instead of enforcing minimum size
	// Only apply minimum size if no explicit constraints are provided
	if g.Constraints.Min.X == 0 && g.Constraints.Min.Y == 0 {
		// No explicit constraints - apply default minimum size
		minWidth := int(float32(b.theme.TextSize) * 1)
		minHeight := int(float32(b.theme.TextSize) * 1)
		minSize := image.Pt(minWidth, minHeight)
		if widgetDims.Size.X < minSize.X {
			widgetDims.Size.X = minSize.X
		}
		if widgetDims.Size.Y < minSize.Y {
			widgetDims.Size.Y = minSize.Y
		}
	} else {
		// Use the explicit constraints provided - completely override widget size
		widgetDims.Size = g.Constraints.Min
		// Also override the dimensions to match the constraints exactly
		widgetDims = D{Size: g.Constraints.Min}
	}

	// Add semantic information for accessibility
	semantic.Button.Add(g.Ops)

	// Register event handler for this button area
	b.eventHandler.AddToOps(g.Ops)
	b.eventHandler.ProcessEvents(g)

	// Draw background - use exact constraints if provided
	finalSize := widgetDims.Size
	if g.Constraints.Min.X > 0 && g.Constraints.Min.Y > 0 {
		finalSize = g.Constraints.Min
	}
	bgDims := b.drawBackground(g, finalSize)

	// Draw content on top - use layout.Center to properly queue the widget
	if b.widget != nil {
		g.Constraints.Min = finalSize
		layout.Center.Layout(g, b.widget)
	}

	return bgDims
}

// drawBackground draws the button background with rounded corners and animations
func (b *ButtonLayout) drawBackground(g C, size image.Point) D {
	// Adjust background color based on state
	bgColor := b.background
	if b.disabled {
		// Make background more transparent when disabled
		bgColor.A = bgColor.A / 2
	} else if b.Hovered() {
		// Lighten background on hover
		bgColor = b.hoveredColor(bgColor)
	}

	// Create rounded rectangle clip
	rect := image.Rectangle{Max: size}

	// Calculate corner radius, handling pill shape and perfect circles
	radius := float32(g.Dp(b.cornerRadius))
	if b.cornerRadius > unit.Dp(500) {
		// This is a pill shape - use half the button height
		radius = float32(size.Y) / 2
	} else if b.cornerRadius > unit.Dp(100) {
		// This is a perfect circle - use half the smaller dimension
		if size.X < size.Y {
			radius = float32(size.X) / 2
		} else {
			radius = float32(size.Y) / 2
		}
	}

	rrect := clip.RRect{
		Rect: rect,
		NW:   int(b.ifCorner(radius, b.corners&CornerNW)),
		NE:   int(b.ifCorner(radius, b.corners&CornerNE)),
		SW:   int(b.ifCorner(radius, b.corners&CornerSW)),
		SE:   int(b.ifCorner(radius, b.corners&CornerSE)),
	}

	// Push clip and fill background
	defer rrect.Push(g.Ops).Pop()
	paint.Fill(g.Ops, bgColor)

	// Draw ink animations for press history (unless disabled)
	// Note: Ink animations are disabled when using EventHandler instead of widget.Clickable
	// This could be re-implemented by tracking press events in the EventHandler

	return D{Size: size}
}

// ifCorner returns the radius if the corner flag is set, otherwise 0
func (b *ButtonLayout) ifCorner(radius float32, corner int) float32 {
	if corner != 0 {
		return radius
	}
	return 0
}

// hoveredColor lightens the color for hover effect
func (b *ButtonLayout) hoveredColor(c color.NRGBA) color.NRGBA {
	// Lighten the color by increasing RGB values
	return color.NRGBA{
		R: uint8(min(255, int(c.R)+30)),
		G: uint8(min(255, int(c.G)+30)),
		B: uint8(min(255, int(c.B)+30)),
		A: c.A,
	}
}

// drawInk draws the animated ink effect for button presses
func (b *ButtonLayout) drawInk(g C, press widget.Press, size image.Point) {
	// Animation durations (matching Gio's material design)
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := g.Now
	t := float32(now.Sub(press.Start).Seconds())

	end := press.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out
		end = now
	}

	endt := float32(end.Sub(press.Start).Seconds())

	// Compute the fade-in/out position in [0;1]
	var alphat float32
	{
		var haste float32
		if press.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}

		// Fade out
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1]
	sizet := t
	if press.Cancelled {
		// Freeze expansion of cancelled presses
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in
	if !press.End.IsZero() || sizet <= 1.0 {
		g.Execute(op.InvalidateCmd{})
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5
	t2 := alphat * 2
	// Bézier ease-in curve
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)

	// Calculate ink size and position
	inkSize := size.X
	if h := size.Y; h > inkSize {
		inkSize = h
	}
	// Cover the entire button and apply curve values to size and color
	inkSize = int(float32(inkSize) * 2 * float32(math.Sqrt(2)) * sizeBezier)
	alpha := 0.7 * alphaBezier

	// Create ink color (white with alpha)
	inkColor := color.NRGBA{
		R: 0xff,
		G: 0xff,
		B: 0xff,
		A: uint8(alpha * 0xff),
	}

	// Draw circular ink effect at press position
	inkRadius := inkSize / 2
	inkRect := image.Rectangle{
		Min: image.Pt(press.Position.X-inkRadius, press.Position.Y-inkRadius),
		Max: image.Pt(press.Position.X+inkRadius, press.Position.Y+inkRadius),
	}

	defer clip.Ellipse(inkRect).Push(g.Ops).Pop()
	paint.Fill(g.Ops, inkColor)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Convenience methods for common button patterns

// PrimaryButton creates a button with primary color
func (t *Theme) PrimaryButton(widget W) *ButtonLayout {
	return t.NewButtonLayout().
		Background(t.Colors.Primary()).
		Widget(widget)
}

// SecondaryButton creates a button with secondary color
func (t *Theme) SecondaryButton(widget W) *ButtonLayout {
	return t.NewButtonLayout().
		Background(t.Colors.Secondary()).
		Widget(widget)
}

// SurfaceButton creates a button with surface color
func (t *Theme) SurfaceButton(widget W) *ButtonLayout {
	return t.NewButtonLayout().
		Background(t.Colors.Surface()).
		Widget(widget)
}

// ErrorButton creates a button with error color
func (t *Theme) ErrorButton(widget W) *ButtonLayout {
	return t.NewButtonLayout().
		Background(t.Colors.Error()).
		Widget(widget)
}

// RoundedButton creates a button with larger corner radius
func (t *Theme) RoundedButton(widget W) *ButtonLayout {
	return t.NewButtonLayout().
		CornerRadius(0.5). // Scale factor for rounded corners
		Widget(widget)
}

// PillButton creates a pill-shaped button
func (t *Theme) PillButton(widget W) *ButtonLayout {
	return t.NewButtonLayout().
		PillRadius(). // Creates true pill shape with rounded sides
		Widget(widget)
}

// TextButton creates a button with text content
func (t *Theme) TextButton(textContent string) *ButtonLayout {
	return t.NewButtonLayout().
		Widget(func(g C) D {
			return t.Body1(textContent).
				Color(t.Colors.OnPrimary()).
				Alignment(text.Middle).
				Layout(g)
		})
}

// IconButton creates a button with an icon (placeholder for now)
func (t *Theme) IconButton(icon string) *ButtonLayout {
	return t.NewButtonLayout().
		PillRadius(). // Creates circular/pill shape for icon
		Widget(func(g C) D {
			// Placeholder for icon - could be extended with actual icon support
			return t.Caption(icon).
				Color(t.Colors.OnPrimary()).
				Alignment(text.Middle).
				Layout(g)
		})
}
