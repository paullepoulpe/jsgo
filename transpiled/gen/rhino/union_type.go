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
// Port of com/google/javascript/rhino/jstype/UnionType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.UnionType
type UnionType struct {
	JSTypeBase
	alt0                          JSType   `gson:"alt0"`
	alt1                          JSType   `gson:"alt1"`
	alt2                          JSType   `gson:"alt2"`
	restAlternates                []JSType `gson:"restAlternates"`
	alternateCount                int32    `gson:"alternateCount"`
	resolvedAlternatesAtLastBuild int32    `gson:"resolvedAlternatesAtLastBuild"`
	alternatesResolvedBeforeBuild bool     `gson:"alternatesResolvedBeforeBuild"`
	resolvedToFirstAlternate      bool     `gson:"resolvedToFirstAlternate"`
}

func (this_ *UnionType) BindSelf_(s any) {
	_ = s
	this_.JSTypeBase.BindSelf_(s)
}
func (this_ *UnionType) IsJavaNil_() bool { return this_ == nil }
func (*UnionType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.UnionType", "UnionType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.UnionType", "com.google.javascript.rhino.jstype.UnionType", func(v_ any) bool {
		return j2gort.InstanceOf[*UnionType](v_)
	})
}
func (*UnionType) Implements_rhino_JSType() {}

var static_UnionType_TYPE_CLASS JSTypeClass

func StaticGet_UnionType_TYPE_CLASS() JSTypeClass {
	clinit_UnionType()
	return static_UnionType_TYPE_CLASS
}
func StaticSet_UnionType_TYPE_CLASS(v JSTypeClass) {
	clinit_UnionType()
	static_UnionType_TYPE_CLASS = v
}
func StaticGet_UnionType_MAX_UNION_SIZE() int32 { return int32(30) }

var clinitDone_UnionType uint32
var clinitState_UnionType j2gort.ClinitState

func clinit_UnionType() {
	if atomic.LoadUint32(&clinitDone_UnionType) != 0 {
		return
	}
	clinitSlow_UnionType()
}

func clinitSlow_UnionType() {
	if !j2gort.EnterClinit(&clinitState_UnionType, &clinitDone_UnionType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_UnionType, &clinitDone_UnionType)
	static_UnionType_TYPE_CLASS = StaticGet_JSTypeClass_UNION()
}

func (this_ *UnionType) initFields_UnionType() {
	_ = this_

	this_.resolvedAlternatesAtLastBuild = int32(-1)
}

// Java constructor: UnionType(com.google.javascript.rhino.jstype.UnionType.Builder builder)
func NewUnionType(builder *UnionType_Builder) *UnionType {
	this_ := &UnionType{}
	this_.BindSelf_(this_)
	InitUnionType(this_, builder)
	return this_
}
func allocPooled_UnionType(slab_ *j2gort.PooledSlab[UnionType], builder *UnionType_Builder) *UnionType {
	clinit_UnionType()
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitUnionType(this_, builder)
	return this_
}

func InitUnionType(this_ *UnionType, builder *UnionType_Builder) {
	_ = builder
	InitJSType_3(&this_.JSTypeBase, builder.registry, nil, int32(64))
	this_.initFields_UnionType()
	this_.priv_UnionType_fillFromBuilder_9566(builder)
	this_.GetJSTypeBase().registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_UnionType_TYPE_CLASS())
}

// Java method: UnionType.getTypeClass()
func (this_ *UnionType) GetTypeClass() JSTypeClass {
	return StaticGet_UnionType_TYPE_CLASS()
}

// Java static method: UnionType.builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func UnionType_S_Builder(registry *JSTypeRegistry) *UnionType_Builder {
	_ = registry
	return NewUnionType_Builder(registry)
}

// Java static method: UnionType.createPooledBuilder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func UnionType_CreatePooledBuilder(registry *JSTypeRegistry) *UnionType_Builder {
	_ = registry
	return NewUnionType_Builder(registry)
}

// Java method: UnionType.getAlternates()
func (this_ *UnionType) GetAlternates() j2gort.Ext_com_google_common_collect_ImmutableList {
	this_.priv_UnionType_ensureAlternatesUpToDate_9564()
	switch this_.alternateCount {
	case int32(0):
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	case int32(1):
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.alt0)
	case int32(2):
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(this_.alt0, this_.alt1)
	case int32(3):
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_3(this_.alt0, this_.alt1, this_.alt2)
	default:
		var b *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, this_.alternateCount)
		_ = b
		b.AddVoid_(this_.alt0)
		b.AddVoid_(this_.alt1)
		b.AddVoid_(this_.alt2)
		var rest []JSType = this_.restAlternates
		_ = rest
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = (this_.alternateCount - int32(3))
			_ = n
			for ; i < n; i = i + 1 {
				b.AddVoid_(rest[i])
			}
		}
		return b.BuildImmList_()
	}
	panic("unreachable switch expression")
}

// Java method: UnionType.getAlternatesCount()
func (this_ *UnionType) GetAlternatesCount() int32 {
	this_.priv_UnionType_ensureAlternatesUpToDate_9564()
	return this_.alternateCount
}

// Java method: UnionType.getAlternateAt(int index)
func (this_ *UnionType) GetAlternateAt(index int32) JSType {
	_ = index
	switch index {
	case int32(0):
		return this_.alt0
	case int32(1):
		return this_.alt1
	case int32(2):
		return this_.alt2
	default:
		return this_.restAlternates[(index - int32(3))]
	}
	panic("unreachable switch expression")
}

// Java method: UnionType.ensureAlternatesUpToDate()
func (this_ *UnionType) priv_UnionType_ensureAlternatesUpToDate_9564() {
	if !(this_.IsResolved()) && !(this_.alternatesResolvedBeforeBuild) {
		var resolvedAlternates int32 = this_.priv_UnionType_countResolvedAlternates_9565()
		_ = resolvedAlternates
		if resolvedAlternates != this_.resolvedAlternatesAtLastBuild {
			this_.resolvedAlternatesAtLastBuild = resolvedAlternates
			var n int32 = this_.alternateCount
			_ = n
			var a0 JSType = this_.alt0
			_ = a0
			var a1 JSType = this_.alt1
			_ = a1
			var a2 JSType = this_.alt2
			_ = a2
			var rest []JSType = this_.restAlternates
			_ = rest
			var b *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeRebuilder(this_)
			_ = b
			defer func() {
				this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(b)
			}()
			{
				if n >= int32(1) {
					b.AddAlternate(a0)
				}
				if n >= int32(2) {
					b.AddAlternate(a1)
				}
				if n >= int32(3) {
					b.AddAlternate(a2)
				}
				{
					var i int32 = int32(3)
					_ = i
					for ; i < n; i = i + 1 {
						b.AddAlternate(rest[(i - int32(3))])
					}
				}
				if !(this_.IsResolved()) && !(this_.alternatesResolvedBeforeBuild) {
					b.Build()
					this_.resolvedAlternatesAtLastBuild = this_.priv_UnionType_countResolvedAlternates_9565()
				}
			}
		}
	}
}

// Java method: UnionType.countResolvedAlternates()
func (this_ *UnionType) priv_UnionType_countResolvedAlternates_9565() int32 {
	var count int32 = int32(0)
	_ = count
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).GetJSTypeBase().IsResolved() {
				count = count + 1
			}
		}
	}
	return count
}

