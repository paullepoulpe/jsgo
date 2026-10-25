// Port of com/google/javascript/jscomp/InferConsts.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// InferConsts attaches the CONST_VAR annotation to any variable that's
// 1) Provably well-defined and assigned once in its lifetime.
// 2) Annotated 'const'
// 3) Declared with the 'const' keyword.
//
// These 3 are considered semantically equivalent. Notice that a variable in
// a loop is never considered const.
//
// Note that criteria (1) is only used for normal code, not externs.
type InferConsts struct {
	compiler AbstractCompiler
}

// Java: InferConsts(AbstractCompiler)
func NewInferConsts(compiler AbstractCompiler) *InferConsts {
	return &InferConsts{compiler: compiler}
}

func (*InferConsts) onlyNameLeaves() {}

// Java: InferConsts.process(Node, Node)
func (c *InferConsts) Process(externs rhino.Node, js rhino.Node) {
	scopeCreator := NewSyntacticScopeCreator(c.compiler)
	collector := newInferConstsCollector(c, c.compiler, scopeCreator)
	collector.processRoot(js)

	collector.flushLocalStates()
	for i := 1; i < len(collector.globalStates); i++ {
		st := &collector.globalStates[i]
		st.v.nameHash = 0
		c.considerVarState(st)
	}

	globalExternsScope := scopeCreator.CreateScopeOfScope(externs, nil)
	for v := range globalExternsScope.GetAllSymbols() {
		c.considerExternVar(v)
	}
}

func (c *InferConsts) considerExternVar(v *Var) {
	if v.IsImplicitGoogNamespace() {
		return // no name node for provided variables.
	}
	nameNode := v.GetNameNode()
	if !nameNode.IsNil() && nameNode.GetParent().IsConst() {
		nameNode.SetDeclaredConstantVar(true)
	} else if docInfo := v.GetJSDocInfo(); docInfo != nil && docInfo.IsConstant() {
		nameNode.SetDeclaredConstantVar(true)
	}
}

func (c *InferConsts) considerVarState(st *inferConstVarState) {
	v := st.v
	if v.IsImplicitGoogNamespace() {
		return // no name node for provided variables.
	}
	nameNode := v.GetNameNode()
	if !nameNode.IsNil() && nameNode.GetParent().IsConst() {
		nameNode.SetDeclaredConstantVar(true)
	} else if docInfo := v.GetJSDocInfo(); docInfo != nil && docInfo.IsConstant() {
		nameNode.SetDeclaredConstantVar(true)
	}
	if inferConsts_isInferredConstState(v, nameNode, st) {
		nameNode.SetInferredConstantVar(true)
	}
}

func inferConsts_isInferredConstState(v *Var, nameNode rhino.Node, st *inferConstVarState) bool {
	if nameNode.IsNil() || st.assignCount != 1 || !st.isAssignedOnceInLifetime() {
		return false
	}
	declarationType := v.DeclarationType()

	// A switch on a null Token throws NullPointerException.
	switch *declarationType {
	case rhino.Token_LET:
		// Check that non-destructuring let names are actually assigned at declaration.
		return !nameNode.GetParent().IsLet() || nameNode.HasChildren()
	case rhino.Token_CONST,
		rhino.Token_CATCH,
		rhino.Token_CLASS,
		rhino.Token_PARAM_LIST, // Parameters cannot be referenced before declaration.
		rhino.Token_FUNCTION:   // Function hoisting means no references before declaration.
		return true
	case rhino.Token_IMPORT:
		// ES module exports are mutable.
		// TODO(lharker): make this check smarter if we start optimizing unrewritten modules.
		return false
	case rhino.Token_VAR:
		// var hoisting requires this extra work to make sure the 'declaration' is also the first
		// reference.
		if !st.firstIsInitDecl {
			return false
		}
		jutil.CheckState(st.firstIsDecl)
		return st.allAfterFirst
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unrecognized declaration type " + declarationType.String()))
	}
}

type inferConstVarState struct {
	v               *Var
	firstBlock      *BasicBlock
	assignBlock     *BasicBlock
	assignHoistRoot rhino.Node
	assignCount     int32
	firstIsInitDecl bool
	firstIsDecl     bool
	allAfterFirst   bool
}

