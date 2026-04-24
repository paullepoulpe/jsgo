// Port of com/google/javascript/jscomp/DiagnosticGroups.java
// @ closure-compiler v20260415

package diag

import "sync"

// DiagnosticGroups is the registry of named diagnostic groups exposed by the
// compiler. Mirrors Java's DiagnosticGroups class. Each public field
// corresponds to a Java public-static-final field.
//
// Groups are registered at init time and can be looked up by name via ForName.
type DiagnosticGroups struct{}

var (
	groupsByNameMu sync.RWMutex
	groupsByName   = make(map[string]*DiagnosticGroup)
)

// RegisterGroup registers a named group built from the given DiagnosticTypes.
// Mirrors Java's DiagnosticGroups.registerGroup(String, DiagnosticType...).
func RegisterGroup(name string, types ...*DiagnosticType) *DiagnosticGroup {
	g := NewDiagnosticGroup(name, types...)
	groupsByNameMu.Lock()
	groupsByName[name] = g
	groupsByNameMu.Unlock()
	return g
}

// RegisterGroupComposite registers a named composite group from sub-groups.
// Mirrors Java's DiagnosticGroups.registerGroup(String, DiagnosticGroup...).
func RegisterGroupComposite(name string, groups ...*DiagnosticGroup) *DiagnosticGroup {
	g := NewDiagnosticGroupComposite(name, groups...)
	groupsByNameMu.Lock()
	groupsByName[name] = g
	groupsByNameMu.Unlock()
	return g
}

// RegisterDeprecatedGroup creates and registers a group that maps a deprecated
// name to a single UNUSED warning. Mirrors Java's registerDeprecatedGroup.
func RegisterDeprecatedGroup(name string) *DiagnosticGroup {
	return RegisterGroup(name, Unused)
}

// RegisterUnsuppressibleGroup creates an unnamed (unsuppressible) group.
// Mirrors Java's registerUnsuppressibleGroup.
func RegisterUnsuppressibleGroup(types ...*DiagnosticType) *DiagnosticGroup {
	return NewDiagnosticGroupAnonymous(types...)
}

// ForName looks up a registered group by name. Returns nil if not found.
// Mirrors Java's DiagnosticGroups.forName.
func ForName(name string) *DiagnosticGroup {
	groupsByNameMu.RLock()
	defer groupsByNameMu.RUnlock()
	return groupsByName[name]
}

// GetRegisteredGroups returns a snapshot of all registered groups indexed by name.
func GetRegisteredGroups() map[string]*DiagnosticGroup {
	groupsByNameMu.RLock()
	defer groupsByNameMu.RUnlock()
	out := make(map[string]*DiagnosticGroup, len(groupsByName))
	for k, v := range groupsByName {
		out[k] = v
	}
	return out
}

// Unused is a sentinel DiagnosticType used by deprecated groups.
// Mirrors Java's DiagnosticGroups.UNUSED.
var Unused = DiagnosticTypeWarning("JSC_UNUSED", "{0}")

// ---------------------------------------------------------------------------
// Named groups — public constants mirroring Java public-static-final fields.
// Only groups relevant to the passes jsgo currently runs are wired to real
// DiagnosticType sets; the rest use the UNUSED sentinel so that
// --jscomp_warning=<name> parsing won't panic on unknown names.
// ---------------------------------------------------------------------------

