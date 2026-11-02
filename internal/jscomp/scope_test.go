package jscomp

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
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
	if sz := reflect.TypeFor[TypedVar]().Size(); sz > 48 {
		t.Fatalf("TypedVar size = %d, want <= 48", sz)
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
	if id := rhino.RhinoStringPool_LookupID(transientQName); id != 0 {
		t.Fatalf("transient qualified slot %q was interned into strpool (id=%d)", transientQName, id)
	}
}
