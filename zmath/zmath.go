package zmath

import "math"

// Zero Vector
var (
	ZV  = Vec{0, 0}
	ZVI = VecInt{0, 0}
)

//				   //
// - - - VEC - - - //
//				   //

// Vec is a 2D float64 vector.
// Note that Vec's member functions will NOT modify the underlying Vec when called!
type Vec struct {
	X, Y float64
}

// V returns a new (x, y) float64 vector
func V(x, y float64) Vec {
	var v Vec = Vec{
		X: x,
		Y: y,
	}
	return v
}

// VC returns a Vec from a complex number
func VC(c complex128) Vec {
	return Vec{
		X: real(c),
		Y: imag(c),
	}
}

// Add adds
func (v Vec) Add(addend Vec) Vec {
	return Vec{
		X: v.X + addend.X,
		Y: v.Y + addend.Y,
	}
}

// AddXY adds two float64's to a Vec (x adds to X; y adds to Y)
func (v Vec) AddXY(x, y float64) Vec {
	return v.Add(V(x, y))
}

// Subtract subtracts
func (v Vec) Subtract(subtrahend Vec) Vec {
	return Vec{
		X: v.X - subtrahend.X,
		Y: v.Y - subtrahend.Y,
	}
}

// Multiply multiplies
func (v Vec) Multiply(multiplicand Vec) Vec {
	return Vec{
		X: v.X * multiplicand.X,
		Y: v.Y * multiplicand.Y,
	}
}

// Divide divides. No divide-by-zero-checking included; you gotta do that yourself!
func (v Vec) Divide(divisor Vec) Vec {
	return Vec{
		X: v.X / divisor.X,
		Y: v.Y / divisor.Y,
	}
}

// Scale scales a Vec by some vlaue.
func (v Vec) Scale(by float64) Vec {
	return Vec{
		X: v.X * by,
		Y: v.Y * by,
	}
}

// Min returns the minimum of a Vec's two components
func (v Vec) Min() float64 {
	return math.Min(v.X, v.Y)
}

// Max returns the maximum of a Vec's two components
func (v Vec) Max() float64 {
	return math.Max(v.X, v.Y)
}

// Dot returns the dot product of two vectors
func (v Vec) Dot(by Vec) float64 {
	return v.X*by.X + v.Y*by.Y
}

// Slope returns the slope of the line connecting two points
func (v Vec) Slope(pt Vec) float64 {
	return (pt.Y - v.Y) / (pt.X - v.X)
}

// GreaterThan returns whether both X and Y of the called Vec are greater than the X and Y of the passed Vec
func (v Vec) GreaterThan(pt Vec) bool {
	return v.X > pt.X && v.Y > pt.Y
}

// LessThan returns whether both X and Y of the called Vec are less than the X and Y of the passed Vec
func (v Vec) LessThan(pt Vec) bool {
	return v.X < pt.X && v.Y < pt.Y
}

// Complex converts a Vec to a complex number, with Y being the imaginary component
func (v Vec) Complex() complex128 {
	return complex(v.X, v.Y)
}

// VI converts a Vec to a VecInt
func (v Vec) VI() VecInt {
	return VecInt{
		X: int(v.X),
		Y: int(v.Y),
	}
}

//				      //
// - - - VECINT - - - //
//				      //

// VecInt is a 2D integer vector.
// Note that VecInt's member functions will NOT modify the underlying VecInt when called!
type VecInt struct {
	X, Y int
}

// VI returns a new (x, y) int vector
func VI(x, y int) VecInt {
	var vi VecInt = VecInt{
		X: x,
		Y: y,
	}
	return vi
}

// VIC returns a VecInt from a complex number
func VIC(c complex128) VecInt {
	return VecInt{
		X: int(real(c)),
		Y: int(imag(c)),
	}
}

// Add adds
func (vi VecInt) Add(addend VecInt) VecInt {
	return VecInt{
		X: vi.X + addend.X,
		Y: vi.Y + addend.Y,
	}
}

// AddXY adds two integers to a VecInt (x adds to X; y adds to Y)
func (vi VecInt) AddXY(x, y int) VecInt {
	return vi.Add(VI(x, y))
}

// Subtract subtracts
func (vi VecInt) Subtract(subtrahend VecInt) VecInt {
	return VecInt{
		X: vi.X - subtrahend.X,
		Y: vi.Y - subtrahend.Y,
	}
}

// Multiply multiplies
func (vi VecInt) Multiply(multiplicand VecInt) VecInt {
	return VecInt{
		X: vi.X * multiplicand.X,
		Y: vi.Y * multiplicand.Y,
	}
}

// Divide divides. No divide-by-zero-checking included; you gotta do that yourself!
func (vi VecInt) Divide(divisor VecInt) VecInt {
	return VecInt{
		X: vi.X / divisor.X,
		Y: vi.Y / divisor.Y,
	}
}

// Scale scales a VecInt by some value.
func (vi VecInt) Scale(by float64) VecInt {
	return vi.V().Scale(by).VI()
}

// Min returns the minimum of a VecInt's two components
func (vi VecInt) Min() int {
	return MinInt(vi.X, vi.Y)
}

