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
// Port of com/google/javascript/jscomp/modules/ExportTrace.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.modules.ExportTrace
type ExportTrace struct {
	module UnresolvedModule `gson:"module"`
	exportName string `gson:"exportName"`
}

func (this_ *ExportTrace) BindSelf_(s any) {
	_ = s
}
func (this_ *ExportTrace) IsJavaNil_() bool { return this_ == nil }
func (*ExportTrace) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.modules.ExportTrace", "ExportTrace" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ExportTrace", "com.google.javascript.jscomp.modules.ExportTrace", func(v_ any) bool {
		return j2gort.InstanceOf[*ExportTrace](v_)
	})
}
func (*ExportTrace) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *ExportTrace) initFields_ExportTrace() {
	_ = this_
}

// Java constructor: ExportTrace(com.google.javascript.jscomp.modules.UnresolvedModule module, java.lang.String exportName)
func NewExportTrace(module UnresolvedModule, exportName string) *ExportTrace {
	this_ := &ExportTrace{}
	InitExportTrace(this_, module, exportName)
	return this_
}

func InitExportTrace(this_ *ExportTrace, module UnresolvedModule, exportName string) {
	_ = module
	_ = exportName
	this_.initFields_ExportTrace()
	j2gort.CheckNotNullIface(module)
	j2gort.CheckNotNullString(exportName)
	this_.module = module
	this_.exportName = exportName
}

// Java static method: ExportTrace.create(com.google.javascript.jscomp.modules.UnresolvedModule module, java.lang.String exportName)
func ExportTrace_Create(module UnresolvedModule, exportName string) *ExportTrace {
	_ = module
	_ = exportName
	return NewExportTrace(module, exportName)
}

func (this_ *ExportTrace) Module() UnresolvedModule { return this_.module }
func (this_ *ExportTrace) ExportName() string { return this_.exportName }
func (this_ *ExportTrace) Equals(other_ any) bool {
	o_, ok_ := other_.(*ExportTrace)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.module, o_.module) { return false }
	if this_.exportName != o_.exportName { return false }
	return true
}
func (this_ *ExportTrace) HashCode() int32 {
	var h_ int32 = 0
	{ hi_1 := this_.module; h_ = 31 * h_; if hi_1 != nil && !hi_1.IsJavaNil_() { h_ += hi_1.HashCode() } }
	{ hs_2 := this_.exportName; h_ = 31 * h_; if hs_2 != j2gort.NullString { h_ += javalang.StringHashCode(hs_2) } }
	return h_
}
func (this_ *ExportTrace) String() string {
	_ = fmt.Sprint
	return "ExportTrace[" + "module=" + j2gort.StrVal(this_.module) + ", " + "exportName=" + j2gort.StrVal(j2gort.BoxString(this_.exportName)) + "]"
}

