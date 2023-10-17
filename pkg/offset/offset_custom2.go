package offset

import (
	"log"

	"github.com/zeraye/polygon-editor/pkg/geom"
)

func normsCustom2(i int, points []*geom.Point) *geom.Point {
	return geom.GetUnitNormal(points[i], points[(i+1)%len(points)])
}

func createOffsetPointCustom2(prev, curr, next int, points []*geom.Point, offset float64) []*geom.Point {
	offset_points := []*geom.Point{}
	sin_a := geom.CrossProduct(normsCustom2(curr, points), normsCustom2(prev, points))
	cos_a := geom.DotProduct(normsCustom2(curr, points), normsCustom2(prev, points))
	if sin_a > 1 {
		sin_a = 1
	} else if sin_a < -1 {
		sin_a = -1
	}

	if offset <= 1e-12 {
		offset_points = append(offset_points, points[curr])
		return offset_points
	}

	if cos_a > -0.99 && sin_a*offset < 0 {
		log.Println("concave")
		offset_points = append(offset_points, geom.GetPerpendic(points[curr], normsCustom2(prev, points), offset))
		offset_points = append(offset_points, points[curr])
		offset_points = append(offset_points, geom.GetPerpendic(points[curr], normsCustom2(curr, points), offset))
	} else if cos_a > 0.999 {
		log.Println("miter")
		q := offset / (cos_a + 1)
		offset_points = append(offset_points, geom.NewPoint(
			points[curr].X+(normsCustom2(prev, points).X+normsCustom2(curr, points).X)*q,
			points[curr].Y+(normsCustom2(prev, points).Y+normsCustom2(curr, points).Y)*q,
		))
	} else if cos_a > 0.99 {
		log.Println("bevel")
		pt1 := geom.NewPoint(points[curr].X+offset*normsCustom2(prev, points).X, points[curr].Y+offset*normsCustom2(prev, points).Y)
		pt2 := geom.NewPoint(points[curr].X+offset*normsCustom2(curr, points).X, points[curr].Y+offset*normsCustom2(curr, points).Y)
		offset_points = append(offset_points, pt1, pt2)
	} else {
		log.Println("square")
		vec := geom.GetAvgUnitVector(
			geom.NewPoint(-normsCustom2(prev, points).Y, normsCustom2(prev, points).X),
			geom.NewPoint(normsCustom2(curr, points).Y, -normsCustom2(curr, points).X),
		)
		ptQ := geom.TranslatePoint(points[curr], offset*vec.X, offset*vec.Y)
		pt1 := geom.TranslatePoint(ptQ, offset*vec.Y, offset*(-vec.X))
		pt2 := geom.TranslatePoint(ptQ, offset*(-vec.Y), offset*vec.X)
		pt3 := geom.GetPerpendic(points[curr], normsCustom2(prev, points), offset)
		pt4 := geom.GetPerpendic(points[curr], normsCustom2(prev, points), offset)
		pt := geom.IntersectPoint(pt1, pt2, pt3, pt4)
		offset_points = append(offset_points, pt)
		offset_points = append(offset_points, geom.ReflectPoint(pt, ptQ))

	}

	return offset_points
}

func createOffsetCustom2(poly *geom.Polygon, offset float64) []*geom.Polygon {
	new_poly := geom.NewPolygon()
	new_poly.IsClosed = true
	points := poly.Points
	num_points := len(points)
	for curr := range points {
		prev := (curr + num_points - 1) % num_points
		next := (curr + 1) % num_points
		point := createOffsetPointCustom2(prev, curr, next, points, offset)
		new_poly.Points = append(new_poly.Points, point...)
	}
	return []*geom.Polygon{new_poly}
}
