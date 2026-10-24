// Differential oracle for the jscomp "compiler state" core classes (group g3 of wave 3):
// CompilerOptions, CompilerInput (its DepsFinder), LazyParsedDependencyInfo, JSChunk, JSChunkGraph,
// deps/JsFileRegexParser and deps/JsFileLineParser (internal/jscomp/compiler_options.go,
// compiler_input.go, lazy_parsed_dependency_info.go, js_chunk.go, js_chunk_graph.go,
// deps_js_file_regex_parser.go, deps_js_file_line_parser.go).
//
// Subcommands (output gzipped, every line escaped like ParsingDump: backslash, and every char
// outside printable ASCII as \\uXXXX):
//
//   options <out.gz> <seed> <n>
//     "## statics": the nested enums (values, and their methods / fields), LanguageMode.fromString
//     and InstrumentOption.fromString of fixed strings, validCommandLineNames(), and
//     BrowserFeaturesetYear.minimumRequiredFor(f) for every FeatureSet.Feature.
//     Then n scenarios "## options <i>": a random sequence of setter calls drawn (java.util.Random
//     (seed + i)) from a fixed catalog of ops, each printed as "op <name> <arg>" followed by
//     "ok" or "! <SimpleExceptionClass>: <message>"; then "toString <CompilerOptions.toString()>"
//     and "get <name> <value>" for every zero-argument instance method of CompilerOptions that
//     returns a value (sorted by name). Object.toString() identity hashes print as "@ID".
//
//   deps <out.gz> <root> --list=<listfile> | <rel>...
//     per file "## file <rel>":
//       "regex <DependencyInfo>" + its errors/warnings, for JsFileRegexParser(em) (shortcut mode
//       off, no goog base), then "regexgoog ..." with setShortcutMode(true) and
//       setIncludeGoogBase(true);
//       "ast ..." the dependency fields of a CompilerInput whose compiler is new Compiler()
//       + setErrorManager(collector) + initOptions(new CompilerOptions()) (DepsFinder over the AST),
//       and "lazy ..." the same with compiler.setPreferRegexParser(true) (JsFileRegexParser +
//       LazyParsedDependencyInfo). Errors (not warnings) the compiler collected are printed too.
//
//   depsgen <out.gz> <seed> <n>
//     the same as deps for n generated sources ("## gen <name> <code, newlines as |>"): random
//     sequences of lines drawn from GEN_LINES (goog.provide/require/module lines, ES module
//     import/export lines, JSDoc annotations, comments and strings that hide them) and, one
//     line in 25, from GEN_BAD (fragments that do not parse).
//
//   strings <out.gz> <seed> <n>
//     n random strings over the alphabet ' " \ , : [ ] { } a b space, and for each the results (or
//     exceptions) of JsFileLineParser.parseJsString, parseJsStringArray and parseJsStringMap.
//
//   chunks <out.gz> <seed> <n>
//     n random chunk graphs ("## graph <i>"), each described first (chunks with dependencies,
//     inputs with provides/requires/type requires/kind, queries), then: construction, depths,
//     dependsOn, getDeepestCommonDependency[Inclusive] for every pair, getTransitiveDepsDeepestFirst,
//     getSmallestCoveringSubtree queries, toJson, getRootChunk, and manageDependencies under
//     four DependencyOptions (each on a freshly built graph): result order, chunk contents, kinds,
//     reported errors, or the exception.
//
// Usage: java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.CoreG3Dump <cmd> ...
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.BrowserFeaturesetYear;
import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.deps.DependencyInfo;
import com.google.javascript.jscomp.deps.JsFileLineParser;
import com.google.javascript.jscomp.deps.JsFileRegexParser;
import com.google.javascript.jscomp.deps.ModuleLoader;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.javascript.jscomp.js.RuntimeJsLibManager.RuntimeLibraryMode;
import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.parser.FeatureSet;
import com.google.javascript.rhino.StaticSourceFile.SourceKind;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.BitSet;
import java.util.Collection;
import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Random;
import java.util.function.BiConsumer;
import java.util.regex.Pattern;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class CoreG3Dump {
  private static PrintStream out;

  static String esc(String s) {
    StringBuilder b = new StringBuilder(s.length());
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        b.append("\\\\");
      } else if (c >= 0x20 && c < 0x7f) {
        b.append(c);
      } else {
        b.append(String.format("\\u%04x", (int) c));
      }
    }
    return b.toString();
  }

  static void line(String s) {
    out.print(esc(s));
    out.print('\n');
  }

  private static final Pattern IDENTITY = Pattern.compile("(com\\.google\\.[\\w.$]+)@[0-9a-f]+");

  /** Masks Object.toString() identity hashes. */
  static String mask(String s) {
    return IDENTITY.matcher(s).replaceAll("$1@ID");
  }

  static String exc(Throwable t) {
    return "! " + t.getClass().getSimpleName() + ": " + t.getMessage();
  }

  public static void main(String[] args) throws Exception {
    OutputStream os = new GZIPOutputStream(Files.newOutputStream(Paths.get(args[1])), 1 << 16);
    out = new PrintStream(os, false, StandardCharsets.UTF_8);
    switch (args[0]) {
      case "options" -> options(Long.parseLong(args[2]), Integer.parseInt(args[3]));
      case "deps" -> deps(args);
      case "depsgen" -> depsgen(Long.parseLong(args[2]), Integer.parseInt(args[3]));
      case "strings" -> strings(Long.parseLong(args[2]), Integer.parseInt(args[3]));
      case "chunks" -> chunks(Long.parseLong(args[2]), Integer.parseInt(args[3]));
      default -> throw new IllegalArgumentException(args[0]);
    }
    out.close();
  }

  // ---------------------------------------------------------------------------------------------
  // options

  /** Java's String.valueOf for option values, with collections printed element by element. */
  static String fmt(Object v) {
    if (v == null) {
      return "null";
    }
    if (v instanceof Optional<?> o) {
      return o.isPresent() ? "Optional.of(" + fmt(o.get()) + ")" : "Optional.empty";
    }
    if (v instanceof com.google.javascript.jscomp.jarjar.com.google.common.base.Optional<?> o) {
      return o.isPresent() ? "Optional.of(" + fmt(o.get()) + ")" : "Optional.absent()";
    }
    if (v instanceof Map<?, ?> m) {
      StringBuilder sb = new StringBuilder("{");
      boolean first = true;
      for (Map.Entry<?, ?> e : m.entrySet()) {
        if (!first) {
          sb.append(", ");
        }
        first = false;
        sb.append(fmt(e.getKey())).append('=').append(fmt(e.getValue()));
      }
      return sb.append('}').toString();
    }
    if (v instanceof Collection<?> c) {
      StringBuilder sb = new StringBuilder("[");
      boolean first = true;
      for (Object e : c) {
        if (!first) {
          sb.append(", ");
        }
        first = false;
        sb.append(fmt(e));
      }
      return sb.append(']').toString();
    }
    if (v instanceof Enum<?> e) {
      return e.name();
    }
    return mask(String.valueOf(v));
  }

  record Op(String name, String[] args, BiConsumer<CompilerOptions, String> apply) {}

  static final List<Op> OPS = new ArrayList<>();

  static String[] bools() {
    return new String[] {"true", "false"};
  }

  static <E extends Enum<E>> String[] names(Class<E> c) {
    return Arrays.stream(c.getEnumConstants()).map(Enum::name).toArray(String[]::new);
  }

  static void op(String name, String[] args, BiConsumer<CompilerOptions, String> apply) {
    OPS.add(new Op(name, args, apply));
  }

  static void boolOp(String name, BiConsumer<CompilerOptions, Boolean> apply) {
    op(name, bools(), (o, a) -> apply.accept(o, Boolean.parseBoolean(a)));
  }

  static String nullable(String a) {
    return a.equals("null") ? null : a;
  }

  static List<String> list(String a) {
    return a.isEmpty() ? ImmutableList.of() : ImmutableList.copyOf(a.split(","));
  }

  static ImmutableSet<String> set(String a) {
    return ImmutableSet.copyOf(list(a));
  }

  static List<ModuleIdentifier> entryPoints(String a) {
    List<ModuleIdentifier> l = new ArrayList<>();
    for (String s : list(a)) {
      l.add(s.startsWith("file:") ? ModuleIdentifier.forFile(s.substring(5)) : ModuleIdentifier.forClosure(s));
    }
    return l;
  }

  static {
    String[] lang = names(LanguageMode.class);
    String[] reach = names(CompilerOptions.Reach.class);
    op("setLanguageIn", lang, (o, a) -> o.setLanguageIn(LanguageMode.valueOf(a)));
    op("setLanguageOut", lang, (o, a) -> o.setLanguageOut(LanguageMode.valueOf(a)));
    op("setLanguage", lang, (o, a) -> o.setLanguage(LanguageMode.valueOf(a)));
    op(
        "setBrowserFeaturesetYear",
        new String[] {"2012", "2017", "2018", "2019", "2020", "2021", "2022", "2023", "2024", "2025", "2026", "2027"},
        (o, a) -> o.setBrowserFeaturesetYear(Integer.parseInt(a)));
    boolOp("setCheckTypes", CompilerOptions::setCheckTypes);
    boolOp("setCheckSymbols", CompilerOptions::setCheckSymbols);
    boolOp("setCheckSuspiciousCode", CompilerOptions::setCheckSuspiciousCode);
    boolOp("setFoldConstants", CompilerOptions::setFoldConstants);
    boolOp("setDeadAssignmentElimination", CompilerOptions::setDeadAssignmentElimination);
    boolOp("setDeadPropertyAssignmentElimination", CompilerOptions::setDeadPropertyAssignmentElimination);
    boolOp("setInlineConstantVars", CompilerOptions::setInlineConstantVars);
    boolOp("setCrossChunkCodeMotion", CompilerOptions::setCrossChunkCodeMotion);
    boolOp("setCrossChunkMethodMotion", CompilerOptions::setCrossChunkMethodMotion);
    boolOp("setCoalesceVariableNames", CompilerOptions::setCoalesceVariableNames);
    boolOp("setInlineVariables", (o, b) -> o.setInlineVariables(b));
    op("setInlineVariablesOfReach", reach, (o, a) -> o.setInlineVariables(CompilerOptions.Reach.valueOf(a)));
    op("setRemoveUnusedVariables", reach, (o, a) -> o.setRemoveUnusedVariables(CompilerOptions.Reach.valueOf(a)));
    op("setInlineFunctions", reach, (o, a) -> o.setInlineFunctions(CompilerOptions.Reach.valueOf(a)));
    boolOp("setInlineProperties", CompilerOptions::setInlineProperties);
    boolOp("setSmartNameRemoval", CompilerOptions::setSmartNameRemoval);
    boolOp("setExtractPrototypeMemberDeclarations", (o, b) -> o.setExtractPrototypeMemberDeclarations(b));
    op(
        "setExtractPrototypeMemberDeclarationsOfMode",
        names(CompilerOptions.ExtractPrototypeMemberDeclarationsMode.class),
        (o, a) ->
            o.setExtractPrototypeMemberDeclarations(
                CompilerOptions.ExtractPrototypeMemberDeclarationsMode.valueOf(a)));
    boolOp("setRemoveUnusedPrototypeProperties", CompilerOptions::setRemoveUnusedPrototypeProperties);
    boolOp("setRemoveUnusedClassProperties", CompilerOptions::setRemoveUnusedClassProperties);
    boolOp("setCollapseProperties", CompilerOptions::setCollapseProperties);
    op(
        "setCollapsePropertiesLevel",
        names(CompilerOptions.PropertyCollapseLevel.class),
        (o, a) -> o.setCollapsePropertiesLevel(CompilerOptions.PropertyCollapseLevel.valueOf(a)));
    boolOp("setDisambiguateProperties", CompilerOptions::setDisambiguateProperties);
    boolOp("setAmbiguateProperties", CompilerOptions::setAmbiguateProperties);
    boolOp("setUseTypesForOptimization", CompilerOptions::setUseTypesForOptimization);
    boolOp("setUseTypesForLocalOptimization", CompilerOptions::setUseTypesForLocalOptimization);
    op(
        "setPolymerVersion",
        new String[] {"null", "1", "2", "3"},
        (o, a) -> o.setPolymerVersion(a.equals("null") ? null : Integer.valueOf(a)));
    boolOp("setAngularPass", CompilerOptions::setAngularPass);
    boolOp("setChromePass", CompilerOptions::setChromePass);
    boolOp("setClosurePass", CompilerOptions::setClosurePass);
    boolOp("setPreserveClosurePrimitives", CompilerOptions::setPreserveClosurePrimitives);
    op(
        "setIncrementalChecks",
        names(CompilerOptions.IncrementalCheckMode.class),
        (o, a) -> o.setIncrementalChecks(CompilerOptions.IncrementalCheckMode.valueOf(a)));
    boolOp("setEmitUseStrict", CompilerOptions::setEmitUseStrict);
    boolOp("setStrictModeInput", CompilerOptions::setStrictModeInput);
    boolOp("setIsolatePolyfills", CompilerOptions::setIsolatePolyfills);
    boolOp("setInstrumentAsyncContext", CompilerOptions::setInstrumentAsyncContext);
    boolOp("setRewritePolyfills", CompilerOptions::setRewritePolyfills);
    boolOp("setPreventLibraryInjection", CompilerOptions::setPreventLibraryInjection);
    op(
        "setRuntimeLibraryMode",
        names(RuntimeLibraryMode.class),
        (o, a) -> o.setRuntimeLibraryMode(RuntimeLibraryMode.valueOf(a)));
    op(
        "setDefineToBooleanLiteral",
        new String[] {"A:true", "B:false", "goog.DEBUG:true"},
        (o, a) -> o.setDefineToBooleanLiteral(a.split(":")[0], Boolean.parseBoolean(a.split(":")[1])));
    op(
        "setDefineToNumberLiteral",
        new String[] {"A:3", "C:-7", "goog.FEATURESET_YEAR:2000"},
        (o, a) -> o.setDefineToNumberLiteral(a.split(":")[0], Integer.parseInt(a.split(":")[1])));
    op(
        "setDefineToDoubleLiteral",
        new String[] {"A:0.5", "D:1e21", "E:-0.0"},
        (o, a) -> o.setDefineToDoubleLiteral(a.split(":")[0], Double.parseDouble(a.split(":")[1])));
    op(
        "setDefineToStringLiteral",
        new String[] {"A:x", "F:hello world", "G:\u00e9"},
        (o, a) -> o.setDefineToStringLiteral(a.split(":")[0], a.split(":")[1]));
    op("setRenamePrefix", new String[] {"null", "", "$p_"}, (o, a) -> o.setRenamePrefix(nullable(a)));
    op("setRenamePrefixNamespace", new String[] {"null", "_ns"}, (o, a) -> o.setRenamePrefixNamespace(nullable(a)));
    op("setLocale", new String[] {"null", "en", "fr_CA"}, (o, a) -> o.setLocale(nullable(a)));
    op(
        "setSourceMapOutputPath",
        new String[] {"null", "", "out.map"},
        (o, a) -> o.setSourceMapOutputPath(nullable(a)));
    boolOp("setAlwaysGatherSourceMapInfo", CompilerOptions::setAlwaysGatherSourceMapInfo);
    boolOp("setSkipNonTranspilationPasses", CompilerOptions::setSkipNonTranspilationPasses);
    boolOp("setChecksOnly", CompilerOptions::setChecksOnly);
    op(
        "setOutputJs",
        names(CompilerOptions.OutputJs.class),
        (o, a) -> o.setOutputJs(CompilerOptions.OutputJs.valueOf(a)));
    boolOp("setGenerateExports", CompilerOptions::setGenerateExports);
    boolOp("setPrettyPrint", CompilerOptions::setPrettyPrint);
    boolOp("setLineBreak", CompilerOptions::setLineBreak);
    boolOp("setPrintInputDelimiter", CompilerOptions::setPrintInputDelimiter);
    op("setInputDelimiter", new String[] {"", "// file %name%"}, (o, a) -> o.setInputDelimiter(a));
    op(
        "setTracerMode",
        names(CompilerOptions.TracerMode.class),
        (o, a) -> o.setTracerMode(CompilerOptions.TracerMode.valueOf(a)));
    op("setErrorFormat", names(ErrorFormat.class), (o, a) -> o.setErrorFormat(ErrorFormat.valueOf(a)));
    op(
        "setLineLengthThreshold",
        new String[] {"0", "80", "-1"},
        (o, a) -> o.setLineLengthThreshold(Integer.parseInt(a)));
    op(
        "setVariableRenaming",
        names(VariableRenamingPolicy.class),
        (o, a) -> o.setVariableRenaming(VariableRenamingPolicy.valueOf(a)));
    op(
        "setPropertyRenaming",
        names(PropertyRenamingPolicy.class),
        (o, a) -> o.setPropertyRenaming(PropertyRenamingPolicy.valueOf(a)));
    op("setModuleRoots", new String[] {"", "a", "a,b/c"}, (o, a) -> o.setModuleRoots(list(a)));
    op("setCommonJSModulePathPrefix", new String[] {"lib/"}, (o, a) -> o.setCommonJSModulePathPrefix(a));
    boolOp("setProcessCommonJSModules", CompilerOptions::setProcessCommonJSModules);
    boolOp(
        "setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf",
        CompilerOptions::setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf);
    boolOp("setEnableModuleRewriting", CompilerOptions::setEnableModuleRewriting);
    op(
        "setModuleResolutionMode",
        names(ModuleLoader.ResolutionMode.class),
        (o, a) -> o.setModuleResolutionMode(ModuleLoader.ResolutionMode.valueOf(a)));
    op(
        "setPathEscaper",
        names(ModuleLoader.PathEscaper.class),
        (o, a) -> o.setPathEscaper(ModuleLoader.PathEscaper.valueOf(a)));
    op("setPackageJsonEntryNames", new String[] {"", "main"}, (o, a) -> o.setPackageJsonEntryNames(list(a)));
    op(
        "setDependencyOptions",
        new String[] {"none", "sortOnly", "pruneLegacy:a", "prune:a,file:b.js", "pruneAllow:"},
        (o, a) -> o.setDependencyOptions(depOptions(a)));
    op(
        "setWarningLevel",
        names(CheckLevel.class),
        (o, a) ->
            o.setWarningLevel(
                DiagnosticGroup.forType(RhinoErrorReporter.TYPE_PARSE_ERROR), CheckLevel.valueOf(a)));
    op("resetWarningsGuard", new String[] {""}, (o, a) -> o.resetWarningsGuard());
    op("setIdGenerators", new String[] {"", "goog.getId", "a,b"}, (o, a) -> o.setIdGenerators(set(a)));
    op("setStripTypes", new String[] {"", "goog.debug"}, (o, a) -> o.setStripTypes(set(a)));
    op("setStripNameSuffixes", new String[] {"", "logger,Logger"}, (o, a) -> o.setStripNameSuffixes(set(a)));
    op("setStripNamePrefixes", new String[] {"", "x."}, (o, a) -> o.setStripNamePrefixes(set(a)));
    op(
        "setReplaceStringsConfiguration",
        new String[] {"`:", "#:f(?)"},
        (o, a) -> o.setReplaceStringsConfiguration(a.split(":", 2)[0], list(a.split(":", 2)[1])));
    op(
        "setPropertiesThatMustDisambiguate",
        new String[] {"", "p,q"},
        (o, a) -> o.setPropertiesThatMustDisambiguate(set(a)));
    op("setExtraAnnotationNames", new String[] {"", "foo,bar"}, (o, a) -> o.setExtraAnnotationNames(list(a)));
    op("setForceLibraryInjection", new String[] {"", "es6_runtime"}, (o, a) -> o.setForceLibraryInjection(list(a)));
    op(
        "setUnknownDefinesToIgnore",
        new String[] {"", "X,Y"},
        (o, a) -> o.setUnknownDefinesToIgnore(ImmutableList.copyOf(list(a))));
    op(
        "addUnknownDefinesToIgnore",
        new String[] {"Z"},
        (o, a) -> o.addUnknownDefinesToIgnore(ImmutableList.copyOf(list(a))));
    op(
        "setReplaceMessagesWithChromeI18n",
        new String[] {"true:null", "true:proj", "false:null"},
        (o, a) ->
            o.setReplaceMessagesWithChromeI18n(
                Boolean.parseBoolean(a.split(":")[0]), nullable(a.split(":")[1])));
    boolOp("setDoLateLocalization", CompilerOptions::setDoLateLocalization);
    op(
        "setMessageBundle",
        new String[] {"null", "empty"},
        (o, a) -> o.setMessageBundle(a.equals("null") ? null : new EmptyMessageBundle()));
    boolOp("setValidateRequiredInlinings", CompilerOptions::setValidateRequiredInlinings);
    op(
        "setMaxFunctionSizeAfterInlining",
        new String[] {"0", "100"},
        (o, a) -> o.setMaxFunctionSizeAfterInlining(Integer.parseInt(a)));
    op(
        "setInstrumentForCoverageOption",
        names(CompilerOptions.InstrumentOption.class),
        (o, a) -> o.setInstrumentForCoverageOption(CompilerOptions.InstrumentOption.valueOf(a)));
    boolOp("setInstrumentForCoverageOnly", CompilerOptions::setInstrumentForCoverageOnly);
    op(
        "setProductionInstrumentationArrayName",
        new String[] {"", "arr"},
        (o, a) -> o.setProductionInstrumentationArrayName(a));
    op("setDebugLogFilter", new String[] {"null", "a,,b"}, (o, a) -> o.setDebugLogFilter(nullable(a)));
    boolOp("setTrustedStrings", CompilerOptions::setTrustedStrings);
    boolOp("setPreferSingleQuotes", CompilerOptions::setPreferSingleQuotes);
    boolOp("setQuoteKeywordProperties", CompilerOptions::setQuoteKeywordProperties);
    boolOp("setAssumeStrictThis", CompilerOptions::setAssumeStrictThis);
    boolOp("setProtectHiddenSideEffects", CompilerOptions::setProtectHiddenSideEffects);
    boolOp("setAssumeGettersArePure", CompilerOptions::setAssumeGettersArePure);
    op(
        "setDevMode",
        names(CompilerOptions.DevMode.class),
        (o, a) -> o.setDevMode(CompilerOptions.DevMode.valueOf(a)));
    op(
        "setEnvironment",
        names(CompilerOptions.Environment.class),
        (o, a) -> o.setEnvironment(CompilerOptions.Environment.valueOf(a)));
    op(
        "setTweakProcessing",
        names(CompilerOptions.TweakProcessing.class),
        (o, a) -> o.setTweakProcessing(CompilerOptions.TweakProcessing.valueOf(a)));
    op(
        "setChunkOutputType",
        names(CompilerOptions.ChunkOutputType.class),
        (o, a) -> o.setChunkOutputType(CompilerOptions.ChunkOutputType.valueOf(a)));
    op(
        "setEs6ModuleTranspilation",
        names(CompilerOptions.Es6ModuleTranspilation.class),
        (o, a) -> o.setEs6ModuleTranspilation(CompilerOptions.Es6ModuleTranspilation.valueOf(a)));
    op(
        "setEs6SubclassTranspilation",
        names(CompilerOptions.Es6SubclassTranspilation.class),
        (o, a) -> o.setEs6SubclassTranspilation(CompilerOptions.Es6SubclassTranspilation.valueOf(a)));
    op(
        "setConformanceReportingMode",
        names(CompilerOptions.ConformanceReportingMode.class),
        (o, a) -> o.setConformanceReportingMode(CompilerOptions.ConformanceReportingMode.valueOf(a)));
    boolOp("setContinueAfterErrors", CompilerOptions::setContinueAfterErrors);
    op(
        "setParseJsDocDocumentation",
        names(Config.JsDocParsing.class),
        (o, a) -> o.setParseJsDocDocumentation(Config.JsDocParsing.valueOf(a)));
    boolOp("setPreserveDetailedSourceInfo", CompilerOptions::setPreserveDetailedSourceInfo);
    op(
        "setSourceMapDetailLevel",
        names(SourceMap.DetailLevel.class),
        (o, a) -> o.setSourceMapDetailLevel(SourceMap.DetailLevel.valueOf(a)));
    op(
        "setSourceMapFormat",
        names(SourceMap.Format.class),
        (o, a) -> o.setSourceMapFormat(SourceMap.Format.valueOf(a)));
    op(
        "setSourceMapLocationMappings",
        new String[] {"", "a/|b/", "a/|,c|d"},
        (o, a) -> {
          List<SourceMap.LocationMapping> l = new ArrayList<>();
          for (String m : list(a)) {
            l.add(new SourceMap.PrefixLocationMapping(m.split("\\|", -1)[0], m.split("\\|", -1)[1]));
          }
          o.setSourceMapLocationMappings(l);
        });
    op("setDynamicImportAlias", new String[] {"null", "imp"}, (o, a) -> o.setDynamicImportAlias(nullable(a)));
    boolOp("setAllowDynamicImport", CompilerOptions::setAllowDynamicImport);
    op(
        "setExperimentalOutputFeatureSet",
        new String[] {
          "null", "ES5_WITH_SOME_PERFORMANT_ES2015_AND_ASYNC_FUNCTIONS", "BROWSER_2019_WITHOUT_CLASSES_AND_SPREAD"
        },
        (o, a) ->
            o.setExperimentalOutputFeatureSet(
                a.equals("null") ? null : CompilerOptions.ExperimentalOutputFeatureSet.valueOf(a)));
    boolOp("setTranspilePublicClassFields", CompilerOptions::setTranspilePublicClassFields);
    boolOp("setOptimizeLetAndConst", CompilerOptions::setOptimizeLetAndConst);
    boolOp("setAssumeGlobalScopeIsIsolated", CompilerOptions::setAssumeGlobalScopeIsIsolated);
    op(
        "setOptimizeLocalAccessForGlobalSymbolNamespace",
        names(CompilerOptions.OptimizeLocalAccess.class),
        (o, a) ->
            o.setOptimizeLocalAccessForGlobalSymbolNamespace(
                CompilerOptions.OptimizeLocalAccess.valueOf(a)));
    op("setCssRenamingSkiplist", new String[] {"", "x"}, (o, a) -> o.setCssRenamingSkiplist(set(a)));
    op(
        "setInjectPolyfillsNewerThan",
        new String[] {"null", "ECMASCRIPT5", "ECMASCRIPT_2015"},
        (o, a) -> o.setInjectPolyfillsNewerThan(a.equals("null") ? null : LanguageMode.valueOf(a)));
    boolOp("setRewriteGlobalDeclarationsForTryCatchWrapping", CompilerOptions::setRewriteGlobalDeclarationsForTryCatchWrapping);
    boolOp("setGatherCssNames", CompilerOptions::setGatherCssNames);
    boolOp("setMarkAsCompiled", CompilerOptions::setMarkAsCompiled);
    boolOp("setLabelRenaming", CompilerOptions::setLabelRenaming);
    boolOp("setGeneratePseudoNames", CompilerOptions::setGeneratePseudoNames);
    boolOp("setDevirtualizeMethods", CompilerOptions::setDevirtualizeMethods);
    boolOp("setComputeFunctionSideEffects", CompilerOptions::setComputeFunctionSideEffects);
    boolOp("setOptimizeCalls", CompilerOptions::setOptimizeCalls);
    boolOp("setConvertToDottedProperties", CompilerOptions::setConvertToDottedProperties);
    op(
        "setAliasStringsMode",
        names(CompilerOptions.AliasStringsMode.class),
        (o, a) -> o.setAliasStringsMode(CompilerOptions.AliasStringsMode.valueOf(a)));
    op(
        "setJ2clPass",
        names(CompilerOptions.J2clPassMode.class),
        (o, a) -> o.setJ2clPass(CompilerOptions.J2clPassMode.valueOf(a)));
    op("setSummaryDetailLevel", new String[] {"0", "3"}, (o, a) -> o.setSummaryDetailLevel(Integer.parseInt(a)));
    op("setNumParallelThreads", new String[] {"4"}, (o, a) -> o.setNumParallelThreads(Integer.parseInt(a)));
    op(
        "setConformanceRemoveRegexFromPath",
        new String[] {"absent", "^x/"},
        (o, a) ->
            o.setConformanceRemoveRegexFromPath(
                a.equals("absent")
                    ? com.google.javascript.jscomp.jarjar.com.google.common.base.Optional.absent()
                    : com.google.javascript.jscomp.jarjar.com.google.common.base.Optional.of(
                        Pattern.compile(a))));
  }

  static DependencyOptions depOptions(String a) {
    if (a.equals("none")) {
      return DependencyOptions.none();
    } else if (a.equals("sortOnly")) {
      return DependencyOptions.sortOnly();
    }
    String kind = a.substring(0, a.indexOf(':'));
    List<ModuleIdentifier> eps = entryPoints(a.substring(a.indexOf(':') + 1));
    return switch (kind) {
      case "pruneLegacy" -> DependencyOptions.pruneLegacyForEntryPoints(eps);
      case "prune" -> DependencyOptions.pruneForEntryPoints(eps);
      case "pruneAllow" -> DependencyOptions.pruneAllowNoEntryPoints(eps);
      default -> throw new IllegalArgumentException(a);
    };
  }

  static void options(long seed, int n) throws Exception {
    line("## statics");
    for (LanguageMode m : LanguageMode.values()) {
      String fs;
      try {
        fs = m.toFeatureSet().toString();
      } catch (RuntimeException e) {
        fs = exc(e);
      }
      line("LanguageMode " + m.name() + " " + m.ordinal() + " strict=" + m.isDefaultStrict() + " " + fs);
    }
    line("STABLE_IN " + LanguageMode.STABLE_IN + " STABLE_OUT " + LanguageMode.STABLE_OUT);
    line("validCommandLineNames " + LanguageMode.validCommandLineNames());
    for (String s :
        new String[] {
          "ECMASCRIPT3", "  es3  ", "junk", "ECMASCRIPT_2020", "  es_2020  ", "  es2020  ", "es6",
          "ES6_STRICT", "ecmascript_next", "STABLE", "Es5", "ESNEXT", "ES_NEXT", "", "estable"
        }) {
      line("fromString " + s + " -> " + LanguageMode.fromString(s));
    }
    for (String s : new String[] {"NONE", "LINE", "BRANCH", "PRODUCTION", "LINE_ONLY", "x"}) {
      line("InstrumentOption.fromString " + s + " -> " + CompilerOptions.InstrumentOption.fromString(s));
    }
    for (BrowserFeaturesetYear y : BrowserFeaturesetYear.values()) {
      line("BrowserFeaturesetYear " + y.name() + " " + y.getYear() + " " + y.getFeatureSet());
    }
    for (FeatureSet.Feature f : FeatureSet.Feature.values()) {
      line("minimumRequiredFor " + f.name() + " " + BrowserFeaturesetYear.minimumRequiredFor(f));
    }
    for (CompilerOptions.Reach r : CompilerOptions.Reach.values()) {
      line("Reach " + r + " " + r.isOn() + " " + r.includesGlobals());
    }
    for (CompilerOptions.TracerMode r : CompilerOptions.TracerMode.values()) {
      line("TracerMode " + r + " " + r.isOn() + " " + r.doPruningAnalysis());
    }
    for (CompilerOptions.TweakProcessing r : CompilerOptions.TweakProcessing.values()) {
      line("TweakProcessing " + r + " " + r.isOn() + " " + r.shouldStrip());
    }
    for (CompilerOptions.J2clPassMode r : CompilerOptions.J2clPassMode.values()) {
      line("J2clPassMode " + r + " " + r.shouldAddJ2clPasses());
    }
    List<Class<?>> enums = new ArrayList<>();
    for (Class<?> c : CompilerOptions.class.getDeclaredClasses()) {
      if (c.isEnum()) {
        enums.add(c);
      }
    }
    enums.sort(Comparator.comparing(Class::getSimpleName));
    for (Class<?> c : enums) {
      line("enum " + c.getSimpleName() + " " + Arrays.toString(c.getEnumConstants()));
    }

    List<Method> getters = new ArrayList<>();
    for (Method m : CompilerOptions.class.getDeclaredMethods()) {
      if (m.getParameterCount() == 0
          && !Modifier.isStatic(m.getModifiers())
          && m.getReturnType() != void.class
          && !m.isSynthetic()
          && !m.getName().equals("toString")
          && !m.getName().equals("hashCode")) {
        m.setAccessible(true);
        getters.add(m);
      }
    }
    getters.sort(Comparator.comparing(Method::getName));

    for (int i = 0; i < n; i++) {
      line("## options " + i);
      Random rnd = new Random(seed + i);
      CompilerOptions o = new CompilerOptions();
      int count = i == 0 ? 0 : 1 + rnd.nextInt(i < 50 ? 4 : 25);
      for (int k = 0; k < count; k++) {
        Op op = OPS.get(rnd.nextInt(OPS.size()));
        String arg = op.args()[rnd.nextInt(op.args().length)];
        line("op " + op.name() + " " + arg);
        try {
          op.apply().accept(o, arg);
          line("ok");
        } catch (RuntimeException e) {
          line(exc(e));
        }
      }
      try {
        line("toString " + mask(o.toString()));
      } catch (RuntimeException e) {
        line("toString " + exc(e));
      }
      for (Method m : getters) {
        String v;
        try {
          v = fmt(m.invoke(o));
        } catch (InvocationTargetException e) {
          v = exc(e.getCause());
        }
        line("get " + m.getName() + " " + v);
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // deps

  /** An ErrorManager that records every report. */
  static final class Collector extends BasicErrorManager {
    final List<String> reports = new ArrayList<>();

    @Override
    public void report(CheckLevel level, JSError error) {
      reports.add(level + " " + error);
      super.report(level, error);
    }

    @Override
    public void println(CheckLevel level, JSError error) {}

    @Override
    protected void printSummary() {}
  }

  static String depInfo(DependencyInfo d) {
    return "provides="
        + d.getProvides()
        + " requires="
        + d.getRequires()
        + " typeRequires="
        + d.getTypeRequires()
        + " loadFlags="
        + d.getLoadFlags()
        + " externs="
        + d.getHasExternsAnnotation()
        + " nocompile="
        + d.getHasNoCompileAnnotation()
        + " es6="
        + d.isEs6Module()
        + " goog="
        + d.isGoogModule();
  }

  static void deps(String[] args) throws Exception {
    Path root = Paths.get(args[2]);
    List<String> files = new ArrayList<>();
    for (int i = 3; i < args.length; i++) {
      if (args[i].startsWith("--list=")) {
        files.addAll(Files.readAllLines(Paths.get(args[i].substring(7)), StandardCharsets.UTF_8));
      } else {
        files.add(args[i]);
      }
    }
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    for (String rel : files) {
      if (rel.isEmpty()) {
        continue;
      }
      line("## file " + rel);
      depsOf(rel, new String(Files.readAllBytes(root.resolve(rel)), StandardCharsets.UTF_8));
    }
  }

  /** Lines (some spanning several lines) that parse in a script or in a module. */
  static final String[] GEN_LINES = {
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
  };

  /** Fragments that do not parse. */
  static final String[] GEN_BAD = {
    "let 1x = ;", "/* start", " * @externs", "}", "if (x) {", "export\u0301 = 1;"
  };

  /** Random file contents over GEN_LINES and GEN_BAD (the same sequence as the Go generator). */
  static String genCode(Random rnd) {
    int n = rnd.nextInt(12);
    StringBuilder sb = new StringBuilder();
    String eol = rnd.nextInt(8) == 0 ? "\r\n" : "\n";
    for (int k = 0; k < n; k++) {
      String l =
          rnd.nextInt(25) == 0
              ? GEN_BAD[rnd.nextInt(GEN_BAD.length)]
              : GEN_LINES[rnd.nextInt(GEN_LINES.length)];
      sb.append(l.replace("\n", eol));
      if (k + 1 < n || rnd.nextBoolean()) {
        sb.append(eol);
      }
    }
    return sb.toString();
  }

  static void depsgen(long seed, int n) {
    Random rnd = new Random(seed);
    for (int i = 0; i < n; i++) {
      String rel = "gen" + i + ".js";
      String code = genCode(rnd);
      line("## gen " + rel + " " + code.replace("\n", "|"));
      depsOf(rel, code);
    }
  }

  static void depsOf(String rel, String code) {
    {
      for (int mode = 0; mode < 2; mode++) {
        Collector em = new Collector();
        JsFileRegexParser p = new JsFileRegexParser(em);
        if (mode == 1) {
          p.setShortcutMode(true);
          p.setIncludeGoogBase(true);
        }
        String tag = mode == 0 ? "regex " : "regexgoog ";
        try {
          DependencyInfo d = p.parseFile(rel, rel, code);
          line(tag + d + " succeeded=" + p.didParseSucceed());
        } catch (RuntimeException e) {
          line(tag + exc(e));
        }
        for (String r : em.reports) {
          line(" report " + r);
        }
      }
      for (int mode = 0; mode < 2; mode++) {
        Collector em = new Collector();
        Compiler compiler = new Compiler();
        compiler.setErrorManager(em);
        compiler.initOptions(new CompilerOptions());
        if (mode == 1) {
          compiler.setPreferRegexParser(true);
        }
        String tag = mode == 0 ? "ast " : "lazy ";
        CompilerInput input = new CompilerInput(SourceFile.fromCode(rel, code));
        input.setCompiler(compiler);
        try {
          line(tag + depInfo(input));
          line(tag + "known provides=" + input.getKnownProvides() + " requires=" + input.getKnownRequires());
        } catch (RuntimeException e) {
          line(tag + exc(e));
        }
        for (String r : em.reports) {
          if (r.startsWith("ERROR ")) {
            line(" report " + r);
          }
        }
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // strings

  static void strings(long seed, int n) throws Exception {
    Method parseJsString = JsFileLineParser.class.getDeclaredMethod("parseJsString", String.class);
    Method parseJsStringArray =
        JsFileLineParser.class.getDeclaredMethod("parseJsStringArray", String.class);
    Method parseJsStringMap = JsFileLineParser.class.getDeclaredMethod("parseJsStringMap", String.class);
    parseJsString.setAccessible(true);
    parseJsStringArray.setAccessible(true);
    parseJsStringMap.setAccessible(true);
    JsFileRegexParser p = new JsFileRegexParser(new Collector());
    Random rnd = new Random(seed);
    char[] alphabet = {'\'', '"', '\\', ',', ':', '[', ']', '{', '}', 'a', 'b', ' ', '\t', '\u00e9'};
    String[] fixed = {
      "'a'", "\"a\"", " 'a' ", "'a\\'b'", "'a\\\\'", "'a\\', 'b'", "['a', 'b']", "[]", "[ ]",
      "['a',]", "['a' 'b']", "['a\\', ']", "{'a': 'b'}", "{}", "{'a':'b','c':\"d\"}", "{'a'}",
      "{'a':}", "{'a': 'b',}", "  {'a': 'b'}  ", "{", "}", "'", "", "['a\\', 'b']", "]a[",
    };
    for (int i = 0; i < fixed.length + n; i++) {
      String s;
      if (i < fixed.length) {
        s = fixed[i];
      } else {
        int len = rnd.nextInt(14);
        StringBuilder sb = new StringBuilder();
        if (rnd.nextInt(3) == 0) {
          sb.append(rnd.nextBoolean() ? '[' : '{');
        }
        for (int k = 0; k < len; k++) {
          sb.append(alphabet[rnd.nextInt(alphabet.length)]);
        }
        if (rnd.nextInt(3) == 0) {
          sb.append(rnd.nextBoolean() ? ']' : '}');
        }
        s = sb.toString();
      }
      line("s " + s);
      for (Method m : new Method[] {parseJsString, parseJsStringArray, parseJsStringMap}) {
        String r;
        try {
          r = fmt(m.invoke(p, s));
        } catch (InvocationTargetException e) {
          r = exc(e.getCause());
        }
        line(" " + m.getName() + " " + r);
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // chunks

  record InputSpec(int chunk, String name, SourceKind kind, String code) {}

  record GraphSpec(
      List<List<Integer>> deps,
      boolean explicitWeak,
      List<InputSpec> inputs,
      List<int[]> queries,
      List<String> manage) {}

  static GraphSpec randomGraph(Random rnd) {
    int nChunks = 1 + rnd.nextInt(7);
    List<List<Integer>> deps = new ArrayList<>();
    for (int c = 0; c < nChunks; c++) {
      List<Integer> d = new ArrayList<>();
      if (c > 0) {
        int k = 1 + rnd.nextInt(Math.min(c, 3));
        for (int j = 0; j < k; j++) {
          int dep = rnd.nextInt(c);
          if (!d.contains(dep)) {
            d.add(dep);
          }
        }
      }
      deps.add(d);
    }
    // Rarely: chunks out of dependency order (ChunkDependenceException).
    if (nChunks > 2 && rnd.nextInt(25) == 0) {
      deps.get(1).add(nChunks - 1);
    }
    boolean explicitWeak = rnd.nextInt(8) == 0;
    List<InputSpec> inputs = new ArrayList<>();
    int nInputs = rnd.nextInt(9);
    String[] syms = {"a", "b", "c", "d", "e", "goog"};
    for (int i = 0; i < nInputs; i++) {
      StringBuilder code = new StringBuilder();
      if (rnd.nextInt(10) == 0) {
        code.append("/** @fileoverview\n * @provideGoog */\nvar COMPILED = false;\n");
      }
      int np = rnd.nextInt(3);
      for (int k = 0; k < np; k++) {
        code.append("goog.provide('").append(syms[rnd.nextInt(5)]).append(i % 3).append("');\n");
      }
      int nr = rnd.nextInt(3);
      for (int k = 0; k < nr; k++) {
        String kind = rnd.nextInt(4) == 0 ? "requireType" : "require";
        code.append("goog.").append(kind).append("('").append(syms[rnd.nextInt(5)]).append(rnd.nextInt(3)).append("');\n");
      }
      if (rnd.nextInt(8) == 0) {
        code.append("export {};\n");
      }
      int chunk = explicitWeak && rnd.nextInt(4) == 0 ? nChunks : rnd.nextInt(nChunks);
      SourceKind kind = rnd.nextInt(6) == 0 ? SourceKind.WEAK : SourceKind.STRONG;
      if (explicitWeak) {
        kind = chunk == nChunks ? SourceKind.WEAK : SourceKind.STRONG;
        if (rnd.nextInt(10) == 0) {
          kind = kind == SourceKind.WEAK ? SourceKind.STRONG : SourceKind.WEAK;
        }
      }
      inputs.add(new InputSpec(chunk, "f" + i + ".js", kind, code.toString()));
    }
    List<int[]> queries = new ArrayList<>();
    int total = nChunks + 1;
    for (int q = 0; q < 4; q++) {
      int parent = rnd.nextInt(total);
      int bits = rnd.nextInt(1 << total);
      queries.add(new int[] {parent, bits});
    }
    List<String> manage = new ArrayList<>();
    manage.add("sortOnly");
    manage.add("none");
    String[] eps = {"a0", "b1", "c2", "file:f0.js", "file:f3.js", "chunk:" + "c0" + ":a0", "zz"};
    for (String kind : new String[] {"pruneLegacy", "prune", "pruneAllow"}) {
      StringBuilder sb = new StringBuilder(kind).append(':');
      int ne = rnd.nextInt(3);
      for (int k = 0; k < ne; k++) {
        if (k > 0) {
          sb.append(',');
        }
        sb.append(eps[rnd.nextInt(eps.length)]);
      }
      manage.add(sb.toString());
    }
    return new GraphSpec(deps, explicitWeak, inputs, queries, manage);
  }

  static List<JSChunk> build(GraphSpec g) {
    List<JSChunk> chunks = new ArrayList<>();
    for (int c = 0; c < g.deps().size(); c++) {
      chunks.add(new JSChunk("c" + c));
    }
    for (int c = 0; c < g.deps().size(); c++) {
      for (int d : g.deps().get(c)) {
        chunks.get(c).addDependency(chunks.get(d));
      }
    }
    if (g.explicitWeak()) {
      JSChunk weak = new JSChunk(JSChunk.WEAK_CHUNK_NAME);
      for (JSChunk c : chunks) {
        weak.addDependency(c);
      }
      chunks.add(weak);
    }
    for (InputSpec in : g.inputs()) {
      chunks.get(in.chunk()).add(SourceFile.fromCode(in.name(), in.code(), in.kind()));
    }
    return chunks;
  }

  static String names(Iterable<JSChunk> chunks) {
    List<String> l = new ArrayList<>();
    for (JSChunk c : chunks) {
      l.add(c == null ? "null" : c.getName());
    }
    return l.toString();
  }

  static String chunkContents(JSChunkGraph graph) {
    StringBuilder sb = new StringBuilder();
    for (JSChunk c : graph.getAllChunks()) {
      sb.append(c.getName()).append('[');
      boolean first = true;
      for (CompilerInput in : c.getInputs()) {
        if (!first) {
          sb.append(',');
        }
        first = false;
        sb.append(in.getName()).append(':').append(in.getSourceFile().getKind());
        sb.append(':').append(in.getChunk() == null ? "null" : in.getChunk().getName());
      }
      sb.append("] ");
    }
    return sb.toString();
  }

  static List<ModuleIdentifier> manageEntryPoints(String a) {
    List<ModuleIdentifier> l = new ArrayList<>();
    for (String s : list(a)) {
      if (s.startsWith("file:")) {
        l.add(ModuleIdentifier.forFile(s.substring(5)));
      } else if (s.startsWith("chunk:")) {
        l.add(ModuleIdentifier.forFlagValue(s.substring(6)));
      } else {
        l.add(ModuleIdentifier.forClosure(s));
      }
    }
    return l;
  }

  static void chunks(long seed, int n) {
    Random rnd = new Random(seed);
    for (int gi = 0; gi < n; gi++) {
      GraphSpec g = randomGraph(rnd);
      line("## graph " + gi);
      for (int c = 0; c < g.deps().size(); c++) {
        line("chunk c" + c + " deps=" + g.deps().get(c));
      }
      line("explicitWeak " + g.explicitWeak());
      for (InputSpec in : g.inputs()) {
        line("input " + in.chunk() + " " + in.name() + " " + in.kind() + " " + in.code().replace("\n", "|"));
      }
      for (int[] q : g.queries()) {
        line("query " + q[0] + " " + q[1]);
      }
      for (String m : g.manage()) {
        line("manage " + m);
      }
      line("--");
      JSChunkGraph graph;
      List<JSChunk> chunks = build(g);
      try {
        graph = new JSChunkGraph(chunks);
      } catch (RuntimeException e) {
        line("construct " + exc(e));
        continue;
      }
      List<JSChunk> all = new ArrayList<>();
      graph.getAllChunks().forEach(all::add);
      line("count " + graph.getChunkCount() + " inputs " + graph.getInputCount());
      line("contents " + chunkContents(graph));
      for (JSChunk c : all) {
        StringBuilder sb = new StringBuilder("chunk " + c.getName() + " depth=" + c.getDepth() + " index=" + c.getIndex());
        sb.append(" dependsOn=");
        for (JSChunk d : all) {
          sb.append(graph.dependsOn(c, d) ? '1' : '0');
        }
        sb.append(" transitive=").append(names(graph.getTransitiveDepsDeepestFirst(c)));
        sb.append(" sortedDeps=").append(c.getSortedDependencyNames());
        sb.append(" all=").append(names(c.getAllDependencies()));
        sb.append(" requires=").append(c.getRequires());
        line(sb.toString());
      }
      for (JSChunk a : all) {
        StringBuilder sb = new StringBuilder("dcd " + a.getName());
        for (JSChunk b : all) {
          JSChunk x = graph.getDeepestCommonDependency(a, b);
          JSChunk y = graph.getDeepestCommonDependencyInclusive(a, b);
          sb.append(' ').append(x == null ? "-" : x.getName()).append('/').append(y == null ? "-" : y.getName());
        }
        line(sb.toString());
      }
      line("dcdAll " + graph.getDeepestCommonDependencyInclusive(all).getName());
      for (int[] q : g.queries()) {
        BitSet bs = new BitSet();
        for (int i = 0; i < all.size(); i++) {
          if ((q[1] & (1 << i)) != 0) {
            bs.set(i);
          }
        }
        String r;
        try {
          r = graph.getSmallestCoveringSubtree(all.get(q[0] % all.size()), bs).getName();
        } catch (RuntimeException e) {
          r = exc(e);
        }
        line("cover " + q[0] + " " + bs + " " + r);
      }
      try {
        line("root " + graph.getRootChunk().getName());
      } catch (RuntimeException e) {
        line("root " + exc(e));
      }
      line("json " + graph.toJson());
      for (String m : g.manage()) {
        JSChunkGraph fresh = new JSChunkGraph(build(g));
        Collector em = new Collector();
        Compiler compiler = new Compiler();
        compiler.setErrorManager(em);
        compiler.initOptions(new CompilerOptions());
        for (CompilerInput in : fresh.getAllInputs()) {
          in.setCompiler(compiler);
        }
        try {
          DependencyOptions opts;
          if (m.equals("sortOnly")) {
            opts = DependencyOptions.sortOnly();
          } else if (m.equals("none")) {
            opts = DependencyOptions.none();
          } else {
            String kind = m.substring(0, m.indexOf(':'));
            List<ModuleIdentifier> eps = manageEntryPoints(m.substring(m.indexOf(':') + 1));
            opts =
                switch (kind) {
                  case "pruneLegacy" -> DependencyOptions.pruneLegacyForEntryPoints(eps);
                  case "prune" -> DependencyOptions.pruneForEntryPoints(eps);
                  default -> DependencyOptions.pruneAllowNoEntryPoints(eps);
                };
          }
          List<String> res = new ArrayList<>();
          for (CompilerInput in : fresh.manageDependencies(compiler, opts).orderedInputs()) {
            res.add(in.getName());
          }
          line("managed " + m + " " + res);
          line(" contents " + chunkContents(fresh));
        } catch (Exception e) {
          line("managed " + m + " " + exc(e));
        }
        for (String r : em.reports) {
          if (r.startsWith("ERROR ") || r.contains("WEAK")) {
            line(" report " + r);
          }
        }
      }
    }
  }
}
