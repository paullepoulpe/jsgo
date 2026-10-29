// Port of com/google/javascript/jscomp/J2clClinitPrunerPass.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// J2clClinitPrunerPass is an optimization pass to prune J2CL clinits.
//
// Java: public class J2clClinitPrunerPass implements CompilerPass
type J2clClinitPrunerPass struct {
	emptiedOrder             []strpool.ID
	emptiedClinitMethods     map[strpool.ID]rhino.Node
	compiler                 AbstractCompiler
	initialChangedScopeNodes []rhino.Node
}

// Java: J2clClinitPrunerPass(AbstractCompiler, List<Node>)
func NewJ2clClinitPrunerPass(compiler AbstractCompiler, initialChangedScopeNodes []rhino.Node) *J2clClinitPrunerPass {
	return &J2clClinitPrunerPass{
		compiler:                 compiler,
		initialChangedScopeNodes: initialChangedScopeNodes,
	}
}

// Java: J2clClinitPrunerPass.process(Node, Node)
func (p *J2clClinitPrunerPass) Process(externs rhino.Node, root rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	p.removeRedundantClinits(root, p.initialChangedScopeNodes)
	p.pruneEmptyClinits(root, p.initialChangedScopeNodes)

	if len(p.emptiedOrder) == 0 {
		// Since no clinits are pruned, we don't need to look for more opportunities.
		return
	}

	clinitReferences := p.collectClinitReferences(root)
	for {
		newChangedScopes := p.cleanEmptyClinitReferences(clinitReferences)
		p.pruneEmptyClinits(root, newChangedScopes)
		if len(p.emptiedOrder) == 0 {
			break
		}
	}

	// Update the function side-effect markers on the AST.
	// Removing a clinit from a function may make it side-effect free.
	NewPureFunctionIdentifier_Driver(p.compiler).Process(externs, root)
}

func (p *J2clClinitPrunerPass) recordEmptiedClinit(id strpool.ID, replacement rhino.Node) {
	if p.emptiedClinitMethods == nil {
		p.emptiedClinitMethods = make(map[strpool.ID]rhino.Node)
	}
	if _, exists := p.emptiedClinitMethods[id]; !exists {
		p.emptiedOrder = append(p.emptiedOrder, id)
	}
	p.emptiedClinitMethods[id] = replacement
}

func (p *J2clClinitPrunerPass) removeRedundantClinits(root rhino.Node, changedScopeNodes []rhino.Node) {
	var changedRoots []rhino.Node
	if changedScopeNodes != nil {
		changedRoots = NodeUtil_RemoveNestedChangeScopeNodes(
			NodeUtil_GetParentChangeScopeNodes(changedScopeNodes))
	}
	p.traverseScopeRootsWithFallback(
		root, changedRoots, newRedundantClinitPruner(p, changedRoots), true)
	p.traverseScopeRootsWithFallback(
		root, changedScopeNodes, &lookAheadRedundantClinitPruner{p: p}, false)
}

func (p *J2clClinitPrunerPass) pruneEmptyClinits(root rhino.Node, changedScopes []rhino.Node) {
	p.emptiedOrder = p.emptiedOrder[:0]
	clear(p.emptiedClinitMethods)

	p.traverseScopeRootsWithFallback(root, changedScopes, &emptyClinitPruner{p: p}, false)

	// Make sure replacements point to final destination instead of intermediate ones.
	for _, id := range p.emptiedOrder {
		p.emptiedClinitMethods[id] = p.resolveReplacement(p.emptiedClinitMethods[id])
	}
}

func (p *J2clClinitPrunerPass) traverseScopeRootsWithFallback(
	root rhino.Node,
	scopeRoots []rhino.Node,
	callback NodeTraversal_Callback,
	traverseNested bool) {
	if scopeRoots == nil {
		NodeTraversal_Traverse(p.compiler, root, callback)
	} else {
		NodeTraversal_TraverseScopeRoots(p.compiler, scopeRoots, callback, traverseNested)
	}
}

func (p *J2clClinitPrunerPass) resolveReplacement(node rhino.Node) rhino.Node {
	for !node.IsNil() {
		clinitID := j2clClinitCallID(node)
		if clinitID == 0 {
			break
		}
		next, ok := p.emptiedClinitMethods[clinitID]
		if !ok {
			break
		}
		node = next
	}
	return node
}

