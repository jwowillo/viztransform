// Package parse has functions that turn string representations of primitives in
// package geometry into their respective primitives.
package parse

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/jwowillo/viztransform/geometry"
	"github.com/jwowillo/viztransform/transform"
)

var (
	// ErrBadTransformation ...
	ErrBadTransformation = errors.New("bad trans")
	// ErrBadLine ...
	ErrBadLine = errors.New("bad line")
	// ErrBadPoint ...
	ErrBadPoint = errors.New("bad point")
	// ErrBadVector ...
	ErrBadVector = errors.New("bad vector")
	// ErrBadNumber ...
	ErrBadNumber = errors.New("bad number")
)

// Transformation ...
func Transformation(r io.Reader) (transform.Transformation, error) {
	var t transform.Transformation
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		x := scanner.Text()
		ctor, args, err := constructor(x)
		if err != nil {
			return nil, err
		}
		var nt transform.Transformation
		switch ctor {
		case "NoTransformation":
			nt, err = parseNoTransformation(args)
		case "LineReflection":
			nt, err = parseLineReflection(args)
		case "Translation":
			nt, err = parseTranslation(args)
		case "Rotation":
			nt, err = parseRotation(args)
		case "GlideReflection":
			nt, err = parseGlideReflection(args)
		default:
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

func constructor(x string) (string, []string, error) {
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

func parseNoTransformation(xs []string) (transform.Transformation, error) {
	//
	if len(xs) != 0 {
		return nil, ErrBadTransformation
	}
	return transform.NoTransformation(), nil
}

func parseLineReflection(xs []string) (transform.Transformation, error) {
	// {(%f %f) (%f %f)}
	if len(xs) != 1 {
		return nil, ErrBadTransformation
	}
	l, err := Line(xs[0])
	if err != nil {
		return nil, ErrBadTransformation
	}
	return transform.LineReflection(l), nil
}

func parseTranslation(xs []string) (transform.Transformation, error) {
	// <%f %f>
	if len(xs) != 1 {
		return nil, ErrBadTransformation
	}
	v, err := Vector(xs[0])
	if err != nil {
		return nil, ErrBadTransformation
	}
	return transform.Translation(v), nil
}

func parseRotation(xs []string) (transform.Transformation, error) {
	if len(xs) != 2 {
		return nil, ErrBadTransformation
	}
	p, err := Point(xs[0])
	if err != nil {
		return nil, ErrBadTransformation
	}
	rads, err := Number(xs[1])
	if err != nil {
		return nil, ErrBadTransformation
	}
	return transform.Rotation(p, rads), nil
}

func parseGlideReflection(xs []string) (transform.Transformation, error) {
	if len(xs) != 2 {
		return nil, ErrBadTransformation
	}
	l, err := Line(xs[0])
	if err != nil {
		return nil, ErrBadTransformation
	}
	v, err := Vector(xs[1])
	if err != nil {
		return nil, ErrBadTransformation
	}
	return transform.GlideReflection(l, v), nil
}

// Line ...
func Line(x string) (geometry.Line, error) {
	// {(%f %f) (%f %f)}
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

// Vector ...
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

// Point ...
func Point(sx string) (geometry.Point, error) {
	// (%f %f)
	if sx[0] != '(' || sx[len(sx)-1] != ')' {
		return geometry.Point{}, ErrBadPoint
	}
	fs := strings.Split(sx[1:len(sx)-1], " ")
	if len(fs) != 2 {
		return geometry.Point{}, ErrBadPoint
	}
	x, err := Number(fs[0])
	if err != nil {
		return geometry.Point{}, ErrBadPoint
	}
	y, err := Number(fs[1])
	if err != nil {
		return geometry.Point{}, ErrBadPoint
	}
	return geometry.Point{X: x, Y: y}, nil
}

// Number ...
func Number(x string) (geometry.Number, error) {
	n, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return 0, ErrBadNumber
	}
	return geometry.Number(n), nil
}
