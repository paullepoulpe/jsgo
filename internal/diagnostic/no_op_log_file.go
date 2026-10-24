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
