// Port of com/google/javascript/jscomp/DummyCompileMetricsRecorder.java @ closure-compiler v20261015

/*
 * Copyright 2021 The Closure Compiler Authors.
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

package jscomp

// DummyCompileMetricsRecorder is a do-nothing metrics recorder that
// AbstractCompilerRunner uses by default.
type DummyCompileMetricsRecorder struct{}

// Java: DummyCompileMetricsRecorder()
func NewDummyCompileMetricsRecorder() *DummyCompileMetricsRecorder {
	return &DummyCompileMetricsRecorder{}
}

// Java: DummyCompileMetricsRecorder.recordActionStart()
func (r *DummyCompileMetricsRecorder) RecordActionStart() {}

// Java: DummyCompileMetricsRecorder.recordActionName(String)
func (r *DummyCompileMetricsRecorder) RecordActionName(actionName string) {}

// Java: DummyCompileMetricsRecorder.recordStartState(AbstractCompiler)
func (r *DummyCompileMetricsRecorder) RecordStartState(compiler AbstractCompiler) {}

// Java: DummyCompileMetricsRecorder.recordResultMetrics(AbstractCompiler, Result)
func (r *DummyCompileMetricsRecorder) RecordResultMetrics(compiler AbstractCompiler, result *Result) {
}
