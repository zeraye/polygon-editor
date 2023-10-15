package polygon_editor

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/zeraye/polygon-editor/pkg/config"
	"github.com/zeraye/polygon-editor/pkg/constraint"
	"github.com/zeraye/polygon-editor/pkg/geom"
	"github.com/zeraye/polygon-editor/pkg/sample"
)

type Game struct {
	widget.BaseWidget

	config          *config.Config
	menu            *Menu
	polygons        []*geom.Polygon
	constraints     []*constraint.SegmentConstraint
	draggedPoint    *geom.Point
	draggedSegment  *geom.Segment
	draggedPolygon  *geom.Polygon
	selectedSegment *geom.Segment
	selectedPolygon *geom.Polygon
}

func NewGame(config *config.Config) *Game {
	menu := NewMenu(config)
	polygons, constraints := sample.GenerateSamplePolygons()
	game := &Game{config: config, menu: menu, polygons: polygons, constraints: constraints}
	game.ExtendBaseWidget(game)

	return game
}

func (g *Game) BuildUI() fyne.CanvasObject {
	return container.NewBorder(nil, nil, g.menu.BuildUI(g), g)
}

func (g *Game) CreateRenderer() fyne.WidgetRenderer {
	renderer := &gameRenderer{game: g}
	raster := canvas.NewRaster(renderer.Draw)
	renderer.raster = raster
	renderer.objects = []fyne.CanvasObject{raster}

	return renderer
}

func (g *Game) Tapped(ev *fyne.PointEvent) {
	mouse_pos := geom.NewPoint(float64(ev.Position.X), float64(ev.Position.Y))

	// check if point is tapped
	for _, poly := range g.polygons {
		for _, point := range poly.Points {
			if geom.DistanceToPoint(*mouse_pos, *point) < g.config.UI.PointRadius {
				// if polygon can be closed (over 3 points)
				if len(poly.Points) > 2 {
					// close polygon
					poly.IsClosed = true
					g.Refresh()
					return
				}
			}
		}
	}

	// else if segment is tapped
	for _, poly := range g.polygons {
		for _, seg := range poly.Segments() {
			if geom.DistanceToSegment(*mouse_pos, seg) < g.config.Miscellaneous.LineCatchError {
				// remove old selected segment and polygon
				g.selectedSegment = nil
				g.selectedPolygon = nil
				// enable constraint menu
				g.menu.selector.Enable()
				g.menu.sliderButton.Enable()
				// remove neighbours constraints from options
				// two same constraints can't be next to each other
				allowHorizonstal, allowVertical := true, true
				for _, segConstraint := range g.constraints {
					if segConstraint.P1 == seg.P0 || segConstraint.P0 == seg.P1 {
						if segConstraint.Cstr == constraint.HORIZONTAL {
							allowHorizonstal = false
						} else if segConstraint.Cstr == constraint.VERTICAL {
							allowVertical = false
						}
					}
				}
				options := []string{constraint.ConstraintToString[constraint.NO_CONSTRAINT]}
				if allowHorizonstal {
					options = append(options, constraint.ConstraintToString[constraint.HORIZONTAL])
				}
				if allowVertical {
					options = append(options, constraint.ConstraintToString[constraint.VERTICAL])
				}
				g.menu.selector.SetOptions(options)
				// select current constraint
				g.menu.selector.SetSelectedIndex(0)
				for _, segConstraint := range g.constraints {
					if segConstraint.P0 == seg.P0 && segConstraint.P1 == seg.P1 {
						g.menu.selector.SetSelected(constraint.ConstraintToString[segConstraint.Cstr])
					}
				}
				// set selected segment and polygon
				g.selectedSegment = &seg
				g.selectedPolygon = poly
				g.Refresh()
				return
			}
		}
	}

	// else create new point
	// check if any poly is not closed
	for _, poly := range g.polygons {
		if !poly.IsClosed {
			// add new point to poly
			poly.Points = append(poly.Points, geom.NewPoint(mouse_pos.X, mouse_pos.Y))
			g.Refresh()
			return
		}
	}

	// create new poly with new point
	g.polygons = append(g.polygons, geom.NewPolygon(geom.NewPoint(mouse_pos.X, mouse_pos.Y)))
	g.Refresh()
}

