package jutil

// NewAssertionErrorWithMessageAndCause is new AssertionError(String message,
// Throwable cause).
func NewAssertionErrorWithMessageAndCause(message string, cause Throwable) *AssertionError {
	e := &AssertionError{}
	InitThrowable(e, "java.lang.AssertionError", &message, cause)
	return e
}
