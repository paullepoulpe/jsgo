// Port of com/google/javascript/jscomp/MakeDeclaredNamesUniqueTest.java @ closure-compiler v20260415

package jscomp

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// MakeDeclaredNamesUniqueTest extends CompilerTestCase, which is above the
// core: its test/testSame/srcs/expected/externs helpers skip the test. The
// cases and the pass the Java test installs (getProcessor) are ported in
// full; core_g2_oracle_test.go runs the pass on every JavaScript string of
// this test against Java (mode mdnu).
const makeDeclaredNamesUniqueTest_skip = "needs CompilerTestCase"

const makeDeclaredNamesUniqueTest_LOCAL_NAME_PREFIX = "unique_"

// makeDeclaredNamesUniqueTest is MakeDeclaredNamesUniqueTest's fixture.
type makeDeclaredNamesUniqueTest struct {
	t *testing.T
	// useDefaultRenamer = true; invokes the ContextualRenamer
	// useDefaultRenamer = false; invokes the InlineRenamer
	useDefaultRenamer bool
	// invert = true; treats JavaScript input as normalized code and inverts the renaming
	// invert = false; conducts renaming
	invert bool
	// removeConst = true; removes const-ness of a name (e.g. If the variable name is CONST)
	removeConst bool
	// whether to throw an exception on any names newly made unique.
	assertOnChange bool
}

// Java: MakeDeclaredNamesUniqueTest.customSetUp()
func newMakeDeclaredNamesUniqueTest(t *testing.T) *makeDeclaredNamesUniqueTest {
	return &makeDeclaredNamesUniqueTest{t: t}
}

// makeDeclaredNamesUniqueTest_pass is the anonymous CompilerPass of
// getProcessor.
type makeDeclaredNamesUniqueTest_pass struct {
	x        *makeDeclaredNamesUniqueTest
	compiler AbstractCompiler
}

// Java: MakeDeclaredNamesUniqueTest.getProcessor(...).new CompilerPass() {...}.process(Node, Node)
func (p *makeDeclaredNamesUniqueTest_pass) Process(externs *rhino.Node, root *rhino.Node) {
	renamer := MakeDeclaredNamesUnique_Builder_().WithAssertOnChange(p.x.assertOnChange)
	if !p.x.useDefaultRenamer {
		renamer = renamer.WithRenamer(
			NewMakeDeclaredNamesUnique_InlineRenamer(
				p.compiler.GetCodingConvention(),
				p.compiler.GetUniqueNameIdSupplier(),
				makeDeclaredNamesUniqueTest_LOCAL_NAME_PREFIX,
				p.x.removeConst,
				true,
				nil))
	}
	NodeTraversal_TraverseRoots(p.compiler, renamer.Build(), externs, root)
}

// Java: MakeDeclaredNamesUniqueTest.getProcessor(Compiler)
func (x *makeDeclaredNamesUniqueTest) getProcessor(compiler AbstractCompiler) CompilerPass {
	if !x.invert {
		return &makeDeclaredNamesUniqueTest_pass{x: x, compiler: compiler}
	}
	return MakeDeclaredNamesUnique_GetContextualRenameInverter(compiler)
}

// test is CompilerTestCase.test (every overload).
func (x *makeDeclaredNamesUniqueTest) test(args ...any) {
	x.t.Skip(makeDeclaredNamesUniqueTest_skip)
}

// testSame is CompilerTestCase.testSame (every overload).
func (x *makeDeclaredNamesUniqueTest) testSame(args ...any) {
	x.t.Skip(makeDeclaredNamesUniqueTest_skip)
}

// testNoWarning is CompilerTestCase.testNoWarning(String).
func (x *makeDeclaredNamesUniqueTest) testNoWarning(js string) {
	x.t.Skip(makeDeclaredNamesUniqueTest_skip)
}

// srcs is CompilerTestCase.srcs (every overload).
func (x *makeDeclaredNamesUniqueTest) srcs(args ...any) any {
	x.t.Skip(makeDeclaredNamesUniqueTest_skip)
	return nil
}

// expected is CompilerTestCase.expected (every overload).
func (x *makeDeclaredNamesUniqueTest) expected(args ...any) any {
	x.t.Skip(makeDeclaredNamesUniqueTest_skip)
	return nil
}

// externs is CompilerTestCase.externs (every overload).
func (x *makeDeclaredNamesUniqueTest) externs(args ...any) any {
	x.t.Skip(makeDeclaredNamesUniqueTest_skip)
	return nil
}

// testWithInversion covers testWithInversion(String, String) and
// testWithInversion(String[], String[]).
//
// Java: MakeDeclaredNamesUniqueTest.testWithInversion(String, String)
func (x *makeDeclaredNamesUniqueTest) testWithInversion(original any, expected any) {
	if o, ok := original.([]string); ok {
		e := expected.([]string)
		x.invert = false
		x.test(x.srcs(o), x.expected(e))
		x.invert = true
		x.test(x.srcs(e), x.expected(o))
		x.invert = false
		return
	}
	x.invert = false
	x.test(original, expected)
	x.invert = true
	x.test(expected, original)
	x.invert = false
}

