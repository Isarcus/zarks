package zmath

import "math"

// GetCircleCoords returns an array of vectors representing all points within a circle of the specified radius.
// The list is edge-inclusive.
func GetCircleCoords(radius int) []VecInt {
	points := make([]VecInt, 0, int(radius*radius*4)) // because pi = 4, more or less

	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			pos := VI(x, y)
			if DistanceFormulaInt(VI(x, y), ZVI) <= float64(radius) {
				points = append(points, pos)
			}
		}
	}

	return points[:]
}

// GetLineCoords will return all the integer points contained by a line connecting the two passed points
func GetLineCoords(point1, point2 VecInt) []VecInt {
	points := make([]VecInt, 0)

	slope := point1.Slope(point2)

	if slope == math.Inf(1) {
		var (
			minY = MinInt(point1.Y, point2.Y)
			maxY = MaxInt(point1.Y, point2.Y)
		)

		for y := minY; y < maxY; y++ {
			points = append(points, VI(point1.X, y))
		}
	} else {
		for x := point1.X; x < point2.X; x++ {
			var (
				yNow  = int(float64(point1.Y) + float64(x)*slope)
				yNext = int(float64(point1.Y) + float64(x+1)*slope)
			)

			for y := yNow; y <= yNext; y++ {
				points = append(points, VI(x, y))
			}
		}
	}

	return points
}
