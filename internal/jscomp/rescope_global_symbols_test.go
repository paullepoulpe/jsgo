// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsTest.java @ closure-compiler v20261104
// (test/com/google/javascript/jscomp/RescopeGlobalSymbolsTest.java)

/*
 * Copyright 2011 The Closure Compiler Authors.
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

// RescopeGlobalSymbolsTest tests RescopeGlobalSymbols.
//
// Java: public final class RescopeGlobalSymbolsTest extends CompilerTestCase
type RescopeGlobalSymbolsTest struct {
	CompilerTestCaseBase

	assumeCrossChunkNames bool
	optimizeLocalAccess   CompilerOptions_OptimizeLocalAccess
}

// Java: RescopeGlobalSymbolsTest.NAMESPACE
const rescopeGlobalSymbolsTest_NAMESPACE = "_"

// Java: RescopeGlobalSymbolsTest()
func NewRescopeGlobalSymbolsTest() *RescopeGlobalSymbolsTest {
	c := &RescopeGlobalSymbolsTest{}
	c.InitCompilerTestCaseBase(c)
	c.assumeCrossChunkNames = true
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DISABLED
	return c
}

func TestRescopeGlobalSymbolsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRescopeGlobalSymbolsTest() })
}

// Java: RescopeGlobalSymbolsTest.getProcessor(Compiler)
func (c *RescopeGlobalSymbolsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRescopeGlobalSymbolsWithAddExtern(
		compiler, rescopeGlobalSymbolsTest_NAMESPACE, false, c.assumeCrossChunkNames, c.optimizeLocalAccess)
}

// Java: RescopeGlobalSymbolsTest.setUp()
func (c *RescopeGlobalSymbolsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.assumeCrossChunkNames = true
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DISABLED
}

// Java: RescopeGlobalSymbolsTest.testRescopeGlobalSymbolsRewriteCallback()
func (c *RescopeGlobalSymbolsTest) TestRescopeGlobalSymbolsRewriteCallback() {
	// RescopeGlobalSymbolsRewriteCallback is the callback used under the hood
	// by RescopeGlobalSymbols. We can test it implicitly by checking that
	// global symbols are rewritten to be properties of the namespace.
	c.Test("var a = 1; a + 1;", "_.a = 1; _.a + 1;")
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a; a = _.a = 1; a + 1;", "_.a;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_letNoInitializer()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_letNoInitializer() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("let a; a = 1; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("let a; a = _.a = 1; a + 1;", "_.a;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_loopInitializerMixed()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_loopInitializerMixed() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("for (var a, b; ; ) { b = 1; }").
				AddChunk("a = 2;").
				Build()),
		c.Expected("var b; for (_.a, b; ; ) { b = 1; }", "_.a = 2;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_writeInDifferentChunk()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_writeInDifferentChunk() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1; function getA() { return a; }").
				AddChunk("a = 2;").
				Build()),
		c.Expected("_.a = 1; var getA = function() { return _.a; };", "_.a = 2;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_writeInDifferentChunk_destructuring()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_writeInDifferentChunk_destructuring() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1; function getA() { return a; }").
				AddChunk("var {a} = {a: 2};").
				Build()),
		c.Expected("_.a = 1; var getA = function() { return _.a; };", "({a: _.a} = {a: 2});"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_destructuring()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_destructuring() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var {a} = {a: 1}; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a; ({a} = {a: 1}); _.a = a; a + 1;", "_.a;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_notUsedInDefiningChunk()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_notUsedInDefiningChunk() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("var a = 1;").AddChunk("a;").Build()),
		c.Expected("_.a = 1;", "_.a;"))
}

// Java: RescopeGlobalSymbolsTest.testImplicitGlobal()
func (c *RescopeGlobalSymbolsTest) TestImplicitGlobal() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = true
	c.TestSame("implicitGlobal = 1;")
}

// Java: RescopeGlobalSymbolsTest.testVarDeclarations()
func (c *RescopeGlobalSymbolsTest) TestVarDeclarations() {
	c.Test("var a = 1;", "_.a = 1;")
	c.Test("var a = 1, b = 2, c = 3;", "_.a = 1; _.b = 2; _.c = 3;")
	c.Test(
		"var a = 'str', b = 1, c = { foo: 'bar' }, d = function() {};",
		"_.a = 'str'; _.b = 1; _.c = { foo: 'bar' }; _.d = function() {};")
	c.Test("if(1){var x = 1;}", "if(1){_.x = 1;}")
	c.Test("var x;", "")
	c.Test("var a, b = 1;", "_.b = 1")
}

// Java: RescopeGlobalSymbolsTest.testVarDeclarations_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestVarDeclarations_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.TestSame("var a = 1;")
	c.TestSame("var a = 1, b = 2, c = 3;")
	c.TestSame("var a = 'str', b = 1, c = { foo: 'bar' }, d = function() {};")
	c.TestSame("if(1){var x = 1;}")
	c.TestSame("var x;")
	c.TestSame("var a, b = 1;")
}

// Java: RescopeGlobalSymbolsTest.testVarDeclarations_export()
func (c *RescopeGlobalSymbolsTest) TestVarDeclarations_export() {
	c.assumeCrossChunkNames = false
	c.Test("var _dumpException = 1;", "_._dumpException = 1")
}

// Java: RescopeGlobalSymbolsTest.testVarDeclarations_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestVarDeclarations_acrossChunks() {
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("var a = 1;").AddChunk("a").Build()),
		c.Expected("_.a = 1", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1, b = 2, c = 3;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a = 1; b = 2; _.c = 3;", "_.a;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1, b = 2, c = 3;").
				AddChunk("b;c;").
				Build()),
		c.Expected("var a;a = 1; _.b = 2; _.c = 3;", "_.b;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1, b = 2, c = 3;b;c;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a = 1; b = 2; _.c = 3;b;_.c", "_.a;_.c"))
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("var a, b = 1;").AddChunk("b").Build()),
		c.Expected("var a;_.b = 1;", "_.b"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a, b = 1, c = 2;").
				AddChunk("b").
				Build()),
		c.Expected("var a, c;_.b = 1;c = 2", "_.b"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a, b = 1, c = 2;").
				AddChunk("a").
				Build()),
		c.Expected("var b, c;b = 1;c = 2", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a=1; var b=2,c=3;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a=1;b=2;_.c=3", "_.a;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("1;var a, b = 1, c = 2;").
				AddChunk("b").
				Build()),
		c.Expected("var a, c;1;_.b = 1;c = 2", "_.b"))
}

// Java: RescopeGlobalSymbolsTest.testLetDeclarations()
func (c *RescopeGlobalSymbolsTest) TestLetDeclarations() {
	c.Test("let a = 1;", "_.a = 1;")
	c.Test("let a = 1, b = 2, c = 3;", "_.a = 1; _.b = 2; _.c = 3;")
	c.Test(
		"let a = 'str', b = 1, c = { foo: 'bar' }, d = function() {};",
		"_.a = 'str'; _.b = 1; _.c = { foo: 'bar' }; _.d = function() {};")
	c.TestSame("if(1){let x = 1;}")
	c.Test("let x;", "")
	c.Test("let a, b = 1;", "_.b = 1")
}

// Java: RescopeGlobalSymbolsTest.testLetDeclarations_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestLetDeclarations_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.TestSame("let a = 1;")
	c.TestSame("let a = 1, b = 2, c = 3;")
	c.TestSame("let a = 'str', b = 1, c = { foo: 'bar' }, d = function() {};")
	c.TestSame("if(1){let x = 1;}")
	c.TestSame("let x;")
	c.TestSame("let a, b = 1;")
}

// Java: RescopeGlobalSymbolsTest.testLetDeclarations_export()
func (c *RescopeGlobalSymbolsTest) TestLetDeclarations_export() {
	c.assumeCrossChunkNames = false
	c.Test("let _dumpException = 1;", "_._dumpException = 1")
}

// Java: RescopeGlobalSymbolsTest.testLetDeclarations_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestLetDeclarations_acrossChunks() {
	c.assumeCrossChunkNames = false
	// test references across chunks.
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("let a = 1;").AddChunk("a").Build()),
		c.Expected("_.a = 1", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("let a = 1, b = 2, c = 3;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a = 1; b = 2; _.c = 3;", "_.a;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("let a = 1, b = 2, c = 3;").
				AddChunk("b;c;").
				Build()),
		c.Expected("var a;a = 1; _.b = 2; _.c = 3;", "_.b;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("let a = 1, b = 2, c = 3;b;c;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a = 1; b = 2; _.c = 3;b;_.c", "_.a;_.c"))
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("let a, b = 1;").AddChunk("b").Build()),
		c.Expected("var a;_.b = 1;", "_.b"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("let a, b = 1, c = 2;").
				AddChunk("b").
				Build()),
		c.Expected("var a, c;_.b = 1;c = 2", "_.b"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("let a, b = 1, c = 2;").
				AddChunk("a").
				Build()),
		c.Expected("var b, c;b = 1;c = 2", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("let a=1; let b=2,c=3;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a=1;b=2;_.c=3", "_.a;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("1;let a, b = 1, c = 2;").
				AddChunk("b").
				Build()),
		c.Expected("var a, c;1;_.b = 1;c = 2", "_.b"))
	// test non-globals with same name as cross-chunk globals.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("1;let a, b = 1, c = 2;").
				AddChunk("if (true) { let b = 3; b; }").
				Build()))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("1;let a, b = 1, c = 2;").
				AddChunk("b; if (true) { let b = 3; b; }").
				Build()),
		c.Expected("var a, c; 1;_.b = 1;c = 2", "_.b; if (true) { let b = 3; b; }"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_let()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_let() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("let a = 1; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a; a = _.a = 1; a + 1;", "_.a;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_const()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_const() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("const a = 1; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a; a = _.a = 1; a + 1;", "_.a;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_updateAssignment()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_updateAssignment() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1; a++; a += 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a; a = _.a = 1; a++; _.a = a; a += 1; _.a = a;", "_.a;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_multipleWrites()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_multipleWrites() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1; var b = 1; a = b = 2; a + b;").
				AddChunk("a; b;").
				Build()),
		c.Expected("var a, b; a = _.a = 1; b = _.b = 1; a = _.a = b = _.b = 2; a + b;", "_.a; _.b;"))
}

// Java: RescopeGlobalSymbolsTest.testConstDeclarations()
func (c *RescopeGlobalSymbolsTest) TestConstDeclarations() {
	c.Test("const a = 1;", "_.a = 1;")
	c.Test("const a = 1, b = 2, c = 3;", "_.a = 1; _.b = 2; _.c = 3;")
	c.Test(
		"const a = 'str', b = 1, c = { foo: 'bar' }, d = function() {};",
		"_.a = 'str'; _.b = 1; _.c = { foo: 'bar' }; _.d = function() {};")
	c.TestSame("if(1){const x = 1;}")
}

// Java: RescopeGlobalSymbolsTest.testConstDeclarations_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestConstDeclarations_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.TestSame("const a = 1;")
	c.TestSame("const a = 1, b = 2, c = 3;")
	c.TestSame("const a = 'str', b = 1, c = { foo: 'bar' }, d = function() {};")
	c.TestSame("if(1){const x = 1;}")
}

// Java: RescopeGlobalSymbolsTest.testConstDeclarations_export()
func (c *RescopeGlobalSymbolsTest) TestConstDeclarations_export() {
	c.assumeCrossChunkNames = false
	c.Test("const _dumpException = 1;", "_._dumpException = 1")
}

// Java: RescopeGlobalSymbolsTest.testConstDeclarations_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestConstDeclarations_acrossChunks() {
	c.assumeCrossChunkNames = false
	// test references across chunks.
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("const a = 1;").AddChunk("a").Build()),
		c.Expected("_.a = 1", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("const a = 1, b = 2, c = 3;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a = 1; b = 2; _.c = 3;", "_.a;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("const a = 1, b = 2, c = 3;").
				AddChunk("b;c;").
				Build()),
		c.Expected("var a;a = 1; _.b = 2; _.c = 3;", "_.b;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("const a = 1, b = 2, c = 3;b;c;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a = 1; b = 2; _.c = 3;b;_.c", "_.a;_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("const a=1; const b=2,c=3;").
				AddChunk("a;c;").
				Build()),
		c.Expected("var b;_.a=1;b=2;_.c=3", "_.a;_.c"))
	// test non-globals with same name as cross-chunk globals.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("1;const a = 1, b = 1, c = 2;").
				AddChunk("if (true) { const b = 3; b; }").
				Build()))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("1;const a = 1, b = 1, c = 2;").
				AddChunk("b; if (true) { const b = 3; b; }").
				Build()),
		c.Expected("var a, c; 1;a = 1; _.b = 1;c = 2", "_.b; if (true) { const b = 3; b; }"))
}

// Java: RescopeGlobalSymbolsTest.testObjectDestructuringDeclarations()
func (c *RescopeGlobalSymbolsTest) TestObjectDestructuringDeclarations() {
	c.Test("var {a} = {}; a;", "({a: _.a} = {}); _.a;")
	c.Test("var {a: a} = {}; a;", "({a: _.a} = {}); _.a;")
	c.Test("var {key: a} = {}; a;", "({key: _.a} = {}); _.a;")
	c.Test("var {a: {b}} = {}; b;", "({a: {b: _.b}} = {}); _.b;")
	c.Test("var {a: {key: b}} = {}; b;", "({a: {key: _.b}} = {}); _.b;")
	c.Test("var {['computed']: a} = {}; a;", "({['computed']: _.a} = {}); _.a;")

	c.Test(
		"var {...a} = {a: 1, b: 2, c: 3}; a;", //
		"({..._.a} = {a: 1, b: 2, c: 3}); _.a;")
	c.Test(
		"var {a, ...b} = {x: 1, y: 2, z: 3}; a; b;",
		"({a: _.a, ..._.b} = {x: 1, y: 2, z: 3}); _.a; _.b;")

	c.Test("var {a = 3} = {}; a;", "({a: _.a = 3} = {}); _.a")
	c.Test("var {key: a = 3} = {}; a;", "({key: _.a = 3} = {}); _.a")
	c.Test("var {a} = {}, [b] = [], c = 3", "({a: _.a} = {}); [_.b] = []; _.c = 3;")
}

// Java: RescopeGlobalSymbolsTest.testObjectDestructuringDeclarations_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestObjectDestructuringDeclarations_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.TestSame("var {a} = {}; a;")
	c.TestSame("var {a: a} = {}; a;")
	c.TestSame("var {key: a} = {}; a;")
	c.TestSame("var {a: {b: b}} = {}; b;")
	c.TestSame("var {['computed']: a} = {}; a;")
	c.TestSame("var {a: a = 3} = {}; a;")
	c.TestSame("var {key: a = 3} = {}; a;")
	c.TestSame("var {['computed']: a = 3} = {}; a;")
	c.TestSame("var {a} = {}, b = 3;")
	c.TestSame("var [a] = [], b = 3;")
	c.TestSame("var a = 1, [b] = [], {c} = {};")
	c.TestSame("var {...a} = {x: 1, y: 2, z: 3}; a;")
	c.TestSame("var {a, ...b} = {x: 1, y: 2, z: 3}; a; b;")
}

// Java: RescopeGlobalSymbolsTest.testObjectDestructuringDeclarations_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestObjectDestructuringDeclarations_acrossChunks() {
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("var {a: a} = {};").AddChunk("a").Build()),
		c.Expected("({a: _.a} = {});", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var {a: a, b: b, c: c} = {};").
				AddChunk("a; c;").
				Build()),
		c.Expected("var b;({a: _.a, b: b, c: _.c} = {});", "_.a; _.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var {a: a, b: b, c: c} = {};").
				AddChunk("b; c;").
				Build()),
		c.Expected("var a;({a: a, b: _.b, c: _.c} = {});", "_.b; _.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var {a: a, b: b, c: c} = {}; b; c;").
				AddChunk("a; c;").
				Build()),
		c.Expected("var b;({a: _.a, b: b, c: _.c} = {});b;_.c;", "_.a; _.c"))

	// Test var declarations containing a mix of destructuring and regular names
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("var {a} = {}, b;").AddChunk("a").Build()),
		c.Expected("var b; ({a: _.a} = {});", "_.a"))

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var {a} = {}, b = 3;").
				AddChunk("b").
				Build()),
		c.Expected("var a; ({a} = {}); _.b = 3;", "_.b"))

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var {a: a, ...c} = {};").
				AddChunk("a;").
				Build()),
		c.Expected("var c; ({a: _.a, ...c} = {});", "_.a;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var {a: a, ...c} = {};").
				AddChunk("c;").
				Build()),
		c.Expected("var a; ({a: a, ..._.c} = {});", "_.c"))
}

// Java: RescopeGlobalSymbolsTest.testObjectDestructuringAssignments()
func (c *RescopeGlobalSymbolsTest) TestObjectDestructuringAssignments() {
	c.Test("var a, b; ({key1: a, key2: b} = {}); a; b;", "({key1: _.a, key2: _.b} = {}); _.a; _.b;")
	// Test a destructuring assignment with mixed global and local variables.
	c.Test(
		"var a; if (true) { let b; ({key1: a, key2: b} = {}); b; } a;",
		"if (true) { let b; ({key1: _.a, key2: b} = {}); b; } _.a;")
	c.Test("var obj = {}; ({a: obj.a} = {}); obj.a;", "_.obj = {}; ({a: _.obj.a} = {}); _.obj.a;")
	c.Test(
		"var obj = {}; ({a:   obj['foo bar']} = {});   obj['foo bar'];",
		"  _.obj = {}; ({a: _.obj['foo bar']} = {}); _.obj['foo bar'];")

	c.Test(
		"var a = {}; ({...a.b} = {a: 1, b: 2, c: 3}); a.b;", //
		"_.a = {}; ({..._.a.b} = {a: 1, b: 2, c: 3}); _.a.b;")
	c.Test(
		"var a = {}; ({c, ...a.b} = {a: 1, b: 2, c: 3}); a.b;", //
		"_.a = {}; ({c, ..._.a.b} = {a: 1, b: 2, c: 3}); _.a.b;")
}

// Java: RescopeGlobalSymbolsTest.testArrayDestructuringDeclarations()
func (c *RescopeGlobalSymbolsTest) TestArrayDestructuringDeclarations() {
	c.Test("var [a] = [1]; a", "[_.a] = [1]; _.a;")
	c.Test("var [a, b, c] = [1, 2, 3]; a; b; c;", "[_.a, _.b, _.c] = [1, 2, 3]; _.a; _.b; _.c")
	c.Test("var [[a, b], c] = []; a; b; c", "[[_.a, _.b], _.c] = []; _.a; _.b; _.c;")
	c.Test("var [a = 5] = [1]; a", "[_.a = 5] = [1]; _.a;")
	c.Test("var [a, b = 5] = [1]; a; b;", "[_.a, _.b = 5] = [1]; _.a; _.b;")
	c.Test("var [...a] = [1, 2, 3]; a;", "[..._.a] = [1, 2, 3]; _.a;")
	c.Test("var [a, ...b] = [1, 2, 3]; a; b;", "[_.a, ..._.b] = [1, 2, 3]; _.a; _.b;")
	c.Test("var [a] = 1, b = 2; a; b;", "[_.a] = 1; _.b = 2; _.a; _.b;")
}

// Java: RescopeGlobalSymbolsTest.testArrayDestructuringDeclarations_sameChunk()
func (c *RescopeGlobalSymbolsTest) TestArrayDestructuringDeclarations_sameChunk() {
	c.assumeCrossChunkNames = false
	c.TestSame("var [a] = [1]; a")
	c.TestSame("var [a, b] = [1, 2]; a; b;")
	c.TestSame("var [[a, b], c] = []; a; b; c")
	c.TestSame("var [a = 5] = [1]; a")
	c.TestSame("var [a, b = 5] = [1]; a; b;")
	c.TestSame("var [...a] = [1, 2, 3]; a;")
	c.TestSame("var [a, ...b] = [1, 2, 3]; a; b;")
}

// Java: RescopeGlobalSymbolsTest.testArrayDestructuringDeclarations_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestArrayDestructuringDeclarations_acrossChunks() {
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForUnordered().AddChunk("var [a] = [];").AddChunk("a").Build()),
		c.Expected("[_.a] = [];", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var [a, b, c] = [];").
				AddChunk("a; c;").
				Build()),
		c.Expected("var b; [_.a, b, _.c] = [];", "_.a; _.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var [a, b, c] = [];").
				AddChunk("b; c;").
				Build()),
		c.Expected("var a; [a, _.b, _.c] = [];", "_.b; _.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var [a, b, c] = []; b; c;").
				AddChunk("a; c;").
				Build()),
		c.Expected("var b; [_.a, b, _.c] = []; b; _.c;", "_.a; _.c"))

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var [a, ...c] = [];").
				AddChunk("a;").
				Build()),
		c.Expected("var c; ([_.a, ...c] = []);", "_.a;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var [a, ...c] = [];").
				AddChunk("c;").
				Build()),
		c.Expected("var a; ([a, ..._.c] = []);", "_.c"))
}

// Java: RescopeGlobalSymbolsTest.testArrayDestructuringAssignments()
func (c *RescopeGlobalSymbolsTest) TestArrayDestructuringAssignments() {
	c.Test("var a, b; [a, b] = []; a; b;", "[_.a, _.b] = []; _.a; _.b;")
	// Test a destructuring assignment with mixed global and local variables.
	c.Test(
		"var a; if (true) { let b; [a, b] = []; b; } a;",
		"if (true) { let b; [_.a, b] = []; b; } _.a;")
	// Test assignments to qualified names and quoted properties.
	c.Test("var obj = {}; [obj.a] = []; obj.a;", "_.obj = {}; [_.obj.a] = []; _.obj.a;")
	c.Test(
		"var obj = {}; [  obj['foo bar']] = [];   obj['foo bar'];",
		"  _.obj = {}; [_.obj['foo bar']] = []; _.obj['foo bar'];")

	c.Test(
		"var a = {}; ([...a.b] = [1, 2, 3]); a.b;", //
		"_.a = {}; ([..._.a.b] = [1, 2, 3]); _.a.b;")
	c.Test(
		"var a = {}; ([c, ...a.b] = [1, 2, 3]); a.b;", //
		"_.a = {}; ([c, ..._.a.b] = [1, 2, 3]); _.a.b;")
}

// Java: RescopeGlobalSymbolsTest.testClasses()
func (c *RescopeGlobalSymbolsTest) TestClasses() {
	c.Test("class A {}", "_.A = class {};")
	c.Test("class A {} class B extends A {}", "_.A = class {}; _.B = class extends _.A {}")
	c.Test("class A {} let a = new A;", "_.A = class {}; _.a = new _.A;")
	c.Test(
		`const PI = 3.14;
class A {
  static printPi() {
    console.log(PI);
  }
}
A.printPi();
`,
		`_.PI = 3.14;
_.A = class {
  static printPi() {
    console.log(_.PI);
  }
}
_.A.printPi();
`)

	// Test that class expression names are not rewritten.
	c.Test("var A = class Name {};", "_.A = class Name {};")
	c.Test("var A = class A {};", "_.A = class A {};")
}

// Java: RescopeGlobalSymbolsTest.testClasses_nonGlobal()
func (c *RescopeGlobalSymbolsTest) TestClasses_nonGlobal() {
	c.TestSame("if (true) { class A {} }")
	c.Test("function foo() { class A {} }", "_.foo = function() { class A {} };")
	c.Test("const A = 5; { class A {} }", "_.A = 5; { class A {} }")
}

// Java: RescopeGlobalSymbolsTest.testClasses_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestClasses_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.Test("class A {}", "var A = class {};")
	c.Test("class A {} class B extends A {}", "var A = class {}; var B = class extends A {}")
	c.TestSame("if (true) { class A {} }")
}

// Java: RescopeGlobalSymbolsTest.testForLoops()
func (c *RescopeGlobalSymbolsTest) TestForLoops() {
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("for (var i = 0, c = 2; i < 1000; i++);").
				AddChunk("c").
				Build()),
		c.Expected("var i;for (i = 0, _.c = 2; i < 1000; i++);", "_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("      for (var i = 0, c = 2;   i < 1000;   i++);").
				AddChunk("i").
				Build()),
		c.Expected("var c;for (  _.i = 0, c = 2; _.i < 1000; _.i++);", "_.i"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("       for (var {i: i, c:   c} = {};  i < 1000; i++);").
				AddChunk("c").
				Build()),
		c.Expected("var i; for (   ({i: i, c: _.c} = {}); i < 1000; i++);", "_.c"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("       for (var [i,   c] = [0, 2]; i < 1000; i++);").
				AddChunk("c;").
				Build()),
		c.Expected("var i; for (    [i, _.c] = [0, 2]; i < 1000; i++);", "_.c;"))
}

// Java: RescopeGlobalSymbolsTest.testForLoops_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestForLoops_acrossChunks() {
	c.Test("for (var i = 0; i < 1000; i++);", "for (_.i = 0; _.i < 1000; _.i++);")
	c.Test("for (var i = 0, c = 2; i < 1000; i++);", "for (_.i = 0, _.c = 2; _.i < 1000; _.i++);")
	c.Test(
		"for (var i = 0, c = 2, d = 3; i < 1000; i++);",
		"for (_.i = 0, _.c = 2, _.d = 3; _.i < 1000; _.i++);")
	c.Test(
		"for (var i = 0, c = 2, d = 3, e = 4; i < 1000; i++);",
		"for (_.i = 0, _.c = 2, _.d = 3, _.e = 4; _.i < 1000; _.i++);")
	c.Test("for (var i = 0; i < 1000;)i++;", "for (_.i = 0; _.i < 1000;)_.i++;")
	c.Test("for (var i = 0,b; i < 1000;)i++;b++", "for (_.i = 0,_.b; _.i < 1000;)_.i++;_.b++")
	c.Test("var o={};for (var i in o)i++;", "_.o={};for (_.i in _.o)_.i++;")

	// Test destructuring.
	c.Test("for (var [i] = [0]; i < 1000; i++);", "for ([_.i] = [0]; _.i < 1000; _.i++);")
	c.Test("for (var {i: i} = {}; i < 1000; i++);", " for (({i: _.i} = {}); _.i < 1000; _.i++);")
	c.TestSame("for (let [i] = [0]; i  < 1000; i++);")
}

// Java: RescopeGlobalSymbolsTest.testForInLoops_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestForInLoops_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.TestSame("for (var i in {});")
	c.TestSame("for (var [a] in {});")
	c.TestSame("for (var {a: a} in {});")
}

// Java: RescopeGlobalSymbolsTest.testForInLoops_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestForInLoops_acrossChunks() {
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("for (var i in {});").
				AddChunk("i").
				Build()),
		c.Expected("for (_.i in {});", "_.i"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("       for (var [  a, b] in {});").
				AddChunk("a;").
				Build()),
		c.Expected("var b; for (    [_.a, b] in {});", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("       for (var {i: i, c:   c} in {});").
				AddChunk("c").
				Build()),
		c.Expected("var i; for (    {i: i, c: _.c} in {});", "_.c"))
}

// Java: RescopeGlobalSymbolsTest.testForOfLoops_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestForOfLoops_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.TestSame("for (var i of [1, 2, 3]);")
	c.TestSame("for (var [a] of []);")
	c.TestSame("for (var {a: a} of {});")
}

// Java: RescopeGlobalSymbolsTest.testForOfLoops_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestForOfLoops_acrossChunks() {
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("for (var i of []);").
				AddChunk("i").
				Build()),
		c.Expected("for (_.i of []);", "_.i"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("       for (var [  a, b] of []);").
				AddChunk("a;").
				Build()),
		c.Expected("var b; for (    [_.a, b] of []);", "_.a"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("       for (var {i: i, c:   c} of []);").
				AddChunk("c").
				Build()),
		c.Expected("var i; for (    {i: i, c: _.c} of []);", "_.c"))
}

// Java: RescopeGlobalSymbolsTest.testForAwaitOfLoops_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestForAwaitOfLoops_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.TestSame("async () => { for await (var i of [1, 2, 3]);}")
	c.TestSame("async () => { for await (var [a] of []);}")
	c.TestSame("async () => { for await (var {a: a} of {});}")
}

// Java: RescopeGlobalSymbolsTest.testForAwaitOfLoops_acrossChunks()
func (c *RescopeGlobalSymbolsTest) TestForAwaitOfLoops_acrossChunks() {
	// TODO(b/128938049): re-enable it once we support top-level await.
	c.assumeCrossChunkNames = false
	c.TestError(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("for await (var i of []);").
				AddChunk("i").
				Build()),
		RhinoErrorReporter_PARSE_ERROR)
	c.TestError(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("       for await (var [  a, b] of []);").
				AddChunk("a;").
				Build()),
		RhinoErrorReporter_PARSE_ERROR)
	c.TestError(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("       for await (var {i: i, c:   c} of []);").
				AddChunk("c").
				Build()),
		RhinoErrorReporter_PARSE_ERROR)
}

// Java: RescopeGlobalSymbolsTest.testFunctionStatements()
func (c *RescopeGlobalSymbolsTest) TestFunctionStatements() {
	c.Test("function test(){}", "_.test=function (){}")
	// Ignore block-scoped function declarations.
	c.TestSame("if(1) function test(){}")
	c.Test("async function test() {}", "_.test = async function() {}")
	c.Test("function *test() {}", "_.test = function *() {}")
}

// Java: RescopeGlobalSymbolsTest.testFunctionStatements_allSameChunk()
func (c *RescopeGlobalSymbolsTest) TestFunctionStatements_allSameChunk() {
	c.assumeCrossChunkNames = false
	c.Test("function f() {}", "var f = function() {}")
	c.TestSame("if (true) { function f() {} }")
}

// Java: RescopeGlobalSymbolsTest.testFunctionStatements_freeCallSemantics1()
func (c *RescopeGlobalSymbolsTest) TestFunctionStatements_freeCallSemantics1() {
	c.DisableCompareAsTree()

	// This triggers free call.
	c.Test(
		"function x(){this};var y=function(){var val=x()||{}}",
		"_.x=function(){this};_.y=function(){var val=(0,_.x)()||{}}")
	c.Test("function x(){this;x()}", "_.x=function(){this;(0,_.x)()}")
	c.Test("var a=function(){this};a()", "_.a=function(){this};(0,_.a)()")
	// Always trigger free calls for variables assigned through destructuring.
	c.Test("var {a: a} = {a: function() {}}; a();", "({a:_.a}={a:function(){}});(0,_.a)()")

	c.Test("var ns = {}; ns.a = function() {}; ns.a()", "_.ns={};_.ns.a=function(){};_.ns.a()")
}

// Java: RescopeGlobalSymbolsTest.testFunctionStatements_freeCallSemantics2()
func (c *RescopeGlobalSymbolsTest) TestFunctionStatements_freeCallSemantics2() {
	// Cases where free call forcing through (0, foo)() is not necessary.
	c.Test("var a=function(){};a()", "_.a=function(){};_.a()")
	c.Test("function a(){};a()", "_.a=function(){};;_.a()")
	c.Test("var a;a=function(){};a()", "_.a=function(){};_.a()")

	// Test that calls to arrow functions are not forced to be free calls
	c.Test("var a = () => {}; a();", "_.a = () => {}; _.a();")
	c.Test("var a = () => this; a();", "_.a = () => this; _.a();")
}

// Java: RescopeGlobalSymbolsTest.testFunctionStatements_freeCallSemantics3()
func (c *RescopeGlobalSymbolsTest) TestFunctionStatements_freeCallSemantics3() {
	c.DisableCompareAsTree()

	// Ambiguous cases.
	c.Test("var a=1;a=function(){};a()", "_.a=1;_.a=function(){};(0,_.a)()")
	c.Test("var b;var a=b;a()", "_.a=_.b;(0,_.a)()")
}

// Java: RescopeGlobalSymbolsTest.testFunctionStatements_defaultParameters()
func (c *RescopeGlobalSymbolsTest) TestFunctionStatements_defaultParameters() {
	c.Test("var a = 1; function f(param = a) {}", "_.a = 1; _.f = function(param = _.a) {}")
}

// Java: RescopeGlobalSymbolsTest.testDeeperScopes()
func (c *RescopeGlobalSymbolsTest) TestDeeperScopes() {
	c.Test("var a = function(b){return b}", "_.a = function(b){return b}")
	c.Test("var a = function(b){var a; return a+b}", "_.a = function(b){var a; return a+b}")
	c.Test("var a = function(a,b){return a+b}", "_.a = function(a,b){return a+b}")
	c.Test(
		"var x=1,a = function(b){var a; return a+b+x}",
		"_.x=1;_.a = function(b){var a; return a+b+_.x}")
	c.Test(
		"var x=1,a = function(b){return function(){var a;return a+b+x}}",
		"_.x=1;_.a = function(b){return function(){var a; return a+b+_.x}}")
}

// Java: RescopeGlobalSymbolsTest.testTryCatch()
func (c *RescopeGlobalSymbolsTest) TestTryCatch() {
	c.Test("try{var a = 1}catch(e){throw e}", "try{_.a = 1}catch(e){throw e}")
}

// Java: RescopeGlobalSymbolsTest.testShadowInFunctionScope()
func (c *RescopeGlobalSymbolsTest) TestShadowInFunctionScope() {
	c.Test("var _ = 1; (function () { _ = 2 })()", "_._ = 1; (function () { _._ = 2 })()")
	c.Test(
		"function foo() { var _ = {}; _.foo = foo; _.bar = 1; }",
		"_.foo = function () { var _$ = {}; _$.foo = _.foo; _$.bar = 1}")
	c.Test(
		"function foo() { var {key: _} = {}; _.foo = foo; _.bar = 1; }",
		"_.foo = function () { var {key: _$} = {}; _$.foo = _.foo; _$.bar = 1}")
	c.Test(
		`function foo() { var _ = {}; _.foo = foo; _.bar = 1;
(function() { var _ = 0;})() }
`,
		`_.foo = function () { var _$ = {}; _$.foo = _.foo; _$.bar = 1;
(function() { var _$ = 0;})() }
`)
	c.Test(
		`function foo() { var _ = {}; _.foo = foo; _.bar = 1;
var _$ = 1; }
`,
		`_.foo = function () { var _$ = {}; _$.foo = _.foo; _$.bar = 1;
var _$$ = 1; }
`)
	c.Test(
		`function foo() { var _ = {}; _.foo = foo; _.bar = 1;
var _$ = 1; (function() { _ = _$ })() }
`,
		`_.foo = function () { var _$ = {}; _$.foo = _.foo; _$.bar = 1;
var _$$ = 1; (function() { _$ = _$$ })() }
`)
	c.Test(
		`function foo() { var _ = {}; _.foo = foo; _.bar = 1;
var _$ = 1, _$$ = 2 (function() { _ = _$ = _$$;
var _$, _$$$ })() }
`,
		`_.foo = function () { var _$ = {}; _$.foo = _.foo; _$.bar = 1;
var _$$ = 1, _$$$ = 2 (function() { _$ = _$$ = _$$$;
var _$$, _$$$$ })() }
`)
	c.Test(
		"var a = 5; function foo(_) { return _; }", "_.a = 5; _.foo = function(_$) { return _$; }")
	c.Test(
		"var a = 5; function foo({key: _}) { return _; }",
		"_.a = 5; _.foo = function({key: _$}) { return _$; }")
	c.Test(
		"var a = 5; function foo([_]) { return _; }",
		"_.a = 5; _.foo = function([_$]) { return _$; }")
	// We accept this unnecessary renaming as acceptable to simplify pattern
	// matching in the traversal.
	c.Test("function foo() { var _$a = 1;}", "_.foo = function () { var _$a$ = 1;}")
}

// Java: RescopeGlobalSymbolsTest.testShadowInBlockScope()
func (c *RescopeGlobalSymbolsTest) TestShadowInBlockScope() {
	c.Test(
		"var foo = 1; if (true) { const _ = {}; _.foo = foo; _.bar = 1; }",
		"_.foo = 1; if (true) { const _$ = {}; _$.foo = _.foo; _$.bar = 1}")
	c.Test(
		"var foo = 1; if (true) { const {key: _} = {}; _.foo = foo; _.bar = 1; }",
		"_.foo = 1; if (true) { const {key: _$} = {}; _$.foo = _.foo; _$.bar = 1}")
	c.Test(
		"var foo = 1; if (true) { const _ = {}; _.foo = foo; _.bar = 1; const _$ = 1; }",
		"_.foo = 1; if (true) { const _$ = {}; _$.foo = _.foo; _$.bar = 1;  const _$$ = 1; }")
	c.Test(
		"var foo = 1; if (true) { const [_] = [{}]; _.foo = foo; }",
		"  _.foo = 1; if (true) { const [_$] = [{}]; _$.foo = _.foo; }")
}

// Java: RescopeGlobalSymbolsTest.testExterns()
func (c *RescopeGlobalSymbolsTest) TestExterns() {
	c.Test(c.Externs("var document;"), c.Srcs("document"), c.Expected("document"))
	c.Test(
		c.Externs("var document;"),
		c.Srcs("document.getElementsByTagName('test')"),
		c.Expected("document.getElementsByTagName('test')"))
	c.Test(
		c.Externs("var document;"),
		c.Srcs("document.getElementsByTagName('test')"),
		c.Expected("document.getElementsByTagName('test')"))
	c.Test(
		c.Externs("var document;document.getElementsByTagName"),
		c.Srcs("document.getElementsByTagName('test')"),
		c.Expected("document.getElementsByTagName('test')"))
	c.Test(
		c.Externs("var document,navigator"),
		c.Srcs("document.navigator;navigator"),
		c.Expected("document.navigator;navigator"))
	c.Test(
		c.Externs("var iframes"),
		c.Srcs("function test() { iframes.resize(); }"),
		c.Expected("_.test = function() { iframes.resize(); }"))
	c.Test(c.Externs("var iframes"), c.Srcs("var foo = iframes;"), c.Expected("_.foo = iframes;"))
	c.Test(c.Externs("class A {}"), c.Srcs("A"), c.Expected("A"))
	// Special names.
	c.Test(
		c.Externs("var arguments, window, eval;"),
		c.Srcs("arguments;window;eval;"),
		c.Expected("arguments;window;eval;"))
	// Actually not an extern.
	c.Test(c.Externs(""), c.Srcs("document"), c.Expected("document"))
	// Javascript builtin objects
	c.TestSame(
		`Object;Function;Array;String;Boolean;Number;Math;
Date;RegExp;JSON;Error;EvalError;ReferenceError;
SyntaxError;TypeError;URIError;
`)
}

// Java: RescopeGlobalSymbolsTest.testSameVarDeclaredInExternsAndSource()
func (c *RescopeGlobalSymbolsTest) TestSameVarDeclaredInExternsAndSource() {
	c.Test(
		c.Externs("/** @const */ var ns = {}; function f() {}"),
		c.Srcs("/** @const */ var ns = ns || {};"),
		c.Expected("/** @const */ ns = ns || {};"))

	c.Test(c.Externs("var x;"), c.Srcs("var x = 1; x = 2;"), c.Expected("x = 1; x = 2;"))

	c.Test(
		c.Externs("var x;"),
		c.Srcs("function f() { var x; x = 1; }"),
		c.Expected("_.f = function() { var x; x = 1; }"))

	c.Test(
		c.Externs("var x;"),
		c.Srcs("function f() { x = 1; }"),
		c.Expected("_.f = function() { x = 1; };"))

	c.Test(c.Externs("var x, y;"), c.Srcs("var x = 1, y = 2;"), c.Expected("x = 1; y = 2;"))

	c.Test(c.Externs("var x, y;"), c.Srcs("var x, y = 2;"), c.Expected("y = 2;"))

	c.Test(c.Externs("var x;"), c.Srcs("var x = 1, y = 2;"), c.Expected("x = 1; _.y = 2;"))

	c.Test(c.Externs("var x;"), c.Srcs("var y = 2, x = 1;"), c.Expected("_.y = 2; x = 1;"))

	c.Test(c.Externs("var x;"), c.Srcs("var y, x = 1;"), c.Expected("x = 1;"))

	c.Test(
		c.Externs("var foo;"),
		c.Srcs("var foo = function(x) { if (x > 0) { var y = foo; } };"),
		c.Expected("foo = function(x) { if (x > 0) { var y = foo; } };"))

	// The parameter x doesn't conflict with the x in the source
	c.Test(c.Externs("function f(x) {}"), c.Srcs("var f = 1; var x = 2;"), c.Expected("f = 1; _.x = 2;"))
}

