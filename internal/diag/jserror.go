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

// Make mirrors JSError.make(DiagnosticType, String...) — creates a JSError with no source info.
func Make(dt *DiagnosticType, arguments ...string) JSError {
	return BuilderFor(dt, arguments...).Build()
}

// Go-only: MakeAt is a convenience wrapper for JSError.make(String, int, int, DiagnosticType, String...);
// Java names this overload `make` (same as other overloads), but Go needs a distinct name.
func MakeAt(sourceName string, lineno, charno int, dt *DiagnosticType, arguments ...string) JSError {
	return BuilderFor(dt, arguments...).SetSourceLocation(sourceName, lineno, charno).Build()
}

// Go-only: MakeJSError is the older name for Make; kept for callers that predate the Java-parity rename.
func MakeJSError(dt *DiagnosticType, arguments ...string) JSError {
	return Make(dt, arguments...)
}

// Go-only: MakeJSErrorAt is the older name for MakeAt; kept for backward compatibility.
func MakeJSErrorAt(sourceName string, lineno, charno int, dt *DiagnosticType, arguments ...string) JSError {
	return MakeAt(sourceName, lineno, charno, dt, arguments...)
}

// Go-only: levelFromCheckLevel converts a CheckLevel to the Level enum used by Diagnostic.
func levelFromCheckLevel(cl CheckLevel) Level {
	if cl == CheckLevelError {
		return Error
	}
	return Warning
}

// GetDescription mirrors JSError.getDescription() — returns the formatted error message.
func GetDescription(e JSError) string { return e.Message }

// GetSourceName mirrors JSError.getSourceName() — returns the source file name.
func GetSourceName(e JSError) string { return e.Path }

// GetLineno mirrors JSError.getLineno() — returns the 1-indexed line number.
func GetLineno(e JSError) int { return e.Line }

// GetCharno mirrors JSError.getCharno() — returns the 0-indexed character offset.
func GetCharno(e JSError) int { return e.Col }

// GetLineNumber mirrors JSError.getLineNumber() — alias for GetLineno.
func GetLineNumber(e JSError) int { return e.Line }

// GetLength mirrors JSError.getLength() — returns the error region length in characters.
// Go-only partial: Go's Diagnostic does not track length; returns 0.
func GetLength(e JSError) int { return 0 }

// GetNodeSourceOffset mirrors JSError.getNodeSourceOffset() — returns -1 (no node in Go port).
func GetNodeSourceOffset(e JSError) int { return -1 }

// Format mirrors JSError.format(CheckLevel, MessageFormatter) — formats this error for output.
func Format(e JSError, level CheckLevel, formatter MessageFormatter) string {
	switch level {
	case CheckLevelError:
		return formatter.FormatError(e)
	case CheckLevelWarning:
		return formatter.FormatWarning(e)
	default:
		return ""
	}
}

