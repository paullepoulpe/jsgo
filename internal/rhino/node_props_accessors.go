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
	"math/bits"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// GetNonJSDocCommentString gets the NonJSDoc comment string attached to this
// node.
//
// Returns the information or empty string if no nonJSDoc is attached to this
// node.
//
// Java: Node.getNonJSDocCommentString()
func (n Node) GetNonJSDocCommentString() string {
	if c := n.GetNonJSDocComment(); c != nil {
		return c.GetCommentString()
	}
	return ""
}

// Java: Node.getNonJSDocComment()
func (n Node) GetNonJSDocComment() *NonJSDocComment {
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_NON_JSDOC_COMMENT))
}

// Java: Node.getTrailingNonJSDocCommentString()
func (n Node) GetTrailingNonJSDocCommentString() string {
	if c := n.GetTrailingNonJSDocComment(); c != nil {
		return c.GetCommentString()
	}
	return ""
}

// Java: Node.getTrailingNonJSDocComment()
func (n Node) GetTrailingNonJSDocComment() *NonJSDocComment {
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_TRAILING_NON_JSDOC_COMMENT))
}

// SetNonJSDocComment sets the NonJSDoc comment attached to this node.
//
// Java: Node.setNonJSDocComment(NonJSDocComment)
func (n Node) SetNonJSDocComment(comment *NonJSDocComment) Node {
	n.PutProp(Node_Prop_NON_JSDOC_COMMENT, comment)
	return n
}

// Java: Node.setTrailingNonJSDocComment(NonJSDocComment)
func (n Node) SetTrailingNonJSDocComment(comment *NonJSDocComment) Node {
	n.PutProp(Node_Prop_TRAILING_NON_JSDOC_COMMENT, comment)
	return n
}

// SetIsParenthesized sets whether this node was inside original source-level
// parentheses.
//
// Java: Node.setIsParenthesized(boolean)
func (n Node) SetIsParenthesized(b bool) {
	jutil.CheckState(IR_MayBeExpression(n))
	n.PutBooleanProp(Node_Prop_IS_PARENTHESIZED, b)
}

// GetIsParenthesized checks whether node was inside original source-level
// parentheses.
//
// Java: Node.getIsParenthesized()
func (n Node) GetIsParenthesized() bool {
	return (n.props() & propFlag_IS_PARENTHESIZED) != 0
}

// TODO(sdh): Get rid of these by using accessor methods instead.
// These export instances of a private type, which is awkward but a side
// effect is that it prevents anyone from introducing problemmatic uses of
// the general-purpose accessors.
const (
	Node_INCRDECR_PROP          = Node_Prop_INCRDECR
	Node_QUOTED_PROP            = Node_Prop_QUOTED
	Node_IS_CONSTANT_NAME       = Node_Prop_IS_CONSTANT_NAME
	Node_IS_NAMESPACE           = Node_Prop_IS_NAMESPACE
	Node_DIRECT_EVAL            = Node_Prop_DIRECT_EVAL
	Node_FREE_CALL              = Node_Prop_FREE_CALL
	Node_REFLECTED_OBJECT       = Node_Prop_REFLECTED_OBJECT
	Node_STATIC_MEMBER          = Node_Prop_STATIC_MEMBER
	Node_GENERATOR_FN           = Node_Prop_GENERATOR_FN
	Node_YIELD_ALL              = Node_Prop_YIELD_ALL
	Node_EXPORT_DEFAULT         = Node_Prop_EXPORT_DEFAULT
	Node_EXPORT_ALL_FROM        = Node_Prop_EXPORT_ALL_FROM
	Node_COMPUTED_PROP_METHOD   = Node_Prop_COMPUTED_PROP_METHOD
	Node_COMPUTED_PROP_GETTER   = Node_Prop_COMPUTED_PROP_GETTER
	Node_COMPUTED_PROP_SETTER   = Node_Prop_COMPUTED_PROP_SETTER
	Node_COMPUTED_PROP_VARIABLE = Node_Prop_COMPUTED_PROP_VARIABLE
	Node_OPT_ES6_TYPED          = Node_Prop_OPT_ES6_TYPED
	Node_GENERIC_TYPE_LIST      = Node_Prop_GENERIC_TYPE
	Node_IMPLEMENTS             = Node_Prop_IMPLEMENTS
	Node_CONSTRUCT_SIGNATURE    = Node_Prop_CONSTRUCT_SIGNATURE
	Node_ACCESS_MODIFIER        = Node_Prop_ACCESS_MODIFIER
	Node_GOOG_MODULE            = Node_Prop_GOOG_MODULE
	Node_FEATURE_SET            = Node_Prop_FEATURE_SET
	Node_IS_TYPESCRIPT_ABSTRACT = Node_Prop_IS_TYPESCRIPT_ABSTRACT
	Node_MODULE_ALIAS           = Node_Prop_MODULE_ALIAS
	Node_MODULE_EXPORT          = Node_Prop_MODULE_EXPORT
	Node_IS_SHORTHAND_PROPERTY  = Node_Prop_IS_SHORTHAND_PROPERTY
	Node_ES6_MODULE             = Node_Prop_ES6_MODULE
)

