package erosion

import (
	"math/rand"
	"time"
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
