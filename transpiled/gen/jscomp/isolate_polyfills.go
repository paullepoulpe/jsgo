/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/IsolatePolyfills.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_resources "github.com/paullepoulpe/jsgo/transpiled/gen/resources"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.IsolatePolyfills
type IsolatePolyfills struct {
	compiler                 AbstractCompiler               `gson:"compiler"`
	polyfills                *PolyfillUsageFinder_Polyfills `gson:"polyfills"`
	jscompPolyfillsObject    pkg_rhino.Node                 `gson:"jscompPolyfillsObject"`
	jscompLookupMethod       pkg_rhino.Node                 `gson:"jscompLookupMethod"`
	idHash_                  atomic.Int32
	usedPolyfillMethodLookup bool `gson:"usedPolyfillMethodLookup"`
	isTempVarInitialized     bool `gson:"isTempVarInitialized"`
}

func (this_ *IsolatePolyfills) BindSelf_(s any) {
	_ = s
}
func (this_ *IsolatePolyfills) IsJavaNil_() bool { return this_ == nil }
func (*IsolatePolyfills) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.IsolatePolyfills", "IsolatePolyfills"
}
func (this_ *IsolatePolyfills) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IsolatePolyfills) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.IsolatePolyfills", "com.google.javascript.jscomp.IsolatePolyfills", func(v_ any) bool {
		return j2gort.InstanceOf[*IsolatePolyfills](v_)
	})
}
func (*IsolatePolyfills) Implements_jscomp_CompilerPass() {}

func StaticGet_IsolatePolyfills_POLYFILL_TEMP() string { return "$jscomp$polyfillTmp" }

var static_IsolatePolyfills_FILES_ALLOWED_UNQUALIFIED_POLYFILL_ACCESSES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_IsolatePolyfills_FILES_ALLOWED_UNQUALIFIED_POLYFILL_ACCESSES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_IsolatePolyfills()
	return static_IsolatePolyfills_FILES_ALLOWED_UNQUALIFIED_POLYFILL_ACCESSES
}
func StaticSet_IsolatePolyfills_FILES_ALLOWED_UNQUALIFIED_POLYFILL_ACCESSES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_IsolatePolyfills()
	static_IsolatePolyfills_FILES_ALLOWED_UNQUALIFIED_POLYFILL_ACCESSES = v
}

var clinitDone_IsolatePolyfills uint32
var clinitState_IsolatePolyfills j2gort.ClinitState

func clinit_IsolatePolyfills() {
	if atomic.LoadUint32(&clinitDone_IsolatePolyfills) != 0 {
		return
	}
	clinitSlow_IsolatePolyfills()
}

func clinitSlow_IsolatePolyfills() {
	if !j2gort.EnterClinit(&clinitState_IsolatePolyfills, &clinitDone_IsolatePolyfills) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_IsolatePolyfills, &clinitDone_IsolatePolyfills)
	static_IsolatePolyfills_FILES_ALLOWED_UNQUALIFIED_POLYFILL_ACCESSES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3("src/com/google/javascript/jscomp/js/"+"util/global.js", "src/com/google/javascript/jscomp/js/"+"util/shouldpolyfill.js", "src/com/google/javascript/jscomp/js/"+"es6/util/construct.js")
}

func (this_ *IsolatePolyfills) initFields_IsolatePolyfills() {
	_ = this_
	this_.jscompLookupMethod = pkg_rhino.IR_Name("$jscomp$lookupPolyfilledValue")

}

// Java constructor: IsolatePolyfills(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewIsolatePolyfills_1(compiler AbstractCompiler) *IsolatePolyfills {
	clinit_IsolatePolyfills()
	this_ := &IsolatePolyfills{}
	InitIsolatePolyfills_1(this_, compiler)
	return this_
}

func InitIsolatePolyfills_1(this_ *IsolatePolyfills, compiler AbstractCompiler) {
	_ = compiler
	InitIsolatePolyfills_2(this_, compiler, PolyfillUsageFinder_Polyfills_FromTable(pkg_resources.ResourceLoader_LoadTextResource(j2gort.ClassLiteral("com.google.javascript.jscomp.IsolatePolyfills", "IsolatePolyfills"), "js/polyfills.txt")))
}

