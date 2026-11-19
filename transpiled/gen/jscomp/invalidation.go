/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/disambiguate/Invalidation.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.disambiguate.Invalidation
type Invalidation interface {
	Implements_jscomp_Invalidation()
	IsJavaNil_() bool
	GetInvalidationBase() *InvalidationBase
}

func ToIface_Invalidation[T interface {
	Invalidation
	comparable
}](v T) Invalidation {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Invalidation[T Invalidation](v T) Invalidation {
	var vi Invalidation = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.disambiguate.Invalidation
type InvalidationBase struct {
	idHash_ atomic.Int32
	reason  Invalidation_Reason `gson:"reason"`
}

func (this_ *InvalidationBase) BindSelf_(s any) {
	_ = s
}
func (this_ *InvalidationBase) IsJavaNil_() bool { return this_ == nil }
func (*InvalidationBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.Invalidation", "Invalidation"
}
func (this_ *InvalidationBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.Invalidation", "com.google.javascript.jscomp.disambiguate.Invalidation", func(v_ any) bool {
		return j2gort.InstanceOf[Invalidation](v_)
	})
}
func (b *InvalidationBase) GetInvalidationBase() *InvalidationBase { return b }
func (*InvalidationBase) Implements_jscomp_Invalidation()          {}

var static_Invalidation_WELL_KNOWN_PROPERTY Invalidation

func StaticGet_Invalidation_WELL_KNOWN_PROPERTY() Invalidation {
	clinit_Invalidation()
	return static_Invalidation_WELL_KNOWN_PROPERTY
}
func StaticSet_Invalidation_WELL_KNOWN_PROPERTY(v Invalidation) {
	clinit_Invalidation()
	static_Invalidation_WELL_KNOWN_PROPERTY = v
}

var clinitDone_Invalidation uint32
var clinitState_Invalidation j2gort.ClinitState

func clinit_Invalidation() {
	if atomic.LoadUint32(&clinitDone_Invalidation) != 0 {
		return
	}
	clinitSlow_Invalidation()
}

func clinitSlow_Invalidation() {
	if !j2gort.EnterClinit(&clinitState_Invalidation, &clinitDone_Invalidation) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Invalidation, &clinitDone_Invalidation)
	static_Invalidation_WELL_KNOWN_PROPERTY = NewInvalidation(StaticGet_Invalidation_Reason_WELL_KNOWN_PROPERTY())
}

func (this_ *InvalidationBase) initFields_Invalidation() {
	_ = this_
}

// Java constructor: Invalidation(com.google.javascript.jscomp.disambiguate.Invalidation.Reason reason)
func NewInvalidation(reason Invalidation_Reason) Invalidation {
	clinit_Invalidation()
	this_ := &InvalidationBase{}
	InitInvalidation(this_, reason)
	return this_
}

func InitInvalidation(this_ *InvalidationBase, reason Invalidation_Reason) {
	_ = reason
	this_.initFields_Invalidation()
	this_.reason = reason
}

// Java static method: Invalidation.wellKnownProperty()
func Invalidation_WellKnownProperty() Invalidation {
	clinit_Invalidation()
	return StaticGet_Invalidation_WELL_KNOWN_PROPERTY()
}

// Java static method: Invalidation.invalidatingType(int receiverType)
func Invalidation_InvalidatingType(receiverType int32) Invalidation {
	clinit_Invalidation()
	_ = receiverType
	return NewInvalidation_WithReceiverType(StaticGet_Invalidation_Reason_INVALIDATING_TYPE(), receiverType)
}

// Java static method: Invalidation.undeclaredAccess(int receiverType)
func Invalidation_UndeclaredAccess(receiverType int32) Invalidation {
	clinit_Invalidation()
	_ = receiverType
	return NewInvalidation_WithReceiverType(StaticGet_Invalidation_Reason_UNDECLARED_ACCESS(), receiverType)
}

// Java enum: com.google.javascript.jscomp.disambiguate.Invalidation.Reason
type Invalidation_Reason uint8

const (
	enum_Invalidation_Reason_WELL_KNOWN_PROPERTY Invalidation_Reason = 1
	enum_Invalidation_Reason_INVALIDATING_TYPE   Invalidation_Reason = 2
	enum_Invalidation_Reason_UNDECLARED_ACCESS   Invalidation_Reason = 3
)

var enumBaseTable_Invalidation_Reason = [...]j2gort.EnumBase{
	{Name_: "WELL_KNOWN_PROPERTY", Ordinal_: 0},
	{Name_: "INVALIDATING_TYPE", Ordinal_: 1},
	{Name_: "UNDECLARED_ACCESS", Ordinal_: 2},
}
var enumValues_Invalidation_Reason = []Invalidation_Reason{1, 2, 3}

