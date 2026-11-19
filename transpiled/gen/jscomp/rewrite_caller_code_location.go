/*
 * Copyright 2024 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RewriteCallerCodeLocation.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.RewriteCallerCodeLocation
type RewriteCallerCodeLocation struct {
	compiler                    AbstractCompiler                                                                  `gson:"compiler"`
	astFactory                  *AstFactory                                                                       `gson:"astFactory"`
	callerLocationFunctionNames j2gort.OrderedMap[string, *RewriteCallerCodeLocation_FunctionVarAndParamPosition] `gson:"callerLocationFunctionNames"`
	idHash_                     atomic.Int32
}

func (this_ *RewriteCallerCodeLocation) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteCallerCodeLocation) IsJavaNil_() bool { return this_ == nil }
func (*RewriteCallerCodeLocation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteCallerCodeLocation", "RewriteCallerCodeLocation"
}
func (this_ *RewriteCallerCodeLocation) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteCallerCodeLocation) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteCallerCodeLocation", "com.google.javascript.jscomp.RewriteCallerCodeLocation", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteCallerCodeLocation](v_)
	})
}
func (*RewriteCallerCodeLocation) Implements_jscomp_CompilerPass() {}

var static_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR *DiagnosticType

func StaticGet_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR() *DiagnosticType {
	clinit_RewriteCallerCodeLocation()
	return static_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR
}
func StaticSet_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR(v *DiagnosticType) {
	clinit_RewriteCallerCodeLocation()
	static_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR = v
}

var static_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR *DiagnosticType

func StaticGet_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR() *DiagnosticType {
	clinit_RewriteCallerCodeLocation()
	return static_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR
}
func StaticSet_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR(v *DiagnosticType) {
	clinit_RewriteCallerCodeLocation()
	static_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR = v
}

var static_RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR *DiagnosticType

func StaticGet_RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR() *DiagnosticType {
	clinit_RewriteCallerCodeLocation()
	return static_RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR
}
func StaticSet_RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR(v *DiagnosticType) {
	clinit_RewriteCallerCodeLocation()
	static_RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR = v
}

var static_RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR *DiagnosticType

func StaticGet_RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR() *DiagnosticType {
	clinit_RewriteCallerCodeLocation()
	return static_RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR
}
func StaticSet_RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR(v *DiagnosticType) {
	clinit_RewriteCallerCodeLocation()
	static_RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR = v
}

var static_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME pkg_rhino.QualifiedName

func StaticGet_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME() pkg_rhino.QualifiedName {
	clinit_RewriteCallerCodeLocation()
	return static_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME
}
func StaticSet_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME(v pkg_rhino.QualifiedName) {
	clinit_RewriteCallerCodeLocation()
	static_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME = v
}

var clinitDone_RewriteCallerCodeLocation uint32
var clinitState_RewriteCallerCodeLocation j2gort.ClinitState

func clinit_RewriteCallerCodeLocation() {
	if atomic.LoadUint32(&clinitDone_RewriteCallerCodeLocation) != 0 {
		return
	}
	clinitSlow_RewriteCallerCodeLocation()
}

func clinitSlow_RewriteCallerCodeLocation() {
	if !j2gort.EnterClinit(&clinitState_RewriteCallerCodeLocation, &clinitDone_RewriteCallerCodeLocation) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RewriteCallerCodeLocation, &clinitDone_RewriteCallerCodeLocation)
	static_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR = DiagnosticType_Error("JSC_CALLER_LOCATION_POSITION_ERROR", "Please make sure there is only one goog.callerLocation argument in your function's parameter list, and it is the first optional argument in the list")
	static_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR = DiagnosticType_Error("JSC_CALLER_LOCATION_MISUSE_ERROR", "goog.callerLocation should only be used as a default parameter initializer")
	static_RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR = DiagnosticType_Error("JSC_UNDEFINED_CODE_LOCATION_ERROR", "Do not pass in undefined as an argument to goog.CodeLocation parameter")
	static_RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR = DiagnosticType_Error("JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR", "Do not use goog.callerLocation in an anonymous functions. Functions that use goog.callerLocation should be named.")
	static_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME = pkg_rhino.QualifiedName_Of("goog.callerLocation")
}

func (this_ *RewriteCallerCodeLocation) initFields_RewriteCallerCodeLocation() {
	_ = this_
}

// Java constructor: RewriteCallerCodeLocation(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRewriteCallerCodeLocation(compiler AbstractCompiler) *RewriteCallerCodeLocation {
	clinit_RewriteCallerCodeLocation()
	this_ := &RewriteCallerCodeLocation{}
	InitRewriteCallerCodeLocation(this_, compiler)
	return this_
}

func InitRewriteCallerCodeLocation(this_ *RewriteCallerCodeLocation, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RewriteCallerCodeLocation()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.callerLocationFunctionNames = j2gort.NewOrderedMap[string, *RewriteCallerCodeLocation_FunctionVarAndParamPosition](

	// Java method: RewriteCallerCodeLocation.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
	)
}

func (this_ *RewriteCallerCodeLocation) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer func() {
		this_.callerLocationFunctionNames.Clear()
		sc_.ReleaseSlabs_()
	}()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(NewRewriteCallerCodeLocation_FindCallerLocationFunctions(this_))
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(this_.compiler.GetRoot())
	}

	if !(this_.callerLocationFunctionNames.IsEmpty()) {
		{
			clinit_NodeTraversal()
			var ntb_ NodeTraversal_Builder
			InitNodeTraversal_Builder(&ntb_)
			(&ntb_).SetCompiler(this_.compiler)
			(&ntb_).SetCallback_1_NodeTraversal_Callback(NewRewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls(this_))
			(&ntb_).SetScopeCreator(sc_)
			(&ntb_).Traverse(this_.compiler.GetRoot())
		}

	}
}

// Java class: com.google.javascript.jscomp.RewriteCallerCodeLocation.FindCallerLocationFunctions
type RewriteCallerCodeLocation_FindCallerLocationFunctions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RewriteCallerCodeLocation
}

func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) IsJavaNil_() bool {
	return this_ == nil
}
func (*RewriteCallerCodeLocation_FindCallerLocationFunctions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteCallerCodeLocation$FindCallerLocationFunctions", "FindCallerLocationFunctions"
}
func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RewriteCallerCodeLocation_FindCallerLocationFunctions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteCallerCodeLocation.FindCallerLocationFunctions", "com.google.javascript.jscomp.RewriteCallerCodeLocation$FindCallerLocationFunctions", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteCallerCodeLocation_FindCallerLocationFunctions](v_)
	})
}
func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) SetOuter_(o *RewriteCallerCodeLocation) {
	this_.outer_ = o
}
func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) GetOuter_() *RewriteCallerCodeLocation {
	return this_.outer_
}
func (*RewriteCallerCodeLocation_FindCallerLocationFunctions) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*RewriteCallerCodeLocation_FindCallerLocationFunctions) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) initFields_RewriteCallerCodeLocation_FindCallerLocationFunctions() {
	_ = this_
}

// Java constructor: FindCallerLocationFunctions()
func NewRewriteCallerCodeLocation_FindCallerLocationFunctions(outer_ *RewriteCallerCodeLocation) *RewriteCallerCodeLocation_FindCallerLocationFunctions {
	this_ := &RewriteCallerCodeLocation_FindCallerLocationFunctions{}
	this_.outer_ = outer_
	InitRewriteCallerCodeLocation_FindCallerLocationFunctions(this_)
	return this_
}

func InitRewriteCallerCodeLocation_FindCallerLocationFunctions(this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RewriteCallerCodeLocation_FindCallerLocationFunctions()
}

// Java method: FindCallerLocationFunctions.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.priv_FindCallerLocationFunctions_isGoogCallerLocationMisused_4966(n, parent) {
		this_.outer_.compiler.Report_1(JSError_Make_3(parent.GetParent(), StaticGet_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR(), []string{}))
		return
	}
	if n.IsParamList() {
		this_.priv_FindCallerLocationFunctions_visitParamListAndAddCallerLocationFunctionNames_4967(n, t)
	}
}

// Java method: FindCallerLocationFunctions.isGoogCallerLocationMisused(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) priv_FindCallerLocationFunctions_isGoogCallerLocationMisused_4966(n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = n
	_ = parent
	if !(StaticGet_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME().Matches(n)) {
		return false
	}
	if parent.IsCall() && parent.GetParent().IsDefaultValue() {
		if parent.GetGrandparent().IsStringKey() {
			return true
		}
		return false
	}
	if (parent.IsAssign() && (n.GetNext() != 0)) && n.GetNext().IsFunction() {
		return false
	}
	return true
}

// Java method: FindCallerLocationFunctions.visitParamListAndAddCallerLocationFunctionNames(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) priv_FindCallerLocationFunctions_visitParamListAndAddCallerLocationFunctionNames_4967(n pkg_rhino.Node, t *NodeTraversal) {
	_ = n
	_ = t
	if !(n.IsParamList()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var defaultValuesCount int32 = int32(0)
	_ = defaultValuesCount
	var paramPosition int32 = int32(0)
	_ = paramPosition
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			paramPosition = paramPosition + 1
			if !(c.IsDefaultValue()) {
				continue
			}
			defaultValuesCount = defaultValuesCount + 1
			var call pkg_rhino.Node = c.GetSecondChild()
			_ = call
			if !(call.IsCall()) {
				continue
			}
			var getProp pkg_rhino.Node = call.GetFirstChild()
			_ = getProp
			if StaticGet_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME().Matches(getProp) {
				if defaultValuesCount > int32(1) {
					this_.outer_.compiler.Report_1(JSError_Make_3(c, StaticGet_RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR(), []string{}))
				}
				var functionName string = n.GetParent().GetFirstChild().GetQualifiedName()
				_ = functionName
				if functionName == j2gort.NullString {
					this_.outer_.compiler.Report_1(JSError_Make_3(n.GetParent().GetFirstChild(), StaticGet_RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR(), []string{}))
					return
				}
				var functionVarAndPosition *RewriteCallerCodeLocation_FunctionVarAndParamPosition = NewRewriteCallerCodeLocation_FunctionVarAndParamPosition(j2gort.CastPtr[Var](t.GetScope().GetVar(functionName)), paramPosition)
				_ = functionVarAndPosition
				this_.outer_.callerLocationFunctionNames.Put(functionName, functionVarAndPosition)
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.RewriteCallerCodeLocation.RewriteCallerLocationFunctionCalls
type RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RewriteCallerCodeLocation
}

func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) IsJavaNil_() bool {
	return this_ == nil
}
func (*RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteCallerCodeLocation$RewriteCallerLocationFunctionCalls", "RewriteCallerLocationFunctionCalls"
}
func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteCallerCodeLocation.RewriteCallerLocationFunctionCalls", "com.google.javascript.jscomp.RewriteCallerCodeLocation$RewriteCallerLocationFunctionCalls", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls](v_)
	})
}
func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) SetOuter_(o *RewriteCallerCodeLocation) {
	this_.outer_ = o
}
func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) GetOuter_() *RewriteCallerCodeLocation {
	return this_.outer_
}
func (*RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) initFields_RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls() {
	_ = this_
}

// Java constructor: RewriteCallerLocationFunctionCalls()
func NewRewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls(outer_ *RewriteCallerCodeLocation) *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls {
	this_ := &RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls{}
	this_.outer_ = outer_
	InitRewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls(this_)
	return this_
}

func InitRewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls(this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls()
}

// Java method: RewriteCallerLocationFunctionCalls.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsCall() {
		this_.priv_RewriteCallerLocationFunctionCalls_visitCallNodeAndRewrite_4968(n, t)
	}
}

// Java method: RewriteCallerLocationFunctionCalls.visitCallNodeAndRewrite(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeTraversal t)
func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) priv_RewriteCallerLocationFunctionCalls_visitCallNodeAndRewrite_4968(n pkg_rhino.Node, t *NodeTraversal) {
	_ = n
	_ = t
	var firstChild pkg_rhino.Node = n.GetFirstChild()
	_ = firstChild
	if ((firstChild == 0) || (!(firstChild.IsName()) && !(firstChild.IsGetProp()))) || (firstChild.GetQualifiedName() == j2gort.NullString) {
		return
	}
	var name string = firstChild.GetQualifiedName()
	_ = name
	var moduleContentsName string = javalang.StringReplace(j2gort.JStrReplaceChar(name, uint16(46), uint16(95)), "module$exports$", "module$contents$")
	_ = moduleContentsName
	if !(this_.outer_.callerLocationFunctionNames.ContainsKey(moduleContentsName)) {
		return
	}
	var functionVarAndPosition *RewriteCallerCodeLocation_FunctionVarAndParamPosition = this_.outer_.callerLocationFunctionNames.Get(moduleContentsName)
	_ = functionVarAndPosition
	var callerLocationFunction *Var = functionVarAndPosition.GetFunctionVar()
	_ = callerLocationFunction
	var calleeFunction *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(moduleContentsName))
	_ = calleeFunction
	if ((callerLocationFunction == nil) || (calleeFunction == nil)) || !(callerLocationFunction.GetNameNode() == calleeFunction.GetNameNode()) {
		return
	}
	var positionOfCallerLocationArg int32 = functionVarAndPosition.GetParamPosition()
	_ = positionOfCallerLocationArg
	var numberOfArgs int32 = (n.GetChildCount() - int32(1))
	_ = numberOfArgs
	if numberOfArgs >= positionOfCallerLocationArg {
		var jsType pkg_rhino.JSType = n.GetChildAtIndex(positionOfCallerLocationArg).GetJSType()
		_ = jsType
		if (jsType != nil && !jsType.IsJavaNil_()) && jsType.IsExplicitlyVoidable() {
			this_.outer_.compiler.Report_1(JSError_Make_3(firstChild, StaticGet_RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR(), []string{}))
		}
		return
	}
	var xidCall pkg_rhino.Node = this_.priv_RewriteCallerLocationFunctionCalls_createGoogXidFilePathNode_4969(n)
	_ = xidCall
	this_.outer_.compiler.ReportChangeToEnclosingScope(n)
	n.AddChildToBack(xidCall)
}

// Java method: RewriteCallerLocationFunctionCalls.createGoogXidFilePathNode(com.google.javascript.rhino.Node n)
func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) priv_RewriteCallerLocationFunctionCalls_createGoogXidFilePathNode_4969(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	var googNode pkg_rhino.Node = pkg_rhino.IR_Name("goog")
	_ = googNode
	googNode.SrcrefIfMissing(n)
	var googXid pkg_rhino.Node = this_.outer_.astFactory.CreateGetPropWithUnknownType(googNode, "callerLocationIdInternalDoNotCallOrElse")
	_ = googXid
	googXid.SrcrefIfMissing(n)
	var callNode pkg_rhino.Node = this_.outer_.astFactory.CreateCallWithUnknownType(googXid, []pkg_rhino.Node{})
	_ = callNode
	callNode.SrcrefIfMissing(n)
	var stringNode pkg_rhino.Node = this_.outer_.astFactory.CreateString((j2gort.StrValString(n.GetSourceFileName()) + ":" + j2gort.StrValInt32(n.GetLineno()) + ":" + j2gort.StrValInt32(n.GetCharno())))
	_ = stringNode
	stringNode.SrcrefIfMissing(n)
	callNode.AddChildToBack(stringNode)
	return callNode
}

// Java class: com.google.javascript.jscomp.RewriteCallerCodeLocation.FunctionVarAndParamPosition
type RewriteCallerCodeLocation_FunctionVarAndParamPosition struct {
	functionVar   *Var  `gson:"functionVar"`
	paramPosition int32 `gson:"paramPosition"`
	idHash_       atomic.Int32
}

func (this_ *RewriteCallerCodeLocation_FunctionVarAndParamPosition) BindSelf_(s any) {
	_ = s
}
func (this_ *RewriteCallerCodeLocation_FunctionVarAndParamPosition) IsJavaNil_() bool {
	return this_ == nil
}
func (*RewriteCallerCodeLocation_FunctionVarAndParamPosition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewriteCallerCodeLocation$FunctionVarAndParamPosition", "FunctionVarAndParamPosition"
}
func (this_ *RewriteCallerCodeLocation_FunctionVarAndParamPosition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewriteCallerCodeLocation_FunctionVarAndParamPosition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewriteCallerCodeLocation.FunctionVarAndParamPosition", "com.google.javascript.jscomp.RewriteCallerCodeLocation$FunctionVarAndParamPosition", func(v_ any) bool {
		return j2gort.InstanceOf[*RewriteCallerCodeLocation_FunctionVarAndParamPosition](v_)
	})
}

func (this_ *RewriteCallerCodeLocation_FunctionVarAndParamPosition) initFields_RewriteCallerCodeLocation_FunctionVarAndParamPosition() {
	_ = this_
}

// Java constructor: FunctionVarAndParamPosition(com.google.javascript.jscomp.Var functionVar, int paramPosition)
func NewRewriteCallerCodeLocation_FunctionVarAndParamPosition(functionVar *Var, paramPosition int32) *RewriteCallerCodeLocation_FunctionVarAndParamPosition {
	this_ := &RewriteCallerCodeLocation_FunctionVarAndParamPosition{}
	InitRewriteCallerCodeLocation_FunctionVarAndParamPosition(this_, functionVar, paramPosition)
	return this_
}

func InitRewriteCallerCodeLocation_FunctionVarAndParamPosition(this_ *RewriteCallerCodeLocation_FunctionVarAndParamPosition, functionVar *Var, paramPosition int32) {
	_ = functionVar
	_ = paramPosition
	this_.initFields_RewriteCallerCodeLocation_FunctionVarAndParamPosition()
	this_.functionVar = functionVar
	this_.paramPosition = paramPosition
}

// Java method: FunctionVarAndParamPosition.getFunctionVar()
func (this_ *RewriteCallerCodeLocation_FunctionVarAndParamPosition) GetFunctionVar() *Var {
	return this_.functionVar
}

// Java method: FunctionVarAndParamPosition.getParamPosition()
func (this_ *RewriteCallerCodeLocation_FunctionVarAndParamPosition) GetParamPosition() int32 {
	return this_.paramPosition
}
func (*RewriteCallerCodeLocation_FindCallerLocationFunctions) AlwaysShouldTraverse_()             {}
func (*RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) AlwaysShouldTraverse_()      {}
func (*RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) CanReleaseSyntacticScopes_() {}
