/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/FunctionRewriter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.FunctionRewriter
type FunctionRewriter struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *FunctionRewriter) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter", "FunctionRewriter"
}
func (this_ *FunctionRewriter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionRewriter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter", "com.google.javascript.jscomp.FunctionRewriter", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionRewriter](v_)
	})
}
func (*FunctionRewriter) Implements_jscomp_CompilerPass() {}

func StaticGet_FunctionRewriter_SAVINGS_THRESHOLD() int32 { return int32(16) }
func (this_ *FunctionRewriter) initFields_FunctionRewriter() {
	_ = this_
}

// Java constructor: FunctionRewriter(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewFunctionRewriter(compiler AbstractCompiler) *FunctionRewriter {
	this_ := &FunctionRewriter{}
	InitFunctionRewriter(this_, compiler)
	return this_
}

func InitFunctionRewriter(this_ *FunctionRewriter, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_FunctionRewriter()
	this_.compiler = compiler
}

// Java method: FunctionRewriter.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *FunctionRewriter) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var reducers []FunctionRewriter_Reducer = []FunctionRewriter_Reducer{NewFunctionRewriter_ReturnConstantReducer(), NewFunctionRewriter_GetterReducer(), NewFunctionRewriter_SetterReducer(), NewFunctionRewriter_EmptyFunctionReducer(), NewFunctionRewriter_IdentityReducer()}
	_ = reducers
	var reductionMap j2gort.Ext_com_google_common_collect_Multimap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_HashMultimap_Create_0())
	_ = reductionMap
	NodeTraversal_Traverse(this_.compiler, root, NewFunctionRewriter_ReductionGatherer(this_, j2gort.ImmListFromTypedSlice(reducers), reductionMap))
	{

		for _, reducer := range reducers {

			_ = reducer
			var reductions j2gort.Ext_java_util_Collection = j2gort.AsIface[j2gort.Ext_java_util_Collection](reductionMap.Get(reducer))
			_ = reductions
			if reductions.IsEmpty() {
				continue
			}
			var helperCode pkg_rhino.Node = this_.ParseHelperCode(reducer)
			_ = helperCode
			if helperCode == 0 {
				continue
			}
			var helperCodeCost int32 = InlineCostEstimator_GetCost_1(helperCode)
			_ = helperCodeCost
			var savings int32 = int32(0)
			_ = savings
			{
				cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(reductions))
				for cur_3.HasNext() {
					rawElem_4 := cur_3.Next()
					reduction := j2gort.CastPtr[FunctionRewriter_Reduction](rawElem_4)
					_ = reduction
					savings = (savings + reduction.EstimateSavings())
				}
			}
			if savings > (helperCodeCost + int32(16)) {
				{
					cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(reductions))
					for cur_5.HasNext() {
						rawElem_6 := cur_5.Next()
						reduction_2 := j2gort.CastPtr[FunctionRewriter_Reduction](rawElem_6)
						_ = reduction_2
						reduction_2.Apply()
					}
				}
				var addingRoot pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(nil)
				_ = addingRoot
				NodeUtil_MarkNewScopesChanged_2_Node_AbstractCompiler(helperCode, this_.compiler)
				addingRoot.AddChildToFront(helperCode)
				this_.compiler.ReportChangeToEnclosingScope(addingRoot)
			}
		}
	}
}

// Java method: FunctionRewriter.parseHelperCode(com.google.javascript.jscomp.FunctionRewriter.Reducer reducer)
func (this_ *FunctionRewriter) ParseHelperCode(reducer FunctionRewriter_Reducer) pkg_rhino.Node {
	_ = reducer
	var root pkg_rhino.Node = this_.compiler.ParseSyntheticCode((j2gort.StrVal(j2gort.GetObjectClass(reducer)) + ":helper"), reducer.GetHelperSource())
	_ = root
	if root != 0 {
		return root.RemoveFirstChild()
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java static method: FunctionRewriter.isReduceableFunctionExpression(com.google.javascript.rhino.Node n)
func FunctionRewriter_IsReduceableFunctionExpression(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (((((NodeUtil_IsFunctionExpression(n) && !(NodeUtil_IsGetOrSetKey(parent))) && !(parent.IsMemberFunctionDef())) && !(n.IsAsyncFunction())) && !(n.IsGeneratorFunction())) && !(NodeUtil_HasNonSimpleParameters(n)))
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.Reduction
type FunctionRewriter_Reduction struct {
	outer_   *FunctionRewriter
	oldChild pkg_rhino.Node `gson:"oldChild"`
	newChild pkg_rhino.Node `gson:"newChild"`
	idHash_  atomic.Int32
}

func (this_ *FunctionRewriter_Reduction) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_Reduction) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter_Reduction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$Reduction", "Reduction"
}
func (this_ *FunctionRewriter_Reduction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionRewriter_Reduction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.Reduction", "com.google.javascript.jscomp.FunctionRewriter$Reduction", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionRewriter_Reduction](v_)
	})
}
func (this_ *FunctionRewriter_Reduction) SetOuter_(o *FunctionRewriter) { this_.outer_ = o }
func (this_ *FunctionRewriter_Reduction) GetOuter_() *FunctionRewriter  { return this_.outer_ }

