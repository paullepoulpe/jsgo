// Port of com/google/javascript/jscomp/WarningsGuard.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/StrictWarningsGuard.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ClosureUnawareCodeWarningsGuard.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ijs/CheckTypeSummaryWarningsGuard.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ShowByPathWarningsGuard.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/J2clSuppressWarningsGuard.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ByPathWarningsGuard.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/DiagnosticGroupWarningsGuard.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/DiagnosticGroupPathSuppressingWarningsGuard.java @ closure-compiler v20261021

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

/*
 * Copyright 2016 The Closure Compiler Authors.
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

/*
 * Copyright 2018 The Closure Compiler Authors.
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

/*
 * Copyright 2019 The Closure Compiler Authors.
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

/*
 * Copyright 2024 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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
type WarningsGuardBase struct {
	// Go-only: language-adapter — Java compares guards by identity
	// (ComposeWarningsGuard keys them in a map and orders them in a
	// TreeSet), but pointers to distinct zero-size Go values may be equal,
	// so the base keeps every guard (e.g. StrictWarningsGuard) non-zero-size.
	_ byte
}

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

// String is Enum.name().
func (p WarningsGuard_Priority) String() string { return warningsGuard_Priority_names[p] }

// Ordinal is Enum.ordinal().
func (p WarningsGuard_Priority) Ordinal() int { return int(p) }

// value is the field Priority.value.
func (p WarningsGuard_Priority) value() int { return warningsGuard_Priority_values[p] }

// Java: WarningsGuard.Priority.getValue()
func (p WarningsGuard_Priority) GetValue() int { return p.value() }

func warningsGuardGroupName(s string) *string { return &s }

// StrictWarningsGuard: all warnings should be reported as errors.
type StrictWarningsGuard struct {
	WarningsGuardBase
}

var _ WarningsGuard = (*StrictWarningsGuard)(nil)

var StrictWarningsGuard_UNRAISABLE_WARNING = DiagnosticType_Warning("JSC_UNRAISABLE_WARNING", "{0}")

// Java: StrictWarningsGuard()
func NewStrictWarningsGuard() *StrictWarningsGuard {
	return &StrictWarningsGuard{}
}

// Level may return nil.
//
// Java: StrictWarningsGuard.level(JSError)
func (g *StrictWarningsGuard) Level(error *JSError) *CheckLevel {
	if error.Type() == StrictWarningsGuard_UNRAISABLE_WARNING {
		return nil
	}
	if error.DefaultLevel().IsOn() {
		return CheckLevel_ERROR.Ptr()
	}
	return nil
}

// Java: StrictWarningsGuard.getPriority()
func (g *StrictWarningsGuard) GetPriority() int {
	return WarningsGuard_Priority_STRICT.value() // applied last
}

// String is Object.toString().
func (g *StrictWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.StrictWarningsGuard", g)
}

// ClosureUnawareCodeWarningsGuard is a warnings guard that suppresses
// warnings that are spurious for code that is unaware of Closure compiler's
// requirements.
type ClosureUnawareCodeWarningsGuard struct {
	WarningsGuardBase
}

var _ WarningsGuard = (*ClosureUnawareCodeWarningsGuard)(nil)

// Java: ClosureUnawareCodeWarningsGuard()
func NewClosureUnawareCodeWarningsGuard() *ClosureUnawareCodeWarningsGuard {
	return &ClosureUnawareCodeWarningsGuard{}
}

var closureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS = NewDiagnosticGroupWithName(
	warningsGuardGroupName("closureUnawareCodeJSDocIncompatible"),
	RhinoErrorReporter_TYPE_PARSE_ERROR,
	RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR,
	RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING,
	RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING,
	RhinoErrorReporter_UNNECESSARY_ESCAPE,
	RhinoErrorReporter_BAD_JSDOC_ANNOTATION,
	RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES,
	RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR)

// Level may return nil.
//
// Java: ClosureUnawareCodeWarningsGuard.level(JSError)
func (g *ClosureUnawareCodeWarningsGuard) Level(error *JSError) *CheckLevel {
	if !closureUnawareCodeWarningsGuard_DEFAULT_CLOSURE_UNAWARE_CODE_SUPPRESSIONS.Matches(error) {
		return nil
	}
	node := error.Node()
	if node.IsNil() {
		return nil
	}
	if node.IsClosureUnawareCode() {
		return CheckLevel_OFF.Ptr()
	}
	return nil
}

// Java: ClosureUnawareCodeWarningsGuard.getPriority()
func (g *ClosureUnawareCodeWarningsGuard) GetPriority() int {
	return WarningsGuard_Priority_MAX.GetValue()
}

// String is Object.toString().
func (g *ClosureUnawareCodeWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.ClosureUnawareCodeWarningsGuard", g)
}

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
	src := error.SourceName()
	return src != nil && strings.HasSuffix(*src, ".i.js")
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

// ByPathWarningsGuard is an implementation of a WarningsGuard that can
// modify the CheckLevel based on the file that caused the warning, and
// whether this file matches a set of paths (specified either as include or
// exclude of path name parts).
//
// For example:
//
//	List<String> paths = new ArrayList<String>();
//	paths.add("foo");
//	WarningsGuard guard =
//	    ByPathWarningsGuard.forPath(paths, CheckLevel.ERROR, 1);
//
// This guard will convert any warning that came from a file that contains
// "foo" in its path to an error.
type ByPathWarningsGuard struct {
	WarningsGuardBase
	paths    []string
	include  bool
	priority int
	level    CheckLevel
}

var _ WarningsGuard = (*ByPathWarningsGuard)(nil)

// newByPathWarningsGuard constructs a new instance. The priority is
// determined by the CheckLevel: ERROR have Priority.STRICT, and OFF have
// priority FILTER_BY_PATH.
//
// Use ByPathWarningsGuard_ForPath or ByPathWarningsGuard_ExceptPath to
// actually create a new instance.
//
// Java: ByPathWarningsGuard(List<String>, boolean, CheckLevel)
func newByPathWarningsGuard(
	paths []string, include bool, level CheckLevel) *ByPathWarningsGuard {
	jutil.CheckArgument(paths != nil)
	jutil.CheckArgument(level == CheckLevel_OFF || level == CheckLevel_ERROR)
	g := &ByPathWarningsGuard{paths: paths, include: include, level: level}
	if level == CheckLevel_ERROR {
		g.priority = WarningsGuard_Priority_STRICT.value()
	} else {
		g.priority = WarningsGuard_Priority_FILTER_BY_PATH.value()
	}
	return g
}

// ByPathWarningsGuard_ForPath returns a new ByPathWarningsGuard that would
// affect any file in the given set of paths (paths for matching), applying
// level to the affected files.
//
// Java: ByPathWarningsGuard.forPath(List<String>, CheckLevel)
func ByPathWarningsGuard_ForPath(
	paths []string, level CheckLevel) *ByPathWarningsGuard {
	return newByPathWarningsGuard(paths, true, level)
}

// ByPathWarningsGuard_ExceptPath returns a new ByPathWarningsGuard that
// would affect any file not in the given set of paths (paths for matching),
// applying level to the affected files.
//
// Java: ByPathWarningsGuard.exceptPath(List<String>, CheckLevel)
func ByPathWarningsGuard_ExceptPath(
	paths []string, level CheckLevel) *ByPathWarningsGuard {
	return newByPathWarningsGuard(paths, false, level)
}

// Level may return nil.
//
// Java: ByPathWarningsGuard.level(JSError)
func (g *ByPathWarningsGuard) Level(error *JSError) *CheckLevel {
	errorPath := error.SourceName()
	defaultLevel := error.DefaultLevel()
	if defaultLevel != CheckLevel_ERROR && errorPath != nil {
		src := *errorPath
		inPath := false
		for _, path := range g.paths {
			if strings.Contains(src, path) {
				inPath = true
				break
			}
		}
		if inPath == g.include {
			return g.level.Ptr()
		}
	}
	return nil
}

// Java: ByPathWarningsGuard.getPriority()
func (g *ByPathWarningsGuard) GetPriority() int {
	return g.priority
}

// String is Object.toString().
func (g *ByPathWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.ByPathWarningsGuard", g)
}

// ShowByPathWarningsGuard controls whether warnings should be restricted or
// suppressed for specified paths.
type ShowByPathWarningsGuard struct {
	WarningsGuardBase
	warningsGuard *ByPathWarningsGuard
}

var _ WarningsGuard = (*ShowByPathWarningsGuard)(nil)

// ShowByPathWarningsGuard_ShowType controls whether warnings should be
// restricted to a specified path or suppressed within the specified path.
//
// Java: ShowByPathWarningsGuard.ShowType
type ShowByPathWarningsGuard_ShowType int32

const (
	ShowByPathWarningsGuard_ShowType_INCLUDE ShowByPathWarningsGuard_ShowType = iota // Suppress warnings outside the path.
	ShowByPathWarningsGuard_ShowType_EXCLUDE                                         // Suppress warnings within the path.
)

var showByPathWarningsGuard_ShowType_names = [...]string{"INCLUDE", "EXCLUDE"}

// String is Enum.name().
func (t ShowByPathWarningsGuard_ShowType) String() string {
	return showByPathWarningsGuard_ShowType_names[t]
}

// Java: ShowByPathWarningsGuard(String)
func NewShowByPathWarningsGuard(checkWarningsOnlyForPath string) *ShowByPathWarningsGuard {
	return NewShowByPathWarningsGuardWithShowType(checkWarningsOnlyForPath, ShowByPathWarningsGuard_ShowType_INCLUDE)
}

// Java: ShowByPathWarningsGuard(String[])
func NewShowByPathWarningsGuardWithPaths(checkWarningsOnlyForPath []string) *ShowByPathWarningsGuard {
	return NewShowByPathWarningsGuardWithPathsAndShowType(checkWarningsOnlyForPath, ShowByPathWarningsGuard_ShowType_INCLUDE)
}

// Java: ShowByPathWarningsGuard(String, ShowType)
func NewShowByPathWarningsGuardWithShowType(path string, showType ShowByPathWarningsGuard_ShowType) *ShowByPathWarningsGuard {
	return NewShowByPathWarningsGuardWithPathsAndShowType([]string{path}, showType)
}

// Java: ShowByPathWarningsGuard(String[], ShowType)
func NewShowByPathWarningsGuardWithPathsAndShowType(paths []string, showType ShowByPathWarningsGuard_ShowType) *ShowByPathWarningsGuard {
	jutil.CheckArgument(paths != nil)
	g := &ShowByPathWarningsGuard{}
	pathList := append([]string{}, paths...)
	if showType == ShowByPathWarningsGuard_ShowType_INCLUDE {
		g.warningsGuard = ByPathWarningsGuard_ExceptPath(pathList, CheckLevel_OFF)
	} else {
		g.warningsGuard = ByPathWarningsGuard_ForPath(pathList, CheckLevel_OFF)
	}
	return g
}

// Java: ShowByPathWarningsGuard.level(JSError)
func (g *ShowByPathWarningsGuard) Level(error *JSError) *CheckLevel {
	return g.warningsGuard.Level(error)
}

// Java: ShowByPathWarningsGuard.getPriority()
func (g *ShowByPathWarningsGuard) GetPriority() int {
	return g.warningsGuard.GetPriority()
}

// String is Object.toString().
func (g *ShowByPathWarningsGuard) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.ShowByPathWarningsGuard", g)
}

// J2clSuppressWarningsGuard is a warnings guard that suppresses some
// warnings incompatible with J2CL.
type J2clSuppressWarningsGuard struct {
	WarningsGuardBase
}

var _ WarningsGuard = (*J2clSuppressWarningsGuard)(nil)

// TODO(b/128554878): Cleanup and document all file level suppressions for J2CL generated code.
var J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS *DiagnosticGroup

// Go-only: language-adapter — the static field reads DiagnosticGroups'
// static fields, which are assigned from DiagnosticGroups' init(); this
// file's init() runs after it (Go runs init functions in file name order).
func init() {
	j2clSuppressWarningsGuard_clinit()
}

// Java: J2clSuppressWarningsGuard static initializers
func j2clSuppressWarningsGuard_clinit() {
	J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS =
		NewDiagnosticGroupWithNameFromGroups(
			warningsGuardGroupName("j2clIncomaptible"),
			// Do not warn when static overrides do not match. This is safe, because J2CL generated
			// code directly points to declaration and this is also required with collapse properties
			// pass, which does not support dynamic dispatch for static methods.
			DiagnosticGroups_CHECK_STATIC_OVERRIDES,
			// Do not warn on valid Java constructs like "if(false) {...}" and other situations that
			// may arise from transormation of complex Kotlin constructs.
			DiagnosticGroups_CHECK_USELESS_CODE,
			DiagnosticGroups_CONST,
			DiagnosticGroups_EXTRA_REQUIRE,
			// Kotlin allows provably invalid casts so long as it's via a safe cast. This causes J2CL
			// to generate code of the form:
			//   Foo.$isInstance(value) ? /**@type {!Foo}*/ (value) : null
			// However, if value is obviously never a instance of Foo then this will cause invalid
			// cast error. At runtime it's guarded so it would be safe regardless.
			// This also suppresses casts from non-lambda JsFunction implementations to functions.
			// This particular feature is deprecated and will be removed (b/159954752).
			DiagnosticGroups_INVALID_CASTS,
			DiagnosticGroups_LATE_PROVIDE,
			DiagnosticGroups_MISSING_OVERRIDE,
			DiagnosticGroups_MISSING_REQUIRE,
			DiagnosticGroups_STRICT_MODULE_DEP_CHECK,
			DiagnosticGroups_SUSPICIOUS_CODE,
			DiagnosticGroups_UNUSED_LOCAL_VARIABLE,
			// TODO(b/78521031): J2CL targets are not strict missing property compatible.
			DiagnosticGroups_STRICT_MISSING_PROPERTIES,
			DiagnosticGroups_ForName("transitionalSuspiciousCodeWarnings"))
}

