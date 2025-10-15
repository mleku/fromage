package main

import (
	"context"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/mleku/fromage"
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
	w.Option(app.Title("Drawer Demo"))
	w.Option(app.Size(unit.Dp(800), unit.Dp(600)))
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
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				th.Pool.Reset() // Reset pool at the beginning of each frame
				drawerDemo(gtx, th, fromageWindow)
				e.Frame(gtx.Ops)
			}
		}
	}
}

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
)

// Application state struct to hold persistent widgets
type AppState struct {
	drawerWithControls *fromage.DrawerWithControls
}

var appState *AppState

func drawerDemo(gtx C, th *fromage.Theme, win *fromage.Window) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	// Initialize application state with persistent widgets
	if appState == nil {
		// Create radio group for position selection
		positionRadio := win.NewRadioButtonGroup().
			SetLayout(fromage.LayoutVertical).
			AddButton("Left", true).
			AddButton("Right", false).
			AddButton("Top", false).
			AddButton("Bottom", false)

		appState = &AppState{
			drawerWithControls: win.NewDrawerWithControls().
				Width(unit.Dp(300)).
				Height(unit.Dp(200)).
				Content(func(gtx C) D {
					return th.VFlex().
						SpaceEvenly().
						Rigid(func(gtx C) D {
							return th.H6("Drawer Position").
								Color(th.Colors.OnSurface()).
								Layout(gtx)
						}).
						Rigid(func(gtx C) D {
							return th.Body1("Select drawer position:").
								Color(th.Colors.OnSurface()).
								Layout(gtx)
						}).
						Rigid(func(gtx C) D {
							return positionRadio.Layout(gtx)
						}).
						Rigid(func(gtx C) D {
							btn := th.TextButton("Close Drawer")
							if btn.Clicked(gtx) {
								appState.drawerWithControls.Hide()
							}
							return btn.Layout(gtx)
						}).
						Layout(gtx)
				}).
				OnPositionChange(func(pos fromage.DrawerPosition) {
					// Handle position changes
					switch pos {
					case fromage.DrawerLeft:
						log.Println("Drawer moved to left")
					case fromage.DrawerRight:
						log.Println("Drawer moved to right")
					case fromage.DrawerTop:
						log.Println("Drawer moved to top")
					case fromage.DrawerBottom:
						log.Println("Drawer moved to bottom")
					}
				}),
		}

		// Set up radio group change handler
		positionRadio.SetOnChange(func(index int, label string) {
			var newPos fromage.DrawerPosition
			switch label {
			case "Left":
				newPos = fromage.DrawerLeft
			case "Right":
				newPos = fromage.DrawerRight
			case "Top":
				newPos = fromage.DrawerTop
			case "Bottom":
				newPos = fromage.DrawerBottom
			}

			// Hide current drawer and show in new position
			appState.drawerWithControls.Hide()
			// Use a goroutine to delay the position change
			go func() {
				time.Sleep(350 * time.Millisecond) // Wait for hide animation
				appState.drawerWithControls.SetPosition(newPos)
				appState.drawerWithControls.Show()
			}()
		})
	}

	// Main content
	mainContent := func(gtx C) D {
		return th.VFlex().
			SpaceEvenly().
			Rigid(func(gtx C) D {
				return th.H4("Drawer Demo").
					Color(th.Colors.OnSurface()).
					Layout(gtx)
			}).
			Rigid(func(gtx C) D {
				return th.Body1("Click the button to open the drawer. Use the radio buttons inside the drawer to change its position.").
					Color(th.Colors.OnSurface()).
					Layout(gtx)
			}).
			Rigid(func(gtx C) D {
				// Single button to open drawer
				btn := th.TextButton("Open Drawer")
				if btn.Clicked(gtx) {
					appState.drawerWithControls.Show()
				}
				return btn.Layout(gtx)
			}).
			Rigid(func(gtx C) D {
				return th.Body2("Click outside the drawer or use the close button to hide it.").
					Color(th.Colors.OnSurfaceVariant()).
					Layout(gtx)
			}).
			Layout(gtx)
	}

	// Layout the main content with inset
	win.Inset(20, mainContent).Fn(gtx)

	// Layout the drawer (it handles its own visibility and positioning)
	appState.drawerWithControls.Layout(gtx)
}
