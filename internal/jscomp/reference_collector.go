// Port of com/google/javascript/jscomp/ReferenceCollector.java @ closure-compiler v20261015

/*
 * Copyright 2008 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReferenceCollector is a helper class for passes that want to access all
// information about where a variable is referenced and declared at once and
// then make a decision as to how it should be handled, possibly inlining,
// reordering, or generating warnings. Callers do this by providing Behavior
// and then calling Process.
//
// Java: final class ReferenceCollector implements CompilerPass,
// StaticSymbolTable<Var, Reference>.
//
// Go-only: language-adapter — Java narrows StaticSymbolTable's getScope(Var)
// (to Scope) and getReferences(Var) (Iterable<Reference> to
// ReferenceCollection, which callers need and an iter.Seq cannot give back).
// Go has no covariant overrides, so GetScope and GetReferences keep the
// Java-declared types and AsStaticSymbolTable returns the
// rhino.StaticSymbolTable[*Var, *Reference] view of this collector.
type ReferenceCollector struct {
	// Maps a given variable to a collection of references to that name. Note that
	// Var objects are not stable across multiple traversals (unlike scope root or
	// name).
	//
	// Go-only: language-adapter — Var overrides equals/hashCode (ScopedName),
	// hence the Object map.
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]

	// Go-only: opt-in native lookup; the ordered map remains authoritative.
	nativeReferences map[referenceCollectorScopedKey]*ReferenceCollection

	// The stack of basic blocks and scopes the current traversal is in.
	//
	// Go-only: language-adapter — the ArrayDeque, first element first.
	blockStack []*BasicBlock

	// Source of behavior at various points in the traversal.
	behavior ReferenceCollector_Behavior

	scopeCreator ScopeCreator

	// JavaScript compiler to use in traversing.
	compiler AbstractCompiler

	// Only collect references for filtered variables.
	varFilter func(*Var) bool

	callback *ReferenceCollector_CollectorCallback

	collectedHoistedFunctions *jutil.LinkedHashSet[*rhino.Node]

	narrowScope *Scope
}

// NewReferenceCollector: constructor initializes block stack.
//
// Java: ReferenceCollector(AbstractCompiler, Behavior, ScopeCreator)
func NewReferenceCollector(compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) *ReferenceCollector {
	return NewReferenceCollectorWithVarFilter(compiler, behavior, creator, func(*Var) bool { return true })
}

// NewReferenceCollectorWithVarFilter: constructor only collects references
// that match the given variable.
//
// The test for Var equality uses reference equality, so it's necessary to
// inject a scope when you traverse.
//
// Java: ReferenceCollector(AbstractCompiler, Behavior, ScopeCreator, Predicate<Var>)
func NewReferenceCollectorWithVarFilter(
	compiler AbstractCompiler,
	behavior ReferenceCollector_Behavior,
	creator ScopeCreator,
	varFilter func(*Var) bool) *ReferenceCollector {
	c := &ReferenceCollector{
		referenceMap:              jutil.NewObjectLinkedHashMap[*Var, *ReferenceCollection](),
		blockStack:                nil,
		collectedHoistedFunctions: jutil.NewLinkedHashSet[*rhino.Node](),
	}
	c.callback = &ReferenceCollector_CollectorCallback{c: c}
	c.compiler = compiler
	c.behavior = behavior
	c.scopeCreator = creator
	c.varFilter = varFilter
	c.initNativeReferenceLookup()
	return c
}

// Process is a convenience method for running this pass over a tree with
// this class as a callback.
//
// Java: ReferenceCollector.process(Node, Node)
func (c *ReferenceCollector) Process(externs *rhino.Node, root *rhino.Node) {
	c.CreateTraversalBuilder().TraverseRoots(externs, root)
}

// ProcessRoot is the one-argument process overload.
//
// Go-only: language-adapter — the name: CompilerPass fixes the plain name
// to the two-argument overload.
//
// Java: ReferenceCollector.process(Node)
func (c *ReferenceCollector) ProcessRoot(root *rhino.Node) {
	c.CreateTraversalBuilder().Traverse(root)
}

// ProcessScope targets reference collection to a particular scope.
//
// Java: ReferenceCollector.processScope(Scope)
func (c *ReferenceCollector) ProcessScope(scope *Scope) {
	shouldAddToBlockStack := !scope.IsHoistScope()
	c.narrowScope = scope
	if shouldAddToBlockStack {
		c.PushNewBlock(scope.GetRootNode())
	}
	c.CreateTraversalBuilder().TraverseAtScope(scope)
	if shouldAddToBlockStack {
		c.PopLastBlock(scope.GetRootNode())
	}
	c.narrowScope = nil
}

// CreateTraversalBuilder: same as process but only runs on a part of AST
// associated to one script.
//
// Java: ReferenceCollector.createTraversalBuilder()
func (c *ReferenceCollector) CreateTraversalBuilder() *NodeTraversal_Builder {
	return NodeTraversal_Builder_().
		SetCompiler(c.compiler).
		SetCallback(c.callback).
		SetScopeCreator(c.scopeCreator).
		SetObeyDestructuringAndDefaultValueExecutionOrder(true)
}

// GetAllSymbols gets the variables that were referenced in this callback.
//
// Java: ReferenceCollector.getAllSymbols()
func (c *ReferenceCollector) GetAllSymbols() iter.Seq[*Var] {
	return c.referenceMap.Keys()
}

// GetScope narrows StaticSymbolTable.getScope; see ReferenceCollector.
//
// Java: ReferenceCollector.getScope(Var)
func (c *ReferenceCollector) GetScope(var_ *Var) *Scope {
	return var_.getScopeS()
}

// GetReferences gets the reference collection for the given variable. It
// narrows StaticSymbolTable.getReferences; see ReferenceCollector.
//
// Java: ReferenceCollector.getReferences(Var)
func (c *ReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	if c.nativeReferences != nil {
		return c.nativeReferences[referenceCollectorKey(v)]
	}
	return c.referenceMap.Get(v)
}

// AsStaticSymbolTable returns this collector as the Java
// StaticSymbolTable<Var, Reference> it implements (for
// SymbolTable.addSymbolsFrom).
//
// Go-only: language-adapter — see ReferenceCollector.
func (c *ReferenceCollector) AsStaticSymbolTable() rhino.StaticSymbolTable[*Var, *Reference] {
	return referenceCollector_staticSymbolTable{c}
}

// referenceCollector_staticSymbolTable is the
// rhino.StaticSymbolTable[*Var, *Reference] view of a ReferenceCollector:
// the same three methods, with StaticSymbolTable's Go result types.
//
// Go-only: language-adapter — see ReferenceCollector.
type referenceCollector_staticSymbolTable struct {
	c *ReferenceCollector
}

// Java: ReferenceCollector.getReferences(Var)
func (t referenceCollector_staticSymbolTable) GetReferences(symbol *Var) iter.Seq[*Reference] {
	refs := t.c.GetReferences(symbol)
	if refs == nil {
		return nil
	}
	return refs.All()
}

// Java: ReferenceCollector.getScope(Var)
func (t referenceCollector_staticSymbolTable) GetScope(symbol *Var) rhino.StaticScope {
	if s := t.c.GetScope(symbol); s != nil {
		return s
	}
	return nil
}

// Java: ReferenceCollector.getAllSymbols()
func (t referenceCollector_staticSymbolTable) GetAllSymbols() iter.Seq[*Var] {
	return t.c.GetAllSymbols()
}

// Java: ReferenceCollector.maybeJumpToHoistedFunction(Var, NodeTraversal)
func (c *ReferenceCollector) MaybeJumpToHoistedFunction(var_ *Var, t *NodeTraversal) {
	fnNode := var_.GetParentNode()
	varScope := var_.getScopeS()

	if fnNode == nil ||
		!NodeUtil_IsHoistedFunctionDeclaration(fnNode) ||
		// If we're only traversing a narrow scope, do not try to climb outside.
		(c.narrowScope != nil && c.narrowScope.GetDepth() > varScope.GetDepth()) ||
		c.collectedHoistedFunctions.Contains(fnNode) {
		return
	}

	/*
	 * Replace the block stack with a new one matching the hoist position.
	 *
	 * <p>This algorithm only works because we know hoisted functions cannot be inside loops. It
	 * will have to change if we ever do general function continuations.
	 *
	 * <p>This is tricky to compute because of the weird traverseAtScope call for
	 * AggressiveInlineAliases.
	 */
	oldBlockStack := c.blockStack
	c.blockStack = nil

	if varScope.IsGlobal() {
		c.blockStack = append(c.blockStack, referenceCollector_getFirst(oldBlockStack))
	} else {
		for _, b := range oldBlockStack {
			c.blockStack = append(c.blockStack, b)
			if b.GetRoot() == varScope.GetRootNode() {
				break
			}
		}
	}

	/*
	 * Record the function declaration reference explicitly because it is a bleeding name. The
	 * reference must be recorded with the containing scope, and traverseAtScope skips bleeding
	 * function NAME nodes
	 */
	c.AddReference(var_, fnNode.GetFirstChild(), t)
	c.CreateTraversalBuilder().TraverseAtScope(c.scopeCreator.CreateScope(fnNode, varScope))

	c.blockStack = oldBlockStack
}