// SetDeclaredTypeExpression sets the syntactical type specified on this
// node.
//
// Java: Node.setDeclaredTypeExpression(Node)
func (n Node) SetDeclaredTypeExpression(typeExpression Node) {
	n.PutProp(Node_Prop_DECLARED_TYPE_EXPR, typeExpression)
}

// GetDeclaredTypeExpression returns the syntactical type specified on this
// node. Not to be confused with getJSType() which returns the
// compiler-inferred type.
//
// Java: Node.getDeclaredTypeExpression()
func (n Node) GetDeclaredTypeExpression() Node {
	return node_cast[Node](n.GetProp(Node_Prop_DECLARED_TYPE_EXPR))
}

// SetJSTypeBeforeCast sets the type of this node before casting.
//
// Java: Node.setJSTypeBeforeCast(JSType)
func (n Node) SetJSTypeBeforeCast(type_ JSType) {
	if type_ == nil {
		n.RemoveProp(Node_Prop_TYPE_BEFORE_CAST)
		return
	}
	box := n.ensureExtBox()
	if box.typeBeforeCast == 0 && box.propListHead != nil {
		box.propListHead = Node_RebuildListWithoutProp(box.propListHead, Node_Prop_TYPE_BEFORE_CAST)
	}
	box.typeBeforeCast = n.store().typeHandle(jsTypeBase(type_))
}

// GetJSTypeBeforeCast returns the type of this node before casting. This
// annotation will only exist on the first child of a CAST node after type
// checking.
//
// Java: Node.getJSTypeBeforeCast()
func (n Node) GetJSTypeBeforeCast() JSType {
	if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
		return nil
	}
	return n.getJSTypeBeforeCastBoxed()
}

//go:noinline
func (n Node) getJSTypeBeforeCastBoxed() JSType {
	box := n.extBoxOrNil()
	if box == nil {
		return nil
	}
	if box.typeBeforeCast != 0 {
		return jsTypeOf(n.store().typeAt(box.typeBeforeCast))
	}
	if box.propListHead != nil {
		return node_cast[JSType](n.GetProp(Node_Prop_TYPE_BEFORE_CAST))
	}
	return nil
}

// SetColorFromTypeCast indicates that this node's color comes from a type
// assertion. Only set when colors are present; when JSTypes are on the AST
// we instead preserve the actual JSType before the type assertion.
//
// Java: Node.setColorFromTypeCast()
func (n Node) SetColorFromTypeCast() {
	jutil.CheckStateMsg(n.GetColor() != nil, "Only use on nodes with colors present")
	n.PutBooleanProp(Node_Prop_COLOR_FROM_CAST, true)
}

// IsColorFromTypeCast indicates that this node's color comes from a type
// assertion. Only set when colors are present.
//
// Java: Node.isColorFromTypeCast()
func (n Node) IsColorFromTypeCast() bool {
	return n.GetBooleanProp(Node_Prop_COLOR_FROM_CAST)
}

