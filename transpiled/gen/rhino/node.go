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
// Port of com/google/javascript/rhino/Node.java (generated mechanically by j2go).

package rhino

import (
	"math"
	"strconv"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// Node is backed directly by transpiled/shims/nodestore's flat 32-byte NodeData arena.
type Node nodestore.Node
type Node_NumberNode = Node
type Node_BigIntNode = Node
type Node_StringNode = Node
type Node_TemplateLiteralSubstringNode = Node

var nodeNullPtrEx any = javalang.NewNullPointerException()

type nodeExtraFields struct {
	bigint *j2gort.Ext_java_math_BigInteger
	cooked string
	raw    string
}

func allocNodeHandle(kind uint8) Node {
	return Node(nodestore.AllocNodeInit(kind))
}

func allocNodeHandleToken(kind uint8, token uint8) Node {
	return Node(nodestore.AllocNodeToken(kind, token))
}

func tokenToByte(v Token) uint8 {
	return uint8(v) - 1
}

//go:noinline
func (this_ Node) checkDetachedSlow_() {
	rec := nodestore.Node(this_).Rec()
	if rec.Parent != 0 {
		j2gort.CheckStateFmt1(false, "Has parent: %s", this_)
	}
	if rec.Next != 0 {
		j2gort.CheckStateFmt1(false, "Has next: %s", this_)
	}
	if rec.Previous != 0 {
		j2gort.CheckStateFmt1(false, "Has previous: %s", this_)
	}
}

func checkDetachedRec_(n Node) *nodestore.NodeData {
	if n == 0 {
		panic(nodeNullPtrEx)
	}
	rec := nodestore.Node(n).Rec()
	if (rec.Parent | rec.Next | rec.Previous) != 0 {
		n.checkDetachedSlow_()
	}
	return rec
}

func getNodeExtra(n Node) *nodeExtraFields {
	aux, ok := nodestore.Node(n).GetAux()
	if !ok || aux.Extra == nil {
		return nil
	}
	return aux.Extra.(*nodeExtraFields)
}

func ensureNodeExtra(n Node) *nodeExtraFields {
	aux, ok := nodestore.Node(n).GetAux()
	if ok && aux.Extra != nil {
		return aux.Extra.(*nodeExtraFields)
	}
	ex := &nodeExtraFields{cooked: j2gort.NullString, raw: j2gort.NullString}
	aux.Extra = ex
	nodestore.Node(n).SetAux(aux)
	return ex
}

func (this_ Node) nodeGet_token() Token {
	return Token(nodestore.Node(this_).Rec().Token + 1)
}
func (this_ Node) nodeSet_token(v Token) {
	nodestore.Node(this_).Rec().Token = tokenToByte(v)
}
func (this_ Node) TokenOrd_() int32 {
	return int32(nodestore.Node(this_).Rec().Token+1) - 1
}
func (this_ Node) RawToken_() uint8 {
	return nodestore.Node(this_).Rec().Token
}
func (this_ Node) TokenSwitchOrd_() int32 {
	t := nodestore.Node(this_).Rec().Token
	if t == 0xFF {
		panic(nodeNullPtrEx)
	}
	return int32(t)
}
func (this_ Node) FirstNextTok_() (Node, Node, uint8) {
	f, nx, tok := nodestore.Node(this_).FirstNextTok()
	return Node(f), Node(nx), tok
}
func (this_ Node) HasColdProps_() bool {
	return nodestore.Node(this_).HasColdProps()
}
func (this_ Node) nodeGet_first() Node {
	return Node(nodestore.Node(this_).Rec().First)
}
func (this_ Node) nodeSet_first(v Node) {
	nodestore.Node(this_).Rec().First = nodestore.Node(v)
}
func (this_ Node) nodeGet_next() Node {
	return Node(nodestore.Node(this_).Rec().Next)
}
func (this_ Node) nodeSet_next(v Node) {
	nodestore.Node(this_).Rec().Next = nodestore.Node(v)
}
func (this_ Node) nodeGet_previous() Node {
	return Node(nodestore.Node(this_).Rec().Previous)
}
func (this_ Node) nodeSet_previous(v Node) {
	nodestore.Node(this_).Rec().Previous = nodestore.Node(v)
}
func (this_ Node) nodeGet_parent() Node {
	return Node(nodestore.Node(this_).Rec().Parent)
}
func (this_ Node) nodeSet_parent(v Node) {
	nodestore.Node(this_).Rec().Parent = nodestore.Node(v)
}
func (this_ Node) nodeGet_linenoCharno() int32 {
	return nodestore.Node(this_).Rec().LinenoCharno
}
func (this_ Node) nodeSet_linenoCharno(v int32) {
	nodestore.Node(this_).Rec().LinenoCharno = v
}
func (this_ Node) nodeGet_number() float64 {
	return nodestore.Node(this_).NumberVal()
}
func (this_ Node) nodeSet_number(v float64) {
	nodestore.Node(this_).SetNumberVal(v)
}
func (this_ Node) nodeGet_sourceFile() StaticSourceFile {
	if sf := nodestore.SourceFileByHandle(nodestore.Node(this_).SourceFile()); sf != nil {
		return sf.(StaticSourceFile)
	}
	return nil
}
func (this_ Node) nodeSet_sourceFile(v StaticSourceFile) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if v == nil || v.IsJavaNil_() {
		nodestore.Node(this_).SetSourceFileObj(nil)
	} else {
		nodestore.Node(this_).SetSourceFileObj(v)
	}
}
func (this_ Node) nodeGet_propFlags() int32 {
	return int32(nodestore.Node(this_).Rec().PropFlags)
}
func (this_ Node) nodeSet_propFlags(v int32) {
	nodestore.Node(this_).Rec().PropFlags = uint32(v)
}
func (this_ Node) nodeGet_str() string {
	p := nodestore.StringID(nodestore.Node(this_).Rec().Payload)
	if p == 0 {
		return j2gort.NullString
	}
	return nodestore.String(p)
}
func (this_ Node) nodeSet_str(v string) {
	if v == j2gort.NullString {
		nodestore.Node(this_).SetStringPayload(0)
	} else {
		nodestore.Node(this_).SetStringPayload(nodestore.Intern(v))
	}
}
func (this_ Node) nodeGet_originalName() string {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	id := nodestore.Node(this_).OrigNameID()
	if id == 0 {
		return j2gort.NullString
	}
	return nodestore.String(id)
}
func (this_ Node) nodeSet_originalName(v string) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if v == j2gort.NullString {
		nodestore.Node(this_).SetOrigNameID(0)
	} else {
		nodestore.Node(this_).SetOrigNameID(nodestore.Intern(v))
	}
}
func (this_ Node) nodeGet_length() int32 {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	return nodestore.Node(this_).Length()
}
func (this_ Node) nodeSet_length(v int32) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	nodestore.Node(this_).SetLength(v)
}
func (this_ Node) nodeGet_changeTime() int32 {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	return nodestore.Node(this_).ChangeTime()
}
func (this_ Node) nodeSet_changeTime(v int32) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	nodestore.Node(this_).SetChangeTime(v)
}
func (this_ Node) nodeGet_jstypeOrColor() any {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	return nodestore.Node(this_).TypeOrColor()
}
func (this_ Node) nodeSet_jstypeOrColor(v any) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	nodestore.Node(this_).SetTypeOrColor(v)
}
func (this_ Node) nodeGet_propListHead() *Node_PropListItem {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if p := nodestore.Node(this_).PropList(); p != nil {
		return p.(*Node_PropListItem)
	}
	return nil
}
func (this_ Node) nodeSet_propListHead(v *Node_PropListItem) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if v == nil {
		if !nodestore.Node(this_).HasColdProps() {
			return
		}
		nodestore.Node(this_).SetPropList(nil)
	} else {
		nodestore.Node(this_).SetPropList(v)
	}
}
func (this_ Node) nodeGet_jsDocInfo() *JSDocInfo {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if d := nodestore.Node(this_).JSDocObj(); d != nil {
		return d.(*JSDocInfo)
	}
	return nil
}
func (this_ Node) nodeSet_jsDocInfo(v *JSDocInfo) {
	if v == nil {
		nodestore.Node(this_).ClearJSDoc()
	} else {
		nodestore.Node(this_).SetJSDocObj(v)
	}
}
func (this_ Node) nodeGet_typeBeforeCast() JSType {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if t := nodestore.Node(this_).JSTypeBeforeCastObj(); t != nil {
		return t.(JSType)
	}
	return nil
}
func (this_ Node) nodeSet_typeBeforeCast(v JSType) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if v == nil || v.IsJavaNil_() {
		nodestore.Node(this_).ClearJSTypeBeforeCast()
	} else {
		nodestore.Node(this_).SetJSTypeBeforeCastObj(v)
	}
}
func (this_ Node) nodeGet_bigint() *j2gort.Ext_java_math_BigInteger {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if ex := getNodeExtra(this_); ex != nil {
		return ex.bigint
	}
	return nil
}
func (this_ Node) nodeSet_bigint(v *j2gort.Ext_java_math_BigInteger) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	ensureNodeExtra(this_).bigint = v
}
func (this_ Node) nodeGet_cooked() string {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if ex := getNodeExtra(this_); ex != nil {
		return ex.cooked
	}
	return j2gort.NullString
}
func (this_ Node) nodeSet_cooked(v string) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	ensureNodeExtra(this_).cooked = v
}
func (this_ Node) nodeGet_raw() string {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if ex := getNodeExtra(this_); ex != nil {
		return ex.raw
	}
	return j2gort.NullString
}
func (this_ Node) nodeSet_raw(v string) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	ensureNodeExtra(this_).raw = v
}

var static_Node_propValues []Node_Prop

func StaticGet_Node_PropValues() []Node_Prop {
	clinit_Node()
	return static_Node_propValues
}
func StaticSet_Node_PropValues(v []Node_Prop) {
	clinit_Node()
	static_Node_propValues = v
}

var static_Node_INCRDECR_PROP Node_Prop

func StaticGet_Node_INCRDECR_PROP() Node_Prop {
	clinit_Node()
	return static_Node_INCRDECR_PROP
}
func StaticSet_Node_INCRDECR_PROP(v Node_Prop) {
	clinit_Node()
	static_Node_INCRDECR_PROP = v
}

var static_Node_QUOTED_PROP Node_Prop

func StaticGet_Node_QUOTED_PROP() Node_Prop {
	clinit_Node()
	return static_Node_QUOTED_PROP
}
func StaticSet_Node_QUOTED_PROP(v Node_Prop) {
	clinit_Node()
	static_Node_QUOTED_PROP = v
}

var static_Node_IS_CONSTANT_NAME Node_Prop

func StaticGet_Node_IS_CONSTANT_NAME() Node_Prop {
	clinit_Node()
	return static_Node_IS_CONSTANT_NAME
}
func StaticSet_Node_IS_CONSTANT_NAME(v Node_Prop) {
	clinit_Node()
	static_Node_IS_CONSTANT_NAME = v
}

var static_Node_IS_NAMESPACE Node_Prop

func StaticGet_Node_IS_NAMESPACE() Node_Prop {
	clinit_Node()
	return static_Node_IS_NAMESPACE
}
func StaticSet_Node_IS_NAMESPACE(v Node_Prop) {
	clinit_Node()
	static_Node_IS_NAMESPACE = v
}

var static_Node_DIRECT_EVAL Node_Prop

func StaticGet_Node_DIRECT_EVAL() Node_Prop {
	clinit_Node()
	return static_Node_DIRECT_EVAL
}
func StaticSet_Node_DIRECT_EVAL(v Node_Prop) {
	clinit_Node()
	static_Node_DIRECT_EVAL = v
}

var static_Node_FREE_CALL Node_Prop

func StaticGet_Node_FREE_CALL() Node_Prop {
	clinit_Node()
	return static_Node_FREE_CALL
}
func StaticSet_Node_FREE_CALL(v Node_Prop) {
	clinit_Node()
	static_Node_FREE_CALL = v
}

var static_Node_REFLECTED_OBJECT Node_Prop

func StaticGet_Node_REFLECTED_OBJECT() Node_Prop {
	clinit_Node()
	return static_Node_REFLECTED_OBJECT
}
func StaticSet_Node_REFLECTED_OBJECT(v Node_Prop) {
	clinit_Node()
	static_Node_REFLECTED_OBJECT = v
}

var static_Node_STATIC_MEMBER Node_Prop

func StaticGet_Node_STATIC_MEMBER() Node_Prop {
	clinit_Node()
	return static_Node_STATIC_MEMBER
}
func StaticSet_Node_STATIC_MEMBER(v Node_Prop) {
	clinit_Node()
	static_Node_STATIC_MEMBER = v
}

var static_Node_GENERATOR_FN Node_Prop

func StaticGet_Node_GENERATOR_FN() Node_Prop {
	clinit_Node()
	return static_Node_GENERATOR_FN
}
func StaticSet_Node_GENERATOR_FN(v Node_Prop) {
	clinit_Node()
	static_Node_GENERATOR_FN = v
}

var static_Node_YIELD_ALL Node_Prop

func StaticGet_Node_YIELD_ALL() Node_Prop {
	clinit_Node()
	return static_Node_YIELD_ALL
}
func StaticSet_Node_YIELD_ALL(v Node_Prop) {
	clinit_Node()
	static_Node_YIELD_ALL = v
}

var static_Node_EXPORT_DEFAULT Node_Prop

func StaticGet_Node_EXPORT_DEFAULT() Node_Prop {
	clinit_Node()
	return static_Node_EXPORT_DEFAULT
}
func StaticSet_Node_EXPORT_DEFAULT(v Node_Prop) {
	clinit_Node()
	static_Node_EXPORT_DEFAULT = v
}

var static_Node_EXPORT_ALL_FROM Node_Prop

func StaticGet_Node_EXPORT_ALL_FROM() Node_Prop {
	clinit_Node()
	return static_Node_EXPORT_ALL_FROM
}
func StaticSet_Node_EXPORT_ALL_FROM(v Node_Prop) {
	clinit_Node()
	static_Node_EXPORT_ALL_FROM = v
}

var static_Node_COMPUTED_PROP_METHOD Node_Prop

func StaticGet_Node_COMPUTED_PROP_METHOD() Node_Prop {
	clinit_Node()
	return static_Node_COMPUTED_PROP_METHOD
}
func StaticSet_Node_COMPUTED_PROP_METHOD(v Node_Prop) {
	clinit_Node()
	static_Node_COMPUTED_PROP_METHOD = v
}

var static_Node_COMPUTED_PROP_GETTER Node_Prop

func StaticGet_Node_COMPUTED_PROP_GETTER() Node_Prop {
	clinit_Node()
	return static_Node_COMPUTED_PROP_GETTER
}
func StaticSet_Node_COMPUTED_PROP_GETTER(v Node_Prop) {
	clinit_Node()
	static_Node_COMPUTED_PROP_GETTER = v
}

var static_Node_COMPUTED_PROP_SETTER Node_Prop

func StaticGet_Node_COMPUTED_PROP_SETTER() Node_Prop {
	clinit_Node()
	return static_Node_COMPUTED_PROP_SETTER
}
func StaticSet_Node_COMPUTED_PROP_SETTER(v Node_Prop) {
	clinit_Node()
	static_Node_COMPUTED_PROP_SETTER = v
}

var static_Node_COMPUTED_PROP_VARIABLE Node_Prop

func StaticGet_Node_COMPUTED_PROP_VARIABLE() Node_Prop {
	clinit_Node()
	return static_Node_COMPUTED_PROP_VARIABLE
}
func StaticSet_Node_COMPUTED_PROP_VARIABLE(v Node_Prop) {
	clinit_Node()
	static_Node_COMPUTED_PROP_VARIABLE = v
}

var static_Node_OPT_ES6_TYPED Node_Prop

func StaticGet_Node_OPT_ES6_TYPED() Node_Prop {
	clinit_Node()
	return static_Node_OPT_ES6_TYPED
}
func StaticSet_Node_OPT_ES6_TYPED(v Node_Prop) {
	clinit_Node()
	static_Node_OPT_ES6_TYPED = v
}

var static_Node_GENERIC_TYPE_LIST Node_Prop

func StaticGet_Node_GENERIC_TYPE_LIST() Node_Prop {
	clinit_Node()
	return static_Node_GENERIC_TYPE_LIST
}
func StaticSet_Node_GENERIC_TYPE_LIST(v Node_Prop) {
	clinit_Node()
	static_Node_GENERIC_TYPE_LIST = v
}

var static_Node_IMPLEMENTS Node_Prop

func StaticGet_Node_IMPLEMENTS() Node_Prop {
	clinit_Node()
	return static_Node_IMPLEMENTS
}
func StaticSet_Node_IMPLEMENTS(v Node_Prop) {
	clinit_Node()
	static_Node_IMPLEMENTS = v
}

var static_Node_CONSTRUCT_SIGNATURE Node_Prop

func StaticGet_Node_CONSTRUCT_SIGNATURE() Node_Prop {
	clinit_Node()
	return static_Node_CONSTRUCT_SIGNATURE
}
func StaticSet_Node_CONSTRUCT_SIGNATURE(v Node_Prop) {
	clinit_Node()
	static_Node_CONSTRUCT_SIGNATURE = v
}

var static_Node_ACCESS_MODIFIER Node_Prop

func StaticGet_Node_ACCESS_MODIFIER() Node_Prop {
	clinit_Node()
	return static_Node_ACCESS_MODIFIER
}
func StaticSet_Node_ACCESS_MODIFIER(v Node_Prop) {
	clinit_Node()
	static_Node_ACCESS_MODIFIER = v
}

var static_Node_GOOG_MODULE Node_Prop

func StaticGet_Node_GOOG_MODULE() Node_Prop {
	clinit_Node()
	return static_Node_GOOG_MODULE
}
func StaticSet_Node_GOOG_MODULE(v Node_Prop) {
	clinit_Node()
	static_Node_GOOG_MODULE = v
}

var static_Node_FEATURE_SET Node_Prop

func StaticGet_Node_FEATURE_SET() Node_Prop {
	clinit_Node()
	return static_Node_FEATURE_SET
}
func StaticSet_Node_FEATURE_SET(v Node_Prop) {
	clinit_Node()
	static_Node_FEATURE_SET = v
}

var static_Node_IS_TYPESCRIPT_ABSTRACT Node_Prop

func StaticGet_Node_IS_TYPESCRIPT_ABSTRACT() Node_Prop {
	clinit_Node()
	return static_Node_IS_TYPESCRIPT_ABSTRACT
}
func StaticSet_Node_IS_TYPESCRIPT_ABSTRACT(v Node_Prop) {
	clinit_Node()
	static_Node_IS_TYPESCRIPT_ABSTRACT = v
}

var static_Node_MODULE_ALIAS Node_Prop

func StaticGet_Node_MODULE_ALIAS() Node_Prop {
	clinit_Node()
	return static_Node_MODULE_ALIAS
}
func StaticSet_Node_MODULE_ALIAS(v Node_Prop) {
	clinit_Node()
	static_Node_MODULE_ALIAS = v
}

var static_Node_MODULE_EXPORT Node_Prop

func StaticGet_Node_MODULE_EXPORT() Node_Prop {
	clinit_Node()
	return static_Node_MODULE_EXPORT
}
func StaticSet_Node_MODULE_EXPORT(v Node_Prop) {
	clinit_Node()
	static_Node_MODULE_EXPORT = v
}

var static_Node_IS_SHORTHAND_PROPERTY Node_Prop

func StaticGet_Node_IS_SHORTHAND_PROPERTY() Node_Prop {
	clinit_Node()
	return static_Node_IS_SHORTHAND_PROPERTY
}
func StaticSet_Node_IS_SHORTHAND_PROPERTY(v Node_Prop) {
	clinit_Node()
	static_Node_IS_SHORTHAND_PROPERTY = v
}

var static_Node_ES6_MODULE Node_Prop

func StaticGet_Node_ES6_MODULE() Node_Prop {
	clinit_Node()
	return static_Node_ES6_MODULE
}
func StaticSet_Node_ES6_MODULE(v Node_Prop) {
	clinit_Node()
	static_Node_ES6_MODULE = v
}
func StaticGet_Node_PROP_FLAG_IS_CONSTANT_NAME() int32              { return int32(1) }
func StaticGet_Node_PROP_FLAG_QUOTED() int32                        { return int32(2) }
func StaticGet_Node_PROP_FLAG_FREE_CALL() int32                     { return int32(4) }
func StaticGet_Node_PROP_FLAG_NON_INDEXABLE() int32                 { return int32(8) }
func StaticGet_Node_PROP_FLAG_ADDED_BLOCK() int32                   { return int32(16) }
func StaticGet_Node_PROP_FLAG_TRAILING_COMMA() int32                { return int32(32) }
func StaticGet_Node_PROP_FLAG_IS_PARENTHESIZED() int32              { return int32(64) }
func StaticGet_Node_PROP_FLAG_ARROW_FN() int32                      { return int32(128) }
func StaticGet_Node_PROP_FLAG_INCRDECR() int32                      { return int32(256) }
func StaticGet_Node_PROP_FLAG_SYNTHETIC() int32                     { return int32(512) }
func StaticGet_Node_PROP_FLAG_IS_SHORTHAND_PROPERTY() int32         { return int32(1024) }
func StaticGet_Node_PROP_FLAG_IS_NAMESPACE() int32                  { return int32(2048) }
func StaticGet_Node_PROP_FLAG_USE_STRICT() int32                    { return int32(4096) }
func StaticGet_Node_PROP_FLAG_STATIC_MEMBER() int32                 { return int32(8192) }
func StaticGet_Node_PROP_FLAG_GENERATOR_FN() int32                  { return int32(16384) }
func StaticGet_Node_PROP_FLAG_ASYNC_FN() int32                      { return int32(32768) }
func StaticGet_Node_PROP_FLAG_COLOR_FROM_CAST() int32               { return int32(65536) }
func StaticGet_Node_PROP_FLAG_GOOG_MODULE() int32                   { return int32(131072) }
func StaticGet_Node_PROP_FLAG_ES6_MODULE() int32                    { return int32(262144) }
func StaticGet_Node_PROP_FLAG_DELETED() int32                       { return int32(524288) }
func StaticGet_Node_PROP_FLAG_MODULE_ALIAS() int32                  { return int32(1048576) }
func StaticGet_Node_PROP_FLAG_IS_UNUSED_PARAMETER() int32           { return int32(2097152) }
func StaticGet_Node_PROP_FLAG_MODULE_EXPORT() int32                 { return int32(4194304) }
func StaticGet_Node_PROP_FLAG_START_OF_OPT_CHAIN() int32            { return int32(8388608) }
func StaticGet_Node_PROP_FLAG_IS_IN_CLOSURE_UNAWARE_SUBTREE() int32 { return int32(16777216) }
func StaticGet_Node_PROP_FLAG_COMPUTED_PROP_METHOD() int32          { return int32(33554432) }
func StaticGet_Node_INLINE_BOOL_COUNT() int32                       { return int32(26) }
func StaticGet_Node_INLINE_BOOL_ALL_MASK() int32                    { return int32(67108863) }
func StaticGet_Node_CONST_VAR_SHIFT() int32                         { return int32(26) }
func StaticGet_Node_CONST_VAR_MAX() int32                           { return int32(3) }
func StaticGet_Node_CONST_VAR_BITS_MASK() int32                     { return int32(201326592) }
func StaticGet_Node_SIDE_EFFECT_SHIFT() int32                       { return int32(28) }
func StaticGet_Node_SIDE_EFFECT_MAX() int32                         { return int32(15) }
func StaticGet_Node_SIDE_EFFECT_BITS_MASK() int32                   { return int32(-268435456) }
func StaticGet_Node_VALIDATED_INLINE_PROPS_MASK() int32             { return int32(209814208) }
func StaticGet_Node_PROP_EQUALITY_INLINE_MASK() int32               { return int32(42000774) }

var static_Node_INLINE_BIT_TO_PROP []Node_Prop

func StaticGet_Node_INLINE_BIT_TO_PROP() []Node_Prop {
	clinit_Node()
	return static_Node_INLINE_BIT_TO_PROP
}
func StaticSet_Node_INLINE_BIT_TO_PROP(v []Node_Prop) {
	clinit_Node()
	static_Node_INLINE_BIT_TO_PROP = v
}

var static_Node_INLINE_BOOL_MASKS []int32

func StaticGet_Node_INLINE_BOOL_MASKS() []int32 {
	clinit_Node()
	return static_Node_INLINE_BOOL_MASKS
}
func StaticSet_Node_INLINE_BOOL_MASKS(v []int32) {
	clinit_Node()
	static_Node_INLINE_BOOL_MASKS = v
}

var static_Node_STATIC_SINGLE_TRUE_ITEMS []*Node_PropListItem

func StaticGet_Node_STATIC_SINGLE_TRUE_ITEMS() []*Node_PropListItem {
	clinit_Node()
	return static_Node_STATIC_SINGLE_TRUE_ITEMS
}
func StaticSet_Node_STATIC_SINGLE_TRUE_ITEMS(v []*Node_PropListItem) {
	clinit_Node()
	static_Node_STATIC_SINGLE_TRUE_ITEMS = v
}

var static_Node_STATIC_INLINE_BOOL_ITEMS []*Node_PropListItem

func StaticGet_Node_STATIC_INLINE_BOOL_ITEMS() []*Node_PropListItem {
	clinit_Node()
	return static_Node_STATIC_INLINE_BOOL_ITEMS
}
func StaticSet_Node_STATIC_INLINE_BOOL_ITEMS(v []*Node_PropListItem) {
	clinit_Node()
	static_Node_STATIC_INLINE_BOOL_ITEMS = v
}

var static_Node_STATIC_CONST_VAR_ITEMS []*Node_PropListItem

func StaticGet_Node_STATIC_CONST_VAR_ITEMS() []*Node_PropListItem {
	clinit_Node()
	return static_Node_STATIC_CONST_VAR_ITEMS
}
func StaticSet_Node_STATIC_CONST_VAR_ITEMS(v []*Node_PropListItem) {
	clinit_Node()
	static_Node_STATIC_CONST_VAR_ITEMS = v
}

var static_Node_STATIC_SIDE_EFFECT_ITEMS []*Node_PropListItem

func StaticGet_Node_STATIC_SIDE_EFFECT_ITEMS() []*Node_PropListItem {
	clinit_Node()
	return static_Node_STATIC_SIDE_EFFECT_ITEMS
}
func StaticSet_Node_STATIC_SIDE_EFFECT_ITEMS(v []*Node_PropListItem) {
	clinit_Node()
	static_Node_STATIC_SIDE_EFFECT_ITEMS = v
}
func StaticGet_Node_CHARNO_BITS() int32       { return int32(12) }
func StaticGet_Node_MAX_COLUMN_NUMBER() int32 { return int32(4095) }

var static_Node_PROP_EQUALITY_ORDINAL_MASK int64

func StaticGet_Node_PROP_EQUALITY_ORDINAL_MASK() int64 {
	clinit_Node()
	return static_Node_PROP_EQUALITY_ORDINAL_MASK
}
func StaticSet_Node_PROP_EQUALITY_ORDINAL_MASK(v int64) {
	clinit_Node()
	static_Node_PROP_EQUALITY_ORDINAL_MASK = v
}

var static_Node_POW31 []int32

func StaticGet_Node_POW31() []int32 {
	clinit_Node()
	return static_Node_POW31
}
func StaticSet_Node_POW31(v []int32) {
	clinit_Node()
	static_Node_POW31 = v
}

var clinitDone_Node uint32
var clinitState_Node j2gort.ClinitState

func clinit_Node() {
	if atomic.LoadUint32(&clinitDone_Node) != 0 {
		return
	}
	clinitSlow_Node()
}

func clinitSlow_Node() {
	if !j2gort.EnterClinit(&clinitState_Node, &clinitDone_Node) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Node, &clinitDone_Node)
	static_Node_propValues = Node_Prop_Values()
	static_Node_INCRDECR_PROP = StaticGet_Node_Prop_INCRDECR()
	static_Node_QUOTED_PROP = StaticGet_Node_Prop_QUOTED()
	static_Node_IS_CONSTANT_NAME = StaticGet_Node_Prop_IS_CONSTANT_NAME()
	static_Node_IS_NAMESPACE = StaticGet_Node_Prop_IS_NAMESPACE()
	static_Node_DIRECT_EVAL = StaticGet_Node_Prop_DIRECT_EVAL()
	static_Node_FREE_CALL = StaticGet_Node_Prop_FREE_CALL()
	static_Node_REFLECTED_OBJECT = StaticGet_Node_Prop_REFLECTED_OBJECT()
	static_Node_STATIC_MEMBER = StaticGet_Node_Prop_STATIC_MEMBER()
	static_Node_GENERATOR_FN = StaticGet_Node_Prop_GENERATOR_FN()
	static_Node_YIELD_ALL = StaticGet_Node_Prop_YIELD_ALL()
	static_Node_EXPORT_DEFAULT = StaticGet_Node_Prop_EXPORT_DEFAULT()
	static_Node_EXPORT_ALL_FROM = StaticGet_Node_Prop_EXPORT_ALL_FROM()
	static_Node_COMPUTED_PROP_METHOD = StaticGet_Node_Prop_COMPUTED_PROP_METHOD()
	static_Node_COMPUTED_PROP_GETTER = StaticGet_Node_Prop_COMPUTED_PROP_GETTER()
	static_Node_COMPUTED_PROP_SETTER = StaticGet_Node_Prop_COMPUTED_PROP_SETTER()
	static_Node_COMPUTED_PROP_VARIABLE = StaticGet_Node_Prop_COMPUTED_PROP_VARIABLE()
	static_Node_OPT_ES6_TYPED = StaticGet_Node_Prop_OPT_ES6_TYPED()
	static_Node_GENERIC_TYPE_LIST = StaticGet_Node_Prop_GENERIC_TYPE()
	static_Node_IMPLEMENTS = StaticGet_Node_Prop_IMPLEMENTS()
	static_Node_CONSTRUCT_SIGNATURE = StaticGet_Node_Prop_CONSTRUCT_SIGNATURE()
	static_Node_ACCESS_MODIFIER = StaticGet_Node_Prop_ACCESS_MODIFIER()
	static_Node_GOOG_MODULE = StaticGet_Node_Prop_GOOG_MODULE()
	static_Node_FEATURE_SET = StaticGet_Node_Prop_FEATURE_SET()
	static_Node_IS_TYPESCRIPT_ABSTRACT = StaticGet_Node_Prop_IS_TYPESCRIPT_ABSTRACT()
	static_Node_MODULE_ALIAS = StaticGet_Node_Prop_MODULE_ALIAS()
	static_Node_MODULE_EXPORT = StaticGet_Node_Prop_MODULE_EXPORT()
	static_Node_IS_SHORTHAND_PROPERTY = StaticGet_Node_Prop_IS_SHORTHAND_PROPERTY()
	static_Node_ES6_MODULE = StaticGet_Node_Prop_ES6_MODULE()
	static_Node_INLINE_BIT_TO_PROP = make([]Node_Prop, int32(26))
	static_Node_INLINE_BOOL_MASKS = make([]int32, int32(len(StaticGet_Node_PropValues())))
	static_Node_STATIC_SINGLE_TRUE_ITEMS = make([]*Node_PropListItem, int32(len(StaticGet_Node_PropValues())))
	static_Node_STATIC_INLINE_BOOL_ITEMS = make([]*Node_PropListItem, int32(26))
	static_Node_STATIC_CONST_VAR_ITEMS = make([]*Node_PropListItem, (int32(3) + int32(1)))
	static_Node_STATIC_SIDE_EFFECT_ITEMS = make([]*Node_PropListItem, (int32(15) + int32(1)))
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(StaticGet_Node_PropValues())); i = i + 1 {
			{
				a_1 := StaticGet_Node_STATIC_SINGLE_TRUE_ITEMS()
				i_2 := i
				a_1[i_2] = NewNode_PropListItem(int8(i), int32(1), nil, nil)
			}
		}
	}
	Node_RegisterInlineBoolProp(int32(0), StaticGet_Node_Prop_IS_CONSTANT_NAME())
	Node_RegisterInlineBoolProp(int32(1), StaticGet_Node_Prop_QUOTED())
	Node_RegisterInlineBoolProp(int32(2), StaticGet_Node_Prop_FREE_CALL())
	Node_RegisterInlineBoolProp(int32(3), StaticGet_Node_Prop_NON_INDEXABLE())
	Node_RegisterInlineBoolProp(int32(4), StaticGet_Node_Prop_ADDED_BLOCK())
	Node_RegisterInlineBoolProp(int32(5), StaticGet_Node_Prop_TRAILING_COMMA())
	Node_RegisterInlineBoolProp(int32(6), StaticGet_Node_Prop_IS_PARENTHESIZED())
	Node_RegisterInlineBoolProp(int32(7), StaticGet_Node_Prop_ARROW_FN())
	Node_RegisterInlineBoolProp(int32(8), StaticGet_Node_Prop_INCRDECR())
	Node_RegisterInlineBoolProp(int32(9), StaticGet_Node_Prop_SYNTHETIC())
	Node_RegisterInlineBoolProp(int32(10), StaticGet_Node_Prop_IS_SHORTHAND_PROPERTY())
	Node_RegisterInlineBoolProp(int32(11), StaticGet_Node_Prop_IS_NAMESPACE())
	Node_RegisterInlineBoolProp(int32(12), StaticGet_Node_Prop_USE_STRICT())
	Node_RegisterInlineBoolProp(int32(13), StaticGet_Node_Prop_STATIC_MEMBER())
	Node_RegisterInlineBoolProp(int32(14), StaticGet_Node_Prop_GENERATOR_FN())
	Node_RegisterInlineBoolProp(int32(15), StaticGet_Node_Prop_ASYNC_FN())
	Node_RegisterInlineBoolProp(int32(16), StaticGet_Node_Prop_COLOR_FROM_CAST())
	Node_RegisterInlineBoolProp(int32(17), StaticGet_Node_Prop_GOOG_MODULE())
	Node_RegisterInlineBoolProp(int32(18), StaticGet_Node_Prop_ES6_MODULE())
	Node_RegisterInlineBoolProp(int32(19), StaticGet_Node_Prop_DELETED())
	Node_RegisterInlineBoolProp(int32(20), StaticGet_Node_Prop_MODULE_ALIAS())
	Node_RegisterInlineBoolProp(int32(21), StaticGet_Node_Prop_IS_UNUSED_PARAMETER())
	Node_RegisterInlineBoolProp(int32(22), StaticGet_Node_Prop_MODULE_EXPORT())
	Node_RegisterInlineBoolProp(int32(23), StaticGet_Node_Prop_START_OF_OPT_CHAIN())
	Node_RegisterInlineBoolProp(int32(24), StaticGet_Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE())
	Node_RegisterInlineBoolProp(int32(25), StaticGet_Node_Prop_COMPUTED_PROP_METHOD())
	{
		var v int32 = int32(1)
		_ = v
		for ; v <= int32(3); v = v + 1 {
			{
				a_3 := StaticGet_Node_STATIC_CONST_VAR_ITEMS()
				i_4 := v
				a_3[i_4] = Node_PropListItem_IntItem(int8((StaticGet_Node_Prop_CONSTANT_VAR_FLAGS()).Ordinal()), v, nil)
			}
		}
	}
	{
		var v_2 int32 = int32(1)
		_ = v_2
		for ; v_2 <= int32(15); v_2 = v_2 + 1 {
			{
				a_5 := StaticGet_Node_STATIC_SIDE_EFFECT_ITEMS()
				i_6 := v_2
				a_5[i_6] = Node_PropListItem_IntItem(int8((StaticGet_Node_Prop_SIDE_EFFECT_FLAGS()).Ordinal()), v_2, nil)
			}
		}
	}
	static_Node_PROP_EQUALITY_ORDINAL_MASK = Node_ComputePropEqualityOrdinalMask()
	static_Node_POW31 = Node_InitPow31()
}

