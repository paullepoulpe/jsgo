package jutil

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"regexp/syntax"
	"runtime"
	"strconv"
	"strings"
	"syscall"
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

// NewIllegalArgumentExceptionWithMessageAndCause is new
// IllegalArgumentException(String message, Throwable cause).
func NewIllegalArgumentExceptionWithMessageAndCause(message string, cause Throwable) *IllegalArgumentException {
	e := &IllegalArgumentException{}
	InitThrowable(e, "java.lang.IllegalArgumentException", &message, cause)
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

// NewStringIndexOutOfBoundsExceptionForIndex is the exception thrown by
// String.charAt(int): "Index 3 out of bounds for length 3".
func NewStringIndexOutOfBoundsExceptionForIndex(index, length int) *StringIndexOutOfBoundsException {
	e := &StringIndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.StringIndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("Index %d out of bounds for length %d", index, length)), nil)
	return e
}

// NewStringIndexOutOfBoundsExceptionForRange is the exception thrown by
// String.substring(int, int) and String.substring(int):
// "Range [2, 1) out of bounds for length 3".
func NewStringIndexOutOfBoundsExceptionForRange(begin, end, length int) *StringIndexOutOfBoundsException {
	e := &StringIndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.StringIndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("Range [%d, %d) out of bounds for length %d", begin, end, length)), nil)
	return e
}

// NewIndexOutOfBoundsExceptionNoMessage is new IndexOutOfBoundsException().
func NewIndexOutOfBoundsExceptionNoMessage() *IndexOutOfBoundsException {
	e := &IndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.IndexOutOfBoundsException", nil, nil)
	return e
}

// NewIndexOutOfBoundsExceptionForRange is the exception thrown by
// Preconditions.checkFromToIndex with the IOOBE formatter:
// "Range [2, 5) out of bounds for length 3".
func NewIndexOutOfBoundsExceptionForRange(from, to, length int) *IndexOutOfBoundsException {
	e := &IndexOutOfBoundsException{}
	InitThrowable(e, "java.lang.IndexOutOfBoundsException",
		msgPtr(fmt.Sprintf("Range [%d, %d) out of bounds for length %d", from, to, length)), nil)
	return e
}

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

// JavaLangError is java.lang.Error (named so that embedding it does not
// shadow the Go error interface's Error method).
type JavaLangError struct{ ThrowableBase }

func (*JavaLangError) javaLangError() {}

// IError matches java.lang.Error and all its subclasses.
type IError interface {
	Throwable
	javaLangError()
}

