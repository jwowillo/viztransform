package main

import (
	"errors"
	"flag"
	"fmt"
	"image/png"
	"os"

	"github.com/jwowillo/viztransform/parse"
	"github.com/jwowillo/viztransform/viz"
)

func main() {
	if len(os.Args) != 2 {
		fail(errArgs)
	}
	arg := os.Args[1]
	f, err := os.OpenFile(arg+".png", os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fail(err)
	}
	defer f.Close()
	t, err := parse.Transformation(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := png.Encode(f, viz.Transformation(t)); err != nil {
		fail(err)
	}
}

var errArgs = errors.New("must pass file name to save to")

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func usage() {

}
