package geom

type Segment struct {
	P0, P1 *Point
}

func NewSegment(p0, p1 *Point) Segment {
	return Segment{p0, p1}
}
