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
// Port of com/google/javascript/jscomp/FlowSensitiveInlineVariables.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_flow_sensitive_inline_variables_1 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_flow_sensitive_inline_variables_2 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_flow_sensitive_inline_variables_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_3 uint32) bool {
		input := pkg_rhino.Node(rawNode_3)
		_ = input
		switch input.TokenSwitchOrd_() {
		case int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(105) /* Token.CLASS */, int32(46) /* Token.ARRAYLIT */, int32(47) /* Token.OBJECTLIT */, int32(42) /* Token.REGEXP */, int32(23) /* Token.NEW */ :
			return true
		default:
		}
		return false
	})
	lambdaSingleton_flow_sensitive_inline_variables_2 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_4 uint32) bool {
		input := pkg_rhino.Node(rawNode_4)
		_ = input
		return !(input.IsFunction())
	})
}

// Java class: com.google.javascript.jscomp.FlowSensitiveInlineVariables
type FlowSensitiveInlineVariables struct {
	compiler            AbstractCompiler                                  `gson:"compiler"`
	sideEffectPredicate *FlowSensitiveInlineVariables_SideEffectPredicate `gson:"sideEffectPredicate"`
	cfg                 ControlFlowGraph[pkg_rhino.Node]                  `gson:"cfg"`
	candidates          *j2gort.Ext_java_util_LinkedHashSet               `gson:"candidates"`
	reachingDef         *MustBeReachingVariableDef                        `gson:"reachingDef"`
	reachingUses        *MaybeReachingVariableUse                         `gson:"reachingUses"`
	idHash_             atomic.Int32
}

func (this_ *FlowSensitiveInlineVariables) BindSelf_(s any) {
	_ = s
}
func (this_ *FlowSensitiveInlineVariables) IsJavaNil_() bool { return this_ == nil }
func (*FlowSensitiveInlineVariables) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FlowSensitiveInlineVariables", "FlowSensitiveInlineVariables"
}
func (this_ *FlowSensitiveInlineVariables) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FlowSensitiveInlineVariables) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FlowSensitiveInlineVariables", "com.google.javascript.jscomp.FlowSensitiveInlineVariables", func(v_ any) bool {
		return j2gort.InstanceOf[*FlowSensitiveInlineVariables](v_)
	})
}
func (*FlowSensitiveInlineVariables) Implements_jscomp_CompilerPass()                 {}
func (*FlowSensitiveInlineVariables) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*FlowSensitiveInlineVariables) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *FlowSensitiveInlineVariables) initFields_FlowSensitiveInlineVariables() {
	_ = this_
}

// Java constructor: FlowSensitiveInlineVariables(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewFlowSensitiveInlineVariables(compiler AbstractCompiler) *FlowSensitiveInlineVariables {
	this_ := &FlowSensitiveInlineVariables{}
	InitFlowSensitiveInlineVariables(this_, compiler)
	return this_
}

func InitFlowSensitiveInlineVariables(this_ *FlowSensitiveInlineVariables, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_FlowSensitiveInlineVariables()
	this_.compiler = compiler
	this_.sideEffectPredicate = NewFlowSensitiveInlineVariables_SideEffectPredicate(this_)
}

// Java static method: FlowSensitiveInlineVariables.isTopLevelAssignTarget(com.google.javascript.rhino.Node n)
func FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n pkg_rhino.Node) bool {
	_ = n
	var ancestor pkg_rhino.Node = n.GetParent()
	_ = ancestor
	for ancestor.IsAssign() {
		ancestor = ancestor.GetParent()
	}
	return ancestor.IsExprResult()
}

// Java method: FlowSensitiveInlineVariables.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FlowSensitiveInlineVariables) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsScript()) || !(t.GetInput().IsExtern()))
}

