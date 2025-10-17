package fromage

import (
	"image/color"

	"gio.mleku.dev/font"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget"
)

// Label is a text widget with fluent method chaining
type Label struct {
	// Theme reference
	theme *Theme
	// Text content
	text string
	// Font specification
	font font.Font
	// Text color
	color color.NRGBA
	// Text alignment
	alignment text.Alignment
	// Maximum number of lines (0 = unlimited)
	maxLines int
	// Text size
	textSize unit.Sp
	// Text shaper
	shaper *text.Shaper
}

// NewLabel creates a new label with default settings
func (t *Theme) NewLabel() *Label {
	return &Label{
		theme:     t,
		text:      "",
		font:      font.Font{},
		color:     t.Colors.OnBackground(),
		alignment: text.Start,
		maxLines:  0,
		textSize:  unit.Sp(t.TextSize),
		shaper:    t.Shaper,
	}
}

// Text sets the text content
func (l *Label) Text(text string) *Label {
	l.text = text
	return l
}

// Font sets the font
func (l *Label) Font(font font.Font) *Label {
	l.font = font
	return l
}

// Color sets the text color
func (l *Label) Color(color color.NRGBA) *Label {
	l.color = color
	return l
}

// Alignment sets the text alignment
func (l *Label) Alignment(alignment text.Alignment) *Label {
	l.alignment = alignment
	return l
}

// MaxLines sets the maximum number of lines
func (l *Label) MaxLines(maxLines int) *Label {
	l.maxLines = maxLines
	return l
}

// TextSize sets the text size
func (l *Label) TextSize(size unit.Sp) *Label {
	l.textSize = size
	return l
}

// TextScale sets the text size relative to theme's base text size
func (l *Label) TextScale(scale float32) *Label {
	l.textSize = unit.Sp(float32(l.theme.TextSize) * scale)
	return l
}

// Layout renders the label
func (l *Label) Layout(g C) D {
	// Create the underlying Gio label
	label := widget.Label{
		Alignment: l.alignment,
		MaxLines:  l.maxLines,
	}

	// Record color operation
	textColorMacro := op.Record(g.Ops)
	paint.ColorOp{Color: l.color}.Add(g.Ops)
	textColor := textColorMacro.Stop()

	// Layout the label
	return label.Layout(g, l.shaper, l.font, l.textSize, l.text, textColor)
}

// Convenience methods for common text styles

// H1 creates a large heading
func (t *Theme) H1(text string) *Label {
	return t.NewLabel().Text(text).TextScale(2) // 96/16
}

// H2 creates a medium heading
func (t *Theme) H2(text string) *Label {
	return t.NewLabel().Text(text).TextScale(1.75) // 60/16
}

// H3 creates a small heading
func (t *Theme) H3(text string) *Label {
	return t.NewLabel().Text(text).TextScale(1.5) // 48/16
}

// H4 creates a smaller heading
func (t *Theme) H4(text string) *Label {
	return t.NewLabel().Text(text).TextScale(1.25) // 34/16
}

// H5 creates a small heading
func (t *Theme) H5(text string) *Label {
	return t.NewLabel().Text(text).TextScale(1) // 24/16
}

// H6 creates the smallest heading
func (t *Theme) H6(text string) *Label {
	return t.NewLabel().Text(text).TextScale(0.75) // 20/16
}

// Body1 creates normal body text
func (t *Theme) Body1(text string) *Label {
	return t.NewLabel().Text(text).TextScale(1.0)
}

// Body2 creates smaller body text
func (t *Theme) Body2(text string) *Label {
	return t.NewLabel().Text(text).TextScale(0.75) // 14/16
}

// Caption creates caption text
func (t *Theme) Caption(text string) *Label {
	return t.NewLabel().Text(text).TextScale(0.5) // 12/16
}