// Java constructor: IsolatePolyfills(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.PolyfillUsageFinder.Polyfills polyfills)
func NewIsolatePolyfills_2(compiler AbstractCompiler, polyfills *PolyfillUsageFinder_Polyfills) *IsolatePolyfills {
	clinit_IsolatePolyfills()
	this_ := &IsolatePolyfills{}
	InitIsolatePolyfills_2(this_, compiler, polyfills)
	return this_
}

func InitIsolatePolyfills_2(this_ *IsolatePolyfills, compiler AbstractCompiler, polyfills *PolyfillUsageFinder_Polyfills) {
	_ = compiler
	_ = polyfills
	this_.initFields_IsolatePolyfills()
	this_.compiler = compiler
	this_.polyfills = polyfills
	var hasPropertyCollapsingRun bool
	{
		re_1 := compiler.GetOptions().GetPropertyCollapseLevel()
		ae_2 := StaticGet_CompilerOptions_PropertyCollapseLevel_ALL()
		if re_1 == 0 {
			j2gort.ThrowNPE()
		}
		hasPropertyCollapsingRun = re_1 == ae_2
	}

	_ = hasPropertyCollapsingRun
	if hasPropertyCollapsingRun {
		this_.jscompPolyfillsObject = IsolatePolyfills_CreateCollapsedName()
	} else {
		this_.jscompPolyfillsObject = IsolatePolyfills_CreateJSCompPolyfillsAccess()
	}
}

// Java static method: IsolatePolyfills.createCollapsedName()
func IsolatePolyfills_CreateCollapsedName() pkg_rhino.Node {
	clinit_IsolatePolyfills()
	var collapsedName pkg_rhino.Node = pkg_rhino.IR_Name("$jscomp$polyfills")
	_ = collapsedName
	collapsedName.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	return collapsedName
}

// Java static method: IsolatePolyfills.createJSCompPolyfillsAccess()
func IsolatePolyfills_CreateJSCompPolyfillsAccess() pkg_rhino.Node {
	clinit_IsolatePolyfills()
	var jscomp pkg_rhino.Node = pkg_rhino.IR_Name("$jscomp")
	_ = jscomp
	jscomp.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	return pkg_rhino.IR_Getprop_2(jscomp, "polyfills")
}

// Java method: IsolatePolyfills.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *IsolatePolyfills) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var injectedPolyfills j2gort.Ext_com_google_common_collect_ImmutableSet = this_.priv_IsolatePolyfills_findAllInjectedPolyfills_3449()
	_ = injectedPolyfills
	var polyfillUsages j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = polyfillUsages
	var shouldDisambiguateProperties bool = this_.compiler.GetOptions().ShouldDisambiguateProperties()
	_ = shouldDisambiguateProperties
	{
		clinit_PolyfillUsageFinder()
		var val_3 PolyfillUsageFinder
		InitPolyfillUsageFinder(&val_3, this_.compiler, this_.polyfills, shouldDisambiguateProperties)
		(&val_3).TraverseIncludingGuarded(root, func() j2gort.Ext_java_util_function_Consumer {
			bound_4 := polyfillUsages
			return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_4.Add_1(mrp0_) })
		}())
	}
	var visitedNodes j2gort.LocalOrderedSet[pkg_rhino.Node] = j2gort.NewLocalOrderedSet[pkg_rhino.Node]()
	_ = visitedNodes
	{
		cur_5 := j2gort.NewIterCursor(polyfillUsages)
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			usage := j2gort.CastPtr[PolyfillUsageFinder_PolyfillUsage](rawElem_6)
			_ = usage
			if (visitedNodes.Contains(usage.Node()) || (j2gort.CheckNotNullString(usage.Polyfill().library) == "")) || !(j2gort.SetContainsStr(injectedPolyfills, usage.Polyfill().nativeSymbol)) {
				continue
			}
			this_.priv_IsolatePolyfills_rewritePolyfill_3451(usage)
			visitedNodes.Add(usage.Node())
		}
	}
	this_.priv_IsolatePolyfills_cleanUpJscompLookupPolyfilledValue_3455()
}

