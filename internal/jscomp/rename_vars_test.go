// Port of com/google/javascript/jscomp/RenameVarsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/RenameVarsTest.java)

/*
 * Copyright 2005 The Closure Compiler Authors.
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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RenameVarsTest: tests for RenameVars.
type RenameVarsTest struct {
	CompilerTestCaseBase

	prefix string

	previouslyUsedMap         *VariableMap
	renameVars                *RenameVars
	withClosurePass           bool
	localRenamingOnly         bool
	useGoogleCodingConvention bool
	generatePseudoNames       bool
	preferStableNames         bool
	withNormalize             bool

	// NameGenerator to use, or null for a default.
	nameGenerator *DefaultNameGenerator
}

const RenameVarsTest_DEFAULT_PREFIX = ""

// Java: RenameVarsTest()
func NewRenameVarsTest() *RenameVarsTest {
	c := &RenameVarsTest{}
	c.InitCompilerTestCaseBase(c)
	c.prefix = RenameVarsTest_DEFAULT_PREFIX
	c.previouslyUsedMap = NewVariableMap(jutil.NewLinkedHashMap[string, string]())
	c.withClosurePass = false
	c.localRenamingOnly = false
	c.useGoogleCodingConvention = true
	c.generatePseudoNames = false
	c.preferStableNames = false
	c.withNormalize = false
	c.nameGenerator = nil
	return c
}

func TestRenameVarsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRenameVarsTest() })
}

// Java: RenameVarsTest.getCodingConvention()
func (c *RenameVarsTest) GetCodingConvention() CodingConvention {
	if c.useGoogleCodingConvention {
		return NewGoogleCodingConvention()
	} else {
		return CodingConventions_GetDefault()
	}
}

// Java: RenameVarsTest.getProcessor(Compiler)
func (c *RenameVarsTest) GetProcessor(compiler *Compiler) CompilerPass {
	var pass CompilerPass
	if c.withClosurePass {
		pass = NewRenameVarsTest_ClosurePassAndRenameVars(c, compiler)
	} else if c.nameGenerator != nil {
		c.renameVars =
			NewRenameVars(
				compiler,
				&c.prefix,
				c.localRenamingOnly,
				c.generatePseudoNames,
				c.preferStableNames,
				c.previouslyUsedMap,
				jutil.NewLinkedHashSet[uint16](),
				nil,
				c.nameGenerator)
		pass = c.renameVars
	} else {
		c.renameVars =
			NewRenameVars(
				compiler,
				&c.prefix,
				c.localRenamingOnly,
				c.generatePseudoNames,
				c.preferStableNames,
				c.previouslyUsedMap,
				jutil.NewLinkedHashSet[uint16](),
				nil,
				NewDefaultNameGenerator())
		pass = c.renameVars
	}

	if c.withNormalize {
		// Don't use the standard CompilerTestCase normalization options
		// as renaming is a post denormalization operation, but we do still
		// want to run the normal normalizations on the input in some cases.
		pass = NewRenameVarsTest_NormalizePassWrapper(compiler, pass)
	}

	return pass
}

// Java: RenameVarsTest.setUp()
func (c *RenameVarsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.DisableValidateAstChangeMarking()
	c.previouslyUsedMap = NewVariableMap(jutil.NewLinkedHashMap[string, string]())
	c.prefix = RenameVarsTest_DEFAULT_PREFIX
	c.withClosurePass = false
	c.withNormalize = false
	c.localRenamingOnly = false
	c.generatePseudoNames = false
	c.preferStableNames = false
	c.nameGenerator = nil
}

// Java: RenameVarsTest.testRenameSimple()
func (c *RenameVarsTest) TestRenameSimple() {
	c.Test(
		"function foo(v1, v2) {\n"+
			"  return v1;\n"+
			"}\n"+
			"foo();\n",
		"function a(b, c) {\n"+
			"  return b;\n"+
			"}\n"+
			"a();\n")

	// Do a sanity check on the source info
	lastScript :=
		c.GetLastCompiler().
			GetRoot().
			GetLastChild(). // first child is externs tree, last / second is sources
			GetLastChild()
	lastScriptName := lastScript.GetSourceFileName()

	// function foo(v1, v2) { ... }
	fnSubject :=
		AssertNode(lastScript.GetFirstChild()).
			IsFunction().
			HasSourceFileName(lastScriptName).
			HasLineno(1).
			HasCharno(0).
			HasLength(37)

	// function foo(v1, v2) {
	//          ^^^
	fnSubject.
		HasFirstChildThat().
		IsName("a").
		HasSourceFileName(lastScriptName).
		HasLineno(1).
		HasCharno(9).
		HasLength(3)

	// function foo(v1, v2) {
	//             ^^^^^^^^
	paramListSubject :=
		fnSubject.
			HasSecondChildThat().
			IsParamList().
			HasSourceFileName(lastScriptName).
			HasLineno(1).
			HasCharno(12).
			HasLength(8).
			HasXChildren(2)

	// function foo(v1, v2) {
	//              ^^
	paramListSubject.
		HasFirstChildThat().
		IsName("b").
		HasSourceFileName(lastScriptName).
		HasLineno(1).
		HasCharno(13).
		HasLength(2)

	// function foo(v1, v2) {
	//                  ^^
	paramListSubject.
		HasSecondChildThat().
		IsName("c").
		HasSourceFileName(lastScriptName).
		HasLineno(1).
		HasCharno(17).
		HasLength(2)

	//  { ... }
	fnBodySubject :=
		fnSubject.
			HasLastChildThat().
			IsBlock().
			HasSourceFileName(lastScriptName).
			HasLineno(1).
			HasCharno(21).
			HasLength(16)

	// return v1;
	returnSubject :=
		fnBodySubject.
			HasOneChildThat().
			IsReturn().
			HasSourceFileName(lastScriptName).
			HasLineno(2).
			HasCharno(2).
			HasLength(10)

	// return v1;
	//        ^^
	returnSubject.HasOneChildThat().IsName("b").HasLineno(2).HasCharno(9).HasLength(2)

	// foo();
	AssertNode(lastScript).
		HasLastChildThat().
		IsExprResult(). // foo(); (statement)
		HasSourceFileName(lastScriptName).
		HasLineno(4).
		HasCharno(0).
		HasLength(6).
		HasOneChildThat(). // foo() (expression)
		IsCall().
		HasSourceFileName(lastScriptName).
		HasLineno(4).
		HasCharno(0).
		HasLength(5).
		HasOneChildThat(). // foo (function name)
		IsName("a").
		HasLineno(4).
		HasCharno(0).
		HasLength(3)
}

// Java: RenameVarsTest.testRenameGlobals()
func (c *RenameVarsTest) TestRenameGlobals() {
	c.Test(
		"var Foo; var Bar, y; function x() { Bar++; }", //
		"var   a; var   b, c; function d() {   b++; }")
}

// Java: RenameVarsTest.testRenameLocals()
func (c *RenameVarsTest) TestRenameLocals() {
	c.Test(
		"(function (v1, v2) {}); (function (v3, v4) {});",
		"(function ( a,  b) {}); (function ( a,  b) {});")
	c.Test(
		"function f1(v1, v2) {}; function f2(v3, v4) {};",
		"function  c( a,  b) {}; function  d( a,  b) {};")
}

// Java: RenameVarsTest.testRenameLocals_let()
func (c *RenameVarsTest) TestRenameLocals_let() {
	c.Test(
		"(function () { let var1 = 0; let another = 1; });",
		"(function () { let a = 0; let b = 1; });")
}

// Java: RenameVarsTest.testRenameLocals_const()
func (c *RenameVarsTest) TestRenameLocals_const() {
	c.Test(
		"(function () { const var1 = 0; const another = 1; });",
		"(function () { const a = 0; const b = 1; });")
}

// Java: RenameVarsTest.testRenameParamsWithLeadingUnderscores()
func (c *RenameVarsTest) TestRenameParamsWithLeadingUnderscores() {
	c.Test("(function (_v1, _v2) {});", "(function (a, b) {});")
}

// Java: RenameVarsTest.testRenameLocalsToSame()
func (c *RenameVarsTest) TestRenameLocalsToSame() {
	c.preferStableNames = true
	c.TestSame("(function(a) {})")
	c.TestSame("(function(a, b) {})")
	c.TestSame("(function(a, b, c) {})")
	c.TestSame("(function() { var a; })")
	c.TestSame("(function() { var a, b; })")
	c.TestSame("(function() { var a, b, c; })")
}

// Java: RenameVarsTest.testRenameRedeclaredGlobals()
func (c *RenameVarsTest) TestRenameRedeclaredGlobals() {
	c.Test(
		"function f1(v1, v2) {f1()};\n"+
			"/** @suppress {duplicate} */\n"+
			"function f1(v3, v4) {f1()};\n",
		"function a(b, c) {a()};\n"+
			"/** @suppress {duplicate} */\n"+
			"function a(b, c) {a()};\n")

	c.localRenamingOnly = true

	c.Test(
		"function f1(v1, v2) {f1()};\n"+
			"/** @suppress {duplicate} */\n"+
			"function f1(v3, v4) {f1()};\n",
		"function f1(a, b) {f1()};\n"+
			"/** @suppress {duplicate} */\n"+
			"function f1(a, b) {f1()};\n")
}

