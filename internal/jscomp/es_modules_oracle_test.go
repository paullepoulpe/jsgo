package jscomp

// Go-only: test-helper — the Go side of tools/oracle/java/EsModulesDump.java, the differential
// oracle for the w5-es-modules passes (GatherModuleMetadata, ProcessCommonJSModules,
// FindModuleDependencies, Es6RewriteScriptsToModules, Es6RewriteModules, ModuleRenaming,
// Es6RewriteModulesToCommonJsModules, Es6RelativizeImportPaths, RewriteDynamicImports,
// ijs/ConvertToTypedInterface). It runs the same sections the same way (a fresh Compiler per
// section, Compiler.init + Compiler.parse, the prerequisites DefaultPassConfig runs) and prints the
// same text, which must match the Java output byte for byte.
//
// The checked-in golden set: testdata/es_modules/src/<rel> with testdata/es_modules/golden.txt.gz,
// produced by
//
//	java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.EsModulesDump \
//	  --full internal/jscomp/testdata/es_modules/golden.txt.gz \
//	  internal/jscomp/testdata/es_modules/src --list=internal/jscomp/testdata/es_modules/files.txt
//
// and the checked-in fixtures run: testdata/fixtures (every .js file, listed in
// testdata/es_modules/fixtures-files.txt) with testdata/es_modules/fixtures.txt.gz, produced by
//
//	java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.EsModulesDump \
//	  --max-size=2000000 internal/jscomp/testdata/es_modules/fixtures.txt.gz testdata/fixtures \
//	  --list=internal/jscomp/testdata/es_modules/fixtures-files.txt
//
// The corpus run (env-gated): ES_MODULES_ORACLE=<the Java .gz> ES_MODULES_ROOT=<corpus root>
// ES_MODULES_FILES=<list file> [ES_MODULES_SECTIONS=a,b,..., as --sections]
// [ES_MODULES_FULL=1, as --full] [ES_MODULES_MAX_SIZE=<B>, as --max-size]
// [ES_MODULES_DIFF_DIR=<dir> to write the differing file sections]. Several Java outputs (e.g. one
// per chunk of the file list) may be concatenated: they are gzip members. With
// ES_MODULES_SECTIONS, only those sections of the Java output are compared.

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var esModulesAllSections = []string{
	"gather",
	"gatherCjs",
	"rewriteCommonJs",
	"es6RewriteModules",
	"es6ToCjs",
	"relativize",
	"dynamicImports",
	"ijs",
}

type esModulesDump struct {
	sb       *strings.Builder
	full     bool
	sections []string
}

type esModulesRun func(compiler *Compiler, externs *rhino.Node, js *rhino.Node, extra *[]string)

func esModulesEscape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		if c == '\\' {
			b.WriteString(`\\`)
		} else if c >= 0x20 && c < 0x7f {
			b.WriteByte(byte(c))
		} else {
			fmt.Fprintf(&b, `\u%04x`, c)
		}
	}
	return b.String()
}

func esModulesLine(sb *strings.Builder, s string) {
	sb.WriteString(esModulesEscape(s))
	sb.WriteByte('\n')
}

func esModulesSha256(ascii string) string {
	d := sha256.Sum256([]byte(ascii))
	return fmt.Sprintf("%x", d[:8])
}

func esModulesNewCompiler(rel string, src string, setup func(*CompilerOptions)) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroup_ForType(ModuleLoader_INVALID_MODULE_PATH), CheckLevel_WARNING)
	setup(options)
	compiler.Init([]*SourceFile{}, []*SourceFile{SourceFile_FromCode(rel, src)}, options)
	compiler.Parse()
	return compiler
}

func esModulesDiag(level string, e *JSError) string {
	return "diag " + level + " " + e.GetType().key + " " + strconv.Itoa(e.GetLineno()) + ":" +
		strconv.Itoa(e.GetCharno()) + " \"" + e.GetDescription() + "\""
}

func esModulesDiagnostics(compiler *Compiler) []string {
	var out []string
	for _, e := range compiler.GetErrors() {
		out = append(out, esModulesDiag("ERROR", e))
	}
	for _, e := range compiler.GetWarnings() {
		out = append(out, esModulesDiag("WARNING", e))
	}
	return out
}

func esModulesNodeRef(n *rhino.Node) string {
	if n == nil {
		return "null"
	}
	return n.GetToken().String() + "@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

func esModulesNoSetup(*CompilerOptions) {}

func esModulesCjs(o *CompilerOptions) { o.SetProcessCommonJSModules(true) }

// esModulesGatherAndCreateModuleMap is the Java tool's gatherAndCreateModuleMap.
func esModulesGatherAndCreateModuleMap(compiler *Compiler, externs *rhino.Node, js *rhino.Node, processCommonJsModules bool) {
	compiler.GetSynthesizedExternsInput()
	NewGatherModuleMetadata(compiler, processCommonJsModules, ModuleLoader_ResolutionMode_BROWSER).Process(externs, js)
	NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).Process(externs, js)
}

