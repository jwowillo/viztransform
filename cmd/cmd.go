// Package cmd is the parent package of all viztransform commands.
package cmd

import (
	"flag"
	"fmt"
	"os"
)

// Fail with error err.
func Fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

// Init command with description u by setting a usage func.
func Init(u string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, u)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, transformations)
	}
	flag.Parse()
}

// transformations usage string.
const transformations = `
Transformations:
	- NoTransformation(): Does nothing.
	- LineReflection({(ax ay) (bx by)}): Reflects points across the line.
	- Translation(<i j>): Translates points by the vector.
	- Rotation((x y), rads): Rotates points counter-clockwise by radians
	  around the point.
	- GlideReflection({(ax ay) (bx by)}, <i j>): Reflects points across the
	  line and translates by the vector.
`
