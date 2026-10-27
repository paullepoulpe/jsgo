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
	"math"
	"strconv"
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

	badKeyCheckedTypes *jutil.ObjectLinkedHashSet[rhino.JSType]

	// Logs types for @logTypeInCompiler.
	debugTypeLogger *TypeCheck_DebugTypeLogger
}

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
	tc.debugTypeLogger = nil
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
	jutil.CheckNotNull(tc.scopeCreator)
	jutil.CheckNotNull(tc.topScope)

	externsAndJs := jsRoot.GetParent()
	jutil.CheckState(!externsAndJs.IsNil())
	jutil.CheckState(externsRoot.IsNil() || externsAndJs.HasChild(externsRoot))

	func() {
		debugTypeLogger := NewTypeCheck_DebugTypeLogger(tc)
		defer debugTypeLogger.Close()
		tc.debugTypeLogger = debugTypeLogger

		if !externsRoot.IsNil() {
			tc.Check(externsRoot, true)
		}
		tc.Check(jsRoot, false)
	}()
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
	tc.debugTypeLogger.MaybeStartLoggingAt(n)

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
		switch nodeType.GetTypeClass() {
		case rhino.JSTypeClass_ENUM_ELEMENT,
			rhino.JSTypeClass_FUNCTION,
			rhino.JSTypeClass_NAMED,
			rhino.JSTypeClass_UNION,
			rhino.JSTypeClass_TEMPLATIZED,
			rhino.JSTypeClass_TEMPLATE:
			tc.validator.ExpectWellFormedTemplatizedType(n)
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
		tc.debugTypeLogger.MaybeLogTypeOfNode(n)
		tc.DoPercentTypedAccounting(n)
	}

	tc.CheckJsdocInfoContainsObjectWithBadKey(n)

	// If this is the node which started logging types, stop logging.
	tc.debugTypeLogger.StopLoggingIfThisIsWhereWeStarted(n)
}

// Java: TypeCheck.visitUnaryPlus(Node)
func (tc *TypeCheck) VisitUnaryPlus(op ast.Unary) {
	n := op.Node()
	if tc.GetJSType(n).IsNoType() {
		tc.Report(
			n,
			TypeCheck_UNARY_OPERATION,
			jutil.ToJavaString(NodeUtil_OpToStr(op.Token())),
			tc.GetJSType(op.Operand().Node()).String())
	} else {
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
	}
}

// Java: TypeCheck.visitUnaryMinus(Node)
func (tc *TypeCheck) VisitUnaryMinus(op ast.Unary) {
	n := op.Node()
	operatorType := tc.GetJSType(n)
	childType := tc.GetJSType(op.Operand().Node())
	if operatorType.IsNumber() {
		// This condition is meant to catch any old cases (where bigint isn't involved)
		tc.validator.ExpectNumber(n, childType, "sign operator")
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
	} else {
		tc.validator.ExpectBigIntOrNumber(n, childType, "unary minus operator")
	}
}

// Java: TypeCheck.visitBitwiseNOT(Node)
func (tc *TypeCheck) VisitBitwiseNOT(op ast.Unary) {
	n := op.Node()
	operatorType := tc.GetJSType(n)
	childType := tc.GetJSType(op.Operand().Node())
	if operatorType.IsNumber() {
		// This condition is meant to catch any old cases (where bigint isn't involved)
		if !childType.MatchesNumberContext() {
			tc.Report(n, TypeCheck_BIT_OPERATION, jutil.ToJavaString(NodeUtil_OpToStr(op.Token())), childType.String())
		} else {
			tc.validator.ExpectNumberStrict(n, childType, "bitwise NOT")
		}
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
	} else {
		tc.validator.ExpectBigIntOrNumber(n, childType, "bitwise NOT")
	}
}

// Java: TypeCheck.checkSpread(Node)
func (tc *TypeCheck) CheckSpread(spreadNode rhino.Node) {
	target := spreadNode.GetOnlyChild()
	tc.EnsureTyped(target)
	targetType := tc.GetJSType(target)

	switch spreadNode.GetParent().GetToken() {
	case rhino.Token_OBJECTLIT:
		// Case: `var x = {A: a, B: b, ...obj}`.
		// Nothing to check about object spread.
	case rhino.Token_ARRAYLIT, rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL, rhino.Token_NEW:
		// Case: `var x = [a, b, ...itr]`
		// Case: `var x = fn(a, b, ...itr)`
		// Case: `var x = fn?.(a, b, ...itr)`
		tc.validator.ExpectAutoboxesToIterable(
			target, targetType, "Spread operator only applies to Iterable types")
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected parent of SPREAD: " + spreadNode.GetParent().ToStringTree()))
	}
}

