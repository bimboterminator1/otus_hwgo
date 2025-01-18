package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func deepCompare(file1, file2 string) bool {
	sf, err := os.Open(file1)
	if err != nil {
		return false
	}

	defer sf.Close()
	df, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer df.Close()
	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return false
		}
	}

	return true
}

func TestCopy(t *testing.T) {
	tests := []struct {
		testname string
		fromPath string
		toPath   string
		limit    int64
		offset   int64
		expected string
	}{
		{
			"Simple copy",
			"testdata/input.txt",
			"out.txt",
			0,
			0,
			"testdata/out_offset0_limit0.txt",
		},
		{
			"limit 10",
			"testdata/input.txt",
			"out3.txt",
			10,
			0,
			"testdata/out_offset0_limit10.txt",
		},
		{
			"limit 1000",
			"testdata/input.txt",
			"out4.txt",
			1000,
			0,
			"testdata/out_offset0_limit1000.txt",
		},
		{
			"limit 10000",
			"testdata/input.txt",
			"out5.txt",
			10000,
			0,
			"testdata/out_offset0_limit10000.txt",
		},
		{
			"offset 100 limit 1000",
			"testdata/input.txt",
			"out6.txt",
			1000,
			100,
			"testdata/out_offset100_limit1000.txt",
		},
		{
			"offset 6000 limit 1000",
			"testdata/input.txt",
			"out7.txt",
			1000,
			6000,
			"testdata/out_offset6000_limit1000.txt",
		},
	}
	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			err := Copy(test.fromPath, test.toPath, test.offset, test.limit)
			require.NoError(t, err)
			defer os.Remove(test.toPath)
			require.True(t, deepCompare(test.expected, test.toPath))
		})
	}
	t.Run("limit is larger than fileSize", func(t *testing.T) {
		fromPath := "testdata/input.txt"
		fromFile, err := os.Open(fromPath)
		require.NoError(t, err)
		info, err := fromFile.Stat()
		require.NoError(t, err)
		err = Copy(fromPath, "out.txt", 0, info.Size()+10000)
		defer os.Remove("out.txt")
		require.NoError(t, err)
	})
	t.Run("offset is larger than file Size", func(t *testing.T) {
		fromPath := "testdata/input.txt"
		fromFile, err := os.Open(fromPath)
		require.NoError(t, err)
		info, err := fromFile.Stat()
		require.NoError(t, err)
		err = Copy(fromPath, "out.txt", info.Size()+10000, 0)
		defer os.Remove("out.txt")
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})
	t.Run("file length unknown", func(t *testing.T) {
		fromPath := "/dev/urandom"
		err := Copy(fromPath, "out.txt", 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
}
