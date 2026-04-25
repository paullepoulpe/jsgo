// Port of test/com/google/javascript/jscomp/ScopeTest.java @ closure-compiler v20260415
//
// ScopeTest.java was not available at the closure-compiler v20260415 tag
// (HTTP 404 on both the tag and master branch). The tests below exercise
// the Go TypedScope API surface that corresponds to Java's Scope / TypedScope
// contract: variable declaration, scope-chain lookup, refinement overlays,
// type-alias lookup, namespace-prop recording, and parent-scope traversal.

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// newTestScope creates a minimal TypedScope with no parent for unit tests.
func newTestScope() *TypedScope {
	reg := jstype.NewRegistry()
	return &TypedScope{
		declared: map[string]jstype.JSType{},
		registry: reg,
	}
}

// newChildScope creates a scope whose parent is parent.
func newChildScope(parent *TypedScope) *TypedScope {
	return &TypedScope{
		declared: map[string]jstype.JSType{},
		registry: parent.registry,
		parent:   parent,
	}
}

// ---------------------------------------------------------------------------
// Variable declaration and lookup
// ---------------------------------------------------------------------------

// TestGetDeclaredType_local: declaring a variable in a scope and looking it up
// in the same scope returns the declared type.
func TestGetDeclaredType_local(t *testing.T) {
	s := newTestScope()
	s.declared["x"] = jstype.Number
	got := s.GetDeclaredType("x")
	if got != jstype.Number {
		t.Errorf("GetDeclaredType(\"x\") = %v; want number", got)
	}
}

// TestGetDeclaredType_parentScope: a name declared in the parent scope is
// visible to the child scope through the scope chain.
func TestGetDeclaredType_parentScope(t *testing.T) {
	parent := newTestScope()
	parent.declared["y"] = jstype.String
	child := newChildScope(parent)
	got := child.GetDeclaredType("y")
	if got != jstype.String {
		t.Errorf("GetDeclaredType(\"y\") via parent = %v; want string", got)
	}
}

// TestGetDeclaredType_undeclared: looking up a name that was never declared
// returns nil from a scope with no parent.
func TestGetDeclaredType_undeclared(t *testing.T) {
	s := newTestScope()
	got := s.GetDeclaredType("undeclared")
	if got != nil {
		t.Errorf("GetDeclaredType of undeclared name = %v; want nil", got)
	}
}

// TestGetDeclaredType_childShadowsParent: a child declaration shadows the
// parent's declaration for the same name.
func TestGetDeclaredType_childShadowsParent(t *testing.T) {
	parent := newTestScope()
	parent.declared["x"] = jstype.String
	child := newChildScope(parent)
	child.declared["x"] = jstype.Number
	got := child.GetDeclaredType("x")
	if got != jstype.Number {
		t.Errorf("child should shadow parent: got %v; want number", got)
	}
}

// ---------------------------------------------------------------------------
// GetRawDeclaredType (skips refinements)
// ---------------------------------------------------------------------------

// TestGetRawDeclaredType_ignoresRefinement: GetRawDeclaredType should return
// the declared type, not the flow-refined type.
func TestGetRawDeclaredType_ignoresRefinement(t *testing.T) {
	s := newTestScope()
	s.declared["z"] = jstype.NewUnion(jstype.String, jstype.Null)
	s.PushRefinement(map[string]jstype.JSType{"z": jstype.String})
	raw := s.GetRawDeclaredType("z")
	if _, ok := raw.(*jstype.UnionType); !ok {
		t.Errorf("GetRawDeclaredType should bypass refinement; got %T %v", raw, raw)
	}
}

// ---------------------------------------------------------------------------
// Refinement overlay (PushRefinement / PopRefinement)
// ---------------------------------------------------------------------------

// TestPushPopRefinement: pushing a refinement makes the refined type visible;
// popping it restores the declared type.
func TestPushPopRefinement(t *testing.T) {
	s := newTestScope()
	s.declared["x"] = jstype.NewUnion(jstype.String, jstype.Null)

	s.PushRefinement(map[string]jstype.JSType{"x": jstype.String})
	if s.GetDeclaredType("x") != jstype.String {
		t.Error("after PushRefinement, GetDeclaredType should return refined type")
	}

	s.PopRefinement()
	if _, ok := s.GetDeclaredType("x").(*jstype.UnionType); !ok {
		t.Error("after PopRefinement, GetDeclaredType should return the original union")
	}
}

