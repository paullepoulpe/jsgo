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

package linter

import (
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

// LintPassConfig is a PassConfig for the standalone linter, which runs on a
// single file at a time. This runs a similar set of checks to what you would
// get when running the compiler with the lintChecks DiagnosticGroup enabled,
// but some of the lint checks depend on type information, which is not
// available when looking at a single file, so those are omitted here.
//
// Java: class LintPassConfig extends PassConfig.PassConfigDelegate
type LintPassConfig struct {
	jscomp.PassConfigDelegate

	gatherModuleMetadataPass *jscomp.PassFactory
	earlyLintChecks          *jscomp.PassFactory
	variableReferenceCheck   *jscomp.PassFactory
	lateLintChecks           *jscomp.PassFactory
}

// Java: LintPassConfig(CompilerOptions)
func NewLintPassConfig(options *jscomp.CompilerOptions) *LintPassConfig {
	c := &LintPassConfig{
		PassConfigDelegate: jscomp.NewPassConfigDelegate(jscomp.NewDefaultPassConfig(options)),
	}

	c.gatherModuleMetadataPass =
		jscomp.PassFactory_Builder_().
			SetName(jscomp.PassNames_GATHER_MODULE_METADATA).
			SetInternalFactory(
				func(compiler jscomp.AbstractCompiler) jscomp.CompilerPass {
					return jscomp.NewGatherModuleMetadata(
						compiler,
						compiler.GetOptions().GetProcessCommonJSModules(),
						compiler.GetOptions().GetModuleResolutionMode())
				}).
			Build()

	c.earlyLintChecks =
		jscomp.PassFactory_Builder_().
			SetName("earlyLintChecks").
			SetInternalFactory(
				func(compiler jscomp.AbstractCompiler) jscomp.CompilerPass {
					return jscomp.NewCombinedCompilerPassOfList(
						compiler,
						[]jscomp.NodeTraversal_Callback{
							jscomp.NewCheckConstPrivateProperties(compiler),
							jscomp.NewCheckConstantCaseNames(compiler),
							jscomp.NewCheckDuplicateCase(compiler),
							jscomp.NewCheckEmptyStatements(compiler),
							jscomp.NewCheckEnums(compiler),
							jscomp.NewCheckExtraRequires(compiler, c.Options().GetUnusedImportsToRemove()),
							jscomp.NewCheckGoogModuleTypeScriptName(compiler),
							jscomp.NewCheckJSDocStyle(compiler),
							jscomp.NewCheckJSDoc(compiler),
							jscomp.NewCheckMissingSemicolon(compiler),
							jscomp.NewCheckSuper(compiler),
							jscomp.NewCheckPrimitiveAsObject(compiler),
							jscomp.NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap()),
							jscomp.NewCheckNullabilityModifiers(compiler),
							jscomp.NewCheckProvidesSorted(jscomp.CheckProvidesSorted_Mode_COLLECT_AND_REPORT),
							jscomp.NewCheckRequiresSorted(jscomp.CheckRequiresSorted_Mode_COLLECT_AND_REPORT),
							jscomp.NewCheckSideEffects(
								compiler /* report= */, true /* protectSideEffectFreeCode= */, false),
							jscomp.NewCheckTypeImportCodeReferences(compiler),
							jscomp.NewCheckUnusedLabels(compiler),
							jscomp.NewCheckUnusedPrivateProperties(compiler),
							jscomp.NewCheckUselessBlocks(compiler),
							jscomp.NewCheckVar(compiler)})
				}).
			Build()

	c.variableReferenceCheck =
		jscomp.PassFactory_Builder_().
			SetName("variableReferenceCheck").
			SetRunInFixedPointLoop(true).
			SetInternalFactory(func(compiler jscomp.AbstractCompiler) jscomp.CompilerPass {
				return jscomp.NewVariableReferenceCheck(compiler)
			}).
			Build()

	c.lateLintChecks =
		jscomp.PassFactory_Builder_().
			SetName("lateLintChecks").
			SetInternalFactory(
				func(compiler jscomp.AbstractCompiler) jscomp.CompilerPass {
					return jscomp.NewCombinedCompilerPassOfList(
						compiler,
						[]jscomp.NodeTraversal_Callback{
							jscomp.NewCheckInterfaces(compiler), jscomp.NewCheckPrototypeProperties(compiler)})
				}).
			Build()
	return c
}

// Java: LintPassConfig.getChecks()
func (c *LintPassConfig) GetChecks() *jscomp.PassListBuilder {
	passes := jscomp.NewPassListBuilder(c.Options())
	passes.MaybeAdd(c.gatherModuleMetadataPass)
	passes.MaybeAdd(c.earlyLintChecks)
	passes.MaybeAdd(c.variableReferenceCheck)
	passes.MaybeAdd(c.lateLintChecks)
	return passes
}

// Java: LintPassConfig.getOptimizations()
func (c *LintPassConfig) GetOptimizations() *jscomp.PassListBuilder {
	return jscomp.NewPassListBuilder(c.Options())
}

// Java: LintPassConfig.getFinalizations()
func (c *LintPassConfig) GetFinalizations() *jscomp.PassListBuilder {
	return jscomp.NewPassListBuilder(c.Options())
}
