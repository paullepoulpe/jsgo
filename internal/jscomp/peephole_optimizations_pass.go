// Port of com/google/javascript/jscomp/PeepholeOptimizationsPass.java @ closure-compiler v20261021

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
	"math/bits"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	pendingFeatures       uint64
	visitTokenSet         ast.TokenSet
	activeMaskByToken     [256]uint16
	useActiveMask         bool
	retraverseOnChange    bool
}

var (
	_ NodeTraversal_Callback         = (*PeepholeOptimizationsPass)(nil)
	_ nodeTraversal_noShouldTraverse = (*PeepholeOptimizationsPass)(nil)
	_ nodeTraversal_visitTokens      = (*PeepholeOptimizationsPass)(nil)
)

func (*PeepholeOptimizationsPass) nodeTraversalNoShouldTraverse() {}

func (p *PeepholeOptimizationsPass) visitTokens() *ast.TokenSet {
	if p.useActiveMask {
		return &p.visitTokenSet
	}
	return nil
}

// Java: PeepholeOptimizationsPass.PeepCallback.shouldTraverse(NodeTraversal, Node, Node)
func (p *PeepholeOptimizationsPass) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
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
	p.initDispatchTable()
	return p
}

func (p *PeepholeOptimizationsPass) initDispatchTable() {
	if len(p.peepholeOptimizations) > 16 {
		p.useActiveMask = false
		return
	}
	p.useActiveMask = true
	clear(p.activeMaskByToken[:])

	for idx, optim := range p.peepholeOptimizations {
		bit := uint16(1) << idx
		cfg := optim.peepholeConfig()
		if !cfg.isStandard {
			for tok := range p.activeMaskByToken {
				p.activeMaskByToken[tok] |= bit
			}
			continue
		}
		for w, word := range cfg.tokens {
			for word != 0 {
				b := bits.TrailingZeros64(word)
				word &^= uint64(1) << b
				p.activeMaskByToken[(w<<6)|b] |= bit
			}
		}
	}

	p.visitTokenSet = ast.TokenSet{}
	p.visitTokenSet.Add(rhino.Token_SCRIPT)
	p.visitTokenSet.Add(rhino.Token_FUNCTION)
	for tok, mask := range p.activeMaskByToken {
		if mask != 0 {
			p.visitTokenSet.Add(rhino.Token(tok))
		}
	}
}

// Java: PeepholeOptimizationsPass.setRetraverseOnChange(boolean)
func (p *PeepholeOptimizationsPass) SetRetraverseOnChange(retraverse bool) {
	p.retraverseOnChange = retraverse
}

// Java: PeepholeOptimizationsPass.process(Node, Node)
func (p *PeepholeOptimizationsPass) Process(externs rhino.Node, root rhino.Node) {
	parallel := p.canCloneForParallel() && parallelThreadCount(p.compiler, 2) > 1
	if parallel {
		prewarmCompilerForParallelReads(p.compiler)
		restoreReg := enableParallelTypeRegistryReads(p.compiler)
		defer restoreReg()
	}

	p.BeginTraversal()

	// Repeat to an internal fixed point.
	for changedScopeNodes := p.compiler.GetChangeTracker().GetChangedScopeNodesForPass(p.passName); changedScopeNodes == nil || len(changedScopeNodes) != 0; changedScopeNodes = p.compiler.GetChangeTracker().GetChangedScopeNodesForPass(p.passName) {
		if changedScopeNodes == nil {
			// changedScopeNodes is null if this is the first run of peepholeOptimizationsPass.
			if !parallel || !p.processFullRootParallel(root) {
				NodeTraversal_Traverse(p.compiler, root, p)
			}
		} else if !parallel || !p.processChangedScopesParallel(changedScopeNodes) {
			NodeTraversal_TraverseScopeRoots(
				p.compiler, changedScopeNodes, p, false /* traverseNested */)
		}

		// Cancel the fixed point if requested.
		if !p.retraverseOnChange {
			break
		}
	}

	p.EndTraversal()
}

func (p *PeepholeOptimizationsPass) canCloneForParallel() bool {
	for _, optim := range p.peepholeOptimizations {
		if !optim.peepholeConfig().isStandard {
			return false
		}
	}
	return true
}

func (p *PeepholeOptimizationsPass) cloneForWorker(rc *recordingCompiler) *PeepholeOptimizationsPass {
	clonedOpts := make([]AbstractPeepholeOptimization, len(p.peepholeOptimizations))
	for i, optim := range p.peepholeOptimizations {
		clonedOpts[i] = optim.cloneForWorker()
	}
	cp := NewPeepholeOptimizationsPassOfList(rc, p.passName, clonedOpts)
	cp.retraverseOnChange = p.retraverseOnChange
	return cp
}

func enclosingScriptForPeephole(n rhino.Node) rhino.Node {
	for cur := n; !cur.IsNil(); cur = cur.GetParent() {
		if cur.IsScript() {
			return cur
		}
	}
	return rhino.NilNode
}

type peepholeScriptBucket struct {
	origIndices []int
}

