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
	"sync"
	"sync/atomic"

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

	references := acquireOptimizeCallsRefMap()
	defer releaseOptimizeCallsRefMap(references)
	sc := NewSyntacticScopeCreator(o.compiler)
	sc.EnableArenaPool()
	releaseEarly := len(o.passes) == 1
	if releaseEarly {
		if _, isPFI := o.passes[0].(*PureFunctionIdentifier); !isPFI {
			releaseEarly = false
		}
	}
	if !releaseEarly {
		defer sc.ReleaseArena()
	}
	cb := acquireOptimizeCallsCallback(o, references)
	b := NodeTraversal_Builder{
		compiler:     o.compiler,
		callback:     cb,
		scopeCreator: sc,
	}
	b.TraverseRoots(externs, root)
	releaseOptimizeCallsCallback(cb)
	if releaseEarly {
		references.globalScope = nil
		sc.ReleaseArena()
	}
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
	references.lastPropList = nil
}

const ocRefSlotChunkSize = 4096

type optimizeCalls_refListSlot struct {
	list jutil.ArrayList[*rhino.Node]
	buf  [8]*rhino.Node
}

// OptimizeCalls_ReferenceMap is a reference map for global symbols and
// properties.
type OptimizeCalls_ReferenceMap struct {
	globalScope  *Scope
	names        *jutil.LinkedHashMap[string, *jutil.ArrayList[*rhino.Node]]
	props        *jutil.LinkedHashMap[string, *jutil.ArrayList[*rhino.Node]]
	slotChunks   [][]optimizeCalls_refListSlot
	slotChunkIdx int
	slotUsed     int
	pooled       bool
	lastNameKey  string
	lastNameList *jutil.ArrayList[*rhino.Node]
	lastPropKey  string
	lastPropList *jutil.ArrayList[*rhino.Node]
}

var optimizeCallsRefMapPool sync.Pool
var optimizeCallsCallbackPool sync.Pool

var cachedOptimizeCallsRefMap atomic.Pointer[OptimizeCalls_ReferenceMap]
var cachedOptimizeCallsCallback atomic.Pointer[OptimizeCalls_ReferenceMapBuildingCallback]

func evictOptimizeCallsLoopCache() {
	cachedOptimizeCallsRefMap.Store(nil)
	cachedOptimizeCallsCallback.Store(nil)
}

func acquireOptimizeCallsRefMap() *OptimizeCalls_ReferenceMap {
	if m := cachedOptimizeCallsRefMap.Swap(nil); m != nil {
		return m
	}
	if v := optimizeCallsRefMapPool.Get(); v != nil {
		return v.(*OptimizeCalls_ReferenceMap)
	}
	return &OptimizeCalls_ReferenceMap{
		names: jutil.NewLinkedHashMapWithCapacity[string, *jutil.ArrayList[*rhino.Node]](4096),
		props: jutil.NewLinkedHashMapWithCapacity[string, *jutil.ArrayList[*rhino.Node]](4096),
	}
}

func releaseOptimizeCallsRefMap(m *OptimizeCalls_ReferenceMap) {
	total := m.names.Size() + m.props.Size()
	if total < 32 || total > 524288 {
		return
	}
	m.globalScope = nil
	m.lastNameKey = ""
	m.lastNameList = nil
	m.lastPropKey = ""
	m.lastPropList = nil
	m.names.Clear()
	m.props.Clear()
	for ci := 0; ci <= m.slotChunkIdx && ci < len(m.slotChunks); ci++ {
		chunk := m.slotChunks[ci]
		limit := len(chunk)
		if ci == m.slotChunkIdx {
			limit = m.slotUsed
		}
		for i := 0; i < limit; i++ {
			slot := &chunk[i]
			clear(slot.list.Slice())
			if cap(slot.list.Slice()) <= 65536 {
				slot.list.Clear()
			} else {
				slot.buf = [8]*rhino.Node{}
				slot.list.InitWithBuffer(slot.buf[:])
			}
		}
	}
	const maxRetainedChunks = 128
	if len(m.slotChunks) > maxRetainedChunks {
		clear(m.slotChunks[maxRetainedChunks:])
		m.slotChunks = m.slotChunks[:maxRetainedChunks]
	}
	m.slotChunkIdx = 0
	m.slotUsed = 0
	m.pooled = true
	if cachedOptimizeCallsRefMap.CompareAndSwap(nil, m) {
		return
	}
	optimizeCallsRefMapPool.Put(m)
}

