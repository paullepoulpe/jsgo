// Port of com/google/javascript/jscomp/OptimizeCalls.java @ closure-compiler v20261015

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// OptimizeCalls is a root pass that is a container for other passes that
// should run on with a single call graph.
//
// Known passes include:
//
//   - OptimizeParameters (remove unused and inline constant parameters)
//   - OptimizeReturns (remove unused)
//   - DevirtualizeMethods
type OptimizeCalls struct {
	compiler        AbstractCompiler
	passes          []OptimizeCalls_CallGraphCompilerPass
	considerExterns bool
}

// Java: OptimizeCalls(AbstractCompiler, ImmutableList<CallGraphCompilerPass>, boolean)
func NewOptimizeCalls(
	compiler AbstractCompiler,
	passes []OptimizeCalls_CallGraphCompilerPass,
	considerExterns bool) *OptimizeCalls {
	return &OptimizeCalls{compiler: compiler, passes: passes, considerExterns: considerExterns}
}

// OptimizeCalls_Builder_ is the static factory builder() (the trailing
// underscore keeps it apart from the nested type Builder).
//
// Java: OptimizeCalls.builder()
func OptimizeCalls_Builder_() *OptimizeCalls_Builder {
	return NewOptimizeCalls_Builder()
}

// OptimizeCalls_CallGraphCompilerPass is the nested interface
// CallGraphCompilerPass.
type OptimizeCalls_CallGraphCompilerPass interface {
	// Java: OptimizeCalls.CallGraphCompilerPass.process(Node, Node, ReferenceMap)
	Process(externs *rhino.Node, root *rhino.Node, references *OptimizeCalls_ReferenceMap)
}

// OptimizeCalls_CallGraphCompilerPassFunc lets a closure (a Java lambda) be
// used as a CallGraphCompilerPass.
//
// Go-only: language-adapter — Java lambdas implement the functional
// interface directly.
type OptimizeCalls_CallGraphCompilerPassFunc func(externs *rhino.Node, root *rhino.Node, references *OptimizeCalls_ReferenceMap)

// Process calls f.
func (f OptimizeCalls_CallGraphCompilerPassFunc) Process(externs *rhino.Node, root *rhino.Node, references *OptimizeCalls_ReferenceMap) {
	f(externs, root, references)
}

// OptimizeCalls_Builder is the nested class Builder.
type OptimizeCalls_Builder struct {
	compiler        AbstractCompiler
	passes          *jutil.ImmutableList_Builder[OptimizeCalls_CallGraphCompilerPass]
	considerExterns *bool // Nullable to force users to specify a value.
}

// Java: OptimizeCalls.Builder.setCompiler(AbstractCompiler)
func (b *OptimizeCalls_Builder) SetCompiler(compiler AbstractCompiler) *OptimizeCalls_Builder {
	b.compiler = compiler
	return b
}

// SetConsiderExterns sets whether or not to include references to extern
// names and properties in the ReferenceMap being generated.
//
// If considered, references to externs in both extern code and executable
// code will be collected. Otherwise, neither will be.
//
// This setting allows extern references to be effectively invisible to
// passes that should not mutate them.
//
// Java: OptimizeCalls.Builder.setConsiderExterns(boolean)
func (b *OptimizeCalls_Builder) SetConsiderExterns(v bool) *OptimizeCalls_Builder {
	b.considerExterns = &v
	return b
}

// Java: OptimizeCalls.Builder.addPass(CallGraphCompilerPass)
func (b *OptimizeCalls_Builder) AddPass(pass OptimizeCalls_CallGraphCompilerPass) *OptimizeCalls_Builder {
	b.passes.Add(pass)
	return b
}

// Java: OptimizeCalls.Builder.build()
func (b *OptimizeCalls_Builder) Build() *OptimizeCalls {
	jutil.CheckNotNull(b.compiler)
	jutil.CheckNotNull(b.considerExterns)

	return NewOptimizeCalls(b.compiler, b.passes.Build(), *b.considerExterns)
}

