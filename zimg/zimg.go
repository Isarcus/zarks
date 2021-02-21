package zimg

import (
	"image"
	"image/color"
	"math"

	"github.com/Isarcus/zarks/system"
	"github.com/Isarcus/zarks/zmath"
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
	RGBA256 [4]zmath.Map // zmath.Map representations of the image
	RGBA32  *image.RGBA  // a typical, 32-bit color representation of an image
	ptr     zptr
}

// zptr is used to point to a rectangular section of a ZImage. (TODO)
type zptr struct {
	active bool
	img    *ZImage
	area   *zmath.RectInt
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
		area:   zmath.RI(zmath.ZVI, bounds),
	}
	return &ZImage{
		RGBA256: rgba,
		RGBA32:  rep32,
		ptr:     src,
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

// ZImageFromPath loads a ZImage from the file at the specified path.
func ZImageFromPath(path string) *ZImage {
	return ZImageFromRGBA(ImageToRGBA(system.LoadImage(path))) // basically so you don't have to type all this
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
		R: uint8(zi.RGBA256[R].At(pos)),
		G: uint8(zi.RGBA256[G].At(pos)),
		B: uint8(zi.RGBA256[B].At(pos)),
		A: uint8(zi.RGBA256[A].At(pos)),
	}
}

// RedAt returns the Red value at the desired location
func (zi *ZImage) RedAt(pos zmath.VecInt) float64 { return zi.RGBA256[Red].At(pos) }

// GreenAt returns the Red value at the desired location
func (zi *ZImage) GreenAt(pos zmath.VecInt) float64 { return zi.RGBA256[Green].At(pos) }

// BlueAt returns the Red value at the desired location
func (zi *ZImage) BlueAt(pos zmath.VecInt) float64 { return zi.RGBA256[Blue].At(pos) }

// AlphaAt returns the Red value at the desired location
func (zi *ZImage) AlphaAt(pos zmath.VecInt) float64 { return zi.RGBA256[Alpha].At(pos) }

// Set sets the color at the desired location
func (zi *ZImage) Set(pos zmath.VecInt, col color.RGBA) *ZImage {
	zi.RGBA256[R].Set(pos, float64(col.R))
	zi.RGBA256[G].Set(pos, float64(col.G))
	zi.RGBA256[B].Set(pos, float64(col.B))
	zi.RGBA256[A].Set(pos, float64(col.A))
	return zi
}

// ScaleDim scales the ZImage proportionally by the provided factor
func (zi *ZImage) ScaleDim(by float64) *ZImage {
	for i, m := range zi.RGBA256 {
		zi.RGBA256[i] = m.ScaleDim(by)
	}
	zi.RGBA32 = image.NewRGBA(image.Rect(0, 0, zi.Bounds().X, zi.Bounds().Y))
	zi.Update()
	return zi
}

// SetMaxBounds scales the ZImage proportionally so that it will be within the new maximum bounds.
// If bounds larger than the current bounds are provided, nothing will be changed.
func (zi *ZImage) SetMaxBounds(newMax zmath.VecInt) *ZImage {
	var (
		bounds = zi.Bounds()
		xScale = float64(newMax.X) / float64(bounds.X)
		yScale = float64(newMax.Y) / float64(bounds.Y)
	)
	if xScale < 1 || yScale < 1 {
		zi.ScaleDim(math.Min(xScale, yScale))
	}

	return zi
}

// Clear sets every pixel in the ZImage to the desired color
func (zi *ZImage) Clear(col color.Color) *ZImage {
	c := toUint8(col)
	for i, m := range zi.RGBA256 {
		m.Clear(float64(c[i]))
	}
	return zi
}

// Zero wipes an entire ZImage to black
func (zi *ZImage) Zero() *ZImage {
	zi.Clear(color.RGBA{})
	return zi
}

// MakeOpaque sets the underlying Alpha map to 255
func (zi *ZImage) MakeOpaque() *ZImage {
	zi.RGBA256[Alpha].Clear(255)
	return zi
}

// BlurGaussian blurs gaussianly! (and very slowly, unfortunately)
func (zi *ZImage) BlurGaussian(radius int) *ZImage {
	for _, m := range zi.RGBA256 {
		m.BlurGaussian(radius)
	}
	return zi
}

// Bounds returns the bounds of the ZImage!
func (zi *ZImage) Bounds() zmath.VecInt {
	return zi.RGBA256[Red].Bounds()
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
			zi.RGBA32.Set(x, y, c)
		}
	}
	return zi
}

// Interpolate interpolates between a new min and max on the desired colors.
// If no colors are passed in, RGB (but not A) will be interpolated.
func (zi *ZImage) Interpolate(newMin, newMax float64, onColors ...ColorType) *ZImage {
	if len(onColors) == 0 {
		for _, m := range ColorsRGB {
			zi.RGBA256[m].Interpolate(newMin, newMax)
		}
		return zi
	}
	for _, m := range onColors {
		zi.RGBA256[m].Interpolate(newMin, newMax)
	}
	return zi
}

// FlipHorizontal flips the image across its Y-axis. By default, all colors are flipped, but if arguments are
// provided, then only the specified ColorTypes are flipped.
func (zi *ZImage) FlipHorizontal(onColors ...ColorType) *ZImage {
	if len(onColors) == 0 {
		for _, m := range ColorsAll {
			zi.RGBA256[m].FlipHorizontal()
		}
		return zi
	}
	for _, m := range onColors {
		zi.RGBA256[m].FlipHorizontal()
	}
	return zi
}

// FlipVertical flips the image across its X-axis. By default, all colors are flipped, but if arguments are
// provided, then only the specified ColorTypes are flipped.
func (zi *ZImage) FlipVertical(onColors ...ColorType) *ZImage {
	if len(onColors) == 0 {
		for _, m := range ColorsAll {
			zi.RGBA256[m].FlipVertical()
		}
		return zi
	}
	for _, m := range onColors {
		zi.RGBA256[m].FlipVertical()
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

// CustomModAt lets you pass in any function that takes in a color.RGBA and returns a color.RGBA, which it
// then applies to the specified pixels in the ZImage.
func (zi *ZImage) CustomModAt(modFunc func(color.RGBA) color.RGBA, at []zmath.VecInt) *ZImage {
	for _, pt := range at {
		zi.Set(pt, modFunc(zi.At(pt)))
	}
	return zi
}

// Save saves an image!
func (zi *ZImage) Save(path string) {
	zi.Update()
	system.SaveImage(path, zi.RGBA32)
}

// GetGrayscale returns an image.RGBA grayscale of the current ZImage, without modifying the underlying ZImage.
func (zi *ZImage) GetGrayscale() *image.RGBA {
	bounds := zi.Bounds()
	gs := image.NewRGBA(image.Rect(0, 0, bounds.X, bounds.Y))
	for x := 0; x < bounds.X; x++ {
		for y := 0; y < bounds.Y; y++ {
			bright := uint8(BrightnessOf(zi.At(zmath.VI(x, y))) * 255.0)
			gs.Set(x, y, color.RGBA{bright, bright, bright, 255})
		}
	}
	return gs
}
