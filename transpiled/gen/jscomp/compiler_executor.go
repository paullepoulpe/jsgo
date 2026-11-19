/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CompilerExecutor.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CompilerExecutor
type CompilerExecutor struct {
	compilerThread *j2gort.Ext_java_lang_Thread `gson:"compilerThread"`
	debugMessage   string                       `gson:"debugMessage"`
	timeout        int32                        `gson:"timeout"`
	idHash_        atomic.Int32
	useThreads     bool `gson:"useThreads"`
}

func (this_ *CompilerExecutor) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerExecutor) IsJavaNil_() bool { return this_ == nil }
func (*CompilerExecutor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerExecutor", "CompilerExecutor"
}
func (this_ *CompilerExecutor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CompilerExecutor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CompilerExecutor", "com.google.javascript.jscomp.CompilerExecutor", func(v_ any) bool {
		return j2gort.InstanceOf[*CompilerExecutor](v_)
	})
}

func StaticGet_CompilerExecutor_COMPILER_STACK_SIZE() int64 { return int64(67108864) }
func (this_ *CompilerExecutor) initFields_CompilerExecutor() {
	_ = this_

	this_.useThreads = true

	this_.debugMessage = j2gort.NullString
}

// Java constructor: CompilerExecutor()
func NewCompilerExecutor() *CompilerExecutor {
	this_ := &CompilerExecutor{}
	InitCompilerExecutor(this_)
	return this_
}

func InitCompilerExecutor(this_ *CompilerExecutor) {
	this_.initFields_CompilerExecutor()
}

// Java method: CompilerExecutor.getExecutorService()
func (this_ *CompilerExecutor) GetExecutorService() j2gort.Ext_java_util_concurrent_ExecutorService {
	return CompilerExecutor_GetDefaultExecutorService()
}

// Java static method: CompilerExecutor.getDefaultExecutorService()
func CompilerExecutor_GetDefaultExecutorService() j2gort.Ext_java_util_concurrent_ExecutorService {
	return j2gort.Ext_java_util_concurrent_Executors_NewSingleThreadExecutor_1(NewCompilerExecutor__Anon_14())
}

// Java method: CompilerExecutor.disableThreads()
func (this_ *CompilerExecutor) DisableThreads() {
	this_.useThreads = false
}

// Java method: CompilerExecutor.setTimeout(int timeout)
func (this_ *CompilerExecutor) SetTimeout(timeout int32) {
	_ = timeout
	this_.timeout = timeout
}

// Java method: CompilerExecutor.setDebugMessage(java.lang.String debugMessage)
func (this_ *CompilerExecutor) SetDebugMessage(debugMessage string) {
	_ = debugMessage
	this_.debugMessage = debugMessage
}

