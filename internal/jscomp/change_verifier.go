// Port of com/google/javascript/jscomp/ChangeVerifier.java @ closure-compiler v20261104

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ChangeVerifier is a class to assist in AST change tracking verification.
// To validate a "snapshot" is taken before and "checkRecordedChanges" at the
// desired check point.
type ChangeVerifier struct {
	changeTracker   *ChangeTracker
	clonesByCurrent map[rhino.Node]rhino.Node
	currentByClone  map[rhino.Node]rhino.Node
	snapshotChange  int32
}

// Java: ChangeVerifier(AbstractCompiler)
func NewChangeVerifier(compiler AbstractCompiler) *ChangeVerifier {
	return &ChangeVerifier{
		clonesByCurrent: make(map[rhino.Node]rhino.Node),
		currentByClone:  make(map[rhino.Node]rhino.Node),
		changeTracker:   compiler.GetChangeTracker(),
	}
}

// Java: ChangeVerifier.snapshot(Node)
func (v *ChangeVerifier) Snapshot(root rhino.Node) *ChangeVerifier {
	// remove any existing snapshot data.
	clear(v.clonesByCurrent)
	clear(v.currentByClone)
	v.snapshotChange = v.changeTracker.GetChangeStamp()

	snapshot := root.CloneTree()
	v.AssociateClones(root, snapshot)
	return v
}

// Java: ChangeVerifier.checkRecordedChanges(Node)
func (v *ChangeVerifier) CheckRecordedChanges(current rhino.Node) {
	v.CheckRecordedChangesWithPassName("", current)
}

// Java: ChangeVerifier.checkRecordedChanges(String, Node)
func (v *ChangeVerifier) CheckRecordedChangesWithPassName(passName string, root rhino.Node) {
	v.VerifyScopeChangesHaveBeenRecorded(passName, root)
}

// AssociateClones: given an AST and its copy, map the root node of each
// scope of main to the corresponding root node of clone
//
// Java: ChangeVerifier.associateClones(Node, Node)
func (v *ChangeVerifier) AssociateClones(n rhino.Node, snapshot rhino.Node) {
	// TODO(johnlenz): determine if MODULE_BODY is useful here.
	if n.IsRoot() || ChangeTracker_IsChangeScopeRoot(n) {
		v.clonesByCurrent[n] = snapshot
		v.currentByClone[snapshot] = n
	}

	child := n.GetFirstChild()
	snapshotChild := snapshot.GetFirstChild()
	for !child.IsNil() {
		v.AssociateClones(child, snapshotChild)
		child = child.GetNext()
		snapshotChild = snapshotChild.GetNext()
	}
}

// VerifyScopeChangesHaveBeenRecorded checks that the scope roots marked as
// changed have indeed changed
//
// Java: ChangeVerifier.verifyScopeChangesHaveBeenRecorded(String, Node)
func (v *ChangeVerifier) VerifyScopeChangesHaveBeenRecorded(passName string, root rhino.Node) {
	passNameMsg := ""
	if passName != "" {
		passNameMsg = passName + ": "
	}

	// Gather all the scope nodes that existed when the snapshot was taken.
	snapshotScopeNodes := make(map[rhino.Node]struct{})
	NodeUtil_VisitPreOrder(
		v.clonesByCurrent[root],
		func(oldNode rhino.Node) {
			if ChangeTracker_IsChangeScopeRoot(oldNode) {
				snapshotScopeNodes[oldNode] = struct{}{}
			}
		})

	NodeUtil_VisitPreOrder(
		root,
		func(n rhino.Node) {
			if n.IsRoot() {
				v.VerifyRoot(n)
			} else if ChangeTracker_IsChangeScopeRoot(n) {
				clone := v.clonesByCurrent[n]
				// Remove any scope nodes that still exist.
				delete(snapshotScopeNodes, clone)
				v.VerifyNode(passNameMsg, n)
				if clone.IsNil() {
					v.VerifyNewNode(passNameMsg, n)
				} else {
					v.VerifyNodeChange(passNameMsg, n, clone)
				}
			}
		})

	// Only actually deleted snapshot scope nodes should remain.
	v.VerifyDeletedScopeNodes(passNameMsg, snapshotScopeNodes)
}

