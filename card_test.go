package fromage

import (
	"image/color"
	"testing"

	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/unit"
)

func TestCardCreation(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		nil,
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test basic card creation
	card := th.NewCard(nil)
	if card == nil {
		t.Fatal("NewCard returned nil")
	}

	// Test card with color
	cardWithColor := th.NewCardWithColor(color.NRGBA{R: 255, G: 0, B: 0, A: 255}, nil)
	if cardWithColor == nil {
		t.Fatal("NewCardWithColor returned nil")
	}

	// Test convenience methods
	primaryCard := th.CardPrimary(nil)
	if primaryCard == nil {
		t.Fatal("CardPrimary returned nil")
	}

	surfaceCard := th.CardSurface(nil)
	if surfaceCard == nil {
		t.Fatal("CardSurface returned nil")
	}

	errorCard := th.CardError(nil)
	if errorCard == nil {
		t.Fatal("CardError returned nil")
	}
}

func TestCardWithTitle(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		nil,
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test card with title
	content := func(gtx C) D {
		return D{Size: gtx.Constraints.Max}
	}

	titleCard := th.CardWithTitle("Test Title", content)
	if titleCard == nil {
		t.Fatal("CardWithTitle returned nil")
	}

	// Test card with title and color
	customCard := th.CardWithTitleAndColor("Custom Title", color.NRGBA{R: 0, G: 255, B: 0, A: 255}, content)
	if customCard == nil {
		t.Fatal("CardWithTitleAndColor returned nil")
	}
}

func TestCardList(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		nil,
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test card list
	widget1 := func(gtx C) D {
		return D{Size: gtx.Constraints.Max}
	}
	widget2 := func(gtx C) D {
		return D{Size: gtx.Constraints.Max}
	}

	cardList := th.CardList(widget1, widget2)
	if cardList == nil {
		t.Fatal("CardList returned nil")
	}

	// Test card list with spacing
	cardListWithSpacing := th.CardListWithSpacing(unit.Dp(10), widget1, widget2)
	if cardListWithSpacing == nil {
		t.Fatal("CardListWithSpacing returned nil")
	}
}

func TestCardConfiguration(t *testing.T) {
	// Create a test theme
	th := NewThemeWithMode(
		nil,
		NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		ThemeModeLight,
	)

	// Test card configuration methods
	card := th.NewCard(nil)

	// Test method chaining
	configuredCard := card.
		BackgroundColor(color.NRGBA{R: 255, G: 255, B: 0, A: 255}).
		Padding(unit.Dp(20)).
		CornerRadius(12).
		AllCorners(false)

	if configuredCard == nil {
		t.Fatal("Card configuration returned nil")
	}

	// Verify configuration was applied
	if configuredCard.backgroundColor.R != 255 {
		t.Error("Background color not set correctly")
	}
	if configuredCard.padding != unit.Dp(20) {
		t.Error("Padding not set correctly")
	}
	if configuredCard.cornerRadius != 12 {
		t.Error("Corner radius not set correctly")
	}
	if configuredCard.allCorners != false {
		t.Error("All corners not set correctly")
	}
}
