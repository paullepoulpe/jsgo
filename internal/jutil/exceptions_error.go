package jutil

// Go-only: language-adapter — java.lang.Error and AssertionError (not
// RuntimeExceptions: `catch (RuntimeException e)` sites do not see them).

// JavaLangError is java.lang.Error (named so that embedding it does not
// shadow the Go error interface's Error method).
type JavaLangError struct{ ThrowableBase }

func (*JavaLangError) javaLangError() {}

// IError matches java.lang.Error and all its subclasses.
type IError interface {
	Throwable
	javaLangError()
}

// AssertionError is java.lang.AssertionError.
type AssertionError struct{ JavaLangError }

// NewAssertionError is new AssertionError(Object detailMessage) with a
// string message.
func NewAssertionError(message string) *AssertionError {
	e := &AssertionError{}
	InitThrowable(e, "java.lang.AssertionError", msgPtr(message), nil)
	return e
}

// NewAssertionErrorNoMessage is new AssertionError().
func NewAssertionErrorNoMessage() *AssertionError {
	e := &AssertionError{}
	InitThrowable(e, "java.lang.AssertionError", nil, nil)
	return e
}
