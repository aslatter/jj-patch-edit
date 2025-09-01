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

	var startedDiff bool
	var currentPatch [][]byte
	for line := range diffLines(os.Args[1], os.Args[2], &iterErr) {
		if startedDiff {
			if stoppingLine(line) {
				// prompt for inclusion of everything before this line
				fmt.Println("\nInclude change? ")
				_, _ = fmt.Scanln()
				fmt.Println()
				// apply current patch
				currentPatch = nil
				startedDiff = false
			}
		}
		if isStartOfDiff(line) {
			startedDiff = true
		}
		currentPatch = append(currentPatch, line)
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

func isStartOfDiff(line []byte) bool {
	if len(line) == 0 {
		return false
	}
	return line[0] == '@'
}

func stoppingLine(line []byte) bool {
	if len(line) == 0 {
		// ??
		return false
	}
	c := line[0]
	switch c {
	case '-', '+', ' ':
		return false
	}
	return true
}
