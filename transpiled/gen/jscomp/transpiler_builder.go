/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/transpile/TranspilerBuilder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.transpile.TranspilerBuilder
type TranspilerBuilder struct {
	transpiler Transpile_Transpiler `gson:"transpiler"`
	idHash_    atomic.Int32
}

func (this_ *TranspilerBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *TranspilerBuilder) IsJavaNil_() bool { return this_ == nil }
func (*TranspilerBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.TranspilerBuilder", "TranspilerBuilder"
}
func (this_ *TranspilerBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TranspilerBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.transpile.TranspilerBuilder", "com.google.javascript.jscomp.transpile.TranspilerBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspilerBuilder](v_)
	})
}

var static_TranspilerBuilder_TO_ES5 *TranspilerBuilder

func StaticGet_TranspilerBuilder_TO_ES5() *TranspilerBuilder {
	clinit_TranspilerBuilder()
	return static_TranspilerBuilder_TO_ES5
}
func StaticSet_TranspilerBuilder_TO_ES5(v *TranspilerBuilder) {
	clinit_TranspilerBuilder()
	static_TranspilerBuilder_TO_ES5 = v
}
func StaticGet_TranspilerBuilder_DEFAULT_CACHE_SPEC() string { return "maximumSize=10000" }

var clinitDone_TranspilerBuilder uint32
var clinitState_TranspilerBuilder j2gort.ClinitState

func clinit_TranspilerBuilder() {
	if atomic.LoadUint32(&clinitDone_TranspilerBuilder) != 0 {
		return
	}
	clinitSlow_TranspilerBuilder()
}

func clinitSlow_TranspilerBuilder() {
	if !j2gort.EnterClinit(&clinitState_TranspilerBuilder, &clinitDone_TranspilerBuilder) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TranspilerBuilder, &clinitDone_TranspilerBuilder)
	static_TranspilerBuilder_TO_ES5 = NewTranspilerBuilder(NewBaseTranspiler(NewBaseTranspiler_CompilerSupplier(), "es6_runtime"))
}

func (this_ *TranspilerBuilder) initFields_TranspilerBuilder() {
	_ = this_
}

// Java constructor: TranspilerBuilder(com.google.javascript.jscomp.transpile.Transpiler transpiler)
func NewTranspilerBuilder(transpiler Transpile_Transpiler) *TranspilerBuilder {
	clinit_TranspilerBuilder()
	this_ := &TranspilerBuilder{}
	InitTranspilerBuilder(this_, transpiler)
	return this_
}

func InitTranspilerBuilder(this_ *TranspilerBuilder, transpiler Transpile_Transpiler) {
	_ = transpiler
	this_.initFields_TranspilerBuilder()
	this_.transpiler = transpiler
}

// Java static method: TranspilerBuilder.toEs5()
func TranspilerBuilder_ToEs5() *TranspilerBuilder {
	clinit_TranspilerBuilder()
	return StaticGet_TranspilerBuilder_TO_ES5()
}

// Java method: TranspilerBuilder.caching()
func (this_ *TranspilerBuilder) Caching_0() *TranspilerBuilder {
	return this_.Caching_1_string("maximumSize=10000")
}

// Java method: TranspilerBuilder.caching(java.lang.String spec)
func (this_ *TranspilerBuilder) Caching_1_string(spec string) *TranspilerBuilder {
	_ = spec
	return this_.Caching_1_Ext_com_google_common_cache_CacheBuilder(j2gort.Ext_com_google_common_cache_CacheBuilder_From_1_string(spec))
}

// Java method: TranspilerBuilder.caching(com.google.common.cache.CacheBuilder<java.lang.Object,? super com.google.javascript.jscomp.transpile.TranspileResult> builder)
func (this_ *TranspilerBuilder) Caching_1_Ext_com_google_common_cache_CacheBuilder(builder *j2gort.Ext_com_google_common_cache_CacheBuilder) *TranspilerBuilder {
	_ = builder
	return NewTranspilerBuilder(NewCachingTranspiler(this_.transpiler, builder))
}

// Java method: TranspilerBuilder.build()
func (this_ *TranspilerBuilder) Build() Transpile_Transpiler {
	return this_.transpiler
}
