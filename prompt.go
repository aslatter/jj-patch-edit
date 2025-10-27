package main

import (
	"fmt"
	"iter"
	"slices"
	"strings"
)

// promptUser asks the user what to do with every change in every patch-file,
// annotating the hunks which they user wishes to keep.
func promptUser(files iter.Seq[*file], outErr *error) iter.Seq[*file] {
	return func(yield func(*file) bool) {

		var didQuit bool

		for f := range files {
			if didQuit {
				if !yield(f) {
					return
				}
				continue
			}
			printFileHeader(f)

		hunkLoop:
			for i := 0; i < len(f.hunks); i++ {
				h := &f.hunks[i]

				printHunk(h)

			promptLoop:
				for {
					fmt.Print("\nInclude change? [y, n, s, a, d, q, ?] ")

					var includeStr string
					_, _ = fmt.Scanln(&includeStr)
					includeStr = strings.ToLower(includeStr)

					switch includeStr {
					case "y":
						h.selected = true
						break promptLoop

					case "n":
						break promptLoop

					case "s":
						newHunks, err := splitHunk(h)
						if err != nil {
							outErr = &err
							return
						}
						if len(newHunks) > 1 {
							// swap in new hunks for the old one
							f.hunks = slices.Replace(f.hunks, i, i+1, newHunks...)
							// we may have a new backing-slice, so our old pointer
							// may be bad.
							h = &f.hunks[i]
							// the old printout has been split up, so re-print the
							// current hunk
							printHunk(h)
							continue promptLoop
						}

					case "q", "quit":
						// quiting is hard, because we need to include
						// all proposed diffs in the output for later
						// processing.
						didQuit = true
						break hunkLoop

					case "a":
						for j := i; j < len(f.hunks); j++ {
							f.hunks[j].selected = true
						}
						break hunkLoop

					case "d":
						break hunkLoop

					case "?":
						printHelp()

					default:
						fmt.Println("unknown command")
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
 s - split changes
 q - do not include any remaining changes
 a - include this change and all remaining changes in this file
 d - do not include this change or any remaining changes in this file
 ? - print this help message`)

	fmt.Println()
}
