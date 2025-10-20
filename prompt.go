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
				fmt.Print("\nInclude change? [y, n, a, q] ")

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
