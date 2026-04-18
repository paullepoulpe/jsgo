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
		// Strip generic brackets: `Array<string>` → `Array`. Proper
		// template resolution is a Tier-3 port; for now we at least
		// recover the nominal head so `@type {Array<T>}` resolves to
		// Array instead of collapsing the whole expression to Unknown.
		if i := indexOfByte(p, '<'); i >= 0 {
			p = trimSpace(p[:i])
		}
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

func indexOfByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
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
	return c.createScopeIn(script, c.builtins())
}

// builtins returns the synthetic parent scope carrying a minimal set
// of externs. Mirrors how Java's TypedScopeCreator pre-populates the
// scope from DefaultExterns.js — we can't load real externs yet, so
// each builtin is narrowly scoped to what jsgo can actually reason
// about without regressing on the long tail of tests that rely on
// the real (much richer) spec type.
func (c *TypedScopeCreator) builtins() *TypedScope {
	s := &TypedScope{
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
	}
	// Array: minimal surface so tests that reference `@type {Array}`
	// or touch `.length` / `.push` on an array variable don't
	// collapse to Unknown. We intentionally don't register Object /
	// Function / Number / Boolean / String stubs: Java's handling of
	// their boxed vs primitive vs object forms is subtle and any
	// narrow stub we ship would false-positive more than it fixes.
	arr := jstype.NewConstructorType("Array", nil)
	arr.Instance.DefineProperty("length", jstype.Number)
	arr.Instance.DefineProperty("push", jstype.NewFunctionType([]jstype.JSType{jstype.Unknown}, jstype.Number))
	arr.Instance.DefineProperty("pop", jstype.NewFunctionType(nil, jstype.Unknown))
	arr.Instance.DefineProperty("shift", jstype.NewFunctionType(nil, jstype.Unknown))
	arr.Instance.DefineProperty("indexOf", jstype.NewFunctionType([]jstype.JSType{jstype.Unknown}, jstype.Number))
	arr.Instance.DefineProperty("slice", jstype.NewFunctionType(nil, arr.Instance))
	arr.Instance.DefineProperty("concat", jstype.NewFunctionType([]jstype.JSType{jstype.Unknown}, arr.Instance))
	s.declared["Array"] = arr
	return s
}

func (c *TypedScopeCreator) createScopeIn(block *rhino.Node, parent *TypedScope) *TypedScope {
	scope := &TypedScope{
		parent:   parent,
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
	}
	// Pass 1: register every constructor-like binding (`@constructor`
	// functions + `class` declarations + `@record`/`@interface`
	// declarations) under its name so later `@type {Foo}` references
	// resolve nominally even when they appear textually before the
	// declaration.
	for _, child := range block.Children() {
		switch child.Token() {
		case rhino.FUNCTION:
			c.preregisterConstructor(child, scope)
		case rhino.CLASS:
			c.preregisterClass(child, scope)
		}
	}
	// Pass 2: record declarations + functions properly.
	for _, child := range block.Children() {
		switch {
		case isDeclToken(child.Token()):
			c.recordDecl(child, scope)
		case child.Token() == rhino.FUNCTION:
			c.recordFunction(child, scope)
		case child.Token() == rhino.CLASS:
			c.recordClass(child, scope)
		case child.Token() == rhino.EXPR_RESULT:
			// Closure-library-style prototype method attachment:
			//   Foo.prototype.bar = function(…) {…};
			// Register `bar` on Foo's InstanceType so later property
			// access + call-site checks can see it.
			if a := child.FirstChild(); a != nil && a.Token() == rhino.ASSIGN {
				c.recordPrototypeAssign(a, scope)
			}
		}
	}
	return scope
}

