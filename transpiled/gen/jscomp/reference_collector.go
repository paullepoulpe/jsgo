/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ReferenceCollector.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ReferenceCollector
type ReferenceCollector struct {
	referenceMap                            j2gort.OrderedMap[*Var, *ReferenceCollection] `gson:"referenceMap"`
	blockStack                              *j2gort.Deque[*BasicBlock]                    `gson:"blockStack"`
	behavior                                ReferenceCollector_Behavior                   `gson:"behavior"`
	scopeCreator                            ScopeCreator                                  `gson:"scopeCreator"`
	compiler                                AbstractCompiler                              `gson:"compiler"`
	varFilter                               j2gort.Ext_com_google_common_base_Predicate   `gson:"varFilter"`
	callback                                *ReferenceCollector_CollectorCallback         `gson:"callback"`
	collectedHoistedFunctions               j2gort.OrderedSet[pkg_rhino.Node]             `gson:"collectedHoistedFunctions"`
	narrowScope                             *Scope                                        `gson:"narrowScope"`
	idHash_                                 atomic.Int32
	collectHoistedFunctionsInReferenceOrder bool `gson:"collectHoistedFunctionsInReferenceOrder"`
	slab_BasicBlock                         j2gort.PooledSlab[BasicBlock]
	slab_ReferenceCollection                j2gort.PooledSlab[ReferenceCollection]
	slab_Reference                          j2gort.PooledSlab[Reference]
	ownsScopeCreator_                       bool
}

func (this_ *ReferenceCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *ReferenceCollector) IsJavaNil_() bool { return this_ == nil }
func (*ReferenceCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReferenceCollector", "ReferenceCollector"
}
func (this_ *ReferenceCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReferenceCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReferenceCollector", "com.google.javascript.jscomp.ReferenceCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*ReferenceCollector](v_)
	})
}
func (*ReferenceCollector) Implements_jscomp_CompilerPass()     {}
func (*ReferenceCollector) Implements_jscomp_ReferenceMap()     {}
func (*ReferenceCollector) Implements_rhino_StaticSymbolTable() {}

var static_ReferenceCollector_DO_NOTHING_BEHAVIOR ReferenceCollector_Behavior

func StaticGet_ReferenceCollector_DO_NOTHING_BEHAVIOR() ReferenceCollector_Behavior {
	clinit_ReferenceCollector()
	return static_ReferenceCollector_DO_NOTHING_BEHAVIOR
}
func StaticSet_ReferenceCollector_DO_NOTHING_BEHAVIOR(v ReferenceCollector_Behavior) {
	clinit_ReferenceCollector()
	static_ReferenceCollector_DO_NOTHING_BEHAVIOR = v
}

var clinitDone_ReferenceCollector uint32
var clinitState_ReferenceCollector j2gort.ClinitState

func clinit_ReferenceCollector() {
	if atomic.LoadUint32(&clinitDone_ReferenceCollector) != 0 {
		return
	}
	clinitSlow_ReferenceCollector()
}

func clinitSlow_ReferenceCollector() {
	if !j2gort.EnterClinit(&clinitState_ReferenceCollector, &clinitDone_ReferenceCollector) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReferenceCollector, &clinitDone_ReferenceCollector)
	static_ReferenceCollector_DO_NOTHING_BEHAVIOR = NewReferenceCollector__Anon_76()
}

func (this_ *ReferenceCollector) initFields_ReferenceCollector() {
	_ = this_
	this_.referenceMap = j2gort.AcquireOrderedMap[*Var, *ReferenceCollection]()
	this_.blockStack = &j2gort.Deque[*BasicBlock]{}
	this_.callback = NewReferenceCollector_CollectorCallback(this_)
	this_.collectedHoistedFunctions = j2gort.NewOrderedSet[pkg_rhino.Node]()
}

