package main

import (
	"context"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
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

// No longer needed - removed button placement functionality

// RightClickPopup represents a popup that appears on right-click
type RightClickPopup struct {
	visible     bool
	position    image.Point
	closeButton *fromage.ButtonLayout
	scrimClick  *widget.Clickable
	theme       *fromage.Theme
}

// Application state
type AppState struct {
	rightClickGesture gesture.Click
	theme             *fromage.Theme
	popup             *RightClickPopup
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
		popup: &RightClickPopup{
			visible:    false,
			theme:      th,
			scrimClick: &widget.Clickable{},
			closeButton: th.NewButtonLayout().
				Background(th.Colors.Error()).
				CornerRadius(0.5).
				Widget(func(g C) D {
					return th.Caption("×").
						Color(th.Colors.OnError()).
						Alignment(text.Middle).
						Layout(g)
				}),
		},
	}

	w.Option(
		app.Size(unit.Dp(1200), unit.Dp(1200)),
		app.Title("Right-Click Popup Demo"),
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

	// Register click gesture area for the entire screen (for right-click detection)
	area := image.Rectangle{Max: gtx.Constraints.Max}
	defer clip.Rect(area).Push(gtx.Ops).Pop()

	// Handle right-click gestures for showing popup
	for {
		ev, ok := appState.rightClickGesture.Update(gtx.Source)
		if !ok {
			break
		}

		if ev.Kind == gesture.KindClick {
			// Show popup at the click position
			clickPos := image.Pt(int(ev.Position.X), int(ev.Position.Y))
			appState.popup.ShowPopup(clickPos, gtx.Constraints.Max)
		}
	}

	// Register right-click gesture area for the entire screen
	appState.rightClickGesture.Add(gtx.Ops)

	// No buttons to layout - this demo only shows right-click popups

	// Add instructions at the top
	instructions := th.Body1("Right-click anywhere to show a popup menu. Click the scrim or close button to dismiss it.").
		Color(th.Colors.OnBackground()).
		Alignment(text.Middle)

	// Position instructions at the top center
	instructionsOffset := op.Offset(image.Pt(gtx.Constraints.Max.X/2, 20)).Push(gtx.Ops)
	instructions.Layout(gtx)
	instructionsOffset.Pop()

	// Layout the popup on top of everything
	appState.popup.Layout(gtx)
}

// Removed button placement functions - no longer needed

// ShowPopup shows the popup at the specified position
func (p *RightClickPopup) ShowPopup(position image.Point, screenSize image.Point) {
	p.visible = true
	p.position = p.calculatePopupPosition(position, screenSize)
	log.I.F("Showing popup at position (%d, %d)", p.position.X, p.position.Y)
}

// HidePopup hides the popup
func (p *RightClickPopup) HidePopup() {
	p.visible = false
	log.I.F("Hiding popup")
}

// calculatePopupPosition calculates where to position the popup so the corner faces away from center
func (p *RightClickPopup) calculatePopupPosition(clickPos image.Point, screenSize image.Point) image.Point {
	centerX := screenSize.X / 2
	centerY := screenSize.Y / 2

	popupWidth := 200  // Approximate popup width
	popupHeight := 100 // Approximate popup height

	// Determine which corner should face away from center
	if clickPos.X < centerX {
		// Click is on left side, position popup to the right
		if clickPos.Y < centerY {
			// Click is in top-left, position popup bottom-right of click
			return image.Pt(clickPos.X, clickPos.Y)
		} else {
			// Click is in bottom-left, position popup top-right of click
			return image.Pt(clickPos.X, clickPos.Y-popupHeight)
		}
	} else {
		// Click is on right side, position popup to the left
		if clickPos.Y < centerY {
			// Click is in top-right, position popup bottom-left of click
			return image.Pt(clickPos.X-popupWidth, clickPos.Y)
		} else {
			// Click is in bottom-right, position popup top-left of click
			return image.Pt(clickPos.X-popupWidth, clickPos.Y-popupHeight)
		}
	}
}

// Layout renders the popup if it's visible
func (p *RightClickPopup) Layout(gtx C) D {
	if !p.visible {
		return D{}
	}

	// Handle scrim clicks
	if p.scrimClick.Clicked(gtx) {
		p.HidePopup()
		return D{}
	}

	// Handle close button clicks
	if p.closeButton.Clicked(gtx) {
		p.HidePopup()
		return D{}
	}

	// Create scrim (dimmed background)
	scrimColor := color.NRGBA{R: 0, G: 0, B: 0, A: 128} // 50% opacity black
	paint.Fill(gtx.Ops, scrimColor)

	// Layout scrim clickable area
	p.scrimClick.Layout(gtx, func(gtx C) D {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	})

	// Position the popup
	offset := op.Offset(p.position).Push(gtx.Ops)
	defer offset.Pop()

	// Constrain popup size
	gtx.Constraints.Min.X = 200
	gtx.Constraints.Max.X = 200
	gtx.Constraints.Min.Y = 100
	gtx.Constraints.Max.Y = 100

	// Create popup background
	return p.theme.NewCard(
		func(g C) D {
			return p.theme.VFlex().
				Rigid(func(gtx C) D {
					// Title
					return p.theme.Body2("Right-click Popup").
						Color(p.theme.Colors.OnSurface()).
						Alignment(text.Middle).
						Layout(gtx)
				}).
				Rigid(func(gtx C) D {
					// Content
					return p.theme.Caption("This popup appeared because you right-clicked!").
						Color(p.theme.Colors.OnSurfaceVariant()).
						Alignment(text.Middle).
						Layout(gtx)
				}).
				Rigid(func(gtx C) D {
					// Close button
					return p.closeButton.Layout(gtx)
				}).
				Layout(g)
		},
	).CornerRadius(8).Padding(unit.Dp(12)).Layout(gtx)
}