// recordPrototypeAssign detects the `Foo.prototype.<name> = <expr>`
// pattern and adds `<name>` to Foo's InstanceType member map, mirroring
// how TypedScopeCreator in Java inspects StaticSlot assignments that
// write to `X.prototype.*`.
func (c *TypedScopeCreator) recordPrototypeAssign(a *rhino.Node, scope *TypedScope) {
	kids := a.Children()
	if len(kids) < 2 {
		return
	}
	lhs, rhs := kids[0], kids[1]
	if lhs.Token() != rhino.GETPROP {
		return
	}
	lhsKids := lhs.Children()
	if len(lhsKids) < 2 {
		return
	}
	inner := lhsKids[0]
	propName := lhsKids[1].String_()
	if inner.Token() != rhino.GETPROP {
		return
	}
	innerKids := inner.Children()
	if len(innerKids) < 2 {
		return
	}
	if innerKids[0].Token() != rhino.NAME || innerKids[1].String_() != "prototype" {
		return
	}
	ownerName := innerKids[0].String_()
	fn, ok := scope.declared[ownerName].(*jstype.FunctionType)
	if !ok || fn.Instance == nil {
		return
	}
	memberType := c.inferPrototypeValueType(a.JSDocInfo(), rhs, scope)
	fn.Instance.DefineProperty(propName, memberType)
}

// inferPrototypeValueType computes the member type for a
// `X.prototype.m = value` assignment. Function expressions yield a
// FunctionType derived from the ASSIGN's @param/@return JSDoc; other
// values fall back to @type if present, otherwise Unknown.
func (c *TypedScopeCreator) inferPrototypeValueType(info *rhino.JSDocInfo, rhs *rhino.Node, scope *TypedScope) jstype.JSType {
	if rhs.Token() == rhino.FUNCTION {
		kids := rhs.Children()
		if len(kids) < 2 {
			return jstype.Unknown
		}
		params := kids[1]
		var paramTypes []jstype.JSType
		for _, p := range params.Children() {
			var t jstype.JSType = jstype.Unknown
			if info != nil {
				if expr := info.ParameterType(p.String_()); expr != nil {
					t = scope.ResolveExpr(expr.TypeName)
				}
			}
			paramTypes = append(paramTypes, t)
		}
		var ret jstype.JSType = jstype.Unknown
		if info != nil && info.ReturnType() != nil {
			ret = scope.ResolveExpr(info.ReturnType().TypeName)
		}
		return jstype.NewFunctionType(paramTypes, ret)
	}
	if info != nil && info.HasType() {
		return scope.ResolveExpr(info.Type().TypeName)
	}
	return jstype.Unknown
}

