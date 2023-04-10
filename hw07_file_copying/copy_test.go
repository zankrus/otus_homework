package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("empty files", func(t *testing.T) {
		err := Copy("", "", 0, 0)
		require.EqualError(t, err, "empty path")

		err = Copy("/etc/hosts", "", 0, 0)
		require.EqualError(t, err, "empty path")
	})

	t.Run("unlimit file", func(t *testing.T) {
		err := Copy("/dev/urandom", "/some", 0, 0)

		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write([]byte("")); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}

		err = Copy(tmpfile.Name(), "/some", 1, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})
}