func acquireOptimizeCallsCallback(outer *OptimizeCalls, references *OptimizeCalls_ReferenceMap) *OptimizeCalls_ReferenceMapBuildingCallback {
	c := cachedOptimizeCallsCallback.Swap(nil)
	if c == nil {
		if v := optimizeCallsCallbackPool.Get(); v != nil {
			c = v.(*OptimizeCalls_ReferenceMapBuildingCallback)
		}
	}
	if c != nil {
		c.outer = outer
		c.codingConvention = outer.compiler.GetCodingConvention()
		c.externProps = OptimizeCalls_SafeSet(outer.compiler.GetExternProperties())
		c.references = references
		return c
	}
	return NewOptimizeCalls_ReferenceMapBuildingCallback(outer, references)
}

func releaseOptimizeCallsCallback(c *OptimizeCalls_ReferenceMapBuildingCallback) {
	c.outer = nil
	c.codingConvention = nil
	c.externProps = nil
	c.references = nil
	c.globalScope = nil
	c.lastSkippedName = ""
	c.lastExternProp = ""
	c.hasExternProp = false
	clear(c.nameCache[:])
	clear(c.propCache[:])
	if cachedOptimizeCallsCallback.CompareAndSwap(nil, c) {
		return
	}
	optimizeCallsCallbackPool.Put(c)
}

// Java: OptimizeCalls.ReferenceMap()
func NewOptimizeCalls_ReferenceMap() *OptimizeCalls_ReferenceMap {
	return &OptimizeCalls_ReferenceMap{
		names: jutil.NewLinkedHashMapWithCapacity[string, *jutil.ArrayList[*rhino.Node]](256),
		props: jutil.NewLinkedHashMapWithCapacity[string, *jutil.ArrayList[*rhino.Node]](512),
	}
}

func (m *OptimizeCalls_ReferenceMap) allocRefList() *jutil.ArrayList[*rhino.Node] {
	if len(m.slotChunks) == 0 {
		m.slotChunks = append(m.slotChunks, make([]optimizeCalls_refListSlot, ocRefSlotChunkSize))
		m.slotChunkIdx = 0
		m.slotUsed = 0
	} else if m.slotUsed == len(m.slotChunks[m.slotChunkIdx]) {
		m.slotChunkIdx++
		if m.slotChunkIdx == len(m.slotChunks) {
			m.slotChunks = append(m.slotChunks, make([]optimizeCalls_refListSlot, ocRefSlotChunkSize))
		}
		m.slotUsed = 0
	}
	slot := &m.slotChunks[m.slotChunkIdx][m.slotUsed]
	m.slotUsed++
	if cap(slot.list.Slice()) == 0 {
		slot.list.InitWithBuffer(slot.buf[:])
	}
	return &slot.list
}

func optimizeCallsAddNode(list *jutil.ArrayList[*rhino.Node], n *rhino.Node) {
	s := list.Slice()
	if len(s) == cap(s) {
		newCap := max(cap(s)*4, 32)
		next := make([]*rhino.Node, newCap)
		list.InitWithBuffer(next)
		list.AddAll(s...)
		if cap(s) == 8 {
			clear(s)
		}
	}
	list.Add(n)
}

// Java: OptimizeCalls.ReferenceMap.addReference(LinkedHashMap<String, ArrayList<Node>>, String, Node)
func (m *OptimizeCalls_ReferenceMap) AddReference(data *jutil.LinkedHashMap[string, *jutil.ArrayList[*rhino.Node]], name string, n *rhino.Node) {
	refs := data.Get(name)
	if refs == nil {
		refs = m.allocRefList()
		data.PutNew(name, refs)
	}
	optimizeCallsAddNode(refs, n)
}

