// Port of com/google/javascript/jscomp/PolymerPassFindExterns.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const polymerPassFindExterns_POLYMER_ELEMENT_NAME = "PolymerElement"

// PolymerPassFindExterns finds the externs for the PolymerElement base class
// and all of its properties in the externs.
//
// Java: final class PolymerPassFindExterns implements NodeTraversal.Callback
type PolymerPassFindExterns struct {
	polymerElementExterns *rhino.Node
	polymerElementProps   []*rhino.Node
}

// Java: PolymerPassFindExterns()
func NewPolymerPassFindExterns() *PolymerPassFindExterns {
	f := &PolymerPassFindExterns{}
	f.polymerElementProps = nil // ImmutableList.builder()
	return f
}

// Java: PolymerPassFindExterns.shouldTraverse(NodeTraversal, Node, Node)
func (f *PolymerPassFindExterns) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	/*
	 * The externs subtree also includes @typeSummary files, which are slightly different from
	 * normal externs. The Polymer externs may exist in either file type. However,
	 * some @typeSummaries are goog.modules, which definitely don't contain Polymer externs, and
	 * which we don't want to search.
	 */
	return !n.IsModuleBody()
}

// Java: PolymerPassFindExterns.visit(NodeTraversal, Node, Node)
func (f *PolymerPassFindExterns) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if polymerPassFindExterns_isPolymerElementExterns(n) {
		f.polymerElementExterns = n
	} else if polymerPassFindExterns_isPolymerElementPropExpr(n) {
		f.polymerElementProps = append(f.polymerElementProps, n)
	}
}

// GetPolymerElementProps returns all properties of the PolymerElement class
// in externs.
//
// Java: PolymerPassFindExterns.getPolymerElementProps()
func (f *PolymerPassFindExterns) GetPolymerElementProps() []*rhino.Node {
	// ImmutableList.Builder.build() snapshots the elements added so far.
	return append([]*rhino.Node(nil), f.polymerElementProps...)
}

// GetPolymerElementExterns returns the root Node PolymerElement class in
// externs.
//
// Java: PolymerPassFindExterns.getPolymerElementExterns()
func (f *PolymerPassFindExterns) GetPolymerElementExterns() *rhino.Node {
	return f.polymerElementExterns
}

// polymerPassFindExterns_isPolymerElementExterns returns whether the node is
// the declaration of PolymerElement.
//
// Java: PolymerPassFindExterns.isPolymerElementExterns(Node)
func polymerPassFindExterns_isPolymerElementExterns(value *rhino.Node) bool {
	// Explicitly allow var, let or const.
	return NodeUtil_IsNameDeclaration(value) &&
		value.GetFirstChild().MatchesQualifiedName(polymerPassFindExterns_POLYMER_ELEMENT_NAME)
}

// polymerPassFindExterns_isPolymerElementPropExpr returns whether the node is
// an expression result of an assignment to a property of PolymerElement.
//
// Java: PolymerPassFindExterns.isPolymerElementPropExpr(Node)
func polymerPassFindExterns_isPolymerElementPropExpr(value *rhino.Node) bool {
	return value != nil &&
		value.IsExprResult() &&
		value.GetFirstFirstChild() != nil &&
		value.GetFirstFirstChild().IsGetProp() &&
		value.GetFirstFirstChild().IsQualifiedName() &&
		NodeUtil_GetRootOfQualifiedName(
			value.GetFirstFirstChild()).MatchesQualifiedName(polymerPassFindExterns_POLYMER_ELEMENT_NAME)
}
