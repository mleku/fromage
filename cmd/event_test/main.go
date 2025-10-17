package main

import (
	"context"
	"fmt"
	"image"

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

// TestState tracks events for demonstration
type TestState struct {
	buttonClicks     int
	scrollbarScrolls int
	scrollEvents     int
	clickEvents      int
	theme            *fromage.Theme
	eventHandler     *fromage.EventHandler
}

var testState *TestState

func main() {
	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeDark,
	)

	w := fromage.NewWindow(th)

	// Initialize test state
	testState = &TestState{
		theme: th,
	}

	// Create event handler with callbacks
	testState.eventHandler = fromage.NewEventHandler(func(event string) {
		log.I.F("[Test] %s", event)
	}).SetOnClick(func(e pointer.Event) {
		testState.clickEvents++
	}).SetOnScroll(func(distance float32) {
		testState.scrollEvents++
	})

	w.Option(
		app.Size(unit.Dp(800), unit.Dp(600)),
		app.Title("EventHandler Test - Simultaneous Events"),
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
	testState.eventHandler.AddToOps(gtx.Ops)
	testState.eventHandler.ProcessEvents(gtx)

	// Create a button
	button := th.NewButtonLayout().
		Background(th.Colors.Primary()).
		CornerRadius(0.5).
		Widget(func(g C) D {
			return th.Body1("Click Me!").
				Color(th.Colors.OnPrimary()).
				Alignment(text.Middle).
				Layout(g)
		}).
		OnClick(func() {
			testState.buttonClicks++
		})

	// Create a scrollbar
	scrollbar := th.NewScrollbar(fromage.Vertical).
		SetViewport(0.3).
		SetPosition(0.5).
		SetHook(func(pos float32) {
			testState.scrollbarScrolls++
		})

	// Layout the UI
	th.CenteredColumn().
		Rigid(func(gtx C) D {
			return th.H1("EventHandler Test").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body1("This test demonstrates that widgets can receive scroll and click events simultaneously.").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2("Try clicking the button and scrolling at the same time!").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			// Button with EventHandler
			return button.Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			// Scrollbar with EventHandler
			gtx.Constraints.Min = image.Pt(200, 300)
			return scrollbar.Layout(gtx, th)
		}).
		Rigid(func(gtx C) D {
			return th.Body2("Event Counters:").
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Button Clicks: %d", testState.buttonClicks)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Scrollbar Scrolls: %d", testState.scrollbarScrolls)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Global Click Events: %d", testState.clickEvents)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Rigid(func(gtx C) D {
			return th.Body2(fmt.Sprintf("Global Scroll Events: %d", testState.scrollEvents)).
				Color(th.Colors.OnBackground()).
				Layout(gtx)
		}).
		Layout(gtx)
}
