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

	// Use the new Flex API with fluent method chaining
	th.CenteredColumn().
		Rigid(func(g C) D {
			return th.H1("Theme Demo").Alignment(text.Middle).Layout(g)
		}).
		Rigid(func(g C) D {
			themeText := "Current Theme: Light"
			if th.IsDark() {
				themeText = "Current Theme: Dark"
			}
			return th.Body1(themeText).Alignment(text.Middle).Layout(g)
		}).
		Rigid(func(g C) D {
			// Nested flex for text samples
			return th.VFlex().
				SpaceEnd().
				Rigid(func(g C) D {
					return th.H2("Heading 2").Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Body1("This is body text with fluent API").Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Body2("This is smaller body text").Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Caption("This is caption text").Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			return th.NewLabel().
				Text("Custom styled text").
				Color(th.Colors.Primary()).
				TextSize(unit.Sp(18)).
				Alignment(text.Middle).
				Layout(g)
		}).
		Layout(gtx)
}
