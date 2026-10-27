// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Norris Boyd
 *   Roger Lawrence
 *   Mike McCabe
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ==========================================================================
// Custom annotations

// GetJSType returns the compiler inferred type on this node. Not to be
// confused with getDeclaredTypeExpression() which returns the syntactically
// specified type.
//
// Java: Node.getJSType()
func (n Node) GetJSType() JSType {
	if n.tc() == 0 || (n.rec().kindFlags&kindFlagTypeIsColor) != 0 {
		return nil
	}
	return jsTypeOf(n.store().typeAt(n.tc()))
}

// GetJSTypeBase returns the JSTypeBase on this node without constructing the
// JSType interface value, or nil if no JSType is attached.
func (n Node) GetJSTypeBase() *JSTypeBase {
	if n.tc() == 0 || (n.rec().kindFlags&kindFlagTypeIsColor) != 0 {
		return nil
	}
	return n.store().typeAt(n.tc())
}

// GetJSTypeRequired returns the compiled inferred type on this node, or
// throws an NPE if there isn't one.
//
// Java: Node.getJSTypeRequired()
func (n Node) GetJSTypeRequired() JSType {
	return jutil.CheckNotNullMsg1(n.GetJSType(), "no jstypeOrColor: %s", n)
}

// Java: Node.setJSType(JSType)
func (n Node) SetJSType(x JSType) Node {
	jutil.CheckStateMsg1(n.tc() == 0 || (n.rec().kindFlags&kindFlagTypeIsColor) == 0, "%s", n)
	if b := jsTypeBase(x); b != nil {
		n.setTC(n.store().typeHandle(b))
	} else {
		n.setTC(0)
	}
	n.rec().kindFlags &^= kindFlagTypeIsColor
	return n
}

// GetColor returns the compiled inferred type on this node. Not to be
// confused with getDeclaredTypeExpression() which returns the syntactically
// specified type.
//
// Java: Node.getColor()
func (n Node) GetColor() *colors.Color {
	if n.tc() == 0 || (n.rec().kindFlags&kindFlagTypeIsColor) == 0 {
		return nil
	}
	return n.store().colorAt(n.tc())
}

// Java: Node.setColor(Color)
func (n Node) SetColor(x *colors.Color) Node {
	jutil.CheckStateMsg1(n.tc() == 0 || (n.rec().kindFlags&kindFlagTypeIsColor) != 0, "%s", n)
	if x != nil {
		n.setTC(n.store().colorHandle(x))
		n.rec().kindFlags |= kindFlagTypeIsColor
	} else {
		n.setTC(0)
		n.rec().kindFlags &^= kindFlagTypeIsColor
	}
	return n
}

// CopyTypeFrom copies a nodes JSType or Color (if present).
//
// Java: Node.copyTypeFrom(Node)
func (n Node) CopyTypeFrom(other Node) Node {
	n.setTC(n.store().translateTC(other.store(), other.tc(), other.rec().kindFlags&kindFlagTypeIsColor != 0))
	n.rec().kindFlags = (n.rec().kindFlags &^ kindFlagTypeIsColor) | (other.rec().kindFlags & kindFlagTypeIsColor)
	return n
}

// GetJSDocInfo gets the JSDocInfo attached to this node; null if no JSDoc is
// attached to this node.
//
// Java: Node.getJSDocInfo()
func (n Node) GetJSDocInfo() *JSDocInfo {
	if (n.props() & propFlagHasJSDoc) == 0 {
		return nil
	}
	a, ok := n.auxOf()
	if !ok {
		return nil
	}
	if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
		return a.jsDoc
	}
	return n.getJSDocInfoBoxed(a.ext)
}

//go:noinline
func (n Node) getJSDocInfoBoxed(box *nodeExtBox) *JSDocInfo {
	if box == nil {
		return nil
	}
	if box.jsDocInfo != nil {
		return box.jsDocInfo
	}
	return node_cast[*JSDocInfo](n.GetProp(Node_Prop_JSDOC_INFO))
}

