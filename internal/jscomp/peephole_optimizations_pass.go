// Port of com/google/javascript/jscomp/PeepholeOptimizationsPass.java @ closure-compiler v20261028

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
		} else if !parallel || !p.processChangedScopesParallel(root, changedScopeNodes) {
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

func (p *PeepholeOptimizationsPass) newWorkerPass(rc *Compiler) *PeepholeOptimizationsPass {
	workerPass := p.cloneForWorker(rc)
	workerPass.BeginTraversal()
	return workerPass
}

func finishPeepholeWorker(_ int, workerPass *PeepholeOptimizationsPass) {
	workerPass.pendingFeatures = 0
	workerPass.EndTraversal()
}

type peepholeScopeOutcome struct {
	events   []recordedCompilerEvent
	panicVal any
}

type peepholeWaveCallback struct {
	pass        *PeepholeOptimizationsPass
	rootScope   rhino.Node
	slotByScope map[rhino.Node]int
	outcomes    []peepholeScopeOutcome
	rc          *recordingCompiler
	childPanic  any
}

var (
	_ NodeTraversal_Callback                            = (*peepholeWaveCallback)(nil)
	_ nodeTraversal_onlyScriptAndFunctionShouldTraverse = (*peepholeWaveCallback)(nil)
	_ nodeTraversal_visitTokens                         = (*peepholeWaveCallback)(nil)
)

func (*peepholeWaveCallback) shouldTraverseOnlyScriptAndFunction() bool { return true }

func (w *peepholeWaveCallback) visitTokens() *ast.TokenSet { return w.pass.visitTokens() }

func (w *peepholeWaveCallback) Visit(t *NodeTraversal, n, parent rhino.Node) {
	if w.childPanic != nil {
		return
	}
	w.pass.Visit(t, n, parent)
}

func (w *peepholeWaveCallback) ShouldTraverse(t *NodeTraversal, n, _ rhino.Node) bool {
	if w.childPanic != nil {
		return false
	}
	if n == w.rootScope || !n.IsFunction() {
		return true
	}
	if slot, ok := w.slotByScope[n]; ok {
		out := &w.outcomes[slot]
		if len(out.events) > 0 {
			w.rc.parallelRecorder.events = append(w.rc.parallelRecorder.events, out.events...)
			out.events = nil
		}
		if out.panicVal != nil {
			w.childPanic = out.panicVal
			out.panicVal = nil
			t.noVisit = true
			t.noBranchShouldTraverse = false
			return false
		}
	}
	return false
}

type peepholeWorkerSlot struct {
	rc     recordingCompiler
	pass   *PeepholeOptimizationsPass
	waveCb peepholeWaveCallback
	ready  bool
}

func (s *peepholeWorkerSlot) prepare(p *PeepholeOptimizationsPass) {
	if !s.ready {
		initRecordingCompiler(&s.rc, p.compiler)
		s.pass = p.newWorkerPass(&s.rc)
		s.ready = true
		return
	}
	resetRecordingCompiler(&s.rc, p.compiler)
	s.pass.pendingFeatures = 0
}

func cleanupPeepholeWorkerPool(pool []peepholeWorkerSlot) {
	for i := range pool {
		if pool[i].ready {
			finishPeepholeWorker(i, pool[i].pass)
		}
	}
}

func runPeepholeWave(
	p *PeepholeOptimizationsPass,
	pool []peepholeWorkerSlot,
	itemCount int,
	runItem func(slot *peepholeWorkerSlot, idx int),
	afterItem func(slot *peepholeWorkerSlot, idx int, panicVal any),
) {
	if itemCount <= 0 {
		return
	}
	if itemCount == 1 || len(pool) <= 1 {
		slot := &pool[0]
		for idx := 0; idx < itemCount; idx++ {
			slot.prepare(p)
			var itemPanic any
			func() {
				defer func() {
					if r := recover(); r != nil {
						itemPanic = r
					}
				}()
				runItem(slot, idx)
			}()
			if itemPanic != nil {
				slot.pass.pendingFeatures = 0
			}
			afterItem(slot, idx, itemPanic)
		}
		return
	}

	numWorkers := min(len(pool), itemCount)
	workerPanic := spawnIndexedWorkers(itemCount, numWorkers, func(wIdx int, claim func() (int, bool)) {
		slot := &pool[wIdx]
		for {
			idx, ok := claim()
			if !ok {
				return
			}
			slot.prepare(p)
			var itemPanic any
			func() {
				defer func() {
					if r := recover(); r != nil {
						itemPanic = r
					}
				}()
				runItem(slot, idx)
			}()
			if itemPanic != nil {
				slot.pass.pendingFeatures = 0
			}
			afterItem(slot, idx, itemPanic)
		}
	})
	if workerPanic != nil {
		rethrowParallelPanic(p.compiler, workerPanic)
	}
}

func peepholeHasHugeScript(compiler *Compiler, scripts []rhino.Node) bool {
	for _, script := range scripts {
		if script.GetChildCount() >= 512 {
			return true
		}
		if compiler != nil {
			if inp := compiler.GetInput(script.GetInputId()); inp != nil {
				if sf := inp.GetSourceFile(); sf != nil {
					if sf.numBytes >= 256*1024 || len(sf.lineOffsets) >= 4096 {
						return true
					}
					if cp := sf.code.Load(); cp != nil && len(*cp) >= 256*1024 {
						return true
					}
				}
			}
		}
	}
	return false
}

func peepholeNodeAttachedToRoot(n, root rhino.Node) bool {
	if n.IsNil() || n.IsDeleted() {
		return false
	}
	for anc := n; !anc.IsNil(); anc = anc.GetParent() {
		if anc == root {
			return true
		}
	}
	return false
}

func (p *PeepholeOptimizationsPass) traverseFullScriptsParallel(root rhino.Node) bool {
	return traverseScriptsParallel(p.compiler, root, parallelScriptOptions[*PeepholeOptimizationsPass]{
		initWorker:   func(_ int, rc *Compiler) *PeepholeOptimizationsPass { return p.newWorkerPass(rc) },
		finishWorker: finishPeepholeWorker,
	})
}

func (p *PeepholeOptimizationsPass) processFullRootParallel(root rhino.Node) bool {
	scripts, ok := collectRootScripts(root)
	if !ok || len(scripts) == 0 {
		return false
	}
	opts := p.compiler.GetOptions()
	if opts == nil || opts.GetNumParallelThreads() <= 1 {
		return false
	}
	maxThreads := opts.GetNumParallelThreads()

	if len(scripts) >= maxThreads*2 && !peepholeHasHugeScript(p.compiler, scripts) {
		return p.traverseFullScriptsParallel(root)
	}

	var waves [][]rhino.Node
	waves = append(waves, scripts)
	totalFunctions := 0
	var collectFns func(n rhino.Node, depth int)
	collectFns = func(n rhino.Node, depth int) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if c.IsFunction() {
				totalFunctions++
				for len(waves) <= depth {
					waves = append(waves, nil)
				}
				waves[depth] = append(waves[depth], c)
				collectFns(c, depth+1)
				continue
			}
			if !c.GetFirstChild().IsNil() {
				collectFns(c, depth)
			}
		}
	}
	for _, script := range scripts {
		collectFns(script, 1)
	}
	if totalFunctions < 32 {
		return p.traverseFullScriptsParallel(root)
	}

	maxWaveSize := 0
	for _, w := range waves {
		if len(w) > maxWaveSize {
			maxWaveSize = len(w)
		}
	}
	if maxWaveSize < 2 {
		return p.traverseFullScriptsParallel(root)
	}

	totalScopes := len(scripts) + totalFunctions
	outcomes := make([]peepholeScopeOutcome, totalScopes)
	slotByScope := make(map[rhino.Node]int, totalScopes)
	nextSlot := 0
	for _, w := range waves {
		for _, sc := range w {
			slotByScope[sc] = nextSlot
			nextSlot++
		}
	}

	pool := make([]peepholeWorkerSlot, min(maxThreads, maxWaveSize))
	defer cleanupPeepholeWorkerPool(pool)

	jsErrorPinsDeferred.Add(1)
	defer jsErrorPinsDeferred.Add(-1)

	deepest := len(waves) - 1
	for d := deepest; d >= 0; d-- {
		wave := waves[d]
		isLeafWave := d == deepest
		runPeepholeWave(
			p,
			pool,
			len(wave),
			func(slot *peepholeWorkerSlot, idx int) {
				u := wave[idx]
				if isLeafWave {
					b := NodeTraversal_Builder{compiler: &slot.rc, callback: slot.pass}
					b.traverseScopedRoot(u, false)
					return
				}
				slot.waveCb = peepholeWaveCallback{
					pass:        slot.pass,
					rootScope:   u,
					slotByScope: slotByScope,
					outcomes:    outcomes,
					rc:          &slot.rc,
				}
				b := NodeTraversal_Builder{compiler: &slot.rc, callback: &slot.waveCb}
				b.traverseScopedRoot(u, true)
			},
			func(slot *peepholeWorkerSlot, idx int, panicVal any) {
				uSlot := slotByScope[wave[idx]]
				outcomes[uSlot].events = slot.rc.takeEvents()
				if !isLeafWave && slot.waveCb.childPanic != nil {
					outcomes[uSlot].panicVal = slot.waveCb.childPanic
					slot.waveCb.childPanic = nil
				} else {
					outcomes[uSlot].panicVal = panicVal
				}
			},
		)
	}

	for i := 0; i < len(scripts); i++ {
		if len(outcomes[i].events) != 0 {
			replayCompilerEvents(p.compiler, outcomes[i].events)
		}
		if outcomes[i].panicVal != nil {
			rethrowParallelPanic(p.compiler, outcomes[i].panicVal)
		}
	}
	return true
}

