package parse

import (
	"fmt"

	"github.com/jwowillo/viztransform/transform"
)

// TODO: Same UnparseFloat stuf.

// UnparseSimplified ...
func UnparseSimplified(t transform.Transformation) string {
	t = transform.Simplify(t)
	var unparsed string
	switch transform.TypeOf(t) {
	case transform.TypeNoTransformation:
		unparsed = unparseNoTransformation()
	case transform.TypeLineReflection:
		unparsed = unparseLineReflection(t[0])
	case transform.TypeTranslation:
		unparsed = unparseTranslation(t[0], t[1])
	case transform.TypeRotation:
		unparsed = unparseRotation(t[0], t[1])
	case transform.TypeGlideReflection:
		unparsed = unparseGlideReflection(t[0], t[1], t[2])
	}
	return unparsed
}

func unparseNoTransformation() string {
	return "NoTransformation()"
}

func unparseLineReflection(l transform.Line) string {
	return fmt.Sprintf("LineReflection(%s)", l)
}

func unparseTranslation(a, b transform.Line) string {
	// TODO: Refactor the getting of this normal Line somehow. Do like one
	// in GlideReflection. Vector is key.
	n := transform.MustLine(transform.NewLineFromPoints(
		a.A(),
		transform.MustPoint(transform.Intersection(
			b,
			transform.Perpendicular(b, a.A())),
		),
	))
	dist := transform.Distance(
		transform.MustPoint(transform.Intersection(n, a)),
		transform.MustPoint(transform.Intersection(n, b)),
	)
	return fmt.Sprintf("Translation(%s, %s)", n, 2*dist)
}

func unparseRotation(a, b transform.Line) string {
	return fmt.Sprintf(
		"Rotation(%s, %s)",
		transform.MustPoint(transform.Intersection(a, b)),
		2*transform.Angle(a, b),
	)
}

func unparseGlideReflection(a, b, c transform.Line) string {
	var parallelA, parallelB, not transform.Line
	if transform.AreParallel(a, b) {
		parallelA = a
		parallelB = b
		not = c
	} else {
		not = a
		parallelA = a
		parallelB = b
	}
	n := transform.MustLine(transform.NewLineFromPoints(
		transform.MustPoint(transform.Intersection(
			parallelA,
			transform.Perpendicular(parallelA, not.A()),
		)),
		transform.MustPoint(transform.Intersection(
			parallelB,
			transform.Perpendicular(parallelB, not.A())),
		),
	))
	dist := transform.Distance(
		transform.MustPoint(transform.Intersection(n, parallelA)),
		transform.MustPoint(transform.Intersection(n, parallelB)),
	)
	return fmt.Sprintf("GlideReflection(%s, %s)", n, 2*dist)
}
