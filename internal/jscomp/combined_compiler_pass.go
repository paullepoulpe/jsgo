// Port of com/google/javascript/jscomp/CombinedCompilerPass.java @ closure-compiler v20261015

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
	"slices"
	"strconv"

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
	callbacks             []combinedCompilerPass_CallbackWrapper
	shouldTraverseIndices []int
	visitCallbacks        []NodeTraversal_Callback
	scopedIndices         []int
	compiler              AbstractCompiler
	concreteCompiler      *Compiler
	cachedErrorManager    ErrorManager
	haltingErrorCount     *int
	continueAfterErrors   bool
	allNoVisit            bool
	inactiveCount         int
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
	p := &CombinedCompilerPass{compiler: compiler, allNoVisit: true}
	p.callbacks = make([]combinedCompilerPass_CallbackWrapper, len(callbacks))
	for i := 0; i < len(callbacks); i++ {
		w := initCombinedCompilerPass_CallbackWrapper(callbacks[i])
		if !w.noShouldTraverse {
			p.shouldTraverseIndices = append(p.shouldTraverseIndices, i)
		}
		if !w.noVisit {
			p.allNoVisit = false
			p.visitCallbacks = append(p.visitCallbacks, w.callback)
		}
		if w.scopedCallback != nil {
			p.scopedIndices = append(p.scopedIndices, i)
		}
		p.callbacks[i] = w
	}
	p.refreshErrorState()
	return p
}

func (p *CombinedCompilerPass) refreshErrorState() {
	p.continueAfterErrors = false
	p.cachedErrorManager = nil
	p.haltingErrorCount = nil
	c, ok := p.compiler.(*Compiler)
	if !ok || c == nil {
		return
	}
	p.concreteCompiler = c
	if opts := c.GetOptions(); opts != nil && opts.CanContinueAfterErrors() {
		p.continueAfterErrors = true
		return
	}
	em := c.errorManager
	p.cachedErrorManager = em
	if ts, ok := em.(*ThreadSafeDelegatingErrorManager); ok {
		em = ts.delegated
	}
	if sem, ok := em.(SortingErrorManager); ok {
		p.haltingErrorCount = &sem.GetSortingErrorManagerBase().originalErrorCount
	}
}

func (p *CombinedCompilerPass) hasHaltingErrors() bool {
	if p.concreteCompiler != nil {
		if p.continueAfterErrors {
			return false
		}
		if p.haltingErrorCount != nil && p.concreteCompiler.errorManager == p.cachedErrorManager {
			return *p.haltingErrorCount != 0
		}
	}
	return p.compiler.HasHaltingErrors()
}

