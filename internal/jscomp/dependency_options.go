// Port of com/google/javascript/jscomp/DependencyOptions.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DependencyOptions holds options for how to manage dependencies between
// input files.
//
// Dependency information is pulled out from the JS code by looking for
// import and export primitives (like ES import and export statements,
// Closure Library's goog.module, goog.provide and goog.require calls, or
// CommonJS require calls). The compiler can optionally use this information
// to sort input files in dependency order and/or prune unnecessary input
// files.
//
// Also see CodingConvention#extractClassNameIfProvide(Node, Node) and
// CodingConvention#extractClassNameIfRequire(Node, Node), which affect what
// the compiler considers to be goog.provide and goog.require statements.
//
// mode: the dependency management mode. entryPoints: the list of explicit
// entry points.
type DependencyOptions struct {
	mode        DependencyOptions_DependencyMode
	entryPoints []*ModuleIdentifier
}

// Java: DependencyOptions(DependencyMode, ImmutableList<ModuleIdentifier>)
func NewDependencyOptions(mode DependencyOptions_DependencyMode, entryPoints []*ModuleIdentifier) *DependencyOptions {
	if entryPoints == nil {
		panic(jutil.NewNullPointerExceptionMsg("entryPoints"))
	}
	return &DependencyOptions{mode: mode, entryPoints: entryPoints}
}

// Java: DependencyOptions.mode()
func (o *DependencyOptions) Mode() DependencyOptions_DependencyMode { return o.mode }

// Java: DependencyOptions.entryPoints()
func (o *DependencyOptions) EntryPoints() []*ModuleIdentifier { return o.entryPoints }

// Java: DependencyOptions.getMode()
func (o *DependencyOptions) GetMode() DependencyOptions_DependencyMode { return o.Mode() }

// Java: DependencyOptions.getEntryPoints()
func (o *DependencyOptions) GetEntryPoints() []*ModuleIdentifier { return o.EntryPoints() }

// Java: DependencyOptions.equals(Object) (record)
func (o *DependencyOptions) Equals(obj any) bool {
	that, ok := obj.(*DependencyOptions)
	return ok && that != nil && o.mode == that.mode &&
		slices.EqualFunc(o.entryPoints, that.entryPoints, func(a, b *ModuleIdentifier) bool { return a.Equals(b) })
}

// Java: DependencyOptions.hashCode() (record)
//
// Go-only: language-adapter — Java's Enum.hashCode() is the identity hash;
// the ordinal stands in for it.
func (o *DependencyOptions) HashCode() int32 {
	eh := int32(1)
	for _, ep := range o.entryPoints {
		eh = 31*eh + ep.HashCode()
	}
	return 31*int32(o.mode) + eh
}

// Java: DependencyOptions.toString() (record)
func (o *DependencyOptions) String() string {
	return "DependencyOptions[mode=" + o.mode.String() + ", entryPoints=" + jutil.CollectionToString(o.entryPoints) + "]"
}

// DependencyOptions_DependencyMode describes how the compiler should manage
// dependencies.
type DependencyOptions_DependencyMode int32

const (
	// All input files will be included in the compilation in the order they
	// were specified in.
	DependencyOptions_DependencyMode_NONE DependencyOptions_DependencyMode = iota

	// All input files will be included in the compilation in dependency
	// order.
	DependencyOptions_DependencyMode_SORT_ONLY

	// Input files that are transitive dependencies of the entry points will
	// be included in the compilation in dependency order. All other input
	// files will be dropped.
	//
	// In addition to the explicitly defined entry points, moochers (see
	// below) are implicit entry points.
	DependencyOptions_DependencyMode_PRUNE_LEGACY

	// Input files that are transitive dependencies of the entry points will
	// be included in the compilation in dependency order. All other input
	// files will be dropped.
	//
	// All entry points must be explicitly defined.
	DependencyOptions_DependencyMode_PRUNE

	// Identical to PRUNE unless there are no entry points specifified, in
	// which case drops all other input files instead of crashing with an
	// error about invalid flag usage.
	//
	// PRUNE mode will report an error if used without a specified entry
	// point because that usually indicates user error, where the user
	// actually wants pruning starting from one/more entry points.
	//
	// This mode is only recommended for use when using some macro to
	// generate JSCompiler flags for a variety of targets, some of which have
	// entry points and some of which do not.
	DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS
)

var dependencyOptions_DependencyMode_names = [...]string{
	"NONE", "SORT_ONLY", "PRUNE_LEGACY", "PRUNE", "PRUNE_ALLOW_NO_ENTRY_POINTS",
}

// DependencyOptions_DependencyMode_values is DependencyMode.values().
func DependencyOptions_DependencyMode_values() []DependencyOptions_DependencyMode {
	return []DependencyOptions_DependencyMode{
		DependencyOptions_DependencyMode_NONE, DependencyOptions_DependencyMode_SORT_ONLY,
		DependencyOptions_DependencyMode_PRUNE_LEGACY, DependencyOptions_DependencyMode_PRUNE,
		DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS,
	}
}

