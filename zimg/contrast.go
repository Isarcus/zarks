package zimg

import (
	"image"
	"image/color"
	"math"
	"zarks/zmath"
)

// GetContrast outputs a black and white contrast image.
func GetContrast(inputImg *image.RGBA, radius int) *image.RGBA {
	var (
		width  int = inputImg.Bounds().Dx()
		height int = inputImg.Bounds().Dy()
	)
	outputImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// For the loop
	circle := zmath.GetCircleCoords(radius)
	bounds := zmath.NewBoxInt(0, 0, width, height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			contrast := FuncContrastGaussian(
				inputImg,
				bounds,
				zmath.NewVecInt(x, y),
				circle,
				float64(radius)/3.0,
			)

			RGBVal := uint8(int((1.0 - contrast) * 255.0)) // probably don't need double-cast
			RGBA := color.RGBA{
				R: RGBVal,
				G: RGBVal,
				B: RGBVal,
				A: 255,
			}
			outputImg.Set(x, y, RGBA)

		}
	}

	return outputImg
}

// TODO: make a gauss.go file that has a function that calculates the gaussian weights of a given circle
// so that it need only be done ONCE. (The method below is unnecessarily slow)

// FuncContrastGaussian returns the gaussian blur for the given point
func FuncContrastGaussian(img *image.RGBA, bounds zmath.BoxInt, pos zmath.VecInt, circle []zmath.VecInt, sigma float64) float64 {
	var pxColor = img.RGBAAt(pos.X, pos.Y)

	var pixelSum float64 = 0
	var weightSum float64 = 0
	var gaussCoeff float64 = 1.0 / (2.0 * math.Pi * sigma * sigma)

	for _, cpt := range circle {
		testPos := pos.Add(cpt)
		if zmath.IsWithinBounds(testPos, bounds) && cpt != zmath.ZVI {
			// increment pixel counter
			pixelSum++

			// color difference calculations
			testColor := img.RGBAAt(testPos.X, testPos.Y)
			var dCol float64 = 0
			dCol += math.Abs(float64(int(pxColor.R)-int(testColor.R))) / 255.0
			dCol += math.Abs(float64(int(pxColor.G)-int(testColor.G))) / 255.0
			dCol += math.Abs(float64(int(pxColor.B)-int(testColor.B))) / 255.0
			dCol /= 3.0

			// weight calculations
			weight := math.Pow(math.E, -1.0*float64(cpt.X*cpt.X+cpt.Y*cpt.Y)/(2*sigma*sigma))
			weightSum += weight * dCol
		}
	}

	return weightSum * gaussCoeff * (float64(len(circle)-1) / pixelSum)
}