// testSameWithInversion covers testSameWithInversion(String, String) and
// testSameWithInversion(String).
//
// Java: MakeDeclaredNamesUniqueTest.testSameWithInversion(String, String)
func (x *makeDeclaredNamesUniqueTest) testSameWithInversion(args ...string) {
	externs, original := "", args[0]
	if len(args) == 2 {
		externs, original = args[0], args[1]
	}
	x.invert = false
	x.testSame(x.externs(externs), x.srcs(original))
	x.invert = true
	x.testSame(x.externs(externs), x.srcs(original))
	x.invert = false
}

// Java: MakeDeclaredNamesUniqueTest.wrapInFunction(String)
func makeDeclaredNamesUniqueTest_wrapInFunction(s string) string {
	return "function f(){" + s + "}"
}

// Java: MakeDeclaredNamesUniqueTest.testInFunctionWithInversion(String, String)
func (x *makeDeclaredNamesUniqueTest) testInFunctionWithInversion(original string, expected string) {
	x.testWithInversion(
		makeDeclaredNamesUniqueTest_wrapInFunction(original), makeDeclaredNamesUniqueTest_wrapInFunction(expected))
}

// Go-only: test-helper — exercise the renamers directly (no Java test).
func TestMakeDeclaredNamesUnique_renamers(t *testing.T) {
	root := NewMakeDeclaredNamesUnique_ContextualRenamer()
	root.AddDeclaredName("a", false)
	if got := root.GetReplacementName("a"); got != nil {
		t.Errorf("global a renamed to %q", *got)
	}
	fn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), rhino.IR_Block())
	child := root.CreateForChildScope(fn, false)
	child.AddDeclaredName("a", false)
	child.AddDeclaredName("arguments", false)
	if got := child.GetReplacementName("a"); got == nil || *got != "a$jscomp$1" {
		t.Errorf("local a: got %v", got)
	}
	if got := child.GetReplacementName("arguments"); got != nil {
		t.Errorf("arguments renamed to %q", *got)
	}
	if !strings.HasPrefix(root.String(), "ContextualRenamer{scopeRoot=null, nameUsage=[") {
		t.Errorf("toString: %s", root.String())
	}
	if got := MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName("x$jscomp$12"); got != "x" {
		t.Errorf("getOriginalName: %q", got)
	}
	if got := MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName("$jscomp$12"); got != "$jscomp$12" {
		t.Errorf("getOriginalName: %q", got)
	}
	targetted := NewMakeDeclaredNamesUnique_TargettedRenamer(
		NewMakeDeclaredNamesUnique_ContextualRenamer(), jutil.NewLinkedHashSet("t"))
	targetted.AddDeclaredName("t", false)
	targetted.AddDeclaredName("u", false)
	if targetted.GetReplacementName("u") != nil {
		t.Errorf("untargetted name replaced")
	}
}

// Java: MakeDeclaredNamesUniqueTest.testMakeDeclaredNamesUniqueNullishCoalesce()
func TestMakeDeclaredNamesUnique_testMakeDeclaredNamesUniqueNullishCoalesce(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true

	x.test(
		"var foo; var x = function foo(){var foo = false ?? {};}",
		"var foo; var x = function foo$jscomp$1(){var foo$jscomp$2 = false ?? {}}")
	x.testSameWithInversion("var a = b ?? c;")
}

