// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261028

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

	Node_LAST_PROP = Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE
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

// Avoid cloning "values" repeatedly in hot code, we save it off now.
//
// Java: Node.propValues
var Node_propValues = Node_Prop_values()

const (
	// Hot boolean properties (bits 0..24) live directly in nodeData.propFlags.
	propFlag_IS_CONSTANT_NAME uint64 = 1 << iota
	propFlag_QUOTED
	propFlag_FREE_CALL
	propFlag_NON_INDEXABLE
	propFlag_ADDED_BLOCK
	propFlag_TRAILING_COMMA
	propFlag_IS_PARENTHESIZED
	propFlag_ARROW_FN
	propFlag_INCRDECR
	propFlag_SYNTHETIC
	propFlag_IS_SHORTHAND_PROPERTY
	propFlag_IS_NAMESPACE
	propFlag_USE_STRICT
	propFlag_STATIC_MEMBER
	propFlag_GENERATOR_FN
	propFlag_ASYNC_FN
	propFlag_COLOR_FROM_CAST
	propFlag_GOOG_MODULE
	propFlag_ES6_MODULE
	propFlag_DELETED
	propFlag_MODULE_ALIAS
	propFlag_IS_UNUSED_PARAMETER
	propFlag_MODULE_EXPORT
	propFlag_START_OF_OPT_CHAIN
	propFlag_IS_IN_CLOSURE_UNAWARE_SUBTREE

	propHotCount = iota
)

const propColdShift = 32

const (
	// Rare boolean properties (bits 32..46) live in nodeExtBox.coldProps.
	propFlag_DIRECT_EVAL uint64 = 1 << (propColdShift + iota)
	propFlag_YIELD_ALL
	propFlag_EXPORT_DEFAULT
	propFlag_EXPORT_ALL_FROM
	propFlag_IS_GENERATOR_MARKER
	propFlag_IS_GENERATOR_SAFE
	propFlag_COMPUTED_PROP_METHOD
	propFlag_COMPUTED_PROP_GETTER
	propFlag_COMPUTED_PROP_SETTER
	propFlag_COMPUTED_PROP_VARIABLE
	propFlag_OPT_ES6_TYPED
	propFlag_CONSTRUCT_SIGNATURE
	propFlag_IS_TYPESCRIPT_ABSTRACT
	propFlag_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
	propFlag_PRIVATE_IDENTIFIER

	propColdCount = iota
	propColdEnd   = propColdShift + propColdCount
)

const (
	nodePropBoolMask uint64 = ((uint64(1) << propHotCount) - 1) | (((uint64(1) << propColdCount) - 1) << propColdShift)

	// Bits 25..31 complete the inline 32-bit propWord in nodeData.propFlags.
	propFlagHasJSDoc uint64 = 1 << 25

	nodePropConstVarShift          = 26
	nodePropConstVarMask    uint64 = uint64(3) << nodePropConstVarShift
	nodePropSideEffectShift        = 28
	nodePropSideEffectMask  uint64 = uint64(15) << nodePropSideEffectShift
)

// propKind classifies how a Node_Prop is stored on a compact Node record.
type propKind uint8

const (
	propKindCold       propKind = iota // stored exclusively in coldPropListHead
	propKindBool                       // 1-bit flag in n.props(), shadowed in coldPropListHead when kindFlagHasExtra
	propKindPackedInt                  // small unsigned int in n.props(), shadowed in coldPropListHead on overflow/object
	propKindChangeTime                 // uint32 payload on node_kindNode, else coldPropListHead
	propKindSourceFile                 // sourceFile ID in nodeData.sourceFile
	propKindJSDoc                      // inline *JSDocInfo in nodeAux / nodeExtBox.jsDocInfo + propFlagHasJSDoc
	propKindTypeCast                   // type handle in nodeExtBox.typeBeforeCast, else coldPropListHead
)

type propSpec struct {
	mask       uint64
	kind       propKind
	shift      uint8
	maxInline  int8
	staticBase uint8
}

var _ [len(node_Prop_names) - 1 - int(Node_LAST_PROP)]struct{}
var _ [int(Node_LAST_PROP) - (len(node_Prop_names) - 1)]struct{}

