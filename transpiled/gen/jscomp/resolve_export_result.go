/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/modules/ResolveExportResult.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.modules.ResolveExportResult
type ResolveExportResult struct {
	binding *Binding `gson:"binding"`
	idHash_ atomic.Int32
	state ResolveExportResult_State `gson:"state"`
}

func (this_ *ResolveExportResult) BindSelf_(s any) {
	_ = s
}
func (this_ *ResolveExportResult) IsJavaNil_() bool { return this_ == nil }
func (*ResolveExportResult) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.modules.ResolveExportResult", "ResolveExportResult" }
func (this_ *ResolveExportResult) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ResolveExportResult) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ResolveExportResult", "com.google.javascript.jscomp.modules.ResolveExportResult", func(v_ any) bool {
		return j2gort.InstanceOf[*ResolveExportResult](v_)
	})
}

var static_ResolveExportResult_AMBIGUOUS *ResolveExportResult
func StaticGet_ResolveExportResult_AMBIGUOUS() *ResolveExportResult {
	clinit_ResolveExportResult()
	return static_ResolveExportResult_AMBIGUOUS
}
func StaticSet_ResolveExportResult_AMBIGUOUS(v *ResolveExportResult) {
	clinit_ResolveExportResult()
	static_ResolveExportResult_AMBIGUOUS = v
}
var static_ResolveExportResult_NOT_FOUND *ResolveExportResult
func StaticGet_ResolveExportResult_NOT_FOUND() *ResolveExportResult {
	clinit_ResolveExportResult()
	return static_ResolveExportResult_NOT_FOUND
}
func StaticSet_ResolveExportResult_NOT_FOUND(v *ResolveExportResult) {
	clinit_ResolveExportResult()
	static_ResolveExportResult_NOT_FOUND = v
}
var static_ResolveExportResult_ERROR *ResolveExportResult
func StaticGet_ResolveExportResult_ERROR() *ResolveExportResult {
	clinit_ResolveExportResult()
	return static_ResolveExportResult_ERROR
}
func StaticSet_ResolveExportResult_ERROR(v *ResolveExportResult) {
	clinit_ResolveExportResult()
	static_ResolveExportResult_ERROR = v
}
var clinitDone_ResolveExportResult uint32
var clinitState_ResolveExportResult j2gort.ClinitState
func clinit_ResolveExportResult() {
	if atomic.LoadUint32(&clinitDone_ResolveExportResult) != 0 {
		return
	}
	clinitSlow_ResolveExportResult()
}

func clinitSlow_ResolveExportResult() {
	if !j2gort.EnterClinit(&clinitState_ResolveExportResult, &clinitDone_ResolveExportResult) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ResolveExportResult, &clinitDone_ResolveExportResult)
	static_ResolveExportResult_AMBIGUOUS = NewResolveExportResult(nil, StaticGet_ResolveExportResult_State_AMBIGUOUS())
	static_ResolveExportResult_NOT_FOUND = NewResolveExportResult(nil, StaticGet_ResolveExportResult_State_NOT_FOUND())
	static_ResolveExportResult_ERROR = NewResolveExportResult(nil, StaticGet_ResolveExportResult_State_ERROR())
}

func (this_ *ResolveExportResult) initFields_ResolveExportResult() {
	_ = this_
}

// Java constructor: ResolveExportResult(com.google.javascript.jscomp.modules.@org.jspecify.annotations.Nullable Binding binding, com.google.javascript.jscomp.modules.ResolveExportResult.State state)
func NewResolveExportResult(binding *Binding, state ResolveExportResult_State) *ResolveExportResult {
	clinit_ResolveExportResult()
	this_ := &ResolveExportResult{}
	InitResolveExportResult(this_, binding, state)
	return this_
}

func InitResolveExportResult(this_ *ResolveExportResult, binding *Binding, state ResolveExportResult_State) {
	_ = binding
	_ = state
	this_.initFields_ResolveExportResult()
	this_.binding = binding
	this_.state = state
}

// Java method: ResolveExportResult.copy(com.google.javascript.rhino.Node sourceNode, com.google.javascript.jscomp.modules.Binding.CreatedBy createdBy)
func (this_ *ResolveExportResult) Copy(sourceNode pkg_rhino.Node, createdBy Binding_CreatedBy) *ResolveExportResult {
	_ = sourceNode
	_ = createdBy
	j2gort.CheckNotNullNode(sourceNode)
	if ((this_.binding == nil)) {
		return this_
	}
	return NewResolveExportResult(this_.binding.Copy(sourceNode, createdBy), this_.state)
}

// Java method: ResolveExportResult.hadError()
func (this_ *ResolveExportResult) HadError() bool {
	return (this_.state == StaticGet_ResolveExportResult_State_ERROR())
}

