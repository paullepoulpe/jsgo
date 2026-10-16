// Port of com/google/debugging/sourcemap/SourceMapParseException.java @ closure-compiler v20260415

package sourcemap

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceMapParseException is thrown if an invalid or unknown source map is
// encountered. It is a checked exception (java.lang.Exception).
type SourceMapParseException struct{ jutil.JavaException }

func (*SourceMapParseException) sourceMapParseException() {}

// ISourceMapParseException matches SourceMapParseException.
type ISourceMapParseException interface {
	jutil.IException
	sourceMapParseException()
}

// NewSourceMapParseException constructs a SourceMapParseException.
//
// Java: SourceMapParseException(String)
func NewSourceMapParseException(message string) *SourceMapParseException {
	e := &SourceMapParseException{}
	jutil.InitThrowable(e, "com.google.debugging.sourcemap.SourceMapParseException", &message, nil)
	return e
}
