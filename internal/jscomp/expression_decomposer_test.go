// Port of com/google/javascript/jscomp/ExpressionDecomposerTest.java @ closure-compiler v20260415

package jscomp

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Unit tests for ExpressionDecomposer. Note: functions "foo" and "goo" are
// external functions in the helper.
//
// Every helper parses with Compiler.parseTestCode and
// Normalize.NormalizeStatements, finds nodes and compares results through
// Compiler.toSource (CodePrinter), type-checks with TypeCheck and validates
// source positions with SourceInfoCheck: all of them are above the core, so
// the helpers skip the test. The cases themselves are ported in full; the
// helpers' Java bodies are summarized on each.
const expressionDecomposerTest_skip = "needs Compiler.parseTestCode, Normalize.NormalizeStatements, " +
	"CodePrinter (Compiler.toSource), TypeCheck, SourceInfoCheck"

// expressionDecomposerTest is ExpressionDecomposerTest's fixture.
type expressionDecomposerTest struct {
	t              *testing.T
	knownConstants *jutil.LinkedHashSet[string]

	// The language out to set in the compiler options. If null, use the default.
	languageOut *CompilerOptions_LanguageMode

	// Whether we should run type checking and test the type information in the output expression
	shouldTestTypes bool
}

// Java: ExpressionDecomposerTest.setUp()
func newExpressionDecomposerTest(t *testing.T) *expressionDecomposerTest {
	x := &expressionDecomposerTest{t: t, knownConstants: jutil.NewLinkedHashSet[string]()}
	x.knownConstants.Clear()
	x.shouldTestTypes = true
	x.languageOut = nil
	return x
}

func (x *expressionDecomposerTest) setLanguageOut(languageOut CompilerOptions_LanguageMode) {
	x.languageOut = &languageOut
}

// expressionDecomposerTest_nodeFinder is the tests'
// Function<AbstractCompiler, Function<Node, Node>>.
type expressionDecomposerTest_nodeFinder func(compiler AbstractCompiler) func(root *rhino.Node) *rhino.Node

// helperCanExposeExpression parses code (and the externs "function goo()
// {} function foo() {}"), finds the expression and checks
// decomposer.canExposeExpression(expression) == expectedResult.
//
// Java: ExpressionDecomposerTest.helperCanExposeExpression(DecompositionType, String, Function<AbstractCompiler, Function<Node, Node>>)
func (x *expressionDecomposerTest) helperCanExposeExpression(
	expectedResult ExpressionDecomposer_DecompositionType, code string, nodeFinderFn expressionDecomposerTest_nodeFinder) {
	x.t.Skip(expressionDecomposerTest_skip)
}

// expressionDecomposerTest_nodeToSource is NodeToSource.toString(): the
// source code a node represents where possible, or some explanatory text
// when not possible.
//
// Java: ExpressionDecomposerTest.nodeToSource(AbstractCompiler, Node)
func expressionDecomposerTest_nodeToSource(compiler AbstractCompiler, node *rhino.Node) string {
	if node.IsTemplateLitString() {
		// A string part of a template literal cannot be printed as code on its own.
		return "[Template literal string: '" + node.GetRawString() + "']"
	} else if node.IsTemplateLitSub() {
		// The template literal substitution node cannot itself be turned into source code,
		// but we can do that for the expression inside of it.
		return "[Template literal substitution: '" + compiler.ToSourceOfNode(node.GetOnlyChild()) + "']"
	} else {
		return compiler.ToSourceOfNode(node)
	}
}

// Java: ExpressionDecomposerTest.exprMatchesStr(String)
func exprMatchesStr(exprString string) expressionDecomposerTest_nodeFinder {
	return func(compiler AbstractCompiler) func(root *rhino.Node) *rhino.Node {
		return func(root *rhino.Node) *rhino.Node {
			// When matching trim off the trailing newline added by the compiler's pretty-print
			// option.
			isAMatch := func(node *rhino.Node) bool {
				return exprString == jutil.StringTrim(expressionDecomposerTest_nodeToSource(compiler, node))
			}
			containsAMatch := func(node *rhino.Node) bool {
				return strings.Contains(expressionDecomposerTest_nodeToSource(compiler, node), exprString)
			}
			matchingNode := NodeUtil_FindPreorder(root, isAMatch, containsAMatch)
			if matchingNode == nil {
				panic(jutil.NewAssertionError(
					"Expected node `" + exprString + "` was not found in `" + compiler.ToSourceOfNode(root) + "`"))
			}
			return matchingNode
		}
	}
}

// helperExposeExpressionThenTypeCheck parses code, exposes the found
// expression (temp prefix "temp", result prefix "result") and type-checks
// the result.
//
// Java: ExpressionDecomposerTest.helperExposeExpressionThenTypeCheck(String, Function<Node, Node>)
func (x *expressionDecomposerTest) helperExposeExpressionThenTypeCheck(
	code string, nodeFinder func(*rhino.Node) *rhino.Node) *rhino.Node {
	x.t.Skip(expressionDecomposerTest_skip)
	return nil
}

// helperExposeExpression parses code (type-checking it when
// shouldTestTypes), checks the found expression is DECOMPOSABLE, calls
// maybeExposeExpression, validates source info and compares the result
// with expectedResult as source; when shouldTestTypes, the types must match
// those of decompose-then-type-check.
//
// Java: ExpressionDecomposerTest.helperExposeExpression(String, Function<AbstractCompiler, Function<Node, Node>>, String)
func (x *expressionDecomposerTest) helperExposeExpression(
	code string, compilerToNodeFinder expressionDecomposerTest_nodeFinder, expectedResult string) {
	x.t.Skip(expressionDecomposerTest_skip)
}

// createDecomposer is compiler.createExpressionDecomposer(
// compiler.getUniqueNameIdSupplier(), ImmutableSet.copyOf(knownConstants),
// newScope()).
//
// Java: ExpressionDecomposerTest.createDecomposer(Compiler)
func (x *expressionDecomposerTest) createDecomposer(compiler AbstractCompiler) *ExpressionDecomposer {
	return compiler.CreateExpressionDecomposer(
		compiler.GetUniqueNameIdSupplier(), jutil.NewLinkedHashSet(x.knownConstants.ToSlice()...), x.newScope())
}

// helperMoveExpression parses code (type-checking it when shouldTestTypes),
// calls moveExpression on the found expression, validates source info and
// compares the result with expectedResult as source; when shouldTestTypes,
// the types must match those of move-then-type-check.
//
// Java: ExpressionDecomposerTest.helperMoveExpression(String, Function<AbstractCompiler, Function<Node, Node>>, String)
func (x *expressionDecomposerTest) helperMoveExpression(
	code string, compilerToNodeFinder expressionDecomposerTest_nodeFinder, expectedResult string) {
	x.t.Skip(expressionDecomposerTest_skip)
}

// Java: ExpressionDecomposerTest.findClass(Node)
func expressionDecomposerTest_findClass(n *rhino.Node) *rhino.Node {
	if n.IsClass() {
		return n
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		maybeClass := expressionDecomposerTest_findClass(child)
		if maybeClass != nil {
			return maybeClass
		}
	}
	return nil
}

// expressionDecomposerTest_assertToken is assertNode(n).hasToken(token).
func expressionDecomposerTest_assertToken(t *testing.T, n *rhino.Node, token rhino.Token) {
	t.Helper()
	if n.GetToken() != token {
		t.Fatalf("%v: expected token %s", n, token)
	}
}

// Java: ExpressionDecomposerTest.newScope()
func (x *expressionDecomposerTest) newScope() *Scope {
	return Scope_CreateGlobalScope(rhino.NewNode(rhino.Token_ROOT))
}

