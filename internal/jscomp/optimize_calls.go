// Port of com/google/javascript/jscomp/OptimizeCalls.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	Process(externs rhino.Node, root rhino.Node, references *OptimizeCalls_ReferenceMap)
}

// OptimizeCalls_CallGraphCompilerPassFunc lets a closure (a Java lambda) be
// used as a CallGraphCompilerPass.
//
// Go-only: language-adapter — Java lambdas implement the functional
// interface directly.
type OptimizeCalls_CallGraphCompilerPassFunc func(externs rhino.Node, root rhino.Node, references *OptimizeCalls_ReferenceMap)

// Process calls f.
func (f OptimizeCalls_CallGraphCompilerPassFunc) Process(externs rhino.Node, root rhino.Node, references *OptimizeCalls_ReferenceMap) {
	f(externs, root, references)
}

// OptimizeCalls_Builder is the nested class Builder.
type OptimizeCalls_Builder struct {
	compiler           AbstractCompiler
	passes             []OptimizeCalls_CallGraphCompilerPass
	considerExterns    bool
	considerExternsSet bool
}

// Java: OptimizeCalls.Builder.setCompiler(AbstractCompiler)
func (b *OptimizeCalls_Builder) SetCompiler(compiler AbstractCompiler) *OptimizeCalls_Builder {
	b.compiler = compiler
	return b
}

// SetConsiderExterns sets whether or not to include references to extern
// names and properties in the ReferenceMap being generated.
//
// Java: OptimizeCalls.Builder.setConsiderExterns(boolean)
func (b *OptimizeCalls_Builder) SetConsiderExterns(v bool) *OptimizeCalls_Builder {
	b.considerExterns = v
	b.considerExternsSet = true
	return b
}

// Java: OptimizeCalls.Builder.addPass(CallGraphCompilerPass)
func (b *OptimizeCalls_Builder) AddPass(pass OptimizeCalls_CallGraphCompilerPass) *OptimizeCalls_Builder {
	b.passes = append(b.passes, jutil.CheckNotNull(pass))
	return b
}

// Java: OptimizeCalls.Builder.build()
func (b *OptimizeCalls_Builder) Build() *OptimizeCalls {
	jutil.CheckNotNull(b.compiler)
	jutil.CheckState(b.considerExternsSet)

	return NewOptimizeCalls(b.compiler, slices.Clone(b.passes), b.considerExterns)
}

// Java: OptimizeCalls.Builder()
func NewOptimizeCalls_Builder() *OptimizeCalls_Builder {
	return &OptimizeCalls_Builder{}
}

// Java: OptimizeCalls.process(Node, Node)
func (o *OptimizeCalls) Process(externs rhino.Node, root rhino.Node) {
	// Only global names are collected, which is insufficient if names have not been normalized.
	jutil.CheckState(o.compiler.GetLifeCycleStage().IsNormalized())

	if len(o.passes) == 0 {
		return
	}

	references := NewOptimizeCalls_ReferenceMap()
	sc := NewSyntacticScopeCreator(o.compiler)
	releaseEarly := len(o.passes) == 1
	if releaseEarly {
		if _, isPFI := o.passes[0].(*PureFunctionIdentifier); !isPFI {
			releaseEarly = false
		}
	}
	cb := newOptimizeCallsRefCollector(o, references)
	scopeRoot := externs.GetParent()
	if scopeRoot.IsNil() {
		failNodeTraversalNull()
	}
	jutil.CheckState(root.GetParent() == scopeRoot)
	cb.initGlobalScope(sc.CreateScopeOfScope(scopeRoot, nil))
	cb.traverseNode(externs)
	cb.traverseNode(root)
	if releaseEarly {
		references.globalScope = nil
	}
	o.eliminateAccessorsFrom(references)

	for _, pass := range o.passes {
		pass.Process(externs, root, references)
	}
}

// eliminateAccessorsFrom deletes getter and setter names from references.
//
// Java: OptimizeCalls.eliminateAccessorsFrom(ReferenceMap)
func (o *OptimizeCalls) eliminateAccessorsFrom(references *OptimizeCalls_ReferenceMap) {
	for name := range o.compiler.GetAccessorSummary().GetAccessors().Keys() {
		references.props.Remove(name)
	}
}

