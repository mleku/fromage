package fromage

import (
	"context"
	"testing"

	"gioui.org/layout"
)

func TestFlex(t *testing.T) {
	// Create a theme
	th := NewThemeWithMode(
		context.Background(), // context not needed for this test
		NewColors,
		nil, // shaper not needed for this test
		16,  // text size
		ThemeModeLight,
	)

	// Test basic flex creation
	flex := th.NewFlex()
	if flex == nil {
		t.Error("NewFlex should return a non-nil flex")
	}

	// Test vertical flex
	vflex := th.VFlex()
	if vflex.flex.Axis != layout.Vertical {
		t.Error("VFlex should create a vertical flex")
	}

	// Test horizontal flex
	hflex := th.HFlex()
	if hflex.flex.Axis != layout.Horizontal {
		t.Error("HFlex should create a horizontal flex")
	}
}

func TestFlexAlignment(t *testing.T) {
	th := NewThemeWithMode(context.Background(), NewColors, nil, 16, ThemeModeLight)

	// Test alignment methods
	flex := th.NewFlex()

	// Test AlignStart
	flex.AlignStart()
	if flex.flex.Alignment != layout.Start {
		t.Error("AlignStart should set alignment to Start")
	}

	// Test AlignMiddle
	flex.AlignMiddle()
	if flex.flex.Alignment != layout.Middle {
		t.Error("AlignMiddle should set alignment to Middle")
	}

	// Test AlignEnd
	flex.AlignEnd()
	if flex.flex.Alignment != layout.End {
		t.Error("AlignEnd should set alignment to End")
	}

	// Test AlignBaseline
	flex.AlignBaseline()
	if flex.flex.Alignment != layout.Baseline {
		t.Error("AlignBaseline should set alignment to Baseline")
	}
}

func TestFlexSpacing(t *testing.T) {
	th := NewThemeWithMode(context.Background(), NewColors, nil, 16, ThemeModeLight)

	// Test spacing methods
	flex := th.NewFlex()

	// Test SpaceStart
	flex.SpaceStart()
	if flex.flex.Spacing != layout.SpaceStart {
		t.Error("SpaceStart should set spacing to SpaceStart")
	}

	// Test SpaceEvenly
	flex.SpaceEvenly()
	if flex.flex.Spacing != layout.SpaceEvenly {
		t.Error("SpaceEvenly should set spacing to SpaceEvenly")
	}

	// Test SpaceBetween
	flex.SpaceBetween()
	if flex.flex.Spacing != layout.SpaceBetween {
		t.Error("SpaceBetween should set spacing to SpaceBetween")
	}

	// Test SpaceAround
	flex.SpaceAround()
	if flex.flex.Spacing != layout.SpaceAround {
		t.Error("SpaceAround should set spacing to SpaceAround")
	}

	// Test SpaceSides
	flex.SpaceSides()
	if flex.flex.Spacing != layout.SpaceSides {
		t.Error("SpaceSides should set spacing to SpaceSides")
	}

	// Test SpaceEnd
	flex.SpaceEnd()
	if flex.flex.Spacing != layout.SpaceEnd {
		t.Error("SpaceEnd should set spacing to SpaceEnd")
	}
}

func TestFlexChildren(t *testing.T) {
	th := NewThemeWithMode(context.Background(), NewColors, nil, 16, ThemeModeLight)

	// Test adding children
	flex := th.NewFlex()

	// Create a dummy widget
	dummyWidget := func(g C) D {
		return D{}
	}

	// Test Rigid
	flex.Rigid(dummyWidget)
	if len(flex.children) != 1 {
		t.Errorf("Expected 1 child after Rigid, got %d", len(flex.children))
	}

	// Test Flexed
	flex.Flexed(0.5, dummyWidget)
	if len(flex.children) != 2 {
		t.Errorf("Expected 2 children after Flexed, got %d", len(flex.children))
	}
}

func TestFlexConvenienceMethods(t *testing.T) {
	th := NewThemeWithMode(context.Background(), NewColors, nil, 16, ThemeModeLight)

	// Test Column
	column := th.Column()
	if column.flex.Axis != layout.Vertical {
		t.Error("Column should create a vertical flex")
	}
	if column.flex.Spacing != layout.SpaceEvenly {
		t.Error("Column should use SpaceEvenly spacing")
	}

	// Test Row
	row := th.Row()
	if row.flex.Axis != layout.Horizontal {
		t.Error("Row should create a horizontal flex")
	}
	if row.flex.Spacing != layout.SpaceEvenly {
		t.Error("Row should use SpaceEvenly spacing")
	}

	// Test CenteredColumn
	centeredColumn := th.CenteredColumn()
	if centeredColumn.flex.Axis != layout.Vertical {
		t.Error("CenteredColumn should create a vertical flex")
	}
	if centeredColumn.flex.Alignment != layout.Middle {
		t.Error("CenteredColumn should use Middle alignment")
	}
	if centeredColumn.flex.Spacing != layout.SpaceEvenly {
		t.Error("CenteredColumn should use SpaceEvenly spacing")
	}

	// Test CenteredRow
	centeredRow := th.CenteredRow()
	if centeredRow.flex.Axis != layout.Horizontal {
		t.Error("CenteredRow should create a horizontal flex")
	}
	if centeredRow.flex.Alignment != layout.Middle {
		t.Error("CenteredRow should use Middle alignment")
	}
	if centeredRow.flex.Spacing != layout.SpaceEvenly {
		t.Error("CenteredRow should use SpaceEvenly spacing")
	}
}
