// Port of com/google/javascript/jscomp/InlineObjectLiterals.java @ closure-compiler v20261028

/*
 * Copyright 2011 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// InlineObjectLiterals: using the infrastructure provided by
// ReferenceCollector, identify variables that are only ever assigned to
// object literals and that are never used in entirety, and expand the
// objects into individual variables.
//
// Based on the InlineVariables pass
type InlineObjectLiterals struct {
	compiler AbstractCompiler

	safeNameIdSupplier func() string

	// Variables whose reference information is out of sync with the AST.
	staleVars map[*Var]struct{}

	validProps  smallSet[strpool.ID]
	propSlots   smallMap[strpool.ID, iolPropSlot]
	assignNodes []rhino.Node

	scanState      iolCandidateScanner
	scanner        *iolCandidateScanner
	localBlooms    rcFnBloomCache
	blooms         *rcFnBloomCache
	useLocalPrefix bool
	localPrefix    string
	localCounter   int
}

const InlineObjectLiterals_VAR_PREFIX = "JSCompiler_object_inline_"
const InlineObjectLiterals_STRING_KEY_IDENTIFIER = "string_key"

const iolTimelineKey = "__jsgo_iol"

// Java: InlineObjectLiterals(AbstractCompiler, Supplier<String>)
func NewInlineObjectLiterals(compiler AbstractCompiler, safeNameIdSupplier func() string) *InlineObjectLiterals {
	return &InlineObjectLiterals{compiler: compiler, safeNameIdSupplier: safeNameIdSupplier}
}

func makeIOLSkipFn(blooms *rcFnBloomCache, s *iolCandidateScanner) func(rhino.Node) bool {
	if blooms == nil || s == nil || s.hasInvalidDecl {
		return nil
	}
	active := s.candIDBloom
	active.or(&s.hoistedFnBloom)
	return func(fn rhino.Node) bool {
		if s.hasCandSubtree[fn] {
			return false
		}
		e := blooms.get(fn)
		if (e.nameID > uint32(strpool.EmptyID) && s.candIDBloom.hasID(e.nameID)) || e.bodyBloom.intersects(&active) {
			active.or(&e.hoistedBloom)
			return false
		}
		return true
	}
}

type iolCandidateScanner struct {
	candidateBloom       uint64
	candIDBloom          rcNameBloom
	hoistedFnBloom       rcNameBloom
	candidateScopeBlooms map[rhino.Node]uint64
	fnHasOwnCandidate    map[rhino.Node]bool
	hasCandSubtree       map[rhino.Node]bool
	hasNonFnCandidate    bool
	anyCandFnInHoistedFn bool
	hasInvalidDecl       bool
	candFns              []rhino.Node
}

func (s *iolCandidateScanner) reset() {
	s.candidateBloom = 0
	s.candIDBloom = rcNameBloom{}
	s.hoistedFnBloom = rcNameBloom{}
	clear(s.candidateScopeBlooms)
	clear(s.fnHasOwnCandidate)
	clear(s.hasCandSubtree)
	s.hasNonFnCandidate = false
	s.anyCandFnInHoistedFn = false
	s.hasInvalidDecl = false
	clear(s.candFns)
	s.candFns = s.candFns[:0]
}

func iolHasParamName(fnNode rhino.Node, nameID uint32) bool {
	fn, ok := ast.AsFunction(fnNode)
	if !ok {
		return false
	}
	for p := range fn.Params().Params() {
		n := p.Node()
		if p.Token() == rhino.Token_DEFAULT_VALUE {
			n = p.MustDefaultValue().Target().Node()
		} else if p.Token() == rhino.Token_ITER_REST {
			n = p.MustIterRest().Target().Node()
		}
		if n.IsName() && n.GetStringID() == nameID {
			return true
		}
	}
	return false
}

func iolHasDirectNameDecl(block rhino.Node, nameID uint32, includeVar bool) bool {
	for c := block.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		tok := c.GetToken()
		if tok == rhino.Token_LET || tok == rhino.Token_CONST || (includeVar && tok == rhino.Token_VAR) {
			for lhs := range ast.MustNameDeclaration(c).Declarators() {
				if lhs.Node().IsName() && lhs.Node().GetStringID() == nameID {
					return true
				}
			}
		}
	}
	return false
}

func (s *iolCandidateScanner) addScopeCandidate(scopeRoot rhino.Node, varNameID uint32) {
	if s.candidateScopeBlooms == nil {
		s.candidateScopeBlooms = make(map[rhino.Node]uint64, 16)
		s.fnHasOwnCandidate = make(map[rhino.Node]bool, 16)
		s.hasCandSubtree = make(map[rhino.Node]bool, 64)
	}
	mask := scopeBloomFromID(varNameID)
	s.candidateBloom |= mask
	s.candIDBloom.addID(varNameID)
	s.candidateScopeBlooms[scopeRoot] |= mask
	var ownerFn rhino.Node
	for p := scopeRoot; !p.IsNil(); p = p.GetParent() {
		if p.IsFunction() {
			ownerFn = p
			break
		}
	}
	if !ownerFn.IsNil() {
		if !s.fnHasOwnCandidate[ownerFn] {
			s.fnHasOwnCandidate[ownerFn] = true
			for p := ownerFn; !p.IsNil(); p = p.GetParent() {
				if p.IsFunction() && rcIsHoistedFunctionDeclaration(p) {
					s.anyCandFnInHoistedFn = true
				}
				if s.hasCandSubtree[p] {
					break
				}
				s.hasCandSubtree[p] = true
				if p.IsRoot() {
					break
				}
			}
		}
	} else {
		s.hasNonFnCandidate = true
	}
}

func (s *iolCandidateScanner) mergeFrom(other *iolCandidateScanner) {
	s.hoistedFnBloom.or(&other.hoistedFnBloom)
	if other.hasInvalidDecl {
		s.hasInvalidDecl = true
	}
	if other.candidateBloom == 0 && len(other.candidateScopeBlooms) == 0 {
		return
	}
	if s.candidateScopeBlooms == nil {
		s.candidateScopeBlooms = make(map[rhino.Node]uint64, len(other.candidateScopeBlooms))
		s.fnHasOwnCandidate = make(map[rhino.Node]bool, len(other.fnHasOwnCandidate))
		s.hasCandSubtree = make(map[rhino.Node]bool, len(other.hasCandSubtree))
	}
	s.candidateBloom |= other.candidateBloom
	s.candIDBloom.or(&other.candIDBloom)
	if other.hasNonFnCandidate {
		s.hasNonFnCandidate = true
	}
	if other.anyCandFnInHoistedFn {
		s.anyCandFnInHoistedFn = true
	}
	// order-independent: bitwise OR of per-scope candidate blooms into s.candidateScopeBlooms
	for scopeRoot, mask := range other.candidateScopeBlooms {
		s.candidateScopeBlooms[scopeRoot] |= mask
	}
	// order-independent: boolean union into s.fnHasOwnCandidate
	for fn := range other.fnHasOwnCandidate {
		s.fnHasOwnCandidate[fn] = true
	}
	// order-independent: boolean union into s.hasCandSubtree
	for n := range other.hasCandSubtree {
		s.hasCandSubtree[n] = true
	}
}

func (s *iolCandidateScanner) checkObjectLit(objLit rhino.Node) {
	parent := objLit.GetParent()
	if parent.IsNil() {
		return
	}
	var varNameID uint32
	var isDecl bool
	var gp rhino.Node
	if name, ok := ast.AsName(parent); ok && parent.GetFirstChild() == objLit {
		gp = parent.GetParent()
		if gp.IsNil() || !NodeUtil_IsNameDeclaration(gp) {
			return
		}
		ggp := gp.GetParent()
		if ggp.IsNil() || NodeUtil_IsAnyFor(ggp) {
			return
		}
		varNameID = name.Node().GetStringID()
		isDecl = true
	} else if assign, ok := ast.AsAssign(parent); ok && assign.Value().Is(objLit) {
		lhs, isName := assign.Target().AsName()
		if !isName {
			return
		}
		gp = parent.GetParent()
		if gp.IsNil() || !gp.IsExprResult() {
			return
		}
		varNameID = lhs.Node().GetStringID()
	} else {
		return
	}

	for member := range ast.MustObjectLit(objLit).Members() {
		switch member.Token() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_OBJECT_SPREAD:
			return
		}
	}

	if isDecl {
		if gp.IsVar() {
			for anc := gp.GetParent(); !anc.IsNil() && !anc.IsScript() && !anc.IsRoot(); anc = anc.GetParent() {
				if NodeUtil_IsFunctionBlock(anc) || anc.IsModuleBody() || NodeUtil_IsClassStaticBlock(anc) {
					if !NodeUtil_IsFunctionBlock(anc) || !iolHasParamName(anc.GetParent(), varNameID) {
						s.addScopeCandidate(anc, varNameID)
					}
					break
				}
			}
		} else {
			for anc := gp.GetParent(); !anc.IsNil() && !anc.IsScript() && !anc.IsRoot(); anc = anc.GetParent() {
				if (anc.IsBlock() && scanVarsBlockCreatesScope(anc)) || anc.IsSwitchBody() || NodeUtil_IsAnyFor(anc) || anc.IsModuleBody() {
					s.addScopeCandidate(anc, varNameID)
				}
				if NodeUtil_IsFunctionBlock(anc) || anc.IsModuleBody() || NodeUtil_IsClassStaticBlock(anc) {
					break
				}
			}
		}
	} else {
		for anc := gp.GetParent(); !anc.IsNil() && !anc.IsScript() && !anc.IsRoot(); anc = anc.GetParent() {
			if NodeUtil_IsFunctionBlock(anc) || anc.IsModuleBody() || NodeUtil_IsClassStaticBlock(anc) {
				if NodeUtil_IsFunctionBlock(anc) && iolHasParamName(anc.GetParent(), varNameID) {
					break
				}
				s.addScopeCandidate(anc, varNameID)
				if iolHasDirectNameDecl(anc, varNameID, true) {
					break
				}
			} else if (anc.IsBlock() && scanVarsBlockCreatesScope(anc)) || anc.IsSwitchBody() || NodeUtil_IsAnyFor(anc) {
				s.addScopeCandidate(anc, varNameID)
				if anc.IsBlock() && iolHasDirectNameDecl(anc, varNameID, false) {
					break
				}
			}
		}
	}
}

func (s *iolCandidateScanner) checkNode(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_OBJECTLIT:
		s.checkObjectLit(n)
	case rhino.Token_FUNCTION:
		if rcIsHoistedFunctionDeclaration(n) {
			s.hoistedFnBloom.addID(n.GetFirstChild().GetStringID())
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if !c.IsName() && !c.IsDestructuringLhs() {
				s.hasInvalidDecl = true
			}
		}
	}
}

func (s *iolCandidateScanner) scan(n rhino.Node) {
	s.checkNode(n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.GetToken() == rhino.Token_OBJECTLIT {
				s.checkObjectLit(c)
			}
			continue
		}
		s.scan(c)
	}
}

func (s *iolCandidateScanner) scanIncremental(root rhino.Node, hasChangedDescendant map[rhino.Node]bool) {
	for scr := root.GetFirstChild(); !scr.IsNil(); scr = scr.GetNext() {
		if scr.IsScript() && !hasChangedDescendant[scr] {
			for c := scr.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				fn := c
				if fn.IsExport() {
					fn = fn.GetFirstChild()
				}
				if !fn.IsNil() && fn.IsFunction() && rcIsHoistedFunctionDeclaration(fn) {
					s.hoistedFnBloom.addID(fn.GetFirstChild().GetStringID())
				}
			}
			continue
		}
		s.scanIncrementalNode(scr, hasChangedDescendant)
	}
}

func (s *iolCandidateScanner) scanIncrementalNode(n rhino.Node, hasChangedDescendant map[rhino.Node]bool) {
	s.checkNode(n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.IsFunction() && !hasChangedDescendant[c] {
			if rcIsHoistedFunctionDeclaration(c) {
				s.hoistedFnBloom.addID(c.GetFirstChild().GetStringID())
			}
			continue
		}
		if c.GetFirstChild().IsNil() {
			if c.GetToken() == rhino.Token_OBJECTLIT {
				s.checkObjectLit(c)
			}
			continue
		}
		s.scanIncrementalNode(c, hasChangedDescendant)
	}
}

func iolBuildChangedDescendantSet(root rhino.Node, changed []rhino.Node) map[rhino.Node]bool {
	var hasChangedDescendant map[rhino.Node]bool
	for _, ch := range changed {
		if ch.IsNil() || ch.IsDeleted() {
			continue
		}
		attached := false
		for p := ch; !p.IsNil(); p = p.GetParent() {
			if p == root {
				attached = true
				break
			}
		}
		if !attached {
			continue
		}
		if hasChangedDescendant == nil {
			hasChangedDescendant = make(map[rhino.Node]bool, len(changed)*4)
		}
		for p := ch; !p.IsNil() && !hasChangedDescendant[p]; p = p.GetParent() {
			hasChangedDescendant[p] = true
			if p == root {
				break
			}
		}
	}
	return hasChangedDescendant
}

func (o *InlineObjectLiterals) scanFull(s *iolCandidateScanner, root rhino.Node) {
	if scripts, ok := collectRootScripts(root); ok && len(scripts) >= 8 {
		if numWorkers := parallelThreadCount(o.compiler, len(scripts)); numWorkers > 1 {
			prewarmCompilerForParallelReads(o.compiler)
			shards := make([]iolCandidateScanner, numWorkers)
			runParallelWorkersIndexed(len(scripts), numWorkers,
				func(i int) *iolCandidateScanner { return &shards[i] },
				func(w **iolCandidateScanner, idx int) { (*w).scan(scripts[idx]) }, nil)
			for i := range shards {
				s.mergeFrom(&shards[i])
			}
			return
		}
	}
	s.scan(root)
}

func (o *InlineObjectLiterals) scanCandidates(s *iolCandidateScanner, root rhino.Node) bool {
	o.blooms = acquireRCFnBloomCache(o.compiler, root, &o.localBlooms)
	c := o.compiler
	if c != nil && c.parallelRecorder == nil && root.IsRoot() && root == c.GetJsRoot() {
		if ct := c.GetChangeTracker(); ct != nil {
			if changed := ct.GetChangedScopeNodesForPass(iolTimelineKey); changed != nil {
				if len(changed) == 0 {
					return false
				}
				hasChangedDescendant := iolBuildChangedDescendantSet(root, changed)
				if len(hasChangedDescendant) == 0 {
					return false
				}
				o.blooms.invalidateChanged(hasChangedDescendant)
				s.scanIncremental(root, hasChangedDescendant)
				return s.candidateBloom != 0 || s.hasInvalidDecl
			}
		}
	}
	o.blooms.clear()
	o.scanFull(s, root)
	return s.candidateBloom != 0 || s.hasInvalidDecl
}

func (s *iolCandidateScanner) collectCandFns(n rhino.Node) {
	if !s.hasCandSubtree[n] {
		return
	}
	if n.IsFunction() && s.fnHasOwnCandidate[n] {
		s.candFns = append(s.candFns, n)
		return
	}
	switch n.GetToken() {
	case rhino.Token_CLASS:
		cls := ast.MustClass(n)
		s.collectCandFns(cls.ExtendsNode())
		s.collectCandFns(cls.NameNode())
		s.collectCandFns(cls.Body().Node())
	case rhino.Token_CLASS_MEMBERS:
		for ch := range ast.MustClassBody(n).EvalOrder() {
			s.collectCandFns(ch)
		}
	default:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			s.collectCandFns(c)
		}
	}
}

func newIOLDummyTraversal(rc *ReferenceCollector, root rhino.Node) *NodeTraversal {
	s := newScope(root)
	s.scanDone = true
	s.hasImplicitVars = false
	tb := rc.builderVal()
	t := newNodeTraversal(&tb)
	t.currentNode = root
	t.PushScopeOfScope(s)
	return t
}

func traverseIOLCandFn(t *NodeTraversal, rc *ReferenceCollector, candFn rhino.Node) {
	defer func() {
		if rec := recover(); rec != nil {
			t.ThrowUnexpectedException(jutil.AsThrowable(rec))
		}
	}()
	rc.blocks.resetHoistedCollected()
	rc.lastRejectedVar = nil
	t.currentScript = t.GetEnclosingScript(candFn)
	t.ClearScriptState()
	t.TraverseBranch(candFn, candFn.GetParent())
}

func iolVarFilter(s *iolCandidateScanner, v *Var) bool {
	if v == nil || !v.IsLocal() || v.nameNode.IsNil() {
		return false
	}
	parent := v.nameNode.GetParent()
	if parent.IsNil() || !NodeUtil_IsNameDeclaration(parent) {
		return false
	}
	bloom := scopeBloomFromID(v.nameID)
	return (s.candidateBloom&bloom) == bloom && (s.candidateScopeBlooms[v.scope.GetRootNode()]&bloom) == bloom
}

func iolCandidateScopeOnExit(t *NodeTraversal, s *iolCandidateScanner) *Scope {
	if t.InGlobalScope() || (len(t.scopes) > 0 && t.scopes[len(t.scopes)-1].scope == nil) {
		return nil
	}
	scope := t.GetScope()
	if scope.GetVarCount() == 0 {
		return nil
	}
	if s != nil {
		scopeBloom := s.candidateScopeBlooms[scope.GetRootNode()]
		if scopeBloom == 0 || (scope.varBloom&scopeBloom) == 0 {
			return nil
		}
	}
	return scope
}

func (o *InlineObjectLiterals) varFilter(v *Var) bool { return iolVarFilter(o.scanner, v) }

func (o *InlineObjectLiterals) clearStaleVars() bool {
	if len(o.staleVars) > 0 {
		clear(o.staleVars)
	}
	return true
}

func (o *InlineObjectLiterals) newFilteredCollector(sc *SyntacticScopeCreator) *ReferenceCollector {
	rc := NewReferenceCollectorWithVarFilter(o.compiler, o, sc, o.varFilter)
	rc.skipFn = makeIOLSkipFn(o.blooms, o.scanner)
	rc.enableLocalReclamation(o.clearStaleVars)
	return rc
}

type iolFnRunner struct {
	iol *InlineObjectLiterals
	rc  *ReferenceCollector
	t   *NodeTraversal
}

func newIOLFnRunner(iol *InlineObjectLiterals, s *iolCandidateScanner, root rhino.Node, useLocalPrefix bool) iolFnRunner {
	iol.scanner = s
	iol.useLocalPrefix = useLocalPrefix
	rc := iol.newFilteredCollector(NewSyntacticScopeCreator(iol.compiler))
	return iolFnRunner{iol: iol, rc: rc, t: newIOLDummyTraversal(rc, root)}
}

func (r *iolFnRunner) traverseCandFn(candFn rhino.Node, runPrefix string, fnIdx int) {
	if r.iol.useLocalPrefix {
		r.iol.localPrefix = runPrefix + "_" + strconv.Itoa(fnIdx) + "_"
		r.iol.localCounter = 0
	}
	r.iol.clearStaleVars()
	traverseIOLCandFn(r.t, r.rc, candFn)
}

func (r *iolFnRunner) close() {
	r.iol.useLocalPrefix = false
	for len(r.t.scopes) > 0 {
		r.t.PopScope()
	}
}

type iolScopeBloomEntry struct {
	scopeRoot rhino.Node
	mask      uint64
	idBloom   rcNameBloom
}

type iolCandCheckResult struct {
	actionable            bool
	workerID, start, size int
}

type iolReadOnlyChecker struct {
	convention CodingConvention
	scanner    *iolCandidateScanner
	rc         *ReferenceCollector
	t          *NodeTraversal
	validProps smallSet[strpool.ID]
	scopes     []iolScopeBloomEntry
	actionable bool
}

func newIOLReadOnlyChecker(compiler AbstractCompiler, s *iolCandidateScanner, blooms *rcFnBloomCache, root rhino.Node) *iolReadOnlyChecker {
	c := &iolReadOnlyChecker{convention: compiler.GetCodingConvention(), scanner: s}
	rc := NewReferenceCollectorWithVarFilter(compiler, c, NewSyntacticScopeCreator(compiler), c.varFilter)
	rc.skipFn = makeIOLSkipFn(blooms, s)
	rc.enableLocalReclamation(nil)
	c.rc = rc
	c.t = newIOLDummyTraversal(rc, root)
	return c
}

func (c *iolReadOnlyChecker) close() {
	if c != nil && c.t != nil {
		for len(c.t.scopes) > 0 {
			c.t.PopScope()
		}
	}
}

func (c *iolReadOnlyChecker) varFilter(v *Var) bool { return iolVarFilter(c.scanner, v) }

func (c *iolReadOnlyChecker) AfterExitScope(t *NodeTraversal, referenceMap *ReferenceCollector) {
	scope := iolCandidateScopeOnExit(t, c.scanner)
	if scope == nil {
		return
	}
	var actionableMask uint64
	var idBloom rcNameBloom
	for _, v := range scope.varList {
		if v == nil || iolIsVarForbiddenByConvention(c.convention, v) {
			continue
		}
		if info := referenceMap.GetReferences(v); info != nil && iolIsInlinableObject(info.references, &c.validProps) {
			actionableMask |= scopeBloomFromID(v.nameID)
			idBloom.addID(v.nameID)
		}
	}
	if actionableMask != 0 {
		c.actionable = true
		c.scopes = append(c.scopes, iolScopeBloomEntry{scopeRoot: scope.GetRootNode(), mask: actionableMask, idBloom: idBloom})
	}
}

func (c *iolReadOnlyChecker) checkCandFn(candFn rhino.Node) (actionable bool, start, size int) {
	c.actionable = false
	start = len(c.scopes)
	traverseIOLCandFn(c.t, c.rc, candFn)
	return c.actionable, start, len(c.scopes) - start
}

func iolInspectActionableFn(fn rhino.Node) (script rhino.Node, inHoisted, inGlobalHoisted bool) {
	for p := fn; !p.IsNil() && !p.IsRoot(); p = p.GetParent() {
		if p.IsScript() {
			script = p
		} else if p.IsFunction() && rcIsHoistedFunctionDeclaration(p) {
			inHoisted = true
			if parent := p.GetParent(); !parent.IsNil() &&
				(parent.IsScript() || (parent.IsExport() && !parent.GetParent().IsNil() && parent.GetParent().IsScript())) {
				inGlobalHoisted = true
			}
		}
	}
	return script, inHoisted, inGlobalHoisted
}

func (o *InlineObjectLiterals) filterActionableCandFns(
	s *iolCandidateScanner,
	root rhino.Node,
) (allInDistinctScripts, anyGlobalHoisted bool) {
	n := len(s.candFns)
	results := make([]iolCandCheckResult, n)
	var workerScopes [][]iolScopeBloomEntry

	numWorkers := 1
	if n >= 2 {
		numWorkers = parallelThreadCount(o.compiler, n)
	}
	if numWorkers > 1 {
		prewarmCompilerForParallelReads(o.compiler)
		workerScopes = make([][]iolScopeBloomEntry, numWorkers)
		runParallelWork(o.compiler, numWorkers, n, func(workerID int) (func(int, *recordingCompiler), func()) {
			var checker *iolReadOnlyChecker
			return func(fnIdx int, rc *recordingCompiler) {
					if checker == nil {
						checker = newIOLReadOnlyChecker(rc, s, o.blooms, root)
					}
					if act, start, size := checker.checkCandFn(s.candFns[fnIdx]); act {
						results[fnIdx] = iolCandCheckResult{actionable: true, workerID: workerID, start: start, size: size}
					}
				}, func() {
					if checker != nil {
						workerScopes[workerID] = checker.scopes
						checker.close()
					}
				}
		}, false)
	} else {
		checker := newIOLReadOnlyChecker(o.compiler, s, o.blooms, root)
		for fnIdx, candFn := range s.candFns {
			if act, start, size := checker.checkCandFn(candFn); act {
				results[fnIdx] = iolCandCheckResult{actionable: true, workerID: 0, start: start, size: size}
			}
		}
		workerScopes = [][]iolScopeBloomEntry{checker.scopes}
		checker.close()
	}

	s.candidateBloom = 0
	s.candIDBloom = rcNameBloom{}
	clear(s.candidateScopeBlooms)
	writeIdx := 0
	anyHoisted := false
	allInDistinctScripts = true
	var prevScript rhino.Node
	for fnIdx, candFn := range s.candFns {
		res := results[fnIdx]
		if !res.actionable {
			continue
		}
		s.candFns[writeIdx] = candFn
		writeIdx++
		for _, e := range workerScopes[res.workerID][res.start : res.start+res.size] {
			s.candidateBloom |= e.mask
			s.candIDBloom.or(&e.idBloom)
			s.candidateScopeBlooms[e.scopeRoot] |= e.mask
		}
		script, inHoisted, inGlobalHoisted := iolInspectActionableFn(candFn)
		if inHoisted {
			anyHoisted = true
		}
		if inGlobalHoisted {
			anyGlobalHoisted = true
		}
		if script.IsNil() || script == prevScript {
			allInDistinctScripts = false
		}
		prevScript = script
	}
	clear(s.candFns[writeIdx:])
	s.candFns = s.candFns[:writeIdx]
	s.anyCandFnInHoistedFn = anyHoisted
	return allInDistinctScripts, anyGlobalHoisted
}

func (o *InlineObjectLiterals) runSequentialCandFns(s *iolCandidateScanner, root rhino.Node) {
	r := newIOLFnRunner(o, s, root, false)
	for fnIdx, candFn := range s.candFns {
		r.traverseCandFn(candFn, "", fnIdx)
	}
	r.close()
}

func (o *InlineObjectLiterals) runFastCandFns(s *iolCandidateScanner, root rhino.Node) {
	runPrefix := o.safeNameIdSupplier()
	numWorkers := parallelThreadCount(o.compiler, len(s.candFns))
	if numWorkers <= 1 {
		r := newIOLFnRunner(o, s, root, true)
		for fnIdx, candFn := range s.candFns {
			r.traverseCandFn(candFn, runPrefix, fnIdx)
		}
		r.close()
		return
	}
	prewarmCompilerForParallelReads(o.compiler)
	runParallelWork(o.compiler, numWorkers, len(s.candFns), func(_ int) (func(int, *recordingCompiler), func()) {
		var r iolFnRunner
		return func(fnIdx int, rc *recordingCompiler) {
				if r.iol == nil {
					worker := &InlineObjectLiterals{compiler: rc, safeNameIdSupplier: o.safeNameIdSupplier, blooms: o.blooms}
					r = newIOLFnRunner(worker, s, root, true)
				}
				r.traverseCandFn(s.candFns[fnIdx], runPrefix, fnIdx)
			}, func() {
				if r.t != nil {
					r.close()
				}
			}
	}, false)
}

// Java: InlineObjectLiterals.process(Node, Node)
func (o *InlineObjectLiterals) Process(externs rhino.Node, root rhino.Node) {
	if root.IsNil() {
		return
	}
	s := &o.scanState
	s.reset()
	if !o.scanCandidates(s, root) {
		return
	}

	fastMode := o.compiler.GetOptions().isFastMode()
	if !s.hasNonFnCandidate && !s.hasInvalidDecl {
		s.collectCandFns(root)
		if len(s.candFns) == 0 {
			return
		}
		if fastMode {
			o.runFastCandFns(s, root)
			return
		}
		if len(s.candFns) == 1 ||
			(!s.anyCandFnInHoistedFn && (len(s.candFns) < 8 || parallelThreadCount(o.compiler, len(s.candFns)) <= 1)) {
			o.runSequentialCandFns(s, root)
			return
		}
		allInDistinctScripts, anyGlobalHoisted := o.filterActionableCandFns(s, root)
		if len(s.candFns) == 0 {
			return
		}
		if len(s.candFns) == 1 || !s.anyCandFnInHoistedFn || (allInDistinctScripts && !anyGlobalHoisted) {
			o.runSequentialCandFns(s, root)
			return
		}
	}

	sc := NewSyntacticScopeCreator(o.compiler)
	o.scanner = s
	o.useLocalPrefix = false
	o.clearStaleVars()
	o.newFilteredCollector(sc).Process(externs, root)
}

type iolPropSlot struct {
	varName  string
	initVal  rhino.Node
	assigned bool
}

func (o *InlineObjectLiterals) nextID() string {
	if o.useLocalPrefix {
		id := o.localPrefix + strconv.Itoa(o.localCounter)
		o.localCounter++
		return id
	}
	return o.safeNameIdSupplier()
}

func (o *InlineObjectLiterals) addStaleVar(v *Var) {
	if v == nil {
		return
	}
	if o.staleVars == nil {
		o.staleVars = make(map[*Var]struct{}, 8)
	}
	o.staleVars[v] = struct{}{}
}

// Java: InlineObjectLiterals.InliningBehavior.afterExitScope(NodeTraversal, ReferenceMap)
func (o *InlineObjectLiterals) AfterExitScope(t *NodeTraversal, referenceMap *ReferenceCollector) {
	scope := iolCandidateScopeOnExit(t, o.scanner)
	if scope == nil {
		return
	}
	for _, v := range scope.varList {
		if v == nil || o.isVarInlineForbidden(v) {
			continue
		}
		referenceInfo := referenceMap.GetReferences(v)
		if referenceInfo != nil && o.isInlinableObject(referenceInfo.references) {
			o.addStaleVar(v)
			o.splitObject(v, referenceInfo.GetInitializingReference(), referenceInfo)
		}
	}
}

// Java: InlineObjectLiterals.InliningBehavior.recordStaleVarReferencesInTree(Node, Scope)
func (o *InlineObjectLiterals) recordStaleVarReferencesInTree(root rhino.Node, scope *Scope) {
	if root.IsName() {
		o.addStaleVar(scope.GetVarByID(root.GetStringID()))
	}
	if !root.IsFunction() {
		for c := root.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			o.recordStaleVarReferencesInTree(c, scope)
		}
	}
}

func iolIsVarForbiddenByConvention(cc CodingConvention, var_ *Var) bool {
	return var_.IsGlobal() ||
		var_.IsExtern() ||
		cc.IsExportedWithLocal(var_.GetName(), true /* local= */) ||
		cc.IsPropertyRenameFunction(var_.GetNameNode())
}

