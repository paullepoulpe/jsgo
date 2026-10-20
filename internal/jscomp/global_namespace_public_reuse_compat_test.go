package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func globalNamespacePublicReuseCompiler(t *testing.T) *Compiler {
	t.Helper()
	c := NewCompiler()
	options := NewCompilerOptions()
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetSkipNonTranspilationPasses(true)
	options.SetWrapGoogModulesForWhitespaceOnly(false)
	result := c.Compile(nil, []*SourceFile{
		SourceFile_FromCode("first.js", "var first = 1;"),
		SourceFile_FromCode("second.js", "var second = 2;"),
	}, options)
	if !result.Success || len(c.GetErrors()) != 0 {
		t.Fatal("valid namespace fixture did not compile")
	}
	return c
}

func TestGlobalNamespacePredicateCallbackTrace(t *testing.T) {
	c := globalNamespacePublicReuseCompiler(t)
	g := NewGlobalNamespaceWithRoots(c, nil, c.GetJsRoot())
	calls := make(map[*rhino.Node]int)
	var trace []*rhino.Node
	g.SetShouldTraverseScriptPredicate(func(script *rhino.Node) bool {
		calls[script]++
		trace = append(trace, script)
		return calls[script] == 1
	})
	g.GetNameForest()
	var want []*rhino.Node
	for script := c.GetJsRoot().GetFirstChild(); script != nil; script = script.GetNext() {
		want = append(want, script)
	}
	if len(trace) != len(want) {
		t.Fatalf("predicate trace has %d callbacks, want %d traversal callbacks", len(trace), len(want))
	}
	for i, script := range want {
		if trace[i] != script || calls[script] != 1 {
			t.Fatalf("predicate callback order/count changed at script %d", i)
		}
	}
	if g.GetOwnSlot("first") == nil || g.GetOwnSlot("second") == nil {
		t.Fatal("sizing consumed the stateful predicate before traversal")
	}
}

func TestGlobalNamespaceExposedNameIndexCaches(t *testing.T) {
	for _, requested := range []string{"fresh", "fresh.child"} {
		t.Run(requested, func(t *testing.T) {
			c := globalNamespacePublicReuseCompiler(t)
			g := NewGlobalNamespaceWithRoots(c, nil, c.GetJsRoot())
			index := g.GetNameIndex()
			builder := NewGlobalNamespace_BuildGlobalNamespace(g)
			builder.GetOrCreateName(requested, nil)
			old := index.Get("fresh")
			index.Remove("fresh")
			got := builder.GetOrCreateName("fresh", nil)
			if got == old || index.Get("fresh") != got {
				t.Fatal("cached name bypassed removal through the public name index")
			}
			replacement := NewGlobalNamespace_Name(nil, "fresh", nil, GlobalNamespace_SourceKind_CODE)
			index.Put("fresh", replacement)
			if builder.GetOrCreateName("fresh", nil) != replacement {
				t.Fatal("cached name bypassed replacement through the public name index")
			}
		})
	}
}

func globalNamespacePublicReuseRefs() (*GlobalNamespace_Name, []*rhino.Node, []*GlobalNamespace_Ref) {
	name := NewGlobalNamespace_Name(nil, "x", nil, GlobalNamespace_SourceKind_CODE)
	var nodes []*rhino.Node
	var refs []*GlobalNamespace_Ref
	for i := 0; i < 9; i++ {
		node := rhino.IR_Name("x")
		nodes = append(nodes, node)
		refs = append(refs, name.AddSingleRefForTesting(node, GlobalNamespace_Ref_Type_DIRECT_GET))
	}
	return name, nodes, refs
}

func TestGlobalNamespacePromotedRefsSnapshot(t *testing.T) {
	for _, beforePromotion := range []int{5, 8} {
		name := NewGlobalNamespace_Name(nil, "x", nil, GlobalNamespace_SourceKind_CODE)
		var first *GlobalNamespace_Ref
		for i := 0; i < beforePromotion; i++ {
			ref := name.AddSingleRefForTesting(rhino.IR_Name("x"), GlobalNamespace_Ref_Type_DIRECT_GET)
			if i == 0 {
				first = ref
			}
		}
		smallView := name.GetRefs()
		for i := beforePromotion; i < 9; i++ {
			name.AddSingleRefForTesting(rhino.IR_Name("x"), GlobalNamespace_Ref_Type_DIRECT_GET)
		}
		smallView[0] = nil
		if name.GetRefs()[0] != first {
			t.Fatalf("retained public slice of %d refs changed later promoted storage", beforePromotion)
		}
	}
	for _, keep := range []int{9, 1} {
		name, _, refs := globalNamespacePublicReuseRefs()
		for i := len(refs) - 1; i >= keep; i-- {
			name.RemoveRef(refs[i])
		}
		snapshot := name.GetRefs()
		if len(snapshot) != keep {
			t.Fatalf("after shrinking, got %d refs, want %d", len(snapshot), keep)
		}
		snapshot[0] = nil
		if name.GetRefs()[0] != refs[0] {
			t.Fatalf("public snapshot write changed promoted refs after shrinking to %d", keep)
		}
	}
}

func TestGlobalNamespacePromotedRefsMutableMap(t *testing.T) {
	name, nodes, refs := globalNamespacePublicReuseRefs()
	view := name.CastRefsForNodeMap()
	if view == nil {
		t.Fatal("promoted refs lost their public mutable map")
	}
	if name.CastRefsForNodeMap() != view {
		t.Fatal("repeated access returned a detached map")
	}
	view.Remove(nodes[0])
	if name.GetRefForNode(nodes[0]) != nil || len(name.GetRefs()) != 8 {
		t.Fatal("public map removal was not visible to Name")
	}
	replacement := NewGlobalNamespace_Ref(nil, nil, nodes[1], GlobalNamespace_Ref_Type_DIRECT_GET)
	view.Put(nodes[1], replacement)
	if name.GetRefForNode(nodes[1]) != replacement || name.GetFirstRef() != replacement {
		t.Fatal("public map replacement/order was not visible to Name")
	}
	newNode := rhino.IR_Name("x")
	name.UpdateRefNode(refs[2], newNode)
	if view.Get(nodes[2]) != nil || view.Get(newNode) != refs[2] {
		t.Fatal("Name update did not change the escaped map")
	}
	name.RemoveRefFromNodeMap(refs[3])
	if view.ContainsKey(nodes[3]) {
		t.Fatal("Name removal did not change the escaped map")
	}
	view.Clear()
	if len(name.GetRefs()) != 0 {
		t.Fatal("public map clear was not visible to Name")
	}
	name.PutRef(nodes[0], refs[0])
	if view.Get(nodes[0]) != refs[0] || name.GetFirstRef() != refs[0] {
		t.Fatal("Name insert after public map clear lost map ownership")
	}
}