// Go-only: JSErrorToString implements JSError.toString() as a package-level function;
// Java uses the record's toString() method but Go can't add methods to a type alias.
func JSErrorToString(e JSError) string {
	srcName := e.Path
	if srcName == "" {
		srcName = "(unknown source)"
	}
	lineno := "(unknown line)"
	if e.Line != DefaultLineno {
		lineno = fmt.Sprintf("%d", e.Line)
	}
	charno := "(unknown column)"
	if e.Col != DefaultCharno {
		charno = fmt.Sprintf("%d", e.Col)
	}
	return e.Code + ". " + e.Message + " at " + srcName + " line " + lineno + " : " + charno
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

// Builder mirrors JSError.builder(DiagnosticType, String...) — returns a new JSErrorBuilder.
func Builder(dt *DiagnosticType, args ...string) *JSErrorBuilder {
	return &JSErrorBuilder{
		dt:     dt,
		args:   args,
		level:  dt.Level,
		lineno: DefaultLineno,
		charno: DefaultCharno,
		length: DefaultLength,
	}
}

// Go-only: BuilderFor is the older name for Builder; kept for callers that predate the rename.
func BuilderFor(dt *DiagnosticType, args ...string) *JSErrorBuilder {
	return Builder(dt, args...)
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

// SetNode mirrors JSError.Builder.setNode — sets source location from a node.
// Go-only partial: Go's Node equivalent is not yet fully ported; stores source name and position.
func (b *JSErrorBuilder) SetNode(sourceName string, lineno, charno int) *JSErrorBuilder {
	b.sourceName = sourceName
	b.lineno = lineno
	b.charno = charno
	return b
}

// SetNodeRange mirrors JSError.Builder.setNodeRange — sets source location from a node range.
// Go-only partial: Go's Node equivalent is not yet fully ported.
func (b *JSErrorBuilder) SetNodeRange(sourceName string, lineno, charno int) *JSErrorBuilder {
	b.sourceName = sourceName
	b.lineno = lineno
	b.charno = charno
	return b
}

// SetRequirement mirrors JSError.Builder.setRequirement — attaches a conformance requirement.
// Go-only partial: conformance requirements are not yet ported; this is a no-op.
func (b *JSErrorBuilder) SetRequirement() *JSErrorBuilder {
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

// ---------------------------------------------------------------------------
// JSErrorMessageFormatter — a MessageFormatter implementation for JSErrorFull
// ---------------------------------------------------------------------------

// JSErrorMessageFormatter formats JSErrorFull values for display.
// Mirrors the general contract of Java's MessageFormatter interface but
// operates on JSErrorFull rather than the basic Diagnostic type.
type JSErrorMessageFormatter struct {
	// IncludeSourceCode, if true, adds a source excerpt to the output.
	IncludeSourceCode bool
}

// FormatErrorFull formats an error-level JSErrorFull. Mirrors the behavior
// of Java's LightweightMessageFormatter.formatError.
func (f *JSErrorMessageFormatter) FormatErrorFull(e JSErrorFull) string {
	return f.formatFull(e, "ERROR")
}

// FormatWarningFull formats a warning-level JSErrorFull. Mirrors the behavior
// of Java's LightweightMessageFormatter.formatWarning.
func (f *JSErrorMessageFormatter) FormatWarningFull(e JSErrorFull) string {
	return f.formatFull(e, "WARNING")
}

func (f *JSErrorMessageFormatter) formatFull(e JSErrorFull, level string) string {
	src := e.Path
	if src == "" {
		src = "(unknown source)"
	}
	return fmt.Sprintf("%s:%d:%d: %s - [%s] %s",
		src, e.Line, e.Col, level, e.Code, e.Message)
}

// FormatError implements MessageFormatter.FormatError for the basic Diagnostic type.
// Mirrors Java's LightweightMessageFormatter.formatError(JSError).
func (f *JSErrorMessageFormatter) FormatError(d Diagnostic) string {
	return fmt.Sprintf("%s:%d:%d: ERROR - [%s] %s",
		d.Path, d.Line, d.Col, d.Code, d.Message)
}

// FormatWarning implements MessageFormatter.FormatWarning for the basic Diagnostic type.
// Mirrors Java's LightweightMessageFormatter.formatWarning(JSError).
func (f *JSErrorMessageFormatter) FormatWarning(d Diagnostic) string {
	return fmt.Sprintf("%s:%d:%d: WARNING - [%s] %s",
		d.Path, d.Line, d.Col, d.Code, d.Message)
}

// ---------------------------------------------------------------------------
// JSErrorFull equality / comparison helpers
// ---------------------------------------------------------------------------

// SameLocation reports whether two JSErrorFull values refer to the same
// source location (same file, line, and column). Useful for deduplication.
func SameLocation(a, b JSErrorFull) bool {
	return a.Path == b.Path && a.Line == b.Line && a.Col == b.Col
}

// SameType reports whether two JSErrorFull values have the same diagnostic
// code. Useful for grouping errors by type.
func SameType(a, b JSErrorFull) bool {
	return a.Code == b.Code
}

// ---------------------------------------------------------------------------
// Make helpers for common JSError construction patterns
// ---------------------------------------------------------------------------

// MakeJSErrorFullAt creates a JSErrorFull at the given location.
// Convenience wrapper over NewJSErrorFullAt.
func MakeJSErrorFullAt(sourceName string, lineno, charno int, dt *DiagnosticType, args ...string) JSErrorFull {
	return NewJSErrorFullAt(sourceName, lineno, charno, dt, args...)
}

// MakeJSErrorFullWithLength creates a JSErrorFull with a specific error-region
// length. Mirrors Java's JSError.Builder.setNodeRange(start, end).
func MakeJSErrorFullWithLength(sourceName string, lineno, charno, length int, dt *DiagnosticType, args ...string) JSErrorFull {
	e := NewJSErrorFullAt(sourceName, lineno, charno, dt, args...)
	e.Length = length
	return e
}

// ---------------------------------------------------------------------------
// WarningsGuard interface (minimal) — mirrors Java's WarningsGuard
// ---------------------------------------------------------------------------

// WarningsGuard adjusts the reporting level of JSErrors before they are
// reported to the error manager. Mirrors Java's WarningsGuard abstract class.
type WarningsGuard interface {
	// Level returns the overridden CheckLevel for the given JSError, and
	// whether the override applies. If the second return value is false,
	// the error's default level is used.
	Level(d Diagnostic) (CheckLevel, bool)
	// GetPriority returns the priority of this guard. Higher numbers apply
	// later. Mirrors Java's WarningsGuard.getPriority().
	GetPriority() int
}

// CompositeWarningsGuard chains multiple WarningsGuards. Mirrors Java's
// ComposeWarningsGuard.
type CompositeWarningsGuard struct {
	guards []WarningsGuard
}

// NewCompositeWarningsGuard creates a CompositeWarningsGuard from the given
// guards, ordered by priority (lowest first). Mirrors Java's
// ComposeWarningsGuard(List<WarningsGuard>).
func NewCompositeWarningsGuard(guards ...WarningsGuard) *CompositeWarningsGuard {
	return &CompositeWarningsGuard{guards: guards}
}

// Level returns the first guard's overriding level, if any.
// Mirrors Java's ComposeWarningsGuard.level(JSError).
func (g *CompositeWarningsGuard) Level(d Diagnostic) (CheckLevel, bool) {
	for _, guard := range g.guards {
		if level, ok := guard.Level(d); ok {
			return level, true
		}
	}
	return 0, false
}

// GetPriority returns the maximum priority of all guards.
func (g *CompositeWarningsGuard) GetPriority() int {
	max := 0
	for _, guard := range g.guards {
		if p := guard.GetPriority(); p > max {
			max = p
		}
	}
	return max
}

// AddGuard appends a WarningsGuard to the composite.
func (g *CompositeWarningsGuard) AddGuard(guard WarningsGuard) {
	g.guards = append(g.guards, guard)
}