// Java: RescopeGlobalSymbolsTest.testSameVarDeclaredInExternsAndSource2()
func (c *RescopeGlobalSymbolsTest) TestSameVarDeclaredInExternsAndSource2() {
	c.assumeCrossChunkNames = false

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk(
					`Foo = function() { this.b = ns; };
var f = function(a) {
  if (a instanceof Foo && a.b === ns) {}
},
ns = {},
g = function(a) { var b = new Foo; };
`).
				AddChunk("f; g;").
				Build()),
		c.Expected(
			`var ns;
Foo = function() { this.b = ns; };
_.f = function(a) {
  if (a instanceof Foo && a.b === ns) {}
};
ns = {};
_.g = function(a) { var b = new Foo; };
`,
			"_.f; _.g;"))

	c.Test(
		c.Externs("var y;"),
		c.Srcs("var x = 1, y = 2; function f() { return x + y; }"),
		c.Expected("var x; x = 1; y = 2; var f = function() { return x + y; }"))
}

// Java: RescopeGlobalSymbolsTest.testArrowFunctions()
func (c *RescopeGlobalSymbolsTest) TestArrowFunctions() {
	c.Test("const fn = () => 3;", "_.fn = () => 3;")
	c.Test("const PI = 3.14; const fn = () => PI;", "_.PI = 3.14; _.fn = () => _.PI")
	c.Test("let a = 3; const fn = () => a = 4;", "_.a = 3; _.fn = () => _.a = 4;")
	c.Test("const PI = 3.14; (() => PI)()", "_.PI = 3.14; (() => _.PI)();")
	c.Test("const PI = 3.14; (() => { return PI; })()", "_.PI = 3.14; (() => { return _.PI; })()")
}

