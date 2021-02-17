package zmath

import (
	"math"
	"strconv"
)

// Quat represents a quaternion, a 4D complex number.
// Note that Quat's member functions will NOT modify the underlying Quat when called!
type Quat struct {
	A, I, J, K float64
}

// ZQ is a zero quaternion
var ZQ = Quat{0, 0, 0, 0}

// Q returns a new quaternion
func Q(a, i, j, k float64) Quat {
	return Quat{
		A: a,
		I: i,
		J: j,
		K: k,
	}
}

// Add adds two quaternions
func (q Quat) Add(v Quat) Quat {
	q.A += v.A
	q.I += v.I
	q.J += v.J
	q.K += v.K
	return q
}

// Sub subtracts two quaternions
func (q Quat) Sub(v Quat) Quat {
	q.A -= v.A
	q.I -= v.I
	q.J -= v.J
	q.K -= v.K
	return q
}

// Mult returns the result of multiplication of the called quaternion by the passed quaternion
func (q Quat) Mult(v Quat) Quat {
	return Quat{
		A: q.A*v.A - q.I*v.I - q.J*v.J - q.K*v.K,
		I: q.A*v.I + q.I*v.A + q.J*v.K - q.K*v.J,
		J: q.A*v.J - q.I*v.K + q.J*v.A + q.K*v.I,
		K: q.A*v.K + q.I*v.J - q.J*v.I + q.K*v.A,
	}
	// wow, that was tedious!
}

// Div returns the result of division of the called quaternion by the passed quaternion
func (q Quat) Div(v Quat) Quat {
	return q.Mult(v.Conj())
}

// PowInt returns the result of raising a quaternion to some integer power
func (q Quat) PowInt(exp int) Quat {
	if exp == 0 {
		return Quat{1, 0, 0, 0}
	}

	ret := q
	for p := 1; p < exp; p++ {
		ret = ret.Mult(q)
	}

	return ret
}

// Conj returns the conjugate of a quaternion
func (q Quat) Conj() Quat {
	q.I *= -1
	q.J *= -1
	q.K *= -1
	return q
}

// Abs returns the absolute value of a quat
func (q Quat) Abs() float64 {
	return math.Sqrt(q.A*q.A + q.I*q.I + q.J*q.J + q.K*q.K)
}

// String returns a string-formatted quaternion as a + bi + cj + dk to 3 decimal places per component
func (q Quat) String() string {
	return strconv.FormatFloat(q.A, 'E', 3, 64) + " + " +
		strconv.FormatFloat(q.I, 'E', 3, 64) + "i + " +
		strconv.FormatFloat(q.J, 'E', 3, 64) + "j + " +
		strconv.FormatFloat(q.K, 'E', 3, 64) + "k"
}
