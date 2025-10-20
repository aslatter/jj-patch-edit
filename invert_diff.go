package main

import (
	"bytes"
	"fmt"
	"iter"
)

func invertDiff(tokens iter.Seq[token]) iter.Seq[token] {
	// we show the user the diff between left and right,
	// but we actually need to *remove* changes from right
	// when doing a split - so we need to invert the orientation
	// of the diffs we chose not to keep.
	//
	// It's a bit backwards from what I'm used to, but it works.
	return func(yield func(token) bool) {
		for t := range tokens {
			if t.kind == tokenKindFile {
				invertAddOrDelete(&t)

				if !yield(t) {
					return
				}
				continue
			}

			swapRangeBlock(&t)

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

func swapRangeBlock(t *token) {
	body := t.body
	if len(body) == 0 {
		// ?!
		return
	}
	rangeBlock := string(body[0])
	if len(rangeBlock) == 0 {
		// ?!
		return
	}

	var leftRange string
	var rightRange string
	n, err := fmt.Sscanf(rangeBlock, "@@ -%s +%s @@", &leftRange, &rightRange)
	if n != 2 || err != nil {
		// ?!
		return
	}

	var remainder []byte
	remainderIndex := bytes.Index(body[0][1:], []byte("@@"))
	if remainderIndex != -1 {
		remainder = body[0][remainderIndex+1:]
	}

	body[0] = append([]byte("@@ -"+rightRange+" +"+leftRange+" @@"), remainder...)
	t.body = body
}

func invertAddOrDelete(t *token) {
	// a new/delete should have at least 5 lines
	if len(t.body) < 5 {
		return
	}

	body := t.body

	isDelete := bytes.HasPrefix(body[1], []byte("deleted "))
	isNew := bytes.HasPrefix(body[1], []byte("new "))

	if !isDelete && !isNew {
		return
	}

	// lets get to work!

	if isDelete {
		body[1] = bytes.Replace(body[1], []byte("deleted "), []byte("new "), 1)
	} else {
		body[1] = bytes.Replace(body[1], []byte("new "), []byte("deleted "), 1)
	}

	headerFields := bytes.Fields(body[0])
	if len(headerFields) < 2 {
		// ?!
		return
	}

	leftFile := bytes.Clone(headerFields[len(headerFields)-2])
	rightFile := bytes.Clone(headerFields[len(headerFields)-1])

	// re-build +++/--- lines
	if isDelete {
		body[3] = []byte("--- /dev/null")
		body[4] = append([]byte("+++ "), rightFile...)
	} else {
		body[3] = append([]byte("--- "), leftFile...)
		body[4] = []byte("+++ /dev/null")
	}

	t.body = body
}
