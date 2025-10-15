package main

import (
	"context"
	"image"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
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

// PlacedButton represents a button that has been placed at a specific position
type PlacedButton struct {
	Position image.Point
	Button   *fromage.ButtonLayout
	ID       int
}

// Application state
type AppState struct {
	buttons      []PlacedButton
	nextID       int
	clickGesture gesture.Click
	theme        *fromage.Theme
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
		buttons: make([]PlacedButton, 0),
		nextID:  1,
		theme:   th,
	}

	w.Option(
		app.Size(unit.Dp(1200), unit.Dp(1200)),
		app.Title("Click to Place Buttons Demo"),
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

	// Handle click gestures for placing buttons
	for {
		ev, ok := appState.clickGesture.Update(gtx.Source)
		if !ok {
			break
		}

		if ev.Kind == gesture.KindClick {
			// Place a new button at the click position
			clickPos := image.Pt(int(ev.Position.X), int(ev.Position.Y))
			placeButton(gtx, th, clickPos)
		}
	}

	// Register click gesture area for the entire screen
	area := image.Rectangle{Max: gtx.Constraints.Max}
	defer clip.Rect(area).Push(gtx.Ops).Pop()
	appState.clickGesture.Add(gtx.Ops)

	// Layout all placed buttons
	for i := range appState.buttons {
		button := &appState.buttons[i]

		// Position the button at its stored position
		offset := op.Offset(button.Position).Push(gtx.Ops)

		// Check if this button was clicked
		if button.Button.Clicked(gtx) {
			log.I.F("Button %d clicked - removing this button", button.ID)
			removeButton(i)
			offset.Pop()
			break // Exit the loop since we modified the slice
		}

		// Layout the button
		button.Button.Layout(gtx)

		offset.Pop()
	}

	// Add instructions at the top
	instructions := th.Body1("Click anywhere to place a small button. Click any button to remove just that button.").
		Color(th.Colors.OnBackground()).
		Alignment(text.Middle)

	// Position instructions at the top center
	instructionsOffset := op.Offset(image.Pt(gtx.Constraints.Max.X/2, 20)).Push(gtx.Ops)
	instructions.Layout(gtx)
	instructionsOffset.Pop()
}

// placeButton creates a new button at the specified position
func placeButton(gtx layout.Context, th *fromage.Theme, position image.Point) {
	// Calculate button size (3 text heights wide and tall)
	textHeight := gtx.Dp(unit.Dp(float32(th.TextSize)))
	buttonSize := textHeight * 3

	// Create a new button with fixed size
	button := th.NewButtonLayout().
		Background(th.Colors.Primary()).
		CornerRadius(0.5).
		Widget(func(g C) D {
			// Constrain the button to the desired size
			g.Constraints.Min.X = buttonSize
			g.Constraints.Max.X = buttonSize
			g.Constraints.Min.Y = buttonSize
			g.Constraints.Max.Y = buttonSize

			return th.Caption("•").
				Color(th.Colors.OnPrimary()).
				Alignment(text.Middle).
				Layout(g)
		})

	// Create the placed button
	placedButton := PlacedButton{
		Position: position,
		Button:   button,
		ID:       appState.nextID,
	}

	// Add to the list
	appState.buttons = append(appState.buttons, placedButton)
	appState.nextID++

	log.I.F("Placed button %d at position (%d, %d)", placedButton.ID, position.X, position.Y)
}

// removeButton removes a button at the specified index
func removeButton(index int) {
	if index < 0 || index >= len(appState.buttons) {
		return
	}

	// Remove the button at the specified index
	appState.buttons = append(appState.buttons[:index], appState.buttons[index+1:]...)
	log.I.F("Removed button at index %d", index)
}

// clearAllButtons removes all placed buttons
func clearAllButtons() {
	appState.buttons = appState.buttons[:0] // Clear the slice
	log.I.F("Cleared all buttons")
}
