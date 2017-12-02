package transform

import g "github.com/jwowillo/viztransform/geometry"

// Apply the Transformation to the Point by reflecting the Point about each Line
// in order.
func Apply(t Transformation, p g.Point) g.Point {
	for _, l := range t {
		p = apply(l, p)
	}
	return p
}

// apply a Transformation with a single Line by reflecting the Point about the
// Line.
func apply(l g.Line, p g.Point) g.Point {
	// MustPoint can be used since a Line and its perpendicular must
	// intersect.
	i := g.MustPoint(g.Intersection(l, g.Perpendicular(l, p)))
	return g.Point{X: p.X + 2*(i.X-p.X), Y: p.Y + 2*(i.Y-p.Y)}
}

// Compose Transformations into a single Transformation by appending them in the
// order they were given.
func Compose(ts ...Transformation) Transformation {
	var composed Transformation
	for _, t := range ts {
		composed = append(composed, t...)
	}
	return composed
}

// NoTransformation creates a Transformation with TypeNoTransformation that does
// nothing.
func NoTransformation() Transformation {
	return Transformation{}
}

// LineReflection creates a Transformation with TypeLineReflection that reflects
// about l.
func LineReflection(l g.Line) Transformation {
	return Transformation{l}
}

// Translation creates a Transformation with TypeTranslation that translates
// dist in the direction dir.
//
// Order of Points in dir decides the translation-direction. Negative dist will
// translate in the opposite direction.
func Translation(v g.Vector) Transformation {
	length := g.Length(v)
	if g.IsZero(length) {
		return NoTransformation()
	}
	v = g.Scale(v, length/2)
	a, b := g.Point{X: 0, Y: 0}, g.Point{X: v.I, Y: v.J}
	l := g.MustLine(g.NewLineFromPoints(a, b))
	return Transformation{g.Perpendicular(l, a), g.Perpendicular(l, b)}
}

// Rotation creates a Transformation with TypeRotation that rotates rads radians
// about p.
//
// Negative rads will rotate in the opposite direction.
func Rotation(p g.Point, rads g.Number) Transformation {
	a := g.MustLine(g.NewLineFromPoints(p, g.Point{X: p.X + 1, Y: p.Y}))
	b := g.Rotate(a, p, rads/2)
	return Transformation{a, b}
}

// GlideReflection creates a Transformation with TypeGlideReflection that
// reflects about ref then translates dist distance in direction parallel to
// ref.
//
// Negative distances will translate in the opposite direction.
func GlideReflection(ref g.Line, v g.Vector) Transformation {
	// TODO: Is this right?
	a := g.PerpendicularThroughPoint(ref, g.Point{X: 0, Y: 0})
	b := g.PerpendicularThroughPoint(ref, g.Point{X: v.I, Y: v.J})
	return Compose(LineReflection{ref}, Translation(g.ShortestVector(a, b)))
}