// Java: ChangeVerifier.verifyDeletedScopeNodes(String, Set<Node>)
func (v *ChangeVerifier) VerifyDeletedScopeNodes(passNameMsg string, deletedScopeNodes map[rhino.Node]struct{}) {
	// order-independent: verification passes iff every remaining snapshot scope node is marked deleted
	for snapshotScopeNode := range deletedScopeNodes {
		currentNode := v.currentByClone[snapshotScopeNode]
		// If the original was marked deleted, that's fine, ignore it.
		if currentNode.IsDeleted() {
			continue
		}

		// But if it was deleted but not marked deleted, that's a problem.
		panic(jutil.NewIllegalStateExceptionMsg(
			passNameMsg + "deleted scope was not reported:\n" + currentNode.ToStringTree()))
	}
}

// Java: ChangeVerifier.verifyNode(String, Node)
func (v *ChangeVerifier) VerifyNode(passNameMsg string, n rhino.Node) {
	if n.IsDeleted() {
		panic(jutil.NewIllegalStateExceptionMsg(
			passNameMsg + "existing scope is improperly marked as deleted:\n" + n.ToStringTree()))
	}
}

// Java: ChangeVerifier.verifyNewNode(String, Node)
func (v *ChangeVerifier) VerifyNewNode(passNameMsg string, n rhino.Node) {
	changeTime := n.GetChangeTime()
	if changeTime == 0 || changeTime < v.snapshotChange {
		panic(jutil.NewIllegalStateExceptionMsg(
			passNameMsg + "new scope not explicitly marked as changed:\n" + n.ToStringTree()))
	}
}

// Java: ChangeVerifier.verifyRoot(Node)
func (v *ChangeVerifier) VerifyRoot(root rhino.Node) {
	jutil.CheckState(root.IsRoot())
	if root.GetChangeTime() != 0 {
		panic(jutil.NewIllegalStateExceptionMsg("Root nodes should never be marked as changed."))
	}
}

// Java: ChangeVerifier.verifyNodeChange(String, Node, Node)
func (v *ChangeVerifier) VerifyNodeChange(passNameMsg string, n rhino.Node, snapshot rhino.Node) {
	if n.IsRoot() {
		return
	}
	result := inequivalenceReasonExcludingFunctions(n, snapshot)
	if n.GetChangeTime() > snapshot.GetChangeTime() {
		// If the current node is marked as changed (changeTime >
		// snapshot.getChangeTime) but is actually equal to the snapshot,
		// that's an error.
		if result.equals {
			panic(jutil.NewIllegalStateExceptionMsg(
				passNameMsg + "unchanged scope marked as changed: " + v.GetNameForNode(n)))
		}
	} else {
		// If the current node is NOT marked as changed (changeTime <=
		// snapshot.getChangeTime) but is actually different from the
		// snapshot, that's an error.
		if !result.equals {
			panic(jutil.NewIllegalStateExceptionMsg(
				jutil.StringFormat(
					"\"%schanged scope not marked as changed: %s.\n"+
						"%s\n"+
						"Ancestor nodes:\n"+
						"%s\n",
					passNameMsg,
					v.GetNameForNode(n),
					result.errorMessage(),
					v.Path(n, result.errorNode))))
		}
	}
}

// Java: ChangeVerifier.getNameForNode(Node)
func (v *ChangeVerifier) GetNameForNode(n rhino.Node) string {
	sourceName := jutil.FormatOptionalString(NodeUtil_SourceName(n))
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return "SCRIPT: " + sourceName
	case rhino.Token_FUNCTION:
		name, ok := NodeUtil_NearestFunctionName(n)
		if !ok {
			name = "anonymous@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
		}
		return "FUNCTION: " + name + " in " + sourceName
	default:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected Node type"))
	}
}

