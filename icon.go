package fromage

import (
	"image"
	"image/color"
	"image/draw"

	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/unit"
	"github.com/gio-eui/ivgconv"
	"golang.org/x/exp/shiny/iconvg"
)

// Icon represents an icon widget that can render IconVG data
type Icon struct {
	theme *Theme
	color color.NRGBA
	src   *[]byte
	size  unit.Dp
	// Cached values
	sz       int
	op       paint.ImageOp
	imgSize  int
	imgColor color.NRGBA
}

// IconByColor maps colors to image operations
type IconByColor map[color.NRGBA]paint.ImageOp

// IconBySize maps sizes to color maps
type IconBySize map[float32]IconByColor

// IconCache maps icon sources to size and color combinations
type IconCache map[*[]byte]IconBySize

// NewIcon creates a new icon widget
func (t *Theme) NewIcon() *Icon {
	return &Icon{
		theme: t,
		color: t.Colors.OnSurface(),
		size:  t.TextSize,
	}
}

// Color sets the color of the icon
func (i *Icon) Color(color color.NRGBA) *Icon {
	i.color = color
	return i
}

// Src sets the icon source data (IconVG format)
func (i *Icon) Src(data *[]byte) *Icon {
	if data == nil {
		// Don't set source if data is nil
		return i
	}
	_, err := iconvg.DecodeMetadata(*data)
	if err != nil {
		// Return the icon without setting source if data is invalid
		return i
	}
	i.src = data
	return i
}

// Scale sets the size relative to the theme's text size
func (i *Icon) Scale(scale float32) *Icon {
	i.size = unit.Dp(float32(i.theme.TextSize) * scale)
	return i
}

// Size sets the absolute size of the icon
func (i *Icon) Size(size unit.Dp) *Icon {
	i.size = size
	return i
}

// Layout renders the icon
func (i *Icon) Layout(g C) D {
	if i.src == nil {
		// Return empty dimensions if no source is set
		return D{}
	}

	ico := i.image(g.Dp(i.size))
	ico.Add(g.Ops)
	paint.PaintOp{}.Add(g.Ops)
	return D{Size: ico.Size()}
}

// image creates or retrieves a cached image operation for the icon
func (i *Icon) image(sz int) paint.ImageOp {
	// Check if we have a cached version
	if ico, ok := i.theme.iconCache[i.src]; ok {
		if isz, ok := ico[float32(i.size)]; ok {
			if icl, ok := isz[i.color]; ok {
				return icl
			}
		}
	}

	// Decode the IconVG metadata
	m, err := iconvg.DecodeMetadata(*i.src)
	if err != nil {
		// Return empty image operation if decode fails
		return paint.ImageOp{}
	}

	// Calculate aspect ratio
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{
		X: sz,
		Y: int(float32(sz) * dy / dx),
	}})

	// Create rasterizer
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)

	// Set the color in the palette
	m.Palette[0] = color.RGBA(i.color)

	// Decode the icon
	err = iconvg.Decode(&ico, *i.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	if err != nil {
		// Return empty image operation if decode fails
		return paint.ImageOp{}
	}

	// Create image operation
	operation := paint.NewImageOp(img)

	// Cache the result
	if _, ok := i.theme.iconCache[i.src]; !ok {
		i.theme.iconCache[i.src] = make(IconBySize)
	}
	if _, ok := i.theme.iconCache[i.src][float32(i.size)]; !ok {
		i.theme.iconCache[i.src][float32(i.size)] = make(IconByColor)
	}
	i.theme.iconCache[i.src][float32(i.size)][i.color] = operation

	return operation
}

// Convenience methods for common icon colors

// IconPrimary creates an icon with primary color
func (t *Theme) IconPrimary() *Icon {
	return t.NewIcon().Color(t.Colors.Primary())
}

// IconSecondary creates an icon with secondary color
func (t *Theme) IconSecondary() *Icon {
	return t.NewIcon().Color(t.Colors.Secondary())
}

// IconTertiary creates an icon with tertiary color
func (t *Theme) IconTertiary() *Icon {
	return t.NewIcon().Color(t.Colors.Tertiary())
}

// IconOnSurface creates an icon with on-surface color
func (t *Theme) IconOnSurface() *Icon {
	return t.NewIcon().Color(t.Colors.OnSurface())
}

// IconOnBackground creates an icon with on-background color
func (t *Theme) IconOnBackground() *Icon {
	return t.NewIcon().Color(t.Colors.OnBackground())
}

// IconError creates an icon with error color
func (t *Theme) IconError() *Icon {
	return t.NewIcon().Color(t.Colors.Error())
}

// IconOutline creates an icon with outline color
func (t *Theme) IconOutline() *Icon {
	return t.NewIcon().Color(t.Colors.Outline())
}

// IconSmall creates a small icon (0.75x text size)
func (t *Theme) IconSmall() *Icon {
	return t.NewIcon().Scale(0.75)
}

// IconLarge creates a large icon (1.5x text size)
func (t *Theme) IconLarge() *Icon {
	return t.NewIcon().Scale(1.5)
}

// IconExtraLarge creates an extra large icon (2x text size)
func (t *Theme) IconExtraLarge() *Icon {
	return t.NewIcon().Scale(2.0)
}

// SVGIconCache maps SVG data to IconVG data
type SVGIconCache map[string]*[]byte

// svgIconCache stores converted SVG to IconVG data
var svgIconCache SVGIconCache = make(SVGIconCache)

// NewIconFromSVG creates a new icon from SVG data
func (t *Theme) NewIconFromSVG(svgData string) *Icon {
	// Check if we already have this SVG converted
	if iconvgData, exists := svgIconCache[svgData]; exists {
		return t.NewIcon().Src(iconvgData)
	}

	// Convert SVG to IconVG
	iconvgData, err := convertSVGToIconVG(svgData)
	if err != nil {
		// If conversion fails, return empty icon
		return t.NewIcon()
	}

	// Cache the converted data
	svgIconCache[svgData] = &iconvgData

	return t.NewIcon().Src(&iconvgData)
}

// convertSVGToIconVG converts SVG data to IconVG format using ivgconv
func convertSVGToIconVG(svgData string) ([]byte, error) {
	// Use ivgconv to convert SVG to IconVG
	iconvgData, err := ivgconv.FromContent([]byte(svgData))
	if err != nil {
		return nil, err
	}

	return iconvgData, nil
}

// NewIconFromSVGFile creates a new icon from SVG file content
func (t *Theme) NewIconFromSVGFile(svgContent string) *Icon {
	return t.NewIconFromSVG(svgContent)
}