// Java: OptimizeCalls.Builder()
func NewOptimizeCalls_Builder() *OptimizeCalls_Builder {
	return &OptimizeCalls_Builder{passes: jutil.ImmutableList_Builder_[OptimizeCalls_CallGraphCompilerPass]()}
}

// Java: OptimizeCalls.process(Node, Node)
func (o *OptimizeCalls) Process(externs *rhino.Node, root *rhino.Node) {
	// Only global names are collected, which is insufficient if names have not been normalized.
	jutil.CheckState(o.compiler.GetLifeCycleStage().IsNormalized())

	if len(o.passes) == 0 {
		return
	}

	references := NewOptimizeCalls_ReferenceMap()
	NodeTraversal_TraverseRoots(
		o.compiler, NewOptimizeCalls_ReferenceMapBuildingCallback(o, references), externs, root)
	o.EliminateAccessorsFrom(references)

	for _, pass := range o.passes {
		pass.Process(externs, root, references)
	}
}

// EliminateAccessorsFrom deletes getter and setter names from references.
//
// Accessor names are disqualified from being in the ReferenceMap. We don't
// intentionally collect them, but other properties may share the same
// names. One reason why we do this is exemplified below:
//
//	class A {
//	  pure() { }
//	}
//
//	class B {
//	  get pure() { return impure; }
//	}
//
//	var x = (Math.random() > 0.5) ? new A() : new B();
//	x.pure(); // We can't safely optimize this call.
//
// Java: OptimizeCalls.eliminateAccessorsFrom(ReferenceMap)
func (o *OptimizeCalls) EliminateAccessorsFrom(references *OptimizeCalls_ReferenceMap) {
	// references.props.keySet().removeAll(accessors.keySet())
	for name := range o.compiler.GetAccessorSummary().GetAccessors().Keys() {
		references.props.Remove(name)
	}
}

type optimizeCalls_refListSlot struct {
	list jutil.ArrayList[*rhino.Node]
	buf  [4]*rhino.Node
}

// OptimizeCalls_ReferenceMap is a reference map for global symbols and
// properties.
type OptimizeCalls_ReferenceMap struct {
	globalScope *Scope
	names       *jutil.LinkedHashMap[string, *jutil.ArrayList[*rhino.Node]]
	props       *jutil.LinkedHashMap[string, *jutil.ArrayList[*rhino.Node]]
	listSlab    []optimizeCalls_refListSlot
}

// Java: OptimizeCalls.ReferenceMap()
func NewOptimizeCalls_ReferenceMap() *OptimizeCalls_ReferenceMap {
	return &OptimizeCalls_ReferenceMap{
		names: jutil.NewLinkedHashMapWithCapacity[string, *jutil.ArrayList[*rhino.Node]](1024),
		props: jutil.NewLinkedHashMapWithCapacity[string, *jutil.ArrayList[*rhino.Node]](1024),
	}
}

func (m *OptimizeCalls_ReferenceMap) allocRefList() *jutil.ArrayList[*rhino.Node] {
	if len(m.listSlab) == 0 {
		m.listSlab = make([]optimizeCalls_refListSlot, 512)
	}
	slot := &m.listSlab[0]
	m.listSlab = m.listSlab[1:]
	slot.list.InitWithBuffer(slot.buf[:])
	return &slot.list
}

// Java: OptimizeCalls.ReferenceMap.addReference(LinkedHashMap<String, ArrayList<Node>>, String, Node)
func (m *OptimizeCalls_ReferenceMap) AddReference(data *jutil.LinkedHashMap[string, *jutil.ArrayList[*rhino.Node]], name string, n *rhino.Node) {
	refs := data.Get(name)
	if refs == nil {
		refs = m.allocRefList()
		data.PutNew(name, refs)
	}
	refs.Add(n)
}

