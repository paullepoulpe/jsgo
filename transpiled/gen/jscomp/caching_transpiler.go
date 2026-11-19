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
// Port of com/google/javascript/jscomp/transpile/CachingTranspiler.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.transpile.CachingTranspiler
type CachingTranspiler struct {
	cache   j2gort.Ext_com_google_common_cache_LoadingCache `gson:"cache"`
	runtime j2gort.Ext_com_google_common_base_Supplier      `gson:"runtime"`
	idHash_ atomic.Int32
}

func (this_ *CachingTranspiler) BindSelf_(s any) {
	_ = s
}
func (this_ *CachingTranspiler) IsJavaNil_() bool { return this_ == nil }
func (*CachingTranspiler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.CachingTranspiler", "CachingTranspiler"
}
func (this_ *CachingTranspiler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CachingTranspiler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.transpile.CachingTranspiler", "com.google.javascript.jscomp.transpile.CachingTranspiler", func(v_ any) bool {
		return j2gort.InstanceOf[*CachingTranspiler](v_)
	})
}
func (*CachingTranspiler) Implements_jscomp_Transpile_Transpiler() {}

func (this_ *CachingTranspiler) initFields_CachingTranspiler() {
	_ = this_
}

// Java constructor: CachingTranspiler(com.google.javascript.jscomp.transpile.Transpiler delegate, com.google.common.cache.CacheBuilder<java.lang.Object,? super com.google.javascript.jscomp.transpile.TranspileResult> builder)
func NewCachingTranspiler(delegate Transpile_Transpiler, builder *j2gort.Ext_com_google_common_cache_CacheBuilder) *CachingTranspiler {
	this_ := &CachingTranspiler{}
	InitCachingTranspiler(this_, delegate, builder)
	return this_
}

func InitCachingTranspiler(this_ *CachingTranspiler, delegate Transpile_Transpiler, builder *j2gort.Ext_com_google_common_cache_CacheBuilder) {
	_ = delegate
	_ = builder
	this_.initFields_CachingTranspiler()
	j2gort.CheckNotNullIface(delegate)
	this_.cache = builder.Build_1(NewCachingTranspiler__Anon_127(this_, delegate))
	this_.runtime = j2gort.Ext_com_google_common_base_Suppliers_Memoize(func() j2gort.Ext_com_google_common_base_Supplier {
		bound_1 := delegate
		return j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any { return j2gort.BoxString(bound_1.Runtime()) })
	}())
}

// Java method: CachingTranspiler.transpile(java.net.URI path, java.lang.String code)
func (this_ *CachingTranspiler) Transpile(path *j2gort.Ext_java_net_URI, code string) *TranspileResult {
	_ = path
	_ = code
	var didRet_2 bool
	var retVal_3 *TranspileResult
	func() {
		defer func() {
			if rec_4 := recover(); rec_4 != nil {
				jex_5 := javalang.RecoverException(rec_4)
				_ = jex_5
				if j2gort.InstanceOf[j2gort.Ext_com_google_common_util_concurrent_UncheckedExecutionException](jex_5) {
					e := j2gort.AsIface[j2gort.Ext_com_google_common_util_concurrent_UncheckedExecutionException](jex_5)
					_ = e
					if j2gort.InstanceOf[*TranspilationException](e.GetCause()) {
						panic(NewTranspilationException_1(j2gort.AsIface[javalang.IException](e)))
					} else {
						panic(j2gort.CheckNotNull(e))
					}
				} else {
					panic(rec_4)
				}
			}
		}()
		retVal_3 = j2gort.CastPtr[TranspileResult](this_.cache.GetUnchecked(NewCachingTranspiler_Key(path, code)))
		didRet_2 = true
		return
	}()
	if didRet_2 {
		return retVal_3
	}
	panic("unreachable")
}

// Java method: CachingTranspiler.runtime()
func (this_ *CachingTranspiler) Runtime() string {
	return j2gort.CastString(this_.runtime.Get())
}

// Java class: com.google.javascript.jscomp.transpile.CachingTranspiler.Key
type CachingTranspiler_Key struct {
	path *j2gort.Ext_java_net_URI `gson:"path"`
	code string                   `gson:"code"`
}

