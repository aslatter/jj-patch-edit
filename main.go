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
	var promptError error

	lines := diffLines(os.Args[1], os.Args[2], &iterErr)

	tokens := tokenize(lines)
	for range promptUser(tokens, &promptError) {
		// noop
	}

	err := errors.Join(iterErr, promptError)
	if err != nil {
		return err
	}

	fmt.Println("Press enter to continue")
	_, _ = fmt.Scanln()

	// until we know what we're doing, exit with failure
	return errors.New("unimplemented")
}
