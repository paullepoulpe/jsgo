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
// Port of com/google/javascript/jscomp/PureFunctionIdentifier.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_pure_function_identifier_1 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_pure_function_identifier_2 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_pure_function_identifier_3 func(any) string
var lambdaSingleton_pure_function_identifier_4 j2gort.Ext_java_util_function_Predicate
var lambdaSingleton_pure_function_identifier_5 j2gort.Ext_java_util_function_Predicate
var lambdaSingleton_pure_function_identifier_6 j2gort.Ext_java_util_function_Predicate
var lambdaSingleton_pure_function_identifier_7 j2gort.Ext_java_util_function_Function
var lambdaSingleton_pure_function_identifier_8 pkg_graph.FixedPointGraphTraversal_EdgeCallback[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]

func init() {
	lambdaSingleton_pure_function_identifier_1 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(mrp0_ any) int32 { return j2gort.CastNode[pkg_rhino.Node](mrp0_).GetCharno() })
	lambdaSingleton_pure_function_identifier_2 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(mrp0_ any) int32 { return j2gort.CastNode[pkg_rhino.Node](mrp0_).GetLineno() })
	lambdaSingleton_pure_function_identifier_3 = func(mrp0_ any) string { return j2gort.CastNode[pkg_rhino.Node](mrp0_).GetSourceFileName() }
	lambdaSingleton_pure_function_identifier_4 = j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_1 uint32) bool {
		lhs := pkg_rhino.Node(rawNode_1)
		_ = lhs
		return true
	})
	lambdaSingleton_pure_function_identifier_5 = j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_2 uint32) bool {
		lhs := pkg_rhino.Node(rawNode_2)
		_ = lhs
		return false
	})
	lambdaSingleton_pure_function_identifier_6 = j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_3 uint32) bool {
		lhs := pkg_rhino.Node(rawNode_3)
		_ = lhs
		var rhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(lhs)
		_ = rhs
		return ((rhs == 0) || NodeUtil_EvaluatesToLocalValue(rhs))
	})
	lambdaSingleton_pure_function_identifier_7 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.BoxString(PureFunctionIdentifier_FormatSourceLocation(j2gort.CastNode[pkg_rhino.Node](mrp0_)))
	})
	lambdaSingleton_pure_function_identifier_8 = pkg_graph.NewLambda_FixedPointGraphTraversal_EdgeCallback[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation](func(source *PureFunctionIdentifier_AmbiguatedFunctionSummary, edge *PureFunctionIdentifier_SideEffectPropagation, destination *PureFunctionIdentifier_AmbiguatedFunctionSummary) bool {
		_ = source
		_ = edge
		_ = destination
		return edge.Propagate(source, destination)
	})
}

// Java class: com.google.javascript.jscomp.PureFunctionIdentifier
type PureFunctionIdentifier struct {
	compiler                             AbstractCompiler                                                                                                                `gson:"compiler"`
	astAnalyzer                          *AstAnalyzer                                                                                                                    `gson:"astAnalyzer"`
	summariesByName                      j2gort.OrderedMap[string, *PureFunctionIdentifier_AmbiguatedFunctionSummary]                                                    `gson:"summariesByName"`
	summariesByPropName                  j2gort.OrderedMap[string, *PureFunctionIdentifier_AmbiguatedFunctionSummary]                                                    `gson:"summariesByPropName"`
	summariesForAllNamesOfFunctionByNode j2gort.OrderedMap[pkg_rhino.Node, *PureFunctionIdentifier_SummaryList]                                                          `gson:"summariesForAllNamesOfFunctionByNode"`
	artificiallyPureLiteralsForDebugging j2gort.OrderedSet[pkg_rhino.Node]                                                                                               `gson:"artificiallyPureLiteralsForDebugging"`
	allFunctionCalls                     *OptimizeCalls_NodeList                                                                                                         `gson:"allFunctionCalls"`
	scratchCallees                       *OptimizeCalls_NodeList                                                                                                         `gson:"scratchCallees"`
	scratchRvalueLeaves                  *OptimizeCalls_NodeList                                                                                                         `gson:"scratchRvalueLeaves"`
	scratchSummaries                     []*PureFunctionIdentifier_AmbiguatedFunctionSummary                                                                             `gson:"scratchSummaries"`
	reusableFlags                        *pkg_rhino.Node_SideEffectFlags                                                                                                 `gson:"reusableFlags"`
	reverseCallGraph                     pkg_graph.LinkedDirectedGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation] `gson:"reverseCallGraph"`
	unknownFunctionSummary               *PureFunctionIdentifier_AmbiguatedFunctionSummary                                                                               `gson:"unknownFunctionSummary"`
	idHash_                              atomic.Int32
	validateArtificialPurity             bool `gson:"validateArtificialPurity"`
	lastInvocationArgsMatchParams        bool `gson:"lastInvocationArgsMatchParams"`
	assumeGettersArePure                 bool `gson:"assumeGettersArePure"`
	hasProcessed                         bool `gson:"hasProcessed"`
	slab_SideEffectPropagation           j2gort.PooledSlab[PureFunctionIdentifier_SideEffectPropagation]
}

func (this_ *PureFunctionIdentifier) BindSelf_(s any) {
	_ = s
}
func (this_ *PureFunctionIdentifier) IsJavaNil_() bool { return this_ == nil }
func (*PureFunctionIdentifier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PureFunctionIdentifier", "PureFunctionIdentifier"
}
func (this_ *PureFunctionIdentifier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PureFunctionIdentifier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PureFunctionIdentifier", "com.google.javascript.jscomp.PureFunctionIdentifier", func(v_ any) bool {
		return j2gort.InstanceOf[*PureFunctionIdentifier](v_)
	})
}
func (*PureFunctionIdentifier) Implements_jscomp_OptimizeCalls_CallGraphCompilerPass() {}

func StaticGet_PureFunctionIdentifier_PROP_NAME_PREFIX() string { return "." }

var static_PureFunctionIdentifier_IMPLICIT_PURE_FN pkg_rhino.Node

func StaticGet_PureFunctionIdentifier_IMPLICIT_PURE_FN() pkg_rhino.Node {
	clinit_PureFunctionIdentifier()
	return static_PureFunctionIdentifier_IMPLICIT_PURE_FN
}
func StaticSet_PureFunctionIdentifier_IMPLICIT_PURE_FN(v pkg_rhino.Node) {
	clinit_PureFunctionIdentifier()
	static_PureFunctionIdentifier_IMPLICIT_PURE_FN = v
}

var static_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION *DiagnosticType

func StaticGet_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION() *DiagnosticType {
	clinit_PureFunctionIdentifier()
	return static_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION
}
func StaticSet_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION(v *DiagnosticType) {
	clinit_PureFunctionIdentifier()
	static_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION = v
}

var static_PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR j2gort.Ext_java_util_Comparator

func StaticGet_PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR() j2gort.Ext_java_util_Comparator {
	clinit_PureFunctionIdentifier()
	return static_PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR
}
func StaticSet_PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR(v j2gort.Ext_java_util_Comparator) {
	clinit_PureFunctionIdentifier()
	static_PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR = v
}

var static_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL j2gort.Ext_java_util_function_Predicate

func StaticGet_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL() j2gort.Ext_java_util_function_Predicate {
	clinit_PureFunctionIdentifier()
	return static_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL
}
func StaticSet_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL(v j2gort.Ext_java_util_function_Predicate) {
	clinit_PureFunctionIdentifier()
	static_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL = v
}

var static_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL j2gort.Ext_java_util_function_Predicate

func StaticGet_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL() j2gort.Ext_java_util_function_Predicate {
	clinit_PureFunctionIdentifier()
	return static_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL
}
func StaticSet_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL(v j2gort.Ext_java_util_function_Predicate) {
	clinit_PureFunctionIdentifier()
	static_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL = v
}

var static_PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE j2gort.Ext_java_util_function_Predicate

func StaticGet_PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE() j2gort.Ext_java_util_function_Predicate {
	clinit_PureFunctionIdentifier()
	return static_PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE
}
func StaticSet_PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE(v j2gort.Ext_java_util_function_Predicate) {
	clinit_PureFunctionIdentifier()
	static_PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE = v
}

var clinitDone_PureFunctionIdentifier uint32
var clinitState_PureFunctionIdentifier j2gort.ClinitState

func clinit_PureFunctionIdentifier() {
	if atomic.LoadUint32(&clinitDone_PureFunctionIdentifier) != 0 {
		return
	}
	clinitSlow_PureFunctionIdentifier()
}

func clinitSlow_PureFunctionIdentifier() {
	if !j2gort.EnterClinit(&clinitState_PureFunctionIdentifier, &clinitDone_PureFunctionIdentifier) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PureFunctionIdentifier, &clinitDone_PureFunctionIdentifier)
	static_PureFunctionIdentifier_IMPLICIT_PURE_FN = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{}), pkg_rhino.IR_Block_0())
	static_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION = DiagnosticType_Disabled("JSC_UNUSED_ARTIFICIAL_PURE_ANNOTATION", "Artificial @nosideeffects annotation cannot be enforced: found ambiguous definitions of {0}.\n\nSide-effectful definitions:\n  {1}")
	static_PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR = j2gort.ComparatorComparingStr(lambdaSingleton_pure_function_identifier_3).ThenComparingInt(lambdaSingleton_pure_function_identifier_2).ThenComparingInt(lambdaSingleton_pure_function_identifier_1)
	static_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL = lambdaSingleton_pure_function_identifier_4
	static_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL = lambdaSingleton_pure_function_identifier_5
	static_PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE = lambdaSingleton_pure_function_identifier_6
}

func (this_ *PureFunctionIdentifier) initFields_PureFunctionIdentifier() {
	_ = this_
	this_.summariesByName = j2gort.NewOrderedMap[string, *PureFunctionIdentifier_AmbiguatedFunctionSummary]()
	this_.summariesByPropName = j2gort.NewOrderedMap[string, *PureFunctionIdentifier_AmbiguatedFunctionSummary]()
	this_.summariesForAllNamesOfFunctionByNode = j2gort.NewOrderedMap[pkg_rhino.Node, *PureFunctionIdentifier_SummaryList]()
	this_.artificiallyPureLiteralsForDebugging = j2gort.NewOrderedSet[pkg_rhino.Node]()
	this_.allFunctionCalls = NewOptimizeCalls_NodeList_1(int32(64))
	this_.scratchCallees = NewOptimizeCalls_NodeList_1(int32(8))
	this_.scratchRvalueLeaves = NewOptimizeCalls_NodeList_1(int32(16))
	this_.scratchSummaries = make([]*PureFunctionIdentifier_AmbiguatedFunctionSummary, int32(8))

	this_.reusableFlags = pkg_rhino.NewNode_SideEffectFlags()
	this_.reverseCallGraph = pkg_graph.LinkedDirectedGraph_CreateWithoutAnnotations[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]()
	this_.unknownFunctionSummary = PureFunctionIdentifier_AmbiguatedFunctionSummary_CreateInGraph(j2gort.AsIface[pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]](this_.reverseCallGraph), "<unknown>").SetMutatesGlobalStateAndAllOtherFlags(pkg_rhino.Node(0))

}

// Java constructor: PureFunctionIdentifier(com.google.javascript.jscomp.AbstractCompiler compiler, boolean assumeGettersArePure, boolean validateArtificialPurity)
func NewPureFunctionIdentifier(compiler AbstractCompiler, assumeGettersArePure bool, validateArtificialPurity bool) *PureFunctionIdentifier {
	clinit_PureFunctionIdentifier()
	this_ := &PureFunctionIdentifier{}
	InitPureFunctionIdentifier(this_, compiler, assumeGettersArePure, validateArtificialPurity)
	return this_
}

