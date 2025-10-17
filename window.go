package fromage

import (
	"gio.mleku.dev/app"
)

type Window struct {
	*app.Window
	opts []app.Option
	*Theme
}

func NewWindow(th *Theme) *Window {
	return &Window{Window: &app.Window{}, Theme: th}
}

// Invalidate requests a new frame to be drawn
func (w *Window) Invalidate() {
	if w.Window != nil {
		w.Window.Invalidate()
	}
}

// Run starts the window event loop with the provided function
func (w *Window) Run(fn func()) {
	if w.Window != nil {
		w.Window.Run(fn)
	}
}
