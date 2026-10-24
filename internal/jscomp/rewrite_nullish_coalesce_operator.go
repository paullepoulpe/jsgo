// Port of com/google/javascript/jscomp/RewriteNullishCoalesceOperator.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteNullishCoalesceOperator replaces the ES2020 `??` operator with
// conditional (?:).
type RewriteNullishCoalesceOperator struct {
	compiler         AbstractCompiler
	astFactory       *AstFactory
	uniqueIdSupplier *UniqueIdSupplier
}

// Java: RewriteNullishCoalesceOperator.TEMP_VAR_NAME_PREFIX
const RewriteNullishCoalesceOperator_TEMP_VAR_NAME_PREFIX = "$jscomp$nullish$tmp"

// Java: RewriteNullishCoalesceOperator(AbstractCompiler)
func NewRewriteNullishCoalesceOperator(compiler AbstractCompiler) *RewriteNullishCoalesceOperator {
	r := &RewriteNullishCoalesceOperator{}
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	r.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	return r
}

// Java: RewriteNullishCoalesceOperator.process(Node, Node)
func (r *RewriteNullishCoalesceOperator) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, r)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(r.compiler, root, parser.FeatureSet_Feature_NULL_COALESCE_OP)
}

// Java: RewriteNullishCoalesceOperator.shouldTraverse(NodeTraversal, Node, Node)
func (r *RewriteNullishCoalesceOperator) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		return scriptFeatures == nil || scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP)
	}
	return true
}

// Java: RewriteNullishCoalesceOperator.visit(NodeTraversal, Node, Node)
func (r *RewriteNullishCoalesceOperator) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsNullishCoalesce() {
		r.VisitNullishCoalesce(t, n)
		t.ReportCodeChange()
	}
}

// Java: RewriteNullishCoalesceOperator.visitNullishCoalesce(NodeTraversal, Node)
func (r *RewriteNullishCoalesceOperator) VisitNullishCoalesce(t *NodeTraversal, n rhino.Node) {
	// a() ?? b()
	// let temp;
	// ((temp = a()) != null) ? temp : b()
	tempVarName := RewriteNullishCoalesceOperator_TEMP_VAR_NAME_PREFIX + r.uniqueIdSupplier.GetUniqueId(t.GetInput())
	enclosingStatement := NodeUtil_GetEnclosingStatement(n)

	left := n.RemoveFirstChild()
	right := n.GetLastChild().Detach()

	let := r.astFactory.CreateSingleLetNameDeclaration(tempVarName)
	assignName := r.astFactory.CreateName(tempVarName, AstFactory_Type_(left))
	assign := r.astFactory.CreateAssign(assignName, left)
	ne := r.astFactory.CreateNe(assign, r.astFactory.CreateNull())
	hookName := r.astFactory.CreateName(tempVarName, AstFactory_Type_(left))
	hook := r.astFactory.CreateHook(ne, hookName, right)

	let.SrcrefTreeIfMissing(left)
	assignName.SrcrefTreeIfMissing(left)
	assign.SrcrefTreeIfMissing(left)
	ne.SrcrefTreeIfMissing(left)
	hookName.SrcrefTreeIfMissing(left)

	let.InsertBefore(enclosingStatement)
	n.ReplaceWith(hook)

	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, r.compiler)
	r.compiler.ReportChangeToEnclosingScope(hook)
}
