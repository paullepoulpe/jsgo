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
	"fmt"
	"io"
	"iter"
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Node implements the root of the intermediate representation.
//
// Go-only: language-adapter — Java's private subclasses NumberNode,
// BigIntNode, StringNode and TemplateLiteralSubstringNode only add a payload
// and override isEquivalentTo/cloneNode, so they are folded into Node
// (PORTING.md §4): `kind` records which class a node is, the payload fields
// are only meaningful for that kind, and the accessors throw the
// ClassCastException Java's casts throw on any other kind.

// Node_Prop is a node property kind.
//
// Java: Node.Prop
type Node_Prop int32

const (
	// Is this Node within parentheses
	Node_Prop_IS_PARENTHESIZED Node_Prop = iota
	// Contains non-JSDoc comment
	Node_Prop_NON_JSDOC_COMMENT
	Node_Prop_TRAILING_NON_JSDOC_COMMENT
	// Contains a JSDocInfo object
	Node_Prop_JSDOC_INFO
	// Whether incrdecr is pre (false) or post (true)
	Node_Prop_INCRDECR
	// Set to indicate a quoted object lit key
	Node_Prop_QUOTED
	// A synthetic block. Used to make processing simpler, and does not represent a real block in
	// the source.
	Node_Prop_SYNTHETIC
	// Used to indicate BLOCK that is added
	Node_Prop_ADDED_BLOCK
	// Function or constructor call side effect flags.
	Node_Prop_SIDE_EFFECT_FLAGS
	// The variable or property is constant.
	// TODO(lukes): either document the differences or otherwise reconcile with CONSTANT_VAR_FLAGS
	Node_Prop_IS_CONSTANT_NAME
	// The variable creates a namespace.
	Node_Prop_IS_NAMESPACE
	// The presence of the "use strict" directive on this node.
	Node_Prop_USE_STRICT
	// ES5 distinguishes between direct and indirect calls to eval.
	Node_Prop_DIRECT_EVAL
	// A CALL without an explicit "this" value.
	Node_Prop_FREE_CALL
	// A StaticSourceFile indicating the file where this node lives.
	Node_Prop_SOURCE_FILE
	// The id of the input associated with this node.
	Node_Prop_INPUT_ID
	// For passes that work only on changed funs.
	Node_Prop_CHANGE_TIME
	// An object that's used for goog.object.reflect-style reflection.
	Node_Prop_REFLECTED_OBJECT
	// Set if class member definition is static
	Node_Prop_STATIC_MEMBER
	// Set if the node is a Generator function or member method.
	Node_Prop_GENERATOR_FN
	// Set if the node is an arrow function.
	Node_Prop_ARROW_FN
	// http://tc39.github.io/ecmascript-asyncawait/
	Node_Prop_ASYNC_FN
	// Set if a yield is a "yield all"
	Node_Prop_YIELD_ALL
	// Set if a export is a "default" export
	Node_Prop_EXPORT_DEFAULT
	// Set if an export is a "*"
	Node_Prop_EXPORT_ALL_FROM
	// A variable is inferred or declared as const meaning it is only ever assigned once at its
	// declaration site.
	// This is an int prop that holds a bitset of {@link ConstantVarFlags} values.
	Node_Prop_CONSTANT_VAR_FLAGS
	// Used by the ES6-to-ES3 translator.
	Node_Prop_IS_GENERATOR_MARKER
	// Used by the ES6-to-ES3 translator.
	Node_Prop_IS_GENERATOR_SAFE
	// A computed property that has the method syntax
	//   ( [prop]() {...} )
	// rather than the property definition syntax
	//   ( [prop]: value ).
	Node_Prop_COMPUTED_PROP_METHOD
	// A computed property in a getter, e.g. var obj = { get [prop]() {...} };
	Node_Prop_COMPUTED_PROP_GETTER
	// A computed property in a setter, e.g. var obj = 32;
	Node_Prop_COMPUTED_PROP_SETTER
	// A computed property that's a variable, e.g. [prop]: string;
	Node_Prop_COMPUTED_PROP_VARIABLE
	// Used to attach TypeDeclarationNode ASTs to Nodes which represent a typed NAME or FUNCTION.
	Node_Prop_DECLARED_TYPE_EXPR
	// The type of an expression before the cast. This will be present only if the expression is
	// casted.
	Node_Prop_TYPE_BEFORE_CAST
	// Indicates that this epxression was casted but we don't necessarily know to which type
	Node_Prop_COLOR_FROM_CAST
	// The node is an optional parameter or property in ES6 Typed syntax.
	Node_Prop_OPT_ES6_TYPED
	// Generic type list in ES6 typed syntax.
	Node_Prop_GENERIC_TYPE
	// "implements" clause in ES6 typed syntax.
	Node_Prop_IMPLEMENTS
	// This node is a TypeScript ConstructSignature
	Node_Prop_CONSTRUCT_SIGNATURE
	// TypeScript accessibility modifiers (public, protected, private)
	Node_Prop_ACCESS_MODIFIER
	// Indicates the node should not be indexed by analysis tools.
	Node_Prop_NON_INDEXABLE
	// Indicates that a SCRIPT node is a goog.module. Remains set after the goog.module is
	// desugared.
	Node_Prop_GOOG_MODULE
	// Attaches a FeatureSet to SCRIPT nodes.
	Node_Prop_FEATURE_SET
	// Indicates a TypeScript abstract method or class, for use in Migrants
	Node_Prop_IS_TYPESCRIPT_ABSTRACT
	// For passes that work only on deleted funs.
	Node_Prop_DELETED
	// Indicates that the node is an alias or a name from goog.require'd module or ES6
	// module. Aliases are desugared and inlined by compiler passes but we need to preserve them for
	// building index.
	Node_Prop_MODULE_ALIAS
	// Mark a parameter as unused. Used to defer work from RemovedUnusedVars to OptimizeParameters.
	Node_Prop_IS_UNUSED_PARAMETER
	// Mark a property as a module export so that collase properties can act on it.
	Node_Prop_MODULE_EXPORT
	// Indicates that a property {x:x} was originally parsed as {x}.
	Node_Prop_IS_SHORTHAND_PROPERTY
	// Indicates that a SCRIPT node is or was an ES module. Remains set after the module is
	// rewritten.
	Node_Prop_ES6_MODULE
	// Record the type associated with a @typedef to enable looking up typedef in the AST possible
	// without saving the type scope.
	Node_Prop_TYPEDEF_TYPE
	// Indicate that a OPTCHAIN_GETPROP, OPTCHAIN_GETELEM, or OPTCHAIN_CALL is the start of an
	// optional chain.
	Node_Prop_START_OF_OPT_CHAIN
	// Indicates a trailing comma in an array literal, object literal, parameter list, or argument
	// list
	Node_Prop_TRAILING_COMMA
	// Indicates that a variable declaration was synthesized to provide a declaration of some
	// name referenced in code but never defined, as most compiler passes expect that to be an
	// invariant.
	// Only present in the "synthetic externs file". Builds initialized using a
	// "TypedAST filesystem" will delete any such declarations present in a different compilation
	// shard
	Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
	// This prop holds a reference to a closure-unaware sub-AST.
	Node_Prop_CLOSURE_UNAWARE_SHADOW
	// Indicates that a string node is a private identifier (e.g. `class { #privateProp }`).
	Node_Prop_PRIVATE_IDENTIFIER
	// Indicates that this node is part of a subtree detached because it originated
	// from a Closure-unaware file.
	Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE
)

var node_Prop_names = [...]string{
	"IS_PARENTHESIZED",
	"NON_JSDOC_COMMENT",
	"TRAILING_NON_JSDOC_COMMENT",
	"JSDOC_INFO",
	"INCRDECR",
	"QUOTED",
	"SYNTHETIC",
	"ADDED_BLOCK",
	"SIDE_EFFECT_FLAGS",
	"IS_CONSTANT_NAME",
	"IS_NAMESPACE",
	"USE_STRICT",
	"DIRECT_EVAL",
	"FREE_CALL",
	"SOURCE_FILE",
	"INPUT_ID",
	"CHANGE_TIME",
	"REFLECTED_OBJECT",
	"STATIC_MEMBER",
	"GENERATOR_FN",
	"ARROW_FN",
	"ASYNC_FN",
	"YIELD_ALL",
	"EXPORT_DEFAULT",
	"EXPORT_ALL_FROM",
	"CONSTANT_VAR_FLAGS",
	"IS_GENERATOR_MARKER",
	"IS_GENERATOR_SAFE",
	"COMPUTED_PROP_METHOD",
	"COMPUTED_PROP_GETTER",
	"COMPUTED_PROP_SETTER",
	"COMPUTED_PROP_VARIABLE",
	"DECLARED_TYPE_EXPR",
	"TYPE_BEFORE_CAST",
	"COLOR_FROM_CAST",
	"OPT_ES6_TYPED",
	"GENERIC_TYPE",
	"IMPLEMENTS",
	"CONSTRUCT_SIGNATURE",
	"ACCESS_MODIFIER",
	"NON_INDEXABLE",
	"GOOG_MODULE",
	"FEATURE_SET",
	"IS_TYPESCRIPT_ABSTRACT",
	"DELETED",
	"MODULE_ALIAS",
	"IS_UNUSED_PARAMETER",
	"MODULE_EXPORT",
	"IS_SHORTHAND_PROPERTY",
	"ES6_MODULE",
	"TYPEDEF_TYPE",
	"START_OF_OPT_CHAIN",
	"TRAILING_COMMA",
	"SYNTHESIZED_UNFULFILLED_NAME_DECLARATION",
	"CLOSURE_UNAWARE_SHADOW",
	"PRIVATE_IDENTIFIER",
	"IS_IN_CLOSURE_UNAWARE_SUBTREE",
}

// Node_Prop_values is Prop.values().
func Node_Prop_values() []Node_Prop {
	out := make([]Node_Prop, len(node_Prop_names))
	for i := range out {
		out[i] = Node_Prop(i)
	}
	return out
}

// String is Enum.name() / toString().
func (p Node_Prop) String() string { return node_Prop_names[p] }

// Name is Enum.name().
func (p Node_Prop) Name() string { return node_Prop_names[p] }

// Ordinal is Enum.ordinal().
func (p Node_Prop) Ordinal() int { return int(p) }

// Avoid cloning "values" repeatedly in hot code, we save it off now.
//
// Java: Node.propValues
var Node_propValues = Node_Prop_values()

// GetNonJSDocCommentString gets the NonJSDoc comment string attached to this
// node.
//
// Returns the information or empty string if no nonJSDoc is attached to this
// node.
//
// Java: Node.getNonJSDocCommentString()
func (n Node) GetNonJSDocCommentString() string {
	if n.GetProp(Node_Prop_NON_JSDOC_COMMENT) == nil {
		return ""
	}
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_NON_JSDOC_COMMENT)).GetCommentString()
}

// Java: Node.getNonJSDocComment()
func (n Node) GetNonJSDocComment() *NonJSDocComment {
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_NON_JSDOC_COMMENT))
}

// Java: Node.getTrailingNonJSDocCommentString()
func (n Node) GetTrailingNonJSDocCommentString() string {
	if n.GetProp(Node_Prop_TRAILING_NON_JSDOC_COMMENT) == nil {
		return ""
	}
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_TRAILING_NON_JSDOC_COMMENT)).GetCommentString()
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

// node_kind records which Java class a Node is (see the Node doc comment).
type node_kind uint8

const (
	node_kindNode node_kind = iota
	node_kindString
	node_kindNumber
	node_kindBigInt
	node_kindTemplateLiteralSubstring
)

const (
	kindMask            uint8 = 7
	kindFlagTypeIsColor uint8 = 1 << 3
	kindFlagNoBlockDecl uint8 = 1 << 4
	// kindFlagHasExtra: the cold property list may hold the value of a
	// property that is normally stored inline (e.g. an int value other than
	// 1 for a boolean property, or CHANGE_TIME of a non-Node kind).
	kindFlagHasExtra uint8 = 1 << 5
	// kindFlagOrigNameIsString: the original name of this string node is its
	// own string (node_payload.go).
	kindFlagOrigNameIsString uint8 = 1 << 6
	// kindFlagHasExtBox: the aux entry holds a *nodeExtBox (otherwise the
	// *JSDocInfo when propFlagHasJSDoc is set).
	kindFlagHasExtBox uint8 = 1 << 7
)

const (
	// The first propHotCount booleans are stored in the record's property
	// word; the others (rare: under 1 % of the nodes each on the benchmark
	// workloads) in the node's aux entry (see propWord).
	propBit_IS_CONSTANT_NAME uint = iota
	propBit_QUOTED
	propBit_FREE_CALL
	propBit_NON_INDEXABLE
	propBit_ADDED_BLOCK
	propBit_TRAILING_COMMA
	propBit_IS_PARENTHESIZED
	propBit_ARROW_FN
	propBit_INCRDECR
	propBit_SYNTHETIC
	propBit_IS_NAMESPACE
	propBit_USE_STRICT
	propBit_DIRECT_EVAL
	propBit_STATIC_MEMBER
	propBit_GENERATOR_FN
	propBit_ASYNC_FN
	propBit_YIELD_ALL
	propBit_EXPORT_DEFAULT
	propBit_EXPORT_ALL_FROM
	propBit_IS_GENERATOR_MARKER
	propBit_IS_GENERATOR_SAFE
	propBit_COMPUTED_PROP_METHOD
	propBit_COMPUTED_PROP_GETTER
	propBit_COMPUTED_PROP_SETTER
	propBit_COMPUTED_PROP_VARIABLE
	propBit_COLOR_FROM_CAST
	propBit_OPT_ES6_TYPED
	propBit_CONSTRUCT_SIGNATURE
	propBit_GOOG_MODULE
	propBit_IS_TYPESCRIPT_ABSTRACT
	propBit_DELETED
	propBit_MODULE_ALIAS
	propBit_IS_UNUSED_PARAMETER
	propBit_MODULE_EXPORT
	propBit_IS_SHORTHAND_PROPERTY
	propBit_ES6_MODULE
	propBit_START_OF_OPT_CHAIN
	propBit_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
	propBit_PRIVATE_IDENTIFIER
	propBit_IS_IN_CLOSURE_UNAWARE_SUBTREE
	propBitCount
)

// propHotCount is the number of boolean properties kept in the record.
const propHotCount = 8

const (
	propFlag_IS_PARENTHESIZED                         uint64 = 1 << propBit_IS_PARENTHESIZED
	propFlag_INCRDECR                                 uint64 = 1 << propBit_INCRDECR
	propFlag_QUOTED                                   uint64 = 1 << propBit_QUOTED
	propFlag_SYNTHETIC                                uint64 = 1 << propBit_SYNTHETIC
	propFlag_ADDED_BLOCK                              uint64 = 1 << propBit_ADDED_BLOCK
	propFlag_IS_CONSTANT_NAME                         uint64 = 1 << propBit_IS_CONSTANT_NAME
	propFlag_IS_NAMESPACE                             uint64 = 1 << propBit_IS_NAMESPACE
	propFlag_USE_STRICT                               uint64 = 1 << propBit_USE_STRICT
	propFlag_DIRECT_EVAL                              uint64 = 1 << propBit_DIRECT_EVAL
	propFlag_FREE_CALL                                uint64 = 1 << propBit_FREE_CALL
	propFlag_STATIC_MEMBER                            uint64 = 1 << propBit_STATIC_MEMBER
	propFlag_GENERATOR_FN                             uint64 = 1 << propBit_GENERATOR_FN
	propFlag_ARROW_FN                                 uint64 = 1 << propBit_ARROW_FN
	propFlag_ASYNC_FN                                 uint64 = 1 << propBit_ASYNC_FN
	propFlag_YIELD_ALL                                uint64 = 1 << propBit_YIELD_ALL
	propFlag_EXPORT_DEFAULT                           uint64 = 1 << propBit_EXPORT_DEFAULT
	propFlag_EXPORT_ALL_FROM                          uint64 = 1 << propBit_EXPORT_ALL_FROM
	propFlag_IS_GENERATOR_MARKER                      uint64 = 1 << propBit_IS_GENERATOR_MARKER
	propFlag_IS_GENERATOR_SAFE                        uint64 = 1 << propBit_IS_GENERATOR_SAFE
	propFlag_COMPUTED_PROP_METHOD                     uint64 = 1 << propBit_COMPUTED_PROP_METHOD
	propFlag_COMPUTED_PROP_GETTER                     uint64 = 1 << propBit_COMPUTED_PROP_GETTER
	propFlag_COMPUTED_PROP_SETTER                     uint64 = 1 << propBit_COMPUTED_PROP_SETTER
	propFlag_COMPUTED_PROP_VARIABLE                   uint64 = 1 << propBit_COMPUTED_PROP_VARIABLE
	propFlag_COLOR_FROM_CAST                          uint64 = 1 << propBit_COLOR_FROM_CAST
	propFlag_OPT_ES6_TYPED                            uint64 = 1 << propBit_OPT_ES6_TYPED
	propFlag_CONSTRUCT_SIGNATURE                      uint64 = 1 << propBit_CONSTRUCT_SIGNATURE
	propFlag_NON_INDEXABLE                            uint64 = 1 << propBit_NON_INDEXABLE
	propFlag_GOOG_MODULE                              uint64 = 1 << propBit_GOOG_MODULE
	propFlag_IS_TYPESCRIPT_ABSTRACT                   uint64 = 1 << propBit_IS_TYPESCRIPT_ABSTRACT
	propFlag_DELETED                                  uint64 = 1 << propBit_DELETED
	propFlag_MODULE_ALIAS                             uint64 = 1 << propBit_MODULE_ALIAS
	propFlag_IS_UNUSED_PARAMETER                      uint64 = 1 << propBit_IS_UNUSED_PARAMETER
	propFlag_MODULE_EXPORT                            uint64 = 1 << propBit_MODULE_EXPORT
	propFlag_IS_SHORTHAND_PROPERTY                    uint64 = 1 << propBit_IS_SHORTHAND_PROPERTY
	propFlag_ES6_MODULE                               uint64 = 1 << propBit_ES6_MODULE
	propFlag_START_OF_OPT_CHAIN                       uint64 = 1 << propBit_START_OF_OPT_CHAIN
	propFlag_TRAILING_COMMA                           uint64 = 1 << propBit_TRAILING_COMMA
	propFlag_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION uint64 = 1 << propBit_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
	propFlag_PRIVATE_IDENTIFIER                       uint64 = 1 << propBit_PRIVATE_IDENTIFIER
	propFlag_IS_IN_CLOSURE_UNAWARE_SUBTREE            uint64 = 1 << propBit_IS_IN_CLOSURE_UNAWARE_SUBTREE

	nodePropBoolMask uint64 = (uint64(1) << propBitCount) - 1

	propFlagHasJSDoc uint64 = 1 << 40

	nodePropConstVarShift          = 41
	nodePropConstVarMask    uint64 = uint64(3) << nodePropConstVarShift
	nodePropSideEffectShift        = 43
	nodePropSideEffectMask  uint64 = uint64(15) << nodePropSideEffectShift
)

