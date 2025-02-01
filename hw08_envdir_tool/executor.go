package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	prevEnv := os.Environ()
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
	cmdExec.Stdin = os.Stdin
	cmdExec.Env = os.Environ()
	cmdExec.Run()
	os.Clearenv()
	for _, envVar := range prevEnv {
		parts := strings.Split(envVar, "=")
		os.Setenv(parts[0], parts[1])
	}
	return cmdExec.ProcessState.ExitCode()
}
