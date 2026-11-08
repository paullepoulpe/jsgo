// Port of com/google/javascript/jscomp/AstValidator.java @ closure-compiler v20261104

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

// AstValidator walks the AST and validates that its structure, node properties,
// FeatureSet annotations, and optional type/color information are valid.
//
// Validator methods check node token and child arity first, using typed
// internal/ast views when structure matches and falling back to raw child
// slots when arity or token is malformed so cascading violations and
// NullPointerExceptions match Java's AstValidator on corrupted ASTs.
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

	// Go-only: the violation handler is the default one, which panics, so
	// the scripts may be validated on parallel workers (processParallel); a
	// caller's handler would be called concurrently and out of order.
	defaultViolationHandler bool
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

// NewAstValidatorWithViolationFunc constructs an AstValidator with a violation
// callback function directly without allocating an interface wrapper.
func NewAstValidatorWithViolationFunc(
	compiler AbstractCompiler,
	handler func(message string, n rhino.Node),
	validateScriptFeatures bool,
	shouldValidateRequiredInlinings bool,
) *AstValidator {
	defaultHandler := handler == nil
	if handler == nil {
		handler = astValidator_defaultViolationHandler
	}
	return &AstValidator{
		typeValidationMode:               AstValidator_TypeInfoValidation_NONE,
		compiler:                         compiler,
		violationHandler:                 handler,
		isScriptFeatureValidationEnabled: validateScriptFeatures,
		shouldValidateRequiredInlinings:  shouldValidateRequiredInlinings,
		defaultViolationHandler:          defaultHandler,
	}
}

// Java: AstValidator(AbstractCompiler, ViolationHandler, boolean, boolean)
func NewAstValidatorWithHandler(
	compiler AbstractCompiler,
	handler interface {
		HandleViolation(message string, n rhino.Node)
	},
	validateScriptFeatures bool,
	shouldValidateRequiredInlinings bool,
) *AstValidator {
	var fn func(message string, n rhino.Node)
	if handler != nil {
		fn = handler.HandleViolation
	}
	return NewAstValidatorWithViolationFunc(
		compiler,
		fn,
		validateScriptFeatures,
		shouldValidateRequiredInlinings,
	)
}

// Java: AstValidator(AbstractCompiler)
func NewAstValidator(compiler AbstractCompiler) *AstValidator {
	return NewAstValidatorWithValidateScriptFeatures(compiler, false)
}

// Java: AstValidator(AbstractCompiler, boolean)
func NewAstValidatorWithValidateScriptFeatures(compiler AbstractCompiler, validateScriptFeatures bool) *AstValidator {
	v := NewAstValidatorWithViolationFunc(
		compiler,
		astValidator_defaultViolationHandler,
		validateScriptFeatures,
		compiler.GetOptions().GetShouldValidateRequiredInlinings() == base.Tri_TRUE,
	)
	v.defaultViolationHandler = true
	return v
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
	v.defaultViolationHandler = handler == nil
	if handler == nil {
		handler = astValidator_defaultViolationHandler
	}
	v.violationHandler = handler
	return v
}

// SetTypeValidationMode enables or disables validation of type information.
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
	if v.processParallel(externs, root) {
		return
	}
	if !externs.IsNil() {
		v.ValidateCodeRoot(externs)
		v.checkGrammar(externs)
	}
	if !root.IsNil() {
		v.ValidateCodeRoot(root)
		v.checkGrammar(root)
	}
}

// Go-only: with more than one thread, the scripts of the externs and the
// program are validated on workers. Process validates them one by one and
// ValidateScript keeps nothing across scripts but the current script, so
// each script's validation is independent; the default handler panics at the
// first violation, and the workers' panics are rethrown in script order, so
// the first violation in sequential order is the one that panics. Not with a
// caller's violation handler, with type validation (it reads types), or with
// the grammar check (a walk of each root after its scripts).
func (v *AstValidator) processParallel(externs rhino.Node, root rhino.Node) bool {
	if !v.defaultViolationHandler || v.typeValidationMode != AstValidator_TypeInfoValidation_NONE || astValidatorChecksGrammar {
		return false
	}
	return traverseScriptsParallel(v.compiler, root, parallelScriptOptions[*AstValidator]{
		hasExterns: true,
		externs:    externs,
		initWorker: func(_ int, rc *Compiler) *AstValidator {
			return NewAstValidatorWithViolationFunc(rc, nil, v.isScriptFeatureValidationEnabled, v.shouldValidateRequiredInlinings)
		},
		runScript: func(w *AstValidator, _ int, _, script rhino.Node, _ *Compiler) {
			w.ValidateScript(script)
		},
	})
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

// Java: AstValidator.violation(String, Node)
func (v *AstValidator) Violation(message string, n rhino.Node) {
	v.violationHandler(message, n)
}

// validateTypeInformation enforces that n has a type attached when validating
// JSTypes or Colors.
//
// Java: AstValidator.validateTypeInformation(Node)
func (v *AstValidator) validateTypeInformation(n rhino.Node) {
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
		if typ := n.GetJSType(); typ != nil && !typ.IsResolved() {
			// Null types are checked in the switch below.
			v.Violation("Found unresolved type "+typ.String(), n)
		}
	}

	if call, ok := ast.AsCall(n); ok {
		if !call.Callee().Node().IsSuper() {
			// TODO(sdh): need to validate super() using validateNewType()
			// instead, if it existed.
			v.validateCallType(call)
		}
		return
	}
	v.expectSomeTypeInformation(n)
}