func InitPureFunctionIdentifier(this_ *PureFunctionIdentifier, compiler AbstractCompiler, assumeGettersArePure bool, validateArtificialPurity bool) {
	_ = compiler
	_ = assumeGettersArePure
	_ = validateArtificialPurity
	this_.initFields_PureFunctionIdentifier()
	this_.compiler = j2gort.CheckNotNullIface(compiler)
	this_.assumeGettersArePure = assumeGettersArePure
	this_.astAnalyzer = compiler.GetAstAnalyzer()
	this_.validateArtificialPurity = validateArtificialPurity
}

// Java method: PureFunctionIdentifier.putFunctionSummary(com.google.javascript.rhino.Node fn, com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary summary)
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_putFunctionSummary_4566(fn pkg_rhino.Node, summary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	_ = fn
	_ = summary
	var list *PureFunctionIdentifier_SummaryList = this_.summariesForAllNamesOfFunctionByNode.Get(fn)
	_ = list
	if list == nil {
		list = NewPureFunctionIdentifier_SummaryList()
		this_.summariesForAllNamesOfFunctionByNode.Put(fn, list)
	}
	list.Add(summary)
}

// Java method: PureFunctionIdentifier.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root, com.google.javascript.jscomp.OptimizeCalls.ReferenceMap references)
func (this_ *PureFunctionIdentifier) Process_3(externs pkg_rhino.Node, root pkg_rhino.Node, references *OptimizeCalls_ReferenceMap) {
	_ = externs
	_ = root
	_ = references
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	if !(!(this_.hasProcessed)) {
		j2gort.CheckStateStr(false, "PureFunctionIdentifier::process may only be called once per instance.")
	}
	this_.hasProcessed = true
	this_.priv_PureFunctionIdentifier_populateDatastructuresForAnalysisTraversal_4571(references)
	NodeTraversal_Traverse(this_.compiler, externs, NewPureFunctionIdentifier_ExternFunctionAnnotationAnalyzer(this_))
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		if this_.reverseCallGraph != nil && !this_.reverseCallGraph.IsJavaNil_() {
			if r_, ok_ := any(this_.reverseCallGraph).(interface{ ReleaseSlabs_() }); ok_ {
				r_.ReleaseSlabs_()
			}
			this_.reverseCallGraph = nil
		}
		this_.summariesByName.Clear()
		this_.summariesForAllNamesOfFunctionByNode.Clear()
		sc_.ReleaseSlabs_()
		this_.slab_SideEffectPropagation.Release()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(NewPureFunctionIdentifier_FunctionBodyAnalyzer(this_))
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(root)
	}

	this_.priv_PureFunctionIdentifier_propagateSideEffects_4573()
	this_.priv_PureFunctionIdentifier_validateArtificialPurity_4567()
	this_.priv_PureFunctionIdentifier_markPureFunctionCalls_4574()
	this_.summariesByPropName.Clear()
}

// Java method: PureFunctionIdentifier.validateArtificialPurity()
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_validateArtificialPurity_4567() {
	if !(this_.validateArtificialPurity) {
		return
	}
	{
		cur_4 := this_.artificiallyPureLiteralsForDebugging.Iter()
		for cur_4.HasNext() {
			function := cur_4.
				Next()
			_ = function
			var summaries *PureFunctionIdentifier_SummaryList = this_.summariesForAllNamesOfFunctionByNode.Get(function)
			_ = summaries
			if summaries == nil {
				continue
			}
			{
				var i int32 = int32(0)
				_ = i
				var count int32 = summaries.Size()
				_ = count
				for ; i < count; i = i + 1 {
					var summary *PureFunctionIdentifier_AmbiguatedFunctionSummary = summaries.Get(i)
					_ = summary
					if summary.HasNoFlagsSet() {
						continue
					}
					if j2gort.CheckNotNullString(summary.name) == ".constructor" {
						continue
					}
					var sideEffectfulDefinitions string = j2gort.CastString(summary.impureFunctionReasonsForDebugging.Stream().Sorted_1(StaticGet_PureFunctionIdentifier_SOURCE_LOCATION_COMPARATOR()).Map(lambdaSingleton_pure_function_identifier_7).Collect_1(j2gort.Ext_java_util_stream_Collectors_Joining_1("\n  ")))
					_ = sideEffectfulDefinitions
					this_.compiler.Report_1(JSError_Make_3(function, StaticGet_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION(), []string{summary.name, sideEffectfulDefinitions}))
				}
			}
		}
	}
}

// Java static method: PureFunctionIdentifier.formatSourceLocation(com.google.javascript.rhino.Node node)
func PureFunctionIdentifier_FormatSourceLocation(node pkg_rhino.Node) string {
	clinit_PureFunctionIdentifier()
	_ = node
	return j2gort.Ext_java_lang_String_Format_2("%s:%s:%s", []any{j2gort.BoxString(node.GetSourceFileName()), j2gort.ValueOfInt(node.GetLineno()), j2gort.ValueOfInt(node.GetCharno())})
}

// Java static method: PureFunctionIdentifier.collectCallableLeavesInternal(com.google.javascript.rhino.Node expr, com.google.javascript.jscomp.OptimizeCalls.NodeList results, com.google.javascript.jscomp.AstAnalyzer astAnalyzer)
func PureFunctionIdentifier_CollectCallableLeavesInternal(expr pkg_rhino.Node, results *OptimizeCalls_NodeList, astAnalyzer *AstAnalyzer) bool {
	clinit_PureFunctionIdentifier()
	_ = expr
	_ = results
	_ = astAnalyzer
	switch expr.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */, int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(32) /* Token.NAME */ :
		results.Add(expr)
		return true
	case int32(110) /* Token.SUPER */ :
		{
			var clazz pkg_rhino.Node = j2gort.CheckNotNullNode(NodeUtil_GetEnclosingClass(expr))
			_ = clazz
			return PureFunctionIdentifier_CollectCallableLeavesInternal(clazz.GetSecondChild(), results, astAnalyzer)
		}
	case int32(105) /* Token.CLASS */ :
		{
			var members pkg_rhino.Node = expr.GetLastChild()
			_ = members
			if (members != 0) && members.IsClassMembers() {
				{
					var member pkg_rhino.Node = members.GetFirstChild()
					_ = member
					for ; member != 0; member = member.GetNext() {
						var initializer pkg_rhino.Node = member.GetLastChild()
						_ = initializer
						if ((initializer != 0) && NodeUtil_IsNonStaticClassFieldInitializer(initializer)) && astAnalyzer.MayHaveSideEffects(initializer) {
							return false
						}
					}
				}
			}
			var ctorDef pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(expr)
			_ = ctorDef
			if ctorDef != 0 {
				return PureFunctionIdentifier_CollectCallableLeavesInternal(ctorDef.GetOnlyChild(), results, astAnalyzer)
			} else if expr.GetSecondChild().IsEmpty() {
				results.Add(StaticGet_PureFunctionIdentifier_IMPLICIT_PURE_FN())
				return true
			} else {
				return PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetSecondChild(), results, astAnalyzer)
			}
		}
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */ :
		return (PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetFirstChild(), results, astAnalyzer) && PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetSecondChild(), results, astAnalyzer))
	case int32(50) /* Token.COMMA */, int32(51) /* Token.ASSIGN */ :
		return PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetSecondChild(), results, astAnalyzer)
	case int32(67) /* Token.HOOK */ :
		return (PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetSecondChild(), results, astAnalyzer) && PureFunctionIdentifier_CollectCallableLeavesInternal(expr.GetChildAtIndex(int32(2)), results, astAnalyzer))
	case int32(135) /* Token.NEW_TARGET */, int32(37) /* Token.THIS */ :
		fallthrough
	default:
		return false
	}
	panic("unreachable")
}

