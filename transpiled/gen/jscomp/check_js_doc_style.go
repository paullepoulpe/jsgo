/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckJSDocStyle.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.lint.CheckJSDocStyle
type CheckJSDocStyle struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckJSDocStyle) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckJSDocStyle) IsJavaNil_() bool { return this_ == nil }
func (*CheckJSDocStyle) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckJSDocStyle", "CheckJSDocStyle"
}
func (this_ *CheckJSDocStyle) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckJSDocStyle) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckJSDocStyle", "com.google.javascript.jscomp.lint.CheckJSDocStyle", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckJSDocStyle](v_)
	})
}
func (*CheckJSDocStyle) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckJSDocStyle) Implements_jscomp_CompilerPass()                            {}
func (*CheckJSDocStyle) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckJSDocStyle_CLASS_DISALLOWED_JSDOC *DiagnosticType

func StaticGet_CheckJSDocStyle_CLASS_DISALLOWED_JSDOC() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_CLASS_DISALLOWED_JSDOC
}
func StaticSet_CheckJSDocStyle_CLASS_DISALLOWED_JSDOC(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_CLASS_DISALLOWED_JSDOC = v
}

var static_CheckJSDocStyle_MISSING_JSDOC *DiagnosticType

func StaticGet_CheckJSDocStyle_MISSING_JSDOC() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_MISSING_JSDOC
}
func StaticSet_CheckJSDocStyle_MISSING_JSDOC(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_MISSING_JSDOC = v
}

var static_CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER *DiagnosticType

func StaticGet_CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER
}
func StaticSet_CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER = v
}

var static_CheckJSDocStyle_MISSING_PARAMETER_JSDOC *DiagnosticType

func StaticGet_CheckJSDocStyle_MISSING_PARAMETER_JSDOC() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_MISSING_PARAMETER_JSDOC
}
func StaticSet_CheckJSDocStyle_MISSING_PARAMETER_JSDOC(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_MISSING_PARAMETER_JSDOC = v
}

var static_CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES *DiagnosticType

func StaticGet_CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES
}
func StaticSet_CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES = v
}

var static_CheckJSDocStyle_MISSING_RETURN_JSDOC *DiagnosticType

func StaticGet_CheckJSDocStyle_MISSING_RETURN_JSDOC() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_MISSING_RETURN_JSDOC
}
func StaticSet_CheckJSDocStyle_MISSING_RETURN_JSDOC(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_MISSING_RETURN_JSDOC = v
}

var static_CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL *DiagnosticType

func StaticGet_CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL
}
func StaticSet_CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL = v
}

var static_CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS *DiagnosticType

func StaticGet_CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS
}
func StaticSet_CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS = v
}

var static_CheckJSDocStyle_INCORRECT_PARAM_NAME *DiagnosticType

func StaticGet_CheckJSDocStyle_INCORRECT_PARAM_NAME() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_INCORRECT_PARAM_NAME
}
func StaticSet_CheckJSDocStyle_INCORRECT_PARAM_NAME(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_INCORRECT_PARAM_NAME = v
}

var static_CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED *DiagnosticType

func StaticGet_CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED
}
func StaticSet_CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED = v
}

var static_CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS *DiagnosticType

func StaticGet_CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS
}
func StaticSet_CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS = v
}

var static_CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE *DiagnosticType

func StaticGet_CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE() *DiagnosticType {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE
}
func StaticSet_CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE(v *DiagnosticType) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE = v
}

var static_CheckJSDocStyle_LINT_DIAGNOSTICS *DiagnosticGroup

func StaticGet_CheckJSDocStyle_LINT_DIAGNOSTICS() *DiagnosticGroup {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_LINT_DIAGNOSTICS
}
func StaticSet_CheckJSDocStyle_LINT_DIAGNOSTICS(v *DiagnosticGroup) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_LINT_DIAGNOSTICS = v
}

var static_CheckJSDocStyle_ALL_DIAGNOSTICS *DiagnosticGroup

