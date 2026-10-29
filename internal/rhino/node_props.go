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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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
	"IS_PARENTHESIZED", "NON_JSDOC_COMMENT", "TRAILING_NON_JSDOC_COMMENT",
	"JSDOC_INFO", "INCRDECR", "QUOTED",
	"SYNTHETIC", "ADDED_BLOCK", "SIDE_EFFECT_FLAGS",
	"IS_CONSTANT_NAME", "IS_NAMESPACE", "USE_STRICT",
	"DIRECT_EVAL", "FREE_CALL", "SOURCE_FILE",
	"INPUT_ID", "CHANGE_TIME", "REFLECTED_OBJECT",
	"STATIC_MEMBER", "GENERATOR_FN", "ARROW_FN",
	"ASYNC_FN", "YIELD_ALL", "EXPORT_DEFAULT",
	"EXPORT_ALL_FROM", "CONSTANT_VAR_FLAGS", "IS_GENERATOR_MARKER",
	"IS_GENERATOR_SAFE", "COMPUTED_PROP_METHOD", "COMPUTED_PROP_GETTER",
	"COMPUTED_PROP_SETTER", "COMPUTED_PROP_VARIABLE", "DECLARED_TYPE_EXPR",
	"TYPE_BEFORE_CAST", "COLOR_FROM_CAST", "OPT_ES6_TYPED",
	"GENERIC_TYPE", "IMPLEMENTS", "CONSTRUCT_SIGNATURE",
	"ACCESS_MODIFIER", "NON_INDEXABLE", "GOOG_MODULE",
	"FEATURE_SET", "IS_TYPESCRIPT_ABSTRACT", "DELETED",
	"MODULE_ALIAS", "IS_UNUSED_PARAMETER", "MODULE_EXPORT",
	"IS_SHORTHAND_PROPERTY", "ES6_MODULE", "TYPEDEF_TYPE",
	"START_OF_OPT_CHAIN", "TRAILING_COMMA", "SYNTHESIZED_UNFULFILLED_NAME_DECLARATION",
	"CLOSURE_UNAWARE_SHADOW", "PRIVATE_IDENTIFIER", "IS_IN_CLOSURE_UNAWARE_SUBTREE",
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

func (n Node) lookupPropListHead(prop Node_Prop) *Node_PropListItem {
	propType := int8(prop)
	for x := n.coldPropListHead(); x != nil; x = x.next {
		if x.propType == propType {
			return x
		}
	}
	return nil
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

func (n Node) removeExtraColdProp(prop Node_Prop) {
	if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
		if head := n.coldPropListHead(); head != nil {
			n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
		}
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
			n.removeExtraColdProp(propType)
		} else {
			if (n.props() & mask) == 0 {
				return
			}
			n.andNotProps(mask)
			n.removeExtraColdProp(propType)
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
			n.removeExtraColdProp(prop)
			return
		}
		if value == 1 {
			n.orProps(mask)
			n.removeExtraColdProp(prop)
			return
		}
		n.orProps(mask)
		n.rec().kindFlags |= kindFlagHasExtra
	} else if prop == Node_Prop_CHANGE_TIME && n.nodeKind() == node_kindNode {
		n.rec().payload = uint32(value)
		n.removeExtraColdProp(prop)
		return
	} else if prop == Node_Prop_CONSTANT_VAR_FLAGS {
		if value == 0 {
			n.andNotProps(nodePropConstVarMask)
			n.removeExtraColdProp(prop)
			return
		}
		if value&^3 == 0 {
			n.setProps((n.props() &^ nodePropConstVarMask) | (uint64(value) << nodePropConstVarShift))
			n.removeExtraColdProp(prop)
			return
		}
		n.andNotProps(nodePropConstVarMask)
		n.rec().kindFlags |= kindFlagHasExtra
	} else if prop == Node_Prop_SIDE_EFFECT_FLAGS {
		if value == 0 {
			n.andNotProps(nodePropSideEffectMask)
			n.removeExtraColdProp(prop)
			return
		}
		if value&^15 == 0 {
			n.setProps((n.props() &^ nodePropSideEffectMask) | (uint64(value) << nodePropSideEffectShift))
			n.removeExtraColdProp(prop)
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
		n.removeExtraColdProp(prop)
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
