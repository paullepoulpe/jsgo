// Port of com/google/javascript/jscomp/PrintStreamErrorManager.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// PrintStreamErrorManager is an error manager that prints errors and
// warnings to the print stream provided in addition to the functionality of
// the BasicErrorManager.
//
// It collaborates with a SourceExcerptProvider via a MessageFormatter to
// display error messages with source context.
//
// Deprecated: Please use the SortingErrorManager with a
// PrintStreamErrorReportGenerator instead.
type PrintStreamErrorManager struct {
	BasicErrorManagerBase
	formatter          MessageFormatter
	stream             *jutil.PrintStream
	summaryDetailLevel int
}

var _ BasicErrorManager = (*PrintStreamErrorManager)(nil)

// NewPrintStreamErrorManagerWithFormatter creates an error manager.
// formatter is the message formatter used to format the messages; stream is
// the stream on which the errors and warnings should be printed (this class
// does not close the stream).
//
// Java: PrintStreamErrorManager(MessageFormatter, PrintStream)
func NewPrintStreamErrorManagerWithFormatter(formatter MessageFormatter,
	stream *jutil.PrintStream) *PrintStreamErrorManager {
	m := &PrintStreamErrorManager{summaryDetailLevel: 1}
	m.InitBasicErrorManagerBase(m)
	m.formatter = formatter
	m.stream = stream
	return m
}

// NewPrintStreamErrorManager creates an instance with a source-less error
// formatter.
//
// Java: PrintStreamErrorManager(PrintStream)
func NewPrintStreamErrorManager(stream *jutil.PrintStream) *PrintStreamErrorManager {
	return NewPrintStreamErrorManagerWithFormatter(ErrorFormat_SOURCELESS.ToFormatter(nil, false), stream)
}

// Java: PrintStreamErrorManager.println(CheckLevel, JSError)
func (m *PrintStreamErrorManager) Println(level *CheckLevel, error *JSError) {
	if level == nil {
		panic(jutil.NewNullPointerException()) // switch on a null CheckLevel in JSError.format
	}
	m.stream.PrintlnNullable(error.Format(*level, m.formatter))
}

// Java: PrintStreamErrorManager.setSummaryDetailLevel(int)
func (m *PrintStreamErrorManager) SetSummaryDetailLevel(summaryDetailLevel int) {
	m.summaryDetailLevel = summaryDetailLevel
}

// Java: PrintStreamErrorManager.printSummary()
func (m *PrintStreamErrorManager) PrintSummary() {
	if m.summaryDetailLevel >= 3 ||
		(m.summaryDetailLevel >= 1 && m.GetErrorCount()+m.GetWarningCount() > 0) ||
		(m.summaryDetailLevel >= 2 && m.GetTypedPercent() > 0.0) {
		if m.GetTypedPercent() > 0.0 {
			m.stream.Printf(
				"%d error(s), %d warning(s), %s%% typed%n",
				m.GetErrorCount(), m.GetWarningCount(), jutil.Formatter_FormatDecimalFloat(m.GetTypedPercent(), 1))
		} else {
			m.stream.Printf("%d error(s), %d warning(s)%n", m.GetErrorCount(), m.GetWarningCount())
		}
	}
}
