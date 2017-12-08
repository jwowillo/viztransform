// Package parse has functions that turn string-representations of geometry
// package primitives into their respective primitives following the same
// pattern set by each primitive's string-representation.
package parse

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/jwowillo/viztransform/geometry"
	"github.com/jwowillo/viztransform/transform"
)

var (
	// ErrBadTransformation is returned when a transform.Transformation's
	// string is bad.
	ErrBadTransformation = errors.New("bad geometry.Transformation-string")
	// ErrBadLine is returned when a geometry.Line's string is bad.
	ErrBadLine = errors.New("bad geometry.Line-string")
	// ErrBadPoint is returned when a geometry.Point's string is bad.
	ErrBadPoint = errors.New("bad geometry.Point-string")
	// ErrBadVector is returned when a geometry.Vector's string is bad.
	ErrBadVector = errors.New("bad geometry.Vector-string")
	// ErrBadNumber is returned when a geometry.Number's string is bad.
	ErrBadNumber = errors.New("bad geometry.Number-string")
	// ErrBadAngle is returned when a geometry.Angle's string is bad.
	ErrBadAngle = errors.New("bad geometry.Angle-string")
)

// Transformation parses a transform.Transformation from the io.Reader r.
//
// A transform.Transformation's string is a newline separated string-list where
// string is a string-representation of a called
// transform.Transformation-constructor. Each string is turned into its
// respective transform.Transformations and then composed together.
//
// Returns an error if any string can't be parsed depending on the reason.
// Returns ErrBadTransformation if the constructor name isn't recognized, the
// calling syntax is bad, or the wrong number of arguments are passed to the
// constructor. Returns a corresponding ErrBad error if an argument can't be
// parsed to its geometry package primitive. All of these errors stem from parts
// of the string not fitting corresponding string-representation patterns.
func Transformation(r io.Reader) (transform.Transformation, error) {
	var t transform.Transformation
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		name, args, err := split(scanner.Text())
		if err != nil {
			return nil, err
		}
		var nt transform.Transformation
		switch name {
		case "NoTransformation":
			nt, err = noTransformation(args)
		case "LineReflection":
			nt, err = lineReflection(args)
		case "Translation":
			nt, err = translation(args)
		case "Rotation":
			nt, err = rotation(args)
		case "GlideReflection":
			nt, err = glideReflection(args)
		}
		if nt == nil {
			return nil, ErrBadTransformation
		}
		if err != nil {
			return nil, err
		}
		t = transform.Compose(t, nt)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return t, nil
}

// split x which is a transform.Transformation's string-representation into the
// constructor name name and arguments list.
//
// Returns ErrBadTransformation if x isn't formatted properly.
func split(x string) (string, []string, error) {
	i := strings.Index(x, "(")
	if i == -1 || x[len(x)-1] != ')' {
		return "", nil, ErrBadTransformation
	}
	args := strings.Split(x[i+1:len(x)-1], ", ")
	if len(args) == 1 && args[0] == "" {
		args = nil
	}
	return x[:i], args, nil
}

// noTransformation parses a transform.Transformation with
// transform.TypeNoTransformation from constructor arguments xs.
//
// Returns ErrBadTransformation if any arguments are passed.
func noTransformation(xs []string) (transform.Transformation, error) {
	if len(xs) != 0 {
		return nil, ErrBadTransformation
	}
	return transform.NoTransformation(), nil
}

// lineReflection parses a transform.Transformation with
// transform.TypeLineReflection from constructor arguments xs.
//
// Returns ErrBadTransformation if there isn't exactly 1 argument passed.
// Returns ErrBadLine if that argument can't be parsed to a geometry.Line.
func lineReflection(xs []string) (transform.Transformation, error) {
	if len(xs) != 1 {
		return nil, ErrBadTransformation
	}
	l, err := Line(xs[0])
	if err != nil {
		return nil, err
	}
	return transform.LineReflection(l), nil
}

// translation parses a transform.Transformation with transform.TypeTranslation
// from constructor arguments xs.
//
// Returns ErrBadTransformation if there isn't exactly 1 argument passed.
// Returns ErrBadVector if that argument can't be parsed to a geometry.Vector.
func translation(xs []string) (transform.Transformation, error) {
	if len(xs) != 1 {
		return nil, ErrBadTransformation
	}
	v, err := Vector(xs[0])
	if err != nil {
		return nil, err
	}
	return transform.Translation(v), nil
}

