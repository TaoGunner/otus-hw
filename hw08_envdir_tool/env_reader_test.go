package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	envs, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	require.Equal(t, "bar", envs["BAR"].Value)
	require.False(t, envs["BAR"].NeedRemove)

	require.Equal(t, "", envs["EMPTY"].Value)
	require.True(t, envs["EMPTY"].NeedRemove)

	require.Equal(t, "   foo\nwith new line", envs["FOO"].Value)
	require.False(t, envs["FOO"].NeedRemove)

	require.Equal(t, "\"hello\"", envs["HELLO"].Value)
	require.False(t, envs["HELLO"].NeedRemove)

	require.Equal(t, "", envs["UNSET"].Value)
	require.True(t, envs["UNSET"].NeedRemove)
}
