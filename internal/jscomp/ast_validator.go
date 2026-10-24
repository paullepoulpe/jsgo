// Port of com/google/javascript/jscomp/AstValidator.java @ closure-compiler v20261015

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
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AstValidator walks the AST and validates that the structure is correct.
type AstValidator struct {
	// Possible enhancements:
	// * verify NAME, LABEL_NAME, GETPROP property name and unquoted
	// object-literal keys are valid JavaScript identifiers.
	// * optionally verify every node has source location information.

	compiler         AbstractCompiler
	violationHandler AstValidator_ViolationHandler
	currentScript    rhino.Node

	// Perform type validation if this is enabled.
	typeValidationMode AstValidator_TypeInfoValidation

	// Validate that a SCRIPT's FeatureSet property includes all features if
	// this is enabled.
	isScriptFeatureValidationEnabled bool

	// TODO: varomodt - make this the default.
	// Validate that all required inlinings were performed.
	shouldValidateRequiredInlinings bool
}

// AstValidator_ViolationHandler is the violation handler.
type AstValidator_ViolationHandler interface {
	// Java: AstValidator.ViolationHandler.handleViolation(String, Node)
	HandleViolation(message string, n rhino.Node)
}

// AstValidator_TypeInfoValidation selects the type information that is
// validated.
type AstValidator_TypeInfoValidation int32

const (
	AstValidator_TypeInfoValidation_JSTYPE AstValidator_TypeInfoValidation = iota
	AstValidator_TypeInfoValidation_COLOR
	AstValidator_TypeInfoValidation_NONE
)

var astValidator_TypeInfoValidation_names = [...]string{"JSTYPE", "COLOR", "NONE"}

// AstValidator_TypeInfoValidation_values is TypeInfoValidation.values().
func AstValidator_TypeInfoValidation_values() []AstValidator_TypeInfoValidation {
	return []AstValidator_TypeInfoValidation{
		AstValidator_TypeInfoValidation_JSTYPE,
		AstValidator_TypeInfoValidation_COLOR,
		AstValidator_TypeInfoValidation_NONE,
	}
}

// Name is Enum.name().
func (v AstValidator_TypeInfoValidation) Name() string {
	return astValidator_TypeInfoValidation_names[v]
}

// String is Enum.toString().
func (v AstValidator_TypeInfoValidation) String() string { return v.Name() }

// Ordinal is Enum.ordinal().
func (v AstValidator_TypeInfoValidation) Ordinal() int { return int(v) }

// Java: AstValidator(AbstractCompiler, ViolationHandler, boolean, boolean)
func NewAstValidatorWithHandler(
	compiler AbstractCompiler,
	handler AstValidator_ViolationHandler,
	validateScriptFeatures bool,
	shouldValidateRequiredInlinings bool) *AstValidator {
	return &AstValidator{
		typeValidationMode:               AstValidator_TypeInfoValidation_NONE,
		compiler:                         compiler,
		violationHandler:                 handler,
		isScriptFeatureValidationEnabled: validateScriptFeatures,
		shouldValidateRequiredInlinings:  shouldValidateRequiredInlinings,
	}
}

// Java: AstValidator(AbstractCompiler)
func NewAstValidator(compiler AbstractCompiler) *AstValidator {
	return NewAstValidatorWithValidateScriptFeatures(compiler /* validateScriptFeatures= */, false)
}

// Java: AstValidator(AbstractCompiler, boolean)
func NewAstValidatorWithValidateScriptFeatures(compiler AbstractCompiler, validateScriptFeatures bool) *AstValidator {
	return NewAstValidatorWithHandler(
		compiler,
		astValidator_throwingViolationHandler{},
		validateScriptFeatures,
		compiler.GetOptions().GetShouldValidateRequiredInlinings() == base.Tri_TRUE)
}

// astValidator_throwingViolationHandler is the anonymous ViolationHandler
// of AstValidator(AbstractCompiler, boolean).
type astValidator_throwingViolationHandler struct{}

// Java: AstValidator(AbstractCompiler, boolean).new ViolationHandler() {...}.handleViolation(String, Node)
func (astValidator_throwingViolationHandler) HandleViolation(message string, n rhino.Node) {
	parent := " no parent "
	if n.HasParent() {
		parent = n.GetParent().ToStringTree()
	}
	panic(jutil.NewIllegalStateExceptionMsg(
		message +
			". Reference node:\n" +
			n.ToStringTree() +
			"\n Parent node:\n" +
			parent))
}

// SetTypeValidationMode enables or disables validation of type
// information.
//
// TODO(b/74537281): Currently only expressions are checked for type
// information. Do we need to do more?
//
// Java: AstValidator.setTypeValidationMode(TypeInfoValidation)
func (v *AstValidator) SetTypeValidationMode(mode AstValidator_TypeInfoValidation) *AstValidator {
	v.typeValidationMode = mode
	return v
}

// Java: AstValidator.process(Node, Node)
func (v *AstValidator) Process(externs rhino.Node, root rhino.Node) {
	if !externs.IsNil() {
		v.ValidateCodeRoot(externs)
		v.checkGrammar(externs)
	}
	if !root.IsNil() {
		v.ValidateCodeRoot(root)
		v.checkGrammar(root)
	}
}

// checkGrammar reports the nodes that do not follow the grammar of the
// typed AST view (internal/ast/grammar.txt), so that the grammar stays true
// to the trees the compiler builds.
//
// Go-only: the typed view is not in Java.
func (v *AstValidator) checkGrammar(root rhino.Node) {
	ast.Check(root, func(n rhino.Node, problem string) {
		v.Violation("AST grammar (internal/ast/grammar.txt): "+problem, n)
	})
}

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
				"SCRIPT node contains these unallowable features:"+differentFeatures.GetFeatures().String(),
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

// Java: AstValidator.validateExpression(Node)
func (v *AstValidator) ValidateExpression(n rhino.Node) {
	v.ValidateTypeInformation(n)
	v.ValidateRequiredInlinings(n)

	switch n.GetToken() {
	// Childless expressions
	case rhino.Token_NEW_TARGET:
		v.ValidateFeature(parser.FeatureSet_Feature_NEW_TARGET, n)
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_IMPORT_META:
		v.ValidateFeature(parser.FeatureSet_Feature_IMPORT_META, n)
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_FALSE, rhino.Token_NULL, rhino.Token_THIS, rhino.Token_TRUE:
		v.ValidateProperties(n)
		v.ValidateChildless(n)

	// General unary ops
	case rhino.Token_DELPROP, rhino.Token_POS, rhino.Token_NEG, rhino.Token_NOT, rhino.Token_TYPEOF,
		rhino.Token_VOID, rhino.Token_BITNOT, rhino.Token_CAST:
		v.ValidateUnaryOp(n)
	case rhino.Token_INC, rhino.Token_DEC:
		v.ValidateIncDecOp(n)

	// Assignments
	case rhino.Token_ASSIGN:
		v.ValidateAssignmentExpression(n)
	case rhino.Token_ASSIGN_EXPONENT:
		v.ValidateFeature(parser.FeatureSet_Feature_EXPONENT_OP, n)
		v.ValidateCompoundAssignmentExpression(n)
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD:
		v.ValidateCompoundAssignmentExpression(n)
	case rhino.Token_ASSIGN_COALESCE:
		v.ValidateFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP, n)
		v.ValidateFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT, n)
		v.ValidateCompoundAssignmentExpression(n)
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND:
		v.ValidateFeature(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT, n)
		v.ValidateCompoundAssignmentExpression(n)

	case rhino.Token_HOOK:
		v.ValidateTrinaryOp(n)

	// Node types that require special handling
	case rhino.Token_STRINGLIT:
		v.ValidateStringLit(n)
	case rhino.Token_NUMBER:
		v.ValidateNumber(n)
	case rhino.Token_BIGINT:
		v.ValidateBigInt(n)
	case rhino.Token_NAME:
		v.ValidateName(n)

	// General binary ops
	case rhino.Token_EXPONENT:
		v.ValidateFeature(parser.FeatureSet_Feature_EXPONENT_OP, n)
		v.ValidateBinaryOp(n)
	case rhino.Token_COALESCE:
		v.ValidateFeature(parser.FeatureSet_Feature_NULL_COALESCE_OP, n)
		v.ValidateBinaryOp(n)
	case rhino.Token_COMMA,
		rhino.Token_OR,
		rhino.Token_AND,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_BITAND,
		rhino.Token_EQ,
		rhino.Token_NE,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_LT,
		rhino.Token_GT,
		rhino.Token_LE,
		rhino.Token_GE,
		rhino.Token_INSTANCEOF,
		rhino.Token_IN,
		rhino.Token_LSH,
		rhino.Token_RSH,
		rhino.Token_URSH,
		rhino.Token_SUB,
		rhino.Token_ADD,
		rhino.Token_MUL,
		rhino.Token_MOD,
		rhino.Token_DIV:
		v.ValidateBinaryOp(n)

	case rhino.Token_GETELEM:
		v.ValidateGetElem(n)
	case rhino.Token_OPTCHAIN_GETELEM:
		v.ValidateOptChainGetElem(n)

	case rhino.Token_GETPROP:
		v.ValidateGetProp(n)
	case rhino.Token_OPTCHAIN_GETPROP:
		v.ValidateOptChainGetProp(n)

	case rhino.Token_ARRAYLIT:
		v.ValidateArrayLit(n)
	case rhino.Token_OBJECTLIT:
		v.ValidateObjectLit(n)
	case rhino.Token_REGEXP:
		v.ValidateRegExpLit(n)

	case rhino.Token_CALL:
		v.ValidateCall(n)
	case rhino.Token_OPTCHAIN_CALL:
		v.ValidateOptChainCall(n)
	case rhino.Token_NEW:
		v.ValidateNew(n)

	case rhino.Token_FUNCTION:
		v.ValidateRequiredInlinings(n)
		v.ValidateFunctionExpression(n)
	case rhino.Token_CLASS:
		v.ValidateClass(n)

	case rhino.Token_TEMPLATELIT:
		v.ValidateTemplateLit(n)
	case rhino.Token_TAGGED_TEMPLATELIT:
		v.ValidateTaggedTemplateLit(n)

	case rhino.Token_YIELD:
		v.ValidateYield(n)
	case rhino.Token_AWAIT:
		v.ValidateAwait(n)
	case rhino.Token_DYNAMIC_IMPORT:
		v.ValidateFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT, n)
		v.ValidateUnaryOp(n)

	default:
		v.Violation("Expected expression but was "+n.GetToken().String(), n)
	}
}

