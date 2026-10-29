// Port of com/google/javascript/jscomp/TypeCheck.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast). Visit keeps the
// Java's token switch, because every node also goes through the checks
// before and after it, and hands each visit method the kind it checks
// (ast.Unary, ast.NormalOrOptChainCall, ast.Assignment, ast.Return...).
// The methods read children by role (op.Left(), call.Callee(),
// assign.Target(), class.Extends(), ret.Value()...) instead of by position,
// and optional children comma-ok. The checks themselves work on the JSTypes
// of nodes and report on nodes, so they take rhino.Nodes (x.Node()).
//
// Left on the rhino API: VisitObjectOrClassLiteralKey and the inheritance
// checks it calls, which serve the members of object literals and classes
// alike (two categories, eleven kinds) and read the first child as the key
// of a computed member or the value of the others; the missing-property
// checks (CheckPropertyAccess and below), whose node is a property access or
// a key of an object pattern; CheckSpread, which reads the only child and
// switches on the parent's token; VisitFunction and the constructor and
// interface checks, which read no children.

// TypeCheck checks the types of JS expressions against any declared type
// information.
type TypeCheck struct {
	compiler  AbstractCompiler
	validator *TypeValidator

	reverseInterpreter ReverseAbstractInterpreter

	typeRegistry *rhino.JSTypeRegistry
	topScope     *TypedScope

	scopeCreator *TypedScopeCreator

	reportUnknownTypes bool
	subtypingMode      rhino.JSType_SubtypingMode

	// This may be expensive, so don't emit these warnings if they're
	// explicitly turned off.
	reportMissingProperties bool

	inferJSDocInfo *InferJSDocInfo

	// These fields are used to calculate the percentage of expressions typed.
	typedCount   int
	nullCount    int
	unknownCount int
	inExterns    bool

	badKeyCheckedTypes map[rhino.JSType]struct{}

	// Logs types for @logTypeInCompiler.
	debugTypeLogger TypeCheck_DebugTypeLogger
}

var typeCheckPrototypeID = rhino.IR_Name("prototype").GetStringID()

//
// Internal errors
//

var TypeCheck_UNEXPECTED_TOKEN = DiagnosticType_Error(
	"JSC_INTERNAL_ERROR_UNEXPECTED_TOKEN",
	"Internal Error: TypeCheck doesn''t know how to handle {0}")

//
// User warnings
//

var TypeCheck_DETERMINISTIC_TEST = DiagnosticType_Warning(
	"JSC_DETERMINISTIC_TEST",
	"condition always evaluates to {2}\n"+
		"left : {0}\n"+
		"right: {1}")

var TypeCheck_INEXISTENT_ENUM_ELEMENT = DiagnosticType_Warning(
	"JSC_INEXISTENT_ENUM_ELEMENT", "element {0} does not exist on this enum")

var TypeCheck_INEXISTENT_PROPERTY = DiagnosticType_Warning("JSC_INEXISTENT_PROPERTY", "Property {0} never defined on {1}")

const TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION = "\n\n" +
	"This property is accessed on a \"loose\" type, but is not defined anywhere in the" +
	" program, so it must not exist."

// disabled by default. This one only makes sense if you're using
// well-typed externs.
var TypeCheck_POSSIBLE_INEXISTENT_PROPERTY = DiagnosticType_Disabled(
	"JSC_POSSIBLE_INEXISTENT_PROPERTY",
	"Property {0} never defined on {1}"+
		strings.ReplaceAll(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION, "'", "''"))

var TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION = DiagnosticType_Warning(
	"JSC_INEXISTENT_PROPERTY_WITH_SUGGESTION",
	"Property {0} never defined on {1}. Did you mean {2}?")

var TypeCheck_STRICT_INEXISTENT_PROPERTY = DiagnosticType_Disabled(
	"JSC_STRICT_INEXISTENT_PROPERTY", "Property {0} never defined on {1}")

var TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY = DiagnosticType_Disabled(
	"JSC_STRICT_INEXISTENT_UNION_PROPERTY",
	"Property {0} not defined on all member types of {1}")

var TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION = DiagnosticType_Disabled(
	"JSC_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION",
	"Property {0} never defined on {1}. Did you mean {2}?")

var TypeCheck_NOT_A_CONSTRUCTOR = DiagnosticType_Warning(
	"JSC_NOT_A_CONSTRUCTOR", "cannot instantiate non-constructor, found type: {0}")

var TypeCheck_INSTANTIATE_ABSTRACT_CLASS = DiagnosticType_Warning("JSC_INSTANTIATE_ABSTRACT_CLASS", "cannot instantiate abstract class")

var TypeCheck_BIT_OPERATION = DiagnosticType_Warning(
	"JSC_BAD_TYPE_FOR_BIT_OPERATION", "operator {0} cannot be applied to {1}")

var TypeCheck_UNARY_OPERATION = DiagnosticType_Warning(
	"JSC_BAD_TYPE_FOR_UNARY_OPERATION", "unary operator {0} cannot be applied to {1}")

var TypeCheck_BINARY_OPERATION = DiagnosticType_Warning(
	"JSC_BAD_TYPES_FOR_BINARY_OPERATION", "operator {0} cannot be applied to {1} and {2}")

var TypeCheck_NOT_CALLABLE = DiagnosticType_Warning("JSC_NOT_FUNCTION_TYPE", "{0} expressions are not callable")

var TypeCheck_CONSTRUCTOR_NOT_CALLABLE = DiagnosticType_Warning(
	"JSC_CONSTRUCTOR_NOT_CALLABLE",
	"Constructor {0} should be called with the \"new\" keyword")

var TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE = DiagnosticType_Warning(
	"JSC_ABSTRACT_SUPER_METHOD_NOT_USABLE",
	"Abstract super method {0} cannot be dereferenced")

var TypeCheck_FUNCTION_MASKS_VARIABLE = DiagnosticType_Warning("JSC_FUNCTION_MASKS_VARIABLE", "function {0} masks variable (IE bug)")

var TypeCheck_MULTIPLE_VAR_DEF = DiagnosticType_Warning(
	"JSC_MULTIPLE_VAR_DEF", "declaration of multiple variables with shared type information")

var TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION = DiagnosticType_Warning(
	"JSC_INVALID_INTERFACE_MEMBER_DECLARATION",
	"interface members can only be empty property declarations,"+" empty functions{0}")

var TypeCheck_INTERFACE_METHOD_NOT_EMPTY = DiagnosticType_Warning(
	"JSC_INTERFACE_METHOD_NOT_EMPTY", "interface member functions must have an empty body")

var TypeCheck_CONFLICTING_EXTENDED_TYPE = DiagnosticType_Warning(
	"JSC_CONFLICTING_EXTENDED_TYPE",
	"{1} cannot extend this type; {0}s can only extend {0}s")

var TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS = DiagnosticType_Warning(
	"JSC_ES5_CLASS_EXTENDING_ES6_CLASS", "ES5 class {0} cannot extend ES6 class {1}")

var TypeCheck_DICT_EXTEND_STRUCT_TYPE = DiagnosticType_Warning(
	"JSC_DICT_EXTEND_STRUCT_TYPE", "@dict class {0} cannot extend @struct class {1}")

var TypeCheck_STRUCT_EXTEND_DICT_TYPE = DiagnosticType_Warning(
	"JSC_DICT_EXTEND_STRUCT_TYPE", "@struct class {0} cannot extend @dict class {1}")

var TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS = DiagnosticType_Warning(
	"JSC_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS",
	"Do not use goog.inherits with ES6 classes. Use the ES6 `extends` keyword to inherit"+
		" instead.")

var TypeCheck_INTERFACE_EXTENDS_LOOP = DiagnosticType_Warning("JSC_INTERFACE_EXTENDS_LOOP", "extends loop involving {0}, loop: {1}")

var TypeCheck_CONFLICTING_IMPLEMENTED_TYPE = DiagnosticType_Warning(
	"JSC_CONFLICTING_IMPLEMENTED_TYPE",
	"{0} cannot implement this type; "+
		"an interface can only extend, but not implement interfaces")

var TypeCheck_BAD_IMPLEMENTED_TYPE = DiagnosticType_Warning("JSC_IMPLEMENTS_NON_INTERFACE", "can only implement interfaces")

// disabled by default.
var TypeCheck_HIDDEN_SUPERCLASS_PROPERTY = DiagnosticType_Disabled(
	"JSC_HIDDEN_SUPERCLASS_PROPERTY",
	"property {0} already defined on superclass {1}; use @override to override it")

var TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY = DiagnosticType_Disabled(
	"JSC_PROTOTYPAL_HIDDEN_SUPERCLASS_PROPERTY",
	"property {0} already defined on supertype {1}; use @override to override it")

// disabled by default.
var TypeCheck_HIDDEN_INTERFACE_PROPERTY = DiagnosticType_Disabled(
	"JSC_HIDDEN_INTERFACE_PROPERTY",
	"property {0} already defined on interface {1}; use @override to override it")

var TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH = DiagnosticType_Warning(
	"JSC_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH",
	"mismatch of the {0} property type and the type "+
		"of the property it overrides from supertype {1}\n"+
		"original: {2}\n"+
		"override: {3}")

var TypeCheck_UNKNOWN_OVERRIDE = DiagnosticType_Warning(
	"JSC_UNKNOWN_OVERRIDE", "property {0} not defined on any superclass of {1}")

var TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE = DiagnosticType_Warning(
	"JSC_UNKNOWN_PROTOTYPAL_OVERRIDE", //
	"property {0} not defined on any supertype of {1}")

var TypeCheck_UNKNOWN_EXPR_TYPE = DiagnosticType_Warning(
	"JSC_UNKNOWN_EXPR_TYPE", "could not determine the type of this expression")

var TypeCheck_WRONG_ARGUMENT_COUNT = DiagnosticType_Warning(
	"JSC_WRONG_ARGUMENT_COUNT",
	"Function {0}: called with {1} argument(s). "+
		"Function requires at least {2} argument(s){3}.")

var TypeCheck_ILLEGAL_IMPLICIT_CAST = DiagnosticType_Warning(
	"JSC_ILLEGAL_IMPLICIT_CAST",
	"Illegal annotation on {0}. @implicitCast may only be used in externs.")

var TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE = DiagnosticType_Warning(
	"JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE",
	"Interface {0} has a property {1} with incompatible types in "+
		"its super interfaces {2} and {3}")

var TypeCheck_EXPECTED_THIS_TYPE = DiagnosticType_Warning(
	"JSC_EXPECTED_THIS_TYPE", "\"{0}\" must be called with a \"this\" type")

var TypeCheck_IN_USED_WITH_STRUCT = DiagnosticType_Warning("JSC_IN_USED_WITH_STRUCT", "Cannot use the IN operator with structs")

var TypeCheck_ILLEGAL_PROPERTY_CREATION = DiagnosticType_Warning(
	"JSC_ILLEGAL_PROPERTY_CREATION",
	"Cannot add a property to a struct instance after it is constructed."+
		" (If you already declared the property, make sure to give it a type.)")

var TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE = DiagnosticType_Disabled(
	"JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE",
	"Cannot add a property to an instance of union type.")

var TypeCheck_ILLEGAL_OBJLIT_KEY = DiagnosticType_Warning("JSC_ILLEGAL_OBJLIT_KEY", "Illegal key, the object literal is a {0}")

var TypeCheck_ILLEGAL_CLASS_KEY = DiagnosticType_Warning("JSC_ILLEGAL_CLASS_KEY", "Illegal key, the class is a {0}")

var TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY = DiagnosticType_Warning(
	"JSC_NON_STRINGIFIABLE_OBJECT_KEY",
	"Object type \"{0}\" contains non-stringifiable key and it may lead to an "+
		"error. Please use ES6 Map instead or implement your own Map structure.")

var TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS = DiagnosticType_Warning(
	"JSC_ABSTRACT_METHOD_IN_CONCRETE_CLASS",
	"Abstract methods can only appear in abstract classes. Please declare the class as "+
		"@abstract")

var TypeCheck_CONFLICTING_GETTER_SETTER_TYPE = DiagnosticType_Warning(
	"JSC_CONFLICTING_GETTER_SETTER_TYPE",
	"The types of the getter and setter for property ''{0}'' do not match.\n"+
		"getter type is: {1}\n"+
		"setter type is: {2}")

var TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS = DiagnosticType_Warning(
	"JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS",
	"Cannot @implement the same interface more than once\nRepeated interface: {0}")

var TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE = DiagnosticType_Error(
	"JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE",
	"Should not assign to a property of readonly type ''{0}''")

// If a diagnostic is disabled by default, do not add it in this list
// TODO(dimvar): Either INEXISTENT_PROPERTY shouldn't be here, or we should
// change DiagnosticGroups.setWarningLevel to not accidentally enable it.
var TypeCheck_ALL_DIAGNOSTICS = NewDiagnosticGroup(
	TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS,
	TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE,
	TypeCheck_BAD_IMPLEMENTED_TYPE,
	TypeCheck_BINARY_OPERATION,
	TypeCheck_BIT_OPERATION,
	TypeCheck_CONFLICTING_EXTENDED_TYPE,
	TypeCheck_CONFLICTING_GETTER_SETTER_TYPE,
	TypeCheck_CONFLICTING_IMPLEMENTED_TYPE,
	TypeCheck_CONSTRUCTOR_NOT_CALLABLE,
	TypeCheck_DETERMINISTIC_TEST,
	TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS,
	TypeCheck_EXPECTED_THIS_TYPE,
	TypeCheck_FUNCTION_MASKS_VARIABLE,
	TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH,
	TypeCheck_ILLEGAL_CLASS_KEY,
	TypeCheck_ILLEGAL_IMPLICIT_CAST,
	TypeCheck_ILLEGAL_OBJLIT_KEY,
	TypeCheck_ILLEGAL_PROPERTY_CREATION,
	TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE,
	TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE,
	TypeCheck_INEXISTENT_ENUM_ELEMENT,
	TypeCheck_INEXISTENT_PROPERTY,
	TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION,
	TypeCheck_INSTANTIATE_ABSTRACT_CLASS,
	TypeCheck_INTERFACE_METHOD_NOT_EMPTY,
	TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION,
	TypeCheck_IN_USED_WITH_STRUCT,
	TypeCheck_MULTIPLE_VAR_DEF,
	TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY,
	TypeCheck_NOT_A_CONSTRUCTOR,
	TypeCheck_NOT_CALLABLE,
	TypeCheck_STRUCT_EXTEND_DICT_TYPE,
	TypeCheck_DICT_EXTEND_STRUCT_TYPE,
	TypeCheck_POSSIBLE_INEXISTENT_PROPERTY,
	TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE,
	RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR,
	RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS,
	RhinoErrorReporter_TYPE_PARSE_ERROR,
	RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR,
	TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS,
	TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH,
	TypedScopeCreator_CTOR_INITIALIZER,
	TypedScopeCreator_IFACE_INITIALIZER,
	TypedScopeCreator_LENDS_ON_NON_OBJECT,
	TypedScopeCreator_UNKNOWN_LENDS,
	TypeCheck_UNARY_OPERATION,
	TypeCheck_UNKNOWN_OVERRIDE,
	TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE,
	TypeCheck_WRONG_ARGUMENT_COUNT)

var TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP = NewDiagnosticGroup(TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS)

// TypeCheck_SuggestionPair is a property-name suggestion and its edit
// distance.
type TypeCheck_SuggestionPair struct {
	suggestion string
	distance   int
}

// Java: TypeCheck.SuggestionPair(String, int)
func NewTypeCheck_SuggestionPair(suggestion string, distance int) *TypeCheck_SuggestionPair {
	p := &TypeCheck_SuggestionPair{}
	p.suggestion = suggestion
	p.distance = distance
	return p
}

// Java: TypeCheck(AbstractCompiler, ReverseAbstractInterpreter, JSTypeRegistry, TypedScope, TypedScopeCreator)
func NewTypeCheckWithScope(
	compiler AbstractCompiler,
	reverseInterpreter ReverseAbstractInterpreter,
	typeRegistry *rhino.JSTypeRegistry,
	topScope *TypedScope,
	scopeCreator *TypedScopeCreator) *TypeCheck {
	tc := &TypeCheck{}
	// Field initializers.
	tc.reportUnknownTypes = false
	tc.subtypingMode = rhino.JSType_SubtypingMode_NORMAL
	tc.reportMissingProperties = true
	tc.inferJSDocInfo = nil
	tc.typedCount = 0
	tc.nullCount = 0
	tc.unknownCount = 0
	// Constructor body.
	tc.compiler = compiler
	tc.validator = compiler.GetTypeValidator()
	tc.reverseInterpreter = reverseInterpreter
	tc.typeRegistry = typeRegistry
	tc.topScope = topScope
	tc.scopeCreator = scopeCreator
	tc.inferJSDocInfo = NewInferJSDocInfo(compiler)
	return tc
}

// Java: TypeCheck(AbstractCompiler, ReverseAbstractInterpreter, JSTypeRegistry)
func NewTypeCheck(
	compiler AbstractCompiler,
	reverseInterpreter ReverseAbstractInterpreter,
	typeRegistry *rhino.JSTypeRegistry) *TypeCheck {
	return NewTypeCheckWithScope(compiler, reverseInterpreter, typeRegistry, nil, nil)
}

// ReportMissingProperties turns on the missing property check. Returns this
// for easy chaining.
//
// Java: TypeCheck.reportMissingProperties(boolean)
func (tc *TypeCheck) ReportMissingProperties(report bool) *TypeCheck {
	tc.reportMissingProperties = report
	return tc
}

// ReportUnknownTypes turns on the unknown types check. Returns this for easy
// chaining.
//
// Java: TypeCheck.reportUnknownTypes(boolean)
func (tc *TypeCheck) ReportUnknownTypes(report bool) *TypeCheck {
	tc.reportUnknownTypes = report
	return tc
}