// NewJavaLangError is new Error(String message).
func NewJavaLangError(message string) *JavaLangError {
	e := &JavaLangError{}
	InitThrowable(e, "java.lang.Error", msgPtr(message), nil)
	return e
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

// ConcurrentModificationException is
// java.util.ConcurrentModificationException: thrown by fail-fast iterators
// when a collection is structurally modified other than through the iterator
// itself.
type ConcurrentModificationException struct{ RuntimeException }

func (*ConcurrentModificationException) concurrentModificationException() {}

type IConcurrentModificationException interface {
	IRuntimeException
	concurrentModificationException()
}

// NewConcurrentModificationException is new
// ConcurrentModificationException().
func NewConcurrentModificationException() *ConcurrentModificationException {
	e := &ConcurrentModificationException{}
	InitThrowable(e, "java.util.ConcurrentModificationException", nil, nil)
	return e
}

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

// InterruptedException is java.lang.InterruptedException (a checked
// exception).
type InterruptedException struct{ JavaException }

func (*InterruptedException) interruptedException() {}

// IInterruptedException matches InterruptedException.
type IInterruptedException interface {
	IException
	interruptedException()
}

// NewInterruptedException is new InterruptedException().
func NewInterruptedException() *InterruptedException {
	e := &InterruptedException{}
	InitThrowable(e, "java.lang.InterruptedException", nil, nil)
	return e
}

// GetMessageOrNull is Throwable.getMessage() with Java null as nil:
// GetMessage returns "" for both a null and an empty message, which is not
// enough where Java concatenates the message ("null" + ...).
func GetMessageOrNull(t Throwable) *string {
	msg := t.GetMessage()
	if msg == "" && !t.throwable().hasMsg {
		return nil
	}
	return &msg
}

// Throwable_PrintStackTrace is Throwable.printStackTrace(PrintStream): the
// throwable's toString(), then "Caused by: " and the toString() of each
// cause in turn.
//
// Go-only: language-adapter — a Go panic carries no Java stack frames, so
// the "\tat ..." and "\t... n more" lines that Java prints after each of
// these lines are absent.
func Throwable_PrintStackTrace(t Throwable, s *PrintStream) {
	s.Println(t.Error())
	seen := map[Throwable]bool{t: true}
	for cause := t.GetCause(); cause != nil; cause = cause.GetCause() {
		if seen[cause] {
			s.Println("\t[CIRCULAR REFERENCE: " + cause.Error() + "]")
			break
		}
		seen[cause] = true
		s.Println("Caused by: " + cause.Error())
	}
}

// The JVM's "helpful" NullPointerException messages (JEP 358), written from
// the JEP and checked black-box against the pinned compiler's output on
// JDK 21 (tools/clie2e).
//
// When the JVM itself throws a NullPointerException (a method call, a field
// or array access on null), getMessage() describes the bytecode that failed,
// for example `Cannot invoke "com.google.javascript.rhino.Node.getSecondChild()"
// because "clazz" is null`. With -XX:-ShowCodeDetailsInExceptionMessages it
// is null instead. The option is on by default since JDK 15, so
// `java -jar closure-compiler.jar` prints the message wherever one reaches
// its output, as in an INTERNAL COMPILER ERROR report.
//
// A Go nil dereference carries no message (AsThrowable): the JVM derives
// it from bytecode. So a port spells the message out, with
// NewHelpfulNullPointerException, at each site where Java's output shows
// it. The oracle dumps made with the option off (or that replace these
// messages by "null") turn it off with SetShowCodeDetailsInExceptionMessages.

// showCodeDetailsInExceptionMessages is -XX:[+-]ShowCodeDetailsInExceptionMessages.
var showCodeDetailsInExceptionMessages = true

// SetShowCodeDetailsInExceptionMessages sets what
// -XX:+ShowCodeDetailsInExceptionMessages (on) or
// -XX:-ShowCodeDetailsInExceptionMessages (off) sets on the JVM, and returns
// the previous setting. It is on by default, as on JDK 21.
func SetShowCodeDetailsInExceptionMessages(on bool) (was bool) {
	was = showCodeDetailsInExceptionMessages
	showCodeDetailsInExceptionMessages = on
	return was
}

// NewHelpfulNullPointerException is the NullPointerException that the JVM
// throws on a null dereference whose helpful message is detail: its message
// is detail, or null when the option is off.
func NewHelpfulNullPointerException(detail string) *NullPointerException {
	if !showCodeDetailsInExceptionMessages {
		return NewNullPointerException()
	}
	return NewNullPointerExceptionMsg(detail)
}

// Go-only: language-adapter — java.nio.charset.CharacterCodingException and
// its subclasses MalformedInputException and UnmappableCharacterException
// (all IOExceptions).

// CharacterCodingException is java.nio.charset.CharacterCodingException.
type CharacterCodingException struct{ IOException }

func (*CharacterCodingException) characterCodingException() {}

// ICharacterCodingException matches CharacterCodingException and its
// subclasses.
type ICharacterCodingException interface {
	IIOException
	characterCodingException()
}

// MalformedInputException is java.nio.charset.MalformedInputException.
type MalformedInputException struct {
	CharacterCodingException
	inputLength int
}

// NewMalformedInputException is new MalformedInputException(inputLength);
// its message is "Input length = <inputLength>".
func NewMalformedInputException(inputLength int) *MalformedInputException {
	e := &MalformedInputException{inputLength: inputLength}
	msg := "Input length = " + strconv.Itoa(inputLength)
	InitThrowable(e, "java.nio.charset.MalformedInputException", &msg, nil)
	return e
}

// GetInputLength is MalformedInputException.getInputLength().
func (e *MalformedInputException) GetInputLength() int { return e.inputLength }

// UnmappableCharacterException is
// java.nio.charset.UnmappableCharacterException.
type UnmappableCharacterException struct {
	CharacterCodingException
	inputLength int
}

// NewUnmappableCharacterException is new
// UnmappableCharacterException(inputLength); its message is
// "Input length = <inputLength>".
func NewUnmappableCharacterException(inputLength int) *UnmappableCharacterException {
	e := &UnmappableCharacterException{inputLength: inputLength}
	msg := "Input length = " + strconv.Itoa(inputLength)
	InitThrowable(e, "java.nio.charset.UnmappableCharacterException", &msg, nil)
	return e
}

// GetInputLength is UnmappableCharacterException.getInputLength().
func (e *UnmappableCharacterException) GetInputLength() int { return e.inputLength }

// FileNotFoundException is java.io.FileNotFoundException.
type FileNotFoundException struct{ IOException }

func (*FileNotFoundException) fileNotFoundException() {}

// IFileNotFoundException is the catch marker for FileNotFoundException.
type IFileNotFoundException interface {
	IIOException
	fileNotFoundException()
}

// NewFileNotFoundException is new FileNotFoundException(String).
func NewFileNotFoundException(message string) *FileNotFoundException {
	e := &FileNotFoundException{}
	InitThrowable(e, "java.io.FileNotFoundException", &message, nil)
	return e
}

// FileNotFoundExceptionFromError is the exception java.io's
// FileInputStream/FileOutputStream constructors throw when opening path
// fails: FileNotFoundException("<path> (<reason>)"), the reason being the
// C library's strerror text (e.g. "No such file or directory", "Is a
// directory").
func FileNotFoundExceptionFromError(err error, path string) *FileNotFoundException {
	reason := err.Error()
	var errno syscall.Errno
	if errors.As(err, &errno) {
		reason = errno.Error()
	} else if pe := (*fs.PathError)(nil); errors.As(err, &pe) {
		reason = pe.Err.Error()
	}
	// Go's errno strings are strerror's, lowercased in the first letter.
	if reason != "" {
		reason = strings.ToUpper(reason[:1]) + reason[1:]
	}
	return NewFileNotFoundException(path + " (" + reason + ")")
}

// ZipException is java.util.zip.ZipException: an I/O error reading a ZIP
// archive.
type ZipException struct{ IOException }

func (*ZipException) zipException() {}

// IZipException is the catch marker for ZipException.
type IZipException interface {
	IIOException
	zipException()
}

// NewZipException is new ZipException(String).
func NewZipException(message string) *ZipException {
	e := &ZipException{}
	InitThrowable(e, "java.util.zip.ZipException", &message, nil)
	return e
}

// Go-only: language-adapter — java.util.regex.PatternSyntaxException,
// written from its Javadoc ("Unchecked exception thrown to indicate a syntax
// error in a regular-expression pattern"; constructor (desc, regex, index),
// with index -1 if the index is not known; getDescription, getIndex,
// getPattern, and getMessage: "a multi-line string containing the
// description of the syntax error and its index, the erroneous
// regular-expression pattern, and a visual indication of the error index
// within the pattern"). The exact layout of that message was checked
// black-box against JDK 21 (testdata/pattern_syntax.tsv, from
// tools/oracle/java/JutilPatternSyntaxVectors.java).

// PatternSyntaxException is java.util.regex.PatternSyntaxException (an
// IllegalArgumentException).
type PatternSyntaxException struct {
	IllegalArgumentException
	desc    string
	pattern string
	index   int
}

func (*PatternSyntaxException) patternSyntaxException() {}

// IPatternSyntaxException matches PatternSyntaxException.
type IPatternSyntaxException interface {
	IIllegalArgumentException
	patternSyntaxException()
}

// NewPatternSyntaxException is new PatternSyntaxException(desc, regex,
// index); index is -1 if it is not known.
func NewPatternSyntaxException(desc string, regex string, index int) *PatternSyntaxException {
	e := &PatternSyntaxException{desc: desc, pattern: regex, index: index}
	var sb strings.Builder
	sb.WriteString(desc)
	if index >= 0 {
		sb.WriteString(" near index ")
		sb.WriteString(strconv.Itoa(index))
	}
	sb.WriteString("\n")
	sb.WriteString(regex)
	if index >= 0 && index < StringLength(regex) {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat(" ", index))
		sb.WriteString("^")
	}
	InitThrowable(e, "java.util.regex.PatternSyntaxException", msgPtr(sb.String()), nil)
	return e
}