// referenceCollector_getFirst is ArrayDeque.getFirst().
//
// Go-only: language-adapter — see ReferenceCollector.blockStack.
func referenceCollector_getFirst(stack []*BasicBlock) *BasicBlock {
	if len(stack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return stack[0]
}

// ReferenceCollector_CollectorCallback is the traversal callback.
type ReferenceCollector_CollectorCallback struct {
	c *ReferenceCollector
}

// ShouldTraverse updates block stack.
//
// Java: ReferenceCollector.CollectorCallback.shouldTraverse(NodeTraversal, Node, Node)
func (cb *ReferenceCollector_CollectorCallback) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	// We automatically traverse a hoisted function body when that function
	// is first referenced, so that the reference lists are in the right order.
	//
	// TODO(nicksantos): Maybe generalize this to a continuation mechanism
	// like in RemoveUnusedCode.
	if NodeUtil_IsHoistedFunctionDeclaration(n) {
		if !cb.c.collectedHoistedFunctions.Add(n) {
			return false
		}
	}

	// If node is a new basic block, put on basic block stack
	if referenceCollector_isBlockBoundary(n, parent) {
		cb.c.PushNewBlock(n)
	}

	return true
}

// EnterScope updates block stack and invokes any additional behavior.
//
// Java: ReferenceCollector.CollectorCallback.enterScope(NodeTraversal)
func (cb *ReferenceCollector_CollectorCallback) EnterScope(t *NodeTraversal) {
	// Don't add all ES6 scope roots to blockStack, only those that are also scopes according to
	// the ES5 scoping rules. Other nodes that ought to be considered the root of a BasicBlock
	// are added in shouldTraverse() or processScope() and removed in visit().
	if t.IsHoistScope() {
		cb.c.PushNewBlock(t.GetScopeRoot())
	}
}

