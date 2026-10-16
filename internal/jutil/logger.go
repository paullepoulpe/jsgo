package jutil

import (
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// Go-only: language-adapter — java.util.logging: Logger, Level, LogRecord
// and Handler, with the default configuration of the JDK (a root logger at
// INFO whose ConsoleHandler writes SimpleFormatter output to stderr).
//
// Differences that Go forces: SimpleFormatter's "source" is the logger name
// (Java infers the calling class and method by walking the stack), and a
// thrown exception is rendered as its toString() without stack frames.
// Message parameters are substituted like java.text.MessageFormat for
// "{n}" placeholders with their Java string values.

// Level is java.util.logging.Level.
type Level struct {
	name  string
	value int32
}

var (
	Level_OFF     = &Level{"OFF", math.MaxInt32}
	Level_SEVERE  = &Level{"SEVERE", 1000}
	Level_WARNING = &Level{"WARNING", 900}
	Level_INFO    = &Level{"INFO", 800}
	Level_CONFIG  = &Level{"CONFIG", 700}
	Level_FINE    = &Level{"FINE", 500}
	Level_FINER   = &Level{"FINER", 400}
	Level_FINEST  = &Level{"FINEST", 300}
	Level_ALL     = &Level{"ALL", math.MinInt32}
)

// GetName is Level.getName().
func (l *Level) GetName() string { return l.name }

// IntValue is Level.intValue().
func (l *Level) IntValue() int32 { return l.value }

// String is Level.toString().
func (l *Level) String() string { return l.name }

// LogRecord is java.util.logging.LogRecord.
type LogRecord struct {
	level      *Level
	message    string
	parameters []any
	thrown     Throwable
	loggerName string
	millis     int64
}

// NewLogRecord is new LogRecord(Level, String).
func NewLogRecord(level *Level, msg string) *LogRecord {
	return &LogRecord{level: level, message: msg, millis: time.Now().UnixMilli()}
}

func (r *LogRecord) GetLevel() *Level       { return r.level }
func (r *LogRecord) GetMessage() string     { return r.message }
func (r *LogRecord) GetParameters() []any   { return r.parameters }
func (r *LogRecord) SetParameters(p []any)  { r.parameters = p }
func (r *LogRecord) GetThrown() Throwable   { return r.thrown }
func (r *LogRecord) SetThrown(t Throwable)  { r.thrown = t }
func (r *LogRecord) GetLoggerName() string  { return r.loggerName }
func (r *LogRecord) SetLoggerName(n string) { r.loggerName = n }
func (r *LogRecord) GetMillis() int64       { return r.millis }

// Handler is java.util.logging.Handler (the part loggers call).
type Handler interface {
	Publish(record *LogRecord)
}

// Logger is java.util.logging.Logger.
type Logger struct {
	name              string
	level             *Level
	handlers          []Handler
	useParentHandlers bool
}

var (
	loggerRoot = &Logger{name: "", level: Level_INFO, useParentHandlers: true,
		handlers: []Handler{&ConsoleHandler{out: os.Stderr, level: Level_INFO}}}
	loggers = map[string]*Logger{"": loggerRoot}
)

// Logger_GetLogger is Logger.getLogger(String): one logger per name.
func Logger_GetLogger(name string) *Logger {
	if l, ok := loggers[name]; ok {
		return l
	}
	l := &Logger{name: name, useParentHandlers: true}
	loggers[name] = l
	return l
}

// Logger_GetGlobal is Logger.getGlobal().
func Logger_GetGlobal() *Logger { return Logger_GetLogger("global") }

// GetName is Logger.getName().
func (l *Logger) GetName() string { return l.name }

// GetParent is Logger.getParent(): the nearest existing ancestor in the
// dotted name hierarchy, else the root logger (nil for the root).
func (l *Logger) GetParent() *Logger {
	if l == loggerRoot {
		return nil
	}
	name := l.name
	for {
		i := strings.LastIndexByte(name, '.')
		if i < 0 {
			return loggerRoot
		}
		name = name[:i]
		if p, ok := loggers[name]; ok {
			return p
		}
	}
}

// SetLevel is Logger.setLevel(Level); nil inherits the parent's level.
func (l *Logger) SetLevel(level *Level) { l.level = level }

// GetLevel is Logger.getLevel().
func (l *Logger) GetLevel() *Level { return l.level }

func (l *Logger) effectiveLevel() *Level {
	for x := l; x != nil; x = x.GetParent() {
		if x.level != nil {
			return x.level
		}
	}
	return Level_INFO
}

// IsLoggable is Logger.isLoggable(Level).
func (l *Logger) IsLoggable(level *Level) bool {
	eff := l.effectiveLevel()
	return level.value >= eff.value && eff.value != Level_OFF.value
}

// AddHandler is Logger.addHandler(Handler).
func (l *Logger) AddHandler(h Handler) { l.handlers = append(l.handlers, h) }

// RemoveHandler is Logger.removeHandler(Handler).
func (l *Logger) RemoveHandler(h Handler) {
	for i, x := range l.handlers {
		if x == h {
			l.handlers = append(l.handlers[:i:i], l.handlers[i+1:]...)
			return
		}
	}
}

// GetHandlers is Logger.getHandlers().
func (l *Logger) GetHandlers() []Handler { return append([]Handler(nil), l.handlers...) }

// SetUseParentHandlers is Logger.setUseParentHandlers(boolean).
func (l *Logger) SetUseParentHandlers(b bool) { l.useParentHandlers = b }

// LogRecordTo is Logger.log(LogRecord).
func (l *Logger) LogRecordTo(r *LogRecord) {
	if !l.IsLoggable(r.level) {
		return
	}
	for x := l; x != nil; x = x.GetParent() {
		for _, h := range x.handlers {
			h.Publish(r)
		}
		if !x.useParentHandlers {
			break
		}
	}
}

func (l *Logger) newRecord(level *Level, msg string) *LogRecord {
	r := NewLogRecord(level, msg)
	r.loggerName = l.name
	return r
}

// Log is Logger.log(Level, String).
func (l *Logger) Log(level *Level, msg string) {
	if !l.IsLoggable(level) {
		return
	}
	l.LogRecordTo(l.newRecord(level, msg))
}

// LogWithParam is Logger.log(Level, String, Object).
func (l *Logger) LogWithParam(level *Level, msg string, param any) {
	if !l.IsLoggable(level) {
		return
	}
	r := l.newRecord(level, msg)
	r.parameters = []any{param}
	l.LogRecordTo(r)
}

// LogWithParams is Logger.log(Level, String, Object[]).
func (l *Logger) LogWithParams(level *Level, msg string, params []any) {
	if !l.IsLoggable(level) {
		return
	}
	r := l.newRecord(level, msg)
	r.parameters = params
	l.LogRecordTo(r)
}

// LogWithThrown is Logger.log(Level, String, Throwable).
func (l *Logger) LogWithThrown(level *Level, msg string, thrown Throwable) {
	if !l.IsLoggable(level) {
		return
	}
	r := l.newRecord(level, msg)
	r.thrown = thrown
	l.LogRecordTo(r)
}

// Severe is Logger.severe(String).
func (l *Logger) Severe(msg string) { l.Log(Level_SEVERE, msg) }

// Warning is Logger.warning(String).
func (l *Logger) Warning(msg string) { l.Log(Level_WARNING, msg) }

// Info is Logger.info(String).
func (l *Logger) Info(msg string) { l.Log(Level_INFO, msg) }

// Config is Logger.config(String).
func (l *Logger) Config(msg string) { l.Log(Level_CONFIG, msg) }

// Fine is Logger.fine(String).
func (l *Logger) Fine(msg string) { l.Log(Level_FINE, msg) }

// Finer is Logger.finer(String).
func (l *Logger) Finer(msg string) { l.Log(Level_FINER, msg) }

// Finest is Logger.finest(String).
func (l *Logger) Finest(msg string) { l.Log(Level_FINEST, msg) }

// ConsoleHandler is java.util.logging.ConsoleHandler with SimpleFormatter.
type ConsoleHandler struct {
	out   io.Writer
	level *Level
}

// Publish is ConsoleHandler.publish(LogRecord).
func (h *ConsoleHandler) Publish(r *LogRecord) {
	if r.level.value < h.level.value || h.level.value == Level_OFF.value {
		return
	}
	io.WriteString(h.out, SimpleFormatter_Format(r))
}

// SimpleFormatter_Format is SimpleFormatter.format(LogRecord) with the
// default format "%1$tb %1$td, %1$tY %1$tl:%1$tM:%1$tS %1$Tp %2$s%n%4$s:
// %5$s%6$s%n" (source = logger name).
func SimpleFormatter_Format(r *LogRecord) string {
	var sb strings.Builder
	sb.WriteString(time.UnixMilli(r.millis).Format("Jan 02, 2006 3:04:05 PM"))
	sb.WriteByte(' ')
	sb.WriteString(r.loggerName)
	sb.WriteByte('\n')
	sb.WriteString(r.level.name)
	sb.WriteString(": ")
	sb.WriteString(Formatter_FormatMessage(r))
	if r.thrown != nil {
		sb.WriteByte('\n')
		sb.WriteString(r.thrown.Error())
	}
	sb.WriteByte('\n')
	return sb.String()
}

// Formatter_FormatMessage is Formatter.formatMessage(LogRecord): the
// parameters are substituted only if the message contains "{0" to "{3".
func Formatter_FormatMessage(r *LogRecord) string {
	format := r.message
	if len(r.parameters) == 0 {
		return format
	}
	if !strings.Contains(format, "{0") && !strings.Contains(format, "{1") &&
		!strings.Contains(format, "{2") && !strings.Contains(format, "{3") {
		return format
	}
	return MessageFormat_Format(format, r.parameters)
}

// MessageFormat_Format is java.text.MessageFormat.format(String, Object...)
// for simple "{n}" arguments (no format types), with MessageFormat's
// quoting: ” is a quote and text between single quotes is literal.
func MessageFormat_Format(pattern string, args []any) string {
	var sb strings.Builder
	inQuote := false
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch {
		case c == '\'':
			if i+1 < len(pattern) && pattern[i+1] == '\'' {
				sb.WriteByte('\'')
				i++
			} else {
				inQuote = !inQuote
			}
		case c == '{' && !inQuote:
			end := strings.IndexByte(pattern[i:], '}')
			if end < 0 {
				panic(NewIllegalArgumentExceptionMsg("Unmatched braces in the pattern."))
			}
			n, err := strconv.Atoi(strings.TrimSpace(pattern[i+1 : i+end]))
			if err != nil {
				panic(NewIllegalArgumentExceptionMsg("can't parse argument number: " + pattern[i+1:i+end]))
			}
			if n < len(args) {
				sb.WriteString(ToJavaString(args[n]))
			} else {
				sb.WriteString(pattern[i : i+end+1])
			}
			i += end
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}