// Java constructor: ReferenceCollector(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ReferenceCollector.Behavior behavior, com.google.javascript.jscomp.ScopeCreator creator)
func NewReferenceCollector_3(compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) *ReferenceCollector {
	clinit_ReferenceCollector()
	this_ := &ReferenceCollector{}
	InitReferenceCollector_3(this_, compiler, behavior, creator)
	return this_
}

func InitReferenceCollector_3(this_ *ReferenceCollector, compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator) {
	_ = compiler
	_ = behavior
	_ = creator
	InitReferenceCollector_4(this_, compiler, behavior, creator, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())
}

// Java constructor: ReferenceCollector(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ReferenceCollector.Behavior behavior, com.google.javascript.jscomp.ScopeCreator creator, com.google.common.base.Predicate<com.google.javascript.jscomp.Var> varFilter)
func NewReferenceCollector_4(compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator, varFilter j2gort.Ext_com_google_common_base_Predicate) *ReferenceCollector {
	clinit_ReferenceCollector()
	this_ := &ReferenceCollector{}
	InitReferenceCollector_4(this_, compiler, behavior, creator, varFilter)
	return this_
}

func InitReferenceCollector_4(this_ *ReferenceCollector, compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator, varFilter j2gort.Ext_com_google_common_base_Predicate) {
	_ = compiler
	_ = behavior
	_ = creator
	_ = varFilter
	InitReferenceCollector_5(this_, compiler, behavior, creator, varFilter, true)
}

// Java constructor: ReferenceCollector(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ReferenceCollector.Behavior behavior, com.google.javascript.jscomp.ScopeCreator creator, com.google.common.base.Predicate<com.google.javascript.jscomp.Var> varFilter, boolean collectHoistedFunctionsInReferenceOrder)
func NewReferenceCollector_5(compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator, varFilter j2gort.Ext_com_google_common_base_Predicate, collectHoistedFunctionsInReferenceOrder bool) *ReferenceCollector {
	clinit_ReferenceCollector()
	this_ := &ReferenceCollector{}
	InitReferenceCollector_5(this_, compiler, behavior, creator, varFilter, collectHoistedFunctionsInReferenceOrder)
	return this_
}

func InitReferenceCollector_5(this_ *ReferenceCollector, compiler AbstractCompiler, behavior ReferenceCollector_Behavior, creator ScopeCreator, varFilter j2gort.Ext_com_google_common_base_Predicate, collectHoistedFunctionsInReferenceOrder bool) {
	_ = compiler
	_ = behavior
	_ = creator
	_ = varFilter
	_ = collectHoistedFunctionsInReferenceOrder
	this_.initFields_ReferenceCollector()
	this_.compiler = compiler
	this_.behavior = behavior
	this_.scopeCreator = creator
	this_.varFilter = varFilter
	this_.collectHoistedFunctionsInReferenceOrder = collectHoistedFunctionsInReferenceOrder
}

// Java method: ReferenceCollector.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ReferenceCollector) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.referenceMap.Size() == 0 {
		this_.referenceMap.ReleaseToPool()
		this_.referenceMap = j2gort.AcquireOrderedMapCap[*Var, *ReferenceCollection](int32(1024))
	}

	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.
			compiler)

		(&ntb_).SetCallback_1_NodeTraversal_Callback(ToIface_NodeTraversal_Callback(this_.callback))

		(&ntb_).SetScopeCreator(this_.
			scopeCreator)

		(&ntb_).SetObeyDestructuringAndDefaultValueExecutionOrder(true)

		(&ntb_).TraverseRoots(externs, root)
	}
}

// Java method: ReferenceCollector.process(com.google.javascript.rhino.Node root)
func (this_ *ReferenceCollector) Process_1(root pkg_rhino.Node) {
	_ = root
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.
			compiler)

		(&ntb_).SetCallback_1_NodeTraversal_Callback(ToIface_NodeTraversal_Callback(this_.callback))

		(&ntb_).SetScopeCreator(this_.
			scopeCreator)

		(&ntb_).SetObeyDestructuringAndDefaultValueExecutionOrder(true)

		(&ntb_).Traverse(root)
	}
}

