// Port of com/google/javascript/jscomp/CompilerOptions.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

// CompilerOptions_Reach is a common enum for compiler passes that can run either globally or locally.
type CompilerOptions_Reach int32

const (
	CompilerOptions_Reach_ALL CompilerOptions_Reach = iota
	CompilerOptions_Reach_LOCAL_ONLY
	CompilerOptions_Reach_NONE
)

var compilerOptions_Reach_names = [...]string{"ALL", "LOCAL_ONLY", "NONE"}

// CompilerOptions_Reach_values is Reach.values().
func CompilerOptions_Reach_values() []CompilerOptions_Reach {
	return []CompilerOptions_Reach{CompilerOptions_Reach_ALL, CompilerOptions_Reach_LOCAL_ONLY, CompilerOptions_Reach_NONE}
}

// String is Enum.name().
func (e CompilerOptions_Reach) String() string { return compilerOptions_Reach_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_Reach) Ordinal() int { return int(e) }

// Java: CompilerOptions.Reach.isOn()
func (e CompilerOptions_Reach) IsOn() bool {
	return e != CompilerOptions_Reach_NONE
}

// Java: CompilerOptions.Reach.includesGlobals()
func (e CompilerOptions_Reach) IncludesGlobals() bool {
	return e == CompilerOptions_Reach_ALL
}

// CompilerOptions_PropertyCollapseLevel is CompilerOptions.PropertyCollapseLevel.
type CompilerOptions_PropertyCollapseLevel int32

const (
	CompilerOptions_PropertyCollapseLevel_ALL CompilerOptions_PropertyCollapseLevel = iota
	CompilerOptions_PropertyCollapseLevel_NONE
	CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT
)

var compilerOptions_PropertyCollapseLevel_names = [...]string{"ALL", "NONE", "MODULE_EXPORT"}

// CompilerOptions_PropertyCollapseLevel_values is PropertyCollapseLevel.values().
func CompilerOptions_PropertyCollapseLevel_values() []CompilerOptions_PropertyCollapseLevel {
	return []CompilerOptions_PropertyCollapseLevel{CompilerOptions_PropertyCollapseLevel_ALL, CompilerOptions_PropertyCollapseLevel_NONE, CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT}
}

// String is Enum.name().
func (e CompilerOptions_PropertyCollapseLevel) String() string {
	return compilerOptions_PropertyCollapseLevel_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_PropertyCollapseLevel) Ordinal() int { return int(e) }

// CompilerOptions_BrowserFeaturesetYear represents browser feature set year to use for compilation. It tells the
// JSCompiler to output code that works on the releases of major browsers that
// were current as of January 1 of the given year, without including
// transpilation or other workarounds for browsers older than that
type CompilerOptions_BrowserFeaturesetYear int32

const (
	CompilerOptions_BrowserFeaturesetYear_YEAR_2012 CompilerOptions_BrowserFeaturesetYear = iota
	CompilerOptions_BrowserFeaturesetYear_YEAR_2018
	CompilerOptions_BrowserFeaturesetYear_YEAR_2019
	CompilerOptions_BrowserFeaturesetYear_YEAR_2020
	CompilerOptions_BrowserFeaturesetYear_YEAR_2021
	CompilerOptions_BrowserFeaturesetYear_YEAR_2022
	CompilerOptions_BrowserFeaturesetYear_YEAR_2023
	CompilerOptions_BrowserFeaturesetYear_YEAR_2024
	CompilerOptions_BrowserFeaturesetYear_YEAR_2025
	CompilerOptions_BrowserFeaturesetYear_YEAR_2026
)

var compilerOptions_BrowserFeaturesetYear_names = [...]string{"YEAR_2012", "YEAR_2018", "YEAR_2019", "YEAR_2020", "YEAR_2021", "YEAR_2022", "YEAR_2023", "YEAR_2024", "YEAR_2025", "YEAR_2026"}

// CompilerOptions_BrowserFeaturesetYear_values is BrowserFeaturesetYear.values().
func CompilerOptions_BrowserFeaturesetYear_values() []CompilerOptions_BrowserFeaturesetYear {
	return []CompilerOptions_BrowserFeaturesetYear{CompilerOptions_BrowserFeaturesetYear_YEAR_2012, CompilerOptions_BrowserFeaturesetYear_YEAR_2018, CompilerOptions_BrowserFeaturesetYear_YEAR_2019, CompilerOptions_BrowserFeaturesetYear_YEAR_2020, CompilerOptions_BrowserFeaturesetYear_YEAR_2021, CompilerOptions_BrowserFeaturesetYear_YEAR_2022, CompilerOptions_BrowserFeaturesetYear_YEAR_2023, CompilerOptions_BrowserFeaturesetYear_YEAR_2024, CompilerOptions_BrowserFeaturesetYear_YEAR_2025, CompilerOptions_BrowserFeaturesetYear_YEAR_2026}
}

// String is Enum.name().
func (e CompilerOptions_BrowserFeaturesetYear) String() string {
	return compilerOptions_BrowserFeaturesetYear_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_BrowserFeaturesetYear) Ordinal() int { return int(e) }

var compilerOptions_BrowserFeaturesetYear_ptrs = func() []CompilerOptions_BrowserFeaturesetYear {
	v := CompilerOptions_BrowserFeaturesetYear_values()
	return v
}()

// Ptr returns the canonical pointer to e (Go-only: language-adapter — a
// nullable Java BrowserFeaturesetYear reference is a *CompilerOptions_BrowserFeaturesetYear; nil is null).
func (e CompilerOptions_BrowserFeaturesetYear) Ptr() *CompilerOptions_BrowserFeaturesetYear {
	return &compilerOptions_BrowserFeaturesetYear_ptrs[e]
}

// Enum fields, indexed by ordinal (Java: BrowserFeaturesetYear(int, FeatureSet)).
var compilerOptions_BrowserFeaturesetYear_year = [...]int{2012, 2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025, 2026}

var compilerOptions_BrowserFeaturesetYear_featureSet = [...]*parser.FeatureSet{
	CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT.ToFeatureSet(),
	CompilerOptions_LanguageMode_ECMASCRIPT_2016.ToFeatureSet(),
	parser.FeatureSet_BROWSER_2019,
	parser.FeatureSet_BROWSER_2020,
	parser.FeatureSet_BROWSER_2021,
	parser.FeatureSet_BROWSER_2022,
	parser.FeatureSet_BROWSER_2023,
	parser.FeatureSet_BROWSER_2024,
	parser.FeatureSet_BROWSER_2025,
	parser.FeatureSet_BROWSER_2026,
}