// Well-known DiagnosticType stubs used solely to populate group memberships
// below. These are declared here rather than in per-pass files because
// DiagnosticGroups depends on no pass, matching Java's package structure.
//
// A full DiagnosticType registry (with format strings) lives in each pass
// file (typecheck.go, varcheck.go, etc.). The stubs below carry only keys.
var (
	// --- GlobalThis ---
	dtGlobalThis = DiagnosticTypeWarning("JSC_USED_GLOBAL_THIS", "")

	// --- CheckAccessControls ---
	dtBadPrivateGlobalAccess          = DiagnosticTypeWarning("JSC_BAD_PRIVATE_GLOBAL_ACCESS", "")
	dtBadPrivatePropertyAccess        = DiagnosticTypeWarning("JSC_BAD_PRIVATE_PROPERTY_ACCESS", "")
	dtBadPackagePropertyAccess        = DiagnosticTypeWarning("JSC_BAD_PACKAGE_PROPERTY_ACCESS", "")
	dtBadProtectedPropertyAccess      = DiagnosticTypeWarning("JSC_BAD_PROTECTED_PROPERTY_ACCESS", "")
	dtExtendFinalClass                = DiagnosticTypeWarning("JSC_EXTEND_FINAL_CLASS", "")
	dtFinalPropertyOverridden         = DiagnosticTypeWarning("JSC_FINAL_PROPERTY_OVERRIDDEN", "")
	dtPrivateOverride                 = DiagnosticTypeWarning("JSC_PRIVATE_OVERRIDE", "")
	dtVisibilityMismatch              = DiagnosticTypeWarning("JSC_VISIBILITY_MISMATCH", "")
	dtDeprecatedName                  = DiagnosticTypeWarning("JSC_DEPRECATED_VAR", "")
	dtDeprecatedNameReason            = DiagnosticTypeWarning("JSC_DEPRECATED_VAR_REASON", "")
	dtDeprecatedProp                  = DiagnosticTypeWarning("JSC_DEPRECATED_PROP", "")
	dtDeprecatedPropReason            = DiagnosticTypeWarning("JSC_DEPRECATED_PROP_REASON", "")
	dtDeprecatedClass                 = DiagnosticTypeWarning("JSC_DEPRECATED_CLASS", "")
	dtDeprecatedClassReason           = DiagnosticTypeWarning("JSC_DEPRECATED_CLASS_REASON", "")
	dtConstPropertyReassigned         = DiagnosticTypeWarning("JSC_CONSTANT_PROPERTY_REASSIGNED_VALUE", "")
	dtConstPropertyDeleted            = DiagnosticTypeWarning("JSC_CONSTANT_PROPERTY_DELETED", "")

	// --- TypeValidator ---
	dtInvalidCast                 = DiagnosticTypeWarning("JSC_INVALID_CAST", "")
	dtDupVarDeclaration           = DiagnosticTypeWarning("JSC_DUP_VAR_DECLARATION", "")
	dtDupVarDeclarationTypeMismatch = DiagnosticTypeWarning("JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH", "")
	dtInvalidOperandType          = DiagnosticTypeWarning("JSC_INVALID_OPERAND_TYPE", "")

	// --- TypeCheck ---
	dtTypeMismatch                    = DiagnosticTypeError("JSC_TYPE_MISMATCH", "")
	dtInexistentProperty              = DiagnosticTypeWarning("JSC_INEXISTENT_PROPERTY", "")
	dtPossibleInexistentProperty      = DiagnosticTypeWarning("JSC_POSSIBLE_INEXISTENT_PROPERTY", "")
	dtNotAConstructor                 = DiagnosticTypeError("JSC_NOT_A_CONSTRUCTOR", "")
	dtIllegalPropertyAccess           = DiagnosticTypeError("JSC_ILLEGAL_PROPERTY_ACCESS", "")
	dtHiddenSuperclassPropertyMismatch = DiagnosticTypeWarning("JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH", "")
	dtHiddenInterfacePropertyMismatch  = DiagnosticTypeWarning("JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH", "")
	dtFunctionMasksVariable            = DiagnosticTypeWarning("JSC_FUNCTION_MASKS_VARIABLE", "")
	dtDeterministicTest                = DiagnosticTypeWarning("JSC_DETERMINISTIC_TEST", "")
	dtUnknownExprType                  = DiagnosticTypeWarning("JSC_UNKNOWN_TYPEOF_VALUE", "")
	dtStrictInexistentProperty         = DiagnosticTypeWarning("JSC_POSSIBLE_INEXISTENT_PROPERTY", "")
	dtIllegalPropertyCreationOnUnion   = DiagnosticTypeWarning("JSC_ILLEGAL_PROPERTY_CREATION", "")

	// --- VarCheck ---
	dtUndefinedVarError        = DiagnosticTypeError("JSC_UNDEFINED_VARIABLE", "")
	dtVarMultiplyDeclaredError = DiagnosticTypeError("JSC_VAR_MULTIPLY_DECLARED_ERROR", "")

	// --- VariableReferenceCheck ---
	dtEarlyReference    = DiagnosticTypeWarning("JSC_REFERENCE_BEFORE_DECLARE", "")
	dtRedeclaredVariable = DiagnosticTypeWarning("JSC_DUP_VAR_DECLARATION", "")

	// --- CheckMissingReturn ---
	dtMissingReturnStatement = DiagnosticTypeWarning("JSC_MISSING_RETURN_STATEMENT", "")

	// --- CheckDebuggerStatement ---
	dtDebuggerStatementPresent = DiagnosticTypeDisabled("JSC_DEBUGGER_STATEMENT_PRESENT", "")

	// --- CheckRegExp ---
	dtRegexpReference = DiagnosticTypeWarning("JSC_REGEXP_REFERENCE", "")
	dtMalformedRegexp = DiagnosticTypeError("JSC_MALFORMED_REGEXP", "")

	// --- CheckSideEffects / CheckUnreachableCode ---
	dtUselessCode   = DiagnosticTypeWarning("JSC_USELESS_CODE", "")
	dtUnreachableCode = DiagnosticTypeWarning("JSC_UNREACHABLE_CODE", "")

	// --- CheckSuspiciousCode ---
	dtSuspiciousSemicolon                         = DiagnosticTypeWarning("JSC_SUSPICIOUS_SEMICOLON", "")
	dtSuspiciousBreakingOutOfOptionalChain        = DiagnosticTypeWarning("SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN", "")
	dtSuspiciousNan                               = DiagnosticTypeWarning("JSC_SUSPICIOUS_NAN", "")
	dtSuspiciousIn                                = DiagnosticTypeWarning("JSC_SUSPICIOUS_IN", "")
	dtSuspiciousInstanceofLeft                    = DiagnosticTypeWarning("JSC_SUSPICIOUS_INSTANCEOF_LEFT", "")
	dtSuspiciousLeftOperandOfLogicalOperator      = DiagnosticTypeWarning("JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR", "")
	dtSuspiciousNegatedLeftOperandOfInOperator    = DiagnosticTypeWarning("JSC_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR", "")

	// --- CheckJSDoc ---
	dtJsdocInBlockComment            = DiagnosticTypeWarning("JSC_JSDOC_IN_BLOCK_COMMENT", "")
	dtArrowFunctionAsConstructor     = DiagnosticTypeWarning("JSC_ARROW_FUNCTION_AS_CONSTRUCTOR", "")
	dtBadRestParameterAnnotation     = DiagnosticTypeWarning("JSC_BAD_REST_PARAMETER_ANNOTATION", "")
	dtDefaultParamMustBeMarkedOptional = DiagnosticTypeWarning("JSC_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL", "")
	dtInvalidNoSideEffectAnnotation  = DiagnosticTypeWarning("JSC_INVALID_NO_SIDE_EFFECT_ANNOTATION", "")
	dtInvalidModifiesAnnotation      = DiagnosticTypeWarning("JSC_INVALID_MODIFIES_ANNOTATION", "")
	dtJsdocOnReturn                  = DiagnosticTypeWarning("JSC_JSDOC_ON_RETURN", "")
	dtMisplacedAnnotation            = DiagnosticTypeWarning("JSC_MISPLACED_ANNOTATION", "")
	dtMisplacedMsgAnnotation         = DiagnosticTypeWarning("JSC_MISPLACED_MSG_ANNOTATION", "")
	dtAnnotationDeprecated           = DiagnosticTypeWarning("JSC_ANNOTATION_DEPRECATED", "")
	dtDisallowedMemberJsdoc          = DiagnosticTypeWarning("JSC_DISALLOWED_MEMBER_JSDOC", "")
	dtMisplacedSuppress              = DiagnosticTypeWarning("JSC_MISPLACED_SUPPRESS", "")

	// --- CheckMissingRequires ---
	dtMissingRequire                    = DiagnosticTypeWarning("JSC_MISSING_REQUIRE", "")
	dtMissingRequireInProvidesFile      = DiagnosticTypeWarning("JSC_MISSING_REQUIRE_IN_PROVIDES_FILE", "")
	dtMissingRequireInGoogScope         = DiagnosticTypeWarning("JSC_MISSING_REQUIRE_IN_GOOG_SCOPE", "")
	dtMissingRequireTypeInGoogScope     = DiagnosticTypeWarning("JSC_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE", "")
	dtMissingRequireType                = DiagnosticTypeWarning("JSC_MISSING_REQUIRE_TYPE", "")
	dtMissingRequireTypeInProvidesFile  = DiagnosticTypeWarning("JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE", "")
	dtIncorrectNamespaceAliasRequire    = DiagnosticTypeWarning("JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE", "")
	dtIncorrectNamespaceAliasRequireType = DiagnosticTypeWarning("JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE", "")
	dtIndirectNamespaceRefRequire       = DiagnosticTypeWarning("JSC_INDIRECT_NAMESPACE_REF_REQUIRE", "")
	dtIndirectNamespaceRefRequireType   = DiagnosticTypeWarning("JSC_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE", "")

	// --- CheckClosureImports ---
	dtLateProvidesError      = DiagnosticTypeError("JSC_LATE_PROVIDE_ERROR", "")
	dtCrossChunkRequireError = DiagnosticTypeError("JSC_XMODULE_REQUIRE_ERROR", "")

	// --- CheckTypeImportCodeReferences ---
	dtTypeImportCodeReference = DiagnosticTypeWarning("JSC_TYPE_IMPORT_CODE_REFERENCE", "")

	// --- CheckConformance ---
	dtConformanceViolation         = DiagnosticTypeWarning("JSC_CONFORMANCE_VIOLATION", "")
	dtConformancePossibleViolation = DiagnosticTypeWarning("JSC_CONFORMANCE_POSSIBLE_VIOLATION", "")

	// --- RhinoErrorReporter ---
	dtBadJsdocAnnotation         = DiagnosticTypeWarning("JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES", "")
	dtUnrecognizedTypeError       = DiagnosticTypeWarning("JSC_UNRECOGNIZED_TYPE_ERROR", "")
	dtMissingModuleOrProvide      = DiagnosticTypeWarning("JSC_MISSING_MODULE_OR_PROVIDE", "")
	dtMissingModuleOrProvideForFD = DiagnosticTypeWarning("JSC_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE", "")
)

