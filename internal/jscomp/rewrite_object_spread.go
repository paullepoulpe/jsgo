// Port of com/google/javascript/jscomp/RewriteObjectSpread.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) to read the
// object literal it rewrites: its members (ast.ObjectMember), which of them
// are spreads, and the expression of each spread (spread.Expr()). The
// Object.assign call and the object literals it receives come from the
// AstFactory (which sets their colors), so they stay rhino.Nodes.

// RewriteObjectSpread converts object spread to valid ES2017 code.
//
// Currently this class converts Object spread properties as documented in
// tc39. https://github.com/tc39/proposal-object-rest-spread. For example:
//
// `const bar = {a: 1, ...foo};`
//
// Note that object rest is handled by Es6RewriteDestructuring.
type RewriteObjectSpread struct {
	compiler AbstractCompiler

	astFactory *AstFactory
	namespace  rhino.StaticScope
}

// Java: RewriteObjectSpread.transpiledFeatures
var RewriteObjectSpread_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)

// Java: RewriteObjectSpread(AbstractCompiler)
func NewRewriteObjectSpread(compiler AbstractCompiler) *RewriteObjectSpread {
	r := &RewriteObjectSpread{}
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	r.namespace = r.compiler.GetTranspilationNamespace()
	return r
}

// Java: RewriteObjectSpread.process(Node, Node)
func (r *RewriteObjectSpread) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(r.compiler, root, RewriteObjectSpread_transpiledFeatures, r)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(r.compiler, root, RewriteObjectSpread_transpiledFeatures)
}

// Java: RewriteObjectSpread.shouldTraverse(NodeTraversal, Node, Node)
func (r *RewriteObjectSpread) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Java: RewriteObjectSpread.visit(NodeTraversal, Node, Node)
func (r *RewriteObjectSpread) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if obj, ok := ast.AsObjectLit(n); ok {
		r.VisitObject(obj)
	}
}

// Java: RewriteObjectSpread.visitObject(Node)
func (r *RewriteObjectSpread) VisitObject(obj ast.ObjectLit) {
	for member := range obj.Members() {
		if _, ok := member.AsObjectSpread(); ok {
			r.VisitObjectWithSpread(obj)
			return
		}
	}
}

// Convert '{first: b, c, ...spread, d: e, last}' to:
//
// Object.assign({}, {first:b, c}, spread, {d:e, last});
//
// Java: RewriteObjectSpread.visitObjectWithSpread(Node)
func (r *RewriteObjectSpread) VisitObjectWithSpread(obj ast.ObjectLit) {
	// Add an empty target object literal so changes made by Object.assign will not affect any other
	// variables.
	result :=
		r.astFactory.CreateObjectDotAssignCall(r.namespace, AstFactory_Type_(obj.Node()), r.astFactory.CreateObjectLit())

	// An indicator whether the current last thing in the param list is an object literal to which
	// properties may be added.  Initialized to null since nothing should be added to the empty
	// object literal in first position of the param list.
	var trailingObjectLiteral rhino.Node = rhino.NilNode

	// The next member is read before the current one moves.
	next, more := obj.FirstMember()
	for more {
		member := next
		next, more = member.Next()
		if spread, ok := member.AsObjectSpread(); ok {
			// Add the object directly to the param list.
			spreaded := spread.Expr().Detach()
			result.AddChildToBack(spreaded.Node())

			// Properties should not be added to the trailing object.
			trailingObjectLiteral = rhino.NilNode
		} else {
			if trailingObjectLiteral.IsNil() {
				// Add a new object to which properties may be added.
				trailingObjectLiteral = r.astFactory.CreateObjectLit()
				result.AddChildToBack(trailingObjectLiteral)
			}
			// Add the property to the object literal.
			trailingObjectLiteral.AddChildToBack(member.Detach().Node())
		}
	}

	result.SrcrefTreeIfMissing(obj.Node())
	obj.Node().ReplaceWith(result)
	r.compiler.ReportChangeToEnclosingScope(result)
}
