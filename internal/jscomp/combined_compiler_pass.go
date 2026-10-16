// Port of com/google/javascript/jscomp/CombinedCompilerPass.java @ closure-compiler v20260415

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
	callbacks []*combinedCompilerPass_CallbackWrapper
	compiler  AbstractCompiler
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
	p := &CombinedCompilerPass{compiler: compiler}
	p.callbacks = make([]*combinedCompilerPass_CallbackWrapper, len(callbacks))
	for i := 0; i < len(callbacks); i++ {
		p.callbacks[i] = newCombinedCompilerPass_CallbackWrapper(callbacks[i])
	}
	return p
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
}

// Java: CombinedCompilerPass.CallbackWrapper(NodeTraversal.Callback)
func newCombinedCompilerPass_CallbackWrapper(callback NodeTraversal_Callback) *combinedCompilerPass_CallbackWrapper {
	w := &combinedCompilerPass_CallbackWrapper{callback: callback}
	if scopedCallback, ok := callback.(NodeTraversal_ScopedCallback); ok {
		w.scopedCallback = scopedCallback
	} else {
		w.scopedCallback = nil
	}
	return w
}

// VisitOrMaybeActivate visits the node unless the wrapped callback is
// inactive. Activates the callback if appropriate.
//
// Java: CombinedCompilerPass.CallbackWrapper.visitOrMaybeActivate(NodeTraversal, Node, Node)
func (w *combinedCompilerPass_CallbackWrapper) VisitOrMaybeActivate(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if w.IsActive() {
		w.callback.Visit(t, n, parent)
	} else if w.waiting == n {
		w.waiting = nil
	}
}

// Java: CombinedCompilerPass.CallbackWrapper.shouldTraverseIfActive(NodeTraversal, Node, Node)
func (w *combinedCompilerPass_CallbackWrapper) ShouldTraverseIfActive(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if w.IsActive() && !w.callback.ShouldTraverse(t, n, parent) {
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
	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: CombinedCompilerPass.shouldTraverse(NodeTraversal, Node, Node)
func (p *CombinedCompilerPass) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if p.compiler.HasHaltingErrors() {
		return false
	}

	for _, callback := range p.callbacks {
		callback.ShouldTraverseIfActive(t, n, parent)
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
	if p.compiler.HasHaltingErrors() {
		return
	}

	for _, callback := range p.callbacks {
		callback.VisitOrMaybeActivate(t, n, parent)
	}
}

// Java: CombinedCompilerPass.enterScope(NodeTraversal)
func (p *CombinedCompilerPass) EnterScope(t *NodeTraversal) {
	for _, callback := range p.callbacks {
		callback.EnterScopeIfActive(t)
	}
}

// Java: CombinedCompilerPass.exitScope(NodeTraversal)
func (p *CombinedCompilerPass) ExitScope(t *NodeTraversal) {
	for _, callback := range p.callbacks {
		callback.ExitScopeIfActive(t)
	}
}
