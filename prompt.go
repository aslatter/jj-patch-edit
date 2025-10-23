package main

import (
	"fmt"
	"iter"
	"strings"
)

func promptUser(files iter.Seq[*file]) iter.Seq[*file] {
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
			for i := range f.hunks {
				h := &f.hunks[i]

				printHunk(h)

			promptLoop:
				for {
					fmt.Print("\nInclude change? [y, n, a, d, q, ?] ")

					var includeStr string
					_, _ = fmt.Scanln(&includeStr)
					includeStr = strings.ToLower(includeStr)

					switch includeStr {
					case "y":
						h.selected = true
						break promptLoop

					case "n":
						break promptLoop

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
 q - do not include any remaining changes
 a - include this change and all remaining changes in this file
 d - do not include this change or any remaining changes in this file
 ? - print this help message`)

	fmt.Println()
}
