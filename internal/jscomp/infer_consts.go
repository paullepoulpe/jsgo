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

	"github.com/paullepoulpe/jsgo/internal/ast"
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
func (c *InferConsts) Process(externs rhino.Node, js rhino.Node) {
	scopeCreator := NewSyntacticScopeCreator(c.compiler)
	collector := newInferConstsCollector(c, c.compiler, scopeCreator)
	collector.processRoot(js)

	for i := range collector.globalStates {
		if collector.globalStates[i].v != nil {
			c.considerVarState(&collector.globalStates[i])
		}
	}

	globalExternsScope := scopeCreator.CreateScopeOfScope(externs, nil)
	for v := range globalExternsScope.GetAllSymbols() {
		c.considerExternVar(v)
	}
}

func markDeclaredConst(v *Var) rhino.Node {
	if v.IsImplicitGoogNamespace() {
		return rhino.NilNode // no name node for provided variables.
	}
	nameNode := v.GetNameNode()
	if !nameNode.IsNil() && nameNode.GetParent().IsConst() {
		nameNode.SetDeclaredConstantVar(true)
	} else if docInfo := v.GetJSDocInfo(); docInfo != nil && docInfo.IsConstant() {
		nameNode.SetDeclaredConstantVar(true)
	}
	return nameNode
}

func (c *InferConsts) considerExternVar(v *Var) {
	markDeclaredConst(v)
}

func (c *InferConsts) considerVarState(st *inferConstVarState) {
	nameNode := markDeclaredConst(st.v)
	if !nameNode.IsNil() && inferConsts_isInferredConstState(st.v, nameNode, st) {
		nameNode.SetInferredConstantVar(true)
	}
}

