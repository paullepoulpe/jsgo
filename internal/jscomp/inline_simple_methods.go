// Port of com/google/javascript/jscomp/InlineSimpleMethods.java @ closure-compiler v20261021

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
	externPropBits  inlineBitSet
	externPropsInit bool
	props           map[uint32]ismPropSlot
	hasDefs         bool

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
	props := m.ensureProps()
	slot := props[propID]
	if slot.nonInlineable {
		return
	}
	slot.nonInlineable = true
	slot.extraDefs = nil
	props[propID] = slot
}

// Java: InlineSimpleMethods.process(Node, Node)
func (m *InlineSimpleMethods) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckStateMsg1(m.compiler.GetLifeCycleStage().IsNormalized(), "%s", m.compiler.GetLifeCycleStage())
	jutil.CheckState(!m.hasDefs)
	jutil.CheckState(!externs.IsNil())

	m.ensureProps()
	defer func() {
		m.props = nil
		m.externPropBits = inlineBitSet{}
		m.externPropsInit = false
		m.hasDefs = false
	}()

	NodeTraversal_TraverseRootsVisitor(m.compiler, externs, root, &ast.Visitor[*NodeTraversal]{
		GetProp:   func(_ *NodeTraversal, x ast.GetProp) { m.gatherNormalGet(x.Node()) },
		GetElem:   func(_ *NodeTraversal, x ast.GetElem) { m.gatherNormalGet(x.Node()) },
		ObjectLit: func(_ *NodeTraversal, x ast.ObjectLit) { m.gatherMemberContainer(x.Node()) },
		ClassBody: func(_ *NodeTraversal, x ast.ClassBody) { m.gatherMemberContainer(x.Node()) },
		Call:      m.gatherCall,
	})
	if m.hasDefs {
		NodeTraversal_TraverseVisitor(m.compiler, root, &ast.Visitor[*NodeTraversal]{
			Call: m.inlineCall,
		})
	}
}

// inlineCall checks whether a method call is a candidate for inlining and
// performs the inlining if so.
//
// Java: InlineSimpleMethods.InlineTrivialAccessors.visit(NodeTraversal, Node, Node, String)
func (m *InlineSimpleMethods) inlineCall(t *NodeTraversal, call ast.Call) {
	callee, ok := call.Callee().AsGetProp()
	if !ok {
		return
	}
	propID := callee.Node().GetStringID()
	slot, ok := m.props[propID]
	if !ok || slot.nonInlineable || !slot.hasFirst {
		return
	}

	firstDefinition := slot.firstDef
	if !slot.checked {
		if (len(slot.extraDefs) > 0 && !m.allDefinitionsEquivalent(firstDefinition, slot.extraDefs)) ||
			m.anyDefinitionsNotInlineable(firstDefinition, slot.extraDefs) {
			slot.nonInlineable = true
			slot.extraDefs = nil
			m.props[propID] = slot
			return
		}
		slot.checked = true
		m.props[propID] = slot
	}

	callNode := call.Node()

	// Do not inline if the callsite is a derived class calling a base method using `super.`
	if !m.argsMayHaveSideEffects(call) && !NodeUtil_ReferencesSuper(callNode) {
		if returned, hasRet := inlineSimpleMethods_returnedExpression(firstDefinition); hasRet {
			if inlineSimpleMethods_isPropertyTree(returned) && !firstDefinition.IsArrow() {
				m.inlinePropertyReturn(call, callee, returned.Node())
			} else if NodeUtil_IsLiteralValue(returned.Node(), false) &&
				!m.astAnalyzer.MayHaveSideEffects(callee.Node()) {
				m.inlineConstReturn(callNode, returned.Node())
			}
		} else if inlineSimpleMethods_isEmptyMethod(firstDefinition) &&
			!m.astAnalyzer.MayHaveSideEffects(callee.Node()) {
			m.inlineEmptyMethod(t, callNode.GetParent(), callNode)
		}
	}
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
func (m *InlineSimpleMethods) inlinePropertyReturn(call ast.Call, callee ast.GetProp, returnedValue rhino.Node) {
	getProp := returnedValue.CloneTree()
	inlineSimpleMethods_replaceThis(getProp, callee.Object().Detach().Node())
	call.Node().ReplaceWith(getProp)
	m.compiler.ReportChangeToEnclosingScope(getProp)
}

// Java: InlineSimpleMethods.inlineConstReturn(Node, Node)
func (m *InlineSimpleMethods) inlineConstReturn(call rhino.Node, returnedValue rhino.Node) {
	retValue := returnedValue.CloneTree()
	call.ReplaceWith(retValue)
	m.compiler.ReportChangeToEnclosingScope(retValue)
}

// Java: InlineSimpleMethods.inlineEmptyMethod(NodeTraversal, Node, Node)
func (m *InlineSimpleMethods) inlineEmptyMethod(t *NodeTraversal, parent rhino.Node, call rhino.Node) {
	// If the return value of the method call is read,
	// replace it with "void 0". Otherwise, remove the call entirely.
	if NodeUtil_IsExprCall(parent) {
		parent.ReplaceWith(rhino.IR_Empty())
		NodeUtil_MarkFunctionsDeleted(parent, m.compiler)
	} else {
		call.ReplaceWith(NodeUtil_NewUndefinedNode(call))
		NodeUtil_MarkFunctionsDeleted(call, m.compiler)
	}
	t.ReportCodeChange()
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
	props := m.ensureProps()
	slot := props[propID]
	if slot.nonInlineable {
		return
	}
	if m.hasExternProp(propID) {
		slot.nonInlineable = true
		slot.extraDefs = nil
		props[propID] = slot
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

// Java: InlineSimpleMethods.GatherSignatures.visit (GETPROP / GETELEM)
func (m *InlineSimpleMethods) gatherNormalGet(n rhino.Node) {
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
		m.processPrototypeParent(parent)
	} else if assign, ok := ast.AsAssign(parent); ok && assign.Target().Is(n) {
		// Static methods of the form Foo.bar = function() {} or
		// Static methods of the form Foo.bar = baz (where baz is a function name).
		m.addPossibleSignature(propID, assign.Value().Node())
	}
}

// Java: InlineSimpleMethods.GatherSignatures.visit (OBJECTLIT / CLASS_MEMBERS)
func (m *InlineSimpleMethods) gatherMemberContainer(container rhino.Node) {
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

// Java: InlineSimpleMethods.GatherSignatures.visit (CALL)
func (m *InlineSimpleMethods) gatherCall(_ *NodeTraversal, call ast.Call) {
	// If a goog.reflect.objectProperty is used for a method's name, we can't assume that the
	// method can be safely inlined.
	callee := call.Callee().Node()
	if firstArg, ok := call.FirstArg(); ok && m.codingConvention.IsPropertyRenameFunction(callee) {
		// Other code guarantees that getSecondChild() is a STRINGLIT
		m.markNonInlineable(firstArg.Node().GetStringID())
	}
}

// Java: InlineSimpleMethods.GatherSignatures.processPrototypeParent(Node)
func (m *InlineSimpleMethods) processPrototypeParent(n rhino.Node) {
	propID, hasName := ismStaticPropertyID(n)
	if !hasName {
		return
	}
	if assign, isAssign := ast.AsAssign(n.GetGrandparent()); isAssign {
		m.addPossibleSignature(propID, assign.Value().Node())
	}
}
