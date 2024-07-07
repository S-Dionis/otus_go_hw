package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("correct read of files", func(t *testing.T) {
		os.Setenv("UNSET", "MUR")
		os.Setenv("EMPTY", "")

		expected := Environment{
			"BAR":     {"bar", false},
			"EMPTY":   {"", false},
			"FOO":     {"   foo\nwith new line", false},
			"HELLO":   {"\"hello\"", false},
			"MURMEOW": {"without equals sign", false},
			"UNSET":   {"", true},
		}

		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Equal(t, expected, env)
	})

	t.Run("path not found test", func(t *testing.T) {
		_, err := ReadDir("./mur")
		require.ErrorIs(t, err, ErrDirRead)
	})
}