// Java: RescopeGlobalSymbolsTest.testEnhancedObjectLiterals()
func (c *RescopeGlobalSymbolsTest) TestEnhancedObjectLiterals() {
	c.Test("var a = 3; var obj = {[a]: a};", "_.a = 3; _.obj = {[_.a]: _.a};")
	c.Test(
		"var g = 3; var obj = {a() {}, b() { return g; }};",
		"_.g = 3; _.obj = {a() {}, b() { return _.g; }};")
	c.Test("var a = 1; var obj = {a}", "_.a = 1; _.obj = {a: _.a};")
}

// Java: RescopeGlobalSymbolsTest.testEs6Modules()
func (c *RescopeGlobalSymbolsTest) TestEs6Modules() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	// Test that this pass does nothing to ES6 modules.
	c.TestSame("var a = 3; a; export default a;")

	c.assumeCrossChunkNames = false
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 3; export {a};").
				AddChunk("import {a} from './input0';").
				Build()))
}

// Java: RescopeGlobalSymbolsTest.testEmptyDestructuring()
func (c *RescopeGlobalSymbolsTest) TestEmptyDestructuring() {
	c.TestSame("var {} = {};")
}

// Java: RescopeGlobalSymbolsTest.testAllChunksLocalAccessOptimization()
func (c *RescopeGlobalSymbolsTest) TestAllChunksLocalAccessOptimization() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("var a = 1; a + 1;").AddChunk("a;").Build()),
		c.Expected("var a; a = _.a = 1; a + 1;", "var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksLocalAccessOptimization_reassignedInDefiningChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksLocalAccessOptimization_reassignedInDefiningChunk() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a = 2; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a; a = _.a = 1; a = _.a = 2; a + 1;", "var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksLocalAccessOptimization_reassignedOutsideStaticExecution()
func (c *RescopeGlobalSymbolsTest) TestAllChunksLocalAccessOptimization_reassignedOutsideStaticExecution() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; function f() { a = 2; } a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a; a = _.a = 1; var f = function() { a = _.a = 2; }; a + 1;", "_.a;"))

	// Writes inside an `if` block occur in a block scope (!t.inGlobalScope()), so other chunks do
	// not alias `a` (`var {a} = _;`) and instead read `_.a` directly when `readA()` is called.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function readA() { return a; }").
				AddChunk("var a = 1; if (a + 3) { a++; readA(); }").
				Build()),
		c.Expected(
			"_.readA = function() { return _.a; };",
			"var {readA} = _; var a; a = _.a = 1; if (a + 3) { a++; _.a = a; readA(); }"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksLocalAccessOptimization_reassignedInOtherChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksLocalAccessOptimization_reassignedInOtherChunk() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1;").
				AddChunk("a = 2; a;").
				Build()),
		c.Expected("_.a = 1; _.a + 1;", "_.a = 2; _.a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksLocalAccessOptimization_multipleChunks()
func (c *RescopeGlobalSymbolsTest) TestAllChunksLocalAccessOptimization_multipleChunks() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1;").
				AddChunk("a; a + 2;").
				AddChunk("a; a + 3;").
				Build()),
		c.Expected("var a; a = _.a = 1; a + 1;", "var {a} = _; a; a + 2;", "var {a} = _; a; a + 3;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("var a = 1; a + 1;").AddChunk("a;").Build()),
		c.Expected("var a = _.a = 1; a + 1;", "var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_let()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_let() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("let a = 1; a + 1;").AddChunk("a;").Build()),
		c.Expected("var a = _.a = 1; a + 1;", "var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_const()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_const() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("const a = 1; a + 1;").AddChunk("a;").Build()),
		c.Expected("var a = _.a = 1; a + 1;", "var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_mixedVarDeclaration()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_mixedVarDeclaration() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1, b = 2, c = 3, d = 4; a + c + d;").
				AddChunk("a; b; d;").
				Build()),
		c.Expected(
			"var a = _.a = 1; _.b = 2; var c = 3, d = _.d = 4; a + c + d;",
			"var {a, b, d} = _; a; b; d;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_reassignedInDefiningChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_reassignedInDefiningChunk() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a = 2; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a = _.a = 1; a = _.a = 2; a + 1;", "var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_multipleWrites()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_multipleWrites() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; var b = 1; a = b = 2; a + b;").
				AddChunk("a; b;").
				Build()),
		c.Expected(
			"var a = _.a = 1; var b = _.b = 1; a = _.a = b = _.b = 2; a + b;",
			"var {a, b} = _; a; b;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_es5()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_es5() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("var a = 1; a + 1;").AddChunk("a;").Build()),
		c.Expected("var a = _.a = 1; a + 1;", "var a = _.a; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_destructuring()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_destructuring() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var [a, b, c] = [1, 2, 3]; a + b;").
				AddChunk("a; c;").
				Build()),
		c.Expected("var a, b; [a, b, _.c] = [1, 2, 3]; _.a = a; a + b;", "var {a, c} = _; a; c;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1, [b, c] = [2, 3], d = 4; a + b + d;").
				AddChunk("a; c; d;").
				Build()),
		c.Expected(
			"var a = _.a = 1, b; [b, _.c] = [2, 3]; var d = _.d = 4; a + b + d;",
			"var {a, c, d} = _; a; c; d;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_forLoop()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_forLoop() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("for (var a = 0, b = 1, c = 2; a < 10; a++) { b++; }").
				AddChunk("a; c;").
				Build()),
		c.Expected(
			"_.a = {}; var c = _.c = {}; var _a, b; for (_a = _.a._ = 0, b = 1, c._ = 2; _a < 10;"+
				" ++_.a._, _a++) { b++; }",
			"var {a, c} = _; a._; c._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_wrappedReassignableComposition()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_wrappedReassignableComposition() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1, b = 2; function f() { b = 3; } a + b;").
				AddChunk("a; b;").
				Build()),
		c.Expected(
			"_.b = {}; var a = _.a = 1, _b = _.b._ = 2; var f = function() { _b = _.b._ = 3; }; a +"+
				" _b;",
			"var {a, b} = _; a; b._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_uninitializedExternAndNonLocal()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_uninitializedExternAndNonLocal() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Externs("var externVar;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var externVar, a = 1; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a = _.a = 1; a + 1;", "var {a} = _; a;"))
	c.Test(
		c.Externs("var externVar;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var externVar; var a = 1; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var a = _.a = 1; a + 1;", "var {a} = _; a;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var crossChunkNoInit, a = 1; a + 1;").
				AddChunk("crossChunkNoInit; a;").
				Build()),
		c.Expected("var a = _.a = 1; a + 1;", "var {crossChunkNoInit, a} = _; crossChunkNoInit; a;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var crossChunkNoInit, a = 1; a + 1;").
				AddChunk("crossChunkNoInit = 2; a;").
				Build()),
		c.Expected(
			"var crossChunkNoInit = _.crossChunkNoInit = {}; var a = _.a = 1; a + 1;",
			"var {a, crossChunkNoInit} = _; crossChunkNoInit._ = 2; a;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().AddChunk("let a; a = 1; a + 1;").AddChunk("a;").Build()),
		c.Expected("let a; a = _.a = 1; a + 1;", "var {a} = _; a;"))
	c.Test(
		c.Externs("var externVar;"),
		c.Srcs("var externVar, localVar = 1;"),
		c.Expected("var localVar = 1;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let localVar = 1, a; a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected("var localVar = 1, a; a + 1;", "var {a} = _; a;"))

	// Standalone statement writes: compound assignments, prefix/postfix increments/decrements,
	// and destructuring assignments.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a++; ++a; a--; --a; a += 2; [a] = [3]; var [b] = [4]; b;").
				AddChunk("a; b;").
				Build()),
		c.Expected(
			"var a = _.a = 1, b; ++_.a, a++; _.a = ++a; --_.a, a--; _.a = --a; _.a = a += 2;"+
				" [a] = [3]; _.a = a; [b] = [4]; _.b = b; b;",
			"var {a, b} = _; a; b;"))

	// Top-level writes inside expression contexts (binary, ternary, and control-flow conditions).
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var a = 1;"+
						" var x = (a++) + (++a) + (a--) + (--a) + (a += 2);"+
						" var y = a ? a++ : ++a;"+
						" if (a += 3) {}"+
						" while (a--) {}"+
						" do {} while (--a);"+
						" x + y;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"var a = _.a = 1;"+
				" var x = (++_.a, a++) + (_.a = ++a) + (--_.a, a--) + (_.a = --a) + (_.a = a +="+
				" 2);"+
				" var y = a ? (++_.a, a++) : (_.a = ++a);"+
				" if (_.a = a += 3) {}"+
				" while (--_.a, a--) {}"+
				" do {} while (_.a = --a);"+
				" x + y;",
			"var {a} = _; a;"))

	// Function-scoped writes inside expression contexts (return statements and concise arrow
	// bodies): because these writes occur in a nested scope (!t.inGlobalScope()), the symbol is
	// wrapped in a holder object (a._), while keeping a direct local variable (_a) in the defining
	// chunk.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var a = 1;"+
						" function postInc() { return a++; }"+
						" function preInc() { return ++a; }"+
						" function postDec() { return a--; }"+
						" function preDec() { return --a; }"+
						" var addAssign = () => a += 1;"+
						" postInc() + preInc() + postDec() + preDec() + addAssign() + a;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"_.a = {}; var _a = _.a._ = 1;"+
				" var postInc = function() { return ++_.a._, _a++; };"+
				" var preInc = function() { return _.a._ = ++_a; };"+
				" var postDec = function() { return --_.a._, _a--; };"+
				" var preDec = function() { return _.a._ = --_a; };"+
				" var addAssign = () => _.a._ = _a += 1;"+
				" postInc() + preInc() + postDec() + preDec() + addAssign() + _a;",
			"var {a} = _; a._;"))

	// Block-scoped writes (such as inside an `if` body) also have !t.inGlobalScope(), so the
	// symbol is wrapped in a holder object (a._) and cross-chunk function calls inside the block
	// observe the updated value.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; if (a + 3) { a++; }").
				AddChunk("a;").
				Build()),
		c.Expected(
			"_.a = {}; var _a = _.a._ = 1; if (_a + 3) { ++_.a._, _a++; }", "var {a} = _; a._;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function readA() { return a; }").
				AddChunk("var a = 1; if (a + 3) { a++; readA(); }").
				Build()),
		c.Expected(
			"var a = _.a = {}; _.readA = function() { return a._; };",
			"var {readA} = _; var _a = _.a._ = 1; if (_a + 3) { ++_.a._, _a++; readA(); }"))
	// Uninitialized var declarations followed by top-level compound assignment or increment.
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("var a; a += 2; a;").AddChunk("a;").Build()),
		c.Expected("var a; _.a = a += 2; a;", "var {a} = _; a;"))
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("var a; a++; a;").AddChunk("a;").Build()),
		c.Expected("var a; ++_.a, a++; a;", "var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testLocalAccessOptimization_uninitializedDeclarationWriteThrows()
func (c *RescopeGlobalSymbolsTest) TestLocalAccessOptimization_uninitializedDeclarationWriteThrows() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY
	c.assumeCrossChunkNames = false
	sources :=
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("for (var a of [1, 2]) {} a;").
				AddChunk("a;").
				Build())
	expected := c.Expected("", "")
	e := AssertThrows[jutil.IRuntimeException](func() { c.Test(sources, expected) })
	IsInstanceOf[jutil.IIllegalStateException](AssertThat(e.GetCause()))
	AssertThat(e.GetCause().GetMessage()).Contains("Expected initialized declaration for NAME a")
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappingAndInitAssignment_unsafeWriteToLocalAliasThrows()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappingAndInitAssignment_unsafeWriteToLocalAliasThrows() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	for _, chunk0 := range []string{
		"var a; label: [a] = [1]; a;",
		"var a; (([a] = [1]), 2); a;",
		"var a; if (([a] = [1])) {} a;",
	} {
		sources := c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk(chunk0).AddChunk("a;").Build())
		expected := c.Expected("", "")
		e := AssertThrows[jutil.IRuntimeException](func() { c.Test(sources, expected) })
		IsInstanceOf[jutil.IIllegalStateException](AssertThat(e.GetCause()))
		AssertThat(e.GetCause().GetMessage()).Contains("Unexpected unsafe write to local alias: NAME a")
	}
}

