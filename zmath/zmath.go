package zmath

import "math"

//								//
// - - - TYPE DEFINITIONS - - - //
//								//

// Vec is a vector. Nuff said.
type Vec struct {
	X, Y float64
}

// VecInt is an integer vector.
type VecInt struct {
	X, Y int
}

// Box is a 4-cornered box of float64 vertices.
type Box struct {
	MinX, MinY, MaxX, MaxY float64
}

// BoxInt is a 4-cornered box of integer vertices.
type BoxInt struct {
	MinX, MinY, MaxX, MaxY int
}

//								  //
// - - - 'NEW___' FUNCTIONS - - - //
//								  //

// NewVec does what you think it does
func NewVec(x, y float64) Vec {
	var v Vec = Vec{
		X: x,
		Y: y,
	}
	return v
}

// NewVecInt does what you think it does
func NewVecInt(x, y int) VecInt {
	var vi VecInt = VecInt{
		X: x,
		Y: y,
	}
	return vi
}

// NewBox does what you think it does
func NewBox(minX, minY, maxX, maxY float64) Box {
	var b Box = Box{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
	}
	return b
}

// NewBoxInt does what you think it does
func NewBoxInt(minX, minY, maxX, maxY int) BoxInt {
	var bi BoxInt = BoxInt{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
	}
	return bi
}

//							    //
// - - - MEMBER FUNCTIONS - - - //
//							    //

// Dot returns the dot product of two vectors
func (v Vec) Dot(by Vec) float64 {
	return v.X*by.X + v.Y*by.Y
}

// Add adds
func (v Vec) Add(addend Vec) Vec {
	return Vec{
		X: v.X + addend.X,
		Y: v.Y + addend.Y,
	}
}

// Add adds
func (vi VecInt) Add(addend VecInt) VecInt {
	return VecInt{
		X: vi.X + addend.X,
		Y: vi.Y + addend.Y,
	}
}

// AddXY adds
func (vi VecInt) AddXY(x, y int) VecInt {
	return VecInt{
		X: vi.X + x,
		Y: vi.Y + y,
	}
}

// Subtract subtracts
func (v Vec) Subtract(subtrahend Vec) Vec {
	return Vec{
		X: v.X - subtrahend.X,
		Y: v.Y - subtrahend.Y,
	}
}

// Scale scales a vector by some value. Scale by 0 to zero the vector
func (v Vec) Scale(by float64) Vec {
	return Vec{
		X: v.X * by,
		Y: v.Y * by,
	}
}

// ToInt converts a Vec to a VecInt
func (v Vec) ToInt() VecInt {
	return VecInt{
		X: int(v.X),
		Y: int(v.Y),
	}
}

// ToFloat64 converts a VecInt to a Vec
func (vi VecInt) ToFloat64() Vec {
	return Vec{
		X: float64(vi.X),
		Y: float64(vi.Y),
	}
}

// ByVec returns a Box scaled by a Vec
func (b Box) ByVec(v Vec) Box {
	return Box{
		MinX: b.MinX * v.X,
		MinY: b.MinY * v.Y,
		MaxX: b.MaxX * v.X,
		MaxY: b.MaxY * v.Y,
	}
}

// ToInt converts a Box to BoxInt
func (b Box) ToInt() BoxInt {
	return BoxInt{
		MinX: int(b.MinX),
		MinY: int(b.MinY),
		MaxX: int(b.MaxX),
		MaxY: int(b.MaxY),
	}
}

// Expand returns a COPY of the called BoxInt, but expanded. (Or shrunk, if negative numbers for VecInt)
func (b BoxInt) Expand(by VecInt) BoxInt {
	b.MinX -= by.X
	b.MinY -= by.Y
	b.MaxX += by.X
	b.MaxY += by.Y

	return b
}

// Zero Vector
var (
	ZV  = Vec{0, 0}
	ZVI = VecInt{0, 0}
)

//							  //
// - - - MATH FUNCTIONS - - - //
//							  //

// GetBounds calculates a BoxInt, given a position and a size.
func GetBounds(pos, size Vec) BoxInt {
	return BoxInt{
		MinX: int(math.Floor(pos.X)),
		MinY: int(math.Floor(pos.Y)),
		MaxX: int(math.Ceil(pos.X+size.X)) - 1,
		MaxY: int(math.Ceil(pos.Y+size.Y)) - 1,
	}
}

// MinInt returns the minimum of two integers a and b, as an integer.
// Because Go won't do it for you.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MinIntList returns the minimum of an arbitrarily long list of integers
func MinIntList(numbers ...int) int {
	min := numbers[0]

	for i := 1; i < len(numbers); i++ {
		MinInt(min, numbers[i])
	}

	return min
}

// IsWithinBounds checks whether the provited vector is within the provided BoxInt
func IsWithinBounds(vi VecInt, bi BoxInt) bool {
	return vi.X >= bi.MinX && vi.X < bi.MaxX && vi.Y >= bi.MinY && vi.Y < bi.MaxY
}

// AreInRange tells you whether the passed float64 is within the range specified.
// The range is min-inclusive, and max-non-inclusive.
func AreInRange(min, max float64, numbers ...float64) bool {
	for _, num := range numbers {
		if num < min || num >= max {
			return false
		}
	}
	return true
}

// AreInRangeInt tells you whether the passed integer is within the range specified.
// The range is min-inclusive, and max-non-inclusive.
func AreInRangeInt(min, max int, numbers ...int) bool {
	for _, num := range numbers {
		if num < min || num >= max {
			return false
		}
	}
	return true
}

// DistanceFormula applies the distance formula to two Vecs
func DistanceFormula(v1, v2 Vec) float64 {
	return math.Sqrt((v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y))
}

// DistanceFormulaInt applies the distance formula to two VecInts
func DistanceFormulaInt(v1, v2 VecInt) float64 {
	return math.Sqrt(float64((v1.X-v2.X)*(v1.X-v2.X)) + float64((v1.Y-v2.Y)*(v1.Y-v2.Y)))
}

// MinMax returns the num if within range, or the min or max if it exceeds the range
func MinMax(min, max, num float64) float64 {
	return math.Max(min, math.Min(max, num))
}