// Java method: UnionType.fillFromBuilder(com.google.javascript.rhino.jstype.UnionType.Builder builder)
func (this_ *UnionType) priv_UnionType_fillFromBuilder_9566(builder *UnionType_Builder) {
	_ = builder
	if !(!(this_.alternatesResolvedBeforeBuild)) {
		j2gort.CheckState(false)
	}
	var n int32 = builder.altCount
	_ = n
	if !(n > int32(0)) {
		j2gort.CheckState(false)
	}
	this_.alternateCount = n
	this_.alt0 = builder.alt0
	if n >= int32(2) {
		this_.alt1 = builder.alt1
	} else {
		this_.alt1 = nil
	}
	if n >= int32(3) {
		this_.alt2 = builder.alt2
	} else {
		this_.alt2 = nil
	}
	if n <= int32(3) {
		this_.restAlternates = nil
	} else {
		var rest []JSType = make([]JSType, (n - int32(3)))
		_ = rest
		{
			var i int32 = int32(3)
			_ = i
			for ; i < n; i = i + 1 {
				{
					a_1 := rest
					i_2 := (i - int32(3))
					a_1[i_2] = builder.priv_Builder_getAlternateAtInternal_9573(i)
				}
			}
		}
		this_.restAlternates = rest
	}
	this_.alternatesResolvedBeforeBuild = builder.alternatesResolvedBeforeBuild
	if this_.GetJSTypeBase().registry.IsResolverClosed() {
		var combined int8 = int8(32)
		_ = combined
		combined = int8((int32(combined) & int32(this_.alt0.GetJSTypeBase().jstypeFlags)))
		var wellFormedOk bool = (this_.alt0.GetJSTypeBase().GetWellFormedTemplatizedKind() == int32(0))
		_ = wellFormedOk
		if n >= int32(2) {
			combined = int8((int32(combined) & int32(this_.alt1.GetJSTypeBase().jstypeFlags)))
			{

				ba_3 := wellFormedOk
				bb_4 := (this_.alt1.GetJSTypeBase().GetWellFormedTemplatizedKind() == int32(0))
				wellFormedOk = ba_3 && bb_4
			}

		}
		if n >= int32(3) {
			combined = int8((int32(combined) & int32(this_.alt2.GetJSTypeBase().jstypeFlags)))
			{

				ba_5 := wellFormedOk
				bb_6 := (this_.alt2.GetJSTypeBase().GetWellFormedTemplatizedKind() == int32(0))
				wellFormedOk = ba_5 && bb_6
			}

		}
		if (n > int32(3)) && ((int32(combined) != int32(0)) || wellFormedOk) {
			var rest_2 []JSType = this_.restAlternates
			_ = rest_2
			{
				var i_2 int32 = int32(0)
				_ = i_2
				var m int32 = (n - int32(3))
				_ = m
				for ; i_2 < m; i_2 = i_2 + 1 {
					var r JSType = rest_2[i_2]
					_ = r
					combined = int8((int32(combined) & int32(r.GetJSTypeBase().jstypeFlags)))
					if wellFormedOk && (r.GetJSTypeBase().GetWellFormedTemplatizedKind() != int32(0)) {
						wellFormedOk = false
					}
				}
			}
		}
		this_.GetJSTypeBase().jstypeFlags = int8(((int32(this_.GetJSTypeBase().jstypeFlags) & ^(int32(int8(32)))) | int32(combined)))
		if wellFormedOk {
			this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) | int32(int8(-128))))
		} else {
			this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) & ^(int32(int8(-128)))))
		}
	} else {
		this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) & ^(int32(int8(32)))))
		this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) & ^(int32(int8(-128)))))
	}
}

// Java method: UnionType.matchesNumberContext()
func (this_ *UnionType) MatchesNumberContext() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).MatchesNumberContext() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.matchesStringContext()
func (this_ *UnionType) MatchesStringContext() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).MatchesStringContext() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.matchesSymbolContext()
func (this_ *UnionType) MatchesSymbolContext() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).MatchesSymbolContext() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.matchesObjectContext()
func (this_ *UnionType) MatchesObjectContext() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).MatchesObjectContext() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.findPropertyTypeWithoutConsideringTemplateTypes(java.lang.String propertyName)
func (this_ *UnionType) FindPropertyTypeWithoutConsideringTemplateTypes_1_string(propertyName string) JSType {
	_ = propertyName
	var propertyType JSType = nil
	_ = propertyType
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			var alternate JSType = this_.GetAlternateAt(i)
			_ = alternate
			if alternate.GetJSTypeBase().IsNullType() || alternate.GetJSTypeBase().IsVoidType() {
				continue
			}
			var altPropertyType JSType = alternate.GetJSTypeBase().FindPropertyType_1_string(propertyName)
			_ = altPropertyType
			if altPropertyType == nil || altPropertyType.IsJavaNil_() {
				continue
			}
			if propertyType == nil || propertyType.IsJavaNil_() {
				propertyType = altPropertyType
			} else {
				propertyType = propertyType.GetLeastSupertype(altPropertyType)
			}
		}
	}
	return propertyType
}

// Java method: UnionType.findPropertyTypeWithoutConsideringTemplateTypes(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *UnionType) FindPropertyTypeWithoutConsideringTemplateTypes_1_Property_Key(propertyName Property_Key) JSType {
	_ = propertyName
	var propertyType JSType = nil
	_ = propertyType
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			var alternate JSType = this_.GetAlternateAt(i)
			_ = alternate
			if alternate.GetJSTypeBase().IsNullType() || alternate.GetJSTypeBase().IsVoidType() {
				continue
			}
			var altPropertyType JSType = alternate.GetJSTypeBase().FindPropertyType_1_Property_Key(propertyName)
			_ = altPropertyType
			if altPropertyType == nil || altPropertyType.IsJavaNil_() {
				continue
			}
			if propertyType == nil || propertyType.IsJavaNil_() {
				propertyType = altPropertyType
			} else {
				propertyType = propertyType.GetLeastSupertype(altPropertyType)
			}
		}
	}
	return propertyType
}

// Java method: UnionType.canBeCalled()
func (this_ *UnionType) CanBeCalled() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if !(this_.GetAlternateAt(i).CanBeCalled()) {
				return false
			}
		}
	}
	return true
}

// Java method: UnionType.autobox()
func (this_ *UnionType) Autobox() JSType {
	var n int32 = this_.GetAlternatesCount()
	_ = n
	var builder *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(builder)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				builder.AddAlternate(this_.GetAlternateAt(i).Autobox())
			}
		}
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: UnionType.tryFastRestrict2(com.google.javascript.rhino.jstype.JSType a0, com.google.javascript.rhino.jstype.JSType a1, com.google.javascript.rhino.jstype.JSType r0, com.google.javascript.rhino.jstype.JSType r1)
func (this_ *UnionType) priv_UnionType_tryFastRestrict2_9567(a0 JSType, a1 JSType, r0 JSType, r1 JSType) JSType {
	_ = a0
	_ = a1
	_ = r0
	_ = r1
	if (r0 == a0 || ((r0 == nil || r0.IsJavaNil_()) && (a0 == nil || a0.IsJavaNil_()))) && (r1 == a1 || ((r1 == nil || r1.IsJavaNil_()) && (a1 == nil || a1.IsJavaNil_()))) {
		return this_
	}
	var noType JSType = this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
	_ = noType
	if r0 == noType || ((r0 == nil || r0.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_())) {
		if r1 == noType || ((r1 == nil || r1.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_())) {
			return noType
		}
		if r1 == a1 || ((r1 == nil || r1.IsJavaNil_()) && (a1 == nil || a1.IsJavaNil_())) {
			return r1
		}
	} else if (r1 == noType || ((r1 == nil || r1.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_()))) && (r0 == a0 || ((r0 == nil || r0.IsJavaNil_()) && (a0 == nil || a0.IsJavaNil_()))) {
		return r0
	}
	return nil
}