var nodePropSpecs = [Node_LAST_PROP + 1]propSpec{
	Node_Prop_IS_PARENTHESIZED:                         {mask: propFlag_IS_PARENTHESIZED, kind: propKindBool},
	Node_Prop_JSDOC_INFO:                               {kind: propKindJSDoc},
	Node_Prop_INCRDECR:                                 {mask: propFlag_INCRDECR, kind: propKindBool},
	Node_Prop_QUOTED:                                   {mask: propFlag_QUOTED, kind: propKindBool},
	Node_Prop_SYNTHETIC:                                {mask: propFlag_SYNTHETIC, kind: propKindBool},
	Node_Prop_ADDED_BLOCK:                              {mask: propFlag_ADDED_BLOCK, kind: propKindBool},
	Node_Prop_SIDE_EFFECT_FLAGS:                        {mask: nodePropSideEffectMask, kind: propKindPackedInt, shift: nodePropSideEffectShift, maxInline: 15, staticBase: 4},
	Node_Prop_IS_CONSTANT_NAME:                         {mask: propFlag_IS_CONSTANT_NAME, kind: propKindBool},
	Node_Prop_IS_NAMESPACE:                             {mask: propFlag_IS_NAMESPACE, kind: propKindBool},
	Node_Prop_USE_STRICT:                               {mask: propFlag_USE_STRICT, kind: propKindBool},
	Node_Prop_DIRECT_EVAL:                              {mask: propFlag_DIRECT_EVAL, kind: propKindBool},
	Node_Prop_FREE_CALL:                                {mask: propFlag_FREE_CALL, kind: propKindBool},
	Node_Prop_SOURCE_FILE:                              {kind: propKindSourceFile},
	Node_Prop_CHANGE_TIME:                              {kind: propKindChangeTime},
	Node_Prop_STATIC_MEMBER:                            {mask: propFlag_STATIC_MEMBER, kind: propKindBool},
	Node_Prop_GENERATOR_FN:                             {mask: propFlag_GENERATOR_FN, kind: propKindBool},
	Node_Prop_ARROW_FN:                                 {mask: propFlag_ARROW_FN, kind: propKindBool},
	Node_Prop_ASYNC_FN:                                 {mask: propFlag_ASYNC_FN, kind: propKindBool},
	Node_Prop_YIELD_ALL:                                {mask: propFlag_YIELD_ALL, kind: propKindBool},
	Node_Prop_EXPORT_DEFAULT:                           {mask: propFlag_EXPORT_DEFAULT, kind: propKindBool},
	Node_Prop_EXPORT_ALL_FROM:                          {mask: propFlag_EXPORT_ALL_FROM, kind: propKindBool},
	Node_Prop_CONSTANT_VAR_FLAGS:                       {mask: nodePropConstVarMask, kind: propKindPackedInt, shift: nodePropConstVarShift, maxInline: 3, staticBase: 0},
	Node_Prop_IS_GENERATOR_MARKER:                      {mask: propFlag_IS_GENERATOR_MARKER, kind: propKindBool},
	Node_Prop_IS_GENERATOR_SAFE:                        {mask: propFlag_IS_GENERATOR_SAFE, kind: propKindBool},
	Node_Prop_COMPUTED_PROP_METHOD:                     {mask: propFlag_COMPUTED_PROP_METHOD, kind: propKindBool},
	Node_Prop_COMPUTED_PROP_GETTER:                     {mask: propFlag_COMPUTED_PROP_GETTER, kind: propKindBool},
	Node_Prop_COMPUTED_PROP_SETTER:                     {mask: propFlag_COMPUTED_PROP_SETTER, kind: propKindBool},
	Node_Prop_COMPUTED_PROP_VARIABLE:                   {mask: propFlag_COMPUTED_PROP_VARIABLE, kind: propKindBool},
	Node_Prop_TYPE_BEFORE_CAST:                         {kind: propKindTypeCast},
	Node_Prop_COLOR_FROM_CAST:                          {mask: propFlag_COLOR_FROM_CAST, kind: propKindBool},
	Node_Prop_OPT_ES6_TYPED:                            {mask: propFlag_OPT_ES6_TYPED, kind: propKindBool},
	Node_Prop_CONSTRUCT_SIGNATURE:                      {mask: propFlag_CONSTRUCT_SIGNATURE, kind: propKindBool},
	Node_Prop_NON_INDEXABLE:                            {mask: propFlag_NON_INDEXABLE, kind: propKindBool},
	Node_Prop_GOOG_MODULE:                              {mask: propFlag_GOOG_MODULE, kind: propKindBool},
	Node_Prop_IS_TYPESCRIPT_ABSTRACT:                   {mask: propFlag_IS_TYPESCRIPT_ABSTRACT, kind: propKindBool},
	Node_Prop_DELETED:                                  {mask: propFlag_DELETED, kind: propKindBool},
	Node_Prop_MODULE_ALIAS:                             {mask: propFlag_MODULE_ALIAS, kind: propKindBool},
	Node_Prop_IS_UNUSED_PARAMETER:                      {mask: propFlag_IS_UNUSED_PARAMETER, kind: propKindBool},
	Node_Prop_MODULE_EXPORT:                            {mask: propFlag_MODULE_EXPORT, kind: propKindBool},
	Node_Prop_IS_SHORTHAND_PROPERTY:                    {mask: propFlag_IS_SHORTHAND_PROPERTY, kind: propKindBool},
	Node_Prop_ES6_MODULE:                               {mask: propFlag_ES6_MODULE, kind: propKindBool},
	Node_Prop_START_OF_OPT_CHAIN:                       {mask: propFlag_START_OF_OPT_CHAIN, kind: propKindBool},
	Node_Prop_TRAILING_COMMA:                           {mask: propFlag_TRAILING_COMMA, kind: propKindBool},
	Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION: {mask: propFlag_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION, kind: propKindBool},
	Node_Prop_PRIVATE_IDENTIFIER:                       {mask: propFlag_PRIVATE_IDENTIFIER, kind: propKindBool},
	Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE:            {mask: propFlag_IS_IN_CLOSURE_UNAWARE_SUBTREE, kind: propKindBool},
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
	nodePropBitMask       [64]uint64
	nodeBitToProp         [propColdEnd]Node_Prop
	staticBoolPropItems   [Node_LAST_PROP + 1]Node_PropListItem
	staticPackedPropItems [20]Node_PropListItem
)

