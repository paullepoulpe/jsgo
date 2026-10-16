// Port of com/google/javascript/jscomp/ErrorFormat.java @ closure-compiler v20260415

package jscomp

// ErrorFormat lists the error formats available.
type ErrorFormat int32

const (
	ErrorFormat_SINGLELINE ErrorFormat = iota
	ErrorFormat_FULL
	ErrorFormat_MULTILINE
	ErrorFormat_SOURCELESS
)

var errorFormat_names = [...]string{"SINGLELINE", "FULL", "MULTILINE", "SOURCELESS"}

// Java: ErrorFormat.values()
func ErrorFormat_values() []ErrorFormat {
	return []ErrorFormat{ErrorFormat_SINGLELINE, ErrorFormat_FULL, ErrorFormat_MULTILINE, ErrorFormat_SOURCELESS}
}

// String is Enum.name().
func (f ErrorFormat) String() string { return errorFormat_names[f] }

// Ordinal is Enum.ordinal().
func (f ErrorFormat) Ordinal() int { return int(f) }

// ToFormatter converts to a concrete formatter.
//
// Java: ErrorFormat.toFormatter(SourceExcerptProvider, boolean)
func (f ErrorFormat) ToFormatter(source SourceExcerptProvider, colorize bool) MessageFormatter {
	switch f {
	case ErrorFormat_SINGLELINE:
		formatter := NewLightweightMessageFormatter(
			source)
		formatter.SetColorize(colorize)
		return formatter
	case ErrorFormat_FULL:
		formatter :=
			NewLightweightMessageFormatterWithFormat(source, SourceExcerptProvider_SourceExcerpt_FULL)
		formatter.SetColorize(colorize)
		return formatter
	case ErrorFormat_MULTILINE:
		formatter :=
			NewLightweightMessageFormatterWithFormat(source, SourceExcerptProvider_SourceExcerpt_REGION)
		formatter.SetColorize(colorize)
		return formatter
	default: // SOURCELESS
		formatter :=
			LightweightMessageFormatter_WithoutSource()
		formatter.SetColorize(colorize)
		return formatter
	}
}