// Java method: FlowSensitiveInlineVariables.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *FlowSensitiveInlineVariables) EnterScope(t *NodeTraversal) {
	_ = t
	if t.InGlobalScope() {
		return
	}
	if !(t.GetScope().IsFunctionBlockScope()) {
		return
	}
	var functionScopeRoot pkg_rhino.Node = t.GetScopeRoot().GetParent()
	_ = functionScopeRoot
	if !(this_.priv_FlowSensitiveInlineVariables_isCandidateFunction_2864(functionScopeRoot)) {
		return
	}
	if int32(100) < t.GetScope().GetVarCount() {
		return
	}
	var scopeCreator *SyntacticScopeCreator = j2gort.CastPtr[SyntacticScopeCreator](t.GetScopeCreator())
	_ = scopeCreator
	var localCfg ControlFlowGraph[pkg_rhino.Node] = ControlFlowAnalysis_ComputeCfgDirect(this_.compiler, functionScopeRoot, false, true)
	_ = localCfg
	defer func() {
		this_.cfg = nil
		this_.reachingDef = nil
		this_.reachingUses = nil
		this_.candidates = nil
		ControlFlowAnalysis_ReleaseCfg(localCfg)
	}()
	{
		this_.cfg = localCfg
		var escaped *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
		_ = escaped
		var scope *Scope = t.GetScope()
		_ = scope
		var allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction = NodeUtil_GetAllVarsDeclaredInFunction(this_.compiler, ToIface_ScopeCreator(scopeCreator), j2gort.CastPtr[Scope](scope.GetParent()))
		_ = allVarsDeclaredInFunction
		var allVarsInFn j2gort.Ext_java_util_Map = allVarsDeclaredInFunction.GetAllVariables()
		_ = allVarsInFn
		DataFlowAnalysis_ComputeEscaped(j2gort.CastPtr[Scope](scope.GetParent()), escaped, this_.compiler, ToIface_ScopeCreator(scopeCreator), allVarsInFn)
		this_.reachingDef = NewMustBeReachingVariableDef(this_.cfg, this_.compiler, escaped, allVarsInFn)
		this_.reachingDef.Analyze()
		this_.candidates = nil
		NodeTraversal_Traverse(this_.compiler, t.GetScopeRoot(), NewFlowSensitiveInlineVariables_GatherCandidates(this_))
		this_.reachingUses = NewMaybeReachingVariableUse(this_.cfg, escaped, allVarsInFn)
		this_.reachingUses.Analyze()
		for (this_.candidates != nil) && !(this_.candidates.IsEmpty()) {
			var c *FlowSensitiveInlineVariables_Candidate = j2gort.CastPtr[FlowSensitiveInlineVariables_Candidate](j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.candidates.Iterator()).Next())
			_ = c
			var candidateVar *Var = j2gort.CheckNotNullPtr(j2gort.CastPtr[Var](j2gort.MapGetStr(allVarsInFn, c.varName)))
			_ = candidateVar
			if c.priv_Candidate_canInline_2868(j2gort.CastPtr[Scope](candidateVar.GetScope())) {
				c.priv_Candidate_inlineVariable_2870()
				this_.candidates.Remove_1_any(j2gort.PtrToAny(c))
				if !(c.defMetadata.IsEmptyDepends()) {
					var cVar *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(c.varName))
					_ = cVar
					{
						var it j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.candidates.Iterator())
						_ = it
						for it.HasNext() {
							var other *FlowSensitiveInlineVariables_Candidate = j2gort.CastPtr[FlowSensitiveInlineVariables_Candidate](it.Next())
							_ = other
							if other.defMetadata.HasDependency(cVar) && !(other.defMetadata.ContainsAllDependencies(c.defMetadata)) {
								it.Remove()
							}
						}
					}
				}
			} else {
				this_.candidates.Remove_1_any(j2gort.PtrToAny(c))
			}
		}
	}
}

// Java method: FlowSensitiveInlineVariables.isCandidateFunction(com.google.javascript.rhino.Node fn)
func (this_ *FlowSensitiveInlineVariables) priv_FlowSensitiveInlineVariables_isCandidateFunction_2864(fn pkg_rhino.Node) bool {
	_ = fn
	var fnBody pkg_rhino.Node = fn.GetLastChild()
	_ = fnBody
	return FlowSensitiveInlineVariables_ContainsCandidateExpressions(fnBody)
}

// Java static method: FlowSensitiveInlineVariables.containsCandidateExpressions(com.google.javascript.rhino.Node n)
func FlowSensitiveInlineVariables_ContainsCandidateExpressions(n pkg_rhino.Node) bool {
	_ = n
	if n.IsFunction() {
		return false
	}
	if NodeUtil_IsNameDeclaration(n) || FlowSensitiveInlineVariables_IsAssignmentToName(n) {
		if n.GetFirstChild().IsName() {
			return true
		}
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if FlowSensitiveInlineVariables_ContainsCandidateExpressions(c) {
				return true
			}
		}
	}
	return false
}

// Java static method: FlowSensitiveInlineVariables.isAssignmentToName(com.google.javascript.rhino.Node n)
func FlowSensitiveInlineVariables_IsAssignmentToName(n pkg_rhino.Node) bool {
	_ = n
	if (NodeUtil_IsAssignmentOp(n) || n.IsDec()) || n.IsInc() {
		return (n.GetFirstChild().IsName())
	}
	return false
}

// Java method: FlowSensitiveInlineVariables.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *FlowSensitiveInlineVariables) ExitScope(t *NodeTraversal) {
	_ = t
}

// Java method: FlowSensitiveInlineVariables.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *FlowSensitiveInlineVariables) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(this_)
		(&ntb_).TraverseRoots(externs, root)
	}
}

// Java method: FlowSensitiveInlineVariables.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FlowSensitiveInlineVariables) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

// Java static method: FlowSensitiveInlineVariables.checkPostExpressions(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node expressionRoot, com.google.common.base.Predicate<com.google.javascript.rhino.Node> predicate)
func FlowSensitiveInlineVariables_CheckPostExpressions(n pkg_rhino.Node, expressionRoot pkg_rhino.Node, predicate j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = n
	_ = expressionRoot
	_ = predicate
	{
		var p pkg_rhino.Node = n
		_ = p
		for ; p != expressionRoot; p = p.GetParent() {
			{
				var cur pkg_rhino.Node = p.GetNext()
				_ = cur
				for ; cur != 0; cur = cur.GetNext() {
					if j2gort.PredTestNode(predicate, uint32(cur)) {
						return true
					}
				}
			}
		}
	}
	return false
}

