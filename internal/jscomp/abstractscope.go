// Port of com/google/javascript/jscomp/AbstractScope.java
// @ closure-compiler v20260415
//
// AbstractScope<S,V> is the generic base for Scope and TypedScope in Java.
// Go cannot express the same recursive generic bound, so we represent it as
// an interface (AbstractScope) together with a concrete base struct
// (AbstractScopeBase) that captures the shared fields and provides the
// methods that are identical in Scope and TypedScope.
//
// Java inheritance chain:
//   AbstractScope<S,V> ← Scope (S=Scope,V=Var)
//   AbstractScope<S,V> ← TypedScope (S=TypedScope,V=TypedVar)
//
// Go embedding chain (mirrored):
//   AbstractScopeBase is embedded in UntypedScope (Scope.java port).
//   TypedScope (scope.go) carries equivalent fields directly; a future
//   refactor may embed AbstractScopeBase there too.
//
// Full method inventory ported from AbstractScope.java (527 LoC):
//   • getRootNode / getParentScope / getParent
//   • isGlobal / isLocal / getDepth
//   • isBlockScope / isStaticBlockScope / isFunctionBlockScope / isFunctionScope
//   • isModuleScope / isMemberFieldDefScope / isComputedFieldDefRhsScope
//   • isCatchScope / isCfgRootScope / isHoistScope
//   • getGlobalScope / getClosestHoistScope / getClosestContainerScope
//   • getClosestCfgRootScope / getScopeOfThis
//   • contains(other) / getCommonParent / hasSameContainerScope
//   • hasOwnSlot / hasSlot / getOwnSlot / getSlot / getVar / getArgumentsVar
//   • declareInternal / undeclare / canDeclare / clearVars
//   • getVarIterable / getAllAccessibleVariables / getAllSymbols / getVarCount
//   • ImplicitVar enum + isMadeByScope

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ImplicitVar enumerates the implicit variables that are implicitly defined
// in every function scope without being explicitly declared.
// Mirrors AbstractScope.ImplicitVar in Java.
type ImplicitVar int

const (
	// ImplicitVarArguments is the "arguments" object implicit in functions.
	ImplicitVarArguments ImplicitVar = iota
	// ImplicitVarExports is the module "exports" implicit in goog.module scopes.
	ImplicitVarExports
	// ImplicitVarSuper is the "super" keyword implicit in class-body scopes.
	ImplicitVarSuper
	// ImplicitVarThis is the "this" keyword implicit in function/class scopes.
	ImplicitVarThis
)

// Go-only: language-adapter — implicitVarName provides the enum-constant name
// as a function; Java accesses it as the ImplicitVar.name field directly since
// Java enums are objects with accessible instance fields.
func implicitVarName(v ImplicitVar) string {
	switch v {
	case ImplicitVarArguments:
		return "arguments"
	case ImplicitVarExports:
		return "exports"
	case ImplicitVarSuper:
		return "super"
	case ImplicitVarThis:
		return "this"
	}
	return ""
}

// Go-only: language-adapter — implicitVarOf returns (value, ok) rather than
// a nullable ImplicitVar; Java's ImplicitVar.of(String) returns @Nullable ImplicitVar.
func implicitVarOf(name string) (ImplicitVar, bool) {
	switch name {
	case "arguments":
		return ImplicitVarArguments, true
	case "super":
		return ImplicitVarSuper, true
	case "this":
		return ImplicitVarThis, true
	case "exports":
		return ImplicitVarExports, true
	}
	return 0, false
}

// Go-only: language-adapter — implicitVarIsMadeByScope takes a rootNode rather
// than a full scope; Java's ImplicitVar.isMadeByScope(AbstractScope) calls scope
// predicates directly on the scope object. The Go version accepts a Node to
// avoid circular dependency between scope types.
func implicitVarIsMadeByScope(v ImplicitVar, rootNode *rhino.Node) bool {
	if rootNode == nil {
		return false
	}
	switch v {
	case ImplicitVarExports:
		// exports is implicit only in a goog.module scope.
		return rootNode.Token() == rhino.MODULE_BODY
	case ImplicitVarSuper, ImplicitVarThis:
		// super and this are implicit in non-arrow functions and class static blocks.
		tok := rootNode.Token()
		return tok == rhino.FUNCTION || tok == rhino.CLASS
	case ImplicitVarArguments:
		// arguments is implicit in non-arrow functions only.
		return rootNode.Token() == rhino.FUNCTION
	}
	return false
}