// Java: RenameVarsTest.testRecursiveFunctions1()
func (c *RenameVarsTest) TestRecursiveFunctions1() {
	c.Test(
		"var walk = function walk(node, aFunction) { walk(node, aFunction); };",
		"var    a = function    a(   b,         c) {    a(   b,         c); };")

	c.localRenamingOnly = true

	c.Test(
		"var walk = function walk(node, aFunction) {  walk(node, aFunction); };",
		"var walk = function walk(   a,         b) {  walk(   a,         b); };")
}

// Java: RenameVarsTest.testRenameLocalsClashingWithGlobals()
func (c *RenameVarsTest) TestRenameLocalsClashingWithGlobals() {
	c.Test(
		"function a(v1, v2) {return v1;} a();", //
		"function a( b,  c) {return  b;} a();")
}

// Java: RenameVarsTest.testRenameNested()
func (c *RenameVarsTest) TestRenameNested() {
	c.Test(
		"function f1(v1, v2) { (function(v3, v4) {}) }",
		"function  a( b,  c) { (function( d,  e) {}) }")
	c.Test(
		"function f1(v1, v2) { function f2(v3, v4) {} }",
		"function  a( b,  c) { function  d( e,  f) {} }")
}

// Java: RenameVarsTest.testBleedingRecursiveFunctions1()
func (c *RenameVarsTest) TestBleedingRecursiveFunctions1() {
	// On IE, bleeding functions will interfere with each other if
	// they are in the same scope. In the below example, we want to be
	// sure that a and b get separate names.
	c.Test(
		"var x = function a(x) { return x ? 1 : a(1); };\n"+
			"var y = function b(x) { return x ? 2 : b(2); };\n",
		"var c = function b(a) { return a ? 1 : b(1); };\n"+
			"var e = function d(a) { return a ? 2 : d(2); };\n")
}

// Java: RenameVarsTest.testBleedingRecursiveFunctions2()
func (c *RenameVarsTest) TestBleedingRecursiveFunctions2() {
	c.Test(
		"function f() {\n"+
			"  var x = function a(x) { return x ? 1 : a(1); };\n"+
			"  var y = function b(x) { return x ? 2 : b(2); };\n"+
			"}\n",
		"function d() {\n"+
			"  var e = function a(b) { return b ? 1 : a(1); };\n"+
			"  var f = function c(a) { return a ? 2 : c(2); };\n"+
			"}\n")
}

