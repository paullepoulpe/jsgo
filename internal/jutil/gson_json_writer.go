// Port of com/google/gson/stream/JsonWriter.java (Gson 2.9.1)

/*
 * Copyright (C) 2010 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jutil

import "io"

// Go-only: language-adapter — Gson is a library dependency of Closure; it is
// ported here as Java-semantics runtime. The Java Writer is an
// io.StringWriter; a write error becomes the IOException Writer.write
// throws. The value(double/float/long/Boolean) overloads are represented by
// ValueBoolean and ValueNumber; flush() is not ported.

var (
	jsonWriter_REPLACEMENT_CHARS           [128]*string
	jsonWriter_HTML_SAFE_REPLACEMENT_CHARS [128]*string
)

func init() {
	for i := 0; i <= 0x1f; i++ {
		jsonWriter_REPLACEMENT_CHARS[i] = StringPtr(StringFormat("\\u%04x", int32(i)))
	}
	jsonWriter_REPLACEMENT_CHARS['"'] = StringPtr("\\\"")
	jsonWriter_REPLACEMENT_CHARS['\\'] = StringPtr("\\\\")
	jsonWriter_REPLACEMENT_CHARS['\t'] = StringPtr("\\t")
	jsonWriter_REPLACEMENT_CHARS['\b'] = StringPtr("\\b")
	jsonWriter_REPLACEMENT_CHARS['\n'] = StringPtr("\\n")
	jsonWriter_REPLACEMENT_CHARS['\r'] = StringPtr("\\r")
	jsonWriter_REPLACEMENT_CHARS['\f'] = StringPtr("\\f")
	jsonWriter_HTML_SAFE_REPLACEMENT_CHARS = jsonWriter_REPLACEMENT_CHARS
	jsonWriter_HTML_SAFE_REPLACEMENT_CHARS['<'] = StringPtr("\\u003c")
	jsonWriter_HTML_SAFE_REPLACEMENT_CHARS['>'] = StringPtr("\\u003e")
	jsonWriter_HTML_SAFE_REPLACEMENT_CHARS['&'] = StringPtr("\\u0026")
	jsonWriter_HTML_SAFE_REPLACEMENT_CHARS['='] = StringPtr("\\u003d")
	jsonWriter_HTML_SAFE_REPLACEMENT_CHARS['\''] = StringPtr("\\u0027")
}

// JsonWriter is com.google.gson.stream.JsonWriter: writes a JSON (RFC 7159)
// encoded value to a stream, one token at a time.
type JsonWriter struct {
	// The output data, containing at most one top-level array or object.
	out io.StringWriter

	stack     []int
	stackSize int

	// A string containing a full set of spaces for a single level of
	// indentation, or null for no pretty printing.
	indent *string

	// The name/value separator; either ":" or ": ".
	separator string

	lenient bool

	htmlSafe bool

	deferredName *string

	serializeNulls bool
}

// NewJsonWriter creates a new instance that writes a JSON-encoded stream to
// out.
//
// Java: JsonWriter(Writer)
func NewJsonWriter(out io.StringWriter) *JsonWriter {
	w := &JsonWriter{stack: make([]int, 32), separator: ":", serializeNulls: true}
	w.push(jsonScope_EMPTY_DOCUMENT)
	if out == nil {
		panic(NewNullPointerExceptionMsg("out == null"))
	}
	w.out = out
	return w
}

func (w *JsonWriter) write(s string) {
	if _, err := w.out.WriteString(s); err != nil {
		panic(NewIOException(err.Error()))
	}
}

// SetIndent sets the indentation string to be repeated for each level of
// indentation in the encoded document.
//
// Java: JsonWriter.setIndent(String)
func (w *JsonWriter) SetIndent(indent string) {
	if len(indent) == 0 {
		w.indent = nil
		w.separator = ":"
	} else {
		w.indent = &indent
		w.separator = ": "
	}
}

// SetLenient configures this writer to relax its syntax rules.
//
// Java: JsonWriter.setLenient(boolean)
func (w *JsonWriter) SetLenient(lenient bool) {
	w.lenient = lenient
}

// IsLenient returns true if this writer has relaxed syntax rules.
//
// Java: JsonWriter.isLenient()
func (w *JsonWriter) IsLenient() bool {
	return w.lenient
}

// SetHtmlSafe configures this writer to emit JSON that's safe for direct
// inclusion in HTML and XML documents.
//
// Java: JsonWriter.setHtmlSafe(boolean)
func (w *JsonWriter) SetHtmlSafe(htmlSafe bool) {
	w.htmlSafe = htmlSafe
}

// IsHtmlSafe returns true if this writer writes JSON that's safe for
// inclusion in HTML and XML documents.
//
// Java: JsonWriter.isHtmlSafe()
func (w *JsonWriter) IsHtmlSafe() bool {
	return w.htmlSafe
}

// SetSerializeNulls sets whether object members are serialized when their
// value is null.
//
// Java: JsonWriter.setSerializeNulls(boolean)
func (w *JsonWriter) SetSerializeNulls(serializeNulls bool) {
	w.serializeNulls = serializeNulls
}

// GetSerializeNulls returns true if object members are serialized when
// their value is null.
//
// Java: JsonWriter.getSerializeNulls()
func (w *JsonWriter) GetSerializeNulls() bool {
	return w.serializeNulls
}

// BeginArray begins encoding a new array.
//
// Java: JsonWriter.beginArray()
func (w *JsonWriter) BeginArray() *JsonWriter {
	w.writeDeferredName()
	return w.open(jsonScope_EMPTY_ARRAY, "[")
}

// EndArray ends encoding the current array.
//
// Java: JsonWriter.endArray()
func (w *JsonWriter) EndArray() *JsonWriter {
	return w.close(jsonScope_EMPTY_ARRAY, jsonScope_NONEMPTY_ARRAY, "]")
}

// BeginObject begins encoding a new object.
//
// Java: JsonWriter.beginObject()
func (w *JsonWriter) BeginObject() *JsonWriter {
	w.writeDeferredName()
	return w.open(jsonScope_EMPTY_OBJECT, "{")
}

// EndObject ends encoding the current object.
//
// Java: JsonWriter.endObject()
func (w *JsonWriter) EndObject() *JsonWriter {
	return w.close(jsonScope_EMPTY_OBJECT, jsonScope_NONEMPTY_OBJECT, "}")
}

// open enters a new scope by appending any necessary whitespace and the
// given bracket.
func (w *JsonWriter) open(empty int, openBracket string) *JsonWriter {
	w.beforeValue()
	w.push(empty)
	w.write(openBracket)
	return w
}

// close closes the current scope by appending any necessary whitespace and
// the given bracket.
func (w *JsonWriter) close(empty, nonempty int, closeBracket string) *JsonWriter {
	context := w.peek()
	if context != nonempty && context != empty {
		panic(NewIllegalStateExceptionMsg("Nesting problem."))
	}
	if w.deferredName != nil {
		panic(NewIllegalStateExceptionMsg("Dangling name: " + *w.deferredName))
	}

	w.stackSize--
	if context == nonempty {
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

// peek returns the value on the top of the stack.
func (w *JsonWriter) peek() int {
	if w.stackSize == 0 {
		panic(NewIllegalStateExceptionMsg("JsonWriter is closed."))
	}
	return w.stack[w.stackSize-1]
}

// replaceTop replaces the value on the top of the stack with the given
// value.
func (w *JsonWriter) replaceTop(topOfStack int) {
	w.stack[w.stackSize-1] = topOfStack
}

// Name encodes the property name.
//
// Java: JsonWriter.name(String)
func (w *JsonWriter) Name(name string) *JsonWriter {
	if w.deferredName != nil {
		panic(NewIllegalStateException())
	}
	if w.stackSize == 0 {
		panic(NewIllegalStateExceptionMsg("JsonWriter is closed."))
	}
	w.deferredName = &name
	return w
}

func (w *JsonWriter) writeDeferredName() {
	if w.deferredName != nil {
		w.beforeName()
		w.string(*w.deferredName)
		w.deferredName = nil
	}
}

// ValueString encodes value.
//
// Java: JsonWriter.value(String)
func (w *JsonWriter) ValueString(value *string) *JsonWriter {
	if value == nil {
		return w.NullValue()
	}
	w.writeDeferredName()
	w.beforeValue()
	w.string(*value)
	return w
}

// JsonValue writes value directly to the writer without quoting or
// escaping.
//
// Java: JsonWriter.jsonValue(String)
func (w *JsonWriter) JsonValue(value *string) *JsonWriter {
	if value == nil {
		return w.NullValue()
	}
	w.writeDeferredName()
	w.beforeValue()
	w.write(*value)
	return w
}

// NullValue encodes null.
//
// Java: JsonWriter.nullValue()
func (w *JsonWriter) NullValue() *JsonWriter {
	if w.deferredName != nil {
		if w.serializeNulls {
			w.writeDeferredName()
		} else {
			w.deferredName = nil
			return w // skip the name and the value
		}
	}
	w.beforeValue()
	w.write("null")
	return w
}

// ValueBoolean encodes value.
//
// Java: JsonWriter.value(boolean)
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

// isTrustedNumberType reports whether the Number's class always creates a
// valid JSON number string (Integer, Long, Double, ...). LazilyParsedNumber
// is not trusted because it could contain an arbitrary malformed string.
func jsonWriter_isTrustedNumberType(value JavaNumber) bool {
	return false
}

// javaNumberClass is getClass().toString() of a JavaNumber.
func javaNumberClass(value JavaNumber) string {
	switch value.(type) {
	case *LazilyParsedNumber:
		return "class " + GsonPackage + ".internal.LazilyParsedNumber"
	}
	panic(NewIllegalStateException())
}

// jsonWriter_isValidJsonNumber matches VALID_JSON_NUMBER_PATTERN,
// -?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?(?:[eE][-+]?[0-9]+)?, against the whole
// string.
func jsonWriter_isValidJsonNumber(s string) bool {
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

// ValueNumber encodes value.
//
// Java: JsonWriter.value(Number)
func (w *JsonWriter) ValueNumber(value JavaNumber) *JsonWriter {
	if value == nil {
		return w.NullValue()
	}

	w.writeDeferredName()
	s := value.String()
	if s == "-Infinity" || s == "Infinity" || s == "NaN" {
		if !w.lenient {
			panic(NewIllegalArgumentExceptionMsg("Numeric values must be finite, but was " + s))
		}
	} else {
		// Validate that string is valid before writing it directly to JSON
		// output
		if !jsonWriter_isTrustedNumberType(value) && !jsonWriter_isValidJsonNumber(s) {
			panic(NewIllegalArgumentExceptionMsg("String created by " + javaNumberClass(value) + " is not a valid JSON number: " + s))
		}
	}

	w.beforeValue()
	w.write(s)
	return w
}

// Close flushes and closes this writer and the underlying Writer.
//
// Go-only: language-adapter — the Writer is an io.StringWriter: it is
// closed through its Close method when it has one (a strings.Builder, like
// Java's StringWriter, has nothing to close).
//
// Java: JsonWriter.close()
func (w *JsonWriter) Close() {
	switch out := w.out.(type) {
	case interface{ Close() }:
		out.Close()
	case io.Closer:
		if err := out.Close(); err != nil {
			panic(NewIOException(err.Error()))
		}
	}
	size := w.stackSize
	if size > 1 || size == 1 && w.stack[size-1] != jsonScope_NONEMPTY_DOCUMENT {
		panic(NewIOException("Incomplete document"))
	}
	w.stackSize = 0
}

func (w *JsonWriter) string(value string) {
	replacements := &jsonWriter_REPLACEMENT_CHARS
	if w.htmlSafe {
		replacements = &jsonWriter_HTML_SAFE_REPLACEMENT_CHARS
	}
	w.write("\"")
	units := StringToUTF16(value)
	last := 0
	length := len(units)
	for i := 0; i < length; i++ {
		c := units[i]
		var replacement string
		if c < 128 {
			r := replacements[c]
			if r == nil {
				continue
			}
			replacement = *r
		} else if c == 0x2028 {
			replacement = "\\u2028"
		} else if c == 0x2029 {
			replacement = "\\u2029"
		} else {
			continue
		}
		if last < i {
			w.write(StringFromUTF16(units[last:i]))
		}
		w.write(replacement)
		last = i + 1
	}
	if last < length {
		w.write(StringFromUTF16(units[last:length]))
	}
	w.write("\"")
}

func (w *JsonWriter) newline() {
	if w.indent == nil {
		return
	}

	w.write("\n")
	for i, size := 1, w.stackSize; i < size; i++ {
		w.write(*w.indent)
	}
}

// beforeName inserts any necessary separators and whitespace before a name.
// Also adjusts the stack to expect the name's value.
func (w *JsonWriter) beforeName() {
	context := w.peek()
	if context == jsonScope_NONEMPTY_OBJECT { // first in object
		w.write(",")
	} else if context != jsonScope_EMPTY_OBJECT { // not in an object!
		panic(NewIllegalStateExceptionMsg("Nesting problem."))
	}
	w.newline()
	w.replaceTop(jsonScope_DANGLING_NAME)
}

// beforeValue inserts any necessary separators and whitespace before a
// literal value, inline array, or inline object. Also adjusts the stack to
// expect either a closing bracket or another element.
func (w *JsonWriter) beforeValue() {
	switch w.peek() {
	case jsonScope_NONEMPTY_DOCUMENT:
		if !w.lenient {
			panic(NewIllegalStateExceptionMsg("JSON must have only one top-level value."))
		}
		// fall-through
		w.replaceTop(jsonScope_NONEMPTY_DOCUMENT)
	case jsonScope_EMPTY_DOCUMENT: // first in document
		w.replaceTop(jsonScope_NONEMPTY_DOCUMENT)
	case jsonScope_EMPTY_ARRAY: // first in array
		w.replaceTop(jsonScope_NONEMPTY_ARRAY)
		w.newline()
	case jsonScope_NONEMPTY_ARRAY: // another in array
		w.write(",")
		w.newline()
	case jsonScope_DANGLING_NAME: // value for name
		w.write(w.separator)
		w.replaceTop(jsonScope_NONEMPTY_OBJECT)
	default:
		panic(NewIllegalStateExceptionMsg("Nesting problem."))
	}
}
