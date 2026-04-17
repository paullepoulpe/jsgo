// Port of com/google/javascript/jscomp/TypedScopeCreator.java @ closure-compiler v20260415
//
// TypedScopeCreator walks the AST and builds a tree of TypedScopes
// (one per function plus the global scope) annotated with declared
// JSTypes. Tier-0 needed only the global scope's `var`/`let`/`const`
// bindings; Tier-1 adds top-level function declarations (with
// FunctionType signatures built from @param/@return) plus a per-function
// local scope holding the parameters.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// TypedScope maps identifier names to their declared JSType. The parent
// pointer mirrors the Java scope chain used during name resolution.
type TypedScope struct {
	parent   *TypedScope
	declared map[string]jstype.JSType
}

func (s *TypedScope) GetDeclaredType(name string) jstype.JSType {
	if s == nil {
		return nil
	}
	if t, ok := s.declared[name]; ok {
		return t
	}
	return s.parent.GetDeclaredType(name)
}

type TypedScopeCreator struct {
	registry *jstype.Registry
	// Functions associates each AST FUNCTION node with its local
	// (parameter) scope. TypeCheck reads this to type-check function
	// bodies with the right lexical scope active.
	Functions map[*rhino.Node]*TypedScope
}

func NewTypedScopeCreator(r *jstype.Registry) *TypedScopeCreator {
	return &TypedScopeCreator{registry: r, Functions: map[*rhino.Node]*TypedScope{}}
}

// CreateScope mirrors TypedScopeCreator.createScope for the SCRIPT root.
// Walks top-level declarations and records their types; recurses into
// function bodies to build their parameter scopes.
func (c *TypedScopeCreator) CreateScope(script *rhino.Node) *TypedScope {
	return c.createScopeIn(script, nil)
}

func (c *TypedScopeCreator) createScopeIn(block *rhino.Node, parent *TypedScope) *TypedScope {
	scope := &TypedScope{parent: parent, declared: make(map[string]jstype.JSType)}
	for _, child := range block.Children() {
		switch {
		case isDeclToken(child.Token()):
			c.recordDecl(child, scope)
		case child.Token() == rhino.FUNCTION:
			c.recordFunction(child, scope)
		}
	}
	return scope
}

func (c *TypedScopeCreator) recordDecl(v *rhino.Node, scope *TypedScope) {
	info := v.JSDocInfo()
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return
	}
	if !info.HasType() {
		return
	}
	scope.declared[name.String_()] = c.registry.ResolveName(info.Type().TypeName)
}

func (c *TypedScopeCreator) recordFunction(fn *rhino.Node, scope *TypedScope) {
	info := fn.JSDocInfo()
	kids := fn.Children()
	if len(kids) < 3 {
		return
	}
	name := kids[0]
	params := kids[1]
	body := kids[2]

	// Build a FunctionType from @param names (positional) and @return.
	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		var t jstype.JSType = jstype.Unknown
		if pexpr := info.ParameterType(p.String_()); pexpr != nil {
			t = c.registry.ResolveName(pexpr.TypeName)
		}
		paramTypes = append(paramTypes, t)
	}
	var retType jstype.JSType = jstype.Unknown
	if r := info.ReturnType(); r != nil {
		retType = c.registry.ResolveName(r.TypeName)
	}
	fnType := jstype.NewFunctionType(paramTypes, retType)
	if name.String_() != "" {
		scope.declared[name.String_()] = fnType
	}

	// Per-function local scope contains the parameters.
	localScope := &TypedScope{parent: scope, declared: make(map[string]jstype.JSType)}
	for i, p := range params.Children() {
		if i < len(paramTypes) {
			localScope.declared[p.String_()] = paramTypes[i]
		}
	}
	// Recurse into nested decls/functions so their local scopes are
	// registered too. Statements directly inside the body see
	// localScope as parent.
	nestedScope := c.createScopeIn(body, localScope)
	// Merge nested-local-only additions back into localScope so
	// TypeCheck sees them all.
	for k, v := range nestedScope.declared {
		if _, exists := localScope.declared[k]; !exists {
			localScope.declared[k] = v
		}
	}
	c.Functions[fn] = localScope
}

func isDeclToken(t rhino.Token) bool {
	return t == rhino.VAR || t == rhino.LET || t == rhino.CONST
}
