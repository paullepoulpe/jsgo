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
// Port of com/google/javascript/jscomp/OptimizeReturns.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.OptimizeReturns
type OptimizeReturns struct {
	compiler     AbstractCompiler       `gson:"compiler"`
	decisionsLog pkg_diagnostic.LogFile `gson:"decisionsLog"`
	idHash_      atomic.Int32
}

func (this_ *OptimizeReturns) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeReturns) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeReturns) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeReturns", "OptimizeReturns"
}
func (this_ *OptimizeReturns) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeReturns) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeReturns", "com.google.javascript.jscomp.OptimizeReturns", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeReturns](v_)
	})
}
func (*OptimizeReturns) Implements_jscomp_OptimizeCalls_CallGraphCompilerPass() {}
func (*OptimizeReturns) Implements_jscomp_CompilerPass()                        {}

func (this_ *OptimizeReturns) initFields_OptimizeReturns() {
	_ = this_
}

// Java constructor: OptimizeReturns(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewOptimizeReturns(compiler AbstractCompiler) *OptimizeReturns {
	this_ := &OptimizeReturns{}
	InitOptimizeReturns(this_, compiler)
	return this_
}

func InitOptimizeReturns(this_ *OptimizeReturns, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_OptimizeReturns()
	this_.compiler = compiler
}

// Java method: OptimizeReturns.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *OptimizeReturns) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	OptimizeCalls_S_Builder().SetCompiler(this_.compiler).SetConsiderExterns(false).AddPass(this_).Build().Process_2(externs, root)
}

// Java method: OptimizeReturns.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root, com.google.javascript.jscomp.OptimizeCalls.ReferenceMap definitions)
func (this_ *OptimizeReturns) Process_3(externs pkg_rhino.Node, root pkg_rhino.Node, definitions *OptimizeCalls_ReferenceMap) {
	_ = externs
	_ = root
	_ = definitions
	func() {
		defer func() {
			this_.decisionsLog = nil
		}()
		var logFile pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenIndexedLog(j2gort.GetObjectClass(this_), "decisions.log", []string{})
		_ = logFile
		defer j2gort.CloseResource(logFile)
		this_.decisionsLog = logFile
		var toOptimize []*OptimizeCalls_NodeList = []*OptimizeCalls_NodeList{}
		_ = toOptimize
		{
			cur_1 := definitions.names.IterEntries()
			for cur_1.NextEntry() {
				var key string = cur_1.Key()
				_ = key
				var refs *OptimizeCalls_NodeList = cur_1.Value()
				_ = refs
				if this_.priv_OptimizeReturns_isCandidate_4119(key, refs) {
					if this_.decisionsLog.IsLogging() {
						this_.decisionsLog.Log_2("name %s\tremoving return value", []any{j2gort.BoxString(key)})
					}
					toOptimize = append(toOptimize, refs)
				}
			}
		}
		{
			cur_2 := definitions.props.IterEntries()
			for cur_2.NextEntry() {
				var key_2 string = cur_2.Key()
				_ = key_2
				var refs_2 *OptimizeCalls_NodeList = cur_2.Value()
				_ = refs_2
				if this_.priv_OptimizeReturns_isCandidate_4119(key_2, refs_2) {
					if this_.decisionsLog.IsLogging() {
						this_.decisionsLog.Log_2("property %s\tremoving return value", []any{j2gort.BoxString(key_2)})
					}
					toOptimize = append(toOptimize, refs_2)
				}
			}
		}
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(toOptimize))
			_ = size
			for ; i < size; i = i + 1 {
				var refs_3 *OptimizeCalls_NodeList = toOptimize[i]
				_ = refs_3
				var fns *OptimizeCalls_NodeList = allocPooledVia_OptimizeCalls_ReferenceMap_GetFunctionNodeList(&definitions.slab_OptimizeCalls_NodeList, refs_3)
				_ = fns
				{
					var j int32 = int32(0)
					_ = j
					var fnCount int32 = fns.Size()
					_ = fnCount
					for ; j < fnCount; j = j + 1 {
						this_.priv_OptimizeReturns_rewriteReturns_4122(fns.Get(j))
					}
				}
			}
		}
	}()
}

