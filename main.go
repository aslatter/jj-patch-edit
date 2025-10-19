package main

import (
	"bytes"
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
	tokens = filterFile(tokens, func(f []byte) bool {
		// todo - make better
		return !bytes.Contains(f, []byte("JJ-INSTRUCTIONS"))
	})
	tokens = promptUser(tokens, &promptError)
	tokens = invertDiff(tokens)

	applyErr := apply(os.Args[2], tokens)
	// applyErr := fakeApply(tokens)

	err := errors.Join(iterErr, promptError, applyErr)
	if err != nil {
		return err
	}

	return nil
}