func (this_ *FunctionRewriter_Reduction) initFields_FunctionRewriter_Reduction() {
	_ = this_
}

// Java constructor: Reduction(com.google.javascript.rhino.Node oldChild, com.google.javascript.rhino.Node newChild)
func NewFunctionRewriter_Reduction(outer_ *FunctionRewriter, oldChild pkg_rhino.Node, newChild pkg_rhino.Node) *FunctionRewriter_Reduction {
	this_ := &FunctionRewriter_Reduction{}
	this_.outer_ = outer_
	InitFunctionRewriter_Reduction(this_, oldChild, newChild)
	return this_
}

func InitFunctionRewriter_Reduction(this_ *FunctionRewriter_Reduction, oldChild pkg_rhino.Node, newChild pkg_rhino.Node) {
	_ = oldChild
	_ = newChild
	this_.initFields_FunctionRewriter_Reduction()
	this_.oldChild = oldChild
	this_.newChild = newChild
}

// Java method: Reduction.apply()
func (this_ *FunctionRewriter_Reduction) Apply() {
	this_.oldChild.ReplaceWith(this_.newChild)
	NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(this_.oldChild, this_.outer_.compiler)
	this_.outer_.compiler.ReportChangeToEnclosingScope(this_.newChild)
}

// Java method: Reduction.estimateSavings()
func (this_ *FunctionRewriter_Reduction) EstimateSavings() int32 {
	return (InlineCostEstimator_GetCost_1(this_.oldChild) - InlineCostEstimator_GetCost_1(this_.newChild))
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.ReductionGatherer
type FunctionRewriter_ReductionGatherer struct {
	outer_     *FunctionRewriter
	reducers   []FunctionRewriter_Reducer                    `gson:"reducers"`
	reductions j2gort.Ext_com_google_common_collect_Multimap `gson:"reductions"`
	idHash_    atomic.Int32
}

func (this_ *FunctionRewriter_ReductionGatherer) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_ReductionGatherer) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter_ReductionGatherer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$ReductionGatherer", "ReductionGatherer"
}
func (this_ *FunctionRewriter_ReductionGatherer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionRewriter_ReductionGatherer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.ReductionGatherer", "com.google.javascript.jscomp.FunctionRewriter$ReductionGatherer", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionRewriter_ReductionGatherer](v_)
	})
}
func (this_ *FunctionRewriter_ReductionGatherer) SetOuter_(o *FunctionRewriter)       { this_.outer_ = o }
func (this_ *FunctionRewriter_ReductionGatherer) GetOuter_() *FunctionRewriter        { return this_.outer_ }
func (*FunctionRewriter_ReductionGatherer) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *FunctionRewriter_ReductionGatherer) initFields_FunctionRewriter_ReductionGatherer() {
	_ = this_
}

// Java constructor: ReductionGatherer(java.util.List<com.google.javascript.jscomp.FunctionRewriter.Reducer> reducers, com.google.common.collect.Multimap<com.google.javascript.jscomp.FunctionRewriter.Reducer,com.google.javascript.jscomp.FunctionRewriter.Reduction> reductions)
func NewFunctionRewriter_ReductionGatherer(outer_ *FunctionRewriter, reducers j2gort.Ext_java_util_List, reductions j2gort.Ext_com_google_common_collect_Multimap) *FunctionRewriter_ReductionGatherer {
	this_ := &FunctionRewriter_ReductionGatherer{}
	this_.outer_ = outer_
	InitFunctionRewriter_ReductionGatherer(this_, reducers, reductions)
	return this_
}