// Java: OptimizeCalls.ReferenceMap.addNameReference(String, Node)
func (m *OptimizeCalls_ReferenceMap) AddNameReference(name string, n *rhino.Node) {
	if m.lastNameList != nil && name == m.lastNameKey {
		optimizeCallsAddNode(m.lastNameList, n)
		return
	}
	refs := m.names.Get(name)
	if refs == nil {
		refs = m.allocRefList()
		m.names.PutNew(name, refs)
	}
	m.lastNameKey = name
	m.lastNameList = refs
	optimizeCallsAddNode(refs, n)
}

// Java: OptimizeCalls.ReferenceMap.addPropReference(String, Node)
func (m *OptimizeCalls_ReferenceMap) AddPropReference(name string, n *rhino.Node) {
	if m.lastPropList != nil && name == m.lastPropKey {
		optimizeCallsAddNode(m.lastPropList, n)
		return
	}
	refs := m.props.Get(name)
	if refs == nil {
		refs = m.allocRefList()
		m.props.PutNew(name, refs)
	}
	m.lastPropKey = name
	m.lastPropList = refs
	optimizeCallsAddNode(refs, n)
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

var optimizeCalls_emptyFunctionNodes = jutil.NewImmutableListMultimapBuilder[*rhino.Node, *rhino.Node]().Build()

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
	var result *jutil.ImmutableListMultimapBuilder[*rhino.Node, *rhino.Node]
	for _, def := range definitionSites {
		optimizeCalls_putDefinitionFunctionNodesFor(&result, def)
	}
	if result == nil {
		return optimizeCalls_emptyFunctionNodes
	}
	return result.Build()
}

func optimizeCalls_putFn(result **jutil.ImmutableListMultimapBuilder[*rhino.Node, *rhino.Node], def *rhino.Node, fn *rhino.Node) {
	if *result == nil {
		*result = jutil.NewImmutableListMultimapBuilder[*rhino.Node, *rhino.Node]()
	}
	(*result).Put(def, fn)
}

func optimizeCalls_putDefinitionFunctionNodesFor(result **jutil.ImmutableListMultimapBuilder[*rhino.Node, *rhino.Node], definitionSite *rhino.Node) {
	if definitionSite.IsGetterDef() || definitionSite.IsSetterDef() {
		return
	}

	parent := definitionSite.GetParent()
	if parent == nil {
		return
	}

	switch parent.GetToken() {
	case rhino.Token_CLASS:
		if definitionSite.IsFirstChildOf(parent) {
			constructorFnDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(parent)
			if constructorFnDef != nil {
				optimizeCalls_putFn(result, definitionSite, constructorFnDef.GetOnlyChild())
			}
		}
	case rhino.Token_FUNCTION:
		optimizeCalls_putFn(result, definitionSite, parent)
	case rhino.Token_CLASS_MEMBERS:
		if definitionSite.IsMemberFunctionDef() {
			optimizeCalls_putFn(result, definitionSite, definitionSite.GetLastChild())
		} else {
			jutil.CheckArgumentMsg1(definitionSite.IsMemberFieldDef(), "%s", definitionSite)
			value := definitionSite.GetFirstChild()
			if value != nil {
				optimizeCalls_putValueFunctionNodes(result, definitionSite, value)
			}
		}
	case rhino.Token_OBJECTLIT:
		jutil.CheckArgumentMsg1(
			definitionSite.IsStringKey() || definitionSite.IsMemberFunctionDef(), //
			"%s", definitionSite)
		optimizeCalls_putValueFunctionNodes(result, definitionSite, definitionSite.GetLastChild())
	case rhino.Token_ASSIGN:
		target := parent.GetFirstChild()
		value := parent.GetLastChild()
		if definitionSite == target {
			optimizeCalls_putValueFunctionNodes(result, definitionSite, value)
		}
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		if definitionSite.IsName() && definitionSite.HasChildren() {
			optimizeCalls_putValueFunctionNodes(result, definitionSite, definitionSite.GetFirstChild())
		}
	default:
	}
}

