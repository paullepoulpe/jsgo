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
// Port of com/google/javascript/rhino/QualifiedName.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.rhino.QualifiedName
type QualifiedName interface {
	Implements_rhino_QualifiedName()
	IsJavaNil_() bool
	GetQualifiedNameBase() *QualifiedNameBase
	GetOwner() QualifiedName
	GetComponent() string
	IsSimple() bool
	AppendTo(*j2gort.Ext_java_lang_StringBuilder, *j2gort.BoxChar)
	Matches(Node) bool
	GetRoot() string
	Components() j2gort.Ext_java_lang_Iterable
	Join_0() string
	Join_1(*j2gort.BoxChar) string
	GetComponentCount() int32
	Getprop(string) QualifiedName
}

func ToIface_QualifiedName[T interface {
	QualifiedName
	comparable
}](v T) QualifiedName {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.rhino.QualifiedName
type QualifiedNameBase struct {
	self_   QualifiedName
	idHash_ atomic.Int32
}

func (this_ *QualifiedNameBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(QualifiedName)
}
func (this_ *QualifiedNameBase) IsJavaNil_() bool { return this_ == nil }
func (*QualifiedNameBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.QualifiedName", "QualifiedName"
}
func (this_ *QualifiedNameBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.QualifiedName", "com.google.javascript.rhino.QualifiedName", func(v_ any) bool {
		return j2gort.InstanceOf[QualifiedName](v_)
	})
}
func (b *QualifiedNameBase) GetQualifiedNameBase() *QualifiedNameBase { return b }
func (*QualifiedNameBase) Implements_rhino_QualifiedName()            {}

var static_QualifiedName_THIS string

func StaticGet_QualifiedName_THIS() string {
	clinit_QualifiedName()
	return static_QualifiedName_THIS
}
func StaticSet_QualifiedName_THIS(v string) {
	clinit_QualifiedName()
	static_QualifiedName_THIS = v
}

var static_QualifiedName_SUPER string

func StaticGet_QualifiedName_SUPER() string {
	clinit_QualifiedName()
	return static_QualifiedName_SUPER
}
func StaticSet_QualifiedName_SUPER(v string) {
	clinit_QualifiedName()
	static_QualifiedName_SUPER = v
}

var clinitDone_QualifiedName uint32
var clinitState_QualifiedName j2gort.ClinitState

func clinit_QualifiedName() {
	if atomic.LoadUint32(&clinitDone_QualifiedName) != 0 {
		return
	}
	clinitSlow_QualifiedName()
}

func clinitSlow_QualifiedName() {
	if !j2gort.EnterClinit(&clinitState_QualifiedName, &clinitDone_QualifiedName) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_QualifiedName, &clinitDone_QualifiedName)
	static_QualifiedName_THIS = RhinoStringPool_AddOrGet("this")
	static_QualifiedName_SUPER = RhinoStringPool_AddOrGet("super")
}

func (this_ *QualifiedNameBase) initFields_QualifiedName() {
	_ = this_
}

func InitQualifiedName(this_ *QualifiedNameBase) {
	this_.initFields_QualifiedName()
}

// Java static method: QualifiedName.of(java.lang.String string)
func QualifiedName_Of(string_ string) QualifiedName {
	clinit_QualifiedName()
	_ = string_
	var lastIndex int32 = int32(0)
	_ = lastIndex
	var index int32 = int32(0)
	_ = index
	var builder []string = []string{}
	_ = builder
	for doFirst_1 := true; doFirst_1 || (index >= int32(0)); doFirst_1 = false {
		index = j2gort.JStrIndexOfCpFrom(string_, int32(uint16(46)), lastIndex)
		var term string = j2gort.JStrSub(string_, lastIndex, func() int32 {
			if index < int32(0) {
				return j2gort.JStrLen(string_)
			}
			return index
		}())
		_ = term
		builder = append(builder, j2gort.CheckNotNullString(RhinoStringPool_AddOrGet(term)))
		lastIndex = (index + int32(1))
	}
	var terms []string = builder
	_ = terms
	return NewQualifiedName_StringListQname(j2gort.ImmListFromTypedSlice(terms), int32(len(terms)))
}

// Java method: QualifiedName.getOwner()
func (this_ *QualifiedNameBase) GetOwner() QualifiedName {
	panic("abstract or native method invoked")
}

// Java method: QualifiedName.getComponent()
func (this_ *QualifiedNameBase) GetComponent() string {
	panic("abstract or native method invoked")
}

// Java method: QualifiedName.isSimple()
func (this_ *QualifiedNameBase) IsSimple() bool {
	panic("abstract or native method invoked")
}