// SetJSDocInfo sets the JSDocInfo attached to this node.
//
// Java: Node.setJSDocInfo(JSDocInfo)
func (n Node) SetJSDocInfo(info *JSDocInfo) Node {
	if info == nil {
		n.RemoveProp(Node_Prop_JSDOC_INFO)
		return n
	}
	if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
		n.orProps(propFlagHasJSDoc)
		n.setAux(nodeAux{jsDoc: info})
		return n
	}
	box := n.extBox()
	if (n.props()&propFlagHasJSDoc) != 0 && box.jsDocInfo == nil && box.propListHead != nil {
		box.propListHead = Node_RebuildListWithoutProp(box.propListHead, Node_Prop_JSDOC_INFO)
	}
	n.orProps(propFlagHasJSDoc)
	box.jsDocInfo = info
	return n
}

// SetChangeTime records that this node was last changed at `time`.
//
// Java: Node.setChangeTime(int)
func (n Node) SetChangeTime(time int32) {
	if n.nodeKind() == node_kindNode && (n.rec().kindFlags&kindFlagHasExtra) == 0 {
		n.rec().payload = uint32(time)
		return
	}
	n.PutIntProp(Node_Prop_CHANGE_TIME, time)
}

// GetChangeTime returns the time of the last change for this node.
//
// Java: Node.getChangeTime()
func (n Node) GetChangeTime() int32 {
	if n.nodeKind() == node_kindNode && (n.rec().kindFlags&kindFlagHasExtra) == 0 {
		return int32(uint32(n.rec().payload))
	}
	return n.GetIntProp(Node_Prop_CHANGE_TIME)
}

// Java: Node.setDeleted(boolean)
func (n Node) SetDeleted(deleted bool) {
	n.PutBooleanProp(Node_Prop_DELETED, deleted)
}

// Java: Node.isDeleted()
func (n Node) IsDeleted() bool {
	return n.GetBooleanProp(Node_Prop_DELETED)
}

// SetTypedefTypeProp sets, if this node represents a typedef declaration,
// the associated JSType.
//
// Java: Node.setTypedefTypeProp(JSType)
func (n Node) SetTypedefTypeProp(type_ JSType) {
	n.PutProp(Node_Prop_TYPEDEF_TYPE, type_)
}

// GetTypedefTypeProp returns, if this node represents a typedef declaration,
// the associated JSType.
//
// Java: Node.getTypedefTypeProp()
func (n Node) GetTypedefTypeProp() JSType {
	return node_cast[JSType](n.GetProp(Node_Prop_TYPEDEF_TYPE))
}

// TakeTypesAndJSDocForColorConversion clears the node's JSType and removes
// TYPE_BEFORE_CAST, DECLARED_TYPE_EXPR, and TYPEDEF_TYPE properties in a single
// pass over propListHead, returning the old JSTypeBase, whether TYPE_BEFORE_CAST
// was present, and the node's JSDocInfo if any.
func (n Node) TakeTypesAndJSDocForColorConversion() (oldType *JSTypeBase, hadBeforeCast bool, jsdoc *JSDocInfo) {
	if n.tc() != 0 {
		if (n.rec().kindFlags & kindFlagTypeIsColor) != 0 {
			jutil.CheckStateMsg1(false, "%s", n)
		}
		oldType = n.store().typeAt(n.tc())
		n.setTC(0)
	}
	if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
		if (n.props() & propFlagHasJSDoc) != 0 {
			jsdoc = n.jsDocAux()
		}
		return oldType, false, jsdoc
	}
	box := n.extBox()
	if box.typeBeforeCast != 0 {
		hadBeforeCast = true
		box.typeBeforeCast = 0
	}
	jsdoc = box.jsDocInfo
	if head := box.propListHead; head != nil {
		var removeNeedsRebuild bool
		for x := head; x != nil; x = x.next {
			switch Node_Prop(x.propType) {
			case Node_Prop_JSDOC_INFO:
				if jsdoc == nil {
					jsdoc = node_cast[*JSDocInfo](x.objectValue)
				}
			case Node_Prop_TYPE_BEFORE_CAST:
				hadBeforeCast = true
				removeNeedsRebuild = true
			case Node_Prop_DECLARED_TYPE_EXPR, Node_Prop_TYPEDEF_TYPE:
				removeNeedsRebuild = true
			}
		}
		if removeNeedsRebuild {
			box.propListHead = node_rebuildListWithoutColorConversionProps(head)
		}
	}
	n.maybeCollapseExtBox(box)
	return oldType, hadBeforeCast, jsdoc
}