// Java static method: FlowSensitiveInlineVariables.checkPreExpressions(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node expressionRoot, com.google.common.base.Predicate<com.google.javascript.rhino.Node> predicate)
func FlowSensitiveInlineVariables_CheckPreExpressions(n pkg_rhino.Node, expressionRoot pkg_rhino.Node, predicate j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = n
	_ = expressionRoot
	_ = predicate
	{
		var p pkg_rhino.Node = n
		_ = p
		for ; p != expressionRoot; p = p.GetParent() {
			var oldestSibling pkg_rhino.Node = p.GetParent().GetFirstChild()
			_ = oldestSibling
			if oldestSibling.IsDestructuringPattern() {
				if p.IsDestructuringPattern() {
					if (p.GetNext() != 0) && j2gort.PredTestNode(predicate, uint32(p.GetNext())) {
						return true
					}
				}
				continue
			}
			{
				var cur pkg_rhino.Node = oldestSibling
				_ = cur
				for ; cur != p; cur = cur.GetNext() {
					if j2gort.PredTestNode(predicate, uint32(cur)) {
						return true
					}
				}
			}
		}
	}
	return false
}

// Java class: com.google.javascript.jscomp.FlowSensitiveInlineVariables.SideEffectPredicate
type FlowSensitiveInlineVariables_SideEffectPredicate struct {
	outer_       *FlowSensitiveInlineVariables
	namesToCheck j2gort.Ext_java_util_Set `gson:"namesToCheck"`
	idHash_      atomic.Int32
}

func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) BindSelf_(s any) {
	_ = s
}
func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) IsJavaNil_() bool { return this_ == nil }
func (*FlowSensitiveInlineVariables_SideEffectPredicate) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FlowSensitiveInlineVariables$SideEffectPredicate", "SideEffectPredicate"
}
func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FlowSensitiveInlineVariables_SideEffectPredicate) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FlowSensitiveInlineVariables.SideEffectPredicate", "com.google.javascript.jscomp.FlowSensitiveInlineVariables$SideEffectPredicate", func(v_ any) bool {
		return j2gort.InstanceOf[*FlowSensitiveInlineVariables_SideEffectPredicate](v_)
	})
}
func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) SetOuter_(o *FlowSensitiveInlineVariables) {
	this_.outer_ = o
}
func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) GetOuter_() *FlowSensitiveInlineVariables {
	return this_.outer_
}
func (*FlowSensitiveInlineVariables_SideEffectPredicate) Implements_j2gort_Ext_com_google_common_base_Predicate() {
}
func (*FlowSensitiveInlineVariables_SideEffectPredicate) Implements_j2gort_Ext_java_util_function_Predicate() {
}

func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) initFields_FlowSensitiveInlineVariables_SideEffectPredicate() {
	_ = this_
}

// Java constructor: SideEffectPredicate()
func NewFlowSensitiveInlineVariables_SideEffectPredicate(outer_ *FlowSensitiveInlineVariables) *FlowSensitiveInlineVariables_SideEffectPredicate {
	this_ := &FlowSensitiveInlineVariables_SideEffectPredicate{}
	this_.outer_ = outer_
	InitFlowSensitiveInlineVariables_SideEffectPredicate(this_)
	return this_
}

func InitFlowSensitiveInlineVariables_SideEffectPredicate(this_ *FlowSensitiveInlineVariables_SideEffectPredicate) {
	this_.initFields_FlowSensitiveInlineVariables_SideEffectPredicate()
	this_.namesToCheck = nil
}

// Java constructor: SideEffectPredicate(java.util.Set<java.lang.String> names)
func NewFlowSensitiveInlineVariables_SideEffectPredicate_1(outer_ *FlowSensitiveInlineVariables, names j2gort.Ext_java_util_Set) *FlowSensitiveInlineVariables_SideEffectPredicate {
	this_ := &FlowSensitiveInlineVariables_SideEffectPredicate{}
	this_.outer_ = outer_
	InitFlowSensitiveInlineVariables_SideEffectPredicate_1(this_, names)
	return this_
}

func InitFlowSensitiveInlineVariables_SideEffectPredicate_1(this_ *FlowSensitiveInlineVariables_SideEffectPredicate, names j2gort.Ext_java_util_Set) {
	_ = names
	this_.initFields_FlowSensitiveInlineVariables_SideEffectPredicate()
	this_.namesToCheck = names
}