// Java: CompilerOptions.BrowserFeaturesetYear.YEAR_MAP
var compilerOptions_BrowserFeaturesetYear_YEAR_MAP = func() *jutil.LinkedHashMap[int, CompilerOptions_BrowserFeaturesetYear] {
	m := jutil.NewLinkedHashMap[int, CompilerOptions_BrowserFeaturesetYear]()
	// go/keep-sorted start
	m.Put(2012, CompilerOptions_BrowserFeaturesetYear_YEAR_2012)
	m.Put(2018, CompilerOptions_BrowserFeaturesetYear_YEAR_2018)
	m.Put(2019, CompilerOptions_BrowserFeaturesetYear_YEAR_2019)
	m.Put(2020, CompilerOptions_BrowserFeaturesetYear_YEAR_2020)
	m.Put(2021, CompilerOptions_BrowserFeaturesetYear_YEAR_2021)
	m.Put(2022, CompilerOptions_BrowserFeaturesetYear_YEAR_2022)
	m.Put(2023, CompilerOptions_BrowserFeaturesetYear_YEAR_2023)
	m.Put(2024, CompilerOptions_BrowserFeaturesetYear_YEAR_2024)
	m.Put(2025, CompilerOptions_BrowserFeaturesetYear_YEAR_2025)
	m.Put(2026, CompilerOptions_BrowserFeaturesetYear_YEAR_2026)
	// go/keep-sorted end
	return m
}()

// Java: CompilerOptions.BrowserFeaturesetYear.from(int)
func CompilerOptions_BrowserFeaturesetYear_From(year int) CompilerOptions_BrowserFeaturesetYear {
	jutil.CheckStateMsg1(
		compilerOptions_BrowserFeaturesetYear_YEAR_MAP.ContainsKey(year),
		"Illegal browser_featureset_year=%s. We support values 2012, or 2018..2026 only",
		year)
	return compilerOptions_BrowserFeaturesetYear_YEAR_MAP.Get(year)
}

// Java: CompilerOptions.BrowserFeaturesetYear.setDependentValuesFromYear(CompilerOptions)
func (e CompilerOptions_BrowserFeaturesetYear) SetDependentValuesFromYear(options *CompilerOptions) {
	year := compilerOptions_BrowserFeaturesetYear_year[e]
	options.SetOutputFeatureSet(compilerOptions_BrowserFeaturesetYear_featureSet[e])
	// All values targeted by browser featureset year are default strict.
	defaultStrict := true
	options.languageOutIsDefaultStrict = &defaultStrict
	options.SetDefineToNumberLiteral("goog.FEATURESET_YEAR", year)
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES5", year > 2012)
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES6", year >= 2018)
	options.SetDefineToBooleanLiteral("$jscomp.ASSUME_ES2020", year >= 2021)
}

// Java: CompilerOptions.BrowserFeaturesetYear.getFeatureSet()
func (e CompilerOptions_BrowserFeaturesetYear) GetFeatureSet() *parser.FeatureSet {
	return compilerOptions_BrowserFeaturesetYear_featureSet[e]
}

// Java: CompilerOptions.BrowserFeaturesetYear.getYear()
func (e CompilerOptions_BrowserFeaturesetYear) GetYear() int {
	return compilerOptions_BrowserFeaturesetYear_year[e]
}

// CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor returns the
// minimum browser featureset year required for the given feature, or nil if
// the feature is in a language level higher than that corresponding to the
// most recent BrowserFeaturesetYear.
//
// Java: CompilerOptions.BrowserFeaturesetYear.minimumRequiredFor(FeatureSet.Feature)
func CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(feature parser.FeatureSet_Feature) *CompilerOptions_BrowserFeaturesetYear {
	// Depends on YEAR_TO_FEATURE_SET being created with keys in increasing order
	for year := range compilerOptions_BrowserFeaturesetYear_YEAR_MAP.Keys() {
		yearObject := compilerOptions_BrowserFeaturesetYear_YEAR_MAP.Get(year)
		if yearObject.GetFeatureSet().ContainsFeature(feature) {
			return yearObject.Ptr()
		}
	}
	return nil
}

// CompilerOptions_IncrementalCheckMode is CompilerOptions.IncrementalCheckMode.
type CompilerOptions_IncrementalCheckMode int32

const (
	// Normal mode
	CompilerOptions_IncrementalCheckMode_OFF CompilerOptions_IncrementalCheckMode = iota
	// The compiler should generate an output file that represents the type-only
	// interface of the code being compiled. This is useful for incremental type
	// checking.
	CompilerOptions_IncrementalCheckMode_GENERATE_IJS
)

var compilerOptions_IncrementalCheckMode_names = [...]string{"OFF", "GENERATE_IJS"}

// CompilerOptions_IncrementalCheckMode_values is IncrementalCheckMode.values().
func CompilerOptions_IncrementalCheckMode_values() []CompilerOptions_IncrementalCheckMode {
	return []CompilerOptions_IncrementalCheckMode{CompilerOptions_IncrementalCheckMode_OFF, CompilerOptions_IncrementalCheckMode_GENERATE_IJS}
}

// String is Enum.name().
func (e CompilerOptions_IncrementalCheckMode) String() string {
	return compilerOptions_IncrementalCheckMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_IncrementalCheckMode) Ordinal() int { return int(e) }

// CompilerOptions_ExtractPrototypeMemberDeclarationsMode is CompilerOptions.ExtractPrototypeMemberDeclarationsMode.
type CompilerOptions_ExtractPrototypeMemberDeclarationsMode int32

const (
	CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF CompilerOptions_ExtractPrototypeMemberDeclarationsMode = iota
	CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP
	CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP
	CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE
)

var compilerOptions_ExtractPrototypeMemberDeclarationsMode_names = [...]string{"OFF", "USE_GLOBAL_TEMP", "USE_CHUNK_TEMP", "USE_IIFE"}

