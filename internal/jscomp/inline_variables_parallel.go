// Port of com/google/javascript/jscomp/InlineVariables.java @ closure-compiler v20261021

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

import "github.com/paullepoulpe/jsgo/internal/rhino"

type rcStreamOpKind uint8

const (
	rcOpCompilerEvents rcStreamOpKind = iota + 1
	rcOpGlobalRef
	rcOpGlobalAliasRetry
	rcOpJumpHoisted
	rcOpDeclHoisted
)

type rcStreamOp struct {
	kind    rcStreamOpKind
	events  []recordedCompilerEvent
	v       *Var
	node    rhino.Node
	block   *BasicBlock
	scope   *Scope
	inputId *rhino.InputId
	retry   *InlineVariables_InliningBehavior_AliasInlineRetryHandler
}

type rcStreamSegment struct {
	fnNode rhino.Node
	ops    []rcStreamOp
}

type rcParallelScriptState struct {
	globalScope   *Scope
	globalBlock   *BasicBlock
	rc            *recordingCompiler
	curSeg        *rcStreamSegment
	rootSeg       rcStreamSegment
	hoistedSegs   []*rcStreamSegment
	inHoistedRoot rhino.Node
}

func (ps *rcParallelScriptState) flushCompilerEvents() {
	if ps.rc != nil && len(ps.rc.events) > 0 {
		ps.curSeg.ops = append(ps.curSeg.ops, rcStreamOp{
			kind:   rcOpCompilerEvents,
			events: ps.rc.takeEvents(),
		})
	}
}

func (ps *rcParallelScriptState) recordHoistedFunctionDecl(
	c *ReferenceCollector,
	t *NodeTraversal,
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
	ps.flushCompilerEvents()
	declInputId := t.inputId
	if t.compilerInput == nil {
		declInputId = t.GetInput().GetInputId()
	}
	ps.curSeg.ops = append(ps.curSeg.ops, rcStreamOp{
		kind:    rcOpDeclHoisted,
		v:       gv,
		node:    n,
		scope:   ps.globalScope,
		inputId: declInputId,
	})
	hSeg := &rcStreamSegment{fnNode: n}
	ps.hoistedSegs = append(ps.hoistedSegs, hSeg)

	prevSeg := ps.curSeg
	prevHoisted := ps.inHoistedRoot
	savedBlocks := c.blocks.pushRootBlockFrame(ps.globalBlock)
	oldLocalDepth := c.localScopeDepth

	ps.curSeg = hSeg
	ps.inHoistedRoot = n
	c.localScopeDepth = 0

	b := c.builderVal()
	b.TraverseAtScope(c.scopeCreator.CreateScope(n, ps.globalScope))
	ps.flushCompilerEvents()

	c.blocks.restoreFrame(savedBlocks)
	c.localScopeDepth = oldLocalDepth
	ps.inHoistedRoot = prevHoisted
	ps.curSeg = prevSeg
	return true
}

func (ps *rcParallelScriptState) maybeRecordJumpHoisted(v *Var, t *NodeTraversal) {
	if v.nameNode.IsNil() {
		return
	}
	fnNode := v.nameNode.GetParent()
	if fnNode.IsNil() || !fnNode.IsFunction() || !rcIsHoistedFunctionDeclaration(fnNode) {
		return
	}
	ps.flushCompilerEvents()
	inputId := t.inputId
	if t.compilerInput == nil {
		inputId = t.GetInput().GetInputId()
	}
	ps.curSeg.ops = append(ps.curSeg.ops, rcStreamOp{
		kind:    rcOpJumpHoisted,
		v:       v,
		node:    fnNode,
		scope:   t.GetScope(),
		inputId: inputId,
	})
}

func (ps *rcParallelScriptState) recordGlobalRef(
	c *ReferenceCollector,
	v *Var,
	n rhino.Node,
	t *NodeTraversal,
	scope *Scope,
) {
	block := c.blocks.currentBlock()
	if scope == nil {
		scope = t.GetScope()
	}
	inputId := t.inputId
	if t.compilerInput == nil {
		inputId = t.GetInput().GetInputId()
	}
	ps.flushCompilerEvents()
	ps.curSeg.ops = append(ps.curSeg.ops, rcStreamOp{
		kind:    rcOpGlobalRef,
		v:       v,
		node:    n,
		block:   block,
		scope:   scope,
		inputId: inputId,
	})
}

func (c *ReferenceCollector) addGlobalReferenceDirect(
	v *Var,
	n rhino.Node,
	block *BasicBlock,
	scope *Scope,
	inputId *rhino.InputId,
) {
	if v != c.lastVar && c.varFilter != nil {
		if v == c.lastRejectedVar || !c.varFilter(v) {
			c.lastRejectedVar = v
			return
		}
	}
	collection := c.getOrCreateReferences(v)
	r := c.refSlab[rcGlobal].New()
	*r = Reference{nameNode: n, basicBlock: block, scope: scope, inputId: inputId}
	c.addToCollection(collection, r, rcGlobal)
}

