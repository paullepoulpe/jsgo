// Port of com/google/javascript/jscomp/AbstractMessageFormatter.java
// @ closure-compiler v20260415

package diag

// MessageFormatter formats a JSError (represented here as a Diagnostic) for
// output. Mirrors Java's MessageFormatter interface.
type MessageFormatter interface {
	FormatError(d Diagnostic) string
	FormatWarning(d Diagnostic) string
}

// SourceExcerptProvider returns a line of source text given a file and line
// number. Used by AbstractMessageFormatter to render source context.
// Mirrors Java's SourceExcerptProvider.
type SourceExcerptProvider interface {
	// GetSourceLine returns the text of the given 1-indexed line in the named
	// source file, or "" if not available.
	GetSourceLine(sourceName string, lineNumber int) string
}

// Color ANSI escape codes mirroring Java's inner Color enum.
const (
	colorError   = "[31m"
	colorWarning = "[35m"
	colorNoColor = "[39m"
	colorBold    = "[1m"
	colorUnbold  = "[0m"
)

// supportedColorTerminals mirrors Java's SUPPORTED_COLOR_TERMINALS.
var supportedColorTerminals = map[string]bool{
	"xterm":       true,
	"xterm-color": true,
	"xterm-256color": true,
	"screen-bce":  true,
}

// TermSupportsColor reports whether the given TERM value supports ANSI color.
// Mirrors Java's AbstractMessageFormatter.termSupportsColor.
func TermSupportsColor(term string) bool {
	return supportedColorTerminals[term]
}

// AbstractMessageFormatter provides shared colorization logic for message
// formatters. Mirrors Java's AbstractMessageFormatter abstract class.
type AbstractMessageFormatter struct {
	source   SourceExcerptProvider
	colorize bool
}

// NewAbstractMessageFormatter constructs an AbstractMessageFormatter backed by
// the given SourceExcerptProvider.
func NewAbstractMessageFormatter(source SourceExcerptProvider) *AbstractMessageFormatter {
	return &AbstractMessageFormatter{source: source}
}

// SetColorize enables or disables ANSI colorization in output.
func (f *AbstractMessageFormatter) SetColorize(colorize bool) {
	f.colorize = colorize
}

// GetSource returns the underlying SourceExcerptProvider.
func (f *AbstractMessageFormatter) GetSource() SourceExcerptProvider {
	return f.source
}

// GetLevelName returns a possibly-colorized level name ("ERROR" or "WARNING").
// Mirrors Java's AbstractMessageFormatter.getLevelName.
func (f *AbstractMessageFormatter) GetLevelName(level CheckLevel) string {
	switch level {
	case CheckLevelError:
		return f.maybeColorize("ERROR", colorError, colorNoColor)
	case CheckLevelWarning:
		return f.maybeColorize("WARNING", colorWarning, colorNoColor)
	default:
		return level.String()
	}
}

// MaybeEmbolden wraps text in ANSI bold codes when colorize is on.
// Mirrors Java's AbstractMessageFormatter.maybeEmbolden.
func (f *AbstractMessageFormatter) MaybeEmbolden(text string) string {
	if !f.colorize {
		return text
	}
	return colorBold + text + colorUnbold
}

func (f *AbstractMessageFormatter) maybeColorize(text, start, end string) string {
	if !f.colorize {
		return text
	}
	return start + text + end
}