func InitFunctionRewriter_ReductionGatherer(this_ *FunctionRewriter_ReductionGatherer, reducers j2gort.Ext_java_util_List, reductions j2gort.Ext_com_google_common_collect_Multimap) {
	_ = reducers
	_ = reductions
	this_.initFields_FunctionRewriter_ReductionGatherer()
	this_.reducers = func() []FunctionRewriter_Reducer {
		arr_10 := reducers.ToArray_1_slice_any(func() []any {
			arr_7 := make([]FunctionRewriter_Reducer, int32(0))
			if arr_7 == nil {
				return nil
			}
			out_8 := make([]any, len(arr_7))
			for i_, el_9 := range arr_7 {
				out_8[i_] = el_9
			}
			return out_8
		}())
		if arr_10 == nil {
			return nil
		}
		out_11 := make([]FunctionRewriter_Reducer, len(arr_10))
		for i_, el_12 := range arr_10 {
			out_11[i_] = j2gort.AsIface[FunctionRewriter_Reducer](el_12)
		}
		return out_11
	}()
	this_.reductions = reductions
}

// Java method: ReductionGatherer.shouldTraverse(com.google.javascript.jscomp.NodeTraversal raversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *FunctionRewriter_ReductionGatherer) ShouldTraverse_3(raversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = raversal
	_ = node
	_ = parent
	for _, el_13 := range j2gort.CheckNotNullSlice(this_.reducers) {
		reducer := el_13
		_ = reducer
		var replacement pkg_rhino.Node = reducer.Reduce(node)
		_ = replacement
		if replacement != node {
			this_.reductions.Put(reducer, NewFunctionRewriter_Reduction(this_.outer_, node, replacement))
			return false
		}
	}
	return true
}

// Java method: ReductionGatherer.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *FunctionRewriter_ReductionGatherer) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
}

// Java polymorphic class interface: com.google.javascript.jscomp.FunctionRewriter.Reducer
type FunctionRewriter_Reducer interface {
	Implements_jscomp_FunctionRewriter_Reducer()
	IsJavaNil_() bool
	GetFunctionRewriter_ReducerBase() *FunctionRewriter_ReducerBase
	GetHelperSource() string
	Reduce(pkg_rhino.Node) pkg_rhino.Node
	BuildCallNode(string, pkg_rhino.Node, pkg_rhino.Node) pkg_rhino.Node
}

func ToIface_FunctionRewriter_Reducer[T interface {
	FunctionRewriter_Reducer
	comparable
}](v T) FunctionRewriter_Reducer {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.Reducer
type FunctionRewriter_ReducerBase struct {
	idHash_ atomic.Int32
}

func (this_ *FunctionRewriter_ReducerBase) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_ReducerBase) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter_ReducerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$Reducer", "Reducer"
}
func (this_ *FunctionRewriter_ReducerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.Reducer", "com.google.javascript.jscomp.FunctionRewriter$Reducer", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionRewriter_Reducer](v_)
	})
}
func (b *FunctionRewriter_ReducerBase) GetFunctionRewriter_ReducerBase() *FunctionRewriter_ReducerBase {
	return b
}
func (*FunctionRewriter_ReducerBase) Implements_jscomp_FunctionRewriter_Reducer() {}

func (this_ *FunctionRewriter_ReducerBase) initFields_FunctionRewriter_Reducer() {
	_ = this_
}

func InitFunctionRewriter_Reducer(this_ *FunctionRewriter_ReducerBase) {
	this_.initFields_FunctionRewriter_Reducer()
}

// Java method: Reducer.getHelperSource()
func (this_ *FunctionRewriter_ReducerBase) GetHelperSource() string {
	panic("abstract or native method invoked")
}

// Java method: Reducer.reduce(com.google.javascript.rhino.Node node)
func (this_ *FunctionRewriter_ReducerBase) Reduce(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	panic("abstract or native method invoked")
}