// Process is the main entry point for this phase of processing. This follows
// the pattern for JSCompiler phases.
//
// externsRoot is the root of the externs parse tree; jsRoot the root of the
// input parse tree to be checked.
//
// Java: TypeCheck.process(Node, Node)
func (tc *TypeCheck) Process(externsRoot rhino.Node, jsRoot rhino.Node) {
	jutil.CheckNotNullPtr(tc.scopeCreator)
	jutil.CheckNotNullPtr(tc.topScope)

	externsAndJs := jsRoot.GetParent()
	jutil.CheckState(!externsAndJs.IsNil())
	jutil.CheckState(externsRoot.IsNil() || externsAndJs.HasChild(externsRoot))

	tc.debugTypeLogger.init(tc.compiler)
	defer tc.debugTypeLogger.Close()

	if !externsRoot.IsNil() {
		tc.Check(externsRoot, true)
	}
	tc.Check(jsRoot, false)
}

// ProcessForTesting is the main entry point of this phase for testing code.
//
// externsRoot may be null or a ROOT node. If null the externs are not
// typechecked. Note: the externs node must always exist in the AST, even if
// not typechecked. jsRoot must be a ROOT node and the second child of the
// global ROOT.
//
// Java: TypeCheck.processForTesting(Node, Node)
func (tc *TypeCheck) ProcessForTesting(externsRoot rhino.Node, jsRoot rhino.Node) *TypedScope {
	jutil.CheckState(tc.scopeCreator == nil)
	jutil.CheckState(tc.topScope == nil)

	jutil.CheckArgumentMsg1(externsRoot.IsNil() || externsRoot.IsRoot(), "%s", externsRoot)
	jutil.CheckArgumentMsg1(jsRoot.IsRoot(), "%s", jsRoot)

	jutil.CheckStateMsg1(jsRoot.HasParent() && jsRoot.GetParent().IsRoot(), "%s", jsRoot.GetParent())
	jutil.CheckStateMsg(
		externsRoot.IsNil() || externsRoot.GetNext() == jsRoot,
		"externs root must be the preceding sibling of the js root")

	tc.scopeCreator = NewTypedScopeCreator(tc.compiler)

	tc.topScope =
		NewTypeInferencePass(tc.compiler, tc.reverseInterpreter, tc.scopeCreator).
			InferAllScopes(jsRoot.GetParent())

	tc.Process(externsRoot, jsRoot)

	tc.compiler.SetTypeCheckingHasRun(true)
	tc.compiler.SetTopScope(tc.topScope)

	return tc.topScope
}

// Java: TypeCheck.check(Node, boolean)
func (tc *TypeCheck) Check(node rhino.Node, externs bool) {
	rhino.CheckNotNull(node)
	tc.inExterns = externs
	b := NodeTraversal_Builder{
		compiler:     tc.compiler,
		callback:     tc,
		scopeCreator: tc.scopeCreator,
	}
	b.TraverseWithScope(node, tc.topScope)
	if externs {
		tc.inferJSDocInfo.Process(node, rhino.NilNode)
	} else {
		tc.inferJSDocInfo.Process(rhino.NilNode, node)
	}
}

// Java: TypeCheck.report(Node, DiagnosticType, String...)
func (tc *TypeCheck) Report(n rhino.Node, diagnosticType *DiagnosticType, arguments ...string) {
	tc.compiler.Report(JSError_MakeWithNode(n, diagnosticType, arguments...))
}

// Java: TypeCheck.shouldTraverse(NodeTraversal, Node, Node)
func (tc *TypeCheck) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Start logging types if we're in debug mode and this node was annotated.
	if tc.debugTypeLogger.isLogging {
		tc.debugTypeLogger.MaybeStartLoggingAt(n)
	}

	if n.IsScript() {
		if NodeUtil_IsFromTypeSummary(n) {
			// Errors in type summary files are suppressed, so no use traversing them.
			return false
		}
		if J2clSourceUtils_IsJ2clSourceOfNode(n) {
			tc.subtypingMode = rhino.JSType_SubtypingMode_IGNORE_NULL_UNDEFINED
		} else {
			tc.subtypingMode = rhino.JSType_SubtypingMode_NORMAL
		}
		tc.validator.SetSubtypingMode(tc.subtypingMode)
	} else if fn, ok := ast.AsFunction(n); ok {
		// normal type checking
		// Check for a bug in IE9 (quirks mode) and earlier where bleeding function names would
		// refer to the wrong variable (i.e. "var x; var y = function x() { use(x); }" would pass
		// the 'x' from the outer scope, rather than the local alias bled into the function).
		if fnName, ok := fn.Name(); ok {
			outerScope := t.GetTypedScope()
			var_ := outerScope.GetVar(fnName.Ident())
			if var_ != nil &&
				var_.GetScope().(*TypedScope).HasSameContainerScope(outerScope) &&
				// Ideally, we would want to check whether the type in the scope
				// differs from the type being defined, but then the extern
				// redeclarations of built-in types generates spurious warnings.
				!typeCheck_isFunctionType(var_.GetType()) &&
				!TypeValidator_HasDuplicateDeclarationSuppression(tc.compiler, var_.GetNameNode()) {
				tc.Report(n, TypeCheck_FUNCTION_MASKS_VARIABLE, var_.GetName())
			}
		}

		// TODO(user): Only traverse the function's body. The function's
		// name and arguments are traversed by the scope creator, and ideally
		// should not be traversed by the type checker.
	}
	return true
}

// typeCheck_isFunctionType is Java's `type instanceof FunctionType`.
//
// Go-only: language-adapter — instanceof on an interface type.
func typeCheck_isFunctionType(type_ rhino.JSType) bool {
	_, ok := type_.(rhino.FunctionType)
	return ok
}

