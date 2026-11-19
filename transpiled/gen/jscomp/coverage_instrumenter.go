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
// Port of com/google/javascript/jscomp/bundle/CoverageInstrumenter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.bundle.CoverageInstrumenter
type CoverageInstrumenter struct {
	CompilerBasedTransformerBase
}

func (this_ *CoverageInstrumenter) BindSelf_(s any) {
	_ = s
	this_.CompilerBasedTransformerBase.BindSelf_(s)
}
func (this_ *CoverageInstrumenter) IsJavaNil_() bool { return this_ == nil }
func (*CoverageInstrumenter) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.bundle.CoverageInstrumenter", "CoverageInstrumenter" }
func (this_ *CoverageInstrumenter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CompilerBasedTransformerBase.IdentityHash_()
}
func (*CoverageInstrumenter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.CoverageInstrumenter", "com.google.javascript.jscomp.bundle.CoverageInstrumenter", func(v_ any) bool {
		return j2gort.InstanceOf[*CoverageInstrumenter](v_)
	})
}
func (*CoverageInstrumenter) Implements_jscomp_CompilerBasedTransformer() {}
func (*CoverageInstrumenter) Implements_jscomp_Source_Transformer() {}

var static_CoverageInstrumenter_INSTRUMENTER *CoverageInstrumenter
func StaticGet_CoverageInstrumenter_INSTRUMENTER() *CoverageInstrumenter {
	clinit_CoverageInstrumenter()
	return static_CoverageInstrumenter_INSTRUMENTER
}
func StaticSet_CoverageInstrumenter_INSTRUMENTER(v *CoverageInstrumenter) {
	clinit_CoverageInstrumenter()
	static_CoverageInstrumenter_INSTRUMENTER = v
}
var clinitDone_CoverageInstrumenter uint32
var clinitState_CoverageInstrumenter j2gort.ClinitState
func clinit_CoverageInstrumenter() {
	if atomic.LoadUint32(&clinitDone_CoverageInstrumenter) != 0 {
		return
	}
	clinitSlow_CoverageInstrumenter()
}

func clinitSlow_CoverageInstrumenter() {
	if !j2gort.EnterClinit(&clinitState_CoverageInstrumenter, &clinitDone_CoverageInstrumenter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CoverageInstrumenter, &clinitDone_CoverageInstrumenter)
	static_CoverageInstrumenter_INSTRUMENTER = NewCoverageInstrumenter(ToIface_CompilerBasedTransformer_CompilerSupplier(CoverageInstrumenter_S_CompilerSupplier()))
}

func (this_ *CoverageInstrumenter) initFields_CoverageInstrumenter() {
	_ = this_
}

// Java constructor: CoverageInstrumenter(com.google.javascript.jscomp.bundle.CompilerBasedTransformer.CompilerSupplier compilerSupplier)
func NewCoverageInstrumenter(compilerSupplier CompilerBasedTransformer_CompilerSupplier) *CoverageInstrumenter {
	clinit_CoverageInstrumenter()
	this_ := &CoverageInstrumenter{}
	this_.BindSelf_(this_)
	InitCoverageInstrumenter(this_, compilerSupplier)
	return this_
}

func InitCoverageInstrumenter(this_ *CoverageInstrumenter, compilerSupplier CompilerBasedTransformer_CompilerSupplier) {
	_ = compilerSupplier
	InitCompilerBasedTransformer(&this_.CompilerBasedTransformerBase, compilerSupplier)
	this_.initFields_CoverageInstrumenter()
}

// Java method: CoverageInstrumenter.getRuntime()
func (this_ *CoverageInstrumenter) GetRuntime() *j2gort.Ext_java_util_Optional {
	return j2gort.Ext_java_util_Optional_Empty()
}

// Java method: CoverageInstrumenter.getTranformationName()
func (this_ *CoverageInstrumenter) GetTranformationName() string {
	return "Coverage Instrumentation"
}

// Java static method: CoverageInstrumenter.compilerSupplier()
func CoverageInstrumenter_S_CompilerSupplier() *CoverageInstrumenter_CompilerSupplier {
	clinit_CoverageInstrumenter()
	return NewCoverageInstrumenter_CompilerSupplier()
}

// Java class: com.google.javascript.jscomp.bundle.CoverageInstrumenter.CompilerSupplier
type CoverageInstrumenter_CompilerSupplier struct {
	CompilerBasedTransformer_CompilerSupplierBase
}

func (this_ *CoverageInstrumenter_CompilerSupplier) BindSelf_(s any) {
	_ = s
	this_.CompilerBasedTransformer_CompilerSupplierBase.BindSelf_(s)
}
func (this_ *CoverageInstrumenter_CompilerSupplier) IsJavaNil_() bool { return this_ == nil }
func (*CoverageInstrumenter_CompilerSupplier) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.bundle.CoverageInstrumenter$CompilerSupplier", "CompilerSupplier" }
func (this_ *CoverageInstrumenter_CompilerSupplier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CompilerBasedTransformer_CompilerSupplierBase.IdentityHash_()
}
func (*CoverageInstrumenter_CompilerSupplier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.CoverageInstrumenter.CompilerSupplier", "com.google.javascript.jscomp.bundle.CoverageInstrumenter$CompilerSupplier", func(v_ any) bool {
		return j2gort.InstanceOf[*CoverageInstrumenter_CompilerSupplier](v_)
	})
}
func (*CoverageInstrumenter_CompilerSupplier) Implements_jscomp_CompilerBasedTransformer_CompilerSupplier() {}

func (this_ *CoverageInstrumenter_CompilerSupplier) initFields_CoverageInstrumenter_CompilerSupplier() {
	_ = this_
}

// Java constructor: CompilerSupplier()
func NewCoverageInstrumenter_CompilerSupplier() *CoverageInstrumenter_CompilerSupplier {
	this_ := &CoverageInstrumenter_CompilerSupplier{}
	this_.BindSelf_(this_)
	InitCoverageInstrumenter_CompilerSupplier(this_)
	return this_
}

func InitCoverageInstrumenter_CompilerSupplier(this_ *CoverageInstrumenter_CompilerSupplier) {
	InitCompilerBasedTransformer_CompilerSupplier(&this_.CompilerBasedTransformer_CompilerSupplierBase)
	this_.initFields_CoverageInstrumenter_CompilerSupplier()
}

// Java method: CompilerSupplier.setOptions(com.google.javascript.jscomp.CompilerOptions options)
func (this_ *CoverageInstrumenter_CompilerSupplier) SetOptions(options *CompilerOptions) {
	_ = options
	options.SetCoalesceVariableNames(false)
	options.SetLanguageOut(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT5())
	options.SetLanguageIn(StaticGet_CompilerOptions_LanguageMode_ECMASCRIPT_NEXT())
	options.SetStrictModeInput(false)
	options.SetSourceMapOutputPath("/dev/null")
	options.SetVariableRenaming(StaticGet_VariableRenamingPolicy_OFF())
	options.SetInstrumentForCoverageOption(StaticGet_CompilerOptions_InstrumentOption_LINE_ONLY())
	options.SetInstrumentForCoverageOnly(true)
}

// Java method: CompilerSupplier.transformed(com.google.javascript.jscomp.Result result)
func (this_ *CoverageInstrumenter_CompilerSupplier) Transformed(result *Result) bool {
	_ = result
	return true
}

