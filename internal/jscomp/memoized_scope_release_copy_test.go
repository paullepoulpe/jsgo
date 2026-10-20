package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReleaseArena is new in abcf69a0, so this root belongs to the new-main
// negative control and fixed candidate, not the historical 1fe overlay.
func TestMemoizedScopeCreatorValueCopiesShareRelease(t *testing.T) {
	delegate := &auditMemoizedScopeDelegate{}
	original := NewMemoizedScopeCreator(delegate)
	n := rhino.IR_Root()
	old := original.CreateScope(n, nil)
	copied := *original
	original.ReleaseArena()
	fresh := copied.CreateScope(n, nil)
	if fresh == old || original.CreateScope(n, nil) != fresh || delegate.calls != 2 {
		t.Fatal("value copies kept a stale cache after a shared release")
	}
}

func TestScopeCallbackEscapePreservesDifferentCreatorParents(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	parentCreator := NewSyntacticScopeCreatorWithRedeclarationHandler(compiler, NewSyntacticScopeCreator_DefaultRedeclarationHandler())
	parentRoot := rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_Var(rhino.IR_Name("outer"))))
	parent := parentCreator.CreateScopeOfScope(parentRoot, nil)
	outer := parent.GetVar("outer")
	h := &auditRetainingRedeclarationHandler{}
	childCreator := NewSyntacticScopeCreatorWithRedeclarationHandler(compiler, h)
	childRoot := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(rhino.IR_Name("x"), rhino.IR_Name("x")), rhino.IR_Block())
	child := childCreator.CreateScopeOfScope(childRoot, parent)
	childCreator.ReleaseArena()
	parentCreator.ReleaseArena()
	if h.scope != child || child.GetRootNode() != childRoot || child.GetParent() != parent || parent.GetRootNode() != parentRoot || parent.GetVar("outer") != outer || outer.GetName() != "outer" {
		t.Fatal("escaping child callback failed to protect a parent owned by another creator")
	}
}