func (st *inferConstVarState) isAssignedOnceInLifetime() bool {
	varScope := st.v.scope
	for block := st.assignBlock; block != nil; block = block.parent {
		if block.isFunction {
			if varScope.GetClosestHoistScope().GetRootNode() != st.assignHoistRoot {
				return false
			}
			break
		}

		if block.parent == nil && !block.isLoop {
			break
		}

		if block.root == varScope.GetRootNode() {
			break
		}

		if block.isLoop {
			return false
		}
	}
	return true
}

type inferConstsCollector struct {
	outer            *InferConsts
	compiler         AbstractCompiler
	scopeCreator     ScopeCreator
	sc               *SyntacticScopeCreator
	localStates      []inferConstVarState
	globalStates     []inferConstVarState
	blockStack       []basicBlockEntry
	blocks           jutil.Slab[BasicBlock] // one pass's blocks die together
	collectedHoisted map[rhino.Node]struct{}
	jumpDepth        int
}

func newInferConstsCollector(outer *InferConsts, compiler AbstractCompiler, sc *SyntacticScopeCreator) *inferConstsCollector {
	return &inferConstsCollector{
		outer:        outer,
		compiler:     compiler,
		scopeCreator: sc,
		sc:           sc,
		localStates:  make([]inferConstVarState, 1),
		globalStates: make([]inferConstVarState, 1),
	}
}

func (c *inferConstsCollector) flushLocalStates() {
	for i := 1; i < len(c.localStates); i++ {
		st := &c.localStates[i]
		st.v.nameHash = 0
		c.outer.considerVarState(st)
	}
	clear(c.localStates[1:])
	c.localStates = c.localStates[:1]
}

func (c *inferConstsCollector) builderVal() NodeTraversal_Builder {
	return NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: c.scopeCreator,
		obeyDestructuringAndDefaultValueExecutionOrder: true,
	}
}

func (c *inferConstsCollector) processRoot(root rhino.Node) {
	b := c.builderVal()
	b.Traverse(root)
}

func (c *inferConstsCollector) pushNewBlock(root rhino.Node) {
	c.blockStack = append(c.blockStack, basicBlockEntry{root: root})
}

func (c *inferConstsCollector) popLastBlock(root rhino.Node) {
	idx := len(c.blockStack) - 1
	lastRoot := c.blockStack[idx].root
	c.blockStack[idx] = basicBlockEntry{}
	c.blockStack = c.blockStack[:idx]
	jutil.CheckStateMsg1(root == lastRoot, "%s", root)
}

func (c *inferConstsCollector) getCurrentBlock() *BasicBlock {
	idx := len(c.blockStack) - 1
	if b := c.blockStack[idx].block; b != nil {
		return b
	}
	return c.materializeBlock(idx)
}

func (c *inferConstsCollector) materializeBlock(idx int) *BasicBlock {
	i := idx
	for i > 0 && c.blockStack[i-1].block == nil {
		i--
	}
	var parent *BasicBlock
	if i > 0 {
		parent = c.blockStack[i-1].block
	}
	for ; i <= idx; i++ {
		parent = initBasicBlock(c.blocks.New(), parent, c.blockStack[i].root)
		c.blockStack[i].block = parent
	}
	return parent
}

func (c *inferConstsCollector) addCollectedHoistedFunction(n rhino.Node) bool {
	if c.collectedHoisted == nil {
		c.collectedHoisted = make(map[rhino.Node]struct{})
	}
	if _, exists := c.collectedHoisted[n]; exists {
		return false
	}
	c.collectedHoisted[n] = struct{}{}
	return true
}

func (c *inferConstsCollector) hasCollectedHoistedFunction(n rhino.Node) bool {
	_, exists := c.collectedHoisted[n]
	return exists
}

func (c *inferConstsCollector) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.GetFirstChild().IsNil() {
		switch n.GetToken() {
		case rhino.Token_NAME, rhino.Token_IMPORT_STAR, rhino.Token_BLOCK, rhino.Token_SWITCH_BODY:
			if referenceCollector_isBlockBoundary(n, parent) {
				c.pushNewBlock(n)
			}
			return true
		default:
			return false
		}
	}
	if n.IsFunction() && rcIsHoistedFunctionDeclaration(n) {
		if !c.addCollectedHoistedFunction(n) {
			return false
		}
	}
	if referenceCollector_isBlockBoundary(n, parent) {
		c.pushNewBlock(n)
	}
	return true
}

