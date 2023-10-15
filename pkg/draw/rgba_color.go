package draw

import "image/color"

func RGBAToColor(RGBA [4]uint8) color.Color {
	return color.RGBA{RGBA[0], RGBA[1], RGBA[2], RGBA[3]}
}
