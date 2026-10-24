// Port of com/google/javascript/jscomp/LintPassConfig.java @ closure-compiler v20261021

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

package jscomp

// LintPassConfig is a PassConfig for the standalone linter, which runs on a
// single file at a time. This runs a similar set of checks to what you would
// get when running the compiler with the lintChecks DiagnosticGroup enabled,
// but some of the lint checks depend on type information, which is not
// available when looking at a single file, so those are omitted here.
//
// Java: class LintPassConfig extends PassConfig.PassConfigDelegate
type LintPassConfig struct {
	PassConfig_PassConfigDelegateBase

	gatherModuleMetadataPass *PassFactory
	earlyLintChecks          *PassFactory
	variableReferenceCheck   *PassFactory
	lateLintChecks           *PassFactory
}

// Java: LintPassConfig(CompilerOptions)
func NewLintPassConfig(options *CompilerOptions) *LintPassConfig {
	c := &LintPassConfig{}
	c.self = c
	c.InitPassConfig_PassConfigDelegateBase(NewDefaultPassConfig(options))

	c.gatherModuleMetadataPass =
		PassFactory_Builder_().
			SetName(PassNames_GATHER_MODULE_METADATA).
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewGatherModuleMetadata(
						compiler,
						compiler.GetOptions().GetProcessCommonJSModules(),
						compiler.GetOptions().GetModuleResolutionMode())
				}).
			Build()

	c.earlyLintChecks =
		PassFactory_Builder_().
			SetName("earlyLintChecks").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCombinedCompilerPassOfList(
						compiler,
						[]NodeTraversal_Callback{
							NewCheckConstPrivateProperties(compiler),
							NewCheckConstantCaseNames(compiler),
							NewCheckDuplicateCase(compiler),
							NewCheckEmptyStatements(compiler),
							NewCheckEnums(compiler),
							NewCheckExtraRequires(compiler, c.options.GetUnusedImportsToRemove()),
							NewCheckGoogModuleTypeScriptName(compiler),
							NewCheckJSDocStyle(compiler),
							NewCheckJSDoc(compiler),
							NewCheckMissingSemicolon(compiler),
							NewCheckSuper(compiler),
							NewCheckPrimitiveAsObject(compiler),
							NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap()),
							NewCheckNullabilityModifiers(compiler),
							NewCheckProvidesSorted(CheckProvidesSorted_Mode_COLLECT_AND_REPORT),
							NewCheckRequiresSorted(CheckRequiresSorted_Mode_COLLECT_AND_REPORT),
							NewCheckSideEffects(
								compiler /* report= */, true /* protectSideEffectFreeCode= */, false),
							NewCheckTypeImportCodeReferences(compiler),
							NewCheckUnusedLabels(compiler),
							NewCheckUnusedPrivateProperties(compiler),
							NewCheckUselessBlocks(compiler),
							NewCheckVar(compiler)})
				}).
			Build()

	c.variableReferenceCheck =
		PassFactory_Builder_().
			SetName("variableReferenceCheck").
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(compiler AbstractCompiler) CompilerPass {
				return NewVariableReferenceCheck(compiler)
			}).
			Build()

	c.lateLintChecks =
		PassFactory_Builder_().
			SetName("lateLintChecks").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCombinedCompilerPassOfList(
						compiler,
						[]NodeTraversal_Callback{
							NewCheckInterfaces(compiler), NewCheckPrototypeProperties(compiler)})
				}).
			Build()
	return c
}

// Java: LintPassConfig.getChecks()
func (c *LintPassConfig) GetChecks() *PassListBuilder {
	passes := NewPassListBuilder(c.options)
	passes.MaybeAdd(c.gatherModuleMetadataPass)
	passes.MaybeAdd(c.earlyLintChecks)
	passes.MaybeAdd(c.variableReferenceCheck)
	passes.MaybeAdd(c.lateLintChecks)
	return passes
}

// Java: LintPassConfig.getOptimizations()
func (c *LintPassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}

// Java: LintPassConfig.getFinalizations()
func (c *LintPassConfig) GetFinalizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}