// Java method: Node.getNonJSDocCommentString()
func (this_ Node) GetNonJSDocCommentString() string {
	if this_.GetProp(StaticGet_Node_Prop_NON_JSDOC_COMMENT()) == nil {
		return ""
	}
	return (j2gort.CastPtr[NonJSDocComment](this_.GetProp(StaticGet_Node_Prop_NON_JSDOC_COMMENT()))).GetCommentString()
}

// Java method: Node.getNonJSDocComment()
func (this_ Node) GetNonJSDocComment() *NonJSDocComment {
	return j2gort.CastPtr[NonJSDocComment](this_.GetProp(StaticGet_Node_Prop_NON_JSDOC_COMMENT()))
}

// Java method: Node.getTrailingNonJSDocCommentString()
func (this_ Node) GetTrailingNonJSDocCommentString() string {
	if this_.GetProp(StaticGet_Node_Prop_TRAILING_NON_JSDOC_COMMENT()) == nil {
		return ""
	}
	return (j2gort.CastPtr[NonJSDocComment](this_.GetProp(StaticGet_Node_Prop_TRAILING_NON_JSDOC_COMMENT()))).GetCommentString()
}

// Java method: Node.getTrailingNonJSDocComment()
func (this_ Node) GetTrailingNonJSDocComment() *NonJSDocComment {
	return j2gort.CastPtr[NonJSDocComment](this_.GetProp(StaticGet_Node_Prop_TRAILING_NON_JSDOC_COMMENT()))
}

// Java method: Node.setNonJSDocComment(com.google.javascript.rhino.NonJSDocComment comment)
func (this_ Node) SetNonJSDocComment(comment *NonJSDocComment) Node {
	_ = comment
	this_.PutProp(StaticGet_Node_Prop_NON_JSDOC_COMMENT(), j2gort.PtrToAny(comment))
	return this_
}

// Java method: Node.setTrailingNonJSDocComment(com.google.javascript.rhino.NonJSDocComment comment)
func (this_ Node) SetTrailingNonJSDocComment(comment *NonJSDocComment) Node {
	_ = comment
	this_.PutProp(StaticGet_Node_Prop_TRAILING_NON_JSDOC_COMMENT(), j2gort.PtrToAny(comment))
	return this_
}

// Java method: Node.setIsParenthesized(boolean b)
func (this_ Node) SetIsParenthesized(b bool) {
	_ = b
	if !(IR_MayBeExpression(this_)) {
		j2gort.CheckState(false)
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_IS_PARENTHESIZED(), b)
}

// Java method: Node.getIsParenthesized()
func (this_ Node) GetIsParenthesized() bool {
	return ((this_.nodeGet_propFlags() & int32(64)) != int32(0))
}

// Java static method: Node.registerInlineBoolProp(int bit, com.google.javascript.rhino.Node.Prop prop)
func Node_RegisterInlineBoolProp(bit int32, prop Node_Prop) {
	clinit_Node()
	_ = bit
	_ = prop
	{
		a_7 := StaticGet_Node_INLINE_BIT_TO_PROP()
		i_8 := bit
		a_7[i_8] = prop
	}
	{
		a_9 := StaticGet_Node_INLINE_BOOL_MASKS()
		i_10 := (prop).Ordinal()
		a_9[i_10] = javalang.Shl32(int32(1), int(bit))
	}
	{
		a_11 := StaticGet_Node_STATIC_INLINE_BOOL_ITEMS()
		i_12 := bit
		a_11[i_12] = StaticGet_Node_STATIC_SINGLE_TRUE_ITEMS()[(prop).Ordinal()]
	}
}

func NewNode_1(token Token) Node {
	clinit_Node()
	tok_ := tokenToByte(token)
	raw_, rec_ := nodestore.AllocNodeFastRec()
	this_ := Node(raw_)
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	return this_
}

func initNodeNode_1(this_ Node, token Token) {
	_ = token
	this_.nodeSet_token(token)
}

func NewNode_2(token Token, child Node) Node {
	clinit_Node()
	_ = token
	_ = child
	tok_ := tokenToByte(token)
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.First = nodestore.Node(child)
	crec_ := checkDetachedRec_(child)
	crec_.Previous = nodestore.Node(child)
	crec_.Parent = raw_
	return Node(raw_)
}

func NewNode_3(token Token, left Node, right Node) Node {
	clinit_Node()
	_ = token
	_ = left
	_ = right
	tok_ := tokenToByte(token)
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.First = nodestore.Node(left)
	lrec_ := checkDetachedRec_(left)
	lrec_.Next = nodestore.Node(right)
	lrec_.Previous = nodestore.Node(right)
	lrec_.Parent = raw_
	rrec_ := checkDetachedRec_(right)
	rrec_.Previous = nodestore.Node(left)
	rrec_.Parent = raw_
	return Node(raw_)
}

func NewNode_4(token Token, left Node, mid Node, right Node) Node {
	clinit_Node()
	_ = token
	_ = left
	_ = mid
	_ = right
	tok_ := tokenToByte(token)
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.First = nodestore.Node(left)
	lrec_ := checkDetachedRec_(left)
	lrec_.Next = nodestore.Node(mid)
	lrec_.Previous = nodestore.Node(right)
	lrec_.Parent = raw_
	mrec_ := checkDetachedRec_(mid)
	mrec_.Next = nodestore.Node(right)
	mrec_.Previous = nodestore.Node(left)
	mrec_.Parent = raw_
	rrec_ := checkDetachedRec_(right)
	rrec_.Previous = nodestore.Node(mid)
	rrec_.Parent = raw_
	return Node(raw_)
}

// Java static method: Node.newNumber(double number)
func Node_NewNumber(number float64) Node {
	clinit_Node()
	_ = number
	return NewNode_NumberNode(number)
}

// Java static method: Node.newBigInt(java.math.BigInteger bigint)
func Node_NewBigInt(bigint *j2gort.Ext_java_math_BigInteger) Node {
	clinit_Node()
	_ = bigint
	return NewNode_BigIntNode(bigint)
}

// Java static method: Node.newString(java.lang.String str)
func Node_NewString_1(str string) Node {
	clinit_Node()
	_ = str
	return NewNode_StringNode_2(StaticGet_Token_STRINGLIT(), str)
}

// Java static method: Node.newString(com.google.javascript.rhino.Token token, java.lang.String str)
func Node_NewString_2(token Token, str string) Node {
	clinit_Node()
	_ = token
	_ = str
	return NewNode_StringNode_2(token, str)
}

// Java static method: Node.newString(com.google.javascript.rhino.Token token, com.google.javascript.rhino.RhinoStringPool.LazyInternedStringList stringPool, int offset)
func Node_NewString_3(token Token, stringPool *RhinoStringPool_LazyInternedStringList, offset int32) Node {
	clinit_Node()
	_ = token
	_ = stringPool
	_ = offset
	return NewNode_StringNode_3(token, stringPool, offset)
}

// Java static method: Node.newTemplateLitString(java.lang.String cooked, java.lang.String raw)
func Node_NewTemplateLitString_2(cooked string, raw string) Node {
	clinit_Node()
	_ = cooked
	_ = raw
	return NewNode_TemplateLiteralSubstringNode_2(cooked, raw)
}

// Java static method: Node.newTemplateLitString(com.google.javascript.rhino.RhinoStringPool.LazyInternedStringList stringPool, int cookedOffsetOrNegativeOne, int rawOffset)
func Node_NewTemplateLitString_3(stringPool *RhinoStringPool_LazyInternedStringList, cookedOffsetOrNegativeOne int32, rawOffset int32) Node {
	clinit_Node()
	_ = stringPool
	_ = cookedOffsetOrNegativeOne
	_ = rawOffset
	return NewNode_TemplateLiteralSubstringNode_3(stringPool, cookedOffsetOrNegativeOne, rawOffset)
}

// Java method: Node.getToken()
func (this_ Node) GetToken() Token {
	return this_.nodeGet_token()
}

// Java method: Node.setToken(com.google.javascript.rhino.Token token)
func (this_ Node) SetToken(token Token) {
	_ = token
	this_.nodeSet_token(token)
}

// Java method: Node.hasChildren()
func (this_ Node) HasChildren() bool {
	return (this_.nodeGet_first() != 0)
}

// Java method: Node.getOnlyChild()
func (this_ Node) GetOnlyChild() Node {
	if !(this_.HasOneChild()) {
		j2gort.CheckState(false)
	}
	return this_.nodeGet_first()
}

// Java method: Node.getFirstChild()
func (this_ Node) GetFirstChild() Node {
	return this_.nodeGet_first()
}

// Java method: Node.getFirstFirstChild()
func (this_ Node) GetFirstFirstChild() Node {
	return this_.nodeGet_first().nodeGet_first()
}

// Java method: Node.getSecondChild()
func (this_ Node) GetSecondChild() Node {
	return this_.nodeGet_first().nodeGet_next()
}

// Java method: Node.getLastChild()
func (this_ Node) GetLastChild() Node {
	if this_.nodeGet_first() != 0 {
		return this_.nodeGet_first().nodeGet_previous()
	} else {
		return Node(0)
	}
}

// Java method: Node.getNext()
func (this_ Node) GetNext() Node {
	return this_.nodeGet_next()
}

// Java method: Node.getPrevious()
func (this_ Node) GetPrevious() Node {
	if this_ == this_.nodeGet_parent().nodeGet_first() {
		return Node(0)
	} else {
		return this_.nodeGet_previous()
	}
}

// Java method: Node.setClosureUnawareShadow(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node shadowRoot)
func (this_ Node) SetClosureUnawareShadow(shadowRoot Node) {
	_ = shadowRoot
	if !(this_.nodeGet_first() == 0) {
		j2gort.CheckStateStr(false, "Cannot set shadow root on a node with children")
	}
	if !(this_.RawToken_() == uint8(32)) {
		j2gort.CheckStateStr(false, "Only NAME nodes can be used as shadows")
	}
	this_.PutProp(StaticGet_Node_Prop_CLOSURE_UNAWARE_SHADOW(), j2gort.BoxNode(shadowRoot))
}

// Java method: Node.getClosureUnawareShadow()
func (this_ Node) GetClosureUnawareShadow() Node {
	if this_.RawToken_() != uint8(32) {
		return Node(0)
	}
	return j2gort.CastNode[Node](this_.GetProp(StaticGet_Node_Prop_CLOSURE_UNAWARE_SHADOW()))
}

// Java method: Node.setIsInClosureUnawareSubtree(boolean value)
func (this_ Node) SetIsInClosureUnawareSubtree(value bool) {
	_ = value
	if ((((j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil) || (this_.nodeGet_jsDocInfo() != nil)) || !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil))) || (this_.nodeGet_propFlags() != int32(0))) || (this_.nodeGet_changeTime() != int32(0))) || (this_.nodeGet_propListHead() != nil) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("Cannot set IS_IN_CLOSURE_UNAWARE_SUBTREE property on a node which proplist is not a single SOURCE_FILE prop: %s %s", []any{j2gort.BoxNode(this_), j2gort.BoxString(this_.GetPropListDebugString())})))))
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE(), value)
}

// Java method: Node.getIsInClosureUnawareSubtree()
func (this_ Node) GetIsInClosureUnawareSubtree() bool {
	return ((this_.nodeGet_propFlags() & int32(16777216)) != int32(0))
}

// Java method: Node.getChildAtIndex(int i)
func (this_ Node) GetChildAtIndex(i int32) Node {
	_ = i
	var n Node = this_.nodeGet_first()
	_ = n
	for i > int32(0) {
		n = n.nodeGet_next()
		i = i - 1
	}
	return n
}

// Java method: Node.getIndexOfChild(com.google.javascript.rhino.Node child)
func (this_ Node) GetIndexOfChild(child Node) int32 {
	_ = child
	var n Node = this_.nodeGet_first()
	_ = n
	var i int32 = int32(0)
	_ = i
	for n != 0 {
		if child == n {
			return i
		}
		n = n.nodeGet_next()
		i = i + 1
	}
	return int32(-1)
}

// Java method: Node.addChildToFront(com.google.javascript.rhino.Node child)
func (this_ Node) AddChildToFront(child Node) {
	_ = child
	if !(child.nodeGet_parent() == 0) {
		j2gort.CheckArgument(false)
	}
	if !(child.nodeGet_next() == 0) {
		j2gort.CheckArgument(false)
	}
	if !(child.nodeGet_previous() == 0) {
		j2gort.CheckArgument(false)
	}
	{
		r_13 := child
		r_13.nodeSet_parent(this_)
	}
	{
		r_14 := child
		r_14.nodeSet_next(this_.nodeGet_first())
	}
	if this_.nodeGet_first() == 0 {
		{
			r_15 := child
			r_15.nodeSet_previous(child)
		}
	} else {
		var last Node = this_.nodeGet_first().nodeGet_previous()
		_ = last
		{
			r_16 := child
			r_16.nodeSet_previous(last)
		}
		{
			r_17 := child
			r_17.nodeSet_next(this_.nodeGet_first())
		}
		{
			r_18 := this_.nodeGet_first()
			r_18.nodeSet_previous(child)
		}
	}
	this_.nodeSet_first(child)
}

// Java method: Node.addChildToBack(com.google.javascript.rhino.Node child)
func (this_ Node) AddChildToBack(child Node) {
	_ = child
	if !(child.nodeGet_parent() == 0) {
		j2gort.CheckArgumentFmt(false, "Cannot add already-owned child node.\nChild: %s\nExisting parent: %s\nNew parent: %s", j2gort.BoxNode(child), j2gort.BoxNode(child.nodeGet_parent()), j2gort.BoxNode(this_))
	}
	if !(child.nodeGet_next() == 0) {
		j2gort.CheckArgument(false)
	}
	if !(child.nodeGet_previous() == 0) {
		j2gort.CheckArgument(false)
	}
	if this_.nodeGet_first() == 0 {
		{
			r_19 := child
			r_19.nodeSet_previous(child)
		}
		this_.nodeSet_first(child)
	} else {
		var last Node = this_.nodeGet_first().nodeGet_previous()
		_ = last
		{
			r_20 := last
			r_20.nodeSet_next(child)
		}
		{
			r_21 := child
			r_21.nodeSet_previous(last)
		}
		{
			r_22 := this_.nodeGet_first()
			r_22.nodeSet_previous(child)
		}
	}
	{
		r_23 := child
		r_23.nodeSet_parent(this_)
	}
}

// Java method: Node.addChildrenToFront(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node children)
func (this_ Node) AddChildrenToFront(children Node) {
	_ = children
	if children == 0 {
		return
	}
	j2gort.CheckNotNullNode(children.nodeGet_previous())
	{
		var child Node = children
		_ = child
		for ; child != 0; child = child.nodeGet_next() {
			if !(child.nodeGet_parent() == 0) {
				j2gort.CheckArgument(false)
			}
			{
				r_24 := child
				r_24.nodeSet_parent(this_)
			}
		}
	}
	var lastSib Node = children.nodeGet_previous()
	_ = lastSib
	if this_.nodeGet_first() != 0 {
		var last Node = this_.nodeGet_first().nodeGet_previous()
		_ = last
		{
			r_25 := children
			r_25.nodeSet_previous(last)
		}
		{
			r_26 := lastSib
			r_26.nodeSet_next(this_.nodeGet_first())
		}
		{
			r_27 := this_.nodeGet_first()
			r_27.nodeSet_previous(lastSib)
		}
	}
	this_.nodeSet_first(children)
}

// Java method: Node.addChildrenToBack(com.google.javascript.rhino.Node children)
func (this_ Node) AddChildrenToBack(children Node) {
	_ = children
	this_.AddChildrenAfter(children, this_.GetLastChild())
}

// Java method: Node.insertAfter(com.google.javascript.rhino.Node existing)
func (this_ Node) InsertAfter(existing Node) {
	_ = existing
	existing.priv_Node_checkAttached_8412()
	this_.priv_Node_checkDetached_8413()
	var existingParent Node = existing.nodeGet_parent()
	_ = existingParent
	var existingNext Node = existing.nodeGet_next()
	_ = existingNext
	this_.nodeSet_parent(existingParent)
	{
		r_28 := existing
		r_28.nodeSet_next(this_)
	}
	this_.nodeSet_previous(existing)
	if existingNext == 0 {
		{
			r_29 := existingParent.nodeGet_first()
			r_29.nodeSet_previous(this_)
		}
	} else {
		{
			r_30 := existingNext
			r_30.nodeSet_previous(this_)
		}
		this_.nodeSet_next(existingNext)
	}
}

// Java method: Node.insertBefore(com.google.javascript.rhino.Node existing)
func (this_ Node) InsertBefore(existing Node) {
	_ = existing
	existing.priv_Node_checkAttached_8412()
	this_.priv_Node_checkDetached_8413()
	var existingParent Node = existing.nodeGet_parent()
	_ = existingParent
	var existingPrevious Node = existing.nodeGet_previous()
	_ = existingPrevious
	this_.nodeSet_parent(existingParent)
	this_.nodeSet_next(existing)
	{
		r_31 := existing
		r_31.nodeSet_previous(this_)
	}
	this_.nodeSet_previous(existingPrevious)
	if existingPrevious.nodeGet_next() == 0 {
		{
			r_32 := existingParent
			r_32.nodeSet_first(this_)
		}
	} else {
		{
			r_33 := existingPrevious
			r_33.nodeSet_next(this_)
		}
	}
}

// Java method: Node.addChildrenAfter(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node children, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node)
func (this_ Node) AddChildrenAfter(children Node, node Node) {
	_ = children
	_ = node
	if children == 0 {
		return
	}
	if !((node == 0) || (node.nodeGet_parent() == this_)) {
		j2gort.CheckArgument(false)
	}
	j2gort.CheckNotNullNode(children.nodeGet_previous())
	if node == 0 {
		this_.AddChildrenToFront(children)
		return
	}
	{
		var child Node = children
		_ = child
		for ; child != 0; child = child.nodeGet_next() {
			if !(child.nodeGet_parent() == 0) {
				j2gort.CheckArgument(false)
			}
			{
				r_34 := child
				r_34.nodeSet_parent(this_)
			}
		}
	}
	var lastSibling Node = children.nodeGet_previous()
	_ = lastSibling
	var nodeAfter Node = node.nodeGet_next()
	_ = nodeAfter
	{
		r_35 := lastSibling
		r_35.nodeSet_next(nodeAfter)
	}
	if nodeAfter == 0 {
		{
			r_36 := this_.nodeGet_first()
			r_36.nodeSet_previous(lastSibling)
		}
	} else {
		{
			r_37 := nodeAfter
			r_37.nodeSet_previous(lastSibling)
		}
	}
	{
		r_38 := node
		r_38.nodeSet_next(children)
	}
	{
		r_39 := children
		r_39.nodeSet_previous(node)
	}
}

// Java method: Node.replaceWith(com.google.javascript.rhino.Node replacement)
func (this_ Node) ReplaceWith(replacement Node) {
	_ = replacement
	this_.priv_Node_checkAttached_8412()
	replacement.priv_Node_checkDetached_8413()
	var existingParent Node = this_.nodeGet_parent()
	_ = existingParent
	var existingNext Node = this_.nodeGet_next()
	_ = existingNext
	var existingPrevious Node = this_.nodeGet_previous()
	_ = existingPrevious
	replacement.SrcrefIfMissing(this_)
	this_.nodeSet_parent(Node(0))
	{
		r_40 := replacement
		r_40.nodeSet_parent(existingParent)
	}
	this_.nodeSet_previous(Node(0))
	{
		r_41 := replacement
		r_41.nodeSet_previous(existingPrevious)
	}
	if existingPrevious.nodeGet_next() == 0 {
		{
			r_42 := existingParent
			r_42.nodeSet_first(replacement)
		}
	} else {
		{
			r_43 := existingPrevious
			r_43.nodeSet_next(replacement)
		}
	}
	if existingNext == 0 {
		{
			r_44 := existingParent.nodeGet_first()
			r_44.nodeSet_previous(replacement)
		}
	} else {
		this_.nodeSet_next(Node(0))
		{
			r_45 := existingNext
			r_45.nodeSet_previous(replacement)
		}
		{
			r_46 := replacement
			r_46.nodeSet_next(existingNext)
		}
	}
}

// Java method: Node.detach()
func (this_ Node) Detach() Node {
	this_.priv_Node_checkAttached_8412()
	var existingParent Node = this_.nodeGet_parent()
	_ = existingParent
	var existingNext Node = this_.nodeGet_next()
	_ = existingNext
	var existingPrevious Node = this_.nodeGet_previous()
	_ = existingPrevious
	this_.nodeSet_parent(Node(0))
	if existingNext == 0 {
		{
			r_47 := existingParent.nodeGet_first()
			r_47.nodeSet_previous(existingPrevious)
		}
	} else {
		this_.nodeSet_next(Node(0))
		{
			r_48 := existingNext
			r_48.nodeSet_previous(existingPrevious)
		}
	}
	this_.nodeSet_previous(Node(0))
	if existingPrevious.nodeGet_next() == 0 {
		{
			r_49 := existingParent
			r_49.nodeSet_first(existingNext)
		}
	} else {
		{
			r_50 := existingPrevious
			r_50.nodeSet_next(existingNext)
		}
	}
	return this_
}

// Java method: Node.checkAttached()
func (this_ Node) priv_Node_checkAttached_8412() {
	if !(this_.nodeGet_parent() != 0) {
		j2gort.CheckStateFmt(false, "Has no parent: %s", j2gort.BoxNode(this_))
	}
}

// Java method: Node.checkDetached()
func (this_ Node) priv_Node_checkDetached_8413() {
	if !(this_.nodeGet_parent() == 0) {
		j2gort.CheckStateFmt(false, "Has parent: %s", j2gort.BoxNode(this_))
	}
	if !(this_.nodeGet_next() == 0) {
		j2gort.CheckStateFmt(false, "Has next: %s", j2gort.BoxNode(this_))
	}
	if !(this_.nodeGet_previous() == 0) {
		j2gort.CheckStateFmt(false, "Has previous: %s", j2gort.BoxNode(this_))
	}
}

// Java method: Node.removeFirstChild()
func (this_ Node) RemoveFirstChild() Node {
	var child Node = this_.nodeGet_first()
	_ = child
	if child != 0 {
		child.Detach()
	}
	return child
}

// Java method: Node.removeChildren()
func (this_ Node) RemoveChildren() Node {
	var children Node = this_.nodeGet_first()
	_ = children
	{
		var child Node = this_.nodeGet_first()
		_ = child
		for ; child != 0; child = child.nodeGet_next() {
			{
				r_51 := child
				r_51.nodeSet_parent(Node(0))
			}
		}
	}
	this_.nodeSet_first(Node(0))
	return children
}

// Java method: Node.detachChildren()
func (this_ Node) DetachChildren() {
	{
		var child Node = this_.nodeGet_first()
		_ = child
		for child != 0 {
			var nextChild Node = child.nodeGet_next()
			_ = nextChild
			{
				r_52 := child
				r_52.nodeSet_parent(Node(0))
			}
			{
				r_53 := child
				r_53.nodeSet_next(Node(0))
			}
			{
				r_54 := child
				r_54.nodeSet_previous(Node(0))
			}
			child = nextChild
		}
	}
	this_.nodeSet_first(Node(0))
}

// Java method: Node.lookupColdProperty(com.google.javascript.rhino.Node.Prop prop)
func (this_ Node) priv_Node_lookupColdProperty_8417(prop Node_Prop) *Node_PropListItem {
	_ = prop
	var x *Node_PropListItem = this_.nodeGet_propListHead()
	_ = x
	if x == nil {
		return nil
	}
	var propType int8 = int8((prop).Ordinal())
	_ = propType
	for (x != nil) && (int32(propType) != int32(x.propType)) {
		x = x.next
	}
	return x
}

// Java method: Node.lookupProperty(com.google.javascript.rhino.Node.Prop prop)
func (this_ Node) LookupProperty(prop Node_Prop) *Node_PropListItem {
	_ = prop
	if prop == StaticGet_Node_Prop_SOURCE_FILE() {
		if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
			return Node_PropListItem_ObjectItem(int8((StaticGet_Node_Prop_SOURCE_FILE()).Ordinal()), this_.nodeGet_sourceFile(), nil)
		}
		return this_.priv_Node_lookupColdProperty_8417(prop)
	}
	if prop == StaticGet_Node_Prop_JSDOC_INFO() {
		if this_.nodeGet_jsDocInfo() != nil {
			return Node_PropListItem_ObjectItem(int8((StaticGet_Node_Prop_JSDOC_INFO()).Ordinal()), j2gort.PtrToAny(this_.nodeGet_jsDocInfo()), nil)
		}
		return this_.priv_Node_lookupColdProperty_8417(prop)
	}
	if prop == StaticGet_Node_Prop_TYPE_BEFORE_CAST() {
		if !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
			return Node_PropListItem_ObjectItem(int8((StaticGet_Node_Prop_TYPE_BEFORE_CAST()).Ordinal()), this_.nodeGet_typeBeforeCast(), nil)
		}
		return this_.priv_Node_lookupColdProperty_8417(prop)
	}
	var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(prop).Ordinal()]
	_ = mask
	if mask != int32(0) {
		if (this_.nodeGet_propFlags() & mask) != int32(0) {
			if this_.nodeGet_propListHead() != nil {
				var cold *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(prop)
				_ = cold
				if cold != nil {
					return cold
				}
			}
			return StaticGet_Node_STATIC_INLINE_BOOL_ITEMS()[j2gort.IntTrailingZeros(mask)]
		}
		return this_.priv_Node_lookupColdProperty_8417(prop)
	}
	if prop == StaticGet_Node_Prop_CONSTANT_VAR_FLAGS() {
		var v int32 = ((javalang.Ushr32(this_.nodeGet_propFlags(), int(int32(26)))) & int32(3))
		_ = v
		if v != int32(0) {
			return StaticGet_Node_STATIC_CONST_VAR_ITEMS()[v]
		}
		return this_.priv_Node_lookupColdProperty_8417(prop)
	}
	if prop == StaticGet_Node_Prop_SIDE_EFFECT_FLAGS() {
		var v_2 int32 = ((javalang.Ushr32(this_.nodeGet_propFlags(), int(int32(28)))) & int32(15))
		_ = v_2
		if v_2 != int32(0) {
			return StaticGet_Node_STATIC_SIDE_EFFECT_ITEMS()[v_2]
		}
		return this_.priv_Node_lookupColdProperty_8417(prop)
	}
	if prop == StaticGet_Node_Prop_CHANGE_TIME() {
		if this_.nodeGet_changeTime() != int32(0) {
			return Node_PropListItem_IntItem(int8((StaticGet_Node_Prop_CHANGE_TIME()).Ordinal()), this_.nodeGet_changeTime(), nil)
		}
		return this_.priv_Node_lookupColdProperty_8417(prop)
	}
	return this_.priv_Node_lookupColdProperty_8417(prop)
}

// Java method: Node.clonePropsFrom(com.google.javascript.rhino.Node other)
func (this_ Node) ClonePropsFrom(other Node) Node {
	if !nodestore.Node(this_).ClonePropsFrom(nodestore.Node(other)) {
		j2gort.CheckStateStr(false, "Node has existing properties.")
	}
	return this_
}

// Java method: Node.hasPropertyViolations()
func (this_ Node) HasPropertyViolations() bool {
	var flags int32 = this_.nodeGet_propFlags()
	_ = flags
	var head *Node_PropListItem = this_.nodeGet_propListHead()
	_ = head
	if this_.RawToken_() == uint8(89) {
		return (((((!(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) || (this_.nodeGet_jsDocInfo() != nil)) || !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil))) || (flags != int32(0))) || (this_.nodeGet_changeTime() != int32(0))) || (head != nil))
	}
	var inlineFlags int32 = (flags & int32(209814208))
	_ = inlineFlags
	if (inlineFlags == int32(0)) && (head == nil) {
		return false
	}
	if (head != nil) && (inlineFlags != int32(0)) {
		{
			var x *Node_PropListItem = head
			_ = x
			for ; x != nil; x = x.next {
				inlineFlags = (inlineFlags & ^(StaticGet_Node_INLINE_BOOL_MASKS()[int32(x.propType)]))
			}
		}
	}
	if inlineFlags != int32(0) {
		if ((inlineFlags & int32(64)) != int32(0)) && !(IR_MayBeExpression(this_)) {
			return true
		}
		if ((inlineFlags & int32(128)) != int32(0)) && !(this_.IsFunction()) {
			return true
		}
		if ((inlineFlags & int32(32768)) != int32(0)) && !(this_.IsFunction()) {
			return true
		}
		if ((inlineFlags & int32(512)) != int32(0)) && !(this_.IsBlock()) {
			return true
		}
		if ((inlineFlags & int32(65536)) != int32(0)) && (this_.GetColor() == nil) {
			return true
		}
		if ((inlineFlags & int32(8388608)) != int32(0)) && !((this_.IsOptChainCall() || this_.IsOptChainGetElem()) || this_.IsOptChainGetProp()) {
			return true
		}
		if ((inlineFlags & int32(201326592)) != int32(0)) && !(this_.IsName() || this_.IsImportStar()) {
			return true
		}
	}
	{
		var propListItem *Node_PropListItem = head
		_ = propListItem
		for ; propListItem != nil; propListItem = propListItem.next {
			var prop Node_Prop = StaticGet_Node_PropValues()[int32(propListItem.propType)]
			_ = prop
			if !((prop).Ordinal() == int32(propListItem.propType)) {
				j2gort.CheckStateFmt(false, "ordinal doesn't match: %s", j2gort.BoxNullableVal(prop))
			}
			switch (prop).Ordinal() {
			case int32(0) /* Prop.IS_PARENTHESIZED */ :
				if !(IR_MayBeExpression(this_)) {
					return true
				}
			case int32(20) /* Prop.ARROW_FN */ :
				if !(this_.IsFunction()) {
					return true
				}
			case int32(21) /* Prop.ASYNC_FN */ :
				if !(this_.IsFunction()) {
					return true
				}
			case int32(6) /* Prop.SYNTHETIC */ :
				if !(this_.IsBlock()) {
					return true
				}
			case int32(34) /* Prop.COLOR_FROM_CAST */ :
				if this_.GetColor() == nil {
					return true
				}
			case int32(51) /* Prop.START_OF_OPT_CHAIN */ :
				if !((this_.IsOptChainCall() || this_.IsOptChainGetElem()) || this_.IsOptChainGetProp()) {
					return true
				}
			case int32(25) /* Prop.CONSTANT_VAR_FLAGS */ :
				if !(this_.IsName() || this_.IsImportStar()) {
					return true
				}
			case int32(53) /* Prop.SYNTHESIZED_UNFULFILLED_NAME_DECLARATION */ :
				if !(this_.IsVar()) {
					return true
				}
			case int32(54) /* Prop.CLOSURE_UNAWARE_SHADOW */ :
				var shadowProp *Node_PropListItem = this_.LookupProperty(StaticGet_Node_Prop_CLOSURE_UNAWARE_SHADOW())
				_ = shadowProp
				if ((shadowProp == nil) || !(shadowProp.IsObjectItem())) || !(j2gort.InstanceOf[Node](shadowProp.GetObjectValue())) {
					return true
				}
			default:
			}
		}
	}
	return false
}

