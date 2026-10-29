// Port of com/google/javascript/jscomp/NodeUtil.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// NodeUtil_GetVarsDeclaredInBranch retrieves vars declared in the current
// node tree, excluding descent scopes.
//
// Java: NodeUtil.getVarsDeclaredInBranch(Node)
func NodeUtil_GetVarsDeclaredInBranch(root rhino.Node) []rhino.Node {
	var seenBuf [16]uint32
	seen := seenBuf[:0]
	var out []rhino.Node
	var walk func(n rhino.Node)
	walk = func(n rhino.Node) {
		if n.IsVar() {
			NodeUtil_GetLhsNodesHelper(n, func(name rhino.Node) {
				if !name.IsName() {
					_ = name.GetString()
				}
				id := name.GetStringID()
				if !slices.Contains(seen, id) {
					seen = append(seen, id)
					out = append(out, name)
				}
			})
		}
		if !n.IsFunction() {
			for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				walk(c)
			}
		}
	}
	walk(root)
	return out
}

// Java: NodeUtil.getLhsNodesHelper(Node, Consumer<Node>)
func NodeUtil_GetLhsNodesHelper(n rhino.Node, consumer func(rhino.Node)) {
	switch n.GetToken() {
	case rhino.Token_IMPORT:
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
		NodeUtil_GetLhsNodesHelper(n.GetSecondChild(), consumer)
		return
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET, rhino.Token_OBJECT_PATTERN, rhino.Token_ARRAY_PATTERN,
		rhino.Token_PARAM_LIST, rhino.Token_IMPORT_SPECS:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			NodeUtil_GetLhsNodesHelper(child, consumer)
		}
		return
	case rhino.Token_DESTRUCTURING_LHS, rhino.Token_DEFAULT_VALUE, rhino.Token_CATCH, rhino.Token_ITER_REST,
		rhino.Token_OBJECT_REST, rhino.Token_CAST:
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
	case rhino.Token_IMPORT_SPEC, rhino.Token_COMPUTED_PROP, rhino.Token_STRING_KEY:
		NodeUtil_GetLhsNodesHelper(n.GetLastChild(), consumer)
	case rhino.Token_NAME, rhino.Token_IMPORT_STAR:
		consumer(n)
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		// Not valid in declarations but may appear in assignments.
		consumer(n)
	case rhino.Token_EMPTY: // nothing to do.
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// Enhanced for loops assign to variables in their first child
		// e.g.
		// for (some.prop in someObj) {...
		// for ({a, b} of someIterable) {...
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
	default:
		if NodeUtil_IsAssignmentOp(n) {
			NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
		} else {
			panic(jutil.NewIllegalStateExceptionMsg("Invalid node in lhs: " + n.String()))
		}
	}
}

// NodeUtil_MatchDeclaration is a predicate for matching var, let, const,
// class or function declarations.
//
// Java: NodeUtil.MatchDeclaration
type NodeUtil_MatchDeclaration struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: NodeUtil.MatchDeclaration()
func NewNodeUtil_MatchDeclaration() *NodeUtil_MatchDeclaration {
	return &NodeUtil_MatchDeclaration{}
}

// Java: NodeUtil.MatchDeclaration.apply(Node)
func (m *NodeUtil_MatchDeclaration) Apply(n rhino.Node) bool {
	return NodeUtil_IsDeclaration(n)
}

// Java: NodeUtil.MATCH_NOT_FUNCTION
var NodeUtil_MATCH_NOT_FUNCTION = func(n rhino.Node) bool { return !n.IsFunction() }

// NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION is a predicate for matching
// anything except for a non-arrow function.
//
// Useful to avoid traversing into scopes that don't share the same values
// for this, super, or arguments.
//
// Java: NodeUtil.MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION
var NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION = func(n rhino.Node) bool { return !NodeUtil_IsNonArrowFunction(n) }