func init() {
	for prop, s := range nodePropSpecs {
		staticBoolPropItems[prop] = Node_PropListItem{propType: int8(prop), intValue: 1, shared: true}
		if s.kind == propKindBool {
			nodePropBitMask[prop] = s.mask
			nodeBitToProp[bits.TrailingZeros64(s.mask)] = Node_Prop(prop)
		} else if s.kind == propKindPackedInt {
			for v := int8(0); v <= s.maxInline; v++ {
				staticPackedPropItems[int(s.staticBase)+int(v)] = Node_PropListItem{
					propType: int8(prop),
					intValue: int32(v),
					shared:   true,
				}
			}
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
	} else if Node_Prop(item.propType) == prop {
		return item.next
	} else {
		result := Node_RebuildListWithoutProp(item.next, prop)
		if result == item.next {
			return item
		}
		return item.Chain(result)
	}
}

func (n Node) inlineJSDoc() *JSDocInfo {
	if (n.rec().kindFlags & kindFlagHasExtBox) == 0 {
		return n.jsDocAux()
	}
	if box := n.extBoxOrNil(); box != nil {
		return box.jsDocInfo
	}
	return nil
}

func (n Node) extraColdPropItem(prop Node_Prop) *Node_PropListItem {
	if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
		return n.lookupPropListHead(prop)
	}
	return nil
}

func (n Node) removeColdProp(prop Node_Prop) {
	if head := n.coldPropListHead(); head != nil {
		n.setColdPropListHead(Node_RebuildListWithoutProp(head, prop))
	}
}

func (n Node) removeExtraColdProp(prop Node_Prop) {
	if (n.rec().kindFlags & kindFlagHasExtra) != 0 {
		n.removeColdProp(prop)
	}
}

func (n Node) putColdObjectProp(prop Node_Prop, value any) {
	head := n.coldPropListHead()
	for cur := head; cur != nil && !cur.shared; cur = cur.next {
		if Node_Prop(cur.propType) == prop && cur.isObject {
			cur.objectValue = value
			return
		}
	}
	head = Node_RebuildListWithoutProp(head, prop)
	n.setColdPropListHead(NewNode_ObjectPropListItem(int8(prop), value, head))
}

