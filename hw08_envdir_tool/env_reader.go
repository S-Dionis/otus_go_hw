package main

import (
	"bufio"
	"bytes"
	"errors"
	"log"
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
		if strings.Contains(fileName, "=") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, ErrFileRead
		}

		line, err := readFirstLine(dir, info)

		if errors.Is(err, ErrFileIsEmpty) {
			remove = true
		} else if err != nil {
			return nil, ErrFileRead
		}

		value := replaceNullWithNewlineBytes(strings.TrimRight(line, " \t\r"))

		envs[fileName] = EnvValue{
			value, remove,
		}
	}

	return envs, nil
}

func readFirstLine(dir string, fileInfo os.FileInfo) (string, error) {
	if fileInfo.Size() == 0 {
		return "", ErrFileIsEmpty
	}

	filePath := filepath.Join(dir, fileInfo.Name())

	file, err := os.Open(filePath)
	// Доступен ли на чтение
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing file:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return "", scanner.Err()
	}

	return scanner.Text(), nil
}

func replaceNullWithNewlineBytes(b string) string {
	return string(bytes.ReplaceAll([]byte(b), []byte{0}, []byte{'\n'}))
}
