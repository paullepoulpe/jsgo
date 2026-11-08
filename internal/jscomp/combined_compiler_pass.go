// Port of com/google/javascript/jscomp/CombinedCompilerPass.java @ closure-compiler v20261104

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CombinedCompilerPass is a compiler pass combining multiple Callback and
// ScopedCallback objects. This pass can be used to separate logically
// different verifications without incurring any additional traversal and
// CFG generation costs.
//
// Due to this compiler pass' nature, none of the callbacks may mutate the
// parse tree.
//
// TODO(user): This combined pass is currently limited in the type of
// callbacks it can combine due to the difficulty of handling
// NodeTraversal's methods that initiate more recursion (e.g.,
// NodeTraversal#traverse(Node) and NodeTraversal#traverseInnerNode(Node,
// Node, Scope)). The NodeTraversal object passed to the individual
// callbacks should be instrumented to emulate the correct behavior. For
// instance, one could create a NodeTraversal whose
// NodeTraversal#traverseInnerNode(Node, Node, Scope) ties back into this
// compiler pass to give it context about what combined passes are doing.
type CombinedCompilerPass struct {
	// The callbacks that this pass combines.
	callbacks                   []combinedCompilerPass_CallbackWrapper
	shouldTraverseIndices       []int
	visitCallbacks              []NodeTraversal_Callback
	scopedIndices               []int
	visitTokenMask              [256]uint32
	waitingTokenMask            [256]uint32
	unionVisitTokens            ast.TokenSet
	activeVisitTokens           ast.TokenSet
	visitJSDocMask              uint32
	shouldTraverseMask          uint32
	onlyScriptAndFunctionSTMask uint32
	unreachableCodeMask         uint32
	scopedMask                  uint32
	inactiveMask                uint32
	useVisitTokenMask           bool
	hasVisitTokenFilter         bool
	compiler                    AbstractCompiler
	cachedErrorManager          ErrorManager
	sortingErrorManager         *SortingErrorManager
	continueAfterErrors         bool
	allNoVisit                  bool
	inactiveCount               int
}

var (
	_ nodeTraversal_visitTokens = (*CombinedCompilerPass)(nil)
	_ nodeTraversal_visitJSDoc  = (*CombinedCompilerPass)(nil)
)

func (p *CombinedCompilerPass) visitTokens() *ast.TokenSet {
	if p.hasVisitTokenFilter || (p.allNoVisit && p.useVisitTokenMask) {
		return &p.activeVisitTokens
	}
	return nil
}

func (p *CombinedCompilerPass) visitJSDoc() bool {
	return p.hasVisitTokenFilter && (p.visitJSDocMask&^p.inactiveMask) != 0
}

// NewCombinedCompilerPass creates a combined compiler pass.
//
// Java: CombinedCompilerPass(AbstractCompiler, NodeTraversal.Callback...)
func NewCombinedCompilerPass(compiler AbstractCompiler, callbacks ...NodeTraversal_Callback) *CombinedCompilerPass {
	// ImmutableList.copyOf(callbacks) rejects null elements.
	for i, c := range callbacks {
		if c == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + strconv.Itoa(i)))
		}
	}
	return NewCombinedCompilerPassOfList(compiler, slices.Clone(callbacks))
}