// NodeUtil_IsShallowStatementTree: whether the given node's subtree may
// contain statements, excepting nested functions.
//
// This is useful for traversing all statements in a given script or
// function without traversing into nested functions, so it will return true
// for e.g. a BLOCK or SWITCH statement.
//
// Java: NodeUtil.isShallowStatementTree(Node)
func NodeUtil_IsShallowStatementTree(n rhino.Node) bool {
	return n.IsNil() ||
		NodeUtil_IsControlStructure(n) ||
		NodeUtil_IsStatementBlock(n) ||
		n.IsSwitchBody()
}

// NodeUtil_MatchShallowStatement is a predicate for matching statements
// without exiting the current scope.
//
// Java: NodeUtil.MatchShallowStatement
type NodeUtil_MatchShallowStatement struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: NodeUtil.MatchShallowStatement()
func NewNodeUtil_MatchShallowStatement() *NodeUtil_MatchShallowStatement {
	return &NodeUtil_MatchShallowStatement{}
}

// Java: NodeUtil.MatchShallowStatement.apply(Node)
func (m *NodeUtil_MatchShallowStatement) Apply(n rhino.Node) bool {
	parent := n.GetParent()
	return n.IsRoot() || n.IsBlock() || (!n.IsFunction() && NodeUtil_IsShallowStatementTree(parent))
}

// NodeUtil_GetNodeTypeReferenceCount finds the number of times a type is
// referenced within the node tree.
//
// Java: NodeUtil.getNodeTypeReferenceCount(Node, Token, Predicate<Node>)
func NodeUtil_GetNodeTypeReferenceCount(
	node rhino.Node, type_ rhino.Token, traverseChildrenPred func(rhino.Node) bool) int {
	total := 0
	if node.GetToken() == type_ {
		total++
	}
	if traverseChildrenPred(node) {
		for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			total += NodeUtil_GetNodeTypeReferenceCount(c, type_, traverseChildrenPred)
		}
	}
	return total
}

// NodeUtil_IsNameReferencedWithPred: whether a simple name is referenced
// within the node tree.
//
// Java: NodeUtil.isNameReferenced(Node, String, Predicate<Node>)
func NodeUtil_IsNameReferencedWithPred(node rhino.Node, name string, traverseChildrenPred func(rhino.Node) bool) bool {
	id := uint32(strpool.Lookup(name))
	if id == 0 {
		_ = node.GetToken()
		_ = traverseChildrenPred(node)
		return false
	}
	return nodeUtil_hasNameIDWithPred(node, id, traverseChildrenPred)
}

func nodeUtil_hasNameIDWithPred(node rhino.Node, id uint32, traverseChildrenPred func(rhino.Node) bool) bool {
	if node.IsName() && node.GetStringID() == id {
		return true
	}
	if !traverseChildrenPred(node) {
		return false
	}
	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if nodeUtil_hasNameIDWithPred(c, id, traverseChildrenPred) {
			return true
		}
	}
	return false
}

// NodeUtil_IsNameReferenced: whether a simple name is referenced within the
// node tree.
//
// Java: NodeUtil.isNameReferenced(Node, String)
func NodeUtil_IsNameReferenced(node rhino.Node, name string) bool {
	_ = node.GetToken()
	id := uint32(strpool.Lookup(name))
	if id == 0 {
		return false
	}
	return nodeUtil_hasNameID(node, id)
}

func nodeUtil_hasNameID(node rhino.Node, id uint32) bool {
	if node.IsName() && node.GetStringID() == id {
		return true
	}
	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if nodeUtil_hasNameID(c, id) {
			return true
		}
	}
	return false
}

// NodeUtil_GetNameReferenceCount finds the number of times a simple name is
// referenced within the node tree.
//
// Java: NodeUtil.getNameReferenceCount(Node, String)
func NodeUtil_GetNameReferenceCount(node rhino.Node, name string) int {
	_ = node.GetToken()
	id := uint32(strpool.Lookup(name))
	if id == 0 {
		return 0
	}
	return nodeUtil_countNameID(node, id)
}

func nodeUtil_countNameID(node rhino.Node, id uint32) int {
	total := 0
	if node.IsName() && node.GetStringID() == id {
		total++
	}
	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		total += nodeUtil_countNameID(c, id)
	}
	return total
}

