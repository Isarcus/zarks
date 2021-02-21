package brots

import (
	"math/cmplx"

	"github.com/Isarcus/zarks/zmath"
)

// NewBuddhaBrot will return a strange-looking map of a modified Mandelbrot algorithm!
func NewBuddhaBrot(cfg Config) zmath.Map {
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

			isInside, _, path := TestPointTrace(coord, cfg.Pow, cfg.Iter)

			if !isInside {
				for _, ptComplex := range path {
					pt := zmath.VC(ptComplex)

					ptVI := zmath.VI(
						int((pt.X-cfg.Bounds.Min.X)/dx*float64(cfg.Res.X)),
						int((pt.Y-cfg.Bounds.Min.Y)/dy*float64(cfg.Res.Y)),
					)

					if m.ContainsCoord(ptVI) {
						m.Set(ptVI, m.At(ptVI)+1)
					}
				}
			}
		}
	}

	return m
}

// TestPointTrace will apply the Mandelbrot algorithm on the point applied, and return three values:
// 1. Whether the tested point is inside the Mandelbrot set
// 2. How many tries it took to determine it was NOT in the Mandelbrot set (applies only to non-Mandelbrot points)
// 3. The path the point took through space
func TestPointTrace(point complex128, pow complex128, iter int) (isInside bool, tries int, path []complex128) {
	path = make([]complex128, 0, iter)
	isInside = true

	var val complex128
	for tries < iter {
		tries++
		val = cmplx.Pow(val, pow) + point

		if cmplx.Abs(val) > 2 {
			isInside = false
			break
		}
		path = append(path, val)
	}

	return
}

func addLine(val float64) float64 { return val + 1 }
