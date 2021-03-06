package zmath

import (
	"encoding/binary"
	"fmt"
	"image"
	"math"
	"os"

	"github.com/Isarcus/zarks/system"
	"github.com/Isarcus/zarks/zmath/zbits"
)

// Map is a set of 2D raster data, with some helpful member functions
type Map [][]float64

// NewMap returns a 2D array of float64 of the given bounds, with all cells set to the given initial value
func NewMap(bounds VecInt, initValue float64) Map {
	x, y := bounds.X, bounds.Y
	data := make([][]float64, x)
	for i := 0; i < x; i++ {
		data[i] = make([]float64, y)
		for j := 0; j < y; j++ {
			data[i][j] = initValue
		}
	}

	return data
}

// At returns the value of the map at the coordinates specified by the passed VecInt. It does NOT bounds-check!
func (m Map) At(pos VecInt) float64 {
	return m[pos.X][pos.Y]
}

// Set sets the point in the map at the desired coordinates to the passed value. It does NOT bounds-check!
func (m Map) Set(pos VecInt, value float64) {
	m[pos.X][pos.Y] = value
}

// PtrTo returns a pointer to the map index at the desired coordinates
func (m Map) PtrTo(pos VecInt) *float64 {
	return &(m[pos.X][pos.Y])
}

// Bounds returns the bounds of a map
func (m Map) Bounds() VecInt {
	return VecInt{
		X: len(m),
		Y: len(m[0]),
	}
}

// Area returns the area of the map as determined by its width and height.
func (m Map) Area() float64 {
	b := m.Bounds()
	return float64(b.X * b.Y)
}

// Clear sets all points on the map equal to the passed value
func (m Map) Clear(value float64) Map {
	for x, row := range m {
		for y := range row {
			m[x][y] = value
		}
	}
	return m
}

// Zero wipes the called map entirely, setting all values to 0
func (m Map) Zero() Map {
	return m.Clear(0)
}

// GetSum returns the sum of all map elements
func (m Map) GetSum() float64 {
	sum := 0.0

	for x, row := range m {
		for y := range row {
			sum += m[x][y]
		}
	}

	return sum
}

// GetMean returns the mean of all map elements
func (m Map) GetMean() float64 {
	sum := m.GetSum()
	return sum / float64(m.Bounds().X*m.Bounds().Y)
}

// GetMin returns the minimum of all map elements
func (m Map) GetMin() float64 {
	if m.Bounds().X > 0 && m.Bounds().Y > 0 {
		min := m[0][0]
		for _, row := range m {
			for i := range row {
				min = math.Min(min, row[i])
			}
		}

		return min
	}

	return 0
}

// GetMax returns the maximum of all map elements
func (m Map) GetMax() float64 {
	if m.Bounds().X > 0 && m.Bounds().Y > 0 {
		max := m[0][0]
		for _, row := range m {
			for i := range row {
				max = math.Max(max, row[i])
			}
		}

		return max
	}

	return 0
}

// GetMinMax returns the min and max of the called Map.
// This is computationally faster than calling GetMin and GetMax independently.
func (m Map) GetMinMax() (min, max float64) {
	if m.Bounds().X > 0 && m.Bounds().Y > 0 {
		min, max = m[0][0], m[0][0]
		for _, row := range m {
			for _, val := range row {
				min = math.Min(min, val)
				max = math.Max(max, val)
			}
		}
	}
	return
}

// GetRange returns the range of the called Map.
func (m Map) GetRange() float64 {
	min, max := m.GetMinMax()
	return max - min
}

// Add the passed value to every datapoint
func (m Map) Add(addend float64) Map {
	for x, row := range m {
		for y := range row {
			m[x][y] += addend
		}
	}
	return m
}

// Subtract the passed value from every datapoint
func (m Map) Subtract(subtrahend float64) Map {
	for x, row := range m {
		for y := range row {
			m[x][y] -= subtrahend
		}
	}
	return m
}

