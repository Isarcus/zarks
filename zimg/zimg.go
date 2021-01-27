package zimg

import (
	"image"
	"image/color"
	"zarks/zmath"
)

var (
	// ColorsAll represents R, G, B, A for convenient looping!
	ColorsAll = [4]ColorType{Red, Green, Blue, Alpha}
	// ColorsRGB represents R, G, B for convenient looping!
	ColorsRGB = [3]ColorType{Red, Green, Blue}
)

// ZImage is an RGBA image with 64-bit floating point accuracy for each color. It also has some really cool functions.
// Expect to need about 32 MB of RAM for a 1000x1000 ZImage.
type ZImage struct {
	RGBA  [4]zmath.Map // zmath.Map representations of the image
	rep32 *image.RGBA  // a typical, 32-bit color representation of an image
	src   zptr
}

// zptr is used to point to a rectangular section of a ZImage. (TODO)
type zptr struct {
	active   bool
	img      *ZImage
	min, max zmath.VecInt
}

// NewZImage creates returns a new, empty ZImage of the desired size
func NewZImage(bounds zmath.VecInt) *ZImage {
	rep32 := image.NewRGBA(image.Rect(0, 0, bounds.X, bounds.Y))
	rgba := [4]zmath.Map{
		zmath.NewMap(bounds, 0),
		zmath.NewMap(bounds, 0),
		zmath.NewMap(bounds, 0),
		zmath.NewMap(bounds, 0),
	}
	src := zptr{
		active: false,
		img:    nil,
	}
	return &ZImage{
		RGBA:  rgba,
		rep32: rep32,
		src:   src,
	}
}

// ZImageFromRGBA creates a new ZImage from an image.RGBA
func ZImageFromRGBA(img *image.RGBA) *ZImage {
	bounds := zmath.VI(img.Bounds().Max.X, img.Bounds().Max.Y)
	zi := NewZImage(bounds)
	for x := 0; x < bounds.X; x++ {
		for y := 0; y < bounds.Y; y++ {
			color := img.RGBAAt(x, y)
			zi.Set(zmath.VI(x, y), color)
		}
	}

	return zi
}

// At returns the color at the desired location
func (zi *ZImage) At(pos zmath.VecInt) color.RGBA {
	return color.RGBA{
		R: uint8(zi.RGBA[R].At(pos)),
		G: uint8(zi.RGBA[G].At(pos)),
		B: uint8(zi.RGBA[B].At(pos)),
		A: uint8(zi.RGBA[A].At(pos)),
	}
}

// Set sets the color at the desired location
func (zi *ZImage) Set(pos zmath.VecInt, col color.Color) *ZImage {
	c := toUint8(col)
	for i, m := range zi.RGBA {
		m.Set(pos, float64(c[i]))
	}
	return zi
}

// Clear sets every pixel in the ZImage to the desired color
func (zi *ZImage) Clear(col color.Color) *ZImage {
	c := toUint8(col)
	for i, m := range zi.RGBA {
		m.Clear(float64(c[i]))
	}
	return zi
}

// Zero wipes an entire ZImage to black
func (zi *ZImage) Zero() *ZImage {
	zi.Clear(color.RGBA{})
	return zi
}

// BlurGaussian blurs gaussianly! (and very slowly, unfortunately)
func (zi *ZImage) BlurGaussian(radius int) *ZImage {
	for _, m := range zi.RGBA {
		m.BlurGaussian(radius)
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
