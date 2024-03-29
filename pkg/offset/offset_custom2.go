package offset

import (
	"math"

	"github.com/zeraye/polygon-editor/pkg/geom"
)

func normsCustom2(i int, points []*geom.Point) *geom.Point {
	return geom.GetUnitNormal(points[i], points[(i+1)%len(points)])
}

func createOffsetPointCustom2(prev, curr int, points []*geom.Point, offset float64, joinType string) []*geom.Point {
	offset_points := []*geom.Point{}
	sin_a := geom.CrossProduct(normsCustom2(curr, points), normsCustom2(prev, points))
	cos_a := geom.DotProduct(normsCustom2(curr, points), normsCustom2(prev, points))
	if sin_a > 1 {
		sin_a = 1
	} else if sin_a < -1 {
		sin_a = -1
	}

	if math.Abs(offset) <= 1e-12 {
		offset_points = append(offset_points, points[curr])
		return offset_points
	}

	if cos_a > -0.99 && sin_a*offset < 0 {
		offset_points = append(offset_points, geom.GetPerpendic(points[curr], normsCustom2(prev, points), offset))
		offset_points = append(offset_points, points[curr])
		offset_points = append(offset_points, geom.GetPerpendic(points[curr], normsCustom2(curr, points), offset))
	} else if cos_a > 0.999 || joinType == "miter" {
		q := offset / (cos_a + 1)
		offset_points = append(offset_points, geom.NewPoint(
			points[curr].X+(normsCustom2(prev, points).X+normsCustom2(curr, points).X)*q,
			points[curr].Y+(normsCustom2(prev, points).Y+normsCustom2(curr, points).Y)*q,
		))
	}

	return offset_points
}

func createOffsetCustom2(poly *geom.Polygon, offset float64, joinType string) []*geom.Polygon {
	new_poly := geom.NewPolygon()
	new_poly.IsClosed = true
	pts := new_poly.Points
	points := poly.Points
	num_points := len(points)
	// check if offset is reversed
	cos := geom.DotProduct(normsCustom2(1, points), normsCustom2(0, points))
	p := geom.NewPoint(
		points[1].X+(normsCustom2(0, points).X+normsCustom2(1, points).X)*(offset/(cos+1)),
		points[1].Y+(normsCustom2(0, points).Y+normsCustom2(1, points).Y)*(offset/(cos+1)),
	)
	if geom.IsPointInsidePolygon(*poly, *p) {
		offset *= -1
	}
	for curr := range points {
		prev := (curr + num_points - 1) % num_points
		point := createOffsetPointCustom2(prev, curr, points, offset, joinType)
		pts = append(pts, point...)
	}

	// fix self-intersections
	intersectBool := make([]bool, len(pts))
	intersectInt := make([]int, len(pts))
	intersectCurrs := []int{}
	for curr := range pts {
		next := (curr + 1) % len(pts)
		for _curr := range pts {
			_next := (_curr + 1) % len(pts)
			skip := false
			for _, intersectCurr := range intersectCurrs {
				if intersectCurr == _curr {
					skip = true
				}
			}
			if skip {
				continue
			}
			if geom.SegmentsIntersect(pts[curr], pts[next], pts[_curr], pts[_next]) {
				intersectBool[curr] = true
				intersectInt[curr] = _curr
				intersectCurrs = append(intersectCurrs, curr)
			}
		}
	}

	editCounter := 0
	for i := range intersectBool {
		if intersectBool[i] {
			curr, next, _curr, _next := i, i+1, intersectInt[i], intersectInt[i]+1
			curr = (curr + editCounter) % len(pts)
			next = (next + editCounter) % len(pts)
			_curr = (_curr + editCounter) % len(pts)
			_next = (_next + editCounter) % len(pts)
			if curr < 0 || next < 0 || _curr < 0 || _next < 0 {
				continue
			}
			pt := geom.IntersectPoint(pts[curr], pts[next], pts[_curr], pts[_next])
			pts[next] = pt
			pts = append(pts[:next+1], pts[_next:]...)
			editCounter = editCounter - (_curr - next)
		}
	}
	new_poly.Points = pts
	return []*geom.Polygon{new_poly}
}
