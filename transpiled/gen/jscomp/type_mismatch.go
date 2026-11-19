/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TypeMismatch.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.TypeMismatch
type TypeMismatch struct {
	found    pkg_rhino.JSType `gson:"found"`
	required pkg_rhino.JSType `gson:"required"`
	location pkg_rhino.Node   `gson:"location"`
}

func (this_ *TypeMismatch) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeMismatch) IsJavaNil_() bool { return this_ == nil }
func (*TypeMismatch) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeMismatch", "TypeMismatch"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeMismatch", "com.google.javascript.jscomp.TypeMismatch", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeMismatch](v_)
	})
}
func (*TypeMismatch) Implements_j2gort_Ext_java_lang_Record() {}

var static_TypeMismatch_TEST_LOCATION pkg_rhino.Node

func StaticGet_TypeMismatch_TEST_LOCATION() pkg_rhino.Node {
	clinit_TypeMismatch()
	return static_TypeMismatch_TEST_LOCATION
}
func StaticSet_TypeMismatch_TEST_LOCATION(v pkg_rhino.Node) {
	clinit_TypeMismatch()
	static_TypeMismatch_TEST_LOCATION = v
}

var clinitDone_TypeMismatch uint32
var clinitState_TypeMismatch j2gort.ClinitState

func clinit_TypeMismatch() {
	if atomic.LoadUint32(&clinitDone_TypeMismatch) != 0 {
		return
	}
	clinitSlow_TypeMismatch()
}

func clinitSlow_TypeMismatch() {
	if !j2gort.EnterClinit(&clinitState_TypeMismatch, &clinitDone_TypeMismatch) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypeMismatch, &clinitDone_TypeMismatch)
	static_TypeMismatch_TEST_LOCATION = pkg_rhino.IR_Empty()
}

func (this_ *TypeMismatch) initFields_TypeMismatch() {
	_ = this_
}

// Java constructor: TypeMismatch(com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required, com.google.javascript.rhino.Node location)
func NewTypeMismatch(found pkg_rhino.JSType, required pkg_rhino.JSType, location pkg_rhino.Node) *TypeMismatch {
	clinit_TypeMismatch()
	this_ := &TypeMismatch{}
	InitTypeMismatch(this_, found, required, location)
	return this_
}

func InitTypeMismatch(this_ *TypeMismatch, found pkg_rhino.JSType, required pkg_rhino.JSType, location pkg_rhino.Node) {
	_ = found
	_ = required
	_ = location
	this_.initFields_TypeMismatch()
	j2gort.CheckNotNullIface(found)
	j2gort.CheckNotNullIface(required)
	j2gort.CheckNotNullNode(location)
	this_.found = found
	this_.required = required
	this_.location = location
}

// Java method: TypeMismatch.getFound()
func (this_ *TypeMismatch) GetFound() pkg_rhino.JSType {
	return this_.Found()
}

// Java method: TypeMismatch.getRequired()
func (this_ *TypeMismatch) GetRequired() pkg_rhino.JSType {
	return this_.Required()
}

// Java method: TypeMismatch.getLocation()
func (this_ *TypeMismatch) GetLocation() pkg_rhino.Node {
	return this_.Location()
}

// Java static method: TypeMismatch.create(com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required, com.google.javascript.rhino.Node location)
func TypeMismatch_Create(found pkg_rhino.JSType, required pkg_rhino.JSType, location pkg_rhino.Node) *TypeMismatch {
	clinit_TypeMismatch()
	_ = found
	_ = required
	_ = location
	return NewTypeMismatch(found, required, location)
}

// Java static method: TypeMismatch.createForTesting(com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required)
func TypeMismatch_CreateForTesting(found pkg_rhino.JSType, required pkg_rhino.JSType) *TypeMismatch {
	clinit_TypeMismatch()
	_ = found
	_ = required
	return TypeMismatch_Create(found, required, StaticGet_TypeMismatch_TEST_LOCATION())
}

func (this_ *TypeMismatch) Found() pkg_rhino.JSType    { return this_.found }
func (this_ *TypeMismatch) Required() pkg_rhino.JSType { return this_.required }
func (this_ *TypeMismatch) Location() pkg_rhino.Node   { return this_.location }
func (this_ *TypeMismatch) Equals(other_ any) bool {
	o_, ok_ := other_.(*TypeMismatch)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.found, o_.found) {
		return false
	}
	if !javalang.ObjectsEquals(this_.required, o_.required) {
		return false
	}
	if this_.location != o_.location {
		return false
	}
	return true
}
func (this_ *TypeMismatch) HashCode() int32 {
	var h_ int32 = 0
	{
		hi_1 := this_.found
		h_ = 31 * h_
		if hi_1 != nil && !hi_1.IsJavaNil_() {
			h_ += hi_1.GetJSTypeBase().HashCode()
		}
	}
	{
		hi_2 := this_.required
		h_ = 31 * h_
		if hi_2 != nil && !hi_2.IsJavaNil_() {
			h_ += hi_2.GetJSTypeBase().HashCode()
		}
	}
	h_ = 31*h_ + int32(this_.location)
	return h_
}
func (this_ *TypeMismatch) String() string {
	_ = fmt.Sprint
	return "TypeMismatch[" + "found=" + j2gort.StrVal(this_.found) + ", " + "required=" + j2gort.StrVal(this_.required) + ", " + "location=" + j2gort.StrVal(j2gort.BoxNode(this_.location)) + "]"
}

