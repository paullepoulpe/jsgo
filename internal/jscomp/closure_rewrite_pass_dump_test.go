package jscomp

// Go-only: test-helper — the Go side of tools/oracle/java/ClosureRewritePassDump.java, the
// differential oracle for the w5-closure-modules passes (ClosureRewriteModule, ScopedAliases,
// CheckClosureImports, RewriteGoogJsImports, ModuleMapCreator and its module processors,
// ProcessDefines, ClosureCodeRemoval, GenerateExports, ExternExportsPass). It runs the same passes
// the same way (a fresh Compiler per pass, Compiler.init + Compiler.parse, the synthesized externs
// input, GatherModuleMetadata and ModuleMapCreator first, as DefaultPassConfig does) and prints the
// same text, which must match the Java output byte for byte.
//
// The checked-in golden set is testdata/closure_rewrite: src/custom (hand-written goog.module,
// goog.loadModule, goog.scope, @define, @export, ES module and goog.js inputs; printed with --full),
// src/javatests (the JS literals of the ported Java unit tests, one file each; sha only) and
// src/prog (one multi-file program, --program --full), produced by
//
//	java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.ClosureRewritePassDump \
//	  [--full] [--program] internal/jscomp/testdata/closure_rewrite/<golden>.txt.gz \
//	  internal/jscomp/testdata/closure_rewrite/src --list=internal/jscomp/testdata/closure_rewrite/<list>.txt
//
// The corpus run (env-gated): CLOSURE_REWRITE_ORACLE=<the Java .gz> CLOSURE_REWRITE_ROOT=<corpus
// root> CLOSURE_REWRITE_FILES=<list file> [CLOSURE_REWRITE_PROGRAM=1, as --program]
// [CLOSURE_REWRITE_FULL=1, as --full] [CLOSURE_REWRITE_DIFF_DIR=<dir> to write the differing
// sections].

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type closureRewriteDump struct {
	sb   *strings.Builder
	full bool
}

// closureRewriteRun is the Run interface of the Java tool.
type closureRewriteRun func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string)

func closureRewriteNewCompiler(inputs []*SourceFile) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	compiler.Init([]*SourceFile{}, inputs, options)
	compiler.Parse()
	return compiler
}

// closureRewriteExceptionText is the Java tool's exceptionText: Throwable.toString(), except
// that a NullPointerException prints its class only.
func closureRewriteExceptionText(e jutil.IRuntimeException) string {
	if _, ok := e.(jutil.INullPointerException); ok {
		return "java.lang.NullPointerException"
	}
	return e.Error()
}

func closureRewriteGatherModuleMetadata(compiler *Compiler, externs rhino.Node, js rhino.Node) {
	// Force creation of the synthetic input so that we create metadata for it (as
	// DefaultPassConfig.gatherModuleMetadataPass does).
	compiler.GetSynthesizedExternsInput()
	NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(externs, js)
}

func closureRewriteCreateModuleMap(compiler *Compiler, externs rhino.Node, js rhino.Node) {
	NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).Process(externs, js)
}

// closureRewriteDefineReplacements is the Java tool's define map.
func closureRewriteDefineReplacements() *jutil.LinkedHashMap[string, rhino.Node] {
	m := jutil.NewLinkedHashMap[string, rhino.Node]()
	m.Put("COMPILED", rhino.NewNode(rhino.Token_TRUE))
	m.Put("goog.DEBUG", rhino.NewNode(rhino.Token_FALSE))
	m.Put("goog.LOCALE", rhino.IR_String("de"))
	m.Put("goog.TRUSTED_SITE", rhino.NewNode(rhino.Token_FALSE))
	m.Put("goog.dom.ASSUME_STANDARDS_MODE", rhino.NewNode(rhino.Token_TRUE))
	m.Put("goog.userAgent.ASSUME_ANY_VERSION", rhino.NewNode(rhino.Token_TRUE))
	m.Put("goog.json.USE_NATIVE_JSON", rhino.NewNode(rhino.Token_TRUE))
	m.Put("goog.asserts.ENABLE_ASSERTS", rhino.NewNode(rhino.Token_FALSE))
	m.Put("goog.FEATURESET_YEAR", rhino.IR_Number(2021))
	m.Put("UNKNOWN_DEFINE_X", rhino.IR_Number(1))
	return m
}

