package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func crossChunkPublicReuseVar(scope *Scope, name string) *Var {
	return NewVar(name, rhino.IR_Name(name), scope, 0, nil, nil)
}

func crossChunkPublicReuseReference(scope *Scope, name string) *Reference {
	return newReferenceWithScopeAndInputId(rhino.IR_Name(name), nil, scope, nil)
}

func TestCrossChunkCollectorPublicEqualVarAddition(t *testing.T) {
	c := NewCrossChunkReferenceCollector(NewCompiler(), nil)
	scope := Scope_CreateGlobalScope(rhino.IR_Root())
	v1 := crossChunkPublicReuseVar(scope, "x")
	v2 := crossChunkPublicReuseVar(scope, "x")
	if v1 == v2 || !v1.Equals(v2) {
		t.Fatal("fixture needs distinct equal Vars")
	}
	r1 := crossChunkPublicReuseReference(scope, "x")
	r2 := crossChunkPublicReuseReference(scope, "x")
	c.AddReferenceToCollection(v1, r1)
	c.AddReferenceToCollection(v2, r2)
	var got []*Reference
	for ref := range c.GetReferences(v1).All() {
		got = append(got, ref)
	}
	if len(got) != 2 || got[0] != r1 || got[1] != r2 || c.GetReferences(v2) != c.GetReferences(v1) {
		t.Fatal("equal Vars did not share one ordered reference collection")
	}
	var symbols []*Var
	for v := range c.GetAllSymbols() {
		symbols = append(symbols, v)
	}
	if len(symbols) != 1 || symbols[0] != v1 {
		t.Fatal("equal Var insertion changed original symbol identity/order")
	}
	// ObjectLinkedHashMap also accepts a nil key; the public helper did too.
	c.AddReferenceToCollection(nil, r1)
	if c.GetReferences(nil) == nil {
		t.Fatal("nil key no longer has a reference collection")
	}
}

func TestCrossChunkCollectorDeferredSymbolSequence(t *testing.T) {
	c := NewCrossChunkReferenceCollector(NewCompiler(), nil)
	scope := Scope_CreateGlobalScope(rhino.IR_Root())
	v1 := crossChunkPublicReuseVar(scope, "x")
	v2 := crossChunkPublicReuseVar(scope, "y")
	seq := c.GetAllSymbols()
	c.AddReferenceToCollection(v1, crossChunkPublicReuseReference(scope, "x"))
	var got []*Var
	for v := range seq {
		got = append(got, v)
		if v == v1 {
			c.AddReferenceToCollection(v2, crossChunkPublicReuseReference(scope, "y"))
		}
	}
	if len(got) != 2 || got[0] != v1 || got[1] != v2 {
		t.Fatal("delayed/live symbol enumeration missed later insertions")
	}
}

func TestCrossChunkCollectorPublicDraftOwnership(t *testing.T) {
	c := NewCrossChunkReferenceCollector(NewCompiler(), nil)
	n1 := rhino.NewNode(rhino.Token_EMPTY)
	n2 := rhino.NewNode(rhino.Token_EMPTY)
	d1 := c.InitializeDraftStatement(nil, n1)
	d2 := c.InitializeDraftStatement(nil, n2)
	if d1 == d2 {
		t.Fatal("public drafts alias reusable collector state")
	}
	s1 := NewCrossChunkReferenceCollector_TopLevelStatement(c, d1)
	s2 := NewCrossChunkReferenceCollector_TopLevelStatement(c, d2)
	if s1.GetStatementNode() != n1 || s1.GetOriginalOrder() != 0 || s2.GetStatementNode() != n2 || s2.GetOriginalOrder() != 1 {
		t.Fatal("later draft initialization changed retained draft contents")
	}
}