// OptimizeCalls_ReferenceMap is a reference map for global symbols and
// properties.
type OptimizeCalls_ReferenceMap struct {
	globalScope *Scope
	names       optimizeCalls_refTable
	props       optimizeCalls_refTable
	refLists    jutil.Slab[optimizeCalls_refList] // one pass's lists die together
}

// Java: OptimizeCalls.ReferenceMap()
func NewOptimizeCalls_ReferenceMap() *OptimizeCalls_ReferenceMap {
	m := &OptimizeCalls_ReferenceMap{}
	m.names.initWithCapacity(256)
	m.props.initWithCapacity(512)
	return m
}

// optimizeCalls_refList is a reference list with room for its first nodes in
// the same allocation.
type optimizeCalls_refList struct {
	name    string
	nodes   []rhino.Node
	buf     [4]rhino.Node
	removed bool
}

type optimizeCalls_refTable struct {
	byName map[string]*optimizeCalls_refList
	order  []*optimizeCalls_refList
	size   int
}

func (t *optimizeCalls_refTable) initWithCapacity(capHint int) {
	t.byName = make(map[string]*optimizeCalls_refList, capHint)
	t.order = make([]*optimizeCalls_refList, 0, min(capHint, 4096))
	t.size = 0
}

func (t *optimizeCalls_refTable) Size() int     { return t.size }
func (t *optimizeCalls_refTable) IsEmpty() bool { return t.size == 0 }

func (t *optimizeCalls_refTable) Remove(name string) {
	if l := t.byName[name]; l != nil {
		delete(t.byName, name)
		l.removed = true
		l.nodes = nil
		t.size--
	}
}

func (t *optimizeCalls_refTable) all() iter.Seq2[string, []rhino.Node] {
	return func(yield func(string, []rhino.Node) bool) {
		for _, l := range t.order {
			if !l.removed && !yield(l.name, l.nodes) {
				return
			}
		}
	}
}

// appendDoubling is append that doubles the capacity of a full large slice.
//
// Go-only: performance — append already doubles up to 256 elements, then grows
// by only 1.25x, which allocates about five times the final size of a long list.
func appendDoubling[T any](s []T, e T) []T {
	if len(s) == cap(s) && len(s) >= 256 {
		grown := make([]T, len(s), 2*len(s))
		copy(grown, s)
		s = grown
	}
	return append(s, e)
}

// Java: OptimizeCalls.ReferenceMap.addReference(LinkedHashMap<String, ArrayList<Node>>, String, Node)
func (m *OptimizeCalls_ReferenceMap) AddReference(data *optimizeCalls_refTable, name string, n rhino.Node) {
	refs := data.byName[name]
	if refs == nil {
		refs = m.refLists.New()
		refs.name = name
		refs.nodes = refs.buf[:0:len(refs.buf)]
		data.byName[name] = refs
		data.order = append(data.order, refs)
		data.size++
	}
	refs.nodes = appendDoubling(refs.nodes, n)
}

// Java: OptimizeCalls.ReferenceMap.addNameReference(String, Node)
func (m *OptimizeCalls_ReferenceMap) AddNameReference(name string, n rhino.Node) {
	m.AddReference(&m.names, name, n)
}

// Java: OptimizeCalls.ReferenceMap.addPropReference(String, Node)
func (m *OptimizeCalls_ReferenceMap) AddPropReference(name string, n rhino.Node) {
	m.AddReference(&m.props, name, n)
}

// Java: OptimizeCalls.ReferenceMap.getGlobalScope()
func (m *OptimizeCalls_ReferenceMap) GetGlobalScope() *Scope {
	return m.globalScope
}

// Java: OptimizeCalls.ReferenceMap.getNameReferences()
func (m *OptimizeCalls_ReferenceMap) GetNameReferences() iter.Seq2[string, []rhino.Node] {
	return m.names.all()
}

// Java: OptimizeCalls.ReferenceMap.getPropReferences()
func (m *OptimizeCalls_ReferenceMap) GetPropReferences() iter.Seq2[string, []rhino.Node] {
	return m.props.all()
}

// OptimizeCalls_ReferenceMap_AppendFunctionNodes appends all function nodes
// associated with definitionSites to dst and returns the updated slice.
func OptimizeCalls_ReferenceMap_AppendFunctionNodes(dst []rhino.Node, definitionSites []rhino.Node) []rhino.Node {
	for _, def := range definitionSites {
		dst = optimizeCalls_appendDefinitionFunctionNodesFor(dst, def)
	}
	return dst
}

