// Port of com/google/javascript/jscomp/Es6RewriteArrowFunction.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it reads
// shapes: the functions it rewrites or opens a this/arguments context for
// (ast.Function), the body of an arrow function (a block after
// normalization: arrow.Body().AsBlock()) and the callee of a super(...) call.
// ShouldTraverse still switches on the token: it runs on every node and
// hands typed nodes to the helpers.

// Es6RewriteArrowFunction converts ES6 arrow functions to standard anonymous
// ES3 functions.
type Es6RewriteArrowFunction struct {
	compiler   AbstractCompiler
	astFactory *AstFactory
	updater    ThisAndArgumentsReferenceUpdater
	// Contiguous value stack (top is last) avoiding per-function/per-script
	// heap allocations during traversal.
	contextStack    []ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext
	contextStackBuf [8]ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext
}

var es6RewriteArrowVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FUNCTION)
	s.Add(rhino.Token_SCRIPT)
	s.Add(rhino.Token_BLOCK)
	return s
}()

// Java: Es6RewriteArrowFunction.transpiledFeatures
var Es6RewriteArrowFunction_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_ARROW_FUNCTIONS)

// Java: Es6RewriteArrowFunction(AbstractCompiler)
func NewEs6RewriteArrowFunction(compiler AbstractCompiler) *Es6RewriteArrowFunction {
	astFactory := compiler.CreateAstFactory()
	e := &Es6RewriteArrowFunction{
		compiler:   compiler,
		astFactory: astFactory,
		updater: ThisAndArgumentsReferenceUpdater{
			compiler:   compiler,
			astFactory: astFactory,
		},
	}
	e.contextStack = e.contextStackBuf[:0]
	return e
}

func (*Es6RewriteArrowFunction) visitTokens() *ast.TokenSet {
	return &es6RewriteArrowVisitTokens
}

// Java: Es6RewriteArrowFunction.process(Node, Node)
func (e *Es6RewriteArrowFunction) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, Es6RewriteArrowFunction_transpiledFeatures, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6RewriteArrowFunction_transpiledFeatures)
}

// contextStackPeek is Deque.peek(): pointer to the top element, or nil when empty.
func (e *Es6RewriteArrowFunction) contextStackPeek() *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	if len(e.contextStack) == 0 {
		return nil
	}
	return &e.contextStack[len(e.contextStack)-1]
}

// Java: Es6RewriteArrowFunction.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteArrowFunction) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		e.contextStack = append(e.contextStack, e.makeScriptContext(n, t.GetInput()))
	case rhino.Token_FUNCTION:
		if fn, ok := ast.AsFunction(n); ok && !fn.IsArrow() {
			e.contextStack = append(e.contextStack, e.makeFunctionContext(fn, t.GetInput()))
		}
	case rhino.Token_SUPER:
		context := jutil.CheckNotNull(e.contextStackPeek())
		// super(...) within a constructor.
		call, isCall := ast.AsCall(parent)
		if context.isConstructor && isCall && call.Callee().Is(n) {
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

	if fn, ok := ast.AsFunction(n); ok && fn.IsArrow() {
		e.VisitArrowFunction(t, fn, jutil.CheckNotNull(context))
	} else if context != nil && context.scopeBody == n {
		popped := *context
		*context = ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext{}
		e.contextStack = e.contextStack[:len(e.contextStack)-1]
		popped.AddVarDeclarations(e.compiler, e.astFactory, t.GetCurrentScript())
	}
}

// Java: Es6RewriteArrowFunction.visitArrowFunction(NodeTraversal, Node, ThisAndArgumentsContext)
func (e *Es6RewriteArrowFunction) VisitArrowFunction(t *NodeTraversal, arrow ast.Function, context *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext) {
	arrow.SetArrow(false)

	body, ok := arrow.Body().AsBlock()
	jutil.CheckStateMsg(ok, "Arrow function body must be a block after normalization")

	e.updater.context = context
	NodeTraversal_Traverse(e.compiler, body.Node(), &e.updater)
	e.updater.context = nil

	t.ReportCodeChange()
}

func (e *Es6RewriteArrowFunction) makeFunctionContext(fn ast.Function, input *CompilerInput) ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	return makeThisAndArgumentsContext(
		fn.Body().Node(),
		NodeUtil_IsEs6Constructor(fn.Node()),
		e.compiler.GetUniqueIdSupplier().GetUniqueId(input))
}

func (e *Es6RewriteArrowFunction) makeScriptContext(scriptNode rhino.Node, input *CompilerInput) ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	return makeThisAndArgumentsContext(
		scriptNode, false, e.compiler.GetUniqueIdSupplier().GetUniqueId(input))
}

// Java: Es6RewriteArrowFunction.contextForFunction(Node, CompilerInput)
func (e *Es6RewriteArrowFunction) ContextForFunction(fn ast.Function, input *CompilerInput) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	c := e.makeFunctionContext(fn, input)
	return &c
}

// Java: Es6RewriteArrowFunction.contextForScript(Node, CompilerInput)
func (e *Es6RewriteArrowFunction) ContextForScript(scriptNode rhino.Node, input *CompilerInput) *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	c := e.makeScriptContext(scriptNode, input)
	return &c
}
