package main

import (
	"context"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/mleku/fromage"
	"lol.mleku.dev/chk"
)

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
)

// CornerButton represents a button that shows a modal
type CornerButton struct {
	theme      *fromage.Theme
	trigger    *fromage.ButtonLayout
	modal      *Modal
	label      string
	buttonPos  image.Point
	buttonSize image.Point
	buttonRect image.Rectangle
}

// Modal represents a simple modal overlay
type Modal struct {
	theme      *fromage.Theme
	visible    bool
	position   image.Point
	label      string
	scrimClick *widget.Clickable
}

// Application state
type AppState struct {
	theme   *fromage.Theme
	button1 *CornerButton // NW
	button2 *CornerButton // NE
	button3 *CornerButton // SW
	button4 *CornerButton // SE
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

	w := fromage.NewWindow(th)

	// Initialize application state
	appState = &AppState{
		theme:   th,
		button1: NewCornerButton(th, "NW", "1"),
		button2: NewCornerButton(th, "NE", "2"),
		button3: NewCornerButton(th, "SW", "3"),
		button4: NewCornerButton(th, "SE", "4"),
	}

	w.Option(
		app.Size(unit.Dp(800), unit.Dp(600)),
		app.Title("Modal Positioning Demo"),
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

	// Draw center border lines (red cross)
	drawCenterLines(gtx, th)

	// Create flex-column with two flex-row containers
	th.VFlex().
		Flexed(1, func(g C) D {
			// Top row container
			return th.HFlex().
				Flexed(1, func(g C) D {
					return th.Direction().Center().Embed(func(g C) D {
						// NW button
						if appState.button1.trigger.Clicked(g) {
							appState.button1.ToggleModal()
						}
						dims := appState.button1.Layout(g)
						// Record actual button dimensions
						appState.button1.buttonSize = dims.Size
						return dims
					}).Fn(g)
				}).
				Flexed(1, func(g C) D {
					return th.Direction().Center().Embed(func(g C) D {
						// NE button
						if appState.button2.trigger.Clicked(g) {
							appState.button2.ToggleModal()
						}
						dims := appState.button2.Layout(g)
						// Record actual button dimensions
						appState.button2.buttonSize = dims.Size
						return dims
					}).Fn(g)
				}).
				Layout(g)
		}).
		Flexed(1, func(g C) D {
			// Bottom row container
			return th.HFlex().
				Flexed(1, func(g C) D {
					return th.Direction().Center().Embed(func(g C) D {
						// SW button
						if appState.button3.trigger.Clicked(g) {
							appState.button3.ToggleModal()
						}
						dims := appState.button3.Layout(g)
						// Record actual button dimensions
						appState.button3.buttonSize = dims.Size
						return dims
					}).Fn(g)
				}).
				Flexed(1, func(g C) D {
					return th.Direction().Center().Embed(func(g C) D {
						// SE button
						if appState.button4.trigger.Clicked(g) {
							appState.button4.ToggleModal()
						}
						dims := appState.button4.Layout(g)
						// Record actual button dimensions
						appState.button4.buttonSize = dims.Size
						return dims
					}).Fn(g)
				}).
				Layout(g)
		}).
		Layout(gtx)

	// Record button positions after layout with actual dimensions
	recordButtonPositions(gtx)

	// Layout the modals on top of everything
	appState.button1.LayoutModal(gtx)
	appState.button2.LayoutModal(gtx)
	appState.button3.LayoutModal(gtx)
	appState.button4.LayoutModal(gtx)
}

// drawCenterLines draws red cross lines in the center of the view
func drawCenterLines(gtx layout.Context, th *fromage.Theme) {
	screenWidth := gtx.Constraints.Max.X
	screenHeight := gtx.Constraints.Max.Y

	// Red color for the lines
	redColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	// Draw horizontal line (top to bottom)
	horizontalRect := image.Rect(screenWidth/2-1, 0, screenWidth/2+1, screenHeight)
	paint.FillShape(gtx.Ops, redColor, clip.Rect(horizontalRect).Op())

	// Draw vertical line (left to right)
	verticalRect := image.Rect(0, screenHeight/2-1, screenWidth, screenHeight/2+1)
	paint.FillShape(gtx.Ops, redColor, clip.Rect(verticalRect).Op())
}

// recordButtonPositions records the absolute positions of all buttons
func recordButtonPositions(gtx layout.Context) {
	screenWidth := gtx.Constraints.Max.X
	screenHeight := gtx.Constraints.Max.Y

	// Calculate button positions based on the 2x2 grid layout
	// Each button is centered in its quadrant using actual button dimensions

	// Top-left button (NW) - centered in top-left quadrant
	button1Size := appState.button1.buttonSize
	appState.button1.buttonPos = image.Pt(screenWidth/4-button1Size.X/2, screenHeight/4-button1Size.Y/2)
	appState.button1.buttonRect = image.Rect(
		screenWidth/4-button1Size.X/2, screenHeight/4-button1Size.Y/2,
		screenWidth/4+button1Size.X/2, screenHeight/4+button1Size.Y/2,
	)

	// Top-right button (NE) - centered in top-right quadrant
	button2Size := appState.button2.buttonSize
	appState.button2.buttonPos = image.Pt(3*screenWidth/4-button2Size.X/2, screenHeight/4-button2Size.Y/2)
	appState.button2.buttonRect = image.Rect(
		3*screenWidth/4-button2Size.X/2, screenHeight/4-button2Size.Y/2,
		3*screenWidth/4+button2Size.X/2, screenHeight/4+button2Size.Y/2,
	)

	// Bottom-left button (SW) - centered in bottom-left quadrant
	button3Size := appState.button3.buttonSize
	appState.button3.buttonPos = image.Pt(screenWidth/4-button3Size.X/2, 3*screenHeight/4-button3Size.Y/2)
	appState.button3.buttonRect = image.Rect(
		screenWidth/4-button3Size.X/2, 3*screenHeight/4-button3Size.Y/2,
		screenWidth/4+button3Size.X/2, 3*screenHeight/4+button3Size.Y/2,
	)

	// Bottom-right button (SE) - centered in bottom-right quadrant
	button4Size := appState.button4.buttonSize
	appState.button4.buttonPos = image.Pt(3*screenWidth/4-button4Size.X/2, 3*screenHeight/4-button4Size.Y/2)
	appState.button4.buttonRect = image.Rect(
		3*screenWidth/4-button4Size.X/2, 3*screenHeight/4-button4Size.Y/2,
		3*screenWidth/4+button4Size.X/2, 3*screenHeight/4+button4Size.Y/2,
	)
}

// NewCornerButton creates a new corner button with modal
func NewCornerButton(theme *fromage.Theme, buttonLabel, modalLabel string) *CornerButton {
	cb := &CornerButton{
		theme: theme,
		label: buttonLabel,
	}

	// Create trigger button
	cb.trigger = theme.NewButtonLayout().
		Background(theme.Colors.Surface()).
		CornerRadius(0.25).
		Widget(func(g C) D {
			return theme.Body1(buttonLabel).
				Color(theme.Colors.OnSurface()).
				Alignment(text.Middle).
				Layout(g)
		})

	// Create modal
	cb.modal = NewModal(theme, modalLabel)

	return cb
}

// ToggleModal toggles the visibility of the modal
func (cb *CornerButton) ToggleModal() {
	cb.modal.visible = !cb.modal.visible
}

// Layout renders the corner button
func (cb *CornerButton) Layout(gtx C) D {
	return cb.trigger.Layout(gtx)
}

// LayoutModal renders the modal with specific positioning rules
func (cb *CornerButton) LayoutModal(gtx C) {
	if !cb.modal.visible {
		return
	}

	screenWidth := gtx.Constraints.Max.X
	screenHeight := gtx.Constraints.Max.Y

	// Make modal exactly 16x16 text heights (square)
	textHeight := 16             // Base font height in pixels
	modalSize := textHeight * 16 // 256 pixels square
	modalWidth := modalSize
	modalHeight := modalSize

	// Use the stored button rectangle
	buttonRect := cb.buttonRect
	buttonCenterY := buttonRect.Min.Y + (buttonRect.Max.Y-buttonRect.Min.Y)/2
	viewportCenterY := screenHeight / 2

	// Determine vertical position based on button position relative to center
	var modalY int
	if buttonCenterY < viewportCenterY {
		// Button is above center, modal appears below button
		// Position so top-left corner of modal sits against bottom edge of button
		modalY = buttonRect.Max.Y
	} else {
		// Button is below center, modal appears above button
		// Position so bottom-left corner of modal sits against top edge of button
		modalY = buttonRect.Min.Y - modalHeight
	}

	// Horizontal position: left edge of modal aligned with left edge of button
	modalX := buttonRect.Min.X

	// Check if modal would overflow and clamp to edge while retaining dimensions
	if modalX+modalWidth > screenWidth {
		// Clamp to right edge but keep modal dimensions
		modalX = screenWidth - modalWidth
	}
	if modalX < 0 {
		// Clamp to left edge but keep modal dimensions
		modalX = 0
	}

	// Ensure modal doesn't go off the top or bottom while retaining dimensions
	if modalY < 0 {
		modalY = 0
	}
	if modalY+modalHeight > screenHeight {
		modalY = screenHeight - modalHeight
	}

	cb.modal.position = image.Pt(modalX, modalY)
	cb.modal.Layout(gtx)
}

// NewModal creates a new modal overlay
func NewModal(theme *fromage.Theme, label string) *Modal {
	return &Modal{
		theme:      theme,
		visible:    false,
		label:      label,
		scrimClick: &widget.Clickable{},
	}
}

// Layout renders the modal overlay
func (m *Modal) Layout(gtx C) {
	if !m.visible {
		return
	}

	// Handle scrim clicks
	if m.scrimClick.Clicked(gtx) {
		m.visible = false
		return
	}

	// Create scrim (dimmed background)
	scrimColor := color.NRGBA{R: 0, G: 0, B: 0, A: 64} // 25% opacity black
	paint.Fill(gtx.Ops, scrimColor)

	// Layout scrim clickable area
	m.scrimClick.Layout(gtx, func(gtx C) D {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	})

	// Position the modal content
	offset := op.Offset(m.position).Push(gtx.Ops)
	defer offset.Pop()

	// Constrain modal to exactly 16x16 text heights (256x256 pixels)
	modalSize := 256 // 16 * 16 pixels
	gtx.Constraints.Min.X = modalSize
	gtx.Constraints.Max.X = modalSize
	gtx.Constraints.Min.Y = modalSize
	gtx.Constraints.Max.Y = modalSize

	// Create modal content - exactly 16x16 text height square
	m.theme.NewCard(
		func(g C) D {
			return m.theme.Body2(m.label).
				Color(m.theme.Colors.OnSurface()).
				Alignment(text.Middle).
				Layout(g)
		},
	).CornerRadius(8).Padding(unit.Dp(8)).Layout(gtx)
}
