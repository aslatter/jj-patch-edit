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

	changesLoop:
		for t := range changes {
			for _, line := range t.body {
				fmt.Println(string(line))
			}

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
				switch includeStr {
				case "y", "yes":
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
				case "n", "no":
					fmt.Println()
					continue changesLoop
				case "q", "quit":
					return
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
