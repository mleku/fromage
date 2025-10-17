package fromage

import (
	"gio.mleku.dev/layout"
	"gio.mleku.dev/unit"
)

// Inset creates a padded empty space around a widget
type Inset struct {
	*Window
	in layout.Inset
	w  layout.Widget
}

// Inset creates a padded empty space around a widget
func (w *Window) Inset(pad float32, embed layout.Widget) (out *Inset) {
	out = &Inset{
		Window: w,
		in:     layout.UniformInset(unit.Dp(float32(w.TextSize) * pad)),
		w:      embed,
	}
	return
}

// Embed sets the widget that will be inside the inset
func (in *Inset) Embed(w layout.Widget) *Inset {
	in.w = w
	return in
}

// Fn lays out the given widget with the configured context and padding
func (in *Inset) Fn(gtx layout.Context) layout.Dimensions {
	return in.in.Layout(gtx, in.w)
}