func (p *J2clClinitPrunerPass) collectClinitReferences(root rhino.Node) map[strpool.ID][]rhino.Node {
	refs := make(map[strpool.ID][]rhino.Node)
	NodeTraversal_Traverse(p.compiler, root, &collectClinitRefsCallback{refs: refs})
	return refs
}

type collectClinitRefsCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	refs map[strpool.ID][]rhino.Node
}

func (c *collectClinitRefsCallback) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if id := j2clClinitCallID(node); id != 0 {
		c.refs[id] = append(c.refs[id], node)
	}
}

func (p *J2clClinitPrunerPass) cleanEmptyClinitReferences(clinitReferences map[strpool.ID][]rhino.Node) []rhino.Node {
	newChangedScopes := []rhino.Node{}
	for _, clinitID := range p.emptiedOrder {
		replacement := p.emptiedClinitMethods[clinitID]
		references := clinitReferences[clinitID]
		delete(clinitReferences, clinitID)
		for _, reference := range references {
			changedScope := ChangeTracker_GetEnclosingChangeScopeRoot(reference.GetParent())
			if replacement.IsNil() {
				NodeUtil_DeleteFunctionCall(reference, p.compiler)
			} else {
				cloned := replacement.CloneTree()
				reference.ReplaceWith(cloned)
				p.compiler.ReportChangeToChangeScope(changedScope)
				if replID := j2clClinitCallID(cloned); replID != 0 {
					clinitReferences[replID] = append(clinitReferences[replID], cloned)
				}
			}
			newChangedScopes = append(newChangedScopes, changedScope)
		}
	}
	return newChangedScopes
}

// clinitBranchTracker tracks clinits guaranteed to have run along the active
// control-flow path using a flat undo log instead of per-branch heap sets.
type clinitBranchTracker struct {
	log   []strpool.ID
	marks []int
	set   map[strpool.ID]struct{}
}

func (ct *clinitBranchTracker) reset() {
	ct.log = ct.log[:0]
	ct.marks = ct.marks[:0]
	if ct.set != nil {
		clear(ct.set)
	}
}

func (ct *clinitBranchTracker) pushBranch() {
	ct.marks = append(ct.marks, len(ct.log))
}

func (ct *clinitBranchTracker) popBranch() {
	mark := ct.marks[len(ct.marks)-1]
	ct.marks = ct.marks[:len(ct.marks)-1]
	if ct.set != nil {
		for _, id := range ct.log[mark:] {
			delete(ct.set, id)
		}
	}
	ct.log = ct.log[:mark]
}

// add records id as called in the current branch and returns true if this is
// the first time id has been seen in the current branch or any ancestor branch.
func (ct *clinitBranchTracker) add(id strpool.ID) bool {
	if ct.set != nil {
		if _, exists := ct.set[id]; exists {
			return false
		}
		ct.set[id] = struct{}{}
		ct.log = append(ct.log, id)
		return true
	}
	if slices.Contains(ct.log, id) {
		return false
	}
	ct.log = append(ct.log, id)
	if len(ct.log) > 16 {
		ct.set = make(map[strpool.ID]struct{}, 32)
		for _, existing := range ct.log {
			ct.set[existing] = struct{}{}
		}
	}
	return true
}

// redundantClinitPruner removes redundant clinit calls inside a method body if
// the clinit is guaranteed to have been called earlier.
type redundantClinitPruner struct {
	p          *J2clClinitPrunerPass
	roots      map[rhino.Node]struct{}
	stateStack []clinitBranchTracker
	freePool   []clinitBranchTracker
	current    clinitBranchTracker
}

func newRedundantClinitPruner(p *J2clClinitPrunerPass, roots []rhino.Node) *redundantClinitPruner {
	r := &redundantClinitPruner{p: p}
	if len(roots) > 0 {
		r.roots = make(map[rhino.Node]struct{}, len(roots))
		for _, root := range roots {
			r.roots[root] = struct{}{}
		}
	}
	return r
}

func (r *redundantClinitPruner) pushIsolatedState() {
	r.stateStack = append(r.stateStack, r.current)
	if n := len(r.freePool); n > 0 {
		r.current = r.freePool[n-1]
		r.freePool = r.freePool[:n-1]
		r.current.reset()
	} else {
		r.current = clinitBranchTracker{}
	}
}

