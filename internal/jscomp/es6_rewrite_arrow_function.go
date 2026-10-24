// Port of com/google/javascript/jscomp/Es6RewriteArrowFunction.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteArrowFunction converts ES6 arrow functions to standard anonymous
// ES3 functions.
type Es6RewriteArrowFunction struct {
	compiler   AbstractCompiler
	astFactory *AstFactory
	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose last element is the top (Deque.peek/push/pop).
	contextStack []*ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext
}

// Java: Es6RewriteArrowFunction.transpiledFeatures
var Es6RewriteArrowFunction_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_ARROW_FUNCTIONS)

// Java: Es6RewriteArrowFunction(AbstractCompiler)
func NewEs6RewriteArrowFunction(compiler AbstractCompiler) *Es6RewriteArrowFunction {
	e := &Es6RewriteArrowFunction{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.contextStack = nil
	return e
}

// Java: Es6RewriteArrowFunction.process(Node, Node)
func (e *Es6RewriteArrowFunction) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, Es6RewriteArrowFunction_transpiledFeatures, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6RewriteArrowFunction_transpiledFeatures)
}

// contextStackPeek is Deque.peek(): the top, or null when empty.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (e *Es6RewriteArrowFunction) contextStackPeek() *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	if len(e.contextStack) == 0 {
		return nil
	}
	return e.contextStack[len(e.contextStack)-1]
}

// Java: Es6RewriteArrowFunction.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteArrowFunction) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		e.contextStack = append(e.contextStack, e.ContextForScript(n, t.GetInput()))
	case rhino.Token_FUNCTION:
		if !n.IsArrowFunction() {
			e.contextStack = append(e.contextStack, e.ContextForFunction(n, t.GetInput()))
		}
	case rhino.Token_SUPER:
		context := jutil.CheckNotNull(e.contextStackPeek())
		// super(...) within a constructor.
		if context.isConstructor && parent.IsCall() && parent.GetFirstChild() == n {
			context.lastSuperStatement = e.GetEnclosingStatement(parent, context.scopeBody)
		}
	default:
	}
	return true
}

// GetEnclosingStatement returns the statement Node that is a child of block
// and contains n.
//
// Java: Es6RewriteArrowFunction.getEnclosingStatement(Node, Node)
func (e *Es6RewriteArrowFunction) GetEnclosingStatement(n rhino.Node, block rhino.Node) rhino.Node {
	for rhino.CheckNotNull(n.GetParent()) != block {
		n = rhino.CheckNotNull(NodeUtil_GetEnclosingStatement(n.GetParent()))
	}
	return n
}

// Java: Es6RewriteArrowFunction.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteArrowFunction) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	context := e.contextStackPeek()

	if n.IsArrowFunction() {
		e.VisitArrowFunction(t, n, jutil.CheckNotNull(context))
	} else if context != nil && context.scopeBody == n {
		e.contextStack = e.contextStack[:len(e.contextStack)-1]
		context.AddVarDeclarations(e.compiler, e.astFactory, t.GetCurrentScript())
	}
}

// Java: Es6RewriteArrowFunction.visitArrowFunction(NodeTraversal, Node, ThisAndArgumentsContext)
func (e *Es6RewriteArrowFunction) VisitArrowFunction(t *NodeTraversal, n rhino.Node, context *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) {
	n.SetIsArrowFunction(false)

	body := n.GetLastChild()
	jutil.CheckStateMsg(body.IsBlock(), "Arrow function body must be a block after normalization")

	updater :=
		NewThisAndArgumentsReferenceUpdater(e.compiler, context, e.astFactory)
	NodeTraversal_Traverse(e.compiler, body, updater)

	t.ReportCodeChange()
}

// Java: Es6RewriteArrowFunction.contextForFunction(Node, CompilerInput)
func (e *Es6RewriteArrowFunction) ContextForFunction(functionNode rhino.Node, input *CompilerInput) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	scopeBody := functionNode.GetLastChild()
	return NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
		scopeBody,
		NodeUtil_IsEs6Constructor(functionNode),
		e.compiler.GetUniqueIdSupplier().GetUniqueId(input))
}

// Java: Es6RewriteArrowFunction.contextForScript(Node, CompilerInput)
func (e *Es6RewriteArrowFunction) ContextForScript(scriptNode rhino.Node, input *CompilerInput) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	return NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
		scriptNode /* isConstructor= */, false, e.compiler.GetUniqueIdSupplier().GetUniqueId(input))
}