// Java method: QualifiedName.appendTo(java.lang.StringBuilder sb, java.lang.Character separator)
func (this_ *QualifiedNameBase) AppendTo(sb *j2gort.Ext_java_lang_StringBuilder, separator *j2gort.BoxChar) {
	_ = sb
	_ = separator
	panic("abstract or native method invoked")
}

// Java method: QualifiedName.matches(com.google.javascript.rhino.Node n)
func (this_ *QualifiedNameBase) Matches(n Node) bool {
	_ = n
	panic("abstract or native method invoked")
}

// Java method: QualifiedName.getRoot()
func (this_ *QualifiedNameBase) GetRoot() string {
	var name QualifiedName = this_.self_
	_ = name
	for !(name.IsSimple()) {
		name = name.GetOwner()
	}
	return name.GetComponent()
}

// Java method: QualifiedName.components()
func (this_ *QualifiedNameBase) Components() j2gort.Ext_java_lang_Iterable {
	var components *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = components
	this_.priv_QualifiedName_buildComponents_8770(components)
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](components.BuildImmList_())
}

// Java method: QualifiedName.buildComponents(com.google.common.collect.ImmutableList.Builder<java.lang.String> builder)
func (this_ *QualifiedNameBase) priv_QualifiedName_buildComponents_8770(builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder) {
	_ = builder
	var owner QualifiedName = this_.self_.GetOwner()
	_ = owner
	if owner != nil && !owner.IsJavaNil_() {
		owner.GetQualifiedNameBase().priv_QualifiedName_buildComponents_8770(builder)
	}
	builder.AddStrVoid_(this_.self_.GetComponent())
}

// Java method: QualifiedName.join()
func (this_ *QualifiedNameBase) Join_0() string {
	return this_.Join_1(j2gort.ValueOfChar(uint16(46)))
}

// Java method: QualifiedName.join(java.lang.Character separator)
func (this_ *QualifiedNameBase) Join_1(separator *j2gort.BoxChar) string {
	_ = separator
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	this_.self_.AppendTo(sb, separator)
	return sb.String()
}

// Java method: QualifiedName.getComponentCount()
func (this_ *QualifiedNameBase) GetComponentCount() int32 {
	panic("abstract or native method invoked")
}

// Java method: QualifiedName.getprop(java.lang.String propertyName)
func (this_ *QualifiedNameBase) Getprop(propertyName string) QualifiedName {
	_ = propertyName
	return NewQualifiedName_GetpropQname(this_.self_, propertyName)
}

// Java class: com.google.javascript.rhino.QualifiedName.StringListQname
type QualifiedName_StringListQname struct {
	QualifiedNameBase
	terms []string `gson:"terms"`
	size  int32    `gson:"size"`
}

func (this_ *QualifiedName_StringListQname) BindSelf_(s any) {
	_ = s
	this_.QualifiedNameBase.BindSelf_(s)
}
func (this_ *QualifiedName_StringListQname) IsJavaNil_() bool { return this_ == nil }
func (*QualifiedName_StringListQname) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.QualifiedName$StringListQname", "StringListQname"
}
func (this_ *QualifiedName_StringListQname) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.QualifiedNameBase.IdentityHash_()
}
func (*QualifiedName_StringListQname) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.QualifiedName.StringListQname", "com.google.javascript.rhino.QualifiedName$StringListQname", func(v_ any) bool {
		return j2gort.InstanceOf[*QualifiedName_StringListQname](v_)
	})
}
func (*QualifiedName_StringListQname) Implements_rhino_QualifiedName() {}

func (this_ *QualifiedName_StringListQname) initFields_QualifiedName_StringListQname() {
	_ = this_
}

// Java constructor: StringListQname(com.google.common.collect.ImmutableList<java.lang.String> terms, int size)
func NewQualifiedName_StringListQname(terms j2gort.Ext_com_google_common_collect_ImmutableList, size int32) *QualifiedName_StringListQname {
	this_ := &QualifiedName_StringListQname{}
	this_.BindSelf_(this_)
	InitQualifiedName_StringListQname(this_, terms, size)
	return this_
}

func InitQualifiedName_StringListQname(this_ *QualifiedName_StringListQname, terms j2gort.Ext_com_google_common_collect_ImmutableList, size int32) {
	_ = terms
	_ = size
	InitQualifiedName(&this_.QualifiedNameBase)
	this_.initFields_QualifiedName_StringListQname()
	this_.terms = j2gort.ImmListToSliceNullable[string](terms)
	this_.size = size
}

// Java method: StringListQname.getOwner()
func (this_ *QualifiedName_StringListQname) GetOwner() QualifiedName {
	if this_.size > int32(1) {
		return NewQualifiedName_StringListQname(j2gort.ImmListFromTypedSlice(this_.terms), (this_.size - int32(1)))
	} else {
		return nil
	}
}

// Java method: StringListQname.getComponent()
func (this_ *QualifiedName_StringListQname) GetComponent() string {
	return this_.terms[(this_.size - int32(1))]
}

