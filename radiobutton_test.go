package fromage

import (
	"testing"

	"gioui.org/unit"
)

func TestRadioButtonCreation(t *testing.T) {
	// Create a mock theme for testing
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Test creating a radio button
	rb := w.NewRadioButton(true)
	if rb == nil {
		t.Fatal("NewRadioButton returned nil")
	}

	if !rb.IsChecked() {
		t.Error("Radio button should be checked")
	}

	// Test setting label
	rb.Label("Test Label")
	if rb.label != "Test Label" {
		t.Error("Label not set correctly")
	}

	// Test setting size
	rb.Size(unit.Dp(20))
	if rb.size != unit.Dp(20) {
		t.Error("Size not set correctly")
	}
}

func TestRadioButtonGroup(t *testing.T) {
	// Create a mock theme for testing
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Test creating a radio button group
	rbg := w.NewRadioButtonGroup()
	if rbg == nil {
		t.Fatal("NewRadioButtonGroup returned nil")
	}

	// Test adding buttons
	rbg.AddButton("Option 1", true).
		AddButton("Option 2", false).
		AddButton("Option 3", false)

	if len(rbg.buttons) != 3 {
		t.Errorf("Expected 3 buttons, got %d", len(rbg.buttons))
	}

	// Test getting selected
	index, label := rbg.GetSelected()
	if index != 0 || label != "Option 1" {
		t.Errorf("Expected selected (0, 'Option 1'), got (%d, '%s')", index, label)
	}

	// Test setting selected
	rbg.SetSelected(2)
	index, label = rbg.GetSelected()
	if index != 2 || label != "Option 3" {
		t.Errorf("Expected selected (2, 'Option 3'), got (%d, '%s')", index, label)
	}

	// Test layout directions
	rbg.SetLayout(LayoutHorizontal)
	if rbg.layout != LayoutHorizontal {
		t.Error("Layout not set correctly")
	}

	rbg.SetLayout(LayoutVertical)
	if rbg.layout != LayoutVertical {
		t.Error("Layout not set correctly")
	}
}

func TestConvenienceMethods(t *testing.T) {
	// Create a mock theme for testing
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Test vertical radio group
	verticalGroup := w.VerticalRadioGroup()
	if verticalGroup.layout != LayoutVertical {
		t.Error("VerticalRadioGroup should have vertical layout")
	}

	// Test horizontal radio group
	horizontalGroup := w.HorizontalRadioGroup()
	if horizontalGroup.layout != LayoutHorizontal {
		t.Error("HorizontalRadioGroup should have horizontal layout")
	}
}

func TestRadioButtonCentering(t *testing.T) {
	// Create a mock theme for testing
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Test creating a checked radio button
	rb := w.NewRadioButton(true)
	if rb == nil {
		t.Fatal("NewRadioButton returned nil")
	}

	// Test that the radio button is checked
	if !rb.IsChecked() {
		t.Error("Radio button should be checked")
	}

	// Test setting and getting checked state
	rb.SetChecked(false)
	if rb.IsChecked() {
		t.Error("Radio button should not be checked after SetChecked(false)")
	}

	rb.SetChecked(true)
	if !rb.IsChecked() {
		t.Error("Radio button should be checked after SetChecked(true)")
	}
}

func TestRadioButtonGroupSelection(t *testing.T) {
	// Create a mock theme for testing
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Create a radio button group
	rbg := w.NewRadioButtonGroup()

	// Add buttons to the group
	rbg.AddButton("Option 1", true). // This should be selected initially
						AddButton("Option 2", false).
						AddButton("Option 3", false)

	// Test initial state
	if len(rbg.buttons) != 3 {
		t.Errorf("Expected 3 buttons, got %d", len(rbg.buttons))
	}

	// Test that only the first button is checked initially
	if !rbg.buttons[0].IsChecked() {
		t.Error("First button should be checked initially")
	}
	if rbg.buttons[1].IsChecked() {
		t.Error("Second button should not be checked initially")
	}
	if rbg.buttons[2].IsChecked() {
		t.Error("Third button should not be checked initially")
	}

	// Test setting selected button
	rbg.SetSelected(2)

	// Test that only the third button is checked after selection change
	if rbg.buttons[0].IsChecked() {
		t.Error("First button should not be checked after selecting third")
	}
	if rbg.buttons[1].IsChecked() {
		t.Error("Second button should not be checked after selecting third")
	}
	if !rbg.buttons[2].IsChecked() {
		t.Error("Third button should be checked after selection")
	}

	// Test getting selected
	index, label := rbg.GetSelected()
	if index != 2 || label != "Option 3" {
		t.Errorf("Expected selected (2, 'Option 3'), got (%d, '%s')", index, label)
	}
}

func TestRadioButtonGroupCallbacks(t *testing.T) {
	// Create a mock theme for testing
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Track callback calls
	var callbackCalls []struct {
		index int
		label string
	}

	// Create a radio button group with callback
	rbg := w.NewRadioButtonGroup().
		SetOnChange(func(index int, label string) {
			callbackCalls = append(callbackCalls, struct {
				index int
				label string
			}{index, label})
		})

	// Add buttons to the group
	rbg.AddButton("Option 1", true). // This should trigger callback
						AddButton("Option 2", false).
						AddButton("Option 3", false)

	// Test that initial selection triggers callback
	if len(callbackCalls) != 1 {
		t.Errorf("Expected 1 callback call, got %d", len(callbackCalls))
	}
	if callbackCalls[0].index != 0 || callbackCalls[0].label != "Option 1" {
		t.Errorf("Expected callback (0, 'Option 1'), got (%d, '%s')",
			callbackCalls[0].index, callbackCalls[0].label)
	}

	// Test setting selected button triggers callback
	rbg.SetSelected(2)

	// Test that selection change triggers callback
	if len(callbackCalls) != 2 {
		t.Errorf("Expected 2 callback calls, got %d", len(callbackCalls))
	}
	if callbackCalls[1].index != 2 || callbackCalls[1].label != "Option 3" {
		t.Errorf("Expected callback (2, 'Option 3'), got (%d, '%s')",
			callbackCalls[1].index, callbackCalls[1].label)
	}
}
