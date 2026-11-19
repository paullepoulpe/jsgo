//go:build unix

package javanio

import "syscall"

const hasFastFDRead = true

func openReadOnlyFD(path string) (int, error) {
	for {
		fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
		if err == syscall.EINTR {
			continue
		}
		return fd, err
	}
}

func readFD(fd int, p []byte) (int, error) {
	for {
		n, err := syscall.Read(fd, p)
		if err == syscall.EINTR {
			continue
		}
		return n, err
	}
}

func closeFD(fd int) {
	_ = syscall.Close(fd)
}

func fstatRegularSize(fd int) (int64, bool) {
	var st syscall.Stat_t
	if err := syscall.Fstat(fd, &st); err == nil && (st.Mode&syscall.S_IFMT) == syscall.S_IFREG {
		return st.Size, true
	}
	return 0, false
}