// Java: OptimizeCalls.ReferenceMap.addNameReference(String, Node)
func (m *OptimizeCalls_ReferenceMap) AddNameReference(name string, n *rhino.Node) {
	m.AddReference(m.names, name, n)
}

// Java: OptimizeCalls.ReferenceMap.addPropReference(String, Node)
func (m *OptimizeCalls_ReferenceMap) AddPropReference(name string, n *rhino.Node) {
	m.AddReference(m.props, name, n)
}

// Java: OptimizeCalls.ReferenceMap.getGlobalScope()
func (m *OptimizeCalls_ReferenceMap) GetGlobalScope() *Scope {
	return m.globalScope
}

// Java: OptimizeCalls.ReferenceMap.getNameReferences()
func (m *OptimizeCalls_ReferenceMap) GetNameReferences() iter.Seq2[string, *jutil.ArrayList[*rhino.Node]] {
	return m.names.All()
}

// Java: OptimizeCalls.ReferenceMap.getPropReferences()
func (m *OptimizeCalls_ReferenceMap) GetPropReferences() iter.Seq2[string, *jutil.ArrayList[*rhino.Node]] {
	return m.props.All()
}

// OptimizeCalls_ReferenceMap_GetFunctionNodes, given a set of references,
// returns the set of known definitions; specifically, those of the form:
// `function x() { }` or `x = ...;`
//
// As much as possible, functions are collected from conditional
// definitions. This is useful for optimizations that can be performed when
// the callers are known but all definitions may not be (unused call results,
// parameters that are never provided). Examples expressions:
//
//   - `(a(), function() {})`
//   - `a && function(){}`
//   - `b || function(){}`
//   - `a ? function() {} : function() {}`
//
// definitionSites are the definition site nodes to search for associated
// functions. These should taken from a ReferenceMap since only the node
// types collected by ReferenceMap are supported. Returns a mapping from the
// input definitionSites to each of their associated functions.
//
// Java: OptimizeCalls.ReferenceMap.getFunctionNodes(List<Node>)
func OptimizeCalls_ReferenceMap_GetFunctionNodes(definitionSites []*rhino.Node) *jutil.ListMultimap[*rhino.Node, *rhino.Node] {
	result := jutil.NewImmutableListMultimapBuilder[*rhino.Node, *rhino.Node]()
	for _, def := range definitionSites {
		// result.putAll(def, definitionFunctionNodesFor(def))
		for _, fn := range OptimizeCalls_ReferenceMap_DefinitionFunctionNodesFor(def) {
			result.Put(def, fn)
		}
	}
	return result.Build()
}

// OptimizeCalls_ReferenceMap_DefinitionFunctionNodesFor collects potential
// definition FUNCTIONs associated with a method definition site.
//
// See getFunctionNodes.
//
// Java: OptimizeCalls.ReferenceMap.definitionFunctionNodesFor(Node)
func OptimizeCalls_ReferenceMap_DefinitionFunctionNodesFor(definitionSite *rhino.Node) []*rhino.Node {
	if definitionSite.IsGetterDef() || definitionSite.IsSetterDef() {
		// TODO(nickreid): Support getters and setters. Ignore them for now since they aren't
		// "called".
		return nil
	}

	// Ignore detached nodes.
	parent := definitionSite.GetParent()
	if parent == nil {
		return nil
	}

	var fns []*rhino.Node
	switch parent.GetToken() {
	case rhino.Token_CLASS:
		if definitionSite.IsFirstChildOf(parent) {
			constructorFnDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(parent)
			if constructorFnDef != nil {
				fns = append(fns, constructorFnDef.GetOnlyChild())
			}
		}
	case rhino.Token_FUNCTION:
		fns = append(fns, parent)
	case rhino.Token_CLASS_MEMBERS:
		if definitionSite.IsMemberFunctionDef() {
			fns = append(fns, definitionSite.GetLastChild())
		} else {
			jutil.CheckArgumentMsg1(definitionSite.IsMemberFieldDef(), "%s", definitionSite)
			value := definitionSite.GetFirstChild()
			if value != nil {
				optimizeCalls_addValueFunctionNodes(&fns, value)
			}
		}
	case rhino.Token_OBJECTLIT:
		jutil.CheckArgumentMsg1(
			definitionSite.IsStringKey() || definitionSite.IsMemberFunctionDef(), //
			"%s", definitionSite)
		optimizeCalls_addValueFunctionNodes(&fns, definitionSite.GetLastChild())
	case rhino.Token_ASSIGN:
		// Only a candidate if the assign isn't consumed.
		target := parent.GetFirstChild()
		value := parent.GetLastChild()
		if definitionSite == target {
			optimizeCalls_addValueFunctionNodes(&fns, value)
		}
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		if definitionSite.IsName() && definitionSite.HasChildren() {
			optimizeCalls_addValueFunctionNodes(&fns, definitionSite.GetFirstChild())
		}
	default:
	}
	return fns
}