// Java method: ReferenceCollector.processScope(com.google.javascript.jscomp.Scope scope)
func (this_ *ReferenceCollector) ProcessScope(scope *Scope) {
	_ = scope
	var shouldAddToBlockStack bool = !(scope.IsHoistScope())
	_ = shouldAddToBlockStack
	this_.narrowScope = scope
	if shouldAddToBlockStack {
		this_.priv_ReferenceCollector_pushNewBlock_4695(scope.GetRootNode_0())
	}
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.
			compiler)

		(&ntb_).SetCallback_1_NodeTraversal_Callback(ToIface_NodeTraversal_Callback(this_.callback))

		(&ntb_).SetScopeCreator(this_.
			scopeCreator)

		(&ntb_).SetObeyDestructuringAndDefaultValueExecutionOrder(true)

		(&ntb_).TraverseAtScope(ToIface_AbstractScope(scope))
	}
	if shouldAddToBlockStack {
		this_.priv_ReferenceCollector_popLastBlock_4696(scope.GetRootNode_0())
	}
	this_.narrowScope = nil
}

// Java method: ReferenceCollector.createTraversalBuilder()
func (this_ *ReferenceCollector) priv_ReferenceCollector_createTraversalBuilder_4693() *NodeTraversal_Builder {
	return NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(ToIface_NodeTraversal_Callback(this_.callback)).SetScopeCreator(this_.scopeCreator).SetObeyDestructuringAndDefaultValueExecutionOrder(true)
}

// Java method: ReferenceCollector.getAllSymbols()
func (this_ *ReferenceCollector) GetAllSymbols() any {
	return j2gort.AsIface[j2gort.Ext_java_util_Set](this_.referenceMap.LiveKeySet())
}

// Java method: ReferenceCollector.getScope(com.google.javascript.jscomp.Var var)
func (this_ *ReferenceCollector) GetScope(arg_1 any) any {
	var_ := j2gort.CastPtr[Var](arg_1)
	_ = var_
	return j2gort.PtrToAny(j2gort.CastPtr[Scope](var_.GetScope()))
}

// Java method: ReferenceCollector.getReferences(com.google.javascript.jscomp.Var v)
func (this_ *ReferenceCollector) GetReferences(arg_2 any) any {
	v := j2gort.CastPtr[Var](arg_2)
	_ = v
	rc_ := this_.referenceMap.Get(v)
	if rc_ == nil && v != nil && v.index < 0 {
		if sc_, ok_ := v.scope.(*Scope); ok_ && sc_ != nil && sc_.implicitVars != nil {
			if cv_, okV_ := sc_.GetOwnSlot(v.name).(*Var); okV_ && cv_ != nil && cv_ != v {
				rc_ = this_.referenceMap.Get(cv_)
			}
		}
	}
	return j2gort.PtrToAny(rc_)

}

