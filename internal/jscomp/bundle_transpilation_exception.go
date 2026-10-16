// Port of com/google/javascript/jscomp/bundle/TranspilationException.java @ closure-compiler v20260415

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TranspilationException is an unchecked exception thrown when
// transpilation fails due to one or more errors in the input script.
type TranspilationException struct {
	jutil.RuntimeException
	errors   []*JSError
	warnings []*JSError
}

func (*TranspilationException) transpilationException() {}

// ITranspilationException matches TranspilationException.
type ITranspilationException interface {
	jutil.IRuntimeException
	transpilationException()
}

// NewTranspilationExceptionWithCause takes cause, an Exception.
//
// Java: TranspilationException(Exception)
func NewTranspilationExceptionWithCause(cause jutil.Throwable) *TranspilationException {
	return newTranspilationExceptionFromRoot(transpilationException_tryCastToTranspilationException(cause.GetCause()), cause)
}

// Java: TranspilationException(SourceExcerptProvider, ImmutableList<JSError>, ImmutableList<JSError>)
func NewTranspilationException(
	source SourceExcerptProvider,
	errors []*JSError,
	warnings []*JSError) *TranspilationException {
	formatted := transpilationException_format(source, errors, warnings)
	return newTranspilationException(errors, warnings, &formatted, nil)
}

// Java: TranspilationException(TranspilationException, Exception)
func newTranspilationExceptionFromRoot(root *TranspilationException, cause jutil.Throwable) *TranspilationException {
	if root == nil {
		panic(jutil.NewNullPointerException())
	}
	// root's message is the formatted report, never null.
	message := root.GetMessage()
	return newTranspilationException(root.errors, root.warnings, &message, cause)
}

// Java: TranspilationException(ImmutableList<JSError>, ImmutableList<JSError>, String, Exception)
func newTranspilationException(
	errors []*JSError,
	warnings []*JSError,
	formatted *string,
	cause jutil.Throwable) *TranspilationException {
	e := &TranspilationException{errors: errors, warnings: warnings}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.bundle.TranspilationException", formatted, cause)
	return e
}

// Java: TranspilationException.errors()
func (e *TranspilationException) Errors() []*JSError {
	return e.errors
}

// Java: TranspilationException.warnings()
func (e *TranspilationException) Warnings() []*JSError {
	return e.warnings
}

// Java: TranspilationException.format(SourceExcerptProvider, ImmutableList<JSError>, ImmutableList<JSError>)
func transpilationException_format(
	source SourceExcerptProvider,
	errors []*JSError,
	warnings []*JSError) string {
	var sb strings.Builder
	sb.WriteString("Transpilation failed:\n")
	var formatter MessageFormatter
	if source != nil {
		formatter = ErrorFormat_SINGLELINE.ToFormatter(source, false)
	} else {
		formatter = ErrorFormat_SOURCELESS.ToFormatter(source, false)
	}
	for _, error := range errors {
		sb.WriteString("\n")
		sb.WriteString(formatter.FormatError(error))
	}
	for _, warning := range warnings {
		sb.WriteString("\n")
		sb.WriteString(formatter.FormatError(warning))
	}
	return sb.String()
}

// transpilationException_tryCastToTranspilationException may return nil.
//
// Java: TranspilationException.tryCastToTranspilationException(Throwable)
func transpilationException_tryCastToTranspilationException(t jutil.Throwable) *TranspilationException {
	if transpilationException, ok := t.(*TranspilationException); ok {
		return transpilationException
	}
	return nil
}