// Max returns the maximum of a VecInt's two components
func (vi VecInt) Max() int {
	return MaxInt(vi.X, vi.Y)
}

// Slope returns the slope of the line connecting two points
func (vi VecInt) Slope(pt VecInt) float64 {
	if pt.X == vi.X {
		return math.Inf(1)
	}
	return float64(pt.Y-vi.Y) / float64(pt.X-vi.X)
}

// GreaterThan returns whether both X and Y of the called VecInt are greater than the X and Y of the passed VecInt
func (vi VecInt) GreaterThan(pt VecInt) bool {
	return vi.X > pt.X && vi.Y > pt.Y
}

// LessThan returns whether both X and Y of the called VecInt are less than the X and Y of the passed VecInt
func (vi VecInt) LessThan(pt VecInt) bool {
	return vi.X < pt.X && vi.Y < pt.Y
}

// Complex converts a VecInt to a complex number, with Y being the imaginary component
func (vi VecInt) Complex() complex128 {
	return complex(float64(vi.X), float64(vi.Y))
}

// V converts a VecInt to a Vec
func (vi VecInt) V() Vec {
	return Vec{
		X: float64(vi.X),
		Y: float64(vi.Y),
	}
}

//                  //
// - - - RECT - - - //
//                  //

// Rect us a good old rectangle with float64 vertices.
// Note that member functions WILL reference and modify the called Rect diRectly.
type Rect struct {
	Min, Max Vec
}

// R returns a new Rect. *ANY* two coordinates can be passed in in ANY order, and a correctly organized
// Rect will still be returned.
func R(min, max Vec) *Rect {
	return &Rect{
		Min: V(math.Min(min.X, max.X), math.Min(min.Y, max.Y)),
		Max: V(math.Max(min.X, max.X), math.Max(min.Y, max.Y)),
	}
}

// Dx returns the X width of a Rect
func (r *Rect) Dx() float64 { return r.Max.X - r.Min.X }

// Dy returns the Y width of a Rect
func (r *Rect) Dy() float64 { return r.Max.Y - r.Min.Y }

// Diag returns the length of the diagonal of the RectInt
func (r *Rect) Diag() float64 { return DistanceFormula(r.Min, r.Max) }

// Area returns the Rect's area
func (r *Rect) Area() float64 { return r.Dx() * r.Dy() }

// Contains returns whether the called Rectangle contains the passed point. It is min-inclusive and max-exclusive,
// so a Rect from (0, 0) to (1, 2) ONLY contains the points (0, 0) and (0, 1). A rect where min == max will contain
// no points and always return false.
func (r *Rect) Contains(point Vec) bool {
	return point.X < r.Max.X && point.X >= r.Min.X && point.Y < r.Max.Y && point.Y >= r.Min.Y
}

// Overlaps returns whether the called Rect and the passed Rect overlap at all in 2D space
func (r *Rect) Overlaps(w Rect) bool {
	return !(r.Min.X >= w.Max.X || r.Max.X <= w.Min.X || r.Min.Y >= w.Max.Y || r.Max.Y <= w.Min.Y)
}

// Scale scales the called Rect by the passed value, around (0, 0)
func (r *Rect) Scale(by float64) *Rect {
	r.Max = r.Max.Scale(by)
	r.Min = r.Min.Scale(by)
	return r
}

// Shift moves a Rect by the desired amount
func (r *Rect) Shift(by Vec) *Rect {
	r.Min = r.Min.Add(by)
	r.Max = r.Max.Add(by)
	return r
}

// Expand will expand the called rectangle by the desired amount. If either value in the argument is negative, the
// map will still be expanded, but out from it minimum coordinate(s).
func (r *Rect) Expand(by Vec) *Rect {
	if by.X > 0 {
		r.Max.X += by.X
	} else {
		r.Min.X += by.X
	}
	if by.Y > 0 {
		r.Max.Y += by.Y
	} else {
		r.Min.Y += by.Y
	}
	return r
}

// Shrink will shrink the called rectangle by the desired amount. If either value in the argument is negative, the
// map will still be shrunk, but in from it minimum coordinate(s). It is not possible to make a Rect's Min exceed
// its Max or its Max be less than its Min using this function, but it may be shrunk such that Min == Max.
func (r *Rect) Shrink(by Vec) *Rect {
	if by.X > 0 {
		r.Max.X = math.Max(r.Min.X, r.Max.X-by.X)
	} else {
		r.Min.X = math.Min(r.Max.X, r.Min.X-by.X)
	}
	if by.Y > 0 {
		r.Max.Y = math.Max(r.Min.Y, r.Max.Y-by.Y)
	} else {
		r.Min.Y = math.Min(r.Max.Y, r.Min.Y-by.Y)
	}
	return r
}

//				       //
// - - - RECTINT - - - //
//				       //

// RectInt is a good old rectangle with integer vertices.
// Note that member functions WILL reference and modify the called RectInt diRectly.
type RectInt struct {
	Min, Max VecInt
}

