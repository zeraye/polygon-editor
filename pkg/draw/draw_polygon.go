package draw

import (
	"image"
	"image/color"

	"github.com/zeraye/polygon-editor/pkg/geom"
)

func DrawPolygon(poly geom.Polygon, s *geom.Segment, pointColor, segmentColor, selectedSegmentColor color.Color, circleRadius float64, drawPoints, fillPoints bool, img *image.RGBA) {
	// draw segments (without last)
	for i := 1; i < len(poly.Points); i++ {
		// if segment is selected, color it properly
		if s != nil && s.P0 == poly.Points[i-1] && s.P1 == poly.Points[i] {
			BresenhamDrawLine(*poly.Points[i-1], *poly.Points[i], selectedSegmentColor, img)
		} else {
			BresenhamDrawLine(*poly.Points[i-1], *poly.Points[i], segmentColor, img)
		}
	}

	// draw closing segment
	if poly.IsClosed {
		// if segment is selected, color it properly
		if s != nil && s.P0 == poly.Points[len(poly.Points)-1] && s.P1 == poly.Points[0] {
			BresenhamDrawLine(*poly.Points[len(poly.Points)-1], *poly.Points[0], selectedSegmentColor, img)
		} else {
			BresenhamDrawLine(*poly.Points[len(poly.Points)-1], *poly.Points[0], segmentColor, img)
		}
	}

	// draw points
	if drawPoints {
		for _, p := range poly.Points {
			DrawCircle(*p, circleRadius, pointColor, fillPoints, img)
		}
	}
}