// Java: RescopeGlobalSymbolsTest.testAllChunksLocalAccessOptimization_independentChunks()
func (c *RescopeGlobalSymbolsTest) TestAllChunksLocalAccessOptimization_independentChunks() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 1;").
				AddChunk("function f() { return a; }").
				Build()),
		c.Expected("_.a = 1;", "var f = function() { return _.a; };"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksLocalAccessOptimization_definingChunkDependsOnUsingChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksLocalAccessOptimization_definingChunkDependsOnUsingChunk() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function f() { return a; }").
				AddChunk("var a = 1;").
				Build()),
		c.Expected("var f = function() { return _.a; };", "_.a = 1;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_reassignedOutsideStaticExecution()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_reassignedOutsideStaticExecution() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; function f() { a = 2; } a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"var a = _.a = {}; a._ = 1; var f = function() { a._ = 2; }; a._ + 1;",
			"var {a} = _; a._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_reassignedInOtherChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_reassignedInOtherChunk() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1;").
				AddChunk("a = 2; a;").
				Build()),
		c.Expected("var a = _.a = {}; a._ = 1; a._ + 1;", "var {a} = _; a._ = 2; a._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_multipleHolderSymbols()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_multipleHolderSymbols() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1;").
				AddChunk("var b = 1; b + 1;").
				AddChunk("a = 2; b = 2;").
				Build()),
		c.Expected(
			"var a = _.a = {}; a._ = 1; a._ + 1;",
			"var b = _.b = {}; b._ = 1; b._ + 1;",
			"var {a, b} = _; a._ = 2; b._ = 2;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_parentAccessesChildHolder()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_parentAccessesChildHolder() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function f() { return a; }").
				AddChunk("var a = 1; function g() { a = 2; } a + 1;").
				Build()),
		c.Expected(
			"var a = _.a = {}; var f = function() { return a._; };",
			"var {a} = _; a._ = 1; var g = function() { a._ = 2; }; a._ + 1;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_branchingCommonParent()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_branchingCommonParent() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForTree().
				AddChunk("var x = 1;").
				AddChunk("var a = 1; a + 1;").
				AddChunk("a = 2; a + 2;").
				Build()),
		c.Expected(
			"_.a = {}; var x = 1;",
			"var {a} = _; a._ = 1; a._ + 1;",
			"var {a} = _; a._ = 2; a._ + 2;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksLocalAccessOptimization_es5()
func (c *RescopeGlobalSymbolsTest) TestAllChunksLocalAccessOptimization_es5() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("var a = 1; a + 1;").AddChunk("a;").Build()),
		c.Expected("var a; a = _.a = 1; a + 1;", "var a = _.a; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_multipleHolderSymbols_es5()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_multipleHolderSymbols_es5() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1;").
				AddChunk("var b = 1; b + 1;").
				AddChunk("a = 2; b = 2;").
				Build()),
		c.Expected(
			"var a = _.a = {}; a._ = 1; a._ + 1;",
			"var b = _.b = {}; b._ = 1; b._ + 1;",
			"var a = _.a; var b = _.b; a._ = 2; b._ = 2;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_definingChunkDependsOnUsingChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_definingChunkDependsOnUsingChunk() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function f() { return a; }").
				AddChunk("var a = 1;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"var a = _.a = {}; var f = function() { return a._; };",
			"var {a} = _; a._ = 1;",
			"var {a} = _; a._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_independentSiblingChunks()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_independentSiblingChunks() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForTree().
				AddChunk("var x = 1;").
				AddChunk("var a = 1; a + 1;").
				AddChunk("function f() { return a; }").
				Build()),
		c.Expected(
			"_.a = {}; var x = 1;",
			"var {a} = _; a._ = 1; a._ + 1;",
			"var {a} = _; var f = function() { return a._; };"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_notUsedInDefiningChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_notUsedInDefiningChunk() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("var a = 1;").AddChunk("a;").Build()),
		c.Expected("_.a = 1;", "var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_externs()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_externs() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Externs("var window; var customExtern;"),
		c.Srcs(
			JSChunkGraphBuilder_ForTree().
				AddChunk("window.onload = function() { return customExtern + a; };").
				AddChunk("var a = customExtern; window.foo = a;").
				AddChunk("customExtern = window.foo + a;").
				Build()),
		c.Expected(
			"var a = _.a = {}; window.onload = function() { return customExtern + a._; };",
			"var {a} = _; a._ = customExtern; window.foo = a._;",
			"var {a} = _; customExtern = window.foo + a._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_sameVarDeclaredInExternsAndSource()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_sameVarDeclaredInExternsAndSource() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Externs("var ext;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function f() { return ext; }").
				AddChunk("var ext = 1; ext + 1;").
				Build()),
		c.Expected("var f = function() { return ext; };", "ext = 1; ext + 1;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_functionAndClassUsedInNonDependentChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_functionAndClassUsedInNonDependentChunk() {
	c.DisableCompareAsTree()
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function f() { return g() + h() + new C(); }").
				AddChunk("function g() { return 1; } function h() { return this; } class C {}").
				Build()),
		c.Expected(
			"var g=_.g={};var h=_.h={};var C=_.C={};"+
				"var f=function(){return g._()+(0,h._)()+new C._};",
			"var {g,h,C}=_;g._=function(){return 1};h._=function(){return this};C._=class{}"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_mixedWrappedAndUnwrappedLocalAliases()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_mixedWrappedAndUnwrappedLocalAliases() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1; function f() { return b; }").
				AddChunk("var b, c; b = 2; c = 3; a + b + c;").
				Build()),
		c.Expected(
			"var b = _.b = {}; var a; a = _.a = 1; a + 1; var f = function() { return b._; };",
			"var {a, b} = _; var c; b._ = 2; c = 3; a + b._ + c;"))
}

