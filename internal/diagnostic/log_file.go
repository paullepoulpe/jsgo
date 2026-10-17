// Port of com/google/javascript/jscomp/diagnostic/LogFile.java @ closure-compiler v20260415

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

// LogFile is a simple interface for writing to a human readable log file.
//
// This API is designed to be compatible with J2CL. In the future it may be
// worth implementing this in terms of Flogger
// (https://github.com/google/flogger); however, at the time of writing,
// Flogger was not J2CL compatible.
//
// Java: public abstract class LogFile implements AutoCloseable. Its
// subclasses embed LogFileBase (PORTING.md §4).
type LogFile interface {
	// Java: LogFile.log(Object)
	Log(value any) LogFile
	// Provides a disambiguating overload.
	//
	// Java: LogFile.log(String)
	LogString(value string) LogFile
	// Java: LogFile.log(Supplier<String>)
	LogSupplier(value func() string) LogFile
	// Java: LogFile.log(String, Object...)
	LogFormat(template string, values ...any) LogFile
	// Java: LogFile.logJson(Object)
	LogJson(value any) LogFile
	// Java: LogFile.logJson(Supplier<Object>)
	LogJsonSupplier(value func() any) LogFile
	// Java: LogFile.logJson(StreamedJsonProducer)
	LogJsonProducer(producer LogFile_StreamedJsonProducer) LogFile
	// IsLogging reports whether or not the results are being recorded.
	//
	// Java: LogFile.isLogging()
	IsLogging() bool
	// Java: LogFile.close()
	Close()
	GetLogFileBase() *LogFileBase
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

// LogFileBase holds LogFile's base implementation (LogFile has no state).
type LogFileBase struct {
	self LogFile
}

// LogFile_CreateOrReopen creates the log file (and its directory) if
// needed and opens it for appending.
//
// Java: LogFile.createOrReopen(Path)
func LogFile_CreateOrReopen(file *jutil.Path) LogFile {
	return WritingLogFile_Create(file)
}

// LogFile_CreateNoOp returns a LogFile that discards everything.
//
// Java: LogFile.createNoOp()
func LogFile_CreateNoOp() LogFile {
	return NewNoOpLogFile()
}

// InitLogFileBase is the Go form of the LogFile constructor for subclasses.
// All subclasses are package classes.
//
// Java: LogFile()
func (b *LogFileBase) InitLogFileBase(self LogFile) {
	b.self = self
}

// GetLogFileBase returns the base struct.
//
// Go-only: language-adapter — field access through the interface
// (PORTING.md §4).
func (b *LogFileBase) GetLogFileBase() *LogFileBase {
	return b
}
