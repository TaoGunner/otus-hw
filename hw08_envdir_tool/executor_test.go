package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	exitCode := RunCmd([]string{"go", "env"}, nil)
	require.Equal(t, 0, exitCode)

	exitCode = RunCmd([]string{"go", "env1"}, nil)
	require.NotEqual(t, 0, exitCode)
}