func closureRewriteKeys[V any](m *jutil.LinkedHashMap[string, V]) string {
	return "[" + strings.Join(m.KeySlice(), ", ") + "]"
}

func closureRewriteStr(s *string) string { return jutil.ToJavaString(s) }

func closureRewritePath(p *ModuleLoader_ModulePath) string {
	if p == nil {
		return "null"
	}
	return p.String()
}

func closureRewriteMetadata(md *ModuleMetadataMap_ModuleMetadata) string {
	if md == nil {
		return "metadata=null"
	}
	return "type=" + md.ModuleType().String() + " goog=" + md.GoogNamespaces().String() + " mpath=" + closureRewritePath(md.Path())
}

func closureRewriteModule(m *Module) string {
	return "path=" + closureRewritePath(m.Path()) +
		" " + closureRewriteMetadata(m.Metadata()) +
		" closureNamespace=" + closureRewriteStr(m.ClosureNamespace()) +
		" namespace=" + closureRewriteKeys(m.Namespace()) +
		" boundNames=" + closureRewriteKeys(m.BoundNames()) +
		" localExports=" + closureRewriteKeys(m.LocalNameToLocalExport())
}

func closureRewriteBinding(b *Binding) string {
	boundName, mutated := "-", "-"
	if !b.IsModuleNamespace() {
		boundName = closureRewriteStr(b.BoundName())
		mutated = strconv.FormatBool(b.IsMutated())
	}
	export := "null"
	if b.OriginatingExport() != nil {
		export = closureRewriteExport(b.OriginatingExport())
	}
	return "{" + closureRewriteMetadata(b.Metadata()) +
		" source=" + closureModulesNodeRef(b.SourceNode()) +
		" nsObj=" + strconv.FormatBool(b.IsModuleNamespace()) +
		" closureNs=" + closureRewriteStr(b.ClosureNamespace()) +
		" createdBy=" + b.CreatedBy().String() +
		" boundName=" + boundName +
		" mutated=" + mutated +
		" export=" + export + "}"
}

func closureRewriteExport(x *Export) string {
	return "[" + closureRewriteStr(x.ExportName()) +
		" req=" + closureRewriteStr(x.ModuleRequest()) +
		" import=" + closureRewriteStr(x.ImportName()) +
		" local=" + closureRewriteStr(x.LocalName()) +
		" path=" + closureRewritePath(x.ModulePath()) +
		" exportNode=" + closureModulesNodeRef(x.ExportNode()) +
		" nameNode=" + closureModulesNodeRef(x.NameNode()) +
		" closureNs=" + closureRewriteStr(x.ClosureNamespace()) +
		" mutated=" + strconv.FormatBool(x.Mutated()) + "]"
}

func closureRewriteModuleMap(m *ModuleMap, extra *[]string) {
	if m == nil {
		*extra = append(*extra, "moduleMap null")
		return
	}
	for k, v := range m.GetModulesByPath().All() {
		*extra = append(*extra, "module "+k+" "+closureRewriteModule(v))
		for name, b := range v.Namespace().All() {
			*extra = append(*extra, "  ns "+name+" "+closureRewriteBinding(b))
		}
		for name, b := range v.BoundNames().All() {
			*extra = append(*extra, "  bound "+name+" "+closureRewriteBinding(b))
		}
		for name, e := range v.LocalNameToLocalExport().All() {
			*extra = append(*extra, "  local "+name+" "+closureRewriteExport(e))
		}
	}
	for k, v := range m.GetModulesByClosureNamespace().All() {
		*extra = append(*extra, "closureModule "+k+" "+closureRewriteModule(v))
	}
}

