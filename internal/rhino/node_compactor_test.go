package rhino

import (
	"testing"
)

// Go-only: in-place compaction of the node store (node_compact.go).
func TestCompactorKeepsLiveTreesAndForwardsHolders(t *testing.T) {
	prev := SetCurrentNodeStore(NewNodeStore())
	defer SetCurrentNodeStore(prev)
	registry := NewJSTypeRegistry(nil)
	num := registry.GetNativeType(JSTypeNative_NUMBER_TYPE)

	for i := 0; i < 1000; i++ {
		IR_Name("garbage") // dead nodes, before the live ones
	}
	script := IR_ScriptOfNodes(
		IR_Var(IR_Name("a")),
		IR_ExprResult(IR_Call(IR_Name("f"), IR_Name("a"), Node_NewNumber(1))),
		IR_Function(IR_Name("g"), IR_ParamList(IR_Name("x")), IR_BlockOfNode(IR_ReturnNodeWithExpr(IR_Name("x")))))
	root := IR_Root(script)
	expected := root.CloneTree()
	script.GetFirstChild().GetFirstChild().SetJSType(num)
	IR_Block().AddChildToBack(IR_Name("dead")) // a dead tree between live ones

	fn := script.GetLastChild()
	jsdoc := NewJSDocInfo_Builder()
	jsdoc.RecordType(NewJSTypeExpression(IR_Name("number"), "test.js"))
	fn.SetJSDocInfo(jsdoc.Build())
	typeRoot := fn.GetJSDocInfo().GetType().GetRoot()

	// A detached subtree held through an interior node (like a module binding).
	detached := IR_Var(IR_Name("d"))
	heldInterior := detached.GetFirstChild()
	// A detached node of its own.
	held := IR_Name("h")
	before := CurrentNodeStore().NodeCount()

	k := NewCompactor()
	holders := []Node{root, fn, heldInterior, held, expected}
	for _, n := range holders {
		k.Mark(n)
	}
	k.Slide()
	k.ForwardSlice(holders)
	newRoot, newFn, newInterior, newHeld, newExpected := holders[0], holders[1], holders[2], holders[3], holders[4]
	if k.Forward(newFn) != newFn {
		t.Fatal("Forward is not idempotent")
	}
	k.Finish()

	if live := CurrentNodeStore().NodeCount(); live >= before-1000 || live != k.Live {
		t.Fatalf("store has %d nodes after compaction (%d before, %d live)", live, before, k.Live)
	}
	if !newRoot.IsEquivalentTo(newExpected) {
		t.Fatalf("compacted tree differs:\n%s\nvs\n%s", newRoot.ToStringTree(), newExpected.ToStringTree())
	}
	if newFn.GetParent().GetParent() != newRoot || newFn.GetFirstChild().GetString() != "g" {
		t.Fatal("forwarded function is not in the compacted tree")
	}
	if newInterior.GetParent().IsNil() || !newInterior.GetParent().IsVar() || newInterior.GetString() != "d" {
		t.Fatal("interior holder was not forwarded with its detached subtree")
	}
	if newHeld.GetString() != "h" || !newHeld.GetParent().IsNil() {
		t.Fatal("detached holder was not forwarded")
	}
	if nt := newFn.GetJSDocInfo().GetType().GetRoot(); nt == typeRoot || nt.GetString() != "number" || nt.GetParent() != NilNode {
		t.Fatal("JSDoc type expression root was not forwarded")
	}
	if newRoot.GetFirstChild().GetFirstChild().GetFirstChild().GetJSType() != num {
		t.Fatal("the JSType of a node was lost")
	}
	if k.ExtraRoots != 5 { // root, expected, the two held subtrees, and the JSDoc type expression
		t.Fatalf("expected 5 extra roots, got %d", k.ExtraRoots)
	}
	// New nodes start zeroed after compaction.
	n := NewNode(Token_BLOCK)
	if !n.GetFirstChild().IsNil() || !n.GetParent().IsNil() || n.GetJSType() != nil || n.GetJSDocInfo() != nil {
		t.Fatal("a node allocated after compaction is not empty")
	}
	newRoot.GetFirstChild().AddChildToBack(n)
	// Stale use of an old node fails loudly.
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected a panic from a stale node")
			}
		}()
		fn.GetNext()
	}()
}
