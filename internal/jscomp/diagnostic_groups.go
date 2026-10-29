// Port of com/google/javascript/jscomp/DiagnosticGroups.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// DiagnosticGroups holds the named groups of DiagnosticTypes exposed by
// Compiler.
type DiagnosticGroups struct{}

var DiagnosticGroups_UNUSED = DiagnosticType_Warning("JSC_UNUSED", "{0}")

type diagnosticGroupNameSet map[string]struct{}

func (s diagnosticGroupNameSet) Contains(name string) bool {
	_, ok := s[name]
	return ok
}

func (s diagnosticGroupNameSet) Size() int {
	return len(s)
}

var DiagnosticGroups_wildcardExcludedGroups = diagnosticGroupNameSet{
	"reportUnknownTypes":                  {},
	"analyzerChecks":                      {},
	"missingSourcesWarnings":              {},
	"closureUnawareCodeAnnotationPresent": {},
}

// Java: DiagnosticGroups()
func NewDiagnosticGroups() *DiagnosticGroups {
	return &DiagnosticGroups{}
}

var DiagnosticGroups_groupsByName = jutil.NewLinkedHashMap[string, *DiagnosticGroup]()

// Java: DiagnosticGroups.registerDeprecatedGroup(String)
func DiagnosticGroups_RegisterDeprecatedGroup(name string) *DiagnosticGroup {
	return DiagnosticGroups_RegisterGroup(name, NewDiagnosticGroupWithName(&name, DiagnosticGroups_UNUSED))
}

// DiagnosticGroups_RegisterUnsuppressibleGroup creates a group that is
// unsuppressible via the command line or in code.
//
// The resulting group is also undocumented.
//
// Java: DiagnosticGroups.registerUnsuppressibleGroup(DiagnosticType...)
func DiagnosticGroups_RegisterUnsuppressibleGroup(types ...*DiagnosticType) *DiagnosticGroup {
	return NewDiagnosticGroup(types...)
}

// DiagnosticGroups_RegisterGroup registers group itself under name (Java
// picks this overload for a single DiagnosticGroup argument).
//
// Java: DiagnosticGroups.registerGroup(String, DiagnosticGroup)
func DiagnosticGroups_RegisterGroup(name string, group *DiagnosticGroup) *DiagnosticGroup {
	DiagnosticGroups_groupsByName.Put(name, group)
	return group
}

// Java: DiagnosticGroups.registerGroup(String, DiagnosticType...)
func DiagnosticGroups_RegisterGroupOfTypes(name string, types ...*DiagnosticType) *DiagnosticGroup {
	group := NewDiagnosticGroupWithName(&name, types...)
	DiagnosticGroups_groupsByName.Put(name, group)
	return group
}

// Java: DiagnosticGroups.registerGroup(String, DiagnosticGroup...)
func DiagnosticGroups_RegisterGroupOfGroups(name string, groups ...*DiagnosticGroup) *DiagnosticGroup {
	group := NewDiagnosticGroupWithNameFromGroups(&name, groups...)
	DiagnosticGroups_groupsByName.Put(name, group)
	return group
}

// DiagnosticGroups_GetRegisteredGroups gets the registered diagnostic
// groups, indexed by name.
//
// Java: DiagnosticGroups.getRegisteredGroups()
func DiagnosticGroups_GetRegisteredGroups() *jutil.LinkedHashMap[string, *DiagnosticGroup] {
	return jutil.NewLinkedHashMapFrom(DiagnosticGroups_groupsByName)
}

// DiagnosticGroups_ForName finds the diagnostic group registered under the
// given name, or nil.
//
// Java: DiagnosticGroups.forName(String)
func DiagnosticGroups_ForName(name string) *DiagnosticGroup {
	return DiagnosticGroups_groupsByName.Get(name)
}

// A bit of a hack to display the available groups on the command-line.
// New groups should be added to this list if they are public and should
// be listed on the command-line as an available option.
//
// If a group is suppressible on a per-file basis, it should be added
// to parsing/ParserConfig.properties
const DiagnosticGroups_DIAGNOSTIC_GROUP_NAMES = "accessControls, " +
	"checkPrototypalTypes, " +
	"checkRegExp, " +
	"checkTypes, " +
	"checkVars, " +
	"conformanceViolations, " +
	"const, " +
	"constantProperty, " +
	"deprecated, " +
	"deprecatedAnnotations, " +
	"defineWithoutGoogDefine, " +
	"duplicateMessage, " +
	"es5Strict, " +
	"externsValidation, " +
	"functionParams, " +
	"globalThis, " +
	"invalidCasts, " +
	"lintVarDeclarations, " +
	"misplacedTypeAnnotation, " +
	"missingCallToSuper, " +
	"missingOverride, " +
	"missingPolyfill, " +
	"missingProperties, " +
	"missingProvide, " +
	"missingRequire, " +
	"missingReturn, " +
	"missingSourcesWarnings, " +
	"moduleLoad, " +
	"moduleImport, " +
	"msgDescriptions, " +
	"nonStandardJsDocs, " +
	"partialAlias, " +
	"polymer, " +
	"reportUnknownTypes, " +
	"strictCheckTypes, " +
	"strictMissingProperties, " +
	"strictModuleDepCheck, " +
	"strictPrimitiveOperators, " +
	"suspiciousCode, " +
	"typeInvalidation, " +
	"undefinedVars, " +
	"underscore, " +
	"unknownDefines, " +
	"unusedLocalVariables, " +
	"uselessCode, " +
	"untranspilableFeatures," +
	"visibility, " +
	"weakModuleGet"