// Java method: Node.validateProperties(java.util.function.Consumer<java.lang.String> violationMessageConsumer)
func (this_ Node) ValidateProperties(violationMessageConsumer j2gort.Ext_java_util_function_Consumer) {
	_ = violationMessageConsumer
	var flags int32 = this_.nodeGet_propFlags()
	_ = flags
	var head *Node_PropListItem = this_.nodeGet_propListHead()
	_ = head
	if ((((j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil) && (this_.nodeGet_jsDocInfo() == nil)) && j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) && (flags == int32(0))) && (this_.nodeGet_changeTime() == int32(0))) && (head == nil) {
		return
	}
	if this_.RawToken_() == uint8(89) {
		violationMessageConsumer.Accept("ROOT has properties")
	}
	var inlineFlags int32 = flags
	_ = inlineFlags
	if (head != nil) && ((inlineFlags & int32(209814208)) != int32(0)) {
		{
			var x *Node_PropListItem = head
			_ = x
			for ; x != nil; x = x.next {
				inlineFlags = (inlineFlags & ^(StaticGet_Node_INLINE_BOOL_MASKS()[int32(x.propType)]))
			}
		}
	}
	if (inlineFlags & int32(209814208)) != int32(0) {
		if ((inlineFlags & int32(64)) != int32(0)) && !(IR_MayBeExpression(this_)) {
			violationMessageConsumer.Accept("non-expression is parenthesized")
		}
		if ((inlineFlags & int32(128)) != int32(0)) && !(this_.IsFunction()) {
			violationMessageConsumer.Accept("invalid ARROW_FN prop")
		}
		if ((inlineFlags & int32(32768)) != int32(0)) && !(this_.IsFunction()) {
			violationMessageConsumer.Accept("invalid ASYNC_FN prop")
		}
		if ((inlineFlags & int32(512)) != int32(0)) && !(this_.IsBlock()) {
			violationMessageConsumer.Accept("invalid SYNTHETIC prop")
		}
		if ((inlineFlags & int32(65536)) != int32(0)) && (this_.GetColor() == nil) {
			violationMessageConsumer.Accept("COLOR_FROM_CAST with no Color")
		}
		if ((inlineFlags & int32(8388608)) != int32(0)) && !((this_.IsOptChainCall() || this_.IsOptChainGetElem()) || this_.IsOptChainGetProp()) {
			violationMessageConsumer.Accept("START_OF_OPT_CHAIN on non-optional Node")
		}
		if ((inlineFlags & int32(201326592)) != int32(0)) && !(this_.IsName() || this_.IsImportStar()) {
			violationMessageConsumer.Accept("invalid CONST_VAR_FLAGS")
		}
	}
	{
		var propListItem *Node_PropListItem = head
		_ = propListItem
		for ; propListItem != nil; propListItem = propListItem.next {
			var prop Node_Prop = StaticGet_Node_PropValues()[int32(propListItem.propType)]
			_ = prop
			if !((prop).Ordinal() == int32(propListItem.propType)) {
				j2gort.CheckStateFmt(false, "ordinal doesn't match: %s", j2gort.BoxNullableVal(prop))
			}
			switch (prop).Ordinal() {
			case int32(0) /* Prop.IS_PARENTHESIZED */ :
				if !(IR_MayBeExpression(this_)) {
					violationMessageConsumer.Accept("non-expression is parenthesized")
				}
			case int32(20) /* Prop.ARROW_FN */ :
				if !(this_.IsFunction()) {
					violationMessageConsumer.Accept("invalid ARROW_FN prop")
				}
			case int32(21) /* Prop.ASYNC_FN */ :
				if !(this_.IsFunction()) {
					violationMessageConsumer.Accept("invalid ASYNC_FN prop")
				}
			case int32(6) /* Prop.SYNTHETIC */ :
				if !(this_.IsBlock()) {
					violationMessageConsumer.Accept("invalid SYNTHETIC prop")
				}
			case int32(34) /* Prop.COLOR_FROM_CAST */ :
				if this_.GetColor() == nil {
					violationMessageConsumer.Accept("COLOR_FROM_CAST with no Color")
				}
			case int32(51) /* Prop.START_OF_OPT_CHAIN */ :
				if !((this_.IsOptChainCall() || this_.IsOptChainGetElem()) || this_.IsOptChainGetProp()) {
					violationMessageConsumer.Accept("START_OF_OPT_CHAIN on non-optional Node")
				}
			case int32(25) /* Prop.CONSTANT_VAR_FLAGS */ :
				if !(this_.IsName() || this_.IsImportStar()) {
					violationMessageConsumer.Accept("invalid CONST_VAR_FLAGS")
				}
			case int32(53) /* Prop.SYNTHESIZED_UNFULFILLED_NAME_DECLARATION */ :
				if !(this_.IsVar()) {
					violationMessageConsumer.Accept("Expected all synthetic unfulfilled declarations to be `var <name>`")
				}
			case int32(54) /* Prop.CLOSURE_UNAWARE_SHADOW */ :
				var shadowProp *Node_PropListItem = this_.LookupProperty(StaticGet_Node_Prop_CLOSURE_UNAWARE_SHADOW())
				_ = shadowProp
				if ((shadowProp == nil) || !(shadowProp.IsObjectItem())) || !(j2gort.InstanceOf[Node](shadowProp.GetObjectValue())) {
					violationMessageConsumer.Accept("CLOSURE_UNAWARE_SHADOW property must point to a Node")
				}
			default:
			}
		}
	}
}

// Java static method: Node.rebuildListWithoutProp(com.google.javascript.rhino.Node.@org.jspecify.annotations.Nullable PropListItem item, com.google.javascript.rhino.Node.Prop prop)
func Node_RebuildListWithoutProp(item *Node_PropListItem, prop Node_Prop) *Node_PropListItem {
	clinit_Node()
	_ = item
	_ = prop
	if item == nil {
		return nil
	} else if int32(item.propType) == (prop).Ordinal() {
		return item.next
	} else {
		var result *Node_PropListItem = Node_RebuildListWithoutProp(item.next, prop)
		_ = result
		if result == item.next {
			return item
		} else {
			return item.Chain(result)
		}
	}
}

// Java method: Node.removeColdProp(com.google.javascript.rhino.Node.Prop prop)
func (this_ Node) priv_Node_removeColdProp_8422(prop Node_Prop) {
	_ = prop
	var head *Node_PropListItem = this_.nodeGet_propListHead()
	_ = head
	if head != nil {
		var nextHead *Node_PropListItem = Node_RebuildListWithoutProp(head, prop)
		_ = nextHead
		if nextHead != head {
			this_.nodeSet_propListHead(nextHead)
		}
	}
}

// Java method: Node.getProp(com.google.javascript.rhino.Node.Prop propType)
func (this_ Node) GetProp(propType Node_Prop) any {
	_ = propType
	if propType == StaticGet_Node_Prop_SOURCE_FILE() {
		if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
			return this_.nodeGet_sourceFile()
		}
		var item *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(propType)
		_ = item
		if item == nil {
			return nil
		} else {
			return item.GetObjectValue()
		}
	}
	if propType == StaticGet_Node_Prop_JSDOC_INFO() {
		if this_.nodeGet_jsDocInfo() != nil {
			return j2gort.PtrToAny(this_.nodeGet_jsDocInfo())
		}
		var item_2 *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(propType)
		_ = item_2
		if item_2 == nil {
			return nil
		} else {
			return item_2.GetObjectValue()
		}
	}
	if propType == StaticGet_Node_Prop_TYPE_BEFORE_CAST() {
		if !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
			return this_.nodeGet_typeBeforeCast()
		}
		var item_3 *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(propType)
		_ = item_3
		if item_3 == nil {
			return nil
		} else {
			return item_3.GetObjectValue()
		}
	}
	var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(propType).Ordinal()]
	_ = mask
	if mask != int32(0) {
		if (this_.nodeGet_propFlags() & mask) != int32(0) {
			panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
		}
	} else if propType == StaticGet_Node_Prop_CONSTANT_VAR_FLAGS() {
		if (this_.nodeGet_propFlags() & int32(201326592)) != int32(0) {
			panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
		}
	} else if propType == StaticGet_Node_Prop_SIDE_EFFECT_FLAGS() {
		if (this_.nodeGet_propFlags() & int32(-268435456)) != int32(0) {
			panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
		}
	} else if propType == StaticGet_Node_Prop_CHANGE_TIME() {
		if this_.nodeGet_changeTime() != int32(0) {
			panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
		}
	}
	var item_4 *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(propType)
	_ = item_4
	if item_4 == nil {
		return nil
	}
	return item_4.GetObjectValue()
}

// Java method: Node.getBooleanProp(com.google.javascript.rhino.Node.Prop propType)
func (this_ Node) GetBooleanProp(propType Node_Prop) bool {
	_ = propType
	var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(propType).Ordinal()]
	_ = mask
	if mask != int32(0) {
		return ((this_.nodeGet_propFlags() & mask) != int32(0))
	}
	return this_.priv_Node_getColdBooleanProp_8425(propType)
}

// Java method: Node.getColdBooleanProp(com.google.javascript.rhino.Node.Prop propType)
func (this_ Node) priv_Node_getColdBooleanProp_8425(propType Node_Prop) bool {
	_ = propType
	if (((propType == StaticGet_Node_Prop_SOURCE_FILE()) && !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil))) || ((propType == StaticGet_Node_Prop_JSDOC_INFO()) && (this_.nodeGet_jsDocInfo() != nil))) || ((propType == StaticGet_Node_Prop_TYPE_BEFORE_CAST()) && !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil))) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
	}
	if (propType == StaticGet_Node_Prop_CHANGE_TIME()) && (this_.nodeGet_changeTime() != int32(0)) {
		return true
	}
	var item *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(propType)
	_ = item
	return ((item != nil) && (item.GetIntValue() != int32(0)))
}

// Java method: Node.getColdIntProp(com.google.javascript.rhino.Node.Prop propType)
func (this_ Node) priv_Node_getColdIntProp_8426(propType Node_Prop) int32 {
	_ = propType
	var item *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(propType)
	_ = item
	if item == nil {
		return int32(0)
	} else {
		return item.GetIntValue()
	}
}

// Java method: Node.getIntProp(com.google.javascript.rhino.Node.Prop propType)
func (this_ Node) priv_Node_getIntProp_8427(propType Node_Prop) int32 {
	_ = propType
	var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(propType).Ordinal()]
	_ = mask
	if mask != int32(0) {
		if (this_.nodeGet_propFlags() & mask) != int32(0) {
			if this_.nodeGet_propListHead() != nil {
				var cold *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(propType)
				_ = cold
				if cold != nil {
					return cold.GetIntValue()
				}
			}
			return int32(1)
		}
		return this_.priv_Node_getColdIntProp_8426(propType)
	}
	if propType == StaticGet_Node_Prop_CONSTANT_VAR_FLAGS() {
		var v int32 = ((javalang.Ushr32(this_.nodeGet_propFlags(), int(int32(26)))) & int32(3))
		_ = v
		if v != int32(0) {
			return v
		}
		return this_.priv_Node_getColdIntProp_8426(propType)
	}
	if propType == StaticGet_Node_Prop_SIDE_EFFECT_FLAGS() {
		var v_2 int32 = ((javalang.Ushr32(this_.nodeGet_propFlags(), int(int32(28)))) & int32(15))
		_ = v_2
		if v_2 != int32(0) {
			return v_2
		}
		return this_.priv_Node_getColdIntProp_8426(propType)
	}
	if propType == StaticGet_Node_Prop_CHANGE_TIME() {
		if this_.nodeGet_changeTime() != int32(0) {
			return this_.nodeGet_changeTime()
		}
		if this_.nodeGet_propListHead() != nil {
			return this_.priv_Node_getColdIntProp_8426(propType)
		} else {
			return int32(0)
		}
	}
	if (propType == StaticGet_Node_Prop_SOURCE_FILE()) && !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
	}
	if (propType == StaticGet_Node_Prop_JSDOC_INFO()) && (this_.nodeGet_jsDocInfo() != nil) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
	}
	if (propType == StaticGet_Node_Prop_TYPE_BEFORE_CAST()) && !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
	}
	return this_.priv_Node_getColdIntProp_8426(propType)
}

// Java method: Node.putProp(com.google.javascript.rhino.Node.Prop prop, java.lang.@org.jspecify.annotations.Nullable Object value)
func (this_ Node) PutProp(prop Node_Prop, value any) {
	_ = prop
	_ = value
	if prop == StaticGet_Node_Prop_SOURCE_FILE() {
		if value == nil {
			this_.nodeSet_sourceFile(nil)
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		if j2gort.InstanceOf[StaticSourceFile](value) {
			this_.nodeSet_sourceFile(j2gort.AsIface[StaticSourceFile](value))
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		this_.nodeSet_sourceFile(nil)
	} else if prop == StaticGet_Node_Prop_JSDOC_INFO() {
		if value == nil {
			this_.nodeSet_jsDocInfo(nil)
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		var info *JSDocInfo
		_ = info
		var ok_55 bool
		info, ok_55 = j2gort.TryCast[*JSDocInfo](value)
		if ok_55 {
			this_.nodeSet_jsDocInfo(info)
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		this_.nodeSet_jsDocInfo(nil)
	} else if prop == StaticGet_Node_Prop_TYPE_BEFORE_CAST() {
		if value == nil {
			this_.nodeSet_typeBeforeCast(nil)
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		var jsType JSType
		_ = jsType
		var ok_56 bool
		jsType, ok_56 = j2gort.TryCast[JSType](value)
		if ok_56 {
			this_.nodeSet_typeBeforeCast(jsType)
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		this_.nodeSet_typeBeforeCast(nil)
	} else {
		var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(prop).Ordinal()]
		_ = mask
		if mask != int32(0) {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(mask)))
		} else if prop == StaticGet_Node_Prop_CONSTANT_VAR_FLAGS() {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(int32(201326592))))
		} else if prop == StaticGet_Node_Prop_SIDE_EFFECT_FLAGS() {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(int32(-268435456))))
		} else if prop == StaticGet_Node_Prop_CHANGE_TIME() {
			this_.nodeSet_changeTime(int32(0))
		}
	}
	var oldHead *Node_PropListItem = this_.nodeGet_propListHead()
	_ = oldHead
	var head *Node_PropListItem = oldHead
	_ = head
	if head == nil {
		if value == nil {
			return
		}
	} else {
		if value != nil {
			var existing *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(prop)
			_ = existing
			if ((existing != nil) && existing.IsObjectItem()) && j2gort.RefEquals(existing.objectValue, value) {
				return
			}
		}
		head = Node_RebuildListWithoutProp(head, prop)
	}
	if value != nil {
		head = Node_PropListItem_ObjectItem(int8((prop).Ordinal()), value, head)
	}
	if head != oldHead {
		this_.nodeSet_propListHead(head)
	}
}

// Java method: Node.putBooleanProp(com.google.javascript.rhino.Node.Prop propType, boolean value)
func (this_ Node) PutBooleanProp(propType Node_Prop, value bool) {
	_ = propType
	_ = value
	var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(propType).Ordinal()]
	_ = mask
	if mask != int32(0) {
		if value {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() | mask))
		} else {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(mask)))
		}
		this_.priv_Node_removeColdProp_8422(propType)
		return
	}
	this_.PutIntProp(propType, func() int32 {
		if value {
			return int32(1)
		}
		return int32(0)
	}())
}

// Java method: Node.putIntProp(com.google.javascript.rhino.Node.Prop prop, int value)
func (this_ Node) PutIntProp(prop Node_Prop, value int32) {
	_ = prop
	_ = value
	var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(prop).Ordinal()]
	_ = mask
	if mask != int32(0) {
		if value == int32(0) {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(mask)))
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		this_.nodeSet_propFlags((this_.nodeGet_propFlags() | mask))
		if value == int32(1) {
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
	} else if prop == StaticGet_Node_Prop_CONSTANT_VAR_FLAGS() {
		if (value & ^(int32(3))) == int32(0) {
			this_.nodeSet_propFlags(((this_.nodeGet_propFlags() & ^(int32(201326592))) | (javalang.Shl32(value, int(int32(26))))))
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(int32(201326592))))
	} else if prop == StaticGet_Node_Prop_SIDE_EFFECT_FLAGS() {
		if (value & ^(int32(15))) == int32(0) {
			this_.nodeSet_propFlags(((this_.nodeGet_propFlags() & ^(int32(-268435456))) | (javalang.Shl32(value, int(int32(28))))))
			this_.priv_Node_removeColdProp_8422(prop)
			return
		}
		this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(int32(-268435456))))
	} else if prop == StaticGet_Node_Prop_CHANGE_TIME() {
		this_.nodeSet_changeTime(value)
		this_.priv_Node_removeColdProp_8422(prop)
		return
	} else if prop == StaticGet_Node_Prop_SOURCE_FILE() {
		this_.nodeSet_sourceFile(nil)
	} else if prop == StaticGet_Node_Prop_JSDOC_INFO() {
		this_.nodeSet_jsDocInfo(nil)
	} else if prop == StaticGet_Node_Prop_TYPE_BEFORE_CAST() {
		this_.nodeSet_typeBeforeCast(nil)
	}
	var oldHead *Node_PropListItem = this_.nodeGet_propListHead()
	_ = oldHead
	var head *Node_PropListItem = oldHead
	_ = head
	if head == nil {
		if value == int32(0) {
			return
		}
	} else {
		if value != int32(0) {
			var existing *Node_PropListItem = this_.priv_Node_lookupColdProperty_8417(prop)
			_ = existing
			if ((existing != nil) && existing.IsIntItem()) && (existing.intValue == value) {
				return
			}
		}
		head = Node_RebuildListWithoutProp(head, prop)
	}
	if value != int32(0) {
		head = Node_PropListItem_IntItem(int8((prop).Ordinal()), value, head)
	}
	if head != oldHead {
		this_.nodeSet_propListHead(head)
	}
}

// Java method: Node.removeProp(com.google.javascript.rhino.Node.Prop prop)
func (this_ Node) RemoveProp(prop Node_Prop) {
	_ = prop
	if prop == StaticGet_Node_Prop_SOURCE_FILE() {
		this_.nodeSet_sourceFile(nil)
	} else if prop == StaticGet_Node_Prop_JSDOC_INFO() {
		this_.nodeSet_jsDocInfo(nil)
	} else if prop == StaticGet_Node_Prop_TYPE_BEFORE_CAST() {
		this_.nodeSet_typeBeforeCast(nil)
	} else {
		var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(prop).Ordinal()]
		_ = mask
		if mask != int32(0) {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(mask)))
		} else if prop == StaticGet_Node_Prop_CONSTANT_VAR_FLAGS() {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(int32(201326592))))
		} else if prop == StaticGet_Node_Prop_SIDE_EFFECT_FLAGS() {
			this_.nodeSet_propFlags((this_.nodeGet_propFlags() & ^(int32(-268435456))))
		} else if prop == StaticGet_Node_Prop_CHANGE_TIME() {
			this_.nodeSet_changeTime(int32(0))
		}
	}
	this_.priv_Node_removeColdProp_8422(prop)
}

// Java static method: Node.nodePropertyToBit(com.google.javascript.jscomp.serialization.NodeProperty prop)
func Node_NodePropertyToBit(prop *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty) int64 {
	clinit_Node()
	_ = prop
	return javalang.Shl64(int64(1), int(prop.GetNumber()))
}

// Java static method: Node.setNodePropertyBit(long bitset, com.google.javascript.jscomp.serialization.NodeProperty prop)
func Node_SetNodePropertyBit(bitset int64, prop *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty) int64 {
	clinit_Node()
	_ = bitset
	_ = prop
	return (bitset | Node_NodePropertyToBit(prop))
}

// Java static method: Node.removeNodePropertyBit(long bitset, com.google.javascript.jscomp.serialization.NodeProperty prop)
func Node_RemoveNodePropertyBit(bitset int64, prop *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty) int64 {
	clinit_Node()
	_ = bitset
	_ = prop
	return (bitset & ^(Node_NodePropertyToBit(prop)))
}

// Java static method: Node.hasNodePropertyBitSet(long bitset, com.google.javascript.jscomp.serialization.NodeProperty prop)
func Node_HasNodePropertyBitSet(bitset int64, prop *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty) bool {
	clinit_Node()
	_ = bitset
	_ = prop
	return ((bitset & Node_NodePropertyToBit(prop)) != int64(int32(0)))
}

// Java static method: Node.hasBitSet(long bitset, int bit)
func Node_HasBitSet(bitset int64, bit int32) bool {
	clinit_Node()
	_ = bitset
	_ = bit
	return ((bitset & javalang.Shl64(int64(1), int(bit))) != int64(int32(0)))
}

// Java method: Node.serializeProperties()
func (this_ Node) SerializeProperties() int64 {
	var propSet int64 = int64(int32(0))
	_ = propSet
	var flags int32 = this_.nodeGet_propFlags()
	_ = flags
	if flags != int32(0) {
		var boolBits int32 = (flags & int32(67108863))
		_ = boolBits
		for boolBits != int32(0) {
			var bit int32 = j2gort.IntTrailingZeros(boolBits)
			_ = bit
			boolBits = (boolBits & ^(javalang.Shl32(int32(1), int(bit))))
			var prop Node_Prop = StaticGet_Node_INLINE_BIT_TO_PROP()[bit]
			_ = prop
			var nodeProperty *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty = PropTranslator_Serialize(prop)
			_ = nodeProperty
			if nodeProperty != nil {
				propSet = Node_SetNodePropertyBit(propSet, nodeProperty)
			}
		}
		var constFlags int32 = ((javalang.Ushr32(flags, int(int32(26)))) & int32(3))
		_ = constFlags
		if constFlags != int32(0) {
			if Node_AnyBitSet(constFlags, int32(2)) {
				propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_INFERRED_CONSTANT)
			}
			if Node_AnyBitSet(constFlags, int32(1)) {
				propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_DECLARED_CONSTANT)
			}
		}
		var seFlags int32 = ((javalang.Ushr32(flags, int(int32(28)))) & int32(15))
		_ = seFlags
		if seFlags != int32(0) {
			propSet = this_.priv_Node_setNodePropertySideEffectFlags_8433(propSet, seFlags)
		}
	}
	if !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
		propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_COLOR_FROM_CAST)
	}
	{
		var propListItem *Node_PropListItem = this_.nodeGet_propListHead()
		_ = propListItem
		for ; propListItem != nil; propListItem = propListItem.next {
			var prop_2 Node_Prop = StaticGet_Node_PropValues()[int32(propListItem.propType)]
			_ = prop_2
			switch (prop_2).Ordinal() {
			case int32(33) /* Prop.TYPE_BEFORE_CAST */ :
				propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_COLOR_FROM_CAST)
			case int32(25) /* Prop.CONSTANT_VAR_FLAGS */ :
				var intVal int32 = propListItem.GetIntValue()
				_ = intVal
				if Node_AnyBitSet(intVal, int32(2)) {
					propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_INFERRED_CONSTANT)
				}
				if Node_AnyBitSet(intVal, int32(1)) {
					propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_IS_DECLARED_CONSTANT)
				}
			case int32(8) /* Prop.SIDE_EFFECT_FLAGS */ :
				propSet = this_.priv_Node_setNodePropertySideEffectFlags_8433(propSet, propListItem.GetIntValue())
			case int32(54) /* Prop.CLOSURE_UNAWARE_SHADOW */ :
				propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_CLOSURE_UNAWARE_SHADOW)
			default:
				if propListItem.IsIntItem() {
					var nodeProperty_2 *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty = PropTranslator_Serialize(prop_2)
					_ = nodeProperty_2
					if nodeProperty_2 != nil {
						propSet = Node_SetNodePropertyBit(propSet, nodeProperty_2)
					}
				}
			}
		}
	}
	return propSet
}

// Java method: Node.setNodePropertySideEffectFlags(long propSet, int sideEffectFlags)
func (this_ Node) priv_Node_setNodePropertySideEffectFlags_8433(propSet int64, sideEffectFlags int32) int64 {
	_ = propSet
	_ = sideEffectFlags
	if Node_AnyBitSet(sideEffectFlags, int32(1)) {
		propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_MUTATES_GLOBAL_STATE)
	}
	if Node_AnyBitSet(sideEffectFlags, int32(2)) {
		propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_MUTATES_THIS)
	}
	if Node_AnyBitSet(sideEffectFlags, int32(4)) {
		propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_MUTATES_ARGUMENTS)
	}
	if Node_AnyBitSet(sideEffectFlags, int32(8)) {
		propSet = Node_SetNodePropertyBit(propSet, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_THROWS)
	}
	return propSet
}

// Java method: Node.getPropListDebugString()
func (this_ Node) GetPropListDebugString() string {
	var propListItem *Node_PropListItem = this_.nodeGet_propListHead()
	_ = propListItem
	var propListStrRepr []string = []string{}
	_ = propListStrRepr
	for propListItem != nil {
		propListStrRepr = append(propListStrRepr, (StaticGet_Node_PropValues()[int32(propListItem.propType)]).Name())
		propListItem = propListItem.next
	}
	if this_.nodeGet_jsDocInfo() != nil {
		propListStrRepr = append(propListStrRepr, (StaticGet_Node_Prop_JSDOC_INFO()).Name())
	}
	if !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
		propListStrRepr = append(propListStrRepr, (StaticGet_Node_Prop_TYPE_BEFORE_CAST()).Name())
	}
	if this_.nodeGet_changeTime() != int32(0) {
		propListStrRepr = append(propListStrRepr, (StaticGet_Node_Prop_CHANGE_TIME()).Name())
	}
	var flags int32 = this_.nodeGet_propFlags()
	_ = flags
	if flags != int32(0) {
		if ((javalang.Ushr32(flags, int(int32(28)))) & int32(15)) != int32(0) {
			propListStrRepr = append(propListStrRepr, (StaticGet_Node_Prop_SIDE_EFFECT_FLAGS()).Name())
		}
		if ((javalang.Ushr32(flags, int(int32(26)))) & int32(3)) != int32(0) {
			propListStrRepr = append(propListStrRepr, (StaticGet_Node_Prop_CONSTANT_VAR_FLAGS()).Name())
		}
		var boolBits int32 = (flags & (int32(67108863) & ^(int32(16777216))))
		_ = boolBits
		if (this_.nodeGet_propListHead() != nil) && (boolBits != int32(0)) {
			{
				var x *Node_PropListItem = this_.nodeGet_propListHead()
				_ = x
				for ; x != nil; x = x.next {
					boolBits = (boolBits & ^(StaticGet_Node_INLINE_BOOL_MASKS()[int32(x.propType)]))
				}
			}
		}
		for boolBits != int32(0) {
			var bit int32 = j2gort.IntTrailingZeros(boolBits)
			_ = bit
			boolBits = (boolBits & ^(javalang.Shl32(int32(1), int(bit))))
			propListStrRepr = append(propListStrRepr, (StaticGet_Node_INLINE_BIT_TO_PROP()[bit]).Name())
		}
		if ((flags & int32(16777216)) != int32(0)) && ((this_.nodeGet_propListHead() == nil) || (this_.priv_Node_lookupColdProperty_8417(StaticGet_Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE()) == nil)) {
			propListStrRepr = append(propListStrRepr, (StaticGet_Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE()).Name())
		}
	}
	if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
		propListStrRepr = append(propListStrRepr, (StaticGet_Node_Prop_SOURCE_FILE()).Name())
	}
	var joiner *j2gort.Ext_com_google_common_base_Joiner = j2gort.Ext_com_google_common_base_Joiner_On_1_string("->")
	_ = joiner
	return joiner.Join_1_Ext_java_lang_Iterable(j2gort.SliceReadOnlyIterable(propListStrRepr))
}

// Java method: Node.deserializeProperties(long propSet, boolean isInClosureUnawareSubtree)
func (this_ Node) DeserializeProperties(propSet int64, isInClosureUnawareSubtree bool) {
	_ = propSet
	_ = isInClosureUnawareSubtree
	if this_.IsRoot() {
		if !(((((j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil) && (this_.nodeGet_jsDocInfo() == nil)) && j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) && (this_.nodeGet_propFlags() == int32(0))) && (this_.nodeGet_changeTime() == int32(0))) && (this_.nodeGet_propListHead() == nil)) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(this_.nodeGet_propListHead()))
		}
	} else {
		if !(this_.priv_Node_validatePropListTailOrdering_8524(isInClosureUnawareSubtree)) {
			j2gort.CheckStateStr(false, "Node not ready to receive properties")
		}
	}
	var constantVarFlags int32 = int32(0)
	_ = constantVarFlags
	var sideEffectFlags int32 = int32(0)
	_ = sideEffectFlags
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(63); i = i + 1 {
			if !(Node_HasBitSet(propSet, i)) {
				continue
			}
			var nodeProperty *j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty = j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_ForNumber(i)
			_ = nodeProperty
			switch j2gort.CheckNotNullPtr(nodeProperty).Ordinal() {
			case int32(36) /* NodeProperty.IS_DECLARED_CONSTANT */ :
				constantVarFlags = (constantVarFlags | int32(1))
			case int32(35) /* NodeProperty.IS_INFERRED_CONSTANT */ :
				constantVarFlags = (constantVarFlags | int32(2))
			case int32(38) /* NodeProperty.MUTATES_GLOBAL_STATE */ :
				sideEffectFlags = (sideEffectFlags | int32(1))
			case int32(39) /* NodeProperty.MUTATES_THIS */ :
				sideEffectFlags = (sideEffectFlags | int32(2))
			case int32(40) /* NodeProperty.MUTATES_ARGUMENTS */ :
				sideEffectFlags = (sideEffectFlags | int32(4))
			case int32(41) /* NodeProperty.THROWS */ :
				sideEffectFlags = (sideEffectFlags | int32(8))
			case int32(42) /* NodeProperty.CLOSURE_UNAWARE_SHADOW */ :
			default:
				var prop Node_Prop = PropTranslator_Deserialize(nodeProperty)
				_ = prop
				if prop == 0 {
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Can not translate "+j2gort.StrVal(nodeProperty)+" to AST Prop")))
				}
				var mask int32 = StaticGet_Node_INLINE_BOOL_MASKS()[(prop).Ordinal()]
				_ = mask
				if mask != int32(0) {
					this_.nodeSet_propFlags((this_.nodeGet_propFlags() | mask))
				} else {
					this_.nodeSet_propListHead(Node_PropListItem_IntItem(int8((prop).Ordinal()), int32(1), this_.nodeGet_propListHead()))
				}
			}
		}
	}
	if constantVarFlags != int32(0) {
		this_.nodeSet_propFlags(((this_.nodeGet_propFlags() & ^(int32(201326592))) | (javalang.Shl32((constantVarFlags & int32(3)), int(int32(26))))))
	}
	if sideEffectFlags != int32(0) {
		this_.nodeSet_propFlags(((this_.nodeGet_propFlags() & ^(int32(-268435456))) | (javalang.Shl32((sideEffectFlags & int32(15)), int(int32(28))))))
	}
	if this_.HasPropertyViolations() {
		this_.ValidateProperties(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_57 any) {
			errorMessage := j2gort.CastString(lp_57)
			_ = errorMessage
			if !(errorMessage == j2gort.NullString) {
				j2gort.CheckStateFmt(false, "deserialize error: %s: %s", j2gort.BoxString(errorMessage), j2gort.BoxNode(this_))
			}
		}))
	}
}

// Java method: Node.setDeclaredTypeExpression(com.google.javascript.rhino.Node typeExpression)
func (this_ Node) SetDeclaredTypeExpression(typeExpression Node) {
	_ = typeExpression
	this_.PutProp(StaticGet_Node_Prop_DECLARED_TYPE_EXPR(), j2gort.BoxNode(typeExpression))
}

// Java method: Node.getReflectedObject()
func (this_ Node) GetReflectedObject() any {
	return this_.GetProp(StaticGet_Node_Prop_REFLECTED_OBJECT())
}

// Java method: Node.setReflectedObject(java.lang.Object object)
func (this_ Node) SetReflectedObject(object any) {
	_ = object
	this_.PutProp(StaticGet_Node_Prop_REFLECTED_OBJECT(), object)
}

// Java method: Node.getGenericType()
func (this_ Node) GetGenericType() Node {
	return j2gort.CastNode[Node](this_.GetProp(StaticGet_Node_Prop_GENERIC_TYPE()))
}

// Java method: Node.setGenericType(com.google.javascript.rhino.Node node)
func (this_ Node) SetGenericType(node Node) {
	_ = node
	this_.PutProp(StaticGet_Node_Prop_GENERIC_TYPE(), j2gort.BoxNode(node))
}

// Java method: Node.getImplements()
func (this_ Node) GetImplements() Node {
	return j2gort.CastNode[Node](this_.GetProp(StaticGet_Node_Prop_IMPLEMENTS()))
}

// Java method: Node.setImplements(com.google.javascript.rhino.Node node)
func (this_ Node) SetImplements(node Node) {
	_ = node
	this_.PutProp(StaticGet_Node_Prop_IMPLEMENTS(), j2gort.BoxNode(node))
}

// Java method: Node.getConstructSignature()
func (this_ Node) GetConstructSignature() bool {
	return this_.GetBooleanProp(StaticGet_Node_Prop_CONSTRUCT_SIGNATURE())
}

// Java method: Node.setConstructSignature(boolean value)
func (this_ Node) SetConstructSignature(value bool) {
	_ = value
	this_.PutBooleanProp(StaticGet_Node_Prop_CONSTRUCT_SIGNATURE(), value)
}

// Java method: Node.getAccessModifier()
func (this_ Node) GetAccessModifier() JSDocInfo_Visibility {
	return j2gort.Cast[JSDocInfo_Visibility](this_.GetProp(StaticGet_Node_Prop_ACCESS_MODIFIER()))
}

// Java method: Node.setAccessModifier(com.google.javascript.rhino.JSDocInfo.Visibility visibility)
func (this_ Node) SetAccessModifier(visibility JSDocInfo_Visibility) {
	_ = visibility
	this_.PutProp(StaticGet_Node_Prop_ACCESS_MODIFIER(), j2gort.BoxNullableVal(visibility))
}

