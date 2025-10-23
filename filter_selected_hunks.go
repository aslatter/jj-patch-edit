package main

import "iter"

func filterSelectedHunks(files iter.Seq[*file]) iter.Seq[*file] {
	// we presented the user with a diff of left vs right.
	// but jujitsu expects us to modify what's on right - so we
	// actually only keep the hunks the user *didn't* select.
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