// ScaleDim scales a map's dimensions by an amount. For example, scaling by 2 will double an image's dimensions.
// This will lose data when shrinking an image
func (m Map) ScaleDim(by float64) Map {
	newBounds := m.Bounds().V().Scale(by).VI()
	newMap := NewMap(newBounds, 0)
	inv := 1.0 / by
	for x, row := range newMap {
		for y := range row {
			newMap[x][y] = m.At(VI(x, y).V().Scale(inv).VI())
		}
	}
	return newMap
}

// Multiply every data point by the passed value
func (m Map) Multiply(multiplicand float64) Map {
	for x, row := range m {
		for y := range row {
			m[x][y] *= multiplicand
		}
	}
	return m
}

// AddMap adds the passed map to the called map
func (m Map) AddMap(addend Map) Map {
	if m.Bounds() != addend.Bounds() {
		fmt.Println("Bounds don't match")
		return m
	}
	for x, row := range m {
		for y := range row {
			m[x][y] += addend[x][y]
		}
	}
	return m
}

// SubtractMap subtracts the passed map from the called map
func (m Map) SubtractMap(subtrahend Map) Map {
	if m.Bounds() != subtrahend.Bounds() {
		return m
	}
	for x, row := range m {
		for y := range row {
			m[x][y] -= subtrahend[x][y]
		}
	}
	return m
}

// GeometricMean calculates the geometric mean of two maps
func (m Map) GeometricMean(by Map) Map {
	if m.Bounds() != by.Bounds() {
		return m
	}

	for x, row := range m {
		for y := range row {
			m[x][y] = math.Sqrt(math.Abs(m[x][y] * by[x][y]))
		}
	}

	return m
}

// Interpolate interpolates a map between two values
func (m Map) Interpolate(newMin, newMax float64) Map {
	linear := m.ToLinear()
	linear.Interpolate(newMin, newMax)
	linear.To2D(m)
	return m
}

// MakeUniform makes the called Map follow a uniform distribution.
// This will not change the original min and max values of the called map.
func (m Map) MakeUniform() Map {
	min, max := m.GetMinMax()
	linear := m.ToLinear()
	linear.MakeUniform()
	linear.To2D(m)
	m.Interpolate(min, max)

	return m
}

// FlipVertical flips the map across its X-axis..
func (m Map) FlipVertical() Map {
	for x, row := range m {
		maxY := m.Bounds().Y - 1
		for y := range row {
			if y >= maxY {
				break
			}
			m[x][y], m[x][maxY] = m[x][maxY], m[x][y]
			maxY--
		}
	}
	return m
}

// FlipHorizontal flips the map across its Y-axis.
func (m Map) FlipHorizontal() Map {
	maxX := m.Bounds().X - 1
	for x, row := range m {
		if x >= maxX {
			break
		}
		for y := range row {
			m[x][y], m[maxX][y] = m[maxX][y], m[x][y]
		}
		maxX--
	}
	return m
}

// RotateCW90 rotates a map clockwise, 90 degrees. If the called map is square, the underlying array will be used.
// Otherwise, a new Map will be created and the original discarded.
// func (m Map) RotateCW90() Map {
// 	bounds := m.Bounds()

// 	// if square
// 	if bounds.X == bounds.Y {

// 	} else {
// 		newBounds := VecInt{
// 			X: bounds.Y,
// 			Y: bounds.X,
// 		}
// 		rotated := NewMap(newBounds, 0)

// 		for x, row := range m {
// 			for y := range row {
// 				rotated[y][x]
// 			}
// 		}
// 	}

// 	return m
// }

// SetMin sets every value in the map that is less than the passed value TO the passed value.
// This function does *NOT* interpolate the map between its maximum and a new minimum.
func (m Map) SetMin(value float64) Map {
	for x, row := range m {
		for y := range row {
			if m[x][y] < value {
				m[x][y] = value
			}
		}
	}
	return m
}

// SetMax sets every value in the map that is greater than the passed value TO the passed value.
// This function does *NOT* interpolate the map between its minimum and a new maximum.
func (m Map) SetMax(value float64) Map {
	for x, row := range m {
		for y := range row {
			if m[x][y] > value {
				m[x][y] = value
			}
		}
	}
	return m
}

