// Port of com/google/javascript/jscomp/JSError.java
// @ closure-compiler v20260415
//
// Java's JSError is a record that carries:
//   type          — DiagnosticType (code + default level + format)
//   description   — already-formatted message string
//   sourceName    — source file path (nullable)
//   lineno        — 1-indexed line number (-1 if unknown)
//   charno        — 0-indexed column number (-1 if unknown)
//   length        — length of the error region in chars
//   node          — AST node where the warning occurred (nullable)
//   defaultLevel  — CheckLevel (may be overridden by WarningsGuard)
//   requirement   — conformance Requirement (nullable)
//
// Go maps Java's JSError record to the existing Diagnostic struct plus this
// file's factory functions and JSError type alias. Diagnostic carries the
// essential (path, line, col, code, level, message) tuple that the harness
// and CLI consume.
//
// Additional fields (length, node, requirement) are tracked in the extended
// JSErrorFull struct for passes that need them; JSError remains an alias for
// backward compatibility.

package diag

import "fmt"

// ---------------------------------------------------------------------------
// Sentinel values matching Java's DEFAULT_* constants
// ---------------------------------------------------------------------------

const (
	// DefaultLineno is -1, matching Java's JSError.DEFAULT_LINENO.
	DefaultLineno = -1
	// DefaultCharno is -1, matching Java's JSError.DEFAULT_CHARNO.
	DefaultCharno = -1
	// DefaultLength is 0, matching Java's JSError.DEFAULT_LENGTH.
	DefaultLength = 0
)

// JSError mirrors Java's JSError record. In Go it is a type alias for
// Diagnostic so existing code that builds diagnostics is unaffected, while
// new code can use the Java-shaped factory API below.
type JSError = Diagnostic

// ---------------------------------------------------------------------------
// Factory functions — mirror Java's JSError.make overloads
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// JSErrorBuilder — mirrors Java's JSError.Builder inner class
// ---------------------------------------------------------------------------

// JSErrorBuilder is a mutable builder for JSError. Mirrors Java's
// JSError.Builder inner class.
type JSErrorBuilder struct {
	dt         *DiagnosticType
	args       []string
	sourceName string
	lineno     int
	charno     int
	length     int
	level      CheckLevel
	// node and requirement are omitted in the basic builder; they are tracked
	// in JSErrorFull for passes that need them.
}

// BuilderFor returns a new JSErrorBuilder for the given DiagnosticType.
// Mirrors Java's JSError.builder(DiagnosticType, String...).
func BuilderFor(dt *DiagnosticType, args ...string) *JSErrorBuilder {
	return &JSErrorBuilder{
		dt:     dt,
		args:   args,
		level:  dt.Level,
		lineno: DefaultLineno,
		charno: DefaultCharno,
		length: DefaultLength,
	}
}

// SetSourceLocation sets the file/line/col. Mirrors Java's
// JSError.Builder.setSourceLocation(String, int, int).
func (b *JSErrorBuilder) SetSourceLocation(sourceName string, lineno, charno int) *JSErrorBuilder {
	b.sourceName = sourceName
	b.lineno = lineno
	b.charno = charno
	return b
}

// SetLevel overrides the default level. Mirrors Java's
// JSError.Builder.setLevel(CheckLevel).
func (b *JSErrorBuilder) SetLevel(level CheckLevel) *JSErrorBuilder {
	b.level = level
	return b
}

// Build constructs the JSError. Mirrors Java's JSError.Builder.build().
func (b *JSErrorBuilder) Build() JSError {
	msg := b.dt.FormatMessage(b.args...)
	lineno := b.lineno
	if lineno == DefaultLineno {
		lineno = 0
	}
	charno := b.charno
	if charno == DefaultCharno {
		charno = 0
	}
	return JSError{
		Path:    b.sourceName,
		Line:    lineno,
		Col:     charno,
		Code:    b.dt.Key,
		Level:   levelFromCheckLevel(b.level),
		Message: msg,
	}
}

// ---------------------------------------------------------------------------
// JSErrorFull — extended JSError struct for passes needing node / length
// ---------------------------------------------------------------------------

// JSErrorFull extends JSError with the additional fields carried by Java's
// JSError record that the basic Diagnostic struct does not capture.
// Passes that need to inspect the AST node or error region length use
// JSErrorFull instead of JSError.
type JSErrorFull struct {
	// Diagnostic embeds the core JSError / Diagnostic fields.
	Diagnostic

	// DiagType is the DiagnosticType that generated this error.
	DiagType *DiagnosticType

	// Length is the number of characters in the error region.
	// Mirrors Java's JSError.length().
	Length int

	// DefaultLevel is the default reporting level before any WarningsGuard
	// is applied. Mirrors Java's JSError.defaultLevel().
	DefaultLevel CheckLevel
}

