package fromage

import (
	"context"
	"image"
	"image/color"
	"testing"

	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/unit"
)

func TestBorder(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)
	border := th.NewBorder()

	if border == nil {
		t.Fatal("NewBorder returned nil")
	}

	if border.theme != th {
		t.Error("Border theme not set correctly")
	}

	if border.color != th.Colors.Outline() {
		t.Error("Border color not set to default outline color")
	}
}

func TestBorderFluentAPI(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test fluent API methods
	border := th.NewBorder().
		Color(color.NRGBA{R: 255, G: 0, B: 0, A: 255}).
		CornerRadius(unit.Dp(8)).
		Width(unit.Dp(2))

	if border.color.R != 255 {
		t.Error("Color not set correctly")
	}

	if border.cornerRadius != unit.Dp(8) {
		t.Error("CornerRadius not set correctly")
	}

	if border.width != unit.Dp(2) {
		t.Error("Width not set correctly")
	}
}

func TestBorderConvenienceMethods(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test convenience methods
	borderPrimary := th.BorderPrimary()
	if borderPrimary.color != th.Colors.Primary() {
		t.Error("BorderPrimary color not set correctly")
	}

	borderSecondary := th.BorderSecondary()
	if borderSecondary.color != th.Colors.Secondary() {
		t.Error("BorderSecondary color not set correctly")
	}

	borderOutline := th.BorderOutline()
	if borderOutline.color != th.Colors.Outline() {
		t.Error("BorderOutline color not set correctly")
	}

	borderError := th.BorderError()
	if borderError.color != th.Colors.Error() {
		t.Error("BorderError color not set correctly")
	}

	borderSurface := th.BorderSurface()
	if borderSurface.color != th.Colors.Surface() {
		t.Error("BorderSurface color not set correctly")
	}

	borderRounded := th.BorderRounded()
	expectedRadius := unit.Dp(float32(th.TextSize) * 0.5)
	if borderRounded.cornerRadius != expectedRadius {
		t.Error("BorderRounded corner radius not set correctly")
	}

	borderThick := th.BorderThick()
	expectedThickWidth := unit.Dp(float32(th.TextSize) * 0.25)
	if borderThick.width != expectedThickWidth {
		t.Error("BorderThick width not set correctly")
	}

	borderThin := th.BorderThin()
	expectedThinWidth := unit.Dp(float32(th.TextSize) * 0.0625)
	if borderThin.width != expectedThinWidth {
		t.Error("BorderThin width not set correctly")
	}
}

func TestBorderWidget(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Create a simple widget
	widget := func(g C) D {
		return D{Size: image.Pt(100, 50)}
	}

	border := th.NewBorder().Widget(widget)
	if border.widget == nil {
		t.Error("Widget not set correctly")
	}
}

func TestBorderLayout(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Create a simple widget
	widget := func(g C) D {
		return D{Size: image.Pt(100, 50)}
	}

	border := th.NewBorder().Widget(widget)

	// Create a mock context (this would normally be provided by Gio)
	// For testing purposes, we'll just verify the border is created correctly
	if border.widget == nil {
		t.Error("Border widget not set for layout")
	}
}
