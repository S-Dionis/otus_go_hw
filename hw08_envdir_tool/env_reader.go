package main

import (
	"bytes"
	"errors"
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

var (
	ErrFileIsEmpty = errors.New("file is empty")
	ErrFileRead    = errors.New("file read error")
	ErrDirRead     = errors.New("directory read error")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrDirRead
	}

	envs := make(Environment)

	remove := false

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		filePath := filepath.Join(dir, fileName)
		line, err := readFirstLine(filePath)

		if errors.Is(err, ErrFileIsEmpty) {
			remove = true
		} else if err != nil {
			return nil, ErrFileRead
		}

		value := replaceNullWithNewlineBytes(strings.TrimRight(line, " \t\r"))
		fileName = strings.ReplaceAll(fileName, "=", "")

		envs[fileName] = EnvValue{
			value, remove,
		}
	}

	return envs, nil
}

func readFirstLine(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", ErrFileIsEmpty
	}

	return strings.Split(string(content), "\n")[0], nil
}

func replaceNullWithNewlineBytes(b string) string {
	return string(bytes.ReplaceAll([]byte(b), []byte{0}, []byte{'\n'}))
}
