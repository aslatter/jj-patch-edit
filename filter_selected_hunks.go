package main

import "iter"

// filterSelectedHunks interprets the 'selected' field on each hunk, to re-write
// the stream of files to match what should be applied.
//
// The logic is backwards, because jj expects us to re-shape 'right' to look
// like the new commit - so we have to 'select' the diffs that the user did not
// want.
func filterSelectedHunks(files iter.Seq[*file]) iter.Seq[*file] {
	return func(yield func(*file) bool) {
		for f := range files {
			var unslectedHunks []hunk
			var droppedAdds int
			var droppedRemoves int
			for i := range f.hunks {
				if !f.hunks[i].selected {
					h := f.hunks[i]
					h.header.newOffset += droppedAdds
					h.header.newOffset -= droppedRemoves
					unslectedHunks = append(unslectedHunks, h)
				} else {
					adds, removes := f.hunks[i].addRemoveCount()
					droppedAdds += adds
					droppedRemoves += removes
				}
			}
			if len(unslectedHunks) > 0 {
				f.hunks = unslectedHunks
				if !yield(f) {
					return
				}
			}
		}
	}
}

func (h *hunk) addRemoveCount() (adds int, removes int) {
	for _, ln := range h.changes {
		if len(ln) == 0 {
			// ?!
			continue
		}
		switch ln[0] {
		case '+':
			adds++
		case '-':
			removes--
		default:
		}
	}
	return
}