// GetDescription is PatternSyntaxException.getDescription().
func (e *PatternSyntaxException) GetDescription() string { return e.desc }

// GetIndex is PatternSyntaxException.getIndex() (-1 if not known).
func (e *PatternSyntaxException) GetIndex() int { return e.index }

// GetPattern is PatternSyntaxException.getPattern().
func (e *PatternSyntaxException) GetPattern() string { return e.pattern }

// compileTranslatedJavaRegex compiles translated, the RE2 form of the Java
// regex regex. A syntax error that RE2 reports is a Java syntax error too
// (the translation keeps the shared constructs' syntax), so it becomes a
// PatternSyntaxException, except for RE2's own limits and its stricter
// rules on repeated repetitions (Java accepts x{2}{3}, a{1001}): those are
// constructs this adapter cannot express (UnsupportedOperationException).
// RE2 does not report Java's index or description, so the exception has
// RE2's description and index -1.
func compileTranslatedJavaRegex(regex string, translated string) *regexp.Regexp {
	re, err := regexp.Compile(translated)
	if err == nil {
		return re
	}
	if se, ok := err.(*syntax.Error); ok {
		switch se.Code {
		case syntax.ErrInvalidRepeatSize, syntax.ErrInvalidRepeatOp, syntax.ErrLarge,
			syntax.ErrNestingDepth, syntax.ErrInternalError:
			unsupportedJavaRegex(regex, "RE2: "+se.Code.String())
		}
		panic(NewPatternSyntaxException(se.Code.String(), regex, -1))
	}
	panic(err)
}