// Java: InlineObjectLiterals.InliningBehavior.isVarInlineForbidden(Var)
func (o *InlineObjectLiterals) isVarInlineForbidden(var_ *Var) bool {
	if _, stale := o.staleVars[var_]; stale {
		return true
	}
	return iolIsVarForbiddenByConvention(o.compiler.GetCodingConvention(), var_)
}

// Java: InlineObjectLiterals.InliningBehavior.isInlinableObject(List<Reference>)
func (o *InlineObjectLiterals) isInlinableObject(refs []*Reference) bool {
	return iolIsInlinableObject(refs, &o.validProps)
}

func iolIsInlinableObject(refs []*Reference, validProps *smallSet[strpool.ID]) bool {
	ret := false
	validProps.clear()
	for _, ref := range refs {
		name := ref.GetNode()
		parent := ref.GetParent()
		grandparent := ref.GetGrandparent()

		// Matches Java v20261028 NodeUtil.isNormalOrOptChainGetProp(parent) parity (b/571220152).
		if gp, ok := ast.AsNormalOrOptChainGetProp(parent); ok {
			if !gp.Object().Is(name) {
				panic(jutil.NewIllegalStateException())
			}
			// A call target may be using the object as a 'this' value.
			if NodeUtil_IsInvocationTarget(parent) {
				return false
			}
			// Deleting a property has different semantics from deleting a variable.
			if grandparent.IsDelProp() {
				return false
			}

			propID := nodeStringID(parent)
			if !validProps.contains(propID) {
				if NodeUtil_IsNameDeclOrSimpleAssignLhs(parent, grandparent) {
					validProps.add(propID)
				} else {
					return false
				}
			}
			continue
		}

		// Only rewrite VAR/LET/CONST declarations or simple assignment statements.
		if !iolIsVarOrAssignExprLhs(name) {
			return false
		}

		// Don't rewrite declarations inside for-loop headers.
		if NodeUtil_IsNameDeclaration(parent) && NodeUtil_IsAnyFor(grandparent) {
			return false
		}

		val := ref.GetAssignedValue()
		if val.IsNil() {
			continue
		}

		objLit, ok := ast.AsObjectLit(val)
		if !ok {
			return false
		}

		// Disallow getters/setters, computed keys, spreads, and self-referential values.
		for member := range objLit.Members() {
			switch member.Token() {
			case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_OBJECT_SPREAD:
				return false
			case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_STRING_KEY:
			default:
				panic(jutil.NewIllegalStateExceptionMsg(
					"Unexpected child of OBJECTLIT: " + member.Node().ToStringTree()))
			}

			pkm := ast.MustPlainKeyMember(member.Node())
			validProps.add(nodeStringID(member.Node()))

			childVal := pkm.Value().Node()
			for _, t := range refs {
				for refNode := t.GetParent(); !NodeUtil_IsStatementBlock(refNode); refNode = refNode.GetParent() {
					if refNode == childVal {
						return false
					}
				}
			}
		}

		ret = true
	}
	return ret
}