// Java static method: PureFunctionIdentifier.isDefinitelyRValue(com.google.javascript.rhino.Node rvalue)
func PureFunctionIdentifier_IsDefinitelyRValue(rvalue pkg_rhino.Node) bool {
	clinit_PureFunctionIdentifier()
	_ = rvalue
	var parent pkg_rhino.Node = rvalue.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(69) /* Token.AND */, int32(46) /* Token.ARRAYLIT */, int32(28) /* Token.CALL */, int32(70) /* Token.COALESCE */, int32(50) /* Token.COMMA */, int32(4) /* Token.EQ */, int32(27) /* Token.GETELEM */, int32(26) /* Token.GETPROP */, int32(67) /* Token.HOOK */, int32(45) /* Token.INSTANCEOF */, int32(23) /* Token.NEW */, int32(19) /* Token.NOT */, int32(32) /* Token.NAME */, int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(68) /* Token.OR */, int32(0) /* Token.RETURN */, int32(40) /* Token.SHEQ */, int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(25) /* Token.TYPEOF */, int32(114) /* Token.YIELD */ :
		return true
	case int32(76) /* Token.CASE */, int32(74) /* Token.IF */, int32(75) /* Token.SWITCH */, int32(78) /* Token.WHILE */ :
		return rvalue.IsFirstChildOf(parent)
	case int32(93) /* Token.EXPR_RESULT */ :
		return !(rvalue.IsFromExterns())
	case int32(51) /* Token.ASSIGN */, int32(105) /* Token.CLASS */ :
		return rvalue.IsSecondChildOf(parent)
	case int32(100) /* Token.STRING_KEY */ :
		return parent.GetParent().IsObjectLit()
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: PureFunctionIdentifier.appendSummariesForCalleeNode(com.google.javascript.rhino.Node callee, int outCount)
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_appendSummariesForCalleeNode_4568(callee pkg_rhino.Node, outCount int32) int32 {
	_ = callee
	_ = outCount
	if callee.IsFunction() {
		var summariesForFunction *PureFunctionIdentifier_SummaryList = this_.summariesForAllNamesOfFunctionByNode.Get(callee)
		_ = summariesForFunction
		if !((summariesForFunction != nil) && (summariesForFunction.Size() > int32(0))) {
			j2gort.CheckStateFmt(false, "Function missed during analysis: %s", j2gort.BoxNode(callee))
		}
		var addCount int32 = summariesForFunction.Size()
		_ = addCount
		this_.priv_PureFunctionIdentifier_ensureScratchSummariesCapacity_4569((outCount + addCount))
		var out []*PureFunctionIdentifier_AmbiguatedFunctionSummary = this_.scratchSummaries
		_ = out
		{
			var i int32 = int32(0)
			_ = i
			for ; i < addCount; i = i + 1 {
				{
					a_5 := out
					i_6 := func() int32 { old_7 := outCount; nw_8 := old_7 + 1; outCount = nw_8; return old_7 }()
					a_5[i_6] = summariesForFunction.Get(i)
				}
			}
		}
		return outCount
	}
	this_.priv_PureFunctionIdentifier_ensureScratchSummariesCapacity_4569((outCount + int32(1)))
	{
		a_9 := this_.scratchSummaries
		i_10 := outCount
		a_9[i_10] = this_.priv_PureFunctionIdentifier_summaryForReference_4575(callee)
	}
	return (outCount + int32(1))
}

// Java method: PureFunctionIdentifier.ensureScratchSummariesCapacity(int minCap)
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_ensureScratchSummariesCapacity_4569(minCap int32) {
	_ = minCap
	if minCap > int32(len(this_.scratchSummaries)) {
		this_.scratchSummaries = func() []*PureFunctionIdentifier_AmbiguatedFunctionSummary {
			arrSrc_11 := this_.scratchSummaries
			arrLen_12 := max((int32(len(this_.scratchSummaries)) * int32(2)), minCap)
			j2gort.CheckNotNullSlice(arrSrc_11)
			if arrLen_12 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_13 := make([]*PureFunctionIdentifier_AmbiguatedFunctionSummary, arrLen_12)
			copy(arrOut_13, arrSrc_11)
			return arrOut_13
		}()
	}
}

// Java method: PureFunctionIdentifier.collectSummariesForCallee(com.google.javascript.rhino.Node invocation)
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_collectSummariesForCallee_4570(invocation pkg_rhino.Node) int32 {
	_ = invocation
	if !(NodeUtil_IsInvocation(invocation)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(invocation))
	}
	var cacheCall *CodingConvention_Cache = this_.compiler.GetCodingConvention().DescribeCachingCall(invocation)
	_ = cacheCall
	if cacheCall != nil {
		this_.lastInvocationArgsMatchParams = false
		var callees *OptimizeCalls_NodeList = this_.scratchCallees
		_ = callees
		callees.Clear()
		if !(PureFunctionIdentifier_CollectCallableLeavesInternal(cacheCall.valueFn, callees, this_.astAnalyzer)) {
			panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException")))
		}
		if cacheCall.keyFn != 0 {
			var beforeKeyCount int32 = callees.Size()
			_ = beforeKeyCount
			if !(PureFunctionIdentifier_CollectCallableLeavesInternal(cacheCall.keyFn, callees, this_.astAnalyzer)) {
				panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException")))
			}
			if !(callees.Size() >= beforeKeyCount) {
				j2gort.CheckState(false)
			}
		}
		if !(!(callees.IsEmpty())) {
			j2gort.CheckStateStr(false, "Unexpected empty callees for valid result")
		}
		var outCount int32 = int32(0)
		_ = outCount
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = callees.Size()
			_ = len_
			for ; i < len_; i = i + 1 {
				outCount = this_.priv_PureFunctionIdentifier_appendSummariesForCalleeNode_4568(callees.Get(i), outCount)
			}
		}
		return outCount
	}
	if PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation) {
		this_.lastInvocationArgsMatchParams = false
		return this_.priv_PureFunctionIdentifier_appendSummariesForCalleeNode_4568(invocation.GetFirstFirstChild(), int32(0))
	}
	this_.lastInvocationArgsMatchParams = !(invocation.IsTaggedTemplateLit())
	var target pkg_rhino.Node = invocation.GetFirstChild()
	_ = target
	switch target.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */, int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(32) /* Token.NAME */ :
		return this_.priv_PureFunctionIdentifier_appendSummariesForCalleeNode_4568(target, int32(0))
	default:
	}
	var callees_2 *OptimizeCalls_NodeList = this_.scratchCallees
	_ = callees_2
	callees_2.Clear()
	if !(PureFunctionIdentifier_CollectCallableLeavesInternal(target, callees_2, this_.astAnalyzer)) {
		{
			a_14 := this_.scratchSummaries
			i_15 := int32(0)
			a_14[i_15] = this_.unknownFunctionSummary
		}
		return int32(1)
	}
	if !(!(callees_2.IsEmpty())) {
		j2gort.CheckStateStr(false, "Unexpected empty callees for valid result")
	}
	var outCount_2 int32 = int32(0)
	_ = outCount_2
	{
		var i_2 int32 = int32(0)
		_ = i_2
		var len__2 int32 = callees_2.Size()
		_ = len__2
		for ; i_2 < len__2; i_2 = i_2 + 1 {
			outCount_2 = this_.priv_PureFunctionIdentifier_appendSummariesForCalleeNode_4568(callees_2.Get(i_2), outCount_2)
		}
	}
	return outCount_2
}

// Java method: PureFunctionIdentifier.populateDatastructuresForAnalysisTraversal(com.google.javascript.jscomp.OptimizeCalls.ReferenceMap referenceMap)
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_populateDatastructuresForAnalysisTraversal_4571(referenceMap *OptimizeCalls_ReferenceMap) {
	_ = referenceMap
	{
		cur_16 := referenceMap.names.IterEntries()
		for cur_16.NextEntry() {
			var refs *OptimizeCalls_NodeList = cur_16.Value()
			_ = refs
			if !(refs.IsEmpty()) {
				var name string = cur_16.Key()
				_ = name
				if !(!(j2gort.CheckNotNullString(name) == "")) {
					j2gort.CheckState(false)
				}
				this_.summariesByName.Put(name, PureFunctionIdentifier_AmbiguatedFunctionSummary_CreateInGraph(j2gort.AsIface[pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]](this_.reverseCallGraph), name))
			}
		}
	}
	{
		cur_17 := referenceMap.props.IterEntries()
		for cur_17.NextEntry() {
			var refs_2 *OptimizeCalls_NodeList = cur_17.Value()
			_ = refs_2
			if !(refs_2.IsEmpty()) {
				var propName string = cur_17.Key()
				_ = propName
				if !(!(j2gort.CheckNotNullString(propName) == "")) {
					j2gort.CheckState(false)
				}
				this_.summariesByPropName.Put(propName, PureFunctionIdentifier_AmbiguatedFunctionSummary_CreateInGraph(j2gort.AsIface[pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]](this_.reverseCallGraph), ("."+j2gort.StrValString(propName))))
			}
		}
	}
	{
		cur_18 := referenceMap.names.IterEntries()
		for cur_18.NextEntry() {
			var refs_3 *OptimizeCalls_NodeList = cur_18.Value()
			_ = refs_3
			if !(refs_3.IsEmpty()) {
				this_.priv_PureFunctionIdentifier_populateFunctionDefinitions_4572(j2gort.CheckNotNullPtr(this_.summariesByName.Get(cur_18.Key())), refs_3)
			}
		}
	}
	{
		cur_19 := referenceMap.props.IterEntries()
		for cur_19.NextEntry() {
			var refs_4 *OptimizeCalls_NodeList = cur_19.Value()
			_ = refs_4
			if !(refs_4.IsEmpty()) {
				this_.priv_PureFunctionIdentifier_populateFunctionDefinitions_4572(j2gort.CheckNotNullPtr(this_.summariesByPropName.Get(cur_19.Key())), refs_4)
			}
		}
	}
	this_.priv_PureFunctionIdentifier_putFunctionSummary_4566(StaticGet_PureFunctionIdentifier_IMPLICIT_PURE_FN(), PureFunctionIdentifier_AmbiguatedFunctionSummary_CreateInGraph(j2gort.AsIface[pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]](this_.reverseCallGraph), "<implicit pure class ctor"))
}

// Java method: PureFunctionIdentifier.populateFunctionDefinitions(com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary summaryForName, com.google.javascript.jscomp.OptimizeCalls.NodeList references)
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_populateFunctionDefinitions_4572(summaryForName *PureFunctionIdentifier_AmbiguatedFunctionSummary, references *OptimizeCalls_NodeList) {
	_ = summaryForName
	_ = references
	var invalid bool = false
	_ = invalid
	var rvalueLeaves *OptimizeCalls_NodeList = this_.scratchRvalueLeaves
	_ = rvalueLeaves
	rvalueLeaves.Clear()
	{
		var i int32 = int32(0)
		_ = i
		var refCount int32 = references.Size()
		_ = refCount
		for ; i < refCount; i = i + 1 {
			var reference pkg_rhino.Node = references.Get(i)
			_ = reference
			if !(PureFunctionIdentifier_IsDefinitelyRValue(reference)) {
				var rvalue pkg_rhino.Node = NodeUtil_GetRValueOfLValue(reference)
				_ = rvalue
				if rvalue == 0 {
					invalid = true
					break
				} else {
					var startIdx int32 = rvalueLeaves.Size()
					_ = startIdx
					if !(PureFunctionIdentifier_CollectCallableLeavesInternal(rvalue, rvalueLeaves, this_.astAnalyzer)) {
						invalid = true
						break
					}
					if this_.validateArtificialPurity && !(reference.IsFromExterns()) {
						{
							var cIdx int32 = startIdx
							_ = cIdx
							var endIdx int32 = rvalueLeaves.Size()
							_ = endIdx
							for ; cIdx < endIdx; cIdx = cIdx + 1 {
								var callable pkg_rhino.Node = rvalueLeaves.Get(cIdx)
								_ = callable
								var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(callable)
								_ = jsdoc
								var isArtificiallyPure bool = ((jsdoc != nil) && jsdoc.IsNoSideEffects())
								_ = isArtificiallyPure
								if isArtificiallyPure {
									this_.artificiallyPureLiteralsForDebugging.Add(callable)
									summaryForName.SetCollectImpureDebuggingReason()
								}
							}
						}
					}
				}
			}
		}
	}
	if rvalueLeaves.IsEmpty() || invalid {
		summaryForName.SetMutatesGlobalStateAndAllOtherFlags(pkg_rhino.Node(0))
	} else {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var len_ int32 = rvalueLeaves.Size()
			_ = len_
			for ; i_2 < len_; i_2 = i_2 + 1 {
				var rvalue_2 pkg_rhino.Node = rvalueLeaves.Get(i_2)
				_ = rvalue_2
				if rvalue_2.IsFunction() {
					this_.priv_PureFunctionIdentifier_putFunctionSummary_4566(rvalue_2, summaryForName)
				} else if NodeUtil_IsUndefined(rvalue_2) {
				} else {
					var rvalueSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary = this_.priv_PureFunctionIdentifier_summaryForReference_4575(rvalue_2)
					_ = rvalueSummary
					this_.reverseCallGraph.Connect_3_DiGraph_DiGraphNode_any_any__any_DiGraph_DiGraphNode_any_any_(rvalueSummary.graphNode, allocPooledVia_PureFunctionIdentifier_SideEffectPropagation_ForAlias(&this_.slab_SideEffectPropagation), summaryForName.graphNode)
				}
			}
		}
	}
}

// Java method: PureFunctionIdentifier.propagateSideEffects()
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_propagateSideEffects_4573() {
	pkg_graph.FixedPointGraphTraversal_NewTraversal[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation](lambdaSingleton_pure_function_identifier_8).ComputeFixedPoint_1(j2gort.AsIface[pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]](this_.reverseCallGraph))
}

