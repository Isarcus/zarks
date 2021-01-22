package zimg

import (
	"image"
	"image/color"
	"math"
	"zarks/zmath"
)

// BlurGaussian blurs an image!
func BlurGaussian(inputImg *image.RGBA, radius int) *image.RGBA {
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
			blurred := FuncBlurGaussian(
				inputImg,
				bounds,
				zmath.NewVecInt(x, y),
				circle,
				float64(radius)/3.0,
			)
			outputImg.Set(x, y, blurred)

		}
	}

	return outputImg
}

// FuncBlurGaussian is for gaussian blurs
func FuncBlurGaussian(img *image.RGBA, bounds zmath.BoxInt, pos zmath.VecInt, circle []zmath.VecInt, sigma float64) color.RGBA {
	var r, g, b float64
	var weightSum, trueWeightSum float64
	var gaussCoeff float64 = 1.0 / (2.0 * math.Pi * sigma * sigma)

	for _, cpt := range circle {
		testPos := pos.Add(cpt)
		weight := math.Pow(math.E, -1.0*float64(cpt.X*cpt.X+cpt.Y*cpt.Y)/(2*sigma*sigma))
		trueWeightSum += weight
		if zmath.IsWithinBounds(testPos, bounds) {
			// weight summation
			weightSum += weight

			// color summing
			color := img.RGBAAt(testPos.X, testPos.Y)
			r += weight * float64(color.R)
			g += weight * float64(color.G)
			b += weight * float64(color.B)
		}
	}

	ratio := gaussCoeff * trueWeightSum / weightSum

	r *= ratio
	g *= ratio
	b *= ratio

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}
