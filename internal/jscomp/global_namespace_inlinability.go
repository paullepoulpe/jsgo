// Port of com/google/javascript/jscomp/GlobalNamespace.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// GlobalNamespace_Inlinability is how much to inline a Name.
//
// The `Inlinability#INLINE_BUT_KEEP_DECLARATION_*` cass are really an
// indicator that something 'unsafe' is happening in order to not break
// CollapseProperties as badly. Sadly Inlinability #INLINE_COMPLETELY may
// *also* be unsafe.
type GlobalNamespace_Inlinability int32

const (
	GlobalNamespace_Inlinability_INLINE_COMPLETELY GlobalNamespace_Inlinability = iota
	GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
	GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY
	GlobalNamespace_Inlinability_DO_NOT_INLINE
)

var globalNamespace_Inlinability_names = [...]string{
	"INLINE_COMPLETELY",
	"INLINE_BUT_KEEP_DECLARATION",
	"INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY",
	"DO_NOT_INLINE",
}

// The enum fields, indexed by ordinal.
var (
	globalNamespace_Inlinability_shouldInlineUsages      = [...]bool{true, true, true, false}
	globalNamespace_Inlinability_shouldRemoveDeclaration = [...]bool{true, false, false, false}
	globalNamespace_Inlinability_canCollapse             = [...]bool{true, true, true, false}
)

// Java: GlobalNamespace.Inlinability.name()
func (i GlobalNamespace_Inlinability) String() string { return globalNamespace_Inlinability_names[i] }

// Java: GlobalNamespace.Inlinability.ordinal()
func (i GlobalNamespace_Inlinability) Ordinal() int { return int(i) }

// Java: GlobalNamespace.Inlinability.shouldInlineUsages()
func (i GlobalNamespace_Inlinability) ShouldInlineUsages() bool {
	return globalNamespace_Inlinability_shouldInlineUsages[i]
}

// Java: GlobalNamespace.Inlinability.shouldRemoveDeclaration()
func (i GlobalNamespace_Inlinability) ShouldRemoveDeclaration() bool {
	return globalNamespace_Inlinability_shouldRemoveDeclaration[i]
}

// Java: GlobalNamespace.Inlinability.canCollapse()
func (i GlobalNamespace_Inlinability) CanCollapse() bool {
	return globalNamespace_Inlinability_canCollapse[i]
}