// OptimizeCalls_ReferenceMap_GetSingleFunctionNode returns the single function
// node associated with definitionSites if there is exactly 1, or nil if 0 or > 1.
func OptimizeCalls_ReferenceMap_GetSingleFunctionNode(definitionSites []rhino.Node) rhino.Node {
	fn, count := optimizeCalls_getSingleFunctionNode(definitionSites)
	if count == 1 {
		return fn
	}
	return rhino.NilNode
}

func optimizeCalls_getSingleFunctionNode(definitionSites []rhino.Node) (rhino.Node, int) {
	var buf [2]rhino.Node
	fns := buf[:0]
	for _, def := range definitionSites {
		fns = optimizeCalls_appendDefinitionFunctionNodesLimit(fns, def, 2)
		if len(fns) > 1 {
			return rhino.NilNode, len(fns)
		}
	}
	if len(fns) == 1 {
		return fns[0], 1
	}
	return rhino.NilNode, 0
}

// OptimizeCalls_ReferenceMap_DefinitionFunctionNodesFor collects potential
// definition FUNCTIONs associated with a method definition site.
//
// Java: OptimizeCalls.ReferenceMap.definitionFunctionNodesFor(Node)
func OptimizeCalls_ReferenceMap_DefinitionFunctionNodesFor(definitionSite rhino.Node) []rhino.Node {
	return optimizeCalls_appendDefinitionFunctionNodesFor(nil, definitionSite)
}

func optimizeCalls_appendDefinitionFunctionNodesFor(fns []rhino.Node, definitionSite rhino.Node) []rhino.Node {
	return optimizeCalls_appendDefinitionFunctionNodesLimit(fns, definitionSite, 0)
}

// optimizeCalls_appendDefinitionFunctionNodesLimit appends function nodes
// associated with definitionSite up to limit (0 means unlimited).
func optimizeCalls_appendDefinitionFunctionNodesLimit(fns []rhino.Node, definitionSite rhino.Node, limit int) []rhino.Node {
	if limit > 0 && len(fns) >= limit {
		return fns
	}
	if definitionSite.IsGetterDef() || definitionSite.IsSetterDef() {
		// TODO(nickreid): Support getters and setters. Ignore them for now since they aren't
		// "called".
		return fns
	}

	// Ignore detached nodes.
	parent := definitionSite.GetParent()
	if parent.IsNil() {
		return fns
	}

	switch parent.GetToken() {
	case rhino.Token_CLASS:
		cls := ast.MustClass(parent)
		if cls.NameNode() == definitionSite {
			if ctor, ok := cls.Constructor(); ok {
				fns = append(fns, ctor.Function().Node())
			}
		}
	case rhino.Token_FUNCTION:
		fns = append(fns, parent)
	case rhino.Token_CLASS_MEMBERS:
		if m, ok := ast.AsMemberFunctionDef(definitionSite); ok {
			fns = append(fns, m.Function().Node())
		} else {
			mf := ast.MustMemberFieldDef(definitionSite)
			if value, ok := mf.Value(); ok {
				fns = optimizeCalls_appendValueFunctionNodesLimit(fns, value, limit)
			}
		}
	case rhino.Token_OBJECTLIT:
		if pkm, ok := ast.AsPlainKeyMember(definitionSite); ok {
			if expr, isExpr := pkm.Value().AsExpr(); isExpr {
				fns = optimizeCalls_appendValueFunctionNodesLimit(fns, expr, limit)
			}
		} else {
			panic(jutil.NewIllegalArgumentExceptionMsg(definitionSite.String()))
		}
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(parent)
		if assign.Target().Is(definitionSite) {
			fns = optimizeCalls_appendValueFunctionNodesLimit(fns, assign.Value(), limit)
		}
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		if name, ok := ast.AsName(definitionSite); ok {
			if init, hasInit := name.Init(); hasInit {
				fns = optimizeCalls_appendValueFunctionNodesLimit(fns, init, limit)
			}
		}
	}
	return fns
}

