package fromage

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

// Int is a slider widget for selecting integer values
type Int struct {
	value      int
	clickable  *widget.Clickable
	pos        float32 // position normalized to [0, 1]
	length     float32
	changed    bool
	changeHook func(int)
	min, max   int
	dragging   bool
	drag       gesture.Drag
}

// NewInt creates a new integer slider
func (t *Theme) NewInt() *Int {
	return &Int{
		clickable:  t.Pool.GetClickable(),
		changeHook: func(int) {},
		min:        0,
		max:        100,
	}
}

// SetValue sets the current value
func (i *Int) SetValue(value int) *Int {
	i.value = value
	return i
}

// Value returns the current value
func (i *Int) Value() int {
	return i.value
}

// SetRange sets the min and max values
func (i *Int) SetRange(min, max int) *Int {
	i.min = min
	i.max = max
	return i
}

// SetHook sets the change callback
func (i *Int) SetHook(fn func(int)) *Int {
	i.changeHook = fn
	return i
}

// Layout renders the integer slider
func (i *Int) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	// Ensure minimum size
	minSize := gtx.Dp(unit.Dp(200))
	if gtx.Constraints.Min.X < minSize {
		gtx.Constraints.Min.X = minSize
	}
	if gtx.Constraints.Min.Y < gtx.Dp(unit.Dp(40)) {
		gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(40))
	}

	size := gtx.Constraints.Min
	i.length = float32(size.X)

	// Update position based on current value
	if i.min != i.max {
		i.pos = float32(i.value-i.min) / float32(i.max-i.min)
	}

	if i.pos < 0 {
		i.pos = 0
	} else if i.pos > 1 {
		i.pos = 1
	}

	// Handle drag gestures
	for {
		ev, ok := i.drag.Update(gtx.Metric, gtx.Source, gesture.Horizontal)
		if !ok {
			break
		}

		switch ev.Kind {
		case pointer.Press:
			i.dragging = true
			// Update value based on press position
			newPos := float32(ev.Position.X) / i.length
			if newPos < 0 {
				newPos = 0
			} else if newPos > 1 {
				newPos = 1
			}
			i.pos = newPos
			// Convert to integer value
			i.value = i.min + int(i.pos*float32(i.max-i.min)+0.5)
			i.changed = true
			i.changeHook(i.value)
		case pointer.Drag:
			if i.dragging {
				// Update value based on drag position
				newPos := float32(ev.Position.X) / i.length
				if newPos < 0 {
					newPos = 0
				} else if newPos > 1 {
					newPos = 1
				}
				i.pos = newPos
				// Convert to integer value
				i.value = i.min + int(i.pos*float32(i.max-i.min)+0.5)
				i.changed = true
				i.changeHook(i.value)
			}
		case pointer.Release:
			i.dragging = false
		}
	}

	// Register drag gesture area
	area := clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops)
	i.drag.Add(gtx.Ops)
	area.Pop()

	// Draw track
	trackHeight := gtx.Dp(unit.Dp(4))
	trackRect := image.Rectangle{
		Min: image.Pt(0, (size.Y-trackHeight)/2),
		Max: image.Pt(size.X, (size.Y+trackHeight)/2),
	}
	defer clip.RRect{Rect: trackRect, NW: trackHeight / 2, NE: trackHeight / 2, SW: trackHeight / 2, SE: trackHeight / 2}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, th.Colors.OutlineVariant())

	// Draw thumb
	thumbSize := gtx.Dp(unit.Dp(20))
	thumbX := int(i.pos * i.length)
	// Ensure thumb stays within bounds
	if thumbX < thumbSize/2 {
		thumbX = thumbSize / 2
	} else if thumbX > size.X-thumbSize/2 {
		thumbX = size.X - thumbSize/2
	}

	thumbRect := image.Rectangle{
		Min: image.Pt(thumbX-thumbSize/2, (size.Y-thumbSize)/2),
		Max: image.Pt(thumbX+thumbSize/2, (size.Y+thumbSize)/2),
	}
	defer clip.Ellipse(thumbRect).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, th.Colors.Primary())

	return layout.Dimensions{Size: size}
}

// Pos returns the normalized position [0, 1]
func (i *Int) Pos() float32 {
	return i.pos
}

// Changed returns true if the value has changed since last call
func (i *Int) Changed() bool {
	changed := i.changed
	i.changed = false
	return changed
}
