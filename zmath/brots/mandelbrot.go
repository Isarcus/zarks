package brots

import (
	"math/cmplx"

	"github.com/Isarcus/zarks/zmath"
)

// NewMandelbrot returns a new, zoomed-out mandelbrot set
func NewMandelbrot(cfg Config) zmath.Map {
	m := zmath.NewMap(cfg.Res, 0)

	var (
		dx = cfg.Bounds.Dx()
		dy = cfg.Bounds.Dy()
	)

	for x := 0; x < cfg.Res.X; x++ {
		for y := 0; y < cfg.Res.Y; y++ {
			coord := zmath.Vec{
				X: cfg.Bounds.Min.X + (dx * float64(x) / float64(cfg.Res.X)),
				Y: cfg.Bounds.Min.Y + (dy * float64(y) / float64(cfg.Res.Y)),
			}.Complex()

			isInside, tries := TestPoint(coord, cfg.Pow, cfg.Iter)
			if !isInside {
				m[x][y] = float64(tries)
			}
		}
	}

	return m
}

// TestPoint will return whether the passed point is part of the Mandelbrot set. If the point is not part
// of the Mandelbrot set, TestPoint's returned int is the number if iterations it took to determine that.
func TestPoint(point complex128, pow complex128, iter int) (isInside bool, tries int) {
	var val complex128
	isInside = true
	for tries < iter {
		tries++

		val = cmplx.Pow(val, pow) + point

		if cmplx.Abs(val) > 2 {
			isInside = false
			return
		}
	}

	return
}
