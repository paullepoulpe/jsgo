package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"math/big"
	"runtime"
	"testing"
	"weak"
)

func TestPrivateNodeCloneMatchesOrdinaryClone(t *testing.T) {
	text := "cooked"
	for _, n := range []*Node{NewNode(Token_ROOT), Node_NewStringWithToken(Token_NAME, "x"), Node_NewNumber(3.5), Node_NewBigInt(big.NewInt(123)), Node_NewTemplateLitString(&text, "raw")} {
		n.SetOriginalName(jutil.StringPtr("original"))
		n.SetLinenoCharno(8, 4)
		n.SetLength(2)
		var slab jutil.Slab[Node]
		got, want := Node_CloneTreeInSlab(n, &slab), n.CloneTree()
		if !got.IsEquivalentTo(want) || got.String() != want.String() {
			t.Fatal(got, want)
		}
		got.SetOriginalName(jutil.StringPtr("changed"))
		if *n.GetOriginalName() != "original" {
			t.Fatal("mutable extension shared")
		}
	}
}

//go:noinline
func boundaryDiscardedTree() weak.Pointer[Node] {
	Node_AllocationBoundary()
	n := IR_Root()
	n.AddChildToBack(IR_Script())
	return weak.Make(n)
}

func TestNodeAllocationBoundaryDoesNotLinkCompilations(t *testing.T) {
	discarded := boundaryDiscardedTree()
	Node_AllocationBoundary()
	retained := NewNode(Token_ROOT)
	for i := 0; i < 10 && discarded.Value() != nil; i++ {
		runtime.GC()
	}
	if discarded.Value() != nil {
		t.Fatal("partial chunks linked independent compilations")
	}
	runtime.KeepAlive(retained)
}
