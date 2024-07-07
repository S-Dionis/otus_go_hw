package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <envdir> <command> [args...]", os.Args[0])
	}
	envDir := os.Args[1]
	envs, err := ReadDir(envDir)
	if err != nil {
		panic(err) // TODO replace with log maybe?
	}

	code := RunCmd(os.Args, envs)
	os.Exit(code)
}