// Java method: UnionType.tryFastRestrict3(com.google.javascript.rhino.jstype.JSType a0, com.google.javascript.rhino.jstype.JSType a1, com.google.javascript.rhino.jstype.JSType a2, com.google.javascript.rhino.jstype.JSType r0, com.google.javascript.rhino.jstype.JSType r1, com.google.javascript.rhino.jstype.JSType r2)
func (this_ *UnionType) priv_UnionType_tryFastRestrict3_9568(a0 JSType, a1 JSType, a2 JSType, r0 JSType, r1 JSType, r2 JSType) JSType {
	_ = a0
	_ = a1
	_ = a2
	_ = r0
	_ = r1
	_ = r2
	if ((r0 == a0 || ((r0 == nil || r0.IsJavaNil_()) && (a0 == nil || a0.IsJavaNil_()))) && (r1 == a1 || ((r1 == nil || r1.IsJavaNil_()) && (a1 == nil || a1.IsJavaNil_())))) && (r2 == a2 || ((r2 == nil || r2.IsJavaNil_()) && (a2 == nil || a2.IsJavaNil_()))) {
		return this_
	}
	var noType JSType = this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
	_ = noType
	if r0 == noType || ((r0 == nil || r0.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_())) {
		if r1 == noType || ((r1 == nil || r1.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_())) {
			if r2 == noType || ((r2 == nil || r2.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_())) {
				return noType
			}
			if r2 == a2 || ((r2 == nil || r2.IsJavaNil_()) && (a2 == nil || a2.IsJavaNil_())) {
				return r2
			}
		} else if r1 == a1 || ((r1 == nil || r1.IsJavaNil_()) && (a1 == nil || a1.IsJavaNil_())) {
			if r2 == noType || ((r2 == nil || r2.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_())) {
				return r1
			}
			if r2 == a2 || ((r2 == nil || r2.IsJavaNil_()) && (a2 == nil || a2.IsJavaNil_())) {
				return this_.GetJSTypeBase().registry.CreateUnionType_2(r1, r2)
			}
		}
	} else if r0 == a0 || ((r0 == nil || r0.IsJavaNil_()) && (a0 == nil || a0.IsJavaNil_())) {
		if r1 == noType || ((r1 == nil || r1.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_())) {
			if r2 == noType || ((r2 == nil || r2.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_())) {
				return r0
			}
			if r2 == a2 || ((r2 == nil || r2.IsJavaNil_()) && (a2 == nil || a2.IsJavaNil_())) {
				return this_.GetJSTypeBase().registry.CreateUnionType_2(r0, r2)
			}
		} else if (r1 == a1 || ((r1 == nil || r1.IsJavaNil_()) && (a1 == nil || a1.IsJavaNil_()))) && (r2 == noType || ((r2 == nil || r2.IsJavaNil_()) && (noType == nil || noType.IsJavaNil_()))) {
			return this_.GetJSTypeBase().registry.CreateUnionType_2(r0, r1)
		}
	}
	return nil
}

// Java method: UnionType.restrictByNotNullOrUndefined()
func (this_ *UnionType) RestrictByNotNullOrUndefined() JSType {
	var n int32 = this_.GetAlternatesCount()
	_ = n
	if (int32(this_.GetJSTypeBase().jstypeFlags) & int32(int8(8))) != int32(0) {
		if n == int32(2) {
			var a0 JSType = this_.alt0
			_ = a0
			var a1 JSType = this_.alt1
			_ = a1
			var fast JSType = this_.priv_UnionType_tryFastRestrict2_9567(a0, a1, a0.RestrictByNotNullOrUndefined(), a1.RestrictByNotNullOrUndefined())
			_ = fast
			if fast != nil && !fast.IsJavaNil_() {
				return fast
			}
		} else if n == int32(3) {
			var a0_2 JSType = this_.alt0
			_ = a0_2
			var a1_2 JSType = this_.alt1
			_ = a1_2
			var a2 JSType = this_.alt2
			_ = a2
			var fast_2 JSType = this_.priv_UnionType_tryFastRestrict3_9568(a0_2, a1_2, a2, a0_2.RestrictByNotNullOrUndefined(), a1_2.RestrictByNotNullOrUndefined(), a2.RestrictByNotNullOrUndefined())
			_ = fast_2
			if fast_2 != nil && !fast_2.IsJavaNil_() {
				return fast_2
			}
		}
	}
	var builder *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(builder)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				builder.AddAlternate(this_.GetAlternateAt(i).RestrictByNotNullOrUndefined())
			}
		}
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: UnionType.restrictByNotUndefined()
func (this_ *UnionType) RestrictByNotUndefined() JSType {
	var n int32 = this_.GetAlternatesCount()
	_ = n
	if (int32(this_.GetJSTypeBase().jstypeFlags) & int32(int8(8))) != int32(0) {
		if n == int32(2) {
			var a0 JSType = this_.alt0
			_ = a0
			var a1 JSType = this_.alt1
			_ = a1
			var fast JSType = this_.priv_UnionType_tryFastRestrict2_9567(a0, a1, a0.RestrictByNotUndefined(), a1.RestrictByNotUndefined())
			_ = fast
			if fast != nil && !fast.IsJavaNil_() {
				return fast
			}
		} else if n == int32(3) {
			var a0_2 JSType = this_.alt0
			_ = a0_2
			var a1_2 JSType = this_.alt1
			_ = a1_2
			var a2 JSType = this_.alt2
			_ = a2
			var fast_2 JSType = this_.priv_UnionType_tryFastRestrict3_9568(a0_2, a1_2, a2, a0_2.RestrictByNotUndefined(), a1_2.RestrictByNotUndefined(), a2.RestrictByNotUndefined())
			_ = fast_2
			if fast_2 != nil && !fast_2.IsJavaNil_() {
				return fast_2
			}
		}
	}
	var builder *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(builder)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				builder.AddAlternate(this_.GetAlternateAt(i).RestrictByNotUndefined())
			}
		}
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: UnionType.restrictByNotNull()
func (this_ *UnionType) RestrictByNotNull() JSType {
	var n int32 = this_.GetAlternatesCount()
	_ = n
	if (int32(this_.GetJSTypeBase().jstypeFlags) & int32(int8(8))) != int32(0) {
		if n == int32(2) {
			var a0 JSType = this_.alt0
			_ = a0
			var a1 JSType = this_.alt1
			_ = a1
			var fast JSType = this_.priv_UnionType_tryFastRestrict2_9567(a0, a1, a0.RestrictByNotNull(), a1.RestrictByNotNull())
			_ = fast
			if fast != nil && !fast.IsJavaNil_() {
				return fast
			}
		} else if n == int32(3) {
			var a0_2 JSType = this_.alt0
			_ = a0_2
			var a1_2 JSType = this_.alt1
			_ = a1_2
			var a2 JSType = this_.alt2
			_ = a2
			var fast_2 JSType = this_.priv_UnionType_tryFastRestrict3_9568(a0_2, a1_2, a2, a0_2.RestrictByNotNull(), a1_2.RestrictByNotNull(), a2.RestrictByNotNull())
			_ = fast_2
			if fast_2 != nil && !fast_2.IsJavaNil_() {
				return fast_2
			}
		}
	}
	var builder *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(builder)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				builder.AddAlternate(this_.GetAlternateAt(i).RestrictByNotNull())
			}
		}
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: UnionType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *UnionType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	var result pkg_base.Tri = pkg_base.Tri(0)
	_ = result
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			var t JSType = this_.GetAlternateAt(i)
			_ = t
			var test pkg_base.Tri = t.TestForEquality(that)
			_ = test
			if result == 0 {
				result = test
			} else if !(j2gort.CheckNotNullEnum(result) == test) {
				return pkg_base.StaticGet_Tri_UNKNOWN()
			}
		}
	}
	return result
}

