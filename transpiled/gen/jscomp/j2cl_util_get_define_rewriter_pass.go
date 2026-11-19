/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/J2clUtilGetDefineRewriterPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.J2clUtilGetDefineRewriterPass
type J2clUtilGetDefineRewriterPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler                                  `gson:"compiler"`
	defines  j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"defines"`
}

func (this_ *J2clUtilGetDefineRewriterPass) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clUtilGetDefineRewriterPass) IsJavaNil_() bool { return this_ == nil }
func (*J2clUtilGetDefineRewriterPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clUtilGetDefineRewriterPass", "J2clUtilGetDefineRewriterPass"
}
func (this_ *J2clUtilGetDefineRewriterPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clUtilGetDefineRewriterPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clUtilGetDefineRewriterPass", "com.google.javascript.jscomp.J2clUtilGetDefineRewriterPass", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clUtilGetDefineRewriterPass](v_)
	})
}
func (*J2clUtilGetDefineRewriterPass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*J2clUtilGetDefineRewriterPass) Implements_jscomp_CompilerPass()                            {}
func (*J2clUtilGetDefineRewriterPass) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME *DiagnosticType

func StaticGet_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME() *DiagnosticType {
	clinit_J2clUtilGetDefineRewriterPass()
	return static_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME
}
func StaticSet_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME(v *DiagnosticType) {
	clinit_J2clUtilGetDefineRewriterPass()
	static_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME = v
}

var static_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY *DiagnosticType

func StaticGet_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY() *DiagnosticType {
	clinit_J2clUtilGetDefineRewriterPass()
	return static_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY
}
func StaticSet_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY(v *DiagnosticType) {
	clinit_J2clUtilGetDefineRewriterPass()
	static_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY = v
}

var static_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME *DiagnosticType

func StaticGet_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME() *DiagnosticType {
	clinit_J2clUtilGetDefineRewriterPass()
	return static_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME
}
func StaticSet_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME(v *DiagnosticType) {
	clinit_J2clUtilGetDefineRewriterPass()
	static_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME = v
}

var static_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE *DiagnosticType

func StaticGet_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE() *DiagnosticType {
	clinit_J2clUtilGetDefineRewriterPass()
	return static_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE
}
func StaticSet_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE(v *DiagnosticType) {
	clinit_J2clUtilGetDefineRewriterPass()
	static_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE = v
}

var clinitDone_J2clUtilGetDefineRewriterPass uint32
var clinitState_J2clUtilGetDefineRewriterPass j2gort.ClinitState

func clinit_J2clUtilGetDefineRewriterPass() {
	if atomic.LoadUint32(&clinitDone_J2clUtilGetDefineRewriterPass) != 0 {
		return
	}
	clinitSlow_J2clUtilGetDefineRewriterPass()
}

func clinitSlow_J2clUtilGetDefineRewriterPass() {
	if !j2gort.EnterClinit(&clinitState_J2clUtilGetDefineRewriterPass, &clinitDone_J2clUtilGetDefineRewriterPass) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_J2clUtilGetDefineRewriterPass, &clinitDone_J2clUtilGetDefineRewriterPass)
	static_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME = DiagnosticType_Error("JSC_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME", "Calls to System.getProperty must use a String literal for the property name")
	static_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY = DiagnosticType_Error("JSC_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY", "Unknown system property \"{0}\", are you missing a call to jre.addSystemPropertyFromGoogDefine?")
	static_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME = DiagnosticType_Error("JSC_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME", "Calls to addSystemPropertyFromGoogDefine must use a String literal for the define name")
	static_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE = DiagnosticType_Error("JSC_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE", "Unknown goog.define with name: {0}")
}

func (this_ *J2clUtilGetDefineRewriterPass) initFields_J2clUtilGetDefineRewriterPass() {
	_ = this_
}

// Java constructor: J2clUtilGetDefineRewriterPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewJ2clUtilGetDefineRewriterPass(compiler AbstractCompiler) *J2clUtilGetDefineRewriterPass {
	clinit_J2clUtilGetDefineRewriterPass()
	this_ := &J2clUtilGetDefineRewriterPass{}
	InitJ2clUtilGetDefineRewriterPass(this_, compiler)
	return this_
}

func InitJ2clUtilGetDefineRewriterPass(this_ *J2clUtilGetDefineRewriterPass, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clUtilGetDefineRewriterPass()
	this_.compiler = compiler
}

// Java method: J2clUtilGetDefineRewriterPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *J2clUtilGetDefineRewriterPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(J2clSourceFileChecker_ShouldRunJ2clPasses(this_.compiler)) {
		return
	}
	this_.defines = this_.compiler.GetDefineNames()
	var checkGetDefineCalls *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls = NewJ2clUtilGetDefineRewriterPass_CheckGetDefineCalls(this_)
	_ = checkGetDefineCalls
	NodeTraversal_Traverse(this_.compiler, root, checkGetDefineCalls)
	checkGetDefineCalls.priv_CheckGetDefineCalls_validateRegisteredDefines_3500()
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: J2clUtilGetDefineRewriterPass.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *J2clUtilGetDefineRewriterPass) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if J2clUtilGetDefineRewriterPass_IsUtilGetDefineCall(n) {
		this_.priv_J2clUtilGetDefineRewriterPass_rewriteUtilGetDefine_3498(t, n)
	}
}

// Java method: J2clUtilGetDefineRewriterPass.rewriteUtilGetDefine(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode)
func (this_ *J2clUtilGetDefineRewriterPass) priv_J2clUtilGetDefineRewriterPass_rewriteUtilGetDefine_3498(t *NodeTraversal, callNode pkg_rhino.Node) {
	_ = t
	_ = callNode
	var firstExpr pkg_rhino.Node = callNode.GetSecondChild()
	_ = firstExpr
	var secondExpr pkg_rhino.Node = callNode.GetLastChild()
	_ = secondExpr
	if secondExpr != firstExpr {
		secondExpr.Detach()
	} else {
		secondExpr = pkg_rhino.IR_NullNode()
	}
	var replacement pkg_rhino.Node = this_.priv_J2clUtilGetDefineRewriterPass_getDefineReplacement_3499(firstExpr, secondExpr)
	_ = replacement
	replacement.SrcrefTreeIfMissing(callNode)
	callNode.ReplaceWith(replacement)
	t.ReportCodeChange_0()
}

// Java method: J2clUtilGetDefineRewriterPass.getDefineReplacement(com.google.javascript.rhino.Node firstExpr, com.google.javascript.rhino.Node secondExpr)
func (this_ *J2clUtilGetDefineRewriterPass) priv_J2clUtilGetDefineRewriterPass_getDefineReplacement_3499(firstExpr pkg_rhino.Node, secondExpr pkg_rhino.Node) pkg_rhino.Node {
	_ = firstExpr
	_ = secondExpr
	var defineName string = firstExpr.GetString()
	_ = defineName
	if j2gort.SetContainsStr(this_.defines, defineName) {
		var define pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, ProcessDefines_GetGlobalDefineAlias(defineName))
		_ = define
		var defineStringValue pkg_rhino.Node = NodeUtil_NewCallNode(pkg_rhino.IR_Name("String"), []pkg_rhino.Node{define})
		_ = defineStringValue
		return pkg_rhino.IR_Comma(secondExpr, defineStringValue)
	} else {
		return secondExpr
	}
}

// Java static method: J2clUtilGetDefineRewriterPass.isUtilGetDefineCall(com.google.javascript.rhino.Node n)
func J2clUtilGetDefineRewriterPass_IsUtilGetDefineCall(n pkg_rhino.Node) bool {
	clinit_J2clUtilGetDefineRewriterPass()
	_ = n
	return (n.IsCall() && J2clUtilGetDefineRewriterPass_IsUtilGetDefineMethodName(n.GetFirstChild().GetQualifiedName()))
}

// Java static method: J2clUtilGetDefineRewriterPass.isUtilGetDefineMethodName(java.lang.String fnName)
func J2clUtilGetDefineRewriterPass_IsUtilGetDefineMethodName(fnName string) bool {
	clinit_J2clUtilGetDefineRewriterPass()
	_ = fnName
	return (((fnName != j2gort.NullString) && javalang.StringEndsWith(j2gort.CheckNotNullString(fnName), ".$getDefine")) && (j2gort.JStrIndexOf(fnName, "Util") >= 0))
}

// Java class: com.google.javascript.jscomp.J2clUtilGetDefineRewriterPass.CheckGetDefineCalls
type J2clUtilGetDefineRewriterPass_CheckGetDefineCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                            *J2clUtilGetDefineRewriterPass
	defineNameToAddSystemPropertyCall j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"defineNameToAddSystemPropertyCall"`
	defineNameToGetDefineCall         j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"defineNameToGetDefineCall"`
}

func (this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) BindSelf_(s any) {
	_ = s
}
func (this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) IsJavaNil_() bool {
	return this_ == nil
}
func (*J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.J2clUtilGetDefineRewriterPass$CheckGetDefineCalls", "CheckGetDefineCalls"
}
func (this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.J2clUtilGetDefineRewriterPass.CheckGetDefineCalls", "com.google.javascript.jscomp.J2clUtilGetDefineRewriterPass$CheckGetDefineCalls", func(v_ any) bool {
		return j2gort.InstanceOf[*J2clUtilGetDefineRewriterPass_CheckGetDefineCalls](v_)
	})
}
func (this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) SetOuter_(o *J2clUtilGetDefineRewriterPass) {
	this_.outer_ = o
}
func (this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) GetOuter_() *J2clUtilGetDefineRewriterPass {
	return this_.outer_
}
func (*J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) initFields_J2clUtilGetDefineRewriterPass_CheckGetDefineCalls() {
	_ = this_
	this_.defineNameToAddSystemPropertyCall = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	this_.defineNameToGetDefineCall = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
}

// Java constructor: CheckGetDefineCalls()
func NewJ2clUtilGetDefineRewriterPass_CheckGetDefineCalls(outer_ *J2clUtilGetDefineRewriterPass) *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls {
	this_ := &J2clUtilGetDefineRewriterPass_CheckGetDefineCalls{}
	this_.outer_ = outer_
	InitJ2clUtilGetDefineRewriterPass_CheckGetDefineCalls(this_)
	return this_
}

func InitJ2clUtilGetDefineRewriterPass_CheckGetDefineCalls(this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_J2clUtilGetDefineRewriterPass_CheckGetDefineCalls()
}

// Java method: CheckGetDefineCalls.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if J2clUtilGetDefineRewriterPass_CheckGetDefineCalls_IsAddSystemPropertyFromGoogDefine(n) {
		var defineName string = J2clUtilGetDefineRewriterPass_CheckGetDefineCalls_GetDefineNameParameter(n)
		_ = defineName
		if defineName != j2gort.NullString {
			this_.defineNameToAddSystemPropertyCall.Put(defineName, n)
		} else {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_CONSTANT_NAME(), []string{}))
		}
	} else if J2clUtilGetDefineRewriterPass_IsUtilGetDefineCall(n) {
		var defineName_2 string = J2clUtilGetDefineRewriterPass_CheckGetDefineCalls_GetDefineNameParameter(n)
		_ = defineName_2
		if defineName_2 != j2gort.NullString {
			this_.defineNameToGetDefineCall.Put(defineName_2, n)
		} else {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_CONSTANT_NAME(), []string{}))
		}
	}
}

// Java method: CheckGetDefineCalls.validateRegisteredDefines()
func (this_ *J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) priv_CheckGetDefineCalls_validateRegisteredDefines_3500() {
	{
		cur_1 := this_.defineNameToAddSystemPropertyCall.IterEntries()
		for cur_1.NextEntry() {
			if !(j2gort.SetContainsStr(this_.outer_.defines, cur_1.Key())) {
				this_.outer_.compiler.Report_1(JSError_Make_3(cur_1.Value(), StaticGet_J2clUtilGetDefineRewriterPass_J2CL_ADD_SYSTEM_PROPERTY_UNKNOWN_DEFINE(), []string{cur_1.Key()}))
			}
		}
	}
	{
		cur_2 := this_.defineNameToGetDefineCall.IterEntries()
		for cur_2.NextEntry() {
			if !(j2gort.SetContainsStr(this_.outer_.defines, cur_2.Key())) {
				continue
			}
			if !(this_.defineNameToAddSystemPropertyCall.ContainsKey(cur_2.Key())) {
				this_.outer_.compiler.Report_1(JSError_Make_3(cur_2.Value(), StaticGet_J2clUtilGetDefineRewriterPass_J2CL_SYSTEM_GET_PROPERTY_UNKNOWN_PROPERTY(), []string{cur_2.Key()}))
			}
		}
	}
}

// Java static method: CheckGetDefineCalls.getDefineNameParameter(com.google.javascript.rhino.Node callNode)
func J2clUtilGetDefineRewriterPass_CheckGetDefineCalls_GetDefineNameParameter(callNode pkg_rhino.Node) string {
	_ = callNode
	var nameNode pkg_rhino.Node = callNode.GetSecondChild()
	_ = nameNode
	if (nameNode == 0) || !(nameNode.IsString()) {
		return j2gort.NullString
	}
	return nameNode.GetString()
}

// Java static method: CheckGetDefineCalls.isAddSystemPropertyFromGoogDefine(com.google.javascript.rhino.Node n)
func J2clUtilGetDefineRewriterPass_CheckGetDefineCalls_IsAddSystemPropertyFromGoogDefine(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsCall() && (n.GetFirstChild().GetQualifiedName() != j2gort.NullString)) && javalang.StringEndsWith(j2gort.CheckNotNullString(n.GetFirstChild().GetQualifiedName()), "addSystemPropertyFromGoogDefine"))
}
func (*J2clUtilGetDefineRewriterPass) AlwaysShouldTraverse_()                          {}
func (*J2clUtilGetDefineRewriterPass) CanReleaseSyntacticScopes_()                     {}
func (*J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) AlwaysShouldTraverse_()      {}
func (*J2clUtilGetDefineRewriterPass_CheckGetDefineCalls) CanReleaseSyntacticScopes_() {}
