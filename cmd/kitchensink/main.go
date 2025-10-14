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
	"lol.mleku.dev/log"
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
				th.Pool.Reset() // Reset pool at the beginning of each frame
				mainUI(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	// Use the new Flex API with enhanced button demonstrations
	th.CenteredColumn().
		Rigid(func(g C) D {
			// Title with primary color fill
			return th.FillPrimary(
				func(g C) D {
					return th.H1("Interactive Button Demo").Alignment(text.Middle).Layout(g)
				},
			).CornerRadius(8).Layout(g)
		}).
		Rigid(func(g C) D {
			themeText := "Current Theme: Light"
			if th.IsDark() {
				themeText = "Current Theme: Dark"
			}
			// Theme info with surface color fill
			return th.FillSurface(
				func(g C) D {
					return th.Body1(themeText).Alignment(text.Middle).Layout(g)
				},
			).CornerRadius(4).Layout(g)
		}).
		Rigid(func(g C) D {
			// Main interactive button with theme toggle
			button := th.PrimaryButton(
				func(g C) D {
					return th.Body1("Toggle Theme").
						Color(th.Colors.OnPrimary()).
						Alignment(text.Middle).
						Layout(g)
				},
			)

			// Check for clicks BEFORE layout (this is the key fix!)
			if button.Clicked(g) {
				log.I.F("Toggle theme button clicked")
				th.ToggleTheme()
			}

			return button.Layout(g)
		}).
		Rigid(func(g C) D {
			// Button style showcase with hover effects
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Secondary button
					btn := th.SecondaryButton(
						func(g C) D {
							return th.Body2("Secondary").
								Color(th.Colors.OnSecondary()).
								Alignment(text.Middle).
								Layout(g)
						},
					)
					if btn.Clicked(g) {
						log.I.F("Secondary button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Surface button
					btn := th.SurfaceButton(
						func(g C) D {
							return th.Body2("Surface").
								Color(th.Colors.OnSurface()).
								Alignment(text.Middle).
								Layout(g)
						},
					)
					if btn.Clicked(g) {
						log.I.F("Surface button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Error button
					btn := th.ErrorButton(
						func(g C) D {
							return th.Body2("Error").
								Color(th.Colors.OnError()).
								Alignment(text.Middle).
								Layout(g)
						},
					)
					if btn.Clicked(g) {
						log.I.F("Error button clicked")
					}
					return btn.Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Shape showcase with different corner styles
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Rounded button
					btn := th.RoundedButton(
						func(g C) D {
							return th.Caption("Rounded").
								Color(th.Colors.OnPrimary()).
								Alignment(text.Middle).
								Layout(g)
						},
					)
					if btn.Clicked(g) {
						log.I.F("rounded button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Pill button
					btn := th.PillButton(
						func(g C) D {
							return th.Caption("Pill Shape").
								Color(th.Colors.OnPrimary()).
								Alignment(text.Middle).
								Layout(g)
						},
					)
					if btn.Clicked(g) {
						log.I.F("pill button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Icon button (using text as placeholder)
					btn := th.IconButton("★").
						Background(th.Colors.Tertiary())
					if btn.Clicked(g) {
						log.I.F("icon button clicked")
					}
					return btn.Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Text buttons showcase
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Text button
					btn := th.TextButton("Text Button")
					if btn.Clicked(g) {
						log.I.F("text button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Custom styled button
					btn := th.NewButtonLayout().
						Background(th.Colors.Tertiary()).
						CornerRadius(0.3). // 30% of text size
						Corners(fromage.CornerNW | fromage.CornerNE).
						Widget(func(g C) D {
							return th.Caption("Custom Style").
								Color(th.Colors.OnTertiary()).
								Alignment(text.Middle).
								Layout(g)
						})
					if btn.Clicked(g) {
						log.I.F("custom button clicked")
					}
					return btn.Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Disabled button example
			btn := th.PrimaryButton(
				func(g C) D {
					return th.Body2("Disabled Button").
						Color(th.Colors.OnPrimary()).
						Alignment(text.Middle).
						Layout(g)
				},
			).Disabled(true) // This button is disabled

			return btn.Layout(g)
		}).
		Layout(gtx)
}
