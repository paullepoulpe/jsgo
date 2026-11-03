// Port of com/google/javascript/jscomp/AstValidator.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: AstValidator.validateRoot(Node)
func (v *AstValidator) ValidateRoot(n rhino.Node) {
	isRoot := v.validateNodeType(rhino.Token_ROOT, n)
	if v.validatePropsAndCount(n, 2) && isRoot {
		root := ast.MustRoot(n)
		externsRoot, _ := root.FirstItem()
		codeRoot, _ := root.LastItem()
		v.ValidateCodeRoot(externsRoot.Node())
		v.ValidateCodeRoot(codeRoot.Node())
	} else {
		v.ValidateCodeRoot(n.GetFirstChild())
		v.ValidateCodeRoot(n.GetLastChild())
	}
	v.checkGrammar(n)
}

// Java: AstValidator.validateCodeRoot(Node)
func (v *AstValidator) ValidateCodeRoot(n rhino.Node) {
	if v.validateNodeType(rhino.Token_ROOT, n) {
		for item := range ast.MustRoot(n).Items() {
			v.ValidateScript(item.Node())
		}
		return
	}
	for c := range n.Children() {
		v.ValidateScript(c)
	}
}

// Java: AstValidator.validateScript(Node)
func (v *AstValidator) ValidateScript(n rhino.Node) {
	isScript := v.validateNodeType(rhino.Token_SCRIPT, n)
	v.validateHasSourceName(n)
	v.validateHasInputId(n)
	v.currentScript = n
	if isScript {
		script := ast.MustScript(n)
		if first, ok := script.FirstItem(); ok && first.Node().IsModuleBody() {
			v.validatePropsAndCount(n, 1)
			v.validateModuleContents(ast.MustModuleBody(first.Node()))
		} else {
			for item := range script.Items() {
				v.ValidateStatement(item.Node())
			}
		}
	} else if n.HasChildren() && n.GetFirstChild().IsModuleBody() {
		v.validatePropsAndCount(n, 1)
		v.validateModuleContents(ast.MustModuleBody(n.GetFirstChild()))
	} else {
		for c := range n.Children() {
			v.ValidateStatement(c)
		}
	}
	if v.isScriptFeatureValidationEnabled {
		v.validateScriptFeatureSet(n)
	}
}

// validateScriptFeatureSet confirms that every SCRIPT node's FEATURE_SET <=
// compiler's allowable featureSet.
//
// Java: AstValidator.validateScriptFeatureSet(Node)
func (v *AstValidator) validateScriptFeatureSet(scriptNode rhino.Node) {
	if !scriptNode.IsScript() {
		// Unit tests pass non-SCRIPT nodes to ValidateScript and expect a
		// violation rather than a panic in NodeUtil_GetFeatureSetOfScript.
		v.Violation("Not a script node", scriptNode)
		return
	}
	scriptFeatures := NodeUtil_GetFeatureSetOfScript(v.currentScript)
	allowableFeatures := v.compiler.GetAllowableFeatures()
	if scriptFeatures.IsZero() || allowableFeatures.IsZero() {
		return
	}
	if !allowableFeatures.Contains(scriptFeatures) && !scriptNode.IsFromExterns() {
		differentFeatures := scriptFeatures.Without(allowableFeatures)
		v.Violation(
			"SCRIPT node contains these unallowable features:"+differentFeatures.String(),
			v.currentScript)
	}
}

// Java: AstValidator.validateModuleContents(Node)
func (v *AstValidator) validateModuleContents(mod ast.ModuleBody) {
	for stmt := range mod.Statements() {
		v.ValidateStatement(stmt.Node())
	}
}

// Java: AstValidator.validateStatement(Node)
func (v *AstValidator) ValidateStatement(n rhino.Node) {
	v.validateStatementWithAmbient(n, false)
}

