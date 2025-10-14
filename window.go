package fromage

import "gioui.org/app"

type Window struct {
	*app.Window
	th *Theme
}

func NewWindow(th *Theme) *Window {
	return &Window{Window: &app.Window{}, th: th}
}