// Java: AstValidator.validateCallType(Node)
func (v *AstValidator) validateCallType(call ast.Call) {
	callNode := call.Node()
	callee := call.Callee().Node()
	switch v.typeValidationMode {
	case AstValidator_TypeInfoValidation_JSTYPE:
		// TODO(b/74537281): Shouldn't CALL nodes always have a type, even if
		// it is unknown?
		calleeType := callee.GetJSType()
		if calleeType == nil {
			// Go-only: language-adapter — callNode.ToStringTree() is only built
			// when the check fails.
			jutil.CheckNotNullMsg1(calleeType, "Callee of\n\n%s\nhas no type.", callNode.ToStringTree())
		}
		if calleeType.IsFunctionType() {
			returnType := calleeType.ToMaybeFunctionType().GetReturnType()
			// Skip this check if the call node was originally in a cast,
			// because the cast type may be narrower than the return type. Also
			// skip the check if the function's return type is the any
			// (formerly unknown) type, since we may have inferred a better
			// type.
			if callNode.GetJSTypeBeforeCast() == nil && !returnType.IsUnknownType() {
				v.expectMatchingTypeInformation(callNode, returnType)
			}
		}
	case AstValidator_TypeInfoValidation_COLOR:
		if callee.GetColor() == nil {
			// Go-only: language-adapter — callNode.ToStringTree() is only built
			// when the check fails.
			jutil.CheckNotNullMsg1(callee.GetColor(), "Callee of\n\n%s\nhas no color.", callNode.ToStringTree())
		}
		// Skip additional validation of return types, since optimization
		// colors don't include call signature types.
	case AstValidator_TypeInfoValidation_NONE:
		panic(jutil.NewAssertionErrorNoMessage())
	}
}

