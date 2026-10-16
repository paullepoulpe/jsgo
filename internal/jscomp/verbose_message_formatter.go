// Port of com/google/javascript/jscomp/VerboseMessageFormatter.java @ closure-compiler v20260415

package jscomp

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// VerboseMessageFormatter is a verbose message formatter. This formatter
// generates very loud and long messages with multi-line source excerpts.
type VerboseMessageFormatter struct {
	AbstractMessageFormatterBase
}

var _ AbstractMessageFormatter = (*VerboseMessageFormatter)(nil)

// Java: VerboseMessageFormatter(SourceExcerptProvider)
func NewVerboseMessageFormatter(source SourceExcerptProvider) *VerboseMessageFormatter {
	f := &VerboseMessageFormatter{}
	f.InitAbstractMessageFormatterBase(source)
	return f
}

// Java: VerboseMessageFormatter.formatError(JSError)
func (f *VerboseMessageFormatter) FormatError(error *JSError) string {
	return f.GetLevelName(CheckLevel_ERROR) + ": " + f.Format(error)
}

// Java: VerboseMessageFormatter.formatWarning(JSError)
func (f *VerboseMessageFormatter) FormatWarning(warning *JSError) string {
	return f.GetLevelName(CheckLevel_WARNING) + ": " + f.Format(warning)
}

// Java: VerboseMessageFormatter.format(JSError)
func (f *VerboseMessageFormatter) Format(message *JSError) string {
	description := message.Description()
	sourceName := message.SourceName()
	lineNumber := message.GetLineNumber()
	sourceRegion := f.GetSource().GetSourceRegion(sourceName, lineNumber)
	var lineSource *string
	if sourceRegion != nil {
		excerpt := sourceRegion.GetSourceExcerpt()
		lineSource = &excerpt
	}
	source := "(unknown source)"
	if sourceName != nil && *sourceName != "" {
		source = *sourceName
	}
	line := "(unknown line)"
	if lineNumber < 0 {
		line = strconv.Itoa(lineNumber)
	}
	excerpt := "."
	if lineSource != nil {
		excerpt = ":\n\n" + *lineSource
	}
	return jutil.StringFormat(
		"%s at %s line %s %s",
		description,
		source,
		line,
		excerpt)
}