// Java method: Node.getFeatureSet()
func (this_ Node) GetFeatureSet() any {
	return this_.GetProp(StaticGet_Node_Prop_FEATURE_SET())
}

// Java method: Node.setFeatureSet(java.lang.Object featureSet)
func (this_ Node) SetFeatureSet(featureSet any) {
	_ = featureSet
	this_.PutProp(StaticGet_Node_Prop_FEATURE_SET(), featureSet)
}

// Java method: Node.removeDirectEvalProp()
func (this_ Node) RemoveDirectEvalProp() {
	this_.RemoveProp(StaticGet_Node_Prop_DIRECT_EVAL())
}

// Java method: Node.getDeclaredTypeExpression()
func (this_ Node) GetDeclaredTypeExpression() Node {
	return j2gort.CastNode[Node](this_.GetProp(StaticGet_Node_Prop_DECLARED_TYPE_EXPR()))
}

// Java method: Node.setJSTypeBeforeCast(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func (this_ Node) SetJSTypeBeforeCast(type_ JSType) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if type_ == nil || type_.IsJavaNil_() {
		nodestore.Node(this_).ClearJSTypeBeforeCast()
		return
	}

	_ = type_
	nodestore.Node(this_).SetJSTypeBeforeCastObj(type_)

}

// Java method: Node.getJSTypeBeforeCast()
func (this_ Node) GetJSTypeBeforeCast() JSType {
	return this_.nodeGet_typeBeforeCast()
}

// Java method: Node.setColorFromTypeCast()
func (this_ Node) SetColorFromTypeCast() {
	if !(this_.GetColor() != nil) {
		j2gort.CheckStateStr(false, "Only use on nodes with colors present")
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_COLOR_FROM_CAST(), true)
}

// Java method: Node.isColorFromTypeCast()
func (this_ Node) IsColorFromTypeCast() bool {
	return ((this_.nodeGet_propFlags() & int32(65536)) != int32(0))
}

// Java method: Node.getSortedPropTypes()
func (this_ Node) priv_Node_getSortedPropTypes_8455() []int8 {
	var flags int32 = this_.nodeGet_propFlags()
	_ = flags
	var head *Node_PropListItem = this_.nodeGet_propListHead()
	_ = head
	if ((((j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil) && (this_.nodeGet_jsDocInfo() == nil)) && j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) && (flags == int32(0))) && (this_.nodeGet_changeTime() == int32(0))) && (head == nil) {
		return make([]int8, int32(0))
	}
	var inlineBoolBits int32 = (flags & int32(67108863))
	_ = inlineBoolBits
	if (head != nil) && (inlineBoolBits != int32(0)) {
		{
			var x *Node_PropListItem = head
			_ = x
			for ; x != nil; x = x.next {
				inlineBoolBits = (inlineBoolBits & ^(StaticGet_Node_INLINE_BOOL_MASKS()[int32(x.propType)]))
			}
		}
	}
	var count int32 = j2gort.IntBitCount(inlineBoolBits)
	_ = count
	if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
		count = count + 1
	}
	if this_.nodeGet_jsDocInfo() != nil {
		count = count + 1
	}
	if !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
		count = count + 1
	}
	if this_.nodeGet_changeTime() != int32(0) {
		count = count + 1
	}
	if (flags & int32(201326592)) != int32(0) {
		count = count + 1
	}
	if (flags & int32(-268435456)) != int32(0) {
		count = count + 1
	}
	{
		var x_2 *Node_PropListItem = head
		_ = x_2
		for ; x_2 != nil; x_2 = x_2.next {
			count = count + 1
		}
	}
	var keys []int8 = make([]int8, count)
	_ = keys
	var idx int32 = int32(0)
	_ = idx
	if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
		{
			a_58 := keys
			i_59 := func() int32 { old_60 := idx; nw_61 := old_60 + 1; idx = nw_61; return old_60 }()
			a_58[i_59] = int8((StaticGet_Node_Prop_SOURCE_FILE()).Ordinal())
		}
	}
	if this_.nodeGet_jsDocInfo() != nil {
		{
			a_62 := keys
			i_63 := func() int32 { old_64 := idx; nw_65 := old_64 + 1; idx = nw_65; return old_64 }()
			a_62[i_63] = int8((StaticGet_Node_Prop_JSDOC_INFO()).Ordinal())
		}
	}
	if !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
		{
			a_66 := keys
			i_67 := func() int32 { old_68 := idx; nw_69 := old_68 + 1; idx = nw_69; return old_68 }()
			a_66[i_67] = int8((StaticGet_Node_Prop_TYPE_BEFORE_CAST()).Ordinal())
		}
	}
	if this_.nodeGet_changeTime() != int32(0) {
		{
			a_70 := keys
			i_71 := func() int32 { old_72 := idx; nw_73 := old_72 + 1; idx = nw_73; return old_72 }()
			a_70[i_71] = int8((StaticGet_Node_Prop_CHANGE_TIME()).Ordinal())
		}
	}
	for inlineBoolBits != int32(0) {
		var bit int32 = j2gort.IntTrailingZeros(inlineBoolBits)
		_ = bit
		inlineBoolBits = (inlineBoolBits & ^(javalang.Shl32(int32(1), int(bit))))
		{
			a_74 := keys
			i_75 := func() int32 { old_76 := idx; nw_77 := old_76 + 1; idx = nw_77; return old_76 }()
			a_74[i_75] = int8((StaticGet_Node_INLINE_BIT_TO_PROP()[bit]).Ordinal())
		}
	}
	if (flags & int32(201326592)) != int32(0) {
		{
			a_78 := keys
			i_79 := func() int32 { old_80 := idx; nw_81 := old_80 + 1; idx = nw_81; return old_80 }()
			a_78[i_79] = int8((StaticGet_Node_Prop_CONSTANT_VAR_FLAGS()).Ordinal())
		}
	}
	if (flags & int32(-268435456)) != int32(0) {
		{
			a_82 := keys
			i_83 := func() int32 { old_84 := idx; nw_85 := old_84 + 1; idx = nw_85; return old_84 }()
			a_82[i_83] = int8((StaticGet_Node_Prop_SIDE_EFFECT_FLAGS()).Ordinal())
		}
	}
	{
		var x_3 *Node_PropListItem = head
		_ = x_3
		for ; x_3 != nil; x_3 = x_3.next {
			{
				a_86 := keys
				i_87 := func() int32 { old_88 := idx; nw_89 := old_88 + 1; idx = nw_89; return old_88 }()
				a_86[i_87] = x_3.propType
			}
		}
	}
	j2gort.Ext_java_util_Arrays_Sort_1_slice_int8(keys)
	return keys
}

// Java method: Node.getDouble()
func (this_ Node) GetDouble() float64 {
	return (this_).nodeGet_number()
}

// Java method: Node.setDouble(double x)
func (this_ Node) SetDouble(x float64) {
	_ = x
	if !(!(math.IsNaN(x))) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfDouble(x))
	}
	if !(pkg_base.JSCompDoubles_IsPositive(x)) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfDouble(x))
	}
	{
		r_90 := (this_)
		r_90.nodeSet_number(x)
	}
}

// Java method: Node.getBigInt()
func (this_ Node) GetBigInt() *j2gort.Ext_java_math_BigInteger {
	return (this_).nodeGet_bigint()
}

// Java method: Node.setBigInt(java.math.BigInteger number)
func (this_ Node) SetBigInt(number *j2gort.Ext_java_math_BigInteger) {
	_ = number
	j2gort.CheckNotNullPtr(number)
	if !(number.Signum() >= int32(0)) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(number))
	}
	{
		r_91 := (this_)
		r_91.nodeSet_bigint(number)
	}
}

// Java method: Node.getString()
func (this_ Node) GetString() string {
	return (this_).nodeGet_str()
}

// Java method: Node.setString(java.lang.String str)
func (this_ Node) SetString(str string) {
	_ = str
	{
		r_92 := (this_)
		r_92.nodeSet_str(RhinoStringPool_AddOrGet(str))
	}
}

// Java method: Node.setStringFromStringPool(LazyInternedStringList, int)
func (this_ Node) SetStringFromStringPool(stringPool *RhinoStringPool_LazyInternedStringList, offset int32) {
	nodestore.Node(this_).SetStringPayload(stringPool.InternID_(offset))
}

// Java method: Node.getRawString()
func (this_ Node) GetRawString() string {
	return (this_).nodeGet_raw()
}

// Java method: Node.getCookedString()
func (this_ Node) GetCookedString() string {
	return (this_).nodeGet_cooked()
}

// Java method: Node.toString()
func (this_ Node) String() string {
	return this_.ToString_3(true, true, true)
}

// Java method: Node.toString(boolean printSource, boolean printAnnotations, boolean printType)
func (this_ Node) ToString_3(printSource bool, printAnnotations bool, printType bool) string {
	_ = printSource
	_ = printAnnotations
	_ = printType
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	this_.priv_Node_toString_8467(sb, printSource, printAnnotations, printType)
	return sb.String()
}

// Java method: Node.toString(java.lang.StringBuilder sb, boolean printSource, boolean printAnnotations, boolean printType)
func (this_ Node) priv_Node_toString_8467(sb *j2gort.Ext_java_lang_StringBuilder, printSource bool, printAnnotations bool, printType bool) {
	_ = sb
	_ = printSource
	_ = printAnnotations
	_ = printType
	sb.Append_1_any_9850(j2gort.BoxNullableVal(this_.nodeGet_token()))
	if this_ != 0 && nodestore.Node(this_).KindFlags() == 3 {
		sb.Append_1_uint16(uint16(32))
		sb.Append_1_string(this_.GetString())
	} else if this_.RawToken_() == uint8(73) {
		sb.Append_1_uint16(uint16(32))
		if (this_.nodeGet_first() == 0) || (this_.nodeGet_first().RawToken_() != uint8(32)) {
			sb.Append_1_string("<invalid>")
		} else {
			sb.Append_1_string(this_.nodeGet_first().GetString())
		}
	} else if this_.RawToken_() == uint8(33) {
		sb.Append_1_uint16(uint16(32))
		sb.Append_1_float64(this_.GetDouble())
	}
	if printSource {
		var lineno int32 = this_.GetLineno()
		_ = lineno
		if lineno != int32(-1) {
			sb.Append_1_uint16(uint16(32))
			sb.Append_1_int32(lineno)
			sb.Append_1_uint16(uint16(58))
			sb.Append_1_int32(this_.GetCharno())
			sb.Append_1_uint16(uint16(32))
		}
		if this_.nodeGet_length() != int32(0) {
			sb.Append_1_string(" [length: ")
			sb.Append_1_int32(this_.nodeGet_length())
			sb.Append_1_uint16(uint16(93))
		}
	}
	if printAnnotations {
		var keys []int8 = this_.priv_Node_getSortedPropTypes_8455()
		_ = keys
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(keys)); i = i + 1 {
				var type_ Node_Prop = StaticGet_Node_PropValues()[int32(keys[i])]
				_ = type_
				if type_ == StaticGet_Node_Prop_CLOSURE_UNAWARE_SHADOW() {
					sb.Append_1_string(" [is_shadow_host]")
					continue
				}
				sb.Append_1_string(" [")
				sb.Append_1_string(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(j2gort.StrVal(type_)))
				sb.Append_1_string(": ")
				if (type_ == StaticGet_Node_Prop_SOURCE_FILE()) && !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
					sb.Append_1_any_9850(this_.nodeGet_sourceFile())
				} else if (type_ == StaticGet_Node_Prop_JSDOC_INFO()) && (this_.nodeGet_jsDocInfo() != nil) {
					sb.Append_1_any_9850(j2gort.PtrToAny(this_.nodeGet_jsDocInfo()))
				} else if (type_ == StaticGet_Node_Prop_TYPE_BEFORE_CAST()) && !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
					sb.Append_1_any_9850(this_.nodeGet_typeBeforeCast())
				} else if (type_ == StaticGet_Node_Prop_CHANGE_TIME()) && (this_.nodeGet_changeTime() != int32(0)) {
					sb.Append_1_int32(this_.nodeGet_changeTime())
				} else {
					sb.Append_1_any_9850(j2gort.PtrToAny(this_.LookupProperty(type_)))
				}
				sb.Append_1_uint16(uint16(93))
			}
		}
		if this_.nodeGet_originalName() != j2gort.NullString {
			sb.Append_1_string(" [original_name: ")
			sb.Append_1_string(this_.nodeGet_originalName())
			sb.Append_1_uint16(uint16(93))
		}
	}
	if printType && (this_.nodeGet_jstypeOrColor() != nil) {
		var typeString string = javalang.ToJavaString(j2gort.CheckNotNull(this_.nodeGet_jstypeOrColor()))
		_ = typeString
		if typeString != j2gort.NullString {
			sb.Append_1_string(" : ")
			sb.Append_1_string(typeString)
		}
	}
}

// Java static method: Node.createJsonPair(java.lang.String name, java.lang.String value)
func Node_CreateJsonPair(name string, value string) string {
	clinit_Node()
	_ = name
	_ = value
	return ("\"" + j2gort.StrValString(name) + "\":\"" + j2gort.StrValString(value) + "\"")
}

// Java static method: Node.createJsonPairRawValue(java.lang.String name, java.lang.String value)
func Node_CreateJsonPairRawValue(name string, value string) string {
	clinit_Node()
	_ = name
	_ = value
	return ("\"" + j2gort.StrValString(name) + "\":" + j2gort.StrValString(value))
}

// Java method: Node.toJson(java.lang.Appendable sb)
func (this_ Node) priv_Node_toJson_8468(sb j2gort.Ext_java_lang_Appendable) {
	_ = sb
	sb.Append_1_uint16(uint16(123))
	sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPair("token", j2gort.CheckNotNullEnum(this_.nodeGet_token()).String())))
	if this_ != 0 && nodestore.Node(this_).KindFlags() == 3 {
		sb.Append_1_uint16(uint16(44))
		sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPair("string", this_.GetString())))
	} else if this_.RawToken_() == uint8(73) {
		sb.Append_1_uint16(uint16(44))
		if (this_.nodeGet_first() == 0) || (this_.nodeGet_first().RawToken_() != uint8(32)) {
			sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPair("functionName", "<invalid>")))
		} else {
			sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPair("functionName", this_.nodeGet_first().GetString())))
		}
	} else if this_.RawToken_() == uint8(33) {
		sb.Append_1_uint16(uint16(44))
		sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPair("number", javalang.DoubleToString(this_.GetDouble()))))
	}
	var lineno int32 = this_.GetLineno()
	_ = lineno
	if lineno != int32(-1) {
		sb.Append_1_uint16(uint16(44))
		sb.Append_1_any("\"sourceLocation\":{")
		sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPairRawValue("line", strconv.FormatInt(int64(lineno), 10))))
		sb.Append_1_uint16(uint16(44))
		sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPairRawValue("col", strconv.FormatInt(int64(this_.GetCharno()), 10))))
		if this_.nodeGet_length() != int32(0) {
			sb.Append_1_uint16(uint16(44))
			sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPairRawValue("length", strconv.FormatInt(int64(this_.nodeGet_length()), 10))))
		}
		sb.Append_1_any("}")
	}
	if this_.nodeGet_originalName() != j2gort.NullString {
		sb.Append_1_any(",")
		sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPair("original_name", this_.nodeGet_originalName())))
	}
	var keys []int8 = this_.priv_Node_getSortedPropTypes_8455()
	_ = keys
	if int32(len(keys)) != int32(0) {
		sb.Append_1_any(",")
		sb.Append_1_any("\"props\":{")
		{
			var i int32 = int32(0)
			_ = i
			for ; i < int32(len(keys)); i = i + 1 {
				var type_ Node_Prop = StaticGet_Node_PropValues()[int32(keys[i])]
				_ = type_
				var valStr string
				_ = valStr
				if (type_ == StaticGet_Node_Prop_SOURCE_FILE()) && !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
					valStr = j2gort.StrVal(this_.nodeGet_sourceFile())
				} else if (type_ == StaticGet_Node_Prop_JSDOC_INFO()) && (this_.nodeGet_jsDocInfo() != nil) {
					valStr = j2gort.StrVal(this_.nodeGet_jsDocInfo())
				} else if (type_ == StaticGet_Node_Prop_TYPE_BEFORE_CAST()) && !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
					valStr = j2gort.StrVal(this_.nodeGet_typeBeforeCast())
				} else if (type_ == StaticGet_Node_Prop_CHANGE_TIME()) && (this_.nodeGet_changeTime() != int32(0)) {
					valStr = strconv.FormatInt(int64(this_.nodeGet_changeTime()), 10)
				} else {
					valStr = this_.LookupProperty(type_).String()
				}
				sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPair(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(j2gort.StrVal(type_)), javalang.StringReplace(javalang.StringReplace(valStr, "\n", "\\n"), "\"", "\\\""))))
				if (i + int32(1)) < int32(len(keys)) {
					sb.Append_1_uint16(uint16(44))
				}
			}
		}
		sb.Append_1_any("}")
	}
	if this_.nodeGet_jstypeOrColor() != nil {
		var typeString string = javalang.ToJavaString(j2gort.CheckNotNull(this_.nodeGet_jstypeOrColor()))
		_ = typeString
		if typeString != j2gort.NullString {
			sb.Append_1_uint16(uint16(44))
			sb.Append_1_any(j2gort.BoxString(Node_CreateJsonPair("typeString", javalang.StringReplace(typeString, "\"", "\\\""))))
		}
	}
	if this_.nodeGet_first() != 0 {
		sb.Append_1_uint16(uint16(44))
		sb.Append_1_any("\"children\":[")
		{
			var child Node = this_.nodeGet_first()
			_ = child
			for ; child != 0; child = child.nodeGet_next() {
				child.priv_Node_toJson_8468(sb)
				if child.nodeGet_next() != 0 {
					sb.Append_1_uint16(uint16(44))
				}
			}
		}
		sb.Append_1_any("]")
	}
	sb.Append_1_uint16(uint16(125))
}

// Java method: Node.toStringTree()
func (this_ Node) ToStringTree() string {
	return this_.priv_Node_toStringTreeImpl_8470()
}

// Java method: Node.toStringTreeImpl()
func (this_ Node) priv_Node_toStringTreeImpl_8470() string {
	var didRet_93 bool
	var retVal_94 string
	func() {
		defer func() {
			if rec_95 := recover(); rec_95 != nil {
				jex_96 := javalang.RecoverException(rec_95)
				_ = jex_96
				if j2gort.InstanceOf[javaio.IIOException](jex_96) {
					e := j2gort.AsIface[javaio.IIOException](jex_96)
					_ = e
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Should not happen\n"+j2gort.StrVal(e))))
				} else {
					panic(rec_95)
				}
			}
		}()
		var s *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = s
		this_.AppendStringTree(s)
		retVal_94 = s.String()
		didRet_93 = true
		return
	}()
	if didRet_93 {
		return retVal_94
	}
	panic("unreachable")
}

// Java method: Node.appendStringTree(java.lang.Appendable appendable)
func (this_ Node) AppendStringTree(appendable j2gort.Ext_java_lang_Appendable) {
	_ = appendable
	Node_ToStringTreeHelper(this_, int32(0), appendable)
}

// Java static method: Node.toStringTreeHelper(com.google.javascript.rhino.Node n, int level, java.lang.Appendable sb)
func Node_ToStringTreeHelper(n Node, level int32, sb j2gort.Ext_java_lang_Appendable) {
	clinit_Node()
	_ = n
	_ = level
	_ = sb
	{
		var i int32 = int32(0)
		_ = i
		for ; i != level; i = i + 1 {
			sb.Append_1_any("    ")
		}
	}
	sb.Append_1_any(j2gort.BoxString(n.String()))
	sb.Append_1_uint16(uint16(10))
	{
		var cursor Node = n.nodeGet_first()
		_ = cursor
		for ; cursor != 0; cursor = cursor.nodeGet_next() {
			Node_ToStringTreeHelper(cursor, (level + int32(1)), sb)
		}
	}
	var shadow Node = n.GetClosureUnawareShadow()
	_ = shadow
	if shadow != 0 {
		Node_ToStringTreeHelper(shadow, (level + int32(1)), sb)
	}
}

// Java method: Node.appendJsonTree(java.lang.Appendable appendable)
func (this_ Node) AppendJsonTree(appendable j2gort.Ext_java_lang_Appendable) {
	_ = appendable
	Node_ToJsonTreeHelper(this_, appendable)
}

// Java static method: Node.toJsonTreeHelper(com.google.javascript.rhino.Node n, java.lang.Appendable sb)
func Node_ToJsonTreeHelper(n Node, sb j2gort.Ext_java_lang_Appendable) {
	clinit_Node()
	_ = n
	_ = sb
	n.priv_Node_toJson_8468(sb)
}

// Java method: Node.setStaticSourceFileFrom(com.google.javascript.rhino.Node other)
func (this_ Node) SetStaticSourceFileFrom(other Node) {
	_ = other
	this_.nodeSet_sourceFile(other.GetStaticSourceFile())
	this_.priv_Node_removeColdProp_8422(StaticGet_Node_Prop_SOURCE_FILE())
}

// Java method: Node.setStaticSourceFile(com.google.javascript.rhino.@org.jspecify.annotations.Nullable StaticSourceFile file)
func (this_ Node) SetStaticSourceFile(file StaticSourceFile) Node {
	_ = file
	this_.nodeSet_sourceFile(file)
	this_.priv_Node_removeColdProp_8422(StaticGet_Node_Prop_SOURCE_FILE())
	return this_
}

// Java method: Node.setSourceFileForTesting(java.lang.String name)
func (this_ Node) SetSourceFileForTesting(name string) {
	_ = name
	this_.SetStaticSourceFile(NewSimpleSourceFile(name, StaticGet_StaticSourceFile_SourceKind_STRONG()))
}

// Java method: Node.getSourceFileName()
func (this_ Node) GetSourceFileName() string {
	var file StaticSourceFile = this_.GetStaticSourceFile()
	_ = file
	if file == nil || file.IsJavaNil_() {
		return j2gort.NullString
	} else {
		return file.GetName()
	}
}

// Java method: Node.getStaticSourceFile()
func (this_ Node) GetStaticSourceFile() StaticSourceFile {
	if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
		return this_.nodeGet_sourceFile()
	}
	if this_.nodeGet_propListHead() != nil {
		return j2gort.AsIface[StaticSourceFile](this_.GetProp(StaticGet_Node_Prop_SOURCE_FILE()))
	}
	return nil
}

// Java method: Node.setInputId(com.google.javascript.rhino.InputId inputId)
func (this_ Node) SetInputId(inputId *InputId) {
	_ = inputId
	this_.PutProp(StaticGet_Node_Prop_INPUT_ID(), j2gort.PtrToAny(inputId))
}

// Java method: Node.getInputId()
func (this_ Node) GetInputId() *InputId {
	return (j2gort.CastPtr[InputId](this_.GetProp(StaticGet_Node_Prop_INPUT_ID())))
}

// Java method: Node.getOriginalName()
func (this_ Node) GetOriginalName() string {
	return this_.nodeGet_originalName()
}

// Java method: Node.setOriginalName(java.lang.String s)
func (this_ Node) SetOriginalName(s string) {
	_ = s
	if s == j2gort.NullString {
		this_.nodeSet_originalName(j2gort.NullString)
	} else {
		this_.nodeSet_originalName(RhinoStringPool_AddOrGet(s))
	}
}

// Java method: Node.setOriginalNameFromStringPool(LazyInternedStringList, int)
func (this_ Node) SetOriginalNameFromStringPool(stringPool *RhinoStringPool_LazyInternedStringList, offset int32) {
	id := stringPool.InternID_(offset)
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	nodestore.Node(this_).SetOrigNameID(id)
}

// Java method: Node.setOriginalNameFromName(com.google.javascript.rhino.Node name)
func (this_ Node) SetOriginalNameFromName(name Node) {
	_ = name
	this_.nodeSet_originalName(name.GetString())
}

// Java method: Node.isIndexable()
func (this_ Node) IsIndexable() bool {
	return ((this_.nodeGet_propFlags() & int32(8)) == int32(0))
}

// Java method: Node.makeNonIndexable()
func (this_ Node) MakeNonIndexable() {
	this_.PutBooleanProp(StaticGet_Node_Prop_NON_INDEXABLE(), true)
}

// Java method: Node.makeNonIndexableRecursive()
func (this_ Node) MakeNonIndexableRecursive() {
	this_.MakeNonIndexable()
	{
		var child Node = this_.nodeGet_first()
		_ = child
		for ; child != 0; child = child.GetNext() {
			child.MakeNonIndexableRecursive()
		}
	}
}

// Java method: Node.isFromExterns()
func (this_ Node) IsFromExterns() bool {
	var file StaticSourceFile = this_.GetStaticSourceFile()
	_ = file
	if file == nil || file.IsJavaNil_() {
		return false
	} else {
		return file.IsExtern()
	}
}

// Java method: Node.isClosureUnawareCode()
func (this_ Node) IsClosureUnawareCode() bool {
	var file StaticSourceFile = this_.GetStaticSourceFile()
	_ = file
	return ((file != nil && !file.IsJavaNil_()) && file.IsClosureUnawareCode())
}

// Java method: Node.getLength()
func (this_ Node) GetLength() int32 {
	return this_.nodeGet_length()
}

// Java method: Node.setLength(int length)
func (this_ Node) SetLength(length int32) {
	_ = length
	this_.nodeSet_length(length)
}

// Java method: Node.getLineno()
func (this_ Node) GetLineno() int32 {
	if this_.nodeGet_linenoCharno() == int32(-1) {
		return int32(-1)
	} else {
		return javalang.Ushr32(this_.nodeGet_linenoCharno(), int(int32(12)))
	}
}

// Java method: Node.getCharno()
func (this_ Node) GetCharno() int32 {
	if this_.nodeGet_linenoCharno() == int32(-1) {
		return int32(-1)
	} else {
		return (this_.nodeGet_linenoCharno() & int32(4095))
	}
}

// Java method: Node.getLocation()
func (this_ Node) GetLocation() string {
	return (j2gort.StrValString(this_.GetSourceFileName()) + ":" + j2gort.StrValInt32(this_.GetLineno()) + ":" + j2gort.StrValInt32(this_.GetCharno()))
}

// Java method: Node.getSourceOffset()
func (this_ Node) GetSourceOffset() int32 {
	var file StaticSourceFile = this_.GetStaticSourceFile()
	_ = file
	if file == nil || file.IsJavaNil_() {
		return int32(-1)
	}
	var lineno int32 = this_.GetLineno()
	_ = lineno
	if lineno == int32(-1) {
		return int32(-1)
	}
	return (file.GetLineOffset(lineno) + this_.GetCharno())
}

// Java method: Node.getSourcePosition()
func (this_ Node) GetSourcePosition() int32 {
	return this_.nodeGet_linenoCharno()
}

// Java method: Node.setLinenoCharno(int lineno, int charno)
func (this_ Node) SetLinenoCharno(lineno int32, charno int32) Node {
	_ = lineno
	_ = charno
	if (lineno < int32(0)) || (charno < int32(0)) {
		this_.nodeSet_linenoCharno(int32(-1))
		return this_
	}
	if charno > int32(4095) {
		charno = int32(4095)
	}
	this_.nodeSet_linenoCharno(((javalang.Shl32(lineno, int(int32(12)))) | charno))
	return this_
}

// Java method: Node.children()
func (this_ Node) Children() j2gort.Ext_java_lang_Iterable {
	if this_.nodeGet_first() == 0 {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	} else {
		return NewNode_SiblingNodeIterable(this_.nodeGet_first())
	}
}

// Java method: Node.getPropListHeadForTesting()
func (this_ Node) GetPropListHeadForTesting() *Node_PropListItem {
	if this_.nodeGet_propListHead() != nil {
		return this_.nodeGet_propListHead()
	}
	if this_.nodeGet_jsDocInfo() != nil {
		return Node_PropListItem_ObjectItem(int8((StaticGet_Node_Prop_JSDOC_INFO()).Ordinal()), j2gort.PtrToAny(this_.nodeGet_jsDocInfo()), nil)
	}
	if !(j2gort.IfaceEquals[JSType](this_.nodeGet_typeBeforeCast(), nil)) {
		return Node_PropListItem_ObjectItem(int8((StaticGet_Node_Prop_TYPE_BEFORE_CAST()).Ordinal()), this_.nodeGet_typeBeforeCast(), nil)
	}
	if this_.nodeGet_changeTime() != int32(0) {
		return Node_PropListItem_IntItem(int8((StaticGet_Node_Prop_CHANGE_TIME()).Ordinal()), this_.nodeGet_changeTime(), nil)
	}
	var flags int32 = this_.nodeGet_propFlags()
	_ = flags
	if flags != int32(0) {
		var se int32 = ((javalang.Ushr32(flags, int(int32(28)))) & int32(15))
		_ = se
		if se != int32(0) {
			return StaticGet_Node_STATIC_SIDE_EFFECT_ITEMS()[se]
		}
		var cv int32 = ((javalang.Ushr32(flags, int(int32(26)))) & int32(3))
		_ = cv
		if cv != int32(0) {
			return StaticGet_Node_STATIC_CONST_VAR_ITEMS()[cv]
		}
		var boolBits int32 = (flags & int32(67108863))
		_ = boolBits
		if boolBits != int32(0) {
			return StaticGet_Node_STATIC_INLINE_BOOL_ITEMS()[j2gort.IntTrailingZeros(boolBits)]
		}
	}
	if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
		return Node_PropListItem_ObjectItem(int8((StaticGet_Node_Prop_SOURCE_FILE()).Ordinal()), this_.nodeGet_sourceFile(), nil)
	}
	return nil
}

// Java method: Node.applyPropListHeadItems(com.google.javascript.rhino.Node.@org.jspecify.annotations.Nullable PropListItem item)
func (this_ Node) priv_Node_applyPropListHeadItems_8499(item *Node_PropListItem) {
	_ = item
	if item == nil {
		return
	}
	this_.priv_Node_applyPropListHeadItems_8499(item.next)
	var prop Node_Prop = StaticGet_Node_PropValues()[int32(item.propType)]
	_ = prop
	if item.IsObjectItem() {
		this_.PutProp(prop, item.GetObjectValue())
	} else {
		this_.PutIntProp(prop, item.GetIntValue())
	}
}

// Java method: Node.setPropListHead(com.google.javascript.rhino.Node.@org.jspecify.annotations.Nullable PropListItem propListHead)
func (this_ Node) SetPropListHead(propListHead *Node_PropListItem) {
	_ = propListHead
	this_.nodeSet_sourceFile(nil)
	this_.nodeSet_jsDocInfo(nil)
	this_.nodeSet_typeBeforeCast(nil)
	this_.nodeSet_propFlags(int32(0))
	this_.nodeSet_changeTime(int32(0))
	this_.nodeSet_propListHead(nil)
	this_.priv_Node_applyPropListHeadItems_8499(propListHead)
}

// Java method: Node.getParent()
func (this_ Node) GetParent() Node {
	return this_.nodeGet_parent()
}

// Java method: Node.hasParent()
func (this_ Node) HasParent() bool {
	return (this_.nodeGet_parent() != 0)
}

// Java method: Node.getGrandparent()
func (this_ Node) GetGrandparent() Node {
	if this_.nodeGet_parent() == 0 {
		return Node(0)
	} else {
		return this_.nodeGet_parent().nodeGet_parent()
	}
}

// Java method: Node.getAncestor(int level)
func (this_ Node) GetAncestor(level int32) Node {
	_ = level
	if !(level >= int32(0)) {
		j2gort.CheckArgument(false)
	}
	var node Node = this_
	_ = node
	for (node != 0) && (func() int32 { old_97 := level; nw_98 := old_97 - 1; level = nw_98; return old_97 }() > int32(0)) {
		node = node.GetParent()
	}
	return node
}

// Java method: Node.isDescendantOf(com.google.javascript.rhino.Node node)
func (this_ Node) IsDescendantOf(node Node) bool {
	_ = node
	{
		var n Node = this_
		_ = n
		for ; n != 0; n = n.nodeGet_parent() {
			if n == node {
				return true
			}
		}
	}
	return false
}

// Java method: Node.isOnlyChildOf(com.google.javascript.rhino.Node possibleParent)
func (this_ Node) IsOnlyChildOf(possibleParent Node) bool {
	_ = possibleParent
	return (((possibleParent == this_.GetParent()) && (this_.GetPrevious() == 0)) && (this_.GetNext() == 0))
}

// Java method: Node.isFirstChildOf(com.google.javascript.rhino.Node possibleParent)
func (this_ Node) IsFirstChildOf(possibleParent Node) bool {
	_ = possibleParent
	return ((possibleParent == this_.GetParent()) && (this_.GetPrevious() == 0))
}

// Java method: Node.isSecondChildOf(com.google.javascript.rhino.Node possibleParent)
func (this_ Node) IsSecondChildOf(possibleParent Node) bool {
	_ = possibleParent
	var previousNode Node = this_.GetPrevious()
	_ = previousNode
	return ((previousNode != 0) && previousNode.IsFirstChildOf(possibleParent))
}

// Java method: Node.getAncestors()
func (this_ Node) GetAncestors() *Node_AncestorIterable {
	return NewNode_AncestorIterable(this_.GetParent())
}

// Java method: Node.hasOneChild()
func (this_ Node) HasOneChild() bool {
	return ((this_.nodeGet_first() != 0) && (this_.nodeGet_first().nodeGet_next() == 0))
}

// Java method: Node.hasTwoChildren()
func (this_ Node) HasTwoChildren() bool {
	return (((this_.nodeGet_first() != 0) && (this_.nodeGet_first().nodeGet_next() != 0)) && (this_.nodeGet_first().nodeGet_next() == this_.GetLastChild()))
}

