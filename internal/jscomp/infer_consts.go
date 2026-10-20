// Port of com/google/javascript/jscomp/InferConsts.java @ closure-compiler v20261015

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

// Java: InferConsts.process(Node, Node)
func (c *InferConsts) Process(externs *rhino.Node, js *rhino.Node) {
	scopeCreator := NewSyntacticScopeCreator(c.compiler)
	collector := newInferConstsCollector(c.compiler, scopeCreator)
	collector.processRoot(js)

	for i := range collector.states {
		c.considerVarState(&collector.states[i])
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
	if nameNode != nil && nameNode.GetParent().IsConst() {
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
	if nameNode != nil && nameNode.GetParent().IsConst() {
		nameNode.SetDeclaredConstantVar(true)
	} else if docInfo := v.GetJSDocInfo(); docInfo != nil && docInfo.IsConstant() {
		nameNode.SetDeclaredConstantVar(true)
	}
	if inferConsts_isInferredConstState(v, nameNode, st) {
		nameNode.SetInferredConstantVar(true)
	}
}

func inferConsts_isInferredConstState(v *Var, nameNode *rhino.Node, st *inferConstVarState) bool {
	if nameNode == nil || st.assignCount != 1 || !st.isAssignedOnceInLifetime() {
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
	assignScope     *Scope
	assignCount     int32
	firstIsInitDecl bool
	firstIsDecl     bool
	allAfterFirst   bool
}

func (st *inferConstVarState) isAssignedOnceInLifetime() bool {
	varScope := st.v.scope
	for block := st.assignBlock; block != nil; block = block.parent {
		if block.isFunction {
			if varScope.GetClosestHoistScope() != st.assignScope.GetClosestHoistScope() {
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
	compiler            AbstractCompiler
	scopeCreator        ScopeCreator
	states              []inferConstVarState
	blockStack          []*BasicBlock
	blockSlab           []BasicBlock
	collectedHoistedMap map[*rhino.Node]struct{}
	blockStackBuf       [16]*BasicBlock
	blockBuf            [32]BasicBlock
	statesBuf           [64]inferConstVarState
}

func newInferConstsCollector(compiler AbstractCompiler, scopeCreator ScopeCreator) *inferConstsCollector {
	c := &inferConstsCollector{
		compiler:     compiler,
		scopeCreator: scopeCreator,
	}
	c.blockStack = c.blockStackBuf[:0]
	c.blockSlab = c.blockBuf[:0]
	c.states = c.statesBuf[:0]
	return c
}

func (c *inferConstsCollector) processRoot(root *rhino.Node) {
	c.createTraversalBuilder().Traverse(root)
}

func (c *inferConstsCollector) createTraversalBuilder() *NodeTraversal_Builder {
	return NodeTraversal_Builder_().
		SetCompiler(c.compiler).
		SetCallback(c).
		SetScopeCreator(c.scopeCreator).
		SetObeyDestructuringAndDefaultValueExecutionOrder(true)
}

func (c *inferConstsCollector) allocBasicBlock(parent *BasicBlock, root *rhino.Node) *BasicBlock {
	if len(c.blockSlab) == cap(c.blockSlab) {
		newCap := 64
		if cap(c.blockSlab) >= 64 {
			newCap = min(cap(c.blockSlab)*2, 1024)
		}
		c.blockSlab = make([]BasicBlock, 0, newCap)
	}
	idx := len(c.blockSlab)
	c.blockSlab = c.blockSlab[:idx+1]
	return initBasicBlock(&c.blockSlab[idx], parent, root)
}

func (c *inferConstsCollector) pushNewBlock(root *rhino.Node) {
	var peekLast *BasicBlock
	if len(c.blockStack) > 0 {
		peekLast = c.blockStack[len(c.blockStack)-1]
	}
	c.blockStack = append(c.blockStack, c.allocBasicBlock(peekLast, root))
}

func (c *inferConstsCollector) popLastBlock(root *rhino.Node) {
	last := c.blockStack[len(c.blockStack)-1]
	c.blockStack[len(c.blockStack)-1] = nil
	c.blockStack = c.blockStack[:len(c.blockStack)-1]
	jutil.CheckStateMsg1(root == last.GetRoot(), "%s", root)
}

func (c *inferConstsCollector) addCollectedHoistedFunction(n *rhino.Node) bool {
	if c.collectedHoistedMap == nil {
		c.collectedHoistedMap = make(map[*rhino.Node]struct{}, 8)
	}
	if _, exists := c.collectedHoistedMap[n]; exists {
		return false
	}
	c.collectedHoistedMap[n] = struct{}{}
	return true
}

func (c *inferConstsCollector) hasCollectedHoistedFunction(n *rhino.Node) bool {
	if c.collectedHoistedMap == nil {
		return false
	}
	_, exists := c.collectedHoistedMap[n]
	return exists
}

func (c *inferConstsCollector) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsFunction() && NodeUtil_IsHoistedFunctionDeclaration(n) {
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
	if t.IsHoistScope() {
		c.pushNewBlock(t.GetScopeRoot())
	}
}

func (c *inferConstsCollector) ExitScope(t *NodeTraversal) {
	if t.IsHoistScope() {
		c.popLastBlock(t.GetScopeRoot())
	}
}

func (c *inferConstsCollector) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	tok := n.GetToken()
	if tok == rhino.Token_NAME || tok == rhino.Token_IMPORT_STAR {
		if (parent.IsImportSpec() && n != parent.GetLastChild()) ||
			(parent.IsExportSpec() && n != parent.GetFirstChild()) {
			return
		}
		name := n.GetString()
		if name != "" {
			scope := t.GetScope()
			v := scope.GetVar(name)
			if v != nil {
				c.addReference(v, n, parent, scope)
				c.maybeJumpToHoistedFunction(v, t)
			}
		}
	}

	if referenceCollector_isBlockBoundary(n, parent) {
		c.popLastBlock(n)
	}
}

func (c *inferConstsCollector) maybeJumpToHoistedFunction(var_ *Var, t *NodeTraversal) {
	fnNode := var_.GetParentNode()
	if fnNode == nil || !fnNode.IsFunction() || !NodeUtil_IsHoistedFunctionDeclaration(fnNode) {
		return
	}
	if c.hasCollectedHoistedFunction(fnNode) {
		return
	}
	varScope := var_.scope
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

	firstChild := fnNode.GetFirstChild()
	c.addReference(var_, firstChild, fnNode, t.GetScope())
	c.createTraversalBuilder().TraverseAtScope(c.scopeCreator.CreateScope(fnNode, varScope))

	c.blockStack = oldBlockStack
}

func inferConsts_isLvalue(n *rhino.Node, parent *rhino.Node) bool {
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

func (c *inferConstsCollector) addReference(v *Var, n *rhino.Node, parent *rhino.Node, scope *Scope) {
	currentBlock := c.blockStack[len(c.blockStack)-1]
	idx := v.nameHash
	if idx == 0 {
		idx = int32(len(c.states) + 1)
		v.nameHash = idx
		isDecl := reference_isDeclarationHelper(n)
		parentTok := parent.GetToken()
		isInitDecl := (isDecl && parentTok != rhino.Token_VAR && parentTok != rhino.Token_LET) || n.HasChildren()
		isAssign := isInitDecl || inferConsts_isLvalue(n, parent)
		st := inferConstVarState{
			v:               v,
			firstBlock:      currentBlock,
			firstIsInitDecl: isInitDecl,
			firstIsDecl:     isDecl,
			allAfterFirst:   true,
		}
		if isAssign {
			st.assignCount = 1
			st.assignBlock = currentBlock
			st.assignScope = scope
		}
		c.states = append(c.states, st)
		return
	}

	st := &c.states[idx-1]
	if st.assignCount > 1 {
		return
	}
	parentTok := parent.GetToken()
	isInitDecl := n.HasChildren() || (parentTok != rhino.Token_VAR && parentTok != rhino.Token_LET && reference_isDeclarationHelper(n))
	if isInitDecl || inferConsts_isLvalue(n, parent) {
		st.assignCount++
		if st.assignCount == 1 {
			st.assignBlock = currentBlock
			st.assignScope = scope
		}
	}
	if st.assignCount <= 1 && st.firstIsInitDecl && st.allAfterFirst {
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
	if nameNode != nil && nameNode.GetParent().IsConst() {
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
	if nameNode == nil || refCollection == nil || !refCollection.IsAssignedOnceInLifetime() {
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
