package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envKey, enVal := range env {
		if enVal.NeedRemove {
			os.Unsetenv(envKey)
			continue
		}
		os.Setenv(envKey, enVal.Value)
	}
	cmdExec := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Run()
	return cmdExec.ProcessState.ExitCode()
}
