// Package geometry defines geometric primitives like Number, Point, Vector, and
// Line.
//
// The package also defines operations and predicates on those primitives.
package geometry

import (
	"errors"
	"fmt"
	"strconv"
)

// Epsilon is the Number that two Numbers must have a difference with each other
// less than for them to be considered equal.
const Epsilon = Number(0.0000001)

var (
	// ErrNoIntersection is returned when parallel Lines are given to
	// something expecting intersecting Lines.
	ErrNoIntersection = errors.New("parallel Lines don't intersect")
	// ErrNoLine is returned when a Line is created with two Points that are
	// the same.
	ErrNoLine = errors.New("given information doesn't determine a Line")
	// ErrNoVector is returned when a length 0 Vector is given to something
	// expecting a Vector with a direction.
	ErrNoVector = errors.New("given information doesn't determine a Vector")
)

// MustLine panics if err isn't nil and returns the Line otherwise.
//
// Useful if a Line must exist based on preconditions and handling the error
// would be messy.
func MustLine(l Line, err error) Line {
	if err != nil {
		panic(err)
	}
	return l
}

// MustVector panics if err isn't nil and returns the Vector otherwise.
//
// Useful if a Vector must exist based on preconditions and handling the error
// would be messy.
func MustVector(v Vector, err error) Vector {
	if err != nil {
		panic(err)
	}
	return v
}

// MustPoint panics if the err isn't nil and returns the Line otherwise.
//
// Useful if a Point must exist based on preconditions and handling the error
// would be messy.
func MustPoint(p Point, err error) Point {
	if err != nil {
		panic(err)
	}
	return p
}

// Point in the xy-plane defined by x and y coordinates.
type Point struct{ X, Y Number }

// String-representation of the Point.
//
// Looks like '(X Y)' where X and Y are the Point's X and Y values.
func (p Point) String() string {
	return fmt.Sprintf("(%s %s)", p.X, p.Y)
}

// Line is a straight curve through 2 different Points.
type Line struct{ a, b Point }

// NewLineFromPointAndSlope creates a Line from a Point on the Line and the
// Line's slope.
//
// Returns ErrNoLine if rise and run are 0 since the 2 Points on the Line will
// be the same.
func NewLineFromPointAndSlope(p Point, rise, run Number) (Line, error) {
	return NewLineFromPoints(p, Point{X: p.X + run, Y: p.Y + rise})
}

// NewLineFromPoints creates a Line from 2 Points on the Line.
//
// Returns ErrNoLine if both Points are the same.
func NewLineFromPoints(a, b Point) (Line, error) {
	if AreSamePoint(a, b) {
		return Line{}, ErrNoLine
	}
	return Line{a: a, b: b}, nil
}

// String-representation of the Line.
//
// Looks like '{A B}' where A and B are two different Points on the Line.
func (l Line) String() string {
	return fmt.Sprintf("{%s %s}", l.a, l.b)
}

// Vector is a representation of a direction and a magnitude where the direction
// is the direction of an arrow from the origin to the Point with X value
// same as the Vector's I value and Y value same as the Vector's J value and
// and the magnitude is the length of the arrow.
type Vector struct{ I, J Number }

// String-representation of the Vector.
//
// Looks like '<I J>' where I and J are the Vector's I and J values.
func (v Vector) String() string {
	return fmt.Sprintf("<%s %s>", v.I, v.J)
}

// Number is an element of the real-numbers.
type Number float64

// String-representation of the Number.
//
// Looks like the Number truncated to 32 bits and with -0 turned to 0.
func (n Number) String() string {
	if AreEqual(n, 0) {
		n = 0
	}
	return strconv.FormatFloat(float64(n), 'f', -1, 32)
}
