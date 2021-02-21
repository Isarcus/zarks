package noise

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/Isarcus/zarks/zmath"
)

// NewPerlinMap generates a new perlin noise map according to the specified configuration.
func NewPerlinMap(cfg Config) zmath.Map {
	cfg.checkDefaults()
	perlinMap := zmath.NewMap(cfg.Dimensions, 0)

	var cornerVecs = [4]zmath.Vec{
		zmath.V(0, 0),
		zmath.V(1, 0),
		zmath.V(1, 1),
		zmath.V(0, 1),
	}

	// The Loop
	for oct := 1.0; oct <= float64(cfg.Octaves); oct++ {
		boxSize := cfg.BoxSizeInitial / math.Pow(2, oct)
		boxesX := math.Ceil(float64(cfg.Dimensions.X)/boxSize) + 1
		boxesY := math.Ceil(float64(cfg.Dimensions.Y)/boxSize) + 1
		vectors := zmath.NewMapVec(zmath.VI(int(boxesX+1), int(boxesY+1)))

		for bx := 0; bx < int(boxesX); bx++ {
			for by := 0; by < int(boxesY); by++ {
				// Generate random vectors
				for _, cornerVec := range cornerVecs {
					if vectors[bx+int(cornerVec.X)][by+int(cornerVec.Y)] == zmath.ZV {
						theta := 2.0 * rand.Float64() * math.Pi
						vectors[bx+int(cornerVec.X)][by+int(cornerVec.Y)] = zmath.Vec{
							X: math.Cos(theta),
							Y: math.Sin(theta),
						}
						//fmt.Println(theta)
					}
				}

				// Calculate current box size
				boxStart := zmath.Vec{
					X: math.Floor(float64(bx) * boxSize),
					Y: math.Floor(float64(by) * boxSize),
				}
				boxSize := zmath.Vec{
					X: math.Floor(float64(bx+1)*boxSize) - boxStart.X,
					Y: math.Floor(float64(by+1)*boxSize) - boxStart.Y,
				}

				// Loop within one box
				for dx := 0.0; dx < boxSize.X; dx++ {
					dataX := int(boxStart.X + dx)
					if dataX < 0 || dataX >= cfg.Dimensions.X {
						continue
					}
					for dy := 0.0; dy < boxSize.Y; dy++ {
						dataY := int(boxStart.Y + dy)
						if dataY < 0 || dataY >= cfg.Dimensions.Y {
							continue
						}

						boxPos := zmath.Vec{
							X: dx / boxSize.X,
							Y: dy / boxSize.Y,
						}

						dot00 := boxPos.Subtract(cornerVecs[0]).Dot(vectors[bx][by])
						dot10 := boxPos.Subtract(cornerVecs[1]).Dot(vectors[bx+1][by])
						dot11 := boxPos.Subtract(cornerVecs[2]).Dot(vectors[bx+1][by+1])
						dot01 := boxPos.Subtract(cornerVecs[3]).Dot(vectors[bx][by+1])

						x0 := interpolatePow5(dot00, dot10, boxPos.X)
						x1 := interpolatePow5(dot01, dot11, boxPos.X)

						Z := interpolatePow5(x0, x1, boxPos.Y)

						perlinMap[dataX][dataY] += Z * math.Pow(0.5, oct)
					}
				}
			}
		}
		fmt.Println("Octave ", oct, " finished.")
	}

	// Post-analysis
	if cfg.Normalize {
		perlinMap.Interpolate(0, 1)
	}

	return perlinMap
}

func interpolatePow5(i0, i1, t float64) float64 {
	weight := 6*math.Pow(t, 5) - 15*math.Pow(t, 4) + 10*math.Pow(t, 3)
	return weight*i1 + (1.0-weight)*i0
}