// nodeUtil_alwaysTrue is Guava's Predicates.alwaysTrue() for nodes.
func nodeUtil_alwaysTrue(rhino.Node) bool { return true }

// NodeUtil_Has returns whether the predicate is true for the node or any of
// its descendants.
//
// Java: NodeUtil.has(Node, Predicate<Node>, Predicate<Node>)
func NodeUtil_Has(node rhino.Node, pred func(rhino.Node) bool, traverseChildrenPred func(rhino.Node) bool) bool {
	if pred(node) {
		return true
	}

	if !traverseChildrenPred(node) {
		return false
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if NodeUtil_Has(c, pred, traverseChildrenPred) {
			return true
		}
	}

	return false
}

// NodeUtil_FindPreorder returns the first Node matching the given pred via a
// pre-order traversal.
//
// Java: NodeUtil.findPreorder(Node, Predicate<Node>, Predicate<Node>)
func NodeUtil_FindPreorder(
	node rhino.Node, pred func(rhino.Node) bool, traverseChildrenPred func(rhino.Node) bool) rhino.Node {
	if pred(node) {
		return node
	}

	if !traverseChildrenPred(node) {
		return rhino.NilNode
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := NodeUtil_FindPreorder(c, pred, traverseChildrenPred)
		if !result.IsNil() {
			return result
		}
	}

	return rhino.NilNode
}

// NodeUtil_GetCount returns the number of times the predicate is true for
// the node or any of its descendants.
//
// Java: NodeUtil.getCount(Node, Predicate<Node>, Predicate<Node>)
func NodeUtil_GetCount(n rhino.Node, pred func(rhino.Node) bool, traverseChildrenPred func(rhino.Node) bool) int {
	total := 0

	if pred(n) {
		total++
	}

	if traverseChildrenPred(n) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			total += NodeUtil_GetCount(c, pred, traverseChildrenPred)
		}
	}

	return total
}

// NodeUtil_VisitPreOrder is a pre-order traversal, calling visitor for each
// decendent.
//
// Java: NodeUtil.visitPreOrder(Node, Visitor)
func NodeUtil_VisitPreOrder(node rhino.Node, visitor func(rhino.Node)) {
	NodeUtil_VisitPreOrderWithPred(node, visitor, nodeUtil_alwaysTrue)
}

// NodeUtil_VisitPreOrderWithPred is a pre-order traversal, calling visitor for
// each node in the tree. Children of nodes that do not match the predicate
// will not be visited.
//
// Java: NodeUtil.visitPreOrder(Node, Visitor, Predicate<Node>)
func NodeUtil_VisitPreOrderWithPred(
	node rhino.Node, visitor func(rhino.Node), traverseChildrenPred func(rhino.Node) bool) {
	visitor(node)

	if traverseChildrenPred(node) {
		for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			NodeUtil_VisitPreOrderWithPred(c, visitor, traverseChildrenPred)
		}
	}
}

// NodeUtil_VisitPostOrder is a post-order traversal, calling visitor for each
// decendent.
//
// Java: NodeUtil.visitPostOrder(Node, Visitor)
func NodeUtil_VisitPostOrder(node rhino.Node, visitor func(rhino.Node)) {
	NodeUtil_VisitPostOrderWithPred(node, visitor, nodeUtil_alwaysTrue)
}

// NodeUtil_VisitPostOrderWithPred is a post-order traversal, calling visitor
// for each node in the tree. Children of nodes that do not match the predicate
// will not be visited.
//
// Java: NodeUtil.visitPostOrder(Node, Visitor, Predicate<Node>)
func NodeUtil_VisitPostOrderWithPred(
	node rhino.Node, visitor func(rhino.Node), traverseChildrenPred func(rhino.Node) bool) {
	if traverseChildrenPred(node) {
		for c := node.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			NodeUtil_VisitPostOrderWithPred(c, visitor, traverseChildrenPred)
			c = next
		}
	}

	visitor(node)
}

