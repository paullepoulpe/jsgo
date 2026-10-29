// Port of com/google/javascript/jscomp/AbstractScope.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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
//   - Java narrows three rhino.StaticScope methods to S/V: getParentScope(),
//     getSlot(String) and getOwnSlot(String). They keep rhino.StaticScope's
//     Go signatures (as rhino.StaticTypedScope does), so that *Scope is a
//     rhino.StaticScope and *TypedScope a rhino.StaticTypedScope; they never
//     return a typed nil. Callers that need the Java type assert, with the
//     comma-ok form so that Java null stays nil:
//       v, _ := scope.GetOwnSlot(name).(*Var)
//       tv, _ := typedScope.GetSlot(name).(*TypedVar)
//     The same holds for AbstractVar's getScope() (rhino.StaticSlot),
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

// Java: AbstractScope.getDepth()
func (b *scopeCommon) GetDepth() int {
	return int(b.depth)
}

// Java: AbstractScope.toString()
func (b *scopeCommon) String() string {
	return "Scope@" + jutil.ToJavaString(b.rootNode)
}

// Java: AbstractScope.untyped()
func (b *scopeCommon) Untyped() *Scope {
	panic(jutil.NewIllegalStateExceptionMsg("untyped() called, but not an untyped scope."))
}

// Java: AbstractScope.typed()
func (b *scopeCommon) Typed() *TypedScope {
	panic(jutil.NewIllegalStateExceptionMsg("typed() called, but not a typed scope."))
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

func (b *scopeCommon) checkChildScopeRoot(parentRoot rhino.Node) {
	if !NodeUtil_CreatesScope(b.rootNode) {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", b.rootNode)))
	}
	if b.rootNode == parentRoot {
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
	implicitArgumentsID = rhino.RhinoStringPool_InternID("arguments")
	implicitExportsID   = rhino.RhinoStringPool_InternID("exports")
	implicitSuperID     = rhino.RhinoStringPool_InternID("super")
	implicitThisID      = rhino.RhinoStringPool_InternID("this")

	abstractScope_ImplicitVar_id = [...]uint32{
		implicitArgumentsID,
		implicitExportsID,
		implicitSuperID,
		implicitThisID,
	}
)

// Go-only: language-adapter — nullable *AbstractScope_ImplicitVar callers share
// static pointers into this table so no call allocates.
var abstractScope_ImplicitVar_ptrs = [...]AbstractScope_ImplicitVar{
	AbstractScope_ImplicitVar_ARGUMENTS, AbstractScope_ImplicitVar_EXPORTS,
	AbstractScope_ImplicitVar_SUPER, AbstractScope_ImplicitVar_THIS,
}

// String is Enum.name().
func (v AbstractScope_ImplicitVar) String() string { return abstractScope_ImplicitVar_names[v] }

// Ordinal is Enum.ordinal().
func (v AbstractScope_ImplicitVar) Ordinal() int { return int(v) }

// name is the field ImplicitVar.name.
func (v AbstractScope_ImplicitVar) name() string { return abstractScope_ImplicitVar_name[v] }

func (v AbstractScope_ImplicitVar) nameID() uint32 { return abstractScope_ImplicitVar_id[v] }

// ptr returns the canonical non-null *AbstractScope_ImplicitVar for v.
func (v AbstractScope_ImplicitVar) ptr() *AbstractScope_ImplicitVar {
	return &abstractScope_ImplicitVar_ptrs[v]
}

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

// Java: AbstractScope.ImplicitVar.of(String)
func AbstractScope_ImplicitVar_Of(name string) *AbstractScope_ImplicitVar {
	if imp, ok := AbstractScope_ImplicitVar_Lookup(name); ok {
		return imp.ptr()
	}
	return nil
}

func AbstractScope_ImplicitVar_OfID(id uint32) *AbstractScope_ImplicitVar {
	if imp, ok := AbstractScope_ImplicitVar_LookupID(id); ok {
		return imp.ptr()
	}
	return nil
}
