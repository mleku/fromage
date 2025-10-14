package fromage

import (
	"gioui.org/app"
)

type Window struct {
	*app.Window
	opts []app.Option
	*Theme
}

func NewWindow(th *Theme) *Window {
	return &Window{Window: &app.Window{}, Theme: th}
}