// Java: Node.lookupProperty(Prop)
func (n Node) LookupProperty(prop Node_Prop) *Node_PropListItem {
	if mask := nodePropBitMask[prop]; mask != 0 {
		if (n.props() & mask) == 0 {
			return nil
		}
		if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
			if x := n.lookupPropListHead(prop); x != nil {
				return x
			}
		}
		return &staticBoolPropItems[prop]
	}
	switch prop {
	case Node_Prop_SOURCE_FILE:
		if sfID := n.rec().sourceFile; sfID != 0 {
			return sourceFilePropItem(sfID)
		}
		return nil
	case Node_Prop_CONSTANT_VAR_FLAGS:
		if v := (n.props() >> nodePropConstVarShift) & 3; v != 0 {
			return &staticConstVarPropItems[v]
		}
		if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
			return nil
		}
	case Node_Prop_SIDE_EFFECT_FLAGS:
		if v := (n.props() >> nodePropSideEffectShift) & 15; v != 0 {
			return &staticSideEffectPropItems[v]
		}
		if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
			return nil
		}
	case Node_Prop_JSDOC_INFO:
		if (n.props() & propFlagHasJSDoc) == 0 {
			return nil
		}
		var jsdoc *JSDocInfo
		if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
			jsdoc = n.jsDocAux()
		} else if box := n.extBoxOrNil(); box != nil {
			jsdoc = box.jsDocInfo
		}
		if jsdoc != nil {
			return &Node_PropListItem{
				propType:    int8(Node_Prop_JSDOC_INFO),
				objectValue: jsdoc,
				isObject:    true,
				shared:      true,
			}
		}
	case Node_Prop_TYPE_BEFORE_CAST:
		box := n.extBoxOrNil()
		if box == nil {
			return nil
		}
		if tbc := n.store().typeAt(box.typeBeforeCast); tbc != nil {
			return &Node_PropListItem{
				propType:    int8(Node_Prop_TYPE_BEFORE_CAST),
				objectValue: tbc.self(),
				isObject:    true,
				shared:      true,
			}
		}
	case Node_Prop_CHANGE_TIME:
		if n.nodeKind() == node_kindNode {
			if v := int32(uint32(n.rec().payload)); v != 0 {
				return &Node_PropListItem{
					propType: int8(Node_Prop_CHANGE_TIME),
					intValue: v,
					shared:   true,
				}
			}
			if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
				return nil
			}
		}
	}
	return n.lookupPropListHead(prop)
}

// ClonePropsFrom clones the properties from the provided node without
// copying the property object. The receiving node may not have any existing
// properties.
//
// Java: Node.clonePropsFrom(Node)
func (n Node) ClonePropsFrom(other Node) Node {
	jutil.CheckStateMsg(
		n.coldPropListHead() == nil && n.typeBeforeCastBase() == nil && n.props() == 0 && n.rec().sourceFile == 0 && (n.nodeKind() != node_kindNode || uint32(n.rec().payload) == 0),
		"Node has existing properties.")
	if otherBox := other.extBoxOrNil(); otherBox != nil {
		if otherBox.propListHead != nil {
			otherBox.propListHead.markShared()
		}
		if otherBox.jsDocInfo != nil && otherBox.typeBeforeCast == 0 && otherBox.propListHead == nil && (n.rec().kindFlags&kindFlagHasExtBox) == 0 {
			n.setAux(nodeAux{jsDoc: otherBox.jsDocInfo})
		} else if otherBox.jsDocInfo != nil || otherBox.typeBeforeCast != 0 || otherBox.propListHead != nil {
			box := n.ensureExtBox()
			box.jsDocInfo = otherBox.jsDocInfo
			box.typeBeforeCast = n.store().translateTypeHandle(other.store(), otherBox.typeBeforeCast)
			box.propListHead = otherBox.propListHead
		}
	} else if (other.props() & propFlagHasJSDoc) != 0 {
		if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
			n.setAux(nodeAux{jsDoc: other.jsDocAux()})
		} else {
			n.ensureExtBox().jsDocInfo = other.jsDocAux()
		}
	}
	n.setPropsAndSourceFile(other)
	if (other.rec().kindFlags&kindFlagHasExtra) != 0 && other.coldPropListHead() != nil {
		n.rec().kindFlags |= kindFlagHasExtra
	}
	if other.nodeKind() == node_kindNode {
		if ct := int32(uint32(other.rec().payload)); ct != 0 {
			n.SetChangeTime(ct)
		}
	} else if n.nodeKind() == node_kindNode && other.coldPropListHead() != nil {
		if x := n.lookupPropListHead(Node_Prop_CHANGE_TIME); x != nil && !x.isObject {
			n.rec().payload = uint32(x.intValue)
			n.setColdPropListHead(Node_RebuildListWithoutProp(n.coldPropListHead(), Node_Prop_CHANGE_TIME))
		}
	}
	return n
}

const nodeValidatedPropsMask = propFlag_IS_PARENTHESIZED |
	propFlag_ARROW_FN |
	propFlag_ASYNC_FN |
	propFlag_SYNTHETIC |
	propFlag_COLOR_FROM_CAST |
	propFlag_START_OF_OPT_CHAIN |
	nodePropConstVarMask |
	propFlag_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION

// ValidateProperties checks for invalid or missing properties and feeds
// error messages for any violations to the given `Consumer`.
//
// We use a `Consumer` to avoid the cost of building a usually-empty list
// every time this method is called.
//
// Java: Node.validateProperties(Consumer<String>)
func (n Node) ValidateProperties(violationMessageConsumer func(string)) {
	if n.props() == 0 && n.rec().sourceFile == 0 && (n.rec().kindFlags&kindFlagHasExtBox) == 0 && (n.nodeKind() != node_kindNode || uint32(n.rec().payload) == 0) {
		return
	}

	if n.rec().token == Token_ROOT {
		// ROOT tokens should never have properties
		violationMessageConsumer("ROOT has properties")
	}

	flags := n.props()
	if (flags & nodeValidatedPropsMask) != 0 {
		if (flags&propFlag_IS_PARENTHESIZED) != 0 && !IR_MayBeExpression(n) {
			violationMessageConsumer("non-expression is parenthesized")
		}
		if (flags&propFlag_ARROW_FN) != 0 && !n.IsFunction() {
			violationMessageConsumer("invalid ARROW_FN prop")
		}
		if (flags&propFlag_ASYNC_FN) != 0 && !n.IsFunction() {
			violationMessageConsumer("invalid ASYNC_FN prop")
		}
		if (flags&propFlag_SYNTHETIC) != 0 && !n.IsBlock() {
			violationMessageConsumer("invalid SYNTHETIC prop")
		}
		if (flags&propFlag_COLOR_FROM_CAST) != 0 && n.GetColor() == nil {
			violationMessageConsumer("COLOR_FROM_CAST with no Color")
		}
		if (flags&propFlag_START_OF_OPT_CHAIN) != 0 && !(n.IsOptChainCall() || n.IsOptChainGetElem() || n.IsOptChainGetProp()) {
			violationMessageConsumer("START_OF_OPT_CHAIN on non-optional Node")
		}
		if (flags&nodePropConstVarMask) != 0 && !(n.IsName() || n.IsImportStar()) {
			violationMessageConsumer("invalid CONST_VAR_FLAGS")
		}
		if (flags&propFlag_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION) != 0 && !n.IsVar() {
			violationMessageConsumer(
				"Expected all synthetic unfulfilled declarations to be `var <name>`")
		}
	}
	if (n.rec().kindFlags & kindFlagHasExtBox) != 0 {
		for x := n.coldPropListHead(); x != nil; x = x.next {
			switch Node_Prop(x.propType) {
			case Node_Prop_CONSTANT_VAR_FLAGS:
				if !(n.IsName() || n.IsImportStar()) {
					violationMessageConsumer("invalid CONST_VAR_FLAGS")
				}
			case Node_Prop_CLOSURE_UNAWARE_SHADOW:
				if !x.isObject || !node_isNode(x.GetObjectValue()) {
					violationMessageConsumer("CLOSURE_UNAWARE_SHADOW property must point to a Node")
				}
			}
		}
	}
}

func node_isNode(v any) bool {
	_, ok := v.(Node)
	return ok
}

func (n Node) presentNodePropsMask() uint64 {
	if n.props() == 0 && n.rec().sourceFile == 0 && (n.rec().kindFlags&kindFlagHasExtBox) == 0 && (n.nodeKind() != node_kindNode || uint32(n.rec().payload) == 0) {
		return 0
	}
	var mask uint64
	boolBits := n.props() & nodePropBoolMask
	for boolBits != 0 {
		b := bits.TrailingZeros64(boolBits)
		boolBits &^= uint64(1) << uint(b)
		mask |= uint64(1) << uint(nodeBitToProp[b])
	}
	if (n.props() & propFlagHasJSDoc) != 0 {
		mask |= uint64(1) << uint(Node_Prop_JSDOC_INFO)
	}
	if (n.props() & nodePropConstVarMask) != 0 {
		mask |= uint64(1) << uint(Node_Prop_CONSTANT_VAR_FLAGS)
	}
	if (n.props() & nodePropSideEffectMask) != 0 {
		mask |= uint64(1) << uint(Node_Prop_SIDE_EFFECT_FLAGS)
	}
	if n.rec().sourceFile != 0 {
		mask |= uint64(1) << uint(Node_Prop_SOURCE_FILE)
	}
	if n.nodeKind() == node_kindNode && uint32(n.rec().payload) != 0 {
		mask |= uint64(1) << uint(Node_Prop_CHANGE_TIME)
	}
	if box := n.extBoxOrNil(); box != nil {
		if box.typeBeforeCast != 0 {
			mask |= uint64(1) << uint(Node_Prop_TYPE_BEFORE_CAST)
		}
		for x := box.propListHead; x != nil; x = x.next {
			mask |= uint64(1) << uint(x.propType)
		}
	}
	return mask
}

