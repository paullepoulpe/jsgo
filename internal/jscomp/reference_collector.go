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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type basicBlockEntry struct {
	root  *rhino.Node
	block *BasicBlock
}

// scopedVarKey is what makes two Vars equal (ScopedName equality).
type scopedVarKey struct {
	scopeRoot *rhino.Node
	nameID    uint32
}

func scopedVarKeyOf(v *Var) scopedVarKey {
	return scopedVarKey{scopeRoot: v.scope.GetRootNode(), nameID: v.nameID}
}

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
	// Go-only: language-adapter — the LinkedHashMap<Var, ReferenceCollection>
	// is refs, keyed by what makes Vars equal, plus vars, its keys in
	// insertion order.
	refs map[scopedVarKey]*ReferenceCollection
	vars []*Var

	// Go-only: performance — with reclaimLocals set, the vars of local scopes
	// are listed in localVars instead of vars, and those past localBase are
	// forgotten when the traversal leaves an outermost local scope (unless
	// canReclaimLocals says otherwise), so that a pass does not hold the
	// references of every function at once.
	reclaimLocals    bool
	canReclaimLocals func() bool
	localVars        []*Var
	localBase        int
	localScopeDepth  int

	lastVar         *Var
	lastCollection  *ReferenceCollection
	lastRejectedVar *Var

	// The stack of basic blocks and scopes the current traversal is in.
	//
	// Go-only: language-adapter — the ArrayDeque, first element first. Go-only:
	// performance — a block is created only when a reference needs it.
	blockStack []basicBlockEntry

	// Source of behavior at various points in the traversal.
	behavior ReferenceCollector_Behavior

	scopeCreator ScopeCreator

	// JavaScript compiler to use in traversing.
	compiler AbstractCompiler

	// Only collect references for filtered variables.
	varFilter func(*Var) bool

	callback            ReferenceCollector_CollectorCallback
	referenceMapWrapper ReferenceCollector_ReferenceMapWrapper

	collectedHoistedFunctions map[*rhino.Node]struct{}

	narrowScope *Scope
}

// NewReferenceCollector: constructor initializes block stack.
//
// Java: ReferenceCollector(AbstractCompiler, Behavior, ScopeCreator)
func NewReferenceCollector(compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) *ReferenceCollector {
	return NewReferenceCollectorWithVarFilter(compiler, behavior, creator, nil)
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
		compiler:     compiler,
		behavior:     behavior,
		scopeCreator: creator,
		varFilter:    varFilter,
	}
	c.callback.c = c
	c.referenceMapWrapper.c = c
	return c
}

func (c *ReferenceCollector) enableLocalReclamation(canReclaim func() bool) {
	c.reclaimLocals = true
	c.canReclaimLocals = canReclaim
}

// clearLocalBase lets the next reclamation forget every local var.
func (c *ReferenceCollector) clearLocalBase() {
	c.localBase = 0
}

func (c *ReferenceCollector) maybeReclaimLocals() {
	if c.narrowScope != nil {
		return
	}
	if c.canReclaimLocals != nil && !c.canReclaimLocals() {
		c.localBase = len(c.localVars)
		return
	}
	for _, v := range c.localVars[c.localBase:] {
		delete(c.refs, scopedVarKeyOf(v))
	}
	clear(c.localVars[c.localBase:])
	c.localVars = c.localVars[:c.localBase]
	c.lastVar = nil
	c.lastCollection = nil
}

func (c *ReferenceCollector) getOrCreateReferences(v *Var) *ReferenceCollection {
	if v == c.lastVar && v != nil {
		return c.lastCollection
	}
	key := scopedVarKeyOf(v)
	col := c.refs[key]
	if col == nil {
		col = &ReferenceCollection{varScope: v.scope}
		if c.refs == nil {
			c.refs = make(map[scopedVarKey]*ReferenceCollection)
		}
		c.refs[key] = col
		if c.reclaimLocals && v.scope.parent != nil {
			c.localVars = append(c.localVars, v)
		} else {
			c.vars = append(c.vars, v)
		}
	}
	c.lastVar = v
	c.lastCollection = col
	return col
}

func (c *ReferenceCollector) getReferencesByScopedNameID(scopeRoot *rhino.Node, nameID uint32) *ReferenceCollection {
	return c.refs[scopedVarKey{scopeRoot: scopeRoot, nameID: nameID}]
}

func (c *ReferenceCollector) builderVal() NodeTraversal_Builder {
	return NodeTraversal_Builder{
		compiler:                                       c.compiler,
		callback:                                       &c.callback,
		scopeCreator:                                   c.scopeCreator,
		obeyDestructuringAndDefaultValueExecutionOrder: true,
	}
}

