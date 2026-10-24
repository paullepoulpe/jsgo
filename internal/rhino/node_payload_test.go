package rhino

import (
	"math"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: the 32-bit payload of a node (node_payload.go).

func TestNumberPayload(t *testing.T) {
	// Number nodes are never negative (that is a NEG node), but the
	// encoding handles any value.
	for _, x := range []float64{-1, -(1 << 30), -(1 << 30) - 1, math.Copysign(0, -1), math.NaN()} {
		if got := payloadNumber(numberPayload(x)); math.Float64bits(got) != math.Float64bits(x) {
			t.Errorf("payload of %v reads back as %v", x, got)
		}
	}
	prev := SetCurrentNodeStore(NewNodeStore())
	defer SetCurrentNodeStore(prev)
	for _, x := range []float64{0, 1, 2, 42, 1<<30 - 1, 1 << 30, 1<<31 - 1, 0.5, 1e300, math.Inf(1), 4294967296, 3.14159} {
		n := Node_NewNumber(x)
		got := n.GetDouble()
		if math.Float64bits(got) != math.Float64bits(x) {
			t.Errorf("number %v reads back as %v", x, got)
		}
		if c := n.CloneNode(); math.Float64bits(c.GetDouble()) != math.Float64bits(x) || !c.IsEquivalentTo(n) {
			t.Errorf("clone of number %v is %v", x, c.GetDouble())
		}
	}
	if Node_NewNumber(0.5).IsEquivalentTo(Node_NewNumber(0.25)) {
		t.Error("different numbers are equivalent")
	}
}

func TestOriginalNameKeptAcrossRenaming(t *testing.T) {
	prev := SetCurrentNodeStore(NewNodeStore())
	defer SetCurrentNodeStore(prev)
	n := IR_Name("foo")
	n.SetOriginalNameFromName(n) // what SourceInformationAnnotator does
	if n.rec().kindFlags&kindFlagOrigNameIsString == 0 || origNameField.get(n) != 0 {
		t.Fatal("an original name equal to the node's string is not stored as the flag")
	}
	n.SetString("a") // renaming
	if o := n.GetOriginalName(); o == nil || *o != "foo" || n.GetString() != "a" {
		t.Fatalf("after renaming: string %q, original name %v", n.GetString(), o)
	}
	n.SetString("foo")
	if *n.GetOriginalName() != "foo" {
		t.Fatal("original name lost")
	}
	n.SetOriginalName(nil)
	if n.GetOriginalName() != nil || origNameField.get(n) != 0 || n.rec().kindFlags&kindFlagOrigNameIsString != 0 {
		t.Fatal("original name not cleared")
	}
	// Any kind of node can have an original name; copies keep it.
	call := IR_Call(IR_Name("f"))
	call.SetOriginalName(jutil.StringPtr("g"))
	clone := call.CloneNode()
	other := NewNode(Token_EMPTY).Srcref(call)
	if *clone.GetOriginalName() != "g" || *other.GetOriginalName() != "g" {
		t.Fatal("original name not copied")
	}

	// Compaction moves the original names with their nodes.
	root := IR_Root(IR_ScriptOfNodes(IR_ExprResult(call)))
	n2 := IR_Name("x")
	n2.SetOriginalNameFromName(n2)
	n2.SetString("y")
	held := []Node{root, n2}
	k := NewCompactor()
	for _, h := range held {
		k.Mark(h)
	}
	k.Slide()
	k.ForwardSlice(held)
	k.Finish()
	newCall := held[0].GetFirstChild().GetFirstChild().GetFirstChild()
	if *newCall.GetOriginalName() != "g" || *held[1].GetOriginalName() != "x" || held[1].GetString() != "y" {
		t.Fatal("compaction lost an original name")
	}
}

func TestTemplateLitCookedString(t *testing.T) {
	prev := SetCurrentNodeStore(NewNodeStore())
	defer SetCurrentNodeStore(prev)
	same := newNode_TemplateLiteralSubstringNode(jutil.StringPtr("abc"), "abc")
	escaped := newNode_TemplateLiteralSubstringNode(jutil.StringPtr("\n"), `\n`)
	invalid := newNode_TemplateLiteralSubstringNode(nil, `\u{`)
	if same.hasAux() {
		t.Error("a cooked string equal to the raw one takes an ext box")
	}
	if *same.GetCookedString() != "abc" || same.GetRawString() != "abc" {
		t.Error("cooked == raw read back wrong")
	}
	if *escaped.GetCookedString() != "\n" || escaped.GetRawString() != `\n` {
		t.Error("escaped cooked string read back wrong")
	}
	if invalid.GetCookedString() != nil || invalid.GetCookedStringID() != 0 || invalid.GetRawString() != `\u{` {
		t.Error("missing cooked string read back wrong")
	}
	for _, n := range []Node{same, escaped, invalid} {
		if c := n.CloneNode(); !c.IsEquivalentTo(n) || c.GetCookedStringID() != n.GetCookedStringID() {
			t.Errorf("clone of %q differs", n.GetRawString())
		}
	}
	if newNode_TemplateLiteralSubstringNode(jutil.StringPtr("x"), `\x78`).IsEquivalentTo(newNode_TemplateLiteralSubstringNode(jutil.StringPtr(`\x78`), `\x78`)) {
		t.Error("template strings with different cooked strings are equivalent")
	}
}
