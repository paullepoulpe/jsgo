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
 *   Annie Wang
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/Outcome.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.rhino.Outcome
type Outcome struct {
	j2gort.EnumBase
	fn_IsTruthy func(*Outcome) bool
	fn_IsNullish func(*Outcome) pkg_base.Tri
	fn_Not func(*Outcome) *Outcome
}

func (this_ *Outcome) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *Outcome) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*Outcome) Implements_j2gort_Ext_java_lang_Enum() {}
func (*Outcome) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *Outcome) initFields_Outcome() {
	_ = this_
}

func initEnumOutcome(this_ *Outcome, enumName_ string, enumOrd_ int32) *Outcome {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_Outcome()
	return this_
}

var enum_Outcome_TRUE *Outcome
func StaticGet_Outcome_TRUE() *Outcome {
	clinit_Outcome()
	return enum_Outcome_TRUE
}
var enum_Outcome_FALSE *Outcome
func StaticGet_Outcome_FALSE() *Outcome {
	clinit_Outcome()
	return enum_Outcome_FALSE
}
var enum_Outcome_FALSE_NOT_NULL *Outcome
func StaticGet_Outcome_FALSE_NOT_NULL() *Outcome {
	clinit_Outcome()
	return enum_Outcome_FALSE_NOT_NULL
}
var enum_Outcome_NULLISH *Outcome
func StaticGet_Outcome_NULLISH() *Outcome {
	clinit_Outcome()
	return enum_Outcome_NULLISH
}
var clinitDone_Outcome uint32
var clinitState_Outcome j2gort.ClinitState
func clinit_Outcome() {
	if atomic.LoadUint32(&clinitDone_Outcome) != 0 {
		return
	}
	clinitSlow_Outcome()
}

func clinitSlow_Outcome() {
	if !j2gort.EnterClinit(&clinitState_Outcome, &clinitDone_Outcome) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Outcome, &clinitDone_Outcome)
	enum_Outcome_TRUE = &Outcome{}
	enum_Outcome_TRUE.fn_IsTruthy = func(this_ *Outcome) bool {
	return true
}
	enum_Outcome_TRUE.fn_IsNullish = func(this_ *Outcome) pkg_base.Tri {
	return pkg_base.StaticGet_Tri_FALSE()
}
	enum_Outcome_TRUE.fn_Not = func(this_ *Outcome) *Outcome {
	return StaticGet_Outcome_FALSE()
}
	initEnumOutcome(enum_Outcome_TRUE, "TRUE", int32(0))
	enum_Outcome_FALSE = &Outcome{}
	enum_Outcome_FALSE.fn_IsTruthy = func(this_ *Outcome) bool {
	return false
}
	enum_Outcome_FALSE.fn_IsNullish = func(this_ *Outcome) pkg_base.Tri {
	return pkg_base.StaticGet_Tri_UNKNOWN()
}
	enum_Outcome_FALSE.fn_Not = func(this_ *Outcome) *Outcome {
	return StaticGet_Outcome_TRUE()
}
	initEnumOutcome(enum_Outcome_FALSE, "FALSE", int32(1))
	enum_Outcome_FALSE_NOT_NULL = &Outcome{}
	enum_Outcome_FALSE_NOT_NULL.fn_IsTruthy = func(this_ *Outcome) bool {
	return false
}
	enum_Outcome_FALSE_NOT_NULL.fn_IsNullish = func(this_ *Outcome) pkg_base.Tri {
	return pkg_base.StaticGet_Tri_FALSE()
}
	enum_Outcome_FALSE_NOT_NULL.fn_Not = func(this_ *Outcome) *Outcome {
	return StaticGet_Outcome_TRUE()
}
	initEnumOutcome(enum_Outcome_FALSE_NOT_NULL, "FALSE_NOT_NULL", int32(2))
	enum_Outcome_NULLISH = &Outcome{}
	enum_Outcome_NULLISH.fn_IsTruthy = func(this_ *Outcome) bool {
	return false
}
	enum_Outcome_NULLISH.fn_IsNullish = func(this_ *Outcome) pkg_base.Tri {
	return pkg_base.StaticGet_Tri_TRUE()
}
	enum_Outcome_NULLISH.fn_Not = func(this_ *Outcome) *Outcome {
	return StaticGet_Outcome_TRUE()
}
	initEnumOutcome(enum_Outcome_NULLISH, "NULLISH", int32(3))
	enumValues_Outcome = []*Outcome{enum_Outcome_TRUE, enum_Outcome_FALSE, enum_Outcome_FALSE_NOT_NULL, enum_Outcome_NULLISH}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Outcome", []any{enum_Outcome_TRUE, enum_Outcome_FALSE, enum_Outcome_FALSE_NOT_NULL, enum_Outcome_NULLISH})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.rhino.Outcome", "com.google.javascript.rhino.Outcome", func() {
		clinit_Outcome()
	})
}

var enumValues_Outcome []*Outcome
func Outcome_Values() []*Outcome {
	clinit_Outcome()
	return j2gort.CloneSlice(enumValues_Outcome)
}

func Outcome_ValueOf(name string) *Outcome {
	clinit_Outcome()
	for _, c := range enumValues_Outcome {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.Outcome." + name))
}

// Java method: Outcome.isTruthy()
func (this_ *Outcome) IsTruthy() bool {
	if this_.fn_IsTruthy != nil {
		return this_.fn_IsTruthy(this_)
	}
	return this_.baseEnum_IsTruthy()
}

func (this_ *Outcome) baseEnum_IsTruthy() bool {
	panic("abstract or native method invoked")
}

// Java method: Outcome.isNullish()
func (this_ *Outcome) IsNullish() pkg_base.Tri {
	if this_.fn_IsNullish != nil {
		return this_.fn_IsNullish(this_)
	}
	return this_.baseEnum_IsNullish()
}

func (this_ *Outcome) baseEnum_IsNullish() pkg_base.Tri {
	panic("abstract or native method invoked")
}

// Java method: Outcome.not()
func (this_ *Outcome) Not() *Outcome {
	if this_.fn_Not != nil {
		return this_.fn_Not(this_)
	}
	return this_.baseEnum_Not()
}

func (this_ *Outcome) baseEnum_Not() *Outcome {
	panic("abstract or native method invoked")
}

// Java static method: Outcome.forBoolean(boolean val)
func Outcome_ForBoolean(val bool) *Outcome {
	clinit_Outcome()
	_ = val
	if val {
		return StaticGet_Outcome_TRUE()
	} else {
		return StaticGet_Outcome_FALSE()
	}
}

func (this_ *Outcome) IsJavaNil_() bool { return this_ == nil }
func (*Outcome) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.Outcome", "Outcome" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Outcome", "com.google.javascript.rhino.Outcome", func(v_ any) bool {
		return j2gort.InstanceOf[*Outcome](v_)
	})
}
