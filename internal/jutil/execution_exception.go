package jutil

// ExecutionException is java.util.concurrent.ExecutionException: the
// checked exception thrown when retrieving the result of a task that
// aborted by throwing an exception.
type ExecutionException struct{ JavaException }

func (*ExecutionException) executionException() {}

// IExecutionException matches ExecutionException and its subclasses.
type IExecutionException interface {
	IException
	executionException()
}

// NewExecutionException is new ExecutionException(Throwable cause): the
// message is cause.toString().
func NewExecutionException(cause Throwable) *ExecutionException {
	e := &ExecutionException{}
	InitThrowable(e, "java.util.concurrent.ExecutionException", nil, cause)
	return e
}