// Java: RenameVarsTest.testBleedingRecursiveFunctions3()
func (c *RenameVarsTest) TestBleedingRecursiveFunctions3() {
	c.Test(
		"function f() {\n"+
			"  var x = function a(x) { return x ? 1 : a(1); };\n"+
			"  var y = function b(x) { return x ? 2 : b(2); };\n"+
			"  var z = function c(x) { return x ? y : c(2); };\n"+
			"}\n",
		"function f() {\n"+
			"  var g = function a(c) { return c ? 1 : a(1); };\n"+
			"  var d = function b(a) { return a ? 2 : b(2); };\n"+
			"  var h = function e(b) { return b ? d : e(2); };\n"+
			"}\n")
}

// Java: RenameVarsTest.testBleedingFunctionInBlocks()
func (c *RenameVarsTest) TestBleedingFunctionInBlocks() {
	c.Test(
		"if (true) {\n"+
			"   var x = function a(x) {return x;}\n"+
			"}\n",
		"if (true) {\n"+
			"   var c = function b(a) {return a;}\n"+
			"}\n")
}

// Java: RenameVarsTest.testRenameWithExterns1()
func (c *RenameVarsTest) TestRenameWithExterns1() {
	externs := "var foo;"
	c.Test(
		c.Externs(externs), //
		c.Srcs("    var bar; foo(bar);"),
		c.Expected("var   a; foo(  a);"))
}

// Java: RenameVarsTest.testRenameWithExterns2()
func (c *RenameVarsTest) TestRenameWithExterns2() {
	externs := "var a;"
	c.Test(
		c.Externs(externs), //
		c.Srcs("    var b = 5"),
		c.Expected("var b = 5"))
}

// Java: RenameVarsTest.testDoNotRenameExportedName()
func (c *RenameVarsTest) TestDoNotRenameExportedName() {
	c.TestSame("_foo()")
}

// Java: RenameVarsTest.testDoNotRenameArguments()
func (c *RenameVarsTest) TestDoNotRenameArguments() {
	c.TestSame("function a() { arguments; }")
}

// Java: RenameVarsTest.testRenameWithNameOverlap()
func (c *RenameVarsTest) TestRenameWithNameOverlap() {
	c.TestSame("var a = 1; var b = 2; b + b;")
}

// Java: RenameVarsTest.testRenameWithPrefix1()
func (c *RenameVarsTest) TestRenameWithPrefix1() {
	c.prefix = "PRE_"
	c.Test(
		"function  Foo(v1, v2) {return v1}  Foo();", //
		"function PRE_( a,  b) {return  a} PRE_();")
	c.prefix = RenameVarsTest_DEFAULT_PREFIX
}

// Java: RenameVarsTest.testRenameWithPrefix2()
func (c *RenameVarsTest) TestRenameWithPrefix2() {
	c.prefix = "PRE_"
	c.Test(
		"function  Foo(v1, v2) {var v3 = v1 + v2; return v3;}  Foo();",
		"function PRE_( a,  b) {var  c =  a +  b; return  c;} PRE_();")
	c.prefix = RenameVarsTest_DEFAULT_PREFIX
}

// Java: RenameVarsTest.testRenameWithPrefix3()
func (c *RenameVarsTest) TestRenameWithPrefix3() {
	c.prefix = "a"
	c.Test(
		"function Foo() {return 1;}\n"+
			"function Bar() {\n"+
			"  var a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,\n"+
			"      A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,aa,ab;\n"+
			"  Foo();\n"+
			"} Bar();\n",
		"function a() {return 1;}\n"+
			"function aa() {\n"+
			"  var b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,A,\n"+
			"      B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,$,ba,ca;\n"+
			"  a();\n"+
			"} aa();\n")
	c.prefix = RenameVarsTest_DEFAULT_PREFIX
}

// Java: RenameVarsTest.testNamingBasedOnOrderOfOccurrence()
func (c *RenameVarsTest) TestNamingBasedOnOrderOfOccurrence() {
	c.Test(
		"var q,p,m,n,l,k; try { } catch(r) {try {} catch(s) {}}; var t = q + q;",
		"var a,b,c,d,e,f; try { } catch(g) {try {} catch(h) {}}; var i = a + a;")
	c.Test(
		"(function(A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,\n"+
			"a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,$){});\n"+
			"var a4,a3,a2,a1,b4,b3,b2,b1,ab,ac,ad,fg;function foo(){};\n",
		"(function(a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,\n"+
			"A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,$){});\n"+
			"var aa,ba,ca,da,ea,fa,ga,ha,ia,ja,ka,la;function ma(){};\n")
}

// Java: RenameVarsTest.testTryCatchLifeTime()
func (c *RenameVarsTest) TestTryCatchLifeTime() {
	c.Test(
		"var q,p,m,n,l,k; (function(r) {}); try { } catch(s) {}; var t = q + q;",
		"var a,c,d,e,f,g; (function(b) {}); try { } catch(b) {}; var h = a + a;")

	c.Test(
		"try {try {} catch(p) {}} catch(s) {};", //
		"try {try {} catch(a) {}} catch(a) {};")

	c.Test(
		"try {\n"+
			"  try {\n"+
			"  } catch(p) {\n"+
			"    try {\n"+
			"    } catch(r) {}\n"+
			"  }\n"+
			"} catch(s) {\n"+
			"  try {\n"+
			"  } catch(q) {}\n"+
			"};\n",
		"try {\n"+
			"  try {\n"+
			"  } catch(a) {\n"+
			"    try {\n"+
			"    } catch(b) {}\n"+
			"  }\n"+
			"} catch(a) {\n"+
			"  try {\n"+
			"  } catch(b) {}\n"+
			"};\n")
}

