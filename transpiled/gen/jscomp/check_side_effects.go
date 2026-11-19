/*
 * Copyright 2006 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CheckSideEffects.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CheckSideEffects
type CheckSideEffects struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	problemNodes              []pkg_rhino.Node          `gson:"problemNodes"`
	noSideEffectExterns       j2gort.OrderedSet[string] `gson:"noSideEffectExterns"`
	compiler                  AbstractCompiler          `gson:"compiler"`
	report                    bool                      `gson:"report"`
	protectSideEffectFreeCode bool                      `gson:"protectSideEffectFreeCode"`
	preserveFunctionInjected  bool                      `gson:"preserveFunctionInjected"`
}

func (this_ *CheckSideEffects) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSideEffects) IsJavaNil_() bool { return this_ == nil }
func (*CheckSideEffects) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSideEffects", "CheckSideEffects"
}
func (this_ *CheckSideEffects) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckSideEffects) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSideEffects", "com.google.javascript.jscomp.CheckSideEffects", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSideEffects](v_)
	})
}
func (*CheckSideEffects) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckSideEffects) Implements_jscomp_CompilerPass()                            {}
func (*CheckSideEffects) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckSideEffects_USELESS_CODE_ERROR *DiagnosticType

func StaticGet_CheckSideEffects_USELESS_CODE_ERROR() *DiagnosticType {
	clinit_CheckSideEffects()
	return static_CheckSideEffects_USELESS_CODE_ERROR
}
func StaticSet_CheckSideEffects_USELESS_CODE_ERROR(v *DiagnosticType) {
	clinit_CheckSideEffects()
	static_CheckSideEffects_USELESS_CODE_ERROR = v
}
func StaticGet_CheckSideEffects_PROTECTOR_FN() string      { return "JSCOMPILER_PRESERVE" }
func StaticGet_CheckSideEffects_J2CL_PROTECTOR_FN() string { return "$J2CL_PRESERVE$" }

var clinitDone_CheckSideEffects uint32
var clinitState_CheckSideEffects j2gort.ClinitState

func clinit_CheckSideEffects() {
	if atomic.LoadUint32(&clinitDone_CheckSideEffects) != 0 {
		return
	}
	clinitSlow_CheckSideEffects()
}

func clinitSlow_CheckSideEffects() {
	if !j2gort.EnterClinit(&clinitState_CheckSideEffects, &clinitDone_CheckSideEffects) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckSideEffects, &clinitDone_CheckSideEffects)
	static_CheckSideEffects_USELESS_CODE_ERROR = DiagnosticType_Warning("JSC_USELESS_CODE", "Suspicious code. {0}")
}

func (this_ *CheckSideEffects) initFields_CheckSideEffects() {
	_ = this_
	this_.problemNodes = []pkg_rhino.Node{}
	this_.noSideEffectExterns = j2gort.NewOrderedSet[string]()

}

// Java constructor: CheckSideEffects(com.google.javascript.jscomp.AbstractCompiler compiler, boolean report, boolean protectSideEffectFreeCode)
func NewCheckSideEffects(compiler AbstractCompiler, report bool, protectSideEffectFreeCode bool) *CheckSideEffects {
	clinit_CheckSideEffects()
	this_ := &CheckSideEffects{}
	InitCheckSideEffects(this_, compiler, report, protectSideEffectFreeCode)
	return this_
}

func InitCheckSideEffects(this_ *CheckSideEffects, compiler AbstractCompiler, report bool, protectSideEffectFreeCode bool) {
	_ = compiler
	_ = report
	_ = protectSideEffectFreeCode
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckSideEffects()
	this_.compiler = compiler
	this_.report = report
	this_.protectSideEffectFreeCode = protectSideEffectFreeCode
}

// Java method: CheckSideEffects.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckSideEffects) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, externs, NewCheckSideEffects_GetNoSideEffectExterns(this_))
	NodeTraversal_Traverse(this_.compiler, root, this_)
	if this_.protectSideEffectFreeCode {
		this_.priv_CheckSideEffects_protectSideEffects_1105()
	}
}

// Java method: CheckSideEffects.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckSideEffects) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsEmpty() || n.IsComma() {
		return
	}
	if parent == 0 {
		return
	}
	if n.IsExprResult() || n.IsBlock() {
		return
	}
	if n.IsQualifiedName() && (n.GetJSDocInfo() != nil) {
		return
	}
	if (((n.IsStringLit() && n.GetParent().IsExprResult()) && n.GetGrandparent().IsScript()) && javalang.StringStartsWith(j2gort.CheckNotNullString(n.GetSourceFileName()), "src/com/google/javascript/jscomp/js/")) && javalang.StringStartsWith(j2gort.CheckNotNullString(n.GetString()), "require ") {
		return
	}
	var isResultUsed bool = NodeUtil_IsExpressionResultUsed(n)
	_ = isResultUsed
	var isSimpleOp bool = ((NodeUtil_IsSimpleOperator(n) || n.IsGetProp()) || n.IsGetElem())
	_ = isSimpleOp
	if !(isResultUsed) && !(n.IsVoid()) {
		if isSimpleOp || !(t.GetCompiler().GetAstAnalyzer().MayHaveSideEffects(n)) {
			if this_.report {
				var msg string = "This code lacks side-effects. Is there a bug?"
				_ = msg
				if n.IsStringLit() || n.IsTemplateLit() {
					msg = "Is there a missing '+' on the previous line?"
				} else if isSimpleOp {
					msg = ("The result of the '" + j2gort.StrValString(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(j2gort.CheckNotNullEnum(n.GetToken()).String())) + "' operator is not being used.")
				}
				t.Report_3(n, StaticGet_CheckSideEffects_USELESS_CODE_ERROR(), []string{msg})
			}
			if !(NodeUtil_IsStatement(n)) {
				this_.problemNodes = append(this_.problemNodes, n)
			}
		} else if n.IsCall() && ((n.GetFirstChild().IsGetProp() || n.GetFirstChild().IsName()) || n.GetFirstChild().IsStringLit()) {
			var qname string = n.GetFirstChild().GetQualifiedName()
			_ = qname
			var isDefinedInSrc bool = false
			_ = isDefinedInSrc
			if qname != j2gort.NullString {
				if n.GetFirstChild().IsGetProp() {
					var rootNameNode pkg_rhino.Node = NodeUtil_GetRootOfQualifiedName_1_Node(n.GetFirstChild())
					_ = rootNameNode
					isDefinedInSrc = (((rootNameNode != 0) && rootNameNode.IsName()) && (j2gort.CastPtr[Var](t.GetScope().GetVar(rootNameNode.GetString())) != nil))
				} else {
					isDefinedInSrc = (j2gort.CastPtr[Var](t.GetScope().GetVar(qname)) != nil)
				}
			}
			if ((qname != j2gort.NullString) && this_.noSideEffectExterns.Contains(qname)) && !(isDefinedInSrc) {
				this_.problemNodes = append(this_.problemNodes, n)
				if this_.report {
					var msg_2 string = ("The result of the extern function call '" + j2gort.StrValString(qname) + "' is not being used.")
					_ = msg_2
					t.Report_3(n, StaticGet_CheckSideEffects_USELESS_CODE_ERROR(), []string{msg_2})
				}
			}
		}
	}
}

// Java method: CheckSideEffects.protectSideEffects()
func (this_ *CheckSideEffects) priv_CheckSideEffects_protectSideEffects_1105() {
	if len(this_.problemNodes) != 0 {
		if !(this_.preserveFunctionInjected) {
			CheckSideEffects_AddExtern(this_.compiler)
		}
		for _, n := range this_.problemNodes {
			_ = n
			var name pkg_rhino.Node = pkg_rhino.IR_Name("JSCOMPILER_PRESERVE").Srcref(n)
			_ = name
			name.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
			var replacement pkg_rhino.Node = pkg_rhino.IR_Call(name, []pkg_rhino.Node{}).Srcref(n)
			_ = replacement
			replacement.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
			n.ReplaceWith(replacement)
			replacement.AddChildToBack(n)
			this_.compiler.ReportChangeToEnclosingScope(replacement)
		}

	}
}

// Java static method: CheckSideEffects.addExtern(com.google.javascript.jscomp.AbstractCompiler compiler)
func CheckSideEffects_AddExtern(compiler AbstractCompiler) {
	clinit_CheckSideEffects()
	_ = compiler
	NodeUtil_CreateSynthesizedExternsSymbol(compiler, "JSCOMPILER_PRESERVE")
}

// Java class: com.google.javascript.jscomp.CheckSideEffects.StripProtection
type CheckSideEffects_StripProtection struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckSideEffects_StripProtection) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSideEffects_StripProtection) IsJavaNil_() bool { return this_ == nil }
func (*CheckSideEffects_StripProtection) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSideEffects$StripProtection", "StripProtection"
}
func (this_ *CheckSideEffects_StripProtection) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckSideEffects_StripProtection) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSideEffects.StripProtection", "com.google.javascript.jscomp.CheckSideEffects$StripProtection", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSideEffects_StripProtection](v_)
	})
}
func (*CheckSideEffects_StripProtection) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*CheckSideEffects_StripProtection) Implements_jscomp_CompilerPass()           {}
func (*CheckSideEffects_StripProtection) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *CheckSideEffects_StripProtection) initFields_CheckSideEffects_StripProtection() {
	_ = this_
}

// Java constructor: StripProtection(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckSideEffects_StripProtection(compiler AbstractCompiler) *CheckSideEffects_StripProtection {
	this_ := &CheckSideEffects_StripProtection{}
	InitCheckSideEffects_StripProtection(this_, compiler)
	return this_
}

func InitCheckSideEffects_StripProtection(this_ *CheckSideEffects_StripProtection, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckSideEffects_StripProtection()
	this_.compiler = compiler
}

// Java method: StripProtection.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckSideEffects_StripProtection) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: StripProtection.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckSideEffects_StripProtection) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsCall() {
		var target pkg_rhino.Node = n.GetFirstChild()
		_ = target
		if !(target.IsName()) {
			return
		}
		{
			var iifeCond_1 bool
			{

				sr_2 := target.GetString()
				se_3 := "$J2CL_PRESERVE$"
				iifeCond_1 = j2gort.CheckNotNullString(sr_2) == se_3
			}
			if iifeCond_1 {
				if !(J2clSourceUtils_IsJ2clSource_1_string(n.GetSourceFileName())) {
					j2gort.CheckStateStr(false, "Only allowed for J2CL code")
				}
				NodeUtil_DeleteFunctionCall(n, this_.compiler)
				return
			}
		}
		{
			var iifeCond_1 bool
			{

				sr_4 := target.GetString()
				se_5 := "JSCOMPILER_PRESERVE"
				iifeCond_1 = j2gort.CheckNotNullString(sr_4) == se_5
			}
			if iifeCond_1 {
				var expr pkg_rhino.Node = n.GetLastChild()
				_ = expr
				n.DetachChildren()
				n.ReplaceWith(expr)
				t.ReportCodeChange_0()
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.CheckSideEffects.GetNoSideEffectExterns
type CheckSideEffects_GetNoSideEffectExterns struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *CheckSideEffects
}

func (this_ *CheckSideEffects_GetNoSideEffectExterns) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckSideEffects_GetNoSideEffectExterns) IsJavaNil_() bool { return this_ == nil }
func (*CheckSideEffects_GetNoSideEffectExterns) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckSideEffects$GetNoSideEffectExterns", "GetNoSideEffectExterns"
}
func (this_ *CheckSideEffects_GetNoSideEffectExterns) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckSideEffects_GetNoSideEffectExterns) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckSideEffects.GetNoSideEffectExterns", "com.google.javascript.jscomp.CheckSideEffects$GetNoSideEffectExterns", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckSideEffects_GetNoSideEffectExterns](v_)
	})
}
func (this_ *CheckSideEffects_GetNoSideEffectExterns) SetOuter_(o *CheckSideEffects) {
	this_.outer_ = o
}
func (this_ *CheckSideEffects_GetNoSideEffectExterns) GetOuter_() *CheckSideEffects {
	return this_.outer_
}
func (*CheckSideEffects_GetNoSideEffectExterns) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*CheckSideEffects_GetNoSideEffectExterns) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *CheckSideEffects_GetNoSideEffectExterns) initFields_CheckSideEffects_GetNoSideEffectExterns() {
	_ = this_
}

// Java constructor: GetNoSideEffectExterns()
func NewCheckSideEffects_GetNoSideEffectExterns(outer_ *CheckSideEffects) *CheckSideEffects_GetNoSideEffectExterns {
	this_ := &CheckSideEffects_GetNoSideEffectExterns{}
	this_.outer_ = outer_
	InitCheckSideEffects_GetNoSideEffectExterns(this_)
	return this_
}

func InitCheckSideEffects_GetNoSideEffectExterns(this_ *CheckSideEffects_GetNoSideEffectExterns) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckSideEffects_GetNoSideEffectExterns()
}

// Java method: GetNoSideEffectExterns.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckSideEffects_GetNoSideEffectExterns) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsFunction() {
		var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = jsDoc
		if (jsDoc != nil) && jsDoc.IsNoSideEffects() {
			var name string = NodeUtil_GetName(n)
			_ = name
			this_.outer_.noSideEffectExterns.Add(name)
		}
	}
	if n.IsName() && func() bool {
		sr_6 := "JSCOMPILER_PRESERVE"
		se_7 := n.GetString()
		return j2gort.CheckNotNullString(sr_6) == se_7
	}() {
		this_.outer_.preserveFunctionInjected = true
	}
}
func (*CheckSideEffects) AlwaysShouldTraverse_()                             {}
func (*CheckSideEffects) CanReleaseSyntacticScopes_()                        {}
func (*CheckSideEffects_StripProtection) AlwaysShouldTraverse_()             {}
func (*CheckSideEffects_StripProtection) CanReleaseSyntacticScopes_()        {}
func (*CheckSideEffects_GetNoSideEffectExterns) AlwaysShouldTraverse_()      {}
func (*CheckSideEffects_GetNoSideEffectExterns) CanReleaseSyntacticScopes_() {}
