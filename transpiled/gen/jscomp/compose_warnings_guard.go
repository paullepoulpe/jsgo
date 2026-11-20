/*
 * Copyright 2008 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/ComposeWarningsGuard.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.ComposeWarningsGuard
type ComposeWarningsGuard struct {
	WarningsGuardBase
	orderOfAddition j2gort.OrderedMap[WarningsGuard, *j2gort.BoxInt] `gson:"orderOfAddition"`
	guardComparator j2gort.Ext_java_util_Comparator                  `gson:"guardComparator"`
	guards          *j2gort.Ext_java_util_TreeSet                    `gson:"guards"`
	numberOfAdds    int32                                            `gson:"numberOfAdds"`
	demoteErrors    bool                                             `gson:"demoteErrors"`
}

func (this_ *ComposeWarningsGuard) BindSelf_(s any) {
	_ = s
}
func (this_ *ComposeWarningsGuard) IsJavaNil_() bool { return this_ == nil }
func (*ComposeWarningsGuard) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ComposeWarningsGuard", "ComposeWarningsGuard"
}
func (this_ *ComposeWarningsGuard) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.WarningsGuardBase.IdentityHash_()
}
func (*ComposeWarningsGuard) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ComposeWarningsGuard", "com.google.javascript.jscomp.ComposeWarningsGuard", func(v_ any) bool {
		return j2gort.InstanceOf[*ComposeWarningsGuard](v_)
	})
}
func (*ComposeWarningsGuard) Implements_jscomp_WarningsGuard() {}

func StaticGet_ComposeWarningsGuard_SerialVersionUID() int64 { return int64(1) }
func (this_ *ComposeWarningsGuard) initFields_ComposeWarningsGuard() {
	_ = this_
	this_.orderOfAddition = j2gort.NewOrderedMap[WarningsGuard, *j2gort.BoxInt]()

	this_.guardComparator = NewComposeWarningsGuard_GuardComparator(&this_.orderOfAddition)

	this_.guards = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(this_.guardComparator)
}

// Java constructor: ComposeWarningsGuard(java.util.List<com.google.javascript.jscomp.WarningsGuard> guards)
func NewComposeWarningsGuard_1_Ext_java_util_List(guards j2gort.Ext_java_util_List) *ComposeWarningsGuard {
	this_ := &ComposeWarningsGuard{}
	InitComposeWarningsGuard_1_Ext_java_util_List(this_, guards)
	return this_
}

func InitComposeWarningsGuard_1_Ext_java_util_List(this_ *ComposeWarningsGuard, guards j2gort.Ext_java_util_List) {
	_ = guards
	InitWarningsGuard(&this_.WarningsGuardBase)
	this_.initFields_ComposeWarningsGuard()
	this_.priv_ComposeWarningsGuard_addGuards_2117(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](guards))
}

// Java constructor: ComposeWarningsGuard(com.google.javascript.jscomp.WarningsGuard[] guards)
func NewComposeWarningsGuard_1_WarningsGuard(guards []WarningsGuard) *ComposeWarningsGuard {
	this_ := &ComposeWarningsGuard{}
	InitComposeWarningsGuard_1_WarningsGuard(this_, guards)
	return this_
}

func InitComposeWarningsGuard_1_WarningsGuard(this_ *ComposeWarningsGuard, guards []WarningsGuard) {
	_ = guards
	var iifeArg_1 []any
iifeSwitch_1:
	switch {
	default:

		arr_1 := guards
		if arr_1 == nil {
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = el_3
		}
		iifeArg_1 = out_2
		break iifeSwitch_1
	}
	InitComposeWarningsGuard_1_Ext_java_util_List(this_, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(iifeArg_1)))
}

// Java method: ComposeWarningsGuard.addGuard(com.google.javascript.jscomp.WarningsGuard guard)
func (this_ *ComposeWarningsGuard) AddGuard(guard WarningsGuard) {
	_ = guard
	var composeGuard *ComposeWarningsGuard
	_ = composeGuard
	var ok_4 bool
	composeGuard, ok_4 = j2gort.TryCast[*ComposeWarningsGuard](guard)
	if ok_4 {
		if composeGuard.demoteErrors {
			this_.demoteErrors = composeGuard.demoteErrors
		}
		this_.priv_ComposeWarningsGuard_addGuards_2117(j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_NavigableSet](composeGuard.guards.DescendingSet())))
	} else {
		this_.numberOfAdds = this_.numberOfAdds + 1
		this_.orderOfAddition.Put(guard, j2gort.ValueOfInt(this_.numberOfAdds))
		this_.guards.Remove_1_any(guard)
		this_.guards.Add_1(guard)
	}
}

// Java method: ComposeWarningsGuard.addGuards(java.lang.Iterable<com.google.javascript.jscomp.WarningsGuard> guards)
func (this_ *ComposeWarningsGuard) priv_ComposeWarningsGuard_addGuards_2117(guards j2gort.Ext_java_lang_Iterable) {
	_ = guards
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(guards))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			guard := j2gort.AsIface[WarningsGuard](rawElem_6)
			_ = guard
			this_.AddGuard(guard)
		}
	}
}

// Java method: ComposeWarningsGuard.level(com.google.javascript.jscomp.JSError error)
func (this_ *ComposeWarningsGuard) Level(error_ *JSError) *CheckLevel {
	_ = error_
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(this_.guards))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			guard := j2gort.AsIface[WarningsGuard](rawElem_8)
			_ = guard
			var newLevel *CheckLevel = guard.Level(error_)
			_ = newLevel
			if newLevel != nil {
				if this_.demoteErrors && (newLevel == StaticGet_CheckLevel_ERROR()) {
					return StaticGet_CheckLevel_WARNING()
				}
				return newLevel
			}
		}
	}
	return nil
}

// Java method: ComposeWarningsGuard.mustRunChecks(com.google.javascript.jscomp.DiagnosticGroup group)
func (this_ *ComposeWarningsGuard) MustRunChecks(group *DiagnosticGroup) pkg_base.Tri {
	_ = group
	var enable bool = this_.priv_ComposeWarningsGuard_enables_2120(group)
	_ = enable
	var disable bool = this_.priv_ComposeWarningsGuard_disables_2119(group)
	_ = disable
	if !(!(enable) || !(disable)) {
		j2gort.CheckStateFmt(false, "%s applied to %s", this_, j2gort.PtrToAny(group))
	}
	if enable {
		return pkg_base.StaticGet_Tri_TRUE()
	} else if disable {
		return pkg_base.StaticGet_Tri_FALSE()
	} else {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
}

// Java method: ComposeWarningsGuard.disables(com.google.javascript.jscomp.DiagnosticGroup group)
func (this_ *ComposeWarningsGuard) priv_ComposeWarningsGuard_disables_2119(group *DiagnosticGroup) bool {
	_ = group
	{
		cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(group.GetTypes()))
	nextSingleton:
		for cur_9.HasNext() {
			rawElem_10 := cur_9.Next()
			type_ := j2gort.CastPtr[DiagnosticType](rawElem_10)
			_ = type_
			var singleton *DiagnosticGroup = DiagnosticGroup_ForType(type_)
			_ = singleton
			{
				cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(this_.guards))
				for cur_11.HasNext() {
					rawElem_12 := cur_11.Next()
					guard := j2gort.AsIface[WarningsGuard](rawElem_12)
					_ = guard
					switch (guard.MustRunChecks(singleton)).Ordinal() {
					case int32(2) /* Tri.TRUE */ :
						return false
					case int32(0) /* Tri.FALSE */ :
						continue nextSingleton
					case int32(1) /* Tri.UNKNOWN */ :
					}
				}
			}
			return false
		}
	}
	return true
}