func (n Node) putColdIntProp(prop Node_Prop, value int32) {
	head := n.coldPropListHead()
	if head != nil && !head.shared && Node_Prop(head.propType) == prop && !head.isObject {
		head.intValue = value
		return
	}
	head = Node_RebuildListWithoutProp(head, prop)
	n.setColdPropListHead(NewNode_IntPropListItem(int8(prop), value, head))
}

// Java: Node.getProp(Prop)
func (n Node) GetProp(propType Node_Prop) any {
	s := nodePropSpecs[propType]
	switch s.kind {
	case propKindBool:
		if (n.props() & s.mask) == 0 {
			return nil
		}
		if x := n.extraColdPropItem(propType); x != nil {
			return x.GetObjectValue()
		}
		panic(jutil.NewUnsupportedOperationException())
	case propKindPackedInt:
		if (n.props() & s.mask) != 0 {
			panic(jutil.NewUnsupportedOperationException())
		}
	case propKindChangeTime:
		r := n.rec()
		if node_kind(r.kindFlags&kindMask) == node_kindNode && r.payload != 0 {
			panic(jutil.NewUnsupportedOperationException())
		}
	case propKindSourceFile:
		if sfID := n.rec().sourceFile; sfID != 0 {
			return sourceFilePropItem(sfID).objectValue
		}
		return nil
	case propKindJSDoc:
		if (n.rec().propFlags & propWordJSDoc) == 0 {
			return nil
		}
		if jsdoc := n.inlineJSDoc(); jsdoc != nil {
			return jsdoc
		}
	case propKindTypeCast:
		box := n.extBoxOrNil()
		if box == nil {
			return nil
		}
		if tbc := n.store().typeAt(box.typeBeforeCast); tbc != nil {
			return tbc.self()
		}
	}
	if x := n.lookupPropListHead(propType); x != nil {
		return x.GetObjectValue()
	}
	return nil
}

// Java: Node.getBooleanProp(Prop)
func (n Node) GetBooleanProp(propType Node_Prop) bool {
	if mask := nodePropBitMask[propType]; mask != 0 {
		if w := propWord(mask); w != 0 {
			return n.rec().propFlags&w != 0
		}
		return (n.props() & mask) != 0
	}
	return n.GetIntProp(propType) != 0
}

