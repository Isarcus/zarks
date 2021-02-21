package noise

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/Isarcus/zarks/zmath"
)

// dimensional transformation constants
var (
	F2D float64 = (math.Sqrt(3) - 1.0) / 2.0
	G2D float64 = (1.0 - 1.0/math.Sqrt(3)) / 2.0
)

// NewSimplexMap generates a new simplex noise map according to the specified configuration.
func NewSimplexMap(cfg Config) zmath.Map {
	cfg.checkDefaults()
	r2 := cfg.R * cfg.R
	simplexMap := zmath.NewMap(cfg.Dimensions, 0)

	for oct := 0.0; oct < float64(cfg.Octaves); oct++ {
		// create a new hashmap for angles
		vecMap := make(map[string]float64)
		octInfluence := math.Pow(0.5, oct)

		// dive in
		for x := 0; x < cfg.Dimensions.X; x++ {
			for y := 0; y < cfg.Dimensions.Y; y++ {
				ipt := zmath.V(
					float64(x)/(cfg.BoxSizeInitial*octInfluence),
					float64(y)/(cfg.BoxSizeInitial*octInfluence),
				)

				// skew input coordinates
				skewed := skew(ipt)

				// find internal coordinates of simplex
				var corners [3]zmath.Vec
				var vectors [3]zmath.Vec
				corners[0] = zmath.V(math.Floor(skewed.X), math.Floor(skewed.Y))
				internal := zmath.V(skewed.X-corners[0].X, skewed.Y-corners[0].Y)
				if internal.X > internal.Y {
					corners[1] = corners[0].Add(zmath.V(1, 0))
				} else {
					corners[1] = corners[0].Add(zmath.V(0, 1))
				}
				corners[2] = corners[0].Add(zmath.V(1, 1))

				// get vectors of the three corners
				for i := 0; i < 3; i++ {
					key := getKey(corners[i].VI())
					if angle, ok := vecMap[key]; ok { // if random angle already determined, use it
						vectors[i] = zmath.V(math.Cos(angle), math.Sin(angle))
					} else { // otherwise, generate one for this and future use
						angle = rand.Float64() * 2.0 * math.Pi
						vectors[i] = zmath.V(math.Cos(angle), math.Sin(angle))
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
					influence := math.Pow(math.Max(0, r2-dist*dist), 4) // HEY, REDUCE POWER FOR COOL EFFECT!!!
					Z += influence * disp[i].Dot(vectors[i])
				}

				simplexMap[x][y] += Z * octInfluence
			}
		}
		fmt.Println("Octave ", oct+1, " finished.")
	}

	if cfg.Normalize {
		simplexMap.Interpolate(0, 1)
	}

	return simplexMap
}

func skew(vec zmath.Vec) zmath.Vec {
	return zmath.Vec{
		X: vec.X + (vec.X+vec.Y)*F2D,
		Y: vec.Y + (vec.X+vec.Y)*F2D,
	}
}

func unskew(vec zmath.Vec) zmath.Vec {
	return zmath.Vec{
		X: vec.X - (vec.X+vec.Y)*G2D,
		Y: vec.Y - (vec.X+vec.Y)*G2D,
	}
}

// getKey returns a unique string key for the given integer coordinates
func getKey(vec zmath.VecInt) string {
	return strconv.Itoa(vec.X) + "." + strconv.Itoa(vec.Y)
}
