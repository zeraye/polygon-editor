package polygon_editor

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/zeraye/polygon-editor/pkg/draw"
)

type gameRenderer struct {
	raster  *canvas.Raster
	objects []fyne.CanvasObject
	game    *Game
}

func (gr *gameRenderer) Destroy() {
}

func (gr *gameRenderer) Layout(size fyne.Size) {
	gr.raster.Resize(size)
}

func (gr *gameRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(gr.game.config.UI.RasterWidth), float32(gr.game.config.UI.RasterHeight))
}

func (gr *gameRenderer) Objects() []fyne.CanvasObject {
	return gr.objects
}

func (gr *gameRenderer) Refresh() {
	canvas.Refresh(gr.raster)
}

// Draw game raster (polygon editor area, not menus)
func (gr *gameRenderer) Draw(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(gr.game.config.UI.RasterWidth), int(gr.game.config.UI.RasterHeight)))

	// draw raster background
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, draw.RGBAToColor(gr.game.config.UI.BackgroundColorRGBA))
		}
	}

	// draw raster border
	for x := 0; x < img.Bounds().Dx(); x++ {
		img.Set(x, 0, draw.RGBAToColor(gr.game.config.UI.RasterBorderColorRGBA))
		img.Set(x, img.Bounds().Dy()-1, draw.RGBAToColor(gr.game.config.UI.RasterBorderColorRGBA))
	}
	for y := 0; y < img.Bounds().Dx(); y++ {
		img.Set(0, y, draw.RGBAToColor(gr.game.config.UI.RasterBorderColorRGBA))
		img.Set(img.Bounds().Dx()-1, y, draw.RGBAToColor(gr.game.config.UI.RasterBorderColorRGBA))
	}

	// draw offset polygon
	if len(gr.game.offsetPolygons) > 0 {
		for _, poly := range gr.game.offsetPolygons {
			draw.DrawPolygon(
				*poly,
				nil,
				draw.RGBAToColor(gr.game.config.UI.PointColorRGBA),
				draw.RGBAToColor(gr.game.config.UI.OffsetSegmentColorRGBA),
				draw.RGBAToColor(gr.game.config.UI.SelectedSegmentColorRGBA),
				gr.game.config.UI.PointRadius,
				false,
				gr.game.config.UI.FillPoints,
				img,
				gr.game.config.Miscellaneous.LineDrawAlgorithm,
				nil,
			)
		}
	}

	// draw polygons
	for _, poly := range gr.game.polygons {
		draw.DrawPolygon(
			*poly,
			gr.game.selectedSegment,
			draw.RGBAToColor(gr.game.config.UI.PointColorRGBA),
			draw.RGBAToColor(gr.game.config.UI.SegmentColorRGBA),
			draw.RGBAToColor(gr.game.config.UI.SelectedSegmentColorRGBA),
			gr.game.config.UI.PointRadius,
			true,
			gr.game.config.UI.FillPoints,
			img,
			// gr.game.config.Miscellaneous.LineDrawAlgorithm,
			gr.game.lineAlgorithm,
			gr.game.widths,
		)
	}

	// draw constraints
	for _, constraint := range gr.game.constraints {
		draw.DrawConstraint(*constraint, draw.RGBAToColor(gr.game.config.UI.ConstraintCharColorRGBA), img)
	}

	return img
}
