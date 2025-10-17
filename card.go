package fromage

import (
	"image/color"

	"gio.mleku.dev/layout"
	"gio.mleku.dev/unit"
)

// Card provides a card widget with background color, padding, and rounded corners
type Card struct {
	// Theme reference
	theme *Theme
	// Background color
	backgroundColor color.NRGBA
	// Embedded widget
	widget W
	// Padding around the content
	padding unit.Dp
	// Corner radius for rounded corners
	cornerRadius float32
	// Whether to use all corners for rounding
	allCorners bool
}

// NewCard creates a new card widget
func (t *Theme) NewCard(widget W) *Card {
	return &Card{
		theme:           t,
		backgroundColor: t.Colors.Surface(),
		widget:          widget,
		padding:         t.TextSize,
		cornerRadius:    8,
		allCorners:      true,
	}
}

// NewCardWithColor creates a new card widget with a specific background color
func (t *Theme) NewCardWithColor(backgroundColor color.NRGBA, widget W) *Card {
	return &Card{
		theme:           t,
		backgroundColor: backgroundColor,
		widget:          widget,
		padding:         t.TextSize,
		cornerRadius:    8,
		allCorners:      true,
	}
}

// BackgroundColor sets the background color of the card
func (c *Card) BackgroundColor(color color.NRGBA) *Card {
	c.backgroundColor = color
	return c
}

// Widget sets the embedded widget
func (c *Card) Widget(widget W) *Card {
	c.widget = widget
	return c
}

// Padding sets the padding around the content
func (c *Card) Padding(padding unit.Dp) *Card {
	c.padding = padding
	return c
}

// CornerRadius sets the corner radius for rounded corners
func (c *Card) CornerRadius(radius float32) *Card {
	c.cornerRadius = radius
	return c
}

// AllCorners sets whether to round all corners
func (c *Card) AllCorners(allCorners bool) *Card {
	c.allCorners = allCorners
	return c
}

// Layout renders the card widget
func (c *Card) Layout(g C) D {
	// Create padding inset
	inset := layout.UniformInset(c.padding)

	// Create the fill widget with rounded corners
	var corners int
	if c.allCorners {
		corners = CornerAll
	}

	fill := c.theme.NewFillWithRadius(c.backgroundColor, c.cornerRadius, corners, nil)

	// Create the padded content
	paddedContent := func(gtx C) D {
		return inset.Layout(gtx, c.widget)
	}

	// Set the padded content as the fill widget
	fill.Widget(paddedContent)

	// Layout the fill
	return fill.Layout(g)
}

// Convenience methods for common card patterns

// CardPrimary creates a card with primary color background
func (t *Theme) CardPrimary(widget W) *Card {
	return t.NewCardWithColor(t.Colors.Primary(), widget)
}

// CardSecondary creates a card with secondary color background
func (t *Theme) CardSecondary(widget W) *Card {
	return t.NewCardWithColor(t.Colors.Secondary(), widget)
}

// CardSurface creates a card with surface color background (default)
func (t *Theme) CardSurface(widget W) *Card {
	return t.NewCard(widget)
}

// CardError creates a card with error color background
func (t *Theme) CardError(widget W) *Card {
	return t.NewCardWithColor(t.Colors.Error(), widget)
}

// CardWithTitle creates a card with a title and content
func (t *Theme) CardWithTitle(title string, content W) *Card {
	// Create a vertical layout with title and content
	cardContent := func(gtx C) D {
		return t.VFlex().
			SpaceStart().
			Rigid(func(gtx C) D {
				// Title label
				label := t.NewLabel().Text(title).Color(t.Colors.OnSurface()).TextScale(1.2) // Slightly larger for title
				return label.Layout(gtx)
			}).
			Rigid(func(gtx C) D {
				// Content
				return content(gtx)
			}).
			Layout(gtx)
	}

	return t.NewCard(cardContent)
}

// CardWithTitleAndColor creates a card with a title, content, and custom background color
func (t *Theme) CardWithTitleAndColor(title string, backgroundColor color.NRGBA, content W) *Card {
	// Create a vertical layout with title and content
	cardContent := func(gtx C) D {
		return t.VFlex().
			SpaceStart().
			Rigid(func(gtx C) D {
				// Title label
				label := t.NewLabel().Text(title).Color(t.Colors.OnSurface()).TextScale(1.2) // Slightly larger for title
				return label.Layout(gtx)
			}).
			Rigid(func(gtx C) D {
				// Content
				return content(gtx)
			}).
			Layout(gtx)
	}

	return t.NewCardWithColor(backgroundColor, cardContent)
}

// CardList creates a list of cards from multiple widgets
func (t *Theme) CardList(widgets ...W) W {
	return func(gtx C) D {
		// Create a vertical list of cards
		flex := t.VFlex().SpaceStart()
		for _, widget := range widgets {
			card := t.NewCard(widget)
			flex = flex.Rigid(card.Layout)
		}
		return flex.Layout(gtx)
	}
}

// CardListWithSpacing creates a list of cards with custom spacing
func (t *Theme) CardListWithSpacing(spacing unit.Dp, widgets ...W) W {
	return func(gtx C) D {
		// Create a vertical list of cards with spacing
		flex := t.VFlex().SpaceStart()
		for i, widget := range widgets {
			if i > 0 {
				// Add spacing between cards
				flex = flex.Rigid(func(gtx C) D {
					return layout.Spacer{Height: spacing}.Layout(gtx)
				})
			}
			card := t.NewCard(widget)
			flex = flex.Rigid(card.Layout)
		}
		return flex.Layout(gtx)
	}
}