// Java method: OptimizeReturns.isCandidate(java.lang.String name, com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeReturns) priv_OptimizeReturns_isCandidate_4119(name string, refs *OptimizeCalls_NodeList) bool {
	_ = name
	_ = refs
	if !(OptimizeCalls_MayBeOptimizableName(this_.compiler, name)) {
		if this_.decisionsLog.IsLogging() {
			this_.decisionsLog.Log_2("%s\tnot an optimizable name", []any{j2gort.BoxString(name)})
		}
		return false
	}
	var seenCandidateDefiniton bool = false
	_ = seenCandidateDefiniton
	var seenUse bool = false
	_ = seenUse
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = refs.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var n pkg_rhino.Node = refs.Get(i)
			_ = n
			if OptimizeCalls_ReferenceMap_IsCallTarget(n) || OptimizeCalls_ReferenceMap_IsOptChainCallTarget(n) {
				var callNode pkg_rhino.Node = OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n)
				_ = callNode
				if NodeUtil_IsExpressionResultUsed(callNode) {
					if this_.decisionsLog.IsLogging() {
						this_.decisionsLog.Log_2("%s\treturn value used: %s", []any{j2gort.BoxString(name), j2gort.BoxString(callNode.GetLocation())})
					}
					return false
				}
				seenUse = true
			} else if this_.priv_OptimizeReturns_isCandidateDefinition_4120(n) {
				seenCandidateDefiniton = true
			} else {
				if !(OptimizeCalls_IsAllowedReference(n)) {
					if this_.decisionsLog.IsLogging() {
						this_.decisionsLog.Log_2("%s\tdisallowed reference: %s", []any{j2gort.BoxString(name), j2gort.BoxString(n.GetLocation())})
					}
					return false
				}
			}
		}
	}
	if !(seenUse) {
		if this_.decisionsLog.IsLogging() {
			this_.decisionsLog.Log_2("%s\tno usage seen", []any{j2gort.BoxString(name)})
		}
		return false
	}
	if !(seenCandidateDefiniton) {
		if this_.decisionsLog.IsLogging() {
			this_.decisionsLog.Log_2("%s\tno definition seen", []any{j2gort.BoxString(name)})
		}
		return false
	}
	return true
}

// Java method: OptimizeReturns.isCandidateDefinition(com.google.javascript.rhino.Node n)
func (this_ *OptimizeReturns) priv_OptimizeReturns_isCandidateDefinition_4120(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
		return true
	} else if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		if OptimizeReturns_IsCandidateFunction(parent.GetLastChild()) {
			return true
		}
	} else if n.IsName() {
		if n.HasChildren() && OptimizeReturns_IsCandidateFunction(n.GetFirstChild()) {
			return true
		}
	} else if this_.priv_OptimizeReturns_isClassMemberDefinition_4121(n) {
		return true
	}
	return false
}

// Java method: OptimizeReturns.isClassMemberDefinition(com.google.javascript.rhino.Node n)
func (this_ *OptimizeReturns) priv_OptimizeReturns_isClassMemberDefinition_4121(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsMemberFunctionDef() && n.GetParent().IsClassMembers())
}

