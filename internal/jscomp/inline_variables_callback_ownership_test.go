package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

func inlineExpertOwnershipFixture() (*InlineVariables_InliningBehavior, *Scope, *Var, *ReferenceCollection, *Var, *ReferenceCollection) {
	c := NewCompiler()
	c.InitOptions(NewCompilerOptions())
	first, second := rhino.IR_Name("first"), rhino.IR_Name("second")
	script := rhino.IR_ScriptOfNodes(rhino.IR_Var(first), rhino.IR_VarWithValue(second, rhino.IR_Number(1)))
	scope := Scope_CreateGlobalScope(rhino.IR_Root(script))
	input := NewCompilerInput(SourceFile_FromCode("experts.js", "var first; var second=1;"))
	a, b := scope.Declare("first", first, input), scope.Declare("second", second, input)
	ar, br := NewReferenceCollection(), NewReferenceCollection()
	ar.Add(newReferenceWithScopeAndInputId(first, nil, scope, input.GetInputId()))
	br.Add(newReferenceWithScopeAndInputId(second, nil, scope, input.GetInputId()))
	return NewInlineVariables_InliningBehavior(NewInlineVariables(c, InlineVariables_Mode_ALL)), scope, a, ar, b, br
}

func TestInlineVariablesPublicExpertOwnership(t *testing.T) {
	behavior, _, a, ar, b, br := inlineExpertOwnershipFixture()
	first := behavior.CreateVarExpert(a, ar, false).(*InlineVariables_InliningBehavior_StandardVarExpert)
	second := behavior.CreateVarExpert(b, br, false).(*InlineVariables_InliningBehavior_StandardVarExpert)
	if first == second || !first.IsNeverAssigned() || second.IsNeverAssigned() {
		t.Fatal("a later public expert replaced the earlier variable or its analysis")
	}
}
