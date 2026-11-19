/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Xid.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)


var lambdaSingleton_xid_1 Xid_HashFunction
func init() {
	lambdaSingleton_xid_1 = NewLambda_Xid_HashFunction(func(mrp0_ string) int32 { return javalang.StringHashCode(j2gort.CheckNotNullString(mrp0_)) })
}
// Java class: com.google.javascript.jscomp.Xid
type Xid struct {
	hasher Xid_HashFunction `gson:"hasher"`
	idHash_ atomic.Int32
}

func (this_ *Xid) BindSelf_(s any) {
	_ = s
}
func (this_ *Xid) IsJavaNil_() bool { return this_ == nil }
func (*Xid) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Xid", "Xid" }
func (this_ *Xid) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Xid) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Xid", "com.google.javascript.jscomp.Xid", func(v_ any) bool {
		return j2gort.InstanceOf[*Xid](v_)
	})
}

func StaticGet_Xid_START_CHARS() string { return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" }
func StaticGet_Xid_CHARS() string { return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" }
var static_Xid_START_RADIX int32
func StaticGet_Xid_START_RADIX() int32 {
	clinit_Xid()
	return static_Xid_START_RADIX
}
func StaticSet_Xid_START_RADIX(v int32) {
	clinit_Xid()
	static_Xid_START_RADIX = v
}
var static_Xid_RADIX int32
func StaticGet_Xid_RADIX() int32 {
	clinit_Xid()
	return static_Xid_RADIX
}
func StaticSet_Xid_RADIX(v int32) {
	clinit_Xid()
	static_Xid_RADIX = v
}
var clinitDone_Xid uint32
var clinitState_Xid j2gort.ClinitState
func clinit_Xid() {
	if atomic.LoadUint32(&clinitDone_Xid) != 0 {
		return
	}
	clinitSlow_Xid()
}

func clinitSlow_Xid() {
	if !j2gort.EnterClinit(&clinitState_Xid, &clinitDone_Xid) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Xid, &clinitDone_Xid)
	static_Xid_START_RADIX = j2gort.JStrLen("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	static_Xid_RADIX = j2gort.JStrLen("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

func (this_ *Xid) initFields_Xid() {
	_ = this_
}

// Java constructor: Xid()
func NewXid() *Xid {
	clinit_Xid()
	this_ := &Xid{}
	InitXid(this_)
	return this_
}

func InitXid(this_ *Xid) {
	this_.initFields_Xid()
	this_.hasher = lambdaSingleton_xid_1
}

// Java constructor: Xid(com.google.javascript.jscomp.Xid.HashFunction hasher)
func NewXid_1(hasher Xid_HashFunction) *Xid {
	clinit_Xid()
	this_ := &Xid{}
	InitXid_1(this_, hasher)
	return this_
}

func InitXid_1(this_ *Xid, hasher Xid_HashFunction) {
	_ = hasher
	this_.initFields_Xid()
	this_.hasher = hasher
}

// Java method: Xid.get(java.lang.String key)
func (this_ *Xid) Get(key string) string {
	_ = key
	return Xid_ToString(this_.GetAsInt(key))
}

// Java method: Xid.getAsInt(java.lang.String key)
func (this_ *Xid) GetAsInt(key string) int32 {
	_ = key
	return this_.hasher.M_HashCode(key)
}

// Java static method: Xid.toString(int i)
func Xid_ToString(i int32) string {
	clinit_Xid()
	_ = i
	var buf []uint16 = make([]uint16, int32(6))
	_ = buf
	var len_ int32 = int32(0)
	_ = len_
	var l int64 = (int64(i) - int64(int32(-2147483648)))
	_ = l
	{ a_1 := buf; i_2 := func() int32 { old_3 := len_; nw_4 := old_3 + 1; len_ = nw_4; return old_3 }(); a_1[i_2] = j2gort.JStrCharAt("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", int32(((l % int64(StaticGet_Xid_START_RADIX()))))) }
	i = int32(((l / int64(StaticGet_Xid_START_RADIX()))))
	for ((i > int32(0))) {
		{ a_5 := buf; i_6 := func() int32 { old_7 := len_; nw_8 := old_7 + 1; len_ = nw_8; return old_7 }(); a_5[i_6] = j2gort.JStrCharAt("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", (i % StaticGet_Xid_RADIX())) }
		i = (i / StaticGet_Xid_RADIX())
	}
	return func() string { ca_9 := buf; co_10 := int32(0); cc_11 := len_; return javalang.StringFromUTF16(ca_9[co_10 : co_10+cc_11]) }()
}

// Java interface: com.google.javascript.jscomp.Xid.HashFunction
type Xid_HashFunction interface {
	Implements_jscomp_Xid_HashFunction()
	IsJavaNil_() bool
	M_HashCode(string) int32
}
func ToIface_Xid_HashFunction[T interface {
	Xid_HashFunction
	comparable
}](v T) Xid_HashFunction {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Xid_HashFunction[T Xid_HashFunction](v T) Xid_HashFunction {
	var vi Xid_HashFunction = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Xid.HashFunction", "com.google.javascript.jscomp.Xid$HashFunction", func(v_ any) bool {
		return j2gort.InstanceOf[Xid_HashFunction](v_)
	})
}
type lambda_Xid_HashFunction struct {
	fn func(string) int32
}

func (l *lambda_Xid_HashFunction) IsJavaNil_() bool { return l == nil }
func NewLambda_Xid_HashFunction(fn func(string) int32) Xid_HashFunction {
	if fn == nil {
		return nil
	}
	return &lambda_Xid_HashFunction{fn: fn}
}

func (*lambda_Xid_HashFunction) Implements_jscomp_Xid_HashFunction() {}
func (l *lambda_Xid_HashFunction) M_HashCode(p0_ string) int32 {
	return l.fn(p0_)
}