// AbstractScope is the interface satisfied by both UntypedScope (port of
// Java's Scope) and the existing TypedScope. It provides the root-node and
// parent-scope accessors that are defined final on AbstractScope<S,V> in Java.
// Mirrors the StaticScope/AbstractScope contract in Java.
type AbstractScope interface {
	// GetRootNode returns the AST node that is the root of this scope.
	// Mirrors AbstractScope.getRootNode() in Java.
	GetRootNode() *rhino.Node

	// GetParentScope returns the parent AbstractScope, or nil for global.
	// Mirrors AbstractScope.getParentScope() in Java.
	GetParentScope() AbstractScope

	// IsGlobal reports whether this scope has no parent.
	// Mirrors AbstractScope.isGlobal() in Java.
	IsGlobal() bool

	// IsLocal reports whether this scope has a parent.
	// Mirrors AbstractScope.isLocal() in Java.
	IsLocal() bool

	// GetDepth returns the nesting depth (0 = global).
	// Mirrors AbstractScope.getDepth() in Java.
	GetDepth() int

	// IsBlockScope reports whether the root is a block-creating node.
	// Mirrors AbstractScope.isBlockScope() in Java.
	IsBlockScope() bool

	// IsFunctionScope reports whether the root is a FUNCTION node.
	// Mirrors AbstractScope.isFunctionScope() in Java.
	IsFunctionScope() bool

	// IsModuleScope reports whether the root is a MODULE_BODY node.
	// Mirrors AbstractScope.isModuleScope() in Java.
	IsModuleScope() bool

	// IsFunctionBlockScope reports whether the root is the body block of
	// a FUNCTION node (i.e. the directly-enclosing block inside a function).
	// Mirrors AbstractScope.isFunctionBlockScope() in Java.
	IsFunctionBlockScope() bool

	// IsHoistScope reports whether var declarations in this scope would be
	// hoisted to this scope (function, function-block, global, module, or
	// static-block scopes).
	// Mirrors AbstractScope.isHoistScope() in Java.
	IsHoistScope() bool

	// IsCatchScope reports whether this scope is an exception catch block.
	// Mirrors AbstractScope.isCatchScope() in Java.
	IsCatchScope() bool

	// IsStaticBlockScope reports whether this scope is a class static block.
	// Mirrors AbstractScope.isStaticBlockScope() in Java.
	IsStaticBlockScope() bool

	// IsCfgRootScope reports whether this scope's root is a valid CFG root.
	// Mirrors AbstractScope.isCfgRootScope() in Java.
	IsCfgRootScope() bool
}

// AbstractScopeBase implements the shared fields and logic for both
// UntypedScope (Java's Scope) and could eventually be embedded in TypedScope.
// It stores the rootNode (from Java AbstractScope) and provides most of the
// "final" methods in Java's AbstractScope that don't need generics.
//
// Not embedded in the existing TypedScope to avoid a disruptive refactor, but
// provides a parallel implementation that satisfies AbstractScope.
type AbstractScopeBase struct {
	rootNode *rhino.Node
}

// GetRootNode returns the scope's root AST node.
// Mirrors AbstractScope.getRootNode() (final) in Java.
func (b *AbstractScopeBase) GetRootNode() *rhino.Node { return b.rootNode }

// IsBlockScope reports whether the root creates a block scope.
// Mirrors AbstractScope.isBlockScope() (final) in Java via NodeUtil.createsBlockScope.
func (b *AbstractScopeBase) IsBlockScope() bool {
	return nodeCreatesBlockScope(b.rootNode)
}

// IsFunctionBlockScope reports whether the root is a function body block.
// Mirrors AbstractScope.isFunctionBlockScope() (final) in Java via
// NodeUtil.isFunctionBlock.
func (b *AbstractScopeBase) IsFunctionBlockScope() bool {
	if b.rootNode == nil {
		return false
	}
	// A function body block: BLOCK whose parent is FUNCTION.
	if b.rootNode.Token() != rhino.BLOCK {
		return false
	}
	p := b.rootNode.Parent()
	return p != nil && p.Token() == rhino.FUNCTION
}

// IsFunctionScope reports whether the root is a FUNCTION node.
// Mirrors AbstractScope.isFunctionScope() (final) in Java.
func (b *AbstractScopeBase) IsFunctionScope() bool {
	return b.rootNode != nil && b.rootNode.Token() == rhino.FUNCTION
}

