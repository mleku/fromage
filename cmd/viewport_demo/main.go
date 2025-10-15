package main

import (
	"context"
	"image"

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
		unit.Dp(800), unit.Dp(800)),
		app.Title("Viewport Demo"),
	)
	w.Run(loop(w.Window, th, w))
}

func loop(w *app.Window, th *fromage.Theme, window *fromage.Window) func() {
	// Horizontal scrollbar for bottom edge
	horizontalScrollbar := th.NewScrollbar(fromage.Horizontal)

	// Vertical scrollbar for right edge
	verticalScrollbar := th.NewScrollbar(fromage.Vertical)

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

				// Update scrollbar positions if they changed
				if horizontalScrollbar.Changed() {
					horizontalPos = horizontalScrollbar.Position()
				}
				if verticalScrollbar.Changed() {
					verticalPos = verticalScrollbar.Position()
				}

				mainUI(gtx, th, window, horizontalScrollbar, verticalScrollbar, horizontalPos, verticalPos)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, window *fromage.Window,
	horizontalScrollbar, verticalScrollbar *fromage.Scrollbar,
	horizontalPos, verticalPos float32) {

	// Fill background with theme background color
	th.FillBackground(nil).Layout(gtx)

	// Main layout with scrollbars on edges
	th.VFlex().
		Flexed(1, func(g C) D {
			// Main content area with vertical scrollbar on right
			return th.HFlex().
				Flexed(1, func(g C) D {
					// Main content area - for now just a placeholder
					return th.CenteredColumn().
						Rigid(func(g C) D {
							return th.H4("Viewport Demo").Alignment(text.Middle).Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1("This is the main content area").Layout(g)
						}).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Vertical scrollbar on right edge
					return verticalScrollbar.Layout(g, th)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Bottom area with horizontal scrollbar and corner space
			return th.HFlex().
				Flexed(1, func(g C) D {
					// Horizontal scrollbar on bottom edge (stops before corner)
					return horizontalScrollbar.Layout(g, th)
				}).
				Rigid(func(g C) D {
					// Square corner space (same size as scrollbar width)
					scrollbarWidth := th.TextSize
					return layout.Dimensions{
						Size: image.Pt(gtx.Dp(scrollbarWidth), gtx.Dp(scrollbarWidth)),
					}
				}).
				Layout(g)
		}).
		Layout(gtx)
}
