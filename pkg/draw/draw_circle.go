package draw

import (
	"image"
	"image/color"

	"github.com/zeraye/polygon-editor/pkg/geom"
)

// Midpoint circle algorithm for drawing circle: https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func DrawCircle(centre geom.Point, radius float64, color color.Color, fill bool, img *image.RGBA) {
	x := radius
	y := 0.0
	err := 0.0
	for x >= y {
		drawCircleSet(centre.X+x, centre.Y+y, centre.X-x, centre.Y+y, color, fill, img)
		drawCircleSet(centre.X+y, centre.Y+x, centre.X-y, centre.Y+x, color, fill, img)
		drawCircleSet(centre.X+x, centre.Y-y, centre.X-x, centre.Y-y, color, fill, img)
		drawCircleSet(centre.X+y, centre.Y-x, centre.X-y, centre.Y-x, color, fill, img)

		if err <= 0 {
			y++
			err += 2*y + 1
		} else {
			x--
			err += -2*x + 1
		}
	}

}

func drawCircleSet(x0, y0, x1, y1 float64, color color.Color, fill bool, img *image.RGBA) {
	if fill {
		BresenhamDrawLine(*geom.NewPoint(x0, y0), *geom.NewPoint(x1, y1), color, img, 1)
	} else {
		img.Set(int(x0), int(y0), color)
		img.Set(int(x1), int(y1), color)
	}
}
