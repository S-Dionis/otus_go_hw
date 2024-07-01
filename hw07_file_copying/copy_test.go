package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
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
	fromPath := "testdata/input.txt"
	toPath := "testdata/testFile.txt"

	t.Run("Errors on copying", func(t *testing.T) {
		err := Copy(fromPath, toPath, -1, 10000)
		require.ErrorIs(t, err, ErrWrongOffsetValue, "actual err - %v", err)
		err = Copy(fromPath, toPath, 0, -1)
		require.ErrorIs(t, err, ErrWrongLimitValue, "actual err - %v", err)
	})

	t.Run("Errors on pass the same file", func(t *testing.T) {
		err := Copy(fromPath, fromPath, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile, "actual err - %v", err)
	})

	t.Run("Tasks without errors", func(t *testing.T) {
		err := Copy(fromPath, toPath, 0, 0)
		require.NoError(t, err)
		inputFile, err := os.Open(fromPath)
		require.NoError(t, err)
		logErr(err)

		testFile, err := os.Open(toPath)
		require.NoError(t, err)
		inputFileStat, err := inputFile.Stat()
		require.NoError(t, err)
		testFileStat, err := testFile.Stat()
		require.NoError(t, err)
		expectedSize := inputFileStat.Size()
		actualSize := testFileStat.Size()

		ActionWithLogError(func() error { return inputFile.Close() })
		ActionWithLogError(func() error { return testFile.Close() })
		ActionWithLogError(func() error { return os.Remove(toPath) })

		require.Truef(t, expectedSize == actualSize, "Equality of file size")
	})

	t.Run("Offset and limit test", func(t *testing.T) {
		err := Copy(fromPath, toPath, 1000, 100)
		require.NoError(t, err)

		// dd if=input.txt of=output_1000_100.txt bs=1 skip=1000 count=100
		expected, err := os.ReadFile("testdata/output_1000_100.txt")
		require.NoError(t, err)
		file, err := os.ReadFile(toPath)
		require.NoError(t, err)
		actualContent := string(file)
		expectedContent := string(expected)

		assert.Equal(t, expectedContent, actualContent)

		err = os.Remove(toPath)
		require.NoError(t, err)
	})
}