// Java method: Node.hasZeroOrOneChild()
func (this_ Node) HasZeroOrOneChild() bool {
	return (this_.nodeGet_first() == this_.GetLastChild())
}

// Java method: Node.hasMoreThanOneChild()
func (this_ Node) HasMoreThanOneChild() bool {
	return ((this_.nodeGet_first() != 0) && (this_.nodeGet_first().nodeGet_next() != 0))
}

// Java method: Node.hasXChildren(int x)
func (this_ Node) HasXChildren(x int32) bool {
	_ = x
	var c int32 = int32(0)
	_ = c
	{
		var n Node = this_.nodeGet_first()
		_ = n
		for ; (n != 0) && (c <= x); n = n.nodeGet_next() {
			c = c + 1
		}
	}
	return (c == x)
}

// Java method: Node.getChildCount()
func (this_ Node) GetChildCount() int32 {
	var c int32 = int32(0)
	_ = c
	{
		var n Node = this_.nodeGet_first()
		_ = n
		for ; n != 0; n = n.nodeGet_next() {
			c = c + 1
		}
	}
	return c
}

// Java method: Node.hasChild(com.google.javascript.rhino.Node child)
func (this_ Node) HasChild(child Node) bool {
	_ = child
	{
		var n Node = this_.nodeGet_first()
		_ = n
		for ; n != 0; n = n.nodeGet_next() {
			if child == n {
				return true
			}
		}
	}
	return false
}

// Java method: Node.isEquivalentToShallow(com.google.javascript.rhino.Node node)
func (this_ Node) IsEquivalentToShallow(node Node) bool {
	_ = node
	return this_.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, StaticGet_Node_RecursionMode_SHALLOW(), StaticGet_Node_TypeComparison_IGNORE(), StaticGet_Node_JsDocComparison_IGNORE(), StaticGet_Node_SideEffectComparison_IGNORE())
}

// Java method: Node.isEquivalentWithSideEffectsTo(com.google.javascript.rhino.Node node)
func (this_ Node) IsEquivalentWithSideEffectsTo(node Node) bool {
	_ = node
	return this_.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, StaticGet_Node_RecursionMode_DEEP(), StaticGet_Node_TypeComparison_IGNORE(), StaticGet_Node_JsDocComparison_IGNORE(), StaticGet_Node_SideEffectComparison_COMPARE())
}

// Java method: Node.isEquivalentWithSideEffectsToShallow(com.google.javascript.rhino.Node node)
func (this_ Node) IsEquivalentWithSideEffectsToShallow(node Node) bool {
	_ = node
	return this_.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, StaticGet_Node_RecursionMode_SHALLOW(), StaticGet_Node_TypeComparison_IGNORE(), StaticGet_Node_JsDocComparison_IGNORE(), StaticGet_Node_SideEffectComparison_COMPARE())
}

// Java method: Node.isEquivalentToTyped(com.google.javascript.rhino.Node node)
func (this_ Node) IsEquivalentToTyped(node Node) bool {
	_ = node
	return this_.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, StaticGet_Node_RecursionMode_DEEP(), StaticGet_Node_TypeComparison_COMPARE(), StaticGet_Node_JsDocComparison_COMPARE(), StaticGet_Node_SideEffectComparison_IGNORE())
}

// Java method: Node.isEquivalentTo(com.google.javascript.rhino.Node node)
func (this_ Node) IsEquivalentTo_1(node Node) bool {
	_ = node
	return this_.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, StaticGet_Node_RecursionMode_DEEP(), StaticGet_Node_TypeComparison_IGNORE(), StaticGet_Node_JsDocComparison_IGNORE(), StaticGet_Node_SideEffectComparison_IGNORE())
}

// Java method: Node.getSecondToLastPropListItem()
func (this_ Node) priv_Node_getSecondToLastPropListItem_8522() *Node_PropListItem {
	if this_.GetIsInClosureUnawareSubtree() {
		return this_.LookupProperty(StaticGet_Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE())
	}
	var propListItem *Node_PropListItem = this_.nodeGet_propListHead()
	_ = propListItem
	if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
		if propListItem == nil {
			return nil
		}
		for propListItem.next != nil {
			propListItem = propListItem.next
		}
		return propListItem
	}
	if (propListItem == nil) || (propListItem.next == nil) {
		return nil
	}
	for propListItem.next.next != nil {
		propListItem = propListItem.next
	}
	return propListItem
}

// Java method: Node.getLastPropListItem()
func (this_ Node) priv_Node_getLastPropListItem_8523() *Node_PropListItem {
	if !(j2gort.IfaceEquals[StaticSourceFile](this_.nodeGet_sourceFile(), nil)) {
		return Node_PropListItem_ObjectItem(int8((StaticGet_Node_Prop_SOURCE_FILE()).Ordinal()), this_.nodeGet_sourceFile(), nil)
	}
	var propListItem *Node_PropListItem = this_.nodeGet_propListHead()
	_ = propListItem
	if propListItem == nil {
		return nil
	}
	for propListItem.next != nil {
		propListItem = propListItem.next
	}
	return propListItem
}

// Java method: Node.validatePropListTailOrdering(boolean shouldNodeBeClosureUnaware)
func (this_ Node) priv_Node_validatePropListTailOrdering_8524(shouldNodeBeClosureUnaware bool) bool {
	_ = shouldNodeBeClosureUnaware
	if j2gort.IfaceEquals[StaticSourceFile](this_.GetStaticSourceFile(), nil) {
		return false
	}
	if shouldNodeBeClosureUnaware && !(this_.GetIsInClosureUnawareSubtree()) {
		return false
	}
	return true
}

// Java static method: Node.validateMemorySensitivePropertyGuarantees(com.google.javascript.rhino.Node node1, boolean shouldNode1BeClosureUnaware, com.google.javascript.rhino.Node node2, boolean shouldNode2BeClosureUnaware)
func Node_ValidateMemorySensitivePropertyGuarantees(node1 Node, shouldNode1BeClosureUnaware bool, node2 Node, shouldNode2BeClosureUnaware bool) bool {
	clinit_Node()
	_ = node1
	_ = shouldNode1BeClosureUnaware
	_ = node2
	_ = shouldNode2BeClosureUnaware
	if !(node1.priv_Node_validatePropListTailOrdering_8524(shouldNode1BeClosureUnaware)) {
		return false
	}
	if !(node2.priv_Node_validatePropListTailOrdering_8524(shouldNode2BeClosureUnaware)) {
		return false
	}
	if !(j2gort.IfaceEquals[StaticSourceFile](node1.GetStaticSourceFile(), node2.GetStaticSourceFile())) {
		return false
	}
	if shouldNode1BeClosureUnaware && shouldNode2BeClosureUnaware {
		if node1.GetIsInClosureUnawareSubtree() != node2.GetIsInClosureUnawareSubtree() {
			return false
		}
	}
	return true
}

// Java method: Node.isEquivalentTo(com.google.javascript.rhino.Node node, boolean compareType, boolean recurse, boolean jsDoc)
func (this_ Node) IsEquivalentTo_4(node Node, compareType bool, recurse bool, jsDoc bool) bool {
	_ = node
	_ = compareType
	_ = recurse
	_ = jsDoc
	return this_.IsEquivalentTo_5_Node_bool_bool_bool_bool(node, compareType, recurse, jsDoc, false)
}

// Java method: Node.isEquivalentTo(com.google.javascript.rhino.Node node, boolean compareType, boolean recurse, boolean jsDoc, boolean sideEffect)
func (this_ Node) IsEquivalentTo_5_Node_bool_bool_bool_bool(node Node, compareType bool, recurse bool, jsDoc bool, sideEffect bool) bool {
	_ = node
	_ = compareType
	_ = recurse
	_ = jsDoc
	_ = sideEffect
	return this_.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, func() Node_RecursionMode {
		if recurse {
			return StaticGet_Node_RecursionMode_DEEP()
		}
		return StaticGet_Node_RecursionMode_SHALLOW()
	}(), func() Node_TypeComparison {
		if compareType {
			return StaticGet_Node_TypeComparison_COMPARE()
		}
		return StaticGet_Node_TypeComparison_IGNORE()
	}(), func() Node_JsDocComparison {
		if jsDoc {
			return StaticGet_Node_JsDocComparison_COMPARE()
		}
		return StaticGet_Node_JsDocComparison_IGNORE()
	}(), func() Node_SideEffectComparison {
		if sideEffect {
			return StaticGet_Node_SideEffectComparison_COMPARE()
		}
		return StaticGet_Node_SideEffectComparison_IGNORE()
	}())
}

// Java method: Node.nodeSubclassKind()
func (this_ Node) base_NodeSubclassKind() int32 {
	return int32(0)
}

func (this_ Node) NodeSubclassKind() int32 {
	switch nodestore.Node(this_).KindFlags() {
	case 1:
		return this_.nodeSub1_NodeSubclassKind()
	case 2:
		return this_.nodeSub2_NodeSubclassKind()
	case 3:
		return this_.nodeSub3_NodeSubclassKind()
	case 4:
		return this_.nodeSub4_NodeSubclassKind()
	default:
		return this_.base_NodeSubclassKind()
	}
}

// Java method: Node.isEquivalentTo(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node.RecursionMode recursion, com.google.javascript.rhino.Node.TypeComparison typeConfig, com.google.javascript.rhino.Node.JsDocComparison jsDocConfig, com.google.javascript.rhino.Node.SideEffectComparison sideEffectConfig)
func (this_ Node) base_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node Node, recursion Node_RecursionMode, typeConfig Node_TypeComparison, jsDocConfig Node_JsDocComparison, sideEffectConfig Node_SideEffectComparison) bool {
	_ = node
	_ = recursion
	_ = typeConfig
	_ = jsDocConfig
	_ = sideEffectConfig
	if ((this_.RawToken_() != node.RawToken_()) || (this_.GetChildCount() != node.GetChildCount())) || (this_.NodeSubclassKind() != node.NodeSubclassKind()) {
		return false
	}
	if (typeConfig == StaticGet_Node_TypeComparison_COMPARE()) && !(javalang.ObjectsEquals(this_.nodeGet_jstypeOrColor(), node.nodeGet_jstypeOrColor())) {
		return false
	}
	if (jsDocConfig == StaticGet_Node_JsDocComparison_COMPARE()) && !(JSDocInfo_AreEquivalent(this_.GetJSDocInfo(), node.GetJSDocInfo())) {
		return false
	}
	var thisDte Node = this_.GetDeclaredTypeExpression()
	_ = thisDte
	var thatDte Node = node.GetDeclaredTypeExpression()
	_ = thatDte
	if thisDte == thatDte {
	} else if (thisDte == 0) || (thatDte == 0) {
		return false
	} else if !(thisDte.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(thatDte, recursion, typeConfig, jsDocConfig, sideEffectConfig)) {
		return false
	}
	var eqMask int32
	_ = eqMask
	if this_.IsFunction() {
		eqMask = int32(42000774)
	} else {
		eqMask = (int32(42000774) & ^(int32(128) | int32(32768)))
	}
	if ((this_.nodeGet_propFlags() ^ node.nodeGet_propFlags()) & eqMask) != int32(0) {
		return false
	}
	var thisHead *Node_PropListItem = this_.nodeGet_propListHead()
	_ = thisHead
	var nodeHead *Node_PropListItem = node.nodeGet_propListHead()
	_ = nodeHead
	if (thisHead != nil) || (nodeHead != nil) {
		var coldPropBits int64 = int64(0)
		_ = coldPropBits
		{
			var propListItem *Node_PropListItem = thisHead
			_ = propListItem
			for ; propListItem != nil; propListItem = propListItem.next {
				coldPropBits = (coldPropBits | javalang.Shl64(int64(1), int(int32(propListItem.propType))))
			}
		}
		{
			var propListItem_2 *Node_PropListItem = nodeHead
			_ = propListItem_2
			for ; propListItem_2 != nil; propListItem_2 = propListItem_2.next {
				coldPropBits = (coldPropBits | javalang.Shl64(int64(1), int(int32(propListItem_2.propType))))
			}
		}
		coldPropBits = (coldPropBits & StaticGet_Node_PROP_EQUALITY_ORDINAL_MASK())
		for coldPropBits != int64(0) {
			var ord int32 = j2gort.BoxLong_NumberOfTrailingZeros(coldPropBits)
			_ = ord
			coldPropBits = (coldPropBits & ^(javalang.Shl64(int64(1), int(ord))))
			if !(Node_PropEqualsForEquivalence(StaticGet_Node_PropValues()[ord], this_, node)) {
				return false
			}
		}
	}
	if sideEffectConfig == StaticGet_Node_SideEffectComparison_COMPARE() {
		if this_.GetSideEffectFlags() != node.GetSideEffectFlags() {
			return false
		}
		if this_.IsUnusedParameter() != node.IsUnusedParameter() {
			return false
		}
	}
	var thisIsShadowHost bool = (this_.GetClosureUnawareShadow() != 0)
	_ = thisIsShadowHost
	var nodeIsShadowHost bool = (node.GetClosureUnawareShadow() != 0)
	_ = nodeIsShadowHost
	if thisIsShadowHost != nodeIsShadowHost {
		return false
	}
	if (recursion.RecurseShadowAst() && thisIsShadowHost) && nodeIsShadowHost {
		var thisShadow Node = this_.GetClosureUnawareShadow()
		_ = thisShadow
		var nodeShadow Node = node.GetClosureUnawareShadow()
		_ = nodeShadow
		if !(thisShadow.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(nodeShadow, recursion, typeConfig, jsDocConfig, sideEffectConfig)) {
			return false
		}
	}
	if recursion.RecurseChildren() {
		{
			var n Node = this_.nodeGet_first()
			_ = n
			var n2 Node = node.nodeGet_first()
			_ = n2
			for ; n != 0; func() { n = n.nodeGet_next(); n2 = n2.nodeGet_next() }() {
				if !(n.IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(n2, recursion, typeConfig, jsDocConfig, sideEffectConfig)) {
					return false
				}
			}
		}
	}
	return true
}

func (this_ Node) IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(p0_ Node, p1_ Node_RecursionMode, p2_ Node_TypeComparison, p3_ Node_JsDocComparison, p4_ Node_SideEffectComparison) bool {
	switch nodestore.Node(this_).KindFlags() {
	case 1:
		return this_.nodeSub1_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(p0_, p1_, p2_, p3_, p4_)
	case 2:
		return this_.nodeSub2_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(p0_, p1_, p2_, p3_, p4_)
	case 3:
		return this_.nodeSub3_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(p0_, p1_, p2_, p3_, p4_)
	case 4:
		return this_.nodeSub4_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(p0_, p1_, p2_, p3_, p4_)
	default:
		return this_.base_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(p0_, p1_, p2_, p3_, p4_)
	}
}

// Java static method: Node.propEqualsForEquivalence(com.google.javascript.rhino.Node.Prop prop, com.google.javascript.rhino.Node a, com.google.javascript.rhino.Node b)
func Node_PropEqualsForEquivalence(prop Node_Prop, a Node, b Node) bool {
	clinit_Node()
	_ = prop
	_ = a
	_ = b
	switch (prop).Ordinal() {
	case int32(20) /* Prop.ARROW_FN */ :
		return (a.IsArrowFunction() == b.IsArrowFunction())
	case int32(21) /* Prop.ASYNC_FN */ :
		return (a.IsAsyncFunction() == b.IsAsyncFunction())
	case int32(19) /* Prop.GENERATOR_FN */ :
		return (a.IsGeneratorFunction() == b.IsGeneratorFunction())
	case int32(51) /* Prop.START_OF_OPT_CHAIN */ :
		return (a.IsOptionalChainStart() == b.IsOptionalChainStart())
	case int32(18) /* Prop.STATIC_MEMBER */ :
		return (a.IsStaticMember() == b.IsStaticMember())
	case int32(22) /* Prop.YIELD_ALL */ :
		return (a.IsYieldAll() == b.IsYieldAll())
	case int32(23) /* Prop.EXPORT_DEFAULT */, int32(24) /* Prop.EXPORT_ALL_FROM */, int32(4) /* Prop.INCRDECR */, int32(5) /* Prop.QUOTED */ :
		return (a.priv_Node_getIntProp_8427(prop) == b.priv_Node_getIntProp_8427(prop))
	case int32(13) /* Prop.FREE_CALL */, int32(12) /* Prop.DIRECT_EVAL */, int32(28) /* Prop.COMPUTED_PROP_METHOD */, int32(29) /* Prop.COMPUTED_PROP_GETTER */, int32(30) /* Prop.COMPUTED_PROP_SETTER */ :
		return (a.GetBooleanProp(prop) == b.GetBooleanProp(prop))
	default:
		return true
	}
	panic("unreachable switch expression")
}