// Java: RenameVarsTest.testStableRenameSimple()
func (c *RenameVarsTest) TestStableRenameSimple() {
	expectedVariableMap := c.MakeVariableMap("Foo", "a", "L 0", "b", "L 1", "c")
	c.TestRenameMap(
		"function Foo(v1, v2) {return v1;} Foo();",
		"function   a( b,  c) {return  b;}   a();",
		expectedVariableMap)

	expectedVariableMap = c.MakeVariableMap("Foo", "a", "L 0", "b", "L 1", "c", "L 2", "d")
	c.TestRenameMapUsingOldMap(
		"function Foo(v1, v2, v3) {return v1;} Foo();",
		"function   a( b,  c,  d) {return  b;}   a();",
		expectedVariableMap)
}

// Java: RenameVarsTest.testStableRenameGlobals()
func (c *RenameVarsTest) TestStableRenameGlobals() {
	expectedVariableMap := c.MakeVariableMap("Foo", "a", "Bar", "b", "y", "c", "x", "d")
	c.TestRenameMap(
		"var Foo; var Bar, y; function x() { Bar++; }",
		"var   a; var   b, c; function d() {   b++; }",
		expectedVariableMap)

	expectedVariableMap =
		c.MakeVariableMap("Foo", "a", "Bar", "b", "y", "c", "x", "d", "Baz", "f", "L 0", "e")
	c.TestRenameMapUsingOldMap(
		"var Foo, Baz; var Bar, y; function x(R) { return R + Bar++; }",
		"var   a,   f; var   b, c; function d(e) { return e +   b++; }",
		expectedVariableMap)
}

// Java: RenameVarsTest.testStableRenameWithPointlesslyAnonymousFunctions()
func (c *RenameVarsTest) TestStableRenameWithPointlesslyAnonymousFunctions() {
	expectedVariableMap := c.MakeVariableMap("L 0", "a", "L 1", "b")
	c.TestRenameMap(
		"(function (v1, v2) {}); (function (v3, v4) {});",
		"(function ( a,  b) {}); (function ( a,  b) {});",
		expectedVariableMap)

	expectedVariableMap = c.MakeVariableMap("L 0", "a", "L 1", "b", "L 2", "c")
	c.TestRenameMapUsingOldMap(
		"(function (v0, v1, v2) {}); (function (v3, v4) {});",
		"(function ( a,  b,  c) {}); (function ( a,  b) {});",
		expectedVariableMap)
}

// Java: RenameVarsTest.testStableRenameLocalsClashingWithGlobals()
func (c *RenameVarsTest) TestStableRenameLocalsClashingWithGlobals() {
	c.Test("function a(v1, v2) {return v1;} a();", "function a(b, c) {return b;} a();")
	c.previouslyUsedMap = c.renameVars.GetVariableMap()
	c.Test(
		"function bar(){return;}function a(v1, v2) {return v1;} a();",
		"function   d(){return;}function a( b,  c) {return  b;} a();")
}

// Java: RenameVarsTest.testStableRenameNested()
func (c *RenameVarsTest) TestStableRenameNested() {
	expectedVariableMap :=
		c.MakeVariableMap("f1", "a", "L 0", "b", "L 1", "c", "L 2", "d", "L 3", "e")
	c.TestRenameMap(
		"function f1(v1, v2) { (function(v3, v4) {}) }",
		"function  a( b,  c) { (function( d,  e) {}) }",
		expectedVariableMap)

	expectedVariableMap =
		c.MakeVariableMap("f1", "a", "L 0", "b", "L 1", "c", "L 2", "d", "L 3", "e", "L 4", "f")
	c.TestRenameMapUsingOldMap(
		"function f1(v1, v2) { (function(v3, v4, v5) {}) }",
		"function  a( b,  c) { (function( d,  e,  f) {}) }",
		expectedVariableMap)
}

// Java: RenameVarsTest.testStableRenameWithExterns1()
func (c *RenameVarsTest) TestStableRenameWithExterns1() {
	externs := "var foo;"
	c.Test(
		c.Externs(externs), //
		c.Srcs("    var bar; foo(bar);"),
		c.Expected("var   a; foo(  a);"))
	c.previouslyUsedMap = c.renameVars.GetVariableMap()
	c.Test(
		c.Externs(externs), //
		c.Srcs("    var bar, baz; foo(bar, baz);"),
		c.Expected("var   a,   b; foo(  a,   b);"))
}

// Java: RenameVarsTest.testStableRenameWithExterns2()
func (c *RenameVarsTest) TestStableRenameWithExterns2() {
	externs := "var a;"
	c.Test(
		c.Externs(externs), //
		c.Srcs("    var b = 5"),
		c.Expected("var b = 5"))
	c.previouslyUsedMap = c.renameVars.GetVariableMap()
	c.Test(
		c.Externs(externs),                  //
		c.Srcs("    var b = 5, catty = 9;"), //
		c.Expected("var b = 5,     c = 9;"))
}

// Java: RenameVarsTest.testStableRenameWithNameOverlap()
func (c *RenameVarsTest) TestStableRenameWithNameOverlap() {
	c.TestSame("var a = 1; var b = 2; b + b;")
	c.previouslyUsedMap = c.renameVars.GetVariableMap()
	c.TestSame("var a = 1; var c, b = 2; b + b;")
}

