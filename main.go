package main

import (
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

	var fake bool
	if os.Getenv("PATCH_FAKE_APPLY") != "" {
		fake = true
	}

	leftFolderName := args[0]
	rightFolderName := args[1]

	var iterErr error
	var parseErr error
	var promptError error

	lines := diffLines(leftFolderName, rightFolderName, &iterErr)

	tokens := tokenize(lines)
	files := parse(tokens, &parseErr)
	files = promptUser(files, &promptError)
	files = filterSelectedHunks(files)
	files = invertDiff(files)

	var applyErr error
	if !fake {
		applyErr = apply(rightFolderName, files)
	} else {
		applyErr = fakeApply(files)
	}

	err := errors.Join(iterErr, parseErr, promptError, applyErr)
	if err != nil {
		return err
	}

	return nil
}