// Java method: SideEffectPredicate.apply(com.google.javascript.rhino.Node n)
func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) TestNode_(rawNode_1 uint32) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	n := pkg_rhino.Node(rawNode_1)
	_ = n
	if n == 0 {
		return false
	}
	if (((this_.namesToCheck != nil) && n.IsName()) && j2gort.SetContainsStr(this_.namesToCheck, n.GetString())) && NodeUtil_IsLValue(n) {
		return !(FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n))
	}
	var astAnalyzer *AstAnalyzer = this_.outer_.compiler.GetAstAnalyzer()
	_ = astAnalyzer
	if ((n.IsCall() || n.IsOptChainCall()) || n.IsTaggedTemplateLit()) && astAnalyzer.FunctionCallHasSideEffects(n) {
		return true
	}
	if n.IsNew() && astAnalyzer.ConstructorCallHasSideEffects(n) {
		return true
	}
	if n.IsDelProp() {
		return true
	}
	if (n.IsGetProp() || n.IsGetElem()) && NodeUtil_IsLValue(n) {
		return ((this_.namesToCheck == nil) || !(FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n)))
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if !(ControlFlowGraph_IsEnteringNewCfgNode(c)) && this_.Apply(j2gort.BoxNode(c)) {
				return true
			}
		}
	}
	return false
}

func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) Apply(arg_ any) bool {
	if this_ == nil {
		j2gort.ThrowNPE()
	}
	return this_.TestNode_(j2gort.CastNode[uint32](arg_))
}

func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_, p0_)
}

func (this_ *FlowSensitiveInlineVariables_SideEffectPredicate) Test(p0_ any) bool {
	return this_.Apply(p0_)
}

// Java class: com.google.javascript.jscomp.FlowSensitiveInlineVariables.GatherCandidatesCfgNodeCallback
type FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	outer_  *FlowSensitiveInlineVariables
	cfgNode pkg_rhino.Node `gson:"cfgNode"`
}

func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) IsJavaNil_() bool {
	return this_ == nil
}
func (*FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FlowSensitiveInlineVariables$GatherCandidatesCfgNodeCallback", "GatherCandidatesCfgNodeCallback"
}
func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase.IdentityHash_()
}
func (*FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FlowSensitiveInlineVariables.GatherCandidatesCfgNodeCallback", "com.google.javascript.jscomp.FlowSensitiveInlineVariables$GatherCandidatesCfgNodeCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback](v_)
	})
}
func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) SetOuter_(o *FlowSensitiveInlineVariables) {
	this_.outer_ = o
}
func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) GetOuter_() *FlowSensitiveInlineVariables {
	return this_.outer_
}
func (*FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) Implements_jscomp_ControlFlowGraph_AbstractCfgNodeTraversalCallback() {
}
func (*FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) initFields_FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback() {
	_ = this_

}

// Java constructor: GatherCandidatesCfgNodeCallback()
func NewFlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback(outer_ *FlowSensitiveInlineVariables) *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback {
	this_ := &FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback{}
	this_.outer_ = outer_
	InitFlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback(this_)
	return this_
}

func InitFlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback(this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) {
	InitControlFlowGraph_AbstractCfgNodeTraversalCallback(&this_.ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase)
	this_.initFields_FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback()
}

// Java method: GatherCandidatesCfgNodeCallback.setCfgNode(com.google.javascript.rhino.Node cfgNode)
func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) SetCfgNode(cfgNode pkg_rhino.Node) {
	_ = cfgNode
	this_.cfgNode = cfgNode
}

// Java method: GatherCandidatesCfgNodeCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		if parent == 0 {
			return
		}
		if ((((((NodeUtil_IsAssignmentOp(parent) && (parent.GetFirstChild() == n)) || NodeUtil_IsNameDeclaration(parent)) || parent.IsInc()) || parent.IsDec()) || parent.IsParamList()) || parent.IsCatch()) || NodeUtil_IsLhsByDestructuring(n) {
			return
		}
		var name string = n.GetString()
		_ = name
		if this_.outer_.compiler.GetCodingConvention().IsExported_2(name, true) {
			return
		}
		var def *MustBeReachingVariableDef_Definition = this_.outer_.reachingDef.GetDef(name, this_.cfgNode)
		_ = def
		if (def != nil) && !(this_.outer_.reachingDef.DependsOnOuterScopeVars(def)) {
			if this_.outer_.candidates == nil {
				this_.outer_.candidates = j2gort.NewExt_java_util_LinkedHashSet()
			}
			this_.outer_.candidates.Add_1(NewFlowSensitiveInlineVariables_Candidate(this_.outer_, name, def, n, this_.cfgNode))
		}
	}
}

// Java class: com.google.javascript.jscomp.FlowSensitiveInlineVariables.GatherCandidates
type FlowSensitiveInlineVariables_GatherCandidates struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer_   *FlowSensitiveInlineVariables
	gatherCb *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback `gson:"gatherCb"`
}