// Java: RenameVarsTest.testStableRenameWithAnonymousFunctions()
func (c *RenameVarsTest) TestStableRenameWithAnonymousFunctions() {
	expectedVariableMap := c.MakeVariableMap("L 0", "a", "foo", "b")
	c.TestRenameMap(
		"function foo(bar){return bar;}foo(function(h){return h;});",
		"function   b(  a){return   a;}  b(function(a){return a;});",
		expectedVariableMap)

	expectedVariableMap = c.MakeVariableMap("foo", "b", "L 0", "a", "L 1", "c")
	c.TestRenameMapUsingOldMap(
		"function foo(bar) {return bar;}foo(function(g,h) {return g+h;});",
		"function b(a){return a}b(function(a,c){return a+c;})",
		expectedVariableMap)
}

// Java: RenameVarsTest.testStableRenameSimpleExternsChanges()
func (c *RenameVarsTest) TestStableRenameSimpleExternsChanges() {
	expectedVariableMap := c.MakeVariableMap("Foo", "a", "L 0", "b", "L 1", "c")
	c.TestRenameMap(
		"function Foo(v1, v2) {return v1;} Foo();",
		"function   a( b,  c) {return  b;}   a();",
		expectedVariableMap)

	expectedVariableMap = c.MakeVariableMap("L 0", "b", "L 1", "c", "L 2", "a")
	externs := "var Foo;"
	c.TestRenameMapUsingOldMapWithExterns(
		externs,
		"function Foo(v1, v2, v0) {return v1;} Foo();",
		"function Foo( b,  c , a) {return  b;} Foo();",
		expectedVariableMap)
}

// Java: RenameVarsTest.testStableRenameSimpleLocalNameExterned()
func (c *RenameVarsTest) TestStableRenameSimpleLocalNameExterned() {
	c.Test(
		"function Foo(v1, v2) {return v1;} Foo();", //
		"function   a( b,  c) {return  b;}   a();")

	c.previouslyUsedMap = c.renameVars.GetVariableMap()

	externs := "var b;"
	c.Test(
		c.Externs(externs),
		c.Srcs("function Foo(v1, v2) {return v1;} Foo(b);"),
		c.Expected("function a(d, c) {return d;} a(b);"))
}

// Java: RenameVarsTest.testStableRenameSimpleGlobalNameExterned()
func (c *RenameVarsTest) TestStableRenameSimpleGlobalNameExterned() {
	c.Test(
		"function Foo(v1, v2) {return v1;} Foo();", //
		"function   a( b,  c) {return  b;}   a();")

	c.previouslyUsedMap = c.renameVars.GetVariableMap()

	externs := "var Foo;"
	c.Test(
		c.Externs(externs),
		c.Srcs("function Foo(v1, v2, v0) {return v1;} Foo();"),
		c.Expected("function Foo(b, c, a) {return b;} Foo();"))
}

// Java: RenameVarsTest.testStableRenameWithPrefix1AndUnstableLocalNames()
func (c *RenameVarsTest) TestStableRenameWithPrefix1AndUnstableLocalNames() {
	c.prefix = "PRE_"
	c.Test(
		"function  Foo(v1, v2) {return v1}  Foo();", //
		"function PRE_( a,  b) {return  a} PRE_();")

	c.previouslyUsedMap = c.renameVars.GetVariableMap()

	c.prefix = "PRE_"
	c.Test(
		"function  Foo(v0, v1, v2) {return v1}  Foo();", //
		"function PRE_( a,  b,  c) {return  b} PRE_();")
}

// Java: RenameVarsTest.testStableRenameWithPrefix2()
func (c *RenameVarsTest) TestStableRenameWithPrefix2() {
	c.prefix = "a"
	c.Test(
		"function Foo() {return 1;}\n"+
			"function Bar() {\n"+
			"  var a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,\n"+
			"      A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,aa,ab;\n"+
			"  Foo();\n"+
			"} Bar();\n",
		"function a() {return 1;}\n"+
			"function aa() {\n"+
			"  var b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,A,\n"+
			"      B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,$,ba,ca;\n"+
			"  a();\n"+
			"} aa();\n")

	c.previouslyUsedMap = c.renameVars.GetVariableMap()

	c.prefix = "a"
	c.Test(
		"function Foo() {return 1;}\n"+
			"function Baz() {return 1;}\n"+
			"function Bar() {\n"+
			"  var a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,\n"+
			"      A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,aa,ab;\n"+
			"  Foo();\n"+
			"} Bar();\n",
		"function a() {return 1;}\n"+
			"function ab() {return 1;}\n"+
			"function aa() {\n"+
			"  var b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,A,\n"+
			"      B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,$,ba,ca;\n"+
			"  a();\n"+
			"} aa();\n")
}

// Java: RenameVarsTest.testContrivedExampleWhereConsistentRenamingIsWorse()
func (c *RenameVarsTest) TestContrivedExampleWhereConsistentRenamingIsWorse() {
	c.previouslyUsedMap = c.MakeVariableMap("Foo", "LongString", "L 0", "b", "L 1", "c")

	c.Test(
		"function        Foo(v1, v2) {return v1;}        Foo();",
		"function LongString( b,  c) {return  b;} LongString();")

	c.previouslyUsedMap = c.renameVars.GetVariableMap()
	expectedVariableMap := c.MakeVariableMap("Foo", "LongString", "L 0", "b", "L 1", "c")
	RenameVarsTest_AssertVariableMapsEqual(expectedVariableMap, c.previouslyUsedMap)
}

// Java: RenameVarsTest.testPrevUsedMapWithDuplicates()
func (c *RenameVarsTest) TestPrevUsedMapWithDuplicates() {
	// try { ...; throw new AssertionError(); } catch (IllegalArgumentException expected) {}
	AssertThrows[*jutil.IllegalArgumentException](func() {
		c.MakeVariableMap("Foo", "z", "Bar", "z")
		c.TestSame("")
	})
}

