// Port of com/google/javascript/jscomp/lint/CheckJSDocStyle.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckJSDocStyle checks for various JSDoc-related style issues, such as
// function definitions without JsDoc, params with no corresponding @param
// annotation, coding conventions not being respected, etc.
type CheckJSDocStyle struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

// Java: CheckJSDocStyle.CLASS_DISALLOWED_JSDOC
var CheckJSDocStyle_CLASS_DISALLOWED_JSDOC = DiagnosticType_Disabled("JSC_CLASS_DISALLOWED_JSDOC",
	"@constructor annotations are redundant on classes.")

// Java: CheckJSDocStyle.MISSING_JSDOC
var CheckJSDocStyle_MISSING_JSDOC = DiagnosticType_Disabled("JSC_MISSING_JSDOC", "Function must have JSDoc.")

// Java: CheckJSDocStyle.INCORRECT_ANNOTATION_ON_GETTER_SETTER
var CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER = DiagnosticType_Disabled(
	"JSC_TYPE_ON_GETTER_SETTER",
	"Getters and setters must not have @type annotations. Did you mean @return or @param"+
		" instead?")

// Java: CheckJSDocStyle.MISSING_PARAMETER_JSDOC
var CheckJSDocStyle_MISSING_PARAMETER_JSDOC = DiagnosticType_Disabled("JSC_MISSING_PARAMETER_JSDOC", "Parameter must have JSDoc.{0}")

// Java: CheckJSDocStyle.MIXED_PARAM_JSDOC_STYLES
var CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES = DiagnosticType_Disabled("JSC_MIXED_PARAM_JSDOC_STYLES",
	"Functions may not use both @param annotations and inline JSDoc")

// Java: CheckJSDocStyle.MISSING_RETURN_JSDOC
var CheckJSDocStyle_MISSING_RETURN_JSDOC = DiagnosticType_Disabled(
	"JSC_MISSING_RETURN_JSDOC",
	"Function that returns a value must have JSDoc indicating the return type.{0}")

// Java: CheckJSDocStyle.OPTIONAL_PARAM_NOT_MARKED_OPTIONAL
var CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL = DiagnosticType_Disabled(
	"JSC_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL",
	"Parameter {0} is optional so it must have a JSDoc type ending with ''=''")

// Java: CheckJSDocStyle.WRONG_NUMBER_OF_PARAMS
var CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS = DiagnosticType_Disabled("JSC_WRONG_NUMBER_OF_PARAMS", "Wrong number of @param annotations")

// Java: CheckJSDocStyle.INCORRECT_PARAM_NAME
var CheckJSDocStyle_INCORRECT_PARAM_NAME = DiagnosticType_Disabled("JSC_INCORRECT_PARAM_NAME",
	"Incorrect param name. Are your @param annotations in the wrong order?")

// Java: CheckJSDocStyle.EXTERNS_FILES_SHOULD_BE_ANNOTATED
var CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED = DiagnosticType_Disabled("JSC_EXTERNS_FILES_SHOULD_BE_ANNOTATED",
	"Externs files should be annotated with @externs in the @fileoverview block.")

// Java: CheckJSDocStyle.LICENSE_CONTAINS_AT_EXTERNS
var CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS = DiagnosticType_Disabled(
	"JSC_LICENSE_CONTAINS_AT_EXTERNS",
	"@license block contains an @externs annotation, which will be parsed as plain "+
		"license text instead of an actual @externs annotation. You probably meant to put "+
		"@externs in a separate @fileoverview block.")

// Java: CheckJSDocStyle.PREFER_BACKTICKS_TO_AT_SIGN_CODE
var CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE = DiagnosticType_Disabled(
	"JSC_PREFER_BACKTICKS_TO_AT_SIGN_CODE",
	"Use `some_code` instead of '{'@code some_code'}'.")

