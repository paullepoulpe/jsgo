package main

// The sweep's flag configurations (group "flags", opt-in): the breadth of
// CommandLineRunner's flags, each value at least once, the important ones
// combined pairwise. Run them on a stratified sample, each configuration
// on its own share of it:
//
//	go run ./tools/clie2e -suite sweep -configs @flags -sample 2000 -per-config 320 ...
//
// A configuration's variants rotate over its files, so a flag with many
// values (--language_in, --browser_featureset_year, the diagnostic groups)
// costs one configuration, not one per value. Output files go to %res%,
// the job's result directory, and are compared like stdout.
//
// Single-file compiles of npm packages mostly stop on errors before the
// optimizations (an ES module's imports, CommonJS's undeclared require and
// exports), so most configurations switch those off (lax and laxAdv), as
// the sweep's -lax configurations do.

import "strings"

var (
	simple   = []string{"--compilation_level", "SIMPLE"}
	advanced = []string{"--compilation_level", "ADVANCED"}
	lax      = []string{"--jscomp_off=moduleLoad"}
	laxAdv   = []string{"--jscomp_off=moduleLoad", "--jscomp_off=undefinedVars"}
	// closureLax lets a Closure Library file compile without the files it
	// requires.
	closureLax = []string{"--jscomp_off=missingSourcesWarnings"}
)

