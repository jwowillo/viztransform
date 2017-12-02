package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/jwowillo/viztransform/parse"
)

func main() {
	if len(os.Args) != 1 {
		fail(errArgs)
	}
	t, err := parse.Transformation(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println(parse.UnparseSimplified(t))
}

var errArgs = errors.New("must not pass any args")

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