// Java method: PureFunctionIdentifier.markPureFunctionCalls()
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_markPureFunctionCalls_4574() {
	var flags *pkg_rhino.Node_SideEffectFlags = this_.reusableFlags
	_ = flags
	{
		var callIdx int32 = int32(0)
		_ = callIdx
		var callCount int32 = this_.allFunctionCalls.Size()
		_ = callCount
		for ; callIdx < callCount; callIdx = callIdx + 1 {
			var callNode pkg_rhino.Node = this_.allFunctionCalls.Get(callIdx)
			_ = callNode
			var summaryCount int32 = this_.priv_PureFunctionIdentifier_collectSummariesForCallee_4570(callNode)
			_ = summaryCount
			var calleeSummaries []*PureFunctionIdentifier_AmbiguatedFunctionSummary = this_.scratchSummaries
			_ = calleeSummaries
			if summaryCount == int32(0) {
				flags.SetAllFlags()
			} else {
				flags.ClearAllFlags()
				var isCallOrTemplate bool = PureFunctionIdentifier_IsCallOrTaggedTemplateLit(callNode)
				_ = isCallOrTemplate
				var viaCallOrApply bool = (isCallOrTemplate && PureFunctionIdentifier_IsInvocationViaCallOrApply(callNode))
				_ = viaCallOrApply
				{
					var s int32 = int32(0)
					_ = s
					for ; s < summaryCount; s = s + 1 {
						var calleeSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary = j2gort.CheckNotNullPtr(calleeSummaries[s])
						_ = calleeSummary
						if calleeSummary.MutatesGlobalState() {
							flags.SetMutatesGlobalState()
						}
						if calleeSummary.MutatesArguments() {
							flags.SetMutatesArguments()
						}
						if calleeSummary.FunctionThrows() {
							flags.SetThrows()
						}
						if isCallOrTemplate && calleeSummary.MutatesThis() {
							if viaCallOrApply {
								flags.SetMutatesArguments()
							} else {
								flags.SetMutatesThis()
							}
						}
					}
				}
			}
			if callNode.GetFirstChild().IsSuper() {
				flags.SetMutatesThis()
			}
			if PureFunctionIdentifier_IsCallOrTaggedTemplateLit(callNode) {
				if !(this_.astAnalyzer.FunctionCallHasSideEffects(callNode)) {
					flags.ClearAllFlags()
				}
			} else if callNode.IsNew() {
				if !(this_.astAnalyzer.ConstructorCallHasSideEffects(callNode)) {
					flags.ClearAllFlags()
				}
			}
			if callNode.GetSideEffectFlags() != flags.ValueOf() {
				callNode.SetSideEffectFlags_1_Node_SideEffectFlags(flags)
				this_.compiler.ReportChangeToEnclosingScope(callNode)
			}
		}
	}
}

// Java static method: PureFunctionIdentifier.getParamIndex(com.google.javascript.jscomp.Var param)
func PureFunctionIdentifier_GetParamIndex(param *Var) int32 {
	clinit_PureFunctionIdentifier()
	_ = param
	var node pkg_rhino.Node
	_ = node
	if param.IsDefaultParam() {
		node = param.GetParentNode()
	} else {
		node = param.GetNameNode()
	}
	if !(node.GetParent().IsParamList()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
	}
	return node.GetParent().GetIndexOfChild(node)
}

// Java static method: PureFunctionIdentifier.isInvocationViaCallOrApply(com.google.javascript.rhino.Node callSite)
func PureFunctionIdentifier_IsInvocationViaCallOrApply(callSite pkg_rhino.Node) bool {
	clinit_PureFunctionIdentifier()
	_ = callSite
	var receiver pkg_rhino.Node = callSite.GetFirstFirstChild()
	_ = receiver
	if (receiver == 0) || !((receiver.IsName() || receiver.IsGetProp()) || receiver.IsOptChainGetProp()) {
		return false
	}
	return (NodeUtil_IsFunctionObjectCall(callSite) || NodeUtil_IsFunctionObjectApply(callSite))
}

// Java static method: PureFunctionIdentifier.isCallOrTaggedTemplateLit(com.google.javascript.rhino.Node invocation)
func PureFunctionIdentifier_IsCallOrTaggedTemplateLit(invocation pkg_rhino.Node) bool {
	clinit_PureFunctionIdentifier()
	_ = invocation
	return ((invocation.IsCall() || invocation.IsOptChainCall()) || invocation.IsTaggedTemplateLit())
}

// Java method: PureFunctionIdentifier.summaryForReference(com.google.javascript.rhino.Node nameRef)
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_summaryForReference_4575(nameRef pkg_rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = nameRef
	switch nameRef.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return this_.summariesByName.GetOrDefault(nameRef.GetString(), this_.unknownFunctionSummary)
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return this_.summariesByPropName.GetOrDefault(nameRef.GetString(), this_.unknownFunctionSummary)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected name reference: "+j2gort.StrVal(nameRef))))
	}
	panic("unreachable switch expression")
}

// Java method: PureFunctionIdentifier.getPropertyKind(java.lang.String name)
func (this_ *PureFunctionIdentifier) priv_PureFunctionIdentifier_getPropertyKind_4576(name string) *AccessorSummary_PropertyAccessKind {
	_ = name
	if this_.assumeGettersArePure {
		return StaticGet_AccessorSummary_PropertyAccessKind_NORMAL()
	} else {
		return this_.compiler.GetAccessorSummary().GetKind(name)
	}
}

// Java class: com.google.javascript.jscomp.PureFunctionIdentifier.SummaryList
type PureFunctionIdentifier_SummaryList struct {
	first   *PureFunctionIdentifier_AmbiguatedFunctionSummary   `gson:"first"`
	rest    []*PureFunctionIdentifier_AmbiguatedFunctionSummary `gson:"rest"`
	size    int32                                               `gson:"size"`
	idHash_ atomic.Int32
}

func (this_ *PureFunctionIdentifier_SummaryList) BindSelf_(s any) {
	_ = s
}
func (this_ *PureFunctionIdentifier_SummaryList) IsJavaNil_() bool { return this_ == nil }
func (*PureFunctionIdentifier_SummaryList) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PureFunctionIdentifier$SummaryList", "SummaryList"
}
func (this_ *PureFunctionIdentifier_SummaryList) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PureFunctionIdentifier_SummaryList) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PureFunctionIdentifier.SummaryList", "com.google.javascript.jscomp.PureFunctionIdentifier$SummaryList", func(v_ any) bool {
		return j2gort.InstanceOf[*PureFunctionIdentifier_SummaryList](v_)
	})
}

func (this_ *PureFunctionIdentifier_SummaryList) initFields_PureFunctionIdentifier_SummaryList() {
	_ = this_
}

// Java constructor: SummaryList()
func NewPureFunctionIdentifier_SummaryList() *PureFunctionIdentifier_SummaryList {
	this_ := &PureFunctionIdentifier_SummaryList{}
	InitPureFunctionIdentifier_SummaryList(this_)
	return this_
}

func InitPureFunctionIdentifier_SummaryList(this_ *PureFunctionIdentifier_SummaryList) {
	this_.initFields_PureFunctionIdentifier_SummaryList()
}

// Java method: SummaryList.add(com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary summary)
func (this_ *PureFunctionIdentifier_SummaryList) Add(summary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	_ = summary
	var s int32 = this_.size
	_ = s
	if s == int32(0) {
		this_.first = summary
	} else {
		var r []*PureFunctionIdentifier_AmbiguatedFunctionSummary = this_.rest
		_ = r
		if r == nil {
			r = make([]*PureFunctionIdentifier_AmbiguatedFunctionSummary, int32(4))
			this_.rest = r
		} else if (s - int32(1)) == int32(len(r)) {
			r = func() []*PureFunctionIdentifier_AmbiguatedFunctionSummary {
				arrSrc_20 := r
				arrLen_21 := (int32(len(r)) * int32(2))
				j2gort.CheckNotNullSlice(arrSrc_20)
				if arrLen_21 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_22 := make([]*PureFunctionIdentifier_AmbiguatedFunctionSummary, arrLen_21)
				copy(arrOut_22, arrSrc_20)
				return arrOut_22
			}()
			this_.rest = r
		}
		{
			a_23 := r
			i_24 := (s - int32(1))
			a_23[i_24] = summary
		}
	}
	this_.size = (s + int32(1))
}

// Java method: SummaryList.size()
func (this_ *PureFunctionIdentifier_SummaryList) Size() int32 {
	return this_.size
}

// Java method: SummaryList.get(int index)
func (this_ *PureFunctionIdentifier_SummaryList) Get(index int32) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = index
	if index == int32(0) {
		return this_.first
	} else {
		return this_.rest[(index - int32(1))]
	}
}

// Java class: com.google.javascript.jscomp.PureFunctionIdentifier.ExternFunctionAnnotationAnalyzer
type PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer struct {
	outer_  *PureFunctionIdentifier
	idHash_ atomic.Int32
}

func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) BindSelf_(s any) {
	_ = s
}
func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) IsJavaNil_() bool {
	return this_ == nil
}
func (*PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PureFunctionIdentifier$ExternFunctionAnnotationAnalyzer", "ExternFunctionAnnotationAnalyzer"
}
func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PureFunctionIdentifier.ExternFunctionAnnotationAnalyzer", "com.google.javascript.jscomp.PureFunctionIdentifier$ExternFunctionAnnotationAnalyzer", func(v_ any) bool {
		return j2gort.InstanceOf[*PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer](v_)
	})
}
func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) SetOuter_(o *PureFunctionIdentifier) {
	this_.outer_ = o
}
func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) GetOuter_() *PureFunctionIdentifier {
	return this_.outer_
}
func (*PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) initFields_PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer() {
	_ = this_
}

// Java constructor: ExternFunctionAnnotationAnalyzer()
func NewPureFunctionIdentifier_ExternFunctionAnnotationAnalyzer(outer_ *PureFunctionIdentifier) *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer {
	this_ := &PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer{}
	this_.outer_ = outer_
	InitPureFunctionIdentifier_ExternFunctionAnnotationAnalyzer(this_)
	return this_
}

func InitPureFunctionIdentifier_ExternFunctionAnnotationAnalyzer(this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) {
	this_.initFields_PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer()
}

// Java method: ExternFunctionAnnotationAnalyzer.shouldTraverse(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) ShouldTraverse_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = traversal
	_ = node
	_ = parent
	return true
}

// Java method: ExternFunctionAnnotationAnalyzer.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
	if !(node.IsFunction()) {
		return
	}
	var summaries *PureFunctionIdentifier_SummaryList = this_.outer_.summariesForAllNamesOfFunctionByNode.Get(node)
	_ = summaries
	if summaries != nil {
		{
			var i int32 = int32(0)
			_ = i
			var count int32 = summaries.Size()
			_ = count
			for ; i < count; i = i + 1 {
				this_.priv_ExternFunctionAnnotationAnalyzer_updateSideEffectsForExternFunction_4580(node, summaries.Get(i))
			}
		}
	}
}

// Java method: ExternFunctionAnnotationAnalyzer.updateSideEffectsForExternFunction(com.google.javascript.rhino.Node externFunction, com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary summary)
func (this_ *PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) priv_ExternFunctionAnnotationAnalyzer_updateSideEffectsForExternFunction_4580(externFunction pkg_rhino.Node, summary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	_ = externFunction
	_ = summary
	if !(externFunction.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	if !(externFunction.IsFromExterns()) {
		j2gort.CheckArgument(false)
	}
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(externFunction)
	_ = info
	if info == nil {
		summary.SetMutatesGlobalStateAndAllOtherFlags(externFunction)
		return
	}
	if info.ModifiesThis() {
		summary.SetMutatesThis(externFunction)
	}
	if info.HasSideEffectsArgumentsAnnotation() {
		summary.SetMutatesArguments(externFunction)
	}
	if !(info.GetThrowsAnnotations().IsEmpty()) {
		summary.SetThrows(externFunction)
	}
	if !(info.IsNoSideEffects()) && summary.HasNoFlagsSet() {
		summary.SetMutatesGlobalStateAndAllOtherFlags(externFunction)
	}
}

// Java class: com.google.javascript.jscomp.PureFunctionIdentifier.FunctionBodyAnalyzer
type PureFunctionIdentifier_FunctionBodyAnalyzer struct {
	outer_             *PureFunctionIdentifier
	functionScopeStack j2gort.Deque[*PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry] `gson:"functionScopeStack"`
	idHash_            atomic.Int32
}

func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) BindSelf_(s any) {
	_ = s
}
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) IsJavaNil_() bool { return this_ == nil }
func (*PureFunctionIdentifier_FunctionBodyAnalyzer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PureFunctionIdentifier$FunctionBodyAnalyzer", "FunctionBodyAnalyzer"
}
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PureFunctionIdentifier_FunctionBodyAnalyzer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PureFunctionIdentifier.FunctionBodyAnalyzer", "com.google.javascript.jscomp.PureFunctionIdentifier$FunctionBodyAnalyzer", func(v_ any) bool {
		return j2gort.InstanceOf[*PureFunctionIdentifier_FunctionBodyAnalyzer](v_)
	})
}
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) SetOuter_(o *PureFunctionIdentifier) {
	this_.outer_ = o
}
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) GetOuter_() *PureFunctionIdentifier {
	return this_.outer_
}
func (*PureFunctionIdentifier_FunctionBodyAnalyzer) Implements_jscomp_NodeTraversal_ScopedCallback() {
}
func (*PureFunctionIdentifier_FunctionBodyAnalyzer) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) initFields_PureFunctionIdentifier_FunctionBodyAnalyzer() {
	_ = this_
	this_.functionScopeStack = j2gort.DequeFromCol[*PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry](j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.ImmListOfPtr(NewPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(pkg_rhino.Node(0), false))))
}