func optimizeCalls_appendValueFunctionNodesLimit(fns []rhino.Node, expr ast.Expr, limit int) []rhino.Node {
	if limit > 0 && len(fns) >= limit {
		return fns
	}
	switch expr.Token() {
	case rhino.Token_CLASS:
		if ctor, ok := expr.MustClass().Constructor(); ok {
			fns = append(fns, ctor.Function().Node())
		}
	case rhino.Token_FUNCTION:
		fns = append(fns, expr.Node())
	case rhino.Token_HOOK:
		h := expr.MustHook()
		fns = optimizeCalls_appendValueFunctionNodesLimit(fns, h.Then(), limit)
		if limit == 0 || len(fns) < limit {
			fns = optimizeCalls_appendValueFunctionNodesLimit(fns, h.Else(), limit)
		}
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		b := expr.MustBinary()
		fns = optimizeCalls_appendValueFunctionNodesLimit(fns, b.Left(), limit)
		if limit == 0 || len(fns) < limit {
			fns = optimizeCalls_appendValueFunctionNodesLimit(fns, b.Right(), limit)
		}
	case rhino.Token_CAST:
		fns = optimizeCalls_appendValueFunctionNodesLimit(fns, expr.MustCast().Expr(), limit)
	case rhino.Token_COMMA:
		fns = optimizeCalls_appendValueFunctionNodesLimit(fns, expr.MustBinary().Right(), limit)
	}
	return fns
}

// OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget reports
// whether the provided node acts as the target function in a new or call or
// optional chain call expression including .call expressions.
//
// Java: OptimizeCalls.ReferenceMap.isNormalOrOptChainCallOrNewTarget(Node)
func OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n rhino.Node) bool {
	return OptimizeCalls_ReferenceMap_IsCallTarget(n) || OptimizeCalls_ReferenceMap_IsNewTarget(n) ||
		OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n)
}

// OptimizeCalls_ReferenceMap_IsCallTarget reports whether the provided node
// acts as the target function in a call expression including .call
// expressions (e.g. 'x' in 'x.call()').
//
// Java: OptimizeCalls.ReferenceMap.isCallTarget(Node)
func OptimizeCalls_ReferenceMap_IsCallTarget(n rhino.Node) bool {
	parent := n.GetParent()
	if call, ok := ast.AsCall(parent); ok && call.Callee().Is(n) {
		return true
	}
	gp, ok := ast.AsGetProp(parent)
	if !ok || gp.Property() != "call" {
		return false
	}
	call, ok := ast.AsCall(parent.GetParent())
	return ok && call.Callee().Is(parent)
}

// OptimizeCalls_ReferenceMap_IsOptChainCallTarget reports whether the
// provided node acts as the target function in an optional chain call
// expression including .call expressions (e.g. 'x' in 'x?.call()').
//
// Java: OptimizeCalls.ReferenceMap.isOptChainCallTarget(Node)
func OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n rhino.Node) bool {
	parent := n.GetParent()
	if call, ok := ast.AsOptChainCall(parent); ok && call.Callee().Is(n) {
		return true
	}
	gp, ok := ast.AsOptChainGetProp(parent)
	if !ok || gp.Property() != "call" {
		return false
	}
	call, ok := ast.AsOptChainCall(parent.GetParent())
	return ok && call.Callee().Is(parent)
}

// OptimizeCalls_ReferenceMap_IsNewTarget reports whether the provided node
// acts as the target function in a new expression.
//
// Java: OptimizeCalls.ReferenceMap.isNewTarget(Node)
func OptimizeCalls_ReferenceMap_IsNewTarget(n rhino.Node) bool {
	newExpr, ok := ast.AsNew(n.GetParent())
	return ok && newExpr.Callee().Is(n)
}

// OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget finds the associated
// call node for a node for which isNormalOrOptChainCallOrNewTarget returns
// true.
//
// Java: OptimizeCalls.ReferenceMap.getCallOrNewNodeForTarget(Node)
func OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n rhino.Node) rhino.Node {
	maybeCall := n.GetParent()
	jutil.CheckStateMsg2(n.IsFirstChildOf(maybeCall), "%s\n\n%s", maybeCall, n)

	if _, ok := ast.AsCallOrNew(maybeCall); ok {
		return maybeCall
	}
	child := maybeCall
	maybeCall = child.GetParent()
	_, isGetProp := ast.AsNormalOrOptChainGetProp(child)
	call, isCall := ast.AsNormalOrOptChainCall(maybeCall)
	jutil.CheckStateMsg1(isGetProp, "%s", child)
	jutil.CheckStateMsg1(isCall, "%s", maybeCall)
	jutil.CheckStateMsg2(call.Callee().Is(child), "%s\n\n%s", maybeCall, child)
	return maybeCall
}

// OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall finds
// the call argument node matching the first parameter of the called
// function for a node for which isNormalOrOptChainCallOrNewTarget returns
// true.
//
// Java: OptimizeCalls.ReferenceMap.getFirstArgumentForCallOrNewOrDotCall(Node)
func OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n rhino.Node) rhino.Node {
	return OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(n, 0)
}

// OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall finds the call
// argument node matching the parameter at the specified index of the called
// function for a node for which isNormalOrOptChainCallOrNewTarget returns
// true, correcting for the receiver argument passed to .call expressions.
//
// Java: OptimizeCalls.ReferenceMap.getArgumentForCallOrNewOrDotCall(Node, int)
func OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(n rhino.Node, index int) rhino.Node {
	adjustedIndex := index
	parent := n.GetParent()
	callOrNew, ok := ast.AsCallOrNew(parent)
	if !ok {
		parent = parent.GetParent()
		callOrNew = ast.MustCallOrNew(parent)
		if NodeUtil_IsFunctionObjectCall(parent) {
			adjustedIndex++
		}
	}
	if arg, found := callOrNew.ArgAt(adjustedIndex); found {
		return arg.Node()
	}
	return rhino.NilNode
}

// Java: OptimizeCalls.ReferenceMap.isSimpleAssignmentTarget(Node)
func OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n rhino.Node) bool {
	parent := n.GetParent()
	assign, ok := ast.AsAssign(parent)
	return ok && assign.Target().Is(n) && parent.GetParent().IsExprResult()
}

var optimizeCalls_emptyStringSet = jutil.NewLinkedHashSet[string]()

func optimizeCalls_safeSet(set *jutil.LinkedHashSet[string]) *jutil.LinkedHashSet[string] {
	if set != nil {
		return set
	}
	return optimizeCalls_emptyStringSet
}

type optimizeCallsRefCollector struct {
	outer            *OptimizeCalls
	codingConvention CodingConvention
	externProps      *jutil.LinkedHashSet[string]
	references       *OptimizeCalls_ReferenceMap
	globalScope      *Scope
}

func newOptimizeCallsRefCollector(outer *OptimizeCalls, references *OptimizeCalls_ReferenceMap) *optimizeCallsRefCollector {
	return &optimizeCallsRefCollector{
		outer:            outer,
		codingConvention: outer.compiler.GetCodingConvention(),
		externProps:      optimizeCalls_safeSet(outer.compiler.GetExternProperties()),
		references:       references,
	}
}

func (c *optimizeCallsRefCollector) initGlobalScope(s *Scope) {
	c.globalScope = s
	c.references.globalScope = s
	// Go-only: performance — presize the maps of a large program, whose
	// growth otherwise allocates about twice their final size.
	varCount := s.GetVarCount()
	if varCount > 4096 && c.references.names.IsEmpty() {
		c.references.names.initWithCapacity(varCount)
	}
	if varCount > 1024 && c.references.props.IsEmpty() {
		c.references.props.initWithCapacity(max(varCount*4, 16384))
	}
}

func (c *optimizeCallsRefCollector) traverseNode(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if !c.outer.considerExterns && n.IsFromExterns() {
			return
		}
	case rhino.Token_CLASS:
		cls := ast.MustClass(n)
		c.traverseNode(cls.ExtendsNode())
		c.traverseNode(cls.NameNode())
		c.traverseNode(cls.Body().Node())
		return
	case rhino.Token_CLASS_MEMBERS:
		body := ast.MustClassBody(n)
		for member := range body.Members() {
			if cm, ok := member.AsComputedMember(); ok {
				c.traverseNode(cm.Key().Node())
			}
		}
		for member := range body.Members() {
			if cm, ok := member.AsComputedMember(); ok {
				if val, hasVal := cm.Value(); hasVal {
					c.traverseNode(val.Node())
				}
			} else {
				c.traverseNode(member.Node())
			}
		}
		return
	}

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if child.GetFirstChild().IsNil() {
			switch child.GetToken() {
			case rhino.Token_NAME:
				c.maybeAddNameReferenceNode(child, child)
			case rhino.Token_SUPER:
				c.visitSuper(child)
			case rhino.Token_MEMBER_FIELD_DEF:
				if !child.IsQuotedStringKey() {
					c.maybeAddPropReference(child.GetString(), child)
				}
			}
			continue
		}
		c.traverseNode(child)
	}

	switch n.GetToken() {
	case rhino.Token_NAME:
		c.maybeAddNameReferenceNode(n, n)
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_GETPROP:
		c.maybeAddPropReference(n.GetString(), n)
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		if firstArg, ok := call.FirstArg(); ok && c.codingConvention.IsPropertyRenameFunction(call.Callee().Node()) {
			c.maybeAddPropReference(firstArg.Node().GetString(), n)
		}
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF:
		nm := ast.MustNamedMember(n)
		if !nm.IsQuoted() {
			c.maybeAddPropReference(nm.Key(), n)
		}
	case rhino.Token_SUPER:
		c.visitSuper(n)
	}
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.visitSuper(Node)
func (c *optimizeCallsRefCollector) visitSuper(superNode rhino.Node) {
	parent := superNode.GetParent()
	call, ok := ast.AsCall(parent)
	if !ok || !call.Callee().Is(superNode) {
		return
	}
	enclosingClass := ast.MustClass(rhino.CheckNotNull(NodeUtil_GetEnclosingClass(parent)))
	extendsExpr, hasExtends := enclosingClass.Extends()
	jutil.CheckStateMsg(hasExtends, "super call appears in class without extends clause")
	if name, isName := extendsExpr.AsName(); isName {
		c.maybeAddNameReferenceNode(name.Node(), superNode)
	} else if gp, isGetProp := extendsExpr.AsGetProp(); isGetProp {
		c.maybeAddPropReference(gp.Property(), superNode)
	}
}

