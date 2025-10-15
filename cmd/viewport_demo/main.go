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
	"github.com/mleku/fromage"
	"lol.mleku.dev/chk"
)

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
)

// RedCornerOutline creates a widget that draws a red 1px square corner outline
// filled with smaller outlined squares, sized to be an even multiple of square size
func RedCornerOutline(gtx layout.Context, th *fromage.Theme, contentSize int) layout.Dimensions {
	// Draw the main red outline
	drawOutline(gtx, 0, 0, contentSize, contentSize, color.NRGBA{R: 255, A: 255})

	// Calculate square size (6 text heights)
	squareSize := gtx.Dp(th.TextSize * 6)

	// Calculate how many squares fit in each direction
	squaresX := contentSize / squareSize
	squaresY := contentSize / squareSize

	// Draw grid of smaller outlined squares
	for y := 0; y < squaresY; y++ {
		for x := 0; x < squaresX; x++ {
			squareX := x * squareSize
			squareY := y * squareSize
			drawOutline(gtx, squareX, squareY, squareSize, squareSize, color.NRGBA{R: 255, A: 255})
		}
	}

	return layout.Dimensions{Size: image.Pt(contentSize, contentSize)}
}

// clipViewport creates a clipped viewport that shows only the portion of content
// corresponding to the scrollbar positions
func clipViewport(gtx layout.Context, th *fromage.Theme, contentSize int, horizontalPos, verticalPos float32, viewportWidth, viewportHeight int, testButton *fromage.Button, popoutMenu *fromage.Popout) layout.Dimensions {
	// Calculate the offset based on scroll position
	// horizontalPos and verticalPos are 0-1, so we multiply by the scrollable distance
	scrollableWidth := contentSize - viewportWidth
	scrollableHeight := contentSize - viewportHeight

	// Calculate offsets - when scrollbar is at 0, show top-left; when at 1, show bottom-right
	offsetX := int(float32(scrollableWidth) * horizontalPos)
	offsetY := int(float32(scrollableHeight) * verticalPos)

	// Ensure we don't have negative scrollable distances
	if scrollableWidth < 0 {
		scrollableWidth = 0
		offsetX = 0
	}
	if scrollableHeight < 0 {
		scrollableHeight = 0
		offsetY = 0
	}

	// Apply translation to move the content based on scroll position
	// Negative offset moves content in opposite direction of scroll
	offsetArea := op.Offset(image.Pt(-offsetX, -offsetY)).Push(gtx.Ops)

	// Draw the red corner outline widget at the translated position
	RedCornerOutline(gtx, th, contentSize)

	// Draw the test button in the center of the canvas
	centerX := contentSize / 2
	centerY := contentSize / 2
	buttonArea := op.Offset(image.Pt(centerX, centerY)).Push(gtx.Ops)

	// Layout the button and popout menu
	th.CenteredColumn().
		Rigid(func(g C) D {
			return testButton.Layout(g, th)
		}).
		Layout(gtx)

	buttonArea.Pop()
	offsetArea.Pop()

	return layout.Dimensions{Size: image.Pt(viewportWidth, viewportHeight)}
}