// Java method: IsolatePolyfills.findAllInjectedPolyfills()
func (this_ *IsolatePolyfills) priv_IsolatePolyfills_findAllInjectedPolyfills_3449() j2gort.Ext_com_google_common_collect_ImmutableSet {
	var actualPolyfills j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = actualPolyfills
	var lastInjectedNode pkg_rhino.Node = this_.compiler.GetNodeForCodeInsertion(nil)
	_ = lastInjectedNode
	NodeTraversal_Traverse(this_.compiler, lastInjectedNode, NewIsolatePolyfills__Anon_53(this_, actualPolyfills))
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](actualPolyfills.Build())
}

// Java method: IsolatePolyfills.isJSCompPolyfillCall(com.google.javascript.rhino.Node call)
func (this_ *IsolatePolyfills) priv_IsolatePolyfills_isJSCompPolyfillCall_3450(call pkg_rhino.Node) bool {
	_ = call
	if !(call.IsCall()) {
		return false
	}
	var jscompPolyfillName string
	_ = jscompPolyfillName
	{
		var iifeCond_1 bool
		{

			re_7 := this_.compiler.GetOptions().GetPropertyCollapseLevel()
			ae_8 := StaticGet_CompilerOptions_PropertyCollapseLevel_ALL()
			if re_7 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_7 == ae_8
		}
		if iifeCond_1 {
			jscompPolyfillName = "$jscomp$polyfill"
		} else {
			jscompPolyfillName = "$jscomp.polyfill"
		}
	}
	return call.GetFirstChild().MatchesQualifiedName_1_string(jscompPolyfillName)
}

// Java method: IsolatePolyfills.rewritePolyfill(com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage polyfillUsage)
func (this_ *IsolatePolyfills) priv_IsolatePolyfills_rewritePolyfill_3451(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	_ = polyfillUsage
	var polyfill *PolyfillUsageFinder_Polyfill = polyfillUsage.Polyfill()
	_ = polyfill
	if this_.compiler.GetOptions().GetOutputFeatureSet().Contains_1_FeatureSet(pkg_parser.FeatureSet_ValueOf(polyfill.nativeVersion)) {
		return
	}
	var polyfillAccess pkg_rhino.Node = polyfillUsage.Node()
	_ = polyfillAccess
	var parent pkg_rhino.Node = polyfillUsage.Node().GetParent()
	_ = parent
	if j2gort.SetContainsStr(StaticGet_IsolatePolyfills_FILES_ALLOWED_UNQUALIFIED_POLYFILL_ACCESSES(), polyfillAccess.GetSourceFileName()) {
		return
	}
	if NodeUtil_IsLValue(polyfillAccess) || (parent.IsAssign() && polyfillAccess.IsFirstChildOf(parent)) {
		return
	}
	var name string = polyfillUsage.Name()
	_ = name
	var isGlobalClass bool = ((j2gort.JStrIndexOfCp(name, int32(uint16(46))) == int32(-1)) && func() bool {
		re_9 := polyfill.kind
		ae_10 := StaticGet_PolyfillUsageFinder_Polyfill_Kind_STATIC()
		if re_9 == 0 {
			j2gort.ThrowNPE()
		}
		return re_9 == ae_10
	}())
	_ = isGlobalClass
	if isGlobalClass {
		polyfillAccess.ReplaceWith(pkg_rhino.IR_Getelem(this_.jscompPolyfillsObject.CloneTree_0(), pkg_rhino.IR_String(name)).SrcrefTree(polyfillAccess))
	} else if (parent.IsCall() || parent.IsOptChainCall()) && polyfillAccess.IsFirstChildOf(parent) {
		this_.priv_IsolatePolyfills_rewritePolyfillInCall_3452(polyfillAccess)
	} else {
		var methodName pkg_rhino.Node = pkg_rhino.IR_String(polyfillAccess.GetString()).Srcref(polyfillAccess)
		_ = methodName
		var receiver pkg_rhino.Node = polyfillAccess.RemoveFirstChild()
		_ = receiver
		polyfillAccess.ReplaceWith(this_.priv_IsolatePolyfills_createPolyfillMethodLookup_3454(receiver, methodName, NodeUtil_IsOptChainNode(polyfillAccess)).SrcrefTree(polyfillAccess))
	}
	this_.compiler.ReportChangeToEnclosingScope(parent)
}