// ValidatePseudoExpression validates an expression or expresison-like
// construct.
//
// An expression-like construct (pseudoexpression) is an AST fragment that
// is valid in some, but not all, of the same contexts as true expressions.
// For example, a VANILLA_FOR permits EMPTY as its condition and increment
// expressions, even though EMPTY is not valid as an expression in general.
//
// allowedPseudoexpressions allows the caller to specify which
// pseudoexpressions are valid for their context. If n is a
// pseudoexpression, it will be considered invalid unless its token is in
// this set.
//
// Java: AstValidator.validatePseudoExpression(Node, Token...)
func (v *AstValidator) ValidatePseudoExpression(n rhino.Node, allowedPseudoexpressions ...rhino.Token) {
	// Go-only: language-adapter — ImmutableSet.copyOf(allowedPseudoexpressions)
	// only serves membership tests, which search the varargs slice.
	allowedTokensSet := allowedPseudoexpressions
	switch n.GetToken() {
	case rhino.Token_EMPTY:
		jutil.CheckArgumentMsg1(slices.Contains(allowedTokensSet, rhino.Token_EMPTY), "Unexpected pseudoexpression %s", n)
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_ITER_SPREAD:
		jutil.CheckArgumentMsg1(
			slices.Contains(allowedTokensSet, rhino.Token_ITER_SPREAD), "Unexpected pseudoexpression %s", n)
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		v.ValidateFeature(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS, n)
		v.ValidateExpression(n.GetFirstChild())
	// The only kinds of potential pseudo-expressions we recognize are EMPTY
	// and ITER_SPREAD. So if the given node is neither, validate that it's a
	// (non-pseudo) legitimate expression
	default:
		v.ValidateExpression(n)
	}
}

// ValidateTypeInformation enforces the given node has a type if we are
// validating JSTypes or Colors
//
// n is a Node which we expect to have a type attached (i.e. not a
// control-flow-only node like a BLOCK or IF)
//
// Java: AstValidator.validateTypeInformation(Node)
func (v *AstValidator) ValidateTypeInformation(n rhino.Node) {
	if n.GetIsInClosureUnawareSubtree() {
		// We don't expect closure-unaware code to have type information.
		// TODO: b/321233583 - Maybe this should be a separate validation
		// step, to ensure that nothing tries to infer type information where
		// we are mostly unsure of it?

		return
	}
	if v.typeValidationMode == AstValidator_TypeInfoValidation_NONE {
		return
	}

	if v.typeValidationMode == AstValidator_TypeInfoValidation_JSTYPE {
		type_ := n.GetJSType()

		if type_ != nil && !type_.IsResolved() { // null types are checked in the switch statement
			v.Violation("Found unresolved type "+type_.String(), n)
		}
	}

	switch n.GetToken() {
	case rhino.Token_CALL:
		if !n.GetFirstChild().IsSuper() {
			// TODO(sdh): need to validate super() using validateNewType()
			// instead, if it existed
			v.ValidateCallType(n)
		}
	default:
		v.ExpectSomeTypeInformation(n)
	}
}

// Java: AstValidator.validateCallType(Node)
func (v *AstValidator) ValidateCallType(callNode rhino.Node) {
	switch v.typeValidationMode {
	case AstValidator_TypeInfoValidation_JSTYPE:
		// TODO(b/74537281): Shouldn't CALL nodes always have a type, even if
		// it is unknown?
		callee := callNode.GetFirstChild()
		calleeType := callee.GetJSType()
		if calleeType == nil {
			// Go-only: language-adapter — the message argument
			// callNode.toStringTree() is only built when the check fails.
			jutil.CheckNotNullMsg1(calleeType, "Callee of\n\n%s\nhas no type.", callNode.ToStringTree())
		}

		if calleeType.IsFunctionType() {
			calleeFunctionType := calleeType.ToMaybeFunctionType()
			returnType := calleeFunctionType.GetReturnType()
			// Skip this check if the call node was originally in a cast,
			// because the cast type may be narrower than the return type. Also
			// skip the check if the function's return type is the any
			// (formerly unknown) type, since we may have inferred a better
			// type.
			if callNode.GetJSTypeBeforeCast() == nil && !returnType.IsUnknownType() {
				v.ExpectMatchingTypeInformation(callNode, returnType)
			}
		}
		// TODO(b/74537281): What other cases should be covered?
	case AstValidator_TypeInfoValidation_COLOR:
		callee := callNode.GetFirstChild()
		if callee.GetColor() == nil {
			// Go-only: language-adapter — the message argument
			// callNode.toStringTree() is only built when the check fails.
			jutil.CheckNotNullMsg1(callee.GetColor(), "Callee of\n\n%s\nhas no color.", callNode.ToStringTree())
		}
		// skip additional validation of return types, since optimization
		// colors don't include call signature types
	case AstValidator_TypeInfoValidation_NONE:
		panic(jutil.NewAssertionErrorNoMessage())
	}
}

// Java: AstValidator.expectSomeTypeInformation(Node)
func (v *AstValidator) ExpectSomeTypeInformation(n rhino.Node) {
	switch v.typeValidationMode {
	case AstValidator_TypeInfoValidation_JSTYPE:
		if n.GetJSType() == nil {
			v.Violation(
				"Type information missing"+
					"\n"+
					v.compiler.ToSourceOfNode(NodeUtil_GetEnclosingStatement(n)),
				n)
		}
	case AstValidator_TypeInfoValidation_COLOR:
		if n.GetColor() == nil {
			v.Violation(
				"Color information missing"+
					"\n"+
					v.compiler.ToSourceOfNode(NodeUtil_GetEnclosingStatement(n)),
				n)
		}
	case AstValidator_TypeInfoValidation_NONE:
		panic(jutil.NewAssertionErrorNoMessage())
	}
}

// Java: AstValidator.expectMatchingTypeInformation(Node, JSType)
func (v *AstValidator) ExpectMatchingTypeInformation(n rhino.Node, expectedTypeI rhino.JSType) {
	typeI := n.GetJSType()
	if !jutil.ObjectsEquals(expectedTypeI, typeI) {
		v.Violation(
			"Expected type: "+
				astValidator_getTypeAnnotationString(expectedTypeI)+
				" Actual type: "+
				astValidator_getTypeAnnotationString(typeI),
			n)
	}
}

// Java: AstValidator.getTypeAnnotationString(JSType)
func astValidator_getTypeAnnotationString(typeI rhino.JSType) string {
	if typeI == nil {
		return "NO TYPE INFORMATION"
	} else {
		return "{" + typeI.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) + "}"
	}
}

