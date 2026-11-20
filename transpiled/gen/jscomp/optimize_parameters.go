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
// Port of com/google/javascript/jscomp/OptimizeParameters.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.OptimizeParameters
type OptimizeParameters struct {
	compiler                    AbstractCompiler               `gson:"compiler"`
	astAnalyzer                 *AstAnalyzer                   `gson:"astAnalyzer"`
	globalScope                 *Scope                         `gson:"globalScope"`
	decisionsLog                pkg_diagnostic.LogFile         `gson:"decisionsLog"`
	scratchCandidateDefinitions j2gort.Ext_java_util_ArrayList `gson:"scratchCandidateDefinitions"`
	idHash_                     atomic.Int32
	slab_CandidateAnalysis      j2gort.PooledSlab[OptimizeParameters_CandidateAnalysis]
	slab_ParameterList          j2gort.PooledSlab[OptimizeParameters_ParameterList]
	slab_Parameter              j2gort.PooledSlab[OptimizeParameters_Parameter]
	slab_OptimizeCalls_NodeList j2gort.PooledSlab[OptimizeCalls_NodeList]
}

func (this_ *OptimizeParameters) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeParameters) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeParameters) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeParameters", "OptimizeParameters"
}
func (this_ *OptimizeParameters) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeParameters) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeParameters", "com.google.javascript.jscomp.OptimizeParameters", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeParameters](v_)
	})
}
func (*OptimizeParameters) Implements_jscomp_CompilerPass()                        {}
func (*OptimizeParameters) Implements_jscomp_OptimizeCalls_CallGraphCompilerPass() {}

func (this_ *OptimizeParameters) initFields_OptimizeParameters() {
	_ = this_
	this_.scratchCandidateDefinitions = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: OptimizeParameters(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewOptimizeParameters(compiler AbstractCompiler) *OptimizeParameters {
	this_ := &OptimizeParameters{}
	InitOptimizeParameters(this_, compiler)
	return this_
}

func InitOptimizeParameters(this_ *OptimizeParameters, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_OptimizeParameters()
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.astAnalyzer = compiler.GetAstAnalyzer()
}

// Java method: OptimizeParameters.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *OptimizeParameters) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(this_.compiler.GetLifeCycleStage() == StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED()) {
		j2gort.CheckState(false)
	}
	OptimizeCalls_S_Builder().SetCompiler(this_.compiler).SetConsiderExterns(false).AddPass(this_).Build().Process_2(externs, root)
}

// Java method: OptimizeParameters.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root, com.google.javascript.jscomp.OptimizeCalls.ReferenceMap refMap)
func (this_ *OptimizeParameters) Process_3(externs pkg_rhino.Node, root pkg_rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	_ = externs
	_ = root
	_ = refMap
	defer func() {
		this_.globalScope = nil
		this_.ReleaseSlabs_()
	}()

	func() {
		defer func() {
			this_.decisionsLog = nil
		}()
		var decisionsLog pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenIndexedLog(j2gort.GetObjectClass(this_), "decisions.log", []string{})
		_ = decisionsLog
		defer j2gort.CloseResource(decisionsLog)
		this_.decisionsLog = decisionsLog
		this_.globalScope = refMap.GetGlobalScope()
		var toOptimize []*OptimizeCalls_NodeList = []*OptimizeCalls_NodeList{}
		_ = toOptimize
		{
			cur_1 := refMap.names.IterEntries()
			for cur_1.NextEntry() {
				var key string = cur_1.Key()
				_ = key
				var refs *OptimizeCalls_NodeList = cur_1.Value()
				_ = refs
				var candidateAnalysis *OptimizeParameters_CandidateAnalysis = this_.priv_OptimizeParameters_analyzeCandidateName_4097(key, refs)
				_ = candidateAnalysis
				if candidateAnalysis.ShouldConvertTaggedTemplateLiterals() {
					candidateAnalysis.ConvertTaggedTemplateLiterals()
				}
				if candidateAnalysis.IsSafeToOptimize() {
					toOptimize = append(toOptimize, refs)
				}
				this_.slab_CandidateAnalysis.Reset()
			}
		}
		{
			cur_2 := refMap.props.IterEntries()
			for cur_2.NextEntry() {
				var key_2 string = cur_2.Key()
				_ = key_2
				var refs_2 *OptimizeCalls_NodeList = cur_2.Value()
				_ = refs_2
				var candidateAnalysis_2 *OptimizeParameters_CandidateAnalysis = this_.priv_OptimizeParameters_analyzeCandidateProperty_4098(key_2, refs_2)
				_ = candidateAnalysis_2
				if candidateAnalysis_2.ShouldConvertTaggedTemplateLiterals() {
					candidateAnalysis_2.ConvertTaggedTemplateLiterals()
				}
				if candidateAnalysis_2.IsSafeToOptimize() {
					toOptimize = append(toOptimize, refs_2)
				}
				this_.slab_CandidateAnalysis.Reset()
			}
		}
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(toOptimize))
			_ = size
			for ; i < size; i = i + 1 {
				this_.priv_OptimizeParameters_tryEliminateOptionalArgs_4105(toOptimize[i])
			}
		}
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var size_2 int32 = int32(len(toOptimize))
			_ = size_2
			for ; i_2 < size_2; i_2 = i_2 + 1 {
				this_.priv_OptimizeParameters_tryEliminateConstantArgs_4106(toOptimize[i_2])
			}
		}
		var val_3 OptimizeParameters_UnusedParameterOptimizer
		val_3.outer_ = this_
		InitOptimizeParameters_UnusedParameterOptimizer(&val_3)
		var optimizer *OptimizeParameters_UnusedParameterOptimizer = &val_3
		_ = optimizer
		{
			var i_3 int32 = int32(0)
			_ = i_3
			var size_3 int32 = int32(len(toOptimize))
			_ = size_3
			for ; i_3 < size_3; i_3 = i_3 + 1 {
				optimizer.TryEliminateUnusedArgs(toOptimize[i_3])
			}
		}
		optimizer.ApplyChanges()
	}()
}

// Java static method: OptimizeParameters.alreadyRemoved(com.google.javascript.rhino.Node n)
func OptimizeParameters_AlreadyRemoved(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	for parent != 0 {
		n = parent
		parent = n.GetParent()
	}
	return !(n.IsRoot())
}

// Java method: OptimizeParameters.analyzeCandidateName(java.lang.String key, com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeParameters) priv_OptimizeParameters_analyzeCandidateName_4097(key string, refs *OptimizeCalls_NodeList) *OptimizeParameters_CandidateAnalysis {
	_ = key
	_ = refs
	return this_.priv_OptimizeParameters_analyzeCandidate_4099("name", key, refs)
}

// Java method: OptimizeParameters.analyzeCandidateProperty(java.lang.String key, com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeParameters) priv_OptimizeParameters_analyzeCandidateProperty_4098(key string, refs *OptimizeCalls_NodeList) *OptimizeParameters_CandidateAnalysis {
	_ = key
	_ = refs
	return this_.priv_OptimizeParameters_analyzeCandidate_4099("property", key, refs)
}