// Java: InlineObjectLiterals.InliningBehavior.isVarOrAssignExprLhs(Node)
func iolIsVarOrAssignExprLhs(n rhino.Node) bool {
	parent := n.GetParent()
	if NodeUtil_IsNameDeclaration(parent) {
		return true
	}
	assign, ok := ast.AsAssign(parent)
	return ok && assign.Target().Is(n) && parent.GetParent().IsExprResult()
}

func (o *InlineObjectLiterals) ensurePropSlot(propID strpool.ID, sanitizeKey bool) {
	if o.propSlots.indexOf(propID) >= 0 {
		return
	}
	varIdent := strpool.String(propID)
	if sanitizeKey && !rhino.TokenStream_IsJSIdentifier(varIdent) {
		varIdent = InlineObjectLiterals_STRING_KEY_IDENTIFIER
	}
	o.propSlots.put(propID, iolPropSlot{varName: InlineObjectLiterals_VAR_PREFIX + varIdent + "_" + o.nextID()})
}

// Java: InlineObjectLiterals.InliningBehavior.computeVarList(ReferenceCollection)
func (o *InlineObjectLiterals) computeVarList(referenceInfo *ReferenceCollection) {
	o.propSlots.clear()

	for _, ref := range referenceInfo.references {
		if ref.IsLvalue() || ref.IsInitializingDeclaration() {
			if val := ref.GetAssignedValue(); !val.IsNil() {
				for member := range ast.MustObjectLit(val).Members() {
					o.ensurePropSlot(nodeStringID(member.Node()), true)
				}
			}
		} else if !NodeUtil_IsNameDeclaration(ref.GetParent()) {
			gp := ast.MustNormalOrOptChainGetProp(ref.GetParent())
			o.ensurePropSlot(nodeStringID(gp.Node()), false)
		}
	}
}

