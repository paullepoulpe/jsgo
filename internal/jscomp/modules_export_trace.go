// Port of com/google/javascript/jscomp/modules/ExportTrace.java @ closure-compiler v20261015

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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ExportTrace is a simple class to keep track of which modules and exports
// have been visited when resolving exports. It is invalid to visit the same
// (module, name) pair more than once when resolving an export (invalid
// cycle).
//
// This is an AutoValue used for its hashCode / equals implementation and
// used in a Set for equality checks. So fields may appear to be "unused".
//
// Go-only: language-adapter — a comparable value struct: Go equality is the
// record's equality (the module by reference, the name by value), so
// ExportTraces can key jutil sets directly.
type ExportTrace struct {
	module     UnresolvedModule
	exportName string
}

// Java: ExportTrace(UnresolvedModule, String)
func NewExportTrace(module UnresolvedModule, exportName string) ExportTrace {
	if module == nil {
		panic(jutil.NewNullPointerExceptionMsg("module"))
	}
	return ExportTrace{module: module, exportName: exportName}
}

// Java: ExportTrace.module()
func (t ExportTrace) Module() UnresolvedModule { return t.module }

// Java: ExportTrace.exportName()
func (t ExportTrace) ExportName() string { return t.exportName }

// Java: ExportTrace.create(UnresolvedModule, String)
func ExportTrace_Create(module UnresolvedModule, exportName string) ExportTrace {
	return NewExportTrace(module, exportName)
}

// Equals is the record's equals.
func (t ExportTrace) Equals(o any) bool {
	other, ok := o.(ExportTrace)
	return ok && other == t
}

// HashCode is the record's hashCode.
func (t ExportTrace) HashCode() int32 {
	return 31*t.module.HashCode() + jutil.StringHashCode(t.exportName)
}

// String is the record's toString.
func (t ExportTrace) String() string {
	return "ExportTrace[module=" + jutil.ToJavaString(t.module) + ", exportName=" + t.exportName + "]"
}