// Java: CombinedCompilerPass(AbstractCompiler, List<NodeTraversal.Callback>)
func NewCombinedCompilerPassOfList(compiler AbstractCompiler, callbacks []NodeTraversal_Callback) *CombinedCompilerPass {
	p := &CombinedCompilerPass{
		compiler:          compiler,
		allNoVisit:        true,
		useVisitTokenMask: len(callbacks) <= 32,
	}
	p.callbacks = make([]combinedCompilerPass_CallbackWrapper, len(callbacks))
	hasUnfilteredVisit := false
	for i, cb := range callbacks {
		w, hasShouldTraverse := initCombinedCompilerPass_CallbackWrapper(cb)
		bit := uint32(1) << i
		if hasShouldTraverse {
			p.shouldTraverseIndices = append(p.shouldTraverseIndices, i)
			if p.useVisitTokenMask {
				p.shouldTraverseMask |= bit
				if w.onlyScriptAndFunction {
					p.onlyScriptAndFunctionSTMask |= bit
				}
				if _, isUC := w.callback.(*CheckUnreachableCode); isUC {
					p.unreachableCodeMask |= bit
				}
			}
		}
		if !w.noVisit {
			p.allNoVisit = false
			p.visitCallbacks = append(p.visitCallbacks, w.callback)
			if p.useVisitTokenMask {
				if f, ok := w.callback.(nodeTraversal_visitTokens); ok && f.visitTokens() != nil {
					ts := f.visitTokens()
					for wordIdx := 0; wordIdx < 4; wordIdx++ {
						for word := ts[wordIdx]; word != 0; word &= word - 1 {
							tok := wordIdx*64 + bits.TrailingZeros64(word)
							p.visitTokenMask[tok] |= bit
						}
						p.unionVisitTokens[wordIdx] |= ts[wordIdx]
					}
					if vj, ok := w.callback.(nodeTraversal_visitJSDoc); ok && vj.visitJSDoc() {
						p.visitJSDocMask |= bit
					}
				} else {
					hasUnfilteredVisit = true
					for tok := range p.visitTokenMask {
						p.visitTokenMask[tok] |= bit
					}
				}
			}
		}
		if w.scopedCallback != nil {
			p.scopedIndices = append(p.scopedIndices, i)
			if p.useVisitTokenMask {
				p.scopedMask |= bit
			}
		}
		p.callbacks[i] = w
	}
	p.hasVisitTokenFilter = p.useVisitTokenMask && !hasUnfilteredVisit && !p.allNoVisit
	p.activeVisitTokens = p.unionVisitTokens
	p.refreshErrorState()
	return p
}

func (p *CombinedCompilerPass) refreshErrorState() {
	p.continueAfterErrors = false
	p.cachedErrorManager = nil
	p.sortingErrorManager = nil
	c := p.compiler
	if c == nil || c.isParallelWorker() {
		return
	}
	if opts := c.GetOptions(); opts != nil && opts.CanContinueAfterErrors() {
		p.continueAfterErrors = true
		return
	}
	em := c.errorManager
	p.cachedErrorManager = em
	if ts, ok := em.(*ThreadSafeDelegatingErrorManager); ok {
		em = ts.delegated
	}
	if sem, ok := em.(interface{ sortingErrorManager() *SortingErrorManager }); ok {
		p.sortingErrorManager = sem.sortingErrorManager()
	}
}

func (p *CombinedCompilerPass) hasHaltingErrors() bool {
	if c := p.compiler; c != nil && !c.isParallelWorker() {
		if p.continueAfterErrors {
			return false
		}
		if p.sortingErrorManager != nil && c.errorManager == p.cachedErrorManager {
			return p.sortingErrorManager.originalErrorCount != 0
		}
	}
	return p.compiler.HasHaltingErrors()
}

// Java: CombinedCompilerPass.traverse(AbstractCompiler, Node, List<NodeTraversal.Callback>)
func CombinedCompilerPass_Traverse(compiler AbstractCompiler, root rhino.Node, callbacks []NodeTraversal_Callback) {
	if len(callbacks) == 1 {
		NodeTraversal_Traverse(compiler, root, callbacks[0])
	} else {
		NewCombinedCompilerPassOfList(compiler, callbacks).Process(rhino.NilNode, root)
	}
}

// combinedCompilerPass_CallbackWrapper maintains information about a
// callback in order to simulate it being the exclusive client of the shared
// NodeTraversal. In particular, this class simulates abbreviating the
// traversal when the wrapped callback returns false for
// Callback#shouldTraverse(NodeTraversal, Node, Node). The callback becomes
// inactive (i.e., traversal messages are not sent to it) until the main
// traversal revisits the node during the post-order visit.
type combinedCompilerPass_CallbackWrapper struct {
	// The callback being wrapped. Never null.
	callback NodeTraversal_Callback
	// if (callback instanceof ScopedCallback), then scopedCallback points to
	// an instance of ScopedCallback, otherwise scopedCallback points to null
	scopedCallback NodeTraversal_ScopedCallback

	// The node that Callback#shouldTraverse(NodeTraversal, Node, Node)
	// returned false for. The wrapped callback doesn't receive messages until
	// after this node is revisited in the post-order traversal.
	waiting               rhino.Node
	noVisit               bool
	onlyScriptAndFunction bool
}