// Java constructor: FunctionBodyAnalyzer()
func NewPureFunctionIdentifier_FunctionBodyAnalyzer(outer_ *PureFunctionIdentifier) *PureFunctionIdentifier_FunctionBodyAnalyzer {
	this_ := &PureFunctionIdentifier_FunctionBodyAnalyzer{}
	this_.outer_ = outer_
	InitPureFunctionIdentifier_FunctionBodyAnalyzer(this_)
	return this_
}

func InitPureFunctionIdentifier_FunctionBodyAnalyzer(this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) {
	this_.initFields_PureFunctionIdentifier_FunctionBodyAnalyzer()
}

// Java method: FunctionBodyAnalyzer.shouldTraverse(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) ShouldTraverse_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = traversal
	_ = node
	_ = parent
	this_.priv_FunctionBodyAnalyzer_addToCatchDepthIfTryBlock_4588(node, int32(1))
	return true
}

// Java method: FunctionBodyAnalyzer.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
	this_.priv_FunctionBodyAnalyzer_addToCatchDepthIfTryBlock_4588(node, int32(-1))
	if NodeUtil_IsInvocation(node) {
		this_.outer_.allFunctionCalls.Add(node)
	}
	var enclosingFunction *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry = this_.functionScopeStack.GetLast()
	_ = enclosingFunction
	var root pkg_rhino.Node = enclosingFunction.root
	_ = root
	if (root == 0) || enclosingFunction.isArtificiallyPure {
		return
	}
	var summaries *PureFunctionIdentifier_SummaryList = this_.outer_.summariesForAllNamesOfFunctionByNode.Get(root)
	_ = summaries
	if summaries != nil {
		{
			var i int32 = int32(0)
			_ = i
			var count int32 = summaries.Size()
			_ = count
			for ; i < count; i = i + 1 {
				this_.priv_FunctionBodyAnalyzer_updateSideEffectsForNode_4581(j2gort.CheckNotNullPtr(summaries.Get(i)), traversal, node)
			}
		}
	}
}

// Java method: FunctionBodyAnalyzer.updateSideEffectsForNode(com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary encloserSummary, com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) priv_FunctionBodyAnalyzer_updateSideEffectsForNode_4581(encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary, traversal *NodeTraversal, node pkg_rhino.Node) {
	_ = encloserSummary
	_ = traversal
	_ = node
	if encloserSummary.MutatesGlobalState() {
		return
	}
	{
		var lhs pkg_rhino.Node = pkg_rhino.Node(0)
		_ = lhs
		var rhsLocality j2gort.Ext_java_util_function_Predicate = nil
		_ = rhsLocality
		switch node.TokenSwitchOrd_() {
		case int32(51) /* Token.ASSIGN */ :
			lhs = node.GetFirstChild()
			rhsLocality = func() j2gort.Ext_java_util_function_Predicate {
				if lhs.IsDestructuringPattern() {
					return StaticGet_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL()
				}
				return StaticGet_PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE()
			}()
			NodeUtil_VisitLhsNodesInNode(node, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_25 any) {
				lhsNode := j2gort.CastNode[pkg_rhino.Node](lp_25)
				_ = lhsNode
				this_.priv_FunctionBodyAnalyzer_visitLhsNode_4585(encloserSummary, traversal.GetScope(), lhsNode, rhsLocality)
			}))
		case int32(71) /* Token.INC */, int32(72) /* Token.DEC */, int32(24) /* Token.DELPROP */ :
			this_.priv_FunctionBodyAnalyzer_visitLhsNode_4585(encloserSummary, traversal.GetScope(), node.GetOnlyChild(), StaticGet_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL())
		case int32(113) /* Token.FOR_AWAIT_OF */ :
			this_.priv_FunctionBodyAnalyzer_deprecatedSetSideEffectsForControlLoss_4583(encloserSummary)
			fallthrough
		case int32(112) /* Token.FOR_OF */ :
			NodeUtil_VisitLhsNodesInNode(node, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_26 any) {
				lhsNode := j2gort.CastNode[pkg_rhino.Node](lp_26)
				_ = lhsNode
				this_.priv_FunctionBodyAnalyzer_visitLhsNode_4585(encloserSummary, traversal.GetScope(), lhsNode, StaticGet_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL())
			}))
			this_.priv_FunctionBodyAnalyzer_checkIteratesImpureIterable_4582(node, encloserSummary)
		case int32(81) /* Token.FOR_IN */ :
			NodeUtil_VisitLhsNodesInNode(node, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_27 any) {
				lhsNode := j2gort.CastNode[pkg_rhino.Node](lp_27)
				_ = lhsNode
				this_.priv_FunctionBodyAnalyzer_visitLhsNode_4585(encloserSummary, traversal.GetScope(), lhsNode, StaticGet_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL())
			}))
		case int32(31) /* Token.OPTCHAIN_CALL */, int32(28) /* Token.CALL */, int32(23) /* Token.NEW */, int32(130) /* Token.TAGGED_TEMPLATELIT */ :
			this_.priv_FunctionBodyAnalyzer_visitCall_4586(encloserSummary, traversal.GetScope(), node)
		case int32(104) /* Token.DESTRUCTURING_LHS */ :
			if NodeUtil_IsAnyFor(node.GetParent()) {
				break
			}
			NodeUtil_VisitLhsNodesInNode(node.GetParent(), j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_28 any) {
				lhsNode := j2gort.CastNode[pkg_rhino.Node](lp_28)
				_ = lhsNode
				this_.priv_FunctionBodyAnalyzer_visitLhsNode_4585(encloserSummary, traversal.GetScope(), lhsNode, StaticGet_PureFunctionIdentifier_RHS_IS_NEVER_LOCAL())
			}))
		case int32(32) /* Token.NAME */ :
			if NodeUtil_IsNameDeclaration(node.GetParent()) {
				var value pkg_rhino.Node = node.GetFirstChild()
				_ = value
				if (value != 0) && !(NodeUtil_EvaluatesToLocalValue(value)) {
					var scope *Scope = traversal.GetScope()
					_ = scope
					var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(node.GetString()))
					_ = var_
					this_.functionScopeStack.GetLast().skiplistedVars.Add(var_)
				}
			}
		case int32(43) /* Token.THROW */ :
			this_.priv_FunctionBodyAnalyzer_recordThrowsBasedOnContext_4584(encloserSummary)
		case int32(114) /* Token.YIELD */ :
			this_.priv_FunctionBodyAnalyzer_checkIteratesImpureIterable_4582(node, encloserSummary)
			this_.priv_FunctionBodyAnalyzer_deprecatedSetSideEffectsForControlLoss_4583(encloserSummary)
		case int32(115) /* Token.AWAIT */ :
			this_.priv_FunctionBodyAnalyzer_deprecatedSetSideEffectsForControlLoss_4583(encloserSummary)
		case int32(126) /* Token.OBJECT_REST */, int32(128) /* Token.OBJECT_SPREAD */ :
			if !(this_.outer_.assumeGettersArePure) {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
			}
		case int32(125) /* Token.ITER_REST */, int32(127) /* Token.ITER_SPREAD */ :
			this_.priv_FunctionBodyAnalyzer_checkIteratesImpureIterable_4582(node, encloserSummary)
		case int32(100) /* Token.STRING_KEY */ :
			if node.GetParent().IsObjectPattern() {
				if this_.outer_.priv_PureFunctionIdentifier_getPropertyKind_4576(node.GetString()).HasGetter() {
					encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
				}
			}
		case int32(29) /* Token.OPTCHAIN_GETPROP */, int32(26) /* Token.GETPROP */ :
			if this_.outer_.priv_PureFunctionIdentifier_getPropertyKind_4576(node.GetString()).HasGetterOrSetter() {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
			}
		case int32(124) /* Token.DYNAMIC_IMPORT */ :
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
		default:
			if NodeUtil_IsCompoundAssignmentOp(node) {
				this_.priv_FunctionBodyAnalyzer_visitLhsNode_4585(encloserSummary, traversal.GetScope(), node.GetFirstChild(), StaticGet_PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL())
				break
			}
			if this_.outer_.compiler.GetAstAnalyzer().NodeTypeMayHaveSideEffects(node) {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unhandled side effect node type "+j2gort.StrVal(node))))
			}
		}
	}
}

// Java method: FunctionBodyAnalyzer.checkIteratesImpureIterable(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary encloserSummary)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) priv_FunctionBodyAnalyzer_checkIteratesImpureIterable_4582(node pkg_rhino.Node, encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	_ = node
	_ = encloserSummary
	if !(NodeUtil_IteratesImpureIterable(node)) {
		return
	}
	encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
}

// Java method: FunctionBodyAnalyzer.deprecatedSetSideEffectsForControlLoss(com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary encloserSummary)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) priv_FunctionBodyAnalyzer_deprecatedSetSideEffectsForControlLoss_4583(encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	_ = encloserSummary
	this_.priv_FunctionBodyAnalyzer_recordThrowsBasedOnContext_4584(encloserSummary)
}

// Java method: FunctionBodyAnalyzer.recordThrowsBasedOnContext(com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary encloserSummary)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) priv_FunctionBodyAnalyzer_recordThrowsBasedOnContext_4584(encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	_ = encloserSummary
	if this_.functionScopeStack.GetLast().catchDepth == int32(0) {
		encloserSummary.SetThrows(this_.functionScopeStack.GetLast().root)
	}
}

// Java method: FunctionBodyAnalyzer.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) EnterScope(t *NodeTraversal) {
	_ = t
	if !(t.GetScope().IsFunctionScope()) {
		return
	}
	var function pkg_rhino.Node = t.GetScopeRoot()
	_ = function
	if !(function.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(function))
	}
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(function)
	_ = jsdoc
	var isArtificiallyPure bool = ((jsdoc != nil) && jsdoc.IsNoSideEffects())
	_ = isArtificiallyPure
	this_.functionScopeStack.AddLast(NewPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(function, isArtificiallyPure))
	if !(this_.outer_.summariesForAllNamesOfFunctionByNode.ContainsKey(function)) {
		var summary *PureFunctionIdentifier_AmbiguatedFunctionSummary = PureFunctionIdentifier_AmbiguatedFunctionSummary_CreateInGraph(j2gort.AsIface[pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]](this_.outer_.reverseCallGraph), "<anonymous>")
		_ = summary
		this_.outer_.priv_PureFunctionIdentifier_putFunctionSummary_4566(function, summary)
	}
}

