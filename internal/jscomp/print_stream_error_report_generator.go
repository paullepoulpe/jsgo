// Port of com/google/javascript/jscomp/PrintStreamErrorReportGenerator.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// PrintStreamErrorReportGenerator is an error report generator that prints
// errors and warnings to the print stream provided.
//
// It collaborates with a SourceExcerptProvider via a MessageFormatter to
// display error messages with source context.
type PrintStreamErrorReportGenerator struct {
	formatter          MessageFormatter
	stream             *jutil.PrintStream
	summaryDetailLevel int
}

var _ SortingErrorManager_ErrorReportGenerator = (*PrintStreamErrorReportGenerator)(nil)

// NewPrintStreamErrorReportGenerator creates an error report generator.
// formatter is the message formatter used to format the messages; stream
// is the stream on which the errors and warnings should be printed (this
// class does not close the stream).
//
// Java: PrintStreamErrorReportGenerator(MessageFormatter, PrintStream, int)
func NewPrintStreamErrorReportGenerator(
	formatter MessageFormatter, stream *jutil.PrintStream, summaryDetailLevel int) *PrintStreamErrorReportGenerator {
	return &PrintStreamErrorReportGenerator{formatter: formatter, stream: stream, summaryDetailLevel: summaryDetailLevel}
}

// Java: PrintStreamErrorReportGenerator.generateReport(SortingErrorManager)
func (g *PrintStreamErrorReportGenerator) GenerateReport(manager SortingErrorManager) {
	for _, e := range manager.GetSortedDiagnostics() {
		g.Println(e.Level, e.Error)
	}
	g.PrintSummary(manager)
}

// Java: PrintStreamErrorReportGenerator.println(CheckLevel, JSError)
func (g *PrintStreamErrorReportGenerator) Println(level *CheckLevel, error *JSError) {
	if level == nil {
		panic(jutil.NewNullPointerException()) // switch on a null CheckLevel in JSError.format
	}
	g.stream.PrintlnNullable(error.Format(*level, g.formatter))
}

// Java: PrintStreamErrorReportGenerator.printSummary(ErrorManager)
func (g *PrintStreamErrorReportGenerator) PrintSummary(manager ErrorManager) {
	if g.summaryDetailLevel >= 3 ||
		(g.summaryDetailLevel >= 1 && manager.GetErrorCount()+manager.GetWarningCount() > 0) ||
		(g.summaryDetailLevel >= 2 && manager.GetTypedPercent() > 0.0) {
		if manager.GetTypedPercent() > 0.0 {
			g.stream.Printf(
				"%d error(s), %d warning(s), %s%% typed%n",
				manager.GetErrorCount(), manager.GetWarningCount(), jutil.Formatter_FormatDecimalFloat(manager.GetTypedPercent(), 1))
		} else {
			g.stream.Printf(
				"%d error(s), %d warning(s)%n", manager.GetErrorCount(), manager.GetWarningCount())
		}
	}
}