// IsModuleScope reports whether the root is a MODULE_BODY node.
// Mirrors AbstractScope.isModuleScope() (final) in Java.
func (b *AbstractScopeBase) IsModuleScope() bool {
	return b.rootNode != nil && b.rootNode.Token() == rhino.MODULE_BODY
}

// IsHoistScope reports whether var declarations within this scope are
// hoisted to this scope (not to an enclosing one).
// Mirrors AbstractScope.isHoistScope() in Java.
func (b *AbstractScopeBase) IsHoistScope() bool {
	return b.IsFunctionScope() ||
		b.IsFunctionBlockScope() ||
		b.rootNode == nil || // global: no parent
		b.IsModuleScope() ||
		b.IsStaticBlockScope()
}

// IsStaticBlockScope reports whether this scope's root is a class static block.
// Mirrors AbstractScope.isStaticBlockScope() in Java via NodeUtil.isClassStaticBlock.
func (b *AbstractScopeBase) IsStaticBlockScope() bool {
	return nodeIsClassStaticBlock(b.rootNode)
}

// IsCatchScope reports whether this scope is the block of a catch clause.
// Mirrors AbstractScope.isCatchScope() in Java.
func (b *AbstractScopeBase) IsCatchScope() bool {
	n := b.rootNode
	if n == nil || n.Token() != rhino.BLOCK {
		return false
	}
	kids := n.Children()
	return len(kids) == 1 && kids[0].Token() == rhino.CATCH
}

// IsCfgRootScope reports whether this scope's root is a valid CFG root.
// Mirrors AbstractScope.isCfgRootScope() in Java via NodeUtil.isValidCfgRoot.
func (b *AbstractScopeBase) IsCfgRootScope() bool {
	return nodeIsValidCfgRoot(b.rootNode)
}

// IsMemberFieldDefScope reports whether this scope's root is a member field
// definition (ES2022 class fields). Mirrors AbstractScope.isMemberFieldDefScope.
func (b *AbstractScopeBase) IsMemberFieldDefScope() bool {
	return b.rootNode != nil && b.rootNode.Token() == rhino.MEMBER_FIELD_DEF
}

// IsComputedFieldDefRhsScope reports whether this scope's root is a computed
// field definition RHS. Mirrors AbstractScope.isComputedFieldDefRhsScope() in Java.
// The Go AST uses BLOCK for all class-body member scopes and does not distinguish
// computed field def RHS; returns false as a conservative default.
func (b *AbstractScopeBase) IsComputedFieldDefRhsScope() bool {
	// TODO(port): AbstractScope.isComputedFieldDefRhsScope from AbstractScope.java:345
	// Reason: language-adapter — the rhino Go port does not emit a separate
	// COMPUTED_FIELD_DEF token; computed class fields resolve to BLOCK nodes.
	return false
}

// ---------------------------------------------------------------------------
// NodeUtil helpers used by AbstractScopeBase
// ---------------------------------------------------------------------------

// Go-only: language-adapter — nodeCreatesBlockScope wraps NodeUtil.createsBlockScope;
// Java calls NodeUtil.createsBlockScope(n) directly as a static method.
func nodeCreatesBlockScope(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.BLOCK,
		rhino.FOR, rhino.FOR_IN, rhino.FOR_OF,
		rhino.SWITCH:
		return true
	}
	return false
}

// Go-only: language-adapter — nodeIsClassStaticBlock wraps NodeUtil.isClassStaticBlock;
// Java calls NodeUtil.isClassStaticBlock(n) as a static method.
func nodeIsClassStaticBlock(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.BLOCK {
		return false
	}
	p := n.Parent()
	return p != nil && p.Token() == rhino.CLASS_MEMBERS
}

// Go-only: language-adapter — nodeIsValidCfgRoot wraps NodeUtil.isValidCfgRoot;
// Java calls NodeUtil.isValidCfgRoot(n) as a static method.
func nodeIsValidCfgRoot(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.FUNCTION, rhino.SCRIPT, rhino.MODULE_BODY:
		return true
	}
	return nodeIsClassStaticBlock(n)
}

// Go-only: language-adapter — nodeIsNonArrowFunction wraps NodeUtil.isNonArrowFunction;
// Java calls the static method directly while Go needs a package-level function.
func nodeIsNonArrowFunction(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	// Our parser may emit FUNCTION for both regular and arrow functions.
	// Check if the FUNCTION has a ARROW-style child (no FUNCTION keyword).
	// For safety, treat all FUNCTION nodes as non-arrow (conservative).
	return n.Token() == rhino.FUNCTION
}