// Java method: FunctionBodyAnalyzer.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) ExitScope(t *NodeTraversal) {
	_ = t
	var functionScope *Scope = j2gort.CastPtr[Scope](t.GetScope().GetClosestContainerScope())
	_ = functionScope
	if !(functionScope.IsFunctionScope()) {
		return
	}
	var functionEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry = this_.functionScopeStack.GetLast()
	_ = functionEntry
	{
		var iifeCond_1 bool
		{

			rn_29 := functionScope.GetRootNode_0()
			an_30 := functionEntry.root
			iifeCond_1 = j2gort.CheckNotNullNode(rn_29) == an_30
		}
		if !(iifeCond_1) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(functionScope.GetRootNode_0()))
		}
	}
	{
		var iifeCond_1 bool
		{

			rn_31 := t.GetScopeRoot()
			an_32 := functionEntry.root
			iifeCond_1 = j2gort.CheckNotNullNode(rn_31) == an_32
		}
		if iifeCond_1 {
			this_.functionScopeStack.RemoveLast()
		}
	}
	var summaries *PureFunctionIdentifier_SummaryList = this_.outer_.summariesForAllNamesOfFunctionByNode.Get(functionEntry.root)
	_ = summaries
	if summaries != nil {
		{
			var i int32 = int32(0)
			_ = i
			var count int32 = summaries.Size()
			_ = count
			for ; i < count; i = i + 1 {
				var sideEffectInfo *PureFunctionIdentifier_AmbiguatedFunctionSummary = summaries.Get(i)
				_ = sideEffectInfo
				j2gort.CheckNotNullPtr(sideEffectInfo)
				if sideEffectInfo.MutatesGlobalState() {
					continue
				}
				{
					cur_33 := (t.GetScope()).iterVars_()
					for cur_33.HasNext() {

						v := cur_33.Next()

						_ = v
						var isFromDestructuring bool = NodeUtil_IsLhsByDestructuring(v.GetNameNode())
						_ = isFromDestructuring
						if ((v.IsParam() && !(isFromDestructuring)) && !(functionEntry.skiplistedVars.Contains(v))) && functionEntry.taintedVars.Contains(v) {
							sideEffectInfo.SetMutatesArgumentAt(PureFunctionIdentifier_GetParamIndex(v), v.GetNode())
							continue
						}
						if !(functionEntry.IsLocalVar(v, t.GetScope())) {
							if functionEntry.taintedVars.Contains(v) {
								sideEffectInfo.SetMutatesGlobalStateAndAllOtherFlags(v.GetNode())
								break
							}
						}
					}
				}
			}
		}
	}
}

// Java static method: FunctionBodyAnalyzer.isVarDeclaredInSameContainerScope(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Var v, com.google.javascript.jscomp.Scope scope)
func PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v *Var, scope *Scope) bool {
	_ = v
	_ = scope
	return ((v != nil) && j2gort.CastPtr[Scope](v.GetScope()).HasSameContainerScope(ToIface_AbstractScope(scope)))
}

// Java method: FunctionBodyAnalyzer.visitLhsNode(com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary encloserSummary, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node lhs, java.util.function.Predicate<com.google.javascript.rhino.Node> hasLocalRhs)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) priv_FunctionBodyAnalyzer_visitLhsNode_4585(encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary, scope *Scope, lhs pkg_rhino.Node, hasLocalRhs j2gort.Ext_java_util_function_Predicate) {
	_ = encloserSummary
	_ = scope
	_ = lhs
	_ = hasLocalRhs
	if NodeUtil_IsNormalOrOptChainGet(lhs) {
		if lhs.GetFirstChild().IsThis() {
			if this_.functionScopeStack.GetLast().root.IsArrowFunction() {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
			} else {
				encloserSummary.SetMutatesThis(lhs)
			}
		} else {
			var objectNode pkg_rhino.Node = lhs.GetFirstChild()
			_ = objectNode
			if objectNode.IsName() {
				var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(objectNode.GetString()))
				_ = var_
				if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
					this_.functionScopeStack.GetLast().taintedVars.Add(var_)
				} else {
					encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
				}
			} else {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
			}
		}
	} else {
		if !(lhs.IsName()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(lhs))
		}
		var var__2 *Var = j2gort.CastPtr[Var](scope.GetVar(lhs.GetString()))
		_ = var__2
		if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var__2, scope) {
			if !(j2gort.PredTestNode(hasLocalRhs, uint32(lhs))) {
				this_.functionScopeStack.GetLast().skiplistedVars.Add(var__2)
			}
		} else {
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
		}
	}
}

// Java method: FunctionBodyAnalyzer.visitCall(com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary callerInfo, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node invocation)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) priv_FunctionBodyAnalyzer_visitCall_4586(callerInfo *PureFunctionIdentifier_AmbiguatedFunctionSummary, scope *Scope, invocation pkg_rhino.Node) {
	_ = callerInfo
	_ = scope
	_ = invocation
	if invocation.IsCall() && !(this_.outer_.astAnalyzer.FunctionCallHasSideEffects(invocation)) {
		return
	}
	if invocation.IsNew() && !(this_.outer_.astAnalyzer.ConstructorCallHasSideEffects(invocation)) {
		return
	}
	var summaryCount int32 = this_.outer_.priv_PureFunctionIdentifier_collectSummariesForCallee_4570(invocation)
	_ = summaryCount
	if summaryCount == int32(0) {
		callerInfo.SetMutatesGlobalStateAndAllOtherFlags(pkg_rhino.Node(0))
		return
	}
	var calleeSummaries []*PureFunctionIdentifier_AmbiguatedFunctionSummary = this_.outer_.scratchSummaries
	_ = calleeSummaries
	var argsMatchParams bool = this_.outer_.lastInvocationArgsMatchParams
	_ = argsMatchParams
	var enclosingFunction *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry = this_.functionScopeStack.GetLast()
	_ = enclosingFunction
	var propagatesThrows bool = (enclosingFunction.catchDepth == int32(0))
	_ = propagatesThrows
	{
		var s int32 = int32(0)
		_ = s
		for ; s < summaryCount; s = s + 1 {
			var calleeInfo *PureFunctionIdentifier_AmbiguatedFunctionSummary = calleeSummaries[s]
			_ = calleeInfo
			var edge *PureFunctionIdentifier_SideEffectPropagation = allocPooledVia_PureFunctionIdentifier_SideEffectPropagation_ForInvocation(&this_.outer_.slab_SideEffectPropagation, invocation, scope, enclosingFunction, argsMatchParams, propagatesThrows)
			_ = edge
			this_.outer_.reverseCallGraph.Connect_3_DiGraph_DiGraphNode_any_any__any_DiGraph_DiGraphNode_any_any_(calleeInfo.graphNode, edge, callerInfo.graphNode)
		}
	}
}

// Java method: FunctionBodyAnalyzer.argumentsMatchParameters(com.google.javascript.rhino.Node invocation)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) priv_FunctionBodyAnalyzer_argumentsMatchParameters_4587(invocation pkg_rhino.Node) bool {
	_ = invocation
	return ((!(invocation.IsTaggedTemplateLit()) && !(PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation))) && (this_.outer_.compiler.GetCodingConvention().DescribeCachingCall(invocation) == nil))
}

// Java method: FunctionBodyAnalyzer.addToCatchDepthIfTryBlock(com.google.javascript.rhino.Node n, int delta)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) priv_FunctionBodyAnalyzer_addToCatchDepthIfTryBlock_4588(n pkg_rhino.Node, delta int32) {
	_ = n
	_ = delta
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if (!(n.IsBlock()) || !(parent.IsTry())) || !(n.IsFirstChildOf(parent)) {
		return
	}
	var jsCatch pkg_rhino.Node = n.GetNext().GetFirstChild()
	_ = jsCatch
	if jsCatch == 0 {
		return
	}
	r_35 := this_.functionScopeStack.GetLast()
	r_35.catchDepth = (r_35.catchDepth + delta)
}

// Java class: com.google.javascript.jscomp.PureFunctionIdentifier.FunctionBodyAnalyzer.FunctionStackEntry
type PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry struct {
	skiplistedVars     j2gort.CustomOrderedSet[*Var] `gson:"skiplistedVars"`
	taintedVars        j2gort.CustomOrderedSet[*Var] `gson:"taintedVars"`
	root               pkg_rhino.Node                `gson:"root"`
	catchDepth         int32                         `gson:"catchDepth"`
	idHash_            atomic.Int32
	isArtificiallyPure bool `gson:"isArtificiallyPure"`
}

func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) BindSelf_(s any) {
	_ = s
}
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) IsJavaNil_() bool {
	return this_ == nil
}
func (*PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PureFunctionIdentifier$FunctionBodyAnalyzer$FunctionStackEntry", "FunctionStackEntry"
}
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PureFunctionIdentifier.FunctionBodyAnalyzer.FunctionStackEntry", "com.google.javascript.jscomp.PureFunctionIdentifier$FunctionBodyAnalyzer$FunctionStackEntry", func(v_ any) bool {
		return j2gort.InstanceOf[*PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry](v_)
	})
}

func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) initFields_PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry() {
	_ = this_
	this_.skiplistedVars = j2gort.NewCustomOrderedSet[*Var]()
	this_.taintedVars = j2gort.NewCustomOrderedSet[*Var]()

}

// Java constructor: FunctionStackEntry(com.google.javascript.rhino.Node root, boolean isArtificiallyPure)
func NewPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(root pkg_rhino.Node, isArtificiallyPure bool) *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry {
	this_ := &PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry{}
	InitPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(this_, root, isArtificiallyPure)
	return this_
}

func InitPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(this_ *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry, root pkg_rhino.Node, isArtificiallyPure bool) {
	_ = root
	_ = isArtificiallyPure
	this_.initFields_PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry()
	this_.root = root
	this_.isArtificiallyPure = isArtificiallyPure
}

// Java method: FunctionStackEntry.isLocalVar(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Var v, com.google.javascript.jscomp.Scope scope)
func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) IsLocalVar(v *Var, scope *Scope) bool {
	_ = v
	_ = scope
	return (((PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v, scope) && !(this_.skiplistedVars.Contains(v))) && !(v.IsParam())) && !(v.IsCatch()))
}

// Java class: com.google.javascript.jscomp.PureFunctionIdentifier.SideEffectPropagation
type PureFunctionIdentifier_SideEffectPropagation struct {
	callScope                  *Scope                                                          `gson:"callScope"`
	callerEntry                *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry `gson:"callerEntry"`
	invocation                 pkg_rhino.Node                                                  `gson:"invocation"`
	idHash_                    atomic.Int32
	callerIsAlias              bool `gson:"callerIsAlias"`
	calleeThisEqualsCallerThis bool `gson:"calleeThisEqualsCallerThis"`
	propagateThrows            bool `gson:"propagateThrows"`
	argumentsMatchParameters   bool `gson:"argumentsMatchParameters"`
}