// CompilerOptions_ExtractPrototypeMemberDeclarationsMode_values is ExtractPrototypeMemberDeclarationsMode.values().
func CompilerOptions_ExtractPrototypeMemberDeclarationsMode_values() []CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return []CompilerOptions_ExtractPrototypeMemberDeclarationsMode{CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF, CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP, CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP, CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE}
}

// String is Enum.name().
func (e CompilerOptions_ExtractPrototypeMemberDeclarationsMode) String() string {
	return compilerOptions_ExtractPrototypeMemberDeclarationsMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_ExtractPrototypeMemberDeclarationsMode) Ordinal() int { return int(e) }

// CompilerOptions_OutputJs is what type of JS file should be output by this compilation
type CompilerOptions_OutputJs int32

const (
	// Don't output anything.
	CompilerOptions_OutputJs_NONE CompilerOptions_OutputJs = iota
	// Output a "sentinel" file containing just a comment.
	CompilerOptions_OutputJs_SENTINEL
	// Output the compiled JS.
	CompilerOptions_OutputJs_NORMAL
)

var compilerOptions_OutputJs_names = [...]string{"NONE", "SENTINEL", "NORMAL"}

// CompilerOptions_OutputJs_values is OutputJs.values().
func CompilerOptions_OutputJs_values() []CompilerOptions_OutputJs {
	return []CompilerOptions_OutputJs{CompilerOptions_OutputJs_NONE, CompilerOptions_OutputJs_SENTINEL, CompilerOptions_OutputJs_NORMAL}
}

// String is Enum.name().
func (e CompilerOptions_OutputJs) String() string { return compilerOptions_OutputJs_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_OutputJs) Ordinal() int { return int(e) }

// CompilerOptions_ConformanceReportingMode is whether to respect the library-level non-allowlisted conformance
// behavior specified in the conformance config.
type CompilerOptions_ConformanceReportingMode int32

const (
	// Used when running conformance checks at the binary level. It ignores the
	// library-level conformance behavior specified in the conformance config and
	// always report the conformance violations. For example, setting this allows
	// the compiler to ignore the "RECORD_ONLY" library-level behavior specified in
	// the config and unconditionally report the non-allowlisted violations when it
	// is running at the binary level. (default)
	CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG CompilerOptions_ConformanceReportingMode = iota
	// Used when running conformance checks at the library level. It respects the
	// library-level conformance behavior specified in the conformance config and
	// only reports the conformance violations according to the behavior. For
	// example, setting this respects the "RECORD_ONLY" library-level behavior
	// specified in the Boq Web configs and only records the violations when it is
	// running at the library level, and also respects the "REPORT_AS_BUILD_ERROR"
	// behavior specified in the global conformance configs to report errors for
	// go/jscp at the library level.
	CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG
)

var compilerOptions_ConformanceReportingMode_names = [...]string{"IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG", "RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG"}

// CompilerOptions_ConformanceReportingMode_values is ConformanceReportingMode.values().
func CompilerOptions_ConformanceReportingMode_values() []CompilerOptions_ConformanceReportingMode {
	return []CompilerOptions_ConformanceReportingMode{CompilerOptions_ConformanceReportingMode_IGNORE_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG, CompilerOptions_ConformanceReportingMode_RESPECT_LIBRARY_LEVEL_BEHAVIOR_SPECIFIED_IN_CONFIG}
}

// String is Enum.name().
func (e CompilerOptions_ConformanceReportingMode) String() string {
	return compilerOptions_ConformanceReportingMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_ConformanceReportingMode) Ordinal() int { return int(e) }

// CompilerOptions_ExperimentalOutputFeatureSet is the option to enable
// specific output features for performance experiments.
//
// Deprecated.
type CompilerOptions_ExperimentalOutputFeatureSet int32

const (
	// If the output feature set is ES2015 or above, downgrade to ES5 output
	// with only a small subset of ES2015 features enabled that are known not to
	// have negative performance impact. Additionally allow ASYNC_FUNCTIONS and
	// ASYNC_GENERATORS output for better async stack traces.
	//
	// Deprecated: Use setBrowserFeaturesetYear or setLanguageOut instead.
	CompilerOptions_ExperimentalOutputFeatureSet_ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS CompilerOptions_ExperimentalOutputFeatureSet = iota
	// If the output feature set is ES2015 or above, downgrade to ES5 output
	// with all BROWSER_2019 features enabled that are known not to have
	// negative performance impact. Features currently excluded because of
	// performance concerns are classes (pending confirmation of actual
	// performance impact), and spread expressions (confirmed known performance
	// impact for Firefox 154 during static script evaluation, see
	// https://bugzilla.mozilla.org/show_bug.cgi?id=2069741).
	//
	// Deprecated: Use setBrowserFeaturesetYear or setLanguageOut instead.
	CompilerOptions_ExperimentalOutputFeatureSet_BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD
)

var compilerOptions_ExperimentalOutputFeatureSet_names = [...]string{
	"ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS",
	"BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD",
}

// CompilerOptions_ExperimentalOutputFeatureSet_values is ExperimentalOutputFeatureSet.values().
func CompilerOptions_ExperimentalOutputFeatureSet_values() []CompilerOptions_ExperimentalOutputFeatureSet {
	return []CompilerOptions_ExperimentalOutputFeatureSet{
		CompilerOptions_ExperimentalOutputFeatureSet_ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS,
		CompilerOptions_ExperimentalOutputFeatureSet_BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD,
	}
}

// String is Enum.name().
func (e CompilerOptions_ExperimentalOutputFeatureSet) String() string {
	return compilerOptions_ExperimentalOutputFeatureSet_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_ExperimentalOutputFeatureSet) Ordinal() int { return int(e) }

// CompilerOptions_Es6ModuleTranspilation is how ES modules should be transformed.
type CompilerOptions_Es6ModuleTranspilation int32

const (
	// Do not touch any Es6 module feature.
	CompilerOptions_Es6ModuleTranspilation_NONE CompilerOptions_Es6ModuleTranspilation = iota
	// Rewrite import paths to resolved, relative paths only.
	CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS
	// Rewrite to common js like modules for bundling.
	CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES
	// Compile ES modules.
	CompilerOptions_Es6ModuleTranspilation_COMPILE
)

var compilerOptions_Es6ModuleTranspilation_names = [...]string{"NONE", "RELATIVIZE_IMPORT_PATHS", "TO_COMMON_JS_LIKE_MODULES", "COMPILE"}