// Java method: UnionType.isNullable()
func (this_ *UnionType) IsNullable() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).IsNullable() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.isVoidable()
func (this_ *UnionType) IsVoidable() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).IsVoidable() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.isExplicitlyVoidable()
func (this_ *UnionType) IsExplicitlyVoidable() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).IsExplicitlyVoidable() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.isUnknownTypeSlow()
func (this_ *UnionType) IsUnknownTypeSlow() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).GetJSTypeBase().IsUnknownType() {
				return true
			}
		}
	}
	if this_.GetJSTypeBase().registry.GetResolver().IsClosed() {
		this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) | int32(int8(32))))
	}
	return false
}

// Java method: UnionType.isStruct()
func (this_ *UnionType) IsStruct() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).IsStruct() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.isDict()
func (this_ *UnionType) IsDict() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).IsDict() {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.getLeastSupertype(com.google.javascript.rhino.jstype.JSType that)
func (this_ *UnionType) GetLeastSupertype(that JSType) JSType {
	_ = that
	base_that_ := that.GetJSTypeBase()
	if !(base_that_.IsUnknownType()) && !(base_that_.IsUnionType()) {
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = this_.alternateCount
			_ = n
			for ; i < n; i = i + 1 {
				var alternate JSType = this_.GetAlternateAt(i)
				_ = alternate
				if !(alternate.GetJSTypeBase().IsUnknownType()) && base_that_.IsSubtypeOf_1(alternate) {
					return this_
				}
			}
		}
	}
	return JSType_GetLeastSupertype(this_, that)
}

// Java static method: UnionType.getGreatestSubtype(com.google.javascript.rhino.jstype.UnionType union, com.google.javascript.rhino.jstype.JSType that)
func UnionType_GetGreatestSubtype(union *UnionType, that JSType) JSType {
	_ = union
	_ = that
	var registry *JSTypeRegistry = union.GetJSTypeBase().registry
	_ = registry
	var builder *UnionType_Builder = registry.AcquireUnionTypeBuilder()
	_ = builder
	var result JSType = nil
	_ = result
	func() {
		defer func() {
			registry.ReleaseUnionTypeBuilder(builder)
		}()
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.alternateCount
			_ = n
			for ; i < n; i = i + 1 {
				var alternate JSType = union.GetAlternateAt(i)
				_ = alternate
				if alternate.GetJSTypeBase().IsSubtypeOf_1(that) {
					builder.AddAlternate(alternate)
				}
			}
		}
		if that.GetJSTypeBase().IsUnionType() {
			var thatUnion *UnionType = that.ToMaybeUnionType()
			_ = thatUnion
			var n_2 int32 = thatUnion.GetAlternatesCount()
			_ = n_2
			{
				var i_2 int32 = int32(0)
				_ = i_2
				for ; i_2 < n_2; i_2 = i_2 + 1 {
					var otherAlternate JSType = thatUnion.GetAlternateAt(i_2)
					_ = otherAlternate
					if otherAlternate.GetJSTypeBase().IsSubtypeOf_1(ToIface_JSType(union)) {
						builder.AddAlternate(otherAlternate)
					}
				}
			}
		} else if that.GetJSTypeBase().IsSubtypeOf_1(ToIface_JSType(union)) {
			builder.AddAlternate(that)
		}
		result = builder.Build()
	}()
	if !(result.GetJSTypeBase().IsNoType()) {
		return result
	} else if union.IsObject() && (that.IsObject() && !(that.GetJSTypeBase().IsNoType())) {
		return registry.GetNativeType(StaticGet_JSTypeNative_NO_OBJECT_TYPE())
	} else {
		return registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
	}
}

// Java method: UnionType.getPropertyKind(java.lang.String pname, boolean autobox)
func (this_ *UnionType) GetPropertyKind_2_string_bool(pname string, autobox bool) JSType_HasPropertyKind {
	_ = pname
	_ = autobox
	j2gort.CheckNotNullString(pname)
	var found bool = false
	_ = found
	var always bool = true
	_ = always
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			var alternate JSType = this_.GetAlternateAt(i)
			_ = alternate
			if alternate.GetJSTypeBase().IsNullType() || alternate.GetJSTypeBase().IsVoidType() {
				continue
			}
			switch (alternate.GetPropertyKind_2_string_bool(pname, autobox)).Ordinal() {
			case int32(1) /* HasPropertyKind.KNOWN_PRESENT */ :
				found = true
			case int32(0) /* HasPropertyKind.ABSENT */ :
				always = false
			case int32(2) /* HasPropertyKind.MAYBE_PRESENT */ :
				found = true
				always = false
			}
			if found && !(always) {
				break
			}
		}
	}
	if found {
		if always {
			return StaticGet_JSType_HasPropertyKind_KNOWN_PRESENT()
		} else {
			return StaticGet_JSType_HasPropertyKind_MAYBE_PRESENT()
		}
	} else {
		return StaticGet_JSType_HasPropertyKind_ABSENT()
	}
}

// Java method: UnionType.getPropertyKind(com.google.javascript.rhino.jstype.Property.Key pname, boolean autobox)
func (this_ *UnionType) GetPropertyKind_2_Property_Key_bool(pname Property_Key, autobox bool) JSType_HasPropertyKind {
	_ = pname
	_ = autobox
	var found bool = false
	_ = found
	var always bool = true
	_ = always
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			var alternate JSType = this_.GetAlternateAt(i)
			_ = alternate
			if alternate.GetJSTypeBase().IsNullType() || alternate.GetJSTypeBase().IsVoidType() {
				continue
			}
			switch (alternate.GetPropertyKind_2_Property_Key_bool(pname, autobox)).Ordinal() {
			case int32(1) /* HasPropertyKind.KNOWN_PRESENT */ :
				found = true
			case int32(0) /* HasPropertyKind.ABSENT */ :
				always = false
			case int32(2) /* HasPropertyKind.MAYBE_PRESENT */ :
				found = true
				always = false
			}
			if found && !(always) {
				break
			}
		}
	}
	if found {
		if always {
			return StaticGet_JSType_HasPropertyKind_KNOWN_PRESENT()
		} else {
			return StaticGet_JSType_HasPropertyKind_MAYBE_PRESENT()
		}
	} else {
		return StaticGet_JSType_HasPropertyKind_ABSENT()
	}
}

