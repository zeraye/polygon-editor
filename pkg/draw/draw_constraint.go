package draw

import (
	"image"
	"image/color"

	"github.com/zeraye/polygon-editor/pkg/constraint"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Draw constraint as char in the middle of segment
func DrawConstraint(segConstraint constraint.SegmentConstraint, color color.Color, img *image.RGBA) {
	// don't draw NO_CONSTRAINT
	cstr := segConstraint.Cstr
	if cstr == constraint.NO_CONSTRAINT {
		return
	}

	// get mid-point
	p0, p1 := segConstraint.P0, segConstraint.P1
	x, y := (p0.X+p1.X)/2, (p0.Y+p1.Y)/2
	point := fixed.Point26_6{X: fixed.I(int(x)), Y: fixed.I(int(y))}
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	// draw only one (first) char
	// eg. for HORIZONTAL -> H
	drawer.DrawString(constraint.ConstraintToString[cstr][:1])
}