// Replace replaces all occurrences of a value with another value
func (m Map) Replace(value, with float64) Map {
	for x, row := range m {
		for y := range row {
			if m[x][y] == value {
				m[x][y] = with
			}
		}
	}
	return m
}

// ReplaceNot replaces all data that is NOT equal to the passed value with another value
func (m Map) ReplaceNot(value, with float64) Map {
	for x, row := range m {
		for y := range row {
			if m[x][y] != value {
				m[x][y] = with
			}
		}
	}
	return m
}

// CustomMod lets you pass in ANY function that takes in a float64 and returns a float64, which it will
// then use to modify every value of the called map.
func (m Map) CustomMod(modFunc func(float64) float64) Map {
	for x, row := range m {
		for y := range row {
			m[x][y] = modFunc(m[x][y])
		}
	}
	return m
}

// CustomModAt lets you pass in ANY function that takes in a float64 and returns a float64, which it will
// then use to modify the values ONLY at the specified points
func (m Map) CustomModAt(points []VecInt, modFunc func(float64) float64) Map {
	for _, p := range points {
		if m.ContainsCoord(p) {
			m[p.X][p.Y] = modFunc(m[p.X][p.Y])
		}
	}
	return m
}

// Copy returns a deepcopy of the portion of the map specified by the min and max coordinates.
// If min > max, a map with zero bounds in either direction will be returned.
// If min or max exceed the bounds of the called map, any out-of-bounds coordinates will be set to zero.
// To copy the entire map, use `m.Copy(zmath.ZVI, m.Bounds())`.
func (m Map) Copy(min, max VecInt) Map {
	var dim = VecInt{}
	if max.X > min.X && max.Y > min.Y {
		dim = VecInt{
			X: max.X - min.X,
			Y: max.Y - min.Y,
		}
	}

	copyMap := NewMap(dim, 0)
	for x := 0; x < dim.X; x++ {
		for y := 0; y < dim.Y; y++ {
			pos := VecInt{
				X: min.X + x,
				Y: min.Y + y,
			}
			if m.ContainsCoord(pos) {
				copyMap[x][y] = m.At(pos)
			}
		}
	}

	return copyMap
}

// CopyAll deepcopies an entire map and returns the copy.
// Equivalent to calling the Copy() function for a map's entire area.
func (m Map) CopyAll() Map {
	return m.Copy(ZVI, m.Bounds())
}

// Paste pastes all of the data from pasteMe into a map, starting at the specified coordinate in the called map.
// If you only want to copy part of the passed Map, Copy() part of it first.
func (m Map) Paste(pasteMe Map, at VecInt) Map {
	b := pasteMe.Bounds()
	for x := 0; x < b.X; x++ {
		for y := 0; y < b.Y; y++ {
			pos := VecInt{
				X: at.X + x,
				Y: at.Y + y,
			}
			if m.ContainsCoord(pos) {
				m.Set(pos, pasteMe[x][y])
			}
		}
	}
	return m
}

// BlurGaussian blurs gaussly (it looks nice)
func (m Map) BlurGaussian(radius int) Map {
	// For later
	var (
		circle     = GetCircleCoords(radius)
		bounds     = m.Bounds()
		sigma2     = math.Pow(float64(radius)/3.0, 2)
		gaussCoeff = 1.0 / (2.0 * math.Pi * sigma2)
	)

	blurMap := NewMap(bounds, 0)

	for x := 0; x < bounds.X; x++ {
		for y := 0; y < bounds.Y; y++ {
			var weightSum, trueWeightSum, sum float64
			pos := VecInt{
				X: x,
				Y: y,
			}

			for _, cpt := range circle {
				blurPos := pos.Add(cpt)
				weight := math.Pow(math.E, -float64(cpt.X*cpt.X+cpt.Y*cpt.Y)/(2.0*sigma2))
				trueWeightSum += weight
				if m.ContainsCoord(blurPos) {
					weightSum += weight
					sum += weight * m.At(blurPos)
				}
			}

			ratio := gaussCoeff * trueWeightSum / weightSum
			sum *= ratio
			blurMap[x][y] = sum
		}
	}
	m.Paste(blurMap, ZVI)
	return m
}