func optimizeCalls_putValueFunctionNodes(result **jutil.ImmutableListMultimapBuilder[*rhino.Node, *rhino.Node], def *rhino.Node, n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		constructorFnDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
		if constructorFnDef != nil {
			optimizeCalls_putFn(result, def, constructorFnDef.GetOnlyChild())
		}
	case rhino.Token_FUNCTION:
		optimizeCalls_putFn(result, def, n)
	case rhino.Token_HOOK:
		optimizeCalls_putValueFunctionNodes(result, def, n.GetSecondChild())
		optimizeCalls_putValueFunctionNodes(result, def, n.GetLastChild())
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		optimizeCalls_putValueFunctionNodes(result, def, n.GetFirstChild())
		optimizeCalls_putValueFunctionNodes(result, def, n.GetLastChild())
	case rhino.Token_CAST, rhino.Token_COMMA:
		optimizeCalls_putValueFunctionNodes(result, def, n.GetLastChild())
	default:
	}
}

// OptimizeCalls_ReferenceMap_AppendFunctionNodes appends all function nodes
// associated with definitionSites to dst and returns the updated slice.
func OptimizeCalls_ReferenceMap_AppendFunctionNodes(dst []*rhino.Node, definitionSites []*rhino.Node) []*rhino.Node {
	for _, def := range definitionSites {
		dst = optimizeCalls_appendDefinitionFunctionNodesFor(dst, def)
	}
	return dst
}

// OptimizeCalls_ReferenceMap_GetSingleFunctionNode returns the single function
// node associated with definitionSites if there is exactly 1, or nil if 0 or > 1.
func OptimizeCalls_ReferenceMap_GetSingleFunctionNode(definitionSites []*rhino.Node) *rhino.Node {
	fn, count := optimizeCalls_getSingleFunctionNode(definitionSites)
	if count == 1 {
		return fn
	}
	return nil
}

func optimizeCalls_getSingleFunctionNode(definitionSites []*rhino.Node) (*rhino.Node, int) {
	var single *rhino.Node
	count := 0
	for _, def := range definitionSites {
		single, count = optimizeCalls_collectSingleDefinitionFunctionNode(def, single, count)
		if count > 1 {
			return nil, count
		}
	}
	return single, count
}

func optimizeCalls_collectSingleDefinitionFunctionNode(definitionSite *rhino.Node, single *rhino.Node, count int) (*rhino.Node, int) {
	if definitionSite.IsGetterDef() || definitionSite.IsSetterDef() {
		return single, count
	}
	parent := definitionSite.GetParent()
	if parent == nil {
		return single, count
	}
	switch parent.GetToken() {
	case rhino.Token_CLASS:
		if definitionSite.IsFirstChildOf(parent) {
			constructorFnDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(parent)
			if constructorFnDef != nil {
				return constructorFnDef.GetOnlyChild(), count + 1
			}
		}
	case rhino.Token_FUNCTION:
		return parent, count + 1
	case rhino.Token_CLASS_MEMBERS:
		if definitionSite.IsMemberFunctionDef() {
			return definitionSite.GetLastChild(), count + 1
		}
		jutil.CheckArgumentMsg1(definitionSite.IsMemberFieldDef(), "%s", definitionSite)
		if value := definitionSite.GetFirstChild(); value != nil {
			return optimizeCalls_collectSingleValueFunctionNode(value, single, count)
		}
	case rhino.Token_OBJECTLIT:
		jutil.CheckArgumentMsg1(
			definitionSite.IsStringKey() || definitionSite.IsMemberFunctionDef(),
			"%s", definitionSite)
		return optimizeCalls_collectSingleValueFunctionNode(definitionSite.GetLastChild(), single, count)
	case rhino.Token_ASSIGN:
		target := parent.GetFirstChild()
		value := parent.GetLastChild()
		if definitionSite == target {
			return optimizeCalls_collectSingleValueFunctionNode(value, single, count)
		}
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		if definitionSite.IsName() && definitionSite.HasChildren() {
			return optimizeCalls_collectSingleValueFunctionNode(definitionSite.GetFirstChild(), single, count)
		}
	default:
	}
	return single, count
}

