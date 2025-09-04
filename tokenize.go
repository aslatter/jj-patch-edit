package main

import "iter"

type token struct {
	kind int
	body [][]byte
}

const (
	tokenKindFile = iota
	tokenKindHunk
)

// tokenize breaks up diff-output into file-headers and patch-hunks.
func tokenize(lines iter.Seq[[]byte]) iter.Seq[token] {
	return func(yield func(token) bool) {
		var currentTok token
		var state int
		// states:
		//  0 - initial. expecting a file-header
		//  1 - in file-header. run until we encounter an '@', transition to '1'
		//  2 - emitting a hunk. run until we encounter a '@' or 'd'. transition to '1' or '2'
		for line := range lines {
			if state == 0 {
				state = 1
				currentTok.kind = tokenKindFile
				currentTok.body = append(currentTok.body, line)
				continue
			}
			if state == 1 {
				if len(line) > 0 && line[0] == '@' {
					// start new token
					if !yield(currentTok) {
						return
					}
					currentTok = token{}
					currentTok.kind = tokenKindHunk
					state = 2
					currentTok.body = append(currentTok.body, line)
					continue
				}
				currentTok.body = append(currentTok.body, line)
				continue
			}
			// state == 2
			if len(line) > 0 && line[0] == '@' {
				// start new token
				if !yield(currentTok) {
					return
				}
				currentTok = token{}
				currentTok.kind = tokenKindHunk
				currentTok.body = append(currentTok.body, line)
				continue
			}
			if len(line) > 0 && line[0] == 'd' {
				// start new token
				if !yield(currentTok) {
					return
				}
				currentTok = token{}
				currentTok.kind = tokenKindFile
				state = 1
				currentTok.body = append(currentTok.body, line)
				continue
			}
			currentTok.body = append(currentTok.body, line)
		}
		if state != 0 {
			yield(currentTok)
		}
	}
}
