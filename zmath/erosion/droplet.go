package erosion

import "github.com/Isarcus/zarks/zmath"

type droplet struct {
	pos zmath.Vec
	dir zmath.Vec
	vel float64

	water    float64
	sediment float64
}