// ExitScope updates block stack and invokes any additional behavior.
//
// Java: ReferenceCollector.CollectorCallback.exitScope(NodeTraversal)
func (cb *ReferenceCollector_CollectorCallback) ExitScope(t *NodeTraversal) {
	if t.IsHoistScope() {
		cb.c.PopLastBlock(t.GetScopeRoot())
	}
	cb.c.behavior.AfterExitScope(t, cb.c.referenceMapWrapper())
}

// Visit: for each node, update the block stack and reference collection as
// appropriate.
//
// Java: ReferenceCollector.CollectorCallback.visit(NodeTraversal, Node, Node)
func (cb *ReferenceCollector_CollectorCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() || n.IsImportStar() {
		if (parent.IsImportSpec() && n != parent.GetLastChild()) ||
			(parent.IsExportSpec() && n != parent.GetFirstChild()) {
			// The n in `import {n as x}` or `export {x as n}` are not references, even though
			// they are represented in the AST as NAME nodes.
			return
		}

		v := t.GetScope().GetVar(n.GetString())
		if v != nil {
			cb.c.AddReference(v, n, t)
			cb.c.MaybeJumpToHoistedFunction(v, t)
		}
	}

	if referenceCollector_isBlockBoundary(n, parent) {
		cb.c.PopLastBlock(n)
	}
}

// Java: ReferenceCollector.pushNewBlock(Node)
func (c *ReferenceCollector) PushNewBlock(root *rhino.Node) {
	var peekLast *BasicBlock
	if len(c.blockStack) > 0 {
		peekLast = c.blockStack[len(c.blockStack)-1]
	}
	c.blockStack = append(c.blockStack, NewBasicBlock(peekLast, root))
}

// Java: ReferenceCollector.popLastBlock(Node)
func (c *ReferenceCollector) PopLastBlock(root *rhino.Node) {
	if len(c.blockStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.removeLast()
	}
	last := c.blockStack[len(c.blockStack)-1]
	c.blockStack[len(c.blockStack)-1] = nil
	c.blockStack = c.blockStack[:len(c.blockStack)-1]
	// Verfiy that the stack is unwound correctly.
	jutil.CheckStateMsg1(base.JSCompObjects_Identical(root, last.GetRoot()), "%s", root)
}

