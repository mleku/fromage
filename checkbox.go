package fromage

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
)

// CheckboxHook is a function type for handling checkbox value changes
type CheckboxHook func(b bool)

// Checkbox represents a checkbox widget with a label
type Checkbox struct {
	// Theme reference
	theme *Theme
	// Current boolean value
	value bool
	// Clickable widget for handling interactions
	clickable *widget.Clickable
	// Whether the value has changed since last check
	changed bool
	// Callback function for value changes
	onChange CheckboxHook
	// Visual styling
	label           string
	labelColor      color.NRGBA
	iconColor       color.NRGBA
	backgroundColor color.NRGBA
	borderColor     color.NRGBA
	cornerRadius    unit.Dp
	// Size constraints
	size     unit.Dp // Size of the checkbox square
	textSize unit.Sp // Size of the label text
	// Font styling
	font font.Font
	// Icon styling
	checkedIcon   *widget.Icon
	uncheckedIcon *widget.Icon
}

// NewCheckbox creates a new checkbox widget
func (t *Theme) NewCheckbox(value bool) *Checkbox {
	return &Checkbox{
		theme:           t,
		value:           value,
		clickable:       &widget.Clickable{},
		changed:         false,
		onChange:        func(b bool) {},
		label:           "Checkbox",
		labelColor:      t.Colors.OnBackground(),
		iconColor:       t.Colors.OnBackground(), // Use theme text color
		backgroundColor: t.Colors.Surface(),
		borderColor:     t.Colors.OnBackground(), // Use theme text color
		cornerRadius:    unit.Dp(2),
		size:            unit.Dp(20),
		textSize:        unit.Sp(float32(t.TextSize) * 14.0 / 16.0),
		font:            font.Font{},
		checkedIcon:     nil, // Will use default checkmark
		uncheckedIcon:   nil, // Will use default empty box
	}
}

// Value sets the checkbox value
func (c *Checkbox) Value(value bool) *Checkbox {
	if c.value != value {
		c.value = value
		c.changed = true
	}
	return c
}

// GetValue returns the current checkbox value
func (c *Checkbox) GetValue() bool {
	return c.value
}

// SetOnChange sets the callback function for value changes
func (c *Checkbox) SetOnChange(fn CheckboxHook) *Checkbox {
	c.onChange = fn
	return c
}

// Label sets the checkbox label text
func (c *Checkbox) Label(label string) *Checkbox {
	c.label = label
	return c
}

// LabelColor sets the color of the label text
func (c *Checkbox) LabelColor(color color.NRGBA) *Checkbox {
	c.labelColor = color
	return c
}

// IconColor sets the color of the checkbox icon
func (c *Checkbox) IconColor(color color.NRGBA) *Checkbox {
	c.iconColor = color
	return c
}

// BackgroundColor sets the background color of the checkbox
func (c *Checkbox) BackgroundColor(color color.NRGBA) *Checkbox {
	c.backgroundColor = color
	return c
}

// BorderColor sets the border color of the checkbox
func (c *Checkbox) BorderColor(color color.NRGBA) *Checkbox {
	c.borderColor = color
	return c
}

// CornerRadius sets the corner radius of the checkbox
func (c *Checkbox) CornerRadius(radius unit.Dp) *Checkbox {
	c.cornerRadius = radius
	return c
}

// Size sets the size of the checkbox square
func (c *Checkbox) Size(size unit.Dp) *Checkbox {
	c.size = size
	return c
}

// TextSize sets the size of the label text
func (c *Checkbox) TextSize(size unit.Sp) *Checkbox {
	c.textSize = size
	return c
}

// Font sets the font for the label
func (c *Checkbox) Font(font font.Font) *Checkbox {
	c.font = font
	return c
}

// Changed returns true if the value has changed since the last call
func (c *Checkbox) Changed() bool {
	changed := c.changed
	c.changed = false
	return changed
}

// Clicked returns true if the widget was clicked
func (c *Checkbox) Clicked(g C) bool {
	return c.clickable.Clicked(g)
}

