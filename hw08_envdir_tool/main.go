package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s : wrong arguments\n", os.Args[0])
		os.Exit(1)
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s : %s\n", os.Args[0], err)
		os.Exit(1)
	}
	rc := RunCmd(os.Args[2:], env)
	os.Exit(rc)
}
