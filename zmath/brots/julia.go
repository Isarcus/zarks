package brots

import (
	"math"
	"math/cmplx"

	"github.com/Isarcus/zarks/zmath"
)

// JuliaSet returns a new Julia Set!
func JuliaSet(cfg Config) zmath.Map {
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

			isInside, tries := TestPointJulia(coord, cfg.Pow, cfg.C, cfg.Iter)
			if !isInside {
				m[x][y] = float64(tries)
			}
		}
	}

	return m
}

// TestPointJulia tests a point using the Julia Set algorithm
func TestPointJulia(point complex128, pow complex128, c complex128, iter int) (isInside bool, tries int) {
	isInside = true

	var (
		val    complex128
		escape = math.Max(2, math.Sqrt(cmplx.Abs(c)))
	)

	for tries < iter {
		tries++

		val = cmplx.Pow(val, pow) + c

		if cmplx.Abs(val) > escape {
			isInside = false
			break
		}
	}
	return
}
