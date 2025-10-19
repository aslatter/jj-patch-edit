# Patch edit for jujitsu vcs

I want to try out the `jj` source-control tool - it seems
to fit how I like to think about managing commits.

However the 'git add -p' workflow is seared deep into my
muscle-memory, and I haven't taken to the built-in TUI
interactive-diff-selector in jujitsu.

This project is an attempt to write a jj-diff-editor which
mimics the Git patch-editor.

It shells out to `git diff` and `git apply` to perform
the actual editing.

