package main

import (
	"bytes"
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
