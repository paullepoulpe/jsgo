// Port of com/google/javascript/jscomp/InferConsts.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
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
	if !c.processParallel(js, scopeCreator) {
		collector := newInferConstsCollector(c, c.compiler, scopeCreator)
		collector.processRoot(js)

		for i := range collector.globalStates {
			if collector.globalStates[i].v != nil {
				c.considerVarState(&collector.globalStates[i])
			}
		}
	}

	globalExternsScope := scopeCreator.CreateScopeOfScope(externs, nil)
	for v := range globalExternsScope.GetAllSymbols() {
		c.considerExternVar(v)
	}
}

const (
	inferConstFlagDeclared uint8 = 1 << iota
	inferConstFlagInferred
)

type inferConstMark struct {
	node  rhino.Node
	flags uint8
}

func computeConstFlags(v *Var, st *inferConstVarState) (rhino.Node, uint8) {
	if v.IsImplicitGoogNamespace() {
		return rhino.NilNode, 0 // no name node for provided variables.
	}
	nameNode := v.GetNameNode()
	if nameNode.IsNil() {
		return rhino.NilNode, 0
	}
	var flags uint8
	if nameNode.GetParent().IsConst() {
		flags |= inferConstFlagDeclared
	} else if docInfo := v.GetJSDocInfo(); docInfo != nil && docInfo.IsConstant() {
		flags |= inferConstFlagDeclared
	}
	if st != nil && inferConsts_isInferredConstState(v, nameNode, st) {
		flags |= inferConstFlagInferred
	}
	return nameNode, flags
}

func applyConstFlags(nameNode rhino.Node, flags uint8) {
	if flags&inferConstFlagDeclared != 0 {
		nameNode.SetDeclaredConstantVar(true)
	}
	if flags&inferConstFlagInferred != 0 {
		nameNode.SetInferredConstantVar(true)
	}
}

func (c *InferConsts) considerExternVar(v *Var) {
	nameNode, flags := computeConstFlags(v, nil)
	applyConstFlags(nameNode, flags)
}

func (c *InferConsts) considerVarState(st *inferConstVarState) {
	nameNode, flags := computeConstFlags(st.v, st)
	applyConstFlags(nameNode, flags)
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
	scope *Scope
	base  int
	count int
}

type inferConstsCollector struct {
	outer         *InferConsts
	compiler      AbstractCompiler
	scopeCreator  ScopeCreator
	localStates   []inferConstVarState
	globalStates  []inferConstVarState
	scopeFrames   []inferConstScopeFrame
	blocks        basicBlockTracker
	parallelState *inferParallelScriptState
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
		if ps := c.parallelState; ps != nil {
			if ps.recordHoistedFunctionDecl(c, n, parent) {
				return false
			}
			if n == ps.inHoistedRoot {
				return true
			}
		}
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
	inGlobal := t.InGlobalScope()
	scopeRoot := t.GetScopeRoot()
	if nodeTraversal_isHoistScopeRootNode(scopeRoot) {
		c.blocks.push(scopeRoot)
		if inGlobal && c.parallelState != nil {
			c.blocks.setTopBlock(c.parallelState.globalBlock)
		}
	}
	if !inGlobal {
		var s *Scope
		if len(t.scopes) > 0 {
			s = t.scopes[len(t.scopes)-1].untypedScope
		}
		if s != nil {
			s.inferFrameIdx = int32(len(c.scopeFrames))
		}
		c.scopeFrames = append(c.scopeFrames, inferConstScopeFrame{root: scopeRoot, scope: s, base: -1})
	}
}

