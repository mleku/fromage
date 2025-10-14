package fromage

import "gioui.org/layout"

type Direction struct {
	D layout.Direction
	w W
}

// Direction creates a directional layout that sets its contents to align according to the configured direction (8
// cardinal directions and centered)
func (w *Theme) Direction() (out *Direction) {
	out = &Direction{}
	return
}

// direction setters

// NW sets the relevant direction for the Direction layout
func (d *Direction) NW() (out *Direction) {
	d.D = layout.NW
	return d
}

// N sets the relevant direction for the Direction layout
func (d *Direction) N() (out *Direction) {
	d.D = layout.N
	return d
}

// NE sets the relevant direction for the Direction layout
func (d *Direction) NE() (out *Direction) {
	d.D = layout.NE
	return d
}

// E sets the relevant direction for the Direction layout
func (d *Direction) E() (out *Direction) {
	d.D = layout.E
	return d
}

// SE sets the relevant direction for the Direction layout
func (d *Direction) SE() (out *Direction) {
	d.D = layout.SE
	return d
}

// S sets the relevant direction for the Direction layout
func (d *Direction) S() (out *Direction) {
	d.D = layout.S
	return d
}

// SW sets the relevant direction for the Direction layout
func (d *Direction) SW() (out *Direction) {
	d.D = layout.SW
	return d
}

// W sets the relevant direction for the Direction layout
func (d *Direction) W() (out *Direction) {
	d.D = layout.W
	return d
}

// Center sets the relevant direction for the Direction layout
func (d *Direction) Center() (out *Direction) {
	d.D = layout.Center
	return d
}

func (d *Direction) Embed(w layout.Widget) *Direction {
	d.w = w
	return d
}

// Fn the given widget given the context and direction
func (d *Direction) Fn(c layout.Context) layout.Dimensions {
	return d.D.Layout(c, d.w)
}
