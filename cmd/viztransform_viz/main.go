// Package main vizualizes a transform.Transformation with more documentation
// from the help flag.
package main

import (
	"errors"
	"image/png"
	"os"

	"github.com/jwowillo/viztransform/cmd"
	"github.com/jwowillo/viztransform/parse"
	"github.com/jwowillo/viztransform/viz"
)

// main outputs a vizualization for the transform.Transfromation read from
// STDIN.
func main() {
	if len(os.Args) != 2 {
		cmd.Fail(errArgs)
	}
	f, err := os.OpenFile(os.Args[1]+".png", os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		cmd.Fail(err)
	}
	defer f.Close()
	t, err := parse.Transformation(os.Stdin)
	if err != nil {
		cmd.Fail(err)
	}
	if err := png.Encode(f, viz.Transformation(t)); err != nil {
		cmd.Fail(err)
	}
}

// errArgs is the error when not a single output-file is passed.
var errArgs = errors.New("must pass output-file")

// init the command.
func init() {
	cmd.Init(usage)
}

const usage = `viztransform_viz usage:

	viztransform_viz

	A vizualization of the transformation read from STDIN as a
	newline-separated and EOF-terimanted list of transformations to be
	composed. The vizualization will consist of 1 panel demonstrating the
	transformation if the transformation is already simplified and 2 panels
	demonstrating the transformation and the isomorphic simplified
	transformation otherwise.`
