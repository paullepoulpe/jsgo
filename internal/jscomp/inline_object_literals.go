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
	behavior *InlineObjectLiterals_InliningBehavior,
	varFilter func(*Var) bool,
) *ReferenceCollector {
	rc := NewReferenceCollectorWithVarFilter(o.compiler, behavior, sc, varFilter)
	rc.enableLocalReclamation(func() bool {
		clear(behavior.staleVars)
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
					behavior := newInlineObjectLiteralsBehavior(o)
					behavior.scanner = &s
					dummyGlobalScope := newIOLDummyGlobalScope(root)
					for {
						fnIdx := int(nextIdx.Add(1) - 1)
						if fnIdx >= len(s.candFns) {
							break
						}
						candFn := s.candFns[fnIdx]
						localCounter := 0
						prefix := runPrefix + "_" + strconv.Itoa(fnIdx) + "_"
						clear(behavior.staleVars)
						behavior.idSupplier = func() string {
							id := prefix + strconv.Itoa(localCounter)
							localCounter++
							return id
						}
						behavior.reportChange = func(n rhino.Node) {
							if n.IsNil() {
								return
							}
							if n.IsScript() {
								changedPerFn[fnIdx] = append(changedPerFn[fnIdx], n)
							} else if cs := ChangeTracker_GetEnclosingChangeScopeRoot(n.GetParent()); !cs.IsNil() {
								changedPerFn[fnIdx] = append(changedPerFn[fnIdx], cs)
							}
						}
						rc := o.newFilteredCollector(sc, behavior, varFilter)
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
			behavior := newInlineObjectLiteralsBehavior(o)
			behavior.scanner = &s
			dummyGlobalScope := newIOLDummyGlobalScope(root)
			if fastMode && len(s.candFns) >= 1 {
				runPrefix := o.safeNameIdSupplier()
				for fnIdx, candFn := range s.candFns {
					localCounter := 0
					prefix := runPrefix + "_" + strconv.Itoa(fnIdx) + "_"
					clear(behavior.staleVars)
					behavior.idSupplier = func() string {
						id := prefix + strconv.Itoa(localCounter)
						localCounter++
						return id
					}
					rc := o.newFilteredCollector(sc, behavior, varFilter)
					traverseIOLCandFn(rc, root, dummyGlobalScope, candFn)
				}
				return
			}
			callback := o.newFilteredCollector(sc, behavior, varFilter)
			tb := callback.builderVal()
			t := newNodeTraversal(&tb)
			t.currentNode = root
			t.PushScopeOfScope(dummyGlobalScope)
			for _, candFn := range s.candFns {
				clear(behavior.staleVars)
				t.currentScript = t.GetEnclosingScript(candFn)
				t.ClearScriptState()
				t.TraverseBranch(candFn, candFn.GetParent())
			}
			t.PopScope()
			return
		}
	}

	sc := NewSyntacticScopeCreator(o.compiler)
	behavior := newInlineObjectLiteralsBehavior(o)
	behavior.scanner = &s
	o.newFilteredCollector(sc, behavior, varFilter).Process(externs, root)
}

// InlineObjectLiterals_InliningBehavior builds up information about nodes in
// each scope. When exiting the scope, inspects all variables in that scope,
// and inlines any that we can.
type InlineObjectLiterals_InliningBehavior struct {
	outer *InlineObjectLiterals

	// Variables whose reference information is out of sync with the AST.
	staleVars map[*Var]struct{}

	validProperties jutil.LinkedHashSet[string]

	scanner      *iolCandidateScanner
	idSupplier   func() string
	reportChange func(rhino.Node)
}

func newInlineObjectLiteralsBehavior(outer *InlineObjectLiterals) *InlineObjectLiterals_InliningBehavior {
	b := &InlineObjectLiterals_InliningBehavior{outer: outer}
	if outer != nil {
		b.idSupplier = outer.safeNameIdSupplier
		if outer.compiler != nil {
			b.reportChange = outer.compiler.ReportChangeToEnclosingScope
		}
	}
	return b
}

func (b *InlineObjectLiterals_InliningBehavior) addStaleVar(v *Var) {
	if v == nil {
		return
	}
	if b.staleVars == nil {
		b.staleVars = make(map[*Var]struct{}, 8)
	}
	b.staleVars[v] = struct{}{}
}

// Java: InlineObjectLiterals.InliningBehavior.afterExitScope(NodeTraversal, ReferenceMap)
func (b *InlineObjectLiterals_InliningBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
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
	if b.scanner != nil {
		scopeBloom := b.scanner.candidateScopeBlooms[scope.GetRootNode()]
		if scopeBloom == 0 || (scope.varBloom&scopeBloom) == 0 {
			return
		}
	}
	for _, v := range scope.varList {
		if v == nil || b.isVarInlineForbidden(v) {
			continue
		}
		referenceInfo := referenceMap.GetReferences(v)
		if referenceInfo == nil {
			continue
		}
		if b.isInlinableObject(referenceInfo.references) {
			b.addStaleVar(v)
			init := referenceInfo.GetInitializingReference()
			b.splitObject(v, init, referenceInfo)
		}
	}
}

// Java: InlineObjectLiterals.InliningBehavior.recordStaleVarReferencesInTree(Node, Scope)
func (b *InlineObjectLiterals_InliningBehavior) recordStaleVarReferencesInTree(root rhino.Node, scope *Scope) {
	if name, ok := ast.AsName(root); ok {
		b.addStaleVar(scope.GetVar(name.Ident()))
	}
	if !root.IsFunction() {
		for c := root.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			b.recordStaleVarReferencesInTree(c, scope)
		}
	}
}