func (d *closureRewriteDump) tree(label string, before string, after string) {
	if after == before {
		closureModulesLine(d.sb, label+" unchanged")
		return
	}
	var tree strings.Builder
	for _, l := range strings.Split(after, "\n") {
		closureModulesLine(&tree, l)
	}
	closureModulesLine(d.sb, label+" changed "+closureModulesSha256(tree.String()))
	if d.full {
		d.sb.WriteString(tree.String())
	}
}

func (d *closureRewriteDump) pass(inputs []*SourceFile, parseDiags []string, name string, run closureRewriteRun) {
	sb := d.sb
	closureModulesLine(sb, "## pass "+name)
	compiler := closureRewriteNewCompiler(inputs)
	externs := compiler.GetExternsRoot()
	js := compiler.GetJsRoot()
	before := js.ToStringTree()
	externsBefore := externs.ToStringTree()
	var extra []string
	var exception *string
	jutil.Try(func() {
		run(compiler, externs, js, &extra)
	}, func(e jutil.IRuntimeException) {
		s := "exception " + closureRewriteExceptionText(e)
		exception = &s
	})
	remaining := append([]string(nil), parseDiags...)
	for _, diag := range closureModulesDiagnostics(compiler) {
		removed := false
		for i, r := range remaining {
			if r == diag {
				remaining = append(remaining[:i], remaining[i+1:]...)
				removed = true
				break
			}
		}
		if !removed {
			closureModulesLine(sb, diag)
		}
	}
	if exception != nil {
		closureModulesLine(sb, *exception)
	}
	d.tree("tree", before, js.ToStringTree())
	d.tree("externs", externsBefore, externs.ToStringTree())
	for _, e := range extra {
		closureModulesLine(sb, e)
	}
}

