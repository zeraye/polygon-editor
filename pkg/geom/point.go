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

func CrossProduct(p0, p1 *Point) float64 {
	return p0.Y*p1.X - p1.Y*p0.X
}

func DotProduct(p0, p1 *Point) float64 {
	return p0.X*p1.X + p0.Y*p1.Y
}

func GetUnitNormal(p0, p1 *Point) *Point {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	if dx == 0 && dy == 0 {
		return NewPoint(0, 0)
	}

	f := 1 / math.Sqrt(dx*dx+dy*dy)
	dx *= f
	dy *= f

	return NewPoint(dy, -dx)
}

func NormalizeVector(p *Point) *Point {
	hypot := math.Sqrt(p.X*p.X + p.Y*p.Y)
	if hypot == 0 {
		return NewPoint(0, 0)
	}
	return NewPoint(p.X/hypot, p.Y/hypot)
}

func GetAvgUnitVector(p0, p1 *Point) *Point {
	return NormalizeVector(NewPoint(p0.X+p1.X, p0.Y+p1.Y))
}

func GetPerpendic(p, norm *Point, delta float64) *Point {
	return NewPoint(p.X+norm.X*delta, p.Y+norm.Y*delta)
}

func TranslatePoint(p *Point, dx, dy float64) *Point {
	return NewPoint(p.X+dx, p.Y+dy)
}

func ReflectPoint(p, pivot *Point) *Point {
	return NewPoint(pivot.X+(pivot.X-p.X), pivot.Y+(pivot.Y-p.Y))
}

func IntersectPoint(pt1a, pt1b, pt2a, pt2b *Point) *Point {
	if pt1a.X == pt1b.X {
		if pt2a.X == pt2b.X {
			return NewPoint(0, 0)
		}
		m2 := (pt2b.Y - pt2a.Y) / (pt2b.X - pt2a.X)
		b2 := pt2a.Y - m2*pt2a.X
		return NewPoint(pt1a.X, m2*pt1a.X+b2)
	} else if pt2a.X == pt2b.X {
		m1 := (pt1b.Y - pt1a.Y) / (pt1b.X - pt1a.X)
		b1 := pt1a.Y - m1*pt1a.X
		return NewPoint(pt2a.X, m1*pt2a.X+b1)
	} else {
		m1 := (pt1b.Y - pt1a.Y) / (pt1b.X - pt1a.X)
		b1 := pt1a.Y - m1*pt1a.X
		m2 := (pt2b.Y - pt2a.Y) / (pt2b.X - pt2a.X)
		b2 := pt2a.Y - m2*pt2a.X
		if m1 == m2 {
			return NewPoint(0, 0)
		}
		x := (b2 - b1) / (m1 - m2)
		return NewPoint(x, m1*x+b1)
	}
}