// CompilerOptions_Es6ModuleTranspilation_values is Es6ModuleTranspilation.values().
func CompilerOptions_Es6ModuleTranspilation_values() []CompilerOptions_Es6ModuleTranspilation {
	return []CompilerOptions_Es6ModuleTranspilation{CompilerOptions_Es6ModuleTranspilation_NONE, CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS, CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES, CompilerOptions_Es6ModuleTranspilation_COMPILE}
}

// String is Enum.name().
func (e CompilerOptions_Es6ModuleTranspilation) String() string {
	return compilerOptions_Es6ModuleTranspilation_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_Es6ModuleTranspilation) Ordinal() int { return int(e) }

// CompilerOptions_Es6SubclassTranspilation configures the transpiled output for ES6 classes when targeting ES5.
//
// This is needed for transpilation of ES6 classes that extend some unknown
// superclass, in case that superclass is still a native ES6 class at runtime
// (not transpiled to ES5).
//
// The default is CONCISE_UNSAFE, which assumes that if a superclass in an
// extends clause cannot be statically analyzed, it's safe to assume it is an
// ES5 class.
//
// The "safe" behavior relies on $jscomp.construct, which is our wrapper for
// Reflect.construct. This is correct, but has the disadvantage of pulling in
// more helper utilities into the compiled output and of the transpiled output
// being slightly larger.
//
// TODO: b/36789413 - always enable the "safe" transpilation.
type CompilerOptions_Es6SubclassTranspilation int32

const (
	CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE CompilerOptions_Es6SubclassTranspilation = iota
	CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT
)

var compilerOptions_Es6SubclassTranspilation_names = [...]string{"CONCISE_UNSAFE", "SAFE_REFLECT_CONSTRUCT"}

// CompilerOptions_Es6SubclassTranspilation_values is Es6SubclassTranspilation.values().
func CompilerOptions_Es6SubclassTranspilation_values() []CompilerOptions_Es6SubclassTranspilation {
	return []CompilerOptions_Es6SubclassTranspilation{CompilerOptions_Es6SubclassTranspilation_CONCISE_UNSAFE, CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT}
}

// String is Enum.name().
func (e CompilerOptions_Es6SubclassTranspilation) String() string {
	return compilerOptions_Es6SubclassTranspilation_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_Es6SubclassTranspilation) Ordinal() int { return int(e) }

// CompilerOptions_InstrumentOption is an option to determine what level of code instrumentation is performed, if any
type CompilerOptions_InstrumentOption int32

const (
	// No coverage instrumentation is performed
	CompilerOptions_InstrumentOption_NONE CompilerOptions_InstrumentOption = iota
	// Collect coverage for every executable statement.
	CompilerOptions_InstrumentOption_LINE_ONLY
	// Collect coverage for control-flow branches.
	CompilerOptions_InstrumentOption_BRANCH_ONLY
	// Collect coverage for functions where code is compiled for production.
	CompilerOptions_InstrumentOption_PRODUCTION
)

var compilerOptions_InstrumentOption_names = [...]string{"NONE", "LINE_ONLY", "BRANCH_ONLY", "PRODUCTION"}

// CompilerOptions_InstrumentOption_values is InstrumentOption.values().
func CompilerOptions_InstrumentOption_values() []CompilerOptions_InstrumentOption {
	return []CompilerOptions_InstrumentOption{CompilerOptions_InstrumentOption_NONE, CompilerOptions_InstrumentOption_LINE_ONLY, CompilerOptions_InstrumentOption_BRANCH_ONLY, CompilerOptions_InstrumentOption_PRODUCTION}
}

// String is Enum.name().
func (e CompilerOptions_InstrumentOption) String() string {
	return compilerOptions_InstrumentOption_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_InstrumentOption) Ordinal() int { return int(e) }

var compilerOptions_InstrumentOption_ptrs = func() []CompilerOptions_InstrumentOption {
	v := CompilerOptions_InstrumentOption_values()
	return v
}()

// Ptr returns the canonical pointer to e (Go-only: language-adapter — a
// nullable Java InstrumentOption reference is a *CompilerOptions_InstrumentOption; nil is null).
func (e CompilerOptions_InstrumentOption) Ptr() *CompilerOptions_InstrumentOption {
	return &compilerOptions_InstrumentOption_ptrs[e]
}

// Java: CompilerOptions.InstrumentOption.fromString(String)
func CompilerOptions_InstrumentOption_FromString(value *string) *CompilerOptions_InstrumentOption {
	if value == nil {
		return nil
	}
	switch *value {
	case "NONE":
		return CompilerOptions_InstrumentOption_NONE.Ptr()
	case "LINE":
		return CompilerOptions_InstrumentOption_LINE_ONLY.Ptr()
	case "BRANCH":
		return CompilerOptions_InstrumentOption_BRANCH_ONLY.Ptr()
	case "PRODUCTION":
		return CompilerOptions_InstrumentOption_PRODUCTION.Ptr()
	default:
		return nil
	}
}

// CompilerOptions_ChunkOutputType is the format for compiler output when multiple chunks are requested.
type CompilerOptions_ChunkOutputType int32

const (
	// Original format - each chunk is a simple script
	CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE CompilerOptions_ChunkOutputType = iota
	// Each chunk is an ES module with dependencies and symbol access set up using
	// import and export statements.
	CompilerOptions_ChunkOutputType_ES_MODULES
)

var compilerOptions_ChunkOutputType_names = [...]string{"GLOBAL_NAMESPACE", "ES_MODULES"}

// CompilerOptions_ChunkOutputType_values is ChunkOutputType.values().
func CompilerOptions_ChunkOutputType_values() []CompilerOptions_ChunkOutputType {
	return []CompilerOptions_ChunkOutputType{CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE, CompilerOptions_ChunkOutputType_ES_MODULES}
}

// String is Enum.name().
func (e CompilerOptions_ChunkOutputType) String() string {
	return compilerOptions_ChunkOutputType_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_ChunkOutputType) Ordinal() int { return int(e) }

// CompilerOptions_OptimizeLocalAccess is the option for optimizing local
// access to the global symbol namespace when using renamePrefixNamespace.
type CompilerOptions_OptimizeLocalAccess int32

