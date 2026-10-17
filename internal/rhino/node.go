// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261015

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
	"slices"
	"strconv"
	"strings"
	"unsafe"

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
func (n *Node) GetNonJSDocCommentString() string {
	if n.GetProp(Node_Prop_NON_JSDOC_COMMENT) == nil {
		return ""
	}
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_NON_JSDOC_COMMENT)).GetCommentString()
}

// Java: Node.getNonJSDocComment()
func (n *Node) GetNonJSDocComment() *NonJSDocComment {
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_NON_JSDOC_COMMENT))
}

// Java: Node.getTrailingNonJSDocCommentString()
func (n *Node) GetTrailingNonJSDocCommentString() string {
	if n.GetProp(Node_Prop_TRAILING_NON_JSDOC_COMMENT) == nil {
		return ""
	}
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_TRAILING_NON_JSDOC_COMMENT)).GetCommentString()
}

// Java: Node.getTrailingNonJSDocComment()
func (n *Node) GetTrailingNonJSDocComment() *NonJSDocComment {
	return node_cast[*NonJSDocComment](n.GetProp(Node_Prop_TRAILING_NON_JSDOC_COMMENT))
}

// SetNonJSDocComment sets the NonJSDoc comment attached to this node.
//
// Java: Node.setNonJSDocComment(NonJSDocComment)
func (n *Node) SetNonJSDocComment(comment *NonJSDocComment) *Node {
	n.PutProp(Node_Prop_NON_JSDOC_COMMENT, comment)
	return n
}

// Java: Node.setTrailingNonJSDocComment(NonJSDocComment)
func (n *Node) SetTrailingNonJSDocComment(comment *NonJSDocComment) *Node {
	n.PutProp(Node_Prop_TRAILING_NON_JSDOC_COMMENT, comment)
	return n
}

// SetIsParenthesized sets whether this node was inside original source-level
// parentheses.
//
// Java: Node.setIsParenthesized(boolean)
func (n *Node) SetIsParenthesized(b bool) {
	jutil.CheckState(IR_MayBeExpression(n))
	n.PutBooleanProp(Node_Prop_IS_PARENTHESIZED, b)
}

// GetIsParenthesized checks whether node was inside original source-level
// parentheses.
//
// Java: Node.getIsParenthesized()
func (n *Node) GetIsParenthesized() bool {
	return n.GetBooleanProp(Node_Prop_IS_PARENTHESIZED)
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
	node_kindNumber
	node_kindBigInt
	node_kindString
	node_kindTemplateLiteralSubstring
)

var node_kindClassNames = [...]string{
	"com.google.javascript.rhino.Node",
	"com.google.javascript.rhino.Node$NumberNode",
	"com.google.javascript.rhino.Node$BigIntNode",
	"com.google.javascript.rhino.Node$StringNode",
	"com.google.javascript.rhino.Node$TemplateLiteralSubstringNode",
}

// GetClassName is `getClass().getName()`: the name of the Java class this
// node is an instance of (Node or one of its folded private subclasses).
//
// Go-only: language-adapter — Java code compares `node.getClass()` with
// `Node.class` (CodeGenerator); the folded subclasses have no Go type of
// their own.
func (n *Node) GetClassName() string {
	return node_kindClassNames[n.kind]
}

// node_castKind is the Java cast `((XNode) this)`: it throws a
// ClassCastException unless n is of the given kind.
func (n *Node) node_castKind(k node_kind) {
	if n.kind != k {
		panic(jutil.NewClassCastException(node_kindClassNames[n.kind], node_kindClassNames[k]))
	}
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
		panic(jutil.NewClassCastException(fmt.Sprintf("%T", value), fmt.Sprintf("%T", any(zero))))
	}
	return t
}

// NumberNode(double)
//
// Java: Node.NumberNode(double)
func newNode_NumberNode(number float64) *Node {
	n := NewNode(Token_NUMBER)
	n.kind = node_kindNumber
	n.SetDouble(number)
	return n
}

// BigIntNode(BigInteger)
//
// Java: Node.BigIntNode(BigInteger)
func newNode_BigIntNode(bigint *big.Int) *Node {
	n := NewNode(Token_BIGINT)
	n.kind = node_kindBigInt
	n.SetBigInt(bigint)
	return n
}

// Only for cloneNode.
//
// Java: Node.StringNode(Token)
func newNode_StringNodeForClone(token Token) *Node {
	n := NewNode(token)
	n.kind = node_kindString
	return n
}

// Java: Node.StringNode(Token, String)
func newNode_StringNode(token Token, str string) *Node {
	n := NewNode(token)
	n.kind = node_kindString
	n.SetString(str)
	return n
}

// Java: Node.StringNode(Token, RhinoStringPool.LazyInternedStringList, int)
func newNode_StringNodeFromStringPool(token Token, stringPool *RhinoStringPool_LazyInternedStringList, offset int) *Node {
	n := NewNode(token)
	n.kind = node_kindString
	n.SetStringFromStringPool(stringPool, offset)
	return n
}

// Java: Node.TemplateLiteralSubstringNode(String, String)
func newNode_TemplateLiteralSubstringNode(cooked *string, raw string) *Node {
	n := NewNode(Token_TEMPLATELIT_STRING)
	n.kind = node_kindTemplateLiteralSubstring
	// RhinoStringPool is null-hostile.
	if cooked == nil {
		// n.cooked = nil: it is.
	} else {
		c := RhinoStringPool_AddOrGet(*cooked)
		n.node_getExt().cooked = &c
	}
	n.str = RhinoStringPool_AddOrGet(raw) // raw
	return n
}

// Java: Node.TemplateLiteralSubstringNode(RhinoStringPool.LazyInternedStringList, int, int)
func newNode_TemplateLiteralSubstringNodeFromStringPool(
	stringPool *RhinoStringPool_LazyInternedStringList,
	cookedOffsetOrNegativeOne int,
	rawOffset int,
) *Node {
	n := NewNode(Token_TEMPLATELIT_STRING)
	n.kind = node_kindTemplateLiteralSubstring
	if cookedOffsetOrNegativeOne == -1 {
		// n.cooked = nil: it is.
	} else {
		c := stringPool.Get(cookedOffsetOrNegativeOne)
		n.node_getExt().cooked = &c
	}
	n.str = stringPool.Get(rawOffset) // raw
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
}

// Java: Node.ObjectPropListItem(byte, Object, PropListItem)
func NewNode_ObjectPropListItem(propType int8, objectValue any, next *Node_PropListItem) *Node_PropListItem {
	return &Node_PropListItem{
		propType:    propType,
		next:        next,
		isObject:    true,
		objectValue: jutil.CheckNotNull(objectValue),
	}
}

// Java: Node.IntPropListItem(byte, int, PropListItem)
func NewNode_IntPropListItem(propType int8, intValue int32, next *Node_PropListItem) *Node_PropListItem {
	item := &Node_PropListItem{propType: propType, next: next, intValue: intValue}
	jutil.CheckState(item.intValue != 0)
	return item
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
func NewNode(token Token) *Node {
	return &Node{token: token, linenoCharno: -1}
}

// Java: Node(Token, Node)
func NewNodeWithChild(token Token, child *Node) *Node {
	n := NewNode(token)
	n.first = child

	child.CheckDetached()
	// child.next remains null;
	child.previous = child
	child.parent = n
	return n
}

// Java: Node(Token, Node, Node)
func NewNodeWith2Children(token Token, left *Node, right *Node) *Node {
	n := NewNode(token)
	n.first = left

	left.CheckDetached()
	left.next = right
	left.previous = right
	left.parent = n

	right.CheckDetached()
	// right.next remains null;
	right.previous = left
	right.parent = n
	return n
}

// Java: Node(Token, Node, Node, Node)
func NewNodeWith3Children(token Token, left *Node, mid *Node, right *Node) *Node {
	n := NewNode(token)
	n.first = left

	left.CheckDetached()
	left.next = mid
	left.previous = right
	left.parent = n

	mid.CheckDetached()
	mid.next = right
	mid.previous = left
	mid.parent = n

	right.CheckDetached()
	// right.next remains null;
	right.previous = mid
	right.parent = n
	return n
}

// Java: Node.newNumber(double)
func Node_NewNumber(number float64) *Node {
	return newNode_NumberNode(number)
}

// Java: Node.newBigInt(BigInteger)
func Node_NewBigInt(bigint *big.Int) *Node {
	return newNode_BigIntNode(bigint)
}

// Java: Node.newString(String)
func Node_NewString(str string) *Node {
	return newNode_StringNode(Token_STRINGLIT, str)
}

// Java: Node.newString(Token, String)
func Node_NewStringWithToken(token Token, str string) *Node {
	return newNode_StringNode(token, str)
}

// Java: Node.newString(Token, RhinoStringPool.LazyInternedStringList, int)
func Node_NewStringFromStringPool(token Token, stringPool *RhinoStringPool_LazyInternedStringList, offset int) *Node {
	return newNode_StringNodeFromStringPool(token, stringPool, offset)
}

// Java: Node.newTemplateLitString(String, String)
func Node_NewTemplateLitString(cooked *string, raw string) *Node {
	return newNode_TemplateLiteralSubstringNode(cooked, raw)
}

// Java: Node.newTemplateLitString(RhinoStringPool.LazyInternedStringList, int, int)
func Node_NewTemplateLitStringFromStringPool(
	stringPool *RhinoStringPool_LazyInternedStringList,
	cookedOffsetOrNegativeOne int,
	rawOffset int,
) *Node {
	return newNode_TemplateLiteralSubstringNodeFromStringPool(stringPool, cookedOffsetOrNegativeOne, rawOffset)
}

// Java: Node.getToken()
func (n *Node) GetToken() Token {
	return n.token
}

// Java: Node.setToken(Token)
func (n *Node) SetToken(token Token) {
	n.token = token
}

// Java: Node.hasChildren()
func (n *Node) HasChildren() bool {
	return n.first != nil
}

// Java: Node.getOnlyChild()
func (n *Node) GetOnlyChild() *Node {
	jutil.CheckState(n.HasOneChild())
	return n.first
}

// Java: Node.getFirstChild()
func (n *Node) GetFirstChild() *Node {
	return n.first
}

// GetFirstFirstChild gets the first child of the first child. This method
// assumes that the first child exists.
//
// Java: Node.getFirstFirstChild()
func (n *Node) GetFirstFirstChild() *Node {
	return n.first.first
}

// Java: Node.getSecondChild()
func (n *Node) GetSecondChild() *Node {
	return n.first.next
}

// Java: Node.getLastChild()
func (n *Node) GetLastChild() *Node {
	if n.first != nil {
		return n.first.previous
	}
	return nil
}

// Java: Node.getNext()
func (n *Node) GetNext() *Node {
	return n.next
}

// Java: Node.getPrevious()
func (n *Node) GetPrevious() *Node {
	if n == n.parent.first {
		return nil
	}
	return n.previous
}

// Java: Node.setClosureUnawareShadow(Node)
func (n *Node) SetClosureUnawareShadow(shadowRoot *Node) {
	jutil.CheckStateMsg(n.first == nil, "Cannot set shadow root on a node with children")
	jutil.CheckStateMsg(n.token == Token_NAME, "Only NAME nodes can be used as shadows")
	n.PutProp(Node_Prop_CLOSURE_UNAWARE_SHADOW, shadowRoot)
}

// Java: Node.getClosureUnawareShadow()
func (n *Node) GetClosureUnawareShadow() *Node {
	if n.token != Token_NAME {
		// As checked in setClosureUnawareShadow, only NAME nodes can have shadow
		// content. We short-circuit the getProp call here because that can be
		// very expensive for nodes that have long proplists (and we know it
		// wouldn't find this prop anyways).
		return nil
	}
	return node_cast[*Node](n.GetProp(Node_Prop_CLOSURE_UNAWARE_SHADOW))
}

// Java: Node.setIsInClosureUnawareSubtree(boolean)
func (n *Node) SetIsInClosureUnawareSubtree(value bool) {
	if n.propListHead == nil ||
		int(n.propListHead.propType) != Node_Prop_SOURCE_FILE.Ordinal() ||
		n.propListHead.next != nil {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.StringFormat(
				"Cannot set IS_IN_CLOSURE_UNAWARE_SUBTREE property on a node which proplist is not a"+
					" single SOURCE_FILE prop: %s %s",
				n, n.GetPropListDebugString())))
	}
	n.PutBooleanProp(Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE, value)
}

// Java: Node.getIsInClosureUnawareSubtree()
func (n *Node) GetIsInClosureUnawareSubtree() bool {
	return n.GetBooleanProp(Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE)
}

// GetChildAtIndex gets the ith child, note that this is O(N) where N is the
// number of children.
//
// Java: Node.getChildAtIndex(int)
func (n *Node) GetChildAtIndex(i int) *Node {
	c := n.first
	for i > 0 {
		c = c.next
		i--
	}
	return c
}

// GetIndexOfChild gets the index of a child, note that this is O(N) where N
// is the number of children.
//
// Java: Node.getIndexOfChild(Node)
func (n *Node) GetIndexOfChild(child *Node) int {
	c := n.first
	i := 0
	for c != nil {
		if child == c {
			return i
		}

		c = c.next
		i++
	}
	return -1
}