// RI returns a new Rect. *ANY* two coordinates can be passed in in ANY order, and a correctly organized
// RectInt will still be returned.
func RI(min, max VecInt) *RectInt {
	return &RectInt{
		Min: VI(MinInt(min.X, max.X), MinInt(min.Y, max.Y)),
		Max: VI(MaxInt(min.X, max.X), MaxInt(min.Y, max.Y)),
	}
}

// Dx returns the X width of a RectInt
func (ri *RectInt) Dx() int { return ri.Max.X - ri.Min.X }

// Dy retursn the Y width of a RectInt
func (ri *RectInt) Dy() int { return ri.Max.Y - ri.Min.Y }

// Diag returns the length of the diagonal of the RectInt
func (ri *RectInt) Diag() float64 { return DistanceFormulaInt(ri.Min, ri.Max) }

// Area returns the RectInt's area
func (ri *RectInt) Area() float64 { return float64(ri.Dx() * ri.Dy()) }

// Contains returns whether the called Rectangle contains the passed point. It is min-inclusive and max-exclusive,
// so a Rect from (0, 0) to (1, 2) ONLY contains the points (0, 0) and (0, 1). A rect where min == max will contain
// no points and always return false.
func (ri *RectInt) Contains(point VecInt) bool {
	return point.X < ri.Max.X && point.X >= ri.Min.X && point.Y < ri.Max.Y && point.Y >= ri.Min.Y
}

// Overlaps returns whether the called rect and the passed rect overlap at all in 2D space
func (ri *RectInt) Overlaps(w RectInt) bool {
	return !(ri.Min.X >= w.Max.X || ri.Max.X <= w.Min.X || ri.Min.Y >= w.Max.Y || ri.Max.Y <= w.Min.Y)
}

// Shift moves a Rect by the desired amount
func (ri *RectInt) Shift(by VecInt) *RectInt {
	ri.Min = ri.Min.Add(by)
	ri.Max = ri.Max.Add(by)
	return ri
}

// Expand will expand the called rectangle by the desired amount. If either value in the argument is negative, the
// map will still be expanded, but out from it minimum coordinate(s).
func (ri *RectInt) Expand(by VecInt) *RectInt {
	if by.X > 0 {
		ri.Max.X += by.X
	} else {
		ri.Min.X += by.X
	}
	if by.Y > 0 {
		ri.Max.Y += by.Y
	} else {
		ri.Min.Y += by.Y
	}
	return ri
}

// Shrink will shrink the called rectangle by the desired amount. If either value in the argument is negative, the
// map will still be shrunk, but in from it minimum coordinate(s). It is not possible to make a Rect's Min exceed
// its Max or its Max be less than its Min using this function, but it may be shrunk such that Min == Max.
func (ri *RectInt) Shrink(by VecInt) *RectInt {
	if by.X > 0 {
		ri.Max.X = MaxInt(ri.Min.X, ri.Max.X-by.X)
	} else {
		ri.Min.X = MinInt(ri.Max.X, ri.Min.X-by.X)
	}
	if by.Y > 0 {
		ri.Max.Y = MaxInt(ri.Min.Y, ri.Max.Y-by.Y)
	} else {
		ri.Min.Y = MinInt(ri.Max.Y, ri.Min.Y-by.Y)
	}
	return ri
}

//				     //
// - - - BOXES - - - //
//				     //

// Box is a 4-cornered box of float64 vertices. [OUTDATED AND BAD; USE RECT INSTEAD]
type Box struct {
	MinX, MinY, MaxX, MaxY float64
}

// BoxInt is a 4-cornered box of integer vertices. [OUTDATED AND BAD; USE RECT INSTEAD]
type BoxInt struct {
	MinX, MinY, MaxX, MaxY int
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
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt returns the maximum of two integers a and b, as an integer.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinIntList returns the minimum of an arbitrarily long list of integers
func MinIntList(numbers ...int) int {
	min := numbers[0]

	for i := 1; i < len(numbers); i++ {
		min = MinInt(min, numbers[i])
	}

	return min
}

// MaxIntList returns the maximum of an arbitrarily long list of integers
func MaxIntList(numbers ...int) int {
	max := numbers[0]

	for i := 1; i < len(numbers); i++ {
		max = MaxInt(max, numbers[i])
	}

	return max
}

// IsWithinBounds checks whether the provited vector is within the provided BoxInt
func IsWithinBounds(vi VecInt, bi BoxInt) bool {
	return vi.X >= bi.MinX && vi.X < bi.MaxX && vi.Y >= bi.MinY && vi.Y < bi.MaxY
}

// AreInRange tells you whether the passed float64 is within the range specified.
// The range is min-inclusive, and max-exclusive.
func AreInRange(min, max float64, numbers ...float64) bool {
	for _, num := range numbers {
		if num < min || num >= max {
			return false
		}
	}
	return true
}

// AreInRangeInt tells you whether the passed integer is within the range specified.
// The range is min-inclusive, and max-exclusive.
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

// MinMaxInt returns the num if within range, or the min or max if it exceeds the range
func MinMaxInt(min, max, num int) int {
	return MaxInt(min, MinInt(max, num))
}
