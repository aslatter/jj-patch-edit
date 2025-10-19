package main

import (
	"fmt"
	"iter"
	"strings"
)

func promptUser(changes iter.Seq[token], outErr *error) iter.Seq[token] {
	return func(yield func(token) bool) {
		var currentFile token
		var emittedCurrentFile bool
		var didQuit bool

	changesLoop:
		for t := range changes {
			if didQuit {
				if !yield(t) {
					return
				}
				continue
			}

			printDiff(t)

			if t.kind == tokenKindFile {
				currentFile = t
				emittedCurrentFile = false
				continue
			}

		promptLoop:
			for {
				fmt.Println("\nInclude change? ")

				var includeStr string
				_, _ = fmt.Scanln(&includeStr)
				includeStr = strings.ToLower(includeStr)

				// we invert the logic of what
				// we're asking the user because
				// we want 'right' to look like
				// the new commit (but everything
				// is already in right)

				switch includeStr {
				case "y", "yes":
					// skip commit
					fmt.Println()
					continue changesLoop
				case "n", "no":
					if !emittedCurrentFile {
						emittedCurrentFile = true
						if !yield(currentFile) {
							return
						}
					}
					if !yield(t) {
						return
					}
					fmt.Println()
					continue changesLoop
				case "q", "quit":
					didQuit = true
					if !emittedCurrentFile {
						emittedCurrentFile = true
						if !yield(currentFile) {
							return
						}
					}
					if !yield(t) {
						return
					}
					continue changesLoop
				case "a", "abort":
					*outErr = fmt.Errorf("aborted")
					return
				default:
					fmt.Println("unknown command")
					continue promptLoop
				}
			}
		}
	}
}

func invertDiff(tokens iter.Seq[token]) iter.Seq[token] {
	// we show the user the diff between left and right,
	// but we actually need to *remove* changes from right
	// when doing a split - so we need to invert the orientation
	// of the diffs we chose not to keep.
	//
	// It's a bit backwards from what I'm used to, but it works.
	return func(yield func(token) bool) {
		for t := range tokens {
			if t.kind != tokenKindHunk {
				if !yield(t) {
					return
				}
				continue
			}

			// change adds to removes and vice-versa
			for i := range t.body {
				if len(t.body[i]) == 0 {
					continue
				}
				switch t.body[i][0] {
				case '+':
					t.body[i][0] = '-'
				case '-':
					t.body[i][0] = '+'
				}
			}
			if !yield(t) {
				return
			}
		}
	}
}
