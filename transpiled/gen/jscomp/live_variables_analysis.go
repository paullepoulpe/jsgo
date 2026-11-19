/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/LiveVariablesAnalysis.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.LiveVariablesAnalysis
type LiveVariablesAnalysis struct {
	DataFlowAnalysisBase[pkg_rhino.Node, *LiveVariablesAnalysis_LiveVariableLattice]
	jsScope                       *Scope                                     `gson:"jsScope"`
	jsScopeChild                  *Scope                                     `gson:"jsScopeChild"`
	escaped                       j2gort.Ext_java_util_Set                   `gson:"escaped"`
	scopeVariables                j2gort.OrderedMap[string, *j2gort.BoxInt]  `gson:"scopeVariables"`
	orderedVars                   j2gort.Ext_java_util_List                  `gson:"orderedVars"`
	allVarsInFn                   j2gort.Ext_java_util_Map                   `gson:"allVarsInFn"`
	emptyLattice                  *LiveVariablesAnalysis_LiveVariableLattice `gson:"emptyLattice"`
	scratchGen                    *LiveVariablesAnalysis_LiveVariableLattice `gson:"scratchGen"`
	scratchKill                   *LiveVariablesAnalysis_LiveVariableLattice `gson:"scratchKill"`
	killLhsConsumer               j2gort.Ext_java_util_function_Consumer     `gson:"killLhsConsumer"`
	killDestructuringNameConsumer j2gort.Ext_java_util_function_Consumer     `gson:"killDestructuringNameConsumer"`
	slab_LiveVariableLattice      j2gort.PooledSlab[LiveVariablesAnalysis_LiveVariableLattice]
	slab_LiveVariableJoinOp       j2gort.PooledSlab[LiveVariablesAnalysis_LiveVariableJoinOp]
}

func (this_ *LiveVariablesAnalysis) BindSelf_(s any) {
	_ = s
	this_.DataFlowAnalysisBase.BindSelf_(s)
}
func (this_ *LiveVariablesAnalysis) IsJavaNil_() bool { return this_ == nil }
func (*LiveVariablesAnalysis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LiveVariablesAnalysis", "LiveVariablesAnalysis"
}
func (this_ *LiveVariablesAnalysis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.DataFlowAnalysisBase.IdentityHash_()
}
func (*LiveVariablesAnalysis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LiveVariablesAnalysis", "com.google.javascript.jscomp.LiveVariablesAnalysis", func(v_ any) bool {
		return j2gort.InstanceOf[*LiveVariablesAnalysis](v_)
	})
}
func (*LiveVariablesAnalysis) Implements_jscomp_DataFlowAnalysis() {}

func StaticGet_LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE() int32 { return int32(100) }
func (this_ *LiveVariablesAnalysis) initFields_LiveVariablesAnalysis() {
	_ = this_
}

// Java constructor: LiveVariablesAnalysis(com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg, com.google.javascript.jscomp.Scope jsScope, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope jsScopeChild, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ScopeCreator scopeCreator, com.google.javascript.jscomp.NodeUtil.AllVarsDeclaredInFunction allVarsDeclaredInFunction)
func NewLiveVariablesAnalysis(cfg ControlFlowGraph[pkg_rhino.Node], jsScope *Scope, jsScopeChild *Scope, compiler AbstractCompiler, scopeCreator ScopeCreator, allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction) *LiveVariablesAnalysis {
	this_ := &LiveVariablesAnalysis{}
	this_.BindSelf_(this_)
	InitLiveVariablesAnalysis(this_, cfg, jsScope, jsScopeChild, compiler, scopeCreator, allVarsDeclaredInFunction)
	return this_
}

func InitLiveVariablesAnalysis(this_ *LiveVariablesAnalysis, cfg ControlFlowGraph[pkg_rhino.Node], jsScope *Scope, jsScopeChild *Scope, compiler AbstractCompiler, scopeCreator ScopeCreator, allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction) {
	_ = cfg
	_ = jsScope
	_ = jsScopeChild
	_ = compiler
	_ = scopeCreator
	_ = allVarsDeclaredInFunction
	InitDataFlowAnalysis[pkg_rhino.Node, *LiveVariablesAnalysis_LiveVariableLattice](&this_.DataFlowAnalysisBase, cfg)
	this_.initFields_LiveVariablesAnalysis()
	if !(jsScope.IsFunctionScope()) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(jsScope))
	}
	this_.jsScope = jsScope
	this_.jsScopeChild = jsScopeChild
	this_.escaped = j2gort.NewExt_java_util_LinkedHashSet()
	this_.orderedVars = allVarsDeclaredInFunction.GetAllVariablesInOrder()
	var varCount int32 = this_.orderedVars.Size()
	_ = varCount
	this_.scopeVariables = j2gort.NewOrderedMapCap[string, *j2gort.BoxInt](varCount)
	this_.allVarsInFn = allVarsDeclaredInFunction.GetAllVariables()
	this_.emptyLattice = allocPooled_LiveVariablesAnalysis_LiveVariableLattice_1_int32(&this_.slab_LiveVariableLattice, varCount)
	this_.scratchGen = allocPooled_LiveVariablesAnalysis_LiveVariableLattice_1_int32(&this_.slab_LiveVariableLattice, varCount)
	this_.scratchKill = allocPooled_LiveVariablesAnalysis_LiveVariableLattice_1_int32(&this_.slab_LiveVariableLattice, varCount)
	DataFlowAnalysis_ComputeEscaped(jsScope, this_.escaped, compiler, scopeCreator, this_.allVarsInFn)
	this_.priv_LiveVariablesAnalysis_addScopeVariables_3780()
}

