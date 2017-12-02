package viz

import (
	"image"
	"image/color"

	"github.com/jwowillo/viztransform/geometry"
	"github.com/jwowillo/viztransform/transform"
)

const (
	width  = 500
	height = 500
)

func newSetImage() setImage {
	// TODO: Fill the canvas with white.
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

// Transformation ...
func Transformation(t transform.Transformation) image.Image {
	img := newSetImage()
	if transform.IsSimplified(t) {
		switch transform.TypeOf(t) {
		case transform.TypeNoTransformation:
			// Something.
		case transform.TypeLineReflection:
			img = lineReflection(img, t[0])
		case transform.TypeTranslation:
			// return translation(t)
		case transform.TypeRotation:
			// return rotation(t)
		case transform.TypeGlideReflection:
			// return glideReflection(t)
		}
	} else {
	}
	return img
}

type setImage interface {
	image.Image
	Set(int, int, color.Color)
}

func lineReflection(img setImage, l geometry.Line) setImage {
	return nil
}

/*
// LineReflection ...
func LineReflection(w io.Writer, l transform.LineReflection) error {
	// Make sure scaling works so I don't have to think super hard about
	// where things go. Set canvas size in var somewhere.

	// Pick 4 points located conviently around the LineReflection's Line.
	// Plot all of them.
	// Plot the Line
	// Plot all their perpendiculars to the Line with a right angle sign on
	// the line.
	// Plot other Reflected Points, also with normal line.

	// Abstract a function that given a Line and a Point, draws the Point
	// and the normal to the Line.
	// If the distance from the Line is 0, skip.

	// Label the distances.

	canvas := image.NewRGBA(image.Rect(0, 0, 500, 500))
	fillWhite(canvas)
	drawLine(canvas, -200, -53, 100, 50)
	return png.Encode(w, canvas)
}

// Translation ...
func Translation(w io.Writer, t transform.Translation) error {
	return nil
}

// Rotation ...
func Rotation(w io.Writer, r transform.Rotation) error {
	return nil
}

// GlideReflection ...
func GlideReflection(w io.Writer, g transform.GlideReflection) error {
	return nil
}

// Composition ...
func Composition(w io.Writer, c transform.Composition) error {
	return nil
}

func fillWhite(canvas *image.RGBA) {
	for i := 0; i <= 500; i++ {
		for j := 0; j <= 500; j++ {
			canvas.Set(
				i, j,
				color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
			)
		}
	}
}

func drawLine(canvas *image.RGBA, x1, y1, x2, y2 float64) {
	y1, y2 = -1*y1, -1*y2
	x, y := x1, y1
	m := (y2 - y1) / (x2 - x1)
	e := 0.001
	for math.Abs(x-x2) > e && math.Abs(y-y2) > e {
		drawPoint(canvas, x, y)
		x++
		y += m
	}
}

func round(x float64) int {
	return int(math.Floor(x + .5))
}

func drawPoint(canvas *image.RGBA, x, y float64) {
	const n = 2
	ix, iy := round(x)+250, round(y)+250
	for i := ix - n; i <= ix+n; i++ {
		for j := iy - n; j <= iy+n; j++ {
			di, dj := i-ix, j-iy
			if di*di+dj*dj > n*n {
				continue
			}
			canvas.Set(i, j, color.RGBA{R: 0xFF, A: 0xFF})
		}
	}
}
*/