// Java: CombinedCompilerPass.CallbackWrapper(NodeTraversal.Callback)
func initCombinedCompilerPass_CallbackWrapper(callback NodeTraversal_Callback) (w combinedCompilerPass_CallbackWrapper, hasShouldTraverse bool) {
	w.callback = callback
	if scopedCallback, ok := callback.(NodeTraversal_ScopedCallback); ok {
		w.scopedCallback = scopedCallback
	}
	_, noShouldTraverse := callback.(nodeTraversal_noShouldTraverse)
	if !noShouldTraverse {
		_, noShouldTraverse = callback.(NodeTraversal_AbstractPostOrderCallbackInterface)
	}
	if !noShouldTraverse {
		if sf, ok := callback.(nodeTraversal_onlyScriptAndFunctionShouldTraverse); ok {
			w.onlyScriptAndFunction = sf.shouldTraverseOnlyScriptAndFunction()
		}
	}
	_, w.noVisit = callback.(nodeTraversal_noVisit)
	return w, !noShouldTraverse
}

func (p *CombinedCompilerPass) resetActiveState() {
	if p.inactiveCount != 0 {
		for i := range p.callbacks {
			if wait := p.callbacks[i].waiting; !wait.IsNil() {
				if p.useVisitTokenMask {
					p.waitingTokenMask[byte(wait.GetToken())] = 0
				}
				p.callbacks[i].waiting = rhino.NilNode
			}
		}
		p.inactiveCount = 0
		p.inactiveMask = 0
	}
	p.activeVisitTokens = p.unionVisitTokens
}

func (p *CombinedCompilerPass) deactivate(idx int, n rhino.Node) {
	p.callbacks[idx].waiting = n
	p.inactiveCount++
	if p.useVisitTokenMask {
		bit := uint32(1) << idx
		p.inactiveMask |= bit
		tok := n.GetToken()
		p.waitingTokenMask[byte(tok)] |= bit
		p.activeVisitTokens.Add(tok)
	}
}

func (p *CombinedCompilerPass) activeScoped(idx int) NodeTraversal_ScopedCallback {
	if w := &p.callbacks[idx]; w.waiting.IsNil() {
		return w.scopedCallback
	}
	return nil
}

func (p *CombinedCompilerPass) syncTraversalVisitFilter(t *NodeTraversal) {
	if t == nil {
		return
	}
	if p.allNoVisit {
		t.noVisit = p.inactiveCount == 0
		if p.useVisitTokenMask {
			t.visitTokens = &p.activeVisitTokens
		}
	} else if p.hasVisitTokenFilter {
		t.visitTokens = &p.activeVisitTokens
		t.visitJSDoc = (p.visitJSDocMask &^ p.inactiveMask) != 0
	}
}

// Java: CombinedCompilerPass.process(Node, Node)
func (p *CombinedCompilerPass) Process(externs rhino.Node, root rhino.Node) {
	p.resetActiveState()
	p.refreshErrorState()
	if p.processParallel(root) {
		return
	}
	NodeTraversal_Traverse(p.compiler, root, p)
}

type parallelCloneableCallback interface {
	cloneForParallel(c *Compiler) NodeTraversal_Callback
}

func (cb *CheckJSDoc) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	out := NewCheckJSDoc(c)
	out.inExterns = cb.inExterns
	return out
}

func (*CheckSuspiciousCode) cloneForParallel(*Compiler) NodeTraversal_Callback {
	return NewCheckSuspiciousCode()
}

func (*CheckGlobalThis) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckGlobalThis(c)
}

func (p *CombinedCompilerPass) canCloneForParallel() bool {
	for i := range p.callbacks {
		if _, cloneable := p.callbacks[i].callback.(parallelCloneableCallback); !cloneable {
			return false
		}
	}
	return true
}

func (p *CombinedCompilerPass) hasUnreachableCodeCallback() bool {
	if p.useVisitTokenMask {
		return p.unreachableCodeMask != 0
	}
	for i := range p.callbacks {
		if _, isUC := p.callbacks[i].callback.(*CheckUnreachableCode); isUC {
			return true
		}
	}
	return false
}

