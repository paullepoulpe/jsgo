// Port of com/google/javascript/jscomp/CheckExtraRequiresPassConfig.java @ closure-compiler v20261015

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

// CheckExtraRequiresPassConfig is a PassConfig to only run the
// CheckExtraRequires linter rule. This allows user to select which
// goog.require imports to be removed, provided by requiresToRemove
//
// Java: final class CheckExtraRequiresPassConfig extends PassConfig.PassConfigDelegate
type CheckExtraRequiresPassConfig struct {
	PassConfig_PassConfigDelegateBase
}

// Java: CheckExtraRequiresPassConfig(CompilerOptions)
func NewCheckExtraRequiresPassConfig(options *CompilerOptions) *CheckExtraRequiresPassConfig {
	c := &CheckExtraRequiresPassConfig{}
	c.self = c
	c.InitPassConfig_PassConfigDelegateBase(NewDefaultPassConfig(options))
	return c
}

// Java: CheckExtraRequiresPassConfig.getChecks()
func (c *CheckExtraRequiresPassConfig) GetChecks() *PassListBuilder {
	passes := NewPassListBuilder(c.options)
	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName("removeUnusedImport").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					return NewCombinedCompilerPassOfList(
						compiler,
						[]NodeTraversal_Callback{
							NewCheckExtraRequires(compiler, c.options.GetUnusedImportsToRemove())})
				}).
			Build())
	return passes
}

// Java: CheckExtraRequiresPassConfig.getOptimizations()
func (c *CheckExtraRequiresPassConfig) GetOptimizations() *PassListBuilder {
	return NewPassListBuilder(c.options)
}
