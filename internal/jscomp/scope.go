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
	// thisType is the instance bound to `this` inside this scope. nil
	// at the global / function scope; populated when entering a
	// constructor body so `this.x = …` lands on the right ObjectType.
	thisType *jstype.ObjectType
	registry *jstype.Registry
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

// ThisType returns the InstanceType `this` refers to in the lexically
// nearest constructor scope, or nil if `this` is unbound.
func (s *TypedScope) ThisType() *jstype.ObjectType {
	if s == nil {
		return nil
	}
	if s.thisType != nil {
		return s.thisType
	}
	return s.parent.ThisType()
}

// LookupNominalType returns the InstanceType associated with a name
// that refers to a known constructor in this scope chain. Used when
// resolving `@type {Foo}` to the instance type rather than the
// constructor function type.
func (s *TypedScope) LookupNominalType(name string) *jstype.ObjectType {
	if s == nil {
		return nil
	}
	if t, ok := s.declared[name]; ok {
		if fn, ok := t.(*jstype.FunctionType); ok && fn.Instance != nil {
			return fn.Instance
		}
	}
	return s.parent.LookupNominalType(name)
}

// ResolveExpr is the scope-aware version of Registry.ResolveExpr: it
// looks up nominal names (constructor → instance type) in the scope
// chain, falling back to the primitive registry. Mirrors how
// JSTypeRegistry.getType(StaticScope, String) dispatches in Java.
func (s *TypedScope) ResolveExpr(expr string) jstype.JSType {
	expr = trimSpace(expr)
	if expr == "" || expr == "?" {
		return jstype.Unknown
	}
	nullable := false
	if startsWith(expr, "?") {
		nullable = true
		expr = trimSpace(expr[1:])
	}
	parts := splitOn(expr, '|')
	var members []jstype.JSType
	for _, p := range parts {
		p = trimSpace(p)
		p = trimPrefix(p, "!")
		if p == "" {
			return jstype.Unknown
		}
		// Try nominal first (constructor → instance), then primitive.
		if inst := s.LookupNominalType(p); inst != nil {
			members = append(members, inst)
			continue
		}
		m := s.registry.ResolveName(p)
		if isJSUnknown(m) {
			return jstype.Unknown
		}
		members = append(members, m)
	}
	if nullable {
		members = append(members, jstype.Null)
	}
	if len(members) == 1 {
		return members[0]
	}
	return jstype.NewUnion(members...)
}