// Java method: LiveVariablesAnalysis.addScopeVariables()
func (this_ *LiveVariablesAnalysis) priv_LiveVariablesAnalysis_addScopeVariables_3780() {
	var num int32 = int32(0)
	_ = num
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.orderedVars))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			v := j2gort.CastPtr[Var](rawElem_2)
			_ = v
			this_.scopeVariables.Put(v.GetName(), j2gort.ValueOfInt(num))
			num = num + 1
		}
	}
}

// Java method: LiveVariablesAnalysis.getEscapedLocals()
func (this_ *LiveVariablesAnalysis) GetEscapedLocals() j2gort.Ext_java_util_Set {
	return this_.escaped
}

// Java method: LiveVariablesAnalysis.getAllVariables()
func (this_ *LiveVariablesAnalysis) GetAllVariables() j2gort.Ext_java_util_Map {
	return this_.allVarsInFn
}

// Java method: LiveVariablesAnalysis.getAllVariablesInOrder()
func (this_ *LiveVariablesAnalysis) GetAllVariablesInOrder() j2gort.Ext_java_util_List {
	return this_.orderedVars
}

// Java method: LiveVariablesAnalysis.getVarIndex(java.lang.String var)
func (this_ *LiveVariablesAnalysis) GetVarIndex(var_ string) int32 {
	_ = var_
	return j2gort.UnboxInt(this_.scopeVariables.Get(var_))
}

// Java method: LiveVariablesAnalysis.isForward()
func (this_ *LiveVariablesAnalysis) IsForward() bool {
	return false
}

// Java method: LiveVariablesAnalysis.createEntryLattice()
func (this_ *LiveVariablesAnalysis) CreateEntryLattice() *LiveVariablesAnalysis_LiveVariableLattice {
	return this_.emptyLattice
}

// Java method: LiveVariablesAnalysis.createInitialEstimateLattice()
func (this_ *LiveVariablesAnalysis) CreateInitialEstimateLattice() *LiveVariablesAnalysis_LiveVariableLattice {
	return this_.emptyLattice
}

// Java method: LiveVariablesAnalysis.createFlowJoiner()
func (this_ *LiveVariablesAnalysis) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*LiveVariablesAnalysis_LiveVariableLattice] {
	return ToIface_DataFlowAnalysis_FlowJoiner[*LiveVariablesAnalysis_LiveVariableLattice](allocPooled_LiveVariablesAnalysis_LiveVariableJoinOp(&this_.slab_LiveVariableJoinOp, this_))
}

// Java method: LiveVariablesAnalysis.flowThrough(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice input)
func (this_ *LiveVariablesAnalysis) FlowThrough(node pkg_rhino.Node, input *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	_ = node
	_ = input
	var gen *LiveVariablesAnalysis_LiveVariableLattice = this_.scratchGen
	_ = gen
	var kill *LiveVariablesAnalysis_LiveVariableLattice = this_.scratchKill
	_ = kill
	gen.ClearAll()
	kill.ClearAll()
	var conditional bool = false
	_ = conditional
	var cfgNode pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = pkg_graph.ToIface_DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch](j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.GetCfg().GetNode(node)))
	_ = cfgNode
	{
		var i int32 = int32(0)
		_ = i
		var outCount int32 = cfgNode.GetOutEdgeCount()
		_ = outCount
		for ; i < outCount; i = i + 1 {
			{
				var iifeCond_1 bool
				{

					re_3 := StaticGet_ControlFlowGraph_Branch_ON_EX()
					ae_4 := cfgNode.GetOutEdgeAt(i).GetValue()
					if re_3 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_3 == ae_4
				}
				if iifeCond_1 {
					conditional = true
					break
				}
			}
		}
	}
	this_.priv_LiveVariablesAnalysis_computeGenKill_3787(node, gen, kill, conditional)
	if ((input.overflow == nil) && (gen.overflow == nil)) && (kill.overflow == nil) {
		var nextW0 int64 = ((input.w0 & ^(kill.w0)) | gen.w0)
		_ = nextW0
		var nextW1 int64 = ((input.w1 & ^(kill.w1)) | gen.w1)
		_ = nextW1
		if (nextW0 == input.w0) && (nextW1 == input.w1) {
			return input
		}
		if (nextW0 == int64(0)) && (nextW1 == int64(0)) {
			return this_.emptyLattice
		}
		return allocPooled_LiveVariablesAnalysis_LiveVariableLattice_3(&this_.slab_LiveVariableLattice, input.numVars, nextW0, nextW1)
	}
	var result *LiveVariablesAnalysis_LiveVariableLattice = this_.slab_LiveVariableLattice.Alloc()
	InitLiveVariablesAnalysis_LiveVariableLattice_1_LiveVariablesAnalysis_LiveVariableLattice(result, input)

	_ = result
	result.AndNotFrom(kill)
	result.OrFrom(gen)
	return result
}

