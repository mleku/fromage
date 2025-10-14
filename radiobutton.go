package fromage

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// RadioButton represents a single radio button
type RadioButton struct {
	*Window
	clickable *widget.Clickable
	checked   bool
	onChange  func(bool)
	label     string
	size      unit.Dp
}

// RadioButtonGroup manages a group of radio buttons where only one can be selected
type RadioButtonGroup struct {
	*Window
	buttons  []*RadioButton
	selected int
	onChange func(int, string)
	layout   LayoutDirection
	spacing  unit.Dp
}

// LayoutDirection specifies how radio buttons should be arranged
type LayoutDirection int

const (
	LayoutVertical LayoutDirection = iota
	LayoutHorizontal
)

// NewRadioButton creates a new radio button
func (w *Window) NewRadioButton(checked bool) *RadioButton {
	return &RadioButton{
		Window:    w,
		clickable: &widget.Clickable{}, // Create clickable directly like checkbox does
		checked:   checked,
		size:      unit.Dp(float32(w.Theme.TextSize) * 2.0), // Make radio buttons larger for better clickability
	}
}

// Label sets the label text for the radio button
func (rb *RadioButton) Label(label string) *RadioButton {
	rb.label = label
	return rb
}

// Size sets the size of the radio button circle
func (rb *RadioButton) Size(size unit.Dp) *RadioButton {
	rb.size = size
	return rb
}

// SetOnChange sets the callback function when the radio button state changes
func (rb *RadioButton) SetOnChange(fn func(bool)) *RadioButton {
	rb.onChange = fn
	return rb
}

// SetChecked sets the checked state of the radio button
func (rb *RadioButton) SetChecked(checked bool) *RadioButton {
	rb.checked = checked
	return rb
}

// IsChecked returns whether the radio button is currently checked
func (rb *RadioButton) IsChecked() bool {
	return rb.checked
}

// Clicked returns true if the radio button was clicked
func (rb *RadioButton) Clicked(gtx C) bool {
	return rb.clickable.Clicked(gtx)
}

// Layout renders the radio button
func (rb *RadioButton) Layout(gtx C) D {
	// Handle clicks first
	if rb.clickable.Clicked(gtx) {
		// Always set to checked when clicked (group logic handles unchecking others)
		rb.checked = true
		if rb.onChange != nil {
			rb.onChange(true)
		}
	}

	// Create the radio button layout with the entire area clickable
	return rb.clickable.Layout(gtx, func(g C) D {
		return rb.Theme.HFlex().
			SpaceEvenly().
			AlignMiddle().
			Rigid(func(g C) D {
				// Radio button circle
				return rb.layoutRadioCircle(g)
			}).
			Rigid(func(g C) D {
				// Label text
				if rb.label != "" {
					return rb.Theme.Body2(rb.label).
						Color(rb.Theme.Colors.OnBackground()).
						Alignment(text.Start).
						Layout(g)
				}
				return D{}
			}).
			Layout(g)
	})
}

// layoutRadioCircle renders the circular radio button
func (rb *RadioButton) layoutRadioCircle(gtx C) D {
	// Calculate circle size
	circleSize := gtx.Dp(rb.size)

	// Create a square constraint for the circle
	circleConstraints := layout.Exact(image.Pt(circleSize, circleSize))
	gtx.Constraints = circleConstraints

	// Create the circle using border with full radius
	circleWidget := rb.Theme.NewBorder().
		Color(rb.Theme.Colors.OnBackground()).              // Always use text color for the outer circle
		CornerRadius(rb.size / 2).                          // Half the size for perfect circle
		Width(unit.Dp(float32(rb.Theme.TextSize) * 0.125)). // Thin border
		Widget(func(g C) D {
			// Inner content - either empty or filled circle
			if rb.checked {
				// Draw filled circle inside, centered
				innerSize := gtx.Dp(rb.size * 0.4) // 40% of outer size

				// Calculate center position
				centerX := (circleSize - innerSize) / 2
				centerY := (circleSize - innerSize) / 2

				// Create centered filled circle
				innerCircle := image.Rectangle{
					Min: image.Pt(centerX, centerY),
					Max: image.Pt(centerX+innerSize, centerY+innerSize),
				}
				paint.FillShape(g.Ops,
					rb.Theme.Colors.Primary(), // Use primary color for the filled circle to make it more visible
					clip.Ellipse{Min: innerCircle.Min, Max: innerCircle.Max}.Op(g.Ops),
				)
			}
			// Debug: print the checked state
			if rb.label != "" {
				// fmt.Printf("Rendering radio button '%s' - checked: %v\n", rb.label, rb.checked)
			}
			return D{Size: image.Pt(circleSize, circleSize)}
		})

	return circleWidget.Layout(gtx)
}

