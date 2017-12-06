// Package main applies a transform.Transformation to a geometry.Point with more
// documentation from the help flag.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jwowillo/viztransform/cmd"
	"github.com/jwowillo/viztransform/parse"
	"github.com/jwowillo/viztransform/transform"
)

// main applies the transform.Transformation read from STDIN to the
// geometry.Point in the args.
func main() {
	if len(os.Args) != 2 {
		cmd.Fail(errArgs)
	}
	arg := os.Args[1]
	p, err := parse.Point(arg)
	if err != nil {
		cmd.Fail(err)
	}
	t, err := parse.Transformation(os.Stdin)
	if err != nil {
		cmd.Fail(err)
	}
	fmt.Println(transform.Apply(t, p))
}

// errArgs is the error when not exactly a single geometry.Point is passed.
var errArgs = errors.New("must pass point to transform")

// init the command.
func init() {
	cmd.Init(usage)
}

// usage to print.
const usage = `viztransform_apply usage:

	viztransform_apply '(x y)'

	The passed point will be transformed by a transformation read from
	STDIN as a newline-separated and EOF-terminated list of transformations
	to be composed.`