const (
	// No local access optimization. All global symbols are rescoped to the
	// global namespace.
	CompilerOptions_OptimizeLocalAccess_DISABLED CompilerOptions_OptimizeLocalAccess = iota

	// Use a local closure scope alias for global symbols within the chunk where
	// they are defined, if they are not reassigned in other chunks.
	CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY

	// In addition to DEFINING_CHUNK_ONLY, create local aliases (e.g. `var a =
	// _.a;`) in other chunks that use the global symbol, and use the local
	// alias instead of the global namespace. This is only done for symbols that
	// are not reassigned anywhere.
	CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS

	// In addition to ALL_CHUNKS, create dedicated singleton holder objects for
	// mutable cross-chunk symbols (e.g. `var a = _.a = {};`) and access them
	// via `a._`.
	CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS
)

var compilerOptions_OptimizeLocalAccess_names = [...]string{
	"DISABLED", "DEFINING_CHUNK_ONLY", "ALL_CHUNKS", "ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS",
}

// CompilerOptions_OptimizeLocalAccess_values is OptimizeLocalAccess.values().
func CompilerOptions_OptimizeLocalAccess_values() []CompilerOptions_OptimizeLocalAccess {
	return []CompilerOptions_OptimizeLocalAccess{
		CompilerOptions_OptimizeLocalAccess_DISABLED,
		CompilerOptions_OptimizeLocalAccess_DEFINING_CHUNK_ONLY,
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS,
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS,
	}
}

// String is Enum.name().
func (e CompilerOptions_OptimizeLocalAccess) String() string {
	return compilerOptions_OptimizeLocalAccess_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_OptimizeLocalAccess) Ordinal() int { return int(e) }

// Java: CompilerOptions.OptimizeLocalAccess.shouldOptimize()
func (e CompilerOptions_OptimizeLocalAccess) ShouldOptimize() bool {
	return e != CompilerOptions_OptimizeLocalAccess_DISABLED
}

// Java: CompilerOptions.OptimizeLocalAccess.shouldOptimizeAllChunks()
func (e CompilerOptions_OptimizeLocalAccess) ShouldOptimizeAllChunks() bool {
	return e.Ordinal() >= CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS.Ordinal()
}

// Java: CompilerOptions.OptimizeLocalAccess.shouldWrapReassignableSymbols()
func (e CompilerOptions_OptimizeLocalAccess) ShouldWrapReassignableSymbols() bool {
	return e.Ordinal() >= CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS.Ordinal()
}

// CompilerOptions_LanguageMode is a language mode applies to the whole compilation job. As a result, the
// compiler does not support mixed strict and non-strict in the same compilation
// job. Therefore, the 'use strict' directive is ignored when the language mode
// is not strict.
type CompilerOptions_LanguageMode int32

const (
	// 90's JavaScript
	CompilerOptions_LanguageMode_ECMASCRIPT3 CompilerOptions_LanguageMode = iota
	// Traditional JavaScript
	CompilerOptions_LanguageMode_ECMASCRIPT5
	// Nitpicky, traditional JavaScript
	CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT
	// ECMAScript standard approved in 2015.
	CompilerOptions_LanguageMode_ECMASCRIPT_2015
	// ECMAScript standard approved in 2016. Adds the exponent operator (**).
	CompilerOptions_LanguageMode_ECMASCRIPT_2016
	// ECMAScript standard approved in 2017. Adds async/await and other syntax
	CompilerOptions_LanguageMode_ECMASCRIPT_2017
	// ECMAScript standard approved in 2018. Adds "..." in object literals/patterns.
	CompilerOptions_LanguageMode_ECMASCRIPT_2018
	// ECMAScript standard approved in 2019. Adds catch blocks with no error binding.
	CompilerOptions_LanguageMode_ECMASCRIPT_2019
	// ECMAScript standard approved in 2020.
	CompilerOptions_LanguageMode_ECMASCRIPT_2020
	// ECMAScript standard approved in 2021.
	CompilerOptions_LanguageMode_ECMASCRIPT_2021
	// ECMAScript standard approved in 2022.
	CompilerOptions_LanguageMode_ECMASCRIPT_2022
	// ECMAScript features from the upcoming standard.
	CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	// Use stable features.
	CompilerOptions_LanguageMode_STABLE
	// For languageOut only. The same language mode as the input.
	CompilerOptions_LanguageMode_NO_TRANSPILE
	// For testing only. Features that can be parsed/checked but cannot necessarily
	// be understood by the rest of the compiler yet.
	CompilerOptions_LanguageMode_UNSTABLE
	// For testing only. Features that can be parsed but cannot necessarily be
	// understood by the rest of the compiler yet.
	CompilerOptions_LanguageMode_UNSUPPORTED
)

var compilerOptions_LanguageMode_names = [...]string{"ECMASCRIPT3", "ECMASCRIPT5", "ECMASCRIPT5_STRICT", "ECMASCRIPT_2015", "ECMASCRIPT_2016", "ECMASCRIPT_2017", "ECMASCRIPT_2018", "ECMASCRIPT_2019", "ECMASCRIPT_2020", "ECMASCRIPT_2021", "ECMASCRIPT_2022", "ECMASCRIPT_NEXT", "STABLE", "NO_TRANSPILE", "UNSTABLE", "UNSUPPORTED"}

// CompilerOptions_LanguageMode_values is LanguageMode.values().
func CompilerOptions_LanguageMode_values() []CompilerOptions_LanguageMode {
	return []CompilerOptions_LanguageMode{CompilerOptions_LanguageMode_ECMASCRIPT3, CompilerOptions_LanguageMode_ECMASCRIPT5, CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT, CompilerOptions_LanguageMode_ECMASCRIPT_2015, CompilerOptions_LanguageMode_ECMASCRIPT_2016, CompilerOptions_LanguageMode_ECMASCRIPT_2017, CompilerOptions_LanguageMode_ECMASCRIPT_2018, CompilerOptions_LanguageMode_ECMASCRIPT_2019, CompilerOptions_LanguageMode_ECMASCRIPT_2020, CompilerOptions_LanguageMode_ECMASCRIPT_2021, CompilerOptions_LanguageMode_ECMASCRIPT_2022, CompilerOptions_LanguageMode_ECMASCRIPT_NEXT, CompilerOptions_LanguageMode_STABLE, CompilerOptions_LanguageMode_NO_TRANSPILE, CompilerOptions_LanguageMode_UNSTABLE, CompilerOptions_LanguageMode_UNSUPPORTED}
}

