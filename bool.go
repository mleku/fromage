package fromage

import (
	"image"
	"image/color"
	"time"

	"gio.mleku.dev/io/semantic"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget"
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
	animationProgress float32   // 0.0 = off, 1.0 = on
	animationStart    time.Time // When the current animation started
	isAnimating       bool      // Whether an animation is currently in progress
	// Color transition state
	colorTransitionStart time.Time   // When color transition started
	isColorTransitioning bool        // Whether color transition is in progress
	oldBackground        color.NRGBA // Previous background color
	oldForeground        color.NRGBA // Previous foreground color
}

// NewBool creates a new boolean widget with its own dedicated clickable
func (t *Theme) NewBool(value bool) *Bool {
	// Material Design switch dimensions
	width := unit.Dp(36)     // Standard switch width
	height := unit.Dp(20)    // Standard switch height
	thumbSize := unit.Dp(16) // Thumb circle size

	// Initialize animation progress based on initial value
	var initialProgress float32
	if value {
		initialProgress = 1.0 // Start in "on" position
	} else {
		initialProgress = 0.0 // Start in "off" position
	}

	return &Bool{
		theme:                t,
		value:                value,
		clickable:            &widget.Clickable{}, // Each bool gets its own dedicated clickable
		changed:              false,
		onChange:             func(b bool) {},
		background:           color.NRGBA{R: 128, G: 128, B: 128, A: 255}, // Dim gray for off state
		foreground:           t.Colors.Surface(),                          // White thumb
		cornerRadius:         unit.Dp(10),                                 // Half of height for pill shape
		width:                width,
		height:               height,
		thumbSize:            thumbSize,
		animationProgress:    initialProgress,
		animationStart:       time.Time{},
		isAnimating:          false,
		colorTransitionStart: time.Time{},
		isColorTransitioning: false,
		oldBackground:        color.NRGBA{R: 128, G: 128, B: 128, A: 255},
		oldForeground:        t.Colors.Surface(),
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
		// Start animation when value changes
		b.startAnimation(g.Now)
	}

	// Update animation progress based on time
	b.updateAnimation(g)

	// Update color transition progress based on time
	b.updateColorTransition(g)

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

		return D{Size: minSize}
	})
}

// drawTrack draws the switch track (background)
func (b *Bool) drawTrack(g C, size image.Point) {
	// Determine track color based on state and theme
	var trackColor color.NRGBA
	if b.value {
		// When on, use the configured background color (text color)
		trackColor = b.background
	} else {
		// When off, use 50% opacity of the text color
		textColor := b.theme.Colors.OnBackground()
		trackColor = color.NRGBA{
			R: textColor.R,
			G: textColor.G,
			B: textColor.B,
			A: textColor.A / 2, // 50% opacity
		}
	}

	// Debug: Print the track color to see what's happening
	// fmt.Printf("Switch state: %v, trackColor: %+v, b.background: %+v\n", b.value, trackColor, b.background)

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

	// Determine thumb color based on theme
	var thumbColor color.NRGBA
	if b.theme.IsLight() {
		// Light mode: white thumb
		thumbColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		// Dark mode: black thumb
		thumbColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	}

	// Draw thumb circle
	defer clip.Ellipse(thumbRect).Push(g.Ops).Pop()
	paint.Fill(g.Ops, thumbColor)
}

// startAnimation begins a new animation
func (b *Bool) startAnimation(now time.Time) {
	b.animationStart = now
	b.isAnimating = true
}

// startColorTransition begins a color transition animation
func (b *Bool) startColorTransition(now time.Time) {
	b.oldBackground = b.background
	b.oldForeground = b.foreground
	b.colorTransitionStart = now
	b.isColorTransitioning = true
}