// Java method: StringListQname.isSimple()
func (this_ *QualifiedName_StringListQname) IsSimple() bool {
	return (this_.size == int32(1))
}

// Java method: StringListQname.appendTo(java.lang.StringBuilder sb, java.lang.Character separator)
func (this_ *QualifiedName_StringListQname) AppendTo(sb *j2gort.Ext_java_lang_StringBuilder, separator *j2gort.BoxChar) {
	_ = sb
	_ = separator
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.size; i = i + 1 {
			if i > int32(0) {
				sb.Append_1_any_9831(j2gort.PtrToAny(separator))
			}
			sb.Append_1_string(this_.terms[i])
		}
	}
}

// Java method: StringListQname.components()
func (this_ *QualifiedName_StringListQname) Components() j2gort.Ext_java_lang_Iterable {
	return j2gort.ImmListFromTypedSlice(j2gort.ImmSliceSubList(this_.terms, int32(0), this_.size))
}

// Java method: StringListQname.getComponentCount()
func (this_ *QualifiedName_StringListQname) GetComponentCount() int32 {
	return this_.size
}

// Java method: StringListQname.matches(com.google.javascript.rhino.Node n)
func (this_ *QualifiedName_StringListQname) Matches(n Node) bool {
	_ = n
	var pos int32 = (this_.size - int32(1))
	_ = pos
	for (pos > int32(0)) && n.IsGetProp() {
		if !(RhinoStringPool_UncheckedEquals(n.GetString(), this_.terms[pos])) {
			return false
		}
		pos = pos - 1
		n = n.GetFirstChild()
	}
	if pos > int32(0) {
		return false
	}
	var term string = this_.terms[int32(0)]
	_ = term
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		return RhinoStringPool_UncheckedEquals(term, n.GetString())
	case int32(37) /* Token.THIS */ :
		return RhinoStringPool_UncheckedEquals(term, StaticGet_QualifiedName_THIS())
	case int32(110) /* Token.SUPER */ :
		return RhinoStringPool_UncheckedEquals(term, StaticGet_QualifiedName_SUPER())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java class: com.google.javascript.rhino.QualifiedName.GetpropQname
type QualifiedName_GetpropQname struct {
	QualifiedNameBase
	owner QualifiedName `gson:"owner"`
	prop  string        `gson:"prop"`
}

func (this_ *QualifiedName_GetpropQname) BindSelf_(s any) {
	_ = s
	this_.QualifiedNameBase.BindSelf_(s)
}
func (this_ *QualifiedName_GetpropQname) IsJavaNil_() bool { return this_ == nil }
func (*QualifiedName_GetpropQname) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.QualifiedName$GetpropQname", "GetpropQname"
}
func (this_ *QualifiedName_GetpropQname) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.QualifiedNameBase.IdentityHash_()
}
func (*QualifiedName_GetpropQname) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.QualifiedName.GetpropQname", "com.google.javascript.rhino.QualifiedName$GetpropQname", func(v_ any) bool {
		return j2gort.InstanceOf[*QualifiedName_GetpropQname](v_)
	})
}
func (*QualifiedName_GetpropQname) Implements_rhino_QualifiedName() {}

func (this_ *QualifiedName_GetpropQname) initFields_QualifiedName_GetpropQname() {
	_ = this_
	this_.prop = j2gort.NullString
}

// Java constructor: GetpropQname(com.google.javascript.rhino.QualifiedName owner, java.lang.String prop)
func NewQualifiedName_GetpropQname(owner QualifiedName, prop string) *QualifiedName_GetpropQname {
	this_ := &QualifiedName_GetpropQname{}
	this_.BindSelf_(this_)
	InitQualifiedName_GetpropQname(this_, owner, prop)
	return this_
}

func InitQualifiedName_GetpropQname(this_ *QualifiedName_GetpropQname, owner QualifiedName, prop string) {
	_ = owner
	_ = prop
	InitQualifiedName(&this_.QualifiedNameBase)
	this_.initFields_QualifiedName_GetpropQname()
	this_.owner = owner
	this_.prop = RhinoStringPool_AddOrGet(prop)
}

// Java method: GetpropQname.getOwner()
func (this_ *QualifiedName_GetpropQname) GetOwner() QualifiedName {
	return this_.owner
}

// Java method: GetpropQname.getComponent()
func (this_ *QualifiedName_GetpropQname) GetComponent() string {
	return this_.prop
}

// Java method: GetpropQname.isSimple()
func (this_ *QualifiedName_GetpropQname) IsSimple() bool {
	return false
}

