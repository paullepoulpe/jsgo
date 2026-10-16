package jutil

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Go-only: language-adapter — Java exceptions as Go panics.
//
// Conventions (see PORTING.md):
//   - `throw new X(...)`  →  panic(NewX(...))
//   - `try { A } catch (X e) { B }`  →  jutil.Try(func() { A }, func(e IX) { B })
//     where IX is the marker interface for X (so subclasses are caught too).
//   - `finally { C }`  →  defer func() { C }() inside the enclosing closure.
//
// Every exception type embeds ThrowableBase (directly or through a parent
// exception struct) so that marker methods are promoted and Java's
// subclass-catch semantics hold.

// Throwable is java.lang.Throwable.
type Throwable interface {
	error
	GetMessage() string
	GetCause() Throwable
	throwable() *ThrowableBase
}

// ThrowableBase holds the state shared by all exceptions.
type ThrowableBase struct {
	className string
	message   string
	hasMsg    bool
	cause     Throwable
}

func (t *ThrowableBase) throwable() *ThrowableBase { return t }

// GetMessage is Throwable.getMessage(); "" stands in for Java null.
func (t *ThrowableBase) GetMessage() string { return t.message }

// GetCause is Throwable.getCause().
func (t *ThrowableBase) GetCause() Throwable { return t.cause }

// InitCause is Throwable.initCause().
func (t *ThrowableBase) InitCause(cause Throwable) { t.cause = cause }

// Error renders like Java's Throwable.toString(): "<class>: <message>".
func (t *ThrowableBase) Error() string {
	if !t.hasMsg {
		return t.className
	}
	return t.className + ": " + t.message
}

// InitThrowable initializes the base of an exception. className is the
// fully qualified Java class name, used by Error()/toString().
func InitThrowable(t Throwable, className string, message *string, cause Throwable) {
	b := t.throwable()
	b.className = className
	if message != nil {
		b.message, b.hasMsg = *message, true
	} else if cause != nil {
		// Java: Throwable(Throwable cause) uses cause.toString() as message.
		b.message, b.hasMsg = cause.Error(), true
	}
	b.cause = cause
}

func msgPtr(s string) *string { return &s }

// RuntimeException is java.lang.RuntimeException.
type RuntimeException struct{ ThrowableBase }

func (*RuntimeException) runtimeException() {}

// IRuntimeException matches RuntimeException and all its subclasses.
type IRuntimeException interface {
	Throwable
	runtimeException()
}

func NewRuntimeException(message string) *RuntimeException {
	e := &RuntimeException{}
	InitThrowable(e, "java.lang.RuntimeException", msgPtr(message), nil)
	return e
}

func NewRuntimeExceptionWithCause(cause Throwable) *RuntimeException {
	e := &RuntimeException{}
	InitThrowable(e, "java.lang.RuntimeException", nil, cause)
	return e
}

// IllegalStateException is java.lang.IllegalStateException.
type IllegalStateException struct{ RuntimeException }

func (*IllegalStateException) illegalStateException() {}

type IIllegalStateException interface {
	IRuntimeException
	illegalStateException()
}

func NewIllegalStateException() *IllegalStateException {
	e := &IllegalStateException{}
	InitThrowable(e, "java.lang.IllegalStateException", nil, nil)
	return e
}

func NewIllegalStateExceptionMsg(message string) *IllegalStateException {
	e := &IllegalStateException{}
	InitThrowable(e, "java.lang.IllegalStateException", msgPtr(message), nil)
	return e
}

// IllegalArgumentException is java.lang.IllegalArgumentException.
type IllegalArgumentException struct{ RuntimeException }

func (*IllegalArgumentException) illegalArgumentException() {}

type IIllegalArgumentException interface {
	IRuntimeException
	illegalArgumentException()
}

func NewIllegalArgumentException() *IllegalArgumentException {
	e := &IllegalArgumentException{}
	InitThrowable(e, "java.lang.IllegalArgumentException", nil, nil)
	return e
}

func NewIllegalArgumentExceptionMsg(message string) *IllegalArgumentException {
	e := &IllegalArgumentException{}
	InitThrowable(e, "java.lang.IllegalArgumentException", msgPtr(message), nil)
	return e
}

// NumberFormatException is java.lang.NumberFormatException.
type NumberFormatException struct{ IllegalArgumentException }

func (*NumberFormatException) numberFormatException() {}

type INumberFormatException interface {
	IIllegalArgumentException
	numberFormatException()
}

func NewNumberFormatException(message string) *NumberFormatException {
	e := &NumberFormatException{}
	InitThrowable(e, "java.lang.NumberFormatException", msgPtr(message), nil)
	return e
}

// NullPointerException is java.lang.NullPointerException.
type NullPointerException struct{ RuntimeException }

func (*NullPointerException) nullPointerException() {}

type INullPointerException interface {
	IRuntimeException
	nullPointerException()
}

