// Port of com/google/javascript/refactoring/JsFlumeCallback.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsFlumeCallback is a compiler node traversal callback that invokes matchers
// against every node and keeps track of any suggested fixes from the
// refactoring.
//
// Java: final class JsFlumeCallback implements NodeTraversal.Callback
type JsFlumeCallback struct {
	scanner            Scanner
	includeFilePattern *jutil.Pattern // may be nil
	matches            []*Match
	fixes              *jutil.ArrayList[*SuggestedFix]
}

// Java: JsFlumeCallback(Scanner, Pattern)
func NewJsFlumeCallback(scanner Scanner, includeFilePattern *jutil.Pattern) *JsFlumeCallback {
	return &JsFlumeCallback{
		scanner:            scanner,
		includeFilePattern: includeFilePattern,
		fixes:              jutil.NewArrayList[*SuggestedFix](),
	}
}

// Java: JsFlumeCallback.getMatches()
func (c *JsFlumeCallback) GetMatches() []*Match {
	return c.matches
}

// Java: JsFlumeCallback.getFixes()
func (c *JsFlumeCallback) GetFixes() *jutil.ArrayList[*SuggestedFix] {
	return c.fixes
}

// Java: JsFlumeCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *JsFlumeCallback) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsFromExterns() {
		return false
	}
	if n.IsScript() {
		filename := n.GetSourceFileName()
		if c.includeFilePattern != nil &&
			c.includeFilePattern.Pattern() != "" &&
			filename != nil && *filename != "" {
			return c.includeFilePattern.Matcher(*filename).Find()
		}
	}
	return true
}

// Java: JsFlumeCallback.visit(NodeTraversal, Node, Node)
func (c *JsFlumeCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	metadata := NodeMetadata_FromTraversal(t)
	if c.scanner.Matches(n, metadata) {
		match := NewMatch(n, metadata)
		for _, fix := range c.scanner.ProcessMatch(match) {
			c.fixes.Add(fix)
		}
		c.matches = append(c.matches, match)
	}
}