// Java method: ComposeWarningsGuard.enables(com.google.javascript.jscomp.DiagnosticGroup group)
func (this_ *ComposeWarningsGuard) priv_ComposeWarningsGuard_enables_2120(group *DiagnosticGroup) bool {
	_ = group
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(this_.guards))
		for cur_13.HasNext() {
			rawElem_14 := cur_13.Next()
			guard := j2gort.AsIface[WarningsGuard](rawElem_14)
			_ = guard
			switch (guard.MustRunChecks(group)).Ordinal() {
			case int32(2) /* Tri.TRUE */ :
				return true
			case int32(0) /* Tri.FALSE */ :
				return false
			case int32(1) /* Tri.UNKNOWN */ :
			}
		}
	}
	return false
}

// Java method: ComposeWarningsGuard.getGuards()
func (this_ *ComposeWarningsGuard) GetGuards() j2gort.Ext_java_util_SortedSet {
	return j2gort.Ext_java_util_Collections_UnmodifiableSortedSet(j2gort.AsIface[j2gort.Ext_java_util_SortedSet](this_.guards))
}

// Java method: ComposeWarningsGuard.toString()
func (this_ *ComposeWarningsGuard) String() string {
	return j2gort.Ext_com_google_common_base_Joiner_On_1_string(", ").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.guards))
}

