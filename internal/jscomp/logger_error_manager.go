// Port of com/google/javascript/jscomp/LoggerErrorManager.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// LoggerErrorManager is an error manager that logs errors and warnings using
// a logger in addition to collecting them in memory. Errors are logged at
// the SEVERE level and warnings are logged at the WARNING level.
type LoggerErrorManager struct {
	BasicErrorManagerBase
	formatter MessageFormatter
	logger    *jutil.Logger
}

var _ BasicErrorManager = (*LoggerErrorManager)(nil)

// NewLoggerErrorManagerWithFormatter creates an instance.
//
// Java: LoggerErrorManager(MessageFormatter, Logger)
func NewLoggerErrorManagerWithFormatter(formatter MessageFormatter, logger *jutil.Logger) *LoggerErrorManager {
	m := &LoggerErrorManager{}
	m.InitBasicErrorManagerBase(m)
	m.formatter = formatter
	m.logger = logger
	return m
}

// NewLoggerErrorManager creates an instance with a source-less error
// formatter.
//
// Java: LoggerErrorManager(Logger)
func NewLoggerErrorManager(logger *jutil.Logger) *LoggerErrorManager {
	return NewLoggerErrorManagerWithFormatter(LightweightMessageFormatter_WithoutSource(), logger)
}

// Java: LoggerErrorManager.println(CheckLevel, JSError)
func (m *LoggerErrorManager) Println(level *CheckLevel, error *JSError) {
	if level == nil {
		panic(jutil.NewNullPointerException()) // switch on a null enum
	}
	switch *level {
	case CheckLevel_ERROR:
		m.logger.Severe(*error.Format(*level, m.formatter))
	case CheckLevel_WARNING:
		m.logger.Warning(*error.Format(*level, m.formatter))
	case CheckLevel_OFF:
	}
}

// Java: LoggerErrorManager.printSummary()
func (m *LoggerErrorManager) PrintSummary() {
	level := jutil.Level_WARNING
	if m.GetErrorCount()+m.GetWarningCount() == 0 {
		level = jutil.Level_INFO
	}
	if m.GetTypedPercent() > 0.0 {
		m.logger.Log(
			level,
			jutil.StringFormat(
				"%d error(s), %d warning(s), %s%% typed",
				m.GetErrorCount(), m.GetWarningCount(), jutil.Formatter_FormatDecimalFloat(m.GetTypedPercent(), 1)))
	} else if m.GetErrorCount()+m.GetWarningCount() > 0 {
		m.logger.Log(
			level, jutil.StringFormat("%d error(s), %d warning(s)", m.GetErrorCount(), m.GetWarningCount()))
	}
}
