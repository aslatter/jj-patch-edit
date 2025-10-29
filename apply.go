package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"os/exec"
	"sync"
)

// apply feeds in patch-contents ('files') into 'git-apply'
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

	var printErr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer inWriter.Close()
		printErr = printFiles(inWriter, files)
	}()

	err = cmd.Wait()
	wg.Wait()

	return errors.Join(err, printErr)
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
