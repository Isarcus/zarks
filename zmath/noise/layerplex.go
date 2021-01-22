package noise

import (
	"fmt"
	"math"
	"math/rand"
	"zarks/zmath"
)

// NewLayerplexMap returns a newly generated simplex map, with a slightly modified algorithm.
// The data within 'layer' is used to modify the magnitude of the simplex vectors, for an interesting effect!
func NewLayerplexMap(cfg Config, layer zmath.Map) zmath.Map {
	cfg.checkDefaults()
	r2 := cfg.R * cfg.R
	layerplexMap := zmath.NewMap(cfg.Dimensions, 0)

	conv := zmath.NewVec(
		float64(layer.Bounds().X)/float64(cfg.Dimensions.X),
		float64(layer.Bounds().Y)/float64(cfg.Dimensions.Y),
	)

	for oct := 0.0; oct < float64(cfg.Octaves); oct++ {
		// create a new hashmap for angles
		vecMap := make(map[string]float64)
		octInfluence := math.Pow(0.5, oct)

		// dive in
		for x := 0; x < cfg.Dimensions.X; x++ {
			for y := 0; y < cfg.Dimensions.Y; y++ {
				iptX := float64(x) / (cfg.BoxSizeInitial * octInfluence)
				iptY := float64(y) / (cfg.BoxSizeInitial * octInfluence)
				ipt := zmath.NewVec(iptX, iptY)

				layerX := int(float64(x) * conv.X)
				layerY := int(float64(y) * conv.Y)

				// skew input coordinates
				skewed := skew(ipt)

				// find internal coordinates of simplex
				var corners [3]zmath.Vec
				var vectors [3]zmath.Vec
				corners[0] = zmath.NewVec(math.Floor(skewed.X), math.Floor(skewed.Y))
				internal := zmath.NewVec(skewed.X-corners[0].X, skewed.Y-corners[0].Y)
				if internal.X > internal.Y {
					corners[1] = corners[0].Add(zmath.NewVec(1, 0))
				} else {
					corners[1] = corners[0].Add(zmath.NewVec(0, 1))
				}
				corners[2] = corners[0].Add(zmath.NewVec(1, 1))

				// get vectors of the three corners
				for i := 0; i < 3; i++ {
					key := getKey(corners[i].ToInt())
					if angle, ok := vecMap[key]; ok { // if random angle already determined, use it
						vectors[i] = zmath.NewVec(math.Cos(angle), math.Sin(angle)).Scale(layer[layerX][layerY])
					} else { // otherwise, generate one for this and future use
						angle = rand.Float64() * 2.0 * math.Pi
						vectors[i] = zmath.NewVec(math.Cos(angle), math.Sin(angle)).Scale(layer[layerX][layerY])
						vecMap[key] = angle
					}
				}

				// determine distance to each unskewed corner
				var disp [3]zmath.Vec
				var distance [3]float64
				for i := range distance {
					disp[i] = ipt.Subtract(unskew(corners[i]))
					distance[i] = math.Sqrt(disp[i].X*disp[i].X + disp[i].Y*disp[i].Y)
				}

				var Z float64
				for i, dist := range distance {
					influence := math.Max(0, r2-dist*dist)
					Z += influence * disp[i].Dot(vectors[i])
				}

				layerplexMap[x][y] += Z * octInfluence
			}
		}
		fmt.Println("Octave ", oct+1, " finished.")
	}

	if cfg.Normalize {
		linearData := zmath.ToLinear(layerplexMap)
		linearData.Interpolate(0, 1)
		linearData.To2D(layerplexMap)
	}

	return layerplexMap
}