// Java method: OptimizeParameters.analyzeCandidate(java.lang.String refKind, java.lang.String key, com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeParameters) priv_OptimizeParameters_analyzeCandidate_4099(refKind string, key string, refs *OptimizeCalls_NodeList) *OptimizeParameters_CandidateAnalysis {
	_ = refKind
	_ = key
	_ = refs
	defer this_.slab_OptimizeCalls_NodeList.Reset()
	val_4 := this_.slab_CandidateAnalysis.Alloc()

	val_4.outer_ = this_
	InitOptimizeParameters_CandidateAnalysis(val_4)
	var analysis *OptimizeParameters_CandidateAnalysis = val_4
	_ = analysis
	if !(OptimizeCalls_MayBeOptimizableName(this_.compiler, key)) {
		if this_.decisionsLog.IsLogging() {
			this_.decisionsLog.Log_2("%s\t%s\tnot an optimizable name", []any{j2gort.BoxString(refKind), j2gort.BoxString(key)})
		}
		return analysis
	}
	var definitions j2gort.Ext_java_util_ArrayList = this_.scratchCandidateDefinitions
	_ = definitions
	definitions.Clear()
	var seenCandidateUse bool = false
	_ = seenCandidateUse
	{
		var i int32 = int32(0)
		_ = i
		var refCount int32 = refs.Size()
		_ = refCount
		for ; i < refCount; i = i + 1 {
			var n pkg_rhino.Node = refs.Get(i)
			_ = n
			if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
				seenCandidateUse = true
			} else if n.GetParent().IsTaggedTemplateLit() {
				analysis.AddTaggedTemplateLiteral(n.GetParent())
			} else if this_.priv_OptimizeParameters_isCandidateDefinition_4101(n) {
				definitions.Add_1(j2gort.BoxNode(n))
			} else if !(OptimizeCalls_IsAllowedReference(n)) {
				if this_.decisionsLog.IsLogging() {
					this_.decisionsLog.Log_2("%s\t%s\tnot an allowed reference: %s", []any{j2gort.BoxString(refKind), j2gort.BoxString(key), j2gort.BoxString(n.GetLocation())})
				}
				return analysis
			}
		}
	}
	if definitions.IsEmpty() {
		if this_.decisionsLog.IsLogging() {
			this_.decisionsLog.Log_2("%s\t%s\tno definition found", []any{j2gort.BoxString(refKind), j2gort.BoxString(key)})
		}
		return analysis
	}
	if analysis.HasTaggedTemplateLiterals() {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var defCount int32 = definitions.Size()
			_ = defCount
			for ; i_2 < defCount; i_2 = i_2 + 1 {
				var fnList *OptimizeCalls_NodeList = allocPooledVia_OptimizeCalls_ReferenceMap_DefinitionFunctionNodeListFor(&this_.slab_OptimizeCalls_NodeList, j2gort.CastNode[pkg_rhino.Node](definitions.Get(i_2)))
				_ = fnList
				{
					var j int32 = int32(0)
					_ = j
					var fnCount int32 = fnList.Size()
					_ = fnCount
					for ; j < fnCount; j = j + 1 {
						var functionNode pkg_rhino.Node = fnList.Get(j)
						_ = functionNode
						var firstParam pkg_rhino.Node = NodeUtil_GetFunctionParameters(functionNode).GetFirstChild()
						_ = firstParam
						if (firstParam != 0) && !(firstParam.IsUnusedParameter()) {
							if this_.decisionsLog.IsLogging() {
								this_.decisionsLog.Log_2("%s\t%s\twill not optimize parameters for tagged template literals", []any{j2gort.BoxString(refKind), j2gort.BoxString(key)})
							}
							return analysis
						}
					}
				}
			}
		}
		if this_.decisionsLog.IsLogging() {
			this_.decisionsLog.Log_2("%s\t%s\t%s", []any{j2gort.BoxString(refKind), j2gort.BoxString(key), "first param is unused, so we will convert tagged template literals to normal calls"})
		}
		seenCandidateUse = true
	}
	if !(seenCandidateUse) {
		if this_.decisionsLog.IsLogging() {
			this_.decisionsLog.Log_2("%s\t%s\tno usage found", []any{j2gort.BoxString(refKind), j2gort.BoxString(key)})
		}
		return analysis
	}
	{
		r_5 := analysis
		r_5.isSafeToOptimize = true
	}
	return analysis
}

// Java method: OptimizeParameters.convertTaggedTemplateLiteralToNormalCall(com.google.javascript.rhino.Node ttlNode)
func (this_ *OptimizeParameters) priv_OptimizeParameters_convertTaggedTemplateLiteralToNormalCall_4100(ttlNode pkg_rhino.Node) {
	_ = ttlNode
	var callee pkg_rhino.Node = j2gort.CheckNotNullNode(ttlNode.GetFirstChild())
	_ = callee
	var templateLiteral pkg_rhino.Node = callee.GetNext()
	_ = templateLiteral
	var firstArg pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}).Srcref(templateLiteral)
	_ = firstArg
	callee.Detach()
	var callNode pkg_rhino.Node = pkg_rhino.IR_Call(callee, []pkg_rhino.Node{firstArg}).Srcref(ttlNode)
	_ = callNode
	callNode.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), (ttlNode.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL()) || !(NodeUtil_IsNormalOrOptChainGet(callee))))
	{
		var tlChild pkg_rhino.Node = templateLiteral.GetFirstChild()
		_ = tlChild
		for ; tlChild != 0; tlChild = tlChild.GetNext() {
			if tlChild.IsTemplateLitSub() {
				var arg pkg_rhino.Node = tlChild.GetOnlyChild()
				_ = arg
				arg.Detach()
				callNode.AddChildToBack(arg)
			}
		}
	}
	ttlNode.ReplaceWith(callNode)
	this_.compiler.ReportChangeToEnclosingScope(callNode)
}

// Java method: OptimizeParameters.isCandidateDefinition(com.google.javascript.rhino.Node n)
func (this_ *OptimizeParameters) priv_OptimizeParameters_isCandidateDefinition_4101(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var functionExpr pkg_rhino.Node = pkg_rhino.Node(0)
	_ = functionExpr
	if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
		functionExpr = parent
	} else if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		functionExpr = parent.GetLastChild()
	} else if n.IsName() && n.HasChildren() {
		functionExpr = n.GetFirstChild()
	} else if this_.priv_OptimizeParameters_isClassMemberDefinition_4102(n) {
		functionExpr = n.GetFirstChild()
	} else if parent.IsClass() && n.IsFirstChildOf(parent) {
		functionExpr = parent
	} else {
		return false
	}
	return this_.priv_OptimizeParameters_allDefinitionsAreCandidateFunctions_4104(functionExpr)
}

// Java method: OptimizeParameters.isClassMemberDefinition(com.google.javascript.rhino.Node n)
func (this_ *OptimizeParameters) priv_OptimizeParameters_isClassMemberDefinition_4102(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsMemberFunctionDef() && n.GetParent().IsClassMembers())
}

// Java method: OptimizeParameters.classContainsClassFieldWithRHSSideEffects(com.google.javascript.rhino.Node classNode)
func (this_ *OptimizeParameters) priv_OptimizeParameters_classContainsClassFieldWithRHSSideEffects_4103(classNode pkg_rhino.Node) bool {
	_ = classNode
	var classMembersNode pkg_rhino.Node = NodeUtil_GetClassMembers(classNode)
	_ = classMembersNode
	{
		var member pkg_rhino.Node = classMembersNode.GetFirstChild()
		_ = member
		for ; member != 0; member = member.GetNext() {
			var initializer pkg_rhino.Node = member.GetLastChild()
			_ = initializer
			if ((initializer != 0) && NodeUtil_IsNonStaticClassFieldInitializer(initializer)) && this_.astAnalyzer.MayHaveSideEffects(initializer) {
				return true
			}
		}
	}
	return false
}

// Java method: OptimizeParameters.allDefinitionsAreCandidateFunctions(com.google.javascript.rhino.Node n)
func (this_ *OptimizeParameters) priv_OptimizeParameters_allDefinitionsAreCandidateFunctions_4104(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */ :
		if NodeUtil_IsNamedClassExpression(n) {
			return false
		} else {
			var constructorMemberFunctionDef pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
			_ = constructorMemberFunctionDef
			if constructorMemberFunctionDef == 0 {
				return false
			} else if this_.priv_OptimizeParameters_classContainsClassFieldWithRHSSideEffects_4103(n) {
				return false
			} else {
				var functionNode pkg_rhino.Node = constructorMemberFunctionDef.GetOnlyChild()
				_ = functionNode
				return ((!(NodeUtil_DoesFunctionReferenceOwnArgumentsObject(functionNode)) && !(OptimizeParameters_MayReferenceParamBeforeBody(functionNode))) && !(OptimizeParameters_IsUsedViaDotConstructor(functionNode)))
			}
		}
	case int32(73) /* Token.FUNCTION */ :
		return (((!(NodeUtil_IsNamedFunctionExpression(n)) && !(NodeUtil_DoesFunctionReferenceOwnArgumentsObject(n))) && !(OptimizeParameters_MayReferenceParamBeforeBody(n))) && !(OptimizeParameters_IsUsedViaDotConstructor(n)))
	case int32(101) /* Token.CAST */, int32(50) /* Token.COMMA */ :
		return this_.priv_OptimizeParameters_allDefinitionsAreCandidateFunctions_4104(n.GetLastChild())
	case int32(67) /* Token.HOOK */ :
		return (this_.priv_OptimizeParameters_allDefinitionsAreCandidateFunctions_4104(n.GetSecondChild()) && this_.priv_OptimizeParameters_allDefinitionsAreCandidateFunctions_4104(n.GetLastChild()))
	case int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(70) /* Token.COALESCE */ :
		return (this_.priv_OptimizeParameters_allDefinitionsAreCandidateFunctions_4104(n.GetFirstChild()) && this_.priv_OptimizeParameters_allDefinitionsAreCandidateFunctions_4104(n.GetLastChild()))
	default:
		return false
	}
	panic("unreachable")
}

// Java static method: OptimizeParameters.isUsedViaDotConstructor(com.google.javascript.rhino.Node function)
func OptimizeParameters_IsUsedViaDotConstructor(function pkg_rhino.Node) bool {
	_ = function
	var docInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(function)
	_ = docInfo
	return ((docInfo != nil) && docInfo.IsUsedViaDotConstructor())
}

