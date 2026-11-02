package jscomp

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: tests of the property-use index (prop_use_index.go).

// withPropIndex sets the index switches for one test.
func withPropIndex(t *testing.T, env string, force, verify bool) {
	t.Helper()
	savedEnv, savedForce, savedVerify := propIndexEnv, propIndexForce, propIndexVerify
	propIndexEnv, propIndexForce, propIndexVerify = env, force, verify
	t.Cleanup(func() { propIndexEnv, propIndexForce, propIndexVerify = savedEnv, savedForce, savedVerify })
}

// The passes' own suites, with every run reading a freshly built index instead
// of walking: the replay must give the walks' results.
func TestPropUseIndex_InlineSimpleMethodsSuite(t *testing.T) {
	withPropIndex(t, "", true, true)
	before := propIndexReads.Load()
	RunJUnit4(t, func() JUnitTestCase { return NewInlineSimpleMethodsTest() })
	if propIndexReads.Load() == before {
		t.Fatal("the index was never read")
	}
}

func TestPropUseIndex_InlinePropertiesSuite(t *testing.T) {
	withPropIndex(t, "", true, true)
	before := propIndexReads.Load()
	RunJUnit4(t, func() JUnitTestCase { return NewInlinePropertiesTest() })
	if propIndexReads.Load() == before {
		t.Fatal("the index was never read")
	}
}

func TestPropUseIndex_DeadPropertyAssignmentEliminationSuite(t *testing.T) {
	withPropIndex(t, "", true, true)
	before := propIndexReads.Load()
	RunJUnit4(t, func() JUnitTestCase { return NewDeadPropertyAssignmentEliminationTest() })
	if propIndexReads.Load() == before {
		t.Fatal("the index was never read")
	}
}

// propIndexISMTest checks InlineSimpleMethods cases the Java suite does not
// have, with the index and without it.
type propIndexISMTest struct {
	CompilerTestCaseBase
}

func newPropIndexISMTest() *propIndexISMTest {
	c := &propIndexISMTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func (c *propIndexISMTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewInlineSimpleMethods(compiler)
}

func (c *propIndexISMTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableGatherExternProperties()
	c.EnableNormalize()
	c.EnableNormalizeExpectedOutput()
}

const propIndexISMDefs = "function A(){}\n" +
	"A.prototype.getF = function() { return this.f; };\n" +
	"A.prototype.getG = function() { return this.g; };\n" +
	"A.prototype.f = function() { return this.b; };\n" +
	"A.prototype.nop = function() {};\n" +
	"var a = new A();\n"

// Inlining a callee that is itself a call (`a.getF()()`) gives the outer call
// a GETPROP callee; the traversal visits the outer call afterwards and inlines
// it too, after the calls in its arguments.
func (c *propIndexISMTest) TestInlinedCalleeCall() {
	c.Test(propIndexISMDefs+"a.getF()();", propIndexISMDefs+"a.b;")
}

func (c *propIndexISMTest) TestInlinedCalleeCallWithArgs() {
	c.Test(propIndexISMDefs+"a.getF()(a.getG());", propIndexISMDefs+"a.b;")
	c.Test(propIndexISMDefs+"x(a.getF()(a.getG()), a.getG());", propIndexISMDefs+"x(a.b, a.g);")
	c.Test(propIndexISMDefs+"a.getF()(a.nop(a.getF()()));", propIndexISMDefs+"a.b;")
}

// An argument's call is inlined before the call it is an argument of, which
// then has no side effects in its arguments.
func (c *propIndexISMTest) TestArgumentsFirst() {
	c.Test(propIndexISMDefs+"a.nop(a.getG()); a.f(a.getF()());", propIndexISMDefs+"; a.b;")
}

func TestPropUseIndex_InlineSimpleMethodsExtra(t *testing.T) {
	for _, mode := range []struct {
		name  string
		force bool
	}{{"walk", false}, {"index", true}} {
		t.Run(mode.name, func(t *testing.T) {
			withPropIndex(t, "", mode.force, mode.force)
			RunJUnit4(t, func() JUnitTestCase { return newPropIndexISMTest() })
		})
	}
}

// propIndexIPTest checks InlineProperties cases where the traversal state
// the replay restores (global hoist scope, enclosing function) decides.
type propIndexIPTest struct {
	CompilerTestCaseBase
}

func newPropIndexIPTest() *propIndexIPTest {
	c := &propIndexIPTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, InlinePropertiesTest_EXTERNS)
	return c
}

