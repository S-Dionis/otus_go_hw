package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	command := exec.Command("/bin/bash", "a")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	err := command.Run()

	fmt.Println(err)

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
