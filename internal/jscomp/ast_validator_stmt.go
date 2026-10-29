// Port of com/google/javascript/jscomp/AstValidator.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: AstValidator.validateRoot(Node)
func (v *AstValidator) ValidateRoot(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_ROOT, n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	v.ValidateCodeRoot(n.GetFirstChild())
	v.ValidateCodeRoot(n.GetLastChild())
	v.checkGrammar(n)
}

// Java: AstValidator.validateCodeRoot(Node)
func (v *AstValidator) ValidateCodeRoot(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_ROOT, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateScript(c)
	}
}

// Java: AstValidator.validateScript(Node)
func (v *AstValidator) ValidateScript(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_SCRIPT, n)
	v.ValidateHasSourceName(n)
	v.ValidateHasInputId(n)
	v.currentScript = n
	if n.HasChildren() && n.GetFirstChild().IsModuleBody() {
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 1)
		v.ValidateModuleContents(n.GetFirstChild())
	} else {
		v.ValidateStatements(n.GetFirstChild())
	}
	if v.isScriptFeatureValidationEnabled {
		v.ValidateScriptFeatureSet(n)
	}
}

// ValidateScriptFeatureSet confirms that every SCRIPT node’s FEATURE_SET <=
// compiler's allowable featureSet. This is possbile because with
// go/accurately-maintain-script-node-featureSet, each transpiler pass
// updates script features anytime it updates the compiler's allowable
// features.
//
// Java: AstValidator.validateScriptFeatureSet(Node)
func (v *AstValidator) ValidateScriptFeatureSet(scriptNode rhino.Node) {
	if !scriptNode.IsScript() {
		v.Violation("Not a script node", scriptNode)
		// unit tests for this pass perform "Negaive Testing" (i.e pass
		// non-script nodes to {@code validateScript}) and expect a violation
		// {@code expectInvalid(n, Check.SCRIPT);} report violation and return
		// here instead of crashing below in {@code
		// NodeUtil.getFeatureSetofScript} for test to complete
		return
	}
	scriptFeatures := NodeUtil_GetFeatureSetOfScript(v.currentScript)
	allowableFeatures := v.compiler.GetAllowableFeatures()
	if scriptFeatures == nil || allowableFeatures == nil {
		return
	}
	if !allowableFeatures.Contains(scriptFeatures) {
		if !scriptNode.IsFromExterns() {
			// Skip this check for externs because we don't need to complete
			// transpilation on externs, and currently only transpile externs
			// so that we can typecheck ES6+ features in externs.
			differentFeatures := scriptFeatures.Without(allowableFeatures)
			v.Violation(
				"SCRIPT node contains these unallowable features:"+differentFeatures.String(),
				v.currentScript)
		}
	}
}

// Java: AstValidator.validateModuleContents(Node)
func (v *AstValidator) ValidateModuleContents(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_MODULE_BODY, n)
	v.ValidateStatements(n.GetFirstChild())
}

// Java: AstValidator.validateStatements(Node)
func (v *AstValidator) ValidateStatements(n rhino.Node) {
	for !n.IsNil() {
		v.ValidateStatement(n)
		n = n.GetNext()
	}
}

// Java: AstValidator.validateStatement(Node)
func (v *AstValidator) ValidateStatement(n rhino.Node) {
	v.ValidateStatementWithAmbient(n, false)
}