// Java: AstValidator.validateYield(Node)
func (v *AstValidator) ValidateYield(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_GENERATORS, n)
	v.ValidateNodeType(rhino.Token_YIELD, n)
	v.ValidateProperties(n)
	v.ValidateChildCountIn(n, 0, 1)
	if n.HasChildren() {
		v.ValidateExpression(n.GetFirstChild())
	}
	v.ValidateYieldWithinGeneratorFunction(n)
}

// Java: AstValidator.validateYieldWithinGeneratorFunction(Node)
func (v *AstValidator) ValidateYieldWithinGeneratorFunction(n rhino.Node) {
	parentFunction := NodeUtil_GetEnclosingFunction(n)
	if parentFunction.IsNil() || !parentFunction.IsGeneratorFunction() {
		v.Violation("'yield' expression is not within a generator function", n)
	} else if v.IsInParameterListOfFunction(n, parentFunction) {
		v.Violation("'yield' expression is not allowed in a parameter list", n)
	}
}

// Java: AstValidator.validateAwait(Node)
func (v *AstValidator) ValidateAwait(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, n)
	v.ValidateNodeType(rhino.Token_AWAIT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateAwaitWithinAsyncFunction(n)
}

// Java: AstValidator.AWAIT_NOT_WITHIN_ASYNC_FUNCTION
const astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION = "'await' expression is not within an async function"

// Java: AstValidator.AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST
const astValidator_AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST = "'await' expression is not allowed in a parameter list"

// Java: AstValidator.validateAwaitWithinAsyncFunction(Node)
func (v *AstValidator) ValidateAwaitWithinAsyncFunction(n rhino.Node) {
	parentFunction := NodeUtil_GetEnclosingFunction(n)

	if parentFunction.IsNil() {
		// Top-level await is only allowed in modules.
		if !NodeUtil_GetEnclosingScript(n).GetBooleanProp(rhino.Node_ES6_MODULE) {
			v.Violation(astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION, n)
		}
		return
	}

	if !parentFunction.IsAsyncFunction() {
		v.Violation(astValidator_AWAIT_NOT_WITHIN_ASYNC_FUNCTION, n)
	} else if v.IsInParameterListOfFunction(n, parentFunction) {
		v.Violation(astValidator_AWAIT_NOT_ALLOWED_IN_PARAMETER_LIST, n)
	}
}

// Java: AstValidator.isInParameterListOfFunction(Node, Node)
func (v *AstValidator) IsInParameterListOfFunction(child rhino.Node, functionNode rhino.Node) bool {
	paramList := rhino.CheckNotNullMsg1(functionNode.GetSecondChild(), "%s", functionNode)
	for parent := child.GetParent(); parent != functionNode; parent = parent.GetParent() {
		rhino.CheckNotNullMsg2(parent, "%s not contained in function %s", child, functionNode)
		if parent == paramList {
			return true
		}
	}
	return false
}

// Java: AstValidator.validateImport(Node)
func (v *AstValidator) ValidateImport(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_MODULES, n)
	v.ValidateNodeType(rhino.Token_IMPORT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)

	if n.GetFirstChild().IsName() {
		v.ValidateName(n.GetFirstChild())
	} else {
		v.ValidateNodeType(rhino.Token_EMPTY, n.GetFirstChild())
	}

	secondChild := n.GetSecondChild()
	switch secondChild.GetToken() {
	case rhino.Token_IMPORT_SPECS:
		v.ValidateImportSpecifiers(secondChild)
	case rhino.Token_IMPORT_STAR:
		v.ValidateNonEmptyString(secondChild)
	default:
		v.ValidateNodeType(rhino.Token_EMPTY, secondChild)
	}

	v.ValidateStringLit(n.GetChildAtIndex(2))
}

// Java: AstValidator.validateImportSpecifiers(Node)
func (v *AstValidator) ValidateImportSpecifiers(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_IMPORT_SPECS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateImportSpecifier(c)
	}
}

// Java: AstValidator.validateImportSpecifier(Node)
func (v *AstValidator) ValidateImportSpecifier(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_IMPORT_SPEC, n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateName(c)
	}
}

// Java: AstValidator.validateExport(Node, boolean)
func (v *AstValidator) ValidateExport(n rhino.Node, isAmbient bool) {
	v.ValidateFeature(parser.FeatureSet_Feature_MODULES, n)
	v.ValidateNodeType(rhino.Token_EXPORT, n)
	if n.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) { // export * from "mod"
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 2)
		v.ValidateNodeType(rhino.Token_EMPTY, n.GetFirstChild())
		v.ValidateStringLit(n.GetSecondChild())
	} else if n.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) { // export default foo = 2
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 1)
		v.ValidateExpression(n.GetFirstChild())
	} else {
		v.ValidateProperties(n)
		v.ValidateChildCountIn(n, 1, 2)
		if n.GetFirstChild().IsExportSpecs() {
			v.ValidateExportSpecifiers(n.GetFirstChild())
		} else {
			v.ValidateStatementWithAmbient(n.GetFirstChild(), isAmbient)
		}
		if n.HasTwoChildren() {
			v.ValidateStringLit(n.GetSecondChild())
		}
	}
}

// Java: AstValidator.validateExportSpecifiers(Node)
func (v *AstValidator) ValidateExportSpecifiers(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_EXPORT_SPECS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateExportSpecifier(c)
	}
}

// Java: AstValidator.validateExportSpecifier(Node)
func (v *AstValidator) ValidateExportSpecifier(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_EXPORT_SPEC, n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateName(c)
	}
}

// Java: AstValidator.validateTaggedTemplateLit(Node)
func (v *AstValidator) ValidateTaggedTemplateLit(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_TEMPLATE_LITERALS, n)
	v.ValidateNodeType(rhino.Token_TAGGED_TEMPLATELIT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateTemplateLit(n.GetLastChild())
}

// Java: AstValidator.validateTemplateLit(Node)
func (v *AstValidator) ValidateTemplateLit(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_TEMPLATE_LITERALS, n)
	v.ValidateNodeType(rhino.Token_TEMPLATELIT, n)
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if child.IsTemplateLitString() {
			v.ValidateTemplateLitString(child)
		} else {
			v.ValidateTemplateLitSub(child)
		}
	}
}

// Java: AstValidator.validateTemplateLitString(Node)
func (v *AstValidator) ValidateTemplateLitString(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_TEMPLATELIT_STRING, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	jutil.Try(func() {
		// Validate that getRawString doesn't throw
		n.GetRawString()
	}, func(e jutil.IUnsupportedOperationException) {
		v.Violation("Invalid TEMPLATELIT_STRING node.", n)
	})
}

// Java: AstValidator.validateTemplateLitSub(Node)
func (v *AstValidator) ValidateTemplateLitSub(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_TEMPLATELIT_SUB, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateInterface(Node)
func (v *AstValidator) ValidateInterface(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_INTERFACE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	name := n.GetFirstChild()
	v.ValidateName(name)
	superTypes := name.GetNext()
	if superTypes.IsEmpty() {
		v.ValidateProperties(superTypes)
		v.ValidateChildless(superTypes)
	} else {
		v.ValidateInterfaceExtends(superTypes)
	}
	v.ValidateInterfaceMembers(n.GetLastChild())
}

// Java: AstValidator.validateInterfaceExtends(Node)
func (v *AstValidator) ValidateInterfaceExtends(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_INTERFACE_EXTENDS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateNamedType(c)
	}
}

// Java: AstValidator.validateInterfaceMembers(Node)
func (v *AstValidator) ValidateInterfaceMembers(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_INTERFACE_MEMBERS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateInterfaceMember(c)
	}
}

// Java: AstValidator.validateInterfaceMember(Node)
func (v *AstValidator) ValidateInterfaceMember(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		v.ValidateFunctionSignature(n.GetFirstChild())
	case rhino.Token_MEMBER_VARIABLE_DEF:
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_INDEX_SIGNATURE:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		child := n.GetFirstChild()
		v.ValidateProperties(child)
		v.ValidateChildless(child)
	case rhino.Token_CALL_SIGNATURE:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
	default:
		v.Violation("Interface contained member of invalid type "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateEnum(Node)
func (v *AstValidator) ValidateEnum(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_ENUM, n)
	v.ValidateName(n.GetFirstChild())
	v.ValidateEnumMembers(n.GetLastChild())
}

// Java: AstValidator.validateEnumMembers(Node)
func (v *AstValidator) ValidateEnumMembers(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_ENUM_MEMBERS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateEnumStringKey(c)
	}
}

// Java: AstValidator.validateEnumStringKey(Node)
func (v *AstValidator) ValidateEnumStringKey(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_STRING_KEY, n)
	v.ValidateObjectLiteralKeyName(n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 0)
}

