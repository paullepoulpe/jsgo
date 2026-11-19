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
// Port of com/google/javascript/jscomp/DiagnosticGroup.java (generated mechanically by j2go).

package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_diagnostic_group_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_diagnostic_group_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return NewDiagnosticGroup_1_DiagnosticType(j2gort.CastPtr[DiagnosticType](mrp0_)) })
}

// Java class: com.google.javascript.jscomp.DiagnosticGroup
type DiagnosticGroup struct {
	types   j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"types"`
	name    string                                            `gson:"name"`
	idHash_ atomic.Int32
}

func (this_ *DiagnosticGroup) BindSelf_(s any) {
	_ = s
}
func (this_ *DiagnosticGroup) IsJavaNil_() bool { return this_ == nil }
func (*DiagnosticGroup) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DiagnosticGroup", "DiagnosticGroup"
}
func (this_ *DiagnosticGroup) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DiagnosticGroup) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DiagnosticGroup", "com.google.javascript.jscomp.DiagnosticGroup", func(v_ any) bool {
		return j2gort.InstanceOf[*DiagnosticGroup](v_)
	})
}

var staticSyncMu_DiagnosticGroup sync.Mutex

func StaticGet_DiagnosticGroup_SerialVersionUID() int64 { return int64(1) }

var static_DiagnosticGroup_singletons j2gort.Ext_java_util_Map

func StaticGet_DiagnosticGroup_Singletons() j2gort.Ext_java_util_Map {
	clinit_DiagnosticGroup()
	return static_DiagnosticGroup_singletons
}
func StaticSet_DiagnosticGroup_Singletons(v j2gort.Ext_java_util_Map) {
	clinit_DiagnosticGroup()
	static_DiagnosticGroup_singletons = v
}

var clinitDone_DiagnosticGroup uint32
var clinitState_DiagnosticGroup j2gort.ClinitState

func clinit_DiagnosticGroup() {
	if atomic.LoadUint32(&clinitDone_DiagnosticGroup) != 0 {
		return
	}
	clinitSlow_DiagnosticGroup()
}

func clinitSlow_DiagnosticGroup() {
	if !j2gort.EnterClinit(&clinitState_DiagnosticGroup, &clinitDone_DiagnosticGroup) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DiagnosticGroup, &clinitDone_DiagnosticGroup)
	static_DiagnosticGroup_singletons = j2gort.NewExt_java_util_LinkedHashMap()
}

func (this_ *DiagnosticGroup) initFields_DiagnosticGroup() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: DiagnosticGroup(java.lang.String name, com.google.javascript.jscomp.DiagnosticType[] types)
func NewDiagnosticGroup_2_string_DiagnosticType(name string, types []*DiagnosticType) *DiagnosticGroup {
	clinit_DiagnosticGroup()
	this_ := &DiagnosticGroup{}
	InitDiagnosticGroup_2_string_DiagnosticType(this_, name, types)
	return this_
}

func InitDiagnosticGroup_2_string_DiagnosticType(this_ *DiagnosticGroup, name string, types []*DiagnosticType) {
	_ = name
	_ = types
	this_.initFields_DiagnosticGroup()
	this_.name = name
	this_.types = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Arrays_AsList(func() []any {
		arr_1 := types
		if arr_1 == nil {
			return nil
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = j2gort.PtrToAny(el_3)
		}
		return out_2
	}())))
}

// Java constructor: DiagnosticGroup(com.google.javascript.jscomp.DiagnosticType[] types)
func NewDiagnosticGroup_1_DiagnosticType_2(types []*DiagnosticType) *DiagnosticGroup {
	clinit_DiagnosticGroup()
	this_ := &DiagnosticGroup{}
	InitDiagnosticGroup_1_DiagnosticType_2(this_, types)
	return this_
}

func InitDiagnosticGroup_1_DiagnosticType_2(this_ *DiagnosticGroup, types []*DiagnosticType) {
	_ = types
	InitDiagnosticGroup_2_string_DiagnosticType(this_, j2gort.NullString, types)
}

// Java constructor: DiagnosticGroup(com.google.javascript.jscomp.DiagnosticType type)
func NewDiagnosticGroup_1_DiagnosticType(type_ *DiagnosticType) *DiagnosticGroup {
	clinit_DiagnosticGroup()
	this_ := &DiagnosticGroup{}
	InitDiagnosticGroup_1_DiagnosticType(this_, type_)
	return this_
}

func InitDiagnosticGroup_1_DiagnosticType(this_ *DiagnosticGroup, type_ *DiagnosticType) {
	_ = type_
	this_.initFields_DiagnosticGroup()
	this_.name = j2gort.NullString
	this_.types = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_1(j2gort.PtrToAny(type_))
}

// Java constructor: DiagnosticGroup(com.google.javascript.jscomp.DiagnosticGroup[] groups)
func NewDiagnosticGroup_1_DiagnosticGroup(groups []*DiagnosticGroup) *DiagnosticGroup {
	clinit_DiagnosticGroup()
	this_ := &DiagnosticGroup{}
	InitDiagnosticGroup_1_DiagnosticGroup(this_, groups)
	return this_
}

func InitDiagnosticGroup_1_DiagnosticGroup(this_ *DiagnosticGroup, groups []*DiagnosticGroup) {
	_ = groups
	InitDiagnosticGroup_2_string_DiagnosticGroup(this_, j2gort.NullString, groups)
}

// Java constructor: DiagnosticGroup(java.lang.String name, com.google.javascript.jscomp.DiagnosticGroup[] groups)
func NewDiagnosticGroup_2_string_DiagnosticGroup(name string, groups []*DiagnosticGroup) *DiagnosticGroup {
	clinit_DiagnosticGroup()
	this_ := &DiagnosticGroup{}
	InitDiagnosticGroup_2_string_DiagnosticGroup(this_, name, groups)
	return this_
}

func InitDiagnosticGroup_2_string_DiagnosticGroup(this_ *DiagnosticGroup, name string, groups []*DiagnosticGroup) {
	_ = name
	_ = groups
	this_.initFields_DiagnosticGroup()
	var set j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = set
	for _, el_4 := range j2gort.CheckNotNullSlice(groups) {
		group := el_4
		_ = group
		set.AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](group.types))
	}
	this_.name = name
	this_.types = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](set.Build())
}

// Java static method: DiagnosticGroup.forType(com.google.javascript.jscomp.DiagnosticType type)
func DiagnosticGroup_ForType(type_ *DiagnosticType) *DiagnosticGroup {
	clinit_DiagnosticGroup()
	staticSyncMu_DiagnosticGroup.Lock()
	defer staticSyncMu_DiagnosticGroup.Unlock()
	_ = type_
	StaticGet_DiagnosticGroup_Singletons().ComputeIfAbsent(j2gort.PtrToAny(type_), lambdaSingleton_diagnostic_group_1)
	return j2gort.CastPtr[DiagnosticGroup](StaticGet_DiagnosticGroup_Singletons().Get(j2gort.PtrToAny(type_)))
}

// Java method: DiagnosticGroup.matches(com.google.javascript.jscomp.JSError error)
func (this_ *DiagnosticGroup) Matches_1_JSError(error_ *JSError) bool {
	_ = error_
	return this_.Matches_1_DiagnosticType(error_.Type())
}

// Java method: DiagnosticGroup.matches(com.google.javascript.jscomp.DiagnosticType type)
func (this_ *DiagnosticGroup) Matches_1_DiagnosticType(type_ *DiagnosticType) bool {
	_ = type_
	return this_.types.Contains(j2gort.PtrToAny(type_))
}

// Java method: DiagnosticGroup.getTypes()
func (this_ *DiagnosticGroup) GetTypes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.types
}

// Java method: DiagnosticGroup.getName()
func (this_ *DiagnosticGroup) GetName() string {
	return this_.name
}

// Java method: DiagnosticGroup.toString()
func (this_ *DiagnosticGroup) String() string {
	if this_.name == j2gort.NullString {
		return javalang.ObjectToString(j2gort.GetObjectClass(this_).GetName(), this_)
	} else {
		return ("DiagnosticGroup<" + j2gort.StrValString(this_.name) + ">")
	}
}
