package rhino

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: the pre-order layout of compaction (Compactor.Permute).

// preorderNodes returns the nodes of the tree rooted at n in pre-order.
func preorderNodes(n Node) []Node {
	out := []Node{n}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		out = append(out, preorderNodes(c)...)
	}
	return out
}

// checkConsecutive fails unless the nodes have consecutive ids from first.
func checkConsecutive(t *testing.T, what string, nodes []Node, first uint32) {
	t.Helper()
	for j, n := range nodes {
		if uint32(n) != first+uint32(j) {
			t.Fatalf("%s: node %d (%s) has id %d, want %d", what, j, n.GetToken(), uint32(n), first+uint32(j))
		}
	}
}

func TestCompactorPermuteLaysOutTreesInPreorder(t *testing.T) {
	var layouts []string
	for _, threads := range []int{1, 4} {
		t.Run(fmt.Sprintf("threads=%d", threads), func(t *testing.T) {
			prev := SetCurrentNodeStore(NewNodeStore())
			defer SetCurrentNodeStore(prev)
			registry := NewJSTypeRegistry(nil)
			num := registry.GetNativeType(JSTypeNative_NUMBER_TYPE)

			for range 100 {
				IR_Name("garbage")
			}
			// The IR builders make children before their parents, so the
			// allocation order is not the pre-order.
			fn := IR_Function(IR_Name("g"), IR_ParamList(IR_Name("x")), IR_BlockOfNode(IR_ReturnNodeWithExpr(IR_Name("x"))))
			jsdoc := NewJSDocInfo_Builder()
			jsdoc.RecordType(NewJSTypeExpression(IR_Name("number"), "test.js"))
			fn.SetJSDocInfo(jsdoc.Build())
			typeRoot := fn.GetJSDocInfo().GetType().GetRoot()
			fn.SetIsGeneratorFunction(true) // a cold property
			call := IR_Call(IR_Name("f"), IR_Name("a"), Node_NewNumber(1))
			script := IR_ScriptOfNodes(IR_Var(IR_Name("a")), IR_ExprResult(call), fn)
			withProp := IR_Name("p")
			propTree := IR_Getprop(IR_Name("T"), "U")
			withProp.PutProp(Node_Prop_DECLARED_TYPE_EXPR, propTree)
			script2 := IR_ScriptOfNodes(IR_ExprResult(withProp))
			// The scripts in an order other than their allocation order
			// (as after dependency sorting).
			root := IR_Root(script2, script)
			IR_Block().AddChildToBack(IR_Name("dead"))
			// A statement added later, at the front: its ids are the
			// highest.
			script.AddChildToFront(IR_ExprResult(IR_Name("late")))
			call.GetFirstChild().SetJSType(num)
			call.GetSecondChild().SetOriginalName(jutil.StringPtr("orig_a"))
			expected := root.CloneTree()

			// A detached subtree held through an interior node, and a
			// detached node of its own.
			detached := IR_Var(IR_Name("d"))
			heldInterior := detached.GetFirstChild()
			held := IR_Name("h")

			k := NewCompactor()
			k.RecordPreorder()
			holders := []Node{root, heldInterior, held, fn, expected}
			for _, n := range holders {
				k.Mark(n)
			}
			k.Permute(threads)
			k.ForwardSlice(holders)
			newRoot, newInterior, newHeld, newFn, newExpected := holders[0], holders[1], holders[2], holders[3], holders[4]
			if k.Forward(newRoot) != newRoot {
				t.Fatal("Forward is not idempotent")
			}
			base := k.newBase
			k.Finish()

			if !newRoot.IsEquivalentTo(newExpected) {
				t.Fatalf("compacted tree differs:\n%s\nvs\n%s", newRoot.ToStringTree(), newExpected.ToStringTree())
			}
			// The detached trees, in the order their holders were marked,
			// then the program's ROOT, then its scripts in allocation
			// order, each in pre-order and followed by the trees hung off
			// its nodes.
			next := base
			lay := func(what string, n Node) {
				nodes := preorderNodes(n)
				checkConsecutive(t, what, nodes, next)
				next += uint32(len(nodes))
			}
			lay("detached subtree", newInterior.GetParent())
			lay("detached node", newHeld)
			lay("clone", newExpected)
			if uint32(newRoot) != next {
				t.Fatalf("root has id %d, want %d", uint32(newRoot), next)
			}
			next++
			newScript, newScript2 := newRoot.GetLastChild(), newRoot.GetFirstChild()
			lay("script", newScript)
			newType := newFn.GetJSDocInfo().GetType().GetRoot()
			if newType == typeRoot || newType.GetString() != "number" {
				t.Fatal("JSDoc type expression root was not forwarded")
			}
			lay("JSDoc type expression", newType)
			lay("script 2", newScript2)
			newWithProp := newScript2.GetFirstChild().GetFirstChild()
			newPropTree, ok := newWithProp.GetProp(Node_Prop_DECLARED_TYPE_EXPR).(Node)
			if !ok || newPropTree == propTree || newPropTree.GetString() != "U" {
				t.Fatalf("node-valued prop was not forwarded: %v", newPropTree)
			}
			lay("prop tree", newPropTree)
			if live := CurrentNodeStore().NodeCount(); live != k.Live || uint32(live) != next-base {
				t.Fatalf("store has %d nodes, %d live, %d laid out", live, k.Live, next-base)
			}
			if k.ExtraRoots != 6 {
				t.Fatalf("expected 6 extra roots, got %d", k.ExtraRoots)
			}

			newCall := newScript.GetChildAtIndex(2).GetFirstChild()
			if newCall.GetFirstChild().GetJSType() != num || newCall.GetSecondChild().GetJSType() != nil {
				t.Fatal("the JSType of a node was lost or moved")
			}
			if o := newCall.GetSecondChild().GetOriginalName(); o == nil || *o != "orig_a" {
				t.Fatalf("the original name of a node was lost: %v", o)
			}
			if !newFn.IsGeneratorFunction() || newFn.GetFirstChild().GetString() != "g" {
				t.Fatal("the cold property of a node was lost")
			}
			var sig []byte
			for _, n := range preorderNodes(newRoot) {
				sig = fmt.Appendf(sig, "%d:%s ", uint32(n)-base, n.GetToken())
			}
			layouts = append(layouts, string(sig))

			n := NewNode(Token_BLOCK)
			if uint32(n) != next || !n.GetFirstChild().IsNil() || !n.GetParent().IsNil() || n.GetJSType() != nil || n.GetJSDocInfo() != nil {
				t.Fatal("a node allocated after compaction is not empty or not next")
			}
			func() {
				defer func() {
					if recover() == nil {
						t.Fatal("expected a panic from a stale node")
					}
				}()
				fn.GetNext()
			}()
		})
	}
	if len(layouts) == 2 && layouts[0] != layouts[1] {
		t.Fatalf("threads change the layout:\n%s\n%s", layouts[0], layouts[1])
	}
}