// NodeUtil_PreOrderIterableWithPred creates an Iterable over the given node
// and all descendents. Nodes are given in depth-first pre-order (document
// source order,
// https://www.w3.org/TR/DOM-Level-2-Traversal-Range/glossary.html#dt-documentorder).
//
// This has the benefit over the visitor patten that it is iterative and can
// be used with Java 8 Stream API for searching, filtering, transforms, and
// other functional processing.
//
// The given predicate determines whether a node's children will be iterated
// over. If a node does not match the predicate, none of its children will be
// visited.
//
// root is the root of the tree; travserseNodePredicate matches nodes in the
// tree whose children should be traversed.
//
// Go-only: language-adapter — the Iterable is an iter.Seq; each traversal
// starts a new PreOrderIterator, like each iterator() call of the Java
// lambda.
//
// Java: NodeUtil.preOrderIterable(Node, Predicate<Node>)
func NodeUtil_PreOrderIterableWithPred(root rhino.Node, travserseNodePredicate func(rhino.Node) bool) iter.Seq[rhino.Node] {
	return func(yield func(rhino.Node) bool) {
		if root.IsNil() {
			panic(jutil.NewNullPointerException())
		}
		for cur := root; !cur.IsNil(); {
			next := nodeUtil_nextPreOrderNode(cur, travserseNodePredicate)
			if !yield(cur) {
				return
			}
			cur = next
		}
	}
}

// NodeUtil_PreOrderIterable is the same as NodeUtil_PreOrderIterableWithPred
// but iterates over all nodes in the tree without exception.
//
// Java: NodeUtil.preOrderIterable(Node)
func NodeUtil_PreOrderIterable(root rhino.Node) iter.Seq[rhino.Node] {
	return NodeUtil_PreOrderIterableWithPred(
		root,
		// Traverse all nodes.
		nodeUtil_alwaysTrue)
}

func nodeUtil_nextPreOrderNode(currentNode rhino.Node, traverseNodePredicate func(rhino.Node) bool) rhino.Node {
	// If node does not match the predicate, do not descend into it.
	if traverseNodePredicate(currentNode) && currentNode.HasChildren() {
		// In prefix order, the next node is the leftmost child.
		return currentNode.GetFirstChild()
	}

	// Walk up the ancestry chain (including current node) and return the first sibling we see.
	for !currentNode.IsNil() {
		if next := currentNode.GetNext(); !next.IsNil() {
			return next
		}
		currentNode = currentNode.GetParent()
	}
	return rhino.NilNode
}

