package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/jwowillo/viztransform/parse"
	"github.com/jwowillo/viztransform/transform"
)

func main() {
	if len(os.Args) != 2 {
		fail(errArgs)
	}
	arg := os.Args[1]
	p, err := parse.Point(arg)
	if err != nil {
		fail(err)
	}
	t, err := parse.Transformation(os.Stdin)
	if err != nil {
		fail(err)
	}
	fmt.Println(parse.UnparsePoint(transform.Apply(t, p)))
}

var errArgs = errors.New("must pass point to transform")

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	out := "Usage of viztransform_apply:\n"
	out += "\tviztransform_apply '(x y)'\n\n"
	out += "\tThe passed point will be transformed by a transformation\n"
	out += "\tread from stdin as a new-line separated list of individual\n"
	out += "\ttransformations to be composed and terminated by an EOF.\n\n"
	out += "\tTransformations:\n"
	out += "\t- NoTransformation()\n"
	out += "\t- LineReflection({(ax ay) (bx by)})\n"
	out += "\t- Translation({(ax ay) (bx by)}, dist)\n"
	out += "\t- Rotation((x y), rads)\n"
	out += "\t- GlideReflection({(ax ay) (bx by)}, dist)\n\n"
	out += "\tExample:\n"
	out += "\tviztransform_apply '(1 1)'\n"
	out += "\tNoTransformation()\n"
	out += "\tLineReflection({(0 0) (0 1)})\n"
	out += "\tTranslation({(0 0) (1 0)}, 1)\n"
	out += "\tRotation((0 0), -1.5707963)\n"
	out += "\tGlideReflection({(0 -0.5) (1 -0.5)}, 1)\n"
	out += "\tEOF\n"
	out += "\t(2.000000 -1.000000)\n"
	fmt.Fprintf(os.Stderr, out)
}
