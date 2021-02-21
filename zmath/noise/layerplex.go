package noise

import (
	"fmt"
	"image"
	"math"
	"math/rand"

	"github.com/Isarcus/zarks/zimg"
	"github.com/Isarcus/zarks/zmath"
)

// NewLayerplexMap returns a newly generated simplex map, with a slightly modified algorithm.
// The data within 'layer' is used to modify the magnitude of the simplex vectors, for an interesting effect!
func NewLayerplexMap(cfg Config, layer zmath.Map) zmath.Map {
	cfg.checkDefaults()
	r2 := cfg.R * cfg.R
	layerplexMap := zmath.NewMap(cfg.Dimensions, 0)

	conv := zmath.V(
		float64(layer.Bounds().X)/float64(cfg.Dimensions.X),
		float64(layer.Bounds().Y)/float64(cfg.Dimensions.Y),
	)

	for oct := 0.0; int(oct) < cfg.Octaves; oct++ {
		// create a new hashmap for angles
		vecMap := make(map[string]float64)
		octInfluence := math.Pow(0.5, oct)

		// dive in
		for x := 0; x < cfg.Dimensions.X; x++ {
			for y := 0; y < cfg.Dimensions.Y; y++ {
				iptX := float64(x) / (cfg.BoxSizeInitial * octInfluence)
				iptY := float64(y) / (cfg.BoxSizeInitial * octInfluence)
				ipt := zmath.V(iptX, iptY)

				layerX := int(float64(x) * conv.X)
				layerY := int(float64(y) * conv.Y)

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
						vectors[i] = zmath.V(math.Cos(angle), math.Sin(angle)).Scale(layer[layerX][layerY])
					} else { // otherwise, generate one for this and future use
						angle = rand.Float64() * 2.0 * math.Pi
						vectors[i] = zmath.V(math.Cos(angle), math.Sin(angle)).Scale(layer[layerX][layerY])
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
		layerplexMap.Interpolate(0, 1)
	}

	return layerplexMap
}

// NewStratoplexMap returns the result of some very cool magic
func NewStratoplexMap(cfg Config, img *image.RGBA) *zimg.ZImage {
	cfg.checkDefaults()
	cfg.Dimensions = zmath.VI(img.Rect.Dx(), img.Rect.Dy())
	r2 := cfg.R * cfg.R

	var (
		inputZi = zimg.ZImageFromRGBA(img)
		zi      = zimg.NewZImage(cfg.Dimensions)
	)

	for oct := 0.0; oct < float64(cfg.Octaves); oct++ {
		octInfluence := math.Pow(0.5, oct)

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
				corners[0] = zmath.V(math.Floor(skewed.X), math.Floor(skewed.Y))
				internal := zmath.V(skewed.X-corners[0].X, skewed.Y-corners[0].Y)
				if internal.X > internal.Y {
					corners[1] = corners[0].Add(zmath.V(1, 0))
				} else {
					corners[1] = corners[0].Add(zmath.V(0, 1))
				}
				corners[2] = corners[0].Add(zmath.V(1, 1))

				// get vectors of the three corners
				var vectors [3][3]zmath.Vec
				for color := 0; color < 3; color++ {
					slope := [3]float64{
						inputZi.RGBA256[color].SlopeAt(corners[0].VI()),
						inputZi.RGBA256[color].SlopeAt(corners[1].VI()),
						inputZi.RGBA256[color].SlopeAt(corners[2].VI()),
					}
					ang := [3]float64{
						inputZi.RGBA256[color].GradientAt(corners[0].VI()) + math.Pi,
						inputZi.RGBA256[color].GradientAt(corners[1].VI()) + math.Pi,
						inputZi.RGBA256[color].GradientAt(corners[2].VI()) + math.Pi,
					}
					vectors[color] = [3]zmath.Vec{
						zmath.V(math.Cos(ang[0]), math.Sin(ang[0])).Scale(slope[0]),
						zmath.V(math.Cos(ang[1]), math.Sin(ang[1])).Scale(slope[1]),
						zmath.V(math.Cos(ang[2]), math.Sin(ang[2])).Scale(slope[2]),
					}
				}

				// determine distance to each unskewed corner
				var disp [3]zmath.Vec
				var distance [3]float64
				for i := range distance {
					disp[i] = ipt.Subtract(unskew(corners[i]))
					distance[i] = math.Sqrt(disp[i].X*disp[i].X + disp[i].Y*disp[i].Y)
				}

				for color := 0; color < 3; color++ {
					var Z float64
					for i, dist := range distance {
						influence := math.Pow(math.Max(0, r2-dist*dist), 4) // HEY, REDUCE POWER FOR COOL EFFECT!!!
						Z += influence * disp[i].Dot(vectors[color][i])
					}
					zi.RGBA256[color][x][y] += Z
				}
			}
		}
		fmt.Println("Octave ", oct, " finished.")
	}

	zi.Interpolate(0, 255)
	zi.MakeOpaque()
	return zi
}