// TestWithoutRefinements: WithoutRefinements temporarily suspends overlays.
func TestWithoutRefinements(t *testing.T) {
	s := newTestScope()
	s.declared["x"] = jstype.Number
	s.PushRefinement(map[string]jstype.JSType{"x": jstype.String})

	var seenInside jstype.JSType
	s.WithoutRefinements(func() {
		seenInside = s.GetDeclaredType("x")
	})
	if seenInside != jstype.Number {
		t.Errorf("WithoutRefinements: inner lookup = %v; want number", seenInside)
	}
	// After the call the refinement must be restored.
	if s.GetDeclaredType("x") != jstype.String {
		t.Errorf("after WithoutRefinements, refinement should be restored; got %v", s.GetDeclaredType("x"))
	}
}

// TestGetQualifiedRefinement: a qualified name refinement is visible from child
// scopes.
func TestGetQualifiedRefinement(t *testing.T) {
	parent := newTestScope()
	parent.PushRefinement(map[string]jstype.JSType{"obj.prop": jstype.Number})
	child := newChildScope(parent)
	got := child.GetQualifiedRefinement("obj.prop")
	if got != jstype.Number {
		t.Errorf("GetQualifiedRefinement(\"obj.prop\") = %v; want number", got)
	}
}

// ---------------------------------------------------------------------------
// Type-alias namespace (LookupTypeAlias)
// ---------------------------------------------------------------------------

// TestLookupTypeAlias_local: a type alias declared in a scope is visible there.
func TestLookupTypeAlias_local(t *testing.T) {
	s := newTestScope()
	s.typeAliases = map[string]jstype.JSType{"MyNum": jstype.Number}
	got := s.LookupTypeAlias("MyNum")
	if got != jstype.Number {
		t.Errorf("LookupTypeAlias(\"MyNum\") = %v; want number", got)
	}
}

// TestLookupTypeAlias_parentScope: aliases declared in a parent are visible
// from a child scope.
func TestLookupTypeAlias_parentScope(t *testing.T) {
	parent := newTestScope()
	parent.typeAliases = map[string]jstype.JSType{"MyStr": jstype.String}
	child := newChildScope(parent)
	got := child.LookupTypeAlias("MyStr")
	if got != jstype.String {
		t.Errorf("LookupTypeAlias via parent = %v; want string", got)
	}
}

// TestLookupTypeAlias_missing: a name not registered as a typedef returns nil.
func TestLookupTypeAlias_missing(t *testing.T) {
	s := newTestScope()
	if s.LookupTypeAlias("Nonexistent") != nil {
		t.Error("LookupTypeAlias of unknown name should return nil")
	}
}

// ---------------------------------------------------------------------------
// Namespace properties (LookupNamespaceProp / recordNamespaceProp)
// ---------------------------------------------------------------------------

// TestLookupNamespaceProp_local: a property registered via recordNamespaceProp
// is retrievable by LookupNamespaceProp.
func TestLookupNamespaceProp_local(t *testing.T) {
	s := newTestScope()
	s.recordNamespaceProp("goog", "x", jstype.Boolean, nil)
	got := s.LookupNamespaceProp("goog", "x")
	if got != jstype.Boolean {
		t.Errorf("LookupNamespaceProp(\"goog\",\"x\") = %v; want boolean", got)
	}
}

// TestLookupNamespaceProp_missing: an unregistered property returns nil.
func TestLookupNamespaceProp_missing(t *testing.T) {
	s := newTestScope()
	if s.LookupNamespaceProp("ns", "missing") != nil {
		t.Error("LookupNamespaceProp of unregistered prop should be nil")
	}
}

// TestLookupNamespaceProp_parentScope: properties from the parent scope are
// visible in a child scope.
func TestLookupNamespaceProp_parentScope(t *testing.T) {
	parent := newTestScope()
	parent.recordNamespaceProp("goog", "bar", jstype.Number, nil)
	child := newChildScope(parent)
	got := child.LookupNamespaceProp("goog", "bar")
	if got != jstype.Number {
		t.Errorf("LookupNamespaceProp via parent = %v; want number", got)
	}
}

// ---------------------------------------------------------------------------
// Enum lookup (LookupEnum)
// ---------------------------------------------------------------------------

// TestLookupEnum_registered: a binding that holds an EnumType is found by
// LookupEnum.
func TestLookupEnum_registered(t *testing.T) {
	s := newTestScope()
	e := jstype.NewEnumType("Color", jstype.Number)
	s.declared["Color"] = e
	got := s.LookupEnum("Color")
	if got != e {
		t.Errorf("LookupEnum(\"Color\") = %v; want %v", got, e)
	}
}

// TestLookupEnum_notEnum: a binding that holds a non-enum type returns nil.
func TestLookupEnum_notEnum(t *testing.T) {
	s := newTestScope()
	s.declared["X"] = jstype.Number
	if s.LookupEnum("X") != nil {
		t.Error("LookupEnum of non-enum binding should return nil")
	}
}

