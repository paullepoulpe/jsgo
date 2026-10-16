// Port of com/google/javascript/jscomp/PrebuildDependencyInfo.java @ closure-compiler v20260415

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

package jscomp

import (
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PrebuildDependencyInfo is a helper class to prebuild DependencyInfo from
// a list of CompilerInput. Dependency info (requires, annotations, etc.) is
// parsed and generated the first time we try to get it for an input. Get
// them all now using multiple threads, so they'll be parsed in parallel and
// already available during the rest of the compilation.
type PrebuildDependencyInfo struct {
	numParallelThreads int
}

// Java: PrebuildDependencyInfo(int)
func NewPrebuildDependencyInfo(numParalleThreads int) *PrebuildDependencyInfo {
	return &PrebuildDependencyInfo{numParallelThreads: numParalleThreads}
}

// Prebuild gets the dependency info of every input.
//
// Go-only: language-adapter — the port is sequential (PORTING.md §5). Java
// submits input::getDependencyInfo for every input to a pool of
// numParallelThreads daemon threads ("jscompiler-PrebuildDependencyInfo",
// with CompilerExecutor.COMPILER_STACK_SIZE stacks), shuts the pool down
// (queued tasks still run) and waits with Futures.allAsList(futureList).get(),
// which throws an ExecutionException wrapping the failure of a task;
// prebuild rethrows it as new RuntimeException(e). Here the tasks run
// inline in submission order, every task runs even after one fails (as in
// the pool), and the first failure is rethrown the same way.
//
// Java: PrebuildDependencyInfo.prebuild(Iterable<CompilerInput>)
func (p *PrebuildDependencyInfo) Prebuild(allInputs iter.Seq[*CompilerInput]) {
	var failure jutil.Throwable
	// TODO(moz): Support canceling all parsing on the first halting error
	for input := range allInputs {
		jutil.Try(func() {
			input.GetDependencyInfo()
		}, func(t jutil.Throwable) {
			if failure == nil {
				failure = t
			}
		})
	}

	if failure != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
	}
}
