package spectrum

import (
	"github.com/lucasb-eyer/go-colorful"
	"math"
)

type ColorWeight struct {
	color  float64
	weight float64
}

func (ctx *ColorWeight) Color() float64 {
	return ctx.color
}

func (ctx *ColorWeight) Weight() float64 {
	return ctx.weight
}

func NewColorWeight(r, g, b uint8) (*ColorWeight, error) {
	const medianValue = 0.5
	const hueValueRange = 360
	var err error = nil
	h, s, l := rgbToHsl(r, g, b)
	weight := s * chiSquare(l, medianValue)
	return &ColorWeight{h / hueValueRange, weight}, err
}

func rgbToHsl(r, g, b uint8) (float64, float64, float64) {
	const uint8ValueRange = 0x100
	normalizedRgb := colorful.Color{
		R: float64(r) / uint8ValueRange,
		G: float64(g) / uint8ValueRange,
		B: float64(b) / uint8ValueRange,
	}
	return normalizedRgb.Hsl()
}

func chiSquare(x float64, expected float64) float64 {
	expectedSquare := math.Pow(expected, 2)
	return (expectedSquare - math.Pow(expected-x, 2)) / expectedSquare
}