// Java: TypeCheck.checkTypeofString(Node, String)
func (tc *TypeCheck) CheckTypeofString(n rhino.Node, s string) {
	if !(s == "number" ||
		s == "string" ||
		s == "boolean" ||
		s == "undefined" ||
		s == "function" ||
		s == "object" ||
		s == "symbol" ||
		s == "unknown" ||
		s == "bigint") {
		tc.validator.ExpectValidTypeofName(n, s)
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

// CheckGetterOrSetterType verifies that a user did not give a getter and
// setter different types, as the type system will arbitrarily take the
// first of the types if different.
//
// Java: TypeCheck.checkGetterOrSetterType(Node, FunctionType)
func (tc *TypeCheck) CheckGetterOrSetterType(child ast.GetterOrSetterDef, classType rhino.FunctionType) {
	propertyName := child.Key()

	methodType := child.Function().Node().GetJSType().ToMaybeFunctionType()
	_, isGetter := child.AsGetterDef()
	var propertyType rhino.JSType
	if isGetter {
		propertyType = tc.DetermineGetterType(methodType)
	} else {
		propertyType = typeCheck_iterablesGetOnlyElement(methodType.GetParameters()).GetJSType()
	}
	var officialPropertyType rhino.JSType
	if child.Node().IsStaticMember() {
		officialPropertyType = classType.GetPropertyType(propertyName)
	} else {
		officialPropertyType = classType.GetPrototype().GetPropertyType(propertyName)
	}
	if !propertyType.Equals(officialPropertyType) {
		// TODO(b/116797078): make this not an error - instead, store the getter and setter types
		// separately
		var getterType, setterType string
		if isGetter {
			getterType = propertyType.String()
			setterType = officialPropertyType.String()
		} else {
			getterType = officialPropertyType.String()
			setterType = propertyType.String()
		}
		tc.Report(
			child.Node(),
			TypeCheck_CONFLICTING_GETTER_SETTER_TYPE,
			propertyName,
			getterType,
			setterType)
	}
}

// typeCheck_iterablesGetOnlyElement is Guava's
// Iterables.getOnlyElement(Iterable): NoSuchElementException when empty,
// IllegalArgumentException listing up to five elements when longer.
//
// Go-only: language-adapter — Guava helper.
func typeCheck_iterablesGetOnlyElement(params []*rhino.FunctionType_Parameter) *rhino.FunctionType_Parameter {
	if len(params) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if len(params) == 1 {
		return params[0]
	}
	var sb strings.Builder
	sb.WriteString("expected one element but was: <")
	sb.WriteString(params[0].String())
	i := 1
	for ; i < 5 && i < len(params); i++ {
		sb.WriteString(", ")
		sb.WriteString(params[i].String())
	}
	if i < len(params) {
		sb.WriteString(", ...")
	}
	sb.WriteByte('>')
	panic(jutil.NewIllegalArgumentExceptionMsg(sb.String()))
}

// Java: TypeCheck.determineGetterType(FunctionType)
func (tc *TypeCheck) DetermineGetterType(methodType rhino.FunctionType) rhino.JSType {
	// TODO(sdh): consider only falling back on unknown if the function body is empty?  But we
	// need to not report a conflicting type error if there's different unknowns.
	if !methodType.IsReturnTypeInferred() {
		return methodType.GetReturnType()
	}
	return tc.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// VisitAssign visits an assignment `lvalue = rvalue`. If the `lvalue` is a
// prototype modification, we change the schema of the object type it is
// referring to.
//
// t is the traversal; assign the assign node (an `=`, `||=`, `&&=` or
// `??=`).
//
// Java: TypeCheck.visitAssign(NodeTraversal, Node)
func (tc *TypeCheck) VisitAssign(t *NodeTraversal, assign ast.Assignment) {
	n := assign.Node()
	info := n.GetJSDocInfo()
	lvalue := assign.Target()
	rvalue := assign.Value()

	rightType := tc.GetJSType(rvalue.Node())
	tc.CheckCanAssignToWithScope(t, n, lvalue, rightType, info, "assignment")
	tc.EnsureTypedWithType(n, rightType)
}

// CheckCanAssignToWithScope checks that we can assign the given right type to
// the given lvalue or destructuring pattern.
//
// See CheckCanAssignToNameGetpropOrGetelem for more details.
//
// nodeToWarn is a node to report type mismatch warnings on; lvalue the lvalue
// to which we're assigning or a destructuring pattern; rightType the type
// we're assigning to the lvalue; msg a message to report along with any type
// mismatch warnings.
//
// Java: TypeCheck.checkCanAssignToWithScope(NodeTraversal, Node, Node, JSType, JSDocInfo, String)
func (tc *TypeCheck) CheckCanAssignToWithScope(
	t *NodeTraversal,
	nodeToWarn rhino.Node,
	lvalue ast.Target,
	rightType rhino.JSType,
	info *rhino.JSDocInfo,
	msg string) {
	if pattern, ok := lvalue.AsPattern(); ok {
		tc.CheckDestructuringAssignment(t, nodeToWarn, pattern, rightType, msg)
	} else {
		tc.CheckCanAssignToNameGetpropOrGetelem(t, nodeToWarn, lvalue, rightType, info, msg)
	}
}

// CheckDestructuringAssignment recursively checks that an assignment to a
// destructuring pattern is valid for all the lvalues contained in the pattern
// (including in nested patterns).
//
// Java: TypeCheck.checkDestructuringAssignment(NodeTraversal, Node, Node, JSType, String)
func (tc *TypeCheck) CheckDestructuringAssignment(
	t *NodeTraversal, nodeToWarn rhino.Node, pattern ast.Pattern, rightType rhino.JSType, msg string) {
	for _, target := range DestructuredTarget_CreateAllNonEmptyTargetsInPattern(tc.typeRegistry, rightType, pattern.Node()) {

		// TODO(b/77597706): this is not very efficient because it re-infers the types below,
		// which we already did once in TypeInference. don't repeat the work.
		tc.CheckCanAssignToWithScope(
			t, nodeToWarn, ast.MustTarget(target.GetNode()), target.InferType() /* info= */, nil, msg)
	}
}

// CheckCanAssignToNameGetpropOrGetelem checks that we can assign the given
// right type to the given lvalue.
//
// If the lvalue is a qualified name, and has a declared type in the given
// scope, uses the declared type of the qualified name instead of the type on
// the node.
//
// nodeToWarn is a node to report type mismatch warnings on; lvalue the lvalue
// to which we're assigning - a NAME, GETELEM, or GETPROP; rightType the type
// we're assigning to the lvalue; msg a message to report along with any type
// mismatch warnings.
//
// Java: TypeCheck.checkCanAssignToNameGetpropOrGetelem(NodeTraversal, Node, Node, JSType, JSDocInfo, String)
func (tc *TypeCheck) CheckCanAssignToNameGetpropOrGetelem(
	t *NodeTraversal, nodeToWarn rhino.Node, lvalue ast.Target, rightType rhino.JSType, info *rhino.JSDocInfo, msg string) {
	_, isPattern := lvalue.AsPattern()
	jutil.CheckArgumentMsg1(!isPattern, "%s", lvalue.Node())

	// Ensure our LHS is not readonly.
	tc.CheckNotReadonlyPropertyAssignment(lvalue)

	if getProp, ok := lvalue.AsGetProp(); ok {
		object := getProp.Object()
		objectJsType := tc.GetJSType(object.Node())
		pname := getProp.Property()

		// the first name in this getprop refers to an interface
		// we perform checks in addition to the ones below
		if objectGetProp, ok := object.AsGetProp(); ok {
			jsType := tc.GetJSType(objectGetProp.Object().Node())
			if jsType.IsInterface() && objectGetProp.Property() == "prototype" {
				tc.VisitInterfacePropertyAssignment(object, getProp)
			}
		}

		tc.CheckEnumAlias(t, info, rightType, nodeToWarn)
		tc.CheckPropCreation(lvalue)

		// Prototype assignments are special, because they actually affect
		// the definition of a class. These are mostly validated
		// during TypedScopeCreator, and we only look for the "dumb" cases here.
		// object.prototype = ...;
		if pname == "prototype" {
			tc.validator.ExpectCanAssignToPrototype(objectJsType, nodeToWarn, rightType)
			return
		}

		// The generic checks for 'object.property' when 'object' is known,
		// and 'property' is declared on it.
		// object.property = ...;
		objectCastType := rhino.ObjectType_Cast(objectJsType.RestrictByNotNullOrUndefined())
		expectedPropertyType := tc.GetPropertyTypeIfDeclared(objectCastType, rhino.NewProperty_StringKey(pname))

		tc.CheckPropertyInheritanceOnAssignment(
			nodeToWarn, object, rhino.NewProperty_StringKey(pname), info, expectedPropertyType)

		// If we successfully found a non-unknown declared type, validate the assignment and don't do
		// any further checks.
		if !expectedPropertyType.IsUnknownType() {
			// Note: if the property has @implicitCast at its declaration, we don't check any
			// assignments to it.
			if !TypeCheck_PropertyIsImplicitCast(objectCastType, pname) {
				tc.validator.ExpectCanAssignToPropertyOf(
					nodeToWarn, rightType, expectedPropertyType, object.Node(), pname)
			}
			return
		}
	} else if getElem, ok := lvalue.AsGetElem(); ok && getElem.Index().Node().GetJSType().IsKnownSymbolValueType() {
		object := getElem.Object()
		objectJsType := tc.GetJSType(object.Node())
		objectCastType := rhino.ObjectType_Cast(objectJsType.RestrictByNotNullOrUndefined())
		property := getElem.Index().Node().GetJSType().ToMaybeKnownSymbolType()
		expectedPropertyType :=
			tc.GetPropertyTypeIfDeclared(objectCastType, rhino.NewProperty_SymbolKey(property))
		tc.CheckPropertyInheritanceOnAssignment(
			nodeToWarn, object, rhino.NewProperty_SymbolKey(property), info, expectedPropertyType)
	}

	// Check qualified name sets to 'object' and 'object.property'.
	// This can sometimes handle cases when the type of 'object' is not known.
	// e.g.,
	// var obj = createUnknownType();
	// /** @type {number} */ obj.foo = true;
	leftType := tc.GetJSType(lvalue.Node())
	if qName, ok := lvalue.Node().GetQualifiedNameString(); ok {
		// variable with inferred type case
		var_ := t.GetTypedScope().GetVar(qName)
		if var_ != nil {
			if var_.IsTypeInferred() {
				return
			}

			if NodeUtil_GetRootOfQualifiedName(lvalue.Node()).IsThis() &&
				rhino.StaticScope(t.GetTypedScope()) != var_.GetScope() {
				// Don't look at "this.foo" variables from other scopes.
				return
			}

			if var_.GetType() != nil {
				leftType = var_.GetType()
			}
		}
	} // Fall through case for arbitrary LHS and arbitrary RHS.

	tc.validator.ExpectCanAssignTo(nodeToWarn, rightType, leftType, msg)
}

// Java: TypeCheck.checkPropCreation(Node)
func (tc *TypeCheck) CheckPropCreation(lvalue ast.Target) {
	if getProp, ok := lvalue.AsGetProp(); ok {
		objType := tc.GetJSType(getProp.Object().Node())
		if !objType.IsEmptyType() && !objType.IsUnknownType() {
			propName := getProp.Property()
			kind := tc.typeRegistry.CanPropertyBeDefined(objType, propName)
			if kind != rhino.JSTypeRegistry_PropDefinitionKind_KNOWN {

				if objType.IsStruct() {
					if objType.RestrictByNotNullOrUndefined().IsUnionType() {
						tc.Report(getProp.Node(), TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE)
					} else {
						tc.Report(getProp.Node(), TypeCheck_ILLEGAL_PROPERTY_CREATION)
					}
				} else {
					// null checks are reported elsewhere
					if !objType.IsNoType() &&
						!objType.IsUnknownType() &&
						objType.IsSubtypeOf(tc.GetNativeType(rhino.JSTypeNative_NULL_VOID)) {
						return
					}

					tc.ReportMissingProperty(getProp.Object().Node(), objType, getProp.Node(), kind, true)
				}
			}
		}
	}
}

// Java: TypeCheck.checkPropertyInheritanceOnAssignment(Node, Node, Property.Key, JSDocInfo, JSType)
func (tc *TypeCheck) CheckPropertyInheritanceOnAssignment(
	assign rhino.Node, object ast.Expr, property rhino.Property_Key, info *rhino.JSDocInfo, propertyType rhino.JSType) {
	// Inheritance checks for prototype properties.
	//
	// TODO(nicksantos): This isn't the right place to do this check. We
	// really want to do this when we're looking at the constructor.
	// We'd find all its properties and make sure they followed inheritance
	// rules, like we currently do for @implements to make sure
	// all the methods are implemented.
	//
	// As-is, this misses many other ways to override a property.

	if prototype, ok := object.AsGetProp(); ok && prototype.Property() == "prototype" {
		// ASSIGN = assign
		//   GETPROP
		//     GETPROP = object
		//       ? = preObject
		//       STRING = "prototype"
		//     STRING = property

		preObject := prototype.Object()
		ctorType := tc.GetJSType(preObject.Node()).ToMaybeFunctionType()
		if ctorType == nil || !ctorType.HasInstanceType() {
			return
		}

		tc.CheckDeclaredPropertyAgainstNominalInheritance(
			assign, ctorType, property, info, propertyType)
		tc.CheckAbstractMethodInConcreteClass(assign, ctorType, info)
	} else {
		// ASSIGN = assign
		//   GETPROP
		//     ? = object
		//     STRING = property
		// or
		//  ASSIGN = assign
		//    GETELEM
		//      ? = object
		//      ? = property

		// We only care about checking a static property assignment.
		ctorType := tc.GetJSType(object.Node()).ToMaybeFunctionType()
		if ctorType == nil || !ctorType.HasInstanceType() {
			return
		}

		tc.CheckDeclaredPropertyAgainstPrototypalInheritance(
			assign, ctorType, property, info, propertyType)
	}
}

// Java: TypeCheck.checkPropertyInheritanceOnPrototypeLitKey(Node, Property.Key, ObjectType)
func (tc *TypeCheck) CheckPropertyInheritanceOnPrototypeLitKey(
	key rhino.Node, propertyName rhino.Property_Key, type_ rhino.ObjectType) {
	// Inheritance checks for prototype objlit properties.
	//
	// TODO(nicksantos): This isn't the right place to do this check. We
	// really want to do this when we're looking at the constructor.
	// We'd find all its properties and make sure they followed inheritance
	// rules, like we currently do for @implements to make sure
	// all the methods are implemented.
	//
	// As-is, this misses many other ways to override a property.
	//
	// object.prototype = { key: function() {} };
	tc.CheckPropertyInheritance(key, propertyName, type_.GetOwnerFunction(), type_)
}

// Java: TypeCheck.checkPropertyInheritanceOnClassMember(Node, Property.Key, FunctionType)
func (tc *TypeCheck) CheckPropertyInheritanceOnClassMember(
	key rhino.Node, propertyName rhino.Property_Key, ctorType rhino.FunctionType) {
	if key.IsStaticMember() {
		tc.CheckDeclaredPropertyAgainstPrototypalInheritance(
			key, ctorType, propertyName, key.GetJSDocInfo(), ctorType.GetPropertyTypeByKey(propertyName))
	} else {
		tc.CheckPropertyInheritance(key, propertyName, ctorType, ctorType.GetInstanceType())
	}
}

// Java: TypeCheck.checkPropertyInheritance(Node, Property.Key, FunctionType, ObjectType)
func (tc *TypeCheck) CheckPropertyInheritance(
	key rhino.Node, propertyName rhino.Property_Key, ctorType rhino.FunctionType, type_ rhino.ObjectType) {
	if ctorType == nil || !ctorType.HasInstanceType() {
		return
	}

	// The node reported on is the first child of key: the value of a string
	// key or a field, the function of a method, getter or setter, the key of
	// a computed property.
	tc.CheckDeclaredPropertyAgainstNominalInheritance(
		key.GetFirstChild(),
		ctorType,
		propertyName,
		key.GetJSDocInfo(),
		type_.GetPropertyTypeByKey(propertyName))
	tc.CheckAbstractMethodInConcreteClass(key, ctorType, key.GetJSDocInfo())
}

// VisitObjectPattern validates all keys in an object pattern.
//
// Validating the types assigned to any lhs nodes in the pattern is done at
// the ASSIGN/VAR/PARAM_LIST/etc. node.
//
// Java: TypeCheck.visitObjectPattern(Node)
func (tc *TypeCheck) VisitObjectPattern(pattern ast.ObjectPattern) {
	n := pattern.Node()
	patternType := tc.GetJSType(n)
	tc.validator.ExpectObject(n, patternType, "cannot destructure 'null' or 'undefined'")
	for prop := range pattern.Props() {
		// Go-only: the Java creates the DestructuredTarget of every prop and asks it for its
		// computed property or string key, which is the prop itself; here the prop is asked
		// for its kind, and the target is created where it is used.
		if computedProperty, ok := prop.AsComputedProp(); ok {
			tc.validator.ExpectIndexMatch(
				computedProperty.Node(), patternType, tc.GetJSType(computedProperty.Key().Node()))
		} else if stringKey, ok := prop.AsStringKey(); ok {
			if !stringKey.IsQuoted() {
				if patternType.IsDict() {
					tc.Report(stringKey.Node(), TypeValidator_ILLEGAL_PROPERTY_ACCESS, "unquoted", "dict")
				}
				// check for missing properties given `const {a} = obj;` but not `const {'a': a} = obj;`
				target := DestructuredTarget_CreateTarget(tc.typeRegistry, patternType, prop.Node())
				tc.CheckPropertyAccessForDestructuring(
					pattern, patternType, stringKey, tc.GetJSType(target.GetNode()))
			} else if patternType.IsStruct() {
				// check that we are not accessing a struct with a quoted string
				tc.Report(stringKey.Node(), TypeValidator_ILLEGAL_PROPERTY_ACCESS, "quoted", "struct")
			}
		}
	}
	tc.EnsureTyped(n)
}

// VisitObjectOrClassLiteralKey visits an object literal field definition
// `key : value`, or a class member definition `key() { ... }` If the
// `lvalue` is a prototype modification, we change the schema of the object
// type it is referring to.
//
// key is the ASSIGN, STRING_KEY, MEMBER_FUNCTION_DEF, SPREAD,
// COMPUTED_PROPERTY, MEMBER_FIELD_DEF, or COMPUTED_FIELD_DEF node; owner the
// parent node, either OBJECTLIT or CLASS_MEMBERS; ownerType the instance type
// of the enclosing object/class.
//
// Java: TypeCheck.visitObjectOrClassLiteralKey(Node, Node, JSType)
func (tc *TypeCheck) VisitObjectOrClassLiteralKey(key rhino.Node, owner rhino.Node, ownerType rhino.JSType) {
	// Do not validate object lit value types in externs. We don't really care,
	// and it makes it easier to generate externs.
	if owner.IsFromExterns() {
		tc.EnsureTyped(key)
		return
	}

	// Validate computed properties similarly to how we validate GETELEMs.
	if key.IsComputedProp() || key.IsComputedFieldDef() {
		keyType := tc.GetJSType(key.GetFirstChild())
		tc.validator.ExpectIndexMatch(key, ownerType, keyType)
		if keyType.IsKnownSymbolValueType() {
			if owner.IsClass() {
				classConstructorType := owner.GetJSType().AssertFunctionType()
				tc.CheckPropertyInheritanceOnClassMember(
					key, rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()), classConstructorType)
			} else {
				jutil.CheckStateMsg1(owner.IsObjectLit(), "Unexpected owner %s", owner)
				// Check if this property has an expected declared type on the owner.
				objectCastType := rhino.ObjectType_Cast(ownerType.RestrictByNotNullOrUndefined())
				expectedPropertyType :=
					tc.GetPropertyTypeIfDeclared(
						objectCastType, rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()))
				tc.validator.ExpectCanAssignToPropertyOf(
					key,
					expectedPropertyType,
					tc.GetJSType(key.GetSecondChild()),
					owner,
					jutil.ToJavaString(keyType.ToMaybeKnownSymbolType().GetDisplayName()))
			}
		}
		return
	}

	if key.IsBlock() {
		return
	}

	if key.IsQuotedStringKey() {
		// NB: this case will never be triggered for member functions, since we store quoted
		// member functions as computed properties. This case does apply to regular string key
		// properties, getters, and setters.
		// See also https://github.com/google/closure-compiler/issues/3071
		if ownerType.IsStruct() {
			if owner.IsClass() {
				tc.Report(key, TypeCheck_ILLEGAL_CLASS_KEY, "struct")
			} else {
				tc.Report(key, TypeCheck_ILLEGAL_OBJLIT_KEY, "struct")
			}
		}
	} else {
		// we have either a non-quoted string or a member function def
		// Neither is allowed for an @dict type except for "constructor" as a special case.
		if ownerType.IsDict() && !NodeUtil_IsEs6ConstructorMemberFunctionDef(key) {
			// Object literals annotated as @dict may only have
			// If you annotate a class with @dict, only the constructor can be a non-computed property.
			if owner.IsClass() {
				tc.Report(key, TypeCheck_ILLEGAL_CLASS_KEY, "dict")
			} else {
				tc.Report(key, TypeCheck_ILLEGAL_OBJLIT_KEY, "dict")
			}
		}
	}

	if key.IsSpread() {
		// Rely on type inference to figure out what the key/types this adds to this object.
		return
	}

	// TODO(johnlenz): Validate get and set function declarations are valid
	// as is the functions can have "extraneous" bits.

	// For getter and setter property definitions the
	// r-value type != the property type.
	rvalue := key.GetFirstChild()
	if rvalue.IsNil() {
		tc.EnsureTyped(key)
		return
	}
	rightType := TypeCheck_GetObjectLitKeyTypeFromValueType(key, tc.GetJSType(rvalue))
	if rightType == nil {
		rightType = tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}

	// Validate value is assignable to the key type.
	keyType := tc.GetJSType(key)
	propertyName := NodeUtil_GetObjectOrClassLitKeyName(key)
	allowedValueType := keyType
	if allowedValueType.IsEnumElementType() {
		allowedValueType = allowedValueType.ToMaybeEnumElementType().GetPrimitiveType()
	}

	valid :=
		tc.validator.ExpectCanAssignToPropertyOf(
			key, rightType, allowedValueType, owner, propertyName)
	if valid {
		tc.EnsureTypedWithType(key, rightType)
	} else {
		tc.EnsureTyped(key)
	}

	// Validate inheritance for classes and object literals used as prototypes
	propertyKey := rhino.NewProperty_StringKey(propertyName)
	if owner.IsClass() {
		classConstructorType := owner.GetJSType().AssertFunctionType()
		tc.CheckPropertyInheritanceOnClassMember(key, propertyKey, classConstructorType)
	} else if ownerType.ToMaybeObjectType() != nil {
		tc.CheckPropertyInheritanceOnPrototypeLitKey(key, propertyKey, ownerType.ToMaybeObjectType())
	}
}

// TypeCheck_PropertyIsImplicitCast returns true if any type in the chain has
// an implicitCast annotation for the given property.
//
// Java: TypeCheck.propertyIsImplicitCast(ObjectType, String)
func TypeCheck_PropertyIsImplicitCast(type_ rhino.ObjectType, prop string) bool {
	for ; type_ != nil; type_ = type_.GetImplicitPrototype() {
		docInfo := type_.GetOwnPropertyJSDocInfo(prop)
		if docInfo != nil && docInfo.IsImplicitCast() {
			return true
		}
	}
	return false
}

// CheckDeclaredPropertyAgainstNominalInheritance: given a ctorType, check that
// the property (propertyName), on the corresponding instance type
// (receiverType), conforms to inheritance rules.
//
// This method only checks nominal inheritance (via extends and implements
// declarations). Compare to CheckDeclaredPropertyAgainstPrototypalInheritance.
//
// To be conformant, the propertyName must
//   - Carry the @override annotation iff it is an override.
//   - Be typed as a subtype of the type of propertyName on all supertypes of
//     receiverType.
//
// Java: TypeCheck.checkDeclaredPropertyAgainstNominalInheritance(Node, FunctionType, Property.Key, JSDocInfo, JSType)
func (tc *TypeCheck) CheckDeclaredPropertyAgainstNominalInheritance(
	n rhino.Node,
	ctorType rhino.FunctionType,
	propertyName rhino.Property_Key,
	info *rhino.JSDocInfo,
	propertyType rhino.JSType) {

	// No need to check special properties; @override is not required for them, nor they are
	// manually typed by the developers.
	if propertyName.Matches("__proto__") || propertyName.Matches("constructor") {
		return
	}

	// Interfaces are checked elsewhere.
	if ctorType.IsInterface() {
		return
	}

	// If the supertype doesn't resolve correctly, we've warned about this already.
	if TypeCheck_HasUnknownOrEmptySupertype(ctorType) {
		return
	}

	foundProperty := false

	superCtor := ctorType.GetSuperClassConstructor()
	if superCtor != nil {
		propSlot := superCtor.GetInstanceType().FindClosestDefinitionByKey(propertyName)
		superClassHasProperty := propSlot != nil && !propSlot.IsOwnedByInterface()
		foundProperty = foundProperty || superClassHasProperty

		if superClassHasProperty {
			superClass := propSlot.GetOwnerInstanceType()
			superClassHasDeclaredProperty := !propSlot.GetValue().IsTypeInferred()
			if superClassHasDeclaredProperty {
				if TypeCheck_IsDeclaredLocally(ctorType, propertyName) && !TypeCheck_DeclaresOverride(info) {
					tc.compiler.Report(
						JSError_MakeWithNode(
							n,
							TypeCheck_HIDDEN_SUPERCLASS_PROPERTY,
							propertyName.HumanReadableName(),
							jutil.ToJavaString(superClass.GetReferenceName())))
				}
				tc.validator.CheckPropertyType(
					n, ctorType.GetTypeOfThis(), superClass, propertyName, propertyType)
			}
		}
	}

	for _, implementedInterface := range ctorType.GetAllImplementedInterfaces() {
		if implementedInterface.IsUnknownType() || implementedInterface.IsEmptyType() {
			continue
		}
		jutil.CheckStateMsg1(implementedInterface.IsInstanceType(), "%s", implementedInterface)
		propSlot := implementedInterface.FindClosestDefinitionByKey(propertyName)
		interfaceHasProperty := propSlot != nil
		foundProperty = foundProperty || interfaceHasProperty

		if interfaceHasProperty && !TypeCheck_DeclaresOverride(info) {
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_HIDDEN_INTERFACE_PROPERTY,
					propertyName.HumanReadableName(),
					jutil.ToJavaString(propSlot.GetOwnerInstanceType().GetReferenceName())))
		}
	}

	if !foundProperty && TypeCheck_DeclaresOverride(info) {
		tc.compiler.Report(
			JSError_MakeWithNode(
				n,
				TypeCheck_UNKNOWN_OVERRIDE,
				propertyName.HumanReadableName(),
				jutil.ToJavaString(ctorType.GetInstanceType().GetReferenceName())))
	}
}

