// Port of com/google/javascript/jscomp/JSError.java +
// com/google/javascript/jscomp/DiagnosticType.java
// @ closure-compiler v20260415 (subset).
//
// Java's JSError carries position, a DiagnosticType (code +
// default-level + format template), substituted args, and links to
// the source file. DiagnosticType is a full registry with
// per-category warning levels. This file tracks only the minimum
// the tuple-level conformance check consumes: (path, line, col,
// code) plus a Java-shaped message so CLI output reads the same.
// A proper DiagnosticType registry lands when we want to honour
// `--jscomp_warning=<group>` style flags.

package diag

import "fmt"

// Go-only: Level is a simplified severity enum used by the pre-JSError diagnostics pipeline.
// Java uses CheckLevel (in diagnostictype.go) for the same purpose.
type Level int

const (
	Warning Level = iota
	Error
)

// Go-only: Level.String is a convenience method; Java uses CheckLevel.toString().
func (l Level) String() string {
	if l == Error {
		return "ERROR"
	}
	return "WARNING"
}

// Go-only: Diagnostic is a lightweight struct used by the pre-JSError diagnostics pipeline.
// The Java equivalent is JSError (see jserror.go).
type Diagnostic struct {
	Path    string
	Line    int // 1-indexed, matches Java Node.getLineno()
	Col     int // 0-indexed, matches Java Node.getCharno()
	Code    string
	Level   Level
	Message string
}

// Go-only: FormatJavaLike reproduces the shape of Closure Compiler's default
// output ("<path>:<line>:<col>: <LEVEL> - [<CODE>] <msg>"). Strings are
// not part of the conformance check (see SPEC.md), but matching the shape
// keeps the CLI comfortable to use.
func (d Diagnostic) FormatJavaLike() string {
	return fmt.Sprintf("%s:%d:%d: %s - [%s] %s", d.Path, d.Line, d.Col, d.Level, d.Code, d.Message)
}

// GetDescription mirrors JSError.getDescription() — returns the formatted error message.
func (d Diagnostic) GetDescription() string { return d.Message }

// GetSourceName mirrors JSError.getSourceName() — returns the source file name.
func (d Diagnostic) GetSourceName() string { return d.Path }

// GetLineno mirrors JSError.getLineno() — returns the 1-indexed line number.
func (d Diagnostic) GetLineno() int { return d.Line }

// GetCharno mirrors JSError.getCharno() — returns the 0-indexed character offset within the line.
func (d Diagnostic) GetCharno() int { return d.Col }

// GetLineNumber mirrors JSError.getLineNumber() — alias for GetLineno().
func (d Diagnostic) GetLineNumber() int { return d.Line }

// ToString mirrors JSError.toString() — returns a human-readable diagnostic string.
func (d Diagnostic) ToString() string {
	srcName := d.Path
	if srcName == "" {
		srcName = "(unknown source)"
	}
	lineno := "(unknown line)"
	if d.Line != 0 {
		lineno = fmt.Sprintf("%d", d.Line)
	}
	charno := "(unknown column)"
	if d.Col >= 0 {
		charno = fmt.Sprintf("%d", d.Col)
	}
	return d.Code + ". " + d.Message + " at " + srcName + " line " + lineno + " : " + charno
}

