// Port of test/com/google/javascript/jscomp/ReferenceCollectorTest.java @ closure-compiler v20260415

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// The behaviors assert with the package's Go test helpers on the running
// test method's *testing.T (c.T()).
//
// Java: public final class ReferenceCollectorTest extends CompilerTestCase
type ReferenceCollectorTest struct {
	CompilerTestCaseBase

	behavior ReferenceCollector_Behavior // @Nullable
}

// Java: ReferenceCollectorTest()
func NewReferenceCollectorTest() *ReferenceCollectorTest {
	c := &ReferenceCollectorTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestReferenceCollectorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewReferenceCollectorTest() })
}

// Java: ReferenceCollectorTest.setUp()
func (c *ReferenceCollectorTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.behavior = nil
}

// Java: ReferenceCollectorTest.getProcessor(Compiler)
func (c *ReferenceCollectorTest) GetProcessor(compiler *Compiler) CompilerPass {
	var scopeCreator ScopeCreator = NewSyntacticScopeCreator(compiler)
	return NewReferenceCollector(compiler, c.behavior, scopeCreator)
}

// Java: ReferenceCollectorTest.getOptions()
func (c *ReferenceCollectorTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// Java: ReferenceCollectorTest.testBehavior(String, Behavior)
func (c *ReferenceCollectorTest) testBehavior(js string, behavior ReferenceCollector_Behavior) {
	c.behavior = behavior
	c.TestSame(js)
}

// referenceCollectorTest_isDeclarations is the IS_DECLARATION correspondence:
// the references' isDeclaration(), in order.
func referenceCollectorTest_isDeclarations(refs *ReferenceCollection) []bool {
	var out []bool
	for ref := range refs.All() {
		out = append(out, ref.IsDeclaration())
	}
	return out
}

func (c *ReferenceCollectorTest) TestIterableRest_declaration() {
	t := c.T()
	c.testBehavior(
		"let [x, ...arr] = [1, 2, 3]; [x, ...arr] = [4, 5, 6];",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			arrReferenceCollection := rm.GetReferences(nt.GetScope().GetVar("arr"))
			assertSliceEqual(t, "isDeclaration", referenceCollectorTest_isDeclarations(arrReferenceCollection),
				[]bool{true, false})
		}))
}

func (c *ReferenceCollectorTest) TestObjectRest_declaration() {
	t := c.T()
	c.testBehavior(
		"let {x, ...obj} = {x: 1, y: 2, z: 3}; ({x, ...obj} = {x: 4, q: 5, u: 6});",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			objReferenceCollection := rm.GetReferences(nt.GetScope().GetVar("obj"))
			assertSliceEqual(t, "isDeclaration", referenceCollectorTest_isDeclarations(objReferenceCollection),
				[]bool{true, false})
		}))
}

func (c *ReferenceCollectorTest) TestClass() {
	t := c.T()
	c.testBehavior(
		`class Foo {}
`,
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("Foo"))

				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
				assertEqual(t, "isWellDefined", x.IsWellDefined(), true)
				assertSliceEqual(t, "isDeclaration", referenceCollectorTest_isDeclarations(x), []bool{true})
			}
		}))
}

func (c *ReferenceCollectorTest) TestClass_withPrototype() {
	t := c.T()
	c.testBehavior(
		`class Foo {}
Foo.prototype.bar = 1;
`,
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("Foo"))

				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
				assertEqual(t, "isWellDefined", x.IsWellDefined(), false) // TODO(b/435019132): Should be true.
			}
		}))
}

func (c *ReferenceCollectorTest) TestClass_referencedAfterDeclaration() {
	t := c.T()
	c.testBehavior(
		`class Foo {}
const x = Foo;
`,
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("Foo"))

				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
				assertEqual(t, "isWellDefined", x.IsWellDefined(), false) // TODO(b/435019132): Should be true.
			}
		}))
}

