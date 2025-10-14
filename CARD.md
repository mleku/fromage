# Card Widget

The Card widget provides a container with background color, padding, and rounded corners, similar to Material Design cards.

## Features

- Background color customization
- Configurable padding
- Rounded corners with adjustable radius
- Method chaining for fluent API
- Integration with fromage theme system
- Support for titles and content
- Card lists with spacing

## Basic Usage

```go
// Create a simple card
card := theme.NewCard(contentWidget)

// Create a card with custom background color
card := theme.NewCardWithColor(color.NRGBA{R: 255, G: 0, B: 0, A: 255}, contentWidget)

// Create themed cards
primaryCard := theme.CardPrimary(contentWidget)
surfaceCard := theme.CardSurface(contentWidget)
errorCard := theme.CardError(contentWidget)
```

## Configuration

```go
card := theme.NewCard(contentWidget).
    BackgroundColor(color.NRGBA{R: 255, G: 255, B: 0, A: 255}).
    Padding(unit.Dp(20)).
    CornerRadius(12).
    AllCorners(false)
```

## Cards with Titles

```go
// Card with title
titleCard := theme.CardWithTitle("My Title", contentWidget)

// Card with title and custom color
customCard := theme.CardWithTitleAndColor("Custom Title", backgroundColor, contentWidget)
```

## Card Lists

```go
// Simple card list
cardList := theme.CardList(widget1, widget2, widget3)

// Card list with spacing
cardListWithSpacing := theme.CardListWithSpacing(unit.Dp(10), widget1, widget2, widget3)
```

## Background Theming

The card demo automatically fills the background with the appropriate theme color:
- **Light Mode**: Uses the light background color
- **Dark Mode**: Uses the dark background color

All text elements automatically use the appropriate "on-background" colors to ensure proper contrast and readability.

## Material Design Switch

The theme toggle uses a Material Design switch with the following features:
- **Pill-shaped track**: Rounded rectangle background
- **Animated thumb**: Circle that slides from left (off) to right (on)
- **Color transitions**: Track changes from dim gray (off) to primary color (on)
- **Smooth animations**: Thumb position and track color animate smoothly
- **Masked touch feedback**: Ink effects are constrained to the switch track area
- **Customizable colors**: Can set arbitrary background colors for the track

### Switch Usage

```go
// Standard switch
switch := theme.Switch(false)

// Switch with custom background color
switch := theme.SwitchWithColor(false, color.NRGBA{R: 255, G: 0, B: 0, A: 255})

// Configure switch dimensions
switch := theme.Switch(false).
    Width(unit.Dp(40)).
    Height(unit.Dp(24)).
    ThumbSize(unit.Dp(20))
```

## Demo

Run the card demo to see all card variations with theme switching:

```bash
go run cmd/cards/main.go
```

The demo includes:
- A Material Design theme toggle switch with animated thumb
- Background that changes color based on the current theme
- Various card styles (Primary, Surface, Error)
- Cards with titles
- Custom colored cards
- Real-time theme switching with full background color changes
- Smooth animations when toggling the switch
- Masked touch ink effects for precise interactive feedback

## API Reference

### Card Methods

- `BackgroundColor(color.NRGBA) *Card` - Set background color
- `Widget(W) *Card` - Set embedded widget
- `Padding(unit.Dp) *Card` - Set padding around content
- `CornerRadius(float32) *Card` - Set corner radius
- `AllCorners(bool) *Card` - Enable/disable all corner rounding
- `Layout(C) D` - Render the card

### Theme Methods

- `NewCard(W) *Card` - Create basic card
- `NewCardWithColor(color.NRGBA, W) *Card` - Create card with color
- `CardPrimary(W) *Card` - Create primary themed card
- `CardSecondary(W) *Card` - Create secondary themed card
- `CardSurface(W) *Card` - Create surface themed card
- `CardError(W) *Card` - Create error themed card
- `CardWithTitle(string, W) *Card` - Create card with title
- `CardWithTitleAndColor(string, color.NRGBA, W) *Card` - Create card with title and color
- `CardList(...W) W` - Create list of cards
- `CardListWithSpacing(unit.Dp, ...W) W` - Create list of cards with spacing
