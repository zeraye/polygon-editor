package offset

import (
	"math"

	"github.com/zeraye/polygon-editor/pkg/geom"
)

func normalizeVecCustom1(x, y float64) (float64, float64) {
	distance := math.Sqrt(x*x + y*y)
	return x / distance, y / distance
}

func createOffsetPointCustom1(prev, curr, next int, points []*geom.Point, offset float64, poly *geom.Polygon) *geom.Point {
	vnX := float64(points[next].X - points[curr].X)
	vnY := float64(points[next].Y - points[curr].Y)
	vnnX, vnnY := normalizeVecCustom1(vnX, vnY)
	nnnX := vnnY
	nnnY := -vnnX

	vpX := float64(points[curr].X - points[prev].X)
	vpY := float64(points[curr].Y - points[prev].Y)
	vpnX, vpnY := normalizeVecCustom1(vpX, vpY)
	npnX := vpnY
	npnY := -vpnX

	bisX := (nnnX + npnX)
	bisY := (nnnY + npnY)

	bisnX, bisnY := normalizeVecCustom1(bisX, bisY)
	bislen := offset / math.Sqrt((1+nnnX*npnX+nnnY*npnY)/2)

	pnt := geom.NewPoint(points[curr].X+bislen*bisnX, points[curr].Y+bislen*bisnY)
	if geom.IsPointInsidePolygon(*poly, *pnt) {
		pnt = geom.NewPoint(points[curr].X-bislen*bisnX, points[curr].Y-bislen*bisnY)
	}
	return pnt
}

func fixInvalidOffsetSegment(points []*geom.Point) {
}

func createOffsetCustom1(poly *geom.Polygon, offset float64) []*geom.Polygon {
	new_poly := geom.NewPolygon()
	new_poly.IsClosed = true
	points := poly.Points
	num_points := len(points)
	for curr := range points {
		prev := (curr + num_points - 1) % num_points
		next := (curr + 1) % num_points
		pnt := createOffsetPointCustom1(prev, curr, next, points, offset, poly)
		new_poly.Points = append(new_poly.Points, pnt)
		if curr > 0 {
			// TODO: check validity of segment new_poly points [curr-1] and [curr]
			fixInvalidOffsetSegment(new_poly.Points)
		}
	}

	return []*geom.Polygon{new_poly}
}
