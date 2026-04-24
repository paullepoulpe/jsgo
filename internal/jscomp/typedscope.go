// Port of com/google/javascript/jscomp/TypedScope.java @ closure-compiler v20260415

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// TODO(port): TypedScope.getDepth from TypedScope.java:119
// Reason: depth tracking not needed for current flat-scope model.

// TODO(port): TypedScope.getVar from TypedScope.java:232
// Reason: Java TypedVar wrapper not ported; use GetDeclaredType instead.

// TODO(port): TypedScope.isBottom / createLatticeBottom from TypedScope.java:73
// Reason: lattice-bottom concept unused in Go port's dataflow model.

// TypedScope maps identifier names to their declared JSType. The parent
// pointer mirrors the Java scope chain used during name resolution.
type TypedScope struct {
	parent   *TypedScope
	// rootNode is the AST node that roots this scope (the SCRIPT/FUNCTION/BLOCK
	// node). It mirrors AbstractScope.getRootNode() in Java and is used by
	// NodeTraversal.getScopeRoot(). May be nil for synthetic scopes.
	rootNode *rhino.Node
	declared map[string]jstype.JSType
	// namespaceProps tracks `/** @type {T} */ obj.prop = …`
	// declarations on plain-object bindings (`var goog = {}`).
	// Keyed by owner-name → property-name → declared type. Kept
	// separate from `declared` so the binding's own slot stays
	// Unknown (avoiding FPs on the original `{}` initializer),
	// while downstream assignments / reads can still honour the
	// declared property type. Mirrors how Java's TypedScopeCreator
	// carries AssignmentInferrer-discovered slot declarations on
	// a parallel map.
	namespaceProps map[string]map[string]jstype.JSType
	// namespacePropNodes mirrors namespaceProps with the AST node
	// that carried the declaring JSDoc, used to produce the
	// cross-referencing position on JSC_DUP_VAR_DECLARATION_TYPE
	// _MISMATCH diagnostics.
	namespacePropNodes map[string]map[string]*rhino.Node
	// inferredNsProps tracks namespace properties whose type was
	// inferred from a literal RHS (no @type JSDoc). Assignment
	// checks skip inferred-only slots — Java treats them as open
	// surfaces in default mode, mirroring the IsPropertyInferred
	// guard on prototype slots.
	inferredNsProps map[string]map[string]bool
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

// WithoutRefinements temporarily clears the refinement stack, calls f,
// then restores it. Used by visitors (e.g. visitEqualityOp) that must
// evaluate declared types rather than flow-narrowed types — matching
// how CheckDeterministicTest used its own scope-walking traversal that
// never consulted flow refinements.
func (s *TypedScope) WithoutRefinements(f func()) {
	if s == nil {
		f()
		return
	}
	saved := s.refinements
	s.refinements = nil
	f()
	s.refinements = saved
}

// GetRootNode returns the AST node that roots this scope (SCRIPT/FUNCTION/BLOCK).
// Mirrors AbstractScope.getRootNode() in Java. May be nil for synthetic scopes.
func (s *TypedScope) GetRootNode() *rhino.Node {
	if s == nil {
		return nil
	}
	return s.rootNode
}

// IsHoistScope reports whether this scope is a hoist scope (function, global,
// or module scope — as opposed to a block scope). Mirrors
// AbstractScope.isHoistScope() in Java.
func (s *TypedScope) IsHoistScope() bool {
	if s == nil {
		return false
	}
	if s.rootNode == nil {
		return true // synthetic global scopes are hoist scopes
	}
	switch s.rootNode.Token() {
	case rhino.FUNCTION, rhino.SCRIPT, rhino.ROOT, rhino.MODULE_BODY:
		return true
	}
	return false
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

// GetQualifiedRefinement walks the refinement stack for a dotted
// qualified name (e.g. "obj.prop", "this.x.y"). Returns the flow-
// narrowed type or nil when the name isn't refined. Distinct from
// GetDeclaredType: qualified names aren't tracked in declared
// (only refinements), so the lookup never falls through to the
// parent scope's declared map.
func (s *TypedScope) GetQualifiedRefinement(qname string) jstype.JSType {
	for cur := s; cur != nil; cur = cur.parent {
		for i := len(cur.refinements) - 1; i >= 0; i-- {
			if t, ok := cur.refinements[i][qname]; ok {
				return t
			}
		}
	}
	return nil
}

// GetRawDeclaredType skips the refinement overlay and returns only
// the binding's declared type from the scope chain. Used by the
// assignment-target lookup in visitAssign so `x = 1` compares the
// RHS against the original declaration, not the flow-narrowed type
// (otherwise a prior `var x = null;` would make x's type null and
// every subsequent assignment would false-positive).
func (s *TypedScope) GetRawDeclaredType(name string) jstype.JSType {
	if s == nil {
		return nil
	}
	if t, ok := s.declared[name]; ok {
		return t
	}
	return s.parent.GetRawDeclaredType(name)
}

// Registry returns the JSType registry for this scope.
// Used by template-instantiation helpers in superclasscheck.go.
func (s *TypedScope) Registry() *jstype.Registry { return s.registry }

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

// LookupNamespaceProp returns the type declared for `owner.prop`
// via a `/** @type {T} */ owner.prop = …` statement earlier in
// the scope chain. Returns nil when no such declaration exists.
func (s *TypedScope) LookupNamespaceProp(owner, prop string) jstype.JSType {
	if s == nil {
		return nil
	}
	if m, ok := s.namespaceProps[owner]; ok {
		if t, ok := m[prop]; ok {
			return t
		}
	}
	return s.parent.LookupNamespaceProp(owner, prop)
}

// recordNamespaceProp remembers `owner.prop : t` on this scope.
// Returns the first-declaration node when this call overrides an
// existing entry — callers use that node to emit the "original
// definition at …" cross-reference on
// JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH diagnostics.
func (s *TypedScope) recordNamespaceProp(owner, prop string, t jstype.JSType, at *rhino.Node) *rhino.Node {
	if s.namespaceProps == nil {
		s.namespaceProps = map[string]map[string]jstype.JSType{}
	}
	if s.namespaceProps[owner] == nil {
		s.namespaceProps[owner] = map[string]jstype.JSType{}
	}
	var firstNode *rhino.Node
	if s.namespacePropNodes != nil && s.namespacePropNodes[owner] != nil {
		firstNode = s.namespacePropNodes[owner][prop]
	}
	s.namespaceProps[owner][prop] = t
	if s.namespacePropNodes == nil {
		s.namespacePropNodes = map[string]map[string]*rhino.Node{}
	}
	if s.namespacePropNodes[owner] == nil {
		s.namespacePropNodes[owner] = map[string]*rhino.Node{}
	}
	if firstNode == nil {
		s.namespacePropNodes[owner][prop] = at
	}
	return firstNode
}

// markNamespacePropInferred flags owner.prop as literal-inferred.
// The assignment checker skips inferred-only namespace props, matching
// Java's open-surface treatment of non-@type namespace assignments.
func (s *TypedScope) markNamespacePropInferred(owner, prop string) {
	if s.inferredNsProps == nil {
		s.inferredNsProps = map[string]map[string]bool{}
	}
	if s.inferredNsProps[owner] == nil {
		s.inferredNsProps[owner] = map[string]bool{}
	}
	s.inferredNsProps[owner][prop] = true
}

// IsNamespacePropInferred reports whether owner.prop was registered
// via a literal-RHS assignment (no @type JSDoc) rather than an
// explicit type declaration. Assignment checks consult this to avoid
// false-positive TYPE_MISMATCH on open-surface re-assignments.
func (s *TypedScope) IsNamespacePropInferred(owner, prop string) bool {
	if s == nil {
		return false
	}
	if m, ok := s.inferredNsProps[owner]; ok {
		if m[prop] {
			return true
		}
	}
	return s.parent.IsNamespacePropInferred(owner, prop)
}

// LookupEnum returns the EnumType associated with `name` in this
// scope chain, or nil when the binding isn't an enum. Used by
// ResolveExpr to route `@type {MyEnum}` to the matching
// EnumElementType.
func (s *TypedScope) LookupEnum(name string) *jstype.EnumType {
	if s == nil {
		return nil
	}
	if t, ok := s.declared[name]; ok {
		if e, ok := t.(*jstype.EnumType); ok {
			return e
		}
	}
	return s.parent.LookupEnum(name)
}

// LookupNominalType returns the InstanceType associated with a name
// that refers to a known constructor in this scope chain. Used when
// resolving `@type {Foo}` to the instance type rather than the
// constructor function type. Template args (e.g. "Foo<string>") are
// stripped before lookup so that @extends {Foo<string>} resolves to Foo.
func (s *TypedScope) LookupNominalType(name string) *jstype.ObjectType {
	if s == nil {
		return nil
	}
	// Strip template args (e.g. "Foo<string>" → "Foo") so that a templatized
	// @extends expression resolves to the raw constructor. Mirrors Java's
	// TypedScopeCreator which passes the raw type name (not the full
	// templatized expression) to scope lookup.
	if i := strings.IndexByte(name, '<'); i >= 0 {
		name = name[:i]
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
	// Inline function type: `function(T1,T2=,...T3):R`. Produced by
	// JsDocInfoParser.parseFunctionType(). Mirrors how Java's
	// JSTypeRegistry converts a FunctionTypeAST into a FunctionType.
	// Note: `Function` (capital F) is the built-in constructor type and
	// must NOT enter this branch — it resolves via the name-lookup path
	// below (just like `Array`, `Object`, etc.). Only bare `function`
	// (keyword) and `function(…)` expressions enter here.
	if startsWith(expr, "function(") || expr == "function" {
		if fn := s.resolveFunctionTypeExpr(expr); fn != nil {
			if nullable {
				return jstype.NewUnion(fn, jstype.Null)
			}
			return fn
		}
		// Fall back to the built-in Function type if resolution fails.
		if fn := s.GetDeclaredType("Function"); fn != nil {
			if nullable {
				return jstype.NewUnion(fn, jstype.Null)
			}
			return fn
		}
		return jstype.Unknown
	}
	// Inline structural: `{k1: T1, k2: T2}`. Only take the early
	// structural branch when the outer braces truly wrap the whole
	// expression (i.e. the opening `{` is matched by the trailing
	// `}` with no higher-level operator in between). Otherwise a
	// union like `{a:number}|{b:string}` would be misparsed as a
	// malformed record.
	if startsWith(expr, "{") && balancedOuterBraces(expr) {
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
	// In Closure's legacy mode (the one the corpus tests run in),
	// nominal (non-primitive) types are nullable by default — `Foo`
	// means `Foo|null`, and `!Foo` strips the null. Primitives stay
	// non-nullable. This flag accumulates across union members so a
	// single nominal piece nullifies the whole expression.
	impliedNullable := false
	for _, p := range parts {
		p = trimSpace(p)
		explicitNonNull := false
		explicitNullable := false
		if startsWith(p, "!") {
			explicitNonNull = true
			p = trimSpace(p[1:])
		} else if startsWith(p, "?") {
			explicitNullable = true
			p = trimSpace(p[1:])
		}
		// Inline structural member: `{…}` inside a union.
		if startsWith(p, "{") && balancedOuterBraces(p) {
			if obj, ok := s.resolveObjectLiteralType(p); ok {
				if !explicitNonNull {
					impliedNullable = true
				}
				members = append(members, obj)
				continue
			}
			return jstype.Unknown
		}
		// Function type expression: `function(T1, T2): R`. Mirrors how
		// Closure JSDoc represents function types with parameter/return
		// type annotations. Explicitly non-null functions don't get a
		// null partner added. Resolves transitively so `function(!Object<!Object,number>)`
		// produces a FunctionType whose param is a TemplatizedType.
		if startsWith(p, "function(") {
			if fn := s.resolveFunctionTypeExpr(p); fn != nil {
				members = append(members, fn)
				continue
			}
			return jstype.Unknown
		}
		// Generic brackets: `Array<string>` → TemplatizedType backed
		// by Array's instance type, with `string` resolved as the
		// type argument. Mirrors how JSTypeRegistry instantiates a
		// templated nominal at JSDoc-resolution time.
		var tplArgs []jstype.JSType
		if i := indexOfByte(p, '<'); i >= 0 {
			closing := indexLastByte(p, '>')
			if closing > i {
				inner := trimSpace(p[i+1 : closing])
				for _, a := range splitTopLevelComma(inner) {
					tplArgs = append(tplArgs, s.ResolveExpr(a))
				}
			}
			p = trimSpace(p[:i])
		}
		if p == "" {
			return jstype.Unknown
		}
		// @typedef aliases take priority over nominal constructors
		// and primitives (matches Java's StaticSourceFile resolution).
		// \`!Foo\` strips nullable members from the resolved alias so
		// \`!?number\` → number.
		if alias := s.LookupTypeAlias(p); alias != nil {
			if explicitNonNull {
				alias = stripNullUndef(alias)
			}
			// If the alias resolved to Unknown (e.g. due to a cyclic typedef
			// forward reference), Unknown absorbs the entire union — return it
			// immediately. Mirrors Java's type lattice invariant.
			if isJSUnknown(alias) {
				return jstype.Unknown
			}
			members = append(members, alias)
			continue
		}
		// `@type {MyEnum}` resolves to the EnumElementType (so it
		// means "any element of MyEnum"), not the bare EnumType.
		// Matches how Java's JSTypeRegistry.getType dispatches
		// a bare enum name to its EnumElementType.
		if enum := s.LookupEnum(p); enum != nil {
			members = append(members, enum.Element)
			continue
		}
		// Try nominal (constructor → instance), then primitive.
		if inst := s.LookupNominalType(p); inst != nil {
			var member jstype.JSType = inst
			if len(tplArgs) > 0 {
				member = jstype.NewTemplatizedType(inst, tplArgs)
			}
			if !explicitNonNull {
				impliedNullable = true
			}
			if explicitNullable {
				impliedNullable = true
			}
			members = append(members, member)
			continue
		}
		// Dotted name: `goog.Foo` / `ns.Class.Inner`. Resolve the
		// first segment as a namespace binding and consult its
		// namespace-prop map for the property's type. A single-step
		// chain that resolves to a constructor's Instance yields a
		// nominal member, matching how Java walks dotted type names
		// through StaticSlot.getInstanceType().
		if idx := indexOfByte(p, '.'); idx > 0 {
			first := p[:idx]
			rest := p[idx+1:]
			if t := s.LookupNamespaceProp(first, rest); t != nil {
				if fn, ok := t.(*jstype.FunctionType); ok && fn.Instance != nil {
					var member jstype.JSType = fn.Instance
					if len(tplArgs) > 0 {
						member = jstype.NewTemplatizedType(fn.Instance, tplArgs)
					}
					if !explicitNonNull {
						impliedNullable = true
					}
					if explicitNullable {
						impliedNullable = true
					}
					members = append(members, member)
					continue
				}
				// `@type {ns.MyEnum}` where ns.MyEnum is an EnumType
				// resolves to the EnumElementType, matching how Java's
				// JSTypeRegistry dispatches a dotted enum name to its
				// element type (same as a bare enum name lookup does via
				// LookupEnum → enum.Element).
				if enum, ok := t.(*jstype.EnumType); ok {
					members = append(members, enum.Element)
					continue
				}
				members = append(members, t)
				continue
			}
			// Static nominal under a known constructor: `C.Inner`.
			if fn, ok := s.GetDeclaredType(first).(*jstype.FunctionType); ok && fn.Instance != nil {
				if t := fn.Instance.GetStaticPropertyType(rest); t != nil {
					if fnT, ok := t.(*jstype.FunctionType); ok && fnT.Instance != nil {
						members = append(members, fnT.Instance)
						if !explicitNonNull {
							impliedNullable = true
						}
						continue
					}
					members = append(members, t)
					continue
				}
			}
		}
		// Primitive fallback.
		if explicitNullable {
			impliedNullable = true
		}
		m := s.registry.ResolveName(p)
		if isJSUnknown(m) {
			return jstype.Unknown
		}
		members = append(members, m)
	}
	if nullable || impliedNullable {
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

func indexLastByte(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// splitTopLevelComma splits `s` on `,` at depth 0 of {} () <>. Used
// for template-arg lists like `Map<string, Array<number>>`.
func splitTopLevelComma(s string) []string {
	var out []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{', '(', '<':
			depth++
		case '}', ')', '>':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, trimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	out = append(out, trimSpace(s[start:]))
	return out
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

// resolveFunctionTypeExpr converts a serialised function-type string
// produced by JsDocInfoParser.parseFunctionType — e.g.
// `function(number,string=,...boolean):Date` — into a *FunctionType
// with correctly computed Params, MinArity, MaxArity, and Return.
//
// Format conventions (mirrors the serialiser in parseFunctionType):
//   - Params separated by commas inside the outer parens.
//   - A trailing `=` suffix on a param marks it optional.
//   - A leading `...` prefix on a param marks it variadic (rest).
//   - An optional `:ReturnType` after the closing `)` gives the return.
//
// Returns nil when the expression is not a recognisable function type
// (empty, `Function`, or malformed), so the caller can fall back.
func (s *TypedScope) resolveFunctionTypeExpr(expr string) *jstype.FunctionType {
	if expr == "Function" || expr == "function" {
		// Bare `Function` / `function` — return nil so the caller falls
		// back to the registered built-in Function constructor type
		// (IsConstructor=true), matching Java's FUNCTION_TYPE which is
		// Kind.CONSTRUCTOR and therefore passes couldBeAConstructor.
		return nil
	}
	// Must start with `function(`.
	if !startsWith(expr, "function(") {
		return nil
	}
	// Find the matching `)` at depth 0 (params may contain nested
	// parens from record / generic types).
	inner := expr[len("function("):]
	depth := 1
	closeIdx := -1
	for i := 0; i < len(inner); i++ {
		switch inner[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				closeIdx = i
			}
		}
		if closeIdx >= 0 {
			break
		}
	}
	if closeIdx < 0 {
		return nil
	}
	paramStr := inner[:closeIdx]
	rest := inner[closeIdx+1:] // `:ReturnType` or ""

	// Resolve return type.
	var retType jstype.JSType = jstype.Void
	if startsWith(rest, ":") {
		retType = s.ResolveExpr(rest[1:])
	}

	// Split param list on top-level commas.
	var paramParts []string
	if trimSpace(paramStr) != "" {
		paramParts = splitTopLevelComma(paramStr)
	}

	// Detect `new:T` and `this:T` pseudo-params emitted by
	// parseFunctionType. `new:T` marks the function as a constructor
	// (FunctionType.Kind.CONSTRUCTOR). `this:T` sets typeOfThis, which
	// drives JSC_EXPECTED_THIS_TYPE and FunctionType.getTypeOfThis().
	isCtor := false
	var typeOfThis jstype.JSType
	var positionalParts []string
	for _, pp := range paramParts {
		trimmed := trimSpace(pp)
		switch {
		case startsWith(trimmed, "new:"):
			isCtor = true
		case startsWith(trimmed, "this:"):
			typeOfThis = s.ResolveExpr(trimmed[5:])
		default:
			positionalParts = append(positionalParts, pp)
		}
	}
	paramParts = positionalParts

	var paramTypes []jstype.JSType
	var parameters []jstype.Parameter
	minArity := 0
	maxArity := 0
	sawOptional := false
	for i, pp := range paramParts {
		pp = trimSpace(pp)
		isRest := false
		isOptional := false
		typeStr := pp
		if startsWith(pp, "...") {
			isRest = true
			typeStr = pp[3:]
		} else if endsWith(pp, "=") {
			isOptional = true
			typeStr = pp[:len(pp)-1]
		}
		t := s.ResolveExpr(typeStr)
		paramTypes = append(paramTypes, t)
		parameters = append(parameters, jstype.NewParameter(t, isOptional, isRest))
		switch {
		case isRest:
			maxArity = jstype.MaxArgsUnlimited
			if !sawOptional {
				minArity = i
			}
			// Variadic must be the last param; stop here.
			sawOptional = true
		case isOptional:
			sawOptional = true
			if maxArity != jstype.MaxArgsUnlimited {
				maxArity = i + 1
			}
		default:
			if !sawOptional {
				minArity = i + 1
			}
			if maxArity != jstype.MaxArgsUnlimited {
				maxArity = i + 1
			}
		}
	}
	if len(paramParts) == 0 {
		maxArity = 0
		minArity = 0
	}

	// Build the FunctionType with the rich parameter slice so callers
	// can consult IsOptional() / IsVariadic() on each param. This
	// mirrors FunctionType.getParameters() in Java and is needed by the
	// optional-parameter check in TypeCheck.visitCall.
	// NewFunctionTypeFromParameters recomputes MinArity / MaxArity via
	// computeArity, which matches the manual computation above.
	_ = paramTypes // superseded by rich parameters slice
	_ = minArity   // computed by NewFunctionTypeFromParameters
	_ = maxArity   // computed by NewFunctionTypeFromParameters
	fn := jstype.NewFunctionTypeFromParameters(parameters, retType)
	// Preserve the @constructor marker from `new:T` pseudo-params.
	// Mirrors FunctionType.builder.forConstructor() when new:T is present.
	if isCtor {
		fn.IsConstructor = true
	}
	// Wire `this:T` through to the function's typeOfThis so
	// EXPECTED_THIS_TYPE checks + call-site this binding fire correctly.
	if typeOfThis != nil {
		fn.SetTypeOfThis(typeOfThis)
	}
	return fn
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

// stripTypeArgs removes a trailing `<…>` suffix from a type name so
// `Foo<string>` → `Foo`. Mirrors TypedScopeCreator's
// stripTemplateArgs helper used when resolving @extends / @implements
// names that carry generic type arguments. Only strips when the first
// `<` is balanced by a matching `>` at the end.
func stripTypeArgs(name string) string {
	i := indexOfByte(name, '<')
	if i <= 0 {
		return name
	}
	return name[:i]
}

// balancedOuterBraces reports whether the outermost `{…}` in `s`
// wraps the entire string. Returns false for shapes like
// `{a}|{b}` where the first `{...}` closes before the end.
func balancedOuterBraces(s string) bool {
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return false
	}
	depth := 0
	for i := 0; i < len(s)-1; i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 && i < len(s)-1 {
				return false
			}
		}
	}
	return true
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

// syntheticReadonlyArrayType returns a synthetic jstype for a ReadonlyArray
// annotation when ReadonlyArray is not registered in the scope (to avoid
// false-positive JSC_TYPE_MISMATCH for `@const {!ReadonlyArray<T>}` decls).
// It creates a bare ObjectType named "ReadonlyArray" for bare annotations, or
// a TemplatizedType{ReadonlyArray, Unknown} for parameterized ones, so that
// checkNotReadonlyPropertyAssignment can identify the ReadonlyArray pattern.
// Returns nil if the raw type string does not mention ReadonlyArray.
func syntheticReadonlyArrayType(rawTypeName string) jstype.JSType {
	// Quick check: does the annotation mention ReadonlyArray?
	base := rawTypeName
	// Strip leading `!`, `?`, spaces.
	for len(base) > 0 && (base[0] == '!' || base[0] == '?' || base[0] == ' ') {
		base = base[1:]
	}
	if base != "ReadonlyArray" && !startsWith(base, "ReadonlyArray<") {
		return nil
	}
	roInst := jstype.NewObjectType("ReadonlyArray")
	roInst.DefineProperty("length", jstype.Number)
	if startsWith(base, "ReadonlyArray<") {
		// Parameterized: extract the type argument for accurate TypeName().
		// e.g. "ReadonlyArray<number>" → targ = "number" → Number type.
		closing := indexLastByte(base, '>')
		if closing > len("ReadonlyArray<") {
			inner := trimSpace(base[len("ReadonlyArray<"):closing])
			// Use the registry/scope to resolve — but we don't have a scope
			// here. Use a simple primitive lookup instead.
			var targ jstype.JSType = jstype.Unknown
			switch inner {
			case "number":
				targ = jstype.Number
			case "string":
				targ = jstype.String
			case "boolean":
				targ = jstype.Boolean
			}
			return jstype.NewTemplatizedType(roInst, []jstype.JSType{targ})
		}
		return jstype.NewTemplatizedType(roInst, []jstype.JSType{jstype.Unknown})
	}
	return roInst
}

// namespacePropNode returns the AST node (typically the property
// STRING_KEY) that carried the first `/** @type {T} */` namespace
// declaration for `owner.prop`. Nil when none exists. Consumed by
// the DUP_VAR_DECLARATION_TYPE_MISMATCH diagnostic.
func (s *TypedScope) namespacePropNode(owner, prop string) *rhino.Node {
	if s == nil {
		return nil
	}
	if m, ok := s.namespacePropNodes[owner]; ok {
		if n, ok := m[prop]; ok {
			return n
		}
	}
	return s.parent.namespacePropNode(owner, prop)
}

// inferExprType is a lightweight type inferencer for simple expressions
// in scope context. Used to get the type of a for-of iterable at
// scope-creation time (before TypeInference runs).
func (s *TypedScope) inferExprType(n *rhino.Node) jstype.JSType {
	if n == nil {
		return jstype.Unknown
	}
	switch n.Token() {
	case rhino.NAME:
		t := s.GetDeclaredType(n.String_())
		if t == nil {
			return jstype.Unknown
		}
		return t
	case rhino.GETPROP:
		// Delegate to the qualified-name lookup.
		if qname := NU.GetBestLValueName(n); qname != "" {
			parts := strings.SplitN(qname, ".", 2)
			if len(parts) == 2 {
				if fn, ok := s.GetDeclaredType(parts[0]).(*jstype.FunctionType); ok && fn.Instance != nil {
					if t := fn.Instance.GetStaticPropertyType(parts[1]); t != nil {
						return t
					}
				}
			}
		}
		return jstype.Unknown
	}
	return jstype.Unknown
}

// ResolveName resolves a single type name against the scope chain,
// checking typedefs, enums, nominal constructors, and primitives in
// that order. Returns Unknown when nothing matches. Used by
// ResolveWorklist as the post-scope-creation resolver.
func (s *TypedScope) ResolveName(name string) jstype.JSType {
	if s == nil {
		return jstype.Unknown
	}
	// Dotted name: check namespace props (e.g. "goog.Foo").
	if dot := strings.IndexByte(name, '.'); dot >= 0 {
		owner := name[:dot]
		prop := name[dot+1:]
		if t := s.LookupNamespaceProp(owner, prop); t != nil {
			if fn, ok := t.(*jstype.FunctionType); ok && fn.Instance != nil {
				return fn.Instance
			}
			// Non-constructor namespace prop — still present, not a type.
			return jstype.Unknown
		}
	}
	if alias := s.LookupTypeAlias(name); alias != nil {
		return alias
	}
	if enum := s.LookupEnum(name); enum != nil {
		return enum.Element
	}
	if inst := s.LookupNominalType(name); inst != nil {
		return inst
	}
	return s.registry.ResolveName(name)
}