func (p *CombinedCompilerPass) cloneForWorker(rc *Compiler) *CombinedCompilerPass {
	cloned := make([]NodeTraversal_Callback, len(p.callbacks))
	var origUC *CheckUnreachableCode
	var clonedUC *CheckUnreachableCode

	for i := range p.callbacks {
		cb := p.callbacks[i].callback
		c := cb.(parallelCloneableCallback).cloneForParallel(rc)
		if uc, ok := c.(*CheckUnreachableCode); ok {
			origUC = cb.(*CheckUnreachableCode)
			clonedUC = uc
		}
		cloned[i] = c
	}
	if origUC != nil {
		for i := range p.callbacks {
			if mr, ok := p.callbacks[i].callback.(*CheckMissingReturn); ok && mr.sharedCfgSource == origUC {
				cloned[i].(*CheckMissingReturn).sharedCfgSource = clonedUC
			}
		}
	}
	return NewCombinedCompilerPassOfList(rc, cloned)
}

func (p *CombinedCompilerPass) prepareParallelScript(script rhino.Node) {
	p.resetActiveState()
	for i := range p.callbacks {
		switch cb := p.callbacks[i].callback.(type) {
		case *CheckUnreachableCode:
			cb.resetForParallelScript(script)
		case *CheckMissingReturn:
			cb.resetCfgStack()
		}
	}
}

func (p *CombinedCompilerPass) parallelScriptFellThrough() bool {
	for i := range p.callbacks {
		if uc, ok := p.callbacks[i].callback.(*CheckUnreachableCode); ok && !uc.scriptFallsThrough {
			return false
		}
	}
	return true
}

func (p *CombinedCompilerPass) processParallel(root rhino.Node) bool {
	if p.hasUnreachableCodeCallback() {
		return p.processParallelWithUnreachableCode(root)
	}
	halted := false
	ok := traverseScriptsParallel(p.compiler, root, parallelScriptOptions[*CombinedCompilerPass]{
		useSharedGlobalScope:    true,
		enableTypeRegistryReads: true,
		stopOnHaltingError:      !p.continueAfterErrors,
		prepare: func(_, _ int) bool {
			if p.hasHaltingErrors() {
				halted = true
				return false
			}
			return p.canCloneForParallel()
		},
		initWorker: func(_ int, rc *Compiler) *CombinedCompilerPass {
			return p.cloneForWorker(rc)
		},
		callback: func(wp *CombinedCompilerPass, _ int, _, _ rhino.Node, _ *Compiler) NodeTraversal_Callback {
			wp.resetActiveState()
			return wp
		},
	})
	return ok || halted
}

func (p *CombinedCompilerPass) processParallelWithUnreachableCode(root rhino.Node) bool {
	if p.hasHaltingErrors() {
		return true
	}
	if !p.canCloneForParallel() {
		return false
	}
	scripts, ok := collectRootScripts(root)
	if !ok {
		return false
	}
	itemCount := len(scripts)
	numWorkers := parallelThreadCount(p.compiler, itemCount)
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(p.compiler)
	restoreReg := enableParallelTypeRegistryReads(p.compiler)
	defer restoreReg()
	lazyGlobal := newLazyGlobalScope(p.compiler, root)

	eventsByItem := make([][]recordedCompilerEvent, itemCount)
	fallsThrough := make([]bool, itemCount)

	type ucWorkerState struct {
		rc recordingCompiler
		wp *CombinedCompilerPass
		sc *SyntacticScopeCreator
	}

	panicsByItem, workerPanic := collectParallelWorkersDeferredPins(
		itemCount,
		numWorkers,
		func(_ int, w *ucWorkerState) {
			initRecordingCompiler(&w.rc, p.compiler)
			w.wp = p.cloneForWorker(&w.rc)
			w.sc = NewSyntacticScopeCreator(&w.rc)
		},
		func(w *ucWorkerState, idx int) {
			resetRecordingCompiler(&w.rc, p.compiler)
			script := scripts[idx]
			w.wp.prepareParallelScript(script)
			t := acquireNodeTraversalDirect(&w.rc, w.wp, w.sc, false)
			traverseScriptUnderRoot(t, root, root, script, lazyGlobal)
			releaseNodeTraversal(t)
			fallsThrough[idx] = w.wp.parallelScriptFellThrough()
		},
		func(w *ucWorkerState, idx int) {
			eventsByItem[idx] = w.rc.takeEvents()
		},
		nil,
		false,
	)

	if workerPanic == nil {
		for i := 0; i < itemCount-1; i++ {
			if panicsByItem[i] == nil && !fallsThrough[i] {
				// A non-final script terminated unconditionally at top level without
				// trailing statements, making the following SCRIPT node unreachable in
				// the whole-program CFG. Fall back before replaying any events.
				return false
			}
		}
	}

	stopOnHaltingError := !p.continueAfterErrors
	for i := 0; i < itemCount; i++ {
		if stopOnHaltingError && p.compiler.HasHaltingErrors() {
			break
		}
		if len(eventsByItem[i]) != 0 {
			replayCompilerEvents(p.compiler, eventsByItem[i])
		}
		if ip := panicsByItem[i]; ip != nil {
			rethrowParallelPanic(p.compiler, ip)
		}
	}
	if workerPanic != nil {
		rethrowParallelPanic(p.compiler, workerPanic)
	}
	return true
}

