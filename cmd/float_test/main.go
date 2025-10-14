package main

import (
	"context"
	"fmt"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/mleku/fromage"
	"lol.mleku.dev/chk"
)

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
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
		unit.Dp(400), unit.Dp(300)),
		app.Title("Float Slider Test"),
	)
	w.Run(loop(w.Window, th, w))
}

func loop(w *app.Window, th *fromage.Theme, window *fromage.Window) func() {
	slider := th.NewFloat().SetRange(0, 100).SetValue(50)
	var currentValue float32 = 50

	return func() {
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				chk.E(e.Err)
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				th.Pool.Reset() // Reset pool at the beginning of each frame

				// Update current value if slider changed
				if slider.Changed() {
					currentValue = slider.Value()
				}

				mainUI(gtx, th, window, slider, currentValue)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, window *fromage.Window, slider *fromage.Float, value float32) {
	// Fill background with theme background color
	th.FillBackground(nil).Layout(gtx)

	th.CenteredColumn().
		Rigid(func(g C) D {
			// Title
			return th.H4("Float Slider Test").Alignment(text.Middle).Layout(g)
		}).
		Rigid(func(g C) D {
			// Current value display
			return th.Body1(fmt.Sprintf("Current Value: %.1f", value)).
				Color(th.Colors.OnBackground()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Slider with padding
			return window.Inset(1.25, func(g C) D {
				return slider.Layout(g, th)
			}).Fn(g)
		}).
		Layout(gtx)
}
