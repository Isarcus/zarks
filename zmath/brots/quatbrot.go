package brots

import "github.com/Isarcus/zarks/zmath"

// NewQuatBrot will use quaternion math to apply a MandelQuat formula!
func NewQuatBrot(cfg Config) zmath.Map {
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
			}

			quat := zmath.Quat{
				A: coord.X,
				I: coord.Y,
				J: -coord.Y,
				K: -coord.X,
			}

			isInside, tries := TestQuat(quat, int(real(cfg.Pow)), cfg.Iter)
			//fmt.Println(tries)
			if !isInside {
				m[x][y] = float64(tries)
			}
		}
	}

	return m
}

// TestQuat will return whether the passed quaternion is part of the MandelQuat set
func TestQuat(point zmath.Quat, pow int, iter int) (isInside bool, tries int) {
	var val = zmath.ZQ
	isInside = true
	for tries < iter {
		tries++

		val = val.PowInt(pow)
		val = val.Add(point)

		if val.Abs() > 2 {
			isInside = false
			return
		}
	}

	return
}