func esModulesRewriteModules(compiler *Compiler, externs *rhino.Node, js *rhino.Node) {
	NewEs6RewriteModulesWithChunkOutputType(
		compiler,
		compiler.GetModuleMetadataMap(),
		compiler.GetModuleMap(),
		/* preprocessorSymbolTable= */ nil,
		/* globalTypedScope= */ nil,
		CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).
		Process(externs, js)
}

func esModulesPathString(p *ModuleLoader_ModulePath) string {
	if p == nil {
		return "null"
	}
	return p.String()
}

// esModulesMetadata prints the ModuleMetadataMap as the Java tool's metadata does.
func esModulesMetadata(m *ModuleMetadataMap, extra *[]string) {
	ids := map[*ModuleMetadataMap_ModuleMetadata]int{}
	var order []*ModuleMetadataMap_ModuleMetadata
	var number func(md *ModuleMetadataMap_ModuleMetadata)
	number = func(md *ModuleMetadataMap_ModuleMetadata) {
		if _, ok := ids[md]; ok {
			return
		}
		ids[md] = len(order)
		order = append(order, md)
		for _, n := range md.NestedModules() {
			number(n)
		}
	}
	for md := range m.GetAllModuleMetadata().All() {
		number(md)
	}
	for _, md := range order {
		var nested strings.Builder
		nested.WriteByte('[')
		for i, n := range md.NestedModules() {
			if i > 0 {
				nested.WriteString(", ")
			}
			nested.WriteString("#" + strconv.Itoa(ids[n]))
		}
		nested.WriteByte(']')
		*extra = append(*extra, "module #"+strconv.Itoa(ids[md])+
			" type="+md.ModuleType().String()+
			" root="+esModulesNodeRef(md.RootNode())+
			" path="+esModulesPathString(md.Path())+
			" usesClosure="+strconv.FormatBool(md.UsesClosure())+
			" testOnly="+strconv.FormatBool(md.IsTestOnly())+
			" ns="+md.GoogNamespaces().String()+
			" strong="+md.StronglyRequiredGoogNamespaces().String()+
			" dynamic="+md.DynamicallyRequiredGoogNamespaces().String()+
			" maybe="+md.MaybeRequiredGoogNamespaces().String()+
			" weak="+md.WeaklyRequiredGoogNamespaces().String()+
			" es6="+md.Es6ImportSpecifiers().String()+
			" toggles="+md.ReadToggles().String()+
			" nested="+nested.String())
	}
	for k, md := range m.GetModulesByPath().All() {
		*extra = append(*extra, "byPath "+k+" #"+strconv.Itoa(ids[md]))
	}
	for k, md := range m.GetModulesByGoogNamespace().All() {
		*extra = append(*extra, "byNs "+k+" #"+strconv.Itoa(ids[md]))
	}
}

func (d *esModulesDump) tree(what string, after string, before *string) {
	if before != nil && after == *before {
		esModulesLine(d.sb, what+" unchanged")
		return
	}
	var tree strings.Builder
	for _, l := range strings.Split(after, "\n") {
		esModulesLine(&tree, l)
	}
	esModulesLine(d.sb, what+" changed "+esModulesSha256(tree.String()))
	if d.full {
		d.sb.WriteString(tree.String())
	}
}

// pass is the Java tool's pass: parseDiags nil means the section parses differently and prints
// its own parse diagnostics.
func (d *esModulesDump) pass(rel string, src string, parseDiags []string, name string, setup func(*CompilerOptions), run esModulesRun) {
	sb := d.sb
	esModulesLine(sb, "## pass "+name)
	compiler := esModulesNewCompiler(rel, src, setup)
	if parseDiags == nil {
		parseDiags = esModulesDiagnostics(compiler)
		for _, diag := range parseDiags {
			esModulesLine(sb, "# parse "+diag)
		}
		if compiler.HasErrors() {
			return
		}
	}
	externs := compiler.GetExternsRoot()
	js := compiler.GetJsRoot()
	before := js.ToStringTree()
	externsBefore := externs.ToStringTree()
	var extra []string
	var exception *string
	jutil.Try(func() {
		run(compiler, externs, js, &extra)
	}, func(e jutil.IRuntimeException) {
		s := "exception " + e.Error()
		exception = &s
	})
	remaining := append([]string(nil), parseDiags...)
	for _, diag := range esModulesDiagnostics(compiler) {
		removed := false
		for i, r := range remaining {
			if r == diag {
				remaining = append(remaining[:i], remaining[i+1:]...)
				removed = true
				break
			}
		}
		if !removed {
			esModulesLine(sb, diag)
		}
	}
	if exception != nil {
		esModulesLine(sb, *exception)
	}
	d.tree("tree", js.ToStringTree(), &before)
	if externsAfter := externs.ToStringTree(); externsAfter != externsBefore {
		d.tree("externs", externsAfter, nil)
	}
	for _, e := range extra {
		esModulesLine(sb, e)
	}
}

