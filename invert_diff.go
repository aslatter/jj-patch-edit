package main

import (
	"bytes"
	"iter"
)

// invertDiff inverts the effect of every selected patch-file and hunk.
//
// file-deletions become file-creation, and we swap '+' and '-' signs
// in the changes themselves.
//
// We have to swap the effect of every patch because we ask the user what
// changes to select by presenting a diff between left and right, but then
// we effect the change by editing *right* in-place (wheras the raw diffs
// would apply cleanly to left).
func invertDiff(files iter.Seq[*file]) iter.Seq[*file] {
	return func(yield func(*file) bool) {
		for f := range files {
			invertAddOrDelete(f)

			for i := range f.hunks {
				h := &f.hunks[i]

				h.header.oldOffset, h.header.newOffset = h.header.newOffset, h.header.oldOffset
				h.header.oldCount, h.header.newCount = h.header.newCount, h.header.oldCount

				// change adds to removes and vice-versa
				for i := range h.changes {
					if len(h.changes[i]) == 0 {
						continue
					}
					switch h.changes[i][0] {
					case '+':
						h.changes[i][0] = '-'
					case '-':
						h.changes[i][0] = '+'
					}
				}
			}

			if !yield(f) {
				return
			}
		}
	}
}

func invertAddOrDelete(f *file) {
	// a new/delete should have at least 5 lines
	if len(f.header) < 5 {
		return
	}

	header := f.header

	isDelete := bytes.HasPrefix(header[1], []byte("deleted "))
	isNew := bytes.HasPrefix(header[1], []byte("new "))

	if !isDelete && !isNew {
		return
	}

	// lets get to work!

	if isDelete {
		header[1] = bytes.Replace(header[1], []byte("deleted "), []byte("new "), 1)
	} else {
		header[1] = bytes.Replace(header[1], []byte("new "), []byte("deleted "), 1)
	}

	headerFields := bytes.Fields(header[0])
	if len(headerFields) < 2 {
		// ?!
		return
	}

	leftFile := bytes.Clone(headerFields[len(headerFields)-2])
	rightFile := bytes.Clone(headerFields[len(headerFields)-1])

	// re-build +++/--- lines
	if isDelete {
		header[3] = []byte("--- /dev/null")
		header[4] = append([]byte("+++ "), rightFile...)
	} else {
		header[3] = append([]byte("--- "), leftFile...)
		header[4] = []byte("+++ /dev/null")
	}

	f.header = header
}
