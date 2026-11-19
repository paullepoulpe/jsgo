// Port of java.lang.Throwable and standard exception hierarchy.
// Clean-room specification based on Java Language Specification & Java SE Javadoc.

package javalang

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

// Throwable is the root interface of the Java exception hierarchy.
type Throwable interface {
	error
	GetMessage() string
	GetCause() Throwable
	throwable() *ThrowableBase
}

// ThrowableBase holds shared state for all Java exceptions.
type ThrowableBase struct {
	className string
	message   string
	hasMsg    bool
	cause     Throwable
}

func (t *ThrowableBase) throwable() *ThrowableBase { return t }

// GetMessage returns the exception detail message (or empty string if none).
func (t *ThrowableBase) GetMessage() string { return t.message }

// HasMessage reports whether an explicit detail message was provided.
func (t *ThrowableBase) HasMessage() bool { return t.hasMsg }

// GetCause returns the underlying cause of this exception, or nil.
func (t *ThrowableBase) GetCause() Throwable { return t.cause }

// InitCause initializes the cause of this exception.
func (t *ThrowableBase) InitCause(c Throwable) { t.cause = c }

// Error formats the exception as "<ClassName>: <message>".
func (t *ThrowableBase) Error() string {
	if !t.hasMsg {
		return t.className
	}
	return t.className + ": " + t.message
}

// InitThrowable initializes the ThrowableBase fields of an exception.
func InitThrowable(t Throwable, className, message string, cause Throwable) {
	b := t.throwable()
	b.className = className
	if message != "" {
		b.message, b.hasMsg = message, true
	} else if cause != nil {
		b.message, b.hasMsg = cause.Error(), true
	}
	b.cause = cause
}

func initThrowable(t Throwable, className, message string, hasMsg bool, cause Throwable) {
	b := t.throwable()
	b.className = className
	if hasMsg {
		b.message, b.hasMsg = message, true
	} else if cause != nil {
		b.message, b.hasMsg = cause.Error(), true
	}
	b.cause = cause
}

// ---------------------------------------------------------------------------
// Standard Exception and Error Hierarchies
// ---------------------------------------------------------------------------

// JavaException corresponds to java.lang.Exception.
type JavaException struct{ ThrowableBase }

func (*JavaException) exception() {}

// IException matches JavaException and all its subclasses.
type IException interface {
	Throwable
	exception()
}

func NewException(msg string) *JavaException {
	e := &JavaException{}
	initThrowable(e, "java.lang.Exception", msg, true, nil)
	return e
}

// RuntimeException corresponds to java.lang.RuntimeException.
type RuntimeException struct{ JavaException }

func (*RuntimeException) runtimeException() {}

// IRuntimeException matches RuntimeException and all its subclasses.
type IRuntimeException interface {
	IException
	runtimeException()
}

func NewRuntimeException() *RuntimeException {
	e := &RuntimeException{}
	initThrowable(e, "java.lang.RuntimeException", "", false, nil)
	return e
}

func NewRuntimeExceptionMsg(msg string) *RuntimeException {
	e := &RuntimeException{}
	initThrowable(e, "java.lang.RuntimeException", msg, true, nil)
	return e
}

func NewRuntimeExceptionCause(cause Throwable) *RuntimeException {
	e := &RuntimeException{}
	initThrowable(e, "java.lang.RuntimeException", "", false, cause)
	return e
}

// IllegalArgumentException
type IllegalArgumentException struct{ RuntimeException }

func (*IllegalArgumentException) illegalArgumentException() {}

type IIllegalArgumentException interface {
	IRuntimeException
	illegalArgumentException()
}

func NewIllegalArgumentException(msg string) *IllegalArgumentException {
	e := &IllegalArgumentException{}
	initThrowable(e, "java.lang.IllegalArgumentException", msg, true, nil)
	return e
}

// IllegalStateException
type IllegalStateException struct{ RuntimeException }

func (*IllegalStateException) illegalStateException() {}

type IIllegalStateException interface {
	IRuntimeException
	illegalStateException()
}