// Java static method: OptimizeParameters.mayReferenceParamBeforeBody(com.google.javascript.rhino.Node function)
func OptimizeParameters_MayReferenceParamBeforeBody(function pkg_rhino.Node) bool {
	_ = function
	var paramList pkg_rhino.Node = function.GetSecondChild()
	_ = paramList
	var firstParam pkg_rhino.Node = paramList.GetFirstChild()
	_ = firstParam
	if firstParam == 0 {
		return false
	}
	if OptimizeParameters_AllSimpleNameParams(firstParam) {
		if firstParam.GetNext() == 0 {
			return false
		}
		{
			var p1 pkg_rhino.Node = firstParam
			_ = p1
			for ; p1 != 0; p1 = p1.GetNext() {
				var s1 string = OptimizeParameters_SimpleParamName(p1)
				_ = s1
				{
					var p2 pkg_rhino.Node = p1.GetNext()
					_ = p2
					for ; p2 != 0; p2 = p2.GetNext() {
						{
							var iifeCond_1 bool
							{

								sr_6 := s1
								se_7 := OptimizeParameters_SimpleParamName(p2)
								iifeCond_1 = j2gort.CheckNotNullString(sr_6) == se_7
							}
							if iifeCond_1 {
								return true
							}
						}
					}
				}
			}
		}
		return false
	}
	var namesByNames j2gort.OrderedMap[string, *OptimizeCalls_NodeList] = j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList]()
	_ = namesByNames
	var slabNodeList_ j2gort.PooledSlab[OptimizeCalls_NodeList]
	defer slabNodeList_.Release()

	allocPooledVia_OptimizeParameters_CollectParamListNames(&slabNodeList_, paramList, &namesByNames)
	{
		cur_8 := namesByNames.IterValues()
		for cur_8.HasNext() {

			names := cur_8.
				Next()
			_ = names
			var count int32 = names.Size()
			_ = count
			if count == int32(1) {
				continue
			}
			{
				var i int32 = int32(0)
				_ = i
				for ; i < count; i = i + 1 {
					if NodeUtil_IsLValue(names.Get(i)) {
						return true
					}
				}
			}
		}
	}
	return false
}

// Java static method: OptimizeParameters.allSimpleNameParams(com.google.javascript.rhino.Node firstParam)
func OptimizeParameters_AllSimpleNameParams(firstParam pkg_rhino.Node) bool {
	_ = firstParam
	{
		var p pkg_rhino.Node = firstParam
		_ = p
		for ; p != 0; p = p.GetNext() {
			if p.IsName() {
				continue
			}
			if p.IsRest() && p.GetOnlyChild().IsName() {
				continue
			}
			return false
		}
	}
	return true
}

// Java static method: OptimizeParameters.simpleParamName(com.google.javascript.rhino.Node param)
func OptimizeParameters_SimpleParamName(param pkg_rhino.Node) string {
	_ = param
	if param.IsName() {
		return param.GetString()
	} else {
		return param.GetOnlyChild().GetString()
	}
}

// Java static method: OptimizeParameters.collectParamListNames(com.google.javascript.rhino.Node n, java.util.LinkedHashMap<java.lang.String,com.google.javascript.jscomp.OptimizeCalls.NodeList> namesByNames)
func OptimizeParameters_CollectParamListNames(n pkg_rhino.Node, namesByNames *j2gort.Ext_java_util_LinkedHashMap) {
	_ = n
	_ = namesByNames
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			OptimizeParameters_CollectParamListNames(c, namesByNames)
		}
	}
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		var list *OptimizeCalls_NodeList = j2gort.CastPtr[OptimizeCalls_NodeList](j2gort.MapGetStr(namesByNames, name))
		_ = list
		if list == nil {
			list = NewOptimizeCalls_NodeList()
			j2gort.MapPutStr(namesByNames, name, j2gort.PtrToAny(list))
		}
		list.Add(n)
	}
}
func allocPooledVia_OptimizeParameters_CollectParamListNames(slab_ *j2gort.PooledSlab[OptimizeCalls_NodeList], n pkg_rhino.Node, namesByNames *j2gort.OrderedMap[string, *OptimizeCalls_NodeList]) {
	_ = n
	_ = namesByNames
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			allocPooledVia_OptimizeParameters_CollectParamListNames(slab_, c, namesByNames)
		}
	}
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		var list *OptimizeCalls_NodeList = namesByNames.Get(name)
		_ = list
		if list == nil {
			list = allocPooled_OptimizeCalls_NodeList(slab_)
			namesByNames.Put(name, list)
		}
		list.Add(n)
	}
}

// Java method: OptimizeParameters.tryEliminateOptionalArgs(com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeParameters) priv_OptimizeParameters_tryEliminateOptionalArgs_4105(refs *OptimizeCalls_NodeList) {
	_ = refs
	defer this_.slab_OptimizeCalls_NodeList.Reset()
	var maxArgs int32 = int32(-1)
	_ = maxArgs
	{
		var i int32 = int32(0)
		_ = i
		var refCount int32 = refs.Size()
		_ = refCount
		for ; i < refCount; i = i + 1 {
			var n pkg_rhino.Node = refs.Get(i)
			_ = n
			if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
				if OptimizeParameters_HasSpreadInFirstCallParam(n) {
					return
				}
				var numArgs int32 = int32(0)
				_ = numArgs
				var firstArg pkg_rhino.Node = OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
				_ = firstArg
				{
					var c pkg_rhino.Node = firstArg
					_ = c
					for ; c != 0; c = c.GetNext() {
						numArgs = numArgs + 1
						if c.IsSpread() {
							return
						}
					}
				}
				if numArgs > maxArgs {
					maxArgs = numArgs
				}
			}
		}
	}
	var fnList *OptimizeCalls_NodeList = allocPooledVia_OptimizeCalls_ReferenceMap_GetFunctionNodeList(&this_.slab_OptimizeCalls_NodeList, refs)
	_ = fnList
	{
		var i_2 int32 = int32(0)
		_ = i_2
		var fnCount int32 = fnList.Size()
		_ = fnCount
		for ; i_2 < fnCount; i_2 = i_2 + 1 {
			this_.priv_OptimizeParameters_eliminateParamsAfter_4119(fnList.Get(i_2), maxArgs)
		}
	}
}

// Java method: OptimizeParameters.tryEliminateConstantArgs(com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeParameters) priv_OptimizeParameters_tryEliminateConstantArgs_4106(refs *OptimizeCalls_NodeList) {
	_ = refs
	defer func() {
		this_.slab_ParameterList.Reset()
		this_.slab_Parameter.Reset()
		this_.slab_OptimizeCalls_NodeList.Reset()
	}()
	var parameters *OptimizeParameters_ParameterList = this_.priv_OptimizeParameters_findFixedArguments_4107(refs)
	_ = parameters
	if parameters == nil {
		return
	}
	var fns *OptimizeCalls_NodeList = allocPooledVia_OptimizeCalls_ReferenceMap_GetFunctionNodeList(&this_.slab_OptimizeCalls_NodeList, refs)
	_ = fns
	if fns.Size() > int32(1) {
		return
	}
	if !(fns.Size() == int32(1)) {
		j2gort.CheckState(false)
	}
	var fn pkg_rhino.Node = fns.Get(int32(0))
	_ = fn
	var continueLooking bool = OptimizeParameters_AdjustForConstraints(fn, parameters)
	_ = continueLooking
	if !(continueLooking) {
		return
	}
	{
		var i int32 = int32(0)
		_ = i
		var refCount int32 = refs.Size()
		_ = refCount
		for ; i < refCount; i = i + 1 {
			var n pkg_rhino.Node = refs.Get(i)
			_ = n
			if !(OptimizeParameters_AlreadyRemoved(n)) && OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
				this_.priv_OptimizeParameters_optimizeCallSite_4112(parameters, n)
			}
		}
	}
	this_.priv_OptimizeParameters_optimizeFunctionDefinition_4111(parameters, fn)
}

// Java method: OptimizeParameters.findFixedArguments(com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeParameters) priv_OptimizeParameters_findFixedArguments_4107(refs *OptimizeCalls_NodeList) *OptimizeParameters_ParameterList {
	_ = refs
	val_10 := this_.slab_ParameterList.Alloc()

	InitOptimizeParameters_ParameterList(val_10)
	var parameters *OptimizeParameters_ParameterList = val_10
	_ = parameters
	var firstCall bool = true
	_ = firstCall
	var continueLooking bool = false
	_ = continueLooking
	{
		var i int32 = int32(0)
		_ = i
		var refCount int32 = refs.Size()
		_ = refCount
		for ; i < refCount; i = i + 1 {
			var n pkg_rhino.Node = refs.Get(i)
			_ = n
			if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
				if OptimizeParameters_HasSpreadInFirstCallParam(n) {
					continueLooking = false
					break
				}
				var cur pkg_rhino.Node = OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
				_ = cur
				if firstCall {
					continueLooking = this_.priv_OptimizeParameters_buildInitialParameterList_4109(parameters, cur)
					firstCall = false
				} else {
					continueLooking = this_.priv_OptimizeParameters_findFixedParameters_4108(parameters, cur)
				}
				if !(continueLooking) {
					break
				}
			}
		}
	}
	if continueLooking {
		return parameters
	} else {
		return nil
	}
}

