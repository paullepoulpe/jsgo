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
// Port of com/google/javascript/jscomp/base/JSCompDoubles.java (generated mechanically by j2go).

package base

import (
	"math"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.base.JSCompDoubles
type JSCompDoubles struct {
	idHash_ atomic.Int32
}

func (this_ *JSCompDoubles) BindSelf_(s any) {
	_ = s
}
func (this_ *JSCompDoubles) IsJavaNil_() bool { return this_ == nil }
func (*JSCompDoubles) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.base.JSCompDoubles", "JSCompDoubles"
}
func (this_ *JSCompDoubles) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSCompDoubles) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.JSCompDoubles", "com.google.javascript.jscomp.base.JSCompDoubles", func(v_ any) bool {
		return j2gort.InstanceOf[*JSCompDoubles](v_)
	})
}

var static_JSCompDoubles_POW_2_32 float64

func StaticGet_JSCompDoubles_POW_2_32() float64 {
	clinit_JSCompDoubles()
	return static_JSCompDoubles_POW_2_32
}
func StaticSet_JSCompDoubles_POW_2_32(v float64) {
	clinit_JSCompDoubles()
	static_JSCompDoubles_POW_2_32 = v
}

var static_JSCompDoubles_POW_2_53 float64

func StaticGet_JSCompDoubles_POW_2_53() float64 {
	clinit_JSCompDoubles()
	return static_JSCompDoubles_POW_2_53
}
func StaticSet_JSCompDoubles_POW_2_53(v float64) {
	clinit_JSCompDoubles()
	static_JSCompDoubles_POW_2_53 = v
}

var clinitDone_JSCompDoubles uint32
var clinitState_JSCompDoubles j2gort.ClinitState

func clinit_JSCompDoubles() {
	if atomic.LoadUint32(&clinitDone_JSCompDoubles) != 0 {
		return
	}
	clinitSlow_JSCompDoubles()
}

func clinitSlow_JSCompDoubles() {
	if !j2gort.EnterClinit(&clinitState_JSCompDoubles, &clinitDone_JSCompDoubles) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSCompDoubles, &clinitDone_JSCompDoubles)
	static_JSCompDoubles_POW_2_32 = math.Pow(float64(int32(2)), float64(int32(32)))
	static_JSCompDoubles_POW_2_53 = math.Pow(float64(int32(2)), float64(int32(53)))
}

func (this_ *JSCompDoubles) initFields_JSCompDoubles() {
	_ = this_
}

// Java constructor: JSCompDoubles()
func NewJSCompDoubles() *JSCompDoubles {
	clinit_JSCompDoubles()
	this_ := &JSCompDoubles{}
	InitJSCompDoubles(this_)
	return this_
}

func InitJSCompDoubles(this_ *JSCompDoubles) {
	this_.initFields_JSCompDoubles()
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java static method: JSCompDoubles.isExactInt32(double x)
func JSCompDoubles_IsExactInt32(x float64) bool {
	clinit_JSCompDoubles()
	_ = x
	return (!(math.IsNaN(x)) && ((float64((javalang.D2I(x)))) == x))
}

// Java static method: JSCompDoubles.isExactInt64(double x)
func JSCompDoubles_IsExactInt64(x float64) bool {
	clinit_JSCompDoubles()
	_ = x
	return (!(math.IsNaN(x)) && ((float64((javalang.D2L(x)))) == x))
}

// Java static method: JSCompDoubles.isMathematicalInteger(double x)
func JSCompDoubles_IsMathematicalInteger(x float64) bool {
	clinit_JSCompDoubles()
	_ = x
	return (math.Mod(x, float64(1.0)) == float64(0.0))
}

// Java static method: JSCompDoubles.isAtLeastIntegerPrecision(double x)
func JSCompDoubles_IsAtLeastIntegerPrecision(x float64) bool {
	clinit_JSCompDoubles()
	_ = x
	return (math.Abs(x) < StaticGet_JSCompDoubles_POW_2_53())
}

// Java static method: JSCompDoubles.isNegative(double x)
func JSCompDoubles_IsNegative(x float64) bool {
	clinit_JSCompDoubles()
	_ = x
	if !(!(math.IsNaN(x))) {
		j2gort.CheckState(false)
	}
	return (j2gort.BoxDouble_Compare(float64(0.0), x) > int32(0))
}

// Java static method: JSCompDoubles.isPositive(double x)
func JSCompDoubles_IsPositive(x float64) bool {
	clinit_JSCompDoubles()
	_ = x
	return !(JSCompDoubles_IsNegative(x))
}

// Java static method: JSCompDoubles.isEitherZero(double x)
func JSCompDoubles_IsEitherZero(x float64) bool {
	clinit_JSCompDoubles()
	_ = x
	return (x == float64(0.0))
}

// Java static method: JSCompDoubles.ecmascriptToInt32(double number)
func JSCompDoubles_EcmascriptToInt32(number float64) int32 {
	clinit_JSCompDoubles()
	_ = number
	if JSCompDoubles_IsExactInt32(number) {
		return javalang.D2I(number)
	}
	if math.IsNaN(number) || math.IsInf(number, 0) {
		return int32(0)
	}
	var posInt float64
	_ = posInt
	if number >= float64(int32(0)) {
		posInt = math.Floor(number)
	} else {
		posInt = math.Ceil(number)
	}
	var int32bit int64 = javalang.D2L((math.Mod(posInt, StaticGet_JSCompDoubles_POW_2_32())))
	_ = int32bit
	return int32(int32bit)
}

// Java static method: JSCompDoubles.ecmascriptToUint32(double number)
func JSCompDoubles_EcmascriptToUint32(number float64) int32 {
	clinit_JSCompDoubles()
	_ = number
	return JSCompDoubles_EcmascriptToInt32(number)
}