var nodePropBitMask = [64]uint64{
	Node_Prop_IS_PARENTHESIZED:                         propFlag_IS_PARENTHESIZED,
	Node_Prop_INCRDECR:                                 propFlag_INCRDECR,
	Node_Prop_QUOTED:                                   propFlag_QUOTED,
	Node_Prop_SYNTHETIC:                                propFlag_SYNTHETIC,
	Node_Prop_ADDED_BLOCK:                              propFlag_ADDED_BLOCK,
	Node_Prop_IS_CONSTANT_NAME:                         propFlag_IS_CONSTANT_NAME,
	Node_Prop_IS_NAMESPACE:                             propFlag_IS_NAMESPACE,
	Node_Prop_USE_STRICT:                               propFlag_USE_STRICT,
	Node_Prop_DIRECT_EVAL:                              propFlag_DIRECT_EVAL,
	Node_Prop_FREE_CALL:                                propFlag_FREE_CALL,
	Node_Prop_STATIC_MEMBER:                            propFlag_STATIC_MEMBER,
	Node_Prop_GENERATOR_FN:                             propFlag_GENERATOR_FN,
	Node_Prop_ARROW_FN:                                 propFlag_ARROW_FN,
	Node_Prop_ASYNC_FN:                                 propFlag_ASYNC_FN,
	Node_Prop_YIELD_ALL:                                propFlag_YIELD_ALL,
	Node_Prop_EXPORT_DEFAULT:                           propFlag_EXPORT_DEFAULT,
	Node_Prop_EXPORT_ALL_FROM:                          propFlag_EXPORT_ALL_FROM,
	Node_Prop_IS_GENERATOR_MARKER:                      propFlag_IS_GENERATOR_MARKER,
	Node_Prop_IS_GENERATOR_SAFE:                        propFlag_IS_GENERATOR_SAFE,
	Node_Prop_COMPUTED_PROP_METHOD:                     propFlag_COMPUTED_PROP_METHOD,
	Node_Prop_COMPUTED_PROP_GETTER:                     propFlag_COMPUTED_PROP_GETTER,
	Node_Prop_COMPUTED_PROP_SETTER:                     propFlag_COMPUTED_PROP_SETTER,
	Node_Prop_COMPUTED_PROP_VARIABLE:                   propFlag_COMPUTED_PROP_VARIABLE,
	Node_Prop_COLOR_FROM_CAST:                          propFlag_COLOR_FROM_CAST,
	Node_Prop_OPT_ES6_TYPED:                            propFlag_OPT_ES6_TYPED,
	Node_Prop_CONSTRUCT_SIGNATURE:                      propFlag_CONSTRUCT_SIGNATURE,
	Node_Prop_NON_INDEXABLE:                            propFlag_NON_INDEXABLE,
	Node_Prop_GOOG_MODULE:                              propFlag_GOOG_MODULE,
	Node_Prop_IS_TYPESCRIPT_ABSTRACT:                   propFlag_IS_TYPESCRIPT_ABSTRACT,
	Node_Prop_DELETED:                                  propFlag_DELETED,
	Node_Prop_MODULE_ALIAS:                             propFlag_MODULE_ALIAS,
	Node_Prop_IS_UNUSED_PARAMETER:                      propFlag_IS_UNUSED_PARAMETER,
	Node_Prop_MODULE_EXPORT:                            propFlag_MODULE_EXPORT,
	Node_Prop_IS_SHORTHAND_PROPERTY:                    propFlag_IS_SHORTHAND_PROPERTY,
	Node_Prop_ES6_MODULE:                               propFlag_ES6_MODULE,
	Node_Prop_START_OF_OPT_CHAIN:                       propFlag_START_OF_OPT_CHAIN,
	Node_Prop_TRAILING_COMMA:                           propFlag_TRAILING_COMMA,
	Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION: propFlag_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION,
	Node_Prop_PRIVATE_IDENTIFIER:                       propFlag_PRIVATE_IDENTIFIER,
	Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE:            propFlag_IS_IN_CLOSURE_UNAWARE_SUBTREE,
}

var nodeBitToProp = [propBitCount]Node_Prop{
	propBit_IS_PARENTHESIZED:                         Node_Prop_IS_PARENTHESIZED,
	propBit_INCRDECR:                                 Node_Prop_INCRDECR,
	propBit_QUOTED:                                   Node_Prop_QUOTED,
	propBit_SYNTHETIC:                                Node_Prop_SYNTHETIC,
	propBit_ADDED_BLOCK:                              Node_Prop_ADDED_BLOCK,
	propBit_IS_CONSTANT_NAME:                         Node_Prop_IS_CONSTANT_NAME,
	propBit_IS_NAMESPACE:                             Node_Prop_IS_NAMESPACE,
	propBit_USE_STRICT:                               Node_Prop_USE_STRICT,
	propBit_DIRECT_EVAL:                              Node_Prop_DIRECT_EVAL,
	propBit_FREE_CALL:                                Node_Prop_FREE_CALL,
	propBit_STATIC_MEMBER:                            Node_Prop_STATIC_MEMBER,
	propBit_GENERATOR_FN:                             Node_Prop_GENERATOR_FN,
	propBit_ARROW_FN:                                 Node_Prop_ARROW_FN,
	propBit_ASYNC_FN:                                 Node_Prop_ASYNC_FN,
	propBit_YIELD_ALL:                                Node_Prop_YIELD_ALL,
	propBit_EXPORT_DEFAULT:                           Node_Prop_EXPORT_DEFAULT,
	propBit_EXPORT_ALL_FROM:                          Node_Prop_EXPORT_ALL_FROM,
	propBit_IS_GENERATOR_MARKER:                      Node_Prop_IS_GENERATOR_MARKER,
	propBit_IS_GENERATOR_SAFE:                        Node_Prop_IS_GENERATOR_SAFE,
	propBit_COMPUTED_PROP_METHOD:                     Node_Prop_COMPUTED_PROP_METHOD,
	propBit_COMPUTED_PROP_GETTER:                     Node_Prop_COMPUTED_PROP_GETTER,
	propBit_COMPUTED_PROP_SETTER:                     Node_Prop_COMPUTED_PROP_SETTER,
	propBit_COMPUTED_PROP_VARIABLE:                   Node_Prop_COMPUTED_PROP_VARIABLE,
	propBit_COLOR_FROM_CAST:                          Node_Prop_COLOR_FROM_CAST,
	propBit_OPT_ES6_TYPED:                            Node_Prop_OPT_ES6_TYPED,
	propBit_CONSTRUCT_SIGNATURE:                      Node_Prop_CONSTRUCT_SIGNATURE,
	propBit_NON_INDEXABLE:                            Node_Prop_NON_INDEXABLE,
	propBit_GOOG_MODULE:                              Node_Prop_GOOG_MODULE,
	propBit_IS_TYPESCRIPT_ABSTRACT:                   Node_Prop_IS_TYPESCRIPT_ABSTRACT,
	propBit_DELETED:                                  Node_Prop_DELETED,
	propBit_MODULE_ALIAS:                             Node_Prop_MODULE_ALIAS,
	propBit_IS_UNUSED_PARAMETER:                      Node_Prop_IS_UNUSED_PARAMETER,
	propBit_MODULE_EXPORT:                            Node_Prop_MODULE_EXPORT,
	propBit_IS_SHORTHAND_PROPERTY:                    Node_Prop_IS_SHORTHAND_PROPERTY,
	propBit_ES6_MODULE:                               Node_Prop_ES6_MODULE,
	propBit_START_OF_OPT_CHAIN:                       Node_Prop_START_OF_OPT_CHAIN,
	propBit_TRAILING_COMMA:                           Node_Prop_TRAILING_COMMA,
	propBit_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION: Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION,
	propBit_PRIVATE_IDENTIFIER:                       Node_Prop_PRIVATE_IDENTIFIER,
	propBit_IS_IN_CLOSURE_UNAWARE_SUBTREE:            Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE,
}

// nodeExtBox holds the rarely-set parts of a node, in its store's cold-data
// map (node_tables.go).
type nodeExtBox struct {
	jsDocInfo      *JSDocInfo         // canonical *JSDocInfo (valid iff propFlagHasJSDoc != 0 && jsDocInfo != nil)
	propListHead   *Node_PropListItem // remaining cold properties
	typeBeforeCast uint32             // node type handle of the canonical TYPE_BEFORE_CAST (0 iff not present as JSType)
	length         int32              // the length, when Node.length16 is 0xFFFF
	cookedID       rhinoStringID      // the cooked string of a TEMPLATELIT_STRING, when it is not the raw one
	cookedNil      bool               // the TEMPLATELIT_STRING has no cooked string
}

// The SOURCE_FILE property of a node is an id into these tables: each
// PutProp(SOURCE_FILE, file) takes a new id, holding the property's list
// item (shared by every node that copies the id, as Java shares the item)
// and the file's name. 16M ids: a few per input file.
const sourceFileIDChunks = 1 << 14

var (
	sourceFilePropMu        sync.Mutex
	sourceFilePropCount     atomic.Uint32
	sourceFilePropChunks    [sourceFileIDChunks]atomic.Pointer[[1024]Node_PropListItem]
	sourceFileNamePtrChunks [sourceFileIDChunks]atomic.Pointer[[1024]*string]
)

func sourceFilePropItem(id uint32) *Node_PropListItem {
	idx := id - 1
	return &sourceFilePropChunks[idx>>10].Load()[idx&1023]
}

// allocSourceFileID returns a new source-file id for file.
func allocSourceFileID(file any) uint32 {
	sourceFilePropMu.Lock()
	id := sourceFilePropCount.Load() + 1
	idx := id - 1
	chunkIdx := idx >> 10
	if chunkIdx >= sourceFileIDChunks {
		sourceFilePropMu.Unlock()
		panic(fmt.Sprintf("rhino: more than %d source-file ids", sourceFileIDChunks*1024))
	}
	chunk := sourceFilePropChunks[chunkIdx].Load()
	if chunk == nil {
		chunk = new([1024]Node_PropListItem)
		sourceFilePropChunks[chunkIdx].Store(chunk)
	}
	nameChunk := sourceFileNamePtrChunks[chunkIdx].Load()
	if nameChunk == nil {
		nameChunk = new([1024]*string)
		sourceFileNamePtrChunks[chunkIdx].Store(nameChunk)
	}
	var namePtr *string
	if fast, ok := file.(StaticSourceFileFast); ok {
		namePtr = fast.GetNamePtr()
	} else if sf, ok := file.(StaticSourceFile); ok {
		s := sf.GetName()
		namePtr = &s
	}
	nameChunk[idx&1023] = namePtr
	chunk[idx&1023] = Node_PropListItem{
		objectValue: file,
		propType:    int8(Node_Prop_SOURCE_FILE),
		isObject:    true,
		shared:      true,
	}
	sourceFilePropCount.Store(id)
	sourceFilePropMu.Unlock()
	return id
}

// jsDocAux is the node's JSDocInfo when it is stored directly in the aux
// entry (no ext box).
func (n Node) jsDocAux() *JSDocInfo {
	if a, ok := n.auxOf(); ok {
		return a.jsDoc
	}
	return nil
}

// extBoxOrNil is the node's ext box, or nil if the node has none. Read
// paths use it; mutators use extBox.
func (n Node) extBoxOrNil() *nodeExtBox {
	if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
		return nil
	}
	if a, ok := n.auxOf(); ok {
		return a.ext
	}
	return nil
}

// extBox is the ext box of a node that has one (kindFlagHasExtBox set).
func (n Node) extBox() *nodeExtBox {
	if box := n.extBoxOrNil(); box != nil {
		return box
	}
	box := &nodeExtBox{}
	n.setAux(nodeAux{ext: box})
	return box
}

func (n Node) ensureExtBox() *nodeExtBox {
	if (n.rec().kindFlags & kindFlagHasExtBox) != 0 {
		return n.extBox()
	}
	box := &nodeExtBox{}
	if (n.props() & propFlagHasJSDoc) != 0 {
		box.jsDocInfo = n.jsDocAux()
	}
	n.rec().kindFlags |= kindFlagHasExtBox
	n.setAux(nodeAux{ext: box})
	return box
}

// maybeCollapseExtBox drops the box when only the JSDocInfo (or nothing) is
// left in it: the aux entry then refers to the JSDocInfo directly.
func (n Node) maybeCollapseExtBox(box *nodeExtBox) {
	if box.typeBeforeCast != 0 || box.propListHead != nil || n.extBoxHoldsNodeFields(box) {
		return
	}
	n.rec().kindFlags &^= kindFlagHasExtBox
	if (n.props()&propFlagHasJSDoc) != 0 && box.jsDocInfo != nil {
		n.setAux(nodeAux{jsDoc: box.jsDocInfo})
	} else {
		n.setAux(nodeAux{})
	}
}

// extBoxHoldsNodeFields reports whether box holds a length or a cooked
// string that does not fit in the Node itself.
func (n Node) extBoxHoldsNodeFields(box *nodeExtBox) bool {
	return n.rec().length16 == 0xFFFF || box.cookedID != 0 || box.cookedNil
}

func (n Node) coldPropListHead() *Node_PropListItem {
	if box := n.extBoxOrNil(); box != nil {
		return box.propListHead
	}
	return nil
}

func (n Node) setColdPropListHead(head *Node_PropListItem) {
	if head == nil {
		if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
			return
		}
		box := n.extBox()
		box.propListHead = nil
		n.maybeCollapseExtBox(box)
		return
	}
	box := n.ensureExtBox()
	box.propListHead = head
}

func (n Node) typeBeforeCastBase() *JSTypeBase {
	if box := n.extBoxOrNil(); box != nil {
		return n.store().typeAt(box.typeBeforeCast)
	}
	return nil
}

func (n Node) copyTemplateColdProps(template Node) {
	if tbox := template.extBoxOrNil(); tbox != nil {
		if tbox.propListHead != nil {
			tbox.propListHead.markShared()
		}
		if tbox.jsDocInfo != nil && tbox.typeBeforeCast == 0 && tbox.propListHead == nil {
			n.setAux(nodeAux{jsDoc: tbox.jsDocInfo})
			return
		}
		if tbox.jsDocInfo != nil || tbox.typeBeforeCast != 0 || tbox.propListHead != nil {
			box := n.ensureExtBox()
			box.jsDocInfo = tbox.jsDocInfo
			box.typeBeforeCast = n.store().translateTypeHandle(template.store(), tbox.typeBeforeCast)
			box.propListHead = tbox.propListHead
		}
		return
	}
	if (template.props() & propFlagHasJSDoc) != 0 {
		n.setAux(nodeAux{jsDoc: template.jsDocAux()})
	}
}

func (n Node) nodeKind() node_kind {
	return node_kind(n.rec().kindFlags & kindMask)
}

func (n Node) setNodeKind(k node_kind) {
	n.rec().kindFlags = (n.rec().kindFlags &^ kindMask) | uint8(k)
}

var node_kindClassNames = [...]string{
	"com.google.javascript.rhino.Node",
	"com.google.javascript.rhino.Node$StringNode",
	"com.google.javascript.rhino.Node$NumberNode",
	"com.google.javascript.rhino.Node$BigIntNode",
	"com.google.javascript.rhino.Node$TemplateLiteralSubstringNode",
}

// GetClassName is `getClass().getName()`: the name of the Java class this
// node is an instance of (Node or one of its folded private subclasses).
//
// Go-only: language-adapter — Java code compares `node.getClass()` with
// `Node.class` (CodeGenerator); the folded subclasses have no Go type of
// their own.
func (n Node) GetClassName() string {
	return node_kindClassNames[n.nodeKind()]
}

//go:noinline
func failNodeCastKind(got, want node_kind) {
	panic(jutil.NewClassCastException(node_kindClassNames[got], node_kindClassNames[want]))
}

// node_castKind is the Java cast `((XNode) this)`: it throws a
// ClassCastException unless n is of the given kind.
func (n Node) node_castKind(k node_kind) {
	if got := n.nodeKind(); got != k {
		failNodeCastKind(got, k)
	}
}

//go:noinline
func failNodeCast(value any, zero any) {
	panic(jutil.NewClassCastException(fmt.Sprintf("%T", value), fmt.Sprintf("%T", zero)))
}

// node_cast is the Java reference cast `(T) value` of a property value:
// null stays null (the zero T), other values must be of type T.
func node_cast[T any](value any) T {
	if value == nil {
		var zero T
		return zero
	}
	t, ok := value.(T)
	if !ok {
		var zero T
		failNodeCast(value, any(zero))
	}
	return t
}

// NumberNode(double)
//
// Java: Node.NumberNode(double)
func newNode_NumberNode(number float64) Node {
	n := NewNode(Token_NUMBER)
	n.rec().kindFlags = uint8(node_kindNumber)
	n.SetDouble(number)
	return n
}

// BigIntNode(BigInteger)
//
// Java: Node.BigIntNode(BigInteger)
func newNode_BigIntNode(bigint *big.Int) Node {
	n := NewNode(Token_BIGINT)
	n.rec().kindFlags = uint8(node_kindBigInt)
	n.SetBigInt(bigint)
	return n
}

// Only for cloneNode.
//
// Java: Node.StringNode(Token)
func newNode_StringNodeForClone(token Token) Node {
	n := NewNode(token)
	n.rec().kindFlags = uint8(node_kindString)
	return n
}

// Java: Node.StringNode(Token, String)
func newNode_StringNode(token Token, str string) Node {
	n := NewNode(token)
	n.rec().kindFlags = uint8(node_kindString)
	n.SetString(str)
	return n
}

// Java: Node.StringNode(Token, RhinoStringPool.LazyInternedStringList, int)
func newNode_StringNodeFromStringPool(token Token, stringPool *RhinoStringPool_LazyInternedStringList, offset int) Node {
	n := NewNode(token)
	n.rec().kindFlags = uint8(node_kindString)
	n.SetStringFromStringPool(stringPool, offset)
	return n
}

// Java: Node.TemplateLiteralSubstringNode(String, String)
func newNode_TemplateLiteralSubstringNode(cooked *string, raw string) Node {
	n := NewNode(Token_TEMPLATELIT_STRING)
	n.rec().kindFlags = uint8(node_kindTemplateLiteralSubstring)
	var cookedID rhinoStringID
	if cooked != nil {
		cookedID = rhinoStringPoolIntern(*cooked)
	}
	rawID := rhinoStringPoolIntern(raw)
	n.initTemplateLitPayload(cookedID, rawID)
	return n
}

// Java: Node.TemplateLiteralSubstringNode(RhinoStringPool.LazyInternedStringList, int, int)
func newNode_TemplateLiteralSubstringNodeFromStringPool(
	stringPool *RhinoStringPool_LazyInternedStringList,
	cookedOffsetOrNegativeOne int,
	rawOffset int,
) Node {
	n := NewNode(Token_TEMPLATELIT_STRING)
	n.rec().kindFlags = uint8(node_kindTemplateLiteralSubstring)
	var cookedID rhinoStringID
	if cookedOffsetOrNegativeOne != -1 {
		cookedID = stringPool.GetID(cookedOffsetOrNegativeOne)
	}
	rawID := stringPool.GetID(rawOffset)
	n.initTemplateLitPayload(cookedID, rawID)
	return n
}

// Node_PropListItem is an element of a node's property list.
//
// Go-only: language-adapter — the abstract PropListItem and its two
// subclasses ObjectPropListItem and IntPropListItem (which only differ in
// the payload and accessors) are one struct; isObject records the class.
//
// Go-only: performance — the small fields share the last word: 32 bytes
// instead of 40 (a 48-byte allocation).
type Node_PropListItem struct {
	next        *Node_PropListItem
	objectValue any
	intValue    int32
	propType    int8
	isObject    bool
	shared      bool
}

func (p *Node_PropListItem) markShared() {
	for ; p != nil && !p.shared; p = p.next {
		p.shared = true
	}
}

var (
	staticBoolPropItems       [64]Node_PropListItem
	staticConstVarPropItems   [4]Node_PropListItem
	staticSideEffectPropItems [16]Node_PropListItem
)

func init() {
	for i := range staticBoolPropItems {
		staticBoolPropItems[i] = Node_PropListItem{
			propType: int8(i),
			intValue: 1,
			shared:   true,
		}
	}
	for i := range staticConstVarPropItems {
		staticConstVarPropItems[i] = Node_PropListItem{
			propType: int8(Node_Prop_CONSTANT_VAR_FLAGS),
			intValue: int32(i),
			shared:   true,
		}
	}
	for i := range staticSideEffectPropItems {
		staticSideEffectPropItems[i] = Node_PropListItem{
			propType: int8(Node_Prop_SIDE_EFFECT_FLAGS),
			intValue: int32(i),
			shared:   true,
		}
	}
}

// Java: Node.ObjectPropListItem(byte, Object, PropListItem)
func NewNode_ObjectPropListItem(propType int8, objectValue any, next *Node_PropListItem) *Node_PropListItem {
	obj := jutil.CheckNotNull(objectValue)
	return &Node_PropListItem{
		propType:    propType,
		next:        next,
		isObject:    true,
		objectValue: obj,
	}
}

// Java: Node.IntPropListItem(byte, int, PropListItem)
func NewNode_IntPropListItem(propType int8, intValue int32, next *Node_PropListItem) *Node_PropListItem {
	jutil.CheckState(intValue != 0)
	return &Node_PropListItem{
		propType: propType,
		next:     next,
		intValue: intValue,
	}
}

// Java: Node.PropListItem.getIntValue()
func (p *Node_PropListItem) GetIntValue() int32 {
	if p.isObject {
		panic(jutil.NewUnsupportedOperationException())
	}
	return p.intValue
}

