package main

import (
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func ActionWithLogError(fn func() error) {
	err := fn()
	if err != nil {
		log.Println(err)
	}
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TestCopy(t *testing.T) {
	path := "testdata/testFile.txt"

	t.Run("Tasks without errors", func(t *testing.T) {
		err := Copy("testdata/input.txt", path, 0, 0)
		require.NoError(t, err)
		inputFile, err := os.Open("testdata/input.txt")
		require.NoError(t, err)
		logErr(err)

		testFile, _ := os.Open(path)
		inputFileStat, _ := inputFile.Stat()
		testFileStat, _ := testFile.Stat()
		expectedSize := inputFileStat.Size()
		actualSize := testFileStat.Size()

		ActionWithLogError(func() error { return inputFile.Close() })
		ActionWithLogError(func() error { return testFile.Close() })
		ActionWithLogError(func() error { return os.Remove(path) })

		require.Truef(t, expectedSize == actualSize, "Equality of file size")
	})

	t.Run("Errors on copying", func(t *testing.T) {
		err := Copy("testdata/input.txt", "./testFile", -1, 10000)
		require.Truef(t, errors.Is(err, ErrWrongOffsetValue), "actual err - %v", err)
		err = Copy("testdata/input.txt", "./testFile", 0, -1)
		require.Truef(t, errors.Is(err, ErrWrongLimitValue), "actual err - %v", err)
	})

	t.Run("Offset and limit test", func(t *testing.T) {
		err := Copy("testdata/input.txt", path, 1000, 100)
		require.NoError(t, err)

		// dd if=input.txt of=output_1000_100.txt bs=1 skip=1000 count=100
		expected, _ := os.ReadFile("testdata/output_1000_100.txt")
		file, _ := os.ReadFile("testdata/testFile.txt")
		content := string(file)
		expectedContent := string(expected)

		compare := strings.Compare(content, expectedContent)
		require.Truef(t, compare == 0, "File content as expected")

		err = os.Remove(path)
		if err != nil {
			log.Fatal(err)
			return
		}
	})
}
