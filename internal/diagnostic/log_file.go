// Port of com/google/javascript/jscomp/diagnostic/LogFile.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/diagnostic/WritingLogFile.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/diagnostic/NoOpLogFile.java @ closure-compiler v20261021

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

package diagnostic

import "github.com/paullepoulpe/jsgo/internal/jutil"

// LogFile is a handle for writing to a human readable log file.
//
// Its zero value (and LogFile_CreateNoOp()) has writer == nil and discards
// everything without evaluating suppliers or producers.
//
// Java: public abstract class LogFile implements AutoCloseable (collapsed with
// package-private subclasses WritingLogFile and NoOpLogFile).
type LogFile struct {
	writer *jutil.BufferedWriter
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
	var result LogFile
	jutil.Try(func() {
		dir := file.GetParent()
		jutil.Files_CreateDirectories(dir)
		// Allow the logged string to contain invalid Unicode (replace invalid character sequences
		// but don't throw). This is useful in cases where the source code contains invalid Unicode
		// (e.g., inside strings).
		// Java: UTF_8.newEncoder() with onMalformedInput(REPLACE) and
		// onUnmappableCharacter(REPLACE), which is jutil.OutputStreamWriter's encoding.
		result = LogFile{
			writer: jutil.NewBufferedWriter(
				jutil.NewOutputStreamWriter(
					jutil.Files_NewOutputStream(
						file,
						jutil.StandardOpenOption_CREATE,
						jutil.StandardOpenOption_APPEND,
						jutil.StandardOpenOption_WRITE))),
		}
	}, func(e jutil.IIOException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	return result
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
	if jutil.IsNil(value) {
		// The JVM's message for value.toString() on null.
		panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "Object.toString()" because "value" is null`))
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
	jutil.Try(func() {
		// Java: try (JsonWriter writer = new JsonWriter(this.writer)) { producer.writeJson(writer); }
		writer := jutil.NewJsonWriter(l.writer)
		closeWriter := func() {
			writer.Close() // also closes this.writer, as Java's JsonWriter.close() does
		}
		completed := false
		defer func() {
			if !completed {
				// The producer threw: close, suppressing close's exception.
				jutil.Try(closeWriter, func(jutil.Throwable) {})
			}
		}()
		producer.WriteJson(writer)
		completed = true
		closeWriter()
	}, func(ex jutil.IIOException) {
		panic(jutil.NewAssertionErrorWithCause(ex))
	})
	return l
}

// logInternal appends value and a newline.
//
// Java: WritingLogFile.logInternal(String)
func (l LogFile) logInternal(value string) LogFile {
	jutil.Try(func() {
		// It's fine to pass a fully rendered string because we know we're going to use it by the
		// time this method is called.
		l.writer.Append(value).AppendChar('\n')
	}, func(e jutil.IIOException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
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
	jutil.Try(func() {
		l.writer.Close()
	}, func(e jutil.IIOException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
}
