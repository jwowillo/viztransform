package transform

import (
	"fmt"

	"github.com/jwowillo/viztransform/geometry"
)

// TODO: Document that the output GlideReflection vector will always be in the
// direction of ref so the length can be trivially calculated.

// Types of Transformations.
//
// All Transformations fall into one of these categories after simplification.
// The Type of a non-simplified Transformation is the Type of its simplified
// form.
const (
	// TypeNoTransformation does nothing.
	//
	// A Transformation with no Lines has this Type.
	TypeNoTransformation Type = iota
	// TypeLineReflection reflects Points across a Line.
	//
	// This is done in such a way that the perpendicular to the Line ending
	// at the Point is found, then mirrored across the Line. The Point is
	// placed at the end of this mirrored perpendicular.
	//
	// A Transformation with a single Line has this Type.
	TypeLineReflection
	// TypeTranslation translates Points.
	//
	// A Transformation with two parallel Lines has this Type.
	TypeTranslation
	// TypeRotation rotates Points about another Point.
	//
	// A Transformation with two intersecting Lines has this Type.
	TypeRotation
	// TypeGlideReflection is a translation followed by a Line-reflection.
	//
	// A Transformation with two parallel Lines perpendicular to another
	// Line has this Type.
	TypeGlideReflection
)

// Transformation is a list of Lines defining Line-reflections meant to be
// performed in order.
//
// All Transformations can be simplified into at most 3 Lines.
type Transformation []geometry.Line

// String representation of the Transformation.
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

func stringNoTransformation() string {
	return "NoTransformation()"
}

func stringLineReflection(l geometry.Line) string {
	return fmt.Sprintf("LineReflection(%s)", l)
}

func stringTranslation(a, b geometry.Line) string {
	v := geometry.ShortestVector(a, b)
	v = geometry.MustVector(geometry.Scale(v, 2*geometry.Length(v)))
	return fmt.Sprintf("Translation(%s)", v)
}

func stringRotation(a, b geometry.Line) string {
	return fmt.Sprintf(
		"Rotation(%s, %s)",
		geometry.MustPoint(geometry.Intersection(a, b)),
		2*geometry.Angle(a, b),
	)
}

func stringGlideReflection(a, b, c geometry.Line) string {
	var parallelA, parallelB, not geometry.Line
	if geometry.AreParallel(a, b) {
		parallelA = a
		parallelB = b
		not = c
	} else {
		not = a
		parallelA = b
		parallelB = c
	}
	v := geometry.ShortestVector(parallelA, parallelB)
	return fmt.Sprintf(
		"GlideReflection(%s, %s)",
		not,
		geometry.MustVector(geometry.Scale(v, 2*geometry.Length(v))),
	)
}

// Type of a Transformation.
type Type int

// TypeOf a Transformation.
//
// Can be any defined value for Type. The Type of a non-simplified
// Transformation is the Type of its simplified form.
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