// Java: TypeCheck.isDeclaredLocally(FunctionType, Property.Key)
func TypeCheck_IsDeclaredLocally(ctorType rhino.FunctionType, propertyName rhino.Property_Key) bool {
	jutil.CheckState(ctorType.IsConstructor())
	return ctorType.GetPrototype().HasOwnPropertyByKey(propertyName) ||
		ctorType.GetInstanceType().HasOwnPropertyByKey(propertyName)
}

// CheckDeclaredPropertyAgainstPrototypalInheritance: given a receiverType,
// check that the property (propertyName) conforms to inheritance rules.
//
// This method only checks prototypal inheritance (via the prototype chain).
// Compare to CheckDeclaredPropertyAgainstNominalInheritance.
//
// To be conformant, the propertyName must
//   - Carry the @override annotation iff it is an override.
//   - Be typed as a subtype of the type of propertyName on all supertypes of
//     receiverType.
//
// Java: TypeCheck.checkDeclaredPropertyAgainstPrototypalInheritance(Node, ObjectType, Property.Key, JSDocInfo, JSType)
func (tc *TypeCheck) CheckDeclaredPropertyAgainstPrototypalInheritance(
	n rhino.Node,
	receiverType rhino.ObjectType,
	propertyName rhino.Property_Key,
	info *rhino.JSDocInfo,
	propertyType rhino.JSType) {
	// TODO(nickreid): Right now this is only expected to run on ctors. However, it wouldn't be bad
	// if it ran on more things. Consider a precondition on the arguments.

	declaredOverride := TypeCheck_DeclaresOverride(info)
	var supertypeWithProperty rhino.ObjectType
	for type_ := receiverType.GetImplicitPrototype(); type_ != nil; type_ = type_.GetImplicitPrototype() {
		// We want to report the supertype that actually had the overridden declaration.
		if type_.HasOwnPropertyByKey(propertyName) {
			// We only care about the lowest match in the chain because it must be the most
			// specific.
			supertypeWithProperty = type_
			break
		}
	}

	if supertypeWithProperty == nil {
		// TODO(b/144327372): stop loosening typechecking for forward supertype references once
		// the type system correctly models this case.
		if declaredOverride && !receiverType.LoosenTypecheckingDueToForwardReferencedSupertype() {
			tc.compiler.Report(
				JSError_MakeWithNode(
					n, //
					TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE,
					propertyName.HumanReadableName(),
					receiverType.String()))
		}
	} else {
		if !declaredOverride {
			tc.compiler.Report(
				JSError_MakeWithNode(
					n, //
					TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY,
					propertyName.HumanReadableName(),
					supertypeWithProperty.String()))
		}

		overriddenPropertyType := supertypeWithProperty.GetPropertyTypeByKey(propertyName)
		if !propertyType.IsSubtypeOf(overriddenPropertyType) {
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH,
					propertyName.HumanReadableName(),
					supertypeWithProperty.String(),
					overriddenPropertyType.String(),
					propertyType.String()))
		}
	}
}

// Java: TypeCheck.checkAbstractMethodInConcreteClass(Node, FunctionType, JSDocInfo)
func (tc *TypeCheck) CheckAbstractMethodInConcreteClass(n rhino.Node, ctorType rhino.FunctionType, info *rhino.JSDocInfo) {
	if info == nil || !info.IsAbstract() {
		return
	}

	if ctorType.IsConstructor() && !ctorType.IsAbstract() {
		tc.Report(n, TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS)
	}
}

// TypeCheck_HasUnknownOrEmptySupertype: given a constructor or an interface
// type, find out whether the unknown type is a supertype of the current type.
//
// Java: TypeCheck.hasUnknownOrEmptySupertype(FunctionType)
func TypeCheck_HasUnknownOrEmptySupertype(ctor rhino.FunctionType) bool {
	jutil.CheckArgument(ctor.IsConstructor() || ctor.IsInterface())
	jutil.CheckArgument(!ctor.IsUnknownType())

	// The type system should notice inheritance cycles on its own
	// and break the cycle.
	for {
		maybeSuperInstanceType := ctor.GetPrototype().GetImplicitPrototype()
		if maybeSuperInstanceType == nil {
			return false
		}
		if maybeSuperInstanceType.IsUnknownType() || maybeSuperInstanceType.IsEmptyType() {
			return true
		}
		ctor = maybeSuperInstanceType.GetConstructor()
		if ctor == nil {
			return false
		}
		jutil.CheckState(ctor.IsConstructor() || ctor.IsInterface())
	}
}

// TypeCheck_GetObjectLitKeyTypeFromValueType returns the type expected when
// using the key.
//
// key is a OBJECTLIT key node.
//
// Java: TypeCheck.getObjectLitKeyTypeFromValueType(Node, JSType)
func TypeCheck_GetObjectLitKeyTypeFromValueType(key rhino.Node, valueType rhino.JSType) rhino.JSType {
	if valueType != nil {
		switch key.GetToken() {
		case rhino.Token_GETTER_DEF:
			// GET must always return a function type.
			if valueType.IsFunctionType() {
				fntype := valueType.ToMaybeFunctionType()
				valueType = fntype.GetReturnType()
			} else {
				return nil
			}
		case rhino.Token_SETTER_DEF:
			if valueType.IsFunctionType() {
				// SET must always return a function type.
				fntype := valueType.ToMaybeFunctionType()
				param := fntype.GetParameters()[0]
				// SET function must always have one parameter.
				valueType = param.GetJSType()
			} else {
				return nil
			}
		default:
		}
	}
	return valueType
}

// VisitInterfacePropertyAssignment visits an lvalue node for cases such as
//
//	interface.prototype.property = ...;
//
// Java: TypeCheck.visitInterfacePropertyAssignment(Node, Node)
func (tc *TypeCheck) VisitInterfacePropertyAssignment(object ast.Expr, lvalue ast.GetProp) {
	assign, ok := ast.AsAssign(lvalue.Node().GetParent())
	if !ok {
		// assignments to interface properties cannot be in destructuring patterns or for-of loops
		tc.ReportInvalidInterfaceMemberDeclaration(object.Node())
		return
	}
	rvalue := assign.Value()
	rvalueType := tc.GetJSType(rvalue.Node())

	// Only 2 values are allowed for interface methods:
	//    goog.abstractMethod
	//    function () {};
	// Other (non-method) interface properties must be stub declarations without assignments, e.g.
	//     someinterface.prototype.nonMethodProperty;
	// which is why we enforce that `rvalueType.isFunctionType()`.
	if !rvalueType.IsFunctionType() {
		tc.ReportInvalidInterfaceMemberDeclaration(object.Node())
	}

	if fn, ok := rvalue.AsFunction(); ok && !NodeUtil_IsEmptyBlock(fn.Body().Node()) {
		abstractMethodName := tc.compiler.GetCodingConvention().GetAbstractMethodName()
		tc.compiler.Report(JSError_MakeWithNode(object.Node(), TypeCheck_INTERFACE_METHOD_NOT_EMPTY, jutil.ToJavaString(abstractMethodName)))
	}
}

// Java: TypeCheck.reportInvalidInterfaceMemberDeclaration(Node)
func (tc *TypeCheck) ReportInvalidInterfaceMemberDeclaration(interfaceNode rhino.Node) {
	abstractMethodName := tc.compiler.GetCodingConvention().GetAbstractMethodName()
	// This is bad i18n style but we don't localize our compiler errors.
	abstractMethodMessage := ""
	if abstractMethodName != nil {
		abstractMethodMessage = ", or " + *abstractMethodName
	}
	tc.compiler.Report(
		JSError_MakeWithNode(interfaceNode, TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION, abstractMethodMessage))
}