// Java: CombinedCompilerPass.shouldTraverse(NodeTraversal, Node, Node)
func (p *CombinedCompilerPass) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if p.hasHaltingErrors() {
		return false
	}

	deactivated := false
	if p.useVisitTokenMask {
		stMask := p.shouldTraverseMask &^ p.inactiveMask
		if p.onlyScriptAndFunctionSTMask != 0 {
			if tok := n.GetToken(); tok != rhino.Token_SCRIPT && tok != rhino.Token_FUNCTION {
				stMask &^= p.onlyScriptAndFunctionSTMask
			}
		}
		for m := stMask; m != 0; m &= m - 1 {
			idx := bits.TrailingZeros32(m)
			if !p.callbacks[idx].callback.ShouldTraverse(t, n, parent) {
				p.deactivate(idx, n)
				deactivated = true
			}
		}
	} else {
		isSF := false
		if tok := n.GetToken(); tok == rhino.Token_SCRIPT || tok == rhino.Token_FUNCTION {
			isSF = true
		}
		for _, idx := range p.shouldTraverseIndices {
			w := &p.callbacks[idx]
			if w.waiting.IsNil() && (!w.onlyScriptAndFunction || isSF) && !w.callback.ShouldTraverse(t, n, parent) {
				p.deactivate(idx, n)
				deactivated = true
			}
		}
	}
	if deactivated || p.allNoVisit || len(p.shouldTraverseIndices) == 0 {
		p.syncTraversalVisitFilter(t)
		p.updateNoShouldTraverse(t)
	} else if t != nil && (t.noShouldTraverse || t.noBranchShouldTraverse) && p.hasHaltingErrors() {
		t.noShouldTraverse = false
		t.noBranchShouldTraverse = false
	}
	// Note that this method could return false if all callbacks are
	// inactive. This apparent optimization would make this method more
	// expensive in the typical case where not all nodes are inactive. It is
	// very unlikely that many all callbacks would be inactive at the same
	// time (indeed, there are several checking passes that never return
	// false).
	return true
}

// Java: CombinedCompilerPass.visit(NodeTraversal, Node, Node)
func (p *CombinedCompilerPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if p.inactiveCount == 0 {
		if p.allNoVisit {
			return
		}
		if p.useVisitTokenMask {
			mask := p.visitTokenMask[byte(n.GetToken())]
			if p.visitJSDocMask != 0 && n.HasJSDocInfo() {
				mask |= p.visitJSDocMask
			}
			if mask == 0 || p.hasHaltingErrors() {
				return
			}
			for ; mask != 0; mask &= mask - 1 {
				p.callbacks[bits.TrailingZeros32(mask)].callback.Visit(t, n, parent)
			}
		} else {
			if p.hasHaltingErrors() {
				return
			}
			for _, cb := range p.visitCallbacks {
				cb.Visit(t, n, parent)
			}
		}
		if t != nil && (t.noShouldTraverse || t.noBranchShouldTraverse) && p.hasHaltingErrors() {
			t.noShouldTraverse = false
			t.noBranchShouldTraverse = false
		}
		return
	}

	if p.useVisitTokenMask {
		tok := n.GetToken()
		btok := byte(tok)
		activeMask := p.visitTokenMask[btok]
		if p.visitJSDocMask != 0 && n.HasJSDocInfo() {
			activeMask |= p.visitJSDocMask
		}
		activeMask &^= p.inactiveMask
		waitMask := p.waitingTokenMask[btok]
		combined := activeMask | waitMask
		if combined == 0 || p.hasHaltingErrors() {
			return
		}
		reactivated := false
		for m := combined; m != 0; m &= m - 1 {
			idx := bits.TrailingZeros32(m)
			bit := uint32(1) << idx
			if (activeMask & bit) != 0 {
				p.callbacks[idx].callback.Visit(t, n, parent)
			} else if p.callbacks[idx].waiting == n {
				p.callbacks[idx].waiting = rhino.NilNode
				p.inactiveMask &^= bit
				p.inactiveCount--
				p.waitingTokenMask[btok] &^= bit
				if p.waitingTokenMask[btok] == 0 && !p.unionVisitTokens.Has(tok) {
					p.activeVisitTokens[btok>>6] &^= uint64(1) << (btok & 63)
				}
				reactivated = true
			}
		}
		if reactivated {
			p.syncTraversalVisitFilter(t)
			p.updateNoShouldTraverse(t)
		} else if t != nil && (t.noShouldTraverse || t.noBranchShouldTraverse) && p.hasHaltingErrors() {
			t.noShouldTraverse = false
			t.noBranchShouldTraverse = false
		}
		return
	}

	if p.hasHaltingErrors() {
		return
	}
	reactivated := false
	for i := range p.callbacks {
		w := &p.callbacks[i]
		if w.waiting.IsNil() {
			if !w.noVisit {
				w.callback.Visit(t, n, parent)
			}
		} else if w.waiting == n {
			w.waiting = rhino.NilNode
			p.inactiveCount--
			reactivated = true
		}
	}
	if reactivated {
		p.syncTraversalVisitFilter(t)
		p.updateNoShouldTraverse(t)
	} else if t != nil && (t.noShouldTraverse || t.noBranchShouldTraverse) && p.hasHaltingErrors() {
		t.noShouldTraverse = false
		t.noBranchShouldTraverse = false
	}
}

