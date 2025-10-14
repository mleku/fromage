package fromage

import (
	"context"

	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Theme struct {
	ctx       context.Context
	Colors    *Colors
	Shaper    *text.Shaper
	TextSize  unit.Dp
	Pool      *Pool
	iconCache IconCache
}

// Pool manages widget instances to avoid creating new ones on every frame
type Pool struct {
	clickables      []*widget.Clickable
	clickablesInUse int
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
		ctx:       ctx,
		Colors:    NewColorsWithMode(mode),
		Shaper:    shaper,
		TextSize:  textSize,
		Pool:      &Pool{},
		iconCache: make(IconCache),
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

// Pool methods

// GetClickable returns a pooled clickable widget
func (p *Pool) GetClickable() *widget.Clickable {
	if len(p.clickables) <= p.clickablesInUse {
		// Allocate more clickables if needed
		for i := 0; i < 10; i++ {
			p.clickables = append(p.clickables, &widget.Clickable{})
		}
	}
	clickable := p.clickables[p.clickablesInUse]
	p.clickablesInUse++
	return clickable
}

// FreeClickable returns a clickable to the pool
func (p *Pool) FreeClickable(c *widget.Clickable) {
	for i := 0; i < p.clickablesInUse; i++ {
		if p.clickables[i] == c {
			if i != p.clickablesInUse-1 {
				// Move the item to the end
				tmp := p.clickables[i]
				p.clickables = append(p.clickables[:i], p.clickables[i+1:]...)
				p.clickables = append(p.clickables, tmp)
				p.clickablesInUse--
				break
			}
		}
	}
}

// Reset resets the pool usage counters
func (p *Pool) Reset() {
	p.clickablesInUse = 0
}
