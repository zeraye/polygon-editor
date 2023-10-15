package geom

import "errors"

// Structure represents polygonal chain: https://en.wikipedia.org/wiki/Polygonal_chain
type Polygon struct {
	Points   []*Point
	IsClosed bool
}

func NewPolygon(p ...*Point) *Polygon {
	return &Polygon{Points: p}
}

// Add point new_p after point p in polygon points.
// If point p isn't in polygon points throw error
func (poly *Polygon) AddPointAfter(p, new_p *Point) error {
	p_index := -1
	for i, point := range poly.Points {
		if p == point {
			p_index = i
		}
	}
	if p_index == -1 {
		return errors.New("Point p isn't in polygon points")
	}

	new_p_index := p_index + 1
	if new_p_index == len(poly.Points) {
		// append at the end of slice
		poly.Points = append(poly.Points, new_p)
	} else {
		// append inside slice
		poly.Points = append(poly.Points[:new_p_index+1], poly.Points[new_p_index:]...)
		poly.Points[new_p_index] = new_p
	}

	return nil
}

// Remove point from polygon points with order remained
func (poly *Polygon) RemovePoint(p *Point) error {
	p_index := -1
	for i, point := range poly.Points {
		if p == point {
			p_index = i
		}
	}
	if p_index == -1 {
		return errors.New("Point p isn't in polygon points")
	}

	poly.Points = append(poly.Points[:p_index], poly.Points[p_index+1:]...)

	return nil
}

// Get all segments from polygon
func (poly *Polygon) Segments() []Segment {
	segments := make([]Segment, len(poly.Points))
	for i := 0; i < len(poly.Points); i++ {
		segments[i] = NewSegment(poly.Points[i], poly.Points[(i+1)%len(poly.Points)])
	}
	return segments
}

// Check if point is inside polygon
func IsPointInsidePolygon(poly Polygon, p Point) bool {
	if !poly.IsClosed {
		return false
	}

	intersections := 0
	for _, s := range poly.Segments() {
		x1, y1, x2, y2 := s.P0.X, s.P0.Y, s.P1.X, s.P1.Y
		if (p.Y < y1) != (p.Y < y2) && p.X < x1+(x2-x1)*(p.Y-y1)/(y2-y1) {
			intersections++
		}
	}

	return intersections%2 == 1
}
