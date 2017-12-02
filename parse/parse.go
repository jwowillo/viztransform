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
	ErrBadLine = errors.New("bad lin")
	// ErrBadPoint ...
	ErrBadPoint = errors.New("bad point")
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
	// {(%f %f) (%f %f)}, %f
	if len(xs) != 2 {
		return nil, ErrBadTransformation
	}
	l, err := Line(xs[0])
	if err != nil {
		return nil, ErrBadTransformation
	}
	dist, err := strconv.ParseFloat(xs[1], 64)
	if err != nil {
		return nil, ErrBadTransformation
	}
	return transform.Translation(l, dist), nil
}

func parseRotation(xs []string) (transform.Transformation, error) {
	// (%f %f), %f
	if len(xs) != 2 {
		return nil, ErrBadTransformation
	}
	p, err := Point(xs[0])
	if err != nil {
		return nil, ErrBadTransformation
	}
	rads, err := strconv.ParseFloat(xs[1], 64)
	if err != nil {
		return nil, ErrBadTransformation
	}
	return transform.Rotation(p, rads), nil
}

func parseGlideReflection(xs []string) (transform.Transformation, error) {
	// {(%f %f) (%f %f)}, %f
	if len(xs) != 2 {
		return nil, ErrBadTransformation
	}
	l, err := Line(xs[0])
	if err != nil {
		return nil, ErrBadTransformation
	}
	dist, err := strconv.ParseFloat(xs[1], 64)
	if err != nil {
		return nil, ErrBadTransformation
	}
	return transform.GlideReflection(l, dist), nil
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
	x, err := strconv.ParseFloat(fs[0], 64)
	if err != nil {
		return geometry.Point{}, ErrBadPoint
	}
	y, err := strconv.ParseFloat(fs[1], 64)
	if err != nil {
		return geometry.Point{}, ErrBadPoint
	}
	return geometry.Point{X: x, Y: y}, nil
}