// Java: Node.nodePropertyToBit(NodeProperty)
func Node_NodePropertyToBit(prop NodeProperty) int64 {
	return jutil.Shl64(1, int64(prop.GetNumber()))
}

// Java: Node.setNodePropertyBit(long, NodeProperty)
func Node_SetNodePropertyBit(bitset int64, prop NodeProperty) int64 {
	return bitset | Node_NodePropertyToBit(prop)
}

// Java: Node.removeNodePropertyBit(long, NodeProperty)
func Node_RemoveNodePropertyBit(bitset int64, prop NodeProperty) int64 {
	return bitset &^ Node_NodePropertyToBit(prop)
}

// Java: Node.hasNodePropertyBitSet(long, NodeProperty)
func Node_HasNodePropertyBitSet(bitset int64, prop NodeProperty) bool {
	return (bitset & Node_NodePropertyToBit(prop)) != 0
}

// Java: Node.hasBitSet(long, int)
func Node_HasBitSet(bitset int64, bit int) bool {
	return (bitset & jutil.Shl64(1, int64(bit))) != 0
}

const nodePropNonSerializedMask = (uint64(1) << uint(Node_Prop_SOURCE_FILE)) |
	(uint64(1) << uint(Node_Prop_JSDOC_INFO)) |
	(uint64(1) << uint(Node_Prop_INPUT_ID)) |
	(uint64(1) << uint(Node_Prop_FEATURE_SET)) |
	(uint64(1) << uint(Node_Prop_DECLARED_TYPE_EXPR)) |
	(uint64(1) << uint(Node_Prop_TYPEDEF_TYPE)) |
	(uint64(1) << uint(Node_Prop_NON_JSDOC_COMMENT)) |
	(uint64(1) << uint(Node_Prop_TRAILING_NON_JSDOC_COMMENT)) |
	(uint64(1) << uint(Node_Prop_GENERIC_TYPE)) |
	(uint64(1) << uint(Node_Prop_IMPLEMENTS)) |
	(uint64(1) << uint(Node_Prop_ACCESS_MODIFIER))

// Java: Node.serializeProperties()
func (n Node) SerializeProperties() int64 {
	if n.props() == 0 && (n.rec().kindFlags&kindFlagHasExtBox) == 0 {
		return 0
	}
	bitsMask := n.presentNodePropsMask() &^ nodePropNonSerializedMask
	if bitsMask == 0 {
		return 0
	}
	var propSet int64 = 0
	for bitsMask != 0 {
		b := bits.TrailingZeros64(bitsMask)
		bitsMask &^= uint64(1) << uint(b)
		prop := Node_Prop(b)

		switch prop {
		case Node_Prop_TYPE_BEFORE_CAST:
			propSet = Node_SetNodePropertyBit(propSet, NodeProperty_COLOR_FROM_CAST)
		case Node_Prop_CONSTANT_VAR_FLAGS:
			intVal := n.GetIntProp(Node_Prop_CONSTANT_VAR_FLAGS)
			if Node_AnyBitSet(intVal, Node_ConstantVarFlags_INFERRED) {
				propSet = Node_SetNodePropertyBit(propSet, NodeProperty_IS_INFERRED_CONSTANT)
			}
			if Node_AnyBitSet(intVal, Node_ConstantVarFlags_DECLARED) {
				propSet = Node_SetNodePropertyBit(propSet, NodeProperty_IS_DECLARED_CONSTANT)
			}
		case Node_Prop_SIDE_EFFECT_FLAGS:
			intVal := n.GetIntProp(Node_Prop_SIDE_EFFECT_FLAGS)
			propSet = n.SetNodePropertySideEffectFlags(propSet, intVal)
		case Node_Prop_CLOSURE_UNAWARE_SHADOW:
			propSet = Node_SetNodePropertyBit(propSet, NodeProperty_CLOSURE_UNAWARE_SHADOW)
		default:
			if x := n.lookupPropListHead(prop); x != nil && x.isObject {
				continue
			}
			nodeProperty := PropTranslator_Serialize(prop)
			if nodeProperty != nil {
				propSet = Node_SetNodePropertyBit(propSet, *nodeProperty)
			}
		}
	}
	return propSet
}

