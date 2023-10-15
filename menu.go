package polygon_editor

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/zeraye/polygon-editor/pkg/config"
	"github.com/zeraye/polygon-editor/pkg/constraint"
)

type Menu struct {
	config       *config.Config
	selector     *widget.Select
	sliderButton *widget.Button
	slider       *widget.Slider
}

func NewMenu(config *config.Config) *Menu {
	return &Menu{config: config}
}

func (m *Menu) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(float32(m.config.Window.Width-m.config.UI.RasterWidth), float32(m.config.UI.RasterHeight))

}

func (m *Menu) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// layout for settingsLabel
	topLeft := fyne.NewPos(0, 0)
	objects[0].Resize(size)
	objects[0].Move(topLeft)

	// layout for other objcets
	padding := theme.Padding()
	for _, child := range objects[1:] {
		childMin := child.MinSize()
		childMin.Width = size.Width - 6*padding // magic number, make UI look nice
		child.Resize(childMin)
		child.Move(fyne.NewPos(float32(size.Width-childMin.Width)/2, float32(size.Height-childMin.Height)/2))
	}
}

func SelectConstraintWrapper(g *Game) func(string) {
	return func(option string) {
		if g.selectedSegment != nil {
			// get constraint type
			cstr, ok := constraint.StringToConstraint[option]
			if !ok {
				log.Fatalln("Option", option, "is not valid constraint")
			}
			// delete old constraint
			for i, segConstraint := range g.constraints {
				if segConstraint.P0 == g.selectedSegment.P0 && segConstraint.P1 == g.selectedSegment.P1 {
					// remove constraint
					g.constraints = append(g.constraints[:i], g.constraints[i+1:]...)
				}
			}
			// add HORIZONTAL or VERTICAL constraint
			if cstr == constraint.HORIZONTAL || cstr == constraint.VERTICAL {
				g.constraints = append(g.constraints, constraint.NewSegmentConstraint(g.selectedSegment.P0, g.selectedSegment.P1, cstr))
			}
			// fix constraints
			constraint.FixSegmentConstraint(g.constraints, g.polygons, g.selectedSegment.P0, g.config.Miscellaneous.MoveOverlapPointLength, g.config.Miscellaneous.AllowMoveOverlapPoint)
			constraint.FixSegmentConstraint(g.constraints, g.polygons, g.selectedSegment.P1, g.config.Miscellaneous.MoveOverlapPointLength, g.config.Miscellaneous.AllowMoveOverlapPoint)
			g.Refresh()
		}
	}
}

// TODO: set offset on selected polygon
func SetOffsetWrapper(g *Game) func() {
	return func() {
		if g.selectedPolygon == nil {
			return
		}
		panic("Not implemented")
	}
}

// TODO: preview offset on selected polygon
func PreviewOffsetWrapper(g *Game) func(float64) {
	return func(offset float64) {
		if g.selectedPolygon == nil {
			g.menu.slider.SetValue(0)
			return
		}
		panic("Not implemented")
	}
}

func (m *Menu) BuildUI(g *Game) fyne.CanvasObject {
	settingsLabel := widget.NewLabelWithStyle("Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	options := make([]string, len(constraint.Constraints))
	for i := 0; i < len(constraint.Constraints); i++ {
		options[i] = constraint.ConstraintToString[constraint.Constraint(i)]
	}
	selector := widget.NewSelect(options, SelectConstraintWrapper(g))
	selector.SetSelectedIndex(0)
	selector.Disable()
	m.selector = selector

	sliderValue := 0.0
	sliderBind := binding.BindFloat(&sliderValue)
	slider := widget.NewSliderWithData(0, 30, sliderBind)
	slider.Step = 0.01
	slider.OnChanged = PreviewOffsetWrapper(g)
	sliderButton := widget.NewButton("Set offset", SetOffsetWrapper(g))
	sliderButton.Disable()
	m.slider = slider
	m.sliderButton = sliderButton
	sliderLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(sliderBind, "Offset (%0.2f)"))

	return container.New(m, settingsLabel, container.NewVBox(selector, sliderLabel, slider, sliderButton))
}
