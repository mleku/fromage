package fromage

import (
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
)

func TestInset(t *testing.T) {
	// Create a test theme
	th := NewTheme(nil, func() *Colors { return NewColors() }, &text.Shaper{}, unit.Dp(16))

	// Create a test window
	w := NewWindow(th)

	// Test creating an inset
	inset := w.Inset(0.5, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{
			Size: image.Point{X: 100, Y: 50},
		}
	})

	if inset == nil {
		t.Fatal("Inset should not be nil")
	}

	if inset.Window != w {
		t.Error("Inset should reference the correct window")
	}

	if inset.w == nil {
		t.Error("Inset should have a widget")
	}
}

func TestInsetEmbed(t *testing.T) {
	// Create a test theme
	th := NewTheme(nil, func() *Colors { return NewColors() }, &text.Shaper{}, unit.Dp(16))

	// Create a test window
	w := NewWindow(th)

	// Create an inset
	inset := w.Inset(0.5, nil)

	// Test embedding a widget
	testWidget := func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{
			Size: image.Point{X: 100, Y: 50},
		}
	}

	result := inset.Embed(testWidget)

	if result != inset {
		t.Error("Embed should return the same inset instance")
	}

	if inset.w == nil {
		t.Error("Inset should have the embedded widget")
	}
}

func TestInsetFn(t *testing.T) {
	// Create a test theme
	th := NewTheme(nil, func() *Colors { return NewColors() }, &text.Shaper{}, unit.Dp(16))

	// Create a test window
	w := NewWindow(th)

	// Create a test widget
	testWidget := func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{
			Size: image.Point{X: 100, Y: 50},
		}
	}

	// Create an inset with the test widget
	inset := w.Inset(0.5, testWidget)

	// Create a test context
	var ops op.Ops
	gtx := layout.Context{
		Ops: &ops,
		Constraints: layout.Constraints{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: 200, Y: 100},
		},
	}

	// Test the Fn method
	dims := inset.Fn(gtx)

	// The dimensions should include the padding
	expectedWidth := 100 + int(float32(unit.Dp(16))*0.5)*2 // widget width + padding on both sides
	expectedHeight := 50 + int(float32(unit.Dp(16))*0.5)*2 // widget height + padding on both sides

	if dims.Size.X < expectedWidth {
		t.Errorf("Expected width to be at least %d, got %d", expectedWidth, dims.Size.X)
	}

	if dims.Size.Y < expectedHeight {
		t.Errorf("Expected height to be at least %d, got %d", expectedHeight, dims.Size.Y)
	}
}