// Java: Node.PropListItem.getObjectValue()
func (p *Node_PropListItem) GetObjectValue() any {
	if !p.isObject {
		panic(jutil.NewUnsupportedOperationException())
	}
	return p.objectValue
}

// Java: Node.PropListItem.toString()
func (p *Node_PropListItem) String() string {
	if p.isObject {
		return jutil.ToJavaString(p.objectValue)
	}
	return strconv.FormatInt(int64(p.intValue), 10)
}

// Java: Node.PropListItem.chain(PropListItem)
func (p *Node_PropListItem) Chain(next *Node_PropListItem) *Node_PropListItem {
	if p.isObject {
		return NewNode_ObjectPropListItem(p.propType, p.objectValue, next)
	}
	return NewNode_IntPropListItem(p.propType, p.intValue, next)
}

// Java: Node(Token)
func NewNode(token Token) Node {
	n := allocNode()
	n.rec().token = token
	n.rec().linenoCharno = -1
	return n
}

// AllocNode returns a fresh zeroed node (to be initialized with one of the
// Init methods).
func AllocNode() Node {
	return allocNode()
}

// Init initializes a zero Node with token and default source position.
func (n Node) Init(token Token) Node {
	n.rec().token = token
	n.rec().linenoCharno = -1
	return n
}

// InitFromTemplate initializes a zero Node from a slab with token and template properties.
func (n Node) InitFromTemplate(token Token, template Node) Node {
	n.rec().token = token
	n.rec().linenoCharno = -1
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	return n
}

// InitWithChildFromTemplate initializes a zero Node from a slab with 1 child and template properties.
func (n Node) InitWithChildFromTemplate(token Token, child Node, template Node) Node {
	if child.rec().parent != 0 || child.rec().next != 0 || child.rec().previous != 0 || template.hasAux() {
		return n.initWithChildFromTemplateSlow(token, child, template)
	}
	n.rec().token = token
	n.rec().linenoCharno = -1
	n.setPropsAndSourceFile(template)
	n.rec().first = n.ref(child)
	child.rec().previous = child.ref(child)
	child.rec().parent = child.ref(n)
	return n
}

//go:noinline
func (n Node) initWithChildFromTemplateSlow(token Token, child Node, template Node) Node {
	n.rec().token = token
	n.rec().linenoCharno = -1
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	n.rec().first = n.ref(child)
	child.CheckDetached()
	child.rec().previous = child.ref(child)
	child.rec().parent = child.ref(n)
	return n
}

// InitWith2ChildrenFromTemplate initializes a zero Node from a slab with 2 children and template properties.
func (n Node) InitWith2ChildrenFromTemplate(token Token, left Node, right Node, template Node) Node {
	n.rec().token = token
	n.rec().linenoCharno = -1
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	n.rec().first = n.ref(left)

	left.CheckDetached()
	left.rec().next = left.ref(right)
	left.rec().previous = left.ref(right)
	left.rec().parent = left.ref(n)

	right.CheckDetached()
	right.rec().previous = right.ref(left)
	right.rec().parent = right.ref(n)
	return n
}

// InitWith3ChildrenFromTemplate initializes a zero Node from a slab with 3 children and template properties.
func (n Node) InitWith3ChildrenFromTemplate(token Token, left Node, mid Node, right Node, template Node) Node {
	n.rec().token = token
	n.rec().linenoCharno = -1
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	n.rec().first = n.ref(left)

	left.CheckDetached()
	left.rec().next = left.ref(mid)
	left.rec().previous = left.ref(right)
	left.rec().parent = left.ref(n)

	mid.CheckDetached()
	mid.rec().next = mid.ref(right)
	mid.rec().previous = mid.ref(left)
	mid.rec().parent = mid.ref(n)

	right.CheckDetached()
	right.rec().previous = right.ref(mid)
	right.rec().parent = right.ref(n)
	return n
}

//go:noinline
func (n Node) initStringPayloadAndCold(str string, template Node) {
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	n.rec().payload = uint32(rhinoStringPoolIntern(str))
}

// InitStringFromTemplate initializes a zero Node from a slab as a StringNode with template properties.
func (n Node) InitStringFromTemplate(token Token, str string, template Node) Node {
	n.rec().token = token
	n.rec().linenoCharno = -1
	n.rec().kindFlags = uint8(node_kindString)
	n.setPropsAndSourceFile(template)
	n.initStringPayloadAndCold(str, template)
	return n
}

// InitStringIDFromTemplate is InitStringFromTemplate for a string already
// interned (its pool id).
func (n Node) InitStringIDFromTemplate(token Token, id rhinoStringID, template Node) Node {
	n.rec().token = token
	n.rec().linenoCharno = -1
	n.rec().kindFlags = uint8(node_kindString)
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	n.rec().payload = uint32(id)
	return n
}

// SetStringID sets the interned string payload on a freshly initialized Node.
func (n Node) SetStringID(id rhinoStringID) Node {
	n.rec().kindFlags |= uint8(node_kindString)
	n.rec().payload = uint32(id)
	return n
}

// SetNumberRaw sets the float64 payload on a freshly initialized Node.
func (n Node) SetNumberRaw(number float64) Node {
	n.rec().kindFlags |= uint8(node_kindNumber)
	n.rec().payload = numberPayload(number)
	return n
}

// InitStringFromStringPool initializes a zero Node from a slab as a StringNode from a LazyInternedStringList.
func (n Node) InitStringFromStringPool(token Token, stringPool *RhinoStringPool_LazyInternedStringList, offset int, template Node) Node {
	id := stringPool.GetID(offset)
	return n.InitFromTemplate(token, template).SetStringID(id)
}

// InitTemplateLitStringFromStringPool initializes a zero Node from a slab as a TemplateLiteralSubstringNode from a LazyInternedStringList.
func (n Node) InitTemplateLitStringFromStringPool(
	stringPool *RhinoStringPool_LazyInternedStringList,
	cookedOffsetOrNegativeOne int,
	rawOffset int,
	template Node,
) Node {
	n.rec().token = Token_TEMPLATELIT_STRING
	n.rec().linenoCharno = -1
	n.rec().kindFlags = uint8(node_kindTemplateLiteralSubstring)
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	var cookedID rhinoStringID
	if cookedOffsetOrNegativeOne != -1 {
		cookedID = stringPool.GetID(cookedOffsetOrNegativeOne)
	}
	rawID := stringPool.GetID(rawOffset)
	n.initTemplateLitPayload(cookedID, rawID)
	return n
}

// InitNumberFromTemplate initializes a zero Node from a slab as a NumberNode with template properties.
func (n Node) InitNumberFromTemplate(number float64, template Node) Node {
	n.InitFromTemplate(Token_NUMBER, template)
	n.rec().kindFlags |= uint8(node_kindNumber)
	n.rec().payload = numberPayload(number)
	return n
}

// InitBigIntFromTemplate initializes a zero Node from a slab as a BigIntNode with template properties.
func (n Node) InitBigIntFromTemplate(bigint *big.Int, template Node) Node {
	n.rec().token = Token_BIGINT
	n.rec().linenoCharno = -1
	n.rec().kindFlags = uint8(node_kindBigInt)
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	n.SetBigInt(bigint)
	return n
}

// InitTemplateLitStringFromTemplate initializes a zero Node from a slab as a TemplateLiteralSubstringNode with template properties.
func (n Node) InitTemplateLitStringFromTemplate(cooked *string, raw string, template Node) Node {
	n.rec().token = Token_TEMPLATELIT_STRING
	n.rec().linenoCharno = -1
	n.rec().kindFlags = uint8(node_kindTemplateLiteralSubstring)
	n.setPropsAndSourceFile(template)
	if template.hasAux() {
		n.copyTemplateColdProps(template)
	}
	var cookedID rhinoStringID
	if cooked != nil {
		cookedID = rhinoStringPoolIntern(*cooked)
	}
	rawID := rhinoStringPoolIntern(raw)
	n.initTemplateLitPayload(cookedID, rawID)
	return n
}

// Java: Node(Token, Node)
func NewNodeWithChild(token Token, child Node) Node {
	n := NewNode(token)
	n.rec().first = n.ref(child)

	child.CheckDetached()
	// child.next remains null;
	child.rec().previous = child.ref(child)
	child.rec().parent = child.ref(n)
	return n
}

// Java: Node(Token, Node, Node)
func NewNodeWith2Children(token Token, left Node, right Node) Node {
	n := NewNode(token)
	n.rec().first = n.ref(left)

	left.CheckDetached()
	left.rec().next = left.ref(right)
	left.rec().previous = left.ref(right)
	left.rec().parent = left.ref(n)

	right.CheckDetached()
	// right.next remains null;
	right.rec().previous = right.ref(left)
	right.rec().parent = right.ref(n)
	return n
}

// Java: Node(Token, Node, Node, Node)
func NewNodeWith3Children(token Token, left Node, mid Node, right Node) Node {
	n := NewNode(token)
	n.rec().first = n.ref(left)

	left.CheckDetached()
	left.rec().next = left.ref(mid)
	left.rec().previous = left.ref(right)
	left.rec().parent = left.ref(n)

	mid.CheckDetached()
	mid.rec().next = mid.ref(right)
	mid.rec().previous = mid.ref(left)
	mid.rec().parent = mid.ref(n)

	right.CheckDetached()
	// right.next remains null;
	right.rec().previous = right.ref(mid)
	right.rec().parent = right.ref(n)
	return n
}

// Java: Node.newNumber(double)
func Node_NewNumber(number float64) Node {
	return newNode_NumberNode(number)
}

// Java: Node.newBigInt(BigInteger)
func Node_NewBigInt(bigint *big.Int) Node {
	return newNode_BigIntNode(bigint)
}

// Java: Node.newString(String)
func Node_NewString(str string) Node {
	return newNode_StringNode(Token_STRINGLIT, str)
}

// Java: Node.newString(Token, String)
func Node_NewStringWithToken(token Token, str string) Node {
	return newNode_StringNode(token, str)
}

// Java: Node.newString(Token, RhinoStringPool.LazyInternedStringList, int)
func Node_NewStringFromStringPool(token Token, stringPool *RhinoStringPool_LazyInternedStringList, offset int) Node {
	return newNode_StringNodeFromStringPool(token, stringPool, offset)
}

// Java: Node.newTemplateLitString(String, String)
func Node_NewTemplateLitString(cooked *string, raw string) Node {
	return newNode_TemplateLiteralSubstringNode(cooked, raw)
}

// Java: Node.newTemplateLitString(RhinoStringPool.LazyInternedStringList, int, int)
func Node_NewTemplateLitStringFromStringPool(
	stringPool *RhinoStringPool_LazyInternedStringList,
	cookedOffsetOrNegativeOne int,
	rawOffset int,
) Node {
	return newNode_TemplateLiteralSubstringNodeFromStringPool(stringPool, cookedOffsetOrNegativeOne, rawOffset)
}

// Java: Node.getToken()
func (n Node) GetToken() Token {
	return n.rec().token
}

// Java: Node.setToken(Token)
func (n Node) SetToken(token Token) {
	n.rec().token = token
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
}

// Java: Node.hasChildren()
func (n Node) HasChildren() bool {
	return n.rec().first != 0
}

// Java: Node.getOnlyChild()
func (n Node) GetOnlyChild() Node {
	jutil.CheckState(n.HasOneChild())
	return n.firstNode()
}

// Java: Node.getFirstChild()
func (n Node) GetFirstChild() Node {
	return n.firstNode()
}

// GetFirstFirstChild gets the first child of the first child. This method
// assumes that the first child exists.
//
// Java: Node.getFirstFirstChild()
func (n Node) GetFirstFirstChild() Node {
	return n.firstNode().firstNode()
}

// Java: Node.getSecondChild()
func (n Node) GetSecondChild() Node {
	return n.firstNode().nextNode()
}

// Java: Node.getLastChild()
func (n Node) GetLastChild() Node {
	// The first child's previous is the last child (a circular list).
	if first := n.firstNode(); first != NilNode {
		return first.prevNode()
	}
	return NilNode
}

// Java: Node.getNext()
func (n Node) GetNext() Node {
	return n.nextNode()
}

// Java: Node.getPrevious()
func (n Node) GetPrevious() Node {
	if uint32(n) == n.parentNode().rec().first {
		return NilNode
	}
	return n.prevNode()
}

// Java: Node.setClosureUnawareShadow(Node)
func (n Node) SetClosureUnawareShadow(shadowRoot Node) {
	jutil.CheckStateMsg(n.rec().first == 0, "Cannot set shadow root on a node with children")
	jutil.CheckStateMsg(n.rec().token == Token_NAME, "Only NAME nodes can be used as shadows")
	n.PutProp(Node_Prop_CLOSURE_UNAWARE_SHADOW, shadowRoot)
}

// Java: Node.getClosureUnawareShadow()
func (n Node) GetClosureUnawareShadow() Node {
	if n.rec().token != Token_NAME || (n.rec().kindFlags&kindFlagHasExtBox) == 0 {
		// As checked in setClosureUnawareShadow, only NAME nodes can have shadow
		// content. We short-circuit the getProp call here because that can be
		// very expensive for nodes that have long proplists (and we know it
		// wouldn't find this prop anyways).
		return NilNode
	}
	return node_cast[Node](n.GetProp(Node_Prop_CLOSURE_UNAWARE_SHADOW))
}

// Java: Node.setIsInClosureUnawareSubtree(boolean)
func (n Node) SetIsInClosureUnawareSubtree(value bool) {
	hasOnlySourceFile := n.props() == 0 &&
		n.typeBeforeCastBase() == nil &&
		(n.nodeKind() != node_kindNode || uint32(n.rec().payload) == 0) &&
		n.rec().sourceFile != 0 && n.coldPropListHead() == nil
	if !hasOnlySourceFile {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.StringFormat(
				"Cannot set IS_IN_CLOSURE_UNAWARE_SUBTREE property on a node which proplist is not a"+
					" single SOURCE_FILE prop: %s %s",
				n, n.GetPropListDebugString())))
	}
	n.PutBooleanProp(Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE, value)
}

// Java: Node.getIsInClosureUnawareSubtree()
func (n Node) GetIsInClosureUnawareSubtree() bool {
	return (n.props() & propFlag_IS_IN_CLOSURE_UNAWARE_SUBTREE) != 0
}

// GetChildAtIndex gets the ith child, note that this is O(N) where N is the
// number of children.
//
// Java: Node.getChildAtIndex(int)
func (n Node) GetChildAtIndex(i int) Node {
	c := n.firstNode()
	for i > 0 {
		c = c.nextNode()
		i--
	}
	return c
}

// GetIndexOfChild gets the index of a child, note that this is O(N) where N
// is the number of children.
//
// Java: Node.getIndexOfChild(Node)
func (n Node) GetIndexOfChild(child Node) int {
	c := n.firstNode()
	i := 0
	for !c.IsNil() {
		if child == c {
			return i
		}

		c = c.nextNode()
		i++
	}
	return -1
}

// Java: Node.addChildToFront(Node)
func (n Node) AddChildToFront(child Node) {
	jutil.CheckArgument(child.rec().parent == 0)
	jutil.CheckArgument(child.rec().next == 0)
	jutil.CheckArgument(child.rec().previous == 0)
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	child.rec().parent = child.ref(n)
	child.rec().next = child.ref(n.firstNode())
	if n.rec().first == 0 {
		// NOTE: child.next remains null
		child.rec().previous = child.ref(child)
	} else {
		last := n.firstNode().prevNode()
		// NOTE: last.next remains null
		child.rec().previous = child.ref(last)
		child.rec().next = child.ref(n.firstNode())
		n.firstNode().rec().previous = n.firstNode().ref(child)
	}
	n.rec().first = n.ref(child)
}

// Java: Node.addChildToBack(Node)
func (n Node) AddChildToBack(child Node) {
	jutil.CheckArgumentMsg3(
		child.rec().parent == 0,
		"Cannot add already-owned child node.\nChild: %s\nExisting parent: %s\nNew parent: %s",
		child,
		child.parentNode(),
		n)

	jutil.CheckArgument(child.rec().next == 0)
	jutil.CheckArgument(child.rec().previous == 0)
	n.rec().kindFlags &^= kindFlagNoBlockDecl

	if n.rec().first == 0 {
		// NOTE: child.next remains null
		child.rec().previous = child.ref(child)
		n.rec().first = n.ref(child)
	} else {
		last := n.firstNode().prevNode()
		last.rec().next = last.ref(child)
		// NOTE: child.next remains null
		child.rec().previous = child.ref(last)
		n.firstNode().rec().previous = n.firstNode().ref(child)
	}

	child.rec().parent = child.ref(n)
}

// AddChildrenToFront adds all children to the front of this node.
//
// children is the first of a list of sibling nodes who have no parent.
// NOTE: Usually you would get this argument from a removeChildren() call. A
// single detached node will not work because its sibling pointers will not
// be correctly initialized.
//
// Java: Node.addChildrenToFront(Node)
func (n Node) AddChildrenToFront(children Node) {
	if children.IsNil() {
		return // removeChildren() returns null when there are none
	}
	// NOTE: If there is only one sibling, its previous pointer must point to
	// itself. Null indicates a fully detached node.
	CheckNotNullMsg1(children.prevNode(), "%s", children)
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	for child := children; !child.IsNil(); child = child.nextNode() {
		jutil.CheckArgument(child.rec().parent == 0)
		child.rec().parent = child.ref(n)
	}

	lastSib := children.prevNode()
	if n.rec().first != 0 {
		last := n.firstNode().prevNode()
		// NOTE: last.next remains null
		children.rec().previous = children.ref(last)
		lastSib.rec().next = lastSib.ref(n.firstNode())
		n.firstNode().rec().previous = n.firstNode().ref(lastSib)
	}
	n.rec().first = n.ref(children)
}

// Java: Node.addChildrenToBack(Node)
func (n Node) AddChildrenToBack(children Node) {
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	n.AddChildrenAfter(children, n.GetLastChild())
}

// Java: Node.insertAfter(Node)
func (n Node) InsertAfter(existing Node) {
	existing.CheckAttached()
	n.CheckDetached()

	existingParent := existing.parentNode()
	existingNext := existing.nextNode()
	existingParent.rec().kindFlags &^= kindFlagNoBlockDecl

	n.rec().parent = n.ref(existingParent)

	existing.rec().next = existing.ref(n)
	n.rec().previous = n.ref(existing)

	if existingNext.IsNil() {
		existingParent.firstNode().rec().previous = existingParent.firstNode().ref(n)
		// this.next remains null
	} else {
		existingNext.rec().previous = existingNext.ref(n)
		n.rec().next = n.ref(existingNext)
	}
}

// Java: Node.insertBefore(Node)
func (n Node) InsertBefore(existing Node) {
	existing.CheckAttached()
	n.CheckDetached()

	existingParent := existing.parentNode()
	existingPrevious := existing.prevNode()
	existingParent.rec().kindFlags &^= kindFlagNoBlockDecl

	n.rec().parent = n.ref(existingParent)

	n.rec().next = n.ref(existing)
	existing.rec().previous = existing.ref(n)

	n.rec().previous = n.ref(existingPrevious)
	if existingPrevious.rec().next == 0 {
		existingParent.rec().first = existingParent.ref(n)
		// existingPrevious.next remains null
	} else {
		// existingParent.first remains existing
		existingPrevious.rec().next = existingPrevious.ref(n)
	}
}

// AddChildrenAfter adds all children after 'node'. If 'node' is null, add
// them to the front of this node.
//
// children is the first of a list of sibling nodes who have no parent.
// NOTE: Usually you would get this argument from a removeChildren() call. A
// single detached node will not work because its sibling pointers will not
// be correctly initialized.
//
// Java: Node.addChildrenAfter(Node, Node)
func (n Node) AddChildrenAfter(children Node, node Node) {
	if children.IsNil() {
		return // removeChildren() returns null when there are none
	}
	jutil.CheckArgument(node.IsNil() || node.parentNode() == n)
	// NOTE: If there is only one sibling, its previous pointer must point to
	// itself. Null indicates a fully detached node.
	CheckNotNullMsg1(children.prevNode(), "%s", children)
	if node.IsNil() {
		n.AddChildrenToFront(children)
		return
	}
	n.rec().kindFlags &^= kindFlagNoBlockDecl

	for child := children; !child.IsNil(); child = child.nextNode() {
		jutil.CheckArgument(child.rec().parent == 0)
		child.rec().parent = child.ref(n)
	}

	lastSibling := children.prevNode()
	nodeAfter := node.nextNode()
	lastSibling.rec().next = lastSibling.ref(nodeAfter)
	if nodeAfter.IsNil() {
		n.firstNode().rec().previous = n.firstNode().ref(lastSibling)
	} else {
		nodeAfter.rec().previous = nodeAfter.ref(lastSibling)
	}
	node.rec().next = node.ref(children)
	children.rec().previous = children.ref(node)
}