// Java static method: OptimizeParameters.hasSpreadInFirstCallParam(com.google.javascript.rhino.Node n)
func OptimizeParameters_HasSpreadInFirstCallParam(n pkg_rhino.Node) bool {
	_ = n
	var call pkg_rhino.Node = OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n)
	_ = call
	if NodeUtil_IsFunctionObjectCall(call) {
		var firstDotCallParam pkg_rhino.Node = call.GetSecondChild()
		_ = firstDotCallParam
		return ((firstDotCallParam != 0) && firstDotCallParam.IsSpread())
	}
	return false
}

// Java static method: OptimizeParameters.adjustForConstraints(com.google.javascript.rhino.Node fn, com.google.javascript.jscomp.OptimizeParameters.ParameterList parameters)
func OptimizeParameters_AdjustForConstraints(fn pkg_rhino.Node, parameters *OptimizeParameters_ParameterList) bool {
	_ = fn
	_ = parameters
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(fn)
	_ = info
	if (info != nil) && info.IsNoInline() {
		return false
	}
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn)
	_ = paramList
	var lastFormal pkg_rhino.Node = paramList.GetLastChild()
	_ = lastFormal
	var restIndex int32 = int32(2147483647)
	_ = restIndex
	var lastNonRestFormal int32 = (paramList.GetChildCount() - int32(1))
	_ = lastNonRestFormal
	var formal pkg_rhino.Node = lastFormal
	_ = formal
	if (lastFormal != 0) && lastFormal.IsRest() {
		restIndex = lastNonRestFormal
		lastNonRestFormal = lastNonRestFormal - 1
		formal = formal.GetPrevious()
	}
	var anyMovable bool = false
	_ = anyMovable
	var seenUnmovableSideEffects bool = false
	_ = seenUnmovableSideEffects
	var seenUnmoveableSideEffected bool = false
	_ = seenUnmoveableSideEffected
	var allRestValueRemovable bool = true
	_ = allRestValueRemovable
	{
		var i int32 = (parameters.Size() - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			var current *OptimizeParameters_Parameter = parameters.Get(i)
			_ = current
			if i <= lastNonRestFormal {
				if formal.IsDefaultValue() && current.mayBeUndefined {
					{
						r_11 := current
						r_11.shouldRemove = false
					}
				}
				formal = formal.GetPrevious()
			}
			if current.shouldRemove && ((seenUnmovableSideEffects && current.CanBeSideEffected()) || (seenUnmoveableSideEffected && current.HasSideEffects())) {
				{
					r_12 := current
					r_12.shouldRemove = false
				}
			}
			if i >= restIndex {
				if allRestValueRemovable {
					if !(current.shouldRemove) {
						anyMovable = false
						allRestValueRemovable = false
						{
							var j int32 = (i + int32(1))
							_ = j
							for ; j < parameters.Size(); j = j + 1 {
								var p *OptimizeParameters_Parameter = parameters.Get(j)
								_ = p
								{
									r_13 := p
									r_13.shouldRemove = false
								}
								if p.canBeSideEffected {
									seenUnmoveableSideEffected = true
								}
								if p.hasSideEffects {
									seenUnmovableSideEffects = true
								}
							}
						}
					}
				} else {
					{
						r_14 := current
						r_14.shouldRemove = false
					}
				}
			}
			if current.shouldRemove {
				anyMovable = true
			} else {
				if current.canBeSideEffected {
					seenUnmoveableSideEffected = true
				}
				if current.hasSideEffects {
					seenUnmovableSideEffects = true
				}
			}
		}
	}
	return anyMovable
}

// Java method: OptimizeParameters.findFixedParameters(com.google.javascript.jscomp.OptimizeParameters.ParameterList parameters, com.google.javascript.rhino.Node cur)
func (this_ *OptimizeParameters) priv_OptimizeParameters_findFixedParameters_4108(parameters *OptimizeParameters_ParameterList, cur pkg_rhino.Node) bool {
	_ = parameters
	_ = cur
	var anyMovable bool = false
	_ = anyMovable
	var index int32 = int32(0)
	_ = index
	for cur != 0 {
		var p *OptimizeParameters_Parameter = nil
		_ = p
		if index >= parameters.Size() {
			p = allocPooled_OptimizeParameters_Parameter(&this_.slab_Parameter, cur, false)
			parameters.Add(p)
			this_.priv_OptimizeParameters_setParameterSideEffectInfo_4110(p, cur)
		} else {
			p = parameters.Get(index)
			if p.ShouldRemove() {
				var value pkg_rhino.Node = p.GetArg()
				_ = value
				if !(cur.IsEquivalentTo_1(value)) {
					p.SetShouldRemove(false)
				} else {
					anyMovable = true
				}
			}
			p.SetHasSideEffects((p.HasSideEffects() || this_.astAnalyzer.MayHaveSideEffects(cur)))
			p.SetCanBeSideEffected((p.CanBeSideEffected() || NodeUtil_CanBeSideEffected_1(cur)))
		}
		if cur.IsSpread() {
			break
		}
		cur = cur.GetNext()
		index = index + 1
	}
	for ; index < parameters.Size(); index = index + 1 {
		parameters.Get(index).SetShouldRemove(false)
	}
	return anyMovable
}

// Java method: OptimizeParameters.buildInitialParameterList(com.google.javascript.jscomp.OptimizeParameters.ParameterList parameters, com.google.javascript.rhino.Node cur)
func (this_ *OptimizeParameters) priv_OptimizeParameters_buildInitialParameterList_4109(parameters *OptimizeParameters_ParameterList, cur pkg_rhino.Node) bool {
	_ = parameters
	_ = cur
	var anyMovable bool = false
	_ = anyMovable
	for cur != 0 {
		var movable bool = OptimizeParameters_IsMovableValue(cur, this_.globalScope)
		_ = movable
		var p *OptimizeParameters_Parameter = allocPooled_OptimizeParameters_Parameter(&this_.slab_Parameter, cur, movable)
		_ = p
		this_.priv_OptimizeParameters_setParameterSideEffectInfo_4110(p, cur)
		parameters.Add(p)
		if movable {
			anyMovable = true
		}
		if cur.IsSpread() {
			break
		}
		cur = cur.GetNext()
	}
	return anyMovable
}

// Java method: OptimizeParameters.setParameterSideEffectInfo(com.google.javascript.jscomp.OptimizeParameters.Parameter p, com.google.javascript.rhino.Node value)
func (this_ *OptimizeParameters) priv_OptimizeParameters_setParameterSideEffectInfo_4110(p *OptimizeParameters_Parameter, value pkg_rhino.Node) {
	_ = p
	_ = value
	p.SetHasSideEffects(this_.astAnalyzer.MayHaveSideEffects(value))
	p.SetCanBeSideEffected(NodeUtil_CanBeSideEffected_1(value))
	if !(value.IsSpread()) {
		p.SetMayBeUndefined(NodeUtil_MayBeUndefined(value))
	}
}

// Java static method: OptimizeParameters.isMovableValue(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Scope globalScope)
func OptimizeParameters_IsMovableValue(n pkg_rhino.Node, globalScope *Scope) bool {
	_ = n
	_ = globalScope
	switch n.TokenSwitchOrd_() {
	case int32(115) /* Token.AWAIT */, int32(114) /* Token.YIELD */, int32(37) /* Token.THIS */, int32(110) /* Token.SUPER */, int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */, int32(135) /* Token.NEW_TARGET */, int32(136) /* Token.IMPORT_META */ :
		return false
	case int32(73) /* Token.FUNCTION */ :
		return false
	case int32(32) /* Token.NAME */ :
		if j2gort.CheckNotNullString(n.GetString()) == "arguments" {
			return false
		} else {
			var v *Var = j2gort.CastPtr[Var](globalScope.GetVar(n.GetString()))
			_ = v
			if v == nil {
				return false
			}
		}
	default:
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if !(OptimizeParameters_IsMovableValue(c, globalScope)) {
				return false
			}
		}
	}
	return true
}