// CalculateInlinability returns whether to treat this alias as completely
// inlineable or to keep the aliasing assignment
//
// This method used to only return true/false, but now returns an enum in
// order to track more information about "unsafely" inlineable names.
//
// CollapseProperties will flatten `@constructor` properties even if they are
// potentially accessed by a reference other than their fully qualified name,
// which breaks those other refs. To avoid breakages AggressiveInlineAliases
// must unsafely inline constructor properties that alias another global
// name. Existing code depends on this behavior, and it's not easily
// determinable where these dependencies are.
//
// However, AggressiveInlineAliases must not also remove the initializtion of
// an alias if it is not safely inlineable. (i.e. if
// Inlinability#shouldRemoveDeclaration()). It's possible that a third name
// aliases the alias - we might later inline the third name (as an alias of
// the original alias) and don't want to set the third name to null.
//
// Java: GlobalNamespace.Name.calculateInlinability()
func (n *GlobalNamespace_Name) CalculateInlinability() GlobalNamespace_Inlinability {
	// Only simple aliases with direct usage are inlineable.
	// Exactly 2 global sets could be OK, if the first is the declaration with no initializer.
	if n.InExterns() || !n.HasOneRealGlobalSet() || n.localSets != 0 {
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	// TODO(lharker): consider separating canCollapseOrInline() into this method, since it
	// duplicates some logic here
	inlinability := n.CanCollapseOrInline()
	if !inlinability.ShouldInlineUsages() {
		// if you can't even inline the usages, do nothing.
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}
	var initialChunk *JSChunk
	if initRef := n.GetInitialization(); initRef != nil {
		initialChunk = initRef.GetChunk()
	} else {
		initialChunk = n.getDeclarationRef().GetChunk()
	}

	// Only allow inlining of simple references.
	for ref := range n.Refs() {
		switch ref.type_ {
		case GlobalNamespace_Ref_Type_SET_FROM_GLOBAL, GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL:
			// Expect one global set
			jutil.CheckStateMsg1(
				ref.IsUninitializedDeclaration() || ref.GetChunk() == initialChunk, "%s", ref)
			continue
		case GlobalNamespace_Ref_Type_SET_FROM_LOCAL, GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL:
			panic(jutil.NewIllegalStateException())
		case GlobalNamespace_Ref_Type_ALIASING_GET,
			GlobalNamespace_Ref_Type_DIRECT_GET,
			GlobalNamespace_Ref_Type_PROTOTYPE_GET,
			GlobalNamespace_Ref_Type_CALL_GET,
			GlobalNamespace_Ref_Type_SUBCLASSING_GET:
			// This name has a reference in a different chunk that is not guaranteed to be loaded
			// before the alias is initialized. In an ideal world, we might just back off inlining
			// this alias entirely. Unfortunately, that causes new runtime errors in practice,
			// because property collapsing is intentionally unsafe and will collapse properties even
			// when it will break code, and so it's safer to still inling usages of this alias
			// except for the specific invalid cross-chunk dependency(s).
			refChunk := ref.GetChunk()
			if inlinability != GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY &&
				refChunk != initialChunk &&
				!n.outer.compiler.GetChunkGraph().DependsOn(refChunk, initialChunk) {
				inlinability = GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY
			}
			continue
		case GlobalNamespace_Ref_Type_DELETE_PROP:
			return GlobalNamespace_Inlinability_DO_NOT_INLINE
		}
		panic(jutil.NewAssertionErrorNoMessage())
	}
	return inlinability
}

// Java: GlobalNamespace.Name.hasOneRealGlobalSet()
func (n *GlobalNamespace_Name) HasOneRealGlobalSet() bool {
	return n.globalSets == 1 || (n.globalSets == 2 && n.getDeclarationRef().IsUninitializedDeclaration())
}

// Java: GlobalNamespace.Name.canCollapse()
func (n *GlobalNamespace_Name) CanCollapse() bool {
	return n.CanCollapseOrInline().CanCollapse()
}

// CanCollapseOrInline determines whether it's safe to collapse properties on
// an objects
//
// For legacy reasons, both CollapseProperties and AggressiveInlineAliases
// share the same logic when deciding whether to inline properties or to
// collapse them.
//
// The main reasons we cannot inline/collapse properties of a name name are:
//
//	a) it is set multiple times or set once in a local scope,
//	b) one or more of the above conditions in canCollapseOrInlineChildNames
//	   applies to the namespace it's on,
//	c) it's annotated at-nocollapse
//	d) it's in the externs,
//	e) or it's a known getter or setter, not a regular property
//	f) it's an ES6 class static method that references `super` or the internal class name
//
// We ignore conditions (a) and (b) on at-constructor and at-enum names in
// CollapseProperties.
//
// In AggressiveInlineAliases we want to do some partial backoff if (a) and
// (b) are false for at-constructor or at-enum names, which is why we return
// an enum value instead of a boolean.
//
// Java: GlobalNamespace.Name.canCollapseOrInline()
func (n *GlobalNamespace_Name) CanCollapseOrInline() GlobalNamespace_Inlinability {
	if n.InExterns() {
		// condition (d)
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "declared in externs")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}
	if n.IsGetOrSetDefinition() {
		// condition (e)
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "getter / setter")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}
	if n.IsCollapsingExplicitlyDenied() {
		// condition (c)
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "@nocollapse")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.ReferencesSuperOrInnerClassName() {
		// condition (f)
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "references super or inner class name")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.IsToStringValueOfInObjectLiteral() {
		n.LogDecision(
			GlobalNamespace_Inlinability_DO_NOT_INLINE,
			"references explicit definition of toString/valueOf functions used implicitly in the JS"+
				" language")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.deleteProps > 0 {
		// If we inline or collapse, then the delete operation will be incorrect.
		n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "delete operator is used on this property")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if declRef := n.getDeclarationRef(); declRef != nil {
		declaration := declRef.GetNode()
		if declaration.GetParent().IsObjectLit() {
			if n.outer.DeclarationHasFollowingObjectSpreadSibling(declaration) {
				// Case: `var x = {a: 0, ...b, c: 2}` where declaration is `a` but not `c`.
				// Following spreads may overwrite the declaration.
				n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "obj lit property followed by spread")
				return GlobalNamespace_Inlinability_DO_NOT_INLINE
			}
			// We may be in a deeply nested object literal like, `{a: {b: {c: 1}}}`, so find the
			// outermost object literal node and check its ancestor expressions to determine
			// whether it is used conditionally.
			child := n.outer.GetOutermostObjectLit(declaration)
			for p := child.GetParent(); !p.IsNil() && !p.GetParent().IsNil() && !NodeUtil_IsStatement(p); child, p = p, p.GetParent() {
				if (p.IsOr() && p.GetFirstChild() != child) ||
					(p.IsHook() && p.GetFirstChild() != child) ||
					(p.IsAnd() && p.GetFirstChild() != child) {
					// Case: `var x = y || {a: b}` or `var x = cond ? y : {a: b}` or `var x = cond && {a:
					// b}`.
					n.LogDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "conditional definition")
					return GlobalNamespace_Inlinability_DO_NOT_INLINE
				}
			}
		}
	}

	// condition (a)
	isUnchangedThroughFullName :=
		(n.globalSets > 0 || n.localSets > 0) && n.localSetsWithNoCollapse == 0 && n.deleteProps == 0
	// additional information about condition (b)
	parentInlinability := GlobalNamespace_Inlinability_INLINE_COMPLETELY
	if parent := n.GetParent(); parent != nil {
		parentInlinability = parent.CanCollapseOrInlineChildNames()
	}

	// if condition (a) or condition (b) is not true, but this is a declared name, we may need
	// to allow inlining usages of a variable but keep the declaration.
	switch parentInlinability {
	case GlobalNamespace_Inlinability_INLINE_COMPLETELY:
		if isUnchangedThroughFullName {
			n.LogDecision(
				GlobalNamespace_Inlinability_INLINE_COMPLETELY, "parent inlineable: unchanged through full name")
			return GlobalNamespace_Inlinability_INLINE_COMPLETELY
		} else {
			// maybe inline usages of this name, but only if a declared type. non-declared-types
			// just
			// back off and don't inline at all
			unsafeInlinablility := n.GetUnsafeInlinabilityBasedOnDeclaredType()
			n.LogDecision(unsafeInlinablility, "parent inlineable: changed through full name")
			return unsafeInlinablility
		}
	case GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION,
		GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY:
		// this is definitely not safe to completely inline/collapse of its parent
		// if it's a declared type, we should still partially inline it and completely collapse it
		// if not a declared type we should partially inline it iff the other conditions hold
		if n.IsDeclaredType() {
			unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
			n.LogDecision(unsafeInlinability, "parent unsafely inlineable & is declared type")
			return unsafeInlinability
		} else if isUnchangedThroughFullName {
			n.LogDecision(
				parentInlinability, "parent unsafely inlineable & unchanged through full name")
			return parentInlinability
		} else {
			// Not a declared type. We may still 'partially' inline it because it must be a property
			// on an @enum or @constructor, but only if it actually matches conditions (a) and (b)
			n.LogDecision(
				GlobalNamespace_Inlinability_DO_NOT_INLINE,
				"parent unsafely inlineable & changed through full name")
			return GlobalNamespace_Inlinability_DO_NOT_INLINE
		}
	case GlobalNamespace_Inlinability_DO_NOT_INLINE:
		// If the parent is unsafely to collapse/inline, we will still inline it if it's on
		// a declaredType (i.e. @constructor or @enum), but we propagate the information that
		// the parent is unsafe. If this is not a declared type, return DO_NOT_INLINE.
		unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
		n.LogDecision(unsafeInlinability, "parent cannot be inlined")
		return unsafeInlinability
	}
	panic(jutil.NewIllegalStateExceptionMsg("unknown enum value " + parentInlinability.String()))
}