// Java: Node.addChildToFront(Node)
func (n *Node) AddChildToFront(child *Node) {
	jutil.CheckArgument(child.parent == nil)
	jutil.CheckArgument(child.next == nil)
	jutil.CheckArgument(child.previous == nil)
	child.parent = n
	child.next = n.first
	if n.first == nil {
		// NOTE: child.next remains null
		child.previous = child
	} else {
		last := n.first.previous
		// NOTE: last.next remains null
		child.previous = last
		child.next = n.first
		n.first.previous = child
	}
	n.first = child
}

// Java: Node.addChildToBack(Node)
func (n *Node) AddChildToBack(child *Node) {
	jutil.CheckArgumentMsg3(
		child.parent == nil,
		"Cannot add already-owned child node.\nChild: %s\nExisting parent: %s\nNew parent: %s",
		child,
		child.parent,
		n)

	jutil.CheckArgument(child.next == nil)
	jutil.CheckArgument(child.previous == nil)

	if n.first == nil {
		// NOTE: child.next remains null
		child.previous = child
		n.first = child
	} else {
		last := n.first.previous
		last.next = child
		// NOTE: child.next remains null
		child.previous = last
		n.first.previous = child
	}

	child.parent = n
}

// AddChildrenToFront adds all children to the front of this node.
//
// children is the first of a list of sibling nodes who have no parent.
// NOTE: Usually you would get this argument from a removeChildren() call. A
// single detached node will not work because its sibling pointers will not
// be correctly initialized.
//
// Java: Node.addChildrenToFront(Node)
func (n *Node) AddChildrenToFront(children *Node) {
	if children == nil {
		return // removeChildren() returns null when there are none
	}
	// NOTE: If there is only one sibling, its previous pointer must point to
	// itself. Null indicates a fully detached node.
	jutil.CheckNotNullMsg1(children.previous, "%s", children)
	for child := children; child != nil; child = child.next {
		jutil.CheckArgument(child.parent == nil)
		child.parent = n
	}

	lastSib := children.previous
	if n.first != nil {
		last := n.first.previous
		// NOTE: last.next remains null
		children.previous = last
		lastSib.next = n.first
		n.first.previous = lastSib
	}
	n.first = children
}

// Java: Node.addChildrenToBack(Node)
func (n *Node) AddChildrenToBack(children *Node) {
	n.AddChildrenAfter(children, n.GetLastChild())
}

// Java: Node.insertAfter(Node)
func (n *Node) InsertAfter(existing *Node) {
	existing.CheckAttached()
	n.CheckDetached()

	existingParent := existing.parent
	existingNext := existing.next

	n.parent = existingParent

	existing.next = n
	n.previous = existing

	if existingNext == nil {
		existingParent.first.previous = n
		// this.next remains null
	} else {
		existingNext.previous = n
		n.next = existingNext
	}
}

