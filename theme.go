package fromage

import (
	"context"

	"gioui.org/text"
)

type Theme struct {
	ctx    context.Context
	Colors *Colors
	Shaper *text.Shaper
}

func NewTheme(
	ctx context.Context,
	colors func() *Colors,
	shaper *text.Shaper) *Theme {

	return &Theme{
		ctx:    ctx,
		Colors: colors(),
		Shaper: shaper,
	}
}