func StaticGet_CheckJSDocStyle_ALL_DIAGNOSTICS() *DiagnosticGroup {
	clinit_CheckJSDocStyle()
	return static_CheckJSDocStyle_ALL_DIAGNOSTICS
}
func StaticSet_CheckJSDocStyle_ALL_DIAGNOSTICS(v *DiagnosticGroup) {
	clinit_CheckJSDocStyle()
	static_CheckJSDocStyle_ALL_DIAGNOSTICS = v
}

var clinitDone_CheckJSDocStyle uint32
var clinitState_CheckJSDocStyle j2gort.ClinitState

func clinit_CheckJSDocStyle() {
	if atomic.LoadUint32(&clinitDone_CheckJSDocStyle) != 0 {
		return
	}
	clinitSlow_CheckJSDocStyle()
}

func clinitSlow_CheckJSDocStyle() {
	if !j2gort.EnterClinit(&clinitState_CheckJSDocStyle, &clinitDone_CheckJSDocStyle) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckJSDocStyle, &clinitDone_CheckJSDocStyle)
	static_CheckJSDocStyle_CLASS_DISALLOWED_JSDOC = DiagnosticType_Disabled("JSC_CLASS_DISALLOWED_JSDOC", "@constructor annotations are redundant on classes.")
	static_CheckJSDocStyle_MISSING_JSDOC = DiagnosticType_Disabled("JSC_MISSING_JSDOC", "Function must have JSDoc.")
	static_CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER = DiagnosticType_Disabled("JSC_TYPE_ON_GETTER_SETTER", "Getters and setters must not have @type annotations. Did you mean @return or @param instead?")
	static_CheckJSDocStyle_MISSING_PARAMETER_JSDOC = DiagnosticType_Disabled("JSC_MISSING_PARAMETER_JSDOC", "Parameter must have JSDoc.{0}")
	static_CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES = DiagnosticType_Disabled("JSC_MIXED_PARAM_JSDOC_STYLES", "Functions may not use both @param annotations and inline JSDoc")
	static_CheckJSDocStyle_MISSING_RETURN_JSDOC = DiagnosticType_Disabled("JSC_MISSING_RETURN_JSDOC", "Function that returns a value must have JSDoc indicating the return type.{0}")
	static_CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL = DiagnosticType_Disabled("JSC_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL", "Parameter {0} is optional so it must have a JSDoc type ending with ''=''")
	static_CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS = DiagnosticType_Disabled("JSC_WRONG_NUMBER_OF_PARAMS", "Wrong number of @param annotations")
	static_CheckJSDocStyle_INCORRECT_PARAM_NAME = DiagnosticType_Disabled("JSC_INCORRECT_PARAM_NAME", "Incorrect param name. Are your @param annotations in the wrong order?")
	static_CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED = DiagnosticType_Disabled("JSC_EXTERNS_FILES_SHOULD_BE_ANNOTATED", "Externs files should be annotated with @externs in the @fileoverview block.")
	static_CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS = DiagnosticType_Disabled("JSC_LICENSE_CONTAINS_AT_EXTERNS", "@license block contains an @externs annotation, which will be parsed as plain license text instead of an actual @externs annotation. You probably meant to put @externs in a separate @fileoverview block.")
	static_CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE = DiagnosticType_Disabled("JSC_PREFER_BACKTICKS_TO_AT_SIGN_CODE", "Use `some_code` instead of '{'@code some_code'}'.")
	static_CheckJSDocStyle_LINT_DIAGNOSTICS = NewDiagnosticGroup_1_DiagnosticType_2([]*DiagnosticType{StaticGet_CheckJSDocStyle_CLASS_DISALLOWED_JSDOC(), StaticGet_CheckJSDocStyle_MISSING_JSDOC(), StaticGet_CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER(), StaticGet_CheckJSDocStyle_MISSING_PARAMETER_JSDOC(), StaticGet_CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES(), StaticGet_CheckJSDocStyle_MISSING_RETURN_JSDOC(), StaticGet_CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL(), StaticGet_CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS(), StaticGet_CheckJSDocStyle_INCORRECT_PARAM_NAME(), StaticGet_CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED(), StaticGet_CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS(), StaticGet_CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE()})
	static_CheckJSDocStyle_ALL_DIAGNOSTICS = NewDiagnosticGroup_1_DiagnosticGroup([]*DiagnosticGroup{StaticGet_CheckJSDocStyle_LINT_DIAGNOSTICS()})
}