// Java static method: Node.computePropEqualityOrdinalMask()
func Node_ComputePropEqualityOrdinalMask() int64 {
	clinit_Node()
	if !(int32(len(StaticGet_Node_PropValues())) <= int32(64)) {
		j2gort.CheckState(false)
	}
	return (((((((((((((((javalang.Shl64(int64(1), int((StaticGet_Node_Prop_ARROW_FN()).Ordinal()))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_ASYNC_FN()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_GENERATOR_FN()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_START_OF_OPT_CHAIN()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_STATIC_MEMBER()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_YIELD_ALL()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_EXPORT_DEFAULT()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_EXPORT_ALL_FROM()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_INCRDECR()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_QUOTED()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_FREE_CALL()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_DIRECT_EVAL()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_COMPUTED_PROP_METHOD()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_COMPUTED_PROP_GETTER()).Ordinal())))) | (javalang.Shl64(int64(1), int((StaticGet_Node_Prop_COMPUTED_PROP_SETTER()).Ordinal()))))
}

// Java method: Node.getQualifiedName()
func (this_ Node) GetQualifiedName() string {
	switch this_.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		var name string = this_.GetString()
		_ = name
		if j2gort.CheckNotNullString(name) == "" {
			return j2gort.NullString
		}
		return name
	case int32(26) /* Token.GETPROP */ :
		var p1 string = this_.GetString()
		_ = p1
		var n1 Node = this_.nodeGet_first()
		_ = n1
		if !(n1.IsGetProp()) {
			var left string = n1.GetQualifiedName()
			_ = left
			if left == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot(j2gort.StrValString(left), j2gort.StrValString(p1))
		}
		var p2 string = n1.GetString()
		_ = p2
		var n2 Node = n1.nodeGet_first()
		_ = n2
		if !(n2.IsGetProp()) {
			var left_2 string = n2.GetQualifiedName()
			_ = left_2
			if left_2 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot3(j2gort.StrValString(left_2), j2gort.StrValString(p2), j2gort.StrValString(p1))
		}
		var p3 string = n2.GetString()
		_ = p3
		var n3 Node = n2.nodeGet_first()
		_ = n3
		if !(n3.IsGetProp()) {
			var left_3 string = n3.GetQualifiedName()
			_ = left_3
			if left_3 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot4(j2gort.StrValString(left_3), j2gort.StrValString(p3), j2gort.StrValString(p2), j2gort.StrValString(p1))
		}
		var p4 string = n3.GetString()
		_ = p4
		var n4 Node = n3.nodeGet_first()
		_ = n4
		if !(n4.IsGetProp()) {
			var left_4 string = n4.GetQualifiedName()
			_ = left_4
			if left_4 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot(nodestore.JoinDot4(j2gort.StrValString(left_4), j2gort.StrValString(p4), j2gort.StrValString(p3), j2gort.StrValString(p2)), j2gort.StrValString(p1))
		}
		var p5 string = n4.GetString()
		_ = p5
		var n5 Node = n4.nodeGet_first()
		_ = n5
		if !(n5.IsGetProp()) {
			var left_5 string = n5.GetQualifiedName()
			_ = left_5
			if left_5 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot3(nodestore.JoinDot4(j2gort.StrValString(left_5), j2gort.StrValString(p5), j2gort.StrValString(p4), j2gort.StrValString(p3)), j2gort.StrValString(p2), j2gort.StrValString(p1))
		}
		var p6 string = n5.GetString()
		_ = p6
		var n6 Node = n5.nodeGet_first()
		_ = n6
		if !(n6.IsGetProp()) {
			var left_6 string = n6.GetQualifiedName()
			_ = left_6
			if left_6 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot4(nodestore.JoinDot4(j2gort.StrValString(left_6), j2gort.StrValString(p6), j2gort.StrValString(p5), j2gort.StrValString(p4)), j2gort.StrValString(p3), j2gort.StrValString(p2), j2gort.StrValString(p1))
		}
		var p7 string = n6.GetString()
		_ = p7
		var n7 Node = n6.nodeGet_first()
		_ = n7
		if !(n7.IsGetProp()) {
			var left_7 string = n7.GetQualifiedName()
			_ = left_7
			if left_7 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot(nodestore.JoinDot4(nodestore.JoinDot4(j2gort.StrValString(left_7), j2gort.StrValString(p7), j2gort.StrValString(p6), j2gort.StrValString(p5)), j2gort.StrValString(p4), j2gort.StrValString(p3), j2gort.StrValString(p2)), j2gort.StrValString(p1))
		}
		var p8 string = n7.GetString()
		_ = p8
		var n8 Node = n7.nodeGet_first()
		_ = n8
		if !(n8.IsGetProp()) {
			var left_8 string = n8.GetQualifiedName()
			_ = left_8
			if left_8 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot3(nodestore.JoinDot4(nodestore.JoinDot4(j2gort.StrValString(left_8), j2gort.StrValString(p8), j2gort.StrValString(p7), j2gort.StrValString(p6)), j2gort.StrValString(p5), j2gort.StrValString(p4), j2gort.StrValString(p3)), j2gort.StrValString(p2), j2gort.StrValString(p1))
		}
		var p9 string = n8.GetString()
		_ = p9
		var n9 Node = n8.nodeGet_first()
		_ = n9
		if !(n9.IsGetProp()) {
			var left_9 string = n9.GetQualifiedName()
			_ = left_9
			if left_9 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot4(nodestore.JoinDot4(nodestore.JoinDot4(j2gort.StrValString(left_9), j2gort.StrValString(p9), j2gort.StrValString(p8), j2gort.StrValString(p7)), j2gort.StrValString(p6), j2gort.StrValString(p5), j2gort.StrValString(p4)), j2gort.StrValString(p3), j2gort.StrValString(p2), j2gort.StrValString(p1))
		}
		var p10 string = n9.GetString()
		_ = p10
		var n10 Node = n9.nodeGet_first()
		_ = n10
		if !(n10.IsGetProp()) {
			var left_10 string = n10.GetQualifiedName()
			_ = left_10
			if left_10 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot(nodestore.JoinDot4(nodestore.JoinDot4(nodestore.JoinDot4(j2gort.StrValString(left_10), j2gort.StrValString(p10), j2gort.StrValString(p9), j2gort.StrValString(p8)), j2gort.StrValString(p7), j2gort.StrValString(p6), j2gort.StrValString(p5)), j2gort.StrValString(p4), j2gort.StrValString(p3), j2gort.StrValString(p2)), j2gort.StrValString(p1))
		}
		var p11 string = n10.GetString()
		_ = p11
		var n11 Node = n10.nodeGet_first()
		_ = n11
		if !(n11.IsGetProp()) {
			var left_11 string = n11.GetQualifiedName()
			_ = left_11
			if left_11 == j2gort.NullString {
				return j2gort.NullString
			}
			return nodestore.JoinDot3(nodestore.JoinDot4(nodestore.JoinDot4(nodestore.JoinDot4(j2gort.StrValString(left_11), j2gort.StrValString(p11), j2gort.StrValString(p10), j2gort.StrValString(p9)), j2gort.StrValString(p8), j2gort.StrValString(p7), j2gort.StrValString(p6)), j2gort.StrValString(p5), j2gort.StrValString(p4), j2gort.StrValString(p3)), j2gort.StrValString(p2), j2gort.StrValString(p1))
		}
		var p12 string = n11.GetString()
		_ = p12
		var n12 Node = n11.nodeGet_first()
		_ = n12
		var left_12 string = n12.GetQualifiedName()
		_ = left_12
		if left_12 == j2gort.NullString {
			return j2gort.NullString
		}
		return nodestore.JoinDot4(nodestore.JoinDot4(nodestore.JoinDot4(nodestore.JoinDot4(j2gort.StrValString(left_12), j2gort.StrValString(p12), j2gort.StrValString(p11), j2gort.StrValString(p10)), j2gort.StrValString(p9), j2gort.StrValString(p8), j2gort.StrValString(p7)), j2gort.StrValString(p6), j2gort.StrValString(p5), j2gort.StrValString(p4)), j2gort.StrValString(p3), j2gort.StrValString(p2), j2gort.StrValString(p1))
	case int32(37) /* Token.THIS */ :
		return "this"
	case int32(110) /* Token.SUPER */ :
		return "super"
	default:
		return j2gort.NullString
	}
	panic("unreachable switch expression")
}

// Java method: Node.getQualifiedNameObject()
func (this_ Node) GetQualifiedNameObject() QualifiedName {
	if this_.IsQualifiedName() {
		return NewQualifiedName_NodeQname(this_)
	} else {
		return nil
	}
}

// Java method: Node.getOriginalQualifiedName()
func (this_ Node) GetOriginalQualifiedName() string {
	if this_.RawToken_() == uint8(32) {
		var name string = this_.GetOriginalName()
		_ = name
		if name == j2gort.NullString {
			name = this_.GetString()
		}
		if j2gort.CheckNotNullString(name) == "" {
			return j2gort.NullString
		} else {
			return name
		}
	} else if this_.RawToken_() == uint8(26) {
		var left string = this_.GetFirstChild().GetOriginalQualifiedName()
		_ = left
		if left == j2gort.NullString {
			return j2gort.NullString
		}
		var right string = this_.GetOriginalName()
		_ = right
		if right == j2gort.NullString {
			right = this_.GetString()
		}
		return nodestore.JoinDot(j2gort.StrValString(left), j2gort.StrValString(right))
	} else if this_.RawToken_() == uint8(37) {
		return "this"
	} else if this_.RawToken_() == uint8(110) {
		return "super"
	} else {
		return j2gort.NullString
	}
}

// Java method: Node.isQualifiedName()
func (this_ Node) IsQualifiedName() bool {
	switch this_.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return !(j2gort.CheckNotNullString(this_.GetString()) == "")
	case int32(37) /* Token.THIS */, int32(110) /* Token.SUPER */ :
		return true
	case int32(26) /* Token.GETPROP */ :
		return this_.GetFirstChild().IsQualifiedName()
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		fallthrough
	default:
		return false
	}
	panic("unreachable")
}

// Java static method: Node.initPow31()
func Node_InitPow31() []int32 {
	clinit_Node()
	var table []int32 = make([]int32, int32(64))
	_ = table
	var p int32 = int32(1)
	_ = p
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(64); i = i + 1 {
			{
				a_99 := table
				i_100 := i
				a_99[i_100] = p
			}
			p = (p * int32(31))
		}
	}
	return table
}

// Java static method: Node.pow31(int n)
func Node_Pow31(n int32) int32 {
	clinit_Node()
	_ = n
	if n < int32(64) {
		return StaticGet_Node_POW31()[n]
	}
	var result int32 = int32(1)
	_ = result
	var base int32 = int32(31)
	_ = base
	var exp int32 = n
	_ = exp
	for exp > int32(0) {
		if (exp & int32(1)) != int32(0) {
			result = (result * base)
		}
		base = (base * base)
		exp = javalang.Ushr32(exp, int(int32(1)))
	}
	return result
}

// Java static method: Node.mixQualifiedHash(int h)
func Node_MixQualifiedHash(h int32) int32 {
	clinit_Node()
	_ = h
	h = (h ^ javalang.Ushr32(h, int(int32(16))))
	h = (h * int32(-2048144789))
	h = (h ^ javalang.Ushr32(h, int(int32(13))))
	h = (h * int32(-1028477387))
	h = (h ^ javalang.Ushr32(h, int(int32(16))))
	if h == int32(0) {
		return int32(1)
	} else {
		return h
	}
}

// Java static method: Node.hashQualifiedNameString(java.lang.@org.jspecify.annotations.Nullable String s)
func Node_HashQualifiedNameString(s string) int32 {
	clinit_Node()
	_ = s
	if (s == j2gort.NullString) || (j2gort.CheckNotNullString(s) == "") {
		return int32(0)
	}
	return Node_MixQualifiedHash(javalang.StringHashCode(j2gort.CheckNotNullString(s)))
}

// Java method: Node.hashQualifiedName()
func (this_ Node) HashQualifiedName() int32 {
	return int32((javalang.Ushr64(this_.QualifiedNameHashes(), int(int32(32)))))
}

// Java method: Node.qualifiedNameHashes()
func (this_ Node) QualifiedNameHashes() int64 {
	var suffixHash int32 = int32(0)
	_ = suffixHash
	var suffixPow int32 = int32(1)
	_ = suffixPow
	var cur Node = this_
	_ = cur
	for cur.RawToken_() == uint8(26) {
		var prop string = cur.GetString()
		_ = prop
		var p int32 = Node_Pow31(j2gort.JStrLen(prop))
		_ = p
		suffixHash = (suffixHash + (((int32(46) * p) + javalang.StringHashCode(j2gort.CheckNotNullString(prop))) * suffixPow))
		suffixPow = (suffixPow * (int32(31) * p))
		cur = cur.nodeGet_first()
		if cur == 0 {
			return int64(0)
		}
	}
	var rawHash int32 = int32(0)
	_ = rawHash
	switch cur.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		var root string = cur.GetString()
		_ = root
		if j2gort.CheckNotNullString(root) == "" {
			return int64(0)
		}
		rawHash = ((javalang.StringHashCode(j2gort.CheckNotNullString(root)) * suffixPow) + suffixHash)
	case int32(37) /* Token.THIS */ :
		rawHash = ((int32(3559070) * suffixPow) + suffixHash)
	case int32(110) /* Token.SUPER */ :
		rawHash = ((int32(109801339) * suffixPow) + suffixHash)
	default:
		return int64(0)
	}
	var mixedHash int32 = Node_MixQualifiedHash(rawHash)
	_ = mixedHash
	return ((javalang.Shl64((int64(mixedHash)), int(int32(32)))) | (int64(rawHash) & int64(4294967295)))
}

// Java method: Node.matchesName(java.lang.String name)
func (this_ Node) MatchesName_1_string(name string) bool {
	_ = name
	if this_.RawToken_() != uint8(32) {
		return false
	}
	var internalString string = this_.GetString()
	_ = internalString
	return (!(j2gort.CheckNotNullString(internalString) == "") && (j2gort.CheckNotNullString(name) == internalString))
}

// Java method: Node.matchesName(com.google.javascript.rhino.Node n)
func (this_ Node) MatchesName_1_Node(n Node) bool {
	_ = n
	if (this_.RawToken_() != uint8(32)) || (n.RawToken_() != uint8(32)) {
		return false
	}
	var internalString string = this_.GetString()
	_ = internalString
	return (!(j2gort.CheckNotNullString(internalString) == "") && RhinoStringPool_UncheckedEquals(internalString, n.GetString()))
}

// Java method: Node.matchesQualifiedName(java.lang.String name)
func (this_ Node) MatchesQualifiedName_1_string(name string) bool {
	_ = name
	return this_.priv_Node_matchesQualifiedName_8538(name, j2gort.JStrLen(name))
}

// Java method: Node.matchesQualifiedName(java.lang.String qname, int endIndex)
func (this_ Node) priv_Node_matchesQualifiedName_8538(qname string, endIndex int32) bool {
	_ = qname
	_ = endIndex
	var start int32 = (j2gort.JStrLastIndexOfCpFrom(qname, int32(uint16(46)), (endIndex-int32(1))) + int32(1))
	_ = start
	{
		var name string = j2gort.NullString
		_ = name
		var prop string = j2gort.NullString
		_ = prop
		switch this_.TokenSwitchOrd_() {
		case int32(32) /* Token.NAME */, int32(119) /* Token.IMPORT_STAR */ :
			name = this_.GetString()
			return ((((start == int32(0)) && !(j2gort.CheckNotNullString(name) == "")) && (j2gort.JStrLen(name) == endIndex)) && javalang.StringStartsWith(j2gort.CheckNotNullString(qname), j2gort.CheckNotNullString(name)))
		case int32(37) /* Token.THIS */ :
			return (((start == int32(0)) && (endIndex == int32(4))) && javalang.StringStartsWith(j2gort.CheckNotNullString(qname), "this"))
		case int32(110) /* Token.SUPER */ :
			return (((start == int32(0)) && (endIndex == int32(5))) && javalang.StringStartsWith(j2gort.CheckNotNullString(qname), "super"))
		case int32(26) /* Token.GETPROP */ :
			prop = this_.GetString()
			return ((((start > int32(1)) && (j2gort.JStrLen(prop) == (endIndex - start))) && j2gort.JStrRegionMatches(prop, false, int32(0), qname, start, (endIndex-start))) && this_.GetFirstChild().priv_Node_matchesQualifiedName_8538(qname, (start-int32(1))))
		case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
			fallthrough
		default:
			return false
		}
	}
	panic("unreachable")
}

// Java method: Node.matchesQualifiedName(com.google.javascript.rhino.Node n)
func (this_ Node) MatchesQualifiedName_1_Node(n Node) bool {
	_ = n
	if n.RawToken_() != this_.RawToken_() {
		return false
	}
	switch this_.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return this_.MatchesName_1_Node(n)
	case int32(37) /* Token.THIS */, int32(110) /* Token.SUPER */ :
		return true
	case int32(26) /* Token.GETPROP */ :
		return (RhinoStringPool_UncheckedEquals(this_.GetString(), n.GetString()) && this_.GetFirstChild().MatchesQualifiedName_1_Node(n.GetFirstChild()))
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		fallthrough
	default:
		return false
	}
	panic("unreachable")
}

// Java method: Node.isUnscopedQualifiedName()
func (this_ Node) IsUnscopedQualifiedName() bool {
	switch this_.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return !(j2gort.CheckNotNullString(this_.GetString()) == "")
	case int32(26) /* Token.GETPROP */ :
		return this_.GetFirstChild().IsUnscopedQualifiedName()
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Node.isValidAssignmentTarget()
func (this_ Node) IsValidAssignmentTarget() bool {
	switch this_.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */, int32(102) /* Token.ARRAY_PATTERN */, int32(103) /* Token.OBJECT_PATTERN */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Node.clone()
func (this_ Node) Clone() any {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Did you mean cloneNode?")))
}

// Java method: Node.cloneNode()
func (this_ Node) CloneNode_0() Node {
	return this_.CloneNode_1(false)
}

// Java method: Node.cloneNode(boolean cloneTypeExprs)
func (this_ Node) base_CloneNode_1(cloneTypeExprs bool) Node {
	_ = cloneTypeExprs
	var clone Node = NewNode_1(this_.nodeGet_token())
	_ = clone
	Node_CopyBaseNodeFields(this_, clone, cloneTypeExprs)
	return clone
}

func (this_ Node) CloneNode_1(p0_ bool) Node {
	switch nodestore.Node(this_).KindFlags() {
	case 1:
		return this_.nodeSub1_CloneNode_1(p0_)
	case 2:
		return this_.nodeSub2_CloneNode_1(p0_)
	case 3:
		return this_.nodeSub3_CloneNode_1(p0_)
	case 4:
		return this_.nodeSub4_CloneNode_1(p0_)
	default:
		return this_.base_CloneNode_1(p0_)
	}
}

// Java static method: Node.copyBaseNodeFields(com.google.javascript.rhino.Node source, com.google.javascript.rhino.Node dest, boolean cloneTypeExprs)
func Node_CopyBaseNodeFields(source Node, dest Node, cloneTypeExprs bool) {
	clinit_Node()
	_ = source
	_ = dest
	_ = cloneTypeExprs
	{
		r_101 := dest
		r_101.nodeSet_linenoCharno(source.nodeGet_linenoCharno())
	}
	{
		r_102 := dest
		r_102.nodeSet_length(source.nodeGet_length())
	}
	{
		r_103 := dest
		r_103.nodeSet_jstypeOrColor(source.nodeGet_jstypeOrColor())
	}
	{
		r_104 := dest
		r_104.nodeSet_originalName(source.nodeGet_originalName())
	}
	{
		r_105 := dest
		r_105.nodeSet_sourceFile(source.nodeGet_sourceFile())
	}
	var info *JSDocInfo = source.nodeGet_jsDocInfo()
	_ = info
	{
		r_106 := dest
		r_106.nodeSet_jsDocInfo(func() *JSDocInfo {
			if cloneTypeExprs && (info != nil) {
				return info.Clone_1(true)
			}
			return info
		}())
	}
	{
		r_107 := dest
		r_107.nodeSet_typeBeforeCast(source.nodeGet_typeBeforeCast())
	}
	{
		r_108 := dest
		r_108.nodeSet_propFlags(source.nodeGet_propFlags())
	}
	{
		r_109 := dest
		r_109.nodeSet_changeTime(source.nodeGet_changeTime())
	}
	{
		r_110 := dest
		r_110.nodeSet_propListHead(source.nodeGet_propListHead())
	}
}

// Java method: Node.cloneTree()
func (this_ Node) CloneTree_0() Node {
	return this_.CloneTree_1(false)
}

// Java method: Node.cloneTree(boolean cloneTypeExprs)
func (this_ Node) CloneTree_1(cloneTypeExprs bool) Node {
	_ = cloneTypeExprs
	var result Node = this_.CloneNode_1(cloneTypeExprs)
	_ = result
	var firstChild Node = Node(0)
	_ = firstChild
	var lastChild Node = Node(0)
	_ = lastChild
	if this_.HasChildren() {
		{
			var n2 Node = this_.GetFirstChild()
			_ = n2
			for ; n2 != 0; n2 = n2.nodeGet_next() {
				var n2clone Node = n2.CloneTree_1(cloneTypeExprs)
				_ = n2clone
				{
					r_111 := n2clone
					r_111.nodeSet_parent(result)
				}
				if firstChild == 0 {
					firstChild = n2clone
					lastChild = firstChild
				} else {
					{
						r_112 := lastChild
						r_112.nodeSet_next(n2clone)
					}
					{
						r_113 := n2clone
						r_113.nodeSet_previous(lastChild)
					}
					lastChild = n2clone
				}
			}
		}
		{
			r_114 := firstChild
			r_114.nodeSet_previous(lastChild)
		}
		{
			r_115 := lastChild
			r_115.nodeSet_next(Node(0))
		}
		{
			r_116 := result
			r_116.nodeSet_first(firstChild)
		}
	}
	var shadow Node = this_.GetClosureUnawareShadow()
	_ = shadow
	if shadow != 0 {
		result.SetClosureUnawareShadow(shadow.CloneTree_1(cloneTypeExprs))
	}
	return result
}

// Java method: Node.srcref(com.google.javascript.rhino.Node other)
func (this_ Node) Srcref(other Node) Node {
	_ = other
	this_.SetStaticSourceFileFrom(other)
	this_.nodeSet_originalName(other.nodeGet_originalName())
	this_.nodeSet_linenoCharno(other.nodeGet_linenoCharno())
	this_.nodeSet_length(other.nodeGet_length())
	return this_
}

// Java method: Node.srcrefTree(com.google.javascript.rhino.Node other)
func (this_ Node) SrcrefTree(other Node) Node {
	_ = other
	this_.Srcref(other)
	{
		var child Node = this_.nodeGet_first()
		_ = child
		for ; child != 0; child = child.nodeGet_next() {
			child.SrcrefTree(other)
		}
	}
	return this_
}

// Java method: Node.srcrefIfMissing(com.google.javascript.rhino.Node other)
func (this_ Node) SrcrefIfMissing(other Node) Node {
	_ = other
	if j2gort.IfaceEquals[StaticSourceFile](this_.GetStaticSourceFile(), nil) {
		this_.SetStaticSourceFileFrom(other)
		this_.nodeSet_linenoCharno(other.nodeGet_linenoCharno())
		this_.nodeSet_length(other.nodeGet_length())
	}
	if this_.nodeGet_originalName() == j2gort.NullString {
		this_.nodeSet_originalName(other.nodeGet_originalName())
	}
	return this_
}

// Java method: Node.srcrefTreeIfMissing(com.google.javascript.rhino.Node other)
func (this_ Node) SrcrefTreeIfMissing(other Node) Node {
	_ = other
	this_.SrcrefIfMissing(other)
	{
		var child Node = this_.nodeGet_first()
		_ = child
		for ; child != 0; child = child.nodeGet_next() {
			child.SrcrefTreeIfMissing(other)
		}
	}
	return this_
}

// Java method: Node.getJSType()
func (this_ Node) GetJSType() JSType {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if v := nodestore.Node(this_).JSTypeObj(); v != nil {
		return v.(JSType)
	}
	return nil
}

// Java method: Node.getJSTypeRequired()
func (this_ Node) GetJSTypeRequired() JSType {
	return j2gort.CheckNotNullIface(this_.GetJSType())
}

// Java method: Node.setJSType(JSType)
func (this_ Node) SetJSType(x JSType) Node {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	var v any
	if x != nil && !x.IsJavaNil_() {
		v = x
	}
	if !nodestore.Node(this_).SetJSTypeObj(v) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(this_))
	}
	return this_
}

// Java method: Node.getColor()
func (this_ Node) GetColor() *pkg_colors.Color {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if v := nodestore.Node(this_).ColorObj(); v != nil {
		return v.(*pkg_colors.Color)
	}
	return nil
}

// Java method: Node.setColor(Color)
func (this_ Node) SetColor(x *pkg_colors.Color) Node {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	var v any
	if x != nil && !x.IsJavaNil_() {
		v = x
	}
	if !nodestore.Node(this_).SetColorObj(v) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(this_))
	}
	return this_
}

// Java method: Node.copyTypeFrom(com.google.javascript.rhino.Node other)
func (this_ Node) CopyTypeFrom(other Node) Node {
	_ = other
	this_.nodeSet_jstypeOrColor(other.nodeGet_jstypeOrColor())
	return this_
}

// Java method: Node.getJSDocInfo()
func (this_ Node) GetJSDocInfo() *JSDocInfo {
	return this_.nodeGet_jsDocInfo()
}

// Java method: Node.setJSDocInfo(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ Node) SetJSDocInfo(info *JSDocInfo) Node {
	if info == nil {
		nodestore.Node(this_).ClearJSDoc()
		return this_
	}

	_ = info
	nodestore.Node(this_).SetJSDocObj(info)

	return this_
}

// Java method: Node.setChangeTime(int time)
func (this_ Node) SetChangeTime(time int32) {
	_ = time
	this_.nodeSet_changeTime(time)
}

// Java method: Node.getChangeTime()
func (this_ Node) GetChangeTime() int32 {
	return this_.nodeGet_changeTime()
}

// Java method: Node.setDeleted(boolean deleted)
func (this_ Node) SetDeleted(deleted bool) {
	_ = deleted
	this_.PutBooleanProp(StaticGet_Node_Prop_DELETED(), deleted)
}

// Java method: Node.isDeleted()
func (this_ Node) IsDeleted() bool {
	return ((this_.nodeGet_propFlags() & int32(524288)) != int32(0))
}

// Java method: Node.setTypedefTypeProp(com.google.javascript.rhino.jstype.JSType type)
func (this_ Node) SetTypedefTypeProp(type_ JSType) {
	_ = type_
	this_.PutProp(StaticGet_Node_Prop_TYPEDEF_TYPE(), type_)
}

// Java method: Node.getTypedefTypeProp()
func (this_ Node) GetTypedefTypeProp() JSType {
	return j2gort.AsIface[JSType](this_.GetProp(StaticGet_Node_Prop_TYPEDEF_TYPE()))
}

// Java method: Node.setUnusedParameter(boolean unused)
func (this_ Node) SetUnusedParameter(unused bool) {
	_ = unused
	this_.PutBooleanProp(StaticGet_Node_Prop_IS_UNUSED_PARAMETER(), unused)
}

// Java method: Node.isUnusedParameter()
func (this_ Node) IsUnusedParameter() bool {
	return ((this_.nodeGet_propFlags() & int32(2097152)) != int32(0))
}

// Java method: Node.setShorthandProperty(boolean shorthand)
func (this_ Node) SetShorthandProperty(shorthand bool) {
	_ = shorthand
	this_.PutBooleanProp(StaticGet_Node_Prop_IS_SHORTHAND_PROPERTY(), shorthand)
}

// Java method: Node.isShorthandProperty()
func (this_ Node) IsShorthandProperty() bool {
	return ((this_.nodeGet_propFlags() & int32(1024)) != int32(0))
}

// Java method: Node.isOptionalEs6Typed()
func (this_ Node) IsOptionalEs6Typed() bool {
	return this_.GetBooleanProp(StaticGet_Node_Prop_OPT_ES6_TYPED())
}

// Java method: Node.setIsSyntheticBlock(boolean val)
func (this_ Node) SetIsSyntheticBlock(val bool) {
	_ = val
	if !(this_.RawToken_() == uint8(90)) {
		j2gort.CheckState(false)
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_SYNTHETIC(), val)
}

// Java method: Node.isSyntheticBlock()
func (this_ Node) IsSyntheticBlock() bool {
	return ((this_.nodeGet_propFlags() & int32(512)) != int32(0))
}

// Java method: Node.setIsSynthesizedUnfulfilledNameDeclaration(boolean val)
func (this_ Node) SetIsSynthesizedUnfulfilledNameDeclaration(val bool) {
	_ = val
	if !(((this_.RawToken_() == uint8(84)) && this_.HasOneChild()) && this_.GetFirstChild().IsName()) {
		j2gort.CheckStateFmt(false, "Expected all synthetic unfulfilled declarations to be `var <name>`, found %s", j2gort.BoxNode(this_))
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION(), val)
}

// Java method: Node.isSynthesizedUnfulfilledNameDeclaration()
func (this_ Node) IsSynthesizedUnfulfilledNameDeclaration() bool {
	return this_.GetBooleanProp(StaticGet_Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION())
}

// Java method: Node.setUseStrict(boolean x)
func (this_ Node) SetUseStrict(x bool) {
	_ = x
	this_.PutBooleanProp(StaticGet_Node_Prop_USE_STRICT(), x)
}

// Java method: Node.isUseStrict()
func (this_ Node) IsUseStrict() bool {
	return ((this_.nodeGet_propFlags() & int32(4096)) != int32(0))
}

// Java method: Node.setIsAddedBlock(boolean val)
func (this_ Node) SetIsAddedBlock(val bool) {
	_ = val
	this_.PutBooleanProp(StaticGet_Node_Prop_ADDED_BLOCK(), val)
}

// Java method: Node.isAddedBlock()
func (this_ Node) IsAddedBlock() bool {
	return ((this_.nodeGet_propFlags() & int32(16)) != int32(0))
}

// Java method: Node.setStaticMember(boolean isStatic)
func (this_ Node) SetStaticMember(isStatic bool) {
	_ = isStatic
	this_.PutBooleanProp(StaticGet_Node_Prop_STATIC_MEMBER(), isStatic)
}

// Java method: Node.isStaticMember()
func (this_ Node) IsStaticMember() bool {
	return ((this_.nodeGet_propFlags() & int32(8192)) != int32(0))
}

// Java method: Node.setIsGeneratorFunction(boolean isGenerator)
func (this_ Node) SetIsGeneratorFunction(isGenerator bool) {
	_ = isGenerator
	this_.PutBooleanProp(StaticGet_Node_Prop_GENERATOR_FN(), isGenerator)
}

// Java method: Node.isGeneratorFunction()
func (this_ Node) IsGeneratorFunction() bool {
	return ((this_.nodeGet_propFlags() & int32(16384)) != int32(0))
}

// Java method: Node.setGeneratorMarker(boolean isGeneratorMarker)
func (this_ Node) SetGeneratorMarker(isGeneratorMarker bool) {
	_ = isGeneratorMarker
	this_.PutBooleanProp(StaticGet_Node_Prop_IS_GENERATOR_MARKER(), isGeneratorMarker)
}

// Java method: Node.isGeneratorMarker()
func (this_ Node) IsGeneratorMarker() bool {
	return this_.GetBooleanProp(StaticGet_Node_Prop_IS_GENERATOR_MARKER())
}

// Java method: Node.setGeneratorSafe(boolean isGeneratorSafe)
func (this_ Node) SetGeneratorSafe(isGeneratorSafe bool) {
	_ = isGeneratorSafe
	this_.PutBooleanProp(StaticGet_Node_Prop_IS_GENERATOR_SAFE(), isGeneratorSafe)
}

// Java method: Node.isGeneratorSafe()
func (this_ Node) IsGeneratorSafe() bool {
	return this_.GetBooleanProp(StaticGet_Node_Prop_IS_GENERATOR_SAFE())
}

// Java method: Node.setIsOptionalChainStart(boolean isOptionalChainStart)
func (this_ Node) SetIsOptionalChainStart(isOptionalChainStart bool) {
	_ = isOptionalChainStart
	if !(((!(isOptionalChainStart) || this_.IsOptChainGetElem()) || this_.IsOptChainGetProp()) || this_.IsOptChainCall()) {
		j2gort.CheckStateStr(false, "cannot make a non-optional node the start of an optional chain.")
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_START_OF_OPT_CHAIN(), isOptionalChainStart)
}

// Java method: Node.isOptionalChainStart()
func (this_ Node) IsOptionalChainStart() bool {
	return ((this_.nodeGet_propFlags() & int32(8388608)) != int32(0))
}

// Java method: Node.setIsArrowFunction(boolean isArrow)
func (this_ Node) SetIsArrowFunction(isArrow bool) {
	_ = isArrow
	if !(this_.IsFunction()) {
		j2gort.CheckState(false)
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_ARROW_FN(), isArrow)
}

// Java method: Node.isArrowFunction()
func (this_ Node) IsArrowFunction() bool {
	return (((this_.nodeGet_propFlags() & int32(128)) != int32(0)) && this_.IsFunction())
}

// Java method: Node.setIsAsyncFunction(boolean isAsync)
func (this_ Node) SetIsAsyncFunction(isAsync bool) {
	_ = isAsync
	if !(this_.IsFunction()) {
		j2gort.CheckState(false)
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_ASYNC_FN(), isAsync)
}

// Java method: Node.isAsyncFunction()
func (this_ Node) IsAsyncFunction() bool {
	return (((this_.nodeGet_propFlags() & int32(32768)) != int32(0)) && this_.IsFunction())
}

// Java method: Node.isAsyncGeneratorFunction()
func (this_ Node) IsAsyncGeneratorFunction() bool {
	return (this_.IsAsyncFunction() && this_.IsGeneratorFunction())
}

// Java method: Node.setYieldAll(boolean isGenerator)
func (this_ Node) SetYieldAll(isGenerator bool) {
	_ = isGenerator
	this_.PutBooleanProp(StaticGet_Node_Prop_YIELD_ALL(), isGenerator)
}

// Java method: Node.isYieldAll()
func (this_ Node) IsYieldAll() bool {
	return this_.GetBooleanProp(StaticGet_Node_Prop_YIELD_ALL())
}

// Java method: Node.setTrailingComma(boolean hasTrailingComma)
func (this_ Node) SetTrailingComma(hasTrailingComma bool) {
	_ = hasTrailingComma
	this_.PutBooleanProp(StaticGet_Node_Prop_TRAILING_COMMA(), hasTrailingComma)
}

// Java method: Node.hasTrailingComma()
func (this_ Node) HasTrailingComma() bool {
	return ((this_.nodeGet_propFlags() & int32(32)) != int32(0))
}

// Java method: Node.setSideEffectFlags(int flags)
func (this_ Node) SetSideEffectFlags_1_int32(flags int32) {
	_ = flags
	if !(((this_.IsCall() || this_.IsOptChainCall()) || this_.IsNew()) || this_.IsTaggedTemplateLit()) {
		j2gort.CheckStateFmt(false, "Side-effect flags can only be set on invocation nodes; got %s", j2gort.BoxNode(this_))
	}
	this_.PutIntProp(StaticGet_Node_Prop_SIDE_EFFECT_FLAGS(), (^(flags) & int32(15)))
}

// Java method: Node.setSideEffectFlags(com.google.javascript.rhino.Node.SideEffectFlags flags)
func (this_ Node) SetSideEffectFlags_1_Node_SideEffectFlags(flags *Node_SideEffectFlags) {
	_ = flags
	this_.SetSideEffectFlags_1_int32(flags.ValueOf())
}

// Java method: Node.getSideEffectFlags()
func (this_ Node) GetSideEffectFlags() int32 {
	var v int32 = ((javalang.Ushr32(this_.nodeGet_propFlags(), int(int32(28)))) & int32(15))
	_ = v
	if (v == int32(0)) && (this_.nodeGet_propListHead() != nil) {
		v = this_.priv_Node_getColdIntProp_8426(StaticGet_Node_Prop_SIDE_EFFECT_FLAGS())
	}
	return (^(v) & int32(15))
}

// Java method: Node.isOnlyModifiesThisCall()
func (this_ Node) IsOnlyModifiesThisCall() bool {
	var sideEffectsBesidesMutatesThis int32 = (this_.GetSideEffectFlags() & ^(int32(2)))
	_ = sideEffectsBesidesMutatesThis
	return (sideEffectsBesidesMutatesThis == int32(0))
}

// Java method: Node.isOnlyModifiesArgumentsCall()
func (this_ Node) IsOnlyModifiesArgumentsCall() bool {
	var sideEffectsBesidesMutatesArguments int32 = (this_.GetSideEffectFlags() & ^(int32(4)))
	_ = sideEffectsBesidesMutatesArguments
	return (sideEffectsBesidesMutatesArguments == int32(0))
}

// Java method: Node.isNoSideEffectsCall()
func (this_ Node) IsNoSideEffectsCall() bool {
	return (this_.GetSideEffectFlags() == int32(0))
}

// Java method: Node.mayMutateArguments()
func (this_ Node) MayMutateArguments() bool {
	return ((this_.GetSideEffectFlags() & int32(4)) == int32(4))
}

// Java method: Node.mayMutateGlobalStateOrThrow()
func (this_ Node) MayMutateGlobalStateOrThrow() bool {
	return ((this_.GetSideEffectFlags() & (int32(1) | int32(8))) != int32(0))
}

// Java static method: Node.anyBitSet(int value, int mask)
func Node_AnyBitSet(value int32, mask int32) bool {
	clinit_Node()
	_ = value
	_ = mask
	return ((value & mask) != int32(0))
}

// Java method: Node.getConstantVarFlags()
func (this_ Node) priv_Node_getConstantVarFlags_8605() int32 {
	var v int32 = ((javalang.Ushr32(this_.nodeGet_propFlags(), int(int32(26)))) & int32(3))
	_ = v
	if (v == int32(0)) && (this_.nodeGet_propListHead() != nil) {
		v = this_.priv_Node_getColdIntProp_8426(StaticGet_Node_Prop_CONSTANT_VAR_FLAGS())
	}
	return v
}

// Java method: Node.setConstantVarFlag(int flag, boolean value)
func (this_ Node) priv_Node_setConstantVarFlag_8606(flag int32, value bool) {
	_ = flag
	_ = value
	var flags int32 = this_.priv_Node_getConstantVarFlags_8605()
	_ = flags
	if value {
		flags = (flags | flag)
	} else {
		flags = (flags & ^(flag))
	}
	this_.PutIntProp(StaticGet_Node_Prop_CONSTANT_VAR_FLAGS(), flags)
}

// Java method: Node.isDeclaredConstantVar()
func (this_ Node) IsDeclaredConstantVar() bool {
	if !(this_.IsName() || this_.IsImportStar()) {
		j2gort.CheckStateFmt(false, "Should only be called on name or import * nodes. Found %s", j2gort.BoxNode(this_))
	}
	return ((this_.priv_Node_getConstantVarFlags_8605() & int32(1)) != int32(0))
}

// Java method: Node.setDeclaredConstantVar(boolean value)
func (this_ Node) SetDeclaredConstantVar(value bool) {
	_ = value
	if !(this_.IsName() || this_.IsImportStar()) {
		j2gort.CheckStateFmt(false, "Should only be called on name or import * nodes. Found %s", j2gort.BoxNode(this_))
	}
	this_.priv_Node_setConstantVarFlag_8606(int32(1), value)
}

// Java method: Node.isInferredConstantVar()
func (this_ Node) IsInferredConstantVar() bool {
	if !(this_.IsName() || this_.IsImportStar()) {
		j2gort.CheckStateFmt(false, "Should only be called on name or import * nodes. Found %s", j2gort.BoxNode(this_))
	}
	return ((this_.priv_Node_getConstantVarFlags_8605() & int32(2)) != int32(0))
}

// Java method: Node.setInferredConstantVar(boolean value)
func (this_ Node) SetInferredConstantVar(value bool) {
	_ = value
	if !(this_.IsName() || this_.IsImportStar()) {
		j2gort.CheckStateFmt(false, "Should only be called on name or import * nodes. Found %s", j2gort.BoxNode(this_))
	}
	this_.priv_Node_setConstantVarFlag_8606(int32(2), value)
}

// Java method: Node.isQuotedStringKey()
func (this_ Node) IsQuotedStringKey() bool {
	return ((this_ != 0 && nodestore.Node(this_).KindFlags() == 3) && ((this_.nodeGet_propFlags() & int32(2)) != int32(0)))
}

// Java method: Node.setQuotedStringKey()
func (this_ Node) SetQuotedStringKey() {
	if !(this_ != 0 && nodestore.Node(this_).KindFlags() == 3) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(this_))
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_QUOTED(), true)
}

// Java method: Node.isPrivateIdentifier()
func (this_ Node) IsPrivateIdentifier() bool {
	return ((this_ != 0 && nodestore.Node(this_).KindFlags() == 3) && this_.GetBooleanProp(StaticGet_Node_Prop_PRIVATE_IDENTIFIER()))
}

// Java method: Node.setPrivateIdentifier()
func (this_ Node) SetPrivateIdentifier() {
	if !(this_ != 0 && nodestore.Node(this_).KindFlags() == 3) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(this_))
	}
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(this_.GetString()), "#")) {
		j2gort.CheckState(false)
	}
	this_.PutBooleanProp(StaticGet_Node_Prop_PRIVATE_IDENTIFIER(), true)
}

// Java method: Node.isAdd()
func (this_ Node) IsAdd() bool {
	return (this_.RawToken_() == uint8(13))
}

// Java method: Node.isSub()
func (this_ Node) IsSub() bool {
	return (this_.RawToken_() == uint8(14))
}

// Java method: Node.isAnd()
func (this_ Node) IsAnd() bool {
	return (this_.RawToken_() == uint8(69))
}

// Java method: Node.isAssignAnd()
func (this_ Node) IsAssignAnd() bool {
	return (this_.RawToken_() == uint8(65))
}

// Java method: Node.isArrayLit()
func (this_ Node) IsArrayLit() bool {
	return (this_.RawToken_() == uint8(46))
}

// Java method: Node.isArrayPattern()
func (this_ Node) IsArrayPattern() bool {
	return (this_.RawToken_() == uint8(102))
}

// Java method: Node.isAssign()
func (this_ Node) IsAssign() bool {
	return (this_.RawToken_() == uint8(51))
}

// Java method: Node.isAssignAdd()
func (this_ Node) IsAssignAdd() bool {
	return (this_.RawToken_() == uint8(58))
}

// Java method: Node.isNormalBlock()
func (this_ Node) IsNormalBlock() bool {
	return this_.IsBlock()
}

// Java method: Node.isBlock()
func (this_ Node) IsBlock() bool {
	return (this_.RawToken_() == uint8(90))
}

// Java method: Node.isRoot()
func (this_ Node) IsRoot() bool {
	return (this_.RawToken_() == uint8(89))
}

// Java method: Node.isAwait()
func (this_ Node) IsAwait() bool {
	return (this_.RawToken_() == uint8(115))
}

// Java method: Node.isBigInt()
func (this_ Node) IsBigInt() bool {
	return (this_.RawToken_() == uint8(34))
}

// Java method: Node.isBitNot()
func (this_ Node) IsBitNot() bool {
	return (this_.RawToken_() == uint8(20))
}

// Java method: Node.isBreak()
func (this_ Node) IsBreak() bool {
	return (this_.RawToken_() == uint8(82))
}

// Java method: Node.isCall()
func (this_ Node) IsCall() bool {
	return (this_.RawToken_() == uint8(28))
}

// Java method: Node.isCase()
func (this_ Node) IsCase() bool {
	return (this_.RawToken_() == uint8(76))
}

// Java method: Node.isCast()
func (this_ Node) IsCast() bool {
	return (this_.RawToken_() == uint8(101))
}

// Java method: Node.isCatch()
func (this_ Node) IsCatch() bool {
	return (this_.RawToken_() == uint8(86))
}

// Java method: Node.isClass()
func (this_ Node) IsClass() bool {
	return (this_.RawToken_() == uint8(105))
}

// Java method: Node.isClassMembers()
func (this_ Node) IsClassMembers() bool {
	return (this_.RawToken_() == uint8(106))
}

// Java method: Node.isComma()
func (this_ Node) IsComma() bool {
	return (this_.RawToken_() == uint8(50))
}

// Java method: Node.isComputedProp()
func (this_ Node) IsComputedProp() bool {
	return (this_.RawToken_() == uint8(129))
}

// Java method: Node.isContinue()
func (this_ Node) IsContinue() bool {
	return (this_.RawToken_() == uint8(83))
}

// Java method: Node.isConst()
func (this_ Node) IsConst() bool {
	return (this_.RawToken_() == uint8(97))
}

// Java method: Node.isDebugger()
func (this_ Node) IsDebugger() bool {
	return (this_.RawToken_() == uint8(98))
}

// Java method: Node.isDec()
func (this_ Node) IsDec() bool {
	return (this_.RawToken_() == uint8(72))
}

// Java method: Node.isDefaultCase()
func (this_ Node) IsDefaultCase() bool {
	return (this_.RawToken_() == uint8(77))
}

// Java method: Node.isDefaultValue()
func (this_ Node) IsDefaultValue() bool {
	return (this_.RawToken_() == uint8(134))
}

// Java method: Node.isDelProp()
func (this_ Node) IsDelProp() bool {
	return (this_.RawToken_() == uint8(24))
}

// Java method: Node.isDestructuringLhs()
func (this_ Node) IsDestructuringLhs() bool {
	return (this_.RawToken_() == uint8(104))
}

// Java method: Node.isDestructuringPattern()
func (this_ Node) IsDestructuringPattern() bool {
	return (this_.IsObjectPattern() || this_.IsArrayPattern())
}

// Java method: Node.isDo()
func (this_ Node) IsDo() bool {
	return (this_.RawToken_() == uint8(79))
}

// Java method: Node.isEmpty()
func (this_ Node) IsEmpty() bool {
	return (this_.RawToken_() == uint8(88))
}

// Java method: Node.isExponent()
func (this_ Node) IsExponent() bool {
	return (this_.RawToken_() == uint8(18))
}

// Java method: Node.isAssignExponent()
func (this_ Node) IsAssignExponent() bool {
	return (this_.RawToken_() == uint8(63))
}

// Java method: Node.isExport()
func (this_ Node) IsExport() bool {
	return (this_.RawToken_() == uint8(120))
}

// Java method: Node.isExportSpec()
func (this_ Node) IsExportSpec() bool {
	return (this_.RawToken_() == uint8(122))
}

// Java method: Node.isExportSpecs()
func (this_ Node) IsExportSpecs() bool {
	return (this_.RawToken_() == uint8(121))
}

// Java method: Node.isExprResult()
func (this_ Node) IsExprResult() bool {
	return (this_.RawToken_() == uint8(93))
}

// Java method: Node.isFalse()
func (this_ Node) IsFalse() bool {
	return (this_.RawToken_() == uint8(38))
}

// Java method: Node.isVanillaFor()
func (this_ Node) IsVanillaFor() bool {
	return (this_.RawToken_() == uint8(80))
}

// Java method: Node.isForIn()
func (this_ Node) IsForIn() bool {
	return (this_.RawToken_() == uint8(81))
}

// Java method: Node.isForOf()
func (this_ Node) IsForOf() bool {
	return (this_.RawToken_() == uint8(112))
}

// Java method: Node.isForAwaitOf()
func (this_ Node) IsForAwaitOf() bool {
	return (this_.RawToken_() == uint8(113))
}

// Java method: Node.isFunction()
func (this_ Node) IsFunction() bool {
	return (this_.RawToken_() == uint8(73))
}

// Java method: Node.isGetterDef()
func (this_ Node) IsGetterDef() bool {
	return (this_.RawToken_() == uint8(95))
}

// Java method: Node.isGetElem()
func (this_ Node) IsGetElem() bool {
	return (this_.RawToken_() == uint8(27))
}

// Java method: Node.isGetProp()
func (this_ Node) IsGetProp() bool {
	return (this_.RawToken_() == uint8(26))
}

// Java method: Node.isHook()
func (this_ Node) IsHook() bool {
	return (this_.RawToken_() == uint8(67))
}

// Java method: Node.isIf()
func (this_ Node) IsIf() bool {
	return (this_.RawToken_() == uint8(74))
}

// Java method: Node.isImport()
func (this_ Node) IsImport() bool {
	return (this_.RawToken_() == uint8(116))
}

// Java method: Node.isImportMeta()
func (this_ Node) IsImportMeta() bool {
	return (this_.RawToken_() == uint8(136))
}

// Java method: Node.isImportStar()
func (this_ Node) IsImportStar() bool {
	return (this_.RawToken_() == uint8(119))
}

// Java method: Node.isImportSpec()
func (this_ Node) IsImportSpec() bool {
	return (this_.RawToken_() == uint8(118))
}

// Java method: Node.isImportSpecs()
func (this_ Node) IsImportSpecs() bool {
	return (this_.RawToken_() == uint8(117))
}

// Java method: Node.isIn()
func (this_ Node) IsIn() bool {
	return (this_.RawToken_() == uint8(44))
}

// Java method: Node.isInc()
func (this_ Node) IsInc() bool {
	return (this_.RawToken_() == uint8(71))
}

// Java method: Node.isInstanceOf()
func (this_ Node) IsInstanceOf() bool {
	return (this_.RawToken_() == uint8(45))
}

// Java method: Node.isInterface()
func (this_ Node) IsInterface() bool {
	return (this_.RawToken_() == uint8(164))
}

// Java method: Node.isInterfaceMembers()
func (this_ Node) IsInterfaceMembers() bool {
	return (this_.RawToken_() == uint8(166))
}

// Java method: Node.isRecordType()
func (this_ Node) IsRecordType() bool {
	return (this_.RawToken_() == uint8(149))
}

// Java method: Node.isCallSignature()
func (this_ Node) IsCallSignature() bool {
	return (this_.RawToken_() == uint8(174))
}

// Java method: Node.isIndexSignature()
func (this_ Node) IsIndexSignature() bool {
	return (this_.RawToken_() == uint8(173))
}

// Java method: Node.isLabel()
func (this_ Node) IsLabel() bool {
	return (this_.RawToken_() == uint8(92))
}

// Java method: Node.isLabelName()
func (this_ Node) IsLabelName() bool {
	return (this_.RawToken_() == uint8(99))
}

// Java method: Node.isLet()
func (this_ Node) IsLet() bool {
	return (this_.RawToken_() == uint8(111))
}

// Java method: Node.isMemberFunctionDef()
func (this_ Node) IsMemberFunctionDef() bool {
	return (this_.RawToken_() == uint8(107))
}

