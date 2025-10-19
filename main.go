package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func mainErr() error {
	// we don't have flags now, but we may
	// want to add some later - so reject
	// invocations with flags.
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		return fmt.Errorf("got %d args, expected 2", len(args))
	}

	leftFolderName := args[0]
	rightFolderName := args[1]

	var iterErr error
	var promptError error

	lines := diffLines(leftFolderName, rightFolderName, &iterErr)

	tokens := tokenize(lines)
	tokens = filterFile(tokens, func(f []byte) bool {
		// todo - make better
		return !bytes.Contains(f, []byte("JJ-INSTRUCTIONS"))
	})
	tokens = promptUser(tokens, &promptError)
	tokens = invertDiff(tokens)
	applyErr := apply(rightFolderName, tokens)

	err := errors.Join(iterErr, promptError, applyErr)
	if err != nil {
		return err
	}

	return nil
}
