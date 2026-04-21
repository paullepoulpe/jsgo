// Port of com/google/javascript/jscomp/CompilerOptions.java
// @ closure-compiler v20260415 (narrow; only the knobs the current
// passes consult).
//
// Java's CompilerOptions is a ~2000-line bag of flags covering every
// command-line switch (language in/out, externs handling, pass
// toggles, warning levels, …). This port carries the minimum jsgo
// needs today: a per-DiagnosticGroup CheckLevel map so a pass can
// consult `opts.LevelFor(group)` before emitting. Defaults match
// Java's `--jscomp_warning=<name>` defaults:
//   - checkTypes        → ERROR   (the type checker is on)
//   - strictCheckTypes  → OFF     (the stricter operand / comparison
//                                  checks stay silent unless opted in)
//
// New flags extend this file as later tiers wire them in.

package jscomp

// CompilerOptions holds the subset of --jscomp_* knobs jsgo honours.
// Additional fields (language-in, module resolution, path renaming,
// …) land as later tiers plumb them through the Compiler shim.
type CompilerOptions struct {
	levels map[*DiagnosticGroup]CheckLevel
}

// NewCompilerOptions returns options pre-populated with Java's
// `CompilerOptions#getDefaultOptions()` defaults for every group
// this port knows about. Callers tweak individual levels via
// SetWarningLevel / SetErrorLevel / SetOff.
func NewCompilerOptions() *CompilerOptions {
	o := &CompilerOptions{levels: map[*DiagnosticGroup]CheckLevel{}}
	o.levels[GroupCheckTypes] = LevelError
	o.levels[GroupStrictCheckTypes] = LevelOff
	return o
}

// SetWarningLevel mirrors CompilerOptions.setWarningLevel(group, lvl)
// in Java — unconditionally sets the level for the group. Absent
// groups default to LevelOff on lookup.
func (o *CompilerOptions) SetWarningLevel(g *DiagnosticGroup, lvl CheckLevel) {
	if o == nil {
		return
	}
	if o.levels == nil {
		o.levels = map[*DiagnosticGroup]CheckLevel{}
	}
	o.levels[g] = lvl
}

// LevelFor returns the current level configured for `g`, or LevelOff
// if the group was never registered on this options instance.
// Matches CompilerOptions.getLevelForGroup(group) for unknown groups.
func (o *CompilerOptions) LevelFor(g *DiagnosticGroup) CheckLevel {
	if o == nil || o.levels == nil {
		return LevelOff
	}
	return o.levels[g]
}

// IsOn reports whether the group is set to WARNING or ERROR. The
// common call pattern inside passes is `if opts.IsOn(group) { … }`
// to gate a whole validator branch.
func (o *CompilerOptions) IsOn(g *DiagnosticGroup) bool {
	return o.LevelFor(g) != LevelOff
}
