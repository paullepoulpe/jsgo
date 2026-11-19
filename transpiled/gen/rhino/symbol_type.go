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
// Port of com/google/javascript/rhino/jstype/SymbolType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.rhino.jstype.SymbolType
type SymbolType interface {
	Implements_rhino_SymbolType()
	IsJavaNil_() bool
	GetSymbolTypeBase() *SymbolTypeBase
	ValueType
}

func ToIface_SymbolType[T interface {
	SymbolType
	comparable
}](v T) SymbolType {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SymbolType[T SymbolType](v T) SymbolType {
	var vi SymbolType = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.jstype.SymbolType
type SymbolTypeBase struct {
	ValueTypeBase
}

func (this_ *SymbolTypeBase) BindSelf_(s any) {
	_ = s
	this_.ValueTypeBase.BindSelf_(s)
}
func (this_ *SymbolTypeBase) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTypeBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.SymbolType", "SymbolType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.SymbolType", "com.google.javascript.rhino.jstype.SymbolType", func(v_ any) bool {
		return j2gort.InstanceOf[SymbolType](v_)
	})
}
func (b *SymbolTypeBase) GetSymbolTypeBase() *SymbolTypeBase { return b }
func (this_ *SymbolTypeBase) self_SymbolType_() SymbolType {
	if s := this_.JSTypeBase.self_; s != nil {
		return s.(SymbolType)
	}
	return nil
}
func (*SymbolTypeBase) Implements_rhino_SymbolType() {}
func (*SymbolTypeBase) Implements_rhino_ValueType()  {}
func (*SymbolTypeBase) Implements_rhino_JSType()     {}

func (this_ *SymbolTypeBase) initFields_SymbolType() {
	_ = this_
}

// Java constructor: SymbolType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewSymbolType_1(registry *JSTypeRegistry) SymbolType {
	this_ := &SymbolTypeBase{}
	this_.BindSelf_(this_)
	InitSymbolType_1(this_, registry)
	return this_
}

func InitSymbolType_1(this_ *SymbolTypeBase, registry *JSTypeRegistry) {
	_ = registry
	InitSymbolType_2(this_, registry, int32(1048576))
}

// Java constructor: SymbolType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, int typeKindMask)
func NewSymbolType_2(registry *JSTypeRegistry, typeKindMask int32) SymbolType {
	this_ := &SymbolTypeBase{}
	this_.BindSelf_(this_)
	InitSymbolType_2(this_, registry, typeKindMask)
	return this_
}

func InitSymbolType_2(this_ *SymbolTypeBase, registry *JSTypeRegistry, typeKindMask int32) {
	_ = registry
	_ = typeKindMask
	InitValueType(&this_.ValueTypeBase, registry, typeKindMask)
	this_.initFields_SymbolType()
}

// Java method: SymbolType.getTypeClass()
func (this_ *SymbolTypeBase) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_SYMBOL()
}

// Java method: SymbolType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *SymbolTypeBase) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	var result pkg_base.Tri = this_.ValueTypeBase.TestForEquality(that)
	_ = result
	if result != 0 {
		return result
	}
	base_that_ := that.GetJSTypeBase()
	if base_that_.CanCastTo(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_TYPE())) || base_that_.CanCastTo(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE())) {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	return pkg_base.StaticGet_Tri_FALSE()
}

// Java method: SymbolType.matchesNumberContext()
func (this_ *SymbolTypeBase) MatchesNumberContext() bool {
	return false
}

// Java method: SymbolType.matchesStringContext()
func (this_ *SymbolTypeBase) MatchesStringContext() bool {
	return false
}

// Java method: SymbolType.matchesSymbolContext()
func (this_ *SymbolTypeBase) MatchesSymbolContext() bool {
	return true
}

// Java method: SymbolType.matchesObjectContext()
func (this_ *SymbolTypeBase) MatchesObjectContext() bool {
	return true
}

// Java method: SymbolType.getDisplayName()
func (this_ *SymbolTypeBase) GetDisplayName() string {
	return "symbol"
}

// Java method: SymbolType.autoboxesTo()
func (this_ *SymbolTypeBase) AutoboxesTo() JSType {
	return this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE())
}

// Java method: SymbolType.getPossibleToBooleanOutcomes()
func (this_ *SymbolTypeBase) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_TRUE()
}

// Java method: SymbolType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *SymbolTypeBase) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseSymbolType(this_.self_SymbolType_())
}
