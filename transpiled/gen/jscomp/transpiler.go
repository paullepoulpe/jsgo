/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/bundle/Transpiler.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.bundle.Transpiler
type Bundle_Transpiler struct {
	CompilerBasedTransformerBase
	runtimeLibraryName string `gson:"runtimeLibraryName"`
}

func (this_ *Bundle_Transpiler) BindSelf_(s any) {
	_ = s
	this_.CompilerBasedTransformerBase.BindSelf_(s)
}
func (this_ *Bundle_Transpiler) IsJavaNil_() bool { return this_ == nil }
func (*Bundle_Transpiler) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.bundle.Transpiler", "Transpiler" }
func (this_ *Bundle_Transpiler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CompilerBasedTransformerBase.IdentityHash_()
}
func (*Bundle_Transpiler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.Transpiler", "com.google.javascript.jscomp.bundle.Transpiler", func(v_ any) bool {
		return j2gort.InstanceOf[*Bundle_Transpiler](v_)
	})
}
func (*Bundle_Transpiler) Implements_jscomp_CompilerBasedTransformer() {}
func (*Bundle_Transpiler) Implements_jscomp_Source_Transformer() {}

var static_Bundle_Transpiler_ES5_TRANSPILER *Bundle_Transpiler
func StaticGet_Bundle_Transpiler_ES5_TRANSPILER() *Bundle_Transpiler {
	clinit_Bundle_Transpiler()
	return static_Bundle_Transpiler_ES5_TRANSPILER
}
func StaticSet_Bundle_Transpiler_ES5_TRANSPILER(v *Bundle_Transpiler) {
	clinit_Bundle_Transpiler()
	static_Bundle_Transpiler_ES5_TRANSPILER = v
}
var static_Bundle_Transpiler_TO_ES5 *Bundle_Transpiler
func StaticGet_Bundle_Transpiler_TO_ES5() *Bundle_Transpiler {
	clinit_Bundle_Transpiler()
	return static_Bundle_Transpiler_TO_ES5
}
func StaticSet_Bundle_Transpiler_TO_ES5(v *Bundle_Transpiler) {
	clinit_Bundle_Transpiler()
	static_Bundle_Transpiler_TO_ES5 = v
}
var clinitDone_Bundle_Transpiler uint32
var clinitState_Bundle_Transpiler j2gort.ClinitState
func clinit_Bundle_Transpiler() {
	if atomic.LoadUint32(&clinitDone_Bundle_Transpiler) != 0 {
		return
	}
	clinitSlow_Bundle_Transpiler()
}

func clinitSlow_Bundle_Transpiler() {
	if !j2gort.EnterClinit(&clinitState_Bundle_Transpiler, &clinitDone_Bundle_Transpiler) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Bundle_Transpiler, &clinitDone_Bundle_Transpiler)
	static_Bundle_Transpiler_ES5_TRANSPILER = NewBundle_Transpiler(NewCompilerBasedTransformer_CompilerSupplier(), "es6_runtime")
	static_Bundle_Transpiler_TO_ES5 = NewBundle_Transpiler(NewCompilerBasedTransformer_CompilerSupplier(), "es6_runtime")
}

func (this_ *Bundle_Transpiler) initFields_Bundle_Transpiler() {
	_ = this_
	this_.runtimeLibraryName = j2gort.NullString
}

// Java constructor: Transpiler(com.google.javascript.jscomp.bundle.CompilerBasedTransformer.CompilerSupplier compilerSupplier, java.lang.String runtimeLibraryName)
func NewBundle_Transpiler(compilerSupplier CompilerBasedTransformer_CompilerSupplier, runtimeLibraryName string) *Bundle_Transpiler {
	clinit_Bundle_Transpiler()
	this_ := &Bundle_Transpiler{}
	this_.BindSelf_(this_)
	InitBundle_Transpiler(this_, compilerSupplier, runtimeLibraryName)
	return this_
}

func InitBundle_Transpiler(this_ *Bundle_Transpiler, compilerSupplier CompilerBasedTransformer_CompilerSupplier, runtimeLibraryName string) {
	_ = compilerSupplier
	_ = runtimeLibraryName
	InitCompilerBasedTransformer(&this_.CompilerBasedTransformerBase, compilerSupplier)
	this_.initFields_Bundle_Transpiler()
	this_.runtimeLibraryName = j2gort.CheckNotNullString(runtimeLibraryName)
}

// Java static method: Transpiler.compilerSupplier()
func Bundle_Transpiler_CompilerSupplier() CompilerBasedTransformer_CompilerSupplier {
	clinit_Bundle_Transpiler()
	return NewCompilerBasedTransformer_CompilerSupplier()
}

// Java method: Transpiler.getRuntime()
func (this_ *Bundle_Transpiler) GetRuntime() *j2gort.Ext_java_util_Optional {
	return j2gort.Ext_java_util_Optional_Of(j2gort.BoxString(this_.runtimeLibraryName))
}

// Java method: Transpiler.getTranformationName()
func (this_ *Bundle_Transpiler) GetTranformationName() string {
	return "Transpilation"
}

// Java static method: Transpiler.toEs5()
func Bundle_Transpiler_ToEs5() *Bundle_Transpiler {
	clinit_Bundle_Transpiler()
	return StaticGet_Bundle_Transpiler_TO_ES5()
}

