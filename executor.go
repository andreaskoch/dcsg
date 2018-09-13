package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Executor interface {
	Run(command ...string) error
}

func newExecutor(in io.Reader, out, err io.Writer, baseDirectory string, dryRun bool) Executor {
	return &CommandExecutor{
		in:  in,
		out: out,
		err: err,

		baseDirectory: baseDirectory,
		dryRun:        dryRun,
	}
}

type CommandExecutor struct {
	in  io.Reader
	out io.Writer
	err io.Writer

	baseDirectory string
	dryRun        bool
}

func (executor *CommandExecutor) InDirectory(directory string) CommandExecutor {

	newWorkingDirectory := filepath.Join(executor.baseDirectory, directory)
	expandedWorkingDirectory := os.ExpandEnv(newWorkingDirectory)

	return CommandExecutor{
		in:  executor.in,
		out: executor.out,
		err: executor.err,

		baseDirectory: expandedWorkingDirectory,
		dryRun:        executor.dryRun,
	}
}

func (executor *CommandExecutor) Run(command ...string) error {

	if len(command) == 0 {
		return fmt.Errorf("No command given")
	}

	workingDirectory := executor.baseDirectory
	expandedWorkingDirectory := os.ExpandEnv(workingDirectory)
	expandedCommandName := os.ExpandEnv(command[0])
	var expandedArguments []string
	for _, argument := range command[1:] {
		expandedArguments = append(expandedArguments, os.ExpandEnv(argument))
	}

	cmd := exec.Command(expandedCommandName, expandedArguments...)

	cmd.Dir = expandedWorkingDirectory
	cmd.Env = os.Environ()

	cmd.Stdout = executor.out
	cmd.Stderr = executor.err
	cmd.Stdin = executor.in

	log.Printf("%s: %s %s", expandedWorkingDirectory, command[0], strings.Join(command[1:], " "))

	if !executor.dryRun {
		err := cmd.Run()
		if err != nil {
			log.Printf("Error running %s: %v", command, err)
			return err
		}
	}

	return nil
}
