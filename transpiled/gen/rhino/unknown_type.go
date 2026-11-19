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
// Port of com/google/javascript/rhino/jstype/UnknownType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.UnknownType
type UnknownType struct {
	ObjectTypeBase
	isChecked bool `gson:"isChecked"`
}

func (this_ *UnknownType) BindSelf_(s any) {
	_ = s
	this_.ObjectTypeBase.BindSelf_(s)
}
func (this_ *UnknownType) IsJavaNil_() bool { return this_ == nil }
func (*UnknownType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.UnknownType", "UnknownType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.UnknownType", "com.google.javascript.rhino.jstype.UnknownType", func(v_ any) bool {
		return j2gort.InstanceOf[*UnknownType](v_)
	})
}
func (*UnknownType) Implements_rhino_ObjectType() {}
func (*UnknownType) Implements_rhino_JSType()     {}

func (this_ *UnknownType) initFields_UnknownType() {
	_ = this_
}

// Java constructor: UnknownType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, boolean isChecked)
func NewUnknownType(registry *JSTypeRegistry, isChecked bool) *UnknownType {
	this_ := &UnknownType{}
	this_.BindSelf_(this_)
	InitUnknownType(this_, registry, isChecked)
	return this_
}

func InitUnknownType(this_ *UnknownType, registry *JSTypeRegistry, isChecked bool) {
	_ = registry
	_ = isChecked
	InitObjectType_2_JSTypeRegistry_int32(&this_.ObjectTypeBase, registry, func() int32 {
		if isChecked {
			return int32(32)
		}
		return int32(0)
	}())
	this_.initFields_UnknownType()
	this_.EagerlyResolveToSelf()
	this_.isChecked = isChecked
}

// Java method: UnknownType.getTypeClass()
func (this_ *UnknownType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_UNKNOWN()
}

// Java method: UnknownType.isUnknownType()
func (this_ *UnknownType) IsUnknownType() bool {
	return true
}

// Java method: UnknownType.canBeCalled()
func (this_ *UnknownType) CanBeCalled() bool {
	return true
}

// Java method: UnknownType.matchesNumberContext()
func (this_ *UnknownType) MatchesNumberContext() bool {
	return true
}

// Java method: UnknownType.matchesObjectContext()
func (this_ *UnknownType) MatchesObjectContext() bool {
	return true
}

// Java method: UnknownType.matchesStringContext()
func (this_ *UnknownType) MatchesStringContext() bool {
	return true
}

// Java method: UnknownType.matchesSymbolContext()
func (this_ *UnknownType) MatchesSymbolContext() bool {
	return true
}

// Java method: UnknownType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *UnknownType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	return pkg_base.StaticGet_Tri_UNKNOWN()
}

// Java method: UnknownType.isNullable()
func (this_ *UnknownType) IsNullable() bool {
	return true
}

// Java method: UnknownType.isVoidable()
func (this_ *UnknownType) IsVoidable() bool {
	return true
}

// Java method: UnknownType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *UnknownType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseUnknownType()
}

// Java method: UnknownType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *UnknownType) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseUnknownType(this_, that)
}

// Java method: UnknownType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *UnknownType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	sb.Append_1_string(this_.GetReferenceName())
}

// Java method: UnknownType.defineProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *UnknownType) DefineProperty_4_string_JSType_bool_Node(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	return true
}

// Java method: UnknownType.defineProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *UnknownType) DefineProperty_4_Property_Key_JSType_bool_Node(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	return true
}

// Java method: UnknownType.getImplicitPrototype()
func (this_ *UnknownType) GetImplicitPrototype() ObjectType {
	return nil
}

// Java method: UnknownType.getConstructor()
func (this_ *UnknownType) GetConstructor() *FunctionType {
	return nil
}

// Java method: UnknownType.getReferenceName()
func (this_ *UnknownType) GetReferenceName() string {
	if this_.isChecked {
		return "??"
	} else {
		return "?"
	}
}

// Java method: UnknownType.getDisplayName()
func (this_ *UnknownType) GetDisplayName() string {
	return "Unknown"
}

// Java method: UnknownType.hasDisplayName()
func (this_ *UnknownType) HasDisplayName() bool {
	return true
}

// Java method: UnknownType.getPossibleToBooleanOutcomes()
func (this_ *UnknownType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_BOTH()
}

// Java method: UnknownType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *UnknownType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: UnknownType.recursionUnsafeHashCode()
func (this_ *UnknownType) RecursionUnsafeHashCode() int32 {
	return j2gort.IdentityHash(this_)
}