// ReplaceWith swaps `replacement` and its subtree into the position of
// `this`.
//
// Java: Node.replaceWith(Node)
func (n Node) ReplaceWith(replacement Node) {
	n.CheckAttached()
	replacement.CheckDetached()

	existingParent := n.parentNode()
	existingNext := n.nextNode()
	existingPrevious := n.prevNode()
	existingParent.rec().kindFlags &^= kindFlagNoBlockDecl

	// Copy over important information.
	// TODO(nickreid): Stop doing this. It's totally unexpected.
	replacement.SrcrefIfMissing(n)

	// The sequence below also has to work when `this` is an only child, ehich
	// can cause many of the variables to point to the same object.

	n.rec().parent = 0
	replacement.rec().parent = replacement.ref(existingParent)

	n.rec().previous = 0
	replacement.rec().previous = replacement.ref(existingPrevious)
	if existingPrevious.rec().next == 0 {
		existingParent.rec().first = existingParent.ref(replacement)
		// existingPrevious.next remains null
	} else {
		// existingParent.first is unchanged;
		existingPrevious.rec().next = existingPrevious.ref(replacement)
	}

	if existingNext.IsNil() {
		// this.next remains null;
		existingParent.firstNode().rec().previous = existingParent.firstNode().ref(replacement)
		// replacement.next remains null
	} else {
		n.rec().next = 0
		existingNext.rec().previous = existingNext.ref(replacement)
		replacement.rec().next = replacement.ref(existingNext)
	}
}

// Detach removes this node from its parent, but retains its subtree.
//
// Java: Node.detach()
func (n Node) Detach() Node {
	n.CheckAttached()

	existingParent := n.parentNode()
	existingNext := n.nextNode()
	existingPrevious := n.prevNode()
	existingParent.rec().kindFlags &^= kindFlagNoBlockDecl

	// The sequence below also has to work when `this` is an only child or has
	// a single sibling, which can cause many of the variables to point to the
	// same object.

	n.rec().parent = 0

	if existingNext.IsNil() {
		// this.next remains null;
		existingParent.firstNode().rec().previous = existingParent.firstNode().ref(existingPrevious)
	} else {
		n.rec().next = 0
		existingNext.rec().previous = existingNext.ref(existingPrevious)
	}

	n.rec().previous = 0
	if existingPrevious.rec().next == 0 {
		existingParent.rec().first = existingParent.ref(existingNext)
		// existingPrevious.next remains null
	} else {
		// existingParent.first is unchanged;
		existingPrevious.rec().next = existingPrevious.ref(existingNext)
	}

	return n
}

// Java: Node.CheckAttached()
func (n Node) CheckAttached() {
	jutil.CheckStateMsg1(n.rec().parent != 0, "Has no parent: %s", n)
}

// Java: Node.CheckDetached()
func (n Node) CheckDetached() {
	jutil.CheckStateMsg1(n.rec().parent == 0, "Has parent: %s", n)
	jutil.CheckStateMsg1(n.rec().next == 0, "Has next: %s", n)
	jutil.CheckStateMsg1(n.rec().previous == 0, "Has previous: %s", n)
}

// RemoveFirstChild removes the first child of Node. Equivalent to:
// node.removeChild(node.getFirstChild());
//
// Java: Node.removeFirstChild()
func (n Node) RemoveFirstChild() Node {
	child := n.firstNode()
	if !child.IsNil() {
		child.Detach()
	}
	return child
}

// RemoveChildren removes all children, but leaves them linked to each other.
// Returns the first child node.
//
// Java: Node.removeChildren()
func (n Node) RemoveChildren() Node {
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	children := n.firstNode()
	for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
		child.rec().parent = 0
	}
	n.rec().first = 0
	return children
}

// DetachChildren removes all children from this node and isolates the
// children from each other.
//
// Java: Node.detachChildren()
func (n Node) DetachChildren() {
	n.rec().kindFlags &^= kindFlagNoBlockDecl
	for child := n.firstNode(); !child.IsNil(); {
		nextChild := child.nextNode()
		child.rec().parent = 0
		child.rec().next = 0
		child.rec().previous = 0
		child = nextChild
	}
	n.rec().first = 0
}

func (n Node) lookupPropListHead(prop Node_Prop) *Node_PropListItem {
	propType := int8(prop)
	for x := n.coldPropListHead(); x != nil; x = x.next {
		if x.propType == propType {
			return x
		}
	}
	return nil
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

// Node_RebuildListWithoutProp returns the replacement list if the property
// was removed, or 'item' otherwise.
//
// Java: Node.rebuildListWithoutProp(PropListItem, Prop)
func Node_RebuildListWithoutProp(item *Node_PropListItem, prop Node_Prop) *Node_PropListItem {
	if item == nil {
		return nil
	} else if int(item.propType) == prop.Ordinal() {
		return item.next
	} else {
		result := Node_RebuildListWithoutProp(item.next, prop)
		if result == item.next {
			return item
		}
		return item.Chain(result)
	}
}

// Java: Node.getProp(Prop)
func (n Node) GetProp(propType Node_Prop) any {
	if propType == Node_Prop_SOURCE_FILE {
		if sfID := n.rec().sourceFile; sfID != 0 {
			return sourceFilePropItem(sfID).objectValue
		}
		return nil
	} else if propType == Node_Prop_JSDOC_INFO {
		if (n.props() & propFlagHasJSDoc) == 0 {
			return nil
		}
		if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
			if jsdoc := n.jsDocAux(); jsdoc != nil {
				return jsdoc
			}
			return nil
		}
		if box := n.extBoxOrNil(); box != nil && box.jsDocInfo != nil {
			return box.jsDocInfo
		}
	} else if propType == Node_Prop_TYPE_BEFORE_CAST {
		box := n.extBoxOrNil()
		if box == nil {
			return nil
		}
		if tbc := n.store().typeAt(box.typeBeforeCast); tbc != nil {
			return tbc.self()
		}
	} else if mask := nodePropBitMask[propType]; mask != 0 {
		if (n.props() & mask) == 0 {
			return nil
		}
		if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
			if x := n.lookupPropListHead(propType); x != nil {
				return x.GetObjectValue()
			}
		}
		panic(jutil.NewUnsupportedOperationException())
	} else if propType == Node_Prop_CONSTANT_VAR_FLAGS {
		if (n.props() & nodePropConstVarMask) != 0 {
			panic(jutil.NewUnsupportedOperationException())
		}
	} else if propType == Node_Prop_SIDE_EFFECT_FLAGS {
		if (n.props() & nodePropSideEffectMask) != 0 {
			panic(jutil.NewUnsupportedOperationException())
		}
	} else if propType == Node_Prop_CHANGE_TIME && n.nodeKind() == node_kindNode && uint32(n.rec().payload) != 0 {
		panic(jutil.NewUnsupportedOperationException())
	}
	if x := n.lookupPropListHead(propType); x != nil {
		return x.GetObjectValue()
	}
	return nil
}

// Java: Node.getBooleanProp(Prop)
func (n Node) GetBooleanProp(propType Node_Prop) bool {
	if mask := nodePropBitMask[propType]; mask != 0 {
		return (n.props() & mask) != 0
	}
	return n.GetIntProp(propType) != 0
}

// GetIntProp returns the integer value for the property, or 0 if the
// property is not defined.
//
// Java: Node.getIntProp(Prop)
func (n Node) GetIntProp(propType Node_Prop) int32 {
	if mask := nodePropBitMask[propType]; mask != 0 {
		if (n.props() & mask) == 0 {
			return 0
		}
		if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
			if x := n.lookupPropListHead(propType); x != nil {
				return x.GetIntValue()
			}
		}
		return 1
	}
	if propType == Node_Prop_CHANGE_TIME && n.nodeKind() == node_kindNode {
		if v := int32(uint32(n.rec().payload)); v != 0 {
			return v
		}
		if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
			return 0
		}
	} else if propType == Node_Prop_CONSTANT_VAR_FLAGS {
		if v := int32((n.props() >> nodePropConstVarShift) & 3); v != 0 {
			return v
		}
		if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
			return 0
		}
	} else if propType == Node_Prop_SIDE_EFFECT_FLAGS {
		if v := int32((n.props() >> nodePropSideEffectShift) & 15); v != 0 {
			return v
		}
		if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
			return 0
		}
	} else if propType == Node_Prop_SOURCE_FILE {
		if n.rec().sourceFile != 0 {
			panic(jutil.NewUnsupportedOperationException())
		}
		return 0
	} else if propType == Node_Prop_JSDOC_INFO {
		if (n.props() & propFlagHasJSDoc) == 0 {
			return 0
		}
		if box := n.extBoxOrNil(); (n.rec().kindFlags&kindFlagHasExtBox) == 0 || (box != nil && box.jsDocInfo != nil) {
			panic(jutil.NewUnsupportedOperationException())
		}
	} else if propType == Node_Prop_TYPE_BEFORE_CAST {
		box := n.extBoxOrNil()
		if box == nil {
			return 0
		}
		if box.typeBeforeCast != 0 {
			panic(jutil.NewUnsupportedOperationException())
		}
	}
	if x := n.lookupPropListHead(propType); x != nil {
		return x.GetIntValue()
	}
	return 0
}

// PutProp sets an object-valued property; a null value removes it.
//
// Go-only: language-adapter — jutil.IsNil treats a typed nil pointer passed
// through `any` as Java null.
//
// Java: Node.putProp(Prop, Object)
func (n Node) PutProp(prop Node_Prop, value any) {
	if jutil.IsNil(value) {
		n.RemoveProp(prop)
		return
	}
	if prop == Node_Prop_SOURCE_FILE {
		n.rec().sourceFile = allocSourceFileID(value)
		return
	} else if prop == Node_Prop_JSDOC_INFO {
		if info, ok := value.(*JSDocInfo); ok {
			n.SetJSDocInfo(info)
			return
		}
		box := n.ensureExtBox()
		box.jsDocInfo = nil
		n.orProps(propFlagHasJSDoc)
	} else if prop == Node_Prop_TYPE_BEFORE_CAST {
		if t, ok := value.(JSType); ok {
			n.SetJSTypeBeforeCast(t)
			return
		}
		box := n.ensureExtBox()
		box.typeBeforeCast = 0
	} else if prop == Node_Prop_CHANGE_TIME && n.nodeKind() == node_kindNode {
		n.rec().payload = 0
		n.rec().kindFlags |= kindFlagHasExtra
	} else if mask := nodePropBitMask[prop]; mask != 0 {
		n.orProps(mask)
		n.rec().kindFlags |= kindFlagHasExtra
	}
	head := n.coldPropListHead()
	for cur := head; cur != nil && !cur.shared; cur = cur.next {
		if int(cur.propType) == prop.Ordinal() && cur.isObject {
			cur.objectValue = value
			return
		}
	}
	head = Node_RebuildListWithoutProp(head, prop)
	n.setColdPropListHead(NewNode_ObjectPropListItem(int8(prop.Ordinal()), value, head))
}

// Java: Node.putBooleanProp(Prop, boolean)
func (n Node) PutBooleanProp(propType Node_Prop, value bool) {
	if mask := nodePropBitMask[propType]; mask != 0 {
		if value {
			n.orProps(mask)
			if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
				if head := n.coldPropListHead(); head != nil {
					n.setColdPropListHead(Node_RebuildListWithoutProp(head, propType))
				}
			}
		} else {
			if (n.props() & mask) == 0 {
				return
			}
			n.andNotProps(mask)
			if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
				if head := n.coldPropListHead(); head != nil {
					n.setColdPropListHead(Node_RebuildListWithoutProp(head, propType))
				}
			}
		}
		return
	}
	if value {
		n.PutIntProp(propType, 1)
	} else {
		n.RemoveProp(propType)
	}
}

// Java: Node.putIntProp(Prop, int)
func (n Node) PutIntProp(prop Node_Prop, value int32) {
	if mask := nodePropBitMask[prop]; mask != 0 {
		if value == 0 {
			if (n.props() & mask) == 0 {
				return
			}
			n.andNotProps(mask)
			if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
				if head := n.coldPropListHead(); head != nil {
					n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
				}
			}
			return
		}
		if value == 1 {
			n.orProps(mask)
			if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
				if head := n.coldPropListHead(); head != nil {
					n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
				}
			}
			return
		}
		n.orProps(mask)
		n.rec().kindFlags |= kindFlagHasExtra
	} else if prop == Node_Prop_CHANGE_TIME && n.nodeKind() == node_kindNode {
		n.rec().payload = uint32(value)
		if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
			if head := n.coldPropListHead(); head != nil {
				n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
			}
		}
		return
	} else if prop == Node_Prop_CONSTANT_VAR_FLAGS {
		if value == 0 {
			n.andNotProps(nodePropConstVarMask)
			if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
				if head := n.coldPropListHead(); head != nil {
					n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
				}
			}
			return
		}
		if value&^3 == 0 {
			n.setProps((n.props() &^ nodePropConstVarMask) | (uint64(value) << nodePropConstVarShift))
			if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
				if head := n.coldPropListHead(); head != nil {
					n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
				}
			}
			return
		}
		n.andNotProps(nodePropConstVarMask)
		n.rec().kindFlags |= kindFlagHasExtra
	} else if prop == Node_Prop_SIDE_EFFECT_FLAGS {
		if value == 0 {
			n.andNotProps(nodePropSideEffectMask)
			if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
				if head := n.coldPropListHead(); head != nil {
					n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
				}
			}
			return
		}
		if value&^15 == 0 {
			n.setProps((n.props() &^ nodePropSideEffectMask) | (uint64(value) << nodePropSideEffectShift))
			if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
				if head := n.coldPropListHead(); head != nil {
					n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
				}
			}
			return
		}
		n.andNotProps(nodePropSideEffectMask)
		n.rec().kindFlags |= kindFlagHasExtra
	} else if prop == Node_Prop_SOURCE_FILE {
		// Go-only: the source file is an object property (an id in the record).
		if value != 0 {
			panic(jutil.NewUnsupportedOperationException())
		}
		n.rec().sourceFile = 0
		return
	} else if prop == Node_Prop_JSDOC_INFO {
		if value == 0 {
			n.RemoveProp(prop)
			return
		}
		box := n.ensureExtBox()
		box.jsDocInfo = nil
		n.orProps(propFlagHasJSDoc)
	} else if prop == Node_Prop_TYPE_BEFORE_CAST {
		if value == 0 {
			n.RemoveProp(prop)
			return
		}
		box := n.ensureExtBox()
		box.typeBeforeCast = 0
	}
	head := n.coldPropListHead()
	if value != 0 && head != nil && !head.shared && int(head.propType) == prop.Ordinal() && !head.isObject {
		head.intValue = value
		return
	}
	head = Node_RebuildListWithoutProp(head, prop)
	if value != 0 {
		head = NewNode_IntPropListItem(int8(prop.Ordinal()), value, head)
	}
	n.setColdPropListHead(head)
}

// Java: Node.removeProp(Prop)
func (n Node) RemoveProp(prop Node_Prop) {
	if mask := nodePropBitMask[prop]; mask != 0 {
		if (n.props() & mask) == 0 {
			return
		}
		n.andNotProps(mask)
		if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
			if head := n.coldPropListHead(); head != nil {
				n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
			}
		}
		return
	}
	switch prop {
	case Node_Prop_SOURCE_FILE:
		n.rec().sourceFile = 0
		return
	case Node_Prop_CHANGE_TIME:
		if n.nodeKind() == node_kindNode {
			n.rec().payload = 0
			if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
				return
			}
		}
	case Node_Prop_CONSTANT_VAR_FLAGS:
		n.andNotProps(nodePropConstVarMask)
		if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
			return
		}
	case Node_Prop_SIDE_EFFECT_FLAGS:
		n.andNotProps(nodePropSideEffectMask)
		if (n.rec().kindFlags & kindFlagHasExtra) == 0 {
			return
		}
	case Node_Prop_JSDOC_INFO:
		if (n.props() & propFlagHasJSDoc) == 0 {
			return
		}
		n.andNotProps(propFlagHasJSDoc)
		if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
			n.setAux(nodeAux{})
			return
		}
		box := n.extBox()
		if box.jsDocInfo != nil {
			box.jsDocInfo = nil
			n.maybeCollapseExtBox(box)
			return
		}
	case Node_Prop_TYPE_BEFORE_CAST:
		if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
			return
		}
		box := n.extBox()
		if box.typeBeforeCast != 0 {
			box.typeBeforeCast = 0
			n.maybeCollapseExtBox(box)
			return
		}
	}
	if head := n.coldPropListHead(); head != nil {
		n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
	}
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

// Java: Node.getDouble()
func (n Node) GetDouble() float64 {
	n.node_castKind(node_kindNumber)
	return n.node_number()
}

// Java: Node.setDouble(double)
func (n Node) SetDouble(x float64) {
	jutil.CheckStateMsg1(!math.IsNaN(x), "%s", x)
	jutil.CheckStateMsg1(base.JSCompDoubles_IsPositive(x), "%s", x)
	n.node_castKind(node_kindNumber)
	n.rec().payload = numberPayload(x)
}

var (
	rhinoBigIntMu   sync.RWMutex
	rhinoBigIntPool []*big.Int
)

func rhinoInternBigInt(v *big.Int) uint32 {
	rhinoBigIntMu.Lock()
	rhinoBigIntPool = append(rhinoBigIntPool, v)
	id := uint32(len(rhinoBigIntPool))
	rhinoBigIntMu.Unlock()
	return id
}

func rhinoBigIntFromID(id uint32) *big.Int {
	if id == 0 {
		return nil
	}
	rhinoBigIntMu.RLock()
	v := rhinoBigIntPool[id-1]
	rhinoBigIntMu.RUnlock()
	return v
}

// Java: Node.getBigInt()
func (n Node) GetBigInt() *big.Int {
	n.node_castKind(node_kindBigInt)
	return rhinoBigIntFromID(uint32(n.rec().payload))
}

// Java: Node.setBigInt(BigInteger)
func (n Node) SetBigInt(number *big.Int) {
	jutil.CheckNotNull(number)
	jutil.CheckStateMsg1(number.Sign() >= 0, "%s", number)
	n.node_castKind(node_kindBigInt)
	n.rec().payload = rhinoInternBigInt(number)
}

// Java: Node.getString()
func (n Node) GetString() string {
	n.node_castKind(node_kindString)
	return n.nodeStringValue()
}

//go:noinline
func (n Node) failGetStringID() {
	if got := n.nodeKind(); got != node_kindString {
		failNodeCastKind(got, node_kindString)
	}
	panic(jutil.NewNullPointerException())
}

//go:noinline
func (n Node) failGetRawStringID() {
	if got := n.nodeKind(); got != node_kindTemplateLiteralSubstring {
		failNodeCastKind(got, node_kindTemplateLiteralSubstring)
	}
	panic(jutil.NewNullPointerException())
}

// GetStringID returns the 32-bit interned string ID of a StringNode.
func (n Node) GetStringID() uint32 {
	id := uint32(n.rec().payload)
	if n.rec().kindFlags&kindMask != uint8(node_kindString) || id == 0 {
		n.failGetStringID()
	}
	return id
}

// Java: Node.setString(String)
func (n Node) SetString(str string) {
	n.node_castKind(node_kindString)
	n.setStringPayload(rhinoStringPoolIntern(str))
}

// Java: Node.setStringFromStringPool(RhinoStringPool.LazyInternedStringList, int)
func (n Node) SetStringFromStringPool(stringPool *RhinoStringPool_LazyInternedStringList, offset int) {
	n.node_castKind(node_kindString)
	n.setStringPayload(stringPool.GetID(offset))
}

// Java: Node.getRawString()
func (n Node) GetRawString() string {
	n.node_castKind(node_kindTemplateLiteralSubstring)
	return n.nodeStringValue() // raw
}

// GetRawStringID returns the 32-bit interned raw string ID of a TemplateLiteralSubstringNode.
func (n Node) GetRawStringID() uint32 {
	id := uint32(n.rec().payload)
	if n.rec().kindFlags&kindMask != uint8(node_kindTemplateLiteralSubstring) || id == 0 {
		n.failGetRawStringID()
	}
	return id
}

