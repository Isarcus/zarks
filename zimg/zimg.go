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
	ptr   zptr
}

// zptr is used to point to a rectangular section of a ZImage. (TODO)
type zptr struct {
	active bool
	img    *ZImage
	area   zmath.Rect
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
		area:   zmath.R(zmath.ZVI, bounds),
	}
	return &ZImage{
		RGBA:  rgba,
		rep32: rep32,
		ptr:   src,
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

// Src returns the underlying source ZImage of a ZImage (or itself, if it is not a reference.)
func (zi *ZImage) Src() *ZImage {
	if zi.Original() {
		return zi
	}
	return zi.ptr.img
}

// Original returns whether the current ZImage is an original ZImage or just a reference to another one.
func (zi *ZImage) Original() bool {
	return zi.ptr.active
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

// RedAt returns the Red value at the desired location
func (zi *ZImage) RedAt(pos zmath.VecInt) float64 { return zi.RGBA[Red].At(pos) }

// GreenAt returns the Red value at the desired location
func (zi *ZImage) GreenAt(pos zmath.VecInt) float64 { return zi.RGBA[Green].At(pos) }

// BlueAt returns the Red value at the desired location
func (zi *ZImage) BlueAt(pos zmath.VecInt) float64 { return zi.RGBA[Blue].At(pos) }

// AlphaAt returns the Red value at the desired location
func (zi *ZImage) AlphaAt(pos zmath.VecInt) float64 { return zi.RGBA[Alpha].At(pos) }

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
	b := zi.Bounds()
	for x := 0; x < b.X; x++ {
		for y := 0; y < b.Y; y++ {
			c := color.RGBA{
				R: uint8(zi.RedAt(zmath.VI(x, y))),
				G: uint8(zi.GreenAt(zmath.VI(x, y))),
				B: uint8(zi.BlueAt(zmath.VI(x, y))),
				A: uint8(zi.AlphaAt(zmath.VI(x, y))),
			}
			zi.rep32.Set(x, y, c)
		}
	}
	return zi
}

// Interpolate interpolates between a new min and max on the desired colors.
// If no colors are passed in, RGB (but not A) will be interpolated.
func (zi *ZImage) Interpolate(newMin, newMax float64, onColors ...ColorType) *ZImage {
	if len(onColors) == 0 {
		for _, m := range ColorsRGB {
			zi.RGBA[m].Interpolate(newMin, newMax)
		}
		return zi
	}
	for _, m := range onColors {
		zi.RGBA[m].Interpolate(newMin, newMax)
	}
	return zi
}

// FlipHorizontal flips the image across its Y-axis. By default, all colors are flipped, but if arguments are
// provided, then only the specified ColorTypes are flipped.
func (zi *ZImage) FlipHorizontal(onColors ...ColorType) *ZImage {
	if len(onColors) == 0 {
		for _, m := range ColorsAll {
			zi.RGBA[m].FlipHorizontal()
		}
		return zi
	}
	for _, m := range onColors {
		zi.RGBA[m].FlipHorizontal()
	}
	return zi
}

// FlipVertical flips the image across its X-axis. By default, all colors are flipped, but if arguments are
// provided, then only the specified ColorTypes are flipped.
func (zi *ZImage) FlipVertical(onColors ...ColorType) *ZImage {
	if len(onColors) == 0 {
		for _, m := range ColorsAll {
			zi.RGBA[m].FlipVertical()
		}
		return zi
	}
	for _, m := range onColors {
		zi.RGBA[m].FlipVertical()
	}
	return zi
}

// Replace replaces all instances of a particular color with another color.
func (zi *ZImage) Replace(col, with color.Color) *ZImage {
	orig := make32bit(col)
	repl := make32bit(with)

	b := zi.Bounds()
	for x := 0; x < b.X; x++ {
		for y := 0; y < b.Y; y++ {
			if zi.At(zmath.VI(x, y)) == orig {
				zi.Set(zmath.VI(x, y), repl)
			}
		}
	}
	return zi
}

// ReplaceNot replaces all pixels that are NOT a particular color with another color.
func (zi *ZImage) ReplaceNot(col, with color.Color) *ZImage {
	orig := make32bit(col)
	repl := make32bit(with)

	b := zi.Bounds()
	for x := 0; x < b.X; x++ {
		for y := 0; y < b.Y; y++ {
			if zi.At(zmath.VI(x, y)) != orig {
				zi.Set(zmath.VI(x, y), repl)
			}
		}
	}
	return zi
}

// CustomMod lets you pass in any function that takes in a color.RGBA and returns a color.RGBA, which it
// then applies to every pixel in the ZImage.
func (zi *ZImage) CustomMod(modFunc func(color.RGBA) color.RGBA) *ZImage {
	b := zi.Bounds()
	for x := 0; x < b.X; x++ {
		for y := 0; y < b.Y; y++ {
			zi.Set(zmath.VI(x, y), modFunc(zi.At(zmath.VI(x, y))))
		}
	}
	return zi
}
