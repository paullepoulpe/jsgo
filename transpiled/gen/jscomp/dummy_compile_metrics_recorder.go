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
// Port of com/google/javascript/jscomp/DummyCompileMetricsRecorder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.DummyCompileMetricsRecorder
type DummyCompileMetricsRecorder struct {
	idHash_ atomic.Int32
}

func (this_ *DummyCompileMetricsRecorder) BindSelf_(s any) {
	_ = s
}
func (this_ *DummyCompileMetricsRecorder) IsJavaNil_() bool { return this_ == nil }
func (*DummyCompileMetricsRecorder) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.DummyCompileMetricsRecorder", "DummyCompileMetricsRecorder" }
func (this_ *DummyCompileMetricsRecorder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DummyCompileMetricsRecorder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DummyCompileMetricsRecorder", "com.google.javascript.jscomp.DummyCompileMetricsRecorder", func(v_ any) bool {
		return j2gort.InstanceOf[*DummyCompileMetricsRecorder](v_)
	})
}
func (*DummyCompileMetricsRecorder) Implements_jscomp_CompileMetricsRecorderInterface() {}

func (this_ *DummyCompileMetricsRecorder) initFields_DummyCompileMetricsRecorder() {
	_ = this_
}

// Java constructor: DummyCompileMetricsRecorder()
func NewDummyCompileMetricsRecorder() *DummyCompileMetricsRecorder {
	this_ := &DummyCompileMetricsRecorder{}
	InitDummyCompileMetricsRecorder(this_)
	return this_
}

func InitDummyCompileMetricsRecorder(this_ *DummyCompileMetricsRecorder) {
	this_.initFields_DummyCompileMetricsRecorder()
}

// Java method: DummyCompileMetricsRecorder.recordActionStart()
func (this_ *DummyCompileMetricsRecorder) RecordActionStart() {
}

// Java method: DummyCompileMetricsRecorder.recordActionName(java.lang.String actionName)
func (this_ *DummyCompileMetricsRecorder) RecordActionName(actionName string) {
	_ = actionName
}

// Java method: DummyCompileMetricsRecorder.recordStartState(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *DummyCompileMetricsRecorder) RecordStartState(compiler AbstractCompiler) {
	_ = compiler
}

// Java method: DummyCompileMetricsRecorder.recordResultMetrics(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.Result result)
func (this_ *DummyCompileMetricsRecorder) RecordResultMetrics(compiler AbstractCompiler, result *Result) {
	_ = compiler
	_ = result
}