func (c *inferConstsCollector) EnterScope(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	if nodeTraversal_isHoistScopeRootNode(scopeRoot) {
		c.pushNewBlock(scopeRoot)
	}
}

func (c *inferConstsCollector) ExitScope(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	if nodeTraversal_isHoistScopeRootNode(scopeRoot) {
		c.popLastBlock(scopeRoot)
	}
	if len(t.scopes) == 2 && c.jumpDepth == 0 {
		c.flushLocalStates()
	}
}

func (c *inferConstsCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	tok := n.GetToken()
	if tok == rhino.Token_NAME || tok == rhino.Token_IMPORT_STAR {
		if (parent.IsImportSpec() && n != parent.GetLastChild()) ||
			(parent.IsExportSpec() && n != parent.GetFirstChild()) {
			return
		}
		id := n.GetStringID()
		if id > rhino.RhinoStringPool_EmptyID {
			v := t.LookupVarByID(id)
			if v != nil {
				c.addReference(v, n, parent, t)
				if v.nameNode != n {
					c.maybeJumpToHoistedFunction(v, t)
				}
			}
		}
	}

	if referenceCollector_isBlockBoundary(n, parent) {
		c.popLastBlock(n)
	}
}

func (c *inferConstsCollector) maybeJumpToHoistedFunction(var_ *Var, t *NodeTraversal) {
	nameNode := var_.nameNode
	if nameNode.IsNil() {
		return
	}
	fnNode := nameNode.GetParent()
	if fnNode.IsNil() || !fnNode.IsFunction() {
		return
	}
	if c.hasCollectedHoistedFunction(fnNode) || !rcIsHoistedFunctionDeclaration(fnNode) {
		return
	}
	varScope := var_.scope
	oldBlockStack := c.blockStack
	var newStack []basicBlockEntry
	if varScope.IsGlobal() {
		if len(oldBlockStack) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		if oldBlockStack[0].block == nil {
			c.materializeBlock(0)
		}
		newStack = []basicBlockEntry{oldBlockStack[0]}
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
		newStack = slices.Clone(oldBlockStack[:end])
	}
	c.blockStack = newStack

	firstChild := fnNode.GetFirstChild()
	c.addReference(var_, firstChild, fnNode, t)
	c.jumpDepth++
	b := c.builderVal()
	b.TraverseAtScope(c.scopeCreator.CreateScope(fnNode, varScope))
	c.jumpDepth--
	if len(t.scopes) == 1 && c.jumpDepth == 0 {
		c.flushLocalStates()
	}

	c.blockStack = oldBlockStack
}

func inferConsts_isLvalue(n rhino.Node, parent rhino.Node) bool {
	switch parent.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		return n.HasChildren() || (NodeUtil_IsEnhancedFor(parent) && parent.GetFirstChild() == n)
	case rhino.Token_DEFAULT_VALUE:
		return parent.GetFirstChild() == n
	case rhino.Token_INC, rhino.Token_DEC, rhino.Token_CATCH, rhino.Token_ITER_REST, rhino.Token_OBJECT_REST,
		rhino.Token_PARAM_LIST:
		return true
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		return NodeUtil_IsEnhancedFor(parent) && parent.GetFirstChild() == n
	case rhino.Token_ARRAY_PATTERN, rhino.Token_STRING_KEY, rhino.Token_COMPUTED_PROP:
		return NodeUtil_IsLhsByDestructuring(n)
	default:
		return NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n
	}
}

