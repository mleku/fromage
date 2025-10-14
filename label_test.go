package fromage

import (
	"testing"

	"gioui.org/text"
	"gioui.org/unit"
)

func TestLabel(t *testing.T) {
	// Create a theme
	th := NewThemeWithMode(
		nil, // context not needed for this test
		NewColors,
		nil, // shaper not needed for this test
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test basic label creation
	label := th.NewLabel()
	if label == nil {
		t.Error("NewLabel should return a non-nil label")
	}

	// Test fluent API
	label = th.NewLabel().
		Text("Hello, World!").
		TextSize(unit.Sp(14)).
		Alignment(text.Middle)

	if label.text != "Hello, World!" {
		t.Errorf("Expected text 'Hello, World!', got '%s'", label.text)
	}

	if label.textSize != unit.Sp(14) {
		t.Errorf("Expected text size 14sp, got %v", label.textSize)
	}

	if label.alignment != text.Middle {
		t.Errorf("Expected middle alignment, got %v", label.alignment)
	}
}

func TestLabelConvenienceMethods(t *testing.T) {
	th := NewThemeWithMode(
		nil,
		NewColors,
		nil,
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test H1
	h1 := th.H1("Heading 1")
	if h1.text != "Heading 1" {
		t.Errorf("H1 text should be 'Heading 1', got '%s'", h1.text)
	}

	// Test Body1
	body1 := th.Body1("Body text")
	if body1.text != "Body text" {
		t.Errorf("Body1 text should be 'Body text', got '%s'", body1.text)
	}

	// Test Caption
	caption := th.Caption("Caption text")
	if caption.text != "Caption text" {
		t.Errorf("Caption text should be 'Caption text', got '%s'", caption.text)
	}
}