// Java: CheckJSDocStyle.LINT_DIAGNOSTICS
var CheckJSDocStyle_LINT_DIAGNOSTICS = NewDiagnosticGroup(
	CheckJSDocStyle_CLASS_DISALLOWED_JSDOC,
	CheckJSDocStyle_MISSING_JSDOC,
	CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER,
	CheckJSDocStyle_MISSING_PARAMETER_JSDOC,
	CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES,
	CheckJSDocStyle_MISSING_RETURN_JSDOC,
	CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL,
	CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS,
	CheckJSDocStyle_INCORRECT_PARAM_NAME,
	CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED,
	CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS,
	CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE)

// Java: CheckJSDocStyle.ALL_DIAGNOSTICS
var CheckJSDocStyle_ALL_DIAGNOSTICS = NewDiagnosticGroupFromGroups(CheckJSDocStyle_LINT_DIAGNOSTICS)

// Java: CheckJSDocStyle(AbstractCompiler)
func NewCheckJSDocStyle(compiler AbstractCompiler) *CheckJSDocStyle {
	return &CheckJSDocStyle{compiler: compiler}
}

func (*CheckJSDocStyle) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckJSDocStyle(c)
}

// Java: CheckJSDocStyle.process(Node, Node)
func (c *CheckJSDocStyle) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
	NodeTraversal_Traverse(c.compiler, externs, &CheckJSDocStyle_ExternsCallback{})
}

// Java: CheckJSDocStyle.visit(NodeTraversal, Node, Node)
func (c *CheckJSDocStyle) Visit(t *NodeTraversal, n rhino.Node, unused rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		c.VisitFunction(t, ast.MustFunction(n))
	case rhino.Token_CLASS:
		c.VisitClass(t, n)
	case rhino.Token_ASSIGN:
		c.CheckStyleForPrivateProperties(t, n)
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_STRING_KEY:
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		// Don't need to call visitFunction because this JSDoc will be visited when the function is
		// visited.
		if !NodeUtil_GetEnclosingClass(n).IsNil() {
			c.CheckStyleForPrivateProperties(t, n)
		}
	case rhino.Token_SCRIPT:
		c.CheckLicenseComment(t, n)
	default:
		c.VisitNonFunction(t, n)
	}
}

// Java: CheckJSDocStyle.checkForAtSignCodePresence(NodeTraversal, Node, JSDocInfo)
func (c *CheckJSDocStyle) CheckForAtSignCodePresence(t *NodeTraversal, n rhino.Node, jsDoc *rhino.JSDocInfo) {
	if jsDoc != nil && jsDoc.IsAtSignCodePresent() {
		t.Report(n, CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE)
	}
}

// Java: CheckJSDocStyle.visitNonFunction(NodeTraversal, Node)
func (c *CheckJSDocStyle) VisitNonFunction(t *NodeTraversal, n rhino.Node) {
	c.CheckForAtSignCodePresence(t, n, n.GetJSDocInfo())
}

// Java: CheckJSDocStyle.checkStyleForPrivateProperties(NodeTraversal, Node)
func (c *CheckJSDocStyle) CheckStyleForPrivateProperties(t *NodeTraversal, n rhino.Node) {
	c.CheckForAtSignCodePresence(t, n, NodeUtil_GetBestJSDocInfo(n))
}

// Java: CheckJSDocStyle.checkNoTypeOnGettersAndSetters(NodeTraversal, Node, JSDocInfo)
func CheckJSDocStyle_CheckNoTypeOnGettersAndSetters(
	t *NodeTraversal, function rhino.Node, jsDoc *rhino.JSDocInfo) {
	if jsDoc != nil && jsDoc.HasType() && function.GetGrandparent().IsClassMembers() {
		memberNode := function.GetParent()
		if memberNode.IsSetterDef() || memberNode.IsGetterDef() {
			t.Report(function, CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER)
		}
	}
}

