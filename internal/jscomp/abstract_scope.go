// Port of com/google/javascript/jscomp/AbstractScope.java @ closure-compiler v20261028

/*
 * Copyright 2018 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"hash/maphash"
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// How the scope hierarchy is mapped (read this before porting code that uses
// scopes or vars).
//
// Java has AbstractScope<S extends AbstractScope<S, V>, V extends
// AbstractVar<S, V>> with two subclasses, Scope (S=Scope, V=Var) and
// TypedScope (S=TypedScope, V=TypedVar), and AbstractVar<S, V> with Var and
// TypedVar. Go-only: language-adapter — Go has neither wildcards nor
// covariant overrides, so:
//
//   - AbstractScope (this interface) is the wildcard type AbstractScope<?, ?>
//     (what NodeTraversal, ScopeCreator and MemoizedScopeCreator handle). It
//     declares every AbstractScope method whose signature mentions neither S
//     nor V.
//   - scopeCommon holds the root node and flags shared by Scope and TypedScope;
//     methods whose signature mentions the concrete scope or var type
//     (GetParent, GetVar, GetGlobalScope, GetClosestHoistScope,
//     GetVarIterable, ...) live directly on *Scope and *TypedScope without
//     F-bounded generic indirection.
//   - varCommon holds the shared fields and methods of Var and TypedVar
//     directly (no AbstractVar interface in production code).
//   - Java narrows three rhino.StaticScope methods to S/V: getParentScope(),
//     getSlot(String) and getOwnSlot(String). They keep rhino.StaticScope's
//     Go signatures (as rhino.StaticTypedScope does), so that *Scope is a
//     rhino.StaticScope and *TypedScope a rhino.StaticTypedScope; they never
//     return a typed nil. Callers that need the Java type assert, with the
//     comma-ok form so that Java null stays nil:
//       v, _ := scope.GetOwnSlot(name).(*Var)
//       tv, _ := typedScope.GetSlot(name).(*TypedVar)
//     The same holds for Var/TypedVar's getScope() (rhino.StaticSlot),
//     getDeclaration() and getSymbol(): `v.GetScope().(*Scope)`.

// AbstractScope contains information about a variable scope in JavaScript.
// Scopes can be nested, a scope points back to its parent scope. A Scope
// contains information about variables defined in that scope.
//
// ES 2015 introduces new scoping rules, which adds some complexity to this
// class. In particular, scopes fall into two mutually exclusive categories:
// block and container. Block scopes are all scopes whose roots are blocks,
// as well as any control structures whose optional blocks are omitted. These
// scopes did not exist at all prior to ES 2015. Container scopes comprise
// function scopes, global scopes, and module scopes, and (aside from modules,
// which didn't exist in ES5) corresponds to the ES5 scope rules. This
// corresponds roughly to one container scope per CFG root (but not exactly,
// due to SCRIPT-level CFGs).
//
// All container scopes, plus the outermost block scope within a function
// (i.e. the function block scope) are considered hoist scopes. All functions
// thus have two hoist scopes: the function scope and the function block
// scope. Hoist scopes are relevant because "var" declarations are hoisted to
// the closest hoist scope, as opposed to "let" and "const" which always apply
// to the specific scope in which they occur.
//
// Note that every function actually has two distinct hoist scopes: a
// container scope on the FUNCTION node, and a block-scope on the top-level
// BLOCK in the function (the "function block"). Local variables are declared
// on the function block, while parameters and optionally the function name
// (if it bleeds, i.e. from a named function expression) are declared on the
// container scope. This is required so that default parameter initializers
// can refer to names from outside the function that could possibly be
// shadowed in the function block. But these scopes are not fully independent
// of one another, since the language does not allow a top-level local
// variable to shadow a parameter name - so in some situations these scopes
// must be treated as a single scope.
//
// See NodeTraversal.
type AbstractScope interface {
	rhino.StaticScope

	// Java: AbstractScope.hasSlot(String)
	HasSlot(name string) bool
}

// scopeCommon holds the root node, depth, and state flags shared by Scope and
// TypedScope (12 bytes, 4-byte aligned).
type scopeCommon struct {
	rootNode        rhino.Node
	depth           int32
	varsInitialized bool
	isBottom        bool
	hasImplicitVars bool
	scanDone        bool
}

func (b *scopeCommon) initScopeCommon(rootNode rhino.Node) {
	*b = scopeCommon{rootNode: rhino.CheckNotNull(rootNode)}
}

func (b *scopeCommon) common() *scopeCommon { return b }

// Java: AbstractScope.getDepth()
func (b *scopeCommon) GetDepth() int {
	return int(b.depth)
}

// Java: AbstractScope.toString()
func (b *scopeCommon) String() string {
	return "Scope@" + jutil.ToJavaString(b.rootNode)
}

// GetRootNode gets the container node of the scope. This is typically the
// FUNCTION node or the global BLOCK/SCRIPT node.
//
// Java: AbstractScope.getRootNode()
func (b *scopeCommon) GetRootNode() rhino.Node {
	return b.rootNode
}

// Java: AbstractScope.isBlockScope()
func (b *scopeCommon) IsBlockScope() bool {
	return NodeUtil_CreatesBlockScope(b.rootNode)
}

// Java: AbstractScope.isStaticBlockScope()
func (b *scopeCommon) IsStaticBlockScope() bool {
	return NodeUtil_IsClassStaticBlock(b.rootNode)
}

// Java: AbstractScope.isFunctionBlockScope()
func (b *scopeCommon) IsFunctionBlockScope() bool {
	return NodeUtil_IsFunctionBlock(b.rootNode)
}

// Java: AbstractScope.isFunctionScope()
func (b *scopeCommon) IsFunctionScope() bool {
	return b.rootNode.IsFunction()
}

// Java: AbstractScope.isModuleScope()
func (b *scopeCommon) IsModuleScope() bool {
	return b.rootNode.IsModuleBody()
}

// Java: AbstractScope.isMemberFieldDefScope()
func (b *scopeCommon) IsMemberFieldDefScope() bool {
	return b.rootNode.IsMemberFieldDef()
}

// Java: AbstractScope.isComputedFieldDefRhsScope()
func (b *scopeCommon) IsComputedFieldDefRhsScope() bool {
	return b.rootNode.IsComputedFieldDef()
}

// Java: AbstractScope.isCatchScope()
func (b *scopeCommon) IsCatchScope() bool {
	return b.rootNode.IsBlock() &&
		b.rootNode.HasOneChild() &&
		b.rootNode.GetFirstChild().IsCatch()
}

// Java: AbstractScope.isCfgRootScope()
func (b *scopeCommon) IsCfgRootScope() bool {
	return NodeUtil_IsValidCfgRoot(b.rootNode)
}

// Java: AbstractScope.isGlobal()
func (b *scopeCommon) IsGlobal() bool {
	return b.depth == 0
}

// Java: AbstractScope.isLocal()
func (b *scopeCommon) IsLocal() bool {
	return b.depth != 0
}

// Java: AbstractScope.isHoistScope()
func (b *scopeCommon) IsHoistScope() bool {
	return b.depth == 0 ||
		b.rootNode.IsFunction() ||
		NodeUtil_IsFunctionBlock(b.rootNode) ||
		b.rootNode.IsModuleBody() ||
		NodeUtil_IsClassStaticBlock(b.rootNode)
}

var scopeHashSeed = maphash.MakeSeed()

func scopeNameHash(s string) uint32 {
	if s == "" {
		return 0
	}
	return uint32(maphash.String(scopeHashSeed, s))
}

func scopeBloomFromHash(h uint32) uint64 {
	return (uint64(1) << (h & 63)) | (uint64(1) << ((h >> 6) & 63))
}

func scopeIDHash(id uint32) uint32 {
	h := id * 0x9e3779b1
	h ^= h >> 16
	h *= 0x85ebca6b
	return h ^ (h >> 13)
}

func scopeBloomFromID(id uint32) uint64 {
	return scopeBloomFromHash(scopeIDHash(id))
}

const (
	scopeSmallVars      = 8
	scopeIndexTagMask   = 0xFF000000
	scopeIndexPosMask   = 0x00FFFFFF
	scopeIndexMinLength = 32
)

// scopeVarIndex holds the open-addressing hash index and 2-bit-per-name Bloom
// filter shared by Scope and TypedScope (16 bytes: pointer + uint64).
//
// Each non-zero table slot stores (h & scopeIndexTagMask) | uint32(pos + 1),
// keeping the table pointer-free for the GC and at most half full.
type scopeVarIndex struct {
	varIndex *[]uint32
	varBloom uint64
}

func (idx *scopeVarIndex) insertSlot(pos int, h uint32) {
	table := *idx.varIndex
	mask := uint32(len(table) - 1)
	slot := h & mask
	for table[slot] != 0 {
		slot = (slot + 1) & mask
	}
	table[slot] = h&scopeIndexTagMask | uint32(pos+1)
}

func (idx *scopeVarIndex) allocTable(nVars int) bool {
	if nVars <= scopeSmallVars {
		idx.varIndex = nil
		return false
	}
	n := scopeIndexMinLength
	for n < 4*nVars {
		n *= 2
	}
	table := make([]uint32, n)
	idx.varIndex = &table
	return true
}

func (idx *scopeVarIndex) recordAdded(nVars int, h uint32, bloom uint64) bool {
	idx.varBloom |= bloom
	if idx.varIndex != nil && nVars*2 <= len(*idx.varIndex) {
		idx.insertSlot(nVars-1, h)
		return false
	}
	return nVars > scopeSmallVars
}

func scopeRebuildIndex[V any](idx *scopeVarIndex, vars []V, hashOf func(V) uint32) {
	if !idx.allocTable(len(vars)) {
		return
	}
	for i, v := range vars {
		idx.insertSlot(i, hashOf(v))
	}
}

func scopeAppendVar[V any](
	c *scopeCommon, idx *scopeVarIndex, vars *[]V, var_ V, h uint32, bloom uint64, hashOf func(V) uint32,
) {
	c.varsInitialized = true
	*vars = append(*vars, var_)
	if idx.recordAdded(len(*vars), h, bloom) {
		scopeRebuildIndex(idx, *vars, hashOf)
	}
}

func scopeUndeclare[V jutil.Equaler](
	var_ V, c *scopeCommon, idx *scopeVarIndex, vars *[]V, pos int, hashOf func(V) uint32,
) {
	if pos < 0 {
		panic(jutil.NewNullPointerException())
	}
	jutil.CheckState((*vars)[pos].Equals(var_))
	if !c.varsInitialized {
		panic(jutil.NewUnsupportedOperationException())
	}
	*vars = slices.Delete(*vars, pos, pos+1)
	if idx.varIndex != nil {
		scopeRebuildIndex(idx, *vars, hashOf)
	}
	if len(*vars) == 0 {
		idx.varBloom = 0
	}
}

func scopeVarSeq[V any](vars *[]V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := 0; i < len(*vars); i++ {
			if !yield((*vars)[i]) {
				return
			}
		}
	}
}

func scopeGetImplicitVar[S any, V comparable](
	s S, cache **[4]V, imp AbstractScope_ImplicitVar, makeVar func(S, AbstractScope_ImplicitVar) V,
) V {
	var zero V
	if *cache == nil {
		*cache = new([4]V)
	}
	if v := (*cache)[imp]; v != zero {
		return v
	}
	if v := makeVar(s, imp); v != zero {
		(*cache)[imp] = v
		return v
	}
	return zero
}

func isBleedingFunctionNode(nameNode rhino.Node) bool {
	return !nameNode.IsNil() && nameNode.GetParent().IsFunction()
}

//go:noinline
func failIllegalShadow(nameNode rhino.Node) {
	panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", nameNode)))
}

type scopeChain[S comparable] interface {
	comparable
	parentScope() S
	lookupParentScope() S
	common() *scopeCommon
}

func scopeGlobal[S scopeChain[S]](s S) S {
	var zero S
	for p := s.parentScope(); p != zero; p = s.parentScope() {
		s = p
	}
	return s
}

func scopeContains[S scopeChain[S]](s, other S) bool {
	var zero S
	if other == zero {
		panic(jutil.NewNullPointerException())
	}
	for cur := other; cur != zero; cur = cur.parentScope() {
		if cur == s {
			return true
		}
	}
	return false
}

func scopeClosestHoist[S scopeChain[S]](s S) S {
	var zero S
	for cur := s; cur != zero; cur = cur.parentScope() {
		if cur.common().IsHoistScope() {
			return cur
		}
	}
	return zero
}

func scopeClosestCfgRoot[S scopeChain[S]](s S) S {
	for cur := s; ; cur = cur.parentScope() {
		if cur.common().IsCfgRootScope() {
			return cur
		}
	}
}

func scopeClosestContainer[S scopeChain[S]](s S) S {
	scope := scopeClosestHoist(s)
	if NodeUtil_IsFunctionBlock(scope.common().rootNode) {
		scope = scope.parentScope()
		if scope.common().IsBlockScope() {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("%s", scope)))
		}
	}
	return scope
}

func scopeHasSameContainer[S scopeChain[S]](s, other S) bool {
	return s == other || scopeClosestContainer(s) == scopeClosestContainer(other)
}

func scopeCommonParent[S scopeChain[S]](s, other S) S {
	var zero S
	left, right := s, other
	for left != zero && right != zero && left != right {
		ld, rd := left.common().depth, right.common().depth
		if ld >= rd {
			left = left.parentScope()
		}
		if ld <= rd {
			right = right.parentScope()
		}
	}
	jutil.CheckState(left != zero && left == right)
	return left
}

func scopeOfThis[S scopeChain[S]](s S) S {
	var zero S
	scope := scopeClosestContainer(s)
	for p := scope.parentScope(); p != zero && !NodeUtil_IsNonArrowFunction(scope.common().rootNode); p = scope.parentScope() {
		scope = scopeClosestContainer(p)
	}
	return scope
}

func scopeArgumentsVar[S scopeChain[S], V comparable](
	s S, getImplicit func(S, AbstractScope_ImplicitVar) V,
) V {
	var zeroS S
	var zeroV V
	for cur := s; cur != zeroS; cur = cur.lookupParentScope() {
		if cur.common().hasImplicitVars {
			if v := getImplicit(cur, AbstractScope_ImplicitVar_ARGUMENTS); v != zeroV {
				return v
			}
		}
	}
	return zeroV
}

func scopeAllAccessibleVars[S scopeChain[S], V any, K comparable](
	s S, vars func(S) []V, key func(V) K,
) iter.Seq[V] {
	return func(yield func(V) bool) {
		var zero S
		var seen map[K]struct{}
		for cur := s; cur != zero; cur = cur.parentScope() {
			list := vars(cur)
			hasParent := cur.parentScope() != zero
			for _, v := range list {
				k := key(v)
				if _, dup := seen[k]; dup {
					continue
				}
				if hasParent || seen != nil {
					if seen == nil {
						seen = make(map[K]struct{}, len(list))
					}
					seen[k] = struct{}{}
				}
				if !yield(v) {
					return
				}
			}
		}
	}
}

func checkChildScope[S scopeChain[S]](b *scopeCommon, parent S) {
	var zero S
	if parent == zero {
		panic(jutil.NewNullPointerException())
	}
	if !NodeUtil_CreatesScope(b.rootNode) {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", b.rootNode)))
	}
	if b.rootNode == parent.common().rootNode {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.LenientFormat("rootNode should not be the parent's root node: %s", b.rootNode)))
	}
}

// CheckRootScope performs simple validity checks when constructing a root scope.
//
// Java: AbstractScope.checkRootScope()
func (b *scopeCommon) CheckRootScope() {
	// TODO(tbreisacher): Can we tighten this to just NodeUtil.createsScope?
	if !NodeUtil_CreatesScope(b.rootNode) && !b.rootNode.IsScript() && !b.rootNode.IsRoot() {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", b.rootNode)))
	}
}

// AbstractScope_ImplicitVar enumerates the implicit var kinds defined in
// scopes without an explicit declaration.
//
// Java: AbstractScope.ImplicitVar
type AbstractScope_ImplicitVar int32

const (
	AbstractScope_ImplicitVar_ARGUMENTS AbstractScope_ImplicitVar = iota
	AbstractScope_ImplicitVar_EXPORTS
	AbstractScope_ImplicitVar_SUPER
	// TODO(sdh): Expand THIS.isMadeByScope to check super.isMadeByScope(scope) || scope.isGlobal()
	// Currently this causes a number of problems (see b/74980936), but could eventually lead to
	// better type information.  We might also want to restrict this so that module-root scopes
	// explicitly *don't* have access to the global this, though I think this is more than just
	// returning false in isMadeByScope - rather, getVar() needs to stop checking and immediately
	// return null.
	AbstractScope_ImplicitVar_THIS
)

var abstractScope_ImplicitVar_names = [...]string{"ARGUMENTS", "EXPORTS", "SUPER", "THIS"}

// abstractScope_ImplicitVar_name is the field ImplicitVar.name.
var abstractScope_ImplicitVar_name = [...]string{"arguments", "exports", "super", "this"}

var (
	implicitArgumentsID = uint32(strpool.Intern("arguments"))
	implicitExportsID   = uint32(strpool.Intern("exports"))
	implicitSuperID     = uint32(strpool.Intern("super"))
	implicitThisID      = uint32(strpool.Intern("this"))

	abstractScope_ImplicitVar_id = [...]uint32{
		implicitArgumentsID,
		implicitExportsID,
		implicitSuperID,
		implicitThisID,
	}
)

// String is Enum.name().
func (v AbstractScope_ImplicitVar) String() string { return abstractScope_ImplicitVar_names[v] }

// name is the field ImplicitVar.name.
func (v AbstractScope_ImplicitVar) name() string { return abstractScope_ImplicitVar_name[v] }

func (v AbstractScope_ImplicitVar) nameID() uint32 { return abstractScope_ImplicitVar_id[v] }

func (v AbstractScope_ImplicitVar) isMadeByRootNode(root rhino.Node) bool {
	switch v {
	case AbstractScope_ImplicitVar_EXPORTS:
		return root.IsModuleBody() &&
			root.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE)
	case AbstractScope_ImplicitVar_SUPER, AbstractScope_ImplicitVar_THIS:
		return NodeUtil_IsClassStaticBlock(root) ||
			NodeUtil_IsNonArrowFunction(root) ||
			root.IsMemberFieldDef() ||
			root.IsComputedFieldDef()
	case AbstractScope_ImplicitVar_ARGUMENTS:
		return NodeUtil_IsNonArrowFunction(root)
	}
	panic(jutil.NewIllegalStateException())
}

// IsMadeByScope reports whether this kind of implicit variable is
// created/owned by the given scope.
//
// Java: AbstractScope.ImplicitVar.isMadeByScope(AbstractScope<?, ?>)
func (v AbstractScope_ImplicitVar) IsMadeByScope(scope AbstractScope) bool {
	return v.isMadeByRootNode(scope.GetRootNode())
}

// AbstractScope_ImplicitVar_Lookup classifies name as an ImplicitVar without
// pointer boxing.
func AbstractScope_ImplicitVar_Lookup(name string) (AbstractScope_ImplicitVar, bool) {
	switch len(name) {
	case 4:
		if name == "this" {
			return AbstractScope_ImplicitVar_THIS, true
		}
	case 5:
		if name == "super" {
			return AbstractScope_ImplicitVar_SUPER, true
		}
	case 7:
		if name == "exports" {
			return AbstractScope_ImplicitVar_EXPORTS, true
		}
	case 9:
		if name == "arguments" {
			return AbstractScope_ImplicitVar_ARGUMENTS, true
		}
	}
	return 0, false
}

// AbstractScope_ImplicitVar_LookupID classifies an interned string ID as an
// ImplicitVar without pointer boxing.
func AbstractScope_ImplicitVar_LookupID(id uint32) (AbstractScope_ImplicitVar, bool) {
	switch id {
	case implicitArgumentsID:
		return AbstractScope_ImplicitVar_ARGUMENTS, true
	case implicitThisID:
		return AbstractScope_ImplicitVar_THIS, true
	case implicitExportsID:
		return AbstractScope_ImplicitVar_EXPORTS, true
	case implicitSuperID:
		return AbstractScope_ImplicitVar_SUPER, true
	}
	return 0, false
}
