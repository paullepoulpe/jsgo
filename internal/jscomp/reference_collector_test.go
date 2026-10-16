// Port of test/com/google/javascript/jscomp/ReferenceCollectorTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// referenceCollectorTest_testBehavior is ReferenceCollectorTest.testBehavior:
// testSame(js) with getProcessor(compiler) = new ReferenceCollector(compiler,
// behavior, new SyntacticScopeCreator(compiler)).
//
// Go-only: test-helper — CompilerTestCase is not ported, so this runs the
// part of testSame the behaviors observe: the externs ("") and the js are
// parsed as CompilerTestCase's inputs (CompilerInputs of a Compiler without
// options, i.e. language UNSUPPORTED), put under ROOT(ROOT(externs),
// ROOT(js)), and the processor runs on (externsRoot, mainRoot). It does not
// check that the code is unchanged or that no warnings were reported.
func referenceCollectorTest_testBehavior(t *testing.T, js string, behavior ReferenceCollector_Behavior) {
	t.Helper()
	compiler := newTestCompilerUninitialized()
	compiler.initCompilerOptionsIfTesting()
	externsInput := NewCompilerInputWithIsExtern(SourceFile_FromCode("externs", ""), true)
	compiler.putCompilerInput(externsInput)
	mainInput := NewCompilerInput(SourceFile_FromCode("testcode", js))
	compiler.putCompilerInput(mainInput)
	externsRoot := rhino.IR_Root(externsInput.GetAstRoot(compiler))
	mainRoot := rhino.IR_Root(mainInput.GetAstRoot(compiler))
	rhino.IR_Root(externsRoot, mainRoot)
	if errors := compiler.GetErrors(); len(errors) != 0 {
		t.Fatalf("parse errors: %v", errors)
	}

	var processor CompilerPass = NewReferenceCollector(compiler, behavior, NewSyntacticScopeCreator(compiler))
	processor.Process(externsRoot, mainRoot)
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

func TestReferenceCollector_IterableRest_declaration(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
		"let [x, ...arr] = [1, 2, 3]; [x, ...arr] = [4, 5, 6];",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			arrReferenceCollection := rm.GetReferences(nt.GetScope().GetVar("arr"))
			assertSliceEqual(t, "isDeclaration", referenceCollectorTest_isDeclarations(arrReferenceCollection),
				[]bool{true, false})
		}))
}

func TestReferenceCollector_ObjectRest_declaration(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
		"let {x, ...obj} = {x: 1, y: 2, z: 3}; ({x, ...obj} = {x: 4, q: 5, u: 6});",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			objReferenceCollection := rm.GetReferences(nt.GetScope().GetVar("obj"))
			assertSliceEqual(t, "isDeclaration", referenceCollectorTest_isDeclarations(objReferenceCollection),
				[]bool{true, false})
		}))
}

func TestReferenceCollector_Class(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_Class_withPrototype(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_Class_referencedAfterDeclaration(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_Import1(t *testing.T) {
	referenceCollectorTest_testBehavior(t, "import x from '/m';", referenceCollectorTest_importBehavior(t))
}

func TestReferenceCollector_Import2(t *testing.T) {
	referenceCollectorTest_testBehavior(t, "import {x} from '/m';", referenceCollectorTest_importBehavior(t))
}

func TestReferenceCollector_Import2_alternate(t *testing.T) {
	referenceCollectorTest_testBehavior(t, "import {x as x} from '/m';", referenceCollectorTest_importBehavior(t))
}

func TestReferenceCollector_Import3(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_Import4(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_VarInBlock(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_VarInLoopNotAssignedOnlyOnceInLifetime(t *testing.T) {
	behavior := ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsBlockScope() {
			x := rm.GetReferences(nt.GetScope().GetVar("x"))
			assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), false)
		}
	})
	referenceCollectorTest_testBehavior(t, "while (true) { var x = 0; }", behavior)
}

func TestReferenceCollector_LetInLoopAssignedOnceInLifetime(t *testing.T) {
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
	referenceCollectorTest_testBehavior(t, "while (true) { let x = 0; }", behavior)
}

// TestReferenceCollector_VarInFunctionNotAssignedOnlyOnceInLifetime:
// although there is only one assignment to x in the code, it's in a function
// which could be called multiple times, so isAssignedOnceInLifetime() returns
// false.
func TestReferenceCollector_VarInFunctionNotAssignedOnlyOnceInLifetime(t *testing.T) {
	behavior := ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsGlobal() {
			x := rm.GetReferences(nt.GetScope().GetVar("x"))
			assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), false)
		}
	})
	referenceCollectorTest_testBehavior(t, "var x; function f() { x = 0; }", behavior)
	referenceCollectorTest_testBehavior(t, "let x; function f() { x = 0; }", behavior)
}