// Java: CheckJSDocStyle.visitFunction(NodeTraversal, Node)
func (c *CheckJSDocStyle) VisitFunction(t *NodeTraversal, fn ast.Function) {
	fnNode := fn.Node()
	jsDoc := NodeUtil_GetBestJSDocInfo(fnNode)

	c.CheckForAtSignCodePresence(t, fnNode, jsDoc)

	hasInline := c.HasAnyInlineJsDoc(fn)
	if jsDoc == nil && !hasInline {
		c.CheckMissingJsDoc(t, fn)
	} else {
		if t.InGlobalScope() ||
			hasInline ||
			!jsDoc.GetParameterNames().IsEmpty() ||
			jsDoc.HasReturnType() ||
			jsDoc.IsOverride() {
			c.CheckParams(t, fn, jsDoc)
		}
		CheckJSDocStyle_CheckNoTypeOnGettersAndSetters(t, fnNode, jsDoc)
		c.CheckReturn(fn, jsDoc)
	}
}

// Java: CheckJSDocStyle.visitClass(NodeTraversal, Node)
func (c *CheckJSDocStyle) VisitClass(t *NodeTraversal, cls rhino.Node) {
	jsDoc := NodeUtil_GetBestJSDocInfo(cls)
	c.CheckForAtSignCodePresence(t, cls, jsDoc)
	if jsDoc != nil && jsDoc.IsConstructor() {
		t.Report(cls, CheckJSDocStyle_CLASS_DISALLOWED_JSDOC)
	}
}

// Java: CheckJSDocStyle.checkMissingJsDoc(NodeTraversal, Node)
func (c *CheckJSDocStyle) CheckMissingJsDoc(t *NodeTraversal, fn ast.Function) {
	if c.IsFunctionThatShouldHaveJsDoc(t, fn) && !c.IsTestMethod(fn.Node()) {
		t.Report(fn.Node(), CheckJSDocStyle_MISSING_JSDOC)
	}
}

// IsFunctionThatShouldHaveJsDoc reports whether the given function should
// have JSDoc. True if it's a function declared in the global scope, or a
// method on a class which is declared in the global scope.
//
// Java: CheckJSDocStyle.isFunctionThatShouldHaveJsDoc(NodeTraversal, Node)
func (c *CheckJSDocStyle) IsFunctionThatShouldHaveJsDoc(t *NodeTraversal, fn ast.Function) bool {
	if !(t.InGlobalHoistScope() || t.InModuleScope()) {
		// TODO(b/233631820): this should check for the module hoist scope instead
		return false
	}
	function := fn.Node()
	if NodeUtil_IsFunctionDeclaration(function) {
		return true
	}
	parent := function.GetParent()
	grandparent := parent.GetParent()
	if NodeUtil_IsNameDeclaration(grandparent) || parent.IsAssign() || parent.IsExport() {
		return true
	}

	if grandparent.IsClassMembers() {
		if parent.IsMemberFunctionDef() {
			// A constructor with no parameters doesn't need JSDoc,
			// but all other member functions do.
			return !c.IsConstructorWithoutParameters(fn)
		} else if parent.IsGetterDef() || parent.IsSetterDef() {
			return true
		}
	}

	return grandparent.IsObjectLit() &&
		NodeUtil_IsCallTo(grandparent.GetParent(), "Polymer")
}

// IsTestMethod reports whether this is a test method (test* or
// setup/teardown) that does not require JSDoc.
//
// Java: CheckJSDocStyle.isTestMethod(Node)
func (c *CheckJSDocStyle) IsTestMethod(function rhino.Node) bool {
	bestLValue := NodeUtil_GetBestLValue(function)
	if bestLValue.IsNil() {
		return false
	}
	name := NodeUtil_GetBestLValueName(bestLValue)
	return name != nil && ExportTestFunctions_IsTestFunction(name)
}

// Java: CheckJSDocStyle.isConstructorWithoutParameters(Node)
func (c *CheckJSDocStyle) IsConstructorWithoutParameters(fn ast.Function) bool {
	_, hasParam := fn.Params().FirstParam()
	return !hasParam && NodeUtil_IsEs6Constructor(fn.Node())
}

