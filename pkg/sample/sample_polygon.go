package sample

import (
	"github.com/zeraye/polygon-editor/pkg/constraint"
	"github.com/zeraye/polygon-editor/pkg/geom"
)

func GenerateSampleSimplePolygon() ([]*geom.Polygon, []*constraint.SegmentConstraint) {
	poly := geom.NewPolygon()
	poly.IsClosed = true
	p0 := geom.NewPoint(50, 50)
	p1 := geom.NewPoint(100, 50)
	p2 := geom.NewPoint(100, 300)
	poly.Points = append(poly.Points, p0, p1, p2)
	c0 := constraint.NewSegmentConstraint(p0, p1, constraint.HORIZONTAL)
	c1 := constraint.NewSegmentConstraint(p1, p2, constraint.VERTICAL)
	return []*geom.Polygon{poly}, []*constraint.SegmentConstraint{c0, c1}
}

func GenerateSampleComplexPolygon() ([]*geom.Polygon, []*constraint.SegmentConstraint) {
	poly := geom.NewPolygon()
	poly.IsClosed = true
	p0 := geom.NewPoint(179, 591)
	p1 := geom.NewPoint(267, 321)
	p2 := geom.NewPoint(78, 120)
	p3 := geom.NewPoint(298, 233)
	p4 := geom.NewPoint(322, 124)
	p5 := geom.NewPoint(677, 124)
	p6 := geom.NewPoint(780, 164)
	p7 := geom.NewPoint(762, 187)
	p8 := geom.NewPoint(396, 338)
	p9 := geom.NewPoint(593, 571)
	p10 := geom.NewPoint(380, 674)
	poly.Points = append(poly.Points, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10)
	c0 := constraint.NewSegmentConstraint(p4, p5, constraint.HORIZONTAL)
	return []*geom.Polygon{poly}, []*constraint.SegmentConstraint{c0}
}

func GenerateSampleComplex2Polygon() ([]*geom.Polygon, []*constraint.SegmentConstraint) {
	poly := geom.NewPolygon()
	poly.IsClosed = true
	p0 := geom.NewPoint(109, 429)
	p1 := geom.NewPoint(423, 443)
	p2 := geom.NewPoint(409, 674)
	p3 := geom.NewPoint(252, 670)
	p4 := geom.NewPoint(258, 616)
	p5 := geom.NewPoint(361, 617)
	p6 := geom.NewPoint(361, 477)
	p7 := geom.NewPoint(156, 466)
	p8 := geom.NewPoint(147, 612)
	p9 := geom.NewPoint(213, 612)
	p10 := geom.NewPoint(210, 668)
	p11 := geom.NewPoint(91, 660)
	poly.Points = append(poly.Points, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11)
	return []*geom.Polygon{poly}, []*constraint.SegmentConstraint{}
}

func GenerateSamplePolygons() ([]*geom.Polygon, []*constraint.SegmentConstraint) {
	polygons1, constraints1 := GenerateSampleSimplePolygon()
	polygons2, constraints2 := GenerateSampleComplexPolygon()
	polygons3, constraints3 := GenerateSampleComplex2Polygon()
	return append(polygons1, append(polygons2, polygons3...)...), append(constraints1, append(constraints2, constraints3...)...)
}