func optimizeCalls_addValueFunctionNodes(fns *[]*rhino.Node, n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		constructorFnDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
		if constructorFnDef != nil {
			*fns = append(*fns, constructorFnDef.GetOnlyChild())
		}
	case rhino.Token_FUNCTION:
		*fns = append(*fns, n)
	case rhino.Token_HOOK:
		optimizeCalls_addValueFunctionNodes(fns, n.GetSecondChild())
		optimizeCalls_addValueFunctionNodes(fns, n.GetLastChild())
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		optimizeCalls_addValueFunctionNodes(fns, n.GetFirstChild())
		optimizeCalls_addValueFunctionNodes(fns, n.GetLastChild())
	case rhino.Token_CAST, rhino.Token_COMMA:
		optimizeCalls_addValueFunctionNodes(fns, n.GetLastChild())
	default:
	}
}

// Java: OptimizeCalls.ReferenceMap.addValueFunctionNodes(ImmutableList.Builder<Node>, Node)
func OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns *jutil.ImmutableList_Builder[*rhino.Node], n *rhino.Node) {
	// TODO(johnlenz): add member definitions
	switch n.GetToken() {
	case rhino.Token_CLASS:
		{
			constructorFnDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
			if constructorFnDef != nil {
				fns.Add(constructorFnDef.GetOnlyChild())
			}
		}
	case rhino.Token_FUNCTION:
		fns.Add(n)
	case rhino.Token_HOOK:
		OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetSecondChild())
		OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetLastChild())
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetFirstChild())
		OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetLastChild())
	case rhino.Token_CAST, rhino.Token_COMMA:
		OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fns, n.GetLastChild())
	default:
		// do nothing.
	}
}

// OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget reports
// whether the provided node acts as the target function in a new or call or
// optional chain call expression including .call expressions. For example,
// returns true for 'x' in 'x?.call()'.
//
// Java: OptimizeCalls.ReferenceMap.isNormalOrOptChainCallOrNewTarget(Node)
func OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n *rhino.Node) bool {
	return OptimizeCalls_ReferenceMap_IsCallTarget(n) || OptimizeCalls_ReferenceMap_IsNewTarget(n) ||
		OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n)
}

// OptimizeCalls_ReferenceMap_IsCallTarget reports whether the provided node
// acts as the target function in a call expression including .call
// expressions. For example, returns true for 'x' in 'x.call()'.
//
// Java: OptimizeCalls.ReferenceMap.isCallTarget(Node)
func OptimizeCalls_ReferenceMap_IsCallTarget(n *rhino.Node) bool {
	parent := n.GetParent()
	if parent.IsCall() && n.IsFirstChildOf(parent) {
		return true
	}

	grandParent := parent.GetParent()
	return parent.IsGetProp() &&
		grandParent.IsCall() &&
		parent.IsFirstChildOf(grandParent) &&
		parent.GetString() == "call"
}