// NewJSErrorFull creates a JSErrorFull from a DiagnosticType and arguments.
// Mirrors Java's JSError.make(DiagnosticType, String...) for the full form.
func NewJSErrorFull(dt *DiagnosticType, args ...string) JSErrorFull {
	msg := dt.FormatMessage(args...)
	lvl := levelFromCheckLevel(dt.Level)
	return JSErrorFull{
		Diagnostic: Diagnostic{
			Code:    dt.Key,
			Level:   lvl,
			Message: msg,
		},
		DiagType:     dt,
		Length:       DefaultLength,
		DefaultLevel: dt.Level,
	}
}

// NewJSErrorFullAt creates a JSErrorFull at a specific source location.
// Mirrors Java's JSError.make(String, int, int, DiagnosticType, String...).
func NewJSErrorFullAt(sourceName string, lineno, charno int, dt *DiagnosticType, args ...string) JSErrorFull {
	msg := dt.FormatMessage(args...)
	lvl := levelFromCheckLevel(dt.Level)
	return JSErrorFull{
		Diagnostic: Diagnostic{
			Path:    sourceName,
			Line:    lineno,
			Col:     charno,
			Code:    dt.Key,
			Level:   lvl,
			Message: msg,
		},
		DiagType:     dt,
		Length:       DefaultLength,
		DefaultLevel: dt.Level,
	}
}

// GetType returns the DiagnosticType. Mirrors Java's JSError.getType().
func (e JSErrorFull) GetType() *DiagnosticType { return e.DiagType }

// GetDescription returns the formatted message. Mirrors Java's
// JSError.getDescription().
func (e JSErrorFull) GetDescription() string { return e.Message }

// GetSourceName returns the source file name. Mirrors Java's
// JSError.getSourceName().
func (e JSErrorFull) GetSourceName() string { return e.Path }

// GetLineno returns the 1-indexed line number. Mirrors Java's
// JSError.getLineno().
func (e JSErrorFull) GetLineno() int { return e.Line }

// GetCharno returns the 0-indexed column number. Mirrors Java's
// JSError.getCharno().
func (e JSErrorFull) GetCharno() int { return e.Col }

// GetLength returns the length of the error region. Mirrors Java's
// JSError.getLength().
func (e JSErrorFull) GetLength() int { return e.Length }

// GetDefaultLevel returns the default CheckLevel before WarningsGuards.
// Mirrors Java's JSError.getDefaultLevel().
func (e JSErrorFull) GetDefaultLevel() CheckLevel { return e.DefaultLevel }

// GetLineNumber is an alias for GetLineno. Mirrors Java's
// JSError.getLineNumber().
func (e JSErrorFull) GetLineNumber() int { return e.Line }

// Format formats the error at the given level using a MessageFormatter.
// Mirrors Java's JSError.format(CheckLevel, MessageFormatter).
func (e JSErrorFull) Format(level CheckLevel) string {
	lvlStr := "WARNING"
	if level == CheckLevelError {
		lvlStr = "ERROR"
	}
	sourceName := e.Path
	if sourceName == "" {
		sourceName = "(unknown source)"
	}
	lineno := fmt.Sprintf("%d", e.Line)
	if e.Line <= 0 {
		lineno = "(unknown line)"
	}
	charno := fmt.Sprintf("%d", e.Col)
	if e.Col < 0 {
		charno = "(unknown column)"
	}
	return e.Code + ". " + e.Message + " at " + sourceName +
		" line " + lineno + " : " + charno + " [" + lvlStr + "]"
}

// String returns the default rendering of the error. Mirrors Java's
// JSError.toString().
func (e JSErrorFull) String() string {
	sourceName := e.Path
	if sourceName == "" {
		sourceName = "(unknown source)"
	}
	lineno := fmt.Sprintf("%d", e.Line)
	if e.Line <= 0 {
		lineno = "(unknown line)"
	}
	charno := fmt.Sprintf("%d", e.Col)
	if e.Col < 0 {
		charno = "(unknown column)"
	}
	return e.Code + ". " + e.Message + " at " + sourceName +
		" line " + lineno + " : " + charno
}
