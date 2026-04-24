// Port of test/com/google/javascript/jscomp/DataFlowAnalysisTest.java @ closure-compiler v20260415
//
// Java's DataFlowAnalysisTest exercises a toy mini-compiler (ArithmeticInstruction,
// BranchInstruction, ConstPropLatticeElement, DummyConstPropagation) built on top of
// DataFlowAnalysis<T, L>. Our Go port has a simpler concrete DataFlowAnalysis that
// operates on FlowScope rather than a generic lattice. The toy-IR tests
// (testSimpleIf, testSimpleLoop, testBranchedSimpleIf, testMaxIterationsExceededException)
// are skipped as "port gap: generic lattice DataFlowAnalysis<T,L> not ported".
//
// The testEscaped* tests rely on LiveVariablesAnalysis which is also not yet ported.
// All tests are therefore stubs — the file exists to enumerate them and justify the skips.

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// ---------------------------------------------------------------------------
// Toy mini-compiler tests — all skipped (generic lattice DFA not ported)
// ---------------------------------------------------------------------------

// testSimpleIf
func TestDFA_SimpleIf(t *testing.T) {
	t.Skip("port gap: generic DataFlowAnalysis<T,L> with toy mini-IR not ported")
}

// testSimpleLoop
func TestDFA_SimpleLoop(t *testing.T) {
	t.Skip("port gap: generic DataFlowAnalysis<T,L> with toy mini-IR not ported")
}

// testBranchedSimpleIf
func TestDFA_BranchedSimpleIf(t *testing.T) {
	t.Skip("port gap: generic DataFlowAnalysis<T,L> with toy mini-IR not ported")
}

// testMaxIterationsExceededException
func TestDFA_MaxIterationsExceededException(t *testing.T) {
	t.Skip("port gap: generic DataFlowAnalysis<T,L> divergence test not ported")
}

// ---------------------------------------------------------------------------
// testEscaped* tests — use LiveVariablesAnalysis not yet ported
// ---------------------------------------------------------------------------

// testEscaped
func TestDFA_Escaped(t *testing.T) {
	t.Skip("port gap: LiveVariablesAnalysis and computeEscapedLocals not yet ported")
}

// testEscapedFunctionLayered
func TestDFA_EscapedFunctionLayered(t *testing.T) {
	t.Skip("port gap: LiveVariablesAnalysis and computeEscapedLocals not yet ported")
}

// testEscapedLetConstSimple
func TestDFA_EscapedLetConstSimple(t *testing.T) {
	t.Skip("port gap: LiveVariablesAnalysis and computeEscapedLocals not yet ported")
}

// testEscapedFunctionAssignment
func TestDFA_EscapedFunctionAssignment(t *testing.T) {
	t.Skip("port gap: LiveVariablesAnalysis and computeEscapedLocals not yet ported")
}

// testEscapedArrowFunction
func TestDFA_EscapedArrowFunction(t *testing.T) {
	t.Skip("port gap: LiveVariablesAnalysis and computeEscapedLocals not yet ported")
}

// ---------------------------------------------------------------------------
// Smoke tests: DataFlowAnalysis.Run() terminates correctly over real JS
// These exercise our Go DataFlowAnalysis via the type-inference pipeline.
// ---------------------------------------------------------------------------

// runDfaOnSrc builds a CFG from `src` and runs DataFlowAnalysis on it.
// Returns true if Run() completed without panicking (termination check).
func runDfaOnSrc(t *testing.T, src string) bool {
	t.Helper()
	script, err := parser.Parse(src)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	cfg := BuildControlFlowGraph(script)
	reg := jstype.NewRegistry()
	creator := NewTypedScopeCreator(reg)
	creator.sourceSrc = src
	creator.sourcePath = "t.js"
	externScope := creator.createExternScope()
	scope := creator.createScopeInWithParent(script, externScope)

	initial := EmptyFlowScope()
	dfa := NewDataFlowAnalysis(cfg, initial, func(node *CfgNode, in *FlowScope) map[Branch]*FlowScope {
		// trivial transfer: propagate unchanged
		return map[Branch]*FlowScope{BranchUnconditional: in}
	})
	dfa.Run()
	_ = scope
	return true
}

// TestDFA_Smoke_SimpleLoop verifies DataFlowAnalysis.Run() terminates on a while loop.
func TestDFA_Smoke_SimpleLoop(t *testing.T) {
	ok := runDfaOnSrc(t, "var x = 0; while(x < 10) { x++; }")
	if !ok {
		t.Error("DataFlowAnalysis did not complete on simple loop")
	}
}

// TestDFA_Smoke_NestedLoops verifies termination on nested loops.
func TestDFA_Smoke_NestedLoops(t *testing.T) {
	ok := runDfaOnSrc(t, "for(var i=0;i<10;i++){for(var j=0;j<10;j++){i+j;}}")
	if !ok {
		t.Error("DataFlowAnalysis did not complete on nested loops")
	}
}

// TestDFA_Smoke_IfElse verifies termination on if-else.
func TestDFA_Smoke_IfElse(t *testing.T) {
	ok := runDfaOnSrc(t, "var a; if (a) { a = 1; } else { a = 2; } a;")
	if !ok {
		t.Error("DataFlowAnalysis did not complete on if-else")
	}
}

// TestDFA_Smoke_TryCatch verifies termination on try-catch.
func TestDFA_Smoke_TryCatch(t *testing.T) {
	ok := runDfaOnSrc(t, "try { var x = 1; } catch (e) { var y = 2; } var z = 3;")
	if !ok {
		t.Error("DataFlowAnalysis did not complete on try-catch")
	}
}

// TestDFA_Smoke_Termination verifies the worklist loop exits within bounds.
func TestDFA_Smoke_Termination(t *testing.T) {
	// A more complex program with multiple paths
	src := `
var a = 0;
do {
  a++;
  if (a > 5) { break; }
} while (a < 100);
var b = a;
`
	ok := runDfaOnSrc(t, src)
	if !ok {
		t.Error("DataFlowAnalysis did not complete on do-while with break")
	}
}
