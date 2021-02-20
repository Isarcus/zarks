package zimg

import (
	"image"
	"image/color"
	"math"
)

// ColorType is a way to refer to each of the four colors: R, G, B, and A!
type ColorType int

// Color Constants
const (
	Red ColorType = iota
	Green
	Blue
	Alpha
	R = Red
	G = Green
	B = Blue
	A = Alpha
)

// RGBA256 is a float64-resolution color
type RGBA256 struct {
	R, G, B, A float64
}

// RGBA returns a new color.RGBA with A set to 255 (opaque)
func RGBA(r, g, b uint8) color.RGBA {
	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}

func make32bit(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r << 8),
		G: uint8(g << 8),
		B: uint8(b << 8),
		A: uint8(a << 8),
	}
}

func toUint8(c color.Color) [4]uint8 {
	col32 := make32bit(c)
	return [4]uint8{col32.R, col32.G, col32.B, col32.A}
}

// BrightnessOf returns the unweighted brightness of a color. It does not take Alpha into account.
func BrightnessOf(c color.RGBA) float64 {
	return (float64(c.R) + float64(c.G) + float64(c.B)) / 765.0
}

// DistanceToColor returns the euclidean "distance" between two colors. It does NOT take alpha into account!
func DistanceToColor(c1, c2 color.RGBA) float64 {
	var (
		dR = float64(c1.R) - float64(c2.R)
		dG = float64(c1.G) - float64(c2.G)
		dB = float64(c1.B) - float64(c2.B)
	)

	return math.Sqrt(dR*dR + dG*dG + dB*dB)
}

// ImageToRGBA converts an image.Image to an *image.RGBA
func ImageToRGBA(img image.Image) *image.RGBA {
	var (
		width  int = img.Bounds().Dx()
		height int = img.Bounds().Dy()
	)

	imgRGBA := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			imgRGBA.Set(x, y, img.At(x, y))
		}
	}

	return imgRGBA
}

// RGBAFromImage is equivalent to ImageToRGBA
var RGBAFromImage = ImageToRGBA
