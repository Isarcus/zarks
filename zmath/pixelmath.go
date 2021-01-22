package zmath

import "github.com/faiface/pixel"

// PVecToZVec converts a pixel.Vec to a zmath.Vec
func PVecToZVec(v pixel.Vec) Vec {
	return Vec{
		X: v.X,
		Y: v.Y,
	}
}

// PVecToZVecInt converts a pixel.Vec to a zmath.Vec
func PVecToZVecInt(v pixel.Vec) VecInt {
	return VecInt{
		X: int(v.X),
		Y: int(v.Y),
	}
}
