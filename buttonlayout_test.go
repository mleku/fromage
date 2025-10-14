package fromage

import (
	"image/color"
	"testing"

	"gioui.org/unit"
)

func TestButtonLayout(t *testing.T) {
	// Create a theme
	th := NewThemeWithMode(
		nil, // context not needed for this test
		NewColors,
		nil, // shaper not needed for this test
		16,  // text size
		ThemeModeLight,
	)

	// Test basic button creation
	button := th.NewButtonLayout()
	if button == nil {
		t.Error("NewButtonLayout should return a non-nil button")
	}

	// Test default values
	if button.background != th.Colors.Primary() {
		t.Errorf("Expected primary color, got %v", button.background)
	}
	expectedRadius := unit.Dp(float32(th.TextSize) * 0.25)
	if button.cornerRadius != expectedRadius {
		t.Errorf("Expected corner radius %v, got %v", expectedRadius, button.cornerRadius)
	}
	if button.corners != CornerAll {
		t.Errorf("Expected CornerAll, got %d", button.corners)
	}
	if button.disabled {
		t.Error("Button should not be disabled by default")
	}
}

func TestButtonLayoutFluentAPI(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)

	// Test fluent API
	newColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	button := th.NewButtonLayout().
		Background(newColor).
		CornerRadius(0.5). // Scale factor
		Corners(CornerNW | CornerNE).
		Disabled(true)

	if button.background != newColor {
		t.Errorf("Expected new color, got %v", button.background)
	}
	expectedRadius := unit.Dp(float32(th.TextSize) * 0.5)
	if button.cornerRadius != expectedRadius {
		t.Errorf("Expected corner radius %v, got %v", expectedRadius, button.cornerRadius)
	}
	if button.corners != (CornerNW | CornerNE) {
		t.Errorf("Expected CornerNW|CornerNE, got %d", button.corners)
	}
	if !button.disabled {
		t.Error("Button should be disabled")
	}
}

func TestButtonLayoutConvenienceMethods(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)

	// Test PrimaryButton
	primaryButton := th.PrimaryButton(nil)
	if primaryButton.background != th.Colors.Primary() {
		t.Errorf("PrimaryButton should use primary color")
	}

	// Test SecondaryButton
	secondaryButton := th.SecondaryButton(nil)
	if secondaryButton.background != th.Colors.Secondary() {
		t.Errorf("SecondaryButton should use secondary color")
	}

	// Test SurfaceButton
	surfaceButton := th.SurfaceButton(nil)
	if surfaceButton.background != th.Colors.Surface() {
		t.Errorf("SurfaceButton should use surface color")
	}

	// Test ErrorButton
	errorButton := th.ErrorButton(nil)
	if errorButton.background != th.Colors.Error() {
		t.Errorf("ErrorButton should use error color")
	}

	// Test RoundedButton
	roundedButton := th.RoundedButton(nil)
	expectedRoundedRadius := unit.Dp(float32(th.TextSize) * 0.5)
	if roundedButton.cornerRadius != expectedRoundedRadius {
		t.Errorf("RoundedButton should have corner radius %v", expectedRoundedRadius)
	}

	// Test PillButton
	pillButton := th.PillButton(nil)
	expectedPillRadius := unit.Dp(1000) // Large value indicating pill shape
	if pillButton.cornerRadius != expectedPillRadius {
		t.Errorf("PillButton should have corner radius %v", expectedPillRadius)
	}
}

func TestButtonLayoutIfCorner(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)
	button := th.NewButtonLayout()

	// Test ifCorner method
	if button.ifCorner(5.0, CornerNW) != 5.0 {
		t.Errorf("ifCorner should return radius when corner flag is set")
	}
	if button.ifCorner(5.0, 0) != 0.0 {
		t.Errorf("ifCorner should return 0 when corner flag is not set")
	}
}

func TestButtonLayoutPillRadius(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)

	// Test PillRadius method
	button := th.NewButtonLayout().PillRadius()
	expectedRadius := unit.Dp(1000) // Large value indicating pill shape
	if button.cornerRadius != expectedRadius {
		t.Errorf("PillRadius should set corner radius to %v, got %v", expectedRadius, button.cornerRadius)
	}
}

func TestButtonLayoutTextButton(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)

	// Test TextButton creation
	textButton := th.TextButton("Test")
	if textButton == nil {
		t.Error("TextButton should return a non-nil button")
	}
	if textButton.widget == nil {
		t.Error("TextButton should have a widget")
	}
}

func TestButtonLayoutIconButton(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)

	// Test IconButton creation
	iconButton := th.IconButton("★")
	if iconButton == nil {
		t.Error("IconButton should return a non-nil button")
	}
	if iconButton.widget == nil {
		t.Error("IconButton should have a widget")
	}
	expectedIconRadius := unit.Dp(1000) // Large value indicating pill shape
	if iconButton.cornerRadius != expectedIconRadius {
		t.Errorf("IconButton should have corner radius %v", expectedIconRadius)
	}
}
