//go:build !unix

package javaio

import (
	"errors"
	"os"
)

const hasRawFileFD = false

func openRawReadFD(string) (int, error)        { return -1, errors.New("unsupported") }
func openRawWriteFD(string, int) (int, error)  { return -1, errors.New("unsupported") }
func readRawFD(int, []byte) (int, error)       { return 0, errors.New("unsupported") }
func writeRawFD(int, []byte) (int, error)      { return 0, errors.New("unsupported") }
func closeRawFD(int) error                     { return nil }
func remainingRegularRawFD(int) (int64, bool)  { return 0, false }
func isRawNotExist(err error) bool             { return os.IsNotExist(err) }
