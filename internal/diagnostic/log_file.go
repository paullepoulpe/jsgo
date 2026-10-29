// Port of com/google/javascript/jscomp/diagnostic/LogFile.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/diagnostic/WritingLogFile.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/diagnostic/NoOpLogFile.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/diagnostic/LogsGson.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package diagnostic is the port of Closure Compiler's
// com.google.javascript.jscomp.diagnostic package: the debug log file
// (LogFile) and the Gson configuration it writes JSON with (LogsGson), on
// jutil's Gson 2.9.1 runtime.
package diagnostic

import (
	"bufio"
	"errors"
	"os"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var errStreamClosed = errors.New("Stream closed")

// logWriter wraps an *os.File and *bufio.Writer, sanitizing lone WTF-8
// surrogates to '?' (matching Java's UTF-8 REPLACE encoder).
type logWriter struct {
	file     *os.File
	buf      *bufio.Writer
	highSurr rune
	closed   bool
}

func hasWTF8Surrogate(s string) bool {
	for i := 0; i+2 < len(s); i++ {
		if s[i] == 0xED && s[i+1] >= 0xA0 && s[i+1] <= 0xBF && s[i+2]&0xC0 == 0x80 {
			return true
		}
	}
	return false
}

func decodeLogRune(s string) (rune, int) {
	c := s[0]
	if c < 0x80 {
		return rune(c), 1
	}
	if c == 0xED && len(s) >= 3 && s[1] >= 0xA0 && s[1] <= 0xBF && s[2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	return utf8.DecodeRuneInString(s)
}

func (w *logWriter) WriteString(s string) (int, error) {
	if w.closed {
		return 0, errStreamClosed
	}
	if len(s) == 0 {
		return 0, nil
	}
	if w.highSurr == 0 && !hasWTF8Surrogate(s) {
		return w.buf.WriteString(s)
	}
	for i := 0; i < len(s); {
		r, width := decodeLogRune(s[i:])
		i += width
		if w.highSurr != 0 {
			high := w.highSurr
			w.highSurr = 0
			if r >= 0xDC00 && r <= 0xDFFF {
				if _, err := w.buf.WriteRune(utf16.DecodeRune(high, r)); err != nil {
					return 0, err
				}
				continue
			}
			if err := w.buf.WriteByte('?'); err != nil {
				return 0, err
			}
		}
		if r >= 0xD800 && r <= 0xDBFF {
			w.highSurr = r
		} else if r >= 0xDC00 && r <= 0xDFFF {
			if err := w.buf.WriteByte('?'); err != nil {
				return 0, err
			}
		} else {
			if _, err := w.buf.WriteRune(r); err != nil {
				return 0, err
			}
		}
	}
	return len(s), nil
}

func (w *logWriter) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true
	var flushErr error
	if w.highSurr != 0 {
		w.highSurr = 0
		flushErr = w.buf.WriteByte('?')
	}
	if err := w.buf.Flush(); err != nil && flushErr == nil {
		flushErr = err
	}
	if err := w.file.Close(); err != nil && flushErr == nil {
		flushErr = err
	}
	return flushErr
}

// LogFile is a handle for writing to a human readable log file.
//
// Its zero value (and LogFile_CreateNoOp()) has writer == nil and discards
// everything without evaluating suppliers or producers.
//
// Java: public abstract class LogFile implements AutoCloseable (collapsed with
// package-private subclasses WritingLogFile and NoOpLogFile).
type LogFile struct {
	writer *logWriter
}

// LogFile_StreamedJsonProducer is an interface allowing streaming JSON to a
// LogFile via a JsonWriter.
//
// Use with LogJsonProducer (logJson(StreamedJsonProducer)).
//
// Java: public interface LogFile.StreamedJsonProducer
type LogFile_StreamedJsonProducer interface {
	// WriteJson throws IOException.
	//
	// Java: StreamedJsonProducer.writeJson(JsonWriter)
	WriteJson(jsonWriter *jutil.JsonWriter)
}

// LogFile_CreateOrReopen creates the log file (and its directory) if
// needed and opens it for appending.
//
// Java: LogFile.createOrReopen(Path) / WritingLogFile.create(Path)
func LogFile_CreateOrReopen(file *jutil.Path) LogFile {
	dir := file.GetParent()
	if dir == nil {
		panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "java.nio.file.Path.getFileSystem()" because "path" is null`))
	}
	if err := os.MkdirAll(dir.String(), 0o777); err != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
	}
	f, err := os.OpenFile(file.String(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
	}
	return LogFile{
		writer: &logWriter{
			file: f,
			buf:  bufio.NewWriter(f),
		},
	}
}

// LogFile_CreateNoOp returns a LogFile that discards everything.
//
// Java: LogFile.createNoOp()
func LogFile_CreateNoOp() LogFile {
	return LogFile{}
}

// Java: LogFile.log(Object)
func (l LogFile) Log(value any) LogFile {
	if l.writer == nil {
		return l
	}
	switch v := value.(type) {
	case nil:
		panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "Object.toString()" because "value" is null`))
	case *string:
		if v == nil {
			panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "Object.toString()" because "value" is null`))
		}
	}
	return l.logInternal(jutil.ToJavaString(value))
}

// Provides a disambiguating overload.
//
// Java: LogFile.log(String)
func (l LogFile) LogString(value string) LogFile {
	if l.writer == nil {
		return l
	}
	return l.logInternal(value)
}

// Java: LogFile.log(Supplier<String>)
func (l LogFile) LogSupplier(value func() string) LogFile {
	if l.writer == nil {
		return l
	}
	return l.logInternal(value())
}

// Java: LogFile.log(String, Object...)
func (l LogFile) LogFormat(template string, values ...any) LogFile {
	if l.writer == nil {
		return l
	}
	return l.logInternal(jutil.StringFormat(template, values...))
}

// Java: LogFile.logJson(Object)
func (l LogFile) LogJson(value any) LogFile {
	if l.writer == nil {
		return l
	}
	return l.logInternal(LogsGson_ToJson(value))
}

// Java: LogFile.logJson(Supplier<Object>)
func (l LogFile) LogJsonSupplier(value func() any) LogFile {
	if l.writer == nil {
		return l
	}
	return l.logInternal(LogsGson_ToJson(value()))
}

// LogJsonProducer streams JSON to the file through a JsonWriter. Like
// Java's try-with-resources, closing the JsonWriter closes the file's
// writer.
//
// Java: LogFile.logJson(StreamedJsonProducer)
func (l LogFile) LogJsonProducer(producer LogFile_StreamedJsonProducer) LogFile {
	if l.writer == nil {
		return l
	}
	writer := jutil.NewJsonWriter(l.writer)
	completed := false
	defer func() {
		r := recover()
		if !completed {
			_ = l.writer.Close()
		}
		if r != nil {
			if ex, ok := r.(jutil.IIOException); ok {
				panic(jutil.NewAssertionErrorWithCause(ex))
			}
			panic(r)
		}
	}()
	producer.WriteJson(writer)
	completed = true
	writer.Close() // also closes this.writer, as Java's JsonWriter.close() does
	return l
}

// logInternal appends value and a newline.
//
// Java: WritingLogFile.logInternal(String)
func (l LogFile) logInternal(value string) LogFile {
	if _, err := l.writer.WriteString(value); err != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
	}
	if _, err := l.writer.WriteString("\n"); err != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
	}
	return l
}

// IsLogging reports whether or not the results are being recorded.
//
// Java: LogFile.isLogging()
func (l LogFile) IsLogging() bool {
	return l.writer != nil
}

// Java: LogFile.close()
func (l LogFile) Close() {
	if l.writer == nil {
		return
	}
	if err := l.writer.Close(); err != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewIOException(err.Error())))
	}
}

// LogsGson_Able is a conversion to a preferred logging representation.
//
// When an implementor is logged, the result of this method will be JSON
// stringified in place of the implementor.
//
// Use this interface when the default Gson serialization for a class is bad
// and the class should always be logged another way. Primitive wrappers and
// entity type are ideal candidates.
//
// Java: public interface LogsGson.Able
type LogsGson_Able interface {
	// Java: Able.toLogsGson()
	ToLogsGson() any
}

type logsGsonMapBacked interface {
	AsMapObject() any
}

// LogsGson_ToJson serializes o with LogsGson_GSON.
//
// Java: LogsGson.toJson(Object)
func LogsGson_ToJson(o any) string {
	return LogsGson_GSON.ToJson(o)
}

// LogsGson_GSON is Gson with two type hierarchy adapters: LogsGson_Able
// values are written as the JSON of their ToLogsGson(), and map-backed
// collections as the JSON of their asMap().
var LogsGson_GSON = jutil.NewGsonBuilder().
	RegisterTypeHierarchyAdapter(
		func(value any) bool { _, ok := value.(LogsGson_Able); return ok }, // Java: Able.class
		&logsGson_ableTypeAdapter{}).
	RegisterTypeHierarchyAdapter(
		func(value any) bool { _, ok := value.(logsGsonMapBacked); return ok }, // Java: Multimap.class
		&logsGson_multimapTypeAdapter{}).
	Create()

// logsGson_ableTypeAdapter is the anonymous TypeAdapter<Able>.
type logsGson_ableTypeAdapter struct{}

// Java: new TypeAdapter<Able>() {...}.read(JsonReader)
func (*logsGson_ableTypeAdapter) Read(in *jutil.JsonReader) any {
	panic(jutil.NewAssertionErrorNoMessage())
}

var logsGsonNullLiteral = "null"

func writeNestedLogsGson(out *jutil.JsonWriter, v any) {
	if _, ok := v.(*jutil.JsonNull); ok || jutil.IsNil(v) {
		out.JsonValue(&logsGsonNullLiteral)
		return
	}
	LogsGson_GSON.ToJsonWriter(v, out)
}

// Java: new TypeAdapter<Able>() {...}.write(JsonWriter, Able)
func (*logsGson_ableTypeAdapter) Write(out *jutil.JsonWriter, value any) {
	writeNestedLogsGson(out, value.(LogsGson_Able).ToLogsGson())
}

// logsGson_multimapTypeAdapter is the anonymous TypeAdapter<Multimap<?, ?>>.
type logsGson_multimapTypeAdapter struct{}

// Java: new TypeAdapter<Multimap<?, ?>>() {...}.read(JsonReader)
func (*logsGson_multimapTypeAdapter) Read(in *jutil.JsonReader) any {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: new TypeAdapter<Multimap<?, ?>>() {...}.write(JsonWriter, Multimap<?, ?>)
func (*logsGson_multimapTypeAdapter) Write(out *jutil.JsonWriter, value any) {
	writeNestedLogsGson(out, value.(logsGsonMapBacked).AsMapObject())
}
