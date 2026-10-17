// Port of com/google/javascript/jscomp/diagnostic/WritingLogFile.java @ closure-compiler v20260415

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

// WritingLogFile is an implementation that adapts a BufferedWriter and
// writes to a real file.
//
// This class is designed to be super-sourceable by J2CL, since JS
// environments often don't have file APIs. Doing so enables logging
// statements to be checked-in.
//
// Java: final class WritingLogFile extends LogFile
type WritingLogFile struct {
	LogFileBase
	writer *jutil.BufferedWriter
}

// WritingLogFile_Create opens file for appending, creating it and its
// directory if needed.
//
// Java: WritingLogFile.create(Path)
func WritingLogFile_Create(file *jutil.Path) LogFile {
	var result LogFile
	jutil.Try(func() {
		dir := file.GetParent()
		jutil.Files_CreateDirectories(dir)
		// Allow the logged string to contain invalid Unicode (replace invalid character sequences
		// but don't throw). This is useful in cases where the source code contains invalid Unicode
		// (e.g., inside strings).
		// Java: UTF_8.newEncoder() with onMalformedInput(REPLACE) and
		// onUnmappableCharacter(REPLACE), which is jutil.OutputStreamWriter's encoding.
		result = NewWritingLogFile(
			jutil.NewBufferedWriter(
				jutil.NewOutputStreamWriter(
					jutil.Files_NewOutputStream(
						file,
						jutil.StandardOpenOption_CREATE,
						jutil.StandardOpenOption_APPEND,
						jutil.StandardOpenOption_WRITE))))
	}, func(e jutil.IIOException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	return result
}

// NewWritingLogFile wraps writer.
//
// Java: WritingLogFile(BufferedWriter)
func NewWritingLogFile(writer *jutil.BufferedWriter) *WritingLogFile {
	w := &WritingLogFile{}
	w.InitLogFileBase(w)
	w.writer = writer
	return w
}

// Java: WritingLogFile.log(Object)
func (w *WritingLogFile) Log(value any) LogFile {
	if jutil.IsNil(value) {
		// The JVM's message for value.toString() on null.
		panic(jutil.NewHelpfulNullPointerException(`Cannot invoke "Object.toString()" because "value" is null`))
	}
	return w.LogInternal(jutil.ToJavaString(value))
}

// Java: WritingLogFile.log(String)
func (w *WritingLogFile) LogString(value string) LogFile {
	return w.LogInternal(value)
}

// Java: WritingLogFile.log(Supplier<String>)
func (w *WritingLogFile) LogSupplier(value func() string) LogFile {
	return w.LogInternal(value())
}

// Java: WritingLogFile.log(String, Object...)
func (w *WritingLogFile) LogFormat(template string, values ...any) LogFile {
	return w.LogInternal(jutil.StringFormat(template, values...))
}

// Java: WritingLogFile.logJson(Object)
func (w *WritingLogFile) LogJson(value any) LogFile {
	return w.LogInternal(LogsGson_ToJson(value))
}

// Java: WritingLogFile.logJson(Supplier<Object>)
func (w *WritingLogFile) LogJsonSupplier(value func() any) LogFile {
	return w.LogInternal(LogsGson_ToJson(value()))
}

// LogJsonProducer streams JSON to the file through a JsonWriter. Like
// Java's try-with-resources, closing the JsonWriter closes the file's
// writer.
//
// Java: WritingLogFile.logJson(StreamedJsonProducer)
func (w *WritingLogFile) LogJsonProducer(producer LogFile_StreamedJsonProducer) LogFile {
	jutil.Try(func() {
		// Java: try (JsonWriter writer = new JsonWriter(this.writer)) { producer.writeJson(writer); }
		writer := jutil.NewJsonWriter(w.writer)
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
	return w
}

// LogInternal appends value and a newline.
//
// Java: WritingLogFile.logInternal(String)
func (w *WritingLogFile) LogInternal(value string) LogFile {
	jutil.Try(func() {
		// It's fine to pass a fully rendered string because we know we're going to use it by the
		// time this method is called.
		w.writer.Append(value).AppendChar('\n')
	}, func(e jutil.IIOException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
	return w
}

// Java: WritingLogFile.close()
func (w *WritingLogFile) Close() {
	jutil.Try(func() {
		w.writer.Close()
	}, func(e jutil.IIOException) {
		panic(jutil.NewRuntimeExceptionWithCause(e))
	})
}

// Java: WritingLogFile.isLogging()
func (w *WritingLogFile) IsLogging() bool {
	return true
}