// Java: ExpressionDecomposerTest.testWindowLocationAssign()
func TestExpressionDecomposer_testWindowLocationAssign(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// avoid decomposing `window.location.assign` when the output code could
	// end up running on IE11. See more explanation in ExpressionDecomposer.java.
	x.setLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "window.location.assign(foo())", exprMatchesStr("foo()"))
	x.helperMoveExpression(
		"window.location.assign(foo())",
		exprMatchesStr("foo()"),
		`var result$jscomp$0 = foo();
window.location.assign(result$jscomp$0)
`)

	// confirm that the default behavior does not treat window.location.assign
	// specially
	x.languageOut = nil
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "window.location.assign(foo())", exprMatchesStr("foo()"))
	x.helperExposeExpression(
		"window.location.assign(foo())",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$1 = window.location;
var temp_const$jscomp$0 = temp_const$jscomp$1.assign;
temp_const$jscomp$0.call(temp_const$jscomp$1, foo());
`)
}

// Java: ExpressionDecomposerTest.testObjectDestructuring_withComputedKey_doesNotCrash()
func TestExpressionDecomposer_testObjectDestructuring_withComputedKey_doesNotCrash(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// computed prop is found to be decomposable
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "var a; ({ [foo()]: a} = obj);", exprMatchesStr("foo()"))

	// TODO(b/339040894): Fix this crash.
	ex := jsaThrows[*jutil.IllegalStateException](t, func() {
		x.helperExposeExpression(
			"var a; ({ [foo()]: a} = obj);", //
			exprMatchesStr("foo()"),
			`var a;
var temp_const$jscomp$0 = obj;
var temp_const$jscomp$1 = foo();
({ [temp_const$jscomp$1]: a} = temp_const$jscomp$0);
`)
	})
	if !strings.Contains(ex.GetMessage(), "exposeExpression exposed nothing") {
		t.Errorf("message: %q", ex.GetMessage())
	}
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression1()
func TestExpressionDecomposer_testCannotExpose_expression1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// Can't move or decompose some classes of expressions.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(foo());", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(x = goo()&&foo()){}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(x += goo()&&foo()){}", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "do{}while(foo());", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "for(;foo(););", exprMatchesStr("foo()"))
	// This case could be supported for loops without conditional continues
	// by moving the increment into the loop body.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "for(;;foo());", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "for(foo();;);", exprMatchesStr("foo()"))

	// This is potentially doable but a bit too complex currently.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "switch(1){case foo():;}", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression2()
func TestExpressionDecomposer_testCanExposeExpression2(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "var x = foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const x = foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "let x = foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "if(foo()){}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "switch(foo()){}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "switch(foo()){}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "function f(){ return foo();}", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo() && 1", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo() || 1", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo() ? 0 : 1", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "(function(a){b = a})(foo())", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "function f(){ throw foo();}", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.nullishCoalesceMovable()
func TestExpressionDecomposer_nullishCoalesceMovable(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo() ?? 1", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.nullishCoalesceDecomposable()
func TestExpressionDecomposer_nullishCoalesceDecomposable(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "var x = null ?? foo()", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.nullishCoalesceUnDecomposable()
func TestExpressionDecomposer_nullishCoalesceUnDecomposable(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(x = goo()??foo()){}", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.assignCoalesceMovable()
func TestExpressionDecomposer_assignCoalesceMovable(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x ??= goo() + foo()", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.optChainMovable()
func TestExpressionDecomposer_optChainMovable(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo()?.x", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo()?.[x]", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo()?.()", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.optChainDecomposable()
func TestExpressionDecomposer_optChainDecomposable(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x?.[foo()]", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x?.(foo())", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.optChainAllowMethodCallDecomposable()
func TestExpressionDecomposer_optChainAllowMethodCallDecomposable(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x?.y(foo())", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.optChainUnDecomposable()
func TestExpressionDecomposer_optChainUnDecomposable(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(x = y?.[foo()]){}", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression3()
func TestExpressionDecomposer_testCanExposeExpression3(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x = 0 && foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x = 1 || foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "var x = 1 ? foo() : 0", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "const x = 1 ? foo() : 0", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "let x = 1 ? foo() : 0", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "goo() && foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x = goo() && foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x += goo() && foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "var x = goo() && foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "const x = goo() && foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "let x = goo() && foo()", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "if(goo() && foo()){}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "switch(goo() && foo()){}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "switch(goo() && foo()){}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "switch(x = goo() && foo()){}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"function f(){ return goo() && foo();}",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression_compoundDeclaration_inForInitializer_firstElement()
func TestExpressionDecomposer_testCanExposeExpression_compoundDeclaration_inForInitializer_firstElement(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// VAR will already be hoisted by `Normalize`.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "for (var x = foo(), y = 5;;) {}", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "for (let x = foo(), y = 5;;) {}", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "for (const x = foo(), y = 5;;) {}", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression_compoundDeclaration_inForInitializer_nthElement()
func TestExpressionDecomposer_testCanExposeExpression_compoundDeclaration_inForInitializer_nthElement(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// TODO(b/121157467) FOR introduces complex scoping that isn't currently `Normalize`d.
	// Since in some cases we'd effectively end up having to `Normalize` these, decomposition just
	// bails for now.

	// VAR will already be hoisted by `Normalize`.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "for (var x = 8, y = foo();;) {}", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"for (let x = 8, y = foo();;) {}",
		exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"for (const x = 8, y = foo();;) {}",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression4()
func TestExpressionDecomposer_testCannotExpose_expression4(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// 'this' must be preserved in call.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "if (goo.a(1, foo()));", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression5()
func TestExpressionDecomposer_testCannotExpose_expression5(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// 'this' must be preserved in call.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "if (goo['a'](foo()));", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression6()
func TestExpressionDecomposer_testCannotExpose_expression6(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "z:if (goo.a(1, foo()));", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression7()
func TestExpressionDecomposer_testCanExposeExpression7(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// Verify calls to function expressions are movable.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		`(function(map){descriptions_=map})(
  function(){
    var ret={};
    ret[INIT]='a';
    ret[MIGRATION_BANNER_DISMISS]='b';
    return ret
  }());
`,
		func(compiler AbstractCompiler) func(*rhino.Node) *rhino.Node {
			return func(rootNode *rhino.Node) *rhino.Node {
				// Dig out the inner IIFE call and return it as the expression we want to ensure is
				// movable.
				expressionDecomposerTest_assertToken(t, rootNode, rhino.Token_SCRIPT)
				exprResult := rootNode.GetOnlyChild()
				expressionDecomposerTest_assertToken(t, exprResult, rhino.Token_EXPR_RESULT)
				outerCall := exprResult.GetOnlyChild()
				expressionDecomposerTest_assertToken(t, outerCall, rhino.Token_CALL)
				innerIifeCall := outerCall.GetSecondChild()
				expressionDecomposerTest_assertToken(t, innerIifeCall, rhino.Token_CALL)
				return innerIifeCall
			}
		})
}

// Java: ExpressionDecomposerTest.testCanExposeExpression8()
func TestExpressionDecomposer_testCanExposeExpression8(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// Can it be decompose?
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		`HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var myUrl = new goog.Uri(
      getDomServices_(self).getDomHelper().getWindow().location.href);
};
`,
		exprMatchesStr("getDomServices_(self)"))

	// Verify it is properly expose the target expression.
	x.helperExposeExpression(
		`HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var myUrl =
      new goog.Uri(getDomServices_(self).getDomHelper().getWindow().location.href);
};
`,
		exprMatchesStr("getDomServices_(self)"),
		`HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var temp_const$jscomp$0 = goog.Uri;
  var myUrl = new temp_const$jscomp$0(
      getDomServices_(self).getDomHelper().getWindow().location.href);
}
`)

	// Verify the results can be properly moved.
	x.helperMoveExpression(
		`HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var temp_const$jscomp$0 = goog.Uri;
  var myUrl = new temp_const$jscomp$0(
      getDomServices_(self).getDomHelper().getWindow().location.href);
}
`,
		exprMatchesStr("getDomServices_(self)"),
		`HangoutStarter.prototype.launchHangout = function() {
  var self=a.b;
  var temp_const$jscomp$0=goog.Uri;
  var result$jscomp$0=getDomServices_(self);
  var myUrl=new temp_const$jscomp$0(
      result$jscomp$0.getDomHelper().getWindow().location.href);
}
`)
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression9()
func TestExpressionDecomposer_testCannotExpose_expression9(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"function *f() { for (let x of yield y) {} }",
		exprMatchesStr("yield y"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_forAwaitOf()
func TestExpressionDecomposer_testCannotExpose_forAwaitOf(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"async function *f() { for await (let x of yield y) {} }",
		exprMatchesStr("yield y"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression10()
func TestExpressionDecomposer_testCanExposeExpression10(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"function *f() { for (let x in yield y) {} }",
		exprMatchesStr("yield y"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression11()
func TestExpressionDecomposer_testCannotExpose_expression11(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// expressions in parameter lists
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "function f(x = foo()) {}", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "function f({[foo()]: x}) {}", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "(function (x = foo()) {})()", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"(function ({[foo()]: x}) {})()",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExpose_aCall_withSpreadSibling()
func TestExpressionDecomposer_testCanExpose_aCall_withSpreadSibling(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...x, y());", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(y(), ...x);", exprMatchesStr("y()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "new D(...x, y());", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "new D(y(), ...x);", exprMatchesStr("y()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "[...x, y()];", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "({...x, z: y()});", exprMatchesStr("y()"))

	// Array- and object-literal instantiations cannot be side-effected.
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "[y(), ...x];", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "({z: y(), ...x});", exprMatchesStr("y()"))

	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...y());", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...y(), x);", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(x, ...y());", exprMatchesStr("y()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...x, x, y());", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...x, ...x, y());", exprMatchesStr("y()"))
}

// Java: ExpressionDecomposerTest.testCanExpose_anExpression_withSpreadRelative_ifInDifferentFunction()
func TestExpressionDecomposer_testCanExpose_anExpression_withSpreadRelative_ifInDifferentFunction(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// TODO(b/121004488): There are potential decompositions that weren't implemented.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(function() { [...x]; }, y());", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(function() { ({...x}); }, y());", exprMatchesStr("y()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(y(), () => [...x]);", exprMatchesStr("y()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(y(), () => ({...x}));", exprMatchesStr("y()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "[() => f(...x), y()];", exprMatchesStr("y()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		`[
   class {
     f(x) { return [...x]; }
   },
  y()
];
`,
		exprMatchesStr("y()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression12()
func TestExpressionDecomposer_testCanExposeExpression12(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// Test destructuring rhs is evaluated before the lhs
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const {a, b = goo()} = foo();", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const [a, b = goo()] = foo();", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "({a, b = goo()} = foo());", exprMatchesStr("foo()"))

	// Default value expressions are conditional, which would make the expressions complex.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		// default value inside array pattern
		"[{ [foo()]: a } = goo()] = arr;",
		exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, //
		// computed property inside object pattern; decomposed
		"({ [foo()]: a = goo()} = arr);",
		exprMatchesStr("foo()"))

	// default value expressions are conditional, which would make the expressions complex
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, //
		// default value inside object pattern
		"({ [foo()]: a = goo()} = arr);",
		exprMatchesStr("goo()"))

	x.helperExposeExpression(
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), [getObj().propName] = CN();
  }
  function CN() {
    return [1];
  }
});
`,
		exprMatchesStr("CN()"),
		`var Di = I(() => {
  function zv() {
    var temp_const$jscomp$0 = JSCOMPILER_PRESERVE(e);
// TODO(b/339701959): We decided to back off here because of b/338660589. But we should
// optimize this to:
//   var temp_const$jscomp$1 = CN();
//   var temp_const$jscomp$2 = getObj();
//   temp_const$jscomp$0, [temp_const$jscomp$2.propName] = temp_const$jscomp$1;
    temp_const$jscomp$0, [getObj().propName] = CN();
  }
  function CN() {
    return [1];
  }
});
`)

	x.helperExposeExpression(
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), ({x: getObj().propName} = CN());
  }
  function CN() {
    return {x: 1};
  }
});
`,
		exprMatchesStr("CN()"),
		`var Di = I(() => {
  function zv() {
    var temp_const$jscomp$0 = JSCOMPILER_PRESERVE(e);
// TODO(b/339701959): We decided to back off here because of b/338246627. But we should
// optimize this to:
//   var temp_const$jscomp$1 = CN();
//   var temp_const$jscomp$2 = getObj();
//   temp_const$jscomp$0, {x: temp_const$jscomp$2.propName} = temp_const$jscomp$1;
    temp_const$jscomp$0, {x:getObj().propName} = CN();
  }
  function CN() {
    return {x: 1};
  }
});
`)

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), [f] = CN();
  }
  function CN() {
    return [1];
  }
});
`,
		exprMatchesStr("CN()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), [f = foo()] = CN();
  }
  function CN() {
    return [1];
  }
});
`,
		exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), ({f: g} = CN());
  }
  function CN() {
    return {f: 1};
  }
});
`,
		exprMatchesStr("CN()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), ({f: g = goo()} = CN());
  }
  function CN() {
    return {f: 1};
  }
});
`,
		exprMatchesStr("goo()"))
}

// Java: ExpressionDecomposerTest.testObjectDestructuring_withDefaultValue_generatesValidAST()
func TestExpressionDecomposer_testObjectDestructuring_withDefaultValue_generatesValidAST(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"var d; ({c: d = 4} = condition ? y() :  {c: 1});",
		exprMatchesStr("y()"),
		`var d;
