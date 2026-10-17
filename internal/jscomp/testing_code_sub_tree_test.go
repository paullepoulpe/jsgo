// Port of com/google/javascript/jscomp/testing/CodeSubTree.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CodeSubTree represents a subtree of the output from a compilation.
type CodeSubTree struct {
	rootNode *rhino.Node
}

// Java: CodeSubTree(Node)
func newCodeSubTree(rootNode *rhino.Node) *CodeSubTree {
	return &CodeSubTree{rootNode: rootNode}
}

// Java: CodeSubTree.getRootNode()
func (t *CodeSubTree) GetRootNode() *rhino.Node {
	return t.rootNode
}

// FindClassDefinition returns the SubTree rooted at the first class
// definition found with the given name.
//
// Java: CodeSubTree.findClassDefinition(String)
func (t *CodeSubTree) FindClassDefinition(wantedClassName string) *CodeSubTree {
	classNode := CodeSubTree_FindFirstNode(
		t.rootNode, func(node *rhino.Node) bool {
			if !node.IsClass() {
				return false
			}
			name := NodeUtil_GetName(node)
			return name != nil && wantedClassName == *name
		})
	return newCodeSubTree(classNode)
}

// CodeSubTree_FindClassDefinition returns a CodeSubTree for the first
// definition of the given class name in the output or externs from the last
// compile.
//
// Java: CodeSubTree.findClassDefinition(AbstractCompiler, String)
func CodeSubTree_FindClassDefinition(compiler AbstractCompiler, wantedClassName string) *CodeSubTree {
	return newCodeSubTree(compiler.GetRoot()).FindClassDefinition(wantedClassName)
}

// FindMethodDefinition returns the first class method definiton found with
// the given name.
//
// Java: CodeSubTree.findMethodDefinition(String)
func (t *CodeSubTree) FindMethodDefinition(wantedMethodName string) *CodeSubTree {
	methodDefinitionNode := CodeSubTree_FindFirstNode(
		t.rootNode,
		func(node *rhino.Node) bool {
			return node.IsMemberFunctionDef() && wantedMethodName == node.GetString()
		})

	return newCodeSubTree(methodDefinitionNode)
}

// CodeSubTree_FindFunctionDefinition returns the first function method
// definition found with the given name.
//
// Java: CodeSubTree.findFunctionDefinition(AbstractCompiler, String)
func CodeSubTree_FindFunctionDefinition(
	compiler AbstractCompiler, wantedMethodName string) *CodeSubTree {
	functionDefinitionNode := CodeSubTree_FindFirstNode(
		compiler.GetRoot(),
		func(node *rhino.Node) bool {
			return node.IsFunction() &&
				node.GetFirstChild().IsName() &&
				wantedMethodName == node.GetFirstChild().GetString()
		})

	return newCodeSubTree(functionDefinitionNode)
}

// FindMatchingQNameReferences executes an action for every instance of a
// given qualified name.
//
// Java: CodeSubTree.findMatchingQNameReferences(String)
func (t *CodeSubTree) FindMatchingQNameReferences(wantedQName string) []*rhino.Node {
	return CodeSubTree_FindNodesAllowEmpty(t.rootNode, func(node *rhino.Node) bool { return node.MatchesQualifiedName(wantedQName) })
}

// CodeSubTree_FindNodesAllowEmpty returns a list of all Nodes matching the
// given predicate starting at the given root.
//
// Java: CodeSubTree.findNodesAllowEmpty(Node, Predicate<Node>)
func CodeSubTree_FindNodesAllowEmpty(rootNode *rhino.Node, predicate func(*rhino.Node) bool) []*rhino.Node {
	listBuilder := []*rhino.Node{}
	NodeUtil_VisitPreOrder(
		rootNode,
		NodeUtil_VisitorFunc(func(node *rhino.Node) {
			if predicate(node) {
				listBuilder = append(listBuilder, node)
			}
		}))
	return listBuilder
}

// CodeSubTree_FindNodesNonEmpty returns a list of all Nodes matching the
// given predicate starting at the given root.
//
// Java: CodeSubTree.findNodesNonEmpty(Node, Predicate<Node>)
func CodeSubTree_FindNodesNonEmpty(rootNode *rhino.Node, predicate func(*rhino.Node) bool) []*rhino.Node {
	results := CodeSubTree_FindNodesAllowEmpty(rootNode, predicate)
	jutil.CheckStateMsg(len(results) != 0, "no nodes found")
	return results
}

// CodeSubTree_FindFirstNode returns the shallowest and earliest of all
// Nodes matching the given predicate starting at the given root.
//
// Throws an exception if none found.
//
// Java: CodeSubTree.findFirstNode(Node, Predicate<Node>)
func CodeSubTree_FindFirstNode(rootNode *rhino.Node, predicate func(*rhino.Node) bool) *rhino.Node {
	allMatchingNodes := CodeSubTree_FindNodesNonEmpty(rootNode, predicate)
	return allMatchingNodes[0]
}