// Java method: OptimizeParameters.optimizeFunctionDefinition(com.google.javascript.jscomp.OptimizeParameters.ParameterList parameters, com.google.javascript.rhino.Node fn)
func (this_ *OptimizeParameters) priv_OptimizeParameters_optimizeFunctionDefinition_4111(parameters *OptimizeParameters_ParameterList, fn pkg_rhino.Node) {
	_ = parameters
	_ = fn
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn)
	_ = paramList
	var maybeRest pkg_rhino.Node = paramList.GetLastChild()
	_ = maybeRest
	var lastParameter int32 = (parameters.Size() - int32(1))
	_ = lastParameter
	if (maybeRest != 0) && maybeRest.IsRest() {
		var restIndex int32 = (paramList.GetChildCount() - int32(1))
		_ = restIndex
		if (parameters.Size() < restIndex) || ((restIndex < parameters.Size()) && parameters.Get(restIndex).ShouldRemove()) {
			var value pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}).Srcref(maybeRest)
			_ = value
			{
				var i int32 = restIndex
				_ = i
				for ; i < parameters.Size(); i = i + 1 {
					var parameter *OptimizeParameters_Parameter = parameters.Get(i)
					_ = parameter
					if !(parameter.ShouldRemove()) {
						j2gort.CheckState(false)
					}
					value.AddChildToBack(parameters.Get(i).GetArg())
				}
			}
			maybeRest.Detach()
			var lhs pkg_rhino.Node = maybeRest.RemoveFirstChild()
			_ = lhs
			this_.priv_OptimizeParameters_addRestVariableToFunction_4113(fn, lhs, value)
		}
		lastParameter = min((parameters.Size() - int32(1)), (restIndex - int32(1)))
	}
	{
		var i_2 int32 = lastParameter
		_ = i_2
		for ; i_2 >= int32(0); i_2 = i_2 - 1 {
			var parameter_2 *OptimizeParameters_Parameter = parameters.Get(i_2)
			_ = parameter_2
			if parameter_2.ShouldRemove() {
				var formalParam pkg_rhino.Node = NodeUtil_GetArgumentForFunction(fn, i_2)
				_ = formalParam
				if formalParam != 0 {
					formalParam.Detach()
					if formalParam.IsDefaultValue() {
						if !(!(parameter_2.mayBeUndefined)) {
							j2gort.CheckState(false)
						}
						formalParam = formalParam.RemoveFirstChild()
					}
				}
				var value_2 pkg_rhino.Node = parameters.Get(i_2).GetArg()
				_ = value_2
				if formalParam != 0 {
					this_.priv_OptimizeParameters_addVariableDeclarationToFunction_4114(fn, formalParam, value_2)
				} else {
					this_.priv_OptimizeParameters_addExpressionToFunction_4115(fn, value_2)
				}
			}
		}
	}
}

// Java method: OptimizeParameters.optimizeCallSite(com.google.javascript.jscomp.OptimizeParameters.ParameterList parameters, com.google.javascript.rhino.Node target)
func (this_ *OptimizeParameters) priv_OptimizeParameters_optimizeCallSite_4112(parameters *OptimizeParameters_ParameterList, target pkg_rhino.Node) {
	_ = parameters
	_ = target
	var call pkg_rhino.Node = OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(target)
	_ = call
	var callMayMutateArgs bool = call.MayMutateArguments()
	_ = callMayMutateArgs
	var callMayMutateGlobalsOrThrow bool = call.MayMutateGlobalStateOrThrow()
	_ = callMayMutateGlobalsOrThrow
	{
		var index int32 = (parameters.Size() - int32(1))
		_ = index
		for ; index >= int32(0); index = index - 1 {
			var p *OptimizeParameters_Parameter = parameters.Get(index)
			_ = p
			if p.ShouldRemove() {
				this_.priv_OptimizeParameters_eliminateCallTargetArgAt_4122(target, index)
				if !(callMayMutateGlobalsOrThrow) && (p.HasSideEffects() || (callMayMutateArgs && !(NodeUtil_IsImmutableValue(p.GetArg())))) {
					callMayMutateGlobalsOrThrow = true
					call.SetSideEffectFlags_1_Node_SideEffectFlags(pkg_rhino.NewNode_SideEffectFlags_1(call.GetSideEffectFlags()).SetMutatesGlobalState())
				}
			}
		}
	}
}

// Java method: OptimizeParameters.addRestVariableToFunction(com.google.javascript.rhino.Node function, com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node value)
func (this_ *OptimizeParameters) priv_OptimizeParameters_addRestVariableToFunction_4113(function pkg_rhino.Node, lhs pkg_rhino.Node, value pkg_rhino.Node) {
	_ = function
	_ = lhs
	_ = value
	if !(lhs.GetParent() == 0) {
		j2gort.CheckState(false)
	}
	this_.priv_OptimizeParameters_addVariableDeclarationToFunction_4114(function, lhs, value)
}

// Java method: OptimizeParameters.addVariableDeclarationToFunction(com.google.javascript.rhino.Node function, com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node value)
func (this_ *OptimizeParameters) priv_OptimizeParameters_addVariableDeclarationToFunction_4114(function pkg_rhino.Node, lhs pkg_rhino.Node, value pkg_rhino.Node) {
	_ = function
	_ = lhs
	_ = value
	if !(value.GetParent() == 0) {
		j2gort.CheckStateStr(false, "Value must be detached")
	}
	if !(lhs != 0) {
		j2gort.CheckStateStr(false, "formal parameter being declared must not be null")
	}
	if !(lhs.GetParent() == 0) {
		j2gort.CheckStateStr(false, "formal parameter being declared must be detached")
	}
	var stmts []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = stmts
	if lhs.IsDestructuringPattern() {
		this_.priv_OptimizeParameters_rewriteDestructuringPattern_4121(lhs, value, &stmts)
	} else {
		stmts = append(stmts, NodeUtil_NewVarNode_2_Node_Node(lhs, value))
	}
	this_.priv_OptimizeParameters_insertStatements_4116(function, &stmts)
}

// Java method: OptimizeParameters.addExpressionToFunction(com.google.javascript.rhino.Node function, com.google.javascript.rhino.Node value)
func (this_ *OptimizeParameters) priv_OptimizeParameters_addExpressionToFunction_4115(function pkg_rhino.Node, value pkg_rhino.Node) {
	_ = function
	_ = value
	var stmts []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = stmts
	stmts = append(stmts, pkg_rhino.IR_ExprResult(value).Srcref(value))
	this_.priv_OptimizeParameters_insertStatements_4116(function, &stmts)
}

// Java method: OptimizeParameters.insertStatements(com.google.javascript.rhino.Node function, java.util.ArrayList<com.google.javascript.rhino.Node> stmts)
func (this_ *OptimizeParameters) priv_OptimizeParameters_insertStatements_4116(function pkg_rhino.Node, stmts *[]pkg_rhino.Node) {
	_ = function
	_ = stmts
	if len(*stmts) == 0 {
		return
	}
	var block pkg_rhino.Node = NodeUtil_GetFunctionBody(function)
	_ = block
	var insertionPoint pkg_rhino.Node = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(block)
	_ = insertionPoint
	if insertionPoint == 0 {
		this_.priv_OptimizeParameters_addStatementsToBack_4117(block, stmts)
	} else {
		this_.priv_OptimizeParameters_addStatementsBefore_4118(insertionPoint, stmts)
	}
	this_.compiler.ReportChangeToEnclosingScope((*stmts)[int32(0)])
}

// Java method: OptimizeParameters.addStatementsToBack(com.google.javascript.rhino.Node block, java.util.ArrayList<com.google.javascript.rhino.Node> stmts)
func (this_ *OptimizeParameters) priv_OptimizeParameters_addStatementsToBack_4117(block pkg_rhino.Node, stmts *[]pkg_rhino.Node) {
	_ = block
	_ = stmts
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(*stmts))
		_ = size
		for ; i < size; i = i + 1 {
			block.AddChildToBack((*stmts)[i])
		}
	}
}

// Java method: OptimizeParameters.addStatementsBefore(com.google.javascript.rhino.Node insertionPoint, java.util.ArrayList<com.google.javascript.rhino.Node> stmts)
func (this_ *OptimizeParameters) priv_OptimizeParameters_addStatementsBefore_4118(insertionPoint pkg_rhino.Node, stmts *[]pkg_rhino.Node) {
	_ = insertionPoint
	_ = stmts
	{
		var i int32 = (int32(len(*stmts)) - int32(1))
		_ = i
		for ; i >= int32(0); i = i - 1 {
			(*stmts)[i].InsertBefore(insertionPoint)
		}
	}
}

// Java method: OptimizeParameters.eliminateParamsAfter(com.google.javascript.rhino.Node fnNode, int argIndex)
func (this_ *OptimizeParameters) priv_OptimizeParameters_eliminateParamsAfter_4119(fnNode pkg_rhino.Node, argIndex int32) {
	_ = fnNode
	_ = argIndex
	var formalArgPtr pkg_rhino.Node = NodeUtil_GetFunctionParameters(fnNode).GetFirstChild()
	_ = formalArgPtr
	for (argIndex != int32(0)) && (formalArgPtr != 0) {
		formalArgPtr = formalArgPtr.GetNext()
		argIndex = argIndex - 1
	}
	this_.priv_OptimizeParameters_eliminateParamsAfter_4120(fnNode, formalArgPtr)
}