func (c *propIndexIPTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(NewDiagnosticGroup(TypeCheck_INEXISTENT_PROPERTY, TypeValidator_TYPE_MISMATCH_WARNING), CheckLevel_OFF)
	return options
}

func (c *propIndexIPTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewInlineProperties(compiler)
}

func (c *propIndexIPTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableNormalize()
	c.EnableGatherExternProperties()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()
}

func (c *propIndexIPTest) TestPrototypeDefinitionScopes() {
	// In the global hoist scope (a block included): a definition.
	c.Test("/** @constructor */ function C() {} { C.prototype.foo = 1; } new C().foo;",
		"/** @constructor */ function C() {} { C.prototype.foo = 1; } new C(), 1;")
	// In a function, a class static block or a field initializer: not one.
	c.TestSame("/** @constructor */ function C() {} function f() { C.prototype.foo = 1; } f(); new C().foo;")
	c.TestSame("/** @constructor */ function C() {} class D { static { C.prototype.foo = 1; } } new C().foo;")
	c.TestSame("/** @constructor */ function C() {} class D { x = (C.prototype.foo = 1); } new D(); new C().foo;")
	c.TestSame("/** @constructor */ function C() {} class D { static x = (C.prototype.foo = 1); } new C().foo;")
}

func (c *propIndexIPTest) TestThisDefinitionScopes() {
	// Directly in the constructor: a definition; in a nested arrow: not one.
	c.Test("/** @constructor */ function C() { this.foo = 1; } new C().foo;",
		"/** @constructor */ function C() { this.foo = 1; } new C(), 1;")
	c.TestSame("/** @constructor */ function C() { (() => { this.foo = 1; })(); } new C().foo;")
	c.TestSame("/** @constructor */ function C() { function g() { this.foo = 1; } } new C().foo;")
}

func TestPropUseIndex_InlinePropertiesExtra(t *testing.T) {
	for _, mode := range []struct {
		name  string
		force bool
	}{{"walk", false}, {"index", true}} {
		t.Run(mode.name, func(t *testing.T) {
			withPropIndex(t, "", mode.force, mode.force)
			RunJUnit4(t, func() JUnitTestCase { return newPropIndexIPTest() })
		})
	}
}

const propIndexLoopExterns = `
/** @param {...*} x */ function alert(x) {}
/** @constructor */ function Element() {}
/** @type {string} */ Element.prototype.id;
/** @param {string} p @param {*=} o @return {string} */ function JSCompiler_renameProperty(p, o) {}
`

// Programs whose main optimization loop runs several iterations with
// inlining, inlined getters and constant properties, moved and deleted
// functions.
var propIndexLoopInputs = []string{`
/** @constructor */
function Point(x, y) { this.x = x; this.y = y; this.kind = 3; }
Point.prototype.getX = function() { return this.x; };
Point.prototype.getY = function() { return this.y; };
Point.prototype.getKind = function() { return this.kind; };
Point.prototype.nop = function() {};
Point.ORIGIN_KIND = 7;
Point.prototype.LIMIT = 100;
/** @return {!Point} */
function make(i) { return new Point(i, i + 1); }
function sum(ps) {
  var s = 0;
  for (var i = 0; i < ps.length; i++) {
    ps[i].nop();
    s += ps[i].getX() + ps[i].getY() * ps[i].getKind() + Point.ORIGIN_KIND + ps[i].LIMIT;
  }
  return s;
}
function helper(p) { return p.getX(); }
function unused() { return make(5).getY(); }
alert(sum([make(1), make(2)]), helper(make(3)));
`, `
class Shape {
  constructor(n) { this.n = n; this.sides = 4; }
  name() { return this.n; }
  area() { return 0; }
  empty() {}
  static create(n) { return new Shape(n); }
}
class Square extends Shape {
  constructor(n, w) { super(n); this.w = w; }
  width() { return this.w; }
}
const shapes = [Shape.create('a'), new Square('b', 3)];
function describe(s) {
  s.empty();
  return s.name() + ':' + s.area() + ':' + s.sides;
}
function inlineMe(s) { return describe(s) + (s instanceof Square ? s.width() : 0); }
alert(shapes.map(inlineMe).join(','));
`, `
var ns = {};
ns.config = { debug: false, level: 2 };
ns.Widget = function(el) { this.el = el; this.count = 0; };
ns.Widget.prototype.getEl = function() { return this.el; };
ns.Widget.prototype.id = function() { return this.getEl().id; };
ns.Widget.prototype.inc = function() { this.count++; return this; };
ns.Widget.prototype.render = function() {
  if (ns.config.debug) { alert('render', this.id()); }
  var f = function(w) { return w.getEl(); };
  return f(this).id + this.inc().count;
};
function wrap(el) { return new ns.Widget(el); }
var w = wrap(new Element());
alert(w.render(), w.getEl().id, JSCompiler_renameProperty('count', w));
`}

