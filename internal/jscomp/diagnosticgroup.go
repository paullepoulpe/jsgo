// Port of com/google/javascript/jscomp/DiagnosticGroup.java +
// DiagnosticGroups.java @ closure-compiler v20260415 (narrow; only
// the groups jsgo's current passes consult).
//
// Closure's DiagnosticGroups is a ~1500-line registry mapping
// well-known group names (`checkTypes`, `strictCheckTypes`, `const`,
// …) to sets of DiagnosticTypes; CompilerOptions stores a per-group
// CheckLevel (ERROR / WARNING / OFF) and every diagnostic emission
// consults `getLevelForGroup(...)` before surfacing it. We port only
// the minimum the current passes need:
//
//   - STRICT_CHECK_TYPES     — gate for arithmetic / bitwise /
//     comparison operand checks that Java's TypeValidator emits only
//     under `--jscomp_warning=strictCheckTypes` (OFF by default, so
//     the extracted corpus's Java oracle doesn't contain them).
//   - CHECK_TYPES            — gate for the core TypeCheck
//     diagnostics (ON by default when the type checker runs).
//
// New groups extend this file as later tiers wire them in.

package jscomp

// DiagnosticGroup is a named set of diagnostic codes that
// CompilerOptions can toggle as one unit. Matches Java's
// DiagnosticGroup class; a group carries a stable string name plus
// the codes it covers.
type DiagnosticGroup struct {
	Name  string
	Codes []string
}

// CheckLevel mirrors com/google/javascript/jscomp/CheckLevel.java.
// ERROR: emit as an Error-level diagnostic. WARNING: emit as
// Warning-level. OFF: suppress entirely.
type CheckLevel int

const (
	LevelOff CheckLevel = iota
	LevelWarning
	LevelError
)

// Well-known diagnostic groups Java surfaces via
// `--jscomp_warning=<name>`. Names match the Java identifiers the
// corpus relies on.
var (
	// GroupCheckTypes is ON by default when the type checker runs.
	// It covers every JSC_TYPE_MISMATCH / JSC_INEXISTENT_PROPERTY /
	// JSC_NOT_A_CONSTRUCTOR / JSC_ILLEGAL_PROPERTY_ACCESS site.
	GroupCheckTypes = &DiagnosticGroup{
		Name: "checkTypes",
		Codes: []string{
			"JSC_TYPE_MISMATCH",
			"JSC_INEXISTENT_PROPERTY",
			"JSC_NOT_A_CONSTRUCTOR",
			"JSC_ILLEGAL_PROPERTY_ACCESS",
			"JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH",
			"JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH",
			"JSC_USED_GLOBAL_THIS",
		},
	}
	// GroupStrictCheckTypes is OFF by default. Enabling it turns on
	// the stricter per-operand checks in TypeValidator — the
	// arithmetic / bitwise / relative-comparison rules that fire as
	// JSC_INVALID_OPERAND_TYPE (WARNING) under Java's
	// `--jscomp_warning=strictCheckTypes`.
	GroupStrictCheckTypes = &DiagnosticGroup{
		Name: "strictCheckTypes",
		Codes: []string{
			"JSC_INVALID_OPERAND_TYPE",
		},
	}
)
