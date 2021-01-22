package erosion

import (
	"math/rand"
	"time"
	"zarks/zmath"
)

// Config contains basic erosion parameters
type Config struct {
	Seed     int64
	Strength float64 // between 0 and 1
	Test     bool
}

var defaultConfig = Config{
	Seed:     0,
	Strength: 0.5,
}

func (cfg *Config) checkDefaults() {
	if cfg.Seed == 0 {
		rand.Seed(time.Now().Unix())
	}

	if cfg.Strength == 0 {
		cfg.Strength = defaultConfig.Strength
	}
}

// BasicErosion does cool erosion stuff on a Map!
func BasicErosion(m zmath.Map, cfg Config, iter int) zmath.Map {
	rand.Seed(cfg.Seed)
	bounds := m.Bounds()

	for i := 0; i < iter; i++ {
		pos := zmath.VecInt{
			X: int(rand.Float64() * float64(bounds.X)),
			Y: int(rand.Float64() * float64(bounds.Y)),
		}

		ct := 0
		for {
			minPos := getPosMinAdjacent(m, pos)
			//fmt.Println(ct, ": ", pos, " ", minPos)
			if minPos == pos { // if already at the lowest point, stop flowing
				break
			}
			dh := m.At(pos) - m.At(minPos)
			if cfg.Test {
				m[pos.X][pos.Y] += 100
			} else {
				m[pos.X][pos.Y] -= dh * cfg.Strength
			}

			pos = minPos
			ct++
		}
	}
	return m
}

func getPosMinAdjacent(m zmath.Map, pos zmath.VecInt) zmath.VecInt {
	bounds := m.Bounds()
	var min = m.At(pos)
	var minPos = pos

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			testPos := zmath.VecInt{
				X: pos.X + x,
				Y: pos.Y + y,
			}

			if testPos.X >= 0 && testPos.X < bounds.X && testPos.Y >= 0 && testPos.Y < bounds.Y && m.At(testPos) < min {
				minPos = testPos
				min = m.At(testPos)
			}
		}
	}

	return minPos
}