// Java method: IsolatePolyfills.rewritePolyfillInCall(com.google.javascript.rhino.Node callee)
func (this_ *IsolatePolyfills) priv_IsolatePolyfills_rewritePolyfillInCall_3452(callee pkg_rhino.Node) {
	_ = callee
	var methodName pkg_rhino.Node = pkg_rhino.IR_String(callee.GetString()).Srcref(callee)
	_ = methodName
	var receiver pkg_rhino.Node = callee.RemoveFirstChild()
	_ = receiver
	var isCalleeOptChain bool = NodeUtil_IsOptChainNode(callee)
	_ = isCalleeOptChain
	var requiresTemp bool = this_.compiler.GetAstAnalyzer().MayEffectMutableState(receiver)
	_ = requiresTemp
	var polyfilledMethod pkg_rhino.Node = pkg_rhino.Node(0)
	_ = polyfilledMethod
	var thisNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = thisNode
	if requiresTemp {
		thisNode = this_.priv_IsolatePolyfills_createTempName_3453(callee)
		polyfilledMethod = pkg_rhino.IR_Comma(pkg_rhino.IR_Assign(thisNode.CloneTree_0(), receiver), this_.priv_IsolatePolyfills_createPolyfillMethodLookup_3454(thisNode.CloneTree_0(), methodName, isCalleeOptChain))
	} else {
		thisNode = receiver
		polyfilledMethod = this_.priv_IsolatePolyfills_createPolyfillMethodLookup_3454(receiver.CloneTree_0(), methodName, isCalleeOptChain)
	}
	var receiverDotCall pkg_rhino.Node = pkg_rhino.Node(0)
	_ = receiverDotCall
	if NodeUtil_IsOptChainNode(callee) {
		receiverDotCall = pkg_rhino.IR_StartOptChainGetprop(polyfilledMethod, "call").SrcrefTree(callee)
	} else {
		receiverDotCall = pkg_rhino.IR_Getprop_2(polyfilledMethod, "call").SrcrefTree(callee)
	}
	callee.ReplaceWith(receiverDotCall)
	thisNode.InsertAfter(receiverDotCall)
}

// Java method: IsolatePolyfills.createTempName(com.google.javascript.rhino.Node srcref)
func (this_ *IsolatePolyfills) priv_IsolatePolyfills_createTempName_3453(srcref pkg_rhino.Node) pkg_rhino.Node {
	_ = srcref
	if !(this_.isTempVarInitialized) {
		this_.isTempVarInitialized = true
		var decl pkg_rhino.Node = pkg_rhino.IR_Var_1(pkg_rhino.IR_Name("$jscomp$polyfillTmp")).SrcrefTree(srcref)
		_ = decl
		this_.compiler.GetNodeForCodeInsertion(nil).AddChildToFront(decl)
		this_.compiler.ReportChangeToEnclosingScope(decl)
	}
	return pkg_rhino.IR_Name("$jscomp$polyfillTmp").Srcref(srcref)
}

// Java method: IsolatePolyfills.createPolyfillMethodLookup(com.google.javascript.rhino.Node receiver, com.google.javascript.rhino.Node methodName, boolean isOptChainNode)
func (this_ *IsolatePolyfills) priv_IsolatePolyfills_createPolyfillMethodLookup_3454(receiver pkg_rhino.Node, methodName pkg_rhino.Node, isOptChainNode bool) pkg_rhino.Node {
	_ = receiver
	_ = methodName
	_ = isOptChainNode
	this_.usedPolyfillMethodLookup = true
	var call pkg_rhino.Node = pkg_rhino.Node(0)
	_ = call
	if isOptChainNode {
		call = pkg_rhino.IR_Call(this_.jscompLookupMethod.CloneTree_0(), []pkg_rhino.Node{receiver, methodName, pkg_rhino.IR_TrueNode()})
	} else {
		call = pkg_rhino.IR_Call(this_.jscompLookupMethod.CloneTree_0(), []pkg_rhino.Node{receiver, methodName})
	}
	call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	return call
}