func NewIllegalStateException(msg string) *IllegalStateException {
	e := &IllegalStateException{}
	initThrowable(e, "java.lang.IllegalStateException", msg, true, nil)
	return e
}

// NullPointerException
type NullPointerException struct{ RuntimeException }

func (*NullPointerException) nullPointerException() {}

type INullPointerException interface {
	IRuntimeException
	nullPointerException()
}

//go:noinline
func NewNullPointerException() *NullPointerException {
	e := &NullPointerException{}
	initThrowable(e, "java.lang.NullPointerException", "", false, nil)
	return e
}

func NewNullPointerExceptionMsg(msg string) *NullPointerException {
	e := &NullPointerException{}
	initThrowable(e, "java.lang.NullPointerException", msg, true, nil)
	return e
}

// IndexOutOfBoundsException
type IndexOutOfBoundsException struct{ RuntimeException }

func (*IndexOutOfBoundsException) indexOutOfBoundsException() {}

type IIndexOutOfBoundsException interface {
	IRuntimeException
	indexOutOfBoundsException()
}

func NewIndexOutOfBoundsException(msg string) *IndexOutOfBoundsException {
	e := &IndexOutOfBoundsException{}
	initThrowable(e, "java.lang.IndexOutOfBoundsException", msg, true, nil)
	return e
}

func NewIndexOutOfBoundsExceptionForIndex(idx, size int) *IndexOutOfBoundsException {
	return NewIndexOutOfBoundsException(fmt.Sprintf("Index %d out of bounds for length %d", idx, size))
}

// ArrayIndexOutOfBoundsException
type ArrayIndexOutOfBoundsException struct{ IndexOutOfBoundsException }

func (*ArrayIndexOutOfBoundsException) arrayIndexOutOfBoundsException() {}

type IArrayIndexOutOfBoundsException interface {
	IIndexOutOfBoundsException
	arrayIndexOutOfBoundsException()
}

func NewArrayIndexOutOfBoundsException(msg string) *ArrayIndexOutOfBoundsException {
	e := &ArrayIndexOutOfBoundsException{}
	initThrowable(e, "java.lang.ArrayIndexOutOfBoundsException", msg, true, nil)
	return e
}

func NewArrayIndexOutOfBoundsExceptionForIndex(idx, size int) *ArrayIndexOutOfBoundsException {
	return NewArrayIndexOutOfBoundsException(fmt.Sprintf("Index %d out of bounds for length %d", idx, size))
}

// StringIndexOutOfBoundsException
type StringIndexOutOfBoundsException struct{ IndexOutOfBoundsException }

func (*StringIndexOutOfBoundsException) stringIndexOutOfBoundsException() {}

type IStringIndexOutOfBoundsException interface {
	IIndexOutOfBoundsException
	stringIndexOutOfBoundsException()
}

func NewStringIndexOutOfBoundsException(msg string) *StringIndexOutOfBoundsException {
	e := &StringIndexOutOfBoundsException{}
	initThrowable(e, "java.lang.StringIndexOutOfBoundsException", msg, true, nil)
	return e
}

func NewStringIndexOutOfBoundsExceptionForIndex(index, length int) *StringIndexOutOfBoundsException {
	return NewStringIndexOutOfBoundsException(fmt.Sprintf("String index out of range: %d", index))
}

func NewStringIndexOutOfBoundsExceptionForRange(begin, end, length int) *StringIndexOutOfBoundsException {
	return NewStringIndexOutOfBoundsException(fmt.Sprintf("String index out of range: begin %d, end %d, length %d", begin, end, length))
}

// UnsupportedOperationException
type UnsupportedOperationException struct{ RuntimeException }

func (*UnsupportedOperationException) unsupportedOperationException() {}

type IUnsupportedOperationException interface {
	IRuntimeException
	unsupportedOperationException()
}

func NewUnsupportedOperationException(msg string) *UnsupportedOperationException {
	e := &UnsupportedOperationException{}
	initThrowable(e, "java.lang.UnsupportedOperationException", msg, true, nil)
	return e
}