func (this_ *CachingTranspiler_Key) BindSelf_(s any) {
	_ = s
}
func (this_ *CachingTranspiler_Key) IsJavaNil_() bool { return this_ == nil }
func (*CachingTranspiler_Key) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.CachingTranspiler$Key", "Key"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.transpile.CachingTranspiler.Key", "com.google.javascript.jscomp.transpile.CachingTranspiler$Key", func(v_ any) bool {
		return j2gort.InstanceOf[*CachingTranspiler_Key](v_)
	})
}

func (this_ *CachingTranspiler_Key) initFields_CachingTranspiler_Key() {
	_ = this_
	this_.code = j2gort.NullString
}

// Java constructor: Key(java.net.URI path, java.lang.String code)
func NewCachingTranspiler_Key(path *j2gort.Ext_java_net_URI, code string) *CachingTranspiler_Key {
	this_ := &CachingTranspiler_Key{}
	InitCachingTranspiler_Key(this_, path, code)
	return this_
}

func InitCachingTranspiler_Key(this_ *CachingTranspiler_Key, path *j2gort.Ext_java_net_URI, code string) {
	_ = path
	_ = code
	this_.initFields_CachingTranspiler_Key()
	this_.path = j2gort.CheckNotNullPtr(path)
	this_.code = j2gort.CheckNotNullString(code)
}

// Java method: Key.equals(java.lang.Object other)
func (this_ *CachingTranspiler_Key) Equals(arg_6 any) bool {
	other := arg_6
	_ = other
	var key *CachingTranspiler_Key
	_ = key
	var ok_7 bool
	key, ok_7 = j2gort.TryCast[*CachingTranspiler_Key](other)
	return ok_7 && (key.path.Equals(j2gort.PtrToAny(this_.path))) && (func() bool { sr_8 := key.code; se_9 := this_.code; return j2gort.CheckNotNullString(sr_8) == se_9 }())
}

// Java method: Key.hashCode()
func (this_ *CachingTranspiler_Key) HashCode() int32 {
	{
		var oh_10 int32 = 1
		oh_10 = 31*oh_10 + j2gort.ObjectHash(j2gort.PtrToAny(this_.path))
		{
			hs_11 := this_.code
			oh_10 = 31 * oh_10
			if hs_11 != j2gort.NullString {
				oh_10 += javalang.StringHashCode(hs_11)
			}
		}
		return oh_10
	}
}

// Java class:
type CachingTranspiler__Anon_127 struct {
	j2gort.Ext_com_google_common_cache_CacheLoaderBase
	outer_         *CachingTranspiler
	cap_0_delegate Transpile_Transpiler
	idHash_        atomic.Int32
}

func (this_ *CachingTranspiler__Anon_127) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_cache_CacheLoaderBase.BindExtSelf_(s)
}
func (this_ *CachingTranspiler__Anon_127) IsJavaNil_() bool { return this_ == nil }
func (*CachingTranspiler__Anon_127) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.transpile.CachingTranspiler$1", ""
}
func (this_ *CachingTranspiler__Anon_127) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (this_ *CachingTranspiler__Anon_127) SetOuter_(o *CachingTranspiler)                       { this_.outer_ = o }
func (this_ *CachingTranspiler__Anon_127) GetOuter_() *CachingTranspiler                        { return this_.outer_ }
func (*CachingTranspiler__Anon_127) Implements_j2gort_Ext_com_google_common_cache_CacheLoader() {}

func (this_ *CachingTranspiler__Anon_127) initFields_CachingTranspiler__Anon_127() {
	_ = this_
}

func NewCachingTranspiler__Anon_127(outer_ *CachingTranspiler, cap_p_delegate Transpile_Transpiler) *CachingTranspiler__Anon_127 {
	this_ := &CachingTranspiler__Anon_127{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	this_.cap_0_delegate = cap_p_delegate
	j2gort.InitExt_com_google_common_cache_CacheLoader(&this_.Ext_com_google_common_cache_CacheLoaderBase)
	this_.initFields_CachingTranspiler__Anon_127()
	return this_
}

// Java method: .load(com.google.javascript.jscomp.transpile.CachingTranspiler.Key key)
func (this_ *CachingTranspiler__Anon_127) Load(arg_12 any) any {
	key := j2gort.CastPtr[CachingTranspiler_Key](arg_12)
	_ = key
	return j2gort.PtrToAny(this_.cap_0_delegate.Transpile(key.path, key.code))
}
