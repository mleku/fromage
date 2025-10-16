package main

import (
	"context"
	"fmt"
	"image"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
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

// Application state
type AppState struct {
	theme      *fromage.Theme
	globalMenu *fromage.GlobalMenu
	widgets    []*SimpleWidget
	nextID     int
}

// SimpleWidget represents a simple widget that can be placed
type SimpleWidget struct {
	id       int
	position image.Point
	theme    *fromage.Theme
	visible  bool
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
		theme: th,
		globalMenu: th.NewGlobalMenu().
			AddItem("Create Widget", func() {
				// Create widget at the right-click position
				clickPos := appState.globalMenu.GetClickPosition()
				appState.createWidgetAt(clickPos, image.Pt(800, 600))
				fmt.Printf("Created widget at (%d, %d)\n", clickPos.X, clickPos.Y)
			}).
			AddItem("Clear All Widgets", func() {
				appState.widgets = make([]*SimpleWidget, 0)
				fmt.Println("Cleared all widgets")
			}).
			AddItem("Widget Count", func() {
				fmt.Printf("Current widget count: %d\n", len(appState.widgets))
			}),
		widgets: make([]*SimpleWidget, 0),
		nextID:  1,
	}

	w.Option(
		app.Size(unit.Dp(800), unit.Dp(600)),
		app.Title("Right-Click Widget Demo"),
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
	th.FillBackground(nil).Layout(gtx)

	// Layout the UI
	th.CenteredColumn().
		Rigid(func(g C) D {
			// Title
			return th.H3("Right-Click Widget Demo").
				Color(th.Colors.OnBackground()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Instructions
			return th.Body1("Right-click anywhere to show context menu. Use the menu to create widgets.").
				Color(th.Colors.OnSurfaceVariant()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Widget count display
			return th.Body2(fmt.Sprintf("Widgets placed: %d", len(appState.widgets))).
				Color(th.Colors.OnSurface()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Layout(gtx)

	// Handle global menu events
	appState.globalMenu.HandleEvents(gtx)

	// Layout the global menu
	if appState.globalMenu.IsVisible() {
		appState.globalMenu.Layout(gtx)
	}

	// Layout all widgets
	for _, widget := range appState.widgets {
		if widget.visible {
			widget.Layout(gtx)
		}
	}
}

// createWidgetAt creates a new widget at the specified position
func (app *AppState) createWidgetAt(position image.Point, viewportSize image.Point) {
	widget := &SimpleWidget{
		position: app.calculateSmartPosition(position, viewportSize),
		theme:    app.theme,
		visible:  true,
		id:       app.nextID,
	}

	app.widgets = append(app.widgets, widget)
	app.nextID++
}

// calculateSmartPosition calculates where to position the widget so the corner faces toward the center
func (app *AppState) calculateSmartPosition(clickPos image.Point, viewportSize image.Point) image.Point {
	centerX := viewportSize.X / 2
	centerY := viewportSize.Y / 2

	widgetWidth := 150 // Widget width
	widgetHeight := 80 // Widget height

	// Determine which corner should face toward the center
	if clickPos.X < centerX {
		// Click is on left side of center
		if clickPos.Y < centerY {
			// Click is in top-left quadrant, position widget bottom-right of click
			return image.Pt(clickPos.X, clickPos.Y)
		} else {
			// Click is in bottom-left quadrant, position widget top-right of click
			return image.Pt(clickPos.X, clickPos.Y-widgetHeight)
		}
	} else {
		// Click is on right side of center
		if clickPos.Y < centerY {
			// Click is in top-right quadrant, position widget bottom-left of click
			return image.Pt(clickPos.X-widgetWidth, clickPos.Y)
		} else {
			// Click is in bottom-right quadrant, position widget top-left of click
			return image.Pt(clickPos.X-widgetWidth, clickPos.Y-widgetHeight)
		}
	}
}

// Layout renders the simple widget at its position
func (w *SimpleWidget) Layout(gtx C) D {
	// Position the widget
	offset := op.Offset(w.position).Push(gtx.Ops)
	defer offset.Pop()

	// Constrain widget size
	gtx.Constraints.Min.X = 150
	gtx.Constraints.Max.X = 150
	gtx.Constraints.Min.Y = 80
	gtx.Constraints.Max.Y = 80

	// Create widget background
	return w.theme.NewCard(
		func(g C) D {
			return w.theme.VFlex().
				Rigid(func(gtx C) D {
					// Header
					return w.theme.Caption(fmt.Sprintf("Widget #%d", w.id)).
						Color(w.theme.Colors.OnSurface()).
						Alignment(text.Start).
						Layout(gtx)
				}).
				Rigid(func(gtx C) D {
					// Content
					return w.theme.Caption("Created via context menu").
						Color(w.theme.Colors.OnSurfaceVariant()).
						Alignment(text.Middle).
						Layout(gtx)
				}).
				Layout(g)
		},
	).CornerRadius(8).Padding(unit.Dp(8)).Layout(gtx)
}
