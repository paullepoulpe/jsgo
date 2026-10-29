package jutil

import (
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Go-only: language-adapter — java.util.logging: Logger, Level, LogRecord
// and Handler, with the default configuration of the JDK (a root logger at
// INFO whose consoleHandler writes SimpleFormatter output to stderr).
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

func (r *LogRecord) GetLevel() *Level { return r.level }

// Handler is java.util.logging.Handler (the part loggers call).
type Handler interface {
	Publish(record *LogRecord)
}

// Logger is java.util.logging.Logger.
type Logger struct {
	name              string
	level             atomic.Pointer[Level]
	handlers          []Handler
	useParentHandlers bool
}

func newRootLogger() *Logger {
	l := &Logger{
		name:              "",
		useParentHandlers: true,
		handlers:          []Handler{&consoleHandler{out: os.Stderr, level: Level_INFO}},
	}
	l.level.Store(Level_INFO)
	return l
}

var (
	loggersMu  sync.RWMutex
	loggerRoot = newRootLogger()
	loggers    = map[string]*Logger{"": loggerRoot}
)

// Logger_GetLogger is Logger.getLogger(String): one logger per name.
func Logger_GetLogger(name string) *Logger {
	loggersMu.RLock()
	if l, ok := loggers[name]; ok {
		loggersMu.RUnlock()
		return l
	}
	loggersMu.RUnlock()

	loggersMu.Lock()
	defer loggersMu.Unlock()
	if l, ok := loggers[name]; ok {
		return l
	}
	l := &Logger{name: name, useParentHandlers: true}
	loggers[name] = l
	return l
}

// GetName is Logger.getName().
func (l *Logger) GetName() string { return l.name }

func (l *Logger) parentLocked() *Logger {
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

// GetParent is Logger.getParent(): the nearest existing ancestor in the
// dotted name hierarchy, else the root logger (nil for the root).
func (l *Logger) GetParent() *Logger {
	if l == loggerRoot {
		return nil
	}
	loggersMu.RLock()
	p := l.parentLocked()
	loggersMu.RUnlock()
	return p
}

// SetLevel is Logger.setLevel(Level); nil inherits the parent's level.
func (l *Logger) SetLevel(level *Level) { l.level.Store(level) }

func (l *Logger) effectiveLevel() *Level {
	if lv := l.level.Load(); lv != nil {
		return lv
	}
	loggersMu.RLock()
	for x := l.parentLocked(); x != nil; x = x.parentLocked() {
		if lv := x.level.Load(); lv != nil {
			loggersMu.RUnlock()
			return lv
		}
	}
	loggersMu.RUnlock()
	return Level_INFO
}

// IsLoggable is Logger.isLoggable(Level).
func (l *Logger) IsLoggable(level *Level) bool {
	eff := l.effectiveLevel()
	return level.value >= eff.value && eff.value != Level_OFF.value
}

// AddHandler is Logger.addHandler(Handler).
func (l *Logger) AddHandler(h Handler) {
	loggersMu.Lock()
	l.handlers = append(l.handlers, h)
	loggersMu.Unlock()
}

// RemoveHandler is Logger.removeHandler(Handler).
func (l *Logger) RemoveHandler(h Handler) {
	loggersMu.Lock()
	defer loggersMu.Unlock()
	for i, x := range l.handlers {
		if x == h {
			l.handlers = append(l.handlers[:i:i], l.handlers[i+1:]...)
			return
		}
	}
}

// SetUseParentHandlers is Logger.setUseParentHandlers(boolean).
func (l *Logger) SetUseParentHandlers(b bool) {
	loggersMu.Lock()
	l.useParentHandlers = b
	loggersMu.Unlock()
}

func (l *Logger) logRecord(r *LogRecord) {
	if !l.IsLoggable(r.level) {
		return
	}
	var buf [4]Handler
	hs := buf[:0]
	loggersMu.RLock()
	for x := l; x != nil; x = x.parentLocked() {
		hs = append(hs, x.handlers...)
		if !x.useParentHandlers {
			break
		}
	}
	loggersMu.RUnlock()
	for _, h := range hs {
		h.Publish(r)
	}
}

func (l *Logger) newRecord(level *Level, msg string) *LogRecord {
	return &LogRecord{
		level:      level,
		message:    msg,
		loggerName: l.name,
		millis:     time.Now().UnixMilli(),
	}
}

// Log is Logger.log(Level, String).
func (l *Logger) Log(level *Level, msg string) {
	if !l.IsLoggable(level) {
		return
	}
	l.logRecord(l.newRecord(level, msg))
}

// LogWithParam is Logger.log(Level, String, Object).
func (l *Logger) LogWithParam(level *Level, msg string, param any) {
	if !l.IsLoggable(level) {
		return
	}
	r := l.newRecord(level, msg)
	r.parameters = []any{param}
	l.logRecord(r)
}

// LogWithThrown is Logger.log(Level, String, Throwable).
func (l *Logger) LogWithThrown(level *Level, msg string, thrown Throwable) {
	if !l.IsLoggable(level) {
		return
	}
	r := l.newRecord(level, msg)
	r.thrown = thrown
	l.logRecord(r)
}

// Severe is Logger.severe(String).
func (l *Logger) Severe(msg string) { l.Log(Level_SEVERE, msg) }

// Warning is Logger.warning(String).
func (l *Logger) Warning(msg string) { l.Log(Level_WARNING, msg) }

// Info is Logger.info(String).
func (l *Logger) Info(msg string) { l.Log(Level_INFO, msg) }

// Fine is Logger.fine(String).
func (l *Logger) Fine(msg string) { l.Log(Level_FINE, msg) }

// Finest is Logger.finest(String).
func (l *Logger) Finest(msg string) { l.Log(Level_FINEST, msg) }

// consoleHandler is java.util.logging.ConsoleHandler with SimpleFormatter.
type consoleHandler struct {
	mu    sync.Mutex
	out   io.Writer
	level *Level
}

// Publish is ConsoleHandler.publish(LogRecord).
func (h *consoleHandler) Publish(r *LogRecord) {
	if r.level.value < h.level.value || h.level.value == Level_OFF.value {
		return
	}
	formatted := simpleFormatterFormat(r)
	h.mu.Lock()
	_, _ = io.WriteString(h.out, formatted)
	h.mu.Unlock()
}

// simpleFormatterFormat is SimpleFormatter.format(LogRecord) with the
// default format "%1$tb %1$td, %1$tY %1$tl:%1$tM:%1$tS %1$Tp %2$s%n%4$s:
// %5$s%6$s%n" (source = logger name).
func simpleFormatterFormat(r *LogRecord) string {
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
	var sb strings.Builder
	sb.Grow(len(format) + 16)
	quoted := false
	for i := 0; i < len(format); i++ {
		c := format[i]
		if c == '\'' {
			if i+1 < len(format) && format[i+1] == '\'' {
				sb.WriteByte('\'')
				i++
			} else {
				quoted = !quoted
			}
			continue
		}
		if c == '{' && !quoted {
			if end := strings.IndexByte(format[i+1:], '}'); end >= 0 {
				if idx, err := strconv.Atoi(format[i+1 : i+1+end]); err == nil && idx >= 0 {
					if idx < len(r.parameters) {
						sb.WriteString(ToJavaString(r.parameters[idx]))
					} else {
						sb.WriteByte('{')
						sb.WriteString(strconv.Itoa(idx))
						sb.WriteByte('}')
					}
					i += 1 + end
					continue
				}
			}
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

// levelKnown are the levels the Level class defines, which Level.parse
// knows by name and by value.
var levelKnown = []*Level{Level_OFF, Level_SEVERE, Level_WARNING, Level_INFO, Level_CONFIG, Level_FINE, Level_FINER, Level_FINEST, Level_ALL}

// Level_Parse is Level.parse(String): a known level name (case-sensitive),
// or an integer value, which yields the known level of that value or a new
// Level named by the string. Anything else throws
// IllegalArgumentException("Bad level \"<name>\"").
func Level_Parse(name string) *Level {
	for _, l := range levelKnown {
		if l.name == name {
			return l
		}
	}
	if v, ok := parseRadixInteger(name, 10, math.MaxInt32); ok {
		value := int32(v)
		for _, l := range levelKnown {
			if l.value == value {
				return l
			}
		}
		return &Level{name: name, value: value}
	}
	panic(NewIllegalArgumentExceptionMsg("Bad level \"" + name + "\""))
}

// Go-only: language-adapter — java.lang.Thread identity. The port runs
// sequentially (PORTING.md §5): there is exactly one thread, so
// Thread.currentThread() always returns the same object and ThreadLocals are
// plain variables.

// Thread is java.lang.Thread (identity only).
type Thread struct {
	name string
}

var thread_main = &Thread{name: "main"}

// Thread_CurrentThread is Thread.currentThread().
func Thread_CurrentThread() *Thread { return thread_main }

// GetName is Thread.getName().
func (t *Thread) GetName() string { return t.name }

// String is Thread.toString().
func (t *Thread) String() string { return "Thread[#1," + t.name + ",5,main]" }