// Java: CombinedCompilerPass.enterScope(NodeTraversal)
func (p *CombinedCompilerPass) EnterScope(t *NodeTraversal) {
	if p.useVisitTokenMask {
		for m := p.scopedMask &^ p.inactiveMask; m != 0; m &= m - 1 {
			p.callbacks[bits.TrailingZeros32(m)].scopedCallback.EnterScope(t)
		}
	} else {
		for _, idx := range p.scopedIndices {
			if cb := p.activeScoped(idx); cb != nil {
				cb.EnterScope(t)
			}
		}
	}
	if p.allNoVisit {
		p.syncTraversalVisitFilter(t)
	}
	p.updateNoShouldTraverse(t)
}

// Java: CombinedCompilerPass.exitScope(NodeTraversal)
func (p *CombinedCompilerPass) ExitScope(t *NodeTraversal) {
	if p.useVisitTokenMask {
		for m := p.scopedMask &^ p.inactiveMask; m != 0; m &= m - 1 {
			p.callbacks[bits.TrailingZeros32(m)].scopedCallback.ExitScope(t)
		}
	} else {
		for _, idx := range p.scopedIndices {
			if cb := p.activeScoped(idx); cb != nil {
				cb.ExitScope(t)
			}
		}
	}
	p.updateNoShouldTraverse(t)
}

func (p *CombinedCompilerPass) updateNoShouldTraverse(t *NodeTraversal) {
	if t == nil {
		return
	}
	if p.hasHaltingErrors() {
		t.noShouldTraverse = false
		t.noBranchShouldTraverse = false
		return
	}
	if p.useVisitTokenMask {
		activeST := p.shouldTraverseMask &^ p.inactiveMask
		if ucMask := activeST & p.unreachableCodeMask; ucMask != 0 {
			for m := ucMask; m != 0; m &= m - 1 {
				idx := bits.TrailingZeros32(m)
				if p.callbacks[idx].callback.(*CheckUnreachableCode).CurrentScopeHasNoUnreachable() {
					activeST &^= uint32(1) << idx
				}
			}
		}
		t.noShouldTraverse = activeST == 0
		t.noBranchShouldTraverse = (activeST &^ p.onlyScriptAndFunctionSTMask) == 0
		return
	}
	noST := true
	noBranchST := true
	for _, idx := range p.shouldTraverseIndices {
		w := &p.callbacks[idx]
		if !w.waiting.IsNil() {
			continue
		}
		if uc, ok := w.callback.(*CheckUnreachableCode); ok && uc.CurrentScopeHasNoUnreachable() {
			continue
		}
		noST = false
		if !w.onlyScriptAndFunction {
			noBranchST = false
			break
		}
	}
	t.noShouldTraverse = noST
	t.noBranchShouldTraverse = noBranchST
}
