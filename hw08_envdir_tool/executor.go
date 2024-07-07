package main

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdName := cmd[2]
	cmdArgs := cmd[3:]

	for k, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				return 1
			}
		} else {
			err := os.Setenv(k, v.Value)
			if err != nil {
				return 1
			}
		}
	}

	command := exec.Command(cmdName, cmdArgs...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	err := command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			return 1
		}
	}

	return 0
}