func TestCompactorPermuteMultiBlockReusesOldBlocks(t *testing.T) {
	for _, tc := range []struct {
		order   string
		threads int
	}{{"local", 1}, {"local", 3}, {"local", 4}, {"global", 1}, {"global", 4}} {
		t.Run(fmt.Sprintf("%s/threads=%d", tc.order, tc.threads), func(t *testing.T) {
			threads := tc.threads
			prev := SetCurrentNodeStore(NewNodeStore())
			defer SetCurrentNodeStore(prev)
			registry := NewJSTypeRegistry(nil)
			num := registry.GetNativeType(JSTypeNative_NUMBER_TYPE)

			// Live names interleaved with dead ones that have a type and an
			// original name (stale column values for reused blocks), then
			// put in the script in an order other than their ids': reversed
			// in runs of 1000 ("local", so that the old blocks empty one
			// after the other and are reused), or scattered over the whole
			// range ("global": every old block is needed until the end).
			const total = 5*blockLen + 777
			names := make([]Node, 0, total)
			for idx := range total {
				if idx%4 == 0 {
					d := IR_Name("dead")
					d.SetJSType(num)
					d.SetOriginalName(jutil.StringPtr("dead_orig"))
				}
				n := IR_Name(fmt.Sprintf("v_%d", idx))
				if idx%3 == 0 {
					n.SetJSType(num)
				}
				if idx%5 == 0 {
					n.SetOriginalName(jutil.StringPtr(fmt.Sprintf("orig_%d", idx)))
				}
				names = append(names, n)
			}
			// The trees hung off nodes are laid out after the program: an old
			// block that holds one is reused only at the end.
			for idx := 0; idx < total; idx += 257 {
				names[idx].PutProp(Node_Prop_DECLARED_TYPE_EXPR, IR_Name(fmt.Sprintf("prop_%d", idx)))
			}
			script := IR_Script()
			root := IR_Root(script)
			var want []int
			for j := range total {
				idx := (j * 7919) % total // 7919 is prime, so this visits every index once
				if tc.order == "local" {
					idx = j/1000*1000 + 999 - j%1000
					if idx >= total {
						idx = j
					}
				}
				script.AddChildToBack(names[idx])
				want = append(want, idx)
			}

			k := NewCompactor()
			k.RecordPreorder()
			k.Mark(root)
			k.Permute(threads)
			newRoot := k.Forward(root)
			base := k.newBase
			k.Finish()

			nodes := preorderNodes(newRoot)
			checkConsecutive(t, "program", nodes, base)
			keep := (k.Live + blockLen - 1) / blockLen
			if k.NewBlocks+k.ReusedBlocks != keep || (tc.order == "local") != (k.ReusedBlocks > 0) {
				t.Fatalf("%d new and %d reused blocks for %d", k.NewBlocks, k.ReusedBlocks, keep)
			}
			if CurrentNodeStore().NodeCount() != k.Live || k.Live != len(nodes)+total/257+1 {
				t.Fatalf("store has %d nodes, %d live, %d in the program", CurrentNodeStore().NodeCount(), k.Live, len(nodes))
			}
			j := 0
			for c := newRoot.GetFirstChild().GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				idx := want[j]
				if c.GetString() != fmt.Sprintf("v_%d", idx) {
					t.Fatalf("child %d = %q, want v_%d", j, c.GetString(), idx)
				}
				if (idx%3 == 0) != (c.GetJSType() == num) {
					t.Fatalf("child %d (v_%d): JSType %v", j, idx, c.GetJSType())
				}
				o := c.GetOriginalName()
				if idx%5 == 0 {
					if o == nil || *o != fmt.Sprintf("orig_%d", idx) {
						t.Fatalf("child %d (v_%d): original name %v", j, idx, o)
					}
				} else if o != nil && *o != c.GetString() {
					t.Fatalf("child %d (v_%d): unexpected original name %q", j, idx, *o)
				}
				if idx%257 == 0 {
					p, ok := c.GetProp(Node_Prop_DECLARED_TYPE_EXPR).(Node)
					if !ok || p.GetString() != fmt.Sprintf("prop_%d", idx) || uint32(p) < base+uint32(len(nodes)) {
						t.Fatalf("child %d (v_%d): prop node %v", j, idx, p)
					}
				}
				j++
			}
			if j != total {
				t.Fatalf("got %d children, want %d", j, total)
			}
			// The tail of the last block is clear for new nodes.
			for range 3 {
				n := NewNode(Token_NAME)
				if n.GetJSType() != nil || n.GetOriginalName() != nil || !n.GetParent().IsNil() || n.GetToken() != Token_NAME {
					t.Fatal("a node allocated after compaction is not empty")
				}
			}
		})
	}
}