func node_rebuildListWithoutColorConversionProps(item *Node_PropListItem) *Node_PropListItem {
	if item == nil {
		return nil
	}
	switch Node_Prop(item.propType) {
	case Node_Prop_TYPE_BEFORE_CAST, Node_Prop_DECLARED_TYPE_EXPR, Node_Prop_TYPEDEF_TYPE:
		return node_rebuildListWithoutColorConversionProps(item.next)
	default:
		result := node_rebuildListWithoutColorConversionProps(item.next)
		if result == item.next {
			return item
		}
		return item.Chain(result)
	}
}

// SetUnusedParameter sets the value for isUnusedParameter().
//
// Java: Node.setUnusedParameter(boolean)
func (n Node) SetUnusedParameter(unused bool) {
	n.PutBooleanProp(Node_Prop_IS_UNUSED_PARAMETER, unused)
}

// IsUnusedParameter reports whether this node is an unused function
// parameter declaration.
//
// # Set by RemoveUnusedVars
//
// Java: Node.isUnusedParameter()
func (n Node) IsUnusedParameter() bool {
	return n.GetBooleanProp(Node_Prop_IS_UNUSED_PARAMETER)
}

// SetShorthandProperty sets the isShorthandProperty annotation.
//
// Java: Node.setShorthandProperty(boolean)
func (n Node) SetShorthandProperty(shorthand bool) {
	n.PutBooleanProp(Node_Prop_IS_SHORTHAND_PROPERTY, shorthand)
}

// IsShorthandProperty reports whether this {x:x} property was originally
// parsed as {x}.
//
// Java: Node.isShorthandProperty()
func (n Node) IsShorthandProperty() bool {
	return n.GetBooleanProp(Node_Prop_IS_SHORTHAND_PROPERTY)
}

// IsOptionalEs6Typed returns whether this node is an optional node in the
// ES6 Typed syntax.
//
// Java: Node.isOptionalEs6Typed()
func (n Node) IsOptionalEs6Typed() bool {
	return n.GetBooleanProp(Node_Prop_OPT_ES6_TYPED)
}

// SetIsSyntheticBlock sets whether this is a synthetic block that should not
// be considered a real source block.
//
// Java: Node.setIsSyntheticBlock(boolean)
func (n Node) SetIsSyntheticBlock(val bool) {
	jutil.CheckState(n.rec().token == Token_BLOCK)
	n.PutBooleanProp(Node_Prop_SYNTHETIC, val)
}

// IsSyntheticBlock returns whether this is a synthetic block that should not
// be considered a real source block.
//
// Java: Node.isSyntheticBlock()
func (n Node) IsSyntheticBlock() bool {
	return n.GetBooleanProp(Node_Prop_SYNTHETIC)
}

// Java: Node.setIsSynthesizedUnfulfilledNameDeclaration(boolean)
func (n Node) SetIsSynthesizedUnfulfilledNameDeclaration(val bool) {
	jutil.CheckStateMsg1(
		n.rec().token == Token_VAR && n.HasOneChild() && n.GetFirstChild().IsName(),

		"Expected all synthetic unfulfilled declarations to be `var <name>`, found %s",
		n)

	n.PutBooleanProp(Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION, val)
}

// Java: Node.isSynthesizedUnfulfilledNameDeclaration()
func (n Node) IsSynthesizedUnfulfilledNameDeclaration() bool {
	return n.GetBooleanProp(Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION)
}

// SetUseStrict sets whether this node contained the "use strict" directive.
//
// Java: Node.setUseStrict(boolean)
func (n Node) SetUseStrict(x bool) {
	n.PutBooleanProp(Node_Prop_USE_STRICT, x)
}

// IsUseStrict returns whether this node contained the "use strict"
// directive.
//
// Java: Node.isUseStrict()
func (n Node) IsUseStrict() bool {
	return n.GetBooleanProp(Node_Prop_USE_STRICT)
}

// SetIsAddedBlock sets whether this is an added block that should not be
// considered a real source block. Eg: In "if (true) x;", the "x;" is put
// under an added block in the AST.
//
// Java: Node.setIsAddedBlock(boolean)
func (n Node) SetIsAddedBlock(val bool) {
	n.PutBooleanProp(Node_Prop_ADDED_BLOCK, val)
}

