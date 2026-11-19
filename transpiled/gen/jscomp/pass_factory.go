/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PassFactory.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)


var lambdaSingleton_pass_factory_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_pass_factory_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_pass_factory_3 CompilerPass
var lambdaSingleton_pass_factory_4 j2gort.Ext_java_util_function_Function
func init() {
	lambdaSingleton_pass_factory_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
	o := j2gort.CastPtr[CompilerOptions](lp_1)
	_ = o
	return j2gort.PtrToAny(StaticGet_PassFactory_PreconditionResult_SUCCESS())
})
	lambdaSingleton_pass_factory_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
	o := j2gort.CastPtr[CompilerOptions](lp_2)
	_ = o
	return j2gort.ValueOfBoolean(true)
})
	lambdaSingleton_pass_factory_3 = NewLambda_CompilerPass(func(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
})
	lambdaSingleton_pass_factory_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
	c := j2gort.AsIface[AbstractCompiler](lp_3)
	_ = c
	return lambdaSingleton_pass_factory_3
})
}
// Java polymorphic class interface: com.google.javascript.jscomp.PassFactory
type PassFactory interface {
	Implements_jscomp_PassFactory()
	IsJavaNil_() bool
	GetPassFactoryBase() *PassFactoryBase
	GetName() string
	GetCondition() j2gort.Ext_java_util_function_Function
	GetPreconditionCheck() j2gort.Ext_java_util_function_Function
	IsRunInFixedPointLoop() bool
	GetInternalFactory() j2gort.Ext_java_util_function_Function
	ToBuilder() PassFactory_Builder
	ValidatePreconditions(*CompilerOptions)
	Create(AbstractCompiler) CompilerPass
}
func ToIface_PassFactory[T interface {
	PassFactory
	comparable
}](v T) PassFactory {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.PassFactory
type PassFactoryBase struct {
	self_ PassFactory
	idHash_ atomic.Int32
}

func (this_ *PassFactoryBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(PassFactory)
}
func (this_ *PassFactoryBase) IsJavaNil_() bool { return this_ == nil }
func (*PassFactoryBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.PassFactory", "PassFactory" }
func (this_ *PassFactoryBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PassFactory", "com.google.javascript.jscomp.PassFactory", func(v_ any) bool {
		return j2gort.InstanceOf[PassFactory](v_)
	})
}
func (b *PassFactoryBase) GetPassFactoryBase() *PassFactoryBase { return b }
func (*PassFactoryBase) Implements_jscomp_PassFactory() {}

func (this_ *PassFactoryBase) initFields_PassFactory() {
	_ = this_
}

func InitPassFactory(this_ *PassFactoryBase) {
	this_.initFields_PassFactory()
}

// Java method: PassFactory.getName()
func (this_ *PassFactoryBase) GetName() string {
	panic("abstract or native method invoked")
}

// Java method: PassFactory.getCondition()
func (this_ *PassFactoryBase) GetCondition() j2gort.Ext_java_util_function_Function {
	panic("abstract or native method invoked")
}

// Java method: PassFactory.getPreconditionCheck()
func (this_ *PassFactoryBase) GetPreconditionCheck() j2gort.Ext_java_util_function_Function {
	panic("abstract or native method invoked")
}

// Java method: PassFactory.isRunInFixedPointLoop()
func (this_ *PassFactoryBase) IsRunInFixedPointLoop() bool {
	panic("abstract or native method invoked")
}

// Java method: PassFactory.getInternalFactory()
func (this_ *PassFactoryBase) GetInternalFactory() j2gort.Ext_java_util_function_Function {
	panic("abstract or native method invoked")
}

// Java method: PassFactory.toBuilder()
func (this_ *PassFactoryBase) ToBuilder() PassFactory_Builder {
	panic("abstract or native method invoked")
}

// Java static method: PassFactory.builder()
func PassFactory_S_Builder() PassFactory_Builder {
	return NewAutoValue_PassFactory_Builder().SetRunInFixedPointLoop(false).SetCondition(lambdaSingleton_pass_factory_2).SetPreconditionCheck(lambdaSingleton_pass_factory_1)
}

