/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckExtraRequiresPassConfig.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.CheckExtraRequiresPassConfig
type CheckExtraRequiresPassConfig struct {
	PassConfig_PassConfigDelegateBase
}

func (this_ *CheckExtraRequiresPassConfig) BindSelf_(s any) {
	_ = s
	this_.PassConfig_PassConfigDelegateBase.BindSelf_(s)
}
func (this_ *CheckExtraRequiresPassConfig) IsJavaNil_() bool { return this_ == nil }
func (*CheckExtraRequiresPassConfig) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckExtraRequiresPassConfig", "CheckExtraRequiresPassConfig"
}
func (this_ *CheckExtraRequiresPassConfig) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PassConfig_PassConfigDelegateBase.IdentityHash_()
}
func (*CheckExtraRequiresPassConfig) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckExtraRequiresPassConfig", "com.google.javascript.jscomp.CheckExtraRequiresPassConfig", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckExtraRequiresPassConfig](v_)
	})
}
func (*CheckExtraRequiresPassConfig) Implements_jscomp_PassConfig_PassConfigDelegate() {}
func (*CheckExtraRequiresPassConfig) Implements_jscomp_PassConfig()                    {}

func (this_ *CheckExtraRequiresPassConfig) initFields_CheckExtraRequiresPassConfig() {
	_ = this_
}

// Java constructor: CheckExtraRequiresPassConfig(com.google.javascript.jscomp.CompilerOptions options)
func NewCheckExtraRequiresPassConfig(options *CompilerOptions) *CheckExtraRequiresPassConfig {
	this_ := &CheckExtraRequiresPassConfig{}
	this_.BindSelf_(this_)
	InitCheckExtraRequiresPassConfig(this_, options)
	return this_
}

func InitCheckExtraRequiresPassConfig(this_ *CheckExtraRequiresPassConfig, options *CompilerOptions) {
	_ = options
	InitPassConfig_PassConfigDelegate(&this_.PassConfig_PassConfigDelegateBase, NewDefaultPassConfig(options))
	this_.initFields_CheckExtraRequiresPassConfig()
}

// Java method: CheckExtraRequiresPassConfig.getChecks()
func (this_ *CheckExtraRequiresPassConfig) GetChecks() *PassListBuilder {
	var val_1 PassListBuilder
	InitPassListBuilder(&val_1, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_1
	_ = passes
	passes.MaybeAdd(PassFactory_S_Builder().SetName("removeUnusedImport").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_2)
		_ = compiler
		return NewCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(compiler, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(NewCheckExtraRequires(compiler, this_.GetPassConfigBase().options.GetUnusedImportsToRemove()))))
	})).Build())
	return passes
}

// Java method: CheckExtraRequiresPassConfig.getOptimizations()
func (this_ *CheckExtraRequiresPassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(this_.GetPassConfigBase().options)
}