func (d *esModulesDump) dumpFile(rel string, src string) {
	sb := d.sb
	esModulesLine(sb, "## file "+rel)
	parsed := esModulesNewCompiler(rel, src, esModulesNoSetup)
	parseDiags := esModulesDiagnostics(parsed)
	if parsed.HasErrors() {
		esModulesLine(sb, "# parse errors")
		for _, diag := range parseDiags {
			esModulesLine(sb, diag)
		}
		return
	}
	esModulesLine(sb, "# parse ok")
	for _, diag := range parseDiags {
		esModulesLine(sb, diag)
	}
	if parseDiags == nil {
		parseDiags = []string{}
	}

	for _, section := range d.sections {
		switch section {
		case "gather":
			d.pass(rel, src, parseDiags, section, esModulesNoSetup, func(compiler *Compiler, externs, js *rhino.Node, extra *[]string) {
				NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(externs, js)
				esModulesMetadata(compiler.GetModuleMetadataMap(), extra)
			})
		case "gatherCjs":
			d.pass(rel, src, nil, section, esModulesCjs, func(compiler *Compiler, externs, js *rhino.Node, extra *[]string) {
				for _, input := range compiler.GetChunkGraph().GetAllInputs() {
					*extra = append(*extra, "input "+input.GetName()+" "+input.GetJsModuleType().String())
				}
				NewGatherModuleMetadata(compiler, true, ModuleLoader_ResolutionMode_BROWSER).Process(externs, js)
				esModulesMetadata(compiler.GetModuleMetadataMap(), extra)
			})
		case "rewriteCommonJs":
			d.pass(rel, src, nil, section, esModulesCjs, func(compiler *Compiler, externs, js *rhino.Node, extra *[]string) {
				esModulesGatherAndCreateModuleMap(compiler, externs, js, true)
				NewProcessCommonJSModules(compiler).Process(externs, js)
			})
		case "es6RewriteModules":
			d.pass(rel, src, parseDiags, section, esModulesNoSetup, func(compiler *Compiler, externs, js *rhino.Node, extra *[]string) {
				compiler.GetSynthesizedExternsInput()
				NewEs6RewriteScriptsToModules(compiler).Process(externs, js)
				esModulesGatherAndCreateModuleMap(compiler, externs, js, false)
				esModulesRewriteModules(compiler, externs, js)
			})
		case "es6ToCjs":
			d.pass(rel, src, parseDiags, section, esModulesNoSetup, func(compiler *Compiler, externs, js *rhino.Node, extra *[]string) {
				esModulesGatherAndCreateModuleMap(compiler, externs, js, false)
				NewEs6RewriteModulesToCommonJsModules(compiler).Process(externs, js)
			})
		case "relativize":
			d.pass(rel, src, parseDiags, section, esModulesNoSetup, func(compiler *Compiler, externs, js *rhino.Node, extra *[]string) {
				NewEs6RelativizeImportPaths(compiler).Process(externs, js)
			})
		case "dynamicImports":
			d.pass(rel, src, parseDiags, section, func(o *CompilerOptions) { o.SetAllowDynamicImport(true) }, func(compiler *Compiler, externs, js *rhino.Node, extra *[]string) {
				esModulesGatherAndCreateModuleMap(compiler, externs, js, false)
				esModulesRewriteModules(compiler, externs, js)
				NewRewriteDynamicImports(compiler, nil, CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).Process(externs, js)
			})
		case "ijs":
			d.pass(rel, src, parseDiags, section, esModulesNoSetup, func(compiler *Compiler, externs, js *rhino.Node, extra *[]string) {
				NewConvertToTypedInterface(compiler).Process(externs, js)
			})
		default:
			panic("unknown section " + section)
		}
	}
}

// esModulesSplitFiles splits a dump into its "## file" sections.
func esModulesSplitFiles(dump string) map[string]string {
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
		if strings.HasPrefix(l, "## file ") {
			flush()
			cur = strings.TrimSuffix(strings.TrimPrefix(l, "## file "), "\n")
		}
		sb.WriteString(l)
	}
	flush()
	return out
}