// Java method: LiveVariablesAnalysis.getKillLhsConsumer()
func (this_ *LiveVariablesAnalysis) priv_LiveVariablesAnalysis_getKillLhsConsumer_3785() j2gort.Ext_java_util_function_Consumer {
	var c j2gort.Ext_java_util_function_Consumer = this_.killLhsConsumer
	_ = c
	if c == nil {
		c = j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_6 any) {
			lhsNode := j2gort.CastNode[pkg_rhino.Node](lp_6)
			_ = lhsNode
			this_.priv_LiveVariablesAnalysis_addToSetIfLocal_3788(lhsNode, this_.scratchKill)
		})
		this_.killLhsConsumer = c
	}
	return c
}

// Java method: LiveVariablesAnalysis.getKillDestructuringNameConsumer()
func (this_ *LiveVariablesAnalysis) priv_LiveVariablesAnalysis_getKillDestructuringNameConsumer_3786() j2gort.Ext_java_util_function_Consumer {
	var c j2gort.Ext_java_util_function_Consumer = this_.killDestructuringNameConsumer
	_ = c
	if c == nil {
		c = j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_7 any) {
			child := j2gort.CastNode[pkg_rhino.Node](lp_7)
			_ = child
			if child.IsName() {
				this_.priv_LiveVariablesAnalysis_addToSetIfLocal_3788(child, this_.scratchKill)
			}
		})
		this_.killDestructuringNameConsumer = c
	}
	return c
}

