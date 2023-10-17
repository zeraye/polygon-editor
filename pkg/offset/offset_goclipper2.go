package offset

import (
	"github.com/epit3d/goclipper2/goclipper2"
	"github.com/zeraye/polygon-editor/pkg/geom"
)

func createOffsetGoclipper2(poly *geom.Polygon, offset float64) []*geom.Polygon {
	var cb goclipper2.ClipperOffsetCallback = func(
		path *goclipper2.Path64,
		path_normals *goclipper2.PathD,
		curr_idx int,
		prev_idx int,
	) float64 {
		// log.Println("my callback is called with params ", curr_idx, prev_idx)
		return offset
	}

	p := goclipper2.NewPath64()
	for _, pt := range poly.Points {
		p.AddPoint(*goclipper2.NewPoint64(int64(pt.X), int64(pt.Y)))
	}

	co := goclipper2.NewClipperoffset(offset, 0, 0, 0)
	co.AddPath64(*p, goclipper2.SquareJoin, goclipper2.PolygonEnd)

	result := co.ExecuteCallback(cb)
	if result.Length() == 0 {
		return nil
	}

	polys := []*geom.Polygon{}

	for j := 0; j < len(result.Lengths()); j++ {
		new_poly := geom.NewPolygon()
		new_poly.IsClosed = true
		for i := int64(0); i < result.Lengths()[j]; i++ {
			pt := result.GetPoint(int64(j), i)
			new_poly.Points = append(new_poly.Points, geom.NewPoint(float64(pt.X()), float64(pt.Y())))
		}
		polys = append(polys, new_poly)
	}
	return polys
}
