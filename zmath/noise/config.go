package noise

import (
	"fmt"
	"math/rand"
	"time"
	"zarks/zmath"
)

// Config contains all of the necessary information to generate a noise map.
// Set seed to 0 for a random seed
type Config struct {
	Dimensions     zmath.VecInt // size of image
	Octaves        int          // number of octaves to iterate
	Normalize      bool         // whether to normalize data between [0,1] after generation
	Seed           int64        // what to seed the random number generator with. 0 for random seed!
	BoxSizeInitial float64      // Initial size of box or triangle in pixels

	R float64 // Simplex only - size of radius of influence for vector
	N int     // Worley only - how many nearest points to include in calculations
}

// DefaultConfig will give you a randomly seeded, 512x512 image with 4 octaves.
var DefaultConfig = Config{
	Dimensions:     zmath.NewVecInt(512, 512),
	Octaves:        4,
	Seed:           0,
	BoxSizeInitial: 128,
	Normalize:      true,

	R: 0.6,
	N: 3,
}

func (cfg *Config) checkDefaults() {
	// Initial calculations and initialization
	if cfg.Dimensions.X == 0 || cfg.Dimensions.Y == 0 {
		cfg.Dimensions = DefaultConfig.Dimensions
	}
	if cfg.Octaves == 0 {
		cfg.Octaves = DefaultConfig.Octaves
	}
	if cfg.Seed == 0 {
		seed := time.Now().Unix()
		rand.Seed(seed)
		fmt.Println("Using random seed: ", seed)
	} else {
		rand.Seed(cfg.Seed)
	}
	if cfg.BoxSizeInitial == 0 {
		cfg.BoxSizeInitial = DefaultConfig.BoxSizeInitial
	}
	if cfg.R == 0 {
		cfg.R = DefaultConfig.R
	}
	if cfg.N == 0 {
		cfg.N = DefaultConfig.N
	}
}