// Java method: ResolveExportResult.isAmbiguous()
func (this_ *ResolveExportResult) IsAmbiguous() bool {
	return (this_.state == StaticGet_ResolveExportResult_State_AMBIGUOUS())
}

// Java method: ResolveExportResult.resolved()
func (this_ *ResolveExportResult) Resolved() bool {
	return (this_.state == StaticGet_ResolveExportResult_State_RESOLVED())
}

// Java method: ResolveExportResult.found()
func (this_ *ResolveExportResult) Found() bool {
	return (this_.state != StaticGet_ResolveExportResult_State_NOT_FOUND())
}

// Java method: ResolveExportResult.getBinding()
func (this_ *ResolveExportResult) GetBinding() *Binding {
	return this_.binding
}

// Java static method: ResolveExportResult.of(com.google.javascript.jscomp.modules.Binding binding)
func ResolveExportResult_Of(binding *Binding) *ResolveExportResult {
	clinit_ResolveExportResult()
	_ = binding
	j2gort.CheckNotNullPtr(binding)
	return NewResolveExportResult(binding, StaticGet_ResolveExportResult_State_RESOLVED())
}

// Java enum: com.google.javascript.jscomp.modules.ResolveExportResult.State
type ResolveExportResult_State uint8

const (
	enum_ResolveExportResult_State_RESOLVED ResolveExportResult_State = 1
	enum_ResolveExportResult_State_AMBIGUOUS ResolveExportResult_State = 2
	enum_ResolveExportResult_State_NOT_FOUND ResolveExportResult_State = 3
	enum_ResolveExportResult_State_ERROR ResolveExportResult_State = 4
)

var enumBaseTable_ResolveExportResult_State = [...]j2gort.EnumBase{
	{Name_: "RESOLVED", Ordinal_: 0},
	{Name_: "AMBIGUOUS", Ordinal_: 1},
	{Name_: "NOT_FOUND", Ordinal_: 2},
	{Name_: "ERROR", Ordinal_: 3},
}
var enumValues_ResolveExportResult_State = []ResolveExportResult_State{1, 2, 3, 4}

func StaticGet_ResolveExportResult_State_RESOLVED() ResolveExportResult_State {
	return enum_ResolveExportResult_State_RESOLVED
}
func StaticGet_ResolveExportResult_State_AMBIGUOUS() ResolveExportResult_State {
	return enum_ResolveExportResult_State_AMBIGUOUS
}
func StaticGet_ResolveExportResult_State_NOT_FOUND() ResolveExportResult_State {
	return enum_ResolveExportResult_State_NOT_FOUND
}
func StaticGet_ResolveExportResult_State_ERROR() ResolveExportResult_State {
	return enum_ResolveExportResult_State_ERROR
}

func ResolveExportResult_State_Values() []ResolveExportResult_State {
	return j2gort.CloneSlice(enumValues_ResolveExportResult_State)
}

func ResolveExportResult_State_ValueOf(name string) ResolveExportResult_State {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ResolveExportResult_State {
		if enumBaseTable_ResolveExportResult_State[i_].Name_ == name {
			return ResolveExportResult_State(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.modules.ResolveExportResult.State." + name))
}

func (this_ ResolveExportResult_State) IsJavaNil_() bool { return this_ == 0 }
func (ResolveExportResult_State) HasIdentityEquals_() {}
func (ResolveExportResult_State) Implements_j2gort_Ext_java_lang_Enum() {}
func (ResolveExportResult_State) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ResolveExportResult_State) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (ResolveExportResult_State) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.modules.ResolveExportResult$State", "State" }
func (this_ ResolveExportResult_State) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ResolveExportResult_State[this_-1]
}
func (this_ ResolveExportResult_State) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ResolveExportResult_State) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ResolveExportResult_State) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ResolveExportResult_State[this_-1].Name_
}
func (this_ ResolveExportResult_State) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ResolveExportResult_State[this_-1].Name_
}
func (this_ ResolveExportResult_State) Equals(other_ any) bool {
	o_, ok_ := other_.(ResolveExportResult_State)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ResolveExportResult_State) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ResolveExportResult_State[this_-1].HashCode()
}
func (this_ ResolveExportResult_State) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ResolveExportResult_State[this_-1].CompareTo(other_)
}
func (this_ ResolveExportResult_State) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ResolveExportResult_State[this_-1].Name_)
}
func init() {
	vals_ := []any{ResolveExportResult_State(1), ResolveExportResult_State(2), ResolveExportResult_State(3), ResolveExportResult_State(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.modules.ResolveExportResult.State", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.modules.ResolveExportResult$State", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.modules.ResolveExportResult.State", "com.google.javascript.jscomp.modules.ResolveExportResult$State", func(v_ any) bool {
		return j2gort.InstanceOf[ResolveExportResult_State](v_)
	})
}