// Java: NodeUtil.countAstSize(Node)
func NodeUtil_CountAstSize(n rhino.Node) int {
	count := 1
	if n.GetFirstChild().IsNil() {
		return 1
	}
	var stackBuf [64]rhino.Node
	stack := stackBuf[:0]
	cur := n
	for {
		for c := cur.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			count++
			if gc := c.GetFirstChild(); !gc.IsNil() {
				if gc.GetNext().IsNil() && gc.GetFirstChild().IsNil() {
					count++
				} else {
					stack = append(stack, c)
				}
			}
		}
		if len(stack) == 0 {
			return count
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}

// Java: NodeUtil.createConstantJsDoc()
func NodeUtil_CreateConstantJsDoc() *rhino.JSDocInfo {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordConstancy()
	return builder.Build()
}

var externVarNameVisitTokens = makeVisitTokenSet(rhino.Token_NAME, rhino.Token_IMPORT_STAR)

type externVarNameCollector struct {
	orderedSet  *jutil.LinkedHashSet[string]
	externNames map[string]struct{}
	reserved    stringMapSet
}

var (
	_ NodeTraversal_Callback         = (*externVarNameCollector)(nil)
	_ nodeTraversal_noShouldTraverse = (*externVarNameCollector)(nil)
	_ nodeTraversal_visitTokens      = (*externVarNameCollector)(nil)
)

func (*externVarNameCollector) nodeTraversalNoShouldTraverse() {}
func (*externVarNameCollector) visitTokens() *ast.TokenSet     { return &externVarNameVisitTokens }
func (*externVarNameCollector) ShouldTraverse(*NodeTraversal, rhino.Node, rhino.Node) bool {
	return true
}

func (c *externVarNameCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !parent.IsNil() {
		if pt := parent.GetToken(); (pt == rhino.Token_IMPORT_SPEC && n != parent.GetLastChild()) ||
			(pt == rhino.Token_EXPORT_SPEC && n != parent.GetFirstChild()) {
			return
		}
	}
	id := n.GetStringID()
	var v *Var
	if id == implicitArgumentsID {
		v = t.GetScope().GetVarByID(id)
	} else {
		v = t.LookupVarByID(id)
	}
	if v != nil && !v.IsParam() {
		name := v.GetName()
		if c.orderedSet != nil {
			c.orderedSet.Add(name)
		} else {
			c.externNames[name] = struct{}{}
			if c.reserved != nil {
				c.reserved[name] = struct{}{}
			}
		}
	}
}

func collectExternVariableNamesInto(
	compiler AbstractCompiler,
	externs rhino.Node,
	externNames map[string]struct{},
	reserved stringMapSet,
) {
	cb := externVarNameCollector{externNames: externNames, reserved: reserved}
	b := NodeTraversal_Builder{
		compiler: compiler,
		callback: &cb,
		obeyDestructuringAndDefaultValueExecutionOrder: true,
	}
	b.Traverse(externs)
}

// Java: NodeUtil.collectExternVariableNames(AbstractCompiler, Node)
func NodeUtil_CollectExternVariableNames(
	compiler AbstractCompiler, externs rhino.Node) *jutil.LinkedHashSet[string] {
	externsNames := jutil.NewLinkedHashSet[string]()
	cb := externVarNameCollector{orderedSet: externsNames}
	b := NodeTraversal_Builder{
		compiler: compiler,
		callback: &cb,
		obeyDestructuringAndDefaultValueExecutionOrder: true,
	}
	b.Traverse(externs)
	return externsNames
}

// NodeUtil_GetAllVarNamesDeclaredInModule records a mapping of names to vars
// of everything reachable in a module. Should only be called with a module
// scope.
//
// nameVarMap: an empty map that gets populated with the keys being variable
// names and values being variable objects. orderedVars: an empty list that
// gets populated with variable objects in the order that they appear in the
// module.
//
// Java: NodeUtil.getAllVarNamesDeclaredInModule(Node, AbstractCompiler, ScopeCreator, Scope)
func NodeUtil_GetAllVarNamesDeclaredInModule(
	moduleNode rhino.Node,
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
	globalScope *Scope) *jutil.LinkedHashSet[string] {

	jutil.CheckStateMsg(moduleNode.IsModuleBody(), "getAllVarsDeclaredInModule expects a module body node")
	jutil.CheckStateMsg1(globalScope.IsGlobal(), "%s", globalScope)

	nameVars := jutil.NewLinkedHashSet[string]()

	finder := &nodeUtil_getAllVarNamesDeclaredInModule_finder{nameVars: nameVars}
	NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(finder).
		SetScopeCreator(scopeCreator).
		TraverseWithScope(moduleNode, globalScope)
	return nameVars
}

// nodeUtil_getAllVarNamesDeclaredInModule_finder is the anonymous
// ScopedCallback of getAllVarNamesDeclaredInModule.
type nodeUtil_getAllVarNamesDeclaredInModule_finder struct {
	nameVars *jutil.LinkedHashSet[string]
}

var _ nodeTraversal_noVisit = (*nodeUtil_getAllVarNamesDeclaredInModule_finder)(nil)

func (*nodeUtil_getAllVarNamesDeclaredInModule_finder) nodeTraversalNoVisit() {}

// Java: NodeUtil.getAllVarNamesDeclaredInModule(...).finder.enterScope(NodeTraversal)
func (f *nodeUtil_getAllVarNamesDeclaredInModule_finder) EnterScope(t *NodeTraversal) {
	currentScope := t.GetScope()
	if currentScope.IsModuleScope() {
		for v := range currentScope.GetVarIterable() {
			f.nameVars.Add(v.GetName())
		}
	}
}

// Java: NodeUtil.getAllVarNamesDeclaredInModule(...).finder.exitScope(NodeTraversal)
func (f *nodeUtil_getAllVarNamesDeclaredInModule_finder) ExitScope(t *NodeTraversal) {}

// Java: NodeUtil.getAllVarNamesDeclaredInModule(...).finder.shouldTraverse(NodeTraversal, Node, Node)
func (f *nodeUtil_getAllVarNamesDeclaredInModule_finder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return n.IsModuleBody()
}

// Java: NodeUtil.getAllVarNamesDeclaredInModule(...).finder.visit(NodeTraversal, Node, Node)
func (f *nodeUtil_getAllVarNamesDeclaredInModule_finder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

// NodeUtil_GetAllVarsDeclaredInFunction returns a mapping of names to vars
// of everything reachable in a function. Should only be called with a
// function scope. Does not enter new control flow areas aka embedded
// functions.
//
// Java: NodeUtil.getAllVarsDeclaredInFunction(AbstractCompiler, ScopeCreator, Scope)
func NodeUtil_GetAllVarsDeclaredInFunction(
	compiler AbstractCompiler, scopeCreator ScopeCreator, scope *Scope) *NodeUtil_AllVarsDeclaredInFunction {
	jutil.CheckStateMsg1(scope.IsFunctionScope(), "%s", scope)

	finder := &nodeUtil_getAllVarsDeclaredInFunction_finder{
		scope: scope,
	}

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     finder,
		scopeCreator: scopeCreator,
	}
	b.TraverseAtScope(scope)

	finder.result.table.finalize()
	jutil.CheckState((finder.result.table.uniqueCount == 0) == (len(finder.result.table.vars) == 0))
	return &finder.result
}

// nodeUtil_getAllVarsDeclaredInFunction_finder is the anonymous
// ScopedCallback of getAllVarsDeclaredInFunction.
type nodeUtil_getAllVarsDeclaredInFunction_finder struct {
	result NodeUtil_AllVarsDeclaredInFunction
	scope  *Scope
}

func nodeUtilSubtreeHasBlockScopeDecls(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_CATCH, rhino.Token_IMPORT:
		return true
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(n)
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.GetFirstChild().IsNil() && nodeUtilSubtreeHasBlockScopeDecls(c) {
			return true
		}
	}
	return false
}