func NewNullPointerException() *NullPointerException {
	e := &NullPointerException{}
	InitThrowable(e, "java.lang.NullPointerException", nil, nil)
	return e
}

func NewNullPointerExceptionMsg(message string) *NullPointerException {
	e := &NullPointerException{}
	InitThrowable(e, "java.lang.NullPointerException", msgPtr(message), nil)
	return e
}

// UnsupportedOperationException is java.lang.UnsupportedOperationException.
type UnsupportedOperationException struct{ RuntimeException }

func (*UnsupportedOperationException) unsupportedOperationException() {}

type IUnsupportedOperationException interface {
	IRuntimeException
	unsupportedOperationException()
}

func NewUnsupportedOperationException() *UnsupportedOperationException {
	e := &UnsupportedOperationException{}
	InitThrowable(e, "java.lang.UnsupportedOperationException", nil, nil)
	return e
}

func NewUnsupportedOperationExceptionMsg(message string) *UnsupportedOperationException {
	e := &UnsupportedOperationException{}
	InitThrowable(e, "java.lang.UnsupportedOperationException", msgPtr(message), nil)
	return e
}

// IndexOutOfBoundsException is java.lang.IndexOutOfBoundsException.
type IndexOutOfBoundsException struct{ RuntimeException }

func (*IndexOutOfBoundsException) indexOutOfBoundsException() {}

type IIndexOutOfBoundsException interface {
	IRuntimeException
	indexOutOfBoundsException()
}

func NewIndexOutOfBoundsException(message string) *IndexOutOfBoundsException {
	e := &IndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.IndexOutOfBoundsException", msgPtr(message), nil)
	return e
}

// StringIndexOutOfBoundsException is java.lang.StringIndexOutOfBoundsException.
type StringIndexOutOfBoundsException struct{ IndexOutOfBoundsException }

func NewStringIndexOutOfBoundsException(begin, end int) *StringIndexOutOfBoundsException {
	e := &StringIndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.StringIndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("begin %d, end %d", begin, end)), nil)
	return e
}

// AsThrowable converts a recovered panic value into a Throwable. Go runtime
// errors (nil dereference, index out of range, ...) correspond to Java's
// unchecked NullPointerException / IndexOutOfBoundsException and are wrapped
// as RuntimeExceptions so that `catch (RuntimeException e)` sites see them.
func AsThrowable(r any) Throwable {
	switch v := r.(type) {
	case Throwable:
		return v
	case runtime.Error:
		return runtimeErrorAsThrowable(v)
	case error:
		e := &RuntimeException{}
		InitThrowable(e, "java.lang.RuntimeException", msgPtr(v.Error()), nil)
		return e
	default:
		e := &RuntimeException{}
		InitThrowable(e, "java.lang.RuntimeException", msgPtr(fmt.Sprint(v)), nil)
		return e
	}
}

var goIndexOutOfRange = regexp.MustCompile(`index out of range \[(-?\d+)\] with length (\d+)`)

// runtimeErrorAsThrowable maps a Go runtime fault to the Java exception the
// same faulty operation throws on the JVM, so that `catch
// (NullPointerException e)` / `catch (IndexOutOfBoundsException e)` sites
// behave as in Java. Messages follow the JDK where the JDK's message is fixed;
// the JVM's "helpful" NullPointerException messages describe Java bytecode
// and cannot be reproduced, so a nil dereference carries no message.
func runtimeErrorAsThrowable(v runtime.Error) Throwable {
	msg := v.Error()
	switch {
	case strings.Contains(msg, "nil pointer dereference"):
		return NewNullPointerException()
	case strings.Contains(msg, "index out of range"):
		if m := goIndexOutOfRange.FindStringSubmatch(msg); m != nil {
			i, _ := strconv.Atoi(m[1])
			n, _ := strconv.Atoi(m[2])
			return NewArrayIndexOutOfBoundsExceptionForLength(i, n)
		}
		e := &ArrayIndexOutOfBoundsException{}
		InitThrowable(e, "java.lang.ArrayIndexOutOfBoundsException", msgPtr(msg), nil)
		return e
	case strings.Contains(msg, "slice bounds out of range"):
		e := &ArrayIndexOutOfBoundsException{}
		InitThrowable(e, "java.lang.ArrayIndexOutOfBoundsException", msgPtr(msg), nil)
		return e
	case strings.Contains(msg, "integer divide by zero"):
		return NewArithmeticException("/ by zero")
	case strings.Contains(msg, "interface conversion"):
		e := &ClassCastException{}
		InitThrowable(e, "java.lang.ClassCastException", msgPtr(msg), nil)
		return e
	}
	e := &RuntimeException{}
	InitThrowable(e, "java.lang.RuntimeException", msgPtr(msg), nil)
	return e
}

// Try runs body and, if it throws an exception assignable to E, runs handler
// with it. Other panics propagate unchanged. E is normally one of the I*
// marker interfaces (e.g. jutil.IRuntimeException) or a concrete exception
// pointer type.
func Try[E any](body func(), handler func(E)) {
	caught, ok := tryCatch[E](body)
	if ok {
		handler(caught)
	}
}

