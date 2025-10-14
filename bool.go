package fromage

import (
	"image"
	"image/color"

	"gioui.org/io/semantic"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

// BoolHook is a function type for handling boolean value changes
type BoolHook func(b bool)

// Bool represents a boolean toggle widget
type Bool struct {
	// Theme reference
	theme *Theme
	// Current boolean value
	value bool
	// Clickable widget for handling interactions
	clickable *widget.Clickable
	// Whether the value has changed since last check
	changed bool
	// Callback function for value changes
	onChange BoolHook
	// Visual styling
	background   color.NRGBA
	foreground   color.NRGBA
	cornerRadius unit.Dp
	// Size constraints
	size unit.Dp
}

// NewBool creates a new boolean widget with its own dedicated clickable
func (t *Theme) NewBool(value bool) *Bool {
	return &Bool{
		theme:        t,
		value:        value,
		clickable:    &widget.Clickable{}, // Each bool gets its own dedicated clickable
		changed:      false,
		onChange:     func(b bool) {},
		background:   t.Colors.Primary(),
		foreground:   t.Colors.OnPrimary(),
		cornerRadius: unit.Dp(float32(t.TextSize) * 0.25),
		size:         unit.Dp(float32(t.TextSize) * 1.5), // Scale with text size
	}
}

// Value sets the boolean value
func (b *Bool) Value(value bool) *Bool {
	if b.value != value {
		b.value = value
		b.changed = true
	}
	return b
}

// GetValue returns the current boolean value
func (b *Bool) GetValue() bool {
	return b.value
}

// SetOnChange sets the callback function for value changes
func (b *Bool) SetOnChange(fn BoolHook) *Bool {
	b.onChange = fn
	return b
}

// Background sets the background color
func (b *Bool) Background(color color.NRGBA) *Bool {
	b.background = color
	return b
}

// Foreground sets the foreground color
func (b *Bool) Foreground(color color.NRGBA) *Bool {
	b.foreground = color
	return b
}

// CornerRadius sets the corner radius
func (b *Bool) CornerRadius(radius unit.Dp) *Bool {
	b.cornerRadius = radius
	return b
}

// Size sets the widget size
func (b *Bool) Size(size unit.Dp) *Bool {
	b.size = size
	return b
}

// Changed returns true if the value has changed since the last call
func (b *Bool) Changed() bool {
	changed := b.changed
	b.changed = false
	return changed
}

// Clicked returns true if the widget was clicked
func (b *Bool) Clicked(g C) bool {
	return b.clickable.Clicked(g)
}

// Hovered returns true if the widget is being hovered
func (b *Bool) Hovered() bool {
	return b.clickable.Hovered()
}

// Pressed returns true if the widget is being pressed
func (b *Bool) Pressed() bool {
	return b.clickable.Pressed()
}

// Layout renders the boolean widget
func (b *Bool) Layout(g C) D {
	// Handle click events BEFORE layout (like regular buttons)
	if b.clickable.Clicked(g) {
		b.value = !b.value
		b.changed = true
		if b.onChange != nil {
			b.onChange(b.value)
		}
	}

	// Calculate size
	size := g.Dp(b.size)
	minSize := image.Pt(size, size)

	// Create the layout using the clickable's Layout method
	return b.clickable.Layout(g, func(g C) D {
		// Add semantic information for accessibility
		semantic.Button.Add(g.Ops)

		// Draw background
		bgDims := b.drawBackground(g, minSize)

		// Draw checkmark or content
		b.drawContent(g, minSize)

		return bgDims
	})
}

// drawBackground draws the widget background
func (b *Bool) drawBackground(g C, size image.Point) D {
	// Adjust background color based on state
	bgColor := b.background
	if b.Hovered() {
		// Lighten background on hover
		bgColor = b.hoveredColor(bgColor)
	}

	// Create rounded rectangle clip
	rect := image.Rectangle{Max: size}
	radius := float32(g.Dp(b.cornerRadius))

	// Create rounded rectangle
	rrect := clip.RRect{
		Rect: rect,
		NW:   int(radius),
		NE:   int(radius),
		SW:   int(radius),
		SE:   int(radius),
	}

	// Push clip and fill background
	defer rrect.Push(g.Ops).Pop()
	paint.Fill(g.Ops, bgColor)

	// Draw ink animations for press history
	for _, press := range b.clickable.History() {
		b.drawInk(g, press, size)
	}

	return D{Size: size}
}

// drawContent draws the checkmark or other visual indicator
func (b *Bool) drawContent(g C, size image.Point) {
	if !b.value {
		return // Don't draw anything when false
	}

	// Draw a simple checkmark using lines
	// This is a basic implementation - could be enhanced with better graphics
	center := image.Pt(size.X/2, size.Y/2)
	checkSize := int(float32(size.X) * 0.4) // 40% of widget size

	// Create a simple checkmark path
	// This is a placeholder - in a real implementation you might use
	// a proper vector graphics library or pre-rendered icon

	// For now, we'll draw a simple filled circle to indicate "on" state
	circleRadius := checkSize / 2
	circleRect := image.Rectangle{
		Min: image.Pt(center.X-circleRadius, center.Y-circleRadius),
		Max: image.Pt(center.X+circleRadius, center.Y+circleRadius),
	}

	defer clip.Ellipse(circleRect).Push(g.Ops).Pop()
	paint.Fill(g.Ops, b.foreground)
}

// hoveredColor lightens the color for hover effect
func (b *Bool) hoveredColor(c color.NRGBA) color.NRGBA {
	return color.NRGBA{
		R: uint8(min(255, int(c.R)+30)),
		G: uint8(min(255, int(c.G)+30)),
		B: uint8(min(255, int(c.B)+30)),
		A: c.A,
	}
}

// drawInk draws the animated ink effect for button presses
func (b *Bool) drawInk(g C, press widget.Press, size image.Point) {
	// Animation durations (matching Gio's material design)
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := g.Now
	t := float32(now.Sub(press.Start).Seconds())

	end := press.End
	if end.IsZero() {
		end = now
	}

	endt := float32(end.Sub(press.Start).Seconds())

	// Compute the fade-in/out position in [0;1]
	var alphat float32
	{
		var haste float32
		if press.Cancelled {
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
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1]
	sizet := t
	if press.Cancelled {
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
		alphat = 1.0 - alphat
	}

	t2 := alphat * 2
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)

	// Calculate ink size and position
	inkSize := size.X
	if h := size.Y; h > inkSize {
		inkSize = h
	}
	inkSize = int(float32(inkSize) * 2 * sizeBezier)
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

// Convenience methods for common boolean widget patterns

// PrimaryBool creates a boolean widget with primary colors
func (t *Theme) PrimaryBool(value bool) *Bool {
	return t.NewBool(value).
		Background(t.Colors.Primary()).
		Foreground(t.Colors.OnPrimary())
}

// SecondaryBool creates a boolean widget with secondary colors
func (t *Theme) SecondaryBool(value bool) *Bool {
	return t.NewBool(value).
		Background(t.Colors.Secondary()).
		Foreground(t.Colors.OnSecondary())
}

// SurfaceBool creates a boolean widget with surface colors
func (t *Theme) SurfaceBool(value bool) *Bool {
	return t.NewBool(value).
		Background(t.Colors.Surface()).
		Foreground(t.Colors.OnSurface())
}

// Checkbox creates a checkbox-style boolean widget
func (t *Theme) Checkbox(value bool) *Bool {
	return t.NewBool(value).
		Background(t.Colors.Surface()).
		Foreground(t.Colors.Primary()).
		CornerRadius(unit.Dp(2)) // Small corner radius for checkbox
}

// Switch creates a switch-style boolean widget (pill-shaped)
func (t *Theme) Switch(value bool) *Bool {
	return t.NewBool(value).
		Background(t.Colors.Surface()).
		Foreground(t.Colors.Primary()).
		CornerRadius(unit.Dp(1000)) // Large radius for pill shape
}