// Java method: OptimizeParameters.eliminateParamsAfter(com.google.javascript.rhino.Node fnNode, com.google.javascript.rhino.Node formal)
func (this_ *OptimizeParameters) priv_OptimizeParameters_eliminateParamsAfter_4120(fnNode pkg_rhino.Node, formal pkg_rhino.Node) {
	_ = fnNode
	_ = formal
	if formal != 0 {
		this_.priv_OptimizeParameters_eliminateParamsAfter_4120(fnNode, formal.GetNext())
		var stmts []pkg_rhino.Node = []pkg_rhino.Node{}
		_ = stmts
		formal.Detach()
		if formal.IsRest() {
			if !(formal.GetNext() == 0) {
				j2gort.CheckState(false)
			}
			var lhs pkg_rhino.Node = formal.RemoveFirstChild()
			_ = lhs
			var value pkg_rhino.Node = pkg_rhino.IR_Arraylit_1_Node([]pkg_rhino.Node{}).Srcref(formal)
			_ = value
			this_.priv_OptimizeParameters_addVariableDeclarationToFunction_4114(fnNode, lhs, value)
		} else if formal.IsDefaultValue() {
			var lhs_2 pkg_rhino.Node = formal.RemoveFirstChild()
			_ = lhs_2
			var value_2 pkg_rhino.Node = formal.GetLastChild().Detach()
			_ = value_2
			this_.priv_OptimizeParameters_addVariableDeclarationToFunction_4114(fnNode, lhs_2, value_2)
		} else if formal.IsDestructuringPattern() {
			var value_3 pkg_rhino.Node = NodeUtil_NewUndefinedNode(formal)
			_ = value_3
			this_.priv_OptimizeParameters_rewriteDestructuringPattern_4121(formal, value_3, &stmts)
			this_.priv_OptimizeParameters_insertStatements_4116(fnNode, &stmts)
		} else {
			stmts = append(stmts, pkg_rhino.IR_Var_1(formal).SrcrefIfMissing(formal))
			this_.priv_OptimizeParameters_insertStatements_4116(fnNode, &stmts)
		}
	}
}

// Java method: OptimizeParameters.rewriteDestructuringPattern(com.google.javascript.rhino.Node destructuringPattern, com.google.javascript.rhino.Node value, java.util.ArrayList<com.google.javascript.rhino.Node> stmts)
func (this_ *OptimizeParameters) priv_OptimizeParameters_rewriteDestructuringPattern_4121(destructuringPattern pkg_rhino.Node, value pkg_rhino.Node, stmts *[]pkg_rhino.Node) {
	_ = destructuringPattern
	_ = value
	_ = stmts
	if !(destructuringPattern.IsDestructuringPattern()) {
		j2gort.CheckState(false)
	}
	if !(!(destructuringPattern.HasParent())) {
		j2gort.CheckStateStr(false, "Formal parameter must be detached for rewriting")
	}
	NodeUtil_VisitLhsNodesInDestructuringPattern(destructuringPattern, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_15 any) {
		name := j2gort.CastNode[pkg_rhino.Node](lp_15)
		_ = name
		if !(name.IsName()) {
			j2gort.CheckStateFmt(false, "lhs in destructuring declaration should be a simple name. (%s)", j2gort.BoxNode(name))
		}
		var newName pkg_rhino.Node = pkg_rhino.IR_Name(name.GetString()).Srcref(name)
		_ = newName
		var newVar pkg_rhino.Node = pkg_rhino.IR_Var_1(newName).Srcref(name)
		_ = newVar
		*stmts = append(*stmts, newVar)
	}))
	var expr pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Assign(destructuringPattern, value).Srcref(destructuringPattern)).Srcref(destructuringPattern)
	_ = expr
	*stmts = append(*stmts, expr)
}

// Java method: OptimizeParameters.eliminateCallTargetArgAt(com.google.javascript.rhino.Node ref, int argIndex)
func (this_ *OptimizeParameters) priv_OptimizeParameters_eliminateCallTargetArgAt_4122(ref pkg_rhino.Node, argIndex int32) {
	_ = ref
	_ = argIndex
	var callArgNode pkg_rhino.Node = OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(ref, argIndex)
	_ = callArgNode
	if callArgNode != 0 {
		NodeUtil_DeleteNode(callArgNode, this_.compiler)
	}
}

// Java class: com.google.javascript.jscomp.OptimizeParameters.UnusedParameterOptimizer
type OptimizeParameters_UnusedParameterOptimizer struct {
	outer_            *OptimizeParameters
	toRemove          []pkg_rhino.Node `gson:"toRemove"`
	toReplaceWithZero []pkg_rhino.Node `gson:"toReplaceWithZero"`
	idHash_           atomic.Int32
}

func (this_ *OptimizeParameters_UnusedParameterOptimizer) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeParameters_UnusedParameterOptimizer) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeParameters_UnusedParameterOptimizer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeParameters$UnusedParameterOptimizer", "UnusedParameterOptimizer"
}
func (this_ *OptimizeParameters_UnusedParameterOptimizer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeParameters_UnusedParameterOptimizer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeParameters.UnusedParameterOptimizer", "com.google.javascript.jscomp.OptimizeParameters$UnusedParameterOptimizer", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeParameters_UnusedParameterOptimizer](v_)
	})
}
func (this_ *OptimizeParameters_UnusedParameterOptimizer) SetOuter_(o *OptimizeParameters) {
	this_.outer_ = o
}
func (this_ *OptimizeParameters_UnusedParameterOptimizer) GetOuter_() *OptimizeParameters {
	return this_.outer_
}

func (this_ *OptimizeParameters_UnusedParameterOptimizer) initFields_OptimizeParameters_UnusedParameterOptimizer() {
	_ = this_
	this_.toRemove = []pkg_rhino.Node{}
	this_.toReplaceWithZero = []pkg_rhino.Node{}
}

// Java constructor: UnusedParameterOptimizer()
func NewOptimizeParameters_UnusedParameterOptimizer(outer_ *OptimizeParameters) *OptimizeParameters_UnusedParameterOptimizer {
	this_ := &OptimizeParameters_UnusedParameterOptimizer{}
	this_.outer_ = outer_
	InitOptimizeParameters_UnusedParameterOptimizer(this_)
	return this_
}

func InitOptimizeParameters_UnusedParameterOptimizer(this_ *OptimizeParameters_UnusedParameterOptimizer) {
	this_.initFields_OptimizeParameters_UnusedParameterOptimizer()
}

// Java method: UnusedParameterOptimizer.tryEliminateUnusedArgs(com.google.javascript.jscomp.OptimizeCalls.NodeList refs)
func (this_ *OptimizeParameters_UnusedParameterOptimizer) TryEliminateUnusedArgs(refs *OptimizeCalls_NodeList) {
	_ = refs
	defer this_.outer_.slab_OptimizeCalls_NodeList.Reset()
	var fns *OptimizeCalls_NodeList = allocPooledVia_OptimizeCalls_ReferenceMap_GetFunctionNodeList(&this_.outer_.slab_OptimizeCalls_NodeList, refs)
	_ = fns
	if !(!(fns.IsEmpty())) {
		j2gort.CheckState(false)
	}
	var maxFormalsCount int32 = int32(0)
	_ = maxFormalsCount
	var lowestUsedRest int32 = int32(2147483647)
	_ = lowestUsedRest
	var used *j2gort.Ext_java_util_BitSet = j2gort.NewExt_java_util_BitSet()
	_ = used
	{
		var fi int32 = int32(0)
		_ = fi
		var fnCount int32 = fns.Size()
		_ = fnCount
		for ; fi < fnCount; fi = fi + 1 {
			var fn pkg_rhino.Node = fns.Get(fi)
			_ = fn
			var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn)
			_ = paramList
			var index int32 = int32(-1)
			_ = index
			{
				var c pkg_rhino.Node = paramList.GetFirstChild()
				_ = c
				for ; c != 0; c = c.GetNext() {
					index = index + 1
					if !(c.IsUnusedParameter()) {
						used.Set_1(index)
						if c.IsRest() {
							lowestUsedRest = min(lowestUsedRest, index)
							if lowestUsedRest == int32(0) {
								return
							}
						}
					}
				}
			}
			maxFormalsCount = max(maxFormalsCount, (index + int32(1)))
		}
	}
	if lowestUsedRest < maxFormalsCount {
		used.Set_2_int32_int32(lowestUsedRest, maxFormalsCount)
	}
	var unused *j2gort.Ext_java_util_BitSet = (j2gort.CastPtr[j2gort.Ext_java_util_BitSet](used.Clone()))
	_ = unused
	unused.Flip_2(int32(0), maxFormalsCount)
	if (lowestUsedRest < maxFormalsCount) && (unused.Cardinality() == int32(0)) {
		return
	}
	var unremovable *j2gort.Ext_java_util_BitSet = (j2gort.CastPtr[j2gort.Ext_java_util_BitSet](used.Clone()))
	_ = unremovable
	var lowestSpread int32 = int32(2147483647)
	_ = lowestSpread
	{
		var i int32 = int32(0)
		_ = i
		var refCount int32 = refs.Size()
		_ = refCount
		for ; i < refCount; i = i + 1 {
			var n pkg_rhino.Node = refs.Get(i)
			_ = n
			if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
				if OptimizeParameters_HasSpreadInFirstCallParam(n) {
					return
				}
				var param pkg_rhino.Node = OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
				_ = param
				var paramIndex int32 = int32(0)
				_ = paramIndex
				for param != 0 {
					if paramIndex >= maxFormalsCount {
						break
					}
					if param.IsSpread() {
						lowestSpread = min(lowestSpread, paramIndex)
						break
					}
					if !(unremovable.Get_1(paramIndex)) && this_.outer_.astAnalyzer.MayHaveSideEffects(param) {
						unremovable.Set_1(paramIndex)
					}
					param = param.GetNext()
					paramIndex = paramIndex + 1
				}
			}
		}
	}
	if lowestSpread < maxFormalsCount {
		unremovable.Set_2_int32_int32(lowestSpread, maxFormalsCount)
	}
	var removeAllAfterIndex int32 = int32(2147483647)
	_ = removeAllAfterIndex
	if lowestUsedRest == int32(2147483647) {
		removeAllAfterIndex = (maxFormalsCount - int32(1))
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		var refCount_2 int32 = refs.Size()
		_ = refCount_2
		for ; i_2 < refCount_2; i_2 = i_2 + 1 {
			var n_2 pkg_rhino.Node = refs.Get(i_2)
			_ = n_2
			if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n_2) && !(OptimizeParameters_AlreadyRemoved(n_2)) {
				var arg pkg_rhino.Node = OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n_2)
				_ = arg
				this_.RecordRemovalCallArguments(lowestUsedRest, removeAllAfterIndex, unused, unremovable, arg, int32(0))
			}
		}
	}
	{
		var fi_2 int32 = int32(0)
		_ = fi_2
		var fnCount_2 int32 = fns.Size()
		_ = fnCount_2
		for ; fi_2 < fnCount_2; fi_2 = fi_2 + 1 {
			var fn_2 pkg_rhino.Node = fns.Get(fi_2)
			_ = fn_2
			var paramList_2 pkg_rhino.Node = NodeUtil_GetFunctionParameters(fn_2)
			_ = paramList_2
			var param_2 pkg_rhino.Node = paramList_2.GetFirstChild()
			_ = param_2
			this_.RemoveUnusedFunctionParameters(unremovable, param_2, int32(0))
		}
	}
}

