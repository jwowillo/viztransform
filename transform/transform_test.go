package transform_test

import (
	"fmt"
	"testing"

	"github.com/jwowillo/viztransform/transform"
)

func TestTranslation(t *testing.T) {
	l := transform.MustLine(transform.NewLineFromPoints(
		transform.Point{X: 0, Y: 1},
		transform.Point{X: 1, Y: 0},
	))
	fmt.Println(transform.Translation(l, 1)[0])
	fmt.Println(transform.Translation(l, 1)[1])
}