func (d *closureRewriteDump) dump(header string, inputs []*SourceFile) {
	sb := d.sb
	closureModulesLine(sb, header)
	var parsed *Compiler
	var parseException *string
	jutil.Try(func() {
		parsed = closureRewriteNewCompiler(inputs)
	}, func(e jutil.IRuntimeException) {
		s := "# parse exception " + closureRewriteExceptionText(e)
		parseException = &s
	})
	if parseException != nil {
		closureModulesLine(sb, *parseException)
		return
	}
	parseDiags := closureModulesDiagnostics(parsed)
	if parsed.HasErrors() {
		closureModulesLine(sb, "# parse errors")
		for _, diag := range parseDiags {
			closureModulesLine(sb, diag)
		}
		return
	}
	closureModulesLine(sb, "# parse ok")
	for _, diag := range parseDiags {
		closureModulesLine(sb, diag)
	}

	d.pass(inputs, parseDiags, "moduleMap", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		closureRewriteGatherModuleMetadata(compiler, externs, js)
		closureRewriteCreateModuleMap(compiler, externs, js)
		closureRewriteModuleMap(compiler.GetModuleMap(), extra)
	})
	d.pass(inputs, parseDiags, "checkClosureImports", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		closureRewriteGatherModuleMetadata(compiler, externs, js)
		closureRewriteCreateModuleMap(compiler, externs, js)
		NewCheckClosureImports(compiler, compiler.GetModuleMetadataMap()).Process(externs, js)
	})
	d.pass(inputs, parseDiags, "rewriteGoogJsImports", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		closureRewriteGatherModuleMetadata(compiler, externs, js)
		closureRewriteCreateModuleMap(compiler, externs, js)
		NewRewriteGoogJsImports(compiler, RewriteGoogJsImports_Mode_LINT_AND_REWRITE, compiler.GetModuleMap()).
			Process(externs, js)
		closureRewriteModuleMap(compiler.GetModuleMap(), extra)
	})
	d.pass(inputs, parseDiags, "closureGoogScopeAliases", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		closureRewriteGatherModuleMetadata(compiler, externs, js)
		closureRewriteCreateModuleMap(compiler, externs, js)
		ScopedAliases_Builder_(compiler).
			SetPreprocessorSymbolTable(nil).
			SetModuleMetadataMap(compiler.GetModuleMetadataMap()).
			SetInvalidModuleGetHandling(ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME).
			Build().
			Process(externs, js)
	})
	d.pass(inputs, parseDiags, "closureGoogScopeAliasesForIjs", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		closureRewriteGatherModuleMetadata(compiler, externs, js)
		closureRewriteCreateModuleMap(compiler, externs, js)
		ScopedAliases_Builder_(compiler).Build().Process(externs, js)
	})
	d.pass(inputs, parseDiags, "closureRewriteModule", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		closureRewriteGatherModuleMetadata(compiler, externs, js)
		closureRewriteCreateModuleMap(compiler, externs, js)
		NewClosureRewriteModule(compiler, nil, compiler.GetTopScope()).Process(externs, js)
	})
	d.pass(inputs, parseDiags, "closureGoogScopeAliases+closureRewriteModule", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		closureRewriteGatherModuleMetadata(compiler, externs, js)
		closureRewriteCreateModuleMap(compiler, externs, js)
		ScopedAliases_Builder_(compiler).
			SetModuleMetadataMap(compiler.GetModuleMetadataMap()).
			SetInvalidModuleGetHandling(ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME).
			Build().
			Process(externs, js)
		NewClosureRewriteModule(compiler, nil, compiler.GetTopScope()).Process(externs, js)
	})
	for _, mode := range ProcessDefines_Mode_values() {
		d.pass(inputs, parseDiags, "processDefines_"+mode.String(), func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
			NewProcessDefines_Builder(compiler).
				PutReplacements(closureRewriteDefineReplacements()).
				SetMode(mode).
				SetRecognizeClosureDefines(true).
				SetEnableZonesDefineName(compiler.GetOptions().GetEnableZonesDefineName()).
				SetZoneInputPattern(compiler.GetOptions().GetZoneInputPattern()).
				SetUnknownDefinesToIgnore(compiler.GetOptions().GetUnknownDefinesToIgnore()).
				Build().
				Process(externs, js)
			*extra = append(*extra, "defineNames "+compiler.GetDefineNames().String())
		})
	}
	d.pass(inputs, parseDiags, "closureCodeRemoval", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		NewClosureCodeRemoval(compiler, true, true).Process(externs, js)
	})
	for _, allowNonGlobalExports := range []bool{false, true} {
		d.pass(inputs, parseDiags, "generateExports nonGlobal="+strconv.FormatBool(allowNonGlobalExports), func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
			g := NewGenerateExports(compiler, allowNonGlobalExports, jutil.StringPtr("goog.exportSymbol"), jutil.StringPtr("goog.exportProperty"))
			g.Process(externs, js)
			*extra = append(*extra, "exported "+g.GetExportedVariableNames().String())
		})
	}
	d.pass(inputs, parseDiags, "externExports", func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
		NewExternExportsPass(compiler).Process(externs, js)
		result := compiler.GetResult().ExternExport
		if result == nil {
			*extra = append(*extra, "externExport null")
		} else {
			for _, l := range strings.Split(*result, "\n") {
				*extra = append(*extra, "| "+l)
			}
		}
	})
}

// closureRewriteSplit splits a dump into its "## file" / "## program" sections.
func closureRewriteSplit(dump string) map[string]string {
	out := map[string]string{}
	var cur string
	var sb strings.Builder
	flush := func() {
		if cur != "" {
			out[cur] = sb.String()
		}
		sb.Reset()
	}
	for _, l := range strings.SplitAfter(dump, "\n") {
		if strings.HasPrefix(l, "## file ") || strings.HasPrefix(l, "## program ") {
			flush()
			cur = strings.TrimSuffix(l, "\n")
		}
		sb.WriteString(l)
	}
	flush()
	return out
}