// Java class: com.google.javascript.jscomp.TypeMismatch.Accumulator
type TypeMismatch_Accumulator struct {
	mismatches []*TypeMismatch `gson:"mismatches"`
	idHash_    atomic.Int32
}

func (this_ *TypeMismatch_Accumulator) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeMismatch_Accumulator) IsJavaNil_() bool { return this_ == nil }
func (*TypeMismatch_Accumulator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeMismatch$Accumulator", "Accumulator"
}
func (this_ *TypeMismatch_Accumulator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeMismatch_Accumulator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeMismatch.Accumulator", "com.google.javascript.jscomp.TypeMismatch$Accumulator", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeMismatch_Accumulator](v_)
	})
}

func (this_ *TypeMismatch_Accumulator) initFields_TypeMismatch_Accumulator() {
	_ = this_
	this_.mismatches = []*TypeMismatch{}
}

// Java constructor: Accumulator()
func NewTypeMismatch_Accumulator() *TypeMismatch_Accumulator {
	this_ := &TypeMismatch_Accumulator{}
	InitTypeMismatch_Accumulator(this_)
	return this_
}

func InitTypeMismatch_Accumulator(this_ *TypeMismatch_Accumulator) {
	this_.initFields_TypeMismatch_Accumulator()
}

// Java method: Accumulator.registerMismatch(com.google.javascript.rhino.Node location, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required)
func (this_ *TypeMismatch_Accumulator) RegisterMismatch(location pkg_rhino.Node, found pkg_rhino.JSType, required pkg_rhino.JSType) {
	_ = location
	_ = found
	_ = required
	found = TypeMismatch_Accumulator_RemoveNullUndefinedAndTemplates(found)
	required = TypeMismatch_Accumulator_RemoveNullUndefinedAndTemplates(required)
	if found.GetJSTypeBase().IsSubtypeOf_1(required) || required.GetJSTypeBase().IsSubtypeOf_1(found) {
		return
	}
	if required.GetJSTypeBase().IsUnionType() {
		var union *pkg_rhino.UnionType = required.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				if union.GetAlternateAt(i).GetJSTypeBase().IsSubtypeOf_1(found) {
					return
				}
			}
		}
	}
	if TypeMismatch_Accumulator_BothAreNotTemplateTypes(found, required) {
		this_.mismatches = append(this_.mismatches, TypeMismatch_Create(found, required, location))
	}
	if found.GetJSTypeBase().IsFunctionType() && required.GetJSTypeBase().IsFunctionType() {
		var fnTypeA *pkg_rhino.FunctionType = found.ToMaybeFunctionType()
		_ = fnTypeA
		var fnTypeB *pkg_rhino.FunctionType = required.ToMaybeFunctionType()
		_ = fnTypeB
		var paramsA []pkg_rhino.FunctionType_Parameter = fnTypeA.GetParameterArray()
		_ = paramsA
		var paramsB []pkg_rhino.FunctionType_Parameter = fnTypeB.GetParameterArray()
		_ = paramsB
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var len_ int32 = min(int32(len(paramsA)), int32(len(paramsB)))
			_ = len_
			for ; i_2 < len_; i_2 = i_2 + 1 {
				this_.priv_Accumulator_registerIfMismatch_5583(location, paramsA[i_2].GetJSType(), paramsB[i_2].GetJSType())
			}
		}
		this_.priv_Accumulator_registerIfMismatch_5583(location, fnTypeA.GetReturnType(), fnTypeB.GetReturnType())
	}
}

// Java method: Accumulator.getMismatches()
func (this_ *TypeMismatch_Accumulator) GetMismatches() j2gort.Ext_com_google_common_collect_ImmutableCollection {
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](j2gort.ImmListFromSlice(this_.mismatches))
}

// Java method: Accumulator.registerIfMismatch(com.google.javascript.rhino.Node location, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required)
func (this_ *TypeMismatch_Accumulator) priv_Accumulator_registerIfMismatch_5583(location pkg_rhino.Node, found pkg_rhino.JSType, required pkg_rhino.JSType) {
	_ = location
	_ = found
	_ = required
	if ((found != nil && !found.IsJavaNil_()) && (required != nil && !required.IsJavaNil_())) && !(found.GetJSTypeBase().IsSubtypeOf_1(required)) {
		this_.RegisterMismatch(location, found, required)
	}
}

// Java static method: Accumulator.bothAreNotTemplateTypes(com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required)
func TypeMismatch_Accumulator_BothAreNotTemplateTypes(found pkg_rhino.JSType, required pkg_rhino.JSType) bool {
	_ = found
	_ = required
	return (!(found.GetJSTypeBase().IsTemplateType()) && !(required.GetJSTypeBase().IsTemplateType()))
}

// Java static method: Accumulator.removeNullUndefinedAndTemplates(com.google.javascript.rhino.jstype.JSType t)
func TypeMismatch_Accumulator_RemoveNullUndefinedAndTemplates(t pkg_rhino.JSType) pkg_rhino.JSType {
	_ = t
	var result pkg_rhino.JSType = t.RestrictByNotNullOrUndefined()
	_ = result
	var obj pkg_rhino.ObjectType = result.GetJSTypeBase().ToMaybeObjectType()
	_ = obj
	if (obj != nil && !obj.IsJavaNil_()) && obj.GetObjectTypeBase().IsTemplatizedType() {
		return pkg_rhino.Upcast_JSType(obj.ToMaybeTemplatizedType().GetRawType())
	}
	return result
}
