package main

import (
	"context"
	"fmt"
	"image"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/mleku/fromage"
	"lol.mleku.dev/log"
)

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
)

// MouseEvent represents a mouse event with coordinates
type MouseEvent struct {
	Type      string
	Position  image.Point
	Timestamp string
}

// GestureEvent represents a gesture event
type GestureEvent struct {
	Type      string
	Position  image.Point
	Details   string
	Timestamp string
}

// Application state
type AppState struct {
	lastMouseEvent   MouseEvent
	lastGestureEvent GestureEvent
	theme            *fromage.Theme
	eventHandler     *fromage.EventHandler
	lastPressed      pointer.Buttons // Track the last pressed button
}

var appState *AppState

func main() {
	// Debug: Print button constants
	log.I.F("Button constants - Primary: %d, Secondary: %d, Tertiary: %d",
		pointer.ButtonPrimary, pointer.ButtonSecondary, pointer.ButtonTertiary)

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
		theme:       th,
		lastPressed: 0, // No buttons pressed initially
		lastMouseEvent: MouseEvent{
			Type:      "No mouse events yet",
			Position:  image.Pt(0, 0),
			Timestamp: "",
		},
		lastGestureEvent: GestureEvent{
			Type:      "No gesture events yet",
			Position:  image.Pt(0, 0),
			Details:   "",
			Timestamp: "",
		},
	}

	// Create event handler with callbacks
	appState.eventHandler = fromage.NewEventHandler(func(event string) {
		log.I.F("[App] %s", event)
	}).SetOnClick(func(e pointer.Event) {
		clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))
		appState.lastMouseEvent = MouseEvent{
			Type:      "Click",
			Position:  clickPos,
			Timestamp: fmt.Sprintf("%v", e.Time),
		}
	}).SetOnPress(func(e pointer.Event) {
		clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))
		appState.lastPressed = e.Buttons
		var buttonType string
		switch e.Buttons {
		case pointer.ButtonPrimary:
			buttonType = "Left Press"
		case pointer.ButtonSecondary:
			buttonType = "Right Press"
		case pointer.ButtonTertiary:
			buttonType = "Middle Press"
		default:
			buttonType = fmt.Sprintf("Button %d Press", e.Buttons)
		}
		appState.lastMouseEvent = MouseEvent{
			Type:      buttonType,
			Position:  clickPos,
			Timestamp: fmt.Sprintf("%v", e.Time),
		}
	}).SetOnRelease(func(e pointer.Event) {
		clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))
		var buttonType string
		switch appState.lastPressed {
		case pointer.ButtonPrimary:
			buttonType = "Left Release"
		case pointer.ButtonSecondary:
			buttonType = "Right Release"
		case pointer.ButtonTertiary:
			buttonType = "Middle Release"
		default:
			buttonType = fmt.Sprintf("Button %d Release", appState.lastPressed)
		}
		appState.lastMouseEvent = MouseEvent{
			Type:      buttonType,
			Position:  clickPos,
			Timestamp: fmt.Sprintf("%v", e.Time),
		}
		appState.lastPressed = 0
	}).SetOnScroll(func(distance float32) {
		appState.lastGestureEvent = GestureEvent{
			Type:      fmt.Sprintf("Scroll %.2f", distance),
			Position:  image.Pt(0, 0), // Scroll doesn't have a specific position
			Details:   fmt.Sprintf("Distance: %.2f", distance),
			Timestamp: fmt.Sprintf("%v", time.Now()),
		}
	})

	w.Option(
		app.Size(unit.Dp(800), unit.Dp(800)), // 800x800 Dp window as requested
		app.Title("Event Handler Demo - Mouse Coordinate Tracker"),
	)

	w.Run(loop(w.Window, th, w))
}

func loop(w *app.Window, th *fromage.Theme, window *fromage.Window) func() {
	return func() {
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				mainUI(gtx, th, window)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, w *fromage.Window) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	// Register event handler for the entire window area
	appState.eventHandler.AddToOps(gtx.Ops)
	appState.eventHandler.ProcessEvents(gtx)

	// Layout the UI
	th.CenteredColumn().
		Rigid(func(gtx C) D {
			return th.H1("Event Handler Demo").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body1("This demo shows how all widgets can receive scroll and click events simultaneously using the EventHandler.").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2("Mouse Events:").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Type: %s", appState.lastMouseEvent.Type)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Position: (%d, %d)", appState.lastMouseEvent.Position.X, appState.lastMouseEvent.Position.Y)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Timestamp: %s", appState.lastMouseEvent.Timestamp)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2("Scroll Events:").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Type: %s", appState.lastGestureEvent.Type)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Details: %s", appState.lastGestureEvent.Details)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Timestamp: %s", appState.lastGestureEvent.Timestamp)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2("Try clicking and scrolling anywhere on this window!").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Layout(gtx)
}