func cat(parts ...[]string) []string {
	var out []string
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// split turns "a b c" into a variant.
func split(s string) []string { return strings.Fields(s) }

func variants(ss ...string) [][]string {
	out := make([][]string, len(ss))
	for i, s := range ss {
		out[i] = split(s)
	}
	return out
}

// diagnosticGroups are DiagnosticGroups.getRegisteredGroups() of the
// pinned compiler, in registration order: what --jscomp_error,
// --jscomp_warning and --jscomp_off accept besides "*".
var diagnosticGroups = split(`untranspilableFeatures moduleLoad moduleImport globalThis deprecated
underscore visibility accessControls nonStandardJsDocs invalidCasts strictModuleDepCheck
violatedModuleDep externsValidation unknownDefines tweakValidation missingOverride
missingProperties globallyMissingProperties j2clChecks missingReturn undefinedVars
checkDebuggerStatement checkRegExp checkTypes checkEs5InheritanceCorrectnessConditions
checkPrototypalTypes tooManyTypeParams strictMissingProperties strictPrimitiveOperators
strictCheckTypes reportUnknownTypes checkVars uselessCode const constantProperty
accessControlsConst typeInvalidation duplicate es5Strict missingProvide unrecognizedTypeError
dangerousUnrecognizedTypeError missingRequire missingSourcesWarnings extraRequire
duplicateMessage msgDescriptions messageConventions misplacedTypeAnnotation
misplacedMsgAnnotation misplacedSuppress suspiciousCode functionParams deprecatedAnnotations
unusedLocalVariables jsdocMissingType typeImportCodeReferences partialAlias useOfGoogProvide
lintVarDeclarations lintChecks strictModuleChecks analyzerChecks closureDepMethodUsageChecks
closureClassChecks transitionalSuspiciousCodeWarnings conformanceViolations lateProvide
missingPolyfill polymer boundedGenerics closureUnawareCodeAnnotationPresent polymerBehavior
invalidProvide conflictingIjsFile implicitWeakEntryPoint artificialFunctionPurityValidation`)

// jscompVariants deals the diagnostic groups into 7 blocks; variant i makes
// block i errors, block i+2 warnings and block i+4 off, the flags
// interleaved. Over the 7 variants each group is once an error, once a
// warning and once off.
func jscompVariants() [][]string {
	const k = 7
	blocks := make([][]string, k)
	for i, g := range diagnosticGroups {
		blocks[i%k] = append(blocks[i%k], g)
	}
	out := make([][]string, k)
	for i := range k {
		e, w, o := blocks[i], blocks[(i+2)%k], blocks[(i+4)%k]
		var v []string
		for j := 0; j < len(e) || j < len(w) || j < len(o); j++ {
			if j < len(e) {
				v = append(v, "--jscomp_error="+e[j])
			}
			if j < len(w) {
				v = append(v, "--jscomp_warning="+w[j])
			}
			if j < len(o) {
				v = append(v, "--jscomp_off="+o[j])
			}
		}
		out[i] = v
	}
	return out
}

// languageIns are --language_in's values, with the abbreviations that
// LanguageMode.fromString accepts.
var languageIns = []string{"ECMASCRIPT3", "ES5", "ECMASCRIPT5_STRICT", "ECMASCRIPT_2015", "ES6",
	"ECMASCRIPT_2016", "ES_2017", "ECMASCRIPT_2018", "ECMASCRIPT_2019", "ECMASCRIPT_2020",
	"ECMASCRIPT_2021", "ECMASCRIPT_2022", "ECMASCRIPT_NEXT", "STABLE", "UNSTABLE"}

var languageOuts = []string{"ECMASCRIPT3", "ECMASCRIPT5", "ECMASCRIPT_2015", "ECMASCRIPT_2016",
	"ECMASCRIPT_2017", "ECMASCRIPT_2018", "ECMASCRIPT_2019", "ECMASCRIPT_2020", "ECMASCRIPT_2021",
	"STABLE", "NO_TRANSPILE", "ECMASCRIPT_NEXT"}

// languagePairs cover every input and output mode at least once, the
// down-levelling pairs first.
var languagePairs = variants(
	"--language_in=ECMASCRIPT3 --language_out=ECMASCRIPT3",
	"--language_in=ECMASCRIPT5 --language_out=ECMASCRIPT3",
	"--language_in=ECMASCRIPT5_STRICT --language_out=ECMASCRIPT5",
	"--language_in=ECMASCRIPT_2015 --language_out=ECMASCRIPT5",
	"--language_in=ECMASCRIPT_2016 --language_out=ECMASCRIPT_2015",
	"--language_in=ECMASCRIPT_2017 --language_out=ECMASCRIPT5",
	"--language_in=ECMASCRIPT_2018 --language_out=ECMASCRIPT_2017",
	"--language_in=ECMASCRIPT_2019 --language_out=ECMASCRIPT_2015",
	"--language_in=ECMASCRIPT_2020 --language_out=ECMASCRIPT_2019",
	"--language_in=ECMASCRIPT_2021 --language_out=ECMASCRIPT_2017",
	"--language_in=ECMASCRIPT_2022 --language_out=ECMASCRIPT_2019",
	"--language_in=ECMASCRIPT_NEXT --language_out=ECMASCRIPT_2021",
	"--language_in=ECMASCRIPT_NEXT --language_out=ECMASCRIPT3",
	"--language_in=STABLE --language_out=NO_TRANSPILE",
	"--language_in=UNSTABLE --language_out=ECMASCRIPT5",
	"--language_in=UNSTABLE --language_out=ECMASCRIPT_2020",
)

func prefixed(flag string, values []string) [][]string {
	out := make([][]string, len(values))
	for i, v := range values {
		out[i] = []string{flag + "=" + v}
	}
	return out
}

func init() {
	years := []string{"2012", "2018", "2019", "2020", "2021", "2022", "2023", "2024", "2025", "2026"}
	closureBase := "in/google-closure-library@20230802.0.0/package/closure/goog/base.js"
	flagConfigs := []config{
		// --language_in, --language_out.
		{name: "f-lang-in", args: cat(simple, lax), variants: prefixed("--language_in", languageIns)},
		{name: "f-lang-out", args: cat(simple, lax), variants: prefixed("--language_out", languageOuts)},
		{name: "f-lang-pairs", args: cat(simple, lax), variants: languagePairs},
		{name: "f-lang-pairs-adv", args: cat(advanced, laxAdv), variants: languagePairs},
		{name: "f-es5", args: cat(simple, lax, split("--language_out=ECMASCRIPT5"))},
		{name: "f-es3-adv", args: cat(advanced, laxAdv, split("--language_out=ECMASCRIPT3"))},
		{name: "f-es2015-adv", args: cat(advanced, laxAdv, split("--language_out=ECMASCRIPT_2015"))},

		// --compilation_level=WHITESPACE_ONLY, --formatting.
		{name: "f-whitespace", args: split("--compilation_level=WHITESPACE_ONLY")},
		{name: "f-whitespace-lang", args: split("--compilation_level=WHITESPACE_ONLY"),
			variants: prefixed("--language_out", []string{"ECMASCRIPT3", "ECMASCRIPT5", "ECMASCRIPT_2015", "ECMASCRIPT_2017", "ECMASCRIPT_2019", "ECMASCRIPT_2021", "NO_TRANSPILE"})},
		{name: "f-format", args: cat(simple, lax), variants: variants(
			"--formatting=PRETTY_PRINT",
			"--formatting=PRINT_INPUT_DELIMITER",
			"--formatting=SINGLE_QUOTES",
			"--formatting=PRETTY_PRINT --formatting=PRINT_INPUT_DELIMITER --formatting=SINGLE_QUOTES")},
		{name: "f-format-adv", args: cat(advanced, laxAdv, split("--formatting PRETTY_PRINT --formatting PRINT_INPUT_DELIMITER --formatting SINGLE_QUOTES"))},
		{name: "f-format-ws", args: split("-O WHITESPACE_ONLY --formatting=PRETTY_PRINT --formatting=SINGLE_QUOTES")},

		// --warning_level, --checks_only.
		{name: "f-warning-level", args: simple, variants: variants("--warning_level=QUIET", "--warning_level=DEFAULT", "--warning_level=VERBOSE", "-W VERBOSE --language_out=ECMASCRIPT5")},
		{name: "f-adv-verbose", args: cat(advanced, split("--warning_level=VERBOSE"))},
		{name: "f-checks-verbose", args: split("--checks_only --warning_level=VERBOSE")},

		// --jscomp_error, --jscomp_warning, --jscomp_off.
		{name: "f-jscomp-checks", args: split("--checks_only -W VERBOSE"), variants: jscompVariants()},
		{name: "f-jscomp-simple", args: simple, variants: jscompVariants()},
		{name: "f-jscomp-adv", args: cat(advanced, laxAdv), variants: jscompVariants()},
		{name: "f-jscomp-star", args: split("--checks_only"), variants: variants(
			"--jscomp_warning=*",
			"--jscomp_error=*",
			"--jscomp_off=*",
			"-W VERBOSE --jscomp_off=* --jscomp_error=checkTypes",
			"--jscomp_error=* --jscomp_off=moduleLoad --jscomp_warning=undefinedVars")},

		// Polyfills and runtime libraries.
		{name: "f-no-rewrite-polyfills", args: cat(simple, lax, split("--language_out=ECMASCRIPT5 --rewrite_polyfills=false"))},
		{name: "f-no-inject-libraries", args: cat(advanced, laxAdv, split("--language_out=ECMASCRIPT5 --inject_libraries=false"))},
		{name: "f-force-inject", args: cat(simple, lax), variants: variants(
			"--force_inject_library=base",
			"--force_inject_library=es6_runtime",
			"--force_inject_library=es6_runtime --language_out=ECMASCRIPT5 --rewrite_polyfills=false",
			"--force_inject_library=es6/promise/promise --force_inject_library=util/global",
			"--force_inject_library=es6_runtime --inject_libraries=false",
			"--force_inject_library=base --isolate_polyfills --language_out=ECMASCRIPT5")},

		// --define, -D.
		{name: "f-define", args: cat(simple, lax), variants: variants(
			"--define=goog.DEBUG=false",
			"-D COMPILED",
			"--define goog.FEATURESET_YEAR=2019",
			"-D goog.LOCALE='de'",
			"--define=ns.NUM=1.5 --define=ns.NEG=-3 --D ns.HEX=0x10",
			"--D ns.STR='a b' -D ns.FLAG=true --define=ns.OFF=false")},
		{name: "f-define-adv", args: cat(advanced, laxAdv), variants: variants(
			"-D goog.DEBUG=false -D goog.LOCALE='fr'",
			"--define=COMPILED=true --define=goog.FEATURESET_YEAR=2021",
			"--define=ns.STR='' -D ns.NUM=0")},
		// After Closure Library's base.js, with missingSourcesWarnings off
		// (a namespace that no input provides is then not an error).
		{name: "f-define-closure", args: cat(lax, closureLax, []string{"--js", closureBase}), variants: variants(
			"-O SIMPLE --define=goog.DEBUG=false --define=goog.LOCALE='ja'",
			"-O ADVANCED -D goog.FEATURESET_YEAR=2018 -D goog.TRUSTED_SITE=false",
			"-O SIMPLE --define goog.DISALLOW_TEST_ONLY_CODE --define=goog.ENABLE_DEBUG_LOADER=true --formatting=PRETTY_PRINT",
			"--checks_only -W VERBOSE -D goog.DEBUG=false")},

		// --debug, --strict_mode_input, --emit_use_strict.
		{name: "f-debug", args: cat(advanced, laxAdv), variants: variants("--debug", "--debug=true --language_out=ECMASCRIPT5", "--debug --formatting=PRETTY_PRINT")},
		{name: "f-strict-input-off", args: cat(simple, lax, split("--strict_mode_input=false")), variants: variants("", "--language_out=ECMASCRIPT5")},
		{name: "f-emit-use-strict", args: cat(simple, lax), variants: variants(
			"--emit_use_strict=false",
			"--emit_use_strict",
			"--emit_use_strict=true --language_out=ECMASCRIPT5",
			"--emit_use_strict=false --language_out=ECMASCRIPT5")},
		{name: "f-strict-adv", args: cat(advanced, laxAdv, split("--strict_mode_input=false --emit_use_strict"))},

		// --browser_featureset_year.
		{name: "f-featureset-year", args: cat(simple, lax), variants: prefixed("--browser_featureset_year", years)},
		{name: "f-featureset-year-adv", args: cat(advanced, laxAdv), variants: prefixed("--browser_featureset_year", years)},

		// Types, static inheritance, Polymer, Angular, Closure primitives.
		{name: "f-use-types", args: cat(advanced, laxAdv), variants: variants(
			"--use_types_for_optimization=false",
			"--use_types_for_optimization --jscomp_warning=checkTypes",
			"--use_types_for_optimization=false --jscomp_warning=checkTypes",
			"--use_types_for_optimization=false --jscomp_off=missingSourcesWarnings --js "+closureBase)},
		{name: "f-static-inheritance", args: cat(advanced, laxAdv), variants: variants(
			"--assume_static_inheritance_is_not_used=false",
			"--assume_static_inheritance_is_not_used=true --language_out=ECMASCRIPT5",
			"--assume_static_inheritance_is_not_used=false --language_out=ECMASCRIPT5")},
		{name: "f-polymer", args: lax, variants: variants(
			"-O ADVANCED --jscomp_off=undefinedVars --polymer_version=1",
			"-O ADVANCED --jscomp_off=undefinedVars --polymer_version=2",
			"-O SIMPLE --polymer_version=1",
			"-O SIMPLE --polymer_version=2 --language_out=ECMASCRIPT5",
			"--checks_only --polymer_version=2")},
		{name: "f-angular", args: lax, variants: variants(
			"-O SIMPLE --angular_pass",
			"-O ADVANCED --jscomp_off=undefinedVars --angular_pass",
			"--checks_only --angular_pass=true")},
		{name: "f-no-closure-primitives", args: lax, variants: variants(
			"-O SIMPLE --process_closure_primitives=false",
			"-O ADVANCED --jscomp_off=undefinedVars --process_closure_primitives=false",
			"--checks_only --process_closure_primitives=false",
			"-O SIMPLE --jscomp_off=missingSourcesWarnings --process_closure_primitives=false --js "+closureBase)},

		// Renaming reports.
		{name: "f-renaming-reports", args: cat(advanced, laxAdv, split("--variable_renaming_report %res%/vars.map --property_renaming_report %res%/props.map")),
			variants: variants("", "--debug", "--formatting=PRETTY_PRINT --language_out=ECMASCRIPT5")},
		{name: "f-renaming-reports-simple", args: cat(simple, lax, split("--variable_renaming_report=%res%/vars.map --property_renaming_report=%res%/props.map"))},

		// Pairwise combinations of the important ones.
		{name: "f-combo-adv-es5-debug", args: cat(advanced, laxAdv, split("--language_out=ECMASCRIPT5 --formatting=PRETTY_PRINT --debug --use_types_for_optimization=false"))},
		{name: "f-combo-checks-strict", args: split("--checks_only -W VERBOSE --language_in=ECMASCRIPT_2017 --jscomp_error=checkTypes --jscomp_off=lintChecks --strict_mode_input=false")},
		{name: "f-combo-ws-es5", args: split("-O WHITESPACE_ONLY --language_out=ECMASCRIPT5 --formatting=SINGLE_QUOTES --emit_use_strict --rewrite_polyfills=false")},
		{name: "f-combo-simple-quiet", args: cat(simple, lax, split("-W QUIET --language_out=ECMASCRIPT_2015 --define=goog.DEBUG=false --strict_mode_input=false --formatting=PRINT_INPUT_DELIMITER"))},
		{name: "f-combo-adv-year", args: cat(advanced, laxAdv, split("--browser_featureset_year=2018 --assume_static_inheritance_is_not_used=false --angular_pass --polymer_version=2 -W VERBOSE"))},
		{name: "f-combo-adv-es3", args: cat(advanced, laxAdv, split("--process_closure_primitives=false --inject_libraries=false --language_out=ECMASCRIPT3 --variable_renaming_report=%res%/vars.map"))},

		// Flags that fail before the inputs are read: a few files each.
		{name: "f-flag-errors", files: 40, variants: variants(
			"--polymer_pass",
			"--polymer_version=3",
			"--browser_featureset_year=2015",
			"--browser_featureset_year=2020 --language_out=ECMASCRIPT5",
			"--browser_featureset_year=2020 --language_out=ECMASCRIPT_NEXT",
			"--language_in=ES2015",
			"--language_in=NO_TRANSPILE",
			"--language_out=UNSUPPORTED",
			"--formatting=PRETTY_PRINT,SINGLE_QUOTES",
			"--define=goog.LOCALE=de",
			"--define=ns.X=1+1",
			"--force_inject_library=nosuchlibrary",
			"--jscomp_warning=nosuchgroup",
			"--warning_level=verbose",
			"--use_types_for_optimization=maybe",
			"--emit_use_strict --chunk_output_type=ES_MODULES",
		)},

		// Beyond the list: the other flags that neither the project-build
		// sweep (inputs, modules, chunks, output wrappers) nor the ones
		// above cover. Fewer files each.
		{name: "f-levels", files: 160, args: lax, variants: variants(
			"-O TRANSPILE_ONLY",
			"-O TRANSPILE_ONLY --language_out=ECMASCRIPT5",
			"-O BUNDLE",
			"-O BUNDLE --language_out=ECMASCRIPT5 --formatting=PRETTY_PRINT")},
		{name: "f-env-custom", files: 160, variants: variants(
			"--env=CUSTOM",
			"--env=CUSTOM -O ADVANCED",
			"--env=CUSTOM --checks_only -W VERBOSE",
			"--env=BROWSER -O ADVANCED --jscomp_off=undefinedVars")},
		{name: "f-error-format", files: 160, variants: variants(
			"--error_format=JSON",
			"--error_format=JSON -W VERBOSE",
			"--error_format=JSON -O ADVANCED",
			"--error_format=STANDARD --checks_only -W VERBOSE")},
		{name: "f-check-flags", files: 160, args: split("-W VERBOSE"), variants: variants(
			"--extra_annotation_name=foo --extra_annotation_name=nocollapse2",
			"--hide_warnings_for=in/",
			"--hide_warnings_for=/package/src/ --checks_only",
			"--third_party",
			"--third_party --checks_only")},
		{name: "f-exports", files: 160, args: cat(advanced, laxAdv), variants: variants(
			"--generate_exports",
			"--generate_exports=false",
			"--generate_exports --export_local_property_definitions=false",
			"--assume_function_wrapper",
			"--assume_function_wrapper --generate_exports --language_out=ECMASCRIPT5")},
		{name: "f-rename-prefix", files: 160, args: cat(advanced, laxAdv), variants: variants(
			"--rename_variable_prefix=pfx_",
			"--rename_prefix_namespace=NS",
			"--rename_variable_prefix=$x --variable_renaming_report=%res%/vars.map")},
		{name: "f-polyfill-misc", files: 160, args: cat(simple, lax), variants: variants(
			"--isolate_polyfills --language_out=ECMASCRIPT5",
			"--isolate_polyfills --language_out=ECMASCRIPT_2015 -O ADVANCED --jscomp_off=undefinedVars",
			"--dynamic_import_alias=imp",
			"--dynamic_import_alias=import --language_out=ECMASCRIPT_2015")},
		{name: "f-charset-json", files: 160, args: cat(simple, lax), variants: variants(
			"--charset=UTF-8",
			"--charset=ISO-8859-1",
			"--charset=US-ASCII --formatting=PRETTY_PRINT",
			"--json_streams=OUT",
			"--json_streams=OUT --create_source_map=%res%/out.js.map")},

		// A pairwise covering array of the important flags: every pair of
		// values of two of them is in some variant.
		{name: "f-pairwise", args: lax, variants: pairwise([][]string{
			{"-O WHITESPACE_ONLY", "-O SIMPLE", "-O ADVANCED --jscomp_off=undefinedVars"},
			{"", "--language_out=ECMASCRIPT3", "--language_out=ECMASCRIPT5", "--language_out=ECMASCRIPT_2015", "--language_out=ECMASCRIPT_2017", "--language_out=ECMASCRIPT_2019", "--language_out=ECMASCRIPT_2021", "--language_out=NO_TRANSPILE"},
			{"", "--formatting=PRETTY_PRINT", "--formatting=SINGLE_QUOTES", "--formatting=PRINT_INPUT_DELIMITER"},
			{"-W QUIET", "", "-W VERBOSE"},
			{"", "--rewrite_polyfills=false"},
			{"", "--inject_libraries=false"},
			{"", "--debug"},
			{"", "--strict_mode_input=false"},
			{"", "--emit_use_strict"},
			{"", "--use_types_for_optimization=false"},
			{"", "--assume_static_inheritance_is_not_used=false"},
			{"", "--process_closure_primitives=false"},
			{"", "-D goog.DEBUG=false", "--define=ns.N=3", "--define=ns.S='s'"},
		})},
	}
	for i := range flagConfigs {
		flagConfigs[i].group = "flags"
		flagConfigs[i].optIn = true
	}
	sweepConfigs = append(sweepConfigs, flagConfigs...)
}

// pairwise returns rows (variants) that cover every pair of values of two
// factors, greedily: each row starts from the first pair not covered yet
// and takes, for every other factor in turn, the value that covers the most
// pairs not covered yet with the values already chosen. Deterministic.
func pairwise(factors [][]string) [][]string {
	type pair struct{ f1, v1, f2, v2 int }
	covered := map[pair]bool{}
	left := 0
	for f1 := range factors {
		for f2 := f1 + 1; f2 < len(factors); f2++ {
			left += len(factors[f1]) * len(factors[f2])
		}
	}
	var rows [][]string
	for left > 0 {
		row := make([]int, len(factors))
		for i := range row {
			row[i] = -1
		}
		// The first pair not covered.
	first:
		for f1 := range factors {
			for f2 := f1 + 1; f2 < len(factors); f2++ {
				for v1 := range factors[f1] {
					for v2 := range factors[f2] {
						if !covered[pair{f1, v1, f2, v2}] {
							row[f1], row[f2] = v1, v2
							break first
						}
					}
				}
			}
		}
		for f := range factors {
			if row[f] >= 0 {
				continue
			}
			best, bestGain := 0, -1
			for v := range factors[f] {
				gain := 0
				for g := range factors {
					if g == f || row[g] < 0 {
						continue
					}
					p := pair{g, row[g], f, v}
					if f < g {
						p = pair{f, v, g, row[g]}
					}
					if !covered[p] {
						gain++
					}
				}
				if gain > bestGain {
					best, bestGain = v, gain
				}
			}
			row[f] = best
		}
		var args []string
		for f1 := range factors {
			args = append(args, split(factors[f1][row[f1]])...)
			for f2 := f1 + 1; f2 < len(factors); f2++ {
				p := pair{f1, row[f1], f2, row[f2]}
				if !covered[p] {
					covered[p] = true
					left--
				}
			}
		}
		rows = append(rows, args)
	}
	return rows
}
