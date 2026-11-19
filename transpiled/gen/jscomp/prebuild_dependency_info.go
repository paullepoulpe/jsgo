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
// Port of com/google/javascript/jscomp/PrebuildDependencyInfo.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PrebuildDependencyInfo
type PrebuildDependencyInfo struct {
	numParallelThreads int32 `gson:"numParallelThreads"`
	idHash_            atomic.Int32
}

func (this_ *PrebuildDependencyInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *PrebuildDependencyInfo) IsJavaNil_() bool { return this_ == nil }
func (*PrebuildDependencyInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PrebuildDependencyInfo", "PrebuildDependencyInfo"
}
func (this_ *PrebuildDependencyInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PrebuildDependencyInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PrebuildDependencyInfo", "com.google.javascript.jscomp.PrebuildDependencyInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*PrebuildDependencyInfo](v_)
	})
}

func (this_ *PrebuildDependencyInfo) initFields_PrebuildDependencyInfo() {
	_ = this_
}

// Java constructor: PrebuildDependencyInfo(int numParalleThreads)
func NewPrebuildDependencyInfo(numParalleThreads int32) *PrebuildDependencyInfo {
	this_ := &PrebuildDependencyInfo{}
	InitPrebuildDependencyInfo(this_, numParalleThreads)
	return this_
}

func InitPrebuildDependencyInfo(this_ *PrebuildDependencyInfo, numParalleThreads int32) {
	_ = numParalleThreads
	this_.initFields_PrebuildDependencyInfo()
	this_.numParallelThreads = numParalleThreads
}

// Java method: PrebuildDependencyInfo.prebuild(java.lang.Iterable<com.google.javascript.jscomp.CompilerInput> allInputs)
func (this_ *PrebuildDependencyInfo) Prebuild(allInputs j2gort.Ext_java_lang_Iterable) {
	_ = allInputs
	if this_.numParallelThreads <= int32(1) {
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(allInputs))
			for cur_1.HasNext() {
				rawElem_2 := cur_1.Next()
				input := j2gort.CastPtr[CompilerInput](rawElem_2)
				_ = input
				input.GetDependencyInfo()
			}
		}
		return
	}
	var threadFactory j2gort.Ext_java_util_concurrent_ThreadFactory = NewPrebuildDependencyInfo__Anon_71(this_)
	_ = threadFactory
	var poolExecutor *j2gort.Ext_java_util_concurrent_ThreadPoolExecutor = j2gort.NewExt_java_util_concurrent_ThreadPoolExecutor_6_int32_int32_int64_Ext_java_util_concurrent_TimeUnit_Ext_java_util_concurrent_BlockingQueue_Ext_java_util_concurrent_ThreadFactory(this_.numParallelThreads, this_.numParallelThreads, int64(int32(2147483647)), j2gort.Ext_java_util_concurrent_TimeUnit_F_SECONDS, j2gort.NewExt_java_util_concurrent_LinkedBlockingQueue(), threadFactory)
	_ = poolExecutor
	var executorService j2gort.Ext_com_google_common_util_concurrent_ListeningExecutorService = j2gort.Ext_com_google_common_util_concurrent_MoreExecutors_S_ListeningDecorator_1_Ext_java_util_concurrent_ExecutorService(poolExecutor)
	_ = executorService
	var futureList j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_int32(j2gort.Ext_com_google_common_collect_Iterables_Size(allInputs))
	_ = futureList
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(allInputs))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			input_2 := j2gort.CastPtr[CompilerInput](rawElem_4)
			_ = input_2
			futureList.Add_1(j2gort.AsIface[j2gort.Ext_com_google_common_util_concurrent_ListenableFuture](executorService.Submit_1_Ext_java_util_concurrent_Callable(func() j2gort.Ext_java_util_concurrent_Callable {
				bound_5 := input_2
				return j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any { return bound_5.GetDependencyInfo() })
			}())))
		}
	}
	poolExecutor.Shutdown()
	func() {
		defer func() {
			if rec_6 := recover(); rec_6 != nil {
				jex_7 := javalang.RecoverException(rec_6)
				_ = jex_7
				if j2gort.InstanceOf[j2gort.Ext_java_lang_InterruptedException](jex_7) || j2gort.InstanceOf[j2gort.Ext_java_util_concurrent_ExecutionException](jex_7) {
					e := jex_7
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_6)
				}
			}
		}()
		j2gort.Ext_com_google_common_util_concurrent_Futures_AllAsList_1_Ext_java_lang_Iterable(futureList).Get_0()
	}()
}

// Java class:
type PrebuildDependencyInfo__Anon_71 struct {
	outer_  *PrebuildDependencyInfo
	idHash_ atomic.Int32
}

func (this_ *PrebuildDependencyInfo__Anon_71) BindSelf_(s any) {
	_ = s
}
func (this_ *PrebuildDependencyInfo__Anon_71) IsJavaNil_() bool { return this_ == nil }
func (*PrebuildDependencyInfo__Anon_71) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PrebuildDependencyInfo$1", ""
}
func (this_ *PrebuildDependencyInfo__Anon_71) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PrebuildDependencyInfo__Anon_71) HasIdentityEquals_()                       {}
func (this_ *PrebuildDependencyInfo__Anon_71) SetOuter_(o *PrebuildDependencyInfo) { this_.outer_ = o }
func (this_ *PrebuildDependencyInfo__Anon_71) GetOuter_() *PrebuildDependencyInfo {
	return this_.outer_
}
func (*PrebuildDependencyInfo__Anon_71) Implements_j2gort_Ext_java_util_concurrent_ThreadFactory() {}

func (this_ *PrebuildDependencyInfo__Anon_71) initFields_PrebuildDependencyInfo__Anon_71() {
	_ = this_
}

func NewPrebuildDependencyInfo__Anon_71(outer_ *PrebuildDependencyInfo) *PrebuildDependencyInfo__Anon_71 {
	this_ := &PrebuildDependencyInfo__Anon_71{}
	this_.outer_ = outer_
	this_.initFields_PrebuildDependencyInfo__Anon_71()
	return this_
}

// Java method: .newThread(java.lang.Runnable r)
func (this_ *PrebuildDependencyInfo__Anon_71) NewThread(r j2gort.Ext_java_lang_Runnable) *j2gort.Ext_java_lang_Thread {
	_ = r
	var t *j2gort.Ext_java_lang_Thread = j2gort.NewExt_java_lang_Thread_4(nil, r, "jscompiler-PrebuildDependencyInfo", int64(67108864))
	_ = t
	t.SetDaemon(true)
	return t
}
