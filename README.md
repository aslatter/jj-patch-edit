# Patch edit for jujitsu vcs

I want to try out the `jj` source-control tool - it seems
to fit how I like to think about managing commits.

However the 'git add -p' workflow is seared deep into my
muscle-memory, and I haven't taken to the built-in TUI
interactive-diff-selector in jujitsu.

This project is an attempt to write a jj diff-editing
tool which mimics the Git patch-editor.

It shells out to `git diff` and `git apply` to perform
the actual editing.

## Building

If you have "just" and Go installed run:

```
just build
```

This will produce a binary at `./bin/jj-patch-edit`.

Otherwise run:

```
go build -o bin/jj-patch-edit ./
```

## Using

Run an interactive `jj` command specifying a tool:

```
jj split -i --tool [path/to/]jj-patch-edit
```

You can also update Jujitsu configurations to use a
custom tool by default:

  https://jj-vcs.github.io/jj/latest/config/#editing-diffs

