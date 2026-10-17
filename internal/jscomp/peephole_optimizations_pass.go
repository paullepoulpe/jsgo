// Port of com/google/javascript/jscomp/PeepholeOptimizationsPass.java @ closure-compiler v20260415

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PeepholeOptimizationsPass is a compiler pass to run various peephole
// optimizations (e.g. constant folding, some useless code removal, some
// minimizations).
type PeepholeOptimizationsPass struct {
	compiler AbstractCompiler
	passName string
	// NOTE: Use a native array rather than a List to avoid creating iterators
	// for every node in the AST.
	peepholeOptimizations []AbstractPeepholeOptimization
	retraverseOnChange    bool
}

// NewPeepholeOptimizationsPass creates a peephole optimization pass that
// runs the given optimizations.
//
// Java: PeepholeOptimizationsPass(AbstractCompiler, String, AbstractPeepholeOptimization...)
func NewPeepholeOptimizationsPass(
	compiler AbstractCompiler, passName string, optimizations ...AbstractPeepholeOptimization) *PeepholeOptimizationsPass {
	return NewPeepholeOptimizationsPassOfList(compiler, passName, optimizations)
}

// Java: PeepholeOptimizationsPass(AbstractCompiler, String, List<AbstractPeepholeOptimization>)
func NewPeepholeOptimizationsPassOfList(
	compiler AbstractCompiler,
	passName string,
	optimizations []AbstractPeepholeOptimization) *PeepholeOptimizationsPass {
	p := &PeepholeOptimizationsPass{}
	p.compiler = compiler
	p.passName = passName
	p.peepholeOptimizations = append([]AbstractPeepholeOptimization(nil), optimizations...)
	p.retraverseOnChange = true
	return p
}

// Java: PeepholeOptimizationsPass.setRetraverseOnChange(boolean)
func (p *PeepholeOptimizationsPass) SetRetraverseOnChange(retraverse bool) {
	p.retraverseOnChange = retraverse
}

// Java: PeepholeOptimizationsPass.process(Node, Node)
func (p *PeepholeOptimizationsPass) Process(externs *rhino.Node, root *rhino.Node) {
	p.BeginTraversal()

	// Repeat to an internal fixed point.
	for changedScopeNodes := p.compiler.GetChangeTracker().GetChangedScopeNodesForPass(p.passName); changedScopeNodes == nil || len(changedScopeNodes) != 0; changedScopeNodes = p.compiler.GetChangeTracker().GetChangedScopeNodesForPass(p.passName) {

		if changedScopeNodes == nil {
			// changedScopeNodes is null if this is the first run of peepholeOptimizationsPass.
			NodeTraversal_Traverse(p.compiler, root, newPeepholeOptimizationsPass_PeepCallback(p))
		} else {
			NodeTraversal_TraverseScopeRoots(
				p.compiler, changedScopeNodes, newPeepholeOptimizationsPass_PeepCallback(p), false /* traverseNested */)
		}

		// Cancel the fixed point if requested.
		if !p.retraverseOnChange {
			break
		}
	}

	p.EndTraversal()
}

// peepholeOptimizationsPass_PeepCallback is the inner class PeepCallback.
type peepholeOptimizationsPass_PeepCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
	this *PeepholeOptimizationsPass
}

func newPeepholeOptimizationsPass_PeepCallback(outer *PeepholeOptimizationsPass) *peepholeOptimizationsPass_PeepCallback {
	return &peepholeOptimizationsPass_PeepCallback{this: outer}
}

// Java: PeepholeOptimizationsPass.PeepCallback.visit(NodeTraversal, Node, Node)
func (c *peepholeOptimizationsPass_PeepCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	currentNode := n
	for _, optim := range c.this.peepholeOptimizations {
		currentNode = optim.OptimizeSubtree(currentNode)
		if currentNode == nil {
			return
		}
	}
}

// Java: PeepholeOptimizationsPass.PeepCallback.enterScope(NodeTraversal)
func (c *peepholeOptimizationsPass_PeepCallback) EnterScope(t *NodeTraversal) {}

// Java: PeepholeOptimizationsPass.PeepCallback.exitScope(NodeTraversal)
func (c *peepholeOptimizationsPass_PeepCallback) ExitScope(t *NodeTraversal) {
	// Call updateFeatures() here, instead of immediately after
	// `optim.optimizeSubtree(currentNode)` in visit, to avoid repeating work for
	// every AST node x every peephole optimization. In practice, profiling shows a small
	// but measurable improvement for large projects. See comments on cl/856752797.
	c.UpdateFeatures(t)
}

// Java: PeepholeOptimizationsPass.PeepCallback.updateFeatures(NodeTraversal)
func (c *peepholeOptimizationsPass_PeepCallback) UpdateFeatures(t *NodeTraversal) {
	for _, optim := range c.this.peepholeOptimizations {
		newFeatures := optim.GetNewFeatures()
		if !newFeatures.IsEmpty() {
			// Go-only: language-adapter — FeatureSet.with(Set<Feature>) takes the Go EnumSet form.
			NodeUtil_AddFeaturesToScript(
				t.GetCurrentScript(),
				parser.FeatureSet_BARE_MINIMUM.WithEnumSet(jutil.NewEnumSet(newFeatures.ToSlice()...)),
				c.this.compiler)
		}
		optim.ClearNewFeatures()
	}
}

// BeginTraversal makes sure that all the optimizations have the current
// compiler so they can report errors.
//
// Java: PeepholeOptimizationsPass.beginTraversal()
func (p *PeepholeOptimizationsPass) BeginTraversal() {
	for _, optimization := range p.peepholeOptimizations {
		optimization.BeginTraversal(p.compiler)
	}
}

// EndTraversal ends the traversal.
//
// Java: PeepholeOptimizationsPass.endTraversal()
func (p *PeepholeOptimizationsPass) EndTraversal() {
	for _, optimization := range p.peepholeOptimizations {
		optimization.EndTraversal()
	}
}