// Java method: UnionType.recursionUnsafeHashCode()
func (this_ *UnionType) RecursionUnsafeHashCode() int32 {
	var n int32 = this_.alternateCount
	_ = n
	var hashCode int32 = n
	_ = hashCode
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			hashCode = (hashCode * this_.GetAlternateAt(i).GetJSTypeBase().HashCode())
		}
	}
	return hashCode
}

// Java method: UnionType.toMaybeUnionType()
func (this_ *UnionType) ToMaybeUnionType() *UnionType {
	return this_
}

// Java method: UnionType.isObject()
func (this_ *UnionType) IsObject() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if !(this_.GetAlternateAt(i).IsObject()) {
				return false
			}
		}
	}
	return true
}

// Java method: UnionType.contains(com.google.javascript.rhino.jstype.JSType type)
func (this_ *UnionType) Contains(type_ JSType) bool {
	_ = type_
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if type_.GetJSTypeBase().EqualsJSType(this_.GetAlternateAt(i)) {
				return true
			}
		}
	}
	return false
}

// Java method: UnionType.getRestrictedUnion(com.google.javascript.rhino.jstype.JSType type)
func (this_ *UnionType) GetRestrictedUnion(type_ JSType) JSType {
	_ = type_
	var restricted *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = restricted
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(restricted)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = this_.alternateCount
			_ = n
			for ; i < n; i = i + 1 {
				var t JSType = this_.GetAlternateAt(i)
				_ = t
				if t.GetJSTypeBase().IsUnknownType() || !(t.GetJSTypeBase().IsSubtypeOf_1(type_)) {
					restricted.AddAlternate(t)
				}
			}
		}
		return restricted.Build()
	}
	panic("unreachable")
}

// Java method: UnionType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *UnionType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	sb.Append_1_string("(")
	var n int32 = this_.alternateCount
	_ = n
	var names []string = j2gort.MakeNonNilSlice[string](n)
	_ = names
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			names = append(names, sb.CloneWithConfig().Append_1_JSType(this_.GetAlternateAt(i)).Build())
		}
	}
	j2gort.SliceSortNatural(names)
	sb.AppendAll(j2gort.SliceReadOnlyIterable(names), "|")
	sb.Append_1_string(")")
}

// Java method: UnionType.getRestrictedTypeGivenOutcome(com.google.javascript.rhino.Outcome outcome)
func (this_ *UnionType) GetRestrictedTypeGivenOutcome(outcome *Outcome) JSType {
	_ = outcome
	var n int32 = this_.GetAlternatesCount()
	_ = n
	var builder *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(builder)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				builder.AddAlternate(this_.GetAlternateAt(i).GetRestrictedTypeGivenOutcome(outcome))
			}
		}
		return builder.Build()
	}
	panic("unreachable")
}

// Java method: UnionType.getPossibleToBooleanOutcomes()
func (this_ *UnionType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	var literals BooleanLiteralSet = StaticGet_BooleanLiteralSet_EMPTY()
	_ = literals
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			var element JSType = this_.GetAlternateAt(i)
			_ = element
			literals = literals.Union(element.GetPossibleToBooleanOutcomes())
			if literals == StaticGet_BooleanLiteralSet_BOTH() {
				break
			}
		}
	}
	return literals
}

// Java method: UnionType.getTypesUnderEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *UnionType) GetTypesUnderEquality(that JSType) JSType_TypePair {
	_ = that
	var thisRestricted *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = thisRestricted
	var thatRestricted *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = thatRestricted
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(thatRestricted)
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(thisRestricted)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = this_.alternateCount
			_ = n
			for ; i < n; i = i + 1 {
				var element JSType = this_.GetAlternateAt(i)
				_ = element
				var p JSType_TypePair = element.GetTypesUnderEquality(that)
				_ = p
				if p.typeA != nil && !p.typeA.IsJavaNil_() {
					thisRestricted.AddAlternate(p.typeA)
				}
				if p.typeB != nil && !p.typeB.IsJavaNil_() {
					thatRestricted.AddAlternate(p.typeB)
				}
			}
		}
		return NewJSType_TypePair(thisRestricted.Build(), thatRestricted.Build())
	}
	panic("unreachable")
}

// Java method: UnionType.getTypesUnderInequality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *UnionType) GetTypesUnderInequality(that JSType) JSType_TypePair {
	_ = that
	var thisRestricted *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = thisRestricted
	var thatRestricted *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = thatRestricted
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(thatRestricted)
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(thisRestricted)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = this_.alternateCount
			_ = n
			for ; i < n; i = i + 1 {
				var element JSType = this_.GetAlternateAt(i)
				_ = element
				var p JSType_TypePair = element.GetTypesUnderInequality(that)
				_ = p
				if p.typeA != nil && !p.typeA.IsJavaNil_() {
					thisRestricted.AddAlternate(p.typeA)
				}
				if p.typeB != nil && !p.typeB.IsJavaNil_() {
					thatRestricted.AddAlternate(p.typeB)
				}
			}
		}
		return NewJSType_TypePair(thisRestricted.Build(), thatRestricted.Build())
	}
	panic("unreachable")
}

// Java method: UnionType.getTypesUnderShallowInequality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *UnionType) GetTypesUnderShallowInequality(that JSType) JSType_TypePair {
	_ = that
	var thisRestricted *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = thisRestricted
	var thatRestricted *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
	_ = thatRestricted
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(thatRestricted)
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(thisRestricted)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = this_.alternateCount
			_ = n
			for ; i < n; i = i + 1 {
				var element JSType = this_.GetAlternateAt(i)
				_ = element
				var p JSType_TypePair = element.GetTypesUnderShallowInequality(that)
				_ = p
				if p.typeA != nil && !p.typeA.IsJavaNil_() {
					thisRestricted.AddAlternate(p.typeA)
				}
				if p.typeB != nil && !p.typeB.IsJavaNil_() {
					thatRestricted.AddAlternate(p.typeB)
				}
			}
		}
		return NewJSType_TypePair(thisRestricted.Build(), thatRestricted.Build())
	}
	panic("unreachable")
}

// Java method: UnionType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *UnionType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseUnionType(this_)
}

// Java method: UnionType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *UnionType) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseUnionType(this_, that)
}

// Java method: UnionType.getCachedResolveResult()
func (this_ *UnionType) GetCachedResolveResult() JSType {
	if !(this_.IsResolved()) {
		return nil
	}
	if this_.resolvedToFirstAlternate {
		return this_.alt0
	} else {
		return this_
	}
}

