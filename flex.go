package fromage

import (
	"gio.mleku.dev/layout"
)

// Flex provides a fluent API for creating flex layouts
type Flex struct {
	flex     layout.Flex
	children []layout.FlexChild
}

// NewFlex creates a new flex layout
func (t *Theme) NewFlex() *Flex {
	return &Flex{
		flex:     layout.Flex{},
		children: make([]layout.FlexChild, 0),
	}
}

// VFlex creates a new vertical flex layout
func (t *Theme) VFlex() *Flex {
	return t.NewFlex().Vertical()
}

// HFlex creates a new horizontal flex layout
func (t *Theme) HFlex() *Flex {
	return t.NewFlex()
}

// Alignment setters

// AlignStart sets alignment for layout from Start
func (f *Flex) AlignStart() *Flex {
	f.flex.Alignment = layout.Start
	return f
}

// AlignEnd sets alignment for layout from End
func (f *Flex) AlignEnd() *Flex {
	f.flex.Alignment = layout.End
	return f
}

// AlignMiddle sets alignment for layout from Middle
func (f *Flex) AlignMiddle() *Flex {
	f.flex.Alignment = layout.Middle
	return f
}

// AlignBaseline sets alignment for layout from Baseline
func (f *Flex) AlignBaseline() *Flex {
	f.flex.Alignment = layout.Baseline
	return f
}

// Axis setters

// Vertical sets axis to vertical, otherwise it is horizontal
func (f *Flex) Vertical() *Flex {
	f.flex.Axis = layout.Vertical
	return f
}

// Horizontal sets axis to horizontal (default)
func (f *Flex) Horizontal() *Flex {
	f.flex.Axis = layout.Horizontal
	return f
}

// Spacing setters

// SpaceStart sets the corresponding flex spacing parameter
func (f *Flex) SpaceStart() *Flex {
	f.flex.Spacing = layout.SpaceStart
	return f
}

// SpaceEnd sets the corresponding flex spacing parameter
func (f *Flex) SpaceEnd() *Flex {
	f.flex.Spacing = layout.SpaceEnd
	return f
}

// SpaceSides sets the corresponding flex spacing parameter
func (f *Flex) SpaceSides() *Flex {
	f.flex.Spacing = layout.SpaceSides
	return f
}

// SpaceAround sets the corresponding flex spacing parameter
func (f *Flex) SpaceAround() *Flex {
	f.flex.Spacing = layout.SpaceAround
	return f
}

// SpaceBetween sets the corresponding flex spacing parameter
func (f *Flex) SpaceBetween() *Flex {
	f.flex.Spacing = layout.SpaceBetween
	return f
}

// SpaceEvenly sets the corresponding flex spacing parameter
func (f *Flex) SpaceEvenly() *Flex {
	f.flex.Spacing = layout.SpaceEvenly
	return f
}

// Child management

// Rigid inserts a rigid widget into the flex
func (f *Flex) Rigid(w W) *Flex {
	f.children = append(f.children, layout.Rigid(w))
	return f
}

// Flexed inserts a flexed widget into the flex
func (f *Flex) Flexed(weight float32, w W) *Flex {
	f.children = append(f.children, layout.Flexed(weight, w))
	return f
}

// Layout renders the flex layout
func (f *Flex) Layout(g C) D {
	return f.flex.Layout(g, f.children...)
}

// Convenience methods for common layouts

// Column creates a vertical flex layout with even spacing
func (t *Theme) Column() *Flex {
	return t.VFlex().SpaceEvenly()
}

// Row creates a horizontal flex layout with even spacing
func (t *Theme) Row() *Flex {
	return t.HFlex().SpaceEvenly()
}

// CenteredColumn creates a vertical flex layout with centered alignment
func (t *Theme) CenteredColumn() *Flex {
	return t.VFlex().AlignMiddle().SpaceEvenly()
}

// CenteredRow creates a horizontal flex layout with centered alignment
func (t *Theme) CenteredRow() *Flex {
	return t.HFlex().AlignMiddle().SpaceEvenly()
}