// String is Enum.name().
func (e CompilerOptions_LanguageMode) String() string { return compilerOptions_LanguageMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_LanguageMode) Ordinal() int { return int(e) }

var compilerOptions_LanguageMode_ptrs = func() []CompilerOptions_LanguageMode {
	v := CompilerOptions_LanguageMode_values()
	return v
}()

// Ptr returns the canonical pointer to e (Go-only: language-adapter — a
// nullable Java LanguageMode reference is a *CompilerOptions_LanguageMode; nil is null).
func (e CompilerOptions_LanguageMode) Ptr() *CompilerOptions_LanguageMode {
	return &compilerOptions_LanguageMode_ptrs[e]
}

// CompilerOptions_LanguageMode_STABLE_IN is the default input language
// level of JSCompiler.
//
// Java: CompilerOptions.LanguageMode.STABLE_IN
const CompilerOptions_LanguageMode_STABLE_IN = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT

// CompilerOptions_LanguageMode_STABLE_OUT is the default output language
// level of JSCompiler.
//
// Java: CompilerOptions.LanguageMode.STABLE_OUT
const CompilerOptions_LanguageMode_STABLE_OUT = CompilerOptions_LanguageMode_ECMASCRIPT5

// IsDefaultStrict reports whether this language mode defaults to strict
// mode
//
// Java: CompilerOptions.LanguageMode.isDefaultStrict()
func (e CompilerOptions_LanguageMode) IsDefaultStrict() bool {
	switch e {
	case CompilerOptions_LanguageMode_ECMASCRIPT3, CompilerOptions_LanguageMode_ECMASCRIPT5:
		return false
	default:
		return true
	}
}

// CompilerOptions_LanguageMode_ValidCommandLineNames returns a list of valid
// names used to select a `LanguageMode` on the command line.
//
// Java: CompilerOptions.LanguageMode.validCommandLineNames()
func CompilerOptions_LanguageMode_ValidCommandLineNames() []string {
	names := []string{}
	for _, mode := range CompilerOptions_LanguageMode_values() {
		if mode != CompilerOptions_LanguageMode_UNSUPPORTED {
			names = append(names, mode.String())
			if strings.HasPrefix(mode.String(), "ECMASCRIPT") {
				names = append(names, strings.Replace(mode.String(), "ECMASCRIPT", "ES", 1))
			}
		}
	}
	names = append(names, "ECMASCRIPT6")
	names = append(names, "ES6")
	names = append(names, "ECMASCRIPT6_STRICT")
	names = append(names, "ES6_STRICT")
	return names
}

// Java: CompilerOptions.LanguageMode.fromString(String)
func CompilerOptions_LanguageMode_FromString(value *string) *CompilerOptions_LanguageMode {
	if value == nil {
		return nil
	}
	// Trim spaces, disregard case, and allow abbreviation of ECMASCRIPT for convenience.
	canonicalizedName := jutil.StringToUpperCaseASCII(jutil.StringTrim(*value))
	if strings.HasPrefix(canonicalizedName, "ES") {
		canonicalizedName = "ECMASCRIPT" + canonicalizedName[2:]
	}

	if canonicalizedName == "ECMASCRIPT6" ||
		canonicalizedName == "ECMASCRIPT6_STRICT" {
		return CompilerOptions_LanguageMode_ECMASCRIPT_2015.Ptr()
	}

	// LanguageMode.valueOf(canonicalizedName), or null for an unknown name.
	for _, mode := range CompilerOptions_LanguageMode_values() {
		if mode.String() == canonicalizedName {
			return mode.Ptr()
		}
	}
	return nil // unknown name.
}

// Java: CompilerOptions.LanguageMode.toFeatureSet()
func (e CompilerOptions_LanguageMode) ToFeatureSet() *parser.FeatureSet {
	switch e {
	case CompilerOptions_LanguageMode_ECMASCRIPT3:
		return parser.FeatureSet_ES3
	case CompilerOptions_LanguageMode_ECMASCRIPT5, CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT:
		return parser.FeatureSet_ES5
	case CompilerOptions_LanguageMode_ECMASCRIPT_2015:
		return parser.FeatureSet_ES2015_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2016:
		return parser.FeatureSet_ES2016_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2017:
		return parser.FeatureSet_ES2017_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2018:
		return parser.FeatureSet_ES2018_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2019:
		return parser.FeatureSet_ES2019_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2020:
		return parser.FeatureSet_ES2020_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2021:
		return parser.FeatureSet_ES2021_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_2022:
		return parser.FeatureSet_ES2022_MODULES
	case CompilerOptions_LanguageMode_ECMASCRIPT_NEXT:
		return parser.FeatureSet_ES_NEXT
	case CompilerOptions_LanguageMode_NO_TRANSPILE, CompilerOptions_LanguageMode_UNSTABLE:
		return parser.FeatureSet_ES_UNSTABLE
	case CompilerOptions_LanguageMode_UNSUPPORTED:
		return parser.FeatureSet_ES_UNSUPPORTED
	case CompilerOptions_LanguageMode_STABLE:
		panic(jutil.NewUnsupportedOperationExceptionMsg(
			"STABLE has different feature sets for language in and out. " +
				"Use STABLE_IN or STABLE_OUT."))
	}
	panic(jutil.NewIllegalStateException())
}

// CompilerOptions_DevMode is when to do the extra validity checks
type CompilerOptions_DevMode int32

const (
	// Don't do any extra checks.
	CompilerOptions_DevMode_OFF CompilerOptions_DevMode = iota
	// After the initial parse
	CompilerOptions_DevMode_START
	// At the start and at the end of all optimizations.
	CompilerOptions_DevMode_START_AND_END
	// After every pass
	CompilerOptions_DevMode_EVERY_PASS
)

var compilerOptions_DevMode_names = [...]string{"OFF", "START", "START_AND_END", "EVERY_PASS"}

// CompilerOptions_DevMode_values is DevMode.values().
func CompilerOptions_DevMode_values() []CompilerOptions_DevMode {
	return []CompilerOptions_DevMode{CompilerOptions_DevMode_OFF, CompilerOptions_DevMode_START, CompilerOptions_DevMode_START_AND_END, CompilerOptions_DevMode_EVERY_PASS}
}