func (this_ *FlowSensitiveInlineVariables_GatherCandidates) BindSelf_(s any) {
	_ = s
}
func (this_ *FlowSensitiveInlineVariables_GatherCandidates) IsJavaNil_() bool { return this_ == nil }
func (*FlowSensitiveInlineVariables_GatherCandidates) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FlowSensitiveInlineVariables$GatherCandidates", "GatherCandidates"
}
func (this_ *FlowSensitiveInlineVariables_GatherCandidates) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowCallbackBase.IdentityHash_()
}
func (*FlowSensitiveInlineVariables_GatherCandidates) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FlowSensitiveInlineVariables.GatherCandidates", "com.google.javascript.jscomp.FlowSensitiveInlineVariables$GatherCandidates", func(v_ any) bool {
		return j2gort.InstanceOf[*FlowSensitiveInlineVariables_GatherCandidates](v_)
	})
}
func (this_ *FlowSensitiveInlineVariables_GatherCandidates) SetOuter_(o *FlowSensitiveInlineVariables) {
	this_.outer_ = o
}
func (this_ *FlowSensitiveInlineVariables_GatherCandidates) GetOuter_() *FlowSensitiveInlineVariables {
	return this_.outer_
}
func (*FlowSensitiveInlineVariables_GatherCandidates) Implements_jscomp_NodeTraversal_AbstractShallowCallback() {
}
func (*FlowSensitiveInlineVariables_GatherCandidates) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *FlowSensitiveInlineVariables_GatherCandidates) initFields_FlowSensitiveInlineVariables_GatherCandidates() {
	_ = this_
	this_.gatherCb = NewFlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback(this_.outer_)
}

// Java constructor: GatherCandidates()
func NewFlowSensitiveInlineVariables_GatherCandidates(outer_ *FlowSensitiveInlineVariables) *FlowSensitiveInlineVariables_GatherCandidates {
	this_ := &FlowSensitiveInlineVariables_GatherCandidates{}
	this_.outer_ = outer_
	InitFlowSensitiveInlineVariables_GatherCandidates(this_)
	return this_
}

func InitFlowSensitiveInlineVariables_GatherCandidates(this_ *FlowSensitiveInlineVariables_GatherCandidates) {
	InitNodeTraversal_AbstractShallowCallback(&this_.NodeTraversal_AbstractShallowCallbackBase)
	this_.initFields_FlowSensitiveInlineVariables_GatherCandidates()
}

// Java method: GatherCandidates.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FlowSensitiveInlineVariables_GatherCandidates) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var graphNode pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = pkg_graph.ToIface_DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch](j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.outer_.cfg.GetNode(n)))
	_ = graphNode
	if graphNode == nil || graphNode.IsJavaNil_() {
		return
	}
	var cfgNode pkg_rhino.Node = n
	_ = cfgNode
	this_.gatherCb.SetCfgNode(cfgNode)
	NodeTraversal_Traverse(this_.outer_.compiler, cfgNode, ToIface_NodeTraversal_Callback(this_.gatherCb))
}

// Java class: com.google.javascript.jscomp.FlowSensitiveInlineVariables.Candidate
type FlowSensitiveInlineVariables_Candidate struct {
	outer_               *FlowSensitiveInlineVariables
	varName              string                                `gson:"varName"`
	defMetadata          *MustBeReachingVariableDef_Definition `gson:"defMetadata"`
	def                  pkg_rhino.Node                        `gson:"def"`
	use                  pkg_rhino.Node                        `gson:"use"`
	useCfgNode           pkg_rhino.Node                        `gson:"useCfgNode"`
	numUsesWithinCfgNode int32                                 `gson:"numUsesWithinCfgNode"`
	idHash_              atomic.Int32
}

func (this_ *FlowSensitiveInlineVariables_Candidate) BindSelf_(s any) {
	_ = s
}
func (this_ *FlowSensitiveInlineVariables_Candidate) IsJavaNil_() bool { return this_ == nil }
func (*FlowSensitiveInlineVariables_Candidate) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FlowSensitiveInlineVariables$Candidate", "Candidate"
}
func (this_ *FlowSensitiveInlineVariables_Candidate) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FlowSensitiveInlineVariables_Candidate) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FlowSensitiveInlineVariables.Candidate", "com.google.javascript.jscomp.FlowSensitiveInlineVariables$Candidate", func(v_ any) bool {
		return j2gort.InstanceOf[*FlowSensitiveInlineVariables_Candidate](v_)
	})
}
func (this_ *FlowSensitiveInlineVariables_Candidate) SetOuter_(o *FlowSensitiveInlineVariables) {
	this_.outer_ = o
}
func (this_ *FlowSensitiveInlineVariables_Candidate) GetOuter_() *FlowSensitiveInlineVariables {
	return this_.outer_
}

func (this_ *FlowSensitiveInlineVariables_Candidate) initFields_FlowSensitiveInlineVariables_Candidate() {
	_ = this_
	this_.varName = j2gort.NullString
}

// Java constructor: Candidate(java.lang.String varName, com.google.javascript.jscomp.MustBeReachingVariableDef.Definition defMetadata, com.google.javascript.rhino.Node use, com.google.javascript.rhino.Node useCfgNode)
func NewFlowSensitiveInlineVariables_Candidate(outer_ *FlowSensitiveInlineVariables, varName string, defMetadata *MustBeReachingVariableDef_Definition, use pkg_rhino.Node, useCfgNode pkg_rhino.Node) *FlowSensitiveInlineVariables_Candidate {
	this_ := &FlowSensitiveInlineVariables_Candidate{}
	this_.outer_ = outer_
	InitFlowSensitiveInlineVariables_Candidate(this_, varName, defMetadata, use, useCfgNode)
	return this_
}

