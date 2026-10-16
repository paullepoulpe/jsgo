package jutil

// Go-only: language-adapter — Guava's UncheckedExecutionException and
// CacheLoader.InvalidCacheLoadException (shaded class names, as they appear
// in the pinned compiler jar).

// UncheckedExecutionException is
// com.google.common.util.concurrent.UncheckedExecutionException.
type UncheckedExecutionException struct{ RuntimeException }

func (*UncheckedExecutionException) uncheckedExecutionException() {}

// IUncheckedExecutionException matches UncheckedExecutionException.
type IUncheckedExecutionException interface {
	IRuntimeException
	uncheckedExecutionException()
}

// NewUncheckedExecutionException is new UncheckedExecutionException(cause).
func NewUncheckedExecutionException(cause Throwable) *UncheckedExecutionException {
	e := &UncheckedExecutionException{}
	InitThrowable(e, "com.google.javascript.jscomp.jarjar.com.google.common.util.concurrent.UncheckedExecutionException", nil, cause)
	return e
}

// CacheLoader_InvalidCacheLoadException is
// com.google.common.cache.CacheLoader.InvalidCacheLoadException.
type CacheLoader_InvalidCacheLoadException struct{ RuntimeException }

// NewCacheLoader_InvalidCacheLoadException is new InvalidCacheLoadException(message).
func NewCacheLoader_InvalidCacheLoadException(message string) *CacheLoader_InvalidCacheLoadException {
	e := &CacheLoader_InvalidCacheLoadException{}
	InitThrowable(e, "com.google.javascript.jscomp.jarjar.com.google.common.cache.CacheLoader$InvalidCacheLoadException", &message, nil)
	return e
}

// ExecutionError is com.google.common.util.concurrent.ExecutionError.
type ExecutionError struct{ JavaLangError }

// NewExecutionError is new ExecutionError(cause).
func NewExecutionError(cause Throwable) *ExecutionError {
	e := &ExecutionError{}
	InitThrowable(e, "com.google.javascript.jscomp.jarjar.com.google.common.util.concurrent.ExecutionError", nil, cause)
	return e
}
