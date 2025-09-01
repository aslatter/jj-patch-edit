package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
	"os/exec"
	"path/filepath"
	"slices"
)

// diffLines produces the raw diff-output between two folders.
func diffLines(leftPath string, rightPath string, outErr *error) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {

		wd, leftPath, rightPath := getFolders(leftPath, rightPath)

		cmd := exec.Command("diff",
			"-N", // treat absent files as empty
			"-r", // recursively compare any subdirectories found
			"-u", // unified diff
			"--",
			leftPath,
			rightPath,
		)
		cmd.Dir = wd

		errReader, err := cmd.StderrPipe()
		if err != nil {
			*outErr = fmt.Errorf("setting up error-output pipe: %s", err)
			return
		}
		outReader, err := cmd.StdoutPipe()
		if err != nil {
			*outErr = fmt.Errorf("setting up output pipe: %s", err)
			return
		}

		err = cmd.Start()
		if err != nil {
			*outErr = fmt.Errorf("starting diff subprocess: %s", err)
			return
		}

		defer func() {
			err := cmd.Wait()

			// ignore non-zero exit-code from diff
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if exitErr.ExitCode() > 0 {
					err = nil
				}
			}

			if *outErr == nil && err != nil {
				*outErr = err
			}
		}()

		defer func() {
			if *outErr != nil {
				// TODO - custom buffer to avoid blowing out memory
				// if stderr is too big?
				errBytes, err := io.ReadAll(errReader)
				if err != nil {
					*outErr = fmt.Errorf("reading error-output from diff: %s", err)
					return
				}
				if len(errBytes) != 0 {
					*outErr = errors.New(string(errBytes))
				}
			}
		}()

		s := bufio.NewScanner(outReader)
		for s.Scan() {
			if !yield(slices.Clone(s.Bytes())) {
				return
			}
		}

		if err = s.Err(); err != nil {
			*outErr = err
		}
	}
}

func getFolders(left, right string) (wd string, newLeft string, newRight string) {
	var err error

	left = filepath.Clean(left)
	right = filepath.Clean(right)

	wd = filepath.Dir(left)
	newLeft, err = filepath.Rel(wd, left)
	if err != nil {
		return "", left, right
	}
	newRight, err = filepath.Rel(wd, right)
	if err != nil {
		return "", left, right
	}
	return wd, newLeft, newRight
}