// Java method: UnusedParameterOptimizer.recordRemovalCallArguments(int firstRestIndex, int removeAllAfterIndex, java.util.BitSet unused, java.util.BitSet unremovable, com.google.javascript.rhino.Node arg, int index)
func (this_ *OptimizeParameters_UnusedParameterOptimizer) RecordRemovalCallArguments(firstRestIndex int32, removeAllAfterIndex int32, unused *j2gort.Ext_java_util_BitSet, unremovable *j2gort.Ext_java_util_BitSet, arg pkg_rhino.Node, index int32) bool {
	_ = firstRestIndex
	_ = removeAllAfterIndex
	_ = unused
	_ = unremovable
	_ = arg
	_ = index
	if arg == 0 {
		return true
	}
	if index > removeAllAfterIndex {
		return this_.RemoveArgAndFollowing(arg)
	}
	if arg.IsSpread() {
		return false
	}
	var removedAllTrailing bool = this_.RecordRemovalCallArguments(firstRestIndex, removeAllAfterIndex, unused, unremovable, arg.GetNext(), (index + int32(1)))
	_ = removedAllTrailing
	if index < firstRestIndex {
		var isRemovableTrailingUndefined bool = (NodeUtil_IsUndefined(arg) && removedAllTrailing)
		_ = isRemovableTrailingUndefined
		if isRemovableTrailingUndefined && !(this_.outer_.astAnalyzer.MayHaveSideEffects(arg)) {
			this_.toRemove = append(this_.toRemove, arg)
			return true
		}
		if unused.Get_1(index) {
			var hasSideEffects bool = (isRemovableTrailingUndefined || this_.outer_.astAnalyzer.MayHaveSideEffects(arg))
			_ = hasSideEffects
			if !(hasSideEffects) {
				if unremovable.Get_1(index) {
					this_.priv_UnusedParameterOptimizer_recordReplaceWithZero_4126(arg)
				} else {
					this_.toRemove = append(this_.toRemove, arg)
					return removedAllTrailing
				}
			} else {
				this_.RemoveFromCommaIfNoSideEffects(arg)
			}
		}
	}
	return false
}

// Java method: UnusedParameterOptimizer.removeFromCommaIfNoSideEffects(com.google.javascript.rhino.Node arg)
func (this_ *OptimizeParameters_UnusedParameterOptimizer) RemoveFromCommaIfNoSideEffects(arg pkg_rhino.Node) {
	_ = arg
	if arg.IsComma() {
		var secondChild pkg_rhino.Node = arg.GetSecondChild()
		_ = secondChild
		if !(this_.outer_.astAnalyzer.MayHaveSideEffects(secondChild)) {
			this_.priv_UnusedParameterOptimizer_recordReplaceWithZero_4126(secondChild)
		} else {
			this_.RemoveFromCommaIfNoSideEffects(secondChild)
		}
	}
}

// Java method: UnusedParameterOptimizer.recordReplaceWithZero(com.google.javascript.rhino.Node arg)
func (this_ *OptimizeParameters_UnusedParameterOptimizer) priv_UnusedParameterOptimizer_recordReplaceWithZero_4126(arg pkg_rhino.Node) {
	_ = arg
	if !(arg.IsNumber()) || (arg.GetDouble() != float64(int32(0))) {
		this_.toReplaceWithZero = append(this_.toReplaceWithZero, arg)
	}
}

// Java method: UnusedParameterOptimizer.removeArgAndFollowing(com.google.javascript.rhino.Node arg)
func (this_ *OptimizeParameters_UnusedParameterOptimizer) RemoveArgAndFollowing(arg pkg_rhino.Node) bool {
	_ = arg
	if arg != 0 {
		var removedAll bool = this_.RemoveArgAndFollowing(arg.GetNext())
		_ = removedAll
		if !(this_.outer_.astAnalyzer.MayHaveSideEffects(arg)) {
			this_.toRemove = append(this_.toRemove, arg)
			return removedAll
		}
		return false
	}
	return true
}

// Java method: UnusedParameterOptimizer.removeUnusedFunctionParameters(java.util.BitSet unremovable, com.google.javascript.rhino.Node param, int index)
func (this_ *OptimizeParameters_UnusedParameterOptimizer) RemoveUnusedFunctionParameters(unremovable *j2gort.Ext_java_util_BitSet, param pkg_rhino.Node, index int32) {
	_ = unremovable
	_ = param
	_ = index
	if param != 0 {
		this_.RemoveUnusedFunctionParameters(unremovable, param.GetNext(), (index + int32(1)))
		if !(unremovable.Get_1(index)) {
			if !(param.IsName()) {
				j2gort.CheckState(false)
			}
			this_.outer_.compiler.ReportChangeToEnclosingScope(param)
			param.Detach()
		}
	}
}

// Java method: UnusedParameterOptimizer.applyChanges()
func (this_ *OptimizeParameters_UnusedParameterOptimizer) ApplyChanges() {
	for _, n := range this_.toRemove {
		_ = n
		if OptimizeParameters_AlreadyRemoved(n) {
			continue
		}
		this_.outer_.compiler.ReportChangeToEnclosingScope(n)
		n.Detach()
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n, this_.outer_.compiler)
	}

	for _, n_2 := range this_.toReplaceWithZero {
		_ = n_2
		if !(!(n_2.IsNumber()) || (n_2.GetDouble() != float64(0.0))) {
			j2gort.CheckState(false)
		}
		if OptimizeParameters_AlreadyRemoved(n_2) {
			continue
		}
		this_.outer_.compiler.ReportChangeToEnclosingScope(n_2)
		n_2.ReplaceWith(pkg_rhino.IR_Number(float64(int32(0))).Srcref(n_2))
		NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(n_2, this_.outer_.compiler)
	}

}

// Java class: com.google.javascript.jscomp.OptimizeParameters.CandidateAnalysis
type OptimizeParameters_CandidateAnalysis struct {
	outer_                 *OptimizeParameters
	taggedTemplateLiterals []pkg_rhino.Node `gson:"taggedTemplateLiterals"`
	idHash_                atomic.Int32
	isSafeToOptimize       bool `gson:"isSafeToOptimize"`
}

func (this_ *OptimizeParameters_CandidateAnalysis) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeParameters_CandidateAnalysis) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeParameters_CandidateAnalysis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeParameters$CandidateAnalysis", "CandidateAnalysis"
}
func (this_ *OptimizeParameters_CandidateAnalysis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeParameters_CandidateAnalysis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeParameters.CandidateAnalysis", "com.google.javascript.jscomp.OptimizeParameters$CandidateAnalysis", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeParameters_CandidateAnalysis](v_)
	})
}
func (this_ *OptimizeParameters_CandidateAnalysis) SetOuter_(o *OptimizeParameters) { this_.outer_ = o }
func (this_ *OptimizeParameters_CandidateAnalysis) GetOuter_() *OptimizeParameters {
	return this_.outer_
}

