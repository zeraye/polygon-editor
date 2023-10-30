package offset

import (
	"github.com/zeraye/polygon-editor/pkg/geom"
)

func CreateOffset(poly *geom.Polygon, offset float64, offsetAlgorithm string) []*geom.Polygon {
	if offsetAlgorithm == "custom1" {
		return createOffsetCustom1(poly, offset)
	} else if offsetAlgorithm == "custom2" {
		return createOffsetCustom2(poly, offset, "miter")
	}
	panic("Invalid offset algorithm")
}
