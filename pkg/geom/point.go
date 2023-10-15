package geom

import (
	"math"
)

type Point struct {
	X, Y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{X: x, Y: y}
}

// Distance between points p0 and p1
func DistanceToPoint(p0, p1 Point) float64 {
	return math.Sqrt(math.Pow(p0.X-p1.X, 2) + math.Pow(p0.Y-p1.Y, 2))
}

// Distance between point p0 and segment (p1, p2) based on: https://stackoverflow.com/a/6853926
func DistanceToSegment(p0 Point, s Segment) float64 {
	p1, p2 := *s.P0, *s.P1
	a := p0.X - p1.X
	b := p0.Y - p1.Y
	c := p2.X - p1.X
	d := p2.Y - p1.Y
	dot := a*c + b*d
	len_sq := c*c + d*d
	param := -1.0

	if len_sq != 0 {
		param = dot / len_sq
	}

	p := *NewPoint(p1.X+param*c, p1.Y+param*d)

	if param < 0 {
		p = p1
	} else if param > 1 {
		p = p2
	}

	x, y := p0.X-p.X, p0.Y-p.Y
	return math.Sqrt(x*x + y*y)
}