func (c *optimizeCallsRefCollector) maybeAddNameReferenceNode(nameNode rhino.Node, n rhino.Node) {
	id := nameNode.GetStringID()
	if id <= rhino.RhinoStringPool_EmptyID {
		return
	}
	var_ := c.globalScope.GetVarByID(id)
	if var_ != nil && (c.outer.considerExterns || !var_.IsExtern()) {
		c.references.AddNameReference(nameNode.GetString(), n)
	}
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.maybeAddPropReference(String, Node)
func (c *optimizeCallsRefCollector) maybeAddPropReference(name string, n rhino.Node) {
	if !c.outer.considerExterns && c.externProps.Contains(name) {
		return
	}
	c.references.AddPropReference(name, n)
}

// OptimizeCalls_MayBeOptimizableName returns whether the provided name may be
// a candidate for call optimizations.
//
// Java: OptimizeCalls.mayBeOptimizableName(AbstractCompiler, String)
func OptimizeCalls_MayBeOptimizableName(compiler AbstractCompiler, name string) bool {
	if compiler.GetCodingConvention().IsExported(name) {
		return false
	}

	// Avoid modifying a few special case functions. Specifically, $jscomp.inherits to
	// recognize 'inherits' calls. (b/27244988)
	switch name {
	case NodeUtil_JSC_PROPERTY_NAME_FN, "inherits", "$jscomp$inherits", "goog$inherits":
		return false
	default:
		return true
	}
}

// OptimizeCalls_IsAllowedReference returns whether the reference is a known
// non-aliasing reference.
//
// Java: OptimizeCalls.isAllowedReference(Node)
func OptimizeCalls_IsAllowedReference(n rhino.Node) bool {
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// inspecting the properties is allowed.
		return ast.MustEnhancedFor(parent).Collection().Is(n)
	case rhino.Token_INSTANCEOF, rhino.Token_TYPEOF, rhino.Token_IN:
		return true
	case rhino.Token_GETELEM, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM:
		// Calls escape the "this" value. a.foo() aliases "a" as "this" but general
		// property references do not.
		if ast.MustNormalOrOptChainGet(parent).Object().Is(n) {
			grandparent := parent.GetParent()
			if grandparent.IsCall() || grandparent.IsOptChainCall() || grandparent.IsTaggedTemplateLit() {
				return false // `a.foo()` or `a?.foo()` or `a?.[foo]()` or `a.foo`bar``
			}
		}
		return true
	case rhino.Token_CLASS:
		cls := ast.MustClass(parent)
		if cls.NameNode() == n {
			// class Name { — this is a definition, not a read reference
			return false
		}
		jutil.CheckStateMsg1(cls.ExtendsNode() == n, "%s", parent)
		// We can optimize the constructor of the class being extended as long as all child
		// classes have explicit constructors, so we can see and update their `super()` calls.
		_, hasExplicitCtor := cls.Constructor()
		return hasExplicitCtor
	default:
		return NodeUtil_IsNameDeclaration(parent) && !n.HasChildren()
	}
}
