package jutil

// NewRuntimeExceptionWithMessageAndCause is new RuntimeException(String
// message, Throwable cause).
func NewRuntimeExceptionWithMessageAndCause(message string, cause Throwable) *RuntimeException {
	e := &RuntimeException{}
	InitThrowable(e, "java.lang.RuntimeException", &message, cause)
	return e
}

// NewAssertionErrorWithCause is new AssertionError(Object detailMessage)
// with a Throwable: the message is the cause's toString() and the cause is
// set.
func NewAssertionErrorWithCause(cause Throwable) *AssertionError {
	e := &AssertionError{}
	InitThrowable(e, "java.lang.AssertionError", msgPtr(cause.Error()), cause)
	return e
}

// NewIllegalStateExceptionWithMessageAndCause is new
// IllegalStateException(String message, Throwable cause).
func NewIllegalStateExceptionWithMessageAndCause(message string, cause Throwable) *IllegalStateException {
	e := &IllegalStateException{}
	InitThrowable(e, "java.lang.IllegalStateException", &message, cause)
	return e
}