// Java: Node.getCookedString()
func (n Node) GetCookedString() *string {
	n.node_castKind(node_kindTemplateLiteralSubstring)
	return n.node_cooked()
}

// GetCookedStringID returns the 32-bit interned cooked string ID of a TemplateLiteralSubstringNode, or 0 if unset.
func (n Node) GetCookedStringID() uint32 {
	n.node_castKind(node_kindTemplateLiteralSubstring)
	return uint32(n.node_cookedID())
}

// Java: Node.toString()
func (n Node) String() string {
	return n.ToStringWithFlags(true, true, true)
}

// Java: Node.toString(boolean, boolean, boolean)
func (n Node) ToStringWithFlags(printSource bool, printAnnotations bool, printType bool) string {
	var sb strings.Builder
	n.ToStringInto(&sb, printSource, printAnnotations, printType)
	return sb.String()
}

// Java: Node.toString(StringBuilder, boolean, boolean, boolean)
func (n Node) ToStringInto(sb *strings.Builder, printSource bool, printAnnotations bool, printType bool) {
	sb.WriteString(n.rec().token.String())
	if n.nodeKind() == node_kindString {
		sb.WriteByte(' ')
		sb.WriteString(n.GetString())
	} else if n.rec().token == Token_FUNCTION {
		sb.WriteByte(' ')
		// In the case of JsDoc trees, the first child is often not a string
		// which causes exceptions to be thrown when calling toString or
		// toStringTree.
		if n.rec().first == 0 || n.firstNode().rec().token != Token_NAME {
			sb.WriteString("<invalid>")
		} else {
			sb.WriteString(n.firstNode().GetString())
		}
	} else if n.rec().token == Token_NUMBER {
		sb.WriteByte(' ')
		sb.WriteString(jutil.DoubleToString(n.GetDouble()))
	}
	if printSource {
		lineno := n.GetLineno()
		if lineno != -1 {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(lineno))
			sb.WriteByte(':')
			sb.WriteString(strconv.Itoa(n.GetCharno()))
			sb.WriteByte(' ')
		}
		if length := n.GetLength(); length != 0 {
			sb.WriteString(" [length: ")
			sb.WriteString(strconv.Itoa(length))
			sb.WriteByte(']')
		}
	}

	if printAnnotations {
		keys := n.GetSortedPropTypes()
		for i := 0; i < len(keys); i++ {
			type_ := Node_propValues[keys[i]]
			if type_ == Node_Prop_CLOSURE_UNAWARE_SHADOW {
				sb.WriteString(" [is_shadow_host]")
				continue
			}
			x := n.LookupProperty(type_)
			sb.WriteString(" [")
			sb.WriteString(jutil.StringToLowerCaseASCII(type_.String()))
			sb.WriteString(": ")
			sb.WriteString(jutil.ToJavaString(x))
			sb.WriteByte(']')
		}
		if origID := n.node_originalNameID(); origID != 0 {
			sb.WriteString(" [original_name: ")
			sb.WriteString(rhinoStringFromID(origID))
			sb.WriteByte(']')
		}
	}

	if printType && n.tc() != 0 {
		typeString := jutil.ToJavaString(n.typeOrColorValue())
		sb.WriteString(" : ")
		sb.WriteString(typeString)
	}
}

// Java: Node.createJsonPair(String, String)
func Node_CreateJsonPair(name string, value string) string {
	return "\"" + name + "\":\"" + value + "\""
}

// Java: Node.createJsonPairRawValue(String, String)
func Node_CreateJsonPairRawValue(name string, value string) string {
	return "\"" + name + "\":" + value
}

// node_appender is Java's Appendable over an io.StringWriter: write errors
// become the IOException Appendable.append throws.
type node_appender struct{ w io.StringWriter }

func (a node_appender) append(s string) {
	if _, err := a.w.WriteString(s); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

// ToJson writes the JSON form of this node and its subtree.
//
// Java: Node.toJson(Appendable)
func (n Node) ToJson(appendable io.StringWriter) {
	sb := node_appender{appendable}
	sb.append("{")
	sb.append(Node_CreateJsonPair("token", n.rec().token.String()))
	if n.nodeKind() == node_kindString {
		sb.append(",")
		sb.append(Node_CreateJsonPair("string", n.GetString()))
	} else if n.rec().token == Token_FUNCTION {
		sb.append(",")
		// In the case of JsDoc trees, the first child is often not a string
		// which causes exceptions to be thrown when calling toString or
		// toStringTree.
		if n.rec().first == 0 || n.firstNode().rec().token != Token_NAME {
			sb.append(Node_CreateJsonPair("functionName", "<invalid>"))
		} else {
			sb.append(Node_CreateJsonPair("functionName", n.firstNode().GetString()))
		}
	} else if n.rec().token == Token_NUMBER {
		sb.append(",")
		sb.append(Node_CreateJsonPair("number", jutil.DoubleToString(n.GetDouble())))
	}
	lineno := n.GetLineno()
	if lineno != -1 {
		sb.append(",")
		sb.append("\"sourceLocation\":{")
		sb.append(Node_CreateJsonPairRawValue("line", strconv.Itoa(lineno)))
		sb.append(",")
		sb.append(Node_CreateJsonPairRawValue("col", strconv.Itoa(n.GetCharno())))
		if length := n.GetLength(); length != 0 {
			sb.append(",")
			sb.append(Node_CreateJsonPairRawValue("length", strconv.Itoa(length)))
		}
		sb.append("}")
	}

	if origID := n.node_originalNameID(); origID != 0 {
		sb.append(",")
		sb.append(Node_CreateJsonPair("original_name", rhinoStringFromID(origID)))
	}

	keys := n.GetSortedPropTypes()
	if len(keys) != 0 {
		sb.append(",")
		sb.append("\"props\":{")
		for i := 0; i < len(keys); i++ {
			type_ := Node_propValues[keys[i]]
			x := n.LookupProperty(type_)
			sb.append(
				Node_CreateJsonPair(
					jutil.StringToLowerCaseASCII(type_.String()),
					jutil.StringReplace(jutil.StringReplace(x.String(), "\n", "\\n"), "\"", "\\\"")))
			if i+1 < len(keys) {
				sb.append(",")
			}
		}
		sb.append("}")
	}

	if n.tc() != 0 {
		typeString := jutil.ToJavaString(n.typeOrColorValue())
		sb.append(",")
		sb.append(Node_CreateJsonPair("typeString", jutil.StringReplace(typeString, "\"", "\\\"")))
	}

	if n.rec().first != 0 {
		sb.append(",")
		sb.append("\"children\":[")
		for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
			child.ToJson(appendable)
			if child.rec().next != 0 {
				sb.append(",")
			}
		}
		sb.append("]")
	}
	sb.append("}")
}

// Java: Node.toStringTree()
func (n Node) ToStringTree() string {
	return n.ToStringTreeImpl()
}

// Java: Node.toStringTreeImpl()
func (n Node) ToStringTreeImpl() string {
	var result string
	jutil.Try(func() {
		var s strings.Builder
		n.AppendStringTree(&s)
		result = s.String()
	}, func(e jutil.IIOException) {
		panic(jutil.NewIllegalStateExceptionMsg("Should not happen\n" + e.Error()))
	})
	return result
}

// Java: Node.appendStringTree(Appendable)
func (n Node) AppendStringTree(appendable io.StringWriter) {
	Node_ToStringTreeHelper(n, 0, appendable)
}

// Java: Node.toStringTreeHelper(Node, int, Appendable)
func Node_ToStringTreeHelper(n Node, level int, appendable io.StringWriter) {
	sb := node_appender{appendable}
	for i := 0; i != level; i++ {
		sb.append("    ")
	}
	sb.append(n.String())
	sb.append("\n")
	for cursor := n.firstNode(); !cursor.IsNil(); cursor = cursor.nextNode() {
		Node_ToStringTreeHelper(cursor, level+1, appendable)
	}
	shadow := n.GetClosureUnawareShadow()
	if !shadow.IsNil() {
		Node_ToStringTreeHelper(shadow, level+1, appendable)
	}
}

// Java: Node.appendJsonTree(Appendable)
func (n Node) AppendJsonTree(appendable io.StringWriter) {
	Node_ToJsonTreeHelper(n, appendable)
}

// Java: Node.toJsonTreeHelper(Node, Appendable)
func Node_ToJsonTreeHelper(n Node, appendable io.StringWriter) {
	n.ToJson(appendable)
}

// nodeData is the record of an AST node (fields as declared in Java, plus
// the folded subclass payloads); a Node is a handle to one (node_store.go).
//
// Go-only: performance. 36 bytes, no Go pointers, 4-byte aligned: the
// tree links are node ids; the leaf payload is one 32-bit value whose
// meaning depends on the kind (a string id, a number, a change time...;
// node_payload.go); the source file is a 32-bit id (allocSourceFileID); the
// 8 common boolean properties, CONSTANT_VAR_FLAGS and SIDE_EFFECT_FLAGS
// share a 32-bit property word, and the 32 rare boolean properties are cold
// data (propWord).
// What most nodes do not have is kept beside the record (node_tables.go):
// the JSType or Color handle and the original name in per-block columns,
// the cold data (JSDocInfo, remaining properties, overflow length, a cooked
// template string that differs from the raw one: about 1% of nodes) in a
// map of the store, looked up only when the record's flags say the node has
// some.
type nodeData struct {
	next     nodeID // next sibling, a linked list (0 = none)
	first    nodeID // first element of a linked list of children
	parent   nodeID // parent node
	previous nodeID // previous sibling, a circular linked list

	payload      uint32   // one value whose meaning depends on the kind (node_payload.go)
	sourceFile   uint32   // the SOURCE_FILE property: an id from allocSourceFileID (0 = none)
	propFlags    propWord // the property bits (Node.props): bool props, JSDoc, ConstVar, SideEffect
	linenoCharno int32    // line (upper 20 bits) and column (low 12 bits)
	length16     uint16   // code length (0..65534 inline; 65535: in the nodeExtBox)
	token        Token    // token kind (uint8)
	kindFlags    uint8    // bits 0..2: node_kind; bit 3: typeIsColor; bit 4: noBlockDecls; bit 5: hasExtra; bit 6: origNameIsString; bit 7: hasExtBox
}

func (n Node) HasNoBlockScopeDecls() bool {
	return n.rec().kindFlags&kindFlagNoBlockDecl != 0
}

func (n Node) SetNoBlockScopeDecls() {
	n.rec().kindFlags |= kindFlagNoBlockDecl
}

func (n Node) ClearNoBlockScopeDecls() {
	n.rec().kindFlags &^= kindFlagNoBlockDecl
}

// node_number is the NumberNode.number field.
func (n Node) node_number() float64 {
	return payloadNumber(n.rec().payload)
}

// node_cooked is the TemplateLiteralSubstringNode.cooked field.
func (n Node) node_cooked() *string {
	return rhinoStringPtrFromID(n.node_cookedID())
}

// ==========================================================================
// Source position management

// Java: Node.setStaticSourceFileFrom(Node)
func (n Node) SetStaticSourceFileFrom(other Node) {
	n.rec().sourceFile = other.rec().sourceFile
}

// Java: Node.setStaticSourceFile(StaticSourceFile)
func (n Node) SetStaticSourceFile(file StaticSourceFile) Node {
	n.PutProp(Node_Prop_SOURCE_FILE, file)
	return n
}

// SetSourceFileForTesting sets the source file to a non-extern file of the
// given name.
//
// Java: Node.setSourceFileForTesting(String)
func (n Node) SetSourceFileForTesting(name string) {
	n.PutProp(Node_Prop_SOURCE_FILE, NewSimpleSourceFile(name, StaticSourceFile_SourceKind_STRONG))
}

// TODO(johnlenz): make this final
//
// Java: Node.getSourceFileName()
func (n Node) GetSourceFileName() *string {
	sfID := n.rec().sourceFile
	if sfID == 0 {
		return nil
	}
	idx := sfID - 1
	if p := sourceFileNamePtrChunks[idx>>10].Load()[idx&1023]; p != nil {
		return p
	}
	_ = node_cast[StaticSourceFile](sourceFilePropItem(sfID).objectValue)
	return nil
}

// GetStaticSourceFile returns the source file associated with this input.
//
// Java: Node.getStaticSourceFile()
func (n Node) GetStaticSourceFile() StaticSourceFile {
	sfID := n.rec().sourceFile
	if sfID == 0 {
		return nil
	}
	return node_cast[StaticSourceFile](sourceFilePropItem(sfID).objectValue)
}

// SetInputId sets the ID of the input this Node came from.
//
// Java: Node.setInputId(InputId)
func (n Node) SetInputId(inputId *InputId) {
	n.PutProp(Node_Prop_INPUT_ID, inputId)
}

// GetInputId returns the Id of the CompilerInput associated with this Node.
//
// Java: Node.getInputId()
func (n Node) GetInputId() *InputId {
	return node_cast[*InputId](n.GetProp(Node_Prop_INPUT_ID))
}

// GetOriginalName returns the original name of this node, if the node has
// been renamed.
//
// Deprecated: "original name" is poorly defined.
//
// Java: Node.getOriginalName()
func (n Node) GetOriginalName() *string {
	return rhinoStringPtrFromID(n.node_originalNameID())
}

// GetOriginalNameID returns the 32-bit interned originalName ID of n, or 0 if unset.
func (n Node) GetOriginalNameID() uint32 {
	return uint32(n.node_originalNameID())
}

func (n Node) getOriginalNameString() (string, bool) {
	id := n.node_originalNameID()
	if id == 0 {
		return "", false
	}
	return rhinoStringFromID(id), true
}

// Java: Node.setOriginalName(String)
func (n Node) SetOriginalName(s *string) {
	if s == nil {
		n.node_setOriginalNameID(0)
	} else {
		n.node_setOriginalNameID(rhinoStringPoolIntern(*s))
	}
}

// SetOriginalNameString interns s in RhinoStringPool and sets it as the node's
// original name without boxing s on the heap.
func (n Node) SetOriginalNameString(s string) {
	n.node_setOriginalNameID(rhinoStringPoolIntern(s))
}

// Java: Node.setOriginalNameFromStringPool(RhinoStringPool.LazyInternedStringList, int)
func (n Node) SetOriginalNameFromStringPool(stringPool *RhinoStringPool_LazyInternedStringList, offset int) {
	n.node_setOriginalNameID(stringPool.GetID(offset))
}

// SetOriginalNameFromName copies the original name.
//
// Java: Node.setOriginalNameFromName(Node)
func (n Node) SetOriginalNameFromName(name Node) {
	name.node_castKind(node_kindString)
	n.node_setOriginalNameID(rhinoStringID(name.rec().payload))
}

// IsIndexable reports whether this node should be indexed by static
// analysis / code indexing tools.
//
// Java: Node.isIndexable()
func (n Node) IsIndexable() bool {
	return !n.GetBooleanProp(Node_Prop_NON_INDEXABLE)
}

// Java: Node.makeNonIndexable()
func (n Node) MakeNonIndexable() {
	n.PutBooleanProp(Node_Prop_NON_INDEXABLE, true)
}

// Java: Node.makeNonIndexableRecursive()
func (n Node) MakeNonIndexableRecursive() {
	n.MakeNonIndexable()
	for child := n.firstNode(); !child.IsNil(); child = child.GetNext() {
		child.MakeNonIndexableRecursive()
	}
}

// Java: Node.isFromExterns()
func (n Node) IsFromExterns() bool {
	file := n.GetStaticSourceFile()
	if file == nil {
		return false
	}
	return file.IsExtern()
}

// IsClosureUnawareCode indicates that this node is for source that was
// written without any understanding of the compiler's various type-checking
// or optimization limitations, and should be optimized very cautiously by
// the compiler.
//
// This attribute is per source file, not per node / subset of an source
// file's AST.
//
// Java: Node.isClosureUnawareCode()
func (n Node) IsClosureUnawareCode() bool {
	file := n.GetStaticSourceFile()
	return file != nil && file.IsClosureUnawareCode()
}

// Java: Node.getLength()
func (n Node) GetLength() int {
	if n.rec().length16 != 0xFFFF {
		return int(n.rec().length16)
	}
	if box := n.extBoxOrNil(); box != nil {
		return int(box.length)
	}
	return 0
}

// Java: Node.setLength(int)
func (n Node) SetLength(length int) {
	if length >= 0 && length < 0xFFFF {
		wasInBox := n.rec().length16 == 0xFFFF
		n.rec().length16 = uint16(length)
		if wasInBox {
			n.maybeCollapseExtBox(n.extBox())
		}
		return
	}
	n.ensureExtBox().length = int32(length)
	n.rec().length16 = 0xFFFF
}

func (n Node) copyLengthFrom(other Node) {
	if other.rec().length16 != 0xFFFF && n.rec().length16 != 0xFFFF {
		n.rec().length16 = other.rec().length16
		return
	}
	n.SetLength(other.GetLength())
}

// Java: Node.getLineno()
func (n Node) GetLineno() int {
	if n.rec().linenoCharno == -1 {
		return -1
	}
	return int(jutil.Ushr32(n.rec().linenoCharno, Node_CHARNO_BITS))
}

// GetCharno returns the 0-based column number.
//
// Java: Node.getCharno()
func (n Node) GetCharno() int {
	if n.rec().linenoCharno == -1 {
		return -1
	}
	return int(n.rec().linenoCharno & Node_MAX_COLUMN_NUMBER)
}

