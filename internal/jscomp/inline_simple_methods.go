// Port of com/google/javascript/jscomp/InlineSimpleMethods.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	// List of property names that we know are not safe to inline
	// This includes:
	//   - extern properties (we can't see the actual method definition for externs)
	//   - non-method properties
	//   - methods with @noinline
	//   - methods with multiple, non-equivalent definitions
	nonInlineableProperties *jutil.LinkedHashSet[string]

	// Use a linked map here to keep the output deterministic.  Otherwise,
	// the choice of method bodies is random when multiple identical definitions
	// are found which causes problems in the source maps.
	methodDefinitions *jutil.SetMultimap[string, *rhino.Node]

	compiler    AbstractCompiler
	astAnalyzer *AstAnalyzer
}

// Java: InlineSimpleMethods(AbstractCompiler)
func NewInlineSimpleMethods(compiler AbstractCompiler) *InlineSimpleMethods {
	m := &InlineSimpleMethods{
		nonInlineableProperties: jutil.NewLinkedHashSet[string](),
		methodDefinitions:       jutil.NewSetMultimap[string, *rhino.Node](),
	}
	m.compiler = compiler
	m.astAnalyzer = compiler.GetAstAnalyzer()
	m.nonInlineableProperties.AddAll(compiler.GetExternProperties().ToSlice()...)
	return m
}

// Java: InlineSimpleMethods.process(Node, Node)
func (m *InlineSimpleMethods) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckStateMsg1(m.compiler.GetLifeCycleStage().IsNormalized(), "%s", m.compiler.GetLifeCycleStage())
	jutil.CheckState(m.methodDefinitions.IsEmpty())
	jutil.CheckState(externs != nil)

	NodeTraversal_TraverseRoots(m.compiler, NewInlineSimpleMethods_GatherSignatures(m), externs, root)
	NodeTraversal_Traverse(m.compiler, root, NewInlineSimpleMethods_InlineTrivialAccessors(m))
}

// InlineSimpleMethods_InlineTrivialAccessors: for each method call, see if
// it is a candidate for inlining, and do the inlining if so.
type InlineSimpleMethods_InlineTrivialAccessors struct {
	InvocationsCallbackBase
	outer *InlineSimpleMethods
}

// Java: InlineSimpleMethods.InlineTrivialAccessors()
func NewInlineSimpleMethods_InlineTrivialAccessors(outer *InlineSimpleMethods) *InlineSimpleMethods_InlineTrivialAccessors {
	a := &InlineSimpleMethods_InlineTrivialAccessors{outer: outer}
	a.InitInvocationsCallbackBase(a)
	return a
}

// Java: InlineSimpleMethods.InlineTrivialAccessors.visit(NodeTraversal, Node, Node, String)
func (a *InlineSimpleMethods_InlineTrivialAccessors) VisitWithCallName(t *NodeTraversal, callNode *rhino.Node, parent *rhino.Node, callName string) {
	m := a.outer
	if m.nonInlineableProperties.Contains(callName) {
		return
	}

	definitions := m.methodDefinitions.Get(callName)
	if definitions == nil || definitions.IsEmpty() {
		return
	}

	// Exit early if all method definitions are not equivalent, and mark this method as not
	// inlineable.
	// NOTE: we could also cache the 'good' result of this method having all equivalent
	// definitions and avoid recalculating later, but profile data suggests that's not too useful
	if definitions.Size() > 1 && !m.AllDefinitionsEquivalent(definitions) {
		m.nonInlineableProperties.Add(callName)
		return
	}
	// Exit early if any definitions are incompatible with the InlineSimpleMethods inlining
	// NOTE: we could also cache the 'good' result of these methods being inlineable to
	// avoid recalculating later, but profile data suggests that's not too useful.
	if m.AnyDefinitionsNotInlineable(definitions) {
		m.nonInlineableProperties.Add(callName)
		return
	}

	// Do check of arity, complexity, and consistency in what we think is
	// the order from least to most complex
	var firstDefinition *rhino.Node
	for d := range definitions.All() {
		firstDefinition = d
		break
	}

	// Do not inline if the callsite is a derived class calling a base method using `super.`
	if !m.ArgsMayHaveSideEffects(callNode) && !NodeUtil_ReferencesSuper(callNode) {
		// Verify this is a trivial return
		returned := InlineSimpleMethods_ReturnedExpression(firstDefinition)
		if returned != nil {
			if InlineSimpleMethods_IsPropertyTree(returned) && !firstDefinition.IsArrowFunction() {
				m.InlinePropertyReturn(callNode, returned)
			} else if NodeUtil_IsLiteralValue(returned, false) &&
				!m.astAnalyzer.MayHaveSideEffects(callNode.GetFirstChild()) {
				m.InlineConstReturn(callNode, returned)
			}
		} else if InlineSimpleMethods_IsEmptyMethod(firstDefinition) &&
			!m.astAnalyzer.MayHaveSideEffects(callNode.GetFirstChild()) {
			m.InlineEmptyMethod(t, parent, callNode)
		}
	}
}

