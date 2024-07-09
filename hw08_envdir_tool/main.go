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
		log.Fatalf("Error reading dir %s: %s", envDir, err)
	}

	code := RunCmd(os.Args[2:], envs)
	os.Exit(code)
}
