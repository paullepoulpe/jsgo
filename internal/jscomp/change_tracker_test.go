// Port of com/google/javascript/jscomp/ChangeTrackerTest.java @ closure-compiler v20261015

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestChangeTracker_ReportChangeNoScopeFails(t *testing.T) {
	changeTracker := NewChangeTracker()

	detachedNode := rhino.IR_Var(rhino.IR_Name("foo"))

	assertThrows[jutil.IIllegalStateException](t, func() {
		changeTracker.ReportChangeToEnclosingScope(detachedNode)
	})
}

func TestChangeTracker_ReportChangeWithScopeSucceeds(t *testing.T) {
	changeTracker := NewChangeTracker()

	attachedNode := rhino.IR_Var(rhino.IR_Name("foo"))
	rhino.IR_Function(rhino.IR_Name("bar"), rhino.IR_ParamList(), rhino.IR_BlockOfNode(attachedNode))

	// Succeeds without throwing an exception.
	changeTracker.ReportChangeToEnclosingScope(attachedNode)
}

// See TimelineTest for the many timeline behavior tests that don't make
// sense to duplicate here.
func TestChangeTracker_GetChangesAndDeletionsBaseline(t *testing.T) {
	changeTracker := NewChangeTracker()

	// In the initial state nothing has been marked changed or deleted.
	if got := changeTracker.GetChangedScopeNodesForPass("FunctionInliner"); got != nil {
		t.Errorf("got %v, want null", got)
	}
}

func newChangeTrackerTestFunctions() (rhino.Node, rhino.Node) {
	function1 := rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(), rhino.IR_Block())
	function2 := rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(), rhino.IR_Block())
	rhino.IR_Root(rhino.IR_ScriptOfNodes(function1, function2))
	return function1, function2
}

func TestChangeTracker_GetChangesAndDeletionsChangeReportsVisible(t *testing.T) {
	changeTracker := NewChangeTracker()
	function1, function2 := newChangeTrackerTestFunctions()

	// Mark original baseline.
	changeTracker.GetChangedScopeNodesForPass("FunctionInliner")

	// Mark both functions changed.
	changeTracker.ReportChangeToChangeScope(function1)
	changeTracker.ReportChangeToChangeScope(function2)

	// Both function1 and function2 are seen as changed and nothing is seen
	// as deleted.
	assertSameElements(t, "changed", changeTracker.GetChangedScopeNodesForPass("FunctionInliner"),
		[]rhino.Node{function1, function2})
}

func TestChangeTracker_GetChangesAndDeletionsDeleteOverridesChange(t *testing.T) {
	changeTracker := NewChangeTracker()
	function1, function2 := newChangeTrackerTestFunctions()

	// Mark original baseline.
	changeTracker.GetChangedScopeNodesForPass("FunctionInliner")

	// Mark both functions changed, then delete function2 and mark it
	// deleted.
	changeTracker.ReportChangeToChangeScope(function1)
	changeTracker.ReportChangeToChangeScope(function2)
	function2.Detach()
	changeTracker.ReportFunctionDeleted(function2)

	// Now function1 will be seen as changed and function2 will be seen as
	// deleted, since delete overrides change.
	assertSameElements(t, "changed", changeTracker.GetChangedScopeNodesForPass("FunctionInliner"),
		[]rhino.Node{function1})
}

func TestChangeTracker_GetChangesAndDeletionsChangeDoesntOverrideDelete(t *testing.T) {
	changeTracker := NewChangeTracker()
	function1, function2 := newChangeTrackerTestFunctions()

	// Mark original baseline.
	changeTracker.GetChangedScopeNodesForPass("FunctionInliner")

	// Mark function1 changed and function2 deleted, then try to mark
	// function2 changed.
	changeTracker.ReportChangeToChangeScope(function1)
	function2.Detach()
	changeTracker.ReportFunctionDeleted(function2)
	changeTracker.ReportChangeToChangeScope(function2)

	// Now function1 will be seen as changed and function2 will be seen as
	// deleted, since change does not override delete.
	assertSameElements(t, "changed", changeTracker.GetChangedScopeNodesForPass("FunctionInliner"),
		[]rhino.Node{function1})
}

func TestChangeTracker_GetChangesAndDeletionsOnlySeesChangesSinceLastRequest(t *testing.T) {
	changeTracker := NewChangeTracker()
	function1, function2 := newChangeTrackerTestFunctions()

	// Mark original baseline.
	changeTracker.GetChangedScopeNodesForPass("FunctionInliner")

	// Mark function1 changed and function2 deleted.
	changeTracker.ReportChangeToChangeScope(function1)
	function2.Detach()
	changeTracker.ReportFunctionDeleted(function2)

	// Verify their respective states are seen.
	assertSameElements(t, "changed", changeTracker.GetChangedScopeNodesForPass("FunctionInliner"),
		[]rhino.Node{function1})

	// Check states again. Should find nothing since nothing has changed
	// since the last 'FunctionInliner' request.
	got := changeTracker.GetChangedScopeNodesForPass("FunctionInliner")
	if got == nil || len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}
