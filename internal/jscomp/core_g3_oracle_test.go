package jscomp

// Go-only: differential test — the Go side of tools/oracle/java/CoreG3Dump.java
// (core group g3: CompilerOptions, CompilerInput, LazyParsedDependencyInfo,
// JSChunk, JSChunkGraph, deps/JsFileRegexParser, deps/JsFileLineParser).
// coreG3Dump reproduces, byte for byte, what the Java tool prints (see its
// header for the format of each mode).
//
//   - TestCoreG3OracleGolden compares against the checked-in Java dumps in
//     testdata/core_g3/: options.txt.gz (seed 0, 300 scenarios),
//     strings.txt.gz (seed 0, 3000 strings), chunks.txt.gz (seed 0, 300
//     graphs), depsgen.txt.gz (seed 0, 1500 generated sources) and
//     fixtures-deps.txt.gz (deps over the repository's testdata/fixtures).
//     Regenerate with
//       java -Xss1g -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.CoreG3Dump options \
//         internal/jscomp/testdata/core_g3/options.txt.gz 0 300
//     (strings ... 0 3000, chunks ... 0 300, depsgen ... 0 1500, and
//     `deps internal/jscomp/testdata/core_g3/fixtures-deps.txt.gz
//     testdata/fixtures --list=.cache/fixtures-files.txt`).
//   - TestCoreG3OracleCorpus compares the deps mode over a corpus when
//     CORE_G3_ROOT (the root of the sources) and CORE_G3_ORACLE (the Java
//     `deps` dump of files under that root; the files compared are the
//     dump's "## file" sections) are set. CORE_G3_EVERY=n compares every
//     n-th file only, and CORE_G3_DIFF_DIR, if set, receives the Go and Java
//     sections of every mismatching file. The Java dump of
//     /home/user/jsgo/.cache/corpus is produced with
//     `deps <out.gz> .cache/corpus --list=.cache/corpus-files.txt`. The
//     held-out corpus (.cache/corpus-holdout, corpus-holdout-files.txt) is
//     for a final check only.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// Shared helpers

// coreG3Escape is CoreG3Dump.esc (ParsingDump.escape).
func coreG3Escape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			fmt.Fprintf(&b, "\\u%04x", c)
		}
	}
	return b.String()
}

type coreG3Out struct {
	sb strings.Builder
}

func (o *coreG3Out) line(s string) {
	o.sb.WriteString(coreG3Escape(s))
	o.sb.WriteByte('\n')
}

var coreG3Identity = regexp.MustCompile(`(com\.google\.[\w.$]+)@[0-9a-f]+`)

// coreG3Mask is CoreG3Dump.mask: Object.toString() identity hashes print as
// "@ID".
func coreG3Mask(s string) string {
	return coreG3Identity.ReplaceAllString(s, "${1}@ID")
}

// coreG3Exc is CoreG3Dump.exc: "! <getSimpleName()>: <getMessage()>".
func coreG3Exc(t jutil.Throwable) string {
	s, msg := t.Error(), t.GetMessage()
	class, message := s, "null"
	if strings.HasSuffix(s, ": "+msg) {
		class, message = s[:len(s)-len(msg)-2], msg
	}
	if i := strings.LastIndexAny(class, ".$"); i >= 0 {
		class = class[i+1:]
	}
	return "! " + class + ": " + message
}

// coreG3Try runs body; a thrown Java exception is returned (nil if none).
func coreG3Try(body func()) (thrown jutil.Throwable) {
	jutil.Try(body, func(e jutil.Throwable) { thrown = e })
	return thrown
}

