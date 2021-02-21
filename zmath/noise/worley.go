package noise

import "github.com/Isarcus/zarks/zmath"

// NewWorleyMap returns a new Worley Noise map (WORK IN PROGRESS)
func NewWorleyMap(cfg Config) zmath.Map {
	cfg.checkDefaults()
	noiseMap := zmath.NewMap(cfg.Dimensions, 0)

	var neighbors = make([]zmath.VecInt, 21)
	idx := 0
	for x := -2; x <= 2; x++ {
		for y := -2; y <= 2; y++ {
			if x*y != 4 && x*y != -4 { //&& (x != 0 || y != 0) {
				neighbors[idx] = zmath.VI(x, y)
				idx++
			}
		}
	}

	//ptMap := make(map[string]zmath.Vec)

	for oct := 0; oct < cfg.Octaves; oct++ {
		//octInfluence := math.Pow(0.5, float64(oct))
		//boxSize := cfg.BoxSizeInitial * octInfluence

		for x := 0; x < cfg.Dimensions.X; x++ {
			for y := 0; y < cfg.Dimensions.Y; y++ {
				/*box := zmath.VecInt{
					X: int(float64(x) / boxSize),
					Y: int(float64(y) / boxSize),
				}
				//pos := zmath*/

			}
		}
	}

	return noiseMap
}
