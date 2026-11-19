//go:build !unix

package javanio

import "errors"

const hasFastFDRead = false

func openReadOnlyFD(string) (int, error) { return -1, errors.New("unsupported") }
func readFD(int, []byte) (int, error)   { return 0, errors.New("unsupported") }
func closeFD(int)                       {}
func fstatRegularSize(int) (int64, bool) { return 0, false }
