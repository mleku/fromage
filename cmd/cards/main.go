package main

import (
	"context"
	"time"

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

// Application state struct to hold persistent widgets
type AppState struct {
	themeToggle *fromage.Bool
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

	// Initialize application state with theme toggle
	appState = &AppState{
		themeToggle: th.Switch(false).SetOnChange(func(b bool) {
			// Toggle theme when switch is clicked
			if b {
				th.SetThemeMode(fromage.ThemeModeDark)
			} else {
				th.SetThemeMode(fromage.ThemeModeLight)
			}
			// Trigger color transition for the switch itself
			appState.themeToggle.UpdateThemeColors(time.Now())
		}),
	}

	w := fromage.NewWindow(th)
	w.Option(app.Size(
		unit.Dp(800), unit.Dp(600)),
		app.Title("Card Widget Demo"),
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
	// Fill the background with the theme's background color
	backgroundFill := th.NewFill(th.Colors.Background(), nil)
	backgroundFill.Layout(gtx)

	// Create a simple card with text content
	cardContent := func(gtx C) D {
		return th.Body1("This is a simple card with some text content.").Layout(gtx)
	}

	// Create cards with different styles
	primaryCard := th.CardPrimary(cardContent)
	surfaceCard := th.CardSurface(cardContent)
	errorCard := th.CardError(cardContent)

	// Create a card with title
	titleCard := th.CardWithTitle("Card with Title", cardContent)

	// Create a card with custom color
	customCard := th.CardWithTitleAndColor("Custom Card", th.Colors.Secondary(), cardContent)

	// Layout all cards in a vertical list
	th.VFlex().
		SpaceStart().
		Rigid(func(gtx C) D {
			return th.H1("Card Widget Demo").Color(th.Colors.OnBackground()).Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			// Theme toggle section
			return th.HFlex().
				SpaceStart().
				Rigid(func(gtx C) D {
					return appState.themeToggle.Layout(gtx)
				}).
				Rigid(func(gtx C) D {
					return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
				}).
				Rigid(func(gtx C) D {
					// Show current theme mode
					themeText := "Light"
					if th.IsDark() {
						themeText = "Dark"
					}
					return th.Body1(themeText).Color(th.Colors.OnBackground()).Layout(gtx)
				}).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.H2("Primary Card").Color(th.Colors.OnBackground()).Layout(gtx)
		}).
		Rigid(primaryCard.Layout).
		Rigid(func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.H3("Surface Card").Color(th.Colors.OnBackground()).Layout(gtx)
		}).
		Rigid(surfaceCard.Layout).
		Rigid(func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.H4("Error Card").Color(th.Colors.OnBackground()).Layout(gtx)
		}).
		Rigid(errorCard.Layout).
		Rigid(func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.H3("Card with Title").Color(th.Colors.OnBackground()).Layout(gtx)
		}).
		Rigid(titleCard.Layout).
		Rigid(func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.H3("Custom Card").Color(th.Colors.OnBackground()).Layout(gtx)
		}).
		Rigid(customCard.Layout).
		Rigid(func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.H3("Card with Inset").Color(th.Colors.OnBackground()).Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			// Create a card with inset padding
			insetCard := th.CardPrimary(func(gtx C) D {
				// Create an inset with 1.0 padding (scaled by text size)
				w := fromage.NewWindow(th)
				inset := w.Inset(1.0, func(gtx C) D {
					return th.Body1("This card content has inset padding around it. The inset creates space between the card border and the content.").Layout(gtx)
				})
				return inset.Fn(gtx)
			})
			return insetCard.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.H3("Card with Different Inset Values").Color(th.Colors.OnBackground()).Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			// Create a card with different inset values
			insetCard := th.CardSurface(func(gtx C) D {
				// Create an inset with 0.5 padding (smaller padding)
				w := fromage.NewWindow(th)
				inset := w.Inset(0.5, func(gtx C) D {
					return th.Body1("This card has smaller inset padding (0.5x text size).").Layout(gtx)
				})
				return inset.Fn(gtx)
			})
			return insetCard.Layout(gtx)
		}).
		Layout(gtx)
}