// Hovered returns true if the widget is being hovered
func (c *Checkbox) Hovered() bool {
	return c.clickable.Hovered()
}

// Pressed returns true if the widget is being pressed
func (c *Checkbox) Pressed() bool {
	return c.clickable.Pressed()
}

// Layout renders the checkbox widget
func (c *Checkbox) Layout(g C) D {
	// Handle click events BEFORE layout
	if c.clickable.Clicked(g) {
		c.value = !c.value
		c.changed = true
		if c.onChange != nil {
			c.onChange(c.value)
		}
	}

	// Calculate dimensions
	checkboxSize := g.Dp(c.size)
	textSize := g.Sp(c.textSize)

	// Create the layout using the clickable's Layout method
	return c.clickable.Layout(g, func(g C) D {
		// Add semantic information for accessibility
		semantic.CheckBox.Add(g.Ops)

		// Layout checkbox and label horizontally
		return c.theme.HFlex().
			AlignMiddle().
			SpaceStart().
			Rigid(func(g C) D {
				// Draw the checkbox square
				return c.drawCheckbox(g, checkboxSize)
			}).
			Rigid(func(g C) D {
				// Add spacing between checkbox and label
				return layout.Inset{Left: unit.Dp(8)}.Layout(g, func(g C) D {
					// Draw the label text
					return c.drawLabel(g, textSize)
				})
			}).
			Layout(g)
	})
}

// drawCheckbox draws the checkbox square and checkmark
func (c *Checkbox) drawCheckbox(g C, size int) D {
	// Create a widget function that draws the checkmark if checked
	checkboxWidget := func(g C) D {
		// Draw checkmark if checked
		if c.value {
			c.drawCheckmark(g, size)
		}
		return D{Size: image.Pt(size, size)}
	}

	// Use the existing border widget to draw the square outline
	border := c.theme.NewBorder().
		Color(c.borderColor).
		CornerRadius(c.cornerRadius).
		Width(unit.Dp(2)).
		Widget(checkboxWidget)

	return border.Layout(g)
}

// drawCheckmark draws a checkmark using the SVG icon
func (c *Checkbox) drawCheckmark(g C, size int) {
	// SVG content with only the checkmark path (removed the background rectangle)
	svgContent := `<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" fill="currentColor"/></svg>`

	// Create icon from SVG using ivgconv
	icon := c.theme.NewIconFromSVG(svgContent).
		Color(c.iconColor).
		Size(unit.Dp(float32(size) * 0.5)) // Scale icon to half the checkbox size

	// Center the icon in the checkbox
	centerX := size / 2
	centerY := size / 2

	// Layout the icon centered in the checkbox
	iconDims := icon.Layout(g)

	// Calculate offset to center the icon
	offsetX := centerX - iconDims.Size.X/2
	offsetY := centerY - iconDims.Size.Y/2

	// Apply the offset
	defer op.Offset(image.Pt(offsetX, offsetY)).Push(g.Ops).Pop()
	icon.Layout(g)
}

// drawLabel draws the checkbox label text
func (c *Checkbox) drawLabel(g C, textSize int) D {
	// Create label widget
	label := c.theme.NewLabel().
		Text(c.label).
		Color(c.labelColor).
		Font(c.font).
		TextSize(c.textSize)

	return label.Layout(g)
}

// Convenience methods for common checkbox patterns

// PrimaryCheckbox creates a checkbox with primary color styling
func (t *Theme) PrimaryCheckbox(value bool) *Checkbox {
	return t.NewCheckbox(value)
}

// SecondaryCheckbox creates a checkbox with secondary color styling
func (t *Theme) SecondaryCheckbox(value bool) *Checkbox {
	return t.NewCheckbox(value)
}

// SurfaceCheckbox creates a checkbox with surface color styling
func (t *Theme) SurfaceCheckbox(value bool) *Checkbox {
	return t.NewCheckbox(value)
}

// Helper function for max integer
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
