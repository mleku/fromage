package main

import (
	"context"
	"fmt"
	"image"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
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

// MouseEvent represents a mouse event
type MouseEvent struct {
	Type      string
	Position  image.Point
	Timestamp string
}

// Application state
type AppState struct {
	lastEvent  MouseEvent
	theme      *fromage.Theme
	pointerTag interface{}
}

var appState *AppState

func main() {
	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeDark,
	)

	w := fromage.NewWindow(th)

	// Initialize application state
	appState = &AppState{
		theme:      th,
		pointerTag: &struct{}{}, // Unique tag for pointer events
		lastEvent: MouseEvent{
			Type:      "No events yet",
			Position:  image.Pt(0, 0),
			Timestamp: "",
		},
	}

	w.Option(
		app.Size(unit.Dp(400), unit.Dp(300)),
		app.Title("Mouse Events Demo"),
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

	// Register for pointer events over the entire window area
	r := image.Rectangle{Max: gtx.Constraints.Max}
	area := clip.Rect(r).Push(gtx.Ops)
	event.Op(gtx.Ops, appState.pointerTag)
	area.Pop()

	// Handle pointer events
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: appState.pointerTag,
			Kinds:  pointer.Press | pointer.Release | pointer.Move,
		})
		if !ok {
			break
		}
		if e, ok := ev.(pointer.Event); ok {
			clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))

			// Determine which button was pressed
			var buttonType string
			switch e.Kind {
			case pointer.Press:
				switch {
				case e.Buttons == pointer.ButtonPrimary:
					buttonType = "Left Click"
				case e.Buttons == pointer.ButtonSecondary:
					buttonType = "Right Click"
				case e.Buttons == pointer.ButtonTertiary:
					buttonType = "Middle Click"
				default:
					buttonType = fmt.Sprintf("Button %d", e.Buttons)
				}
			case pointer.Release:
				switch {
				case e.Buttons == pointer.ButtonPrimary:
					buttonType = "Left Release"
				case e.Buttons == pointer.ButtonSecondary:
					buttonType = "Right Release"
				case e.Buttons == pointer.ButtonTertiary:
					buttonType = "Middle Release"
				default:
					buttonType = fmt.Sprintf("Button %d Release", e.Buttons)
				}
			case pointer.Move:
				buttonType = "Mouse Move"
			}

			appState.lastEvent = MouseEvent{
				Type:      buttonType,
				Position:  clickPos,
				Timestamp: fmt.Sprintf("%v", e.Time),
			}

			log.I.F("Mouse event: %s at (%d, %d)", appState.lastEvent.Type, clickPos.X, clickPos.Y)
		}
	}

	// Layout the UI
	th.CenteredColumn().
		Rigid(func(g C) D {
			// Title
			return th.H3("Mouse Events Demo").
				Color(th.Colors.OnBackground()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Instructions
			return th.Body1("Click anywhere in this window to see mouse events").
				Color(th.Colors.OnSurfaceVariant()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Event display area
			return th.NewCard(
				func(g C) D {
					return th.VFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							return th.Body2("Last Mouse Event:").
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1(fmt.Sprintf("Type: %s", appState.lastEvent.Type)).
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1(fmt.Sprintf("Position: (%d, %d)", appState.lastEvent.Position.X, appState.lastEvent.Position.Y)).
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1(fmt.Sprintf("Time: %s", appState.lastEvent.Timestamp)).
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Layout(g)
				},
			).CornerRadius(8).Padding(unit.Dp(16)).Layout(g)
		}).
		Layout(gtx)
}
