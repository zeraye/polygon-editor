package constraint

import (
	"github.com/zeraye/polygon-editor/pkg/geom"
)

type Constraint uint8

const (
	NO_CONSTRAINT Constraint = iota
	HORIZONTAL
	VERTICAL
)

type SegmentConstraint struct {
	P0, P1 *geom.Point
	Cstr   Constraint
}

func NewSegmentConstraint(p0, p1 *geom.Point, constraint Constraint) *SegmentConstraint {
	return &SegmentConstraint{p0, p1, constraint}
}

var Constraints = []Constraint{
	NO_CONSTRAINT,
	HORIZONTAL,
	VERTICAL,
}

var StringToConstraint = map[string]Constraint{
	"NO_CONSTRAINT": NO_CONSTRAINT,
	"HORIZONTAL":    HORIZONTAL,
	"VERTICAL":      VERTICAL,
}

var ConstraintToString = map[Constraint]string{
	NO_CONSTRAINT: "NO_CONSTRAINT",
	HORIZONTAL:    "HORIZONTAL",
	VERTICAL:      "VERTICAL",
}

// Force polygons to obey all constraints.
// Invalid (misplaced) points will be fixed
func FixSegmentConstraint(constraints []*SegmentConstraint, polygons []*geom.Polygon, p *geom.Point, moveOverlapPointLength float64, allowMoveOverlapPoint bool) {
	for _, segConstraint := range constraints {
		// fix invalid point
		if segConstraint.P0 == p {
			if segConstraint.Cstr == HORIZONTAL {
				segConstraint.P1.Y = p.Y
			} else if segConstraint.Cstr == VERTICAL {
				segConstraint.P1.X = p.X
			}
		} else if segConstraint.P1 == p {
			if segConstraint.Cstr == HORIZONTAL {
				segConstraint.P0.Y = p.Y
			} else if segConstraint.Cstr == VERTICAL {
				segConstraint.P0.X = p.X
			}
		} else {
			// if we didn't fix point in this iteration
			// don't go after if-statement
			continue
		}
		// check if moving overlapped points is necessary
		if !allowMoveOverlapPoint {
			continue
		}
		// if two points in segment overlap
		if *segConstraint.P0 == *segConstraint.P1 {
			// move one point (not p) with remained constrain
			if segConstraint.P0 == p {
				if segConstraint.Cstr == HORIZONTAL {
					segConstraint.P1.X += moveOverlapPointLength
				} else if segConstraint.Cstr == VERTICAL {
					segConstraint.P1.Y += moveOverlapPointLength
				}
			} else if segConstraint.P1 == p {
				if segConstraint.Cstr == HORIZONTAL {
					segConstraint.P0.X += moveOverlapPointLength
				} else if segConstraint.Cstr == VERTICAL {
					segConstraint.P0.Y += moveOverlapPointLength
				}
			}
		}
	}
}
