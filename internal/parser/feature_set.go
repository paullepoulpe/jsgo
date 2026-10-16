// Port of com/google/javascript/jscomp/parsing/parser/FeatureSet.java @ closure-compiler v20260415

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// FeatureSet represents various aspects of language version and support.
//
// This is somewhat redundant with LanguageMode, but is separate for two
// reasons: (1) it's used for parsing, which cannot depend on LanguageMode,
// and (2) it's concerned with slightly different nuances: implemented
// features and modules rather than strictness.
//
// In the long term, it would be good to disentangle all these concerns and
// pull out a single LanguageSyntax enum with a separate strict mode flag, and
// then these could possibly be unified.
//
// Instances of this class are immutable.
type FeatureSet struct {
	features *jutil.EnumSet[FeatureSet_Feature]
}

var (
	// The bare minimum set of features.
	FeatureSet_BARE_MINIMUM = NewFeatureSet(FeatureSet_EmptyEnumSet())

	// Features from ES3.
	FeatureSet_ES3 = FeatureSet_BARE_MINIMUM.WithEnumSet(FeatureSet_LangVersion_ES3.Features())

	// Features from ES5 only.
	FeatureSet_ES5 = FeatureSet_ES3.WithEnumSet(FeatureSet_LangVersion_ES5.Features())

	// All ES2015 features, including modules.
	FeatureSet_ES2015_MODULES = FeatureSet_ES5.WithEnumSet(FeatureSet_LangVersion_ES2015.Features())

	// The full set of ES2015 features, not including modules.
	FeatureSet_ES2015 = FeatureSet_ES2015_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)

	FeatureSet_ES2016_MODULES = FeatureSet_ES2015_MODULES.WithEnumSet(FeatureSet_LangVersion_ES2016.Features())

	FeatureSet_ES2016 = FeatureSet_ES2016_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)

	FeatureSet_ES2017_MODULES = FeatureSet_ES2016_MODULES.WithEnumSet(FeatureSet_LangVersion_ES2017.Features())

	FeatureSet_ES2017 = FeatureSet_ES2017_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)

	FeatureSet_ES2018_MODULES = FeatureSet_ES2017_MODULES.WithEnumSet(FeatureSet_LangVersion_ES2018.Features())

	FeatureSet_ES2018 = FeatureSet_ES2018_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)

	FeatureSet_ES2019_MODULES = FeatureSet_ES2018_MODULES.WithEnumSet(FeatureSet_LangVersion_ES2019.Features())

	FeatureSet_ES2019 = FeatureSet_ES2019_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)

	FeatureSet_ES2020_MODULES = FeatureSet_ES2019_MODULES.WithEnumSet(FeatureSet_LangVersion_ES2020.Features())

	FeatureSet_ES2020 = FeatureSet_ES2020_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)

	FeatureSet_ES2021_MODULES = FeatureSet_ES2020_MODULES.WithEnumSet(FeatureSet_LangVersion_ES2021.Features())

	FeatureSet_ES2021 = FeatureSet_ES2021_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)

	FeatureSet_ES2022_MODULES = FeatureSet_ES2021_MODULES.WithEnumSet(FeatureSet_LangVersion_ES2022.Features())

	FeatureSet_ES2022 = FeatureSet_ES2022_MODULES.WithoutFeatures(FeatureSet_Feature_MODULES)

	// NOTE: when ES2023 is added, the BROWSER_2024 and BROWSER_2025 FeatureSets defined below should
	// be updated to include it.

	// Set of all fully supported features, even those part of language versions not fully supported
	FeatureSet_ES_NEXT = FeatureSet_ES2022_MODULES.WithEnumSet(FeatureSet_LangVersion_ES_NEXT.Features())

	// Set of features fully supported in checks, even those not fully supported in optimizations
	FeatureSet_ES_UNSTABLE = FeatureSet_ES_NEXT.WithEnumSet(FeatureSet_LangVersion_ES_UNSTABLE.Features())

	// Set of all features that can be parsed, even those not yet fully supported in checks.
	FeatureSet_ES_UNSUPPORTED = FeatureSet_ES_UNSTABLE.WithEnumSet(FeatureSet_LangVersion_ES_UNSUPPORTED.Features())

	// NOTE: The BROWSER_20XX FeatureSets are the set of features supported by major browsers as of
	// January 1st of the given year. So, typically, take the year of the BFSY, subtract 1 to get
	// the correct ES20XX_MODULES FeatureSet as the base, then exclude any features as needed.
	//
	// Features and their supported browser versions can be found at
	// https://compat-table.github.io/compat-table/es2016plus/ or https://caniuse.com/
	//
	// The format for excluded features is: [Browser name] [version] (version release year)
	// Or, if a browser has no support yet: [Browser name] (unsupported)
	//
	// All excluded features should have a (year) this is greater than or equal to the BFSY.

	FeatureSet_BROWSER_2019 = FeatureSet_ES2018_MODULES.WithoutFeatures(
		// NOTE: These 4 are excluded because, historically, we defined BROWSER_2019 as ES2017.
		// See cl/266708942 and cl/328623467 for context.
		// Chrome 60 (2017). Firefox 55 (2017). Safari 11.1 (2018).
		FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD,
		// Chrome 60 (2017). Firefox 55 (2017). Safari 11.1 (2018).
		FeatureSet_Feature_OBJECT_PATTERN_REST,
		// Chrome 63 (2017). Firefox 57 (2017). Safari 12 (2018).
		FeatureSet_Feature_ASYNC_GENERATORS,
		// Chrome 63 (2017). Firefox 57 (2017). Safari 12 (2018).
		FeatureSet_Feature_FOR_AWAIT_OF,

		// Chrome 62 (2017). Firefox 78 (2020). Safari 11.1 (2018).
		FeatureSet_Feature_REGEXP_FLAG_S,
		// Chrome 62 (2017). Firefox 78 (2020). Safari 16.4 (2023).
		FeatureSet_Feature_REGEXP_LOOKBEHIND,
		// Chrome 64 (2018). Firefox 78 (2020). Safari 11.1 (2018).
		FeatureSet_Feature_REGEXP_NAMED_GROUPS,
		// Chrome 64 (2018). Firefox 78 (2020). Safari 11.1 (2018).
		FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE)

	FeatureSet_BROWSER_2020 = FeatureSet_ES2019_MODULES.WithoutFeatures(
		// Chrome 62 (2017). Firefox 78 (2020). Safari 11.1 (2018).
		FeatureSet_Feature_REGEXP_FLAG_S,
		// Chrome 64 (2018). Firefox 78 (2020). Safari 11.1 (2018).
		FeatureSet_Feature_REGEXP_NAMED_GROUPS,
		// Chrome 64 (2018). Firefox 78 (2020). Safari 11.1 (2018).
		FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE,
		// Chrome 62 (2017). Firefox 78 (2020). Safari 16.4 (2023).
		FeatureSet_Feature_REGEXP_LOOKBEHIND)

	FeatureSet_BROWSER_2021 = FeatureSet_ES2020_MODULES.WithoutFeatures(
		// Chrome 62 (2017). Firefox 78 (2020). Safari 16.4 (2023).
		FeatureSet_Feature_REGEXP_LOOKBEHIND)

	FeatureSet_BROWSER_2022 = FeatureSet_ES2021_MODULES.WithoutFeatures(
		// Chrome 62 (2017). Firefox 78 (2020). Safari 16.4 (2023).
		FeatureSet_Feature_REGEXP_LOOKBEHIND)

	FeatureSet_BROWSER_2023 = FeatureSet_ES2022_MODULES.WithoutFeatures(
		// Chrome 62 (2017). Firefox 78 (2020). Safari 16.4 (2023).
		FeatureSet_Feature_REGEXP_LOOKBEHIND,
		// Chrome 94 (2021). Firefox 93 (2021). Safari 16.4 (2023).
		FeatureSet_Feature_CLASS_STATIC_BLOCK)

	// Note: Only "Hashbang Grammar" is an ES2023 syntax feature. Its versions are:
	// Chrome 74 (2019). Firefox 67 (2019). Safari 13.1 (2020).
	FeatureSet_BROWSER_2024 = FeatureSet_ES2022_MODULES

	// Note: Only "RegExp `v` flag" is an ES2024 syntax feature. Its versions are:
	// Chrome 112 (2023). Firefox 116 (2023). Safari 17 (2023).
	FeatureSet_BROWSER_2025 = FeatureSet_ES2022_MODULES

	// Note: Only "RegExp syntax features again for ES2025. The 'i', 'm', and 's' flags, as well
	// as support for duplicate named capture groups. As of Jan 5 2026, Safari hasn't yet
	// implemented the new flags.
	FeatureSet_BROWSER_2026 = FeatureSet_ES2022_MODULES

	FeatureSet_ALL = FeatureSet_ES_UNSUPPORTED.WithEnumSet(FeatureSet_LangVersion_TYPESCRIPT.Features())
)