func (this_ *CheckJSDocStyle) initFields_CheckJSDocStyle() {
	_ = this_
}

// Java constructor: CheckJSDocStyle(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckJSDocStyle(compiler AbstractCompiler) *CheckJSDocStyle {
	clinit_CheckJSDocStyle()
	this_ := &CheckJSDocStyle{}
	InitCheckJSDocStyle(this_, compiler)
	return this_
}

func InitCheckJSDocStyle(this_ *CheckJSDocStyle, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckJSDocStyle()
	this_.compiler = compiler
}

// Java method: CheckJSDocStyle.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckJSDocStyle) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	NodeTraversal_Traverse(this_.compiler, externs, NewCheckJSDocStyle_ExternsCallback())
}

// Java method: CheckJSDocStyle.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node unused)
func (this_ *CheckJSDocStyle) Visit_3(t *NodeTraversal, n pkg_rhino.Node, unused pkg_rhino.Node) {
	_ = t
	_ = n
	_ = unused
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_CheckJSDocStyle_visitFunction_6573(t, n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_CheckJSDocStyle_visitClass_6574(t, n)
	case int32(51) /* Token.ASSIGN */ :
		this_.priv_CheckJSDocStyle_checkStyleForPrivateProperties_6572(t, n)
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(100) /* Token.STRING_KEY */ :
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		if NodeUtil_GetEnclosingClass(n) != 0 {
			this_.priv_CheckJSDocStyle_checkStyleForPrivateProperties_6572(t, n)
		}
	case int32(94) /* Token.SCRIPT */ :
		this_.priv_CheckJSDocStyle_checkLicenseComment_6584(t, n)
	default:
		this_.priv_CheckJSDocStyle_visitNonFunction_6571(t, n)
	}
}

// Java method: CheckJSDocStyle.checkForAtSignCodePresence(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo jsDoc)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_checkForAtSignCodePresence_6570(t *NodeTraversal, n pkg_rhino.Node, jsDoc *pkg_rhino.JSDocInfo) {
	_ = t
	_ = n
	_ = jsDoc
	if jsDoc == nil {
		return
	}
	if jsDoc.IsAtSignCodePresent() {
		t.Report_3(n, StaticGet_CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE(), []string{})
	}
}

// Java method: CheckJSDocStyle.visitNonFunction(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_visitNonFunction_6571(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var jsDoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsDoc
	this_.priv_CheckJSDocStyle_checkForAtSignCodePresence_6570(t, n, jsDoc)
}

// Java method: CheckJSDocStyle.checkStyleForPrivateProperties(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_checkStyleForPrivateProperties_6572(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = jsDoc
	this_.priv_CheckJSDocStyle_checkForAtSignCodePresence_6570(t, n, jsDoc)
}

// Java static method: CheckJSDocStyle.checkNoTypeOnGettersAndSetters(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node function, com.google.javascript.rhino.JSDocInfo jsDoc)
func CheckJSDocStyle_CheckNoTypeOnGettersAndSetters(t *NodeTraversal, function pkg_rhino.Node, jsDoc *pkg_rhino.JSDocInfo) {
	clinit_CheckJSDocStyle()
	_ = t
	_ = function
	_ = jsDoc
	if function.GetGrandparent().IsClassMembers() {
		var memberNode pkg_rhino.Node = function.GetParent()
		_ = memberNode
		if memberNode.IsSetterDef() || memberNode.IsGetterDef() {
			if (jsDoc != nil) && jsDoc.HasType() {
				t.Report_3(function, StaticGet_CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER(), []string{})
			}
		}
	}
}