// Java method: ReferenceCollector.maybeJumpToHoistedFunction(com.google.javascript.jscomp.Var var, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ReferenceCollector) priv_ReferenceCollector_maybeJumpToHoistedFunction_4694(var_ *Var, t *NodeTraversal) {
	_ = var_
	_ = t
	var fnNode pkg_rhino.Node = var_.GetParentNode()
	_ = fnNode
	var varScope *Scope = j2gort.CastPtr[Scope](var_.GetScope())
	_ = varScope
	if (((fnNode == 0) || !(NodeUtil_IsHoistedFunctionDeclaration(fnNode))) || ((this_.narrowScope != nil) && (this_.narrowScope.GetDepth() > varScope.GetDepth()))) || this_.collectedHoistedFunctions.Contains(fnNode) {
		return
	}
	var oldBlockStack *j2gort.Deque[*BasicBlock] = this_.blockStack
	_ = oldBlockStack
	this_.blockStack = &j2gort.Deque[*BasicBlock]{}
	if varScope.IsGlobal() {
		this_.blockStack.AddLast(oldBlockStack.GetFirst())
	} else {
		{
			cur_3 := oldBlockStack.Iter()
			for cur_3.HasNext() {

				b := cur_3.Next()
				_ = b
				this_.blockStack.AddLast(b)
				if b.GetRoot() == varScope.GetRootNode_0() {
					break
				}
			}
		}
	}
	this_.priv_ReferenceCollector_addReference_4697(var_, fnNode.GetFirstChild(), t)
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.
			compiler)

		(&ntb_).SetCallback_1_NodeTraversal_Callback(ToIface_NodeTraversal_Callback(this_.callback))

		(&ntb_).SetScopeCreator(this_.
			scopeCreator)

		(&ntb_).SetObeyDestructuringAndDefaultValueExecutionOrder(true)

		(&ntb_).TraverseAtScope(j2gort.AsIface[AbstractScope](this_.scopeCreator.CreateScope_2_Node_AbstractScope(fnNode, ToIface_AbstractScope(varScope))))
	}
	this_.blockStack = oldBlockStack
}

// Java static method: ReferenceCollector.isNonLocalModuleSpecifierName(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func ReferenceCollector_IsNonLocalModuleSpecifierName(n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	clinit_ReferenceCollector()
	_ = n
	_ = parent
	return ((parent.IsImportSpec() && !(n == parent.GetLastChild())) || (parent.IsExportSpec() && !(n == parent.GetFirstChild())))
}

// Java method: ReferenceCollector.pushNewBlock(com.google.javascript.rhino.Node root)
func (this_ *ReferenceCollector) priv_ReferenceCollector_pushNewBlock_4695(root pkg_rhino.Node) {
	_ = root
	this_.blockStack.AddLast(allocPooled_BasicBlock(&this_.slab_BasicBlock, this_.blockStack.PeekLast(), root))
}

// Java method: ReferenceCollector.popLastBlock(com.google.javascript.rhino.Node root)
func (this_ *ReferenceCollector) priv_ReferenceCollector_popLastBlock_4696(root pkg_rhino.Node) {
	_ = root
	var last *BasicBlock = this_.blockStack.RemoveLast()
	_ = last
	if !(root == last.GetRoot()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(root))
	}
}

// Java static method: ReferenceCollector.isBlockBoundary(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func ReferenceCollector_IsBlockBoundary(n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	clinit_ReferenceCollector()
	_ = n
	_ = parent
	if parent != 0 {
		switch parent.TokenSwitchOrd_() {
		case int32(79) /* Token.DO */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(48) /* Token.TRY */, int32(78) /* Token.WHILE */, int32(85) /* Token.WITH */, int32(105) /* Token.CLASS */, int32(91) /* Token.SWITCH_BODY */ :
			return true
		case int32(69) /* Token.AND */, int32(67) /* Token.HOOK */, int32(74) /* Token.IF */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(31) /* Token.OPTCHAIN_CALL */, int32(134) /* Token.DEFAULT_VALUE */ :
			return (n != parent.GetFirstChild())
		default:
		}
	}
	return n.IsCase()
}

// Java method: ReferenceCollector.addReference(com.google.javascript.jscomp.Var v, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ReferenceCollector) priv_ReferenceCollector_addReference_4697(v *Var, n pkg_rhino.Node, t *NodeTraversal) *Reference {
	_ = v
	_ = n
	_ = t
	if !(j2gort.RefEquals(this_.varFilter, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue())) && !(this_.varFilter.Apply(j2gort.PtrToAny(v))) {
		return nil
	}
	var collection *ReferenceCollection = this_.referenceMap.Get(v)
	_ = collection
	if collection == nil {
		collection = allocPooled_ReferenceCollection(&this_.slab_ReferenceCollection)
		this_.referenceMap.Put(v, collection)
	}
	val_5 := this_.slab_Reference.Alloc()

	InitReference_3(val_5, n, t, this_.blockStack.GetLast())
	var reference *Reference = val_5
	_ = reference
	collection.Add(reference)
	return reference
}