// Java: GlobalNamespace.Name.logDecision(Inlinability, String)
func (n *GlobalNamespace_Name) LogDecision(inlinability GlobalNamespace_Inlinability, reason string) {
	if n.outer.decisionsLog.IsLogging() {
		n.outer.decisionsLog.LogFormat("%s: %s: %s", n.GetFullName(), inlinability.String(), reason)
	}
}

// Java: GlobalNamespace.Name.getUnsafeInlinabilityBasedOnDeclaredType()
func (n *GlobalNamespace_Name) GetUnsafeInlinabilityBasedOnDeclaredType() GlobalNamespace_Inlinability {
	if n.GetBooleanProperty(GlobalNamespace_NameProp_CONSTRUCTOR_TYPE) {
		return GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
	} else if n.GetBooleanProperty(GlobalNamespace_NameProp_INTERFACE_TYPE) {
		return GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
	} else if n.GetBooleanProperty(GlobalNamespace_NameProp_ENUM_TYPE) {
		return GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
	} else if n.GetBooleanProperty(GlobalNamespace_NameProp_NOT_A_TYPE) {
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}
	panic(jutil.NewIllegalStateExceptionMsg("name missing declaredType value: " + n.String()))
}

// ReferencesSuperOrInnerClassName examines ES6 class members for some syntax
// that blocks collapsing
//
// Specifically, this looks for super references and references to inner
// class names. These are unique to ES6 static class members so we don't
// need more general handling.
//
// TODO(b/122665204): also return false on `super` in an object lit method
//
// Java: GlobalNamespace.Name.referencesSuperOrInnerClassName()
func (n *GlobalNamespace_Name) ReferencesSuperOrInnerClassName() bool {
	ref := n.getDeclarationRef()
	if ref == nil {
		return false
	}
	member := ref.GetNode()
	if member.IsNil() || !(member.IsStaticMember() && member.GetParent().IsClassMembers()) {
		return false
	}
	// Get either the function body or the member field initializer (which may be null if there is
	// no initializer).
	var body rhino.Node
	if member.IsMemberFieldDef() {
		body = member.GetFirstChild()
	} else {
		body = NodeUtil_GetFunctionBody(member.GetFirstChild())
	}
	if !body.IsNil() && NodeUtil_ReferencesSuper(body) {
		return true
	}

	classNode := member.GetGrandparent()
	if NodeUtil_IsClassDeclaration(classNode) {
		return false // e.g. class C {}
	}

	innerNameNode := classNode.GetFirstChild()
	return !innerNameNode.IsEmpty() && // e.g. const C = class {};
		NodeUtil_IsNameReferenced(member, innerNameNode.GetString())
}

