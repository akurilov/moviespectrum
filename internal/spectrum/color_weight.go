package spectrum

import (
	"github.com/akurilov/moviespectrum/internal/util"
	"github.com/lucasb-eyer/go-colorful"
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
	var err error = nil
	nr, ng, nb := util.NormalizeRgb(r, g, b)
	h, s, l := colorful.Color{R: nr, G: ng, B: nb}.Hsl()
	weight, err := util.MedianChiSquare(l)
	if err == nil {
		weight = s * weight
	}
	return &ColorWeight{h / 360, weight}, err
}