// Java method: Reducer.buildCallNode(java.lang.String methodName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node argumentNode, com.google.javascript.rhino.Node srcref)
func (this_ *FunctionRewriter_ReducerBase) BuildCallNode(methodName string, argumentNode pkg_rhino.Node, srcref pkg_rhino.Node) pkg_rhino.Node {
	_ = methodName
	_ = argumentNode
	_ = srcref
	var call pkg_rhino.Node = pkg_rhino.IR_Call(pkg_rhino.IR_Name(methodName), []pkg_rhino.Node{}).SrcrefTree(srcref)
	_ = call
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	if argumentNode != 0 {
		call.AddChildToBack(argumentNode.CloneTree_0())
	}
	return call
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.EmptyFunctionReducer
type FunctionRewriter_EmptyFunctionReducer struct {
	FunctionRewriter_ReducerBase
}

func (this_ *FunctionRewriter_EmptyFunctionReducer) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_EmptyFunctionReducer) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter_EmptyFunctionReducer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$EmptyFunctionReducer", "EmptyFunctionReducer"
}
func (this_ *FunctionRewriter_EmptyFunctionReducer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionRewriter_ReducerBase.IdentityHash_()
}
func (*FunctionRewriter_EmptyFunctionReducer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.EmptyFunctionReducer", "com.google.javascript.jscomp.FunctionRewriter$EmptyFunctionReducer", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionRewriter_EmptyFunctionReducer](v_)
	})
}
func (*FunctionRewriter_EmptyFunctionReducer) Implements_jscomp_FunctionRewriter_Reducer() {}

func StaticGet_FunctionRewriter_EmptyFunctionReducer_FACTORY_METHOD_NAME() string {
	return "JSCompiler_emptyFn"
}
func StaticGet_FunctionRewriter_EmptyFunctionReducer_HELPER_SOURCE() string {
	return "function JSCompiler_emptyFn() {  return function() {}}"
}
func (this_ *FunctionRewriter_EmptyFunctionReducer) initFields_FunctionRewriter_EmptyFunctionReducer() {
	_ = this_
}

// Java constructor: EmptyFunctionReducer()
func NewFunctionRewriter_EmptyFunctionReducer() *FunctionRewriter_EmptyFunctionReducer {
	this_ := &FunctionRewriter_EmptyFunctionReducer{}
	InitFunctionRewriter_EmptyFunctionReducer(this_)
	return this_
}

func InitFunctionRewriter_EmptyFunctionReducer(this_ *FunctionRewriter_EmptyFunctionReducer) {
	InitFunctionRewriter_Reducer(&this_.FunctionRewriter_ReducerBase)
	this_.initFields_FunctionRewriter_EmptyFunctionReducer()
}

// Java method: EmptyFunctionReducer.getHelperSource()
func (this_ *FunctionRewriter_EmptyFunctionReducer) GetHelperSource() string {
	return "function JSCompiler_emptyFn() {  return function() {}}"
}

// Java method: EmptyFunctionReducer.reduce(com.google.javascript.rhino.Node node)
func (this_ *FunctionRewriter_EmptyFunctionReducer) Reduce(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if FunctionRewriter_IsReduceableFunctionExpression(node) && NodeUtil_IsEmptyFunctionExpression(node) {
		return this_.BuildCallNode("JSCompiler_emptyFn", pkg_rhino.Node(0), node)
	} else {
		return node
	}
}

// Java polymorphic class interface: com.google.javascript.jscomp.FunctionRewriter.SingleReturnStatementReducer
type FunctionRewriter_SingleReturnStatementReducer interface {
	Implements_jscomp_FunctionRewriter_SingleReturnStatementReducer()
	IsJavaNil_() bool
	GetFunctionRewriter_SingleReturnStatementReducerBase() *FunctionRewriter_SingleReturnStatementReducerBase
	FunctionRewriter_Reducer
	MaybeGetSingleReturnRValue(pkg_rhino.Node) pkg_rhino.Node
}

func ToIface_FunctionRewriter_SingleReturnStatementReducer[T interface {
	FunctionRewriter_SingleReturnStatementReducer
	comparable
}](v T) FunctionRewriter_SingleReturnStatementReducer {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.SingleReturnStatementReducer
type FunctionRewriter_SingleReturnStatementReducerBase struct {
	FunctionRewriter_ReducerBase
}

func (this_ *FunctionRewriter_SingleReturnStatementReducerBase) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_SingleReturnStatementReducerBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*FunctionRewriter_SingleReturnStatementReducerBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$SingleReturnStatementReducer", "SingleReturnStatementReducer"
}
func (this_ *FunctionRewriter_SingleReturnStatementReducerBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionRewriter_ReducerBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.SingleReturnStatementReducer", "com.google.javascript.jscomp.FunctionRewriter$SingleReturnStatementReducer", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionRewriter_SingleReturnStatementReducer](v_)
	})
}
func (b *FunctionRewriter_SingleReturnStatementReducerBase) GetFunctionRewriter_SingleReturnStatementReducerBase() *FunctionRewriter_SingleReturnStatementReducerBase {
	return b
}
func (*FunctionRewriter_SingleReturnStatementReducerBase) Implements_jscomp_FunctionRewriter_SingleReturnStatementReducer() {
}
func (*FunctionRewriter_SingleReturnStatementReducerBase) Implements_jscomp_FunctionRewriter_Reducer() {
}

