package brots

import (
	"math"

	"github.com/Isarcus/zarks/zmath"
)

// NewTestBrot brots a brot brotly
func NewTestBrot(cfg Config) zmath.Map {
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

					dist := math.Min(3, 1.0/zmath.DistanceFormula(prev, current))

					if m.ContainsCoord(prevVI) {
						m.Set(prevVI, m.At(prevVI)+dist)
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
