package jscomp

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestLiveVariableIndicesKeepDeclaredNamesAfterPublicVarReinitialization(t *testing.T) {
	for _, count := range []int{2, 17} {
		t.Run(fmt.Sprintf("variables_%d", count), func(t *testing.T) {
			compiler := NewCompiler()
			compiler.InitOptions(NewCompilerOptions())
			body := rhino.IR_Block()
			fn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), body)
			root := rhino.IR_Root(rhino.IR_ScriptOfNodes(fn))
			global := Scope_CreateGlobalScope(root)
			scope := Scope_CreateChildScope(global, fn)
			byName := jutil.NewLinkedHashMap[string, *Var]()
			ordered := make([]*Var, 0, count)
			for i := 0; i < count; i++ {
				name := fmt.Sprintf("v%d", i)
				n := rhino.IR_Name(name)
				body.AddChildToBack(rhino.IR_Var(n))
				v := scope.Declare(name, n, nil)
				byName.Put(name, v)
				ordered = append(ordered, v)
			}
			allVars := NewNodeUtil_AllVarsDeclaredInFunction(byName, ordered)
			analysis := NewLiveVariablesAnalysis(
				NewControlFlowGraph(fn, true, true), scope, nil, compiler,
				NewSyntacticScopeCreator(compiler), allVars)
			if analysis.GetVarIndex("v0") != 0 {
				t.Fatal("fixture did not capture the original variable index")
			}
			v := analysis.GetAllVariablesInOrder()[0]
			if analysis.GetAllVariables().Get("v0") != v {
				t.Fatal("public analysis views disagree before reinitialization")
			}
			v.InitAbstractVarBase("renamed", v.GetNameNode(), scope, 0, nil, nil)
			if v.GetName() != "renamed" || scope.GetVar("v0") != v {
				t.Fatal("public Var reinitialization did not preserve its declaration key")
			}
			if got := analysis.GetVarIndex("v0"); got != 0 {
				t.Fatalf("original declaration index = %d, want 0", got)
			}
			var panicValue any
			func() {
				defer func() { panicValue = recover() }()
				analysis.GetVarIndex("renamed")
			}()
			if _, ok := panicValue.(jutil.INullPointerException); !ok {
				t.Fatalf("new variable name became an analysis key: panic=%v", panicValue)
			}
		})
	}
}
