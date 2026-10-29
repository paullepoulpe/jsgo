// Port of com/google/javascript/jscomp/InlineObjectLiterals.java @ closure-compiler v20261021

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
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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

	validPropInline [8]uint32
	validPropLen    int
	validPropMap    map[uint32]struct{}

	propSlots   []iolPropSlot
	propIndex   map[uint32]int
	assignNodes []rhino.Node

	scanner      *iolCandidateScanner
	idSupplier   func() string
	reportChange func(rhino.Node)
}

const InlineObjectLiterals_VAR_PREFIX = "JSCompiler_object_inline_"
const InlineObjectLiterals_STRING_KEY_IDENTIFIER = "string_key"

// Java: InlineObjectLiterals(AbstractCompiler, Supplier<String>)
func NewInlineObjectLiterals(compiler AbstractCompiler, safeNameIdSupplier func() string) *InlineObjectLiterals {
	return &InlineObjectLiterals{compiler: compiler, safeNameIdSupplier: safeNameIdSupplier}
}

type iolCandidateScanner struct {
	candidateBloom       uint64
	candidateScopeBlooms map[rhino.Node]uint64
	fnHasOwnCandidate    map[rhino.Node]bool
	hasCandSubtree       map[rhino.Node]bool
	hasNonFnCandidate    bool
	anyCandFnInHoistedFn bool
	hasInvalidDecl       bool
	candFns              []rhino.Node
}

func iolHasParamName(fnNode rhino.Node, nameID uint32) bool {
	fn, ok := ast.AsFunction(fnNode)
	if !ok {
		return false
	}
	for p := range fn.Params().Params() {
		switch p.Token() {
		case rhino.Token_NAME:
			if p.Node().GetStringID() == nameID {
				return true
			}
		case rhino.Token_DEFAULT_VALUE:
			if target := p.MustDefaultValue().Target().Node(); target.IsName() && target.GetStringID() == nameID {
				return true
			}
		case rhino.Token_ITER_REST:
			if target := p.MustIterRest().Target().Node(); target.IsName() && target.GetStringID() == nameID {
				return true
			}
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

func (s *iolCandidateScanner) addScopeCandidate(scopeRoot rhino.Node, mask uint64) {
	if s.candidateScopeBlooms == nil {
		s.candidateScopeBlooms = make(map[rhino.Node]uint64, 16)
		s.fnHasOwnCandidate = make(map[rhino.Node]bool, 16)
		s.hasCandSubtree = make(map[rhino.Node]bool, 64)
	}
	s.candidateBloom |= mask
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
			}
		}
	} else {
		s.hasNonFnCandidate = true
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

	mask := scopeBloomFromID(varNameID)
	if isDecl {
		if gp.IsVar() {
			for anc := gp.GetParent(); !anc.IsNil() && !anc.IsScript() && !anc.IsRoot(); anc = anc.GetParent() {
				if NodeUtil_IsFunctionBlock(anc) {
					if !iolHasParamName(anc.GetParent(), varNameID) {
						s.addScopeCandidate(anc, mask)
					}
					break
				}
				if anc.IsModuleBody() || NodeUtil_IsClassStaticBlock(anc) {
					s.addScopeCandidate(anc, mask)
					break
				}
			}
		} else {
			for anc := gp.GetParent(); !anc.IsNil() && !anc.IsScript() && !anc.IsRoot(); anc = anc.GetParent() {
				if (anc.IsBlock() && scanVarsBlockCreatesScope(anc)) || anc.IsSwitchBody() || NodeUtil_IsAnyFor(anc) || anc.IsModuleBody() {
					s.addScopeCandidate(anc, mask)
				}
				if NodeUtil_IsFunctionBlock(anc) || anc.IsModuleBody() || NodeUtil_IsClassStaticBlock(anc) {
					break
				}
			}
		}
	} else {
		for anc := gp.GetParent(); !anc.IsNil() && !anc.IsScript() && !anc.IsRoot(); anc = anc.GetParent() {
			if NodeUtil_IsFunctionBlock(anc) {
				if iolHasParamName(anc.GetParent(), varNameID) {
					break
				}
				s.addScopeCandidate(anc, mask)
				if iolHasDirectNameDecl(anc, varNameID, true) {
					break
				}
			} else if anc.IsModuleBody() || NodeUtil_IsClassStaticBlock(anc) {
				s.addScopeCandidate(anc, mask)
				if iolHasDirectNameDecl(anc, varNameID, true) {
					break
				}
			} else if (anc.IsBlock() && scanVarsBlockCreatesScope(anc)) || anc.IsSwitchBody() || NodeUtil_IsAnyFor(anc) {
				s.addScopeCandidate(anc, mask)
				if anc.IsBlock() && iolHasDirectNameDecl(anc, varNameID, false) {
					break
				}
			}
		}
	}
}

func (s *iolCandidateScanner) scan(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_OBJECTLIT:
		s.checkObjectLit(n)
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if !c.IsName() && !c.IsDestructuringLhs() {
				s.hasInvalidDecl = true
			}
		}
	}
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
		body := ast.MustClassBody(n)
		for member := range body.Members() {
			if cm, ok := member.AsComputedMember(); ok {
				s.collectCandFns(cm.Key().Node())
			}
		}
		for member := range body.Members() {
			if cm, ok := member.AsComputedMember(); ok {
				if val, hasVal := cm.Value(); hasVal {
					s.collectCandFns(val.Node())
				}
			} else {
				s.collectCandFns(member.Node())
			}
		}
	default:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			s.collectCandFns(c)
		}
	}
}

