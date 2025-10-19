package main

import (
	"fmt"
	"iter"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

func apply(rightPath string, tokens iter.Seq[token]) (retErr error) {
	cmd := exec.Command("git",
		"apply",
		"--",
	)
	cmd.Dir = rightPath

	// leave stdin and stdout wired up to current?
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	inWriter, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("setting up input pipe: %s", err)
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	var wg errgroup.Group

	wg.Go(func() error {
		return cmd.Wait()
	})
	wg.Go(func() error {
		for t := range tokens {
			for _, ln := range t.body {
				_, err := fmt.Fprintln(inWriter, string(ln))
				if err != nil {
					return err
				}
			}
		}
		return inWriter.Close()
	})

	err = wg.Wait()

	return err
}