// Tiny string helpers — kept here to avoid pulling all of strings into
// every callsite, and to match the Java side's defensive scope code.

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func trimPrefix(s, prefix string) string {
	if startsWith(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func splitOn(s string, sep byte) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

func isJSUnknown(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	return t.TypeName() == "?"
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
	scope := &TypedScope{
		parent:   parent,
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
	}
	// Pass 1: register every constructor under its name so later
	// `@type {Foo}` references resolve nominally even when they
	// appear textually before `function Foo() {}`.
	for _, child := range block.Children() {
		if child.Token() == rhino.FUNCTION {
			c.preregisterConstructor(child, scope)
		}
	}
	// Pass 2: record declarations + functions properly.
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

// preregisterConstructor reserves a slot in `scope` for any function
// carrying `@constructor`, so forward `@type {Foo}` references resolve
// to the real InstanceType rather than collapsing to Unknown. The
// param/return types stay Unknown until pass 2 fills them in.
func (c *TypedScopeCreator) preregisterConstructor(fn *rhino.Node, scope *TypedScope) {
	info := fn.JSDocInfo()
	if !info.IsConstructor() {
		return
	}
	kids := fn.Children()
	if len(kids) < 3 {
		return
	}
	name := kids[0]
	if name.String_() == "" {
		return
	}
	// Stub: real param/return types are computed in recordFunction.
	scope.declared[name.String_()] = jstype.NewConstructorType(name.String_(), nil)
}

func (c *TypedScopeCreator) recordDecl(v *rhino.Node, scope *TypedScope) {
	info := v.JSDocInfo()
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return
	}
	if info.HasType() {
		scope.declared[name.String_()] = scope.ResolveExpr(info.Type().TypeName)
		return
	}
	// No @type: try to infer the binding's type from a simple
	// initializer. Mirrors TypedScopeCreator's "implicit declaration"
	// branch — only the shape `var x = new Foo();` is needed for
	// Tier-2.
	init := name.FirstChild()
	if init == nil || init.Token() != rhino.NEW {
		return
	}
	callee := init.FirstChild()
	if callee == nil || callee.Token() != rhino.NAME {
		return
	}
	if fn, ok := scope.declared[callee.String_()].(*jstype.FunctionType); ok && fn.Instance != nil {
		scope.declared[name.String_()] = fn.Instance
	}
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

	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		var t jstype.JSType = jstype.Unknown
		if pexpr := info.ParameterType(p.String_()); pexpr != nil {
			t = scope.ResolveExpr(pexpr.TypeName)
		}
		paramTypes = append(paramTypes, t)
	}
	var retType jstype.JSType = jstype.Unknown
	if r := info.ReturnType(); r != nil {
		retType = scope.ResolveExpr(r.TypeName)
	}

	var fnType *jstype.FunctionType
	if info.IsConstructor() && name.String_() != "" {
		// Pass 1 stub already registered the InstanceType under this
		// name; reuse the same Instance pointer so forward
		// `@type {Foo}` refs resolve to the same object as `new Foo()`.
		preStub, _ := scope.declared[name.String_()].(*jstype.FunctionType)
		fnType = &jstype.FunctionType{
			Params:        paramTypes,
			Return:        preStub.Instance,
			IsConstructor: true,
			Instance:      preStub.Instance,
		}
		// Wire @extends: look up the base ctor's Instance and set as
		// super on this ctor's Instance, mirroring
		// FunctionType.setExtendedClass in Java.
		if base := info.BaseType(); base != nil {
			if super := scope.LookupNominalType(base.TypeName); super != nil {
				fnType.Instance.SetSuper(super)
			}
		}
	} else {
		fnType = jstype.NewFunctionType(paramTypes, retType)
	}
	if name.String_() != "" {
		scope.declared[name.String_()] = fnType
	}

	// Per-function local scope contains the parameters and a binding
	// for `this` (only meaningful inside constructors).
	localScope := &TypedScope{
		parent:   scope,
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
		thisType: fnType.Instance, // nil for non-ctors
	}
	for i, p := range params.Children() {
		if i < len(paramTypes) {
			localScope.declared[p.String_()] = paramTypes[i]
		}
	}
	// Walk the body for `this.x = …` patterns so we can populate the
	// instance type's member map. Mirrors TypedScopeCreator's
	// AbstractScopeBuilder.attachLiteralTypes path.
	if fnType.Instance != nil {
		c.collectThisProperties(body, fnType.Instance)
	}
	nestedScope := c.createScopeIn(body, localScope)
	for k, v := range nestedScope.declared {
		if _, exists := localScope.declared[k]; !exists {
			localScope.declared[k] = v
		}
	}
	c.Functions[fn] = localScope
}

// collectThisProperties walks a constructor body looking for
// `this.<name> = <expr>` statements (with optional `@type` JSDoc) and
// adds each property to the instance's member map. Mirrors how
// TypedScopeCreator's AssignmentInferrer attaches inferred property
// types to `this`.
func (c *TypedScopeCreator) collectThisProperties(node *rhino.Node, inst *jstype.ObjectType) {
	if node == nil {
		return
	}
	if node.Token() == rhino.EXPR_RESULT {
		if a := node.FirstChild(); a != nil && a.Token() == rhino.ASSIGN {
			kids := a.Children()
			if len(kids) >= 2 && kids[0].Token() == rhino.GETPROP {
				obj := kids[0].FirstChild()
				if obj != nil && obj.Token() == rhino.THIS {
					prop := kids[0].Children()[1]
					var t jstype.JSType = jstype.Unknown
					if info := a.JSDocInfo(); info != nil && info.HasType() {
						t = inst.GetPropertyType(prop.String_()) // keep prior decl if any
						if t == nil {
							t = jstype.Unknown
						}
						declared := jstype.NewRegistry().ResolveExpr(info.Type().TypeName)
						if !isJSUnknown(declared) {
							t = declared
						}
					}
					inst.DefineProperty(prop.String_(), t)
				}
			}
		}
	}
	for _, k := range node.Children() {
		c.collectThisProperties(k, inst)
	}
}

func isDeclToken(t rhino.Token) bool {
	return t == rhino.VAR || t == rhino.LET || t == rhino.CONST
}
