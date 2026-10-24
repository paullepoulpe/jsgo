// Port of com/google/javascript/jscomp/InjectTranspilationRuntimeLibrariesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/InjectTranspilationRuntimeLibrariesTest.java)

/*
 * Copyright 2018 The Closure Compiler Authors.
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
)

type InjectTranspilationRuntimeLibrariesTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	compiler               *Compiler
	languageOut            CompilerOptions_LanguageMode
	instrumentAsyncContext bool
}

// Java: InjectTranspilationRuntimeLibrariesTest()
func NewInjectTranspilationRuntimeLibrariesTest() *InjectTranspilationRuntimeLibrariesTest {
	return &InjectTranspilationRuntimeLibrariesTest{}
}

func TestInjectTranspilationRuntimeLibrariesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInjectTranspilationRuntimeLibrariesTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *InjectTranspilationRuntimeLibrariesTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method setup().
func (c *InjectTranspilationRuntimeLibrariesTest) JUnitBefore() { c.Setup() }

// JUnitAfter: the class has no @After method.
func (c *InjectTranspilationRuntimeLibrariesTest) JUnitAfter() {}

// Java: InjectTranspilationRuntimeLibrariesTest.setup()
func (c *InjectTranspilationRuntimeLibrariesTest) Setup() {
	c.compiler = NewCompiler()
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT5
	c.instrumentAsyncContext = false
}

// ParseAndRunInjectionPass parses the given code and runs the
// InjectTranspilationRuntimeLibraries pass over the resulting AST.
//
// Returns the set of paths to all the injected libraries.
//
// Java: InjectTranspilationRuntimeLibrariesTest.parseAndRunInjectionPass(String)
func (c *InjectTranspilationRuntimeLibrariesTest) ParseAndRunInjectionPass(js string) *jutil.LinkedHashSet[string] {
	options := NewCompilerOptions()
	options.SetLanguageOut(c.languageOut)
	options.SetInstrumentAsyncContext(c.instrumentAsyncContext)
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)

	c.compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", "")},
		[]*SourceFile{SourceFile_FromCode("testcode", js)},
		options)
	c.compiler.Parse()
	NewInjectTranspilationRuntimeLibraries(c.compiler).
		Process(c.compiler.GetExternsRoot(), c.compiler.GetJsRoot())

	return jutil.NewLinkedHashSet(c.compiler.GetRuntimeJsLibManager().GetInjectedLibraries()...)
}

// Java: InjectTranspilationRuntimeLibrariesTest.testEmptyInjected()
func (c *InjectTranspilationRuntimeLibrariesTest) TestEmptyInjected() {
	injected := c.ParseAndRunInjectionPass("")

	AssertThat(injected).IsEmpty()
}

// Java: InjectTranspilationRuntimeLibrariesTest.testMakeIteratorAndObjectAssignInjectedForObjectPatternRest()
func (c *InjectTranspilationRuntimeLibrariesTest) TestMakeIteratorAndObjectAssignInjectedForObjectPatternRest() {
	injected := c.ParseAndRunInjectionPass("const {a, ...rest} = something();")
	AssertThat(injected).ContainsExactly("es6/util/makeiterator", "es6/object/assign")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testObjectAssignInjectedForObjectSpread()
func (c *InjectTranspilationRuntimeLibrariesTest) TestObjectAssignInjectedForObjectSpread() {
	injected := c.ParseAndRunInjectionPass("const obj = {a, ...rest};")
	AssertThat(injected).ContainsExactly("es6/object/assign")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testForOf_injectsMakeIterator()
func (c *InjectTranspilationRuntimeLibrariesTest) TestForOf_injectsMakeIterator() {
	injected := c.ParseAndRunInjectionPass("for (x of []) {}")

	AssertThat(injected).ContainsExactly("es6/util/makeiterator", "es6/util/iteratorclose")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testArrayPattern_injectsMakeIterator()
func (c *InjectTranspilationRuntimeLibrariesTest) TestArrayPattern_injectsMakeIterator() {
	injected := c.ParseAndRunInjectionPass("var [a] = [];")

	AssertThat(injected).ContainsExactly("es6/util/makeiterator")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testObjectPattern_injectsNothing()
func (c *InjectTranspilationRuntimeLibrariesTest) TestObjectPattern_injectsNothing() {
	injected := c.ParseAndRunInjectionPass("var {a} = {};")

	AssertThat(injected).IsEmpty()
}

// Java: InjectTranspilationRuntimeLibrariesTest.testArrayPatternRest_injectsArrayFromIterator()
func (c *InjectTranspilationRuntimeLibrariesTest) TestArrayPatternRest_injectsArrayFromIterator() {
	injected := c.ParseAndRunInjectionPass("var [...a] = [];")

	AssertThat(injected).ContainsExactly("es6/util/makeiterator", "es6/util/arrayfromiterator")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testInjectsExecuteAsyncFunctionSupport_es5Out_includesGeneratorEngine()
func (c *InjectTranspilationRuntimeLibrariesTest) TestInjectsExecuteAsyncFunctionSupport_es5Out_includesGeneratorEngine() {
	injected := c.ParseAndRunInjectionPass("async function foo() {}")

	AssertThat(injected).ContainsExactly("es6/execute_async_generator", "es6/generator_engine")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testInjectsExecuteAsyncFunctionSupport_es6Out_doesNotIncludeGeneratorEngine()
func (c *InjectTranspilationRuntimeLibrariesTest) TestInjectsExecuteAsyncFunctionSupport_es6Out_doesNotIncludeGeneratorEngine() {
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT_2015
	injected := c.ParseAndRunInjectionPass("async function foo() {}")

	AssertThat(injected).ContainsExactly("es6/execute_async_generator")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testTaggedTemplateFirstArgCreaterInjected()
func (c *InjectTranspilationRuntimeLibrariesTest) TestTaggedTemplateFirstArgCreaterInjected() {
	injected := c.ParseAndRunInjectionPass("function tag(...a) {}; tag`hello`;")
	AssertThat(injected).ContainsExactly("es6/util/createtemplatetagfirstarg", "es6/util/restarguments")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testClassInheritance_injectsInheritsAndConstruct()
func (c *InjectTranspilationRuntimeLibrariesTest) TestClassInheritance_injectsInheritsAndConstruct() {
	injected := c.ParseAndRunInjectionPass("class A {} class B extends A {}")
	AssertThat(injected).ContainsExactly("es6/util/inherits", "es6/util/construct", "es6/util/arrayfromiterable")
}

// Java: InjectTranspilationRuntimeLibrariesTest.testClass_noInheritances_doesNotInject()
func (c *InjectTranspilationRuntimeLibrariesTest) TestClass_noInheritances_doesNotInject() {
	injected := c.ParseAndRunInjectionPass("class A {}")
	AssertThat(injected).IsEmpty()
}

// Java: InjectTranspilationRuntimeLibrariesTest.testInjectRuntimeLibraries_forceInject()
func (c *InjectTranspilationRuntimeLibrariesTest) TestInjectRuntimeLibraries_forceInject() {
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)

	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", "")},
		[]*SourceFile{SourceFile_FromCode("testcode", "")},
		options)
	compiler.Parse()

	pass :=
		NewInjectRuntimeLibraries(compiler, jutil.NewLinkedHashSet("es6/util/makeiterator"))
	pass.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())

	AssertThat(jutil.NewLinkedHashSet(compiler.GetRuntimeJsLibManager().GetInjectedLibraries()...)).
		Contains("es6/util/makeiterator")
}