// referenceCollectorTest_importBehavior is the behavior shared by the
// testImport* tests.
func referenceCollectorTest_importBehavior(t *testing.T) ReferenceCollector_Behavior {
	return ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsModuleScope() {
			x := rm.GetReferences(nt.GetScope().GetVar("x"))

			assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
			assertEqual(t, "isWellDefined", x.IsWellDefined(), true)
			assertSliceEqual(t, "isDeclaration", referenceCollectorTest_isDeclarations(x), []bool{true})
		}
	})
}

func (c *ReferenceCollectorTest) TestImport1() {
	t := c.T()
	c.testBehavior("import x from '/m';", referenceCollectorTest_importBehavior(t))
}

func (c *ReferenceCollectorTest) TestImport2() {
	t := c.T()
	c.testBehavior("import {x} from '/m';", referenceCollectorTest_importBehavior(t))
}

func (c *ReferenceCollectorTest) TestImport2_alternate() {
	t := c.T()
	c.testBehavior("import {x as x} from '/m';", referenceCollectorTest_importBehavior(t))
}

func (c *ReferenceCollectorTest) TestImport3() {
	t := c.T()
	c.testBehavior(
		"import {y as x} from '/m';",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsModuleScope() {
				if y := nt.GetScope().GetVar("y"); y != nil {
					t.Errorf("getVar(y) = %v, want null", y)
				}
				x := rm.GetReferences(nt.GetScope().GetVar("x"))

				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
				assertEqual(t, "isWellDefined", x.IsWellDefined(), true)
				assertSliceEqual(t, "isDeclaration", referenceCollectorTest_isDeclarations(x), []bool{true})
			}
		}))
}

func (c *ReferenceCollectorTest) TestImport4() {
	t := c.T()
	c.testBehavior(
		"import * as x from '/m';",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsModuleScope() {
				var_ := nt.GetScope().GetVar("x")
				jutil.CheckNotNull(var_)
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				jutil.CheckNotNull(x)

				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
				assertEqual(t, "isWellDefined", x.IsWellDefined(), true)
				assertSliceEqual(t, "isDeclaration", referenceCollectorTest_isDeclarations(x), []bool{true})
			}
		}))
}

func (c *ReferenceCollectorTest) TestVarInBlock() {
	t := c.T()
	c.testBehavior(
		`function f(x) {
  if (true) {
    var y = x;
    y;
    y;
  }
}
`,
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsBlockScope() && nt.GetScope().GetParent().IsFunctionBlockScope() {
				y := rm.GetReferences(nt.GetScope().GetVar("y"))
				assertEqual(t, "isAssignedOnceInLifetime", y.IsAssignedOnceInLifetime(), true)
				assertEqual(t, "isWellDefined", y.IsWellDefined(), true)
			}
		}))
}

func (c *ReferenceCollectorTest) TestVarInLoopNotAssignedOnlyOnceInLifetime() {
	t := c.T()
	behavior := ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsBlockScope() {
			x := rm.GetReferences(nt.GetScope().GetVar("x"))
			assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), false)
		}
	})
	c.testBehavior("while (true) { var x = 0; }", behavior)
}

func (c *ReferenceCollectorTest) TestLetInLoopAssignedOnceInLifetime() {
	t := c.T()
	behavior := ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsBlockScope() {
			xVar := nt.GetScope().GetVar("x")
			if xVar == nil {
				t.Fatal("getVar(x) is null")
			}
			assertEqual(t, "isLet", xVar.IsLet(), true)
			x := rm.GetReferences(xVar)
			assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
		}
	})
	c.testBehavior("while (true) { let x = 0; }", behavior)
}

// TestReferenceCollector_VarInFunctionNotAssignedOnlyOnceInLifetime:
// although there is only one assignment to x in the code, it's in a function
// which could be called multiple times, so isAssignedOnceInLifetime() returns
// false.
func (c *ReferenceCollectorTest) TestVarInFunctionNotAssignedOnlyOnceInLifetime() {
	t := c.T()
	behavior := ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsGlobal() {
			x := rm.GetReferences(nt.GetScope().GetVar("x"))
			assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), false)
		}
	})
	c.testBehavior("var x; function f() { x = 0; }", behavior)
	c.testBehavior("let x; function f() { x = 0; }", behavior)
}

