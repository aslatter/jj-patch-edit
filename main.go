package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

func mainErr() error {

	var iterErr error

	for line := range diffLines(os.Args[1], os.Args[2], &iterErr) {
		fmt.Println(string(line))
	}
	if iterErr != nil {
		return iterErr
	}

	fmt.Println("Press enter to continue")
	_, _ = fmt.Scanln()

	// until we know what we're doing, exit with failure
	return errors.New("unimplemented")
}

func diffLines(leftPath string, rightPath string, outErr *error) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {

		wd, leftPath, rightPath := getFolders(leftPath, rightPath)

		cmd := exec.Command("diff",
			"-N",
			"-r",
			"-u",
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
			if *outErr == nil && err != nil {
				*outErr = err
			}
		}()

		defer func() {
			if *outErr != nil {
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
