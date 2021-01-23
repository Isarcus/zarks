package zmath

// GetCircleCoords returns an array of vectors representing all points within a circle of the specified radius.
// The list is edge-inclusive.
func GetCircleCoords(radius int) []VecInt {
	points := make([]VecInt, 0, int(radius*radius*4)) // pi = 4, more or less

	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			pos := NewVecInt(x, y)
			if DistanceFormulaInt(NewVecInt(x, y), ZVI) <= float64(radius) {
				points = append(points, pos)
			}
		}
	}

	return points[:]
}
