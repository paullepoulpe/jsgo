//go:build !unix

package jutil

import "os"

func fileIsDirectoryRaw(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}
