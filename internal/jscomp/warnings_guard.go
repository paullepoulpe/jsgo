// Port of com/google/javascript/jscomp/WarningsGuard.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/base"

// WarningsGuard allows to flexibly manage what to do with a reported
// warning/error.
//
// Guard has several choices: - return OFF - suppress the warning/error -
// return WARNING - return ERROR report it with high severity - return null.
// Does not know what to do with it. Lets the other guard decide what to do
// with it.
//
// Although the interface is very simple, it allows you easily customize what
// warnings you are interested in.
//
// For example there are could be several implementations: StrictGuard -
// {return ERROR}. All warnings should be treat as errors. SilentGuard - {if
// (WARNING) return OFF}. Suppress all warnings but still fail if JS has
// errors. AllowlistGuard (if !allowlistErrors.contains(error) return ERROR)
// return error if it does not present in the allowlist.
type WarningsGuard interface {
	// Level returns a new check level for a given error.
	//
	// nil means that this guard does not know what to do with the error. nil
	// can be used it chain multiple guards; if current guard returns nil,
	// then the next in the chain should process it.
	//
	// Java: WarningsGuard.level(JSError)
	Level(error *JSError) *CheckLevel

	// MustRunChecks: do checks for group still need to be run if this guard
	// is installed?
	//
	//  1. TRUE: Enables one or more types in the group, so it must be checked.
	//  2. FALSE: Disables all types in the group, so it need not be checked.
	//  3. UNKNOWN: Does not affect or only partially disables the group, so
	//     checking is undecided.
	//
	// Java: WarningsGuard.mustRunChecks(DiagnosticGroup)
	MustRunChecks(group *DiagnosticGroup) base.Tri

	// GetPriority is the priority in which warnings guards are applied.
	// Lower means the guard will be applied sooner. Expressed on a scale of
	// 1 to 100.
	//
	// Java: WarningsGuard.getPriority()
	GetPriority() int
}

// WarningsGuardBase holds the base implementations of WarningsGuard;
// subclasses embed it.
type WarningsGuardBase struct{}

const warningsGuard_serialVersionUID int64 = 1

// Java: WarningsGuard.mustRunChecks(DiagnosticGroup)
func (WarningsGuardBase) MustRunChecks(group *DiagnosticGroup) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: WarningsGuard.getPriority()
func (WarningsGuardBase) GetPriority() int {
	return WarningsGuard_Priority_DEFAULT.value()
}

// WarningsGuard_Priority is the priority of a guard.
//
// Java: WarningsGuard.Priority
type WarningsGuard_Priority int32

const (
	WarningsGuard_Priority_MAX WarningsGuard_Priority = iota
	WarningsGuard_Priority_MIN
	WarningsGuard_Priority_STRICT
	WarningsGuard_Priority_DEFAULT
	WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST
	WarningsGuard_Priority_SUPPRESS_DOC
	WarningsGuard_Priority_FILTER_BY_PATH
)

var warningsGuard_Priority_names = [...]string{
	"MAX", "MIN", "STRICT", "DEFAULT", "SUPPRESS_BY_ALLOWLIST", "SUPPRESS_DOC", "FILTER_BY_PATH",
}

var warningsGuard_Priority_values = [...]int{1, 100, 100, 50, 40, 20, 1}

// Java: WarningsGuard.Priority.values()
func WarningsGuard_Priority_values() []WarningsGuard_Priority {
	return []WarningsGuard_Priority{
		WarningsGuard_Priority_MAX, WarningsGuard_Priority_MIN, WarningsGuard_Priority_STRICT,
		WarningsGuard_Priority_DEFAULT, WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST,
		WarningsGuard_Priority_SUPPRESS_DOC, WarningsGuard_Priority_FILTER_BY_PATH,
	}
}

// String is Enum.name().
func (p WarningsGuard_Priority) String() string { return warningsGuard_Priority_names[p] }

// Ordinal is Enum.ordinal().
func (p WarningsGuard_Priority) Ordinal() int { return int(p) }

// value is the field Priority.value.
func (p WarningsGuard_Priority) value() int { return warningsGuard_Priority_values[p] }

// Java: WarningsGuard.Priority.getValue()
func (p WarningsGuard_Priority) GetValue() int { return p.value() }
