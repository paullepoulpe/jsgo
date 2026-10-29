// Port of com/google/javascript/jscomp/CompilerExecutor.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CompilerExecutor runs the compiler in a separate thread with a larger
// stack.
//
// Go-only: language-adapter — the port is sequential (PORTING.md §5). In
// Java, runInCompilerThread submits the callable to a single-thread
// executor whose daemon thread "jscompiler" has a 64MB stack (the compiler
// recurses O(depth) on the AST), waits for it with future.get() (or
// future.get(timeout, SECONDS)), and shuts the executor down. Here the
// callable runs inline on the calling goroutine: Go stacks grow on demand,
// so no dedicated stack is needed, and the timeout is recorded but cannot
// interrupt the inline call. The per-thread Tracer trace of the Java
// compiler thread is the single Tracer trace of the port. Everything else
// (the re-entrancy check, the bookkeeping of compilerThread, which
// exceptions are caught and how they are rethrown) is Java's.
type CompilerExecutor struct {
	// Use a dedicated compiler thread per Compiler instance.
	compilerThread *jutil.Thread

	// Whether to use threads.
	useThreads bool

	timeout int

	debugMessage *string
}

// CompilerExecutor_COMPILER_STACK_SIZE: we use many recursive algorithms
// that use O(d) memory in the depth of the tree. Also, (de)serialization
// between phases can involve a lot of recursion.
const CompilerExecutor_COMPILER_STACK_SIZE int64 = 1 << 26 // About 64MB

// NewCompilerExecutor is new CompilerExecutor().
func NewCompilerExecutor() *CompilerExecutor {
	return &CompilerExecutor{useThreads: true}
}

// Java: CompilerExecutor.disableThreads()
func (e *CompilerExecutor) DisableThreads() {
	e.useThreads = false
}

// Java: CompilerExecutor.setTimeout(int)
func (e *CompilerExecutor) SetTimeout(timeout int) {
	e.timeout = timeout
}

// Java: CompilerExecutor.setDebugMessage(String)
func (e *CompilerExecutor) SetDebugMessage(debugMessage *string) {
	e.debugMessage = debugMessage
}

// CompilerExecutor_RunInCompilerThread is runInCompilerThread.
//
// Go-only: language-adapter — a generic method, so a package function;
// Callable<T> is a func() T.
//
// Java: CompilerExecutor.runInCompilerThread(Callable<T>, boolean)
func CompilerExecutor_RunInCompilerThread[T any](e *CompilerExecutor, callable func() T, dumpTraceReport bool) (result T) {
	if !(e.compilerThread == nil || e.compilerThread == jutil.Thread_CurrentThread()) {
		// checkState(boolean, Object errorMessage)
		panic(jutil.NewIllegalStateExceptionMsg("Please do not share the Compiler across threads"))
	}

	bootThread := e.useThreads && e.compilerThread == nil
	defer func() {
		r := recover()
		if bootThread {
			e.compilerThread = nil
			if dumpTraceReport {
				Tracer_LogCurrentThreadTrace()
			}
			Tracer_ClearCurrentThreadTrace()
		}
		if r == nil {
			return
		}
		exception, ok := r.(jutil.Throwable)
		if !ok {
			exception = jutil.AsThrowable(r)
		}
		if !bootThread {
			// Java catches Exception: an Error propagates.
			if _, isError := exception.(jutil.IError); isError {
				panic(exception)
			}
		}
		if e.debugMessage != nil {
			panic(jutil.NewRuntimeExceptionWithMessageAndCause("Exception during compilation: "+*e.debugMessage, exception))
		}
		// Throwables.throwIfUnchecked
		switch exception.(type) {
		case jutil.IRuntimeException, jutil.IError:
			panic(exception)
		}
		panic(jutil.NewRuntimeExceptionWithCause(exception))
	}()

	if bootThread {
		e.compilerThread = jutil.Thread_CurrentThread()
		if dumpTraceReport {
			Tracer_InitCurrentThreadTrace()
		}
	}
	return callable()
}
