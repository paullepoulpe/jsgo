/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckEnums.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.lint.CheckEnums
type CheckEnums struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckEnums) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckEnums) IsJavaNil_() bool { return this_ == nil }
func (*CheckEnums) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckEnums", "CheckEnums"
}
func (this_ *CheckEnums) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckEnums) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckEnums", "com.google.javascript.jscomp.lint.CheckEnums", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckEnums](v_)
	})
}
func (*CheckEnums) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckEnums) Implements_jscomp_CompilerPass()                            {}
func (*CheckEnums) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckEnums_DUPLICATE_ENUM_VALUE *DiagnosticType

func StaticGet_CheckEnums_DUPLICATE_ENUM_VALUE() *DiagnosticType {
	clinit_CheckEnums()
	return static_CheckEnums_DUPLICATE_ENUM_VALUE
}
func StaticSet_CheckEnums_DUPLICATE_ENUM_VALUE(v *DiagnosticType) {
	clinit_CheckEnums()
	static_CheckEnums_DUPLICATE_ENUM_VALUE = v
}

var static_CheckEnums_COMPUTED_PROP_NAME_IN_ENUM *DiagnosticType

func StaticGet_CheckEnums_COMPUTED_PROP_NAME_IN_ENUM() *DiagnosticType {
	clinit_CheckEnums()
	return static_CheckEnums_COMPUTED_PROP_NAME_IN_ENUM
}
func StaticSet_CheckEnums_COMPUTED_PROP_NAME_IN_ENUM(v *DiagnosticType) {
	clinit_CheckEnums()
	static_CheckEnums_COMPUTED_PROP_NAME_IN_ENUM = v
}

var static_CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM *DiagnosticType

func StaticGet_CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM() *DiagnosticType {
	clinit_CheckEnums()
	return static_CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM
}
func StaticSet_CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM(v *DiagnosticType) {
	clinit_CheckEnums()
	static_CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM = v
}

var static_CheckEnums_ENUM_PROP_NOT_CONSTANT *DiagnosticType

func StaticGet_CheckEnums_ENUM_PROP_NOT_CONSTANT() *DiagnosticType {
	clinit_CheckEnums()
	return static_CheckEnums_ENUM_PROP_NOT_CONSTANT
}
func StaticSet_CheckEnums_ENUM_PROP_NOT_CONSTANT(v *DiagnosticType) {
	clinit_CheckEnums()
	static_CheckEnums_ENUM_PROP_NOT_CONSTANT = v
}

var static_CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER *DiagnosticType

func StaticGet_CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER() *DiagnosticType {
	clinit_CheckEnums()
	return static_CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER
}
func StaticSet_CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER(v *DiagnosticType) {
	clinit_CheckEnums()
	static_CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER = v
}

var static_CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM *DiagnosticType

func StaticGet_CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM() *DiagnosticType {
	clinit_CheckEnums()
	return static_CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM
}
func StaticSet_CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM(v *DiagnosticType) {
	clinit_CheckEnums()
	static_CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM = v
}

var clinitDone_CheckEnums uint32
var clinitState_CheckEnums j2gort.ClinitState

func clinit_CheckEnums() {
	if atomic.LoadUint32(&clinitDone_CheckEnums) != 0 {
		return
	}
	clinitSlow_CheckEnums()
}

func clinitSlow_CheckEnums() {
	if !j2gort.EnterClinit(&clinitState_CheckEnums, &clinitDone_CheckEnums) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckEnums, &clinitDone_CheckEnums)
	static_CheckEnums_DUPLICATE_ENUM_VALUE = DiagnosticType_Disabled("JSC_DUPLICATE_ENUM_VALUE", "The value {0} is duplicated in this enum.")
	static_CheckEnums_COMPUTED_PROP_NAME_IN_ENUM = DiagnosticType_Disabled("JSC_COMPUTED_PROP_NAME_IN_ENUM", "Computed property name used in enum.")
	static_CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM = DiagnosticType_Disabled("JSC_SHORTHAND_ASSIGNMENT_IN_ENUM", "Shorthand assignment used in enum.")
	static_CheckEnums_ENUM_PROP_NOT_CONSTANT = DiagnosticType_Disabled("JSC_ENUM_PROP_NOT_CONSTANT", "enum key {0} must be in ALL_CAPS.")
	static_CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER = DiagnosticType_Disabled("JSC_ENUM_VALUE_NOT_STRING_OR_NUMBER", "enum type must be either string or number.")
	static_CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM = DiagnosticType_Disabled("JSC_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM", "Enum string values must be statically initialized as per the style guide.")
}

func (this_ *CheckEnums) initFields_CheckEnums() {
	_ = this_
}

// Java constructor: CheckEnums(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckEnums(compiler AbstractCompiler) *CheckEnums {
	clinit_CheckEnums()
	this_ := &CheckEnums{}
	InitCheckEnums(this_, compiler)
	return this_
}

func InitCheckEnums(this_ *CheckEnums, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckEnums()
	this_.compiler = compiler
}