// Java: RenameVarsTest.testExportSimpleSymbolReservesName()
func (c *RenameVarsTest) TestExportSimpleSymbolReservesName() {
	c.Test(
		"var goog, x; goog.exportSymbol('a', x);", //
		"var    a, b;    a.exportSymbol('a', b);")
	c.withClosurePass = true
	c.Test(
		"var goog, x; goog.exportSymbol('a', x);", //
		"var    b, c;    b.exportSymbol('a', c);")
}

// Java: RenameVarsTest.testExportComplexSymbolReservesName()
func (c *RenameVarsTest) TestExportComplexSymbolReservesName() {
	c.Test(
		"var goog, x; goog.exportSymbol('a.b', x);", //
		"var    a, b;    a.exportSymbol('a.b', b);")
	c.withClosurePass = true
	c.Test(
		"var goog, x; goog.exportSymbol('a.b', x);", //
		"var    b, c;    b.exportSymbol('a.b', c);")
}

// Java: RenameVarsTest.testExportToNonStringDoesntExplode()
func (c *RenameVarsTest) TestExportToNonStringDoesntExplode() {
	c.withClosurePass = true
	c.Test(
		"var goog, a, b; goog.exportSymbol(a, b);", //
		"var    a, b, c;    a.exportSymbol(b, c);")
}

// Java: RenameVarsTest.testDollarSignSuperExport1()
func (c *RenameVarsTest) TestDollarSignSuperExport1() {
	c.useGoogleCodingConvention = false
	// See http://blickly.github.io/closure-compiler-issues/#32
	c.Test(
		"var x = function($super,duper,$fantastic){}",
		"var c = function($super,    a,        b){}")

	c.localRenamingOnly = false
	c.Test("var $super = 1", "var a = 1")

	c.useGoogleCodingConvention = true
	c.Test(
		"var x = function($super,duper,$fantastic){}", //
		"var c = function($super,    a,         b){}")
}

// Java: RenameVarsTest.testDollarSignSuperExport2()
func (c *RenameVarsTest) TestDollarSignSuperExport2() {
	c.withNormalize = true

	c.useGoogleCodingConvention = false
	// See http://blickly.github.io/closure-compiler-issues/#32
	c.Test(
		"var x = function($super,duper,$fantastic){}; var y = function($super,duper){};",
		"var c = function($super,    a,         b){}; var d = function($super,    a){};")

	c.localRenamingOnly = false
	c.Test("var $super = 1", "var a = 1")

	c.useGoogleCodingConvention = true
	c.Test(
		"var x = function($super,duper,$fantastic){}; var y = function($super,duper){};",
		"var c = function($super,    a,         b){}; var d = function($super,    a){};")
}

// Java: RenameVarsTest.testBias()
func (c *RenameVarsTest) TestBias() {
	c.nameGenerator =
		NewDefaultNameGeneratorWithReserved(jutil.NewLinkedHashSet[string](), "", jutil.NewLinkedHashSet[uint16]())
	c.nameGenerator.Favors("AAAAAAAAHH")
	c.Test("var x, y", "var A, H")
}

// Java: RenameVarsTest.testPseudoNames()
func (c *RenameVarsTest) TestPseudoNames() {
	c.generatePseudoNames = false
	// See http://blickly.github.io/closure-compiler-issues/#32
	c.Test(
		"var foo = function(a, b, c){}", //
		"var   d = function(a, b, c){}")

	c.generatePseudoNames = true
	c.Test(
		"var    foo = function(   a,    b,    c){}", //
		"var $foo$$ = function($a$$, $b$$, $c$$){}")

	c.Test(
		"var    a = function(   a,    b,    c){}", //
		"var $a$$ = function($a$$, $b$$, $c$$){}")
}

// Java: RenameVarsTest.testArrowFunctions()
func (c *RenameVarsTest) TestArrowFunctions() {
	c.Test(
		"foo => {return foo + 3;}", //
		"  a => {return   a + 3;}")

	c.Test(
		"(foo, bar) => {return foo + bar + 3;}", //
		"(  a,   b) => {return   a +   b + 3;}")
}

// Java: RenameVarsTest.testClasses()
func (c *RenameVarsTest) TestClasses() {
	c.Test(
		"class fooBar {}", //
		"class      a {}")

	c.Test(
		"class fooBar {\n"+
			"  constructor(foo, bar) {\n"+
			"    this.foo = foo;\n"+
			"    this.bar = bar;\n"+
			"  }\n"+
			"}\n"+
			"var x = new fooBar(2, 3);\n",
		"class a {\n"+
			"  constructor(b, c) {\n"+
			"    this.foo = b;\n"+
			"    this.bar = c;\n"+
			"  }\n"+
			"}\n"+
			"var d = new a(2, 3);\n")

	c.Test(
		"class fooBar {\n"+
			"  constructor(foo, bar) {\n"+
			"    this.foo = foo;\n"+
			"    this.bar = bar;\n"+
			"  }\n"+
			"  func(x) {\n"+
			"    return this.foo + x;\n"+
			"  }\n"+
			"}\n"+
			"var x = new fooBar(2,3);\n"+
			"var abcd = x.func(5);\n",
		"class b {\n"+
			"  constructor(a, c) {\n"+
			"    this.foo = a;\n"+
			"    this.bar = c;\n"+
			"  }\n"+
			"  func(a) {\n"+
			"    return this.foo + a;\n"+
			"  }\n"+
			"}\n"+
			"var d = new b(2,3);\n"+
			"var e = d.func(5);\n")
}

// Java: RenameVarsTest.testLetConst()
func (c *RenameVarsTest) TestLetConst() {
	c.Test(
		"let xyz;", //
		"let   a;")

	c.Test(
		"const xyz = 1;", //
		"const   a = 1")

	c.Test(
		"let zyx = 1; {\n"+
			"  const xyz = 1;\n"+
			"  let zyx = 2;\n"+
			"  zyx = 3;\n"+
			"}\n"+
			"let xyz = 'potato';\n"+
			"zyx = 4;\n",
		"let a = 1; {\n"+
			"  const c = 1;\n"+
			"  let b = 2;\n"+
			"  b = 3;\n"+
			"}\n"+
			"let d = 'potato';\n"+
			"a = 4;\n")
}

