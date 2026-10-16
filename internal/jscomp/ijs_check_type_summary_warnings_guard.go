// Port of com/google/javascript/jscomp/ijs/CheckTypeSummaryWarningsGuard.java @ closure-compiler v20260415

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CheckTypeSummaryWarningsGuard is a warnings guard that demotes the errors
// found in type summary files to be less severe, leaving only the errors
// found in the original source.
type CheckTypeSummaryWarningsGuard struct {
	WarningsGuardBase
	level CheckLevel
}

var _ WarningsGuard = (*CheckTypeSummaryWarningsGuard)(nil)

// Java: CheckTypeSummaryWarningsGuard(CheckLevel)
func NewCheckTypeSummaryWarningsGuard(level CheckLevel) *CheckTypeSummaryWarningsGuard {
	return &CheckTypeSummaryWarningsGuard{level: level}
}

// Level may return nil.
//
// Java: CheckTypeSummaryWarningsGuard.level(JSError)
func (g *CheckTypeSummaryWarningsGuard) Level(error *JSError) *CheckLevel {
	jutil.CheckNotNull(error)
	if !g.ShouldDemoteIfTypeSummary(error) {
		return nil
	}
	if g.InTypeSummary(error) {
		return g.level.Ptr()
	}
	return nil
}

// Java: CheckTypeSummaryWarningsGuard.getPriority()
func (g *CheckTypeSummaryWarningsGuard) GetPriority() int {
	// Treat warnings in .i.js files as though they are allowlisted.
	return WarningsGuard_Priority_SUPPRESS_BY_ALLOWLIST.GetValue()
}

// InTypeSummary returns whether the given error was produced inside a type
// summary file.
//
// Java: CheckTypeSummaryWarningsGuard.inTypeSummary(JSError)
func (g *CheckTypeSummaryWarningsGuard) InTypeSummary(error *JSError) bool {
	return error.SourceName() != nil && strings.HasSuffix(*error.SourceName(), ".i.js")
}

// ShouldDemoteIfTypeSummary returns whether the given error should be
// demoted even if in a type summary file.
//
// Java: CheckTypeSummaryWarningsGuard.shouldDemoteIfTypeSummary(JSError)
func (g *CheckTypeSummaryWarningsGuard) ShouldDemoteIfTypeSummary(error *JSError) bool {
	// TODO(b/215776774): also skip demoting DUPLICATE_MODULE and DUPLICATE_NAMESPACE errors.
	// DUPLICATE_NAMESPACE_AND_MODULE is only special-cased because it was easier to land.
	return !error.Type().Equals(ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE)
}

// String is Object.toString().
func (g *CheckTypeSummaryWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.ijs.CheckTypeSummaryWarningsGuard", g)
}