func (this_ *PureFunctionIdentifier_SideEffectPropagation) BindSelf_(s any) {
	_ = s
}
func (this_ *PureFunctionIdentifier_SideEffectPropagation) IsJavaNil_() bool { return this_ == nil }
func (*PureFunctionIdentifier_SideEffectPropagation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PureFunctionIdentifier$SideEffectPropagation", "SideEffectPropagation"
}
func (this_ *PureFunctionIdentifier_SideEffectPropagation) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PureFunctionIdentifier_SideEffectPropagation) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PureFunctionIdentifier.SideEffectPropagation", "com.google.javascript.jscomp.PureFunctionIdentifier$SideEffectPropagation", func(v_ any) bool {
		return j2gort.InstanceOf[*PureFunctionIdentifier_SideEffectPropagation](v_)
	})
}

func (this_ *PureFunctionIdentifier_SideEffectPropagation) initFields_PureFunctionIdentifier_SideEffectPropagation() {
	_ = this_
}

// Java constructor: SideEffectPropagation(boolean callerIsAlias, boolean calleeThisEqualsCallerThis, boolean propagateThrows, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node invocation, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope callScope, com.google.javascript.jscomp.PureFunctionIdentifier.FunctionBodyAnalyzer.@org.jspecify.annotations.Nullable FunctionStackEntry callerEntry, boolean argumentsMatchParameters)
func NewPureFunctionIdentifier_SideEffectPropagation(callerIsAlias bool, calleeThisEqualsCallerThis bool, propagateThrows bool, invocation pkg_rhino.Node, callScope *Scope, callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry, argumentsMatchParameters bool) *PureFunctionIdentifier_SideEffectPropagation {
	this_ := &PureFunctionIdentifier_SideEffectPropagation{}
	InitPureFunctionIdentifier_SideEffectPropagation(this_, callerIsAlias, calleeThisEqualsCallerThis, propagateThrows, invocation, callScope, callerEntry, argumentsMatchParameters)
	return this_
}
func allocPooled_PureFunctionIdentifier_SideEffectPropagation(slab_ *j2gort.PooledSlab[PureFunctionIdentifier_SideEffectPropagation], callerIsAlias bool, calleeThisEqualsCallerThis bool, propagateThrows bool, invocation pkg_rhino.Node, callScope *Scope, callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry, argumentsMatchParameters bool) *PureFunctionIdentifier_SideEffectPropagation {
	this_ := slab_.Alloc()
	InitPureFunctionIdentifier_SideEffectPropagation(this_, callerIsAlias, calleeThisEqualsCallerThis, propagateThrows, invocation, callScope, callerEntry, argumentsMatchParameters)
	return this_
}

func InitPureFunctionIdentifier_SideEffectPropagation(this_ *PureFunctionIdentifier_SideEffectPropagation, callerIsAlias bool, calleeThisEqualsCallerThis bool, propagateThrows bool, invocation pkg_rhino.Node, callScope *Scope, callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry, argumentsMatchParameters bool) {
	_ = callerIsAlias
	_ = calleeThisEqualsCallerThis
	_ = propagateThrows
	_ = invocation
	_ = callScope
	_ = callerEntry
	_ = argumentsMatchParameters
	this_.initFields_PureFunctionIdentifier_SideEffectPropagation()
	if !((invocation == 0) || NodeUtil_IsInvocation(invocation)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(invocation))
	}
	this_.callerIsAlias = callerIsAlias
	this_.calleeThisEqualsCallerThis = calleeThisEqualsCallerThis
	this_.propagateThrows = propagateThrows
	this_.invocation = invocation
	this_.callScope = callScope
	this_.callerEntry = callerEntry
	this_.argumentsMatchParameters = argumentsMatchParameters
}

// Java static method: SideEffectPropagation.forAlias()
func PureFunctionIdentifier_SideEffectPropagation_ForAlias() *PureFunctionIdentifier_SideEffectPropagation {
	return NewPureFunctionIdentifier_SideEffectPropagation(true, false, true, pkg_rhino.Node(0), nil, nil, false)
}
func allocPooledVia_PureFunctionIdentifier_SideEffectPropagation_ForAlias(slab_ *j2gort.PooledSlab[PureFunctionIdentifier_SideEffectPropagation]) *PureFunctionIdentifier_SideEffectPropagation {
	return allocPooled_PureFunctionIdentifier_SideEffectPropagation(slab_, true, false, true, pkg_rhino.Node(0), nil, nil, false)
}

// Java static method: SideEffectPropagation.forInvocation(com.google.javascript.rhino.Node invocation, com.google.javascript.jscomp.Scope callScope, com.google.javascript.jscomp.PureFunctionIdentifier.FunctionBodyAnalyzer.FunctionStackEntry callerEntry, boolean argumentsMatchParameters, boolean propagateThrows)
func PureFunctionIdentifier_SideEffectPropagation_ForInvocation(invocation pkg_rhino.Node, callScope *Scope, callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry, argumentsMatchParameters bool, propagateThrows bool) *PureFunctionIdentifier_SideEffectPropagation {
	_ = invocation
	_ = callScope
	_ = callerEntry
	_ = argumentsMatchParameters
	_ = propagateThrows
	if !(NodeUtil_IsInvocation(invocation)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(invocation))
	}
	return NewPureFunctionIdentifier_SideEffectPropagation(false, (PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation) && !(callerEntry.root.IsArrowFunction())), propagateThrows, invocation, callScope, callerEntry, argumentsMatchParameters)
}
func allocPooledVia_PureFunctionIdentifier_SideEffectPropagation_ForInvocation(slab_ *j2gort.PooledSlab[PureFunctionIdentifier_SideEffectPropagation], invocation pkg_rhino.Node, callScope *Scope, callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry, argumentsMatchParameters bool, propagateThrows bool) *PureFunctionIdentifier_SideEffectPropagation {
	_ = invocation
	_ = callScope
	_ = callerEntry
	_ = argumentsMatchParameters
	_ = propagateThrows
	if !(NodeUtil_IsInvocation(invocation)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(invocation))
	}
	return allocPooled_PureFunctionIdentifier_SideEffectPropagation(slab_, false, (PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation) && !(callerEntry.root.IsArrowFunction())), propagateThrows, invocation, callScope, callerEntry, argumentsMatchParameters)
}

// Java method: SideEffectPropagation.holdsLocalValue(com.google.javascript.rhino.Node expr)
func (this_ *PureFunctionIdentifier_SideEffectPropagation) priv_SideEffectPropagation_holdsLocalValue_4590(expr pkg_rhino.Node) bool {
	_ = expr
	return (NodeUtil_EvaluatesToLocalValue(expr) || (expr.IsName() && this_.callerEntry.IsLocalVar(j2gort.CastPtr[Var](this_.callScope.GetVar(expr.GetString())), this_.callScope)))
}

// Java method: SideEffectPropagation.mutatedArgumentsAreLocal(long mutatedParamMask)
func (this_ *PureFunctionIdentifier_SideEffectPropagation) priv_SideEffectPropagation_mutatedArgumentsAreLocal_4591(mutatedParamMask int64) bool {
	_ = mutatedParamMask
	if !(this_.argumentsMatchParameters) {
		return NodeUtil_AllArgsUnescapedLocal(this_.invocation)
	}
	var index int32 = int32(0)
	_ = index
	{
		var arg pkg_rhino.Node = this_.invocation.GetSecondChild()
		_ = arg
		for ; arg != 0; func() { arg = arg.GetNext(); index = index + 1 }() {
			if arg.IsSpread() {
				return false
			}
			if ((index >= int32(64)) || ((mutatedParamMask & (javalang.Shl64(int64(1), int(index)))) != int64(int32(0)))) && !(this_.priv_SideEffectPropagation_holdsLocalValue_4590(arg)) {
				return false
			}
		}
	}
	return true
}

// Java method: SideEffectPropagation.receiverIsLocal()
func (this_ *PureFunctionIdentifier_SideEffectPropagation) priv_SideEffectPropagation_receiverIsLocal_4592() bool {
	var callee pkg_rhino.Node = this_.invocation.GetFirstChild()
	_ = callee
	return ((!(PureFunctionIdentifier_IsInvocationViaCallOrApply(this_.invocation)) && NodeUtil_IsNormalOrOptChainGetProp(callee)) && this_.priv_SideEffectPropagation_holdsLocalValue_4590(callee.GetFirstChild()))
}

// Java static method: SideEffectPropagation.calleeAndCallerShareThis(com.google.javascript.rhino.Node invocation)
func PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation pkg_rhino.Node) bool {
	_ = invocation
	if !(PureFunctionIdentifier_IsCallOrTaggedTemplateLit(invocation)) {
		return false
	}
	var callee pkg_rhino.Node = invocation.GetFirstChild()
	_ = callee
	if callee.IsSuper() {
		return true
	}
	var thisArg pkg_rhino.Node = pkg_rhino.Node(0)
	_ = thisArg
	if PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation) {
		thisArg = invocation.GetSecondChild()
	} else if callee.IsGetProp() || callee.IsOptChainGetProp() {
		thisArg = callee.GetFirstChild()
	} else {
		thisArg = pkg_rhino.Node(0)
	}
	if thisArg == 0 {
		return false
	} else if thisArg.IsThis() || thisArg.IsSuper() {
		return true
	}
	return false
}

// Java method: SideEffectPropagation.propagate(com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary callee, com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary caller)
func (this_ *PureFunctionIdentifier_SideEffectPropagation) Propagate(callee *PureFunctionIdentifier_AmbiguatedFunctionSummary, caller *PureFunctionIdentifier_AmbiguatedFunctionSummary) bool {
	_ = callee
	_ = caller
	var initialCallerFlags int32 = caller.bitmask
	_ = initialCallerFlags
	var initialCallerParamMask int64 = caller.mutatedParamMask
	_ = initialCallerParamMask
	if this_.callerIsAlias {
		caller.priv_AmbiguatedFunctionSummary_setMask_4594(callee.bitmask)
		r_36 := caller
		r_36.mutatedParamMask = (r_36.mutatedParamMask | callee.mutatedParamMask)
	} else if callee.MutatesGlobalState() {
		caller.SetMutatesGlobalStateAndAllOtherFlags(this_.invocation)
	} else {
		if this_.propagateThrows && callee.FunctionThrows() {
			caller.SetThrows(this_.invocation)
		}
		if callee.MutatesArguments() && !(this_.priv_SideEffectPropagation_mutatedArgumentsAreLocal_4591(callee.mutatedParamMask)) {
			caller.SetMutatesGlobalStateAndAllOtherFlags(this_.invocation)
		}
		if callee.MutatesThis() {
			if this_.invocation.IsNew() {
			} else if this_.calleeThisEqualsCallerThis {
				caller.SetMutatesThis(this_.invocation)
			} else if this_.priv_SideEffectPropagation_receiverIsLocal_4592() {
			} else {
				caller.SetMutatesGlobalStateAndAllOtherFlags(this_.invocation)
			}
		}
	}
	return ((caller.bitmask != initialCallerFlags) || (caller.mutatedParamMask != initialCallerParamMask))
}

// Java class: com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary
type PureFunctionIdentifier_AmbiguatedFunctionSummary struct {
	name                              string                                                                                                                          `gson:"name"`
	graphNode                         pkg_graph.DiGraph_DiGraphNode[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation] `gson:"graphNode"`
	impureFunctionReasonsForDebugging *j2gort.Ext_java_util_LinkedHashSet                                                                                             `gson:"impureFunctionReasonsForDebugging"`
	mutatedParamMask                  int64                                                                                                                           `gson:"mutatedParamMask"`
	bitmask                           int32                                                                                                                           `gson:"bitmask"`
	idHash_                           atomic.Int32
	collectImpureDebuggingReason      bool `gson:"collectImpureDebuggingReason"`
}