func optimizeCalls_collectSingleValueFunctionNode(n *rhino.Node, single *rhino.Node, count int) (*rhino.Node, int) {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		constructorFnDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
		if constructorFnDef != nil {
			return constructorFnDef.GetOnlyChild(), count + 1
		}
	case rhino.Token_FUNCTION:
		return n, count + 1
	case rhino.Token_HOOK:
		single, count = optimizeCalls_collectSingleValueFunctionNode(n.GetSecondChild(), single, count)
		if count > 1 {
			return nil, count
		}
		return optimizeCalls_collectSingleValueFunctionNode(n.GetLastChild(), single, count)
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		single, count = optimizeCalls_collectSingleValueFunctionNode(n.GetFirstChild(), single, count)
		if count > 1 {
			return nil, count
		}
		return optimizeCalls_collectSingleValueFunctionNode(n.GetLastChild(), single, count)
	case rhino.Token_CAST, rhino.Token_COMMA:
		return optimizeCalls_collectSingleValueFunctionNode(n.GetLastChild(), single, count)
	default:
	}
	return single, count
}

// OptimizeCalls_ReferenceMap_DefinitionFunctionNodesFor collects potential
// definition FUNCTIONs associated with a method definition site.
//
// See getFunctionNodes.
//
// Java: OptimizeCalls.ReferenceMap.definitionFunctionNodesFor(Node)
func OptimizeCalls_ReferenceMap_DefinitionFunctionNodesFor(definitionSite *rhino.Node) []*rhino.Node {
	return optimizeCalls_appendDefinitionFunctionNodesFor(nil, definitionSite)
}

func optimizeCalls_appendDefinitionFunctionNodesFor(fns []*rhino.Node, definitionSite *rhino.Node) []*rhino.Node {
	if definitionSite.IsGetterDef() || definitionSite.IsSetterDef() {
		// TODO(nickreid): Support getters and setters. Ignore them for now since they aren't
		// "called".
		return fns
	}

	// Ignore detached nodes.
	parent := definitionSite.GetParent()
	if parent == nil {
		return fns
	}

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
				fns = optimizeCalls_appendValueFunctionNodes(fns, value)
			}
		}
	case rhino.Token_OBJECTLIT:
		jutil.CheckArgumentMsg1(
			definitionSite.IsStringKey() || definitionSite.IsMemberFunctionDef(), //
			"%s", definitionSite)
		fns = optimizeCalls_appendValueFunctionNodes(fns, definitionSite.GetLastChild())
	case rhino.Token_ASSIGN:
		// Only a candidate if the assign isn't consumed.
		target := parent.GetFirstChild()
		value := parent.GetLastChild()
		if definitionSite == target {
			fns = optimizeCalls_appendValueFunctionNodes(fns, value)
		}
	case rhino.Token_CONST, rhino.Token_LET, rhino.Token_VAR:
		if definitionSite.IsName() && definitionSite.HasChildren() {
			fns = optimizeCalls_appendValueFunctionNodes(fns, definitionSite.GetFirstChild())
		}
	default:
	}
	return fns
}

func optimizeCalls_appendValueFunctionNodes(fns []*rhino.Node, n *rhino.Node) []*rhino.Node {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		constructorFnDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
		if constructorFnDef != nil {
			fns = append(fns, constructorFnDef.GetOnlyChild())
		}
	case rhino.Token_FUNCTION:
		fns = append(fns, n)
	case rhino.Token_HOOK:
		fns = optimizeCalls_appendValueFunctionNodes(fns, n.GetSecondChild())
		fns = optimizeCalls_appendValueFunctionNodes(fns, n.GetLastChild())
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		fns = optimizeCalls_appendValueFunctionNodes(fns, n.GetFirstChild())
		fns = optimizeCalls_appendValueFunctionNodes(fns, n.GetLastChild())
	case rhino.Token_CAST, rhino.Token_COMMA:
		fns = optimizeCalls_appendValueFunctionNodes(fns, n.GetLastChild())
	default:
	}
	return fns
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

type ocRefCacheEntry struct {
	key  string
	list *jutil.ArrayList[*rhino.Node]
	skip bool
}

func ocCacheIdx(s string) uint32 {
	n := len(s)
	h := uint32(n)*2654435761 ^ uint32(s[0])*374761393 ^ uint32(s[n-1])*668265263
	if n > 2 {
		h ^= uint32(s[n>>1]) * 2246822519
	}
	return (h ^ (h >> 9)) & 511
}