// Java class: com.google.javascript.jscomp.ComposeWarningsGuard.GuardComparator
type ComposeWarningsGuard_GuardComparator struct {
	orderOfAddition *j2gort.OrderedMap[WarningsGuard, *j2gort.BoxInt] `gson:"orderOfAddition"`
	idHash_         atomic.Int32
}

func (this_ *ComposeWarningsGuard_GuardComparator) BindSelf_(s any) {
	_ = s
}
func (this_ *ComposeWarningsGuard_GuardComparator) IsJavaNil_() bool { return this_ == nil }
func (*ComposeWarningsGuard_GuardComparator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ComposeWarningsGuard$GuardComparator", "GuardComparator"
}
func (this_ *ComposeWarningsGuard_GuardComparator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ComposeWarningsGuard_GuardComparator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ComposeWarningsGuard.GuardComparator", "com.google.javascript.jscomp.ComposeWarningsGuard$GuardComparator", func(v_ any) bool {
		return j2gort.InstanceOf[*ComposeWarningsGuard_GuardComparator](v_)
	})
}
func (*ComposeWarningsGuard_GuardComparator) Implements_j2gort_Ext_java_util_Comparator() {}

func StaticGet_ComposeWarningsGuard_GuardComparator_SerialVersionUID() int64 { return int64(1) }
func (this_ *ComposeWarningsGuard_GuardComparator) initFields_ComposeWarningsGuard_GuardComparator() {
	_ = this_
}

// Java constructor: GuardComparator(java.util.Map<com.google.javascript.jscomp.WarningsGuard,java.lang.Integer> orderOfAddition)
func NewComposeWarningsGuard_GuardComparator(orderOfAddition *j2gort.OrderedMap[WarningsGuard, *j2gort.BoxInt]) *ComposeWarningsGuard_GuardComparator {
	this_ := &ComposeWarningsGuard_GuardComparator{}
	InitComposeWarningsGuard_GuardComparator(this_, orderOfAddition)
	return this_
}

func InitComposeWarningsGuard_GuardComparator(this_ *ComposeWarningsGuard_GuardComparator, orderOfAddition *j2gort.OrderedMap[WarningsGuard, *j2gort.BoxInt]) {
	_ = orderOfAddition
	this_.initFields_ComposeWarningsGuard_GuardComparator()
	this_.orderOfAddition = orderOfAddition
}

// Java method: GuardComparator.compare(com.google.javascript.jscomp.WarningsGuard a, com.google.javascript.jscomp.WarningsGuard b)
func (this_ *ComposeWarningsGuard_GuardComparator) Compare(arg_15 any, arg_16 any) int32 {
	a := j2gort.AsIface[WarningsGuard](arg_15)
	_ = a
	b := j2gort.AsIface[WarningsGuard](arg_16)
	_ = b
	var priorityDiff int32 = (a.GetPriority() - b.GetPriority())
	_ = priorityDiff
	if priorityDiff != int32(0) {
		return priorityDiff
	}
	return (this_.orderOfAddition.Get(b).IntValue() - this_.orderOfAddition.Get(a).IntValue())
}

func (this_ *ComposeWarningsGuard_GuardComparator) Reversed() j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_Reversed(this_)
}

func (this_ *ComposeWarningsGuard_GuardComparator) ThenComparing_1_Ext_java_util_Comparator(p0_ j2gort.Ext_java_util_Comparator) j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_ThenComparing(this_, p0_)
}

func (this_ *ComposeWarningsGuard_GuardComparator) ThenComparing_1_Ext_java_util_function_Function(p0_ j2gort.Ext_java_util_function_Function) j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_ThenComparingKey(this_, p0_)
}

func (this_ *ComposeWarningsGuard_GuardComparator) ThenComparingInt(p0_ j2gort.Ext_java_util_function_ToIntFunction) j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_ThenComparingInt(this_, p0_)
}

func (this_ *ComposeWarningsGuard_GuardComparator) ThenComparingLong(p0_ j2gort.Ext_java_util_function_ToLongFunction) j2gort.Ext_java_util_Comparator {
	return j2gort.ExtDefault_Comparator_ThenComparingLong(this_, p0_)
}