func (p *PeepholeOptimizationsPass) processFullRootParallel(root rhino.Node) bool {
	scripts, ok := collectRootScripts(root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(p.compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}
	runParallelWork(p.compiler, numWorkers, len(scripts), func(_ int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var workerPass *PeepholeOptimizationsPass

		cleanup := func() {
			if workerPass != nil {
				workerPass.EndTraversal()
			}
		}
		runItem := func(idx int, rc *recordingCompiler) {
			if workerPass == nil || workerRC != rc {
				workerRC = rc
				workerPass = p.cloneForWorker(rc)
				workerPass.BeginTraversal()
			}
			b := NodeTraversal_Builder{
				compiler: rc,
				callback: workerPass,
			}
			b.traverseScriptUnderRoot(root, root, scripts[idx], nil)
		}
		return runItem, cleanup
	}, false)
	return true
}

func (p *PeepholeOptimizationsPass) processChangedScopesParallel(changedScopeNodes []rhino.Node) bool {
	if len(changedScopeNodes) < 2 {
		return false
	}
	bucketByScript := make(map[rhino.Node]int, len(changedScopeNodes))
	var buckets []peepholeScriptBucket
	for origIdx, scopeNode := range changedScopeNodes {
		script := enclosingScriptForPeephole(scopeNode)
		if script.IsNil() {
			return false
		}
		bIdx, exists := bucketByScript[script]
		if !exists {
			bIdx = len(buckets)
			bucketByScript[script] = bIdx
			buckets = append(buckets, peepholeScriptBucket{origIndices: []int{origIdx}})
		} else {
			buckets[bIdx].origIndices = append(buckets[bIdx].origIndices, origIdx)
		}
	}
	numWorkers := parallelThreadCount(p.compiler, len(buckets))
	if numWorkers <= 1 {
		return false
	}

	eventsByOrigIdx := make([][]recordedCompilerEvent, len(changedScopeNodes))
	panicsByOrigIdx := make([]any, len(changedScopeNodes))

	runParallelWork(p.compiler, numWorkers, len(buckets), func(_ int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var workerPass *PeepholeOptimizationsPass

		cleanup := func() {
			if workerPass != nil {
				workerPass.EndTraversal()
			}
		}
		runItem := func(bIdx int, rc *recordingCompiler) {
			if workerPass == nil || workerRC != rc {
				workerRC = rc
				workerPass = p.cloneForWorker(rc)
				workerPass.BeginTraversal()
			}
			for _, origIdx := range buckets[bIdx].origIndices {
				resetRecordingCompiler(rc, p.compiler)
				panicked := runParallelItemIndexed(rc, origIdx, panicsByOrigIdx, func(rc *recordingCompiler, itemIdx int) {
					b := NodeTraversal_Builder{
						compiler: rc,
						callback: workerPass,
					}
					b.traverseScopedRoot(changedScopeNodes[itemIdx], false)
				})
				eventsByOrigIdx[origIdx] = rc.takeEvents()
				if panicked {
					workerPass.pendingFeatures = 0
					return
				}
			}
		}
		return runItem, cleanup
	}, false)

	for i := range changedScopeNodes {
		if len(eventsByOrigIdx[i]) != 0 {
			replayCompilerEvents(p.compiler, eventsByOrigIdx[i])
		}
		if r := panicsByOrigIdx[i]; r != nil {
			if ie, ok := r.(parallelInternalError); ok {
				p.compiler.ThrowInternalError(ie.msg, ie.cause)
			}
			panic(r)
		}
	}
	return true
}

// Java: PeepholeOptimizationsPass.PeepCallback.visit(NodeTraversal, Node, Node)
func (p *PeepholeOptimizationsPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	p.visitOptimize(n)
	if p.pendingFeatures != 0 && (n.IsScript() || n == t.scopeTraversalRoot) {
		script := t.GetCurrentScript()
		jutil.CheckStateMsg1(script.IsScript(), "%s", script)
		for w := p.pendingFeatures; w != 0; {
			b := bits.TrailingZeros64(w)
			w &^= uint64(1) << uint(b)
			NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature(b), p.compiler)
		}
		p.pendingFeatures = 0
	}
}

func (p *PeepholeOptimizationsPass) visitOptimize(n rhino.Node) {
	if !p.useActiveMask {
		currentNode := n
		for _, optim := range p.peepholeOptimizations {
			currentNode = optim.OptimizeSubtree(currentNode)
			if currentNode.IsNil() {
				return
			}
		}
		return
	}

	currentNode := n
	tok := currentNode.GetToken()
	mask := p.activeMaskByToken[tok]
	for mask != 0 {
		idx := bits.TrailingZeros16(mask)
		bit := uint16(1) << idx
		mask &^= bit
		currentNode = p.peepholeOptimizations[idx].OptimizeSubtree(currentNode)
		if currentNode.IsNil() {
			return
		}
		if newTok := currentNode.GetToken(); newTok != tok {
			tok = newTok
			higherMask := ^((bit << 1) - 1)
			mask = p.activeMaskByToken[tok] & higherMask
		}
	}
}

// BeginTraversal makes sure that all the optimizations have the current
// compiler so they can report errors.
//
// Java: PeepholeOptimizationsPass.beginTraversal()
func (p *PeepholeOptimizationsPass) BeginTraversal() {
	p.pendingFeatures = 0
	for _, optimization := range p.peepholeOptimizations {
		optimization.bindPendingFeatures(&p.pendingFeatures)
		optimization.BeginTraversal(p.compiler)
	}
	p.initDispatchTable()
}

// EndTraversal ends the traversal.
//
// Java: PeepholeOptimizationsPass.endTraversal()
func (p *PeepholeOptimizationsPass) EndTraversal() {
	if p.pendingFeatures != 0 {
		jutil.CheckStateMsg1(
			false,
			"Expected getNewFeatures() to be empty in endTraversal() but found %s",
			peepholeFeaturesToSet(p.pendingFeatures))
	}
	for _, optimization := range p.peepholeOptimizations {
		optimization.bindPendingFeatures(nil)
		optimization.EndTraversal()
	}
}