// Java: GlobalNamespace.Name.isSetInLoop()
func (n *GlobalNamespace_Name) IsSetInLoop() bool {
	ref := n.getDeclarationRef()
	if ref != nil {
		node := ref.GetNode()
		if !node.IsNil() {
			return NodeUtil_IsWithinLoop(node)
		}
	}
	return false
}

// Java: GlobalNamespace.Name.isGetOrSetDefinition()
func (n *GlobalNamespace_Name) IsGetOrSetDefinition() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_GET_SET)
}

// Java: GlobalNamespace.Name.canCollapseUnannotatedChildNames()
func (n *GlobalNamespace_Name) CanCollapseUnannotatedChildNames() bool {
	return n.CanCollapseOrInlineChildNames().CanCollapse()
}

// IsToStringValueOfInObjectLiteral: toString/valueOf are implicitly used as
// part of the JS language should not be collapsed.
//
// Java: GlobalNamespace.Name.isToStringValueOfInObjectLiteral()
func (n *GlobalNamespace_Name) IsToStringValueOfInObjectLiteral() bool {
	parent := n.GetParent()
	baseName := n.GetBaseName()
	return n.IsFunction() &&
		parent != nil &&
		parent.IsObjectLiteral() &&
		(baseName == "toString" || baseName == "valueOf")
}