func InitFlowSensitiveInlineVariables_Candidate(this_ *FlowSensitiveInlineVariables_Candidate, varName string, defMetadata *MustBeReachingVariableDef_Definition, use pkg_rhino.Node, useCfgNode pkg_rhino.Node) {
	_ = varName
	_ = defMetadata
	_ = use
	_ = useCfgNode
	this_.initFields_FlowSensitiveInlineVariables_Candidate()
	if !(use.IsName()) {
		j2gort.CheckArgument(false)
	}
	this_.varName = varName
	this_.defMetadata = defMetadata
	this_.use = use
	this_.useCfgNode = useCfgNode
}

// Java method: Candidate.getDefCfgNode()
func (this_ *FlowSensitiveInlineVariables_Candidate) priv_Candidate_getDefCfgNode_2867() pkg_rhino.Node {
	return this_.defMetadata.node
}

// Java method: Candidate.canInline(com.google.javascript.jscomp.Scope scope)
func (this_ *FlowSensitiveInlineVariables_Candidate) priv_Candidate_canInline_2868(scope *Scope) bool {
	_ = scope
	if this_.priv_Candidate_getDefCfgNode_2867().IsFunction() {
		return false
	}
	this_.priv_Candidate_getDefinition_2871(this_.priv_Candidate_getDefCfgNode_2867())
	this_.priv_Candidate_getNumUseInUseCfgNode_2872(this_.useCfgNode)
	if this_.def == 0 {
		return false
	}
	if this_.def.IsAssign() && !(NodeUtil_IsExprAssign(this_.def.GetParent())) {
		return false
	}
	var namesToCheck j2gort.Ext_java_util_Set = nil
	_ = namesToCheck
	if this_.defMetadata.IsEmptyDepends() {
		namesToCheck = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	} else {
		namesToCheck = j2gort.NewExt_java_util_LinkedHashSet_1_int32(int32(4))
		this_.defMetadata.AddDependencyNamesTo(namesToCheck)
	}
	var val_2 FlowSensitiveInlineVariables_SideEffectPredicate
	val_2.outer_ = this_.outer_
	InitFlowSensitiveInlineVariables_SideEffectPredicate_1(&val_2, namesToCheck)
	var sideEffectPredicateWithNames *FlowSensitiveInlineVariables_SideEffectPredicate = &val_2
	_ = sideEffectPredicateWithNames
	if FlowSensitiveInlineVariables_CheckPostExpressions(this_.def, this_.priv_Candidate_getDefCfgNode_2867(), sideEffectPredicateWithNames) {
		return false
	}
	if FlowSensitiveInlineVariables_CheckPreExpressions(this_.use, this_.useCfgNode, sideEffectPredicateWithNames) {
		return false
	}
	if this_.outer_.compiler.GetAstAnalyzer().MayHaveSideEffects(this_.def.GetLastChild()) {
		return false
	}
	if this_.numUsesWithinCfgNode != int32(1) {
		return false
	}
	if NodeUtil_IsWithinLoop(this_.use) {
		return false
	}
	if !(this_.outer_.reachingUses.HasSingleUse(this_.varName, this_.priv_Candidate_getDefCfgNode_2867())) {
		return false
	}
	if !(this_.priv_Candidate_isRhsSafeToInline_2873(scope)) {
		return false
	}
	if NodeUtil_IsStatementBlock(this_.priv_Candidate_getDefCfgNode_2867().GetParent()) && (this_.priv_Candidate_getDefCfgNode_2867().GetNext() != this_.useCfgNode) {
		var pathCheck *pkg_graph.CheckPathsBetweenNodes[pkg_rhino.Node, ControlFlowGraph_Branch] = pkg_graph.NewCheckPathsBetweenNodes_6[pkg_rhino.Node, ControlFlowGraph_Branch](j2gort.AsIface[pkg_graph.DiGraph[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.outer_.cfg), pkg_graph.ToIface_DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch](j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.outer_.cfg.GetNode(this_.priv_Candidate_getDefCfgNode_2867()))), pkg_graph.ToIface_DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch](j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.outer_.cfg.GetNode(this_.useCfgNode))), j2gort.ToIface_Ext_com_google_common_base_Predicate(this_.outer_.sideEffectPredicate), j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue(), false)
		_ = pathCheck
		if pathCheck.SomePathsSatisfyPredicate() {
			return false
		}
	}
	return true
}

// Java method: Candidate.hasExactlyOne(java.lang.Iterable<com.google.javascript.rhino.Node> iterable)
func (this_ *FlowSensitiveInlineVariables_Candidate) HasExactlyOne(iterable j2gort.Ext_java_lang_Iterable) bool {
	_ = iterable
	var iterator j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](iterable.Iterator())
	_ = iterator
	if iterator.HasNext() {
		iterator.Next()
		if !(iterator.HasNext()) {
			return true
		}
	}
	return false
}