// Well-known diagnostic codes. Add new ones as the corpus grows. Codes
// must match the JSC_* identifiers used by Java Closure Compiler (see
// DiagnosticType registrations in TypeCheck.java / VarCheck.java /
// VariableReferenceCheck.java).
const (
	JSC_TYPE_MISMATCH               = "JSC_TYPE_MISMATCH"
	JSC_INEXISTENT_PROPERTY         = "JSC_INEXISTENT_PROPERTY"
	JSC_UNDEFINED_VARIABLE          = "JSC_UNDEFINED_VARIABLE"
	JSC_VAR_MULTIPLY_DECLARED_ERROR = "JSC_VAR_MULTIPLY_DECLARED_ERROR"
	JSC_FUNCTION_MASKS_VARIABLE     = "JSC_FUNCTION_MASKS_VARIABLE"
	JSC_MISSING_RETURN_STATEMENT    = "JSC_MISSING_RETURN_STATEMENT"
	JSC_USELESS_CODE                = "JSC_USELESS_CODE"
	JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH = "JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH"
	JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH  = "JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH"
	JSC_ILLEGAL_PROPERTY_ACCESS             = "JSC_ILLEGAL_PROPERTY_ACCESS"
	JSC_NOT_A_CONSTRUCTOR                   = "JSC_NOT_A_CONSTRUCTOR"
	JSC_INVALID_OPERAND_TYPE                = "JSC_INVALID_OPERAND_TYPE"
	JSC_INVALID_CAST                        = "JSC_INVALID_CAST"
	JSC_INEXISTENT_ENUM_ELEMENT             = "JSC_INEXISTENT_ENUM_ELEMENT"
	JSC_WRONG_ARGUMENT_COUNT                = "JSC_WRONG_ARGUMENT_COUNT"
	JSC_NON_STRINGIFIABLE_OBJECT_KEY        = "JSC_NON_STRINGIFIABLE_OBJECT_KEY"
	JSC_INTERFACE_METHOD_NOT_IMPLEMENTED    = "JSC_INTERFACE_METHOD_NOT_IMPLEMENTED"
	JSC_DETERMINISTIC_TEST                  = "JSC_DETERMINISTIC_TEST"
	JSC_TYPE_PARSE_ERROR                    = "JSC_TYPE_PARSE_ERROR"
	JSC_UNRECOGNIZED_TYPE_ERROR             = "JSC_UNRECOGNIZED_TYPE_ERROR"
	JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE = "JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE"
	JSC_UNRECOGNIZED_TYPEOF_ERROR           = "JSC_UNRECOGNIZED_TYPEOF_ERROR"
	JSC_USED_GLOBAL_THIS                    = "JSC_USED_GLOBAL_THIS"
	JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH   = "JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH"
	JSC_IMPLEMENTS_NON_INTERFACE            = "JSC_IMPLEMENTS_NON_INTERFACE"
	JSC_EXTENDS_NON_OBJECT                 = "JSC_EXTENDS_NON_OBJECT"
	JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS  = "JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS"
	JSC_DUPLICATE_MEMBER                    = "JSC_DUPLICATE_MEMBER"
	JSC_DUP_VAR_DECLARATION                 = "JSC_DUP_VAR_DECLARATION"
	JSC_INCOMPATIBLE_ALIAS_ANNOTATION       = "JSC_INCOMPATIBLE_ALIAS_ANNOTATION"
	JSC_POSSIBLE_INEXISTENT_PROPERTY        = "JSC_POSSIBLE_INEXISTENT_PROPERTY"
	JSC_ILLEGAL_PROPERTY_CREATION           = "JSC_ILLEGAL_PROPERTY_CREATION"
	JSC_CONFLICTING_EXTENDED_TYPE           = "JSC_CONFLICTING_EXTENDED_TYPE"
	JSC_CONSTRUCTOR_REQUIRED                = "JSC_CONSTRUCTOR_REQUIRED"
	JSC_TEMPLATE_TRANSFORMATION_ON_CLASS    = "JSC_TEMPLATE_TRANSFORMATION_ON_CLASS"
	JSC_EXPECTED_THIS_TYPE                  = "JSC_EXPECTED_THIS_TYPE"
	JSC_IN_USED_WITH_STRUCT                 = "JSC_IN_USED_WITH_STRUCT"
	JSC_DICT_EXTEND_STRUCT_TYPE             = "JSC_DICT_EXTEND_STRUCT_TYPE"
	JSC_REASSIGN_CLASS_PROTOTYPE                   = "JSC_REASSIGN_CLASS_PROTOTYPE"
	JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE      = "JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE"
	JSC_BAD_TYPE_FOR_BIT_OPERATION          = "JSC_BAD_TYPE_FOR_BIT_OPERATION"
	JSC_NOT_FUNCTION_TYPE                   = "JSC_NOT_FUNCTION_TYPE"
	JSC_INEXISTENT_PARAM                    = "JSC_INEXISTENT_PARAM"
	JSC_UNKNOWN_OVERRIDE                    = "JSC_UNKNOWN_OVERRIDE"
	JSC_CONSTRUCTOR_NOT_CALLABLE            = "JSC_CONSTRUCTOR_NOT_CALLABLE"
	JSC_INVALID_PARAM                       = "JSC_INVALID_PARAM"
	JSC_INSTANTIATE_ABSTRACT_CLASS          = "JSC_INSTANTIATE_ABSTRACT_CLASS"
	JSC_DEBUGGER_STATEMENT_PRESENT          = "JSC_DEBUGGER_STATEMENT_PRESENT"
	JSC_REGEXP_REFERENCE                    = "JSC_REGEXP_REFERENCE"
	JSC_MALFORMED_REGEXP                    = "JSC_MALFORMED_REGEXP"
	JSC_UNREACHABLE_CODE                    = "JSC_UNREACHABLE_CODE"
	JSC_ILLEGAL_OBJLIT_KEY                  = "JSC_ILLEGAL_OBJLIT_KEY"
	JSC_ILLEGAL_CLASS_KEY                   = "JSC_ILLEGAL_CLASS_KEY"
	JSC_CYCLIC_INHERITANCE_ERROR            = "JSC_CYCLIC_INHERITANCE_ERROR"
	JSC_ENUM_INITIALIZER_NOT_ENUM           = "JSC_ENUM_INITIALIZER_NOT_ENUM"
	// RhinoErrorReporter / JsDocInfoParser: bounded-generic (@template {B} T)
	// is still in development. Port of RhinoErrorReporter.UNSUPPORTED_BOUNDED_GENERIC_TYPES.
	JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES = "JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES"

	// CheckSuspiciousCode diagnostics (all warnings).
	JSC_SUSPICIOUS_SEMICOLON                             = "JSC_SUSPICIOUS_SEMICOLON"
	JSC_SUSPICIOUS_NAN                                   = "JSC_SUSPICIOUS_NAN"
	JSC_SUSPICIOUS_IN                                    = "JSC_SUSPICIOUS_IN"
	JSC_SUSPICIOUS_INSTANCEOF_LEFT                       = "JSC_SUSPICIOUS_INSTANCEOF_LEFT"
	JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR      = "JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR"
	JSC_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR   = "JSC_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR"
	SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN            = "SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN"

	// CheckSuper diagnostics.
	JSC_MISSING_CALL_TO_SUPER                 = "JSC_MISSING_CALL_TO_SUPER"
	JSC_THIS_BEFORE_SUPER                     = "JSC_THIS_BEFORE_SUPER"
	JSC_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR = "JSC_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR"
	JSC_INVALID_SUPER_CALL                    = "JSC_INVALID_SUPER_CALL"
	JSC_SUPER_CALL_IN_ARROW                   = "JSC_SUPER_CALL_IN_ARROW"
	JSC_INVALID_SUPER_USAGE                   = "JSC_INVALID_SUPER_USAGE"
	JSC_INVALID_SUPER_ACCESS                  = "JSC_INVALID_SUPER_ACCESS"
	JSC_INVALID_SUPER_CALL_WITH_SUGGESTION    = "JSC_INVALID_SUPER_CALL_WITH_SUGGESTION"

	// CheckAccessControls diagnostics.
	JSC_BAD_PRIVATE_PROPERTY_ACCESS        = "JSC_BAD_PRIVATE_PROPERTY_ACCESS"
	JSC_BAD_PROTECTED_PROPERTY_ACCESS      = "JSC_BAD_PROTECTED_PROPERTY_ACCESS"
	JSC_BAD_PACKAGE_PROPERTY_ACCESS        = "JSC_BAD_PACKAGE_PROPERTY_ACCESS"
	JSC_BAD_PRIVATE_GLOBAL_ACCESS          = "JSC_BAD_PRIVATE_GLOBAL_ACCESS"
	JSC_PRIVATE_OVERRIDE                   = "JSC_PRIVATE_OVERRIDE"
	JSC_VISIBILITY_MISMATCH                = "JSC_VISIBILITY_MISMATCH"
	JSC_DEPRECATED_VAR                     = "JSC_DEPRECATED_VAR"
	JSC_DEPRECATED_VAR_REASON              = "JSC_DEPRECATED_VAR_REASON"
	JSC_DEPRECATED_PROP                    = "JSC_DEPRECATED_PROP"
	JSC_DEPRECATED_PROP_REASON             = "JSC_DEPRECATED_PROP_REASON"
	JSC_DEPRECATED_CLASS                   = "JSC_DEPRECATED_CLASS"
	JSC_DEPRECATED_CLASS_REASON            = "JSC_DEPRECATED_CLASS_REASON"
	JSC_CONSTANT_PROPERTY_REASSIGNED_VALUE = "JSC_CONSTANT_PROPERTY_REASSIGNED_VALUE"
	JSC_FINAL_PROPERTY_OVERRIDDEN          = "JSC_FINAL_PROPERTY_OVERRIDDEN"
	JSC_EXTEND_FINAL_CLASS                 = "JSC_EXTEND_FINAL_CLASS"
	JSC_CONSTANT_PROPERTY_DELETED          = "JSC_CONSTANT_PROPERTY_DELETED"

	// CheckClosureImports diagnostics (from CheckClosureImports.java).
	JSC_INVALID_CLOSURE_IMPORT_DESTRUCTURING         = "JSC_INVALID_CLOSURE_IMPORT_DESTRUCTURING"
	JSC_ONE_CLOSURE_IMPORT_PER_DECLARATION           = "JSC_ONE_CLOSURE_IMPORT_PER_DECLARATION"
	JSC_INVALID_CLOSURE_IMPORT_CALL                  = "JSC_INVALID_CLOSURE_IMPORT_CALL"
	JSC_LATE_PROVIDE_ERROR                           = "JSC_LATE_PROVIDE_ERROR"
	JSC_LET_CLOSURE_IMPORT                           = "JSC_LET_CLOSURE_IMPORT"
	JSC_NO_CLOSURE_IMPORT_DESTRUCTURING              = "JSC_NO_CLOSURE_IMPORT_DESTRUCTURING"
	JSC_LHS_OF_CLOUSRE_IMPORT_MUST_BE_CONST_IN_ES_MODULE = "JSC_LHS_OF_CLOUSRE_IMPORT_MUST_BE_CONST_IN_ES_MODULE"
	JSC_XMODULE_REQUIRE_ERROR                        = "JSC_XMODULE_REQUIRE_ERROR"

	// Closure primitive errors referenced by CheckClosureImports.
	JSC_INVALID_CLOSURE_CALL_ERROR                       = "JSC_INVALID_CLOSURE_CALL_ERROR"
	JSC_MODULE_USES_GOOG_MODULE_GET                      = "JSC_MODULE_USES_GOOG_MODULE_GET"
	JSC_MISSING_MODULE_OR_PROVIDE                        = "JSC_MISSING_MODULE_OR_PROVIDE"
	JSC_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE    = "JSC_MISSING_MODULE_OR_PROVIDE_FOR_FORWARD_DECLARE"
	JSC_GOOG_MODULE_INVALID_GET_CALL_SCOPE               = "JSC_GOOG_MODULE_INVALID_GET_CALL_SCOPE"
	JSC_GOOG_MODULE_INVALID_GET_ALIAS                    = "JSC_GOOG_MODULE_INVALID_GET_ALIAS"

	// ClosureCheckModule diagnostic.
	JSC_INCORRECT_SHORTNAME_CAPITALIZATION = "JSC_INCORRECT_SHORTNAME_CAPITALIZATION"

	// Es6RewriteModules diagnostic.
	JSC_SHOULD_IMPORT_ES6_MODULE = "JSC_SHOULD_IMPORT_ES6_MODULE"

	// Module path resolution error (ModuleLoader).
	JSC_INVALID_MODULE_PATH = "JSC_INVALID_MODULE_PATH"

	// CheckJSDoc diagnostics.
	JSC_ILLEGAL_IMPLICIT_CAST                 = "JSC_ILLEGAL_IMPLICIT_CAST"
	JSC_MISPLACED_MSG_ANNOTATION              = "JSC_MISPLACED_MSG_ANNOTATION"
	JSC_MISPLACED_ANNOTATION                  = "JSC_MISPLACED_ANNOTATION"
	JSC_ANNOTATION_DEPRECATED                 = "JSC_ANNOTATION_DEPRECATED"
	JSC_DISALLOWED_MEMBER_JSDOC               = "JSC_DISALLOWED_MEMBER_JSDOC"
	JSC_ARROW_FUNCTION_AS_CONSTRUCTOR         = "JSC_ARROW_FUNCTION_AS_CONSTRUCTOR"
	JSC_BAD_REST_PARAMETER_ANNOTATION         = "JSC_BAD_REST_PARAMETER_ANNOTATION"
	JSC_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL = "JSC_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL"
	JSC_INVALID_NO_SIDE_EFFECT_ANNOTATION     = "JSC_INVALID_NO_SIDE_EFFECT_ANNOTATION"
	JSC_INVALID_MODIFIES_ANNOTATION           = "JSC_INVALID_MODIFIES_ANNOTATION"
	JSC_INVALID_DEFINE_ON_LET                 = "JSC_INVALID_DEFINE_ON_LET"
	JSC_MISPLACED_SUPPRESS                    = "JSC_MISPLACED_SUPPRESS"
	JSC_JSDOC_IN_BLOCK_COMMENT                = "JSC_JSDOC_IN_BLOCK_COMMENT"
	JSC_JSDOC_ON_RETURN                       = "JSC_JSDOC_ON_RETURN"

	// CheckMissingRequires codes.
	JSC_MISSING_REQUIRE                          = "JSC_MISSING_REQUIRE"
	JSC_MISSING_REQUIRE_TYPE                     = "JSC_MISSING_REQUIRE_TYPE"
	JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE        = "JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE"
	JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE   = "JSC_INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE"
	JSC_INDIRECT_NAMESPACE_REF_REQUIRE           = "JSC_INDIRECT_NAMESPACE_REF_REQUIRE"
	JSC_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE      = "JSC_INDIRECT_NAMESPACE_REF_REQUIRE_TYPE"
	JSC_MISSING_REQUIRE_IN_PROVIDES_FILE         = "JSC_MISSING_REQUIRE_IN_PROVIDES_FILE"
	JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE    = "JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE"
	JSC_NON_LEGACY_GOOG_MODULE_REFERENCE         = "JSC_NON_LEGACY_GOOG_MODULE_REFERENCE"
	JSC_MISSING_REQUIRE_IN_GOOG_SCOPE            = "JSC_MISSING_REQUIRE_IN_GOOG_SCOPE"
	JSC_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE       = "JSC_MISSING_REQUIRE_TYPE_IN_GOOG_SCOPE"
	JSC_MISSING_REQUIRE_FOR_GOOG_MODULE_GET      = "JSC_MISSING_REQUIRE_FOR_GOOG_MODULE_GET"
	JSC_TYPE_IMPORT_CODE_REFERENCE               = "JSC_TYPE_IMPORT_CODE_REFERENCE"
	JSC_REFERENCE_BEFORE_DECLARE                 = "JSC_REFERENCE_BEFORE_DECLARE"

	// TypeCheck / TypeValidator diagnostics.
	JSC_MULTIPLE_VAR_DEF     = "JSC_MULTIPLE_VAR_DEF"
	JSC_UNKNOWN_TYPEOF_VALUE = "JSC_UNKNOWN_TYPEOF_VALUE"

	// SuperclassCheck / TypedScopeCreator: @extends tag resolves to nothing.
	JSC_RESOLVED_TAG_EMPTY = "JSC_RESOLVED_TAG_EMPTY"

	// TypedScopeCreator / TypeCheck diagnostics for getter/setter and @lends.
	JSC_CONFLICTING_GETTER_SETTER_TYPE = "JSC_CONFLICTING_GETTER_SETTER_TYPE"
	JSC_UNKNOWN_LENDS                  = "JSC_UNKNOWN_LENDS"
	JSC_LENDS_ON_NON_OBJECT            = "JSC_LENDS_ON_NON_OBJECT"

	// TypeCheck: abstract-method diagnostics.
	JSC_ABSTRACT_METHOD_IN_CONCRETE_CLASS = "JSC_ABSTRACT_METHOD_IN_CONCRETE_CLASS"
	JSC_ABSTRACT_SUPER_METHOD_NOT_USABLE  = "JSC_ABSTRACT_SUPER_METHOD_NOT_USABLE"
)