// Java method: PassFactory.validatePreconditions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *PassFactoryBase) ValidatePreconditions(options *CompilerOptions) {
	_ = options
	var canRun *PassFactory_PreconditionResult = j2gort.CastPtr[PassFactory_PreconditionResult](this_.self_.GetPreconditionCheck().Apply(j2gort.PtrToAny(options)))
	_ = canRun
	if (!(canRun.success)) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("Precondition for pass %s failed: %s", []any{j2gort.BoxString(this_.self_.GetName()), j2gort.BoxString(canRun.Message())})))))
	}
}

// Java static method: PassFactory.createEmptyPass(java.lang.String name)
func PassFactory_CreateEmptyPass(name string) PassFactory {
	_ = name
	return PassFactory_S_Builder().SetName(name).SetInternalFactory(lambdaSingleton_pass_factory_4).Build()
}

// Java method: PassFactory.create(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PassFactoryBase) Create(compiler AbstractCompiler) CompilerPass {
	_ = compiler
	return j2gort.AsIface[CompilerPass](this_.self_.GetInternalFactory().Apply(compiler))
}

// Java record: com.google.javascript.jscomp.PassFactory.PreconditionResult
type PassFactory_PreconditionResult struct {
	message string `gson:"message"`
	success bool `gson:"success"`
}

func (this_ *PassFactory_PreconditionResult) BindSelf_(s any) {
	_ = s
}
func (this_ *PassFactory_PreconditionResult) IsJavaNil_() bool { return this_ == nil }
func (*PassFactory_PreconditionResult) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.PassFactory$PreconditionResult", "PreconditionResult" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PassFactory.PreconditionResult", "com.google.javascript.jscomp.PassFactory$PreconditionResult", func(v_ any) bool {
		return j2gort.InstanceOf[*PassFactory_PreconditionResult](v_)
	})
}
func (*PassFactory_PreconditionResult) Implements_j2gort_Ext_java_lang_Record() {}

var static_PassFactory_PreconditionResult_SUCCESS *PassFactory_PreconditionResult
func StaticGet_PassFactory_PreconditionResult_SUCCESS() *PassFactory_PreconditionResult {
	clinit_PassFactory_PreconditionResult()
	return static_PassFactory_PreconditionResult_SUCCESS
}
func StaticSet_PassFactory_PreconditionResult_SUCCESS(v *PassFactory_PreconditionResult) {
	clinit_PassFactory_PreconditionResult()
	static_PassFactory_PreconditionResult_SUCCESS = v
}
var clinitDone_PassFactory_PreconditionResult uint32
var clinitState_PassFactory_PreconditionResult j2gort.ClinitState
func clinit_PassFactory_PreconditionResult() {
	if atomic.LoadUint32(&clinitDone_PassFactory_PreconditionResult) != 0 {
		return
	}
	clinitSlow_PassFactory_PreconditionResult()
}

func clinitSlow_PassFactory_PreconditionResult() {
	if !j2gort.EnterClinit(&clinitState_PassFactory_PreconditionResult, &clinitDone_PassFactory_PreconditionResult) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PassFactory_PreconditionResult, &clinitDone_PassFactory_PreconditionResult)
	static_PassFactory_PreconditionResult_SUCCESS = NewPassFactory_PreconditionResult(true, j2gort.NullString)
}

func (this_ *PassFactory_PreconditionResult) initFields_PassFactory_PreconditionResult() {
	_ = this_
}

// Java constructor: PreconditionResult(boolean success, java.lang.String message)
func NewPassFactory_PreconditionResult(success bool, message string) *PassFactory_PreconditionResult {
	clinit_PassFactory_PreconditionResult()
	this_ := &PassFactory_PreconditionResult{}
	InitPassFactory_PreconditionResult(this_, success, message)
	return this_
}

func InitPassFactory_PreconditionResult(this_ *PassFactory_PreconditionResult, success bool, message string) {
	_ = success
	_ = message
	this_.initFields_PassFactory_PreconditionResult()
	this_.success = success
	this_.message = message
}

func (this_ *PassFactory_PreconditionResult) Success() bool { return this_.success }
func (this_ *PassFactory_PreconditionResult) Message() string { return this_.message }
func (this_ *PassFactory_PreconditionResult) Equals(other_ any) bool {
	o_, ok_ := other_.(*PassFactory_PreconditionResult)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.success != o_.success { return false }
	if this_.message != o_.message { return false }
	return true
}
func (this_ *PassFactory_PreconditionResult) HashCode() int32 {
	var h_ int32 = 0
	if this_.success { h_ = 31*h_ + 1231 } else { h_ = 31*h_ + 1237 }
	{ hs_4 := this_.message; h_ = 31 * h_; if hs_4 != j2gort.NullString { h_ += javalang.StringHashCode(hs_4) } }
	return h_
}
func (this_ *PassFactory_PreconditionResult) String() string {
	_ = fmt.Sprint
	return "PreconditionResult[" + "success=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.success)) + ", " + "message=" + j2gort.StrVal(j2gort.BoxString(this_.message)) + "]"
}

