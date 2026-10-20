package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type traversalScopeAPIWrapper struct {
	*Scope
	untyped func() *Scope
}

func (s *traversalScopeAPIWrapper) Untyped() *Scope { return s.untyped() }

// The slice makes this valid custom creator a noncomparable interface value.
type traversalScopeAPIValueCreator struct {
	marker []byte
	scope  AbstractScope
	calls  *int
}

func (c traversalScopeAPIValueCreator) CreateScope(_ *rhino.Node, _ AbstractScope) AbstractScope {
	*c.calls++
	return c.scope
}

func TestNodeTraversal_PublishedDefaultCreatorSurvivesRelease(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("function f() { var seed; }")
	var retained *SyntacticScopeCreator
	var arena *scopeArena
	var copied SyntacticScopeCreator
	cb := &nodeTraversalTestCallback{visit: func(nt *NodeTraversal, n, _ *rhino.Node) {
		if !n.IsName() || n.GetString() != "seed" {
			return
		}
		current := nt.GetScope()
		retained = nt.GetScopeCreator().(*SyntacticScopeCreator)
		if nt.GetScopeCreator().(*SyntacticScopeCreator) != retained || nt.GetScope() != current {
			t.Fatal("publishing changed creator or canonical scope identity")
		}
		arena = retained.arena
		if arena == nil || current.arena != arena {
			t.Fatal("fixture did not create the shared child-scope arena")
		}
		copied = *retained
	}}
	NodeTraversal_Traverse(fixture.compiler, root, cb)
	if retained == nil || retained.compiler != fixture.compiler || retained.arena != arena {
		t.Fatal("published creator lost its compiler or arena on traversal release")
	}
	if _, ok := retained.redeclarationHandler.(*SyntacticScopeCreator_DefaultRedeclarationHandler); !ok {
		t.Fatal("published creator lost its handler on traversal release")
	}
	// A subsequent public traversal may reuse the same pool entry with an
	// entirely different compiler. The published creator must remain usable.
	other := newSyntacticScopeCreatorTest(t)
	NodeTraversal_Traverse(other.compiler, other.getRoot("function other() { var next; }"), &nodeTraversalTestCallback{visit: func(nt *NodeTraversal, _ *rhino.Node, _ *rhino.Node) {
		nt.GetScope()
		nt.GetScopeCreator()
	}})
	if retained.compiler != fixture.compiler || retained.arena != arena {
		t.Fatal("another pooled traversal changed the published creator")
	}
	for _, creator := range []*SyntacticScopeCreator{retained, &copied} {
		global := creator.CreateScopeOfScope(root, nil)
		function := creator.CreateScopeOfScope(root.GetFirstChild(), global)
		body := creator.CreateScopeOfScope(root.GetFirstChild().GetLastChild(), function)
		if global.GetVar("f") == nil || body.GetVar("seed") == nil || body.arena != arena {
			t.Fatal("retained or value-copied creator cannot create valid scopes")
		}
	}
}

func TestNodeTraversal_PublishedDefaultCreatorSurvivesReentry(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("var outer;")
	var retained *SyntacticScopeCreator
	entries := 0
	cb := &nodeTraversalTestCallback{shouldTraverse: func(nt *NodeTraversal, _ *rhino.Node, _ *rhino.Node) bool {
		entries++
		creator := nt.GetScopeCreator().(*SyntacticScopeCreator)
		if retained == nil {
			retained = creator
			nt.Traverse(root)
		} else if creator != retained {
			t.Fatal("reentry changed the published creator")
		}
		return false
	}}
	NodeTraversal_Traverse(fixture.compiler, root, cb)
	if entries != 2 || retained == nil || retained.compiler != fixture.compiler {
		t.Fatal("reentrant traversal did not retain a usable creator")
	}
	if retained.CreateScopeOfScope(root, nil).GetVar("outer") == nil {
		t.Fatal("retained creator failed after reentrant traversal returned")
	}
}