// Java method: Node.isMemberVariableDef()
func (this_ Node) IsMemberVariableDef() bool {
	return (this_.RawToken_() == uint8(172))
}

// Java method: Node.isMemberFieldDef()
func (this_ Node) IsMemberFieldDef() bool {
	return (this_.RawToken_() == uint8(108))
}

// Java method: Node.isComputedFieldDef()
func (this_ Node) IsComputedFieldDef() bool {
	return (this_.RawToken_() == uint8(109))
}

// Java method: Node.isModuleBody()
func (this_ Node) IsModuleBody() bool {
	return (this_.RawToken_() == uint8(123))
}

// Java method: Node.isName()
func (this_ Node) IsName() bool {
	return (this_.RawToken_() == uint8(32))
}

// Java method: Node.isNE()
func (this_ Node) IsNE() bool {
	return (this_.RawToken_() == uint8(5))
}

// Java method: Node.isSHNE()
func (this_ Node) IsSHNE() bool {
	return (this_.RawToken_() == uint8(41))
}

// Java method: Node.isEQ()
func (this_ Node) IsEQ() bool {
	return (this_.RawToken_() == uint8(4))
}

// Java method: Node.isSHEQ()
func (this_ Node) IsSHEQ() bool {
	return (this_.RawToken_() == uint8(40))
}

// Java method: Node.isNeg()
func (this_ Node) IsNeg() bool {
	return (this_.RawToken_() == uint8(22))
}

// Java method: Node.isNew()
func (this_ Node) IsNew() bool {
	return (this_.RawToken_() == uint8(23))
}

// Java method: Node.isNot()
func (this_ Node) IsNot() bool {
	return (this_.RawToken_() == uint8(19))
}

// Java method: Node.isNull()
func (this_ Node) IsNull() bool {
	return (this_.RawToken_() == uint8(36))
}

// Java method: Node.isNullishCoalesce()
func (this_ Node) IsNullishCoalesce() bool {
	return (this_.RawToken_() == uint8(70))
}

// Java method: Node.isAssignNullishCoalesce()
func (this_ Node) IsAssignNullishCoalesce() bool {
	return (this_.RawToken_() == uint8(66))
}

// Java method: Node.isNumber()
func (this_ Node) IsNumber() bool {
	return (this_.RawToken_() == uint8(33))
}

// Java method: Node.isObjectLit()
func (this_ Node) IsObjectLit() bool {
	return (this_.RawToken_() == uint8(47))
}

// Java method: Node.isObjectPattern()
func (this_ Node) IsObjectPattern() bool {
	return (this_.RawToken_() == uint8(103))
}

// Java method: Node.isOptChainCall()
func (this_ Node) IsOptChainCall() bool {
	return (this_.RawToken_() == uint8(31))
}

// Java method: Node.isOptChainGetElem()
func (this_ Node) IsOptChainGetElem() bool {
	return (this_.RawToken_() == uint8(30))
}

// Java method: Node.isOptChainGetProp()
func (this_ Node) IsOptChainGetProp() bool {
	return (this_.RawToken_() == uint8(29))
}

// Java method: Node.isOr()
func (this_ Node) IsOr() bool {
	return (this_.RawToken_() == uint8(68))
}

// Java method: Node.isAssignOr()
func (this_ Node) IsAssignOr() bool {
	return (this_.RawToken_() == uint8(64))
}

// Java method: Node.isParamList()
func (this_ Node) IsParamList() bool {
	return (this_.RawToken_() == uint8(49))
}

// Java method: Node.isPos()
func (this_ Node) IsPos() bool {
	return (this_.RawToken_() == uint8(21))
}

// Java method: Node.isRegExp()
func (this_ Node) IsRegExp() bool {
	return (this_.RawToken_() == uint8(42))
}

// Java method: Node.isRest()
func (this_ Node) IsRest() bool {
	return ((this_.RawToken_() == uint8(125)) || (this_.RawToken_() == uint8(126)))
}

// Java method: Node.isObjectRest()
func (this_ Node) IsObjectRest() bool {
	return (this_.RawToken_() == uint8(126))
}

// Java method: Node.isReturn()
func (this_ Node) IsReturn() bool {
	return (this_.RawToken_() == uint8(0))
}

// Java method: Node.isScript()
func (this_ Node) IsScript() bool {
	return (this_.RawToken_() == uint8(94))
}

// Java method: Node.isSetterDef()
func (this_ Node) IsSetterDef() bool {
	return (this_.RawToken_() == uint8(96))
}

// Java method: Node.isSpread()
func (this_ Node) IsSpread() bool {
	return ((this_.RawToken_() == uint8(127)) || (this_.RawToken_() == uint8(128)))
}

// Java method: Node.isString()
func (this_ Node) IsString() bool {
	return (this_.RawToken_() == uint8(35))
}

// Java method: Node.isStringKey()
func (this_ Node) IsStringKey() bool {
	return (this_.RawToken_() == uint8(100))
}

// Java method: Node.isStringLit()
func (this_ Node) IsStringLit() bool {
	return (this_.RawToken_() == uint8(35))
}

// Java method: Node.isSuper()
func (this_ Node) IsSuper() bool {
	return (this_.RawToken_() == uint8(110))
}

// Java method: Node.isSwitch()
func (this_ Node) IsSwitch() bool {
	return (this_.RawToken_() == uint8(75))
}

// Java method: Node.isSwitchBody()
func (this_ Node) IsSwitchBody() bool {
	return (this_.RawToken_() == uint8(91))
}

// Java method: Node.isTaggedTemplateLit()
func (this_ Node) IsTaggedTemplateLit() bool {
	return (this_.RawToken_() == uint8(130))
}

// Java method: Node.isTemplateLit()
func (this_ Node) IsTemplateLit() bool {
	return (this_.RawToken_() == uint8(131))
}

// Java method: Node.isTemplateLitString()
func (this_ Node) IsTemplateLitString() bool {
	return (this_.RawToken_() == uint8(133))
}

// Java method: Node.isTemplateLitSub()
func (this_ Node) IsTemplateLitSub() bool {
	return (this_.RawToken_() == uint8(132))
}

// Java method: Node.isThis()
func (this_ Node) IsThis() bool {
	return (this_.RawToken_() == uint8(37))
}

// Java method: Node.isThrow()
func (this_ Node) IsThrow() bool {
	return (this_.RawToken_() == uint8(43))
}

// Java method: Node.isTrue()
func (this_ Node) IsTrue() bool {
	return (this_.RawToken_() == uint8(39))
}

// Java method: Node.isTry()
func (this_ Node) IsTry() bool {
	return (this_.RawToken_() == uint8(48))
}

// Java method: Node.isTypeOf()
func (this_ Node) IsTypeOf() bool {
	return (this_.RawToken_() == uint8(25))
}

// Java method: Node.isVar()
func (this_ Node) IsVar() bool {
	return (this_.RawToken_() == uint8(84))
}

// Java method: Node.isVoid()
func (this_ Node) IsVoid() bool {
	return (this_.RawToken_() == uint8(87))
}

// Java method: Node.isWhile()
func (this_ Node) IsWhile() bool {
	return (this_.RawToken_() == uint8(78))
}

// Java method: Node.isWith()
func (this_ Node) IsWith() bool {
	return (this_.RawToken_() == uint8(85))
}

// Java method: Node.isYield()
func (this_ Node) IsYield() bool {
	return (this_.RawToken_() == uint8(114))
}

// Java method: Node.isDeclare()
func (this_ Node) IsDeclare() bool {
	return (this_.RawToken_() == uint8(171))
}

func (this_ Node) Equals(other any) bool {
	o, ok := other.(Node)
	return ok && this_ == o
}
func (this_ Node) HashCode() int32    { return int32(this_) }
func (this_ Node) IsJavaNil_() bool   { return this_ == 0 }
func (this_ Node) RawNodeId_() uint32 { return uint32(this_) }
func (this_ Node) JavaClassInfo_() (string, string) {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	switch nodestore.Node(this_).KindFlags() {
	case 1:
		return "com.google.javascript.rhino.Node$NumberNode", "NumberNode"
	case 2:
		return "com.google.javascript.rhino.Node$BigIntNode", "BigIntNode"
	case 3:
		return "com.google.javascript.rhino.Node$StringNode", "StringNode"
	case 4:
		return "com.google.javascript.rhino.Node$TemplateLiteralSubstringNode", "TemplateLiteralSubstringNode"
	default:
		return "com.google.javascript.rhino.Node", "Node"
	}
}
func init() {
	j2gort.RegisterNodeBoxer(func(id_ uint32) any { return Node(id_) })
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node", "com.google.javascript.rhino.Node", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 })
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.NumberNode", "com.google.javascript.rhino.Node$NumberNode", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 && nodestore.Node(n).KindFlags() == 1 })
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.BigIntNode", "com.google.javascript.rhino.Node$BigIntNode", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 && nodestore.Node(n).KindFlags() == 2 })
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.StringNode", "com.google.javascript.rhino.Node$StringNode", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 && nodestore.Node(n).KindFlags() == 3 })
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.TemplateLiteralSubstringNode", "com.google.javascript.rhino.Node$TemplateLiteralSubstringNode", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 && nodestore.Node(n).KindFlags() == 4 })
}

// Java enum: com.google.javascript.rhino.Node.Prop
type Node_Prop uint8

const (
	enum_Node_Prop_IS_PARENTHESIZED                         Node_Prop = 1
	enum_Node_Prop_NON_JSDOC_COMMENT                        Node_Prop = 2
	enum_Node_Prop_TRAILING_NON_JSDOC_COMMENT               Node_Prop = 3
	enum_Node_Prop_JSDOC_INFO                               Node_Prop = 4
	enum_Node_Prop_INCRDECR                                 Node_Prop = 5
	enum_Node_Prop_QUOTED                                   Node_Prop = 6
	enum_Node_Prop_SYNTHETIC                                Node_Prop = 7
	enum_Node_Prop_ADDED_BLOCK                              Node_Prop = 8
	enum_Node_Prop_SIDE_EFFECT_FLAGS                        Node_Prop = 9
	enum_Node_Prop_IS_CONSTANT_NAME                         Node_Prop = 10
	enum_Node_Prop_IS_NAMESPACE                             Node_Prop = 11
	enum_Node_Prop_USE_STRICT                               Node_Prop = 12
	enum_Node_Prop_DIRECT_EVAL                              Node_Prop = 13
	enum_Node_Prop_FREE_CALL                                Node_Prop = 14
	enum_Node_Prop_SOURCE_FILE                              Node_Prop = 15
	enum_Node_Prop_INPUT_ID                                 Node_Prop = 16
	enum_Node_Prop_CHANGE_TIME                              Node_Prop = 17
	enum_Node_Prop_REFLECTED_OBJECT                         Node_Prop = 18
	enum_Node_Prop_STATIC_MEMBER                            Node_Prop = 19
	enum_Node_Prop_GENERATOR_FN                             Node_Prop = 20
	enum_Node_Prop_ARROW_FN                                 Node_Prop = 21
	enum_Node_Prop_ASYNC_FN                                 Node_Prop = 22
	enum_Node_Prop_YIELD_ALL                                Node_Prop = 23
	enum_Node_Prop_EXPORT_DEFAULT                           Node_Prop = 24
	enum_Node_Prop_EXPORT_ALL_FROM                          Node_Prop = 25
	enum_Node_Prop_CONSTANT_VAR_FLAGS                       Node_Prop = 26
	enum_Node_Prop_IS_GENERATOR_MARKER                      Node_Prop = 27
	enum_Node_Prop_IS_GENERATOR_SAFE                        Node_Prop = 28
	enum_Node_Prop_COMPUTED_PROP_METHOD                     Node_Prop = 29
	enum_Node_Prop_COMPUTED_PROP_GETTER                     Node_Prop = 30
	enum_Node_Prop_COMPUTED_PROP_SETTER                     Node_Prop = 31
	enum_Node_Prop_COMPUTED_PROP_VARIABLE                   Node_Prop = 32
	enum_Node_Prop_DECLARED_TYPE_EXPR                       Node_Prop = 33
	enum_Node_Prop_TYPE_BEFORE_CAST                         Node_Prop = 34
	enum_Node_Prop_COLOR_FROM_CAST                          Node_Prop = 35
	enum_Node_Prop_OPT_ES6_TYPED                            Node_Prop = 36
	enum_Node_Prop_GENERIC_TYPE                             Node_Prop = 37
	enum_Node_Prop_IMPLEMENTS                               Node_Prop = 38
	enum_Node_Prop_CONSTRUCT_SIGNATURE                      Node_Prop = 39
	enum_Node_Prop_ACCESS_MODIFIER                          Node_Prop = 40
	enum_Node_Prop_NON_INDEXABLE                            Node_Prop = 41
	enum_Node_Prop_GOOG_MODULE                              Node_Prop = 42
	enum_Node_Prop_FEATURE_SET                              Node_Prop = 43
	enum_Node_Prop_IS_TYPESCRIPT_ABSTRACT                   Node_Prop = 44
	enum_Node_Prop_DELETED                                  Node_Prop = 45
	enum_Node_Prop_MODULE_ALIAS                             Node_Prop = 46
	enum_Node_Prop_IS_UNUSED_PARAMETER                      Node_Prop = 47
	enum_Node_Prop_MODULE_EXPORT                            Node_Prop = 48
	enum_Node_Prop_IS_SHORTHAND_PROPERTY                    Node_Prop = 49
	enum_Node_Prop_ES6_MODULE                               Node_Prop = 50
	enum_Node_Prop_TYPEDEF_TYPE                             Node_Prop = 51
	enum_Node_Prop_START_OF_OPT_CHAIN                       Node_Prop = 52
	enum_Node_Prop_TRAILING_COMMA                           Node_Prop = 53
	enum_Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION Node_Prop = 54
	enum_Node_Prop_CLOSURE_UNAWARE_SHADOW                   Node_Prop = 55
	enum_Node_Prop_PRIVATE_IDENTIFIER                       Node_Prop = 56
	enum_Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE            Node_Prop = 57
)

var enumBaseTable_Node_Prop = [...]j2gort.EnumBase{
	{Name_: "IS_PARENTHESIZED", Ordinal_: 0},
	{Name_: "NON_JSDOC_COMMENT", Ordinal_: 1},
	{Name_: "TRAILING_NON_JSDOC_COMMENT", Ordinal_: 2},
	{Name_: "JSDOC_INFO", Ordinal_: 3},
	{Name_: "INCRDECR", Ordinal_: 4},
	{Name_: "QUOTED", Ordinal_: 5},
	{Name_: "SYNTHETIC", Ordinal_: 6},
	{Name_: "ADDED_BLOCK", Ordinal_: 7},
	{Name_: "SIDE_EFFECT_FLAGS", Ordinal_: 8},
	{Name_: "IS_CONSTANT_NAME", Ordinal_: 9},
	{Name_: "IS_NAMESPACE", Ordinal_: 10},
	{Name_: "USE_STRICT", Ordinal_: 11},
	{Name_: "DIRECT_EVAL", Ordinal_: 12},
	{Name_: "FREE_CALL", Ordinal_: 13},
	{Name_: "SOURCE_FILE", Ordinal_: 14},
	{Name_: "INPUT_ID", Ordinal_: 15},
	{Name_: "CHANGE_TIME", Ordinal_: 16},
	{Name_: "REFLECTED_OBJECT", Ordinal_: 17},
	{Name_: "STATIC_MEMBER", Ordinal_: 18},
	{Name_: "GENERATOR_FN", Ordinal_: 19},
	{Name_: "ARROW_FN", Ordinal_: 20},
	{Name_: "ASYNC_FN", Ordinal_: 21},
	{Name_: "YIELD_ALL", Ordinal_: 22},
	{Name_: "EXPORT_DEFAULT", Ordinal_: 23},
	{Name_: "EXPORT_ALL_FROM", Ordinal_: 24},
	{Name_: "CONSTANT_VAR_FLAGS", Ordinal_: 25},
	{Name_: "IS_GENERATOR_MARKER", Ordinal_: 26},
	{Name_: "IS_GENERATOR_SAFE", Ordinal_: 27},
	{Name_: "COMPUTED_PROP_METHOD", Ordinal_: 28},
	{Name_: "COMPUTED_PROP_GETTER", Ordinal_: 29},
	{Name_: "COMPUTED_PROP_SETTER", Ordinal_: 30},
	{Name_: "COMPUTED_PROP_VARIABLE", Ordinal_: 31},
	{Name_: "DECLARED_TYPE_EXPR", Ordinal_: 32},
	{Name_: "TYPE_BEFORE_CAST", Ordinal_: 33},
	{Name_: "COLOR_FROM_CAST", Ordinal_: 34},
	{Name_: "OPT_ES6_TYPED", Ordinal_: 35},
	{Name_: "GENERIC_TYPE", Ordinal_: 36},
	{Name_: "IMPLEMENTS", Ordinal_: 37},
	{Name_: "CONSTRUCT_SIGNATURE", Ordinal_: 38},
	{Name_: "ACCESS_MODIFIER", Ordinal_: 39},
	{Name_: "NON_INDEXABLE", Ordinal_: 40},
	{Name_: "GOOG_MODULE", Ordinal_: 41},
	{Name_: "FEATURE_SET", Ordinal_: 42},
	{Name_: "IS_TYPESCRIPT_ABSTRACT", Ordinal_: 43},
	{Name_: "DELETED", Ordinal_: 44},
	{Name_: "MODULE_ALIAS", Ordinal_: 45},
	{Name_: "IS_UNUSED_PARAMETER", Ordinal_: 46},
	{Name_: "MODULE_EXPORT", Ordinal_: 47},
	{Name_: "IS_SHORTHAND_PROPERTY", Ordinal_: 48},
	{Name_: "ES6_MODULE", Ordinal_: 49},
	{Name_: "TYPEDEF_TYPE", Ordinal_: 50},
	{Name_: "START_OF_OPT_CHAIN", Ordinal_: 51},
	{Name_: "TRAILING_COMMA", Ordinal_: 52},
	{Name_: "SYNTHESIZED_UNFULFILLED_NAME_DECLARATION", Ordinal_: 53},
	{Name_: "CLOSURE_UNAWARE_SHADOW", Ordinal_: 54},
	{Name_: "PRIVATE_IDENTIFIER", Ordinal_: 55},
	{Name_: "IS_IN_CLOSURE_UNAWARE_SUBTREE", Ordinal_: 56},
}
var enumValues_Node_Prop = []Node_Prop{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57}

func StaticGet_Node_Prop_IS_PARENTHESIZED() Node_Prop {
	return enum_Node_Prop_IS_PARENTHESIZED
}
func StaticGet_Node_Prop_NON_JSDOC_COMMENT() Node_Prop {
	return enum_Node_Prop_NON_JSDOC_COMMENT
}
func StaticGet_Node_Prop_TRAILING_NON_JSDOC_COMMENT() Node_Prop {
	return enum_Node_Prop_TRAILING_NON_JSDOC_COMMENT
}
func StaticGet_Node_Prop_JSDOC_INFO() Node_Prop {
	return enum_Node_Prop_JSDOC_INFO
}
func StaticGet_Node_Prop_INCRDECR() Node_Prop {
	return enum_Node_Prop_INCRDECR
}
func StaticGet_Node_Prop_QUOTED() Node_Prop {
	return enum_Node_Prop_QUOTED
}
func StaticGet_Node_Prop_SYNTHETIC() Node_Prop {
	return enum_Node_Prop_SYNTHETIC
}
func StaticGet_Node_Prop_ADDED_BLOCK() Node_Prop {
	return enum_Node_Prop_ADDED_BLOCK
}
func StaticGet_Node_Prop_SIDE_EFFECT_FLAGS() Node_Prop {
	return enum_Node_Prop_SIDE_EFFECT_FLAGS
}
func StaticGet_Node_Prop_IS_CONSTANT_NAME() Node_Prop {
	return enum_Node_Prop_IS_CONSTANT_NAME
}
func StaticGet_Node_Prop_IS_NAMESPACE() Node_Prop {
	return enum_Node_Prop_IS_NAMESPACE
}
func StaticGet_Node_Prop_USE_STRICT() Node_Prop {
	return enum_Node_Prop_USE_STRICT
}
func StaticGet_Node_Prop_DIRECT_EVAL() Node_Prop {
	return enum_Node_Prop_DIRECT_EVAL
}
func StaticGet_Node_Prop_FREE_CALL() Node_Prop {
	return enum_Node_Prop_FREE_CALL
}
func StaticGet_Node_Prop_SOURCE_FILE() Node_Prop {
	return enum_Node_Prop_SOURCE_FILE
}
func StaticGet_Node_Prop_INPUT_ID() Node_Prop {
	return enum_Node_Prop_INPUT_ID
}
func StaticGet_Node_Prop_CHANGE_TIME() Node_Prop {
	return enum_Node_Prop_CHANGE_TIME
}
func StaticGet_Node_Prop_REFLECTED_OBJECT() Node_Prop {
	return enum_Node_Prop_REFLECTED_OBJECT
}
func StaticGet_Node_Prop_STATIC_MEMBER() Node_Prop {
	return enum_Node_Prop_STATIC_MEMBER
}
func StaticGet_Node_Prop_GENERATOR_FN() Node_Prop {
	return enum_Node_Prop_GENERATOR_FN
}
func StaticGet_Node_Prop_ARROW_FN() Node_Prop {
	return enum_Node_Prop_ARROW_FN
}
func StaticGet_Node_Prop_ASYNC_FN() Node_Prop {
	return enum_Node_Prop_ASYNC_FN
}
func StaticGet_Node_Prop_YIELD_ALL() Node_Prop {
	return enum_Node_Prop_YIELD_ALL
}
func StaticGet_Node_Prop_EXPORT_DEFAULT() Node_Prop {
	return enum_Node_Prop_EXPORT_DEFAULT
}
func StaticGet_Node_Prop_EXPORT_ALL_FROM() Node_Prop {
	return enum_Node_Prop_EXPORT_ALL_FROM
}
func StaticGet_Node_Prop_CONSTANT_VAR_FLAGS() Node_Prop {
	return enum_Node_Prop_CONSTANT_VAR_FLAGS
}
func StaticGet_Node_Prop_IS_GENERATOR_MARKER() Node_Prop {
	return enum_Node_Prop_IS_GENERATOR_MARKER
}
func StaticGet_Node_Prop_IS_GENERATOR_SAFE() Node_Prop {
	return enum_Node_Prop_IS_GENERATOR_SAFE
}
func StaticGet_Node_Prop_COMPUTED_PROP_METHOD() Node_Prop {
	return enum_Node_Prop_COMPUTED_PROP_METHOD
}
func StaticGet_Node_Prop_COMPUTED_PROP_GETTER() Node_Prop {
	return enum_Node_Prop_COMPUTED_PROP_GETTER
}
func StaticGet_Node_Prop_COMPUTED_PROP_SETTER() Node_Prop {
	return enum_Node_Prop_COMPUTED_PROP_SETTER
}
func StaticGet_Node_Prop_COMPUTED_PROP_VARIABLE() Node_Prop {
	return enum_Node_Prop_COMPUTED_PROP_VARIABLE
}
func StaticGet_Node_Prop_DECLARED_TYPE_EXPR() Node_Prop {
	return enum_Node_Prop_DECLARED_TYPE_EXPR
}
func StaticGet_Node_Prop_TYPE_BEFORE_CAST() Node_Prop {
	return enum_Node_Prop_TYPE_BEFORE_CAST
}
func StaticGet_Node_Prop_COLOR_FROM_CAST() Node_Prop {
	return enum_Node_Prop_COLOR_FROM_CAST
}
func StaticGet_Node_Prop_OPT_ES6_TYPED() Node_Prop {
	return enum_Node_Prop_OPT_ES6_TYPED
}
func StaticGet_Node_Prop_GENERIC_TYPE() Node_Prop {
	return enum_Node_Prop_GENERIC_TYPE
}
func StaticGet_Node_Prop_IMPLEMENTS() Node_Prop {
	return enum_Node_Prop_IMPLEMENTS
}
func StaticGet_Node_Prop_CONSTRUCT_SIGNATURE() Node_Prop {
	return enum_Node_Prop_CONSTRUCT_SIGNATURE
}
func StaticGet_Node_Prop_ACCESS_MODIFIER() Node_Prop {
	return enum_Node_Prop_ACCESS_MODIFIER
}
func StaticGet_Node_Prop_NON_INDEXABLE() Node_Prop {
	return enum_Node_Prop_NON_INDEXABLE
}
func StaticGet_Node_Prop_GOOG_MODULE() Node_Prop {
	return enum_Node_Prop_GOOG_MODULE
}
func StaticGet_Node_Prop_FEATURE_SET() Node_Prop {
	return enum_Node_Prop_FEATURE_SET
}
func StaticGet_Node_Prop_IS_TYPESCRIPT_ABSTRACT() Node_Prop {
	return enum_Node_Prop_IS_TYPESCRIPT_ABSTRACT
}
func StaticGet_Node_Prop_DELETED() Node_Prop {
	return enum_Node_Prop_DELETED
}
func StaticGet_Node_Prop_MODULE_ALIAS() Node_Prop {
	return enum_Node_Prop_MODULE_ALIAS
}
func StaticGet_Node_Prop_IS_UNUSED_PARAMETER() Node_Prop {
	return enum_Node_Prop_IS_UNUSED_PARAMETER
}
func StaticGet_Node_Prop_MODULE_EXPORT() Node_Prop {
	return enum_Node_Prop_MODULE_EXPORT
}
func StaticGet_Node_Prop_IS_SHORTHAND_PROPERTY() Node_Prop {
	return enum_Node_Prop_IS_SHORTHAND_PROPERTY
}
func StaticGet_Node_Prop_ES6_MODULE() Node_Prop {
	return enum_Node_Prop_ES6_MODULE
}
func StaticGet_Node_Prop_TYPEDEF_TYPE() Node_Prop {
	return enum_Node_Prop_TYPEDEF_TYPE
}
func StaticGet_Node_Prop_START_OF_OPT_CHAIN() Node_Prop {
	return enum_Node_Prop_START_OF_OPT_CHAIN
}
func StaticGet_Node_Prop_TRAILING_COMMA() Node_Prop {
	return enum_Node_Prop_TRAILING_COMMA
}
func StaticGet_Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION() Node_Prop {
	return enum_Node_Prop_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
}
func StaticGet_Node_Prop_CLOSURE_UNAWARE_SHADOW() Node_Prop {
	return enum_Node_Prop_CLOSURE_UNAWARE_SHADOW
}
func StaticGet_Node_Prop_PRIVATE_IDENTIFIER() Node_Prop {
	return enum_Node_Prop_PRIVATE_IDENTIFIER
}
func StaticGet_Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE() Node_Prop {
	return enum_Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE
}

func Node_Prop_Values() []Node_Prop {
	return j2gort.CloneSlice(enumValues_Node_Prop)
}

func Node_Prop_ValueOf(name string) Node_Prop {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Node_Prop {
		if enumBaseTable_Node_Prop[i_].Name_ == name {
			return Node_Prop(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.Node.Prop." + name))
}

func (this_ Node_Prop) IsJavaNil_() bool                              { return this_ == 0 }
func (Node_Prop) HasIdentityEquals_()                                 {}
func (Node_Prop) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Node_Prop) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Node_Prop) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Node_Prop) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$Prop", "Prop"
}
func (this_ Node_Prop) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Node_Prop[this_-1]
}
func (this_ Node_Prop) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Node_Prop) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Node_Prop) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_Prop[this_-1].Name_
}
func (this_ Node_Prop) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Node_Prop[this_-1].Name_
}
func (this_ Node_Prop) Equals(other_ any) bool {
	o_, ok_ := other_.(Node_Prop)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Node_Prop) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Node_Prop[this_-1].HashCode()
}
func (this_ Node_Prop) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_Prop[this_-1].CompareTo(other_)
}
func (this_ Node_Prop) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Node_Prop[this_-1].Name_)
}
func init() {
	vals_ := []any{Node_Prop(1), Node_Prop(2), Node_Prop(3), Node_Prop(4), Node_Prop(5), Node_Prop(6), Node_Prop(7), Node_Prop(8), Node_Prop(9), Node_Prop(10), Node_Prop(11), Node_Prop(12), Node_Prop(13), Node_Prop(14), Node_Prop(15), Node_Prop(16), Node_Prop(17), Node_Prop(18), Node_Prop(19), Node_Prop(20), Node_Prop(21), Node_Prop(22), Node_Prop(23), Node_Prop(24), Node_Prop(25), Node_Prop(26), Node_Prop(27), Node_Prop(28), Node_Prop(29), Node_Prop(30), Node_Prop(31), Node_Prop(32), Node_Prop(33), Node_Prop(34), Node_Prop(35), Node_Prop(36), Node_Prop(37), Node_Prop(38), Node_Prop(39), Node_Prop(40), Node_Prop(41), Node_Prop(42), Node_Prop(43), Node_Prop(44), Node_Prop(45), Node_Prop(46), Node_Prop(47), Node_Prop(48), Node_Prop(49), Node_Prop(50), Node_Prop(51), Node_Prop(52), Node_Prop(53), Node_Prop(54), Node_Prop(55), Node_Prop(56), Node_Prop(57)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node.Prop", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node$Prop", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.Prop", "com.google.javascript.rhino.Node$Prop", func(v_ any) bool {
		return j2gort.InstanceOf[Node_Prop](v_)
	})
}

func StaticGet_Node_NumberNode_SerialVersionUID() int64 { return int64(1) }
func NewNode_NumberNode(number float64) Node {
	_ = number
	tok_ := tokenToByte(StaticGet_Token_NUMBER())
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.KindFlags = 1
	nv_ := number
	if !(nv_ == 0 && 1/nv_ > 0) {
		Node(raw_).SetDouble(nv_)
	}
	return Node(raw_)
}

// Java method: NumberNode.nodeSubclassKind()
func (this_ Node) nodeSub1_NodeSubclassKind() int32 {
	return int32(1)
}

// Java method: NumberNode.isEquivalentTo(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node.RecursionMode recursion, com.google.javascript.rhino.Node.TypeComparison typeConfig, com.google.javascript.rhino.Node.JsDocComparison jsDocConfig, com.google.javascript.rhino.Node.SideEffectComparison sideEffectConfig)
func (this_ Node) nodeSub1_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node Node, recursion Node_RecursionMode, typeConfig Node_TypeComparison, jsDocConfig Node_JsDocComparison, sideEffectConfig Node_SideEffectComparison) bool {
	_ = node
	_ = recursion
	_ = typeConfig
	_ = jsDocConfig
	_ = sideEffectConfig
	return (this_.base_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) && (this_.nodeGet_number() == (node).nodeGet_number()))
}

// Java method: NumberNode.cloneNode(boolean cloneTypeExprs)
func (this_ Node) nodeSub1_CloneNode_1(cloneTypeExprs bool) Node {
	_ = cloneTypeExprs
	var clone Node = NewNode_NumberNode(this_.nodeGet_number())
	_ = clone
	Node_CopyBaseNodeFields(this_, clone, cloneTypeExprs)
	return clone
}

func StaticGet_Node_BigIntNode_SerialVersionUID() int64 { return int64(1) }
func NewNode_BigIntNode(bigint *j2gort.Ext_java_math_BigInteger) Node {
	_ = bigint
	tok_ := tokenToByte(StaticGet_Token_BIGINT())
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.KindFlags = 2
	this_ := Node(raw_)
	this_.SetBigInt(bigint)
	return this_
}

// Java method: BigIntNode.nodeSubclassKind()
func (this_ Node) nodeSub2_NodeSubclassKind() int32 {
	return int32(2)
}

// Java method: BigIntNode.isEquivalentTo(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node.RecursionMode recursion, com.google.javascript.rhino.Node.TypeComparison typeConfig, com.google.javascript.rhino.Node.JsDocComparison jsDocConfig, com.google.javascript.rhino.Node.SideEffectComparison sideEffectConfig)
func (this_ Node) nodeSub2_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node Node, recursion Node_RecursionMode, typeConfig Node_TypeComparison, jsDocConfig Node_JsDocComparison, sideEffectConfig Node_SideEffectComparison) bool {
	_ = node
	_ = recursion
	_ = typeConfig
	_ = jsDocConfig
	_ = sideEffectConfig
	return (this_.base_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) && this_.GetBigInt().Equals(j2gort.PtrToAny(node.GetBigInt())))
}

// Java method: BigIntNode.cloneNode(boolean cloneTypeExprs)
func (this_ Node) nodeSub2_CloneNode_1(cloneTypeExprs bool) Node {
	_ = cloneTypeExprs
	var clone Node = NewNode_BigIntNode(this_.nodeGet_bigint())
	_ = clone
	Node_CopyBaseNodeFields(this_, clone, cloneTypeExprs)
	return clone
}

func StaticGet_Node_StringNode_SerialVersionUID() int64 { return int64(1) }
func NewNode_StringNode_1(token Token) Node {
	_ = token
	tok_ := tokenToByte(token)
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.KindFlags = 3
	return Node(raw_)
}

func NewNode_StringNode_2(token Token, str string) Node {
	_ = token
	_ = str
	tok_ := tokenToByte(token)
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.KindFlags = 3
	s_ := str
	if s_ == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	rec_.Payload = uint32(nodestore.Intern(s_))
	return Node(raw_)
}

func NewNode_StringNode_3(token Token, stringPool *RhinoStringPool_LazyInternedStringList, offset int32) Node {
	_ = token
	_ = stringPool
	_ = offset
	tok_ := tokenToByte(token)
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.KindFlags = 3
	rec_.Payload = uint32(stringPool.InternID_(offset))
	return Node(raw_)
}

// Java method: StringNode.nodeSubclassKind()
func (this_ Node) nodeSub3_NodeSubclassKind() int32 {
	return int32(3)
}