// NewRadioButtonGroup creates a new radio button group
func (w *Window) NewRadioButtonGroup() *RadioButtonGroup {
	return &RadioButtonGroup{
		Window:  w,
		layout:  LayoutVertical,
		spacing: unit.Dp(float32(w.Theme.TextSize) * 0.5), // Half text size spacing
	}
}

// AddButton adds a radio button to the group
func (rbg *RadioButtonGroup) AddButton(label string, checked bool) *RadioButtonGroup {
	button := rbg.NewRadioButton(checked).Label(label)

	// Set up the button's onChange to handle group selection
	originalOnChange := button.onChange
	buttonIndex := len(rbg.buttons)

	button.SetOnChange(func(b bool) {
		if b {
			// Uncheck all other buttons in the group
			for i, btn := range rbg.buttons {
				if i != buttonIndex {
					btn.SetChecked(false)
				}
			}
			rbg.selected = buttonIndex
			if rbg.onChange != nil {
				rbg.onChange(buttonIndex, label)
			}
			// Request a new frame to update the display
			rbg.Window.Invalidate()
		}
		// Call the original onChange if it exists
		if originalOnChange != nil {
			originalOnChange(b)
		}
	})

	rbg.buttons = append(rbg.buttons, button)
	if checked {
		rbg.selected = buttonIndex
		// Trigger callback for initial selection
		if rbg.onChange != nil {
			rbg.onChange(buttonIndex, label)
		}
	}

	return rbg
}

// SetLayout sets the layout direction for the group
func (rbg *RadioButtonGroup) SetLayout(direction LayoutDirection) *RadioButtonGroup {
	rbg.layout = direction
	return rbg
}

// SetSpacing sets the spacing between radio buttons
func (rbg *RadioButtonGroup) SetSpacing(spacing unit.Dp) *RadioButtonGroup {
	rbg.spacing = spacing
	return rbg
}

// SetOnChange sets the callback function when the selected radio button changes
func (rbg *RadioButtonGroup) SetOnChange(fn func(int, string)) *RadioButtonGroup {
	rbg.onChange = fn
	return rbg
}

// GetSelected returns the index and label of the currently selected radio button
func (rbg *RadioButtonGroup) GetSelected() (int, string) {
	if rbg.selected >= 0 && rbg.selected < len(rbg.buttons) {
		return rbg.selected, rbg.buttons[rbg.selected].label
	}
	return -1, ""
}

// SetSelected sets the selected radio button by index
func (rbg *RadioButtonGroup) SetSelected(index int) *RadioButtonGroup {
	if index >= 0 && index < len(rbg.buttons) {
		// Uncheck all buttons
		for _, btn := range rbg.buttons {
			btn.SetChecked(false)
		}
		// Check the selected button
		rbg.buttons[index].SetChecked(true)
		rbg.selected = index
		// Trigger callback for programmatic selection
		if rbg.onChange != nil {
			rbg.onChange(index, rbg.buttons[index].label)
		}
		// Request a new frame to update the display
		rbg.Window.Invalidate()
	}
	return rbg
}

// Layout renders the radio button group
func (rbg *RadioButtonGroup) Layout(gtx C) D {
	if len(rbg.buttons) == 0 {
		return D{}
	}

	// Synchronize button states with group selection
	rbg.syncButtonStates()

	if rbg.layout == LayoutVertical {
		// Vertical layout
		flex := rbg.Theme.VFlex().SpaceEvenly()
		for _, button := range rbg.buttons {
			flex = flex.Rigid(button.Layout)
		}
		return flex.Layout(gtx)
	} else {
		// Horizontal layout
		flex := rbg.Theme.HFlex().SpaceEvenly()
		for _, button := range rbg.buttons {
			flex = flex.Rigid(button.Layout)
		}
		return flex.Layout(gtx)
	}
}

// syncButtonStates ensures that button states match the group's selected state
func (rbg *RadioButtonGroup) syncButtonStates() {
	for i, button := range rbg.buttons {
		shouldBeChecked := (i == rbg.selected)
		if button.IsChecked() != shouldBeChecked {
			button.SetChecked(shouldBeChecked)
		}
	}
}

// Convenience methods for common radio button group configurations

// VerticalRadioGroup creates a vertical radio button group
func (w *Window) VerticalRadioGroup() *RadioButtonGroup {
	return w.NewRadioButtonGroup().SetLayout(LayoutVertical)
}

// HorizontalRadioGroup creates a horizontal radio button group
func (w *Window) HorizontalRadioGroup() *RadioButtonGroup {
	return w.NewRadioButtonGroup().SetLayout(LayoutHorizontal)
}
