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
// Port of com/google/javascript/jscomp/modules/GoogEsImports.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.GoogEsImports
type GoogEsImports struct {
	idHash_ atomic.Int32
}

func (this_ *GoogEsImports) BindSelf_(s any) {
	_ = s
}
func (this_ *GoogEsImports) IsJavaNil_() bool { return this_ == nil }
func (*GoogEsImports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.modules.GoogEsImports", "GoogEsImports"
}
func (this_ *GoogEsImports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GoogEsImports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.GoogEsImports", "com.google.javascript.jscomp.modules.GoogEsImports", func(v_ any) bool {
		return j2gort.InstanceOf[*GoogEsImports](v_)
	})
}

func StaticGet_GoogEsImports_GOOG_IMPORT_PREFIX() string { return "goog:" }
func (this_ *GoogEsImports) initFields_GoogEsImports() {
	_ = this_
}

// Java constructor: GoogEsImports()
func NewGoogEsImports() *GoogEsImports {
	this_ := &GoogEsImports{}
	InitGoogEsImports(this_)
	return this_
}

func InitGoogEsImports(this_ *GoogEsImports) {
	this_.initFields_GoogEsImports()
}

// Java static method: GoogEsImports.isGoogImportSpecifier(java.lang.String moduleSpecifier)
func GoogEsImports_IsGoogImportSpecifier(moduleSpecifier string) bool {
	_ = moduleSpecifier
	return javalang.StringStartsWith(j2gort.CheckNotNullString(moduleSpecifier), "goog:")
}

// Java static method: GoogEsImports.getClosureIdFromGoogImportSpecifier(java.lang.String moduleSpecifier)
func GoogEsImports_GetClosureIdFromGoogImportSpecifier(moduleSpecifier string) string {
	_ = moduleSpecifier
	return j2gort.JStrSubFrom(moduleSpecifier, j2gort.JStrLen("goog:"))
}
