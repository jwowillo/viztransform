// Package transform defines a Transformation and provides constructors for,
// ways to find the types of, and ways to simplify Transformations.
package transform

import (
	"fmt"

	"github.com/jwowillo/viztransform/geometry"
)

// Types of Transformations.
//
// All Transformations fall into one of these categories after simplification.
// The Type of a non-simplified Transformation is the Type of its simplified
// form.
const (
	// TypeNoTransformation belongs to Transformations that do nothing.
	//
	// A Transformation with no geometry.Lines has this Type.
	TypeNoTransformation Type = iota
	// TypeLineReflection belongs to Transformations that reflects
	// geometry.Points across a geometry.Line by mirroring the
	// geometry.Point across the geometry.Line through the perpendicular
	// to the geometry.Line passing through the geometry.Point.
	//
	// A Transformation with a single geometry.Line has this Type where the
	// geometry.Line is the one being reflected across.
	TypeLineReflection
	// TypeTranslation belongs to Transformations that translate
	// geometry.Points by a geometry.Vector.
	//
	// A Transformation with 2 parallel geometry.Lines has this Type where
	// the geometry.Vector is the shortest one from the first geometry.Line
	// to the second.
	TypeTranslation
	// TypeRotation belongs to Transformations that rotate geometry.Points
	// by an angle around a geometry.Point.
	//
	// A Transformation with 2 intersecting geometry.Lines has this Type
	// where the geometry.Point is the intersection of the geometry.Lines
	// and the angle is the angle between the first geometry.Line and the
	// second.
	TypeRotation
	// TypeGlideReflection belongs to Transformations that both translate
	// and reflect geometry.Points in any order.
	//
	// A Transformation with 2 parallel geometry.Lines perpendicular to
	// another geometry.Line with the parallel geometry.Lines next to each
	// other has this Type where the parallel geometry.Lines define the
	// corresponding Transformation with TypeTranslation and the remaining
	// geometry.Line defines the corresponding Transformation with
	// TypeLineReflection.
	TypeGlideReflection
)

// Transformation is a list of geometry.Lines each representing an individual
// Transformation with TypeLineReflection that are all composed together.
//
// The arrangement of the geometry.Lines creates different Transformation-Types.
type Transformation []geometry.Line

// String-representation of the Transformation.
//
// Looks like a called Transformation-constructor with arguments in their
// respective string-representations. Examples are:
//
// 	NoTransformation()
// 	LineReflection(geometry.Line)
// 	Translation(geometry.Vector)
// 	Rotation(geometry.Point, geometry.Number)
// 	GlideReflection(geometry.Line, geometry.Vector)
func (t Transformation) String() string {
	t = Simplify(t)
	var out string
	switch TypeOf(t) {
	case TypeNoTransformation:
		out = stringNoTransformation()
	case TypeLineReflection:
		out = stringLineReflection(t[0])
	case TypeTranslation:
		out = stringTranslation(t[0], t[1])
	case TypeRotation:
		out = stringRotation(t[0], t[1])
	case TypeGlideReflection:
		out = stringGlideReflection(t[0], t[1], t[2])
	}
	return out
}

// stringNoTransformation returns the string-representation of a Transformation
// with TypeNoTransformation.
//
// Looks like 'NoTransformation()'.
func stringNoTransformation() string {
	return "NoTransformation()"
}

// stringLineReflection returns the string-representation of the Transformation
// with TypeLineReflection over geometry.Line l.
//
// Looks like 'LineReflection(geometry.Line)' where geometry.Line is the passed
// geometry.Line.
func stringLineReflection(l geometry.Line) string {
	return fmt.Sprintf("LineReflection(%s)", l)
}

// stringTranslation returns the string-representation of the Transformation
// with TypeTranslation created by parallel geometry.Lines a and b.
//
// Looks like 'Translation(geometry.Vector)' where geometry.Vector is the
// shortest geometry.Vector from a to b scaled by 2 since a translation from 2
// line-reflections translates by 2 times the shortest distance from the first
// geometry.Line to the second.
func stringTranslation(a, b geometry.Line) string {
	v := geometry.ShortestVector(a, b)
	v = geometry.MustVector(geometry.Scale(v, 2*geometry.Length(v)))
	return fmt.Sprintf("Translation(%s)", v)
}

// stringRotation returns the string-representation of the Transformation with
// TypeRotation created by intersecting geometry.Lines a and b.
//
// Looks like 'Rotation(geometry.Point, geometry.Angle)' where geometry.Point
// is the intersection of a and b and geometry.Angle is 2 times the angle from
// a to b mod 2pi since a rotation from 2 line-reflections rotates 2 times the
// angle from the first geometry.Line to the second.
func stringRotation(a, b geometry.Line) string {
	return fmt.Sprintf(
		"Rotation(%s, %s)",
		geometry.MustPoint(geometry.Intersection(a, b)),
		2*geometry.AngleBetween(a, b),
	)
}

// stringGlideReflection returns the string-representation of the Transformation
// with TypeGlieReflection created by geometry.Lines a, b, and c with a and be
// parallel and b and c perpendicular or a and b perpendicular or b and c
// parallel.
//
// Looks like 'GlideReflection(geometry.Line, geometry.Vector)' where
// geometry.Line is the geometry.Line the reflection is happening over and the
// projection of the geometry.Vector onto the geometry.Line is the
// geometry.Vector of translation.
func stringGlideReflection(a, b, c geometry.Line) string {
	if geometry.AreParallel(b, c) {
		a, b, c = b, c, a
	}
	v := geometry.ShortestVector(a, b)
	return fmt.Sprintf(
		"GlideReflection(%s, %s)",
		c,
		geometry.MustVector(geometry.Scale(v, 2*geometry.Length(v))),
	)
}

// Type of a Transformation in terms of how it transforms geometry.Points.
type Type int

// TypeOf a Transformation from the defined Transformation-Types.
func TypeOf(t Transformation) Type {
	t = Simplify(t)
	if len(t) == 0 {
		return TypeNoTransformation
	} else if len(t) == 1 {
		return TypeLineReflection
	} else if len(t) == 2 && geometry.AreParallel(t[0], t[1]) {
		return TypeTranslation
	} else if len(t) == 2 {
		return TypeRotation
	}
	return TypeGlideReflection
}
