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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clClinitPrunerPass is an optimization pass to prune J2CL clinits.
//
// Java: public class J2clClinitPrunerPass implements CompilerPass
type J2clClinitPrunerPass struct {
	emptiedClinitMethods *jutil.LinkedHashMap[string, *rhino.Node]
	compiler             AbstractCompiler
	// initialChangedScopeNodes is a nullable Java List: nil is null.
	initialChangedScopeNodes []*rhino.Node
}

// Java: J2clClinitPrunerPass(AbstractCompiler, List<Node>)
func NewJ2clClinitPrunerPass(compiler AbstractCompiler, initialChangedScopeNodes []*rhino.Node) *J2clClinitPrunerPass {
	return &J2clClinitPrunerPass{
		emptiedClinitMethods:     jutil.NewLinkedHashMap[string, *rhino.Node](),
		compiler:                 compiler,
		initialChangedScopeNodes: initialChangedScopeNodes,
	}
}

// Java: J2clClinitPrunerPass.process(Node, Node)
func (p *J2clClinitPrunerPass) Process(externs *rhino.Node, root *rhino.Node) {
	if !J2clSourceFileChecker_ShouldRunJ2clPasses(p.compiler) {
		return
	}

	p.removeRedundantClinits(root, p.initialChangedScopeNodes)
	p.pruneEmptyClinits(root, p.initialChangedScopeNodes)

	if p.emptiedClinitMethods.IsEmpty() {
		// Since no clinits are pruned, we don't need look for more opportunities.
		return
	}

	clinitReferences := p.collectClinitReferences(root)
	for {
		newChangedScopes := p.cleanEmptyClinitReferences(clinitReferences)
		p.pruneEmptyClinits(root, newChangedScopes)
		if p.emptiedClinitMethods.IsEmpty() {
			break
		}
	}

	// Update the function side-effect markers on the AST.
	// Removing a clinit from a function may make it side-effect free.
	NewPureFunctionIdentifier_Driver(p.compiler).Process(externs, root)
}

// Java: J2clClinitPrunerPass.removeRedundantClinits(Node, List<Node>)
func (p *J2clClinitPrunerPass) removeRedundantClinits(root *rhino.Node, changedScopeNodes []*rhino.Node) {
	changedRoots := j2clClinitPrunerPass_getNonNestedParentScopeNodes(changedScopeNodes)
	p.traverseScopeRootsWithFallback(
		root, changedRoots, p.newRedundantClinitPruner(changedRoots), true)

	p.traverseScopeRootsWithFallback(
		root, changedScopeNodes, p.newLookAheadRedundantClinitPruner(), false)
}

// Java: J2clClinitPrunerPass.pruneEmptyClinits(Node, List<Node>)
func (p *J2clClinitPrunerPass) pruneEmptyClinits(root *rhino.Node, changedScopes []*rhino.Node) {
	// Clear emptiedClinitMethods before EmptyClinitPruner to populate only with new ones.
	p.emptiedClinitMethods.Clear()
	p.traverseScopeRootsWithFallback(root, changedScopes, p.newEmptyClinitPruner(), false)

	// Make sure replacements are to final destination instead of pointing intermediate ones.
	for key, value := range p.emptiedClinitMethods.All() {
		p.emptiedClinitMethods.Put(key, p.resolveReplacement(value)) // entry.setValue
	}
}

// Java: J2clClinitPrunerPass.traverseScopeRootsWithFallback(Node, List<Node>, Callback, boolean)
func (p *J2clClinitPrunerPass) traverseScopeRootsWithFallback(
	root *rhino.Node,
	scopeRoots []*rhino.Node,
	callback NodeTraversal_Callback,
	traverseNested bool) {
	if scopeRoots == nil {
		NodeTraversal_Traverse(p.compiler, root, callback)
	} else {
		NodeTraversal_TraverseScopeRoots(p.compiler, scopeRoots, callback, traverseNested)
	}
}

// Java: J2clClinitPrunerPass.resolveReplacement(Node)
func (p *J2clClinitPrunerPass) resolveReplacement(node *rhino.Node) *rhino.Node {
	if node == nil {
		return nil
	}
	clinitName := j2clClinitPrunerPass_getClinitMethodName(node)
	if clinitName != nil && p.emptiedClinitMethods.ContainsKey(*clinitName) {
		return p.resolveReplacement(p.emptiedClinitMethods.Get(*clinitName))
	}
	return node
}