// IsAddedBlock returns whether this is an added block that should not be
// considered a real source block.
//
// Java: Node.isAddedBlock()
func (n Node) IsAddedBlock() bool {
	return n.GetBooleanProp(Node_Prop_ADDED_BLOCK)
}

// SetStaticMember sets whether this node is a static member node. This
// method is meaningful only on GETTER_DEF, SETTER_DEF or MEMBER_FUNCTION_DEF
// nodes contained within CLASS.
//
// Java: Node.setStaticMember(boolean)
func (n Node) SetStaticMember(isStatic bool) {
	n.PutBooleanProp(Node_Prop_STATIC_MEMBER, isStatic)
}

// IsStaticMember returns whether this node is a static member node. This
// method is meaningful only on GETTER_DEF, SETTER_DEF or MEMBER_FUNCTION_DEF
// nodes contained within CLASS.
//
// Java: Node.isStaticMember()
func (n Node) IsStaticMember() bool {
	return n.GetBooleanProp(Node_Prop_STATIC_MEMBER)
}

// SetIsGeneratorFunction sets whether this node is a generator node. This
// method is meaningful only on FUNCTION or MEMBER_FUNCTION_DEF nodes.
//
// Java: Node.setIsGeneratorFunction(boolean)
func (n Node) SetIsGeneratorFunction(isGenerator bool) {
	n.PutBooleanProp(Node_Prop_GENERATOR_FN, isGenerator)
}

// IsGeneratorFunction returns whether this node is a generator function
// node.
//
// Java: Node.isGeneratorFunction()
func (n Node) IsGeneratorFunction() bool {
	return n.GetBooleanProp(Node_Prop_GENERATOR_FN)
}

// SetGeneratorMarker sets whether this node subtree contains YIELD nodes.
//
// It's used in the translation of generators.
//
// Java: Node.setGeneratorMarker(boolean)
func (n Node) SetGeneratorMarker(isGeneratorMarker bool) {
	n.PutBooleanProp(Node_Prop_IS_GENERATOR_MARKER, isGeneratorMarker)
}

// IsGeneratorMarker returns whether this node was marked as containing YIELD
// nodes.
//
// It's used in the translation of generators.
//
// Java: Node.isGeneratorMarker()
func (n Node) IsGeneratorMarker() bool {
	return n.GetBooleanProp(Node_Prop_IS_GENERATOR_MARKER)
}

// SetGeneratorSafe sets the value for isGeneratorSafe().
//
// Java: Node.setGeneratorSafe(boolean)
func (n Node) SetGeneratorSafe(isGeneratorSafe bool) {
	n.PutBooleanProp(Node_Prop_IS_GENERATOR_SAFE, isGeneratorSafe)
}

// IsGeneratorSafe is used when translating ES6 generators. If this returns
// true, this Node was generated by the compiler, and it is safe to copy this
// node to the transpiled output with no further changes.
//
// Java: Node.isGeneratorSafe()
func (n Node) IsGeneratorSafe() bool {
	return n.GetBooleanProp(Node_Prop_IS_GENERATOR_SAFE)
}

// SetIsOptionalChainStart sets whether this node is the start of an optional
// chain. This method is meaningful only on OPTCHAIN_GETELEM,
// OPTCHAIN_GETPROP, OPTCHAIN_CALL.
//
// Java: Node.setIsOptionalChainStart(boolean)
func (n Node) SetIsOptionalChainStart(isOptionalChainStart bool) {
	jutil.CheckStateMsg(
		!isOptionalChainStart || n.IsOptChainGetElem() || n.IsOptChainGetProp() || n.IsOptChainCall(),
		"cannot make a non-optional node the start of an optional chain.")
	n.PutBooleanProp(Node_Prop_START_OF_OPT_CHAIN, isOptionalChainStart)
}

// IsOptionalChainStart returns whether this node is an optional chaining
// node.
//
// Java: Node.isOptionalChainStart()
func (n Node) IsOptionalChainStart() bool {
	return n.GetBooleanProp(Node_Prop_START_OF_OPT_CHAIN)
}