// The static DiagnosticGroup fields. They are assigned, in Java declaration
// order, by diagnosticGroups_clinit.
var (
	// TODO(b/123768968) remove this diagnostic group, do not allow this suppression. Instead the only
	// work around should be to raise the output language to a high enough level. We need to suppress
	// right now because we don't have any language output level higher than ES5.
	DiagnosticGroups_UNTRANSPILABLE_FEATURES *DiagnosticGroup

	DiagnosticGroups_MODULE_LOAD *DiagnosticGroup

	DiagnosticGroups_MODULE_IMPORT *DiagnosticGroup

	DiagnosticGroups_GLOBAL_THIS *DiagnosticGroup

	DiagnosticGroups_DEPRECATED *DiagnosticGroup

	DiagnosticGroups_UNDERSCORE *DiagnosticGroup

	DiagnosticGroups_VISIBILITY *DiagnosticGroup

	// TODO(tbreisacher): Deprecate this and keep just the "visibility" group.
	DiagnosticGroups_ACCESS_CONTROLS *DiagnosticGroup

	DiagnosticGroups_NON_STANDARD_JSDOC *DiagnosticGroup

	DiagnosticGroups_INVALID_CASTS *DiagnosticGroup

	DiagnosticGroups_STRICT_MODULE_DEP_CHECK *DiagnosticGroup

	DiagnosticGroups_VIOLATED_MODULE_DEP *DiagnosticGroup

	DiagnosticGroups_EXTERNS_VALIDATION *DiagnosticGroup

	DiagnosticGroups_UNKNOWN_DEFINES *DiagnosticGroup

	DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE *DiagnosticGroup

	DiagnosticGroups_TWEAKS *DiagnosticGroup

	DiagnosticGroups_MISSING_CALL_TO_SUPER *DiagnosticGroup

	DiagnosticGroups_MISSING_OVERRIDE *DiagnosticGroup

	DiagnosticGroups_MISSING_PROPERTIES *DiagnosticGroup

	DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES *DiagnosticGroup

	DiagnosticGroups_J2CL_CHECKS *DiagnosticGroup

	DiagnosticGroups_MISSING_RETURN *DiagnosticGroup

	DiagnosticGroups_UNDEFINED_VARIABLES *DiagnosticGroup

	DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT *DiagnosticGroup

	DiagnosticGroups_CHECK_REGEXP *DiagnosticGroup

	// NOTE(dimvar): it'd be nice to add TypedScopeCreator.ALL_DIAGNOSTICS here,
	// but we would first need to cleanup projects that would break because
	// they set --jscomp_error=checkTypes.
	DiagnosticGroups_CHECK_TYPES *DiagnosticGroup

	DiagnosticGroups_ES5_INHERITANCE_DIAGNOSTIC_GROUP *DiagnosticGroup

	DiagnosticGroups_CHECK_PROTOTYPAL_TYPES *DiagnosticGroup

	// This group exists for the J2CL team to suppress the associated diagnostics using Java code
	// rather than `@suppress` annotations.
	DiagnosticGroups_CHECK_STATIC_OVERRIDES *DiagnosticGroup

	DiagnosticGroups_TOO_MANY_TYPE_PARAMS *DiagnosticGroup

	DiagnosticGroups_STRICT_MISSING_PROPERTIES *DiagnosticGroup

	DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS *DiagnosticGroup

	DiagnosticGroups_STRICT_CHECK_TYPES *DiagnosticGroup

	DiagnosticGroups_REPORT_UNKNOWN_TYPES *DiagnosticGroup

	DiagnosticGroups_CHECK_VARIABLES *DiagnosticGroup

	DiagnosticGroups_CHECK_USELESS_CODE *DiagnosticGroup

	DiagnosticGroups_CONST *DiagnosticGroup

	DiagnosticGroups_CONSTANT_PROPERTY *DiagnosticGroup

	DiagnosticGroups_ACCESS_CONTROLS_CONST *DiagnosticGroup

	DiagnosticGroups_TYPE_INVALIDATION *DiagnosticGroup

	DiagnosticGroups_DUPLICATE_VARS *DiagnosticGroup

	DiagnosticGroups_ES5_STRICT *DiagnosticGroup

	DiagnosticGroups_WEAK_MODULE_GET *DiagnosticGroup

	DiagnosticGroups_MISSING_PROVIDE *DiagnosticGroup

	DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR *DiagnosticGroup

	// identical to UNRECOGNIZED_TYPE_ERROR, but allowed in @suppress tags
	DiagnosticGroups_DANGEROUS_UNRECOGNIZED_TYPE_ERROR *DiagnosticGroup

	DiagnosticGroups_MISSING_REQUIRE *DiagnosticGroup

	// A set of diagnostics expected when parsing and type checking partial programs. Useful for clutz
	// (tool that extracts TypeScript definitions from JS code).
	DiagnosticGroups_MISSING_SOURCES_WARNINGS *DiagnosticGroup

	DiagnosticGroups_EXTRA_REQUIRE *DiagnosticGroup

	DiagnosticGroups_DUPLICATE_MESSAGE *DiagnosticGroup

	DiagnosticGroups_MESSAGE_DESCRIPTIONS *DiagnosticGroup

	// Warnings that only apply to people who use MSG_ to denote messages. Note that this doesn't
	// include warnings about proper use of goog.getMsg
	DiagnosticGroups_MSG_CONVENTIONS *DiagnosticGroup

	DiagnosticGroups_MISPLACED_TYPE_ANNOTATION *DiagnosticGroup

	DiagnosticGroups_MISPLACED_MSG_ANNOTATION *DiagnosticGroup

	DiagnosticGroups_MISPLACED_SUPPRESS *DiagnosticGroup

	DiagnosticGroups_SUSPICIOUS_CODE *DiagnosticGroup

	DiagnosticGroups_FUNCTION_PARAMS *DiagnosticGroup

	DiagnosticGroups_DEPRECATED_ANNOTATIONS *DiagnosticGroup

	DiagnosticGroups_UNUSED_LOCAL_VARIABLE *DiagnosticGroup

	DiagnosticGroups_JSDOC_MISSING_TYPE *DiagnosticGroup

	DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES *DiagnosticGroup

	DiagnosticGroups_PARTIAL_ALIAS *DiagnosticGroup

	// This lint is given its own diagnostic group because it's harder to fix than other
	// lint errors, and we want to discourage users from doing a blanket @suppress {lintChecks}.
	// This is intentionally not public. It should not be enabled directly; instead enable lintChecks.
	DiagnosticGroups_USE_OF_GOOG_PROVIDE *DiagnosticGroup

	// This is intended to be used to suppress warnings in code that cannot for some reason be updated
	// to use `let` and `const` instead of `var`. It should not be enabled directly, instead enable
	// "lintChecks".
	DiagnosticGroups_LINT_VAR_DECLARATIONS *DiagnosticGroup

	// Warnings reported by the linter. If you enable these as errors in your build targets,
	// the JS Compiler team will break your build and not rollback.
	DiagnosticGroups_LINT_CHECKS *DiagnosticGroup

	DiagnosticGroups_STRICT_MODULE_CHECKS *DiagnosticGroup

	// Similar to the lintChecks group above, but includes things that cannot be done on a single
	// file at a time, for example because they require typechecking. If you enable these as errors
	// in your build targets, the JS Compiler team will break your build and not rollback.
	DiagnosticGroups_ANALYZER_CHECKS *DiagnosticGroup

	DiagnosticGroups_CLOSURE_DEP_METHOD_USAGE_CHECKS *DiagnosticGroup

	DiagnosticGroups_CLOSURE_CLASS_CHECKS *DiagnosticGroup

	// This group exists so that tests can check for these warnings. It is intentionally not
	// named so that it is is not suppressible via the command line or in code.
	DiagnosticGroups_MALFORMED_GOOG_MODULE *DiagnosticGroup

	// This diagnostic group is intentionally absent in ParserConfig.properties.
	// Conformance checks are supposed to be enforced project-wide, so we don't
	// allow suppressions on individual functions.
	// In the future, we may carve out a subset of the conformance checks that is
	// OK to suppress.
	// For now, the only way to suppress a check at a granularity smaller than
	// the file level is by using a allowlist file.
	DiagnosticGroups_CONFORMANCE_VIOLATIONS *DiagnosticGroup

	DiagnosticGroups_LATE_PROVIDE *DiagnosticGroup

	DiagnosticGroups_DUPLICATE_NAMESPACES *DiagnosticGroup

	DiagnosticGroups_INVALID_DEFINES *DiagnosticGroup

	DiagnosticGroups_MISSING_POLYFILL *DiagnosticGroup

	DiagnosticGroups_POLYMER *DiagnosticGroup

	DiagnosticGroups_BOUNDED_GENERICS *DiagnosticGroup

	// This diagnostic group is intentionally absent in ParserConfig.properties and unnamed. User code
	// should never suppress parse errors but it is useful occasionally for tooling to check whether a
	// given error is from parsing.
	DiagnosticGroups_PARSING *DiagnosticGroup

	DiagnosticGroups_CLOSURE_UNAWARE_CODE_ANNOTATION_PRESENT *DiagnosticGroup

	DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION *DiagnosticGroup
)

