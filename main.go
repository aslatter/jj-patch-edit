package main

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

func main() {
	//fmt.Printf("%+v\n", os.Args)

	var iterErr error
	for diffFile := range walkDirPair(os.Args[1], os.Args[2], &iterErr) {
		fmt.Printf("%+v\n", diffFile)
	}
	if iterErr != nil {
		fmt.Fprintln(os.Stderr, "error: ", iterErr)
	}

	fmt.Println("Press enter to continue")
	_, _ = fmt.Scanln()

	// until we know what we're doing, exit with failure
	os.Exit(1)
}

type diffFile struct {
	relPath string
	left    string
	right   string
}

// walkDirPair iterates over two directories in lock-step. A 'diffFile' is
// returned for each unique relative path between the two folders.
func walkDirPair(left string, right string, outErr *error) iter.Seq[diffFile] {
	return func(yield func(diffFile) bool) {
		var leftErr error
		var rightErr error

		leftNext, leftStop := iter.Pull(walkDir(left, &leftErr))
		defer leftStop()

		rightNext, rightStop := iter.Pull(walkDir(right, &rightErr))
		defer rightStop()

		var left walkResult
		var leftOkay bool

		var right walkResult
		var rightOkay bool

		for {
			switch {
			case left.relPath == right.relPath:
				{
					// advance both
					left, leftOkay = leftNext()
					right, rightOkay = rightNext()
				}
			case left.relPath > right.relPath || !leftOkay:
				right, rightOkay = rightNext()

			case left.relPath < right.relPath || !rightOkay:
				left, leftOkay = leftNext()
			}

			if leftErr != nil {
				outErr = &leftErr
				return
			}
			if rightErr != nil {
				outErr = &rightErr
				return
			}

			if !leftOkay && !rightOkay {
				return
			}

			if !leftOkay || right.relPath < left.relPath {
				if !yield(diffFile{
					relPath: right.relPath,
					left:    "",
					right:   right.absPath,
				}) {
					return
				}

			}
			if !rightOkay || right.relPath > left.relPath {
				if !yield(diffFile{
					relPath: left.relPath,
					left:    left.absPath,
					right:   "",
				}) {
					return
				}
			}

			if left.relPath == right.relPath {
				if !yield(diffFile{
					relPath: left.relPath,
					left:    left.absPath,
					right:   right.absPath,
				}) {
					return
				}
			}
		}
	}
}

type walkResult struct {
	relPath string
	absPath string
}

// walkDir iterates over files in a folder, recursively, returning
// both the relative-path and absolute path. walkDir only returns regular
// files.
func walkDir(dir string, outErr *error) iter.Seq[walkResult] {
	stopErr := errors.New("stop")

	return func(yield func(walkResult) bool) {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			// ??
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}

			if !yield(walkResult{
				relPath: relPath,
				absPath: path,
			}) {
				return stopErr
			}

			return nil
		})
		if err == stopErr {
			err = nil
		}
		if err != nil {
			*outErr = err
		}
	}
}