// ClassCastException
type ClassCastException struct{ RuntimeException }

func (*ClassCastException) classCastException() {}

type IClassCastException interface {
	IRuntimeException
	classCastException()
}

func NewClassCastException(msg string) *ClassCastException {
	e := &ClassCastException{}
	initThrowable(e, "java.lang.ClassCastException", msg, true, nil)
	return e
}

// NoSuchElementException
type NoSuchElementException struct{ RuntimeException }

func (*NoSuchElementException) noSuchElementException() {}

type INoSuchElementException interface {
	IRuntimeException
	noSuchElementException()
}

func NewNoSuchElementException(msg string) *NoSuchElementException {
	e := &NoSuchElementException{}
	initThrowable(e, "java.util.NoSuchElementException", msg, true, nil)
	return e
}

// ArithmeticException
type ArithmeticException struct{ RuntimeException }

func (*ArithmeticException) arithmeticException() {}

type IArithmeticException interface {
	IRuntimeException
	arithmeticException()
}

func NewArithmeticException(msg string) *ArithmeticException {
	e := &ArithmeticException{}
	initThrowable(e, "java.lang.ArithmeticException", msg, true, nil)
	return e
}

// NegativeArraySizeException
type NegativeArraySizeException struct{ RuntimeException }

func (*NegativeArraySizeException) negativeArraySizeException() {}

type INegativeArraySizeException interface {
	IRuntimeException
	negativeArraySizeException()
}

func NewNegativeArraySizeException(msg string) *NegativeArraySizeException {
	e := &NegativeArraySizeException{}
	initThrowable(e, "java.lang.NegativeArraySizeException", msg, true, nil)
	return e
}

// NumberFormatException
type NumberFormatException struct{ IllegalArgumentException }

func (*NumberFormatException) numberFormatException() {}

type INumberFormatException interface {
	IIllegalArgumentException
	numberFormatException()
}

func NewNumberFormatException(msg string) *NumberFormatException {
	e := &NumberFormatException{}
	initThrowable(e, "java.lang.NumberFormatException", msg, true, nil)
	return e
}

// JavaLangError corresponds to java.lang.Error.
type JavaLangError struct{ ThrowableBase }

func (*JavaLangError) error_() {}

type IError interface {
	Throwable
	error_()
}

func NewError(msg string) *JavaLangError {
	e := &JavaLangError{}
	initThrowable(e, "java.lang.Error", msg, true, nil)
	return e
}

// AssertionError
type AssertionError struct{ JavaLangError }

func (*AssertionError) assertionError() {}

type IAssertionError interface {
	IError
	assertionError()
}

func NewAssertionError(detail any) *AssertionError {
	e := &AssertionError{}
	msg := ""
	if detail != nil {
		msg = fmt.Sprint(detail)
	}
	initThrowable(e, "java.lang.AssertionError", msg, detail != nil, nil)
	return e
}

// RecoverException converts a recovered panic value into a Throwable.
// Fail-fast UNIMPLEMENTED_EXT panics bypass Java catch blocks so missing bindings
// always surface immediately with a stack trace.
func RecoverException(r any) Throwable {
	if r == nil {
		return nil
	}
	if s, ok := r.(string); ok && strings.HasPrefix(s, "UNIMPLEMENTED_EXT") {
		panic(r)
	}
	if t, ok := r.(Throwable); ok {
		if os.Getenv("JSGO_TRACE_THROWABLE") != "" {
			fmt.Fprintf(os.Stderr, "[JAVA_THROW] %s\n%s", t.Error(), debug.Stack())
		}
		return t
	}
	if os.Getenv("JSGO_STACKTRACE") != "" {
		fmt.Fprintf(os.Stderr, "[GO_PANIC] %v\n%s", r, debug.Stack())
	}
	if err, ok := r.(error); ok {
		return NewRuntimeExceptionMsg(err.Error())
	}
	return NewRuntimeExceptionMsg(fmt.Sprint(r))
}