func (c *ReferenceCollectorTest) TestParameterAssignedOnlyOnceInLifetime() {
	t := c.T()
	c.testBehavior(
		"function f(x) { x; }",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsFunctionScope() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
			}
		}))
}

func (c *ReferenceCollectorTest) TestModifiedParameterNotAssignedOnlyOnceInLifetime() {
	t := c.T()
	c.testBehavior(
		"function f(x) { x = 3; }",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsFunctionScope() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), false)
			}
		}))
}

func (c *ReferenceCollectorTest) TestVarAssignedOnceInLifetime1() {
	t := c.T()
	behavior := ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsFunctionBlockScope() {
			x := rm.GetReferences(nt.GetScope().GetVar("x"))
			assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
		}
	})
	c.testBehavior("function f() { var x = 0; }", behavior)
	c.testBehavior("function f() { let x = 0; }", behavior)
}

func (c *ReferenceCollectorTest) TestVarAssignedOnceInLifetime2() {
	t := c.T()
	c.testBehavior(
		"function f() { { let x = 0; } }",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsBlockScope() && !nt.GetScope().IsFunctionBlockScope() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
			}
		}))
}

// referenceCollectorTest_catchBehavior is the behavior of
// testVarAssignedOnceInLifetime3 and testLetAssignedOnceInLifetime1/2.
func referenceCollectorTest_catchBehavior(t *testing.T) ReferenceCollector_Behavior {
	return ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsCatchScope() {
			e := rm.GetReferences(nt.GetScope().GetVar("e"))
			assertEqual(t, "e.isAssignedOnceInLifetime", e.IsAssignedOnceInLifetime(), true)
			y := rm.GetReferences(nt.GetScope().GetVar("y"))
			assertEqual(t, "y.isAssignedOnceInLifetime", y.IsAssignedOnceInLifetime(), true)
			assertEqual(t, "y.isWellDefined", y.IsWellDefined(), true)
		}
	})
}

func (c *ReferenceCollectorTest) TestVarAssignedOnceInLifetime3() {
	t := c.T()
	behavior := referenceCollectorTest_catchBehavior(t)
	c.testBehavior(
		`try {
} catch (e) {
  var y = e;
  g();
  y;y;
}
`,
		behavior)
	c.testBehavior(
		`try {
} catch (e) {
  var y; y = e;
  g();
  y;y;
}
`,
		behavior)
}

func (c *ReferenceCollectorTest) TestLetAssignedOnceInLifetime1() {
	t := c.T()
	c.testBehavior(
		`try {
} catch (e) {
  let y = e;
  g();
  y;y;
}
`,
		referenceCollectorTest_catchBehavior(t))
}

func (c *ReferenceCollectorTest) TestLetAssignedOnceInLifetime2() {
	t := c.T()
	c.testBehavior(
		`try {
} catch (e) {
  let y; y = e;
  g();
  y;y;
}
`,
		referenceCollectorTest_catchBehavior(t))
}

// referenceCollectorTest_assertBlockRoots checks the tokens of the roots of
// the basic blocks of the references to x.
func referenceCollectorTest_assertBlockRoots(t *testing.T, x *ReferenceCollection, want ...rhino.Token) {
	t.Helper()
	if x.references.Size() != len(want) {
		t.Fatalf("%d references, want %d", x.references.Size(), len(want))
	}
	for i, tok := range want {
		assertEqual(t, "basic block root token", x.references.Get(i).GetBasicBlock().GetRoot().GetToken(), tok)
	}
}