// Java: Node.getLocation()
func (n Node) GetLocation() string {
	return jutil.ToJavaString(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

// TODO(johnlenz): make this final
//
// Java: Node.getSourceOffset()
func (n Node) GetSourceOffset() int {
	file := n.GetStaticSourceFile()
	if file == nil {
		return -1
	}
	lineno := n.GetLineno()
	if lineno == -1 {
		return -1
	}
	return int(int32(file.GetLineOffset(lineno)) + int32(n.GetCharno()))
}

// Java: Node.getSourcePosition()
func (n Node) GetSourcePosition() int {
	return int(n.rec().linenoCharno)
}

// Node_CHARNO_BITS represents how many of the lower-order bits of
// linenoCharno are reserved for storing the column number. Bits above these
// store the line number. This gives us decent position information for
// everything except files already passed through a minimizer, where lines
// might be longer than 4096 characters.
const Node_CHARNO_BITS = 12

// Node_MAX_COLUMN_NUMBER represents the maximum column number that can be
// represented. JSCompiler's modifications to Rhino cause all tokens located
// beyond the maximum column to MAX_COLUMN_NUMBER.
const Node_MAX_COLUMN_NUMBER = (1 << Node_CHARNO_BITS) - 1

// SetLinenoCharno merges the line number and character number in one
// integer.
//
// The charno takes the first 12 bits and the line number takes the rest. If
// the charno is greater than (2^12)-1 it is adjusted to (2^12)-1
//
// Java: Node.setLinenoCharno(int, int)
func (n Node) SetLinenoCharno(lineno int, charno int) Node {
	if lineno < 0 || charno < 0 {
		n.rec().linenoCharno = -1
		return n
	}

	if charno > Node_MAX_COLUMN_NUMBER {
		charno = Node_MAX_COLUMN_NUMBER
	}
	n.rec().linenoCharno = jutil.Shl32(int32(lineno), Node_CHARNO_BITS) | int32(charno)

	return n
}

// ==========================================================================
// Iteration

// Children returns an iterable object that iterates over this node's
// children. The iterator does not support the optional operation
// Iterator#remove().
//
// Deprecated: see above.
//
// Java: Node.children()
func (n Node) Children() iter.Seq[Node] {
	if n.rec().first == 0 {
		return func(yield func(Node) bool) {}
	}
	start := n.firstNode()
	return func(yield func(Node) bool) {
		current := start
		for !current.IsNil() {
			c := current
			current = current.GetNext()
			if !yield(c) {
				return
			}
		}
	}
}

// ==========================================================================
// Accessors

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

// Java: Node.getParent()
func (n Node) GetParent() Node {
	return n.parentNode()
}

// Java: Node.hasParent()
func (n Node) HasParent() bool {
	return n.rec().parent != 0
}

// Java: Node.getGrandparent()
func (n Node) GetGrandparent() Node {
	if n.rec().parent == 0 {
		return NilNode
	}
	return n.parentNode().parentNode()
}

// GetAncestor gets the ancestor node relative to this.
//
// level: 0 = this, 1 = the parent, etc.
//
// Java: Node.getAncestor(int)
func (n Node) GetAncestor(level int) Node {
	jutil.CheckArgument(level >= 0)
	node := n
	for !node.IsNil() {
		cont := level > 0
		level--
		if !cont {
			break
		}
		node = node.GetParent()
	}
	return node
}

// IsDescendantOf reports whether this Node is the same as `node` or a
// descendant of `node`.
//
// Java: Node.isDescendantOf(Node)
func (n Node) IsDescendantOf(node Node) bool {
	for c := n; !c.IsNil(); c = c.parentNode() {
		if c == node {
			return true
		}
	}
	return false
}

// Java: Node.isOnlyChildOf(Node)
func (n Node) IsOnlyChildOf(possibleParent Node) bool {
	return possibleParent == n.GetParent() && n.GetPrevious().IsNil() && n.GetNext().IsNil()
}

// Java: Node.isFirstChildOf(Node)
func (n Node) IsFirstChildOf(possibleParent Node) bool {
	return possibleParent == n.GetParent() && n.GetPrevious().IsNil()
}

// Java: Node.isSecondChildOf(Node)
func (n Node) IsSecondChildOf(possibleParent Node) bool {
	previousNode := n.GetPrevious()

	return !previousNode.IsNil() && previousNode.IsFirstChildOf(possibleParent)
}

// GetAncestors iterates all of the node's ancestors excluding itself.
//
// Java: Node.getAncestors()
func (n Node) GetAncestors() *Node_AncestorIterable {
	return NewNode_AncestorIterable(n.GetParent())
}

// Node_AncestorIterable is an iterator to go up the ancestor tree.
type Node_AncestorIterable struct {
	cur Node
}

// NewNode_AncestorIterable starts at `cur`.
//
// Java: Node.AncestorIterable(Node)
func NewNode_AncestorIterable(cur Node) *Node_AncestorIterable {
	return &Node_AncestorIterable{cur: cur}
}

// All is the iterable's iterator(). Like the Java iterator, it advances the
// iterable's own cursor, so an AncestorIterable can be traversed only once.
//
// Java: Node.AncestorIterable.iterator()
func (a *Node_AncestorIterable) All() iter.Seq[Node] {
	return func(yield func(Node) bool) {
		for !a.cur.IsNil() {
			c := a.cur
			a.cur = a.cur.GetParent()
			if !yield(c) {
				return
			}
		}
	}
}

// HasOneChild checks for one child more efficiently than by iterating over
// all the children as is done with Node.getChildCount().
//
// Java: Node.hasOneChild()
func (n Node) HasOneChild() bool {
	return n.rec().first != 0 && n.firstNode().rec().next == 0
}

// HasTwoChildren checks for two children more efficiently than
// `getChildCount() == 2`.
//
// Java: Node.hasTwoChildren()
func (n Node) HasTwoChildren() bool {
	return n.rec().first != 0 && n.firstNode().rec().next != 0 && n.firstNode().nextNode() == n.GetLastChild()
}

// HasZeroOrOneChild checks for zero or one child more efficiently than by
// iterating over all the children as is done with Node.getChildCount().
//
// Java: Node.hasZeroOrOneChild()
func (n Node) HasZeroOrOneChild() bool {
	return n.firstNode() == n.GetLastChild()
}

// HasMoreThanOneChild checks for more than one child more efficiently than
// by iterating over all the children as is done with Node.getChildCount().
//
// Java: Node.hasMoreThanOneChild()
func (n Node) HasMoreThanOneChild() bool {
	return n.rec().first != 0 && n.firstNode().rec().next != 0
}

// HasXChildren checks for has exactly the number of specified children.
//
// Java: Node.hasXChildren(int)
func (n Node) HasXChildren(x int) bool {
	c := 0
	for c2 := n.firstNode(); !c2.IsNil() && c <= x; c2 = c2.nextNode() {
		c++
	}
	return c == x
}

// Java: Node.getChildCount()
func (n Node) GetChildCount() int {
	c := 0
	for c2 := n.firstNode(); !c2.IsNil(); c2 = c2.nextNode() {
		c++
	}
	return c
}

// HasChild is intended for testing and verification only.
//
// Java: Node.hasChild(Node)
func (n Node) HasChild(child Node) bool {
	for c := n.firstNode(); !c.IsNil(); c = c.nextNode() {
		if child == c {
			return true
		}
	}
	return false
}

// Node_RecursionMode is the recursion mode for isEquivalentTo.
//
// Java: Node.RecursionMode
type Node_RecursionMode int32

const (
	// Do not compare children or shadow AST.
	Node_RecursionMode_SHALLOW Node_RecursionMode = iota
	// Recursively compare children but not shadow AST.
	Node_RecursionMode_DEEP_NO_SHADOW
	// Recursively compare children and shadow AST.
	Node_RecursionMode_DEEP
)

var node_RecursionMode_names = [...]string{"SHALLOW", "DEEP_NO_SHADOW", "DEEP"}

// Node_RecursionMode_values is RecursionMode.values().
func Node_RecursionMode_values() []Node_RecursionMode {
	return []Node_RecursionMode{Node_RecursionMode_SHALLOW, Node_RecursionMode_DEEP_NO_SHADOW, Node_RecursionMode_DEEP}
}

// String is Enum.name().
func (m Node_RecursionMode) String() string { return node_RecursionMode_names[m] }

// Java: Node.RecursionMode.recurseShadowAst()
func (m Node_RecursionMode) RecurseShadowAst() bool {
	return m == Node_RecursionMode_DEEP
}

// Java: Node.RecursionMode.recurseChildren()
func (m Node_RecursionMode) RecurseChildren() bool {
	return m == Node_RecursionMode_DEEP || m == Node_RecursionMode_DEEP_NO_SHADOW
}

// Node_TypeComparison is the type comparison mode for isEquivalentTo.
//
// Java: Node.TypeComparison
type Node_TypeComparison int32

const (
	// Ignore JSTypes/Colors.
	Node_TypeComparison_IGNORE Node_TypeComparison = iota
	// Compare JSTypes/Colors.
	Node_TypeComparison_COMPARE
)

// Node_TypeComparison_values is TypeComparison.values().
func Node_TypeComparison_values() []Node_TypeComparison {
	return []Node_TypeComparison{Node_TypeComparison_IGNORE, Node_TypeComparison_COMPARE}
}

// String is Enum.name().
func (c Node_TypeComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// Node_JsDocComparison is the JsDoc comparison mode for isEquivalentTo.
//
// Java: Node.JsDocComparison
type Node_JsDocComparison int32

const (
	// Ignore JSDoc.
	Node_JsDocComparison_IGNORE Node_JsDocComparison = iota
	// Compare JSDoc.
	Node_JsDocComparison_COMPARE
)

// Node_JsDocComparison_values is JsDocComparison.values().
func Node_JsDocComparison_values() []Node_JsDocComparison {
	return []Node_JsDocComparison{Node_JsDocComparison_IGNORE, Node_JsDocComparison_COMPARE}
}

// String is Enum.name().
func (c Node_JsDocComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// Node_SideEffectComparison is the side effect comparison mode for
// isEquivalentTo.
//
// Java: Node.SideEffectComparison
type Node_SideEffectComparison int32

const (
	// Ignore side effect flags.
	Node_SideEffectComparison_IGNORE Node_SideEffectComparison = iota
	// Compare side effect flags.
	Node_SideEffectComparison_COMPARE
)

// Node_SideEffectComparison_values is SideEffectComparison.values().
func Node_SideEffectComparison_values() []Node_SideEffectComparison {
	return []Node_SideEffectComparison{Node_SideEffectComparison_IGNORE, Node_SideEffectComparison_COMPARE}
}

// String is Enum.name().
func (c Node_SideEffectComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// IsEquivalentToShallow checks equivalence without going into child nodes.
//
// Java: Node.isEquivalentToShallow(Node)
func (n Node) IsEquivalentToShallow(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_SHALLOW,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_IGNORE)
}

// IsEquivalentWithSideEffectsTo returns true if this node is equivalent
// semantically to another including side effects.
//
// Java: Node.isEquivalentWithSideEffectsTo(Node)
func (n Node) IsEquivalentWithSideEffectsTo(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_COMPARE)
}

// IsEquivalentWithSideEffectsToShallow returns true if this node is
// equivalent semantically to another including side effects.
//
// Java: Node.isEquivalentWithSideEffectsToShallow(Node)
func (n Node) IsEquivalentWithSideEffectsToShallow(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_SHALLOW,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_COMPARE)
}

// IsEquivalentToTyped returns true if this node is equivalent semantically
// to another and the types are equivalent.
//
// Java: Node.isEquivalentToTyped(Node)
func (n Node) IsEquivalentToTyped(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_COMPARE,
		Node_JsDocComparison_COMPARE,
		Node_SideEffectComparison_IGNORE)
}

// IsEquivalentTo returns true if this node is equivalent semantically to
// another.
//
// Java: Node.isEquivalentTo(Node)
func (n Node) IsEquivalentTo(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_IGNORE)
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

// IsEquivalentToWithFlags returns whether this node is equivalent
// semantically to the provided node.
//
// compareType: whether to compare the JSTypes of the nodes. recurse: whether
// to compare the children of the current node. If not, only the count of
// the children are compared. jsDoc: whether to check that the JsDoc of the
// nodes are equivalent.
//
// Deprecated: Use IsEquivalentToWithModes instead.
//
// Java: Node.isEquivalentTo(Node, boolean, boolean, boolean)
func (n Node) IsEquivalentToWithFlags(node Node, compareType bool, recurse bool, jsDoc bool) bool {
	return n.IsEquivalentToWithFlagsAndSideEffect(node, compareType, recurse, jsDoc, false)
}

// IsEquivalentToWithFlagsAndSideEffect returns whether this node is
// equivalent semantically to the provided node.
//
// sideEffect: whether to check that the side-effect flags of the nodes are
// equivalent.
//
// Deprecated: Use IsEquivalentToWithModes instead.
//
// Java: Node.isEquivalentTo(Node, boolean, boolean, boolean, boolean)
func (n Node) IsEquivalentToWithFlagsAndSideEffect(node Node, compareType bool, recurse bool, jsDoc bool, sideEffect bool) bool {
	recursion := Node_RecursionMode_SHALLOW
	if recurse {
		recursion = Node_RecursionMode_DEEP
	}
	typeConfig := Node_TypeComparison_IGNORE
	if compareType {
		typeConfig = Node_TypeComparison_COMPARE
	}
	jsDocConfig := Node_JsDocComparison_IGNORE
	if jsDoc {
		jsDocConfig = Node_JsDocComparison_COMPARE
	}
	sideEffectConfig := Node_SideEffectComparison_IGNORE
	if sideEffect {
		sideEffectConfig = Node_SideEffectComparison_COMPARE
	}
	return n.IsEquivalentToWithModes(node, recursion, typeConfig, jsDocConfig, sideEffectConfig)
}

// IsEquivalentToWithModes returns whether this node is equivalent
// semantically to the provided node.
//
// Go-only: language-adapter — dispatches to the overrides of the folded
// subclasses (NumberNode, BigIntNode, StringNode,
// TemplateLiteralSubstringNode), each of which calls the base
// implementation first.
//
// Java: Node.isEquivalentTo(Node, RecursionMode, TypeComparison, JsDocComparison, SideEffectComparison)
func (n Node) IsEquivalentToWithModes(
	node Node,
	recursion Node_RecursionMode,
	typeConfig Node_TypeComparison,
	jsDocConfig Node_JsDocComparison,
	sideEffectConfig Node_SideEffectComparison,
) bool {
	switch n.nodeKind() {
	case node_kindNumber:
		// Java: NumberNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			(n.node_number() == node.node_number()) // -0.0 and NaN are forbidden.
	case node_kindBigInt:
		// Java: BigIntNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			jutil.BigIntegerEquals(n.GetBigInt(), node.GetBigInt())
	case node_kindString:
		// Java: StringNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			uint32(n.rec().payload) == uint32(node.rec().payload)
	case node_kindTemplateLiteralSubstring:
		// Java: TemplateLiteralSubstringNode.isEquivalentTo
		if !n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
			return false
		}
		return n.rec().payload == node.rec().payload && n.node_cookedID() == node.node_cookedID()
	}
	return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig)
}

// isEquivalentToBase is Node's own isEquivalentTo (the super call of the
// folded subclasses).
func (n Node) isEquivalentToBase(
	node Node,
	recursion Node_RecursionMode,
	typeConfig Node_TypeComparison,
	jsDocConfig Node_JsDocComparison,
	sideEffectConfig Node_SideEffectComparison,
) bool {
	if n.rec().token != node.rec().token ||
		n.GetChildCount() != node.GetChildCount() ||
		n.nodeKind() != node.nodeKind() {
		return false
	}

	if typeConfig == Node_TypeComparison_COMPARE &&
		!jutil.ObjectsEquals(n.typeOrColorValue(), node.typeOrColorValue()) {
		return false
	}

	if jsDocConfig == Node_JsDocComparison_COMPARE &&
		!JSDocInfo_AreEquivalent(n.GetJSDocInfo(), node.GetJSDocInfo()) {
		return false
	}

	thisDte := n.GetDeclaredTypeExpression()
	thatDte := node.GetDeclaredTypeExpression()
	if thisDte == thatDte {
		// Do nothing
	} else if thisDte.IsNil() || thatDte.IsNil() {
		return false
	} else if !thisDte.IsEquivalentToWithModes(
		thatDte, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
		return false
	}

	if (n.props() & node_propEqualityMask) != (node.props() & node_propEqualityMask) {
		return false
	}
	if ((n.rec().kindFlags | node.rec().kindFlags) & kindFlagHasExtra) != 0 {
		if n.GetIntProp(Node_Prop_EXPORT_DEFAULT) != node.GetIntProp(Node_Prop_EXPORT_DEFAULT) ||
			n.GetIntProp(Node_Prop_EXPORT_ALL_FROM) != node.GetIntProp(Node_Prop_EXPORT_ALL_FROM) ||
			n.GetIntProp(Node_Prop_INCRDECR) != node.GetIntProp(Node_Prop_INCRDECR) ||
			n.GetIntProp(Node_Prop_QUOTED) != node.GetIntProp(Node_Prop_QUOTED) {
			return false
		}
	}

	if sideEffectConfig == Node_SideEffectComparison_COMPARE {
		if n.GetSideEffectFlags() != node.GetSideEffectFlags() {
			return false
		}

		if n.IsUnusedParameter() != node.IsUnusedParameter() {
			return false
		}
	}

	thisIsShadowHost := !n.GetClosureUnawareShadow().IsNil()
	nodeIsShadowHost := !node.GetClosureUnawareShadow().IsNil()
	if thisIsShadowHost != nodeIsShadowHost {
		return false
	}

	if recursion.RecurseShadowAst() && thisIsShadowHost && nodeIsShadowHost {
		thisShadow := n.GetClosureUnawareShadow()
		nodeShadow := node.GetClosureUnawareShadow()
		if !thisShadow.IsEquivalentToWithModes(
			nodeShadow, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
			return false
		}
	}

	if recursion.RecurseChildren() {
		for c, c2 := n.firstNode(), node.firstNode(); !c.IsNil(); c, c2 = c.nextNode(), c2.nextNode() {
			if !c.IsEquivalentToWithModes(c2, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
				return false
			}
		}
	}

	return true
}

// Node_PROP_MAP_FOR_EQUALITY holds accessors for Node properties that should
// also be compared when comparing nodes for equality.
//
// We'd prefer to list the props that should be ignored rather than the ones
// that should be checked, but that was too difficult initially. In general
// these properties are ones that show up as keywords / symbols which don't
// have their own nodes.
//
// Accessor functions are used rather than Props to encode the correct way
// of reading the prop.
//
// Java: Node.PROP_MAP_FOR_EQUALITY
var Node_PROP_MAP_FOR_EQUALITY = func() *jutil.LinkedHashMap[Node_Prop, func(Node) any] {
	m := jutil.NewLinkedHashMap[Node_Prop, func(Node) any]()
	m.Put(Node_Prop_ARROW_FN, func(n Node) any { return n.IsArrowFunction() })
	m.Put(Node_Prop_ASYNC_FN, func(n Node) any { return n.IsAsyncFunction() })
	m.Put(Node_Prop_GENERATOR_FN, func(n Node) any { return n.IsGeneratorFunction() })
	m.Put(Node_Prop_START_OF_OPT_CHAIN, func(n Node) any { return n.IsOptionalChainStart() })
	m.Put(Node_Prop_STATIC_MEMBER, func(n Node) any { return n.IsStaticMember() })
	m.Put(Node_Prop_YIELD_ALL, func(n Node) any { return n.IsYieldAll() })
	m.Put(Node_Prop_EXPORT_DEFAULT, func(n Node) any { return n.GetIntProp(Node_Prop_EXPORT_DEFAULT) })
	m.Put(Node_Prop_EXPORT_ALL_FROM, func(n Node) any { return n.GetIntProp(Node_Prop_EXPORT_ALL_FROM) })
	m.Put(Node_Prop_INCRDECR, func(n Node) any { return n.GetIntProp(Node_Prop_INCRDECR) })
	m.Put(Node_Prop_QUOTED, func(n Node) any { return n.GetIntProp(Node_Prop_QUOTED) })
	m.Put(Node_Prop_FREE_CALL, func(n Node) any { return n.GetBooleanProp(Node_Prop_FREE_CALL) })
	m.Put(Node_Prop_DIRECT_EVAL, func(n Node) any { return n.GetBooleanProp(Node_Prop_DIRECT_EVAL) })
	m.Put(Node_Prop_COMPUTED_PROP_METHOD, func(n Node) any { return n.GetBooleanProp(Node_Prop_COMPUTED_PROP_METHOD) })
	m.Put(Node_Prop_COMPUTED_PROP_GETTER, func(n Node) any { return n.GetBooleanProp(Node_Prop_COMPUTED_PROP_GETTER) })
	m.Put(Node_Prop_COMPUTED_PROP_SETTER, func(n Node) any { return n.GetBooleanProp(Node_Prop_COMPUTED_PROP_SETTER) })
	return m
}()

// Node_propMapForEqualityKeys is used for faster Map.containsKey() lookups
// in PROP_MAP_FOR_EQUALITY.
//
// Java: Node.propMapForEqualityKeys
var Node_propMapForEqualityKeys = jutil.NewEnumSet(Node_PROP_MAP_FOR_EQUALITY.KeySlice()...)

const node_propEqualityMask uint64 = propFlag_ARROW_FN |
	propFlag_ASYNC_FN |
	propFlag_GENERATOR_FN |
	propFlag_START_OF_OPT_CHAIN |
	propFlag_STATIC_MEMBER |
	propFlag_YIELD_ALL |
	propFlag_EXPORT_DEFAULT |
	propFlag_EXPORT_ALL_FROM |
	propFlag_INCRDECR |
	propFlag_QUOTED |
	propFlag_FREE_CALL |
	propFlag_DIRECT_EVAL |
	propFlag_COMPUTED_PROP_METHOD |
	propFlag_COMPUTED_PROP_GETTER |
	propFlag_COMPUTED_PROP_SETTER

// GetQualifiedName takes a set of GETPROP nodes and produces a string that
// is each property separated by dots. If the node ultimately under the left
// sub-tree is not a simple name, this is not a valid qualified name.
//
// Returns null if this is not a qualified name, or a dot-separated string of
// the name and properties.
//
// Java: Node.getQualifiedName()
//
// Go-only: performance — the body is in getQualifiedNameString, so that this
// wrapper inlines: a caller that only reads the result keeps the string
// header on its stack instead of the heap (these headers were 6% of
// checks-only's allocations).
func (n Node) GetQualifiedName() *string {
	s, ok := n.getQualifiedNameString()
	if !ok {
		return nil
	}
	return &s
}

// GetQualifiedNameString is the body of GetQualifiedName; ok is false
// where Java returns null.
func (n Node) GetQualifiedNameString() (string, bool) {
	if n.IsNil() {
		return "", false
	}
	return n.getQualifiedNameString()
}

