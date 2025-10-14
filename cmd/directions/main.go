package main

import (
	"context"
	"image"

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
		unit.Dp(600), unit.Dp(600)),
		app.Title("Direction Layout Showcase"),
	)
	w.Run(loop(w.Window, th, w))
}

func loop(w *app.Window, th *fromage.Theme, window *fromage.Window) func() {
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
				mainUI(gtx, th, window)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, window *fromage.Window) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	th.CenteredColumn().
		Rigid(func(g C) D {
			// Title with primary color fill
			return th.H4("Direction Layout Showcase").Alignment(text.Middle).Layout(g)
		}).
		Rigid(func(g C) D {
			// 3x3 Grid of Direction Examples
			return th.VFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Row 1: NW, N, NE
					return th.HFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "NW", window.Direction().NW())(g)
						}).
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "N", window.Direction().N())(g)
						}).
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "NE", window.Direction().NE())(g)
						}).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Row 2: W, Center, E
					return th.HFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "W", window.Direction().W())(g)
						}).
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "Center", window.Direction().Center())(g)
						}).
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "E", window.Direction().E())(g)
						}).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Row 3: SW, S, SE
					return th.HFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "SW", window.Direction().SW())(g)
						}).
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "S", window.Direction().S())(g)
						}).
						Rigid(func(g C) D {
							return createDirectionBox(th, window, "SE", window.Direction().SE())(g)
						}).
						Layout(g)
				}).
				Layout(g)
		}).
		Layout(gtx)
}

// createDirectionBox creates a bordered box demonstrating a specific direction alignment
func createDirectionBox(th *fromage.Theme, window *fromage.Window, label string, direction *fromage.Direction) W {
	return func(g C) D {
		// Calculate box size: 6 text heights square
		boxSize := unit.Dp(float32(th.TextSize) * 6)

		// Constrain the layout to the desired size
		gtx := g
		gtx.Constraints.Min = image.Pt(g.Dp(boxSize), g.Dp(boxSize))
		gtx.Constraints.Max = image.Pt(g.Dp(boxSize), g.Dp(boxSize))

		return th.NewBorder().
			Color(th.Colors.Primary()).
			Width(unit.Dp(2)).
			CornerRadius(unit.Dp(4)).
			Widget(func(g C) D {
				return direction.Embed(func(g C) D {
					return th.Body1(label).
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).Fn(gtx)
			}).
			Layout(gtx)
	}
}