// Java: J2clClinitPrunerPass.collectClinitReferences(Node)
func (p *J2clClinitPrunerPass) collectClinitReferences(root *rhino.Node) *jutil.LinkedHashMultimap[*string, *rhino.Node] {
	// Go-only: language-adapter — Java HashMultimap: its keys are only looked
	// up (never iterated) and its values are Nodes, which Java hashes by
	// identity, so Java's own value order is not deterministic; insertion
	// order is used.
	clinitReferences := jutil.NewLinkedHashMultimap[*string, *rhino.Node]()
	NodeTraversal_Traverse(
		p.compiler,
		root,
		&j2clClinitPrunerPass_collectClinitReferencesCallback{clinitReferences: clinitReferences})
	return clinitReferences
}

// j2clClinitPrunerPass_collectClinitReferencesCallback is the anonymous
// AbstractPostOrderCallback of collectClinitReferences.
type j2clClinitPrunerPass_collectClinitReferencesCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	clinitReferences *jutil.LinkedHashMultimap[*string, *rhino.Node]
}

// Java: J2clClinitPrunerPass.collectClinitReferences(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *j2clClinitPrunerPass_collectClinitReferencesCallback) Visit(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	clinitName := j2clClinitPrunerPass_getClinitMethodName(node)
	if clinitName != nil {
		c.clinitReferences.Put(clinitName, node)
	}
}

// Java: J2clClinitPrunerPass.cleanEmptyClinitReferences(Multimap<String, Node>)
func (p *J2clClinitPrunerPass) cleanEmptyClinitReferences(clinitReferences *jutil.LinkedHashMultimap[*string, *rhino.Node]) []*rhino.Node {
	newChangedScopes := []*rhino.Node{}
	// resolveReplacement step above should not require any particular iteration order of the
	// emptiedClinitMethods but we are using LinkedHashMap to be extra safe.
	for clinitName, replacement := range p.emptiedClinitMethods.All() {
		references := clinitReferences.RemoveAll(&clinitName)
		for _, reference := range references {
			changedScope := ChangeTracker_GetEnclosingChangeScopeRoot(reference.GetParent())
			if replacement == nil {
				NodeUtil_DeleteFunctionCall(reference, p.compiler)
			} else {
				replacement = replacement.CloneTree()
				reference.ReplaceWith(replacement)
				p.compiler.ReportChangeToChangeScope(changedScope)
				clinitReferences.Put(j2clClinitPrunerPass_getClinitMethodName(replacement), replacement)
			}
			newChangedScopes = append(newChangedScopes, changedScope)
		}
	}
	return newChangedScopes
}

// Java: J2clClinitPrunerPass.getNonNestedParentScopeNodes(List<Node>)
func j2clClinitPrunerPass_getNonNestedParentScopeNodes(changedScopeNodes []*rhino.Node) []*rhino.Node {
	if changedScopeNodes == nil {
		return nil
	}
	return NodeUtil_RemoveNestedChangeScopeNodes(
		NodeUtil_GetParentChangeScopeNodes(changedScopeNodes))
}

	// J2clClinitPrunerPass_RedundantClinitPruner removes redundant clinit calls
// inside method body if it is guaranteed to be called earlier.
//
// Java: private final class J2clClinitPrunerPass.RedundantClinitPruner implements NodeTraversal.Callback
type J2clClinitPrunerPass_RedundantClinitPruner struct {
	outer                 *J2clClinitPrunerPass
	roots                 *jutil.LinkedHashSet[*rhino.Node]
	stateStack            []*J2clClinitPrunerPass_HierarchicalSet[string]
	clinitsCalledAtBranch *J2clClinitPrunerPass_HierarchicalSet[string]
	setPool               []*J2clClinitPrunerPass_HierarchicalSet[string]
}

func (r *J2clClinitPrunerPass_RedundantClinitPruner) allocSet(parent *J2clClinitPrunerPass_HierarchicalSet[string]) *J2clClinitPrunerPass_HierarchicalSet[string] {
	if n := len(r.setPool); n > 0 {
		s := r.setPool[n-1]
		r.setPool = r.setPool[:n-1]
		if s.currentSet != nil && !s.currentSet.IsEmpty() {
			s.currentSet.Clear()
		}
		s.parent = parent
		return s
	}
	return &J2clClinitPrunerPass_HierarchicalSet[string]{parent: parent}
}