// ValidateStatementWithAmbient validates a statement node and its
// children.
//
// isAmbient is whether this statement comes from TS ambient `declare
// [...]`.
//
// Java: AstValidator.validateStatement(Node, boolean)
func (v *AstValidator) ValidateStatementWithAmbient(n rhino.Node, isAmbient bool) {
	switch n.GetToken() {
	case rhino.Token_LABEL:
		v.ValidateLabel(n)
	case rhino.Token_BLOCK:
		v.ValidateBlock(n)
	case rhino.Token_FUNCTION:
		if isAmbient {
			v.ValidateFunctionSignature(n)
		} else {
			v.ValidateFunctionStatement(n)
		}
	case rhino.Token_WITH:
		v.ValidateWith(n)
	case rhino.Token_FOR:
		v.ValidateFor(n)
	case rhino.Token_FOR_IN:
		v.ValidateForIn(n)
	case rhino.Token_FOR_OF:
		v.ValidateForOf(n)
	case rhino.Token_FOR_AWAIT_OF:
		v.ValidateForAwaitOf(n)
	case rhino.Token_WHILE:
		v.ValidateWhile(n)
	case rhino.Token_DO:
		v.ValidateDo(n)
	case rhino.Token_SWITCH:
		v.ValidateSwitch(n)
	case rhino.Token_IF:
		v.ValidateIf(n)
	case rhino.Token_CONST, rhino.Token_VAR, rhino.Token_LET:
		v.ValidateNameDeclarationHelper(n, n.GetToken(), n)
	case rhino.Token_EXPR_RESULT:
		v.ValidateExprStmt(n)
	case rhino.Token_RETURN:
		v.ValidateReturn(n)
	case rhino.Token_THROW:
		v.ValidateThrow(n)
	case rhino.Token_TRY:
		v.ValidateTry(n)
	case rhino.Token_BREAK:
		v.ValidateBreak(n)
	case rhino.Token_CONTINUE:
		v.ValidateContinue(n)
	case rhino.Token_EMPTY, rhino.Token_DEBUGGER:
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_CLASS:
		v.ValidateClassDeclaration(n, isAmbient)
	case rhino.Token_IMPORT:
		v.ValidateImport(n)
	case rhino.Token_EXPORT:
		v.ValidateExport(n, isAmbient)
	case rhino.Token_INTERFACE:
		v.ValidateInterface(n)
	case rhino.Token_ENUM:
		v.ValidateEnum(n)
	case rhino.Token_TYPE_ALIAS:
		v.ValidateTypeAlias(n)
	case rhino.Token_DECLARE:
		v.ValidateAmbientDeclaration(n)
	case rhino.Token_NAMESPACE:
		v.ValidateNamespace(n, isAmbient)
	case rhino.Token_MODULE_BODY:
		// Uncommon case where a module body is not the first child of a
		// script. This may happen in a specific circumstance where the {@code
		// LateEs6ToEs3Rewriter} pass injects code above a module body. Valid
		// only when skipNonTranspilationPasses=true and
		// setWrapGoogModulesForWhitespaceOnly=false
		// TODO: b/294420383 Ideally the LateEs6ToEs3Rewriter pass should not
		// inject code above the module body node
		if v.compiler.GetOptions().GetSkipNonTranspilationPasses() {
			if !v.compiler.GetOptions().ShouldWrapGoogModulesForWhitespaceOnly() {
				v.ValidateModuleContents(n)
				return
			}
		}
		v.Violation("Expected statement but was "+n.GetToken().String()+".", n)
	default:
		if n.IsModuleBody() && v.compiler.GetOptions().GetSkipNonTranspilationPasses() {
			jutil.CheckStateMsg(
				!v.compiler.GetOptions().ShouldWrapGoogModulesForWhitespaceOnly(),
				"Modules can exist in transpiler only if setWrapGoogModulesForWhitespaceOnly is"+
					" false")
			v.ValidateModuleContents(n)
			return
		}
		v.Violation("Expected statement but was "+n.GetToken().String()+".", n)
	}
}

// Java: AstValidator.validateBlock(Node)
func (v *AstValidator) ValidateBlock(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_BLOCK, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateStatement(c)
	}
}

// Java: AstValidator.validateLabel(Node)
func (v *AstValidator) ValidateLabel(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_LABEL, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateLabelName(n.GetFirstChild())
	v.ValidateStatement(n.GetLastChild())
}

// Java: AstValidator.validateLabelName(Node)
func (v *AstValidator) ValidateLabelName(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_LABEL_NAME, n)
	v.ValidateNonEmptyString(n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
}

// Java: AstValidator.validateFor(Node)
func (v *AstValidator) ValidateFor(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_FOR, n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 4)
	target := n.GetFirstChild()
	if NodeUtil_IsNameDeclaration(target) {
		v.ValidateNameDeclarationHelper(n, target.GetToken(), target)
	} else {
		v.ValidatePseudoExpression(target, rhino.Token_EMPTY)
	}
	v.ValidatePseudoExpression(n.GetSecondChild(), rhino.Token_EMPTY)
	v.ValidatePseudoExpression(n.GetChildAtIndex(2), rhino.Token_EMPTY)
	v.ValidateBlock(n.GetLastChild())
}

// Java: AstValidator.validateForIn(Node)
func (v *AstValidator) ValidateForIn(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_FOR_IN, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateEnhancedForVarOrAssignmentTarget(n, n.GetFirstChild())
	v.ValidateExpression(n.GetSecondChild())
	v.ValidateBlock(n.GetLastChild())
}

// Java: AstValidator.validateForOf(Node)
func (v *AstValidator) ValidateForOf(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_FOR_OF, n)
	v.ValidateNodeType(rhino.Token_FOR_OF, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateEnhancedForVarOrAssignmentTarget(n, n.GetFirstChild())
	v.ValidateExpression(n.GetSecondChild())
	v.ValidateBlock(n.GetLastChild())
}

// Java: AstValidator.validateForAwaitOf(Node)
func (v *AstValidator) ValidateForAwaitOf(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_FOR_AWAIT_OF, n)
	v.ValidateNodeType(rhino.Token_FOR_AWAIT_OF, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateEnhancedForVarOrAssignmentTarget(n, n.GetFirstChild())
	v.ValidateExpression(n.GetSecondChild())
	v.ValidateBlock(n.GetLastChild())
	v.ValidateAwaitWithinAsyncFunction(n)
}

// Java: AstValidator.validateEnhancedForVarOrAssignmentTarget(Node, Node)
func (v *AstValidator) ValidateEnhancedForVarOrAssignmentTarget(forNode rhino.Node, n rhino.Node) {
	if NodeUtil_IsNameDeclaration(n) {
		// Only one NAME can be declared for FOR-IN and FOR_OF expressions.
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 1)
		v.ValidateNameDeclarationHelper(forNode, n.GetToken(), n)
	} else {
		v.ValidateLHS(n.GetParent().GetToken(), n)
	}
}