var temp$jscomp$0;
if (condition) {
  temp$jscomp$0 = y();
} else {
  temp$jscomp$0 = {c: 1};
}
({c: d = 4} = temp$jscomp$0);
`)
}

// Java: ExpressionDecomposerTest.testObjectDestructuring_withDefaultValue_withComputedKey()
func TestExpressionDecomposer_testObjectDestructuring_withDefaultValue_withComputedKey(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// default value expressions are conditional, which would make the expressions complex
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"var a; ({ [foo()]: a = bar()} = baz());",
		exprMatchesStr("bar()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		"var a; ({ [foo()]: a = bar()} = baz());",
		exprMatchesStr("baz()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"var a; ({ [foo()]: a = bar()} = baz());",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testArrayDestructuring_withDefaultValue_generatesValidAST()
func TestExpressionDecomposer_testArrayDestructuring_withDefaultValue_generatesValidAST(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"var [c = 4] = condition ? y() :  [c = 2];",
		exprMatchesStr("y()"),
		`var c;
var temp$jscomp$0;
if (condition) {
  temp$jscomp$0 = y();
} else {
  temp$jscomp$0 = [c = 2];
}
[c = 4] = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testCanExposeExpressionInTemplateLiteralSubstitution()
func TestExpressionDecomposer_testCanExposeExpressionInTemplateLiteralSubstitution(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const result = `${foo()}`;", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"const obj = {f(x) {}}; obj.f(`${foo()}`);",
		exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const result = `${foo()} ${goo()}`;", exprMatchesStr("foo()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"const result = `${foo()} ${goo()}`;",
		exprMatchesStr("goo()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_defaultValueInParamList()
func TestExpressionDecomposer_testCannotExpose_defaultValueInParamList(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "function fn(a = g()) {}", exprMatchesStr("g()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_defaultValueInDestructuring()
func TestExpressionDecomposer_testCannotExpose_defaultValueInDestructuring(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "let {x = fn()} = y;", exprMatchesStr("fn()"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "let [x = fn()] = y;", exprMatchesStr("fn()"))
}

// Java: ExpressionDecomposerTest.testMoveExpression1()
func TestExpressionDecomposer_testMoveExpression1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// There isn't a reason to do this, but it works.
	x.helperMoveExpression(
		"foo()", exprMatchesStr("foo()"), "var result$jscomp$0 = foo(); result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression2()
func TestExpressionDecomposer_testMoveExpression2(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"x = foo()", exprMatchesStr("foo()"), "var result$jscomp$0 = foo(); x = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression3()
func TestExpressionDecomposer_testMoveExpression3(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"var x = foo()",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); var x = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression4()
func TestExpressionDecomposer_testMoveExpression4(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"const x = foo()",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); const x = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression5()
func TestExpressionDecomposer_testMoveExpression5(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"let x = foo()",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); let x = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression6()
func TestExpressionDecomposer_testMoveExpression6(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"if(foo()){}",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); if (result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testMoveExpression7()
func TestExpressionDecomposer_testMoveExpression7(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"switch(foo()){}",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); switch(result$jscomp$0){}")
}

// Java: ExpressionDecomposerTest.testMoveExpression8()
func TestExpressionDecomposer_testMoveExpression8(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"switch(1 + foo()){}",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); switch(1 + result$jscomp$0){}")
}

// Java: ExpressionDecomposerTest.testMoveExpression9()
func TestExpressionDecomposer_testMoveExpression9(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"function f(){ return foo();}",
		exprMatchesStr("foo()"),
		"function f(){ var result$jscomp$0 = foo(); return result$jscomp$0;}")
}

// Java: ExpressionDecomposerTest.testMoveExpression10()
func TestExpressionDecomposer_testMoveExpression10(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"x = foo() && 1",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); x = result$jscomp$0 && 1")
}

// Java: ExpressionDecomposerTest.testMoveExpression11()
func TestExpressionDecomposer_testMoveExpression11(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"x = foo() || 1",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); x = result$jscomp$0 || 1")
}

// Java: ExpressionDecomposerTest.testMoveExpression12()
func TestExpressionDecomposer_testMoveExpression12(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"x = foo() ? 0 : 1",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); x = result$jscomp$0 ? 0 : 1")
}

// Java: ExpressionDecomposerTest.testMoveExpressionNullishCoalesce()
func TestExpressionDecomposer_testMoveExpressionNullishCoalesce(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"x = foo() ?? 0",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); x = result$jscomp$0 ?? 0")
}

