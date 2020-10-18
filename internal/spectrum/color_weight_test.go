package spectrum

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRgbToHslToRgb(t *testing.T) {
	{
		h, s, l := rgbToHsl(0, 255, 0)
		c := colorful.Hsl(h, s, l)
		assert.InDelta(t, c.R, 0, 0.01)
		assert.InDelta(t, c.G, 1, 0.01)
		assert.InDelta(t, c.B, 0, 0.01)
	}
	{
		h, s, l := rgbToHsl(255, 0, 0)
		c := colorful.Hsl(h, s, l)
		assert.InDelta(t, c.R, 1, 0.01)
		assert.InDelta(t, c.G, 0, 0.01)
		assert.InDelta(t, c.B, 0, 0.01)
	}
	{
		h, s, l := rgbToHsl(0, 0, 255)
		c := colorful.Hsl(h, s, l)
		assert.InDelta(t, c.R, 0, 0.01)
		assert.InDelta(t, c.G, 0, 0.01)
		assert.InDelta(t, c.B, 1, 0.01)
	}
	{
		h, s, l := rgbToHsl(128, 0, 0)
		c := colorful.Hsl(h, s, l)
		assert.InDelta(t, c.R, 0.5, 0.01)
		assert.InDelta(t, c.G, 0, 0.01)
		assert.InDelta(t, c.B, 0, 0.01)
	}
	{
		h, s, l := rgbToHsl(0, 128, 0)
		c := colorful.Hsl(h, s, l)
		assert.InDelta(t, c.R, 0, 0.01)
		assert.InDelta(t, c.G, 0.5, 0.01)
		assert.InDelta(t, c.B, 0, 0.01)
	}
	{
		h, s, l := rgbToHsl(0, 128, 128)
		c := colorful.Hsl(h, s, l)
		assert.InDelta(t, c.R, 0, 0.01)
		assert.InDelta(t, c.G, 0.5, 0.01)
		assert.InDelta(t, c.B, 0.5, 0.01)
	}
	{
		h, s, l := rgbToHsl(0, 0, 0)
		c := colorful.Hsl(h, s, l)
		assert.InDelta(t, c.R, 0, 0.01)
		assert.InDelta(t, c.G, 0, 0.01)
		assert.InDelta(t, c.B, 0, 0.01)
	}
}
