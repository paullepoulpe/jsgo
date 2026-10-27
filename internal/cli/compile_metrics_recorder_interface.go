// Port of com/google/javascript/jscomp/CompileMetricsRecorderInterface.java @ closure-compiler v20261021

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

package cli

import (
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

// CompileMetricsRecorderInterface is used by AbstractCompilerRunner to
// report metrics about the compilation action it performs.
type CompileMetricsRecorderInterface interface {
	// Java: CompileMetricsRecorderInterface.recordActionStart()
	RecordActionStart()

	// Java: CompileMetricsRecorderInterface.recordActionName(String)
	RecordActionName(actionName string)

	// Java: CompileMetricsRecorderInterface.recordStartState(AbstractCompiler)
	RecordStartState(compiler jscomp.AbstractCompiler)

	// Java: CompileMetricsRecorderInterface.recordResultMetrics(AbstractCompiler, Result)
	RecordResultMetrics(compiler jscomp.AbstractCompiler, result *jscomp.Result)
}