// Go-only: language-adapter — Java runs the static initializers of
// DiagnosticGroups top to bottom, and the order in which they register
// groups is observable through getRegisteredGroups(). Go's package-level
// variable initialization follows dependencies across files instead (a
// group whose DiagnosticTypes live in a later file would be registered
// late), so the static fields and static blocks run here, in Java order,
// from init(), after every DiagnosticType has been initialized. No
// package-level initializer outside this file reads these fields (in Java,
// only J2clSuppressWarningsGuard's static field does, and it is initialized
// from its own init(), which runs after this one).
func init() {
	diagnosticGroups_clinit()
}

// Java: DiagnosticGroups static initializers
func diagnosticGroups_clinit() {
	// TODO(b/123768968) remove this diagnostic group, do not allow this suppression. Instead the only
	// work around should be to raise the output language to a high enough level. We need to suppress
	// right now because we don't have any language output level higher than ES5.
	DiagnosticGroups_UNTRANSPILABLE_FEATURES = DiagnosticGroups_RegisterGroupOfTypes(
		"untranspilableFeatures",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
	)

	DiagnosticGroups_MODULE_LOAD = DiagnosticGroups_RegisterGroupOfTypes(
		"moduleLoad",
		ModuleLoader_LOAD_WARNING,
		ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT,
		ProcessCommonJSModules_UNKNOWN_REQUIRE_ENSURE,
	)

	DiagnosticGroups_MODULE_IMPORT = DiagnosticGroups_RegisterGroupOfTypes(
		"moduleImport",
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT,
		ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS,
	)

	DiagnosticGroups_GLOBAL_THIS = DiagnosticGroups_RegisterGroupOfTypes(
		"globalThis",
		CheckGlobalThis_GLOBAL_THIS,
	)

	DiagnosticGroups_DEPRECATED = DiagnosticGroups_RegisterGroupOfTypes(
		"deprecated",
		CheckAccessControls_DEPRECATED_NAME,
		CheckAccessControls_DEPRECATED_NAME_REASON,
		CheckAccessControls_DEPRECATED_PROP,
		CheckAccessControls_DEPRECATED_PROP_REASON,
		CheckAccessControls_DEPRECATED_CLASS,
		CheckAccessControls_DEPRECATED_CLASS_REASON,
	)

	DiagnosticGroups_UNDERSCORE = DiagnosticGroups_RegisterDeprecatedGroup("underscore")

	DiagnosticGroups_VISIBILITY = DiagnosticGroups_RegisterGroupOfTypes(
		"visibility",
		CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS,
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS,
		CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS,
		CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS,
		CheckAccessControls_EXTEND_FINAL_CLASS,
		CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN,
		CheckAccessControls_PRIVATE_OVERRIDE,
		CheckAccessControls_VISIBILITY_MISMATCH,
	)

	// TODO(tbreisacher): Deprecate this and keep just the "visibility" group.
	DiagnosticGroups_ACCESS_CONTROLS = DiagnosticGroups_RegisterGroup(
		"accessControls",
		DiagnosticGroups_VISIBILITY,
	)

	DiagnosticGroups_NON_STANDARD_JSDOC = DiagnosticGroups_RegisterGroupOfTypes(
		"nonStandardJsDocs",
		RhinoErrorReporter_BAD_JSDOC_ANNOTATION,
		RhinoErrorReporter_INVALID_PARAM,
		RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING,
		CheckJSDoc_JSDOC_IN_BLOCK_COMMENT,
	)

	DiagnosticGroups_INVALID_CASTS = DiagnosticGroups_RegisterGroupOfTypes(
		"invalidCasts",
		TypeValidator_INVALID_CAST,
	)

	DiagnosticGroups_STRICT_MODULE_DEP_CHECK = DiagnosticGroups_RegisterGroupOfTypes(
		"strictModuleDepCheck",
		VarCheck_STRICT_CHUNK_DEP_ERROR,
		CheckClosureImports_CROSS_CHUNK_REQUIRE_ERROR,
	)

	DiagnosticGroups_VIOLATED_MODULE_DEP = DiagnosticGroups_RegisterGroupOfTypes(
		"violatedModuleDep",
		VarCheck_VIOLATED_CHUNK_DEP_ERROR,
	)

	DiagnosticGroups_EXTERNS_VALIDATION = DiagnosticGroups_RegisterGroupOfTypes(
		"externsValidation",
		VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR,
		VarCheck_UNDEFINED_EXTERN_VAR_ERROR,
	)

	DiagnosticGroups_UNKNOWN_DEFINES = DiagnosticGroups_RegisterGroupOfTypes(
		"unknownDefines",
		ProcessDefines_UNKNOWN_DEFINE_WARNING,
	)

	DiagnosticGroups_DEFINE_WITHOUT_GOOG_DEFINE = DiagnosticGroups_RegisterGroupOfTypes(
		"defineWithoutGoogDefine",
		ProcessDefines_DEFINE_WITHOUT_GOOG_DEFINE,
	)

	DiagnosticGroups_TWEAKS = DiagnosticGroups_RegisterGroupOfTypes(
		"tweakValidation",
		ProcessTweaks_INVALID_TWEAK_DEFAULT_VALUE_WARNING,
		ProcessTweaks_TWEAK_WRONG_GETTER_TYPE_WARNING,
	)

	DiagnosticGroups_MISSING_CALL_TO_SUPER = DiagnosticGroups_RegisterGroupOfTypes(
		"missingCallToSuper", CheckSuper_MISSING_ROOT_LEVEL_CALL_TO_SUPER,
	)

	DiagnosticGroups_MISSING_OVERRIDE = DiagnosticGroups_RegisterGroupOfTypes(
		"missingOverride",
		TypeCheck_HIDDEN_INTERFACE_PROPERTY,
		TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY,
		TypeCheck_HIDDEN_SUPERCLASS_PROPERTY,
	)

	DiagnosticGroups_MISSING_PROPERTIES = DiagnosticGroups_RegisterGroupOfTypes(
		"missingProperties",
		TypeCheck_INEXISTENT_PROPERTY,
		TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION,
		TypeCheck_POSSIBLE_INEXISTENT_PROPERTY,
	)

	DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES = DiagnosticGroups_RegisterGroupOfTypes(
		"globallyMissingProperties",
		TypeCheck_POSSIBLE_INEXISTENT_PROPERTY,
	)

	DiagnosticGroups_J2CL_CHECKS = DiagnosticGroups_RegisterGroupOfTypes(
		"j2clChecks",
		J2clChecksPass_J2CL_REFERENCE_EQUALITY,
	)

	DiagnosticGroups_MISSING_RETURN = DiagnosticGroups_RegisterGroupOfTypes(
		"missingReturn",
		CheckMissingReturn_MISSING_RETURN_STATEMENT,
	)

	DiagnosticGroups_UNDEFINED_VARIABLES = DiagnosticGroups_RegisterGroupOfTypes(
		"undefinedVars",
		VarCheck_UNDEFINED_VAR_ERROR,
	)

	DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT = DiagnosticGroups_RegisterGroupOfTypes(
		"checkDebuggerStatement",
		CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT,
	)

	DiagnosticGroups_CHECK_REGEXP = DiagnosticGroups_RegisterGroupOfTypes(
		"checkRegExp",
		CheckRegExp_REGEXP_REFERENCE,
		CheckRegExp_MALFORMED_REGEXP,
	)

	// NOTE(dimvar): it'd be nice to add TypedScopeCreator.ALL_DIAGNOSTICS here,
	// but we would first need to cleanup projects that would break because
	// they set --jscomp_error=checkTypes.
	DiagnosticGroups_CHECK_TYPES = DiagnosticGroups_RegisterGroupOfGroups(
		"checkTypes",
		TypeValidator_ALL_DIAGNOSTICS,
		TypeCheck_ALL_DIAGNOSTICS,
		FunctionTypeBuilder_ALL_DIAGNOSTICS,
		DiagnosticGroups_GLOBAL_THIS,
		NewDiagnosticGroup(
			CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL,
			CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION,
		),
	)

	DiagnosticGroups_ES5_INHERITANCE_DIAGNOSTIC_GROUP = DiagnosticGroups_RegisterGroup(
		"checkEs5InheritanceCorrectnessConditions",
		TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP,
	)

	DiagnosticGroups_CHECK_PROTOTYPAL_TYPES = DiagnosticGroups_RegisterGroupOfTypes(
		"checkPrototypalTypes",
		TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE,
		TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY,
		TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH,
	)

	// This group exists for the J2CL team to suppress the associated diagnostics using Java code
	// rather than `@suppress` annotations.
	DiagnosticGroups_CHECK_STATIC_OVERRIDES = DiagnosticGroups_CHECK_PROTOTYPAL_TYPES

	DiagnosticGroups_TOO_MANY_TYPE_PARAMS = DiagnosticGroups_RegisterGroupOfTypes(
		"tooManyTypeParams",
		RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS,
	)

	DiagnosticGroups_STRICT_MISSING_PROPERTIES = DiagnosticGroups_RegisterGroupOfTypes(
		"strictMissingProperties",
		TypeCheck_STRICT_INEXISTENT_PROPERTY,
		TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION,
		TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY,
		TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE,
	)

	DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS = DiagnosticGroups_RegisterGroupOfTypes(
		"strictPrimitiveOperators",
		TypeValidator_INVALID_OPERAND_TYPE,
	)

	DiagnosticGroups_STRICT_CHECK_TYPES = DiagnosticGroups_RegisterGroupOfGroups(
		"strictCheckTypes",
		DiagnosticGroups_STRICT_MISSING_PROPERTIES,
		DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS,
	)

	DiagnosticGroups_REPORT_UNKNOWN_TYPES = DiagnosticGroups_RegisterGroupOfTypes(
		"reportUnknownTypes",
		TypeCheck_UNKNOWN_EXPR_TYPE,
	)

	DiagnosticGroups_CHECK_VARIABLES = DiagnosticGroups_RegisterGroupOfTypes(
		"checkVars",
		VarCheck_UNDEFINED_VAR_ERROR,
		VarCheck_VAR_MULTIPLY_DECLARED_ERROR,
		VariableReferenceCheck_EARLY_REFERENCE,
		VariableReferenceCheck_REDECLARED_VARIABLE,
	)

	DiagnosticGroups_CHECK_USELESS_CODE = DiagnosticGroups_RegisterGroupOfTypes(
		"uselessCode",
		CheckSideEffects_USELESS_CODE_ERROR,
		CheckUnreachableCode_UNREACHABLE_CODE,
	)

	DiagnosticGroups_CONST = DiagnosticGroups_RegisterGroupOfTypes(
		"const",
		CheckAccessControls_CONST_PROPERTY_DELETED,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE,
		ConstCheck_CONST_REASSIGNED_VALUE_ERROR,
	)

	DiagnosticGroups_CONSTANT_PROPERTY = DiagnosticGroups_RegisterGroupOfTypes(
		"constantProperty",
		CheckAccessControls_CONST_PROPERTY_DELETED,
		CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE,
	)

	DiagnosticGroups_ACCESS_CONTROLS_CONST = DiagnosticGroups_RegisterGroup(
		"accessControlsConst",
		DiagnosticGroups_CONSTANT_PROPERTY,
	)

	DiagnosticGroups_TYPE_INVALIDATION = DiagnosticGroups_RegisterGroupOfTypes(
		"typeInvalidation",
		DisambiguateProperties_PROPERTY_INVALIDATION,
	)

	DiagnosticGroups_DUPLICATE_VARS = DiagnosticGroups_RegisterGroupOfTypes(
		"duplicate",
		InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING,
		TypeValidator_DUP_VAR_DECLARATION,
		TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH,
		TypeCheck_FUNCTION_MASKS_VARIABLE,
		VarCheck_VAR_MULTIPLY_DECLARED_ERROR,
		VariableReferenceCheck_REDECLARED_VARIABLE,
	)

	DiagnosticGroups_ES5_STRICT = DiagnosticGroups_RegisterGroupOfTypes(
		"es5Strict",
		RhinoErrorReporter_INVALID_OCTAL_LITERAL,
		RhinoErrorReporter_DUPLICATE_PARAM,
		StrictModeCheck_ARGUMENTS_ASSIGNMENT,
		StrictModeCheck_ARGUMENTS_DECLARATION,
		StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN,
		StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN,
		StrictModeCheck_DELETE_VARIABLE,
		StrictModeCheck_DUPLICATE_MEMBER,
		StrictModeCheck_EVAL_ASSIGNMENT,
		StrictModeCheck_EVAL_DECLARATION,
		StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN,
		StrictModeCheck_FUNCTION_CALLER_FORBIDDEN,
		StrictModeCheck_USE_OF_WITH,
	)

	DiagnosticGroups_WEAK_MODULE_GET = DiagnosticGroups_RegisterGroupOfTypes(
		"weakModuleGet",
		ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE,
	)

	DiagnosticGroups_MISSING_PROVIDE = DiagnosticGroups_RegisterGroupOfTypes(
		"missingProvide",
		ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE,
	)

	DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR = DiagnosticGroups_RegisterGroupOfTypes(
		"unrecognizedTypeError",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR,
	)

	// identical to UNRECOGNIZED_TYPE_ERROR, but allowed in @suppress tags
	DiagnosticGroups_DANGEROUS_UNRECOGNIZED_TYPE_ERROR = DiagnosticGroups_RegisterGroup(
		"dangerousUnrecognizedTypeError",
		DiagnosticGroups_UNRECOGNIZED_TYPE_ERROR,
	)

	DiagnosticGroups_MISSING_REQUIRE = DiagnosticGroups_RegisterGroupOfTypes(
		"missingRequire",
		CheckMissingRequires_MISSING_REQUIRE,
		CheckMissingRequires_MISSING_REQUIRE_IN_PROVIDES_FILE,
		CheckMissingRequires_MISSING_REQUIRE_IN_GOOG_SCOPE,
		CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE,
		CheckMissingRequires_MISSING_REQUIRE_TYPE,
		CheckMissingRequires_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE,
		CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE,
		CheckMissingRequires_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE,
		CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE,
		CheckMissingRequires_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE,
	)

	// A set of diagnostics expected when parsing and type checking partial programs. Useful for clutz
	// (tool that extracts TypeScript definitions from JS code).
	DiagnosticGroups_MISSING_SOURCES_WARNINGS = DiagnosticGroups_RegisterGroupOfGroups(
		"missingSourcesWarnings",
		DiagnosticGroups_REPORT_UNKNOWN_TYPES,
		DiagnosticGroups_UNDEFINED_VARIABLES,
		DiagnosticGroups_MISSING_PROVIDE,
		DiagnosticGroup_ForType(ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE),
		DiagnosticGroups_MISSING_PROPERTIES,
		DiagnosticGroups_DUPLICATE_VARS,
		DiagnosticGroup_ForType(ProcessDefines_INVALID_DEFINE_VALUE),
	)

	DiagnosticGroups_EXTRA_REQUIRE = DiagnosticGroups_RegisterGroupOfTypes(
		"extraRequire",
		CheckExtraRequires_EXTRA_REQUIRE_WARNING,
	)

	DiagnosticGroups_DUPLICATE_MESSAGE = DiagnosticGroups_RegisterGroupOfTypes(
		"duplicateMessage",
		JsMessageVisitor_MESSAGE_DUPLICATE_KEY,
	)

	DiagnosticGroups_MESSAGE_DESCRIPTIONS = DiagnosticGroups_RegisterGroupOfTypes(
		"msgDescriptions",
		JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION,
	)

	// Warnings that only apply to people who use MSG_ to denote messages. Note that this doesn't
	// include warnings about proper use of goog.getMsg
	DiagnosticGroups_MSG_CONVENTIONS = DiagnosticGroups_RegisterGroupOfTypes(
		"messageConventions",
		JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION,
		JsMessageVisitor_MESSAGE_HAS_NO_TEXT,
		JsMessageVisitor_MESSAGE_TREE_MALFORMED,
		JsMessageVisitor_MESSAGE_HAS_NO_VALUE,
		JsMessageVisitor_MESSAGE_DUPLICATE_KEY,
		JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY,
	)

	DiagnosticGroups_MISPLACED_TYPE_ANNOTATION = DiagnosticGroups_RegisterGroupOfTypes(
		"misplacedTypeAnnotation",
		CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR,
		CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION,
		CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL,
		CheckJSDoc_DISALLOWED_MEMBER_JSDOC,
		CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION,
		CheckJSDoc_INVALID_MODIFIES_ANNOTATION,
		CheckJSDoc_JSDOC_ON_RETURN,
		CheckJSDoc_MISPLACED_ANNOTATION,
		CheckJSDoc_MISPLACED_MSG_ANNOTATION,
	)

	DiagnosticGroups_MISPLACED_MSG_ANNOTATION = DiagnosticGroups_RegisterGroupOfTypes(
		"misplacedMsgAnnotation",
		CheckJSDoc_MISPLACED_MSG_ANNOTATION,
	)

	DiagnosticGroups_MISPLACED_SUPPRESS = DiagnosticGroups_RegisterGroupOfTypes(
		"misplacedSuppress",
		CheckJSDoc_MISPLACED_SUPPRESS,
	)

	DiagnosticGroups_SUSPICIOUS_CODE = DiagnosticGroups_RegisterGroupOfTypes(
		"suspiciousCode",
		CheckDuplicateCase_DUPLICATE_CASE,
		CheckSuspiciousCode_SUSPICIOUS_SEMICOLON,
		CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN,
		CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN,
		CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR,
		CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND,
		CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR,
		CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR,
		ProcessCommonJSModules_SUSPICIOUS_EXPORTS_ASSIGNMENT,
		TypeCheck_DETERMINISTIC_TEST,
	)

	DiagnosticGroups_FUNCTION_PARAMS = DiagnosticGroups_RegisterGroupOfTypes(
		"functionParams",
		FunctionTypeBuilder_INEXISTENT_PARAM,
		FunctionTypeBuilder_OPTIONAL_ARG_AT_END,
	)

	DiagnosticGroups_DEPRECATED_ANNOTATIONS = DiagnosticGroups_RegisterGroupOfTypes(
		"deprecatedAnnotations",
		CheckJSDoc_ANNOTATION_DEPRECATED,
	)

	DiagnosticGroups_UNUSED_LOCAL_VARIABLE = DiagnosticGroups_RegisterGroupOfTypes(
		"unusedLocalVariables",
		VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT,
	)

	DiagnosticGroups_JSDOC_MISSING_TYPE = DiagnosticGroups_RegisterGroupOfTypes(
		"jsdocMissingType",
		RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING,
	)

	DiagnosticGroups_TYPE_IMPORT_CODE_REFERENCES = DiagnosticGroups_RegisterGroupOfTypes(
		"typeImportCodeReferences",
		CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE,
	)

	DiagnosticGroups_PARTIAL_ALIAS = DiagnosticGroups_RegisterGroupOfTypes(
		"partialAlias",
		InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING,
	)

	// This lint is given its own diagnostic group because it's harder to fix than other
	// lint errors, and we want to discourage users from doing a blanket @suppress {lintChecks}.
	// This is intentionally not public. It should not be enabled directly; instead enable lintChecks.
	DiagnosticGroups_USE_OF_GOOG_PROVIDE = DiagnosticGroups_RegisterGroupOfTypes(
		"useOfGoogProvide",
		ClosureCheckModule_USE_OF_GOOG_PROVIDE,
	)

	// This is intended to be used to suppress warnings in code that cannot for some reason be updated
	// to use `let` and `const` instead of `var`. It should not be enabled directly, instead enable
	// "lintChecks".
	DiagnosticGroups_LINT_VAR_DECLARATIONS = DiagnosticGroups_RegisterGroupOfTypes(
		"lintVarDeclarations",
		CheckVar_VAR,
	)

	// Warnings reported by the linter. If you enable these as errors in your build targets,
	// the JS Compiler team will break your build and not rollback.
	DiagnosticGroups_LINT_CHECKS = DiagnosticGroups_RegisterGroupOfGroups(
		"lintChecks",
		CheckJSDocStyle_LINT_DIAGNOSTICS,
		DiagnosticGroups_USE_OF_GOOG_PROVIDE,
		DiagnosticGroups_LINT_VAR_DECLARATIONS,
		NewDiagnosticGroup(
			CheckClosureImports_LET_CLOSURE_IMPORT,
			CheckConstPrivateProperties_MISSING_CONST_PROPERTY,
			CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME,
			CheckConstantCaseNames_MISSING_CONST_PROPERTY,
			CheckEmptyStatements_USELESS_EMPTY_STATEMENT,
			CheckEnums_COMPUTED_PROP_NAME_IN_ENUM,
			CheckEnums_DUPLICATE_ENUM_VALUE,
			CheckEnums_ENUM_PROP_NOT_CONSTANT,
			CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER,
			CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM,
			CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM,
			CheckEs6ModuleFileStructure_MUST_COME_BEFORE,
			CheckEs6Modules_DUPLICATE_IMPORT,
			CheckEs6Modules_NO_DEFAULT_EXPORT,
			CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE,
			CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY,
			CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS,
			CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS,
			CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE,
			CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT,
			CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS,
			CheckMissingSemicolon_MISSING_SEMICOLON,
			CheckNoMutatedEs6Exports_MUTATED_EXPORT,
			CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC,
			CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC,
			CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC,
			CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT,
			CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION,
			CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER,
			CheckProvidesSorted_PROVIDES_NOT_SORTED,
			CheckRequiresSorted_REQUIRES_NOT_SORTED,
			CheckUnusedLabels_UNUSED_LABEL,
			CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY,
			CheckUselessBlocks_USELESS_BLOCK,
			ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE,
			ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE,
			ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION,
			ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE,
			RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING,
			RhinoErrorReporter_UNNECESSARY_ESCAPE,
			RhinoErrorReporter_STRING_CONTINUATION,
		),
	)

	DiagnosticGroups_STRICT_MODULE_CHECKS = DiagnosticGroups_RegisterGroupOfTypes(
		"strictModuleChecks",
		ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE,
		ClosureCheckModule_LET_GOOG_REQUIRE,
		ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME,
		ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME,
	)

	// Similar to the lintChecks group above, but includes things that cannot be done on a single
	// file at a time, for example because they require typechecking. If you enable these as errors
	// in your build targets, the JS Compiler team will break your build and not rollback.
	DiagnosticGroups_ANALYZER_CHECKS = DiagnosticGroups_RegisterGroupOfTypes(
		"analyzerChecks",
		CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT,
		ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC,
		ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC,
		CheckNestedNames_NESTED_NAME_IN_GOOG_MODULE,
	)

	DiagnosticGroups_CLOSURE_DEP_METHOD_USAGE_CHECKS = DiagnosticGroups_RegisterGroupOfTypes(
		"closureDepMethodUsageChecks",
		ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_ERROR,
		ProcessClosurePrimitives_CLOSURE_CALL_CANNOT_BE_ALIASED_OUTSIDE_MODULE_ERROR,
		ClosurePrimitiveErrors_INVALID_CLOSURE_CALL_SCOPE_ERROR,
		ClosurePrimitiveErrors_INVALID_GET_CALL_SCOPE,
	)

	DiagnosticGroups_CLOSURE_CLASS_CHECKS = DiagnosticGroups_RegisterGroupOfTypes(
		"closureClassChecks",
		ProcessClosurePrimitives_POSSIBLE_BASE_CLASS_ERROR,
	)

	// This group exists so that tests can check for these warnings. It is intentionally not
	// named so that it is is not suppressible via the command line or in code.
	DiagnosticGroups_MALFORMED_GOOG_MODULE = DiagnosticGroups_RegisterUnsuppressibleGroup(
		ClosureCheckModule_GOOG_MODULE_MISPLACED,
		ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE,
	)

	// This group exists so that generated code can suppress these
	// warnings. Not for general use. These diagnostics will most likely
	// be moved to the suspiciousCode group.
	// Java: static {
	DiagnosticGroups_RegisterGroupOfTypes(
		"transitionalSuspiciousCodeWarnings",
		PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND,
	)
	// Java: } (end of static block)

	// This diagnostic group is intentionally absent in ParserConfig.properties.
	// Conformance checks are supposed to be enforced project-wide, so we don't
	// allow suppressions on individual functions.
	// In the future, we may carve out a subset of the conformance checks that is
	// OK to suppress.
	// For now, the only way to suppress a check at a granularity smaller than
	// the file level is by using a allowlist file.
	DiagnosticGroups_CONFORMANCE_VIOLATIONS = DiagnosticGroups_RegisterGroupOfTypes(
		"conformanceViolations",
		CheckConformance_CONFORMANCE_VIOLATION,
		CheckConformance_CONFORMANCE_POSSIBLE_VIOLATION,
	)

	DiagnosticGroups_LATE_PROVIDE = DiagnosticGroups_RegisterGroupOfTypes(
		"lateProvide",
		CheckClosureImports_LATE_PROVIDE_ERROR,
	)

	DiagnosticGroups_DUPLICATE_NAMESPACES = DiagnosticGroups_RegisterUnsuppressibleGroup(
		ClosurePrimitiveErrors_DUPLICATE_MODULE,
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE,
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE,
	)

	DiagnosticGroups_INVALID_DEFINES = DiagnosticGroups_RegisterUnsuppressibleGroup(
		ProcessDefines_INVALID_DEFINE_VALUE,
		ProcessDefines_INVALID_DEFINE_TYPE,
	)

	DiagnosticGroups_MISSING_POLYFILL = DiagnosticGroups_RegisterGroupOfTypes(
		"missingPolyfill",
		RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR,
	)

	DiagnosticGroups_POLYMER = DiagnosticGroups_RegisterGroupOfTypes(
		"polymer",
		PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID,
	)

	DiagnosticGroups_BOUNDED_GENERICS = DiagnosticGroups_RegisterGroupOfTypes(
		"boundedGenerics",
		RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES,
		RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR,
	)

	// This diagnostic group is intentionally absent in ParserConfig.properties and unnamed. User code
	// should never suppress parse errors but it is useful occasionally for tooling to check whether a
	// given error is from parsing.
	DiagnosticGroups_PARSING = DiagnosticGroup_ForType(RhinoErrorReporter_PARSE_ERROR)

	DiagnosticGroups_CLOSURE_UNAWARE_CODE_ANNOTATION_PRESENT = DiagnosticGroups_RegisterGroupOfTypes(
		"closureUnawareCodeAnnotationPresent",
		RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT,
	)

	// For internal use only, so there are no constants for these groups.
	// Java: static {
	DiagnosticGroups_RegisterGroupOfTypes(
		"polymerBehavior",
		PolymerPassErrors_POLYMER_UNQUALIFIED_BEHAVIOR,
	)
	DiagnosticGroups_RegisterGroupOfTypes(
		"invalidProvide",
		ProcessClosurePrimitives_INVALID_PROVIDE_ERROR,
	)
	DiagnosticGroups_RegisterGroupOfTypes(
		"conflictingIjsFile",
		IjsErrors_CONFLICTING_IJS_FILE,
	)
	DiagnosticGroups_RegisterGroupOfTypes(
		"implicitWeakEntryPoint",
		JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR,
	)
	// Java: } (end of static block)

	DiagnosticGroups_ARTIFICIAL_FUNCTION_PURITY_VALIDATION = DiagnosticGroups_RegisterGroupOfTypes(
		"artificialFunctionPurityValidation",
		PureFunctionIdentifier_UNUSED_ARTIFICIAL_PURE_ANNOTATION,
	)
}

// SetWarningLevel adds warning levels by name.
//
// Java: DiagnosticGroups.setWarningLevel(CompilerOptions, String, CheckLevel)
func (d *DiagnosticGroups) SetWarningLevel(options *CompilerOptions, name string, level CheckLevel) {
	group := DiagnosticGroups_ForName(name)
	jutil.CheckNotNullMsg1(group, "No warning class for name: %s", name)
	options.SetWarningLevel(group, level)
}