// ToLinear converts a Map to a Set
func (m Map) ToLinear() Set {
	linear := make(Set, 0)
	for i := range m {
		linear = append(linear, m[i]...)
	}
	return linear
}

// ContainsCoord tells you whether the specified coordinate is inside the called map
func (m Map) ContainsCoord(pos VecInt) bool {
	b := m.Bounds()
	return pos.X >= 0 && pos.Y >= 0 && pos.X < b.X && pos.Y < b.Y
}

// DerivativeAt returns the derivative at the desired point on the map as a Vec{X: dh/dx, Y: dh/dy}
// This function has built-in bounds-checking and will return {0, 0} if out of bounds.
// The formula used for derivative calculations can be found at:
// https://desktop.arcgis.com/en/arcmap/10.3/tools/spatial-analyst-toolbox/how-slope-works.htm
func (m Map) DerivativeAt(pos VecInt) Vec {
	if !m.ContainsCoord(pos) {
		return Vec{}
	}

	var (
		posX = m.ContainsCoord(pos.AddXY(1, 0))
		negX = m.ContainsCoord(pos.AddXY(-1, 0))
		posY = m.ContainsCoord(pos.AddXY(0, 1))
		negY = m.ContainsCoord(pos.AddXY(0, -1))

		dh     = Vec{}
		weight = Vec{}

		val = m.At(pos)
	)

	if posX {
		dh.X += 2.0 * (m.At(pos.AddXY(1, 0)) - val)
		weight.X += 2

		if posY {
			dh.X += m.At(pos.AddXY(1, 1)) - val
			dh.Y += m.At(pos.AddXY(1, 1)) - val

			weight.X++
			weight.Y++
		}
		if negY {
			dh.X += m.At(pos.AddXY(1, -1)) - val
			dh.Y -= m.At(pos.AddXY(1, -1)) - val

			weight.X++
			weight.Y++
		}
	}
	if negX {
		dh.X -= 2.0 * (m.At(pos.AddXY(-1, 0)) - val)
		weight.X += 2

		if posY {
			dh.X -= m.At(pos.AddXY(-1, 1)) - val
			dh.Y += m.At(pos.AddXY(-1, 1)) - val

			weight.X++
			weight.Y++
		}
		if negY {
			dh.X -= m.At(pos.AddXY(-1, -1)) - val
			dh.Y -= m.At(pos.AddXY(-1, -1)) - val

			weight.X++
			weight.Y++
		}
	}
	if posY {
		dh.Y += 2.0 * (m.At(pos.AddXY(0, 1)) - val)
		weight.Y += 2
	}
	if negY {
		dh.Y -= 2.0 * (m.At(pos.AddXY(0, -1)) - val)
		weight.Y += 2
	}

	dh.X /= weight.X
	dh.Y /= weight.Y

	return dh
}

// GradientAt returns the angle of the gradient at the desired coordinate using Atan2
func (m Map) GradientAt(pos VecInt) float64 {
	d := m.DerivativeAt(pos)
	return math.Atan2(d.Y, d.X)
}

// SlopeAt returns the slope at the point
func (m Map) SlopeAt(pos VecInt) float64 {
	dh := m.DerivativeAt(pos)
	return DistanceFormula(dh, ZV)
}

// GetSlopeMap returns a NEW map of slopes.
func (m Map) GetSlopeMap() Map {
	slopeMap := NewMap(m.Bounds(), 0)

	for x, row := range m {
		for y := range row {
			slopeMap[x][y] = m.SlopeAt(VecInt{x, y})
		}
	}

	return slopeMap
}

//                              //
// - - - IMAGE CONVERSION - - - //
//                              //

// Color identifies colors
type Color int

// Colors
const (
	Red Color = iota // I don't really like having enumerated colors in the zmath package, but eh they can stay for now
	Green
	Blue
	Black
)

