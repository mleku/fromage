package fromage

import (
	"context"
	"image/color"
	"testing"

	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
)

func TestIcon(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)
	icon := th.NewIcon()

	if icon == nil {
		t.Fatal("NewIcon returned nil")
	}

	if icon.theme != th {
		t.Error("Icon theme not set correctly")
	}

	if icon.color != th.Colors.OnSurface() {
		t.Error("Icon color not set to default on-surface color")
	}

	if icon.size != th.TextSize {
		t.Error("Icon size not set to default text size")
	}
}

func TestIconFluentAPI(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test fluent API methods
	icon := th.NewIcon().
		Color(color.NRGBA{R: 255, G: 0, B: 0, A: 255}).
		Scale(1.5).
		Size(unit.Dp(24))

	if icon.color.R != 255 {
		t.Error("Color not set correctly")
	}

	if icon.size != unit.Dp(24) {
		t.Error("Size not set correctly")
	}
}

func TestIconConvenienceMethods(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test convenience methods
	iconPrimary := th.IconPrimary()
	if iconPrimary.color != th.Colors.Primary() {
		t.Error("IconPrimary color not set correctly")
	}

	iconSecondary := th.IconSecondary()
	if iconSecondary.color != th.Colors.Secondary() {
		t.Error("IconSecondary color not set correctly")
	}

	iconTertiary := th.IconTertiary()
	if iconTertiary.color != th.Colors.Tertiary() {
		t.Error("IconTertiary color not set correctly")
	}

	iconOnSurface := th.IconOnSurface()
	if iconOnSurface.color != th.Colors.OnSurface() {
		t.Error("IconOnSurface color not set correctly")
	}

	iconOnBackground := th.IconOnBackground()
	if iconOnBackground.color != th.Colors.OnBackground() {
		t.Error("IconOnBackground color not set correctly")
	}

	iconError := th.IconError()
	if iconError.color != th.Colors.Error() {
		t.Error("IconError color not set correctly")
	}

	iconOutline := th.IconOutline()
	if iconOutline.color != th.Colors.Outline() {
		t.Error("IconOutline color not set correctly")
	}
}

func TestIconSizeConvenienceMethods(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test size convenience methods
	iconSmall := th.IconSmall()
	expectedSmallSize := unit.Dp(float32(th.TextSize) * 0.75)
	if iconSmall.size != expectedSmallSize {
		t.Error("IconSmall size not set correctly")
	}

	iconLarge := th.IconLarge()
	expectedLargeSize := unit.Dp(float32(th.TextSize) * 1.5)
	if iconLarge.size != expectedLargeSize {
		t.Error("IconLarge size not set correctly")
	}

	iconExtraLarge := th.IconExtraLarge()
	expectedExtraLargeSize := unit.Dp(float32(th.TextSize) * 2.0)
	if iconExtraLarge.size != expectedExtraLargeSize {
		t.Error("IconExtraLarge size not set correctly")
	}
}

func TestIconSrc(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test with nil data (should not crash)
	icon := th.NewIcon().Src(nil)
	if icon.src != nil {
		t.Error("Icon source should remain nil for invalid data")
	}

	// Test with empty data (should not crash)
	emptyData := []byte{}
	icon = th.NewIcon().Src(&emptyData)
	if icon.src != nil {
		t.Error("Icon source should remain nil for empty data")
	}
}

func TestIconLayout(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test layout with no source (should return empty dimensions)
	icon := th.NewIcon()
	// Note: We can't actually call Layout without a proper context,
	// but we can verify the icon is created correctly
	if icon.src != nil {
		t.Error("Icon should have nil source by default")
	}
}

func TestIconCache(t *testing.T) {
	th := NewThemeWithMode(
		context.Background(),
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test that icon cache is initialized
	if th.iconCache == nil {
		t.Error("Icon cache not initialized")
	}

	// Test that icon cache is empty initially
	if len(th.iconCache) != 0 {
		t.Error("Icon cache should be empty initially")
	}
}
