package main

import "errors"

func splitHunk(h *hunk) ([]hunk, error) {
	var newHunks []hunk
	// find spans of at least two non-change-lines
	var currentHunk hunk

	currentHunk.header.newOffset = h.header.newOffset
	currentHunk.header.oldOffset = h.header.oldOffset
	currentHunk.header.bonusContent = h.header.bonusContent

	var currentHunkHasChanges bool
	var contextSpanLen int

	for ln := range h.changes {
		changeBytes := h.changes[ln]
		currentHunk.changes = append(currentHunk.changes, changeBytes)
		if len(changeBytes) == 0 {
			return nil, errors.New("found empty change-line")
		}
		isAdd := changeBytes[0] == '+'
		isRemove := changeBytes[0] == '-'
		if !isAdd && !isRemove {
			currentHunk.header.newCount++
			currentHunk.header.oldCount++
			if currentHunkHasChanges {
				contextSpanLen++
			}
		} else if isAdd {
			currentHunk.header.newCount++
			currentHunkHasChanges = true
			contextSpanLen = 0
		} else if isRemove {
			currentHunk.header.oldCount++
			currentHunkHasChanges = true
			contextSpanLen = 0
		}
		if contextSpanLen == 2 {
			var newHunk hunk
			newHunk.header.newOffset = currentHunk.header.newOffset + currentHunk.header.newCount - 2
			newHunk.header.oldOffset = currentHunk.header.oldOffset + currentHunk.header.oldCount - 2
			newHunk.header.newCount = 2
			newHunk.header.oldCount = 2

			newHunk.changes = append(newHunk.changes, currentHunk.changes[len(currentHunk.changes)-2])
			newHunk.changes = append(newHunk.changes, currentHunk.changes[len(currentHunk.changes)-1])

			newHunks = append(newHunks, currentHunk)
			currentHunk = newHunk
			currentHunkHasChanges = false
			contextSpanLen = 0
		}
	}
	if currentHunkHasChanges {
		newHunks = append(newHunks, currentHunk)
	}
	return newHunks, nil
}
