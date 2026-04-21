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

type Level int

const (
	Warning Level = iota
	Error
)

func (l Level) String() string {
	if l == Error {
		return "ERROR"
	}
	return "WARNING"
}

type Diagnostic struct {
	Path    string
	Line    int // 1-indexed, matches Java Node.getLineno()
	Col     int // 0-indexed, matches Java Node.getCharno()
	Code    string
	Level   Level
	Message string
}

// FormatJavaLike reproduces the shape of Closure Compiler's default
// output ("<path>:<line>:<col>: <LEVEL> - [<CODE>] <msg>"). Strings are
// not part of the conformance check (see SPEC.md), but matching the shape
// keeps the CLI comfortable to use.
func (d Diagnostic) FormatJavaLike() string {
	return fmt.Sprintf("%s:%d:%d: %s - [%s] %s", d.Path, d.Line, d.Col, d.Level, d.Code, d.Message)
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
	JSC_MISSING_RETURN_STATEMENT    = "JSC_MISSING_RETURN_STATEMENT"
	JSC_USELESS_CODE                = "JSC_USELESS_CODE"
	JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH = "JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH"
	JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH  = "JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH"
	JSC_ILLEGAL_PROPERTY_ACCESS             = "JSC_ILLEGAL_PROPERTY_ACCESS"
	JSC_BAD_PRIVATE_PROPERTY_ACCESS         = "JSC_BAD_PRIVATE_PROPERTY_ACCESS"
	JSC_BAD_PROTECTED_PROPERTY_ACCESS       = "JSC_BAD_PROTECTED_PROPERTY_ACCESS"
	JSC_NOT_A_CONSTRUCTOR                   = "JSC_NOT_A_CONSTRUCTOR"
	JSC_INVALID_OPERAND_TYPE                = "JSC_INVALID_OPERAND_TYPE"
	JSC_INEXISTENT_ENUM_ELEMENT             = "JSC_INEXISTENT_ENUM_ELEMENT"
	JSC_WRONG_ARGUMENT_COUNT                = "JSC_WRONG_ARGUMENT_COUNT"
	JSC_NON_STRINGIFIABLE_OBJECT_KEY        = "JSC_NON_STRINGIFIABLE_OBJECT_KEY"
	JSC_INTERFACE_METHOD_NOT_IMPLEMENTED    = "JSC_INTERFACE_METHOD_NOT_IMPLEMENTED"
	JSC_DETERMINISTIC_TEST                  = "JSC_DETERMINISTIC_TEST"
	JSC_UNRECOGNIZED_TYPE_ERROR             = "JSC_UNRECOGNIZED_TYPE_ERROR"
	JSC_TYPE_PARSE_ERROR                    = "JSC_TYPE_PARSE_ERROR"
	JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE = "JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE"
	JSC_UNRECOGNIZED_TYPEOF_ERROR           = "JSC_UNRECOGNIZED_TYPEOF_ERROR"
	JSC_USED_GLOBAL_THIS                    = "JSC_USED_GLOBAL_THIS"
	JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH   = "JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH"
	JSC_IMPLEMENTS_NON_INTERFACE            = "JSC_IMPLEMENTS_NON_INTERFACE"
	JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS  = "JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS"
	JSC_DUPLICATE_MEMBER                    = "JSC_DUPLICATE_MEMBER"
	JSC_DUP_VAR_DECLARATION                 = "JSC_DUP_VAR_DECLARATION"
	JSC_ILLEGAL_PROPERTY_CREATION           = "JSC_ILLEGAL_PROPERTY_CREATION"
	JSC_CONFLICTING_EXTENDED_TYPE           = "JSC_CONFLICTING_EXTENDED_TYPE"
	JSC_CONSTRUCTOR_REQUIRED                = "JSC_CONSTRUCTOR_REQUIRED"
	JSC_CONSTRUCTOR_NOT_CALLABLE            = "JSC_CONSTRUCTOR_NOT_CALLABLE"
	JSC_IN_USED_WITH_STRUCT                 = "JSC_IN_USED_WITH_STRUCT"
	JSC_DICT_EXTEND_STRUCT_TYPE             = "JSC_DICT_EXTEND_STRUCT_TYPE"
	JSC_REASSIGN_CLASS_PROTOTYPE            = "JSC_REASSIGN_CLASS_PROTOTYPE"
	JSC_BAD_TYPE_FOR_BIT_OPERATION          = "JSC_BAD_TYPE_FOR_BIT_OPERATION"
	JSC_NOT_FUNCTION_TYPE                   = "JSC_NOT_FUNCTION_TYPE"
	JSC_INEXISTENT_PARAM                    = "JSC_INEXISTENT_PARAM"
	// JSC_INCOMPATIBLE_ALIAS_ANNOTATION — emitted by TypedScopeCreator
	// when a `/** @typedef */` alias is re-declared with a different
	// body (e.g. `/** @typedef {number} */ A; /** @typedef {string} */ A;`).
	// Matches the Java TypedScopeCreator INCOMPATIBLE_ALIAS_ANNOTATION
	// registration.
	JSC_INCOMPATIBLE_ALIAS_ANNOTATION  = "JSC_INCOMPATIBLE_ALIAS_ANNOTATION"
	// JSC_CONFLICTING_GETTER_SETTER_TYPE — emitted when a paired get /
	// set accessor pair has incompatible declared read-vs-write types.
	// Mirrors the Java TypedScopeCreator CONFLICTING_GETTER_SETTER_TYPE
	// diagnostic.
	JSC_CONFLICTING_GETTER_SETTER_TYPE = "JSC_CONFLICTING_GETTER_SETTER_TYPE"
)