// Java method: CheckJSDocStyle.visitFunction(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node function)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_visitFunction_6573(t *NodeTraversal, function pkg_rhino.Node) {
	_ = t
	_ = function
	var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(function)
	_ = jsDoc
	this_.priv_CheckJSDocStyle_checkForAtSignCodePresence_6570(t, function, jsDoc)
	if (jsDoc == nil) && !(this_.priv_CheckJSDocStyle_hasAnyInlineJsDoc_6582(function)) {
		this_.priv_CheckJSDocStyle_checkMissingJsDoc_6575(t, function)
	} else {
		if (((t.InGlobalScope() || this_.priv_CheckJSDocStyle_hasAnyInlineJsDoc_6582(function)) || (jsDoc.GetParameterCount() > int32(0))) || jsDoc.HasReturnType()) || jsDoc.IsOverride() {
			this_.priv_CheckJSDocStyle_checkParams_6579(t, function, jsDoc)
		}
		CheckJSDocStyle_CheckNoTypeOnGettersAndSetters(t, function, jsDoc)
		this_.priv_CheckJSDocStyle_checkReturn_6583(function, jsDoc)
	}
}

// Java method: CheckJSDocStyle.visitClass(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node cls)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_visitClass_6574(t *NodeTraversal, cls pkg_rhino.Node) {
	_ = t
	_ = cls
	var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(cls)
	_ = jsDoc
	this_.priv_CheckJSDocStyle_checkForAtSignCodePresence_6570(t, cls, jsDoc)
	if jsDoc == nil {
		return
	}
	if jsDoc.IsConstructor() {
		t.Report_3(cls, StaticGet_CheckJSDocStyle_CLASS_DISALLOWED_JSDOC(), []string{})
	}
}

// Java method: CheckJSDocStyle.checkMissingJsDoc(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node function)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_checkMissingJsDoc_6575(t *NodeTraversal, function pkg_rhino.Node) {
	_ = t
	_ = function
	if this_.priv_CheckJSDocStyle_isFunctionThatShouldHaveJsDoc_6576(t, function) && !(this_.priv_CheckJSDocStyle_isTestMethod_6577(function)) {
		t.Report_3(function, StaticGet_CheckJSDocStyle_MISSING_JSDOC(), []string{})
	}
}

// Java method: CheckJSDocStyle.isFunctionThatShouldHaveJsDoc(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node function)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_isFunctionThatShouldHaveJsDoc_6576(t *NodeTraversal, function pkg_rhino.Node) bool {
	_ = t
	_ = function
	if !(t.InGlobalHoistScope() || t.InModuleScope()) {
		return false
	}
	if NodeUtil_IsFunctionDeclaration(function) {
		return true
	}
	if NodeUtil_IsNameDeclaration(function.GetGrandparent()) || function.GetParent().IsAssign() {
		return true
	}
	if function.GetParent().IsExport() {
		return true
	}
	if function.GetGrandparent().IsClassMembers() {
		var memberNode pkg_rhino.Node = function.GetParent()
		_ = memberNode
		if memberNode.IsMemberFunctionDef() {
			return !(this_.priv_CheckJSDocStyle_isConstructorWithoutParameters_6578(function))
		} else if memberNode.IsGetterDef() || memberNode.IsSetterDef() {
			return true
		}
	}
	return (function.GetGrandparent().IsObjectLit() && NodeUtil_IsCallTo_2_Node_string(function.GetGrandparent().GetParent(), "Polymer"))
}

// Java method: CheckJSDocStyle.isTestMethod(com.google.javascript.rhino.Node function)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_isTestMethod_6577(function pkg_rhino.Node) bool {
	_ = function
	var bestLValue pkg_rhino.Node = NodeUtil_GetBestLValue(function)
	_ = bestLValue
	var name string
	_ = name
	if bestLValue != 0 {
		name = NodeUtil_GetBestLValueName(bestLValue)
	} else {
		name = j2gort.NullString
	}
	return ((name != j2gort.NullString) && ExportTestFunctions_IsTestFunction(name))
}

// Java method: CheckJSDocStyle.isConstructorWithoutParameters(com.google.javascript.rhino.Node function)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_isConstructorWithoutParameters_6578(function pkg_rhino.Node) bool {
	_ = function
	return (NodeUtil_IsEs6Constructor(function) && !(NodeUtil_GetFunctionParameters(function).HasChildren()))
}