func (r *redundantClinitPruner) popIsolatedState() {
	if len(r.stateStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	r.current.reset()
	r.freePool = append(r.freePool, r.current)
	r.current = r.stateStack[len(r.stateStack)-1]
	r.stateStack = r.stateStack[:len(r.stateStack)-1]
}

func (r *redundantClinitPruner) ShouldTraverse(t *NodeTraversal, node rhino.Node, parent rhino.Node) bool {
	if _, ok := r.roots[node]; ok {
		for i := range r.stateStack {
			r.stateStack[i].reset()
			r.freePool = append(r.freePool, r.stateStack[i])
		}
		r.stateStack = r.stateStack[:0]
		r.current.reset()
	}

	if NodeUtil_IsFunctionDeclaration(node) || node.IsScript() {
		// Function declarations and script roots may execute out of program order.
		r.pushIsolatedState()
	}

	if isNewControlBranch(parent) {
		r.current.pushBranch()
		if fnID := j2clClinitFunctionID(parent); fnID != 0 {
			// Inside a clinit method itself, assume it is already called.
			r.current.add(fnID)
		}
	}

	return true
}

func (r *redundantClinitPruner) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if clinitID := j2clClinitCallID(node); clinitID != 0 && !r.current.add(clinitID) {
		NodeUtil_DeleteFunctionCall(node, r.p.compiler)
	}

	if isNewControlBranch(parent) {
		r.current.popBranch()
	}

	if NodeUtil_IsFunctionDeclaration(node) || node.IsScript() {
		r.popIsolatedState()
	}
}

func isNewControlBranch(n rhino.Node) bool {
	return !n.IsNil() &&
		(NodeUtil_IsControlStructure(n) ||
			n.IsDefaultValue() ||
			n.IsHook() ||
			n.IsAnd() ||
			n.IsOr() ||
			n.IsFunction())
}

// lookAheadRedundantClinitPruner prunes clinit calls which immediately precede
// calls to a static function that calls the same clinit.
type lookAheadRedundantClinitPruner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	p *J2clClinitPrunerPass
}

func (l *lookAheadRedundantClinitPruner) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if !node.IsExprResult() {
		return
	}

	clinitID := j2clClinitCallID(node.GetFirstChild())
	if clinitID == 0 {
		return
	}

	callOrNewNode := getCallOrNewNode(node.GetNext())
	if callOrNewNode.IsNil() || !callOrNewNode.GetFirstChild().IsName() {
		return
	}

	enclosingFunction := NodeUtil_GetEnclosingFunction(node)
	if enclosingFunction.IsNil() {
		return
	}
	calleeName := callOrNewNode.GetFirstChild().GetString()
	if nearest := NodeUtil_GetNearestFunctionName(enclosingFunction); nearest != nil && calleeName == *nearest {
		return
	}

	var_ := t.GetScope().GetVar(calleeName)
	if var_ == nil || var_.GetInitialValue().IsNil() || !var_.GetInitialValue().IsFunction() {
		return
	}

	staticFnNode := var_.GetInitialValue()
	if callsClinit(staticFnNode, clinitID) && hasSafeArguments(t, callOrNewNode) {
		node.Detach()
		l.p.compiler.ReportChangeToEnclosingScope(parent)
	}
}

func hasSafeArguments(t *NodeTraversal, callOrNewNode rhino.Node) bool {
	for child := callOrNewNode.GetSecondChild(); !child.IsNil(); child = child.GetNext() {
		if !NodeUtil_IsLiteralValue(child, false) && !isParameter(t, child) {
			return false
		}
	}
	return true
}

func isParameter(t *NodeTraversal, n rhino.Node) bool {
	if !n.IsName() {
		return false
	}
	return t.GetScope().GetVarByNode(n).GetParentNode().IsParamList()
}

func getCallOrNewNode(n rhino.Node) rhino.Node {
	for !n.IsNil() {
		switch n.GetToken() {
		case rhino.Token_EXPR_RESULT, rhino.Token_RETURN:
			n = n.GetFirstChild()
		case rhino.Token_CALL, rhino.Token_NEW:
			return n
		case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
			if !n.HasOneChild() {
				return rhino.NilNode
			}
			n = n.GetFirstFirstChild()
		default:
			return rhino.NilNode
		}
	}
	return rhino.NilNode
}

func callsClinit(fnNode rhino.Node, clinitID strpool.ID) bool {
	child := fnNode.GetLastChild().GetFirstChild()
	if child.IsNil() || !child.IsExprResult() {
		return false
	}
	return j2clClinitCallID(child.GetFirstChild()) == clinitID
}

// emptyClinitPruner removes the body of empty clinits.
type emptyClinitPruner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	p *J2clClinitPrunerPass
}

