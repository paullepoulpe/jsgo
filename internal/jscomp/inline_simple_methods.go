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
)

type ismPropSlot struct {
	defs          []ast.Function
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
	externProperties *jutil.LinkedHashSet[string]
	props            map[string]*ismPropSlot
	hasDefs          bool

	compiler         AbstractCompiler
	astAnalyzer      *AstAnalyzer
	codingConvention CodingConvention
}

// Java: InlineSimpleMethods(AbstractCompiler)
func NewInlineSimpleMethods(compiler AbstractCompiler) *InlineSimpleMethods {
	return &InlineSimpleMethods{
		externProperties: compiler.GetExternProperties(),
		compiler:         compiler,
		astAnalyzer:      compiler.GetAstAnalyzer(),
		codingConvention: compiler.GetCodingConvention(),
	}
}

func (m *InlineSimpleMethods) ensureProps() map[string]*ismPropSlot {
	if m.props == nil {
		m.props = make(map[string]*ismPropSlot)
	}
	return m.props
}

func (m *InlineSimpleMethods) markNonInlineable(name string) {
	props := m.ensureProps()
	slot := props[name]
	if slot == nil {
		slot = &ismPropSlot{}
		props[name] = slot
	}
	slot.nonInlineable = true
}

// Java: InlineSimpleMethods.process(Node, Node)
func (m *InlineSimpleMethods) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckStateMsg1(m.compiler.GetLifeCycleStage().IsNormalized(), "%s", m.compiler.GetLifeCycleStage())
	jutil.CheckState(!m.hasDefs)
	jutil.CheckState(!externs.IsNil())

	m.ensureProps()
	defer func() {
		m.props = nil
		m.hasDefs = false
	}()

	NodeTraversal_TraverseRootsVisitor(m.compiler, externs, root, &ast.Visitor[*NodeTraversal]{
		GetProp:   func(_ *NodeTraversal, x ast.GetProp) { m.gatherNormalGet(x.AsNormalGet()) },
		GetElem:   func(_ *NodeTraversal, x ast.GetElem) { m.gatherNormalGet(x.AsNormalGet()) },
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
	slot := m.props[callee.Property()]
	if slot == nil || slot.nonInlineable || len(slot.defs) == 0 {
		return
	}

	defs := slot.defs
	if !slot.checked {
		if (len(defs) > 1 && !m.allDefinitionsEquivalent(defs)) || m.anyDefinitionsNotInlineable(defs) {
			slot.nonInlineable = true
			return
		}
		slot.checked = true
	}

	firstDefinition := defs[0]
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
func (m *InlineSimpleMethods) allDefinitionsEquivalent(definitions []ast.Function) bool {
	first := definitions[0].Node()
	for _, fn := range definitions[1:] {
		if !m.compiler.AreNodesEqualForInlining(first, fn.Node()) {
			return false
		}
	}
	return true
}

// Java: InlineSimpleMethods.anyDefinitionsNotInlineable(Set<Node>)
func (m *InlineSimpleMethods) anyDefinitionsNotInlineable(definitions []ast.Function) bool {
	for _, fn := range definitions {
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
		if jsDocInfo != nil && jsDocInfo.IsNoInline() {
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
func (m *InlineSimpleMethods) addPossibleSignature(name string, node rhino.Node) {
	if fn, ok := ast.AsFunction(node); ok {
		m.addSignature(name, fn)
	} else {
		m.markNonInlineable(name)
	}
}

// Java: InlineSimpleMethods.addSignature(String, Node)
func (m *InlineSimpleMethods) addSignature(name string, function ast.Function) {
	props := m.ensureProps()
	slot := props[name]
	if slot != nil && slot.nonInlineable {
		return
	}
	if m.externProperties != nil && m.externProperties.Contains(name) {
		if slot == nil {
			slot = &ismPropSlot{}
			props[name] = slot
		}
		slot.nonInlineable = true
		return
	}
	if slot == nil {
		slot = &ismPropSlot{defs: []ast.Function{function}}
		props[name] = slot
		m.hasDefs = true
		return
	}
	if slices.Contains(slot.defs, function) {
		return
	}
	slot.defs = append(slot.defs, function)
	m.hasDefs = true
}

// Java: InlineSimpleMethods.GatherSignatures.visit (GETPROP / GETELEM)
func (m *InlineSimpleMethods) gatherNormalGet(get ast.NormalGet) {
	n := get.Node()
	parent := n.GetParent()
	pt := parent.GetToken()
	if pt != rhino.Token_ASSIGN && pt != rhino.Token_GETPROP && pt != rhino.Token_GETELEM {
		return
	}
	// TODO(b/251573710): Handle ES2015+ language features that can reference properties
	name, ok := get.StaticPropertyName()
	if !ok {
		return
	}

	if name == "prototype" {
		m.processPrototypeParent(parent)
	} else if assign, ok := ast.AsAssign(parent); ok && assign.Target().Is(n) {
		// Static methods of the form Foo.bar = function() {} or
		// Static methods of the form Foo.bar = baz (where baz is a function name).
		m.addPossibleSignature(name, assign.Value().Node())
	}
}

// Java: InlineSimpleMethods.GatherSignatures.visit (OBJECTLIT / CLASS_MEMBERS)
func (m *InlineSimpleMethods) gatherMemberContainer(container rhino.Node) {
	for key := container.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		switch key.GetToken() {
		case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_STRING_KEY:
			m.addPossibleSignature(key.GetString(), key.GetFirstChild())
		case rhino.Token_SETTER_DEF, rhino.Token_GETTER_DEF:
			m.markNonInlineable(key.GetString())
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
		m.markNonInlineable(firstArg.Node().GetString())
	}
}

// Java: InlineSimpleMethods.GatherSignatures.processPrototypeParent(Node)
func (m *InlineSimpleMethods) processPrototypeParent(n rhino.Node) {
	get, ok := ast.AsNormalGet(n)
	if !ok {
		return
	}
	name, hasName := get.StaticPropertyName()
	if !hasName {
		return
	}
	if assign, isAssign := ast.AsAssign(n.GetGrandparent()); isAssign {
		m.addPossibleSignature(name, assign.Value().Node())
	}
}
