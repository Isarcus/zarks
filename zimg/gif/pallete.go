package gif

import "image/color"

// Pallete represents a set of colors for a GIF to choose from.
type Pallete [256]color.RGBA

// GetPalleteGrayscale returns a basic grayscale pallete
func GetPalleteGrayscale() Pallete {
	pal := [256]color.RGBA{}

	for i := range pal {
		b := uint8(i)
		pal[i] = color.RGBA{b, b, b, 255}
	}

	return pal
}

// Contains tells you whether the pallete contains the queried color
func (p Pallete) Contains(c color.RGBA) bool {
	for _, pCol := range p {
		if pCol == c {
			return true
		}
	}
	return false
}
