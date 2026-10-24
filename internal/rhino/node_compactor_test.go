package rhino

import (
	"testing"
)

// Go-only: copying compaction of the node store (node_compact.go).
func TestCompactorCopiesTreeAndForwardsHolders(t *testing.T) {
	prev := SetCurrentNodeStore(NewNodeStore())
	defer SetCurrentNodeStore(prev)
	script := IR_ScriptOfNodes(
		IR_Var(IR_Name("a")),
		IR_ExprResult(IR_Call(IR_Name("f"), IR_Name("a"), Node_NewNumber(1))),
		IR_Function(IR_Name("g"), IR_ParamList(IR_Name("x")), IR_BlockOfNode(IR_ReturnNodeWithExpr(IR_Name("x")))))
	root := IR_Root(script)
	expected := root.CloneTree()

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

	k := NewCompactor()
	newRoot := k.CopyTree(root)
	newFn := k.Forward(fn)
	newInterior := k.Forward(heldInterior)
	newHeld := k.Forward(held)
	expected = k.Forward(expected) // a holder of a whole detached tree
	if k.Forward(newFn) != newFn || k.Forward(fn) != newFn {
		t.Fatal("Forward is not idempotent")
	}
	if CurrentNodeStore() != k.Store() || newRoot.Store() != k.Store() || root.Store() == k.Store() {
		t.Fatal("the copy is not in the compactor's new, current store")
	}
	k.Finish()
	if !newRoot.IsEquivalentTo(expected) {
		t.Fatalf("copy differs:\n%s\nvs\n%s", newRoot.ToStringTree(), expected.ToStringTree())
	}
	if newFn.GetParent().GetParent() != newRoot || newFn.GetFirstChild().GetString() != "g" {
		t.Fatal("forwarded function is not in the copied tree")
	}
	if newInterior.GetParent().IsNil() || !newInterior.GetParent().IsVar() || newInterior.GetString() != "d" {
		t.Fatal("interior holder was not forwarded with its detached subtree")
	}
	if newHeld.GetString() != "h" || !newHeld.GetParent().IsNil() {
		t.Fatal("detached holder was not forwarded")
	}
	if nt := newFn.GetJSDocInfo().GetType().GetRoot(); nt == typeRoot || nt.GetString() != "number" {
		t.Fatal("JSDoc type expression root was not forwarded")
	}
	if k.ExtraRoots != 4 { // the two held subtrees, expected, and the JSDoc type expression
		t.Fatalf("expected 4 extra roots, got %d", k.ExtraRoots)
	}
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