// InlineSimpleMethods_IsPropertyTree returns true if the provided node is a
// getprop for which the left child is this or a valid property tree and for
// which the right side is a string.
//
// Java: InlineSimpleMethods.isPropertyTree(Node)
func InlineSimpleMethods_IsPropertyTree(expectedGetprop *rhino.Node) bool {
	if !expectedGetprop.IsGetProp() {
		return false
	}

	getpropLhs := expectedGetprop.GetFirstChild()
	return getpropLhs.IsThis() || InlineSimpleMethods_IsPropertyTree(getpropLhs)
}

// InlineSimpleMethods_ReplaceThis finds the occurrence of "this" in the
// provided property tree and replaces it with replacement.
//
// Java: InlineSimpleMethods.replaceThis(Node, Node)
func InlineSimpleMethods_ReplaceThis(expectedGetprop *rhino.Node, replacement *rhino.Node) {
	leftChild := expectedGetprop.GetFirstChild()
	if leftChild.IsThis() {
		leftChild.ReplaceWith(replacement)
	} else {
		InlineSimpleMethods_ReplaceThis(leftChild, replacement)
	}
}

// InlineSimpleMethods_ReturnedExpression returns the node that represents
// the expression returned by the method, given a FUNCTION node.
//
// Java: InlineSimpleMethods.returnedExpression(Node)
func InlineSimpleMethods_ReturnedExpression(fn *rhino.Node) *rhino.Node {
	expectedBlock := NodeUtil_GetFunctionBody(fn)
	if !expectedBlock.HasOneChild() {
		return nil
	}

	expectedReturn := expectedBlock.GetFirstChild()
	if !expectedReturn.IsReturn() {
		return nil
	}

	if !expectedReturn.HasOneChild() {
		return nil
	}

	return expectedReturn.GetOnlyChild()
}

// InlineSimpleMethods_IsEmptyMethod returns whether the given FUNCTION node
// is an empty method definition.
//
// Must be private, or moved to NodeUtil.
//
// Java: InlineSimpleMethods.isEmptyMethod(Node)
func InlineSimpleMethods_IsEmptyMethod(fn *rhino.Node) bool {
	return NodeUtil_IsEmptyBlock(NodeUtil_GetFunctionBody(fn))
}

// AllDefinitionsEquivalent verifies, given a set of method definitions, that
// they are the same.
//
// Java: InlineSimpleMethods.allDefinitionsEquivalent(Set<Node>)
func (m *InlineSimpleMethods) AllDefinitionsEquivalent(definitions *jutil.ObjectLinkedHashSet[*rhino.Node]) bool {
	var first *rhino.Node
	for n := range definitions.All() {
		if first == nil {
			first = n
		} else if !m.compiler.AreNodesEqualForInlining(first, n) {
			return false
		} // else continue
	}
	return true
}