// FeatureSet_LangVersion is FeatureSet.LangVersion.
type FeatureSet_LangVersion int32

const (
	FeatureSet_LangVersion_ES3 FeatureSet_LangVersion = iota
	FeatureSet_LangVersion_ES5
	FeatureSet_LangVersion_ES2015
	FeatureSet_LangVersion_ES2016
	FeatureSet_LangVersion_ES2017
	FeatureSet_LangVersion_ES2018
	FeatureSet_LangVersion_ES2019
	FeatureSet_LangVersion_ES2020
	FeatureSet_LangVersion_ES2021
	FeatureSet_LangVersion_ES2022
	FeatureSet_LangVersion_ES_NEXT
	FeatureSet_LangVersion_ES_UNSTABLE
	FeatureSet_LangVersion_ES_UNSUPPORTED
	FeatureSet_LangVersion_TYPESCRIPT
)

var featureSet_LangVersion_names = [...]string{
	"ES3", "ES5", "ES2015", "ES2016", "ES2017", "ES2018", "ES2019", "ES2020", "ES2021", "ES2022",
	"ES_NEXT", "ES_UNSTABLE", "ES_UNSUPPORTED", "TYPESCRIPT",
}

// FeatureSet_LangVersion_values is LangVersion.values().
//
// Java: FeatureSet.LangVersion.values()
func FeatureSet_LangVersion_values() []FeatureSet_LangVersion {
	values := make([]FeatureSet_LangVersion, len(featureSet_LangVersion_names))
	for i := range values {
		values[i] = FeatureSet_LangVersion(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: FeatureSet.LangVersion.name()
func (v FeatureSet_LangVersion) Name() string { return featureSet_LangVersion_names[v] }

// String is Enum.toString().
//
// Java: FeatureSet.LangVersion.toString()
func (v FeatureSet_LangVersion) String() string { return v.Name() }

// Features returns the features introduced by this version.
//
// Java: FeatureSet.LangVersion.features()
func (v FeatureSet_LangVersion) Features() *jutil.EnumSet[FeatureSet_Feature] {
	set := jutil.NewEnumSet[FeatureSet_Feature]()
	for _, feature := range FeatureSet_Feature_values() {
		if featureSet_Feature_version[feature] == v {
			set.Add(feature)
		}
	}
	return set
}

// FeatureSet_Feature is FeatureSet.Feature: specific features that can be included in a FeatureSet.
type FeatureSet_Feature int32

const (

	// ES3 features
	FeatureSet_Feature_REGEXP_SYNTAX FeatureSet_Feature = iota

	// ES5 features
	FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS
	FeatureSet_Feature_GETTER
	FeatureSet_Feature_KEYWORDS_AS_PROPERTIES
	FeatureSet_Feature_SETTER
	FeatureSet_Feature_STRING_CONTINUATION
	FeatureSet_Feature_TRAILING_COMMA

	// ES2015 features (besides modules): all stable browsers are now fully compliant
	// go/keep-sorted start
	FeatureSet_Feature_ARRAY_DESTRUCTURING
	FeatureSet_Feature_ARRAY_PATTERN_REST
	FeatureSet_Feature_ARROW_FUNCTIONS
	FeatureSet_Feature_BINARY_LITERALS
	FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION
	FeatureSet_Feature_CLASSES
	FeatureSet_Feature_CLASS_GETTER_SETTER
	FeatureSet_Feature_COMPUTED_PROPERTIES
	FeatureSet_Feature_CONST_DECLARATIONS
	FeatureSet_Feature_DEFAULT_PARAMETERS
	FeatureSet_Feature_FOR_OF
	FeatureSet_Feature_GENERATORS
	FeatureSet_Feature_LET_DECLARATIONS
	FeatureSet_Feature_MEMBER_DECLARATIONS
	FeatureSet_Feature_NEW_TARGET
	FeatureSet_Feature_OBJECT_DESTRUCTURING
	FeatureSet_Feature_OCTAL_LITERALS
	FeatureSet_Feature_REGEXP_FLAG_U
	FeatureSet_Feature_REGEXP_FLAG_Y
	FeatureSet_Feature_REST_PARAMETERS
	FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES
	FeatureSet_Feature_SPREAD_EXPRESSIONS
	FeatureSet_Feature_SUPER
	FeatureSet_Feature_TEMPLATE_LITERALS
	// go/keep-sorted end

	// ES modules
	FeatureSet_Feature_MODULES

	// ES 2016 only added one new feature:
	FeatureSet_Feature_EXPONENT_OP

	// ES 2017 features:
	FeatureSet_Feature_ASYNC_FUNCTIONS

	// ES 2018 adds https://github.com/tc39/proposal-object-rest-spread
	FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD
	FeatureSet_Feature_OBJECT_PATTERN_REST

	// https://github.com/tc39/proposal-async-iteration
	FeatureSet_Feature_ASYNC_GENERATORS
	FeatureSet_Feature_FOR_AWAIT_OF

	// ES 2018 adds Regex Features:
	// https://github.com/tc39/proposal-regexp-dotall-flag
	// Note: Untranspilable.
	FeatureSet_Feature_REGEXP_FLAG_S
	// https://github.com/tc39/proposal-regexp-named-groups
	// Note: Untranspilable.
	FeatureSet_Feature_REGEXP_NAMED_GROUPS
	// https://github.com/tc39/proposal-regexp-unicode-property-escapes
	// Note: Untranspilable.
	FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE
	// https://github.com/tc39/proposal-regexp-lookbehind
	// Note: Untranspilable.
	FeatureSet_Feature_REGEXP_LOOKBEHIND

	// ES 2019 adds https://github.com/tc39/proposal-json-superset
	// Note: Untranspilable.
	FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP

	// ES 2019 adds optional catch bindings:
	// https://github.com/tc39/proposal-optional-catch-binding
	FeatureSet_Feature_OPTIONAL_CATCH_BINDING

	// ES 2020 Stage 4
	FeatureSet_Feature_DYNAMIC_IMPORT
	// Note: Untranspilable.
	FeatureSet_Feature_BIGINT
	FeatureSet_Feature_IMPORT_META
	FeatureSet_Feature_NULL_COALESCE_OP
	FeatureSet_Feature_OPTIONAL_CHAINING

	// ES 2021 Stage 4
	FeatureSet_Feature_NUMERIC_SEPARATOR
	FeatureSet_Feature_LOGICAL_ASSIGNMENT

	// ES 2022 adds https://github.com/tc39/proposal-class-fields
	FeatureSet_Feature_PUBLIC_CLASS_FIELDS

	// ES 2022 adds https://github.com/tc39/proposal-class-static-block
	FeatureSet_Feature_CLASS_STATIC_BLOCK

	// ES 2022 adds https://github.com/tc39/proposal-regexp-match-indices
	// Note: Untranspilable.
	FeatureSet_Feature_REGEXP_FLAG_D

	// ES 2022 adds https://github.com/tc39/proposal-top-level-await
	// Note: Untranspilable. Needs to stay UNSUPPORTED because pass-through of top-level await does
	// not make sense as we don't emit ES6 modules so top-level await cannot be used in the output.
	FeatureSet_Feature_TOP_LEVEL_AWAIT

	// ES_NEXT: Features that are fully supported, but part of a language version that is not yet
	// fully supported

	// Polyfill implementations can target "es_next" as their fromLang so we need to ensure that
	// the "es_next" version name is distinct from the latest dated version so we don't incorrectly
	// prune those polyfills.
	FeatureSet_Feature_ES_NEXT_RUNTIME

	// ES_UNSTABLE: Features fully supported in checks, but not fully supported everywhere else

	// Polyfill implementations can target "es_unstable" as their fromLang so we need to ensure that
	// the "es_unstable" version name is distinct from the latest dated version so we don't
	// incorrectly prune those polyfills.
	FeatureSet_Feature_ES_UNSTABLE_RUNTIME

	// ES_UNSUPPORTED: Features that we can parse, but not yet supported in all checks

	// ES 2022 adds https://github.com/tc39/proposal-class-fields
	FeatureSet_Feature_PRIVATE_CLASS_PROPERTIES

	// TypeScript type syntax that will never be implemented in browsers. Only used as an indicator
	// to the CodeGenerator that it should handle type syntax.
	FeatureSet_Feature_TYPE_ANNOTATION
)

var featureSet_Feature_names = [...]string{
	"REGEXP_SYNTAX",
	"ES3_KEYWORDS_AS_IDENTIFIERS",
	"GETTER",
	"KEYWORDS_AS_PROPERTIES",
	"SETTER",
	"STRING_CONTINUATION",
	"TRAILING_COMMA",
	"ARRAY_DESTRUCTURING",
	"ARRAY_PATTERN_REST",
	"ARROW_FUNCTIONS",
	"BINARY_LITERALS",
	"BLOCK_SCOPED_FUNCTION_DECLARATION",
	"CLASSES",
	"CLASS_GETTER_SETTER",
	"COMPUTED_PROPERTIES",
	"CONST_DECLARATIONS",
	"DEFAULT_PARAMETERS",
	"FOR_OF",
	"GENERATORS",
	"LET_DECLARATIONS",
	"MEMBER_DECLARATIONS",
	"NEW_TARGET",
	"OBJECT_DESTRUCTURING",
	"OCTAL_LITERALS",
	"REGEXP_FLAG_U",
	"REGEXP_FLAG_Y",
	"REST_PARAMETERS",
	"SHORTHAND_OBJECT_PROPERTIES",
	"SPREAD_EXPRESSIONS",
	"SUPER",
	"TEMPLATE_LITERALS",
	"MODULES",
	"EXPONENT_OP",
	"ASYNC_FUNCTIONS",
	"OBJECT_LITERALS_WITH_SPREAD",
	"OBJECT_PATTERN_REST",
	"ASYNC_GENERATORS",
	"FOR_AWAIT_OF",
	"REGEXP_FLAG_S",
	"REGEXP_NAMED_GROUPS",
	"REGEXP_UNICODE_PROPERTY_ESCAPE",
	"REGEXP_LOOKBEHIND",
	"UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP",
	"OPTIONAL_CATCH_BINDING",
	"DYNAMIC_IMPORT",
	"BIGINT",
	"IMPORT_META",
	"NULL_COALESCE_OP",
	"OPTIONAL_CHAINING",
	"NUMERIC_SEPARATOR",
	"LOGICAL_ASSIGNMENT",
	"PUBLIC_CLASS_FIELDS",
	"CLASS_STATIC_BLOCK",
	"REGEXP_FLAG_D",
	"TOP_LEVEL_AWAIT",
	"ES_NEXT_RUNTIME",
	"ES_UNSTABLE_RUNTIME",
	"PRIVATE_CLASS_PROPERTIES",
	"TYPE_ANNOTATION",
}

// featureSet_Feature_name is the Java field Feature.name.
var featureSet_Feature_name = [...]string{
	"RegExp syntax",
	"ES3 keywords as identifiers",
	"getters",
	"reserved words as properties",
	"setters",
	"string continuation",
	"trailing comma",
	"array destructuring",
	"array pattern rest",
	"arrow function",
	"binary literal",
	"block-scoped function declaration",
	"class",
	"class getters/setters",
	"computed property",
	"const declaration",
	"default parameter",
	"for-of loop",
	"generator",
	"let declaration",
	"member declaration",
	"new.target",
	"object destructuring",
	"octal literal",
	"RegExp flag 'u'",
	"RegExp flag 'y'",
	"rest parameter",
	"shorthand object property",
	"spread expression",
	"super",
	"template literal",
	"modules",
	"exponent operator (**)",
	"async function",
	"object literals with spread",
	"object pattern rest",
	"async generator functions",
	"for-await-of loop",
	"RegExp flag 's'",
	"RegExp named groups",
	"RegExp unicode property escape",
	"RegExp Lookbehind",
	"Unescaped unicode line or paragraph separator",
	"Optional catch binding",
	"Dynamic module import",
	"bigint",
	"import.meta",
	"Nullish coalescing",
	"Optional chaining",
	"numeric separator",
	"Logical assignments",
	"Public class fields",
	"Class static block",
	"RegExp flag 'd'",
	"Top-level await",
	"es_next runtime",
	"es_unstable runtime",
	"Private class properties",
	"type annotation",
}

// featureSet_Feature_version is the Java field Feature.version.
var featureSet_Feature_version = [...]FeatureSet_LangVersion{
	FeatureSet_LangVersion_ES3,
	FeatureSet_LangVersion_ES5,
	FeatureSet_LangVersion_ES5,
	FeatureSet_LangVersion_ES5,
	FeatureSet_LangVersion_ES5,
	FeatureSet_LangVersion_ES5,
	FeatureSet_LangVersion_ES5,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2015,
	FeatureSet_LangVersion_ES2016,
	FeatureSet_LangVersion_ES2017,
	FeatureSet_LangVersion_ES2018,
	FeatureSet_LangVersion_ES2018,
	FeatureSet_LangVersion_ES2018,
	FeatureSet_LangVersion_ES2018,
	FeatureSet_LangVersion_ES2018,
	FeatureSet_LangVersion_ES2018,
	FeatureSet_LangVersion_ES2018,
	FeatureSet_LangVersion_ES2018,
	FeatureSet_LangVersion_ES2019,
	FeatureSet_LangVersion_ES2019,
	FeatureSet_LangVersion_ES2020,
	FeatureSet_LangVersion_ES2020,
	FeatureSet_LangVersion_ES2020,
	FeatureSet_LangVersion_ES2020,
	FeatureSet_LangVersion_ES2020,
	FeatureSet_LangVersion_ES2021,
	FeatureSet_LangVersion_ES2021,
	FeatureSet_LangVersion_ES2022,
	FeatureSet_LangVersion_ES2022,
	FeatureSet_LangVersion_ES2022,
	FeatureSet_LangVersion_ES_UNSUPPORTED,
	FeatureSet_LangVersion_ES_NEXT,
	FeatureSet_LangVersion_ES_UNSTABLE,
	FeatureSet_LangVersion_ES_UNSUPPORTED,
	FeatureSet_LangVersion_TYPESCRIPT,
}

// FeatureSet_Feature_values is Feature.values().
//
// Java: FeatureSet.Feature.values()
func FeatureSet_Feature_values() []FeatureSet_Feature {
	values := make([]FeatureSet_Feature, len(featureSet_Feature_names))
	for i := range values {
		values[i] = FeatureSet_Feature(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: FeatureSet.Feature.name()
func (f FeatureSet_Feature) Name() string { return featureSet_Feature_names[f] }

// String returns the feature's human-readable name.
//
// Java: FeatureSet.Feature.toString()
func (f FeatureSet_Feature) String() string {
	return featureSet_Feature_name[f]
}

// NewFeatureSet constructs a FeatureSet.
//
// Java: FeatureSet(EnumSet<Feature>)
func NewFeatureSet(features *jutil.EnumSet[FeatureSet_Feature]) *FeatureSet {
	f := &FeatureSet{}
	// ImmutableSet will only use an EnumSet if the set starts as an EnumSet.
	f.features = jutil.EnumSetCopyOf(features)
	return f
}

// Version returns a string representation suitable for use in polyfill
// definition files and encoding in depgraph and deps.js files.
//
// Java: FeatureSet.version()
func (f *FeatureSet) Version() string {
	if FeatureSet_ES3.Contains(f) {
		return "es3"
	}
	if FeatureSet_ES5.Contains(f) {
		return "es5"
	}
	if FeatureSet_ES2015_MODULES.Contains(f) {
		return "es6"
	}
	if FeatureSet_ES2016_MODULES.Contains(f) {
		return "es7"
	}
	if FeatureSet_ES2017_MODULES.Contains(f) {
		return "es8"
	}
	if FeatureSet_ES2018_MODULES.Contains(f) {
		return "es9"
	}
	if FeatureSet_ES2019_MODULES.Contains(f) {
		return "es_2019"
	}
	if FeatureSet_ES2020_MODULES.Contains(f) {
		return "es_2020"
	}
	if FeatureSet_ES2021_MODULES.Contains(f) {
		return "es_2021"
	}
	if FeatureSet_ES2022_MODULES.Contains(f) {
		return "es_2022"
	}
	if FeatureSet_ES_NEXT.Contains(f) {
		return "es_next"
	}
	if FeatureSet_ES_UNSTABLE.Contains(f) {
		return "es_unstable"
	}
	if FeatureSet_ES_UNSUPPORTED.Contains(f) {
		return "es_unsupported"
	}
	if FeatureSet_ALL.Contains(f) {
		return "all"
	}
	panic(jutil.NewIllegalStateExceptionMsg(f.String()))
}

// VersionForDebugging returns a string representation useful for debugging.
//
// Deprecated: Please use Version() instead.
//
// Java: FeatureSet.versionForDebugging()
func (f *FeatureSet) VersionForDebugging() string {
	return f.Version()
}

// WithoutFeatures returns this set minus the given features.
//
// Java: FeatureSet.without(Feature, Feature...)
func (f *FeatureSet) WithoutFeatures(featureToRemove FeatureSet_Feature, moreFeaturesToRemove ...FeatureSet_Feature) *FeatureSet {
	removed := jutil.NewEnumSet(featureToRemove)
	for _, feature := range moreFeaturesToRemove {
		removed.Add(feature)
	}
	return NewFeatureSet(FeatureSet_Difference(f.features, removed))
}

// Without returns this set minus the features of other.
//
// Java: FeatureSet.without(FeatureSet)
func (f *FeatureSet) Without(other *FeatureSet) *FeatureSet {
	return NewFeatureSet(FeatureSet_Difference(f.features, other.features))
}

// WithoutTypes returns this set minus the TypeScript features.
//
// Java: FeatureSet.withoutTypes()
func (f *FeatureSet) WithoutTypes() *FeatureSet {
	return NewFeatureSet(FeatureSet_Difference(f.features, FeatureSet_LangVersion_TYPESCRIPT.Features()))
}

// Union returns a new FeatureSet including all features of both this and
// other.
//
// Java: FeatureSet.union(FeatureSet)
func (f *FeatureSet) Union(other *FeatureSet) *FeatureSet {
	return NewFeatureSet(FeatureSet_Union(f.features, other.features))
}

// ContainsAtLeastOneOf reports whether this FeatureSet contains at least one
// of the features of other.
//
// Java: FeatureSet.containsAtLeastOneOf(FeatureSet)
func (f *FeatureSet) ContainsAtLeastOneOf(other *FeatureSet) bool {
	for otherFeature := range other.features.All() {
		if f.features.Contains(otherFeature) {
			return true
		}
	}
	return false
}

// Contains reports whether this FeatureSet contains all of the features of
// other.
//
// Java: FeatureSet.contains(FeatureSet)
func (f *FeatureSet) Contains(other *FeatureSet) bool {
	return f.features.ContainsAll(other.features)
}

// ContainsFeature reports whether this FeatureSet contains the given
// feature.
//
// Java: FeatureSet.contains(Feature)
func (f *FeatureSet) ContainsFeature(feature FeatureSet_Feature) bool {
	return f.features.ContainsAll(jutil.NewEnumSet(feature))
}

// FeatureSet_EmptyEnumSet returns an empty feature set.
//
// Java: FeatureSet.emptyEnumSet()
func FeatureSet_EmptyEnumSet() *jutil.EnumSet[FeatureSet_Feature] {
	return jutil.NewEnumSet[FeatureSet_Feature]()
}

// FeatureSet_EnumSetOf copies set.
//
// Java: FeatureSet.enumSetOf(Set<Feature>)
func FeatureSet_EnumSetOf(set *jutil.EnumSet[FeatureSet_Feature]) *jutil.EnumSet[FeatureSet_Feature] {
	if set.IsEmpty() {
		return FeatureSet_EmptyEnumSet()
	}
	return jutil.EnumSetCopyOf(set)
}

// FeatureSet_Add returns features plus feature.
//
// Java: FeatureSet.add(Set<Feature>, Feature)
func FeatureSet_Add(features *jutil.EnumSet[FeatureSet_Feature], feature FeatureSet_Feature) *jutil.EnumSet[FeatureSet_Feature] {
	result := FeatureSet_EnumSetOf(features)
	result.Add(feature)
	return result
}

// FeatureSet_Union returns features plus newFeatures.
//
// Java: FeatureSet.union(Set<Feature>, Set<Feature>)
func FeatureSet_Union(features, newFeatures *jutil.EnumSet[FeatureSet_Feature]) *jutil.EnumSet[FeatureSet_Feature] {
	result := FeatureSet_EnumSetOf(features)
	result.AddAll(newFeatures)
	return result
}

// FeatureSet_Difference returns features minus removedFeatures.
//
// Java: FeatureSet.difference(Set<Feature>, Set<Feature>)
func FeatureSet_Difference(features, removedFeatures *jutil.EnumSet[FeatureSet_Feature]) *jutil.EnumSet[FeatureSet_Feature] {
	result := FeatureSet_EnumSetOf(features)
	result.RemoveAll(removedFeatures)
	return result
}

// With returns a feature set combining all the features from this and
// feature.
//
// Java: FeatureSet.with(Feature)
func (f *FeatureSet) With(feature FeatureSet_Feature) *FeatureSet {
	if f.features.Contains(feature) {
		return f
	}
	return NewFeatureSet(FeatureSet_Add(f.features, feature))
}

// WithFeatures returns a feature set combining all the features from this
// and newFeatures.
//
// Java: FeatureSet.with(Feature...)
func (f *FeatureSet) WithFeatures(newFeatures ...FeatureSet_Feature) *FeatureSet {
	return NewFeatureSet(FeatureSet_Union(f.features, jutil.NewEnumSet(newFeatures...)))
}

// WithEnumSet returns a feature set combining all the features from this
// and newFeatures.
//
// Java: FeatureSet.with(Set<Feature>)
func (f *FeatureSet) WithEnumSet(newFeatures *jutil.EnumSet[FeatureSet_Feature]) *FeatureSet {
	return NewFeatureSet(FeatureSet_Union(f.features, newFeatures))
}

// WithFeatureSet returns a feature set combining all the features from this
// and newFeatures.
//
// Java: FeatureSet.with(FeatureSet)
func (f *FeatureSet) WithFeatureSet(newFeatures *FeatureSet) *FeatureSet {
	return NewFeatureSet(FeatureSet_Union(f.features, newFeatures.features))
}

// Has reports whether this FeatureSet includes feature.
//
// Java: FeatureSet.has(Feature)
func (f *FeatureSet) Has(feature FeatureSet_Feature) bool {
	return f.features.Contains(feature)
}

// GetFeatures returns the (immutable; callers must not modify it) set of
// features.
//
// Java: FeatureSet.getFeatures()
func (f *FeatureSet) GetFeatures() *jutil.EnumSet[FeatureSet_Feature] {
	return f.features
}

// Equals is FeatureSet.equals(Object).
//
// Java: FeatureSet.equals(Object)
func (f *FeatureSet) Equals(other any) bool {
	featureSet, ok := other.(*FeatureSet)
	return ok && featureSet != nil && featureSet.features.Equals(f.features)
}

// HashCode is FeatureSet.hashCode().
//
// Java: FeatureSet.hashCode()
func (f *FeatureSet) HashCode() int32 {
	return f.features.HashCode()
}

// String renders the features, e.g. "[RegExp syntax, getters]".
//
// Java: FeatureSet.toString()
func (f *FeatureSet) String() string {
	return f.features.String()
}

// FeatureSet_ValueOf parses known strings into feature sets.
//
// Java: FeatureSet.valueOf(String)
func FeatureSet_ValueOf(name string) *FeatureSet {
	switch name {
	case "es3":
		return FeatureSet_ES3
	case "es5":
		return FeatureSet_ES5
	case "es_2015", "es6":
		return FeatureSet_ES2015
	case "es_2016", "es7":
		return FeatureSet_ES2016
	case "es_2017", "es8":
		return FeatureSet_ES2017
	case "es_2018", "es9":
		return FeatureSet_ES2018
	case "es_2019":
		return FeatureSet_ES2019
	case "es_2020":
		return FeatureSet_ES2020
	case "es_2021":
		return FeatureSet_ES2021
	case "es_2022":
		return FeatureSet_ES2022
	case "es_next":
		return FeatureSet_ES_NEXT
	case "es_unstable":
		return FeatureSet_ES_UNSTABLE
	case "es_unsupported":
		return FeatureSet_ES_UNSUPPORTED
	case "all":
		return FeatureSet_ALL
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("No such FeatureSet: " + name))
	}
}

// FeatureSet_BrowserFeatureSetValueOf is a variant of FeatureSet_ValueOf
// that outputs a BROWSER_20XX feature set equivalent. This will generally
// exclude some not-yet-implemented features from a given ES20XX while also
// including Feature MODULES.
//
// Note: The returned feature set will include the MODULES feature which you
// may want to exclude manually.
//
// Java: FeatureSet.browserFeatureSetValueOf(String)
func FeatureSet_BrowserFeatureSetValueOf(name string) *FeatureSet {
	featureSet := FeatureSet_ValueOf(name)
	if featureSet.Contains(FeatureSet_ES_NEXT) {
		// For ES_NEXT and ES_UNSTABLE, we need to return the raw featureSet to retain the RUNTIME
		// features which pin polyfills to future versions.
		return featureSet
	}
	if featureSet.Contains(FeatureSet_ES2022) {
		return FeatureSet_BROWSER_2023
	}
	if featureSet.Contains(FeatureSet_ES2021) {
		return FeatureSet_BROWSER_2022
	}
	if featureSet.Contains(FeatureSet_ES2020) {
		return FeatureSet_BROWSER_2021
	}
	if featureSet.Contains(FeatureSet_ES2019) {
		return FeatureSet_BROWSER_2020
	}
	if featureSet.Contains(FeatureSet_ES2018) {
		return FeatureSet_BROWSER_2019
	}
	return featureSet
}

// FeatureSet_All returns a FeatureSet containing all known features.
//
// NOTE: PassFactory classes that claim to support FeatureSet.all() should
// be only those that cannot be broken by new features being added to the
// language. Mainly these are passes that don't have to actually look at the
// AST at all, like empty marker passes.
//
// Java: FeatureSet.all()
func FeatureSet_All() *FeatureSet {
	return FeatureSet_ALL
}

// FeatureSet_Latest returns the latest parseable feature set.
//
// Java: FeatureSet.latest()
func FeatureSet_Latest() *FeatureSet {
	return FeatureSet_ES_UNSUPPORTED
}