// Java: InlineObjectLiterals.InliningBehavior.fillInitialValues(Reference, Map<String, Node>)
func (o *InlineObjectLiterals) fillInitialValues(init *Reference) {
	slots := o.propSlots.values()
	objLit := ast.MustObjectLit(init.GetAssignedValue())
	for member := range objLit.Members() {
		key := member.Node()
		idx := o.propSlots.indexOf(nodeStringID(key))
		slots[idx].initVal = key.RemoveFirstChild()
	}
}

// replaceAssignmentExpression replaces `x = {...}` with `t1=a, t2=b, true`.
//
// Java: InlineObjectLiterals.InliningBehavior.replaceAssignmentExpression(Var, Reference, Map<String, String>)
func (o *InlineObjectLiterals) replaceAssignmentExpression(v *Var, ref *Reference) {
	nodes := o.assignNodes[:0]
	val := ref.GetAssignedValue()
	o.recordStaleVarReferencesInTree(val, v.GetScope().(*Scope))
	objLit := ast.MustObjectLit(val)
	slots := o.propSlots.values()
	for i := range slots {
		slots[i].assigned = false
	}
	for member := range objLit.Members() {
		keyNode := member.Node()
		idx := o.propSlots.indexOf(nodeStringID(keyNode))
		value := keyNode.RemoveFirstChild()
		nodes = append(nodes, rhino.IR_Assign(rhino.IR_Name(slots[idx].varName), value))
		slots[idx].assigned = true
	}

	for i := range slots {
		if !slots[i].assigned {
			nodes = append(nodes, rhino.IR_Assign(rhino.IR_Name(slots[i].varName), NodeUtil_NewUndefinedNode(rhino.NilNode)))
		}
	}

	var replacement rhino.Node
	if len(nodes) == 0 {
		replacement = rhino.IR_TrueNode()
	} else {
		// All object-literal assignments evaluate to truthy; preserve truthiness with trailing `true`
		// folded into a left-associative COMMA tree matching parser shape.
		nodes = append(nodes, rhino.IR_TrueNode())
		replacement = nodes[0]
		for _, next := range nodes[1:] {
			replacement = rhino.IR_Comma(replacement, next)
		}
	}
	clear(nodes)
	o.assignNodes = nodes[:0]

	replace := ref.GetParent()
	replacement.SrcrefTreeIfMissing(replace)

	if NodeUtil_IsNameDeclaration(replace) {
		replace.ReplaceWith(NodeUtil_NewExpr(replacement))
	} else {
		replace.ReplaceWith(replacement)
	}
}