// Java: InlineSimpleMethods.anyDefinitionsNotInlineable(Set<Node>)
func (m *InlineSimpleMethods) AnyDefinitionsNotInlineable(definitions *jutil.ObjectLinkedHashSet[*rhino.Node]) bool {
	for n := range definitions.All() {
		if n.IsAsyncFunction() || n.IsGeneratorFunction() {
			// async and generator functions cannot be trivially inlined because they return a Promise
			// or generator object, not the return value directly. For now just don't inline them -
			// it's unclear it's worth the complexity.
			return true
		}
		jsDocInfo := NodeUtil_GetBestJSDocInfo(n.GetParent())
		if jsDocInfo != nil && jsDocInfo.IsNoInline() {
			return true
		}
	}
	return false
}

// InlinePropertyReturn replaces the provided method call with the tree
// specified in returnedValue.
//
// Parse tree of a call is name call getprop obj string.
//
// Java: InlineSimpleMethods.inlinePropertyReturn(Node, Node)
func (m *InlineSimpleMethods) InlinePropertyReturn(call *rhino.Node, returnedValue *rhino.Node) {
	getProp := returnedValue.CloneTree()
	InlineSimpleMethods_ReplaceThis(getProp, call.GetFirstChild().RemoveFirstChild())
	call.ReplaceWith(getProp)
	m.compiler.ReportChangeToEnclosingScope(getProp)
}

// InlineConstReturn replaces the provided object and its method call with
// the tree specified in returnedValue. Should be called only if the object
// reference has no side effects.
//
// Java: InlineSimpleMethods.inlineConstReturn(Node, Node)
func (m *InlineSimpleMethods) InlineConstReturn(call *rhino.Node, returnedValue *rhino.Node) {
	retValue := returnedValue.CloneTree()
	call.ReplaceWith(retValue)
	m.compiler.ReportChangeToEnclosingScope(retValue)
}

// InlineEmptyMethod removes the provided object and its method call.
//
// Java: InlineSimpleMethods.inlineEmptyMethod(NodeTraversal, Node, Node)
func (m *InlineSimpleMethods) InlineEmptyMethod(t *NodeTraversal, parent *rhino.Node, call *rhino.Node) {
	// If the return value of the method call is read,
	// replace it with "void 0". Otherwise, remove the call entirely.

	if NodeUtil_IsExprCall(parent) {
		parent.ReplaceWith(rhino.IR_Empty())
		NodeUtil_MarkFunctionsDeleted(parent, m.compiler)
	} else {
		srcLocation := call
		call.ReplaceWith(NodeUtil_NewUndefinedNode(srcLocation))
		NodeUtil_MarkFunctionsDeleted(call, m.compiler)
	}
	t.ReportCodeChange()
}

// ArgsMayHaveSideEffects checks whether the given method call's arguments
// have side effects.
//
// call is the call node of a method invocation.
//
// Java: InlineSimpleMethods.argsMayHaveSideEffects(Node)
func (m *InlineSimpleMethods) ArgsMayHaveSideEffects(call *rhino.Node) bool {
	for currentChild := call.GetSecondChild(); currentChild != nil; currentChild = currentChild.GetNext() {
		if m.astAnalyzer.MayHaveSideEffects(currentChild) {
			return true
		}
	}

	return false
}

// AddPossibleSignature adds a node that may represent a function signature
// (if it's a function itself or the name of a function).
//
// Java: InlineSimpleMethods.addPossibleSignature(String, Node)
func (m *InlineSimpleMethods) AddPossibleSignature(name string, node *rhino.Node) {
	if node != nil && node.IsFunction() {
		// The node we're looking at is a function, so we can add it directly
		m.AddSignature(name, node)
	} else {
		m.nonInlineableProperties.Add(name)
	}
}

// Java: InlineSimpleMethods.addSignature(String, Node)
func (m *InlineSimpleMethods) AddSignature(name string, function *rhino.Node) {
	if m.nonInlineableProperties.Contains(name) {
		return
	}

	m.methodDefinitions.Put(name, function)
}

// InlineSimpleMethods_GatherSignatures gathers signatures from the source to
// be compiled.
type InlineSimpleMethods_GatherSignatures struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineSimpleMethods
}