// Java method: Candidate.inlineVariable()
func (this_ *FlowSensitiveInlineVariables_Candidate) priv_Candidate_inlineVariable_2870() {
	var defParent pkg_rhino.Node = this_.def.GetParent()
	_ = defParent
	var useParent pkg_rhino.Node = this_.use.GetParent()
	_ = useParent
	if this_.def.IsAssign() {
		var rhs pkg_rhino.Node = this_.def.GetLastChild()
		_ = rhs
		rhs.Detach()
		if !(defParent.IsExprResult()) {
			j2gort.CheckState(false)
		}
		for defParent.GetParent().IsLabel() {
			defParent = defParent.GetParent()
		}
		this_.outer_.compiler.ReportChangeToEnclosingScope(defParent)
		defParent.Detach()
		this_.use.ReplaceWith(rhs)
	} else if NodeUtil_IsNameDeclaration(defParent) {
		var rhs_2 pkg_rhino.Node = this_.def.GetLastChild()
		_ = rhs_2
		if defParent.IsConst() {
			rhs_2.ReplaceWith(pkg_rhino.Node_NewString_2(pkg_rhino.StaticGet_Token_NAME(), "undefined"))
			this_.use.ReplaceWith(rhs_2)
		} else {
			rhs_2.Detach()
			this_.use.ReplaceWith(rhs_2)
		}
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "No other definitions can be inlined.")))
	}
	this_.outer_.compiler.ReportChangeToEnclosingScope(useParent)
}

// Java method: Candidate.getDefinition(com.google.javascript.rhino.Node n)
func (this_ *FlowSensitiveInlineVariables_Candidate) priv_Candidate_getDefinition_2871(n pkg_rhino.Node) {
	_ = n
	var gatherCb ControlFlowGraph_AbstractCfgNodeTraversalCallback = NewFlowSensitiveInlineVariables__Anon_43(this_)
	_ = gatherCb
	NodeTraversal_Traverse(this_.outer_.compiler, n, gatherCb)
}

// Java method: Candidate.getNumUseInUseCfgNode(com.google.javascript.rhino.Node cfgNode)
func (this_ *FlowSensitiveInlineVariables_Candidate) priv_Candidate_getNumUseInUseCfgNode_2872(cfgNode pkg_rhino.Node) {
	_ = cfgNode
	this_.numUsesWithinCfgNode = int32(0)
	var gatherCb ControlFlowGraph_AbstractCfgNodeTraversalCallback = NewFlowSensitiveInlineVariables__Anon_44(this_, cfgNode)
	_ = gatherCb
	NodeTraversal_Traverse(this_.outer_.compiler, cfgNode, gatherCb)
}

// Java method: Candidate.isRhsSafeToInline(com.google.javascript.jscomp.Scope usageScope)
func (this_ *FlowSensitiveInlineVariables_Candidate) priv_Candidate_isRhsSafeToInline_2873(usageScope *Scope) bool {
	_ = usageScope
	if NodeUtil_Has(this_.def.GetLastChild(), lambdaSingleton_flow_sensitive_inline_variables_1, lambdaSingleton_flow_sensitive_inline_variables_2) {
		return false
	}
	if nodeUtilHasPredAlwaysTrueFn_(this_.def.GetLastChild(), func(rawNode_5 uint32) bool {
		input := pkg_rhino.Node(rawNode_5)
		_ = input
		if input.IsName() {
			var name string = input.GetString()
			_ = name
			if !(j2gort.CheckNotNullString(name) == "") && !(usageScope.HasSlot(name)) {
				return true
			}
		}
		return false
	}) {
		return false
	}
	return true
}

// Java class:
type FlowSensitiveInlineVariables__Anon_43 struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	outer_ *FlowSensitiveInlineVariables_Candidate
}

func (this_ *FlowSensitiveInlineVariables__Anon_43) BindSelf_(s any) {
	_ = s
}
func (this_ *FlowSensitiveInlineVariables__Anon_43) IsJavaNil_() bool { return this_ == nil }
func (*FlowSensitiveInlineVariables__Anon_43) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FlowSensitiveInlineVariables$Candidate$1", ""
}
func (this_ *FlowSensitiveInlineVariables__Anon_43) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase.IdentityHash_()
}
func (*FlowSensitiveInlineVariables__Anon_43) HasIdentityEquals_() {}
func (this_ *FlowSensitiveInlineVariables__Anon_43) SetOuter_(o *FlowSensitiveInlineVariables_Candidate) {
	this_.outer_ = o
}
func (this_ *FlowSensitiveInlineVariables__Anon_43) GetOuter_() *FlowSensitiveInlineVariables_Candidate {
	return this_.outer_
}
func (*FlowSensitiveInlineVariables__Anon_43) Implements_jscomp_ControlFlowGraph_AbstractCfgNodeTraversalCallback() {
}
func (*FlowSensitiveInlineVariables__Anon_43) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *FlowSensitiveInlineVariables__Anon_43) initFields_FlowSensitiveInlineVariables__Anon_43() {
	_ = this_
}

