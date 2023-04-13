package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// serviceRunner stores all the
// data needed to run a go service
type serviceRunner struct {
	name        string
	workDir     string
	envs        []string
	initScripts []string
}

// newServiceRunner creates a new serviceRunner
func newServiceRunner(s Service) (serviceRunner, error) {
	var (
		workDir string
		envs    []string
	)

	if s.Dir != "" {
		wd, err := os.Getwd()
		if err != nil {
			return serviceRunner{}, nil
		}
		workDir = filepath.Join(wd, s.Dir)
	}

	if len(s.Envs) > 0 {
		envs = append(
			os.Environ(),
			s.Envs...,
		)
	}

	return serviceRunner{
		name:        s.Name,
		workDir:     workDir,
		envs:        envs,
		initScripts: s.InitScripts,
	}, nil
}

// run starts an instance of a go routine
func (sr *serviceRunner) run(ctx context.Context) error {
	cmd := exec.Command("go", "run", sr.workDir)
	cmd.Dir = sr.workDir
	cmd.Env = sr.envs

	prefix := fmt.Sprintf("%s service: ", sr.name)
	cmd.Stdout = &prefixer{prefix: prefix, w: os.Stdout}
	cmd.Stderr = &prefixer{prefix: prefix, w: os.Stderr}

	if err := cmd.Start(); err != nil {
		log.Printf("ERROR starting %s service: %v\n", sr.name, err)
		return err
	}

	log.Printf("%s service PID: %d\n", sr.name, cmd.Process.Pid)

	go func() {
		<-ctx.Done()
		// Now I need something to catch the output from
		// the rest of the go service to make sure it gracefully shuts down? Because I don't get Stdout or Stderr aftewards
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("ERROR signaling %s service for graceful shutdown: %s \n", sr.name, err.Error())
		}

		if err := ctx.Err(); err != nil {
			if !errors.Is(ctx.Err(), context.Canceled) {
				log.Printf("CONTEXT ERROR: %s service: %s\n", sr.name, err)
			}
		}
	}()

	return cmd.Wait()
}

// intialize runs all the init scripts for
// a service before running the service
func (sr *serviceRunner) initialize() error {
	if len(sr.initScripts) == 0 {
		return nil
	}

	for _, script := range sr.initScripts {
		if err := runSH(script); err != nil {
			return err
		}
	}

	return nil
}

// runSH runs a shell or bash script.
// Pass it the script's name as the first argument,
// and then pass in additional parameters.
func runSH(args ...string) error {
	if len(args) == 0 {
		return errors.New("must pass in at least one argument")
	}

	cmd := exec.Command("/bin/sh", args...)

	prefix := fmt.Sprintf("SH %s: ", args[0])
	cmd.Stdout = &prefixer{prefix: prefix, w: os.Stdout}
	cmd.Stderr = &prefixer{prefix: prefix, w: os.Stderr}

	// This doesn't always catch the scripts error
	// if the script doesn't return an exit code > 0?
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// result is used to store the results
// of running a script in a channel
type result struct {
	output string
	err    error
}

func runInitScripts(args ...string) <-chan result {
	resChan := make(chan result)

	go func() {
		cmd := exec.Command("/bin/sh", args...)

		// I don't think this will error if the script fails,
		// so it has no impact on the functions that depend
		// on this script running successfully.
		// Maybe I'll get the stdout and stderr, and use the stderr to return an error
		// A script with an exit code 1 might trigger an error
		output, err := cmd.CombinedOutput()
		resChan <- result{
			output: string(output),
			err:    err,
		}

		close(resChan)
	}()

	return resChan
}

type stopContainer func() error

// startContainer is calls a script to start a docker container,
// and returns a function that is used to clean up the docker container.
func startContainer(args ...string) (stopContainer, error) {
	resChan := runInitScripts(args...)

	result := <-resChan
	if err := result.err; err != nil {
		return nil, err
	}

	containerID := strings.TrimSpace(result.output)

	log.Println("received container id:", containerID)

	return func() error {
		return runSH("./teardown_container.sh", containerID)
	}, nil
}
