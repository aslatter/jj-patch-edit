package main

import (
	"fmt"
)

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	green = "\033[32m"
	red   = "\033[31m"
	cyan  = "\033[36m"
)

// printFileHeader prints the header of a file in color.
func printFileHeader(f *file) {
	for _, ln := range f.header {
		fmt.Println(bold + string(ln) + reset)
	}
}

// printHunk prints a colorized patch-hunk.
func printHunk(h *hunk) {
	fmt.Printf("%s@@ -%d,%d +%d,%d @@%s%s\n",
		cyan,
		h.header.oldOffset,
		h.header.oldCount,
		h.header.newOffset,
		h.header.newCount,
		reset,
		string(h.header.bonusContent),
	)
	for _, ln := range h.changes {
		var prefix string
		var suffix string

		switch ln[0] {
		case '+':
			prefix = green
			suffix = reset
		case '-':
			prefix = red
			suffix = reset
		}
		fmt.Println(prefix + string(ln) + suffix)
	}
}
