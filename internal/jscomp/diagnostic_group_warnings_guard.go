// Port of com/google/javascript/jscomp/DiagnosticGroupWarningsGuard.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/base"

// DiagnosticGroupWarningsGuard sets the level for a particular
// DiagnosticGroup.
type DiagnosticGroupWarningsGuard struct {
	WarningsGuardBase
	group *DiagnosticGroup
	level CheckLevel
}

var _ WarningsGuard = (*DiagnosticGroupWarningsGuard)(nil)

const diagnosticGroupWarningsGuard_serialVersionUID int64 = 1

// Java: DiagnosticGroupWarningsGuard(DiagnosticGroup, CheckLevel)
func NewDiagnosticGroupWarningsGuard(
	group *DiagnosticGroup, level CheckLevel) *DiagnosticGroupWarningsGuard {
	return &DiagnosticGroupWarningsGuard{group: group, level: level}
}

// Level may return nil.
//
// Java: DiagnosticGroupWarningsGuard.level(JSError)
func (g *DiagnosticGroupWarningsGuard) Level(error *JSError) *CheckLevel {
	if g.group.Matches(error) {
		return g.level.Ptr()
	}
	return nil
}

// Java: DiagnosticGroupWarningsGuard.mustRunChecks(DiagnosticGroup)
func (g *DiagnosticGroupWarningsGuard) MustRunChecks(otherGroup *DiagnosticGroup) base.Tri {
	if g.level.IsOn() {
		for t := range otherGroup.GetTypes().All() {
			if g.group.MatchesType(t) {
				return base.Tri_TRUE
			}
		}
		return base.Tri_UNKNOWN
	} else {
		for t := range otherGroup.GetTypes().All() {
			if !g.group.MatchesType(t) {
				return base.Tri_UNKNOWN
			}
		}
		return base.Tri_FALSE
	}
}

// Java: DiagnosticGroupWarningsGuard.toString()
func (g *DiagnosticGroupWarningsGuard) String() string {
	return g.group.String() + "(" + g.level.String() + ")"
}
