package fromage

import (
	"image"
	"image/color"

	"gio.mleku.dev/layout"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
)

// Fill provides a widget that fills the background with a specified color and corner radius
type Fill struct {
	// Theme reference
	theme *Theme
	// Background color
	color color.NRGBA
	// Embedded widget
	widget W
	// Direction for layout
	direction layout.Direction
	// Corner radius
	cornerRadius float32
	// Corner flags
	corners int
}

// Corner flags for rounded corners
const (
	CornerNW = 1 << iota // North West
	CornerNE             // North East
	CornerSW             // South West
	CornerSE             // South East
)

// CornerAll sets all corners
const CornerAll = CornerNW | CornerNE | CornerSW | CornerSE

// NewFill creates a new fill widget
func (t *Theme) NewFill(color color.NRGBA, widget W) *Fill {
	return &Fill{
		theme:        t,
		color:        color,
		widget:       widget,
		direction:    layout.Center,
		cornerRadius: 0,
		corners:      0,
	}
}

// NewFillWithRadius creates a new fill widget with corner radius
func (t *Theme) NewFillWithRadius(color color.NRGBA, radius float32, corners int, widget W) *Fill {
	return &Fill{
		theme:        t,
		color:        color,
		widget:       widget,
		direction:    layout.Center,
		cornerRadius: radius,
		corners:      corners,
	}
}

// Color sets the background color
func (f *Fill) Color(color color.NRGBA) *Fill {
	f.color = color
	return f
}

// Widget sets the embedded widget
func (f *Fill) Widget(widget W) *Fill {
	f.widget = widget
	return f
}

// Direction sets the layout direction
func (f *Fill) Direction(direction layout.Direction) *Fill {
	f.direction = direction
	return f
}

// CornerRadius sets the corner radius
func (f *Fill) CornerRadius(radius float32) *Fill {
	f.cornerRadius = radius
	return f
}

// Corners sets which corners should be rounded
func (f *Fill) Corners(corners int) *Fill {
	f.corners = corners
	return f
}

// Layout renders the fill widget
func (f *Fill) Layout(g C) D {
	// Get the dimensions of the embedded widget
	var widgetDims D
	if f.widget != nil {
		widgetDims = f.widget(g)
	} else {
		// If no widget, use the available space
		widgetDims = D{Size: g.Constraints.Max}
	}

	// Fill the background
	f.fillBackground(g, widgetDims.Size)

	// Layout the embedded widget
	if f.widget != nil {
		return f.direction.Layout(g, f.widget)
	}

	return widgetDims
}

// fillBackground draws the background fill with rounded corners
func (f *Fill) fillBackground(g C, bounds image.Point) {
	rect := image.Rectangle{
		Max: bounds,
	}

	// Create rounded rectangle clip
	rrect := clip.RRect{
		Rect: rect,
		NW:   int(f.ifCorner(f.cornerRadius, f.corners&CornerNW)),
		NE:   int(f.ifCorner(f.cornerRadius, f.corners&CornerNE)),
		SW:   int(f.ifCorner(f.cornerRadius, f.corners&CornerSW)),
		SE:   int(f.ifCorner(f.cornerRadius, f.corners&CornerSE)),
	}

	// Push clip and fill
	defer rrect.Push(g.Ops).Pop()
	paint.Fill(g.Ops, f.color)
}

// ifCorner returns the radius if the corner flag is set, otherwise 0
func (f *Fill) ifCorner(radius float32, corner int) float32 {
	if corner != 0 {
		return radius
	}
	return 0
}

// Convenience methods for common fill patterns

// FillPrimary creates a fill with primary color
func (t *Theme) FillPrimary(widget W) *Fill {
	return t.NewFill(t.Colors.Primary(), widget)
}

// FillSecondary creates a fill with secondary color
func (t *Theme) FillSecondary(widget W) *Fill {
	return t.NewFill(t.Colors.Secondary(), widget)
}

// FillSurface creates a fill with surface color
func (t *Theme) FillSurface(widget W) *Fill {
	return t.NewFill(t.Colors.Surface(), widget)
}

// FillBackground creates a fill with background color
func (t *Theme) FillBackground(widget W) *Fill {
	return t.NewFill(t.Colors.Background(), widget)
}

// FillError creates a fill with error color
func (t *Theme) FillError(widget W) *Fill {
	return t.NewFill(t.Colors.Error(), widget)
}

// FillRounded creates a fill with rounded corners
func (t *Theme) FillRounded(color color.NRGBA, radius float32, widget W) *Fill {
	return t.NewFillWithRadius(color, radius, CornerAll, widget)
}

// FillCard creates a card-like fill with surface color and rounded corners
func (t *Theme) FillCard(widget W) *Fill {
	return t.NewFillWithRadius(t.Colors.Surface(), 8, CornerAll, widget)
}

// FillButton creates a button-like fill with primary color and rounded corners
func (t *Theme) FillButton(widget W) *Fill {
	return t.NewFillWithRadius(t.Colors.Primary(), 4, CornerAll, widget)
}
