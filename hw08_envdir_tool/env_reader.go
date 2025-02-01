package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrUnsupportedFormat error

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, entry := range dirEntries {
		fi, _ := entry.Info()
		if strings.Contains(entry.Name(), "=") {
			return nil, ErrUnsupportedFormat
		}
		if fi.Size() == 0 {
			env[entry.Name()] = EnvValue{"", true}
			continue
		}
		file, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			envVal := scanner.Text()
			envVal = strings.ReplaceAll(envVal, "\000", "\n")
			envVal = strings.TrimRight(envVal, " \t\r\v\f")
			env[entry.Name()] = EnvValue{envVal, false}
			break
		}
		file.Close()
	}

	return env, nil
}