// Process is a convenience method for running this pass over a tree with
// this class as a callback.
//
// Java: ReferenceCollector.process(Node, Node)
func (c *ReferenceCollector) Process(externs *rhino.Node, root *rhino.Node) {
	b := c.builderVal()
	b.TraverseRoots(externs, root)
}

// ProcessRoot is the one-argument process overload.
//
// Go-only: language-adapter — the name: CompilerPass fixes the plain name
// to the two-argument overload.
//
// Java: ReferenceCollector.process(Node)
func (c *ReferenceCollector) ProcessRoot(root *rhino.Node) {
	b := c.builderVal()
	b.Traverse(root)
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
	b := c.builderVal()
	b.TraverseAtScope(scope)
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
		SetCallback(&c.callback).
		SetScopeCreator(c.scopeCreator).
		SetObeyDestructuringAndDefaultValueExecutionOrder(true)
}

// GetAllSymbols gets the variables that were referenced in this callback.
//
// Java: ReferenceCollector.getAllSymbols()
func (c *ReferenceCollector) GetAllSymbols() iter.Seq[*Var] {
	if len(c.localVars) == 0 {
		return slices.Values(c.vars)
	}
	return func(yield func(*Var) bool) {
		for _, v := range c.vars {
			if !yield(v) {
				return
			}
		}
		for _, v := range c.localVars {
			if !yield(v) {
				return
			}
		}
	}
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
	if v == nil {
		return nil
	}
	if v == c.lastVar {
		return c.lastCollection
	}
	col := c.refs[scopedVarKeyOf(v)]
	if col != nil {
		c.lastVar = v
		c.lastCollection = col
	}
	return col
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

func rcIsHoistedFunctionDeclaration(fnNode *rhino.Node) bool {
	parent := fnNode.GetParent()
	if parent == nil {
		return false
	}
	switch parent.GetToken() {
	case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_EXPORT:
		first := fnNode.GetFirstChild()
		return first != nil && first.IsName() && first.GetString() != ""
	case rhino.Token_BLOCK, rhino.Token_LABEL, rhino.Token_DECLARE, rhino.Token_NAMESPACE_ELEMENTS:
		if gp := parent.GetParent(); gp != nil && gp.IsFunction() {
			first := fnNode.GetFirstChild()
			return first != nil && first.IsName() && first.GetString() != ""
		}
	}
	return false
}

// Java: ReferenceCollector.maybeJumpToHoistedFunction(Var, NodeTraversal)
func (c *ReferenceCollector) MaybeJumpToHoistedFunction(var_ *Var, t *NodeTraversal) {
	nameNode := var_.nameNode
	if nameNode == nil {
		return
	}
	fnNode := nameNode.GetParent()
	if fnNode == nil || !fnNode.IsFunction() {
		return
	}
	if _, ok := c.collectedHoistedFunctions[fnNode]; ok || !rcIsHoistedFunctionDeclaration(fnNode) {
		return
	}
	varScope := var_.getScopeS()

	// If we're only traversing a narrow scope, do not try to climb outside.
	if c.narrowScope != nil && c.narrowScope.GetDepth() > varScope.GetDepth() {
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
	if varScope.IsGlobal() {
		if len(oldBlockStack) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		if oldBlockStack[0].block == nil {
			c.materializeBlock(0)
		}
		c.blockStack = append(oldBlockStack[len(oldBlockStack):len(oldBlockStack)], oldBlockStack[0])
	} else {
		rootNode := varScope.GetRootNode()
		end := len(oldBlockStack)
		for i := range oldBlockStack {
			if oldBlockStack[i].root == rootNode {
				end = i + 1
				break
			}
		}
		if end > 0 && oldBlockStack[end-1].block == nil {
			c.materializeBlock(end - 1)
		}
		c.blockStack = append(oldBlockStack[len(oldBlockStack):len(oldBlockStack)], oldBlockStack[:end]...)
	}

	/*
	 * Record the function declaration reference explicitly because it is a bleeding name. The
	 * reference must be recorded with the containing scope, and traverseAtScope skips bleeding
	 * function NAME nodes
	 */
	c.AddReference(var_, fnNode.GetFirstChild(), t)
	b := c.builderVal()
	b.TraverseAtScope(c.scopeCreator.CreateScope(fnNode, varScope))

	clear(c.blockStack)
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

func (c *ReferenceCollector) getCurrentBlock() *BasicBlock {
	idx := len(c.blockStack) - 1
	if b := c.blockStack[idx].block; b != nil {
		return b
	}
	return c.materializeBlock(idx)
}

func (c *ReferenceCollector) materializeBlock(idx int) *BasicBlock {
	i := idx
	for i > 0 && c.blockStack[i-1].block == nil {
		i--
	}
	var parent *BasicBlock
	if i > 0 {
		parent = c.blockStack[i-1].block
	}
	for ; i <= idx; i++ {
		parent = NewBasicBlock(parent, c.blockStack[i].root)
		c.blockStack[i].block = parent
	}
	return parent
}

// ReferenceCollector_CollectorCallback is the traversal callback.
type ReferenceCollector_CollectorCallback struct {
	c *ReferenceCollector
}

// ShouldTraverse updates block stack.
//
// Java: ReferenceCollector.CollectorCallback.shouldTraverse(NodeTraversal, Node, Node)
func (cb *ReferenceCollector_CollectorCallback) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.GetFirstChild() == nil {
		switch n.GetToken() {
		case rhino.Token_NAME, rhino.Token_IMPORT_STAR, rhino.Token_BLOCK, rhino.Token_SWITCH_BODY:
			if referenceCollector_isBlockBoundary(n, parent) {
				cb.c.PushNewBlock(n)
			}
			return true
		default:
			return false
		}
	}

	// We automatically traverse a hoisted function body when that function
	// is first referenced, so that the reference lists are in the right order.
	//
	// TODO(nicksantos): Maybe generalize this to a continuation mechanism
	// like in RemoveUnusedCode.
	if n.IsFunction() && rcIsHoistedFunctionDeclaration(n) {
		c := cb.c
		if _, ok := c.collectedHoistedFunctions[n]; ok {
			return false
		}
		if c.collectedHoistedFunctions == nil {
			c.collectedHoistedFunctions = make(map[*rhino.Node]struct{})
		}
		c.collectedHoistedFunctions[n] = struct{}{}
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
	if !t.InGlobalScope() {
		cb.c.localScopeDepth++
	}
	// Don't add all ES6 scope roots to blockStack, only those that are also scopes according to
	// the ES5 scoping rules. Other nodes that ought to be considered the root of a BasicBlock
	// are added in shouldTraverse() or processScope() and removed in visit().
	if root := t.GetScopeRoot(); nodeTraversal_isHoistScopeRootNode(root) {
		cb.c.PushNewBlock(root)
	}
}

// ExitScope updates block stack and invokes any additional behavior.
//
// Java: ReferenceCollector.CollectorCallback.exitScope(NodeTraversal)
func (cb *ReferenceCollector_CollectorCallback) ExitScope(t *NodeTraversal) {
	if root := t.GetScopeRoot(); nodeTraversal_isHoistScopeRootNode(root) {
		cb.c.PopLastBlock(root)
	}
	cb.c.behavior.AfterExitScope(t, &cb.c.referenceMapWrapper)
	if !t.InGlobalScope() {
		cb.c.localScopeDepth--
		if cb.c.localScopeDepth == 0 && cb.c.reclaimLocals {
			cb.c.maybeReclaimLocals()
		}
	}
}

// Visit: for each node, update the block stack and reference collection as
// appropriate.
//
// Java: ReferenceCollector.CollectorCallback.visit(NodeTraversal, Node, Node)
func (cb *ReferenceCollector_CollectorCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if tok := n.GetToken(); tok == rhino.Token_NAME || tok == rhino.Token_IMPORT_STAR {
		if parent != nil {
			if pt := parent.GetToken(); (pt == rhino.Token_IMPORT_SPEC && n != parent.GetLastChild()) ||
				(pt == rhino.Token_EXPORT_SPEC && n != parent.GetFirstChild()) {
				// The n in `import {n as x}` or `export {x as n}` are not references, even though
				// they are represented in the AST as NAME nodes.
				return
			}
		}

		id := n.GetStringID()
		var v *Var
		var scope *Scope
		if id == implicitArgumentsID {
			scope = t.GetScope()
			v = scope.GetVarByID(id)
		} else {
			v = t.LookupVarByID(id)
		}
		if v != nil {
			cb.c.addReferenceWithScope(v, n, t, scope)
			if v.nameNode != n {
				cb.c.MaybeJumpToHoistedFunction(v, t)
			}
		}
	}

	if referenceCollector_isBlockBoundary(n, parent) {
		cb.c.PopLastBlock(n)
	}
}

// Java: ReferenceCollector.pushNewBlock(Node)
func (c *ReferenceCollector) PushNewBlock(root *rhino.Node) {
	c.blockStack = append(c.blockStack, basicBlockEntry{root: root})
}

// Java: ReferenceCollector.popLastBlock(Node)
func (c *ReferenceCollector) PopLastBlock(root *rhino.Node) {
	if len(c.blockStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.removeLast()
	}
	idx := len(c.blockStack) - 1
	lastRoot := c.blockStack[idx].root
	c.blockStack[idx] = basicBlockEntry{}
	c.blockStack = c.blockStack[:idx]
	// Verfiy that the stack is unwound correctly.
	if lastRoot != root {
		jutil.CheckStateMsg1(base.JSCompObjects_Identical(root, lastRoot), "%s", root)
	}
}

var refCollectorParentBoundaryKind = [256]uint8{
	rhino.Token_DO:               1,
	rhino.Token_FOR:              1,
	rhino.Token_FOR_IN:           1,
	rhino.Token_FOR_OF:           1,
	rhino.Token_FOR_AWAIT_OF:     1,
	rhino.Token_TRY:              1,
	rhino.Token_WHILE:            1,
	rhino.Token_WITH:             1,
	rhino.Token_CLASS:            1,
	rhino.Token_SWITCH_BODY:      1,
	rhino.Token_AND:              2,
	rhino.Token_HOOK:             2,
	rhino.Token_IF:               2,
	rhino.Token_OR:               2,
	rhino.Token_COALESCE:         2,
	rhino.Token_OPTCHAIN_GETPROP: 2,
	rhino.Token_OPTCHAIN_GETELEM: 2,
	rhino.Token_OPTCHAIN_CALL:    2,
	rhino.Token_DEFAULT_VALUE:    2,
}

// referenceCollector_isBlockBoundary returns true if this node marks the
// start of a new basic block
//
// Java: ReferenceCollector.isBlockBoundary(Node, Node)
func referenceCollector_isBlockBoundary(n *rhino.Node, parent *rhino.Node) bool {
	if parent != nil {
		switch refCollectorParentBoundaryKind[uint8(parent.GetToken())] {
		case 1:
			return true
		case 2:
			return n != parent.GetFirstChild()
		}
	}

	return n.IsCase()
}

// Java: ReferenceCollector.addReference(Var, Node, NodeTraversal)
func (c *ReferenceCollector) AddReference(v *Var, n *rhino.Node, t *NodeTraversal) {
	c.addReferenceWithScope(v, n, t, nil)
}

func (c *ReferenceCollector) addReferenceWithScope(v *Var, n *rhino.Node, t *NodeTraversal, scope *Scope) {
	if v != c.lastVar && c.varFilter != nil {
		if v == c.lastRejectedVar || !c.varFilter(v) {
			c.lastRejectedVar = v
			return
		}
	}
	collection := c.getOrCreateReferences(v)

	// Add this particular reference
	if len(c.blockStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.getLast()
	}
	if scope == nil {
		scope = t.GetScope()
	}
	inputId := t.inputId
	if t.compilerInput == nil {
		inputId = t.GetInput().GetInputId()
	}
	collection.Add(newReferenceWithScopeAndInputId(n, c.getCurrentBlock(), scope, inputId))
}

// ReferenceCollector_ReferenceMapWrapper is the ReferenceMap handed to
// Behavior.afterExitScope.
type ReferenceCollector_ReferenceMapWrapper struct {
	c            *ReferenceCollector
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]
}

// Java: ReferenceCollector.ReferenceMapWrapper(Map<Var, ReferenceCollection>)
func NewReferenceCollector_ReferenceMapWrapper(
	referenceMap *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection]) *ReferenceCollector_ReferenceMapWrapper {
	return &ReferenceCollector_ReferenceMapWrapper{referenceMap: referenceMap}
}

// Java: ReferenceCollector.ReferenceMapWrapper.getReferences(Var)
func (w *ReferenceCollector_ReferenceMapWrapper) GetReferences(var_ *Var) *ReferenceCollection {
	if w.c != nil {
		return w.c.GetReferences(var_)
	}
	return w.referenceMap.Get(var_)
}

func (w *ReferenceCollector_ReferenceMapWrapper) getReferencesByScopedName(scopeRoot *rhino.Node, name string) *ReferenceCollection {
	if w.c != nil {
		id := rhino.RhinoStringPool_LookupID(name)
		if id == 0 {
			return nil
		}
		return w.c.getReferencesByScopedNameID(scopeRoot, id)
	}
	return nil
}

// Java: ReferenceCollector.ReferenceMapWrapper.getRawReferenceMap()
func (w *ReferenceCollector_ReferenceMapWrapper) GetRawReferenceMap() *jutil.ObjectLinkedHashMap[*Var, *ReferenceCollection] {
	if w.c == nil {
		return w.referenceMap
	}
	// Go-only: language-adapter — a copy, as the collector does not keep an
	// ObjectLinkedHashMap.
	m := jutil.NewObjectLinkedHashMap[*Var, *ReferenceCollection]()
	for v := range w.c.GetAllSymbols() {
		m.Put(v, w.c.refs[scopedVarKeyOf(v)])
	}
	return m
}

// Java: ReferenceCollector.ReferenceMapWrapper.toString()
func (w *ReferenceCollector_ReferenceMapWrapper) String() string {
	return w.GetRawReferenceMap().String()
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
