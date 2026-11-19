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
// Port of com/google/javascript/rhino/jstype/JSTypeIterations.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.JSTypeIterations
type JSTypeIterations struct {
	idHash_ atomic.Int32
}

func (this_ *JSTypeIterations) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeIterations) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeIterations) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.JSTypeIterations", "JSTypeIterations" }
func (this_ *JSTypeIterations) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeIterations) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeIterations", "com.google.javascript.rhino.jstype.JSTypeIterations", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeIterations](v_)
	})
}

func (this_ *JSTypeIterations) initFields_JSTypeIterations() {
	_ = this_
}

// Java constructor: JSTypeIterations()
func NewJSTypeIterations() *JSTypeIterations {
	this_ := &JSTypeIterations{}
	InitJSTypeIterations(this_)
	return this_
}

func InitJSTypeIterations(this_ *JSTypeIterations) {
	this_.initFields_JSTypeIterations()
}

// Java static method: JSTypeIterations.anyTypeMatches(java.util.function.Predicate<? super com.google.javascript.rhino.jstype.JSType> predicate, com.google.common.collect.ImmutableList<? extends com.google.javascript.rhino.jstype.JSType> types)
func JSTypeIterations_AnyTypeMatches_2_Ext_java_util_function_Predicate_Ext_com_google_common_collect_ImmutableList(predicate j2gort.Ext_java_util_function_Predicate, types j2gort.Ext_com_google_common_collect_ImmutableList) bool {
	_ = predicate
	_ = types
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = types.Size()
		_ = n
		for ; (i < n); i = i + 1 {
			if (predicate.Test(j2gort.AsIface[JSType](types.Get(i)))) {
				return true
			}
		}
	}
	return false
}

// Java static method: JSTypeIterations.anyTypeMatches(java.util.function.Predicate<? super com.google.javascript.rhino.jstype.JSType> predicate, com.google.javascript.rhino.jstype.UnionType union)
func JSTypeIterations_AnyTypeMatches_2_Ext_java_util_function_Predicate_UnionType(predicate j2gort.Ext_java_util_function_Predicate, union *UnionType) bool {
	_ = predicate
	_ = union
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = union.GetAlternatesCount()
		_ = n
		for ; (i < n); i = i + 1 {
			if (predicate.Test(union.GetAlternateAt(i))) {
				return true
			}
		}
	}
	return false
}

// Java static method: JSTypeIterations.allTypesMatch(java.util.function.Predicate<? super com.google.javascript.rhino.jstype.JSType> predicate, com.google.common.collect.ImmutableList<? extends com.google.javascript.rhino.jstype.JSType> types)
func JSTypeIterations_AllTypesMatch_2_Ext_java_util_function_Predicate_Ext_com_google_common_collect_ImmutableList(predicate j2gort.Ext_java_util_function_Predicate, types j2gort.Ext_com_google_common_collect_ImmutableList) bool {
	_ = predicate
	_ = types
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = types.Size()
		_ = n
		for ; (i < n); i = i + 1 {
			if (!(predicate.Test(j2gort.AsIface[JSType](types.Get(i))))) {
				return false
			}
		}
	}
	return true
}

// Java static method: JSTypeIterations.allTypesMatch(java.util.function.Predicate<? super com.google.javascript.rhino.jstype.JSType> predicate, com.google.javascript.rhino.jstype.UnionType union)
func JSTypeIterations_AllTypesMatch_2_Ext_java_util_function_Predicate_UnionType(predicate j2gort.Ext_java_util_function_Predicate, union *UnionType) bool {
	_ = predicate
	_ = union
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = union.GetAlternatesCount()
		_ = n
		for ; (i < n); i = i + 1 {
			if (!(predicate.Test(union.GetAlternateAt(i)))) {
				return false
			}
		}
	}
	return true
}

// Java static method: JSTypeIterations.mapTypes(java.util.function.Function<? super com.google.javascript.rhino.jstype.JSType,T> mapper, com.google.common.collect.ImmutableList<? extends com.google.javascript.rhino.jstype.JSType> types)
func JSTypeIterations_MapTypes_2_Ext_java_util_function_Function_Ext_com_google_common_collect_ImmutableList(mapper j2gort.Ext_java_util_function_Function, types j2gort.Ext_com_google_common_collect_ImmutableList) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = mapper
	_ = types
	var n int32 = types.Size()
	_ = n
	if ((n == int32(0))) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, n)
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; (i < n); i = i + 1 {
			builder.AddVoid_(j2gort.AsIface[JSType](mapper.Apply(j2gort.AsIface[JSType](types.Get(i)))))
		}
	}
	return builder.BuildImmList_()
}

// Java static method: JSTypeIterations.mapTypes(java.util.function.Function<? super com.google.javascript.rhino.jstype.JSType,? extends com.google.javascript.rhino.jstype.JSType> mapper, com.google.javascript.rhino.jstype.UnionType union)
func JSTypeIterations_MapTypes_2_Ext_java_util_function_Function_UnionType(mapper j2gort.Ext_java_util_function_Function, union *UnionType) JSType {
	_ = mapper
	_ = union
	var n int32 = union.GetAlternatesCount()
	_ = n
	var registry *JSTypeRegistry = union.GetJSTypeBase().registry
	_ = registry
	var builder *UnionType_Builder = registry.AcquireUnionTypeBuilder()
	_ = builder
	defer func() {
		registry.ReleaseUnionTypeBuilder(builder)
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; (i < n); i = i + 1 {
				builder.AddAlternate(j2gort.AsIface[JSType](mapper.Apply(union.GetAlternateAt(i))))
			}
		}
		return builder.Build()
	}
	panic("unreachable")
}