// Java method: LiveVariablesAnalysis.computeGenKill(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice gen, com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice kill, boolean conditional)
func (this_ *LiveVariablesAnalysis) priv_LiveVariablesAnalysis_computeGenKill_3787(n pkg_rhino.Node, gen *LiveVariablesAnalysis_LiveVariableLattice, kill *LiveVariablesAnalysis_LiveVariableLattice, conditional bool) {
	_ = n
	_ = gen
	_ = kill
	_ = conditional
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */, int32(89) /* Token.ROOT */, int32(73) /* Token.FUNCTION */, int32(90) /* Token.BLOCK */ :
		return
	case int32(78) /* Token.WHILE */, int32(79) /* Token.DO */, int32(74) /* Token.IF */, int32(80) /* Token.FOR */ :
		this_.priv_LiveVariablesAnalysis_computeGenKill_3787(NodeUtil_GetConditionExpression(n), gen, kill, conditional)
		return
	case int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(81) /* Token.FOR_IN */ :
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		if NodeUtil_IsNameDeclaration(lhs) {
			lhs = lhs.GetLastChild()
		}
		this_.priv_LiveVariablesAnalysis_computeGenKill_3787(lhs, gen, kill, conditional)
		return
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(84) /* Token.VAR */ :
		{
			var c pkg_rhino.Node = n.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				if c.IsName() {
					if c.HasChildren() {
						this_.priv_LiveVariablesAnalysis_computeGenKill_3787(c.GetFirstChild(), gen, kill, conditional)
						if !(conditional) {
							this_.priv_LiveVariablesAnalysis_addToSetIfLocal_3788(c, kill)
						}
					}
				} else {
					if !(c.IsDestructuringLhs()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(c))
					}
					if !(conditional) {
						NodeUtil_VisitLhsNodesInNode(c, this_.priv_LiveVariablesAnalysis_getKillLhsConsumer_3785())
					}
					this_.priv_LiveVariablesAnalysis_computeGenKill_3787(c.GetFirstChild(), gen, kill, conditional)
					this_.priv_LiveVariablesAnalysis_computeGenKill_3787(c.GetSecondChild(), gen, kill, conditional)
				}
			}
		}
		return
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_LiveVariablesAnalysis_computeGenKill_3787(n.GetFirstChild(), gen, kill, conditional)
		this_.priv_LiveVariablesAnalysis_computeGenKill_3787(n.GetLastChild(), gen, kill, true)
		return
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		this_.priv_LiveVariablesAnalysis_computeGenKill_3787(n.GetFirstChild(), gen, kill, conditional)
		{
			var c_2 pkg_rhino.Node = n.GetSecondChild()
			_ = c_2
			for ; c_2 != 0; c_2 = c_2.GetNext() {
				this_.priv_LiveVariablesAnalysis_computeGenKill_3787(c_2, gen, kill, true)
			}
		}
		return
	case int32(67) /* Token.HOOK */ :
		this_.priv_LiveVariablesAnalysis_computeGenKill_3787(n.GetFirstChild(), gen, kill, conditional)
		this_.priv_LiveVariablesAnalysis_computeGenKill_3787(n.GetSecondChild(), gen, kill, true)
		this_.priv_LiveVariablesAnalysis_computeGenKill_3787(n.GetLastChild(), gen, kill, true)
		return
	case int32(32) /* Token.NAME */ :
		if j2gort.CheckNotNullString(n.GetString()) == "arguments" {
			this_.MarkAllParametersEscaped()
		} else if !(NodeUtil_IsLhsByDestructuring(n)) {
			this_.priv_LiveVariablesAnalysis_addToSetIfLocal_3788(n, gen)
		}
		return
	default:
		if NodeUtil_IsAssignmentOp(n) && n.GetFirstChild().IsName() {
			var lhs_2 pkg_rhino.Node = n.GetFirstChild()
			_ = lhs_2
			if !(conditional) {
				this_.priv_LiveVariablesAnalysis_addToSetIfLocal_3788(lhs_2, kill)
			}
			if !(n.IsAssign()) {
				this_.priv_LiveVariablesAnalysis_addToSetIfLocal_3788(lhs_2, gen)
			}
			this_.priv_LiveVariablesAnalysis_computeGenKill_3787(lhs_2.GetNext(), gen, kill, conditional)
		} else if n.IsAssign() && n.GetFirstChild().IsDestructuringPattern() {
			if !(conditional) {
				NodeUtil_VisitLhsNodesInNode(n, this_.priv_LiveVariablesAnalysis_getKillDestructuringNameConsumer_3786())
			}
			this_.priv_LiveVariablesAnalysis_computeGenKill_3787(n.GetFirstChild(), gen, kill, conditional)
			this_.priv_LiveVariablesAnalysis_computeGenKill_3787(n.GetSecondChild(), gen, kill, conditional)
		} else {
			{
				var c_3 pkg_rhino.Node = n.GetFirstChild()
				_ = c_3
				for ; c_3 != 0; c_3 = c_3.GetNext() {
					this_.priv_LiveVariablesAnalysis_computeGenKill_3787(c_3, gen, kill, conditional)
				}
			}
		}
		return
	}
}

// Java method: LiveVariablesAnalysis.addToSetIfLocal(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice set)
func (this_ *LiveVariablesAnalysis) priv_LiveVariablesAnalysis_addToSetIfLocal_3788(node pkg_rhino.Node, set *LiveVariablesAnalysis_LiveVariableLattice) {
	_ = node
	_ = set
	if !(node.IsName()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
	}
	var name string = node.GetString()
	_ = name
	var var_ *Var = j2gort.CastPtr[Var](j2gort.MapGetStr(this_.allVarsInFn, name))
	_ = var_
	if var_ == nil {
		return
	}
	var local bool = false
	_ = local
	var localScope *Scope = j2gort.CastPtr[Scope](var_.GetScope())
	_ = localScope
	if localScope.IsFunctionBlockScope() {
		local = LiveVariablesAnalysis_IsDeclaredInFunctionBlockOrParameter(localScope, name)
	} else if (localScope == this_.jsScope) && (this_.jsScopeChild != nil) {
		local = LiveVariablesAnalysis_IsDeclaredInFunctionBlockOrParameter(this_.jsScopeChild, name)
	} else {
		local = localScope.HasOwnSlot(name)
	}
	if !(local) {
		return
	}
	if !(this_.escaped.Contains(j2gort.PtrToAny(var_))) {
		set.SetLive(this_.GetVarIndex(var_.GetName()))
	}
}

// Java static method: LiveVariablesAnalysis.isDeclaredInFunctionBlockOrParameter(com.google.javascript.jscomp.Scope scope, java.lang.String name)
func LiveVariablesAnalysis_IsDeclaredInFunctionBlockOrParameter(scope *Scope, name string) bool {
	_ = scope
	_ = name
	if !(scope.IsFunctionBlockScope()) {
		j2gort.CheckState(false)
	}
	return (scope.HasOwnSlot(name) || j2gort.CastPtr[Scope](scope.GetParent()).HasOwnSlot(name))
}

