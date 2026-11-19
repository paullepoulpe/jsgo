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
// Port of com/google/javascript/rhino/jstype/AllType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.AllType
type AllType struct {
	JSTypeBase
}

func (this_ *AllType) BindSelf_(s any) {
	_ = s
	this_.JSTypeBase.BindSelf_(s)
}
func (this_ *AllType) IsJavaNil_() bool { return this_ == nil }
func (*AllType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.AllType", "AllType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.AllType", "com.google.javascript.rhino.jstype.AllType", func(v_ any) bool {
		return j2gort.InstanceOf[*AllType](v_)
	})
}
func (*AllType) Implements_rhino_JSType() {}

func (this_ *AllType) initFields_AllType() {
	_ = this_
}

// Java constructor: AllType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewAllType(registry *JSTypeRegistry) *AllType {
	this_ := &AllType{}
	this_.BindSelf_(this_)
	InitAllType(this_, registry)
	return this_
}

func InitAllType(this_ *AllType, registry *JSTypeRegistry) {
	_ = registry
	InitJSType_3(&this_.JSTypeBase, registry, nil, int32(16))
	this_.initFields_AllType()
	this_.EagerlyResolveToSelf()
}

// Java method: AllType.getTypeClass()
func (this_ *AllType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_ALL()
}

// Java method: AllType.matchesStringContext()
func (this_ *AllType) MatchesStringContext() bool {
	return true
}

// Java method: AllType.matchesObjectContext()
func (this_ *AllType) MatchesObjectContext() bool {
	return true
}

// Java method: AllType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *AllType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	return pkg_base.StaticGet_Tri_UNKNOWN()
}

// Java method: AllType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *AllType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	sb.Append_1_string("*")
}

// Java method: AllType.getDisplayName()
func (this_ *AllType) GetDisplayName() string {
	return "<Any Type>"
}

// Java method: AllType.hasDisplayName()
func (this_ *AllType) HasDisplayName() bool {
	return true
}

// Java method: AllType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *AllType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseAllType()
}

// Java method: AllType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *AllType) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseAllType(that)
}

// Java method: AllType.getPossibleToBooleanOutcomes()
func (this_ *AllType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_BOTH()
}

// Java method: AllType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *AllType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: AllType.recursionUnsafeHashCode()
func (this_ *AllType) RecursionUnsafeHashCode() int32 {
	return j2gort.IdentityHash(this_)
}

// Java method: AllType.isNullable()
func (this_ *AllType) IsNullable() bool {
	return true
}

// Java method: AllType.isVoidable()
func (this_ *AllType) IsVoidable() bool {
	return true
}
