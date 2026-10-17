package jutil

import (
	"io"
	"os"
)

// Go-only: language-adapter — the java.nio.file.Files methods ported code
// uses, on the local file system. The java.nio.file exception subclasses
// (FileAlreadyExistsException, AccessDeniedException, ...) and their
// messages are not reproduced: a failure is an IOException carrying the Go
// error's message.

// Files_CreateDirectories is Files.createDirectories(Path): creates dir and
// its missing parents; an existing directory is not an error.
func Files_CreateDirectories(dir *Path) *Path {
	if dir == nil {
		// The JVM's message for the NullPointerException that
		// Files.createDirectories(null) throws.
		panic(NewHelpfulNullPointerException(`Cannot invoke "java.nio.file.Path.getFileSystem()" because "path" is null`))
	}
	if err := os.MkdirAll(dir.String(), 0o777); err != nil {
		panic(NewIOException(err.Error()))
	}
	return dir
}

// Files_NewOutputStream is Files.newOutputStream(Path, OpenOption...): with
// no options, CREATE, TRUNCATE_EXISTING and WRITE. READ, DELETE_ON_CLOSE,
// SPARSE, SYNC and DSYNC are not supported.
func Files_NewOutputStream(path *Path, options ...StandardOpenOption) io.WriteCloser {
	if len(options) == 0 {
		options = []StandardOpenOption{
			StandardOpenOption_CREATE, StandardOpenOption_TRUNCATE_EXISTING, StandardOpenOption_WRITE,
		}
	}
	flag := os.O_WRONLY
	for _, o := range options {
		switch o {
		case StandardOpenOption_WRITE:
		case StandardOpenOption_APPEND:
			flag |= os.O_APPEND
		case StandardOpenOption_TRUNCATE_EXISTING:
			flag |= os.O_TRUNC
		case StandardOpenOption_CREATE:
			flag |= os.O_CREATE
		case StandardOpenOption_CREATE_NEW:
			flag |= os.O_CREATE | os.O_EXCL
		default:
			panic(NewUnsupportedOperationExceptionMsg(o.String()))
		}
	}
	if flag&os.O_APPEND != 0 && flag&os.O_TRUNC != 0 {
		panic(NewIllegalArgumentExceptionMsg("APPEND + TRUNCATE_EXISTING not allowed"))
	}
	f, err := os.OpenFile(path.String(), flag, 0o666)
	if err != nil {
		panic(NewIOException(err.Error()))
	}
	return f
}