// Java method: CheckEnums.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckEnums) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckEnums.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckEnums) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsObjectLit() {
		var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = jsdoc
		if (jsdoc != nil) && jsdoc.HasEnumParameterType() {
			this_.priv_CheckEnums_checkNamingAndAssignmentUsage_6515(t, n)
			CheckEnums_CheckDuplicateEnumValues(t, n)
			CheckEnums_CheckEnumTypeAndInitializerValues(t, n, jsdoc)
		}
	}
}

// Java static method: CheckEnums.checkEnumTypeAndInitializerValues(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSDocInfo jsDocInfo)
func CheckEnums_CheckEnumTypeAndInitializerValues(t *NodeTraversal, n pkg_rhino.Node, jsDocInfo *pkg_rhino.JSDocInfo) {
	clinit_CheckEnums()
	_ = t
	_ = n
	_ = jsDocInfo
	if !(n.IsObjectLit()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	var enumTypeExpr pkg_rhino.JSTypeExpression = jsDocInfo.GetEnumParameterType()
	_ = enumTypeExpr
	var enumType pkg_rhino.Node = enumTypeExpr.GetRoot()
	_ = enumType
	var isStringEnum bool = (enumType.IsStringLit() && (j2gort.CheckNotNullString(enumType.GetString()) == "string"))
	_ = isStringEnum
	var isNumberEnum bool = (enumType.IsStringLit() && (j2gort.CheckNotNullString(enumType.GetString()) == "number"))
	_ = isNumberEnum
	if !(isStringEnum) && !(isNumberEnum) {
		t.Report_3(n, StaticGet_CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER(), []string{})
	}
	if isStringEnum {
		CheckEnums_CheckStringEnumInitializerValues(t, n)
	}
}

// Java static method: CheckEnums.checkStringEnumInitializerValues(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node enumNode)
func CheckEnums_CheckStringEnumInitializerValues(t *NodeTraversal, enumNode pkg_rhino.Node) {
	clinit_CheckEnums()
	_ = t
	_ = enumNode
	if !(enumNode.IsObjectLit()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(enumNode))
	}
	{
		var prop pkg_rhino.Node = enumNode.GetFirstChild()
		_ = prop
		for ; prop != 0; prop = prop.GetNext() {
			var valueNode pkg_rhino.Node = prop.GetLastChild()
			_ = valueNode
			if !(valueNode.IsStringLit()) && !(valueNode.IsTemplateLit() && valueNode.HasOneChild()) {
				t.Report_3(valueNode, StaticGet_CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM(), []string{})
			}
		}
	}
}

// Java method: CheckEnums.checkNamingAndAssignmentUsage(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node objLit)
func (this_ *CheckEnums) priv_CheckEnums_checkNamingAndAssignmentUsage_6515(t *NodeTraversal, objLit pkg_rhino.Node) {
	_ = t
	_ = objLit
	{
		var child pkg_rhino.Node = objLit.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_CheckEnums_checkName_6516(t, child)
		}
	}
}

// Java method: CheckEnums.checkName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node prop)
func (this_ *CheckEnums) priv_CheckEnums_checkName_6516(t *NodeTraversal, prop pkg_rhino.Node) {
	_ = t
	_ = prop
	if prop.IsComputedProp() {
		t.Report_3(prop, StaticGet_CheckEnums_COMPUTED_PROP_NAME_IN_ENUM(), []string{})
		return
	}
	if prop.IsStringKey() && prop.IsShorthandProperty() {
		t.Report_3(prop, StaticGet_CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM(), []string{})
	}
	if !(this_.compiler.GetCodingConvention().IsValidEnumKey(prop.GetString())) {
		t.Report_3(prop, StaticGet_CheckEnums_ENUM_PROP_NOT_CONSTANT(), []string{})
	}
}

// Java static method: CheckEnums.checkDuplicateEnumValues(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node enumNode)
func CheckEnums_CheckDuplicateEnumValues(t *NodeTraversal, enumNode pkg_rhino.Node) {
	clinit_CheckEnums()
	_ = t
	_ = enumNode
	var values j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = values
	{
		var prop pkg_rhino.Node = enumNode.GetFirstChild()
		_ = prop
		for ; prop != 0; prop = prop.GetNext() {
			var valueNode pkg_rhino.Node = prop.GetLastChild()
			_ = valueNode
			var value string = j2gort.NullString
			_ = value
			if valueNode == 0 {
				return
			} else if valueNode.IsStringLit() {
				value = valueNode.GetString()
			} else if valueNode.IsNumber() {
				value = javalang.DoubleToString(valueNode.GetDouble())
			} else {
				return
			}
			if !(values.Add(value)) {
				t.Report_3(valueNode, StaticGet_CheckEnums_DUPLICATE_ENUM_VALUE(), []string{value})
			}
		}
	}
}
func (*CheckEnums) AlwaysShouldTraverse_()      {}
func (*CheckEnums) CanReleaseSyntacticScopes_() {}
