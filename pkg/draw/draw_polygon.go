package draw

import (
	"image"
	"image/color"

	"github.com/zeraye/polygon-editor/pkg/geom"
)

func GetWidth(widths map[*geom.Point]float64, p *geom.Point) float64 {
	if widths == nil {
		return 1
	}
	val, ok := widths[p]
	if ok {
		return val
	}
	return 1
}

func DrawPolygon(poly geom.Polygon, s *geom.Segment, pointColor, segmentColor, selectedSegmentColor color.Color, circleRadius float64, drawPoints, fillPoints bool, img *image.RGBA, drawLineAlgorithm string, widths map[*geom.Point]float64) {
	// draw segments (without last)
	for i := 1; i < len(poly.Points); i++ {
		// if segment is selected, color it properly
		if s != nil && s.P0 == poly.Points[i-1] && s.P1 == poly.Points[i] {
			DrawLine(*poly.Points[i-1], *poly.Points[i], selectedSegmentColor, img, drawLineAlgorithm, GetWidth(widths, poly.Points[i-1]))
		} else {
			DrawLine(*poly.Points[i-1], *poly.Points[i], segmentColor, img, drawLineAlgorithm, GetWidth(widths, poly.Points[i-1]))
		}
	}

	// draw closing segment
	if poly.IsClosed {
		// if segment is selected, color it properly
		if s != nil && s.P0 == poly.Points[len(poly.Points)-1] && s.P1 == poly.Points[0] {
			DrawLine(*poly.Points[len(poly.Points)-1], *poly.Points[0], selectedSegmentColor, img, drawLineAlgorithm, GetWidth(widths, poly.Points[len(poly.Points)-1]))
		} else {
			DrawLine(*poly.Points[len(poly.Points)-1], *poly.Points[0], segmentColor, img, drawLineAlgorithm, GetWidth(widths, poly.Points[len(poly.Points)-1]))
		}
	}

	// draw points
	if drawPoints {
		for _, p := range poly.Points {
			DrawCircle(*p, circleRadius, pointColor, fillPoints, img)
		}
	}
}
