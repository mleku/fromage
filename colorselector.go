package fromage

import (
	"fmt"
	"image/color"
	"math"

	"gio.mleku.dev/layout"
	"gio.mleku.dev/unit"
)

// ColorSelector provides tone, hue, and saturation sliders
type ColorSelector struct {
	tone       *Float // Tone slider (0-1)
	hue        *Float // Hue slider (0-1, represents 0-360 degrees)
	saturation *Float // Saturation slider (0-1)
	onChange   func(color.NRGBA)
}

// NewColorSelector creates a new color selector
func (t *Theme) NewColorSelector() *ColorSelector {
	cs := &ColorSelector{
		tone:       t.NewFloat().SetRange(0, 1).SetValue(0.5),
		hue:        t.NewFloat().SetRange(0, 1).SetValue(0), // 0 = red
		saturation: t.NewFloat().SetRange(0, 1).SetValue(1),
		onChange:   func(color.NRGBA) {},
	}

	// Set up change handlers
	cs.tone.SetHook(func(t float32) {
		cs.updateColor()
	})
	cs.hue.SetHook(func(h float32) {
		cs.updateColor()
	})
	cs.saturation.SetHook(func(s float32) {
		cs.updateColor()
	})

	return cs
}

// SetOnChange sets the callback for when color changes
func (cs *ColorSelector) SetOnChange(fn func(color.NRGBA)) *ColorSelector {
	cs.onChange = fn
	return cs
}

// GetColor returns the current selected color
func (cs *ColorSelector) GetColor() color.NRGBA {
	// Convert hue from 0-1 to 0-360 degrees
	hueDegrees := cs.hue.Value() * 360
	return cs.hsvToRgb(hueDegrees, cs.saturation.Value(), cs.tone.Value())
}

// SetColor sets the color from RGB values
func (cs *ColorSelector) SetColor(c color.NRGBA) *ColorSelector {
	h, s, v := cs.rgbToHsv(c)
	cs.hue.SetValue(h / 360) // Convert to 0-1 range
	cs.saturation.SetValue(s)
	cs.tone.SetValue(v)
	return cs
}

// GetTone returns the current tone value (0-1)
func (cs *ColorSelector) GetTone() float32 {
	return cs.tone.Value()
}

// GetHue returns the current hue value (0-1)
func (cs *ColorSelector) GetHue() float32 {
	return cs.hue.Value()
}

// GetSaturation returns the current saturation value (0-1)
func (cs *ColorSelector) GetSaturation() float32 {
	return cs.saturation.Value()
}

// SetTone sets the tone value (0-1)
func (cs *ColorSelector) SetTone(tone float32) *ColorSelector {
	cs.tone.SetValue(tone)
	return cs
}

// SetHue sets the hue value (0-1)
func (cs *ColorSelector) SetHue(hue float32) *ColorSelector {
	cs.hue.SetValue(hue)
	return cs
}

// SetSaturation sets the saturation value (0-1)
func (cs *ColorSelector) SetSaturation(saturation float32) *ColorSelector {
	cs.saturation.SetValue(saturation)
	return cs
}

// Layout renders the color selector
func (cs *ColorSelector) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	return th.VFlex().
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Tone slider
			return cs.createSliderRow(gtx, th, "Tone", cs.tone)
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Spacer
			return layout.Spacer{Height: unit.Dp(5)}.Layout(gtx)
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Hue slider
			return cs.createSliderRow(gtx, th, "Hue", cs.hue)
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Spacer
			return layout.Spacer{Height: unit.Dp(5)}.Layout(gtx)
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Saturation slider
			return cs.createSliderRow(gtx, th, "Saturation", cs.saturation)
		}).
		Layout(gtx)
}

// createSliderRow creates a row with label, slider, and value
func (cs *ColorSelector) createSliderRow(gtx layout.Context, th *Theme, label string, slider *Float) layout.Dimensions {
	return th.HFlex().
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Label
			return th.Caption(label).Layout(gtx)
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Spacer
			return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
		}).
		Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// Slider - constrain to full width minus padding
			textHeight := gtx.Dp(th.TextSize)
			gtx.Constraints.Min.X = gtx.Constraints.Max.X - textHeight*2 // Leave space for padding
			return slider.Layout(gtx, th)
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Spacer
			return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
		}).
		Rigid(func(gtx layout.Context) layout.Dimensions {
			// Value display with 4 decimal places
			value := slider.Value()
			valueStr := fmt.Sprintf("%.4f", value)
			return th.Caption(valueStr).Layout(gtx)
		}).
		Layout(gtx)
}

// hsvToRgb converts HSV to RGB
func (cs *ColorSelector) hsvToRgb(h, s, v float32) color.NRGBA {
	h = h / 360.0 // Normalize to [0, 1]
	c := v * s
	x := c * (1 - float32(math.Abs(float64(math.Mod(float64(h*6), 2))-1)))
	m := v - c

	var r, g, b float32
	if h < 1.0/6.0 {
		r, g, b = c, x, 0
	} else if h < 2.0/6.0 {
		r, g, b = x, c, 0
	} else if h < 3.0/6.0 {
		r, g, b = 0, c, x
	} else if h < 4.0/6.0 {
		r, g, b = 0, x, c
	} else if h < 5.0/6.0 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}

	return color.NRGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}

// rgbToHsv converts RGB to HSV
func (cs *ColorSelector) rgbToHsv(c color.NRGBA) (float32, float32, float32) {
	r := float32(c.R) / 255.0
	g := float32(c.G) / 255.0
	b := float32(c.B) / 255.0

	max := float32(math.Max(float64(r), math.Max(float64(g), float64(b))))
	min := float32(math.Min(float64(r), math.Min(float64(g), float64(b))))
	delta := max - min

	var h float32
	if delta == 0 {
		h = 0
	} else if max == r {
		h = float32(60 * math.Mod(float64((g-b)/delta), 6))
	} else if max == g {
		h = 60 * ((b-r)/delta + 2)
	} else {
		h = 60 * ((r-g)/delta + 4)
	}

	if h < 0 {
		h += 360
	}

	s := float32(0)
	if max != 0 {
		s = delta / max
	}

	v := max

	return h, s, v
}

// updateColor calls the onChange callback with the current color
func (cs *ColorSelector) updateColor() {
	cs.onChange(cs.GetColor())
}