// ValidateClassDeclaration: in a class declaration, unlike a class
// expression, the class name is required.
//
// Java: AstValidator.validateClassDeclaration(Node, boolean)
func (v *AstValidator) ValidateClassDeclaration(n rhino.Node, isAmbient bool) {
	v.ValidateClassHelper(n, isAmbient)
	v.ValidateName(n.GetFirstChild())
}

// Java: AstValidator.validateClass(Node)
func (v *AstValidator) ValidateClass(n rhino.Node) {
	v.ValidateClassHelper(n, false)
}

// Java: AstValidator.validateClassHelper(Node, boolean)
func (v *AstValidator) ValidateClassHelper(n rhino.Node, isAmbient bool) {
	v.ValidateFeature(parser.FeatureSet_Feature_CLASSES, n)
	v.ValidateNodeType(rhino.Token_CLASS, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)

	name := n.GetFirstChild()
	if name.IsEmpty() {
		v.ValidateProperties(name)
		v.ValidateChildless(name)
	} else {
		v.ValidateName(name)
	}

	superClass := name.GetNext()
	if superClass.IsEmpty() {
		v.ValidateProperties(superClass)
		v.ValidateChildless(superClass)
	} else {
		v.ValidateExpression(superClass)
	}

	v.ValidateClassMembers(n.GetLastChild(), isAmbient)
}

// Java: AstValidator.validateClassMembers(Node, boolean)
func (v *AstValidator) ValidateClassMembers(n rhino.Node, isAmbient bool) {
	v.ValidateNodeType(rhino.Token_CLASS_MEMBERS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateClassMember(c, isAmbient)
	}
}

// Java: AstValidator.validateClassMember(Node, boolean)
func (v *AstValidator) ValidateClassMember(n rhino.Node, isAmbient bool) {
	switch n.GetToken() {
	case rhino.Token_MEMBER_FUNCTION_DEF:
		v.ValidateFeature(parser.FeatureSet_Feature_MEMBER_DECLARATIONS, n)
		if n.IsPrivateIdentifier() {
			v.ValidateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
		}
		v.ValidateObjectLiteralKeyName(n)
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		v.ValidateMemberFunction(n, isAmbient)
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		v.ValidateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
		if n.IsPrivateIdentifier() {
			v.ValidateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
		}
		v.ValidateObjectLiteralKeyName(n)
		v.ValidateObjectLitKey(n)
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		v.ValidateMemberFunction(n, isAmbient)
	case rhino.Token_MEMBER_VARIABLE_DEF:
		v.ValidateProperties(n)
		v.ValidateChildless(n)
	case rhino.Token_COMPUTED_PROP:
		v.ValidateComputedPropClassMethod(n)
	case rhino.Token_MEMBER_FIELD_DEF:
		v.ValidateClassField(n)
	case rhino.Token_COMPUTED_FIELD_DEF:
		v.ValidateComputedPropClassField(n)
	case rhino.Token_INDEX_SIGNATURE:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
		child := n.GetFirstChild()
		v.ValidateProperties(child)
		v.ValidateChildless(child)
	case rhino.Token_CALL_SIGNATURE:
		v.ValidateProperties(n)
		v.ValidateChildCount(n)
	case rhino.Token_BLOCK:
		v.ValidateFeature(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK, n)
		v.ValidateBlock(n)
	case rhino.Token_EMPTY:
		// Empty is allowed too.
	default:
		v.Violation("Class contained member of invalid type "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateMemberFunction(Node, boolean)
func (v *AstValidator) ValidateMemberFunction(n rhino.Node, isAmbient bool) {
	function := n.GetFirstChild()
	if isAmbient {
		v.ValidateFunctionSignature(function)
	} else {
		v.ValidateFunctionExpression(function)
	}
}

// Java: AstValidator.validateClassField(Node)
func (v *AstValidator) ValidateClassField(n rhino.Node) {
	if n.IsPrivateIdentifier() {
		v.ValidateFeature(parser.FeatureSet_Feature_PRIVATE_ELEMENTS, n)
	} else {
		v.ValidateFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS, n)
	}
	v.ValidateNonEmptyString(n)
	if n.HasChildren() {
		v.ValidateExpression(n.GetFirstChild())
	}
}

// Java: AstValidator.validateComputedPropClassField(Node)
func (v *AstValidator) ValidateComputedPropClassField(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS, n)
	v.ValidateExpression(n.GetFirstChild())
	if !n.GetSecondChild().IsNil() {
		v.ValidateExpression(n.GetSecondChild())
	}
}

// Java: AstValidator.validateBlock(Node)
func (v *AstValidator) ValidateBlock(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_BLOCK, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateStatement(c)
	}
}

// Java: AstValidator.validateHasSourceName(Node)
func (v *AstValidator) ValidateHasSourceName(n rhino.Node) {
	sourceName := n.GetSourceFileName()
	if sourceName == nil || *sourceName == "" {
		v.Violation("Missing 'source name' annotation.", n)
	}
}