// OptimizeCalls_ReferenceMap_IsOptChainCallTarget reports whether the
// provided node acts as the target function in an optional chain call
// expression including .call expressions. For example, returns true for 'x'
// in 'x?.call()'.
//
// Java: OptimizeCalls.ReferenceMap.isOptChainCallTarget(Node)
func OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n *rhino.Node) bool {
	parent := n.GetParent()
	if parent.IsOptChainCall() && n.IsFirstChildOf(parent) {
		return true
	}

	grandParent := parent.GetParent()
	return parent.IsOptChainGetProp() &&
		grandParent.IsOptChainCall() &&
		parent.IsFirstChildOf(grandParent) &&
		parent.GetString() == "call"
}

// OptimizeCalls_ReferenceMap_IsNewTarget reports whether the provided node
// acts as the target function in a new expression.
//
// Java: OptimizeCalls.ReferenceMap.isNewTarget(Node)
func OptimizeCalls_ReferenceMap_IsNewTarget(n *rhino.Node) bool {
	parent := n.GetParent()
	return parent.IsNew() && parent.GetFirstChild() == n
}

// OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget finds the associated
// call node for a node for which isNormalOrOptChainCallOrNewTarget returns
// true.
//
// Java: OptimizeCalls.ReferenceMap.getCallOrNewNodeForTarget(Node)
func OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n *rhino.Node) *rhino.Node {
	maybeCall := n.GetParent()
	jutil.CheckStateMsg2(n.IsFirstChildOf(maybeCall), "%s\n\n%s", maybeCall, n)

	if NodeUtil_IsCallOrNew(maybeCall) {
		// e.g. `n` input param is the `a` in `a()` or `a?.()`
		return maybeCall
	} else {
		// e.g. `n` input param is the `a` in `a.b()` or `a?.b()`.
		child := maybeCall
		maybeCall = child.GetParent()

		jutil.CheckStateMsg1(NodeUtil_IsNormalOrOptChainGetProp(child), "%s", child)
		jutil.CheckStateMsg1(NodeUtil_IsNormalOrOptChainCall(maybeCall), "%s", maybeCall)
		jutil.CheckStateMsg2(child.IsFirstChildOf(maybeCall), "%s\n\n%s", maybeCall, child)

		return maybeCall
	}
}

// OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall finds
// the call argument node matching the first parameter of the called
// function for a node for which isNormalOrOptChainCallOrNewTarget returns
// true. Specifically, corrects for the additional argument provided to
// .call expressions.
//
// Java: OptimizeCalls.ReferenceMap.getFirstArgumentForCallOrNewOrDotCall(Node)
func OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n *rhino.Node) *rhino.Node {
	return OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(n, 0)
}

// OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall finds the call
// argument node matching the parameter at the specified index of the called
// function for a node for which isNormalOrOptChainCallOrNewTarget returns
// true. Specifically, corrects for the additional argument provided to
// .call expressions.
//
// Java: OptimizeCalls.ReferenceMap.getArgumentForCallOrNewOrDotCall(Node, int)
func OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(n *rhino.Node, index int) *rhino.Node {
	adjustedIndex := index
	parent := n.GetParent()
	if !(parent.IsCall() || parent.IsOptChainCall() || parent.IsNew()) {
		parent = parent.GetParent()
		if NodeUtil_IsFunctionObjectCall(parent) {
			adjustedIndex++
		}
	}
	return NodeUtil_GetArgumentForCallOrNew(parent, adjustedIndex)
}

// Java: OptimizeCalls.ReferenceMap.isSimpleAssignmentTarget(Node)
func OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n *rhino.Node) bool {
	parent := n.GetParent()
	// `ref = value;`
	return parent.IsAssign() && n.IsFirstChildOf(parent) && parent.GetParent().IsExprResult()
}

var optimizeCalls_emptyStringSet = jutil.NewLinkedHashSet[string]()

// Java: OptimizeCalls.safeSet(Set<String>)
func OptimizeCalls_SafeSet(set *jutil.LinkedHashSet[string]) *jutil.LinkedHashSet[string] {
	if set != nil {
		return set
	}
	return optimizeCalls_emptyStringSet
}

