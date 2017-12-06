// Package main simplifies a transform.Transformation with more documentation
// from the help flag.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jwowillo/viztransform/cmd"
	"github.com/jwowillo/viztransform/parse"
)

// main simplifies the transform.Transformation read from STDIN.
func main() {
	if len(os.Args) != 1 {
		cmd.Fail(errArgs)
	}
	t, err := parse.Transformation(os.Stdin)
	if err != nil {
		cmd.Fail(err)
	}
	fmt.Println(t)
}

// errArgs is the error when any arguments are passed.
var errArgs = errors.New("must not pass any args")

// init the command.
func init() {
	cmd.Init(usage)
}

// usage to print.
const usage = `viztransform_simplify usage:

	viztransform_simplify

	The transformation read from STDIN as a newline-separated and
	EOF-terminated list of transformations to be composed will be simplified
	into a single transformation.`