// Java: CheckJSDocStyle.checkParams(NodeTraversal, Node, JSDocInfo)
func (c *CheckJSDocStyle) CheckParams(t *NodeTraversal, fn ast.Function, jsDoc *rhino.JSDocInfo) {
	if jsDoc != nil && jsDoc.GetType() != nil {
		// Sometimes functions are declared with @type {function(Foo, Bar)} instead of
		//   @param {Foo} foo
		//   @param {Bar} bar
		// which is fine.
		return
	}

	if jsDoc == nil || jsDoc.GetParameterNames().IsEmpty() {
		c.CheckInlineParams(t, fn)
		return
	}

	paramList := fn.Params()
	paramNames := jsDoc.GetParameterNames()
	if !paramList.Node().HasXChildren(paramNames.Size()) {
		c.compiler.Report(
			JSError_MakeWithNode(
				paramList.Node(),
				CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS,
				""))
		return
	}

	param := paramList.Node().GetFirstChild()
	for name := range paramNames.All() {
		if param.GetJSDocInfo() != nil {
			t.Report(param, CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES)
		}
		paramType := jsDoc.GetParameterType(name)
		if c.CheckParam(t, param, name, true, paramType) {
			return
		}
		param = param.GetNext()
	}
}

// CheckInlineParams checks that the inline type annotations are correct.
//
// Java: CheckJSDocStyle.checkInlineParams(NodeTraversal, Node, JSDocInfo)
func (c *CheckJSDocStyle) CheckInlineParams(t *NodeTraversal, fn ast.Function) {
	for param := range fn.Params().Params() {
		paramNode := param.Node()
		var jsDoc *rhino.JSDocInfo
		if defVal, ok := param.AsDefaultValue(); ok {
			jsDoc = defVal.Target().Node().GetJSDocInfo()
		} else {
			jsDoc = paramNode.GetJSDocInfo()
		}
		if jsDoc == nil {
			c.compiler.Report(
				JSError_MakeWithNode(
					paramNode,
					CheckJSDocStyle_MISSING_PARAMETER_JSDOC,
					""))
			return
		}
		paramType := jsDoc.GetType()
		jutil.CheckNotNullMsg(paramType, "Inline JSDoc info should always have a type")
		c.CheckParam(t, paramNode, "", false, paramType)
	}
}

// CheckParam checks that the given parameter node has the given name, and
// that the given type is compatible.
//
// param: if this is a non-NAME node, such as a destructuring pattern, skip
// the name check. hasName: if false, skip the name check. Returns whether a
// warning was reported.
//
// Java: CheckJSDocStyle.checkParam(NodeTraversal, Node, String, JSTypeExpression)
func (c *CheckJSDocStyle) CheckParam(
	t *NodeTraversal, param rhino.Node, name string, hasName bool, paramType *rhino.JSTypeExpression) bool {
	var nameOptional bool
	nodeToCheck := param
	if defVal, ok := ast.AsDefaultValue(param); ok {
		nodeToCheck = defVal.Target().Node()
		nameOptional = true
	} else if paramName, ok := ast.AsName(param); ok {
		nameOptional = strings.HasPrefix(paramName.Ident(), "opt_")
	} else {
		jutil.CheckStateMsg1(param.IsDestructuringPattern() || param.IsRest(), "%s", param)
	}

	reportedName := name
	if !hasName || !nodeToCheck.IsName() {
		// Skip the name check, but use "<unknown name>" for other errors that might be reported.
		reportedName = "<unknown name>"
	} else if !nodeToCheck.MatchesQualifiedName(name) {
		t.Report(nodeToCheck, CheckJSDocStyle_INCORRECT_PARAM_NAME)
		return true
	}

	if !nameOptional {
		return false
	}

	if paramType != nil && paramType.IsOptionalArg() {
		return false
	}

	errorSource := nodeToCheck
	if paramType != nil {
		errorSource = paramType.GetRoot()
	}
	t.Report(errorSource, CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL, reportedName)
	return true
}