// Named groups — public variables matching Java static-final DiagnosticGroup fields.
// Assigned in init() below.
var (
	GroupGlobalThis       *DiagnosticGroup
	GroupDeprecated       *DiagnosticGroup
	GroupUnderscore       *DiagnosticGroup
	GroupVisibility       *DiagnosticGroup
	GroupAccessControls   *DiagnosticGroup
	GroupNonStandardJsDoc *DiagnosticGroup
	GroupInvalidCasts     *DiagnosticGroup
	GroupStrictModuleDepCheck *DiagnosticGroup
	GroupMissingOverride  *DiagnosticGroup
	GroupMissingProperties *DiagnosticGroup
	GroupMissingReturn    *DiagnosticGroup
	GroupUndefinedVariables *DiagnosticGroup
	GroupDebuggerStatementPresent *DiagnosticGroup
	GroupCheckRegExp      *DiagnosticGroup
	GroupCheckTypes       *DiagnosticGroup
	GroupStrictMissingProperties *DiagnosticGroup
	GroupStrictPrimitiveOperators *DiagnosticGroup
	GroupStrictCheckTypes *DiagnosticGroup
	GroupReportUnknownTypes *DiagnosticGroup
	GroupCheckVariables   *DiagnosticGroup
	GroupCheckUselessCode *DiagnosticGroup
	GroupConst            *DiagnosticGroup
	GroupConstantProperty *DiagnosticGroup
	GroupDuplicateVars    *DiagnosticGroup
	GroupSuspiciousCode   *DiagnosticGroup
	GroupMisplacedTypeAnnotation *DiagnosticGroup
	GroupMisplacedMsgAnnotation  *DiagnosticGroup
	GroupMisplacedSuppress       *DiagnosticGroup
	GroupMissingRequire   *DiagnosticGroup
	GroupTypeImportCodeReferences *DiagnosticGroup
	GroupConformanceViolations    *DiagnosticGroup
	GroupDeprecatedAnnotations    *DiagnosticGroup
	GroupLateProvide              *DiagnosticGroup
	GroupUnrecognizedTypeError    *DiagnosticGroup
	GroupMissingProvide           *DiagnosticGroup
	GroupMissingSourcesWarnings   *DiagnosticGroup
)