// CanCollapseOrInlineChildNames returns whether to assume that child
// properties of this name are collapsible/inlineable
//
// For legacy reasons, both CollapseProperties and AggressiveInlineAliases
// share the same logic when deciding whether to inline properties or to
// collapse them.
//
// The main reasons we cannot inline/collapse properties of a name name are:
//
//	a) it is set multiple times
//	b) its properties might not be referred to by their full qname but on a different object
//	c) one or more of the above conditions applies to a parent name
//	d) it's annotated @nocollapse
//	e) it's in the externs.
//	f) it's assigned a value that supports being aliased (e.g. ctors, so their static
//	   properties can be accessed via `this`)
//
// However, in some cases for properties of `@constructor` or `@enum` names,
// we ignore some of these conditions in order to more aggressively collapse
// `@constructor`s used in goog.provide namespace chains.
//
// Java: GlobalNamespace.Name.canCollapseOrInlineChildNames()
func (n *GlobalNamespace_Name) CanCollapseOrInlineChildNames() GlobalNamespace_Inlinability {
	// condition (a) and (b)
	if n.GetBooleanProperty(GlobalNamespace_NameProp_OTHER_OBJECT) {
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "NameProp.OTHER_OBJECT")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	} else if n.IsGetOrSetDefinition() {
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "getter/setter")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	} else if n.globalSets != 1 {
		if n.outer.decisionsLog.IsLogging() {
			n.LogChildNamesDecisionSupplier(
				GlobalNamespace_Inlinability_DO_NOT_INLINE, func() string { return "set " + strconv.Itoa(int(n.globalSets)) + " times globally" })
		}
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	} else if n.localSets != 0 {
		if n.outer.decisionsLog.IsLogging() {
			n.LogChildNamesDecisionSupplier(
				GlobalNamespace_Inlinability_DO_NOT_INLINE, func() string { return "set " + strconv.Itoa(int(n.localSets)) + " times locally" })
		}
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	} else if n.deleteProps != 0 {
		if n.outer.decisionsLog.IsLogging() {
			n.LogChildNamesDecisionSupplier(
				GlobalNamespace_Inlinability_DO_NOT_INLINE,
				func() string { return "properties are deleted " + strconv.Itoa(int(n.deleteProps)) + " times" })
		}
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	// Don't try to collapse if the one global set is a twin reference.
	// We could theoretically handle this case in CollapseProperties, but
	// it's probably not worth the effort.
	declRef := jutil.CheckNotNull(n.getDeclarationRef())
	if declRef.IsTwin() {
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "twinned declaration")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.IsCollapsingExplicitlyDenied() {
		// condition (d)
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "@nocollapse")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.IsSetInLoop() {
		// condition (a)
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "set in a loop")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.UsesHasOwnProperty() {
		// condition (b)
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "hasOwnProperty() call exists")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	if n.ValueImplicitlySupportsAliasing() {
		// condition (f)
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_DO_NOT_INLINE, "value implicitly supports aliasing")
		return GlobalNamespace_Inlinability_DO_NOT_INLINE
	}

	// If this is a key of an aliased object literal, then it will be aliased
	// later. So we won't be able to collapse its properties.
	// condition (b)
	parent := n.GetParent()
	if parent != nil && parent.ShouldKeepKeys() {
		unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
		n.LogChildNamesDecision(unsafeInlinability, "parent.shouldKeepKeys()")
		return unsafeInlinability
	}

	// If this is aliased, then its properties can't be collapsed either. but we may do so anyway
	// if it's a declared type.
	// condition (b)
	if n.aliasingGets > 0 {
		unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
		if n.outer.decisionsLog.IsLogging() {
			n.LogChildNamesDecisionSupplier(
				unsafeInlinability, func() string { return strconv.Itoa(int(n.aliasingGets)) + " aliasing gets exist" })
		}
		return unsafeInlinability
	}

	if parent == nil {
		// this is completely safe to inline! yay
		n.LogChildNamesDecision(GlobalNamespace_Inlinability_INLINE_COMPLETELY, "no reason not to inline")
		return GlobalNamespace_Inlinability_INLINE_COMPLETELY
	}

	// Cases are:
	//  - parent is safe to completely inline. then same for this name
	//  - parent is unsafe but should still be inlined. then same for this name
	//  - parent is unsafe, should not be inlined at all. then return either DO_NOT_INLINE,
	//    or maybe unsafely inline if this is a ctor property
	parentInlinability := parent.CanCollapseOrInlineChildNames()
	if parentInlinability == GlobalNamespace_Inlinability_DO_NOT_INLINE {
		// the parent name is used in a way making this unsafe to inline, but we might want to
		// inline usages of this name
		unsafeInlinability := n.GetUnsafeInlinabilityBasedOnDeclaredType()
		n.LogChildNamesDecision(unsafeInlinability, "parent is not inlineable")
		return unsafeInlinability
	}
	n.LogChildNamesDecision(parentInlinability, "inherited from parent")
	return parentInlinability
}

// Java: GlobalNamespace.Name.logChildNamesDecision(Inlinability, String)
func (n *GlobalNamespace_Name) LogChildNamesDecision(inlinability GlobalNamespace_Inlinability, reason string) {
	if n.outer.decisionsLog.IsLogging() {
		n.outer.decisionsLog.LogFormat("%s: children: %s: %s", n.GetFullName(), inlinability.String(), reason)
	}
}

// Java: GlobalNamespace.Name.logChildNamesDecision(Inlinability, Supplier<String>)
func (n *GlobalNamespace_Name) LogChildNamesDecisionSupplier(inlinability GlobalNamespace_Inlinability, reasonSupplier func() string) {
	if n.outer.decisionsLog.IsLogging() {
		n.outer.decisionsLog.LogSupplier(
			func() string {
				return n.GetFullName() + ": children: " + inlinability.String() + ": " + reasonSupplier()
			})
	}
}

