package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"time"

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
	"gioui.org/widget"
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

// ClickableWidget represents a widget that can be placed at a click position
type ClickableWidget struct {
	position     image.Point
	theme        *fromage.Theme
	closeButton  *fromage.ButtonLayout
	clickable    *widget.Clickable
	scrimClick   *widget.Clickable
	visible      bool
	scrimVisible bool
	id           int
	// Animation state
	showTime     time.Time
	hideTime     time.Time
	isHiding     bool
	shouldRemove bool
}

// Application state
type AppState struct {
	theme        *fromage.Theme
	pointerTag   interface{}
	scrimClick   *widget.Clickable
	widgets      []*ClickableWidget
	nextWidgetID int
	animating    bool
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
		theme:        th,
		pointerTag:   &struct{}{}, // Unique tag for pointer events
		scrimClick:   &widget.Clickable{},
		widgets:      make([]*ClickableWidget, 0),
		nextWidgetID: 1,
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
	paint.Fill(gtx.Ops, th.Colors.Background())

	// Register for pointer events over the entire window area
	r := image.Rectangle{Max: gtx.Constraints.Max}
	area := clip.Rect(r).Push(gtx.Ops)
	event.Op(gtx.Ops, appState.pointerTag)
	area.Pop()

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
			return th.Body1("Right-click anywhere to place a widget. Click the × button or click the scrim to remove widgets.").
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

	// Check if any widget has a visible scrim and handle scrim events
	hasVisibleScrim := false
	for _, widget := range appState.widgets {
		if widget.visible && widget.scrimVisible && !widget.shouldRemove && !widget.isHiding {
			hasVisibleScrim = true
			break
		}
	}

	// If scrim is visible, create a scrim layer that consumes all events BEFORE widget layout
	if hasVisibleScrim {
		// Create scrim area that covers the entire window and consumes all events
		scrimArea := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)
		event.Op(gtx.Ops, appState.scrimClick)

		// Consume all pointer events on the scrim
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: appState.scrimClick,
				Kinds:  pointer.Press,
			})
			if !ok {
				break
			}
			if e, ok := ev.(pointer.Event); ok {
				clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))

				if e.Kind == pointer.Press {
					widgetClosed := false

					// Check which widget's scrim was clicked
					for _, widget := range appState.widgets {
						if widget.visible && widget.scrimVisible && !widget.shouldRemove && !widget.isHiding {
							widgetRect := image.Rectangle{
								Min: widget.position,
								Max: widget.position.Add(image.Pt(150, 80)), // Widget size
							}

							// If clicking outside the widget, start hide animation
							if !clickPos.In(widgetRect) && e.Buttons == pointer.ButtonPrimary {
								// Start hide animation instead of immediate removal
								widget.startHideAnimation()
								widget.scrimVisible = false // Hide scrim immediately
								log.I.F("Widget %d hide animation started", widget.id)

								// Update scrim visibility immediately after hiding scrim
								hasVisibleScrim = false
								for _, w := range appState.widgets {
									if w.visible && w.scrimVisible && !w.shouldRemove && !w.isHiding {
										hasVisibleScrim = true
										break
									}
								}

								widgetClosed = true

								// IMMEDIATELY break out of scrim event loop to prevent catching next click
								log.I.F("Breaking scrim event loop immediately")
								break // Exit the widget loop immediately
							}
						}
					}

					// If we closed a widget, break out of the event loop immediately
					if widgetClosed {
						log.I.F("Widget closed, breaking scrim event loop immediately")
						break // Break out of the outer event loop
					}

					// Re-evaluate scrim visibility after potential widget removal
					hasVisibleScrim = false
					for _, widget := range appState.widgets {
						if widget.visible && widget.scrimVisible && !widget.shouldRemove && !widget.isHiding {
							hasVisibleScrim = true
							break
						}
					}

					// If no scrim is visible anymore, break out of scrim event handling
					if !hasVisibleScrim {
						log.I.F("No scrim visible, breaking scrim event loop")
						break
					}

					// Consume all other events - no events pass through scrim
					log.I.F("Event consumed by scrim at (%d, %d)", clickPos.X, clickPos.Y)
				}
			}
		}

		scrimArea.Pop()

		// After scrim event handling, check if we still have visible scrims
		// If not, we need to process right-click events in this same frame
		hasVisibleScrim = false
		for _, widget := range appState.widgets {
			if widget.visible && widget.scrimVisible && !widget.shouldRemove && !widget.isHiding {
				hasVisibleScrim = true
				break
			}
		}
	}

	// Check if any widgets are animating and clean up completed animations
	appState.animating = false

	// First pass: check animations and mark widgets for removal
	for _, widget := range appState.widgets {
		// Check if widget is in animation (regardless of visibility)
		now := time.Now()
		if !widget.isHiding {
			// Fade in animation
			elapsed := now.Sub(widget.showTime)
			if elapsed < 250*time.Millisecond {
				appState.animating = true
			}
		} else {
			// Fade out animation
			elapsed := now.Sub(widget.hideTime)
			if elapsed < 250*time.Millisecond {
				appState.animating = true
			} else {
				// Fade out animation has completed, mark widget for removal
				if !widget.shouldRemove {
					widget.shouldRemove = true
					log.I.F("Fade-out animation completed for widget %d, marking for removal", widget.id)
				}
			}
		}
	}

	// Second pass: remove widgets that are marked for removal
	for i := len(appState.widgets) - 1; i >= 0; i-- {
		widget := appState.widgets[i]
		if widget.shouldRemove {
			appState.widgets = append(appState.widgets[:i], appState.widgets[i+1:]...)
			log.I.F("Widget %d removed after animation completion", widget.id)
		}
	}

	// Third pass: layout remaining visible widgets
	for _, widget := range appState.widgets {
		// Only layout visible widgets
		if widget.visible {
			widget.Layout(gtx)
		}
	}

	// Handle right-click detection for widget creation (only when no scrim is visible)
	if !hasVisibleScrim {
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: appState.pointerTag,
				Kinds:  pointer.Press,
			})
			if !ok {
				break
			}
			if e, ok := ev.(pointer.Event); ok {
				if e.Kind == pointer.Press && e.Buttons == pointer.ButtonSecondary {
					clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))
					appState.createWidgetAt(clickPos, gtx.Constraints.Max)
					log.I.F("Right-click detected at (%d, %d)", clickPos.X, clickPos.Y)
				}
			}
		}
	} else {
		log.I.F("Scrim still visible, skipping right-click detection")
	}

	// Request animation frame if any widgets are animating
	if appState.animating {
		gtx.Execute(op.InvalidateCmd{})
	}
}