// String is Enum.name().
func (e CompilerOptions_DevMode) String() string { return compilerOptions_DevMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_DevMode) Ordinal() int { return int(e) }

// CompilerOptions_TracerMode is how much tracing we want to do
type CompilerOptions_TracerMode int32

const (
	// Collect all timing and size metrics. Very slow.
	CompilerOptions_TracerMode_ALL CompilerOptions_TracerMode = iota
	// Collect all timing and size metrics, except gzipped size. Slow.
	CompilerOptions_TracerMode_RAW_SIZE
	// Collect AST size data and metrics on the impact of dependency pruning.
	CompilerOptions_TracerMode_AST_SIZE_AND_PRUNING
	// For size data, don't serialize the AST, just count the number of nodes.
	CompilerOptions_TracerMode_AST_SIZE
	// Collect timing metrics only.
	CompilerOptions_TracerMode_TIMING_ONLY
	// Collect no timing and size metrics.
	CompilerOptions_TracerMode_OFF
)

var compilerOptions_TracerMode_names = [...]string{"ALL", "RAW_SIZE", "AST_SIZE_AND_PRUNING", "AST_SIZE", "TIMING_ONLY", "OFF"}

// CompilerOptions_TracerMode_values is TracerMode.values().
func CompilerOptions_TracerMode_values() []CompilerOptions_TracerMode {
	return []CompilerOptions_TracerMode{CompilerOptions_TracerMode_ALL, CompilerOptions_TracerMode_RAW_SIZE, CompilerOptions_TracerMode_AST_SIZE_AND_PRUNING, CompilerOptions_TracerMode_AST_SIZE, CompilerOptions_TracerMode_TIMING_ONLY, CompilerOptions_TracerMode_OFF}
}

// String is Enum.name().
func (e CompilerOptions_TracerMode) String() string { return compilerOptions_TracerMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_TracerMode) Ordinal() int { return int(e) }

// Java: CompilerOptions.TracerMode.isOn()
func (e CompilerOptions_TracerMode) IsOn() bool {
	return e != CompilerOptions_TracerMode_OFF
}

// Java: CompilerOptions.TracerMode.doPruningAnalysis()
func (e CompilerOptions_TracerMode) DoPruningAnalysis() bool {
	return e == CompilerOptions_TracerMode_AST_SIZE_AND_PRUNING || e == CompilerOptions_TracerMode_ALL || e == CompilerOptions_TracerMode_RAW_SIZE
}

// CompilerOptions_TweakProcessing is the option for the ProcessTweaks pass
type CompilerOptions_TweakProcessing int32

const (
	// Do not run the ProcessTweaks pass.
	CompilerOptions_TweakProcessing_OFF CompilerOptions_TweakProcessing = iota
	// Run the pass, but do not strip out the calls.
	CompilerOptions_TweakProcessing_CHECK
	// Strip out all calls to goog.tweak.*.
	CompilerOptions_TweakProcessing_STRIP
)

var compilerOptions_TweakProcessing_names = [...]string{"OFF", "CHECK", "STRIP"}

// CompilerOptions_TweakProcessing_values is TweakProcessing.values().
func CompilerOptions_TweakProcessing_values() []CompilerOptions_TweakProcessing {
	return []CompilerOptions_TweakProcessing{CompilerOptions_TweakProcessing_OFF, CompilerOptions_TweakProcessing_CHECK, CompilerOptions_TweakProcessing_STRIP}
}

// String is Enum.name().
func (e CompilerOptions_TweakProcessing) String() string {
	return compilerOptions_TweakProcessing_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_TweakProcessing) Ordinal() int { return int(e) }

// Java: CompilerOptions.TweakProcessing.isOn()
func (e CompilerOptions_TweakProcessing) IsOn() bool {
	return e != CompilerOptions_TweakProcessing_OFF
}

// Java: CompilerOptions.TweakProcessing.shouldStrip()
func (e CompilerOptions_TweakProcessing) ShouldStrip() bool {
	return e == CompilerOptions_TweakProcessing_STRIP
}

// CompilerOptions_IsolationMode is what kind of isolation is going to be used
type CompilerOptions_IsolationMode int32

const (
	// output does not include additional isolation.
	CompilerOptions_IsolationMode_NONE CompilerOptions_IsolationMode = iota
	// The output should be wrapped in an IIFE to isolate global variables.
	CompilerOptions_IsolationMode_IIFE
)

var compilerOptions_IsolationMode_names = [...]string{"NONE", "IIFE"}

// CompilerOptions_IsolationMode_values is IsolationMode.values().
func CompilerOptions_IsolationMode_values() []CompilerOptions_IsolationMode {
	return []CompilerOptions_IsolationMode{CompilerOptions_IsolationMode_NONE, CompilerOptions_IsolationMode_IIFE}
}

// String is Enum.name().
func (e CompilerOptions_IsolationMode) String() string { return compilerOptions_IsolationMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_IsolationMode) Ordinal() int { return int(e) }

// CompilerOptions_SegmentOfCompilationToRun is what segment of the compilation to run.
type CompilerOptions_SegmentOfCompilationToRun int32

const (
	// all JSC compilation in a single action.
	CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION CompilerOptions_SegmentOfCompilationToRun = iota
	// all checks
	CompilerOptions_SegmentOfCompilationToRun_CHECKS
	// first half of optimizations
	CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF
	// second half of optimizations
	CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF
	// all optimizations
	CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS
	// all optimizations and finalizations in one action
	CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_AND_FINALIZATIONS
	// all finalizations
	CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS
)

var compilerOptions_SegmentOfCompilationToRun_names = [...]string{"ENTIRE_COMPILATION", "CHECKS", "OPTIMIZATIONS_FIRST_HALF", "OPTIMIZATIONS_SECOND_HALF", "OPTIMIZATIONS", "OPTIMIZATIONS_AND_FINALIZATIONS", "FINALIZATIONS"}

// CompilerOptions_SegmentOfCompilationToRun_values is SegmentOfCompilationToRun.values().
func CompilerOptions_SegmentOfCompilationToRun_values() []CompilerOptions_SegmentOfCompilationToRun {
	return []CompilerOptions_SegmentOfCompilationToRun{CompilerOptions_SegmentOfCompilationToRun_ENTIRE_COMPILATION, CompilerOptions_SegmentOfCompilationToRun_CHECKS, CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF, CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF, CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS, CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_AND_FINALIZATIONS, CompilerOptions_SegmentOfCompilationToRun_FINALIZATIONS}
}

