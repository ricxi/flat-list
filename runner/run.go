package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

type goService struct {
	name    string
	workDir string
	envs    []string
}

// run starts an instance of a go routine
func (gs *goService) run(wg *sync.WaitGroup) error {
	defer wg.Done()

	cmd := exec.Command("go", "run", gs.workDir)
	cmd.Dir = gs.workDir
	cmd.Env = gs.envs

	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	if err := cmd.Start(); err != nil {
		log.Printf("problem starting %s service: %v\n", gs.name, err)
		return err
	}

	log.Printf("%s service PID: %d\n", gs.name, cmd.Process.Pid)

	return cmd.Wait()
}

// runSH runs a shell or bash script.
// Pass it the script's name as the first argument,
// and then pass in additional parameters.
func runSH(args ...string) error {
	if len(args) == 0 {
		return errors.New("must pass in at least one argument")
	}

	cmd := exec.Command("/bin/sh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Printf("container: %s\n", output)
	return nil
}