func tryCatch[E any](body func()) (caught E, ok bool) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if _, uncatchable := r.(Uncatchable); uncatchable {
			panic(r)
		}
		if e, match := r.(E); match {
			caught, ok = e, true
			return
		}
		if _, isThrowable := r.(Throwable); !isThrowable {
			if e, match := any(AsThrowable(r)).(E); match {
				caught, ok = e, true
				return
			}
		}
		panic(r)
	}()
	body()
	return caught, false
}

// Uncatchable marks a panic value that no Java catch block catches: Try
// re-panics it whatever its handler's type.
//
// Go-only: language-adapter — the declarations generated for not-yet-ported
// Java members (internal/jscomp/pending_gen.go) panic with one, so that
// ported catch blocks do not turn "not ported yet" into a Java exception.
type Uncatchable interface {
	error
	Uncatchable()
}

// Throw panics with t. It exists so that call sites read like Java and so
// that Go vet sees a terminating statement after it via the panic.
func Throw(t Throwable) {
	panic(t)
}

// ArrayIndexOutOfBoundsException is java.lang.ArrayIndexOutOfBoundsException.
type ArrayIndexOutOfBoundsException struct{ IndexOutOfBoundsException }

func (*ArrayIndexOutOfBoundsException) arrayIndexOutOfBoundsException() {}

type IArrayIndexOutOfBoundsException interface {
	IIndexOutOfBoundsException
	arrayIndexOutOfBoundsException()
}

// NewArrayIndexOutOfBoundsException is new ArrayIndexOutOfBoundsException(int index).
func NewArrayIndexOutOfBoundsException(index int) *ArrayIndexOutOfBoundsException {
	e := &ArrayIndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.ArrayIndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("Array index out of range: %d", index)), nil)
	return e
}

// NewArrayIndexOutOfBoundsExceptionForLength is the exception the JVM throws
// for an out-of-range array access ("Index i out of bounds for length n").
func NewArrayIndexOutOfBoundsExceptionForLength(index, length int) *ArrayIndexOutOfBoundsException {
	e := &ArrayIndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.ArrayIndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("Index %d out of bounds for length %d", index, length)), nil)
	return e
}

// ClassCastException is java.lang.ClassCastException.
type ClassCastException struct{ RuntimeException }

func (*ClassCastException) classCastException() {}

type IClassCastException interface {
	IRuntimeException
	classCastException()
}

// NewClassCastException is the exception the JVM throws for a failed
// reference cast `(To) obj`; the message names both classes like the JVM's.
func NewClassCastException(from, to string) *ClassCastException {
	e := &ClassCastException{}
	InitThrowable(e, "java.lang.ClassCastException",
		msgPtr("class "+from+" cannot be cast to class "+to), nil)
	return e
}

// NoSuchElementException is java.util.NoSuchElementException.
type NoSuchElementException struct{ RuntimeException }

func (*NoSuchElementException) noSuchElementException() {}

type INoSuchElementException interface {
	IRuntimeException
	noSuchElementException()
}

func NewNoSuchElementException() *NoSuchElementException {
	e := &NoSuchElementException{}
	InitThrowable(e, "java.util.NoSuchElementException", nil, nil)
	return e
}

func NewNoSuchElementExceptionMsg(message string) *NoSuchElementException {
	e := &NoSuchElementException{}
	InitThrowable(e, "java.util.NoSuchElementException", msgPtr(message), nil)
	return e
}

// ArithmeticException is java.lang.ArithmeticException.
type ArithmeticException struct{ RuntimeException }

func (*ArithmeticException) arithmeticException() {}

type IArithmeticException interface {
	IRuntimeException
	arithmeticException()
}

func NewArithmeticException(message string) *ArithmeticException {
	e := &ArithmeticException{}
	InitThrowable(e, "java.lang.ArithmeticException", msgPtr(message), nil)
	return e
}

// JavaException is java.lang.Exception, as the base of checked exceptions
// (named like JavaLangError so that it reads unambiguously next to Go's
// error). jutil's RuntimeException does not embed it: a Java `catch
// (Exception e)` that must also see RuntimeExceptions catches Throwable.
type JavaException struct{ ThrowableBase }

func (*JavaException) javaException() {}

// IException matches the checked exceptions (JavaException subclasses).
type IException interface {
	Throwable
	javaException()
}

// IOException is java.io.IOException.
type IOException struct{ JavaException }

func (*IOException) ioException() {}

type IIOException interface {
	IException
	ioException()
}

func NewIOException(message string) *IOException {
	e := &IOException{}
	InitThrowable(e, "java.io.IOException", msgPtr(message), nil)
	return e
}

// NewIOExceptionWithCause is new IOException(Throwable cause).
func NewIOExceptionWithCause(cause Throwable) *IOException {
	e := &IOException{}
	InitThrowable(e, "java.io.IOException", nil, cause)
	return e
}
