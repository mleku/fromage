package fromage

import (
	"context"
	"testing"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
)

func TestModalStack(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		context.TODO(),
		NewColors,
		text.NewShaper(text.WithCollection(nil)),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Create a modal stack
	modalStack := th.NewModalStack()

	// Test initial state
	if !modalStack.IsEmpty() {
		t.Error("Modal stack should be empty initially")
	}

	if modalStack.Count() != 0 {
		t.Error("Modal stack count should be 0 initially")
	}

	// Test scrim darkness
	modalStack.ScrimDarkness(0.8)
	if modalStack.scrimDark != 0.8 {
		t.Error("Scrim darkness should be set to 0.8")
	}

	// Test scrim darkness bounds
	modalStack.ScrimDarkness(-0.1)
	if modalStack.scrimDark != 0.0 {
		t.Error("Scrim darkness should be clamped to 0.0")
	}

	modalStack.ScrimDarkness(1.1)
	if modalStack.scrimDark != 1.0 {
		t.Error("Scrim darkness should be clamped to 1.0")
	}

	// Create a simple modal content
	modalContent := func(g C) D {
		return layout.Dimensions{Size: g.Constraints.Max}
	}

	// Test pushing a modal
	modal := modalStack.Push(modalContent, func() {})
	if modal == nil {
		t.Error("Push should return a modal")
	}

	if modalStack.IsEmpty() {
		t.Error("Modal stack should not be empty after pushing")
	}

	if modalStack.Count() != 1 {
		t.Error("Modal stack count should be 1 after pushing")
	}

	// Test animation fields are initialized correctly
	if modal.animationProgress != 0.0 {
		t.Error("Modal should start with 0.0 animation progress")
	}

	if modal.isAnimating {
		t.Error("Modal should not be animating initially")
	}

	if modal.animationStarted {
		t.Error("Modal should not have animation started initially")
	}

	// Test pushing another modal
	modal2 := modalStack.Push(modalContent, func() {})
	if modal2 == nil {
		t.Error("Push should return a modal")
	}

	if modalStack.Count() != 2 {
		t.Error("Modal stack count should be 2 after pushing second modal")
	}

	// Test popping a modal (starts fade-out animation)
	modalStack.Pop()
	// Modal should still be in stack until fade-out completes
	if modalStack.Count() != 2 {
		t.Error("Modal stack count should still be 2 after starting fade-out")
	}

	// Check that the top modal is now fading out
	topModal := modalStack.modals[len(modalStack.modals)-1]
	if !topModal.isFadingOut {
		t.Error("Top modal should be fading out after Pop()")
	}

	// Test clearing all modals
	modalStack.Clear()
	if !modalStack.IsEmpty() {
		t.Error("Modal stack should be empty after clearing")
	}

	if modalStack.Count() != 0 {
		t.Error("Modal stack count should be 0 after clearing")
	}
}

func TestModal(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		context.TODO(),
		NewColors,
		text.NewShaper(text.WithCollection(nil)),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Create a simple modal content
	modalContent := func(g C) D {
		return layout.Dimensions{Size: g.Constraints.Max}
	}

	// Create a modal
	modal := th.NewModal(modalContent, func() {})
	if modal == nil {
		t.Error("NewModal should return a modal")
	}

	if modal.theme != th {
		t.Error("Modal should have the correct theme reference")
	}

	if modal.content == nil {
		t.Error("Modal should have content")
	}

	if modal.scrimClickable == nil {
		t.Error("Modal should have a scrim clickable")
	}

	if !modal.blocking {
		t.Error("Modal should be blocking by default")
	}

	if modal.animationProgress != 0.0 {
		t.Error("Modal should start with 0.0 animation progress")
	}

	if modal.isAnimating {
		t.Error("Modal should not be animating initially")
	}

	// Test setting blocking state
	modal.Blocking(false)
	if modal.blocking {
		t.Error("Modal should not be blocking after setting to false")
	}

	// Test setting close handler
	closeCalled := false
	modal.OnClose(func() {
		closeCalled = true
	})
	if modal.onClose == nil {
		t.Error("Modal should have a close handler")
	}

	// Test calling close handler
	modal.onClose()
	if !closeCalled {
		t.Error("Close handler should be called")
	}
}