// Path returns the path from the ancestor to the child node.
//
// Java: ChangeVerifier.path(Node, Node)
func (v *ChangeVerifier) Path(ancestor rhino.Node, child rhino.Node) string {
	var childToAncestor []string
	for current := child; current != ancestor; current = current.GetParent() {
		childToAncestor = append(childToAncestor, current.String())
	}
	childToAncestor = append(childToAncestor, ancestor.String())
	slices.Reverse(childToAncestor)
	var result strings.Builder
	for i, elem := range childToAncestor {
		for k := 0; k < i*2; k++ {
			result.WriteByte(' ')
		}
		result.WriteString(elem)
		result.WriteByte('\n')
	}
	return result.String()
}

// changeEqualsResult corresponds to Java's private record EqualsResult(boolean
// equals, Node errorNode, Supplier<String> errorMessage), stored by value and
// formatting its message lazily only when an inequivalence is reported.
type changeEqualsResult struct {
	equals    bool
	errorNode rhino.Node
	reason    string
	before    any
	after     any
}

func (r changeEqualsResult) errorMessage() string {
	if r.equals {
		return ""
	}
	return jutil.StringFormat(
		"%s\n"+
			"Before: %s\n"+
			"After:  %s\n",
		r.reason, r.before, r.after)
}

func equalChangeResult() changeEqualsResult {
	return changeEqualsResult{equals: true}
}

// Java: ChangeVerifier.EqualsResult.notEqual(Node, String, Object, Object)
func notEqualChangeResult(errorNode rhino.Node, reason string, after any, before any) changeEqualsResult {
	return changeEqualsResult{
		errorNode: errorNode,
		reason:    reason,
		before:    before,
		after:     after,
	}
}

// inequivalenceReasonExcludingFunctions checks whether the two given nodes are
// equivalent, while ignoring differences in descendant functions.
//
// Java: ChangeVerifier.getInequivalenceReasonExcludingFunctions(Node, Node)
func inequivalenceReasonExcludingFunctions(thisNode rhino.Node, thatNode rhino.Node) changeEqualsResult {
	rhino.CheckNotNull(thisNode)
	rhino.CheckNotNull(thatNode)
	if thisNode.GetChildCount() != thatNode.GetChildCount() {
		return notEqualChangeResult(
			thisNode,
			"differing child count",
			thisNode.String()+": "+strconv.Itoa(thisNode.GetChildCount()),
			thatNode.String()+": "+strconv.Itoa(thatNode.GetChildCount()))
	}
	if !thisNode.IsEquivalentWithSideEffectsToShallow(thatNode) {
		return notEqualChangeResult(thisNode, "shallow inequivalence", thisNode, thatNode)
	}
	if thisNode.IsFunction() && thatNode.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(thisNode) != NodeUtil_IsFunctionDeclaration(thatNode) {
			return notEqualChangeResult(
				thisNode, "mismatched isFunctionDeclaration", thisNode, thatNode)
		}
	}

	thisChild := thisNode.GetFirstChild()
	thatChild := thatNode.GetFirstChild()
	for !thisChild.IsNil() && !thatChild.IsNil() {
		if thisChild.IsFunction() || thisChild.IsScript() {
			// Don't compare function expression name, parameters or bodies.
			// But do check that that the node is there.
			if thatChild.GetToken() != thisChild.GetToken() {
				return notEqualChangeResult(thisNode, "different tokens", thisChild, thatChild)
			}
			// Only compare function names for function declarations (not
			// function expressions) as they change the outer scope definition.
			if thisFn, ok := ast.AsFunction(thisChild); ok && NodeUtil_IsFunctionDeclaration(thisChild) {
				thisName := thisFn.NameNode().Ident()
				thatName := ast.MustFunction(thatChild).NameNode().Ident()
				if thisName != thatName {
					return notEqualChangeResult(thisNode, "function name changed", thisName, thatName)
				}
			}
		} else {
			result := inequivalenceReasonExcludingFunctions(thisChild, thatChild)
			if !result.equals {
				return result
			}
		}
		thisChild = thisChild.GetNext()
		thatChild = thatChild.GetNext()
	}

	return equalChangeResult()
}
