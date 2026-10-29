// Port of com/google/javascript/jscomp/NodeTraversal.java @ closure-compiler v20261021

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
	"strconv"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeTraversal allows an iteration through the nodes in the parse tree, and
// facilitates the optimizations on the parse tree.
type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCallback                                  NodeTraversal_ScopedCallback
	scopeCreator                                   ScopeCreator
	syntacticScopeCreator                          *SyntacticScopeCreator
	lazyGlobalScope                                *lazyGlobalScope
	defaultScopeCreator                            SyntacticScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
	mayContainSyntheticBlocks                      bool
	noShouldTraverse                               bool
	noBranchShouldTraverse                         bool
	noVisit                                        bool
	noVisitLeaves                                  bool
	onlyNameLeaves                                 bool
	// visitTokens, if not nil, holds the only tokens Visit acts on (see
	// nodeTraversal_visitTokens): the traversal skips the call for others.
	visitTokens *ast.TokenSet

	// Contains the current node
	currentNode rhino.Node

	currentHoistScopeRoot rhino.Node

	// Contains the current FUNCTION node if there is one, otherwise null.
	currentFunction rhino.Node

	// Contains the enclosing SCRIPT node if there is one, otherwise null.
	currentScript rhino.Node

	// The change scope for the current node being visiteds
	currentChangeScope rhino.Node

	// The chain scope for the currentNode being visited. Scopes are relatively
	// expensive to build so they are built lazily. Each entry holds the root
	// Node of the scope and, once instantiated, its AbstractScope.
	//
	// If `scope` is non-nil all previous entries will also have been
	// instantiated due to the hierarchical nature of scopes.
	//
	// Go-only: performance — concrete struct instead of any avoids interface
	// boxing and runtime type assertions on every GetScopeRoot / GetScope call.
	scopes    []scopeStackEntry
	scopesBuf [8]scopeStackEntry

	// The current source file name
	sourceName *string

	// The current input
	inputId *rhino.InputId

	compilerInput *CompilerInput
}

var nodeTraversalPool = sync.Pool{
	New: func() any { return &NodeTraversal{} },
}

func acquireNodeTraversal(builder *NodeTraversal_Builder) *NodeTraversal {
	if builder.compiler == nil || builder.callback == nil {
		failNodeTraversalNull()
	}
	t := nodeTraversalPool.Get().(*NodeTraversal)
	if cap(t.scopes) == 0 {
		t.scopes = t.scopesBuf[:0]
	}
	initNodeTraversalFromBuilder(t, builder)
	return t
}

func releaseNodeTraversal(t *NodeTraversal) {
	scopes := t.scopes
	clear(scopes)
	*t = NodeTraversal{}
	if cap(scopes) > len(t.scopesBuf) {
		t.scopes = scopes[:0]
	} else {
		t.scopes = t.scopesBuf[:0]
	}
	nodeTraversalPool.Put(t)
}

// Java: NodeTraversal(Builder)
func newNodeTraversal(builder *NodeTraversal_Builder) *NodeTraversal {
	if builder.compiler == nil || builder.callback == nil {
		failNodeTraversalNull()
	}
	t := &NodeTraversal{}
	t.scopes = t.scopesBuf[:0]
	initNodeTraversalFromBuilder(t, builder)
	return t
}

func initNodeTraversalFromBuilder(t *NodeTraversal, builder *NodeTraversal_Builder) {
	t.compiler = builder.compiler
	cb := builder.callback
	t.callback = cb
	_, t.noShouldTraverse = cb.(nodeTraversal_noShouldTraverse)
	if !t.noShouldTraverse {
		_, t.noShouldTraverse = cb.(NodeTraversal_AbstractPostOrderCallbackInterface)
	}
	t.noBranchShouldTraverse = t.noShouldTraverse
	if !t.noBranchShouldTraverse {
		if sf, ok := cb.(nodeTraversal_onlyScriptAndFunctionShouldTraverse); ok {
			t.noBranchShouldTraverse = sf.shouldTraverseOnlyScriptAndFunction()
		}
	}
	_, t.noVisit = cb.(nodeTraversal_noVisit)
	if vt, ok := cb.(nodeTraversal_visitTokens); ok {
		t.visitTokens = vt.visitTokens()
	} else {
		t.visitTokens = nil
	}
	_, t.onlyNameLeaves = cb.(nodeTraversal_onlyNameLeaves)
	if !t.onlyNameLeaves {
		_, t.noVisitLeaves = cb.(nodeTraversal_noVisitLeaves)
	} else {
		t.noVisitLeaves = false
	}
	if scopedCallback, ok := cb.(NodeTraversal_ScopedCallback); ok {
		t.scopeCallback = scopedCallback
	} else {
		t.scopeCallback = nil
	}
	t.scopeCreator = builder.scopeCreator
	if builder.scopeCreator == nil {
		t.defaultScopeCreator = SyntacticScopeCreator{
			compiler:             builder.compiler,
			redeclarationHandler: SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER,
		}
		t.syntacticScopeCreator = &t.defaultScopeCreator
	} else if ssc, ok := builder.scopeCreator.(*SyntacticScopeCreator); ok {
		t.syntacticScopeCreator = ssc
	} else {
		t.syntacticScopeCreator = nil
	}
	t.obeyDestructuringAndDefaultValueExecutionOrder = builder.obeyDestructuringAndDefaultValueExecutionOrder
	// TODO(johnlenz): refactor into compiler options helper.
	options := t.compiler.GetOptions()
	t.mayContainSyntheticBlocks =
		options == nil || options.GetSyntheticBlockStartMarker() != nil
}

