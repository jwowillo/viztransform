package transform

import (
	"math"

	"github.com/jwowillo/viztransform/geometry"
)

// IsSimplified returns true if the Transformation t uses the minimum number of
// line-reflections to achieve the result of the Transformation.
func IsSimplified(t Transformation) bool {
	if len(t) < 2 {
		return true
	}
	if len(t) == 2 {
		return !geometry.AreSameLine(t[0], t[1])
	}
	if len(t) == 3 {
		a, b, c := t[0], t[1], t[2]
		return geometry.ArePerpendicular(a, b) &&
			geometry.AreParallel(b, c) ||
			geometry.AreParallel(a, b) &&
				geometry.ArePerpendicular(b, c)
	}
	return false
}

// Simplify Transformation t into its simplest form that expresses the same
// Transformation.
//
// The algorithm is documented at
// https://github.com/jwowillo/viztransform/blob/master/doc/algorithm.pdf.
func Simplify(t Transformation) Transformation {
	if len(t) < 2 {
		return t
	}
	if len(t) == 2 {
		return simplify2(t[0], t[1])
	}
	if len(t) == 3 {
		return simplify3(t[0], t[1], t[2])
	}
	return Simplify(Compose(
		Simplify(t[:len(t)-4]),
		simplify4(t[len(t)-4], t[len(t)-3], t[len(t)-2], t[len(t)-1]),
	))
}

// simplify2 simplifies a Transformation represented by geometry.Lines a and b
// into its simplest form.
func simplify2(a, b geometry.Line) Transformation {
	if geometry.AreSameLine(a, b) {
		return Transformation{}
	}
	return Transformation{a, b}
}

// simplify3 simplifies a Transformation represented by geometry.Lines a, b, and
// c into its simplest form.
func simplify3(a, b, c geometry.Line) Transformation {
	if len(simplify2(a, b)) == 0 {
		return Transformation{c}
	}
	if len(simplify2(b, c)) == 0 {
		return Transformation{a}
	}
	if geometry.AreParallel(a, b) && geometry.AreParallel(b, c) {
		return Transformation{shiftBToC(a, b, c)}
	}
	a, b, c = rotateToParallelAndPerpendicular(a, b, c)
	return Compose(simplify2(a, b), Transformation{c})
}

// simplify4 simplifies a Transformation represented by geometry.Lines a, b, c,
// and d into its simplest form.
func simplify4(a, b, c, d geometry.Line) Transformation {
	f3 := simplify3(a, b, c)
	if len(f3) < 3 {
		return Simplify(Compose(f3, Transformation{d}))
	}
	l3 := simplify3(b, c, d)
	if len(l3) < 3 {
		return Simplify(Compose(Transformation{a}, l3))
	}
	a, b, c = f3[0], f3[1], f3[2]
	if geometry.AreParallel(b, d) {
		return Transformation{shiftBToC(a, b, d), c}
	}
	a, d = rotateBCToSame(a, c, b, d)
	return Transformation{a, d}
}

// rotateBCToSame takes geometry.Lines a, b, c, and d representing a rotation
// with a and b and a rotation with c and d and simplifies them to a single
// rotation by turning the rotations so b and c are the same and cancel.
func rotateBCToSame(a, b, c, d geometry.Line) (geometry.Line, geometry.Line) {
	ia := geometry.MustPoint(geometry.Intersection(a, b))
	ib := geometry.MustPoint(geometry.Intersection(c, d))
	l := geometry.MustLine(geometry.NewLineFromPoints(ia, ib))
	radsa, radsb := geometry.Angle(b, l), geometry.Angle(c, l)
	return geometry.Rotate(a, ia, radsa), geometry.Rotate(d, ib, radsb)
}

// shiftBToC takes geometry.Lines a, b, and c representing line-reflections and
// simplifies them to a single line-reflection by shifting a and b by the
// geometry.Vector that makes b the same as c causing b and c to cancel.
func shiftBToC(a, b, c geometry.Line) geometry.Line {
	return geometry.Shift(a, geometry.ShortestVector(b, c))
}

// rotateToParallelAndPerpendicular takes geometry.Lines a, b, and c
// representing line-reflections with at least one pair of geometry.Lines
// intersecting and rotates them so that the first two returned geometry.Lines
// are parallel and the second two are perpendicular.
func rotateToParallelAndPerpendicular(
	a, b, c geometry.Line,
) (geometry.Line, geometry.Line, geometry.Line) {
	if geometry.AreParallel(a, b) {
		a, b, c = c, a, b
	}
	rads := geometry.Angle(b, c)
	i := geometry.MustPoint(geometry.Intersection(a, b))
	a = geometry.Rotate(a, i, math.Pi/2-rads)
	b = geometry.Rotate(b, i, math.Pi/2-rads)
	rads = geometry.Angle(b, a)
	i = geometry.MustPoint(geometry.Intersection(b, c))
	return a, geometry.Rotate(b, i, rads), geometry.Rotate(c, i, rads)
}