// Java: InlineObjectLiterals.InliningBehavior.splitObject(Var, Reference, ReferenceCollection)
func (o *InlineObjectLiterals) splitObject(v *Var, init *Reference, referenceInfo *ReferenceCollection) {
	o.computeVarList(referenceInfo)

	var vnode rhino.Node
	defined := referenceInfo.IsWellDefined() && NodeUtil_IsNameDeclaration(init.GetParent())
	if defined {
		vnode = init.GetParent()
		o.fillInitialValues(init)
	} else if v.GetParentNode().IsLet() || v.GetParentNode().IsConst() {
		scopeRoot := v.GetScope().(*Scope).GetRootNode()
		vnode = scopeRoot.GetFirstChild()
		if !NodeUtil_IsStatement(vnode) && NodeUtil_IsStatement(scopeRoot) {
			vnode = scopeRoot
		} else if scopeRoot.IsSwitchBody() {
			vnode = scopeRoot.GetParent()
		}
	} else {
		vnode = v.GetScope().(*Scope).GetClosestHoistScope().GetRootNode().GetFirstChild()
	}

	jutil.CheckStateMsg1(NodeUtil_IsStatement(vnode), "%s", vnode)

	slots := o.propSlots.values()
	for i := range slots {
		slot := &slots[i]
		val := slot.initVal
		newVarNode := NodeUtil_NewVarNode(slot.varName, val)
		if val.IsNil() {
			newVarNode.SrcrefTreeIfMissing(vnode)
		} else {
			o.recordStaleVarReferencesInTree(val, v.GetScope().(*Scope))
		}
		newVarNode.InsertBefore(vnode)
		o.compiler.ReportChangeToEnclosingScope(vnode)
	}

	if defined {
		o.compiler.ReportChangeToEnclosingScope(vnode.GetParent())
		vnode.Detach()
	}

	for _, ref := range referenceInfo.references {
		if defined && ref == init {
			continue
		}
		o.compiler.ReportChangeToEnclosingScope(ref.GetNode())

		if ref.IsLvalue() {
			o.replaceAssignmentExpression(v, ref)
		} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
			ref.GetParent().Detach()
		} else {
			gpNode := ref.GetParent()
			idx := o.propSlots.indexOf(nodeStringID(gpNode))
			jutil.CheckState(idx >= 0)

			replacement := rhino.IR_Name(slots[idx].varName).SrcrefIfMissing(gpNode)
			gpNode.ReplaceWith(replacement)
		}
	}
}
