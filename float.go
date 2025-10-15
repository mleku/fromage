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

// Float is a slider widget for selecting float values
type Float struct {
	value      float32
	clickable  *widget.Clickable
	pos        float32 // position normalized to [0, 1]
	length     float32
	changed    bool
	changeHook func(float32)
	min, max   float32
	dragging   bool
	drag       gesture.Drag
}

// NewFloat creates a new float slider
func (t *Theme) NewFloat() *Float {
	return &Float{
		clickable:  t.Pool.GetClickable(),
		changeHook: func(float32) {},
		min:        0,
		max:        1,
	}
}

// SetValue sets the current value
func (f *Float) SetValue(value float32) *Float {
	f.value = value
	return f
}

// Value returns the current value
func (f *Float) Value() float32 {
	return f.value
}

// SetRange sets the min and max values
func (f *Float) SetRange(min, max float32) *Float {
	f.min = min
	f.max = max
	return f
}

// SetHook sets the change callback
func (f *Float) SetHook(fn func(float32)) *Float {
	f.changeHook = fn
	return f
}

// Layout renders the float slider
func (f *Float) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	// Ensure minimum size
	minSize := gtx.Dp(unit.Dp(200))
	if gtx.Constraints.Min.X < minSize {
		gtx.Constraints.Min.X = minSize
	}
	if gtx.Constraints.Min.Y < gtx.Dp(unit.Dp(40)) {
		gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(40))
	}

	size := gtx.Constraints.Min
	f.length = float32(size.X)

	// Update position based on current value
	if f.min != f.max {
		f.pos = (f.value - f.min) / (f.max - f.min)
	}

	if f.pos < 0 {
		f.pos = 0
	} else if f.pos > 1 {
		f.pos = 1
	}

	// Handle drag gestures
	for {
		ev, ok := f.drag.Update(gtx.Metric, gtx.Source, gesture.Horizontal)
		if !ok {
			break
		}

		switch ev.Kind {
		case pointer.Press:
			f.dragging = true
			// Update value based on press position
			newPos := float32(ev.Position.X) / f.length
			if newPos < 0 {
				newPos = 0
			} else if newPos > 1 {
				newPos = 1
			}
			f.pos = newPos
			f.value = f.min + f.pos*(f.max-f.min)
			f.changed = true
			f.changeHook(f.value)
		case pointer.Drag:
			if f.dragging {
				// Update value based on drag position
				newPos := float32(ev.Position.X) / f.length
				if newPos < 0 {
					newPos = 0
				} else if newPos > 1 {
					newPos = 1
				}
				f.pos = newPos
				f.value = f.min + f.pos*(f.max-f.min)
				f.changed = true
				f.changeHook(f.value)
			}
		case pointer.Release:
			f.dragging = false
		}
	}

	// Register drag gesture area
	area := clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops)
	f.drag.Add(gtx.Ops)
	area.Pop()

	// Draw track
	trackHeight := gtx.Dp(unit.Dp(4))
	centerY := size.Y / 2
	trackRect := image.Rectangle{
		Min: image.Pt(0, centerY-trackHeight/2),
		Max: image.Pt(size.X, centerY+trackHeight/2),
	}
	defer clip.RRect{Rect: trackRect, NW: trackHeight / 2, NE: trackHeight / 2, SW: trackHeight / 2, SE: trackHeight / 2}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, th.Colors.OutlineVariant())

	// Draw thumb
	thumbSize := gtx.Dp(unit.Dp(20))
	thumbX := int(f.pos * f.length)
	// Ensure thumb stays within bounds
	if thumbX < thumbSize/2 {
		thumbX = thumbSize / 2
	} else if thumbX > size.X-thumbSize/2 {
		thumbX = size.X - thumbSize/2
	}

	thumbRect := image.Rectangle{
		Min: image.Pt(thumbX-thumbSize/2, centerY-thumbSize/2),
		Max: image.Pt(thumbX+thumbSize/2, centerY+thumbSize/2),
	}
	defer clip.Ellipse(thumbRect).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, th.Colors.Primary())

	return layout.Dimensions{Size: size}
}

// Pos returns the normalized position [0, 1]
func (f *Float) Pos() float32 {
	return f.pos
}

// Changed returns true if the value has changed since last call
func (f *Float) Changed() bool {
	changed := f.changed
	f.changed = false
	return changed
}