func inferConsts_isInferredConstState(v *Var, nameNode rhino.Node, st *inferConstVarState) bool {
	if nameNode.IsNil() || st.assignCount != 1 || !st.assignedOnceInLifetime {
		return false
	}
	declarationType, ok := v.DeclarationToken()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	switch declarationType {
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
	v                      *Var
	firstBlock             *BasicBlock
	assignCount            int32
	firstIsInitDecl        bool
	firstIsDecl            bool
	allAfterFirst          bool
	assignedOnceInLifetime bool
}

type inferConstScopeFrame struct {
	root  rhino.Node
	base  int
	count int
}

type inferConstsCollector struct {
	outer        *InferConsts
	compiler     AbstractCompiler
	scopeCreator ScopeCreator
	localStates  []inferConstVarState
	globalStates []inferConstVarState
	scopeFrames  []inferConstScopeFrame
	blocks       basicBlockTracker
}

var (
	_ NodeTraversal_ScopedCallback = (*inferConstsCollector)(nil)
	_ nodeTraversal_visitTokens    = (*inferConstsCollector)(nil)
)

func (*inferConstsCollector) visitTokens() *ast.TokenSet {
	return &refCollectorVisitTokens
}

func newInferConstsCollector(outer *InferConsts, compiler AbstractCompiler, sc *SyntacticScopeCreator) *inferConstsCollector {
	return &inferConstsCollector{
		outer:        outer,
		compiler:     compiler,
		scopeCreator: sc,
	}
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

func isLoopBlockRoot(root rhino.Node) bool {
	p := root.GetParent()
	if p.IsNil() {
		return false
	}
	switch p.GetToken() {
	case rhino.Token_DO, rhino.Token_WHILE, rhino.Token_FOR,
		rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_FOR_IN:
		return true
	default:
		return false
	}
}

func (c *inferConstsCollector) isCurrentBlockSingleExec(varScope *Scope, t *NodeTraversal) bool {
	entries := c.blocks.entries
	base := c.blocks.frameBase
	if len(entries) <= base {
		panic(jutil.NewNoSuchElementException())
	}
	varScopeRoot := varScope.GetRootNode()
	for i := len(entries) - 1; i >= base; i-- {
		root := entries[i].root
		if root.IsFunction() {
			_, assignHoistRoot := t.closestScopeHoistRoot()
			return varScope.GetClosestHoistScope().GetRootNode() == assignHoistRoot
		}
		isLoop := isLoopBlockRoot(root)
		if i == base && !isLoop {
			return true
		}
		if root == varScopeRoot {
			return true
		}
		if isLoop {
			return false
		}
	}
	return true
}

func (c *inferConstsCollector) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.GetFirstChild().IsNil() {
		switch n.GetToken() {
		case rhino.Token_NAME, rhino.Token_IMPORT_STAR, rhino.Token_BLOCK, rhino.Token_SWITCH_BODY:
			if referenceCollector_isBlockBoundary(n, parent) {
				c.blocks.push(n)
			}
			return true
		default:
			return false
		}
	}
	if n.IsFunction() && rcIsHoistedFunctionDeclaration(n) {
		if !c.blocks.markHoistedCollected(n) {
			return false
		}
	}
	if referenceCollector_isBlockBoundary(n, parent) {
		c.blocks.push(n)
	}
	return true
}

func (c *inferConstsCollector) EnterScope(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	if nodeTraversal_isHoistScopeRootNode(scopeRoot) {
		c.blocks.push(scopeRoot)
	}
	if !t.InGlobalScope() {
		c.scopeFrames = append(c.scopeFrames, inferConstScopeFrame{root: scopeRoot, base: -1})
	}
}

func (c *inferConstsCollector) ExitScope(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	if nodeTraversal_isHoistScopeRootNode(scopeRoot) {
		c.blocks.pop(scopeRoot)
	}
	if !t.InGlobalScope() {
		top := len(c.scopeFrames) - 1
		f := c.scopeFrames[top]
		c.scopeFrames = c.scopeFrames[:top]
		if f.base >= 0 {
			end := f.base + f.count
			for i := f.base; i < end; i++ {
				if c.localStates[i].v != nil {
					c.outer.considerVarState(&c.localStates[i])
				}
			}
			clear(c.localStates[f.base:end])
			if end == len(c.localStates) {
				c.localStates = c.localStates[:f.base]
			}
		}
		if len(c.scopeFrames) == 0 && len(c.localStates) > 0 {
			clear(c.localStates)
			c.localStates = c.localStates[:0]
		}
	}
}

func (c *inferConstsCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	tok := n.GetToken()
	if tok == rhino.Token_NAME || tok == rhino.Token_IMPORT_STAR {
		switch ast.RoleIn(n, parent) {
		case ast.RoleImportSpecImported, ast.RoleExportSpecExported:
			return
		}
		id := n.GetStringID()
		if id > rhino.RhinoStringPool_EmptyID {
			v := t.LookupVarByID(id)
			if v != nil && !v.nameNode.IsNil() {
				c.addReference(v, n, parent, t)
				if v.nameNode != n {
					c.maybeJumpToHoistedFunction(v, t)
				}
			}
		}
	}

	if referenceCollector_isBlockBoundary(n, parent) {
		c.blocks.pop(n)
	}
}

func (c *inferConstsCollector) maybeJumpToHoistedFunction(var_ *Var, t *NodeTraversal) {
	fnNode, ok := c.blocks.uncollectedHoistedTarget(var_)
	if !ok {
		return
	}
	varScope := var_.scope
	saved := c.blocks.saveHoistedFrame(varScope)

	c.addReference(var_, fnNode.GetFirstChild(), fnNode, t)
	b := c.builderVal()
	b.TraverseAtScope(c.scopeCreator.CreateScope(fnNode, varScope))

	c.blocks.restoreFrame(saved)
}

func inferConsts_isLvalue(n rhino.Node, parent rhino.Node) bool {
	switch ast.RoleIn(n, parent) {
	case ast.RoleNameDeclarationDeclarators:
		return n.HasChildren()
	case ast.RoleDefaultValueTarget,
		ast.RoleIncDecTarget,
		ast.RoleCatchParam,
		ast.RoleIterRestTarget,
		ast.RoleObjectRestTarget,
		ast.RoleParamListParams,
		ast.RoleForInTarget,
		ast.RoleForOfTarget,
		ast.RoleForAwaitOfTarget,
		ast.RoleAssignTarget,
		ast.RoleAssignOpTarget:
		return true
	case ast.RoleArrayPatternElements,
		ast.RoleStringKeyValue,
		ast.RoleComputedPropKey,
		ast.RoleComputedPropValue:
		return NodeUtil_IsLhsByDestructuring(n)
	default:
		return false
	}
}

func (c *inferConstsCollector) stateForVar(v *Var) *inferConstVarState {
	slot := scopeVarSlot(v.scope, v)
	if slot < 0 {
		return nil
	}
	if !v.IsLocal() {
		if slot >= len(c.globalStates) {
			oldLen := len(c.globalStates)
			n := max(len(v.scope.varList), slot+1)
			c.globalStates = slices.Grow(c.globalStates, n-oldLen)[:n]
			clear(c.globalStates[oldLen:n])
		}
		return &c.globalStates[slot]
	}
	targetRoot := v.scope.rootNode
	for fi := len(c.scopeFrames) - 1; fi >= 0; fi-- {
		f := &c.scopeFrames[fi]
		if f.root == targetRoot {
			if f.base < 0 {
				count := max(len(v.scope.varList), slot+1)
				base := len(c.localStates)
				f.base = base
				f.count = count
				c.localStates = slices.Grow(c.localStates, count)[:base+count]
				clear(c.localStates[base : base+count])
			} else if slot >= f.count {
				return nil
			}
			return &c.localStates[f.base+slot]
		}
	}
	return nil
}

func (c *inferConstsCollector) addReference(v *Var, n rhino.Node, parent rhino.Node, t *NodeTraversal) {
	st := c.stateForVar(v)
	if st == nil {
		return
	}
	if st.v == nil {
		isDecl := reference_isDeclarationHelper(n)
		parentTok := parent.GetToken()
		isInitDecl := (isDecl && parentTok != rhino.Token_VAR && parentTok != rhino.Token_LET) || n.HasChildren()
		isAssign := isInitDecl || inferConsts_isLvalue(n, parent)
		declType, hasDeclType := v.DeclarationToken()
		needOrderCheck := isInitDecl && hasDeclType && declType == rhino.Token_VAR
		var firstBlock *BasicBlock
		if needOrderCheck {
			firstBlock = c.blocks.currentBlock()
		}
		*st = inferConstVarState{
			v:               v,
			firstBlock:      firstBlock,
			firstIsInitDecl: isInitDecl,
			firstIsDecl:     isDecl,
			allAfterFirst:   needOrderCheck,
		}
		if isAssign {
			st.assignCount = 1
			st.assignedOnceInLifetime = c.isCurrentBlockSingleExec(v.scope, t)
		}
		return
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
	if isAssign {
		st.assignCount++
		if st.assignCount == 1 {
			st.assignedOnceInLifetime = c.isCurrentBlockSingleExec(v.scope, t)
		}
	}
	if st.assignCount <= 1 && st.allAfterFirst {
		currentBlock := c.blocks.currentBlock()
		if st.firstBlock != currentBlock && !st.firstBlock.ProvablyExecutesBefore(currentBlock) {
			st.allAfterFirst = false
		}
	}
}
