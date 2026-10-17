// Port of com/google/javascript/refactoring/NodeMetadata.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// NodeMetadata is the class that holds metadata (or meta objects) for use by
// JsFlume that aren't contained within the Node itself.
type NodeMetadata struct {
	compiler AbstractCompiler

	traversal *NodeTraversal
}

// Java: NodeMetadata(AbstractCompiler)
func NewNodeMetadata(compiler AbstractCompiler) *NodeMetadata {
	return newNodeMetadataWithTraversal(compiler, nil)
}

// Java: NodeMetadata(AbstractCompiler, NodeTraversal)
func newNodeMetadataWithTraversal(compiler AbstractCompiler, traversal *NodeTraversal) *NodeMetadata {
	return &NodeMetadata{compiler: compiler, traversal: traversal}
}

// Java: NodeMetadata.fromTraversal(NodeTraversal)
func NodeMetadata_FromTraversal(traversal *NodeTraversal) *NodeMetadata {
	return newNodeMetadataWithTraversal(traversal.GetCompiler(), traversal)
}

// Java: NodeMetadata.getCompiler()
func (m *NodeMetadata) GetCompiler() AbstractCompiler {
	return m.compiler
}

// Java: NodeMetadata.getTraversal()
func (m *NodeMetadata) GetTraversal() *NodeTraversal {
	return m.traversal
}