// Java method: LiveVariablesAnalysis.markAllParametersEscaped()
func (this_ *LiveVariablesAnalysis) MarkAllParametersEscaped() {
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(this_.jsScope.GetRootNode_0())
	_ = paramList
	{
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		for ; param != 0; param = param.GetNext() {
			if param.IsName() {
				this_.escaped.Add_1(j2gort.PtrToAny(j2gort.CastPtr[Var](this_.jsScope.GetVar(param.GetString()))))
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableJoinOp
type LiveVariablesAnalysis_LiveVariableJoinOp struct {
	outer_  *LiveVariablesAnalysis
	result  *LiveVariablesAnalysis_LiveVariableLattice `gson:"result"`
	idHash_ atomic.Int32
}

func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) BindSelf_(s any) {
	_ = s
}
func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) IsJavaNil_() bool { return this_ == nil }
func (*LiveVariablesAnalysis_LiveVariableJoinOp) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LiveVariablesAnalysis$LiveVariableJoinOp", "LiveVariableJoinOp"
}
func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LiveVariablesAnalysis_LiveVariableJoinOp) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableJoinOp", "com.google.javascript.jscomp.LiveVariablesAnalysis$LiveVariableJoinOp", func(v_ any) bool {
		return j2gort.InstanceOf[*LiveVariablesAnalysis_LiveVariableJoinOp](v_)
	})
}
func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) SetOuter_(o *LiveVariablesAnalysis) {
	this_.outer_ = o
}
func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) GetOuter_() *LiveVariablesAnalysis {
	return this_.outer_
}
func (*LiveVariablesAnalysis_LiveVariableJoinOp) Implements_jscomp_DataFlowAnalysis_FlowJoiner() {}

func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) initFields_LiveVariablesAnalysis_LiveVariableJoinOp() {
	_ = this_
	this_.result = allocPooled_LiveVariablesAnalysis_LiveVariableLattice_1_int32(&this_.outer_.slab_LiveVariableLattice, this_.outer_.orderedVars.Size())
}

// Java constructor: LiveVariableJoinOp()
func NewLiveVariablesAnalysis_LiveVariableJoinOp(outer_ *LiveVariablesAnalysis) *LiveVariablesAnalysis_LiveVariableJoinOp {
	this_ := &LiveVariablesAnalysis_LiveVariableJoinOp{}
	this_.outer_ = outer_
	InitLiveVariablesAnalysis_LiveVariableJoinOp(this_)
	return this_
}
func allocPooled_LiveVariablesAnalysis_LiveVariableJoinOp(slab_ *j2gort.PooledSlab[LiveVariablesAnalysis_LiveVariableJoinOp], outer_ *LiveVariablesAnalysis) *LiveVariablesAnalysis_LiveVariableJoinOp {
	this_ := slab_.Alloc()
	this_.outer_ = outer_
	InitLiveVariablesAnalysis_LiveVariableJoinOp(this_)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableJoinOp(this_ *LiveVariablesAnalysis_LiveVariableJoinOp) {
	this_.initFields_LiveVariablesAnalysis_LiveVariableJoinOp()
}

// Java method: LiveVariableJoinOp.joinFlow(com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice x)
func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) JoinFlow(x *LiveVariablesAnalysis_LiveVariableLattice) {
	_ = x
	this_.result.OrFrom(x)
}

// Java method: LiveVariableJoinOp.finish()
func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) Finish() *LiveVariablesAnalysis_LiveVariableLattice {
	if this_.result.IsEmpty() {
		return this_.outer_.emptyLattice
	}
	return this_.result
}

// Java class: com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice
type LiveVariablesAnalysis_LiveVariableLattice struct {
	overflow *j2gort.Ext_java_util_BitSet `gson:"overflow"`
	w0       int64                        `gson:"w0"`
	w1       int64                        `gson:"w1"`
	numVars  int32                        `gson:"numVars"`
}

func (this_ *LiveVariablesAnalysis_LiveVariableLattice) BindSelf_(s any) {
	_ = s
}
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) IsJavaNil_() bool { return this_ == nil }
func (*LiveVariablesAnalysis_LiveVariableLattice) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LiveVariablesAnalysis$LiveVariableLattice", "LiveVariableLattice"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice", "com.google.javascript.jscomp.LiveVariablesAnalysis$LiveVariableLattice", func(v_ any) bool {
		return j2gort.InstanceOf[*LiveVariablesAnalysis_LiveVariableLattice](v_)
	})
}
func (*LiveVariablesAnalysis_LiveVariableLattice) Implements_graph_LatticeElement() {}
func (*LiveVariablesAnalysis_LiveVariableLattice) Implements_graph_Annotation()     {}