// Java: RenameVarsTest.testGenerators()
func (c *RenameVarsTest) TestGenerators() {
	c.Test(
		"function* gen() {\n"+
			"  var xyz = 3;\n"+
			"  yield xyz + 4;\n"+
			"}\n"+
			"gen().next()\n",
		"function* a() {\n"+
			"  var b = 3;\n"+
			"  yield b + 4;\n"+
			"}\n"+
			"a().next()\n")
}

// Java: RenameVarsTest.testForOf()
func (c *RenameVarsTest) TestForOf() {
	c.Test(
		"for (var item of items) {}", //
		"for (var    a of items) {}")
}

// Java: RenameVarsTest.testTemplateStrings()
func (c *RenameVarsTest) TestTemplateStrings() {
	c.Test(
		"var name = 'Foo';\n"+
			"`My name is ${name}`;\n",
		"var a = 'Foo';\n"+
			"`My name is ${a}`;\n")
}

// Java: RenameVarsTest.testArrayDestructuring()
func (c *RenameVarsTest) TestArrayDestructuring() {
	c.Test(
		"var [x, y, z] = [1, 2, 3];", //
		"var [a, b, c] = [1, 2, 3];")
}

// Java: RenameVarsTest.testObjectDestructuring()
func (c *RenameVarsTest) TestObjectDestructuring() {
	// TODO(sdh): Teach RenameVars to take advantage of shorthand properties by
	// building up a Map from var name strings to property name multisets.  We
	// should be able to treat this similar to the "previous names" map, where
	// we preferentially pick names with the most lined-up properties, provided
	// the property names are short (should be easy enough to do the math).
	// Note, the same property name could get different var names in different
	// scopes, so we probably need to do the comparison per scope.
	// Also, this is only relevant if language_out >= ES6.
	c.Test(
		"var obj = {p: 5, h: false};\n"+
			"var {p, h} = obj;\n",
		"var a = {p: 5, h: false};\n"+
			"var {p: b, h: c} = a;\n")

	c.Test(
		"var obj = {p: 5, h: false};\n"+
			"var {p: x, h: y} = obj;\n",
		"var a = {p: 5, h: false};\n"+
			"var {p: b, h: c} = a;\n")
}

// Java: RenameVarsTest.testDefaultFunction()
func (c *RenameVarsTest) TestDefaultFunction() {
	c.Test(
		"function f(x, y=12) {\n"+
			"  return x * y;\n"+
			"}\n",
		"function c(a, b=12) {\n"+
			"  return a * b;\n"+
			"}\n")
}

// Java: RenameVarsTest.testRestFunction()
func (c *RenameVarsTest) TestRestFunction() {
	c.Test(
		"function f(x, ...y) {\n"+
			"  return x * y[0];\n"+
			"}\n",
		"function c(a, ...b) {\n"+
			"  return a * b[0];\n"+
			"}\n")
}

// Java: RenameVarsTest.testObjectLiterals()
func (c *RenameVarsTest) TestObjectLiterals() {
	c.Test(
		"var objSuper = {\n"+
			"  f: 'potato'\n"+
			"};\n"+
			"var obj = {\n"+
			"  __proto__: objSuper,\n"+
			"  g: false,\n"+
			"  x() {\n"+
			"    return super.f;\n"+
			"  }\n"+
			"};\n"+
			"obj.x();\n",
		"var a = {\n"+
			"  f: 'potato'\n"+
			"};\n"+
			"var b = {\n"+
			"  __proto__: a,\n"+
			"  g: false,\n"+
			"  x() {\n"+
			"    return super.f;\n"+
			"  }\n"+
			"};\n"+
			"b.x();\n")
}

// Java: RenameVarsTest.testImport1()
func (c *RenameVarsTest) TestImport1() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.Test("import name from './other.js'; use(name);", "import a from './other.js'; use(a);")

	c.Test(
		"import * as name from './other.js'; use(name);",
		"import * as a from './other.js'; use(a);")

	c.Test(
		"import {default as name} from './other.js'; use(name);",
		"import {default as a} from './other.js'; use(a);")
}

// Java: RenameVarsTest.testImport2()
func (c *RenameVarsTest) TestImport2() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.withNormalize = true
	c.Test(
		"import {name} from './other.js'; use(name);",
		"import {name as a} from './other.js'; use(a);")
}

// Java: RenameVarsTest.testPreferStableNames()
func (c *RenameVarsTest) TestPreferStableNames() {
	c.preferStableNames = true
	// Locals in scopes with too many local variables (>1000) should
	// not receive temporary names (eg, 'L 123').  These locals will
	// appear in the name maps with the same name as in the code (eg,
	// 'a0' in this case).
	c.Test(RenameVarsTest_CreateManyVarFunction(1000), (*string)(nil))
	AssertThat(c.renameVars.GetVariableMap().LookupNewName("a0")).IsNull()
	AssertThat(c.renameVars.GetVariableMap().LookupNewName("L 0")).IsEqualTo("b")
	c.Test(RenameVarsTest_CreateManyVarFunction(1001), (*string)(nil))
	AssertThat(c.renameVars.GetVariableMap().LookupNewName("a0")).IsEqualTo("b")
	AssertThat(c.renameVars.GetVariableMap().LookupNewName("L 0")).IsNull()

	// With {@code preferStableNames} off locals should
	// unconditionally receive temporary names.
	c.preferStableNames = false
	c.Test(RenameVarsTest_CreateManyVarFunction(1000), (*string)(nil))
	AssertThat(c.renameVars.GetVariableMap().LookupNewName("a0")).IsNull()
	AssertThat(c.renameVars.GetVariableMap().LookupNewName("L 0")).IsEqualTo("b")
	c.Test(RenameVarsTest_CreateManyVarFunction(1001), (*string)(nil))
	AssertThat(c.renameVars.GetVariableMap().LookupNewName("a0")).IsNull()
	AssertThat(c.renameVars.GetVariableMap().LookupNewName("L 0")).IsEqualTo("b")
}

