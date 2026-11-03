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
	useLocalPrefix bool
	localPrefix    string
	localCounter   int
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

func (s *iolCandidateScanner) reset() {
	s.candidateBloom = 0
	if len(s.candidateScopeBlooms) > 0 {
		clear(s.candidateScopeBlooms)
	}
	if len(s.fnHasOwnCandidate) > 0 {
		clear(s.fnHasOwnCandidate)
	}
	if len(s.hasCandSubtree) > 0 {
		clear(s.hasCandSubtree)
	}
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
				if NodeUtil_IsFunctionBlock(anc) || anc.IsModuleBody() || NodeUtil_IsClassStaticBlock(anc) {
					if !NodeUtil_IsFunctionBlock(anc) || !iolHasParamName(anc.GetParent(), varNameID) {
						s.addScopeCandidate(anc, mask)
					}
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
			if NodeUtil_IsFunctionBlock(anc) || anc.IsModuleBody() || NodeUtil_IsClassStaticBlock(anc) {
				if NodeUtil_IsFunctionBlock(anc) && iolHasParamName(anc.GetParent(), varNameID) {
					break
				}
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
		for ch := range ast.MustClassBody(n).EvalOrder() {
			s.collectCandFns(ch)
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

func (o *InlineObjectLiterals) varFilter(v *Var) bool {
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
	s := o.scanner
	bloom := scopeBloomFromID(v.nameID)
	if (s.candidateBloom & bloom) != bloom {
		return false
	}
	return (s.candidateScopeBlooms[v.scope.GetRootNode()] & bloom) == bloom
}

func (o *InlineObjectLiterals) clearStaleVars() bool {
	if len(o.staleVars) > 0 {
		clear(o.staleVars)
	}
	return true
}

func (o *InlineObjectLiterals) newFilteredCollector(sc *SyntacticScopeCreator) *ReferenceCollector {
	rc := NewReferenceCollectorWithVarFilter(o.compiler, o, sc, o.varFilter)
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
	tb := rc.builderVal()
	t := newNodeTraversal(&tb)
	t.currentNode = root
	t.PushScopeOfScope(newIOLDummyGlobalScope(root))
	return iolFnRunner{iol: iol, rc: rc, t: t}
}

func (r *iolFnRunner) traverseCandFn(candFn rhino.Node, runPrefix string, fnIdx int) {
	if r.iol.useLocalPrefix {
		r.iol.localPrefix = runPrefix + "_" + strconv.Itoa(fnIdx) + "_"
		r.iol.localCounter = 0
		r.rc.blocks.resetHoistedCollected()
		r.rc.lastRejectedVar = nil
	}
	r.iol.clearStaleVars()
	r.t.currentScript = r.t.GetEnclosingScript(candFn)
	r.t.ClearScriptState()
	r.t.TraverseBranch(candFn, candFn.GetParent())
}

func (r *iolFnRunner) close() {
	r.iol.useLocalPrefix = false
	r.t.PopScope()
}

// Java: InlineObjectLiterals.process(Node, Node)
func (o *InlineObjectLiterals) Process(externs rhino.Node, root rhino.Node) {
	if root.IsNil() {
		return
	}
	s := &o.scanState
	s.reset()
	s.scan(root)
	if s.candidateBloom == 0 && !s.hasInvalidDecl {
		return
	}

	fastMode := o.compiler.GetOptions().isFastMode()
	if !s.hasNonFnCandidate && !s.hasInvalidDecl {
		s.collectCandFns(root)
		if len(s.candFns) == 0 {
			return
		}
		if fastMode || len(s.candFns) == 1 || !s.anyCandFnInHoistedFn {
			var runPrefix string
			numWorkers := 1
			if fastMode {
				runPrefix = o.safeNameIdSupplier()
				numWorkers = parallelThreadCount(o.compiler, len(s.candFns))
			}
			if numWorkers <= 1 {
				r := newIOLFnRunner(o, s, root, fastMode)
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
							worker := &InlineObjectLiterals{
								compiler:           rc,
								safeNameIdSupplier: o.safeNameIdSupplier,
							}
							r = newIOLFnRunner(worker, s, root, true)
						}
						r.traverseCandFn(s.candFns[fnIdx], runPrefix, fnIdx)
					}, func() {
						if r.t != nil {
							r.close()
						}
					}
			}, false)
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

// Java: InlineObjectLiterals.InliningBehavior.isInlinableObject(List<Reference>)
func (o *InlineObjectLiterals) isInlinableObject(refs []*Reference) bool {
	ret := false
	o.validProps.clear()
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

			propID := nodeStringID(parent)
			if !o.validProps.contains(propID) {
				if NodeUtil_IsNameDeclOrSimpleAssignLhs(parent, grandparent) {
					o.validProps.add(propID)
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
			o.validProps.add(nodeStringID(member.Node()))

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
