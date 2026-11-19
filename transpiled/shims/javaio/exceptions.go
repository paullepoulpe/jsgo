// Clean-room specification based on java.io.IOException Javadoc.

package javaio

import (
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// IOException is the base class for I/O exceptions in Java.
type IOException struct{ javalang.JavaException }

func (*IOException) ioException() {}

type IIOException interface {
	javalang.IException
	ioException()
}

func NewIOException(msg string) *IOException {
	e := &IOException{}
	javalang.InitThrowable(e, "java.io.IOException", msg, nil)
	return e
}

func NewIOExceptionCause(cause javalang.Throwable) *IOException {
	e := &IOException{}
	javalang.InitThrowable(e, "java.io.IOException", "", cause)
	return e
}

// FileNotFoundException
type FileNotFoundException struct{ IOException }

func (*FileNotFoundException) fileNotFoundException() {}

type IFileNotFoundException interface {
	IIOException
	fileNotFoundException()
}

func NewFileNotFoundException(msg string) *FileNotFoundException {
	e := &FileNotFoundException{}
	javalang.InitThrowable(e, "java.io.FileNotFoundException", msg, nil)
	return e
}

// EOFException
type EOFException struct{ IOException }

func (*EOFException) eofException() {}

type IEOFException interface {
	IIOException
	eofException()
}

func NewEOFException(msg string) *EOFException {
	e := &EOFException{}
	javalang.InitThrowable(e, "java.io.EOFException", msg, nil)
	return e
}