func (e *emptyClinitPruner) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	fnID := j2clClinitFunctionID(node)
	if fnID == 0 {
		return
	}
	e.trySubstituteEmptyFunction(node, fnID)
}

func (e *emptyClinitPruner) trySubstituteEmptyFunction(fnNode rhino.Node, fnID strpool.ID) {
	body := fnNode.GetLastChild()
	if !body.HasChildren() {
		return
	}

	// Ensure that the first expression in the body sets itself to the empty function.
	firstExpr := body.GetFirstChild()
	if !isAssignToEmptyFn(firstExpr, fnID) {
		return
	}

	secondExpr := firstExpr.GetNext()
	var replaceReferencesWith rhino.Node
	if secondExpr.IsNil() {
		// No expressions in clinit so it is a no-op; remove all references.
		replaceReferencesWith = rhino.NilNode
	} else if secondExpr.GetNext().IsNil() &&
		secondExpr.IsExprResult() &&
		j2clClinitCallID(secondExpr.GetFirstChild()) != 0 {
		// Only expression in clinit is a call to another clinit; forward references to it.
		replaceReferencesWith = secondExpr.GetFirstChild()
	} else {
		return
	}

	e.p.recordEmptiedClinit(fnID, replaceReferencesWith)
	NodeUtil_DeleteNode(firstExpr, e.p.compiler)
}

func isAssignToEmptyFn(node rhino.Node, enclosingFnID strpool.ID) bool {
	if !NodeUtil_IsExprAssign(node) {
		return false
	}
	assign := node.GetFirstChild()
	lhs := assign.GetFirstChild()
	if lhs.IsGetElem() {
		return false
	}
	jutil.CheckStateMsg1(lhs.IsName() || lhs.IsGetProp(), "%s", lhs)

	rhs := assign.GetLastChild()
	if !NodeUtil_IsEmptyFunctionExpression(rhs) {
		return false
	}
	if lhs.IsName() && !lhs.GetParent().IsNil() {
		return strpool.ID(lhs.GetStringID()) == enclosingFnID
	}
	lValueName, ok := nodeUtil_getBestLValueNameString(lhs)
	return ok && lValueName == strpool.String(enclosingFnID)
}

func j2clIsClinitName(name string) bool {
	// The '.$clinit' case only happens when collapseProperties is off.
	return strings.HasSuffix(name, "$$0clinit") || strings.HasSuffix(name, ".$clinit")
}

func j2clLeafMayBeClinit(leaf string) bool {
	return leaf == "$clinit" || strings.HasSuffix(leaf, "$$0clinit") || strings.HasSuffix(leaf, ".$clinit")
}

// j2clClinitCallID returns the interned clinit method ID if node is a clinit
// CALL node, or 0 otherwise.
func j2clClinitCallID(node rhino.Node) strpool.ID {
	if !node.IsCall() {
		return 0
	}
	callee := node.GetFirstChild()
	switch callee.GetToken() {
	case rhino.Token_NAME:
		if j2clIsClinitName(callee.GetString()) {
			return strpool.ID(callee.GetStringID())
		}
		return 0
	case rhino.Token_GETPROP:
		if !j2clLeafMayBeClinit(callee.GetString()) {
			return 0
		}
		if qname, ok := callee.GetQualifiedNameString(); ok && j2clIsClinitName(qname) {
			return strpool.Intern(qname)
		}
		return 0
	default:
		fnName, ok := nodeUtil_getBestLValueNameString(callee)
		if ok && j2clIsClinitName(fnName) {
			return strpool.Intern(fnName)
		}
		return 0
	}
}

// j2clClinitFunctionID returns the interned qualified name ID if node is a
// clinit FUNCTION node, or 0 otherwise.
func j2clClinitFunctionID(node rhino.Node) strpool.ID {
	if !node.IsFunction() {
		return 0
	}
	lValue := NodeUtil_GetBestLValue(node)
	if lValue.IsNil() {
		return 0
	}
	switch lValue.GetToken() {
	case rhino.Token_NAME:
		if !lValue.GetParent().IsNil() && j2clIsClinitName(lValue.GetString()) {
			return strpool.ID(lValue.GetStringID())
		}
		return 0
	case rhino.Token_GETPROP, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FIELD_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		if !j2clLeafMayBeClinit(lValue.GetString()) {
			return 0
		}
	}
	fnName, ok := nodeUtil_getBestLValueNameString(lValue)
	if !ok || !j2clIsClinitName(fnName) {
		return 0
	}
	return strpool.Intern(fnName)
}