// coreG3Fmt is CoreG3Dump.fmt: String.valueOf for option values, with
// collections printed element by element and enums by name.
func coreG3Fmt(v any) string {
	if jutil.IsNil(v) {
		return "null"
	}
	switch x := v.(type) {
	case *string:
		return *x
	case string, bool, int, int32, int64, float64, uint16:
		return jutil.ToJavaString(x)
	}
	rv := reflect.ValueOf(v)
	if e, ok := v.(interface{ Name() string }); ok && rv.CanInt() {
		// An enum whose toString() is not its name (e.g. Tri).
		return e.Name()
	}
	typeName := rv.Type().String()
	switch {
	case strings.HasPrefix(typeName, "*jutil.LinkedHashMap["):
		keys := rv.MethodByName("KeySlice").Call(nil)[0]
		get := rv.MethodByName("Get")
		parts := make([]string, keys.Len())
		for i := range parts {
			k := keys.Index(i)
			parts[i] = coreG3Fmt(k.Interface()) + "=" + coreG3Fmt(get.Call([]reflect.Value{k})[0].Interface())
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case strings.HasPrefix(typeName, "*jutil.LinkedHashSet["), strings.HasPrefix(typeName, "*jutil.ArrayList["):
		return coreG3FmtSlice(rv.MethodByName("ToSlice").Call(nil)[0])
	case rv.Kind() == reflect.Slice:
		return coreG3FmtSlice(rv)
	}
	return coreG3Mask(compilerOptions_toJavaString(v))
}

func coreG3FmtSlice(rv reflect.Value) string {
	parts := make([]string, rv.Len())
	for i := range parts {
		parts[i] = coreG3Fmt(rv.Index(i).Interface())
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// coreG3List is CoreG3Dump.list: "" is the empty list, otherwise the
// comma-separated parts (the catalog has no trailing empty parts, which
// String.split would drop).
func coreG3List(a string) []string {
	if a == "" {
		return []string{}
	}
	return strings.Split(a, ",")
}

func coreG3Set(a string) *jutil.LinkedHashSet[string] {
	return jutil.NewLinkedHashSet(coreG3List(a)...)
}

func coreG3Nullable(a string) *string {
	if a == "null" {
		return nil
	}
	return &a
}

func coreG3Names[E fmt.Stringer](values []E) []string {
	names := make([]string, len(values))
	for i, v := range values {
		names[i] = v.String()
	}
	return names
}

// coreG3ValueOf is Enum.valueOf(name).
func coreG3ValueOf[E fmt.Stringer](values []E, name string) E {
	for _, v := range values {
		if v.String() == name {
			return v
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("No enum constant " + name))
}

// ---------------------------------------------------------------------------
// options

type coreG3Op struct {
	name  string
	args  []string
	apply func(o *CompilerOptions, a string)
}

func coreG3BoolOp(name string, apply func(o *CompilerOptions, b bool)) coreG3Op {
	return coreG3Op{name, []string{"true", "false"}, func(o *CompilerOptions, a string) {
		apply(o, strings.EqualFold(a, "true"))
	}}
}

func coreG3EnumOp[E fmt.Stringer](name string, values []E, apply func(o *CompilerOptions, e E)) coreG3Op {
	return coreG3Op{name, coreG3Names(values), func(o *CompilerOptions, a string) {
		apply(o, coreG3ValueOf(values, a))
	}}
}

func coreG3EntryPoints(a string) []*ModuleIdentifier {
	l := []*ModuleIdentifier{}
	for _, s := range coreG3List(a) {
		if strings.HasPrefix(s, "file:") {
			l = append(l, ModuleIdentifier_ForFile(s[5:]))
		} else {
			l = append(l, ModuleIdentifier_ForClosure(s))
		}
	}
	return l
}

func coreG3DepOptions(a string) *DependencyOptions {
	switch a {
	case "none":
		return DependencyOptions_None()
	case "sortOnly":
		return DependencyOptions_SortOnly()
	}
	kind, rest, _ := strings.Cut(a, ":")
	eps := coreG3EntryPoints(rest)
	switch kind {
	case "pruneLegacy":
		return DependencyOptions_PruneLegacyForEntryPoints(eps)
	case "prune":
		return DependencyOptions_PruneForEntryPoints(eps)
	case "pruneAllow":
		return DependencyOptions_PruneAllowNoEntryPoints(eps)
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(a))
}

func coreG3Atoi(a string) int {
	n, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return n
}

// coreG3Ops is CoreG3Dump.OPS, in the same order.
func coreG3Ops() []coreG3Op {
	lang := CompilerOptions_LanguageMode_values()
	reach := CompilerOptions_Reach_values()
	ops := []coreG3Op{
		coreG3EnumOp("setLanguageIn", lang, (*CompilerOptions).SetLanguageIn),
		coreG3EnumOp("setLanguageOut", lang, (*CompilerOptions).SetLanguageOut),
		coreG3EnumOp("setLanguage", lang, (*CompilerOptions).SetLanguage),
		{"setBrowserFeaturesetYear",
			[]string{"2012", "2017", "2018", "2019", "2020", "2021", "2022", "2023", "2024", "2025", "2026", "2027"},
			func(o *CompilerOptions, a string) { o.SetBrowserFeaturesetYear(coreG3Atoi(a)) }},
		coreG3BoolOp("setCheckTypes", (*CompilerOptions).SetCheckTypes),
		coreG3BoolOp("setCheckSymbols", (*CompilerOptions).SetCheckSymbols),
		coreG3BoolOp("setCheckSuspiciousCode", (*CompilerOptions).SetCheckSuspiciousCode),
		coreG3BoolOp("setFoldConstants", (*CompilerOptions).SetFoldConstants),
		coreG3BoolOp("setDeadAssignmentElimination", (*CompilerOptions).SetDeadAssignmentElimination),
		coreG3BoolOp("setDeadPropertyAssignmentElimination", (*CompilerOptions).SetDeadPropertyAssignmentElimination),
		coreG3BoolOp("setInlineConstantVars", (*CompilerOptions).SetInlineConstantVars),
		coreG3BoolOp("setCrossChunkCodeMotion", (*CompilerOptions).SetCrossChunkCodeMotion),
		coreG3BoolOp("setCrossChunkMethodMotion", (*CompilerOptions).SetCrossChunkMethodMotion),
		coreG3BoolOp("setCoalesceVariableNames", (*CompilerOptions).SetCoalesceVariableNames),
		coreG3BoolOp("setInlineVariables", (*CompilerOptions).SetInlineVariables),
		coreG3EnumOp("setInlineVariablesOfReach", reach, (*CompilerOptions).SetInlineVariablesOfReach),
		coreG3EnumOp("setRemoveUnusedVariables", reach, (*CompilerOptions).SetRemoveUnusedVariables),
		coreG3EnumOp("setInlineFunctions", reach, (*CompilerOptions).SetInlineFunctions),
		coreG3BoolOp("setInlineProperties", (*CompilerOptions).SetInlineProperties),
		coreG3BoolOp("setSmartNameRemoval", (*CompilerOptions).SetSmartNameRemoval),
		coreG3BoolOp("setExtractPrototypeMemberDeclarations", (*CompilerOptions).SetExtractPrototypeMemberDeclarations),
		coreG3EnumOp("setExtractPrototypeMemberDeclarationsOfMode",
			CompilerOptions_ExtractPrototypeMemberDeclarationsMode_values(),
			(*CompilerOptions).SetExtractPrototypeMemberDeclarationsOfMode),
		coreG3BoolOp("setRemoveUnusedPrototypeProperties", (*CompilerOptions).SetRemoveUnusedPrototypeProperties),
		coreG3BoolOp("setRemoveUnusedClassProperties", (*CompilerOptions).SetRemoveUnusedClassProperties),
		coreG3BoolOp("setCollapseProperties", (*CompilerOptions).SetCollapseProperties),
		coreG3EnumOp("setCollapsePropertiesLevel", CompilerOptions_PropertyCollapseLevel_values(),
			(*CompilerOptions).SetCollapsePropertiesLevel),
		coreG3BoolOp("setDisambiguateProperties", (*CompilerOptions).SetDisambiguateProperties),
		coreG3BoolOp("setAmbiguateProperties", (*CompilerOptions).SetAmbiguateProperties),
		coreG3BoolOp("setUseTypesForOptimization", (*CompilerOptions).SetUseTypesForOptimization),
		coreG3BoolOp("setUseTypesForLocalOptimization", (*CompilerOptions).SetUseTypesForLocalOptimization),
		{"setPolymerVersion", []string{"null", "1", "2", "3"}, func(o *CompilerOptions, a string) {
			if a == "null" {
				o.SetPolymerVersion(nil)
			} else {
				v := int32(coreG3Atoi(a))
				o.SetPolymerVersion(&v)
			}
		}},
		coreG3BoolOp("setAngularPass", (*CompilerOptions).SetAngularPass),
		coreG3BoolOp("setChromePass", (*CompilerOptions).SetChromePass),
		coreG3BoolOp("setClosurePass", (*CompilerOptions).SetClosurePass),
		coreG3BoolOp("setPreserveClosurePrimitives", (*CompilerOptions).SetPreserveClosurePrimitives),
		coreG3EnumOp("setIncrementalChecks", CompilerOptions_IncrementalCheckMode_values(),
			(*CompilerOptions).SetIncrementalChecks),
		coreG3BoolOp("setEmitUseStrict", func(o *CompilerOptions, b bool) { o.SetEmitUseStrict(b) }),
		coreG3BoolOp("setStrictModeInput", func(o *CompilerOptions, b bool) { o.SetStrictModeInput(b) }),
		coreG3BoolOp("setIsolatePolyfills", (*CompilerOptions).SetIsolatePolyfills),
		coreG3BoolOp("setInstrumentAsyncContext", (*CompilerOptions).SetInstrumentAsyncContext),
		coreG3BoolOp("setRewritePolyfills", (*CompilerOptions).SetRewritePolyfills),
		coreG3BoolOp("setPreventLibraryInjection", (*CompilerOptions).SetPreventLibraryInjection),
		coreG3EnumOp("setRuntimeLibraryMode", RuntimeJsLibManager_RuntimeLibraryMode_values(),
			(*CompilerOptions).SetRuntimeLibraryMode),
		{"setDefineToBooleanLiteral", []string{"A:true", "B:false", "goog.DEBUG:true"}, func(o *CompilerOptions, a string) {
			k, v, _ := strings.Cut(a, ":")
			o.SetDefineToBooleanLiteral(k, strings.EqualFold(v, "true"))
		}},
		{"setDefineToNumberLiteral", []string{"A:3", "C:-7", "goog.FEATURESET_YEAR:2000"}, func(o *CompilerOptions, a string) {
			k, v, _ := strings.Cut(a, ":")
			o.SetDefineToNumberLiteral(k, coreG3Atoi(v))
		}},
		{"setDefineToDoubleLiteral", []string{"A:0.5", "D:1e21", "E:-0.0"}, func(o *CompilerOptions, a string) {
			k, v, _ := strings.Cut(a, ":")
			d, err := strconv.ParseFloat(v, 64)
			if err != nil {
				panic(err)
			}
			o.SetDefineToDoubleLiteral(k, d)
		}},
		{"setDefineToStringLiteral", []string{"A:x", "F:hello world", "G:\u00e9"}, func(o *CompilerOptions, a string) {
			k, v, _ := strings.Cut(a, ":")
			o.SetDefineToStringLiteral(k, v)
		}},
		{"setRenamePrefix", []string{"null", "", "$p_"}, func(o *CompilerOptions, a string) { o.SetRenamePrefix(coreG3Nullable(a)) }},
		{"setRenamePrefixNamespace", []string{"null", "_ns"}, func(o *CompilerOptions, a string) {
			o.SetRenamePrefixNamespace(coreG3Nullable(a))
		}},
		{"setLocale", []string{"null", "en", "fr_CA"}, func(o *CompilerOptions, a string) { o.SetLocale(coreG3Nullable(a)) }},
		{"setSourceMapOutputPath", []string{"null", "", "out.map"}, func(o *CompilerOptions, a string) {
			o.SetSourceMapOutputPath(coreG3Nullable(a))
		}},
		coreG3BoolOp("setAlwaysGatherSourceMapInfo", (*CompilerOptions).SetAlwaysGatherSourceMapInfo),
		coreG3BoolOp("setSkipNonTranspilationPasses", (*CompilerOptions).SetSkipNonTranspilationPasses),
		coreG3BoolOp("setChecksOnly", (*CompilerOptions).SetChecksOnly),
		coreG3EnumOp("setOutputJs", CompilerOptions_OutputJs_values(), (*CompilerOptions).SetOutputJs),
		coreG3BoolOp("setGenerateExports", (*CompilerOptions).SetGenerateExports),
		coreG3BoolOp("setPrettyPrint", (*CompilerOptions).SetPrettyPrint),
		coreG3BoolOp("setLineBreak", (*CompilerOptions).SetLineBreak),
		coreG3BoolOp("setPrintInputDelimiter", (*CompilerOptions).SetPrintInputDelimiter),
		{"setInputDelimiter", []string{"", "// file %name%"}, (*CompilerOptions).SetInputDelimiter},
		coreG3EnumOp("setTracerMode", CompilerOptions_TracerMode_values(), (*CompilerOptions).SetTracerMode),
		coreG3EnumOp("setErrorFormat", ErrorFormat_values(), (*CompilerOptions).SetErrorFormat),
		{"setLineLengthThreshold", []string{"0", "80", "-1"}, func(o *CompilerOptions, a string) {
			o.SetLineLengthThreshold(coreG3Atoi(a))
		}},
		coreG3EnumOp("setVariableRenaming", VariableRenamingPolicy_values(), (*CompilerOptions).SetVariableRenaming),
		coreG3EnumOp("setPropertyRenaming", PropertyRenamingPolicy_values(), (*CompilerOptions).SetPropertyRenaming),
		{"setModuleRoots", []string{"", "a", "a,b/c"}, func(o *CompilerOptions, a string) { o.SetModuleRoots(coreG3List(a)) }},
		{"setCommonJSModulePathPrefix", []string{"lib/"}, (*CompilerOptions).SetCommonJSModulePathPrefix},
		coreG3BoolOp("setProcessCommonJSModules", (*CompilerOptions).SetProcessCommonJSModules),
		coreG3BoolOp("setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf",
			(*CompilerOptions).SetBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf),
		coreG3BoolOp("setEnableModuleRewriting", (*CompilerOptions).SetEnableModuleRewriting),
		coreG3EnumOp("setModuleResolutionMode", ModuleLoader_ResolutionMode_values(),
			(*CompilerOptions).SetModuleResolutionMode),
		coreG3EnumOp("setPathEscaper", ModuleLoader_PathEscaper_values(), (*CompilerOptions).SetPathEscaper),
		{"setPackageJsonEntryNames", []string{"", "main"}, func(o *CompilerOptions, a string) {
			o.SetPackageJsonEntryNames(coreG3List(a))
		}},
		{"setDependencyOptions", []string{"none", "sortOnly", "pruneLegacy:a", "prune:a,file:b.js", "pruneAllow:"},
			func(o *CompilerOptions, a string) { o.SetDependencyOptions(coreG3DepOptions(a)) }},
		coreG3EnumOp("setWarningLevel", CheckLevel_values(), func(o *CompilerOptions, l CheckLevel) {
			o.SetWarningLevel(DiagnosticGroup_ForType(RhinoErrorReporter_TYPE_PARSE_ERROR), l)
		}),
		{"resetWarningsGuard", []string{""}, func(o *CompilerOptions, a string) { o.ResetWarningsGuard() }},
		{"setIdGenerators", []string{"", "goog.getId", "a,b"}, func(o *CompilerOptions, a string) { o.SetIdGenerators(coreG3Set(a)) }},
		{"setStripTypes", []string{"", "goog.debug"}, func(o *CompilerOptions, a string) { o.SetStripTypes(coreG3Set(a)) }},
		{"setStripNameSuffixes", []string{"", "logger,Logger"}, func(o *CompilerOptions, a string) {
			o.SetStripNameSuffixes(coreG3Set(a))
		}},
		{"setStripNamePrefixes", []string{"", "x."}, func(o *CompilerOptions, a string) { o.SetStripNamePrefixes(coreG3Set(a)) }},
		{"setReplaceStringsConfiguration", []string{"`:", "#:f(?)"}, func(o *CompilerOptions, a string) {
			k, v, _ := strings.Cut(a, ":")
			o.SetReplaceStringsConfiguration(k, coreG3List(v))
		}},
		{"setPropertiesThatMustDisambiguate", []string{"", "p,q"}, func(o *CompilerOptions, a string) {
			o.SetPropertiesThatMustDisambiguate(coreG3Set(a))
		}},
		{"setExtraAnnotationNames", []string{"", "foo,bar"}, func(o *CompilerOptions, a string) {
			o.SetExtraAnnotationNames(coreG3List(a))
		}},
		{"setForceLibraryInjection", []string{"", "es6_runtime"}, func(o *CompilerOptions, a string) {
			o.SetForceLibraryInjection(coreG3List(a))
		}},
		{"setUnknownDefinesToIgnore", []string{"", "X,Y"}, func(o *CompilerOptions, a string) {
			o.SetUnknownDefinesToIgnore(coreG3List(a))
		}},
		{"addUnknownDefinesToIgnore", []string{"Z"}, func(o *CompilerOptions, a string) {
			o.AddUnknownDefinesToIgnore(coreG3List(a))
		}},
		{"setReplaceMessagesWithChromeI18n", []string{"true:null", "true:proj", "false:null"}, func(o *CompilerOptions, a string) {
			k, v, _ := strings.Cut(a, ":")
			o.SetReplaceMessagesWithChromeI18n(strings.EqualFold(k, "true"), coreG3Nullable(v))
		}},
		coreG3BoolOp("setDoLateLocalization", (*CompilerOptions).SetDoLateLocalization),
		{"setMessageBundle", []string{"null", "empty"}, func(o *CompilerOptions, a string) {
			if a == "null" {
				o.SetMessageBundle(nil)
			} else {
				o.SetMessageBundle(NewEmptyMessageBundle())
			}
		}},
		coreG3BoolOp("setValidateRequiredInlinings", (*CompilerOptions).SetValidateRequiredInlinings),
		{"setMaxFunctionSizeAfterInlining", []string{"0", "100"}, func(o *CompilerOptions, a string) {
			o.SetMaxFunctionSizeAfterInlining(coreG3Atoi(a))
		}},
		coreG3EnumOp("setInstrumentForCoverageOption", CompilerOptions_InstrumentOption_values(),
			(*CompilerOptions).SetInstrumentForCoverageOption),
		coreG3BoolOp("setInstrumentForCoverageOnly", (*CompilerOptions).SetInstrumentForCoverageOnly),
		{"setProductionInstrumentationArrayName", []string{"", "arr"}, (*CompilerOptions).SetProductionInstrumentationArrayName},
		{"setDebugLogFilter", []string{"null", "a,,b"}, func(o *CompilerOptions, a string) { o.SetDebugLogFilter(coreG3Nullable(a)) }},
		coreG3BoolOp("setTrustedStrings", (*CompilerOptions).SetTrustedStrings),
		coreG3BoolOp("setPreferSingleQuotes", (*CompilerOptions).SetPreferSingleQuotes),
		coreG3BoolOp("setQuoteKeywordProperties", (*CompilerOptions).SetQuoteKeywordProperties),
		coreG3BoolOp("setAssumeStrictThis", (*CompilerOptions).SetAssumeStrictThis),
		coreG3BoolOp("setProtectHiddenSideEffects", (*CompilerOptions).SetProtectHiddenSideEffects),
		coreG3BoolOp("setAssumeGettersArePure", (*CompilerOptions).SetAssumeGettersArePure),
		coreG3EnumOp("setDevMode", CompilerOptions_DevMode_values(), (*CompilerOptions).SetDevMode),
		coreG3EnumOp("setEnvironment", CompilerOptions_Environment_values(), (*CompilerOptions).SetEnvironment),
		coreG3EnumOp("setTweakProcessing", CompilerOptions_TweakProcessing_values(), (*CompilerOptions).SetTweakProcessing),
		coreG3EnumOp("setChunkOutputType", CompilerOptions_ChunkOutputType_values(), (*CompilerOptions).SetChunkOutputType),
		coreG3EnumOp("setEs6ModuleTranspilation", CompilerOptions_Es6ModuleTranspilation_values(),
			(*CompilerOptions).SetEs6ModuleTranspilation),
		coreG3EnumOp("setEs6SubclassTranspilation", CompilerOptions_Es6SubclassTranspilation_values(),
			(*CompilerOptions).SetEs6SubclassTranspilation),
		coreG3EnumOp("setConformanceReportingMode", CompilerOptions_ConformanceReportingMode_values(),
			(*CompilerOptions).SetConformanceReportingMode),
		coreG3BoolOp("setContinueAfterErrors", (*CompilerOptions).SetContinueAfterErrors),
		coreG3EnumOp("setParseJsDocDocumentation", parsing.Config_JsDocParsing_values(),
			(*CompilerOptions).SetParseJsDocDocumentation),
		coreG3BoolOp("setPreserveDetailedSourceInfo", (*CompilerOptions).SetPreserveDetailedSourceInfo),
		coreG3EnumOp("setSourceMapDetailLevel", SourceMap_DetailLevel_values(), (*CompilerOptions).SetSourceMapDetailLevel),
		coreG3EnumOp("setSourceMapFormat", SourceMap_Format_values(), (*CompilerOptions).SetSourceMapFormat),
		{"setSourceMapLocationMappings", []string{"", "a/|b/", "a/|,c|d"}, func(o *CompilerOptions, a string) {
			l := []SourceMap_LocationMapping{}
			for _, m := range coreG3List(a) {
				parts := strings.Split(m, "|")
				l = append(l, NewSourceMap_PrefixLocationMapping(parts[0], parts[1]))
			}
			o.SetSourceMapLocationMappings(l)
		}},
		{"setDynamicImportAlias", []string{"null", "imp"}, func(o *CompilerOptions, a string) {
			o.SetDynamicImportAlias(coreG3Nullable(a))
		}},
		coreG3BoolOp("setAllowDynamicImport", (*CompilerOptions).SetAllowDynamicImport),
		{"setExperimentalOutputFeatureSet",
			[]string{"null", "ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS", "BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD"},
			func(o *CompilerOptions, a string) {
				if a == "null" {
					o.SetExperimentalOutputFeatureSet(nil)
				} else {
					v := coreG3ValueOf(CompilerOptions_ExperimentalOutputFeatureSet_values(), a)
					o.SetExperimentalOutputFeatureSet(&v)
				}
			}},
		coreG3BoolOp("setTranspilePublicClassFields", (*CompilerOptions).SetTranspilePublicClassFields),
		coreG3BoolOp("setOptimizeLetAndConst", (*CompilerOptions).SetOptimizeLetAndConst),
		coreG3BoolOp("setAssumeGlobalScopeIsIsolated", (*CompilerOptions).SetAssumeGlobalScopeIsIsolated),
		coreG3EnumOp("setOptimizeLocalAccessForGlobalSymbolNamespace", CompilerOptions_OptimizeLocalAccess_values(),
			(*CompilerOptions).SetOptimizeLocalAccessForGlobalSymbolNamespace),
		{"setCssRenamingSkiplist", []string{"", "x"}, func(o *CompilerOptions, a string) { o.SetCssRenamingSkiplist(coreG3Set(a)) }},
		{"setInjectPolyfillsNewerThan", []string{"null", "ECMASCRIPT5", "ECMASCRIPT_2015"}, func(o *CompilerOptions, a string) {
			if a == "null" {
				o.SetInjectPolyfillsNewerThan(nil)
			} else {
				o.SetInjectPolyfillsNewerThan(coreG3ValueOf(lang, a).Ptr())
			}
		}},
		coreG3BoolOp("setRewriteGlobalDeclarationsForTryCatchWrapping",
			(*CompilerOptions).SetRewriteGlobalDeclarationsForTryCatchWrapping),
		coreG3BoolOp("setGatherCssNames", (*CompilerOptions).SetGatherCssNames),
		coreG3BoolOp("setMarkAsCompiled", (*CompilerOptions).SetMarkAsCompiled),
		coreG3BoolOp("setLabelRenaming", (*CompilerOptions).SetLabelRenaming),
		coreG3BoolOp("setGeneratePseudoNames", (*CompilerOptions).SetGeneratePseudoNames),
		coreG3BoolOp("setDevirtualizeMethods", (*CompilerOptions).SetDevirtualizeMethods),
		coreG3BoolOp("setComputeFunctionSideEffects", (*CompilerOptions).SetComputeFunctionSideEffects),
		coreG3BoolOp("setOptimizeCalls", (*CompilerOptions).SetOptimizeCalls),
		coreG3BoolOp("setConvertToDottedProperties", (*CompilerOptions).SetConvertToDottedProperties),
		coreG3EnumOp("setAliasStringsMode", CompilerOptions_AliasStringsMode_values(), (*CompilerOptions).SetAliasStringsMode),
		coreG3EnumOp("setJ2clPass", CompilerOptions_J2clPassMode_values(), (*CompilerOptions).SetJ2clPass),
		{"setSummaryDetailLevel", []string{"0", "3"}, func(o *CompilerOptions, a string) { o.SetSummaryDetailLevel(coreG3Atoi(a)) }},
		{"setNumParallelThreads", []string{"4"}, func(o *CompilerOptions, a string) { o.SetNumParallelThreads(coreG3Atoi(a)) }},
		{"setConformanceRemoveRegexFromPath", []string{"absent", "^x/"}, func(o *CompilerOptions, a string) {
			if a == "absent" {
				o.SetConformanceRemoveRegexFromPath(nil)
			} else {
				o.SetConformanceRemoveRegexFromPath(jutil.Pattern_Compile(a))
			}
		}},
	}
	return ops
}

// coreG3EnumLists are the enums declared in CompilerOptions, by simple name
// (CoreG3Dump lists CompilerOptions.class.getDeclaredClasses() sorted by
// simple name).
func coreG3EnumLists() map[string][]string {
	return map[string][]string{
		"AliasStringsMode":                       coreG3Names(CompilerOptions_AliasStringsMode_values()),
		"BrowserFeaturesetYear":                  coreG3Names(CompilerOptions_BrowserFeaturesetYear_values()),
		"ChunkOutputType":                        coreG3Names(CompilerOptions_ChunkOutputType_values()),
		"ConformanceReportingMode":               coreG3Names(CompilerOptions_ConformanceReportingMode_values()),
		"DevMode":                                coreG3Names(CompilerOptions_DevMode_values()),
		"Environment":                            coreG3Names(CompilerOptions_Environment_values()),
		"Es6ModuleTranspilation":                 coreG3Names(CompilerOptions_Es6ModuleTranspilation_values()),
		"Es6SubclassTranspilation":               coreG3Names(CompilerOptions_Es6SubclassTranspilation_values()),
		"ExperimentalOutputFeatureSet":           coreG3Names(CompilerOptions_ExperimentalOutputFeatureSet_values()),
		"ExtractPrototypeMemberDeclarationsMode": coreG3Names(CompilerOptions_ExtractPrototypeMemberDeclarationsMode_values()),
		"IncrementalCheckMode":                   coreG3Names(CompilerOptions_IncrementalCheckMode_values()),
		"InstrumentOption":                       coreG3Names(CompilerOptions_InstrumentOption_values()),
		"IsolationMode":                          coreG3Names(CompilerOptions_IsolationMode_values()),
		"J2clPassMode":                           coreG3Names(CompilerOptions_J2clPassMode_values()),
		"JsonStreamMode":                         coreG3Names(CompilerOptions_JsonStreamMode_values()),
		"LanguageMode":                           coreG3Names(CompilerOptions_LanguageMode_values()),
		"OptimizeLocalAccess":                    coreG3Names(CompilerOptions_OptimizeLocalAccess_values()),
		"OutputJs":                               coreG3Names(CompilerOptions_OutputJs_values()),
		"PropertyCollapseLevel":                  coreG3Names(CompilerOptions_PropertyCollapseLevel_values()),
		"Reach":                                  coreG3Names(CompilerOptions_Reach_values()),
		"SegmentOfCompilationToRun":              coreG3Names(CompilerOptions_SegmentOfCompilationToRun_values()),
		"TracerMode":                             coreG3Names(CompilerOptions_TracerMode_values()),
		"TweakProcessing":                        coreG3Names(CompilerOptions_TweakProcessing_values()),
	}
}

// coreG3Getter is one zero-argument CompilerOptions method that returns a
// value; javaName is the Java method name.
type coreG3Getter struct {
	javaName string
	method   reflect.Method
}

// coreG3Getters are the Go counterparts of CompilerOptions' declared
// zero-argument instance methods returning a value (other than toString and
// hashCode), sorted by Java name.
func coreG3Getters() []coreG3Getter {
	var getters []coreG3Getter
	t := reflect.TypeOf((*CompilerOptions)(nil))
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.Type.NumIn() != 1 || m.Type.NumOut() != 1 || m.Name == "String" || m.Name == "HashCode" {
			continue
		}
		getters = append(getters, coreG3Getter{strings.ToLower(m.Name[:1]) + m.Name[1:], m})
	}
	sort.Slice(getters, func(i, j int) bool { return getters[i].javaName < getters[j].javaName })
	return getters
}

func (o *coreG3Out) options(seed int64, n int) {
	o.line("## statics")
	for _, m := range CompilerOptions_LanguageMode_values() {
		var fs string
		if t := coreG3Try(func() { fs = m.ToFeatureSet().String() }); t != nil {
			fs = coreG3Exc(t)
		}
		o.line("LanguageMode " + m.String() + " " + strconv.Itoa(m.Ordinal()) + " strict=" +
			strconv.FormatBool(m.IsDefaultStrict()) + " " + fs)
	}
	o.line("STABLE_IN " + CompilerOptions_LanguageMode_STABLE_IN.String() +
		" STABLE_OUT " + CompilerOptions_LanguageMode_STABLE_OUT.String())
	o.line("validCommandLineNames " + jutil.CollectionToString(CompilerOptions_LanguageMode_ValidCommandLineNames()))
	for _, s := range []string{
		"ECMASCRIPT3", "  es3  ", "junk", "ECMASCRIPT_2020", "  es_2020  ", "  es2020  ", "es6",
		"ES6_STRICT", "ecmascript_next", "STABLE", "Es5", "ESNEXT", "ES_NEXT", "", "estable",
	} {
		o.line("fromString " + s + " -> " + coreG3Fmt(CompilerOptions_LanguageMode_FromString(&s)))
	}
	for _, s := range []string{"NONE", "LINE", "BRANCH", "PRODUCTION", "LINE_ONLY", "x"} {
		o.line("InstrumentOption.fromString " + s + " -> " + coreG3Fmt(CompilerOptions_InstrumentOption_FromString(&s)))
	}
	for _, y := range CompilerOptions_BrowserFeaturesetYear_values() {
		o.line("BrowserFeaturesetYear " + y.String() + " " + strconv.Itoa(y.GetYear()) + " " + y.GetFeatureSet().String())
	}
	for _, f := range parser.FeatureSet_Feature_values() {
		o.line("minimumRequiredFor " + f.Name() + " " + coreG3Fmt(CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(f)))
	}
	for _, r := range CompilerOptions_Reach_values() {
		o.line("Reach " + r.String() + " " + strconv.FormatBool(r.IsOn()) + " " + strconv.FormatBool(r.IncludesGlobals()))
	}
	for _, r := range CompilerOptions_TracerMode_values() {
		o.line("TracerMode " + r.String() + " " + strconv.FormatBool(r.IsOn()) + " " + strconv.FormatBool(r.DoPruningAnalysis()))
	}
	for _, r := range CompilerOptions_TweakProcessing_values() {
		o.line("TweakProcessing " + r.String() + " " + strconv.FormatBool(r.IsOn()) + " " + strconv.FormatBool(r.ShouldStrip()))
	}
	for _, r := range CompilerOptions_J2clPassMode_values() {
		o.line("J2clPassMode " + r.String() + " " + strconv.FormatBool(r.ShouldAddJ2clPasses()))
	}
	enums := coreG3EnumLists()
	var enumNames []string
	for name := range enums {
		enumNames = append(enumNames, name)
	}
	sort.Strings(enumNames)
	for _, name := range enumNames {
		o.line("enum " + name + " " + jutil.CollectionToString(enums[name]))
	}

	ops := coreG3Ops()
	getters := coreG3Getters()
	for i := 0; i < n; i++ {
		o.line("## options " + strconv.Itoa(i))
		rnd := jutil.NewRandomWithSeed(seed + int64(i))
		opts := NewCompilerOptions()
		count := 0
		if i != 0 {
			bound := int32(25)
			if i < 50 {
				bound = 4
			}
			count = 1 + int(rnd.NextIntBound(bound))
		}
		for k := 0; k < count; k++ {
			op := ops[rnd.NextIntBound(int32(len(ops)))]
			arg := op.args[rnd.NextIntBound(int32(len(op.args)))]
			o.line("op " + op.name + " " + arg)
			if t := coreG3TryRuntime(func() { op.apply(opts, arg) }); t != nil {
				o.line(coreG3Exc(t))
			} else {
				o.line("ok")
			}
		}
		var s string
		if t := coreG3TryRuntime(func() { s = coreG3Mask(opts.String()) }); t != nil {
			s = coreG3Exc(t)
		}
		o.line("toString " + s)
		for _, g := range getters {
			var v string
			if t := coreG3Try(func() {
				r := g.method.Func.Call([]reflect.Value{reflect.ValueOf(opts)})[0].Interface()
				if g.javaName == "getExperimentalOutputFeatureSet" {
					// Go-only: the Java getter returns a Guava Optional; the port
					// returns a nil-able pointer.
					if jutil.IsNil(r) {
						v = "Optional.absent()"
					} else {
						v = "Optional.of(" + coreG3Fmt(reflect.ValueOf(r).Elem().Interface()) + ")"
					}
					return
				}
				if g.javaName == "getConformanceRemoveRegexFromPath" {
					// Go-only: the Java getter returns a Guava Optional<Pattern>;
					// the port returns a nil-able *jutil.Pattern.
					if jutil.IsNil(r) {
						v = "Optional.absent()"
					} else {
						v = "Optional.of(" + coreG3Fmt(r) + ")"
					}
					return
				}
				v = coreG3Fmt(r)
			}); t != nil {
				v = coreG3Exc(t)
			}
			o.line("get " + g.javaName + " " + v)
		}
	}
}

// coreG3TryRuntime is coreG3Try for `catch (RuntimeException e)`.
func coreG3TryRuntime(body func()) (thrown jutil.Throwable) {
	jutil.Try(body, func(e jutil.IRuntimeException) { thrown = e })
	return thrown
}

// ---------------------------------------------------------------------------
// deps

// coreG3Collector is CoreG3Dump.Collector: an ErrorManager that records every
// report.
type coreG3Collector struct {
	BasicErrorManagerBase
	reports []string
}

func newCoreG3Collector() *coreG3Collector {
	c := &coreG3Collector{}
	c.InitBasicErrorManagerBase(c)
	return c
}

func (c *coreG3Collector) Report(level *CheckLevel, error *JSError) {
	c.reports = append(c.reports, level.String()+" "+error.String())
	c.BasicErrorManagerBase.Report(level, error)
}

func (c *coreG3Collector) Println(level *CheckLevel, error *JSError) {}

func (c *coreG3Collector) PrintSummary() {}

// coreG3DepInfo is CoreG3Dump.depInfo for a CompilerInput.
func coreG3DepInfo(d *CompilerInput) string {
	return "provides=" + jutil.CollectionToString(d.GetProvides()) +
		" requires=" + jutil.CollectionToString(d.GetRequires()) +
		" typeRequires=" + jutil.CollectionToString(d.GetTypeRequires()) +
		" loadFlags=" + d.GetLoadFlags().String() +
		" externs=" + strconv.FormatBool(d.GetHasExternsAnnotation()) +
		" nocompile=" + strconv.FormatBool(d.GetHasNoCompileAnnotation()) +
		" es6=" + strconv.FormatBool(d.IsEs6Module()) +
		" goog=" + strconv.FormatBool(d.IsGoogModule())
}

// depsOf is CoreG3Dump.depsOf.
func (o *coreG3Out) depsOf(rel string, code string) {
	for mode := 0; mode < 2; mode++ {
		em := newCoreG3Collector()
		p := NewJsFileRegexParser(em)
		if mode == 1 {
			p.SetShortcutMode(true)
			p.SetIncludeGoogBase(true)
		}
		tag := "regex "
		if mode == 1 {
			tag = "regexgoog "
		}
		var d DependencyInfo
		if t := coreG3TryRuntime(func() { d = p.ParseFile(rel, rel, code) }); t != nil {
			o.line(tag + coreG3Exc(t))
		} else {
			o.line(tag + jutil.ToJavaString(d) + " succeeded=" + strconv.FormatBool(p.DidParseSucceed()))
		}
		for _, r := range em.reports {
			o.line(" report " + r)
		}
	}
	for mode := 0; mode < 2; mode++ {
		em := newCoreG3Collector()
		compiler := NewCompiler()
		compiler.SetErrorManager(em)
		compiler.InitOptions(NewCompilerOptions())
		if mode == 1 {
			compiler.SetPreferRegexParser(true)
		}
		tag := "ast "
		if mode == 1 {
			tag = "lazy "
		}
		input := NewCompilerInput(SourceFile_FromCode(rel, code))
		input.SetCompiler(compiler)
		if t := coreG3TryRuntime(func() {
			o.line(tag + coreG3DepInfo(input))
			o.line(tag + "known provides=" + jutil.CollectionToString(input.GetKnownProvides()) +
				" requires=" + jutil.CollectionToString(input.GetKnownRequires()))
		}); t != nil {
			o.line(tag + coreG3Exc(t))
		}
		for _, r := range em.reports {
			if strings.HasPrefix(r, "ERROR ") {
				o.line(" report " + r)
			}
		}
	}
}

// coreG3GenLines is CoreG3Dump.GEN_LINES.
var coreG3GenLines = []string{
	"goog.provide('a.b');",
	"goog.provide(\"x\");",
	"goog.provide('a.b');goog.provide('c');",
	"goog.require('a.b');",
	"const y = goog.require('m.y');",
	"let {A, B} = goog.require('ab');",
	"var z=goog.require( 'z' );",
	"goog.requireType('t');",
	"const T = goog.requireType(\"t2\");",
	"goog.forwardDeclare('f');",
	"goog.module('m.z');",
	"goog.module(\"m.q\");",
	"goog.module.declareLegacyNamespace();",
	"goog.declareModuleId('es.m');",
	"goog.setTestOnly();",
	"goog.addDependency('x.js', ['a'], ['b']);",
	"goog.require('a.b')",
	"  goog.require('f');  // trailing",
	"goog.require('g'); goog.require('h');",
	"goog .require('sp');",
	"goog.require('a', 'b');",
	"goog.require(x);",
	"goog.require('');",
	"goog.require('a.b'); goog.provide('a.b');",
	"import {a} from './dep.js';",
	"import * as b from 'goog:a.b';",
	"import 'side.js';",
	"import x from \"./y.js\";",
	"import {c as d} from '../up/c.js'",
	"import\t'tab.js';",
	"export const z = 1;",
	"export {a};",
	"export * from './x.js';",
	"export default 1;",
	"exported = 1;",
	"importx = 1;",
	"import(`dyn.js`);",
	"export\u00e9 = 1;",
	"/** @fileoverview\n * @externs\n */",
	"/**\n * @fileoverview\n * @nocompile\n */",
	"/** @fileoverview\n * @provideGoog */",
	"/** @fileoverview @externs @nocompile */",
	"/** @externs */",
	"/* goog.provide('c'); */",
	"/* start\nend */ goog.provide('i');",
	"// goog.require('d');",
	"var s = 'goog.require(\"e\")';",
	"var t = \"/* not a comment\";",
	"var u = '/*'; goog.require('after.q');",
	"var v = `goog.require('tpl')\n`;",
	"var w = 'a\\\ngoog.require(\"cont\")';",
	"function f() {}",
	"var COMPILED = false;",
	"var goog = goog || {};",
	"",
	"   ",
	"'use strict';",
	"\"use strict\";",
	"goog.loadModule(function(exports) {\n  goog.module('lm');\n  return exports;\n});",
	"if (x) {\n  goog.require('in.block');\n}",
	"class C { m() { return goog.require; } }",
}

// coreG3GenBad is CoreG3Dump.GEN_BAD.
var coreG3GenBad = []string{
	"let 1x = ;", "/* start", " * @externs", "}", "if (x) {", "export\u0301 = 1;",
}

// coreG3GenCode is CoreG3Dump.genCode.
func coreG3GenCode(rnd *jutil.Random) string {
	n := int(rnd.NextIntBound(12))
	var sb strings.Builder
	eol := "\n"
	if rnd.NextIntBound(8) == 0 {
		eol = "\r\n"
	}
	for k := 0; k < n; k++ {
		var l string
		if rnd.NextIntBound(25) == 0 {
			l = coreG3GenBad[rnd.NextIntBound(int32(len(coreG3GenBad)))]
		} else {
			l = coreG3GenLines[rnd.NextIntBound(int32(len(coreG3GenLines)))]
		}
		sb.WriteString(strings.ReplaceAll(l, "\n", eol))
		if k+1 < n || rnd.NextBoolean() {
			sb.WriteString(eol)
		}
	}
	return sb.String()
}

func (o *coreG3Out) depsgen(seed int64, n int) {
	rnd := jutil.NewRandomWithSeed(seed)
	for i := 0; i < n; i++ {
		rel := "gen" + strconv.Itoa(i) + ".js"
		code := coreG3GenCode(rnd)
		o.line("## gen " + rel + " " + strings.ReplaceAll(code, "\n", "|"))
		o.depsOf(rel, code)
	}
}

// ---------------------------------------------------------------------------
// strings

func (o *coreG3Out) strings(seed int64, n int) {
	p := NewJsFileRegexParser(newCoreG3Collector())
	rnd := jutil.NewRandomWithSeed(seed)
	alphabet := []rune{'\'', '"', '\\', ',', ':', '[', ']', '{', '}', 'a', 'b', ' ', '\t', '\u00e9'}
	fixed := []string{
		"'a'", "\"a\"", " 'a' ", "'a\\'b'", "'a\\\\'", "'a\\', 'b'", "['a', 'b']", "[]", "[ ]",
		"['a',]", "['a' 'b']", "['a\\', ']", "{'a': 'b'}", "{}", "{'a':'b','c':\"d\"}", "{'a'}",
		"{'a':}", "{'a': 'b',}", "  {'a': 'b'}  ", "{", "}", "'", "", "['a\\', 'b']", "]a[",
	}
	for i := 0; i < len(fixed)+n; i++ {
		var s string
		if i < len(fixed) {
			s = fixed[i]
		} else {
			length := int(rnd.NextIntBound(14))
			var sb strings.Builder
			if rnd.NextIntBound(3) == 0 {
				if rnd.NextBoolean() {
					sb.WriteByte('[')
				} else {
					sb.WriteByte('{')
				}
			}
			for k := 0; k < length; k++ {
				sb.WriteRune(alphabet[rnd.NextIntBound(int32(len(alphabet)))])
			}
			if rnd.NextIntBound(3) == 0 {
				if rnd.NextBoolean() {
					sb.WriteByte(']')
				} else {
					sb.WriteByte('}')
				}
			}
			s = sb.String()
		}
		o.line("s " + s)
		methods := []struct {
			name string
			call func() any
		}{
			{"parseJsString", func() any { return p.ParseJsString(s) }},
			{"parseJsStringArray", func() any { return p.ParseJsStringArray(s) }},
			{"parseJsStringMap", func() any { return p.ParseJsStringMap(s) }},
		}
		for _, m := range methods {
			var r string
			if t := coreG3Try(func() { r = coreG3Fmt(m.call()) }); t != nil {
				r = coreG3Exc(t)
			}
			o.line(" " + m.name + " " + r)
		}
	}
}

// ---------------------------------------------------------------------------
// chunks

type coreG3InputSpec struct {
	chunk int
	name  string
	kind  rhino.StaticSourceFile_SourceKind
	code  string
}

type coreG3GraphSpec struct {
	deps         [][]int
	explicitWeak bool
	inputs       []coreG3InputSpec
	queries      [][2]int
	manage       []string
}

// coreG3RandomGraph is CoreG3Dump.randomGraph.
func coreG3RandomGraph(rnd *jutil.Random) *coreG3GraphSpec {
	nChunks := 1 + int(rnd.NextIntBound(7))
	deps := [][]int{}
	for c := 0; c < nChunks; c++ {
		d := []int{}
		if c > 0 {
			k := 1 + int(rnd.NextIntBound(int32(min(c, 3))))
			for j := 0; j < k; j++ {
				dep := int(rnd.NextIntBound(int32(c)))
				if !slices.Contains(d, dep) {
					d = append(d, dep)
				}
			}
		}
		deps = append(deps, d)
	}
	// Rarely: chunks out of dependency order (ChunkDependenceException).
	if nChunks > 2 && rnd.NextIntBound(25) == 0 {
		deps[1] = append(deps[1], nChunks-1)
	}
	explicitWeak := rnd.NextIntBound(8) == 0
	inputs := []coreG3InputSpec{}
	nInputs := int(rnd.NextIntBound(9))
	syms := []string{"a", "b", "c", "d", "e", "goog"}
	for i := 0; i < nInputs; i++ {
		var code strings.Builder
		if rnd.NextIntBound(10) == 0 {
			code.WriteString("/** @fileoverview\n * @provideGoog */\nvar COMPILED = false;\n")
		}
		np := int(rnd.NextIntBound(3))
		for k := 0; k < np; k++ {
			code.WriteString("goog.provide('" + syms[rnd.NextIntBound(5)] + strconv.Itoa(i%3) + "');\n")
		}
		nr := int(rnd.NextIntBound(3))
		for k := 0; k < nr; k++ {
			kind := "require"
			if rnd.NextIntBound(4) == 0 {
				kind = "requireType"
			}
			sym := syms[rnd.NextIntBound(5)]
			code.WriteString("goog." + kind + "('" + sym + strconv.Itoa(int(rnd.NextIntBound(3))) + "');\n")
		}
		if rnd.NextIntBound(8) == 0 {
			code.WriteString("export {};\n")
		}
		var chunk int
		if explicitWeak && rnd.NextIntBound(4) == 0 {
			chunk = nChunks
		} else {
			chunk = int(rnd.NextIntBound(int32(nChunks)))
		}
		kind := rhino.StaticSourceFile_SourceKind_STRONG
		if rnd.NextIntBound(6) == 0 {
			kind = rhino.StaticSourceFile_SourceKind_WEAK
		}
		if explicitWeak {
			if chunk == nChunks {
				kind = rhino.StaticSourceFile_SourceKind_WEAK
			} else {
				kind = rhino.StaticSourceFile_SourceKind_STRONG
			}
			if rnd.NextIntBound(10) == 0 {
				if kind == rhino.StaticSourceFile_SourceKind_WEAK {
					kind = rhino.StaticSourceFile_SourceKind_STRONG
				} else {
					kind = rhino.StaticSourceFile_SourceKind_WEAK
				}
			}
		}
		inputs = append(inputs, coreG3InputSpec{chunk, "f" + strconv.Itoa(i) + ".js", kind, code.String()})
	}
	queries := [][2]int{}
	total := nChunks + 1
	for q := 0; q < 4; q++ {
		parent := int(rnd.NextIntBound(int32(total)))
		bits := int(rnd.NextIntBound(int32(1) << total))
		queries = append(queries, [2]int{parent, bits})
	}
	manage := []string{"sortOnly", "none"}
	eps := []string{"a0", "b1", "c2", "file:f0.js", "file:f3.js", "chunk:" + "c0" + ":a0", "zz"}
	for _, kind := range []string{"pruneLegacy", "prune", "pruneAllow"} {
		var sb strings.Builder
		sb.WriteString(kind + ":")
		ne := int(rnd.NextIntBound(3))
		for k := 0; k < ne; k++ {
			if k > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(eps[rnd.NextIntBound(int32(len(eps)))])
		}
		manage = append(manage, sb.String())
	}
	return &coreG3GraphSpec{deps, explicitWeak, inputs, queries, manage}
}

// build is CoreG3Dump.build.
func (g *coreG3GraphSpec) build() []*JSChunk {
	chunks := []*JSChunk{}
	for c := range g.deps {
		chunks = append(chunks, NewJSChunk("c"+strconv.Itoa(c)))
	}
	for c, ds := range g.deps {
		for _, d := range ds {
			chunks[c].AddDependency(chunks[d])
		}
	}
	if g.explicitWeak {
		weak := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
		for _, c := range chunks {
			weak.AddDependency(c)
		}
		chunks = append(chunks, weak)
	}
	for _, in := range g.inputs {
		chunks[in.chunk].AddSourceFile(SourceFile_FromCodeWithKind(in.name, in.code, in.kind))
	}
	return chunks
}

func coreG3ChunkNames[S ~[]*JSChunk](chunks S) string {
	l := make([]string, len(chunks))
	for i, c := range chunks {
		if c == nil {
			l[i] = "null"
		} else {
			l[i] = c.GetName()
		}
	}
	return jutil.CollectionToString(l)
}

func coreG3ChunkContents(graph *JSChunkGraph) string {
	var sb strings.Builder
	for _, c := range graph.GetAllChunks() {
		sb.WriteString(c.GetName() + "[")
		for i, in := range c.GetInputs() {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(in.GetName() + ":" + in.GetSourceFile().GetKind().String() + ":")
			if in.GetChunk() == nil {
				sb.WriteString("null")
			} else {
				sb.WriteString(in.GetChunk().GetName())
			}
		}
		sb.WriteString("] ")
	}
	return sb.String()
}

func coreG3ManageEntryPoints(a string) []*ModuleIdentifier {
	l := []*ModuleIdentifier{}
	for _, s := range coreG3List(a) {
		switch {
		case strings.HasPrefix(s, "file:"):
			l = append(l, ModuleIdentifier_ForFile(s[5:]))
		case strings.HasPrefix(s, "chunk:"):
			l = append(l, ModuleIdentifier_ForFlagValue(s[6:]))
		default:
			l = append(l, ModuleIdentifier_ForClosure(s))
		}
	}
	return l
}

func coreG3ChunkName(c *JSChunk) string {
	if c == nil {
		return "-"
	}
	return c.GetName()
}

func (o *coreG3Out) chunks(seed int64, n int) {
	rnd := jutil.NewRandomWithSeed(seed)
	for gi := 0; gi < n; gi++ {
		g := coreG3RandomGraph(rnd)
		o.line("## graph " + strconv.Itoa(gi))
		for c, d := range g.deps {
			o.line("chunk c" + strconv.Itoa(c) + " deps=" + jutil.CollectionToString(d))
		}
		o.line("explicitWeak " + strconv.FormatBool(g.explicitWeak))
		for _, in := range g.inputs {
			o.line("input " + strconv.Itoa(in.chunk) + " " + in.name + " " + in.kind.String() + " " +
				strings.ReplaceAll(in.code, "\n", "|"))
		}
		for _, q := range g.queries {
			o.line("query " + strconv.Itoa(q[0]) + " " + strconv.Itoa(q[1]))
		}
		for _, m := range g.manage {
			o.line("manage " + m)
		}
		o.line("--")
		chunks := g.build()
		var graph *JSChunkGraph
		if t := coreG3TryRuntime(func() { graph = NewJSChunkGraph(chunks) }); t != nil {
			o.line("construct " + coreG3Exc(t))
			continue
		}
		all := graph.GetAllChunks()
		o.line("count " + strconv.Itoa(graph.GetChunkCount()) + " inputs " + strconv.Itoa(graph.GetInputCount()))
		o.line("contents " + coreG3ChunkContents(graph))
		for _, c := range all {
			var sb strings.Builder
			sb.WriteString("chunk " + c.GetName() + " depth=" + strconv.Itoa(c.GetDepth()) + " index=" + strconv.Itoa(c.GetIndex()))
			sb.WriteString(" dependsOn=")
			for _, d := range all {
				if graph.DependsOn(c, d) {
					sb.WriteByte('1')
				} else {
					sb.WriteByte('0')
				}
			}
			sb.WriteString(" transitive=" + coreG3ChunkNames(graph.GetTransitiveDepsDeepestFirst(c)))
			sb.WriteString(" sortedDeps=" + jutil.CollectionToString(c.GetSortedDependencyNames()))
			sb.WriteString(" all=" + coreG3ChunkNames(c.GetAllDependencies().ToSlice()))
			sb.WriteString(" requires=" + jutil.CollectionToString(c.GetRequires()))
			o.line(sb.String())
		}
		for _, a := range all {
			var sb strings.Builder
			sb.WriteString("dcd " + a.GetName())
			for _, b := range all {
				x := graph.GetDeepestCommonDependency(a, b)
				y := graph.GetDeepestCommonDependencyInclusive(a, b)
				sb.WriteString(" " + coreG3ChunkName(x) + "/" + coreG3ChunkName(y))
			}
			o.line(sb.String())
		}
		o.line("dcdAll " + graph.GetDeepestCommonDependencyInclusiveOfCollection(all).GetName())
		for _, q := range g.queries {
			bs := jutil.NewBitSet()
			for i := range all {
				if q[1]&(1<<i) != 0 {
					bs.Set(i)
				}
			}
			var r string
			if t := coreG3TryRuntime(func() {
				r = graph.GetSmallestCoveringSubtree(all[q[0]%len(all)], bs).GetName()
			}); t != nil {
				r = coreG3Exc(t)
			}
			o.line("cover " + strconv.Itoa(q[0]) + " " + bs.String() + " " + r)
		}
		var root string
		if t := coreG3TryRuntime(func() { root = graph.GetRootChunk().GetName() }); t != nil {
			root = coreG3Exc(t)
		}
		o.line("root " + root)
		o.line("json " + graph.ToJson().String())
		for _, m := range g.manage {
			fresh := NewJSChunkGraph(g.build())
			em := newCoreG3Collector()
			compiler := NewCompiler()
			compiler.SetErrorManager(em)
			compiler.InitOptions(NewCompilerOptions())
			for _, in := range fresh.GetAllInputs() {
				in.SetCompiler(compiler)
			}
			if t := coreG3Try(func() {
				var opts *DependencyOptions
				switch m {
				case "sortOnly":
					opts = DependencyOptions_SortOnly()
				case "none":
					opts = DependencyOptions_None()
				default:
					kind, rest, _ := strings.Cut(m, ":")
					eps := coreG3ManageEntryPoints(rest)
					switch kind {
					case "pruneLegacy":
						opts = DependencyOptions_PruneLegacyForEntryPoints(eps)
					case "prune":
						opts = DependencyOptions_PruneForEntryPoints(eps)
					default:
						opts = DependencyOptions_PruneAllowNoEntryPoints(eps)
					}
				}
				res := []string{}
				for _, in := range fresh.ManageDependencies(compiler, opts).OrderedInputs() {
					res = append(res, in.GetName())
				}
				o.line("managed " + m + " " + jutil.CollectionToString(res))
				o.line(" contents " + coreG3ChunkContents(fresh))
			}); t != nil {
				o.line("managed " + m + " " + coreG3Exc(t))
			}
			for _, r := range em.reports {
				if strings.HasPrefix(r, "ERROR ") || strings.Contains(r, "WEAK") {
					o.line(" report " + r)
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Comparison

// coreG3ReadDump reads a gzipped dump and splits it into sections, each
// starting at a "## " line (the text before the first one is section "").
func coreG3ReadDump(path string) (map[string]string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, nil, err
	}
	var sb strings.Builder
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	for sc.Scan() {
		sb.WriteString(sc.Text())
		sb.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		return nil, nil, err
	}
	sections, order := coreG3Sections(sb.String())
	return sections, order, nil
}

func coreG3Sections(text string) (map[string]string, []string) {
	sections := map[string]string{}
	var order []string
	current := ""
	var sb strings.Builder
	for _, l := range strings.SplitAfter(text, "\n") {
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "## ") {
			sections[current] += sb.String()
			sb.Reset()
			current = strings.TrimSuffix(l, "\n")
			order = append(order, current)
		}
		sb.WriteString(l)
	}
	sections[current] += sb.String()
	return sections, order
}

// coreG3FirstDiff describes the first differing line.
func coreG3FirstDiff(got, want string) string {
	g := strings.Split(got, "\n")
	w := strings.Split(want, "\n")
	for i := 0; i < len(g) || i < len(w); i++ {
		var gl, wl string
		if i < len(g) {
			gl = g[i]
		}
		if i < len(w) {
			wl = w[i]
		}
		if gl != wl {
			if len(gl) > 200 || len(wl) > 200 {
				// Show the neighbourhood of the first differing byte.
				c := 0
				for c < len(gl) && c < len(wl) && gl[c] == wl[c] {
					c++
				}
				from := max(0, c-80)
				gl, wl = "..."+gl[from:min(len(gl), c+120)], "..."+wl[from:min(len(wl), c+120)]
			}
			return fmt.Sprintf("line %d:\n  go:   %s\n  java: %s", i+1, gl, wl)
		}
	}
	return "identical?"
}

// coreG3Compare compares the sections of got against the Java dump's
// sections named in order.
func coreG3Compare(t *testing.T, name string, got string, want map[string]string, order []string, diffDir string) {
	t.Helper()
	gotSections, gotOrder := coreG3Sections(got)
	if want[""] != gotSections[""] {
		t.Errorf("%s: preamble differs: %s", name, coreG3FirstDiff(gotSections[""], want[""]))
	}
	if len(gotOrder) != len(order) {
		t.Errorf("%s: %d sections, want %d", name, len(gotOrder), len(order))
	}
	mismatches := 0
	for _, sec := range order {
		if gotSections[sec] == want[sec] {
			continue
		}
		mismatches++
		if mismatches <= 10 {
			t.Errorf("%s: %s: %s", name, sec, coreG3FirstDiff(gotSections[sec], want[sec]))
		}
		if diffDir != "" {
			base := filepath.Join(diffDir, name, strings.NewReplacer("/", "_", " ", "_").Replace(sec))
			_ = os.MkdirAll(filepath.Dir(base), 0o755)
			_ = os.WriteFile(base+".go.txt", []byte(gotSections[sec]), 0o644)
			_ = os.WriteFile(base+".java.txt", []byte(want[sec]), 0o644)
		}
	}
	if mismatches != 0 {
		t.Errorf("%s: %d of %d sections differ", name, mismatches, len(order))
	}
}

// coreG3DepsFiles dumps the deps mode for the "## file" sections of a Java
// dump, every every-th one if every > 1, and returns the Go dump with the
// Java sections it covers.
func coreG3DepsFiles(t *testing.T, root string, want map[string]string, order []string, every int) (string, []string) {
	t.Helper()
	var o coreG3Out
	var covered []string
	for i, sec := range order {
		if every > 1 && i%every != 0 {
			continue
		}
		rel, ok := strings.CutPrefix(sec, "## file ")
		if !ok {
			t.Fatalf("unexpected section %q", sec)
		}
		b, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		o.line("## file " + rel)
		o.depsOf(rel, jutil.NewStringFromUTF8Bytes(b))
		covered = append(covered, sec)
	}
	return o.sb.String(), covered
}

func TestCoreG3OracleGolden(t *testing.T) {
	dir := filepath.Join("testdata", "core_g3")
	for _, tc := range []struct {
		name string
		dump func(o *coreG3Out)
	}{
		{"options", func(o *coreG3Out) { o.options(0, 300) }},
		{"strings", func(o *coreG3Out) { o.strings(0, 3000) }},
		{"chunks", func(o *coreG3Out) { o.chunks(0, 300) }},
		{"depsgen", func(o *coreG3Out) { o.depsgen(0, 1500) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want, order, err := coreG3ReadDump(filepath.Join(dir, tc.name+".txt.gz"))
			if err != nil {
				t.Fatal(err)
			}
			var o coreG3Out
			tc.dump(&o)
			coreG3Compare(t, tc.name, o.sb.String(), want, order, "")
		})
	}
	t.Run("fixtures-deps", func(t *testing.T) {
		want, order, err := coreG3ReadDump(filepath.Join(dir, "fixtures-deps.txt.gz"))
		if err != nil {
			t.Fatal(err)
		}
		got, covered := coreG3DepsFiles(t, filepath.Join("..", "..", "testdata", "fixtures"), want, order, 1)
		coreG3Compare(t, "fixtures-deps", got, want, covered, "")
	})
}

func TestCoreG3OracleCorpus(t *testing.T) {
	root, oracle := os.Getenv("CORE_G3_ROOT"), os.Getenv("CORE_G3_ORACLE")
	if root == "" || oracle == "" {
		t.Skip("set CORE_G3_ROOT and CORE_G3_ORACLE to compare a corpus")
	}
	want, order, err := coreG3ReadDump(oracle)
	if err != nil {
		t.Fatal(err)
	}
	every, _ := strconv.Atoi(os.Getenv("CORE_G3_EVERY"))
	got, covered := coreG3DepsFiles(t, root, want, order, every)
	coreG3Compare(t, "corpus-deps", got, want, covered, os.Getenv("CORE_G3_DIFF_DIR"))
}
