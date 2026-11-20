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
// Port of com/google/javascript/jscomp/transpile/Transpiler.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.transpile.Transpiler
type Transpile_Transpiler interface {
	Implements_jscomp_Transpile_Transpiler()
	IsJavaNil_() bool
	Transpile(*j2gort.Ext_java_net_URI, string) *TranspileResult
	Runtime() string
}

func ToIface_Transpile_Transpiler[T interface {
	Transpile_Transpiler
	comparable
}](v T) Transpile_Transpiler {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Transpile_Transpiler[T Transpile_Transpiler](v T) Transpile_Transpiler {
	var vi Transpile_Transpiler = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.transpile.Transpiler", "com.google.javascript.jscomp.transpile.Transpiler", func(v_ any) bool {
		return j2gort.InstanceOf[Transpile_Transpiler](v_)
	})
}

var static_Transpile_Transpiler_NULL Transpile_Transpiler

func StaticGet_Transpile_Transpiler_NULL() Transpile_Transpiler {
	clinit_Transpile_Transpiler()
	return static_Transpile_Transpiler_NULL
}
func StaticSet_Transpile_Transpiler_NULL(v Transpile_Transpiler) {
	clinit_Transpile_Transpiler()
	static_Transpile_Transpiler_NULL = v
}

var clinitDone_Transpile_Transpiler uint32
var clinitState_Transpile_Transpiler j2gort.ClinitState

func clinit_Transpile_Transpiler() {
	if atomic.LoadUint32(&clinitDone_Transpile_Transpiler) != 0 {
		return
	}
	clinitSlow_Transpile_Transpiler()
}

func clinitSlow_Transpile_Transpiler() {
	if !j2gort.EnterClinit(&clinitState_Transpile_Transpiler, &clinitDone_Transpile_Transpiler) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Transpile_Transpiler, &clinitDone_Transpile_Transpiler)
	static_Transpile_Transpiler_NULL = NewTranspile_Transpiler__Anon_129()
}

// Java class:
type Transpile_Transpiler__Anon_129 struct {
	idHash_ atomic.Int32
}

func (this_ *Transpile_Transpiler__Anon_129) BindSelf_(s any) {
	_ = s
}
func (this_ *Transpile_Transpiler__Anon_129) IsJavaNil_() bool { return this_ == nil }
func (*Transpile_Transpiler__Anon_129) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.Transpiler$1", ""
}
func (this_ *Transpile_Transpiler__Anon_129) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Transpile_Transpiler__Anon_129) HasIdentityEquals_()                     {}
func (*Transpile_Transpiler__Anon_129) Implements_jscomp_Transpile_Transpiler() {}

func (this_ *Transpile_Transpiler__Anon_129) initFields_Transpile_Transpiler__Anon_129() {
	_ = this_
}

func NewTranspile_Transpiler__Anon_129() *Transpile_Transpiler__Anon_129 {
	this_ := &Transpile_Transpiler__Anon_129{}
	this_.initFields_Transpile_Transpiler__Anon_129()
	return this_
}

// Java method: .transpile(java.net.URI path, java.lang.String code)
func (this_ *Transpile_Transpiler__Anon_129) Transpile(path *j2gort.Ext_java_net_URI, code string) *TranspileResult {
	_ = path
	_ = code
	return NewTranspileResult(path, code, code, "")
}

// Java method: .runtime()
func (this_ *Transpile_Transpiler__Anon_129) Runtime() string {
	return ""
}
