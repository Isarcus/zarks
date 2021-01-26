package zimg

import (
	"image"
	"zarks/zmath"
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

// Colors represents R, G, B, A for convenient looping!
var (
	Colors = [4]ColorType{
		Red, Green, Blue, Alpha,
	}
)

// ZImage is an RGBA image with 64-bit floating point accuracy for each color. It also has some really cool functions.
type ZImage struct {
	RGBA  [4]zmath.Map // zmath.Map representations of the image
	rep32 *image.RGBA  // a typical, 32-bit color representation of an image
}

// NewZImage returns an empty new ZImage of the desired size
func NewZImage(bounds zmath.VecInt) *ZImage {
	rep32 := image.NewRGBA(image.Rect(0, 0, bounds.X, bounds.Y))
	rgba := [4]zmath.Map{
		zmath.NewMap(bounds, 0),
		zmath.NewMap(bounds, 0),
		zmath.NewMap(bounds, 0),
		zmath.NewMap(bounds, 0),
	}
	return &ZImage{
		RGBA:  rgba,
		rep32: rep32,
	}
}

// ZImageFromRGBA creates a new ZImage from an image.RGBA
func ZImageFromRGBA(img *image.RGBA) *ZImage {
	bounds := zmath.VI(img.Bounds().Max.X, img.Bounds().Max.Y)

	zi := NewZImage(bounds)

	for x := 0; x < bounds.X; x++ {
		for y := 0; y < bounds.Y; y++ {

		}
	}

	return zi
}

// Bounds returns the bounds of the ZImage!
func (zi *ZImage) Bounds() zmath.VecInt {
	return zi.RGBA[Red].Bounds()
}

// Update refreshes the underlying 32-bit representation of the image
func (zi *ZImage) Update() *ZImage {
	return zi
}