// updateAnimation updates the animation progress based on elapsed time
func (b *Bool) updateAnimation(g C) {
	if !b.isAnimating {
		return
	}

	const animationDuration = 250 * time.Millisecond
	elapsed := g.Now.Sub(b.animationStart)

	if elapsed >= animationDuration {
		// Animation complete
		if b.value {
			b.animationProgress = 1.0
		} else {
			b.animationProgress = 0.0
		}
		b.isAnimating = false
		return
	}

	// Calculate progress (0.0 to 1.0)
	progress := float32(elapsed) / float32(animationDuration)

	// Apply easing function (ease-out for smooth deceleration)
	progress = 1.0 - (1.0-progress)*(1.0-progress)

	if b.value {
		// Animating to ON state (0.0 -> 1.0)
		b.animationProgress = progress
	} else {
		// Animating to OFF state (1.0 -> 0.0)
		b.animationProgress = 1.0 - progress
	}

	// Request invalidation to continue animation
	g.Execute(op.InvalidateCmd{})
}

// updateColorTransition updates the color transition progress based on elapsed time
func (b *Bool) updateColorTransition(g C) {
	if !b.isColorTransitioning {
		return
	}

	const colorTransitionDuration = 250 * time.Millisecond
	elapsed := g.Now.Sub(b.colorTransitionStart)

	if elapsed >= colorTransitionDuration {
		// Color transition complete
		b.isColorTransitioning = false
		return
	}

	// Calculate progress (0.0 to 1.0)
	progress := float32(elapsed) / float32(colorTransitionDuration)

	// Apply easing function (ease-out for smooth deceleration)
	progress = 1.0 - (1.0-progress)*(1.0-progress)

	// Interpolate colors
	b.background = b.interpolateColor(b.oldBackground, b.background, progress)
	b.foreground = b.interpolateColor(b.oldForeground, b.foreground, progress)

	// Request invalidation to continue color transition
	g.Execute(op.InvalidateCmd{})
}

// interpolateColor interpolates between two colors based on progress (0.0 to 1.0)
func (b *Bool) interpolateColor(from, to color.NRGBA, progress float32) color.NRGBA {
	return color.NRGBA{
		R: uint8(float32(from.R) + (float32(to.R)-float32(from.R))*progress),
		G: uint8(float32(from.G) + (float32(to.G)-float32(from.G))*progress),
		B: uint8(float32(from.B) + (float32(to.B)-float32(from.B))*progress),
		A: uint8(float32(from.A) + (float32(to.A)-float32(from.A))*progress),
	}
}

// UpdateThemeColors updates the switch colors to match the current theme and starts color transition
func (b *Bool) UpdateThemeColors(now time.Time) {
	// Store current colors as old colors for transition
	b.oldBackground = b.background
	b.oldForeground = b.foreground

	// Update to new theme colors
	b.background = b.theme.Colors.OnBackground() // Update to new text color

	// Update thumb color based on new theme
	if b.theme.IsLight() {
		// Light mode: white thumb
		b.foreground = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		// Dark mode: black thumb
		b.foreground = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	}

	// Start color transition
	b.startColorTransition(now)
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

// Convenience methods for common boolean widget patterns

// PrimaryBool creates a boolean widget with primary color as track when active and background color thumb
func (t *Theme) PrimaryBool(value bool) *Bool {
	// Use explicit white color for thumb
	thumbColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White thumb
	return t.NewBool(value).
		Background(t.Colors.Primary()).
		Foreground(thumbColor)
}

// SecondaryBool creates a boolean widget with secondary color as track when active and background color thumb
func (t *Theme) SecondaryBool(value bool) *Bool {
	// Use explicit white color for thumb
	thumbColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White thumb
	return t.NewBool(value).
		Background(t.Colors.Secondary()).
		Foreground(thumbColor)
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

// Switch creates a switch-style boolean widget with proper light/dark mode colors
func (t *Theme) Switch(value bool) *Bool {
	var thumbColor color.NRGBA
	if t.IsLight() {
		// Light mode: white thumb
		thumbColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		// Dark mode: black thumb
		thumbColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	}
	return t.NewBool(value).
		Background(t.Colors.OnBackground()). // Text color as track
		Foreground(thumbColor)
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