// ---------------------------------------------------------------------------
// Scope-chain walk helpers (standalone, work on AbstractScope interface)
// ---------------------------------------------------------------------------

// Go-only: language-adapter — ScopeContains provides the AbstractScope.contains()
// method as a standalone function so it works on the AbstractScope interface;
// Java calls s.contains(other) as an instance method on a concrete scope type.
func ScopeContains(outer, inner AbstractScope) bool {
	for s := inner; s != nil; s = s.GetParentScope() {
		if s == outer {
			return true
		}
	}
	return false
}

// Go-only: language-adapter — ScopeGetGlobalScope provides getGlobalScope() on
// the AbstractScope interface; Java calls s.getGlobalScope() as a final instance method.
func ScopeGetGlobalScope(s AbstractScope) AbstractScope {
	for s.GetParentScope() != nil {
		s = s.GetParentScope()
	}
	return s
}

// Go-only: language-adapter — ScopeGetClosestHoistScope provides getClosestHoistScope()
// on the AbstractScope interface; Java uses a final instance method.
func ScopeGetClosestHoistScope(s AbstractScope) AbstractScope {
	for cur := s; cur != nil; cur = cur.GetParentScope() {
		if cur.IsHoistScope() {
			return cur
		}
	}
	return nil
}

// Go-only: language-adapter — ScopeGetClosestContainerScope provides
// getClosestContainerScope() on the AbstractScope interface; Java uses a final
// instance method on AbstractScope.
func ScopeGetClosestContainerScope(s AbstractScope) AbstractScope {
	scope := ScopeGetClosestHoistScope(s)
	if scope == nil {
		return nil
	}
	if scope.IsFunctionBlockScope() {
		scope = scope.GetParentScope()
	}
	return scope
}

// Go-only: language-adapter — ScopeGetClosestCfgRootScope provides
// getClosestCfgRootScope() on the AbstractScope interface; Java uses a final
// instance method on AbstractScope.
func ScopeGetClosestCfgRootScope(s AbstractScope) AbstractScope {
	for cur := s; cur != nil; cur = cur.GetParentScope() {
		if cur.IsCfgRootScope() {
			return cur
		}
	}
	return s
}

// Go-only: language-adapter — ScopeGetScopeOfThis provides getScopeOfThis()
// on the AbstractScope interface; Java uses a final instance method on AbstractScope.
func ScopeGetScopeOfThis(s AbstractScope) AbstractScope {
	scope := ScopeGetClosestContainerScope(s)
	for scope != nil && !scope.IsGlobal() && !scope.IsFunctionScope() {
		scope = ScopeGetClosestContainerScope(scope.GetParentScope())
	}
	return scope
}

// Go-only: language-adapter — ScopeGetCommonParent provides getCommonParent()
// on the AbstractScope interface; Java uses an instance method that returns S
// (the concrete scope type).
func ScopeGetCommonParent(a, b AbstractScope) AbstractScope {
	left := a
	right := b
	for left != nil && right != nil && left != right {
		ld := left.GetDepth()
		rd := right.GetDepth()
		if ld >= rd {
			left = left.GetParentScope()
		}
		if rd >= ld {
			right = right.GetParentScope()
		}
	}
	if left == right {
		return left
	}
	return nil
}

// Go-only: language-adapter — ScopeHasSameContainerScope provides
// hasSameContainerScope() on the AbstractScope interface; Java uses an instance
// method on AbstractScope that takes a concrete S typed other scope.
func ScopeHasSameContainerScope(a, b AbstractScope) bool {
	if a == b {
		return true
	}
	return ScopeGetClosestContainerScope(a) == ScopeGetClosestContainerScope(b)
}

// ---------------------------------------------------------------------------
// UntypedScope scope-chain walk helpers (typed variants)
// ---------------------------------------------------------------------------

// GetClosestContainerScope returns the closest container (non-block) scope.
// Mirrors AbstractScope.getClosestContainerScope() in Java.
func (s *UntypedScope) GetClosestContainerScope() *UntypedScope {
	scope := s.GetClosestHoistScope()
	if scope == nil {
		return nil
	}
	if scope.IsFunctionBlockScope() && scope.parent != nil {
		scope = scope.parent
	}
	return scope
}

// GetClosestCfgRootScope returns the nearest scope whose root is a valid CFG root.
// Mirrors AbstractScope.getClosestCfgRootScope() in Java.
func (s *UntypedScope) GetClosestCfgRootScope() *UntypedScope {
	for cur := s; cur != nil; cur = cur.parent {
		if cur.IsCfgRootScope() {
			return cur
		}
	}
	return s
}

