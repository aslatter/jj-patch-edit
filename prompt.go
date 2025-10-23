package main

import (
	"fmt"
	"iter"
	"strings"
)

func promptUser(files iter.Seq[*file], outErr *error) iter.Seq[*file] {
	return func(yield func(*file) bool) {

	fileLoop:
		for f := range files {
			printFileHeader(f)

		hunkLoop:
			for i := range f.hunks {
				h := &f.hunks[i]

				printHunk(h)

			promptLoop:
				for {
					fmt.Print("\nInclude change? [y, n, a, q, ?] ")

					var includeStr string
					_, _ = fmt.Scanln(&includeStr)
					includeStr = strings.ToLower(includeStr)

					switch includeStr {
					case "y", "yes":
						h.selected = true
						continue hunkLoop

					case "n", "no":
						continue hunkLoop

					case "q", "quit":
						continue fileLoop

					case "a", "abort":
						*outErr = fmt.Errorf("aborted")
						return

					case "?":
						printHelp()
						continue promptLoop

					default:
						fmt.Println("unknown command")
						continue promptLoop
					}
				}
			}

			if !yield(f) {
				return
			}
		}
	}
}

func printHelp() {
	fmt.Println(`Select if this change should be included.

 y - include this change
 n - do not include this change
 q - do not include any remaining changes
 a - abort this operation as a whole
 ? - print this help message`)

	fmt.Println()
}