func TestModalAnimation(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		context.TODO(),
		NewColors,
		text.NewShaper(text.WithCollection(nil)),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Create a simple modal content
	modalContent := func(g C) D {
		return layout.Dimensions{Size: g.Constraints.Max}
	}

	// Create a modal
	modal := th.NewModal(modalContent, func() {})

	// Test initial animation state
	if modal.animationProgress != 0.0 {
		t.Error("Modal should start with 0.0 animation progress")
	}

	if modal.isAnimating {
		t.Error("Modal should not be animating initially")
	}

	if modal.animationStarted {
		t.Error("Modal should not have animation started initially")
	}

	// Test starting animation
	now := time.Now()
	modal.startAnimation(now)

	if !modal.isAnimating {
		t.Error("Modal should be animating after startAnimation")
	}

	if modal.animationStart != now {
		t.Error("Modal should have correct animation start time")
	}

	if !modal.animationStarted {
		t.Error("Modal should have animation started flag set")
	}

	// Test animation progress calculation
	// Create a mock context with time
	var ops op.Ops
	gtx := app.NewContext(&ops, app.FrameEvent{
		Now: now.Add(125 * time.Millisecond), // Halfway through 250ms animation
	})

	modal.updateAnimation(gtx)

	// Should be approximately 0.75 progress (ease-out curve at halfway point)
	// Ease-out: 1 - (1-0.5)^2 = 1 - 0.25 = 0.75
	if modal.animationProgress < 0.7 || modal.animationProgress > 0.8 {
		t.Errorf("Animation progress should be around 0.75 (ease-out), got %f", modal.animationProgress)
	}

	// Test animation completion
	gtx = app.NewContext(&ops, app.FrameEvent{
		Now: now.Add(300 * time.Millisecond), // Past animation duration
	})

	modal.updateAnimation(gtx)

	if modal.animationProgress != 1.0 {
		t.Error("Animation progress should be 1.0 when complete")
	}

	if modal.isAnimating {
		t.Error("Modal should not be animating when complete")
	}
}

func TestModalFadeOut(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		context.TODO(),
		NewColors,
		text.NewShaper(text.WithCollection(nil)),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Create a simple modal content
	modalContent := func(g C) D {
		return layout.Dimensions{Size: g.Constraints.Max}
	}

	// Create a modal
	modal := th.NewModal(modalContent, func() {})

	// Start fade-in animation
	now := time.Now()
	modal.startAnimation(now)

	// Complete the fade-in animation
	var ops op.Ops
	gtx := app.NewContext(&ops, app.FrameEvent{
		Now: now.Add(300 * time.Millisecond), // Past animation duration
	})
	modal.updateAnimation(gtx)

	// Modal should be fully visible
	if modal.animationProgress != 1.0 {
		t.Error("Modal should be fully visible after fade-in")
	}

	if modal.isAnimating {
		t.Error("Modal should not be animating after fade-in complete")
	}

	// Start fade-out animation
	modal.startFadeOut()

	if !modal.isFadingOut {
		t.Error("Modal should be fading out")
	}

	if !modal.isAnimating {
		t.Error("Modal should be animating during fade-out")
	}

	// Test fade-out progress
	gtx = app.NewContext(&ops, app.FrameEvent{
		Now: now.Add(400 * time.Millisecond), // 100ms into fade-out
	})
	modal.updateAnimation(gtx)

	// Should be less than 1.0 (fading out)
	if modal.animationProgress >= 1.0 {
		t.Error("Modal should be fading out (progress < 1.0)")
	}

	// Complete the fade-out animation
	gtx = app.NewContext(&ops, app.FrameEvent{
		Now: now.Add(600 * time.Millisecond), // Past fade-out duration
	})
	modal.updateAnimation(gtx)

	// Modal should be invisible
	if modal.animationProgress != 0.0 {
		t.Error("Modal should be invisible after fade-out")
	}

	if modal.isAnimating {
		t.Error("Modal should not be animating after fade-out complete")
	}
}
