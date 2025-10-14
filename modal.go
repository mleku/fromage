package fromage

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

// ModalStack manages a stack of modal dialogs
type ModalStack struct {
	theme     *Theme
	modals    []*Modal
	scrimDark float32 // 0.0 = transparent, 1.0 = fully opaque
}

// Modal represents a single modal dialog
type Modal struct {
	theme             *Theme
	content           W
	scrimClickable    *widget.Clickable
	onClose           func()
	blocking          bool
	animationProgress float32   // 0.0 = invisible, 1.0 = fully visible
	animationStart    time.Time // When the animation started
	isAnimating       bool      // Whether animation is in progress
	animationStarted  bool      // Whether animation has ever been started
	isFadingOut       bool      // Whether we're fading out (true) or fading in (false)
	slideIn           bool      // Whether to use slide-in animation instead of fade
}

// NewModalStack creates a new modal stack
func (t *Theme) NewModalStack() *ModalStack {
	return &ModalStack{
		theme:     t,
		modals:    make([]*Modal, 0),
		scrimDark: 0.5, // Default 50% opacity
	}
}

// ScrimDarkness sets the darkness of the scrim (0.0 = transparent, 1.0 = fully opaque)
func (ms *ModalStack) ScrimDarkness(darkness float32) *ModalStack {
	ms.scrimDark = darkness
	return ms
}

// Push adds a new modal to the stack
func (ms *ModalStack) Push(content W, onClose func()) {
	ms.PushWithSlide(content, onClose, false)
}

// PushWithSlide adds a new modal to the stack with slide-in option
func (ms *ModalStack) PushWithSlide(content W, onClose func(), slideIn bool) {
	modal := &Modal{
		theme:             ms.theme,
		content:           content,
		scrimClickable:    ms.theme.Pool.GetClickable(),
		onClose:           onClose,
		blocking:          true,
		animationProgress: 0.0,
		animationStart:    time.Time{},
		isAnimating:       false,
		animationStarted:  false,
		isFadingOut:       false,
		slideIn:           slideIn,
	}
	ms.modals = append(ms.modals, modal)
}

// Pop removes the top modal from the stack
func (ms *ModalStack) Pop() {
	if len(ms.modals) == 0 {
		return
	}

	// Start fade-out animation instead of immediate removal
	topModal := ms.modals[len(ms.modals)-1]
	topModal.startFadeOut()
}

// Clear removes all modals from the stack
func (ms *ModalStack) Clear() {
	ms.modals = ms.modals[:0]
}

// IsEmpty returns true if there are no modals in the stack
func (ms *ModalStack) IsEmpty() bool {
	return len(ms.modals) == 0
}

// Count returns the number of modals in the stack
func (ms *ModalStack) Count() int {
	return len(ms.modals)
}

