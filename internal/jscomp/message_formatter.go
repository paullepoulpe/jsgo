// Port of com/google/javascript/jscomp/MessageFormatter.java @ closure-compiler v20260415

package jscomp

// MessageFormatter formats warnings and errors. This interface may be used
// by implementations of ErrorManager to request message formatting
// capabilities.
type MessageFormatter interface {
	// FormatError formats an error.
	FormatError(error *JSError) string

	// FormatWarning formats a warning.
	FormatWarning(warning *JSError) string
}
