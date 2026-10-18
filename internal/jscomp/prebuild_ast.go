// Port of com/google/javascript/jscomp/PrebuildAst.java @ closure-compiler v20261015

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

package jscomp

import (
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PrebuildAst is a helper class to prebuild ASTs from a list of
// CompilerInput. Inputs are parsed into ASTs the first time we try to get
// them. Get them all now using multiple threads, so they'll be parsed in
// parallel and already available during the rest of the compilation.
type PrebuildAst struct {
	compiler           AbstractCompiler
	numParallelThreads int
}

// Java: PrebuildAst(AbstractCompiler, int)
func NewPrebuildAst(compiler AbstractCompiler, numParalleThreads int) *PrebuildAst {
	return &PrebuildAst{compiler: compiler, numParallelThreads: numParalleThreads}
}

// Prebuild gets the AST root of every input.
//
// Go-only: language-adapter — independent files are parsed by a bounded worker
// pool. Compiler state and diagnostics are published in input order after the
// workers finish. Custom compilers, shared sources and TypedAST loading use the
// sequential compatibility path. Every task completes even if another fails;
// the first failure in input order has Java's ExecutionException wrapping.
//
// Java: PrebuildAst.prebuild(Iterable<CompilerInput>)
func (p *PrebuildAst) Prebuild(allInputs iter.Seq[*CompilerInput]) {
	if p.numParallelThreads > 1 {
		inputs := slices.Collect(allInputs)
		if p.prebuildParallel(inputs) {
			return
		}
		allInputs = slices.Values(inputs)
	}

	var failure jutil.Throwable
	// TODO(moz): Support canceling all parsing on the first halting error
	for input := range allInputs {
		jutil.Try(func() {
			input.GetAstRoot(p.compiler)
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