// Java: RescopeGlobalSymbolsTest.testDeleteGlobalsOnLastAccess_chain()
func (c *RescopeGlobalSymbolsTest) TestDeleteGlobalsOnLastAccess_chain() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DELETE_GLOBALS_ON_LAST_ACCESS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1; var b = 1; b + 1; var c = 1; c + 1;").
				AddChunk("a; b; c;").
				AddChunk("b;").
				Build()),
		c.Expected(
			"_._ = function(a) { for (let b of a.split(',')) { delete _[b]; } };"+
				" var a = _.a = 1; a + 1; var b = _.b = 1; b + 1; var c = _.c = 1; c + 1;",
			"var {b, a, c} = _; _._('a,c'); a; b; c;",
			"var {b} = _; _._('b'); b;"))
}

// Java: RescopeGlobalSymbolsTest.testDeleteGlobalsOnLastAccess_es5()
func (c *RescopeGlobalSymbolsTest) TestDeleteGlobalsOnLastAccess_es5() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DELETE_GLOBALS_ON_LAST_ACCESS
	c.assumeCrossChunkNames = false
	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5, CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1; var b = 1; b + 1;").
				AddChunk("a; b;").
				AddChunk("b;").
				Build()),
		c.Expected(
			"var a = _.a = 1; a + 1; var b = _.b = 1; b + 1;",
			"var b = _.b; var a = _.a; delete _.a; a; b;",
			"var b = _.b; delete _.b; b;"))
}