func TestNodeTraversal_GetScopePreservesCustomUntyped(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("var x;")
	a := Scope_CreateGlobalScope(root)
	b := Scope_CreateGlobalScope(root)
	child := Scope_CreateGlobalScope(root)
	rootCalls, childCalls, creatorCalls := 0, 0, 0
	wrapper := &traversalScopeAPIWrapper{Scope: a}
	wrapper.untyped = func() *Scope {
		rootCalls++
		switch rootCalls {
		case 1, 4:
			return a
		case 2:
			return b
		default:
			return nil
		}
	}
	childWrapper := &traversalScopeAPIWrapper{Scope: child, untyped: func() *Scope {
		childCalls++
		return child
	}}
	creator := traversalScopeAPIValueCreator{marker: []byte{1}, scope: wrapper, calls: &creatorCalls}
	callback := &nodeTraversalTestCallback{shouldTraverse: func(nt *NodeTraversal, _ *rhino.Node, _ *rhino.Node) bool {
		for i, want := range []*Scope{a, b, nil} {
			if got := nt.GetScope(); got != want {
				t.Fatalf("custom Untyped call %d = %p, want %p", i+1, got, want)
			}
		}
		nt.PushScopeOfScope(childWrapper)
		if nt.GetScope() != child || nt.GetScope() != child || childCalls != 2 {
			t.Fatal("pushed custom scope's Untyped calls were cached")
		}
		nt.PopScope()
		if nt.GetScope() != a || rootCalls != 4 {
			t.Fatal("popping scope did not preserve custom Untyped behavior")
		}
		gotCreator, ok := nt.GetScopeCreator().(traversalScopeAPIValueCreator)
		if !ok || gotCreator.calls != &creatorCalls || len(gotCreator.marker) != 1 {
			t.Fatal("custom creator was replaced")
		}
		return false
	}}
	NodeTraversal_Builder_().SetCompiler(fixture.compiler).SetCallback(callback).SetScopeCreator(creator).Traverse(root)
	if creatorCalls != 1 {
		t.Fatalf("custom scope creator called %d times, want 1", creatorCalls)
	}
}

func TestNodeTraversal_GetScopeCustomUntypedReentry(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("var x;")
	a, b := Scope_CreateGlobalScope(root), Scope_CreateGlobalScope(root)
	var active *NodeTraversal
	calls, creatorCalls := 0, 0
	wrapper := &traversalScopeAPIWrapper{Scope: a}
	wrapper.untyped = func() *Scope {
		calls++
		if calls == 1 {
			if active.GetScope() != b {
				t.Fatal("reentrant Untyped result changed")
			}
			return a
		}
		return b
	}
	creator := traversalScopeAPIValueCreator{marker: []byte{1}, scope: wrapper, calls: &creatorCalls}
	callback := &nodeTraversalTestCallback{shouldTraverse: func(nt *NodeTraversal, _ *rhino.Node, _ *rhino.Node) bool {
		active = nt
		if nt.GetScope() != a || nt.GetScope() != b || calls != 3 {
			t.Fatal("custom scope reentry or subsequent Untyped call was cached")
		}
		return false
	}}
	NodeTraversal_Builder_().SetCompiler(fixture.compiler).SetCallback(callback).SetScopeCreator(creator).Traverse(root)
}

func TestNodeTraversal_GetScopeTypedScopeStillRejectsUntyped(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("var x;")
	typed := TypedScope_CreateGlobalScope(root)
	creatorCalls := 0
	creator := traversalScopeAPIValueCreator{marker: []byte{1}, scope: typed, calls: &creatorCalls}
	callback := &nodeTraversalTestCallback{shouldTraverse: func(nt *NodeTraversal, _ *rhino.Node, _ *rhino.Node) bool {
		for i := 0; i < 2; i++ {
			thrown := scopeTestCatch(func() { nt.GetScope() })
			if thrown == nil || thrown.GetMessage() != "untyped() called, but not an untyped scope." {
				t.Fatalf("GetScope TypedScope error = %v, want original untyped-scope rejection", thrown)
			}
		}
		return false
	}}
	NodeTraversal_Builder_().SetCompiler(fixture.compiler).SetCallback(callback).SetScopeCreator(creator).Traverse(root)
}

func TestNodeTraversal_ProvidedSyntacticCreatorIsUnchanged(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("var x;")
	provided := NewSyntacticScopeCreator(fixture.compiler)
	callback := &nodeTraversalTestCallback{shouldTraverse: func(nt *NodeTraversal, _ *rhino.Node, _ *rhino.Node) bool {
		got, ok := nt.GetScopeCreator().(*SyntacticScopeCreator)
		if !ok || got != provided {
			t.Fatal("GetScopeCreator replaced the provided creator")
		}
		return false
	}}
	NodeTraversal_Builder_().SetCompiler(fixture.compiler).SetCallback(callback).SetScopeCreator(provided).Traverse(root)
	if provided.compiler != fixture.compiler || provided.CreateScopeOfScope(root, nil).GetVar("x") == nil {
		t.Fatal("provided creator was changed by traversal release")
	}
}