func (iv *InlineVariables) processParallel(externs rhino.Node, root rhino.Node) bool {
	scopeRoot, items, ok := collectExternAndRootScripts(externs, root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(iv.compiler, len(items))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(iv.compiler)
	restoreReg := enableParallelTypeRegistryReads(iv.compiler)
	defer restoreReg()

	mainSC := NewSyntacticScopeCreator(iv.compiler)
	globalScope := mainSC.CreateScopeOfScope(scopeRoot, nil)
	freezeGlobalScopeForParallelReads(globalScope)
	lazyGlobal := &lazyGlobalScope{
		compiler: iv.compiler,
		root:     scopeRoot,
		scope:    globalScope,
	}
	lazyGlobal.once.Do(func() {})
	globalBlock := NewBasicBlock(nil, scopeRoot)
	varFilter := inlineVariables_Mode_varPredicate[iv.mode]

	statesByItem := make([]*rcParallelScriptState, len(items))
	runParallelWork(iv.compiler, numWorkers, len(items), func(_ int) (func(int, *recordingCompiler), func()) {
		var reusableBeh *InlineVariables
		var reusableCollector *ReferenceCollector
		var reusableSC *SyntacticScopeCreator

		runItem := func(idx int, rc *recordingCompiler) {
			if reusableSC == nil || reusableSC.compiler != rc {
				reusableSC = NewSyntacticScopeCreator(rc)
			}
			var beh *InlineVariables
			var col *ReferenceCollector
			if reusableBeh != nil && !reusableBeh.hasGlobalAliasRetries && len(reusableBeh.aliasRetryHandlers) == 0 {
				beh = reusableBeh
				beh.compiler = rc
				col = reusableCollector
				col.compiler = rc
				col.scopeCreator = reusableSC
				col.clearLocalBase()
				col.maybeReclaimLocals()
				col.blocks.resetHoistedCollected()
			} else {
				beh = NewInlineVariables(rc, iv.mode)
				beh.resetBehavior()
				col = NewReferenceCollectorWithVarFilter(rc, beh, reusableSC, varFilter)
				bRef := beh
				cRef := col
				col.enableLocalReclamation(func() bool {
					if len(bRef.aliasRetryHandlers) == 0 && !bRef.hasGlobalAliasRetries {
						bRef.unpinnedRetries = 0
						cRef.clearLocalBase()
						return true
					}
					if bRef.unpinnedRetries == 0 {
						return true
					}
					bRef.unpinnedRetries = 0
					return false
				})
				reusableBeh = beh
				reusableCollector = col
			}

			ps := &rcParallelScriptState{
				globalScope: globalScope,
				globalBlock: globalBlock,
				rc:          rc,
			}
			ps.curSeg = &ps.rootSeg
			col.parallelState = ps
			beh.onGlobalAliasRetry = func(aliasedVar *Var, retryHandler *InlineVariables_InliningBehavior_AliasInlineRetryHandler) {
				ps.flushCompilerEvents()
				ps.curSeg.ops = append(ps.curSeg.ops, rcStreamOp{
					kind:  rcOpGlobalAliasRetry,
					v:     aliasedVar,
					retry: retryHandler,
				})
			}

			item := items[idx]
			tb := col.builderVal()
			tb.traverseScriptUnderRoot(scopeRoot, item.parent, item.script, lazyGlobal)
			ps.flushCompilerEvents()

			col.parallelState = nil
			beh.onGlobalAliasRetry = nil
			if beh.hasGlobalAliasRetries || len(beh.aliasRetryHandlers) > 0 {
				// Keep beh and col slabs alive until globalScope finishes on main goroutine.
				reusableBeh = nil
				reusableCollector = nil
			}
			statesByItem[idx] = ps
		}
		return runItem, nil
	}, false)

	var hoistedByFn map[rhino.Node]*rcStreamSegment
	for _, ps := range statesByItem {
		if len(ps.hoistedSegs) > 0 {
			if hoistedByFn == nil {
				hoistedByFn = make(map[rhino.Node]*rcStreamSegment, len(ps.hoistedSegs)*2)
			}
			for _, hSeg := range ps.hoistedSegs {
				hoistedByFn[hSeg.fnNode] = hSeg
			}
		}
	}

	iv.resetBehavior()
	globalCollector := NewReferenceCollectorWithVarFilter(iv.compiler, iv, mainSC, varFilter)
	var collectedHoisted map[rhino.Node]struct{}
	if len(hoistedByFn) > 0 {
		collectedHoisted = make(map[rhino.Node]struct{}, len(hoistedByFn))
	}

	var replaySeg func(seg *rcStreamSegment)
	replaySeg = func(seg *rcStreamSegment) {
		if seg == nil {
			return
		}
		for i := range seg.ops {
			op := &seg.ops[i]
			switch op.kind {
			case rcOpCompilerEvents:
				replayCompilerEvents(iv.compiler, op.events)
			case rcOpGlobalRef:
				globalCollector.addGlobalReferenceDirect(op.v, op.node, op.block, op.scope, op.inputId)
			case rcOpGlobalAliasRetry:
				op.retry.outer.outer = iv
				iv.putAliasRetryHandler(op.v, op.retry)
			case rcOpJumpHoisted, rcOpDeclHoisted:
				if _, done := collectedHoisted[op.node]; done {
					continue
				}
				if collectedHoisted == nil {
					collectedHoisted = make(map[rhino.Node]struct{})
				}
				collectedHoisted[op.node] = struct{}{}
				globalCollector.addGlobalReferenceDirect(op.v, op.node.GetFirstChild(), globalBlock, op.scope, op.inputId)
				if hoistedByFn != nil {
					replaySeg(hoistedByFn[op.node])
				}
			}
		}
	}

	for _, ps := range statesByItem {
		replaySeg(&ps.rootSeg)
	}

	gb := globalCollector.builderVal()
	globalTraversal := acquireNodeTraversal(&gb)
	defer releaseNodeTraversal(globalTraversal)
	globalTraversal.InitTraversal(globalScope.GetRootNode())
	globalTraversal.currentNode = globalScope.GetRootNode()
	globalTraversal.PushScopeOfScopeWithQuietly(globalScope, true)
	iv.AfterExitScope(globalTraversal, globalCollector)
	return true
}