// Java: CheckJSDocStyle.isDefaultAssignedParamWithInlineJsDoc(Node)
func CheckJSDocStyle_IsDefaultAssignedParamWithInlineJsDoc(param rhino.Node) bool {
	if defVal, ok := ast.AsDefaultValue(param); ok {
		if nameTarget, ok := defVal.Target().AsName(); ok {
			return nameTarget.Node().GetJSDocInfo() != nil
		}
	}
	return false
}

// Java: CheckJSDocStyle.hasAnyInlineJsDoc(Node)
func (c *CheckJSDocStyle) HasAnyInlineJsDoc(fn ast.Function) bool {
	if fn.NameNode().Node().GetJSDocInfo() != nil {
		// Inline return annotation.
		return true
	}
	for param := range fn.Params().Params() {
		pNode := param.Node()
		if pNode.GetJSDocInfo() != nil || CheckJSDocStyle_IsDefaultAssignedParamWithInlineJsDoc(pNode) {
			return true
		}
	}
	return false
}

// Java: CheckJSDocStyle.checkReturn(Node, JSDocInfo)
func (c *CheckJSDocStyle) CheckReturn(fn ast.Function, jsDoc *rhino.JSDocInfo) {
	if jsDoc != nil && (jsDoc.HasType() || jsDoc.IsConstructor() || jsDoc.HasReturnType()) {
		return
	}

	fnNode := fn.Node()
	if NodeUtil_IsEs6Constructor(fnNode) {
		// ES6 class constructors should never have "@return".
		return
	}

	if fn.NameNode().Node().GetJSDocInfo() != nil {
		return
	}

	finder := &CheckJSDocStyle_FindNonTrivialReturn{}
	NodeTraversal_Traverse(c.compiler, fn.Body().Node(), finder)
	if finder.found {
		c.compiler.Report(
			JSError_MakeWithNode(
				fnNode,
				CheckJSDocStyle_MISSING_RETURN_JSDOC,
				""))
	}
}

// Java: CheckJSDocStyle.FindNonTrivialReturn
type CheckJSDocStyle_FindNonTrivialReturn struct {
	NodeTraversal_AbstractPreOrderCallbackBase

	found bool
}

// Java: CheckJSDocStyle.FindNonTrivialReturn.shouldTraverse(NodeTraversal, Node, Node)
func (f *CheckJSDocStyle_FindNonTrivialReturn) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if f.found {
		return false
	}

	// Shallow traversal, since we don't need to inspect within functions or expressions.
	if NodeUtil_IsShallowStatementTree(parent) {
		if ret, ok := ast.AsReturn(n); ok {
			if _, hasVal := ret.Value(); hasVal {
				f.found = true
				return false
			}
		}
		return true
	}
	return false
}

// Java: CheckJSDocStyle.ExternsCallback
type CheckJSDocStyle_ExternsCallback struct{}

// Java: CheckJSDocStyle.ExternsCallback.shouldTraverse(NodeTraversal, Node, Node)
func (*CheckJSDocStyle_ExternsCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return parent.IsNil() || n.IsScript()
}

// Java: CheckJSDocStyle.ExternsCallback.visit(NodeTraversal, Node, Node)
func (*CheckJSDocStyle_ExternsCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsScript(n); ok {
		info := n.GetJSDocInfo()
		if info == nil || !info.IsExterns() {
			t.Report(n, CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED)
		}
	}
}

// Java: CheckJSDocStyle.checkLicenseComment(NodeTraversal, Node)
func (c *CheckJSDocStyle) CheckLicenseComment(t *NodeTraversal, n rhino.Node) {
	jsDoc := n.GetJSDocInfo()
	if jsDoc == nil {
		return
	}
	license := jsDoc.GetLicense()
	if license != nil && strings.Contains(*license, "@externs") {
		t.Report(n, CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS)
	}
}
