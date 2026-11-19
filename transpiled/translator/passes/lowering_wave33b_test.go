package passes_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave33B_PurgeTrimming_DeadPropSlab_And_TypedScopeVarSlabs(t *testing.T) {
	tmpDir := t.TempDir()
	mustWrite := func(rel, content string) string {
		t.Helper()
		full := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return full
	}

	mustWrite("go.mod", "module example.com/w33btest\n\ngo 1.24\n")

	mustWrite("gen/rhino/rhino.go", `package rhino

type Node uint32

func (n Node) IsJavaNil_() bool     { return n == 0 }
func (n Node) IsRoot() bool         { return n == 1 }
func (n Node) IsFunction() bool     { return n >= 10 && n < 20 }
func (n Node) IsGetProp() bool      { return n >= 20 && n < 30 }
func (n Node) IsQualifiedName() bool { return n >= 20 }
func (n Node) HasChildren() bool    { return n != 0 }
func (n Node) GetFirstChild() Node  {
	if n >= 20 && n < 30 {
		return 30
	}
	return 0
}
func (n Node) GetString() string {
	if n == 21 {
		return "propA"
	}
	if n == 22 {
		return "propB"
	}
	return "obj"
}
func (n Node) GetQualifiedName() string {
	if n == 21 {
		return "obj.propA"
	}
	if n == 22 {
		return "obj.propB"
	}
	return "obj"
}

type JSType interface {
	IsJavaNil_() bool
}
`)

	mustWrite("gen/j2gort/j2gort.go", `package j2gort

import "sync/atomic"

const NullString = "\x00"

var (
	PurgeCalls           atomic.Int64
	SlabAllocCalls       atomic.Int64
	PooledAllocCalls     atomic.Int64
	PooledResetCalls     atomic.Int64
	PooledReleaseCalls   atomic.Int64
)

func ResetWave33BMetrics() {
	PurgeCalls.Store(0)
	SlabAllocCalls.Store(0)
	PooledAllocCalls.Store(0)
	PooledResetCalls.Store(0)
	PooledReleaseCalls.Store(0)
}

func PurgeAllRuntimePools() { PurgeCalls.Add(1) }
func TrimStaleSlabPools(int) {}

type Slab[T any] struct {
	buf     []T
	nextCap int32
}

func (s *Slab[T]) Alloc() *T {
	SlabAllocCalls.Add(1)
	if i := len(s.buf); i < cap(s.buf) {
		b := s.buf[:i+1]
		s.buf = b
		return &b[i]
	}
	n := s.nextCap
	if n == 0 {
		n = 16
		s.nextCap = 32
	} else if n < 256 {
		s.nextCap = n * 2
	} else {
		n = 256
	}
	c := make([]T, 1, n)
	s.buf = c
	return &c[0]
}

func (s *Slab[T]) Release() {
	if s == nil {
		return
	}
	clear(s.buf)
	s.buf = nil
	s.nextCap = 0
}

type PooledSlab[T any] struct {
	items []*T
}

func (s *PooledSlab[T]) Alloc() *T {
	PooledAllocCalls.Add(1)
	v := new(T)
	s.items = append(s.items, v)
	return v
}

func (s *PooledSlab[T]) Reset() {
	if s == nil {
		return
	}
	PooledResetCalls.Add(1)
	for _, it := range s.items {
		var zero T
		*it = zero
	}
	s.items = s.items[:0]
}

func (s *PooledSlab[T]) Release() {
	if s == nil {
		return
	}
	PooledReleaseCalls.Add(1)
	for _, it := range s.items {
		var zero T
		*it = zero
	}
	s.items = nil
}

func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	if s == nil {
		return
	}
	PooledReleaseCalls.Add(1)
	for _, it := range s.items {
		fn(it)
		var zero T
		*it = zero
	}
	s.items = nil
}

type OrderedMap[K comparable, V any] struct {
	m map[K]V
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V)}
}

func NewOrderedMapCap[K comparable, V any](int32) OrderedMap[K, V] {
	return OrderedMap[K, V]{m: make(map[K]V)}
}

func (o *OrderedMap[K, V]) Get(k K) V {
	if o == nil || o.m == nil {
		var z V
		return z
	}
	return o.m[k]
}

func (o *OrderedMap[K, V]) Put(k K, v V) {
	if o.m == nil {
		o.m = make(map[K]V)
	}
	o.m[k] = v
}

func (o *OrderedMap[K, V]) Clear() {
	if o != nil && o.m != nil {
		clear(o.m)
	}
}

func (o *OrderedMap[K, V]) Size() int32 {
	if o == nil {
		return 0
	}
	return int32(len(o.m))
}

type MapValCursor[V any] struct {
	vals []V
	i    int
}

func (c *MapValCursor[V]) HasNext() bool { return c.i < len(c.vals) }
func (c *MapValCursor[V]) Next() V       { v := c.vals[c.i]; c.i++; return v }

func (o *OrderedMap[K, V]) IterValues() *MapValCursor[V] {
	var vs []V
	for _, v := range o.m {
		vs = append(vs, v)
	}
	return &MapValCursor[V]{vals: vs}
}

type OrderedSet[T comparable] struct {
	m map[T]struct{}
}

func NewOrderedSet[T comparable]() OrderedSet[T] {
	return OrderedSet[T]{m: make(map[T]struct{})}
}

func (s *OrderedSet[T]) Add(v T) bool {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	if _, ok := s.m[v]; ok {
		return false
	}
	s.m[v] = struct{}{}
	return true
}

type Deque[T any] struct {
	s []T
}

func (d *Deque[T]) AddLast(v T)   { d.s = append(d.s, v) }
func (d *Deque[T]) IsEmpty() bool { return len(d.s) == 0 }
func (d *Deque[T]) Size() int32   { return int32(len(d.s)) }

type Ext_java_util_Set interface {
	IsEmpty() bool
}

type Ext_java_util_function_Function interface {
	Apply(any) any
}

type fnImpl struct {
	fn func(any) any
}

func (f *fnImpl) Apply(v any) any { return f.fn(v) }

func NewLambda_Ext_java_util_function_Function(fn func(any) any) Ext_java_util_function_Function {
	return &fnImpl{fn: fn}
}

func OrderedMapComputeIfAbsentFn[K comparable, V comparable](m *OrderedMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
	if v := m.Get(k); v != nullV {
		return v
	}
	raw := fn.Apply(k)
	if raw == nil {
		return nullV
	}
	nv := raw.(V)
	m.Put(k, nv)
	return nv
}

func CastString(v any) string        { return v.(string) }
func CastPtr[T any](v any) *T        { p, _ := v.(*T); return p }
func PtrToAny[T any](p *T) any       { if p == nil { return nil }; return p }
func CheckNotNullPtr[T any](p *T) *T { return p }
`)

	// Part A: Phase-boundary PurgeAllRuntimePools() trimming.
	mustWrite("gen/jscomp/phases.go", `package jscomp

import (
	"example.com/w33btest/gen/j2gort"
	pkg_rhino "example.com/w33btest/gen/rhino"
)

type AbstractCompiler interface {
	IsJavaNil_() bool
}

type Compiler struct{}

func (c *Compiler) IsJavaNil_() bool { return c == nil }

func (this_ *Compiler) ParseInputs() pkg_rhino.Node {
	defer j2gort.PurgeAllRuntimePools()
	j2gort.TrimStaleSlabPools(1)
	return 1
}

func (this_ *Compiler) PerformTranspilationAndOptimizations(seg int) {
	_ = seg
	j2gort.PurgeAllRuntimePools()
}

type TypeInferencePass struct {
	DidWork bool
}

func (this_ *TypeInferencePass) InferAllScopes(inferenceRoot pkg_rhino.Node) *TypedScope {
	_ = inferenceRoot
	j2gort.PurgeAllRuntimePools()
	defer j2gort.PurgeAllRuntimePools()
	this_.DidWork = true
	return nil
}

type TypedAstSerializer struct{}

func (this_ *TypedAstSerializer) SerializeRoots(externsRoot pkg_rhino.Node, jsRoot pkg_rhino.Node) {
	_ = externsRoot
	_ = jsRoot
	defer j2gort.PurgeAllRuntimePools()
}

type ProcessDefines_Mode uint8

const (
	modeCheck ProcessDefines_Mode = 1
	modeOpt   ProcessDefines_Mode = 2
)

func StaticGet_ProcessDefines_Mode_CHECK() ProcessDefines_Mode { return modeCheck }

type ProcessDefines struct {
	mode ProcessDefines_Mode
}

func (this_ *ProcessDefines) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.mode == StaticGet_ProcessDefines_Mode_CHECK() {
		defer j2gort.PurgeAllRuntimePools()
	}
}

type VariableReferenceCheck struct{}

func (this_ *VariableReferenceCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer j2gort.PurgeAllRuntimePools()
}
`)

	// Part B (Shape 1): DeadPropertyAssignmentElimination with FunctionVisitor + FindCandidateAssignmentTraversal + OrderedMapComputeIfAbsentFn.
	mustWrite("gen/jscomp/dead_property_assignment_elimination.go", `package jscomp

import (
	"sync/atomic"

	"example.com/w33btest/gen/j2gort"
	pkg_rhino "example.com/w33btest/gen/rhino"
)

type NodeTraversal struct{}

type NodeTraversal_Callback interface {
	IsJavaNil_() bool
}

func NodeTraversal_Traverse(c AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	if fv, ok := cb.(*DeadPropertyAssignmentElimination_FunctionVisitor); ok {
		fv.EnterChangedScopeRoot(c, 10)
		fv.EnterChangedScopeRoot(c, 11)
		return
	}
	if tr, ok := cb.(*DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal); ok {
		tr.VisitProp(21)
		tr.VisitProp(22)
		tr.VisitProp(21)
	}
	_ = root
}

var lambdaSingleton_dead_property_assignment_elimination_2 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_dead_property_assignment_elimination_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewDeadPropertyAssignmentElimination_Property(j2gort.CastString(mrp0_))
	})
}

type DeadPropertyAssignmentElimination struct {
	compiler    AbstractCompiler
	VisitedProps int
	idHash_     atomic.Int32
}

func (*DeadPropertyAssignmentElimination) Implements_jscomp_CompilerPass() {}

func NewDeadPropertyAssignmentElimination(compiler AbstractCompiler) *DeadPropertyAssignmentElimination {
	this_ := &DeadPropertyAssignmentElimination{}
	InitDeadPropertyAssignmentElimination(this_, compiler)
	return this_
}

func InitDeadPropertyAssignmentElimination(this_ *DeadPropertyAssignmentElimination, compiler AbstractCompiler) {
	this_.compiler = compiler
}

func (this_ *DeadPropertyAssignmentElimination) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.compiler == nil {
		return
	}
	NodeTraversal_Traverse(this_.compiler, root, NewDeadPropertyAssignmentElimination_FunctionVisitor(nil))
}

type DeadPropertyAssignmentElimination_FunctionVisitor struct {
	skiplistedPropNames j2gort.Ext_java_util_Set
}

func (this_ *DeadPropertyAssignmentElimination_FunctionVisitor) IsJavaNil_() bool { return this_ == nil }

func NewDeadPropertyAssignmentElimination_FunctionVisitor(skiplistedPropNames j2gort.Ext_java_util_Set) *DeadPropertyAssignmentElimination_FunctionVisitor {
	this_ := &DeadPropertyAssignmentElimination_FunctionVisitor{}
	InitDeadPropertyAssignmentElimination_FunctionVisitor(this_, skiplistedPropNames)
	return this_
}

func InitDeadPropertyAssignmentElimination_FunctionVisitor(this_ *DeadPropertyAssignmentElimination_FunctionVisitor, skiplistedPropNames j2gort.Ext_java_util_Set) {
	this_.skiplistedPropNames = skiplistedPropNames
}

var TotalUniquePropsObserved atomic.Int64

func (this_ *DeadPropertyAssignmentElimination_FunctionVisitor) EnterChangedScopeRoot(compiler AbstractCompiler, root pkg_rhino.Node) {
	_ = compiler
	if !root.IsFunction() {
		return
	}
	var val_1 DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal
	InitDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(&val_1, this_.skiplistedPropNames)
	var traversal *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal = &val_1
	_ = traversal
	NodeTraversal_Traverse(compiler, root, traversal)
	cur_2 := traversal.propertyMap.IterValues()
	for cur_2.HasNext() {
		prop := cur_2.Next()
		if prop != nil && prop.Name() != "" {
			TotalUniquePropsObserved.Add(1)
		}
	}
}

type DeadPropertyAssignmentElimination_Property struct {
	name     string
	writes   j2gort.Deque[pkg_rhino.Node]
	children j2gort.OrderedSet[*DeadPropertyAssignmentElimination_Property]
	idHash_  atomic.Int32
}

func (this_ *DeadPropertyAssignmentElimination_Property) initFields_DeadPropertyAssignmentElimination_Property() {
	this_.name = j2gort.NullString
	this_.writes = j2gort.Deque[pkg_rhino.Node]{}
	this_.children = j2gort.NewOrderedSet[*DeadPropertyAssignmentElimination_Property]()
}

func NewDeadPropertyAssignmentElimination_Property(name string) *DeadPropertyAssignmentElimination_Property {
	this_ := &DeadPropertyAssignmentElimination_Property{}
	InitDeadPropertyAssignmentElimination_Property(this_, name)
	return this_
}

func InitDeadPropertyAssignmentElimination_Property(this_ *DeadPropertyAssignmentElimination_Property, name string) {
	this_.initFields_DeadPropertyAssignmentElimination_Property()
	this_.name = name
}

func (this_ *DeadPropertyAssignmentElimination_Property) Name() string { return this_.name }

type DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal struct {
	propertyMap         j2gort.OrderedMap[string, *DeadPropertyAssignmentElimination_Property]
	skiplistedPropNames j2gort.Ext_java_util_Set
	idHash_             atomic.Int32
}

func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) IsJavaNil_() bool {
	return this_ == nil
}

func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) initFields_DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal() {
	this_.propertyMap = j2gort.NewOrderedMap[string, *DeadPropertyAssignmentElimination_Property]()
}

func NewDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(skiplistedPropNames j2gort.Ext_java_util_Set) *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal {
	this_ := &DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal{}
	InitDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(this_, skiplistedPropNames)
	return this_
}

func InitDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal, skiplistedPropNames j2gort.Ext_java_util_Set) {
	this_.initFields_DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal()
	this_.skiplistedPropNames = skiplistedPropNames
}

func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) getOrCreateProperty(propNode pkg_rhino.Node) *DeadPropertyAssignmentElimination_Property {
	if !propNode.IsQualifiedName() {
		return nil
	}
	var propName string
	if propNode.IsGetProp() {
		propName = propNode.GetString()
	} else {
		propName = propNode.GetQualifiedName()
	}
	var property *DeadPropertyAssignmentElimination_Property = j2gort.OrderedMapComputeIfAbsentFn(&this_.propertyMap, propName, nil, lambdaSingleton_dead_property_assignment_elimination_2)
	if propNode.IsGetProp() {
		parentProperty := this_.getOrCreateProperty(propNode.GetFirstChild())
		if parentProperty != nil {
			parentProperty.children.Add(property)
		}
	}
	return property
}

func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) VisitProp(n pkg_rhino.Node) {
	p := this_.getOrCreateProperty(n)
	if p != nil {
		p.writes.AddLast(n)
	}
}
`)

	// Part C: TypedScopeCreator + TypedScope + TypedVar + TypeInference.
	mustWrite("gen/jscomp/typed_scope.go", `package jscomp

import (
	"example.com/w33btest/gen/j2gort"
	pkg_rhino "example.com/w33btest/gen/rhino"
)

type AbstractScope interface {
	IsJavaNil_() bool
}

type AbstractScopeBase struct {
	self_    AbstractScope
	rootNode pkg_rhino.Node
}

func (b *AbstractScopeBase) BindSelf_(s any) { b.self_ = s.(AbstractScope) }

type Module struct{}
type CompilerInput struct{}

type TypedScope struct {
	AbstractScopeBase
	parent   *TypedScope
	vars     map[string]*TypedVar
	depth    int32
	isBottom bool
}

func (this_ *TypedScope) BindSelf_(s any)   { this_.AbstractScopeBase.BindSelf_(s) }
func (this_ *TypedScope) IsJavaNil_() bool  { return this_ == nil }
func (this_ *TypedScope) initFields_TypedScope() {}

func NewTypedScope_2_TypedScope_Node(parent *TypedScope, rootNode pkg_rhino.Node) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_2_TypedScope_Node(this_, parent, rootNode)
	return this_
}

func InitTypedScope_2_TypedScope_Node(this_ *TypedScope, parent *TypedScope, rootNode pkg_rhino.Node) {
	InitTypedScope_4_TypedScope_Node_slice_string_Module(this_, parent, rootNode, nil, nil)
}

func NewTypedScope_4_TypedScope_Node_slice_string_Module(parent *TypedScope, rootNode pkg_rhino.Node, reservedNamesOwned []string, module *Module) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_4_TypedScope_Node_slice_string_Module(this_, parent, rootNode, reservedNamesOwned, module)
	return this_
}

func InitTypedScope_4_TypedScope_Node_slice_string_Module(this_ *TypedScope, parent *TypedScope, rootNode pkg_rhino.Node, reservedNamesOwned []string, module *Module) {
	_ = reservedNamesOwned
	_ = module
	this_.initFields_TypedScope()
	this_.rootNode = rootNode
	this_.parent = parent
	this_.depth = parent.depth + 1
	this_.isBottom = false
}

func NewTypedScope_2_Node_bool(rootNode pkg_rhino.Node, isBottom bool) *TypedScope {
	this_ := &TypedScope{}
	this_.BindSelf_(this_)
	InitTypedScope_2_Node_bool(this_, rootNode, isBottom)
	return this_
}

func InitTypedScope_2_Node_bool(this_ *TypedScope, rootNode pkg_rhino.Node, isBottom bool) {
	this_.initFields_TypedScope()
	this_.rootNode = rootNode
	this_.parent = nil
	this_.depth = 0
	this_.isBottom = isBottom
}

func TypedScope_CreateGlobalScope(rootNode pkg_rhino.Node) *TypedScope {
	return NewTypedScope_2_Node_bool(rootNode, false)
}

func TypedScope_CreateLatticeBottom(rootNode pkg_rhino.Node) *TypedScope {
	return NewTypedScope_2_Node_bool(rootNode, true)
}

func (this_ *TypedScope) Declare(name string, nameNode pkg_rhino.Node, type_ pkg_rhino.JSType, input *CompilerInput, inferred bool) *TypedVar {
	var var_ *TypedVar = NewTypedVar(inferred, name, nameNode, type_, this_, int32(len(this_.vars)), input)
	if this_.vars == nil {
		this_.vars = make(map[string]*TypedVar)
	}
	this_.vars[name] = var_
	return var_
}

func (this_ *TypedScope) MakeImplicitVar(name string) *TypedVar {
	return NewTypedVar(false, name, 0, nil, this_, -1, nil)
}

func (this_ *TypedScope) GetOwnVar(name string) *TypedVar {
	return this_.vars[name]
}

var _ = j2gort.NullString
`)

	mustWrite("gen/jscomp/typed_var.go", `package jscomp

import pkg_rhino "example.com/w33btest/gen/rhino"

type TypedVar struct {
	name         string
	nameNode     pkg_rhino.Node
	type_        pkg_rhino.JSType
	scope        *TypedScope
	index        int32
	input        *CompilerInput
	typeInferred bool
}

func (this_ *TypedVar) initFields_TypedVar() {}

func NewTypedVar(inferred bool, name string, nameNode pkg_rhino.Node, type_ pkg_rhino.JSType, scope *TypedScope, index int32, input *CompilerInput) *TypedVar {
	this_ := &TypedVar{}
	InitTypedVar(this_, inferred, name, nameNode, type_, scope, index, input)
	return this_
}

func InitTypedVar(this_ *TypedVar, inferred bool, name string, nameNode pkg_rhino.Node, type_ pkg_rhino.JSType, scope *TypedScope, index int32, input *CompilerInput) {
	this_.initFields_TypedVar()
	this_.name = name
	this_.nameNode = nameNode
	this_.type_ = type_
	this_.scope = scope
	this_.index = index
	this_.input = input
	this_.typeInferred = inferred
}

func (this_ *TypedVar) Name() string { return this_.name }
`)

	mustWrite("gen/jscomp/typed_scope_creator.go", `package jscomp

import (
	"example.com/w33btest/gen/j2gort"
	pkg_rhino "example.com/w33btest/gen/rhino"
)

type TypedScopeCreator struct {
	compiler AbstractCompiler
	memoized j2gort.OrderedMap[pkg_rhino.Node, *TypedScope]
}

func NewTypedScopeCreator_1(compiler AbstractCompiler) *TypedScopeCreator {
	this_ := &TypedScopeCreator{
		compiler: compiler,
		memoized: j2gort.NewOrderedMapCap[pkg_rhino.Node, *TypedScope](16),
	}
	return this_
}

func (this_ *TypedScopeCreator) CreateInitialScope(root pkg_rhino.Node) *TypedScope {
	var s *TypedScope = TypedScope_CreateGlobalScope(root)
	s.Declare("Object", 0, nil, nil, false)
	s.Declare("Array", 0, nil, nil, false)
	this_.memoized.Put(root, s)
	return s
}

func (this_ *TypedScopeCreator) CreateChildScope(root pkg_rhino.Node, parent *TypedScope) *TypedScope {
	newScope := NewTypedScope_4_TypedScope_Node_slice_string_Module(parent, root, nil, nil)
	newScope.Declare("localX", root, nil, nil, true)
	this_.memoized.Put(root, newScope)
	return newScope
}

type TypeInference struct {
	scopeCreator *TypedScopeCreator
	bottomScope  *TypedScope
}

func NewTypeInference(syntacticScope *TypedScope, scopeCreator *TypedScopeCreator) *TypeInference {
	this_ := &TypeInference{}
	InitTypeInference(this_, syntacticScope, scopeCreator)
	return this_
}

func InitTypeInference(this_ *TypeInference, syntacticScope *TypedScope, scopeCreator *TypedScopeCreator) {
	this_.scopeCreator = scopeCreator
	this_.bottomScope = TypedScope_CreateLatticeBottom(syntacticScope.rootNode)
}

func (this_ *TypeInference) ResetForScope(syntacticScope *TypedScope) {
	this_.bottomScope = TypedScope_CreateLatticeBottom(syntacticScope.rootNode)
}
`)

	// Variant 2: Direct (*DeadPropertyAssignmentElimination) methods calling NewDeadPropertyAssignmentElimination_Property + cleaned map.
	mustWrite("gen/jscomp_direct/dpae_direct.go", `package jscomp

import (
	"example.com/w33btest/gen/j2gort"
	pkg_rhino "example.com/w33btest/gen/rhino"
)

type DeadPropertyAssignmentElimination_Property struct {
	name string
}

func (p *DeadPropertyAssignmentElimination_Property) initFields_Prop() {}

func NewDeadPropertyAssignmentElimination_Property(name string) *DeadPropertyAssignmentElimination_Property {
	this_ := &DeadPropertyAssignmentElimination_Property{}
	InitDeadPropertyAssignmentElimination_Property(this_, name)
	return this_
}

func InitDeadPropertyAssignmentElimination_Property(this_ *DeadPropertyAssignmentElimination_Property, name string) {
	this_.initFields_Prop()
	this_.name = name
}

type DeadPropertyAssignmentElimination struct {
	properties j2gort.OrderedMap[string, *DeadPropertyAssignmentElimination_Property]
}

func (this_ *DeadPropertyAssignmentElimination) recordProp(name string) *DeadPropertyAssignmentElimination_Property {
	p := this_.properties.Get(name)
	if p == nil {
		p = NewDeadPropertyAssignmentElimination_Property(name)
		this_.properties.Put(name, p)
	}
	return p
}

func (this_ *DeadPropertyAssignmentElimination) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer this_.properties.Clear()
	_ = this_.recordProp("a")
	_ = this_.recordProp("b")
}
`)

	// Variant 3 (Negative safety guard): DeadPropertyAssignmentElimination retaining properties on struct WITHOUT clearing in Process_2!
	mustWrite("gen/jscomp_unsafe/dpae_unsafe.go", `package jscomp

import (
	"example.com/w33btest/gen/j2gort"
	pkg_rhino "example.com/w33btest/gen/rhino"
)

type DeadPropertyAssignmentElimination_Property struct {
	name string
}

func (p *DeadPropertyAssignmentElimination_Property) initFields_Prop() {}

func NewDeadPropertyAssignmentElimination_Property(name string) *DeadPropertyAssignmentElimination_Property {
	this_ := &DeadPropertyAssignmentElimination_Property{}
	InitDeadPropertyAssignmentElimination_Property(this_, name)
	return this_
}

func InitDeadPropertyAssignmentElimination_Property(this_ *DeadPropertyAssignmentElimination_Property, name string) {
	this_.initFields_Prop()
	this_.name = name
}

type DeadPropertyAssignmentElimination struct {
	leakedProps map[string]*DeadPropertyAssignmentElimination_Property
}

func (this_ *DeadPropertyAssignmentElimination) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.leakedProps == nil {
		this_.leakedProps = make(map[string]*DeadPropertyAssignmentElimination_Property)
	}
	this_.leakedProps["x"] = NewDeadPropertyAssignmentElimination_Property("x")
}

var _ = j2gort.NullString
`)

	allPasses := passes.AllPasses()
	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	directDir := filepath.Join(tmpDir, "gen", "jscomp_direct")
	unsafeDir := filepath.Join(tmpDir, "gen", "jscomp_unsafe")

	for _, tc := range []struct {
		name     string
		dir      string
		wantMin1 int
		wantMax1 int
	}{
		{"jscomp", jscompDir, 6, 100},
		{"jscomp_direct", directDir, 2, 10},
		{"jscomp_unsafe", unsafeDir, 0, 0},
	} {
		s1, err := passes.RunOnPackageDir(tc.dir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 1 on %s failed: %v", tc.name, err)
		}
		if s1.TotalSites < tc.wantMin1 || s1.TotalSites > tc.wantMax1 {
			t.Fatalf("Run 1 on %s: TotalSites = %d, want [%d, %d] (byPass=%+v)",
				tc.name, s1.TotalSites, tc.wantMin1, tc.wantMax1, s1.ByPass)
		}
		s2, err := passes.RunOnPackageDir(tc.dir, allPasses, passes.Options{})
		if err != nil {
			t.Fatalf("Run 2 (idempotency) on %s failed: %v", tc.name, err)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Run 2 on %s NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v",
				tc.name, s2.TotalSites, s2.FilesModified, s2.ByPass)
		}
	}

	readGot := func(rel string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	// Check Part A AST:
	phasesCode := readGot("gen/jscomp/phases.go")
	if strings.Contains(phasesCode, "TrimStaleSlabPools") {
		t.Fatalf("TrimStaleSlabPools was not stripped:\n%s", phasesCode)
	}
	if got := strings.Count(phasesCode, "PurgeAllRuntimePools()"); got != 3 {
		t.Fatalf("expected exactly 3 PurgeAllRuntimePools() calls in phases.go, got %d:\n%s", got, phasesCode)
	}
	if strings.Count(phasesCode, "defer j2gort.PurgeAllRuntimePools()") != 1 {
		t.Fatalf("expected exactly 1 defer j2gort.PurgeAllRuntimePools() (in PerformTranspilationAndOptimizations), got:\n%s", phasesCode)
	}

	normWS := func(s string) string { return strings.Join(strings.Fields(s), " ") }

	// Check Part B AST (both Shape 1 and Shape 2 + negative guard):
	dpaeCode := normWS(readGot("gen/jscomp/dead_property_assignment_elimination.go"))
	for _, want := range []string{
		"slab_DeadPropertyAssignmentElimination_Property j2gort.PooledSlab[DeadPropertyAssignmentElimination_Property]",
		"defer this_.slab_DeadPropertyAssignmentElimination_Property.Release()",
		"allocFunctionVisitorWithSlab_(&this_.slab_DeadPropertyAssignmentElimination_Property,",
		"val_1.propSlab_ = this_.propSlab_",
		"defer this_.propSlab_.Reset()",
		"allocOptPooled_DeadPropertyAssignmentElimination_Property(this_.propSlab_, propName)",
	} {
		if !strings.Contains(dpaeCode, want) {
			t.Fatalf("missing %q in lowered dead_property_assignment_elimination.go:\n%s", want, dpaeCode)
		}
	}

	directCode := normWS(readGot("gen/jscomp_direct/dpae_direct.go"))
	for _, want := range []string{
		"slab_DeadPropertyAssignmentElimination_Property j2gort.PooledSlab[DeadPropertyAssignmentElimination_Property]",
		"allocPooled_DeadPropertyAssignmentElimination_Property(&this_.slab_DeadPropertyAssignmentElimination_Property, name)",
		"defer this_.slab_DeadPropertyAssignmentElimination_Property.Release()",
	} {
		if !strings.Contains(directCode, want) {
			t.Fatalf("missing %q in lowered dpae_direct.go:\n%s", want, directCode)
		}
	}

	unsafeCode := readGot("gen/jscomp_unsafe/dpae_unsafe.go")
	if strings.Contains(unsafeCode, "slab_DeadPropertyAssignmentElimination_Property") {
		t.Fatalf("leaking DeadPropertyAssignmentElimination MUST NOT be slab-lowered:\n%s", unsafeCode)
	}

	// Check TypedScope / TypedVar are NOT arena-slab-lowered (prevents interior-pointer live-heap retention):
	tscCode := normWS(readGot("gen/jscomp/typed_scope_creator.go"))
	tsCode := normWS(readGot("gen/jscomp/typed_scope.go"))
	tvCode := normWS(readGot("gen/jscomp/typed_var.go"))
	for _, unwanted := range []string{
		"slab_TypedScope",
		"varSlab_",
		"allocTypedScope",
	} {
		if strings.Contains(tscCode, unwanted) || strings.Contains(tsCode, unwanted) {
			t.Fatalf("unexpected %q in lowered typed_scope_creator.go / typed_scope.go", unwanted)
		}
	}
	if strings.Contains(tvCode, "allocOpt_TypedVar") {
		t.Fatalf("unexpected allocOpt_TypedVar in lowered typed_var.go:\n%s", tvCode)
	}

	// End-to-end runtime + -race behavioral test!
	mustWrite("w33b_rt_test.go", `package w33b_rt_test

import (
	"testing"

	"example.com/w33btest/gen/j2gort"
	"example.com/w33btest/gen/jscomp"
)

func TestWave33BRuntimeSemantics(t *testing.T) {
	// Part A: Exact 3 purge sites!
	j2gort.ResetWave33BMetrics()
	comp := &jscomp.Compiler{}
	_ = comp.ParseInputs()
	if got := j2gort.PurgeCalls.Load(); got != 0 {
		t.Fatalf("ParseInputs PurgeCalls = %d, want 0", got)
	}
	pd := &jscomp.ProcessDefines{}
	pd.Process_2(0, 1)
	if got := j2gort.PurgeCalls.Load(); got != 0 {
		t.Fatalf("ProcessDefines PurgeCalls = %d, want 0", got)
	}
	tip := &jscomp.TypeInferencePass{}
	_ = tip.InferAllScopes(1)
	if got := j2gort.PurgeCalls.Load(); got != 1 {
		t.Fatalf("InferAllScopes PurgeCalls = %d, want 1 (entry only)", got)
	}
	tas := &jscomp.TypedAstSerializer{}
	tas.SerializeRoots(0, 1)
	if got := j2gort.PurgeCalls.Load(); got != 2 {
		t.Fatalf("after SerializeRoots PurgeCalls = %d, want 2", got)
	}
	comp.PerformTranspilationAndOptimizations(0)
	if got := j2gort.PurgeCalls.Load(); got != 3 {
		t.Fatalf("after PerformTranspilationAndOptimizations PurgeCalls = %d, want 3", got)
	}

	// Part B: DeadPropertyAssignmentElimination PooledSlab alloc, per-function Reset, and Process_2 Release!
	j2gort.ResetWave33BMetrics()
	jscomp.TotalUniquePropsObserved.Store(0)
	dpae := jscomp.NewDeadPropertyAssignmentElimination(comp)
	dpae.Process_2(0, 1)
	// Each of 2 functions creates "propA", "obj", "propB" = 3 unique properties * 2 functions = 6 total!
	if got := jscomp.TotalUniquePropsObserved.Load(); got != 6 {
		t.Fatalf("TotalUniquePropsObserved = %d, want 6", got)
	}
	if got := j2gort.PooledAllocCalls.Load(); got != 6 {
		t.Fatalf("PooledAllocCalls = %d, want 6", got)
	}
	if got := j2gort.PooledResetCalls.Load(); got != 2 {
		t.Fatalf("PooledResetCalls = %d, want 2 (1 per function scope)", got)
	}
	if got := j2gort.PooledReleaseCalls.Load(); got != 1 {
		t.Fatalf("PooledReleaseCalls = %d, want 1 (on Process_2 exit)", got)
	}
	// Standalone callers of NewDeadPropertyAssignmentElimination_Property & FunctionVisitor still work safely!
	standaloneProp := jscomp.NewDeadPropertyAssignmentElimination_Property("standalone")
	if standaloneProp.Name() != "standalone" {
		t.Fatalf("standalone Property.Name() = %q", standaloneProp.Name())
	}
	standaloneFV := jscomp.NewDeadPropertyAssignmentElimination_FunctionVisitor(nil)
	standaloneFV.EnterChangedScopeRoot(comp, 10)

	// TypedScope + TypedVar allocate individually on heap (0 Slab[TypedScope]/Slab[TypedVar] chunk pinning):
	j2gort.ResetWave33BMetrics()
	tsc := jscomp.NewTypedScopeCreator_1(comp)
	globalScope := tsc.CreateInitialScope(1)
	childScope := tsc.CreateChildScope(10, globalScope)
	_ = childScope.MakeImplicitVar("arguments")
	grandChild := jscomp.NewTypedScope_2_TypedScope_Node(childScope, 11)
	grandChild.Declare("gcVar", 11, nil, nil, false)

	ti := jscomp.NewTypeInference(childScope, tsc)
	ti.ResetForScope(childScope)

	if got := j2gort.SlabAllocCalls.Load(); got != 0 {
		t.Fatalf("TypedScope+TypedVar SlabAllocCalls = %d, want 0", got)
	}
	if childScope.GetOwnVar("localX") == nil || grandChild.GetOwnVar("gcVar") == nil {
		t.Fatal("missing declared vars in TypedScopes")
	}

	standaloneScope := jscomp.TypedScope_CreateGlobalScope(1)
	v := standaloneScope.Declare("ext", 0, nil, nil, false)
	if v == nil || v.Name() != "ext" {
		t.Fatalf("standalone TypedScope.Declare failed: %+v", v)
	}
}
`)

	cmd := exec.Command("go", "test", "-race", "-v", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race failed: %v\n%s", err, string(out))
	}
}