func (this_ *FunctionRewriter_SingleReturnStatementReducerBase) initFields_FunctionRewriter_SingleReturnStatementReducer() {
	_ = this_
}

func InitFunctionRewriter_SingleReturnStatementReducer(this_ *FunctionRewriter_SingleReturnStatementReducerBase) {
	InitFunctionRewriter_Reducer(&this_.FunctionRewriter_ReducerBase)
	this_.initFields_FunctionRewriter_SingleReturnStatementReducer()
}

// Java method: SingleReturnStatementReducer.maybeGetSingleReturnRValue(com.google.javascript.rhino.Node functionNode)
func (this_ *FunctionRewriter_SingleReturnStatementReducerBase) MaybeGetSingleReturnRValue(functionNode pkg_rhino.Node) pkg_rhino.Node {
	_ = functionNode
	var body pkg_rhino.Node = functionNode.GetLastChild()
	_ = body
	if !(body.HasOneChild()) {
		return pkg_rhino.Node(0)
	}
	var statement pkg_rhino.Node = body.GetFirstChild()
	_ = statement
	if statement.IsReturn() {
		return statement.GetFirstChild()
	}
	return pkg_rhino.Node(0)
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.IdentityReducer
type FunctionRewriter_IdentityReducer struct {
	FunctionRewriter_SingleReturnStatementReducerBase
}

func (this_ *FunctionRewriter_IdentityReducer) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_IdentityReducer) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter_IdentityReducer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$IdentityReducer", "IdentityReducer"
}
func (this_ *FunctionRewriter_IdentityReducer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionRewriter_SingleReturnStatementReducerBase.IdentityHash_()
}
func (*FunctionRewriter_IdentityReducer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.IdentityReducer", "com.google.javascript.jscomp.FunctionRewriter$IdentityReducer", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionRewriter_IdentityReducer](v_)
	})
}
func (*FunctionRewriter_IdentityReducer) Implements_jscomp_FunctionRewriter_SingleReturnStatementReducer() {
}
func (*FunctionRewriter_IdentityReducer) Implements_jscomp_FunctionRewriter_Reducer() {}

func StaticGet_FunctionRewriter_IdentityReducer_FACTORY_METHOD_NAME() string {
	return "JSCompiler_identityFn"
}
func StaticGet_FunctionRewriter_IdentityReducer_HELPER_SOURCE() string {
	return "function JSCompiler_identityFn() {  return function(JSCompiler_identityFn_value) {return JSCompiler_identityFn_value}}"
}
func (this_ *FunctionRewriter_IdentityReducer) initFields_FunctionRewriter_IdentityReducer() {
	_ = this_
}

// Java constructor: IdentityReducer()
func NewFunctionRewriter_IdentityReducer() *FunctionRewriter_IdentityReducer {
	this_ := &FunctionRewriter_IdentityReducer{}
	InitFunctionRewriter_IdentityReducer(this_)
	return this_
}

func InitFunctionRewriter_IdentityReducer(this_ *FunctionRewriter_IdentityReducer) {
	InitFunctionRewriter_SingleReturnStatementReducer(&this_.FunctionRewriter_SingleReturnStatementReducerBase)
	this_.initFields_FunctionRewriter_IdentityReducer()
}

// Java method: IdentityReducer.getHelperSource()
func (this_ *FunctionRewriter_IdentityReducer) GetHelperSource() string {
	return "function JSCompiler_identityFn() {  return function(JSCompiler_identityFn_value) {return JSCompiler_identityFn_value}}"
}

// Java method: IdentityReducer.reduce(com.google.javascript.rhino.Node node)
func (this_ *FunctionRewriter_IdentityReducer) Reduce(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if !(FunctionRewriter_IsReduceableFunctionExpression(node)) {
		return node
	}
	if this_.priv_IdentityReducer_isIdentityFunction_2949(node) {
		return this_.BuildCallNode("JSCompiler_identityFn", pkg_rhino.Node(0), node)
	} else {
		return node
	}
}