// Java class: com.google.javascript.jscomp.ReferenceCollector.CollectorCallback
type ReferenceCollector_CollectorCallback struct {
	outer_  *ReferenceCollector
	idHash_ atomic.Int32
}

func (this_ *ReferenceCollector_CollectorCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *ReferenceCollector_CollectorCallback) IsJavaNil_() bool { return this_ == nil }
func (*ReferenceCollector_CollectorCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReferenceCollector$CollectorCallback", "CollectorCallback"
}
func (this_ *ReferenceCollector_CollectorCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReferenceCollector_CollectorCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReferenceCollector.CollectorCallback", "com.google.javascript.jscomp.ReferenceCollector$CollectorCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*ReferenceCollector_CollectorCallback](v_)
	})
}
func (this_ *ReferenceCollector_CollectorCallback) SetOuter_(o *ReferenceCollector) { this_.outer_ = o }
func (this_ *ReferenceCollector_CollectorCallback) GetOuter_() *ReferenceCollector {
	return this_.outer_
}
func (*ReferenceCollector_CollectorCallback) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*ReferenceCollector_CollectorCallback) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *ReferenceCollector_CollectorCallback) initFields_ReferenceCollector_CollectorCallback() {
	_ = this_
}

// Java constructor: CollectorCallback()
func NewReferenceCollector_CollectorCallback(outer_ *ReferenceCollector) *ReferenceCollector_CollectorCallback {
	this_ := &ReferenceCollector_CollectorCallback{}
	this_.outer_ = outer_
	InitReferenceCollector_CollectorCallback(this_)
	return this_
}

func InitReferenceCollector_CollectorCallback(this_ *ReferenceCollector_CollectorCallback) {
	this_.initFields_ReferenceCollector_CollectorCallback()
}

// Java method: CollectorCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReferenceCollector_CollectorCallback) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if this_.outer_.collectHoistedFunctionsInReferenceOrder && NodeUtil_IsHoistedFunctionDeclaration(n) {
		if !(this_.outer_.collectedHoistedFunctions.Add(n)) {
			return false
		}
	}
	this_.outer_.behavior.BeforeTraverseNode(nodeTraversal, n, parent)
	if ReferenceCollector_IsBlockBoundary(n, parent) {
		this_.outer_.priv_ReferenceCollector_pushNewBlock_4695(n)
	}
	return true
}

// Java method: CollectorCallback.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ReferenceCollector_CollectorCallback) EnterScope(t *NodeTraversal) {
	_ = t
	if t.IsHoistScope() {
		this_.outer_.priv_ReferenceCollector_pushNewBlock_4695(t.GetScopeRoot())
	}
}

// Java method: CollectorCallback.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *ReferenceCollector_CollectorCallback) ExitScope(t *NodeTraversal) {
	_ = t
	if t.IsHoistScope() {
		this_.outer_.priv_ReferenceCollector_popLastBlock_4696(t.GetScopeRoot())
	}
	this_.outer_.behavior.AfterExitScope(t, ToIface_ReferenceMap(this_.outer_))
}

// Java method: CollectorCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReferenceCollector_CollectorCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var addedReference *Reference = nil
	_ = addedReference
	if (n.IsName() || n.IsImportStar()) && !(ReferenceCollector_IsNonLocalModuleSpecifierName(n, parent)) {
		var v *Var = t.GetScope().LookupVar(n.GetString())
		_ = v
		if v != nil {
			addedReference = this_.outer_.priv_ReferenceCollector_addReference_4697(v, n, t)
			if this_.outer_.collectHoistedFunctionsInReferenceOrder {
				this_.outer_.priv_ReferenceCollector_maybeJumpToHoistedFunction_4694(v, t)
			}
		}
	}
	this_.outer_.behavior.AfterVisitNode(t, n, parent, addedReference)
	if ReferenceCollector_IsBlockBoundary(n, parent) {
		this_.outer_.priv_ReferenceCollector_popLastBlock_4696(n)
	}
}

