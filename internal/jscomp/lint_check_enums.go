// Port of com/google/javascript/jscomp/lint/CheckEnums.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var checkEnumsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_OBJECTLIT)
	return s
}()

var _ nodeTraversal_visitTokens = (*CheckEnums)(nil)

func (*CheckEnums) visitTokens() *ast.TokenSet { return &checkEnumsVisitTokens }

// CheckEnums checks the following:
//
//  1. Whether there are duplicate values in enums.
//  2. Whether enum type (and its initializer values) are string/number or not.
//  3. Whether string enum values are statically initialized or not.
type CheckEnums struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

var CheckEnums_DUPLICATE_ENUM_VALUE = DiagnosticType_Disabled(
	"JSC_DUPLICATE_ENUM_VALUE", "The value {0} is duplicated in this enum.")

var CheckEnums_COMPUTED_PROP_NAME_IN_ENUM = DiagnosticType_Disabled(
	"JSC_COMPUTED_PROP_NAME_IN_ENUM", "Computed property name used in enum.")

var CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM = DiagnosticType_Disabled(
	"JSC_SHORTHAND_ASSIGNMENT_IN_ENUM", "Shorthand assignment used in enum.")

var CheckEnums_ENUM_PROP_NOT_CONSTANT = DiagnosticType_Disabled("JSC_ENUM_PROP_NOT_CONSTANT", "enum key {0} must be in ALL_CAPS.")

var CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER = DiagnosticType_Disabled(
	"JSC_ENUM_VALUE_NOT_STRING_OR_NUMBER",
	"enum type must be either string or number.")

var CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM = DiagnosticType_Disabled(
	"JSC_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM",
	"Enum string values must be statically initialized as per the style guide.")

// NewCheckEnums is new CheckEnums(compiler).
//
// Java: CheckEnums(AbstractCompiler)
func NewCheckEnums(compiler AbstractCompiler) *CheckEnums {
	return &CheckEnums{compiler: compiler}
}

// Java: CheckEnums.process(Node, Node)
func (c *CheckEnums) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckEnums.visit(NodeTraversal, Node, Node)
func (c *CheckEnums) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsObjectLit() {
		jsdoc := NodeUtil_GetBestJSDocInfo(n)
		if jsdoc != nil && jsdoc.HasEnumParameterType() {
			c.CheckNamingAndAssignmentUsage(t, n)
			CheckEnums_CheckDuplicateEnumValues(t, n)
			CheckEnums_CheckEnumTypeAndInitializerValues(t, n, jsdoc)
		}
	}
}

// Java: CheckEnums.checkEnumTypeAndInitializerValues(NodeTraversal, Node, JSDocInfo)
func CheckEnums_CheckEnumTypeAndInitializerValues(
	t *NodeTraversal, n rhino.Node, jsDocInfo *rhino.JSDocInfo) {
	jutil.CheckArgumentMsg1(n.IsObjectLit(), "%s", n)
	enumTypeExpr := jsDocInfo.GetEnumParameterType()

	enumType := enumTypeExpr.GetRoot()
	isStringEnum := enumType.IsStringLit() && enumType.GetString() == "string"
	isNumberEnum := enumType.IsStringLit() && enumType.GetString() == "number"
	if !isStringEnum && !isNumberEnum {
		// warn on `@enum {?}`, `@enum {boolean}`, `@enum {Some|Another}`, `@enum {SomeName}` etc`
		t.Report(n, CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	}
	if isStringEnum {
		CheckEnums_CheckStringEnumInitializerValues(t, n)
	}
}

// CheckEnums_CheckStringEnumInitializerValues reports a warning if the
// string enum value is not statically initialized
//
// Java: CheckEnums.checkStringEnumInitializerValues(NodeTraversal, Node)
func CheckEnums_CheckStringEnumInitializerValues(t *NodeTraversal, enumNode rhino.Node) {
	jutil.CheckArgumentMsg1(enumNode.IsObjectLit(), "%s", enumNode)
	for prop := enumNode.GetFirstChild(); !prop.IsNil(); prop = prop.GetNext() {
		// valueNode is guaranteed to exist by this time, as shorthand `{A}`s are converted to `{A:A}`
		valueNode := prop.GetLastChild()
		if !valueNode.IsStringLit() && !(valueNode.IsTemplateLit() && valueNode.HasOneChild()) {
			// neither string nor substitution-free template literal; report finding.
			t.Report(valueNode, CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
		}
	}
}

// Java: CheckEnums.checkNamingAndAssignmentUsage(NodeTraversal, Node)
func (c *CheckEnums) CheckNamingAndAssignmentUsage(t *NodeTraversal, objLit rhino.Node) {
	for child := objLit.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		c.CheckName(t, child)
	}
}

// Java: CheckEnums.checkName(NodeTraversal, Node)
func (c *CheckEnums) CheckName(t *NodeTraversal, prop rhino.Node) {
	if prop.IsComputedProp() {
		t.Report(prop, CheckEnums_COMPUTED_PROP_NAME_IN_ENUM)
		return
	}

	if prop.IsStringKey() && prop.IsShorthandProperty() {
		t.Report(prop, CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM)
	}

	if !c.compiler.GetCodingConvention().IsValidEnumKey(prop.GetString()) {
		t.Report(prop, CheckEnums_ENUM_PROP_NOT_CONSTANT)
	}
}

// Java: CheckEnums.checkDuplicateEnumValues(NodeTraversal, Node)
func CheckEnums_CheckDuplicateEnumValues(t *NodeTraversal, enumNode rhino.Node) {
	var values map[string]bool
	for prop := enumNode.GetFirstChild(); !prop.IsNil(); prop = prop.GetNext() {
		valueNode := prop.GetLastChild()
		var value string
		if valueNode.IsNil() {
			return
		} else if valueNode.IsStringLit() {
			value = valueNode.GetString()
		} else if valueNode.IsNumber() {
			value = jutil.DoubleToString(valueNode.GetDouble())
		} else {
			return
		}

		if values[value] {
			t.Report(valueNode, CheckEnums_DUPLICATE_ENUM_VALUE, value)
		} else {
			if values == nil {
				values = make(map[string]bool)
			}
			values[value] = true
		}
	}
}