func (this_ *LiveVariablesAnalysis_LiveVariableLattice) initFields_LiveVariablesAnalysis_LiveVariableLattice() {
	_ = this_
}

// Java constructor: LiveVariableLattice(int numVars)
func NewLiveVariablesAnalysis_LiveVariableLattice_1_int32(numVars int32) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := &LiveVariablesAnalysis_LiveVariableLattice{}
	InitLiveVariablesAnalysis_LiveVariableLattice_1_int32(this_, numVars)
	return this_
}
func allocPooled_LiveVariablesAnalysis_LiveVariableLattice_1_int32(slab_ *j2gort.PooledSlab[LiveVariablesAnalysis_LiveVariableLattice], numVars int32) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := slab_.Alloc()
	InitLiveVariablesAnalysis_LiveVariableLattice_1_int32(this_, numVars)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableLattice_1_int32(this_ *LiveVariablesAnalysis_LiveVariableLattice, numVars int32) {
	_ = numVars
	this_.initFields_LiveVariablesAnalysis_LiveVariableLattice()
	if !(numVars >= int32(0)) {
		j2gort.CheckArgument(false)
	}
	this_.numVars = numVars
	if numVars > int32(128) {
		this_.overflow = j2gort.NewExt_java_util_BitSet_1(numVars)
	} else {
		this_.overflow = nil
	}
}

// Java constructor: LiveVariableLattice(int numVars, long w0, long w1)
func NewLiveVariablesAnalysis_LiveVariableLattice_3(numVars int32, w0 int64, w1 int64) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := &LiveVariablesAnalysis_LiveVariableLattice{}
	InitLiveVariablesAnalysis_LiveVariableLattice_3(this_, numVars, w0, w1)
	return this_
}
func allocPooled_LiveVariablesAnalysis_LiveVariableLattice_3(slab_ *j2gort.PooledSlab[LiveVariablesAnalysis_LiveVariableLattice], numVars int32, w0 int64, w1 int64) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := slab_.Alloc()
	InitLiveVariablesAnalysis_LiveVariableLattice_3(this_, numVars, w0, w1)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableLattice_3(this_ *LiveVariablesAnalysis_LiveVariableLattice, numVars int32, w0 int64, w1 int64) {
	_ = numVars
	_ = w0
	_ = w1
	this_.initFields_LiveVariablesAnalysis_LiveVariableLattice()
	this_.numVars = numVars
	this_.w0 = w0
	this_.w1 = w1
	this_.overflow = nil
}

// Java constructor: LiveVariableLattice(com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice other)
func NewLiveVariablesAnalysis_LiveVariableLattice_1_LiveVariablesAnalysis_LiveVariableLattice(other *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := &LiveVariablesAnalysis_LiveVariableLattice{}
	InitLiveVariablesAnalysis_LiveVariableLattice_1_LiveVariablesAnalysis_LiveVariableLattice(this_, other)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableLattice_1_LiveVariablesAnalysis_LiveVariableLattice(this_ *LiveVariablesAnalysis_LiveVariableLattice, other *LiveVariablesAnalysis_LiveVariableLattice) {
	_ = other
	this_.initFields_LiveVariablesAnalysis_LiveVariableLattice()
	j2gort.CheckNotNullPtr(other)
	this_.numVars = other.numVars
	this_.w0 = other.w0
	this_.w1 = other.w1
	if other.overflow != nil {
		this_.overflow = j2gort.CastPtr[j2gort.Ext_java_util_BitSet](other.overflow.Clone())
	} else {
		this_.overflow = nil
	}
}

// Java method: LiveVariableLattice.ensureOverflow()
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) priv_LiveVariableLattice_ensureOverflow_3790() *j2gort.Ext_java_util_BitSet {
	var bs *j2gort.Ext_java_util_BitSet = this_.overflow
	_ = bs
	if bs == nil {
		bs = j2gort.NewExt_java_util_BitSet_1(max(this_.numVars, int32(129)))
		var word0 int64 = this_.w0
		_ = word0
		for word0 != int64(0) {
			var bit int32 = j2gort.BoxLong_NumberOfTrailingZeros(word0)
			_ = bit
			bs.Set_1(bit)
			word0 = (word0 & (word0 - int64(1)))
		}
		var word1 int64 = this_.w1
		_ = word1
		for word1 != int64(0) {
			var bit_2 int32 = j2gort.BoxLong_NumberOfTrailingZeros(word1)
			_ = bit_2
			bs.Set_1((int32(64) + bit_2))
			word1 = (word1 & (word1 - int64(1)))
		}
		this_.w0 = int64(0)
		this_.w1 = int64(0)
		this_.overflow = bs
	}
	return bs
}