// NewJ2clSuppressWarningsGuard is the implicit default constructor.
//
// Java: J2clSuppressWarningsGuard()
func NewJ2clSuppressWarningsGuard() *J2clSuppressWarningsGuard {
	return &J2clSuppressWarningsGuard{}
}

// Java: J2clSuppressWarningsGuard.level(JSError)
func (g *J2clSuppressWarningsGuard) Level(error *JSError) *CheckLevel {
	if !J2clSourceUtils_IsJ2clSource(error.SourceName()) {
		return nil
	}

	// Do not warn about deprecated forwardDeclare usages in J2CL generated code as this particular
	// use-case is still valid, but not elsewhere.
	if DiagnosticGroups_DEPRECATED.Matches(error) &&
		strings.Contains(error.Description(), "forwardDeclare") {
		return CheckLevel_OFF.Ptr()
	}

	if J2clSuppressWarningsGuard_DEFAULT_J2CL_SUPRRESIONS.Matches(error) {
		return CheckLevel_OFF.Ptr()
	}
	return nil
}

// Java: J2clSuppressWarningsGuard.getPriority()
func (g *J2clSuppressWarningsGuard) GetPriority() int {
	return WarningsGuard_Priority_MAX.GetValue()
}

// DiagnosticGroupWarningsGuard sets the level for a particular
// DiagnosticGroup.
type DiagnosticGroupWarningsGuard struct {
	WarningsGuardBase
	group *DiagnosticGroup
	level CheckLevel
}

