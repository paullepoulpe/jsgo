//go:build unix

package jutil

import "syscall"

func fileIsDirectoryRaw(path string) bool {
	if path == "" {
		return false
	}
	var st syscall.Stat_t
	err := syscall.Stat(path, &st)
	return err == nil && (st.Mode&syscall.S_IFMT) == syscall.S_IFDIR
}