// Java: MakeDeclaredNamesUniqueTest.testShadowedBleedingName()
func TestMakeDeclaredNamesUnique_testShadowedBleedingName(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true

	x.test(
		"var foo; var x = function foo(){var foo;}",
		"var foo; var x = function foo$jscomp$1(){var foo$jscomp$2}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext1()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext1(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true

	x.invert = true
	x.test(
		"var a;function foo(){var a$jscomp$inline_1; a = 1}",
		"var a;function foo(){var a$jscomp$0       ; a = 1}")
	x.test(
		"var a;function foo(){var a$jscomp$inline_1;}", //
		"var a;function foo(){var a                ;}")

	x.test(
		"let a;function foo(){let a$jscomp$inline_1; a = 1}",
		"let a;function foo(){let a$jscomp$0       ; a = 1}")
	x.test(
		"const a = 1;function foo(){let a$jscomp$inline_1;}", //
		"const a = 1;function foo(){let a                ;}")
	x.test(
		"class A {} function foo(){class A$jscomp$inline_1 {}}",
		"class A {} function foo(){class A {}}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext2()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext2(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	// Set the test type
	x.useDefaultRenamer = true

	// Verify global names are untouched.
	x.testSameWithInversion("var a;")
	x.testSameWithInversion("let a;")
	x.testSameWithInversion("const a = 0;")

	// Verify global names are untouched.
	x.testSameWithInversion("a;")

	// Local names are made unique.
	x.testWithInversion(
		"var a;function foo(a){var b;a}", "var a;function foo(a$jscomp$1){var b;a$jscomp$1}")
	x.testWithInversion(
		"var a;function foo(){var b;a}function boo(){var b;a}",
		"var a;function foo(){var b;a}function boo(){var b$jscomp$1;a}")
	x.testWithInversion(
		`function foo(a){var b}
function boo(a){var b}
`,
		`function foo(a){var b}
function boo(a$jscomp$1){var b$jscomp$1}
`)
	// variable b is left untouched because it is only declared once
	x.testWithInversion(
		"let a;function foo(a){let b;a}", "let a;function foo(a$jscomp$1){let b;a$jscomp$1}")
	x.testWithInversion(
		"let a;function foo(){let b;a}function boo(){let b;a}",
		"let a;function foo(){let b;a}function boo(){let b$jscomp$1;a}")
	x.testWithInversion(
		`function foo(a){let b}
function boo(a){let b}
`,
		`function foo(a){let b}
function boo(a$jscomp$1){let b$jscomp$1}
`)

	// Verify functions expressions are renamed.
	x.testWithInversion(
		"var a = function foo(){foo()};var b = function foo(){foo()};",
		"var a = function foo(){foo()};var b = function foo$jscomp$1(){foo$jscomp$1()};")
	x.testWithInversion(
		"let a = function foo(){foo()};let b = function foo(){foo()};",
		"let a = function foo(){foo()};let b = function foo$jscomp$1(){foo$jscomp$1()};")

	// Verify catch exceptions names are made unique
	x.testSameWithInversion("try { } catch(e) {e;}")

	// Inversion does not handle exceptions correctly.
	x.test(
		"try { } catch(e) {e;}; try { } catch(e) {e;}",
		"try { } catch(e) {e;}; try { } catch(e$jscomp$1) {e$jscomp$1;}")
	x.test(
		"try { } catch(e) {e; try { } catch(e) {e;}};",
		"try { } catch(e) {e; try { } catch(e$jscomp$1) {e$jscomp$1;} }; ")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext3()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext3(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	// Set the test type
	x.useDefaultRenamer = true

	externs := "var extern1 = {};"

	// Verify global names are untouched.
	x.testSameWithInversion(externs, "var extern1 = extern1 || {};")

	// Verify global names are untouched.
	x.testSame(x.externs(externs), x.srcs("var extern1 = extern1 || {};"))
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext4()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext4(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	// Set the test type
	x.useDefaultRenamer = true

	x.testInFunctionWithInversion(
		"var e; try { } catch(e) {e;}; try { } catch(e) {e;}",
		"var e; try { } catch(e$jscomp$1) {e$jscomp$1;}; try { } catch(e$jscomp$2) {e$jscomp$2;}")
	x.testInFunctionWithInversion(
		"var e; try { } catch(e         ) {         e; try { } catch(e         ) {e;         } };",
		"var e; try { } catch(e$jscomp$1) {e$jscomp$1; try { } catch(e$jscomp$2) {e$jscomp$2;} };")
	x.testInFunctionWithInversion(
		"try { } catch(e         ) {e         ;}; try { } catch(e         ) {e         ;}; var e;",
		"try { } catch(e$jscomp$1) {e$jscomp$1;}; try { } catch(e$jscomp$2) {e$jscomp$2;}; var e;")
	x.testInFunctionWithInversion(
		"try { } catch(e         ) {e         ; try { } catch(e         ) {e         ;} }; var e;",
		"try { } catch(e$jscomp$1) {e$jscomp$1; try { } catch(e$jscomp$2) {e$jscomp$2;} }; var e;")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext5()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext5(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion("function f(){var f; f = 1}", "function f(){var f$jscomp$1; f$jscomp$1 = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext6()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext6(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion("function f(f){f = 1}", "function f(f$jscomp$1){f$jscomp$1 = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext7()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext7(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		"function f(f){var f; f = 1}", "function f(f$jscomp$1){var f$jscomp$1; f$jscomp$1 = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext8()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext8(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.test(
		"var fn = function f(){var f; f = 1}",
		"var fn = function f(){var f$jscomp$1; f$jscomp$1 = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext9()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext9(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSame("var fn = function f(f){f = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext10()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext10(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSame("var fn = function f(f){var f; f = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext11()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithContext11(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	// Changes the parameter name if it's the same as loop object name
	x.test(
		`var loopObjName = {};
for(;;) {
   var fn = (function f(loopObjName){ return function() {}; })(loopObjName);
}
`,
		`var loopObjName = {};
for(;;) {
   var fn = (function f(loopObjName$jscomp$1){ return function() {}; })(loopObjName);
}
`)

	// parameter name is already unique; no change
	x.testSame(
		`var foo = {};
for(;;) {
   var fn = (function f(bar){ return function() {}; })(foo);}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testMakeFunctionsUniqueWithContext()
func TestMakeDeclaredNamesUnique_testMakeFunctionsUniqueWithContext(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSame("function f(){} function f(){}")
	x.testSame("var x = function() {function f(){} function f(){}};")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeFunctionsUniqueWithContext1()
func TestMakeDeclaredNamesUnique_testMakeFunctionsUniqueWithContext1(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.test(
		"if (1) { function f(){} } else { function f(){} }",
		"if (1) { function f(){} } else { function f$jscomp$1(){} }")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeFunctionsUniqueWithContext2()
func TestMakeDeclaredNamesUnique_testMakeFunctionsUniqueWithContext2(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSame("if (1) { function f(){} function f(){} }")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeFunctionsUniqueWithContext3()
func TestMakeDeclaredNamesUnique_testMakeFunctionsUniqueWithContext3(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.test(
		"function f() {} if (1) { function f(){} function f(){} }",
		"function f() {} if (1) { function f$jscomp$1(){} function f$jscomp$1(){} }")
}

// Java: MakeDeclaredNamesUniqueTest.testArguments()
func TestMakeDeclaredNamesUnique_testArguments(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	// Set the test type
	x.useDefaultRenamer = true

	// Don't distinguish between "arguments", it can't be made unique.
	x.testSameWithInversion("function foo(){var arguments;function bar(){var arguments;}}")

	x.invert = true

	// Don't introduce new references to arguments, it is special.
	// Still try to rename it to a name that depends on the shape of the AST rather than
	// the process we happened to follow to reach that shape.
	x.test(
		"function foo(){var arguments$jscomp$1;}", //
		"function foo(){var arguments$jscomp$0;}")
}

// Java: MakeDeclaredNamesUniqueTest.testClassInForLoop()
func TestMakeDeclaredNamesUnique_testClassInForLoop(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSame("for (class a {};;) { break; }")
}

// Java: MakeDeclaredNamesUniqueTest.testFunctionInForLoop()
func TestMakeDeclaredNamesUnique_testFunctionInForLoop(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSame("for (function a() {};;) { break; }")
}

// Java: MakeDeclaredNamesUniqueTest.testLetsInSeparateBlocks()
func TestMakeDeclaredNamesUnique_testLetsInSeparateBlocks(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.test(
		`if (x) {
  let e;
  alert(e);
}
if (y) {
  let e;
  alert(e);
}
`,
		`if (x) {
  let e;
  alert(e);
}
if (y) {
  let e$jscomp$1;
  alert(e$jscomp$1);
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testConstInGlobalHoistScope()
func TestMakeDeclaredNamesUnique_testConstInGlobalHoistScope(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSame(
		`if (true) {
  const x = 1; alert(x);
}
`)

	x.test(
		`if (true) {
  const x = 1; alert(x);
} else {
  const x = 1; alert(x);
}
`,
		`if (true) {
  const x = 1; alert(x);
} else {
  const x$jscomp$1 = 1; alert(x$jscomp$1);
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithoutContext()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithoutContext(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = false

	x.test("var a;", "var a$jscomp$unique_0")
	x.test("let a;", "let a$jscomp$unique_0")

	// Verify undeclared names are untouched.
	x.testSame("a;")

	// Local names are made unique.
	x.test(
		`var a;
function foo(a){var b;a}
`,
		`var a$jscomp$unique_0;
function foo$jscomp$unique_1(a$jscomp$unique_2){
  var b$jscomp$unique_3;a$jscomp$unique_2}
`)
	x.test(
		`var a;
function foo(){var b;a}
function boo(){var b;a}
`,
		`var a$jscomp$unique_0;
function foo$jscomp$unique_1(){var b$jscomp$unique_3;a$jscomp$unique_0}
function boo$jscomp$unique_2(){var b$jscomp$unique_4;a$jscomp$unique_0}
`)

	x.test(
		"let a; function foo(a) {let b; a; }",
		`let a$jscomp$unique_0;
function foo$jscomp$unique_1(a$jscomp$unique_2) {
  let b$jscomp$unique_3;
  a$jscomp$unique_2;
}
`)

	x.test(
		`let a;
function foo() { let b; a; }
function boo() { let b; a; }
`,
		`let a$jscomp$unique_0;
function foo$jscomp$unique_1() {
  let b$jscomp$unique_3;
  a$jscomp$unique_0;
}
function boo$jscomp$unique_2() {
  let b$jscomp$unique_4;
  a$jscomp$unique_0;
}
`)

	// Verify function expressions are renamed.
	x.test(
		"var a = function foo(){foo()};",
		"var a$jscomp$unique_0 = function foo$jscomp$unique_1(){foo$jscomp$unique_1()};")
	x.test(
		"const a = function foo(){foo()};",
		"const a$jscomp$unique_0 = function foo$jscomp$unique_1(){foo$jscomp$unique_1()};")

	// Verify catch exceptions names are made unique
	x.test("try { } catch(e) {e;}", "try { } catch(e$jscomp$unique_0) {e$jscomp$unique_0;}")
	x.test(
		`try { } catch(e) {e;};
try { } catch(e) {e;}
`,
		`try { } catch(e$jscomp$unique_0) {e$jscomp$unique_0;};
try { } catch(e$jscomp$unique_1) {e$jscomp$unique_1;}
`)
	x.test(
		`try { } catch(e) {e;
try { } catch(e) {e;}};
`,
		`try { } catch(e$jscomp$unique_0) {e$jscomp$unique_0;
try { } catch(e$jscomp$unique_1) {e$jscomp$unique_1;} }; 
`)
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithoutContext2()
func TestMakeDeclaredNamesUnique_testMakeLocalNamesUniqueWithoutContext2(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	// Set the test type
	x.useDefaultRenamer = false

	x.test("var _a;", "var JSCompiler__a$jscomp$unique_0")
	x.test(
		"var _a = function _b(_c) { var _d; };",
		`var JSCompiler__a$jscomp$unique_0 = function JSCompiler__b$jscomp$unique_2(
JSCompiler__c$jscomp$unique_1) { var JSCompiler__d$jscomp$unique_3; };
`)

	x.test("let _a;", "let JSCompiler__a$jscomp$unique_0")
	x.test(
		"const _a = function _b(_c) { let _d; };",
		`const JSCompiler__a$jscomp$unique_0 = function JSCompiler__b$jscomp$unique_2(
JSCompiler__c$jscomp$unique_1) { let JSCompiler__d$jscomp$unique_3; };
`)
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion()
func TestMakeDeclaredNamesUnique_testOnlyInversion(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.invert = true
	x.test(
		"function f(a, a$jscomp$1) {}", //
		"function f(a, a$jscomp$0) {}")
	x.test(
		"function f(a$jscomp$1, b$jscomp$2) {}", //
		"function f(a         , b         ) {}")
	x.test(
		"function f(a$jscomp$1, a$jscomp$2) {}", //
		"function f(a         , a$jscomp$0) {}")
	x.test(
		"try { } catch(e) {e; try { } catch(e$jscomp$1) {e$jscomp$1;} }; ",
		"try { } catch(e) {e; try { } catch(e) {e;} }; ")
	x.testSame("var a$jscomp$1;")
	x.testSame("const a$jscomp$1 = 1;")
	x.testSame("function f() { var $jscomp$; }")
	x.testSame("var CONST = 3; var b = CONST;")
	x.test(
		"function f() {var CONST = 3; var ACONST$jscomp$1 = 2;}",
		"function f() {var CONST = 3; var ACONST = 2;}")
	x.test(
		"function f() {const CONST = 3; const ACONST$jscomp$1 = 2;}",
		"function f() {const CONST = 3; const ACONST = 2;}")
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion2()
func TestMakeDeclaredNamesUnique_testOnlyInversion2(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.invert = true
	x.test(
		"function f() {try { } catch(e) {e;}; try { } catch(e$jscomp$0) {e$jscomp$0;}}",
		"function f() {try { } catch(e) {e;}; try { } catch(e) {e;}}")
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion3()
func TestMakeDeclaredNamesUnique_testOnlyInversion3(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.invert = true
	x.test(
		`function x1() {
  var a$jscomp$1;
  function x2() {
    var a$jscomp$2;
  }
  function x3() {
    var a$jscomp$3;
  }
}
`,
		`function x1() {
  var a;
  function x2() {
    var a;
  }
  function x3() {
    var a;
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion4()
func TestMakeDeclaredNamesUnique_testOnlyInversion4(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.invert = true
	x.testSame(
		`function x1() {
// The attempt to rename will re-generate this same exact name.
// The purpose of this test is to make sure we don't accidentally report
// this "renaming" to the same name as a change.
  var a$jscomp$0;
  function x2() {
    var a;a$jscomp$0++
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion5()
func TestMakeDeclaredNamesUnique_testOnlyInversion5(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.invert = true
	x.test(
		`function x1() {
  const a$jscomp$1 = 0;
  function x2() {
    const b$jscomp$1 = 0;
  }
}
`,
		`function x1() {
  const a = 0;
  function x2() {
    const b = 0;
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testConstRemovingRename1()
func TestMakeDeclaredNamesUnique_testConstRemovingRename1(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.removeConst = true
	x.test(
		"(function () {var CONST = 3; var ACONST$jscomp$1 = 2;})",
		"(function () {var CONST$jscomp$unique_0 = 3; var ACONST$jscomp$unique_1 = 2;})")
}

// Java: MakeDeclaredNamesUniqueTest.testConstRemovingRename2()
func TestMakeDeclaredNamesUnique_testConstRemovingRename2(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.removeConst = true
	x.test(
		"var CONST = 3; var b = CONST;",
		"var CONST$jscomp$unique_0 = 3; var b$jscomp$unique_1 = CONST$jscomp$unique_0;")
}

// Java: MakeDeclaredNamesUniqueTest.testConstRemovingRenameAlsoRemovesAnnotation()
func TestMakeDeclaredNamesUnique_testConstRemovingRenameAlsoRemovesAnnotation(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.removeConst = true
	x.test(
		"/** @const */ var c = 3; var b = c;",
		"/** blank */ var c$jscomp$unique_0 = 3; var b$jscomp$unique_1 = c$jscomp$unique_0;")
}

// Java: MakeDeclaredNamesUniqueTest.testRestParamWithoutContext()
func TestMakeDeclaredNamesUnique_testRestParamWithoutContext(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.test(
		"function f(...x) { x; }",
		"function f$jscomp$unique_0(...x$jscomp$unique_1) { x$jscomp$unique_1; }")
}

// Java: MakeDeclaredNamesUniqueTest.testRestParamWithContextWithInversion()
func TestMakeDeclaredNamesUnique_testRestParamWithContextWithInversion(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		`let x = 0;
function foo(...x) {
  return x[0];
}
`,
		`let x = 0;
function foo(...x$jscomp$1) {
  return x$jscomp$1[0]
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testVarParamSameName0()
func TestMakeDeclaredNamesUnique_testVarParamSameName0(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.test(
		`function f(x) {
  if (!x) var x = 6;
}
`,
		`function f$jscomp$unique_0(x$jscomp$unique_1) {
  if (!x$jscomp$unique_1) var x$jscomp$unique_1 = 6;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testVarParamSameName1()
func TestMakeDeclaredNamesUnique_testVarParamSameName1(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.test(
		`function f(x) {
  if (!x) x = 6;
}
`,
		`function f$jscomp$unique_0(x$jscomp$unique_1) {
  if (!x$jscomp$unique_1) x$jscomp$unique_1 = 6;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testVarParamSameAsLet0()
func TestMakeDeclaredNamesUnique_testVarParamSameAsLet0(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.test(
		`function f(x) {
  if (!x) { let x = 6; }
}
`,
		`function f$jscomp$unique_0(x$jscomp$unique_1) {
  if (!x$jscomp$unique_1) { let x$jscomp$unique_2 = 6; }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testObjectProperties()
func TestMakeDeclaredNamesUnique_testObjectProperties(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.test("var a = {x : 'a'};", "var a$jscomp$unique_0 = {x : 'a'};")
	x.test("let a = {x : 'a'};", "let a$jscomp$unique_0 = {x : 'a'};")
	x.test("const a = {x : 'a'};", "const a$jscomp$unique_0 = {x : 'a'};")
	x.test("var a = {x : 'a'}; a.x", "var a$jscomp$unique_0 = {x : 'a'}; a$jscomp$unique_0.x")
}

// Java: MakeDeclaredNamesUniqueTest.testClassesWithContextWithInversion()
func TestMakeDeclaredNamesUnique_testClassesWithContextWithInversion(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		`var a;
class Foo {
  constructor(a) {
    this.a = a;
  }
  f() {
    var x = 1;
    return a + x;
  }
}
`,
		`var a;
class Foo {
  constructor(a$jscomp$1) {
    this.a = a$jscomp$1;
  }
  f() {
    var x = 1;
    return a + x;
  }
}
`)

	// class declarations are block-scoped but not hoisted.
	x.testSameWithInversion(
		`{
  let x = new Foo(); // ReferenceError
  class Foo {}
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion1()
func TestMakeDeclaredNamesUnique_testBlockScopesWithContextWithInversion1(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		`{let a;
  {
    let a;
  }}
`,
		`{let a;
  {
  let a$jscomp$1;
  }}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion2()
func TestMakeDeclaredNamesUnique_testBlockScopesWithContextWithInversion2(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	// function declarations are block-scoped
	x.testWithInversion(
		`function foo() {
  function bar() {
    return 1;
  }
}
function boo() {
  function bar() {
    return 2;
  }
}
`,
		`function foo() {
  function bar() {
    return 1;
  }
}
function boo() {
  function bar$jscomp$1() {
    return 2;
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion3()
func TestMakeDeclaredNamesUnique_testBlockScopesWithContextWithInversion3(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.test(
		`function foo() {
  function bar() {
    return 1;
  }
  if (true) {
    function bar() {
      return 2;
    }
  }
}
`,
		`function foo() {
  function bar() {
    return 1;
  }
  if (true) {
    function bar$jscomp$1() {
      return 2;
    }
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion4()
func TestMakeDeclaredNamesUnique_testBlockScopesWithContextWithInversion4(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.test(
		`var f1=function(){
  var x
};
(function() {
  function f2() {
    alert(x)
  }
  {
    var x=0
  }
  f2()
})()
`,
		`var f1=function(){
  var x
};
(function() {
  function f2() {
    alert(x$jscomp$1)
  }
  {
    var x$jscomp$1=0
  }
  f2()
})()
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion5()
func TestMakeDeclaredNamesUnique_testBlockScopesWithContextWithInversion5(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSame(
		`if (true) {
  function f(){};
}
f();
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithoutContext()
func TestMakeDeclaredNamesUnique_testBlockScopesWithoutContext(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = false
	x.test(
		`{
  function foo() {return 1;}
  if (true) {
    function foo() {return 2;}
  }
}
`,
		`{
  function foo$jscomp$unique_0() {return 1;}
  if (true) {
    function foo$jscomp$unique_1() {return 2;}
  }
}
`)

	x.test(
		`function foo(x) {
  return foo(x) - 1;
}
`,
		`function foo$jscomp$unique_0(x$jscomp$unique_1) {
  return foo$jscomp$unique_0(x$jscomp$unique_1) - 1;
}
`)

	x.test(
		`export function foo(x) {
  return foo(x) - 1;
}
`,
		`export function foo$jscomp$unique_0(x$jscomp$unique_1) {
  return foo$jscomp$unique_0(x$jscomp$unique_1) - 1;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testRecursiveFunctionsWithContextWithInversion()
func TestMakeDeclaredNamesUnique_testRecursiveFunctionsWithContextWithInversion(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testSameWithInversion(
		`function foo(x) {
  return foo(x) - 1;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testInvertShadowedParameterNames()
func TestMakeDeclaredNamesUnique_testInvertShadowedParameterNames(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		`var p;
function f(p) {
  return function g(p) {
    return p;
  }
}
`,
		`var p;
function f(p$jscomp$1) {
  return function g(p$jscomp$2) {
    return p$jscomp$2;
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testArrowFunctionWithContextWithInversion()
func TestMakeDeclaredNamesUnique_testArrowFunctionWithContextWithInversion(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		`function foo() {
  var f = (x) => x;
  return f(1);
}
function boo() {
  var f = (x) => x;
  return f(2);
}
`,
		`function foo() {
  var f = (x) => x;
  return f(1);
}
function boo() {
  var f$jscomp$1 = (x$jscomp$1) => x$jscomp$1;
  return f$jscomp$1(2);
}
`)

	x.testWithInversion(
		`function foo() {
  var f = (x, ...y) => x + y[0];
  return f(1, 2);
}
function boo() {
  var f = (x, ...y) => x + y[0];
  return f(1, 2);
}
`,
		`function foo() {
  var f = (x, ...y) => x + y[0];
  return f(1, 2);
}
function boo() {
  var f$jscomp$1 = (x$jscomp$1, ...y$jscomp$1) => x$jscomp$1 + y$jscomp$1[0];
  return f$jscomp$1(1, 2);
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testDefaultParameterWithContextWithInversion1()
func TestMakeDeclaredNamesUnique_testDefaultParameterWithContextWithInversion1(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		`function foo(x = 1) {
  return x;
}
function boo(x = 1) {
  return x;
}
`,
		`function foo(x = 1) {
  return x;
}
function boo(x$jscomp$1 = 1) {
  return x$jscomp$1;
}
`)

	x.testSameWithInversion(
		`function foo(x = 1, y = x) {
  return x + y;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testDefaultParameterWithContextWithInversion2()
func TestMakeDeclaredNamesUnique_testDefaultParameterWithContextWithInversion2(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true

	// Parameter default values don't see the scope of the body
	// Methods or functions defined "inside" parameter default values don't see the local variables
	// of the body.
	x.testWithInversion(
		`let x = 'outer';
function foo(bar = baz => x) {
  let x = 'inner';
  console.log(bar());
}
`,
		`let x = 'outer';
function foo(bar = baz => x) {
  let x$jscomp$1 = 'inner';
  console.log(bar());
}
`)

	x.testWithInversion(
		`const x = 'outer';
function foo(a = x) {
  const x = 'inner';
  return a;
}
`,
		`const x = 'outer';
function foo(a = x) {
  const x$jscomp$1 = 'inner';
  return a;
}
`)

	x.testWithInversion(
		`const x = 'outerouter';
{
  const x = 'outer';
  function foo(a = x) {
    return a;
  }
foo();
}
`,
		`const x = 'outerouter';
{
  const x$jscomp$1 = 'outer';
  function foo(a = x$jscomp$1) {
    return a;
  }
foo();
}
`)

	x.testSameWithInversion(
		`function foo(x, y = x) {
  return x + y;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testObjectLiteralsWithContextWithInversion()
func TestMakeDeclaredNamesUnique_testObjectLiteralsWithContextWithInversion(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		`function foo({x:y}) {
  return y;
}
function boo({x:y}) {
  return y;
}
`,
		`function foo({x:y}) {
  return y;
}
function boo({x:y$jscomp$1}) {
  return y$jscomp$1
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testExportedOrImportedNamesAreUntouched()
func TestMakeDeclaredNamesUnique_testExportedOrImportedNamesAreUntouched(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	// The eventual desired behavior is that none of the 'a's in the following test cases
	// are renamed to a$jscomp$1. Rewrite this test after that behavior is implemented.
	x.useDefaultRenamer = true
	x.test(
		x.srcs("var a;", "let a; export {a as a};"),
		x.expected("var a;", "let a$jscomp$1; export {a$jscomp$1 as a};"))

	x.test(
		x.srcs("var a;", "import {a as a} from './foo.js'; let b = a;"),
		x.expected("var a;", "import {a as a$jscomp$1} from './foo.js'; let b = a$jscomp$1;"))
}

// Java: MakeDeclaredNamesUniqueTest.testTwoMethodsInTheSameFileWithSameLocalNames()
func TestMakeDeclaredNamesUnique_testTwoMethodsInTheSameFileWithSameLocalNames(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	// Verify same local names in 2 different files get different new names.
	// The ContextualRenamer renames an "oldName" by adding "$jscomp$\n" + "id" as a suffix string.
	// So, when another declaration containing "$jscomp$id" exists at any other source location in
	// the entire JS program (due to a prior renaming), the ContextualRenamer should not generate
	// that same name when renaming this declaration.
	x.test(
		"function foo() {var a; a;} function bar() {let a; let a$jscomp$1; a + a$jscomp$1;}",
		`function foo() {var a; a;}
function bar() {
  let a$jscomp$1; let a$jscomp$1$jscomp$1;
  a$jscomp$1 + a$jscomp$1$jscomp$1;
}
`)

	x.test(
		"function bar() {let a; let a$jscomp$1; a + a$jscomp$1;} function foo() {var a; a;}",
		`function bar() {
  let a; let a$jscomp$1;
  a + a$jscomp$1;
}
function foo() {var a$jscomp$2; a$jscomp$2;}
`)

	// tests when name with $jscomp$1 suffix comes first
	x.test(
		"function bar() {let a$jscomp$1; let a; a + a$jscomp$1;} function foo() {var a; a;}",
		`function bar() {
  let a$jscomp$1; let a;
  a + a$jscomp$1;
}
function foo() {var a$jscomp$2; a$jscomp$2;}
`)

	x.test(
		`function bar() {
  let a; let a$jscomp$1;
  a + a$jscomp$1;
}
function foo() {
// tests when a$jscomp$2 declared later in the same scope
  var a; a; var a$jscomp$2; a$jscomp$2;
}
`,
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
  var a$jscomp$2; a$jscomp$2;
  var a$jscomp$2$jscomp$1; a$jscomp$2$jscomp$1;
}
`)

	x.test(
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
// tests when a$jscomp$2 declared first in the same scope
  var a$jscomp$2; a$jscomp$2; var a; a;
}
`,
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
  var a$jscomp$2; a$jscomp$2; var a$jscomp$3; a$jscomp$3;
}
`)

	x.test(
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
// tests when a$jscomp$2 is declared in another scope
  var a; a;
}
function baz() {
  var a$jscomp$2; a$jscomp$2;
}
`,
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
  var a$jscomp$2; a$jscomp$2;
}
function baz() {
  var a$jscomp$2$jscomp$1; a$jscomp$2$jscomp$1;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testTwoFilesWithSameLocalNames()
func TestMakeDeclaredNamesUnique_testTwoFilesWithSameLocalNames(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	// Verify same local names in 2 different files get different new names
	x.test(
		x.srcs(
			"function foo() {var a; a;}",
			"function bar() {let a; let a$jscomp$1; a + a$jscomp$1;}"),
		x.expected(
			"function foo() {var a; a;}",
			`function bar() {let a$jscomp$1; let a$jscomp$1$jscomp$1; a$jscomp$1 +
 a$jscomp$1$jscomp$1;}
`))
}

// Java: MakeDeclaredNamesUniqueTest.testImportStarWithInversion()
func TestMakeDeclaredNamesUnique_testImportStarWithInversion(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.testWithInversion(
		[]string{
			"let a = 5;", //
			"import * as a          from './a.js'; const TAU = 2 * a.PI;",
		},
		[]string{
			"let a = 5;", //
			"import * as a$jscomp$1 from './a.js'; const TAU = 2 * a$jscomp$1.PI",
		})
}

// Java: MakeDeclaredNamesUniqueTest.assertOnChange_throwsException()
func TestMakeDeclaredNamesUnique_assertOnChange_throwsException(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.assertOnChange = true

	e := jsaThrows[jutil.IRuntimeException](t, func() { x.testNoWarning("var a; function foo() { var a = 1; } ") })
	if !strings.Contains(e.GetMessage(), "NAME a") {
		t.Errorf("message: %q", e.GetMessage())
	}
}

// Java: MakeDeclaredNamesUniqueTest.assertOnChange_noExceptionIfNothingChanges()
func TestMakeDeclaredNamesUnique_assertOnChange_noExceptionIfNothingChanges(t *testing.T) {
	x := newMakeDeclaredNamesUniqueTest(t)
	x.useDefaultRenamer = true
	x.assertOnChange = true

	x.testSame("const x = 1; function foo() { const y = 2; }")
}