// Layout renders the modal stack
func (ms *ModalStack) Layout(gtx C) D {
	if len(ms.modals) == 0 {
		return D{}
	}

	// Remove completed fade-outs first
	ms.removeCompletedFadeOuts()

	// Layout all modals in order (bottom to top)
	for i, modal := range ms.modals {
		// Start animation if not already started
		if !modal.animationStarted {
			modal.startAnimation(gtx.Now)
		}

		// Update animation progress
		modal.updateAnimation(gtx)

		// Create scrim color with the specified darkness and animation progress
		// Clamp animation progress to prevent flickering
		progress := modal.animationProgress
		if progress > 1.0 {
			progress = 1.0
		}
		if progress < 0.0 {
			progress = 0.0
		}
		scrimAlpha := uint8(255 * ms.scrimDark * progress)
		scrimColor := color.NRGBA{
			R: 0,
			G: 0,
			B: 0,
			A: scrimAlpha,
		}

		// Create a new context for this modal layer
		modalGtx := gtx

		// Use a stack layout to properly position the scrim and content
		modalResult := layout.Stack{}.Layout(modalGtx,
			// First layer: Fill entire screen with scrim and handle clicks
			layout.Expanded(func(gtx C) D {
				// Fill with scrim color
				paint.Fill(gtx.Ops, scrimColor)

				// Layout clickable area over the entire scrim
				modal.scrimClickable.Layout(gtx, func(gtx C) D {
					return layout.Dimensions{Size: gtx.Constraints.Max}
				})

				// Handle scrim clicks (click outside to close)
				if modal.scrimClickable.Clicked(gtx) {
					if modal.onClose != nil {
						modal.onClose()
					}
				}

				return layout.Dimensions{Size: gtx.Constraints.Max}
			}),
			// Second layer: Layout the modal content at the top
			layout.Stacked(func(gtx C) D {
				if modal.slideIn {
					// For slide-in animation, we need to transform the content position
					// Apply opacity to the content (same as fade animation)
					defer paint.PushOpacity(gtx.Ops, progress).Pop()

					// Calculate the transform offset for slide animation
					// Use a fixed slide distance for consistent animation
					slideDistance := 300 // pixels

					var offsetY int
					if modal.isFadingOut {
						// During slide-out, move content up off-screen
						offsetY = -int(float32(slideDistance) * (1.0 - progress))
					} else {
						// During slide-in, move content from above screen to final position
						offsetY = -int(float32(slideDistance) * (1.0 - progress))
					}

					// Apply transform to move the content vertically
					defer op.Offset(image.Point{X: 0, Y: offsetY}).Push(gtx.Ops).Pop()

					// Layout the modal content at the top, full width
					return layout.NW.Layout(gtx, func(gtx C) D {
						// Constrain to full width
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return modal.content(gtx)
					})
				} else {
					// Apply fade-in opacity to the content
					defer paint.PushOpacity(gtx.Ops, progress).Pop()

					// Layout the modal content at the top, full width
					return layout.NW.Layout(gtx, func(gtx C) D {
						// Constrain to full width
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return modal.content(gtx)
					})
				}
			}),
		)

		// If this is the top modal, we're done
		if i == len(ms.modals)-1 {
			return modalResult
		}
	}

	return D{}
}

// removeCompletedFadeOuts removes modals that have completed their fade-out animation
func (ms *ModalStack) removeCompletedFadeOuts() {
	// Remove modals from the end that have completed fade-out
	for len(ms.modals) > 0 {
		lastModal := ms.modals[len(ms.modals)-1]
		if lastModal.isFadingOut && !lastModal.isAnimating && lastModal.animationProgress <= 0.0 {
			// This modal has completed fade-out, remove it
			ms.modals = ms.modals[:len(ms.modals)-1]
		} else {
			// Stop at the first modal that hasn't completed fade-out
			break
		}
	}
}

// Modal-specific methods

// Blocking sets whether this modal blocks interaction with content behind it
func (m *Modal) Blocking(blocking bool) *Modal {
	m.blocking = blocking
	return m
}

// startAnimation begins the fade-in animation
func (m *Modal) startAnimation(now time.Time) {
	m.animationStart = now
	m.isAnimating = true
	m.animationStarted = true
	m.isFadingOut = false
}

// startFadeOut begins the fade-out animation
func (m *Modal) startFadeOut() {
	m.animationStart = time.Now()
	m.isAnimating = true
	m.isFadingOut = true
}

// updateAnimation updates the animation progress
func (m *Modal) updateAnimation(g C) {
	if !m.isAnimating {
		return
	}

	const animationDuration = 250 * time.Millisecond
	elapsed := g.Now.Sub(m.animationStart)

	if elapsed >= animationDuration {
		if m.isFadingOut {
			m.animationProgress = 0.0
			m.isAnimating = false
		} else {
			m.animationProgress = 1.0
			m.isAnimating = false
		}
		return
	}

	progress := float32(elapsed) / float32(animationDuration)
	if progress > 1.0 {
		progress = 1.0
	}

	// Use ease-out curve: 1 - (1-t)^2
	easedProgress := 1.0 - (1.0-progress)*(1.0-progress)

	if m.isFadingOut {
		m.animationProgress = 1.0 - easedProgress
	} else {
		m.animationProgress = easedProgress
	}

	if m.animationProgress > 1.0 {
		m.animationProgress = 1.0
	}
	if m.animationProgress < 0.0 {
		m.animationProgress = 0.0
	}

	g.Execute(op.InvalidateCmd{})
}