// Java static method: OptimizeReturns.isCandidateFunction(com.google.javascript.rhino.Node n)
func OptimizeReturns_IsCandidateFunction(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		return !(NodeUtil_IsNamedFunctionExpression(n))
	case int32(50) /* Token.COMMA */, int32(101) /* Token.CAST */ :
		return OptimizeReturns_IsCandidateFunction(n.GetLastChild())
	case int32(67) /* Token.HOOK */ :
		return (OptimizeReturns_IsCandidateFunction(n.GetSecondChild()) && OptimizeReturns_IsCandidateFunction(n.GetLastChild()))
	case int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(70) /* Token.COALESCE */ :
		return (OptimizeReturns_IsCandidateFunction(n.GetFirstChild()) && OptimizeReturns_IsCandidateFunction(n.GetLastChild()))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: OptimizeReturns.rewriteReturns(com.google.javascript.rhino.Node fnNode)
func (this_ *OptimizeReturns) priv_OptimizeReturns_rewriteReturns_4122(fnNode pkg_rhino.Node) {
	_ = fnNode
	if !(fnNode.IsFunction()) {
		j2gort.CheckState(false)
	}
	var body pkg_rhino.Node = fnNode.GetLastChild()
	_ = body
	NodeUtil_VisitPostOrder_3(body, NewOptimizeReturns__Anon_68(this_, body), NewNodeUtil_MatchShallowStatement())
}

// Java method: OptimizeReturns.isRemovableValue(com.google.javascript.rhino.Node n)
func (this_ *OptimizeReturns) priv_OptimizeReturns_isRemovableValue_4123(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(131) /* Token.TEMPLATELIT */, int32(46) /* Token.ARRAYLIT */ :
		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if (!(child.IsEmpty())) && !(this_.priv_OptimizeReturns_isRemovableValue_4123(child)) {
					return false
				}
			}
		}
		return true
	case int32(42) /* Token.REGEXP */, int32(35) /* Token.STRINGLIT */, int32(33) /* Token.NUMBER */, int32(36) /* Token.NULL */, int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */, int32(133) /* Token.TEMPLATELIT_STRING */ :
		return true
	case int32(132) /* Token.TEMPLATELIT_SUB */, int32(101) /* Token.CAST */, int32(19) /* Token.NOT */, int32(87) /* Token.VOID */, int32(22) /* Token.NEG */ :
		return this_.priv_OptimizeReturns_isRemovableValue_4123(n.GetFirstChild())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java class:
type OptimizeReturns__Anon_68 struct {
	outer_     *OptimizeReturns
	cap_0_body pkg_rhino.Node
	idHash_    atomic.Int32
}

func (this_ *OptimizeReturns__Anon_68) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeReturns__Anon_68) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeReturns__Anon_68) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeReturns$1", ""
}
func (this_ *OptimizeReturns__Anon_68) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeReturns__Anon_68) HasIdentityEquals_()                 {}
func (this_ *OptimizeReturns__Anon_68) SetOuter_(o *OptimizeReturns)  { this_.outer_ = o }
func (this_ *OptimizeReturns__Anon_68) GetOuter_() *OptimizeReturns   { return this_.outer_ }
func (*OptimizeReturns__Anon_68) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *OptimizeReturns__Anon_68) initFields_OptimizeReturns__Anon_68() {
	_ = this_
}

func NewOptimizeReturns__Anon_68(outer_ *OptimizeReturns, cap_p_body pkg_rhino.Node) *OptimizeReturns__Anon_68 {
	this_ := &OptimizeReturns__Anon_68{}
	this_.outer_ = outer_
	this_.cap_0_body = cap_p_body
	this_.initFields_OptimizeReturns__Anon_68()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node n)
func (this_ *OptimizeReturns__Anon_68) Visit(n pkg_rhino.Node) {
	_ = n
	if n.IsReturn() && n.HasOneChild() {
		var result pkg_rhino.Node = n.GetFirstChild()
		_ = result
		var keepValue bool = !(this_.outer_.priv_OptimizeReturns_isRemovableValue_4123(result))
		_ = keepValue
		result.Detach()
		if keepValue {
			pkg_rhino.IR_ExprResult(result).Srcref(result).InsertBefore(n)
		} else {
			NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(result, this_.outer_.compiler)
		}
		this_.outer_.compiler.ReportChangeToEnclosingScope(this_.cap_0_body)
	}
}