// SetIsArrowFunction sets whether this node is a arrow function node. This
// method is meaningful only on FUNCTION.
//
// Java: Node.setIsArrowFunction(boolean)
func (n Node) SetIsArrowFunction(isArrow bool) {
	jutil.CheckState(n.IsFunction())
	n.PutBooleanProp(Node_Prop_ARROW_FN, isArrow)
}

// IsArrowFunction returns whether this node is a arrow function node.
//
// Java: Node.isArrowFunction()
func (n Node) IsArrowFunction() bool {
	return n.rec().token == Token_FUNCTION && n.GetBooleanProp(Node_Prop_ARROW_FN)
}

// SetIsAsyncFunction sets whether this node is an async function node. This
// method is meaningful only on FUNCTION.
//
// Java: Node.setIsAsyncFunction(boolean)
func (n Node) SetIsAsyncFunction(isAsync bool) {
	jutil.CheckState(n.IsFunction())
	n.PutBooleanProp(Node_Prop_ASYNC_FN, isAsync)
}

// IsAsyncFunction returns whether this is an async function node.
//
// Java: Node.isAsyncFunction()
func (n Node) IsAsyncFunction() bool {
	return n.rec().token == Token_FUNCTION && n.GetBooleanProp(Node_Prop_ASYNC_FN)
}

// IsAsyncGeneratorFunction returns whether this is an async generator
// function node.
//
// Java: Node.isAsyncGeneratorFunction()
func (n Node) IsAsyncGeneratorFunction() bool {
	const mask = propFlag_ASYNC_FN | propFlag_GENERATOR_FN
	return n.rec().token == Token_FUNCTION && (n.props()&mask) == mask
}

// SetYieldAll sets whether this node is a generator node. This method is
// meaningful only on FUNCTION or MEMBER_FUNCTION_DEF nodes.
//
// Java: Node.setYieldAll(boolean)
func (n Node) SetYieldAll(isGenerator bool) {
	n.PutBooleanProp(Node_Prop_YIELD_ALL, isGenerator)
}

// IsYieldAll returns whether this node is a generator node. This method is
// meaningful only on FUNCTION or MEMBER_FUNCTION_DEF nodes.
//
// Java: Node.isYieldAll()
func (n Node) IsYieldAll() bool {
	return n.GetBooleanProp(Node_Prop_YIELD_ALL)
}

// SetTrailingComma indicates that there was a trailing comma in this list.
//
// Java: Node.setTrailingComma(boolean)
func (n Node) SetTrailingComma(hasTrailingComma bool) {
	n.PutBooleanProp(Node_Prop_TRAILING_COMMA, hasTrailingComma)
}

// HasTrailingComma returns true if there was a trailing comma in the orginal
// code.
//
// Java: Node.hasTrailingComma()
func (n Node) HasTrailingComma() bool {
	return n.GetBooleanProp(Node_Prop_TRAILING_COMMA)
}

// SetSideEffectFlags marks this function or constructor call's side effect
// flags. This property is only meaningful for CALL and NEW nodes.
//
// Java: Node.setSideEffectFlags(int)
func (n Node) SetSideEffectFlags(flags int32) {
	jutil.CheckStateMsg1(
		n.IsCall() || n.IsOptChainCall() || n.IsNew() || n.IsTaggedTemplateLit(),
		"Side-effect flags can only be set on invocation nodes; got %s",
		n)

	// We invert the flags before setting because we also invert the them when
	// getting; they go full circle.
	//
	// We apply the mask after inversion so that if all flags are being set
	// (all 1s), it's equivalent to storing a 0, the default value of an
	// int-prop, which has no memory cost.
	n.PutIntProp(Node_Prop_SIDE_EFFECT_FLAGS, ^flags&Node_SideEffectFlags_USED_BITS_MASK)
}

// Java: Node.setSideEffectFlags(SideEffectFlags)
func (n Node) SetSideEffectFlagsObject(flags *Node_SideEffectFlags) {
	n.SetSideEffectFlags(flags.ValueOf())
}

