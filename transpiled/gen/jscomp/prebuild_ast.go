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
// Port of com/google/javascript/jscomp/PrebuildAst.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_prebuild_ast_1 j2gort.Ext_java_util_concurrent_ThreadFactory

func init() {
	lambdaSingleton_prebuild_ast_1 = j2gort.NewLambda_Ext_java_util_concurrent_ThreadFactory(func(r j2gort.Ext_java_lang_Runnable) *j2gort.Ext_java_lang_Thread {
		_ = r
		var t *j2gort.Ext_java_lang_Thread = j2gort.NewExt_java_lang_Thread_4(nil, r, "jscompiler-PrebuildAst", int64(67108864))
		_ = t
		t.SetDaemon(true)
		return t
	})
}

// Java class: com.google.javascript.jscomp.PrebuildAst
type PrebuildAst struct {
	compiler           AbstractCompiler `gson:"compiler"`
	numParallelThreads int32            `gson:"numParallelThreads"`
	idHash_            atomic.Int32
}

func (this_ *PrebuildAst) BindSelf_(s any) {
	_ = s
}
func (this_ *PrebuildAst) IsJavaNil_() bool { return this_ == nil }
func (*PrebuildAst) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PrebuildAst", "PrebuildAst"
}
func (this_ *PrebuildAst) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PrebuildAst) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PrebuildAst", "com.google.javascript.jscomp.PrebuildAst", func(v_ any) bool {
		return j2gort.InstanceOf[*PrebuildAst](v_)
	})
}

func (this_ *PrebuildAst) initFields_PrebuildAst() {
	_ = this_
}

// Java constructor: PrebuildAst(com.google.javascript.jscomp.AbstractCompiler compiler, int numParalleThreads)
func NewPrebuildAst(compiler AbstractCompiler, numParalleThreads int32) *PrebuildAst {
	this_ := &PrebuildAst{}
	InitPrebuildAst(this_, compiler, numParalleThreads)
	return this_
}

func InitPrebuildAst(this_ *PrebuildAst, compiler AbstractCompiler, numParalleThreads int32) {
	_ = compiler
	_ = numParalleThreads
	this_.initFields_PrebuildAst()
	this_.compiler = compiler
	this_.numParallelThreads = numParalleThreads
}

// Java method: PrebuildAst.prebuild(java.lang.Iterable<com.google.javascript.jscomp.CompilerInput> allInputs)
func (this_ *PrebuildAst) Prebuild(allInputs j2gort.Ext_java_lang_Iterable) {
	_ = allInputs
	if this_.numParallelThreads <= int32(1) {
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(allInputs))
			for cur_1.HasNext() {
				rawElem_2 := cur_1.Next()
				input := j2gort.CastPtr[CompilerInput](rawElem_2)
				_ = input
				input.GetAstRoot(this_.compiler)
			}
		}
		return
	}
	var threadFactory j2gort.Ext_java_util_concurrent_ThreadFactory = lambdaSingleton_prebuild_ast_1
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
			futureList.Add_1(j2gort.AsIface[j2gort.Ext_com_google_common_util_concurrent_ListenableFuture](executorService.Submit_1_Ext_java_util_concurrent_Callable(j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
				return j2gort.BoxNode(input_2.GetAstRoot(this_.compiler))
			}))))
		}
	}
	poolExecutor.Shutdown()
	func() {
		defer func() {
			if rec_5 := recover(); rec_5 != nil {
				jex_6 := javalang.RecoverException(rec_5)
				_ = jex_6
				if j2gort.InstanceOf[j2gort.Ext_java_lang_InterruptedException](jex_6) || j2gort.InstanceOf[j2gort.Ext_java_util_concurrent_ExecutionException](jex_6) {
					e := jex_6
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_5)
				}
			}
		}()
		j2gort.Ext_com_google_common_util_concurrent_Futures_AllAsList_1_Ext_java_lang_Iterable(futureList).Get_0()
	}()
}
