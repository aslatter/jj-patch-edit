package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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

func fakeApply(tokens iter.Seq[token]) (retErr error) {
	var buff bytes.Buffer
	for t := range tokens {
		for _, ln := range t.body {
			fmt.Fprintln(&buff, string(ln))
		}
	}
	fmt.Println("-- COLLECTED PATCH --")
	io.Copy(os.Stdout, &buff)
	return errors.New("fake apply - not applying")
}
