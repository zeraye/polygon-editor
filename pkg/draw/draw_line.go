package draw

import (
	"image"
	"image/color"
	"math"

	"github.com/zeraye/polygon-editor/pkg/geom"
)

func ipart(x float64) float64 {
	return float64(int(x))
}

func round(x float64) float64 {
	return float64(int(x + 0.5))
}

func fpart(x float64) float64 {
	return x - ipart(x)
}

func rfpart(x float64) float64 {
	return 1 - fpart(x)
}

func plot(x, y float64, c float64, color color.Color, img *image.RGBA) {
	r, g, b, a := color.RGBA()
	new_color := RGBAToColor([4]uint8{uint8(r), uint8(g), uint8(b), uint8(float64(a) * c)})
	img.Set(int(x), int(y), new_color)
}

func DrawLine(p0, p1 geom.Point, color color.Color, img *image.RGBA, drawLineAlgorithm string, w float64) {
	if drawLineAlgorithm == "bresenham" {
		BresenhamDrawLine(p0, p1, color, img, w)
	} else if drawLineAlgorithm == "xiaolinwu" {
		XiaolinWuDrawLine(p0, p1, color, img, w)
	} else if drawLineAlgorithm == "bresenhamsymmetric" {
		BresenhamDrawLineSymmetric(p0, p1, color, img, w)
	} else {
		panic("Invalid draw line algorithm")
	}
}

// Xiaolin Wu's line algorithm: https://en.wikipedia.org/wiki/Xiaolin_Wu%27s_line_algorithm
func XiaolinWuDrawLine(p0, p1 geom.Point, color color.Color, img *image.RGBA, _ float64) {
	steep := math.Abs(p1.Y-p0.Y) > math.Abs(p1.X-p0.X)

	if steep {
		p0.X, p0.Y = p0.Y, p0.X
		p1.X, p1.Y = p1.Y, p1.X
	}
	if p0.X > p1.X {
		p0, p1 = p1, p0
	}

	dx := p1.X - p0.X
	dy := p1.Y - p0.Y

	gradient := 1.0
	if dx != 0 {
		gradient = dy / dx
	}

	xend := round(p0.X)
	yend := p0.Y + gradient*(xend-p0.X)
	xgap := rfpart(p0.X + 0.5)
	xpxl1 := xend
	ypxl1 := ipart(yend)

	if steep {
		plot(ypxl1, xpxl1, rfpart(yend)*xgap, color, img)
		plot(ypxl1+1, xpxl1, fpart(yend)*xgap, color, img)
	} else {
		plot(xpxl1, ypxl1, rfpart(yend)*xgap, color, img)
		plot(xpxl1, ypxl1+1, fpart(yend)*xgap, color, img)
	}

	intery := yend + gradient

	xend = round(p1.X)
	yend = p1.Y + gradient*(xend-p1.X)
	xgap = fpart(p1.X + 0.5)
	xpxl2 := xend
	ypxl2 := ipart(yend)

	if steep {
		plot(ypxl2, xpxl2, rfpart(yend)*xgap, color, img)
		plot(ypxl2+1, xpxl2, fpart(yend)*xgap, color, img)
	} else {
		plot(xpxl2, ypxl2, rfpart(yend)*xgap, color, img)
		plot(xpxl2, ypxl2+1, fpart(yend)*xgap, color, img)
	}

	if steep {
		for x := xpxl1 + 1; x < xpxl2-1; x++ {
			plot(ipart(intery), x, rfpart(intery), color, img)
			plot(ipart(intery)+1, x, fpart(intery), color, img)
			intery += gradient
		}
	} else {
		for x := xpxl1 + 1; x < xpxl2-1; x++ {
			plot(x, ipart(intery), rfpart(intery), color, img)
			plot(x, ipart(intery)+1, fpart(intery), color, img)
			intery += gradient
		}
	}

}

func BresenhamDrawLineSymmetric(p0, p1 geom.Point, color color.Color, img *image.RGBA, w float64) {
	mid_p := geom.NewPoint((p0.X+p1.X)/2, (p0.Y+p1.Y)/2)
	BresenhamDrawLine(p0, *mid_p, color, img, w)
	BresenhamDrawLine(*mid_p, p1, color, img, w)
}

// Bresenham's line algorithm: https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func BresenhamDrawLine(p0, p1 geom.Point, color color.Color, img *image.RGBA, w float64) {
	if math.Abs(p1.Y-p0.Y) < math.Abs(p1.X-p0.X) {
		if p0.X > p1.X {
			DrawLineLow(p1, p0, color, img, w)
		} else {
			DrawLineLow(p0, p1, color, img, w)
		}
	} else {
		if p0.Y > p1.Y {
			DrawLineHigh(p1, p0, color, img, w)
		} else {
			DrawLineHigh(p0, p1, color, img, w)
		}
	}
}

func DrawLineLow(p0, p1 geom.Point, color color.Color, img *image.RGBA, w float64) {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	yi := 1.0
	if dy < 0 {
		yi = -1
		dy = -dy
	}
	D := 2*dy - dx
	y := p0.Y

	for x := p0.X; x < p1.X; x++ {
		for i := 0; i < int(w); i++ {
			img.Set(int(x), int(y)+i-int(w/2), color)
		}
		if D > 0 {
			y += yi
			D += 2 * (dy - dx)
		} else {
			D += 2 * dy
		}
	}
}

func DrawLineHigh(p0, p1 geom.Point, color color.Color, img *image.RGBA, w float64) {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	xi := 1.0
	if dx < 0 {
		xi = -1
		dx = -dx
	}
	D := 2*dx - dy
	x := p0.X

	for y := p0.Y; y < p1.Y; y++ {
		for i := 0; i < int(w); i++ {
			img.Set(int(x)+i-int(w/2), int(y), color)
		}
		if D > 0 {
			x += xi
			D += 2 * (dx - dy)
		} else {
			D += 2 * dx
		}
	}
}
