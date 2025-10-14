package fromage

import (
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
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
	// Theme reference
	theme *Theme
	// Content widget
	content W
	// Clickable for scrim (background) clicks
	scrimClickable *widget.Clickable
	// Close handler
	onClose func()
	// Whether this modal blocks interaction with content behind it
	blocking bool
	// Animation state
	animationProgress float32   // 0.0 = invisible, 1.0 = fully visible
	animationStart    time.Time // When the animation started
	isAnimating       bool      // Whether animation is in progress
	animationStarted  bool      // Whether animation has ever been started
	isFadingOut       bool      // Whether we're fading out (true) or fading in (false)
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
	if darkness < 0.0 {
		darkness = 0.0
	}
	if darkness > 1.0 {
		darkness = 1.0
	}
	ms.scrimDark = darkness
	return ms
}

// Push adds a new modal to the stack
func (ms *ModalStack) Push(content W, onClose func()) *Modal {
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
	}
	ms.modals = append(ms.modals, modal)
	return modal
}

// Pop removes the top modal from the stack
func (ms *ModalStack) Pop() {
	if len(ms.modals) > 0 {
		// Start fade-out animation instead of immediately removing
		topModal := ms.modals[len(ms.modals)-1]
		if !topModal.isFadingOut {
			topModal.startFadeOut()
		}
	}
}

// Clear removes all modals from the stack
func (ms *ModalStack) Clear() {
	ms.modals = ms.modals[:0]
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
	if ms.IsEmpty() {
		return D{}
	}

	// Remove modals that have completed fade-out animation
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

		// Fill the entire screen with scrim
		paint.Fill(modalGtx.Ops, scrimColor)

		// Handle scrim clicks (click outside to close)
		if modal.scrimClickable.Clicked(modalGtx) {
			if modal.onClose != nil {
				modal.onClose()
			}
		}

		// Layout the scrim clickable area
		modal.scrimClickable.Layout(modalGtx, func(gtx C) D {
			// Fill the entire area
			return layout.Dimensions{Size: gtx.Constraints.Max}
		})

		// Layout the modal content in the center with fade-in animation
		modalGtx.Constraints = layout.Exact(gtx.Constraints.Max)
		modalContent := layout.Center.Layout(modalGtx, func(gtx C) D {
			// Create a new context that doesn't propagate clicks to the scrim
			contentGtx := gtx

			// Create a clip area for the content to prevent clicks from going through
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()

			// Apply fade-in opacity to the content
			// Use the same clamped progress value
			defer paint.PushOpacity(gtx.Ops, progress).Pop()

			// Layout the modal content
			return modal.content(contentGtx)
		})

		// If this is the top modal, we're done
		if i == len(ms.modals)-1 {
			return modalContent
		}
	}

	return D{}
}

// Modal-specific methods

// Blocking sets whether this modal blocks interaction with content behind it
func (m *Modal) Blocking(blocking bool) *Modal {
	m.blocking = blocking
	return m
}

// OnClose sets the close handler for this modal
func (m *Modal) OnClose(handler func()) *Modal {
	m.onClose = handler
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

// updateAnimation updates the animation progress based on elapsed time
func (m *Modal) updateAnimation(g C) {
	if !m.isAnimating {
		return
	}

	const animationDuration = 250 * time.Millisecond
	elapsed := g.Now.Sub(m.animationStart)

	if elapsed >= animationDuration {
		// Animation complete
		if m.isFadingOut {
			// Fade-out complete - modal should be removed
			m.animationProgress = 0.0
			m.isAnimating = false
			// The modal will be removed from the stack in the layout method
		} else {
			// Fade-in complete - ensure we're at 100% opacity
			m.animationProgress = 1.0
			m.isAnimating = false
		}
		return
	}

	// Calculate progress (0.0 to 1.0)
	progress := float32(elapsed) / float32(animationDuration)

	// Clamp progress to prevent floating point issues
	if progress > 1.0 {
		progress = 1.0
	}

	// Apply ease-out timing function (similar to CSS ease-out)
	// This creates a smooth deceleration effect
	easedProgress := 1.0 - (1.0-progress)*(1.0-progress)

	if m.isFadingOut {
		// Fade-out: go from 1.0 to 0.0
		m.animationProgress = 1.0 - easedProgress
	} else {
		// Fade-in: go from 0.0 to 1.0
		m.animationProgress = easedProgress
	}

	// Ensure final progress is clamped
	if m.animationProgress > 1.0 {
		m.animationProgress = 1.0
	}
	if m.animationProgress < 0.0 {
		m.animationProgress = 0.0
	}

	// Request invalidation to continue animation
	g.Execute(op.InvalidateCmd{})
}