// Java: InlineSimpleMethods.GatherSignatures()
func NewInlineSimpleMethods_GatherSignatures(outer *InlineSimpleMethods) *InlineSimpleMethods_GatherSignatures {
	return &InlineSimpleMethods_GatherSignatures{outer: outer}
}

// Java: InlineSimpleMethods.GatherSignatures.visit(NodeTraversal, Node, Node)
func (g *InlineSimpleMethods_GatherSignatures) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	m := g.outer
	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		// TODO(b/251573710): Handle ES2015+ language features that can reference properties
		name := InlineSimpleMethods_GetPropName(n)
		if name == nil {
			return
		}

		if *name == "prototype" {
			g.ProcessPrototypeParent(parent)
		} else {
			// Static methods of the form Foo.bar = function() {} or
			// Static methods of the form Foo.bar = baz (where baz is a
			// function name). Parse tree looks like:
			// assign                 <- parent
			//      getprop           <- n
			//          name Foo
			//          string bar
			//      function or name  <- n.getNext()
			if parent.IsAssign() && n.IsFirstChildOf(parent) {
				m.AddPossibleSignature(*name, n.GetNext())
			}
		}
	case rhino.Token_OBJECTLIT, rhino.Token_CLASS_MEMBERS:
		for key := n.GetFirstChild(); key != nil; key = key.GetNext() {
			switch key.GetToken() {
			case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_STRING_KEY:
				m.AddPossibleSignature(key.GetString(), key.GetFirstChild())
			case rhino.Token_SETTER_DEF, rhino.Token_GETTER_DEF:
				m.nonInlineableProperties.Add(key.GetString())
			case rhino.Token_COMPUTED_PROP, // complicated
				rhino.Token_OBJECT_SPREAD,
				rhino.Token_COMPUTED_FIELD_DEF,
				rhino.Token_BLOCK: // ES2022 static initialization blocks are never methods.
			default:
				panic(jutil.NewIllegalStateExceptionMsg("Unexpected " + n.GetToken().String() + " key: " + key.String()))
			}
		}
	case rhino.Token_CALL:
		// If a goog.reflect.objectProperty is used for a method's name, we can't assume that the
		// method can be safely inlined.
		if m.compiler.GetCodingConvention().IsPropertyRenameFunction(n.GetFirstChild()) {
			// Other code guarantees that getSecondChild() is a STRINGLIT
			m.nonInlineableProperties.Add(n.GetSecondChild().GetString())
		}
	default:
	}
}

// ProcessPrototypeParent processes the parent of a GETPROP prototype, which
// can either be another GETPROP (in the case of Foo.prototype.bar), or can
// be an assignment (in the case of Foo.prototype = ...).
//
// Java: InlineSimpleMethods.GatherSignatures.processPrototypeParent(Node)
func (g *InlineSimpleMethods_GatherSignatures) ProcessPrototypeParent(n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		// Foo.prototype.getBar = function() { ... } or
		// Foo.prototype.getBar = getBaz (where getBaz is a function)
		// parse tree looks like:
		// assign                          <- parent
		//     getprop                     <- n
		//         getprop
		//             name Foo
		//             string prototype
		//         string getBar
		//     function or name            <- assignee
		grandparent := n.GetGrandparent()
		name := InlineSimpleMethods_GetPropName(n)
		if name != nil && grandparent.IsAssign() {
			assignee := grandparent.GetSecondChild()
			g.outer.AddPossibleSignature(*name, assignee)
		}
	default:
	}
}

// InlineSimpleMethods_GetPropName returns the property name, or nil (Java
// null).
//
// Java: InlineSimpleMethods.getPropName(Node)
func InlineSimpleMethods_GetPropName(getPropElem *rhino.Node) *string {
	if getPropElem.IsGetProp() {
		return jutil.StringPtr(getPropElem.GetString())
	} else if getPropElem.GetSecondChild().IsStringLit() {
		return jutil.StringPtr(getPropElem.GetSecondChild().GetString())
	} else {
		return nil
	}
}
