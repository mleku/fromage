package fromage

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Border lays out a widget and draws a border inside it.
type Border struct {
	theme        *Theme
	color        color.NRGBA
	cornerRadius unit.Dp
	width        unit.Dp
	widget       W
}

// NewBorder creates a new border widget
func (t *Theme) NewBorder() *Border {
	return &Border{
		theme:        t,
		color:        t.Colors.Outline(),
		cornerRadius: unit.Dp(float32(t.TextSize) * 0.25),  // Scale corner radius based on text size
		width:        unit.Dp(float32(t.TextSize) * 0.125), // Scale width based on text size
	}
}

// Color sets the color to render the border in
func (b *Border) Color(color color.NRGBA) *Border {
	b.color = color
	return b
}

// CornerRadius sets the radius of the curve on the corners
func (b *Border) CornerRadius(radius unit.Dp) *Border {
	b.cornerRadius = radius
	return b
}

// Width sets the width of the border line
func (b *Border) Width(width unit.Dp) *Border {
	b.width = width
	return b
}

// Widget sets the widget to be bordered
func (b *Border) Widget(widget W) *Border {
	b.widget = widget
	return b
}

// Layout renders the border around the widget
func (b *Border) Layout(g C) D {
	if b.widget == nil {
		return D{}
	}

	// Layout the widget first to get its dimensions
	dims := b.widget(g)
	sz := layout.FPt(dims.Size)

	// Calculate border properties
	rr := float32(g.Dp(b.cornerRadius))
	width := float32(g.Dp(b.width))

	// Adjust size to account for border width
	sz.X -= width
	sz.Y -= width

	// Create rectangle for border
	r := image.Rectangle{Max: image.Pt(int(sz.X), int(sz.Y))}
	r = r.Add(image.Pt(int(width*0.5), int(width*0.5)))

	// Draw the border
	paint.FillShape(g.Ops,
		b.color,
		clip.Stroke{
			Path:  clip.RRect{Rect: r, NW: int(rr), NE: int(rr), SW: int(rr), SE: int(rr)}.Path(g.Ops),
			Width: width,
		}.Op(),
	)

	return dims
}

// Convenience methods for common border styles

// BorderPrimary creates a border with primary color
func (t *Theme) BorderPrimary() *Border {
	return t.NewBorder().Color(t.Colors.Primary())
}

// BorderSecondary creates a border with secondary color
func (t *Theme) BorderSecondary() *Border {
	return t.NewBorder().Color(t.Colors.Secondary())
}

// BorderOutline creates a border with outline color
func (t *Theme) BorderOutline() *Border {
	return t.NewBorder().Color(t.Colors.Outline())
}

// BorderError creates a border with error color
func (t *Theme) BorderError() *Border {
	return t.NewBorder().Color(t.Colors.Error())
}

// BorderSurface creates a border with surface color
func (t *Theme) BorderSurface() *Border {
	return t.NewBorder().Color(t.Colors.Surface())
}

// BorderRounded creates a border with rounded corners
func (t *Theme) BorderRounded() *Border {
	return t.NewBorder().CornerRadius(unit.Dp(float32(t.TextSize) * 0.5))
}

// BorderThick creates a thick border
func (t *Theme) BorderThick() *Border {
	return t.NewBorder().Width(unit.Dp(float32(t.TextSize) * 0.25))
}

// BorderThin creates a thin border
func (t *Theme) BorderThin() *Border {
	return t.NewBorder().Width(unit.Dp(float32(t.TextSize) * 0.0625))
}