func closureRewriteReport(t *testing.T, key string, got string, want string, diffDir string, mismatches int) {
	t.Helper()
	if mismatches <= 20 {
		gl := strings.Split(got, "\n")
		wl := strings.Split(want, "\n")
		for k := 0; k < len(gl) || k < len(wl); k++ {
			var g, w string
			if k < len(gl) {
				g = gl[k]
			}
			if k < len(wl) {
				w = wl[k]
			}
			if g != w {
				t.Errorf("%s: line %d differs:\n  java: %s\n  go  : %s", key, k+1, w, g)
				break
			}
		}
	}
	if diffDir != "" {
		p := filepath.Join(diffDir, strings.NewReplacer("## file ", "", "## program ", "program-", " ", "_").Replace(key))
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		_ = os.WriteFile(p+".go.txt", []byte(got), 0o644)
		_ = os.WriteFile(p+".java.txt", []byte(want), 0o644)
	}
}

// closureRewriteRunDump dumps one file or one program, turning a Go panic into a line.
func closureRewriteRunDump(full bool, header string, inputs []*SourceFile) string {
	d := &closureRewriteDump{sb: &strings.Builder{}, full: full}
	func() {
		defer func() {
			if r := recover(); r != nil {
				d.sb.WriteString(fmt.Sprintf("GO PANIC: %v\n", r))
			}
		}()
		d.dump(header, inputs)
	}()
	return d.sb.String()
}

// closureRewriteCompare dumps every file (or the program) and compares it with its Java section.
func closureRewriteCompare(t *testing.T, root string, files []string, javaDump string, full bool, program bool, diffDir string) {
	t.Helper()
	java := closureRewriteSplit(javaDump)
	var inputs []*SourceFile
	for _, rel := range files {
		// Go-only: test-helper — Files.readString(UTF_8) is strict; the inputs are valid UTF-8.
		srcBytes, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		inputs = append(inputs, SourceFile_FromCode(rel, string(srcBytes)))
	}
	if program {
		header := "## program " + files[0] + " (+" + strconv.Itoa(len(files)-1) + ")"
		if got, want := closureRewriteRunDump(full, header, inputs), java[header]; got != want {
			closureRewriteReport(t, header, got, want, diffDir, 1)
		}
		return
	}
	mismatches := 0
	for i, rel := range files {
		header := "## file " + rel
		got := closureRewriteRunDump(full, header, []*SourceFile{inputs[i]})
		if want := java[header]; got != want {
			mismatches++
			closureRewriteReport(t, header, got, want, diffDir, mismatches)
		}
	}
	if mismatches > 0 {
		t.Errorf("%d of %d files differ", mismatches, len(files))
	}
}

func TestClosureRewritePassDumpGolden(t *testing.T) {
	dir := filepath.Join("testdata", "closure_rewrite")
	diffDir := os.Getenv("CLOSURE_REWRITE_DIFF_DIR")
	for _, c := range []struct {
		golden, list  string
		full, program bool
	}{
		{"custom.txt.gz", "custom.txt", true, false},
		{"javatests.txt.gz", "javatests.txt", false, false},
		{"prog.txt.gz", "prog.txt", true, true},
	} {
		t.Run(c.list, func(t *testing.T) {
			files := closureModulesReadList(t, filepath.Join(dir, c.list))
			closureRewriteCompare(t, filepath.Join(dir, "src"), files, closureModulesReadGz(t, filepath.Join(dir, c.golden)), c.full, c.program, diffDir)
		})
	}
}

func TestClosureRewritePassDumpCorpus(t *testing.T) {
	oracle := os.Getenv("CLOSURE_REWRITE_ORACLE")
	if oracle == "" {
		t.Skip("set CLOSURE_REWRITE_ORACLE, CLOSURE_REWRITE_ROOT and CLOSURE_REWRITE_FILES")
	}
	files := closureModulesReadList(t, os.Getenv("CLOSURE_REWRITE_FILES"))
	closureRewriteCompare(t, os.Getenv("CLOSURE_REWRITE_ROOT"), files, closureModulesReadGz(t, oracle),
		os.Getenv("CLOSURE_REWRITE_FULL") == "1", os.Getenv("CLOSURE_REWRITE_PROGRAM") == "1", os.Getenv("CLOSURE_REWRITE_DIFF_DIR"))
}
