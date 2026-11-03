package jscomp

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

func TestScopeAndVarSizes(t *testing.T) {
	if sz := reflect.TypeFor[Scope]().Size(); sz > 80 {
		t.Fatalf("Scope size = %d, want <= 80", sz)
	}
	if sz := reflect.TypeFor[TypedScope]().Size(); sz > 88 {
		t.Fatalf("TypedScope size = %d, want <= 88", sz)
	}
	if sz := reflect.TypeFor[Var]().Size(); sz > 32 {
		t.Fatalf("Var size = %d, want <= 32", sz)
	}
	tvType := reflect.TypeFor[TypedVar]()
	if sz := tvType.Size(); sz > 48 {
		t.Fatalf("TypedVar size = %d, want <= 48", sz)
	}
	if tvType.Field(0).Type.Kind() != reflect.Pointer || tvType.Field(1).Type.Kind() != reflect.Pointer || tvType.Field(3).Offset != 32 {
		t.Fatalf("expected TypedVar pointer-bearing fields to occupy bytes 0..32")
	}
}

func TestScopeLargeIndexAndQualifiedTypedScopeNoIntern(t *testing.T) {
	root := rhino.IR_Root()
	s := Scope_CreateGlobalScope(root)
	ts := TypedScope_CreateGlobalScope(root)
	for i := range 40 {
		name := "scopeVar_" + strconv.Itoa(i)
		s.Declare(name, rhino.NilNode, nil)
		ts.Declare(name, rhino.NilNode, nil, nil, false)
	}
	for i := range 40 {
		name := "scopeVar_" + strconv.Itoa(i)
		if s.GetVar(name) == nil || ts.GetVar(name) == nil {
			t.Fatalf("missing declared var %s", name)
		}
	}
	s.Undeclare(s.GetVar("scopeVar_17"))
	ts.Undeclare(ts.GetVar("scopeVar_17"))
	if s.GetVar("scopeVar_17") != nil || ts.GetVar("scopeVar_17") != nil {
		t.Fatal("undeclared var still found")
	}

	const transientQName = "scopeVar_0.uninterned_transient_prop_r75_w3"
	ts.Declare(transientQName, rhino.NilNode, nil, nil, true)
	if ts.GetVar(transientQName) == nil {
		t.Fatalf("expected TypedScope.GetVar(%q) to succeed", transientQName)
	}
	if id := strpool.Lookup(transientQName); id != 0 {
		t.Fatalf("transient qualified slot %q was interned into strpool (id=%d)", transientQName, id)
	}
}