// MapFromImage returns a map of the R, G, B, or brightness values of an image
func MapFromImage(img *image.RGBA, color Color) Map {
	imgMap := NewMap(VI(img.Bounds().Max.X, img.Bounds().Max.Y), 0)

	for x := 0; x < imgMap.Bounds().X; x++ {
		for y := 0; y < imgMap.Bounds().Y; y++ {
			var pixel = img.RGBAAt(x, y)
			switch color {
			case Red:
				imgMap[x][y] = float64(pixel.R) / 255.0
			case Green:
				imgMap[x][y] = float64(pixel.G) / 255.0
			case Blue:
				imgMap[x][y] = float64(pixel.B) / 255.0
			case Black:
				imgMap[x][y] = (float64(pixel.R) + float64(pixel.G) + float64(pixel.B)) / 765.0
			}
		}
	}

	return imgMap
}

// ImageToMap returns a map of the R, G, B, or brightness values of an image
var ImageToMap = MapFromImage

// Save saves a Map as binary data at the path specified. File ending should be .zmap
func (m Map) Save(path string) {
	f := system.CreateFile(path)
	defer f.Close()

	// Header (64)
	header := [64]byte{}
	system.WriteBytes(f, header[:])

	// bounds (8)
	var dimX = uint32(m.Bounds().X)
	var dimY = uint32(m.Bounds().Y)
	system.WriteBytes(f, zbits.Uint32ToBytes(dimX, zbits.LE))
	system.WriteBytes(f, zbits.Uint32ToBytes(dimY, zbits.LE))

	// Data (8n)
	for x, row := range m {
		for y := range row {
			bits := math.Float64bits(m[x][y])
			system.WriteBytes(f, zbits.Uint64ToBytes(bits, zbits.LE))
		}
	}
}

// MapFromPath loads a Map from the target path
func MapFromPath(path string) Map {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return Map{{}}
	}

	// Header (64)
	hdr := [64]byte{}
	f.Read(hdr[:])

	// Bounds (8)
	bounds := [8]byte{}
	f.Read(bounds[:])

	dimX := binary.LittleEndian.Uint32(bounds[0:4])
	dimY := binary.LittleEndian.Uint32(bounds[4:8])
	newMap := NewMap(VI(int(dimX), int(dimY)), 0)

	// Data (8n)
	for x, row := range newMap {
		for y := range row {
			bytes := [8]byte{}
			f.Read(bytes[:])
			newMap[x][y] = math.Float64frombits(binary.LittleEndian.Uint64(bytes[:]))
		}
	}

	return newMap
}

//                    //
// - - - MAPVEC - - - //
//                    //

// MapVec is a container for a 2D array of vectors
type MapVec [][]Vec

// NewMapVec returns a zeroed MapVec of the given bounds
func NewMapVec(bounds VecInt) MapVec {
	data := make([][]Vec, bounds.X)
	for i := 0; i < bounds.X; i++ {
		data[i] = make([]Vec, bounds.Y)
		for j := 0; j < bounds.Y; j++ {
			data[i][j] = ZV
		}
	}

	return data
}

//                                //
// - - - NEEDLESS ITERATION - - - //
//                                //

var iterhash = make(map[*float64]VecInt)

// Iterate will iterate through a map if for some reason you really, Really, REALLY don't want to nest two for loops.
// Use the syntax: `for ok, ptr := m.Iterate(); ok; ok, ptr = m.Iterate()`. This is exponentially slower than
// looping through a 2D array the traditional way, but I just wanted to see if this was possible :-)
func (m Map) Iterate() (bool, *float64) {
	var (
		bounds  = m.Bounds()
		zeroPtr = &(m[0][0])
		nextPos = VecInt{}
	)

	// if first iteration
	val, ok := iterhash[zeroPtr]
	if !ok {
		if bounds.X > 1 {
			nextPos.X = 1
		} else {
			nextPos.Y = 1
		}
		iterhash[zeroPtr] = nextPos
		return true, zeroPtr
	}

	if m.ContainsCoord(val) {
		if val.X+1 >= bounds.X {
			nextPos.Y = val.Y + 1
		} else {
			nextPos.X = val.X + 1
			nextPos.Y = val.Y
		}
		iterhash[zeroPtr] = nextPos
		return true, m.PtrTo(val)
	}

	delete(iterhash, zeroPtr)
	return false, nil
}