func newIOLDummyGlobalScope(root rhino.Node) *Scope {
	s := newScope(root)
	s.scanDone = true
	s.hasImplicitVars = false
	return s
}

func (o *InlineObjectLiterals) newFilteredCollector(
	sc *SyntacticScopeCreator,
	varFilter func(*Var) bool,
) *ReferenceCollector {
	rc := NewReferenceCollectorWithVarFilter(o.compiler, o, sc, varFilter)
	rc.enableLocalReclamation(func() bool {
		clear(o.staleVars)
		return true
	})
	return rc
}

func traverseIOLCandFn(rc *ReferenceCollector, root rhino.Node, dummyGlobalScope *Scope, candFn rhino.Node) {
	tb := rc.builderVal()
	t := newNodeTraversal(&tb)
	t.currentNode = root
	t.PushScopeOfScope(dummyGlobalScope)
	t.currentScript = t.GetEnclosingScript(candFn)
	t.ClearScriptState()
	t.TraverseBranch(candFn, candFn.GetParent())
	t.PopScope()
}

// Java: InlineObjectLiterals.process(Node, Node)
func (o *InlineObjectLiterals) Process(externs rhino.Node, root rhino.Node) {
	if root.IsNil() {
		return
	}
	var s iolCandidateScanner
	s.scan(root)
	if s.candidateBloom == 0 && !s.hasInvalidDecl {
		return
	}

	varFilter := func(v *Var) bool {
		if v == nil || !v.IsLocal() {
			return false
		}
		nameNode := v.nameNode
		if nameNode.IsNil() {
			return false
		}
		parent := nameNode.GetParent()
		if parent.IsNil() || !NodeUtil_IsNameDeclaration(parent) {
			return false
		}
		bloom := scopeBloomFromID(v.nameID)
		if (s.candidateBloom & bloom) != bloom {
			return false
		}
		return (s.candidateScopeBlooms[v.scope.GetRootNode()] & bloom) == bloom
	}

	opts := o.compiler.GetOptions()
	fastMode := opts.isFastMode()

	if !s.hasNonFnCandidate && !s.hasInvalidDecl {
		s.collectCandFns(root)
		if len(s.candFns) == 0 {
			return
		}
		if fastMode && len(s.candFns) >= 2 && opts.GetNumParallelThreads() > 1 {
			numWorkers := min(opts.GetNumParallelThreads(), len(s.candFns))
			runPrefix := o.safeNameIdSupplier()
			changedPerFn := make([][]rhino.Node, len(s.candFns))
			var nextIdx atomic.Int32
			var wg sync.WaitGroup
			wg.Add(numWorkers)
			for w := 0; w < numWorkers; w++ {
				go func() {
					defer wg.Done()
					sc := NewSyntacticScopeCreator(o.compiler)
					worker := &InlineObjectLiterals{
						compiler:           o.compiler,
						safeNameIdSupplier: o.safeNameIdSupplier,
						scanner:            &s,
					}
					dummyGlobalScope := newIOLDummyGlobalScope(root)
					for {
						fnIdx := int(nextIdx.Add(1) - 1)
						if fnIdx >= len(s.candFns) {
							break
						}
						candFn := s.candFns[fnIdx]
						localCounter := 0
						prefix := runPrefix + "_" + strconv.Itoa(fnIdx) + "_"
						clear(worker.staleVars)
						worker.idSupplier = func() string {
							id := prefix + strconv.Itoa(localCounter)
							localCounter++
							return id
						}
						worker.reportChange = func(n rhino.Node) {
							if n.IsNil() {
								return
							}
							if n.IsScript() {
								changedPerFn[fnIdx] = append(changedPerFn[fnIdx], n)
							} else if cs := ChangeTracker_GetEnclosingChangeScopeRoot(n.GetParent()); !cs.IsNil() {
								changedPerFn[fnIdx] = append(changedPerFn[fnIdx], cs)
							}
						}
						rc := worker.newFilteredCollector(sc, varFilter)
						traverseIOLCandFn(rc, root, dummyGlobalScope, candFn)
					}
				}()
			}
			wg.Wait()
			for _, list := range changedPerFn {
				for _, cs := range list {
					o.compiler.ReportChangeToChangeScope(cs)
				}
			}
			return
		}
		if fastMode || len(s.candFns) == 1 || !s.anyCandFnInHoistedFn {
			sc := NewSyntacticScopeCreator(o.compiler)
			o.scanner = &s
			o.idSupplier = nil
			o.reportChange = nil
			dummyGlobalScope := newIOLDummyGlobalScope(root)
			if fastMode && len(s.candFns) >= 1 {
				runPrefix := o.safeNameIdSupplier()
				for fnIdx, candFn := range s.candFns {
					localCounter := 0
					prefix := runPrefix + "_" + strconv.Itoa(fnIdx) + "_"
					clear(o.staleVars)
					o.idSupplier = func() string {
						id := prefix + strconv.Itoa(localCounter)
						localCounter++
						return id
					}
					rc := o.newFilteredCollector(sc, varFilter)
					traverseIOLCandFn(rc, root, dummyGlobalScope, candFn)
				}
				o.idSupplier = nil
				return
			}
			callback := o.newFilteredCollector(sc, varFilter)
			tb := callback.builderVal()
			t := newNodeTraversal(&tb)
			t.currentNode = root
			t.PushScopeOfScope(dummyGlobalScope)
			for _, candFn := range s.candFns {
				clear(o.staleVars)
				t.currentScript = t.GetEnclosingScript(candFn)
				t.ClearScriptState()
				t.TraverseBranch(candFn, candFn.GetParent())
			}
			t.PopScope()
			return
		}
	}

	sc := NewSyntacticScopeCreator(o.compiler)
	o.scanner = &s
	o.idSupplier = nil
	o.reportChange = nil
	clear(o.staleVars)
	o.newFilteredCollector(sc, varFilter).Process(externs, root)
}