func TestVarIsDefaultParam(t *testing.T) {
	compiler := NewCompiler()
	src := "var top = 1; var defExpr = 9; function f(a, b = 1, [c = 2] = [], d = defExpr) { return arguments; }"
	script := compiler.ParseTestCode(src)
	externs := rhino.IR_Root()
	jsRoot := rhino.IR_Root(script)
	root := rhino.IR_Root(externs, jsRoot)

	synCreator := NewSyntacticScopeCreator(compiler)
	globalSyn := synCreator.CreateScope(root, nil).(*Scope)
	fnNode := script.GetChildAtIndex(2)
	fnSyn := synCreator.CreateScope(fnNode, globalSyn).(*Scope)
	paramList := fnNode.GetSecondChild()
	bodySyn := synCreator.CreateScope(fnNode.GetLastChild(), fnSyn).(*Scope)

	// Implicit arguments / this / goog namespace vars (NilNode nameNode) must not SIGSEGV.
	if argsVar := fnSyn.GetArgumentsVar(); argsVar == nil || argsVar.IsDefaultParam() {
		t.Fatalf("Syntactic implicit arguments IsDefaultParam() = true or nil (%v), want false", argsVar)
	}
	if thisVar := fnSyn.GetVar("this"); thisVar == nil || thisVar.IsDefaultParam() {
		t.Fatalf("Syntactic implicit this IsDefaultParam() = true or nil (%v), want false", thisVar)
	}
	implicitNs := Var_CreateImplicitGoogNamespace("ns.foo", globalSyn, rhino.IR_String("ns.foo"))
	if implicitNs.IsDefaultParam() {
		t.Fatal("Implicit goog namespace var IsDefaultParam() = true, want false")
	}

	// Top-level var and detached / malformed NAME nodes must return false without panicking.
	if topVar := globalSyn.GetVar("top"); topVar == nil || topVar.IsDefaultParam() {
		t.Fatalf("Top-level var IsDefaultParam() = true or nil (%v), want false", topVar)
	}
	detachedVar := NewVar("detached", rhino.IR_Name("detached"), globalSyn, 99, nil, rhino.NilNode)
	if detachedVar.IsDefaultParam() {
		t.Fatal("Detached NAME var IsDefaultParam() = true, want false")
	}
	orphanDefault := paramList.GetSecondChild().CloneTree()
	orphanVar := NewVar("orphan", orphanDefault.GetFirstChild(), globalSyn, 100, nil, rhino.NilNode)
	if orphanVar.IsDefaultParam() {
		t.Fatal("DEFAULT_VALUE without PARAM_LIST parent IsDefaultParam() = true, want false")
	}
	rhsRefInDefault := NewVar("defExpr", paramList.GetLastChild().GetSecondChild(), globalSyn, 101, nil, rhino.NilNode)
	if rhsRefInDefault.IsDefaultParam() {
		t.Fatal("RHS child of DEFAULT_VALUE IsDefaultParam() = true, want false")
	}

	// Parameter vars in syntactic scope: only direct PARAM_LIST DEFAULT_VALUE targets return true.
	_ = bodySyn
	if aVar := fnSyn.GetVar("a"); aVar == nil || aVar.IsDefaultParam() {
		t.Fatalf("Plain param 'a' IsDefaultParam() = true or nil (%v), want false", aVar)
	}
	if bVar := fnSyn.GetVar("b"); bVar == nil || !bVar.IsDefaultParam() {
		t.Fatalf("Default param 'b = 1' IsDefaultParam() = false or nil (%v), want true", bVar)
	}
	if cVar := fnSyn.GetVar("c"); cVar == nil || cVar.IsDefaultParam() {
		t.Fatalf("Nested pattern default '[c = 2] = []' IsDefaultParam() = true or nil (%v), want false", cVar)
	}
	if dVar := fnSyn.GetVar("d"); dVar == nil || !dVar.IsDefaultParam() {
		t.Fatalf("Default param 'd = defExpr' IsDefaultParam() = false or nil (%v), want true", dVar)
	}

	// Also verify TypedScope / TypedVar implicit arguments, implicit this, detached, and params.
	tsGlobal := TypedScope_CreateGlobalScope(root)
	fnNode.SetJSType(compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_DATE_FUNCTION_TYPE))
	tsFn := NewTypedScope(tsGlobal, fnNode)
	if tvArgs := tsFn.GetArgumentsVar(); tvArgs == nil || tvArgs.IsDefaultParam() {
		t.Fatalf("TypedScope implicit arguments IsDefaultParam() = true or nil (%v), want false", tvArgs)
	}
	if tvThis := tsFn.GetVar("this"); tvThis == nil || tvThis.IsDefaultParam() {
		t.Fatalf("TypedScope implicit this IsDefaultParam() = true or nil (%v), want false", tvThis)
	}
	tvDetached := tsGlobal.Declare("tvDetached", rhino.NilNode, nil, nil, true)
	if tvDetached.IsDefaultParam() {
		t.Fatal("TypedVar with NilNode IsDefaultParam() = true, want false")
	}
	tvB := tsFn.Declare("b", paramList.GetSecondChild().GetFirstChild(), nil, nil, true)
	if !tvB.IsDefaultParam() {
		t.Fatal("TypedVar for 'b = 1' IsDefaultParam() = false, want true")
	}
}
