package main

import (
	"context"
	"fmt"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
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
		unit.Dp(800), unit.Dp(1600)),
		app.Title("Scrollbar Demo"),
	)
	w.Run(loop(w.Window, th, w))
}

func loop(w *app.Window, th *fromage.Theme, window *fromage.Window) func() {
	// Float slider to control viewport proportion
	viewportSlider := th.NewFloat().SetRange(0.1, 1.0).SetValue(0.5)

	// Horizontal scrollbar
	horizontalScrollbar := th.NewScrollbar(fromage.Horizontal)

	// Vertical scrollbar
	verticalScrollbar := th.NewScrollbar(fromage.Vertical)

	var viewportProportion float32 = 0.5
	var horizontalPos float32 = 0.0
	var verticalPos float32 = 0.0

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

				// Update viewport proportion if slider changed
				if viewportSlider.Changed() {
					viewportProportion = viewportSlider.Value()
					horizontalScrollbar.SetViewport(viewportProportion)
					verticalScrollbar.SetViewport(viewportProportion)
				}

				// Update scrollbar positions if they changed
				if horizontalScrollbar.Changed() {
					horizontalPos = horizontalScrollbar.Position()
				}
				if verticalScrollbar.Changed() {
					verticalPos = verticalScrollbar.Position()
				}

				mainUI(gtx, th, window, viewportSlider, horizontalScrollbar, verticalScrollbar, viewportProportion, horizontalPos, verticalPos)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, window *fromage.Window,
	viewportSlider *fromage.Float, horizontalScrollbar, verticalScrollbar *fromage.Scrollbar,
	viewportProportion, horizontalPos, verticalPos float32) {

	// Fill background with theme background color
	th.FillBackground(nil).Layout(gtx)

	// Split into top and bottom halves
	th.VFlex().
		Flexed(1, func(g C) D {
			// Top half - Horizontal scrollbar demo
			return th.CenteredColumn().
				Rigid(func(g C) D {
					return th.H4("Horizontal Scrollbar Demo").Alignment(text.Middle).Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Body1(fmt.Sprintf("Viewport: %.1f%% | Position: %.1f%%",
						viewportProportion*100, horizontalPos*100)).
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Viewport control slider
					return window.Inset(1.25, func(g C) D {
						return th.CenteredColumn().
							Rigid(func(g C) D {
								return th.Body2("Viewport Proportion:").Layout(g)
							}).
							Rigid(func(g C) D {
								return viewportSlider.Layout(g, th)
							}).
							Layout(g)
					}).Fn(g)
				}).
				Rigid(func(g C) D {
					// Horizontal scrollbar
					return window.Inset(1.25, func(g C) D {
						return horizontalScrollbar.Layout(g, th)
					}).Fn(g)
				}).
				Layout(g)
		}).
		Flexed(1, func(g C) D {
			// Bottom half - Vertical scrollbar demo
			return th.CenteredColumn().
				Rigid(func(g C) D {
					return th.H4("Vertical Scrollbar Demo").Alignment(text.Middle).Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Body1(fmt.Sprintf("Viewport: %.1f%% | Position: %.1f%%",
						viewportProportion*100, verticalPos*100)).
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Flexed(1, func(g C) D {
					// Vertical scrollbar taking up 50% of vertical space
					return window.Inset(1.25, func(g C) D {
						return th.HFlex().
							Flexed(1, func(g C) D {
								// Spacer
								return layout.Dimensions{}
							}).
							Rigid(func(g C) D {
								// Vertical scrollbar on the right
								return verticalScrollbar.Layout(g, th)
							}).
							Layout(g)
					}).Fn(g)
				}).
				Layout(g)
		}).
		Layout(gtx)
}
