package main

import (
	"context"
	"fmt"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
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
func clipViewport(gtx layout.Context, th *fromage.Theme, contentSize int, horizontalPos, verticalPos float32, viewportWidth, viewportHeight int, testButton *fromage.ButtonLayout, modalStack *fromage.ModalStack) layout.Dimensions {
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

	// Layout the button (just draw it, click handling is done outside)
	th.CenteredColumn().
		Rigid(func(g C) D {
			return testButton.Layout(g)
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

	// Button with modal for testing
	testButton := th.NewButtonLayout().
		Background(th.Colors.Primary()).
		CornerRadius(0.5).
		Widget(func(g C) D {
			return th.Body1("Test Modal").
				Color(th.Colors.OnPrimary()).
				Alignment(text.Middle).
				Layout(g)
		})
	modalStack := th.NewModalStack()

	// Scroll gesture for handling scroll wheel events
	scrollGesture := gesture.Scroll{}
	// Pointer tag for top-level events
	pointerTag := &struct{}{}
	// Gesture tag for gesture events
	gestureTag := &struct{}{}

	var horizontalPos float32 = 0.0
	var verticalPos float32 = 0.0

	// Keyboard event flags for direct scrolling
	var keyUpPressed, keyDownPressed, keyLeftPressed, keyRightPressed bool

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

				// Register a global key listener for arrow keys
				area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
				event.Op(gtx.Ops, w)

				// Handle keyboard events for scrolling
				for {
					event, ok := gtx.Event(
						key.Filter{Name: key.NameUpArrow},
						key.Filter{Name: key.NameDownArrow},
						key.Filter{Name: key.NameLeftArrow},
						key.Filter{Name: key.NameRightArrow},
					)
					if !ok {
						break
					}
					switch event := event.(type) {
					case key.Event:
						fmt.Printf("Key event: %s, state: %s\n", event.Name, event.State)
						if event.State == key.Press {
							switch event.Name {
							case key.NameUpArrow:
								// Direct upward scroll movement
								keyUpPressed = true
								fmt.Println("UP ARROW PRESSED - Direct upward scroll")
							case key.NameDownArrow:
								// Direct downward scroll movement
								keyDownPressed = true
								fmt.Println("DOWN ARROW PRESSED - Direct downward scroll")
							case key.NameLeftArrow:
								// Direct leftward scroll movement
								keyLeftPressed = true
								fmt.Println("LEFT ARROW PRESSED - Direct leftward scroll")
							case key.NameRightArrow:
								// Direct rightward scroll movement
								keyRightPressed = true
								fmt.Println("RIGHT ARROW PRESSED - Direct rightward scroll")
							}
						}
					}
				}
				area.Pop()

				mainUI(gtx, th, window, horizontalScrollbar, verticalScrollbar, testButton, modalStack, &scrollGesture, pointerTag, gestureTag, &horizontalPos, &verticalPos, &keyUpPressed, &keyDownPressed, &keyLeftPressed, &keyRightPressed)

				// Update scrollbar positions if they changed (after processing animations)
				if horizontalScrollbar.Changed() {
					horizontalPos = horizontalScrollbar.Position()
				}
				if verticalScrollbar.Changed() {
					verticalPos = verticalScrollbar.Position()
				}
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, window *fromage.Window,
	horizontalScrollbar, verticalScrollbar *fromage.Scrollbar,
	testButton *fromage.ButtonLayout, modalStack *fromage.ModalStack,
	scrollGesture *gesture.Scroll, pointerTag, gestureTag interface{},
	horizontalPos, verticalPos *float32, keyUpPressed, keyDownPressed, keyLeftPressed, keyRightPressed *bool) {

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
	// Reduce by 4 squares to make the background widget smaller
	squaresIn1600 -= 4
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
					// Main viewport area with vertical scrollbar on right (if needed)
					return th.HFlex().
						Flexed(1, func(g C) D {
							// Apply clipping mask to the exact area bounded by scrollbars
							clipRect := image.Rect(0, 0, g.Constraints.Max.X, g.Constraints.Max.Y)
							clipArea := clip.Rect(clipRect).Push(g.Ops)

							// Register scroll gesture over the viewport area
							scrollGesture.Add(g.Ops)

							// Clipped viewport for the red corner outline widget
							dims := clipViewport(g, th, contentSize, *horizontalPos, *verticalPos, contentAreaWidth, contentAreaHeight, testButton, modalStack)

							clipArea.Pop()
							return dims
						}).
						Rigid(func(g C) D {
							// Vertical scrollbar on right edge (only show if content is larger than viewport)
							if verticalViewport < 1.0 {
								return verticalScrollbar.Layout(g, th)
							}
							// Return empty space when scrollbar is not needed
							return layout.Dimensions{Size: image.Pt(0, g.Constraints.Max.Y)}
						}).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Bottom area with horizontal scrollbar and corner space (if needed)
					return th.HFlex().
						Flexed(1, func(g C) D {
							// Horizontal scrollbar on bottom edge (only show if content is larger than viewport)
							if horizontalViewport < 1.0 {
								return horizontalScrollbar.Layout(g, th)
							}
							// Return empty space when scrollbar is not needed
							return layout.Dimensions{Size: image.Pt(g.Constraints.Max.X, 0)}
						}).
						Rigid(func(g C) D {
							// Square corner space (only show if both scrollbars are visible)
							if horizontalViewport < 1.0 && verticalViewport < 1.0 {
								return layout.Dimensions{
									Size: image.Pt(scrollbarWidth, scrollbarWidth),
								}
							}
							// Return empty space when corner is not needed
							return layout.Dimensions{Size: image.Pt(0, 0)}
						}).
						Layout(g)
				}).
				Layout(g)
		}).
		Layout(gtx)

	// Handle scroll events - direct 10% jumps
	scrollableWidth := contentSize - contentAreaWidth
	scrollableHeight := contentSize - contentAreaHeight

	// Use large values to allow scrolling in both directions
	min, max := int(-1e6), int(1e6)
	xrange := pointer.ScrollRange{Min: min, Max: max}
	yrange := pointer.ScrollRange{Min: min, Max: max}

	// Update scroll gesture for both vertical and horizontal
	verticalScrollDistance := scrollGesture.Update(gtx.Metric, gtx.Source, gtx.Now, gesture.Vertical, xrange, yrange)
	horizontalScrollDistance := scrollGesture.Update(gtx.Metric, gtx.Source, gtx.Now, gesture.Horizontal, xrange, yrange)

	// Handle vertical scroll input - direct 10% viewport movements
	if verticalScrollDistance != 0 {
		fmt.Printf("Vertical scroll distance: %d\n", verticalScrollDistance)

		if scrollableHeight > 0 {
			// Calculate position change for 10% of visible viewport height
			viewportHeight := contentAreaHeight
			scrollDistancePixels := float32(viewportHeight) * 0.1 // 10% of viewport height
			positionChange := scrollDistancePixels / float32(scrollableHeight)

			// Determine scroll direction and apply directly
			var newPos float32
			if verticalScrollDistance > 0 {
				// Scroll down - move content up (increase position)
				newPos = *verticalPos + positionChange
				fmt.Println("Scroll DOWN - moving content up")
			} else {
				// Scroll up - move content down (decrease position)
				newPos = *verticalPos - positionChange
				fmt.Println("Scroll UP - moving content down")
			}

			// Clamp to valid range
			if newPos < 0.0 {
				newPos = 0.0
			} else if newPos > 1.0 {
				newPos = 1.0
			}

			// Apply the change directly
			*verticalPos = newPos
			verticalScrollbar.SetPosition(newPos)
			fmt.Printf("New vertical position: %.3f\n", newPos)
		}
	}

	// Handle horizontal scroll input - direct 10% viewport movements
	if horizontalScrollDistance != 0 {
		fmt.Printf("Horizontal scroll distance: %d\n", horizontalScrollDistance)

		if scrollableWidth > 0 {
			// Calculate position change for 10% of visible viewport width
			viewportWidth := contentAreaWidth
			scrollDistancePixels := float32(viewportWidth) * 0.1 // 10% of viewport width
			positionChange := scrollDistancePixels / float32(scrollableWidth)

			// Determine scroll direction and apply directly
			var newPos float32
			if horizontalScrollDistance > 0 {
				// Scroll right - move content left (increase position)
				newPos = *horizontalPos + positionChange
				fmt.Println("Scroll RIGHT - moving content left")
			} else {
				// Scroll left - move content right (decrease position)
				newPos = *horizontalPos - positionChange
				fmt.Println("Scroll LEFT - moving content right")
			}

			// Clamp to valid range
			if newPos < 0.0 {
				newPos = 0.0
			} else if newPos > 1.0 {
				newPos = 1.0
			}

			// Apply the change directly
			*horizontalPos = newPos
			horizontalScrollbar.SetPosition(newPos)
			fmt.Printf("New horizontal position: %.3f\n", newPos)
		}
	}

	// Handle keyboard events - direct 10% jumps
	if *keyUpPressed {
		*keyUpPressed = false // Reset flag
		if scrollableHeight > 0 {
			// Calculate position change for 10% of visible viewport height
			viewportHeight := contentAreaHeight
			scrollDistancePixels := float32(viewportHeight) * 0.1 // 10% of viewport height
			positionChange := scrollDistancePixels / float32(scrollableHeight)

			// Up arrow - move content down (decrease position)
			newPos := *verticalPos - positionChange

			// Clamp to valid range
			if newPos < 0.0 {
				newPos = 0.0
			}

			// Apply the change directly
			*verticalPos = newPos
			verticalScrollbar.SetPosition(newPos)
			fmt.Printf("UP ARROW: New vertical position: %.3f\n", newPos)
		}
	}

	if *keyDownPressed {
		*keyDownPressed = false // Reset flag
		if scrollableHeight > 0 {
			// Calculate position change for 10% of visible viewport height
			viewportHeight := contentAreaHeight
			scrollDistancePixels := float32(viewportHeight) * 0.1 // 10% of viewport height
			positionChange := scrollDistancePixels / float32(scrollableHeight)

			// Down arrow - move content up (increase position)
			newPos := *verticalPos + positionChange

			// Clamp to valid range
			if newPos > 1.0 {
				newPos = 1.0
			}

			// Apply the change directly
			*verticalPos = newPos
			verticalScrollbar.SetPosition(newPos)
			fmt.Printf("DOWN ARROW: New vertical position: %.3f\n", newPos)
		}
	}

	// Handle horizontal keyboard events - direct 10% jumps
	if *keyLeftPressed {
		*keyLeftPressed = false // Reset flag
		if scrollableWidth > 0 {
			// Calculate position change for 10% of visible viewport width
			viewportWidth := contentAreaWidth
			scrollDistancePixels := float32(viewportWidth) * 0.1 // 10% of viewport width
			positionChange := scrollDistancePixels / float32(scrollableWidth)

			// Left arrow - move content right (decrease position)
			newPos := *horizontalPos - positionChange

			// Clamp to valid range
			if newPos < 0.0 {
				newPos = 0.0
			}

			// Apply the change directly
			*horizontalPos = newPos
			horizontalScrollbar.SetPosition(newPos)
			fmt.Printf("LEFT ARROW: New horizontal position: %.3f\n", newPos)
		}
	}

	if *keyRightPressed {
		*keyRightPressed = false // Reset flag
		if scrollableWidth > 0 {
			// Calculate position change for 10% of visible viewport width
			viewportWidth := contentAreaWidth
			scrollDistancePixels := float32(viewportWidth) * 0.1 // 10% of viewport width
			positionChange := scrollDistancePixels / float32(scrollableWidth)

			// Right arrow - move content left (increase position)
			newPos := *horizontalPos + positionChange

			// Clamp to valid range
			if newPos > 1.0 {
				newPos = 1.0
			}

			// Apply the change directly
			*horizontalPos = newPos
			horizontalScrollbar.SetPosition(newPos)
			fmt.Printf("RIGHT ARROW: New horizontal position: %.3f\n", newPos)
		}
	}

	// Handle clicks in the center area to show modal
	// Calculate the center area of the viewport
	centerX := contentAreaWidth / 2
	centerY := contentAreaHeight / 2
	buttonSize := gtx.Dp(th.TextSize * 4) // Approximate button size

	// Create a clickable area for the center button
	centerRect := image.Rect(
		centerX-buttonSize/2, centerY-buttonSize/2,
		centerX+buttonSize/2, centerY+buttonSize/2,
	)

	// Check if there's a click in the center area
	clickable := th.Pool.GetClickable()
	clickArea := clip.Rect(centerRect).Push(gtx.Ops)
	event.Op(gtx.Ops, clickable)
	clickArea.Pop()

	if clickable.Clicked(gtx) {
		// Create modal content
		modalContent := func(g C) D {
			return th.NewCard(
				func(g C) D {
					return th.VFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							return th.H3("Test Modal").
								Color(th.Colors.OnSurface()).
								Alignment(text.Middle).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1("This modal appears when you click the button in the center of the scrolled content.").
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1("The modal should maintain its state even when scrolling.").
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Rigid(func(g C) D {
							// Close button
							btn := th.NewButtonLayout().
								Background(th.Colors.Secondary()).
								CornerRadius(0.5).
								Widget(func(g C) D {
									return th.Body2("Close").
										Color(th.Colors.OnSecondary()).
										Alignment(text.Middle).
										Layout(g)
								})
							if btn.Clicked(g) {
								modalStack.Pop()
							}
							return btn.Layout(g)
						}).
						Layout(g)
				},
			).CornerRadius(8).Padding(unit.Dp(16)).Layout(g)
		}

		// Push the modal to the stack
		modalStack.Push(modalContent, func() {
			modalStack.Pop()
		})
	}

	// Layout the modal stack on top of everything
	if !modalStack.IsEmpty() {
		modalStack.Layout(gtx)
	}
}
