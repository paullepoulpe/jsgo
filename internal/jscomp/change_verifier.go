// Port of com/google/javascript/jscomp/ChangeVerifier.java @ closure-compiler v20261021

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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ChangeVerifier is a class to assist in AST change tracking verification.
// To validate a "snapshot" is taken before and "checkRecordedChanges" at the
// desired check point.
type ChangeVerifier struct {
	changeTracker   *ChangeTracker
	clonesByCurrent *jutil.BiMap[*rhino.Node, *rhino.Node]
	snapshotChange  int32
}

// Java: ChangeVerifier(AbstractCompiler)
func NewChangeVerifier(compiler AbstractCompiler) *ChangeVerifier {
	return &ChangeVerifier{
		clonesByCurrent: jutil.HashBiMap_Create[*rhino.Node, *rhino.Node](),
		changeTracker:   compiler.GetChangeTracker(),
	}
}

// Java: ChangeVerifier.snapshot(Node)
func (v *ChangeVerifier) Snapshot(root *rhino.Node) *ChangeVerifier {
	// remove any existing snapshot data.
	v.clonesByCurrent.Clear()
	v.snapshotChange = v.changeTracker.GetChangeStamp()

	snapshot := root.CloneTree()
	v.AssociateClones(root, snapshot)
	return v
}

// Java: ChangeVerifier.checkRecordedChanges(Node)
func (v *ChangeVerifier) CheckRecordedChanges(current *rhino.Node) {
	v.CheckRecordedChangesWithPassName("", current)
}

// Java: ChangeVerifier.checkRecordedChanges(String, Node)
func (v *ChangeVerifier) CheckRecordedChangesWithPassName(passName string, root *rhino.Node) {
	v.VerifyScopeChangesHaveBeenRecorded(passName, root)
}

// AssociateClones: given an AST and its copy, map the root node of each
// scope of main to the corresponding root node of clone
//
// Java: ChangeVerifier.associateClones(Node, Node)
func (v *ChangeVerifier) AssociateClones(n *rhino.Node, snapshot *rhino.Node) {
	// TODO(johnlenz): determine if MODULE_BODY is useful here.
	if n.IsRoot() || ChangeTracker_IsChangeScopeRoot(n) {
		v.clonesByCurrent.Put(n, snapshot)
	}

	child := n.GetFirstChild()
	snapshotChild := snapshot.GetFirstChild()
	for child != nil {
		v.AssociateClones(child, snapshotChild)
		child = child.GetNext()
		snapshotChild = snapshotChild.GetNext()
	}
}

// VerifyScopeChangesHaveBeenRecorded checks that the scope roots marked as
// changed have indeed changed
//
// Java: ChangeVerifier.verifyScopeChangesHaveBeenRecorded(String, Node)
func (v *ChangeVerifier) VerifyScopeChangesHaveBeenRecorded(passName string, root *rhino.Node) {
	passNameMsg := ""
	if passName != "" {
		passNameMsg = passName + ": "
	}

	// Gather all the scope nodes that existed when the snapshot was taken.
	snapshotScopeNodes := jutil.NewLinkedHashSet[*rhino.Node]()
	NodeUtil_VisitPreOrder(
		v.clonesByCurrent.Get(root),
		NodeUtil_VisitorFunc(func(oldNode *rhino.Node) {
			if ChangeTracker_IsChangeScopeRoot(oldNode) {
				snapshotScopeNodes.Add(oldNode)
			}
		}))

	NodeUtil_VisitPreOrder(
		root,
		NodeUtil_VisitorFunc(func(n *rhino.Node) {
			if n.IsRoot() {
				v.VerifyRoot(n)
			} else if ChangeTracker_IsChangeScopeRoot(n) {
				clone := v.clonesByCurrent.Get(n)
				// Remove any scope nodes that still exist.
				snapshotScopeNodes.Remove(clone)
				v.VerifyNode(passNameMsg, n)
				if clone == nil {
					v.VerifyNewNode(passNameMsg, n)
				} else {
					v.VerifyNodeChange(passNameMsg, n, clone)
				}
			}
		}))

	// Only actually deleted snapshot scope nodes should remain.
	v.VerifyDeletedScopeNodes(passNameMsg, snapshotScopeNodes)
}

