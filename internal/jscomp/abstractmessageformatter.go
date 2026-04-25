// Port of com/google/javascript/jscomp/AbstractMessageFormatter.java
// @ closure-compiler v20260415

package jscomp

import "strings"

// messageColor is the private Color enum from AbstractMessageFormatter.
type messageColor int

const (
	messageColorError   messageColor = iota // [31m
	messageColorWarning                     // [35m
	messageColorNoColor                     // [39m
	messageColorBold                        // [1m
	messageColorUnbold                      // [0m
)

var messageColorCodes = map[messageColor]string{
	messageColorError:   "[31m",
	messageColorWarning: "[35m",
	messageColorNoColor: "[39m",
	messageColorBold:    "[1m",
	messageColorUnbold:  "[0m",
}

func (c messageColor) controlChar() string { return messageColorCodes[c] }

// supportedColorTerminals mirrors SUPPORTED_COLOR_TERMINALS in Java.
var supportedColorTerminals = map[string]bool{
	"xterm": true, "xterm-color": true, "xterm-256color": true, "screen-bce": true,
}

// AbstractMessageFormatter is the Go port of AbstractMessageFormatter.java.
// It provides common infrastructure for message formatters that need a
// SourceExcerptProvider.
type AbstractMessageFormatter struct {
	// source is the SourceExcerptProvider (not yet fully ported; stored as
	// interface{} until SourceExcerptProvider is defined).
	source   interface{}
	colorize bool
}

// NewAbstractMessageFormatter creates an AbstractMessageFormatter with the
// given source excerpt provider.
func NewAbstractMessageFormatter(source interface{}) *AbstractMessageFormatter {
	return &AbstractMessageFormatter{source: source}
}

// SetColorize enables or disables ANSI color output.
func (f *AbstractMessageFormatter) SetColorize(colorize bool) { f.colorize = colorize }

// GetSource returns the source excerpt provider.
func (f *AbstractMessageFormatter) GetSource() interface{} { return f.source }

// TermSupportsColor reports whether term is a known color-supporting terminal.
// Mirrors AbstractMessageFormatter.termSupportsColor(String).
func TermSupportsColor(term string) bool { return supportedColorTerminals[term] }

// GetLevelName returns a (possibly colorized) name for the given CheckLevel.
// Mirrors AbstractMessageFormatter.getLevelName(CheckLevel).
func (f *AbstractMessageFormatter) GetLevelName(level CheckLevel) string {
	switch level {
	case LevelError:
		return f.maybeColorize("ERROR", messageColorError)
	case LevelWarning:
		return f.maybeColorize("WARNING", messageColorWarning)
	default:
		return strings.ToUpper(level.String())
	}
}

// MaybeEmbolden wraps text in bold ANSI escape sequences if colorize is on.
func (f *AbstractMessageFormatter) MaybeEmbolden(text string) string {
	if !f.colorize {
		return text
	}
	return messageColorBold.controlChar() + text + messageColorUnbold.controlChar()
}

func (f *AbstractMessageFormatter) maybeColorize(text string, color messageColor) string {
	if !f.colorize {
		return text
	}
	return color.controlChar() + text + messageColorNoColor.controlChar()
}

// String returns the string representation of a CheckLevel, used by
// GetLevelName's default branch.
func (l CheckLevel) String() string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelWarning:
		return "WARNING"
	default:
		return "OFF"
	}
}