func (r *J2clClinitPrunerPass_RedundantClinitPruner) releaseSet(s *J2clClinitPrunerPass_HierarchicalSet[string]) {
	if s == nil {
		return
	}
	s.parent = nil
	r.setPool = append(r.setPool, s)
}

// Java: J2clClinitPrunerPass.RedundantClinitPruner(Iterable<Node>)
func (p *J2clClinitPrunerPass) newRedundantClinitPruner(roots []*rhino.Node) *J2clClinitPrunerPass_RedundantClinitPruner {
	r := &J2clClinitPrunerPass_RedundantClinitPruner{outer: p}
	r.stateStack = nil // new ArrayDeque<>()
	r.clinitsCalledAtBranch = r.allocSet(nil)
	if len(roots) > 0 {
		r.roots = jutil.NewLinkedHashSet(roots...)
	}
	return r
}

// Java: J2clClinitPrunerPass.RedundantClinitPruner.shouldTraverse(NodeTraversal, Node, Node)
func (r *J2clClinitPrunerPass_RedundantClinitPruner) ShouldTraverse(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) bool {
	if r.roots != nil && r.roots.Contains(node) {
		// Reset the clinit call tracking when starting over on a new scope.
		r.clinitsCalledAtBranch = r.allocSet(nil)
		r.stateStack = r.stateStack[:0]
	}

	if NodeUtil_IsFunctionDeclaration(node) || node.IsScript() {
		// In the case of function declarations, unlike function expressions, we don't know when the
		// function will be executed so assume there are no clinits already executed.
		// In the case of script roots, when using modules, we can not assume the scripts are going
		// to be executed in program order.
		r.stateStack = append(r.stateStack, r.clinitsCalledAtBranch)
		r.clinitsCalledAtBranch = r.allocSet(nil)
	}

	if r.isNewControlBranch(parent) {
		r.clinitsCalledAtBranch = r.allocSet(r.clinitsCalledAtBranch)
		if j2clClinitPrunerPass_isClinitMethod(parent) {
			// Adds itself as any of your children can assume clinit is already called.
			r.clinitsCalledAtBranch.Add(*j2clClinitPrunerPass_getQualifiedNameOfFunction(parent))
		}
	}

	return true
}

// Java: J2clClinitPrunerPass.RedundantClinitPruner.visit(NodeTraversal, Node, Node)
func (r *J2clClinitPrunerPass_RedundantClinitPruner) Visit(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	r.tryRemovingClinit(node)

	if r.isNewControlBranch(parent) {
		old := r.clinitsCalledAtBranch
		r.clinitsCalledAtBranch = old.parent
		r.releaseSet(old)
	}

	if NodeUtil_IsFunctionDeclaration(node) || node.IsScript() {
		// In the case of function declarations, unlike function expressions, we don't know when the
		// function will be executed so assume there are no clinits already executed.
		// In the case of script roots, when using modules, we can not assume the scripts are going
		// to be executed in program order.
		if len(r.stateStack) == 0 {
			panic(jutil.NewNoSuchElementException()) // ArrayDeque.removeLast()
		}
		r.releaseSet(r.clinitsCalledAtBranch)
		r.clinitsCalledAtBranch = r.stateStack[len(r.stateStack)-1]
		r.stateStack = r.stateStack[:len(r.stateStack)-1]
	}
}

// Java: J2clClinitPrunerPass.RedundantClinitPruner.tryRemovingClinit(Node)
func (r *J2clClinitPrunerPass_RedundantClinitPruner) tryRemovingClinit(node *rhino.Node) {
	clinitName := j2clClinitPrunerPass_getClinitMethodName(node)
	if clinitName == nil {
		return
	}

	if r.clinitsCalledAtBranch.Add(*clinitName) {
		// This is the first time we are seeing this clinit so cannot remove it.
		return
	}

	NodeUtil_DeleteFunctionCall(node, r.outer.compiler)
}