func (this_ *OptimizeParameters_CandidateAnalysis) initFields_OptimizeParameters_CandidateAnalysis() {
	_ = this_

}

// Java constructor: CandidateAnalysis()
func NewOptimizeParameters_CandidateAnalysis(outer_ *OptimizeParameters) *OptimizeParameters_CandidateAnalysis {
	this_ := &OptimizeParameters_CandidateAnalysis{}
	this_.outer_ = outer_
	InitOptimizeParameters_CandidateAnalysis(this_)
	return this_
}

func InitOptimizeParameters_CandidateAnalysis(this_ *OptimizeParameters_CandidateAnalysis) {
	this_.initFields_OptimizeParameters_CandidateAnalysis()
}

// Java method: CandidateAnalysis.addTaggedTemplateLiteral(com.google.javascript.rhino.Node ttlNode)
func (this_ *OptimizeParameters_CandidateAnalysis) AddTaggedTemplateLiteral(ttlNode pkg_rhino.Node) {
	_ = ttlNode
	if this_.taggedTemplateLiterals == nil {
		this_.taggedTemplateLiterals = []pkg_rhino.Node{}
	}
	this_.taggedTemplateLiterals = append(this_.taggedTemplateLiterals, ttlNode)
}

// Java method: CandidateAnalysis.hasTaggedTemplateLiterals()
func (this_ *OptimizeParameters_CandidateAnalysis) HasTaggedTemplateLiterals() bool {
	return ((this_.taggedTemplateLiterals != nil) && len(this_.taggedTemplateLiterals) != 0)
}

// Java method: CandidateAnalysis.shouldConvertTaggedTemplateLiterals()
func (this_ *OptimizeParameters_CandidateAnalysis) ShouldConvertTaggedTemplateLiterals() bool {
	return (this_.HasTaggedTemplateLiterals() && this_.isSafeToOptimize)
}

// Java method: CandidateAnalysis.convertTaggedTemplateLiterals()
func (this_ *OptimizeParameters_CandidateAnalysis) ConvertTaggedTemplateLiterals() {
	if this_.taggedTemplateLiterals != nil {
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(this_.taggedTemplateLiterals))
			_ = size
			for ; i < size; i = i + 1 {
				this_.outer_.priv_OptimizeParameters_convertTaggedTemplateLiteralToNormalCall_4100(this_.taggedTemplateLiterals[i])
			}
		}
		this_.taggedTemplateLiterals = nil
	}
}

// Java method: CandidateAnalysis.isSafeToOptimize()
func (this_ *OptimizeParameters_CandidateAnalysis) IsSafeToOptimize() bool {
	return (this_.isSafeToOptimize && !(this_.HasTaggedTemplateLiterals()))
}

// Java class: com.google.javascript.jscomp.OptimizeParameters.ParameterList
type OptimizeParameters_ParameterList struct {
	items   []*OptimizeParameters_Parameter `gson:"items"`
	idHash_ atomic.Int32
}

func (this_ *OptimizeParameters_ParameterList) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeParameters_ParameterList) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeParameters_ParameterList) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeParameters$ParameterList", "ParameterList"
}
func (this_ *OptimizeParameters_ParameterList) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeParameters_ParameterList) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeParameters.ParameterList", "com.google.javascript.jscomp.OptimizeParameters$ParameterList", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeParameters_ParameterList](v_)
	})
}

func (this_ *OptimizeParameters_ParameterList) initFields_OptimizeParameters_ParameterList() {
	_ = this_
	this_.items = []*OptimizeParameters_Parameter{}
}

// Java constructor: ParameterList()
func NewOptimizeParameters_ParameterList() *OptimizeParameters_ParameterList {
	this_ := &OptimizeParameters_ParameterList{}
	InitOptimizeParameters_ParameterList(this_)
	return this_
}

func InitOptimizeParameters_ParameterList(this_ *OptimizeParameters_ParameterList) {
	this_.initFields_OptimizeParameters_ParameterList()
}

// Java method: ParameterList.size()
func (this_ *OptimizeParameters_ParameterList) Size() int32 {
	return int32(len(this_.items))
}

// Java method: ParameterList.get(int index)
func (this_ *OptimizeParameters_ParameterList) Get(index int32) *OptimizeParameters_Parameter {
	_ = index
	return this_.items[index]
}

// Java method: ParameterList.add(com.google.javascript.jscomp.OptimizeParameters.Parameter p)
func (this_ *OptimizeParameters_ParameterList) Add(p *OptimizeParameters_Parameter) {
	_ = p
	this_.items = append(this_.items, p)
}

// Java class: com.google.javascript.jscomp.OptimizeParameters.Parameter
type OptimizeParameters_Parameter struct {
	arg               pkg_rhino.Node `gson:"arg"`
	idHash_           atomic.Int32
	shouldRemove      bool `gson:"shouldRemove"`
	hasSideEffects    bool `gson:"hasSideEffects"`
	canBeSideEffected bool `gson:"canBeSideEffected"`
	mayBeUndefined    bool `gson:"mayBeUndefined"`
}

func (this_ *OptimizeParameters_Parameter) BindSelf_(s any) {
	_ = s
}
func (this_ *OptimizeParameters_Parameter) IsJavaNil_() bool { return this_ == nil }
func (*OptimizeParameters_Parameter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.OptimizeParameters$Parameter", "Parameter"
}
func (this_ *OptimizeParameters_Parameter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*OptimizeParameters_Parameter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.OptimizeParameters.Parameter", "com.google.javascript.jscomp.OptimizeParameters$Parameter", func(v_ any) bool {
		return j2gort.InstanceOf[*OptimizeParameters_Parameter](v_)
	})
}

func (this_ *OptimizeParameters_Parameter) initFields_OptimizeParameters_Parameter() {
	_ = this_
}

// Java constructor: Parameter(com.google.javascript.rhino.Node arg, boolean shouldRemove)
func NewOptimizeParameters_Parameter(arg pkg_rhino.Node, shouldRemove bool) *OptimizeParameters_Parameter {
	this_ := &OptimizeParameters_Parameter{}
	InitOptimizeParameters_Parameter(this_, arg, shouldRemove)
	return this_
}
func allocPooled_OptimizeParameters_Parameter(slab_ *j2gort.PooledSlab[OptimizeParameters_Parameter], arg pkg_rhino.Node, shouldRemove bool) *OptimizeParameters_Parameter {
	this_ := slab_.Alloc()
	InitOptimizeParameters_Parameter(this_, arg, shouldRemove)
	return this_
}

func InitOptimizeParameters_Parameter(this_ *OptimizeParameters_Parameter, arg pkg_rhino.Node, shouldRemove bool) {
	_ = arg
	_ = shouldRemove
	this_.initFields_OptimizeParameters_Parameter()
	this_.shouldRemove = shouldRemove
	this_.arg = arg
}

// Java method: Parameter.getArg()
func (this_ *OptimizeParameters_Parameter) GetArg() pkg_rhino.Node {
	return this_.arg
}

// Java method: Parameter.shouldRemove()
func (this_ *OptimizeParameters_Parameter) ShouldRemove() bool {
	return this_.shouldRemove
}

// Java method: Parameter.setShouldRemove(boolean value)
func (this_ *OptimizeParameters_Parameter) SetShouldRemove(value bool) {
	_ = value
	this_.shouldRemove = value
}

// Java method: Parameter.setHasSideEffects(boolean hasSideEffects)
func (this_ *OptimizeParameters_Parameter) SetHasSideEffects(hasSideEffects bool) {
	_ = hasSideEffects
	this_.hasSideEffects = hasSideEffects
}

// Java method: Parameter.hasSideEffects()
func (this_ *OptimizeParameters_Parameter) HasSideEffects() bool {
	return this_.hasSideEffects
}

// Java method: Parameter.setCanBeSideEffected(boolean canBeSideEffected)
func (this_ *OptimizeParameters_Parameter) SetCanBeSideEffected(canBeSideEffected bool) {
	_ = canBeSideEffected
	this_.canBeSideEffected = canBeSideEffected
}

// Java method: Parameter.canBeSideEffected()
func (this_ *OptimizeParameters_Parameter) CanBeSideEffected() bool {
	return this_.canBeSideEffected
}

// Java method: Parameter.setMayBeUndefined(boolean mayBeUndefined)
func (this_ *OptimizeParameters_Parameter) SetMayBeUndefined(mayBeUndefined bool) {
	_ = mayBeUndefined
	this_.mayBeUndefined = mayBeUndefined
}
func (this_ *OptimizeParameters) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_CandidateAnalysis.Release()
	this_.slab_ParameterList.Release()
	this_.slab_Parameter.Release()
	this_.slab_OptimizeCalls_NodeList.Release()
}
