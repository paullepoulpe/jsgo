// Port of com/google/javascript/jscomp/JSError.java
// @ closure-compiler v20260415
//
// Go maps Java's JSError record to the existing Diagnostic struct plus this
// file's factory functions and JSError type alias. Diagnostic carries the
// essential (path, line, col, code, level, message) tuple that the harness
// and CLI consume.

package diag

// JSError mirrors Java's JSError record. In Go it is a type alias for
// Diagnostic so existing code that builds diagnostics is unaffected, while
// new code can use the Java-shaped factory API below.
type JSError = Diagnostic

// MakeJSError creates a JSError (Diagnostic) from a DiagnosticType and
// formatted arguments. Mirrors Java's JSError.make(DiagnosticType, String...).
func MakeJSError(dt *DiagnosticType, arguments ...string) JSError {
	msg := dt.FormatMessage(arguments...)
	lvl := levelFromCheckLevel(dt.Level)
	return JSError{
		Code:    dt.Key,
		Level:   lvl,
		Message: msg,
	}
}

// MakeJSErrorAt creates a JSError at a specific source location.
// Mirrors Java's JSError.make(String, int, int, DiagnosticType, String...).
func MakeJSErrorAt(sourceName string, lineno, charno int, dt *DiagnosticType, arguments ...string) JSError {
	msg := dt.FormatMessage(arguments...)
	lvl := levelFromCheckLevel(dt.Level)
	return JSError{
		Path:    sourceName,
		Line:    lineno,
		Col:     charno,
		Code:    dt.Key,
		Level:   lvl,
		Message: msg,
	}
}

// levelFromCheckLevel converts a CheckLevel to the Level enum used by Diagnostic.
func levelFromCheckLevel(cl CheckLevel) Level {
	if cl == CheckLevelError {
		return Error
	}
	return Warning
}

// JSErrorBuilder is a mutable builder for JSError. Mirrors Java's
// JSError.Builder inner class.
type JSErrorBuilder struct {
	dt         *DiagnosticType
	args       []string
	sourceName string
	lineno     int
	charno     int
	level      CheckLevel
}

// BuilderFor returns a new JSErrorBuilder for the given DiagnosticType.
// Mirrors Java's JSError.builder(DiagnosticType, String...).
func BuilderFor(dt *DiagnosticType, args ...string) *JSErrorBuilder {
	return &JSErrorBuilder{
		dt:    dt,
		args:  args,
		level: dt.Level,
	}
}

// SetSourceLocation sets the file/line/col. Mirrors Java's
// JSError.Builder.setSourceLocation.
func (b *JSErrorBuilder) SetSourceLocation(sourceName string, lineno, charno int) *JSErrorBuilder {
	b.sourceName = sourceName
	b.lineno = lineno
	b.charno = charno
	return b
}

// SetLevel overrides the default level. Mirrors Java's
// JSError.Builder.setLevel.
func (b *JSErrorBuilder) SetLevel(level CheckLevel) *JSErrorBuilder {
	b.level = level
	return b
}

// Build constructs the JSError.
func (b *JSErrorBuilder) Build() JSError {
	msg := b.dt.FormatMessage(b.args...)
	return JSError{
		Path:    b.sourceName,
		Line:    b.lineno,
		Col:     b.charno,
		Code:    b.dt.Key,
		Level:   levelFromCheckLevel(b.level),
		Message: msg,
	}
}
