package jscomp

// Go-only: test-helper — the Go side of tools/oracle/java/ClosureModulesPassDump.java, the
// differential oracle for the l11-closure-modules passes. It runs the same passes the same way
// (a fresh Compiler per pass, Compiler.init + Compiler.parse, GatherModuleMetadata first where
// DefaultPassConfig runs it) and prints the same text, which must match the Java output byte for
// byte.
//
// The checked-in golden set: testdata/closure_modules/src/<rel> (small Closure Library and module
// files) with testdata/closure_modules/golden.txt.gz, produced by
//
//	java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.ClosureModulesPassDump \
//	  --full internal/jscomp/testdata/closure_modules/golden.txt.gz \
//	  internal/jscomp/testdata/closure_modules/src --list=internal/jscomp/testdata/closure_modules/files.txt
//
// The corpus run (env-gated): CLOSURE_MODULES_ORACLE=<the Java .gz> CLOSURE_MODULES_ROOT=<corpus
// root> CLOSURE_MODULES_FILES=<list file> [CLOSURE_MODULES_SAMPLE=<N>, as --sample]
// [CLOSURE_MODULES_FULL=1, as --full] [CLOSURE_MODULES_DIFF_DIR=<dir> to write the differing
// file sections].

import (
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

type closureModulesDump struct {
	sb     *strings.Builder
	full   bool
	sample int
}

// closureModulesRun is the Run interface of the Java tool.
type closureModulesRun func(compiler *Compiler, externs rhino.Node, js rhino.Node, extra *[]string)

func closureModulesEscape(s string) string {
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

func closureModulesLine(sb *strings.Builder, s string) {
	sb.WriteString(closureModulesEscape(s))
	sb.WriteByte('\n')
}

func closureModulesSha256(ascii string) string {
	d := sha256.Sum256([]byte(ascii))
	return fmt.Sprintf("%x", d[:8])
}

func closureModulesNewCompiler(rel string, src string) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	compiler.Init([]*SourceFile{}, []*SourceFile{SourceFile_FromCode(rel, src)}, options)
	compiler.Parse()
	return compiler
}

func closureModulesDiag(level string, e *JSError) string {
	return "diag " + level + " " + e.GetType().key + " " + strconv.Itoa(e.GetLineno()) + ":" +
		strconv.Itoa(e.GetCharno()) + " \"" + e.GetDescription() + "\""
}

func closureModulesDiagnostics(compiler *Compiler) []string {
	var out []string
	for _, e := range compiler.GetErrors() {
		out = append(out, closureModulesDiag("ERROR", e))
	}
	for _, e := range compiler.GetWarnings() {
		out = append(out, closureModulesDiag("WARNING", e))
	}
	return out
}

func closureModulesNodeRef(n rhino.Node) string {
	if n.IsNil() {
		return "null"
	}
	return n.GetToken().String() + "@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

func closureModulesGatherModuleMetadata(compiler *Compiler, externs rhino.Node, js rhino.Node) {
	NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(externs, js)
}

// closureModulesRenameCallback collects the rename table (the anonymous AbstractScopedCallback of
// the Java tool's renameTable).
type closureModulesRenameCallback struct {
	NodeTraversal_AbstractScopedCallbackBase

	table           *jutil.HashBasedTable[rhino.Node, string, string]
	moduleRootsOnly bool
}

func (c *closureModulesRenameCallback) EnterScope(t *NodeTraversal) {
	s := t.GetScope()
	r := s.GetRootNode()
	if c.moduleRootsOnly && !(r.IsModuleBody() || r.IsScript() || NodeUtil_IsFunctionBlock(r)) {
		return
	}
	for v := range s.GetVarIterable() {
		c.table.Put(s.GetRootNode(), v.GetName(), v.GetName()+"$r")
	}
}

func (c *closureModulesRenameCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {}

func closureModulesRenameTable(compiler *Compiler, js rhino.Node, moduleRootsOnly bool) *jutil.HashBasedTable[rhino.Node, string, string] {
	table := jutil.HashBasedTable_Create[rhino.Node, string, string]()
	NodeTraversal_Traverse(compiler, js, &closureModulesRenameCallback{table: table, moduleRootsOnly: moduleRootsOnly})
	return table
}

func (d *closureModulesDump) pass(rel string, src string, parseDiags []string, name string, run closureModulesRun) {
	sb := d.sb
	closureModulesLine(sb, "## pass "+name)
	compiler := closureModulesNewCompiler(rel, src)
	externs := compiler.GetExternsRoot()
	js := compiler.GetJsRoot()
	before := js.ToStringTree()
	var extra []string
	var exception *string
	jutil.Try(func() {
		run(compiler, externs, js, &extra)
	}, func(e jutil.IRuntimeException) {
		s := "exception " + e.Error()
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
	after := js.ToStringTree()
	if after == before {
		closureModulesLine(sb, "tree unchanged")
	} else {
		var tree strings.Builder
		for _, l := range strings.Split(after, "\n") {
			closureModulesLine(&tree, l)
		}
		closureModulesLine(sb, "tree changed "+closureModulesSha256(tree.String()))
		if d.full {
			sb.WriteString(tree.String())
		}
	}
	for _, e := range extra {
		closureModulesLine(sb, e)
	}
}

func (d *closureModulesDump) dumpFile(rel string, src string, sampled bool) {
	sb := d.sb
	closureModulesLine(sb, "## file "+rel)
	parsed := closureModulesNewCompiler(rel, src)
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

	d.pass(rel, src, parseDiags, "closureCheckModule", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		closureModulesGatherModuleMetadata(compiler, externs, js)
		NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap()).Process(externs, js)
	})
	d.pass(rel, src, parseDiags, "processClosurePrimitives", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		closureModulesGatherModuleMetadata(compiler, externs, js)
		NewProcessClosurePrimitives(compiler).Process(externs, js)
	})
	d.pass(rel, src, parseDiags, "collectProvidedNames", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		p := NewProcessClosureProvidesAndRequires(compiler, true)
		names := p.CollectProvidedNames(externs, js)
		for key, pn := range names.All() {
			*extra = append(*extra, "provided "+key+" "+pn.String()+
				" ns="+pn.GetNamespace()+
				" implicitInit="+strconv.FormatBool(pn.HasImplicitInitialization())+
				" legacy="+strconv.FormatBool(pn.IsFromLegacyModule())+
				" first="+closureModulesNodeRef(pn.GetFirstProvideCall())+
				" candidate="+closureModulesNodeRef(pn.GetCandidateDefinition()))
		}
		*extra = append(*extra, "exported "+p.GetExportedVariableNames().String())
	})
	d.pass(rel, src, parseDiags, "processClosureProvidesAndRequires", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		p := NewProcessClosureProvidesAndRequires(compiler, false)
		p.Process(externs, js)
		*extra = append(*extra, "exported "+p.GetExportedVariableNames().String())
	})
	d.pass(rel, src, parseDiags, "closureOptimizePrimitives", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		NewClosureOptimizePrimitives(compiler, true).Process(externs, js)
	})
	d.pass(rel, src, parseDiags, "closureOptimizePrimitivesEs5", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		NewClosureOptimizePrimitives(compiler, false).Process(externs, js)
	})
	d.pass(rel, src, parseDiags, "extraRequireRemover", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		NewExtraRequireRemover(compiler).Process(externs, js)
	})
	d.pass(rel, src, parseDiags, "es6CheckModule", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		NewEs6CheckModule(compiler).Process(externs, js)
	})
	d.pass(rel, src, parseDiags, "processTweaks", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		NewProcessTweaks(compiler, true).Process(externs, js)
	})
	d.pass(rel, src, parseDiags, "exportTestFunctions", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		exportSymbol := "goog.exportSymbol"
		exportProperty := "goog.exportProperty"
		NewExportTestFunctions(compiler, &exportSymbol, &exportProperty).Process(externs, js)
	})
	for _, allowLocalExports := range []bool{false, true} {
		d.pass(rel, src, parseDiags, "findExportableNodes local="+strconv.FormatBool(allowLocalExports), func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
			f := NewFindExportableNodes(compiler, allowLocalExports)
			NodeTraversal_Traverse(compiler, js, f)
			for k, v := range f.GetExports().All() {
				*extra = append(*extra, "export "+k+" "+closureModulesNodeRef(v))
			}
			for s := range f.GetLocalExports().All() {
				*extra = append(*extra, "local "+s)
			}
			for k, v := range f.GetEs6ClassExports().All() {
				*extra = append(*extra, "es6class "+closureModulesNodeRef(k)+" "+v)
			}
		})
	}
	d.pass(rel, src, parseDiags, "gatherRawExports", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
		g := NewGatherRawExports(compiler)
		g.Process(externs, js)
		*extra = append(*extra, "exported "+g.GetExportedVariableNames().String())
	})
	d.pass(rel, src, parseDiags, "whitespaceWrapGoogModules", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		NewWhitespaceWrapGoogModules(compiler).Process(externs, js)
	})
	if sampled {
		d.pass(rel, src, parseDiags, "convertChunksToESModules", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
			NewConvertChunksToESModules(compiler).Process(externs, js)
		})
	}
	d.pass(rel, src, parseDiags, "es6RenameReferences", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		NodeTraversal_Traverse(compiler, js, NewEs6RenameReferences(closureModulesRenameTable(compiler, js, false)))
	})
	d.pass(rel, src, parseDiags, "es6RenameTypeReferences", func(compiler *Compiler, externs, js rhino.Node, extra *[]string) {
		NodeTraversal_Traverse(compiler, js, NewEs6RenameTypeReferences(closureModulesRenameTable(compiler, js, true)))
	})
	d.conventions(parsed.GetJsRoot())
	d.deps(rel, src)
}

