package zimg

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/Isarcus/zarks/zmath"
)

// Colorify returns an image colored according to the provided Map and ColorScheme
func Colorify(noiseMap zmath.Map, scheme ColorScheme) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, noiseMap.Bounds().X, noiseMap.Bounds().Y))

	for x := 0; x < noiseMap.Bounds().X; x++ {
		for y := 0; y < noiseMap.Bounds().Y; y++ {
			img.Set(x, y, scheme(noiseMap[x][y]))
		}
	}

	return img
}

// MapToImage will convert a zmath.Map into a grayscale image
func MapToImage(noiseMap zmath.Map) *image.RGBA {
	return Colorify(noiseMap, SchemeGrayscale)
}

// ImageFromMap is identical to MapToImage
var ImageFromMap = MapToImage

//                           //
// - - - COLOR SCHEMES - - - //
//                           //

// ColorScheme represents all preset coloring functions for a zmath.Map
type ColorScheme func(float64) color.RGBA

// SchemeGrayscale is for standard black and white variation
func SchemeGrayscale(value float64) color.RGBA {
	brightness := uint8(value * 255.0)
	return RGBA(brightness, brightness, brightness)
}

// SchemeChristmas is for red, white, and green!
func SchemeChristmas(value float64) color.RGBA {
	var r, g, b uint8
	brightness := value * 255.0
	r = uint8(zmath.MinMax(0, 255, -brightness*8+1536))               // 1536
	g = uint8(zmath.MinMax(0, 255, brightness*8-512))                 // 512
	b = uint8(zmath.MinMax(0, 255, -math.Abs(brightness*8-1024)+512)) // 512

	return RGBA(r, g, b)
}

// SchemeRainbowFlames was an accident
func SchemeRainbowFlames(value float64) color.RGBA {
	var r, g, b uint8
	brightness := value * 255.0
	r = uint8(zmath.MinMax(0, 255, brightness*8-512))
	g = uint8(zmath.MinMax(0, 255, -brightness*8+1536))
	b = uint8(zmath.MinMax(0, 255, -math.Abs(brightness*8-127)+512))

	return RGBA(r, g, b)
}

// SchemeClouds is for clouds
func SchemeClouds(value float64) color.RGBA {
	brightness := value * 255.0
	return RGBA(
		uint8(brightness),
		uint8(brightness),
		255,
	)
}

// SchemeRivers is best used on a riverplex map
func SchemeRivers(value float64) color.RGBA {
	var r, g, b uint8
	if value < 0.2 { // water
		r = 20
		g = 20
		b = 200
	} else if value < 0.28 { // sand
		r = 200
		g = 200
		b = 100
	} else if value < 0.4 { // grass
		r = 50
		g = 135
		b = 40
	} else if value < 0.6 { // trees
		r = 25
		g = 105
		b = 20
	} else if value < 0.8 { // stone
		r = 160
		g = 160
		b = 160
	} else { // snow
		r = 250
		g = 250
		b = 250
	}

	return RGBA(r, g, b)
}

// CustomDiscreteScheme is a convenient tool for creating discrete color functions from a set of colors, as well as
// thresholds corresponding to those colors.
func CustomDiscreteScheme(colors ColorSetDiscrete, thresholds ThresholdSetDiscrete) ColorScheme {
	var scheme ColorScheme = func(value float64) color.RGBA {
		for i := range colors {
			if i >= len(thresholds) || value < thresholds[i] {
				return colors[i]
			}
		}
		return color.RGBA{}
	}

	return scheme
}

// CustomSmoothScheme is like CustomDiscreteScheme, but for smooth color variation. The thresholds input
// should be 2 items SHORTER than the colors input.
func CustomSmoothScheme(colors ColorSetSmooth, thresholds ThresholdSetSmooth) ColorScheme {
	var scheme ColorScheme = func(value float64) color.RGBA {
		var thr0, thr1 float64
		for i := 0; i <= len(thresholds); i++ {
			if i == 0 {
				thr0 = 0
				if len(thresholds) == 0 {
					thr1 = 0
				} else {
					thr1 = thresholds[0]
				}
			} else if i == len(thresholds) {
				thr0 = thresholds[i-1]
				thr1 = 1
			} else {
				thr0 = thresholds[i-1]
				thr1 = thresholds[i]
			}

			if zmath.AreInRange(thr0, thr1, value) || value == thr1 {
				col0 := colors[i]
				col1 := colors[i+1]
				scale := (value - thr0) / (thr1 - thr0)
				return RGBA(
					uint8(float64(col0.R)*(1.0-scale)+float64(col1.R)*scale),
					uint8(float64(col0.G)*(1.0-scale)+float64(col1.G)*scale),
					uint8(float64(col0.B)*(1.0-scale)+float64(col1.B)*scale),
				)
			}
		}
		fmt.Println("Error in CustomSmoothScheme for value:", value)
		return color.RGBA{}
	}

	return scheme
}