// Java method: UnionType.setCachedResolveResult(com.google.javascript.rhino.jstype.JSType result)
func (this_ *UnionType) SetCachedResolveResult(result JSType) {
	_ = result
	j2gort.CheckNotNullIface(result)
	if result == this_ || ((result == nil || result.IsJavaNil_()) && (this_ == nil || this_.IsJavaNil_())) {
		this_.resolvedToFirstAlternate = false
	} else {
		if !((this_.alternateCount == int32(1)) && (this_.alt0 == result || ((this_.alt0 == nil || this_.alt0.IsJavaNil_()) && (result == nil || result.IsJavaNil_())))) {
			j2gort.CheckState(false)
		}
		this_.resolvedToFirstAlternate = true
	}
	this_.MarkResolved()
}

// Java method: UnionType.clearCachedResolveResult()
func (this_ *UnionType) ClearCachedResolveResult() {
	this_.resolvedToFirstAlternate = false
	this_.JSTypeBase.ClearCachedResolveResult()
}

// Java method: UnionType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *UnionType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	if this_.alternatesResolvedBeforeBuild {
		return this_
	}
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			this_.GetAlternateAt(i).GetJSTypeBase().Resolve(reporter)
		}
	}
	var n_2 int32 = this_.alternateCount
	_ = n_2
	var a0 JSType = this_.alt0
	_ = a0
	var a1 JSType = this_.alt1
	_ = a1
	var a2 JSType = this_.alt2
	_ = a2
	var rest []JSType = this_.restAlternates
	_ = rest
	var b *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeRebuilder(this_)
	_ = b
	defer func() {
		this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(b)
	}()
	{
		if n_2 >= int32(1) {
			b.AddAlternate(a0)
		}
		if n_2 >= int32(2) {
			b.AddAlternate(a1)
		}
		if n_2 >= int32(3) {
			b.AddAlternate(a2)
		}
		{
			var i_2 int32 = int32(3)
			_ = i_2
			for ; i_2 < n_2; i_2 = i_2 + 1 {
				b.AddAlternate(rest[(i_2 - int32(3))])
			}
		}
		return b.Build()
	}
	panic("unreachable")
}

// Java method: UnionType.collapseUnion()
func (this_ *UnionType) CollapseUnion() JSType {
	var currentValue JSType = nil
	_ = currentValue
	var currentCommonSuper ObjectType = nil
	_ = currentCommonSuper
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			var a JSType = this_.GetAlternateAt(i)
			_ = a
			if a.GetJSTypeBase().IsUnknownType() {
				return this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
			}
			var obj ObjectType = a.GetJSTypeBase().ToObjectType()
			_ = obj
			if obj == nil || obj.IsJavaNil_() {
				if (currentValue == nil || currentValue.IsJavaNil_()) && (currentCommonSuper == nil || currentCommonSuper.IsJavaNil_()) {
					currentValue = a
				} else {
					return this_.GetNativeType(StaticGet_JSTypeNative_ALL_TYPE())
				}
			} else if currentValue != nil && !currentValue.IsJavaNil_() {
				return this_.GetNativeType(StaticGet_JSTypeNative_ALL_TYPE())
			} else if currentCommonSuper == nil || currentCommonSuper.IsJavaNil_() {
				currentCommonSuper = obj
			} else {
				currentCommonSuper = this_.GetJSTypeBase().registry.FindCommonSuperObject(currentCommonSuper, obj)
			}
		}
	}
	return Upcast_JSType(currentCommonSuper)
}

// Java method: UnionType.matchConstraint(com.google.javascript.rhino.jstype.JSType constraint)
func (this_ *UnionType) MatchConstraint(constraint JSType) {
	_ = constraint
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.alternateCount
		_ = n
		for ; i < n; i = i + 1 {
			var alternate JSType = this_.GetAlternateAt(i)
			_ = alternate
			alternate.MatchConstraint(constraint)
		}
	}
}

// Java method: UnionType.hasAnyTemplateTypesInternal()
func (this_ *UnionType) HasAnyTemplateTypesInternal() bool {
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			if this_.GetAlternateAt(i).GetJSTypeBase().HasAnyTemplateTypes() {
				return true
			}
		}
	}
	return false
}

// Java class: com.google.javascript.rhino.jstype.UnionType.Builder
type UnionType_Builder struct {
	rebuildTarget                 *UnionType                     `gson:"rebuildTarget"`
	registry                      *JSTypeRegistry                `gson:"registry"`
	alt0                          JSType                         `gson:"alt0"`
	alt1                          JSType                         `gson:"alt1"`
	alt2                          JSType                         `gson:"alt2"`
	alt3                          JSType                         `gson:"alt3"`
	alt4                          JSType                         `gson:"alt4"`
	alt5                          JSType                         `gson:"alt5"`
	restAlternates                j2gort.Ext_java_util_ArrayList `gson:"restAlternates"`
	altCount                      int32                          `gson:"altCount"`
	functionTypePosition          int32                          `gson:"functionTypePosition"`
	idHash_                       atomic.Int32
	hasBuilt                      bool `gson:"hasBuilt"`
	containsVoidType              bool `gson:"containsVoidType"`
	isAllType                     bool `gson:"isAllType"`
	isNativeUnknownType           bool `gson:"isNativeUnknownType"`
	areAllUnknownsChecked         bool `gson:"areAllUnknownsChecked"`
	alternatesResolvedBeforeBuild bool `gson:"alternatesResolvedBeforeBuild"`
}

func (this_ *UnionType_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *UnionType_Builder) IsJavaNil_() bool { return this_ == nil }
func (*UnionType_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.UnionType$Builder", "Builder"
}
func (this_ *UnionType_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*UnionType_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.UnionType.Builder", "com.google.javascript.rhino.jstype.UnionType$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*UnionType_Builder](v_)
	})
}

func StaticGet_UnionType_Builder_NON_SYMBOL_PRIMITIVE_MASK() int32      { return int32(1032192) }
func StaticGet_UnionType_Builder_PRIMITIVE_SUBTYPE_SPECIAL_MASK() int32 { return int32(339) }
func (this_ *UnionType_Builder) initFields_UnionType_Builder() {
	_ = this_

	this_.areAllUnknownsChecked = true
	this_.alternatesResolvedBeforeBuild = true
	this_.functionTypePosition = int32(-1)
}

// Java constructor: Builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewUnionType_Builder(registry *JSTypeRegistry) *UnionType_Builder {
	this_ := &UnionType_Builder{}
	InitUnionType_Builder(this_, registry)
	return this_
}

func InitUnionType_Builder(this_ *UnionType_Builder, registry *JSTypeRegistry) {
	_ = registry
	this_.initFields_UnionType_Builder()
	this_.rebuildTarget = nil
	this_.registry = registry
}

// Java method: Builder.setRebuildTargetForPool(com.google.javascript.rhino.jstype.UnionType rebuildTarget)
func (this_ *UnionType_Builder) SetRebuildTargetForPool(rebuildTarget *UnionType) {
	_ = rebuildTarget
	this_.rebuildTarget = rebuildTarget
}

// Java method: Builder.resetForPool()
func (this_ *UnionType_Builder) ResetForPool() {
	this_.rebuildTarget = nil
	this_.alt0 = nil
	this_.alt1 = nil
	this_.alt2 = nil
	this_.alt3 = nil
	this_.alt4 = nil
	this_.alt5 = nil
	this_.altCount = int32(0)
	if this_.restAlternates != nil {
		this_.restAlternates.Clear()
	}
	this_.hasBuilt = false
	this_.containsVoidType = false
	this_.isAllType = false
	this_.isNativeUnknownType = false
	this_.areAllUnknownsChecked = true
	this_.alternatesResolvedBeforeBuild = true
	this_.functionTypePosition = int32(-1)
}