// GetIntProp returns the integer value for the property, or 0 if the
// property is not defined.
//
// Java: Node.getIntProp(Prop)
func (n Node) GetIntProp(propType Node_Prop) int32 {
	s := nodePropSpecs[propType]
	switch s.kind {
	case propKindBool:
		if (n.props() & s.mask) == 0 {
			return 0
		}
		if x := n.extraColdPropItem(propType); x != nil {
			return x.GetIntValue()
		}
		return 1
	case propKindPackedInt:
		r := n.rec()
		if v := int32((uint64(r.propFlags) & s.mask) >> s.shift); v != 0 {
			return v
		}
		if (r.kindFlags & kindFlagHasExtra) == 0 {
			return 0
		}
	case propKindChangeTime:
		r := n.rec()
		if node_kind(r.kindFlags&kindMask) == node_kindNode {
			if v := int32(r.payload); v != 0 {
				return v
			}
			if (r.kindFlags & kindFlagHasExtra) == 0 {
				return 0
			}
		}
	case propKindSourceFile:
		if n.rec().sourceFile != 0 {
			panic(jutil.NewUnsupportedOperationException())
		}
		return 0
	case propKindJSDoc:
		if (n.rec().propFlags & propWordJSDoc) == 0 {
			return 0
		}
		if n.inlineJSDoc() != nil {
			panic(jutil.NewUnsupportedOperationException())
		}
	case propKindTypeCast:
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
	s := nodePropSpecs[prop]
	switch s.kind {
	case propKindBool:
		n.orProps(s.mask)
		n.rec().kindFlags |= kindFlagHasExtra
	case propKindPackedInt:
		r := n.rec()
		r.propFlags &^= propWord(s.mask)
		r.kindFlags |= kindFlagHasExtra
	case propKindChangeTime:
		r := n.rec()
		if node_kind(r.kindFlags&kindMask) == node_kindNode {
			r.payload = 0
			r.kindFlags |= kindFlagHasExtra
		}
	case propKindSourceFile:
		n.rec().sourceFile = allocSourceFileID(value)
		return
	case propKindJSDoc:
		if info, ok := value.(*JSDocInfo); ok {
			n.SetJSDocInfo(info)
			return
		}
		n.ensureExtBox().jsDocInfo = nil
		n.rec().propFlags |= propWordJSDoc
	case propKindTypeCast:
		if t, ok := value.(JSType); ok {
			n.SetJSTypeBeforeCast(t)
			return
		}
		n.ensureExtBox().typeBeforeCast = 0
	}
	n.putColdObjectProp(prop, value)
}

// Java: Node.putBooleanProp(Prop, boolean)
func (n Node) PutBooleanProp(propType Node_Prop, value bool) {
	if mask := nodePropBitMask[propType]; mask != 0 {
		if w := propWord(mask); w != 0 {
			r := n.rec()
			if value {
				r.propFlags |= w
			} else {
				if r.propFlags&w == 0 {
					return
				}
				r.propFlags &^= w
			}
			if r.kindFlags&kindFlagHasExtra != 0 {
				n.removeColdProp(propType)
			}
			return
		}
		if value {
			n.orProps(mask)
		} else {
			if (n.props() & mask) == 0 {
				return
			}
			n.andNotProps(mask)
		}
		n.removeExtraColdProp(propType)
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
	if value == 0 {
		n.RemoveProp(prop)
		return
	}
	s := nodePropSpecs[prop]
	switch s.kind {
	case propKindBool:
		n.orProps(s.mask)
		if value == 1 {
			n.removeExtraColdProp(prop)
			return
		}
		n.rec().kindFlags |= kindFlagHasExtra
	case propKindPackedInt:
		r := n.rec()
		if value&^int32(s.maxInline) == 0 {
			r.propFlags = (r.propFlags &^ propWord(s.mask)) | propWord(uint32(value)<<s.shift)
			if r.kindFlags&kindFlagHasExtra != 0 {
				n.removeColdProp(prop)
			}
			return
		}
		r.propFlags &^= propWord(s.mask)
		r.kindFlags |= kindFlagHasExtra
	case propKindChangeTime:
		r := n.rec()
		if node_kind(r.kindFlags&kindMask) == node_kindNode {
			r.payload = uint32(value)
			if r.kindFlags&kindFlagHasExtra != 0 {
				n.removeColdProp(prop)
			}
			return
		}
	case propKindSourceFile:
		// Go-only: the source file is an object property (an id in the record).
		panic(jutil.NewUnsupportedOperationException())
	case propKindJSDoc:
		n.ensureExtBox().jsDocInfo = nil
		n.rec().propFlags |= propWordJSDoc
	case propKindTypeCast:
		n.ensureExtBox().typeBeforeCast = 0
	}
	n.putColdIntProp(prop, value)
}

// Java: Node.removeProp(Prop)
func (n Node) RemoveProp(prop Node_Prop) {
	s := nodePropSpecs[prop]
	switch s.kind {
	case propKindBool:
		if (n.props() & s.mask) == 0 {
			return
		}
		n.andNotProps(s.mask)
		n.removeExtraColdProp(prop)
		return
	case propKindPackedInt:
		r := n.rec()
		r.propFlags &^= propWord(s.mask)
		if r.kindFlags&kindFlagHasExtra != 0 {
			n.removeColdProp(prop)
		}
		return
	case propKindChangeTime:
		r := n.rec()
		if node_kind(r.kindFlags&kindMask) == node_kindNode {
			r.payload = 0
			if r.kindFlags&kindFlagHasExtra != 0 {
				n.removeColdProp(prop)
			}
			return
		}
	case propKindSourceFile:
		n.rec().sourceFile = 0
		return
	case propKindJSDoc:
		r := n.rec()
		if (r.propFlags & propWordJSDoc) == 0 {
			return
		}
		r.propFlags &^= propWordJSDoc
		if (r.kindFlags & kindFlagHasExtBox) == 0 {
			n.setAux(nodeAux{})
			return
		}
		box := n.extBox()
		if box.jsDocInfo != nil {
			box.jsDocInfo = nil
			n.maybeCollapseExtBox(box)
			return
		}
	case propKindTypeCast:
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
	n.removeColdProp(prop)
}
