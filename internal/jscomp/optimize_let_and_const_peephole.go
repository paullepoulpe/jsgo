// Port of com/google/javascript/jscomp/OptimizeLetAndConstPeephole.java @ closure-compiler v20261015

/*
 * Copyright 2026 The Closure Compiler Authors.
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

// OptimizeLetAndConstPeephole is a late-stage compiler pass that lowers
// `const`/`let` to `var` where safe, and `const` to `let`, improving keyword
// homogeneity for better compression.
type OptimizeLetAndConstPeephole struct {
	AbstractPeepholeOptimizationBase
	assumeOutputIsWrapped bool
}

// Java: OptimizeLetAndConstPeephole(boolean)
func NewOptimizeLetAndConstPeephole(assumeOutputIsWrapped bool) *OptimizeLetAndConstPeephole {
	p := &OptimizeLetAndConstPeephole{}
	p.InitAbstractPeepholeOptimizationBase(p) // the implicit super()
	p.assumeOutputIsWrapped = assumeOutputIsWrapped
	return p
}

// Java: OptimizeLetAndConstPeephole.optimizeSubtree(Node)
func (p *OptimizeLetAndConstPeephole) OptimizeSubtree(subtree rhino.Node) rhino.Node {
	switch subtree.GetToken() {
	case rhino.Token_LET:
		return p.ProcessLet(subtree)
	case rhino.Token_CONST:
		return p.ProcessConst(subtree)
	default:
		return subtree
	}
}

// Java: OptimizeLetAndConstPeephole.processLet(Node)
func (p *OptimizeLetAndConstPeephole) ProcessLet(n rhino.Node) rhino.Node {
	if !p.assumeOutputIsWrapped && n.GetParent().IsScript() {
		return n
	}

	if p.IsDirectlyInHoistScope(n) {
		n.SetToken(rhino.Token_VAR)
		p.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java: OptimizeLetAndConstPeephole.processConst(Node)
func (p *OptimizeLetAndConstPeephole) ProcessConst(n rhino.Node) rhino.Node {
	if !p.assumeOutputIsWrapped && n.GetParent().IsScript() {
		return n
	}

	if p.IsDirectlyInHoistScope(n) {
		n.SetToken(rhino.Token_VAR)
	} else {
		n.SetToken(rhino.Token_LET)
		p.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_LET_DECLARATIONS)
	}
	p.ReportChangeToEnclosingScope(n)
	return n
}

// Java: OptimizeLetAndConstPeephole.isDirectlyInHoistScope(Node)
func (p *OptimizeLetAndConstPeephole) IsDirectlyInHoistScope(n rhino.Node) bool {
	parent := n.GetParent()
	for !parent.IsNil() && parent.IsLabel() {
		parent = parent.GetParent()
	}
	if parent.IsNil() {
		return false
	}

	if parent.IsScript() || parent.IsModuleBody() {
		return true
	}

	if parent.IsBlock() {
		grandparent := parent.GetParent()
		if !grandparent.IsNil() && (grandparent.IsFunction() || grandparent.IsClassMembers()) {
			return true
		}
	}

	return false
}