func (d *closureModulesDump) conventions(js rhino.Node) {
	sb := d.sb
	closureModulesLine(sb, "## conventions")
	c := CodingConventions_GetDefault()
	for script := js.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		closureModulesLine(sb, "package "+jutil.ToJavaString(c.GetPackageName(script.GetStaticSourceFile())))
	}
	NodeUtil_VisitPreOrder(js, NodeUtil_VisitorFunc(func(n rhino.Node) { d.conventionsAt(c, n) }))
}

func (d *closureModulesDump) conventionsAt(c CodingConvention, n rhino.Node) {
	sb := d.sb
	at := closureModulesNodeRef(n)
	if n.IsCall() {
		rel := c.GetClassesDefinedByCall(n)
		if rel != nil {
			closureModulesLine(sb, at+" subclass "+rel.subclassName+" "+rel.superclassName+" "+rel.type_.String())
		}
		for _, checkTypes := range []bool{false, true} {
			bind := c.DescribeFunctionBind(n, checkTypes)
			if bind != nil {
				tf := "F"
				if checkTypes {
					tf = "T"
				}
				closureModulesLine(sb, at+" bind"+tf+" "+closureModulesNodeRef(bind.target)+" "+
					closureModulesNodeRef(bind.thisValue)+" "+closureModulesNodeRef(bind.parameters)+" "+
					strconv.Itoa(bind.GetBoundParameterCount()))
			}
		}
		if c.IsPropertyTestFunction(n) {
			closureModulesLine(sb, at+" propertyTest")
		}
		if c.IsPropertyRenameFunction(n.GetFirstChild()) {
			closureModulesLine(sb, at+" propertyRename")
		}
		if ModuleImportResolver_IsGoogModuleDependencyCall(n) {
			closureModulesLine(sb, at+" googModuleDependencyCall")
		}
		if ConvertChunksToESModules_IsDynamicImportCallback(n) {
			closureModulesLine(sb, at+" dynamicImportCallback")
		}
	}
	if n.IsExprResult() || n.IsCall() {
		if c.IsFunctionCallThatAlwaysThrows(n) {
			closureModulesLine(sb, at+" alwaysThrows")
		}
		if CodingConventions_DefaultIsFunctionCallThatAlwaysThrows(n, "goog.asserts.fail") {
			closureModulesLine(sb, at+" alwaysThrows goog.asserts.fail")
		}
	}
	if n.IsName() && n.GetString() != "" {
		name := n.GetString()
		if c.IsExported(name) {
			closureModulesLine(sb, at+" exported "+name)
		}
		if c.IsExportedWithLocal(name, false) {
			closureModulesLine(sb, at+" exportedGlobal "+name)
		}
	}
	if n.IsParamList() {
		for p := n.GetFirstChild(); !p.IsNil(); p = p.GetNext() {
			if c.IsVarArgsParameter(p) {
				closureModulesLine(sb, closureModulesNodeRef(p)+" varArgs")
			}
		}
	}
	if NodeUtil_IsNameDeclaration(n) {
		for _, r := range ClosureRequireProcessor_GetAllRequires(n) {
			closureModulesLine(sb, at+" require "+r.LocalName()+" "+
				r.ImportRecord().ModuleRequest()+" "+
				r.ImportRecord().ImportName()+" "+
				r.ImportRecord().LocalName()+" "+
				closureModulesNodeRef(r.ImportRecord().ImportNode())+" "+
				closureModulesNodeRef(r.ImportRecord().NameNode())+" "+
				r.CreatedBy().String())
		}
	}
}