// VisitName visits a NAME node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited. Returns whether the node is typeable or not.
//
// Java: TypeCheck.visitName(NodeTraversal, Node, Node)
func (tc *TypeCheck) VisitName(t *NodeTraversal, name ast.Name, parent rhino.Node) bool {
	n := name.Node()
	// Skip empty function expression names. They don't need a type.
	// (Go-only: the rhino API, because Function.Name() reports an anonymous function's
	// empty NAME as no name.)
	isFunctionName := n.GetParent().IsFunction() && n.IsFirstChildOf(parent)
	if isFunctionName && name.Ident() == "" {
		return false
	}

	// the compiler's unknown type reporting allows certain names to be untyped if declared but
	// never used in an expression. The associated AST nodes still need to be typed to prevent
	// future passes from crashing, though.
	reportIfMissingType := true

	// At this stage, we need to determine whether this is a leaf
	// node in an expression (which therefore needs to have a type
	// assigned for it) versus some other decorative node that we
	// can safely ignore.
	parentNodeType := parent.GetToken()
	if isFunctionName ||
		parentNodeType == rhino.Token_CATCH ||
		parentNodeType == rhino.Token_PARAM_LIST ||
		NodeUtil_IsNameDeclaration(parent) {
		reportIfMissingType = false
	}

	// Not need to type first key in for-in or for-of.
	if loop, ok := ast.AsEnhancedFor(parent); ok && loop.Target().Is(n) {
		reportIfMissingType = false
	}

	type_ := n.GetJSType()
	if type_ == nil {
		type_ = tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		// TODO(b/149843534): crash instead of defaulting to '?' when the var is null.
		var_ := t.GetTypedScope().GetVar(name.Ident())
		if var_ != nil {
			type_ = jutil.CheckNotNull(var_.GetType())
		}
	}
	tc.EnsureTypedWithType(n, type_)
	return reportIfMissingType
}

// CheckForOfTypes visits the loop variable of a FOR_OF and FOR_AWAIT_OF and
// verifies the type being assigned to it.
//
// Java: TypeCheck.checkForOfTypes(NodeTraversal, Node)
func (tc *TypeCheck) CheckForOfTypes(t *NodeTraversal, forOf ast.ForOfLoop) {
	iterable := forOf.Iterable().Node()
	iterableType := tc.GetJSType(iterable)
	var actualType rhino.JSType

	if forOf.Token() == rhino.Token_FOR_AWAIT_OF {
		maybeType :=
			tc.validator.ExpectAutoboxesToIterableOrAsyncIterable(
				iterable,
				iterableType,
				"Can only async iterate over a (non-null) Iterable or AsyncIterable type")

		if maybeType == nil {
			// Not iterable or async iterable, error reported by
			// expectAutoboxesToIterableOrAsyncIterable.
			return
		}

		actualType = maybeType
	} else {
		tc.validator.ExpectAutoboxesToIterable(
			iterable, iterableType, "Can only iterate over a (non-null) Iterable type")

		actualType =
			// Convert primitives to their wrapper type and remove null/undefined
			// If iterable is a union type, autoboxes each member of the union.
			iterableType.
				Autobox().
				GetTemplateTypeMap().
				GetResolvedTemplateType(tc.typeRegistry.GetIterableValueTemplate())
	}

	lhs := typeCheck_forOfTarget(forOf.Target())
	tc.CheckCanAssignToWithScope(
		t,
		forOf.Node(),
		lhs,
		actualType,
		lhs.Node().GetJSDocInfo(),
		"declared type of for-of loop variable does not match inferred type")
}

// typeCheck_forOfTarget returns what a for-of loop assigns each value to:
// its target, or what its declaration declares.
//
// Go-only: the Java steps down from the declaration by position.
func typeCheck_forOfTarget(target ast.ForTarget) ast.Target {
	decl, ok := target.AsNameDeclaration()
	if !ok {
		return ast.MustTarget(target.Node())
	}
	// e.g. get "x" given the VAR in "for (var x of arr) {"
	declarator, ok := decl.OnlyDeclarator()
	jutil.CheckStateMsg1(ok, "%s", decl.Node())
	if destructuring, ok := declarator.AsDestructuringLhs(); ok {
		// e.g. get `[x, y]` given the VAR in `for (var [x, y] of arr) {`
		return destructuring.Pattern().AsTarget()
	}
	return ast.MustName(declarator.Node()).AsTarget()
}

// VisitGetProp visits a GETPROP node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitGetProp(NodeTraversal, Node)
func (tc *TypeCheck) VisitGetProp(t *NodeTraversal, getProp ast.GetProp) {
	// obj.prop or obj.method()
	// Lots of types can appear on the left, a call to a void function can
	// never be on the left. getPropertyType will decide what is acceptable
	// and what isn't.
	n := getProp.Node()
	objNode := getProp.Object().Node()
	childType := tc.GetJSType(objNode)

	if childType.IsDict() {
		tc.Report(n, TypeValidator_ILLEGAL_PROPERTY_ACCESS, "'.'", "dict")
	} else if tc.validator.ExpectNotNullOrUndefined(
		t, n, childType, "No properties on this expression", tc.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)) {
		tc.CheckPropertyAccessForGetProp(getProp.AsNormalOrOptChainGetProp())
	}
	tc.EnsureTyped(n)
}

// VisitOptChainGetProp visits a OPTCHAIN_GETPROP node.
//
// Java: TypeCheck.visitOptChainGetProp(Node)
func (tc *TypeCheck) VisitOptChainGetProp(optChainGetProp ast.OptChainGetProp) {
	// obj?.prop
	// Unlike GETPROP, a call to a void function can also be on the lhs for an OPTCHAIN_GETPROP.
	objNode := optChainGetProp.Object().Node()
	childType := tc.GetJSType(objNode)

	if childType.IsDict() {
		// Go-only: the Java reports on `property`, optChainGetProp.getLastChild(), which is the
		// object: the property name is the node's string, not a child.
		tc.Report(objNode, TypeValidator_ILLEGAL_PROPERTY_ACCESS, "'?.'", "dict")
	} else if !childType.IsUnknownType() {
		tc.CheckPropertyAccessForGetProp(optChainGetProp.AsNormalOrOptChainGetProp())
	}
	tc.EnsureTyped(optChainGetProp.Node())
}

// Java: TypeCheck.checkPropertyAccessForGetProp(Node)
func (tc *TypeCheck) CheckPropertyAccessForGetProp(getProp ast.NormalOrOptChainGetProp) {
	objNode := getProp.Object().Node()
	objType := tc.GetJSType(objNode)
	propType := tc.GetJSType(getProp.Node())

	tc.CheckAbstractPropertyAccess(getProp)
	tc.CheckPropertyAccess(objType, getProp.Node(), propType, objNode)
}

// Java: TypeCheck.checkPropertyAccessForDestructuring(Node, JSType, Node, JSType)
func (tc *TypeCheck) CheckPropertyAccessForDestructuring(
	pattern ast.ObjectPattern, objectType rhino.JSType, stringKey ast.StringKey, inferredPropType rhino.JSType) {

	// Get the object node being destructured when it exists. These cases have an actual node `obj`:
	//   const {a} = obj;
	//   ({a} = obj);
	// while these do not:
	//    for (const {a} of arrayOfObjects) {
	//    const {{a}} = obj;
	var objNode rhino.Node
	patternParent := pattern.Node().GetParent()
	if assign, ok := ast.AsAssign(patternParent); ok {
		objNode = assign.Value().Node()
	} else if declaration, ok := ast.AsDestructuringLhs(patternParent); ok {
		if value, ok := declaration.Init(); ok {
			objNode = value.Node()
		}
	}
	tc.CheckPropertyAccess(objectType, stringKey.Node(), inferredPropType, objNode)
}

// CheckAbstractPropertyAccess warns if @abstract methods are dereferenced,
// with some false negatives.
//
// This method only handles some cases that we are certain are incorrect.
// e.g. we are lenient about union types, and we don't track an abstract
// methods once they are reassigned to new variables.
//
// This method's logic is complicated to avoid spurious warnings. Sometimes a
// reference to something typed @abstract is okay. For example, don't warn on
// `this.foo();` in
// `/** @abstract * / class Base { /** @abstract * / foo() {} bar() { this.foo(); } } }`
//
// The `this` object with which `Base.prototype.bar` is called must be a
// concrete subclass of Base. To avoid false positives, we warn only in the
// following cases:
//   - (a) the function is accessed off `super`
//   - (b) the function is accessed off a .prototype
//   - (c) the function is transpiled from a goog.base superclass reference
//
// Java: TypeCheck.checkAbstractPropertyAccess(Node)
func (tc *TypeCheck) CheckAbstractPropertyAccess(method ast.NormalOrOptChainGetProp) {

	if NodeUtil_IsLhsOfAssign(method.Node()) {
		// Allow declaring abstract methods. (This assumes they are never re-assigned)
		//   /** @abstract */ Foo.prototype.bar = function() {}
		return
	}

	methodType := tc.GetJSType(method.Node()).ToMaybeFunctionType()
	if methodType == nil || !methodType.IsAbstract() || methodType.IsConstructor() {
		// Ignore non-abstract methods and @abstract constructors. An @abstract constructor is still
		// callable.
		return
	}

	object := method.Object()
	if object.Node().IsSuper() {
		// case (a)
		// `super.foo()` definitely refers to `Superclass.prototype.foo`, not an override.
		// At parse time, `Subclass.prototype` becomes a lower bound for what `super` evaluates to,
		// even if the `this` object changes. So `super` will never resolve to a concrete subclass.
		tc.Report(method.Node(), TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE, jutil.ToJavaString(methodType.GetDisplayName()))
	} else if objectGetProp, ok := object.AsGetProp(); ok {
		objectProp := objectGetProp.Property()
		if objectProp == "prototype" || // case (b), e.g. `Foo.prototype.bar`
			tc.compiler.GetCodingConvention().IsSuperClassReference(objectProp) { // case (c)
			tc.Report(method.Node(), TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE, jutil.ToJavaString(methodType.GetDisplayName()))
		}
	}
}

// CheckPropertyAccess emits a warning if we can prove that a property cannot
// possibly be defined on an object. Note the difference between JS and a
// strictly statically typed language: we're checking if the property *cannot
// be defined*, whereas a java compiler would check if the property *can be
// undefined.
//
// This method handles property access in both GETPROPs and object
// destructuring. Consequentially some of its arguments are optional - the
// actual object node and the getprop - while others are required.
//
// objNode is the actual node representing the object we're accessing.
// optional because destructuring accesses MAY not have an actual object node.
//
// Java: TypeCheck.checkPropertyAccess(JSType, Node, JSType, Node)
func (tc *TypeCheck) CheckPropertyAccess(
	childType rhino.JSType, propNode rhino.Node, propType rhino.JSType, objNode rhino.Node) {
	isGetprop := NodeUtil_IsNormalOrOptChainGetProp(propNode)
	propName := propNode.GetString()

	if propType.Equals(tc.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)) {
		childType = childType.Autobox()
		objectType := rhino.ObjectType_Cast(childType)
		if objectType != nil {
			// We special-case object types so that checks on enums can be
			// much stricter, and so that we can use hasProperty (which is much
			// faster in most cases).
			if !objectType.HasProperty(propName) ||
				objectType.Equals(tc.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)) {
				if _, isEnum := objectType.(*rhino.EnumType); isEnum {
					tc.Report(propNode, TypeCheck_INEXISTENT_ENUM_ELEMENT, propName)
				} else {
					tc.CheckPropertyAccessHelper(objectType, propNode, objNode, false)
				}
			}
		} else {
			tc.CheckPropertyAccessHelper(childType, propNode, objNode, false)
		}
	} else if childType.IsUnionType() && isGetprop && !tc.IsLValueGetProp(propNode) {
		// NOTE: strict property assignment checks are done on assignment.
		tc.CheckPropertyAccessHelper(childType, propNode, objNode, true)
	}
}

// Java: TypeCheck.isLValueGetProp(Node)
func (tc *TypeCheck) IsLValueGetProp(n rhino.Node) bool {
	parent := n.GetParent()
	// TODO(b/77597706): this won't work for destructured lvalues
	if incDec, ok := ast.AsIncDec(parent); ok {
		return incDec.Target().Is(n)
	}
	assign, ok := ast.AsAssignment(parent)
	return ok && assign.Target().Is(n)
}

// CheckPropertyAccessHelper: strictCheck is whether this is a check that is
// only performed when "strict missing properties" checks are enabled.
//
// Java: TypeCheck.checkPropertyAccessHelper(JSType, Node, Node, boolean)
func (tc *TypeCheck) CheckPropertyAccessHelper(
	objectType rhino.JSType, propNode rhino.Node, objNode rhino.Node, strictCheck bool) {
	isGetprop := NodeUtil_IsNormalOrOptChainGetProp(propNode)
	propName := propNode.GetString()

	if !tc.reportMissingProperties ||
		objectType.IsEmptyType() ||
		(isGetprop && tc.AllowStrictPropertyAccessOnNode(propNode)) {
		return
	}

	kind := tc.typeRegistry.CanPropertyBeDefined(objectType, propName)
	if kind == rhino.JSTypeRegistry_PropDefinitionKind_KNOWN {
		return
	}
	// If the property definition is known, but only loosely associated,
	// only report a "strict error" which can be optional as code is migrated.
	isLooselyAssociated :=
		kind == rhino.JSTypeRegistry_PropDefinitionKind_LOOSE || kind == rhino.JSTypeRegistry_PropDefinitionKind_LOOSE_UNION
	isUnknownType := objectType.IsUnknownType()
	if isLooselyAssociated && isUnknownType {
		// We still don't want to report this.
		return
	}

	isStruct := objectType.IsStruct()
	loosePropertyDeclaration := !isStruct && isGetprop && tc.IsQNameAssignmentTarget(propNode)
	// always false for destructuring
	maybePropExistenceCheck :=
		!isStruct && isGetprop && tc.AllowLoosePropertyAccessOnNode(propNode)
	// Traditionally, we would not report a warning for "loose" properties, but we want to be
	// able to be more strict, so introduce an optional warning.
	strictReport :=
		strictCheck || isLooselyAssociated || loosePropertyDeclaration || maybePropExistenceCheck

	tc.ReportMissingProperty(objNode, objectType, propNode, kind, strictReport)
}