// SetNodePropertySideEffectFlags updates a bit field to be used for
// serialized node properties to include bits from the `SIDE_EFFECT_FLAGS`
// Node property.
//
// Java: Node.setNodePropertySideEffectFlags(long, int)
func (n Node) SetNodePropertySideEffectFlags(propSet int64, sideEffectFlags int32) int64 {
	if Node_AnyBitSet(sideEffectFlags, Node_SideEffectFlags_MUTATES_GLOBAL_STATE) {
		propSet = Node_SetNodePropertyBit(propSet, NodeProperty_MUTATES_GLOBAL_STATE)
	}
	if Node_AnyBitSet(sideEffectFlags, Node_SideEffectFlags_MUTATES_THIS) {
		propSet = Node_SetNodePropertyBit(propSet, NodeProperty_MUTATES_THIS)
	}
	if Node_AnyBitSet(sideEffectFlags, Node_SideEffectFlags_MUTATES_ARGUMENTS) {
		propSet = Node_SetNodePropertyBit(propSet, NodeProperty_MUTATES_ARGUMENTS)
	}
	if Node_AnyBitSet(sideEffectFlags, Node_SideEffectFlags_THROWS) {
		propSet = Node_SetNodePropertyBit(propSet, NodeProperty_THROWS)
	}
	return propSet
}

// GetPropListDebugString returns a string representation of the properties
// on this node, for debugging.
//
// Java: Node.getPropListDebugString()
func (n Node) GetPropListDebugString() string {
	const constNameOracleMask = (uint64(1) << uint(Node_Prop_IS_CONSTANT_NAME)) |
		(uint64(1) << uint(Node_Prop_IS_NAMESPACE)) |
		(uint64(1) << uint(Node_Prop_CONSTANT_VAR_FLAGS)) |
		(uint64(1) << uint(Node_Prop_MODULE_ALIAS)) |
		(uint64(1) << uint(Node_Prop_MODULE_EXPORT)) |
		(uint64(1) << uint(Node_Prop_IS_UNUSED_PARAMETER)) |
		(uint64(1) << uint(Node_Prop_NON_INDEXABLE)) |
		(uint64(1) << uint(Node_Prop_DELETED)) |
		(uint64(1) << uint(Node_Prop_CHANGE_TIME))
	presentMask := n.presentNodePropsMask()
	if presentMask == constNameOracleMask {
		return "CHANGE_TIME->DELETED->NON_INDEXABLE->IS_UNUSED_PARAMETER->MODULE_EXPORT->MODULE_ALIAS->CONSTANT_VAR_FLAGS->IS_NAMESPACE->IS_CONSTANT_NAME"
	}
	var seen uint64
	var propListStrRepr []string
	for propListItem := n.coldPropListHead(); propListItem != nil; propListItem = propListItem.next {
		seen |= uint64(1) << uint(propListItem.propType)
		propListStrRepr = append(propListStrRepr, Node_propValues[propListItem.propType].Name())
	}
	rem := presentMask &^ seen
	hasSourceFile := (rem & (uint64(1) << uint(Node_Prop_SOURCE_FILE))) != 0
	rem &^= uint64(1) << uint(Node_Prop_SOURCE_FILE)
	for rem != 0 {
		b := 63 - bits.LeadingZeros64(rem)
		rem &^= uint64(1) << uint(b)
		propListStrRepr = append(propListStrRepr, Node_propValues[b].Name())
	}
	if hasSourceFile {
		propListStrRepr = append(propListStrRepr, Node_Prop_SOURCE_FILE.Name())
	}
	return strings.Join(propListStrRepr, "->")
}

