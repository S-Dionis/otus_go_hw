package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	t.Run("correct read of ", func(t *testing.T) {
		os.Setenv("UNSET", "MUR")
		environment := Environment{
			"BAR":   {"bar", false},
			"UNSET": {"", true},
		}

		args := []string{"exit 0", ""}

		_ = RunCmd(args, environment)
		assert.Equal(t, "bar", os.Getenv("BAR"))
		assert.Equal(t, "", os.Getenv("UNSET"))
	})

	t.Run("correct read of ", func(t *testing.T) {
		environment := Environment{}

		args := []string{"exit 1", ""}

		code := RunCmd(args, environment)
		assert.Equal(t, 1, code)
	})
}