// Java method: CheckJSDocStyle.checkParams(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node function, com.google.javascript.rhino.JSDocInfo jsDoc)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_checkParams_6579(t *NodeTraversal, function pkg_rhino.Node, jsDoc *pkg_rhino.JSDocInfo) {
	_ = t
	_ = function
	_ = jsDoc
	if (jsDoc != nil) && !(jsDoc.GetType()).IsJavaNil_() {
		return
	}
	var paramCountFromJsDoc int32
	_ = paramCountFromJsDoc
	if jsDoc == nil {
		paramCountFromJsDoc = int32(0)
	} else {
		paramCountFromJsDoc = jsDoc.GetParameterCount()
	}
	if paramCountFromJsDoc == int32(0) {
		this_.priv_CheckJSDocStyle_checkInlineParams_6580(t, function, jsDoc)
	} else {
		var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(function)
		_ = paramList
		if !(paramList.HasXChildren(paramCountFromJsDoc)) {
			this_.compiler.Report_1(JSError_Make_3(paramList, StaticGet_CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS(), []string{func() string {
				if jsDoc.IsOverride() {
					return ""
				}
				return ""
			}()}))
			return
		}
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		{
			var i int32 = int32(0)
			_ = i
			for ; i < paramCountFromJsDoc; i = i + 1 {
				if param.GetJSDocInfo() != nil {
					t.Report_3(param, StaticGet_CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES(), []string{})
				}
				var name string = jsDoc.GetParameterNameAt(i)
				_ = name
				var paramType pkg_rhino.JSTypeExpression = jsDoc.GetParameterTypeAt(i)
				_ = paramType
				if this_.priv_CheckJSDocStyle_checkParam_6581(t, param, name, paramType) {
					return
				}
				param = param.GetNext()
			}
		}
	}
}

// Java method: CheckJSDocStyle.checkInlineParams(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node function, com.google.javascript.rhino.JSDocInfo fnJSDoc)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_checkInlineParams_6580(t *NodeTraversal, function pkg_rhino.Node, fnJSDoc *pkg_rhino.JSDocInfo) {
	_ = t
	_ = function
	_ = fnJSDoc
	var paramList pkg_rhino.Node = NodeUtil_GetFunctionParameters(function)
	_ = paramList
	{
		var param pkg_rhino.Node = paramList.GetFirstChild()
		_ = param
		for ; param != 0; param = param.GetNext() {
			var jsDoc *pkg_rhino.JSDocInfo
			_ = jsDoc
			if param.IsDefaultValue() {
				jsDoc = param.GetFirstChild().GetJSDocInfo()
			} else {
				jsDoc = param.GetJSDocInfo()
			}
			if jsDoc == nil {
				this_.compiler.Report_1(JSError_Make_3(param, StaticGet_CheckJSDocStyle_MISSING_PARAMETER_JSDOC(), []string{func() string {
					if (fnJSDoc != nil) && fnJSDoc.IsOverride() {
						return ""
					}
					return ""
				}()}))
				return
			} else {
				var paramType pkg_rhino.JSTypeExpression = jsDoc.GetType()
				_ = paramType
				j2gort.CheckNotNullVal(paramType)
				this_.priv_CheckJSDocStyle_checkParam_6581(t, param, j2gort.NullString, paramType)
			}
		}
	}
}