// Java: TypeCheck.reportMissingProperty(Node, JSType, Node, PropDefinitionKind, boolean)
func (tc *TypeCheck) ReportMissingProperty(
	objNode rhino.Node,
	objectType rhino.JSType,
	propNode rhino.Node,
	kind rhino.JSTypeRegistry_PropDefinitionKind,
	strictReport bool) {
	propName := propNode.GetString()

	isObjectType := objectType.Equals(tc.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
	lowConfidence := objectType.IsUnknownType() || objectType.IsAllType() || isObjectType

	isKnownToUnionMember := kind == rhino.JSTypeRegistry_PropDefinitionKind_LOOSE_UNION

	var pair *TypeCheck_SuggestionPair
	if !lowConfidence && !isKnownToUnionMember {
		pair = TypeCheck_GetClosestPropertySuggestion(objectType, propName, (jutil.StringLength(propName)-1)/4)
	}
	var objectTypeName string
	if !objNode.IsNil() {
		objectTypeName = tc.typeRegistry.GetReadableTypeName(objNode)
	} else {
		objectTypeName = objectType.String()
	}
	if pair != nil {
		var reportType *DiagnosticType
		if strictReport {
			reportType = TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION
		} else {
			reportType = TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION
		}
		tc.Report(
			propNode,
			reportType,
			propName,
			objectTypeName,
			pair.suggestion)
	} else {
		var reportType *DiagnosticType
		if strictReport {
			if isKnownToUnionMember {
				reportType = TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY
			} else {
				reportType = TypeCheck_STRICT_INEXISTENT_PROPERTY
			}
		} else if lowConfidence {
			reportType = TypeCheck_POSSIBLE_INEXISTENT_PROPERTY
		} else {
			reportType = TypeCheck_INEXISTENT_PROPERTY
		}
		tc.Report(
			propNode,
			reportType,
			propName,
			objectTypeName)
	}
}

// Java: TypeCheck.allowStrictPropertyAccessOnNode(Node)
func (tc *TypeCheck) AllowStrictPropertyAccessOnNode(n rhino.Node) bool {
	return n.GetParent().IsTypeOf()
}

// AllowLoosePropertyAccessOnNode checks whether a property access is either
// (1) a stub declaration, or (2) a property presence or absence test.
//
// Presence and absence are both allowed here because both are valid for
// conditional use of a property (e.g. `if (x.y != null) use(x.y);` but also
// the opposite, `if (x.y == null) return; use(x.y);`).
//
// Java: TypeCheck.allowLoosePropertyAccessOnNode(Node)
func (tc *TypeCheck) AllowLoosePropertyAccessOnNode(n rhino.Node) bool {
	parent := n.GetParent()
	return NodeUtil_IsPropertyTest(tc.compiler, n) ||
		NodeUtil_IsPropertyAbsenceTest(n) ||
		// Stub property declaration
		(n.IsQualifiedName() && parent.IsExprResult())
}

// Java: TypeCheck.isQNameAssignmentTarget(Node)
func (tc *TypeCheck) IsQNameAssignmentTarget(n rhino.Node) bool {
	assign, ok := ast.AsAssign(n.GetParent())
	return n.IsQualifiedName() && ok && assign.Target().Is(n)
}

// Java: TypeCheck.getClosestPropertySuggestion(JSType, String, int)
func TypeCheck_GetClosestPropertySuggestion(
	objectType rhino.JSType, propName string, maxDistance int) *TypeCheck_SuggestionPair {
	return nil
}

// VisitGetElem visits a GETELEM node.
//
// Java: TypeCheck.visitGetElem(Node)
func (tc *TypeCheck) VisitGetElem(getElem ast.NormalOrOptChainGetElem) {
	n := getElem.Node()
	tc.validator.ExpectIndexMatch(n, tc.GetJSType(getElem.Object().Node()), tc.GetJSType(getElem.Index().Node()))
	tc.EnsureTyped(n)
}

// VisitOptChainGetElem visits a OPTCHAIN_GETELEM node.
//
// Java: TypeCheck.visitOptChainGetElem(Node)
func (tc *TypeCheck) VisitOptChainGetElem(optChainGetElem ast.OptChainGetElem) {
	obj := optChainGetElem.Object().Node()
	objType := tc.GetJSType(obj)
	if objType.IsNullType() || objType.IsVoidType() {
		// no error reported when conditionally checking a null or void lhs object.
		// e.g. `a?.[b]` should not report if `a` is null`
		tc.EnsureTypedWithNative(optChainGetElem.Node(), rhino.JSTypeNative_VOID_TYPE)
	} else {
		tc.VisitGetElem(optChainGetElem.AsNormalOrOptChainGetElem())
	}
}

// VisitVar visits a VAR node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitVar(NodeTraversal, Node)
func (tc *TypeCheck) VisitVar(t *NodeTraversal, decl ast.NameDeclaration) {
	n := decl.Node()
	// Handle var declarations in for-of loops separately from regular var declarations.
	if _, inLoop := ast.AsEnhancedFor(n.GetParent()); inLoop {
		return
	}

	// TODO(nicksantos): Fix this so that the doc info always shows up
	// on the NAME node. We probably want to wait for the parser
	// merge to fix this.
	var varInfo *rhino.JSDocInfo
	if _, ok := decl.OnlyDeclarator(); ok {
		varInfo = n.GetJSDocInfo()
	}
	for declarator := range decl.Declarators() {
		if name, ok := declarator.AsName(); ok {
			if value, ok := name.Init(); ok {
				valueType := tc.GetJSType(value.Node())
				info := name.Node().GetJSDocInfo()
				if info == nil {
					info = varInfo
				}

				tc.CheckEnumAlias(t, info, valueType, value.Node())
				tc.CheckCanAssignToWithScope(t, value.Node(), name.AsTarget(), valueType, info, "initializing variable")
			}
		} else {
			destructuring := ast.MustDestructuringLhs(declarator.Node())
			// Outside a for-in or for-of loop, a destructuring declaration has a value.
			value, ok := destructuring.Init()
			jutil.CheckStateMsg1(ok, "%s", destructuring.Node())
			valueType := tc.GetJSType(value.Node())
			tc.CheckCanAssignToWithScope(
				t, destructuring.Node(), destructuring.Pattern().AsTarget(), valueType /* info= */, nil, "initializing variable")
		}
	}
}

// VisitNew visits a NEW node.
//
// Java: TypeCheck.visitNew(Node)
func (tc *TypeCheck) VisitNew(newNode ast.New) {
	n := newNode.Node()
	constructor := newNode.Callee().Node()
	type_ := tc.GetJSType(constructor).RestrictByNotNullOrUndefined()
	if !tc.CouldBeAConstructor(type_) ||
		type_.Equals(tc.typeRegistry.GetNativeType(rhino.JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE)) ||
		type_.Equals(tc.typeRegistry.GetNativeType(rhino.JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE)) {
		tc.Report(n, TypeCheck_NOT_A_CONSTRUCTOR, type_.String())
		tc.EnsureTyped(n)
		return
	}

	fnType := type_.ToMaybeFunctionType()
	if fnType != nil && fnType.HasInstanceType() {
		if fnType.IsAbstract() {
			tc.Report(n, TypeCheck_INSTANTIATE_ABSTRACT_CLASS)
		}

		tc.VisitArgumentList(newNode.AsCallOrNew(), fnType)

		objType := fnType.GetInstanceType()
		if objType != nil {
			tc.EnsureTypedWithType(n, objType)
		} else {
			tc.EnsureTypedWithType(n, tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
		}
	} else {
		tc.EnsureTyped(n)
	}
}

// Java: TypeCheck.couldBeAConstructor(JSType)
func (tc *TypeCheck) CouldBeAConstructor(type_ rhino.JSType) bool {
	return type_.IsConstructor() || type_.IsEmptyType() || type_.IsUnknownType()
}

// CheckInterfaceConflictProperties checks whether there's any property
// conflict for a particular super interface.
//
// n is the node being visited; functionName the function name being checked;
// properties the property names in the super interfaces that have been
// visited; currentProperties the property names in the super interface that
// have been visited; interfaceType the super interface that is being visited.
//
// Java: TypeCheck.checkInterfaceConflictProperties(Node, String, Map<Property.Key, ObjectType>, Map<Property.Key, ObjectType>, ObjectType)
func (tc *TypeCheck) CheckInterfaceConflictProperties(
	n rhino.Node,
	functionName string,
	properties *jutil.ObjectLinkedHashMap[rhino.Property_Key, rhino.ObjectType],
	currentProperties *jutil.ObjectLinkedHashMap[rhino.Property_Key, rhino.ObjectType],
	interfaceType rhino.ObjectType) {
	implicitProto := interfaceType.GetImplicitPrototype()
	var currentPropertyNames []rhino.Property_Key
	if implicitProto == nil {
		// This can be the case if interfaceType is proxy to a non-existent
		// object (which is a bad type annotation, but shouldn't crash).
		currentPropertyNames = nil
	} else {
		currentPropertyNames = implicitProto.GetOwnPropertyKeys()
	}
	for _, name := range currentPropertyNames {
		oType := properties.Get(name)
		currentProperties.Put(name, interfaceType)
		if oType != nil {
			thisPropType := interfaceType.GetPropertyTypeByKey(name)
			oPropType := oType.GetPropertyTypeByKey(name)
			if thisPropType.IsSubtypeOfWithMode(oPropType, tc.subtypingMode) ||
				oPropType.IsSubtypeOfWithMode(thisPropType, tc.subtypingMode) ||
				(thisPropType.IsFunctionType() &&
					oPropType.IsFunctionType() &&
					thisPropType.
						ToMaybeFunctionType().
						HasEqualCallType(oPropType.ToMaybeFunctionType())) {
				continue
			}
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE,
					functionName,
					name.HumanReadableName(),
					oType.String(),
					interfaceType.String()))
		}
	}
	for _, iType := range interfaceType.GetCtorExtendedInterfaces() {
		tc.CheckInterfaceConflictProperties(n, functionName, properties, currentProperties, iType)
	}
}

// VisitFunction visits a FUNCTION node.
//
// Java: TypeCheck.visitFunction(Node)
func (tc *TypeCheck) VisitFunction(n rhino.Node) {
	if NodeUtil_IsEs6Constructor(n) {
		return // These will be checked via the CLASS node.
	}

	functionType := rhino.JSType_ToMaybeFunctionType(n.GetJSType())
	if functionType.IsConstructor() {
		tc.CheckConstructor(n, functionType)
	} else if functionType.IsInterface() {
		tc.CheckInterface(n, functionType)
	} else if n.IsAsyncGeneratorFunction() {
		// An async generator function must return a AsyncGenerator or supertype of AsyncGenerator
		returnType := functionType.GetReturnType()
		tc.validator.ExpectAsyncGeneratorSupertype(
			n, returnType, "An async generator function must return a (supertype of) AsyncGenerator")
	} else if n.IsGeneratorFunction() {
		// A generator function must return a Generator or supertype of Generator
		returnType := functionType.GetReturnType()
		tc.validator.ExpectGeneratorSupertype(
			n, returnType, "A generator function must return a (supertype of) Generator")

	} else if n.IsAsyncFunction() {
		// An async function must return a Promise or supertype of Promise
		returnType := functionType.GetReturnType()
		tc.validator.ExpectValidAsyncReturnType(n, returnType)
	}
}

// VisitClass visits a CLASS node.
//
// Java: TypeCheck.visitClass(Node)
func (tc *TypeCheck) VisitClass(class ast.Class) {
	n := class.Node()
	functionType := rhino.JSType_ToMaybeFunctionType(n.GetJSType())
	if extendsClause, ok := class.Extends(); ok {
		// Ensure that the `extends` clause is actually a constructor or interface.  If it is, but
		// it's the wrong one then checkConstructor or checkInterface will warn.
		superType := extendsClause.Node().GetJSType()
		if superType.IsConstructor() || superType.IsInterface() {
			tc.validator.ExpectExtends(n, functionType, superType.ToMaybeFunctionType())
		} else if !superType.IsUnknownType() {
			// Only give this error for supertypes *known* to be wrong - unresolved types are OK here.
			kind := "interface"
			if functionType.IsConstructor() {
				kind = "constructor"
			}
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_CONFLICTING_EXTENDED_TYPE,
					kind,
					tc.GetBestFunctionName(n)))
		}
	}
	if functionType.IsConstructor() {
		tc.CheckConstructor(n, functionType)
	} else if functionType.IsInterface() {
		tc.CheckInterface(n, functionType)
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"CLASS node's type must be either constructor or interface: " + jutil.ToJavaString(functionType)))
	}
}