func TestReferenceCollector_ParameterAssignedOnlyOnceInLifetime(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
		"function f(x) { x; }",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsFunctionScope() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
			}
		}))
}

func TestReferenceCollector_ModifiedParameterNotAssignedOnlyOnceInLifetime(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
		"function f(x) { x = 3; }",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsFunctionScope() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), false)
			}
		}))
}

func TestReferenceCollector_VarAssignedOnceInLifetime1(t *testing.T) {
	behavior := ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
		if nt.GetScope().IsFunctionBlockScope() {
			x := rm.GetReferences(nt.GetScope().GetVar("x"))
			assertEqual(t, "isAssignedOnceInLifetime", x.IsAssignedOnceInLifetime(), true)
		}
	})
	referenceCollectorTest_testBehavior(t, "function f() { var x = 0; }", behavior)
	referenceCollectorTest_testBehavior(t, "function f() { let x = 0; }", behavior)
}

func TestReferenceCollector_VarAssignedOnceInLifetime2(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_VarAssignedOnceInLifetime3(t *testing.T) {
	behavior := referenceCollectorTest_catchBehavior(t)
	referenceCollectorTest_testBehavior(t,
		`try {
} catch (e) {
  var y = e;
  g();
  y;y;
}
`,
		behavior)
	referenceCollectorTest_testBehavior(t,
		`try {
} catch (e) {
  var y; y = e;
  g();
  y;y;
}
`,
		behavior)
}

func TestReferenceCollector_LetAssignedOnceInLifetime1(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
		`try {
} catch (e) {
  let y = e;
  g();
  y;y;
}
`,
		referenceCollectorTest_catchBehavior(t))
}

func TestReferenceCollector_LetAssignedOnceInLifetime2(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_BasicBlocks(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_BasicBlocksInConditionals(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_NullishCoalesce(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
		"var x = 0; var y = x ?? (x = 1)",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				// first child of ?? is not a boundary, but the second child is.
				referenceCollectorTest_assertBlockRoots(t, x, rhino.Token_ROOT, rhino.Token_ROOT, rhino.Token_ASSIGN)
			}
		}))
}

func TestReferenceCollector_OptChain(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
		"var x = 0; var y = x?.(x = 1)",
		ReferenceCollector_BehaviorFunc(func(nt *NodeTraversal, rm ReferenceMap) {
			if nt.GetScope().IsGlobal() {
				x := rm.GetReferences(nt.GetScope().GetVar("x"))
				// first child of `?.` is not a boundary, but the second child is.
				referenceCollectorTest_assertBlockRoots(t, x, rhino.Token_ROOT, rhino.Token_ROOT, rhino.Token_ASSIGN)
			}
		}))
}

func TestReferenceCollector_This(t *testing.T) {
	referenceCollectorTest_testBehavior(t,
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

func TestReferenceCollector_ProcessScopeThatsNotABasicBlock(t *testing.T) {
	// Tests the case where the scope we pass in is not really a basic block, but we create a new
	// basic block anyway because ReferenceCollector expects all nodes to be in a block.
	// Java: createCompiler() of CompilerTestCase; the test only parses with it.
	compiler := newTestCompilerUninitialized()
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