// Java method: IdentityReducer.isIdentityFunction(com.google.javascript.rhino.Node functionNode)
func (this_ *FunctionRewriter_IdentityReducer) priv_IdentityReducer_isIdentityFunction_2949(functionNode pkg_rhino.Node) bool {
	_ = functionNode
	var argList pkg_rhino.Node = functionNode.GetSecondChild()
	_ = argList
	var paramNode pkg_rhino.Node = argList.GetFirstChild()
	_ = paramNode
	if (paramNode == 0) || !(paramNode.IsName()) {
		return false
	}
	var value pkg_rhino.Node = this_.MaybeGetSingleReturnRValue(functionNode)
	_ = value
	return (((value != 0) && value.IsName()) && func() bool {
		sr_14 := value.GetString()
		se_15 := paramNode.GetString()
		return j2gort.CheckNotNullString(sr_14) == se_15
	}())
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.ReturnConstantReducer
type FunctionRewriter_ReturnConstantReducer struct {
	FunctionRewriter_SingleReturnStatementReducerBase
}

func (this_ *FunctionRewriter_ReturnConstantReducer) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_ReturnConstantReducer) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter_ReturnConstantReducer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$ReturnConstantReducer", "ReturnConstantReducer"
}
func (this_ *FunctionRewriter_ReturnConstantReducer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionRewriter_SingleReturnStatementReducerBase.IdentityHash_()
}
func (*FunctionRewriter_ReturnConstantReducer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.ReturnConstantReducer", "com.google.javascript.jscomp.FunctionRewriter$ReturnConstantReducer", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionRewriter_ReturnConstantReducer](v_)
	})
}
func (*FunctionRewriter_ReturnConstantReducer) Implements_jscomp_FunctionRewriter_SingleReturnStatementReducer() {
}
func (*FunctionRewriter_ReturnConstantReducer) Implements_jscomp_FunctionRewriter_Reducer() {}

func StaticGet_FunctionRewriter_ReturnConstantReducer_FACTORY_METHOD_NAME() string {
	return "JSCompiler_returnArg"
}
func StaticGet_FunctionRewriter_ReturnConstantReducer_HELPER_SOURCE() string {
	return "function JSCompiler_returnArg(JSCompiler_returnArg_value) {  return function() {return JSCompiler_returnArg_value}}"
}
func (this_ *FunctionRewriter_ReturnConstantReducer) initFields_FunctionRewriter_ReturnConstantReducer() {
	_ = this_
}

// Java constructor: ReturnConstantReducer()
func NewFunctionRewriter_ReturnConstantReducer() *FunctionRewriter_ReturnConstantReducer {
	this_ := &FunctionRewriter_ReturnConstantReducer{}
	InitFunctionRewriter_ReturnConstantReducer(this_)
	return this_
}

func InitFunctionRewriter_ReturnConstantReducer(this_ *FunctionRewriter_ReturnConstantReducer) {
	InitFunctionRewriter_SingleReturnStatementReducer(&this_.FunctionRewriter_SingleReturnStatementReducerBase)
	this_.initFields_FunctionRewriter_ReturnConstantReducer()
}

// Java method: ReturnConstantReducer.getHelperSource()
func (this_ *FunctionRewriter_ReturnConstantReducer) GetHelperSource() string {
	return "function JSCompiler_returnArg(JSCompiler_returnArg_value) {  return function() {return JSCompiler_returnArg_value}}"
}

// Java method: ReturnConstantReducer.reduce(com.google.javascript.rhino.Node node)
func (this_ *FunctionRewriter_ReturnConstantReducer) Reduce(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if !(FunctionRewriter_IsReduceableFunctionExpression(node)) {
		return node
	}
	var valueNode pkg_rhino.Node = this_.priv_ReturnConstantReducer_getValueNode_2950(node)
	_ = valueNode
	if valueNode != 0 {
		return this_.BuildCallNode("JSCompiler_returnArg", valueNode, node)
	} else {
		return node
	}
}

