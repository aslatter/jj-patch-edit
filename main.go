package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

func mainErr() error {

	var iterErr error

	lines := diffLines(os.Args[1], os.Args[2], &iterErr)
	tokens := tokenize(lines)

	for t := range tokens {
		for _, line := range t.body {
			fmt.Println(string(line))
		}

		if t.kind == tokenKindFile {
			// TODO - track current file header? or current patch?
			continue
		}
		fmt.Println("\nInclude change? ")
		_, _ = fmt.Scanln()
		fmt.Println()
	}
	if iterErr != nil {
		return iterErr
	}

	fmt.Println("Press enter to continue")
	_, _ = fmt.Scanln()

	// until we know what we're doing, exit with failure
	return errors.New("unimplemented")
}
