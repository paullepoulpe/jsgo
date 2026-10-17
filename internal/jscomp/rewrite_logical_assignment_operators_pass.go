// Port of com/google/javascript/jscomp/RewriteLogicalAssignmentOperatorsPass.java @ closure-compiler v20261015

/*
 * Copyright 2021 The Closure Compiler Authors.
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

// RewriteLogicalAssignmentOperatorsPass replaces the ES2020 `||=`, `&&=`,
// and `??=` operators.
type RewriteLogicalAssignmentOperatorsPass struct {
	compiler                                AbstractCompiler
	rewriteLogicalAssignmentOperatorsHelper *RewriteLogicalAssignmentOperatorsHelper
}

// Java: RewriteLogicalAssignmentOperatorsPass(AbstractCompiler)
func NewRewriteLogicalAssignmentOperatorsPass(compiler AbstractCompiler) *RewriteLogicalAssignmentOperatorsPass {
	r := &RewriteLogicalAssignmentOperatorsPass{}
	r.compiler = compiler
	r.rewriteLogicalAssignmentOperatorsHelper =
		NewRewriteLogicalAssignmentOperatorsHelper(
			compiler, compiler.CreateAstFactory(), compiler.GetUniqueIdSupplier())
	return r
}

// Java: RewriteLogicalAssignmentOperatorsPass.process(Node, Node)
func (r *RewriteLogicalAssignmentOperatorsPass) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, r)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(
		r.compiler, root, parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
}

// Java: RewriteLogicalAssignmentOperatorsPass.shouldTraverse(NodeTraversal, Node, Node)
func (r *RewriteLogicalAssignmentOperatorsPass) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsScript() {
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		// runs only if logical assignment present
		return scriptFeatures == nil || scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	}
	return true
}

// Java: RewriteLogicalAssignmentOperatorsPass.visit(NodeTraversal, Node, Node)
func (r *RewriteLogicalAssignmentOperatorsPass) Visit(t *NodeTraversal, logicalAssignment *rhino.Node, parent *rhino.Node) {
	switch logicalAssignment.GetToken() {
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		r.rewriteLogicalAssignmentOperatorsHelper.VisitLogicalAssignmentOperator(
			t, logicalAssignment)
	default:
	}
}
