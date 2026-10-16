package jutil

// NewRuntimeExceptionWithMessageAndCause is new RuntimeException(String
// message, Throwable cause).
func NewRuntimeExceptionWithMessageAndCause(message string, cause Throwable) *RuntimeException {
	e := &RuntimeException{}
	InitThrowable(e, "java.lang.RuntimeException", &message, cause)
	return e
}
