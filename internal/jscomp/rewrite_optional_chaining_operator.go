// Port of com/google/javascript/jscomp/RewriteOptionalChainingOperator.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var rewriteOptChainVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_SCRIPT)
	s.Add(rhino.Token_OPTCHAIN_GETPROP)
	s.Add(rhino.Token_OPTCHAIN_GETELEM)
	s.Add(rhino.Token_OPTCHAIN_CALL)
	return s
}()

// RewriteOptionalChainingOperator replaces the ES2020 `?.` operator with
// conditional (? :).
type RewriteOptionalChainingOperator struct {
	compiler                AbstractCompiler
	uniqueIdSupplier        *UniqueIdSupplier
	customTmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator
	callback                RewriteOptionalChainingOperator_TranspilationCallback
}

// NewRewriteOptionalChainingOperator is the constructor to be used for
// actual transpilation.
//
// Java: RewriteOptionalChainingOperator(AbstractCompiler)
func NewRewriteOptionalChainingOperator(compiler AbstractCompiler) *RewriteOptionalChainingOperator {
	return &RewriteOptionalChainingOperator{
		compiler:         compiler,
		uniqueIdSupplier: compiler.GetUniqueIdSupplier(),
	}
}

// NewRewriteOptionalChainingOperatorWithTmpVarNameCreator is the
// constructor for testing.
//
// Java: RewriteOptionalChainingOperator(AbstractCompiler, TmpVarNameCreator)
func NewRewriteOptionalChainingOperatorWithTmpVarNameCreator(
	compiler AbstractCompiler, tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator) *RewriteOptionalChainingOperator {
	return &RewriteOptionalChainingOperator{
		compiler:                compiler,
		customTmpVarNameCreator: jutil.CheckNotNull(tmpVarNameCreator),
	}
}

// Java: RewriteOptionalChainingOperator.process(Node, Node)
func (r *RewriteOptionalChainingOperator) Process(externs rhino.Node, root rhino.Node) {
	r.callback.initForProcess(r)
	NodeTraversal_TraverseRoots(r.compiler, &r.callback, externs, root)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(r.compiler, root, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
}

// RewriteOptionalChainingOperator_TranspilationCallback locates and
// transpiles all optional chains.
type RewriteOptionalChainingOperator_TranspilationCallback struct {
	outer           *RewriteOptionalChainingOperator
	rewriterBuilder *OptionalChainRewriter_Builder
	builderVal      OptionalChainRewriter_Builder
	currentInput    *CompilerInput
	boundTmpCreator OptionalChainRewriter_TmpVarNameCreator
	optChains       []*OptionalChainRewriter
	optChainsBuf    [4]*OptionalChainRewriter
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback()
func NewRewriteOptionalChainingOperator_TranspilationCallback(outer *RewriteOptionalChainingOperator) *RewriteOptionalChainingOperator_TranspilationCallback {
	c := &RewriteOptionalChainingOperator_TranspilationCallback{}
	c.initForProcess(outer)
	return c
}

func (c *RewriteOptionalChainingOperator_TranspilationCallback) initForProcess(outer *RewriteOptionalChainingOperator) {
	c.outer = outer
	c.builderVal = OptionalChainRewriter_Builder{
		compiler:   jutil.CheckNotNull(outer.compiler),
		astFactory: outer.compiler.CreateAstFactory(),
	}
	c.rewriterBuilder = &c.builderVal
	if outer.customTmpVarNameCreator != nil {
		c.rewriterBuilder.SetTmpVarNameCreator(outer.customTmpVarNameCreator)
	} else {
		if c.boundTmpCreator == nil {
			c.boundTmpCreator = c.nextTmpVarName
		}
		c.rewriterBuilder.SetTmpVarNameCreator(c.boundTmpCreator)
	}
	if c.optChains == nil {
		c.optChains = c.optChainsBuf[:0]
	} else {
		clear(c.optChains)
		c.optChains = c.optChains[:0]
	}
}

func (c *RewriteOptionalChainingOperator_TranspilationCallback) nextTmpVarName() string {
	return "$jscomp$optchain$tmp" + c.outer.uniqueIdSupplier.GetUniqueId(c.currentInput)
}

func (*RewriteOptionalChainingOperator_TranspilationCallback) visitTokens() *ast.TokenSet {
	return &rewriteOptChainVisitTokens
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *RewriteOptionalChainingOperator_TranspilationCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		// ensures that the pass early exits if script does not contain the feature
		if scriptFeatures != nil && !scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING) {
			return false
		}
		c.currentInput = t.GetInput()
	}
	return true
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback.visit(NodeTraversal, Node, Node)
func (c *RewriteOptionalChainingOperator_TranspilationCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsScript() {
		// We transpile all of the optional chains in a single script as a batch because,
		// rewriting changes the AST in ways that could interfere with traversal
		// if we change the chains as we visit them.
		if len(c.optChains) > 0 {
			for _, optChain := range c.optChains {
				optChain.Rewrite()
			}
			clear(c.optChains)
			c.optChains = c.optChains[:0]
		}
	} else if NodeUtil_IsEndOfFullOptChain(n) {
		c.optChains = append(c.optChains, c.rewriterBuilder.Build(n))
	}
}
