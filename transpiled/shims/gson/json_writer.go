// Clean-room Gson 2.9.1 streaming JsonWriter (Apache 2.0).

package gson

import (
	"io"
	"math"
	"strconv"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

const (
	jsonScopeEmptyArray       = 1
	jsonScopeNonemptyArray    = 2
	jsonScopeEmptyObject      = 3
	jsonScopeDanglingName     = 4
	jsonScopeNonemptyObject   = 5
	jsonScopeEmptyDocument    = 6
	jsonScopeNonemptyDocument = 7
	jsonScopeClosed           = 8
)

var (
	replacementChars = [128]string{
		0x00: `\u0000`, 0x01: `\u0001`, 0x02: `\u0002`, 0x03: `\u0003`,
		0x04: `\u0004`, 0x05: `\u0005`, 0x06: `\u0006`, 0x07: `\u0007`,
		'\b': `\b`, '\t': `\t`, '\n': `\n`, 0x0b: `\u000b`,
		'\f': `\f`, '\r': `\r`, 0x0e: `\u000e`, 0x0f: `\u000f`,
		0x10: `\u0010`, 0x11: `\u0011`, 0x12: `\u0012`, 0x13: `\u0013`,
		0x14: `\u0014`, 0x15: `\u0015`, 0x16: `\u0016`, 0x17: `\u0017`,
		0x18: `\u0018`, 0x19: `\u0019`, 0x1a: `\u001a`, 0x1b: `\u001b`,
		0x1c: `\u001c`, 0x1d: `\u001d`, 0x1e: `\u001e`, 0x1f: `\u001f`,
		'"': `\"`, '\\': `\\`,
	}
	htmlSafeReplacementChars = func() [128]string {
		t := replacementChars
		t['<'] = `\u003c`
		t['>'] = `\u003e`
		t['&'] = `\u0026`
		t['='] = `\u003d`
		t['\''] = `\u0027`
		return t
	}()
)

type stringSink interface {
	WriteString(s string) (int, error)
}

type javaWriterAdapter struct {
	w javaio.Writer
}

func (a javaWriterAdapter) WriteString(s string) (int, error) {
	a.w.Write(s)
	return len(s), nil
}

func (a javaWriterAdapter) Flush() { a.w.Flush() }
func (a javaWriterAdapter) Close() { a.w.Close() }

// JsonWriter implements com.google.gson.stream.JsonWriter.
type JsonWriter struct {
	out             stringSink
	stack           []int
	stackSize       int
	indent          string
	separator       string
	lenient         bool
	htmlSafe        bool
	deferredName    string
	hasDeferredName bool
	serializeNulls  bool
}

// NewJsonWriter creates a JsonWriter over a javaio.Writer.
func NewJsonWriter(out javaio.Writer) *JsonWriter {
	if out == nil {
		panic(javalang.NewNullPointerExceptionMsg("out == null"))
	}
	return NewJsonWriterSink(javaWriterAdapter{w: out})
}

// NewJsonWriterSink creates a JsonWriter over any io.StringWriter (e.g. *strings.Builder).
func NewJsonWriterSink(out io.StringWriter) *JsonWriter {
	if out == nil {
		panic(javalang.NewNullPointerExceptionMsg("out == null"))
	}
	w := &JsonWriter{
		out:            out,
		stack:          make([]int, 32),
		separator:      ":",
		serializeNulls: true,
	}
	w.push(jsonScopeEmptyDocument)
	return w
}

func (w *JsonWriter) write(s string) {
	if _, err := w.out.WriteString(s); err != nil {
		panic(javaio.NewIOException(err.Error()))
	}
}

func (w *JsonWriter) SetIndent(indent string) {
	if len(indent) == 0 {
		w.indent = ""
		w.separator = ":"
	} else {
		w.indent = indent
		w.separator = ": "
	}
}

func (w *JsonWriter) SetLenient(lenient bool)        { w.lenient = lenient }
func (w *JsonWriter) IsLenient() bool                { return w.lenient }
func (w *JsonWriter) SetHtmlSafe(htmlSafe bool)      { w.htmlSafe = htmlSafe }
func (w *JsonWriter) IsHtmlSafe() bool               { return w.htmlSafe }
func (w *JsonWriter) SetSerializeNulls(v bool)       { w.serializeNulls = v }
func (w *JsonWriter) GetSerializeNulls() bool        { return w.serializeNulls }

func (w *JsonWriter) BeginArray() *JsonWriter {
	w.writeDeferredName()
	return w.open(jsonScopeEmptyArray, "[")
}

func (w *JsonWriter) EndArray() *JsonWriter {
	return w.closeScope(jsonScopeEmptyArray, jsonScopeNonemptyArray, "]")
}

func (w *JsonWriter) BeginObject() *JsonWriter {
	w.writeDeferredName()
	return w.open(jsonScopeEmptyObject, "{")
}

func (w *JsonWriter) EndObject() *JsonWriter {
	return w.closeScope(jsonScopeEmptyObject, jsonScopeNonemptyObject, "}")
}

func (w *JsonWriter) open(empty int, openBracket string) *JsonWriter {
	w.beforeValue()
	w.push(empty)
	w.write(openBracket)
	return w
}

func (w *JsonWriter) closeScope(empty, nonempty int, closeBracket string) *JsonWriter {
	ctx := w.peek()
	if ctx != nonempty && ctx != empty {
		panic(javalang.NewIllegalStateException("Nesting problem."))
	}
	if w.hasDeferredName {
		panic(javalang.NewIllegalStateException("Dangling name: " + w.deferredName))
	}
	w.stackSize--
	if ctx == nonempty {
		w.newline()
	}
	w.write(closeBracket)
	return w
}

func (w *JsonWriter) push(newTop int) {
	if w.stackSize == len(w.stack) {
		w.stack = append(w.stack, make([]int, w.stackSize)...)
	}
	w.stack[w.stackSize] = newTop
	w.stackSize++
}

func (w *JsonWriter) peek() int {
	if w.stackSize == 0 {
		panic(javalang.NewIllegalStateException("JsonWriter is closed."))
	}
	return w.stack[w.stackSize-1]
}

func (w *JsonWriter) replaceTop(topOfStack int) {
	w.stack[w.stackSize-1] = topOfStack
}

func (w *JsonWriter) Name(name string) *JsonWriter {
	if w.hasDeferredName {
		panic(javalang.NewIllegalStateException("Already has deferred name"))
	}
	if w.stackSize == 0 {
		panic(javalang.NewIllegalStateException("JsonWriter is closed."))
	}
	w.deferredName = name
	w.hasDeferredName = true
	return w
}

func (w *JsonWriter) writeDeferredName() {
	if w.hasDeferredName {
		w.beforeName()
		w.writeQuoted(w.deferredName)
		w.deferredName = ""
		w.hasDeferredName = false
	}
}

func (w *JsonWriter) ValueString(value string) *JsonWriter {
	w.writeDeferredName()
	w.beforeValue()
	w.writeQuoted(value)
	return w
}

func (w *JsonWriter) ValueStringPtr(value *string) *JsonWriter {
	if value == nil {
		return w.NullValue()
	}
	return w.ValueString(*value)
}

func (w *JsonWriter) JsonValue(value string) *JsonWriter {
	w.writeDeferredName()
	w.beforeValue()
	w.write(value)
	return w
}

func (w *JsonWriter) NullValue() *JsonWriter {
	if w.hasDeferredName {
		if w.serializeNulls {
			w.writeDeferredName()
		} else {
			w.deferredName = ""
			w.hasDeferredName = false
			return w
		}
	}
	w.beforeValue()
	w.write("null")
	return w
}

func (w *JsonWriter) ValueBoolean(value bool) *JsonWriter {
	w.writeDeferredName()
	w.beforeValue()
	if value {
		w.write("true")
	} else {
		w.write("false")
	}
	return w
}

func (w *JsonWriter) ValueInt(value int32) *JsonWriter {
	return w.ValueLong(int64(value))
}

func (w *JsonWriter) ValueLong(value int64) *JsonWriter {
	w.writeDeferredName()
	w.beforeValue()
	w.write(strconv.FormatInt(value, 10))
	return w
}

func (w *JsonWriter) ValueDouble(value float64) *JsonWriter {
	w.writeDeferredName()
	if !w.lenient && (math.IsNaN(value) || math.IsInf(value, 0)) {
		panic(javalang.NewIllegalArgumentException("Numeric values must be finite, but was " + javalang.DoubleToString(value)))
	}
	w.beforeValue()
	w.write(javalang.DoubleToString(value))
	return w
}

func (w *JsonWriter) ValueNumber(value JavaNumber) *JsonWriter {
	if value == nil {
		return w.NullValue()
	}
	w.writeDeferredName()
	s := value.String()
	if s == "-Infinity" || s == "Infinity" || s == "NaN" {
		if !w.lenient {
			panic(javalang.NewIllegalArgumentException("Numeric values must be finite, but was " + s))
		}
	} else if !isValidJsonNumber(s) {
		panic(javalang.NewIllegalArgumentException("Not a valid JSON number: " + s))
	}
	w.beforeValue()
	w.write(s)
	return w
}

func isValidJsonNumber(s string) bool {
	i, n := 0, len(s)
	digits := func() int {
		start := i
		for i < n && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		return i - start
	}
	if i < n && s[i] == '-' {
		i++
	}
	if i < n && s[i] == '0' {
		i++
	} else if i < n && s[i] >= '1' && s[i] <= '9' {
		digits()
	} else {
		return false
	}
	if i < n && s[i] == '.' {
		save := i
		i++
		if digits() == 0 {
			i = save
		}
	}
	if i < n && (s[i] == 'e' || s[i] == 'E') {
		save := i
		i++
		if i < n && (s[i] == '-' || s[i] == '+') {
			i++
		}
		if digits() == 0 {
			i = save
		}
	}
	return i == n
}

func (w *JsonWriter) Flush() {
	if w.stackSize == 0 {
		panic(javalang.NewIllegalStateException("JsonWriter is closed."))
	}
	if f, ok := w.out.(interface{ Flush() }); ok {
		f.Flush()
	}
}

func (w *JsonWriter) Close() {
	switch out := w.out.(type) {
	case interface{ Close() }:
		out.Close()
	case io.Closer:
		if err := out.Close(); err != nil {
			panic(javaio.NewIOException(err.Error()))
		}
	}
	size := w.stackSize
	if size > 1 || (size == 1 && w.stack[size-1] != jsonScopeNonemptyDocument) {
		panic(javaio.NewIOException("Incomplete document"))
	}
	w.stackSize = 0
}

func (w *JsonWriter) writeQuoted(value string) {
	replacements := &replacementChars
	if w.htmlSafe {
		replacements = &htmlSafeReplacementChars
	}
	w.write("\"")
	last := 0
	length := len(value)
	for i := 0; i < length; i++ {
		b := value[i]
		var replacement string
		step := 1
		if b < 0x80 {
			replacement = replacements[b]
			if replacement == "" {
				continue
			}
		} else if b == 0xe2 && i+2 < length && value[i+1] == 0x80 {
			switch value[i+2] {
			case 0xa8:
				replacement = `\u2028`
				step = 3
			case 0xa9:
				replacement = `\u2029`
				step = 3
			default:
				continue
			}
		} else {
			continue
		}
		if last < i {
			w.write(value[last:i])
		}
		w.write(replacement)
		last = i + step
		i += step - 1
	}
	if last < length {
		w.write(value[last:length])
	}
	w.write("\"")
}

func (w *JsonWriter) newline() {
	if w.indent == "" {
		return
	}
	w.write("\n")
	for i := 1; i < w.stackSize; i++ {
		w.write(w.indent)
	}
}

func (w *JsonWriter) beforeName() {
	ctx := w.peek()
	if ctx == jsonScopeNonemptyObject {
		w.write(",")
	} else if ctx != jsonScopeEmptyObject {
		panic(javalang.NewIllegalStateException("Nesting problem."))
	}
	w.newline()
	w.replaceTop(jsonScopeDanglingName)
}

func (w *JsonWriter) beforeValue() {
	switch w.peek() {
	case jsonScopeNonemptyDocument:
		if !w.lenient {
			panic(javalang.NewIllegalStateException("JSON must have only one top-level value."))
		}
		w.replaceTop(jsonScopeNonemptyDocument)
	case jsonScopeEmptyDocument:
		w.replaceTop(jsonScopeNonemptyDocument)
	case jsonScopeEmptyArray:
		w.replaceTop(jsonScopeNonemptyArray)
		w.newline()
	case jsonScopeNonemptyArray:
		w.write(",")
		w.newline()
	case jsonScopeDanglingName:
		w.write(w.separator)
		w.replaceTop(jsonScopeNonemptyObject)
	default:
		panic(javalang.NewIllegalStateException("Nesting problem."))
	}
}