func nodeUtilStmtHasBlockScopeDecls(stmt rhino.Node) bool {
	switch stmt.GetToken() {
	case rhino.Token_FUNCTION:
		return false
	case rhino.Token_CLASS, rhino.Token_CATCH, rhino.Token_IMPORT:
		return true
	default:
		for c := stmt.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if !c.GetFirstChild().IsNil() && nodeUtilSubtreeHasBlockScopeDecls(c) {
				return true
			}
		}
		return false
	}
}

func nodeUtilBlockScopeMayHaveDecls(n, root rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH, rhino.Token_IMPORT:
		return true
	case rhino.Token_EXPORT:
		first := n.GetFirstChild()
		return !first.IsNil() && nodeUtilBlockScopeMayHaveDecls(first, root)
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(n) && n.GetFirstChild().GetStringID() > rhino.RhinoStringPool_EmptyID
	case rhino.Token_CLASS:
		return !NodeUtil_IsClassExpression(n) && !n.GetFirstChild().IsEmpty()
	case rhino.Token_BLOCK:
		if n != root && scanVarsBlockCreatesScope(n) {
			return false
		}
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY:
		if n != root {
			return false
		}
	case rhino.Token_ROOT, rhino.Token_DO, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_IF,
		rhino.Token_LABEL, rhino.Token_TRY, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
	default:
		return false
	}
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if scanVarsTokenKind[uint8(child.GetToken())] == 1 {
			if nodeUtilBlockScopeMayHaveDecls(child, root) {
				return true
			}
		}
	}
	return false
}