// Java method: CheckJSDocStyle.checkParam(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node param, java.lang.@org.jspecify.annotations.Nullable String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSTypeExpression paramType)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_checkParam_6581(t *NodeTraversal, param pkg_rhino.Node, name string, paramType pkg_rhino.JSTypeExpression) bool {
	_ = t
	_ = param
	_ = name
	_ = paramType
	var nameOptional bool = false
	_ = nameOptional
	var nodeToCheck pkg_rhino.Node = param
	_ = nodeToCheck
	if param.IsDefaultValue() {
		nodeToCheck = param.GetFirstChild()
		nameOptional = true
	} else if param.IsName() {
		nameOptional = javalang.StringStartsWith(j2gort.CheckNotNullString(param.GetString()), "opt_")
	} else {
		if !(param.IsDestructuringPattern() || param.IsRest()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(param))
		}
		nameOptional = false
	}
	if (name == j2gort.NullString) || !(nodeToCheck.IsName()) {
		name = "<unknown name>"
	} else if !(nodeToCheck.MatchesQualifiedName_1_string(name)) {
		t.Report_3(nodeToCheck, StaticGet_CheckJSDocStyle_INCORRECT_PARAM_NAME(), []string{})
		return true
	}
	if !(nameOptional) {
		return false
	}
	var jsDocOptional bool = (!(paramType).IsJavaNil_() && paramType.IsOptionalArg())
	_ = jsDocOptional
	if jsDocOptional {
		return false
	}
	var errorSource pkg_rhino.Node
	_ = errorSource
	if !(paramType).IsJavaNil_() {
		errorSource = paramType.GetRoot()
	} else {
		errorSource = nodeToCheck
	}
	t.Report_3(errorSource, StaticGet_CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL(), []string{name})
	return true
}

// Java static method: CheckJSDocStyle.isDefaultAssignedParamWithInlineJsDoc(com.google.javascript.rhino.Node param)
func CheckJSDocStyle_IsDefaultAssignedParamWithInlineJsDoc(param pkg_rhino.Node) bool {
	clinit_CheckJSDocStyle()
	_ = param
	if param.IsDefaultValue() {
		if param.HasChildren() && param.GetFirstChild().IsName() {
			return (param.GetFirstChild().GetJSDocInfo() != nil)
		}
	}
	return false
}

// Java method: CheckJSDocStyle.hasAnyInlineJsDoc(com.google.javascript.rhino.Node function)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_hasAnyInlineJsDoc_6582(function pkg_rhino.Node) bool {
	_ = function
	if function.GetFirstChild().GetJSDocInfo() != nil {
		return true
	}
	{
		var param pkg_rhino.Node = NodeUtil_GetFunctionParameters(function).GetFirstChild()
		_ = param
		for ; param != 0; param = param.GetNext() {
			if (param.GetJSDocInfo() != nil) || CheckJSDocStyle_IsDefaultAssignedParamWithInlineJsDoc(param) {
				return true
			}
		}
	}
	return false
}

// Java method: CheckJSDocStyle.checkReturn(com.google.javascript.rhino.Node function, com.google.javascript.rhino.JSDocInfo jsDoc)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_checkReturn_6583(function pkg_rhino.Node, jsDoc *pkg_rhino.JSDocInfo) {
	_ = function
	_ = jsDoc
	if (jsDoc != nil) && ((jsDoc.HasType() || jsDoc.IsConstructor()) || jsDoc.HasReturnType()) {
		return
	}
	if NodeUtil_IsEs6Constructor(function) {
		return
	}
	if function.GetFirstChild().GetJSDocInfo() != nil {
		return
	}
	var finder *CheckJSDocStyle_FindNonTrivialReturn = NewCheckJSDocStyle_FindNonTrivialReturn()
	_ = finder
	NodeTraversal_Traverse(this_.compiler, function.GetLastChild(), finder)
	if finder.found {
		this_.compiler.Report_1(JSError_Make_3(function, StaticGet_CheckJSDocStyle_MISSING_RETURN_JSDOC(), []string{func() string {
			if (jsDoc != nil) && jsDoc.IsOverride() {
				return ""
			}
			return ""
		}()}))
	}
}

// Java method: CheckJSDocStyle.checkLicenseComment(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *CheckJSDocStyle) priv_CheckJSDocStyle_checkLicenseComment_6584(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if (n.GetJSDocInfo() == nil) || (n.GetJSDocInfo().GetLicense() == j2gort.NullString) {
		return
	}
	var license string = n.GetJSDocInfo().GetLicense()
	_ = license
	if j2gort.JStrIndexOf(license, "@externs") >= 0 {
		t.Report_3(n, StaticGet_CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS(), []string{})
	}
}

// Java class: com.google.javascript.jscomp.lint.CheckJSDocStyle.FindNonTrivialReturn
type CheckJSDocStyle_FindNonTrivialReturn struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	found bool `gson:"found"`
}