// referenceCollector_isBlockBoundary returns true if this node marks the
// start of a new basic block
//
// Java: ReferenceCollector.isBlockBoundary(Node, Node)
func referenceCollector_isBlockBoundary(n *rhino.Node, parent *rhino.Node) bool {
	if parent != nil {
		switch parent.GetToken() {
		case rhino.Token_DO, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
			rhino.Token_TRY, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_CLASS, rhino.Token_SWITCH_BODY:
			// NOTE: TRY has up to 3 child blocks:
			// TRY
			//   BLOCK
			//   BLOCK
			//     CATCH
			//   BLOCK
			// Note that there is an explicit CATCH token but no explicit
			// FINALLY token. For simplicity, we consider each BLOCK
			// a separate basic BLOCK.
			return true
		case rhino.Token_AND,
			rhino.Token_HOOK,
			rhino.Token_IF,
			rhino.Token_OR,
			rhino.Token_COALESCE,
			rhino.Token_OPTCHAIN_GETPROP,
			rhino.Token_OPTCHAIN_GETELEM,
			rhino.Token_OPTCHAIN_CALL,
			rhino.Token_DEFAULT_VALUE:
			// The first child of a conditional is not a boundary,
			// but all the rest of the children are.
			return n != parent.GetFirstChild()
		default:
		}
	}

	return n.IsCase()
}

// Java: ReferenceCollector.addReference(Var, Node, NodeTraversal)
func (c *ReferenceCollector) AddReference(v *Var, n *rhino.Node, t *NodeTraversal) {
	if !c.varFilter(v) {
		return
	}

	// Create collection if none already
	collection := c.GetReferences(v)
	if collection == nil {
		collection = NewReferenceCollection()
		c.referenceMap.Put(v, collection)
		if c.nativeReferences != nil {
			c.nativeReferences[referenceCollectorKey(v)] = collection
		}
	}

	// Add this particular reference
	if len(c.blockStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.getLast()
	}
	collection.Add(NewReference(n, t, c.blockStack[len(c.blockStack)-1]))
}

// ReferenceCollector_ReferenceMapWrapper is the ReferenceMap handed to
// Behavior.afterExitScope.
type ReferenceCollector_ReferenceMapWrapper struct {
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]
	// A raw-map escape permanently retires the owner's private lookup index.
	// The public constructor has no indexed owner and retains its old behavior.
	owner *ReferenceCollector
}

// Java: ReferenceCollector.ReferenceMapWrapper(Map<Var, ReferenceCollection>)
func NewReferenceCollector_ReferenceMapWrapper(
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]) *ReferenceCollector_ReferenceMapWrapper {
	return &ReferenceCollector_ReferenceMapWrapper{referenceMap: referenceMap}
}

// Java: ReferenceCollector.ReferenceMapWrapper.getReferences(Var)
func (w *ReferenceCollector_ReferenceMapWrapper) GetReferences(var_ *Var) *ReferenceCollection {
	if w.owner != nil {
		return w.owner.GetReferences(var_)
	}
	return w.referenceMap.Get(var_)
}

// Java: ReferenceCollector.ReferenceMapWrapper.getRawReferenceMap()
func (w *ReferenceCollector_ReferenceMapWrapper) GetRawReferenceMap() *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection] {
	if w.owner != nil {
		w.owner.nativeReferences = nil
	}
	return w.referenceMap
}

// Java: ReferenceCollector.ReferenceMapWrapper.toString()
func (w *ReferenceCollector_ReferenceMapWrapper) String() string {
	return w.referenceMap.String()
}

// ReferenceCollector_Behavior is the way for callers to add specific
// behavior during traversal that utilizes the built-up reference
// information.
type ReferenceCollector_Behavior interface {
	// AfterExitScope is called after we finish with a scope.
	//
	// Java: ReferenceCollector.Behavior.afterExitScope(NodeTraversal, ReferenceMap)
	AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap)
}

// ReferenceCollector_BehaviorFunc is a Behavior given as a function (a Java
// lambda for Behavior).
//
// Go-only: language-adapter — the adapter for lambdas implementing the
// single-method interface Behavior (PORTING.md §4).
type ReferenceCollector_BehaviorFunc func(t *NodeTraversal, referenceMap ReferenceMap)

// Java: ReferenceCollector.Behavior.afterExitScope(NodeTraversal, ReferenceMap)
func (f ReferenceCollector_BehaviorFunc) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	f(t, referenceMap)
}

var ReferenceCollector_DO_NOTHING_BEHAVIOR ReferenceCollector_Behavior = &referenceCollector_doNothingBehavior{}

// referenceCollector_doNothingBehavior is the anonymous class of
// DO_NOTHING_BEHAVIOR.
type referenceCollector_doNothingBehavior struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: ReferenceCollector.DO_NOTHING_BEHAVIOR.afterExitScope(NodeTraversal, ReferenceMap)
func (*referenceCollector_doNothingBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
}