// Java: AstValidator.expectSomeTypeInformation(Node)
func (v *AstValidator) expectSomeTypeInformation(n rhino.Node) {
	switch v.typeValidationMode {
	case AstValidator_TypeInfoValidation_JSTYPE:
		if n.GetJSType() == nil {
			v.Violation(
				"Type information missing\n"+
					v.compiler.ToSourceOfNode(NodeUtil_GetEnclosingStatement(n)),
				n)
		}
	case AstValidator_TypeInfoValidation_COLOR:
		if n.GetColor() == nil {
			v.Violation(
				"Color information missing\n"+
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
func (v *AstValidator) expectMatchingTypeInformation(n rhino.Node, expectedType rhino.JSType) {
	actualType := n.GetJSType()
	if !astValidator_jsTypesEqual(expectedType, actualType) {
		v.Violation(
			"Expected type: "+
				astValidator_getTypeAnnotationString(expectedType)+
				" Actual type: "+
				astValidator_getTypeAnnotationString(actualType),
			n)
	}
}

// Java: AstValidator.getTypeAnnotationString(JSType)
func astValidator_getTypeAnnotationString(typ rhino.JSType) string {
	if typ == nil {
		return "NO TYPE INFORMATION"
	}
	return "{" + typ.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT) + "}"
}

// Java: AstValidator.validateHasSourceName(Node)
func (v *AstValidator) validateHasSourceName(n rhino.Node) {
	if sourceName, ok := n.SourceFileName(); !ok || sourceName == "" {
		v.Violation("Missing 'source name' annotation.", n)
	}
}

// Java: AstValidator.validateHasInputId(Node)
func (v *AstValidator) validateHasInputId(n rhino.Node) {
	if n.GetInputId() == nil {
		v.Violation("Missing 'input id' annotation.", n)
	}
}

// Java: AstValidator.validateNonEmptyString(Node)
func (v *AstValidator) validateNonEmptyString(n rhino.Node) {
	if v.validateNonNullString(n) && n.GetString() == "" {
		v.Violation("Expected non-empty string.", n)
	}
}

// Java: AstValidator.validateEmptyString(Node)
func (v *AstValidator) validateEmptyString(n rhino.Node) {
	if v.validateNonNullString(n) && n.GetString() != "" {
		v.Violation("Expected empty string.", n)
	}
}

// Java: AstValidator.validateNonNullString(Node)
func (v *AstValidator) validateNonNullString(n rhino.Node) bool {
	if !n.IsStringNode() {
		v.Violation("Expected non-null string.", n)
		return false
	}
	return true
}

// Java: AstValidator.validateName(Node)
func (v *AstValidator) validateName(n rhino.Node) {
	v.validateNodeType(rhino.Token_NAME, n)
	v.validateNonEmptyString(n)
	v.validatePropsAndArity(n)
	v.validateTypeInformation(n)
	v.validateShadowContentIfPresent(n)
}

// Java: AstValidator.validateOptionalName(Node)
func (v *AstValidator) validateOptionalName(n rhino.Node) {
	v.validateNodeType(rhino.Token_NAME, n)
	v.validateNonNullString(n)
	v.validatePropsAndArity(n)
	if n.GetString() != "" {
		v.validateTypeInformation(n)
	}
}

// Java: AstValidator.validateEmptyName(Node)
func (v *AstValidator) validateEmptyName(n rhino.Node) {
	v.validateNodeType(rhino.Token_NAME, n)
	v.validateEmptyString(n)
	v.validatePropsAndArity(n)
}

// Java: AstValidator.validateShadowContentIfPresent(Node)
func (v *AstValidator) validateShadowContentIfPresent(n rhino.Node) {
	shadow := n.GetClosureUnawareShadow()
	if shadow.IsNil() {
		return
	}
	root, ok := ast.AsRoot(shadow)
	if !ok {
		v.Violation("Shadow reference node is not a ROOT node", shadow)
		return
	}
	firstRootItem, _ := root.FirstItem()
	shadowScript, ok := ast.AsScript(firstRootItem.Node())
	if !ok {
		v.Violation("Shadow root node's child is not a script node", firstRootItem.Node())
		return
	}
	onlyScriptItem, ok := shadowScript.OnlyItem()
	if !ok {
		v.Violation("Shadow SCRIPT node child has more than one child", shadowScript.Node())
		return
	}
	exprResult, ok := onlyScriptItem.AsExprStmt()
	if !ok {
		v.Violation("Shadow SCRIPT node child is not an expr result node", onlyScriptItem.Node())
		return
	}
	if !exprResult.Node().HasOneChild() {
		v.Violation("Shadow EXPR_RESULT node should have exactly one child", exprResult.Node())
		return
	}
	shadowJsCall, ok := exprResult.Expr().AsCall()
	if !ok {
		v.Violation("Shadow node EXPR_RESULT child is not a call", exprResult.Expr().Node())
		return
	}
	var shadowJsFunction rhino.Node
	if lastArg, ok := shadowJsCall.LastArg(); ok {
		shadowJsFunction = lastArg.Node()
	} else {
		shadowJsFunction = shadowJsCall.Callee().Node()
	}
	if !shadowJsFunction.IsFunction() {
		v.Violation("Shadow node CALL child is not a function", shadowJsFunction)
		return
	}
	v.validateFunctionExpression(shadowJsFunction)
}

// validateFirstNodeOfOptChain verifies that the first node of an optional chain
// is marked with START_OF_OPT_CHAIN.
//
// Java: AstValidator.validateFirstNodeOfOptChain(Node)
func (v *AstValidator) validateFirstNodeOfOptChain(opt ast.OptChain) {
	if _, ok := ast.AsOptChain(nodeUtil_optChainReceiver(opt)); !ok && !opt.IsStart() {
		v.Violation("Start of optional chain node "+opt.Token().String()+" is not marked as the start.", opt.Node())
	}
}

// Java: AstValidator.validateNodeType(Token, Node)
func (v *AstValidator) validateNodeType(expected rhino.Token, n rhino.Node) bool {
	if n.GetToken() != expected {
		v.Violation("Expected "+expected.String()+" but was "+n.GetToken().String(), n)
		return false
	}
	return true
}

// Java: AstValidator.validateChildless(Node)
func (v *AstValidator) validateChildless(n rhino.Node) bool {
	return v.validateChildCountWithExpected(n, 0)
}

// validateChildlessNode validates both properties and zero-child arity.
func (v *AstValidator) validateChildlessNode(n rhino.Node) {
	v.validateProperties(n)
	v.validateChildless(n)
}

// validatePropsAndArity validates node properties and its token's fixed arity.
func (v *AstValidator) validatePropsAndArity(n rhino.Node) bool {
	v.validateProperties(n)
	return v.validateChildCount(n)
}

// validatePropsAndCount validates node properties and an explicit child count.
func (v *AstValidator) validatePropsAndCount(n rhino.Node, expected int) bool {
	v.validateProperties(n)
	return v.validateChildCountWithExpected(n, expected)
}

// Java: AstValidator.validateChildCount(Node)
func (v *AstValidator) validateChildCount(n rhino.Node) bool {
	if expectedArity := rhino.Token_Arity(n.GetToken()); expectedArity != -1 {
		return v.validateChildCountWithExpected(n, expectedArity)
	}
	return true
}

// Java: AstValidator.validateChildCount(Node, int)
func (v *AstValidator) validateChildCountWithExpected(n rhino.Node, expected int) bool {
	count := n.GetChildCount()
	if expected != count {
		v.Violation("Expected "+strconv.Itoa(expected)+" children, but was "+strconv.Itoa(count), n)
		return false
	}
	return true
}

// Java: AstValidator.validateChildCountIn(Node, int, int)
func (v *AstValidator) validateChildCountIn(n rhino.Node, min, max int) bool {
	if max == min {
		return v.validateChildCountWithExpected(n, min)
	}
	count := n.GetChildCount()
	if count < min || count > max {
		v.Violation("Expected child count in ["+strconv.Itoa(min)+", "+strconv.Itoa(max)+"], but was "+strconv.Itoa(count), n)
		return false
	}
	return true
}

// Java: AstValidator.validateMinimumChildCount(Node, int)
func (v *AstValidator) validateMinimumChildCount(n rhino.Node, min int) bool {
	var valid bool
	switch min {
	case 1:
		valid = n.HasChildren()
	case 2:
		valid = n.HasMoreThanOneChild()
	default:
		valid = n.GetChildCount() >= min
	}
	if !valid {
		v.Violation("Expected at least "+strconv.Itoa(min)+" children, but was "+strconv.Itoa(n.GetChildCount()), n)
	}
	return valid
}

// Java: AstValidator.validateMaximumChildCount(Node, int)
func (v *AstValidator) validateMaximumChildCount(n rhino.Node, max int) bool {
	var valid bool
	switch max {
	case 1:
		valid = !n.HasMoreThanOneChild()
	case -1:
		valid = true
	default:
		valid = n.GetChildCount() <= max
	}
	if !valid {
		v.Violation("Expected no more than "+strconv.Itoa(max)+" children, but was "+strconv.Itoa(n.GetChildCount()), n)
	}
	return valid
}

// Java: AstValidator.validateFeature(Feature, Node)
func (v *AstValidator) validateFeature(feature parser.FeatureSet_Feature, n rhino.Node) {
	allowableFeatures := v.compiler.GetAllowableFeatures()
	if !n.IsFromExterns() && !allowableFeatures.Has(feature) {
		// Skip this check for externs because we don't need to complete
		// transpilation on externs, and currently only transpile externs so
		// that we can typecheck ES6+ features in externs.
		v.Violation("AST should not contain "+feature.String(), n)
	}
	// Note: currentScript may be nil if someone called ValidateStatement or
	// ValidateExpression directly.
	if !v.isScriptFeatureValidationEnabled || v.currentScript.IsNil() {
		return
	}
	scriptFeatures := NodeUtil_GetFeatureSetOfScript(v.currentScript)
	if scriptFeatures.IsZero() || !scriptFeatures.Has(feature) {
		v.Violation("SCRIPT node should be marked as containing feature "+feature.String(), v.currentScript)
	}
}

// Java: AstValidator.validateProperties(Node)
func (v *AstValidator) validateProperties(n rhino.Node) {
	n.ValidateProperties(func(errorMessage string) { v.Violation(errorMessage, n) })
}

// Java: AstValidator.validateRequiredInlinings(Node)
func (v *AstValidator) validateRequiredInlinings(n rhino.Node) {
	if !v.shouldValidateRequiredInlinings {
		return
	}
	if jsdocInfo := n.GetJSDocInfo(); jsdocInfo != nil && jsdocInfo.IsRequireInlining() {
		v.Violation("@requireInlining node failed to be inlined.", n)
	}
}