func (c *inferConstsCollector) ExitScope(t *NodeTraversal) {
	inGlobal := t.InGlobalScope()
	scopeRoot := t.GetScopeRoot()
	if nodeTraversal_isHoistScopeRootNode(scopeRoot) {
		c.blocks.pop(scopeRoot)
	}
	if !inGlobal {
		top := len(c.scopeFrames) - 1
		f := c.scopeFrames[top]
		c.scopeFrames = c.scopeFrames[:top]
		if f.scope != nil {
			f.scope.inferFrameIdx = -1
		} else if len(t.scopes) > 0 {
			if s := t.scopes[len(t.scopes)-1].untypedScope; s != nil {
				s.inferFrameIdx = -1
			}
		}
		if f.base >= 0 {
			end := f.base + f.count
			if ps := c.parallelState; ps != nil {
				for i := f.base; i < end; i++ {
					if st := &c.localStates[i]; st.v != nil {
						if node, flags := computeConstFlags(st.v, st); flags != 0 {
							ps.localMarks = append(ps.localMarks, inferConstMark{node: node, flags: flags})
						}
					}
				}
			} else {
				for i := f.base; i < end; i++ {
					if c.localStates[i].v != nil {
						c.outer.considerVarState(&c.localStates[i])
					}
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
		if id > strpool.EmptyID {
			v := t.LookupVarByID(id)
			if v != nil && !v.nameNode.IsNil() {
				if ps := c.parallelState; ps != nil && !v.IsLocal() {
					ps.recordGlobalRef(c, v, n, parent, t)
					if v.nameNode != n {
						ps.maybeRecordJumpHoisted(v)
					}
				} else {
					c.addReference(v, n, parent, t)
					if v.nameNode != n {
						c.maybeJumpToHoistedFunction(v, t)
					}
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
	s := v.scope
	fi := int(s.inferFrameIdx)
	if fi < 0 {
		targetRoot := s.rootNode
		for i := len(c.scopeFrames) - 1; i >= 0; i-- {
			if c.scopeFrames[i].root == targetRoot {
				fi = i
				s.inferFrameIdx = int32(i)
				c.scopeFrames[i].scope = s
				break
			}
		}
		if fi < 0 {
			return nil
		}
	} else if uint(fi) >= uint(len(c.scopeFrames)) || c.scopeFrames[fi].root != s.rootNode {
		return nil
	}
	f := &c.scopeFrames[fi]
	if f.base < 0 {
		count := max(len(s.varList), slot+1)
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

type inferGlobalOpKind uint8

const (
	inferOpRef inferGlobalOpKind = iota + 1
	inferOpJumpHoisted
	inferOpDeclHoisted
)

const (
	inferRefIsDecl uint8 = 1 << iota
	inferRefIsInitDecl
	inferRefIsAssign
	inferRefSingleExec
)

type inferGlobalOp struct {
	v      *Var
	block  *BasicBlock
	fnNode rhino.Node
	kind   inferGlobalOpKind
	flags  uint8
}

type inferParallelSegment struct {
	fnNode rhino.Node
	ops    []inferGlobalOp
}

type inferParallelScriptState struct {
	globalScope   *Scope
	globalBlock   *BasicBlock
	curSeg        *inferParallelSegment
	rootSeg       inferParallelSegment
	hoistedSegs   []*inferParallelSegment
	localMarks    []inferConstMark
	inHoistedRoot rhino.Node
}

type inferParallelWorker struct {
	sc        *SyntacticScopeCreator
	collector inferConstsCollector
}

func (ps *inferParallelScriptState) recordHoistedFunctionDecl(
	c *inferConstsCollector,
	n rhino.Node,
	parent rhino.Node,
) bool {
	if n == ps.inHoistedRoot || parent.IsNil() {
		return false
	}
	gv := ps.globalScope.lookupOwnVarByID(n.GetFirstChild().GetStringID())
	if gv == nil || gv.nameNode != n.GetFirstChild() {
		return false
	}
	ps.curSeg.ops = append(ps.curSeg.ops, inferGlobalOp{
		v:      gv,
		fnNode: n,
		kind:   inferOpDeclHoisted,
	})
	hSeg := &inferParallelSegment{fnNode: n}
	ps.hoistedSegs = append(ps.hoistedSegs, hSeg)

	prevSeg := ps.curSeg
	prevHoisted := ps.inHoistedRoot
	savedBlocks := c.blocks.pushRootBlockFrame(ps.globalBlock)

	ps.curSeg = hSeg
	ps.inHoistedRoot = n

	b := c.builderVal()
	b.TraverseAtScope(c.scopeCreator.CreateScope(n, ps.globalScope))

	c.blocks.restoreFrame(savedBlocks)
	ps.inHoistedRoot = prevHoisted
	ps.curSeg = prevSeg
	return true
}

func (ps *inferParallelScriptState) maybeRecordJumpHoisted(v *Var) {
	if v.nameNode.IsNil() {
		return
	}
	fnNode := v.nameNode.GetParent()
	if fnNode.IsNil() || !fnNode.IsFunction() || !rcIsHoistedFunctionDeclaration(fnNode) {
		return
	}
	ps.curSeg.ops = append(ps.curSeg.ops, inferGlobalOp{
		v:      v,
		fnNode: fnNode,
		kind:   inferOpJumpHoisted,
	})
}

func (ps *inferParallelScriptState) recordGlobalRef(
	c *inferConstsCollector,
	v *Var,
	n rhino.Node,
	parent rhino.Node,
	t *NodeTraversal,
) {
	isDecl := reference_isDeclarationHelper(n)
	parentTok := parent.GetToken()
	isInitDecl := (isDecl && parentTok != rhino.Token_VAR && parentTok != rhino.Token_LET) || n.HasChildren()
	isAssign := isInitDecl || inferConsts_isLvalue(n, parent)
	var flags uint8
	if isDecl {
		flags |= inferRefIsDecl
	}
	if isInitDecl {
		flags |= inferRefIsInitDecl
	}
	if isAssign {
		flags |= inferRefIsAssign
		if c.isCurrentBlockSingleExec(v.scope, t) {
			flags |= inferRefSingleExec
		}
	}
	var block *BasicBlock
	if declType, ok := v.DeclarationToken(); ok && declType == rhino.Token_VAR {
		block = c.blocks.currentBlock()
	}
	ps.curSeg.ops = append(ps.curSeg.ops, inferGlobalOp{
		v:     v,
		block: block,
		kind:  inferOpRef,
		flags: flags,
	})
}

func applyGlobalRefOp(globalStates *[]inferConstVarState, v *Var, block *BasicBlock, flags uint8) {
	slot := scopeVarSlot(v.scope, v)
	if slot < 0 {
		return
	}
	if slot >= len(*globalStates) {
		oldLen := len(*globalStates)
		n := max(len(v.scope.varList), slot+1)
		*globalStates = slices.Grow(*globalStates, n-oldLen)[:n]
		clear((*globalStates)[oldLen:n])
	}
	st := &(*globalStates)[slot]
	if st.v == nil {
		isDecl := flags&inferRefIsDecl != 0
		isInitDecl := flags&inferRefIsInitDecl != 0
		isAssign := flags&inferRefIsAssign != 0
		needOrderCheck := isInitDecl && block != nil
		var firstBlock *BasicBlock
		if needOrderCheck {
			firstBlock = block
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
			st.assignedOnceInLifetime = flags&inferRefSingleExec != 0
		}
		return
	}

	if st.assignCount > 1 {
		return
	}
	isAssign := flags&inferRefIsAssign != 0
	if !isAssign && !st.allAfterFirst {
		return
	}
	if isAssign {
		st.assignCount++
		if st.assignCount == 1 {
			st.assignedOnceInLifetime = flags&inferRefSingleExec != 0
		}
	}
	if st.assignCount <= 1 && st.allAfterFirst {
		if st.firstBlock != block && !st.firstBlock.ProvablyExecutesBefore(block) {
			st.allAfterFirst = false
		}
	}
}

func (c *InferConsts) processParallel(js rhino.Node, mainSC *SyntacticScopeCreator) bool {
	scripts, ok := collectRootScripts(js)
	if !ok || len(scripts) < 2 {
		return false
	}
	numWorkers := parallelThreadCount(c.compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(c.compiler)
	globalScope := mainSC.CreateScopeOfScope(js, nil)
	freezeGlobalScopeForParallelReads(globalScope)
	lazyGlobal := &lazyGlobalScope{
		compiler: c.compiler,
		root:     js,
		scope:    globalScope,
	}
	lazyGlobal.once.Do(func() {})
	globalBlock := NewBasicBlock(nil, js)

	states := make([]inferParallelScriptState, len(scripts))
	workers := make([]inferParallelWorker, numWorkers)
	runParallelWork(c.compiler, numWorkers, len(scripts), func(workerID int) (func(int, *recordingCompiler), func()) {
		w := &workers[workerID]
		runItem := func(idx int, rc *recordingCompiler) {
			if w.sc == nil || w.sc.compiler != rc {
				w.sc = NewSyntacticScopeCreator(rc)
			}
			col := &w.collector
			col.outer = c
			col.compiler = rc
			col.scopeCreator = w.sc
			col.scopeFrames = col.scopeFrames[:0]
			col.blocks.entries = col.blocks.entries[:0]
			col.blocks.frameBase = 0
			col.blocks.resetHoistedCollected()

			ps := &states[idx]
			ps.globalScope = globalScope
			ps.globalBlock = globalBlock
			ps.curSeg = &ps.rootSeg
			col.parallelState = ps

			tb := col.builderVal()
			tb.traverseScriptUnderRoot(js, js, scripts[idx], lazyGlobal)
			col.parallelState = nil
		}
		return runItem, nil
	}, false)

	var hoistedByFn map[rhino.Node]*inferParallelSegment
	for i := range states {
		ps := &states[i]
		for _, m := range ps.localMarks {
			applyConstFlags(m.node, m.flags)
		}
		if len(ps.hoistedSegs) > 0 {
			if hoistedByFn == nil {
				hoistedByFn = make(map[rhino.Node]*inferParallelSegment, len(ps.hoistedSegs)*2)
			}
			for _, hSeg := range ps.hoistedSegs {
				hoistedByFn[hSeg.fnNode] = hSeg
			}
		}
	}

	var globalStates []inferConstVarState
	var collectedHoisted map[rhino.Node]struct{}
	if len(hoistedByFn) > 0 {
		collectedHoisted = make(map[rhino.Node]struct{}, len(hoistedByFn))
	}

	const hoistedDeclFlags = inferRefIsDecl | inferRefIsInitDecl | inferRefIsAssign | inferRefSingleExec
	var replaySeg func(seg *inferParallelSegment)
	replaySeg = func(seg *inferParallelSegment) {
		if seg == nil {
			return
		}
		for i := range seg.ops {
			op := &seg.ops[i]
			switch op.kind {
			case inferOpRef:
				applyGlobalRefOp(&globalStates, op.v, op.block, op.flags)
			case inferOpJumpHoisted, inferOpDeclHoisted:
				if _, done := collectedHoisted[op.fnNode]; done {
					continue
				}
				if collectedHoisted == nil {
					collectedHoisted = make(map[rhino.Node]struct{})
				}
				collectedHoisted[op.fnNode] = struct{}{}
				applyGlobalRefOp(&globalStates, op.v, nil, hoistedDeclFlags)
				if hoistedByFn != nil {
					replaySeg(hoistedByFn[op.fnNode])
				}
			}
		}
		seg.ops = nil
	}

	for i := range states {
		replaySeg(&states[i].rootSeg)
	}
	for i := range globalStates {
		if globalStates[i].v != nil {
			c.considerVarState(&globalStates[i])
		}
	}
	return true
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
