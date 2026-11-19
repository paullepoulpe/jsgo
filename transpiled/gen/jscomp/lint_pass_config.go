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
// Port of com/google/javascript/jscomp/LintPassConfig.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

var lambdaSingleton_lint_pass_config_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_lint_pass_config_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_lint_pass_config_3 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_lint_pass_config_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_1)
		_ = compiler
		return NewGatherModuleMetadata(compiler, compiler.GetOptions().GetProcessCommonJSModules(), compiler.GetOptions().GetModuleResolutionMode())
	})
	lambdaSingleton_lint_pass_config_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewVariableReferenceCheck(j2gort.AsIface[AbstractCompiler](mrp0_))
	})
	lambdaSingleton_lint_pass_config_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_3)
		_ = compiler
		return NewCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(compiler, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(NewCheckInterfaces(compiler), NewCheckPrototypeProperties(compiler))))
	})
}

// Java class: com.google.javascript.jscomp.LintPassConfig
type LintPassConfig struct {
	PassConfig_PassConfigDelegateBase
	gatherModuleMetadataPass PassFactory `gson:"gatherModuleMetadataPass"`
	earlyLintChecks          PassFactory `gson:"earlyLintChecks"`
	variableReferenceCheck   PassFactory `gson:"variableReferenceCheck"`
	lateLintChecks           PassFactory `gson:"lateLintChecks"`
}

func (this_ *LintPassConfig) BindSelf_(s any) {
	_ = s
	this_.PassConfig_PassConfigDelegateBase.BindSelf_(s)
}
func (this_ *LintPassConfig) IsJavaNil_() bool { return this_ == nil }
func (*LintPassConfig) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LintPassConfig", "LintPassConfig"
}
func (this_ *LintPassConfig) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PassConfig_PassConfigDelegateBase.IdentityHash_()
}
func (*LintPassConfig) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LintPassConfig", "com.google.javascript.jscomp.LintPassConfig", func(v_ any) bool {
		return j2gort.InstanceOf[*LintPassConfig](v_)
	})
}
func (*LintPassConfig) Implements_jscomp_PassConfig_PassConfigDelegate() {}
func (*LintPassConfig) Implements_jscomp_PassConfig()                    {}

func (this_ *LintPassConfig) initFields_LintPassConfig() {
	_ = this_
	this_.gatherModuleMetadataPass = PassFactory_S_Builder().SetName("gatherModuleMetadata").SetInternalFactory(lambdaSingleton_lint_pass_config_1).GetPassFactory_BuilderBase().Build()
	this_.earlyLintChecks = PassFactory_S_Builder().SetName("earlyLintChecks").SetInternalFactory(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		compiler := j2gort.AsIface[AbstractCompiler](lp_2)
		_ = compiler
		return NewCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(compiler, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_13(NewCheckConstPrivateProperties(compiler), NewCheckConstantCaseNames(compiler), NewCheckDuplicateCase(compiler), NewCheckEmptyStatements(compiler), NewCheckEnums(compiler), NewCheckExtraRequires(compiler, this_.GetPassConfigBase().options.GetUnusedImportsToRemove()), NewCheckGoogModuleTypeScriptName(compiler), NewCheckJSDocStyle(compiler), NewCheckJSDoc(compiler), NewCheckMissingSemicolon(compiler), NewCheckSuper(compiler), NewCheckPrimitiveAsObject(compiler), []any{NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap()), NewCheckNullabilityModifiers(compiler), NewCheckProvidesSorted(StaticGet_CheckProvidesSorted_Mode_COLLECT_AND_REPORT()), NewCheckRequiresSorted(StaticGet_CheckRequiresSorted_Mode_COLLECT_AND_REPORT()), NewCheckSideEffects(compiler, true, false), NewCheckTypeImportCodeReferences(compiler), NewCheckUnusedLabels(compiler), NewCheckUnusedPrivateProperties(compiler), NewCheckUselessBlocks(compiler), NewCheckVar(compiler)})))
	})).GetPassFactory_BuilderBase().Build()
	this_.variableReferenceCheck = PassFactory_S_Builder().SetName("variableReferenceCheck").SetRunInFixedPointLoop(true).SetInternalFactory(lambdaSingleton_lint_pass_config_2).GetPassFactory_BuilderBase().Build()
	this_.lateLintChecks = PassFactory_S_Builder().SetName("lateLintChecks").SetInternalFactory(lambdaSingleton_lint_pass_config_3).GetPassFactory_BuilderBase(

	// Java constructor: LintPassConfig(com.google.javascript.jscomp.CompilerOptions options)
	).Build()
}

func NewLintPassConfig(options *CompilerOptions) *LintPassConfig {
	this_ := &LintPassConfig{}
	this_.BindSelf_(this_)
	InitLintPassConfig(this_, options)
	return this_
}

func InitLintPassConfig(this_ *LintPassConfig, options *CompilerOptions) {
	_ = options
	InitPassConfig_PassConfigDelegate(&this_.PassConfig_PassConfigDelegateBase, NewDefaultPassConfig(options))
	this_.initFields_LintPassConfig()
}

// Java method: LintPassConfig.getChecks()
func (this_ *LintPassConfig) GetChecks() *PassListBuilder {
	var val_4 PassListBuilder
	InitPassListBuilder(&val_4, this_.GetPassConfigBase().options)
	var passes *PassListBuilder = &val_4
	_ = passes
	passes.MaybeAdd(this_.gatherModuleMetadataPass)
	passes.MaybeAdd(this_.earlyLintChecks)
	passes.MaybeAdd(this_.variableReferenceCheck)
	passes.MaybeAdd(this_.lateLintChecks)
	return passes
}

// Java method: LintPassConfig.getOptimizations()
func (this_ *LintPassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(this_.GetPassConfigBase().options)
}

// Java method: LintPassConfig.getFinalizations()
func (this_ *LintPassConfig) GetFinalizations() *PassListBuilder {
	return NewPassListBuilder(this_.GetPassConfigBase().options)
}
