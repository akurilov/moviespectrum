package spectrum

import (
	"errors"
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/lucasb-eyer/go-colorful"
	"io"
)

const (
	HueRange = 360
)

type Spectrum struct {
	colorResolution uint16
	levelResolution uint16
	levels          []float64
}

func NewSpectrum(colorResolution uint16, levelResolution uint16) *Spectrum {
	levels := make([]float64, colorResolution)
	return &Spectrum{
		colorResolution,
		levelResolution,
		levels,
	}
}

func (ctx *Spectrum) AddMeasurement(measurement *ColorWeight) error {
	var err error = nil
	color := measurement.Color()
	if color < 0 || color >= 1 {
		err = errors.New(
			fmt.Sprintf("the measurement color should be in the range of [0, 1), got %f", color))
	} else {
		i := uint16(float64(ctx.colorResolution) * measurement.Color())
		level := measurement.Weight()
		if level < 0 {
			err = errors.New(fmt.Sprintf("the measurement level should not be less than 0, go %f", level))
		} else {
			ctx.levels[i] = ctx.levels[i] + level
		}
	}
	return err
}

func (ctx *Spectrum) normalize() []float64 {
	maxLevel := 0.0
	for _, level := range ctx.levels {
		if level > maxLevel {
			maxLevel = level
		}
	}
	normalizedLevels := make([]float64, ctx.colorResolution)
	if maxLevel > 0 {
		for i, level := range ctx.levels {
			normalizedLevels[i] = level / maxLevel
		}
	}
	return normalizedLevels
}

func (ctx *Spectrum) ToSvgImage(output io.Writer) {
	svgImg := svg.New(output)
	width := int(ctx.colorResolution)
	height := int(ctx.levelResolution)
	svgImg.Start(width, height)
	normalizedLevels := ctx.normalize()
	const s = 1
	polyLineXs := make([]int, 0)
	polyLineYs := make([]int, 0)
	for i, l := range normalizedLevels {
		h := float64(HueRange*i) / float64(ctx.colorResolution)
		colColor := colorful.Hsl(h, s, l/2)
		svgImg.Rect(i, 0, i+1, height, "fill:"+ToCssColor(&colColor)+";stroke:none")
		polyLineXs = append(polyLineXs, i)
		polyLineYs = append(polyLineYs, int(float64(height)*(1-l)))
	}
	svgImg.Polyline(polyLineXs, polyLineYs, "stroke:rgb(255,255,255,0.5);fill:none")
	svgImg.End()
}

func ToCssColor(c *colorful.Color) string {
	const CssColorRange = 256
	r, g, b := CssColorRange*c.R, CssColorRange*c.G, CssColorRange*c.B
	return fmt.Sprintf("rgba(%d, %d, %d, %f)", int(r), int(g), int(b), 0.5)
}