// GetSideEffectFlags returns the side effects flags for this node.
//
// Java: Node.getSideEffectFlags()
func (n Node) GetSideEffectFlags() int32 {
	// Int props default to 0, but we want the default for side-effect flags to
	// be all 1s. Therefore, we invert the value returned here. This is correct
	// for non-defaults because we also invert when setting the flags.
	return ^int32((n.props()>>nodePropSideEffectShift)&15) & Node_SideEffectFlags_USED_BITS_MASK
}

// Node_SideEffectFlags is a helper class for getting and setting invocation
// side-effect flags.
//
// The following values are of interest:
//
//  1. Is global state mutated? (MUTATES_GLOBAL_STATE)
//  2. Is the receiver (`this`) mutated? (MUTATES_THIS)
//  3. Are any arguments mutated? (MUTATES_ARGUMENTS)
//  4. Does the call throw an error? (THROWS)
type Node_SideEffectFlags struct {
	// A bitfield indicating the flag statuses. All used bits set to 1 means
	// "global state changes". A value of 0 means "no side effects"
	value int32
}

const (
	Node_SideEffectFlags_MUTATES_GLOBAL_STATE int32 = 1
	Node_SideEffectFlags_MUTATES_THIS         int32 = 2
	Node_SideEffectFlags_MUTATES_ARGUMENTS    int32 = 4
	Node_SideEffectFlags_THROWS               int32 = 8

	Node_SideEffectFlags_USED_BITS_MASK   int32 = (1 << 4) - 1
	Node_SideEffectFlags_NO_SIDE_EFFECTS  int32 = 0
	Node_SideEffectFlags_ALL_SIDE_EFFECTS       = Node_SideEffectFlags_MUTATES_GLOBAL_STATE | Node_SideEffectFlags_MUTATES_THIS |
		Node_SideEffectFlags_MUTATES_ARGUMENTS | Node_SideEffectFlags_THROWS
)

// Java: Node.SideEffectFlags()
func NewNode_SideEffectFlags() *Node_SideEffectFlags {
	return &Node_SideEffectFlags{value: Node_SideEffectFlags_ALL_SIDE_EFFECTS}
}

// Java: Node.SideEffectFlags(int)
func NewNode_SideEffectFlagsWithValue(value int32) *Node_SideEffectFlags {
	return &Node_SideEffectFlags{value: value}
}

// Java: Node.SideEffectFlags.valueOf()
func (f *Node_SideEffectFlags) ValueOf() int32 {
	return f.value
}

// SetAllFlags records that all side-effect occur and the returned results
// are non-local.
//
// Java: Node.SideEffectFlags.setAllFlags()
func (f *Node_SideEffectFlags) SetAllFlags() *Node_SideEffectFlags {
	f.value = Node_SideEffectFlags_ALL_SIDE_EFFECTS
	return f
}

// ClearAllFlags records that no side-effects occur.
//
// Java: Node.SideEffectFlags.clearAllFlags()
func (f *Node_SideEffectFlags) ClearAllFlags() *Node_SideEffectFlags {
	f.value = Node_SideEffectFlags_NO_SIDE_EFFECTS
	return f
}

// Java: Node.SideEffectFlags.setMutatesGlobalState()
func (f *Node_SideEffectFlags) SetMutatesGlobalState() *Node_SideEffectFlags {
	// Modify global means everything must be assumed to be modified.
	f.value |= Node_SideEffectFlags_MUTATES_GLOBAL_STATE | Node_SideEffectFlags_MUTATES_ARGUMENTS | Node_SideEffectFlags_MUTATES_THIS
	return f
}

// Java: Node.SideEffectFlags.setThrows()
func (f *Node_SideEffectFlags) SetThrows() *Node_SideEffectFlags {
	f.value |= Node_SideEffectFlags_THROWS
	return f
}

// Java: Node.SideEffectFlags.setMutatesThis()
func (f *Node_SideEffectFlags) SetMutatesThis() *Node_SideEffectFlags {
	f.value |= Node_SideEffectFlags_MUTATES_THIS
	return f
}

// Java: Node.SideEffectFlags.setMutatesArguments()
func (f *Node_SideEffectFlags) SetMutatesArguments() *Node_SideEffectFlags {
	f.value |= Node_SideEffectFlags_MUTATES_ARGUMENTS
	return f
}