// Java method: CompilerExecutor.runInCompilerThread(java.util.concurrent.Callable<T> callable, boolean dumpTraceReport)
func (this_ *CompilerExecutor) RunInCompilerThread(callable j2gort.Ext_java_util_concurrent_Callable, dumpTraceReport bool) any {
	_ = callable
	_ = dumpTraceReport
	var executor j2gort.Ext_java_util_concurrent_ExecutorService = this_.GetExecutorService()
	_ = executor
	var result any = nil
	_ = result
	var exception []javalang.Throwable = make([]javalang.Throwable, int32(1))
	_ = exception
	if !((this_.compilerThread == nil) || (this_.compilerThread == j2gort.Ext_java_lang_Thread_CurrentThread())) {
		j2gort.CheckStateStr(false, "Please do not share the Compiler across threads")
	}
	if this_.useThreads && (this_.compilerThread == nil) {
		func() {
			defer func() {
				executor.Shutdown()
			}()
			defer func() {
				if rec_1 := recover(); rec_1 != nil {
					jex_2 := javalang.RecoverException(rec_1)
					_ = jex_2
					if j2gort.InstanceOf[j2gort.Ext_java_lang_InterruptedException](jex_2) || j2gort.InstanceOf[j2gort.Ext_java_util_concurrent_TimeoutException](jex_2) || j2gort.InstanceOf[j2gort.Ext_java_util_concurrent_ExecutionException](jex_2) {
						e := jex_2
						_ = e
						panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
					} else {
						panic(rec_1)
					}
				}
			}()
			var bootCompilerThread j2gort.Ext_java_util_concurrent_Callable = j2gort.NewLambda_Ext_java_util_concurrent_Callable(func() any {
				var didRet_3 bool
				var retVal_4 any
				func() {
					defer func() {
						this_.compilerThread = nil
						if dumpTraceReport {
							Tracer_LogCurrentThreadTrace()
						}
						Tracer_ClearCurrentThreadTrace()
					}()
					defer func() {
						if rec_5 := recover(); rec_5 != nil {
							jex_6 := javalang.RecoverException(rec_5)
							_ = jex_6
							if j2gort.InstanceOf[javalang.Throwable](jex_6) {
								e_2 := j2gort.AsIface[javalang.Throwable](jex_6)
								_ = e_2
								{
									a_7 := exception
									i_8 := int32(0)
									a_7[i_8] = e_2
								}
							} else {
								panic(rec_5)
							}
						}
					}()
					this_.compilerThread = j2gort.Ext_java_lang_Thread_CurrentThread()
					if dumpTraceReport {
						Tracer_InitCurrentThreadTrace_0()
					}
					retVal_4 = callable.Call()
					didRet_3 = true
					return
				}()
				if didRet_3 {
					return retVal_4
				}
				return nil
			})
			_ = bootCompilerThread
			var future j2gort.Ext_java_util_concurrent_Future = j2gort.AsIface[j2gort.Ext_java_util_concurrent_Future](executor.Submit_1_Ext_java_util_concurrent_Callable(bootCompilerThread))
			_ = future
			if this_.timeout > int32(0) {
				result = future.Get_2(int64(this_.timeout), j2gort.Ext_java_util_concurrent_TimeUnit_F_SECONDS)
			} else {
				result = future.Get_0()
			}
		}()
	} else {
		func() {
			defer func() {
				if rec_9 := recover(); rec_9 != nil {
					jex_10 := javalang.RecoverException(rec_9)
					_ = jex_10
					if j2gort.InstanceOf[javalang.IException](jex_10) {
						e_2 := j2gort.AsIface[javalang.IException](jex_10)
						_ = e_2
						{
							a_11 := exception
							i_12 := int32(0)
							a_11[i_12] = j2gort.AsIface[javalang.Throwable](e_2)
						}
					} else {
						panic(rec_9)
					}
				}
			}()
			result = callable.Call()
		}()
	}
	if exception[int32(0)] != nil {
		if this_.debugMessage != j2gort.NullString {
			panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Exception during compilation: "+j2gort.StrValString(this_.debugMessage), exception[int32(0)])))
		}
		j2gort.Ext_com_google_common_base_Throwables_ThrowIfUnchecked(exception[int32(0)])
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", exception[int32(0)])))
	}
	return result
}

// Java class:
type CompilerExecutor__Anon_14 struct {
	idHash_ atomic.Int32
}

func (this_ *CompilerExecutor__Anon_14) BindSelf_(s any) {
	_ = s
}
func (this_ *CompilerExecutor__Anon_14) IsJavaNil_() bool { return this_ == nil }
func (*CompilerExecutor__Anon_14) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CompilerExecutor$1", ""
}
func (this_ *CompilerExecutor__Anon_14) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CompilerExecutor__Anon_14) HasIdentityEquals_()                                       {}
func (*CompilerExecutor__Anon_14) Implements_j2gort_Ext_java_util_concurrent_ThreadFactory() {}

func (this_ *CompilerExecutor__Anon_14) initFields_CompilerExecutor__Anon_14() {
	_ = this_
}

func NewCompilerExecutor__Anon_14() *CompilerExecutor__Anon_14 {
	this_ := &CompilerExecutor__Anon_14{}
	this_.initFields_CompilerExecutor__Anon_14()
	return this_
}

// Java method: .newThread(java.lang.Runnable r)
func (this_ *CompilerExecutor__Anon_14) NewThread(r j2gort.Ext_java_lang_Runnable) *j2gort.Ext_java_lang_Thread {
	_ = r
	var t *j2gort.Ext_java_lang_Thread = j2gort.NewExt_java_lang_Thread_4(nil, r, "jscompiler", int64(67108864))
	_ = t
	t.SetDaemon(true)
	return t
}
