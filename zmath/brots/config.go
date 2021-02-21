package brots

import "github.com/Isarcus/zarks/zmath"

// Config contains all of the information necessary to generate a custom Mandelbrot map
type Config struct {
	Bounds *zmath.Rect  // The area to be tested
	Res    zmath.VecInt // The resolution of the map
	Iter   int          // How many iterations per point

	Pow complex128 // What power to put each point to to test them (2 for original Mandelbrot)
}

// DefaultConfig can be used to generate a classic zoomed-out Mandelbrot set
var DefaultConfig = Config{
	Bounds: zmath.R(zmath.V(-2, -1.5), zmath.V(1, 1.5)),
	Res:    zmath.VI(512, 512),
	Iter:   60,

	Pow: 2,
}