// Java: AstValidator.validateWith(Node)
func (v *AstValidator) ValidateWith(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_WITH, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateBlock(n.GetLastChild())
}

// Java: AstValidator.validateWhile(Node)
func (v *AstValidator) ValidateWhile(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_WHILE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateBlock(n.GetLastChild())
}

// Java: AstValidator.validateDo(Node)
func (v *AstValidator) ValidateDo(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_DO, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateBlock(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateIf(Node)
func (v *AstValidator) ValidateIf(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_IF, n)
	v.ValidateProperties(n)
	v.ValidateChildCountIn(n, 2, 3)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateBlock(n.GetSecondChild())
	if n.HasXChildren(3) {
		v.ValidateBlock(n.GetLastChild())
	}
}

// Java: AstValidator.validateExprStmt(Node)
func (v *AstValidator) ValidateExprStmt(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_EXPR_RESULT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateReturn(Node)
func (v *AstValidator) ValidateReturn(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_RETURN, n)
	v.ValidateProperties(n)
	v.ValidateMaximumChildCount(n, 1)
	if n.HasChildren() {
		v.ValidateExpression(n.GetFirstChild())
	}
}

// Java: AstValidator.validateThrow(Node)
func (v *AstValidator) ValidateThrow(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_THROW, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateBreak(Node)
func (v *AstValidator) ValidateBreak(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_BREAK, n)
	v.ValidateProperties(n)
	v.ValidateMaximumChildCount(n, 1)
	if n.HasChildren() {
		v.ValidateLabelName(n.GetFirstChild())
	}
}

// Java: AstValidator.validateContinue(Node)
func (v *AstValidator) ValidateContinue(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_CONTINUE, n)
	v.ValidateProperties(n)
	v.ValidateMaximumChildCount(n, 1)
	if n.HasChildren() {
		v.ValidateLabelName(n.GetFirstChild())
	}
}

// Java: AstValidator.validateTry(Node)
func (v *AstValidator) ValidateTry(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_TRY, n)
	v.ValidateProperties(n)
	v.ValidateChildCountIn(n, 2, 3)
	v.ValidateBlock(n.GetFirstChild())

	seenCatchOrFinally := false

	// Validate catch
	catches := n.GetSecondChild()
	v.ValidateNodeType(rhino.Token_BLOCK, catches)
	v.ValidateProperties(catches)
	v.ValidateMaximumChildCount(catches, 1)
	if catches.HasChildren() {
		v.ValidateCatch(catches.GetFirstChild())
		seenCatchOrFinally = true
	}

	// Validate finally
	if n.HasXChildren(3) {
		v.ValidateBlock(n.GetLastChild())
		seenCatchOrFinally = true
	}

	if !seenCatchOrFinally {
		v.Violation("Missing catch or finally for try statement.", n)
	}
}

// Java: AstValidator.validateCatch(Node)
func (v *AstValidator) ValidateCatch(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_CATCH, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	caught := n.GetFirstChild()
	if caught.IsName() {
		v.ValidateName(caught)
	} else if caught.IsArrayPattern() {
		v.ValidateArrayPattern(rhino.Token_CATCH, caught)
	} else if caught.IsObjectPattern() {
		v.ValidateObjectPattern(rhino.Token_CATCH, caught)
	} else if caught.IsEmpty() {
		v.ValidateNoCatchBinding(caught)
	} else {
		v.Violation("Unexpected catch binding: "+caught.String(), n)
	}
	v.ValidateBlock(n.GetLastChild())
}

// Java: AstValidator.validateNoCatchBinding(Node)
func (v *AstValidator) ValidateNoCatchBinding(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
}

// Java: AstValidator.validateSwitch(Node)
func (v *AstValidator) ValidateSwitch(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_SWITCH, n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateNodeType(rhino.Token_SWITCH_BODY, n.GetSecondChild())
	cases := n.GetSecondChild()
	defaults := 0
	for c := cases.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateSwitchMember(c)
		if c.IsDefaultCase() {
			defaults++
		}
	}
	if defaults > 1 {
		v.Violation("Expected at most 1 'default' in switch but was "+strconv.Itoa(defaults), n)
	}
}

// Java: AstValidator.validateSwitchMember(Node)
func (v *AstValidator) ValidateSwitchMember(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CASE:
		v.ValidateCase(n)
	case rhino.Token_DEFAULT_CASE:
		v.ValidateDefaultCase(n)
	default:
		v.Violation("Expected switch member but was "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateDefaultCase(Node)
func (v *AstValidator) ValidateDefaultCase(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_DEFAULT_CASE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateBlock(n.GetLastChild())
}

// Java: AstValidator.validateCase(Node)
func (v *AstValidator) ValidateCase(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_CASE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateBlock(n.GetLastChild())
}