func (p *PeepholeOptimizationsPass) processChangedScopesParallel(root rhino.Node, changedScopeNodes []rhino.Node) bool {
	n := len(changedScopeNodes)
	numThreads := parallelThreadCount(p.compiler, n)
	if n < 2 || numThreads <= 1 {
		return false
	}

	inChanged := make(map[rhino.Node]struct{}, n)
	for _, sc := range changedScopeNodes {
		inChanged[sc] = struct{}{}
	}

	var waves [][]int
	for i, sc := range changedScopeNodes {
		if sc.IsNil() || sc.IsDeleted() {
			return false
		}
		ancCount := 0
		reachedRoot := false
		for anc := sc; !anc.IsNil(); anc = anc.GetParent() {
			if anc == root {
				reachedRoot = true
				break
			}
			if anc != sc && ChangeTracker_IsChangeScopeRoot(anc) {
				if _, ok := inChanged[anc]; ok {
					ancCount++
				}
			}
		}
		if !reachedRoot {
			return false
		}
		for len(waves) <= ancCount {
			waves = append(waves, nil)
		}
		waves[ancCount] = append(waves[ancCount], i)
	}

	maxWaveSize := 0
	for _, w := range waves {
		if len(w) > maxWaveSize {
			maxWaveSize = len(w)
		}
	}
	if maxWaveSize < 2 {
		return false
	}

	eventsByItem := make([][]recordedCompilerEvent, n)
	panicsByItem := make([]any, n)
	pool := make([]peepholeWorkerSlot, min(numThreads, maxWaveSize))
	defer cleanupPeepholeWorkerPool(pool)

	jsErrorPinsDeferred.Add(1)
	defer jsErrorPinsDeferred.Add(-1)

	var active []int
	for wIdx, w := range waves {
		if wIdx == 0 {
			active = w
		} else {
			active = active[:0]
			for _, origIdx := range w {
				if peepholeNodeAttachedToRoot(changedScopeNodes[origIdx], root) {
					active = append(active, origIdx)
				}
			}
		}
		if len(active) == 0 {
			continue
		}
		runPeepholeWave(
			p,
			pool,
			len(active),
			func(slot *peepholeWorkerSlot, idx int) {
				origIdx := active[idx]
				b := NodeTraversal_Builder{compiler: &slot.rc, callback: slot.pass}
				b.traverseScopedRoot(changedScopeNodes[origIdx], false)
			},
			func(slot *peepholeWorkerSlot, idx int, panicVal any) {
				origIdx := active[idx]
				eventsByItem[origIdx] = slot.rc.takeEvents()
				panicsByItem[origIdx] = panicVal
			},
		)
	}

	for i := 0; i < n; i++ {
		if len(eventsByItem[i]) != 0 {
			replayCompilerEvents(p.compiler, eventsByItem[i])
		}
		if panicsByItem[i] != nil {
			rethrowParallelPanic(p.compiler, panicsByItem[i])
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
		if p.compiler != nil && p.compiler.parallelRecorder != nil {
			p.compiler.parallelRecorder.recordScriptFeatures(p.compiler, script, p.pendingFeatures)
		} else {
			for w := p.pendingFeatures; w != 0; {
				b := bits.TrailingZeros64(w)
				w &^= uint64(1) << uint(b)
				NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature(b), p.compiler)
			}
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
