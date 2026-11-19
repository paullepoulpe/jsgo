//go:build unix

package javaio

import (
	"errors"
	"io"
	"os"
	"syscall"
)

const hasRawFileFD = true

func openRawReadFD(path string) (int, error) {
	for {
		fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
		if err == syscall.EINTR {
			continue
		}
		return fd, err
	}
}

func openRawWriteFD(path string, osFlags int) (int, error) {
	sysFlags := syscall.O_CLOEXEC
	switch osFlags & (os.O_RDONLY | os.O_WRONLY | os.O_RDWR) {
	case os.O_RDWR:
		sysFlags |= syscall.O_RDWR
	case os.O_WRONLY:
		sysFlags |= syscall.O_WRONLY
	default:
		sysFlags |= syscall.O_RDONLY
	}
	if osFlags&os.O_APPEND != 0 {
		sysFlags |= syscall.O_APPEND
	}
	if osFlags&os.O_CREATE != 0 {
		sysFlags |= syscall.O_CREAT
	}
	if osFlags&os.O_EXCL != 0 {
		sysFlags |= syscall.O_EXCL
	}
	if osFlags&os.O_TRUNC != 0 {
		sysFlags |= syscall.O_TRUNC
	}
	for {
		fd, err := syscall.Open(path, sysFlags, 0o666)
		if err == syscall.EINTR {
			continue
		}
		return fd, err
	}
}

func readRawFD(fd int, p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	for {
		n, err := syscall.Read(fd, p)
		if err == syscall.EINTR {
			continue
		}
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, io.EOF
		}
		return n, nil
	}
}

func writeRawFD(fd int, p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	written := 0
	for written < len(p) {
		n, err := syscall.Write(fd, p[written:])
		if n > 0 {
			written += n
		}
		if err == syscall.EINTR {
			continue
		}
		if err != nil {
			return written, err
		}
		if n == 0 {
			return written, io.ErrShortWrite
		}
	}
	return written, nil
}

func closeRawFD(fd int) error {
	err := syscall.Close(fd)
	if err == nil || err == syscall.EINTR || err == syscall.EBADF {
		return nil
	}
	return err
}

func remainingRegularRawFD(fd int) (int64, bool) {
	var st syscall.Stat_t
	if err := syscall.Fstat(fd, &st); err != nil || (st.Mode&syscall.S_IFMT) != syscall.S_IFREG {
		return 0, false
	}
	pos, err := syscall.Seek(fd, 0, io.SeekCurrent)
	if err != nil || st.Size < pos {
		return 0, false
	}
	return st.Size - pos, true
}

func isRawNotExist(err error) bool {
	return os.IsNotExist(err) || errors.Is(err, syscall.ENOENT)
}