func StaticGet_Invalidation_Reason_WELL_KNOWN_PROPERTY() Invalidation_Reason {
	return enum_Invalidation_Reason_WELL_KNOWN_PROPERTY
}
func StaticGet_Invalidation_Reason_INVALIDATING_TYPE() Invalidation_Reason {
	return enum_Invalidation_Reason_INVALIDATING_TYPE
}
func StaticGet_Invalidation_Reason_UNDECLARED_ACCESS() Invalidation_Reason {
	return enum_Invalidation_Reason_UNDECLARED_ACCESS
}

func Invalidation_Reason_Values() []Invalidation_Reason {
	return j2gort.CloneSlice(enumValues_Invalidation_Reason)
}

func Invalidation_Reason_ValueOf(name string) Invalidation_Reason {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Invalidation_Reason {
		if enumBaseTable_Invalidation_Reason[i_].Name_ == name {
			return Invalidation_Reason(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.disambiguate.Invalidation.Reason." + name))
}

func (this_ Invalidation_Reason) IsJavaNil_() bool                              { return this_ == 0 }
func (Invalidation_Reason) HasIdentityEquals_()                                 {}
func (Invalidation_Reason) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Invalidation_Reason) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Invalidation_Reason) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Invalidation_Reason) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.Invalidation$Reason", "Reason"
}
func (this_ Invalidation_Reason) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Invalidation_Reason[this_-1]
}
func (this_ Invalidation_Reason) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Invalidation_Reason) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Invalidation_Reason) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Invalidation_Reason[this_-1].Name_
}
func (this_ Invalidation_Reason) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Invalidation_Reason[this_-1].Name_
}
func (this_ Invalidation_Reason) Equals(other_ any) bool {
	o_, ok_ := other_.(Invalidation_Reason)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Invalidation_Reason) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Invalidation_Reason[this_-1].HashCode()
}
func (this_ Invalidation_Reason) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Invalidation_Reason[this_-1].CompareTo(other_)
}
func (this_ Invalidation_Reason) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Invalidation_Reason[this_-1].Name_)
}
func init() {
	vals_ := []any{Invalidation_Reason(1), Invalidation_Reason(2), Invalidation_Reason(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.disambiguate.Invalidation.Reason", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.disambiguate.Invalidation$Reason", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.Invalidation.Reason", "com.google.javascript.jscomp.disambiguate.Invalidation$Reason", func(v_ any) bool {
		return j2gort.InstanceOf[Invalidation_Reason](v_)
	})
}

// Java class: com.google.javascript.jscomp.disambiguate.Invalidation.WithReceiverType
type Invalidation_WithReceiverType struct {
	InvalidationBase
	receiverType int32 `gson:"receiverType"`
}

func (this_ *Invalidation_WithReceiverType) BindSelf_(s any) {
	_ = s
}
func (this_ *Invalidation_WithReceiverType) IsJavaNil_() bool { return this_ == nil }
func (*Invalidation_WithReceiverType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.Invalidation$WithReceiverType", "WithReceiverType"
}
func (this_ *Invalidation_WithReceiverType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.InvalidationBase.IdentityHash_()
}
func (*Invalidation_WithReceiverType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.Invalidation.WithReceiverType", "com.google.javascript.jscomp.disambiguate.Invalidation$WithReceiverType", func(v_ any) bool {
		return j2gort.InstanceOf[*Invalidation_WithReceiverType](v_)
	})
}
func (*Invalidation_WithReceiverType) Implements_jscomp_Invalidation() {}

func (this_ *Invalidation_WithReceiverType) initFields_Invalidation_WithReceiverType() {
	_ = this_
}

// Java constructor: WithReceiverType(com.google.javascript.jscomp.disambiguate.Invalidation.Reason reason, int receiverType)
func NewInvalidation_WithReceiverType(reason Invalidation_Reason, receiverType int32) *Invalidation_WithReceiverType {
	this_ := &Invalidation_WithReceiverType{}
	InitInvalidation_WithReceiverType(this_, reason, receiverType)
	return this_
}

func InitInvalidation_WithReceiverType(this_ *Invalidation_WithReceiverType, reason Invalidation_Reason, receiverType int32) {
	_ = reason
	_ = receiverType
	InitInvalidation(&this_.InvalidationBase, reason)
	this_.initFields_Invalidation_WithReceiverType()
	this_.receiverType = receiverType
}
