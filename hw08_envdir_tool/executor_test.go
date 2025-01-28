package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("normal run", func(t *testing.T) {
		os.Setenv("HELLO", "dummy")
		env := make(Environment)
		env["HELLO"] = EnvValue{
			"HELLO",
			false,
		}
		rc := RunCmd([]string{"/bin/bash", "./testdata/echo.sh", "arg"}, env)
		require.Equal(t, 0, rc)
		evar, ok := os.LookupEnv("HELLO")
		require.True(t, ok)
		require.Equal(t, "HELLO", evar)
	})
}
