package geometry

// IsZero returns true if Number n is less than Epsilon difference from 0.
func IsZero(n Number) bool {
	return AreEqual(n, 0)
}

// AreEqual returns true if a and b are less than Epsilon difference of each
// other.
func AreEqual(a, b Number) bool {
	return a-b < Epsilon && b-a < Epsilon
}

// AreSamePoint returns true if Points a and b are the same.
func AreSamePoint(a, b Point) bool {
	return AreEqual(a.X, b.X) && AreEqual(a.Y, b.Y)
}

// AreParallel returns true if Lines a and b are parallel.
func AreParallel(a, b Line) bool {
	return IsZero(det(dx(a), dy(a), dx(b), dy(b)))
}

// ArePerpendicular returns true if Lines a and b are perpendicular.
func ArePerpendicular(a, b Line) bool {
	m, n := Vector{I: dx(a), J: dy(a)}, Vector{I: dx(b), J: dy(b)}
	return IsZero(dot(m, n))
}

// AreSameLine return true if Lines a and b are the same.
func AreSameLine(a, b Line) bool {
	if !AreParallel(a, b) {
		return false
	}
	m, n, c := standard(b)
	return AreEqual(m*a.a.X+n*a.a.Y, c)
}