// Visit is the meat of the type checking. It is basically one big switch,
// with each case representing one type of parse tree node. The individual
// cases are usually pretty straightforward.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited; parent the parent of the node n.
//
// Java: TypeCheck.visit(NodeTraversal, Node, Node)
func (tc *TypeCheck) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// To be explicitly set to false if the node is not typeable.
	typeable := true

	if nodeType := n.GetJSType(); nodeType != nil {
		if tb := nodeType.GetJSTypeBase(); !tb.IsWellFormedTemplateValid() {
			switch nodeType.GetTypeClass() {
			case rhino.JSTypeClass_ENUM_ELEMENT,
				rhino.JSTypeClass_FUNCTION,
				rhino.JSTypeClass_NAMED,
				rhino.JSTypeClass_UNION,
				rhino.JSTypeClass_TEMPLATIZED,
				rhino.JSTypeClass_TEMPLATE:
				tc.validator.checkWellFormedTemplatizedType(n, nodeType)
			default:
				tb.MarkWellFormedTemplateValidIfResolved()
			}
		}
	}

	switch n.GetToken() {
	case rhino.Token_CAST:
		expr := ast.MustCast(n).Expr().Node()
		exprType := tc.GetJSType(expr)
		castType := tc.GetJSType(n)

		// TODO(johnlenz): determine if we can limit object literals in some
		// way.
		if !expr.IsObjectLit() {
			tc.validator.ExpectCanCast(n, castType, exprType)
		}
		tc.EnsureTypedWithType(n, castType)

		expr.SetJSTypeBeforeCast(exprType)
		if castType.RestrictByNotNullOrUndefined().IsSubtypeOf(exprType) || expr.IsObjectLit() {
			expr.SetJSType(castType)
		}
	case rhino.Token_NAME:
		typeable = tc.VisitName(t, ast.MustName(n), parent)
	case rhino.Token_COMMA:
		tc.EnsureTypedWithType(n, tc.GetJSType(ast.MustBinary(n).Right().Node()))
	case rhino.Token_THIS:
		tc.EnsureTypedWithType(n, t.GetTypedScope().GetTypeOfThis())
	case rhino.Token_NULL:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NULL_TYPE)
	case rhino.Token_NUMBER:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
	case rhino.Token_BIGINT:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BIGINT_TYPE)
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		// Object literal keys are handled with OBJECTLIT
	case rhino.Token_ARRAYLIT:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_ARRAY_TYPE)
	case rhino.Token_REGEXP:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_REGEXP_TYPE)
	case rhino.Token_GETPROP:
		tc.VisitGetProp(t, ast.MustGetProp(n))
		assign, ok := ast.AsAssign(parent)
		typeable = !(ok && assign.Target().Is(n))
	case rhino.Token_OPTCHAIN_GETPROP:
		tc.VisitOptChainGetProp(ast.MustOptChainGetProp(n))
	case rhino.Token_OPTCHAIN_GETELEM:
		tc.VisitOptChainGetElem(ast.MustOptChainGetElem(n))
	case rhino.Token_GETELEM:
		tc.VisitGetElem(ast.MustNormalOrOptChainGetElem(n))
		// The type of GETELEM is always unknown, so no point counting that.
		// If that unknown leaks elsewhere (say by an assignment to another
		// variable), then it will be counted.
		typeable = false
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		tc.VisitVar(t, ast.MustNameDeclaration(n))
		typeable = false
	case rhino.Token_NEW:
		tc.VisitNew(ast.MustNew(n))
	case rhino.Token_OPTCHAIN_CALL:
		// We reuse the `visitCall` functionality for OptChain call nodes because we don't report
		// an error for regular calls when the callee is null or undefined. However, we make sure
		// `typeable` isn't explicitly unset as OptChain nodes are always typed during inference.
		tc.VisitCall(t, ast.MustNormalOrOptChainCall(n))
	case rhino.Token_CALL:
		tc.VisitCall(t, ast.MustNormalOrOptChainCall(n))
		typeable = !parent.IsExprResult()
	case rhino.Token_RETURN:
		tc.VisitReturn(t, ast.MustReturn(n))
		typeable = false
	case rhino.Token_YIELD:
		tc.VisitYield(t, ast.MustYield(n))
	case rhino.Token_DEC, rhino.Token_INC:
		incDec := ast.MustIncDec(n)
		tc.CheckPropCreation(incDec.Target())
		target := incDec.Target().Node()
		if tc.GetJSType(n).IsNumber() {
			tc.validator.ExpectNumber(target, tc.GetJSType(target), "increment/decrement")
			tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
		} else {
			tc.validator.ExpectBigIntOrNumber(target, tc.GetJSType(target), "increment/decrement")
		}
	case rhino.Token_VOID:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_VOID_TYPE)
	case rhino.Token_STRINGLIT, rhino.Token_TYPEOF, rhino.Token_TEMPLATELIT, rhino.Token_TEMPLATELIT_STRING:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_STRING_TYPE)
	case rhino.Token_TAGGED_TEMPLATELIT:
		tc.VisitTaggedTemplateLit(ast.MustTaggedTemplateLit(n))
		tc.EnsureTyped(n)
	case rhino.Token_BITNOT:
		tc.VisitBitwiseNOT(ast.MustUnary(n))
	case rhino.Token_POS:
		tc.VisitUnaryPlus(ast.MustUnary(n))
	case rhino.Token_NEG:
		tc.VisitUnaryMinus(ast.MustUnary(n))
	case rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE:
		op := ast.MustBinary(n)
		left, right := op.Left().Node(), op.Right().Node()

		if left.IsTypeOf() {
			if right.IsStringLit() {
				tc.CheckTypeofString(right, right.GetString())
			}
		} else if right.IsTypeOf() && left.IsStringLit() {
			tc.CheckTypeofString(left, left.GetString())
		}

		leftType := tc.GetJSType(left)
		rightType := tc.GetJSType(right)

		// We do not want to warn about explicit comparisons to VOID. People
		// often do this if they think their type annotations screwed up.
		//
		// We do want to warn about cases where people compare things like
		// (Array|null) == (Function|null)
		// because it probably means they screwed up.
		//
		// This heuristic here is not perfect, but should catch cases we
		// care about without too many false negatives.
		leftTypeRestricted := leftType.RestrictByNotNullOrUndefined()
		rightTypeRestricted := rightType.RestrictByNotNullOrUndefined()

		result := base.Tri_UNKNOWN
		if n.IsEQ() || n.IsNE() {
			result = leftTypeRestricted.TestForEquality(rightTypeRestricted)
			if n.IsNE() {
				result = result.Not()
			}
		} else {
			// SHEQ or SHNE
			if !leftTypeRestricted.CanTestForShallowEqualityWith(rightTypeRestricted) {
				if n.IsSHEQ() {
					result = base.Tri_FALSE
				} else {
					result = base.Tri_TRUE
				}
			}
		}

		if result != base.Tri_UNKNOWN {
			tc.Report(
				n, TypeCheck_DETERMINISTIC_TEST, leftType.String(), rightType.String(), result.String())
		}
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_LT, rhino.Token_LE, rhino.Token_GT, rhino.Token_GE:
		op := ast.MustBinary(n)
		leftSide, rightSide := op.Left().Node(), op.Right().Node()
		leftType := tc.GetJSType(leftSide)
		rightType := tc.GetJSType(rightSide)
		if rightType.IsUnknownType() {
			// validate comparable left
			tc.validator.ExpectUnknownOrComparable(leftSide, leftType, "left side of comparison")
		} else if leftType.IsUnknownType() {
			// validate comparable right
			tc.validator.ExpectUnknownOrComparable(rightSide, rightType, "right side of comparison")
		} else if rightType.IsBigIntOrNumber() {
			// validate left operand for numeric comparison
			tc.validator.ExpectBigIntOrNumber(leftSide, leftType, "left side of numeric comparison")
		} else if leftType.IsBigIntOrNumber() {
			// validate right operand for numeric comparison
			tc.validator.ExpectBigIntOrNumber(rightSide, rightType, "right side of numeric comparison")
		} else {
			errorMsg := "expected matching types in comparison"
			tc.validator.ExpectMatchingTypesStrict(n, leftType, rightType, errorMsg)
			if !leftType.MatchesNumberContext() || !rightType.MatchesNumberContext() {
				// Whether the comparison is numeric will be determined at runtime
				// each time the expression is evaluated. Regardless, both operands
				// should match a string context.
				message := "left side of comparison"
				tc.validator.ExpectString(leftSide, leftType, message)
				tc.validator.ExpectNotNullOrUndefined(
					t, leftSide, leftType, message, tc.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
				message = "right side of comparison"
				tc.validator.ExpectString(rightSide, rightType, message)
				tc.validator.ExpectNotNullOrUndefined(
					t, rightSide, rightType, message, tc.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
			}
		}
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_IN:
		op := ast.MustBinary(n)
		left, right := op.Left().Node(), op.Right().Node()
		rightType := tc.GetJSType(right)
		tc.validator.ExpectStringOrSymbol(left, tc.GetJSType(left), "left side of 'in'")
		tc.validator.ExpectObject(n, rightType, "'in' requires an object")
		if rightType.IsStruct() {
			tc.Report(right, TypeCheck_IN_USED_WITH_STRUCT)
		}
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_INSTANCEOF:
		op := ast.MustBinary(n)
		left, right := op.Left().Node(), op.Right().Node()
		rightType := tc.GetJSType(right).RestrictByNotNullOrUndefined()
		tc.validator.ExpectAnyObject(left, tc.GetJSType(left), "deterministic instanceof yields false")
		tc.validator.ExpectActualObject(right, rightType, "instanceof requires an object")
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_ASSIGN, rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		tc.VisitAssign(t, ast.MustAssignment(n))
		typeable = false
	case rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT:
		assignOp := ast.MustAssignOp(n)
		tc.CheckPropCreation(assignOp.Target())
		tc.VisitBinaryOperator(assignOp.Token(), n, assignOp.Target().Node(), assignOp.Value().Node())
	case rhino.Token_LSH, rhino.Token_RSH, rhino.Token_URSH, rhino.Token_DIV, rhino.Token_MOD, rhino.Token_BITOR,
		rhino.Token_BITXOR, rhino.Token_BITAND, rhino.Token_SUB, rhino.Token_ADD, rhino.Token_MUL, rhino.Token_EXPONENT:
		op := ast.MustBinary(n)
		tc.VisitBinaryOperator(op.Token(), n, op.Left().Node(), op.Right().Node())
	case rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_NOT, rhino.Token_DELPROP:
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_BOOLEAN_TYPE)
	case rhino.Token_CASE:
		// The parent of a case is the body of its switch.
		switchNode := ast.MustSwitch(parent.GetParent())
		switchConditionType := tc.GetJSType(switchNode.Discriminant().Node())
		caseType := tc.GetJSType(ast.MustCase(n).Value().Node())
		tc.validator.ExpectSwitchMatchesCase(n, switchConditionType, caseType)
		typeable = false
	case rhino.Token_WITH:
		child := ast.MustWith(n).Object().Node()
		childType := tc.GetJSType(child)
		tc.validator.ExpectObject(child, childType, "with requires an object")
		typeable = false
	case rhino.Token_FUNCTION:
		tc.VisitFunction(n)
	case rhino.Token_CLASS:
		tc.VisitClass(ast.MustClass(n))
	case rhino.Token_MODULE_BODY:
		tc.VisitModuleBody(t, ast.MustModuleBody(n))
	// These nodes have no interesting type behavior.
	// These nodes require data flow analysis.
	case rhino.Token_PARAM_LIST,
		rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_COMPUTED_PROP,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_COMPUTED_FIELD_DEF,
		rhino.Token_LABEL,
		rhino.Token_LABEL_NAME,
		rhino.Token_SWITCH,
		rhino.Token_SWITCH_BODY,
		rhino.Token_BREAK,
		rhino.Token_CATCH,
		rhino.Token_TRY,
		rhino.Token_SCRIPT,
		rhino.Token_EXPORT,
		rhino.Token_EXPORT_SPEC,
		rhino.Token_EXPORT_SPECS,
		rhino.Token_IMPORT,
		rhino.Token_IMPORT_SPEC,
		rhino.Token_IMPORT_SPECS,
		rhino.Token_IMPORT_STAR,
		rhino.Token_EXPR_RESULT,
		rhino.Token_BLOCK,
		rhino.Token_ROOT,
		rhino.Token_EMPTY,
		rhino.Token_DEFAULT_CASE,
		rhino.Token_CONTINUE,
		rhino.Token_DEBUGGER,
		rhino.Token_THROW,
		rhino.Token_DO,
		rhino.Token_IF,
		rhino.Token_WHILE,
		rhino.Token_FOR,
		rhino.Token_TEMPLATELIT_SUB,
		rhino.Token_ITER_REST,
		rhino.Token_OBJECT_REST,
		rhino.Token_DESTRUCTURING_LHS:
		typeable = false
	case rhino.Token_DYNAMIC_IMPORT:
		tc.VisitDynamicImport(t, ast.MustDynamicImport(n))
	case rhino.Token_ARRAY_PATTERN:
		tc.EnsureTyped(n)
		tc.validator.ExpectAutoboxesToIterable(
			n, tc.GetJSType(n), "array pattern destructuring requires an Iterable")
	case rhino.Token_OBJECT_PATTERN:
		tc.VisitObjectPattern(ast.MustObjectPattern(n))
	case rhino.Token_DEFAULT_VALUE:
		defaultValue := ast.MustDefaultValue(n)
		target, value := defaultValue.Target(), defaultValue.Value().Node()
		tc.CheckCanAssignToWithScope(
			t,
			n,
			target,
			tc.GetJSType(value),
			/* info= */ nil,
			"default value has wrong type")

		// Every other usage of a destructuring pattern is checked while visiting the pattern,
		// but default values are different because they are a conditional assignment and the
		// pattern is not given the default value's type
		switch target.Token() {
		case rhino.Token_ARRAY_PATTERN:
			tc.validator.ExpectAutoboxesToIterable(
				value, tc.GetJSType(value), "array pattern destructuring requires an Iterable")
		case rhino.Token_OBJECT_PATTERN:
			// Verify that the value is not null/undefined, since those can't be destructured.
			tc.validator.ExpectObject(
				value, tc.GetJSType(value), "cannot destructure a 'null' or 'undefined' default value")
		}

		typeable = false
	case rhino.Token_CLASS_MEMBERS:
		classType := parent.GetJSType().ToMaybeFunctionType()
		typ := classType.GetInstanceType()
		for member := range ast.MustClassBody(n).Members() {
			tc.VisitObjectOrClassLiteralKey(member.Node(), parent, typ)
			if getterOrSetter, ok := member.AsGetterOrSetterDef(); ok {
				tc.CheckGetterOrSetterType(getterOrSetter, classType)
			}
		}
		typeable = false
	case rhino.Token_FOR_IN:
		obj := ast.MustForIn(n).Object().Node()
		if tc.GetJSType(obj).IsStruct() {
			tc.Report(obj, TypeCheck_IN_USED_WITH_STRUCT)
		}
		typeable = false
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		tc.EnsureTyped(ast.MustForOfLoop(n).Iterable().Node())
		typeable = false
		// These nodes are typed during the type inference.
	case rhino.Token_SUPER, rhino.Token_NEW_TARGET, rhino.Token_IMPORT_META, rhino.Token_AWAIT, rhino.Token_AND,
		rhino.Token_HOOK, rhino.Token_OR, rhino.Token_COALESCE:
		tc.EnsureTyped(n)
	case rhino.Token_OBJECTLIT:
		// If this is an enum, then give that type to the objectlit as well.
		if _, isEnum := parent.GetJSType().(*rhino.EnumType); isEnum {
			tc.EnsureTypedWithType(n, parent.GetJSType())
		} else {
			tc.EnsureTyped(n)
		}
		typ := tc.GetJSType(n)
		for member := range ast.MustObjectLit(n).Members() {
			tc.VisitObjectOrClassLiteralKey(member.Node(), n, typ)
		}
	case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD:
		tc.CheckSpread(n)
		typeable = false
	default:
		tc.Report(n, TypeCheck_UNEXPECTED_TOKEN, n.GetToken().String())
		tc.EnsureTyped(n)
	}

	if p := n.GetParent(); !p.IsNil() {
		// Visit the body of blockless arrow functions
		if fn, ok := ast.AsFunction(p); ok && fn.Body().Is(n) && !n.IsBlock() {
			tc.VisitImplicitReturnExpression(t, n)
		}

		// Visit the loop initializer of a for-of loop
		// We do this check here, instead of when visiting FOR_OF, in order to get the correct
		// TypedScope.
		if forOf, ok := ast.AsForOfLoop(p); ok && forOf.Target().Is(n) {
			tc.CheckForOfTypes(t, forOf)
		}
	}

	// Don't count externs since the user's code may not even use that part.
	typeable = typeable && !tc.inExterns

	// Record typing logs and metrics.
	if typeable {
		if !tc.debugTypeLogger.parentNode.IsNil() {
			tc.debugTypeLogger.MaybeLogTypeOfNode(n)
		}
		tc.DoPercentTypedAccounting(n)
	}

	tc.CheckJsdocInfoContainsObjectWithBadKey(n)

	// If this is the node which started logging types, stop logging.
	if !tc.debugTypeLogger.parentNode.IsNil() {
		tc.debugTypeLogger.StopLoggingIfThisIsWhereWeStarted(n)
	}
}