// OptimizeCalls_ReferenceMapBuildingCallback is the inner class
// ReferenceMapBuildingCallback.
type OptimizeCalls_ReferenceMapBuildingCallback struct {
	outer            *OptimizeCalls
	codingConvention CodingConvention
	externProps      *jutil.LinkedHashSet[string]
	references       *OptimizeCalls_ReferenceMap
	globalScope      *Scope
	lastSkippedName  string
	lastExternProp   string
	hasExternProp    bool
	nameCache        [512]ocRefCacheEntry
	propCache        [512]ocRefCacheEntry
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback(ReferenceMap)
func NewOptimizeCalls_ReferenceMapBuildingCallback(outer *OptimizeCalls, references *OptimizeCalls_ReferenceMap) *OptimizeCalls_ReferenceMapBuildingCallback {
	c := &OptimizeCalls_ReferenceMapBuildingCallback{
		outer:            outer,
		codingConvention: outer.compiler.GetCodingConvention(),
		externProps:      OptimizeCalls_SafeSet(outer.compiler.GetExternProperties()),
		references:       references,
	}
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
		if c.codingConvention.IsPropertyRenameFunction(fnName) {
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
	if name == "" || name == c.lastSkippedName {
		return
	}
	if c.references.lastNameList != nil && name == mLastNameKey(c.references) {
		optimizeCallsAddNode(c.references.lastNameList, n)
		return
	}
	idx := ocCacheIdx(name)
	if e := &c.nameCache[idx]; e.key == name {
		if e.skip {
			c.lastSkippedName = name
			return
		}
		c.references.lastNameKey = name
		c.references.lastNameList = e.list
		optimizeCallsAddNode(e.list, n)
		return
	}
	// TODO(b/129503101): Why are we limiting ourselves to global names?
	var_ := c.globalScope.GetVar(name)
	if var_ != nil && (c.outer.considerExterns || !var_.IsExtern()) {
		// As every name declaration is unique due to normalizations, it is only necessary to build
		// the global scope and ask it if it knows about a name as it can never be shadowed.
		c.references.AddNameReference(name, n)
		c.nameCache[idx] = ocRefCacheEntry{key: name, list: c.references.lastNameList}
	} else {
		c.lastSkippedName = name
		c.nameCache[idx] = ocRefCacheEntry{key: name, skip: true}
	}
}

func mLastNameKey(m *OptimizeCalls_ReferenceMap) string {
	return m.lastNameKey
}

// Java: OptimizeCalls.ReferenceMapBuildingCallback.maybeAddPropReference(String, Node)
func (c *OptimizeCalls_ReferenceMapBuildingCallback) MaybeAddPropReference(name string, n *rhino.Node) {
	if c.references.lastPropList != nil && name == c.references.lastPropKey {
		optimizeCallsAddNode(c.references.lastPropList, n)
		return
	}
	if name != "" {
		idx := ocCacheIdx(name)
		if e := &c.propCache[idx]; e.key == name {
			if e.skip {
				c.lastExternProp = name
				c.hasExternProp = true
				return
			}
			c.references.lastPropKey = name
			c.references.lastPropList = e.list
			optimizeCallsAddNode(e.list, n)
			return
		}
		if !c.outer.considerExterns {
			if c.hasExternProp && name == c.lastExternProp {
				return
			}
			if c.externProps.Contains(name) {
				c.lastExternProp = name
				c.hasExternProp = true
				c.propCache[idx] = ocRefCacheEntry{key: name, skip: true}
				return
			}
		}
		c.references.AddPropReference(name, n)
		c.propCache[idx] = ocRefCacheEntry{key: name, list: c.references.lastPropList}
		return
	}
	if !c.outer.considerExterns {
		if c.hasExternProp && name == c.lastExternProp {
			return
		}
		if c.externProps.Contains(name) {
			c.lastExternProp = name
			c.hasExternProp = true
			return
		}
	}
	c.references.AddPropReference(name, n)
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
		if !c.references.pooled {
			varCount := c.globalScope.GetVarCount()
			if varCount > 4096 && c.references.names.IsEmpty() {
				c.references.names = jutil.NewLinkedHashMapWithCapacity[string, *jutil.ArrayList[*rhino.Node]](varCount)
			}
			if varCount > 1024 && c.references.props.IsEmpty() {
				c.references.props = jutil.NewLinkedHashMapWithCapacity[string, *jutil.ArrayList[*rhino.Node]](max(varCount*4, 16384))
			}
		}
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