// Java polymorphic class interface: com.google.javascript.jscomp.PassFactory.Builder
type PassFactory_Builder interface {
	Implements_jscomp_PassFactory_Builder()
	IsJavaNil_() bool
	GetPassFactory_BuilderBase() *PassFactory_BuilderBase
	SetName(string) PassFactory_Builder
	SetRunInFixedPointLoop(bool) PassFactory_Builder
	SetCondition(j2gort.Ext_java_util_function_Function) PassFactory_Builder
	SetPreconditionCheck(j2gort.Ext_java_util_function_Function) PassFactory_Builder
	SetInternalFactory(j2gort.Ext_java_util_function_Function) PassFactory_Builder
	AutoBuild() PassFactory
	Build() PassFactory
}
func ToIface_PassFactory_Builder[T interface {
	PassFactory_Builder
	comparable
}](v T) PassFactory_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.PassFactory.Builder
type PassFactory_BuilderBase struct {
	self_ PassFactory_Builder
	idHash_ atomic.Int32
}

func (this_ *PassFactory_BuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(PassFactory_Builder)
}
func (this_ *PassFactory_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*PassFactory_BuilderBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.PassFactory$Builder", "Builder" }
func (this_ *PassFactory_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PassFactory.Builder", "com.google.javascript.jscomp.PassFactory$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[PassFactory_Builder](v_)
	})
}
func (b *PassFactory_BuilderBase) GetPassFactory_BuilderBase() *PassFactory_BuilderBase { return b }
func (*PassFactory_BuilderBase) Implements_jscomp_PassFactory_Builder() {}

func (this_ *PassFactory_BuilderBase) initFields_PassFactory_Builder() {
	_ = this_
}

func InitPassFactory_Builder(this_ *PassFactory_BuilderBase) {
	this_.initFields_PassFactory_Builder()
}

// Java method: Builder.setName(java.lang.String x)
func (this_ *PassFactory_BuilderBase) SetName(x string) PassFactory_Builder {
	_ = x
	panic("abstract or native method invoked")
}

// Java method: Builder.setRunInFixedPointLoop(boolean b)
func (this_ *PassFactory_BuilderBase) SetRunInFixedPointLoop(b bool) PassFactory_Builder {
	_ = b
	panic("abstract or native method invoked")
}

// Java method: Builder.setCondition(java.util.function.Function<com.google.javascript.jscomp.CompilerOptions,java.lang.Boolean> cond)
func (this_ *PassFactory_BuilderBase) SetCondition(cond j2gort.Ext_java_util_function_Function) PassFactory_Builder {
	_ = cond
	panic("abstract or native method invoked")
}

// Java method: Builder.setPreconditionCheck(java.util.function.Function<com.google.javascript.jscomp.CompilerOptions,com.google.javascript.jscomp.PassFactory.PreconditionResult> assertion)
func (this_ *PassFactory_BuilderBase) SetPreconditionCheck(assertion j2gort.Ext_java_util_function_Function) PassFactory_Builder {
	_ = assertion
	panic("abstract or native method invoked")
}

// Java method: Builder.setInternalFactory(java.util.function.Function<com.google.javascript.jscomp.AbstractCompiler,? extends com.google.javascript.jscomp.CompilerPass> x)
func (this_ *PassFactory_BuilderBase) SetInternalFactory(x j2gort.Ext_java_util_function_Function) PassFactory_Builder {
	_ = x
	panic("abstract or native method invoked")
}

// Java method: Builder.autoBuild()
func (this_ *PassFactory_BuilderBase) AutoBuild() PassFactory {
	panic("abstract or native method invoked")
}

// Java method: Builder.build()
func (this_ *PassFactory_BuilderBase) Build() PassFactory {
	var result PassFactory = this_.self_.AutoBuild()
	_ = result
	if !(!((j2gort.CheckNotNullString(result.GetName()) == ""))) { j2gort.CheckState(false) }
	return result
}

