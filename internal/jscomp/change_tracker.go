// Port of com/google/javascript/jscomp/ChangeTracker.java @ closure-compiler v20261021

/*
 * Copyright 2025 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ChangeTracker tracks various kind of changes during a single
// compilation.
type ChangeTracker struct {
	// Starts at 0, increases as "interesting" things happen. Nothing happens
	// at time START_TIME, the first pass starts at time 1. The correctness
	// of scope-change tracking relies on Node/getIntProp returning 0 if the
	// custom attribute on a node hasn't been set.
	changeStamp int32

	changeTimeline *Timeline[rhino.Node]

	recentChange *RecentChange

	codeChangeHandlers []CodeChangeHandler
}

// NewChangeTracker is new ChangeTracker().
func NewChangeTracker() *ChangeTracker {
	return &ChangeTracker{
		changeStamp:    1,
		changeTimeline: NewTimeline[rhino.Node](),
		recentChange:   NewRecentChange(),
	}
}

// AddChangeHandler registers a listener for code change events.
//
// Java: ChangeTracker.addChangeHandler(CodeChangeHandler)
func (c *ChangeTracker) AddChangeHandler(handler CodeChangeHandler) {
	c.codeChangeHandlers = append(c.codeChangeHandlers, handler)
}

// RemoveChangeHandler removes a listener for code change events.
//
// Java: ChangeTracker.removeChangeHandler(CodeChangeHandler)
func (c *ChangeTracker) RemoveChangeHandler(handler CodeChangeHandler) {
	if i := slices.Index(c.codeChangeHandlers, handler); i >= 0 {
		c.codeChangeHandlers = slices.Delete(c.codeChangeHandlers, i, i+1)
	}
}

// Java: ChangeTracker.getRecentChange()
func (c *ChangeTracker) GetRecentChange() *RecentChange {
	return c.recentChange
}

// ReportChangeToEnclosingScope marks modifications to the enclosing change
// scope, as defined by isChangeScopeRoot(Node).
//
// Java: ChangeTracker.reportChangeToEnclosingScope(Node)
func (c *ChangeTracker) ReportChangeToEnclosingScope(n rhino.Node) {
	c.recordChange(c.getChangeScopeForNode(n))
	c.notifyChangeHandlers()
}

// ReportChangeToChangeScope marks modifications to a function or script
// node.
//
// Java: ChangeTracker.reportChangeToChangeScope(Node)
func (c *ChangeTracker) ReportChangeToChangeScope(changeScopeRoot rhino.Node) {
	jutil.CheckState(changeScopeRoot.IsScript() || changeScopeRoot.IsFunction())
	c.recordChange(changeScopeRoot)
	c.notifyChangeHandlers()
}

// MarkNewScopesChanged recurses through a tree, marking all function nodes
// as changed.
//
// Java: ChangeTracker.markNewScopesChanged(Node)
func (c *ChangeTracker) MarkNewScopesChanged(node rhino.Node) {
	if node.IsFunction() {
		c.ReportChangeToChangeScope(node)
	}
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		c.MarkNewScopesChanged(child)
	}
}

// ReportFunctionDeleted marks a specific function node as known to be
// deleted. Is part of having accurate change tracking which is necessary
// to streamline optimizations.
//
// Java: ChangeTracker.reportFunctionDeleted(Node)
func (c *ChangeTracker) ReportFunctionDeleted(n rhino.Node) {
	jutil.CheckState(n.IsFunction())
	n.SetDeleted(true)
	c.changeTimeline.Remove(n)
}

// GetChangedScopeNodesForPass returns an accumulation of changed scope
// nodes since the last time the given pass was run.
//
// A returned empty list means no scope nodes have changed since the last
// run and a returned null (nil) means this is the first time the pass has
// run.
//
// Java: ChangeTracker.getChangedScopeNodesForPass(String)
func (c *ChangeTracker) GetChangedScopeNodesForPass(passName string) []rhino.Node {
	changedScopeNodes := c.changeTimeline.GetSince(passName)
	c.changeTimeline.Mark(passName)
	return changedScopeNodes
}

// GetChangeStamp returns a monotonically increasing value to identify a
// change.
//
// Java: ChangeTracker.getChangeStamp()
func (c *ChangeTracker) GetChangeStamp() int32 {
	return c.changeStamp
}

// IncrementChangeStamp indicates that the current change stamp has been
// used.
//
// Java: ChangeTracker.incrementChangeStamp()
func (c *ChangeTracker) IncrementChangeStamp() {
	c.changeStamp++
}

// ResetChangeStamp resets the change stamp to 1.
//
// Java: ChangeTracker.resetChangeStamp()
func (c *ChangeTracker) ResetChangeStamp() {
	c.changeStamp = 1
}

// ChangeTracker_IsChangeScopeRoot: a change scope does not directly
// correspond to a language scope but is an internal grouping of changes.
//
// Returns whether the node represents a change scope root.
//
// Java: ChangeTracker.isChangeScopeRoot(Node)
func ChangeTracker_IsChangeScopeRoot(n rhino.Node) bool {
	return n.IsScript() || n.IsFunction()
}

// ChangeTracker_GetEnclosingChangeScopeRoot returns the change scope root.
//
// Java: ChangeTracker.getEnclosingChangeScopeRoot(Node)
func ChangeTracker_GetEnclosingChangeScopeRoot(n rhino.Node) rhino.Node {
	for !n.IsNil() && !ChangeTracker_IsChangeScopeRoot(n) {
		n = n.GetParent()
	}
	return n
}

// Java: ChangeTracker.notifyChangeHandlers()
func (c *ChangeTracker) notifyChangeHandlers() {
	for _, handler := range c.codeChangeHandlers {
		handler.ReportChange()
	}
}

// Java: ChangeTracker.getChangeScopeForNode(Node)
func (c *ChangeTracker) getChangeScopeForNode(n rhino.Node) rhino.Node {
	// Compiler change reporting usually occurs after the AST change has
	// already occurred. In the case of node removals those nodes are
	// already removed from the tree and so have no parent chain to walk. In
	// these situations changes are reported instead against what (used to
	// be) their parent. If that parent is itself a script node then it's
	// important to be able to recognize it as the enclosing scope without
	// first stepping to its parent as well.
	if n.IsScript() {
		return n
	}

	enclosingScopeNode := ChangeTracker_GetEnclosingChangeScopeRoot(n.GetParent())
	if enclosingScopeNode.IsNil() {
		panic(jutil.NewIllegalStateExceptionMsg(
			"An enclosing scope is required for change reports but node " + n.String() + " doesn't have one."))
	}
	return enclosingScopeNode
}

// Java: ChangeTracker.recordChange(Node)
func (c *ChangeTracker) recordChange(n rhino.Node) {
	if n.IsDeleted() {
		// Some complicated passes (like SmartNameRemoval) might both change
		// and delete a scope in the same pass, and they might even perform
		// the change after the deletion because of internal queueing. Just
		// ignore the spurious attempt to mark changed after already marking
		// deleted. There's no danger of deleted nodes persisting in the AST
		// since this is enforced separately in ChangeVerifier.
		return
	}

	n.SetChangeTime(c.changeStamp)
	// Every code change happens at a different time
	c.IncrementChangeStamp()
	c.changeTimeline.Add(n)
}