// Java: J2clClinitPrunerPass.RedundantClinitPruner.isNewControlBranch(Node)
func (r *J2clClinitPrunerPass_RedundantClinitPruner) isNewControlBranch(n *rhino.Node) bool {
	return n != nil &&
		(NodeUtil_IsControlStructure(n) ||
			n.IsDefaultValue() ||
			n.IsHook() ||
			n.IsAnd() ||
			n.IsOr() ||
			n.IsFunction())
}

// TODO(michaelthomas): Prune clinit calls in functions, if previous functions or the immediate
// next function guarantees clinit call. With that we won't need this pass.

// J2clClinitPrunerPass_LookAheadRedundantClinitPruner prunes clinit calls
// which immediately precede calls to a static function which calls the same
// clinit. e.g. "Foo.clinit(); return new Foo()" -> "return new Foo()"
//
// Java: private final class J2clClinitPrunerPass.LookAheadRedundantClinitPruner extends AbstractPostOrderCallback
type J2clClinitPrunerPass_LookAheadRedundantClinitPruner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *J2clClinitPrunerPass
}

// Java: J2clClinitPrunerPass.LookAheadRedundantClinitPruner()
func (p *J2clClinitPrunerPass) newLookAheadRedundantClinitPruner() *J2clClinitPrunerPass_LookAheadRedundantClinitPruner {
	return &J2clClinitPrunerPass_LookAheadRedundantClinitPruner{outer: p}
}