func TestWave33B_RealGenJscompCopyIdempotencyAndBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping full real gen/jscomp copy test in -short mode")
	}
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is mid-edit; skipping real gen/jscomp copy check")
	}
	checkCmd := exec.Command("go", "build", "./transpiled/gen/j2gort", "./transpiled/gen/rhino", "./transpiled/gen/diagnostic")
	checkCmd.Dir = repoRoot
	if err := checkCmd.Run(); err != nil {
		t.Skipf("gen/j2gort, gen/rhino, or gen/diagnostic is mid-edit (%v); skipping real gen/jscomp copy check", err)
	}

	realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
	if _, err := os.Stat(filepath.Join(realJscompDir, "typed_scope_creator.go")); err != nil {
		t.Skipf("real gen/jscomp not present: %v", err)
	}

	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw33b", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	copyDir := filepath.Join(tmpDir, "jscomp")
	if err := os.MkdirAll(copyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(realJscompDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(realJscompDir, ent.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(copyDir, ent.Name()), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	allPasses := passes.AllPasses()
	s1, err := passes.RunOnPackageDir(copyDir, allPasses, passes.Options{})
	if err != nil {
		if !realGenJ2gortReady(repoRoot) {
			t.Skipf("dependency became mid-edit during Run 1: %v", err)
		}
		t.Fatalf("Run 1 on copied gen/jscomp failed: %v", err)
	}
	t.Logf("Wave 33B Real gen/jscomp Run 1 sites: total=%d files=%d byPass=%+v", s1.TotalSites, s1.FilesModified, s1.ByPass)

	s2, err := passes.RunOnPackageDir(copyDir, allPasses, passes.Options{})
	if err != nil {
		if !realGenJ2gortReady(repoRoot) {
			t.Skipf("dependency became mid-edit during Run 2: %v", err)
		}
		t.Fatalf("Run 2 (idempotency) on copied gen/jscomp failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 on copied gen/jscomp NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v",
			s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	readCopy := func(file string) string {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(copyDir, file))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// Part A checks on real gen/jscomp:
	compCode := readCopy("compiler.go")
	if got := strings.Count(compCode, "PurgeAllRuntimePools()"); got != 1 {
		t.Fatalf("expected compiler.go to have 1 PurgeAllRuntimePools() (in PerformTranspilationAndOptimizations, removed from ParseInputs), got %d", got)
	}
	tipCode := readCopy("type_inference_pass.go")
	if !strings.Contains(tipCode, "j2gort.PurgeAllRuntimePools()") || strings.Contains(tipCode, "defer j2gort.PurgeAllRuntimePools()") {
		t.Fatalf("expected type_inference_pass.go to have entry PurgeAllRuntimePools() ONLY (no defer)")
	}
	tasCode := readCopy("typed_ast_serializer.go")
	if !strings.Contains(tasCode, "j2gort.PurgeAllRuntimePools()") || strings.Contains(tasCode, "defer j2gort.PurgeAllRuntimePools()") {
		t.Fatalf("expected typed_ast_serializer.go to have entry PurgeAllRuntimePools() ONLY")
	}
	pdCode := readCopy("process_defines.go")
	if strings.Contains(pdCode, "PurgeAllRuntimePools") {
		t.Fatalf("expected process_defines.go NOT to contain PurgeAllRuntimePools")
	}

	normWS := func(s string) string { return strings.Join(strings.Fields(s), " ") }

	// Part B checks on real gen/jscomp:
	dpaeCode := normWS(readCopy("dead_property_assignment_elimination.go"))
	for _, want := range []string{
		"slab_DeadPropertyAssignmentElimination_Property j2gort.PooledSlab[DeadPropertyAssignmentElimination_Property]",
		"defer this_.slab_DeadPropertyAssignmentElimination_Property.Release()",
		"allocOptPooled_DeadPropertyAssignmentElimination_Property(",
	} {
		if !strings.Contains(dpaeCode, want) {
			t.Fatalf("expected dead_property_assignment_elimination.go to contain %q", want)
		}
	}

	buildCmd := exec.Command("go", "build", "./jscomp")
	buildCmd.Dir = tmpDir
	if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
		if !realGenJ2gortReady(repoRoot) {
			t.Skipf("dependency became mid-edit during go build ./jscomp: %v", bErr)
		}
		t.Fatalf("go build ./jscomp failed on lowered copy: %v\n%s", bErr, string(out))
	}
}