// Java method: StringNode.isEquivalentTo(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node.RecursionMode recursion, com.google.javascript.rhino.Node.TypeComparison typeConfig, com.google.javascript.rhino.Node.JsDocComparison jsDocConfig, com.google.javascript.rhino.Node.SideEffectComparison sideEffectConfig)
func (this_ Node) nodeSub3_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node Node, recursion Node_RecursionMode, typeConfig Node_TypeComparison, jsDocConfig Node_JsDocComparison, sideEffectConfig Node_SideEffectComparison) bool {
	_ = node
	_ = recursion
	_ = typeConfig
	_ = jsDocConfig
	_ = sideEffectConfig
	return (this_.base_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) && RhinoStringPool_UncheckedEquals(this_.nodeGet_str(), (node).nodeGet_str()))
}

// Java method: StringNode.cloneNode(boolean cloneTypeExprs)
func (this_ Node) nodeSub3_CloneNode_1(cloneTypeExprs bool) Node {
	_ = cloneTypeExprs
	var clone Node = NewNode_StringNode_1(this_.GetToken())
	_ = clone
	Node_CopyBaseNodeFields(this_, clone, cloneTypeExprs)
	{
		r_117 := clone
		r_117.nodeSet_str(this_.nodeGet_str())
	}
	return clone
}

func NewNode_TemplateLiteralSubstringNode_2(cooked string, raw string) Node {
	_ = cooked
	_ = raw
	tok_ := tokenToByte(StaticGet_Token_TEMPLATELIT_STRING())
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.KindFlags = 4
	this_ := Node(raw_)
	if cooked == j2gort.NullString {
		this_.nodeSet_cooked(j2gort.NullString)
	} else {
		this_.nodeSet_cooked(RhinoStringPool_AddOrGet(cooked))
	}
	this_.nodeSet_raw(RhinoStringPool_AddOrGet(raw))
	return this_
}

func NewNode_TemplateLiteralSubstringNode_3(stringPool *RhinoStringPool_LazyInternedStringList, cookedOffsetOrNegativeOne int32, rawOffset int32) Node {
	_ = stringPool
	_ = cookedOffsetOrNegativeOne
	_ = rawOffset
	tok_ := tokenToByte(StaticGet_Token_TEMPLATELIT_STRING())
	raw_, rec_ := nodestore.AllocNodeFastRec()
	rec_.Token = tok_
	rec_.LinenoCharno = -1
	rec_.KindFlags = 4
	this_ := Node(raw_)
	if cookedOffsetOrNegativeOne == int32(-1) {
		this_.nodeSet_cooked(j2gort.NullString)
	} else {
		this_.nodeSet_cooked(stringPool.Get(cookedOffsetOrNegativeOne))
	}
	this_.nodeSet_raw(stringPool.Get(rawOffset))
	return this_
}

// Java method: TemplateLiteralSubstringNode.nodeSubclassKind()
func (this_ Node) nodeSub4_NodeSubclassKind() int32 {
	return int32(4)
}

// Java method: TemplateLiteralSubstringNode.isEquivalentTo(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node.RecursionMode recursion, com.google.javascript.rhino.Node.TypeComparison typeConfig, com.google.javascript.rhino.Node.JsDocComparison jsDocConfig, com.google.javascript.rhino.Node.SideEffectComparison sideEffectConfig)
func (this_ Node) nodeSub4_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node Node, recursion Node_RecursionMode, typeConfig Node_TypeComparison, jsDocConfig Node_JsDocComparison, sideEffectConfig Node_SideEffectComparison) bool {
	_ = node
	_ = recursion
	_ = typeConfig
	_ = jsDocConfig
	_ = sideEffectConfig
	if !(this_.base_IsEquivalentTo_5_Node_Node_RecursionMode_Node_TypeComparison_Node_JsDocComparison_Node_SideEffectComparison(node, recursion, typeConfig, jsDocConfig, sideEffectConfig)) {
		return false
	}
	var castNode Node = node
	_ = castNode
	return (RhinoStringPool_UncheckedEquals(this_.nodeGet_raw(), castNode.nodeGet_raw()) && RhinoStringPool_UncheckedEquals(this_.nodeGet_cooked(), castNode.nodeGet_cooked()))
}

// Java method: TemplateLiteralSubstringNode.cloneNode(boolean cloneTypeExprs)
func (this_ Node) nodeSub4_CloneNode_1(cloneTypeExprs bool) Node {
	_ = cloneTypeExprs
	var clone Node = NewNode_TemplateLiteralSubstringNode_2(this_.nodeGet_cooked(), this_.nodeGet_raw())
	_ = clone
	Node_CopyBaseNodeFields(this_, clone, cloneTypeExprs)
	return clone
}

// Java class: com.google.javascript.rhino.Node.PropListItem
type Node_PropListItem struct {
	next        *Node_PropListItem `gson:"next"`
	objectValue any                `gson:"objectValue"`
	intValue    int32              `gson:"intValue"`
	idHash_     atomic.Int32
	propType    int8 `gson:"propType"`
}

func (this_ *Node_PropListItem) BindSelf_(s any) {
	_ = s
}
func (this_ *Node_PropListItem) IsJavaNil_() bool { return this_ == nil }
func (*Node_PropListItem) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$PropListItem", "PropListItem"
}
func (this_ *Node_PropListItem) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Node_PropListItem) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.PropListItem", "com.google.javascript.rhino.Node$PropListItem", func(v_ any) bool {
		return j2gort.InstanceOf[*Node_PropListItem](v_)
	})
}

func (this_ *Node_PropListItem) initFields_Node_PropListItem() {
	_ = this_
}

// Java constructor: PropListItem(byte propType, int intValue, java.lang.@org.jspecify.annotations.Nullable Object objectValue, com.google.javascript.rhino.Node.@org.jspecify.annotations.Nullable PropListItem next)
func NewNode_PropListItem(propType int8, intValue int32, objectValue any, next *Node_PropListItem) *Node_PropListItem {
	this_ := &Node_PropListItem{}
	InitNode_PropListItem(this_, propType, intValue, objectValue, next)
	return this_
}

func InitNode_PropListItem(this_ *Node_PropListItem, propType int8, intValue int32, objectValue any, next *Node_PropListItem) {
	_ = propType
	_ = intValue
	_ = objectValue
	_ = next
	this_.initFields_Node_PropListItem()
	this_.propType = propType
	this_.intValue = intValue
	this_.objectValue = objectValue
	this_.next = next
}

// Java static method: PropListItem.objectItem(byte propType, java.lang.Object objectValue, com.google.javascript.rhino.Node.@org.jspecify.annotations.Nullable PropListItem next)
func Node_PropListItem_ObjectItem(propType int8, objectValue any, next *Node_PropListItem) *Node_PropListItem {
	_ = propType
	_ = objectValue
	_ = next
	return NewNode_PropListItem(propType, int32(0), j2gort.CheckNotNull(objectValue), next)
}

// Java static method: PropListItem.intItem(byte propType, int intValue, com.google.javascript.rhino.Node.@org.jspecify.annotations.Nullable PropListItem next)
func Node_PropListItem_IntItem(propType int8, intValue int32, next *Node_PropListItem) *Node_PropListItem {
	_ = propType
	_ = intValue
	_ = next
	if !(intValue != int32(0)) {
		j2gort.CheckState(false)
	}
	if (intValue == int32(1)) && (next == nil) {
		return StaticGet_Node_STATIC_SINGLE_TRUE_ITEMS()[(int32(propType) & int32(255))]
	}
	return NewNode_PropListItem(propType, intValue, nil, next)
}

// Java method: PropListItem.isObjectItem()
func (this_ *Node_PropListItem) IsObjectItem() bool {
	return (this_.objectValue != nil)
}

// Java method: PropListItem.isIntItem()
func (this_ *Node_PropListItem) IsIntItem() bool {
	return (this_.objectValue == nil)
}

// Java method: PropListItem.getIntValue()
func (this_ *Node_PropListItem) GetIntValue() int32 {
	if this_.objectValue != nil {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
	}
	return this_.intValue
}

// Java method: PropListItem.getObjectValue()
func (this_ *Node_PropListItem) GetObjectValue() any {
	if this_.objectValue == nil {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
	}
	return this_.objectValue
}

// Java method: PropListItem.toString()
func (this_ *Node_PropListItem) String() string {
	if this_.objectValue != nil {
		return j2gort.StrVal(this_.objectValue)
	} else {
		return strconv.FormatInt(int64(this_.intValue), 10)
	}
}

// Java method: PropListItem.chain(com.google.javascript.rhino.Node.@org.jspecify.annotations.Nullable PropListItem next)
func (this_ *Node_PropListItem) Chain(next *Node_PropListItem) *Node_PropListItem {
	_ = next
	if ((next == nil) && (this_.objectValue == nil)) && (this_.intValue == int32(1)) {
		return StaticGet_Node_STATIC_SINGLE_TRUE_ITEMS()[(int32(this_.propType) & int32(255))]
	}
	return NewNode_PropListItem(this_.propType, this_.intValue, this_.objectValue, next)
}

// Java class: com.google.javascript.rhino.Node.SiblingNodeIterable
type Node_SiblingNodeIterable struct {
	start   Node `gson:"start"`
	idHash_ atomic.Int32
}

func (this_ *Node_SiblingNodeIterable) BindSelf_(s any) {
	_ = s
}
func (this_ *Node_SiblingNodeIterable) IsJavaNil_() bool { return this_ == nil }
func (*Node_SiblingNodeIterable) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$SiblingNodeIterable", "SiblingNodeIterable"
}
func (this_ *Node_SiblingNodeIterable) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Node_SiblingNodeIterable) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.SiblingNodeIterable", "com.google.javascript.rhino.Node$SiblingNodeIterable", func(v_ any) bool {
		return j2gort.InstanceOf[*Node_SiblingNodeIterable](v_)
	})
}
func (*Node_SiblingNodeIterable) Implements_j2gort_Ext_java_lang_Iterable() {}

func (this_ *Node_SiblingNodeIterable) initFields_Node_SiblingNodeIterable() {
	_ = this_
}

// Java constructor: SiblingNodeIterable(com.google.javascript.rhino.Node start)
func NewNode_SiblingNodeIterable(start Node) *Node_SiblingNodeIterable {
	this_ := &Node_SiblingNodeIterable{}
	InitNode_SiblingNodeIterable(this_, start)
	return this_
}

func InitNode_SiblingNodeIterable(this_ *Node_SiblingNodeIterable, start Node) {
	_ = start
	this_.initFields_Node_SiblingNodeIterable()
	this_.start = start
}

// Java method: SiblingNodeIterable.iterator()
func (this_ *Node_SiblingNodeIterable) Iterator() any {
	return NewNode_SiblingNodeIterator(this_.start)
}

func (this_ *Node_SiblingNodeIterable) ForEach(p0_ j2gort.Ext_java_util_function_Consumer) {
	j2gort.ExtDefault_Iterable_ForEach(this_, p0_)
}

// Java class: com.google.javascript.rhino.Node.SiblingNodeIterator
type Node_SiblingNodeIterator struct {
	current Node `gson:"current"`
	idHash_ atomic.Int32
}

func (this_ *Node_SiblingNodeIterator) BindSelf_(s any) {
	_ = s
}
func (this_ *Node_SiblingNodeIterator) IsJavaNil_() bool { return this_ == nil }
func (*Node_SiblingNodeIterator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$SiblingNodeIterator", "SiblingNodeIterator"
}
func (this_ *Node_SiblingNodeIterator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Node_SiblingNodeIterator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.SiblingNodeIterator", "com.google.javascript.rhino.Node$SiblingNodeIterator", func(v_ any) bool {
		return j2gort.InstanceOf[*Node_SiblingNodeIterator](v_)
	})
}
func (*Node_SiblingNodeIterator) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *Node_SiblingNodeIterator) initFields_Node_SiblingNodeIterator() {
	_ = this_
}

// Java constructor: SiblingNodeIterator(com.google.javascript.rhino.Node start)
func NewNode_SiblingNodeIterator(start Node) *Node_SiblingNodeIterator {
	this_ := &Node_SiblingNodeIterator{}
	InitNode_SiblingNodeIterator(this_, start)
	return this_
}

func InitNode_SiblingNodeIterator(this_ *Node_SiblingNodeIterator, start Node) {
	_ = start
	this_.initFields_Node_SiblingNodeIterator()
	this_.current = start
}

// Java method: SiblingNodeIterator.hasNext()
func (this_ *Node_SiblingNodeIterator) HasNext() bool {
	return (this_.current != 0)
}

// Java method: SiblingNodeIterator.next()
func (this_ *Node_SiblingNodeIterator) Next() any {
	if this_.current == 0 {
		panic(j2gort.AsIface[javalang.INoSuchElementException](j2gort.ConstructShimThrowable("java.util.NoSuchElementException")))
	}
	var n Node = this_.current
	_ = n
	this_.current = this_.current.GetNext()
	return j2gort.BoxNode(n)
}

// Java method: SiblingNodeIterator.remove()
func (this_ *Node_SiblingNodeIterator) Remove() {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java class: com.google.javascript.rhino.Node.AncestorIterable
type Node_AncestorIterable struct {
	cur     Node `gson:"cur"`
	idHash_ atomic.Int32
}

func (this_ *Node_AncestorIterable) BindSelf_(s any) {
	_ = s
}
func (this_ *Node_AncestorIterable) IsJavaNil_() bool { return this_ == nil }
func (*Node_AncestorIterable) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$AncestorIterable", "AncestorIterable"
}
func (this_ *Node_AncestorIterable) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Node_AncestorIterable) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.AncestorIterable", "com.google.javascript.rhino.Node$AncestorIterable", func(v_ any) bool {
		return j2gort.InstanceOf[*Node_AncestorIterable](v_)
	})
}
func (*Node_AncestorIterable) Implements_j2gort_Ext_java_lang_Iterable() {}

func (this_ *Node_AncestorIterable) initFields_Node_AncestorIterable() {
	_ = this_
}

// Java constructor: AncestorIterable(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node cur)
func NewNode_AncestorIterable(cur Node) *Node_AncestorIterable {
	this_ := &Node_AncestorIterable{}
	InitNode_AncestorIterable(this_, cur)
	return this_
}

func InitNode_AncestorIterable(this_ *Node_AncestorIterable, cur Node) {
	_ = cur
	this_.initFields_Node_AncestorIterable()
	this_.cur = cur
}

// Java method: AncestorIterable.iterator()
func (this_ *Node_AncestorIterable) Iterator() any {
	return NewNode__Anon_154(this_)
}

func (this_ *Node_AncestorIterable) ForEach(p0_ j2gort.Ext_java_util_function_Consumer) {
	j2gort.ExtDefault_Iterable_ForEach(this_, p0_)
}

// Java enum: com.google.javascript.rhino.Node.RecursionMode
type Node_RecursionMode uint8

const (
	enum_Node_RecursionMode_SHALLOW        Node_RecursionMode = 1
	enum_Node_RecursionMode_DEEP_NO_SHADOW Node_RecursionMode = 2
	enum_Node_RecursionMode_DEEP           Node_RecursionMode = 3
)

var enumBaseTable_Node_RecursionMode = [...]j2gort.EnumBase{
	{Name_: "SHALLOW", Ordinal_: 0},
	{Name_: "DEEP_NO_SHADOW", Ordinal_: 1},
	{Name_: "DEEP", Ordinal_: 2},
}
var enumValues_Node_RecursionMode = []Node_RecursionMode{1, 2, 3}

func StaticGet_Node_RecursionMode_SHALLOW() Node_RecursionMode {
	return enum_Node_RecursionMode_SHALLOW
}
func StaticGet_Node_RecursionMode_DEEP_NO_SHADOW() Node_RecursionMode {
	return enum_Node_RecursionMode_DEEP_NO_SHADOW
}
func StaticGet_Node_RecursionMode_DEEP() Node_RecursionMode {
	return enum_Node_RecursionMode_DEEP
}

func Node_RecursionMode_Values() []Node_RecursionMode {
	return j2gort.CloneSlice(enumValues_Node_RecursionMode)
}

func Node_RecursionMode_ValueOf(name string) Node_RecursionMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Node_RecursionMode {
		if enumBaseTable_Node_RecursionMode[i_].Name_ == name {
			return Node_RecursionMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.Node.RecursionMode." + name))
}

// Java method: RecursionMode.recurseShadowAst()
func (this_ Node_RecursionMode) RecurseShadowAst() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return (this_ == StaticGet_Node_RecursionMode_DEEP())
}

// Java method: RecursionMode.recurseChildren()
func (this_ Node_RecursionMode) RecurseChildren() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((this_ == StaticGet_Node_RecursionMode_DEEP()) || (this_ == StaticGet_Node_RecursionMode_DEEP_NO_SHADOW()))
}

func (this_ Node_RecursionMode) IsJavaNil_() bool                              { return this_ == 0 }
func (Node_RecursionMode) HasIdentityEquals_()                                 {}
func (Node_RecursionMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Node_RecursionMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Node_RecursionMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Node_RecursionMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$RecursionMode", "RecursionMode"
}
func (this_ Node_RecursionMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Node_RecursionMode[this_-1]
}
func (this_ Node_RecursionMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Node_RecursionMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Node_RecursionMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_RecursionMode[this_-1].Name_
}
func (this_ Node_RecursionMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Node_RecursionMode[this_-1].Name_
}
func (this_ Node_RecursionMode) Equals(other_ any) bool {
	o_, ok_ := other_.(Node_RecursionMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Node_RecursionMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Node_RecursionMode[this_-1].HashCode()
}
func (this_ Node_RecursionMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_RecursionMode[this_-1].CompareTo(other_)
}
func (this_ Node_RecursionMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Node_RecursionMode[this_-1].Name_)
}
func init() {
	vals_ := []any{Node_RecursionMode(1), Node_RecursionMode(2), Node_RecursionMode(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node.RecursionMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node$RecursionMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.RecursionMode", "com.google.javascript.rhino.Node$RecursionMode", func(v_ any) bool {
		return j2gort.InstanceOf[Node_RecursionMode](v_)
	})
}

// Java enum: com.google.javascript.rhino.Node.TypeComparison
type Node_TypeComparison uint8

const (
	enum_Node_TypeComparison_IGNORE  Node_TypeComparison = 1
	enum_Node_TypeComparison_COMPARE Node_TypeComparison = 2
)

var enumBaseTable_Node_TypeComparison = [...]j2gort.EnumBase{
	{Name_: "IGNORE", Ordinal_: 0},
	{Name_: "COMPARE", Ordinal_: 1},
}
var enumValues_Node_TypeComparison = []Node_TypeComparison{1, 2}

func StaticGet_Node_TypeComparison_IGNORE() Node_TypeComparison {
	return enum_Node_TypeComparison_IGNORE
}
func StaticGet_Node_TypeComparison_COMPARE() Node_TypeComparison {
	return enum_Node_TypeComparison_COMPARE
}

func Node_TypeComparison_Values() []Node_TypeComparison {
	return j2gort.CloneSlice(enumValues_Node_TypeComparison)
}

func Node_TypeComparison_ValueOf(name string) Node_TypeComparison {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Node_TypeComparison {
		if enumBaseTable_Node_TypeComparison[i_].Name_ == name {
			return Node_TypeComparison(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.Node.TypeComparison." + name))
}

func (this_ Node_TypeComparison) IsJavaNil_() bool                              { return this_ == 0 }
func (Node_TypeComparison) HasIdentityEquals_()                                 {}
func (Node_TypeComparison) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Node_TypeComparison) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Node_TypeComparison) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Node_TypeComparison) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$TypeComparison", "TypeComparison"
}
func (this_ Node_TypeComparison) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Node_TypeComparison[this_-1]
}
func (this_ Node_TypeComparison) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Node_TypeComparison) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Node_TypeComparison) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_TypeComparison[this_-1].Name_
}
func (this_ Node_TypeComparison) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Node_TypeComparison[this_-1].Name_
}
func (this_ Node_TypeComparison) Equals(other_ any) bool {
	o_, ok_ := other_.(Node_TypeComparison)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Node_TypeComparison) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Node_TypeComparison[this_-1].HashCode()
}
func (this_ Node_TypeComparison) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_TypeComparison[this_-1].CompareTo(other_)
}
func (this_ Node_TypeComparison) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Node_TypeComparison[this_-1].Name_)
}
func init() {
	vals_ := []any{Node_TypeComparison(1), Node_TypeComparison(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node.TypeComparison", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node$TypeComparison", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.TypeComparison", "com.google.javascript.rhino.Node$TypeComparison", func(v_ any) bool {
		return j2gort.InstanceOf[Node_TypeComparison](v_)
	})
}

// Java enum: com.google.javascript.rhino.Node.JsDocComparison
type Node_JsDocComparison uint8

const (
	enum_Node_JsDocComparison_IGNORE  Node_JsDocComparison = 1
	enum_Node_JsDocComparison_COMPARE Node_JsDocComparison = 2
)

var enumBaseTable_Node_JsDocComparison = [...]j2gort.EnumBase{
	{Name_: "IGNORE", Ordinal_: 0},
	{Name_: "COMPARE", Ordinal_: 1},
}
var enumValues_Node_JsDocComparison = []Node_JsDocComparison{1, 2}

func StaticGet_Node_JsDocComparison_IGNORE() Node_JsDocComparison {
	return enum_Node_JsDocComparison_IGNORE
}
func StaticGet_Node_JsDocComparison_COMPARE() Node_JsDocComparison {
	return enum_Node_JsDocComparison_COMPARE
}

func Node_JsDocComparison_Values() []Node_JsDocComparison {
	return j2gort.CloneSlice(enumValues_Node_JsDocComparison)
}

func Node_JsDocComparison_ValueOf(name string) Node_JsDocComparison {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Node_JsDocComparison {
		if enumBaseTable_Node_JsDocComparison[i_].Name_ == name {
			return Node_JsDocComparison(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.Node.JsDocComparison." + name))
}

func (this_ Node_JsDocComparison) IsJavaNil_() bool                              { return this_ == 0 }
func (Node_JsDocComparison) HasIdentityEquals_()                                 {}
func (Node_JsDocComparison) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Node_JsDocComparison) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Node_JsDocComparison) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Node_JsDocComparison) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$JsDocComparison", "JsDocComparison"
}
func (this_ Node_JsDocComparison) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Node_JsDocComparison[this_-1]
}
func (this_ Node_JsDocComparison) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Node_JsDocComparison) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Node_JsDocComparison) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_JsDocComparison[this_-1].Name_
}
func (this_ Node_JsDocComparison) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Node_JsDocComparison[this_-1].Name_
}
func (this_ Node_JsDocComparison) Equals(other_ any) bool {
	o_, ok_ := other_.(Node_JsDocComparison)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Node_JsDocComparison) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Node_JsDocComparison[this_-1].HashCode()
}
func (this_ Node_JsDocComparison) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_JsDocComparison[this_-1].CompareTo(other_)
}
func (this_ Node_JsDocComparison) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Node_JsDocComparison[this_-1].Name_)
}
func init() {
	vals_ := []any{Node_JsDocComparison(1), Node_JsDocComparison(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node.JsDocComparison", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node$JsDocComparison", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.JsDocComparison", "com.google.javascript.rhino.Node$JsDocComparison", func(v_ any) bool {
		return j2gort.InstanceOf[Node_JsDocComparison](v_)
	})
}

// Java enum: com.google.javascript.rhino.Node.SideEffectComparison
type Node_SideEffectComparison uint8

const (
	enum_Node_SideEffectComparison_IGNORE  Node_SideEffectComparison = 1
	enum_Node_SideEffectComparison_COMPARE Node_SideEffectComparison = 2
)

var enumBaseTable_Node_SideEffectComparison = [...]j2gort.EnumBase{
	{Name_: "IGNORE", Ordinal_: 0},
	{Name_: "COMPARE", Ordinal_: 1},
}
var enumValues_Node_SideEffectComparison = []Node_SideEffectComparison{1, 2}

func StaticGet_Node_SideEffectComparison_IGNORE() Node_SideEffectComparison {
	return enum_Node_SideEffectComparison_IGNORE
}
func StaticGet_Node_SideEffectComparison_COMPARE() Node_SideEffectComparison {
	return enum_Node_SideEffectComparison_COMPARE
}

func Node_SideEffectComparison_Values() []Node_SideEffectComparison {
	return j2gort.CloneSlice(enumValues_Node_SideEffectComparison)
}

func Node_SideEffectComparison_ValueOf(name string) Node_SideEffectComparison {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Node_SideEffectComparison {
		if enumBaseTable_Node_SideEffectComparison[i_].Name_ == name {
			return Node_SideEffectComparison(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.Node.SideEffectComparison." + name))
}

func (this_ Node_SideEffectComparison) IsJavaNil_() bool                              { return this_ == 0 }
func (Node_SideEffectComparison) HasIdentityEquals_()                                 {}
func (Node_SideEffectComparison) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Node_SideEffectComparison) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Node_SideEffectComparison) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Node_SideEffectComparison) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$SideEffectComparison", "SideEffectComparison"
}
func (this_ Node_SideEffectComparison) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Node_SideEffectComparison[this_-1]
}
func (this_ Node_SideEffectComparison) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Node_SideEffectComparison) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Node_SideEffectComparison) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_SideEffectComparison[this_-1].Name_
}
func (this_ Node_SideEffectComparison) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Node_SideEffectComparison[this_-1].Name_
}
func (this_ Node_SideEffectComparison) Equals(other_ any) bool {
	o_, ok_ := other_.(Node_SideEffectComparison)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Node_SideEffectComparison) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Node_SideEffectComparison[this_-1].HashCode()
}
func (this_ Node_SideEffectComparison) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Node_SideEffectComparison[this_-1].CompareTo(other_)
}
func (this_ Node_SideEffectComparison) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Node_SideEffectComparison[this_-1].Name_)
}
func init() {
	vals_ := []any{Node_SideEffectComparison(1), Node_SideEffectComparison(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node.SideEffectComparison", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Node$SideEffectComparison", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.SideEffectComparison", "com.google.javascript.rhino.Node$SideEffectComparison", func(v_ any) bool {
		return j2gort.InstanceOf[Node_SideEffectComparison](v_)
	})
}

// Java class: com.google.javascript.rhino.Node.SideEffectFlags
type Node_SideEffectFlags struct {
	value   int32 `gson:"value"`
	idHash_ atomic.Int32
}

func (this_ *Node_SideEffectFlags) BindSelf_(s any) {
	_ = s
}
func (this_ *Node_SideEffectFlags) IsJavaNil_() bool { return this_ == nil }
func (*Node_SideEffectFlags) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$SideEffectFlags", "SideEffectFlags"
}
func (this_ *Node_SideEffectFlags) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Node_SideEffectFlags) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.SideEffectFlags", "com.google.javascript.rhino.Node$SideEffectFlags", func(v_ any) bool {
		return j2gort.InstanceOf[*Node_SideEffectFlags](v_)
	})
}

func StaticGet_Node_SideEffectFlags_MUTATES_GLOBAL_STATE() int32 { return int32(1) }
func StaticGet_Node_SideEffectFlags_MUTATES_THIS() int32         { return int32(2) }
func StaticGet_Node_SideEffectFlags_MUTATES_ARGUMENTS() int32    { return int32(4) }
func StaticGet_Node_SideEffectFlags_THROWS() int32               { return int32(8) }
func StaticGet_Node_SideEffectFlags_USED_BITS_MASK() int32       { return int32(15) }
func StaticGet_Node_SideEffectFlags_NO_SIDE_EFFECTS() int32      { return int32(0) }
func StaticGet_Node_SideEffectFlags_ALL_SIDE_EFFECTS() int32     { return int32(15) }
func (this_ *Node_SideEffectFlags) initFields_Node_SideEffectFlags() {
	_ = this_
	this_.value = int32(15)
}

// Java constructor: SideEffectFlags()
func NewNode_SideEffectFlags() *Node_SideEffectFlags {
	this_ := &Node_SideEffectFlags{}
	InitNode_SideEffectFlags(this_)
	return this_
}

func InitNode_SideEffectFlags(this_ *Node_SideEffectFlags) {
	this_.initFields_Node_SideEffectFlags()
}

// Java constructor: SideEffectFlags(int value)
func NewNode_SideEffectFlags_1(value int32) *Node_SideEffectFlags {
	this_ := &Node_SideEffectFlags{}
	InitNode_SideEffectFlags_1(this_, value)
	return this_
}

func InitNode_SideEffectFlags_1(this_ *Node_SideEffectFlags, value int32) {
	_ = value
	this_.initFields_Node_SideEffectFlags()
	this_.value = value
}

// Java method: SideEffectFlags.valueOf()
func (this_ *Node_SideEffectFlags) ValueOf() int32 {
	return this_.value
}

// Java method: SideEffectFlags.setAllFlags()
func (this_ *Node_SideEffectFlags) SetAllFlags() *Node_SideEffectFlags {
	this_.value = int32(15)
	return this_
}

// Java method: SideEffectFlags.clearAllFlags()
func (this_ *Node_SideEffectFlags) ClearAllFlags() *Node_SideEffectFlags {
	this_.value = int32(0)
	return this_
}

// Java method: SideEffectFlags.setMutatesGlobalState()
func (this_ *Node_SideEffectFlags) SetMutatesGlobalState() *Node_SideEffectFlags {
	this_.value = (this_.value | ((int32(1) | int32(4)) | int32(2)))
	return this_
}

// Java method: SideEffectFlags.setThrows()
func (this_ *Node_SideEffectFlags) SetThrows() *Node_SideEffectFlags {
	this_.value = (this_.value | int32(8))
	return this_
}

// Java method: SideEffectFlags.setMutatesThis()
func (this_ *Node_SideEffectFlags) SetMutatesThis() *Node_SideEffectFlags {
	this_.value = (this_.value | int32(2))
	return this_
}

// Java method: SideEffectFlags.setMutatesArguments()
func (this_ *Node_SideEffectFlags) SetMutatesArguments() *Node_SideEffectFlags {
	this_.value = (this_.value | int32(4))
	return this_
}

// Java method: SideEffectFlags.toString()
func (this_ *Node_SideEffectFlags) String() string {
	var builder j2gort.LocalStringBuilder
	builder.InitStr("Side effects: ")

	if (this_.value & int32(2)) != int32(0) {
		builder.AppendStr("this ")
	}
	if (this_.value & int32(1)) != int32(0) {
		builder.AppendStr("global ")
	}
	if (this_.value & int32(8)) != int32(0) {
		builder.AppendStr("throw ")
	}
	if (this_.value & int32(4)) != int32(0) {
		builder.AppendStr("args ")
	}
	return builder.String()
}

// Java class: com.google.javascript.rhino.Node.ConstantVarFlags
type Node_ConstantVarFlags struct {
	idHash_ atomic.Int32
}

func (this_ *Node_ConstantVarFlags) BindSelf_(s any) {
	_ = s
}
func (this_ *Node_ConstantVarFlags) IsJavaNil_() bool { return this_ == nil }
func (*Node_ConstantVarFlags) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$ConstantVarFlags", "ConstantVarFlags"
}
func (this_ *Node_ConstantVarFlags) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Node_ConstantVarFlags) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Node.ConstantVarFlags", "com.google.javascript.rhino.Node$ConstantVarFlags", func(v_ any) bool {
		return j2gort.InstanceOf[*Node_ConstantVarFlags](v_)
	})
}

func StaticGet_Node_ConstantVarFlags_DECLARED() int32 { return int32(1) }
func StaticGet_Node_ConstantVarFlags_INFERRED() int32 { return int32(2) }
func (this_ *Node_ConstantVarFlags) initFields_Node_ConstantVarFlags() {
	_ = this_
}

// Java constructor: ConstantVarFlags()
func NewNode_ConstantVarFlags() *Node_ConstantVarFlags {
	this_ := &Node_ConstantVarFlags{}
	InitNode_ConstantVarFlags(this_)
	return this_
}

func InitNode_ConstantVarFlags(this_ *Node_ConstantVarFlags) {
	this_.initFields_Node_ConstantVarFlags()
}

// Java class:
type Node__Anon_154 struct {
	outer_  *Node_AncestorIterable
	idHash_ atomic.Int32
}

func (this_ *Node__Anon_154) BindSelf_(s any) {
	_ = s
}
func (this_ *Node__Anon_154) IsJavaNil_() bool { return this_ == nil }
func (*Node__Anon_154) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.Node$AncestorIterable$1", ""
}
func (this_ *Node__Anon_154) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Node__Anon_154) HasIdentityEquals_()                       {}
func (this_ *Node__Anon_154) SetOuter_(o *Node_AncestorIterable)  { this_.outer_ = o }
func (this_ *Node__Anon_154) GetOuter_() *Node_AncestorIterable   { return this_.outer_ }
func (*Node__Anon_154) Implements_j2gort_Ext_java_util_Iterator() {}

func (this_ *Node__Anon_154) initFields_Node__Anon_154() {
	_ = this_
}

func NewNode__Anon_154(outer_ *Node_AncestorIterable) *Node__Anon_154 {
	this_ := &Node__Anon_154{}
	this_.outer_ = outer_
	this_.initFields_Node__Anon_154()
	return this_
}

// Java method: .hasNext()
func (this_ *Node__Anon_154) HasNext() bool {
	return (this_.outer_.cur != 0)
}

// Java method: .next()
func (this_ *Node__Anon_154) Next() any {
	if !(this_.HasNext()) {
		panic(j2gort.AsIface[javalang.INoSuchElementException](j2gort.ConstructShimThrowable("java.util.NoSuchElementException")))
	}
	var n Node = this_.outer_.cur
	_ = n
	this_.outer_.cur = this_.outer_.cur.GetParent()
	return j2gort.BoxNode(n)
}

// Java method: .remove()
func (this_ *Node__Anon_154) Remove() {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}
