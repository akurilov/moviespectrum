package spectrum

import (
	"errors"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image"
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

func (ctx *Spectrum) Normalize() *Spectrum {
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
	return &Spectrum{ctx.colorResolution, ctx.levelResolution, normalizedLevels}
}

func (ctx *Spectrum) ToImage() (*image.RGBA, error) {
	width := int(ctx.colorResolution)
	height := int(ctx.levelResolution)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var err error = nil
	const s = 1
	for i, l := range ctx.levels {
		h := float64(360*i) / float64(ctx.colorResolution)
		colColor := colorful.Hsl(h, s, l/2)
		for j := 0; j < int(ctx.levelResolution); j++ {
			img.Set(i, j, colColor)
		}
	}
	return img, err
}