// Java: Node.deserializeProperties(long, boolean)
func (n Node) DeserializeProperties(propSet int64, isInClosureUnawareSubtree bool) {
	if n.IsRoot() {
		jutil.CheckStateMsg1(n.coldPropListHead() == nil && n.props() == 0 && n.rec().sourceFile == 0 && uint32(n.rec().payload) == 0, "%s", n.coldPropListHead())
	} else {
		jutil.CheckStateMsg(
			n.ValidatePropListTailOrdering(isInClosureUnawareSubtree),
			"Node not ready to receive properties")
	}

	var constantVarFlags int32 = 0
	var sideEffectFlags int32 = 0
	rem := uint64(propSet) & 0x7FFFFFFFFFFFFFFF
	for rem != 0 {
		i := bits.TrailingZeros64(rem)
		rem &^= uint64(1) << uint(i)
		nodePropertyOrNull := NodeProperty_ForNumber(int32(i))
		if nodePropertyOrNull == nil {
			// Java: switch on a null enum throws.
			panic(jutil.NewNullPointerException())
		}
		nodeProperty := *nodePropertyOrNull
		switch nodeProperty {
		case NodeProperty_IS_DECLARED_CONSTANT:
			constantVarFlags |= Node_ConstantVarFlags_DECLARED
		case NodeProperty_IS_INFERRED_CONSTANT:
			constantVarFlags |= Node_ConstantVarFlags_INFERRED
		case NodeProperty_MUTATES_GLOBAL_STATE:
			sideEffectFlags |= Node_SideEffectFlags_MUTATES_GLOBAL_STATE
		case NodeProperty_MUTATES_THIS:
			sideEffectFlags |= Node_SideEffectFlags_MUTATES_THIS
		case NodeProperty_MUTATES_ARGUMENTS:
			sideEffectFlags |= Node_SideEffectFlags_MUTATES_ARGUMENTS
		case NodeProperty_THROWS:
			sideEffectFlags |= Node_SideEffectFlags_THROWS
		case NodeProperty_CLOSURE_UNAWARE_SHADOW:
			// Ignore this - we need the deserialized child to actually set the
			// shadow. We'll check for this higher up in the call stack.
		default:
			prop := PropTranslator_Deserialize(nodeProperty)
			if prop == nil {
				panic(jutil.NewIllegalStateExceptionMsg("Can not translate " + nodeProperty.String() + " to AST Prop"))
			}
			if mask := nodePropBitMask[*prop]; mask != 0 {
				n.orProps(mask)
			} else {
				n.PutBooleanProp(*prop, true)
			}
		}
	}

	if constantVarFlags != 0 {
		n.setProps((n.props() &^ nodePropConstVarMask) |
			(uint64(constantVarFlags&3) << nodePropConstVarShift))
	}

	if sideEffectFlags != 0 {
		n.setProps((n.props() &^ nodePropSideEffectMask) |
			(uint64(sideEffectFlags&15) << nodePropSideEffectShift))
	}

	// Make sure the deserialized properties are valid.
	n.ValidateProperties(
		func(errorMessage string) {
			jutil.CheckStateMsg2(false, "deserialize error: %s: %s", errorMessage, n)
		})
}

// GetSortedPropTypes gets all the property types, in sorted order.
//
// Java: Node.getSortedPropTypes()
func (n Node) GetSortedPropTypes() []int8 {
	bitsMask := n.presentNodePropsMask()
	count := bits.OnesCount64(bitsMask)
	if count == 0 {
		return nil
	}
	keys := make([]int8, 0, count)
	for bitsMask != 0 {
		b := bits.TrailingZeros64(bitsMask)
		bitsMask &^= uint64(1) << uint(b)
		keys = append(keys, int8(b))
	}
	return keys
}

// NodePropListHandleForTesting is returned by GetPropListHeadForTesting to check property sharing.
type NodePropListHandleForTesting struct {
	head  *Node_PropListItem
	flags uint64
}

// Java: Node.getPropListHeadForTesting()
func (n Node) GetPropListHeadForTesting() NodePropListHandleForTesting {
	return NodePropListHandleForTesting{head: n.coldPropListHead(), flags: n.props()}
}

// Java: Node.setPropListHead(PropListItem)
func (n Node) SetPropListHead(propListHead *Node_PropListItem) {
	if (n.rec().kindFlags & kindFlagHasExtBox) != 0 {
		box := n.extBox()
		box.jsDocInfo = nil
		box.typeBeforeCast = 0
	} else if (n.props() & propFlagHasJSDoc) != 0 {
		n.setAux(nodeAux{})
	}
	n.setProps(0)
	n.rec().sourceFile = 0
	// Go-only: the source file is an id in the record, not a list item.
	var sourceFile any
	for x := propListHead; x != nil; x = x.next {
		if Node_Prop(x.propType) == Node_Prop_SOURCE_FILE && x.isObject {
			sourceFile = x.objectValue
		}
	}
	if sourceFile != nil {
		propListHead = Node_RebuildListWithoutProp(propListHead, Node_Prop_SOURCE_FILE)
	}
	propListHead.markShared()
	n.setColdPropListHead(propListHead)
	var flags uint64
	for x := propListHead; x != nil; x = x.next {
		prop := Node_Prop(x.propType)
		if x.isObject {
			if !jutil.IsNil(x.objectValue) {
				if prop == Node_Prop_JSDOC_INFO {
					flags |= propFlagHasJSDoc
				} else if mask := nodePropBitMask[prop]; mask != 0 {
					flags |= mask
					n.rec().kindFlags |= kindFlagHasExtra
				}
			}
		} else if x.intValue != 0 {
			if mask := nodePropBitMask[prop]; mask != 0 {
				flags |= mask
				if x.intValue != 1 {
					n.rec().kindFlags |= kindFlagHasExtra
				}
			} else if prop == Node_Prop_CONSTANT_VAR_FLAGS {
				if x.intValue&^3 == 0 {
					flags = (flags &^ nodePropConstVarMask) | (uint64(x.intValue&3) << nodePropConstVarShift)
				} else {
					n.rec().kindFlags |= kindFlagHasExtra
				}
			} else if prop == Node_Prop_SIDE_EFFECT_FLAGS {
				if x.intValue&^15 == 0 {
					flags = (flags &^ nodePropSideEffectMask) | (uint64(x.intValue&15) << nodePropSideEffectShift)
				} else {
					n.rec().kindFlags |= kindFlagHasExtra
				}
			}
		}
	}
	n.setProps(flags)
}

