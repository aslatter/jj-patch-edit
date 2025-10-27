package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
)

type file struct {
	header [][]byte
	hunks  []hunk
}

type hunk struct {
	header   hunkHeader
	changes  [][]byte
	selected bool
}

type hunkHeader struct {
	oldOffset    int
	oldCount     int
	newOffset    int
	newCount     int
	bonusContent []byte
}

// parse converts a stream of tokenized diff-sections (file-headers and hunks)
// into patch-files, with every hunk grouped under the appropriate file.
func parse(tokens iter.Seq[token], outErr *error) iter.Seq[*file] {
	return func(yield func(*file) bool) {
		var currentFile *file
		for t := range tokens {
			switch t.kind {

			case tokenKindFile:
				if currentFile != nil {
					if !yield(currentFile) {
						return
					}
				}
				currentFile = &file{
					header: t.body,
				}

			case tokenKindHunk:
				if currentFile == nil {
					err := errors.New("invalid diff: received change-hunk before file-header")
					outErr = &err
					return
				}

				if len(t.body) == 0 {
					err := errors.New("invalid diff: received change-hunk of zero length")
					outErr = &err
					return
				}

				headerBytes := t.body[0]

				var newHunk hunk
				header := &newHunk.header
				_, err := fmt.Sscanf(
					string(headerBytes),
					"@@ -%d,%d +%d,%d @@",
					&header.oldOffset,
					&header.oldCount,
					&header.newOffset,
					&header.newCount,
				)
				if err != nil {
					err = fmt.Errorf("parsing change-hunk header: %s", err)
					outErr = &err
					return
				}

				headerEndIndex := bytes.Index(headerBytes[1:], []byte("@@"))
				if headerEndIndex != -1 {
					header.bonusContent = headerBytes[1+headerEndIndex+2:]
				}

				newHunk.changes = t.body[1:]

				currentFile.hunks = append(currentFile.hunks, newHunk)
			}
		}
		if currentFile != nil {
			if !yield(currentFile) {
				return
			}
		}
	}
}

func printFiles(w io.Writer, fs iter.Seq[*file]) error {
	for f := range fs {
		for _, ln := range f.header {
			fmt.Fprintln(w, string(ln))
		}
		for _, h := range f.hunks {
			fmt.Fprintf(w, "@@ -%d,%d +%d,%d @@%s\n",
				h.header.oldOffset,
				h.header.oldCount,
				h.header.newOffset,
				h.header.newCount,
				string(h.header.bonusContent),
			)

			for _, ln := range h.changes {
				fmt.Fprintln(w, string(ln))
			}
		}
	}

	return nil
}