// DoPercentTypedAccounting counts the given node in the typed statistics.
//
// n is a node that should be typed.
//
// Java: TypeCheck.doPercentTypedAccounting(Node)
func (tc *TypeCheck) DoPercentTypedAccounting(n rhino.Node) {
	type_ := n.GetJSType()
	if type_ == nil {
		tc.nullCount++
	} else if type_.IsUnknownType() {
		if tc.reportUnknownTypes && !n.GetParent().IsExprResult() {
			tc.compiler.Report(JSError_MakeWithNode(n, TypeCheck_UNKNOWN_EXPR_TYPE))
		}
		tc.unknownCount++
	} else {
		tc.typedCount++
	}
}

// GetJSType gets the JSType from the Node argument and verifies that it is
// present.
//
// Java: TypeCheck.getJSType(Node)
func (tc *TypeCheck) GetJSType(n rhino.Node) rhino.JSType {
	jsType := n.GetJSType()
	if jsType == nil {
		// TODO(nicksantos): This branch indicates a compiler bug, not worthy of
		// halting the compilation but we should log this and analyze to track
		// down why it happens. This is not critical and will be resolved over
		// time as the type checker is extended.
		return tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	} else {
		return jsType
	}
}

// GetPropertyTypeIfDeclared returns the type of the property with the given
// name if declared. Otherwise returns unknown.
//
// Java: TypeCheck.getPropertyTypeIfDeclared(ObjectType, Property.Key)
func (tc *TypeCheck) GetPropertyTypeIfDeclared(
	objectType rhino.ObjectType, propertyName rhino.Property_Key) rhino.JSType {
	if objectType != nil &&
		objectType.HasPropertyByKey(propertyName) &&
		!objectType.IsPropertyTypeInferredByKey(propertyName) {
		return objectType.GetPropertyTypeByKey(propertyName)
	}
	return tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// TODO(nicksantos): TypeCheck should never be attaching types to nodes.
// All types should be attached by TypeInference. This is not true today
// for legacy reasons. There are a number of places where TypeInference
// doesn't attach a type, as a signal to TypeCheck that it needs to check
// that node's type.

// EnsureTyped ensures that the given node has a type. If it does not have
// one, attach the UNKNOWN_TYPE.
//
// Java: TypeCheck.ensureTyped(Node)
func (tc *TypeCheck) EnsureTyped(n rhino.Node) {
	tc.EnsureTypedWithType(n, tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
}

// Java: TypeCheck.ensureTyped(Node, JSTypeNative)
func (tc *TypeCheck) EnsureTypedWithNative(n rhino.Node, type_ rhino.JSTypeNative) {
	tc.EnsureTypedWithType(n, tc.GetNativeType(type_))
}

// EnsureTypedWithType ensures the node is typed.
//
// n is the node getting a type assigned to it; type_ the type to be
// assigned.
//
// Java: TypeCheck.ensureTyped(Node, JSType)
func (tc *TypeCheck) EnsureTypedWithType(n rhino.Node, type_ rhino.JSType) {
	// Make sure FUNCTION nodes always get function type.
	jutil.CheckState(!n.IsFunction() || type_.IsFunctionType() || type_.IsUnknownType())
	if n.GetJSType() == nil {
		n.SetJSType(type_)
	}
}

// GetTypedPercent returns the percentage of nodes typed by the type checker.
//
// Returns a number between 0.0 and 100.0.
//
// Java: TypeCheck.getTypedPercent()
func (tc *TypeCheck) GetTypedPercent() float64 {
	total := tc.nullCount + tc.unknownCount + tc.typedCount
	if total == 0 {
		return 0.0
	}
	return (100.0 * float64(tc.typedCount)) / float64(total)
}

// Java: TypeCheck.getNativeType(JSTypeNative)
func (tc *TypeCheck) GetNativeType(typeId rhino.JSTypeNative) rhino.JSType {
	return tc.typeRegistry.GetNativeType(typeId)
}

// Java: TypeCheck.maybeReferencedType(JSType)
func TypeCheck_MaybeReferencedType(type_ rhino.JSType) rhino.JSType {
	maybeTemplatized := type_.ToMaybeTemplatizedType()
	if maybeTemplatized != nil {
		return maybeTemplatized.GetReferencedType()
	}
	return nil
}

// Java: TypeCheck.flattenUnion(JSType)
func TypeCheck_FlattenUnion(maybeUnion rhino.JSType) []rhino.JSType {
	union := maybeUnion.ToMaybeUnionType()
	if union == nil {
		return []rhino.JSType{maybeUnion}
	}
	return union.GetAlternates()
}

// Java: TypeCheck.declaresOverride(JSDocInfo)
func TypeCheck_DeclaresOverride(jsdoc *rhino.JSDocInfo) bool {
	return (jsdoc != nil) && jsdoc.IsOverride()
}

// TypeCheck_DebugTypeLogger logs types for @logTypeInCompiler.
type TypeCheck_DebugTypeLogger struct {
	// Output log file for @logTypeInCompiler.
	typeLogFile diagnostic.LogFile

	// The node where we started logging, if we're logging.
	parentNode rhino.Node

	isLogging bool
}

// Suppress MustBeClosedChecker because `.close` is called by our `close` method.
//
// Java: TypeCheck.DebugTypeLogger()
func NewTypeCheck_DebugTypeLogger(outer *TypeCheck) *TypeCheck_DebugTypeLogger {
	l := &TypeCheck_DebugTypeLogger{}
	l.init(outer.compiler)
	return l
}

func (l *TypeCheck_DebugTypeLogger) init(compiler AbstractCompiler) {
	*l = TypeCheck_DebugTypeLogger{}
	if compiler.IsDebugLoggingEnabled() {
		l.typeLogFile = compiler.CreateOrReopenLog("DebugTypeLogger", "types_logged_in_compiler.log")
		l.isLogging = l.typeLogFile.IsLogging()
	}
}

// Java: TypeCheck.DebugTypeLogger.maybeStartLoggingAt(Node)
func (l *TypeCheck_DebugTypeLogger) MaybeStartLoggingAt(n rhino.Node) {
	// We don't log outside debug mode, or if we're already logging.
	if !l.isLogging || !l.parentNode.IsNil() {
		return
	}

	// We only start logging on @logTypeInCompiler.
	info := n.GetJSDocInfo()
	if info != nil && info.GetLogTypeInCompiler() {
		l.parentNode = n
	}
}

// Java: TypeCheck.DebugTypeLogger.stopLoggingIfThisIsWhereWeStarted(Node)
func (l *TypeCheck_DebugTypeLogger) StopLoggingIfThisIsWhereWeStarted(n rhino.Node) {
	if n == l.parentNode {
		l.parentNode = rhino.NilNode
	}
}

// Java: TypeCheck.DebugTypeLogger.close()
func (l *TypeCheck_DebugTypeLogger) Close() {
	l.typeLogFile.Close()
	l.isLogging = false
	l.parentNode = rhino.NilNode
}

// Java: TypeCheck.DebugTypeLogger.maybeLogTypeOfNode(Node)
func (l *TypeCheck_DebugTypeLogger) MaybeLogTypeOfNode(n rhino.Node) {
	// If we have no parent, we're not logging.
	if l.parentNode.IsNil() {
		return
	}

	if n.GetJSType() == nil {
		return
	}

	l.typeLogFile.LogString(
		n.ToStringWithFlags( /* printSource= */ true /* printAnnotations= */, true /* printType= */, true))
}

func stringPtrOrNull(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

func jsTypeStringOrNull(t rhino.JSType) string {
	if t == nil {
		return "null"
	}
	return t.String()
}

func nodeStringOrNull(n rhino.Node) string {
	if n.IsNil() {
		return "null"
	}
	return n.String()
}
