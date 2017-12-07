package geometry

import (
	"math"
)

// Shift the Line by the Vector.
func Shift(l Line, v Vector) Line {
	return MustLine(NewLineFromPoints(
		Point{X: l.a.X + v.I, Y: l.a.Y + v.J},
		Point{X: l.b.X + v.I, Y: l.b.Y + v.J},
	))
}

// ShortestVector returns the shortest Vector from Lines a to b.
//
// The Vector has length 0 if the Lines intersect.
func ShortestVector(a, b Line) Vector {
	if !AreParallel(a, b) {
		return Vector{I: 0, J: 0}
	}
	p := Perpendicular(a)
	m, n := MustPoint(Intersection(a, p)), MustPoint(Intersection(b, p))
	return Vector{I: n.X - m.X, J: n.Y - m.Y}
}

// Scale a Vector to a new one of the same direction but given length.
//
// Returns ErrNoVector if v is length 0 since the scaled Vector's direction
// can't be determined.
func Scale(v Vector, l Number) (Vector, error) {
	cl := Length(v)
	if IsZero(cl) {
		return Vector{}, ErrNoVector
	}
	return Vector{I: l * v.I / cl, J: l * v.J / cl}, nil
}

// Length of Vector v.
func Length(v Vector) Number {
	return Number(math.Hypot(float64(v.I), float64(v.J)))
}

// Distance between Points a and b.
func Distance(a, b Point) Number {
	return Length(Vector{I: b.X - a.X, J: b.Y - a.Y})
}

// Angle to rotate Line a counter-clockwise around the intersection with Line b
// so both Lines are parallel.
func Angle(a, b Line) Number {
	if AreParallel(a, b) {
		return 0
	}
	radsa := math.Atan2(float64(dy(a)), float64(dx(a)))
	radsb := math.Atan2(float64(dy(b)), float64(dx(b)))
	rads := Number(math.Mod(radsa-radsb, 2*math.Pi))
	i := MustPoint(Intersection(a, b))
	if AreParallel(Rotate(a, i, rads), b) {
		return rads
	}
	return Number(math.Pi - rads)
}

// Perpendicular Line to Line l through any Point on l.
func Perpendicular(l Line) Line {
	return PerpendicularThroughPoint(l, l.a)
}

// PerpendicularThroughPoint returns a perpendicular Line to Line l that passes
// through Point p.
func PerpendicularThroughPoint(l Line, p Point) Line {
	return MustLine(NewLineFromPointAndSlope(p, -dx(l), dy(l)))
}

// Intersection returns the Point where Lines a and b intersect.
// b.
//
// Returns an ErrParallel if the Lines are parallel since the intersection won't
// exist if the Lines aren't the same or occurs at infinitely many Points if the
// Lines are the same.
func Intersection(a, b Line) (Point, error) {
	if AreParallel(a, b) {
		return Point{}, ErrNoIntersection
	}
	m, n := det(a.a.X, a.a.Y, a.b.X, a.b.Y), det(b.a.X, b.a.Y, b.b.X, b.b.Y)
	xn := det(m, a.a.X-a.b.X, n, b.a.X-b.b.X)
	yn := det(m, a.a.Y-a.b.Y, n, b.a.Y-b.b.Y)
	d := det(a.a.X-a.b.X, a.a.Y-a.b.Y, b.a.X-b.b.X, b.a.Y-b.b.Y)
	return Point{X: xn / d, Y: yn / d}, nil
}

// RotateAroundOrigin rotates Line l counter-clockwise by rads radians around
// the origin.
func RotateAroundOrigin(l Line, rads Number) Line {
	return Rotate(l, Point{X: 0, Y: 0}, rads)
}

// Rotate Line l counter-clockise around Point p by rads radians.
func Rotate(l Line, p Point, rads Number) Line {
	ax, ay, bx, by := l.a.X-p.X, l.a.Y-p.Y, l.b.X-p.X, l.b.Y-p.Y
	cos := Number(math.Cos(float64(rads)))
	sin := Number(math.Sin(float64(rads)))
	a := Point{
		X: det(ax, ay, sin, cos) + p.X,
		Y: det(ay, -ax, sin, cos) + p.Y,
	}
	b := Point{
		X: det(bx, by, sin, cos) + p.X,
		Y: det(by, -bx, sin, cos) + p.Y,
	}
	return MustLine(NewLineFromPoints(a, b))
}

// dx returns the x-difference of the Points on the Line.
func dx(l Line) Number {
	return l.a.X - l.b.X
}

// dy returns the y-difference of the Points on the Line.
func dy(l Line) Number {
	return l.a.Y - l.b.Y
}

// det is determinant of 2x2 matrix with a and b in the first row and c and d in
// the second row.
func det(a, b, c, d Number) Number {
	return a*d - b*c
}

// dot is the dot-product of Vectors a and b.
func dot(a, b Vector) Number {
	return a.I*b.I + a.J*b.J
}

// standard Line-equation for the Line in terms of the coefficients of x and y
// and the value c the equation is equal to.
//
// Line-equation is ax + by = c.
func standard(l Line) (Number, Number, Number) {
	m, n := dy(l), -dx(l)
	return m, n, m*l.a.X + n*l.a.Y
}