func (c *ReferenceCollectorTest) TestBasicBlocks() {
	t := c.T()
	c.testBehavior(
		`var x = 0;
switch (x) {
  case 0:
    x;
}
`,
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				referenceCollectorTest_assertBlockRoots(t, x, rhino.Token_ROOT, rhino.Token_ROOT, rhino.Token_CASE)
			}
		}))
}

func (c *ReferenceCollectorTest) TestBasicBlocksInConditionals() {
	t := c.T()
	c.testBehavior(
		`var x = 0;
x || 3;
3 || x;
const [y = (x = 1)] = [];
`,
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				// first child of || is not a boundary, but the second child is.
				// second child of `y = (x = 1)` is a boundary
				referenceCollectorTest_assertBlockRoots(t, x,
					rhino.Token_ROOT, rhino.Token_ROOT, rhino.Token_NAME, rhino.Token_ASSIGN)
			}
		}))
}

func (c *ReferenceCollectorTest) TestNullishCoalesce() {
	t := c.T()
	c.testBehavior(
		"var x = 0; var y = x ?? (x = 1)",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				// first child of ?? is not a boundary, but the second child is.
				referenceCollectorTest_assertBlockRoots(t, x, rhino.Token_ROOT, rhino.Token_ROOT, rhino.Token_ASSIGN)
			}
		}))
}

func (c *ReferenceCollectorTest) TestOptChain() {
	t := c.T()
	c.testBehavior(
		"var x = 0; var y = x?.(x = 1)",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				// first child of `?.` is not a boundary, but the second child is.
				referenceCollectorTest_assertBlockRoots(t, x, rhino.Token_ROOT, rhino.Token_ROOT, rhino.Token_ASSIGN)
			}
		}))
}

func (c *ReferenceCollectorTest) TestThis() {
	t := c.T()
	c.testBehavior(
		`/** @constructor */
function C() {}

C.prototype.m = function m() {
  var self = this;
  if (true) {
    alert(self);
  }
};
`,
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsFunctionBlockScope() &&
				nt.GetScopeRoot().GetParent().GetFirstChild().MatchesName("m") {
				self := rm.GetReferences(nt.GetScope().GetVar("self"))
				assertEqual(t, "isEscaped", self.IsEscaped(), false)
			}
		}))
}

func (c *ReferenceCollectorTest) TestProcessScopeThatsNotABasicBlock() {
	t := c.T()
	// Tests the case where the scope we pass in is not really a basic block, but we create a new
	// basic block anyway because ReferenceCollector expects all nodes to be in a block.
	compiler := c.CreateCompiler()
	syntacticScopeCreator := NewSyntacticScopeCreator(compiler)
	referenceCollectingCallback := NewReferenceCollector(
		compiler,
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			y := rm.GetReferences(nt.GetScope().GetVar("y"))
			assertEqual(t, "isWellDefined", y.IsWellDefined(), true)
			firstBasicBlock := y.references.Get(0).GetBasicBlock()
			assertEqual(t, "first block root", firstBasicBlock.GetRoot().GetToken(), rhino.Token_BLOCK)
			assertEqual(t, "first block root parent", firstBasicBlock.GetRoot().GetParent().GetToken(), rhino.Token_SCRIPT)

			// We do not create a new BasicBlock for the second { use(y); }
			secondBasicBlock := y.references.Get(1).GetBasicBlock()
			assertEqual(t, "second block root", secondBasicBlock.GetRoot().GetToken(), rhino.Token_BLOCK)
			assertEqual(t, "second block root parent", secondBasicBlock.GetRoot().GetParent().GetToken(), rhino.Token_IF)
		}),
		syntacticScopeCreator)

	js := "let x = 5; { let y = x + 1; if (true) { { use(y); } } }"
	root := compiler.ParseTestCode(js)
	block := root.GetSecondChild()

	globalScope := syntacticScopeCreator.CreateScopeOfScope(root, nil)
	blockScope := syntacticScopeCreator.CreateScopeOfScope(block, globalScope)
	referenceCollectingCallback.ProcessScope(blockScope)
}