// Java: AstValidator.validateHasInputId(Node)
func (v *AstValidator) ValidateHasInputId(n rhino.Node) {
	inputId := n.GetInputId()
	if inputId == nil {
		v.Violation("Missing 'input id' annotation.", n)
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

// Java: AstValidator.validateNonEmptyString(Node)
func (v *AstValidator) ValidateNonEmptyString(n rhino.Node) {
	if v.ValidateNonNullString(n) && n.GetString() == "" {
		v.Violation("Expected non-empty string.", n)
	}
}

// Java: AstValidator.validateEmptyString(Node)
func (v *AstValidator) ValidateEmptyString(n rhino.Node) {
	if v.ValidateNonNullString(n) && n.GetString() != "" {
		v.Violation("Expected empty string.", n)
	}
}

// Java: AstValidator.validateNonNullString(Node)
func (v *AstValidator) ValidateNonNullString(n rhino.Node) bool {
	valid := true
	jutil.Try(func() {
		// Go-only: language-adapter — `if (n.getString() == null)`: a Go
		// string is never null (nor is a Java StringNode's), so only the
		// exception of a node that is not a string node remains.
		n.GetString()
	}, func(e jutil.IRuntimeException) {
		v.Violation("Expected non-null string.", n)
		valid = false
	})
	return valid
}

// Java: AstValidator.validateName(Node)
func (v *AstValidator) ValidateName(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NAME, n)
	v.ValidateNonEmptyString(n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateTypeInformation(n)
	v.ValidateShadowContentIfPresent(n)
}

// Java: AstValidator.validateOptionalName(Node)
func (v *AstValidator) ValidateOptionalName(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NAME, n)
	v.ValidateNonNullString(n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	// `n.getString() != null && n.getString().isEmpty()`: never null.
	isEmpty := n.GetString() == ""
	if !isEmpty {
		v.ValidateTypeInformation(n)
	}
}

// Java: AstValidator.validateShadowContentIfPresent(Node)
func (v *AstValidator) ValidateShadowContentIfPresent(n rhino.Node) {
	shadow := n.GetClosureUnawareShadow()
	if shadow.IsNil() {
		return
	}
	if !shadow.IsRoot() {
		v.Violation("Shadow reference node is not a ROOT node", shadow)
		return
	}
	shadowScript := shadow.GetFirstChild()
	if shadowScript.IsNil() || !shadowScript.IsScript() {
		v.Violation("Shadow root node's child is not a script node", shadowScript)
		return
	}
	if shadowScript.GetChildCount() != 1 {
		v.Violation("Shadow SCRIPT node child has more than one child", shadowScript)
		return
	}
	exprResult := shadowScript.GetFirstChild()
	if exprResult.IsNil() || !exprResult.IsExprResult() {
		v.Violation("Shadow SCRIPT node child is not an expr result node", exprResult)
		return
	}
	if exprResult.GetChildCount() != 1 {
		v.Violation("Shadow EXPR_RESULT node should have exactly one child", exprResult)
		return
	}
	shadowJsCall := exprResult.GetOnlyChild()
	if !shadowJsCall.IsCall() {
		v.Violation("Shadow node EXPR_RESULT child is not a call", shadowJsCall)
		return
	}
	shadowJsFunction := shadowJsCall.GetLastChild()
	if !shadowJsFunction.IsFunction() {
		v.Violation("Shadow node CALL child is not a function", shadowJsFunction)
		return
	}
	v.ValidateFunctionExpression(shadowJsFunction)
}

// Java: AstValidator.validateEmptyName(Node)
func (v *AstValidator) ValidateEmptyName(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NAME, n)
	v.ValidateEmptyString(n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
}

// Java: AstValidator.validateFunctionStatement(Node)
func (v *AstValidator) ValidateFunctionStatement(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_FUNCTION, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateName(n.GetFirstChild())
	v.ValidateParameters(n.GetSecondChild())
	v.ValidateFunctionBody(n.GetLastChild(), false)
	v.ValidateFunctionFeatures(n)
	if n.GetParent().IsBlock() && !n.GetGrandparent().IsFunction() {
		// e.g. if (true) { function f() {} }
		v.ValidateFeature(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, n)
	}
}

// Java: AstValidator.validateFunctionExpression(Node)
func (v *AstValidator) ValidateFunctionExpression(n rhino.Node) {
	v.ValidateFunctionExpressionHelper(n, false)
}

// Java: AstValidator.validateFunctionSignature(Node)
func (v *AstValidator) ValidateFunctionSignature(n rhino.Node) {
	v.ValidateFunctionExpressionHelper(n, true)
}

// Java: AstValidator.validateFunctionExpressionHelper(Node, boolean)
func (v *AstValidator) ValidateFunctionExpressionHelper(n rhino.Node, isAmbient bool) {
	v.ValidateNodeType(rhino.Token_FUNCTION, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)

	v.ValidateParameters(n.GetSecondChild())

	name := n.GetFirstChild()
	body := n.GetLastChild()
	if n.IsArrowFunction() {
		v.ValidateEmptyName(name)
		if body.IsBlock() {
			v.ValidateBlock(body)
		} else {
			v.ValidateExpression(body)
		}
	} else {
		v.ValidateOptionalName(name)
		v.ValidateFunctionBody(body, isAmbient)
	}
	v.ValidateFunctionFeatures(n)
}

// Java: AstValidator.validateFunctionFeatures(Node)
func (v *AstValidator) ValidateFunctionFeatures(n rhino.Node) {
	if n.IsArrowFunction() {
		v.ValidateFeature(parser.FeatureSet_Feature_ARROW_FUNCTIONS, n)
	}
	if n.IsGeneratorFunction() {
		v.ValidateFeature(parser.FeatureSet_Feature_GENERATORS, n)
	}
	if n.IsAsyncFunction() {
		v.ValidateFeature(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, n)
	}
	if n.IsAsyncFunction() && n.IsGeneratorFunction() {
		v.ValidateFeature(parser.FeatureSet_Feature_ASYNC_GENERATORS, n)
	}
}

// Java: AstValidator.validateFunctionBody(Node, boolean)
func (v *AstValidator) ValidateFunctionBody(n rhino.Node, noBlock bool) {
	if noBlock {
		v.ValidateNodeType(rhino.Token_EMPTY, n)
	} else {
		v.ValidateBlock(n)
	}
}

// Java: AstValidator.validateParameters(Node)
func (v *AstValidator) ValidateParameters(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_PARAM_LIST, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.IsRest() {
			v.ValidateRestParameters(rhino.Token_PARAM_LIST, c)
		} else if c.IsDefaultValue() {
			v.ValidateFeature(parser.FeatureSet_Feature_DEFAULT_PARAMETERS, c)
			v.ValidateDefaultValue(rhino.Token_PARAM_LIST, c)
		} else {
			if c.IsName() {
				v.ValidateName(c)
			} else if c.IsArrayPattern() {
				v.ValidateArrayPattern(rhino.Token_PARAM_LIST, c)
			} else {
				v.ValidateObjectPattern(rhino.Token_PARAM_LIST, c)
			}
		}
	}
}

// Java: AstValidator.validateDefaultValue(Token, Node)
func (v *AstValidator) ValidateDefaultValue(contextType rhino.Token, n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateLHS(contextType, n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateCall(Node)
func (v *AstValidator) ValidateCall(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_CALL, n)
	v.ValidateProperties(n)
	v.ValidateMinimumChildCount(n, 1)
	callee := n.GetFirstChild()
	if callee.IsSuper() {
		v.ValidateSuper(callee)
	} else {
		v.ValidateExpression(callee)
	}
	for c := callee.GetNext(); !c.IsNil(); c = c.GetNext() {
		v.ValidatePseudoExpression(c, rhino.Token_ITER_SPREAD)
	}
}

// Java: AstValidator.validateOptChainCall(Node)
func (v *AstValidator) ValidateOptChainCall(node rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, node)
	v.ValidateNodeType(rhino.Token_OPTCHAIN_CALL, node)
	v.ValidateProperties(node)
	v.ValidateMinimumChildCount(node, 1)
	callee := node.GetFirstChild()
	v.ValidateExpression(callee)
	for argument := callee.GetNext(); !argument.IsNil(); argument = argument.GetNext() {
		v.ValidatePseudoExpression(argument, rhino.Token_ITER_SPREAD)
	}
	v.ValidateFirstNodeOfOptChain(node)
}

// Java: AstValidator.validateSuper(Node)
func (v *AstValidator) ValidateSuper(superNode rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_SUPER, superNode)
	v.ValidateProperties(superNode)
	v.ValidateChildless(superNode)
	v.ValidateTypeInformation(superNode)
	superParent := superNode.GetParent()
	methodNode := NodeUtil_GetEnclosingNonArrowFunction(superParent)

	if NodeUtil_IsNormalGet(superParent) && superNode.IsFirstChildOf(superParent) {
		// `super.prop` or `super['prop']`
		if !v.AllowsSuperPropertyReference(superParent) {
			v.Violation(
				"super property references are only allowed in methods, class static blocks and class"+
					" fields.",
				superNode)
		}
	} else if superParent.IsCall() && superNode.IsFirstChildOf(superParent) {
		// super() constructor call
		if methodNode.IsNil() || !NodeUtil_IsEs6Constructor(methodNode) {
			v.Violation("super constructor call is only allowed in a constructor method", superNode)
		} else {
			extendsNode := methodNode.
				GetParent().     // MEMBER_FUNCTION_DEF
				GetParent().     // CLASS_METHODS
				GetParent().     // CLASS
				GetSecondChild() // extends clause
			if extendsNode.IsEmpty() {
				v.Violation("super constructor call in a class that extends nothing", superNode)
			}
		}
	} else {
		v.Violation("`super` is a syntax error here", superNode)
	}
}

// AllowsSuperPropertyReference checks if a super property reference is in
// a method, class static block or class field.
//
// Java: AstValidator.allowsSuperPropertyReference(Node)
func (v *AstValidator) AllowsSuperPropertyReference(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return false
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		return true
	case rhino.Token_FUNCTION:
		if NodeUtil_IsMethodDeclaration(n) {
			return true
		} else if !n.IsArrowFunction() {
			return false
		}
		return v.AllowsSuperPropertyReference(n.GetParent())
	case rhino.Token_BLOCK:
		if NodeUtil_IsClassStaticBlock(n) {
			return true
		}
		return v.AllowsSuperPropertyReference(n.GetParent())
	default:
		return v.AllowsSuperPropertyReference(n.GetParent())
	}
}

// Java: AstValidator.validateRestParameters(Token, Node)
func (v *AstValidator) ValidateRestParameters(contextType rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_REST_PARAMETERS, n)
	v.ValidateRest(contextType, n)
}

// Java: AstValidator.validateArrayPatternRest(Token, Node)
func (v *AstValidator) ValidateArrayPatternRest(contextType rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_ARRAY_PATTERN_REST, n)
	v.ValidateRest(contextType, n)
}

// Java: AstValidator.validateObjectPatternRest(Token, Node)
func (v *AstValidator) ValidateObjectPatternRest(contextType rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OBJECT_PATTERN_REST, n)
	v.ValidateRest(contextType, n)
}

// ValidateRest: contextType is a Token constant value indicating that n
// should be validated appropriately for a descendant of a Node of this
// type.
//
// Java: AstValidator.validateRest(Token, Node)
func (v *AstValidator) ValidateRest(contextType rhino.Token, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
	default:
		v.Violation("Unexpected node type.", n)
		return
	}
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateLHS(contextType, n.GetFirstChild())
	if !n.GetNext().IsNil() {
		v.Violation("Rest parameters must come after all other parameters.", n)
	}
}

// Java: AstValidator.validateObjectSpread(Node)
func (v *AstValidator) ValidateObjectSpread(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateFeature(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD, n)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateNew(Node)
func (v *AstValidator) ValidateNew(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NEW, n)
	v.ValidateProperties(n)
	v.ValidateMinimumChildCount(n, 1)

	v.ValidateExpression(n.GetFirstChild())
	for c := n.GetSecondChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidatePseudoExpression(c, rhino.Token_ITER_SPREAD)
	}
}

