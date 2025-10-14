package fromage

import (
	"image/color"
	"testing"

	"gioui.org/layout"
)

func TestFill(t *testing.T) {
	// Create a theme
	th := NewThemeWithMode(
		nil, // context not needed for this test
		NewColors,
		nil, // shaper not needed for this test
		16,  // text size
		ThemeModeLight,
	)

	// Test basic fill creation
	fill := th.NewFill(color.NRGBA{R: 255, G: 0, B: 0, A: 255}, nil)
	if fill == nil {
		t.Error("NewFill should return a non-nil fill")
	}

	// Test fill with radius
	fillWithRadius := th.NewFillWithRadius(color.NRGBA{R: 0, G: 255, B: 0, A: 255}, 8, CornerAll, nil)
	if fillWithRadius == nil {
		t.Error("NewFillWithRadius should return a non-nil fill")
	}
	if fillWithRadius.cornerRadius != 8 {
		t.Errorf("Expected corner radius 8, got %f", fillWithRadius.cornerRadius)
	}
	if fillWithRadius.corners != CornerAll {
		t.Errorf("Expected CornerAll, got %d", fillWithRadius.corners)
	}
}

func TestFillFluentAPI(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)

	// Test fluent API
	fill := th.NewFill(color.NRGBA{R: 255, G: 0, B: 0, A: 255}, nil).
		CornerRadius(4).
		Corners(CornerNW | CornerNE)

	if fill.cornerRadius != 4 {
		t.Errorf("Expected corner radius 4, got %f", fill.cornerRadius)
	}
	if fill.corners != (CornerNW | CornerNE) {
		t.Errorf("Expected CornerNW|CornerNE, got %d", fill.corners)
	}

	// Test color change
	newColor := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	fill.Color(newColor)
	if fill.color != newColor {
		t.Errorf("Expected new color, got %v", fill.color)
	}

	// Test direction change
	fill.Direction(layout.SE)
	if fill.direction != layout.SE {
		t.Errorf("Expected SE direction, got %v", fill.direction)
	}
}

func TestFillConvenienceMethods(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)

	// Test FillPrimary
	fillPrimary := th.FillPrimary(nil)
	if fillPrimary.color != th.Colors.Primary() {
		t.Errorf("FillPrimary should use primary color")
	}

	// Test FillSurface
	fillSurface := th.FillSurface(nil)
	if fillSurface.color != th.Colors.Surface() {
		t.Errorf("FillSurface should use surface color")
	}

	// Test FillCard
	fillCard := th.FillCard(nil)
	if fillCard.color != th.Colors.Surface() {
		t.Errorf("FillCard should use surface color")
	}
	if fillCard.cornerRadius != 8 {
		t.Errorf("FillCard should have corner radius 8")
	}
	if fillCard.corners != CornerAll {
		t.Errorf("FillCard should have all corners")
	}

	// Test FillButton
	fillButton := th.FillButton(nil)
	if fillButton.color != th.Colors.Primary() {
		t.Errorf("FillButton should use primary color")
	}
	if fillButton.cornerRadius != 4 {
		t.Errorf("FillButton should have corner radius 4")
	}
	if fillButton.corners != CornerAll {
		t.Errorf("FillButton should have all corners")
	}
}

func TestFillCornerFlags(t *testing.T) {
	// Test corner flag constants
	if CornerNW != 1 {
		t.Errorf("CornerNW should be 1, got %d", CornerNW)
	}
	if CornerNE != 2 {
		t.Errorf("CornerNE should be 2, got %d", CornerNE)
	}
	if CornerSW != 4 {
		t.Errorf("CornerSW should be 4, got %d", CornerSW)
	}
	if CornerSE != 8 {
		t.Errorf("CornerSE should be 8, got %d", CornerSE)
	}
	if CornerAll != 15 {
		t.Errorf("CornerAll should be 15, got %d", CornerAll)
	}
}

func TestFillIfCorner(t *testing.T) {
	th := NewThemeWithMode(nil, NewColors, nil, 16, ThemeModeLight)
	fill := th.NewFill(color.NRGBA{}, nil)

	// Test ifCorner method
	if fill.ifCorner(5.0, CornerNW) != 5.0 {
		t.Errorf("ifCorner should return radius when corner flag is set")
	}
	if fill.ifCorner(5.0, 0) != 0.0 {
		t.Errorf("ifCorner should return 0 when corner flag is not set")
	}
}