// validateStatementWithAmbient validates a statement node and its children.
// isAmbient reports whether this statement comes from TS ambient `declare [...]`.
//
// Java: AstValidator.validateStatement(Node, boolean)
func (v *AstValidator) validateStatementWithAmbient(n rhino.Node, isAmbient bool) {
	switch n.GetToken() {
	case rhino.Token_LABEL:
		v.validateLabel(ast.MustLabel(n))
	case rhino.Token_BLOCK:
		v.validateBlock(n)
	case rhino.Token_FUNCTION:
		if isAmbient {
			v.validateFunctionSignature(n)
		} else {
			v.validateFunctionStatement(ast.MustFunction(n))
		}
	case rhino.Token_WITH:
		v.validateWith(ast.MustWith(n))
	case rhino.Token_FOR:
		v.validateFor(ast.MustFor(n))
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		v.validateEnhancedFor(ast.MustEnhancedFor(n))
	case rhino.Token_WHILE:
		v.validateWhile(ast.MustWhile(n))
	case rhino.Token_DO:
		v.validateDo(ast.MustDo(n))
	case rhino.Token_SWITCH:
		v.validateSwitch(ast.MustSwitch(n))
	case rhino.Token_IF:
		v.validateIf(ast.MustIf(n))
	case rhino.Token_CONST, rhino.Token_VAR, rhino.Token_LET:
		v.validateNameDeclaration(n, ast.MustNameDeclaration(n))
	case rhino.Token_EXPR_RESULT:
		v.validateExprStmt(ast.MustExprStmt(n))
	case rhino.Token_RETURN:
		v.validateReturn(ast.MustReturn(n))
	case rhino.Token_THROW:
		v.validateThrow(ast.MustThrow(n))
	case rhino.Token_TRY:
		v.validateTry(ast.MustTry(n))
	case rhino.Token_BREAK, rhino.Token_CONTINUE:
		v.validateJump(ast.MustJump(n))
	case rhino.Token_EMPTY, rhino.Token_DEBUGGER:
		v.validateChildlessNode(n)
	case rhino.Token_CLASS:
		v.validateClassDeclaration(ast.MustClass(n), isAmbient)
	case rhino.Token_IMPORT:
		v.validateImport(ast.MustImport(n))
	case rhino.Token_EXPORT:
		v.validateExport(ast.MustExport(n), isAmbient)
	case rhino.Token_INTERFACE:
		v.validateInterface(n)
	case rhino.Token_ENUM:
		v.validateEnum(n)
	case rhino.Token_TYPE_ALIAS:
		v.validateTypeAlias(n)
	case rhino.Token_DECLARE:
		v.validateAmbientDeclaration(n)
	case rhino.Token_NAMESPACE:
		v.validateNamespace(n, isAmbient)
	case rhino.Token_MODULE_BODY:
		// LateEs6ToEs3Rewriter may inject code above a MODULE_BODY when
		// skipNonTranspilationPasses=true and wrapGoogModulesForWhitespaceOnly=false.
		opts := v.compiler.GetOptions()
		if opts.GetSkipNonTranspilationPasses() && !opts.ShouldWrapGoogModulesForWhitespaceOnly() {
			v.validateModuleContents(ast.MustModuleBody(n))
			return
		}
		v.Violation("Expected statement but was "+n.GetToken().String()+".", n)
	default:
		v.Violation("Expected statement but was "+n.GetToken().String()+".", n)
	}
}

// Java: AstValidator.validateBlock(Node)
func (v *AstValidator) validateBlock(n rhino.Node) {
	if v.validateNodeType(rhino.Token_BLOCK, n) {
		for stmt := range ast.MustBlock(n).Statements() {
			v.ValidateStatement(stmt.Node())
		}
		return
	}
	for c := range n.Children() {
		v.ValidateStatement(c)
	}
}

// Java: AstValidator.validateLabel(Node)
func (v *AstValidator) validateLabel(lbl ast.Label) {
	n := lbl.Node()
	if v.validatePropsAndArity(n) {
		v.validateLabelName(lbl.Name().Node())
		v.ValidateStatement(lbl.Body().Node())
		return
	}
	v.validateLabelName(n.GetFirstChild())
	v.ValidateStatement(n.GetLastChild())
}

// Java: AstValidator.validateLabelName(Node)
func (v *AstValidator) validateLabelName(n rhino.Node) {
	v.validateNodeType(rhino.Token_LABEL_NAME, n)
	v.validateNonEmptyString(n)
	v.validatePropsAndArity(n)
}