// rotation parses a transform.Transformation with transform.TypeRotation from
// constructor arguments xs.
//
// Returns ErrBadTransformation if there aren't exactly 2 arguments passed.
// Returns ErrBadPoint if the first argument can't be parsed to a
// geometry.Point. Returns ErrBadAngle if the second argument can't be parsed
// to a geometry.Angle.
func rotation(xs []string) (transform.Transformation, error) {
	if len(xs) != 2 {
		return nil, ErrBadTransformation
	}
	p, err := Point(xs[0])
	if err != nil {
		return nil, err
	}
	rads, err := Angle(xs[1])
	if err != nil {
		return nil, err
	}
	return transform.Rotation(p, rads), nil
}

// glideReflectin parses a transform.Transformation with
// transform.TypeGlideReflection from constructor arguments xs.
//
// Returns ErrBadTransformation if there aren't exactly 2 arguments passed.
// Returns ErrBadLine if the first argument can't be parsed to a geometry.Line.
// Returns ErrBadVector if the second argument can't be parsed to a
// geometry.Vector.
func glideReflection(xs []string) (transform.Transformation, error) {
	if len(xs) != 2 {
		return nil, ErrBadTransformation
	}
	l, err := Line(xs[0])
	if err != nil {
		return nil, err
	}
	v, err := Vector(xs[1])
	if err != nil {
		return nil, err
	}
	return transform.GlideReflection(l, v), nil
}

// Line parses a geometry.Line from the string x.
//
// Returns ErrBadLine if the string doesn't fit the geometry.Line
// string-representation pattern.
func Line(x string) (geometry.Line, error) {
	if x[0] != '{' || x[len(x)-1] != '}' {
		return geometry.Line{}, ErrBadLine
	}
	x = x[1 : len(x)-1]
	i := strings.Index(x, ")")
	if i == -1 {
		return geometry.Line{}, ErrBadLine
	}
	if x[i+1] != ' ' {
		return geometry.Line{}, ErrBadLine
	}
	sa, sb := x[:i+1], x[i+2:]
	a, err := Point(sa)
	if err != nil {
		return geometry.Line{}, ErrBadLine
	}
	b, err := Point(sb)
	if err != nil {
		return geometry.Line{}, ErrBadLine
	}
	return geometry.NewLineFromPoints(a, b)
}

// Vector parses a geometry.Vector from the string x.
//
// Returns ErrBadVector if the string doesn't fit the geometry.Vector
// string-representation pattern.
func Vector(sx string) (geometry.Vector, error) {
	if sx[0] != '<' || sx[len(sx)-1] != '>' {
		return geometry.Vector{}, ErrBadVector
	}
	fs := strings.Split(sx[1:len(sx)-1], " ")
	if len(fs) != 2 {
		return geometry.Vector{}, ErrBadVector
	}
	i, err := Number(fs[0])
	if err != nil {
		return geometry.Vector{}, ErrBadVector
	}
	j, err := Number(fs[1])
	if err != nil {
		return geometry.Vector{}, ErrBadVector
	}
	return geometry.Vector{I: i, J: j}, nil
}

// Point parses a geometry.Point from the string x.
//
// Returns ErrBadPoint if the string doesn't fit the geometry.Point
// string-representation pattern.
func Point(x string) (geometry.Point, error) {
	if x[0] != '(' || x[len(x)-1] != ')' {
		return geometry.Point{}, ErrBadPoint
	}
	fs := strings.Split(x[1:len(x)-1], " ")
	if len(fs) != 2 {
		return geometry.Point{}, ErrBadPoint
	}
	nx, err := Number(fs[0])
	if err != nil {
		return geometry.Point{}, ErrBadPoint
	}
	ny, err := Number(fs[1])
	if err != nil {
		return geometry.Point{}, ErrBadPoint
	}
	return geometry.Point{X: nx, Y: ny}, nil
}

// Number parses a geometry.Number from the string x.
//
// Returns ErrBadNumber if the string doesn't fit the geometry.Number
// string-representation pattern.
func Number(x string) (geometry.Number, error) {
	n, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return 0, ErrBadNumber
	}
	return geometry.Number(n), nil
}

// Angle parses a geometry.Angle from the string x.
//
// Returns ErrBadAngle if the string doesn't fit the geometry.Angle
// string-representation pattern.
func Angle(x string) (geometry.Angle, error) {
	a, err := Number(x)
	if err != nil {
		return 0, ErrBadAngle
	}
	return geometry.Angle(a), nil
}