// ---------------------------------------------------------------------------
// IsNamespacePropInferred
// ---------------------------------------------------------------------------

// TestIsNamespacePropInferred_flagSet: markNamespacePropInferred causes
// IsNamespacePropInferred to return true.
func TestIsNamespacePropInferred_flagSet(t *testing.T) {
	s := newTestScope()
	s.recordNamespaceProp("ns", "p", jstype.Number, nil)
	s.markNamespacePropInferred("ns", "p")
	if !s.IsNamespacePropInferred("ns", "p") {
		t.Error("IsNamespacePropInferred should be true after marking")
	}
}

// TestIsNamespacePropInferred_flagNotSet: without marking, the flag is false.
func TestIsNamespacePropInferred_flagNotSet(t *testing.T) {
	s := newTestScope()
	s.recordNamespaceProp("ns", "q", jstype.Number, nil)
	if s.IsNamespacePropInferred("ns", "q") {
		t.Error("IsNamespacePropInferred should be false when not marked")
	}
}

// ---------------------------------------------------------------------------
// ThisType
// ---------------------------------------------------------------------------

// TestThisType_set: a scope whose thisType is set returns that type.
func TestThisType_set(t *testing.T) {
	s := newTestScope()
	inst := jstype.NewObjectType("Foo")
	s.thisType = inst
	if s.ThisType() != inst {
		t.Errorf("ThisType() = %v; want %v", s.ThisType(), inst)
	}
}

// TestThisType_inheritsFromParent: if thisType is nil on the child, it falls
// back to the parent's thisType.
func TestThisType_inheritsFromParent(t *testing.T) {
	parent := newTestScope()
	inst := jstype.NewObjectType("Bar")
	parent.thisType = inst
	child := newChildScope(parent)
	if child.ThisType() != inst {
		t.Errorf("ThisType from parent = %v; want %v", child.ThisType(), inst)
	}
}

// TestThisType_nil: no thisType in the scope chain returns nil.
func TestThisType_nil(t *testing.T) {
	s := newTestScope()
	if s.ThisType() != nil {
		t.Error("ThisType should be nil when not set")
	}
}

// ---------------------------------------------------------------------------
// Registry accessor
// ---------------------------------------------------------------------------

// TestRegistry_notNil: Registry() returns the registry passed at construction.
func TestRegistry_notNil(t *testing.T) {
	s := newTestScope()
	if s.Registry() == nil {
		t.Error("Registry() should not be nil")
	}
}

// ---------------------------------------------------------------------------
// ResolveExpr — smoke test through the scope chain
// ---------------------------------------------------------------------------

// TestResolveExpr_primitive: resolving "number" returns Number.
func TestResolveExpr_primitive(t *testing.T) {
	s := newTestScope()
	got := s.ResolveExpr("number")
	if got != jstype.Number {
		t.Errorf("ResolveExpr(\"number\") = %v; want number", got)
	}
}

// TestResolveExpr_union: resolving "number|string" returns a union.
func TestResolveExpr_union(t *testing.T) {
	s := newTestScope()
	got := s.ResolveExpr("number|string")
	if got == nil {
		t.Fatal("ResolveExpr(\"number|string\") returned nil")
	}
	if _, ok := got.(*jstype.UnionType); !ok {
		t.Errorf("ResolveExpr(\"number|string\") = %T %v; want *UnionType", got, got)
	}
}

// TestResolveExpr_nominalViaScope: a name declared as a constructor in the
// scope resolves to its instance type (wrapped in a union with null because
// Closure legacy-mode nominal types are nullable by default, i.e. Foo = Foo|null).
func TestResolveExpr_nominalViaScope(t *testing.T) {
	s := newTestScope()
	inst := jstype.NewObjectType("MyClass")
	fn := jstype.NewFunctionType(nil, nil)
	fn.Instance = inst
	s.declared["MyClass"] = fn
	got := s.ResolveExpr("MyClass")
	if got == nil {
		t.Fatal("ResolveExpr(\"MyClass\") returned nil")
	}
	// In legacy-mode Closure, nominal types are nullable by default so
	// ResolveExpr("MyClass") produces MyClass|null. Accept either the raw
	// instance or a union containing the instance.
	switch v := got.(type) {
	case *jstype.ObjectType:
		if v != inst {
			t.Errorf("ResolveExpr(\"MyClass\") = %v; want instance %v", v, inst)
		}
	case *jstype.UnionType:
		found := false
		for _, m := range v.Members() {
			if m == inst {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ResolveExpr(\"MyClass\") union %v does not contain instance %v", v, inst)
		}
	default:
		t.Errorf("ResolveExpr(\"MyClass\") = %T %v; want ObjectType or Union", got, got)
	}
}