// Java method: LiveVariableLattice.isEmpty()
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) IsEmpty() bool {
	if this_.overflow != nil {
		return this_.overflow.IsEmpty()
	}
	return ((this_.w0 == int64(0)) && (this_.w1 == int64(0)))
}

// Java method: LiveVariableLattice.setLive(int index)
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) SetLive(index int32) {
	_ = index
	if index < int32(0) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", "bitIndex < 0: "+j2gort.StrValInt32(index))))
	}
	if this_.overflow != nil {
		this_.overflow.Set_1(index)
	} else if index < int32(64) {
		this_.w0 = (this_.w0 | (javalang.Shl64(int64(1), int(index))))
	} else if index < int32(128) {
		this_.w1 = (this_.w1 | (javalang.Shl64(int64(1), int((index - int32(64))))))
	} else {
		this_.priv_LiveVariableLattice_ensureOverflow_3790().Set_1(index)
	}
}

// Java method: LiveVariableLattice.clearLive(int index)
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) ClearLive(index int32) {
	_ = index
	if index < int32(0) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", "bitIndex < 0: "+j2gort.StrValInt32(index))))
	}
	if this_.overflow != nil {
		this_.overflow.Clear_1(index)
	} else if index < int32(64) {
		this_.w0 = (this_.w0 & ^(javalang.Shl64(int64(1), int(index))))
	} else if index < int32(128) {
		this_.w1 = (this_.w1 & ^(javalang.Shl64(int64(1), int((index - int32(64))))))
	}
}

// Java method: LiveVariableLattice.clearAll()
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) ClearAll() {
	this_.w0 = int64(0)
	this_.w1 = int64(0)
	if this_.overflow != nil {
		this_.overflow.Clear_0()
	}
}

// Java method: LiveVariableLattice.orFrom(com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice other)
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) OrFrom(other *LiveVariablesAnalysis_LiveVariableLattice) {
	_ = other
	if (this_.overflow == nil) && (other.overflow == nil) {
		this_.w0 = (this_.w0 | other.w0)
		this_.w1 = (this_.w1 | other.w1)
	} else {
		var dst *j2gort.Ext_java_util_BitSet = this_.priv_LiveVariableLattice_ensureOverflow_3790()
		_ = dst
		if other.overflow != nil {
			dst.Or(other.overflow)
		} else {
			var word0 int64 = other.w0
			_ = word0
			for word0 != int64(0) {
				var bit int32 = j2gort.BoxLong_NumberOfTrailingZeros(word0)
				_ = bit
				dst.Set_1(bit)
				word0 = (word0 & (word0 - int64(1)))
			}
			var word1 int64 = other.w1
			_ = word1
			for word1 != int64(0) {
				var bit_2 int32 = j2gort.BoxLong_NumberOfTrailingZeros(word1)
				_ = bit_2
				dst.Set_1((int32(64) + bit_2))
				word1 = (word1 & (word1 - int64(1)))
			}
		}
	}
}

// Java method: LiveVariableLattice.andNotFrom(com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice other)
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) AndNotFrom(other *LiveVariablesAnalysis_LiveVariableLattice) {
	_ = other
	if (this_.overflow == nil) && (other.overflow == nil) {
		this_.w0 = (this_.w0 & ^(other.w0))
		this_.w1 = (this_.w1 & ^(other.w1))
	} else if this_.overflow == nil {
		{
			var i int32 = other.overflow.NextSetBit(int32(0))
			_ = i
			for ; (i >= int32(0)) && (i < int32(128)); i = other.overflow.NextSetBit((i + int32(1))) {
				if i < int32(64) {
					this_.w0 = (this_.w0 & ^(javalang.Shl64(int64(1), int(i))))
				} else {
					this_.w1 = (this_.w1 & ^(javalang.Shl64(int64(1), int((i - int32(64))))))
				}
			}
		}
	} else if other.overflow != nil {
		this_.overflow.AndNot(other.overflow)
	} else {
		var word0 int64 = other.w0
		_ = word0
		for word0 != int64(0) {
			var bit int32 = j2gort.BoxLong_NumberOfTrailingZeros(word0)
			_ = bit
			this_.overflow.Clear_1(bit)
			word0 = (word0 & (word0 - int64(1)))
		}
		var word1 int64 = other.w1
		_ = word1
		for word1 != int64(0) {
			var bit_2 int32 = j2gort.BoxLong_NumberOfTrailingZeros(word1)
			_ = bit_2
			this_.overflow.Clear_1((int32(64) + bit_2))
			word1 = (word1 & (word1 - int64(1)))
		}
	}
}

