// Port of com/google/javascript/jscomp/TypedScopeCreator.java
// @ closure-compiler v20260415 (subset).
//
// Java's TypedScopeCreator is ~3000 lines spread across the outer
// class + AbstractScopeBuilder / ClassScopeBuilder /
// FunctionScopeBuilder inner classes, with callbacks that walk the
// AST in two passes and interact with the global registry. This file
// mirrors the overall shape (TypedScope tree with `this` binding +
// type-alias namespace; two-pass createScopeIn; pre-registration of
// constructor-like bindings) but collapses the inner-builder
// hierarchy into free functions on TypedScopeCreator
// (recordDecl / recordFunction / recordClass / recordTypedef /
// recordPrototypeAssign / collectThisProperties). Method names match
// Java where they exist (CreateScope, LookupNominalType,
// ResolveExpr). The inner-builder port lands when we need per-block
// hoisting or closure-captured scope nuances.

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
	// typeAliases carry `@typedef` declarations — a parallel namespace
	// to value bindings. ResolveExpr checks this before the nominal
	// / primitive registry.
	typeAliases map[string]jstype.JSType
	// thisType is the instance bound to `this` inside this scope. nil
	// at the global / function scope; populated when entering a
	// constructor body so `this.x = …` lands on the right ObjectType.
	thisType *jstype.ObjectType
	registry *jstype.Registry
	// refinements is a stack of overlays applied on top of `declared`
	// when GetDeclaredType is called. TypeCheck pushes a frame when
	// entering an if-true / if-false / hook-branch region so name
	// references inside the region see the narrowed type. The stack
	// is shared with parent scopes via the parent pointer (we walk
	// outward when looking up a name, applying every frame on the
	// way). Mirrors the FlowScope refinement shim Java's
	// TypeInference threads through its dataflow.
	refinements []map[string]jstype.JSType
}

// PushRefinement records a narrowing overlay for the duration of one
// block. Mirrors FlowScope's "refine and restore" pattern used by
// TypeInference when entering an if-true / if-false branch.
func (s *TypedScope) PushRefinement(refs map[string]jstype.JSType) {
	if s == nil || len(refs) == 0 {
		return
	}
	s.refinements = append(s.refinements, refs)
}

// PopRefinement removes the topmost overlay. Safe no-op if the stack
// is empty (matches the FlowScope.popInto contract).
func (s *TypedScope) PopRefinement() {
	if s == nil || len(s.refinements) == 0 {
		return
	}
	s.refinements = s.refinements[:len(s.refinements)-1]
}

