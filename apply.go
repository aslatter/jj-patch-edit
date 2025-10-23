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

func apply(rightPath string, files iter.Seq[*file]) (retErr error) {
	cmd := exec.Command("git",
		"apply",
		"--allow-empty",
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
		defer inWriter.Close()
		return printFiles(inWriter, files)
	})

	err = wg.Wait()

	return err
}

func fakeApply(files iter.Seq[*file]) (retErr error) {
	var buff bytes.Buffer
	err := printFiles(&buff, files)
	if err != nil {
		return err
	}
	fmt.Println("-- COLLECTED PATCH --")
	io.Copy(os.Stdout, &buff)
	return errors.New("fake apply - not applying")
}
