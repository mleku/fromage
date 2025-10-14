package fromage

import (
	"context"
	"image/color"
	"testing"

	"gioui.org/text"
	"gioui.org/unit"
)

func TestColorSelector(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(nil)),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Create a color selector
	cs := th.NewColorSelector()

	// Test initial color (should be red at full saturation)
	initialColor := cs.GetColor()
	if initialColor.R == 0 && initialColor.G == 0 && initialColor.B == 0 {
		t.Error("Initial color should not be black")
	}

	// Test setting a specific color
	testColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Red
	cs.SetColor(testColor)

	// Verify the color was set
	setColor := cs.GetColor()
	if setColor.R != 255 || setColor.G != 0 || setColor.B != 0 {
		t.Errorf("Expected red color, got R=%d G=%d B=%d", setColor.R, setColor.G, setColor.B)
	}
}

func TestHSVToRGB(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(nil)),
		unit.Dp(16),
		ThemeModeLight,
	)

	cs := th.NewColorSelector()

	// Test red (0 degrees)
	red := cs.hsvToRgb(0, 1.0, 1.0)
	if red.R != 255 || red.G != 0 || red.B != 0 {
		t.Errorf("Expected red, got R=%d G=%d B=%d", red.R, red.G, red.B)
	}

	// Test green (120 degrees)
	green := cs.hsvToRgb(120, 1.0, 1.0)
	if green.R != 0 || green.G != 255 || green.B != 0 {
		t.Errorf("Expected green, got R=%d G=%d B=%d", green.R, green.G, green.B)
	}

	// Test blue (240 degrees)
	blue := cs.hsvToRgb(240, 1.0, 1.0)
	if blue.R != 0 || blue.G != 0 || blue.B != 255 {
		t.Errorf("Expected blue, got R=%d G=%d B=%d", blue.R, blue.G, blue.B)
	}
}

func TestSurfaceTint(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(nil)),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test setting surface tint
	testTint := color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Red
	th.Colors.SetSurfaceTint(testTint)

	// Verify the tint was set
	if th.Colors.GetSurfaceTint() != testTint {
		t.Error("Surface tint was not set correctly")
	}

	// Test that surface color is different from background
	surfaceColor := th.Colors.Surface()
	backgroundColor := th.Colors.Background()

	if surfaceColor == backgroundColor {
		t.Error("Surface color should be different from background color when tinted")
	}
}
