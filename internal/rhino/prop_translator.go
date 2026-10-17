// Port of com/google/javascript/rhino/PropTranslator.java @ closure-compiler v20261015

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
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// PropTranslator is a translator for converting between Rhino node properties
// and TypedAST proto node properties used by the Rhino Node class.
type PropTranslator struct{}

// Java: PropTranslator.protoToRhinoProp
var PropTranslator_protoToRhinoProp = make([]*Node_Prop, len(NodeProperty_values()))

// Java: PropTranslator.rhinoToProtoProp
var PropTranslator_rhinoToProtoProp = make([]*NodeProperty, len(Node_Prop_values()))

// Java: PropTranslator.serialize(Prop)
func PropTranslator_Serialize(x Node_Prop) *NodeProperty {
	return PropTranslator_rhinoToProtoProp[x.Ordinal()]
}

// Java: PropTranslator.deserialize(NodeProperty)
func PropTranslator_Deserialize(x NodeProperty) *Node_Prop {
	return PropTranslator_protoToRhinoProp[x.Ordinal()]
}

// Java: PropTranslator()
func NewPropTranslator() *PropTranslator {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: PropTranslator static initializer
func init() {
	PropTranslator_SetProps()
	PropTranslator_CheckUnexpectedNullProtoProps()
}

// Java: PropTranslator.setProps()
func PropTranslator_SetProps() {
	for _, rhinoProp := range Node_Prop_values() {
		protoProp := PropTranslator_SerializeProp(rhinoProp)
		if protoProp != nil {
			// Boolean props are stored as a bitset, see Node#deserializeProperties
			jutil.CheckStateMsg2(
				protoProp.GetNumber() < 63, "enum %s value %s", *protoProp, protoProp.GetNumber())

			p := rhinoProp
			PropTranslator_protoToRhinoProp[protoProp.Ordinal()] = &p
			PropTranslator_rhinoToProtoProp[rhinoProp.Ordinal()] = protoProp
		}
	}
}

// Java: PropTranslator.serializeProp(Prop)
func PropTranslator_SerializeProp(prop Node_Prop) *NodeProperty {
	np := func(p NodeProperty) *NodeProperty { return &p }
	switch prop {
	case Node_Prop_ARROW_FN:
		return np(NodeProperty_ARROW_FN)
	case Node_Prop_ASYNC_FN:
		return np(NodeProperty_ASYNC_FN)
	case Node_Prop_GENERATOR_FN:
		return np(NodeProperty_GENERATOR_FN)
	case Node_Prop_YIELD_ALL:
		return np(NodeProperty_YIELD_ALL)
	case Node_Prop_IS_PARENTHESIZED:
		return np(NodeProperty_IS_PARENTHESIZED)
	case Node_Prop_SYNTHETIC:
		return np(NodeProperty_SYNTHETIC)
	case Node_Prop_ADDED_BLOCK:
		return np(NodeProperty_ADDED_BLOCK)
	case Node_Prop_STATIC_MEMBER:
		return np(NodeProperty_STATIC_MEMBER)
	case Node_Prop_IS_GENERATOR_MARKER:
		return np(NodeProperty_IS_GENERATOR_MARKER)
	case Node_Prop_IS_GENERATOR_SAFE:
		return np(NodeProperty_IS_GENERATOR_SAFE)
	case Node_Prop_COLOR_FROM_CAST:
		return np(NodeProperty_COLOR_FROM_CAST)
	case Node_Prop_NON_INDEXABLE:
		return np(NodeProperty_NON_INDEXABLE)
	case Node_Prop_DELETED:
		return np(NodeProperty_DELETED)
	case Node_Prop_IS_UNUSED_PARAMETER:
		return np(NodeProperty_IS_UNUSED_PARAMETER)
	case Node_Prop_IS_SHORTHAND_PROPERTY:
		return np(NodeProperty_IS_SHORTHAND_PROPERTY)
	case Node_Prop_START_OF_OPT_CHAIN:
		return np(NodeProperty_START_OF_OPT_CHAIN)
	case Node_Prop_TRAILING_COMMA:
		return np(NodeProperty_TRAILING_COMMA)
	case Node_Prop_IS_CONSTANT_NAME:
		return np(NodeProperty_IS_CONSTANT_NAME)
	case Node_Prop_IS_NAMESPACE:
		return np(NodeProperty_IS_NAMESPACE)
	case Node_Prop_DIRECT_EVAL:
		return np(NodeProperty_DIRECT_EVAL)
	case Node_Prop_FREE_CALL:
		return np(NodeProperty_FREE_CALL)
	case Node_Prop_REFLECTED_OBJECT:
		return np(NodeProperty_REFLECTED_OBJECT)
	case Node_Prop_EXPORT_DEFAULT:
		return np(NodeProperty_EXPORT_DEFAULT)
	case Node_Prop_EXPORT_ALL_FROM:
		return np(NodeProperty_EXPORT_ALL_FROM)
	case Node_Prop_COMPUTED_PROP_METHOD:
		return np(NodeProperty_COMPUTED_PROP_METHOD)
	case Node_Prop_COMPUTED_PROP_GETTER:
		return np(NodeProperty_COMPUTED_PROP_GETTER)
	case Node_Prop_COMPUTED_PROP_SETTER:
		return np(NodeProperty_COMPUTED_PROP_SETTER)
	case Node_Prop_COMPUTED_PROP_VARIABLE:
		return np(NodeProperty_COMPUTED_PROP_VARIABLE)
	case Node_Prop_GOOG_MODULE:
		return np(NodeProperty_GOOG_MODULE)
	case Node_Prop_MODULE_ALIAS:
		return np(NodeProperty_MODULE_ALIAS)
	case Node_Prop_MODULE_EXPORT:
		return np(NodeProperty_MODULE_EXPORT)
	case Node_Prop_ES6_MODULE:
		return np(NodeProperty_ES6_MODULE)
	case Node_Prop_CONSTANT_VAR_FLAGS:
		return np(NodeProperty_CONSTANT_VAR_FLAGS)
	case Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION:
		return np(NodeProperty_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION)
	case Node_Prop_PRIVATE_IDENTIFIER:
		return np(NodeProperty_PRIVATE_IDENTIFIER)
	case Node_Prop_CLOSURE_UNAWARE_SHADOW,
		Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE,
		Node_Prop_SIDE_EFFECT_FLAGS,
		Node_Prop_DECLARED_TYPE_EXPR,
		Node_Prop_FEATURE_SET,
		Node_Prop_TYPE_BEFORE_CAST,
		Node_Prop_NON_JSDOC_COMMENT,
		Node_Prop_TRAILING_NON_JSDOC_COMMENT,
		Node_Prop_JSDOC_INFO,
		Node_Prop_INCRDECR,
		Node_Prop_QUOTED,
		Node_Prop_USE_STRICT,
		Node_Prop_SOURCE_FILE,
		Node_Prop_INPUT_ID,
		Node_Prop_CHANGE_TIME,
		Node_Prop_OPT_ES6_TYPED,
		Node_Prop_GENERIC_TYPE,
		Node_Prop_IMPLEMENTS,
		Node_Prop_CONSTRUCT_SIGNATURE,
		Node_Prop_ACCESS_MODIFIER,
		Node_Prop_IS_TYPESCRIPT_ABSTRACT,
		Node_Prop_TYPEDEF_TYPE:
		// (IS_IN_CLOSURE_UNAWARE_SUBTREE is intentionally not included in the
		// serialized TypedAST because it can be re-synthesized when we
		// deserialize nodes later, as it is inferrable from the AST shape.)
		// These cases cannot be translated to a NodeProperty
		return nil
	}
	panic(jutil.NewIllegalStateExceptionMsg("unexpected Prop " + prop.String()))
}

// Java: PropTranslator.checkUnexpectedNullProtoProps()
func PropTranslator_CheckUnexpectedNullProtoProps() {
	for _, protoProp := range NodeProperty_values() {
		switch protoProp {
		case NodeProperty_NODE_PROPERTY_UNSPECIFIED:
			// this is the "no properties are set" value
		case NodeProperty_IS_DECLARED_CONSTANT, NodeProperty_IS_INFERRED_CONSTANT:
			// These are used for the CONSTANT_VAR_FLAGS bit field property
		case NodeProperty_UNRECOGNIZED, NodeProperty_UNUSED_11:
			// unused
		case NodeProperty_MUTATES_GLOBAL_STATE, NodeProperty_MUTATES_THIS, NodeProperty_MUTATES_ARGUMENTS, NodeProperty_THROWS:
			// these are used for the SIDE_EFFECT_FLAGS bit field property
		case NodeProperty_CLOSURE_UNAWARE_SHADOW:
			// this is a special case where the property is a Node pointer, not a boolean
		default:
			// everything else should be a 1-to-1 match with a boolean property
			jutil.CheckStateMsg1(
				PropTranslator_protoToRhinoProp[protoProp.Ordinal()] != nil,
				"Hit unhandled node property: %s",
				protoProp)

		}
	}
}