// getQualifiedNameString is the body of getQualifiedName; ok is false
// where Java returns null.
func (n Node) getQualifiedNameString() (string, bool) {
	switch n.rec().token {
	case Token_NAME:
		name := n.GetString()
		if name == "" {
			return "", false
		}
		return name, true
	case Token_GETPROP:
		first := n.firstNode()
		if !first.IsGetProp() {
			left, ok := first.getQualifiedNameString()
			if !ok {
				return "", false
			}
			return left + "." + n.GetString(), true
		}
		var partsBuf [8]string
		parts := partsBuf[:0]
		totalLen := 0
		cur := n
		for cur.IsGetProp() {
			p := cur.GetString()
			parts = append(parts, p)
			totalLen += 1 + len(p)
			cur = cur.firstNode()
		}
		left, ok := cur.getQualifiedNameString()
		if !ok {
			return "", false
		}
		var builder strings.Builder
		builder.Grow(len(left) + totalLen)
		builder.WriteString(left)
		for i := len(parts) - 1; i >= 0; i-- {
			builder.WriteByte('.')
			builder.WriteString(parts[i])
		}
		return builder.String(), true
	case Token_THIS:
		return "this", true
	case Token_SUPER:
		return "super", true
	default:
		return "", false
	}
}

// Java: Node.getQualifiedNameObject()
func (n Node) GetQualifiedNameObject() QualifiedName {
	if n.IsQualifiedName() {
		return NewQualifiedName_NodeQname(n)
	}
	return nil
}

// GetOriginalQualifiedName takes a set of GETPROP nodes and produces a
// string that is each property separated by dots. If the node ultimately
// under the left sub-tree is not a simple name, this is not a valid
// qualified name. This method returns the original name of each segment
// rather than the renamed version.
//
// Deprecated: "original name" is poorly defined. See #getOriginalName
//
// Java: Node.getOriginalQualifiedName()
func (n Node) GetOriginalQualifiedName() *string {
	if n.rec().token == Token_NAME {
		var name string
		if n.GetOriginalName() == nil {
			name = n.GetString()
		} else {
			name = *n.GetOriginalName()
		}
		if name == "" {
			return nil
		}
		return &name
	} else if n.rec().token == Token_GETPROP {
		left := n.GetFirstChild().GetOriginalQualifiedName()
		if left == nil {
			return nil
		}

		var right string
		if n.GetOriginalName() == nil {
			right = n.GetString()
		} else {
			right = *n.GetOriginalName()
		}

		s := *left + "." + right
		return &s
	} else if n.rec().token == Token_THIS {
		s := "this"
		return &s
	} else if n.rec().token == Token_SUPER {
		s := "super"
		return &s
	} else {
		return nil
	}
}

// IsQualifiedName returns whether a node corresponds to a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.isQualifiedName()
func (n Node) IsQualifiedName() bool {
	for {
		switch n.rec().token {
		case Token_NAME:
			return n.GetStringID() > 1
		case Token_THIS, Token_SUPER:
			return true
		case Token_GETPROP:
			n = n.firstNode()
		default:
			return false
		}
	}
}

// MatchesName returns whether a node matches a simple name, such as `x`,
// returns false if this is not a NAME node.
//
// Java: Node.matchesName(String)
func (n Node) MatchesName(name string) bool {
	if n.rec().token != Token_NAME {
		return false
	}
	internalString := n.GetString()
	return internalString != "" && name == internalString
}

// MatchesNameNode checks that if two NAME node match, returns false if
// either node is not a NAME node. As a empty string is not considered a
// valid Name (it is an AST placeholder), empty strings are never considered
// to be matches.
//
// Java: Node.matchesName(Node)
func (n Node) MatchesNameNode(other Node) bool {
	if n.rec().token != Token_NAME || other.rec().token != Token_NAME {
		return false
	}

	// Compare interned string IDs directly (id > 1 means non-empty!).
	id := rhinoStringID(n.rec().payload)
	return id > 1 && id == rhinoStringID(other.rec().payload)
}

// MatchesQualifiedName returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(String)
func (n Node) MatchesQualifiedName(qname string) bool {
	for {
		switch n.rec().token {
		case Token_NAME, Token_IMPORT_STAR:
			name := n.GetString()
			return name != "" && name == qname && strings.IndexByte(name, '.') < 0
		case Token_THIS:
			return qname == "this"
		case Token_SUPER:
			return qname == "super"
		case Token_GETPROP:
			prop := n.GetString()
			dot := len(qname) - len(prop) - 1
			if dot < 1 || qname[dot] != '.' || qname[dot+1:] != prop || strings.IndexByte(prop, '.') >= 0 {
				return false
			}
			qname = qname[:dot]
			n = n.firstNode()
		default:
			return false
		}
	}
}

// MatchesQualifiedNameWithEnd returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(String, int)
func (n Node) MatchesQualifiedNameWithEnd(qname string, endIndex int) bool {
	start := jutil.StringLastIndexOfCharFrom(qname, '.', endIndex-1) + 1

	switch n.GetToken() {
	case Token_NAME, Token_IMPORT_STAR:
		name := n.GetString()
		return start == 0 && name != "" && jutil.StringLength(name) == endIndex && jutil.StringStartsWith(qname, name)
	case Token_THIS:
		return start == 0 && endIndex == 4 && jutil.StringStartsWith(qname, "this")
	case Token_SUPER:
		return start == 0 && endIndex == 5 && jutil.StringStartsWith(qname, "super")
	case Token_GETPROP:
		prop := n.GetString()
		return start > 1 &&
			jutil.StringLength(prop) == endIndex-start &&
			jutil.StringRegionMatches(prop, 0, qname, start, endIndex-start) &&
			n.GetFirstChild().MatchesQualifiedNameWithEnd(qname, start-1)

	case Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// MatchesQualifiedNameNode returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(Node)
func (n Node) MatchesQualifiedNameNode(other Node) bool {
	if other.rec().token != n.rec().token {
		return false
	}

	switch n.rec().token {
	case Token_NAME:
		return n.MatchesNameNode(other)
	case Token_THIS, Token_SUPER:
		return true
	case Token_GETPROP:
		return uint32(n.rec().payload) == uint32(other.rec().payload) &&
			n.GetFirstChild().MatchesQualifiedNameNode(other.GetFirstChild())

	case Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// IsUnscopedQualifiedName returns whether a node corresponds to a simple or
// a qualified name without a "this" reference, such as `a.b.c`, but not
// `this.a`.
//
// Java: Node.isUnscopedQualifiedName()
func (n Node) IsUnscopedQualifiedName() bool {
	for {
		switch n.rec().token {
		case Token_NAME:
			return n.GetStringID() > 1
		case Token_GETPROP:
			n = n.firstNode()
		default:
			return false
		}
	}
}

// Java: Node.isValidAssignmentTarget()
func (n Node) IsValidAssignmentTarget() bool {
	switch n.GetToken() {
	case Token_NAME, Token_GETPROP, Token_GETELEM, Token_ARRAY_PATTERN, Token_OBJECT_PATTERN:
		return true
	default:
		return false
	}
}

// Clone always throws: use CloneNode.
//
// Java: Node.clone()
func (n Node) Clone() any {
	panic(jutil.NewUnsupportedOperationExceptionMsg("Did you mean cloneNode?"))
}

// CloneNode returns a detached clone of the Node, specifically excluding its
// children.
//
// Java: Node.cloneNode()
func (n Node) CloneNode() Node {
	return n.CloneNodeWithTypeExprs(false)
}

// CloneNodeWithTypeExprs returns a detached clone of the Node, specifically
// excluding its children.
//
// Go-only: language-adapter — dispatches to the cloneNode overrides of the
// folded subclasses.
//
// Java: Node.cloneNode(boolean)
func (n Node) CloneNodeWithTypeExprs(cloneTypeExprs bool) Node {
	clone := allocNode()
	clone.rec().token = n.rec().token
	clone.rec().kindFlags = n.rec().kindFlags & kindMask
	clone.rec().payload = n.rec().payload
	Node_CopyBaseNodeFields(n, clone, cloneTypeExprs)
	return clone
}

// Java: Node.copyBaseNodeFields(Node, Node, boolean)
func Node_CopyBaseNodeFields(source Node, dest Node, cloneTypeExprs bool) {
	dest.rec().linenoCharno = source.rec().linenoCharno
	dest.setTC(dest.store().translateTC(source.store(), source.tc(), source.rec().kindFlags&kindFlagTypeIsColor != 0))
	dest.rec().kindFlags = (dest.rec().kindFlags &^ (kindFlagTypeIsColor | kindFlagHasExtBox | kindFlagHasExtra)) |
		(source.rec().kindFlags & (kindFlagTypeIsColor | kindFlagHasExtBox))
	if srcBox := source.extBoxOrNil(); srcBox != nil {
		srcBox.propListHead.markShared()
		boxCopy := *srcBox
		boxCopy.typeBeforeCast = dest.store().translateTypeHandle(source.store(), boxCopy.typeBeforeCast)
		dest.setAux(nodeAux{ext: &boxCopy})
	} else {
		dest.setAux(source.auxValue())
	}
	if source.rec().kindFlags&kindFlagHasExtra != 0 && source.coldPropListHead() != nil {
		dest.rec().kindFlags |= kindFlagHasExtra
	}
	dest.copyLengthFrom(source)
	dest.node_setOriginalNameID(source.node_originalNameID())
	dest.setPropsAndSourceFile(source)

	// TODO(johnlenz): Remove this once JSTypeExpression are immutable
	if cloneTypeExprs {
		info := source.GetJSDocInfo()
		if info != nil {
			dest.SetJSDocInfo(info.CloneWithTypeNodes(true))
		}
	}
}

// CloneTree returns a detached clone of the Node and all its children.
//
// Java: Node.cloneTree()
func (n Node) CloneTree() Node {
	return n.CloneTreeWithTypeExprs(false)
}

// Java: Node.cloneTree(boolean)
func (n Node) CloneTreeWithTypeExprs(cloneTypeExprs bool) Node {
	result := n.CloneNodeWithTypeExprs(cloneTypeExprs)
	if n.rec().first != 0 {
		var firstChild Node
		var lastChild Node
		for n2 := n.firstNode(); !n2.IsNil(); n2 = n2.nextNode() {
			n2clone := n2.CloneTreeWithTypeExprs(cloneTypeExprs)
			n2clone.rec().parent = n2clone.ref(result)
			if firstChild.IsNil() {
				firstChild = n2clone
				lastChild = firstChild
			} else {
				lastChild.rec().next = lastChild.ref(n2clone)
				n2clone.rec().previous = n2clone.ref(lastChild)
				lastChild = n2clone
			}
		}
		firstChild.rec().previous = firstChild.ref(lastChild)
		lastChild.rec().next = 0
		result.rec().first = result.ref(firstChild)
	}

	if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
		result.SetClosureUnawareShadow(shadow.CloneTreeWithTypeExprs(cloneTypeExprs))
	}

	return result
}

// Srcref copies the source info from `other` onto `this`.
//
// Java: Node.srcref(Node)
func (n Node) Srcref(other Node) Node {
	n.SetStaticSourceFileFrom(other)
	n.node_setOriginalNameID(other.node_originalNameID())
	n.rec().linenoCharno = other.rec().linenoCharno
	n.copyLengthFrom(other)
	return n
}

// SrcrefTree copies, for all Nodes in the subtree of `this`, the source info
// from `other`.
//
// Java: Node.srcrefTree(Node)
func (n Node) SrcrefTree(other Node) Node {
	n.Srcref(other)
	for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
		child.SrcrefTree(other)
	}
	return n
}

// SrcrefIfMissing copies the source info from `other` iff source info is not
// set on `this`.
//
// Java: Node.srcrefIfMissing(Node)
func (n Node) SrcrefIfMissing(other Node) Node {
	if n.rec().sourceFile == 0 {
		n.SetStaticSourceFileFrom(other)
		n.rec().linenoCharno = other.rec().linenoCharno
		n.copyLengthFrom(other)
	}

	// TODO(lharker): should this be inside the above if condition?
	// If the node already has a source file, it seems strange to
	// go ahead and set the original name anyway.
	if n.node_originalNameID() == 0 {
		n.node_setOriginalNameID(other.node_originalNameID())
	}

	return n
}

// SrcrefTreeIfMissing copies, for all Nodes in the subtree of `this`, iff
// source info is not set, the source info from `other`.
//
// Java: Node.srcrefTreeIfMissing(Node)
func (n Node) SrcrefTreeIfMissing(other Node) Node {
	n.SrcrefIfMissing(other)
	for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
		child.SrcrefTreeIfMissing(other)
	}
	return n
}

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

// AST type check methods

// Java: Node.isAdd()
func (n Node) IsAdd() bool {
	return n.rec().token == Token_ADD
}

// Java: Node.isSub()
func (n Node) IsSub() bool {
	return n.rec().token == Token_SUB
}

// Java: Node.isAnd()
func (n Node) IsAnd() bool {
	return n.rec().token == Token_AND
}

// Java: Node.isAssignAnd()
func (n Node) IsAssignAnd() bool {
	return n.rec().token == Token_ASSIGN_AND
}

// Java: Node.isArrayLit()
func (n Node) IsArrayLit() bool {
	return n.rec().token == Token_ARRAYLIT
}

// Java: Node.isArrayPattern()
func (n Node) IsArrayPattern() bool {
	return n.rec().token == Token_ARRAY_PATTERN
}

// Java: Node.isAssign()
func (n Node) IsAssign() bool {
	return n.rec().token == Token_ASSIGN
}

// Java: Node.isAssignAdd()
func (n Node) IsAssignAdd() bool {
	return n.rec().token == Token_ASSIGN_ADD
}

// Java: Node.isNormalBlock()
func (n Node) IsNormalBlock() bool {
	return n.IsBlock()
}

// Java: Node.isBlock()
func (n Node) IsBlock() bool {
	return n.rec().token == Token_BLOCK
}

// Java: Node.isRoot()
func (n Node) IsRoot() bool {
	return n.rec().token == Token_ROOT
}

// Java: Node.isAwait()
func (n Node) IsAwait() bool {
	return n.rec().token == Token_AWAIT
}

// Java: Node.isBigInt()
func (n Node) IsBigInt() bool {
	return n.rec().token == Token_BIGINT
}

// Java: Node.isBitNot()
func (n Node) IsBitNot() bool {
	return n.rec().token == Token_BITNOT
}

// Java: Node.isBreak()
func (n Node) IsBreak() bool {
	return n.rec().token == Token_BREAK
}

// Java: Node.isCall()
func (n Node) IsCall() bool {
	return n.rec().token == Token_CALL
}

// Java: Node.isCase()
func (n Node) IsCase() bool {
	return n.rec().token == Token_CASE
}

// Java: Node.isCast()
func (n Node) IsCast() bool {
	return n.rec().token == Token_CAST
}

// Java: Node.isCatch()
func (n Node) IsCatch() bool {
	return n.rec().token == Token_CATCH
}

// Java: Node.isClass()
func (n Node) IsClass() bool {
	return n.rec().token == Token_CLASS
}

// Java: Node.isClassMembers()
func (n Node) IsClassMembers() bool {
	return n.rec().token == Token_CLASS_MEMBERS
}

// Java: Node.isComma()
func (n Node) IsComma() bool {
	return n.rec().token == Token_COMMA
}

// Java: Node.isComputedProp()
func (n Node) IsComputedProp() bool {
	return n.rec().token == Token_COMPUTED_PROP
}

// Java: Node.isContinue()
func (n Node) IsContinue() bool {
	return n.rec().token == Token_CONTINUE
}

// Java: Node.isConst()
func (n Node) IsConst() bool {
	return n.rec().token == Token_CONST
}

// Java: Node.isDebugger()
func (n Node) IsDebugger() bool {
	return n.rec().token == Token_DEBUGGER
}

// Java: Node.isDec()
func (n Node) IsDec() bool {
	return n.rec().token == Token_DEC
}

// Java: Node.isDefaultCase()
func (n Node) IsDefaultCase() bool {
	return n.rec().token == Token_DEFAULT_CASE
}

// Java: Node.isDefaultValue()
func (n Node) IsDefaultValue() bool {
	return n.rec().token == Token_DEFAULT_VALUE
}

// Java: Node.isDelProp()
func (n Node) IsDelProp() bool {
	return n.rec().token == Token_DELPROP
}

// Java: Node.isDestructuringLhs()
func (n Node) IsDestructuringLhs() bool {
	return n.rec().token == Token_DESTRUCTURING_LHS
}

// Java: Node.isDestructuringPattern()
func (n Node) IsDestructuringPattern() bool {
	return n.IsObjectPattern() || n.IsArrayPattern()
}

// Java: Node.isDo()
func (n Node) IsDo() bool {
	return n.rec().token == Token_DO
}

// Java: Node.isEmpty()
func (n Node) IsEmpty() bool {
	return n.rec().token == Token_EMPTY
}

// Java: Node.isExponent()
func (n Node) IsExponent() bool {
	return n.rec().token == Token_EXPONENT
}

// Java: Node.isAssignExponent()
func (n Node) IsAssignExponent() bool {
	return n.rec().token == Token_ASSIGN_EXPONENT
}

// Java: Node.isExport()
func (n Node) IsExport() bool {
	return n.rec().token == Token_EXPORT
}

// Java: Node.isExportSpec()
func (n Node) IsExportSpec() bool {
	return n.rec().token == Token_EXPORT_SPEC
}

// Java: Node.isExportSpecs()
func (n Node) IsExportSpecs() bool {
	return n.rec().token == Token_EXPORT_SPECS
}

// Java: Node.isExprResult()
func (n Node) IsExprResult() bool {
	return n.rec().token == Token_EXPR_RESULT
}

// Java: Node.isFalse()
func (n Node) IsFalse() bool {
	return n.rec().token == Token_FALSE
}

// Java: Node.isVanillaFor()
func (n Node) IsVanillaFor() bool {
	return n.rec().token == Token_FOR
}

// Java: Node.isForIn()
func (n Node) IsForIn() bool {
	return n.rec().token == Token_FOR_IN
}

// Java: Node.isForOf()
func (n Node) IsForOf() bool {
	return n.rec().token == Token_FOR_OF
}

// Java: Node.isForAwaitOf()
func (n Node) IsForAwaitOf() bool {
	return n.rec().token == Token_FOR_AWAIT_OF
}

// Java: Node.isFunction()
func (n Node) IsFunction() bool {
	return n.rec().token == Token_FUNCTION
}

// Java: Node.isGetterDef()
func (n Node) IsGetterDef() bool {
	return n.rec().token == Token_GETTER_DEF
}

// Java: Node.isGetElem()
func (n Node) IsGetElem() bool {
	return n.rec().token == Token_GETELEM
}

// Java: Node.isGetProp()
func (n Node) IsGetProp() bool {
	return n.rec().token == Token_GETPROP
}

// Java: Node.isHook()
func (n Node) IsHook() bool {
	return n.rec().token == Token_HOOK
}

// Java: Node.isIf()
func (n Node) IsIf() bool {
	return n.rec().token == Token_IF
}

// Java: Node.isImport()
func (n Node) IsImport() bool {
	return n.rec().token == Token_IMPORT
}

// Java: Node.isImportMeta()
func (n Node) IsImportMeta() bool {
	return n.rec().token == Token_IMPORT_META
}

// Java: Node.isImportStar()
func (n Node) IsImportStar() bool {
	return n.rec().token == Token_IMPORT_STAR
}

// Java: Node.isImportSpec()
func (n Node) IsImportSpec() bool {
	return n.rec().token == Token_IMPORT_SPEC
}

// Java: Node.isImportSpecs()
func (n Node) IsImportSpecs() bool {
	return n.rec().token == Token_IMPORT_SPECS
}

// Java: Node.isIn()
func (n Node) IsIn() bool {
	return n.rec().token == Token_IN
}

// Java: Node.isInc()
func (n Node) IsInc() bool {
	return n.rec().token == Token_INC
}

// Java: Node.isInstanceOf()
func (n Node) IsInstanceOf() bool {
	return n.rec().token == Token_INSTANCEOF
}

// Java: Node.isInterface()
func (n Node) IsInterface() bool {
	return n.rec().token == Token_INTERFACE
}

// Java: Node.isInterfaceMembers()
func (n Node) IsInterfaceMembers() bool {
	return n.rec().token == Token_INTERFACE_MEMBERS
}

// Java: Node.isRecordType()
func (n Node) IsRecordType() bool {
	return n.rec().token == Token_RECORD_TYPE
}

// Java: Node.isCallSignature()
func (n Node) IsCallSignature() bool {
	return n.rec().token == Token_CALL_SIGNATURE
}

// Java: Node.isIndexSignature()
func (n Node) IsIndexSignature() bool {
	return n.rec().token == Token_INDEX_SIGNATURE
}

// Java: Node.isLabel()
func (n Node) IsLabel() bool {
	return n.rec().token == Token_LABEL
}

// Java: Node.isLabelName()
func (n Node) IsLabelName() bool {
	return n.rec().token == Token_LABEL_NAME
}

// Java: Node.isLet()
func (n Node) IsLet() bool {
	return n.rec().token == Token_LET
}

// Java: Node.isMemberFunctionDef()
func (n Node) IsMemberFunctionDef() bool {
	return n.rec().token == Token_MEMBER_FUNCTION_DEF
}

// Java: Node.isMemberVariableDef()
func (n Node) IsMemberVariableDef() bool {
	return n.rec().token == Token_MEMBER_VARIABLE_DEF
}

// Java: Node.isMemberFieldDef()
func (n Node) IsMemberFieldDef() bool {
	return n.rec().token == Token_MEMBER_FIELD_DEF
}

// Java: Node.isComputedFieldDef()
func (n Node) IsComputedFieldDef() bool {
	return n.rec().token == Token_COMPUTED_FIELD_DEF
}

// Java: Node.isModuleBody()
func (n Node) IsModuleBody() bool {
	return n.rec().token == Token_MODULE_BODY
}

// Java: Node.isName()
func (n Node) IsName() bool {
	return n.rec().token == Token_NAME
}

// Java: Node.isNE()
func (n Node) IsNE() bool {
	return n.rec().token == Token_NE
}

// Java: Node.isSHNE()
func (n Node) IsSHNE() bool {
	return n.rec().token == Token_SHNE
}

// Java: Node.isEQ()
func (n Node) IsEQ() bool {
	return n.rec().token == Token_EQ
}

// Java: Node.isSHEQ()
func (n Node) IsSHEQ() bool {
	return n.rec().token == Token_SHEQ
}

// Java: Node.isNeg()
func (n Node) IsNeg() bool {
	return n.rec().token == Token_NEG
}

// Java: Node.isNew()
func (n Node) IsNew() bool {
	return n.rec().token == Token_NEW
}

// Java: Node.isNot()
func (n Node) IsNot() bool {
	return n.rec().token == Token_NOT
}

// Java: Node.isNull()
func (n Node) IsNull() bool {
	return n.rec().token == Token_NULL
}

// Java: Node.isNullishCoalesce()
func (n Node) IsNullishCoalesce() bool {
	return n.rec().token == Token_COALESCE
}

// Java: Node.isAssignNullishCoalesce()
func (n Node) IsAssignNullishCoalesce() bool {
	return n.rec().token == Token_ASSIGN_COALESCE
}

// Java: Node.isNumber()
func (n Node) IsNumber() bool {
	return n.rec().token == Token_NUMBER
}

// Java: Node.isObjectLit()
func (n Node) IsObjectLit() bool {
	return n.rec().token == Token_OBJECTLIT
}

// Java: Node.isObjectPattern()
func (n Node) IsObjectPattern() bool {
	return n.rec().token == Token_OBJECT_PATTERN
}

// Java: Node.isOptChainCall()
func (n Node) IsOptChainCall() bool {
	return n.rec().token == Token_OPTCHAIN_CALL
}

// Java: Node.isOptChainGetElem()
func (n Node) IsOptChainGetElem() bool {
	return n.rec().token == Token_OPTCHAIN_GETELEM
}

// Java: Node.isOptChainGetProp()
func (n Node) IsOptChainGetProp() bool {
	return n.rec().token == Token_OPTCHAIN_GETPROP
}

// Java: Node.isOr()
func (n Node) IsOr() bool {
	return n.rec().token == Token_OR
}

// Java: Node.isAssignOr()
func (n Node) IsAssignOr() bool {
	return n.rec().token == Token_ASSIGN_OR
}

// Java: Node.isParamList()
func (n Node) IsParamList() bool {
	return n.rec().token == Token_PARAM_LIST
}

// Java: Node.isPos()
func (n Node) IsPos() bool {
	return n.rec().token == Token_POS
}

// Java: Node.isRegExp()
func (n Node) IsRegExp() bool {
	return n.rec().token == Token_REGEXP
}

// Java: Node.isRest()
func (n Node) IsRest() bool {
	return n.rec().token == Token_ITER_REST || n.rec().token == Token_OBJECT_REST
}

// Java: Node.isObjectRest()
func (n Node) IsObjectRest() bool {
	return n.rec().token == Token_OBJECT_REST
}

// Java: Node.isReturn()
func (n Node) IsReturn() bool {
	return n.rec().token == Token_RETURN
}

// Java: Node.isScript()
func (n Node) IsScript() bool {
	return n.rec().token == Token_SCRIPT
}

// Java: Node.isSetterDef()
func (n Node) IsSetterDef() bool {
	return n.rec().token == Token_SETTER_DEF
}

// Java: Node.isSpread()
func (n Node) IsSpread() bool {
	return n.rec().token == Token_ITER_SPREAD || n.rec().token == Token_OBJECT_SPREAD
}

// Java: Node.isString()
func (n Node) IsString() bool {
	return n.rec().token == Token_STRINGLIT
}

// Java: Node.isStringKey()
func (n Node) IsStringKey() bool {
	return n.rec().token == Token_STRING_KEY
}

// Java: Node.isStringLit()
func (n Node) IsStringLit() bool {
	return n.rec().token == Token_STRINGLIT
}

// Java: Node.isSuper()
func (n Node) IsSuper() bool {
	return n.rec().token == Token_SUPER
}

// Java: Node.isSwitch()
func (n Node) IsSwitch() bool {
	return n.rec().token == Token_SWITCH
}

// Java: Node.isSwitchBody()
func (n Node) IsSwitchBody() bool {
	return n.rec().token == Token_SWITCH_BODY
}

// Java: Node.isTaggedTemplateLit()
func (n Node) IsTaggedTemplateLit() bool {
	return n.rec().token == Token_TAGGED_TEMPLATELIT
}

// Java: Node.isTemplateLit()
func (n Node) IsTemplateLit() bool {
	return n.rec().token == Token_TEMPLATELIT
}

// Java: Node.isTemplateLitString()
func (n Node) IsTemplateLitString() bool {
	return n.rec().token == Token_TEMPLATELIT_STRING
}

// Java: Node.isTemplateLitSub()
func (n Node) IsTemplateLitSub() bool {
	return n.rec().token == Token_TEMPLATELIT_SUB
}

// Java: Node.isThis()
func (n Node) IsThis() bool {
	return n.rec().token == Token_THIS
}

// Java: Node.isThrow()
func (n Node) IsThrow() bool {
	return n.rec().token == Token_THROW
}

// Java: Node.isTrue()
func (n Node) IsTrue() bool {
	return n.rec().token == Token_TRUE
}

// Java: Node.isTry()
func (n Node) IsTry() bool {
	return n.rec().token == Token_TRY
}

// Java: Node.isTypeOf()
func (n Node) IsTypeOf() bool {
	return n.rec().token == Token_TYPEOF
}

// Java: Node.isVar()
func (n Node) IsVar() bool {
	return n.rec().token == Token_VAR
}

// Java: Node.isVoid()
func (n Node) IsVoid() bool {
	return n.rec().token == Token_VOID
}

// Java: Node.isWhile()
func (n Node) IsWhile() bool {
	return n.rec().token == Token_WHILE
}

// Java: Node.isWith()
func (n Node) IsWith() bool {
	return n.rec().token == Token_WITH
}

// Java: Node.isYield()
func (n Node) IsYield() bool {
	return n.rec().token == Token_YIELD
}

// Java: Node.isDeclare()
func (n Node) IsDeclare() bool {
	return n.rec().token == Token_DECLARE
}

// Go-only: language-adapter — recover the Java Object field only for printing
// and equality; hot getters use the tagged pointer directly.
func (n Node) typeOrColorValue() any {
	if n.tc() == 0 {
		return nil
	}
	if (n.rec().kindFlags & kindFlagTypeIsColor) != 0 {
		if c := n.store().colorAt(n.tc()); c != nil {
			return c
		}
		return nil
	}
	return jsTypeOf(n.store().typeAt(n.tc()))
}

func (n Node) nodeStringValue() string {
	return rhinoStringFromID(rhinoStringID(n.rec().payload))
}

// CanResetInPlace reports whether n has no cold extBox, extra side-table entry
// or diagnostic pin, making in-place token/payload mutation exact.
func (n Node) CanResetInPlace() bool {
	return (n.rec().kindFlags & (kindFlagHasExtBox | kindFlagHasExtra)) == 0
}

// PinForDiagnostic marks n so later optimization passes do not recycle n
// in-place, preserving n's original length/token for deferred JSError formatting.
func (n Node) PinForDiagnostic() {
	n.rec().kindFlags |= kindFlagHasExtra
}

// MutateGetPropToNameInPlace rewrites an attached GETPROP node into a NAME node
// in-place (preserving its sibling links, JSType/Color, and Srcref(n)) when n
// has no cold extBox/extra metadata. Interns name before originalName to match
// NodeUtil_NewNameWithOriginalName's exact RhinoStringPool interning order.
func (n Node) MutateGetPropToNameInPlace(name string, originalName string, isConstant bool) bool {
	if !n.CanResetInPlace() {
		return false
	}
	n.DetachChildren()
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	nameID := rhinoStringPoolIntern(name)
	origID := rhinoStringPoolIntern(originalName)
	if (n.props() & propFlagHasJSDoc) != 0 {
		n.setAux(nodeAux{})
	}
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString) | (n.rec().kindFlags & kindFlagTypeIsColor)
	flags := n.props() & propFlag_IS_CONSTANT_NAME
	if isConstant {
		flags |= propFlag_IS_CONSTANT_NAME
	}
	n.setProps(flags)
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(origID)
	return true
}