// Java: InlineObjectLiterals.InliningBehavior.isVarInlineForbidden(Var)
func (b *InlineObjectLiterals_InliningBehavior) isVarInlineForbidden(var_ *Var) bool {
	_, stale := b.staleVars[var_]
	return var_.IsGlobal() ||
		var_.IsExtern() ||
		b.outer.compiler.GetCodingConvention().IsExportedWithLocal(var_.GetName(), true /* local= */) ||
		b.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(var_.GetNameNode()) ||
		stale
}

// Java: InlineObjectLiterals.InliningBehavior.isInlinableObject(List<Reference>)
func (b *InlineObjectLiterals_InliningBehavior) isInlinableObject(refs []*Reference) bool {
	ret := false
	validProperties := &b.validProperties
	if !validProperties.IsEmpty() {
		validProperties.Clear()
	}
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

			propName := gp.Property()
			if !validProperties.Contains(propName) {
				if NodeUtil_IsNameDeclOrSimpleAssignLhs(parent, grandparent) {
					validProperties.Add(propName)
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
			validProperties.Add(pkm.Key())

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
func (b *InlineObjectLiterals_InliningBehavior) computeVarList(referenceInfo *ReferenceCollection) *jutil.LinkedHashMap[string, string] {
	varmap := jutil.NewLinkedHashMap[string, string]()

	for _, ref := range referenceInfo.references {
		if ref.IsLvalue() || ref.IsInitializingDeclaration() {
			val := ref.GetAssignedValue()
			if !val.IsNil() {
				for member := range ast.MustObjectLit(val).Members() {
					varname := member.Node().GetString()
					if varmap.ContainsKey(varname) {
						continue
					}
					varIdent := varname
					if !rhino.TokenStream_IsJSIdentifier(varname) {
						varIdent = InlineObjectLiterals_STRING_KEY_IDENTIFIER
					}
					varmap.Put(varname, InlineObjectLiterals_VAR_PREFIX+varIdent+"_"+b.idSupplier())
				}
			}
		} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
			// Uninitialized declaration; no property keys to register.
		} else {
			gp := ast.MustNormalOrOptChainGetProp(ref.GetParent())
			varname := gp.Property()
			if !varmap.ContainsKey(varname) {
				varmap.Put(varname, InlineObjectLiterals_VAR_PREFIX+varname+"_"+b.idSupplier())
			}
		}
	}

	return varmap
}

// Java: InlineObjectLiterals.InliningBehavior.fillInitialValues(Reference, Map<String, Node>)
func (b *InlineObjectLiterals_InliningBehavior) fillInitialValues(init *Reference, initvals *jutil.LinkedHashMap[string, rhino.Node]) {
	objLit := ast.MustObjectLit(init.GetAssignedValue())
	for member := range objLit.Members() {
		key := member.Node()
		initvals.Put(key.GetString(), key.RemoveFirstChild())
	}
}

// replaceAssignmentExpression replaces `x = {...}` with `t1=a, t2=b, true`.
//
// Java: InlineObjectLiterals.InliningBehavior.replaceAssignmentExpression(Var, Reference, Map<String, String>)
func (b *InlineObjectLiterals_InliningBehavior) replaceAssignmentExpression(v *Var, ref *Reference, varmap *jutil.LinkedHashMap[string, string]) {
	var nodes []rhino.Node
	val := ref.GetAssignedValue()
	b.recordStaleVarReferencesInTree(val, v.GetScope().(*Scope))
	objLit := ast.MustObjectLit(val)
	all := jutil.NewLinkedHashSet(varmap.KeySlice()...)
	for member := range objLit.Members() {
		keyNode := member.Node()
		propName := keyNode.GetString()
		value := keyNode.RemoveFirstChild()
		nodes = append(nodes, rhino.IR_Assign(rhino.IR_Name(varmap.Get(propName)), value))
		all.Remove(propName)
	}

	for propName := range all.All() {
		nodes = append(nodes, rhino.IR_Assign(rhino.IR_Name(varmap.Get(propName)), NodeUtil_NewUndefinedNode(rhino.NilNode)))
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

	replace := ref.GetParent()
	replacement.SrcrefTreeIfMissing(replace)

	if NodeUtil_IsNameDeclaration(replace) {
		replace.ReplaceWith(NodeUtil_NewExpr(replacement))
	} else {
		replace.ReplaceWith(replacement)
	}
}

// Java: InlineObjectLiterals.InliningBehavior.splitObject(Var, Reference, ReferenceCollection)
func (b *InlineObjectLiterals_InliningBehavior) splitObject(v *Var, init *Reference, referenceInfo *ReferenceCollection) {
	varmap := b.computeVarList(referenceInfo)

	initvals := jutil.NewLinkedHashMap[string, rhino.Node]()
	var vnode rhino.Node
	defined := referenceInfo.IsWellDefined() && NodeUtil_IsNameDeclaration(init.GetParent())
	if defined {
		vnode = init.GetParent()
		b.fillInitialValues(init, initvals)
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

	for key, value := range varmap.All() {
		val := initvals.Get(key)
		newVarNode := NodeUtil_NewVarNode(value, val)
		if val.IsNil() {
			newVarNode.SrcrefTreeIfMissing(vnode)
		} else {
			b.recordStaleVarReferencesInTree(val, v.GetScope().(*Scope))
		}
		newVarNode.InsertBefore(vnode)
		b.reportChange(vnode)
	}

	if defined {
		b.reportChange(vnode.GetParent())
		vnode.Detach()
	}

	for _, ref := range referenceInfo.references {
		if defined && ref == init {
			continue
		}
		b.reportChange(ref.GetNode())

		if ref.IsLvalue() {
			b.replaceAssignmentExpression(v, ref, varmap)
		} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
			ref.GetParent().Detach()
		} else {
			gp := ast.MustNormalOrOptChainGetProp(ref.GetParent())
			propName := gp.Property()
			jutil.CheckState(varmap.ContainsKey(propName))

			replacement := rhino.IR_Name(varmap.Get(propName)).SrcrefIfMissing(gp.Node())
			gp.Node().ReplaceWith(replacement)
		}
	}
}