// Java method: LiveVariableLattice.equals(java.lang.Object other)
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) Equals(arg_8 any) bool {
	other := arg_8
	_ = other
	j2gort.CheckNotNull(other)
	var that *LiveVariablesAnalysis_LiveVariableLattice
	_ = that
	var ok_9 bool
	that, ok_9 = j2gort.TryCast[*LiveVariablesAnalysis_LiveVariableLattice](other)
	if !ok_9 {
		return false
	}
	if this_ == that {
		return true
	}
	if (this_.overflow == nil) && (that.overflow == nil) {
		return ((this_.w0 == that.w0) && (this_.w1 == that.w1))
	}
	if (this_.overflow != nil) && (that.overflow != nil) {
		return this_.overflow.Equals(j2gort.PtrToAny(that.overflow))
	}
	var inline *LiveVariablesAnalysis_LiveVariableLattice
	_ = inline
	if this_.overflow == nil {
		inline = this_
	} else {
		inline = that
	}
	var bs *j2gort.Ext_java_util_BitSet
	_ = bs
	if this_.overflow != nil {
		bs = this_.overflow
	} else {
		bs = that.overflow
	}
	if bs.Length() > int32(128) {
		return false
	}
	var words []int64 = bs.ToLongArray()
	_ = words
	var bsW0 int64
	_ = bsW0
	if int32(len(words)) > int32(0) {
		bsW0 = words[int32(0)]
	} else {
		bsW0 = int64(0)
	}
	var bsW1 int64
	_ = bsW1
	if int32(len(words)) > int32(1) {
		bsW1 = words[int32(1)]
	} else {
		bsW1 = int64(0)
	}
	return ((inline.w0 == bsW0) && (inline.w1 == bsW1))
}

// Java method: LiveVariableLattice.isLive(int index)
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) IsLive(index int32) bool {
	_ = index
	if index < int32(0) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", "bitIndex < 0: "+j2gort.StrValInt32(index))))
	}
	if this_.overflow != nil {
		return this_.overflow.Get_1(index)
	}
	if index < int32(64) {
		return ((this_.w0 & (javalang.Shl64(int64(1), int(index)))) != int64(0))
	}
	if index < int32(128) {
		return ((this_.w1 & (javalang.Shl64(int64(1), int((index - int32(64)))))) != int64(0))
	}
	return false
}

// Java method: LiveVariableLattice.toString()
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) String() string {
	if this_.overflow != nil {
		return this_.overflow.String()
	}
	var first int32 = this_.NextSetBit(int32(0))
	_ = first
	if first < int32(0) {
		return "{}"
	}
	var sb j2gort.LocalStringBuilder

	sb.AppendChar(uint16(123)).AppendInt(first)
	{
		var i int32 = this_.NextSetBit((first + int32(1)))
		_ = i
		for ; i >= int32(0); i = this_.NextSetBit((i + int32(1))) {
			sb.AppendStr(", ").AppendInt(i)
		}
	}
	sb.AppendChar(uint16(125))
	return sb.String()
}

// Java method: LiveVariableLattice.nextSetBit(int fromIndex)
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) NextSetBit(fromIndex int32) int32 {
	_ = fromIndex
	if fromIndex < int32(0) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", "fromIndex < 0: "+j2gort.StrValInt32(fromIndex))))
	}
	if this_.overflow != nil {
		return this_.overflow.NextSetBit(fromIndex)
	}
	if fromIndex < int32(64) {
		var word int64 = (this_.w0 & (javalang.Shl64(^(int64(0)), int(fromIndex))))
		_ = word
		if word != int64(0) {
			return j2gort.BoxLong_NumberOfTrailingZeros(word)
		}
		if this_.w1 != int64(0) {
			return (int32(64) + j2gort.BoxLong_NumberOfTrailingZeros(this_.w1))
		}
		return int32(-1)
	}
	if fromIndex < int32(128) {
		var word_2 int64 = (this_.w1 & (javalang.Shl64(^(int64(0)), int((fromIndex - int32(64))))))
		_ = word_2
		if word_2 != int64(0) {
			return (int32(64) + j2gort.BoxLong_NumberOfTrailingZeros(word_2))
		}
	}
	return int32(-1)
}

// Java method: LiveVariableLattice.hashCode()
func (this_ *LiveVariablesAnalysis_LiveVariableLattice) HashCode() int32 {
	if this_.overflow != nil {
		return this_.overflow.HashCode()
	}
	var h int64 = ((int64(1234) ^ (this_.w1 * int64(2))) ^ this_.w0)
	_ = h
	return int32(((javalang.Shr64(h, int(int32(32)))) ^ h))
}
func (this_ *LiveVariablesAnalysis) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_LiveVariableLattice.Release()
	this_.slab_LiveVariableJoinOp.Release()
}
