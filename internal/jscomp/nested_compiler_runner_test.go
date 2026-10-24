// Port of com/google/javascript/jscomp/NestedCompilerRunnerTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/NestedCompilerRunnerTest.java)

/*
 * Copyright 2025 The Closure Compiler Authors.
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

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NestedCompilerRunnerTest: tests for NestedCompilerRunner.
type NestedCompilerRunnerTest struct {
	t *testing.T

	original *Compiler
	options  *CompilerOptions
}

func (c *NestedCompilerRunnerTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *NestedCompilerRunnerTest) JUnitBefore()           { c.Setup() }
func (c *NestedCompilerRunnerTest) JUnitAfter()            {}

func TestNestedCompilerRunnerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &NestedCompilerRunnerTest{} })
}

// Java: NestedCompilerRunnerTest.setup()
func (c *NestedCompilerRunnerTest) Setup() {
	c.original = NewCompiler()
	c.options = NewCompilerOptions()
}

// Java: NestedCompilerRunnerTest.emptyInputListSucceeds()
func (c *NestedCompilerRunnerTest) TestEmptyInputListSucceeds() {
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, NewDefaultPassConfig(c.options)).
		Compile()

	c.verifyCompilationSucceeded()
}

// Java: NestedCompilerRunnerTest.singleEmptyInputSucceeds()
func (c *NestedCompilerRunnerTest) TestSingleEmptyInputSucceeds() {
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)
	script := c.createScript("")

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, NewDefaultPassConfig(c.options)).
		AddScript(script, "test").
		Compile()

	c.verifyCompilationSucceeded()
	AssertNode(script).IsEqualTo(c.parse(""))
}

// Java: NestedCompilerRunnerTest.canRunPeepholeOptimizations()
func (c *NestedCompilerRunnerTest) TestCanRunPeepholeOptimizations() {
	c.options.SetFoldConstants(true)
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)

	script := c.createScript("console.log(1 + 2)")

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, NewDefaultPassConfig(c.options)).
		AddScript(script, "test").
		Compile()

	c.verifyCompilationSucceeded()
	AssertNode(script).IsEqualTo(c.parse("console.log(3)"))
}

// Java: NestedCompilerRunnerTest.delegatesWarningsAndErrorsToOriginalCompiler()
func (c *NestedCompilerRunnerTest) TestDelegatesWarningsAndErrorsToOriginalCompiler() {
	c.options.SetComputeFunctionSideEffects(true)
	c.options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	c.options.SetWarningLevel(
		DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION, CheckLevel_ERROR)
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)

	script :=
		c.createScript(
			`class A {
  /** @nosideeffects */
  error() { throw new Error('A'); }
}
class B {
  error() { throw new Error('B'); }
}
console.log((Math.random() < 0.5 ? new A() : new B()).error());
`)

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, NewDefaultPassConfig(c.options)).
		AddScript(script, "test").
		Compile()

	AssertThat(c.original.GetWarnings()).IsEmpty()
	AssertThat(c.original.GetErrors()).
		ComparingElementsUsing(JSCompCorrespondences_DIAGNOSTIC_EQUALITY).
		ContainsExactly(PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION)
}

// Java: NestedCompilerRunnerTest.canDeleteUnusedVariables()
func (c *NestedCompilerRunnerTest) TestCanDeleteUnusedVariables() {
	CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(c.options)
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)
	c.options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_OFF)
	c.options.SetRemoveUnusedVariables(CompilerOptions_Reach_ALL)

	script :=
		c.createScript(
			`const used = 0;