// Java interface: com.google.javascript.jscomp.ReferenceCollector.Behavior
type ReferenceCollector_Behavior interface {
	Implements_jscomp_ReferenceCollector_Behavior()
	IsJavaNil_() bool
	BeforeTraverseNode(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
	AfterVisitNode(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node, *Reference)
	AfterExitScope(*NodeTraversal, ReferenceMap)
}

func ToIface_ReferenceCollector_Behavior[T interface {
	ReferenceCollector_Behavior
	comparable
}](v T) ReferenceCollector_Behavior {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ReferenceCollector_Behavior[T ReferenceCollector_Behavior](v T) ReferenceCollector_Behavior {
	var vi ReferenceCollector_Behavior = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReferenceCollector.Behavior", "com.google.javascript.jscomp.ReferenceCollector$Behavior", func(v_ any) bool {
		return j2gort.InstanceOf[ReferenceCollector_Behavior](v_)
	})
}

// Java default interface method: Behavior.beforeTraverseNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func IfaceDefault_ReferenceCollector_Behavior_BeforeTraverseNode(this_ ReferenceCollector_Behavior, t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = this_
	_ = t
	_ = n
	_ = parent
}

// Java default interface method: Behavior.afterVisitNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Reference addedReference)
func IfaceDefault_ReferenceCollector_Behavior_AfterVisitNode(this_ ReferenceCollector_Behavior, t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node, addedReference *Reference) {
	_ = this_
	_ = t
	_ = n
	_ = parent
	_ = addedReference
}

// Java class:
type ReferenceCollector__Anon_76 struct {
	idHash_ atomic.Int32
}

func (this_ *ReferenceCollector__Anon_76) BindSelf_(s any) {
	_ = s
}
func (this_ *ReferenceCollector__Anon_76) IsJavaNil_() bool { return this_ == nil }
func (*ReferenceCollector__Anon_76) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReferenceCollector$1", ""
}
func (this_ *ReferenceCollector__Anon_76) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReferenceCollector__Anon_76) HasIdentityEquals_()                            {}
func (*ReferenceCollector__Anon_76) Implements_jscomp_ReferenceCollector_Behavior() {}

func (this_ *ReferenceCollector__Anon_76) initFields_ReferenceCollector__Anon_76() {
	_ = this_
}

func NewReferenceCollector__Anon_76() *ReferenceCollector__Anon_76 {
	this_ := &ReferenceCollector__Anon_76{}
	this_.initFields_ReferenceCollector__Anon_76()
	return this_
}

// Java method: .afterExitScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *ReferenceCollector__Anon_76) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
}

func (this_ *ReferenceCollector__Anon_76) BeforeTraverseNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node) {
	IfaceDefault_ReferenceCollector_Behavior_BeforeTraverseNode(this_, p0_, p1_, p2_)
}

func (this_ *ReferenceCollector__Anon_76) AfterVisitNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node, p3_ *Reference) {
	IfaceDefault_ReferenceCollector_Behavior_AfterVisitNode(this_, p0_, p1_, p2_, p3_)
}
func (this_ *ReferenceCollector) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_BasicBlock.Release()
	this_.slab_ReferenceCollection.ReleaseWith(func(rc_ *ReferenceCollection) {
		j2gort.ReleasePooledSlice(rc_.rest)
		rc_.rest = nil
	})
	this_.slab_Reference.Release()
	this_.referenceMap.ReleaseToPool()
	if this_.ownsScopeCreator_ && this_.scopeCreator != nil {
		if r_, ok_ := any(this_.scopeCreator).(interface{ ReleaseSlabs_() }); ok_ {
			r_.ReleaseSlabs_()
		}
		this_.scopeCreator = nil
	}
}
