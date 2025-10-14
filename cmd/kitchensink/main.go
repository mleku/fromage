package main

import (
	"context"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/mleku/fromage"
	"lol.mleku.dev/chk"
)

func main() {
	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeDark,
	)
	w := fromage.NewWindow(th)
	w.Option(app.Size(
		unit.Dp(800), unit.Dp(450)),
		app.Title("Kitchensink - Theme Demo"),
	)
	w.Run(loop(w.Window, th))
}

func loop(w *app.Window, th *fromage.Theme) func() {
	return func() {
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				chk.E(e.Err)
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				mainUI(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())
}
