package main

import (
	"bytes"
	"fmt"
)

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	green = "\033[32m"
	red   = "\033[31m"
	cyan  = "\033[36m"
)

func printDiff(t token) {
	if t.kind == tokenKindFile {
		for _, ln := range t.body {
			fmt.Println(bold + string(ln) + reset)
		}
		return
	}

	for _, ln := range t.body {
		var prefix string
		var suffix string

		if len(ln) == 0 {
			fmt.Println()
			continue
		}

		if ln[0] == '@' {
			end := bytes.Index(ln[1:], []byte("@@"))
			if end != -1 {
				fmt.Println(cyan + string(ln[:end+3]) + reset + string(ln[end+3:]))
				continue
			}
		}

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