func (this_ *CheckJSDocStyle_FindNonTrivialReturn) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckJSDocStyle_FindNonTrivialReturn) IsJavaNil_() bool { return this_ == nil }
func (*CheckJSDocStyle_FindNonTrivialReturn) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckJSDocStyle$FindNonTrivialReturn", "FindNonTrivialReturn"
}
func (this_ *CheckJSDocStyle_FindNonTrivialReturn) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*CheckJSDocStyle_FindNonTrivialReturn) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckJSDocStyle.FindNonTrivialReturn", "com.google.javascript.jscomp.lint.CheckJSDocStyle$FindNonTrivialReturn", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckJSDocStyle_FindNonTrivialReturn](v_)
	})
}
func (*CheckJSDocStyle_FindNonTrivialReturn) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*CheckJSDocStyle_FindNonTrivialReturn) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *CheckJSDocStyle_FindNonTrivialReturn) initFields_CheckJSDocStyle_FindNonTrivialReturn() {
	_ = this_
}

// Java constructor: FindNonTrivialReturn()
func NewCheckJSDocStyle_FindNonTrivialReturn() *CheckJSDocStyle_FindNonTrivialReturn {
	this_ := &CheckJSDocStyle_FindNonTrivialReturn{}
	InitCheckJSDocStyle_FindNonTrivialReturn(this_)
	return this_
}

func InitCheckJSDocStyle_FindNonTrivialReturn(this_ *CheckJSDocStyle_FindNonTrivialReturn) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_CheckJSDocStyle_FindNonTrivialReturn()
}

// Java method: FindNonTrivialReturn.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckJSDocStyle_FindNonTrivialReturn) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if this_.found {
		return false
	}
	if NodeUtil_IsShallowStatementTree(parent) {
		if n.IsReturn() && n.HasChildren() {
			this_.found = true
			return false
		}
		return true
	}
	return false
}

// Java class: com.google.javascript.jscomp.lint.CheckJSDocStyle.ExternsCallback
type CheckJSDocStyle_ExternsCallback struct {
	idHash_ atomic.Int32
}

func (this_ *CheckJSDocStyle_ExternsCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckJSDocStyle_ExternsCallback) IsJavaNil_() bool { return this_ == nil }
func (*CheckJSDocStyle_ExternsCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckJSDocStyle$ExternsCallback", "ExternsCallback"
}
func (this_ *CheckJSDocStyle_ExternsCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckJSDocStyle_ExternsCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckJSDocStyle.ExternsCallback", "com.google.javascript.jscomp.lint.CheckJSDocStyle$ExternsCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckJSDocStyle_ExternsCallback](v_)
	})
}
func (*CheckJSDocStyle_ExternsCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *CheckJSDocStyle_ExternsCallback) initFields_CheckJSDocStyle_ExternsCallback() {
	_ = this_
}

// Java constructor: ExternsCallback()
func NewCheckJSDocStyle_ExternsCallback() *CheckJSDocStyle_ExternsCallback {
	this_ := &CheckJSDocStyle_ExternsCallback{}
	InitCheckJSDocStyle_ExternsCallback(this_)
	return this_
}

func InitCheckJSDocStyle_ExternsCallback(this_ *CheckJSDocStyle_ExternsCallback) {
	this_.initFields_CheckJSDocStyle_ExternsCallback()
}

// Java method: ExternsCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckJSDocStyle_ExternsCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return ((parent == 0) || n.IsScript())
}

// Java method: ExternsCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckJSDocStyle_ExternsCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = info
		if (info == nil) || !(info.IsExterns()) {
			t.Report_3(n, StaticGet_CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED(), []string{})
		}
	}
}
func (*CheckJSDocStyle) AlwaysShouldTraverse_()                           {}
func (*CheckJSDocStyle) CanReleaseSyntacticScopes_()                      {}
func (*CheckJSDocStyle_FindNonTrivialReturn) NoOpVisit_()                 {}
func (*CheckJSDocStyle_FindNonTrivialReturn) CanReleaseSyntacticScopes_() {}
func (*CheckJSDocStyle_ExternsCallback) CanReleaseSyntacticScopes_()      {}