func (s *TypedScope) GetDeclaredType(name string) jstype.JSType {
	if s == nil {
		return nil
	}
	// Refinement overlays take precedence — TypeCheck pushes them
	// when narrowing inside an if/hook branch.
	for i := len(s.refinements) - 1; i >= 0; i-- {
		if t, ok := s.refinements[i][name]; ok {
			return t
		}
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

// LookupTypeAlias walks the scope chain returning the resolved JSType
// for a @typedef'd name, or nil if none. Lives in a namespace parallel
// to `declared` so `var X;` + `@typedef {Y}` doesn't clash with a
// value binding named X.
func (s *TypedScope) LookupTypeAlias(name string) jstype.JSType {
	if s == nil {
		return nil
	}
	if t, ok := s.typeAliases[name]; ok {
		return t
	}
	return s.parent.LookupTypeAlias(name)
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
// looks up nominal names (constructor → instance type), @typedef
// aliases, inline structural `{k: T}` types, and primitives. Mirrors
// how JSTypeRegistry.getType(StaticScope, String) dispatches in Java.
func (s *TypedScope) ResolveExpr(expr string) jstype.JSType {
	expr = trimSpace(expr)
	if expr == "" || expr == "?" {
		return jstype.Unknown
	}
	// Strip balanced outer grouping parens. Closure's JsDocInfoParser
	// produces expressions like `(number|string)` when the union is
	// grouped with explicit parens — the Java TypeAST drops them
	// structurally; we do it textually here.
	for startsWith(expr, "(") && endsWith(expr, ")") && balancedOuter(expr) {
		expr = trimSpace(expr[1 : len(expr)-1])
		if expr == "" {
			return jstype.Unknown
		}
	}
	nullable := false
	if startsWith(expr, "?") {
		nullable = true
		expr = trimSpace(expr[1:])
	}
	// Inline structural: `{k1: T1, k2: T2}`. Covers the common
	// anonymous record shape that Closure tests use as @param /
	// @return types.
	if startsWith(expr, "{") {
		if obj, ok := s.resolveObjectLiteralType(expr); ok {
			if nullable {
				return jstype.NewUnion(obj, jstype.Null)
			}
			return obj
		}
		return jstype.Unknown
	}
	parts := splitTopLevelUnion(expr)
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
		// @typedef aliases take priority over nominal constructors
		// and primitives (matches Java's StaticSourceFile resolution).
		if alias := s.LookupTypeAlias(p); alias != nil {
			members = append(members, alias)
			continue
		}
		// Try nominal (constructor → instance), then primitive.
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

// splitTopLevelUnion splits `A|B|C` on `|` but respects nesting depth
// of `{}`, `()`, `<>` so `{a: number|string}|boolean` doesn't
// accidentally split inside the object literal.
func splitTopLevelUnion(s string) []string {
	var out []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{', '(', '<':
			depth++
		case '}', ')', '>':
			depth--
		case '|':
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	out = append(out, s[start:])
	return out
}

// resolveObjectLiteralType parses a `{k1: T1, k2: T2}` expression
// into a structural ObjectType (nil name, IsStructural=true) whose
// members are resolved against this scope. Returns false on malformed
// input so the caller can fall back to Unknown.
func (s *TypedScope) resolveObjectLiteralType(expr string) (*jstype.ObjectType, bool) {
	if !startsWith(expr, "{") || !endsWith(expr, "}") {
		return nil, false
	}
	body := expr[1 : len(expr)-1]
	obj := jstype.NewObjectType("")
	obj.IsStructural = true
	// Split on top-level commas; each piece is `key: expr`.
	depth := 0
	start := 0
	parts := []string{}
	for i := 0; i < len(body); i++ {
		switch body[i] {
		case '{', '(', '<':
			depth++
		case '}', ')', '>':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, body[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, body[start:])
	for _, p := range parts {
		p = trimSpace(p)
		if p == "" {
			continue
		}
		colon := indexOfTopLevel(p, ':')
		if colon < 0 {
			return nil, false
		}
		key := trimSpace(p[:colon])
		typ := s.ResolveExpr(p[colon+1:])
		obj.DefineProperty(key, typ)
	}
	return obj, true
}

func indexOfTopLevel(s string, b byte) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{', '(', '<':
			depth++
		case '}', ')', '>':
			depth--
		default:
			if depth == 0 && s[i] == b {
				return i
			}
		}
	}
	return -1
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// balancedOuter reports whether the outermost `(...)` in `s` is a
// complete paren group (the closing `)` is at `len(s)-1`). Tracks
// depth so `(a)|(b)` isn't misread as one big group.
func balancedOuter(s string) bool {
	if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
		return false
	}
	depth := 0
	for i := 0; i < len(s)-1; i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i < len(s)-1 {
				return false
			}
		}
	}
	return true
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
	global := c.createScopeIn(script, c.builtins())
	// Walk the full tree after the top-level pass so function
	// expressions / arrow functions nested inside initializers,
	// returns, and argument lists also get their local scopes. Matches
	// how TypedScopeCreator's AbstractScopeBuilder descends into every
	// syntactic child rather than the first-level only.
	c.registerNestedFunctions(script, global)
	return global
}

// registerNestedFunctions recursively visits every FUNCTION node that
// isn't already associated with a local scope (top-level functions
// were handled by recordFunction) and builds a parameter scope for
// it. Used so arrow / function-expression parameters are reachable to
// VarCheck and TypeCheck.
func (c *TypedScopeCreator) registerNestedFunctions(n *rhino.Node, enclosing *TypedScope) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		if _, seen := c.Functions[n]; !seen {
			c.buildExpressionFunctionScope(n, enclosing)
		}
	}
	next := enclosing
	if local, ok := c.Functions[n]; ok && local != nil {
		next = local
	}
	for _, child := range n.Children() {
		c.registerNestedFunctions(child, next)
	}
}

// buildExpressionFunctionScope creates a local scope for a FUNCTION
// node that was discovered outside the top-level statement walk
// (arrow bodies, anonymous function expressions, object-literal
// method shorthands). Parameters are bound to @param types when
// JSDoc is present, else Unknown. The function body is then walked
// via createScopeIn so nested declarations — `var foo = …`,
// `function bar(){}`, `class Baz{}` — populate the same local slot,
// matching how recordFunction handles top-level functions.
func (c *TypedScopeCreator) buildExpressionFunctionScope(fn *rhino.Node, parent *TypedScope) {
	kids := fn.Children()
	if len(kids) < 2 {
		return
	}
	info := fn.JSDocInfo()
	params := kids[1]
	local := &TypedScope{
		parent:      parent,
		declared:    make(map[string]jstype.JSType),
		typeAliases: make(map[string]jstype.JSType),
		registry:    c.registry,
	}
	for _, p := range params.Children() {
		var t jstype.JSType = jstype.Unknown
		if info != nil {
			if expr := info.ParameterType(p.String_()); expr != nil {
				t = parent.ResolveExpr(expr.TypeName)
			}
		}
		local.declared[p.String_()] = t
	}
	c.Functions[fn] = local
	if len(kids) >= 3 {
		body := kids[2]
		nested := c.createScopeIn(body, local)
		for k, v := range nested.declared {
			if _, exists := local.declared[k]; !exists {
				local.declared[k] = v
			}
		}
		for k, v := range nested.typeAliases {
			if _, exists := local.typeAliases[k]; !exists {
				local.typeAliases[k] = v
			}
		}
	}
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
		parent:      parent,
		declared:    make(map[string]jstype.JSType),
		typeAliases: make(map[string]jstype.JSType),
		registry:    c.registry,
	}
	hoisted := hoistDeclarations(block)
	// Pass 1: register every constructor-like binding (`@constructor`
	// functions + `class` declarations + `@record`/`@interface`
	// declarations) under its name so later `@type {Foo}` references
	// resolve nominally even when they appear textually before the
	// declaration.
	for _, child := range hoisted {
		switch child.Token() {
		case rhino.FUNCTION:
			c.preregisterConstructor(child, scope)
		case rhino.CLASS:
			c.preregisterClass(child, scope)
		}
	}
	// Pass 1b: record @typedef aliases. Typedefs resolve against the
	// scope that now has nominal types pre-registered, so
	// `@typedef {Foo|null}` can reference a Foo declared later.
	for _, child := range hoisted {
		if isDeclToken(child.Token()) {
			c.recordTypedef(child, scope)
		}
	}
	// Pass 2: record declarations + functions properly.
	for _, child := range hoisted {
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

// hoistDeclarations walks `block` collecting every VAR/LET/CONST /
// FUNCTION / CLASS / EXPR_RESULT it and its nested BLOCKs / IF
// branches / loop bodies contain. FUNCTION / CLASS subtrees are not
// descended into — they own their own scope, built separately via
// recordFunction / recordClass. Matches Java's hoistVars pass on
// AbstractScopeBuilder.
func hoistDeclarations(block *rhino.Node) []*rhino.Node {
	var out []*rhino.Node
	var walk func(n *rhino.Node)
	walk = func(n *rhino.Node) {
		if n == nil {
			return
		}
		for _, child := range n.Children() {
			switch child.Token() {
			case rhino.VAR, rhino.LET, rhino.CONST, rhino.FUNCTION, rhino.CLASS, rhino.EXPR_RESULT:
				out = append(out, child)
				// Keep descending into non-function forms so e.g. an
				// IF branch `if (c) { var x = …; }` contributes x.
				if child.Token() == rhino.FUNCTION || child.Token() == rhino.CLASS {
					continue
				}
				walk(child)
			case rhino.BLOCK, rhino.IF, rhino.HOOK:
				walk(child)
			}
		}
	}
	walk(block)
	return out
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

	// Walk body: methods attach to the InstanceType (or the static
	// slot when tagged `static`); fields attach with their declared
	// @type to the same slot.
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
			fnT := jstype.NewFunctionType(paramTypes, ret)
			if info.IsStatic() {
				inst.DefineStaticProperty(methodName, fnT)
			} else {
				inst.DefineProperty(methodName, fnT)
			}
		case rhino.CLASS_FIELD_DEF:
			info := m.JSDocInfo()
			fieldName := m.FirstChild().String_()
			var t jstype.JSType = jstype.Unknown
			if info != nil && info.HasType() {
				t = scope.ResolveExpr(info.Type().TypeName)
			}
			if info.IsStatic() {
				inst.DefineStaticProperty(fieldName, t)
			} else {
				inst.DefineProperty(fieldName, t)
			}
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
		ctorScope := c.buildMethodScope(ctor, inst, scope, ctorParams)
		// Bind `super` inside the derived-class ctor to the base's
		// FunctionType so `super(args)` + `super.method(args)` flow
		// through the normal checkCall / checkGetProp paths.
		c.bindSuper(ctorScope, superNode.String_(), scope)
		c.Functions[ctor] = ctorScope
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
		parent:      parent,
		declared:    make(map[string]jstype.JSType),
		typeAliases: make(map[string]jstype.JSType),
		registry:    c.registry,
		thisType:    inst,
	}
	params := fn.Children()[1]
	for i, p := range params.Children() {
		if i < len(ptypes) {
			local.declared[p.String_()] = ptypes[i]
		}
	}
	return local
}

// templateResolver returns (resolver, templateMap) — the resolver is
// a closure that resolves a type expression against a shim scope
// where each of the function's @template names is bound to a fresh
// TemplateType. The templateMap lets the caller stash those
// TemplateTypes on the resulting FunctionType so call-site
// unification can refer to them later. Port of the per-function
// TemplateTypeMap construction in TypedScopeCreator.java.
func (c *TypedScopeCreator) templateResolver(parent *TypedScope, info *rhino.JSDocInfo) (func(string) jstype.JSType, map[string]*jstype.TemplateType) {
	if info == nil || len(info.TemplateNames()) == 0 {
		return parent.ResolveExpr, nil
	}
	tmpls := make(map[string]*jstype.TemplateType, len(info.TemplateNames()))
	aliases := make(map[string]jstype.JSType, len(info.TemplateNames()))
	for _, tn := range info.TemplateNames() {
		tt := jstype.NewTemplateType(tn)
		tmpls[tn] = tt
		aliases[tn] = tt
	}
	shim := &TypedScope{
		parent:      parent,
		declared:    map[string]jstype.JSType{},
		typeAliases: aliases,
		registry:    c.registry,
	}
	return shim.ResolveExpr, tmpls
}

// bindSuper resolves the named superclass into its constructor's
// FunctionType and stores it under `super` in the given (ctor-local)
// scope. Harmless no-op when the class has no base.
func (c *TypedScopeCreator) bindSuper(local *TypedScope, superName string, outer *TypedScope) {
	if superName == "" {
		return
	}
	baseCtor, ok := outer.declared[superName].(*jstype.FunctionType)
	if !ok {
		return
	}
	local.declared["super"] = baseCtor
}

// recordTypedef stores a `@typedef {T}`-annotated `var NAME;` under
// the scope's typeAliases map. Mirrors how Java's TypedScopeCreator
// registers typedefs ahead of the normal scope-building passes so
// `@type {NAME}` lookups resolve to the aliased type.
func (c *TypedScopeCreator) recordTypedef(v *rhino.Node, scope *TypedScope) {
	info := v.JSDocInfo()
	if info == nil || info.TypedefType() == nil {
		return
	}
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return
	}
	scope.typeAliases[name.String_()] = scope.ResolveExpr(info.TypedefType().TypeName)
}

func (c *TypedScopeCreator) recordDecl(v *rhino.Node, scope *TypedScope) {
	info := v.JSDocInfo()
	name := v.FirstChild()
	if name == nil || name.Token() != rhino.NAME {
		return
	}
	// Typedef declarations primarily populate the typeAliases
	// namespace, but we still record the NAME at Unknown in the
	// value namespace so later `var alias = original;` references
	// don't fire UNDEFINED_VARIABLE. Java's VarCheck treats
	// `var original` (with @typedef) as a normal declaration slot.
	if info != nil && info.TypedefType() != nil {
		scope.declared[name.String_()] = jstype.Unknown
		return
	}
	// `@type {T}` pins the binding's declared type precisely.
	if info.HasType() {
		scope.declared[name.String_()] = scope.ResolveExpr(info.Type().TypeName)
		return
	}
	// No @type: try to recover a sharper type from a simple initializer
	// (`var x = new Foo()` → Foo instance). For everything else, we
	// still record the name at Unknown so VarCheck doesn't flag later
	// references to it as undeclared. Mirrors how TypedScopeCreator
	// always adds a slot for every `var` even when its type stays
	// inferred.
	init := name.FirstChild()
	if init != nil && init.Token() == rhino.NEW {
		if callee := init.FirstChild(); callee != nil && callee.Token() == rhino.NAME {
			if fn, ok := scope.declared[callee.String_()].(*jstype.FunctionType); ok && fn.Instance != nil {
				scope.declared[name.String_()] = fn.Instance
				return
			}
		}
	}
	if _, exists := scope.declared[name.String_()]; !exists {
		scope.declared[name.String_()] = jstype.Unknown
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

	// @template names bind to fresh TemplateType values inside a shim
	// scope; the same TemplateType pointers are recorded on the
	// resulting FunctionType so call-site unification can refer to
	// them. Mirrors how Java's JSTypeRegistry builds a
	// TemplateTypeMap before resolving per-function annotations.
	resolve, templateMap := c.templateResolver(scope, info)
	var paramTypes []jstype.JSType
	for _, p := range params.Children() {
		var t jstype.JSType = jstype.Unknown
		if pexpr := info.ParameterType(p.String_()); pexpr != nil {
			t = resolve(pexpr.TypeName)
		}
		paramTypes = append(paramTypes, t)
	}
	// JS semantics for a function without `return expr` is to return
	// `undefined`. When @return is absent we default to Void instead
	// of Unknown so a `+` / assignment chain that consumes the return
	// value (`a + f()`) still type-checks against `string` /
	// `number`. Constructors keep the existing Instance return.
	var retType jstype.JSType = jstype.Void
	if r := info.ReturnType(); r != nil {
		retType = resolve(r.TypeName)
	} else if functionReturnsExplicitly(fn) {
		// Body has at least one `return expr;` — without dataflow we
		// can't infer the return type cheaply, so stay Unknown.
		retType = jstype.Unknown
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
	if templateMap != nil {
		fnType.TemplateTypes = make([]*jstype.TemplateType, 0, len(templateMap))
		for _, tn := range info.TemplateNames() {
			if tt := templateMap[tn]; tt != nil {
				fnType.TemplateTypes = append(fnType.TemplateTypes, tt)
			}
		}
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

// functionReturnsExplicitly reports whether a FUNCTION body contains
// any `return EXPR` statement. Used so `function f(){}` is treated as
// returning Void while `function g(){return x;}` stays Unknown when
// no @return JSDoc is present. Mirrors the body-walk at the start of
// TypedScopeCreator's recordFunction.
func functionReturnsExplicitly(fn *rhino.Node) bool {
	kids := fn.Children()
	if len(kids) < 3 {
		return false
	}
	body := kids[2]
	return hasReturnExpr(body)
}

func hasReturnExpr(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	if n.Token() == rhino.FUNCTION {
		return false // don't descend into nested functions
	}
	if n.Token() == rhino.RETURN && n.FirstChild() != nil {
		return true
	}
	for _, c := range n.Children() {
		if hasReturnExpr(c) {
			return true
		}
	}
	return false
}
