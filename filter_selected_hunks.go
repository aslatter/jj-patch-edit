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
			for i := range f.hunks {
				if !f.hunks[i].selected {
					unslectedHunks = append(unslectedHunks, f.hunks[i])
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