func init() {
	GroupGlobalThis = RegisterGroup("globalThis", dtGlobalThis)

	GroupDeprecated = RegisterGroup("deprecated",
		dtDeprecatedName, dtDeprecatedNameReason,
		dtDeprecatedProp, dtDeprecatedPropReason,
		dtDeprecatedClass, dtDeprecatedClassReason)

	GroupUnderscore = RegisterDeprecatedGroup("underscore")

	GroupVisibility = RegisterGroup("visibility",
		dtBadPrivateGlobalAccess,
		dtBadPrivatePropertyAccess,
		dtBadPackagePropertyAccess,
		dtBadProtectedPropertyAccess,
		dtExtendFinalClass,
		dtFinalPropertyOverridden,
		dtPrivateOverride,
		dtVisibilityMismatch)

	GroupAccessControls = RegisterGroupComposite("accessControls", GroupVisibility)

	GroupNonStandardJsDoc = RegisterGroup("nonStandardJsDocs",
		dtBadJsdocAnnotation,
		dtJsdocInBlockComment)

	GroupInvalidCasts = RegisterGroup("invalidCasts", dtInvalidCast)

	GroupStrictModuleDepCheck = RegisterGroup("strictModuleDepCheck",
		dtVarMultiplyDeclaredError, // approximation
		dtCrossChunkRequireError)

	GroupMissingOverride = RegisterGroup("missingOverride",
		dtHiddenInterfacePropertyMismatch,
		dtHiddenSuperclassPropertyMismatch)

	GroupMissingProperties = RegisterGroup("missingProperties",
		dtInexistentProperty,
		dtPossibleInexistentProperty)

	GroupMissingReturn = RegisterGroup("missingReturn", dtMissingReturnStatement)

	GroupUndefinedVariables = RegisterGroup("undefinedVars", dtUndefinedVarError)

	GroupDebuggerStatementPresent = RegisterGroup("checkDebuggerStatement",
		dtDebuggerStatementPresent)

	GroupCheckRegExp = RegisterGroup("checkRegExp",
		dtRegexpReference, dtMalformedRegexp)

	GroupCheckTypes = RegisterGroup("checkTypes",
		dtTypeMismatch, dtInexistentProperty, dtPossibleInexistentProperty,
		dtNotAConstructor, dtIllegalPropertyAccess,
		dtHiddenSuperclassPropertyMismatch, dtHiddenInterfacePropertyMismatch,
		dtDefaultParamMustBeMarkedOptional, dtBadRestParameterAnnotation,
		dtGlobalThis)

	GroupStrictMissingProperties = RegisterGroup("strictMissingProperties",
		dtStrictInexistentProperty, dtIllegalPropertyCreationOnUnion)

	GroupStrictPrimitiveOperators = RegisterGroup("strictPrimitiveOperators",
		dtInvalidOperandType)

	GroupStrictCheckTypes = RegisterGroupComposite("strictCheckTypes",
		GroupStrictMissingProperties, GroupStrictPrimitiveOperators)

	GroupReportUnknownTypes = RegisterGroup("reportUnknownTypes", dtUnknownExprType)

	GroupCheckVariables = RegisterGroup("checkVars",
		dtUndefinedVarError, dtVarMultiplyDeclaredError,
		dtEarlyReference, dtRedeclaredVariable)

	GroupCheckUselessCode = RegisterGroup("uselessCode",
		dtUselessCode, dtUnreachableCode)

	GroupConst = RegisterGroup("const",
		dtConstPropertyDeleted, dtConstPropertyReassigned)

	GroupConstantProperty = RegisterGroup("constantProperty",
		dtConstPropertyDeleted, dtConstPropertyReassigned)

	GroupDuplicateVars = RegisterGroup("duplicate",
		dtDupVarDeclaration, dtDupVarDeclarationTypeMismatch,
		dtFunctionMasksVariable, dtVarMultiplyDeclaredError, dtRedeclaredVariable)

	GroupSuspiciousCode = RegisterGroup("suspiciousCode",
		dtSuspiciousSemicolon, dtSuspiciousBreakingOutOfOptionalChain,
		dtSuspiciousNan, dtSuspiciousIn,
		dtSuspiciousInstanceofLeft, dtSuspiciousLeftOperandOfLogicalOperator,
		dtSuspiciousNegatedLeftOperandOfInOperator,
		dtDeterministicTest)

	GroupMisplacedTypeAnnotation = RegisterGroup("misplacedTypeAnnotation",
		dtArrowFunctionAsConstructor, dtBadRestParameterAnnotation,
		dtDefaultParamMustBeMarkedOptional, dtDisallowedMemberJsdoc,
		dtInvalidNoSideEffectAnnotation, dtInvalidModifiesAnnotation,
		dtJsdocOnReturn, dtMisplacedAnnotation, dtMisplacedMsgAnnotation)

	GroupMisplacedMsgAnnotation = RegisterGroup("misplacedMsgAnnotation",
		dtMisplacedMsgAnnotation)

	GroupMisplacedSuppress = RegisterGroup("misplacedSuppress", dtMisplacedSuppress)

	GroupMissingRequire = RegisterGroup("missingRequire",
		dtMissingRequire, dtMissingRequireInProvidesFile,
		dtMissingRequireInGoogScope, dtMissingRequireTypeInGoogScope,
		dtMissingRequireType, dtMissingRequireTypeInProvidesFile,
		dtIncorrectNamespaceAliasRequire, dtIncorrectNamespaceAliasRequireType,
		dtIndirectNamespaceRefRequire, dtIndirectNamespaceRefRequireType)

	GroupTypeImportCodeReferences = RegisterGroup("typeImportCodeReferences",
		dtTypeImportCodeReference)

	GroupConformanceViolations = RegisterGroup("conformanceViolations",
		dtConformanceViolation, dtConformancePossibleViolation)

	GroupDeprecatedAnnotations = RegisterGroup("deprecatedAnnotations",
		dtAnnotationDeprecated)

	GroupLateProvide = RegisterGroup("lateProvide", dtLateProvidesError)

	GroupUnrecognizedTypeError = RegisterGroup("unrecognizedTypeError",
		dtUnrecognizedTypeError)

	GroupMissingProvide = RegisterGroup("missingProvide", dtMissingModuleOrProvideForFD)

	GroupMissingSourcesWarnings = RegisterGroupComposite("missingSourcesWarnings",
		GroupReportUnknownTypes, GroupUndefinedVariables, GroupMissingProvide,
		GroupMissingProperties, GroupDuplicateVars)
}
