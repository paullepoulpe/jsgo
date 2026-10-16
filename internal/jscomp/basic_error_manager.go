// Port of com/google/javascript/jscomp/BasicErrorManager.java @ closure-compiler v20260415

package jscomp

// BasicErrorManager is an error manager that generates a sorted report when
// the GenerateReport method is called.
//
// This error manager does not produce any output, but subclasses can
// override the Println(CheckLevel, JSError) method to generate custom
// output. Consider using the SortingErrorManager with a custom
// SortingErrorManager.ErrorReportGenerator instead.
//
// Deprecated: as in Java.
type BasicErrorManager interface {
	SortingErrorManager

	// Println prints a message with a trailing new line. This method is
	// called by the GenerateReport method when generating messages.
	//
	// Java: BasicErrorManager.println(CheckLevel, JSError)
	Println(level *CheckLevel, error *JSError)

	// PrintSummary prints the summary of the compilation - number of errors
	// and warnings.
	//
	// Java: BasicErrorManager.printSummary()
	PrintSummary()

	GetBasicErrorManagerBase() *BasicErrorManagerBase
}

// BasicErrorManagerBase holds the base implementations of
// BasicErrorManager; subclasses embed it and implement Println and
// PrintSummary.
type BasicErrorManagerBase struct {
	SortingErrorManagerBase
	basicSelf BasicErrorManager
}

// InitBasicErrorManagerBase is the Go form of super(); self is the object
// being constructed.
//
// Java: BasicErrorManager()
func (m *BasicErrorManagerBase) InitBasicErrorManagerBase(self BasicErrorManager) {
	m.InitSortingErrorManagerBase(self, nil)
	m.basicSelf = self
}

// GetBasicErrorManagerBase gives access to the base fields.
func (m *BasicErrorManagerBase) GetBasicErrorManagerBase() *BasicErrorManagerBase { return m }

// Java: BasicErrorManager.generateReport()
func (m *BasicErrorManagerBase) GenerateReport() {
	for _, message := range m.SortingErrorManagerBase.GetSortedDiagnostics() {
		m.basicSelf.Println(message.Level, message.Error)
	}
	m.basicSelf.PrintSummary()
}
