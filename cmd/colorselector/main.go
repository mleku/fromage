package main

import (
	"context"
	"fmt"
	"image"
	"image/color"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
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

// Application state struct to hold persistent widgets
type AppState struct {
	colorSelector *fromage.ColorSelector
}

var appState *AppState

func main() {
	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeLight,
	)

	// Initialize application state with persistent widgets
	appState = &AppState{
		colorSelector: th.NewColorSelector().SetOnChange(func(c color.NRGBA) {
			log.I.F("[HOOK] Color changed to: R=%d G=%d B=%d", c.R, c.G, c.B)
			// Update surface tint using HSV values
			th.Colors.SetSurfaceTintFromHSV(
				appState.colorSelector.GetHue(),
				appState.colorSelector.GetSaturation(),
				appState.colorSelector.GetTone(),
			)
		}),
	}

	// Initialize the color selector with the current surface tint
	currentSurfaceTint := th.Colors.GetSurfaceTint()
	appState.colorSelector.SetColor(currentSurfaceTint)

	w := fromage.NewWindow(th)
	w.Option(app.Size(
		unit.Dp(800), unit.Dp(600)),
		app.Title("Color Selector Demo"),
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
				mainUI(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	th.CenteredColumn().
		Rigid(func(g C) D {
			// Title
			return th.H3("Color Selector Demo").Alignment(text.Middle).Layout(g)
		}).
		Rigid(func(g C) D {
			// Theme info with surface color fill
			themeText := "Current Theme: Light"
			if th.IsDark() {
				themeText = "Current Theme: Dark"
			}
			return th.FillSurface(
				func(g C) D {
					return th.Body1(themeText).Alignment(text.Middle).Layout(g)
				},
			).CornerRadius(4).Layout(g)
		}).
		Rigid(func(g C) D {
			// Theme toggle button
			button := th.PrimaryButton(func(g C) D {
				return th.Body1("Toggle Theme").
					Color(th.Colors.OnPrimary()).
					Alignment(text.Middle).
					Layout(g)
			})

			if button.Clicked(g) {
				log.I.F("Toggle theme button clicked")
				th.ToggleTheme()
			}

			return button.Layout(g)
		}).
		Rigid(func(g C) D {
			// Color selector
			return th.VFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					return th.H3("Surface Tint Color").
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Rigid(func(g C) D {
					return appState.colorSelector.Layout(g, th)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Current color display
			currentColor := appState.colorSelector.GetColor()
			tone := appState.colorSelector.GetTone()
			hue := appState.colorSelector.GetHue()
			saturation := appState.colorSelector.GetSaturation()

			return th.FillSurface(
				func(g C) D {
					return th.VFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							return th.Body1("Current Surface Tint").
								Color(th.Colors.OnSurface()).
								Alignment(text.Middle).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Caption(fmt.Sprintf("RGB(%d, %d, %d)", currentColor.R, currentColor.G, currentColor.B)).
								Color(th.Colors.OnSurface()).
								Alignment(text.Middle).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Caption(fmt.Sprintf("HEX %s", fromage.ColorToHex(currentColor))).
								Color(th.Colors.OnSurface()).
								Alignment(text.Middle).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Caption(fmt.Sprintf("HSV(%.4f, %.4f, %.4f)", hue, saturation, tone)).
								Color(th.Colors.OnSurface()).
								Alignment(text.Middle).
								Layout(g)
						}).
						Layout(g)
				},
			).CornerRadius(4).Layout(g)
		}).
		Rigid(func(g C) D {
			// Color preview
			currentColor := appState.colorSelector.GetColor()
			return th.FillSurface(
				func(g C) D {
					return th.VFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							return th.Body1("Color Preview").
								Color(th.Colors.OnSurface()).
								Alignment(text.Middle).
								Layout(g)
						}).
						Rigid(func(g C) D {
							// Draw the actual color as a rectangle
							rect := image.Rectangle{Max: image.Pt(100, 50)}
							defer clip.Rect(rect).Push(g.Ops).Pop()
							paint.Fill(g.Ops, currentColor)
							return layout.Dimensions{Size: rect.Max}
						}).
						Layout(g)
				},
			).CornerRadius(4).Layout(g)
		}).
		Rigid(func(g C) D {
			// Sample cards to show surface tint effect
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Primary card
					card := th.CardPrimary(func(g C) D {
						return th.Body1("Primary Card").
							Color(th.Colors.OnPrimary()).
							Alignment(text.Middle).
							Layout(g)
					})
					return card.Layout(g)
				}).
				Rigid(func(g C) D {
					// Surface card
					card := th.CardSurface(func(g C) D {
						return th.Body1("Surface Card").
							Color(th.Colors.OnSurface()).
							Alignment(text.Middle).
							Layout(g)
					})
					return card.Layout(g)
				}).
				Rigid(func(g C) D {
					// Secondary card
					card := th.CardSecondary(func(g C) D {
						return th.Body1("Secondary Card").
							Color(th.Colors.OnSecondary()).
							Alignment(text.Middle).
							Layout(g)
					})
					return card.Layout(g)
				}).
				Layout(g)
		}).
		Layout(gtx)
}
