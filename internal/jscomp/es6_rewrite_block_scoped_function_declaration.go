// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedFunctionDeclaration.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteBlockScopedFunctionDeclaration rewrites block-scoped function
// declarations as "let"s. This pass must happen before
// Es6RewriteBlockScopedDeclaration, which rewrites "let" to "var".
type Es6RewriteBlockScopedFunctionDeclaration struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.transpiledFeatures
var es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)

// Java: Es6RewriteBlockScopedFunctionDeclaration(AbstractCompiler)
func NewEs6RewriteBlockScopedFunctionDeclaration(compiler AbstractCompiler) *Es6RewriteBlockScopedFunctionDeclaration {
	e := &Es6RewriteBlockScopedFunctionDeclaration{}
	e.compiler = compiler
	return e
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.process(Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures)
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsFunction() &&
		!parent.IsNil() &&
		parent.IsBlock() &&
		!parent.GetParent().IsFunction() {
		// Only consider declarations (all expressions have non-block parents) that are not directly
		// within a function or top-level.
		e.VisitBlockScopedFunctionDeclaration(t, n, parent)
	}
}

// VisitBlockScopedFunctionDeclaration rewrites the function declaration
// from:
//
//	function f() {}
//	FUNCTION
//	  NAME x
//	  PARAM_LIST
//	  BLOCK
//
// to
//
//	let f = function() {};
//	LET
//	  NAME f
//	    FUNCTION
//	      NAME (w/ empty string)
//	      PARAM_LIST
//	      BLOCK
//
// This is similar to Normalize.NormalizeStatements#rewriteFunctionDeclaration
// but rewrites to "let" instead of "var".
//
// Java: Es6RewriteBlockScopedFunctionDeclaration.visitBlockScopedFunctionDeclaration(NodeTraversal, Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) VisitBlockScopedFunctionDeclaration(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Prepare a spot for the function.
	oldNameNode := n.GetFirstChild()
	fnNameNode := oldNameNode.CloneNode()
	let := rhino.IR_Declaration(fnNameNode, rhino.Token_LET).Srcref(n)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)

	// Prepare the function.
	oldNameNode.SetString("")
	e.compiler.ReportChangeToEnclosingScope(oldNameNode)

	// Move the function to the front of the parent.
	n.Detach()
	parent.AddChildToFront(let)
	e.compiler.ReportChangeToEnclosingScope(let)
	fnNameNode.AddChildToFront(n)
}
