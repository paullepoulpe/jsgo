// Port of com/google/javascript/jscomp/lint/CheckEnums.java @ closure-compiler v20261104

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

var checkEnumsVisitTokens = makeVisitTokenSet(rhino.Token_OBJECTLIT)

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

var CheckEnums_ENUM_PROP_NOT_CONSTANT = DiagnosticType_Disabled(
	"JSC_ENUM_PROP_NOT_CONSTANT", "enum key {0} must be in ALL_CAPS.")

var CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER = DiagnosticType_Disabled(
	"JSC_ENUM_VALUE_NOT_STRING_OR_NUMBER",
	"enum type must be either string or number.")

var CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM = DiagnosticType_Disabled(
	"JSC_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM",
	"Enum string values must be statically initialized as per the style guide.")

// NewCheckEnums constructs a CheckEnums pass.
//
// Java: CheckEnums(AbstractCompiler)
func NewCheckEnums(compiler AbstractCompiler) *CheckEnums {
	return &CheckEnums{compiler: compiler}
}

func (*CheckEnums) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckEnums(c)
}

// Java: CheckEnums.process(Node, Node)
func (c *CheckEnums) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckEnums.visit(NodeTraversal, Node, Node)
func (c *CheckEnums) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	objLit, ok := ast.AsObjectLit(n)
	if !ok {
		return
	}
	if jsdoc := NodeUtil_GetBestJSDocInfo(n); jsdoc != nil && jsdoc.HasEnumParameterType() {
		c.checkNamingAndAssignmentUsage(t, objLit)
		checkDuplicateEnumValues(t, objLit)
		checkEnumTypeAndInitializerValues(t, objLit, jsdoc)
	}
}

func checkEnumTypeAndInitializerValues(t *NodeTraversal, enumLit ast.ObjectLit, jsDocInfo *rhino.JSDocInfo) {
	enumType := jsDocInfo.GetEnumParameterType().GetRoot()
	typeStr, isStr := ast.AsStringLit(enumType)
	isStringEnum := isStr && typeStr.Value() == "string"
	isNumberEnum := isStr && typeStr.Value() == "number"
	if !isStringEnum && !isNumberEnum {
		// Warn on `@enum {?}`, `@enum {boolean}`, `@enum {Some|Another}`, `@enum {SomeName}`, etc.
		t.Report(enumLit.Node(), CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	}
	if isStringEnum {
		checkStringEnumInitializerValues(t, enumLit)
	}
}

// checkStringEnumInitializerValues reports a warning if any string enum value
// is not statically initialized.
func checkStringEnumInitializerValues(t *NodeTraversal, enumLit ast.ObjectLit) {
	for member := range enumLit.Members() {
		// valueNode is guaranteed to exist by this time, as shorthand `{A}`s are converted to `{A:A}`.
		valueNode, _ := enumMemberValue(member)
		if !isStaticStringValue(valueNode) {
			// Neither string nor substitution-free template literal; report finding.
			t.Report(valueNode, CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
		}
	}
}

func isStaticStringValue(n rhino.Node) bool {
	if _, ok := ast.AsStringLit(n); ok {
		return true
	}
	if tmpl, ok := ast.AsTemplateLit(n); ok {
		_, hasOnlyPart := tmpl.OnlyPart()
		return hasOnlyPart
	}
	return false
}

func (c *CheckEnums) checkNamingAndAssignmentUsage(t *NodeTraversal, objLit ast.ObjectLit) {
	for member := range objLit.Members() {
		c.checkMemberName(t, member)
	}
}

func (c *CheckEnums) checkMemberName(t *NodeTraversal, member ast.ObjectMember) {
	prop := member.Node()
	if _, ok := member.AsComputedProp(); ok {
		t.Report(prop, CheckEnums_COMPUTED_PROP_NAME_IN_ENUM)
		return
	}

	if _, ok := member.AsStringKey(); ok && prop.IsShorthandProperty() {
		t.Report(prop, CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM)
	}

	pkm, ok := member.AsPlainKeyMember()
	if !ok || !c.compiler.GetCodingConvention().IsValidEnumKey(pkm.Key()) {
		t.Report(prop, CheckEnums_ENUM_PROP_NOT_CONSTANT)
	}
}

func checkDuplicateEnumValues(t *NodeTraversal, enumLit ast.ObjectLit) {
	var values map[string]bool
	for member := range enumLit.Members() {
		valueNode, ok := enumMemberValue(member)
		if !ok {
			return
		}
		var value string
		if str, ok := ast.AsStringLit(valueNode); ok {
			value = str.Value()
		} else if num, ok := ast.AsNumber(valueNode); ok {
			value = jutil.DoubleToString(num.Value())
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

func enumMemberValue(member ast.ObjectMember) (rhino.Node, bool) {
	if pkm, ok := member.AsPlainKeyMember(); ok {
		return pkm.Value().Node(), true
	}
	if cp, ok := member.AsComputedProp(); ok {
		val, hasVal := cp.Value()
		return val.Node(), hasVal
	}
	if spread, ok := member.AsObjectSpread(); ok {
		return spread.Expr().Node(), true
	}
	return rhino.NilNode, false
}