// Java method: Builder.getAlternateAtInternal(int index)
func (this_ *UnionType_Builder) priv_Builder_getAlternateAtInternal_9573(index int32) JSType {
	_ = index
	switch index {
	case int32(0):
		return this_.alt0
	case int32(1):
		return this_.alt1
	case int32(2):
		return this_.alt2
	case int32(3):
		return this_.alt3
	case int32(4):
		return this_.alt4
	case int32(5):
		return this_.alt5
	default:
		return j2gort.AsIface[JSType](this_.restAlternates.Get((index - int32(6))))
	}
	panic("unreachable")
}

// Java method: Builder.setAlternateAtInternal(int index, com.google.javascript.rhino.jstype.JSType value)
func (this_ *UnionType_Builder) priv_Builder_setAlternateAtInternal_9574(index int32, value JSType) {
	_ = index
	_ = value
	switch index {
	case int32(0):
		this_.alt0 = value
	case int32(1):
		this_.alt1 = value
	case int32(2):
		this_.alt2 = value
	case int32(3):
		this_.alt3 = value
	case int32(4):
		this_.alt4 = value
	case int32(5):
		this_.alt5 = value
	default:
		this_.restAlternates.Set((index - int32(6)), value)
	}
}

// Java method: Builder.appendAlternateInternal(com.google.javascript.rhino.jstype.JSType value)
func (this_ *UnionType_Builder) priv_Builder_appendAlternateInternal_9575(value JSType) {
	_ = value
	switch this_.altCount {
	case int32(0):
		this_.alt0 = value
	case int32(1):
		this_.alt1 = value
	case int32(2):
		this_.alt2 = value
	case int32(3):
		this_.alt3 = value
	case int32(4):
		this_.alt4 = value
	case int32(5):
		this_.alt5 = value
	default:
		if this_.restAlternates == nil {
			this_.restAlternates = j2gort.NewExt_java_util_ArrayList_1_int32(int32(8))
		}
		this_.restAlternates.Add_1(value)
	}
	this_.altCount = this_.altCount + 1
}

// Java method: Builder.removeAlternateAtInternal(int index)
func (this_ *UnionType_Builder) priv_Builder_removeAlternateAtInternal_9576(index int32) {
	_ = index
	if index >= int32(6) {
		this_.restAlternates.Remove_1_int32((index - int32(6)))
		this_.altCount = this_.altCount - 1
		return
	}
	if index <= int32(0) {
		this_.alt0 = this_.alt1
	}
	if index <= int32(1) {
		this_.alt1 = this_.alt2
	}
	if index <= int32(2) {
		this_.alt2 = this_.alt3
	}
	if index <= int32(3) {
		this_.alt3 = this_.alt4
	}
	if index <= int32(4) {
		this_.alt4 = this_.alt5
	}
	if this_.altCount > int32(6) {
		this_.alt5 = j2gort.AsIface[JSType](this_.restAlternates.Remove_1_int32(int32(0)))
	} else {
		this_.alt5 = nil
	}
	this_.altCount = this_.altCount - 1
}

// Java static method: Builder.isSubtype(com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.jstype.JSType leftType)
func UnionType_Builder_IsSubtype(rightType JSType, leftType JSType) bool {
	_ = rightType
	_ = leftType
	return rightType.GetJSTypeBase().IsSubtypeWithoutStructuralTyping(leftType)
}

// Java method: Builder.addAlternates(java.util.Collection<? extends com.google.javascript.rhino.jstype.JSType> c)
func (this_ *UnionType_Builder) AddAlternates_1_Ext_java_util_Collection(c j2gort.Ext_java_util_Collection) *UnionType_Builder {
	_ = c
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(c))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			type_ := j2gort.AsIface[JSType](rawElem_8)
			_ = type_
			this_.AddAlternate(type_)
		}
	}
	return this_
}

// Java method: Builder.addAlternates(java.util.List<? extends com.google.javascript.rhino.jstype.JSType> list)
func (this_ *UnionType_Builder) AddAlternates_1_Ext_java_util_List(list j2gort.Ext_java_util_List) *UnionType_Builder {
	_ = list
	{
		var i int32 = int32(0)
		_ = i
		for ; i < list.Size(); i = i + 1 {
			this_.AddAlternate(j2gort.AsIface[JSType](list.Get(i)))
		}
	}
	return this_
}

