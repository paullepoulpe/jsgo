package jscomp

import (
	"testing"
)

// Go-only: tests of OptimizeCalls' reference map read from the property-use
// index (prop_use_index_oc.go).

// The passes' own suites, with every run reading a freshly built index
// instead of walking: once for their expected outputs alone, once comparing
// every map read with the walk's (verify mode).
func runOptimizeCallsIndexSuite(t *testing.T, newTest func() JUnitTestCase) {
	t.Helper()
	for _, verify := range []bool{false, true} {
		name := "output"
		if verify {
			name = "verify"
		}
		t.Run(name, func(t *testing.T) {
			withPropIndex(t, "", true, verify)
			before := propIndexOCReads.Load()
			RunJUnit4(t, newTest)
			if propIndexOCReads.Load() == before {
				t.Fatal("the index was never read")
			}
		})
	}
}

func TestPropUseIndexOC_OptimizeCallsSuite(t *testing.T) {
	runOptimizeCallsIndexSuite(t, func() JUnitTestCase { return NewOptimizeCallsTest() })
}

func TestPropUseIndexOC_OptimizeParametersSuite(t *testing.T) {
	runOptimizeCallsIndexSuite(t, func() JUnitTestCase { return NewOptimizeParametersTest() })
}

func TestPropUseIndexOC_OptimizeReturnsSuite(t *testing.T) {
	runOptimizeCallsIndexSuite(t, func() JUnitTestCase { return NewOptimizeReturnsTest() })
}

func TestPropUseIndexOC_OptimizeConstructorsSuite(t *testing.T) {
	runOptimizeCallsIndexSuite(t, func() JUnitTestCase { return NewOptimizeConstructorsTest() })
}

func TestPropUseIndexOC_DevirtualizeMethodsSuite(t *testing.T) {
	runOptimizeCallsIndexSuite(t, func() JUnitTestCase { return NewDevirtualizeMethodsTest() })
}

// propIndexOCParamsTest checks OptimizeParameters cases the Java suite does
// not have, with the index and without it.
type propIndexOCParamsTest struct {
	CompilerTestCaseBase
}

func newPropIndexOCParamsTest() *propIndexOCParamsTest {
	c := &propIndexOCParamsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS+"var alert;")
	return c
}

func (c *propIndexOCParamsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewOptimizeParameters(compiler)
}

func (c *propIndexOCParamsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	c.DisableMultistageCompilation()
	c.EnableNormalizeExpectedOutput()
	c.EnableGatherExternProperties()
}

// A method called through an optional chain (an OPTCHAIN_GETPROP reference)
// with another value: its parameter stays.
func (c *propIndexOCParamsTest) TestOptionalChainMethodCall() {
	c.TestSame("class C { m(a) { alert(a); } } new C().m(1); alert?.m(2);")
	c.Test("class C { m(a) { alert(a); } } new C().m(1); alert?.m(1);",
		"class C { m() { var a = 1; alert(a); } } new C().m(); alert?.m();")
}

// A super() call counts as a call of the class it extends, named (a NAME)
// or qualified (a GETPROP).
func (c *propIndexOCParamsTest) TestSuperCalls() {
	c.TestSame("class A { constructor(a) { alert(a); } } class B extends A { constructor() { super(2); } } new A(1); new B();")
	c.TestSame("var ns = {}; ns.A = class { constructor(a) { alert(a); } }; class B extends ns.A { constructor() { super(2); } } new ns.A(1); new B();")
}

func TestPropUseIndexOC_OptimizeParametersExtra(t *testing.T) {
	for _, mode := range []struct {
		name  string
		force bool
	}{{"walk", false}, {"index", true}} {
		t.Run(mode.name, func(t *testing.T) {
			withPropIndex(t, "", mode.force, false)
			RunJUnit4(t, func() JUnitTestCase { return newPropIndexOCParamsTest() })
		})
	}
}
