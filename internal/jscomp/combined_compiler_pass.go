// Port of com/google/javascript/jscomp/CombinedCompilerPass.java @ closure-compiler v20261021

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
	callbacks             []combinedCompilerPass_CallbackWrapper
	shouldTraverseIndices []int
	visitCallbacks        []NodeTraversal_Callback
	scopedIndices         []int
	visitTokenMask        [256]uint32
	unionVisitTokens      ast.TokenSet
	inactiveMask          uint32
	useVisitTokenMask     bool
	hasVisitTokenFilter   bool
	compiler              AbstractCompiler
	cachedErrorManager    ErrorManager
	sortingErrorManager   *SortingErrorManager
	continueAfterErrors   bool
	allNoVisit            bool
	inactiveCount         int
}

var _ nodeTraversal_visitTokens = (*CombinedCompilerPass)(nil)

func (p *CombinedCompilerPass) visitTokens() *ast.TokenSet {
	if p.hasVisitTokenFilter && p.inactiveCount == 0 {
		return &p.unionVisitTokens
	}
	return nil
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
	for i := 0; i < len(callbacks); i++ {
		w := initCombinedCompilerPass_CallbackWrapper(callbacks[i])
		if !w.noShouldTraverse {
			p.shouldTraverseIndices = append(p.shouldTraverseIndices, i)
		}
		if !w.noVisit {
			p.allNoVisit = false
			p.visitCallbacks = append(p.visitCallbacks, w.callback)
			if p.useVisitTokenMask {
				bit := uint32(1) << i
				var ts *ast.TokenSet
				if f, ok := w.callback.(nodeTraversal_visitTokens); ok {
					ts = f.visitTokens()
				}
				if ts != nil {
					for wordIdx := 0; wordIdx < 4; wordIdx++ {
						for word := ts[wordIdx]; word != 0; word &= word - 1 {
							tok := wordIdx*64 + bits.TrailingZeros64(word)
							p.visitTokenMask[tok] |= bit
						}
						p.unionVisitTokens[wordIdx] |= ts[wordIdx]
					}
				} else {
					hasUnfilteredVisit = true
					for tok := 0; tok < 256; tok++ {
						p.visitTokenMask[tok] |= bit
					}
				}
			}
		}
		if w.scopedCallback != nil {
			p.scopedIndices = append(p.scopedIndices, i)
		}
		p.callbacks[i] = w
	}
	p.hasVisitTokenFilter = p.useVisitTokenMask && !hasUnfilteredVisit && !p.allNoVisit
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
	waiting rhino.Node

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
func (w *combinedCompilerPass_CallbackWrapper) VisitOrMaybeActivate(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if w.IsActive() {
		if !w.noVisit {
			w.callback.Visit(t, n, parent)
		}
	} else if w.waiting == n {
		w.waiting = rhino.NilNode
	}
}

// Java: CombinedCompilerPass.CallbackWrapper.shouldTraverseIfActive(NodeTraversal, Node, Node)
func (w *combinedCompilerPass_CallbackWrapper) ShouldTraverseIfActive(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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
	return w.waiting.IsNil()
}

// Java: CombinedCompilerPass.process(Node, Node)
func (p *CombinedCompilerPass) Process(externs rhino.Node, root rhino.Node) {
	p.inactiveCount = 0
	p.inactiveMask = 0
	for i := range p.callbacks {
		p.callbacks[i].waiting = rhino.NilNode
	}
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

func (p *CombinedCompilerPass) canCloneForParallel(root rhino.Node) (rootClean bool, ok bool) {
	for i := range p.callbacks {
		cb := p.callbacks[i].callback
		if _, cloneable := cb.(parallelCloneableCallback); !cloneable {
			return false, false
		}
		if uc, isUC := cb.(*CheckUnreachableCode); isUC {
			if !uc.hasNoUnreachableCodeFast(root) {
				return false, false
			}
			rootClean = true
		}
	}
	return rootClean, true
}

func (p *CombinedCompilerPass) cloneForWorker(rc *Compiler, rootClean bool) *CombinedCompilerPass {
	cloned := make([]NodeTraversal_Callback, len(p.callbacks))
	var origUC *CheckUnreachableCode
	var clonedUC *CheckUnreachableCode

	for i := range p.callbacks {
		cb := p.callbacks[i].callback
		c := cb.(parallelCloneableCallback).cloneForParallel(rc)
		if uc, ok := c.(*CheckUnreachableCode); ok {
			origUC = cb.(*CheckUnreachableCode)
			uc.rootHasNoUnreachable = rootClean
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

func (p *CombinedCompilerPass) processParallel(root rhino.Node) bool {
	var rootClean bool
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
			var canClone bool
			rootClean, canClone = p.canCloneForParallel(root)
			return canClone
		},
		initWorker: func(_ int, rc *Compiler) *CombinedCompilerPass {
			return p.cloneForWorker(rc, rootClean)
		},
		callback: func(wp *CombinedCompilerPass, _ int, _, _ rhino.Node, _ *Compiler) NodeTraversal_Callback {
			wp.inactiveCount = 0
			wp.inactiveMask = 0
			for i := range wp.callbacks {
				wp.callbacks[i].waiting = rhino.NilNode
			}
			return wp
		},
	})
	return ok || halted
}

// Java: CombinedCompilerPass.shouldTraverse(NodeTraversal, Node, Node)
func (p *CombinedCompilerPass) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if p.hasHaltingErrors() {
		return false
	}

	if p.inactiveCount == 0 {
		for _, idx := range p.shouldTraverseIndices {
			w := &p.callbacks[idx]
			if !w.callback.ShouldTraverse(t, n, parent) {
				w.waiting = n
				p.inactiveCount++
				if p.useVisitTokenMask {
					p.inactiveMask |= uint32(1) << idx
				}
			}
		}
	} else {
		for _, idx := range p.shouldTraverseIndices {
			w := &p.callbacks[idx]
			if w.waiting.IsNil() && !w.callback.ShouldTraverse(t, n, parent) {
				w.waiting = n
				p.inactiveCount++
				if p.useVisitTokenMask {
					p.inactiveMask |= uint32(1) << idx
				}
			}
		}
	}
	if t != nil {
		if p.allNoVisit {
			t.noVisit = p.inactiveCount == 0
		} else if p.hasVisitTokenFilter {
			if p.inactiveCount == 0 {
				t.visitTokens = &p.unionVisitTokens
			} else {
				t.visitTokens = nil
			}
		}
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
			if mask == 0 || p.hasHaltingErrors() {
				return
			}
			for ; mask != 0; mask &= mask - 1 {
				idx := bits.TrailingZeros32(mask)
				p.callbacks[idx].callback.Visit(t, n, parent)
			}
			return
		}
		if p.hasHaltingErrors() {
			return
		}
		for _, cb := range p.visitCallbacks {
			cb.Visit(t, n, parent)
		}
		return
	}

	if p.hasHaltingErrors() {
		return
	}
	if p.useVisitTokenMask {
		for mask := p.visitTokenMask[byte(n.GetToken())] | p.inactiveMask; mask != 0; mask &= mask - 1 {
			idx := bits.TrailingZeros32(mask)
			callback := &p.callbacks[idx]
			if callback.waiting.IsNil() {
				callback.callback.Visit(t, n, parent)
			} else if callback.waiting == n {
				callback.waiting = rhino.NilNode
				p.inactiveMask &^= uint32(1) << idx
				p.inactiveCount--
			}
		}
	} else {
		for i := range p.callbacks {
			callback := &p.callbacks[i]
			if callback.waiting.IsNil() {
				if !callback.noVisit {
					callback.callback.Visit(t, n, parent)
				}
			} else if callback.waiting == n {
				callback.waiting = rhino.NilNode
				p.inactiveCount--
			}
		}
	}
	if t != nil {
		if p.allNoVisit {
			t.noVisit = p.inactiveCount == 0
		} else if p.hasVisitTokenFilter && p.inactiveCount == 0 {
			t.visitTokens = &p.unionVisitTokens
		}
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