// Java method: IsolatePolyfills.cleanUpJscompLookupPolyfilledValue()
func (this_ *IsolatePolyfills) priv_IsolatePolyfills_cleanUpJscompLookupPolyfilledValue_3455() {
	var syntheticExternsRoot pkg_rhino.Node = this_.compiler.GetSynthesizedExternsInput().GetAstRoot(this_.compiler)
	_ = syntheticExternsRoot
	scInline_1 := NewSyntacticScopeCreator_1(this_.compiler)
	defer scInline_1.ReleaseSlabs_()
	var syntheticExternsScope *Scope = scInline_1.CreateScope_2_Node_Scope(syntheticExternsRoot, nil)

	_ = syntheticExternsScope
	var externVar *Var = j2gort.CheckNotNullPtr(j2gort.CastPtr[Var](syntheticExternsScope.GetVar(this_.jscompLookupMethod.GetString())))
	_ = externVar
	NodeUtil_DeleteNode(externVar.GetParentNode(), this_.compiler)
	if this_.usedPolyfillMethodLookup {
		return
	}
	scInline_2 := NewSyntacticScopeCreator_1(this_.compiler)
	defer scInline_2.ReleaseSlabs_()
	var syntheticCodeScope *Scope = scInline_2.CreateScope_2_Node_Scope(this_.compiler.GetNodeForCodeInsertion(nil), nil)

	_ = syntheticCodeScope
	var syntheticVar *Var = j2gort.CastPtr[Var](syntheticCodeScope.GetVar(this_.jscompLookupMethod.GetString()))
	_ = syntheticVar
	if syntheticVar != nil {
		NodeUtil_DeleteNode(syntheticVar.GetParentNode(), this_.compiler)
	}
}

// Java class:
type IsolatePolyfills__Anon_53 struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer_                *IsolatePolyfills
	cap_0_actualPolyfills j2gort.Ext_com_google_common_collect_ImmutableSet_Builder
}

func (this_ *IsolatePolyfills__Anon_53) BindSelf_(s any) {
	_ = s
}
func (this_ *IsolatePolyfills__Anon_53) IsJavaNil_() bool { return this_ == nil }
func (*IsolatePolyfills__Anon_53) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.IsolatePolyfills$1", ""
}
func (this_ *IsolatePolyfills__Anon_53) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowCallbackBase.IdentityHash_()
}
func (*IsolatePolyfills__Anon_53) HasIdentityEquals_()                                      {}
func (this_ *IsolatePolyfills__Anon_53) SetOuter_(o *IsolatePolyfills)                      { this_.outer_ = o }
func (this_ *IsolatePolyfills__Anon_53) GetOuter_() *IsolatePolyfills                       { return this_.outer_ }
func (*IsolatePolyfills__Anon_53) Implements_jscomp_NodeTraversal_AbstractShallowCallback() {}
func (*IsolatePolyfills__Anon_53) Implements_jscomp_NodeTraversal_Callback()                {}

func (this_ *IsolatePolyfills__Anon_53) initFields_IsolatePolyfills__Anon_53() {
	_ = this_
}

func NewIsolatePolyfills__Anon_53(outer_ *IsolatePolyfills, cap_p_actualPolyfills j2gort.Ext_com_google_common_collect_ImmutableSet_Builder) *IsolatePolyfills__Anon_53 {
	this_ := &IsolatePolyfills__Anon_53{}
	this_.outer_ = outer_
	this_.cap_0_actualPolyfills = cap_p_actualPolyfills
	InitNodeTraversal_AbstractShallowCallback(&this_.NodeTraversal_AbstractShallowCallbackBase)
	this_.initFields_IsolatePolyfills__Anon_53()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *IsolatePolyfills__Anon_53) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.outer_.priv_IsolatePolyfills_isJSCompPolyfillCall_3450(n) {
		var polyfilledSymbol string = n.GetSecondChild().GetString()
		_ = polyfilledSymbol
		j2gort.ImmSetBuilderAddStrVoid(this_.cap_0_actualPolyfills, polyfilledSymbol)
	}
}
func (*IsolatePolyfills__Anon_53) CanReleaseSyntacticScopes_() {}