// Java: J2clClinitPrunerPass.LookAheadRedundantClinitPruner.visit(NodeTraversal, Node, Node)
func (l *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) Visit(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	if !node.IsExprResult() {
		return
	}

	// Find clinit calls.
	clinitName := j2clClinitPrunerPass_getClinitMethodName(node.GetFirstChild())
	if clinitName == nil {
		return
	}

	// Check for calls to a static method immediately following the clinit call.
	callOrNewNode := l.getCallOrNewNode(node.GetNext())
	if callOrNewNode == nil || !callOrNewNode.GetFirstChild().IsName() {
		return
	}

	// Check that the call isn't a recursive call to the same function.
	enclosingFunction := NodeUtil_GetEnclosingFunction(node)
	if enclosingFunction == nil {
		return
	}
	if nearest := NodeUtil_GetNearestFunctionName(enclosingFunction); nearest != nil &&
		callOrNewNode.GetFirstChild().GetString() == *nearest {
		return
	}

	// Find the definition of the function being called.
	var_ := t.GetScope().GetVar(callOrNewNode.GetFirstChild().GetString())
	if var_ == nil || var_.GetInitialValue() == nil || !var_.GetInitialValue().IsFunction() {
		return
	}

	// Check that the clinit call is safe to prune.
	staticFnNode := var_.GetInitialValue()
	if l.callsClinit(staticFnNode, *clinitName) && l.hasSafeArguments(t, callOrNewNode) {
		node.Detach()
		l.outer.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// hasSafeArguments returns whether the arguments to the specified call/new
// node are "safe". i.e. they are only parameters to the enclosing function
// or literal values (and not e.g. a static field reference which might need
// the clinit we are trying to remove).
//
// Java: J2clClinitPrunerPass.LookAheadRedundantClinitPruner.hasSafeArguments(NodeTraversal, Node)
func (l *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) hasSafeArguments(t *NodeTraversal, callOrNewNode *rhino.Node) bool {
	child := callOrNewNode.GetSecondChild()
	for child != nil {
		if !NodeUtil_IsLiteralValue(child /* includeFunctions= */, false) &&
			!l.isParameter(t, child) {
			return false
		}
		child = child.GetNext()
	}
	return true
}

// isParameter returns whether the specified node is defined as a parameter
// to its enclosing function.
//
// Java: J2clClinitPrunerPass.LookAheadRedundantClinitPruner.isParameter(NodeTraversal, Node)
func (l *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) isParameter(t *NodeTraversal, n *rhino.Node) bool {
	if !n.IsName() {
		return false
	}
	var_ := t.GetScope().GetVar(n.GetString())
	return var_.GetParentNode().IsParamList()
}

// getCallOrNewNode returns the call node associated with the specified node
// if one exists, otherwise returns null.
//
// Java: J2clClinitPrunerPass.LookAheadRedundantClinitPruner.getCallOrNewNode(Node)
func (l *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) getCallOrNewNode(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	switch n.GetToken() {
	case rhino.Token_EXPR_RESULT, rhino.Token_RETURN:
		return l.getCallOrNewNode(n.GetFirstChild())
	case rhino.Token_CALL, rhino.Token_NEW:
		return n
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		if n.HasOneChild() {
			return l.getCallOrNewNode(n.GetFirstFirstChild())
		}
		return nil
	default:
		return nil
	}
}

// callsClinit returns whether the specified function contains a call to the
// specified clinit.
//
// Java: J2clClinitPrunerPass.LookAheadRedundantClinitPruner.callsClinit(Node, String)
func (l *J2clClinitPrunerPass_LookAheadRedundantClinitPruner) callsClinit(fnNode *rhino.Node, clinitName string) bool {
	// TODO(michaelthomas): Consider checking all children, but watch out for return statements
	// that could short-circuit the clinit.
	child := fnNode.GetLastChild().GetFirstChild()
	if child == nil || !child.IsExprResult() {
		return false
	}
	name := j2clClinitPrunerPass_getClinitMethodName(child.GetFirstChild())
	return name != nil && clinitName == *name
}

// J2clClinitPrunerPass_EmptyClinitPruner is a traversal callback that
// removes the body of empty clinits.
//
// Java: private final class J2clClinitPrunerPass.EmptyClinitPruner extends AbstractPostOrderCallback
type J2clClinitPrunerPass_EmptyClinitPruner struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *J2clClinitPrunerPass
}

// Java: J2clClinitPrunerPass.EmptyClinitPruner()
func (p *J2clClinitPrunerPass) newEmptyClinitPruner() *J2clClinitPrunerPass_EmptyClinitPruner {
	return &J2clClinitPrunerPass_EmptyClinitPruner{outer: p}
}

// Java: J2clClinitPrunerPass.EmptyClinitPruner.visit(NodeTraversal, Node, Node)
func (e *J2clClinitPrunerPass_EmptyClinitPruner) Visit(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	if !j2clClinitPrunerPass_isClinitMethod(node) {
		return
	}

	e.trySubstituteEmptyFunction(node)
}

// trySubstituteEmptyFunction clears the body of any functions that are
// equivalent to empty functions.
//
// Java: J2clClinitPrunerPass.EmptyClinitPruner.trySubstituteEmptyFunction(Node)
func (e *J2clClinitPrunerPass_EmptyClinitPruner) trySubstituteEmptyFunction(fnNode *rhino.Node) {
	fnQualifiedName := j2clClinitPrunerPass_getQualifiedNameOfFunction(fnNode)

	// Ignore anonymous/constructor functions.
	if fnQualifiedName == nil || *fnQualifiedName == "" {
		return
	}

	body := fnNode.GetLastChild()
	if !body.HasChildren() {
		return
	}

	// Ensure that the first expression in the body is setting itself to the empty function
	firstExpr := body.GetFirstChild()
	if !e.isAssignToEmptyFn(firstExpr, *fnQualifiedName) {
		return
	}

	secondExpr := firstExpr.GetNext()
	var replaceReferencesWith *rhino.Node
	if secondExpr == nil {
		// There are no expressions in clinit so it is noop; remove all references.
		replaceReferencesWith = nil
	} else if secondExpr.GetNext() == nil &&
		secondExpr.IsExprResult() &&
		j2clClinitPrunerPass_getClinitMethodName(secondExpr.GetFirstChild()) != nil {
		// Only expression in clinit is a call to another clinit. We can safely replace all
		// references with the other clinit.
		replaceReferencesWith = secondExpr.GetFirstChild()
	} else {
		// Clinit is not empty.
		return
	}

	e.outer.emptiedClinitMethods.Put(*fnQualifiedName, replaceReferencesWith)
	NodeUtil_DeleteNode(firstExpr, e.outer.compiler)
}

// Java: J2clClinitPrunerPass.EmptyClinitPruner.isAssignToEmptyFn(Node, String)
func (e *J2clClinitPrunerPass_EmptyClinitPruner) isAssignToEmptyFn(node *rhino.Node, enclosingFnName string) bool {
	if !NodeUtil_IsExprAssign(node) {
		return false
	}

	lhs := node.GetFirstFirstChild()
	if lhs.IsGetElem() {
		// This can happen if code has been inserted, for example for instrumentation.
		return false
	}
	jutil.CheckStateMsg1(lhs.IsName() || lhs.IsGetProp(), "%s", lhs)

	rhs := node.GetFirstChild().GetLastChild()
	if !NodeUtil_IsEmptyFunctionExpression(rhs) {
		return false
	}
	lValueName := NodeUtil_GetBestLValueName(lhs)
	return lValueName != nil && *lValueName == enclosingFnName
}

// Java: J2clClinitPrunerPass.isClinitMethod(Node)
func j2clClinitPrunerPass_isClinitMethod(node *rhino.Node) bool {
	return node.IsFunction() && j2clClinitPrunerPass_isClinitMethodName(j2clClinitPrunerPass_getQualifiedNameOfFunction(node))
}

// Java: J2clClinitPrunerPass.getClinitMethodName(Node)
func j2clClinitPrunerPass_getClinitMethodName(node *rhino.Node) *string {
	if node.IsCall() {
		fnName := NodeUtil_GetBestLValueName(node.GetFirstChild())
		if j2clClinitPrunerPass_isClinitMethodName(fnName) {
			return fnName
		}
		return nil
	}
	return nil
}

// Java: J2clClinitPrunerPass.isClinitMethodName(String)
func j2clClinitPrunerPass_isClinitMethodName(fnName *string) bool {
	// The '.$clinit' case) only happens when collapseProperties is off.
	return fnName != nil && (strings.HasSuffix(*fnName, "$$0clinit") || strings.HasSuffix(*fnName, ".$clinit"))
}

// J2clClinitPrunerPass_HierarchicalSet is a minimalist implelementation of a
// hierarchical Set where an item might exist in current set or any of its
// parents.
//
// Java: private static class J2clClinitPrunerPass.HierarchicalSet<T>
type J2clClinitPrunerPass_HierarchicalSet[T comparable] struct {
	currentSet *jutil.LinkedHashSet[T]
	parent     *J2clClinitPrunerPass_HierarchicalSet[T]
}

// Java: J2clClinitPrunerPass.HierarchicalSet(HierarchicalSet<T>)
func NewJ2clClinitPrunerPass_HierarchicalSet[T comparable](parent *J2clClinitPrunerPass_HierarchicalSet[T]) *J2clClinitPrunerPass_HierarchicalSet[T] {
	return &J2clClinitPrunerPass_HierarchicalSet[T]{
		parent: parent,
	}
}

// Java: J2clClinitPrunerPass.HierarchicalSet.add(T)
func (s *J2clClinitPrunerPass_HierarchicalSet[T]) Add(o T) bool {
	if s.parentsContains(o) {
		return false
	}
	if s.currentSet == nil {
		s.currentSet = jutil.NewLinkedHashSet[T]()
	}
	return s.currentSet.Add(o)
}

// parentsContains returns true either my parent or any of its parents
// contains the item.
//
// Java: J2clClinitPrunerPass.HierarchicalSet.parentsContains(T)
func (s *J2clClinitPrunerPass_HierarchicalSet[T]) parentsContains(o T) bool {
	for p := s.parent; p != nil; p = p.parent {
		if p.currentSet != nil && p.currentSet.Contains(o) {
			return true
		}
	}
	return false
}

// j2clClinitPrunerPass_getQualifiedNameOfFunction returns the qualifed name
// function is being assigned to.
//
// This implementation abstracts over the various ways of naming a function.
// ASSIGN need not be involved and there may not be a sequence of GETPROPs
// representing the name.
//
// TODO(b/123354857): Delete this method when naming concepts have been
// unified. It was created as a temporary measure to support ES6 in this pass
// without becoming blocked on name APIs.
//
// Java: J2clClinitPrunerPass.getQualifiedNameOfFunction(Node)
func j2clClinitPrunerPass_getQualifiedNameOfFunction(function *rhino.Node) *string {
	jutil.CheckArgumentMsg1(function.IsFunction(), "%s", function)
	// The node representing the name (e.g. GETPROP, MEMBER_FUNCTION_DEF, etc.).
	lValue := NodeUtil_GetBestLValue(function)
	return NodeUtil_GetBestLValueName(lValue)
}