// ValidateNameDeclarationHelper: statement is the enclosing statement. Will
// not always match the declaration Token.
//
// Java: AstValidator.validateNameDeclarationHelper(Node, Token, Node)
func (v *AstValidator) ValidateNameDeclarationHelper(statement rhino.Node, declaration rhino.Token, n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateMinimumChildCount(n, 1)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateNameDeclarationChild(statement, declaration, c)
	}
	if declaration == rhino.Token_LET {
		v.ValidateFeature(parser.FeatureSet_Feature_LET_DECLARATIONS, n)
	} else if declaration == rhino.Token_CONST {
		v.ValidateFeature(parser.FeatureSet_Feature_CONST_DECLARATIONS, n)
	}
}

// Java: AstValidator.validateNameDeclarationChild(Node, Token, Node)
func (v *AstValidator) ValidateNameDeclarationChild(statement rhino.Node, declaration rhino.Token, n rhino.Node) {
	inEnhancedFor := NodeUtil_IsEnhancedFor(statement)
	inForIn := statement.IsForIn()
	var minValues int
	var maxValues int
	if inForIn && declaration == rhino.Token_VAR {
		// ECMASCRIPT5 sloppy mode allows for-in initializers.
		minValues = 0
		maxValues = 1
	} else if inEnhancedFor {
		minValues = 0
		maxValues = 0
	} else if n.IsDestructuringLhs() || declaration == rhino.Token_CONST {
		minValues = 1
		maxValues = 1
	} else {
		minValues = 0
		maxValues = 1
	}

	if n.IsName() {
		// Don't use validateName here since this NAME node may have a child.
		v.ValidateNonEmptyString(n)
		v.ValidateProperties(n)
		v.ValidateChildCountIn(n, minValues, maxValues)

		if n.HasChildren() {
			v.ValidateExpression(n.GetFirstChild())
		}
	} else if n.IsDestructuringLhs() {
		v.ValidateProperties(n)
		v.ValidateChildCountIn(n, 1+minValues, 1+maxValues)

		c := n.GetFirstChild()
		switch c.GetToken() {
		case rhino.Token_ARRAY_PATTERN:
			v.ValidateArrayPattern(declaration, c)
		case rhino.Token_OBJECT_PATTERN:
			v.ValidateObjectPattern(declaration, c)
		default:
			v.Violation("Invalid destructuring lhs first child for "+declaration.String()+" node", n)
		}

		if n.HasTwoChildren() {
			v.ValidateExpression(n.GetSecondChild())
		}
	} else {
		v.Violation("Invalid child for "+declaration.String()+" node", n)
	}
}

// ValidateLHS: contextType is a Token constant value indicating that n
// should be validated appropriately for a descendant of a Node of this
// type.
//
// Java: AstValidator.validateLHS(Token, Node)
func (v *AstValidator) ValidateLHS(contextType rhino.Token, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		v.ValidateName(n)
	case rhino.Token_ARRAY_PATTERN:
		v.ValidateArrayPattern(contextType, n)
	case rhino.Token_OBJECT_PATTERN:
		v.ValidateObjectPattern(contextType, n)
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		v.ValidateGetPropGetElemInLHS(contextType, n)
	case rhino.Token_CAST:
		v.ValidateLHS(contextType, n.GetOnlyChild())
	default:
		v.Violation("Invalid child for "+contextType.String()+" node", n)
	}
}

// Java: AstValidator.validateGetPropGetElemInLHS(Token, Node)
func (v *AstValidator) ValidateGetPropGetElemInLHS(contextType rhino.Token, n rhino.Node) {
	if contextType == rhino.Token_CONST ||
		contextType == rhino.Token_LET ||
		contextType == rhino.Token_VAR ||
		contextType == rhino.Token_PARAM_LIST {
		v.Violation("Invalid child for "+contextType.String()+" node", n)
		return
	}
	switch n.GetToken() {
	case rhino.Token_GETPROP:
		v.ValidateGetProp(n)
	case rhino.Token_GETELEM:
		v.ValidateGetElem(n)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Expected GETPROP or GETELEM but instead got node " + n.GetToken().String()))
	}
}

// Java: AstValidator.validateArrayPattern(Token, Node)
func (v *AstValidator) ValidateArrayPattern(type_ rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, n)
	v.ValidateNodeType(rhino.Token_ARRAY_PATTERN, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_DEFAULT_VALUE:
			v.ValidateDefaultValue(type_, c)
		case rhino.Token_ITER_REST:
			v.ValidateArrayPatternRest(type_, c)
		case rhino.Token_EMPTY:
			v.ValidateProperties(c)
			v.ValidateChildless(c)
		default:
			v.ValidateLHS(type_, c)
		}
	}
}

// Java: AstValidator.validateObjectPattern(Token, Node)
func (v *AstValidator) ValidateObjectPattern(type_ rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, n)
	v.ValidateNodeType(rhino.Token_OBJECT_PATTERN, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_STRING_KEY:
			v.ValidateObjectPatternStringKey(type_, c)
		case rhino.Token_OBJECT_REST:
			v.ValidateObjectPatternRest(type_, c)
		case rhino.Token_COMPUTED_PROP:
			v.ValidateObjectPatternComputedPropKey(type_, c)
		default:
			v.Violation("Invalid object pattern child for "+type_.String()+" node", n)
		}
	}
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

// Java: AstValidator.validateChildless(Node)
func (v *AstValidator) ValidateChildless(n rhino.Node) {
	v.ValidateChildCountWithExpected(n, 0)
}

// Java: AstValidator.validateAssignmentExpression(Node)
func (v *AstValidator) ValidateAssignmentExpression(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateLHS(n.GetToken(), n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateCompoundAssignmentExpression(Node)
func (v *AstValidator) ValidateCompoundAssignmentExpression(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	contextType := n.GetToken()
	lhs := n.GetFirstChild()
	v.ValidateAssignmentOpTarget(lhs, contextType)
	v.ValidateExpression(n.GetLastChild())
}

// ValidateAssignmentOpTarget validates the lhs of a compound assignment op,
// inc, or dec
//
// This check is stricter than validateLhs.
//
// Java: AstValidator.validateAssignmentOpTarget(Node, Token)
func (v *AstValidator) ValidateAssignmentOpTarget(lhs rhino.Node, contextType rhino.Token) {
	switch lhs.GetToken() {
	case rhino.Token_NAME:
		v.ValidateName(lhs)
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		v.ValidateGetPropGetElemInLHS(contextType, lhs)
	case rhino.Token_CAST:
		v.ValidateProperties(lhs)
		v.ValidateChildCountWithExpected(lhs, 1)
		v.ValidateAssignmentOpTarget(lhs.GetFirstChild(), contextType)
	default:
		v.Violation("Invalid child for "+contextType.String()+" node", lhs)
	}
}

// Java: AstValidator.validateGetElem(Node)
func (v *AstValidator) ValidateGetElem(n rhino.Node) {
	jutil.CheckArgumentMsg1(n.IsGetElem(), "%s", n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	v.ValidatePropertyReferenceTarget(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateOptChainGetElem(Node)
func (v *AstValidator) ValidateOptChainGetElem(node rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, node)
	jutil.CheckArgumentMsg1(node.IsOptChainGetElem(), "%s", node)
	v.ValidateProperties(node)
	v.ValidateChildCountWithExpected(node, 2)
	v.ValidateExpression(node.GetFirstChild())
	v.ValidateExpression(node.GetLastChild())
	v.ValidateFirstNodeOfOptChain(node)
}

// Java: AstValidator.validateGetProp(Node)
func (v *AstValidator) ValidateGetProp(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_GETPROP, n)
	v.ValidatePropertyReferenceTarget(n.GetFirstChild())
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateNonEmptyString(n)
}

// Java: AstValidator.validateOptChainGetProp(Node)
func (v *AstValidator) ValidateOptChainGetProp(node rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING, node)
	v.ValidateNodeType(rhino.Token_OPTCHAIN_GETPROP, node)
	v.ValidateExpression(node.GetFirstChild())
	v.ValidateFirstNodeOfOptChain(node)
	v.ValidateProperties(node)
	v.ValidateChildCount(node)
	v.ValidateNonEmptyString(node)
}

// Java: AstValidator.validatePropertyReferenceTarget(Node)
func (v *AstValidator) ValidatePropertyReferenceTarget(objectNode rhino.Node) {
	if objectNode.IsSuper() {
		v.ValidateSuper(objectNode)
	} else {
		v.ValidateExpression(objectNode)
	}
}

// Java: AstValidator.validateRegExpLit(Node)
func (v *AstValidator) ValidateRegExpLit(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_REGEXP_SYNTAX, n)
	v.ValidateNodeType(rhino.Token_REGEXP, n)
	v.ValidateProperties(n)
	v.ValidateChildCountIn(n, 1, 2)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateStringLit(c)
	}
}