// String is for debugging only.
//
// Java: Node.SideEffectFlags.toString()
func (f *Node_SideEffectFlags) String() string {
	var builder strings.Builder
	builder.WriteString("Side effects: ")

	if (f.value & Node_SideEffectFlags_MUTATES_THIS) != 0 {
		builder.WriteString("this ")
	}
	if (f.value & Node_SideEffectFlags_MUTATES_GLOBAL_STATE) != 0 {
		builder.WriteString("global ")
	}
	if (f.value & Node_SideEffectFlags_THROWS) != 0 {
		builder.WriteString("throw ")
	}
	if (f.value & Node_SideEffectFlags_MUTATES_ARGUMENTS) != 0 {
		builder.WriteString("args ")
	}

	return builder.String()
}

// IsOnlyModifiesThisCall returns whether the only side-effect is "modifies
// this" or there are no side effects.
//
// Java: Node.isOnlyModifiesThisCall()
func (n Node) IsOnlyModifiesThisCall() bool {
	// TODO(nickreid): Delete this; check if MUTATES_THIS is actually set. This
	// was left in to maintain existing behaviour but it makes the name of this
	// method misleading.
	sideEffectsBesidesMutatesThis := n.GetSideEffectFlags() &^ Node_SideEffectFlags_MUTATES_THIS
	return sideEffectsBesidesMutatesThis == Node_SideEffectFlags_NO_SIDE_EFFECTS
}

// IsOnlyModifiesArgumentsCall returns whether the only side-effect is
// "modifies arguments" or there are no side effects.
//
// Java: Node.isOnlyModifiesArgumentsCall()
func (n Node) IsOnlyModifiesArgumentsCall() bool {
	// TODO(nickreid): Delete this; check if MUTATES_ARGUMENTS is actually set.
	// This was left in to maintain existing behaviour but it makes the name of
	// this method misleading.
	sideEffectsBesidesMutatesArguments := n.GetSideEffectFlags() &^ Node_SideEffectFlags_MUTATES_ARGUMENTS
	return sideEffectsBesidesMutatesArguments == Node_SideEffectFlags_NO_SIDE_EFFECTS
}

// IsNoSideEffectsCall returns true if this node is a function or constructor
// call that has no side effects.
//
// Java: Node.isNoSideEffectsCall()
func (n Node) IsNoSideEffectsCall() bool {
	return n.GetSideEffectFlags() == Node_SideEffectFlags_NO_SIDE_EFFECTS
}

// MayMutateArguments returns true if this is a new/call that may mutate its
// arguments.
//
// Java: Node.mayMutateArguments()
func (n Node) MayMutateArguments() bool {
	return Node_AllBitsSet(n.GetSideEffectFlags(), Node_SideEffectFlags_MUTATES_ARGUMENTS)
}

// MayMutateGlobalStateOrThrow returns true if this is a new/call that may
// mutate global state or throw.
//
// Java: Node.mayMutateGlobalStateOrThrow()
func (n Node) MayMutateGlobalStateOrThrow() bool {
	return Node_AnyBitSet(
		n.GetSideEffectFlags(), Node_SideEffectFlags_MUTATES_GLOBAL_STATE|Node_SideEffectFlags_THROWS)
}

// Node_AllBitsSet returns true iff all the set bits in `mask` are also set
// in `value`.
//
// Java: Node.allBitsSet(int, int)
func Node_AllBitsSet(value int32, mask int32) bool {
	return (value & mask) == mask
}

// Node_AnyBitSet returns true iff any the bit set in `mask` is also set in
// `value`.
//
// Java: Node.anyBitSet(int, int)
func Node_AnyBitSet(value int32, mask int32) bool {
	return (value & mask) != 0
}

// Node_ConstantVarFlags holds the constants for the CONSTANT_VAR_FLAGS bit
// set property.
//
//   - DECLARED means the name was declared using annotation or syntax
//     indicating it must be constant
//   - INFERRED means the compiler can see that it is assigned exactly once,
//     whether or not it was declared to be constant.
//
// Either, both, or neither may be set for any name.
type Node_ConstantVarFlags struct{}

// each constant should be a distinct power of 2.
const (
	Node_ConstantVarFlags_DECLARED int32 = 1
	Node_ConstantVarFlags_INFERRED int32 = 2
)

