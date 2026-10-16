package jutil

import "strconv"

// Go-only: language-adapter — java.nio.file.InvalidPathException (JDK 21).

// InvalidPathException is java.nio.file.InvalidPathException.
type InvalidPathException struct {
	IllegalArgumentException
	input  string
	reason string
	index  int
}

func (*InvalidPathException) invalidPathException() {}

// IInvalidPathException matches InvalidPathException.
type IInvalidPathException interface {
	IIllegalArgumentException
	invalidPathException()
}

// NewInvalidPathException is new InvalidPathException(input, reason).
func NewInvalidPathException(input, reason string) *InvalidPathException {
	return NewInvalidPathExceptionAt(input, reason, -1)
}

// NewInvalidPathExceptionAt is new InvalidPathException(input, reason, index).
func NewInvalidPathExceptionAt(input, reason string, index int) *InvalidPathException {
	e := &InvalidPathException{input: input, reason: reason, index: index}
	msg := reason
	if index > -1 {
		msg += " at index " + strconv.Itoa(index)
	}
	msg += ": " + input
	InitThrowable(e, "java.nio.file.InvalidPathException", &msg, nil)
	return e
}

// GetInput is InvalidPathException.getInput().
func (e *InvalidPathException) GetInput() string { return e.input }

// GetReason is InvalidPathException.getReason().
func (e *InvalidPathException) GetReason() string { return e.reason }

// GetIndex is InvalidPathException.getIndex().
func (e *InvalidPathException) GetIndex() int { return e.index }