func (d *closureModulesDump) deps(rel string, src string) {
	sb := d.sb
	closureModulesLine(sb, "## deps")
	errorManager := NewPrintStreamErrorManager(jutil.NewPrintStream(io.Discard))
	parser := NewDepsFileRegexParser(errorManager)
	parser.SetShortcutMode(false)
	infos := parser.ParseFileWithContents(rel, src)
	for _, info := range infos {
		closureModulesLine(sb, "dep "+info.(*SimpleDependencyInfo).String())
	}
	for _, e := range errorManager.GetErrors() {
		closureModulesLine(sb, closureModulesDiag("ERROR", e))
	}
	for _, e := range errorManager.GetWarnings() {
		closureModulesLine(sb, closureModulesDiag("WARNING", e))
	}
	closureModulesLine(sb, "succeeded "+strconv.FormatBool(parser.DidParseSucceed()))
}

// closureModulesSplitFiles splits a dump into its "## file" sections.
func closureModulesSplitFiles(dump string) map[string]string {
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

func closureModulesReadGz(t *testing.T, path string) string {
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

// closureModulesReadList reads a file list, without the files of
// tools/oracle/gpl-excluded-corpus.txt.
func closureModulesReadList(t *testing.T, path string) []string {
	t.Helper()
	files, err := corpus.ReadList(path)
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// closureModulesCompare dumps every file and compares it with its Java section.
func closureModulesCompare(t *testing.T, root string, files []string, javaDump string, full bool, sample int, diffDir string) {
	t.Helper()
	java := closureModulesSplitFiles(javaDump)
	mismatches := 0
	for i, rel := range files {
		srcBytes, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		d := &closureModulesDump{sb: &strings.Builder{}, full: full, sample: sample}
		// Go-only: test-helper — Files.readString(UTF_8) is strict; the corpus is valid UTF-8.
		func() {
			defer func() {
				if r := recover(); r != nil {
					d.sb.WriteString(fmt.Sprintf("GO PANIC: %v\n", r))
				}
			}()
			d.dumpFile(rel, string(srcBytes), i%sample == 0)
		}()
		got := d.sb.String()
		want := java[rel]
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

func TestClosureModulesPassDumpGolden(t *testing.T) {
	dir := filepath.Join("testdata", "closure_modules")
	files := closureModulesReadList(t, filepath.Join(dir, "files.txt"))
	closureModulesCompare(t, filepath.Join(dir, "src"), files, closureModulesReadGz(t, filepath.Join(dir, "golden.txt.gz")), true, 1, os.Getenv("CLOSURE_MODULES_DIFF_DIR"))
}

func TestClosureModulesPassDumpCorpus(t *testing.T) {
	oracle := os.Getenv("CLOSURE_MODULES_ORACLE")
	if oracle == "" {
		t.Skip("set CLOSURE_MODULES_ORACLE, CLOSURE_MODULES_ROOT and CLOSURE_MODULES_FILES")
	}
	sample := 1
	if s := os.Getenv("CLOSURE_MODULES_SAMPLE"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			t.Fatal(err)
		}
		sample = n
	}
	files := closureModulesReadList(t, os.Getenv("CLOSURE_MODULES_FILES"))
	closureModulesCompare(t, os.Getenv("CLOSURE_MODULES_ROOT"), files, closureModulesReadGz(t, oracle), os.Getenv("CLOSURE_MODULES_FULL") == "1", sample, os.Getenv("CLOSURE_MODULES_DIFF_DIR"))
}
