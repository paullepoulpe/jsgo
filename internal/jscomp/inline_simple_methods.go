// Port of com/google/javascript/jscomp/InlineSimpleMethods.java @ closure-compiler v20261028

/*
 * Copyright 2007 The Closure Compiler Authors.
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

var ismPrototypeID = uint32(strpool.Intern("prototype"))

type ismPropSlot struct {
	firstDef      ast.Function
	extraDefs     []ast.Function
	hasFirst      bool
	nonInlineable bool
	checked       bool
}

// InlineSimpleMethods inlines methods that take no arguments and have only a
// return statement returning a property. Because it works on method names
// rather than type inference, a method with multiple definitions will be
// inlined if each definition is identical.
//
//	A.prototype.foo = function() { return this.b; }
//	B.prototype.foo = function() { return this.b; }
//
// will inline foo, but
//
//	A.prototype.foo = function() { return this.b; }
//	B.prototype.foo = function() { return this.c; }
//
// will not.
//
// Declarations are not removed because we do not find all possible call
// sites. For examples, calls of the form foo["bar"] are not detected.
//
// This pass is not on by default because it is not safe in simple mode. If
// the prototype method is mutated and we don't detect that, inlining it is
// unsafe. We enable it whenever function inlining is enabled.
type InlineSimpleMethods struct {
	externPropsRaw  *jutil.LinkedHashSet[string]
	externPropBits  smallBitSet
	externPropsInit bool
	props           map[uint32]ismPropSlot
	hasDefs         bool
	// Go-only: the properties whose slot is non-inlineable (a slot never
	// becomes inlineable again), to skip the map on the many repeated marks.
	nonInlineableBits smallBitSet

	compiler         AbstractCompiler
	astAnalyzer      *AstAnalyzer
	codingConvention CodingConvention
}

// Java: InlineSimpleMethods(AbstractCompiler)
func NewInlineSimpleMethods(compiler AbstractCompiler) *InlineSimpleMethods {
	return &InlineSimpleMethods{
		externPropsRaw:   compiler.GetExternProperties(),
		compiler:         compiler,
		astAnalyzer:      compiler.GetAstAnalyzer(),
		codingConvention: compiler.GetCodingConvention(),
	}
}

func (m *InlineSimpleMethods) ensureProps() map[uint32]ismPropSlot {
	if m.props == nil {
		m.props = make(map[uint32]ismPropSlot)
	}
	return m.props
}

func (m *InlineSimpleMethods) hasExternProp(id uint32) bool {
	if !m.externPropsInit {
		m.externPropsInit = true
		if m.externPropsRaw != nil {
			for name := range m.externPropsRaw.All() {
				if pid := strpool.Lookup(name); pid != 0 {
					m.externPropBits.Set(int(pid))
				}
			}
		}
	}
	return m.externPropBits.Get(int(id))
}

func (m *InlineSimpleMethods) markNonInlineable(propID uint32) {
	if m.nonInlineableBits.Get(int(propID)) {
		return
	}
	props := m.ensureProps()
	slot := props[propID]
	if slot.nonInlineable {
		return
	}
	slot.nonInlineable = true
	slot.extraDefs = nil
	props[propID] = slot
	m.nonInlineableBits.Set(int(propID))
}

// Java: InlineSimpleMethods.process(Node, Node)
func (m *InlineSimpleMethods) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckStateMsg1(m.compiler.GetLifeCycleStage().IsNormalized(), "%s", m.compiler.GetLifeCycleStage())
	jutil.CheckState(!m.hasDefs)
	jutil.CheckState(!externs.IsNil())

	m.ensureProps()
	defer func() {
		m.props = nil
		m.nonInlineableBits = smallBitSet{}
		m.externPropBits = smallBitSet{}
		m.externPropsInit = false
		m.hasDefs = false
	}()

	if ix := propUseIndexFor(m.compiler, externs, root); ix != nil {
		m.processIndexed(ix, externs, root)
		return
	}

	NodeTraversal_TraverseRootsVisitor(m.compiler, externs, root, &ast.Visitor[*NodeTraversal]{
		GetProp:   func(_ *NodeTraversal, x ast.GetProp) { m.gatherNormalGet(x.Node()) },
		GetElem:   func(_ *NodeTraversal, x ast.GetElem) { m.gatherNormalGet(x.Node()) },
		ObjectLit: func(_ *NodeTraversal, x ast.ObjectLit) { m.gatherMemberContainer(x.Node()) },
		ClassBody: func(_ *NodeTraversal, x ast.ClassBody) { m.gatherMemberContainer(x.Node()) },
		Call:      m.gatherCall,
	})
	if m.hasDefs {
		NodeTraversal_TraverseVisitor(m.compiler, root, &ast.Visitor[*NodeTraversal]{
			Call: func(t *NodeTraversal, call ast.Call) { m.inlineCall(t.currentChangeScope, call) },
		})
	}
}

// processIndexed is Process with the property-use index (prop_use_index.go)
// in place of the two traversals.
//
// Go-only: what GatherSignatures' callbacks find is replayed in the order the
// traversal of externs and root visits them: the index ran the callbacks on
// each change scope (ismSignatureSink) and recorded their findings, or, where
// a finding depends on more than the scope, the node to run them on. The calls
// InlineTrivialAccessors may inline are those whose callee is a GETPROP of a
// property gathered as inlineable (a slot only becomes non-inlineable during
// the inlining), replayed in the traversal's post-order with the change scope
// the traversal would be in. Inlining a call never changes the tree after it
// in that order except at its parent: when the inlined call was the callee of
// another call, that call now has a GETPROP callee and the traversal, which
// visits it after the rest of its subtree, considers it too.
func (m *InlineSimpleMethods) processIndexed(ix *propUseIndex, externs, root rhino.Node) {
	rp := piReplay{compiler: m.compiler}
	defer rp.done()
	gather := func(_ *piScope, it piItem) {
		n := it.n
		if it.bits&piISMEv != 0 {
			if n.IsNil() {
				m.markNonInlineable(it.prop)
			} else {
				m.addSignature(it.prop, ast.MustFunction(n))
			}
			return
		}
		rp.cur = n
		switch n.GetToken() {
		case rhino.Token_GETPROP, rhino.Token_GETELEM:
			m.gatherNormalGet(n)
		case rhino.Token_OBJECTLIT, rhino.Token_CLASS_MEMBERS:
			m.gatherMemberContainer(n)
		}
	}
	ix.each(externs, piListISM, nil, gather)
	ix.each(root, piListISM, nil, gather)
	if !m.hasDefs {
		return
	}

	var candidates smallBitSet
	// order-independent: builds a set
	for propID, slot := range m.props {
		if slot.hasFirst && !slot.nonInlineable {
			candidates.Set(int(propID))
		}
	}
	if candidates.IsEmpty() {
		return
	}

	type pendingCall struct{ call, changeScope rhino.Node }
	var pending []pendingCall
	inline := func(changeScope, call rhino.Node) {
		rp.cur = call
		r := m.inlineCall(changeScope, ast.MustCall(call))
		if !r.IsNil() && r.IsGetProp() {
			if p := r.GetParent(); p.IsCall() && p.GetFirstChild() == r {
				pending = append(pending, pendingCall{p, changeScope})
			}
		}
	}
	// flush inlines the pending calls whose subtree ends before next.
	flush := func(next rhino.Node) {
		for len(pending) > 0 {
			p := pending[len(pending)-1]
			if !next.IsNil() && piIsAncestor(p.call, next) {
				return
			}
			pending = pending[:len(pending)-1]
			inline(p.changeScope, p.call)
		}
	}
	ix.each(root, piListCalls, &candidates, func(s *piScope, it piItem) {
		if len(pending) > 0 {
			flush(it.n)
		}
		inline(s.root, it.n)
	})
	flush(rhino.NilNode)
}

// piIsAncestor reports whether a is a proper ancestor of n.
func piIsAncestor(a, n rhino.Node) bool {
	for p := n.GetParent(); !p.IsNil(); p = p.GetParent() {
		if p == a {
			return true
		}
	}
	return false
}

// inlineCall checks whether a method call is a candidate for inlining and
// performs the inlining if so. changeScope is the change scope the call is
// in. It returns the node that replaced the call, or NilNode.
//
// Java: InlineSimpleMethods.InlineTrivialAccessors.visit(NodeTraversal, Node, Node, String)
func (m *InlineSimpleMethods) inlineCall(changeScope rhino.Node, call ast.Call) rhino.Node {
	callee, ok := call.Callee().AsGetProp()
	if !ok {
		return rhino.NilNode
	}
	propID := callee.Node().GetStringID()
	slot, ok := m.props[propID]
	if !ok || slot.nonInlineable || !slot.hasFirst {
		return rhino.NilNode
	}

	firstDefinition := slot.firstDef
	if !slot.checked {
		if (len(slot.extraDefs) > 0 && !m.allDefinitionsEquivalent(firstDefinition, slot.extraDefs)) ||
			m.anyDefinitionsNotInlineable(firstDefinition, slot.extraDefs) {
			slot.nonInlineable = true
			slot.extraDefs = nil
			m.props[propID] = slot
			m.nonInlineableBits.Set(int(propID))
			return rhino.NilNode
		}
		slot.checked = true
		m.props[propID] = slot
	}

	callNode := call.Node()

	// Do not inline if the callsite is a derived class calling a base method using `super.`
	if !m.argsMayHaveSideEffects(call) && !NodeUtil_ReferencesSuper(callNode) {
		if returned, hasRet := inlineSimpleMethods_returnedExpression(firstDefinition); hasRet {
			if inlineSimpleMethods_isPropertyTree(returned) && !firstDefinition.IsArrow() {
				return m.inlinePropertyReturn(call, callee, returned.Node())
			} else if NodeUtil_IsLiteralValue(returned.Node(), false) &&
				!m.astAnalyzer.MayHaveSideEffects(callee.Node()) {
				return m.inlineConstReturn(callNode, returned.Node())
			}
		} else if inlineSimpleMethods_isEmptyMethod(firstDefinition) &&
			!m.astAnalyzer.MayHaveSideEffects(callee.Node()) {
			return m.inlineEmptyMethod(changeScope, callNode.GetParent(), callNode)
		}
	}
	return rhino.NilNode
}

// inlineSimpleMethods_isPropertyTree reports whether expr is a getprop chain
// rooted at `this` (e.g. `this.a` or `this.a.b`).
//
// Java: InlineSimpleMethods.isPropertyTree(Node)
func inlineSimpleMethods_isPropertyTree(expr ast.Expr) bool {
	gp, ok := expr.AsGetProp()
	if !ok {
		return false
	}
	obj := gp.Object()
	return obj.Node().IsThis() || inlineSimpleMethods_isPropertyTree(obj)
}

// inlineSimpleMethods_replaceThis replaces the root `this` in a cloned
// property-access chain with replacement.
//
// Java: InlineSimpleMethods.replaceThis(Node, Node)
func inlineSimpleMethods_replaceThis(expectedGetprop rhino.Node, replacement rhino.Node) {
	leftChild := ast.MustGetProp(expectedGetprop).Object().Node()
	if leftChild.IsThis() {
		leftChild.ReplaceWith(replacement)
	} else {
		inlineSimpleMethods_replaceThis(leftChild, replacement)
	}
}

// Java: InlineSimpleMethods.returnedExpression(Node)
func inlineSimpleMethods_returnedExpression(fn ast.Function) (ast.Expr, bool) {
	block, ok := fn.Body().AsBlock()
	if !ok {
		return ast.Expr{}, false
	}
	stmt, ok := block.OnlyStatement()
	if !ok {
		return ast.Expr{}, false
	}
	ret, ok := stmt.AsReturn()
	if !ok {
		return ast.Expr{}, false
	}
	return ret.Value()
}

// Java: InlineSimpleMethods.isEmptyMethod(Node)
func inlineSimpleMethods_isEmptyMethod(fn ast.Function) bool {
	return NodeUtil_IsEmptyBlock(fn.Body().Node())
}

// Java: InlineSimpleMethods.allDefinitionsEquivalent(Set<Node>)
func (m *InlineSimpleMethods) allDefinitionsEquivalent(first ast.Function, extraDefs []ast.Function) bool {
	firstNode := first.Node()
	for _, fn := range extraDefs {
		if !m.compiler.AreNodesEqualForInlining(firstNode, fn.Node()) {
			return false
		}
	}
	return true
}

func (m *InlineSimpleMethods) isDefinitionNotInlineable(fn ast.Function) bool {
	if fn.IsAsync() || fn.IsGenerator() {
		// async and generator functions cannot be trivially inlined because they return a Promise
		// or generator object, not the return value directly. For now just don't inline them -
		// it's unclear it's worth the complexity.
		return true
	}
	if NodeUtil_FunctionParametersMayHaveSideEffects(fn.Node(), m.astAnalyzer) {
		// Methods with parameters that may have side effects or throw on destructuring are not
		// inlined to prevent altering runtime behavior.
		return true
	}
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fn.Node().GetParent())
	return jsDocInfo != nil && jsDocInfo.IsNoInline()
}

// Java: InlineSimpleMethods.anyDefinitionsNotInlineable(Set<Node>)
func (m *InlineSimpleMethods) anyDefinitionsNotInlineable(first ast.Function, extraDefs []ast.Function) bool {
	if m.isDefinitionNotInlineable(first) {
		return true
	}
	for _, fn := range extraDefs {
		if m.isDefinitionNotInlineable(fn) {
			return true
		}
	}
	return false
}

// Java: InlineSimpleMethods.inlinePropertyReturn(Node, Node)
func (m *InlineSimpleMethods) inlinePropertyReturn(call ast.Call, callee ast.GetProp, returnedValue rhino.Node) rhino.Node {
	getProp := returnedValue.CloneTree()
	inlineSimpleMethods_replaceThis(getProp, callee.Object().Detach().Node())
	call.Node().ReplaceWith(getProp)
	m.compiler.ReportChangeToEnclosingScope(getProp)
	return getProp
}

// Java: InlineSimpleMethods.inlineConstReturn(Node, Node)
func (m *InlineSimpleMethods) inlineConstReturn(call rhino.Node, returnedValue rhino.Node) rhino.Node {
	retValue := returnedValue.CloneTree()
	call.ReplaceWith(retValue)
	m.compiler.ReportChangeToEnclosingScope(retValue)
	return retValue
}

// changeScope is the traversal's current change scope (Java:
// t.reportCodeChange()).
//
// Java: InlineSimpleMethods.inlineEmptyMethod(NodeTraversal, Node, Node)
func (m *InlineSimpleMethods) inlineEmptyMethod(changeScope rhino.Node, parent rhino.Node, call rhino.Node) rhino.Node {
	// If the return value of the method call is read,
	// replace it with "void 0". Otherwise, remove the call entirely.
	var replacement rhino.Node
	if NodeUtil_IsExprCall(parent) {
		replacement = rhino.IR_Empty()
		parent.ReplaceWith(replacement)
		NodeUtil_MarkFunctionsDeleted(parent, m.compiler)
	} else {
		replacement = NodeUtil_NewUndefinedNode(call)
		call.ReplaceWith(replacement)
		NodeUtil_MarkFunctionsDeleted(call, m.compiler)
	}
	// t.reportCodeChange()
	if changeScope.IsNil() {
		panic(jutil.NewNullPointerException()) // checkNotNull(changeScope)
	}
	jutil.CheckStateMsg1(ChangeTracker_IsChangeScopeRoot(changeScope), "%s", changeScope)
	m.compiler.ReportChangeToChangeScope(changeScope)
	return replacement
}

// Java: InlineSimpleMethods.argsMayHaveSideEffects(Node)
func (m *InlineSimpleMethods) argsMayHaveSideEffects(call ast.Call) bool {
	for arg := range call.Args() {
		if m.astAnalyzer.MayHaveSideEffects(arg.Node()) {
			return true
		}
	}
	return false
}

// Java: InlineSimpleMethods.addPossibleSignature(String, Node)
func (m *InlineSimpleMethods) addPossibleSignature(propID uint32, node rhino.Node) {
	if fn, ok := ast.AsFunction(node); ok {
		m.addSignature(propID, fn)
	} else {
		m.markNonInlineable(propID)
	}
}

// Java: InlineSimpleMethods.addSignature(String, Node)
func (m *InlineSimpleMethods) addSignature(propID uint32, function ast.Function) {
	if m.nonInlineableBits.Get(int(propID)) {
		return
	}
	props := m.ensureProps()
	slot := props[propID]
	if slot.nonInlineable {
		return
	}
	if m.hasExternProp(propID) {
		slot.nonInlineable = true
		slot.extraDefs = nil
		props[propID] = slot
		m.nonInlineableBits.Set(int(propID))
		return
	}
	if !slot.hasFirst {
		slot.firstDef = function
		slot.hasFirst = true
		props[propID] = slot
		m.hasDefs = true
		return
	}
	if slot.firstDef == function || slices.Contains(slot.extraDefs, function) {
		return
	}
	slot.extraDefs = append(slot.extraDefs, function)
	props[propID] = slot
	m.hasDefs = true
}

func ismStaticPropertyID(n rhino.Node) (uint32, bool) {
	switch n.GetToken() {
	case rhino.Token_GETPROP:
		return n.GetStringID(), true
	case rhino.Token_GETELEM:
		last := n.GetLastChild()
		if last.IsStringLit() {
			return last.GetStringID(), true
		}
	}
	return 0, false
}

// ismSignatureSink receives what GatherSignatures finds: the run itself, or
// the property-use index, which records it for the runs that replay it
// (prop_use_index.go).
type ismSignatureSink interface {
	addPossibleSignature(propID uint32, node rhino.Node)
	markNonInlineable(propID uint32)
}

func (m *InlineSimpleMethods) gatherNormalGet(n rhino.Node) { ismGatherNormalGet(m, n) }

func (m *InlineSimpleMethods) gatherMemberContainer(container rhino.Node) {
	ismGatherMemberContainer(m, container)
}

// Java: InlineSimpleMethods.GatherSignatures.visit (CALL)
func (m *InlineSimpleMethods) gatherCall(_ *NodeTraversal, call ast.Call) {
	ismGatherCall(m, m.codingConvention, call)
}

// Java: InlineSimpleMethods.GatherSignatures.visit (GETPROP / GETELEM)
func ismGatherNormalGet(m ismSignatureSink, n rhino.Node) {
	parent := n.GetParent()
	pt := parent.GetToken()
	if pt != rhino.Token_ASSIGN && pt != rhino.Token_GETPROP && pt != rhino.Token_GETELEM {
		return
	}
	// TODO(b/251573710): Handle ES2015+ language features that can reference properties
	propID, ok := ismStaticPropertyID(n)
	if !ok {
		return
	}

	if propID == ismPrototypeID {
		ismProcessPrototypeParent(m, parent)
	} else if assign, ok := ast.AsAssign(parent); ok && assign.Target().Is(n) {
		// Static methods of the form Foo.bar = function() {} or
		// Static methods of the form Foo.bar = baz (where baz is a function name).
		m.addPossibleSignature(propID, assign.Value().Node())
	}
}

// Java: InlineSimpleMethods.GatherSignatures.visit (OBJECTLIT / CLASS_MEMBERS)
func ismGatherMemberContainer(m ismSignatureSink, container rhino.Node) {
	for key := container.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		switch key.GetToken() {
		case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_STRING_KEY:
			m.addPossibleSignature(key.GetStringID(), key.GetFirstChild())
		case rhino.Token_SETTER_DEF, rhino.Token_GETTER_DEF:
			m.markNonInlineable(key.GetStringID())
		case rhino.Token_COMPUTED_PROP, // complicated
			rhino.Token_OBJECT_SPREAD,
			rhino.Token_COMPUTED_FIELD_DEF,
			rhino.Token_BLOCK: // ES2022 static initialization blocks are never methods.
		default:
			panic(jutil.NewIllegalStateExceptionMsg("Unexpected " + container.GetToken().String() + " key: " + key.String()))
		}
	}
}

// ismMemberKeysKnown reports whether ismGatherMemberContainer accepts every
// key of container (it panics on another key).
func ismMemberKeysKnown(container rhino.Node) bool {
	for key := container.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		switch key.GetToken() {
		case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_STRING_KEY,
			rhino.Token_SETTER_DEF, rhino.Token_GETTER_DEF, rhino.Token_COMPUTED_PROP,
			rhino.Token_OBJECT_SPREAD, rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_BLOCK:
		default:
			return false
		}
	}
	return true
}

// Java: InlineSimpleMethods.GatherSignatures.visit (CALL)
func ismGatherCall(m ismSignatureSink, conv CodingConvention, call ast.Call) {
	// If a goog.reflect.objectProperty is used for a method's name, we can't assume that the
	// method can be safely inlined.
	callee := call.Callee().Node()
	if firstArg, ok := call.FirstArg(); ok && conv.IsPropertyRenameFunction(callee) {
		// Other code guarantees that getSecondChild() is a STRINGLIT
		m.markNonInlineable(firstArg.Node().GetStringID())
	}
}

// Java: InlineSimpleMethods.GatherSignatures.processPrototypeParent(Node)
func ismProcessPrototypeParent(m ismSignatureSink, n rhino.Node) {
	propID, hasName := ismStaticPropertyID(n)
	if !hasName {
		return
	}
	if assign, isAssign := ast.AsAssign(n.GetGrandparent()); isAssign {
		m.addPossibleSignature(propID, assign.Value().Node())
	}
}
