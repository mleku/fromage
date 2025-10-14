package fromage

import (
	"context"

	"gioui.org/text"
	"gioui.org/unit"
)

type Theme struct {
	ctx      context.Context
	Colors   *Colors
	Shaper   *text.Shaper
	TextSize unit.Dp
}

func NewTheme(
	ctx context.Context,
	colors func() *Colors,
	shaper *text.Shaper,
	textSize unit.Dp,
) *Theme {
	return NewThemeWithMode(ctx, colors, shaper, textSize, ThemeModeLight)
}

func NewThemeWithMode(
	ctx context.Context,
	colors func() *Colors,
	shaper *text.Shaper,
	textSize unit.Dp,
	mode ThemeMode,
) *Theme {
	return &Theme{
		ctx:      ctx,
		Colors:   NewColorsWithMode(mode),
		Shaper:   shaper,
		TextSize: textSize,
	}
}

// Theme mode methods
func (t *Theme) ThemeMode() ThemeMode {
	return t.Colors.ThemeMode()
}

func (t *Theme) SetThemeMode(mode ThemeMode) {
	t.Colors.SetThemeMode(mode)
}

func (t *Theme) ToggleTheme() {
	t.Colors.ToggleTheme()
}

func (t *Theme) IsDark() bool {
	return t.Colors.ThemeMode() == ThemeModeDark
}

func (t *Theme) IsLight() bool {
	return t.Colors.ThemeMode() == ThemeModeLight
}
