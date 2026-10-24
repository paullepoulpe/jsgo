// Port of com/google/javascript/jscomp/RewriteOptionalChainingOperator.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteOptionalChainingOperator replaces the ES2020 `?.` operator with
// conditional (? :).
type RewriteOptionalChainingOperator struct {
	compiler AbstractCompiler

	// Produces the object responsible for creating temporary variable names for a given
	// CompilerInput.
	getTmpVarNameCreatorForInput func(input *CompilerInput) OptionalChainRewriter_TmpVarNameCreator
}

// NewRewriteOptionalChainingOperator is the constructor to be used for
// actual transpilation.
//
// Java: RewriteOptionalChainingOperator(AbstractCompiler)
func NewRewriteOptionalChainingOperator(compiler AbstractCompiler) *RewriteOptionalChainingOperator {
	r := &RewriteOptionalChainingOperator{}
	r.compiler = compiler
	// Temporary variable names are generated based on the compiler input they will go into and use
	// a recognizable prefix to aid debugging.
	uniqueIdSupplier := compiler.GetUniqueIdSupplier()
	r.getTmpVarNameCreatorForInput =
		func(input *CompilerInput) OptionalChainRewriter_TmpVarNameCreator {
			return OptionalChainRewriter_TmpVarNameCreatorFunc(func() string {
				return "$jscomp$optchain$tmp" + uniqueIdSupplier.GetUniqueId(input)
			})
		}
	return r
}

// NewRewriteOptionalChainingOperatorWithTmpVarNameCreator is the
// constructor for testing.
//
// Java: RewriteOptionalChainingOperator(AbstractCompiler, TmpVarNameCreator)
func NewRewriteOptionalChainingOperatorWithTmpVarNameCreator(
	compiler AbstractCompiler, tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator) *RewriteOptionalChainingOperator {
	r := &RewriteOptionalChainingOperator{}
	r.compiler = compiler
	// The test provides a simple variable name creator that makes for more readable test cases.
	r.getTmpVarNameCreatorForInput = func(input *CompilerInput) OptionalChainRewriter_TmpVarNameCreator {
		return tmpVarNameCreator
	}
	return r
}

// Java: RewriteOptionalChainingOperator.process(Node, Node)
func (r *RewriteOptionalChainingOperator) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(r.compiler, NewRewriteOptionalChainingOperator_TranspilationCallback(r), externs, root)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(r.compiler, root, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
}

// RewriteOptionalChainingOperator_TranspilationCallback locates and
// transpiles all optional chains.
type RewriteOptionalChainingOperator_TranspilationCallback struct {
	outer           *RewriteOptionalChainingOperator
	rewriterBuilder *OptionalChainRewriter_Builder
	optChains       *jutil.ArrayList[*OptionalChainRewriter]
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback()
func NewRewriteOptionalChainingOperator_TranspilationCallback(outer *RewriteOptionalChainingOperator) *RewriteOptionalChainingOperator_TranspilationCallback {
	c := &RewriteOptionalChainingOperator_TranspilationCallback{outer: outer}
	c.rewriterBuilder = OptionalChainRewriter_Builder_(outer.compiler)
	c.optChains = jutil.NewArrayList[*OptionalChainRewriter]()
	return c
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *RewriteOptionalChainingOperator_TranspilationCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		// Set the TmpVarNameCreator to be used when rewriting optional chains in this script.
		c.rewriterBuilder.SetTmpVarNameCreator(c.outer.getTmpVarNameCreatorForInput(t.GetInput()))
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		// ensures that the pass early exits if script does not contain the feature
		return scriptFeatures == nil || scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	}
	return true
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback.visit(NodeTraversal, Node, Node)
func (c *RewriteOptionalChainingOperator_TranspilationCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsEndOfFullOptChain(n) {
		c.optChains.Add(c.rewriterBuilder.Build(n))
	} else if n.IsScript() {
		// We transpile all of the optional chains in a single script as a batch because,
		// rewriting changes the AST in ways that could interfere with traversal
		// if we change the chains as we visit them.
		if !c.optChains.IsEmpty() {
			for optChain := range c.optChains.All() {
				optChain.Rewrite()
			}
			c.optChains.Clear()
		}
	}
}
