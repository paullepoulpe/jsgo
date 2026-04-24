// Port of com/google/javascript/jscomp/TypedScope.java @ closure-compiler v20260415

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// TypedScope contains information about variables and their types. Scopes can
// be nested; a scope points back to its parent scope.
//
// TypedScope is also used as a lattice element for flow-sensitive type
// inference. As a lattice element, a scope is viewed as a map from names to
// types. A name not in the map is considered to have the bottom type.
//
// TypedScopes also have a concept of a set of 'reserved' names. These are
// names that will be declared as actual TypedVars later in the compilation.
// Looking up any of these reserved names will always return null, even if the
// name is available in a parent scope.
//
// Mirrors com/google/javascript/jscomp/TypedScope.java.
type TypedScope struct {
	parent *TypedScope
	// depth is the nesting depth of this scope (0 = global).
	// Mirrors TypedScope.depth in Java.
	depth int
	// isBottom is true when this scope is the lattice bottom for type inference.
	// Mirrors TypedScope.isBottom in Java.
	isBottom bool
	// reservedNames holds names that will be declared later. Looking up reserved
	// names returns nil. Mirrors TypedScope.reservedNames in Java.
	reservedNames map[string]bool
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

// Go-only: language-adapter — CreateGlobalTypedScope is the Go name for
// TypedScope.createGlobalScope(Node); the "Typed" prefix and extra registry
// parameter distinguish it from UntypedScope's CreateGlobalUntypedScope.
func CreateGlobalTypedScope(rootNode *rhino.Node, registry *jstype.Registry) *TypedScope {
	return &TypedScope{
		parent:   nil,
		depth:    0,
		isBottom: false,
		rootNode: rootNode,
		registry: registry,
	}
}

// CreateLatticeBottom creates the lattice-bottom TypedScope used by type inference.
// Mirrors TypedScope.createLatticeBottom(Node) in Java.
func CreateLatticeBottom(rootNode *rhino.Node) *TypedScope {
	return &TypedScope{
		parent:   nil,
		depth:    0,
		isBottom: true,
		rootNode: rootNode,
	}
}

// Typed returns this TypedScope (mirrors AbstractScope.typed() in Java, overridden here).
// Mirrors TypedScope.typed() in Java.
func (s *TypedScope) Typed() *TypedScope { return s }

// IsBottom reports whether this scope is the lattice bottom for type inference.
// Mirrors TypedScope.isBottom() in Java.
func (s *TypedScope) IsBottom() bool {
	if s == nil {
		return false
	}
	return s.isBottom
}

// GetDepth returns the nesting depth (0 = global).
// Mirrors TypedScope.getDepth() in Java.
func (s *TypedScope) GetDepth() int {
	if s == nil {
		return 0
	}
	return s.depth
}

// GetParent returns the parent TypedScope, or nil for the global scope.
// Mirrors TypedScope.getParent() in Java.
func (s *TypedScope) GetParent() *TypedScope {
	if s == nil {
		return nil
	}
	return s.parent
}

// GetParentScope returns the parent as an AbstractScope interface.
// Mirrors AbstractScope.getParentScope() in Java.
// Go-only: language-adapter — Go interface dispatch requires a separate method
// returning AbstractScope; Java uses the same getParent() via generics.
func (s *TypedScope) GetParentScope() AbstractScope {
	if s == nil || s.parent == nil {
		return nil
	}
	return s.parent
}

// ValidateCompletelyBuilt asserts that all reserved names have been declared.
// Mirrors TypedScope.validateCompletelyBuilt() in Java.
func (s *TypedScope) ValidateCompletelyBuilt() {
	// In the Go port, reservedNames is pre-populated and cleared as names are declared.
	// If any remain, it indicates a bug in TypedScopeCreator.
	if len(s.reservedNames) > 0 {
		panic("TypedScope: reserved names not declared: " + joinReservedNames(s.reservedNames))
	}
}

// joinReservedNames formats the reserved name set for an error message.
// Go-only: language-adapter — Go has no built-in set-to-string formatting.
func joinReservedNames(m map[string]bool) string {
	sb := "{"
	first := true
	for k := range m {
		if !first {
			sb += ", "
		}
		sb += k
		first = false
	}
	return sb + "}"
}

// GetTypeOfThis returns the type of `this` in the current scope.
// Mirrors TypedScope.getTypeOfThis() in Java.
func (s *TypedScope) GetTypeOfThis() jstype.JSType {
	if s == nil {
		return nil
	}
	if s.IsGlobal() {
		if s.rootNode != nil {
			if t, ok := s.rootNode.GetJSType().(jstype.JSType); ok {
				return t
			}
		}
		return nil
	}
	if nodeIsNonArrowFunction(s.rootNode) {
		if nodeType, ok := s.rootNode.GetJSType().(jstype.JSType); ok && nodeType != nil {
			if fn, ok := nodeType.(*jstype.FunctionType); ok {
				return fn.GetTypeOfThis()
			}
		}
		return nil
	}
	if s.IsStaticBlockScope() {
		if p := s.GetParent(); p != nil && p.rootNode != nil {
			if t, ok := p.rootNode.GetJSType().(jstype.JSType); ok {
				return t
			}
		}
		return nil
	}
	if s.parent != nil {
		return s.parent.GetTypeOfThis()
	}
	return nil
}

// IsGlobal reports whether this is the global scope (no parent).
// Mirrors AbstractScope.isGlobal() in Java.
func (s *TypedScope) IsGlobal() bool {
	return s == nil || s.parent == nil
}

// Go-only: language-adapter — IsLocal is a final method on AbstractScope in Java;
// Go requires explicit implementation on TypedScope since Go has no inheritance.
func (s *TypedScope) IsLocal() bool {
	return s != nil && s.parent != nil
}

// IsStaticBlockScope reports whether this scope's root is a class static block.
// Mirrors AbstractScope.isStaticBlockScope() in Java.
func (s *TypedScope) IsStaticBlockScope() bool {
	return nodeIsClassStaticBlock(s.rootNode)
}

// IsMemberFieldDefScope reports whether the root is a member field definition.
// Mirrors AbstractScope.isMemberFieldDefScope() in Java.
func (s *TypedScope) IsMemberFieldDefScope() bool {
	return s != nil && s.rootNode != nil && s.rootNode.Token() == rhino.MEMBER_FIELD_DEF
}

// Go-only: language-adapter — IsFunctionBlockScope is a final method on AbstractScope in Java;
// Go requires explicit implementation since TypedScope cannot inherit AbstractScope.
func (s *TypedScope) IsFunctionBlockScope() bool {
	if s == nil || s.rootNode == nil {
		return false
	}
	if s.rootNode.Token() != rhino.BLOCK {
		return false
	}
	p := s.rootNode.Parent()
	return p != nil && p.Token() == rhino.FUNCTION
}

// Go-only: language-adapter — IsFunctionScope is a final method on AbstractScope in Java;
// Go requires explicit implementation since TypedScope cannot inherit AbstractScope.
func (s *TypedScope) IsFunctionScope() bool {
	return s != nil && s.rootNode != nil && s.rootNode.Token() == rhino.FUNCTION
}

// Go-only: language-adapter — IsModuleScope is a final method on AbstractScope in Java;
// Go requires explicit implementation since TypedScope cannot inherit AbstractScope.
func (s *TypedScope) IsModuleScope() bool {
	return s != nil && s.rootNode != nil && s.rootNode.Token() == rhino.MODULE_BODY
}

// IsBlockScope reports whether the root creates a block scope.
// Mirrors AbstractScope.isBlockScope() in Java.
func (s *TypedScope) IsBlockScope() bool {
	if s == nil {
		return false
	}
	return nodeCreatesBlockScope(s.rootNode)
}

// Go-only: language-adapter — IsCatchScope is a final method on AbstractScope in Java;
// Go requires explicit implementation since TypedScope cannot inherit AbstractScope.
func (s *TypedScope) IsCatchScope() bool {
	if s == nil || s.rootNode == nil || s.rootNode.Token() != rhino.BLOCK {
		return false
	}
	kids := s.rootNode.Children()
	return len(kids) == 1 && kids[0].Token() == rhino.CATCH
}

// Go-only: language-adapter — IsCfgRootScope is a final method on AbstractScope in Java;
// Go requires explicit implementation since TypedScope cannot inherit AbstractScope.
func (s *TypedScope) IsCfgRootScope() bool {
	return nodeIsValidCfgRoot(s.rootNode)
}

// GetDeclarativelyUnboundVarsWithoutTypes returns vars declared with 'var' but
// without a known type. These can safely be set to undefined at the start of
// type inference. Parameters and externs are excluded.
// Mirrors TypedScope.getDeclarativelyUnboundVarsWithoutTypes() in Java.
//
// TODO(port): TypedScope.getDeclarativelyUnboundVarsWithoutTypes from TypedScope.java:257
// Reason: GoTypedVar iteration over the declared map requires integrating the
// GoTypedVar registry (currently TypedScope uses jstype.JSType maps, not GoTypedVars).
func (s *TypedScope) GetDeclarativelyUnboundVarsWithoutTypes() []*GoTypedVar {
	if s == nil {
		return nil
	}
	// GoTypedVar iteration not yet wired; declared map holds JSType values.
	// Return empty slice (safe: callers treat nil return as no vars to initialize).
	return []*GoTypedVar{}
}

// GetVar returns the variable for name, considering qualified-name shadowing.
// Mirrors TypedScope.getVar(String) in Java (the overridden version that handles
// qualified names and shadowing by walking the scope chain and checking that
// the root name of a dotted name is not shadowed by a local binding).
//
// The Go port uses GetDeclaredType for type lookups; GoTypedVar map integration
// is tracked by the TODO(port) below.
//
// TODO(port): TypedScope.getVar from TypedScope.java:288
// Reason: GoTypedVar scope map not yet integrated into TypedScope.declared;
// callers use GetDeclaredType for type information instead of getVar().getType().
func (s *TypedScope) GetVar(name string) *GoTypedVar {
	if s == nil {
		return nil
	}
	// Walk the scope chain looking for a simple-name match.
	// Qualified names (dot-delimited) use the root name for shadowing check.
	dot := strings.IndexByte(name, '.')
	rootName := name
	if dot >= 0 {
		rootName = name[:dot]
	}
	// Walk up to find where rootName is declared.
	for sc := s; sc != nil; sc = sc.parent {
		if _, ok := sc.declared[rootName]; ok {
			// Found the root; qualified names belong to the same scope.
			return nil // GoTypedVar not yet wired into declared map
		}
	}
	return nil
}

// GetTypeThroughNamespace resolves a dotted module-id through the scope's
// namespace chain (e.g. "goog.foo" → look up "goog", get property "foo").
// Mirrors TypedScope.getTypeThroughNamespace(String) in Java.
func (s *TypedScope) GetTypeThroughNamespace(moduleId string) jstype.JSType {
	dot := strings.LastIndex(moduleId, ".")
	if dot >= 0 {
		parentName := moduleId[:dot]
		prop := moduleId[dot+1:]
		parentType := s.GetTypeThroughNamespace(parentName)
		if parentType == nil {
			return nil
		}
		if obj, ok := parentType.(*jstype.ObjectType); ok {
			return obj.GetPropertyType(prop)
		}
		return nil
	}
	// Simple name: look up in the scope chain.
	t := s.GetDeclaredType(moduleId)
	return t
}

// GetTopmostScopeOfEventualDeclaration returns the topmost scope that
// declares (or will declare) the given name, walking up the chain.
// Mirrors TypedScope.getTopmostScopeOfEventualDeclaration(String) in Java.
func (s *TypedScope) GetTopmostScopeOfEventualDeclaration(name string) *TypedScope {
	if s == nil {
		return nil
	}
	if _, ok := s.declared[name]; ok {
		return s
	}
	if s.reservedNames[name] {
		return s
	}
	if s.parent == nil {
		return nil
	}
	return s.parent.GetTopmostScopeOfEventualDeclaration(name)
}

// Go-only: language-adapter — PushRefinement replaces Java's FlowScope
// refinement stack, which Java TypeInference manages separately from TypedScope.
// The Go port embeds the refinement stack directly in TypedScope for simplicity.
func (s *TypedScope) PushRefinement(refs map[string]jstype.JSType) {
	if s == nil || len(refs) == 0 {
		return
	}
	s.refinements = append(s.refinements, refs)
}

// Go-only: language-adapter — PopRefinement is the paired cleanup for
// PushRefinement; Java uses FlowScope.popInto instead.
func (s *TypedScope) PopRefinement() {
	if s == nil || len(s.refinements) == 0 {
		return
	}
	s.refinements = s.refinements[:len(s.refinements)-1]
}

// Go-only: language-adapter — WithoutRefinements allows visiting with only
// declared types visible; Java TypeInference uses separate FlowScope objects
// rather than toggling a stack on a single TypedScope.
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

// Go-only: language-adapter — IsHoistScope is a final method on AbstractScope in Java;
// Go requires explicit implementation since TypedScope cannot inherit AbstractScope.
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

// GetDeclaredType looks up name with refinement overlay applied.
// Go-only: language-adapter — combines getOwnSlot+getType+FlowScope into one method.
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

// Go-only: language-adapter — GetQualifiedRefinement walks the Go-specific
// refinement stack for dotted names; Java's FlowScope handles qualified-name
// narrowing through a separate LinkedFlowScope chain.
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

// Go-only: language-adapter — GetRawDeclaredType bypasses the Go refinement stack
// to see the originally declared type; Java does not need this because FlowScope
// and TypedScope are separate objects.
func (s *TypedScope) GetRawDeclaredType(name string) jstype.JSType {
	if s == nil {
		return nil
	}
	if t, ok := s.declared[name]; ok {
		return t
	}
	return s.parent.GetRawDeclaredType(name)
}

// Go-only: language-adapter — Registry() exposes the JSTypeRegistry pointer
// held on the scope; in Java the registry is accessed via compiler.getTypeRegistry().
func (s *TypedScope) Registry() *jstype.Registry { return s.registry }

// Go-only: language-adapter — ThisType() is a Go-specific convenience that
// returns the ObjectType of `this`; Java uses TypedScope.getTypeOfThis()
// which is now also ported as GetTypeOfThis().
func (s *TypedScope) ThisType() *jstype.ObjectType {
	if s == nil {
		return nil
	}
	if s.thisType != nil {
		return s.thisType
	}
	return s.parent.ThisType()
}

// Go-only: language-adapter — LookupTypeAlias provides @typedef resolution
// through the scope chain; Java resolves typedefs through JSTypeRegistry.getType()
// which has access to the full typedef registry directly.
func (s *TypedScope) LookupTypeAlias(name string) jstype.JSType {
	if s == nil {
		return nil
	}
	if t, ok := s.typeAliases[name]; ok {
		return t
	}
	return s.parent.LookupTypeAlias(name)
}

// Go-only: language-adapter — LookupNamespaceProp provides qualified-name type
// lookup through the scope chain; Java resolves obj.prop types through
// ObjectType.getPropertyType() on the object's declared type.
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
// Go-only: language-adapter — recordNamespaceProp stores obj.prop types directly
// on the scope; Java stores them on the ObjectType via ObjectType.defineDeclaredProperty().
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

// Go-only: language-adapter — markNamespacePropInferred tracks literal-inferred
// namespace props; Java uses ObjectType.isPropertyTypeDeclared() to distinguish
// declared from inferred properties.
func (s *TypedScope) markNamespacePropInferred(owner, prop string) {
	if s.inferredNsProps == nil {
		s.inferredNsProps = map[string]map[string]bool{}
	}
	if s.inferredNsProps[owner] == nil {
		s.inferredNsProps[owner] = map[string]bool{}
	}
	s.inferredNsProps[owner][prop] = true
}

// Go-only: language-adapter — IsNamespacePropInferred queries the Go-specific
// inferred-prop map; Java checks ObjectType.isPropertyTypeDeclared(prop) instead.
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

// Go-only: language-adapter — LookupEnum provides enum lookup through the scope
// chain; Java dispatches `@type {MyEnum}` through JSTypeRegistry.getType() which
// stores enum types indexed by name in the registry itself.
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

// Go-only: language-adapter — LookupNominalType maps a constructor name to its
// instance type through the scope chain; Java uses JSTypeRegistry.getType() which
// stores instance types indexed by name directly in the registry.
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

// ResolveExpr resolves a JSDoc type expression string in scope context.
// Go-only: language-adapter — Go port of JSTypeRegistry.getType(StaticScope, String).
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

// Go-only: language-adapter — trimSpace is a zero-allocation string trimmer
// used by ResolveExpr; Java uses String.trim() from the standard library.
func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

// Go-only: stdlib-interop — startsWith avoids importing strings for a single check.
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// Go-only: stdlib-interop — trimPrefix strips a leading substring.
func trimPrefix(s, prefix string) string {
	if startsWith(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// Go-only: stdlib-interop — indexOfByte finds the first occurrence of a byte.
func indexOfByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// Go-only: stdlib-interop — indexLastByte finds the last occurrence of a byte.
func indexLastByte(s string, b byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// splitTopLevelComma splits s on ',' at depth 0 of {} () <> for template args.
// Go-only: language-adapter — Java handles type arg lists structurally via TypeAST.
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

// Go-only: language-adapter — splitTopLevelUnion handles the JSDoc union `A|B`
// split that Java does structurally via TypeAST.
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

// Go-only: language-adapter — resolveObjectLiteralType parses `{k: T}` record
// types textually; Java parses them through JsDocInfoParser.parseRecordType()
// into a TypeAST which is then converted by JSTypeRegistry.
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

// resolveFunctionTypeExpr converts a serialised `function(T1,T2=,...T3):R` string
// produced by JsDocInfoParser.parseFunctionType into a *FunctionType.
// Go-only: language-adapter — Java uses FunctionTypeAST+JSTypeRegistry.createFunctionType.
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

// Go-only: language-adapter — indexOfTopLevel finds the first byte at nesting
// depth 0; used by resolveObjectLiteralType for key:value splitting.
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

// Go-only: stdlib-interop — endsWith avoids importing strings for a single check.
func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// Go-only: language-adapter — stripTypeArgs replicates TypedScopeCreator's
// stripTemplateArgs private method for @extends/@implements name normalization.
func stripTypeArgs(name string) string {
	i := indexOfByte(name, '<')
	if i <= 0 {
		return name
	}
	return name[:i]
}

// Go-only: language-adapter — balancedOuterBraces detects top-level record type
// expressions that Java's TypeAST handles structurally.
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

// Go-only: language-adapter — balancedOuter detects outer-paren grouping in
// type expressions that Java's JsDocInfoParser handles structurally.
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

// Go-only: stdlib-interop — splitOn is a byte-delimited split without importing strings.
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

// Go-only: language-adapter — isJSUnknown centralises the nil/Unknown check;
// Java uses JSType.isUnknownType() which is a virtual method on JSType.
func isJSUnknown(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	return t.TypeName() == "?"
}

// Go-only: language-adapter — syntheticReadonlyArrayType creates a stand-in type
// for ReadonlyArray annotations not registered in the Go type registry; Java
// handles this through the full TypedScopeCreator/JSTypeRegistry type pipeline.
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

// Go-only: language-adapter — namespacePropNode retrieves the declaration node
// for a qualified property type; Java stores this in ObjectType's property map.
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

// Go-only: language-adapter — inferExprType is a lightweight pre-TypeInference
// type lookup; Java runs full TypeInference before scope creation is complete.
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

// Go-only: language-adapter — ResolveName is a Go-specific type-name resolver
// that searches the scope chain before falling back to the registry; Java
// resolves type names through JSTypeRegistry.getType(StaticTypedScope, String).
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
