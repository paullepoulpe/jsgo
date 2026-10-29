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
	"os"
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
	violationHandler func(message string, n rhino.Node)
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

// AstValidator_TypeInfoValidation selects the type information that is
// validated.
type AstValidator_TypeInfoValidation int32

const (
	AstValidator_TypeInfoValidation_JSTYPE AstValidator_TypeInfoValidation = iota
	AstValidator_TypeInfoValidation_COLOR
	AstValidator_TypeInfoValidation_NONE
)

var astValidator_TypeInfoValidation_names = [...]string{"JSTYPE", "COLOR", "NONE"}

// Name is Enum.name().
func (v AstValidator_TypeInfoValidation) Name() string {
	return astValidator_TypeInfoValidation_names[v]
}

// String is Enum.toString().
func (v AstValidator_TypeInfoValidation) String() string { return v.Name() }

// Ordinal is Enum.ordinal().
func (v AstValidator_TypeInfoValidation) Ordinal() int { return int(v) }

// NewAstValidatorWithViolationFunc constructs an AstValidator with a violation
// callback function directly without allocating an interface wrapper.
func NewAstValidatorWithViolationFunc(
	compiler AbstractCompiler,
	handler func(message string, n rhino.Node),
	validateScriptFeatures bool,
	shouldValidateRequiredInlinings bool) *AstValidator {
	if handler == nil {
		handler = astValidator_defaultViolationHandler
	}
	return &AstValidator{
		typeValidationMode:               AstValidator_TypeInfoValidation_NONE,
		compiler:                         compiler,
		violationHandler:                 handler,
		isScriptFeatureValidationEnabled: validateScriptFeatures,
		shouldValidateRequiredInlinings:  shouldValidateRequiredInlinings,
	}
}

// Java: AstValidator(AbstractCompiler, ViolationHandler, boolean, boolean)
func NewAstValidatorWithHandler(
	compiler AbstractCompiler,
	handler interface {
		HandleViolation(message string, n rhino.Node)
	},
	validateScriptFeatures bool,
	shouldValidateRequiredInlinings bool) *AstValidator {
	var fn func(message string, n rhino.Node)
	if handler != nil {
		fn = handler.HandleViolation
	}
	return NewAstValidatorWithViolationFunc(
		compiler,
		fn,
		validateScriptFeatures,
		shouldValidateRequiredInlinings)
}

// Java: AstValidator(AbstractCompiler)
func NewAstValidator(compiler AbstractCompiler) *AstValidator {
	return NewAstValidatorWithValidateScriptFeatures(compiler /* validateScriptFeatures= */, false)
}

// Java: AstValidator(AbstractCompiler, boolean)
func NewAstValidatorWithValidateScriptFeatures(compiler AbstractCompiler, validateScriptFeatures bool) *AstValidator {
	return NewAstValidatorWithViolationFunc(
		compiler,
		astValidator_defaultViolationHandler,
		validateScriptFeatures,
		compiler.GetOptions().GetShouldValidateRequiredInlinings() == base.Tri_TRUE)
}

// Java: AstValidator(AbstractCompiler, boolean).new ViolationHandler() {...}.handleViolation(String, Node)
func astValidator_defaultViolationHandler(message string, n rhino.Node) {
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

// SetViolationHandler sets the violation callback function.
func (v *AstValidator) SetViolationHandler(handler func(message string, n rhino.Node)) *AstValidator {
	if handler == nil {
		handler = astValidator_defaultViolationHandler
	}
	v.violationHandler = handler
	return v
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

// astValidatorChecksGrammar makes the AstValidator also check the trees
// against the grammar of the typed AST view (internal/ast/grammar.txt), so
// that the grammar stays true to the trees the compiler builds. The tests
// turn it on (ast_validator_test.go), and so does JSGO_AST_GRAMMAR_CHECK=1;
// a normal compilation does not pay for the extra walk.
var astValidatorChecksGrammar = os.Getenv("JSGO_AST_GRAMMAR_CHECK") == "1"

// checkGrammar reports the nodes under root that do not follow the grammar
// (see astValidatorChecksGrammar).
//
// Go-only: the typed view is not in Java.
func (v *AstValidator) checkGrammar(root rhino.Node) {
	if !astValidatorChecksGrammar {
		return
	}
	ast.Check(root, func(n rhino.Node, problem string) {
		v.Violation("AST grammar (internal/ast/grammar.txt): "+problem, n)
	})
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

func astValidator_jsTypesEqual(a, b rhino.JSType) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(b)
}

// Java: AstValidator.expectMatchingTypeInformation(Node, JSType)
func (v *AstValidator) ExpectMatchingTypeInformation(n rhino.Node, expectedTypeI rhino.JSType) {
	typeI := n.GetJSType()
	if !astValidator_jsTypesEqual(expectedTypeI, typeI) {
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
	}
	return "{" + typeI.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) + "}"
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
	if !n.IsStringNode() {
		v.Violation("Expected non-null string.", n)
		return false
	}
	return true
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
	// `n.getString() != null && !n.getString().isEmpty()`: never null.
	if n.GetString() != "" {
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

// Java: AstValidator.validateChildless(Node)
func (v *AstValidator) ValidateChildless(n rhino.Node) {
	v.ValidateChildCountWithExpected(n, 0)
}

// Java: AstValidator.violation(String, Node)
func (v *AstValidator) Violation(message string, n rhino.Node) {
	v.violationHandler(message, n)
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
	if scriptFeatures == nil || !scriptFeatures.Has(feature) {
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
