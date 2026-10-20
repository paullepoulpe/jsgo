package jscomp

import (
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Traversal must classify each reference even when an ordinary VarInfo already
// exists. A custom convention can change its answer between references.
func TestRemoveUnusedCodeTraversalReclassifiesCachedVars(t *testing.T) {
	for _, interleave := range []bool{false, true} {
		label := "consecutive_references"
		if interleave {
			label = "interleaved_references"
		}
		t.Run(label, func(t *testing.T) {
			r, convention := round8VarInfoOwner()
			subjectName := rhino.IR_Name("subject")
			otherName := rhino.IR_Name("other")
			root := rhino.IR_Root(rhino.IR_ScriptOfNodes(
				rhino.IR_Var(subjectName), rhino.IR_Var(otherName)))
			scope := Scope_CreateGlobalScope(root)
			input := NewCompilerInput(SourceFile_FromCode("subject.js", "var subject; var other;"))
			subject := NewVar("subject", subjectName, scope, 0, input, nil)
			other := NewVar("other", otherName, scope, 1, input, nil)

			initial := r.TraverseVar(subject)
			if _, ok := initial.(*RemoveUnusedCode_RealVarInfo); !ok {
				t.Fatal("initial ordinary reference did not create a real VarInfo")
			}
			wantCalls := []string{"subject"}
			if interleave {
				r.TraverseVar(other)
				wantCalls = append(wantCalls, "other")
			}
			convention.exported = true
			if got := r.TraverseVar(subject); got != r.canonicalUnremovableVarInfo {
				t.Fatal("cached traversal reference bypassed updated convention classification")
			}
			convention.exported = false
			if got := r.TraverseVar(subject); got != initial {
				t.Fatal("restored ordinary classification lost original VarInfo identity")
			}
			wantCalls = append(wantCalls, "subject", "subject")
			if !reflect.DeepEqual(convention.calls, wantCalls) {
				t.Fatalf("classification callbacks=%v, want %v", convention.calls, wantCalls)
			}
			for _, local := range convention.locals {
				if local {
					t.Fatal("global reference was classified as local")
				}
			}
		})
	}
}
