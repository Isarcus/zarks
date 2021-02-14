package brots

import (
	"math/cmplx"
	"zarks/zmath"
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
				for p := range path {
					if p == 0 {
						continue
					}
					var (
						prev    = zmath.VC(path[p-1])
						current = zmath.VC(path[p])
					)

					prevVI := zmath.VI(
						int((prev.X-cfg.Bounds.Min.X)/dx*float64(cfg.Res.X)),
						int((prev.Y-cfg.Bounds.Min.Y)/dy*float64(cfg.Res.Y)),
					)

					currentVI := zmath.VI(
						int((current.X-cfg.Bounds.Min.X)/dx*float64(cfg.Res.X)),
						int((current.Y-cfg.Bounds.Min.Y)/dy*float64(cfg.Res.Y)),
					)

					if m.ContainsCoord(prevVI) {
						m.Set(prevVI, m.At(prevVI)+1)
					}
					if m.ContainsCoord(currentVI) {
						m.Set(currentVI, m.At(currentVI)+1)
					}
				}
			}
		}
	}

	return m
}

// NewLineBrot will return a strange-looking map of a modified Mandelbrot algorithm!
func NewLineBrot(cfg Config) zmath.Map {
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
				for p := range path {
					if p == 0 {
						continue
					}
					var (
						prev    = zmath.VC(path[p-1])
						current = zmath.VC(path[p])
					)

					prevVI := zmath.VI(
						int((prev.X-cfg.Bounds.Min.X)/dx*float64(cfg.Res.X)),
						int((prev.Y-cfg.Bounds.Min.Y)/dy*float64(cfg.Res.Y)),
					)

					currentVI := zmath.VI(
						int((current.X-cfg.Bounds.Min.X)/dx*float64(cfg.Res.X)),
						int((current.Y-cfg.Bounds.Min.Y)/dy*float64(cfg.Res.Y)),
					)

					lineCoords := zmath.GetLineCoords(prevVI, currentVI)

					m.CustomModAt(lineCoords, addLine)
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
