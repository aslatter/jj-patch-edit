package main

import "errors"

// splitHunk implements the "split" operation during hunk-selection.
// It breaks up a large change into smaller changes. Not every hunk
// can be split - an un-splittable hunk will return the same change
// in the returned slice of hunks.
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
			contextSpanLen++
		} else if isAdd {
			currentHunk.header.newCount++
			currentHunkHasChanges = true
			contextSpanLen = 0
		} else if isRemove {
			currentHunk.header.oldCount++
			currentHunkHasChanges = true
			contextSpanLen = 0
		}

		// we split on two contiguous lines of context, and put one line of context
		// in each split hunk.
		//
		// Ideally we would have more context in each hunk, but the patch-apply
		// process doesn't seem to like overlapping context between hunks - so
		//  we'd need to re-coalesce selected splits which sounds like more
		//  trouble than I want.
		if currentHunkHasChanges && contextSpanLen == 2 {
			// drop the current context-line from the current hunk
			currentHunk.changes = currentHunk.changes[:len(currentHunk.changes)-1]
			currentHunk.header.oldCount--
			currentHunk.header.newCount--

			// start a new hunk
			var newHunk hunk
			newHunk.header.newOffset = currentHunk.header.newOffset + currentHunk.header.newCount
			newHunk.header.oldOffset = currentHunk.header.oldOffset + currentHunk.header.oldCount

			// add the current context-line to the new hunk
			newHunk.changes = append(newHunk.changes, changeBytes)
			newHunk.header.newCount++
			newHunk.header.oldCount++

			// add finished splitHunk to output
			newHunks = append(newHunks, currentHunk)

			currentHunk = newHunk
			currentHunkHasChanges = false
			contextSpanLen = 1
		}
	}
	if currentHunkHasChanges {
		newHunks = append(newHunks, currentHunk)
	}
	return newHunks, nil
}
