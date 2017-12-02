package transform

import (
	"math"

	g "github.com/jwowillo/viztransform/geometry"
)

// IsSimplified ...
func IsSimplified(t Transformation) bool {
	if len(t) > 3 {
		return false
	} else if len(t) == 3 {
		a, b, c := t[0], t[1], t[2]
		return g.ArePerpendicular(a, b) && g.AreParallel(b, c) ||
			g.AreParallel(a, b) && g.ArePerpendicular(b, c)
	} else if len(t) == 2 {
		return !g.AreSameLine(t[0], t[1])
	}
	return true
}

// Simplify ...
func Simplify(t Transformation) Transformation {
	stack := reverse(t)
	for len(stack) > 3 {
		n := len(stack)
		a, b, c, d := stack[n-1], stack[n-2], stack[n-3], stack[n-4]
		stack = stack[:n-4]
		stack = append(stack, reverse(simplify4(a, b, c, d))...)
	}
	if len(stack) == 3 {
		return simplify3(stack[2], stack[1], stack[0])
	}
	if len(stack) == 2 {
		return simplify2(stack[1], stack[0])
	}
	return stack
}

func simplify4(a, b, c, d g.Line) Transformation {
	// ONLY WAY SIMPLIFY3 DOESNT DO ANYTHING ARE WITH 2 GLIDE REFLECTIONS.
	// WHITE BOARD ALL CASES OF THIS AND SEE IF HAPPENS OR HOW TO SIMPLIFY.
	ls := simplify3(a, b, c)
	if len(ls) == 3 {
		// Rotate d to be parallel to the parallel lines in the glide
		// reflection so they can cancel.
		//
		// Need to think out the order of things more, the line
		// reflection will always be 4, so we need to check if the
		// parallel lines are the second thing or not.
	} else if len(ls) == 2 {
		return simplify3(ls[0], ls[1], d)
	} else if len(ls) == 1 {
		return simplify2(ls[0], d)

	}
	return Transformation{d}
}

func simplify3(a, b, c g.Line) Transformation {
	ls := simplify2(a, b)
	if len(ls) == 1 {
		return simplify2(ls[0], c)
	}
	ls = simplify2(b, c)
	if len(ls) == 1 {
		return simplify2(ls[0], a)
	}
	if g.AreParallel(a, b) && g.AreParallel(b, c) {
		a, b, c = shiftBToC(a, b, c)
		return Compose(Transformation{a}, simplify2(b, c))
	}
	if g.AreParallel(a, b) {
		// c and b aren't parallel since a and b are parallel and this
		// would mean all g.Lines are parallel, which would have been
		// caught above.
		// This is done to reduce special case handling.
		a, b, c = c, b, a
	}
	a, b, c = rotateToGlide(a, b, c)
	return Compose(simplify2(a, b), Transformation{c})
}

// Only works for parallel g.Lines.
func shiftBToC(a, b, c g.Line) (g.Line, g.Line, g.Line) {
	l := g.MustLine(g.NewLineFromPoints(
		b.A(),
		g.MustPoint(g.Intersection(c, g.Perpendicular(c, b.A()))),
	))
	// TODO: Maybe this could be shift so DX and DY don't have to be public?
	// Everything happening here feels very primitive.
	return g.MustLine(g.NewLineFromPoints(
			g.Point{X: a.A().X + g.DX(l), Y: a.A().Y + g.DY(l)},
			g.Point{X: a.B().X + g.DX(l), Y: a.B().Y + g.DY(l)},
		)), g.MustLine(g.NewLineFromPoints(
			g.Point{X: b.A().X + g.DX(l), Y: b.A().Y + g.DY(l)},
			g.Point{X: b.B().X + g.DX(l), Y: b.B().Y + g.DY(l)},
		)), c
}

func rotateToGlide(a, b, c g.Line) (g.Line, g.Line, g.Line) {
	rads := g.Angle(b, c)
	i := g.MustPoint(g.Intersection(a, b))
	a, b, c = g.Rotate(a, i, math.Pi/2-rads), g.Rotate(b, i, math.Pi/2-rads), c
	rads = g.Angle(a, b)
	i = g.MustPoint(g.Intersection(b, c))
	return a, g.Rotate(b, i, rads), g.Rotate(c, i, rads)
}

// simplify2 simplifies a Transformation with 2 g.Lines.
func simplify2(a, b g.Line) Transformation {
	// A Transformation with 2 g.Lines can only be simplified if both g.Lines
	// are the same.
	if g.AreSameLine(a, b) {
		return NoTransformation()
	}
	return Transformation{a, b}
}

// reverse a list of g.Lines.
func reverse(ls []g.Line) []g.Line {
	out := make([]g.Line, len(ls))
	for i := 0; i < len(ls); i++ {
		out[i] = ls[len(ls)-i-1]
	}
	return out
}