// String is Enum.name().
func (e CompilerOptions_SegmentOfCompilationToRun) String() string {
	return compilerOptions_SegmentOfCompilationToRun_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_SegmentOfCompilationToRun) Ordinal() int { return int(e) }

// CompilerOptions_AliasStringsMode is a mode enum used to indicate the alias strings policy for the AliasStrings pass
type CompilerOptions_AliasStringsMode int32

const (
	// Do not alias string literals.
	CompilerOptions_AliasStringsMode_NONE CompilerOptions_AliasStringsMode = iota
	// Alias all string literals with a length greater than 100 characters.
	CompilerOptions_AliasStringsMode_LARGE
	// Alias all string literals where it may improve code size
	CompilerOptions_AliasStringsMode_ALL
	// Alias all string regardless of code size
	CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE
)

var compilerOptions_AliasStringsMode_names = [...]string{"NONE", "LARGE", "ALL", "ALL_AGGRESSIVE"}

// CompilerOptions_AliasStringsMode_values is AliasStringsMode.values().
func CompilerOptions_AliasStringsMode_values() []CompilerOptions_AliasStringsMode {
	return []CompilerOptions_AliasStringsMode{CompilerOptions_AliasStringsMode_NONE, CompilerOptions_AliasStringsMode_LARGE, CompilerOptions_AliasStringsMode_ALL, CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE}
}

// String is Enum.name().
func (e CompilerOptions_AliasStringsMode) String() string {
	return compilerOptions_AliasStringsMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_AliasStringsMode) Ordinal() int { return int(e) }

// CompilerOptions_Environment is an environment specifies the built-in externs that are loaded for a given compilation.
type CompilerOptions_Environment int32

const (
	// Hand crafted externs that have traditionally been the default externs.
	CompilerOptions_Environment_BROWSER CompilerOptions_Environment = iota
	// Only language externs are loaded.
	CompilerOptions_Environment_CUSTOM
)

var compilerOptions_Environment_names = [...]string{"BROWSER", "CUSTOM"}

// CompilerOptions_Environment_values is Environment.values().
func CompilerOptions_Environment_values() []CompilerOptions_Environment {
	return []CompilerOptions_Environment{CompilerOptions_Environment_BROWSER, CompilerOptions_Environment_CUSTOM}
}

// String is Enum.name().
func (e CompilerOptions_Environment) String() string { return compilerOptions_Environment_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_Environment) Ordinal() int { return int(e) }

// CompilerOptions_JsonStreamMode is whether standard input or standard output should be an array of JSON encoded files
type CompilerOptions_JsonStreamMode int32

const (
	// stdin/out are both single files.
	CompilerOptions_JsonStreamMode_NONE CompilerOptions_JsonStreamMode = iota
	// stdin is a json stream.
	CompilerOptions_JsonStreamMode_IN
	// stdout is a json stream.
	CompilerOptions_JsonStreamMode_OUT
	// stdin and stdout are both json streams.
	CompilerOptions_JsonStreamMode_BOTH
)

var compilerOptions_JsonStreamMode_names = [...]string{"NONE", "IN", "OUT", "BOTH"}

// CompilerOptions_JsonStreamMode_values is JsonStreamMode.values().
func CompilerOptions_JsonStreamMode_values() []CompilerOptions_JsonStreamMode {
	return []CompilerOptions_JsonStreamMode{CompilerOptions_JsonStreamMode_NONE, CompilerOptions_JsonStreamMode_IN, CompilerOptions_JsonStreamMode_OUT, CompilerOptions_JsonStreamMode_BOTH}
}

// String is Enum.name().
func (e CompilerOptions_JsonStreamMode) String() string {
	return compilerOptions_JsonStreamMode_names[e]
}

// Ordinal is Enum.ordinal().
func (e CompilerOptions_JsonStreamMode) Ordinal() int { return int(e) }

// CompilerOptions_J2clPassMode is a mode enum used to indicate whether J2clPass should be enabled,
// disabled, or enabled automatically if there is any J2cl source file (i.e. in
// the AUTO mode).
type CompilerOptions_J2clPassMode int32

const (
	// J2clPass is disabled.
	CompilerOptions_J2clPassMode_OFF CompilerOptions_J2clPassMode = iota
	// It auto-detects whether there are J2cl generated file. If yes, execute J2clPass.
	CompilerOptions_J2clPassMode_AUTO
)

var compilerOptions_J2clPassMode_names = [...]string{"OFF", "AUTO"}

// CompilerOptions_J2clPassMode_values is J2clPassMode.values().
func CompilerOptions_J2clPassMode_values() []CompilerOptions_J2clPassMode {
	return []CompilerOptions_J2clPassMode{CompilerOptions_J2clPassMode_OFF, CompilerOptions_J2clPassMode_AUTO}
}

// String is Enum.name().
func (e CompilerOptions_J2clPassMode) String() string { return compilerOptions_J2clPassMode_names[e] }

// Ordinal is Enum.ordinal().
func (e CompilerOptions_J2clPassMode) Ordinal() int { return int(e) }

// Java: CompilerOptions.J2clPassMode.shouldAddJ2clPasses()
func (e CompilerOptions_J2clPassMode) ShouldAddJ2clPasses() bool {
	return e == CompilerOptions_J2clPassMode_AUTO
}

// CustomPassExecutionTime is the custom pass type. Controls the compilation
// stage at which each custom pass executes.
type CustomPassExecutionTime int32

const (
	CustomPassExecutionTime_BEFORE_CHECKS CustomPassExecutionTime = iota
	CustomPassExecutionTime_BEFORE_OPTIMIZATIONS
	CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP
	CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP
)

var customPassExecutionTime_names = [...]string{
	"BEFORE_CHECKS", "BEFORE_OPTIMIZATIONS", "BEFORE_OPTIMIZATION_LOOP", "AFTER_OPTIMIZATION_LOOP",
}

// Name is Enum.name().
func (t CustomPassExecutionTime) Name() string { return customPassExecutionTime_names[t] }

// String is Enum.toString().
func (t CustomPassExecutionTime) String() string { return t.Name() }

// Ordinal is Enum.ordinal().
func (t CustomPassExecutionTime) Ordinal() int { return int(t) }
