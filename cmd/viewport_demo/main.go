package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/gesture"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/key"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"github.com/mleku/fromage"
	"lol.mleku.dev/chk"
)

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
)

// WindowState holds current window dimensions and mouse position
type WindowState struct {
	Width         int     // Window width in Dp
	Height        int     // Window height in Dp
	MouseX        float32 // Mouse X position in pixels
	MouseY        float32 // Mouse Y position in pixels
	MouseInWindow bool    // Whether mouse is currently in the window
}

// ViewportState holds viewport-specific coordinate information
type ViewportState struct {
	ViewportX      float32 // Mouse X position relative to viewport (0-1)
	ViewportY      float32 // Mouse Y position relative to viewport (0-1)
	ViewportWidth  int     // Viewport width in Dp
	ViewportHeight int     // Viewport height in Dp
	ContentX       float32 // Mouse X position in content coordinates
	ContentY       float32 // Mouse Y position in content coordinates
	ScrollX        float32 // Current horizontal scroll position (0-1)
	ScrollY        float32 // Current vertical scroll position (0-1)
}

// Physics configuration - tune these values to adjust motion behavior
var (
	// Mass affects how much impulse is needed to achieve the same velocity
	// Higher mass = more impulse needed = less responsive feeling
	// Lower mass = less impulse needed = more responsive feeling
	// Range: 0.1 (very responsive) to 10.0 (very sluggish)
	// Reduced by 1000x for much lighter viewport
	PhysicsMass float32 = 1000

	// Impulse strength per scroll event (pixels per second)
	// This is the base impulse before mass is applied
	// Quadrupled for 4x faster velocity, then 4x more for 4x scroll distance, then 100x more
	BaseImpulseStrength float32 = 610000.0

	// Mouse scroll energy multiplier - makes mouse scrolling more powerful
	// Higher values = more energy per mouse scroll event
	// Range: 1.0 (same as keyboard) to 10.0 (very powerful)
	MouseScrollMultiplier float32 = 0.0002

	// Momentum mass - affects how long motion continues after input stops
	// Higher values = longer momentum, more continuous motion
	// Lower values = motion stops more quickly
	// Range: 0.1 (stops quickly) to 10.0 (very long momentum)
	MomentumMass float32 = 5 // Much lower for very quick stopping
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
func clipViewport(gtx layout.Context, th *fromage.Theme, contentSize int, horizontalPos, verticalPos float32, viewportWidth, viewportHeight int, modalStack *fromage.ModalStack, windowState *WindowState, viewportState *ViewportState, lastScrollEvent *float32) layout.Dimensions {
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

	// Example: Use viewport state for overlay widgets
	// This shows how overlay widgets can access viewport-relative coordinates
	if viewportState.ViewportX > 0 && viewportState.ViewportY > 0 {
		fmt.Printf("🎯 OVERLAY WIDGET: Mouse at viewport (%.3f, %.3f), content (%.1f, %.1f), scroll (%.3f, %.3f)\n",
			viewportState.ViewportX, viewportState.ViewportY,
			viewportState.ContentX, viewportState.ContentY,
			viewportState.ScrollX, viewportState.ScrollY)
	}

	// Draw the test button in the center of the canvas
	centerX := contentSize / 2
	centerY := contentSize / 2
	contentArea := op.Offset(image.Pt(centerX, centerY)).Push(gtx.Ops)

	// Draw some content in the viewport instead of the button
	th.CenteredColumn().
		Rigid(func(g C) D {
			return th.NewLabel().Text("Viewport Content").Layout(g)
		}).
		Layout(gtx)

	contentArea.Pop()
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
	// Removed ButtonLayout - no longer needed
	modalStack := th.NewModalStack()

	// Global menu removed - no right-click context menu

	// Variables for event handling
	var horizontalPos float32 = 0.0
	var verticalPos float32 = 0.0
	var lastScrollEvent float32 = 0.0

	// Create event handler for consistent event processing
	eventHandler := fromage.NewEventHandler(func(eventDesc string) {
		fmt.Printf("📝 EVENT: %s\n", eventDesc)
	}).SetOnScroll(func(scrollY float32) {
		// Handle scroll events with correct direction mapping
		fmt.Printf("🖱️ SCROLL EVENT: Y=%.1f\n", scrollY)
		// Store scroll event for processing in mainUI
		lastScrollEvent = scrollY
	}).SetOnClick(func(e pointer.Event) {
		fmt.Printf("🖱️ CLICK EVENT: Position=(%.1f,%.1f)\n", e.Position.X, e.Position.Y)
	})

	// Pointer tag for top-level events
	pointerTag := &struct{}{}
	// Gesture tag for gesture events
	gestureTag := &struct{}{}

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
	physicsState.friction = 0.89 // Much higher friction for faster deceleration
	physicsState.maxSpeed = 50000.0
	physicsState.bounceRetention = 0.5

	// Keyboard event flags for direct scrolling
	var keyUpPressed, keyDownPressed, keyLeftPressed, keyRightPressed bool
	var pageUpPressed, pageDownPressed bool
	var homePressed, endPressed bool

	// Keyboard state tracking for acceleration with 250ms intervals
	var keyStates struct {
		upPressed, downPressed, leftPressed, rightPressed                     bool
		upLastTime, downLastTime, leftLastTime, rightLastTime                 int64
		upAcceleration, downAcceleration, leftAcceleration, rightAcceleration float32
	}

	// Window state tracking
	var windowState WindowState
	var viewportState ViewportState

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

				// Update window state with current dimensions and mouse position
				windowState.Width = gtx.Dp(unit.Dp(gtx.Constraints.Max.X))
				windowState.Height = gtx.Dp(unit.Dp(gtx.Constraints.Max.Y))

				// Register for pointer events to capture mouse position
				pointerArea := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
				event.Op(gtx.Ops, &windowState) // Use windowState as the event tag

				// Use EventHandler for scroll and click events
				eventHandler.AddToOps(gtx.Ops)
				eventHandler.ProcessEvents(gtx)

				// Process pointer events to update mouse position
				for {
					event, ok := gtx.Event(pointer.Filter{
						Kinds: pointer.Move | pointer.Enter | pointer.Leave,
					})
					if !ok {
						break
					}
					if pointerEvent, ok := event.(pointer.Event); ok {
						windowState.MouseX = pointerEvent.Position.X
						windowState.MouseY = pointerEvent.Position.Y
						windowState.MouseInWindow = pointerEvent.Kind != pointer.Leave
					}
				}
				pointerArea.Pop()

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
						key.Filter{Name: key.NamePageUp},
						key.Filter{Name: key.NamePageDown},
						key.Filter{Name: key.NameHome},
						key.Filter{Name: key.NameEnd},
					)
					if !ok {
						break
					}
					switch event := event.(type) {
					case key.Event:
						fmt.Printf("Key event: %s, state: %s\n", event.Name, event.State)
						now := time.Now().UnixNano()

						switch event.Name {
						case key.NameUpArrow:
							if event.State == key.Press && !keyStates.upPressed {
								// Key pressed down - start tracking
								keyStates.upPressed = true
								keyStates.upLastTime = now
								keyStates.upAcceleration = 1.0
								keyUpPressed = true
								fmt.Println("UP ARROW PRESSED - Starting scroll")
							} else if event.State == key.Release && keyStates.upPressed {
								// Key released - stop tracking
								keyStates.upPressed = false
								keyStates.upAcceleration = 1.0
								keyUpPressed = false
								fmt.Println("UP ARROW RELEASED - Stopping scroll")
							}
						case key.NameDownArrow:
							if event.State == key.Press && !keyStates.downPressed {
								// Key pressed down - start tracking
								keyStates.downPressed = true
								keyStates.downLastTime = now
								keyStates.downAcceleration = 1.0
								keyDownPressed = true
								fmt.Println("DOWN ARROW PRESSED - Starting scroll")
							} else if event.State == key.Release && keyStates.downPressed {
								// Key released - stop tracking
								keyStates.downPressed = false
								keyStates.downAcceleration = 1.0
								keyDownPressed = false
								fmt.Println("DOWN ARROW RELEASED - Stopping scroll")
							}
						case key.NameLeftArrow:
							if event.State == key.Press && !keyStates.leftPressed {
								// Key pressed down - start tracking
								keyStates.leftPressed = true
								keyStates.leftLastTime = now
								keyStates.leftAcceleration = 1.0
								keyLeftPressed = true
								fmt.Println("LEFT ARROW PRESSED - Starting scroll")
							} else if event.State == key.Release && keyStates.leftPressed {
								// Key released - stop tracking
								keyStates.leftPressed = false
								keyStates.leftAcceleration = 1.0
								keyLeftPressed = false
								fmt.Println("LEFT ARROW RELEASED - Stopping scroll")
							}
						case key.NameRightArrow:
							if event.State == key.Press && !keyStates.rightPressed {
								// Key pressed down - start tracking
								keyStates.rightPressed = true
								keyStates.rightLastTime = now
								keyStates.rightAcceleration = 1.0
								keyRightPressed = true
								fmt.Println("RIGHT ARROW PRESSED - Starting scroll")
							} else if event.State == key.Release && keyStates.rightPressed {
								// Key released - stop tracking
								keyStates.rightPressed = false
								keyStates.rightAcceleration = 1.0
								keyRightPressed = false
								fmt.Println("RIGHT ARROW RELEASED - Stopping scroll")
							}
						case key.NamePageUp:
							if event.State == key.Press {
								// Page Up - scroll up one screenful smoothly in 250ms
								pageUpPressed = true
								fmt.Println("PAGE UP PRESSED - Scrolling up one screenful")
							}
						case key.NamePageDown:
							if event.State == key.Press {
								// Page Down - scroll down one screenful smoothly in 250ms
								pageDownPressed = true
								fmt.Println("PAGE DOWN PRESSED - Scrolling down one screenful")
							}
						case key.NameHome:
							if event.State == key.Press {
								// Home - scroll left one screenful smoothly in 250ms
								homePressed = true
								fmt.Println("HOME PRESSED - Scrolling left one screenful")
							}
						case key.NameEnd:
							if event.State == key.Press {
								// End - scroll right one screenful smoothly in 250ms
								endPressed = true
								fmt.Println("END PRESSED - Scrolling right one screenful")
							}
						}
					}
				}
				area.Pop()

				mainUI(gtx, th, window, horizontalScrollbar, verticalScrollbar, modalStack, nil, pointerTag, gestureTag, &horizontalPos, &verticalPos, &keyUpPressed, &keyDownPressed, &keyLeftPressed, &keyRightPressed, &pageUpPressed, &pageDownPressed, &homePressed, &endPressed, &physicsState, &keyStates, &windowState, &viewportState, &lastScrollEvent)

				// Note: Scrollbar positions are now controlled by physics system
				// No need to update from scrollbar changes since physics sets them directly
				e.Frame(gtx.Ops)

				// Invalidate the window to continue animation if physics is active
				// This ensures smooth animation continues until inertia decays
				if physicsState.velocityX != 0.0 || physicsState.velocityY != 0.0 {
					w.Invalidate()
				}
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme, window *fromage.Window,
	horizontalScrollbar, verticalScrollbar *fromage.Scrollbar,
	modalStack *fromage.ModalStack,
	scrollGesture *gesture.Scroll, pointerTag, gestureTag interface{},
	horizontalPos, verticalPos *float32, keyUpPressed, keyDownPressed, keyLeftPressed, keyRightPressed, pageUpPressed, pageDownPressed, homePressed, endPressed *bool,
	physicsState *struct {
		velocityX, velocityY float32
		positionX, positionY float32
		friction             float32
		maxSpeed             float32
		bounceRetention      float32
	},
	keyStates *struct {
		upPressed, downPressed, leftPressed, rightPressed                     bool
		upLastTime, downLastTime, leftLastTime, rightLastTime                 int64
		upAcceleration, downAcceleration, leftAcceleration, rightAcceleration float32
	},
	windowState *WindowState,
	viewportState *ViewportState,
	lastScrollEvent *float32) {

	// Fill background with theme background color
	th.FillBackground(nil).Layout(gtx)

	// Log window state for debugging
	fmt.Printf("🪟 WINDOW STATE: Size=(%d, %d) Dp, Mouse=(%.1f, %.1f) px, InWindow=%t\n",
		windowState.Width, windowState.Height, windowState.MouseX, windowState.MouseY, windowState.MouseInWindow)

	// Log viewport state for debugging
	fmt.Printf("📺 VIEWPORT STATE: Size=(%d, %d) Dp, Viewport=(%.3f, %.3f), Content=(%.1f, %.1f), Scroll=(%.3f, %.3f)\n",
		viewportState.ViewportWidth, viewportState.ViewportHeight,
		viewportState.ViewportX, viewportState.ViewportY,
		viewportState.ContentX, viewportState.ContentY,
		viewportState.ScrollX, viewportState.ScrollY)

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

	// Update viewport state with current information
	viewportState.ViewportWidth = gtx.Dp(unit.Dp(contentAreaWidth))
	viewportState.ViewportHeight = gtx.Dp(unit.Dp(contentAreaHeight))
	viewportState.ScrollX = *horizontalPos
	viewportState.ScrollY = *verticalPos

	// Calculate mouse position relative to viewport (0-1)
	if contentAreaWidth > 0 && contentAreaHeight > 0 {
		viewportState.ViewportX = windowState.MouseX / float32(contentAreaWidth)
		viewportState.ViewportY = windowState.MouseY / float32(contentAreaHeight)

		// Calculate mouse position in content coordinates
		viewportState.ContentX = windowState.MouseX + (*horizontalPos * float32(contentSize-contentAreaWidth))
		viewportState.ContentY = windowState.MouseY + (*verticalPos * float32(contentSize-contentAreaHeight))
	}

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

							// Scroll events are now handled directly via pointer events

							// Clipped viewport for the red corner outline widget
							dims := clipViewport(g, th, contentSize, *horizontalPos, *verticalPos, contentAreaWidth, contentAreaHeight, modalStack, windowState, viewportState, lastScrollEvent)

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

	// Process scroll events from EventHandler
	if *lastScrollEvent != 0.0 {
		fmt.Printf("🖱️ PROCESSING SCROLL: Y=%.1f\n", *lastScrollEvent)

		// Apply scroll to physics state with correct direction mapping
		if scrollableHeight > 0 {
			// Convert scroll Y to impulse with correct direction
			// Positive scroll Y means scroll up (content moves up), negative means scroll down
			mouseImpulseStrength := (BaseImpulseStrength * MouseScrollMultiplier) / (PhysicsMass / 10000.0)
			impulse := *lastScrollEvent * mouseImpulseStrength // No negation - direct mapping
			physicsState.velocityY += impulse
			fmt.Printf("Applied scroll impulse: %.1f, new velocity: %.1f\n", impulse, physicsState.velocityY)
		}

		// Reset scroll event
		*lastScrollEvent = 0.0
	}

	keyboardImpulseStrength := BaseImpulseStrength / PhysicsMass

	// Handle keyboard events with 250ms intervals and 10% acceleration
	now := time.Now().UnixNano()

	// Helper function to check if 250ms have passed since last key press
	checkKeyInterval := func(lastTime int64) bool {
		elapsedMs := (now - lastTime) / 1e6 // Convert nanoseconds to milliseconds
		return elapsedMs >= 250
	}

	// Handle up arrow - apply impulse every 250ms with acceleration
	if *keyUpPressed && keyStates.upPressed {
		if checkKeyInterval(keyStates.upLastTime) {
			// Apply acceleration (10% increase each time)
			keyStates.upAcceleration *= 1.1
			keyStates.upLastTime = now

			if scrollableHeight > 0 {
				// Up arrow - add upward impulse (negative Y velocity)
				impulse := keyboardImpulseStrength * keyStates.upAcceleration
				physicsState.velocityY -= impulse
				fmt.Printf("UP ARROW: Applied upward impulse %.1f, velocity: %.1f px/s (accel: %.2fx)\n",
					impulse, physicsState.velocityY, keyStates.upAcceleration)
			}
		}
	}

	// Handle down arrow - apply impulse every 250ms with acceleration
	if *keyDownPressed && keyStates.downPressed {
		if checkKeyInterval(keyStates.downLastTime) {
			// Apply acceleration (10% increase each time)
			keyStates.downAcceleration *= 1.1
			keyStates.downLastTime = now

			if scrollableHeight > 0 {
				// Down arrow - add downward impulse (positive Y velocity)
				impulse := keyboardImpulseStrength * keyStates.downAcceleration
				physicsState.velocityY += impulse
				fmt.Printf("DOWN ARROW: Applied downward impulse %.1f, velocity: %.1f px/s (accel: %.2fx)\n",
					impulse, physicsState.velocityY, keyStates.downAcceleration)
			}
		}
	}

	// Handle left arrow - apply impulse every 250ms with acceleration
	if *keyLeftPressed && keyStates.leftPressed {
		if checkKeyInterval(keyStates.leftLastTime) {
			// Apply acceleration (10% increase each time)
			keyStates.leftAcceleration *= 1.1
			keyStates.leftLastTime = now

			if scrollableWidth > 0 {
				// Left arrow - add leftward impulse (negative X velocity)
				impulse := keyboardImpulseStrength * keyStates.leftAcceleration
				physicsState.velocityX -= impulse
				fmt.Printf("LEFT ARROW: Applied leftward impulse %.1f, velocity: %.1f px/s (accel: %.2fx)\n",
					impulse, physicsState.velocityX, keyStates.leftAcceleration)
			}
		}
	}

	// Handle right arrow - apply impulse every 250ms with acceleration
	if *keyRightPressed && keyStates.rightPressed {
		if checkKeyInterval(keyStates.rightLastTime) {
			// Apply acceleration (10% increase each time)
			keyStates.rightAcceleration *= 1.1
			keyStates.rightLastTime = now

			if scrollableWidth > 0 {
				// Right arrow - add rightward impulse (positive X velocity)
				impulse := keyboardImpulseStrength * keyStates.rightAcceleration
				physicsState.velocityX += impulse
				fmt.Printf("RIGHT ARROW: Applied rightward impulse %.1f, velocity: %.1f px/s (accel: %.2fx)\n",
					impulse, physicsState.velocityX, keyStates.rightAcceleration)
			}
		}
	}

	// Handle Page Up - scroll up one screenful smoothly in 250ms
	if *pageUpPressed {
		*pageUpPressed = false // Reset flag
		if scrollableHeight > 0 {
			// Calculate velocity needed to scroll one screenful in 250ms
			// One screenful = contentAreaHeight pixels
			// Velocity = distance / time = contentAreaHeight / 0.25 seconds
			pageScrollVelocity := float32(contentAreaHeight) / 0.25
			physicsState.velocityY = -pageScrollVelocity // Negative for upward scroll
			fmt.Printf("PAGE UP: Set velocity to scroll one screenful up: %.1f px/s\n", pageScrollVelocity)
		}
	}

	// Handle Page Down - scroll down one screenful smoothly in 250ms
	if *pageDownPressed {
		*pageDownPressed = false // Reset flag
		if scrollableHeight > 0 {
			// Calculate velocity needed to scroll one screenful in 250ms
			// One screenful = contentAreaHeight pixels
			// Velocity = distance / time = contentAreaHeight / 0.25 seconds
			pageScrollVelocity := float32(contentAreaHeight) / 0.25
			physicsState.velocityY = pageScrollVelocity // Positive for downward scroll
			fmt.Printf("PAGE DOWN: Set velocity to scroll one screenful down: %.1f px/s\n", pageScrollVelocity)
		}
	}

	// Handle Home - scroll left one screenful smoothly in 250ms
	if *homePressed {
		*homePressed = false // Reset flag
		if scrollableWidth > 0 {
			// Calculate velocity needed to scroll one screenful in 250ms
			// One screenful = contentAreaWidth pixels
			// Velocity = distance / time = contentAreaWidth / 0.25 seconds
			pageScrollVelocity := float32(contentAreaWidth) / 0.25
			physicsState.velocityX = -pageScrollVelocity // Negative for leftward scroll
			fmt.Printf("HOME: Set velocity to scroll one screenful left: %.1f px/s\n", pageScrollVelocity)
		}
	}

	// Handle End - scroll right one screenful smoothly in 250ms
	if *endPressed {
		*endPressed = false // Reset flag
		if scrollableWidth > 0 {
			// Calculate velocity needed to scroll one screenful in 250ms
			// One screenful = contentAreaWidth pixels
			// Velocity = distance / time = contentAreaWidth / 0.25 seconds
			pageScrollVelocity := float32(contentAreaWidth) / 0.25
			physicsState.velocityX = pageScrollVelocity // Positive for rightward scroll
			fmt.Printf("END: Set velocity to scroll one screenful right: %.1f px/s\n", pageScrollVelocity)
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
	// Much higher threshold for faster stopping
	if math.Abs(float64(physicsState.velocityX)) < 500.0 {
		physicsState.velocityX = 0.0
	}
	if math.Abs(float64(physicsState.velocityY)) < 500.0 {
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

	// Global menu removed - no right-click context menu handling

	// Layout the modal stack on top of everything
	if !modalStack.IsEmpty() {
		modalStack.Layout(gtx)
	}
}