// esModulesKeepSections drops from a file section the "## pass" sections not in keep.
func esModulesKeepSections(file string, keep []string) string {
	var sb strings.Builder
	in := true
	for _, l := range strings.SplitAfter(file, "\n") {
		if strings.HasPrefix(l, "## pass ") {
			name := strings.TrimSuffix(strings.TrimPrefix(l, "## pass "), "\n")
			in = false
			for _, k := range keep {
				if k == name {
					in = true
				}
			}
		}
		if in {
			sb.WriteString(l)
		}
	}
	return sb.String()
}

func esModulesReadGz(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func esModulesReadList(t *testing.T, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		if l := s.Text(); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// esModulesCompare dumps every file and compares it with its Java section.
func esModulesCompare(t *testing.T, root string, files []string, javaDump string, full bool, sections []string, maxSize int64, diffDir string) {
	t.Helper()
	java := esModulesSplitFiles(javaDump)
	mismatches := 0
	for _, rel := range files {
		srcBytes, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		d := &esModulesDump{sb: &strings.Builder{}, full: full, sections: sections}
		want := esModulesKeepSections(java[rel], sections)
		if int64(len(srcBytes)) > maxSize {
			esModulesLine(d.sb, "## file "+rel)
			esModulesLine(d.sb, "# skipped size "+strconv.Itoa(len(srcBytes)))
			if got := d.sb.String(); got != want {
				t.Errorf("%s: skipped section differs:\n  java: %q\n  go  : %q", rel, want, got)
			}
			continue
		}
		// Go-only: test-helper — Files.readString(UTF_8) is strict; the corpus is valid UTF-8.
		func() {
			defer func() {
				if r := recover(); r != nil {
					d.sb.WriteString(fmt.Sprintf("GO PANIC: %v\n", r))
				}
			}()
			d.dumpFile(rel, string(srcBytes))
		}()
		got := d.sb.String()
		if got != want {
			mismatches++
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
						t.Errorf("%s: line %d differs:\n  java: %s\n  go  : %s", rel, k+1, w, g)
						break
					}
				}
			}
			if diffDir != "" {
				p := filepath.Join(diffDir, rel)
				_ = os.MkdirAll(filepath.Dir(p), 0o755)
				_ = os.WriteFile(p+".go.txt", []byte(got), 0o644)
				_ = os.WriteFile(p+".java.txt", []byte(want), 0o644)
			}
		}
	}
	if mismatches > 0 {
		t.Errorf("%d of %d files differ", mismatches, len(files))
	}
}

func esModulesSections() []string {
	if s := os.Getenv("ES_MODULES_SECTIONS"); s != "" {
		return strings.Split(s, ",")
	}
	return esModulesAllSections
}

func TestEsModulesDumpGolden(t *testing.T) {
	dir := filepath.Join("testdata", "es_modules")
	files := esModulesReadList(t, filepath.Join(dir, "files.txt"))
	esModulesCompare(t, filepath.Join(dir, "src"), files, esModulesReadGz(t, filepath.Join(dir, "golden.txt.gz")), true, esModulesSections(), math.MaxInt64, os.Getenv("ES_MODULES_DIFF_DIR"))
}

func TestEsModulesDumpFixtures(t *testing.T) {
	dir := filepath.Join("testdata", "es_modules")
	files := esModulesReadList(t, filepath.Join(dir, "fixtures-files.txt"))
	esModulesCompare(t, filepath.Join("..", "..", "testdata", "fixtures"), files, esModulesReadGz(t, filepath.Join(dir, "fixtures.txt.gz")), false, esModulesSections(), 2000000, os.Getenv("ES_MODULES_DIFF_DIR"))
}

func TestEsModulesDumpCorpus(t *testing.T) {
	oracle := os.Getenv("ES_MODULES_ORACLE")
	if oracle == "" {
		t.Skip("set ES_MODULES_ORACLE, ES_MODULES_ROOT and ES_MODULES_FILES")
	}
	maxSize := int64(math.MaxInt64)
	if s := os.Getenv("ES_MODULES_MAX_SIZE"); s != "" {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		maxSize = n
	}
	files := esModulesReadList(t, os.Getenv("ES_MODULES_FILES"))
	esModulesCompare(t, os.Getenv("ES_MODULES_ROOT"), files, esModulesReadGz(t, oracle), os.Getenv("ES_MODULES_FULL") == "1", esModulesSections(), maxSize, os.Getenv("ES_MODULES_DIFF_DIR"))
}