// String is Enum.name().
func (m DependencyOptions_DependencyMode) String() string {
	return dependencyOptions_DependencyMode_names[m]
}

// DependencyOptions_DependencyMode_ValueOf is DependencyMode.valueOf(String).
func DependencyOptions_DependencyMode_ValueOf(name string) DependencyOptions_DependencyMode {
	for i, n := range dependencyOptions_DependencyMode_names {
		if n == name {
			return DependencyOptions_DependencyMode(i)
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(
		"No enum constant com.google.javascript.jscomp.DependencyOptions.DependencyMode." + name))
}

// NeedsManagement returns whether dependency management is enabled.
//
// Java: DependencyOptions.needsManagement()
func (o *DependencyOptions) NeedsManagement() bool {
	return o.Mode() != DependencyOptions_DependencyMode_NONE
}

// ShouldSort returns whether files should be sorted.
//
// If true, the input files should be sorted in dependency order. Otherwise,
// input files should not be reordered.
//
// Java: DependencyOptions.shouldSort()
func (o *DependencyOptions) ShouldSort() bool {
	return o.Mode() != DependencyOptions_DependencyMode_NONE
}

// ShouldPrune returns whether files should be pruned.
//
// If true, an input file should be excluded from the compilation if it is
// not a transitive dependency of an entry point. Otherwise, all input files
// should be included.
//
// Java: DependencyOptions.shouldPrune()
func (o *DependencyOptions) ShouldPrune() bool {
	return o.Mode() == DependencyOptions_DependencyMode_PRUNE_LEGACY ||
		o.Mode() == DependencyOptions_DependencyMode_PRUNE ||
		o.Mode() == DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS
}

// ShouldDropMoochers returns whether moochers should be dropped.
//
// A moocher is a strong file that does not goog.provide a namespace and is
// not a goog.module, ES module or CommonJS module. Weak files are never
// considered moochers.
//
// If true, moochers should not be considered implicit entry points.
//
// Java: DependencyOptions.shouldDropMoochers()
func (o *DependencyOptions) ShouldDropMoochers() bool {
	return o.Mode() == DependencyOptions_DependencyMode_PRUNE ||
		o.Mode() == DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS
}

// DependencyOptions_None returns a DependencyOptions using the NONE mode.
//
// Java: DependencyOptions.none()
func DependencyOptions_None() *DependencyOptions {
	return NewDependencyOptions(DependencyOptions_DependencyMode_NONE, []*ModuleIdentifier{})
}

// DependencyOptions_SortOnly returns a DependencyOptions using the SORT_ONLY
// mode.
//
// Java: DependencyOptions.sortOnly()
func DependencyOptions_SortOnly() *DependencyOptions {
	return NewDependencyOptions(DependencyOptions_DependencyMode_SORT_ONLY, []*ModuleIdentifier{})
}

// DependencyOptions_PruneLegacyForEntryPoints returns a DependencyOptions
// using the PRUNE_LEGACY mode with the given entry points.
//
// Deprecated: Prefer DependencyOptions_PruneForEntryPoints with a complete
// list of entry points.
//
// Java: DependencyOptions.pruneLegacyForEntryPoints(Iterable<ModuleIdentifier>)
func DependencyOptions_PruneLegacyForEntryPoints(entryPoints []*ModuleIdentifier) *DependencyOptions {
	return NewDependencyOptions(DependencyOptions_DependencyMode_PRUNE_LEGACY, dependencyOptions_copyOf(entryPoints))
}

// DependencyOptions_PruneAllowNoEntryPoints returns a DependencyOptions
// using the PRUNE_ALLOW_NO_ENTRY_POINTS mode with the given entry points.
//
// Java: DependencyOptions.pruneAllowNoEntryPoints(Iterable<ModuleIdentifier>)
func DependencyOptions_PruneAllowNoEntryPoints(entryPoints []*ModuleIdentifier) *DependencyOptions {
	return NewDependencyOptions(
		DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS, dependencyOptions_copyOf(entryPoints))
}

// DependencyOptions_PruneForEntryPoints returns a DependencyOptions using
// the PRUNE mode with the given entry points.
//
// Java: DependencyOptions.pruneForEntryPoints(Iterable<ModuleIdentifier>)
func DependencyOptions_PruneForEntryPoints(entryPoints []*ModuleIdentifier) *DependencyOptions {
	jutil.CheckStateMsg(len(entryPoints) != 0, "DependencyMode.PRUNE requires at least one entry point")
	return NewDependencyOptions(DependencyOptions_DependencyMode_PRUNE, dependencyOptions_copyOf(entryPoints))
}

// dependencyOptions_copyOf is ImmutableList.copyOf (null elements rejected).
func dependencyOptions_copyOf(entryPoints []*ModuleIdentifier) []*ModuleIdentifier {
	for _, e := range entryPoints {
		jutil.CheckNotNull(e)
	}
	return append([]*ModuleIdentifier{}, entryPoints...)
}

// DependencyOptions_FromFlags is a helper function for validating dependency
// management flags and converting them into a DependencyOptions object.
//
// Returns null when no dependency management flags have been specified.
//
// TODO(tjgq): Simplify this once we deprecate and remove all legacy flags
// and standardize on --dependency_mode and --entry_point.
//
// Java: DependencyOptions.fromFlags(DependencyMode, List<String>, List<String>, String, boolean, boolean)
func DependencyOptions_FromFlags(
	dependencyModeFlag *DependencyOptions_DependencyMode,
	entryPointFlag []string,
	closureEntryPointFlag []string,
	commonJsEntryModuleFlag *string,
	manageClosureDependenciesFlag bool,
	onlyClosureDependenciesFlag bool) *DependencyOptions {

	modeIs := func(m DependencyOptions_DependencyMode) bool {
		return dependencyModeFlag != nil && *dependencyModeFlag == m
	}

	hasEntryPoint := commonJsEntryModuleFlag != nil ||
		len(entryPointFlag) != 0 ||
		len(closureEntryPointFlag) != 0

	if !hasEntryPoint && onlyClosureDependenciesFlag {
		panic(NewFlagUsageException("--only_closure_dependencies requires --entry_point."))
	}

	if !hasEntryPoint && modeIs(DependencyOptions_DependencyMode_PRUNE) {
		panic(NewFlagUsageException("--dependency_mode=PRUNE requires --entry_point."))
	}

	if hasEntryPoint &&
		(modeIs(DependencyOptions_DependencyMode_NONE) ||
			modeIs(DependencyOptions_DependencyMode_SORT_ONLY)) {
		panic(NewFlagUsageException(
			"--dependency_mode=" +
				dependencyModeFlag.String() +
				" cannot be used with --entry_point, --closure_entry_point or " +
				"--common_js_entry_module."))
	}

	if len(entryPointFlag) != 0 && len(closureEntryPointFlag) != 0 {
		panic(NewFlagUsageException("--closure_entry_point cannot be used with --entry_point."))
	}

	if commonJsEntryModuleFlag != nil &&
		(len(entryPointFlag) != 0 || len(closureEntryPointFlag) != 0) {
		panic(NewFlagUsageException(
			"--common_js_entry_module cannot be used with either --entry_point or " +
				"--closure_entry_point."))
	}

	if manageClosureDependenciesFlag && onlyClosureDependenciesFlag {
		panic(NewFlagUsageException(
			"--only_closure_dependencies cannot be used with --manage_closure_dependencies."))
	}

	if manageClosureDependenciesFlag && dependencyModeFlag != nil {
		panic(NewFlagUsageException(
			"--manage_closure_dependencies cannot be used with --dependency_mode."))
	}

	if onlyClosureDependenciesFlag && dependencyModeFlag != nil {
		panic(NewFlagUsageException(
			"--only_closure_dependencies cannot be used with --dependency_mode."))
	}

	var dependencyMode DependencyOptions_DependencyMode
	if modeIs(DependencyOptions_DependencyMode_PRUNE) || onlyClosureDependenciesFlag {
		dependencyMode = DependencyOptions_DependencyMode_PRUNE
	} else if modeIs(DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS) {
		dependencyMode = DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS
	} else if modeIs(DependencyOptions_DependencyMode_PRUNE_LEGACY) ||
		manageClosureDependenciesFlag ||
		hasEntryPoint {
		dependencyMode = DependencyOptions_DependencyMode_PRUNE_LEGACY
	} else if dependencyModeFlag != nil {
		dependencyMode = *dependencyModeFlag
	} else {
		return nil
	}

	var entryPointsBuilder []*ModuleIdentifier
	if commonJsEntryModuleFlag != nil {
		entryPointsBuilder = append(entryPointsBuilder, ModuleIdentifier_ForFile(*commonJsEntryModuleFlag))
	}
	for _, entryPoint := range entryPointFlag {
		entryPointsBuilder = append(entryPointsBuilder, ModuleIdentifier_ForFlagValue(entryPoint))
	}
	for _, closureEntryPoint := range closureEntryPointFlag {
		entryPointsBuilder = append(entryPointsBuilder, ModuleIdentifier_ForClosure(closureEntryPoint))
	}

	switch dependencyMode {
	case DependencyOptions_DependencyMode_NONE:
		return DependencyOptions_None()
	case DependencyOptions_DependencyMode_SORT_ONLY:
		return DependencyOptions_SortOnly()
	case DependencyOptions_DependencyMode_PRUNE_LEGACY:
		return DependencyOptions_PruneLegacyForEntryPoints(entryPointsBuilder)
	case DependencyOptions_DependencyMode_PRUNE_ALLOW_NO_ENTRY_POINTS:
		return DependencyOptions_PruneAllowNoEntryPoints(entryPointsBuilder)
	default: // PRUNE
		return DependencyOptions_PruneForEntryPoints(entryPointsBuilder)
	}
}
