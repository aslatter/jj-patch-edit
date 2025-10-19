package main

import (
	"iter"
)

func filterFile(tokens iter.Seq[token], filter func(f []byte) bool) iter.Seq[token] {
	return func(yield func(token) bool) {
		var isInFilteredFile bool
		for t := range tokens {
			// if token is not a file we did not filter the file
			if t.kind != tokenKindFile {
				if !isInFilteredFile {
					if !yield(t) {
						return
					}
				}
				continue
			}

			// we must have a file at this point
			// we assume the first line is useful somehow to
			// the filter-function
			line := t.body[0]
			if !filter(line) {
				isInFilteredFile = true
				continue
			}
			isInFilteredFile = false
			if !yield(t) {
				return
			}
		}
	}
}