func (g *Game) TappedSecondary(ev *fyne.PointEvent) {
	mouse_pos := geom.NewPoint(float64(ev.Position.X), float64(ev.Position.Y))

	// check if point is tapped
	for poly_idx, poly := range g.polygons {
		for _, point := range poly.Points {
			if geom.DistanceToPoint(*mouse_pos, *point) < g.config.UI.PointRadius {
				// polygon with 2 points cannot be proper closed polygon (polygonal chain),
				// thus polygon should be removed
				if len(poly.Points) == 3 {
					// remove polygon
					g.polygons = append(g.polygons[:poly_idx], g.polygons[poly_idx+1:]...)
				} else {
					// remove point
					err := poly.RemovePoint(point)
					if err != nil {
						log.Fatal(err)
					}
				}
				g.Refresh()
				return
			}
		}
	}

	// else if segment is tapped
	for _, poly := range g.polygons {
		for _, seg := range poly.Segments() {
			if geom.DistanceToSegment(*mouse_pos, seg) < g.config.Miscellaneous.LineCatchError {
				// remove constraints related to segment
				for i, segConstraint := range g.constraints {
					if segConstraint.P0 == seg.P0 && segConstraint.P1 == seg.P1 {
						g.constraints = append(g.constraints[:i], g.constraints[i+1:]...)
					}
				}
				// add point in the middle of segment
				p0, p1 := *seg.P0, *seg.P1
				new_p := geom.NewPoint((p0.X+p1.X)/2, (p0.Y+p1.Y)/2)
				// seg.P0 is always before seg.P1
				poly.AddPointAfter(seg.P0, new_p)
				g.Refresh()
				return
			}
		}
	}
}

func (g *Game) Dragged(ev *fyne.DragEvent) {
	mouse_pos := geom.NewPoint(float64(ev.Position.X), float64(ev.Position.Y))
	dx, dy := float64(ev.Dragged.DX), float64(ev.Dragged.DY)

	// if point is currently dragged
	if g.draggedPoint != nil {
		// move dragged point
		g.draggedPoint.X += dx
		g.draggedPoint.Y += dy
		// fix constraints
		constraint.FixSegmentConstraint(g.constraints, g.polygons, g.draggedPoint, g.config.Miscellaneous.MoveOverlapPointLength, g.config.Miscellaneous.AllowMoveOverlapPoint)
		g.Refresh()
		return
	}

	// else if segment is currently dragged
	if g.draggedSegment != nil {
		// move dragged segment
		g.draggedSegment.P0.X += dx
		g.draggedSegment.P0.Y += dy
		g.draggedSegment.P1.X += dx
		g.draggedSegment.P1.Y += dy
		// fix constraints
		constraint.FixSegmentConstraint(g.constraints, g.polygons, g.draggedSegment.P0, g.config.Miscellaneous.MoveOverlapPointLength, g.config.Miscellaneous.AllowMoveOverlapPoint)
		constraint.FixSegmentConstraint(g.constraints, g.polygons, g.draggedSegment.P1, g.config.Miscellaneous.MoveOverlapPointLength, g.config.Miscellaneous.AllowMoveOverlapPoint)
		g.Refresh()
		return
	}

	// else if polygon is currently dragged
	if g.draggedPolygon != nil {
		// move dragged polygon
		for _, point := range g.draggedPolygon.Points {
			point.X += dx
			point.Y += dy
		}
		g.Refresh()
		return
	}

	// else if point is tapped (dragged)
	for _, poly := range g.polygons {
		for _, point := range poly.Points {
			if geom.DistanceToPoint(*mouse_pos, *point) < g.config.UI.PointRadius {
				// set dragged point and move it
				g.draggedPoint = point
				g.draggedPoint.X += dx
				g.draggedPoint.Y += dy
				// fix constraints
				constraint.FixSegmentConstraint(g.constraints, g.polygons, g.draggedPoint, g.config.Miscellaneous.MoveOverlapPointLength, g.config.Miscellaneous.AllowMoveOverlapPoint)
				g.Refresh()
				return
			}
		}
	}

	// else if segment is tapped (dragged)
	for _, poly := range g.polygons {
		for _, seg := range poly.Segments() {
			if geom.DistanceToSegment(*mouse_pos, seg) < g.config.Miscellaneous.LineCatchError {
				// set dragged segment and move it
				g.draggedSegment = &seg
				g.draggedSegment.P0.X += dx
				g.draggedSegment.P0.Y += dy
				g.draggedSegment.P1.X += dx
				g.draggedSegment.P1.Y += dy
				// fix constraints
				constraint.FixSegmentConstraint(g.constraints, g.polygons, g.draggedSegment.P0, g.config.Miscellaneous.MoveOverlapPointLength, g.config.Miscellaneous.AllowMoveOverlapPoint)
				constraint.FixSegmentConstraint(g.constraints, g.polygons, g.draggedSegment.P1, g.config.Miscellaneous.MoveOverlapPointLength, g.config.Miscellaneous.AllowMoveOverlapPoint)
				g.Refresh()
				return
			}
		}
	}

	// else if polygon is tapped (dragged)
	for _, poly := range g.polygons {
		if geom.IsPointInsidePolygon(*poly, *mouse_pos) {
			// set dragged polygon and move it
			g.draggedPolygon = poly
			for _, point := range g.draggedPolygon.Points {
				point.X += dx
				point.Y += dy
			}
			g.Refresh()
			return
		}
	}
}

func (g *Game) DragEnd() {
	if g.draggedPoint != nil {
		g.draggedPoint = nil
	}
	if g.draggedSegment != nil {
		g.draggedSegment = nil
	}
	if g.draggedPolygon != nil {
		g.draggedPolygon = nil
	}
}