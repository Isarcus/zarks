package zscript

import "github.com/Isarcus/zarks/zmath"

// Object is the interface type of all zscript variables, somewhat a la Python
type Object interface {
	Type() DataType
}

// native Go types
type (
	// INT represents an int
	INT int
	// FLOAT represents a float64
	FLOAT float64
	// BOOL represents a boolean
	BOOL bool
)

// zmath types
type (
	// SET represents a zmath.Set, for linear data
	SET zmath.Set
	// MAP represents a zmath.Map, for 2D data
	MAP zmath.Map
	// VEC represents a zmath.Vec
	VEC zmath.Vec
	// VECINT represents a zmath.VecInt
	VECINT zmath.VecInt
)

// Type returns Int
func (i INT) Type() DataType {
	return Int
}

// Type returns Float
func (f FLOAT) Type() DataType {
	return Float
}

// Type returns Bool
func (b BOOL) Type() DataType {
	return Bool
}

// Type returns Set
func (s SET) Type() DataType {
	return Set
}

// Type returns Map
func (m MAP) Type() DataType {
	return Map
}

// Type returns IVec
func (v VEC) Type() DataType {
	return Vec
}

// Type returns VecInt
func (v VECINT) Type() DataType {
	return VecInt
}
