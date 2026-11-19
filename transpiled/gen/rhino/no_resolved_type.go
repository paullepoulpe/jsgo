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
// Port of com/google/javascript/rhino/jstype/NoResolvedType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.NoResolvedType
type NoResolvedType struct {
	ObjectTypeBase
	referenceName string                                             `gson:"referenceName"`
	templateTypes j2gort.Ext_com_google_common_collect_ImmutableList `gson:"templateTypes"`
}

func (this_ *NoResolvedType) BindSelf_(s any) {
	_ = s
	this_.ObjectTypeBase.BindSelf_(s)
}
func (this_ *NoResolvedType) IsJavaNil_() bool { return this_ == nil }
func (*NoResolvedType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.NoResolvedType", "NoResolvedType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NoResolvedType", "com.google.javascript.rhino.jstype.NoResolvedType", func(v_ any) bool {
		return j2gort.InstanceOf[*NoResolvedType](v_)
	})
}
func (*NoResolvedType) Implements_rhino_ObjectType() {}
func (*NoResolvedType) Implements_rhino_JSType()     {}

func (this_ *NoResolvedType) initFields_NoResolvedType() {
	_ = this_
	this_.referenceName = j2gort.NullString
}

// Java constructor: NoResolvedType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.lang.String referenceName, com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> templateTypes)
func NewNoResolvedType(registry *JSTypeRegistry, referenceName string, templateTypes j2gort.Ext_com_google_common_collect_ImmutableList) *NoResolvedType {
	this_ := &NoResolvedType{}
	this_.BindSelf_(this_)
	InitNoResolvedType(this_, registry, referenceName, templateTypes)
	return this_
}

func InitNoResolvedType(this_ *NoResolvedType, registry *JSTypeRegistry, referenceName string, templateTypes j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = registry
	_ = referenceName
	_ = templateTypes
	InitObjectType_2_JSTypeRegistry_int32(&this_.ObjectTypeBase, registry, int32(8))
	this_.initFields_NoResolvedType()
	this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) | int32(int8(32))))
	this_.referenceName = j2gort.CheckNotNullString(referenceName)
	this_.templateTypes = templateTypes
	this_.EagerlyResolveToSelf()
}

// Java method: NoResolvedType.getTypeClass()
func (this_ *NoResolvedType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_NO_RESOLVED()
}

// Java method: NoResolvedType.getReferenceName()
func (this_ *NoResolvedType) GetReferenceName() string {
	return this_.referenceName
}

// Java method: NoResolvedType.getTemplateTypes()
func (this_ *NoResolvedType) GetTemplateTypes() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.templateTypes
}

// Java method: NoResolvedType.getPossibleToBooleanOutcomes()
func (this_ *NoResolvedType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_BOTH()
}

// Java method: NoResolvedType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *NoResolvedType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	return pkg_base.StaticGet_Tri_UNKNOWN()
}

// Java method: NoResolvedType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *NoResolvedType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: NoResolvedType.getConstructor()
func (this_ *NoResolvedType) GetConstructor() *FunctionType {
	return nil
}

// Java method: NoResolvedType.getImplicitPrototype()
func (this_ *NoResolvedType) GetImplicitPrototype() ObjectType {
	return nil
}

// Java method: NoResolvedType.recursionUnsafeHashCode()
func (this_ *NoResolvedType) RecursionUnsafeHashCode() int32 {
	return j2gort.IdentityHash(this_)
}

// Java method: NoResolvedType.defineProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *NoResolvedType) DefineProperty_4_string_JSType_bool_Node(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	return true
}

// Java method: NoResolvedType.defineProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *NoResolvedType) DefineProperty_4_Property_Key_JSType_bool_Node(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	return true
}

// Java method: NoResolvedType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *NoResolvedType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	sb.Append_1_string(func() string {
		if sb.IsForAnnotations() {
			return "?"
		}
		return ("NoResolvedType<" + j2gort.StrValString(this_.referenceName) + ">")
	}())
}
