package fromage

import (
	"testing"
)

func TestThemeMode(t *testing.T) {
	// Test light theme
	lightColors := NewColorsWithMode(ThemeModeLight)
	if lightColors.ThemeMode() != ThemeModeLight {
		t.Errorf("Expected light theme mode, got %v", lightColors.ThemeMode())
	}

	// Test dark theme
	darkColors := NewColorsWithMode(ThemeModeDark)
	if darkColors.ThemeMode() != ThemeModeDark {
		t.Errorf("Expected dark theme mode, got %v", darkColors.ThemeMode())
	}

	// Test theme switching
	colors := NewColorsWithMode(ThemeModeLight)
	if colors.ThemeMode() != ThemeModeLight {
		t.Errorf("Expected light theme mode, got %v", colors.ThemeMode())
	}

	colors.ToggleTheme()
	if colors.ThemeMode() != ThemeModeDark {
		t.Errorf("Expected dark theme mode after toggle, got %v", colors.ThemeMode())
	}

	colors.ToggleTheme()
	if colors.ThemeMode() != ThemeModeLight {
		t.Errorf("Expected light theme mode after second toggle, got %v", colors.ThemeMode())
	}
}

func TestThemeColors(t *testing.T) {
	// Test that light and dark themes have different background colors
	lightColors := NewColorsWithMode(ThemeModeLight)
	darkColors := NewColorsWithMode(ThemeModeDark)

	lightBg := lightColors.Background()
	darkBg := darkColors.Background()

	if lightBg == darkBg {
		t.Errorf("Light and dark themes should have different background colors")
	}

	// Test that light theme has light background
	if lightBg.R < 200 || lightBg.G < 200 || lightBg.B < 200 {
		t.Errorf("Light theme should have a light background color")
	}

	// Test that dark theme has dark background
	if darkBg.R > 100 || darkBg.G > 100 || darkBg.B > 100 {
		t.Errorf("Dark theme should have a dark background color")
	}
}