const unused = 1;
console.log(used);
`)

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, NewDefaultPassConfig(c.options)).
		AddScript(script, "test").
		Compile()

	c.verifyCompilationSucceeded()
	AssertNode(script).IsEqualTo(c.parse("const used = 0; console.log(used);"))
}

// Java: NestedCompilerRunnerTest.transpiles_exponentialOperatorToES2015()
func (c *NestedCompilerRunnerTest) TestTranspiles_exponentialOperatorToES2015() {
	c.options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.options.SetPreventLibraryInjection(true)
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)

	// Use the exponential operator as a very trivial transpilation test case.
	script := c.createScript("const x = 3 ** 4;")

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, NewDefaultPassConfig(c.options)).
		AddScript(script, "test").
		Compile()

	c.verifyCompilationSucceeded()
	AssertNode(script).IsEqualTo(c.parse("const x = Math.pow(3, 4);"))
}

// Java: NestedCompilerRunnerTest.transpiles_exponentialOperatorToES2016_inTranspileOnlyMode()
func (c *NestedCompilerRunnerTest) TestTranspiles_exponentialOperatorToES2016_inTranspileOnlyMode() {
	c.options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.options.SetPreventLibraryInjection(true)
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)

	// Use the exponential operator as a very trivial transpilation test case.
	script := c.createScript("const x = 3 ** 4;")

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_ONLY, NewDefaultPassConfig(c.options)).
		AddScript(script, "test").
		Compile()

	c.verifyCompilationSucceeded()
	AssertNode(script).IsEqualTo(c.parse("const x = Math.pow(3, 4);"))
}

// Java: NestedCompilerRunnerTest.doesNotMangleOriginalExternProperties()
func (c *NestedCompilerRunnerTest) TestDoesNotMangleOriginalExternProperties() {
	c.options.SetPropertyRenaming(PropertyRenamingPolicy_ALL_UNQUOTED)
	c.original.SetExternProperties(jutil.NewLinkedHashSet("external"))
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)

	script := c.createScript("alert({external: 1, internal: 2});")

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, NewDefaultPassConfig(c.options)).
		AddScript(script, "test").
		Compile()

	c.verifyCompilationSucceeded()
	AssertNode(script).IsEqualTo(c.parse("alert({external: 1, a: 2})"))
}

// nestedCompilerRunnerTest_alwaysThrowsPassConfig is the anonymous
// PassConfig.PassConfigDelegate of stackTrace_specifiesInShadowAstCompilation.
type nestedCompilerRunnerTest_alwaysThrowsPassConfig struct {
	PassConfig_PassConfigDelegateBase
	options *CompilerOptions
}

// Java: NestedCompilerRunnerTest.stackTrace_specifiesInShadowAstCompilation() (anonymous PassConfigDelegate).getOptimizations()
func (p *nestedCompilerRunnerTest_alwaysThrowsPassConfig) GetOptimizations() *PassListBuilder {
	passes := NewPassListBuilder(p.options)
	passes.MaybeAdd(
		PassFactory_Builder_().
			SetName("throw").
			SetInternalFactory(
				func(compiler AbstractCompiler) CompilerPass {
					panic(jutil.NewAssertionError("foobarbaz"))
				}).
			Build())
	return passes
}

// Java: NestedCompilerRunnerTest.stackTrace_specifiesInShadowAstCompilation()
func (c *NestedCompilerRunnerTest) TestStackTrace_specifiesInShadowAstCompilation() {
	alwaysThrows := &nestedCompilerRunnerTest_alwaysThrowsPassConfig{options: c.options}
	alwaysThrows.self = alwaysThrows
	alwaysThrows.InitPassConfig_PassConfigDelegateBase(NewDefaultPassConfig(c.options))
	var alwaysThrowsPassConfig PassConfig = alwaysThrows
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)
	compilerRunner :=
		NestedCompilerRunner_Create(
			c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, alwaysThrowsPassConfig).
			AddScript(c.createScript(""), "test")

	ex := AssertThrows[jutil.IRuntimeException](compilerRunner.Compile)

	AssertThat(ex.GetMessage()).
		IsEqualTo("Exception during compilation: <shadow AST compilation>")
	AssertThat(ex.GetCause().GetMessage()).IsEqualTo("foobarbaz")
}

// TODO: lharker - add a test for debug logging.

// Java: NestedCompilerRunnerTest.doesNotRunTypechecking()
func (c *NestedCompilerRunnerTest) TestDoesNotRunTypechecking() {
	c.options.SetCheckTypes(true)
	c.original.Init([]*SourceFile{}, []*SourceFile{}, c.options)

	script := c.createScript("/** @type {string} */ const x = 0;")

	NestedCompilerRunner_Create(
		c.original, c.options, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE, NewDefaultPassConfig(c.options)).
		AddScript(script, "test").
		Compile()

	c.verifyCompilationSucceeded()
	AssertNode(script).IsEqualTo(c.parse("const x = 0;"))
}

// Java: NestedCompilerRunnerTest.createScript(String)
func (c *NestedCompilerRunnerTest) createScript(code string) rhino.Node {
	file := SourceFile_FromCodeWithKind("test.js", code, rhino.StaticSourceFile_SourceKind_STRONG)

	parsed := c.original.ParseTestCode(code)
	script := rhino.IR_Script()
	script.SetStaticSourceFile(file)
	script.SetIsInClosureUnawareSubtree(true)
	features := NodeUtil_GetFeatureSetOfScript(parsed)
	if features != nil {
		script.PutProp(rhino.Node_FEATURE_SET, features)
	} else {
		script.PutProp(rhino.Node_FEATURE_SET, parser.FeatureSet_BARE_MINIMUM)
	}

	script.AddChildrenToFront(parsed.RemoveChildren())
	return script
}

// Java: NestedCompilerRunnerTest.parse(String)
func (c *NestedCompilerRunnerTest) parse(src string) rhino.Node {
	return c.original.ParseTestCode(src)
}

// Java: NestedCompilerRunnerTest.verifyCompilationSucceeded()
func (c *NestedCompilerRunnerTest) verifyCompilationSucceeded() {
	AssertThat(c.original.GetErrors()).IsEmpty()
	AssertThat(c.original.GetWarnings()).IsEmpty()
}