func nodeUtilScopeRootMayDeclareVars(root, fnRoot rhino.Node) bool {
	switch root.GetToken() {
	case rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_MEMBER_FIELD_DEF:
		return false
	case rhino.Token_CLASS:
		return !root.GetFirstChild().IsEmpty() && NodeUtil_IsClassExpression(root)
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		return syntacticScopeHasBlockDecls(root)
	case rhino.Token_SWITCH_BODY:
		return nodeUtilBlockScopeMayHaveDecls(root, root)
	case rhino.Token_BLOCK:
		parent := root.GetParent()
		if parent == fnRoot || (!parent.IsNil() && parent.IsClassMembers()) {
			return !syntacticFunctionBlockHasNoDecls(root)
		}
		return syntacticScopeHasBlockDecls(root)
	default:
		return true
	}
}

func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) appendScopeVars(currentScope *Scope) {
	f.result.table.appendScopeVars(currentScope)
}

// Java: NodeUtil.getAllVarsDeclaredInFunction(...).finder.enterScope(NodeTraversal)
func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) EnterScope(t *NodeTraversal) {
	root := t.GetScopeRoot()
	var currentScope *Scope
	if root == f.scope.GetRootNode() {
		currentScope = f.scope
	} else {
		if !nodeUtilScopeRootMayDeclareVars(root, f.scope.GetRootNode()) {
			return
		}
		currentScope = t.GetScope()
	}
	f.appendScopeVars(currentScope)
}

func (*nodeUtil_getAllVarsDeclaredInFunction_finder) nodeTraversalNoVisit() {}

// Java: NodeUtil.getAllVarsDeclaredInFunction(...).finder.exitScope(NodeTraversal)
func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) ExitScope(t *NodeTraversal) {}

// Java: NodeUtil.getAllVarsDeclaredInFunction(...).finder.shouldTraverse(NodeTraversal, Node, Node)
func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	fnRoot := f.scope.GetRootNode()
	if n == fnRoot {
		return true
	}
	if n.IsFunction() {
		// Don't enter any new functions
		return false
	}
	if parent == fnRoot {
		if n.IsBlock() {
			return n.HasChildren()
		}
		if n.IsParamList() {
			for p := n.GetFirstChild(); !p.IsNil(); p = p.GetNext() {
				if !p.IsName() && nodeUtilSubtreeHasBlockScopeDecls(p) {
					return true
				}
			}
			return false
		}
		if n == fnRoot.GetFirstChild() {
			return false
		}
		return nodeUtilSubtreeHasBlockScopeDecls(n)
	}
	if !parent.IsNil() && parent.GetParent() == fnRoot && parent.IsBlock() {
		return nodeUtilStmtHasBlockScopeDecls(n)
	}
	return true
}

// Java: NodeUtil.getAllVarsDeclaredInFunction(...).finder.visit(NodeTraversal, Node, Node)
func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

// NodeUtil_AllVarsDeclaredInFunction represents a mapping of names to vars
// of everything reachable in a function.
//
// Java: NodeUtil.AllVarsDeclaredInFunction
type NodeUtil_AllVarsDeclaredInFunction struct {
	table       fnVarTable
	allVarsInFn *jutil.LinkedHashMap[string, *Var]
}

// Java: NodeUtil.AllVarsDeclaredInFunction.getAllVariables()
func (a *NodeUtil_AllVarsDeclaredInFunction) GetAllVariables() *jutil.LinkedHashMap[string, *Var] {
	if a.allVarsInFn == nil {
		m := jutil.NewLinkedHashMapWithCapacity[string, *Var](len(a.table.vars))
		for _, v := range a.table.vars {
			m.Put(v.GetName(), v)
		}
		a.allVarsInFn = m
	}
	return a.allVarsInFn
}

// Java: NodeUtil.AllVarsDeclaredInFunction.getAllVariablesInOrder()
func (a *NodeUtil_AllVarsDeclaredInFunction) GetAllVariablesInOrder() []*Var {
	return a.table.vars
}
