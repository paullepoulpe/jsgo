package jscomp

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// The optimized discovery must agree with the two original read-only walks,
// including order and declaration identity, rather than only the set of names.
func TestFlowSensitiveFunctionFactsDifferential(t *testing.T) {
	cases := []string{
		`function f(a,b) { var x=a; let y=b; return x+y; }`,
		`function f(a) { var x=a; { let y=x; { const z=y; use(z); } } return () => x; }`,
		`function f(a) { let x=a; function g(b=x) { let x=b; return () => a+x; } return g(); }`,
		`function f(a = () => q) { let x=1; let y=2; return () => () => x+y+a; }`,
		`function f(a) { try { use(a); } catch (e) { let x=e; return () => x+a; } }`,
		`function f(a) { for (let x of a) { let y=x; use(() => y); } switch(a) { case 1: let z=2; use(z); } }`,
		`function f(a) { let x=1; class C extends a { [use(x)] = x; static q=x; method(b=x) { return () => a+b+x; } static { let z=x; use(z); } } return C; }`,
		`function f(a) { const C = class Local { q = a; static q = a; }; return C; }`,
		`function f({a,b:[c]}, ...rest) { let [d,e]=rest; return () => a+c+d+e; }`,
		`function f(a) { let _exported=a; var dup=1; var dup=2; return dup; }`,
		`function f(a) { let x=1; { let x=2; use(() => x); } return () => x; }`,
		`function f(a) { return function named(b) { return named(a+b); }; }`,
	}
	for i, src := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			compiler, scope := functionFactsTestScope(t, src)
			functionFactsCompare(t, compiler, scope)
			// No result survives this discovery call: after a mutation the next
			// discovery observes a new lexical declaration without invalidation.
			scope.GetRootNode().GetLastChild().AddChildToBack(rhino.IR_Let(rhino.IR_Name("late")))
			all, _ := functionFactsCompare(t, compiler, scope)
			if all.GetAllVariables().Get("late") == nil {
				t.Fatal("discovery reused stale declarations after AST mutation")
			}
		})
	}
}

func functionFactsTestScope(t *testing.T, src string) (*Compiler, *Scope) {
	t.Helper()
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSTABLE)
	options.SetCodingConvention(NewGoogleCodingConvention())
	compiler.InitOptions(options)
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	function := compiler.ParseTestCode(src).RemoveFirstChild()
	if len(compiler.GetErrors()) != 0 {
		t.Fatalf("parse errors: %v", compiler.GetErrors())
	}
	script := rhino.NewNodeWithChild(rhino.Token_SCRIPT, function)
	script.SetInputId(rhino.NewInputId("test"))
	creator := NewSyntacticScopeCreator(compiler)
	return compiler, creator.CreateScopeOfScope(function, Scope_CreateGlobalScope(script))
}

func functionFactsCompare(t *testing.T, compiler *Compiler, scope *Scope) (*NodeUtil_AllVarsDeclaredInFunction, *jutil.ObjectLinkedHashSet[*Var]) {
	t.Helper()
	creator := NewSyntacticScopeCreator(compiler)
	want := NodeUtil_GetAllVarsDeclaredInFunction(compiler, creator, scope)
	wantEscaped := jutil.NewObjectLinkedHashSet[*Var]()
	DataFlowAnalysis_ComputeEscaped(scope, wantEscaped, compiler, creator, want.GetAllVariables())
	got, gotEscaped := flowSensitiveInlineVariables_collectFunctionFacts(compiler, NewSyntacticScopeCreator(compiler), scope)
	compare := func(label string, want, got []*Var) {
		t.Helper()
		if len(want) != len(got) {
			t.Fatalf("%s length = %d, want %d", label, len(got), len(want))
		}
		for i, v := range want {
			w := got[i]
			if !v.Equals(w) || v.GetNode() != w.GetNode() || v.GetScope().GetRootNode() != w.GetScope().GetRootNode() {
				t.Fatalf("%s[%d] = %v, want %v with identical declaration/scope root", label, i, w, v)
			}
		}
	}
	var wantMap, gotMap, wantEsc, gotEsc []*Var
	for v := range want.GetAllVariables().Values() {
		wantMap = append(wantMap, v)
	}
	for v := range got.GetAllVariables().Values() {
		gotMap = append(gotMap, v)
	}
	for v := range wantEscaped.All() {
		wantEsc = append(wantEsc, v)
	}
	for v := range gotEscaped.All() {
		gotEsc = append(gotEsc, v)
	}
	compare("map order", wantMap, gotMap)
	compare("declaration order", want.GetAllVariablesInOrder(), got.GetAllVariablesInOrder())
	compare("escape order", wantEsc, gotEsc)
	return got, gotEscaped
}

// Effectful custom scanners retain the original callback count and order.
func TestFlowSensitiveFunctionFactsCustomRedeclarations(t *testing.T) {
	compiler, scope := functionFactsTestScope(t, `function f(a) { var x=1; var x=2; return () => x+a; }`)
	wantHandler := &functionFactsRedeclarationHandler{}
	wantCreator := NewSyntacticScopeCreatorWithRedeclarationHandler(compiler, wantHandler)
	want := NodeUtil_GetAllVarsDeclaredInFunction(compiler, wantCreator, scope)
	DataFlowAnalysis_ComputeEscaped(scope, jutil.NewObjectLinkedHashSet[*Var](), compiler, wantCreator, want.GetAllVariables())
	gotHandler := &functionFactsRedeclarationHandler{}
	flowSensitiveInlineVariables_collectFunctionFacts(compiler, NewSyntacticScopeCreatorWithRedeclarationHandler(compiler, gotHandler), scope)
	if len(wantHandler.nodes) == 0 || len(wantHandler.nodes) != len(gotHandler.nodes) {
		t.Fatalf("redeclaration callback counts: got %d, want %d", len(gotHandler.nodes), len(wantHandler.nodes))
	}
	for i, want := range wantHandler.nodes {
		if gotHandler.nodes[i] != want {
			t.Fatalf("redeclaration callback %d changed node", i)
		}
	}
}

type functionFactsRedeclarationHandler struct{ nodes []*rhino.Node }

func (h *functionFactsRedeclarationHandler) OnRedeclaration(s *Scope, name string, n *rhino.Node, input *CompilerInput) {
	h.nodes = append(h.nodes, n)
}