func (c *inferConstsCollector) addReference(v *Var, n rhino.Node, parent rhino.Node, t *NodeTraversal) {
	idx := v.nameHash
	if idx == 0 {
		isDecl := reference_isDeclarationHelper(n)
		parentTok := parent.GetToken()
		isInitDecl := (isDecl && parentTok != rhino.Token_VAR && parentTok != rhino.Token_LET) || n.HasChildren()
		isAssign := isInitDecl || inferConsts_isLvalue(n, parent)
		declType := v.DeclarationType()
		needOrderCheck := isInitDecl && declType != nil && *declType == rhino.Token_VAR
		var currentBlock *BasicBlock
		if isAssign || needOrderCheck {
			currentBlock = c.getCurrentBlock()
		}
		st := inferConstVarState{
			v:               v,
			firstBlock:      currentBlock,
			firstIsInitDecl: isInitDecl,
			firstIsDecl:     isDecl,
			allAfterFirst:   needOrderCheck,
		}
		if isAssign {
			st.assignCount = 1
			st.assignBlock = currentBlock
			_, st.assignHoistRoot = t.closestScopeHoistRoot()
		}
		if v.IsLocal() {
			v.nameHash = int32(len(c.localStates))
			c.localStates = appendDoubling(c.localStates, st)
		} else {
			v.nameHash = -int32(len(c.globalStates))
			c.globalStates = appendDoubling(c.globalStates, st)
		}
		return
	}

	var st *inferConstVarState
	if idx > 0 {
		st = &c.localStates[idx]
	} else {
		st = &c.globalStates[-idx]
	}
	if st.assignCount > 1 {
		return
	}
	parentTok := parent.GetToken()
	isInitDecl := n.HasChildren() || (parentTok != rhino.Token_VAR && parentTok != rhino.Token_LET && reference_isDeclarationHelper(n))
	isAssign := isInitDecl || inferConsts_isLvalue(n, parent)
	if !isAssign && !st.allAfterFirst {
		return
	}
	currentBlock := c.getCurrentBlock()
	if isAssign {
		st.assignCount++
		if st.assignCount == 1 {
			st.assignBlock = currentBlock
			_, st.assignHoistRoot = t.closestScopeHoistRoot()
		}
	}
	if st.assignCount <= 1 && st.allAfterFirst {
		if st.firstBlock != currentBlock && !st.firstBlock.ProvablyExecutesBefore(currentBlock) {
			st.allAfterFirst = false
		}
	}
}

// Java: InferConsts.considerVar(Var, ReferenceCollection)
func (c *InferConsts) ConsiderVar(v *Var, refCollection *ReferenceCollection) {
	if v.IsImplicitGoogNamespace() {
		return // no name node for provided variables.
	}
	nameNode := v.GetNameNode()
	if !nameNode.IsNil() && nameNode.GetParent().IsConst() {
		nameNode.SetDeclaredConstantVar(true)
	} else if docInfo := v.GetJSDocInfo(); docInfo != nil && docInfo.IsConstant() {
		nameNode.SetDeclaredConstantVar(true)
	}
	if InferConsts_IsInferredConst(v, refCollection) {
		nameNode.SetInferredConstantVar(true)
	}
}

// Java: InferConsts.isInferredConst(Var, ReferenceCollection)
func InferConsts_IsInferredConst(v *Var, refCollection *ReferenceCollection) bool {
	nameNode := v.GetNameNode()
	if nameNode.IsNil() || refCollection == nil || !refCollection.IsAssignedOnceInLifetime() {
		return false
	}
	if v.IsImplicitGoogNamespace() {
		return false
	}
	declarationType := v.DeclarationType()

	// A switch on a null Token throws NullPointerException.
	switch *declarationType {
	case rhino.Token_LET:
		// Check that non-destructuring let names are actually assigned at declaration.
		return !nameNode.GetParent().IsLet() || nameNode.HasChildren()
	case rhino.Token_CONST,
		rhino.Token_CATCH,
		rhino.Token_CLASS,
		rhino.Token_PARAM_LIST, // Parameters cannot be referenced before declaration.
		rhino.Token_FUNCTION:   // Function hoisting means no references before declaration.
		return true
	case rhino.Token_IMPORT:
		// ES module exports are mutable.
		// TODO(lharker): make this check smarter if we start optimizing unrewritten modules.
		return false
	case rhino.Token_VAR:
		// var hoisting requires this extra work to make sure the 'declaration' is also the first
		// reference.
		return refCollection.FirstReferenceIsAssigningDeclaration() && refCollection.IsWellDefined()
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unrecognized declaration type " + declarationType.String()))
	}
}