// Java: AstValidator.validateStringLit(Node)
func (v *AstValidator) ValidateStringLit(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_STRINGLIT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	jutil.Try(func() {
		// Validate that getString doesn't throw
		n.GetString()
	}, func(e jutil.IUnsupportedOperationException) {
		v.Violation("Invalid STRING node.", n)
	})
}

// Java: AstValidator.validateNumber(Node)
func (v *AstValidator) ValidateNumber(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NUMBER, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	jutil.Try(func() {
		// Validate that getDouble doesn't throw
		n.GetDouble()
	}, func(e jutil.IUnsupportedOperationException) {
		v.Violation("Invalid NUMBER node.", n)
	})
}

// Java: AstValidator.validateBigInt(Node)
func (v *AstValidator) ValidateBigInt(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_BIGINT, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	jutil.Try(func() {
		// Validate that getBigInt doesn't throw
		n.GetBigInt()
	}, func(e jutil.IUnsupportedOperationException) {
		v.Violation("Invalid BIGINT node.", n)
	})
}

// Java: AstValidator.validateArrayLit(Node)
func (v *AstValidator) ValidateArrayLit(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_ARRAYLIT, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		// Array-literals may have empty slots.
		v.ValidatePseudoExpression(c, rhino.Token_EMPTY, rhino.Token_ITER_SPREAD)
		break
	}
}

// Java: AstValidator.validateObjectLit(Node)
func (v *AstValidator) ValidateObjectLit(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_OBJECTLIT, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		v.ValidateObjectLitKey(c)
	}
}

// Java: AstValidator.validateObjectLitKey(Node)
func (v *AstValidator) ValidateObjectLitKey(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETTER_DEF:
		v.ValidateObjectLitGetKey(n)
	case rhino.Token_SETTER_DEF:
		v.ValidateObjectLitSetKey(n)
	case rhino.Token_STRING_KEY:
		v.ValidateObjectLitStringKey(n)
	case rhino.Token_MEMBER_FUNCTION_DEF:
		v.ValidateClassMember(n, false)
		if n.IsStaticMember() {
			v.Violation("Keys in an object literal should not be static.", n)
		}
	case rhino.Token_COMPUTED_PROP:
		v.ValidateObjectLitComputedPropKey(n)
	case rhino.Token_OBJECT_SPREAD:
		v.ValidateObjectSpread(n)
	default:
		v.Violation("Expected object literal key expression but was "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateObjectLitGetKey(Node)
func (v *AstValidator) ValidateObjectLitGetKey(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_GETTER, n)
	v.ValidateNodeType(rhino.Token_GETTER_DEF, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateObjectLiteralKeyName(n)
	function := n.GetFirstChild()
	v.ValidateFunctionExpression(function)
	// objlit get functions must be nameless, and must have zero parameters.
	if function.GetFirstChild().GetString() != "" {
		v.Violation("Expected unnamed function expression.", n)
	}
	functionParams := function.GetSecondChild()
	if functionParams.HasChildren() {
		v.Violation("get methods must not have parameters.", n)
	}
}

// Java: AstValidator.validateObjectLitSetKey(Node)
func (v *AstValidator) ValidateObjectLitSetKey(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_SETTER, n)
	v.ValidateNodeType(rhino.Token_SETTER_DEF, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateObjectLiteralKeyName(n)
	function := n.GetFirstChild()
	v.ValidateFunctionExpression(function)
	// objlit set functions must be nameless, and must have 1 parameter.
	if function.GetFirstChild().GetString() != "" {
		v.Violation("Expected unnamed function expression.", n)
	}
	functionParams := function.GetSecondChild()
	if !functionParams.HasOneChild() {
		v.Violation("set methods must have exactly one parameter.", n)
	}
}

// Java: AstValidator.validateObjectLitStringKey(Node)
func (v *AstValidator) ValidateObjectLitStringKey(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_STRING_KEY, n)
	v.ValidateObjectLiteralKeyName(n)

	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 1)
	v.ValidateExpression(n.GetFirstChild())
	if n.GetBooleanProp(rhino.Node_IS_SHORTHAND_PROPERTY) {
		v.ValidateFeature(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, n)
	}
}

// Java: AstValidator.validateObjectPatternStringKey(Token, Node)
func (v *AstValidator) ValidateObjectPatternStringKey(type_ rhino.Token, n rhino.Node) {
	v.ValidateNodeType(rhino.Token_STRING_KEY, n)
	v.ValidateObjectLiteralKeyName(n)
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 1)

	c := n.GetFirstChild()
	switch c.GetToken() {
	case rhino.Token_DEFAULT_VALUE:
		v.ValidateDefaultValue(type_, c)
	default:
		v.ValidateLHS(type_, c)
	}
}

// Java: AstValidator.validateObjectLitComputedPropKey(Node)
func (v *AstValidator) ValidateObjectLitComputedPropKey(n rhino.Node) {
	// TODO: b/499356003 - clean up this hack. The right fix is to:
	// 1) model computed prop getter/setters with a separate Feature from
	// COMPUTED_PROPERTIES
	// 2) modify the transpilation passes to *not* remove the computed prop
	// getter/setter in object literals feature from the compiler featureset
	// because it can't actually transpile them. (The transpilation passes
	// report a suppressible error instead).
	if !(n.GetIsInClosureUnawareSubtree() && NodeUtil_IsGetOrSetKey(n)) {
		v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	}
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateObjectPatternComputedPropKey(Token, Node)
func (v *AstValidator) ValidateObjectPatternComputedPropKey(type_ rhino.Token, n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateExpression(n.GetFirstChild())
	if n.GetLastChild().IsDefaultValue() {
		v.ValidateDefaultValue(type_, n.GetLastChild())
	} else {
		v.ValidateLHS(n.GetLastChild().GetToken(), n.GetLastChild())
	}
}

// Java: AstValidator.validateComputedPropClassMethod(Node)
func (v *AstValidator) ValidateComputedPropClassMethod(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateExpression(n.GetFirstChild())
	if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_VARIABLE) {
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 1)
	} else {
		v.ValidateProperties(n)
		v.ValidateChildCountWithExpected(n, 2)
		v.ValidateFunctionExpression(n.GetLastChild())
		if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) {
			v.ValidateObjectLitComputedPropGetKey(n)
			v.ValidateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
		} else if n.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER) {
			v.ValidateObjectLitComputedPropSetKey(n)
			v.ValidateFeature(parser.FeatureSet_Feature_CLASS_GETTER_SETTER, n)
		}
	}
}

// Java: AstValidator.validateObjectLitComputedPropGetKey(Node)
func (v *AstValidator) ValidateObjectLitComputedPropGetKey(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateFeature(parser.FeatureSet_Feature_GETTER, n)
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	function := n.GetLastChild()
	v.ValidateFunctionExpression(function)
	// objlit get functions must be nameless, and must have zero parameters.
	if function.GetFirstChild().GetString() != "" {
		v.Violation("Expected unnamed function expression.", n)
	}
	functionParams := function.GetSecondChild()
	if functionParams.HasChildren() {
		v.Violation("get methods must not have parameters.", n)
	}
}

// Java: AstValidator.validateObjectLitComputedPropSetKey(Node)
func (v *AstValidator) ValidateObjectLitComputedPropSetKey(n rhino.Node) {
	v.ValidateFeature(parser.FeatureSet_Feature_COMPUTED_PROPERTIES, n)
	v.ValidateFeature(parser.FeatureSet_Feature_SETTER, n)
	v.ValidateNodeType(rhino.Token_COMPUTED_PROP, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	function := n.GetLastChild()
	v.ValidateFunctionExpression(function)
	// objlit set functions must be nameless, and must have 1 parameter.
	if function.GetFirstChild().GetString() != "" {
		v.Violation("Expected unnamed function expression.", n)
	}
	functionParams := function.GetSecondChild()
	if !functionParams.HasOneChild() {
		v.Violation("set methods must have exactly one parameter.", n)
	}
}

// Java: AstValidator.validateObjectLiteralKeyName(Node)
func (v *AstValidator) ValidateObjectLiteralKeyName(n rhino.Node) {
	if n.IsQuotedStringKey() {
		jutil.Try(func() {
			// Validate that getString doesn't throw
			n.GetString()
		}, func(e jutil.IUnsupportedOperationException) {
			v.Violation("getString failed for"+n.GetToken().String(), n)
		})
	} else {
		v.ValidateNonEmptyString(n)
	}
}

// Java: AstValidator.validateIncDecOp(Node)
func (v *AstValidator) ValidateIncDecOp(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 1)
	v.ValidateAssignmentOpTarget(n.GetFirstChild(), n.GetToken())
}

