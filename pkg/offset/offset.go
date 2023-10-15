package offset

import (
	"math"

	"github.com/zeraye/polygon-editor/pkg/geom"
)

func normalizeVec(x, y float64) (float64, float64) {
	distance := math.Sqrt(x*x + y*y)
	return x / distance, y / distance
}

func CreateOffset(poly *geom.Polygon, offset float64) *geom.Polygon {
	new_poly := geom.NewPolygon()
	// return new_poly
	new_poly.IsClosed = true
	outer_ccw := 1.0
	points := poly.Points
	num_points := len(points)
	for curr := range points {
		prev := (curr + num_points - 1) % num_points
		next := (curr + 1) % num_points

		vnX := float64(points[next].X - points[curr].X)
		vnY := float64(points[next].Y - points[curr].Y)
		vnnX, vnnY := normalizeVec(vnX, vnY)
		nnnX := vnnY
		nnnY := -vnnX

		vpX := float64(points[curr].X - points[prev].X)
		vpY := float64(points[curr].Y - points[prev].Y)
		vpnX, vpnY := normalizeVec(vpX, vpY)
		npnX := vpnY * outer_ccw
		npnY := -vpnX * outer_ccw

		bisX := (nnnX + npnX) * outer_ccw
		bisY := (nnnY + npnY) * outer_ccw

		bisnX, bisnY := normalizeVec(bisX, bisY)
		bislen := offset / math.Sqrt((1+nnnX*npnX+nnnY*npnY)/2)

		pnt := geom.NewPoint(points[curr].X+bislen*bisnX, points[curr].Y+bislen*bisnY)
		if geom.IsPointInsidePolygon(*poly, *pnt) {
			pnt = geom.NewPoint(points[curr].X-bislen*bisnX, points[curr].Y-bislen*bisnY)
		}
		new_poly.Points = append(new_poly.Points, pnt)
	}

	return new_poly
}