// Java: Node.ConstantVarFlags()
func NewNode_ConstantVarFlags() *Node_ConstantVarFlags { return &Node_ConstantVarFlags{} }

// Java: Node.getConstantVarFlags()
func (n Node) GetConstantVarFlags() int32 {
	return int32((n.props() >> nodePropConstVarShift) & 3)
}

// Java: Node.setConstantVarFlag(int, boolean)
func (n Node) SetConstantVarFlag(flag int32, value bool) {
	flags := int32((n.props() >> nodePropConstVarShift) & 3)
	if value {
		flags = flags | flag
	} else {
		flags = flags &^ flag
	}
	n.PutIntProp(Node_Prop_CONSTANT_VAR_FLAGS, flags)
}

// IsDeclaredConstantVar returns whether this variable is declared as a
// constant.
//
// The compiler considers a variable to be declared if:
//
//   - it is declared with the `const` keyword, or
//   - It is declared with a jsdoc `@const` annotation, or
//   - The current coding convention considers it to be a constant.
//
// Only valid to call on a name node.
//
// Java: Node.isDeclaredConstantVar()
func (n Node) IsDeclaredConstantVar() bool {
	jutil.CheckStateMsg1(
		n.IsName() || n.IsImportStar(),
		"Should only be called on name or import * nodes. Found %s",
		n)

	return (n.props() & (uint64(Node_ConstantVarFlags_DECLARED) << nodePropConstVarShift)) != 0
}

// SetDeclaredConstantVar sets this variable to be a declared constant.
//
// See isDeclaredConstantVar for the rules.
//
// Java: Node.setDeclaredConstantVar(boolean)
func (n Node) SetDeclaredConstantVar(value bool) {
	jutil.CheckStateMsg1(
		n.IsName() || n.IsImportStar(),
		"Should only be called on name or import * nodes. Found %s",
		n)

	n.SetConstantVarFlag(Node_ConstantVarFlags_DECLARED, value)
}

// IsInferredConstantVar returns whether this variable is inferred to be
// constant.
//
// The compiler infers a variable to be a constant if:
//
//   - It is assigned at its declaration site, and
//   - It is never reassigned during its lifetime, and
//   - It is not defined by an extern.
//
// Only valid to call on a name node.
//
// Java: Node.isInferredConstantVar()
func (n Node) IsInferredConstantVar() bool {
	jutil.CheckStateMsg1(
		n.IsName() || n.IsImportStar(),
		"Should only be called on name or import * nodes. Found %s",
		n)

	return (n.props() & (uint64(Node_ConstantVarFlags_INFERRED) << nodePropConstVarShift)) != 0
}

// SetInferredConstantVar sets this variable to be an inferred constant.
//
// See isInferredConstantVar for the rules.
//
// Java: Node.setInferredConstantVar(boolean)
func (n Node) SetInferredConstantVar(value bool) {
	jutil.CheckStateMsg1(
		n.IsName() || n.IsImportStar(),
		"Should only be called on name or import * nodes. Found %s",
		n)

	n.SetConstantVarFlag(Node_ConstantVarFlags_INFERRED, value)
}

// Java: Node.isQuotedStringKey()
func (n Node) IsQuotedStringKey() bool {
	return (n.nodeKind() == node_kindString) && n.GetBooleanProp(Node_Prop_QUOTED)
}

// Java: Node.setQuotedStringKey()
func (n Node) SetQuotedStringKey() {
	jutil.CheckStateMsg1(n.nodeKind() == node_kindString, "%s", n)
	n.PutBooleanProp(Node_Prop_QUOTED, true)
}

// Java: Node.isPrivateIdentifier()
func (n Node) IsPrivateIdentifier() bool {
	return (n.nodeKind() == node_kindString) && n.GetBooleanProp(Node_Prop_PRIVATE_IDENTIFIER)
}

// Java: Node.setPrivateIdentifier()
func (n Node) SetPrivateIdentifier() {
	jutil.CheckStateMsg1(n.nodeKind() == node_kindString, "%s", n)
	jutil.CheckState(strings.HasPrefix(n.GetString(), "#"))
	n.PutBooleanProp(Node_Prop_PRIVATE_IDENTIFIER, true)
}