// Java: AstValidator.validateFor(Node)
func (v *AstValidator) validateFor(forStmt ast.For) {
	n := forStmt.Node()
	v.validatePropsAndCount(n, 4)

	if init, ok := forStmt.Init(); !ok {
		v.validateChildlessNode(n.GetFirstChild())
	} else if decl, ok := init.AsNameDeclaration(); ok {
		v.validateNameDeclaration(n, decl)
	} else {
		v.validatePseudoExpression(init.Node(), rhino.Token_EMPTY)
	}

	if cond, ok := forStmt.Cond(); !ok {
		v.validateChildlessNode(n.GetSecondChild())
	} else {
		v.validatePseudoExpression(cond.Node(), rhino.Token_EMPTY)
	}

	if update, ok := forStmt.Update(); !ok {
		v.validateChildlessNode(n.GetChildAtIndex(2))
	} else {
		v.validatePseudoExpression(update.Node(), rhino.Token_EMPTY)
	}

	v.validateBlock(forStmt.Body().Node())
}

// Java: AstValidator.validateForIn / validateForOf / validateForAwaitOf
func (v *AstValidator) validateEnhancedFor(loop ast.EnhancedFor) {
	n := loop.Node()
	switch loop.Token() {
	case rhino.Token_FOR_OF:
		v.validateFeature(parser.FeatureSet_Feature_FOR_OF, n)
	case rhino.Token_FOR_AWAIT_OF:
		v.validateFeature(parser.FeatureSet_Feature_FOR_AWAIT_OF, n)
	}
	v.validatePropsAndArity(n)
	v.validateEnhancedForTarget(n, loop.Target())
	v.ValidateExpression(loop.Collection().Node())
	v.validateBlock(loop.AsLoop().Body().Node())
	if loop.Token() == rhino.Token_FOR_AWAIT_OF {
		v.validateAwaitWithinAsyncFunction(n)
	}
}

// Java: AstValidator.validateEnhancedForVarOrAssignmentTarget(Node, Node)
func (v *AstValidator) validateEnhancedForTarget(forNode rhino.Node, target ast.ForTarget) {
	if decl, ok := target.AsNameDeclaration(); ok {
		// Only one NAME can be declared for FOR-IN and FOR-OF expressions.
		v.validatePropsAndCount(decl.Node(), 1)
		v.validateNameDeclaration(forNode, decl)
		return
	}
	n := target.Node()
	v.validateLHS(n.GetParent().GetToken(), n)
}

// Java: AstValidator.validateWith(Node)
func (v *AstValidator) validateWith(stmt ast.With) {
	n := stmt.Node()
	if v.validatePropsAndArity(n) {
		v.ValidateExpression(stmt.Object().Node())
		v.validateBlock(stmt.Body().Node())
		return
	}
	v.ValidateExpression(n.GetFirstChild())
	v.validateBlock(n.GetLastChild())
}

// Java: AstValidator.validateWhile(Node)
func (v *AstValidator) validateWhile(stmt ast.While) {
	n := stmt.Node()
	if v.validatePropsAndArity(n) {
		v.ValidateExpression(stmt.Cond().Node())
		v.validateBlock(stmt.Body().Node())
		return
	}
	v.ValidateExpression(n.GetFirstChild())
	v.validateBlock(n.GetLastChild())
}

