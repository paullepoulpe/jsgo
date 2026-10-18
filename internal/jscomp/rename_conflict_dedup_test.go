package jscomp

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

type conflictTestGraph struct {
	contexts []*MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext
	vars     []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo
	ids      map[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo]int
}

func newConflictTestGraph() *conflictTestGraph {
	g := &conflictTestGraph{ids: make(map[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo]int)}
	inverter := NewMakeDeclaredNamesUnique_ContextualRenameInverter(nil)
	for level := 0; level < 4; level++ {
		scope := &Scope{depth: level + 1}
		declared := make([]*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo, []int{3, 12, 24, 3}[level])
		for i := range declared {
			current := fmt.Sprintf("v%d$jscomp$%d", i, level)
			v := NewVar(current, nil, scope, i, nil, nil)
			info := NewMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(v, NewMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(inverter, current, fmt.Sprintf("v%d", i)))
			g.ids[info] = len(g.vars)
			g.vars = append(g.vars, info)
			declared[i] = info
		}
		var parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext
		if level > 0 {
			parent = g.contexts[(level-1)/2]
		}
		g.contexts = append(g.contexts, NewMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(scope, declared, parent))
	}
	return g
}

func originalRecordConflicts(s *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext, referenced MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	if s.scope.GetDepth() >= referenced.GetScopeDepth() {
		for _, declared := range s.declaredVariableInfos {
			referenced.AddPotentialConflict(declared)
			declared.AddPotentialConflict(referenced)
		}
		if s.parent != nil {
			originalRecordConflicts(s.parent, referenced)
		}
	}
}

func (g *conflictTestGraph) state() []any {
	var out []any
	for _, v := range g.vars {
		out = append(out, v.name)
		var conflicts []int
		if v.renamingInfo != nil {
			for item := range v.renamingInfo.potentialShadowVariables.All() {
				conflicts = append(conflicts, g.ids[item])
			}
		}
		out = append(out, conflicts)
	}
	return out
}

func TestRenameConflictDedupPreservesOrderAndNames(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		got, want := newConflictTestGraph(), newConflictTestGraph()
		rng := rand.New(rand.NewSource(seed))
		for step := 0; step < 300; step++ {
			scope, variable := rng.Intn(len(got.contexts)), rng.Intn(len(got.vars))
			if step%37 == 0 {
				got.vars[variable].TryToInvertName()
				want.vars[variable].TryToInvertName()
			} else {
				got.contexts[scope].RecordPotentialConflicts(got.vars[variable])
				originalRecordConflicts(want.contexts[scope], want.vars[variable])
			}
			if !reflect.DeepEqual(got.state(), want.state()) {
				t.Fatalf("seed %d step %d: conflict encounter order or names changed", seed, step)
			}
		}
	}
}