func NewFlowSensitiveInlineVariables__Anon_43(outer_ *FlowSensitiveInlineVariables_Candidate) *FlowSensitiveInlineVariables__Anon_43 {
	this_ := &FlowSensitiveInlineVariables__Anon_43{}
	this_.outer_ = outer_
	InitControlFlowGraph_AbstractCfgNodeTraversalCallback(&this_.ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase)
	this_.initFields_FlowSensitiveInlineVariables__Anon_43()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FlowSensitiveInlineVariables__Anon_43) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		{
			var iifeCond_1 bool
			{

				sr_6 := n.GetString()
				se_7 := this_.outer_.varName
				iifeCond_1 = j2gort.CheckNotNullString(sr_6) == se_7
			}
			if iifeCond_1 && n.HasChildren() {
				this_.outer_.def = n
			}
		}
		return
	case int32(51) /* Token.ASSIGN */ :
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		if lhs.IsName() && func() bool {
			sr_8 := lhs.GetString()
			se_9 := this_.outer_.varName
			return j2gort.CheckNotNullString(sr_8) == se_9
		}() {
			this_.outer_.def = n
		}
		return
	default:
	}
}

// Java class:
type FlowSensitiveInlineVariables__Anon_44 struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	outer_        *FlowSensitiveInlineVariables_Candidate
	cap_0_cfgNode pkg_rhino.Node
}

func (this_ *FlowSensitiveInlineVariables__Anon_44) BindSelf_(s any) {
	_ = s
}
func (this_ *FlowSensitiveInlineVariables__Anon_44) IsJavaNil_() bool { return this_ == nil }
func (*FlowSensitiveInlineVariables__Anon_44) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.FlowSensitiveInlineVariables$Candidate$2", ""
}
func (this_ *FlowSensitiveInlineVariables__Anon_44) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase.IdentityHash_()
}
func (*FlowSensitiveInlineVariables__Anon_44) HasIdentityEquals_() {}
func (this_ *FlowSensitiveInlineVariables__Anon_44) SetOuter_(o *FlowSensitiveInlineVariables_Candidate) {
	this_.outer_ = o
}
func (this_ *FlowSensitiveInlineVariables__Anon_44) GetOuter_() *FlowSensitiveInlineVariables_Candidate {
	return this_.outer_
}
func (*FlowSensitiveInlineVariables__Anon_44) Implements_jscomp_ControlFlowGraph_AbstractCfgNodeTraversalCallback() {
}
func (*FlowSensitiveInlineVariables__Anon_44) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *FlowSensitiveInlineVariables__Anon_44) initFields_FlowSensitiveInlineVariables__Anon_44() {
	_ = this_
}

func NewFlowSensitiveInlineVariables__Anon_44(outer_ *FlowSensitiveInlineVariables_Candidate, cap_p_cfgNode pkg_rhino.Node) *FlowSensitiveInlineVariables__Anon_44 {
	this_ := &FlowSensitiveInlineVariables__Anon_44{}
	this_.outer_ = outer_
	this_.cap_0_cfgNode = cap_p_cfgNode
	InitControlFlowGraph_AbstractCfgNodeTraversalCallback(&this_.ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase)
	this_.initFields_FlowSensitiveInlineVariables__Anon_44()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *FlowSensitiveInlineVariables__Anon_44) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() && func() bool {
		sr_10 := n.GetString()
		se_11 := this_.outer_.varName
		return j2gort.CheckNotNullString(sr_10) == se_11
	}() {
		if (parent.IsAssign() && (parent.GetFirstChild() == n)) && this_.priv_ident__isAssignChain_2874(parent, this_.cap_0_cfgNode) {
			return
		} else {
			this_.outer_.numUsesWithinCfgNode = this_.outer_.numUsesWithinCfgNode + 1
		}
	}
}

// Java method: .isAssignChain(com.google.javascript.rhino.Node child, com.google.javascript.rhino.Node ancestor)
func (this_ *FlowSensitiveInlineVariables__Anon_44) priv_ident__isAssignChain_2874(child pkg_rhino.Node, ancestor pkg_rhino.Node) bool {
	_ = child
	_ = ancestor
	{
		var n pkg_rhino.Node = child
		_ = n
		for ; n != ancestor; n = n.GetParent() {
			if !(n.IsAssign()) {
				return false
			}
		}
	}
	return true
}
func (*FlowSensitiveInlineVariables) NoOpVisit_()                                                 {}
func (*FlowSensitiveInlineVariables) CanReleaseSyntacticScopes_()                                 {}
func (*FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) CanReleaseSyntacticScopes_() {}
func (*FlowSensitiveInlineVariables_GatherCandidates) CanReleaseSyntacticScopes_()                {}
func (*FlowSensitiveInlineVariables__Anon_43) CanReleaseSyntacticScopes_()                        {}
func (*FlowSensitiveInlineVariables__Anon_44) CanReleaseSyntacticScopes_()                        {}
