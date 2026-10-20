package jscomp

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Common historical/current control; no private reuse fields or new APIs.
type crossChunkExportTraceConvention struct {
	CodingConvention
	trace             []string
	firstCalls        int
	alwaysExportFirst bool
}

func (c *crossChunkExportTraceConvention) IsExportedWithLocal(name string, local bool) bool {
	c.trace = append(c.trace, name)
	if name != "first" {
		return false
	}
	c.firstCalls++
	return c.alwaysExportFirst || c.firstCalls == 2
}

func TestCrossChunkCollectorExportPredicateEveryReference(t *testing.T) {
	for _, mode := range []string{"seeded_exported", "stateful", "seeded_stateful"} {
		t.Run(mode, func(t *testing.T) {
			c := globalNamespacePublicReuseCompiler(t)
			firstScript := c.GetJsRoot().GetFirstChild()
			firstScript.AddChildToBack(rhino.IR_ExprResult(rhino.IR_Name("first")))
			firstScript.AddChildToBack(rhino.IR_ExprResult(rhino.IR_Name("first")))
			convention := &crossChunkExportTraceConvention{CodingConvention: c.GetCodingConvention(), alwaysExportFirst: mode == "seeded_exported"}
			c.GetOptions().SetCodingConvention(convention)
			creator := NewMemoizedScopeCreator(NewSyntacticScopeCreator(c))
			scope := creator.CreateScope(c.GetJsRoot(), nil).(*Scope)
			v := scope.GetVar("first")
			collector := NewCrossChunkReferenceCollector(c, creator)
			seeded := mode != "stateful"
			if seeded {
				collector.AddReferenceToCollection(v, newReferenceWithScopeAndInputId(rhino.IR_Name("first"), nil, scope, nil))
			}
			collector.ProcessRoot(c.GetJsRoot())
			if !slices.Equal(convention.trace, []string{"first", "first", "first", "second"}) {
				t.Fatalf("export classification trace = %v; want one callback per global NAME", convention.trace)
			}
			count := 0
			for range collector.GetReferences(v).All() {
				count++
			}
			want := 2
			if seeded {
				want++
			}
			if mode == "seeded_exported" {
				want = 1
			}
			if count != want {
				t.Fatalf("first references = %d, want %d after per-reference export classification", count, want)
			}
			if mode == "seeded_exported" && collector.GetGlobalVariableNamesMap().ContainsKey("first") {
				t.Fatal("public seed caused exported name to enter traversal metadata")
			}
		})
	}
}
