package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type round8VarInfoConvention struct {
	CodingConvention
	calls    []string
	locals   []bool
	exported bool
}

func (c *round8VarInfoConvention) IsExportedWithLocal(name string, local bool) bool {
	c.calls = append(c.calls, name)
	c.locals = append(c.locals, local)
	return c.exported
}

func round8VarInfoOwner() (*RemoveUnusedCode, *round8VarInfoConvention) {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	r := NewRemoveUnusedCode(NewRemoveUnusedCode_Builder(compiler).
		RemoveGlobals(true).RemoveLocalVars(true))
	convention := &round8VarInfoConvention{CodingConvention: compiler.GetCodingConvention()}
	r.codingConvention = convention
	r.unremovableLog = diagnostic.LogFile_CreateNoOp()
	return r, convention
}

// A convention is an extensible callback, including on lookups whose private
// VarInfo storage is already cached. A false-to-true answer is observable.
func TestRemoveUnusedCodeConventionPrecedesVarInfoCaches(t *testing.T) {
	for _, equalCopy := range []bool{false, true} {
		label := "same_var"
		if equalCopy {
			label = "equal_distinct_var"
		}
		t.Run(label, func(t *testing.T) {
			r, convention := round8VarInfoOwner()
			name := rhino.IR_Name("subject")
			root := rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_Var(name)))
			scope := Scope_CreateGlobalScope(root)
			input := NewCompilerInput(SourceFile_FromCode("subject.js", "var subject;"))
			first := NewVar("subject", name, scope, 0, input, nil)
			second := first
			if equalCopy {
				second = NewVar("subject", name, scope, 0, input, nil)
				if second == first || !second.Equals(first) {
					t.Fatal("fixture does not supply equal distinct Vars")
				}
			}
			initial := r.GetVarInfo(first)
			if _, ok := initial.(*RemoveUnusedCode_RealVarInfo); !ok {
				t.Fatal("initial ordinary lookup did not create a real VarInfo")
			}
			convention.exported = true
			if got := r.GetVarInfo(second); got != r.canonicalUnremovableVarInfo {
				t.Fatal("private VarInfo cache bypassed updated convention classification")
			}
			convention.exported = false
			if got := r.GetVarInfo(second); got != initial {
				t.Fatal("restored ordinary classification lost equal-Var cache identity")
			}
			if len(convention.calls) != 3 {
				t.Fatalf("convention calls=%v, want three subject lookups", convention.calls)
			}
			for i, name := range convention.calls {
				if name != "subject" || convention.locals[i] {
					t.Fatalf("classification arguments changed: names=%v locals=%v", convention.calls, convention.locals)
				}
			}
		})
	}
}

func TestRemoveUnusedCodeConventionPrecedesArgumentsButNotExterns(t *testing.T) {
	t.Run("arguments", func(t *testing.T) {
		r, convention := round8VarInfoOwner()
		fn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), rhino.IR_Block())
		root := rhino.IR_Root(rhino.IR_ScriptOfNodes(fn))
		global := Scope_CreateGlobalScope(root)
		scope := Scope_CreateChildScope(global, fn)
		input := NewCompilerInput(SourceFile_FromCode("f.js", "function f() {}"))
		v := NewVar("arguments", nil, scope, -1, input, nil)
		if !v.IsArguments() || v.IsExtern() {
			t.Fatal("fixture does not supply a non-extern arguments Var")
		}
		if got := r.GetVarInfo(v); got != r.canonicalUnremovableVarInfo {
			t.Fatal("arguments should remain unremovable")
		}
		if len(convention.calls) != 1 || convention.calls[0] != "arguments" || !convention.locals[0] {
			t.Fatalf("arguments classification was bypassed: names=%v locals=%v", convention.calls, convention.locals)
		}
	})
	t.Run("extern", func(t *testing.T) {
		r, convention := round8VarInfoOwner()
		global := Scope_CreateGlobalScope(rhino.IR_Root())
		v := NewVar("externSubject", nil, global, 0, nil, nil)
		if !v.IsExtern() {
			t.Fatal("fixture does not supply an extern Var")
		}
		if got := r.GetVarInfo(v); got != r.canonicalUnremovableVarInfo {
			t.Fatal("extern should remain unremovable")
		}
		if len(convention.calls) != 0 {
			t.Fatalf("extern unexpectedly invoked convention: %v", convention.calls)
		}
	})
}