// preregisterConstructor reserves a slot in `scope` for any function
// carrying `@constructor`, `@interface`, or `@record`, so later
// `@type {Foo}` references resolve to the real InstanceType rather
// than collapsing to Unknown. The param/return types stay Unknown
// until pass 2 fills them in.
func (c *TypedScopeCreator) preregisterConstructor(fn *rhino.Node, scope *TypedScope) {
	info := fn.JSDocInfo()
	if !info.IsConstructor() && !info.IsInterface() && !info.IsRecord() {
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

// preregisterClass reserves a slot for a `class` declaration, mirroring
// the @constructor pre-registration path so forward `@type {C}`
// references resolve to C's InstanceType.
func (c *TypedScopeCreator) preregisterClass(cls *rhino.Node, scope *TypedScope) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	name := kids[0]
	if name.String_() == "" {
		return
	}
	scope.declared[name.String_()] = jstype.NewConstructorType(name.String_(), nil)
}

// recordClass fills in a class's constructor params, methods, fields,
// and super link. Mirrors how TypedScopeCreator processes a CLASS node
// in Java (TypedScopeCreator.java, visitClass).
func (c *TypedScopeCreator) recordClass(cls *rhino.Node, scope *TypedScope) {
	kids := cls.Children()
	if len(kids) < 3 {
		return
	}
	nameNode := kids[0]
	superNode := kids[1]
	body := kids[2]
	if nameNode.String_() == "" {
		return // anonymous classes — skip for Tier-2
	}
	preStub, ok := scope.declared[nameNode.String_()].(*jstype.FunctionType)
	if !ok || preStub.Instance == nil {
		return
	}
	inst := preStub.Instance
	// Mark as structural when the class carries @record or @interface.
	if info := cls.JSDocInfo(); info != nil && (info.IsRecord() || info.IsInterface()) {
		inst.IsStructural = true
	}

	// Locate the `constructor` method; otherwise the class has an
	// implicit zero-arg constructor.
	var ctor *rhino.Node
	for _, m := range body.Children() {
		if m.Token() == rhino.FUNCTION && len(m.Children()) >= 1 && m.FirstChild().String_() == "constructor" {
			ctor = m
			break
		}
	}
	var ctorParams []jstype.JSType
	if ctor != nil {
		ctorInfo := ctor.JSDocInfo()
		for _, p := range ctor.Children()[1].Children() {
			var t jstype.JSType = jstype.Unknown
			if ctorInfo != nil {
				if expr := ctorInfo.ParameterType(p.String_()); expr != nil {
					t = scope.ResolveExpr(expr.TypeName)
				}
			}
			ctorParams = append(ctorParams, t)
		}
	}
	preStub.Params = ctorParams

	// Wire @extends if present.
	if superNode.String_() != "" {
		if super := scope.LookupNominalType(superNode.String_()); super != nil {
			inst.SetSuper(super)
		}
	}

	// Walk body: methods attach to the InstanceType; fields attach with
	// their declared @type.
	for _, m := range body.Children() {
		switch m.Token() {
		case rhino.FUNCTION:
			if len(m.Children()) < 2 {
				continue
			}
			methodName := m.FirstChild().String_()
			if methodName == "constructor" {
				continue // handled above
			}
			info := m.JSDocInfo()
			params := m.Children()[1]
			var paramTypes []jstype.JSType
			for _, p := range params.Children() {
				var t jstype.JSType = jstype.Unknown
				if info != nil {
					if expr := info.ParameterType(p.String_()); expr != nil {
						t = scope.ResolveExpr(expr.TypeName)
					}
				}
				paramTypes = append(paramTypes, t)
			}
			var ret jstype.JSType = jstype.Unknown
			if info != nil && info.ReturnType() != nil {
				ret = scope.ResolveExpr(info.ReturnType().TypeName)
			}
			inst.DefineProperty(methodName, jstype.NewFunctionType(paramTypes, ret))
		case rhino.CLASS_FIELD_DEF:
			info := m.JSDocInfo()
			fieldName := m.FirstChild().String_()
			var t jstype.JSType = jstype.Unknown
			if info != nil && info.HasType() {
				t = scope.ResolveExpr(info.Type().TypeName)
			}
			inst.DefineProperty(fieldName, t)
		}
	}

	// Build local scopes for the constructor and each method so
	// TypeCheck can walk their bodies with `this` bound.
	c.Functions[cls] = &TypedScope{
		parent:   scope,
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
		thisType: inst,
	}
	if ctor != nil {
		c.Functions[ctor] = c.buildMethodScope(ctor, inst, scope, ctorParams)
		// Also collect `this.x = …` property writes inside the
		// constructor body.
		if ctorBody := ctor.Children()[2]; ctorBody != nil {
			c.collectThisProperties(ctorBody, inst)
		}
	}
	for _, m := range body.Children() {
		if m.Token() != rhino.FUNCTION || m.FirstChild().String_() == "constructor" {
			continue
		}
		params := m.Children()[1]
		var ptypes []jstype.JSType
		info := m.JSDocInfo()
		for _, p := range params.Children() {
			var t jstype.JSType = jstype.Unknown
			if info != nil {
				if expr := info.ParameterType(p.String_()); expr != nil {
					t = scope.ResolveExpr(expr.TypeName)
				}
			}
			ptypes = append(ptypes, t)
		}
		c.Functions[m] = c.buildMethodScope(m, inst, scope, ptypes)
	}
}

// buildMethodScope wraps the parameter + this binding needed for a
// method body to type-check under the same local-scope model as
// top-level @constructor functions.
func (c *TypedScopeCreator) buildMethodScope(fn *rhino.Node, inst *jstype.ObjectType, parent *TypedScope, ptypes []jstype.JSType) *TypedScope {
	local := &TypedScope{
		parent:   parent,
		declared: make(map[string]jstype.JSType),
		registry: c.registry,
		thisType: inst,
	}
	params := fn.Children()[1]
	for i, p := range params.Children() {
		if i < len(ptypes) {
			local.declared[p.String_()] = ptypes[i]
		}
	}
	return local
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
	// @record / @interface functions are also InstanceType carriers.
	if (info.IsConstructor() || info.IsInterface() || info.IsRecord()) && name.String_() != "" {
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
		// Mark structural when annotated @record / @interface.
		if info.IsInterface() || info.IsRecord() {
			fnType.Instance.IsStructural = true
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