func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) BindSelf_(s any) {
	_ = s
}
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) IsJavaNil_() bool { return this_ == nil }
func (*PureFunctionIdentifier_AmbiguatedFunctionSummary) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PureFunctionIdentifier$AmbiguatedFunctionSummary", "AmbiguatedFunctionSummary"
}
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PureFunctionIdentifier_AmbiguatedFunctionSummary) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary", "com.google.javascript.jscomp.PureFunctionIdentifier$AmbiguatedFunctionSummary", func(v_ any) bool {
		return j2gort.InstanceOf[*PureFunctionIdentifier_AmbiguatedFunctionSummary](v_)
	})
}

func StaticGet_PureFunctionIdentifier_AmbiguatedFunctionSummary_THROWS() int32 { return int32(1) }
func StaticGet_PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_GLOBAL_STATE() int32 {
	return int32(2)
}
func StaticGet_PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_THIS() int32 { return int32(4) }
func StaticGet_PureFunctionIdentifier_AmbiguatedFunctionSummary_MUTATES_ARGUMENTS() int32 {
	return int32(8)
}
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) initFields_PureFunctionIdentifier_AmbiguatedFunctionSummary() {
	_ = this_
	this_.name = j2gort.NullString

}

// Java constructor: AmbiguatedFunctionSummary(com.google.javascript.jscomp.graph.DiGraph<com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary,com.google.javascript.jscomp.PureFunctionIdentifier.SideEffectPropagation> graph, java.lang.String name)
func NewPureFunctionIdentifier_AmbiguatedFunctionSummary(graph pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation], name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	this_ := &PureFunctionIdentifier_AmbiguatedFunctionSummary{}
	InitPureFunctionIdentifier_AmbiguatedFunctionSummary(this_, graph, name)
	return this_
}

func InitPureFunctionIdentifier_AmbiguatedFunctionSummary(this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary, graph pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation], name string) {
	_ = graph
	_ = name
	this_.initFields_PureFunctionIdentifier_AmbiguatedFunctionSummary()
	this_.name = j2gort.CheckNotNullString(name)
	this_.graphNode = j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation]](graph.CreateNode(this_))
}

// Java static method: AmbiguatedFunctionSummary.createInGraph(com.google.javascript.jscomp.graph.DiGraph<com.google.javascript.jscomp.PureFunctionIdentifier.AmbiguatedFunctionSummary,com.google.javascript.jscomp.PureFunctionIdentifier.SideEffectPropagation> graph, java.lang.String name)
func PureFunctionIdentifier_AmbiguatedFunctionSummary_CreateInGraph(graph pkg_graph.DiGraph[*PureFunctionIdentifier_AmbiguatedFunctionSummary, *PureFunctionIdentifier_SideEffectPropagation], name string) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = graph
	_ = name
	return NewPureFunctionIdentifier_AmbiguatedFunctionSummary(graph, name)
}

// Java method: AmbiguatedFunctionSummary.setMask(int mask)
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) priv_AmbiguatedFunctionSummary_setMask_4594(mask int32) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = mask
	this_.bitmask = (this_.bitmask | mask)
	return this_
}

// Java method: AmbiguatedFunctionSummary.getMask(int mask)
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) priv_AmbiguatedFunctionSummary_getMask_4595(mask int32) bool {
	_ = mask
	return ((this_.bitmask & mask) != int32(0))
}

// Java method: AmbiguatedFunctionSummary.maybeCollectImpureDebuggingReason(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node debuggingReason)
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) priv_AmbiguatedFunctionSummary_maybeCollectImpureDebuggingReason_4596(debuggingReason pkg_rhino.Node) {
	_ = debuggingReason
	if !(this_.collectImpureDebuggingReason) || (debuggingReason == 0) {
		return
	}
	if this_.impureFunctionReasonsForDebugging == nil {
		this_.impureFunctionReasonsForDebugging = j2gort.NewExt_java_util_LinkedHashSet()
	}
	j2gort.SetAddNode(this_.impureFunctionReasonsForDebugging, NodeUtil_GetEnclosingFunction(debuggingReason))
}

// Java method: AmbiguatedFunctionSummary.mutatesThis()
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) MutatesThis() bool {
	return this_.priv_AmbiguatedFunctionSummary_getMask_4595(int32(4))
}

// Java method: AmbiguatedFunctionSummary.setMutatesThis(com.google.javascript.rhino.Node debuggingReason)
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesThis(debuggingReason pkg_rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = debuggingReason
	this_.priv_AmbiguatedFunctionSummary_maybeCollectImpureDebuggingReason_4596(debuggingReason)
	return this_.priv_AmbiguatedFunctionSummary_setMask_4594(int32(4))
}

// Java method: AmbiguatedFunctionSummary.functionThrows()
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) FunctionThrows() bool {
	return this_.priv_AmbiguatedFunctionSummary_getMask_4595(int32(1))
}

// Java method: AmbiguatedFunctionSummary.setThrows(com.google.javascript.rhino.Node debuggingReason)
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetThrows(debuggingReason pkg_rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = debuggingReason
	this_.priv_AmbiguatedFunctionSummary_maybeCollectImpureDebuggingReason_4596(debuggingReason)
	return this_.priv_AmbiguatedFunctionSummary_setMask_4594(int32(1))
}

// Java method: AmbiguatedFunctionSummary.mutatesGlobalState()
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) MutatesGlobalState() bool {
	return this_.priv_AmbiguatedFunctionSummary_getMask_4595(int32(2))
}

// Java method: AmbiguatedFunctionSummary.setMutatesGlobalStateAndAllOtherFlags(com.google.javascript.rhino.Node debuggingReason)
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesGlobalStateAndAllOtherFlags(debuggingReason pkg_rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = debuggingReason
	this_.priv_AmbiguatedFunctionSummary_maybeCollectImpureDebuggingReason_4596(debuggingReason)
	this_.mutatedParamMask = ^(int64(0))
	return this_.priv_AmbiguatedFunctionSummary_setMask_4594((((int32(1) | int32(4)) | int32(8)) | int32(2)))
}

// Java method: AmbiguatedFunctionSummary.mutatesArguments()
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) MutatesArguments() bool {
	return this_.priv_AmbiguatedFunctionSummary_getMask_4595(int32(8))
}

// Java method: AmbiguatedFunctionSummary.setMutatesArguments(com.google.javascript.rhino.Node debuggingReason)
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesArguments(debuggingReason pkg_rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = debuggingReason
	this_.priv_AmbiguatedFunctionSummary_maybeCollectImpureDebuggingReason_4596(debuggingReason)
	this_.mutatedParamMask = ^(int64(0))
	return this_.priv_AmbiguatedFunctionSummary_setMask_4594(int32(8))
}

// Java method: AmbiguatedFunctionSummary.setMutatesArgumentAt(int paramIndex, com.google.javascript.rhino.Node debuggingReason)
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetMutatesArgumentAt(paramIndex int32, debuggingReason pkg_rhino.Node) *PureFunctionIdentifier_AmbiguatedFunctionSummary {
	_ = paramIndex
	_ = debuggingReason
	if !(paramIndex >= int32(0)) {
		j2gort.CheckArgumentMsg(false, j2gort.ValueOfInt(paramIndex))
	}
	this_.priv_AmbiguatedFunctionSummary_maybeCollectImpureDebuggingReason_4596(debuggingReason)
	var paramBit int64
	_ = paramBit
	if paramIndex < int32(64) {
		paramBit = javalang.Shl64(int64(1), int(paramIndex))
	} else {
		paramBit = ^(int64(0))
	}
	this_.mutatedParamMask = (this_.mutatedParamMask | paramBit)
	return this_.priv_AmbiguatedFunctionSummary_setMask_4594(int32(8))
}

// Java method: AmbiguatedFunctionSummary.hasNoFlagsSet()
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) HasNoFlagsSet() bool {
	return (this_.bitmask == int32(0))
}

// Java method: AmbiguatedFunctionSummary.setCollectImpureDebuggingReason()
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) SetCollectImpureDebuggingReason() {
	this_.collectImpureDebuggingReason = true
}

// Java method: AmbiguatedFunctionSummary.toString()
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_Class(j2gort.GetObjectClass(this_)).Add_2_string_any("name", j2gort.BoxString(this_.name)).Add_2_string_int32("graphNode", javalang.ObjectHashCode(j2gort.CheckNotNull(this_.graphNode))).Add_2_string_any("sideEffects", j2gort.BoxString(this_.priv_AmbiguatedFunctionSummary_sideEffectsToString_4609())).String()
}

// Java method: AmbiguatedFunctionSummary.sideEffectsToString()
func (this_ *PureFunctionIdentifier_AmbiguatedFunctionSummary) priv_AmbiguatedFunctionSummary_sideEffectsToString_4609() string {
	var status []string = []string{}
	_ = status
	if this_.MutatesThis() {
		status = append(status, "this")
	}
	if this_.MutatesGlobalState() {
		status = append(status, "global")
	}
	if this_.MutatesArguments() {
		status = append(status, "args")
	}
	if this_.FunctionThrows() {
		status = append(status, "throw")
	}
	return j2gort.SliceFormatJavaString(status)
}

// Java class: com.google.javascript.jscomp.PureFunctionIdentifier.Driver
type PureFunctionIdentifier_Driver struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *PureFunctionIdentifier_Driver) BindSelf_(s any) {
	_ = s
}
func (this_ *PureFunctionIdentifier_Driver) IsJavaNil_() bool { return this_ == nil }
func (*PureFunctionIdentifier_Driver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PureFunctionIdentifier$Driver", "Driver"
}
func (this_ *PureFunctionIdentifier_Driver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PureFunctionIdentifier_Driver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PureFunctionIdentifier.Driver", "com.google.javascript.jscomp.PureFunctionIdentifier$Driver", func(v_ any) bool {
		return j2gort.InstanceOf[*PureFunctionIdentifier_Driver](v_)
	})
}
func (*PureFunctionIdentifier_Driver) Implements_jscomp_CompilerPass() {}

func (this_ *PureFunctionIdentifier_Driver) initFields_PureFunctionIdentifier_Driver() {
	_ = this_
}

// Java constructor: Driver(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewPureFunctionIdentifier_Driver(compiler AbstractCompiler) *PureFunctionIdentifier_Driver {
	this_ := &PureFunctionIdentifier_Driver{}
	InitPureFunctionIdentifier_Driver(this_, compiler)
	return this_
}

func InitPureFunctionIdentifier_Driver(this_ *PureFunctionIdentifier_Driver, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_PureFunctionIdentifier_Driver()
	this_.compiler = compiler
}

// Java method: Driver.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *PureFunctionIdentifier_Driver) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var shouldValidate *CheckLevel = this_.compiler.GetOptions().GetWarningsGuard().Level(JSError_Make_2(StaticGet_PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION(), []string{"", ""}))
	_ = shouldValidate
	OptimizeCalls_S_Builder().SetCompiler(this_.compiler).SetConsiderExterns(true).AddPass(NewPureFunctionIdentifier(this_.compiler, this_.compiler.GetOptions().GetAssumeGettersArePure(), (shouldValidate != StaticGet_CheckLevel_OFF()))).Build().Process_2(externs, root)
}
func (*PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) AlwaysShouldTraverse_()      {}
func (*PureFunctionIdentifier_ExternFunctionAnnotationAnalyzer) CanReleaseSyntacticScopes_() {}