// Java: CombinedCompilerPass.traverse(AbstractCompiler, Node, List<NodeTraversal.Callback>)
func CombinedCompilerPass_Traverse(compiler AbstractCompiler, root *rhino.Node, callbacks []NodeTraversal_Callback) {
	if len(callbacks) == 1 {
		NodeTraversal_Traverse(compiler, root, callbacks[0])
	} else {
		NewCombinedCompilerPassOfList(compiler, callbacks).Process(nil, root)
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
	waiting *rhino.Node

	noShouldTraverse bool
	noVisit          bool
}

// Java: CombinedCompilerPass.CallbackWrapper(NodeTraversal.Callback)
func initCombinedCompilerPass_CallbackWrapper(callback NodeTraversal_Callback) combinedCompilerPass_CallbackWrapper {
	w := combinedCompilerPass_CallbackWrapper{callback: callback}
	if scopedCallback, ok := callback.(NodeTraversal_ScopedCallback); ok {
		w.scopedCallback = scopedCallback
	}
	_, w.noShouldTraverse = callback.(nodeTraversal_noShouldTraverse)
	if !w.noShouldTraverse {
		_, w.noShouldTraverse = callback.(NodeTraversal_AbstractPostOrderCallbackInterface)
	}
	_, w.noVisit = callback.(nodeTraversal_noVisit)
	return w
}

// VisitOrMaybeActivate visits the node unless the wrapped callback is
// inactive. Activates the callback if appropriate.
//
// Java: CombinedCompilerPass.CallbackWrapper.visitOrMaybeActivate(NodeTraversal, Node, Node)
func (w *combinedCompilerPass_CallbackWrapper) VisitOrMaybeActivate(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if w.IsActive() {
		if !w.noVisit {
			w.callback.Visit(t, n, parent)
		}
	} else if w.waiting == n {
		w.waiting = nil
	}
}

// Java: CombinedCompilerPass.CallbackWrapper.shouldTraverseIfActive(NodeTraversal, Node, Node)
func (w *combinedCompilerPass_CallbackWrapper) ShouldTraverseIfActive(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if w.IsActive() && !w.noShouldTraverse && !w.callback.ShouldTraverse(t, n, parent) {
		w.waiting = n
	}
}

// Java: CombinedCompilerPass.CallbackWrapper.enterScopeIfActive(NodeTraversal)
func (w *combinedCompilerPass_CallbackWrapper) EnterScopeIfActive(t *NodeTraversal) {
	if w.IsActive() && w.scopedCallback != nil {
		w.scopedCallback.EnterScope(t)
	}
}

// Java: CombinedCompilerPass.CallbackWrapper.exitScopeIfActive(NodeTraversal)
func (w *combinedCompilerPass_CallbackWrapper) ExitScopeIfActive(t *NodeTraversal) {
	if w.IsActive() && w.scopedCallback != nil {
		w.scopedCallback.ExitScope(t)
	}
}

// Java: CombinedCompilerPass.CallbackWrapper.isActive()
func (w *combinedCompilerPass_CallbackWrapper) IsActive() bool {
	return w.waiting == nil
}

// Java: CombinedCompilerPass.process(Node, Node)
func (p *CombinedCompilerPass) Process(externs *rhino.Node, root *rhino.Node) {
	p.inactiveCount = 0
	p.refreshErrorState()
	if p.processParallel(root) {
		return
	}
	NodeTraversal_Traverse(p.compiler, root, p)
}

func (p *CombinedCompilerPass) canCloneForParallel(root *rhino.Node) (rootClean bool, ok bool) {
	for i := range p.callbacks {
		switch cb := p.callbacks[i].callback.(type) {
		case *CheckJSDoc,
			*Es6CheckModule,
			*CheckMissingReturn,
			*CheckSuspiciousCode,
			*CheckDuplicateCase,
			*CheckGlobalThis,
			*CheckDebuggerStatement:
		case *CheckUnreachableCode:
			if !cb.hasNoUnreachableCodeFast(root) {
				return false, false
			}
			rootClean = true
		default:
			return false, false
		}
	}
	return rootClean, true
}

func (p *CombinedCompilerPass) cloneForWorker(rc *recordingCompiler, rootClean bool) *CombinedCompilerPass {
	cloned := make([]NodeTraversal_Callback, len(p.callbacks))
	var origUC *CheckUnreachableCode
	var clonedUC *CheckUnreachableCode

	for i := range p.callbacks {
		if uc, ok := p.callbacks[i].callback.(*CheckUnreachableCode); ok {
			origUC = uc
			c := NewCheckUnreachableCode(rc)
			c.rootHasNoUnreachable = rootClean
			cloned[i] = c
			clonedUC = c
		}
	}
	for i := range p.callbacks {
		switch cb := p.callbacks[i].callback.(type) {
		case *CheckUnreachableCode:
			// Already cloned above.
		case *CheckJSDoc:
			c := NewCheckJSDoc(rc)
			c.inExterns = cb.inExterns
			cloned[i] = c
		case *Es6CheckModule:
			cloned[i] = NewEs6CheckModule(rc)
		case *CheckMissingReturn:
			c := NewCheckMissingReturn(rc)
			if cb.sharedCfgSource != nil && cb.sharedCfgSource == origUC {
				c.sharedCfgSource = clonedUC
			}
			cloned[i] = c
		case *CheckSuspiciousCode:
			cloned[i] = NewCheckSuspiciousCode()
		case *CheckDuplicateCase:
			cloned[i] = NewCheckDuplicateCase(rc)
		case *CheckGlobalThis:
			cloned[i] = NewCheckGlobalThis(rc)
		case *CheckDebuggerStatement:
			cloned[i] = NewCheckDebuggerStatement(rc)
		}
	}
	return NewCombinedCompilerPassOfList(rc, cloned)
}

func (p *CombinedCompilerPass) processParallel(root *rhino.Node) bool {
	scripts, ok := collectRootScripts(root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(p.compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}
	if p.hasHaltingErrors() {
		return true
	}
	rootClean, canClone := p.canCloneForParallel(root)
	if !canClone {
		return false
	}

	prewarmCompilerForParallelReads(p.compiler)
	restoreReg := enableParallelTypeRegistryReads(p.compiler)
	defer restoreReg()

	lazyGlobal := newLazyGlobalScope(p.compiler, root)

	stopOnHalting := !p.continueAfterErrors
	runParallelWork(p.compiler, numWorkers, len(scripts), func(_ int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var workerPass *CombinedCompilerPass
		var sc *SyntacticScopeCreator

		runItem := func(idx int, rc *recordingCompiler) {
			if workerPass == nil || workerRC != rc {
				workerRC = rc
				workerPass = p.cloneForWorker(rc, rootClean)
				sc = NewSyntacticScopeCreator(rc)
			} else {
				workerPass.inactiveCount = 0
				for i := range workerPass.callbacks {
					workerPass.callbacks[i].waiting = nil
				}
			}
			b := NodeTraversal_Builder{
				compiler:     rc,
				callback:     workerPass,
				scopeCreator: sc,
			}
			b.traverseScriptUnderRoot(root, root, scripts[idx], lazyGlobal)
		}
		return runItem, nil
	}, stopOnHalting)
	return true
}

// Java: CombinedCompilerPass.shouldTraverse(NodeTraversal, Node, Node)
func (p *CombinedCompilerPass) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if p.hasHaltingErrors() {
		return false
	}

	if p.inactiveCount == 0 {
		for _, idx := range p.shouldTraverseIndices {
			w := &p.callbacks[idx]
			if !w.callback.ShouldTraverse(t, n, parent) {
				w.waiting = n
				p.inactiveCount++
			}
		}
	} else {
		for _, idx := range p.shouldTraverseIndices {
			w := &p.callbacks[idx]
			if w.waiting == nil && !w.callback.ShouldTraverse(t, n, parent) {
				w.waiting = n
				p.inactiveCount++
			}
		}
	}
	if p.allNoVisit && t != nil {
		t.noVisit = p.inactiveCount == 0
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
func (p *CombinedCompilerPass) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if p.allNoVisit && p.inactiveCount == 0 {
		return
	}
	if p.hasHaltingErrors() {
		return
	}

	if p.inactiveCount == 0 {
		for _, cb := range p.visitCallbacks {
			cb.Visit(t, n, parent)
		}
		return
	}

	for i := range p.callbacks {
		callback := &p.callbacks[i]
		if callback.waiting == nil {
			if !callback.noVisit {
				callback.callback.Visit(t, n, parent)
			}
		} else if callback.waiting == n {
			callback.waiting = nil
			p.inactiveCount--
		}
	}
	if p.allNoVisit && t != nil {
		t.noVisit = p.inactiveCount == 0
	}
}

// Java: CombinedCompilerPass.enterScope(NodeTraversal)
func (p *CombinedCompilerPass) EnterScope(t *NodeTraversal) {
	for _, idx := range p.scopedIndices {
		p.callbacks[idx].EnterScopeIfActive(t)
	}
	p.updateNoShouldTraverse(t)
}

// Java: CombinedCompilerPass.exitScope(NodeTraversal)
func (p *CombinedCompilerPass) ExitScope(t *NodeTraversal) {
	for _, idx := range p.scopedIndices {
		p.callbacks[idx].ExitScopeIfActive(t)
	}
	p.updateNoShouldTraverse(t)
}

func (p *CombinedCompilerPass) updateNoShouldTraverse(t *NodeTraversal) {
	if t == nil || p.hasHaltingErrors() {
		return
	}
	if p.inactiveCount > 0 {
		t.noShouldTraverse = false
		return
	}
	for _, idx := range p.shouldTraverseIndices {
		if uc, ok := p.callbacks[idx].callback.(*CheckUnreachableCode); ok && uc.CurrentScopeHasNoUnreachable() {
			continue
		}
		t.noShouldTraverse = false
		return
	}
	t.noShouldTraverse = true
}