// OptimizeCalls_ReferenceMapBuildingCallback is the inner class
// ReferenceMapBuildingCallback.
type OptimizeCalls_ReferenceMapBuildingCallback struct {
	outer       *OptimizeCalls
	externProps *jutil.LinkedHashSet[string]
	references  *OptimizeCalls_ReferenceMap
	globalScope *Scope
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback(ReferenceMap)
func NewOptimizeCalls_ReferenceMapBuildingCallback(outer *OptimizeCalls, references *OptimizeCalls_ReferenceMap) *OptimizeCalls_ReferenceMapBuildingCallback {
	c := &OptimizeCalls_ReferenceMapBuildingCallback{outer: outer}
	c.externProps = OptimizeCalls_SafeSet(outer.compiler.GetExternProperties())
	c.references = references
	return c
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.visit(NodeTraversal, Node, Node)
func (c *OptimizeCalls_ReferenceMapBuildingCallback) Visit(t *NodeTraversal, n *rhino.Node, unused *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		c.MaybeAddNameReference(n.GetString(), n)
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_GETPROP:
		c.MaybeAddPropReference(n.GetString(), n)
	case rhino.Token_CALL:
		// If we are using goog.reflect.objectProperty on this symbol, we will assume that it
		// gets referenced.
		fnName := n.GetFirstChild()
		if c.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(fnName) {
			propName := NodeUtil_GetArgumentForCallOrNew(n, 0)
			if propName != nil {
				c.MaybeAddPropReference(propName.GetString(), n)
			}
		}

	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF:
		// ignore quoted keys.
		if !n.IsQuotedStringKey() {
			c.MaybeAddPropReference(n.GetString(), n)
		}

	case rhino.Token_SUPER:
		c.VisitSuper(n)
	case
		// Ignore quoted and computed keys.
		// TODO(johnlenz): support symbols.
		rhino.Token_COMPUTED_PROP,
		rhino.Token_COMPUTED_FIELD_DEF,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_GETELEM,
		// Don't worry about invisible accesses using object rest/spread. To be invoked there
		// would need to be downstream references that use the actual name. We'd see those.
		rhino.Token_OBJECT_REST,
		rhino.Token_OBJECT_SPREAD:
	default:
	}
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.visitSuper(Node)
func (c *OptimizeCalls_ReferenceMapBuildingCallback) VisitSuper(superNode *rhino.Node) {
	// Determine whether this is a super() constructor call.
	// If it is, identify the super class and record this as a reference to that.
	parent := superNode.GetParent()
	if parent.IsCall() && superNode.IsFirstChildOf(parent) {
		enclosingClass := jutil.CheckNotNull(NodeUtil_GetEnclosingClass(parent))
		extendsNode := enclosingClass.GetSecondChild()
		jutil.CheckStateMsg(!extendsNode.IsEmpty(), "super call appears in class without extends clause")
		if extendsNode.IsName() {
			c.MaybeAddNameReference(extendsNode.GetString(), superNode)
		} else if extendsNode.IsGetProp() {
			// NOTE: Theoretically we could also include an optional chain getprop here, but
			// A) it's a runtime error if the value ends up being undefined, so that's bad code
			// B) the author is indicating uncertainty, so we should be cautious.
			c.MaybeAddPropReference(extendsNode.GetString(), superNode)
		} // else we cannot tell what super() is referencing (e.g. `class extends getMixin() {`)
	}
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.maybeAddNameReference(String, Node)
func (c *OptimizeCalls_ReferenceMapBuildingCallback) MaybeAddNameReference(name string, n *rhino.Node) {
	// TODO(b/129503101): Why are we limiting ourselves to global names?
	var_, _ := c.globalScope.GetSlot(name).(*Var)
	if var_ != nil && (c.outer.considerExterns || !var_.IsExtern()) {
		// As every name declaration is unique due to normalizations, it is only necessary to build
		// the global scope and ask it if it knows about a name as it can never be shadowed.
		c.references.AddNameReference(name, n)
	}
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.maybeAddPropReference(String, Node)
func (c *OptimizeCalls_ReferenceMapBuildingCallback) MaybeAddPropReference(name string, n *rhino.Node) {
	if c.outer.considerExterns || !c.externProps.Contains(name) {
		c.references.AddPropReference(name, n)
	}
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *OptimizeCalls_ReferenceMapBuildingCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsScript() {
		// Even when considering externs, we only care about top-level identifiers. Dummy function
		// parameters, for example, shouldn't be considered references.
		return (c.outer.considerExterns && t.InGlobalScope()) || !n.IsFromExterns()
	} else {
		return true
	}
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.enterScope(NodeTraversal)
func (c *OptimizeCalls_ReferenceMapBuildingCallback) EnterScope(t *NodeTraversal) {
	if t.InGlobalScope() {
		c.globalScope = t.GetScope()
		c.references.globalScope = c.globalScope
	}
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.exitScope(NodeTraversal)
func (c *OptimizeCalls_ReferenceMapBuildingCallback) ExitScope(t *NodeTraversal) {}

// OptimizeCalls_MayBeOptimizableName returns whether the provide name may be
// a candidate for call optimizations.
//
// Java: OptimizeCalls.mayBeOptimizableName(AbstractCompiler, String)
func OptimizeCalls_MayBeOptimizableName(compiler AbstractCompiler, name string) bool {
	if compiler.GetCodingConvention().IsExported(name) {
		return false
	}

	// Avoid modifying a few special case functions. Specifically, $jscomp.inherits to
	// recognize 'inherits' calls. (b/27244988)
	if name == NodeUtil_JSC_PROPERTY_NAME_FN ||
		name == "inherits" ||
		name == "$jscomp$inherits" ||
		name == "goog$inherits" {
		return false
	}
	return true
}

// OptimizeCalls_IsAllowedReference returns whether the reference is a known
// non-aliasing reference.
//
// Java: OptimizeCalls.isAllowedReference(Node)
func OptimizeCalls_IsAllowedReference(n *rhino.Node) bool {
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// inspecting the properties is allowed.
		return parent.GetSecondChild() == n
	case rhino.Token_INSTANCEOF, rhino.Token_TYPEOF, rhino.Token_IN:
		return true
	case rhino.Token_GETELEM, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM:
		// Calls escape the "this" value. a.foo() aliases "a" as "this" but general
		// property references do not.
		grandparent := parent.GetParent()
		if n == parent.GetFirstChild() &&
			grandparent != nil &&
			(grandparent.IsCall() ||
				grandparent.IsOptChainCall() ||
				grandparent.IsTaggedTemplateLit()) {
			return false // `a.foo()` or `a?.foo()` or `a?.[foo]()` or `a.foo`bar``
		}
		return true
	case rhino.Token_CLASS:
		if n.IsFirstChildOf(parent) {
			// class Name {
			// this is a definition, not a read reference
			return false
		} else {
			// class SubClass extends Name {
			jutil.CheckStateMsg1(n.IsSecondChildOf(parent), "%s", parent)
			// find the constructor
			if NodeUtil_GetEs6ClassConstructorMemberFunctionDef(parent) == nil {
				// The subclass has no explicit constructor, so `new SubClass()` implicitly calls
				// `new Name(...arguments)`. This hidden call makes it harder to safely optimize the
				// `Name` constructor, so we won't do it.
				return false
			} else {
				// We can still optimize the constructor of the class being extended as
				// long as all child classes have explicit constructors, so we can see the
				// `super()` calls in them and update them.
				return true
			}
		}
		//
	default:
		if NodeUtil_IsNameDeclaration(parent) && !n.HasChildren() {
			// allow "let x;"
			return true
		}
	}
	return false
}
