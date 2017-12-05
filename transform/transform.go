package transform

import (
	"github.com/jwowillo/viztransform/geometry"
)

// Apply the Transformation to the geometry.Point by applying each
// line-reflection making up the Transformation in order.
func Apply(t Transformation, p geometry.Point) geometry.Point {
	for _, l := range t {
		p = apply(l, p)
	}
	return p
}

// apply a line-reflection to the geometry.Point as described in
// TypeLineReflection.
func apply(l geometry.Line, p geometry.Point) geometry.Point {
	i := geometry.MustPoint(geometry.Intersection(
		l,
		geometry.PerpendicularThroughPoint(l, p),
	))
	return geometry.Point{X: p.X + 2*(i.X-p.X), Y: p.Y + 2*(i.Y-p.Y)}
}

// Compose Transformations into a single Transformation which is the
// line-reflections in each Transformation appended together in order.
func Compose(ts ...Transformation) Transformation {
	var composed Transformation
	for _, t := range ts {
		composed = append(composed, t...)
	}
	return composed
}

// NoTransformation creates a Transformation with TypeNoTransformation that does
// nothing to geometry.Points.
func NoTransformation() Transformation {
	return Transformation{}
}

// LineReflection creates a Transformation with TypeLineReflection that reflects
// geometry.Points about geometry.Line l as described by TypeLineReflection.
func LineReflection(l geometry.Line) Transformation {
	return Transformation{l}
}

// Translation creates a Transformation with TypeTranslation that translates
// geometry.Points by the geometry.Vector v as described by TypeTranslation.
//
// Returns NoTransformation() if v is length 0.
func Translation(v geometry.Vector) Transformation {
	length := geometry.Length(v)
	if geometry.IsZero(length) {
		return NoTransformation()
	}
	v = geometry.MustVector(geometry.Scale(v, length/2))
	a, b := geometry.Point{X: 0, Y: 0}, geometry.Point{X: v.I, Y: v.J}
	l := geometry.MustLine(geometry.NewLineFromPoints(a, b))
	return Transformation{
		geometry.PerpendicularThroughPoint(l, a),
		geometry.PerpendicularThroughPoint(l, b),
	}
}

// Rotation creates a Transformation with TypeRotation that rotates
// geometry.Points by geometry.Number rads radians counter-clockwise around
// geometry.Point p.
//
// Returns NoTransformation() if rads is 0.
func Rotation(p geometry.Point, rads geometry.Number) Transformation {
	if geometry.IsZero(rads) {
		return NoTransformation()
	}
	a := geometry.MustLine(geometry.NewLineFromPoints(
		p,
		geometry.Point{X: p.X + 1, Y: p.Y},
	))
	b := geometry.Rotate(a, p, rads/2)
	return Transformation{a, b}
}

// GlideReflection creates a Transformation with TypeGlideReflection that
// is a Transformation with TypeLineReflection with ref used as the
// geometry.Line used to create it composed with a Transformation with
// TypeTranslation with the projection of geometry.Vector v onto the
// geometry.Line ref used to create it.
//
// Returns LineReflection(ref) if v is length 0.
func GlideReflection(ref geometry.Line, v geometry.Vector) Transformation {
	a := geometry.PerpendicularThroughPoint(ref, geometry.Point{X: 0, Y: 0})
	b := geometry.PerpendicularThroughPoint(
		ref,
		geometry.Point{X: v.I, Y: v.J},
	)
	return Compose(
		LineReflection(ref),
		Translation(geometry.ShortestVector(a, b)),
	)
}