// Java: ExpressionDecomposerTest.testMoveExpressionOptionalChain()
func TestExpressionDecomposer_testMoveExpressionOptionalChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"foo()?.x", exprMatchesStr("foo()"), "var result$jscomp$0 = foo(); result$jscomp$0?.x")
}

// Java: ExpressionDecomposerTest.testMoveExpression13()
func TestExpressionDecomposer_testMoveExpression13(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"const {a, b} = foo();",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); const {a, b} = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression14()
func TestExpressionDecomposer_testMoveExpression14(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"({a, b} = foo());",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); ({a, b} = result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testExposeExpression1()
func TestExpressionDecomposer_testExposeExpression1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"x = 0 && foo()",
		exprMatchesStr("foo()"),
		"var temp$jscomp$0; if (temp$jscomp$0 = 0) temp$jscomp$0 = foo(); x = temp$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testExposeExpression2()
func TestExpressionDecomposer_testExposeExpression2(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"x = 1 || foo()",
		exprMatchesStr("foo()"),
		"var temp$jscomp$0; if (temp$jscomp$0 = 1); else temp$jscomp$0=foo(); x = temp$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testExposeExpression3()
func TestExpressionDecomposer_testExposeExpression3(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"var x = 1 ? foo() : 0",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (1) temp$jscomp$0 = foo(); else temp$jscomp$0 = 0;
var x = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression4()
func TestExpressionDecomposer_testExposeExpression4(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"const x = 1 ? foo() : 0",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (1) temp$jscomp$0 = foo(); else temp$jscomp$0 = 0;
const x = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression5()
func TestExpressionDecomposer_testExposeExpression5(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"let x = 1 ? foo() : 0",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (1) temp$jscomp$0 = foo(); else temp$jscomp$0 = 0;
let x = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression6()
func TestExpressionDecomposer_testExposeExpression6(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression("goo() && foo()", exprMatchesStr("foo()"), "if (goo()) foo();")
}

// Java: ExpressionDecomposerTest.exposeExpressionNullishCoalesceNoResult()
func TestExpressionDecomposer_exposeExpressionNullishCoalesceNoResult(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"goo() ?? foo()",
		exprMatchesStr("foo()"),
		`var temp$jscomp$1;
if((temp$jscomp$1 = goo()) != null) temp$jscomp$1;
else foo()
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetElem()
func TestExpressionDecomposer_exposeExpressionOptionalGetElem(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  temp$jscomp$1 = temp_const$jscomp$0[foo()];
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptChainCallChain()
func TestExpressionDecomposer_exposeExpressionOptChainCallChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.(a).y.z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0(a).y.z;
  temp$jscomp$1 = temp_const$jscomp$2[foo()];
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptChainCallChainNoResult()
func TestExpressionDecomposer_exposeExpressionOptChainCallChainNoResult(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"x?.(a)[y].z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
if ((temp_const$jscomp$0 = x) == null) {
  void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0(a)[y].z;
  temp_const$jscomp$2[foo()];
}
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetPropChain()
func TestExpressionDecomposer_exposeExpressionOptionalGetPropChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.y.z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0.y.z;
  temp$jscomp$1 = temp_const$jscomp$2[foo()];
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetPropChainNoResult()
func TestExpressionDecomposer_exposeExpressionOptionalGetPropChainNoResult(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"x?.y.z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
if ((temp_const$jscomp$0 = x) == null) {
  void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0.y.z;
  temp_const$jscomp$2[foo()];
}
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetElemChain()
func TestExpressionDecomposer_exposeExpressionOptionalGetElemChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.[y].z[foo()];",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0[y].z;
  temp$jscomp$1 = temp_const$jscomp$2[foo()];
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetElemChainNoResult()
func TestExpressionDecomposer_exposeExpressionOptionalGetElemChainNoResult(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"x?.[y].z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
if ((temp_const$jscomp$0 = x) == null) {
  void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0[y].z;
  temp_const$jscomp$2[foo()];
}
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetElemWithCall()
func TestExpressionDecomposer_exposeExpressionOptionalGetElemWithCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x.y?.[z](foo())",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x.y) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0;
  var temp_const$jscomp$2 = temp_const$jscomp$3[z];
  temp$jscomp$1 = temp_const$jscomp$2.call(temp_const$jscomp$3, foo());
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionGetElemWithOptChainCall()
func TestExpressionDecomposer_exposeExpressionGetElemWithOptChainCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x.y[z]?.(foo(), d)",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$1 = (temp_const$jscomp$0 = x.y)[z]) == null) {
  temp$jscomp$2 = void 0;
} else {
  temp$jscomp$2 = temp_const$jscomp$1.call(temp_const$jscomp$0, foo(), d);
}
a = temp$jscomp$2;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetPropWithCall()
func TestExpressionDecomposer_exposeExpressionOptionalGetPropWithCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x.y?.z(foo(1))",
		exprMatchesStr("foo(1)"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x.y) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0;
  var temp_const$jscomp$2 = temp_const$jscomp$3.z;
  temp$jscomp$1 = temp_const$jscomp$2.call(temp_const$jscomp$3, foo(1));
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetPropWithCallTwiceRewriteCall()
func TestExpressionDecomposer_exposeExpressionOptionalGetPropWithCallTwiceRewriteCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x.y?.z(foo(1))",
		exprMatchesStr("foo(1)"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x.y) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0;
  var temp_const$jscomp$2 = temp_const$jscomp$3.z;
  temp$jscomp$1 = temp_const$jscomp$2.call(temp_const$jscomp$3, foo(1));
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionGetPropWithOptChainCall()
func TestExpressionDecomposer_exposeExpressionGetPropWithOptChainCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x.y.z?.(foo())",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$1 = (temp_const$jscomp$0 = x.y).z) == null) {
  temp$jscomp$2 = void 0;
} else {
  temp$jscomp$2 = temp_const$jscomp$1.call(temp_const$jscomp$0, foo());
}
a = temp$jscomp$2;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionNewOptChainAfterRewriteCall()
func TestExpressionDecomposer_exposeExpressionNewOptChainAfterRewriteCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.y(foo())?.z.q",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$2 = void 0;
} else {
  var temp_const$jscomp$4 = temp_const$jscomp$0;
  var temp_const$jscomp$3 = temp_const$jscomp$4.y;
  temp$jscomp$2 =
      (temp_const$jscomp$1 =
          temp_const$jscomp$3.call(temp_const$jscomp$4, foo())) == null
              ? void 0 : temp_const$jscomp$1.z.q;
}
a = temp$jscomp$2;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionNewOptChainAfter()
func TestExpressionDecomposer_exposeExpressionNewOptChainAfter(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.y[foo()]?.z.q",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$2 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0.y;
  temp$jscomp$2 = (temp_const$jscomp$1 = temp_const$jscomp$3[foo()]) == null ? void 0 : temp_const$jscomp$1.z.q;
}
a = temp$jscomp$2;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionNotImmediatelyFollowedByNewChain()
func TestExpressionDecomposer_exposeExpressionNotImmediatelyFollowedByNewChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.y[foo()].z.q?.b.c",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$2 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0.y;
  temp$jscomp$2 = (temp_const$jscomp$1 = temp_const$jscomp$3[foo()].z.q) == null
      ? void 0 : temp_const$jscomp$1.b.c;
}
a = temp$jscomp$2;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionBreakingOutOfOptionalChain()
func TestExpressionDecomposer_exposeExpressionBreakingOutOfOptionalChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = (x?.y[foo()]).z.q",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0.y;
  temp$jscomp$1 = temp_const$jscomp$2[foo()];
}
a = temp$jscomp$1.z.q;
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionCallBreakingOutOfOptionalChain()
func TestExpressionDecomposer_exposeExpressionCallBreakingOutOfOptionalChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.shouldTestTypes = false
	// Performing a non-optional call on an optional chain is not good coding, because you could
	// end up trying to call `undefined` as a function, but it is allowed and must be supported.
	x.helperExposeExpression(
		"a = (x?.y.z)(foo())",
		exprMatchesStr("foo()"),
		"let temp_const$jscomp$0;\nlet temp_const$jscomp$1;\nvar temp_const$jscomp$3 =\n    (temp_const$jscomp$0 = x) == null\n        ? void 0 : (temp_const$jscomp$1 = temp_const$jscomp$0.y).z;\n// This double .call is unfortunate but not incorrect.\n// Maybe we should make the ExpressionDecomposer recognize `.call` and avoid creating\n// another one? We'd run the risk of breaking \"real\" methods called `.call`, which\n// are allowed.\nvar temp_const$jscomp$2 = temp_const$jscomp$3.call;\n// The temp_const$jscomp$1 argument gets type '?' when we decompose,\n// which is correct, but TypeInference on this expected code appears to give it\n// `undefined`. This is why we've set `shouldTestTypes = false;` above\na = temp_const$jscomp$2.call(temp_const$jscomp$3, temp_const$jscomp$1, foo());\n")
}

// Java: ExpressionDecomposerTest.exposeExpressionFreeCallBreakingOutOfOptionalChain()
func TestExpressionDecomposer_exposeExpressionFreeCallBreakingOutOfOptionalChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// Performing a non-optional call on an optional chain is not good coding, because you could
	// end up trying to call `undefined` as a function, but it is allowed and must be supported.
	x.helperExposeExpression(
		"a = (x?.y.z())(foo())",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = x?.y.z();
a = temp_const$jscomp$0(foo());
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionCallAtEndOfOptionalChain()
func TestExpressionDecomposer_exposeExpressionCallAtEndOfOptionalChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.y.z(foo())",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0.y;
  var temp_const$jscomp$2 = temp_const$jscomp$3.z;
  temp$jscomp$1 = temp_const$jscomp$2.call(temp_const$jscomp$3, foo());
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.testBug117935266_expose_call_target()
func TestExpressionDecomposer_testBug117935266_expose_call_target(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		`function first() {
  alert('first');
  return '';
}
// alert must be preserved before the first side-effect
alert(first().method(alert('second')).method(alert('third')));
`,
		exprMatchesStr("first()"),
		`function first() {
  alert('first');
      return '';
}
var temp_const$jscomp$0 = alert;
temp_const$jscomp$0(first().method(
    alert('second')).method(alert('third')));
`)
}

// Java: ExpressionDecomposerTest.testBug117935266_move_call_target()
func TestExpressionDecomposer_testBug117935266_move_call_target(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		`function first() {
  alert('first');
      return '';
}
var temp_const$jscomp$0 = alert;
temp_const$jscomp$0(first().toString(
    alert('second')).toString(alert('third')));
`,
		exprMatchesStr("first()"),
		`function first() {
  alert('first');
      return '';
}
var temp_const$jscomp$0 = alert;
var result$jscomp$0 = first();
temp_const$jscomp$0(result$jscomp$0.toString(
    alert('second')).toString(alert('third')));
`)
}

// Java: ExpressionDecomposerTest.testBug117935266_expose_call_parameters()
func TestExpressionDecomposer_testBug117935266_expose_call_parameters(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"alert(fn(first(), second(), third()));",
		exprMatchesStr("first()"),
		`var temp_const$jscomp$1 = alert;
var temp_const$jscomp$0 = fn;
temp_const$jscomp$1(temp_const$jscomp$0(first(), second(), third()));
`)

	x.helperExposeExpression(
		"alert(fn(first(), second(), third()));",
		exprMatchesStr("second()"),
		`var temp_const$jscomp$2 = alert;
var temp_const$jscomp$1 = fn;
var temp_const$jscomp$0 = first();
temp_const$jscomp$2(temp_const$jscomp$1(temp_const$jscomp$0, second(), third()));
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionAfterTwoOptionalChains()
func TestExpressionDecomposer_exposeExpressionAfterTwoOptionalChains(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = x?.y.z?.q(foo());",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$2 = void 0;
} else {
  var temp$jscomp$3;
  if ((temp_const$jscomp$1 = temp_const$jscomp$0.y.z) == null) {
    temp$jscomp$3 = void 0;
  } else {
    var temp_const$jscomp$5 = temp_const$jscomp$1;
    var temp_const$jscomp$4 = temp_const$jscomp$5.q;
    temp$jscomp$3 = temp_const$jscomp$4.call(temp_const$jscomp$5, foo());
  }
  temp$jscomp$2 = temp$jscomp$3;
}
a = temp$jscomp$2;
`)
}

// Java: ExpressionDecomposerTest.exposeAnOptionalChain()
func TestExpressionDecomposer_exposeAnOptionalChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"a = foo(arg1, opt?.chain())",
		exprMatchesStr("opt?.chain()"),
		`var temp_const$jscomp$1 = foo;
var temp_const$jscomp$0 = arg1;
a = temp_const$jscomp$1(temp_const$jscomp$0, opt?.chain());
`)
}

// Java: ExpressionDecomposerTest.exposePartOfAnOptionalChain()
func TestExpressionDecomposer_exposePartOfAnOptionalChain(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		// The non-optional part is fine to move.
		"nonOptional.part?.optional.chain?.continues()",
		exprMatchesStr("nonOptional.part"))

	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"nonOptional.part?.optional.chain?.continues()",
		// Cannot just move part of an optional chain.
		exprMatchesStr("nonOptional.part?.optional"))

	x.helperExposeExpression(
		"nonOptional.part?.optional.chain?.continues()",
		exprMatchesStr("nonOptional.part?.optional"),
		"let temp_const$jscomp$0;\nlet temp_const$jscomp$1;\nif ((temp_const$jscomp$0 = nonOptional.part) == null) {\n  void 0;\n} else {\n// `temp_const$jscomp$0.optional` is the expression we were trying to expose, and it\n// is now movable.\n  (temp_const$jscomp$1 = temp_const$jscomp$0.optional.chain) == null\n      ? void 0 : temp_const$jscomp$1.continues();\n}\n")
}

// Java: ExpressionDecomposerTest.canExposeMethodCallee()
func TestExpressionDecomposer_canExposeMethodCallee(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// TODO(b/161802885): This should probably return DECOMPOSABLE, because it is not safe to
	// replace `foo.bar` with a temporary variable containing its value.
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo.bar()", exprMatchesStr("foo.bar"))

	// TODO(b/161802885): This should probably return DECOMPOSABLE, because it is not safe to
	// replace `foo?.bar` with a temporary variable containing its value.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "(foo?.bar)()", exprMatchesStr("foo?.bar"))
}

// Java: ExpressionDecomposerTest.testExposeExpression7()
func TestExpressionDecomposer_testExposeExpression7(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"x = goo() && foo()",
		exprMatchesStr("foo()"),
		"var temp$jscomp$0; if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo(); x = temp$jscomp$0;")
}

// Java: ExpressionDecomposerTest.exposeExpressionNullishCoalesce()
func TestExpressionDecomposer_exposeExpressionNullishCoalesce(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"x = goo() ?? foo()",
		exprMatchesStr("foo()"),
		`var temp$jscomp$1;var temp$jscomp$0;
if((temp$jscomp$1 = goo()) != null) temp$jscomp$0 = temp$jscomp$1;
else temp$jscomp$0=foo(); x = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression8()
func TestExpressionDecomposer_testExposeExpression8(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"var x = 1 + (goo() && foo())",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
var x = 1 + temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression9()
func TestExpressionDecomposer_testExposeExpression9(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"const x = 1 + (goo() && foo())",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
const x = 1 + temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression10()
func TestExpressionDecomposer_testExposeExpression10(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"let x = 1 + (goo() && foo())",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
let x = 1 + temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression11()
func TestExpressionDecomposer_testExposeExpression11(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"if(goo() && foo());",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
if(temp$jscomp$0);
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression12()
func TestExpressionDecomposer_testExposeExpression12(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"switch(goo() && foo()){}",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
switch(temp$jscomp$0){}
`)
}

// Java: ExpressionDecomposerTest.exposeExpressionNullishCoalesceSwitch()
func TestExpressionDecomposer_exposeExpressionNullishCoalesceSwitch(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"switch(goo() ?? foo()){}",
		exprMatchesStr("foo()"),
		`var temp$jscomp$1;var temp$jscomp$0;
if((temp$jscomp$1 = goo()) != null) temp$jscomp$0 = temp$jscomp$1;
else temp$jscomp$0 = foo(); switch(temp$jscomp$0){}
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression13()
func TestExpressionDecomposer_testExposeExpression13(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"switch(1 + goo() + foo()){}",
		exprMatchesStr("foo()"),
		"var temp_const$jscomp$0 = 1 + goo(); switch(temp_const$jscomp$0 + foo()){}")
}

// Java: ExpressionDecomposerTest.testExposeExpression_inVanillaForInitializer_simpleExpression()
func TestExpressionDecomposer_testExposeExpression_inVanillaForInitializer_simpleExpression(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"for (x = goo() + foo();;) {}",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = goo();
for (x = temp_const$jscomp$0 + foo();;) {}
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression_inVanillaForInitializer_usingLabel()
func TestExpressionDecomposer_testExposeExpression_inVanillaForInitializer_usingLabel(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"LABEL: for (x = goo() + foo();;) {}",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = goo();
LABEL: for (x = temp_const$jscomp$0 + foo();;) {}
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression_inVanillaForInitializer_singleDeclaration_withLetOrConst()
func TestExpressionDecomposer_testExposeExpression_inVanillaForInitializer_singleDeclaration_withLetOrConst(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	for _, declarationKeyword := range []string{"let", "const"} {
		x.helperExposeExpression(
			"for ("+declarationKeyword+" x = goo() + foo();;) {}",
			exprMatchesStr("foo()"),
			strings.ReplaceAll(`var temp_const$jscomp$0 = goo();
for (DECLATION_KEYWORD x = temp_const$jscomp$0 + foo();;) {}
`, "DECLATION_KEYWORD", declarationKeyword))
	}
}

// Java: ExpressionDecomposerTest.testExposeExpression_inVanillaForInitializer_firstDeclaration_withLetOrConst()
func TestExpressionDecomposer_testExposeExpression_inVanillaForInitializer_firstDeclaration_withLetOrConst(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	for _, declarationKeyword := range []string{"let", "const"} {
		x.helperExposeExpression(
			"for ("+declarationKeyword+" x = goo() + foo(), y = 5;;) {}",
			exprMatchesStr("foo()"),
			strings.ReplaceAll(`var temp_const$jscomp$0 = goo();
for (DECLATION_KEYWORD x = temp_const$jscomp$0 + foo(), y = 5;;) {}
`, "DECLATION_KEYWORD", declarationKeyword))
	}
}

// Java: ExpressionDecomposerTest.testExposeExpression14()
func TestExpressionDecomposer_testExposeExpression14(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"function f(){ return goo() && foo();}",
		exprMatchesStr("foo()"),
		`function f() {
  var temp$jscomp$0; if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
  return temp$jscomp$0;
}
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression15()
func TestExpressionDecomposer_testExposeExpression15(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// TODO(johnlenz): We really want a constant marking pass.
	// The value "goo" should be constant, but it isn't known to be so.
	x.helperExposeExpression(
		"if (goo(1, goo(2), (1 ? foo() : 0)));",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$1 = goo;
var temp_const$jscomp$0 = goo(2);
var temp$jscomp$2;
if (1) temp$jscomp$2 = foo(); else temp$jscomp$2 = 0;
if (temp_const$jscomp$1(1, temp_const$jscomp$0, temp$jscomp$2));
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression16()
func TestExpressionDecomposer_testExposeExpression16(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"throw bar() && foo();",
		exprMatchesStr("foo()"),
		"var temp$jscomp$0; if (temp$jscomp$0 = bar()) temp$jscomp$0=foo(); throw temp$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testExposeExpression17()
func TestExpressionDecomposer_testExposeExpression17(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"x.foo(y())",
		exprMatchesStr("y()"),
		`var temp_const$jscomp$1 = x;
var temp_const$jscomp$0 = temp_const$jscomp$1.foo;
temp_const$jscomp$0.call(temp_const$jscomp$1, y());
`)
}

// Java: ExpressionDecomposerTest.testExposeFreeCall()
func TestExpressionDecomposer_testExposeFreeCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"(0,x.foo)(y())",
		exprMatchesStr("y()"),
		`var temp_const$jscomp$0 = x.foo;
temp_const$jscomp$0(y());
`)
}

// Java: ExpressionDecomposerTest.testExposeTemplateLiteralFreeCall()
func TestExpressionDecomposer_testExposeTemplateLiteralFreeCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"foo`${x()}${y()}`",
		exprMatchesStr("y()"),
		"var temp_const$jscomp$1 = foo;\nvar temp_const$jscomp$0 = x();\ntemp_const$jscomp$1`${temp_const$jscomp$0}${y()}`;\n")
}

// Java: ExpressionDecomposerTest.testCanExposeTaggedTemplateLiteralInterpolation()
func TestExpressionDecomposer_testCanExposeTaggedTemplateLiteralInterpolation(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x`${y()}`", exprMatchesStr("y()"))
	// TODO(b/251958225): Implement decomposition for this case.
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "x.foo`${y()}`", exprMatchesStr("y()"))
}

// Java: ExpressionDecomposerTest.testExposeExpression18()
func TestExpressionDecomposer_testExposeExpression18(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		`const {a, b, c} = condition ?
  y() :
  {a: 0, b: 0, c: 1};
`,
		exprMatchesStr("y()"),
		`var temp$jscomp$0;
if (condition) {
  temp$jscomp$0 = y();
} else {
  temp$jscomp$0 = {a: 0, b: 0, c: 1};
}
const {a, b, c} = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testMoveClass1()
func TestExpressionDecomposer_testMoveClass1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	x.helperMoveExpression(
		"alert(class X {});",
		func(compiler AbstractCompiler) func(*rhino.Node) *rhino.Node {
			return expressionDecomposerTest_findClass
		},
		"var result$jscomp$0 = class X {}; alert(result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testMoveClass2()
func TestExpressionDecomposer_testMoveClass2(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	x.helperMoveExpression(
		"console.log(1, 2, class X {});",
		func(compiler AbstractCompiler) func(*rhino.Node) *rhino.Node {
			return expressionDecomposerTest_findClass
		},
		"var result$jscomp$0 = class X {}; console.log(1, 2, result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testMoveYieldExpression1()
func TestExpressionDecomposer_testMoveYieldExpression1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"function *f() { return { a: yield 1, c: foo(yield 2, yield 3) }; }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var result$jscomp$0 = yield 1;
  return { a: result$jscomp$0, c: foo(yield 2, yield 3) };
}
`)

	x.helperMoveExpression(
		"function *f() { return { a: 0, c: foo(yield 2, yield 3) }; }",
		exprMatchesStr("yield 2"),
		`function *f() {
  var result$jscomp$0 = yield 2;
  return { a: 0, c: foo(result$jscomp$0, yield 3) };
}
`)

	x.helperMoveExpression(
		"function *f() { return { a: 0, c: foo(1, yield 3) }; }",
		exprMatchesStr("yield 3"),
		`function *f() {
  var result$jscomp$0 = yield 3;
  return { a: 0, c: foo(1, result$jscomp$0) };
}
`)
}

// Java: ExpressionDecomposerTest.testMoveYieldExpression2()
func TestExpressionDecomposer_testMoveYieldExpression2(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"function *f() { return (yield 1) || (yield 2); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var result$jscomp$0 = yield 1;
  return result$jscomp$0 || (yield 2);
}
`)
}

// Java: ExpressionDecomposerTest.testMoveYieldExpression3()
func TestExpressionDecomposer_testMoveYieldExpression3(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"function *f() { return x.y(yield 1); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var result$jscomp$0 = yield 1;
  return x.y(result$jscomp$0);
}
`)
}

// Java: ExpressionDecomposerTest.testExposeYieldExpression1()
func TestExpressionDecomposer_testExposeYieldExpression1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"function *f(x) { return x || (yield 2); }",
		exprMatchesStr("yield 2"),
		`function *f(x) {
  var temp$jscomp$0;
  if (temp$jscomp$0=x); else temp$jscomp$0 = yield 2;
  return temp$jscomp$0
}
`)
}

// Java: ExpressionDecomposerTest.testExposeYieldExpression2()
func TestExpressionDecomposer_testExposeYieldExpression2(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"function *f() { return x.y(yield 1); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var temp_const$jscomp$1 = x;
  var temp_const$jscomp$0 = temp_const$jscomp$1.y;
  return temp_const$jscomp$0.call(temp_const$jscomp$1, yield 1);
}
`)
}

// Java: ExpressionDecomposerTest.testExposeYieldExpression3()
func TestExpressionDecomposer_testExposeYieldExpression3(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"function *f() { return g.call(yield 1); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var temp_const$jscomp$1 = g;
  var temp_const$jscomp$0 = temp_const$jscomp$1.call;
  return temp_const$jscomp$0.call(temp_const$jscomp$1, yield 1);
}
`)
}

// Java: ExpressionDecomposerTest.testExposeYieldExpression4()
func TestExpressionDecomposer_testExposeYieldExpression4(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"function *f() { return g.apply([yield 1, yield 2]); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var temp_const$jscomp$1 = g;
  var temp_const$jscomp$0 = temp_const$jscomp$1.apply;
  return temp_const$jscomp$0.call(temp_const$jscomp$1, [yield 1, yield 2]);
}
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals1()
func TestExpressionDecomposer_testExposePlusEquals1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"var x = 0; x += foo() + 1",
		exprMatchesStr("foo()"),
		"var x = 0; var temp_const$jscomp$0 = x; x = temp_const$jscomp$0 + (foo() + 1);")

	x.helperExposeExpression(
		"var x = 0; y = (x += foo()) + x",
		exprMatchesStr("foo()"),
		"var x = 0; var temp_const$jscomp$0 = x; y = (x = temp_const$jscomp$0 + foo()) + x")
}

// Java: ExpressionDecomposerTest.testExposePlusEquals2()
func TestExpressionDecomposer_testExposePlusEquals2(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"var x = {}; x.a += foo() + 1",
		exprMatchesStr("foo()"),
		`var x = {}; var temp_const$jscomp$0 = x;
var temp_const$jscomp$1 = temp_const$jscomp$0.a;
temp_const$jscomp$0.a = temp_const$jscomp$1 + (foo() + 1);
`)

	x.helperExposeExpression(
		"var x = {}; y = (x.a += foo()) + x.a",
		exprMatchesStr("foo()"),
		`var x = {}; var temp_const$jscomp$0 = x;
var temp_const$jscomp$1 = temp_const$jscomp$0.a;
y = (temp_const$jscomp$0.a = temp_const$jscomp$1 + foo()) + x.a
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals3()
func TestExpressionDecomposer_testExposePlusEquals3(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"/** @const */ var XX = {}; XX.a += foo() + 1",
		exprMatchesStr("foo()"),
		`var XX = {};
var temp_const$jscomp$0 = XX.a;
XX.a = temp_const$jscomp$0 + (foo() + 1);
`)

	x.helperExposeExpression(
		"var XX = {}; y = (XX.a += foo()) + XX.a",
		exprMatchesStr("foo()"),
		`var XX = {};
var temp_const$jscomp$0 = XX.a;
y = (XX.a = temp_const$jscomp$0 + foo()) + XX.a
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals4()
func TestExpressionDecomposer_testExposePlusEquals4(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"var x = {}; goo().a += foo() + 1",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo();
var temp_const$jscomp$1 = temp_const$jscomp$0.a;
temp_const$jscomp$0.a = temp_const$jscomp$1 + (foo() + 1);
`)

	x.helperExposeExpression(
		"var x = {}; y = (goo().a += foo()) + goo().a",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo();
var temp_const$jscomp$1 = temp_const$jscomp$0.a;
y = (temp_const$jscomp$0.a = temp_const$jscomp$1 + foo()) + goo().a
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals5()
func TestExpressionDecomposer_testExposePlusEquals5(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"var x = {}; goo().a.b += foo() + 1",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo().a;
var temp_const$jscomp$1 = temp_const$jscomp$0.b;
temp_const$jscomp$0.b = temp_const$jscomp$1 + (foo() + 1);
`)

	x.helperExposeExpression(
		"var x = {}; y = (goo().a.b += foo()) + goo().a",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo().a;
var temp_const$jscomp$1 = temp_const$jscomp$0.b;
y = (temp_const$jscomp$0.b = temp_const$jscomp$1 + foo()) + goo().a
`)
}

// Java: ExpressionDecomposerTest.testExposeLogicalAssignment1()
func TestExpressionDecomposer_testExposeLogicalAssignment1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// Part of the work here is being done by Normalize, which converts all
	// instances of logical assignment operators into an larger expression
	// that separates the logical operation from the assignment.
	x.helperExposeExpression(
		"let x = 0; x ||= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = 0;
if (x) {
} else {
   x = foo() + 1;
}
`)

	x.helperExposeExpression(
		"let x = 0; x &&= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = 0;
if (x) {
   x = foo() + 1;
}
`)

	x.helperExposeExpression(
		"let x = 0; x ??= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = 0;
var temp$jscomp$1;
if ((temp$jscomp$1 = x) != null) {
   temp$jscomp$1;
} else {
   x = foo() + 1;
}
`)
}

// Java: ExpressionDecomposerTest.testExposeLogicalAssignment2()
func TestExpressionDecomposer_testExposeLogicalAssignment2(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// Part of the work here is being done by Normalize, which converts all
	// instances of logical assignment operators into an larger expression
	// that separates the logical operation from the assignment.
	x.helperExposeExpression(
		"let x = {}; x.a ||= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = {};
let $jscomp$logical$assign$tmpm1146332801$0;
if (($jscomp$logical$assign$tmpm1146332801$0 = x).a) {
} else {
   var temp_const$jscomp$1 = $jscomp$logical$assign$tmpm1146332801$0;
   temp_const$jscomp$1.a = foo() + 1;
}
`)
	x.helperExposeExpression(
		"let x = {}; x[a] &&= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = {};
let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
if (($jscomp$logical$assign$tmpm1146332801$0 = x)
    [$jscomp$logical$assign$tmpindexm1146332801$0 = a]) {
    var temp_const$jscomp$2 = $jscomp$logical$assign$tmpm1146332801$0;
    var temp_const$jscomp$1 = $jscomp$logical$assign$tmpindexm1146332801$0;
    temp_const$jscomp$2[temp_const$jscomp$1] = foo() + 1;
}
`)
}

// Java: ExpressionDecomposerTest.testExposeObjectLit1()
func TestExpressionDecomposer_testExposeObjectLit1(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// Validate that getter and setters methods are seen as side-effect
	// free and that values can move past them.  We don't need to be
	// concerned with exposing the getter or setter here but the
	// decomposer does not have a method of exposing properties, only variables.
	x.helperMoveExpression(
		"var x = {get a() {}, b: foo()};",
		exprMatchesStr("foo()"),
		"var result$jscomp$0=foo();var x = {get a() {}, b: result$jscomp$0};")

	x.helperMoveExpression(
		"var x = {set a(p) {}, b: foo()};",
		exprMatchesStr("foo()"),
		"var result$jscomp$0=foo();var x = {set a(p) {}, b: result$jscomp$0};")
}

// Java: ExpressionDecomposerTest.testMoveSpread_siblingOfCall_outOfArrayLiteral_usesTempArray()
func TestExpressionDecomposer_testMoveSpread_siblingOfCall_outOfArrayLiteral_usesTempArray(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	x.helperExposeExpression(
		"[...x, foo()];",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = [...x];
[...temp_const$jscomp$0, foo()];
`)
}

// Java: ExpressionDecomposerTest.testMoveSpread_siblingOfCall_outOfObjectLiteral_usesTempObject()
func TestExpressionDecomposer_testMoveSpread_siblingOfCall_outOfObjectLiteral_usesTempObject(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	x.helperExposeExpression(
		"({...x, y: foo()});",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = {...x};
({...temp_const$jscomp$0, y: foo()});
`)
}

// Java: ExpressionDecomposerTest.testMoveSpread_siblingOfCall_outOfFunctionCall_usesTempArray()
func TestExpressionDecomposer_testMoveSpread_siblingOfCall_outOfFunctionCall_usesTempArray(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	x.helperExposeExpression(
		`function f() { }
f(...x, foo());
`,
		exprMatchesStr("foo()"),
		`function f() { }
var temp_const$jscomp$1 = f;
var temp_const$jscomp$0 = [...x];
temp_const$jscomp$1(...temp_const$jscomp$0, foo());
`)
}

// Java: ExpressionDecomposerTest.testMoveSpreadParent_siblingOfCall_outOfFunctionCall_usesNoTempArray()
func TestExpressionDecomposer_testMoveSpreadParent_siblingOfCall_outOfFunctionCall_usesNoTempArray(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		`function f() { }
f([...x], foo());
`,
		exprMatchesStr("foo()"),
		`function f() { }
var temp_const$jscomp$1 = f;
var temp_const$jscomp$0 = [...x];
temp_const$jscomp$1(temp_const$jscomp$0, foo());
`)
}

// Java: ExpressionDecomposerTest.testMoveSpreadParent_siblingOfCall_outOfFunctionCall_usesNoTempObject()
func TestExpressionDecomposer_testMoveSpreadParent_siblingOfCall_outOfFunctionCall_usesNoTempObject(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		`function f() { }
f({...x}, foo());
`,
		exprMatchesStr("foo()"),
		`function f() { }
var temp_const$jscomp$1 = f;
var temp_const$jscomp$0 = {...x};
temp_const$jscomp$1(temp_const$jscomp$0, foo());
`)
}

// Java: ExpressionDecomposerTest.testExposeExpressionInTemplateLibSub()
func TestExpressionDecomposer_testExposeExpressionInTemplateLibSub(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"` ${ foo() }  ${ goo() } `;",
		exprMatchesStr("goo()"),
		"var temp_const$jscomp$0 = foo(); ` ${ temp_const$jscomp$0 }  ${ goo() } `;")
}

// Java: ExpressionDecomposerTest.testExposeSubExpressionInTemplateLibSub()
func TestExpressionDecomposer_testExposeSubExpressionInTemplateLibSub(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"` ${ foo() + goo() } `;",
		exprMatchesStr("goo()"),
		"var temp_const$jscomp$0 = foo(); ` ${ temp_const$jscomp$0 + goo() } `;")
}

// Java: ExpressionDecomposerTest.testMoveExpressionInTemplateLibSub()
func TestExpressionDecomposer_testMoveExpressionInTemplateLibSub(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"` ${ foo() }  ${ goo() } `;",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); ` ${ result$jscomp$0 }  ${ goo() } `;")
}

// Java: ExpressionDecomposerTest.testExposeExpression_computedProp_withPureKey()
func TestExpressionDecomposer_testExposeExpression_computedProp_withPureKey(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		`({
  ['a' + 'b']: foo(),
});
`,
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testExposeObjectLitValue_computedProp_withImpureKey()
func TestExpressionDecomposer_testExposeObjectLitValue_computedProp_withImpureKey(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		`({
  [goo()]: foo(),
});
`,
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = goo();
({
  [temp_const$jscomp$0]: foo(),
});
`)
}

// Java: ExpressionDecomposerTest.testExposeObjectLitValue_computedProp_asEarlierSibling_withImpureKeyAndValue()
func TestExpressionDecomposer_testExposeObjectLitValue_computedProp_asEarlierSibling_withImpureKeyAndValue(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		`({
  [goo()]: qux(),
  bar: foo(),
});
`,
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$1 = goo();
var temp_const$jscomp$0 = qux();
({
  [temp_const$jscomp$1]: temp_const$jscomp$0,
  bar: foo(),
});
`)
}

// Java: ExpressionDecomposerTest.testExposeObjectLitValue_memberFunctions_asEarlierSiblings_arePure()
func TestExpressionDecomposer_testExposeObjectLitValue_memberFunctions_asEarlierSiblings_arePure(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		`({
  a() { },
  get b() { },
  set b(v) { },

  bar: foo(),
});
`,
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testMoveSuperCall()
func TestExpressionDecomposer_testMoveSuperCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperMoveExpression(
		"class A { constructor() { super(foo()) } }",
		exprMatchesStr("foo()"),
		"class A{constructor(){var result$jscomp$0=foo();super(result$jscomp$0)}}")
}

// Java: ExpressionDecomposerTest.testMoveSuperCall_noSideEffects()
func TestExpressionDecomposer_testMoveSuperCall_noSideEffects(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// String() is being used since it's known to not have side-effects.
	x.helperMoveExpression(
		"class A { constructor() { super(String()) } }",
		exprMatchesStr("String()"),
		"class A{constructor(){var result$jscomp$0=String();super(result$jscomp$0)}}")
}

// Java: ExpressionDecomposerTest.testExposeSuperCall()
func TestExpressionDecomposer_testExposeSuperCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperExposeExpression(
		"class A { constructor() { super(goo(), foo()) } }",
		exprMatchesStr("foo()"),
		`class A{ constructor(){
   var temp_const$jscomp$0=goo();
   super(temp_const$jscomp$0, foo())
}}
`)
}

// Java: ExpressionDecomposerTest.testExposeSuperCall_noSideEffects()
func TestExpressionDecomposer_testExposeSuperCall_noSideEffects(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	// String() is being used since it's known to not have side-effects.
	x.helperExposeExpression(
		"class A { constructor() { super(goo(), String()) } }",
		exprMatchesStr("String()"),
		`class A{ constructor(){
   var temp_const$jscomp$0=goo();
   super(temp_const$jscomp$0, String())
}}
`)
}

// Java: ExpressionDecomposerTest.testCannotDecomposeSuperMethodCall()
func TestExpressionDecomposer_testCannotDecomposeSuperMethodCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"class A extends B { fn() { super.method(foo()) } }",
		exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"class A extends B { fn() { super['method'](foo()) } }",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.canMovePastFnDotCall()
func TestExpressionDecomposer_canMovePastFnDotCall(t *testing.T) {
	x := newExpressionDecomposerTest(t)
	x.knownConstants.Add("fn")
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "fn.call(foo());", exprMatchesStr("foo()"))
	x.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "unknownIfFn.call(foo());", exprMatchesStr("foo()"))
	x.helperMoveExpression(
		"fn.call(foo());",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); fn.call(result$jscomp$0);")
}
