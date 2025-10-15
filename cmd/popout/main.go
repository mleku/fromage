package main

import (
	"context"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
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
		fromage.ThemeModeLight,
	)

	w := fromage.NewWindow(th)

	w.Option(
		app.Size(unit.Dp(800), unit.Dp(600)),
		app.Title("Modal Positioning Demo"),
	)
	w.Run(loop(w.Window, th))
}

func loop(w *app.Window, th *fromage.Theme) func() {
	return func() {
		var ops op.Ops
		// Create a fromage window wrapper
		fromageWindow := &fromage.Window{Window: w, Theme: th}
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				chk.E(e.Err)
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				th.Pool.Reset() // Reset pool at the beginning of each frame
				mainUI(gtx, th, fromageWindow)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, w *fromage.Window) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	// Draw center border lines (red cross)
	drawCenterLines(gtx, th)

	// Create flex-column with two flex-row containers
	th.VFlex().
		Flexed(1, func(g C) D {
			// Top row container
			return th.HFlex().
				Flexed(1, func(g C) D {
					return th.Direction().Center().Embed(func(g C) D {
						return th.Body1("NW").
							Color(th.Colors.OnBackground()).
							Alignment(text.Middle).
							Layout(g)
					}).Fn(g)
				}).
				Flexed(1, func(g C) D {
					return th.Direction().Center().Embed(func(g C) D {
						return th.Body1("NE").
							Color(th.Colors.OnBackground()).
							Alignment(text.Middle).
							Layout(g)
					}).Fn(g)
				}).
				Layout(g)
		}).
		Flexed(1, func(g C) D {
			// Bottom row container
			return th.HFlex().
				Flexed(1, func(g C) D {
					return th.Direction().Center().Embed(func(g C) D {
						// SW (South-West) label
						return th.Body1("SW").
							Color(th.Colors.OnBackground()).
							Alignment(text.Middle).
							Layout(g)
					}).Fn(g)
				}).
				Flexed(1, func(g C) D {
					return th.Direction().Center().Embed(func(g C) D {
						return th.Body1("SE").
							Color(th.Colors.OnBackground()).
							Alignment(text.Middle).
							Layout(g)
					}).Fn(g)
				}).
				Layout(g)
		}).
		Layout(gtx)
}

// drawCenterLines draws red cross lines in the center of the view
func drawCenterLines(gtx layout.Context, th *fromage.Theme) {
	screenWidth := gtx.Constraints.Max.X
	screenHeight := gtx.Constraints.Max.Y

	// Red color for the lines
	redColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	// Draw horizontal line (top to bottom)
	horizontalRect := image.Rect(screenWidth/2-1, 0, screenWidth/2+1, screenHeight)
	paint.FillShape(gtx.Ops, redColor, clip.Rect(horizontalRect).Op())

	// Draw vertical line (left to right)
	verticalRect := image.Rect(0, screenHeight/2-1, screenWidth, screenHeight/2+1)
	paint.FillShape(gtx.Ops, redColor, clip.Rect(verticalRect).Op())
}