// Java: AstValidator.validateUnaryOp(Node)
func (v *AstValidator) ValidateUnaryOp(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 1)
	v.ValidateExpression(n.GetFirstChild())
}

// Java: AstValidator.validateBinaryOp(Node)
func (v *AstValidator) ValidateBinaryOp(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 2)
	v.ValidateExpression(n.GetFirstChild())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateTrinaryOp(Node)
func (v *AstValidator) ValidateTrinaryOp(n rhino.Node) {
	v.ValidateProperties(n)
	v.ValidateChildCountWithExpected(n, 3)
	first := n.GetFirstChild()
	v.ValidateExpression(first)
	v.ValidateExpression(first.GetNext())
	v.ValidateExpression(n.GetLastChild())
}

// Java: AstValidator.validateNamedType(Node)
func (v *AstValidator) ValidateNamedType(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_NAMED_TYPE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateName(n.GetFirstChild())
}

// Java: AstValidator.validateTypeAlias(Node)
func (v *AstValidator) ValidateTypeAlias(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_TYPE_ALIAS, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
}

// Java: AstValidator.validateAmbientDeclaration(Node)
func (v *AstValidator) ValidateAmbientDeclaration(n rhino.Node) {
	v.ValidateNodeType(rhino.Token_DECLARE, n)
	v.ValidateAmbientDeclarationHelper(n.GetFirstChild())
}

// Java: AstValidator.validateAmbientDeclarationHelper(Node)
func (v *AstValidator) ValidateAmbientDeclarationHelper(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		v.ValidateNameDeclarationHelper(n.GetParent(), n.GetToken(), n)
	case rhino.Token_FUNCTION:
		v.ValidateFunctionSignature(n)
	case rhino.Token_CLASS:
		v.ValidateClassDeclaration(n, true)
	case rhino.Token_ENUM:
		v.ValidateEnum(n)
	case rhino.Token_NAMESPACE:
		v.ValidateNamespace(n, true)
	case rhino.Token_TYPE_ALIAS:
		v.ValidateTypeAlias(n)
	case rhino.Token_EXPORT:
		v.ValidateExport(n, true)
	default:
	}
}

// Java: AstValidator.validateNamespace(Node, boolean)
func (v *AstValidator) ValidateNamespace(n rhino.Node, isAmbient bool) {
	v.ValidateNodeType(rhino.Token_NAMESPACE, n)
	v.ValidateProperties(n)
	v.ValidateChildCount(n)
	v.ValidateNamespaceName(n.GetFirstChild())
	v.ValidateNamespaceElements(n.GetLastChild(), isAmbient)
}

// Java: AstValidator.validateNamespaceName(Node)
func (v *AstValidator) ValidateNamespaceName(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		v.ValidateName(n)
	case rhino.Token_GETPROP:
		v.ValidateGetProp(n)
	default:
	}
}

// Java: AstValidator.validateNamespaceElements(Node, boolean)
func (v *AstValidator) ValidateNamespaceElements(n rhino.Node, isAmbient bool) {
	v.ValidateNodeType(rhino.Token_NAMESPACE_ELEMENTS, n)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if isAmbient {
			v.ValidateAmbientDeclarationHelper(c)
		} else {
			v.ValidateStatement(c)
		}
	}
}

// Java: AstValidator.violation(String, Node)
func (v *AstValidator) Violation(message string, n rhino.Node) {
	v.violationHandler.HandleViolation(message, n)
}

// ValidateFirstNodeOfOptChain: the first node of an opt chain must be
// marked with Prop.START_OF_OPT_CHAIN
//
// Java: AstValidator.validateFirstNodeOfOptChain(Node)
func (v *AstValidator) ValidateFirstNodeOfOptChain(n rhino.Node) {
	if !NodeUtil_IsOptChainNode(n.GetFirstChild()) {
		// if the first child of an opt chain node is not an opt chain node
		// then it is the start of an opt chain
		if !n.IsOptionalChainStart() {
			v.Violation(
				"Start of optional chain node "+n.GetToken().String()+" is not marked as the start.", n)
		}
	}
}

// Java: AstValidator.validateNodeType(Token, Node)
func (v *AstValidator) ValidateNodeType(type_ rhino.Token, n rhino.Node) {
	if n.GetToken() != type_ {
		v.Violation("Expected "+type_.String()+" but was "+n.GetToken().String(), n)
	}
}

// Java: AstValidator.validateChildCount(Node)
func (v *AstValidator) ValidateChildCount(n rhino.Node) {
	expectedArity := rhino.Token_Arity(n.GetToken())
	if expectedArity != -1 {
		v.ValidateChildCountWithExpected(n, expectedArity)
	}
}

// Java: AstValidator.validateChildCount(Node, int)
func (v *AstValidator) ValidateChildCountWithExpected(n rhino.Node, expected int) {
	count := n.GetChildCount()
	if expected != count {
		v.Violation("Expected "+strconv.Itoa(expected)+" children, but was "+strconv.Itoa(count), n)
	}
}

// Java: AstValidator.validateChildCountIn(Node, int, int)
func (v *AstValidator) ValidateChildCountIn(n rhino.Node, min int, max int) {
	if max == min {
		v.ValidateChildCountWithExpected(n, min)
		return
	}
	count := n.GetChildCount()
	if count < min || count > max {
		v.Violation("Expected child count in ["+strconv.Itoa(min)+", "+strconv.Itoa(max)+"], but was "+strconv.Itoa(count), n)
	}
}

// Java: AstValidator.validateMinimumChildCount(Node, int)
func (v *AstValidator) ValidateMinimumChildCount(n rhino.Node, i int) {
	valid := false
	if i == 1 {
		valid = n.HasChildren()
	} else if i == 2 {
		valid = n.HasMoreThanOneChild()
	} else {
		valid = n.GetChildCount() >= i
	}

	if !valid {
		v.Violation("Expected at least "+strconv.Itoa(i)+" children, but was "+strconv.Itoa(n.GetChildCount()), n)
	}
}

// Java: AstValidator.validateMaximumChildCount(Node, int)
func (v *AstValidator) ValidateMaximumChildCount(n rhino.Node, i int) {
	valid := false
	if i == 1 {
		valid = !n.HasMoreThanOneChild()
	} else if i == -1 {
		valid = true // Varying number of children.
	} else {
		valid = n.GetChildCount() <= i
	}
	if !valid {
		v.Violation("Expected no more than "+strconv.Itoa(i)+" children, but was "+strconv.Itoa(n.GetChildCount()), n)
	}
}

// Java: AstValidator.validateFeature(Feature, Node)
func (v *AstValidator) ValidateFeature(feature parser.FeatureSet_Feature, n rhino.Node) {
	allowbleFeatures := v.compiler.GetAllowableFeatures()
	// Checks that feature present in the AST is recorded in the compiler's
	// featureSet.
	if !n.IsFromExterns() && !allowbleFeatures.Has(feature) {
		// Skip this check for externs because we don't need to complete
		// transpilation on externs, and currently only transpile externs so
		// that we can typecheck ES6+ features in externs.
		v.Violation("AST should not contain "+feature.String(), n)
	}
	// Note: currentScript may be null if someone called validateStatement or
	// validateExpression
	if !v.isScriptFeatureValidationEnabled || v.currentScript.IsNil() {
		return
	}
	scriptFeatures := NodeUtil_GetFeatureSetOfScript(v.currentScript)
	// Checks that feature present in the AST is recorded in the SCRIPT node's
	// featureSet.
	if scriptFeatures == nil || !NodeUtil_GetFeatureSetOfScript(v.currentScript).Has(feature) {
		v.Violation("SCRIPT node should be marked as containing feature "+feature.String(), v.currentScript)
	}
}

// Java: AstValidator.validateProperties(Node)
func (v *AstValidator) ValidateProperties(n rhino.Node) {
	n.ValidateProperties(func(errorMessage string) { v.Violation(errorMessage, n) })
}

// Java: AstValidator.validateRequiredInlinings(Node)
func (v *AstValidator) ValidateRequiredInlinings(n rhino.Node) {
	if !v.shouldValidateRequiredInlinings {
		return
	}
	// Any node with JSDoc that says that it was required to be inlined is an
	// error.
	jsdocInfo := n.GetJSDocInfo()
	if jsdocInfo != nil && jsdocInfo.IsRequireInlining() {
		v.Violation("@requireInlining node failed to be inlined.", n)
	}
}