func TestCompactorPermuteLaysOutScriptsInAllocationOrder(t *testing.T) {
	for _, threads := range []int{1, 4} {
		t.Run(fmt.Sprintf("threads=%d", threads), func(t *testing.T) {
			prev := SetCurrentNodeStore(NewNodeStore())
			defer SetCurrentNodeStore(prev)

			// Eight scripts of about half a block each, parsed one after the
			// other (with dead nodes among them), then put in the tree in
			// another order, as dependency sorting does.
			const nScripts, perScript = 8, blockLen / 2
			scripts := make([]Node, nScripts)
			for j := range scripts {
				script := IR_Script()
				for idx := range perScript {
					if idx%5 == 0 {
						IR_Name("dead")
					}
					script.AddChildToBack(IR_ExprResult(IR_Name(fmt.Sprintf("s%d_%d", j, idx))))
				}
				scripts[j] = script
			}
			treeOrder := []int{3, 7, 1, 5, 0, 4, 2, 6}
			root := IR_Root()
			for _, j := range treeOrder {
				root.AddChildToBack(scripts[j])
			}

			k := NewCompactor()
			k.RecordPreorder()
			k.Mark(root)
			k.Permute(threads)
			newRoot := k.Forward(root)
			base := k.newBase
			k.Finish()

			if uint32(newRoot) != base {
				t.Fatalf("root has id %d, want %d", uint32(newRoot), base)
			}
			next := base + 1
			byAlloc := make([]Node, nScripts)
			for c, i := newRoot.GetFirstChild(), 0; !c.IsNil(); c, i = c.GetNext(), i+1 {
				byAlloc[treeOrder[i]] = c
				if want := fmt.Sprintf("s%d_0", treeOrder[i]); c.GetFirstChild().GetFirstChild().GetString() != want {
					t.Fatalf("script %d starts with %q, want %q", i, c.GetFirstChild().GetFirstChild().GetString(), want)
				}
			}
			for j, script := range byAlloc {
				nodes := preorderNodes(script)
				checkConsecutive(t, fmt.Sprintf("script %d", j), nodes, next)
				next += uint32(len(nodes))
			}
			if threads == 1 && k.NewBlocks > 1 {
				t.Fatalf("%d new and %d reused blocks: the copy in allocation order should reuse the old blocks", k.NewBlocks, k.ReusedBlocks)
			}
		})
	}
}