var _ WarningsGuard = (*DiagnosticGroupWarningsGuard)(nil)

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
		for _, t := range otherGroup.GetTypes() {
			if g.group.MatchesType(t) {
				return base.Tri_TRUE
			}
		}
		return base.Tri_UNKNOWN
	}
	for _, t := range otherGroup.GetTypes() {
		if !g.group.MatchesType(t) {
			return base.Tri_UNKNOWN
		}
	}
	return base.Tri_FALSE
}

// Java: DiagnosticGroupWarningsGuard.toString()
func (g *DiagnosticGroupWarningsGuard) String() string {
	return g.group.String() + "(" + g.level.String() + ")"
}

// DiagnosticGroupPathSuppressingWarningsGuard is a warnings guard that
// suppresses warnings for a particular diagnostic group for a file that
// contains the specified substring.
type DiagnosticGroupPathSuppressingWarningsGuard struct {
	WarningsGuardBase
	group *DiagnosticGroup
	part  string
}

var _ WarningsGuard = (*DiagnosticGroupPathSuppressingWarningsGuard)(nil)

// Java: DiagnosticGroupPathSuppressingWarningsGuard(DiagnosticGroup, String)
func NewDiagnosticGroupPathSuppressingWarningsGuard(group *DiagnosticGroup, part string) *DiagnosticGroupPathSuppressingWarningsGuard {
	return &DiagnosticGroupPathSuppressingWarningsGuard{group: group, part: part}
}

// Level does not touch warnings in other paths. It may return nil.
//
// Java: DiagnosticGroupPathSuppressingWarningsGuard.level(JSError)
func (g *DiagnosticGroupPathSuppressingWarningsGuard) Level(error *JSError) *CheckLevel {
	if error.SourceName() == nil || !strings.Contains(*error.SourceName(), g.part) {
		return nil
	}
	if g.group.Matches(error) {
		return CheckLevel_OFF.Ptr()
	}
	return nil
}

// Java: DiagnosticGroupPathSuppressingWarningsGuard.toString()
func (g *DiagnosticGroupPathSuppressingWarningsGuard) String() string {
	return g.group.String() + "(" + g.part + ")"
}