// Java: GlobalNamespace.Name.valueImplicitlySupportsAliasing()
func (n *GlobalNamespace_Name) ValueImplicitlySupportsAliasing() bool {
	if !n.outer.enableImplicitlyAliasedValues {
		return false
	}

	if n.GetBooleanProperty(GlobalNamespace_NameProp_CLASS) {
		// Properties on classes may be referenced via `this` in static methods.
		return true
	}
	if n.GetBooleanProperty(GlobalNamespace_NameProp_FUNCTION) {
		// We want ES5 ctors/interfaces to behave consistently with ES6 because:
		// - transpilation should not change behaviour
		// - updating code shouldn't be hindered by behaviour changes
		jsdoc := n.GetJSDocInfo()
		return jsdoc != nil && jsdoc.IsConstructorOrInterface()
	}
	return false
}

// ShouldKeepKeys reports whether this is an object literal that needs to
// keep its keys.
//
// Java: GlobalNamespace.Name.shouldKeepKeys()
func (n *GlobalNamespace_Name) ShouldKeepKeys() bool {
	return n.IsObjectLiteral() && (n.aliasingGets > 0 || n.IsCollapsingExplicitlyDenied())
}

// Java: GlobalNamespace.Name.needsToBeStubbed()
func (n *GlobalNamespace_Name) NeedsToBeStubbed() bool {
	return n.globalSets == 0 &&
		n.localSets > 0 &&
		n.localSetsWithNoCollapse == 0 &&
		!n.IsCollapsingExplicitlyDenied()
}

// Java: GlobalNamespace.Name.setDeclaredTypeKind(NameProp)
func (n *GlobalNamespace_Name) SetDeclaredTypeKind(declaredType GlobalNamespace_NameProp) {
	jutil.CheckArgumentMsg1(
		(declaredType.Bit()&GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK) != 0,
		"Unexpected NameProp for declaredType %s",
		declaredType)
	n.propertyBitSet =
		(n.propertyBitSet &^ GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK) | declaredType.Bit()
	if declaredType != GlobalNamespace_NameProp_NOT_A_TYPE {
		g := n.outer
		for ancestorID := n.parentID; ancestorID != 0; {
			ancestor := g.nameByID(ancestorID)
			ancestor.SetBooleanProperty(GlobalNamespace_NameProp_IS_DECLARED)
			ancestorID = ancestor.parentID
		}
	}
}

// Java: GlobalNamespace.Name.setNameType(NameProp)
func (n *GlobalNamespace_Name) SetNameType(type_ GlobalNamespace_NameProp) {
	jutil.CheckArgumentMsg1(
		(type_.Bit()&GlobalNamespace_NameProp_NAME_KIND_MASK) != 0, "Unexpected NameProp for nameType %s", type_)
	n.propertyBitSet = (n.propertyBitSet &^ GlobalNamespace_NameProp_NAME_KIND_MASK) | type_.Bit()
}

// Java: GlobalNamespace.Name.isDeclaredType()
func (n *GlobalNamespace_Name) IsDeclaredType() bool {
	jutil.CheckStateMsg1(
		(n.propertyBitSet&GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK) != 0, "%s", n.propertyBitSet)
	return !n.GetBooleanProperty(GlobalNamespace_NameProp_NOT_A_TYPE)
}

// Java: GlobalNamespace.Name.isConstructor()
func (n *GlobalNamespace_Name) IsConstructor() bool {
	declNode := n.getDeclarationRef().node
	rvalueNode := NodeUtil_GetRValueOfLValue(declNode)
	jsdoc := NodeUtil_GetBestJSDocInfo(declNode)
	return !rvalueNode.IsNil() &&
		rvalueNode.IsFunction() &&
		jsdoc != nil &&
		jsdoc.IsConstructor()
}

// IsNamespaceObjectLit determines whether this name is a prefix of at least
// one class or enum name. Because classes and enums are always collapsed,
// the namespace will have different properties in compiled code than in
// uncompiled code.
//
// For example, if foo.bar.DomHelper is a class, then foo and foo.bar are
// considered namespaces.
//
// Java: GlobalNamespace.Name.isNamespaceObjectLit()
func (n *GlobalNamespace_Name) IsNamespaceObjectLit() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_IS_DECLARED) && n.IsObjectLiteral()
}

// IsSimpleName determines whether this is a simple name (as opposed to a
// qualified name).
//
// Java: GlobalNamespace.Name.isSimpleName()
func (n *GlobalNamespace_Name) IsSimpleName() bool {
	return n.parentID == 0
}