// CheckConstructor checks a constructor, which may be either an ES5-style
// FUNCTION node, or a CLASS node.
//
// Java: TypeCheck.checkConstructor(Node, FunctionType)
func (tc *TypeCheck) CheckConstructor(n rhino.Node, functionType rhino.FunctionType) {
	baseConstructor := functionType.GetSuperClassConstructor()
	if !jutil.ObjectsEquals(baseConstructor, tc.GetNativeType(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE)) &&
		baseConstructor != nil &&
		baseConstructor.IsInterface() {
		// Warn if a class extends an interface.
		tc.compiler.Report(
			JSError_MakeWithNode(n, TypeCheck_CONFLICTING_EXTENDED_TYPE, "constructor", tc.GetBestFunctionName(n)))
	} else {
		if n.IsFunction() &&
			baseConstructor != nil &&
			!baseConstructor.GetSource().IsNil() &&
			baseConstructor.GetSource().IsClass() &&
			!functionType.GetSource().IsClass() {
			// Warn if an ES5 class extends an ES6 class.
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS,
					jutil.ToJavaString(functionType.GetDisplayName()),
					jutil.ToJavaString(baseConstructor.GetDisplayName())))
		}
		if baseConstructor != nil {
			tc.CheckStructDictSubtyping(n, functionType, baseConstructor)
		}

		// Warn if any @implemented types are not interfaces or if there are any duplicates
		implementedInterfaces := functionType.GetOwnImplementedInterfaces()
		if len(implementedInterfaces) > 0 {
			alreadySeenInterfaces := jutil.NewObjectLinkedHashSet[rhino.JSType]()

			for _, baseInterfaceObjectType := range implementedInterfaces {
				var baseInterface rhino.JSType = baseInterfaceObjectType
				badImplementedType := false
				baseInterfaceObj := rhino.ObjectType_Cast(baseInterface)
				if baseInterfaceObj != nil {
					interfaceConstructor := baseInterfaceObj.GetConstructor()
					if interfaceConstructor != nil && !interfaceConstructor.IsInterface() {
						badImplementedType = true
					}
				} else {
					badImplementedType = true
				}
				if badImplementedType {
					tc.Report(n, TypeCheck_BAD_IMPLEMENTED_TYPE, tc.GetBestFunctionName(n))
				}
				// Disallow implementing both `Foo<string>` and `Foo<number>`, for example.
				baseInterface = tc.NormalizeTemplatizedType(baseInterface)
				if !alreadySeenInterfaces.Add(baseInterface) {
					tc.Report(n, TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS, baseInterface.String())
				}
			}
		}
		// check properties
		tc.validator.ExpectAllInterfaceProperties(n, functionType)
		if !functionType.IsAbstract() {
			tc.validator.ExpectAbstractMethodsImplemented(n, functionType)
		}
	}
}

// NormalizeTemplatizedType normalizes `Foo<string>` to just `Foo` and is the
// identify function for other types.
//
// Java: TypeCheck.normalizeTemplatizedType(JSType)
func (tc *TypeCheck) NormalizeTemplatizedType(maybeTemplatizedType rhino.JSType) rhino.JSType {
	if !maybeTemplatizedType.IsTemplatizedType() {
		return maybeTemplatizedType
	}
	return maybeTemplatizedType.ToMaybeTemplatizedType().GetRawType()
}

// Java: TypeCheck.checkStructDictSubtyping(Node, FunctionType, FunctionType)
func (tc *TypeCheck) CheckStructDictSubtyping(n rhino.Node, subtype rhino.FunctionType, supertype rhino.FunctionType) {
	if subtype.MakesDicts() && supertype.MakesStructs() {
		tc.compiler.Report(
			JSError_MakeWithNode(
				n, TypeCheck_DICT_EXTEND_STRUCT_TYPE, jutil.ToJavaString(subtype.GetDisplayName()), jutil.ToJavaString(supertype.GetDisplayName())))
	} else if subtype.MakesStructs() && supertype.MakesDicts() {
		tc.compiler.Report(
			JSError_MakeWithNode(
				n, TypeCheck_STRUCT_EXTEND_DICT_TYPE, jutil.ToJavaString(subtype.GetDisplayName()), jutil.ToJavaString(supertype.GetDisplayName())))
	}
}

// CheckInterface checks an interface, which may be either an ES5-style
// FUNCTION node, or a CLASS node.
//
// Java: TypeCheck.checkInterface(Node, FunctionType)
func (tc *TypeCheck) CheckInterface(n rhino.Node, functionType rhino.FunctionType) {
	// Interface must extend only interfaces
	for _, extInterface := range functionType.GetExtendedInterfaces() {
		if extInterface.GetConstructor() != nil && !extInterface.GetConstructor().IsInterface() {
			tc.compiler.Report(
				JSError_MakeWithNode(n, TypeCheck_CONFLICTING_EXTENDED_TYPE, "interface", tc.GetBestFunctionName(n)))
		}
	}

	// Check whether the extended interfaces have any conflicts
	if functionType.GetExtendedInterfacesCount() > 1 {
		// Only check when extending more than one interfaces
		properties := jutil.NewObjectLinkedHashMap[rhino.Property_Key, rhino.ObjectType]()
		currentProperties := jutil.NewObjectLinkedHashMap[rhino.Property_Key, rhino.ObjectType]()
		for _, interfaceType := range functionType.GetExtendedInterfaces() {
			currentProperties.Clear()
			tc.CheckInterfaceConflictProperties(
				n, tc.GetBestFunctionName(n), properties, currentProperties, interfaceType)
			for k, v := range currentProperties.All() {
				properties.Put(k, v)
			}
		}
	}

	loopPath := functionType.CheckExtendsLoop()
	if loopPath != nil {
		names := make([]string, len(loopPath))
		for i, fn := range loopPath {
			names[i] = jutil.ToJavaString(fn.GetDisplayName())
		}
		strPath := strings.Join(names, " -> ")
		tc.compiler.Report(
			JSError_MakeWithNode(n, TypeCheck_INTERFACE_EXTENDS_LOOP, jutil.ToJavaString(loopPath[0].GetDisplayName()), strPath))
	}

	tc.validator.ExpectAllInterfaceProperties(n, functionType)
}

// Java: TypeCheck.getBestFunctionName(Node)
func (tc *TypeCheck) GetBestFunctionName(n rhino.Node) string {
	jutil.CheckState(n.IsClass() || n.IsFunction())
	name := NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
	if name != nil {
		return *name
	}
	return "<anonymous@" + jutil.ToJavaString(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno()) + ">"
}

var TypeCheck_GOOG_INHERITS = rhino.QualifiedName_Of("goog.inherits")

// CheckCallConventions validates class-defining calls. Because JS has no
// 'native' syntax for defining classes, we need to do this manually.
//
// Java: TypeCheck.checkCallConventions(NodeTraversal, Node)
func (tc *TypeCheck) CheckCallConventions(t *NodeTraversal, call ast.NormalOrOptChainCall) {
	n := call.Node()
	relationship := tc.compiler.GetCodingConvention().GetClassesDefinedByCall(n)
	scope := t.GetTypedScope()
	if relationship != nil {
		superClass := rhino.StaticTypedScope_LookupQualifiedNameString(scope, relationship.superclassName)
		superClassInstanceType := TypeValidator_GetInstanceOfCtor(superClass)
		subClass := rhino.StaticTypedScope_LookupQualifiedNameString(scope, relationship.subclassName)
		subClassInstance := TypeValidator_GetInstanceOfCtor(subClass)
		if relationship.type_ == CodingConvention_SubclassType_INHERITS &&
			superClassInstanceType != nil &&
			!superClassInstanceType.IsEmptyType() &&
			subClassInstance != nil &&
			!subClassInstance.IsEmptyType() {
			if callee := call.Callee().Node(); callee.IsQualifiedName() &&
				TypeCheck_GOOG_INHERITS.Matches(callee) &&
				subClass.ToMaybeFunctionType() != nil &&
				!subClass.ToMaybeFunctionType().GetSource().IsNil() &&
				subClass.ToMaybeFunctionType().GetSource().IsClass() {
				tc.compiler.Report(JSError_MakeWithNode(n, TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS))
			}
			tc.validator.ExpectSuperType(n, superClassInstanceType, subClassInstance)
		}
	}
}

// VisitCall visits a CALL node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitCall(NodeTraversal, Node)
func (tc *TypeCheck) VisitCall(t *NodeTraversal, call ast.NormalOrOptChainCall) {
	n := call.Node()
	tc.CheckCallConventions(t, call)

	callee := call.Callee()
	childType := tc.GetJSType(callee.Node()).RestrictByNotNullOrUndefined()

	if !childType.CanBeCalled() {
		tc.Report(n, TypeCheck_NOT_CALLABLE, childType.String())
		tc.EnsureTyped(n)
		return
	}

	// A couple of types can be called as if they were functions.
	// If it is a function type, then validate parameters.
	if childType.IsFunctionType() {
		functionType := childType.ToMaybeFunctionType()

		// Non-native constructors should not be called directly
		// unless they specify a return type
		if functionType.IsConstructor() &&
			!functionType.IsNativeObjectType() &&
			(functionType.GetReturnType().IsUnknownType() ||
				functionType.GetReturnType().IsVoidType()) &&
			!callee.Node().IsSuper() {
			tc.Report(n, TypeCheck_CONSTRUCTOR_NOT_CALLABLE, childType.String())
		}

		// Functions with explicit 'this' types must be called in a GETPROP or GETELEM.
		_, calleeIsGet := callee.AsNormalOrOptChainGet()
		if functionType.IsOrdinaryFunction() && !calleeIsGet {
			receiverType := functionType.GetTypeOfThis()
			if receiverType.IsUnknownType() ||
				receiverType.IsAllType() ||
				receiverType.IsVoidType() ||
				(receiverType.IsObjectType() && receiverType.ToObjectType().IsNativeObjectType()) {
				// Allow these special cases.
			} else {
				tc.Report(n, TypeCheck_EXPECTED_THIS_TYPE, functionType.String())
			}
		}

		tc.VisitArgumentList(call.AsCallOrNew(), functionType)
		tc.EnsureTypedWithType(n, functionType.GetReturnType())
	} else {
		tc.EnsureTyped(n)
	}

	// TODO(nicksantos): Add something to check for calls of RegExp objects,
	// which is not supported by IE. Either say something about the return type
	// or warn about the non-portability of the call or both.
}

// VisitArgumentList visits the parameters of a CALL or a NEW node.
//
// Java: TypeCheck.visitArgumentList(Node, FunctionType)
func (tc *TypeCheck) VisitArgumentList(call ast.CallOrNew, functionType rhino.FunctionType) {
	tc.CheckArgumentsMatchParameters(call.Node(), functionType, functionType.GetParameters(), 0)
}

// CheckArgumentsMatchParameters checks that a list of arguments match a list
// of formal parameters.
//
// If given a TAGGED_TEMPLATE_LIT, the given Iterator should only contain the
// parameters corresponding to the actual template lit sub arguments, skipping
// over the first parameter.
//
// firstParameterIndex is the index of the first parameter in the given
// Iterator in the function type's parameter list.
//
// Go-only: the Java gets the arguments as an Iterator over
// NodeUtil.getInvocationArgsAsIterable(call); here they are read from call,
// a CallOrNew (its arguments) or a TaggedTemplateLit (the substitutions of
// its template), and checkArgument is the body of the Java loop. The
// parameters are a slice (the elements the Java Iterator has left to
// yield).
//
// Java: TypeCheck.checkArgumentsMatchParameters(Node, FunctionType, Iterator<Node>, Iterator<Parameter>, int)
func (tc *TypeCheck) CheckArgumentsMatchParameters(
	call rhino.Node,
	functionType rhino.FunctionType,
	parameters []*rhino.FunctionType_Parameter,
	firstParameterIndex int) {

	spreadArgumentCount := 0
	normalArgumentCount := firstParameterIndex
	checkArgumentTypeAgainstParameter := true
	var parameter *rhino.FunctionType_Parameter
	nextParameter := 0
	checkArgument := func(argument ast.Arg) {
		// Count normal & spread arguments.
		if _, isSpread := argument.AsIterSpread(); isSpread {
			// we have some form of this case
			// someCall(arg1, arg2, ...firstSpreadExpression, argN, ...secondSpreadExpression)
			spreadArgumentCount++
			// Once we see a spread parameter, we can no longer match up arguments with parameters.
			checkArgumentTypeAgainstParameter = false
		} else {
			normalArgumentCount++
		}

		// Get the next parameter, if we're still matching parameters and arguments.
		if checkArgumentTypeAgainstParameter {
			if nextParameter < len(parameters) {
				parameter = parameters[nextParameter]
				nextParameter++
			} else if parameter != nil && parameter.IsVariadic() {
				// use varargs for all remaining parameters
			} else {
				// else we ran out of parameters and will report that after this loop
				parameter = nil
				checkArgumentTypeAgainstParameter = false
			}
		}

		if checkArgumentTypeAgainstParameter {
			tc.validator.ExpectArgumentMatchesParameter(
				argument.Node(), tc.GetJSType(argument.Node()), parameter.GetJSType(), call, normalArgumentCount)
		}
	}

	var callee ast.Expr
	if tagged, ok := ast.AsTaggedTemplateLit(call); ok {
		callee = tagged.Tag()
		for part := range tagged.Template().Parts() {
			if sub, ok := part.AsTemplateLitSub(); ok {
				checkArgument(sub.Expr().AsArg())
			}
		}
	} else {
		invocation := ast.MustCallOrNew(call)
		callee = invocation.Callee()
		for argument := range invocation.Args() {
			checkArgument(argument)
		}
	}

	minArity := functionType.GetMinArity()
	maxArity := functionType.GetMaxArity()

	if spreadArgumentCount > 0 {
		if normalArgumentCount > maxArity {
			// We cannot reliably check whether the total argument count is wrong, but we can at
			// least tell if there are more arguments than the function can handle even ignoring the
			// spreads.
			maxArityText := ""
			if maxArity != math.MaxInt32 {
				maxArityText = " and no more than " + strconv.Itoa(maxArity) + " argument(s)"
			}
			tc.Report(
				call,
				TypeCheck_WRONG_ARGUMENT_COUNT,
				tc.typeRegistry.GetReadableTypeNameNoDeref(callee.Node()),
				"at least "+strconv.Itoa(normalArgumentCount),
				strconv.Itoa(minArity),
				maxArityText)
		}
	} else {
		if minArity > normalArgumentCount || maxArity < normalArgumentCount {
			maxArityText := ""
			if maxArity != math.MaxInt32 {
				maxArityText = " and no more than " + strconv.Itoa(maxArity) + " argument(s)"
			}
			tc.Report(
				call,
				TypeCheck_WRONG_ARGUMENT_COUNT,
				tc.typeRegistry.GetReadableTypeNameNoDeref(callee.Node()),
				strconv.Itoa(normalArgumentCount),
				strconv.Itoa(minArity),
				maxArityText)
		}
	}
}

