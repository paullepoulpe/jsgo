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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/EnumType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.EnumType
type EnumType struct {
	PrototypeObjectTypeBase
	elementType  *EnumElementType         `gson:"elementType"`
	elements     j2gort.Ext_java_util_Set `gson:"elements"`
	googModuleId string                   `gson:"googModuleId"`
}

func (this_ *EnumType) BindSelf_(s any) {
	_ = s
	this_.PrototypeObjectTypeBase.BindSelf_(s)
}
func (this_ *EnumType) IsJavaNil_() bool { return this_ == nil }
func (*EnumType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.EnumType", "EnumType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.EnumType", "com.google.javascript.rhino.jstype.EnumType", func(v_ any) bool {
		return j2gort.InstanceOf[*EnumType](v_)
	})
}
func (*EnumType) Implements_rhino_PrototypeObjectType()  {}
func (*EnumType) Implements_rhino_JSType_WithSourceRef() {}
func (*EnumType) Implements_rhino_ObjectType()           {}
func (*EnumType) Implements_rhino_JSType()               {}

var static_EnumType_TYPE_CLASS JSTypeClass

func StaticGet_EnumType_TYPE_CLASS() JSTypeClass {
	clinit_EnumType()
	return static_EnumType_TYPE_CLASS
}
func StaticSet_EnumType_TYPE_CLASS(v JSTypeClass) {
	clinit_EnumType()
	static_EnumType_TYPE_CLASS = v
}

var clinitDone_EnumType uint32
var clinitState_EnumType j2gort.ClinitState

func clinit_EnumType() {
	if atomic.LoadUint32(&clinitDone_EnumType) != 0 {
		return
	}
	clinitSlow_EnumType()
}

func clinitSlow_EnumType() {
	if !j2gort.EnterClinit(&clinitState_EnumType, &clinitDone_EnumType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_EnumType, &clinitDone_EnumType)
	static_EnumType_TYPE_CLASS = StaticGet_JSTypeClass_ENUM()
}

func (this_ *EnumType) initFields_EnumType() {
	_ = this_
	this_.elements = j2gort.NewExt_java_util_LinkedHashSet()
	this_.googModuleId = j2gort.NullString
}

// Java constructor: EnumType(com.google.javascript.rhino.jstype.EnumType.Builder builder)
func NewEnumType(builder *EnumType_Builder) *EnumType {
	this_ := &EnumType{}
	this_.BindSelf_(this_)
	InitEnumType(this_, builder)
	return this_
}

func InitEnumType(this_ *EnumType, builder *EnumType_Builder) {
	_ = builder
	InitPrototypeObjectType_2(&this_.PrototypeObjectTypeBase, ToIface_PrototypeObjectType_Builder(builder), int32(512))
	this_.initFields_EnumType()
	this_.elementType = NewEnumElementType(builder.GetPrototypeObjectType_BuilderBase().registry, builder.elementType, builder.elementName, this_)
	this_.GetJSTypeBase().subSourceNode = builder.source
	this_.googModuleId = builder.googModuleId
	this_.GetJSTypeBase().registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_EnumType_TYPE_CLASS())
}

// Java static method: EnumType.builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func EnumType_S_Builder(registry *JSTypeRegistry) *EnumType_Builder {
	_ = registry
	return NewEnumType_Builder(registry)
}

// Java method: EnumType.getTypeClass()
func (this_ *EnumType) GetTypeClass() JSTypeClass {
	return StaticGet_EnumType_TYPE_CLASS()
}

// Java method: EnumType.toMaybeEnumType()
func (this_ *EnumType) ToMaybeEnumType() *EnumType {
	return this_
}

// Java method: EnumType.getImplicitPrototype()
func (this_ *EnumType) GetImplicitPrototype() ObjectType {
	return this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE())
}

// Java method: EnumType.getElements()
func (this_ *EnumType) GetElements() j2gort.Ext_java_util_Set {
	return j2gort.Ext_java_util_Collections_UnmodifiableSet(this_.elements)
}

// Java method: EnumType.defineElement(java.lang.String name, com.google.javascript.rhino.Node definingNode)
func (this_ *EnumType) DefineElement(name string, definingNode Node) bool {
	_ = name
	_ = definingNode
	j2gort.SetAddStr(this_.elements, name)
	return this_.DefineDeclaredProperty_3_string_JSType_Node(name, ToIface_JSType(this_.elementType), definingNode)
}

// Java method: EnumType.getElementsType()
func (this_ *EnumType) GetElementsType() *EnumElementType {
	return this_.elementType
}

// Java method: EnumType.getEnumeratedTypeOfEnumObject()
func (this_ *EnumType) GetEnumeratedTypeOfEnumObject() JSType {
	return this_.elementType.GetPrimitiveType()
}

