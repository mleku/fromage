package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"

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

// Physics configuration - tune these values to adjust motion behavior
var (
	// Mass affects how much impulse is needed to achieve the same velocity
	// Higher mass = more impulse needed = less responsive feeling
	// Lower mass = less impulse needed = more responsive feeling
	// Range: 0.1 (very responsive) to 10.0 (very sluggish)
	// Reduced by 1000x for much lighter viewport
	PhysicsMass float32 = 0.001

	// Impulse strength per scroll event (pixels per second)
	// This is the base impulse before mass is applied
	// Quadrupled for 4x faster velocity, then 4x more for 4x scroll distance, then 100x more
	BaseImpulseStrength float32 = 160000.0

	// Mouse scroll energy multiplier - makes mouse scrolling more powerful
	// Higher values = more energy per mouse scroll event
	// Range: 1.0 (same as keyboard) to 10.0 (very powerful)
	MouseScrollMultiplier float32 = 3.0

	// Momentum mass - affects how long motion continues after input stops
	// Higher values = longer momentum, more continuous motion
	// Lower values = motion stops more quickly
	// Range: 0.1 (stops quickly) to 10.0 (very long momentum)
	MomentumMass float32 = 5.0
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

	// Physics state for inertial scrolling
	var physicsState struct {
		// Velocity in pixels per second
		velocityX, velocityY float32
		// Position in pixels (for smooth interpolation)
		positionX, positionY float32
		// Friction coefficient (0.95 = 5% energy loss per frame at 60fps)
		friction float32
		// Maximum speed in pixels per second
		maxSpeed float32
		// Bounce energy retention (0.5 = 50%)
		bounceRetention float32
	}

	// Initialize physics parameters
	physicsState.friction = 0.9999 // Much lower energy loss for longer momentum
	physicsState.maxSpeed = 5000.0
	physicsState.bounceRetention = 0.5

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

				mainUI(gtx, th, window, horizontalScrollbar, verticalScrollbar, testButton, modalStack, &scrollGesture, pointerTag, gestureTag, &horizontalPos, &verticalPos, &keyUpPressed, &keyDownPressed, &keyLeftPressed, &keyRightPressed, &physicsState)

				// Note: Scrollbar positions are now controlled by physics system
				// No need to update from scrollbar changes since physics sets them directly
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, window *fromage.Window,
	horizontalScrollbar, verticalScrollbar *fromage.Scrollbar,
	testButton *fromage.ButtonLayout, modalStack *fromage.ModalStack,
	scrollGesture *gesture.Scroll, pointerTag, gestureTag interface{},
	horizontalPos, verticalPos *float32, keyUpPressed, keyDownPressed, keyLeftPressed, keyRightPressed *bool,
	physicsState *struct {
		velocityX, velocityY float32
		positionX, positionY float32
		friction             float32
		maxSpeed             float32
		bounceRetention      float32
	}) {

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

	// Physics-based inertial scrolling
	scrollableWidth := contentSize - contentAreaWidth
	scrollableHeight := contentSize - contentAreaHeight

	// Use large values to allow scrolling in both directions
	min, max := int(-1e6), int(1e6)
	xrange := pointer.ScrollRange{Min: min, Max: max}
	yrange := pointer.ScrollRange{Min: min, Max: max}

	// Update scroll gesture for both vertical and horizontal
	verticalScrollDistance := scrollGesture.Update(gtx.Metric, gtx.Source, gtx.Now, gesture.Vertical, xrange, yrange)
	horizontalScrollDistance := scrollGesture.Update(gtx.Metric, gtx.Source, gtx.Now, gesture.Horizontal, xrange, yrange)

	// Convert scroll events to impulse (velocity changes)
	// Apply mass to the base impulse: higher mass = less effective impulse
	// Mouse scroll gets extra energy multiplier for faster scrolling
	// Mouse scroll mass effect is divided by 10000 for extremely light inertia
	mouseImpulseStrength := (BaseImpulseStrength * MouseScrollMultiplier) / (PhysicsMass / 10000.0)
	keyboardImpulseStrength := BaseImpulseStrength / PhysicsMass

	if verticalScrollDistance != 0 {
		fmt.Printf("Vertical scroll distance: %d\n", verticalScrollDistance)
		if scrollableHeight > 0 {
			var impulse float32
			if verticalScrollDistance > 0 {
				// Scroll down - add upward impulse
				impulse = mouseImpulseStrength
				fmt.Println("Scroll DOWN - adding upward impulse")
			} else {
				// Scroll up - add downward impulse
				impulse = -mouseImpulseStrength
				fmt.Println("Scroll UP - adding downward impulse")
			}
			physicsState.velocityY += impulse
			fmt.Printf("New vertical velocity: %.1f px/s (mouse: %.1f)\n", physicsState.velocityY, mouseImpulseStrength)
		}
	}

	if horizontalScrollDistance != 0 {
		fmt.Printf("Horizontal scroll distance: %d\n", horizontalScrollDistance)
		if scrollableWidth > 0 {
			var impulse float32
			if horizontalScrollDistance > 0 {
				// Scroll right - add leftward impulse
				impulse = mouseImpulseStrength
				fmt.Println("Scroll RIGHT - adding leftward impulse")
			} else {
				// Scroll left - add rightward impulse
				impulse = -mouseImpulseStrength
				fmt.Println("Scroll LEFT - adding rightward impulse")
			}
			physicsState.velocityX += impulse
			fmt.Printf("New horizontal velocity: %.1f px/s (mouse: %.1f)\n", physicsState.velocityX, mouseImpulseStrength)
		}
	}

	// Handle keyboard events - add impulse to physics
	if *keyUpPressed {
		*keyUpPressed = false // Reset flag
		if scrollableHeight > 0 {
			// Up arrow - add downward impulse
			physicsState.velocityY -= keyboardImpulseStrength
			fmt.Printf("UP ARROW: Added downward impulse, velocity: %.1f px/s (keyboard: %.1f)\n", physicsState.velocityY, keyboardImpulseStrength)
		}
	}

	if *keyDownPressed {
		*keyDownPressed = false // Reset flag
		if scrollableHeight > 0 {
			// Down arrow - add upward impulse
			physicsState.velocityY += keyboardImpulseStrength
			fmt.Printf("DOWN ARROW: Added upward impulse, velocity: %.1f px/s (keyboard: %.1f)\n", physicsState.velocityY, keyboardImpulseStrength)
		}
	}

	if *keyLeftPressed {
		*keyLeftPressed = false // Reset flag
		if scrollableWidth > 0 {
			// Left arrow - add rightward impulse
			physicsState.velocityX -= keyboardImpulseStrength
			fmt.Printf("LEFT ARROW: Added rightward impulse, velocity: %.1f px/s (keyboard: %.1f)\n", physicsState.velocityX, keyboardImpulseStrength)
		}
	}

	if *keyRightPressed {
		*keyRightPressed = false // Reset flag
		if scrollableWidth > 0 {
			// Right arrow - add leftward impulse
			physicsState.velocityX += keyboardImpulseStrength
			fmt.Printf("RIGHT ARROW: Added leftward impulse, velocity: %.1f px/s (keyboard: %.1f)\n", physicsState.velocityX, keyboardImpulseStrength)
		}
	}

	// Physics update loop - apply velocity, friction, speed limits, and bouncing
	deltaTime := float32(1.0 / 60.0) // Assume 60 FPS for consistent physics

	// Apply maximum speed limit
	if physicsState.velocityX > physicsState.maxSpeed {
		physicsState.velocityX = physicsState.maxSpeed
	} else if physicsState.velocityX < -physicsState.maxSpeed {
		physicsState.velocityX = -physicsState.maxSpeed
	}
	if physicsState.velocityY > physicsState.maxSpeed {
		physicsState.velocityY = physicsState.maxSpeed
	} else if physicsState.velocityY < -physicsState.maxSpeed {
		physicsState.velocityY = -physicsState.maxSpeed
	}

	// Update position based on velocity
	physicsState.positionX += physicsState.velocityX * deltaTime
	physicsState.positionY += physicsState.velocityY * deltaTime

	// Convert pixel position to normalized position (0.0 to 1.0)
	var newHorizontalPos, newVerticalPos float32

	if scrollableWidth > 0 {
		newHorizontalPos = physicsState.positionX / float32(scrollableWidth)
	} else {
		newHorizontalPos = 0.0
	}

	if scrollableHeight > 0 {
		newVerticalPos = physicsState.positionY / float32(scrollableHeight)
	} else {
		newVerticalPos = 0.0
	}

	// Handle edge boundaries - stop motion at edges instead of bouncing
	if newHorizontalPos < 0.0 {
		// Hit left edge - stop motion
		physicsState.positionX = 0.0
		physicsState.velocityX = 0.0
		newHorizontalPos = 0.0
		fmt.Printf("HIT LEFT EDGE: Motion stopped\n")
	} else if newHorizontalPos > 1.0 {
		// Hit right edge - stop motion
		physicsState.positionX = float32(scrollableWidth)
		physicsState.velocityX = 0.0
		newHorizontalPos = 1.0
		fmt.Printf("HIT RIGHT EDGE: Motion stopped\n")
	}

	if newVerticalPos < 0.0 {
		// Hit top edge - stop motion
		physicsState.positionY = 0.0
		physicsState.velocityY = 0.0
		newVerticalPos = 0.0
		fmt.Printf("HIT TOP EDGE: Motion stopped\n")
	} else if newVerticalPos > 1.0 {
		// Hit bottom edge - stop motion
		physicsState.positionY = float32(scrollableHeight)
		physicsState.velocityY = 0.0
		newVerticalPos = 1.0
		fmt.Printf("HIT BOTTOM EDGE: Motion stopped\n")
	}

	// Apply friction (Newtonian decay) with momentum mass factor
	// Higher momentum mass = less friction = longer motion
	momentumFriction := physicsState.friction + (1.0-physicsState.friction)*(1.0/MomentumMass)
	physicsState.velocityX *= momentumFriction
	physicsState.velocityY *= momentumFriction

	// Stop very small velocities to prevent infinite tiny movements
	// Increased threshold due to high velocities and low friction
	if math.Abs(float64(physicsState.velocityX)) < 100.0 {
		physicsState.velocityX = 0.0
	}
	if math.Abs(float64(physicsState.velocityY)) < 100.0 {
		physicsState.velocityY = 0.0
	}

	// Check if scrollbar track clicks have changed the position
	// If so, convert the change to physics velocity instead of direct position control
	if horizontalScrollbar.Changed() {
		scrollbarPos := horizontalScrollbar.Position()
		if math.Abs(float64(scrollbarPos-*horizontalPos)) > 0.001 {
			// Scrollbar position changed due to track click - convert to physics
			positionDiff := scrollbarPos - *horizontalPos
			// Convert position difference to velocity (pixels per second)
			// For 250ms animation, we need velocity = distance / 0.25
			velocityChange := (positionDiff * float32(scrollableWidth)) / 0.25
			physicsState.velocityX = velocityChange
			physicsState.positionX = scrollbarPos * float32(scrollableWidth)
			fmt.Printf("TRACK CLICK: Converted to velocity %.1f px/s\n", velocityChange)
		}
	}

	if verticalScrollbar.Changed() {
		scrollbarPos := verticalScrollbar.Position()
		if math.Abs(float64(scrollbarPos-*verticalPos)) > 0.001 {
			// Scrollbar position changed due to track click - convert to physics
			positionDiff := scrollbarPos - *verticalPos
			// Convert position difference to velocity (pixels per second)
			// For 250ms animation, we need velocity = distance / 0.25
			velocityChange := (positionDiff * float32(scrollableHeight)) / 0.25
			physicsState.velocityY = velocityChange
			physicsState.positionY = scrollbarPos * float32(scrollableHeight)
			fmt.Printf("TRACK CLICK: Converted to velocity %.1f px/s\n", velocityChange)
		}
	}

	// Update scrollbar positions and viewport
	*horizontalPos = newHorizontalPos
	*verticalPos = newVerticalPos
	horizontalScrollbar.SetPosition(newHorizontalPos)
	verticalScrollbar.SetPosition(newVerticalPos)

	// Debug output for active physics
	if physicsState.velocityX != 0.0 || physicsState.velocityY != 0.0 {
		fmt.Printf("Physics: pos(%.1f, %.1f) vel(%.1f, %.1f) px/s (momentum friction: %.4f)\n",
			physicsState.positionX, physicsState.positionY,
			physicsState.velocityX, physicsState.velocityY, momentumFriction)
	}

	// Debug output when both velocities are active (potential drift issue)
	if physicsState.velocityX != 0.0 && physicsState.velocityY != 0.0 {
		fmt.Printf("⚠️  DRIFT DETECTED: Both X(%.1f) and Y(%.1f) velocities active!\n",
			physicsState.velocityX, physicsState.velocityY)
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