// VisitImplicitReturnExpression visits an arrow function expression body.
//
// Java: TypeCheck.visitImplicitReturnExpression(NodeTraversal, Node)
func (tc *TypeCheck) VisitImplicitReturnExpression(t *NodeTraversal, exprNode rhino.Node) {
	enclosingFunction := t.GetEnclosingFunction()
	jsType := tc.GetJSType(enclosingFunction)
	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()

		expectedReturnType := functionType.GetReturnType()
		// if no return type is specified, undefined must be returned
		// (it's a void function)
		if expectedReturnType == nil {
			expectedReturnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		} else if enclosingFunction.IsAsyncFunction() {
			// Unwrap the async function's declared return type.
			expectedReturnType = Promises_CreateAsyncReturnableType(tc.typeRegistry, expectedReturnType)
		}

		// Fetch the returned value's type
		actualReturnType := tc.GetJSType(exprNode)

		tc.validator.ExpectCanAssignTo(
			exprNode, actualReturnType, expectedReturnType, "inconsistent return type")
	}
}

// VisitReturn visits a RETURN node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitReturn(NodeTraversal, Node)
func (tc *TypeCheck) VisitReturn(t *NodeTraversal, ret ast.Return) {
	enclosingFunction := t.GetEnclosingFunction()

	jsType := tc.GetJSType(enclosingFunction)

	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()

		returnType := functionType.GetReturnType()
		// if no return type is specified, undefined must be returned
		// (it's a void function)
		if returnType == nil {
			returnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		} else if enclosingFunction.IsGeneratorFunction() {
			// Unwrap the template variable from a generator function's declared return type.
			// e.g. if returnType is "Generator<string, number, void>", the generator should return
			// "number".
			returnType = JsIterables_GetReturnElementType(returnType, tc.typeRegistry)

			if enclosingFunction.IsAsyncGeneratorFunction() {
				// Can return x|IThenable<x> in an AsyncGenerator<x>, no await needed. Note that we must
				// first wrap the type in IThenable as createAsyncReturnableType will map a non-IThenable
				// to `?`.
				returnType =
					Promises_CreateAsyncReturnableType(
						tc.typeRegistry, Promises_WrapInIThenable(tc.typeRegistry, returnType))
			}
		} else if enclosingFunction.IsAsyncFunction() {
			// e.g. `!Promise<string>` => `string|!IThenable<string>`
			// We transform the expected return type rather than the actual return type so that the
			// extual return type is always reported to the user. This was felt to be clearer.
			returnType = Promises_CreateAsyncReturnableType(tc.typeRegistry, returnType)
		} else if returnType.IsVoidType() && functionType.IsConstructor() {
			// Allow constructors to use empty returns for flow control.
			if _, hasValue := ret.Value(); !hasValue {
				return
			}

			// Allow constructors to return its own instance type
			returnType = functionType.GetInstanceType()
		}

		// fetching the returned value's type
		var valueNode rhino.Node
		var actualReturnType rhino.JSType
		if value, ok := ret.Value(); ok {
			valueNode = value.Node()
			actualReturnType = tc.GetJSType(valueNode)
		} else {
			actualReturnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
			valueNode = ret.Node()
		}

		// verifying
		tc.validator.ExpectCanAssignTo(
			valueNode, actualReturnType, returnType, "inconsistent return type")
	}
}

// VisitYield visits a YIELD node.
//
// Java: TypeCheck.visitYield(NodeTraversal, Node)
func (tc *TypeCheck) VisitYield(t *NodeTraversal, yield ast.Yield) {
	n := yield.Node()
	jsType := tc.GetJSType(t.GetEnclosingFunction())

	declaredYieldType := tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()
		returnType := functionType.GetReturnType()
		declaredYieldType = JsIterables_GetElementType(returnType, tc.typeRegistry)

		if t.GetEnclosingFunction().IsAsyncGeneratorFunction() {
			// Can yield x|IThenable<x> in an AsyncGenerator<x>, no await needed. Note that we must
			// first wrap the type in IThenable as createAsyncReturnableType will map a non-IThenable to
			// `?`.
			declaredYieldType =
				Promises_CreateAsyncReturnableType(
					tc.typeRegistry, Promises_WrapInIThenable(tc.typeRegistry, declaredYieldType))
		}
	}

	// fetching the yielded value's type
	var valueNode rhino.Node
	var actualYieldType rhino.JSType
	if value, ok := yield.Value(); ok {
		valueNode = value.Node()
		actualYieldType = tc.GetJSType(valueNode)
	} else {
		actualYieldType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		valueNode = n
	}

	if n.IsYieldAll() {
		if t.GetEnclosingFunction().IsAsyncGeneratorFunction() {
			maybeActualYieldType :=
				tc.validator.ExpectAutoboxesToIterableOrAsyncIterable(
					n, actualYieldType, "Expression yield* expects an iterable or async iterable")
			if maybeActualYieldType == nil {
				// don't do any further typechecking of the yield* type.
				return
			}
			actualYieldType = maybeActualYieldType
		} else {
			if !tc.validator.ExpectAutoboxesToIterable(
				n, actualYieldType, "Expression yield* expects an iterable") {
				// don't do any further typechecking of the yield* type.
				return
			}
			templateType := tc.typeRegistry.GetIterableValueTemplate()
			actualYieldType =
				actualYieldType.Autobox().GetTemplateTypeMap().GetResolvedTemplateType(templateType)
		}
	}

	// verifying
	tc.validator.ExpectCanAssignTo(
		valueNode,
		actualYieldType,
		declaredYieldType,
		"Yielded type does not match declared return type.")
}

// Java: TypeCheck.visitTaggedTemplateLit(Node)
func (tc *TypeCheck) VisitTaggedTemplateLit(tagged ast.TaggedTemplateLit) {
	n := tagged.Node()
	tag := tagged.Tag().Node()
	tagType := tag.GetJSType().RestrictByNotNullOrUndefined()

	if !tagType.CanBeCalled() {
		tc.Report(n, TypeCheck_NOT_CALLABLE, tagType.String())
		return
	} else if !tagType.IsFunctionType() {
		// A few types, like the unknown, regexp, and bottom types, can be called as if they are
		// functions. Return if we have one of those types that is not actually a known function.
		return
	}

	tagFnType := tagType.ToMaybeFunctionType()
	parameters := tagFnType.GetParameters()

	// The tag function gets an array of all the template lit substitutions as its first argument,
	// but there's no actual AST node representing that array so we typecheck it separately from
	// the other tag arguments.

	// Validate that the tag function takes at least one parameter
	if len(parameters) == 0 {
		tc.Report(
			n,
			TypeCheck_WRONG_ARGUMENT_COUNT,
			tc.typeRegistry.GetReadableTypeNameNoDeref(tag),
			strconv.Itoa(NodeUtil_GetInvocationArgsCount(n)),
			"0",
			" and no more than 0 argument(s)")
		return
	}

	// Validate that the first parameter is a supertype of ITemplateArray
	firstParameter := parameters[0]
	parameters = parameters[1:]
	parameterType := firstParameter.GetJSType().RestrictByNotNullOrUndefined()
	if parameterType != nil {
		tc.validator.ExpectITemplateArraySupertype(
			tag, parameterType, "Invalid type for the first parameter of tag function")
	}

	// Validate the remaining parameters (the template literal substitutions)
	tc.CheckArgumentsMatchParameters(n, tagFnType, parameters, 1)
}

// VisitBinaryOperator unifies the type checking involved in the core binary
// operators and the corresponding assignment operators. The representation
// used internally is such that common code can handle both kinds of
// operators easily.
//
// op is the operator; n the node being checked.
//
// Go-only: n is an ast.Binary or an ast.AssignOp, which share no view, so
// the caller passes the operands: left and right, or target and value.
//
// Java: TypeCheck.visitBinaryOperator(Token, Node)
func (tc *TypeCheck) VisitBinaryOperator(op rhino.Token, n rhino.Node, left rhino.Node, right rhino.Node) {
	operatorType := tc.GetJSType(n)
	leftType := tc.GetJSType(left)
	rightType := tc.GetJSType(right)
	if operatorType.IsNoType() {
		// An attempt to mix bigint with other types
		tc.Report(
			n,
			TypeCheck_BINARY_OPERATION,
			jutil.ToJavaString(NodeUtil_OpToStr(n.GetToken())),
			leftType.String(),
			rightType.String())
		return
	}
	switch op {
	case rhino.Token_ASSIGN_LSH, rhino.Token_ASSIGN_RSH, rhino.Token_LSH, rhino.Token_RSH, rhino.Token_ASSIGN_URSH,
		rhino.Token_URSH:
		if operatorType.IsNumber() {
			// TypeInference set the operator type to 'number', so we know bigint isn't involved.
			// NOTE: >>> and >>>= aren't valid operations for bigint, so TypeInference ignores the
			//     operand types for it and always sets the operator's type to 'number'. Thus,
			//     if one of the operands is a bigint, we'll end up reporting that as an error here.
			opStr := jutil.ToJavaString(NodeUtil_OpToStr(n.GetToken()))
			if !leftType.MatchesNumberContext() {
				tc.Report(left, TypeCheck_BIT_OPERATION, opStr, leftType.String())
			} else {
				tc.validator.ExpectNumberStrict(n, leftType, "operator "+opStr)
			}
			if !rightType.MatchesNumberContext() {
				tc.Report(right, TypeCheck_BIT_OPERATION, opStr, rightType.String())
			} else {
				tc.validator.ExpectNumberStrict(n, rightType, "operator "+opStr)
			}
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "left operand")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "right operand")
		}
	case rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_DIV,
		rhino.Token_MOD,
		rhino.Token_MUL,
		rhino.Token_SUB,
		rhino.Token_EXPONENT:
		if operatorType.IsNumber() {
			// TypeInference set the operator type to 'number', so we know bigint isn't involved.
			tc.validator.ExpectNumber(left, leftType, "left operand")
			tc.validator.ExpectNumber(right, rightType, "right operand")
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "left operand")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "right operand")
		}
	case rhino.Token_ASSIGN_BITAND, rhino.Token_ASSIGN_BITXOR, rhino.Token_ASSIGN_BITOR, rhino.Token_BITAND,
		rhino.Token_BITXOR, rhino.Token_BITOR:
		if operatorType.IsNumber() {
			// This condition is meant to catch any old cases (where bigint isn't involved)
			tc.validator.ExpectBitwiseable(left, leftType, "bad left operand to bitwise operator")
			tc.validator.ExpectBitwiseable(right, rightType, "bad right operand to bitwise operator")
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "bad left operand to bitwise operator")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "bad right operand to bitwise operator")
		}
	case rhino.Token_ASSIGN_ADD, rhino.Token_ADD:
	default:
		tc.Report(n, TypeCheck_UNEXPECTED_TOKEN, op.String())
	}
	tc.EnsureTyped(n)
}

// typeCheck_inputProvider is the lambda `(b) -> t.getInput()`, a
// CompilerInputProvider.
//
// Go-only: language-adapter — a lambda implementing a Go interface.
type typeCheck_inputProvider struct {
	t *NodeTraversal
}

func (p typeCheck_inputProvider) GetInput(*rhino.InputId) *CompilerInput {
	return p.t.GetInput()
}

// VisitModuleBody validates the implicit assignment to the global for a
// legacy goog.module.
//
// Java: TypeCheck.visitModuleBody(NodeTraversal, Node)
func (tc *TypeCheck) VisitModuleBody(t *NodeTraversal, body ast.ModuleBody) {
	moduleBody := body.Node()
	associatedModule :=
		ModuleImportResolver_GetModuleFromScopeRoot(
			tc.compiler.GetModuleMap(), typeCheck_inputProvider{t}, moduleBody)
	if !associatedModule.Metadata().IsLegacyGoogModule() {
		return
	}
	moduleName := rhino.QualifiedName_Of(*associatedModule.ClosureNamespace())
	// The first statement of a legacy goog.module is its goog.module call.
	googModuleCall, ok := body.FirstStatement()
	jutil.CheckState(ok)
	if moduleName.IsSimple() {
		globalVar := tc.topScope.GetVar(moduleName.Join())
		tc.validator.ExpectCanAssignTo(
			googModuleCall.Node(), moduleBody.GetJSType(), globalVar.GetType(), "legacy goog.module export")
	} else {
		parentType := tc.topScope.LookupQualifiedName(moduleName.GetOwner())
		var parentObjectType rhino.ObjectType
		if parentType != nil {
			parentObjectType = parentType.ToMaybeObjectType()
		}
		if parentObjectType == nil {
			return
		}
		tc.validator.ExpectCanAssignToPropertyOfOwnerType(
			googModuleCall.Node(),
			moduleBody.GetJSType(),
			parentObjectType.GetPropertyType(moduleName.GetComponent()),
			parentObjectType,
			func() string { return moduleName.GetOwner().Join() },
			moduleName.GetComponent())
	}
}