type (
	// ColorSet is a set of colors
	ColorSet []color.RGBA
	// ColorSetDiscrete designates a ColorSet for use in the CustomDiscreteScheme function only
	ColorSetDiscrete ColorSet
	// ColorSetSmooth designates a ColorSet for use in the CustomSmoothScheme function only
	ColorSetSmooth ColorSet
)
type (
	// ThresholdSet is a set of thresholds
	ThresholdSet []float64
	// ThresholdSetDiscrete designates a ThresholdSet for use in the CustomDiscreteScheme function only
	ThresholdSetDiscrete ThresholdSet
	// ThresholdSetSmooth designates a ThresholdSet for use in the CustomSmoothScheme function only
	ThresholdSetSmooth ThresholdSet
)

// Presets to be passed into CustomDiscreteScheme()
var (
	// Terrain can be used with a number of thresholds for different effects
	ColorSetTerrain = ColorSetDiscrete{
		RGBA(15, 15, 140),   // dark water
		RGBA(20, 20, 200),   // light water
		RGBA(200, 200, 100), // sand
		RGBA(50, 135, 40),   // grass
		RGBA(25, 105, 20),   // trees
		RGBA(160, 160, 160), // stone
		RGBA(250, 250, 250), // snow
	}
	ThresholdSetRivers = ThresholdSetDiscrete{
		0.1,
		0.2,
		0.28,
		0.4,
		0.6,
		0.8,
	}
	ThresholdSetArchipelago = ThresholdSet{
		0.35,
		0.55,
		0.65,
		0.75,
		0.88,
		0.95,
	}

	// The first test of custom smooth variation!!!
	ColorSetTest = ColorSetSmooth{
		RGBA(255, 0, 0),
		RGBA(0, 255, 0),
		RGBA(0, 0, 255),
		RGBA(0, 0, 0),
	}
	ThresholdSetTest = ThresholdSetSmooth{
		0.2,
		0.8,
	}

	// Very nice terrain
	ColorSetTerrain2 = ColorSetSmooth{
		RGBA(10, 10, 150),   // water (deepest)
		RGBA(40, 40, 160),   // water (deep)
		RGBA(75, 110, 235),  // water (shallow)
		RGBA(145, 170, 255), // water (shallowest)
		RGBA(235, 235, 110), // sand
		RGBA(45, 150, 40),   // grass
		RGBA(25, 115, 25),   // woods
		RGBA(130, 153, 130), // sparse woods
		RGBA(135, 135, 135), // stone
		RGBA(250, 250, 250), // snow
	}
	ThresholdSetTerrain2 = ThresholdSetSmooth{
		0.08, // water (deepest)    -> water (deep)
		0.25, // water (deep)       -> water (shallow)
		0.40, // water (shallow)    -> water (shallowest)
		0.45, // water (shallowest) -> sand
		0.52, // sand  -> grass
		0.65, // grass -> woods
		0.75, // woods -> sparse woods
		0.80, // sparse woods -> stone
	}

	// The Rainbow!
	ColorSetRainbow = ColorSetSmooth{
		RGBA(255, 0, 0),
		RGBA(255, 255, 0),
		RGBA(0, 255, 0),
		RGBA(0, 255, 255),
		RGBA(0, 0, 255),
		RGBA(255, 0, 255),
	}
	ThresholdSetRainbow = ThresholdSetSmooth{
		0.3,
		0.45,
		0.55,
		0.7,
	}

	// Rainbows on black background
	ColorSetBlackRainbow = ColorSetSmooth{
		RGBA(0, 0, 0),
		RGBA(255, 0, 0),
		RGBA(0, 0, 0),
		RGBA(0, 255, 0),
		RGBA(0, 0, 0),
		RGBA(0, 0, 255),
		RGBA(0, 0, 0),
	}
	ThresholdSetBlackRainbow = ThresholdSetSmooth{
		0.2,
		0.3,
		0.5,
		0.6,
		0.8,
	}

	// Solid colors on black
	ColorSetRedMiddle = ColorSetSmooth{
		RGBA(0, 0, 0),
		RGBA(255, 0, 0),
		RGBA(0, 0, 0),
	}
	ColorSetGreenMiddle = ColorSetSmooth{
		RGBA(0, 0, 0),
		RGBA(0, 255, 0),
		RGBA(0, 0, 0),
	}
	ColorSetBlueMiddle = ColorSetSmooth{
		RGBA(0, 0, 0),
		RGBA(0, 0, 255),
		RGBA(0, 0, 0),
	}
	ThresholdSetMiddle = ThresholdSetSmooth{
		0.5,
	}

	// Red, White, and Green
	ColorSetChristmas = ColorSetSmooth{
		RGBA(200, 0, 0),
		RGBA(245, 245, 245),
		RGBA(20, 200, 20),
	}
	ThresholdSetChristmas = ThresholdSetSmooth{
		0.5,
	}
)