// drawOutline draws a 1px outline rectangle at the specified position and size
func drawOutline(gtx layout.Context, x, y, width, height int, col color.NRGBA) {
	// Top edge
	topRect := image.Rect(x, y, x+width, y+1)
	topArea := clip.Rect(topRect).Push(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	topArea.Pop()

	// Bottom edge
	bottomRect := image.Rect(x, y+height-1, x+width, y+height)
	bottomArea := clip.Rect(bottomRect).Push(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	bottomArea.Pop()

	// Left edge
	leftRect := image.Rect(x, y, x+1, y+height)
	leftArea := clip.Rect(leftRect).Push(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	leftArea.Pop()

	// Right edge
	rightRect := image.Rect(x+width-1, y, x+width, y+height)
	rightArea := clip.Rect(rightRect).Push(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	rightArea.Pop()
}

func main() {
	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeDark,
	)

	w := fromage.NewWindow(th)
	w.Option(app.Size(
		unit.Dp(800), unit.Dp(800)),
		app.Title("Viewport Demo"),
	)
	w.Run(loop(w.Window, th, w))
}

func loop(w *app.Window, th *fromage.Theme, window *fromage.Window) func() {
	// Horizontal scrollbar for bottom edge
	horizontalScrollbar := th.NewScrollbar(fromage.Horizontal)

	// Vertical scrollbar for right edge
	verticalScrollbar := th.NewScrollbar(fromage.Vertical)

	// Button with popout menu for testing
	testButton := th.NewButton("Test Popout")
	popoutMenu := th.NewPopout()

	var horizontalPos float32 = 0.0
	var verticalPos float32 = 0.0

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

				// Update scrollbar positions if they changed
				if horizontalScrollbar.Changed() {
					horizontalPos = horizontalScrollbar.Position()
				}
				if verticalScrollbar.Changed() {
					verticalPos = verticalScrollbar.Position()
				}

				mainUI(gtx, th, window, horizontalScrollbar, verticalScrollbar, testButton, popoutMenu, horizontalPos, verticalPos)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, window *fromage.Window,
	horizontalScrollbar, verticalScrollbar *fromage.Scrollbar,
	testButton *fromage.Button, popoutMenu *fromage.Popout,
	horizontalPos, verticalPos float32) {

	// Fill background with theme background color
	th.FillBackground(nil).Layout(gtx)

	// Calculate content size to be an even multiple of 6 text height squares
	squareSize := gtx.Dp(th.TextSize * 6) // 6 text heights per square
	// Calculate how many squares fit in 1600 Dp, then round down to even multiple
	squaresIn1600 := 1600 / int(th.TextSize*6)
	// Use an even number of squares (round down to even)
	if squaresIn1600%2 != 0 {
		squaresIn1600--
	}
	contentSize := squaresIn1600 * squareSize

	// Calculate available viewport size (window size minus scrollbar space)
	scrollbarWidth := gtx.Dp(th.TextSize)
	viewportWidth := gtx.Constraints.Max.X - scrollbarWidth
	viewportHeight := gtx.Constraints.Max.Y - scrollbarWidth

	// Get the actual content area size (viewport area minus horizontal scrollbar space)
	contentAreaWidth := viewportWidth
	contentAreaHeight := viewportHeight - scrollbarWidth // Subtract horizontal scrollbar height

	// Calculate viewport proportions (how much of the content is visible)
	horizontalViewport := float32(contentAreaWidth) / float32(contentSize)
	verticalViewport := float32(contentAreaHeight) / float32(contentSize)

	// Clamp viewport proportions to valid range
	if horizontalViewport > 1.0 {
		horizontalViewport = 1.0
	}
	if verticalViewport > 1.0 {
		verticalViewport = 1.0
	}

	// Update scrollbar viewport proportions
	horizontalScrollbar.SetViewport(horizontalViewport)
	verticalScrollbar.SetViewport(verticalViewport)

	// Main layout with headers at top, scrollable viewport in middle, scrollbars on edges
	th.VFlex().
		Rigid(func(g C) D {
			// Headers at top - full width
			return th.CenteredColumn().
				Rigid(func(g C) D {
					return th.H4("Viewport Demo").Alignment(text.Middle).Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Body1("This is the main content area").Layout(g)
				}).
				Layout(g)
		}).
		Flexed(1, func(g C) D {
			// Scrollable viewport area with scrollbars
			return th.VFlex().
				Flexed(1, func(g C) D {
					// Main viewport area with vertical scrollbar on right
					return th.HFlex().
						Flexed(1, func(g C) D {
							// Apply clipping mask to the exact area bounded by scrollbars
							clipRect := image.Rect(0, 0, g.Constraints.Max.X, g.Constraints.Max.Y)
							clipArea := clip.Rect(clipRect).Push(g.Ops)

							// Clipped viewport for the red corner outline widget
							dims := clipViewport(g, th, contentSize, horizontalPos, verticalPos, contentAreaWidth, contentAreaHeight, testButton, popoutMenu)

							clipArea.Pop()
							return dims
						}).
						Rigid(func(g C) D {
							// Vertical scrollbar on right edge
							return verticalScrollbar.Layout(g, th)
						}).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Bottom area with horizontal scrollbar and corner space
					return th.HFlex().
						Flexed(1, func(g C) D {
							// Horizontal scrollbar on bottom edge (stops before corner)
							return horizontalScrollbar.Layout(g, th)
						}).
						Rigid(func(g C) D {
							// Square corner space (same size as scrollbar width)
							return layout.Dimensions{
								Size: image.Pt(scrollbarWidth, scrollbarWidth),
							}
						}).
						Layout(g)
				}).
				Layout(g)
		}).
		Layout(gtx)

	// Layout the popout menu (outside the clipped area so it can appear above other content)
	if testButton.Clicked() {
		popoutMenu.Toggle()
	}

	// Position the popout menu near the button (this will be relative to the window, not the scrolled content)
	if popoutMenu.Visible() {
		// Calculate where the button would be in window coordinates
		buttonWindowX := int(float32(contentSize/2) * (1 - horizontalPos))
		buttonWindowY := int(float32(contentSize/2) * (1 - verticalPos))

		popoutArea := op.Offset(image.Pt(buttonWindowX, buttonWindowY)).Push(gtx.Ops)
		popoutMenu.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
			return th.CenteredColumn().
				Rigid(func(g C) D {
					return th.Body1("Popout Menu Item 1").Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Body1("Popout Menu Item 2").Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Body1("Popout Menu Item 3").Layout(g)
				}).
				Layout(gtx)
		})
		popoutArea.Pop()
	}
}
