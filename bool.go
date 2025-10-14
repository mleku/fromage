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

// Bool represents a boolean toggle widget (Material Design switch)
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
	background   color.NRGBA // Background color of the switch track
	foreground   color.NRGBA // Color of the thumb circle
	cornerRadius unit.Dp     // Corner radius for the track
	// Size constraints
	width     unit.Dp // Width of the switch track
	height    unit.Dp // Height of the switch track
	thumbSize unit.Dp // Size of the thumb circle
	// Animation state
	animationProgress float32 // 0.0 = off, 1.0 = on
}

// NewBool creates a new boolean widget with its own dedicated clickable
func (t *Theme) NewBool(value bool) *Bool {
	// Material Design switch dimensions
	width := unit.Dp(36)     // Standard switch width
	height := unit.Dp(20)    // Standard switch height
	thumbSize := unit.Dp(16) // Thumb circle size

	return &Bool{
		theme:             t,
		value:             value,
		clickable:         &widget.Clickable{}, // Each bool gets its own dedicated clickable
		changed:           false,
		onChange:          func(b bool) {},
		background:        color.NRGBA{R: 128, G: 128, B: 128, A: 255}, // Dim gray for off state
		foreground:        t.Colors.Surface(),                          // White thumb
		cornerRadius:      unit.Dp(10),                                 // Half of height for pill shape
		width:             width,
		height:            height,
		thumbSize:         thumbSize,
		animationProgress: 0.0,
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

// Size sets the widget size (deprecated, use Width/Height instead)
func (b *Bool) Size(size unit.Dp) *Bool {
	b.width = size * 2 // Maintain aspect ratio
	b.height = size
	b.thumbSize = size * 0.8
	return b
}

// Width sets the switch track width
func (b *Bool) Width(width unit.Dp) *Bool {
	b.width = width
	return b
}

// Height sets the switch track height
func (b *Bool) Height(height unit.Dp) *Bool {
	b.height = height
	b.cornerRadius = height / 2 // Maintain pill shape
	return b
}

// ThumbSize sets the thumb circle size
func (b *Bool) ThumbSize(size unit.Dp) *Bool {
	b.thumbSize = size
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

// Layout renders the Material Design switch
func (b *Bool) Layout(g C) D {
	// Handle click events BEFORE layout
	if b.clickable.Clicked(g) {
		b.value = !b.value
		b.changed = true
		if b.onChange != nil {
			b.onChange(b.value)
		}
	}

	// Update animation progress
	if b.value {
		if b.animationProgress < 1.0 {
			b.animationProgress = minFloat32(1.0, b.animationProgress+0.1) // Smooth animation
		}
	} else {
		if b.animationProgress > 0.0 {
			b.animationProgress = maxFloat32(0.0, b.animationProgress-0.1) // Smooth animation
		}
	}

	// Calculate dimensions
	width := g.Dp(b.width)
	height := g.Dp(b.height)
	thumbSize := g.Dp(b.thumbSize)

	minSize := image.Pt(width, height)

	// Create the layout using the clickable's Layout method
	return b.clickable.Layout(g, func(g C) D {
		// Add semantic information for accessibility
		semantic.Button.Add(g.Ops)

		// Draw the switch track (background)
		b.drawTrack(g, minSize)

		// Draw the thumb (circle)
		b.drawThumb(g, minSize, thumbSize)

		// Draw ink animations for press history
		for _, press := range b.clickable.History() {
			b.drawInk(g, press, minSize)
		}

		return D{Size: minSize}
	})
}

// drawTrack draws the switch track (background)
func (b *Bool) drawTrack(g C, size image.Point) {
	// Determine track color based on state
	trackColor := b.background
	if b.value {
		// When on, use primary color with opacity based on animation
		primary := b.theme.Colors.Primary()
		trackColor = color.NRGBA{
			R: primary.R,
			G: primary.G,
			B: primary.B,
			A: uint8(128 + (127 * b.animationProgress)), // Fade from dim to full
		}
	}

	// Create rounded rectangle clip for the track
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

	// Push clip and fill track
	defer rrect.Push(g.Ops).Pop()
	paint.Fill(g.Ops, trackColor)
}

// drawThumb draws the thumb circle
func (b *Bool) drawThumb(g C, size image.Point, thumbSize int) {
	// Calculate thumb position based on animation progress
	// Left position when off (progress = 0), right position when on (progress = 1)
	padding := (size.Y - thumbSize) / 2
	leftPos := padding
	rightPos := size.X - thumbSize - padding

	// Interpolate position based on animation progress
	thumbX := int(float32(leftPos) + float32(rightPos-leftPos)*b.animationProgress)
	thumbY := padding

	// Create thumb circle
	thumbRect := image.Rectangle{
		Min: image.Pt(thumbX, thumbY),
		Max: image.Pt(thumbX+thumbSize, thumbY+thumbSize),
	}

	// Draw thumb circle
	defer clip.Ellipse(thumbRect).Push(g.Ops).Pop()
	paint.Fill(g.Ops, b.foreground)
}

// hoveredColor lightens the color for hover effect
func (b *Bool) hoveredColor(c color.NRGBA) color.NRGBA {
	return color.NRGBA{
		R: uint8(minInt(255, int(c.R)+30)),
		G: uint8(minInt(255, int(c.G)+30)),
		B: uint8(minInt(255, int(c.B)+30)),
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

	// Create a clip mask using the same rounded rectangle as the track
	trackRect := image.Rectangle{Max: size}
	radius := float32(g.Dp(b.cornerRadius))
	trackClip := clip.RRect{
		Rect: trackRect,
		NW:   int(radius),
		NE:   int(radius),
		SW:   int(radius),
		SE:   int(radius),
	}

	// Apply the track clip mask to constrain ink to the switch area
	defer trackClip.Push(g.Ops).Pop()
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

// Switch creates a switch-style boolean widget (Material Design switch)
func (t *Theme) Switch(value bool) *Bool {
	return t.NewBool(value).
		Background(color.NRGBA{R: 128, G: 128, B: 128, A: 255}). // Dim gray for off state
		Foreground(t.Colors.Surface())                           // White thumb
}

// SwitchWithColor creates a switch with custom background color
func (t *Theme) SwitchWithColor(value bool, bgColor color.NRGBA) *Bool {
	return t.NewBool(value).
		Background(bgColor).
		Foreground(t.Colors.Surface()) // White thumb
}

// Helper functions for animation
func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
