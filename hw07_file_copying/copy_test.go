package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name    string
		srcPath string
		dstPath string
		offset  int64
		limit   int64
		err     error
	}{
		{
			name:    "test offset=0 and limit=0",
			srcPath: "testdata/input.txt",
			dstPath: "out_offset0_limit0.txt",
			offset:  0,
			limit:   0,
		},
		{
			name:    "test offset=0 and limit=10",
			srcPath: "testdata/input.txt",
			dstPath: "out_offset0_limit10.txt",
			offset:  0,
			limit:   10,
		},
		{
			name:    "test offset=0 and limit=1000",
			srcPath: "testdata/input.txt",
			dstPath: "out_offset0_limit1000.txt",
			offset:  0,
			limit:   1000,
		},
		{
			name:    "test offset=0 and limit=10000",
			srcPath: "testdata/input.txt",
			dstPath: "out_offset0_limit10000.txt",
			offset:  0,
			limit:   10000,
		},
		{
			name:    "test offset=100 and limit=1000",
			srcPath: "testdata/input.txt",
			dstPath: "out_offset100_limit1000.txt",
			offset:  100,
			limit:   1000,
		},
		{
			name:    "test offset=6000 and limit=1000",
			srcPath: "testdata/input.txt",
			dstPath: "out_offset6000_limit1000.txt",
			offset:  6000,
			limit:   1000,
		},
		{
			name:    "test limit more than file size",
			srcPath: "testdata/input.txt",
			dstPath: "out_offset0_limit0.txt",
			offset:  0,
			limit:   10000000000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpFile, err := os.CreateTemp("", "hw06test_")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			err = Copy(tt.srcPath, tmpFile.Name(), tt.offset, tt.limit)
			require.NoError(t, err)

			expected, err := os.ReadFile("testdata/" + tt.dstPath)
			require.NoError(t, err)

			actual, err := os.ReadFile(tmpFile.Name())
			require.NoError(t, err)

			require.Equal(t, expected, actual)
		})
	}

	t.Run("offset > file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "result.txt", 8000, 6000)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("unsupported source file", func(t *testing.T) {
		err := Copy("/dev/urandom", "result.txt", 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("source file not exists", func(t *testing.T) {
		err := Copy("non_existent_file.txt", "result.txt", 0, 0)
		require.Error(t, err)
	})
}