type iolPropSlot struct {
	propID   uint32
	varName  string
	initVal  rhino.Node
	assigned bool
}

func (o *InlineObjectLiterals) nextID() string {
	if o.idSupplier != nil {
		return o.idSupplier()
	}
	return o.safeNameIdSupplier()
}

func (o *InlineObjectLiterals) reportScopeChange(n rhino.Node) {
	if o.reportChange != nil {
		o.reportChange(n)
		return
	}
	o.compiler.ReportChangeToEnclosingScope(n)
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
	if t.InGlobalScope() {
		return
	}
	if len(t.scopes) > 0 && t.scopes[len(t.scopes)-1].scope == nil {
		return
	}
	scope := t.GetScope()
	if scope.GetVarCount() == 0 {
		return
	}
	if o.scanner != nil {
		scopeBloom := o.scanner.candidateScopeBlooms[scope.GetRootNode()]
		if scopeBloom == 0 || (scope.varBloom&scopeBloom) == 0 {
			return
		}
	}
	for _, v := range scope.varList {
		if v == nil || o.isVarInlineForbidden(v) {
			continue
		}
		referenceInfo := referenceMap.GetReferences(v)
		if referenceInfo == nil {
			continue
		}
		if o.isInlinableObject(referenceInfo.references) {
			o.addStaleVar(v)
			init := referenceInfo.GetInitializingReference()
			o.splitObject(v, init, referenceInfo)
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

// Java: InlineObjectLiterals.InliningBehavior.isVarInlineForbidden(Var)
func (o *InlineObjectLiterals) isVarInlineForbidden(var_ *Var) bool {
	_, stale := o.staleVars[var_]
	return var_.IsGlobal() ||
		var_.IsExtern() ||
		o.compiler.GetCodingConvention().IsExportedWithLocal(var_.GetName(), true /* local= */) ||
		o.compiler.GetCodingConvention().IsPropertyRenameFunction(var_.GetNameNode()) ||
		stale
}

func (o *InlineObjectLiterals) clearValidProps() {
	o.validPropLen = 0
	if len(o.validPropMap) > 0 {
		clear(o.validPropMap)
	}
}

func (o *InlineObjectLiterals) hasValidProp(id uint32) bool {
	if o.validPropLen <= len(o.validPropInline) {
		for i := 0; i < o.validPropLen; i++ {
			if o.validPropInline[i] == id {
				return true
			}
		}
		return false
	}
	_, ok := o.validPropMap[id]
	return ok
}

func (o *InlineObjectLiterals) addValidProp(id uint32) {
	if o.hasValidProp(id) {
		return
	}
	if o.validPropLen < len(o.validPropInline) {
		o.validPropInline[o.validPropLen] = id
		o.validPropLen++
		return
	}
	if o.validPropLen == len(o.validPropInline) {
		if o.validPropMap == nil {
			o.validPropMap = make(map[uint32]struct{}, 16)
		}
		for _, existing := range o.validPropInline {
			o.validPropMap[existing] = struct{}{}
		}
	}
	o.validPropMap[id] = struct{}{}
	o.validPropLen++
}

func (o *InlineObjectLiterals) resetPropSlots() {
	clear(o.propSlots)
	o.propSlots = o.propSlots[:0]
	if len(o.propIndex) > 0 {
		clear(o.propIndex)
	}
}

func (o *InlineObjectLiterals) findPropSlot(propID uint32) int {
	if len(o.propSlots) <= 8 {
		for i := range o.propSlots {
			if o.propSlots[i].propID == propID {
				return i
			}
		}
		return -1
	}
	if idx, ok := o.propIndex[propID]; ok {
		return idx
	}
	return -1
}

func (o *InlineObjectLiterals) addPropSlot(propID uint32, varName string) {
	idx := len(o.propSlots)
	o.propSlots = append(o.propSlots, iolPropSlot{propID: propID, varName: varName})
	if len(o.propSlots) == 9 {
		if o.propIndex == nil {
			o.propIndex = make(map[uint32]int, 16)
		}
		for i := range o.propSlots {
			o.propIndex[o.propSlots[i].propID] = i
		}
	} else if len(o.propSlots) > 9 {
		o.propIndex[propID] = idx
	}
}

// Java: InlineObjectLiterals.InliningBehavior.isInlinableObject(List<Reference>)
func (o *InlineObjectLiterals) isInlinableObject(refs []*Reference) bool {
	ret := false
	o.clearValidProps()
	for _, ref := range refs {
		name := ref.GetNode()
		parent := ref.GetParent()
		grandparent := ref.GetGrandparent()

		if gp, ok := ast.AsNormalOrOptChainGetProp(parent); ok {
			jutil.CheckState(gp.Object().Is(name))
			// A call target may be using the object as a 'this' value.
			if NodeUtil_IsInvocationTarget(parent) {
				return false
			}
			// Deleting a property has different semantics from deleting a variable.
			if grandparent.IsDelProp() {
				return false
			}

			propID := parent.GetStringID()
			if !o.hasValidProp(propID) {
				if NodeUtil_IsNameDeclOrSimpleAssignLhs(parent, grandparent) {
					o.addValidProp(propID)
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
			o.addValidProp(member.Node().GetStringID())

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

// Java: InlineObjectLiterals.InliningBehavior.computeVarList(ReferenceCollection)
func (o *InlineObjectLiterals) computeVarList(referenceInfo *ReferenceCollection) {
	o.resetPropSlots()

	for _, ref := range referenceInfo.references {
		if ref.IsLvalue() || ref.IsInitializingDeclaration() {
			val := ref.GetAssignedValue()
			if !val.IsNil() {
				for member := range ast.MustObjectLit(val).Members() {
					propID := member.Node().GetStringID()
					if o.findPropSlot(propID) >= 0 {
						continue
					}
					varname := rhino.RhinoStringPool_StringFromID(propID)
					varIdent := varname
					if !rhino.TokenStream_IsJSIdentifier(varname) {
						varIdent = InlineObjectLiterals_STRING_KEY_IDENTIFIER
					}
					o.addPropSlot(propID, InlineObjectLiterals_VAR_PREFIX+varIdent+"_"+o.nextID())
				}
			}
		} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
			// Uninitialized declaration; no property keys to register.
		} else {
			gp := ast.MustNormalOrOptChainGetProp(ref.GetParent())
			propID := gp.Node().GetStringID()
			if o.findPropSlot(propID) < 0 {
				varname := rhino.RhinoStringPool_StringFromID(propID)
				o.addPropSlot(propID, InlineObjectLiterals_VAR_PREFIX+varname+"_"+o.nextID())
			}
		}
	}
}

// Java: InlineObjectLiterals.InliningBehavior.fillInitialValues(Reference, Map<String, Node>)
func (o *InlineObjectLiterals) fillInitialValues(init *Reference) {
	objLit := ast.MustObjectLit(init.GetAssignedValue())
	for member := range objLit.Members() {
		key := member.Node()
		idx := o.findPropSlot(key.GetStringID())
		o.propSlots[idx].initVal = key.RemoveFirstChild()
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
	for i := range o.propSlots {
		o.propSlots[i].assigned = false
	}
	for member := range objLit.Members() {
		keyNode := member.Node()
		idx := o.findPropSlot(keyNode.GetStringID())
		value := keyNode.RemoveFirstChild()
		nodes = append(nodes, rhino.IR_Assign(rhino.IR_Name(o.propSlots[idx].varName), value))
		o.propSlots[idx].assigned = true
	}

	for i := range o.propSlots {
		if !o.propSlots[i].assigned {
			nodes = append(nodes, rhino.IR_Assign(rhino.IR_Name(o.propSlots[i].varName), NodeUtil_NewUndefinedNode(rhino.NilNode)))
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

	for i := range o.propSlots {
		slot := &o.propSlots[i]
		val := slot.initVal
		newVarNode := NodeUtil_NewVarNode(slot.varName, val)
		if val.IsNil() {
			newVarNode.SrcrefTreeIfMissing(vnode)
		} else {
			o.recordStaleVarReferencesInTree(val, v.GetScope().(*Scope))
		}
		newVarNode.InsertBefore(vnode)
		o.reportScopeChange(vnode)
	}

	if defined {
		o.reportScopeChange(vnode.GetParent())
		vnode.Detach()
	}

	for _, ref := range referenceInfo.references {
		if defined && ref == init {
			continue
		}
		o.reportScopeChange(ref.GetNode())

		if ref.IsLvalue() {
			o.replaceAssignmentExpression(v, ref)
		} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
			ref.GetParent().Detach()
		} else {
			gpNode := ref.GetParent()
			idx := o.findPropSlot(gpNode.GetStringID())
			jutil.CheckState(idx >= 0)

			replacement := rhino.IR_Name(o.propSlots[idx].varName).SrcrefIfMissing(gpNode)
			gpNode.ReplaceWith(replacement)
		}
	}
}
