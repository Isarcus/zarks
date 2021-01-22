package noise

import (
	"math"
	"zarks/zmath"
)

// Riverplex is useful for lots of rivers
func Riverplex(cfg Config) zmath.Map {
	layer1 := NewSimplexMap(cfg).Interpolate(-1, 1)
	theMap := NewSimplexMap(cfg).Interpolate(1, 2).GeometricMean(layer1).Interpolate(0, 1)

	return theMap
}

// Ridgeplex is useful for mountain ranges
func Ridgeplex(cfg Config) zmath.Map {
	theMap := zmath.NewMap(cfg.Dimensions, 0)
	for oct := 0.0; int(oct) < cfg.Octaves; oct++ {
		octInfluence := math.Pow(0.5, oct)
		newCfg := Config{
			Dimensions:     cfg.Dimensions,
			Octaves:        1,
			Normalize:      false,
			Seed:           cfg.Seed,
			BoxSizeInitial: cfg.BoxSizeInitial * octInfluence,
			R:              cfg.R,
		}
		layer := NewSimplexMap(newCfg).Interpolate(-0.5, 1).CustomMod(math.Abs).Multiply(-1).Add(1).Multiply(octInfluence)
		theMap.AddMap(layer) /*                     ^ this right here should be -0.5 */
	}
	theMap.Interpolate(0, 1)

	return theMap
}