// GetSecondToLastPropListItem returns the second to last item in the
// property list, for nodes in a closure-unaware subtree this will return
// IS_IN_CLOSURE_UNAWARE_SUBTREE -> SOURCE_FILE. For nodes int the main AST
// this will return [ANY_OTHER_PROPERTY] -> SOURCE_FILE. This method is only
// intended to temporarily test the IS_IN_CLOSURE_UNAWARE_SUBTREE ->
// SOURCE_FILE part of the property list is correctly shared between two
// subtrees of the same @closureUnaware annotated source file.
//
// TODO(user): these getters and the following two validation methods are a
// bit redundant. We should consider simplifying them in a follow up CL.
//
// Java: Node.getSecondToLastPropListItem()
func (n Node) GetSecondToLastPropListItem() *Node_PropListItem {
	if n.GetIsInClosureUnawareSubtree() {
		return &staticBoolPropItems[Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE]
	}
	propListItem := n.coldPropListHead()
	if n.rec().sourceFile != 0 {
		if propListItem == nil {
			return nil
		}
		for propListItem.next != nil {
			propListItem = propListItem.next
		}
		return propListItem
	}
	if propListItem == nil || propListItem.next == nil {
		return nil
	}
	for propListItem.next.next != nil {
		propListItem = propListItem.next
	}
	return propListItem
}

// Java: Node.getLastPropListItem()
func (n Node) GetLastPropListItem() *Node_PropListItem {
	if sfID := n.rec().sourceFile; sfID != 0 {
		return sourceFilePropItem(sfID)
	}
	propListItem := n.coldPropListHead()
	if propListItem == nil {
		return nil
	}
	for propListItem.next != nil {
		propListItem = propListItem.next
	}
	return propListItem
}

// Java: Node.validatePropListTailOrdering(boolean)
func (n Node) ValidatePropListTailOrdering(shouldNodeBeClosureUnaware bool) bool {
	curr := n.GetLastPropListItem()
	if curr == nil || int(curr.propType) != Node_Prop_SOURCE_FILE.Ordinal() {
		return false
	}

	if shouldNodeBeClosureUnaware {
		if !n.GetIsInClosureUnawareSubtree() {
			return false
		}
	}
	return true
}

// Java: Node.validateMemorySensitivePropertyGuarantees(Node, boolean, Node, boolean)
func Node_ValidateMemorySensitivePropertyGuarantees(
	node1 Node,
	shouldNode1BeClosureUnaware bool,
	node2 Node,
	shouldNode2BeClosureUnaware bool,
) bool {
	// checks to do here:
	// 1. there should be only one sourcefile proplist item and it should be
	// the last item
	// 2. if there is a proplist item that is IS_INCLOSURE_UNAWARE_SUBTREE it
	// should be followed by a SOURCE_FILE proplist item
	// 3. if there is another node, check that the actual proplist items are
	// the same for the sourcefile and the IS_IN_CLOSURE_UNAWARE_SUBTREE
	// validate that the last two elements of the property list are
	// IS_IN_CLOSURE_UNAWARE_SUBTREE and SOURCE_FILE, also validate that these
	// propelist items are the same actual item and not only the same values.

	if !node1.ValidatePropListTailOrdering(shouldNode1BeClosureUnaware) {
		return false
	}
	if !node2.ValidatePropListTailOrdering(shouldNode2BeClosureUnaware) {
		return false
	}
	if node1.GetLastPropListItem() != node2.GetLastPropListItem() {
		return false
	}
	if shouldNode1BeClosureUnaware && shouldNode2BeClosureUnaware {
		if node1.GetSecondToLastPropListItem() != node2.GetSecondToLastPropListItem() {
			return false
		}
	}
	// return true if all of the property list items are what we expect.
	return true
}