// Java: RenameVarsTest.testRenameMapUsingOldMap(String, String, VariableMap)
func (c *RenameVarsTest) TestRenameMapUsingOldMap(input string, expected string, expectedMap *VariableMap) {
	c.previouslyUsedMap = c.renameVars.GetVariableMap()
	c.TestRenameMapWithExterns("", input, expected, expectedMap)
}

// Java: RenameVarsTest.testRenameMapUsingOldMap(String, String, String, VariableMap)
func (c *RenameVarsTest) TestRenameMapUsingOldMapWithExterns(externs string, input string, expected string, expectedMap *VariableMap) {
	c.previouslyUsedMap = c.renameVars.GetVariableMap()
	c.TestRenameMapWithExterns(externs, input, expected, expectedMap)
}

// Java: RenameVarsTest.testRenameMap(String, String, VariableMap)
func (c *RenameVarsTest) TestRenameMap(input string, expected string, expectedRenameMap *VariableMap) {
	c.TestRenameMapWithExterns("", input, expected, expectedRenameMap)
}

// Java: RenameVarsTest.testRenameMap(String, String, String, VariableMap)
func (c *RenameVarsTest) TestRenameMapWithExterns(externs string, input string, expected string, expectedRenameMap *VariableMap) {
	c.Test(c.Externs(externs), c.Srcs(input), c.Expected(expected))
	renameMap := c.renameVars.GetVariableMap()
	RenameVarsTest_AssertVariableMapsEqual(expectedRenameMap, renameMap)
}

// Java: RenameVarsTest.createManyVarFunction(int)
func RenameVarsTest_CreateManyVarFunction(numVars int) string {
	var locals []string
	for i := 0; i < numVars; i++ {
		locals = append(locals, "a"+jutil.ToJavaString(i))
	}
	return "function foo() { var " + strings.Join(locals, ",") + "; }"
}

// Java: RenameVarsTest.makeVariableMap(String...)
func (c *RenameVarsTest) MakeVariableMap(keyValPairs ...string) *VariableMap {
	jutil.CheckArgument(len(keyValPairs)%2 == 0)

	renameMap := jutil.NewImmutableMapBuilder[string, string]()
	for i := 0; i < len(keyValPairs); i += 2 {
		renameMap.Put(keyValPairs[i], keyValPairs[i+1])
	}

	return NewVariableMap(renameMap.BuildOrThrow())
}

// Java: RenameVarsTest.assertVariableMapsEqual(VariableMap, VariableMap)
func RenameVarsTest_AssertVariableMapsEqual(a *VariableMap, b *VariableMap) {
	ma := a.GetOriginalNameToNewNameMap()
	mb := b.GetOriginalNameToNewNameMap()
	AssertWithMessage("VariableMaps not equal").That(mb).IsEqualTo(ma)
}

// RenameVarsTest_ClosurePassAndRenameVars is the inner class
// ClosurePassAndRenameVars.
type RenameVarsTest_ClosurePassAndRenameVars struct {
	outer    *RenameVarsTest
	compiler *Compiler
}

// Java: RenameVarsTest.ClosurePassAndRenameVars(Compiler)
func NewRenameVarsTest_ClosurePassAndRenameVars(outer *RenameVarsTest, compiler *Compiler) *RenameVarsTest_ClosurePassAndRenameVars {
	return &RenameVarsTest_ClosurePassAndRenameVars{outer: outer, compiler: compiler}
}

// Java: RenameVarsTest.ClosurePassAndRenameVars.process(Node, Node)
func (p *RenameVarsTest_ClosurePassAndRenameVars) Process(externs rhino.Node, root rhino.Node) {
	NewGatherModuleMetadata(
		p.compiler, false /* processCommonJsModules= */, ModuleLoader_ResolutionMode_BROWSER).
		Process(externs, root)
	closurePass := NewProcessClosureProvidesAndRequires(p.compiler, true)
	closurePass.Process(externs, root)
	p.outer.renameVars =
		NewRenameVars(
			p.compiler,
			&p.outer.prefix,
			false,
			false,
			false,
			p.outer.previouslyUsedMap,
			jutil.NewLinkedHashSet[uint16](),
			closurePass.GetExportedVariableNames(),
			NewDefaultNameGenerator())
	p.outer.renameVars.Process(externs, root)
}

// RenameVarsTest_NormalizePassWrapper is the nested class
// NormalizePassWrapper.
type RenameVarsTest_NormalizePassWrapper struct {
	compiler    *Compiler
	wrappedPass CompilerPass
}

// Java: RenameVarsTest.NormalizePassWrapper(Compiler, CompilerPass)
func NewRenameVarsTest_NormalizePassWrapper(compiler *Compiler, wrappedPass CompilerPass) *RenameVarsTest_NormalizePassWrapper {
	return &RenameVarsTest_NormalizePassWrapper{compiler: compiler, wrappedPass: wrappedPass}
}

// Java: RenameVarsTest.NormalizePassWrapper.process(Node, Node)
func (p *RenameVarsTest_NormalizePassWrapper) Process(externs rhino.Node, root rhino.Node) {
	normalize := Normalize_CreateNormalizeForOptimizations(p.compiler)
	normalize.Process(externs, root)

	p.wrappedPass.Process(externs, root)
}