// Java: ChangeVerifier.verifyDeletedScopeNodes(String, Set<Node>)
func (v *ChangeVerifier) VerifyDeletedScopeNodes(passNameMsg string, deletedScopeNodes *jutil.LinkedHashSet[*rhino.Node]) {
	for snapshotScopeNode := range deletedScopeNodes.All() {
		currentNode := v.clonesByCurrent.Inverse().Get(snapshotScopeNode)
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
func (v *ChangeVerifier) VerifyNode(passNameMsg string, n *rhino.Node) {
	if n.IsDeleted() {
		panic(jutil.NewIllegalStateExceptionMsg(
			passNameMsg + "existing scope is improperly marked as deleted:\n" + n.ToStringTree()))
	}
}

// Java: ChangeVerifier.verifyNewNode(String, Node)
func (v *ChangeVerifier) VerifyNewNode(passNameMsg string, n *rhino.Node) {
	changeTime := n.GetChangeTime()
	if changeTime == 0 || changeTime < v.snapshotChange {
		panic(jutil.NewIllegalStateExceptionMsg(
			passNameMsg + "new scope not explicitly marked as changed:\n" + n.ToStringTree()))
	}
}

// Java: ChangeVerifier.verifyRoot(Node)
func (v *ChangeVerifier) VerifyRoot(root *rhino.Node) {
	jutil.CheckState(root.IsRoot())
	if root.GetChangeTime() != 0 {
		panic(jutil.NewIllegalStateExceptionMsg("Root nodes should never be marked as changed."))
	}
}

// Java: ChangeVerifier.verifyNodeChange(String, Node, Node)
func (v *ChangeVerifier) VerifyNodeChange(passNameMsg string, n *rhino.Node, snapshot *rhino.Node) {
	if n.IsRoot() {
		return
	}
	result := changeVerifier_getInequivalenceReasonExcludingFunctions(n, snapshot)
	if n.GetChangeTime() > snapshot.GetChangeTime() {
		// If the current node is marked as changed (changeTime >
		// snapshot.getChangeTime) but is actually equal to the snapshot,
		// that's an error.
		if result.Equals() {
			panic(jutil.NewIllegalStateExceptionMsg(
				passNameMsg + "unchanged scope marked as changed: " + v.GetNameForNode(n)))
		}
	} else {
		// If the current node is NOT marked as changed (changeTime <=
		// snapshot.getChangeTime) but is actually different from the
		// snapshot, that's an error.
		if !result.Equals() {
			panic(jutil.NewIllegalStateExceptionMsg(
				jutil.StringFormat(
					"\"%schanged scope not marked as changed: %s.\n"+
						"%s\n"+
						"Ancestor nodes:\n"+
						"%s\n",
					passNameMsg,
					v.GetNameForNode(n),
					result.errorMessage(),
					v.Path(n, result.ErrorNode()))))
		}
	}
}

// Java: ChangeVerifier.getNameForNode(Node)
func (v *ChangeVerifier) GetNameForNode(n *rhino.Node) string {
	sourceName := NodeUtil_GetSourceName(n)
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return "SCRIPT: " + jutil.ToJavaString(sourceName)
	case rhino.Token_FUNCTION:
		fnName := NodeUtil_GetNearestFunctionName(n)
		if fnName == nil {
			s := "anonymous@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
			fnName = &s
		}
		return "FUNCTION: " + *fnName + " in " + jutil.ToJavaString(sourceName)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected Node type"))
	}
}

// Path returns the path from the ancestor to the child node.
//
// Java: ChangeVerifier.path(Node, Node)
func (v *ChangeVerifier) Path(ancestor *rhino.Node, child *rhino.Node) string {
	childToAncestor := jutil.NewArrayList[string]()
	for current := child; current != ancestor; current = current.GetParent() {
		childToAncestor.Add(current.String())
	}
	childToAncestor.Add(ancestor.String())
	// Collections.reverse(childToAncestor)
	elems := childToAncestor.Slice()
	for i, j := 0, len(elems)-1; i < j; i, j = i+1, j-1 {
		elems[i], elems[j] = elems[j], elems[i]
	}
	result := jutil.NewStringBuilder()
	for i := 0; i < childToAncestor.Size(); i++ {
		// result.repeat(' ', i * 2); // indent
		for k := 0; k < i*2; k++ {
			result.AppendChar(' ')
		}
		result.Append(childToAncestor.Get(i))
		result.Append("\n")
	}
	return result.String()
}

// changeVerifier_EqualsResult is the private record EqualsResult(boolean
// equals, Node errorNode, Supplier<String> errorMessage).
//
// Go-only: language-adapter — the Supplier<String> is a Go func; the
// record's equals/hashCode/toString are never used.
type changeVerifier_EqualsResult struct {
	equals       bool
	errorNode    *rhino.Node
	errorMessage func() *string
}

// Java: ChangeVerifier.EqualsResult.equals()
func (r *changeVerifier_EqualsResult) Equals() bool { return r.equals }

// Java: ChangeVerifier.EqualsResult.errorNode()
func (r *changeVerifier_EqualsResult) ErrorNode() *rhino.Node { return r.errorNode }

// Java: ChangeVerifier.EqualsResult.errorMessage()
func (r *changeVerifier_EqualsResult) ErrorMessage() func() *string { return r.errorMessage }

// Java: ChangeVerifier.EqualsResult.equal()
func changeVerifier_EqualsResult_equal() *changeVerifier_EqualsResult {
	return &changeVerifier_EqualsResult{true, nil, func() *string { return nil }}
}

// Java: ChangeVerifier.EqualsResult.notEqual(Node, String, Object, Object)
func changeVerifier_EqualsResult_notEqual(errorNode *rhino.Node, error_ string, after any, before any) *changeVerifier_EqualsResult {
	return &changeVerifier_EqualsResult{
		false,
		errorNode,
		func() *string {
			s := jutil.StringFormat(
				"%s\n"+
					"Before: %s\n"+
					"After:  %s\n",
				error_, before, after)
			return &s
		}}
}

// changeVerifier_getInequivalenceReasonExcludingFunctions checks whether the
// two given nodes are equivalent, while ignoring differences in descendant
// functions.
//
// Java: ChangeVerifier.getInequivalenceReasonExcludingFunctions(Node, Node)
func changeVerifier_getInequivalenceReasonExcludingFunctions(thisNode *rhino.Node, thatNode *rhino.Node) *changeVerifier_EqualsResult {
	jutil.CheckNotNull(thisNode)
	jutil.CheckNotNull(thatNode)
	if thisNode.GetChildCount() != thatNode.GetChildCount() {
		return changeVerifier_EqualsResult_notEqual(
			thisNode,
			"differing child count",
			thisNode.String()+": "+strconv.Itoa(thisNode.GetChildCount()),
			thatNode.String()+": "+strconv.Itoa(thatNode.GetChildCount()))
	}
	if !thisNode.IsEquivalentWithSideEffectsToShallow(thatNode) {
		return changeVerifier_EqualsResult_notEqual(thisNode, "shallow inequivalence", thisNode, thatNode)
	}
	if thisNode.IsFunction() && thatNode.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(thisNode) != NodeUtil_IsFunctionDeclaration(thatNode) {
			return changeVerifier_EqualsResult_notEqual(
				thisNode, "mismatched isFunctionDeclaration", thisNode, thatNode)
		}
	}

	thisChild := thisNode.GetFirstChild()
	thatChild := thatNode.GetFirstChild()
	for thisChild != nil && thatChild != nil {
		if thisChild.IsFunction() || thisChild.IsScript() {
			// Don't compare function expression name, parameters or bodies.
			// But do check that that the node is there.
			if thatChild.GetToken() != thisChild.GetToken() {
				return changeVerifier_EqualsResult_notEqual(thisNode, "different tokens", thisChild, thatChild)
			}
			// Only compare function names for function declarations (not
			// function expressions) as they change the outer scope definition.
			if thisChild.IsFunction() && NodeUtil_IsFunctionDeclaration(thisChild) {
				thisName := thisChild.GetFirstChild().GetString()
				thatName := thatChild.GetFirstChild().GetString()
				if thisName != thatName {
					return changeVerifier_EqualsResult_notEqual(thisNode, "function name changed", thisName, thatName)
				}
			}
		} else {
			result := changeVerifier_getInequivalenceReasonExcludingFunctions(thisChild, thatChild)
			if !result.Equals() {
				return result
			}
		}
		thisChild = thisChild.GetNext()
		thatChild = thatChild.GetNext()
	}

	return changeVerifier_EqualsResult_equal()
}
