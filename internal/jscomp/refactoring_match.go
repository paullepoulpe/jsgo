// Port of com/google/javascript/refactoring/Match.java @ closure-compiler v20260415

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

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Match is an object that contains the information for a given match.
type Match struct {
	node     *rhino.Node
	metadata *NodeMetadata
}

// Java: Match(Node, NodeMetadata)
func NewMatch(node *rhino.Node, metadata *NodeMetadata) *Match {
	return &Match{node: node, metadata: metadata}
}

// GetNode returns the node that matched the given conditions.
//
// Java: Match.getNode()
func (m *Match) GetNode() *rhino.Node {
	return m.node
}

// GetMetadata returns the metadata for this match.
//
// Java: Match.getMetadata()
func (m *Match) GetMetadata() *NodeMetadata {
	return m.metadata
}
