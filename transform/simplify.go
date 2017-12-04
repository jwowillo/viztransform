package transform

// Algorithm documented in report.

import (
	"math"

	g "github.com/jwowillo/viztransform/geometry"
)

// IsSimplified ...
func IsSimplified(t Transformation) bool {
	if len(t) < 2 {
		return true
	}
	if len(t) == 2 {
		return !g.AreSameLine(t[0], t[1])
	}
	if len(t) == 3 {
		a, b, c := t[0], t[1], t[2]
		return g.ArePerpendicular(a, b) && g.AreParallel(b, c) ||
			g.AreParallel(a, b) && g.ArePerpendicular(b, c)
	}
	return false
}

// Simplify ...
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

// simplify2 simplifies a Transformation with 2 g.Lines.
func simplify2(a, b g.Line) Transformation {
	if g.AreSameLine(a, b) {
		return Transformation{}
	}
	return Transformation{a, b}
}

func simplify3(a, b, c g.Line) Transformation {
	if len(simplify2(a, b)) == 0 {
		return Transformation{c}
	}
	if len(simplify2(b, c)) == 0 {
		return Transformation{a}
	}
	if g.AreParallel(a, b) && g.AreParallel(b, c) {
		return Transformation{shiftBToC(a, b, c)}
	}
	a, b, c = rotateToPerpendicularAndParallel(a, b, c)
	return Compose(simplify2(a, b), Transformation{c})
}

func simplify4(a, b, c, d g.Line) Transformation {
	f3 := simplify3(a, b, c)
	if len(f3) < 3 {
		return Simplify(Compose(f3, Transformation{d}))
	}
	l3 := simplify3(b, c, d)
	if len(l3) < 3 {
		return Simplify(Compose(Transformation{a}, l3))
	}
	a, b, c = f3[0], f3[1], f3[2]
	if g.AreParallel(b, d) {
		return Transformation{shiftBToC(a, b, d), c}
	}
	a, d = rotateBCToSame(a, c, b, d)
	return Transformation{a, d}
}

func rotateBCToSame(a, b, c, d g.Line) (g.Line, g.Line) {
	ia := g.MustPoint(g.Intersection(a, b))
	ib := g.MustPoint(g.Intersection(c, d))
	l := g.MustLine(g.NewLineFromPoints(ia, ib))
	radsa, radsb := g.Angle(b, l), g.Angle(c, l)
	a, d = g.Rotate(a, ia, radsa), g.Rotate(d, ib, radsb)
	return a, d
}

// Only works for parallel g.Lines.
func shiftBToC(a, b, c g.Line) g.Line {
	return g.Shift(a, g.ShortestVector(b, c))
}

// Returned a and b always parallel
func rotateToPerpendicularAndParallel(a, b, c g.Line) (g.Line, g.Line, g.Line) {
	if g.AreParallel(a, b) {
		a, b, c = c, a, b
	}
	rads := g.Angle(b, c)
	i := g.MustPoint(g.Intersection(a, b))
	a, b, c = g.Rotate(a, i, math.Pi/2-rads), g.Rotate(b, i, math.Pi/2-rads), c
	rads = g.Angle(b, a)
	i = g.MustPoint(g.Intersection(b, c))
	a, b, c = a, g.Rotate(b, i, rads), g.Rotate(c, i, rads)
	return a, b, c
}