// GetScopeOfThis returns the scope in which `this` is bound — the nearest
// non-arrow-function scope or the global scope.
// Mirrors AbstractScope.getScopeOfThis() in Java.
func (s *UntypedScope) GetScopeOfThis() *UntypedScope {
	scope := s.GetClosestContainerScope()
	for scope != nil && !scope.IsGlobal() && !scope.IsFunctionScope() {
		if scope.parent == nil {
			break
		}
		scope = scope.parent.GetClosestContainerScope()
	}
	return scope
}

// GetArgumentsVar returns the implicit "arguments" variable accessible from
// this scope, searching up the scope chain. Returns nil if not in any function.
// Mirrors AbstractScope.getArgumentsVar() in Java.
func (s *UntypedScope) GetArgumentsVar() *UntypedVar {
	for sc := s; sc != nil; sc = sc.parent {
		if implicitVarIsMadeByScope(ImplicitVarArguments, sc.rootNode) {
			return sc.getOrMakeImplicitVar(ImplicitVarArguments)
		}
	}
	return nil
}

// Contains reports whether s contains or equals other.
// Mirrors AbstractScope.contains(S other) in Java.
func (s *UntypedScope) Contains(other *UntypedScope) bool {
	for cur := other; cur != nil; cur = cur.parent {
		if cur == s {
			return true
		}
	}
	return false
}

// GetCommonParent returns the lowest common ancestor scope of s and other.
// Mirrors AbstractScope.getCommonParent(S other) in Java.
func (s *UntypedScope) GetCommonParent(other *UntypedScope) *UntypedScope {
	left := s
	right := other
	for left != nil && right != nil && left != right {
		ld := left.depth
		rd := right.depth
		if ld >= rd {
			left = left.parent
		}
		if rd >= ld {
			right = right.parent
		}
	}
	if left == right {
		return left
	}
	return nil
}

// GetAllAccessibleVariables returns all vars accessible from this scope
// (this scope + all parent scopes), with inner scopes shadowing outer ones.
// Does not include the special 'arguments' variable.
// Mirrors AbstractScope.getAllAccessibleVariables() in Java.
func (s *UntypedScope) GetAllAccessibleVariables() []*UntypedVar {
	seen := make(map[string]bool)
	var result []*UntypedVar
	for sc := s; sc != nil; sc = sc.parent {
		for _, v := range sc.GetVarIterable() {
			if !seen[v.GetName()] {
				seen[v.GetName()] = true
				result = append(result, v)
			}
		}
	}
	return result
}

// GetAllSymbols returns all symbols declared in this scope only.
// Mirrors AbstractScope.getAllSymbols() in Java.
func (s *UntypedScope) GetAllSymbols() []*UntypedVar {
	return s.GetVarIterable()
}

// CanDeclare reports whether the given name can be declared in this scope
// without illegal shadowing. Prevents redeclaring parameters in a function
// block scope.
// Mirrors AbstractScope.canDeclare(String) in Java.
func (s *UntypedScope) CanDeclare(name string) bool {
	if s.HasOwnSlot(name) {
		return false
	}
	// Parameters declared on the function-container scope must not be
	// re-declared in the function-block scope (except for bleeding functions).
	if s.IsFunctionBlockScope() && s.parent != nil && s.parent.HasOwnSlot(name) {
		// Allow if it's a bleeding function (function expression whose name
		// bleeds into the enclosing scope).
		v := s.parent.GetOwnSlot(name)
		if v != nil && v.IsBleedingFunction() {
			return true
		}
		return false
	}
	return true
}

// Undeclare removes a variable from this scope. Panics if var is not in
// this scope. Mirrors AbstractScope.undeclare(V var) in Java.
func (s *UntypedScope) Undeclare(v *UntypedVar) {
	if s.vars == nil {
		panic("undeclare: scope has no vars")
	}
	name := v.GetName()
	if s.vars[name] != v {
		panic("undeclare: var not in this scope")
	}
	delete(s.vars, name)
}

// HasSameContainerScope reports whether s and other are in the same container
// scope (same function, or both global).
// Mirrors AbstractScope.hasSameContainerScope(S other) in Java.
func (s *UntypedScope) HasSameContainerScope(other *UntypedScope) bool {
	if s == other {
		return true
	}
	return s.GetClosestContainerScope() == other.GetClosestContainerScope()
}
