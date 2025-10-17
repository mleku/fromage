package fromage

import (
	"context"
	"image"
	"testing"
	"time"

	"gio.mleku.dev/app"
	"gio.mleku.dev/op"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
)

func TestModalStack(t *testing.T) {
	th := NewThemeWithMode(context.TODO(), func() *Colors { return NewColors() }, text.NewShaper(), unit.Dp(16), ThemeModeLight)
	ms := th.NewModalStack()

	// Test initial state
	if !ms.IsEmpty() {
		t.Error("Expected modal stack to be empty initially")
	}
	if ms.Count() != 0 {
		t.Error("Expected modal count to be 0 initially")
	}

	// Test pushing a modal
	content := func(g C) D { return D{Size: image.Pt(100, 50)} }
	onClose := func() {}
	ms.Push(content, onClose)

	if ms.IsEmpty() {
		t.Error("Expected modal stack to not be empty after push")
	}
	if ms.Count() != 1 {
		t.Error("Expected modal count to be 1 after push")
	}

	// Test modal initialization
	modal := ms.modals[0]
	if modal.animationProgress != 0.0 {
		t.Error("Expected animation progress to be 0.0 initially")
	}
	if modal.isAnimating {
		t.Error("Expected modal to not be animating initially")
	}
	if modal.animationStarted {
		t.Error("Expected modal to not have started animation initially")
	}

	// Test pop (should start fade-out, not immediately remove)
	ms.Pop()
	if ms.Count() != 1 {
		t.Error("Expected modal count to remain 1 after pop (fade-out started)")
	}
	if !modal.isFadingOut {
		t.Error("Expected modal to be fading out after pop")
	}

	// Test clear
	ms.Clear()
	if !ms.IsEmpty() {
		t.Error("Expected modal stack to be empty after clear")
	}
}

func TestModal(t *testing.T) {
	th := NewThemeWithMode(context.TODO(), func() *Colors { return NewColors() }, text.NewShaper(), unit.Dp(16), ThemeModeLight)
	ms := th.NewModalStack()

	content := func(g C) D { return D{Size: image.Pt(100, 50)} }
	onClose := func() {}
	ms.Push(content, onClose)

	modal := ms.modals[0]

	// Test modal properties
	if modal.theme != th {
		t.Error("Expected modal theme to match")
	}
	if modal.onClose == nil {
		t.Error("Expected modal onClose to be set")
	}
	if !modal.blocking {
		t.Error("Expected modal to be blocking by default")
	}
	if modal.animationProgress != 0.0 {
		t.Error("Expected animation progress to be 0.0 initially")
	}
	if modal.isAnimating {
		t.Error("Expected modal to not be animating initially")
	}
	if modal.animationStarted {
		t.Error("Expected modal to not have started animation initially")
	}
}

func TestModalAnimation(t *testing.T) {
	th := NewThemeWithMode(context.TODO(), func() *Colors { return NewColors() }, text.NewShaper(), unit.Dp(16), ThemeModeLight)
	ms := th.NewModalStack()

	content := func(g C) D { return D{Size: image.Pt(100, 50)} }
	onClose := func() {}
	ms.Push(content, onClose)

	modal := ms.modals[0]

	// Create a context for testing
	ops := &op.Ops{}
	gtx := app.NewContext(ops, app.FrameEvent{})

	// Start animation
	modal.startAnimation(gtx.Now)

	if !modal.isAnimating {
		t.Error("Expected modal to be animating after start")
	}
	if !modal.animationStarted {
		t.Error("Expected modal to have started animation")
	}
	if modal.isFadingOut {
		t.Error("Expected modal to not be fading out initially")
	}

	// Simulate time passing (halfway through animation)
	halfwayTime := gtx.Now.Add(125 * time.Millisecond)
	gtx.Now = halfwayTime
	modal.updateAnimation(gtx)

	// Check animation progress (should be around 0.75 for ease-out at halfway point)
	if modal.animationProgress < 0.7 || modal.animationProgress > 0.8 {
		t.Errorf("Expected animation progress to be around 0.75, got %f", modal.animationProgress)
	}

	// Complete animation
	completeTime := gtx.Now.Add(125 * time.Millisecond)
	gtx.Now = completeTime
	modal.updateAnimation(gtx)

	if modal.isAnimating {
		t.Error("Expected modal to not be animating after completion")
	}
	if modal.animationProgress != 1.0 {
		t.Error("Expected animation progress to be 1.0 after completion")
	}
}

func TestModalFadeOut(t *testing.T) {
	th := NewThemeWithMode(context.TODO(), func() *Colors { return NewColors() }, text.NewShaper(), unit.Dp(16), ThemeModeLight)
	ms := th.NewModalStack()

	content := func(g C) D { return D{Size: image.Pt(100, 50)} }
	onClose := func() {}
	ms.Push(content, onClose)

	modal := ms.modals[0]

	// Start fade-in animation first
	ops := &op.Ops{}
	gtx := app.NewContext(ops, app.FrameEvent{})
	modal.startAnimation(gtx.Now)

	// Complete fade-in
	gtx.Now = gtx.Now.Add(250 * time.Millisecond)
	modal.updateAnimation(gtx)

	// Start fade-out
	modal.startFadeOut()
	gtx.Now = time.Now()

	if !modal.isAnimating {
		t.Error("Expected modal to be animating during fade-out")
	}
	if !modal.isFadingOut {
		t.Error("Expected modal to be fading out")
	}

	// Simulate time passing (halfway through fade-out)
	halfwayTime := gtx.Now.Add(125 * time.Millisecond)
	gtx.Now = halfwayTime
	modal.updateAnimation(gtx)

	// Check animation progress (should be around 0.25 for fade-out at halfway point)
	if modal.animationProgress < 0.2 || modal.animationProgress > 0.3 {
		t.Errorf("Expected animation progress to be around 0.25, got %f", modal.animationProgress)
	}

	// Complete fade-out
	completeTime := gtx.Now.Add(125 * time.Millisecond)
	gtx.Now = completeTime
	modal.updateAnimation(gtx)

	if modal.isAnimating {
		t.Error("Expected modal to not be animating after fade-out completion")
	}
	if modal.animationProgress != 0.0 {
		t.Error("Expected animation progress to be 0.0 after fade-out completion")
	}
}