// Java method: Builder.addAlternate(com.google.javascript.rhino.jstype.JSType alternate)
func (this_ *UnionType_Builder) AddAlternate(alternate JSType) *UnionType_Builder {
	_ = alternate
	this_.priv_Builder_checkHasNotBuilt_9584()
	if alternate.GetJSTypeBase().IsUnionType() {
		var union *UnionType = alternate.ToMaybeUnionType()
		_ = union
		var n int32 = union.GetAlternatesCount()
		_ = n
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				this_.AddAlternate(union.GetAlternateAt(i))
			}
		}
		return this_
	}
	if this_.altCount > int32(30) {
		return this_
	}
	if !(alternate.GetJSTypeBase().IsResolved()) {
		this_.alternatesResolvedBeforeBuild = false
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var n_2 int32 = this_.altCount
			_ = n_2
			for ; i_2 < n_2; i_2 = i_2 + 1 {
				if j2gort.IfaceEquals[JSType](this_.priv_Builder_getAlternateAtInternal_9573(i_2), alternate) {
					return this_
				}
			}
		}
		this_.priv_Builder_appendAlternateInternal_9575(alternate)
		return this_
	}
	if alternate.GetJSTypeBase().IsNoType() {
		return this_
	}
	this_.isAllType = (this_.isAllType || alternate.GetJSTypeBase().IsAllType())
	this_.containsVoidType = (this_.containsVoidType || alternate.GetJSTypeBase().IsVoidType())
	var isAlternateUnknown bool = j2gort.InstanceOf[*UnknownType](alternate)
	_ = isAlternateUnknown
	this_.isNativeUnknownType = (this_.isNativeUnknownType || isAlternateUnknown)
	if isAlternateUnknown {
		this_.areAllUnknownsChecked = (this_.areAllUnknownsChecked && alternate.GetJSTypeBase().IsCheckedUnknownType())
	}
	if this_.isAllType || this_.isNativeUnknownType {
		return this_
	}
	if alternate.GetJSTypeBase().IsFunctionType() && (this_.functionTypePosition != int32(-1)) {
		var other *FunctionType = this_.priv_Builder_getAlternateAtInternal_9573(this_.functionTypePosition).ToMaybeFunctionType()
		_ = other
		var supremum *FunctionType = alternate.ToMaybeFunctionType().SupAndInfHelper(other, true)
		_ = supremum
		this_.priv_Builder_setAlternateAtInternal_9574(this_.functionTypePosition, ToIface_JSType(supremum))
		return this_
	}
	var altMask int32 = alternate.GetJSTypeBase().typeKindMask
	_ = altMask
	{
		var index int32 = int32(0)
		_ = index
		for ; index < this_.altCount; index = index + 1 {
			var removeCurrent bool = false
			_ = removeCurrent
			var current JSType = this_.priv_Builder_getAlternateAtInternal_9573(index)
			_ = current
			if alternate == current || ((alternate == nil || alternate.IsJavaNil_()) && (current == nil || current.IsJavaNil_())) {
				return this_
			}
			if !(current.GetJSTypeBase().IsResolved()) {
				continue
			}
			if ((alternate.GetJSTypeBase().IsUnknownType() || current.GetJSTypeBase().IsUnknownType()) || alternate.GetJSTypeBase().HasAnyTemplateTypes()) || current.GetJSTypeBase().HasAnyTemplateTypes() {
				if alternate.GetJSTypeBase().EqualsJSType(current) {
					return this_
				}
			} else if alternate.GetJSTypeBase().IsTemplatizedType() || current.GetJSTypeBase().IsTemplatizedType() {
				if !(current.GetJSTypeBase().IsTemplatizedType()) {
					if UnionType_Builder_IsSubtype(alternate, current) {
						return this_
					}
				} else if !(alternate.GetJSTypeBase().IsTemplatizedType()) {
					if UnionType_Builder_IsSubtype(current, alternate) {
						removeCurrent = true
					}
				} else {
					if !(current.GetJSTypeBase().IsTemplatizedType() && alternate.GetJSTypeBase().IsTemplatizedType()) {
						j2gort.CheckState(false)
					}
					var templatizedAlternate *TemplatizedType = alternate.ToMaybeTemplatizedType()
					_ = templatizedAlternate
					var templatizedCurrent *TemplatizedType = current.ToMaybeTemplatizedType()
					_ = templatizedCurrent
					if templatizedCurrent.WrapsSameRawType(ToIface_JSType(templatizedAlternate)) {
						if current.GetJSTypeBase().EqualsJSType(alternate) {
							return this_
						} else if !(templatizedCurrent.GetReferencedType().IsReadonlyArrayType()) {
							var rawType ObjectType = templatizedCurrent.GetReferencedObjTypeInternal()
							_ = rawType
							alternate = ToIface_JSType(this_.registry.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(rawType, j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()))
							altMask = alternate.GetJSTypeBase().typeKindMask
							removeCurrent = true
						}
					}
				}
			} else {
				var combinedMask int32 = (altMask | current.GetJSTypeBase().typeKindMask)
				_ = combinedMask
				if ((combinedMask & int32(1032192)) != int32(0)) && ((combinedMask & int32(339)) == int32(0)) {
					continue
				}
				if (((combinedMask & int32(1048576)) != int32(0)) && (((altMask & current.GetJSTypeBase().typeKindMask) & int32(1048576)) == int32(0))) && ((combinedMask & int32(339)) == int32(0)) {
					continue
				}
				if UnionType_Builder_IsSubtype(alternate, current) {
					this_.priv_Builder_mayRegisterDroppedProperties_9580(alternate, current)
					return this_
				} else if UnionType_Builder_IsSubtype(current, alternate) {
					this_.priv_Builder_mayRegisterDroppedProperties_9580(current, alternate)
					removeCurrent = true
				}
			}
			if removeCurrent {
				this_.priv_Builder_removeAlternateAtInternal_9576(index)
				if index == this_.functionTypePosition {
					this_.functionTypePosition = int32(-1)
				} else if index < this_.functionTypePosition {
					this_.functionTypePosition = this_.functionTypePosition - 1
				}
				index = index - 1
			}
		}
	}
	if alternate.GetJSTypeBase().IsFunctionType() {
		if !(this_.functionTypePosition == int32(-1)) {
			j2gort.CheckState(false)
		}
		this_.functionTypePosition = this_.altCount
	}
	this_.priv_Builder_appendAlternateInternal_9575(alternate)
	return this_
}

// Java method: Builder.mayRegisterDroppedProperties(com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *UnionType_Builder) priv_Builder_mayRegisterDroppedProperties_9580(subtype JSType, supertype JSType) {
	_ = subtype
	_ = supertype
	if subtype.GetJSTypeBase().IsRecordType() && supertype.GetJSTypeBase().IsRecordType() {
		this_.registry.RegisterDroppedPropertiesInUnion(subtype.ToMaybeRecordType(), supertype.ToMaybeRecordType())
	}
}

// Java method: Builder.build()
func (this_ *UnionType_Builder) Build() JSType {
	this_.priv_Builder_checkHasNotBuilt_9584()
	this_.hasBuilt = true
	if this_.rebuildTarget != nil {
		return this_.priv_Builder_buildForRebuildTarget_9582()
	}
	var wildcard JSType = this_.priv_Builder_getNativeWildcardType_9583()
	_ = wildcard
	if wildcard != nil && !wildcard.IsJavaNil_() {
		if this_.containsVoidType {
			this_.altCount = int32(2)
			this_.alt0 = wildcard
			this_.alt1 = this_.registry.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE())
			return ToIface_JSType(allocPooled_UnionType(&this_.registry.slab_UnionType, this_))
		}
		return wildcard
	}
	if this_.altCount == int32(0) {
		return this_.registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
	}
	if this_.altCount == int32(1) {
		return j2gort.CheckNotNullIface(this_.alt0)
	}
	if this_.altCount > int32(30) {
		return this_.registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
	return ToIface_JSType(allocPooled_UnionType(&this_.registry.

		// Java method: Builder.buildForRebuildTarget()
		slab_UnionType, this_))
}

func (this_ *UnionType_Builder) priv_Builder_buildForRebuildTarget_9582() JSType {
	var wildcard JSType = this_.priv_Builder_getNativeWildcardType_9583()
	_ = wildcard
	if wildcard != nil && !wildcard.IsJavaNil_() {
		if this_.containsVoidType {
			this_.altCount = int32(2)
			this_.alt0 = wildcard
			this_.alt1 = this_.registry.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE())
		} else {
			this_.altCount = int32(1)
			this_.alt0 = wildcard
		}
	} else if this_.altCount == int32(0) {
		this_.altCount = int32(1)
		this_.alt0 = this_.registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
	} else if this_.altCount > int32(30) {
		this_.altCount = int32(1)
		this_.alt0 = this_.registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
	this_.rebuildTarget.priv_UnionType_fillFromBuilder_9566(this_)
	if this_.altCount == int32(1) {
		return this_.alt0
	}
	return ToIface_JSType(this_.rebuildTarget)
}

// Java method: Builder.getNativeWildcardType()
func (this_ *UnionType_Builder) priv_Builder_getNativeWildcardType_9583() JSType {
	if this_.isAllType {
		return this_.registry.GetNativeType(StaticGet_JSTypeNative_ALL_TYPE())
	} else if this_.isNativeUnknownType {
		if this_.areAllUnknownsChecked {
			return this_.registry.GetNativeType(StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE())
		} else {
			return this_.registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}
	}
	return nil
}

// Java method: Builder.checkHasNotBuilt()
func (this_ *UnionType_Builder) priv_Builder_checkHasNotBuilt_9584() {
	if !(!(this_.hasBuilt)) {
		j2gort.CheckStateStr(false, "Cannot reuse a `UnionType.Builder` that has already built.")
	}
}
