package draw

import (
	"image"
	"image/color"
)

func Line(img *image.RGBA, c *color.RGBA, start, end *image.Point) {
	x0, y0, x1, y1 := start.X, start.Y, end.X, end.Y
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy
	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}