// Java: NodeTraversal.builder()
func NodeTraversal_Builder_() *NodeTraversal_Builder {
	return newNodeTraversal_Builder()
}

// NodeTraversal_Builder is the builder.
type NodeTraversal_Builder struct {
	callback                                       NodeTraversal_Callback
	compiler                                       AbstractCompiler
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

// Java: NodeTraversal.Builder()
func newNodeTraversal_Builder() *NodeTraversal_Builder {
	return &NodeTraversal_Builder{obeyDestructuringAndDefaultValueExecutionOrder: false}
}

// Java: NodeTraversal.Builder.setCallback(Callback)
func (b *NodeTraversal_Builder) SetCallback(x NodeTraversal_Callback) *NodeTraversal_Builder {
	b.callback = x
	return b
}

// SetCallbackFunc sets a post-order callback given as a function (a Java
// lambda for AbstractPostOrderCallbackInterface).
//
// Java: NodeTraversal.Builder.setCallback(AbstractPostOrderCallbackInterface)
func (b *NodeTraversal_Builder) SetCallbackFunc(x NodeTraversal_AbstractPostOrderCallbackInterface) *NodeTraversal_Builder {
	if x != nil {
		b.callback = x
	} else {
		b.callback = nodeTraversal_nilPostOrderCallback
	}
	return b
}

// Java: NodeTraversal.Builder.setCompiler(AbstractCompiler)
func (b *NodeTraversal_Builder) SetCompiler(x AbstractCompiler) *NodeTraversal_Builder {
	b.compiler = x
	return b
}

// Java: NodeTraversal.Builder.setScopeCreator(ScopeCreator)
func (b *NodeTraversal_Builder) SetScopeCreator(x ScopeCreator) *NodeTraversal_Builder {
	b.scopeCreator = x
	return b
}

// Java: NodeTraversal.Builder.setObeyDestructuringAndDefaultValueExecutionOrder(boolean)
func (b *NodeTraversal_Builder) SetObeyDestructuringAndDefaultValueExecutionOrder(x bool) *NodeTraversal_Builder {
	b.obeyDestructuringAndDefaultValueExecutionOrder = x
	return b
}

// Java: NodeTraversal.Builder.build()
func (b *NodeTraversal_Builder) Build() *NodeTraversal {
	return newNodeTraversal(b)
}

// Java: NodeTraversal.Builder.traverse(Node)
func (b *NodeTraversal_Builder) Traverse(root rhino.Node) {
	t := acquireNodeTraversal(b)
	t.Traverse(root)
	releaseNodeTraversal(t)
}

// Java: NodeTraversal.Builder.traverseAtScope(AbstractScope<?, ?>)
func (b *NodeTraversal_Builder) TraverseAtScope(scope AbstractScope) {
	t := acquireNodeTraversal(b)
	t.TraverseAtScope(scope)
	releaseNodeTraversal(t)
}

// Java: NodeTraversal.Builder.traverseRoots(Node, Node)
func (b *NodeTraversal_Builder) TraverseRoots(externs rhino.Node, root rhino.Node) {
	t := acquireNodeTraversal(b)
	t.TraverseRoots(externs, root)
	releaseNodeTraversal(t)
}

// Java: NodeTraversal.Builder.traverseWithScope(Node, AbstractScope<?, ?>)
func (b *NodeTraversal_Builder) TraverseWithScope(root rhino.Node, s AbstractScope) {
	t := acquireNodeTraversal(b)
	t.TraverseWithScope(root, s)
	releaseNodeTraversal(t)
}

func (b *NodeTraversal_Builder) TraverseScopeRoot(scopeNode rhino.Node) {
	t := acquireNodeTraversal(b)
	t.TraverseScopeRoot(scopeNode)
	releaseNodeTraversal(t)
}

func (b *NodeTraversal_Builder) traverseScriptUnderRoot(
	scopeRoot rhino.Node,
	scriptParent rhino.Node,
	script rhino.Node,
	lazyGlobal *lazyGlobalScope,
) {
	t := acquireNodeTraversal(b)
	traverseScriptUnderRoot(t, scopeRoot, scriptParent, script, lazyGlobal)
	releaseNodeTraversal(t)
}

// Java: NodeTraversal.throwUnexpectedException(Throwable)
func (t *NodeTraversal) ThrowUnexpectedException(unexpectedException jutil.Throwable) {
	// If there's an unexpected exception, try to get the
	// line number of the code that caused it.
	message := jutil.GetMessageOrNull(unexpectedException)

	// TODO(user): It is possible to get more information if currentNode or
	// its parent is missing. We still have the scope stack in which it is still
	// very useful to find out at least which function caused the exception.
	if !t.currentScript.IsNil() {
		parentContext := ""
		if !t.currentNode.IsNil() {
			parentContext = t.FormatNodeContext("Parent", t.currentNode.GetParent())
		}
		m := jutil.ToJavaString(jutil.GetMessageOrNull(unexpectedException)) +
			"\n" +
			t.FormatNodeContext("Node", t.currentNode) +
			parentContext
		message = &m
	}
	// Go-only: language-adapter — AbstractCompiler.throwInternalError takes a
	// Go string. A null message (an exception without one, outside any script)
	// is passed as "null": Compiler.throwInternalError only concatenates it,
	// which renders a null String as "null".
	t.compiler.ThrowInternalError(jutil.ToJavaString(message), unexpectedException)
}

// Java: NodeTraversal.formatNodeContext(String, Node)
func (t *NodeTraversal) FormatNodeContext(label string, n rhino.Node) string {
	if n.IsNil() {
		return "  " + label + ": NULL"
	}
	return "  " + label + "(" + n.ToStringWithFlags(false, false, false) + "): " + t.FormatNodePosition(n)
}

// Traverse traverses a parse tree recursively.
//
// Java: NodeTraversal.traverse(Node)
func (t *NodeTraversal) Traverse(root rhino.Node) {
	defer func() {
		if r := recover(); r != nil {
			t.ThrowUnexpectedException(jutil.AsThrowable(r))
		}
	}()
	t.InitTraversal(root)
	t.currentNode = root
	t.PushScope(root)
	// null parent ensures that the shallow callbacks will traverse root
	t.TraverseBranch(root, rhino.NilNode)
	t.PopScope()
}

// NodeTraversal_Traverse traverses using the SyntacticScopeCreator
//
// Java: NodeTraversal.traverse(AbstractCompiler, Node, Callback)
func NodeTraversal_Traverse(compiler AbstractCompiler, root rhino.Node, cb NodeTraversal_Callback) {
	b := NodeTraversal_Builder{compiler: compiler, callback: cb}
	b.Traverse(root)
}

// Java: NodeTraversal.traverseRoots(Node, Node)
func (t *NodeTraversal) TraverseRoots(externs rhino.Node, root rhino.Node) {
	defer func() {
		if r := recover(); r != nil {
			t.ThrowUnexpectedException(jutil.AsThrowable(r))
		}
	}()
	scopeRoot := externs.GetParent()
	if scopeRoot.IsNil() {
		failNodeTraversalNull()
	}

	t.InitTraversal(scopeRoot)
	t.currentNode = scopeRoot
	t.PushScope(scopeRoot)

	t.TraverseBranch(externs, scopeRoot)
	jutil.CheckState(root.GetParent() == scopeRoot)
	t.TraverseBranch(root, scopeRoot)

	t.PopScope()
}

// Java: NodeTraversal.traverseRoots(AbstractCompiler, Callback, Node, Node)
func NodeTraversal_TraverseRoots(
	compiler AbstractCompiler, cb NodeTraversal_Callback, externs rhino.Node, root rhino.Node) {
	b := NodeTraversal_Builder{compiler: compiler, callback: cb}
	b.TraverseRoots(externs, root)
}

const nodeTraversal_MISSING_SOURCE = "[source unknown]"

// Java: NodeTraversal.formatNodePosition(Node)
func (t *NodeTraversal) FormatNodePosition(n rhino.Node) string {
	sourceFileName := t.GetBestSourceFileName(n)
	if sourceFileName == nil {
		return nodeTraversal_MISSING_SOURCE + "\n"
	}

	lineNumber := n.GetLineno()
	columnNumber := n.GetCharno()
	src := t.compiler.GetSourceLine(sourceFileName, lineNumber)
	if src == nil {
		missing := nodeTraversal_MISSING_SOURCE
		src = &missing
	}
	return *sourceFileName + ":" + strconv.Itoa(lineNumber) + ":" + strconv.Itoa(columnNumber) + "\n" + *src + "\n"
}

// TraverseWithScope traverses a parse tree recursively with a scope,
// starting with the given root. This should only be used in the global
// scope or module scopes. Otherwise, use TraverseAtScope.
//
// Java: NodeTraversal.traverseWithScope(Node, AbstractScope<?, ?>)
func (t *NodeTraversal) TraverseWithScope(root rhino.Node, s AbstractScope) {
	jutil.CheckStateMsg1(s.IsGlobal() || s.IsModuleScope(), "%s", s)
	defer func() {
		if r := recover(); r != nil {
			t.ThrowUnexpectedException(jutil.AsThrowable(r))
		}
	}()
	t.InitTraversal(root)
	t.currentNode = root
	t.PushScopeOfScope(s)
	t.TraverseBranch(root, rhino.NilNode)
	t.PopScope()
}

// Java: NodeTraversal.getCompiler()
func (t *NodeTraversal) GetCompiler() AbstractCompiler {
	return t.compiler
}

var nodeTraversal_emptyString = ""

// GetSourceName gets the current input source name: a string that may be
// empty, but not null.
//
// Go-only: language-adapter — a *string, since the Java String comes from
// Node.getSourceFileName(), which may be null.
//
// Java: NodeTraversal.getSourceName()
func (t *NodeTraversal) GetSourceName() *string {
	if t.sourceName == nil {
		if !t.currentScript.IsNil() {
			t.sourceName = t.currentScript.GetSourceFileName()
		} else {
			t.sourceName = &nodeTraversal_emptyString
		}
	}
	return t.sourceName
}

// GetInput gets the current input source.
//
// Java: NodeTraversal.getInput()
func (t *NodeTraversal) GetInput() *CompilerInput {
	inputId := t.GetInputId()
	if t.compilerInput == nil && inputId != nil {
		t.compilerInput = t.compiler.GetInput(inputId)
	}
	return t.compilerInput
}

// GetChunk gets the current input chunk.
//
// Java: NodeTraversal.getChunk()
func (t *NodeTraversal) GetChunk() *JSChunk {
	input := t.GetInput()
	if input == nil {
		return nil
	}
	return input.GetChunk()
}

// GetCurrentNode returns the node currently being traversed.
//
// Java: NodeTraversal.getCurrentNode()
func (t *NodeTraversal) GetCurrentNode() rhino.Node {
	return t.currentNode
}

func (t *NodeTraversal) shouldTraverseScriptOrFunction(n, parent rhino.Node) bool {
	return t.noShouldTraverse || t.callback.ShouldTraverse(t, n, parent)
}

func (t *NodeTraversal) shouldTraverse(n, parent rhino.Node) bool {
	return t.noBranchShouldTraverse || t.callback.ShouldTraverse(t, n, parent)
}

func (t *NodeTraversal) visitNode(n, parent rhino.Node) {
	if !t.noVisit && (t.visitTokens == nil || t.visitTokens.Has(n.GetToken())) {
		t.callback.Visit(t, n, parent)
	}
}

// Java: NodeTraversal.handleScript(Node, Node)
func (t *NodeTraversal) HandleScript(n rhino.Node, parent rhino.Node) {
	if Platform_IsThreadInterrupted() {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
	}
	t.SetChangeScope(n)

	t.currentNode = n
	t.currentScript = n
	t.ClearScriptState()
	if t.shouldTraverseScriptOrFunction(n, parent) {
		t.TraverseChildren(n)
		t.currentNode = n
		t.visitNode(n, parent)
	}
	t.SetChangeScope(rhino.NilNode)
}

// Java: NodeTraversal.handleFunction(Node, Node)
func (t *NodeTraversal) HandleFunction(n rhino.Node, parent rhino.Node) {
	changeScope := t.currentChangeScope
	t.SetChangeScope(n)
	t.currentNode = n
	if t.shouldTraverseScriptOrFunction(n, parent) {
		t.TraverseFunction(n, parent)
		t.currentNode = n
		t.visitNode(n, parent)
	}
	t.SetChangeScope(changeScope)
}

// HandleModule traverses a module.
//
// Java: NodeTraversal.handleModule(Node, Node)
func (t *NodeTraversal) HandleModule(n rhino.Node, parent rhino.Node) {
	t.currentHoistScopeRoot = n
	t.PushScope(n)
	t.currentNode = n
	if t.shouldTraverse(n, parent) {
		t.TraverseChildren(n)
		t.currentNode = n
		t.visitNode(n, parent)
	}
	t.PopScope()
	// Module bodies don't nest
	t.currentHoistScopeRoot = rhino.NilNode
}

// Java: NodeTraversal.handleDestructuringOrDefaultValue(Node, Node)
func (t *NodeTraversal) HandleDestructuringOrDefaultValue(n rhino.Node, parent rhino.Node) {
	t.currentNode = n
	if t.shouldTraverse(n, parent) {
		first := n.GetFirstChild()
		second := first.GetNext()

		if !second.IsNil() {
			jutil.CheckStateMsg1(second.GetNext().IsNil(), "%s", second)
			t.TraverseBranch(second, n)
		}
		t.TraverseBranch(first, n)

		t.currentNode = n
		t.visitNode(n, parent)
	}
}

const (
	nodeTraversalTokenNormal uint8 = iota
	nodeTraversalTokenBlock
	nodeTraversalTokenBlockScope
	nodeTraversalTokenSpecial
)

var nodeTraversalTokenKind = [256]uint8{
	rhino.Token_BLOCK:             nodeTraversalTokenBlock,
	rhino.Token_FOR:               nodeTraversalTokenBlockScope,
	rhino.Token_FOR_IN:            nodeTraversalTokenBlockScope,
	rhino.Token_FOR_OF:            nodeTraversalTokenBlockScope,
	rhino.Token_FOR_AWAIT_OF:      nodeTraversalTokenBlockScope,
	rhino.Token_SWITCH_BODY:       nodeTraversalTokenBlockScope,
	rhino.Token_SCRIPT:            nodeTraversalTokenSpecial,
	rhino.Token_FUNCTION:          nodeTraversalTokenSpecial,
	rhino.Token_MODULE_BODY:       nodeTraversalTokenSpecial,
	rhino.Token_CLASS:             nodeTraversalTokenSpecial,
	rhino.Token_CLASS_MEMBERS:     nodeTraversalTokenSpecial,
	rhino.Token_DEFAULT_VALUE:     nodeTraversalTokenSpecial,
	rhino.Token_DESTRUCTURING_LHS: nodeTraversalTokenSpecial,
}

// TraverseBranch traverses a branch.
//
// Java: NodeTraversal.traverseBranch(Node, Node)
func (t *NodeTraversal) TraverseBranch(n rhino.Node, parent rhino.Node) {
	tok := n.GetToken()
	pushedScope := false
	switch nodeTraversalTokenKind[uint8(tok)] {
	case nodeTraversalTokenSpecial:
		switch tok {
		case rhino.Token_SCRIPT:
			t.HandleScript(n, parent)
			return
		case rhino.Token_FUNCTION:
			t.HandleFunction(n, parent)
			return
		case rhino.Token_MODULE_BODY:
			t.HandleModule(n, parent)
			return
		case rhino.Token_CLASS:
			t.HandleClass(n, parent)
			return
		case rhino.Token_CLASS_MEMBERS:
			t.HandleClassMembers(n, parent)
			return
		case rhino.Token_DEFAULT_VALUE, rhino.Token_DESTRUCTURING_LHS:
			// TODO(nickreid): Handle ASSIGN with destructuring target.
			if t.obeyDestructuringAndDefaultValueExecutionOrder {
				t.HandleDestructuringOrDefaultValue(n, parent)
				return
			}
		}
	case nodeTraversalTokenBlock:
		pushedScope = t.isBlockScopeBlock(n)
	case nodeTraversalTokenBlockScope:
		pushedScope = true
	}

	cb := t.callback
	t.currentNode = n
	if !t.noBranchShouldTraverse && !cb.ShouldTraverse(t, n, parent) {
		return
	}

	var previousHoistScopeRoot rhino.Node
	if pushedScope {
		previousHoistScopeRoot = t.currentHoistScopeRoot
		t.PushScope(n)
		if tok == rhino.Token_BLOCK && n.GetParent().IsClassMembers() {
			t.currentHoistScopeRoot = n
		}
	}

	/*
	 * Intentionally inlined call to traverseChildren.
	 *
	 * <p>Calling traverseChildren here would double the maximum stack depth seen during
	 * compilation. That can cause stack overflows on some very deep ASTs (e.g. b/188616350).
	 * Inlining side-steps the issue.
	 */
	noVisitLeaves := t.noVisitLeaves
	onlyNameLeaves := t.onlyNameLeaves
	for child := n.GetFirstChild(); !child.IsNil(); {
		// child could be replaced, in which case our child node
		// would no longer point to the true next
		next := child.GetNext()
		if child.GetFirstChild().IsNil() {
			ctok := child.GetToken()
			if nodeTraversalTokenKind[uint8(ctok)] == 0 {
				if !t.noBranchShouldTraverse {
					t.currentNode = child
					if !cb.ShouldTraverse(t, child, n) {
						child = next
						continue
					}
				}
				if !noVisitLeaves && (!onlyNameLeaves || ctok == rhino.Token_NAME) && !t.noVisit &&
					(t.visitTokens == nil || t.visitTokens.Has(ctok)) {
					t.currentNode = child
					cb.Visit(t, child, n)
				}
				child = next
				continue
			}
		}
		t.TraverseBranch(child, n)
		child = next
	}

	if pushedScope {
		t.PopScope()
		t.currentHoistScopeRoot = previousHoistScopeRoot
	}

	if !t.noVisit {
		if (pushedScope || !n.GetFirstChild().IsNil() || (!noVisitLeaves && (!onlyNameLeaves || tok == rhino.Token_NAME)) || tok == rhino.Token_BLOCK) &&
			(t.visitTokens == nil || t.visitTokens.Has(tok)) {
			t.currentNode = n
			cb.Visit(t, n, parent)
		}
	}
}

// TraverseFunction traverses a function.
//
// Java: NodeTraversal.traverseFunction(Node, Node)
func (t *NodeTraversal) TraverseFunction(n rhino.Node, parent rhino.Node) {
	fnName := n.GetFirstChild()
	// NOTE: If a function declaration is the root of a traversal, then we will treat it as a
	// function expression (since 'parent' is null, even though 'n' actually has a parent node) and
	// traverse the function name before entering the scope, rather than afterwards. Removing the
	// null check for 'parent' seems safe, but causes a rare crash when traverseScopeRoots is called
	// by PeepholeOptimizationsPass on a function that somehow doesn't actually have a parent at all
	// (presumably because it's already been removed from the AST?) so that doesn't actually work.
	// This does not actually change anything, though, since rooting a traversal at a function node
	// causes the function scope to be entered twice (unless using #traverseAtScope, which doesn't
	// call this method), so that the name node is just always traversed inside the scope anyway.
	isFunctionDeclaration := !parent.IsNil() && NodeUtil_IsFunctionDeclaration(n)

	if isFunctionDeclaration {
		// Function declarations are in the scope containing the declaration.
		t.TraverseBranch(fnName, n)
	}

	t.currentNode = n

	// may nest within class static blocks
	previousHoistScopeRoot := t.currentHoistScopeRoot
	t.currentHoistScopeRoot = n
	previousFunction := t.currentFunction
	t.currentFunction = n
	t.PushScope(n)

	if !isFunctionDeclaration {
		// Function expression names are only accessible within the function
		// scope.
		t.TraverseBranch(fnName, n)
	}

	args := fnName.GetNext()
	body := args.GetNext()

	// Args
	t.TraverseBranch(args, n)

	// Body
	// ES6 "arrow" function may not have a block as a body.
	t.TraverseBranch(body, n)

	t.PopScope()
	t.currentFunction = previousFunction
	t.currentHoistScopeRoot = previousHoistScopeRoot
}

// HandleClass traverses a class. Note that we traverse some of the child
// nodes slightly out of order to ensure children are visited in the correct
// scope. The following children are in the outer scope: (1) the 'extends'
// clause, (2) any computed method keys, (3) the class name for class
// declarations only (class expression names are traversed in the class
// scope). This requires that we visit the extends node (second child) and
// any computed member keys (grandchildren of the last, body, child) before
// visiting the name (first child) or body (last child).
//
// Java: NodeTraversal.handleClass(Node, Node)
func (t *NodeTraversal) HandleClass(n rhino.Node, parent rhino.Node) {
	t.currentNode = n
	if !t.shouldTraverse(n, parent) {
		return
	}

	className := n.GetFirstChild()
	extendsClause := className.GetNext()
	body := extendsClause.GetNext()

	isClassExpression := NodeUtil_IsClassExpression(n)

	t.TraverseBranch(extendsClause, n)

	if !isClassExpression {
		// Class declarations are in the scope containing the declaration.
		t.TraverseBranch(className, n)
	}

	t.currentNode = n
	t.PushScope(n)

	if isClassExpression {
		// Class expression names are only accessible within the function
		// scope.
		t.TraverseBranch(className, n)
	}

	// Body
	t.TraverseBranch(body, n)

	t.PopScope()

	t.currentNode = n
	t.visitNode(n, parent)
}

// HandleClassMembers traverses class members, excluding keys of computed
// props.
//
// Java: NodeTraversal.handleClassMembers(Node, Node)
func (t *NodeTraversal) HandleClassMembers(n rhino.Node, parent rhino.Node) {
	t.currentNode = n
	if !t.shouldTraverse(n, parent) {
		return
	}

	for child := n.GetFirstChild(); !child.IsNil(); {
		next := child.GetNext() // see traverseChildren
		if child.IsComputedProp() || child.IsComputedFieldDef() {
			t.TraverseBranch(child.GetFirstChild(), child)
		}
		child = next
	}

	for child := n.GetFirstChild(); !child.IsNil(); {
		next := child.GetNext() // see traverseChildren

		switch child.GetToken() {
		case rhino.Token_COMPUTED_PROP:
			t.currentNode = n

			if t.shouldTraverse(child, n) {
				t.TraverseBranch(child.GetLastChild(), child)
				t.currentNode = n
				t.visitNode(child, n)
			}
		case rhino.Token_COMPUTED_FIELD_DEF:
			t.currentNode = n

			previousHoistScopeRoot := t.currentHoistScopeRoot
			t.currentHoistScopeRoot = n
			t.PushScope(child)

			if t.shouldTraverse(child, n) {
				if child.HasTwoChildren() { // No RHS to traverse in `[x];` computed field case
					t.TraverseBranch(child.GetLastChild(), child)
				}
				t.currentNode = n
				t.visitNode(child, n)
			}

			t.PopScope()
			t.currentHoistScopeRoot = previousHoistScopeRoot
		case rhino.Token_MEMBER_FIELD_DEF:
			t.HandleMemberFieldDef(n, child)
		case rhino.Token_BLOCK, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_VARIABLE_DEF,
			rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
			t.TraverseBranch(child, n)
		default:
			panic(jutil.NewIllegalStateExceptionMsg("Invalid class member: " + child.GetToken().String()))
		}
		child = next
	}

	t.currentNode = n
	t.visitNode(n, parent)
}

// Java: NodeTraversal.handleMemberFieldDef(Node, Node)
func (t *NodeTraversal) HandleMemberFieldDef(n rhino.Node, child rhino.Node) {
	previousHoistScopeRoot := t.currentHoistScopeRoot
	t.currentHoistScopeRoot = n
	t.PushScope(child)
	t.TraverseBranch(child, n)
	t.PopScope()
	t.currentHoistScopeRoot = previousHoistScopeRoot
}

// Java: NodeTraversal.traverseChildren(Node)
func (t *NodeTraversal) TraverseChildren(n rhino.Node) {
	cb := t.callback
	noVisitLeaves := t.noVisitLeaves
	onlyNameLeaves := t.onlyNameLeaves
	for child := n.GetFirstChild(); !child.IsNil(); {
		// child could be replaced, in which case our child node
		// would no longer point to the true next
		next := child.GetNext()
		if child.GetFirstChild().IsNil() {
			ctok := child.GetToken()
			if nodeTraversalTokenKind[uint8(ctok)] == 0 {
				if !t.noBranchShouldTraverse {
					t.currentNode = child
					if !cb.ShouldTraverse(t, child, n) {
						child = next
						continue
					}
				}
				if !noVisitLeaves && (!onlyNameLeaves || ctok == rhino.Token_NAME) && !t.noVisit &&
					(t.visitTokens == nil || t.visitTokens.Has(ctok)) {
					t.currentNode = child
					cb.Visit(t, child, n)
				}
				child = next
				continue
			}
		}
		t.TraverseBranch(child, n)
		child = next
	}
}

// Report reports a diagnostic (error or warning)
//
// Java: NodeTraversal.report(Node, DiagnosticType, String...)
func (t *NodeTraversal) Report(n rhino.Node, diagnosticType *DiagnosticType, arguments ...string) {
	error_ := JSError_MakeWithNode(n, diagnosticType, arguments...)
	t.compiler.Report(error_)
}

// ReportWithNodeRange reports a diagnostic (error or warning) This variant
// is particularly useful for reporting on GETPROP nodes whose length is
// otherwise just the last component (`abc` in `foo.abc`).
//
// Java: NodeTraversal.report(Node, Node, DiagnosticType, String...)
func (t *NodeTraversal) ReportWithNodeRange(start rhino.Node, end rhino.Node, diagnosticType *DiagnosticType, arguments ...string) {
	error_ := JSError_MakeWithNodeRange(start, end, diagnosticType, arguments...)
	t.compiler.Report(error_)
}

// Java: NodeTraversal.reportCodeChange()
func (t *NodeTraversal) ReportCodeChange() {
	changeScope := t.currentChangeScope
	if changeScope.IsNil() {
		panic(jutil.NewNullPointerException()) // checkNotNull(changeScope)
	}
	jutil.CheckStateMsg1(ChangeTracker_IsChangeScopeRoot(changeScope), "%s", changeScope)
	t.compiler.ReportChangeToChangeScope(changeScope)
}

// Java: NodeTraversal.reportCodeChange(Node)
func (t *NodeTraversal) ReportCodeChangeWithNode(n rhino.Node) {
	t.compiler.ReportChangeToEnclosingScope(n)
}

// GetCurrentScript returns the SCRIPT node enclosing the current scope.
//
// Throws UnsupportedOperationException if inside a NodeTraversal that does
// not support this operation, or if called while visiting a ROOT node during
// a global traversal. Only traversals that begin at a ROOT node (which
// includes most usages of NodeTraversal) will support this operation. For
// example, NodeTraversal.Builder#traverseAtScope(AbstractScope<?, ?>) is
// unsupported because it can begin at any arbitrary scope root like a
// FUNCTION.
//
// Java: NodeTraversal.getCurrentScript()
func (t *NodeTraversal) GetCurrentScript() rhino.Node {
	if t.currentScript.IsNil() {
		panic(jutil.NewUnsupportedOperationExceptionMsg("getCurrentScript not supported"))
	}
	return t.currentScript
}

// SetChangeScope: n is the current change scope, should be null when the
// traversal is complete.
//
// Java: NodeTraversal.setChangeScope(Node)
func (t *NodeTraversal) SetChangeScope(n rhino.Node) {
	t.currentChangeScope = n
}

// Java: NodeTraversal.getEnclosingScript(Node)
func (t *NodeTraversal) GetEnclosingScript(n rhino.Node) rhino.Node {
	for !n.IsNil() && !n.IsScript() {
		n = n.GetParent()
	}
	return n
}

// Java: NodeTraversal.initTraversal(Node)
func (t *NodeTraversal) InitTraversal(traversalRoot rhino.Node) {
	if Platform_IsThreadInterrupted() {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
	}
	hoistScopeRoot := NodeUtil_GetEnclosingHoistScopeRoot(traversalRoot)
	t.currentHoistScopeRoot = hoistScopeRoot
	// "null" is the global hoist scope root, but we want the current script if any
	// for the "change scope"
	changeScopeStart := traversalRoot
	if !hoistScopeRoot.IsNil() {
		changeScopeStart = hoistScopeRoot
	}
	changeScope := ChangeTracker_GetEnclosingChangeScopeRoot(changeScopeStart)
	t.SetChangeScope(changeScope)

	var enclosingFunction rhino.Node
	if !hoistScopeRoot.IsNil() {
		enclosingFunction = NodeUtil_GetEnclosingFunction(hoistScopeRoot)
	}
	t.currentFunction = enclosingFunction

	script := t.GetEnclosingScript(changeScope)
	t.currentScript = script
	t.ClearScriptState()
}

// ClearScriptState is used to clear any cached state with regard to the
// current script and should be called before traversing a SCRIPT rooted
// subtree.
//
// Java: NodeTraversal.clearScriptState()
func (t *NodeTraversal) ClearScriptState() {
	t.inputId = nil
	t.sourceName = nil
	t.compilerInput = nil
}

// Java: NodeTraversal.getInputId()
func (t *NodeTraversal) GetInputId() *rhino.InputId {
	if !t.currentScript.IsNil() && t.inputId == nil {
		t.inputId = t.currentScript.GetInputId()
	}
	return t.inputId
}

// Java: NodeTraversal.getBestSourceFileName(Node)
func (t *NodeTraversal) GetBestSourceFileName(n rhino.Node) *string {
	if n.IsNil() {
		return t.GetSourceName()
	}
	return n.GetSourceFileName()
}