// Java method: ReturnConstantReducer.getValueNode(com.google.javascript.rhino.Node functionNode)
func (this_ *FunctionRewriter_ReturnConstantReducer) priv_ReturnConstantReducer_getValueNode_2950(functionNode pkg_rhino.Node) pkg_rhino.Node {
	_ = functionNode
	var value pkg_rhino.Node = this_.MaybeGetSingleReturnRValue(functionNode)
	_ = value
	if (value != 0) && NodeUtil_IsImmutableValue(value) {
		return value
	}
	return pkg_rhino.Node(0)
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.GetterReducer
type FunctionRewriter_GetterReducer struct {
	FunctionRewriter_SingleReturnStatementReducerBase
}

func (this_ *FunctionRewriter_GetterReducer) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_GetterReducer) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter_GetterReducer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$GetterReducer", "GetterReducer"
}
func (this_ *FunctionRewriter_GetterReducer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionRewriter_SingleReturnStatementReducerBase.IdentityHash_()
}
func (*FunctionRewriter_GetterReducer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.GetterReducer", "com.google.javascript.jscomp.FunctionRewriter$GetterReducer", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionRewriter_GetterReducer](v_)
	})
}
func (*FunctionRewriter_GetterReducer) Implements_jscomp_FunctionRewriter_SingleReturnStatementReducer() {
}
func (*FunctionRewriter_GetterReducer) Implements_jscomp_FunctionRewriter_Reducer() {}

func StaticGet_FunctionRewriter_GetterReducer_FACTORY_METHOD_NAME() string { return "JSCompiler_get" }
func StaticGet_FunctionRewriter_GetterReducer_HELPER_SOURCE() string {
	return "function JSCompiler_get(JSCompiler_get_name) {  return function() {return this[JSCompiler_get_name]}}"
}
func (this_ *FunctionRewriter_GetterReducer) initFields_FunctionRewriter_GetterReducer() {
	_ = this_
}

// Java constructor: GetterReducer()
func NewFunctionRewriter_GetterReducer() *FunctionRewriter_GetterReducer {
	this_ := &FunctionRewriter_GetterReducer{}
	InitFunctionRewriter_GetterReducer(this_)
	return this_
}

func InitFunctionRewriter_GetterReducer(this_ *FunctionRewriter_GetterReducer) {
	InitFunctionRewriter_SingleReturnStatementReducer(&this_.FunctionRewriter_SingleReturnStatementReducerBase)
	this_.initFields_FunctionRewriter_GetterReducer()
}

// Java method: GetterReducer.getHelperSource()
func (this_ *FunctionRewriter_GetterReducer) GetHelperSource() string {
	return "function JSCompiler_get(JSCompiler_get_name) {  return function() {return this[JSCompiler_get_name]}}"
}

// Java method: GetterReducer.reduce(com.google.javascript.rhino.Node node)
func (this_ *FunctionRewriter_GetterReducer) Reduce(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if !(FunctionRewriter_IsReduceableFunctionExpression(node)) || node.IsArrowFunction() {
		return node
	}
	var getprop pkg_rhino.Node = this_.priv_GetterReducer_getGetPropertyName_2951(node)
	_ = getprop
	if getprop != 0 {
		if !(getprop.IsGetProp()) {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Expected GETPROP, got "+j2gort.StrVal(getprop.GetToken()))))
		}
		var arg pkg_rhino.Node = pkg_rhino.IR_String(getprop.GetString()).Srcref(getprop)
		_ = arg
		return this_.BuildCallNode("JSCompiler_get", arg, node)
	} else {
		return node
	}
}

// Java method: GetterReducer.getGetPropertyName(com.google.javascript.rhino.Node functionNode)
func (this_ *FunctionRewriter_GetterReducer) priv_GetterReducer_getGetPropertyName_2951(functionNode pkg_rhino.Node) pkg_rhino.Node {
	_ = functionNode
	var value pkg_rhino.Node = this_.MaybeGetSingleReturnRValue(functionNode)
	_ = value
	if ((value != 0) && value.IsGetProp()) && value.GetFirstChild().IsThis() {
		return value
	}
	return pkg_rhino.Node(0)
}

// Java class: com.google.javascript.jscomp.FunctionRewriter.SetterReducer
type FunctionRewriter_SetterReducer struct {
	FunctionRewriter_ReducerBase
}

