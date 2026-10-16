// Port of com/google/javascript/jscomp/diagnostic/NoOpLogFile.java @ closure-compiler v20260415

package diagnostic

// NoOpLogFile is an implementation that does nothing to allow logging calls
// to be cheap when disabled.
//
// Java: final class NoOpLogFile extends LogFile
type NoOpLogFile struct {
	LogFileBase
}

// NewNoOpLogFile constructs a NoOpLogFile.
//
// Java: NoOpLogFile()
func NewNoOpLogFile() *NoOpLogFile {
	n := &NoOpLogFile{}
	n.InitLogFileBase(n)
	return n
}

// Java: NoOpLogFile.log(Object)
func (n *NoOpLogFile) Log(value any) LogFile {
	return n
}

// Java: NoOpLogFile.log(String)
func (n *NoOpLogFile) LogString(value string) LogFile {
	return n
}

// Java: NoOpLogFile.log(Supplier<String>)
func (n *NoOpLogFile) LogSupplier(value func() string) LogFile {
	return n
}

// Java: NoOpLogFile.log(String, Object...)
func (n *NoOpLogFile) LogFormat(template string, values ...any) LogFile {
	return n
}

// Java: NoOpLogFile.logJson(Object)
func (n *NoOpLogFile) LogJson(value any) LogFile {
	return n
}

// Java: NoOpLogFile.logJson(Supplier<Object>)
func (n *NoOpLogFile) LogJsonSupplier(value func() any) LogFile {
	return n
}

// Java: NoOpLogFile.logJson(StreamedJsonProducer)
func (n *NoOpLogFile) LogJsonProducer(producer LogFile_StreamedJsonProducer) LogFile {
	return n
}

// Java: NoOpLogFile.close()
func (n *NoOpLogFile) Close() {}

// Java: NoOpLogFile.isLogging()
func (n *NoOpLogFile) IsLogging() bool {
	return false
}