// MutateCallTargetGetPropToNameInPlace rewrites a CALL's first-child GETPROP
// `o.foo` into `NAME(newMethodName)` in-place and inserts detached `o` right
// after it (`foo(o, ...)`), matching DevirtualizeMethods.RewriteCall.
func (n Node) MutateCallTargetGetPropToNameInPlace(call Node, receiver Node, newMethodName string, isConstantName bool) bool {
	if !n.CanResetInPlace() {
		return false
	}
	n.DetachChildren()
	receiver.SrcrefIfMissing(n)
	nameID := rhinoStringPoolIntern(newMethodName)
	origID := n.node_originalNameID()
	if (n.props() & propFlagHasJSDoc) != 0 {
		n.setAux(nodeAux{})
	}
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString) | (n.rec().kindFlags & kindFlagTypeIsColor)
	flags := uint64(0)
	if isConstantName {
		flags |= propFlag_IS_CONSTANT_NAME
	}
	n.setProps(flags)
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(origID)
	receiver.InsertAfter(n)
	return true
}

// MutateThisToNameIDInPlace rewrites an attached leaf THIS node into
// NAME(nameID) in-place, preserving Srcref(n) and CopyTypeFrom(n).
func (n Node) MutateThisToNameIDInPlace(nameID uint32) bool {
	if !n.CanResetInPlace() || n.rec().first != 0 {
		return false
	}
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	if (n.props() & propFlagHasJSDoc) != 0 {
		n.setAux(nodeAux{})
	}
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString) | (n.rec().kindFlags & kindFlagTypeIsColor)
	n.setProps(0)
	n.rec().payload = uint32(nameID)
	return true
}

// ResetDetachedGetPropToNameID resets a detached GETPROP node (CanResetInPlace)
// into IR_Name(nameID).SrcrefIfMissing(self), optionally setting IS_CONSTANT_NAME.
func (n Node) ResetDetachedGetPropToNameID(nameID uint32, isConstantName bool) {
	origID := n.node_originalNameID()
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString)
	flags := uint64(0)
	if isConstantName {
		flags |= propFlag_IS_CONSTANT_NAME
	}
	n.setProps(flags)
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(origID)
}

// MutateAttachedStmtToVarFromSrcref resets an attached childless statement node
// (CanResetInPlace) into IR_Var(nameChild).SrcrefIfMissing(srcref) in-place.
func (n Node) MutateAttachedStmtToVarFromSrcref(nameChild Node, srcref Node, jsdoc *JSDocInfo) {
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	origID := srcref.node_originalNameID() // before n changes: srcref may be n
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_VAR
	n.rec().kindFlags = 0
	n.setProps(0)
	n.rec().sourceFile = srcref.rec().sourceFile
	n.rec().payload = 0
	n.node_setOriginalNameID(origID)
	n.rec().linenoCharno = srcref.rec().linenoCharno
	n.rec().length16 = srcref.rec().length16
	n.AddChildToFront(nameChild)
	if jsdoc != nil {
		n.SetJSDocInfo(jsdoc)
	}
}

// ResetDetachedToFreshNameID resets a detached childless node (CanResetInPlace)
// into a fresh IR_Name(nameID) with unset source location.
func (n Node) ResetDetachedToFreshNameID(nameID uint32) Node {
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString)
	n.setProps(0)
	n.rec().sourceFile = 0
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(0)
	n.rec().linenoCharno = -1
	n.rec().length16 = 0
	return n
}

// CanCloneLeafInPlace reports whether src is a childless leaf node with no cold
// extBox/extra metadata or overflow source file, allowing zero-allocation
// in-place leaf cloning via OverwriteLeafFromClone.
func (n Node) CanCloneLeafInPlace() bool {
	return n.rec().first == 0 && n.CanResetInPlace()
}

// OverwriteLeafFromClone overwrites an attached childless leaf node (which must
// satisfy CanCloneLeafInPlace) with a clone of src (which must also satisfy
// CanCloneLeafInPlace), matching src.CloneTree() + n.ReplaceWith(clone) with
// zero allocations.
func (n Node) OverwriteLeafFromClone(src Node) {
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	oldSF := n.rec().sourceFile
	oldOrigID := n.node_originalNameID()
	srcOrigID := src.node_originalNameID()
	n.rec().token = src.rec().token
	n.rec().kindFlags = src.rec().kindFlags & (kindMask | kindFlagTypeIsColor)
	n.setTC(n.store().translateTC(src.store(), src.tc(), src.rec().kindFlags&kindFlagTypeIsColor != 0))
	n.setAux(src.auxValue())
	n.setProps(src.props())
	n.rec().payload = src.rec().payload
	if src.rec().sourceFile == 0 {
		n.rec().sourceFile = oldSF
	} else {
		n.rec().sourceFile = src.rec().sourceFile
		n.rec().linenoCharno = src.rec().linenoCharno
		n.rec().length16 = src.rec().length16
	}
	if srcOrigID == 0 {
		srcOrigID = oldOrigID
	}
	n.node_setOriginalNameID(srcOrigID)
}

// ResetTemplateUndefinedSrcref updates a detached 2-node `void 0` template in
// place to match NodeUtil_NewUndefinedNode(src). Returns false if either node in
// the template or src has cold/overflow metadata so the caller can allocate a
// fresh template.
func (n Node) ResetTemplateUndefinedSrcref(src Node) bool {
	num := n.firstNode()
	if num.IsNil() || !n.CanResetInPlace() || !num.CanResetInPlace() {
		return false
	}
	if src.IsNil() {
		n.setTC(0)
		n.setAux(nodeAux{})
		n.rec().kindFlags = 0
		n.setProps(0)
		n.rec().sourceFile = 0
		n.rec().payload = 0
		n.rec().linenoCharno = -1
		n.rec().length16 = 0
		num.setTC(0)
		num.setAux(nodeAux{})
		num.rec().kindFlags = uint8(node_kindNumber)
		num.setProps(0)
		num.rec().sourceFile = 0
		num.rec().payload = 0
		num.rec().linenoCharno = -1
		num.rec().length16 = 0
		n.node_setOriginalNameID(0)
		num.node_setOriginalNameID(0)
		return true
	}
	if !src.CanResetInPlace() || src.node_originalNameID() != 0 {
		return false
	}
	sf := src.rec().sourceFile
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().kindFlags = 0
	n.setProps(0)
	n.rec().sourceFile = sf
	n.rec().payload = 0
	n.rec().linenoCharno = src.rec().linenoCharno
	n.rec().length16 = src.rec().length16
	num.setTC(0)
	num.setAux(nodeAux{})
	num.rec().kindFlags = uint8(node_kindNumber)
	num.setProps(0)
	num.rec().sourceFile = sf
	num.rec().payload = 0
	num.rec().linenoCharno = src.rec().linenoCharno
	num.rec().length16 = src.rec().length16
	n.node_setOriginalNameID(0)
	num.node_setOriginalNameID(0)
	return true
}

// MutateAttachedToFreshVarInPlace resets an attached statement node
// (CanResetInPlace) into a Token_VAR node matching n.ReplaceWith(NewNode(Token_VAR))
// (preserving n's Srcref), with optional single child nameChild and optional JSDocInfo.
func (n Node) MutateAttachedToFreshVarInPlace(nameChild Node, jsdoc *JSDocInfo) {
	n.DetachChildren()
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	origID := n.node_originalNameID()
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_VAR
	n.rec().kindFlags = 0
	n.setProps(0)
	n.rec().payload = 0
	n.node_setOriginalNameID(origID)
	if !nameChild.IsNil() {
		n.AddChildToFront(nameChild)
	}
	if jsdoc != nil {
		n.SetJSDocInfo(jsdoc)
	}
}

// ResetDetachedToNameWithOriginalNameInPlace resets a detached node
// (CanResetInPlace) into a NAME(name) node with Srcref(srcref), originalName,
// and IS_CONSTANT_NAME = isConstant (srcref must satisfy CanResetInPlace).
func (n Node) ResetDetachedToNameWithOriginalNameInPlace(name string, originalName string, srcref Node, isConstant bool) Node {
	n.DetachChildren()
	nameID := rhinoStringPoolIntern(name)
	origID := rhinoStringPoolIntern(originalName)
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_NAME
	n.rec().kindFlags = uint8(node_kindString)
	flags := uint64(0)
	n.rec().sourceFile = srcref.rec().sourceFile
	if isConstant {
		flags |= propFlag_IS_CONSTANT_NAME
	}
	n.setProps(flags)
	n.rec().payload = uint32(nameID)
	n.node_setOriginalNameID(origID)
	n.rec().linenoCharno = srcref.rec().linenoCharno
	n.rec().length16 = srcref.rec().length16
	return n
}

// MutateAttachedToExprResultInPlace resets an attached statement node
// (CanResetInPlace) into IR_ExprResult(exprChild).Srcref(exprChild) + ReplaceWith
// (preserving n's SrcrefIfMissing) in-place (exprChild must be detached and
// satisfy CanResetInPlace).
func (n Node) MutateAttachedToExprResultInPlace(exprChild Node) Node {
	n.DetachChildren()
	if n.rec().parent != 0 {
		n.parentNode().rec().kindFlags &^= kindFlagNoBlockDecl
	}
	origID := exprChild.node_originalNameID()
	if origID == 0 {
		origID = n.node_originalNameID()
	}
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = Token_EXPR_RESULT
	n.rec().kindFlags = 0
	n.setProps(0)
	if exprChild.rec().sourceFile != 0 {
		n.rec().sourceFile = exprChild.rec().sourceFile
		n.rec().linenoCharno = exprChild.rec().linenoCharno
		n.rec().length16 = exprChild.rec().length16
	}
	n.rec().payload = 0
	n.node_setOriginalNameID(origID)
	n.AddChildToFront(exprChild)
	return n
}

// ResetDetachedToBinaryOpFromSrcref resets a detached node (CanResetInPlace)
// into NewNodeWith2Children(op, left, right).Srcref(srcref) in-place (srcref
// must satisfy CanResetInPlace; left and right must be detached).
func (n Node) ResetDetachedToBinaryOpFromSrcref(op Token, left Node, right Node, srcref Node) Node {
	n.DetachChildren()
	origID := srcref.node_originalNameID() // before n changes: srcref may be n
	n.setTC(0)
	n.setAux(nodeAux{})
	n.rec().token = op
	n.rec().kindFlags = 0
	n.setProps(0)
	n.rec().sourceFile = srcref.rec().sourceFile
	n.rec().payload = 0
	n.node_setOriginalNameID(origID)
	n.rec().linenoCharno = srcref.rec().linenoCharno
	n.rec().length16 = srcref.rec().length16
	n.AddChildToFront(left)
	n.AddChildToBack(right)
	return n
}