func (this_ *FunctionRewriter_SetterReducer) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionRewriter_SetterReducer) IsJavaNil_() bool { return this_ == nil }
func (*FunctionRewriter_SetterReducer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FunctionRewriter$SetterReducer", "SetterReducer"
}
func (this_ *FunctionRewriter_SetterReducer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.FunctionRewriter_ReducerBase.IdentityHash_()
}
func (*FunctionRewriter_SetterReducer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FunctionRewriter.SetterReducer", "com.google.javascript.jscomp.FunctionRewriter$SetterReducer", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionRewriter_SetterReducer](v_)
	})
}
func (*FunctionRewriter_SetterReducer) Implements_jscomp_FunctionRewriter_Reducer() {}

func StaticGet_FunctionRewriter_SetterReducer_FACTORY_METHOD_NAME() string { return "JSCompiler_set" }
func StaticGet_FunctionRewriter_SetterReducer_HELPER_SOURCE() string {
	return "function JSCompiler_set(JSCompiler_set_name) {  return function(JSCompiler_set_value) {this[JSCompiler_set_name] = JSCompiler_set_value}}"
}
func (this_ *FunctionRewriter_SetterReducer) initFields_FunctionRewriter_SetterReducer() {
	_ = this_
}

// Java constructor: SetterReducer()
func NewFunctionRewriter_SetterReducer() *FunctionRewriter_SetterReducer {
	this_ := &FunctionRewriter_SetterReducer{}
	InitFunctionRewriter_SetterReducer(this_)
	return this_
}

func InitFunctionRewriter_SetterReducer(this_ *FunctionRewriter_SetterReducer) {
	InitFunctionRewriter_Reducer(&this_.FunctionRewriter_ReducerBase)
	this_.initFields_FunctionRewriter_SetterReducer()
}

// Java method: SetterReducer.getHelperSource()
func (this_ *FunctionRewriter_SetterReducer) GetHelperSource() string {
	return "function JSCompiler_set(JSCompiler_set_name) {  return function(JSCompiler_set_value) {this[JSCompiler_set_name] = JSCompiler_set_value}}"
}

// Java method: SetterReducer.reduce(com.google.javascript.rhino.Node node)
func (this_ *FunctionRewriter_SetterReducer) Reduce(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	if !(FunctionRewriter_IsReduceableFunctionExpression(node)) || node.IsArrowFunction() {
		return node
	}
	var propName pkg_rhino.Node = this_.priv_SetterReducer_getSetPropertyName_2952(node)
	_ = propName
	if propName != 0 {
		if !(propName.IsStringLit()) {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Expected STRING, got "+j2gort.StrVal(propName.GetToken()))))
		}
		return this_.BuildCallNode("JSCompiler_set", propName, node)
	} else {
		return node
	}
}

// Java method: SetterReducer.getSetPropertyName(com.google.javascript.rhino.Node functionNode)
func (this_ *FunctionRewriter_SetterReducer) priv_SetterReducer_getSetPropertyName_2952(functionNode pkg_rhino.Node) pkg_rhino.Node {
	_ = functionNode
	var body pkg_rhino.Node = functionNode.GetLastChild()
	_ = body
	if !(body.HasOneChild()) {
		return pkg_rhino.Node(0)
	}
	var argList pkg_rhino.Node = functionNode.GetSecondChild()
	_ = argList
	var paramNode pkg_rhino.Node = argList.GetFirstChild()
	_ = paramNode
	if (paramNode == 0) || !(paramNode.IsName()) {
		return pkg_rhino.Node(0)
	}
	var statement pkg_rhino.Node = body.GetFirstChild()
	_ = statement
	if !(NodeUtil_IsExprAssign(statement)) {
		return pkg_rhino.Node(0)
	}
	var assign pkg_rhino.Node = statement.GetFirstChild()
	_ = assign
	var lhs pkg_rhino.Node = assign.GetFirstChild()
	_ = lhs
	if lhs.IsGetProp() && lhs.GetFirstChild().IsThis() {
		var rhs pkg_rhino.Node = assign.GetLastChild()
		_ = rhs
		if rhs.IsName() && func() bool {
			sr_16 := rhs.GetString()
			se_17 := paramNode.GetString()
			return j2gort.CheckNotNullString(sr_16) == se_17
		}() {
			return pkg_rhino.IR_String(lhs.GetString()).ClonePropsFrom(lhs).Srcref(lhs)
		}
	}
	return pkg_rhino.Node(0)
}
func (*FunctionRewriter_ReductionGatherer) NoOpVisit_()                 {}
func (*FunctionRewriter_ReductionGatherer) CanReleaseSyntacticScopes_() {}