// Java method: EnumType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *EnumType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	var result pkg_base.Tri = this_.PrototypeObjectTypeBase.TestForEquality(that)
	_ = result
	if result != 0 {
		return result
	}
	if this_.GetJSTypeBase().EqualsJSType(that) {
		return pkg_base.StaticGet_Tri_TRUE()
	} else {
		return pkg_base.StaticGet_Tri_FALSE()
	}
}

// Java method: EnumType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *EnumType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	sb.Append_1_string(func() string {
		if sb.IsForAnnotations() {
			return "!Object"
		}
		return this_.GetReferenceName()
	}())
}

// Java method: EnumType.getDisplayName()
func (this_ *EnumType) GetDisplayName() string {
	return this_.elementType.GetDisplayName()
}

// Java method: EnumType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *EnumType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseObjectType(this_)
}

// Java method: EnumType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *EnumType) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseObjectType(this_, that)
}

// Java method: EnumType.getConstructor()
func (this_ *EnumType) GetConstructor() *FunctionType {
	return nil
}

// Java method: EnumType.matchesNumberContext()
func (this_ *EnumType) MatchesNumberContext() bool {
	return false
}

// Java method: EnumType.matchesStringContext()
func (this_ *EnumType) MatchesStringContext() bool {
	return true
}

// Java method: EnumType.matchesObjectContext()
func (this_ *EnumType) MatchesObjectContext() bool {
	return true
}

// Java method: EnumType.getSource()
func (this_ *EnumType) GetSource() Node {
	return this_.GetJSTypeBase().subSourceNode
}

// Java method: EnumType.getGoogModuleId()
func (this_ *EnumType) GetGoogModuleId() string {
	return this_.googModuleId
}

// Java method: EnumType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *EnumType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	this_.elementType = j2gort.CastPtr[EnumElementType](this_.elementType.Resolve(reporter))
	return this_.PrototypeObjectTypeBase.ResolveInternal(reporter)
}

// Java class: com.google.javascript.rhino.jstype.EnumType.Builder
type EnumType_Builder struct {
	PrototypeObjectType_BuilderBase
	elementName  string `gson:"elementName"`
	googModuleId string `gson:"googModuleId"`
	elementType  JSType `gson:"elementType"`
	source       Node   `gson:"source"`
}

func (this_ *EnumType_Builder) BindSelf_(s any) {
	_ = s
	this_.PrototypeObjectType_BuilderBase.BindSelf_(s)
}
func (this_ *EnumType_Builder) IsJavaNil_() bool { return this_ == nil }
func (*EnumType_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.EnumType$Builder", "Builder"
}
func (this_ *EnumType_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PrototypeObjectType_BuilderBase.IdentityHash_()
}
func (*EnumType_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.EnumType.Builder", "com.google.javascript.rhino.jstype.EnumType$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*EnumType_Builder](v_)
	})
}
func (*EnumType_Builder) Implements_rhino_PrototypeObjectType_Builder() {}

func (this_ *EnumType_Builder) initFields_EnumType_Builder() {
	_ = this_
	this_.elementName = j2gort.NullString
	this_.googModuleId = j2gort.NullString
}

// Java constructor: Builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewEnumType_Builder(registry *JSTypeRegistry) *EnumType_Builder {
	this_ := &EnumType_Builder{}
	this_.BindSelf_(this_)
	InitEnumType_Builder(this_, registry)
	return this_
}

func InitEnumType_Builder(this_ *EnumType_Builder, registry *JSTypeRegistry) {
	_ = registry
	InitPrototypeObjectType_Builder(&this_.PrototypeObjectType_BuilderBase, registry)
	this_.initFields_EnumType_Builder()
}

// Java method: Builder.setName(java.lang.String x)
func (this_ *EnumType_Builder) SetName(x string) any {
	_ = x
	this_.PrototypeObjectType_BuilderBase.SetName(("enum{" + j2gort.StrValString(x) + "}"))
	this_.elementName = x
	return this_
}

// Java method: Builder.setSource(com.google.javascript.rhino.Node x)
func (this_ *EnumType_Builder) SetSource(x Node) *EnumType_Builder {
	_ = x
	this_.source = x
	return this_
}

// Java method: Builder.setGoogModuleId(java.lang.String x)
func (this_ *EnumType_Builder) SetGoogModuleId(x string) *EnumType_Builder {
	_ = x
	this_.googModuleId = x
	return this_
}

// Java method: Builder.setElementType(com.google.javascript.rhino.jstype.JSType x)
func (this_ *EnumType_Builder) SetElementType(x JSType) *EnumType_Builder {
	_ = x
	this_.elementType = x
	return this_
}

// Java method: Builder.build()
func (this_ *EnumType_Builder) Build() any {
	return NewEnumType(this_)
}