// Java: GlobalNamespace.Name.getTypeDebugString()
func (n *GlobalNamespace_Name) GetTypeDebugString() string {
	if n.IsClass() {
		return "CLASS"
	} else if n.IsFunction() {
		return "FUNCTION"
	} else if n.IsObjectLiteral() {
		return "OBJECTLIT"
	} else if n.IsGetOrSetDefinition() {
		return "GET_SET"
	} else if n.GetBooleanProperty(GlobalNamespace_NameProp_OTHER_OBJECT) {
		return "OTHER"
	}
	panic(jutil.NewAssertionError("Missing NameProp for name kind in " + strconv.Itoa(int(n.propertyBitSet))))
}

// Java: GlobalNamespace.Name.toString()
func (n *GlobalNamespace_Name) String() string {
	return n.GetFullName() +
		" (" +
		n.GetTypeDebugString() +
		"): " +
		strings.Join([]string{
			"globalSets=" + strconv.Itoa(int(n.globalSets)),
			"localSets=" + strconv.Itoa(int(n.localSets)),
			"totalGets=" + strconv.Itoa(int(n.totalGets)),
			"aliasingGets=" + strconv.Itoa(int(n.aliasingGets)),
			"callGets=" + strconv.Itoa(int(n.callGets)),
			"subclassingGets=" + strconv.Itoa(int(n.subclassingGets)),
		}, ", ")
}

// Java: GlobalNamespace.Name.getJSDocInfo()
func (n *GlobalNamespace_Name) GetJSDocInfo() *rhino.JSDocInfo {
	// e.g.
	// /** @type {string} */ X.numProp;     // could be a declaration, but...
	// /** @type {number} */ X.numProp = 3; // assignment wins
	if n.firstDeclDocID != 0 {
		return n.outer.docByID(n.firstDeclDocID)
	}
	return n.outer.docByID(n.firstQnameDocID)
}

// GetDocInfoForDeclaration tries to get the doc info for a given
// declaration ref.
//
// Java: GlobalNamespace.Name.getDocInfoForDeclaration(Ref)
func (n *GlobalNamespace_Name) GetDocInfoForDeclaration(ref *GlobalNamespace_Ref) *rhino.JSDocInfo {
	if !ref.node.IsNil() {
		refParent := ref.node.GetParent()
		if refParent.IsNil() {
			// May happen when inlineAliases removes refs from the AST.
			return nil
		}
		switch refParent.GetToken() {
		case rhino.Token_FUNCTION, rhino.Token_ASSIGN, rhino.Token_CLASS:
			return refParent.GetJSDocInfo()
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			if ref.node == refParent.GetFirstChild() {
				return refParent.GetJSDocInfo()
			}
			return ref.node.GetJSDocInfo()
		case rhino.Token_OBJECTLIT, rhino.Token_CLASS_MEMBERS:
			return ref.node.GetJSDocInfo()
		default:
		}
	}

	return nil
}

// Java: GlobalNamespace.Name.isModuleExport()
func (n *GlobalNamespace_Name) IsModuleExport() bool {
	return n.GetBooleanProperty(GlobalNamespace_NameProp_IS_MODULE_PROP)
}

// GetOutermostObjectLit, given something like the `'c'` STRING_KEY node in
// `x = {a: {b: {c: 0}}};`, returns the Node for the outermost object
// literal.
//
// objLitKey must be the child of an OBJECT_LIT node. Returns the first
// ancestor that is an OBJECT_LIT and whose grandparent is not an
// OBJECT_LIT.
//
// Java: GlobalNamespace.getOutermostObjectLit(Node)
func (g *GlobalNamespace) GetOutermostObjectLit(objLitKey rhino.Node) rhino.Node {
	outermostObjectLit := objLitKey.GetParent()
	jutil.CheckStateMsg1(outermostObjectLit.IsObjectLit(), "%s", outermostObjectLit)
	for {
		objLitGrandparent := outermostObjectLit.GetGrandparent()
		if !objLitGrandparent.IsNil() && objLitGrandparent.IsObjectLit() {
			outermostObjectLit = objLitGrandparent
		} else {
			return outermostObjectLit
		}
	}
}