func compileForPropIndexTest(t *testing.T, input string) string {
	t.Helper()
	options := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	options.SetCheckTypes(true)
	c := NewCompiler()
	result := c.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", propIndexLoopExterns)},
		[]*SourceFile{SourceFile_FromCode("in.js", input)},
		options)
	if !result.Success {
		t.Fatalf("compilation failed: %v", result.Errors)
	}
	return c.ToSource()
}

// In the loop the index is kept current from the change reports; with the
// verify mode it is checked against a rebuild after every loop pass, and the
// output must be the walks' output.
func TestPropUseIndex_LoopMatchesWalks(t *testing.T) {
	for i, input := range propIndexLoopInputs {
		withPropIndex(t, "0", false, false)
		want := compileForPropIndexTest(t, input)
		withPropIndex(t, "", false, true)
		before := propIndexReads.Load()
		got := compileForPropIndexTest(t, input)
		reads := propIndexReads.Load() - before
		if got != want {
			t.Errorf("input %d: with the index:\n%s\nwalking:\n%s", i, got, want)
		}
		if reads < 2 {
			t.Errorf("input %d: the index was read %d times", i, reads)
		}
	}
}

// The update re-indexes reported scopes, indexes new functions, drops
// deleted ones, and the verify mode catches an unreported change.
func TestPropUseIndex_UpdateAndVerify(t *testing.T) {
	withPropIndex(t, "", false, false)
	options := NewCompilerOptions()
	c := NewCompiler()
	result := c.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", propIndexLoopExterns)},
		[]*SourceFile{SourceFile_FromCode("in.js", "var o = {a: 1}; function f() { o.b = o.a; }\nf();")},
		options)
	if !result.Success {
		t.Fatalf("compilation failed: %v", result.Errors)
	}
	ix := c.refreshPropUseIndex(c.externsRoot, c.jsRoot)
	ix.verify("after build")
	script := c.jsRoot.GetFirstChild()

	// A reported statement in the script, with a new function in it.
	fn := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(),
		rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(rhino.IR_Getprop(rhino.IR_ThisNode(), "c"))))
	stmt := rhino.IR_ExprResult(rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name("o"), "d"), fn))
	script.AddChildToBack(stmt)
	c.ReportChangeToEnclosingScope(stmt)
	c.changeTracker.MarkNewScopesChanged(fn)
	ix = c.refreshPropUseIndex(c.externsRoot, c.jsRoot)
	ix.verify("after an added statement")
	if ix.scopes[fn] == nil {
		t.Fatal("the new function was not indexed")
	}

	// A reported deletion of that function.
	stmt.Detach()
	c.ReportChangeToChangeScope(script)
	NodeUtil_MarkFunctionsDeleted(stmt, c)
	ix = c.refreshPropUseIndex(c.externsRoot, c.jsRoot)
	ix.verify("after a deleted statement")

	// An unreported change.
	script.AddChildToBack(rhino.IR_ExprResult(rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name("o"), "e"), rhino.IR_Number(1))))
	ix = c.refreshPropUseIndex(c.externsRoot, c.jsRoot)
	defer func() {
		r := recover()
		msg, _ := r.(string)
		if !strings.Contains(msg, "differs from a rebuild") {
			t.Fatalf("the unreported change was not caught: %v", r)
		}
	}()
	ix.verify("after an unreported change")
}