// createWidgetAt creates a new widget at the specified position
func (app *AppState) createWidgetAt(position image.Point, viewportSize image.Point) {
	widget := &ClickableWidget{
		position:     app.calculateSmartPosition(position, viewportSize),
		theme:        app.theme,
		visible:      true,
		scrimVisible: true,
		id:           app.nextWidgetID,
		clickable:    &widget.Clickable{},
		scrimClick:   &widget.Clickable{},
		showTime:     time.Now(),
		isHiding:     false,
		shouldRemove: false,
		closeButton: app.theme.NewButtonLayout().
			Background(app.theme.Colors.Error()).
			CornerRadius(0.5).
			Widget(func(g C) D {
				return app.theme.Caption("×").
					Color(app.theme.Colors.OnError()).
					Alignment(text.Middle).
					Layout(g)
			}),
	}

	app.widgets = append(app.widgets, widget)
	app.nextWidgetID++

	// Immediately request animation frame for the new widget
	app.animating = true
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

// Layout renders the widget at its position with scrim and animations
func (w *ClickableWidget) Layout(gtx C) D {
	now := time.Now()

	// Calculate animation progress
	var alpha float32 = 1.0
	if !w.isHiding {
		// Fade in animation
		elapsed := now.Sub(w.showTime)
		if elapsed < 250*time.Millisecond {
			alpha = float32(elapsed) / float32(250*time.Millisecond)
		}
	} else {
		// Fade out animation
		elapsed := now.Sub(w.hideTime)
		if elapsed < 250*time.Millisecond {
			alpha = 1.0 - (float32(elapsed) / float32(250*time.Millisecond))
		} else {
			// Animation complete, mark for removal
			w.shouldRemove = true
			return D{}
		}
	}

	// Handle close button clicks
	if w.closeButton.Clicked(gtx) {
		// Start hide animation instead of immediate removal
		w.startHideAnimation()
		w.scrimVisible = false // Hide scrim immediately
		log.I.F("Widget %d hide animation started via close button", w.id)
		// Note: widget will be removed after animation completes in mainUI
	}

	// Scrim is now handled at the main level, no event handling needed here

	// If scrim was clicked, continue with normal layout but don't render scrim
	// (the widget will be removed in the next frame by mainUI)

	// Only show scrim if it's still visible after handling clicks
	if w.scrimVisible {
		// Create scrim with animation alpha
		scrimAlpha := uint8(128 * alpha) // 50% opacity * animation alpha
		scrimColor := color.NRGBA{R: 0, G: 0, B: 0, A: scrimAlpha}
		paint.Fill(gtx.Ops, scrimColor)
	}

	// Position the widget
	offset := op.Offset(w.position).Push(gtx.Ops)
	defer offset.Pop()

	// Register for pointer events on the widget area to block them from reaching main handler
	widgetArea := clip.Rect(image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(150, 80),
	}).Push(gtx.Ops)
	event.Op(gtx.Ops, w.clickable)

	// Handle pointer events on the widget (only consume right-clicks)
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: w.clickable,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		if e, ok := ev.(pointer.Event); ok {
			if e.Kind == pointer.Press && e.Buttons == pointer.ButtonSecondary {
				log.I.F("Widget %d right-clicked at (%d, %d) - event consumed", w.id, int(e.Position.X), int(e.Position.Y))
				// Only consume right-clicks, let left-clicks pass through to scrim handler
			}
			// Don't consume left-clicks - let them pass through to the main handler
		}
	}

	widgetArea.Pop()

	// Apply widget alpha for fade animation
	widgetAlpha := uint8(255 * alpha)

	// Constrain widget size
	gtx.Constraints.Min.X = 150
	gtx.Constraints.Max.X = 150
	gtx.Constraints.Min.Y = 80
	gtx.Constraints.Max.Y = 80

	// Apply opacity to the entire widget
	opacity := float32(widgetAlpha) / 255.0
	opacityStack := paint.PushOpacity(gtx.Ops, opacity)
	defer opacityStack.Pop()

	// Create widget background with animation alpha
	return w.theme.NewCard(
		func(g C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(g,
				layout.Rigid(func(gtx C) D {
					// Header with close button
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return w.theme.Caption(fmt.Sprintf("Widget #%d", w.id)).
								Color(w.theme.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return w.closeButton.Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					// Content
					return w.theme.Caption("Right-click created this widget").
						Color(w.theme.Colors.OnSurfaceVariant()).
						Alignment(text.Middle).
						Layout(gtx)
				}),
			)
		},
	).CornerRadius(8).Padding(unit.Dp(8)).Layout(gtx)
}

// startHideAnimation starts the fade-out animation
func (w *ClickableWidget) startHideAnimation() {
	w.hideTime = time.Now()
	w.isHiding = true
}

// removeWidget removes this widget from the app state
func (w *ClickableWidget) removeWidget() {
	for i, widget := range appState.widgets {
		if widget.id == w.id {
			appState.widgets = append(appState.widgets[:i], appState.widgets[i+1:]...)
			log.I.F("Removed widget with ID %d", w.id)
			return
		}
	}
}