// Java method: GetpropQname.appendTo(java.lang.StringBuilder sb, java.lang.Character separator)
func (this_ *QualifiedName_GetpropQname) AppendTo(sb *j2gort.Ext_java_lang_StringBuilder, separator *j2gort.BoxChar) {
	_ = sb
	_ = separator
	this_.owner.AppendTo(sb, separator)
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_any_9831(j2gort.PtrToAny(separator))).Append_1_string(this_.prop)
}

// Java method: GetpropQname.matches(com.google.javascript.rhino.Node n)
func (this_ *QualifiedName_GetpropQname) Matches(n Node) bool {
	_ = n
	return ((n.IsGetProp() && RhinoStringPool_UncheckedEquals(n.GetString(), this_.prop)) && this_.owner.Matches(n.GetFirstChild()))
}

// Java method: GetpropQname.getComponentCount()
func (this_ *QualifiedName_GetpropQname) GetComponentCount() int32 {
	return (this_.owner.GetComponentCount() + int32(1))
}

// Java class: com.google.javascript.rhino.QualifiedName.NodeQname
type QualifiedName_NodeQname struct {
	QualifiedNameBase
	node Node `gson:"node"`
}

func (this_ *QualifiedName_NodeQname) BindSelf_(s any) {
	_ = s
	this_.QualifiedNameBase.BindSelf_(s)
}
func (this_ *QualifiedName_NodeQname) IsJavaNil_() bool { return this_ == nil }
func (*QualifiedName_NodeQname) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.QualifiedName$NodeQname", "NodeQname"
}
func (this_ *QualifiedName_NodeQname) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.QualifiedNameBase.IdentityHash_()
}
func (*QualifiedName_NodeQname) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.QualifiedName.NodeQname", "com.google.javascript.rhino.QualifiedName$NodeQname", func(v_ any) bool {
		return j2gort.InstanceOf[*QualifiedName_NodeQname](v_)
	})
}
func (*QualifiedName_NodeQname) Implements_rhino_QualifiedName() {}

func (this_ *QualifiedName_NodeQname) initFields_QualifiedName_NodeQname() {
	_ = this_
}

// Java constructor: NodeQname(com.google.javascript.rhino.Node n)
func NewQualifiedName_NodeQname(n Node) *QualifiedName_NodeQname {
	this_ := &QualifiedName_NodeQname{}
	this_.BindSelf_(this_)
	InitQualifiedName_NodeQname(this_, n)
	return this_
}

func InitQualifiedName_NodeQname(this_ *QualifiedName_NodeQname, n Node) {
	_ = n
	InitQualifiedName(&this_.QualifiedNameBase)
	this_.initFields_QualifiedName_NodeQname()
	this_.node = n
}

// Java method: NodeQname.getOwner()
func (this_ *QualifiedName_NodeQname) GetOwner() QualifiedName {
	if this_.node.IsGetProp() {
		return NewQualifiedName_NodeQname(this_.node.GetFirstChild())
	} else {
		return nil
	}
}

// Java method: NodeQname.getComponent()
func (this_ *QualifiedName_NodeQname) GetComponent() string {
	switch this_.node.TokenSwitchOrd_() {
	case int32(37) /* Token.THIS */ :
		return StaticGet_QualifiedName_THIS()
	case int32(110) /* Token.SUPER */ :
		return StaticGet_QualifiedName_SUPER()
	case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		return this_.node.GetString()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Not a qualified name: "+j2gort.StrVal(this_.node))))
	}
	panic("unreachable switch expression")
}

// Java method: NodeQname.isSimple()
func (this_ *QualifiedName_NodeQname) IsSimple() bool {
	return !(this_.node.IsGetProp())
}

// Java method: NodeQname.appendTo(java.lang.StringBuilder sb, java.lang.Character separator)
func (this_ *QualifiedName_NodeQname) AppendTo(sb *j2gort.Ext_java_lang_StringBuilder, separator *j2gort.BoxChar) {
	_ = sb
	_ = separator
	var components j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.Components().Iterator())
	_ = components
	sb.Append_1_string(j2gort.CastString(components.Next()))
	for components.HasNext() {
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](sb.Append_1_any_9831(j2gort.PtrToAny(separator))).Append_1_string(j2gort.CastString(components.Next()))
	}
}

// Java method: NodeQname.join()
func (this_ *QualifiedName_NodeQname) Join_0() string {
	return this_.node.GetQualifiedName()
}

// Java method: NodeQname.matches(com.google.javascript.rhino.Node n)
func (this_ *QualifiedName_NodeQname) Matches(n Node) bool {
	_ = n
	return n.MatchesQualifiedName_1_Node(this_.node)
}

// Java method: NodeQname.getComponentCount()
func (this_ *QualifiedName_NodeQname) GetComponentCount() int32 {
	var count int32 = int32(1)
	_ = count
	var current Node = this_.node
	_ = current
	for current.IsGetProp() {
		count = count + 1
		current = current.GetFirstChild()
	}
	return count
}