// GlobalNamespace_IsQnameDeclarationWithoutAssignment: true if the given
// Node is the GETPROP in a statement like `some.q.name;`
//
// Such do-nothing statements often have JSDoc on them and are intended to
// declare the qualified name.
//
// node is any Node, or even null.
//
// Java: GlobalNamespace.isQnameDeclarationWithoutAssignment(Node)
func GlobalNamespace_IsQnameDeclarationWithoutAssignment(node rhino.Node) bool {
	return !node.IsNil() && node.IsGetProp() && node.GetParent().IsExprResult()
}

// GlobalNamespace_NameProp is an enum of boolean properties of a Name, to be
// stored in one bit set per name.
//
// Non-private for @AutoValue. Code outside this class shouldn't need this
// enum though.
type GlobalNamespace_NameProp int32

const (
	// Increment the "index" passed to the constructor by 1 for each enum member
	GlobalNamespace_NameProp_IS_DECLARED              GlobalNamespace_NameProp = iota // index 0
	GlobalNamespace_NameProp_IS_MODULE_PROP                                           // index 1
	GlobalNamespace_NameProp_IS_PROVIDED                                              // index 2: If this name was in any goog.provide() calls.
	GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY                                 // index 3
	GlobalNamespace_NameProp_IS_EXTERN                                                // index 4: corresponds to SourceKind.EXTERN

	// Mutually exclusive properties indicating what kind of Closure type this Name is, if any.
	// Corresponds to DECLARED_TYPE_KIND_MASK.
	GlobalNamespace_NameProp_CONSTRUCTOR_TYPE // index 5: a non-interface class {} or a `/** @constructor */ function`
	GlobalNamespace_NameProp_INTERFACE_TYPE   // index 6: a `/** @interface */` or `/** @record */`
	GlobalNamespace_NameProp_ENUM_TYPE        // index 7: an `/** @enum */`
	GlobalNamespace_NameProp_NOT_A_TYPE       // index 8

	// Mutually exclusive properties indicating what kind of JavaScript entity this Name is.
	// Corresponds to NAME_KIND_MASK.
	GlobalNamespace_NameProp_CLASS        // index 9: class C {}
	GlobalNamespace_NameProp_OBJECTLIT    // index 10: var x = {};
	GlobalNamespace_NameProp_FUNCTION     // index 11: function f() {}
	GlobalNamespace_NameProp_GET_SET      // index 12: a getter, setter, or both; e.g. `obj.b` in `const obj = {set b(x) {}};`
	GlobalNamespace_NameProp_OTHER_OBJECT // index 13: anything else, including `var x = 1;`, var x = new Something();`, etc.
)

var globalNamespace_NameProp_names = [...]string{
	"IS_DECLARED",
	"IS_MODULE_PROP",
	"IS_PROVIDED",
	"IS_USED_HAS_OWN_PROPERTY",
	"IS_EXTERN",
	"CONSTRUCTOR_TYPE",
	"INTERFACE_TYPE",
	"ENUM_TYPE",
	"NOT_A_TYPE",
	"CLASS",
	"OBJECTLIT",
	"FUNCTION",
	"GET_SET",
	"OTHER_OBJECT",
}

// Java: GlobalNamespace.NameProp.DECLARED_TYPE_KIND_MASK
const GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK int32 = (int32(1) << GlobalNamespace_NameProp_CONSTRUCTOR_TYPE) |
	(int32(1) << GlobalNamespace_NameProp_INTERFACE_TYPE) |
	(int32(1) << GlobalNamespace_NameProp_ENUM_TYPE) |
	(int32(1) << GlobalNamespace_NameProp_NOT_A_TYPE)

// Java: GlobalNamespace.NameProp.NAME_KIND_MASK
const GlobalNamespace_NameProp_NAME_KIND_MASK int32 = (int32(1) << GlobalNamespace_NameProp_CLASS) |
	(int32(1) << GlobalNamespace_NameProp_OBJECTLIT) |
	(int32(1) << GlobalNamespace_NameProp_FUNCTION) |
	(int32(1) << GlobalNamespace_NameProp_GET_SET) |
	(int32(1) << GlobalNamespace_NameProp_OTHER_OBJECT)

// Java: GlobalNamespace.NameProp.name()
func (p GlobalNamespace_NameProp) String() string { return globalNamespace_NameProp_names[p] }

// Java: GlobalNamespace.NameProp.ordinal()
func (p GlobalNamespace_NameProp) Ordinal() int { return int(p) }

// Bit is the enum field `bit` (some power of 2): `1 << p`.
func (p GlobalNamespace_NameProp) Bit() int32 {
	return int32(1) << p
}