// Java: Node.insertBefore(Node)
func (n *Node) InsertBefore(existing *Node) {
	existing.CheckAttached()
	n.CheckDetached()

	existingParent := existing.parent
	existingPrevious := existing.previous

	n.parent = existingParent

	n.next = existing
	existing.previous = n

	n.previous = existingPrevious
	if existingPrevious.next == nil {
		existingParent.first = n
		// existingPrevious.next remains null
	} else {
		// existingParent.first remains existing
		existingPrevious.next = n
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
func (n *Node) AddChildrenAfter(children *Node, node *Node) {
	if children == nil {
		return // removeChildren() returns null when there are none
	}
	jutil.CheckArgument(node == nil || node.parent == n)
	// NOTE: If there is only one sibling, its previous pointer must point to
	// itself. Null indicates a fully detached node.
	jutil.CheckNotNullMsg1(children.previous, "%s", children)
	if node == nil {
		n.AddChildrenToFront(children)
		return
	}

	for child := children; child != nil; child = child.next {
		jutil.CheckArgument(child.parent == nil)
		child.parent = n
	}

	lastSibling := children.previous
	nodeAfter := node.next
	lastSibling.next = nodeAfter
	if nodeAfter == nil {
		n.first.previous = lastSibling
	} else {
		nodeAfter.previous = lastSibling
	}
	node.next = children
	children.previous = node
}

// ReplaceWith swaps `replacement` and its subtree into the position of
// `this`.
//
// Java: Node.replaceWith(Node)
func (n *Node) ReplaceWith(replacement *Node) {
	n.CheckAttached()
	replacement.CheckDetached()

	existingParent := n.parent
	existingNext := n.next
	existingPrevious := n.previous

	// Copy over important information.
	// TODO(nickreid): Stop doing this. It's totally unexpected.
	replacement.SrcrefIfMissing(n)

	// The sequence below also has to work when `this` is an only child, ehich
	// can cause many of the variables to point to the same object.

	n.parent = nil
	replacement.parent = existingParent

	n.previous = nil
	replacement.previous = existingPrevious
	if existingPrevious.next == nil {
		existingParent.first = replacement
		// existingPrevious.next remains null
	} else {
		// existingParent.first is unchanged;
		existingPrevious.next = replacement
	}

	if existingNext == nil {
		// this.next remains null;
		existingParent.first.previous = replacement
		// replacement.next remains null
	} else {
		n.next = nil
		existingNext.previous = replacement
		replacement.next = existingNext
	}
}

// Detach removes this node from its parent, but retains its subtree.
//
// Java: Node.detach()
func (n *Node) Detach() *Node {
	n.CheckAttached()

	existingParent := n.parent
	existingNext := n.next
	existingPrevious := n.previous

	// The sequence below also has to work when `this` is an only child or has
	// a single sibling, which can cause many of the variables to point to the
	// same object.

	n.parent = nil

	if existingNext == nil {
		// this.next remains null;
		existingParent.first.previous = existingPrevious
	} else {
		n.next = nil
		existingNext.previous = existingPrevious
	}

	n.previous = nil
	if existingPrevious.next == nil {
		existingParent.first = existingNext
		// existingPrevious.next remains null
	} else {
		// existingParent.first is unchanged;
		existingPrevious.next = existingNext
	}

	return n
}

// Java: Node.CheckAttached()
func (n *Node) CheckAttached() {
	jutil.CheckStateMsg1(n.parent != nil, "Has no parent: %s", n)
}

// Java: Node.CheckDetached()
func (n *Node) CheckDetached() {
	jutil.CheckStateMsg1(n.parent == nil, "Has parent: %s", n)
	jutil.CheckStateMsg1(n.next == nil, "Has next: %s", n)
	jutil.CheckStateMsg1(n.previous == nil, "Has previous: %s", n)
}

// RemoveFirstChild removes the first child of Node. Equivalent to:
// node.removeChild(node.getFirstChild());
//
// Java: Node.removeFirstChild()
func (n *Node) RemoveFirstChild() *Node {
	child := n.first
	if child != nil {
		child.Detach()
	}
	return child
}

// RemoveChildren removes all children, but leaves them linked to each other.
// Returns the first child node.
//
// Java: Node.removeChildren()
func (n *Node) RemoveChildren() *Node {
	children := n.first
	for child := n.first; child != nil; child = child.next {
		child.parent = nil
	}
	n.first = nil
	return children
}

// DetachChildren removes all children from this node and isolates the
// children from each other.
//
// Java: Node.detachChildren()
func (n *Node) DetachChildren() {
	for child := n.first; child != nil; {
		nextChild := child.next
		child.parent = nil
		child.next = nil
		child.previous = nil
		child = nextChild
	}
	n.first = nil
}

// Java: Node.lookupProperty(Prop)
func (n *Node) LookupProperty(prop Node_Prop) *Node_PropListItem {
	propType := int8(prop.Ordinal())
	x := n.propListHead
	for x != nil && propType != x.propType {
		x = x.next
	}
	return x
}

// ClonePropsFrom clones the properties from the provided node without
// copying the property object. The receiving node may not have any existing
// properties.
//
// Java: Node.clonePropsFrom(Node)
func (n *Node) ClonePropsFrom(other *Node) *Node {
	jutil.CheckStateMsg(n.propListHead == nil, "Node has existing properties.")
	n.propListHead = other.propListHead
	return n
}

// ValidateProperties checks for invalid or missing properties and feeds
// error messages for any violations to the given `Consumer`.
//
// We use a `Consumer` to avoid the cost of building a usually-empty list
// every time this method is called.
//
// Java: Node.validateProperties(Consumer<String>)
func (n *Node) ValidateProperties(violationMessageConsumer func(string)) {
	if n.propListHead == nil {
		// TODO(bradfordcsmith): Fix the bugs that prevent enabling this
		// validation.
		//
		// In particular:
		//
		// 1. CoverageInstrumentationPass and
		// BranchCoverageInstrumentationCallback create a bunch of nodes without
		// valid source reference info, and it isn't obvious what source
		// reference info they should be using.
		//
		// 2. b/186056977 covers an issue blocking correcting a violation in
		// `VarCheck`
		//
		// 3. Fixing a violation in DeclaredGlobalExternsOnWindow makes an
		// unexpected change to the pre-computed TypedAst for the runtime
		// libraries.
		//
		// if (token != Token.ROOT) {
		//   violationMessageConsumer.accept("non-ROOT has no properties");
		// }
		return
	}

	if n.token == Token_ROOT {
		// ROOT tokens should never have properties
		violationMessageConsumer("ROOT has properties")
	}

	for propListItem := n.propListHead; propListItem != nil; propListItem = propListItem.next {
		prop := Node_propValues[propListItem.propType]
		// Catch it if the definition of Prop ever changes so that the ordinals
		// don't line up.
		jutil.CheckStateMsg1(prop.Ordinal() == int(propListItem.propType), "ordinal doesn't match: %s", prop)

		// TODO(bradfordcsmith): This is not yet an exhaustive list of
		// validations. Other validations should be added as it is found useful
		// to have them. Some property validation is done independently in
		// `AstValidator` and could possibly be moved here.
		//
		// This method was added in response to a bug that created an invalid
		// IS_PARENTHESIZED property that was discovered by a check in
		// `deserializeProperties()`, so initially this method was created to
		// cover the checks previously done there.
		switch prop {
		case Node_Prop_IS_PARENTHESIZED:
			if !IR_MayBeExpression(n) {
				violationMessageConsumer("non-expression is parenthesized")
			}
		case Node_Prop_ARROW_FN:
			if !n.IsFunction() {
				violationMessageConsumer("invalid ARROW_FN prop")
			}
		case Node_Prop_ASYNC_FN:
			if !n.IsFunction() {
				violationMessageConsumer("invalid ASYNC_FN prop")
			}
		case Node_Prop_SYNTHETIC:
			if !n.IsBlock() {
				violationMessageConsumer("invalid SYNTHETIC prop")
			}
		case Node_Prop_COLOR_FROM_CAST:
			if n.GetColor() == nil {
				violationMessageConsumer("COLOR_FROM_CAST with no Color")
			}
		case Node_Prop_START_OF_OPT_CHAIN:
			if !(n.IsOptChainCall() || n.IsOptChainGetElem() || n.IsOptChainGetProp()) {
				violationMessageConsumer("START_OF_OPT_CHAIN on non-optional Node")
			}
		case Node_Prop_CONSTANT_VAR_FLAGS:
			if !(n.IsName() || n.IsImportStar()) {
				violationMessageConsumer("invalid CONST_VAR_FLAGS")
			}
		case Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION:
			// note: we could relax this restriction if VarCheck needed to
			// generate other forms of synthetic externs
			if !n.IsVar() {
				// TODO: this check used to also try and validate that the
				// synthetic externs VAR's had names, but this doesn't work during
				// deserialization because the name is a child node that we
				// haven't even seen yet when we are validating properties.
				// || !hasOneChild() || !getFirstChild().isName()) {
				violationMessageConsumer(
					"Expected all synthetic unfulfilled declarations to be `var <name>`")
			}
		case Node_Prop_CLOSURE_UNAWARE_SHADOW:
			shadowProp := n.LookupProperty(Node_Prop_CLOSURE_UNAWARE_SHADOW)
			if shadowProp == nil || !shadowProp.isObject || !node_isNode(shadowProp.GetObjectValue()) {
				violationMessageConsumer("CLOSURE_UNAWARE_SHADOW property must point to a Node")
			}
		default:
			// No validation is currently done for other properties
		}
	}
}

func node_isNode(v any) bool {
	_, ok := v.(*Node)
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
func (n *Node) GetProp(propType Node_Prop) any {
	item := n.LookupProperty(propType)
	if item == nil {
		return nil
	}
	return item.GetObjectValue()
}

// Java: Node.getBooleanProp(Prop)
func (n *Node) GetBooleanProp(propType Node_Prop) bool {
	return n.GetIntProp(propType) != 0
}

// GetIntProp returns the integer value for the property, or 0 if the
// property is not defined.
//
// Java: Node.getIntProp(Prop)
func (n *Node) GetIntProp(propType Node_Prop) int32 {
	item := n.LookupProperty(propType)
	if item == nil {
		return 0
	}
	return item.GetIntValue()
}

// PutProp sets an object-valued property; a null value removes it.
//
// Go-only: language-adapter — jutil.IsNil treats a typed nil pointer passed
// through `any` as Java null.
//
// Java: Node.putProp(Prop, Object)
func (n *Node) PutProp(prop Node_Prop, value any) {
	n.propListHead = Node_RebuildListWithoutProp(n.propListHead, prop)
	if !jutil.IsNil(value) {
		n.propListHead = NewNode_ObjectPropListItem(int8(prop.Ordinal()), value, n.propListHead)
	}
}

// Java: Node.putBooleanProp(Prop, boolean)
func (n *Node) PutBooleanProp(propType Node_Prop, value bool) {
	if value {
		n.PutIntProp(propType, 1)
	} else {
		n.PutIntProp(propType, 0)
	}
}

// Java: Node.putIntProp(Prop, int)
func (n *Node) PutIntProp(prop Node_Prop, value int32) {
	n.propListHead = Node_RebuildListWithoutProp(n.propListHead, prop)
	if value != 0 {
		n.propListHead = NewNode_IntPropListItem(int8(prop.Ordinal()), value, n.propListHead)
	}
}

// Java: Node.removeProp(Prop)
func (n *Node) RemoveProp(prop Node_Prop) {
	n.propListHead = Node_RebuildListWithoutProp(n.propListHead, prop)
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

// Java: Node.serializeProperties()
func (n *Node) SerializeProperties() int64 {
	var propSet int64 = 0
	for propListItem := n.propListHead; propListItem != nil; propListItem = propListItem.next {
		prop := Node_propValues[propListItem.propType]

		switch prop {
		case Node_Prop_TYPE_BEFORE_CAST:
			propSet = Node_SetNodePropertyBit(propSet, NodeProperty_COLOR_FROM_CAST)
		case Node_Prop_CONSTANT_VAR_FLAGS:
			intVal := propListItem.GetIntValue()
			if Node_AnyBitSet(intVal, Node_ConstantVarFlags_INFERRED) {
				propSet = Node_SetNodePropertyBit(propSet, NodeProperty_IS_INFERRED_CONSTANT)
			}
			if Node_AnyBitSet(intVal, Node_ConstantVarFlags_DECLARED) {
				propSet = Node_SetNodePropertyBit(propSet, NodeProperty_IS_DECLARED_CONSTANT)
			}
		case Node_Prop_SIDE_EFFECT_FLAGS:
			propSet = n.SetNodePropertySideEffectFlags(propSet, propListItem.GetIntValue())
		case Node_Prop_CLOSURE_UNAWARE_SHADOW:
			// This is a bit of an unusual case, because the
			// CLOSURE_UNAWARE_SHADOW Prop is a Node pointer, not a boolean.
			// However, it is treated as a boolean property in the TypedAST
			// representation as a signal that the child ASTNode is shadowed code
			// and not a normal child node.
			// We check for this bit when building in ScriptNodeDeserializer.
			propSet = Node_SetNodePropertyBit(propSet, NodeProperty_CLOSURE_UNAWARE_SHADOW)
		default:
			if !propListItem.isObject {
				nodeProperty := PropTranslator_Serialize(prop)
				if nodeProperty != nil {
					propSet = Node_SetNodePropertyBit(propSet, *nodeProperty)
				}
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
func (n *Node) SetNodePropertySideEffectFlags(propSet int64, sideEffectFlags int32) int64 {
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
func (n *Node) GetPropListDebugString() string {
	propListItem := n.propListHead
	var propListStrRepr []string
	for propListItem != nil {
		propListStrRepr = append(propListStrRepr, Node_Prop_values()[propListItem.propType].Name())
		propListItem = propListItem.next
	}
	return strings.Join(propListStrRepr, "->")
}

// Java: Node.deserializeProperties(long, boolean)
func (n *Node) DeserializeProperties(propSet int64, isInClosureUnawareSubtree bool) {
	if n.IsRoot() {
		jutil.CheckStateMsg1(n.propListHead == nil, "%s", n.propListHead)
	} else {
		jutil.CheckStateMsg(
			n.ValidatePropListTailOrdering(isInClosureUnawareSubtree),
			"Node not ready to receive properties")
	}

	// We'll gather the bits for CONST_VAR_FLAGS and SIDE_EFFECT_FLAGS into
	// these variables.
	var constantVarFlags int32 = 0
	var sideEffectFlags int32 = 0
	// Exclude the sign bit for clarity.
	for i := 0; i < 63; i++ {
		if !Node_HasBitSet(propSet, i) {
			continue
		}
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
			// If we don't ignore this, then a boolean prop is created (and then
			// later clobbered) representing this property, which seems silly at
			// best and potentially dangerous at worst.
		default:
			// All other properties are booleans that are 1-to-1 equivalent with
			// Node properties.
			prop := PropTranslator_Deserialize(nodeProperty)
			if prop == nil {
				panic(jutil.NewIllegalStateExceptionMsg("Can not translate " + nodeProperty.String() + " to AST Prop"))
			}
			n.propListHead = NewNode_IntPropListItem(int8(prop.Ordinal()), 1, n.propListHead)
		}
	}

	// Store the CONSTANT_VAR_FLAGS
	if constantVarFlags != 0 {
		n.propListHead = NewNode_IntPropListItem(
			int8(Node_Prop_CONSTANT_VAR_FLAGS.Ordinal()), constantVarFlags, n.propListHead)
	}

	if sideEffectFlags != 0 {
		n.propListHead = NewNode_IntPropListItem(
			int8(Node_Prop_SIDE_EFFECT_FLAGS.Ordinal()), sideEffectFlags, n.propListHead)
	}

	// Make sure the deserialized properties are valid.
	n.ValidateProperties(
		// NOTE: errorMessage will never be null, but just passing `false` to
		// `checkState()` triggers warning messages from some code analysis
		// tools.
		func(errorMessage string) {
			jutil.CheckStateMsg2(false, "deserialize error: %s: %s", errorMessage, n)
		})
}

// SetDeclaredTypeExpression sets the syntactical type specified on this
// node.
//
// Java: Node.setDeclaredTypeExpression(Node)
func (n *Node) SetDeclaredTypeExpression(typeExpression *Node) {
	n.PutProp(Node_Prop_DECLARED_TYPE_EXPR, typeExpression)
}

// GetDeclaredTypeExpression returns the syntactical type specified on this
// node. Not to be confused with getJSType() which returns the
// compiler-inferred type.
//
// Java: Node.getDeclaredTypeExpression()
func (n *Node) GetDeclaredTypeExpression() *Node {
	return node_cast[*Node](n.GetProp(Node_Prop_DECLARED_TYPE_EXPR))
}

// SetJSTypeBeforeCast sets the type of this node before casting.
//
// Java: Node.setJSTypeBeforeCast(JSType)
func (n *Node) SetJSTypeBeforeCast(type_ JSType) {
	n.PutProp(Node_Prop_TYPE_BEFORE_CAST, type_)
}

// GetJSTypeBeforeCast returns the type of this node before casting. This
// annotation will only exist on the first child of a CAST node after type
// checking.
//
// Java: Node.getJSTypeBeforeCast()
func (n *Node) GetJSTypeBeforeCast() JSType {
	return node_cast[JSType](n.GetProp(Node_Prop_TYPE_BEFORE_CAST))
}

// SetColorFromTypeCast indicates that this node's color comes from a type
// assertion. Only set when colors are present; when JSTypes are on the AST
// we instead preserve the actual JSType before the type assertion.
//
// Java: Node.setColorFromTypeCast()
func (n *Node) SetColorFromTypeCast() {
	jutil.CheckStateMsg(n.GetColor() != nil, "Only use on nodes with colors present")
	n.PutBooleanProp(Node_Prop_COLOR_FROM_CAST, true)
}

// IsColorFromTypeCast indicates that this node's color comes from a type
// assertion. Only set when colors are present.
//
// Java: Node.isColorFromTypeCast()
func (n *Node) IsColorFromTypeCast() bool {
	return n.GetBooleanProp(Node_Prop_COLOR_FROM_CAST)
}

// GetSortedPropTypes gets all the property types, in sorted order.
//
// Java: Node.getSortedPropTypes()
func (n *Node) GetSortedPropTypes() []int8 {
	count := 0
	for x := n.propListHead; x != nil; x = x.next {
		count++
	}

	keys := make([]int8, count)
	for x := n.propListHead; x != nil; x = x.next {
		count--
		keys[count] = x.propType
	}

	slices.Sort(keys)
	return keys
}

// Java: Node.getDouble()
func (n *Node) GetDouble() float64 {
	n.node_castKind(node_kindNumber)
	return n.node_number()
}

// Java: Node.setDouble(double)
func (n *Node) SetDouble(x float64) {
	jutil.CheckStateMsg1(!math.IsNaN(x), "%s", x)
	jutil.CheckStateMsg1(base.JSCompDoubles_IsPositive(x), "%s", x)
	n.node_castKind(node_kindNumber)
	n.node_getExt().number = x
}

// Java: Node.getBigInt()
func (n *Node) GetBigInt() *big.Int {
	n.node_castKind(node_kindBigInt)
	if n.ext == nil {
		return nil
	}
	return n.ext.bigint
}

// Java: Node.setBigInt(BigInteger)
func (n *Node) SetBigInt(number *big.Int) {
	jutil.CheckNotNull(number)
	jutil.CheckStateMsg1(number.Sign() >= 0, "%s", number)
	n.node_castKind(node_kindBigInt)
	n.node_getExt().bigint = number
}

// Java: Node.getString()
func (n *Node) GetString() string {
	n.node_castKind(node_kindString)
	return n.str
}

// Java: Node.setString(String)
func (n *Node) SetString(str string) {
	n.node_castKind(node_kindString)
	n.str = RhinoStringPool_AddOrGet(str) // RhinoStringPool is null-hostile.
}

// Java: Node.setStringFromStringPool(RhinoStringPool.LazyInternedStringList, int)
func (n *Node) SetStringFromStringPool(stringPool *RhinoStringPool_LazyInternedStringList, offset int) {
	// RhinoStringPool.LazyInternedStringList guarantees that the result of
	// calling get(offset) is interned. In some cases this is more performance
	// than RhinoStringPool.addOrGet because the LazyInternedStringList caches
	// interning results.
	n.node_castKind(node_kindString)
	n.str = stringPool.Get(offset)
}

// Java: Node.getRawString()
func (n *Node) GetRawString() string {
	n.node_castKind(node_kindTemplateLiteralSubstring)
	return n.str // raw
}

// Java: Node.getCookedString()
func (n *Node) GetCookedString() *string {
	n.node_castKind(node_kindTemplateLiteralSubstring)
	return n.node_cooked()
}

// Java: Node.toString()
func (n *Node) String() string {
	return n.ToStringWithFlags(true, true, true)
}

// Java: Node.toString(boolean, boolean, boolean)
func (n *Node) ToStringWithFlags(printSource bool, printAnnotations bool, printType bool) string {
	var sb strings.Builder
	n.ToStringInto(&sb, printSource, printAnnotations, printType)
	return sb.String()
}

// Java: Node.toString(StringBuilder, boolean, boolean, boolean)
func (n *Node) ToStringInto(sb *strings.Builder, printSource bool, printAnnotations bool, printType bool) {
	sb.WriteString(n.token.String())
	if n.kind == node_kindString {
		sb.WriteByte(' ')
		sb.WriteString(n.GetString())
	} else if n.token == Token_FUNCTION {
		sb.WriteByte(' ')
		// In the case of JsDoc trees, the first child is often not a string
		// which causes exceptions to be thrown when calling toString or
		// toStringTree.
		if n.first == nil || n.first.token != Token_NAME {
			sb.WriteString("<invalid>")
		} else {
			sb.WriteString(n.first.GetString())
		}
	} else if n.token == Token_NUMBER {
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
		if n.length != 0 {
			sb.WriteString(" [length: ")
			sb.WriteString(strconv.Itoa(int(n.length)))
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
		if originalName := n.node_originalName(); originalName != nil {
			sb.WriteString(" [original_name: ")
			sb.WriteString(*originalName)
			sb.WriteByte(']')
		}
	}

	if printType && n.jstypeOrColor != nil {
		typeString := jutil.ToJavaString(n.jstypeOrColor)
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
func (n *Node) ToJson(appendable io.StringWriter) {
	sb := node_appender{appendable}
	sb.append("{")
	sb.append(Node_CreateJsonPair("token", n.token.String()))
	if n.kind == node_kindString {
		sb.append(",")
		sb.append(Node_CreateJsonPair("string", n.GetString()))
	} else if n.token == Token_FUNCTION {
		sb.append(",")
		// In the case of JsDoc trees, the first child is often not a string
		// which causes exceptions to be thrown when calling toString or
		// toStringTree.
		if n.first == nil || n.first.token != Token_NAME {
			sb.append(Node_CreateJsonPair("functionName", "<invalid>"))
		} else {
			sb.append(Node_CreateJsonPair("functionName", n.first.GetString()))
		}
	} else if n.token == Token_NUMBER {
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
		if n.length != 0 {
			sb.append(",")
			sb.append(Node_CreateJsonPairRawValue("length", strconv.Itoa(int(n.length))))
		}
		sb.append("}")
	}

	if originalName := n.node_originalName(); originalName != nil {
		sb.append(",")
		sb.append(Node_CreateJsonPair("original_name", *originalName))
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

	if n.jstypeOrColor != nil {
		typeString := jutil.ToJavaString(n.jstypeOrColor)
		sb.append(",")
		sb.append(Node_CreateJsonPair("typeString", jutil.StringReplace(typeString, "\"", "\\\"")))
	}

	if n.first != nil {
		sb.append(",")
		sb.append("\"children\":[")
		for child := n.first; child != nil; child = child.next {
			child.ToJson(appendable)
			if child.next != nil {
				sb.append(",")
			}
		}
		sb.append("]")
	}
	sb.append("}")
}

// Java: Node.toStringTree()
func (n *Node) ToStringTree() string {
	return n.ToStringTreeImpl()
}

// Java: Node.toStringTreeImpl()
func (n *Node) ToStringTreeImpl() string {
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
func (n *Node) AppendStringTree(appendable io.StringWriter) {
	Node_ToStringTreeHelper(n, 0, appendable)
}

// Java: Node.toStringTreeHelper(Node, int, Appendable)
func Node_ToStringTreeHelper(n *Node, level int, appendable io.StringWriter) {
	sb := node_appender{appendable}
	for i := 0; i != level; i++ {
		sb.append("    ")
	}
	sb.append(n.String())
	sb.append("\n")
	for cursor := n.first; cursor != nil; cursor = cursor.next {
		Node_ToStringTreeHelper(cursor, level+1, appendable)
	}
	shadow := n.GetClosureUnawareShadow()
	if shadow != nil {
		Node_ToStringTreeHelper(shadow, level+1, appendable)
	}
}

// Java: Node.appendJsonTree(Appendable)
func (n *Node) AppendJsonTree(appendable io.StringWriter) {
	Node_ToJsonTreeHelper(n, appendable)
}

// Java: Node.toJsonTreeHelper(Node, Appendable)
func Node_ToJsonTreeHelper(n *Node, appendable io.StringWriter) {
	n.ToJson(appendable)
}

// Node is the AST node (fields as declared in Java, plus the folded
// subclass payloads).
//
// Go-only: performance — the fields that traversals read first (token, the
// position, next, first and parent) lead, packed into 32 bytes. The int32
// fields share words, the string payload of StringNode and
// TemplateLiteralSubstringNode share one field (a node is only ever one of
// them), and the fields most nodes never set (originalName and the number,
// bigint and cooked payloads) live in a node_ext allocated on first use: a
// Node takes 96 bytes instead of 144, so that a traversal touches a third
// fewer cache lines, and the 32 hot bytes of every node lie in one line.
type Node struct {
	token Token // Type of the token of the node; NAME for example

	// Source position of this node. The position is encoded with the column
	// number in the low 12 bits of the integer, and the line number in the
	// rest. Create some handy constants so we can change this size if we
	// want.
	linenoCharno int32

	next     *Node // next sibling, a linked list
	first    *Node // first element of a linked list of children
	parent   *Node
	previous *Node // previous sibling, a circular linked list

	// We get the last child as first.previous. But last.next is null, not
	// first.

	// The length of the code represented by the node (a Java int).
	length int32

	// Go-only: language-adapter — the folded subclass (see the Node doc
	// comment); its payload is str and ext.
	kind node_kind

	jstypeOrColor any

	// Linked list of properties. Since vast majority of nodes would have no
	// more than 2 properties, linked list saves memory and provides fast
	// lookup. If this does not holds, propListHead can be replaced by UintMap.
	propListHead *Node_PropListItem

	// StringNode.str, or TemplateLiteralSubstringNode.raw: the raw version of
	// the template literal substring, is not null.
	str string

	// originalName and the NumberNode, BigIntNode and
	// TemplateLiteralSubstringNode payloads; nil until one is set.
	ext *node_ext
}

// node_ext holds the fields of a Node that most nodes never set.
//
// Go-only: performance — see the Node doc comment. Each Node has its own
// (never shared), so writes through it affect only that node.
type node_ext struct {
	originalName *string

	number float64  // NumberNode.number
	bigint *big.Int // BigIntNode.bigint
	// TemplateLiteralSubstringNode.cooked: the cooked version of the template
	// literal substring. Null iff the raw string contains an uncookable
	// escape sequence.
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals#es2018_revision_of_illegal_escape_sequences
	cooked *string
}

// node_getExt returns n.ext, allocating it on first use.
func (n *Node) node_getExt() *node_ext {
	if n.ext == nil {
		n.ext = &node_ext{}
	}
	return n.ext
}

// node_originalName is the originalName field.
func (n *Node) node_originalName() *string {
	if n.ext == nil {
		return nil
	}
	return n.ext.originalName
}

// node_setOriginalName assigns the originalName field.
func (n *Node) node_setOriginalName(s *string) {
	if s == nil {
		if n.ext != nil {
			n.ext.originalName = nil
		}
		return
	}
	n.node_getExt().originalName = s
}

// node_number is the NumberNode.number field.
func (n *Node) node_number() float64 {
	if n.ext == nil {
		return 0
	}
	return n.ext.number
}

// node_cooked is the TemplateLiteralSubstringNode.cooked field.
func (n *Node) node_cooked() *string {
	if n.ext == nil {
		return nil
	}
	return n.ext.cooked
}

// ==========================================================================
// Source position management

// Java: Node.setStaticSourceFileFrom(Node)
func (n *Node) SetStaticSourceFileFrom(other *Node) {
	// Make sure source file prop nodes are not duplicated.
	if other.propListHead != nil &&
		(n.propListHead == nil ||
			(int(n.propListHead.propType) == Node_Prop_SOURCE_FILE.Ordinal() &&
				n.propListHead.next == nil)) {
		// Either the node has only Prop.SOURCE_FILE as a property or has not
		// properties.
		tail := other.propListHead
		for tail.next != nil {
			tail = tail.next
		}
		if int(tail.propType) == Node_Prop_SOURCE_FILE.Ordinal() {
			n.propListHead = tail
			return
		}
	}
	n.SetStaticSourceFile(other.GetStaticSourceFile())
}

// Java: Node.setStaticSourceFile(StaticSourceFile)
func (n *Node) SetStaticSourceFile(file StaticSourceFile) *Node {
	n.PutProp(Node_Prop_SOURCE_FILE, file)
	return n
}

// SetSourceFileForTesting sets the source file to a non-extern file of the
// given name.
//
// Java: Node.setSourceFileForTesting(String)
func (n *Node) SetSourceFileForTesting(name string) {
	n.PutProp(Node_Prop_SOURCE_FILE, NewSimpleSourceFile(name, StaticSourceFile_SourceKind_STRONG))
}

// TODO(johnlenz): make this final
//
// Java: Node.getSourceFileName()
func (n *Node) GetSourceFileName() *string {
	file := n.GetStaticSourceFile()
	if file == nil {
		return nil
	}
	name := file.GetName()
	return &name
}

// GetStaticSourceFile returns the source file associated with this input.
//
// Java: Node.getStaticSourceFile()
func (n *Node) GetStaticSourceFile() StaticSourceFile {
	return node_cast[StaticSourceFile](n.GetProp(Node_Prop_SOURCE_FILE))
}

// SetInputId sets the ID of the input this Node came from.
//
// Java: Node.setInputId(InputId)
func (n *Node) SetInputId(inputId *InputId) {
	n.PutProp(Node_Prop_INPUT_ID, inputId)
}

// GetInputId returns the Id of the CompilerInput associated with this Node.
//
// Java: Node.getInputId()
func (n *Node) GetInputId() *InputId {
	return node_cast[*InputId](n.GetProp(Node_Prop_INPUT_ID))
}

// GetOriginalName returns the original name of this node, if the node has
// been renamed.
//
// Do not use original name to make optimization decisions. The original
// intent was to preserve some naming for lightly optimized code to put into
// source maps. It is not rigorously defined and is not a suitable
// replacement for a canonical identifier. "Original name" is not associated
// with any scope and easily transfers to unrelated values. Its existence and
// use beyond its original purpose has delayed creating useful more precise
// alternatives.
//
// Deprecated: "original name" is poorly defined.
//
// Java: Node.getOriginalName()
func (n *Node) GetOriginalName() *string {
	return n.node_originalName()
}

// Java: Node.setOriginalName(String)
func (n *Node) SetOriginalName(s *string) {
	if s == nil {
		n.node_setOriginalName(nil)
	} else {
		interned := RhinoStringPool_AddOrGet(*s)
		n.node_setOriginalName(&interned)
	}
}

// Java: Node.setOriginalNameFromStringPool(RhinoStringPool.LazyInternedStringList, int)
func (n *Node) SetOriginalNameFromStringPool(stringPool *RhinoStringPool_LazyInternedStringList, offset int) {
	// RhinoStringPool.LazyInternedStringList guarantees that the result of
	// calling get(offset) is interned, and sometimes caches results.
	s := stringPool.Get(offset)
	n.node_setOriginalName(&s)
}

// SetOriginalNameFromName copies the original name.
//
// Java: Node.setOriginalNameFromName(Node)
func (n *Node) SetOriginalNameFromName(name *Node) {
	s := name.GetString()
	n.node_setOriginalName(&s)
}

// IsIndexable reports whether this node should be indexed by static
// analysis / code indexing tools.
//
// Java: Node.isIndexable()
func (n *Node) IsIndexable() bool {
	return !n.GetBooleanProp(Node_Prop_NON_INDEXABLE)
}

// Java: Node.makeNonIndexable()
func (n *Node) MakeNonIndexable() {
	n.PutBooleanProp(Node_Prop_NON_INDEXABLE, true)
}

// Java: Node.makeNonIndexableRecursive()
func (n *Node) MakeNonIndexableRecursive() {
	n.MakeNonIndexable()
	for child := n.first; child != nil; child = child.GetNext() {
		child.MakeNonIndexableRecursive()
	}
}

// Java: Node.isFromExterns()
func (n *Node) IsFromExterns() bool {
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
func (n *Node) IsClosureUnawareCode() bool {
	file := n.GetStaticSourceFile()
	return file != nil && file.IsClosureUnawareCode()
}

// Java: Node.getLength()
func (n *Node) GetLength() int {
	return int(n.length)
}

// Java: Node.setLength(int)
func (n *Node) SetLength(length int) {
	n.length = int32(length)
}

// Java: Node.getLineno()
func (n *Node) GetLineno() int {
	if n.linenoCharno == -1 {
		return -1
	}
	return int(jutil.Ushr32(n.linenoCharno, Node_CHARNO_BITS))
}

// GetCharno returns the 0-based column number.
//
// Java: Node.getCharno()
func (n *Node) GetCharno() int {
	if n.linenoCharno == -1 {
		return -1
	}
	return int(n.linenoCharno & Node_MAX_COLUMN_NUMBER)
}

// Java: Node.getLocation()
func (n *Node) GetLocation() string {
	return jutil.ToJavaString(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

// TODO(johnlenz): make this final
//
// Java: Node.getSourceOffset()
func (n *Node) GetSourceOffset() int {
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
func (n *Node) GetSourcePosition() int {
	return int(n.linenoCharno)
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
func (n *Node) SetLinenoCharno(lineno int, charno int) *Node {
	if lineno < 0 || charno < 0 {
		n.linenoCharno = -1
		return n
	}

	if charno > Node_MAX_COLUMN_NUMBER {
		charno = Node_MAX_COLUMN_NUMBER
	}
	n.linenoCharno = jutil.Shl32(int32(lineno), Node_CHARNO_BITS) | int32(charno)

	return n
}

// ==========================================================================
// Iteration

// Children returns an iterable object that iterates over this node's
// children. The iterator does not support the optional operation
// Iterator#remove().
//
// To iterate over a node's children, one can write
//
//	for child := range n.Children() { ...
//
// NOTE: Do not use 'children' for recursive descent of the AST. The overhead
// of using iterators rather then getFirstChild()/getNext() is very
// significant. We have deprecated it as it is easy to misuse.
//
// Go-only: language-adapter — Java's Iterable is an iter.Seq; the private
// SiblingNodeIterable/SiblingNodeIterator classes are the closure below
// (which, like the Java iterator, reads the next sibling before handing out
// the current node).
//
// Deprecated: see above.
//
// Java: Node.children()
func (n *Node) Children() iter.Seq[*Node] {
	if n.first == nil {
		return func(yield func(*Node) bool) {}
	}
	start := n.first
	return func(yield func(*Node) bool) {
		current := start
		for current != nil {
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

// Java: Node.getPropListHeadForTesting()
func (n *Node) GetPropListHeadForTesting() *Node_PropListItem {
	return n.propListHead
}

// Java: Node.setPropListHead(PropListItem)
func (n *Node) SetPropListHead(propListHead *Node_PropListItem) {
	n.propListHead = propListHead
}

// Java: Node.getParent()
func (n *Node) GetParent() *Node {
	return n.parent
}

// Java: Node.hasParent()
func (n *Node) HasParent() bool {
	return n.parent != nil
}

// Java: Node.getGrandparent()
func (n *Node) GetGrandparent() *Node {
	if n.parent == nil {
		return nil
	}
	return n.parent.parent
}

// GetAncestor gets the ancestor node relative to this.
//
// level: 0 = this, 1 = the parent, etc.
//
// Java: Node.getAncestor(int)
func (n *Node) GetAncestor(level int) *Node {
	jutil.CheckArgument(level >= 0)
	node := n
	for node != nil {
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
func (n *Node) IsDescendantOf(node *Node) bool {
	for c := n; c != nil; c = c.parent {
		if c == node {
			return true
		}
	}
	return false
}

// Java: Node.isOnlyChildOf(Node)
func (n *Node) IsOnlyChildOf(possibleParent *Node) bool {
	return possibleParent == n.GetParent() && n.GetPrevious() == nil && n.GetNext() == nil
}

// Java: Node.isFirstChildOf(Node)
func (n *Node) IsFirstChildOf(possibleParent *Node) bool {
	return possibleParent == n.GetParent() && n.GetPrevious() == nil
}

// Java: Node.isSecondChildOf(Node)
func (n *Node) IsSecondChildOf(possibleParent *Node) bool {
	previousNode := n.GetPrevious()

	return previousNode != nil && previousNode.IsFirstChildOf(possibleParent)
}

// GetAncestors iterates all of the node's ancestors excluding itself.
//
// Java: Node.getAncestors()
func (n *Node) GetAncestors() *Node_AncestorIterable {
	return NewNode_AncestorIterable(n.GetParent())
}

// Node_AncestorIterable is an iterator to go up the ancestor tree.
type Node_AncestorIterable struct {
	cur *Node
}

// NewNode_AncestorIterable starts at `cur`.
//
// Java: Node.AncestorIterable(Node)
func NewNode_AncestorIterable(cur *Node) *Node_AncestorIterable {
	return &Node_AncestorIterable{cur: cur}
}

// All is the iterable's iterator(). Like the Java iterator, it advances the
// iterable's own cursor, so an AncestorIterable can be traversed only once.
//
// Java: Node.AncestorIterable.iterator()
func (a *Node_AncestorIterable) All() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for a.cur != nil {
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
func (n *Node) HasOneChild() bool {
	return n.first != nil && n.first.next == nil
}

// HasTwoChildren checks for two children more efficiently than
// `getChildCount() == 2`.
//
// Java: Node.hasTwoChildren()
func (n *Node) HasTwoChildren() bool {
	return n.first != nil && n.first.next != nil && n.first.next == n.GetLastChild()
}

// HasZeroOrOneChild checks for zero or one child more efficiently than by
// iterating over all the children as is done with Node.getChildCount().
//
// Java: Node.hasZeroOrOneChild()
func (n *Node) HasZeroOrOneChild() bool {
	return n.first == n.GetLastChild()
}

// HasMoreThanOneChild checks for more than one child more efficiently than
// by iterating over all the children as is done with Node.getChildCount().
//
// Java: Node.hasMoreThanOneChild()
func (n *Node) HasMoreThanOneChild() bool {
	return n.first != nil && n.first.next != nil
}

// HasXChildren checks for has exactly the number of specified children.
//
// Java: Node.hasXChildren(int)
func (n *Node) HasXChildren(x int) bool {
	c := 0
	for c2 := n.first; c2 != nil && c <= x; c2 = c2.next {
		c++
	}
	return c == x
}

// Java: Node.getChildCount()
func (n *Node) GetChildCount() int {
	c := 0
	for c2 := n.first; c2 != nil; c2 = c2.next {
		c++
	}
	return c
}

// HasChild is intended for testing and verification only.
//
// Java: Node.hasChild(Node)
func (n *Node) HasChild(child *Node) bool {
	for c := n.first; c != nil; c = c.next {
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
func (n *Node) IsEquivalentToShallow(node *Node) bool {
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
func (n *Node) IsEquivalentWithSideEffectsTo(node *Node) bool {
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
func (n *Node) IsEquivalentWithSideEffectsToShallow(node *Node) bool {
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
func (n *Node) IsEquivalentToTyped(node *Node) bool {
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
func (n *Node) IsEquivalentTo(node *Node) bool {
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
func (n *Node) GetSecondToLastPropListItem() *Node_PropListItem {
	propListItem := n.propListHead
	if propListItem == nil || propListItem.next == nil {
		return nil
	}
	for propListItem.next.next != nil {
		propListItem = propListItem.next
	}
	return propListItem
}

// Java: Node.getLastPropListItem()
func (n *Node) GetLastPropListItem() *Node_PropListItem {
	propListItem := n.propListHead
	if propListItem == nil {
		return nil
	}
	for propListItem.next != nil {
		propListItem = propListItem.next
	}
	return propListItem
}

// Java: Node.validatePropListTailOrdering(boolean)
func (n *Node) ValidatePropListTailOrdering(shouldNodeBeClosureUnaware bool) bool {
	var prev *Node_PropListItem
	curr := n.propListHead
	if curr == nil {
		return false
	}
	for curr.next != nil {
		prev = curr
		curr = curr.next
	}
	if int(curr.propType) != Node_Prop_SOURCE_FILE.Ordinal() {
		return false
	}

	if shouldNodeBeClosureUnaware {
		if prev == nil || int(prev.propType) != Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE.Ordinal() {
			return false
		}
	}
	return true
}

// Java: Node.validateMemorySensitivePropertyGuarantees(Node, boolean, Node, boolean)
func Node_ValidateMemorySensitivePropertyGuarantees(
	node1 *Node,
	shouldNode1BeClosureUnaware bool,
	node2 *Node,
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
func (n *Node) IsEquivalentToWithFlags(node *Node, compareType bool, recurse bool, jsDoc bool) bool {
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
func (n *Node) IsEquivalentToWithFlagsAndSideEffect(node *Node, compareType bool, recurse bool, jsDoc bool, sideEffect bool) bool {
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
func (n *Node) IsEquivalentToWithModes(
	node *Node,
	recursion Node_RecursionMode,
	typeConfig Node_TypeComparison,
	jsDocConfig Node_JsDocComparison,
	sideEffectConfig Node_SideEffectComparison,
) bool {
	switch n.kind {
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
			RhinoStringPool_UncheckedEquals(n.str, node.str)
	case node_kindTemplateLiteralSubstring:
		// Java: TemplateLiteralSubstringNode.isEquivalentTo
		if !n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
			return false
		}

		castNode := node
		return RhinoStringPool_UncheckedEquals(n.str, castNode.str) && // raw
			RhinoStringPool_UncheckedEqualsNullable(n.node_cooked(), castNode.node_cooked())
	}
	return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig)
}

// isEquivalentToBase is Node's own isEquivalentTo (the super call of the
// folded subclasses).
func (n *Node) isEquivalentToBase(
	node *Node,
	recursion Node_RecursionMode,
	typeConfig Node_TypeComparison,
	jsDocConfig Node_JsDocComparison,
	sideEffectConfig Node_SideEffectComparison,
) bool {
	if n.token != node.token ||
		n.GetChildCount() != node.GetChildCount() ||
		n.kind != node.kind {
		return false
	}

	if typeConfig == Node_TypeComparison_COMPARE &&
		!jutil.ObjectsEquals(n.jstypeOrColor, node.jstypeOrColor) {
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
	} else if thisDte == nil || thatDte == nil {
		return false
	} else if !thisDte.IsEquivalentToWithModes(
		thatDte, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
		return false
	}

	propSet := jutil.NewEnumSet[Node_Prop]()
	for propListItem := n.propListHead; propListItem != nil; propListItem = propListItem.next {
		prop := Node_propValues[propListItem.propType]
		propSet.Add(prop)
	}
	for propListItem := node.propListHead; propListItem != nil; propListItem = propListItem.next {
		prop := Node_propValues[propListItem.propType]
		propSet.Add(prop)
	}

	for prop := range propSet.All() {
		if Node_propMapForEqualityKeys.Contains(prop) {
			getter := Node_PROP_MAP_FOR_EQUALITY.Get(prop)
			if !jutil.ObjectsEquals(getter(n), getter(node)) {
				return false
			}
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

	thisIsShadowHost := n.GetClosureUnawareShadow() != nil
	nodeIsShadowHost := node.GetClosureUnawareShadow() != nil
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
		for c, c2 := n.first, node.first; c != nil; c, c2 = c.next, c2.next {
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
var Node_PROP_MAP_FOR_EQUALITY = func() *jutil.LinkedHashMap[Node_Prop, func(*Node) any] {
	m := jutil.NewLinkedHashMap[Node_Prop, func(*Node) any]()
	m.Put(Node_Prop_ARROW_FN, func(n *Node) any { return n.IsArrowFunction() })
	m.Put(Node_Prop_ASYNC_FN, func(n *Node) any { return n.IsAsyncFunction() })
	m.Put(Node_Prop_GENERATOR_FN, func(n *Node) any { return n.IsGeneratorFunction() })
	m.Put(Node_Prop_START_OF_OPT_CHAIN, func(n *Node) any { return n.IsOptionalChainStart() })
	m.Put(Node_Prop_STATIC_MEMBER, func(n *Node) any { return n.IsStaticMember() })
	m.Put(Node_Prop_YIELD_ALL, func(n *Node) any { return n.IsYieldAll() })
	m.Put(Node_Prop_EXPORT_DEFAULT, func(n *Node) any { return n.GetIntProp(Node_Prop_EXPORT_DEFAULT) })
	m.Put(Node_Prop_EXPORT_ALL_FROM, func(n *Node) any { return n.GetIntProp(Node_Prop_EXPORT_ALL_FROM) })
	m.Put(Node_Prop_INCRDECR, func(n *Node) any { return n.GetIntProp(Node_Prop_INCRDECR) })
	m.Put(Node_Prop_QUOTED, func(n *Node) any { return n.GetIntProp(Node_Prop_QUOTED) })
	m.Put(Node_Prop_FREE_CALL, func(n *Node) any { return n.GetBooleanProp(Node_Prop_FREE_CALL) })
	m.Put(Node_Prop_DIRECT_EVAL, func(n *Node) any { return n.GetBooleanProp(Node_Prop_DIRECT_EVAL) })
	m.Put(Node_Prop_COMPUTED_PROP_METHOD, func(n *Node) any { return n.GetBooleanProp(Node_Prop_COMPUTED_PROP_METHOD) })
	m.Put(Node_Prop_COMPUTED_PROP_GETTER, func(n *Node) any { return n.GetBooleanProp(Node_Prop_COMPUTED_PROP_GETTER) })
	m.Put(Node_Prop_COMPUTED_PROP_SETTER, func(n *Node) any { return n.GetBooleanProp(Node_Prop_COMPUTED_PROP_SETTER) })
	return m
}()

// Node_propMapForEqualityKeys is used for faster Map.containsKey() lookups
// in PROP_MAP_FOR_EQUALITY.
//
// Java: Node.propMapForEqualityKeys
var Node_propMapForEqualityKeys = jutil.NewEnumSet(Node_PROP_MAP_FOR_EQUALITY.KeySlice()...)

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
func (n *Node) GetQualifiedName() *string {
	s, ok := n.getQualifiedNameString()
	if !ok {
		return nil
	}
	return &s
}

// getQualifiedNameString is the body of getQualifiedName; ok is false
// where Java returns null.
func (n *Node) getQualifiedNameString() (string, bool) {
	switch n.token {
	case Token_NAME:
		name := n.GetString()
		if name == "" {
			return "", false
		}
		return name, true
	case Token_GETPROP:
		builder := n.GetQualifiedNameForGetProp(0)
		if builder != nil {
			// Go-only: performance — builder.toString() without a copy, as
			// strings.Builder.String does: nothing writes to builder after.
			return unsafe.String(unsafe.SliceData(builder), len(builder)), true
		}
		return "", false
	case Token_THIS:
		return "this", true
	case Token_SUPER:
		return "super", true
	default:
		return "", false
	}
}

// Java: Node.getQualifiedNameObject()
func (n *Node) GetQualifiedNameObject() QualifiedName {
	if n.IsQualifiedName() {
		return NewQualifiedName_NodeQname(n)
	}
	return nil
}

// GetQualifiedNameForGetProp is a helper method for getQualifiedName to
// handle GETPROP nodes.
//
// reserve: the number of characters of space to reserve in the
// StringBuilder. Returns null if this is not a qualified name or a
// StringBuilder if it is a complex qualified name.
//
// Go-only: performance — the StringBuilder is a byte slice (nil for null),
// filled with append: one allocation, made at the full size, instead of a
// heap strings.Builder and its buffer; the left name is read without
// boxing it (getQualifiedNameString).
//
// Java: Node.getQualifiedNameForGetProp(int)
func (n *Node) GetQualifiedNameForGetProp(reserve int) []byte {
	propName := n.GetString()
	reserve += 1 + len(propName) // +1 for the '.'
	var builder []byte
	if n.first.IsGetProp() {
		builder = n.first.GetQualifiedNameForGetProp(reserve)
		if builder == nil {
			return nil
		}
	} else {
		left, ok := n.first.getQualifiedNameString()
		if !ok {
			return nil
		}
		builder = make([]byte, 0, len(left)+reserve)
		builder = append(builder, left...)
	}
	builder = append(builder, '.')
	builder = append(builder, propName...)
	return builder
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
func (n *Node) GetOriginalQualifiedName() *string {
	if n.token == Token_NAME {
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
	} else if n.token == Token_GETPROP {
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
	} else if n.token == Token_THIS {
		s := "this"
		return &s
	} else if n.token == Token_SUPER {
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
func (n *Node) IsQualifiedName() bool {
	switch n.GetToken() {
	case Token_NAME:
		return n.GetString() != ""
	case Token_THIS, Token_SUPER:
		return true
	case Token_GETPROP:
		return n.GetFirstChild().IsQualifiedName()

	case Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// MatchesName returns whether a node matches a simple name, such as `x`,
// returns false if this is not a NAME node.
//
// Java: Node.matchesName(String)
func (n *Node) MatchesName(name string) bool {
	if n.token != Token_NAME {
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
func (n *Node) MatchesNameNode(other *Node) bool {
	if n.token != Token_NAME || other.token != Token_NAME {
		return false
	}

	// ==, rather than equal as it is interned in setString
	internalString := n.GetString()
	return internalString != "" &&
		RhinoStringPool_UncheckedEquals(internalString, other.GetString())
}

// MatchesQualifiedName returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(String)
func (n *Node) MatchesQualifiedName(name string) bool {
	return n.MatchesQualifiedNameWithEnd(name, jutil.StringLength(name))
}

// MatchesQualifiedNameWithEnd returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(String, int)
func (n *Node) MatchesQualifiedNameWithEnd(qname string, endIndex int) bool {
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
func (n *Node) MatchesQualifiedNameNode(other *Node) bool {
	if other.token != n.token {
		return false
	}

	switch n.token {
	case Token_NAME:
		return n.MatchesNameNode(other)
	case Token_THIS, Token_SUPER:
		return true
	case Token_GETPROP:
		return RhinoStringPool_UncheckedEquals(n.GetString(), other.GetString()) &&
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
func (n *Node) IsUnscopedQualifiedName() bool {
	switch n.GetToken() {
	case Token_NAME:
		return n.GetString() != ""
	case Token_GETPROP:
		return n.GetFirstChild().IsUnscopedQualifiedName()
	default:
		return false
	}
}

// Java: Node.isValidAssignmentTarget()
func (n *Node) IsValidAssignmentTarget() bool {
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
func (n *Node) Clone() any {
	panic(jutil.NewUnsupportedOperationExceptionMsg("Did you mean cloneNode?"))
}

// CloneNode returns a detached clone of the Node, specifically excluding its
// children.
//
// Java: Node.cloneNode()
func (n *Node) CloneNode() *Node {
	return n.CloneNodeWithTypeExprs(false)
}

// CloneNodeWithTypeExprs returns a detached clone of the Node, specifically
// excluding its children.
//
// Go-only: language-adapter — dispatches to the cloneNode overrides of the
// folded subclasses.
//
// Java: Node.cloneNode(boolean)
func (n *Node) CloneNodeWithTypeExprs(cloneTypeExprs bool) *Node {
	switch n.kind {
	case node_kindNumber:
		// Java: NumberNode.cloneNode(boolean)
		clone := newNode_NumberNode(n.node_number())
		Node_CopyBaseNodeFields(n, clone, cloneTypeExprs)
		return clone
	case node_kindBigInt:
		// Java: BigIntNode.cloneNode(boolean)
		clone := newNode_BigIntNode(n.GetBigInt())
		Node_CopyBaseNodeFields(n, clone, cloneTypeExprs)
		return clone
	case node_kindString:
		// Java: StringNode.cloneNode(boolean)
		clone := newNode_StringNodeForClone(n.GetToken())
		Node_CopyBaseNodeFields(n, clone, cloneTypeExprs)
		clone.str = n.str
		return clone
	case node_kindTemplateLiteralSubstring:
		// Java: TemplateLiteralSubstringNode.cloneNode(boolean)
		clone := newNode_TemplateLiteralSubstringNode(n.node_cooked(), n.str)
		Node_CopyBaseNodeFields(n, clone, cloneTypeExprs)
		return clone
	}
	clone := NewNode(n.token)
	Node_CopyBaseNodeFields(n, clone, cloneTypeExprs)
	return clone
}

// Java: Node.copyBaseNodeFields(Node, Node, boolean)
func Node_CopyBaseNodeFields(source *Node, dest *Node, cloneTypeExprs bool) {
	dest.linenoCharno = source.linenoCharno
	dest.length = source.length
	dest.jstypeOrColor = source.jstypeOrColor
	dest.node_setOriginalName(source.node_originalName())
	dest.propListHead = source.propListHead

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
func (n *Node) CloneTree() *Node {
	return n.CloneTreeWithTypeExprs(false)
}

// Java: Node.cloneTree(boolean)
func (n *Node) CloneTreeWithTypeExprs(cloneTypeExprs bool) *Node {
	result := n.CloneNodeWithTypeExprs(cloneTypeExprs)
	var firstChild *Node
	var lastChild *Node
	if n.HasChildren() {
		for n2 := n.GetFirstChild(); n2 != nil; n2 = n2.next {
			n2clone := n2.CloneTreeWithTypeExprs(cloneTypeExprs)
			n2clone.parent = result
			if firstChild == nil {
				firstChild = n2clone
				lastChild = firstChild
			} else {
				lastChild.next = n2clone
				n2clone.previous = lastChild
				lastChild = n2clone
			}
		}
		firstChild.previous = lastChild
		lastChild.next = nil
		result.first = firstChild
	}

	shadow := n.GetClosureUnawareShadow()
	if shadow != nil {
		result.SetClosureUnawareShadow(shadow.CloneTreeWithTypeExprs(cloneTypeExprs))
	}

	return result
}

// Srcref copies the source info from `other` onto `this`.
//
// Java: Node.srcref(Node)
func (n *Node) Srcref(other *Node) *Node {
	n.SetStaticSourceFileFrom(other)
	n.node_setOriginalName(other.node_originalName())
	n.linenoCharno = other.linenoCharno
	n.length = other.length
	return n
}

// SrcrefTree copies, for all Nodes in the subtree of `this`, the source info
// from `other`.
//
// Java: Node.srcrefTree(Node)
func (n *Node) SrcrefTree(other *Node) *Node {
	n.Srcref(other)
	for child := n.first; child != nil; child = child.next {
		child.SrcrefTree(other)
	}
	return n
}

// SrcrefIfMissing copies the source info from `other` iff source info is not
// set on `this`.
//
// Java: Node.srcrefIfMissing(Node)
func (n *Node) SrcrefIfMissing(other *Node) *Node {
	if n.GetStaticSourceFile() == nil {
		n.SetStaticSourceFileFrom(other)
		n.linenoCharno = other.linenoCharno
		n.length = other.length
	}

	// TODO(lharker): should this be inside the above if condition?
	// If the node already has a source file, it seems strange to
	// go ahead and set the original name anyway.
	if n.node_originalName() == nil {
		n.node_setOriginalName(other.node_originalName())
	}

	return n
}

// SrcrefTreeIfMissing copies, for all Nodes in the subtree of `this`, iff
// source info is not set, the source info from `other`.
//
// Java: Node.srcrefTreeIfMissing(Node)
func (n *Node) SrcrefTreeIfMissing(other *Node) *Node {
	n.SrcrefIfMissing(other)
	for child := n.first; child != nil; child = child.next {
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
func (n *Node) GetJSType() JSType {
	if t, ok := n.jstypeOrColor.(JSType); ok {
		return t
	}
	return nil
}

// GetJSTypeRequired returns the compiled inferred type on this node, or
// throws an NPE if there isn't one.
//
// Java: Node.getJSTypeRequired()
func (n *Node) GetJSTypeRequired() JSType {
	return jutil.CheckNotNullMsg1(n.GetJSType(), "no jstypeOrColor: %s", n)
}

// Java: Node.setJSType(JSType)
func (n *Node) SetJSType(x JSType) *Node {
	_, isJSType := n.jstypeOrColor.(JSType)
	jutil.CheckStateMsg1(n.jstypeOrColor == nil || isJSType, "%s", n)
	if jutil.IsNil(x) {
		n.jstypeOrColor = nil
	} else {
		n.jstypeOrColor = x
	}
	return n
}

// GetColor returns the compiled inferred type on this node. Not to be
// confused with getDeclaredTypeExpression() which returns the syntactically
// specified type.
//
// Java: Node.getColor()
func (n *Node) GetColor() *colors.Color {
	if c, ok := n.jstypeOrColor.(*colors.Color); ok {
		return c
	}
	return nil
}

// Java: Node.setColor(Color)
func (n *Node) SetColor(x *colors.Color) *Node {
	_, isColor := n.jstypeOrColor.(*colors.Color)
	jutil.CheckStateMsg1(n.jstypeOrColor == nil || isColor, "%s", n)
	if x == nil {
		n.jstypeOrColor = nil
	} else {
		n.jstypeOrColor = x
	}
	return n
}

// CopyTypeFrom copies a nodes JSType or Color (if present).
//
// Java: Node.copyTypeFrom(Node)
func (n *Node) CopyTypeFrom(other *Node) *Node {
	n.jstypeOrColor = other.jstypeOrColor
	return n
}

// GetJSDocInfo gets the JSDocInfo attached to this node; null if no JSDoc is
// attached to this node.
//
// Java: Node.getJSDocInfo()
func (n *Node) GetJSDocInfo() *JSDocInfo {
	return node_cast[*JSDocInfo](n.GetProp(Node_Prop_JSDOC_INFO))
}

// SetJSDocInfo sets the JSDocInfo attached to this node.
//
// Java: Node.setJSDocInfo(JSDocInfo)
func (n *Node) SetJSDocInfo(info *JSDocInfo) *Node {
	n.PutProp(Node_Prop_JSDOC_INFO, info)
	return n
}

// SetChangeTime records that this node was last changed at `time`.
//
// Java: Node.setChangeTime(int)
func (n *Node) SetChangeTime(time int32) {
	n.PutIntProp(Node_Prop_CHANGE_TIME, time)
}

// GetChangeTime returns the time of the last change for this node.
//
// Java: Node.getChangeTime()
func (n *Node) GetChangeTime() int32 {
	return n.GetIntProp(Node_Prop_CHANGE_TIME)
}

// Java: Node.setDeleted(boolean)
func (n *Node) SetDeleted(deleted bool) {
	n.PutBooleanProp(Node_Prop_DELETED, deleted)
}

// Java: Node.isDeleted()
func (n *Node) IsDeleted() bool {
	return n.GetBooleanProp(Node_Prop_DELETED)
}

// SetTypedefTypeProp sets, if this node represents a typedef declaration,
// the associated JSType.
//
// Java: Node.setTypedefTypeProp(JSType)
func (n *Node) SetTypedefTypeProp(type_ JSType) {
	n.PutProp(Node_Prop_TYPEDEF_TYPE, type_)
}

// GetTypedefTypeProp returns, if this node represents a typedef declaration,
// the associated JSType.
//
// Java: Node.getTypedefTypeProp()
func (n *Node) GetTypedefTypeProp() JSType {
	return node_cast[JSType](n.GetProp(Node_Prop_TYPEDEF_TYPE))
}

// SetUnusedParameter sets the value for isUnusedParameter().
//
// Java: Node.setUnusedParameter(boolean)
func (n *Node) SetUnusedParameter(unused bool) {
	n.PutBooleanProp(Node_Prop_IS_UNUSED_PARAMETER, unused)
}

// IsUnusedParameter reports whether this node is an unused function
// parameter declaration.
//
// # Set by RemoveUnusedVars
//
// Java: Node.isUnusedParameter()
func (n *Node) IsUnusedParameter() bool {
	return n.GetBooleanProp(Node_Prop_IS_UNUSED_PARAMETER)
}

// SetShorthandProperty sets the isShorthandProperty annotation.
//
// Java: Node.setShorthandProperty(boolean)
func (n *Node) SetShorthandProperty(shorthand bool) {
	n.PutBooleanProp(Node_Prop_IS_SHORTHAND_PROPERTY, shorthand)
}

// IsShorthandProperty reports whether this {x:x} property was originally
// parsed as {x}.
//
// Java: Node.isShorthandProperty()
func (n *Node) IsShorthandProperty() bool {
	return n.GetBooleanProp(Node_Prop_IS_SHORTHAND_PROPERTY)
}

// IsOptionalEs6Typed returns whether this node is an optional node in the
// ES6 Typed syntax.
//
// Java: Node.isOptionalEs6Typed()
func (n *Node) IsOptionalEs6Typed() bool {
	return n.GetBooleanProp(Node_Prop_OPT_ES6_TYPED)
}

// SetIsSyntheticBlock sets whether this is a synthetic block that should not
// be considered a real source block.
//
// Java: Node.setIsSyntheticBlock(boolean)
func (n *Node) SetIsSyntheticBlock(val bool) {
	jutil.CheckState(n.token == Token_BLOCK)
	n.PutBooleanProp(Node_Prop_SYNTHETIC, val)
}

// IsSyntheticBlock returns whether this is a synthetic block that should not
// be considered a real source block.
//
// Java: Node.isSyntheticBlock()
func (n *Node) IsSyntheticBlock() bool {
	return n.GetBooleanProp(Node_Prop_SYNTHETIC)
}

// Java: Node.setIsSynthesizedUnfulfilledNameDeclaration(boolean)
func (n *Node) SetIsSynthesizedUnfulfilledNameDeclaration(val bool) {
	jutil.CheckStateMsg1(
		n.token == Token_VAR && n.HasOneChild() && n.GetFirstChild().IsName(),

		"Expected all synthetic unfulfilled declarations to be `var <name>`, found %s",
		n)

	n.PutBooleanProp(Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION, val)
}

// Java: Node.isSynthesizedUnfulfilledNameDeclaration()
func (n *Node) IsSynthesizedUnfulfilledNameDeclaration() bool {
	return n.GetBooleanProp(Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION)
}

// SetUseStrict sets whether this node contained the "use strict" directive.
//
// Java: Node.setUseStrict(boolean)
func (n *Node) SetUseStrict(x bool) {
	n.PutBooleanProp(Node_Prop_USE_STRICT, x)
}

// IsUseStrict returns whether this node contained the "use strict"
// directive.
//
// Java: Node.isUseStrict()
func (n *Node) IsUseStrict() bool {
	return n.GetBooleanProp(Node_Prop_USE_STRICT)
}

// SetIsAddedBlock sets whether this is an added block that should not be
// considered a real source block. Eg: In "if (true) x;", the "x;" is put
// under an added block in the AST.
//
// Java: Node.setIsAddedBlock(boolean)
func (n *Node) SetIsAddedBlock(val bool) {
	n.PutBooleanProp(Node_Prop_ADDED_BLOCK, val)
}

// IsAddedBlock returns whether this is an added block that should not be
// considered a real source block.
//
// Java: Node.isAddedBlock()
func (n *Node) IsAddedBlock() bool {
	return n.GetBooleanProp(Node_Prop_ADDED_BLOCK)
}

// SetStaticMember sets whether this node is a static member node. This
// method is meaningful only on GETTER_DEF, SETTER_DEF or MEMBER_FUNCTION_DEF
// nodes contained within CLASS.
//
// Java: Node.setStaticMember(boolean)
func (n *Node) SetStaticMember(isStatic bool) {
	n.PutBooleanProp(Node_Prop_STATIC_MEMBER, isStatic)
}

// IsStaticMember returns whether this node is a static member node. This
// method is meaningful only on GETTER_DEF, SETTER_DEF or MEMBER_FUNCTION_DEF
// nodes contained within CLASS.
//
// Java: Node.isStaticMember()
func (n *Node) IsStaticMember() bool {
	return n.GetBooleanProp(Node_Prop_STATIC_MEMBER)
}

// SetIsGeneratorFunction sets whether this node is a generator node. This
// method is meaningful only on FUNCTION or MEMBER_FUNCTION_DEF nodes.
//
// Java: Node.setIsGeneratorFunction(boolean)
func (n *Node) SetIsGeneratorFunction(isGenerator bool) {
	n.PutBooleanProp(Node_Prop_GENERATOR_FN, isGenerator)
}

// IsGeneratorFunction returns whether this node is a generator function
// node.
//
// Java: Node.isGeneratorFunction()
func (n *Node) IsGeneratorFunction() bool {
	return n.GetBooleanProp(Node_Prop_GENERATOR_FN)
}

// SetGeneratorMarker sets whether this node subtree contains YIELD nodes.
//
// It's used in the translation of generators.
//
// Java: Node.setGeneratorMarker(boolean)
func (n *Node) SetGeneratorMarker(isGeneratorMarker bool) {
	n.PutBooleanProp(Node_Prop_IS_GENERATOR_MARKER, isGeneratorMarker)
}

// IsGeneratorMarker returns whether this node was marked as containing YIELD
// nodes.
//
// It's used in the translation of generators.
//
// Java: Node.isGeneratorMarker()
func (n *Node) IsGeneratorMarker() bool {
	return n.GetBooleanProp(Node_Prop_IS_GENERATOR_MARKER)
}

// SetGeneratorSafe sets the value for isGeneratorSafe().
//
// Java: Node.setGeneratorSafe(boolean)
func (n *Node) SetGeneratorSafe(isGeneratorSafe bool) {
	n.PutBooleanProp(Node_Prop_IS_GENERATOR_SAFE, isGeneratorSafe)
}

// IsGeneratorSafe is used when translating ES6 generators. If this returns
// true, this Node was generated by the compiler, and it is safe to copy this
// node to the transpiled output with no further changes.
//
// Java: Node.isGeneratorSafe()
func (n *Node) IsGeneratorSafe() bool {
	return n.GetBooleanProp(Node_Prop_IS_GENERATOR_SAFE)
}

// SetIsOptionalChainStart sets whether this node is the start of an optional
// chain. This method is meaningful only on OPTCHAIN_GETELEM,
// OPTCHAIN_GETPROP, OPTCHAIN_CALL.
//
// Java: Node.setIsOptionalChainStart(boolean)
func (n *Node) SetIsOptionalChainStart(isOptionalChainStart bool) {
	jutil.CheckStateMsg(
		!isOptionalChainStart || n.IsOptChainGetElem() || n.IsOptChainGetProp() || n.IsOptChainCall(),
		"cannot make a non-optional node the start of an optional chain.")
	n.PutBooleanProp(Node_Prop_START_OF_OPT_CHAIN, isOptionalChainStart)
}

// IsOptionalChainStart returns whether this node is an optional chaining
// node.
//
// Java: Node.isOptionalChainStart()
func (n *Node) IsOptionalChainStart() bool {
	return n.GetBooleanProp(Node_Prop_START_OF_OPT_CHAIN)
}

// SetIsArrowFunction sets whether this node is a arrow function node. This
// method is meaningful only on FUNCTION.
//
// Java: Node.setIsArrowFunction(boolean)
func (n *Node) SetIsArrowFunction(isArrow bool) {
	jutil.CheckState(n.IsFunction())
	n.PutBooleanProp(Node_Prop_ARROW_FN, isArrow)
}

// IsArrowFunction returns whether this node is a arrow function node.
//
// Java: Node.isArrowFunction()
func (n *Node) IsArrowFunction() bool {
	return n.IsFunction() && n.GetBooleanProp(Node_Prop_ARROW_FN)
}

// SetIsAsyncFunction sets whether this node is an async function node. This
// method is meaningful only on FUNCTION.
//
// Java: Node.setIsAsyncFunction(boolean)
func (n *Node) SetIsAsyncFunction(isAsync bool) {
	jutil.CheckState(n.IsFunction())
	n.PutBooleanProp(Node_Prop_ASYNC_FN, isAsync)
}

// IsAsyncFunction returns whether this is an async function node.
//
// Java: Node.isAsyncFunction()
func (n *Node) IsAsyncFunction() bool {
	return n.IsFunction() && n.GetBooleanProp(Node_Prop_ASYNC_FN)
}

// IsAsyncGeneratorFunction returns whether this is an async generator
// function node.
//
// Java: Node.isAsyncGeneratorFunction()
func (n *Node) IsAsyncGeneratorFunction() bool {
	return n.IsAsyncFunction() && n.IsGeneratorFunction()
}

// SetYieldAll sets whether this node is a generator node. This method is
// meaningful only on FUNCTION or MEMBER_FUNCTION_DEF nodes.
//
// Java: Node.setYieldAll(boolean)
func (n *Node) SetYieldAll(isGenerator bool) {
	n.PutBooleanProp(Node_Prop_YIELD_ALL, isGenerator)
}

// IsYieldAll returns whether this node is a generator node. This method is
// meaningful only on FUNCTION or MEMBER_FUNCTION_DEF nodes.
//
// Java: Node.isYieldAll()
func (n *Node) IsYieldAll() bool {
	return n.GetBooleanProp(Node_Prop_YIELD_ALL)
}

// SetTrailingComma indicates that there was a trailing comma in this list.
//
// Java: Node.setTrailingComma(boolean)
func (n *Node) SetTrailingComma(hasTrailingComma bool) {
	n.PutBooleanProp(Node_Prop_TRAILING_COMMA, hasTrailingComma)
}

// HasTrailingComma returns true if there was a trailing comma in the orginal
// code.
//
// Java: Node.hasTrailingComma()
func (n *Node) HasTrailingComma() bool {
	return n.GetBooleanProp(Node_Prop_TRAILING_COMMA)
}

// SetSideEffectFlags marks this function or constructor call's side effect
// flags. This property is only meaningful for CALL and NEW nodes.
//
// Java: Node.setSideEffectFlags(int)
func (n *Node) SetSideEffectFlags(flags int32) {
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
func (n *Node) SetSideEffectFlagsObject(flags *Node_SideEffectFlags) {
	n.SetSideEffectFlags(flags.ValueOf())
}

// GetSideEffectFlags returns the side effects flags for this node.
//
// Java: Node.getSideEffectFlags()
func (n *Node) GetSideEffectFlags() int32 {
	// Int props default to 0, but we want the default for side-effect flags to
	// be all 1s. Therefore, we invert the value returned here. This is correct
	// for non-defaults because we also invert when setting the flags.
	return ^n.GetIntProp(Node_Prop_SIDE_EFFECT_FLAGS) & Node_SideEffectFlags_USED_BITS_MASK
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
func (n *Node) IsOnlyModifiesThisCall() bool {
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
func (n *Node) IsOnlyModifiesArgumentsCall() bool {
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
func (n *Node) IsNoSideEffectsCall() bool {
	return n.GetSideEffectFlags() == Node_SideEffectFlags_NO_SIDE_EFFECTS
}

// MayMutateArguments returns true if this is a new/call that may mutate its
// arguments.
//
// Java: Node.mayMutateArguments()
func (n *Node) MayMutateArguments() bool {
	return Node_AllBitsSet(n.GetSideEffectFlags(), Node_SideEffectFlags_MUTATES_ARGUMENTS)
}

// MayMutateGlobalStateOrThrow returns true if this is a new/call that may
// mutate global state or throw.
//
// Java: Node.mayMutateGlobalStateOrThrow()
func (n *Node) MayMutateGlobalStateOrThrow() bool {
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
func (n *Node) GetConstantVarFlags() int32 {
	return n.GetIntProp(Node_Prop_CONSTANT_VAR_FLAGS)
}

// Java: Node.setConstantVarFlag(int, boolean)
func (n *Node) SetConstantVarFlag(flag int32, value bool) {
	flags := n.GetConstantVarFlags()
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
func (n *Node) IsDeclaredConstantVar() bool {
	jutil.CheckStateMsg1(
		n.IsName() || n.IsImportStar(),
		"Should only be called on name or import * nodes. Found %s",
		n)

	return Node_AnyBitSet(n.GetConstantVarFlags(), Node_ConstantVarFlags_DECLARED)
}

// SetDeclaredConstantVar sets this variable to be a declared constant.
//
// See isDeclaredConstantVar for the rules.
//
// Java: Node.setDeclaredConstantVar(boolean)
func (n *Node) SetDeclaredConstantVar(value bool) {
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
func (n *Node) IsInferredConstantVar() bool {
	jutil.CheckStateMsg1(
		n.IsName() || n.IsImportStar(),
		"Should only be called on name or import * nodes. Found %s",
		n)

	return Node_AnyBitSet(n.GetConstantVarFlags(), Node_ConstantVarFlags_INFERRED)
}

// SetInferredConstantVar sets this variable to be an inferred constant.
//
// See isInferredConstantVar for the rules.
//
// Java: Node.setInferredConstantVar(boolean)
func (n *Node) SetInferredConstantVar(value bool) {
	jutil.CheckStateMsg1(
		n.IsName() || n.IsImportStar(),
		"Should only be called on name or import * nodes. Found %s",
		n)

	n.SetConstantVarFlag(Node_ConstantVarFlags_INFERRED, value)
}

// Java: Node.isQuotedStringKey()
func (n *Node) IsQuotedStringKey() bool {
	return (n.kind == node_kindString) && n.GetBooleanProp(Node_Prop_QUOTED)
}

// Java: Node.setQuotedStringKey()
func (n *Node) SetQuotedStringKey() {
	jutil.CheckStateMsg1(n.kind == node_kindString, "%s", n)
	n.PutBooleanProp(Node_Prop_QUOTED, true)
}

// Java: Node.isPrivateIdentifier()
func (n *Node) IsPrivateIdentifier() bool {
	return (n.kind == node_kindString) && n.GetBooleanProp(Node_Prop_PRIVATE_IDENTIFIER)
}

// Java: Node.setPrivateIdentifier()
func (n *Node) SetPrivateIdentifier() {
	jutil.CheckStateMsg1(n.kind == node_kindString, "%s", n)
	jutil.CheckState(strings.HasPrefix(n.GetString(), "#"))
	n.PutBooleanProp(Node_Prop_PRIVATE_IDENTIFIER, true)
}

// AST type check methods

// Java: Node.isAdd()
func (n *Node) IsAdd() bool {
	return n.token == Token_ADD
}

// Java: Node.isSub()
func (n *Node) IsSub() bool {
	return n.token == Token_SUB
}

// Java: Node.isAnd()
func (n *Node) IsAnd() bool {
	return n.token == Token_AND
}

// Java: Node.isAssignAnd()
func (n *Node) IsAssignAnd() bool {
	return n.token == Token_ASSIGN_AND
}

// Java: Node.isArrayLit()
func (n *Node) IsArrayLit() bool {
	return n.token == Token_ARRAYLIT
}

// Java: Node.isArrayPattern()
func (n *Node) IsArrayPattern() bool {
	return n.token == Token_ARRAY_PATTERN
}

// Java: Node.isAssign()
func (n *Node) IsAssign() bool {
	return n.token == Token_ASSIGN
}

// Java: Node.isAssignAdd()
func (n *Node) IsAssignAdd() bool {
	return n.token == Token_ASSIGN_ADD
}

// Java: Node.isNormalBlock()
func (n *Node) IsNormalBlock() bool {
	return n.IsBlock()
}

// Java: Node.isBlock()
func (n *Node) IsBlock() bool {
	return n.token == Token_BLOCK
}

// Java: Node.isRoot()
func (n *Node) IsRoot() bool {
	return n.token == Token_ROOT
}

// Java: Node.isAwait()
func (n *Node) IsAwait() bool {
	return n.token == Token_AWAIT
}

// Java: Node.isBigInt()
func (n *Node) IsBigInt() bool {
	return n.token == Token_BIGINT
}

// Java: Node.isBitNot()
func (n *Node) IsBitNot() bool {
	return n.token == Token_BITNOT
}

// Java: Node.isBreak()
func (n *Node) IsBreak() bool {
	return n.token == Token_BREAK
}

// Java: Node.isCall()
func (n *Node) IsCall() bool {
	return n.token == Token_CALL
}

// Java: Node.isCase()
func (n *Node) IsCase() bool {
	return n.token == Token_CASE
}

// Java: Node.isCast()
func (n *Node) IsCast() bool {
	return n.token == Token_CAST
}

// Java: Node.isCatch()
func (n *Node) IsCatch() bool {
	return n.token == Token_CATCH
}

// Java: Node.isClass()
func (n *Node) IsClass() bool {
	return n.token == Token_CLASS
}

// Java: Node.isClassMembers()
func (n *Node) IsClassMembers() bool {
	return n.token == Token_CLASS_MEMBERS
}

// Java: Node.isComma()
func (n *Node) IsComma() bool {
	return n.token == Token_COMMA
}

// Java: Node.isComputedProp()
func (n *Node) IsComputedProp() bool {
	return n.token == Token_COMPUTED_PROP
}

// Java: Node.isContinue()
func (n *Node) IsContinue() bool {
	return n.token == Token_CONTINUE
}

// Java: Node.isConst()
func (n *Node) IsConst() bool {
	return n.token == Token_CONST
}

// Java: Node.isDebugger()
func (n *Node) IsDebugger() bool {
	return n.token == Token_DEBUGGER
}

// Java: Node.isDec()
func (n *Node) IsDec() bool {
	return n.token == Token_DEC
}

// Java: Node.isDefaultCase()
func (n *Node) IsDefaultCase() bool {
	return n.token == Token_DEFAULT_CASE
}

// Java: Node.isDefaultValue()
func (n *Node) IsDefaultValue() bool {
	return n.token == Token_DEFAULT_VALUE
}

// Java: Node.isDelProp()
func (n *Node) IsDelProp() bool {
	return n.token == Token_DELPROP
}

// Java: Node.isDestructuringLhs()
func (n *Node) IsDestructuringLhs() bool {
	return n.token == Token_DESTRUCTURING_LHS
}

// Java: Node.isDestructuringPattern()
func (n *Node) IsDestructuringPattern() bool {
	return n.IsObjectPattern() || n.IsArrayPattern()
}

// Java: Node.isDo()
func (n *Node) IsDo() bool {
	return n.token == Token_DO
}

// Java: Node.isEmpty()
func (n *Node) IsEmpty() bool {
	return n.token == Token_EMPTY
}

// Java: Node.isExponent()
func (n *Node) IsExponent() bool {
	return n.token == Token_EXPONENT
}

// Java: Node.isAssignExponent()
func (n *Node) IsAssignExponent() bool {
	return n.token == Token_ASSIGN_EXPONENT
}

// Java: Node.isExport()
func (n *Node) IsExport() bool {
	return n.token == Token_EXPORT
}

// Java: Node.isExportSpec()
func (n *Node) IsExportSpec() bool {
	return n.token == Token_EXPORT_SPEC
}

// Java: Node.isExportSpecs()
func (n *Node) IsExportSpecs() bool {
	return n.token == Token_EXPORT_SPECS
}

// Java: Node.isExprResult()
func (n *Node) IsExprResult() bool {
	return n.token == Token_EXPR_RESULT
}

// Java: Node.isFalse()
func (n *Node) IsFalse() bool {
	return n.token == Token_FALSE
}

// Java: Node.isVanillaFor()
func (n *Node) IsVanillaFor() bool {
	return n.token == Token_FOR
}

// Java: Node.isForIn()
func (n *Node) IsForIn() bool {
	return n.token == Token_FOR_IN
}

// Java: Node.isForOf()
func (n *Node) IsForOf() bool {
	return n.token == Token_FOR_OF
}

// Java: Node.isForAwaitOf()
func (n *Node) IsForAwaitOf() bool {
	return n.token == Token_FOR_AWAIT_OF
}

// Java: Node.isFunction()
func (n *Node) IsFunction() bool {
	return n.token == Token_FUNCTION
}

// Java: Node.isGetterDef()
func (n *Node) IsGetterDef() bool {
	return n.token == Token_GETTER_DEF
}

// Java: Node.isGetElem()
func (n *Node) IsGetElem() bool {
	return n.token == Token_GETELEM
}

// Java: Node.isGetProp()
func (n *Node) IsGetProp() bool {
	return n.token == Token_GETPROP
}

// Java: Node.isHook()
func (n *Node) IsHook() bool {
	return n.token == Token_HOOK
}

// Java: Node.isIf()
func (n *Node) IsIf() bool {
	return n.token == Token_IF
}

// Java: Node.isImport()
func (n *Node) IsImport() bool {
	return n.token == Token_IMPORT
}

// Java: Node.isImportMeta()
func (n *Node) IsImportMeta() bool {
	return n.token == Token_IMPORT_META
}

// Java: Node.isImportStar()
func (n *Node) IsImportStar() bool {
	return n.token == Token_IMPORT_STAR
}

// Java: Node.isImportSpec()
func (n *Node) IsImportSpec() bool {
	return n.token == Token_IMPORT_SPEC
}

// Java: Node.isImportSpecs()
func (n *Node) IsImportSpecs() bool {
	return n.token == Token_IMPORT_SPECS
}

// Java: Node.isIn()
func (n *Node) IsIn() bool {
	return n.token == Token_IN
}

// Java: Node.isInc()
func (n *Node) IsInc() bool {
	return n.token == Token_INC
}

// Java: Node.isInstanceOf()
func (n *Node) IsInstanceOf() bool {
	return n.token == Token_INSTANCEOF
}

// Java: Node.isInterface()
func (n *Node) IsInterface() bool {
	return n.token == Token_INTERFACE
}

// Java: Node.isInterfaceMembers()
func (n *Node) IsInterfaceMembers() bool {
	return n.token == Token_INTERFACE_MEMBERS
}

// Java: Node.isRecordType()
func (n *Node) IsRecordType() bool {
	return n.token == Token_RECORD_TYPE
}

// Java: Node.isCallSignature()
func (n *Node) IsCallSignature() bool {
	return n.token == Token_CALL_SIGNATURE
}

// Java: Node.isIndexSignature()
func (n *Node) IsIndexSignature() bool {
	return n.token == Token_INDEX_SIGNATURE
}

// Java: Node.isLabel()
func (n *Node) IsLabel() bool {
	return n.token == Token_LABEL
}

// Java: Node.isLabelName()
func (n *Node) IsLabelName() bool {
	return n.token == Token_LABEL_NAME
}

// Java: Node.isLet()
func (n *Node) IsLet() bool {
	return n.token == Token_LET
}

// Java: Node.isMemberFunctionDef()
func (n *Node) IsMemberFunctionDef() bool {
	return n.token == Token_MEMBER_FUNCTION_DEF
}

// Java: Node.isMemberVariableDef()
func (n *Node) IsMemberVariableDef() bool {
	return n.token == Token_MEMBER_VARIABLE_DEF
}

// Java: Node.isMemberFieldDef()
func (n *Node) IsMemberFieldDef() bool {
	return n.token == Token_MEMBER_FIELD_DEF
}

// Java: Node.isComputedFieldDef()
func (n *Node) IsComputedFieldDef() bool {
	return n.token == Token_COMPUTED_FIELD_DEF
}

// Java: Node.isModuleBody()
func (n *Node) IsModuleBody() bool {
	return n.token == Token_MODULE_BODY
}

// Java: Node.isName()
func (n *Node) IsName() bool {
	return n.token == Token_NAME
}

// Java: Node.isNE()
func (n *Node) IsNE() bool {
	return n.token == Token_NE
}

// Java: Node.isSHNE()
func (n *Node) IsSHNE() bool {
	return n.token == Token_SHNE
}

// Java: Node.isEQ()
func (n *Node) IsEQ() bool {
	return n.token == Token_EQ
}

// Java: Node.isSHEQ()
func (n *Node) IsSHEQ() bool {
	return n.token == Token_SHEQ
}

// Java: Node.isNeg()
func (n *Node) IsNeg() bool {
	return n.token == Token_NEG
}

// Java: Node.isNew()
func (n *Node) IsNew() bool {
	return n.token == Token_NEW
}

// Java: Node.isNot()
func (n *Node) IsNot() bool {
	return n.token == Token_NOT
}

// Java: Node.isNull()
func (n *Node) IsNull() bool {
	return n.token == Token_NULL
}

// Java: Node.isNullishCoalesce()
func (n *Node) IsNullishCoalesce() bool {
	return n.token == Token_COALESCE
}

// Java: Node.isAssignNullishCoalesce()
func (n *Node) IsAssignNullishCoalesce() bool {
	return n.token == Token_ASSIGN_COALESCE
}

// Java: Node.isNumber()
func (n *Node) IsNumber() bool {
	return n.token == Token_NUMBER
}

// Java: Node.isObjectLit()
func (n *Node) IsObjectLit() bool {
	return n.token == Token_OBJECTLIT
}

// Java: Node.isObjectPattern()
func (n *Node) IsObjectPattern() bool {
	return n.token == Token_OBJECT_PATTERN
}

// Java: Node.isOptChainCall()
func (n *Node) IsOptChainCall() bool {
	return n.token == Token_OPTCHAIN_CALL
}

// Java: Node.isOptChainGetElem()
func (n *Node) IsOptChainGetElem() bool {
	return n.token == Token_OPTCHAIN_GETELEM
}

// Java: Node.isOptChainGetProp()
func (n *Node) IsOptChainGetProp() bool {
	return n.token == Token_OPTCHAIN_GETPROP
}

// Java: Node.isOr()
func (n *Node) IsOr() bool {
	return n.token == Token_OR
}

// Java: Node.isAssignOr()
func (n *Node) IsAssignOr() bool {
	return n.token == Token_ASSIGN_OR
}

// Java: Node.isParamList()
func (n *Node) IsParamList() bool {
	return n.token == Token_PARAM_LIST
}

// Java: Node.isPos()
func (n *Node) IsPos() bool {
	return n.token == Token_POS
}

// Java: Node.isRegExp()
func (n *Node) IsRegExp() bool {
	return n.token == Token_REGEXP
}

// Java: Node.isRest()
func (n *Node) IsRest() bool {
	return n.token == Token_ITER_REST || n.token == Token_OBJECT_REST
}

// Java: Node.isObjectRest()
func (n *Node) IsObjectRest() bool {
	return n.token == Token_OBJECT_REST
}

// Java: Node.isReturn()
func (n *Node) IsReturn() bool {
	return n.token == Token_RETURN
}

// Java: Node.isScript()
func (n *Node) IsScript() bool {
	return n.token == Token_SCRIPT
}

// Java: Node.isSetterDef()
func (n *Node) IsSetterDef() bool {
	return n.token == Token_SETTER_DEF
}

// Java: Node.isSpread()
func (n *Node) IsSpread() bool {
	return n.token == Token_ITER_SPREAD || n.token == Token_OBJECT_SPREAD
}

// Java: Node.isString()
func (n *Node) IsString() bool {
	return n.token == Token_STRINGLIT
}

// Java: Node.isStringKey()
func (n *Node) IsStringKey() bool {
	return n.token == Token_STRING_KEY
}

// Java: Node.isStringLit()
func (n *Node) IsStringLit() bool {
	return n.token == Token_STRINGLIT
}

// Java: Node.isSuper()
func (n *Node) IsSuper() bool {
	return n.token == Token_SUPER
}

// Java: Node.isSwitch()
func (n *Node) IsSwitch() bool {
	return n.token == Token_SWITCH
}

// Java: Node.isSwitchBody()
func (n *Node) IsSwitchBody() bool {
	return n.token == Token_SWITCH_BODY
}

// Java: Node.isTaggedTemplateLit()
func (n *Node) IsTaggedTemplateLit() bool {
	return n.token == Token_TAGGED_TEMPLATELIT
}

// Java: Node.isTemplateLit()
func (n *Node) IsTemplateLit() bool {
	return n.token == Token_TEMPLATELIT
}

// Java: Node.isTemplateLitString()
func (n *Node) IsTemplateLitString() bool {
	return n.token == Token_TEMPLATELIT_STRING
}

// Java: Node.isTemplateLitSub()
func (n *Node) IsTemplateLitSub() bool {
	return n.token == Token_TEMPLATELIT_SUB
}

// Java: Node.isThis()
func (n *Node) IsThis() bool {
	return n.token == Token_THIS
}

// Java: Node.isThrow()
func (n *Node) IsThrow() bool {
	return n.token == Token_THROW
}

// Java: Node.isTrue()
func (n *Node) IsTrue() bool {
	return n.token == Token_TRUE
}

// Java: Node.isTry()
func (n *Node) IsTry() bool {
	return n.token == Token_TRY
}

// Java: Node.isTypeOf()
func (n *Node) IsTypeOf() bool {
	return n.token == Token_TYPEOF
}

// Java: Node.isVar()
func (n *Node) IsVar() bool {
	return n.token == Token_VAR
}

// Java: Node.isVoid()
func (n *Node) IsVoid() bool {
	return n.token == Token_VOID
}

// Java: Node.isWhile()
func (n *Node) IsWhile() bool {
	return n.token == Token_WHILE
}

// Java: Node.isWith()
func (n *Node) IsWith() bool {
	return n.token == Token_WITH
}

// Java: Node.isYield()
func (n *Node) IsYield() bool {
	return n.token == Token_YIELD
}

// Java: Node.isDeclare()
func (n *Node) IsDeclare() bool {
	return n.token == Token_DECLARE
}