// VisitDynamicImport validates that a dynamic import statement has a single
// child of type string.
//
// Java: TypeCheck.visitDynamicImport(NodeTraversal, Node)
func (tc *TypeCheck) VisitDynamicImport(t *NodeTraversal, dynamicImport ast.DynamicImport) {
	tc.EnsureTypedWithType(dynamicImport.Node(), tc.GetNativeType(rhino.JSTypeNative_PROMISE_TYPE))

	importSpecifier := dynamicImport.Specifier().Node()
	importSpecifierType := importSpecifier.GetJSType()
	if importSpecifierType == nil {
		tc.EnsureTypedWithNative(importSpecifier, rhino.JSTypeNative_STRING_TYPE)
	} else {
		tc.validator.ExpectNotNullOrUndefined(
			t,
			importSpecifier,
			importSpecifierType,
			"dynamic import specifier",
			tc.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	}
}

// CheckEnumAlias checks enum aliases.
//
// We verify that the enum element type of the enum used for initialization
// is a subtype of the enum element type of the enum the value is being copied
// in.
//
// Example:
//
//	var myEnum = myOtherEnum;
//
// Enum aliases are irregular, so we need special code for this :(
//
// valueType is the type of the value used for initialization of the enum;
// nodeToWarn the node on which to issue warnings on.
//
// Java: TypeCheck.checkEnumAlias(NodeTraversal, JSDocInfo, JSType, Node)
func (tc *TypeCheck) CheckEnumAlias(
	t *NodeTraversal, declInfo *rhino.JSDocInfo, valueType rhino.JSType, nodeToWarn rhino.Node) {
	if declInfo == nil || !declInfo.HasEnumParameterType() {
		return
	}

	if !valueType.IsEnumType() {
		return
	}

	valueEnumType := valueType.ToMaybeEnumType()
	valueEnumPrimitiveType := valueEnumType.GetElementsType().GetPrimitiveType()
	tc.validator.ExpectCanAssignTo(
		nodeToWarn,
		valueEnumPrimitiveType,
		declInfo.GetEnumParameterType().Evaluate(t.GetTypedScope(), tc.typeRegistry),
		"incompatible enum element types")
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

// CheckJsdocInfoContainsObjectWithBadKey checks if current node contains js
// docs and checks all types specified in the js doc whether they have Objects
// with potentially invalid keys. For example: `Object<!Object, number>`. If
// such type is found, a warning is reported for the current node.
//
// Java: TypeCheck.checkJsdocInfoContainsObjectWithBadKey(Node)
func (tc *TypeCheck) CheckJsdocInfoContainsObjectWithBadKey(n rhino.Node) {
	if info := n.GetJSDocInfo(); info != nil {
		tc.checkJsdocInfoContainsObjectWithBadKeySlow(n, info)
	}
}

func (tc *TypeCheck) checkJsdocInfoContainsObjectWithBadKeySlow(n rhino.Node, info *rhino.JSDocInfo) {
	tc.CheckTypeContainsObjectWithBadKey(n, info.GetType())
	tc.CheckTypeContainsObjectWithBadKey(n, info.GetReturnType())
	tc.CheckTypeContainsObjectWithBadKey(n, info.GetTypedefType())
	paramCount := info.GetParameterCount()
	for i := 0; i < paramCount; i++ {
		tc.CheckTypeContainsObjectWithBadKey(n, info.GetParameterTypeAt(i))
	}
}

// Java: TypeCheck.checkTypeContainsObjectWithBadKey(Node, JSTypeExpression)
func (tc *TypeCheck) CheckTypeContainsObjectWithBadKey(n rhino.Node, type_ *rhino.JSTypeExpression) {
	if type_ != nil && type_.GetRoot().GetJSType() != nil {
		realType := type_.GetRoot().GetJSType()
		if !realType.IsTemplatizedType() && !realType.IsUnionType() && !realType.IsOrdinaryFunction() {
			return
		}
		if tc.badKeyCheckedTypes == nil {
			tc.badKeyCheckedTypes = jutil.NewObjectLinkedHashSet[rhino.JSType]()
		} else {
			tc.badKeyCheckedTypes.Clear()
		}
		objectWithBadKey :=
			tc.FindObjectWithNonStringifiableKey(realType, tc.badKeyCheckedTypes)
		tc.badKeyCheckedTypes.Clear()
		if objectWithBadKey != nil {
			tc.compiler.Report(JSError_MakeWithNode(n, TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY, objectWithBadKey.String()))
		}
	}
}

// IsReasonableObjectPropertyKey checks whether type is useful as the key of
// an object property access. This means it should be either stringifiable or
// a symbol. Stringifiable types are types that can be converted to string and
// give unique results for different objects. For example objects have native
// toString() method that on chrome returns "[object Object]" for all objects
// making it useless when used as keys. At the same time native types like
// numbers can be safely converted to strings and used as keys. Also user
// might have provided custom toString() methods for a class making it
// suitable for using as key.
//
// Java: TypeCheck.isReasonableObjectPropertyKey(JSType)
func (tc *TypeCheck) IsReasonableObjectPropertyKey(type_ rhino.JSType) bool {
	// Check built-in types
	if type_.IsUnknownType() ||
		type_.IsNumber() ||
		type_.IsString() ||
		type_.IsSymbol() ||
		type_.IsBooleanObjectType() ||
		type_.IsBooleanValueType() ||
		type_.IsDateType() ||
		type_.IsRegexpType() ||
		type_.IsInterface() ||
		type_.IsRecordType() ||
		type_.IsNullType() ||
		type_.IsVoidType() {
		return true
	}

	// For enums check that underlying type is stringifiable.
	if type_.ToMaybeEnumElementType() != nil {
		return tc.IsReasonableObjectPropertyKey(type_.ToMaybeEnumElementType().GetPrimitiveType())
	}

	// Array is stringifiable if it doesn't have template type or if it does have it, the template
	// type must be also stringifiable.
	// Good: Array, Array.<number>
	// Bad: Array.<!Object>
	if type_.IsArrayType() {
		return true
	}
	if type_.IsTemplatizedType() {
		templatizedType := type_.ToMaybeTemplatizedType()
		if templatizedType.GetReferencedType().IsArrayType() {
			return tc.IsReasonableObjectPropertyKey(templatizedType.GetTemplateTypes()[0])
		}
	}

	// Named types are usually @typedefs. For such types we need to check underlying type specified
	// in @typedef annotation.
	if namedType, ok := type_.(*rhino.NamedType); ok {
		return tc.IsReasonableObjectPropertyKey(namedType.GetReferencedType())
	}

	// For union type every alternate must be stringifiable.
	if type_.IsUnionType() {
		for _, alternateType := range type_.ToMaybeUnionType().GetAlternates() {
			if !tc.IsReasonableObjectPropertyKey(alternateType) {
				return false
			}
		}
		return true
	}

	// Handle interfaces and classes.
	if type_.IsObject() {
		objectType := type_.ToMaybeObjectType()
		constructor := objectType.GetConstructor()
		// Interfaces considered stringifiable as user might implement toString() method in
		// classes-implementations.
		if constructor != nil && constructor.IsInterface() {
			return true
		}
		// This is user-defined class so check if it has custom toString() method.
		return tc.ClassHasToString(objectType)
	}
	return false
}

// IsObjectTypeWithNonStringifiableKey checks whether current type is Object
// type with non-stringifable key.
//
// Java: TypeCheck.isObjectTypeWithNonStringifiableKey(JSType)
func (tc *TypeCheck) IsObjectTypeWithNonStringifiableKey(type_ rhino.JSType) bool {
	if !type_.IsTemplatizedType() {
		// TODO(nickreid): Why don't we care about types like `Foo extends Object<Bar, Qux>`?
		return false
	}

	templateTypeMap := type_.GetTemplateTypeMap()
	objectIndexKey := tc.typeRegistry.GetObjectIndexKey()
	if templateTypeMap.HasTemplateKey(objectIndexKey) {
		return !tc.IsReasonableObjectPropertyKey(
			templateTypeMap.GetResolvedTemplateType(objectIndexKey))
	} else {
		return false
	}
}

// FindObjectWithNonStringifiableKey checks whether type (or one of its
// component if is composed type like union or templatized type) has Object
// with non-stringifiable key. For example `Object.<!Object, number>`.
//
// Returns non-stringifiable type which is used as key or null if all there
// are no such types.
//
// Java: TypeCheck.findObjectWithNonStringifiableKey(JSType, Set<JSType>)
func (tc *TypeCheck) FindObjectWithNonStringifiableKey(
	type_ rhino.JSType, alreadyCheckedTypes *jutil.ObjectLinkedHashSet[rhino.JSType]) rhino.JSType {
	if !type_.IsTemplatizedType() && !type_.IsUnionType() && !type_.IsOrdinaryFunction() {
		return nil
	}
	if !alreadyCheckedTypes.Add(type_) {
		// This can happen in recursive types. Current type already being checked earlier in
		// stacktrace so now we just skip it.
		return nil
	}
	if tc.IsObjectTypeWithNonStringifiableKey(type_) {
		return type_
	}
	if type_.IsUnionType() {
		for _, alternateType := range type_.ToMaybeUnionType().GetAlternates() {
			result := tc.FindObjectWithNonStringifiableKey(alternateType, alreadyCheckedTypes)
			if result != nil {
				return result
			}
		}
	}
	if type_.IsTemplatizedType() {
		for _, templateType := range type_.ToMaybeTemplatizedType().GetTemplateTypes() {
			result := tc.FindObjectWithNonStringifiableKey(templateType, alreadyCheckedTypes)
			if result != nil {
				return result
			}
		}
	}
	if type_.IsOrdinaryFunction() {
		function := type_.ToMaybeFunctionType()
		for _, parameter := range function.GetParameters() {
			result :=
				tc.FindObjectWithNonStringifiableKey(parameter.GetJSType(), alreadyCheckedTypes)
			if result != nil {
				return result
			}
		}
		return tc.FindObjectWithNonStringifiableKey(function.GetReturnType(), alreadyCheckedTypes)
	}
	return nil
}

// ClassHasToString checks whether class has overridden toString() method.
// All objects has native toString() method but we ignore it as it is not
// useful so we need user-provided toString() method.
//
// Java: TypeCheck.classHasToString(ObjectType)
func (tc *TypeCheck) ClassHasToString(type_ rhino.ObjectType) bool {
	toStringProperty := type_.GetOwnSlot("toString")
	if toStringProperty != nil {
		return toStringProperty.GetType().IsFunctionType()
	}
	parent := type_.GetImplicitPrototype()
	if parent != nil && !parent.IsNativeObjectType() {
		return tc.ClassHasToString(parent)
	}
	return false
}

// CheckNotReadonlyPropertyAssignment: given the LHS of a property or element
// assignment, checks that the type that we're assigning into is not readonly
// (ReadonlyArray, in particular).
//
// This is basically unsound since it checks just "ReadonlyArray" so casting
// up to "Iterable" or down to "Array" will defeat it; but it should catch
// basic errors.
//
// Java: TypeCheck.checkNotReadonlyPropertyAssignment(Node)
func (tc *TypeCheck) CheckNotReadonlyPropertyAssignment(lhs ast.Target) {
	// We only care about element or property assignments.
	get, ok := lhs.AsNormalGet()
	if !ok {
		return
	}

	// If we do not have type information, drop out.
	lhsType := get.Object().Node().GetJSType()
	if lhsType == nil {
		return
	}

	// We could be a reference to "ReadonlyArray" or a templatized wrapper.
	roArray := tc.GetNativeType(rhino.JSTypeNative_READONLY_ARRAY_TYPE)
	if !lhsType.IsUnionType() {
		if roArray.Equals(lhsType) || roArray.Equals(TypeCheck_MaybeReferencedType(lhsType)) {
			tc.compiler.Report(JSError_MakeWithNode(get.Node(), TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE, lhsType.String()))
		}
		return
	}
	alternates := TypeCheck_FlattenUnion(lhsType)
	for i := 0; i < len(alternates); i++ {
		type_ := alternates[i]
		if roArray.Equals(type_) || roArray.Equals(TypeCheck_MaybeReferencedType(type_)) {
			tc.compiler.Report(JSError_MakeWithNode(get.Node(), TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE, type_.String()))
			break
		}
	}
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
	outer *TypeCheck

	// Output log file for @logTypeInCompiler.
	typeLogFile diagnostic.LogFile

	// The node where we started logging, if we're logging.
	parentNode rhino.Node
}

// Suppress MustBeClosedChecker because `.close` is called by our `close` method.
//
// Java: TypeCheck.DebugTypeLogger()
func NewTypeCheck_DebugTypeLogger(outer *TypeCheck) *TypeCheck_DebugTypeLogger {
	l := &TypeCheck_DebugTypeLogger{outer: outer}
	if outer.compiler.IsDebugLoggingEnabled() {
		l.typeLogFile = outer.compiler.CreateOrReopenLog("DebugTypeLogger", "types_logged_in_compiler.log")
	}
	return l
}

// Java: TypeCheck.DebugTypeLogger.maybeStartLoggingAt(Node)
func (l *TypeCheck_DebugTypeLogger) MaybeStartLoggingAt(n rhino.Node) {
	// We don't log outside debug mode, or if we're already logging.
	if l.typeLogFile == nil || !l.parentNode.IsNil() {
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
	if l.typeLogFile != nil {
		l.typeLogFile.Close()
	}
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