// Java: RescopeGlobalSymbolsTest.testDeleteGlobalsOnLastAccess_siblingChunksNotDeleted()
func (c *RescopeGlobalSymbolsTest) TestDeleteGlobalsOnLastAccess_siblingChunksNotDeleted() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DELETE_GLOBALS_ON_LAST_ACCESS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForTree().
				AddChunk("var a = 1; a + 1; var b = 1; b + 1;").
				AddChunk("a; b;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"_._ = function(a) { for (let b of a.split(',')) { delete _[b]; } };"+
				" var a = _.a = 1; a + 1; var b = _.b = 1; b + 1;",
			"var {a, b} = _; _._('b'); a; b;",
			"var {a} = _; a;"))
}

// Java: RescopeGlobalSymbolsTest.testDeleteGlobalsOnLastAccess_wrappedAndEarlyUsages()
func (c *RescopeGlobalSymbolsTest) TestDeleteGlobalsOnLastAccess_wrappedAndEarlyUsages() {
	c.optimizeLocalAccess = CompilerOptions_OptimizeLocalAccess_DELETE_GLOBALS_ON_LAST_ACCESS
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function f() { return a; } var b = 1; b + 1; var c = 1;"+
						" function g() { c = 2; } c + 1;").
				AddChunk("var a = 1; a + 1; b = 2; b; c;").
				Build()),
		c.Expected(
			"_._ = function(a) { for (let b of a.split(',')) { delete _[b]; } };"+
				" var b = _.b = {}; var a = _.a = {}; _.c = {}; var f = function() { return a._;"+
				" }; b._ = 1; b._ + 1; var _c = _.c._ = 1; var g = function() { _c = _.c._ = 2;"+
				" }; _c + 1;",
			"var {c, b} = _; _._('b'); var _a = _.a._ = 1; _a + 1; b._ = 2; b._; c._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutablesImproved_reassignedOutsideStaticExecution()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutablesImproved_reassignedOutsideStaticExecution() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; function f() { a = 2; } a + 1;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"_.a = {}; var _a = _.a._ = 1; var f = function() { _a = _.a._ = 2; }; _a + 1;",
			"var {a} = _; a._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutablesImproved_parentAccessesChildHolder()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutablesImproved_parentAccessesChildHolder() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function f() { return a; }").
				AddChunk("var a = 1; function g() { a = 2; } a + 1;").
				Build()),
		c.Expected(
			"var a = _.a = {}; var f = function() { return a._; };",
			"var _a = _.a._ = 1; var g = function() { _a = _.a._ = 2; }; _a + 1;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutablesImproved_independentSiblingChunks()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutablesImproved_independentSiblingChunks() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForTree().
				AddChunk("var x = 1;").
				AddChunk("var a = 1; a + 1;").
				AddChunk("function f() { return a; }").
				Build()),
		c.Expected(
			"_.a = {}; var x = 1;",
			"var _a = _.a._ = 1; _a + 1;",
			"var {a} = _; var f = function() { return a._; };"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutablesImproved_externs()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutablesImproved_externs() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Externs("var window; var customExtern;"),
		c.Srcs(
			JSChunkGraphBuilder_ForTree().
				AddChunk("window.onload = function() { return customExtern + a; };").
				AddChunk("var a = customExtern; window.foo = a;").
				AddChunk("customExtern = window.foo + a;").
				Build()),
		c.Expected(
			"var a = _.a = {}; window.onload = function() { return customExtern + a._; };",
			"var _a = _.a._ = customExtern; window.foo = _a;",
			"var {a} = _; customExtern = window.foo + a._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutablesImproved_mixedWrappedAndUnwrappedLocalAliases()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutablesImproved_mixedWrappedAndUnwrappedLocalAliases() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; a + 1; function f() { return b; }").
				AddChunk("var b, c; b = 2; c = 3; a + b + c;").
				Build()),
		c.Expected(
			"var b = _.b = {}; var a = _.a = 1; a + 1; var f = function() { return b._; };",
			"var {a} = _; var _b, c; _b = _.b._ = 2; c = 3; a + _b + c;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_disambiguationWithDollarSuffixLocalVar()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_disambiguationWithDollarSuffixLocalVar() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var $ = 1; function f(_) { $ = 2; return _; } $;").
				AddChunk("$;").
				Build()),
		c.Expected(
			"_.$ = {}; var __$ = _.$._ = 1;"+
				" var f = function(_$) { __$ = _.$._ = 2; return _$; }; __$;",
			"var {$} = _; $._;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var $ = 1, $$ = 2; function f(_, _$) { $ = 3; $$ = 4; return _ + _$; } $ +"+
						" $$;").
				AddChunk("$; $$;").
				Build()),
		c.Expected(
			"_.$ = {}; _.$$ = {}; var __$ = _.$._ = 1, __$$ = _.$$._ = 2; var f = function(_$,"+
				" _$$) { __$ = _.$._ = 3; __$$ = _.$$._ = 4; return _$ + _$$; }; __$ + __$$;",
			"var {$, $$} = _; $._; $$._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_disambiguationWithExistingUnderscoreNames()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_disambiguationWithExistingUnderscoreNames() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Externs("var ___a;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; function f(_a) { a = 2; return _a; } a;").
				AddChunk("function g(__a) { return a + __a; }").
				Build()),
		c.Expected(
			"_.a = {}; var ____a = _.a._ = 1;"+
				" var f = function(_a) { ____a = _.a._ = 2; return _a; }; ____a;",
			"var {a} = _; var g = function(__a) { return a._ + __a; };"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_reassignedFunctionAndCompoundWriteInDefiningChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_reassignedFunctionAndCompoundWriteInDefiningChunk() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var a = 1, b = 0, c = 0;"+
						" function clinit() { clinit = function() {}; ++a; return b++; }"+
						" const inc = () => c += 1;"+
						" clinit(); a + b + c;").
				AddChunk("clinit(); a; b; c;").
				Build()),
		c.Expected(
			"_.clinit = {}; _.a = {}; _.b = {}; _.c = {};"+
				" var _a = _.a._ = 1, _b = _.b._ = 0, _c = _.c._ = 0;"+
				" var _clinit = _.clinit._ = function() {"+
				" _clinit = _.clinit._ = function() {}; _.a._ = ++_a; return ++_.b._, _b++; };"+
				" const inc = () => _.c._ = _c += 1;"+
				" _clinit(); _a + _b + _c;",
			"var {clinit, a, b, c} = _; clinit._(); a._; b._; c._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_uninitializedVarAndDisambiguationInDefiningChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_uninitializedVarAndDisambiguationInDefiningChunk() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a; var _a = 10; function init() { a = _a; } a;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"_.a = {}; var __a; var _a = _._a = 10; var init = function() { __a = _.a._ = _a; };"+
				" __a;",
			"var {a} = _; a._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_writeOnlyInDefiningChunkNotUnwrapped()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_writeOnlyInDefiningChunkNotUnwrapped() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var a = 1; function init() { a = 2; }").
				AddChunk("a;").
				Build()),
		c.Expected(
			"var a = _.a = {}; a._ = 1; var init = function() { a._ = 2; };", "var {a} = _; a._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_nonMirrorableWriteInDefiningChunkNotUnwrapped()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_nonMirrorableWriteInDefiningChunkNotUnwrapped() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var a = 0, b = 0, c = 0; function setA() { return [a] = [1]; }"+
						" function loopB() { for (b of [1, 2]) {} }"+
						" function setC() { label: [c] = [1]; } a + b + c;").
				AddChunk("a; b; c;").
				Build()),
		c.Expected(
			"var a = _.a = {}; var b = _.b = {}; var c = _.c = {}; a._ = 0; b._ = 0; c._ = 0;"+
				" var setA = function() { return [a._] = [1]; };"+
				" var loopB = function() { for (b._ of [1, 2]) {} };"+
				" var setC = function() { label: [c._] = [1]; }; a._ + b._ + c._;",
			"var {a, b, c} = _; a._; b._; c._;"))
}

// Java: RescopeGlobalSymbolsTest.testAllChunksWithWrappedMutables_mirrorableDestructuringWriteInDefiningChunk()
func (c *RescopeGlobalSymbolsTest) TestAllChunksWithWrappedMutables_mirrorableDestructuringWriteInDefiningChunk() {
	c.optimizeLocalAccess =
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS_IMPROVED
	c.assumeCrossChunkNames = false
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var [a] = [1]; function setA() { [a] = [2]; } a;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"_.a = {}; var _a; [_a] = [1]; _.a._ = _a;"+
				" var setA = function() { [_a] = [2]; _.a._ = _a; }; _a;",
			"var {a} = _; a._;"))
}