// Java: AstValidator.validateDo(Node)
func (v *AstValidator) validateDo(stmt ast.Do) {
	n := stmt.Node()
	if v.validatePropsAndArity(n) {
		v.validateBlock(stmt.Body().Node())
		v.ValidateExpression(stmt.Cond().Node())
		return
	}
	v.validateBlock(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateIf(Node)
func (v *AstValidator) validateIf(stmt ast.If) {
	n := stmt.Node()
	v.validateProperties(n)
	valid := v.validateChildCountIn(n, 2, 3)
	v.ValidateExpression(stmt.Cond().Node())
	v.validateBlock(stmt.Then().Node())
	if valid {
		if elseBlock, ok := stmt.Else(); ok {
			v.validateBlock(elseBlock.Node())
		}
	}
}

// Java: AstValidator.validateExprStmt(Node)
func (v *AstValidator) validateExprStmt(stmt ast.ExprStmt) {
	v.validatePropsAndArity(stmt.Node())
	v.ValidateExpression(stmt.Expr().Node())
}

// Java: AstValidator.validateReturn(Node)
func (v *AstValidator) validateReturn(stmt ast.Return) {
	n := stmt.Node()
	v.validateProperties(n)
	v.validateMaximumChildCount(n, 1)
	if val, ok := stmt.Value(); ok {
		v.ValidateExpression(val.Node())
	}
}

// Java: AstValidator.validateThrow(Node)
func (v *AstValidator) validateThrow(stmt ast.Throw) {
	v.validatePropsAndArity(stmt.Node())
	v.ValidateExpression(stmt.Value().Node())
}

// Java: AstValidator.validateBreak / validateContinue
func (v *AstValidator) validateJump(stmt ast.Jump) {
	n := stmt.Node()
	v.validateProperties(n)
	v.validateMaximumChildCount(n, 1)
	if lbl, ok := stmt.Label(); ok {
		v.validateLabelName(lbl.Node())
	}
}

// Java: AstValidator.validateTry(Node)
func (v *AstValidator) validateTry(stmt ast.Try) {
	n := stmt.Node()
	v.validateProperties(n)
	valid := v.validateChildCountIn(n, 2, 3)
	v.validateBlock(stmt.Body().Node())

	seenCatchOrFinally := false

	catches := stmt.Catches()
	v.validateNodeType(rhino.Token_BLOCK, catches.Node())
	v.validateProperties(catches.Node())
	v.validateMaximumChildCount(catches.Node(), 1)
	if catchClause, ok := catches.Catch(); ok {
		v.validateCatch(catchClause.Node())
		seenCatchOrFinally = true
	}

	if valid {
		if finallyBlock, ok := stmt.Finally(); ok {
			v.validateBlock(finallyBlock.Node())
			seenCatchOrFinally = true
		}
	}

	if !seenCatchOrFinally {
		v.Violation("Missing catch or finally for try statement.", n)
	}
}

// Java: AstValidator.validateCatch(Node)
func (v *AstValidator) validateCatch(n rhino.Node) {
	isCatch := v.validateNodeType(rhino.Token_CATCH, n)
	valid := v.validatePropsAndArity(n)

	var caught, body rhino.Node
	if isCatch && valid {
		c := ast.MustCatch(n)
		body = c.Body().Node()
		if param, ok := c.Param(); ok {
			caught = param.Node()
		} else {
			caught = n.GetFirstChild()
		}
	} else {
		caught = n.GetFirstChild()
		body = n.GetLastChild()
	}

	switch caught.GetToken() {
	case rhino.Token_NAME:
		v.validateName(caught)
	case rhino.Token_ARRAY_PATTERN:
		v.validateArrayPattern(rhino.Token_CATCH, caught)
	case rhino.Token_OBJECT_PATTERN:
		v.validateObjectPattern(rhino.Token_CATCH, caught)
	case rhino.Token_EMPTY:
		v.validateNoCatchBinding(caught)
	default:
		v.Violation("Unexpected catch binding: "+caught.String(), n)
	}
	v.validateBlock(body)
}

// Java: AstValidator.validateNoCatchBinding(Node)
func (v *AstValidator) validateNoCatchBinding(n rhino.Node) {
	v.validateFeature(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING, n)
	v.validatePropsAndArity(n)
}

// Java: AstValidator.validateSwitch(Node)
func (v *AstValidator) validateSwitch(sw ast.Switch) {
	n := sw.Node()
	v.validatePropsAndCount(n, 2)
	v.ValidateExpression(sw.Discriminant().Node())
	body := sw.Body()
	v.validateNodeType(rhino.Token_SWITCH_BODY, body.Node())
	defaults := 0
	for c := range body.Cases() {
		v.validateSwitchMember(c.Node())
		if c.Node().IsDefaultCase() {
			defaults++
		}
	}
	if defaults > 1 {
		v.Violation("Expected at most 1 'default' in switch but was "+strconv.Itoa(defaults), n)
	}
}

// Java: AstValidator.validateSwitchMember(Node)
func (v *AstValidator) validateSwitchMember(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CASE:
		v.validateCase(ast.MustCase(n))
	case rhino.Token_DEFAULT_CASE:
		v.validateDefaultCase(ast.MustDefaultCase(n))
	default:
		v.Violation("Expected switch member but was "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateDefaultCase(Node)
func (v *AstValidator) validateDefaultCase(c ast.DefaultCase) {
	n := c.Node()
	if v.validatePropsAndArity(n) {
		v.validateBlock(c.Body().Node())
		return
	}
	v.validateBlock(n.GetLastChild())
}

// Java: AstValidator.validateCase(Node)
func (v *AstValidator) validateCase(c ast.Case) {
	n := c.Node()
	if v.validatePropsAndArity(n) {
		v.ValidateExpression(c.Value().Node())
		v.validateBlock(c.Body().Node())
		return
	}
	v.ValidateExpression(n.GetFirstChild())
	v.validateBlock(n.GetLastChild())
}
